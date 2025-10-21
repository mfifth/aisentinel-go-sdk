package governor

import (
	"sync"
	"time"
)

// cacheEntry stores rulepack bytes and their expiration timestamp.
type cacheEntry[T any] struct {
	value     T
	expiresAt time.Time
}

// RuleCache provides a threadsafe TTL cache tailored to rulepacks. It uses a
// combination of generics and RWMutexes to provide contention free reads.
type RuleCache[T any] struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry[T]
	clock   func() time.Time
	ttl     time.Duration
}

// NewRuleCache constructs a RuleCache with the supplied TTL.
func NewRuleCache[T any](ttl time.Duration) *RuleCache[T] {
	return &RuleCache[T]{
		entries: make(map[string]cacheEntry[T]),
		clock:   time.Now,
		ttl:     ttl,
	}
}

// Get returns the cached value when it is still valid.
func (c *RuleCache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	entry, ok := c.entries[key]
	if ok && !c.expired(entry) {
		c.mu.RUnlock()
		return entry.value, true
	}
	c.mu.RUnlock()
	
	var zero T
	if ok {
		c.mu.Lock()
		// Re-check the entry after acquiring write lock to avoid race with Set
		if currentEntry, exists := c.entries[key]; exists && c.expired(currentEntry) {
			delete(c.entries, key)
		}
		c.mu.Unlock()
	}
	return zero, false
}

// Set stores a value with an optional per-value TTL.
func (c *RuleCache[T]) Set(key string, value T, ttlOverride ...time.Duration) {
	ttl := c.ttl
	if len(ttlOverride) > 0 {
		ttl = ttlOverride[0]
	}
	entry := cacheEntry[T]{
		value:     value,
		expiresAt: c.clock().Add(ttl),
	}
	c.mu.Lock()
	c.entries[key] = entry
	c.mu.Unlock()
}

// Invalidate removes an entry from the cache.
func (c *RuleCache[T]) Invalidate(key string) {
	c.mu.Lock()
	delete(c.entries, key)
	c.mu.Unlock()
}

func (c *RuleCache[T]) expired(entry cacheEntry[T]) bool {
	return !entry.expiresAt.IsZero() && c.clock().After(entry.expiresAt)
}

// Len returns the number of active entries. Mainly used for metrics.
func (c *RuleCache[T]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
