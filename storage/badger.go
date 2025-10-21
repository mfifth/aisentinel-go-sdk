package storage

import (
	"context"
	"errors"
	"sync"
)

// BadgerStore acts as an in-memory approximation of a BadgerDB store. The
// simplified implementation keeps the Go module dependency free for CI
// environments without network access while remaining API compatible.
type BadgerStore struct {
	mu   sync.RWMutex
	data map[string][]byte
}

// NewBadger creates an in-memory BadgerStore.
func NewBadger(_ string, _ any) (*BadgerStore, error) {
	return &BadgerStore{data: make(map[string][]byte)}, nil
}

// Put stores a record.
func (s *BadgerStore) Put(_ context.Context, record Record) error {
	s.mu.Lock()
	s.data[record.Key] = append([]byte(nil), record.Value...)
	s.mu.Unlock()
	return nil
}

// Get retrieves a record by key.
func (s *BadgerStore) Get(_ context.Context, key string) (Record, error) {
	s.mu.RLock()
	value, ok := s.data[key]
	s.mu.RUnlock()
	if !ok {
		return Record{}, errRecordNotFound
	}
	return Record{Key: key, Value: append([]byte(nil), value...)}, nil
}

// Iter iterates over all records.
func (s *BadgerStore) Iter(_ context.Context, fn func(Record) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, v := range s.data {
		if err := fn(Record{Key: k, Value: append([]byte(nil), v...)}); err != nil {
			return err
		}
	}
	return nil
}

// Delete removes a record.
func (s *BadgerStore) Delete(_ context.Context, key string) error {
	s.mu.Lock()
	delete(s.data, key)
	s.mu.Unlock()
	return nil
}

// Close releases resources.
func (s *BadgerStore) Close() error {
	if s.data == nil {
		return errors.New("badger store closed")
	}
	s.mu.Lock()
	s.data = nil
	s.mu.Unlock()
	return nil
}

// DefaultBadgerOptions returns a nil placeholder for API compatibility.
func DefaultBadgerOptions() any { return nil }
