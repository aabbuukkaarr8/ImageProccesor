package handler

import (
	"context"
	"github.com/aabbuukkaarr8/internal/model"
	"github.com/google/uuid"
	"io"
)

type Service interface {
	Get(ctx context.Context, id uuid.UUID) (model.Image, io.ReadCloser, error)
	MachineForImg(ctx context.Context, image model.Image) (uuid.UUID, error)
	SaveImage(ctx context.Context, subdir, filename string, file io.Reader, action model.Action) (uuid.UUID, string, error)
}
