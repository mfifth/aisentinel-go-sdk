package governor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/aisentinel/aisentinel-go-sdk/storage"
)

// ErrOffline indicates the Governor is operating in offline mode.
var ErrOffline = errors.New("governor: offline mode enabled")

// ErrRuleNotFound occurs when the requested rule is not found in the cache.
var ErrRuleNotFound = errors.New("governor: rule not found")

// DecisionRequest describes an authorization decision request.
type DecisionRequest struct {
	RulepackID string
	Payload    json.RawMessage
}

// DecisionResult represents the outcome of a decision evaluation.
type DecisionResult struct {
	Allowed bool
	Reason  string
	Latency time.Duration
}

// Option configures Governor construction.
type Option func(*Governor) error

// Governor coordinates configuration, caching, storage and evaluation.
type Governor struct {
	cfg         Config
	httpClient  *http.Client
	cache       *RuleCache[*Rulepack]
	evaluator   *Evaluator
	storage     storage.Store
	offline     bool
	offlineChan chan DecisionRequest
	mu          sync.RWMutex
}

// NewGovernor constructs a Governor instance using the provided configuration.
func NewGovernor(ctx context.Context, cfg Config, opts ...Option) (*Governor, error) {
	cfg = DefaultConfig().Merge(cfg)
	if err := cfg.ApplyEnv(); err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: cfg.HTTPTimeout,
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			DialContext:         (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
			IdleConnTimeout:     90 * time.Second,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
		},
	}

	cache := NewRuleCache[*Rulepack](cfg.CacheTTL)
	evaluator := NewEvaluator()

	store, err := buildStore(cfg)
	if err != nil {
		return nil, err
	}

	g := &Governor{
		cfg:         cfg,
		httpClient:  client,
		cache:       cache,
		evaluator:   evaluator,
		storage:     store,
		offline:     cfg.OfflineMode,
		offlineChan: make(chan DecisionRequest, cfg.OfflineQueueSize),
	}

	for _, opt := range opts {
		if err := opt(g); err != nil {
			return nil, err
		}
	}

	if g.offline {
		go g.drainOfflineQueue(ctx)
	}

	return g, nil
}

// buildStore creates a storage backend from configuration.
func buildStore(cfg Config) (storage.Store, error) {
	switch storage.BackendType(cfg.StorageBackend) {
	case storage.BackendBolt:
		if cfg.StorageDSN == "" {
			return nil, fmt.Errorf("bolt backend selected but StorageDSN empty")
		}
		return storage.NewBolt(cfg.StorageDSN, nil)
	case storage.BackendBadger:
		if cfg.StorageDSN == "" {
			return nil, fmt.Errorf("badger backend selected but StorageDSN empty")
		}
		return storage.NewBadger(cfg.StorageDSN, storage.DefaultBadgerOptions())
	default:
		return storage.NewMemory(), nil
	}
}

// WithHTTPClient overrides the default HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(g *Governor) error {
		if client == nil {
			return fmt.Errorf("http client cannot be nil")
		}
		g.httpClient = client
		return nil
	}
}

// WithStorage injects a preconstructed storage backend.
func WithStorage(store storage.Store) Option {
	return func(g *Governor) error {
		if store == nil {
			return fmt.Errorf("storage cannot be nil")
		}
		g.storage = store
		return nil
	}
}

// Rulepack holds compiled rule evaluation metadata.
type Rulepack struct {
	ID        string           `json:"id"`
	Version   string           `json:"version"`
	Rules     []RuleDefinition `json:"rules"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// Evaluate performs a governance decision against the current rulepack.
func (g *Governor) Evaluate(ctx context.Context, req DecisionRequest) (DecisionResult, error) {
	start := time.Now()
	pack, err := g.loadRulepack(ctx, req.RulepackID)
	if err != nil {
		return DecisionResult{}, err
	}

	allowed, reason, err := g.evaluator.Evaluate(ctx, pack, req.Payload)
	if err != nil {
		return DecisionResult{}, err
	}

	result := DecisionResult{Allowed: allowed, Reason: reason, Latency: time.Since(start)}
	_ = g.persistAudit(ctx, req, result)
	return result, nil
}

// loadRulepack retrieves a rulepack from cache or remote.
func (g *Governor) loadRulepack(ctx context.Context, id string) (*Rulepack, error) {
	if pack, ok := g.cache.Get(id); ok {
		return pack, nil
	}

	if g.offline {
		return nil, fmt.Errorf("%w: rulepack %s unavailable", ErrOffline, id)
	}

	pack, err := g.fetchRulepack(ctx, id)
	if err != nil {
		return nil, err
	}
	g.cache.Set(id, pack)
	return pack, nil
}

// fetchRulepack downloads the rulepack from the control plane. A minimal
// implementation is provided to keep the SDK functional in offline examples.
func (g *Governor) fetchRulepack(ctx context.Context, id string) (*Rulepack, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/rulepacks/%s", g.cfg.APIBaseURL, id), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+g.cfg.APIKey)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch rulepack: unexpected status %d", resp.StatusCode)
	}
	var pack Rulepack
	if err := json.NewDecoder(resp.Body).Decode(&pack); err != nil {
		return nil, err
	}
	return &pack, nil
}

func (g *Governor) persistAudit(ctx context.Context, req DecisionRequest, result DecisionResult) error {
	if g.storage == nil {
		return nil
	}
	record := storage.Record{
		Key: fmt.Sprintf("%s:%d", req.RulepackID, time.Now().UnixNano()),
		Value: mustJSON(map[string]any{
			"rulepack_id": req.RulepackID,
			"payload":     req.Payload,
			"allowed":     result.Allowed,
			"reason":      result.Reason,
			"latency_ms":  result.Latency.Milliseconds(),
		}),
	}
	return g.storage.Put(ctx, record)
}

func mustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		return []byte("{}")
	}
	return b
}

func (g *Governor) drainOfflineQueue(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-g.offlineChan:
			_, _ = g.Evaluate(context.Background(), req)
		}
	}
}

// Queue stores a request for later replay in offline mode.
func (g *Governor) Queue(req DecisionRequest) error {
	if !g.offline {
		return fmt.Errorf("queue requires offline mode")
	}
	select {
	case g.offlineChan <- req:
		return nil
	default:
		return fmt.Errorf("offline queue full")
	}
}

// Close releases resources used by the Governor.
func (g *Governor) Close() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.storage != nil {
		return g.storage.Close()
	}
	return nil
}

// WithOffline toggles offline mode after construction.
func (g *Governor) WithOffline(enabled bool) {
	g.mu.Lock()
	g.offline = enabled
	g.mu.Unlock()
}
