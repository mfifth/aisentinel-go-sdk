package governor

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestRuleCache(t *testing.T) {
	cache := NewRuleCache[int](time.Millisecond)
	cache.Set("a", 42)
	if value, ok := cache.Get("a"); !ok || value != 42 {
		t.Fatalf("unexpected cache read: %v %v", value, ok)
	}
	time.Sleep(2 * time.Millisecond)
	if _, ok := cache.Get("a"); ok {
		t.Fatal("expected cache entry to expire")
	}
}

func TestGovernorEvaluateOffline(t *testing.T) {
	cfg := Config{APIKey: "test", OfflineMode: true}
	ctx := context.Background()
	gov, err := NewGovernor(ctx, cfg)
	if err != nil {
		t.Fatalf("expected governor: %v", err)
	}
	t.Cleanup(func() { _ = gov.Close() })

	payload, _ := json.Marshal(map[string]string{"rule-1": "value"})
	_, err = gov.Evaluate(ctx, DecisionRequest{RulepackID: "local", Payload: payload})
	if err == nil {
		t.Fatal("expected error in offline mode with missing cache")
	}
}
