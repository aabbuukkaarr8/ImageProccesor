package service

import (
	"context"
	"github.com/aabbuukkaarr8/internal/model"
	"github.com/google/uuid"
	"io"
)

func (s *Service) Get(ctx context.Context, id uuid.UUID) (model.Image, io.ReadCloser, error) {
	c := model.Image{}
	return c, nil, nil
}
