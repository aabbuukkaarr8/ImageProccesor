package service

import (
	"context"
	"github.com/aabbuukkaarr8/internal/model"
	"github.com/google/uuid"
	"io"
)

type StorageMinio interface {
	Save(ctx context.Context, subdir, filename string, src io.Reader) (string, error)
	Load(ctx context.Context, path string) (io.ReadCloser, error)
	Delete(ctx context.Context, path string) error
}

type Producer interface {
	Produce(ctx context.Context, img model.Image) error
}

type Machine interface {
	Machine(ctx context.Context, img model.Image) (model.Image, error)
}
type Repository interface {
	Save(ctx context.Context, img model.Image) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (model.Image, error)
	Update(ctx context.Context, id uuid.UUID, path, status string) error
	Delete(ctx context.Context, id uuid.UUID) error
}
