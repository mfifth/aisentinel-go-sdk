package storage

import (
	"context"
	"errors"
	"sync"
)

var errRecordNotFound = errors.New("storage: record not found")

// MemoryStore is a concurrency safe in-memory store. It is primarily used for
// unit tests and ephemeral deployments.
type MemoryStore struct {
	mu     sync.RWMutex
	buffer map[string][]byte
}

// NewMemory creates a new MemoryStore instance.
func NewMemory() *MemoryStore {
	return &MemoryStore{buffer: make(map[string][]byte)}
}

// Put stores a record.
func (s *MemoryStore) Put(_ context.Context, record Record) error {
	s.mu.Lock()
	s.buffer[record.Key] = append([]byte(nil), record.Value...)
	s.mu.Unlock()
	return nil
}

// Get retrieves a record by key.
func (s *MemoryStore) Get(_ context.Context, key string) (Record, error) {
	s.mu.RLock()
	value, ok := s.buffer[key]
	s.mu.RUnlock()
	if !ok {
		return Record{}, errRecordNotFound
	}
	return Record{Key: key, Value: append([]byte(nil), value...)}, nil
}

// Iter iterates over all records.
func (s *MemoryStore) Iter(_ context.Context, fn func(Record) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, v := range s.buffer {
		if err := fn(Record{Key: k, Value: append([]byte(nil), v...)}); err != nil {
			return err
		}
	}
	return nil
}

// Delete removes a record by key.
func (s *MemoryStore) Delete(_ context.Context, key string) error {
	s.mu.Lock()
	delete(s.buffer, key)
	s.mu.Unlock()
	return nil
}

// Close releases resources. It is a no-op for the in-memory backend.
func (s *MemoryStore) Close() error { return nil }

// ErrNotFound exposes the not found error for external consumption.
func ErrNotFound() error { return errRecordNotFound }
