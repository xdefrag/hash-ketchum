package types

import "context"

// Hash domain type.
type Hash struct {
	Login     string
	Hash      string
	Timestamp int64
}

// HashStorager interface for storing hash object.
type HashStorager interface {
	Store(ctx context.Context, hash Hash) error
}
