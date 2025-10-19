package repository

import (
	"github.com/aabbuukkaarr8/internal/storage"
)

type Repository struct {
	store *storage.Store
}

func NewRepository(store *storage.Store) *Repository {
	return &Repository{
		store: store,
	}
}
