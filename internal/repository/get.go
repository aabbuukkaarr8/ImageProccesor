package repository

import (
	"context"
	"github.com/aabbuukkaarr8/internal/model"
	"github.com/google/uuid"
)

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (model.Image, error)
