package storage

import "context"

// BackendType enumerates available storage backends.
type BackendType string

const (
	BackendMemory BackendType = "memory"
	BackendBolt   BackendType = "bolt"
	BackendBadger BackendType = "badger"
)

// Record represents an audit log entry saved to embedded storage.
type Record struct {
	Key   string
	Value []byte
}

// Store defines the persistence behaviour needed by the Governor. Backends
// must be safe for concurrent usage.
type Store interface {
	Put(ctx context.Context, record Record) error
	Get(ctx context.Context, key string) (Record, error)
	Iter(ctx context.Context, fn func(Record) error) error
	Delete(ctx context.Context, key string) error
	Close() error
}
