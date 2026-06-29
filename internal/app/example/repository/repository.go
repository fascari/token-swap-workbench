package repository

import (
	"context"

	"sync"

	"github.com/fascari/token-swap-workbench/internal/app/example/domain"
)

// Repository is an in-memory implementation of the example data store.
type Repository struct {
	mu      sync.Mutex
	nextID  uint
	storage map[uint]domain.Example
}

func New() *Repository {
	return &Repository{
		nextID:  1,
		storage: make(map[uint]domain.Example),
	}
}

func (r *Repository) Create(ctx context.Context, example domain.Example) (domain.Example, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	example.ID = r.nextID
	r.nextID++
	r.storage[example.ID] = example
	return example, nil
}
