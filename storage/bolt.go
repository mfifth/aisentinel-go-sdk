package storage

import (
	"context"
	"errors"
	"sync"
)

// BoltStore simulates a BoltDB backed storage. In constrained environments the
// SDK falls back to an in-memory map while preserving the public API expected by
// the higher level components.
type BoltStore struct {
	mu     sync.RWMutex
	bucket map[string][]byte
}

// NewBolt creates a BoltStore backed by an in-memory map.
func NewBolt(_ string, _ any) (*BoltStore, error) {
	return &BoltStore{bucket: make(map[string][]byte)}, nil
}

// Put stores a record in the pseudo Bolt bucket.
func (s *BoltStore) Put(_ context.Context, record Record) error {
	s.mu.Lock()
	s.bucket[record.Key] = append([]byte(nil), record.Value...)
	s.mu.Unlock()
	return nil
}

// Get retrieves a record by key.
func (s *BoltStore) Get(_ context.Context, key string) (Record, error) {
	s.mu.RLock()
	value, ok := s.bucket[key]
	s.mu.RUnlock()
	if !ok {
		return Record{}, errRecordNotFound
	}
	return Record{Key: key, Value: append([]byte(nil), value...)}, nil
}

// Iter iterates over all stored records.
func (s *BoltStore) Iter(_ context.Context, fn func(Record) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, v := range s.bucket {
		if err := fn(Record{Key: k, Value: append([]byte(nil), v...)}); err != nil {
			return err
		}
	}
	return nil
}

// Delete removes a record by key.
func (s *BoltStore) Delete(_ context.Context, key string) error {
	s.mu.Lock()
	delete(s.bucket, key)
	s.mu.Unlock()
	return nil
}

// Close releases resources.
func (s *BoltStore) Close() error {
	if s.bucket == nil {
		return errors.New("bolt store closed")
	}
	s.mu.Lock()
	s.bucket = nil
	s.mu.Unlock()
	return nil
}
