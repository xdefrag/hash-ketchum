package usecase

import (
	"context"
	"errors"

	"github.com/xdefrag/hash-ketchum/pkg/types"
)

// Hash contains use cases for hashes.
type Hash struct {
	storage types.HashStorager
}

// ErrNoLeadingZeroes returns from Store when hash has no leading zeroes.
var ErrNoLeadingZeroes = errors.New("Hash must contains leading zeroes otherwise it invalid")

// NewHash create Hash usecase with dependencies.
func NewHash(storage types.HashStorager) *Hash {
	return &Hash{storage}
}

// Store checks hash for two leading zeroes and saves it
// to store, otherwise returns errNoLeadingZeroes error
func (h Hash) Store(ctx context.Context, hash types.Hash) error {

	if "00" != hash.Hash[:2] {
		return ErrNoLeadingZeroes
	}

	return h.storage.Store(ctx, hash)
}
