package service

import (
	"context"
	"fmt"
	"github.com/aabbuukkaarr8/internal/model"
	"github.com/google/uuid"
	"io"
)

func (s *Service) SaveImage(ctx context.Context, subdir, filename string, file io.Reader, action model.Action) (uuid.UUID, string, error) {
	dst, err := s.minio.Save(ctx, subdir, filename, file)
	if err != nil {
		return uuid.Nil, "", err
	}
	img := model.Image{
		Filename: filename,
		Path:     dst,
		Action:   action,
		Status:   "pending",
	}
	id, err := s.repo.SaveImage(ctx, img)
	if err != nil {
		return uuid.Nil, "", err
	}
	img.ID = id
	err = s.producer.Produce(ctx, img)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("failed to produce image: %w", err)
	}
	return id, dst, nil
}
