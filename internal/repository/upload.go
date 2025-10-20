package repository

import (
	"context"
	"github.com/google/uuid"
)

func (r *Repository) Update(ctx context.Context, id uuid.UUID, path, status string) error
