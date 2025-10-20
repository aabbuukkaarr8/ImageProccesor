package service

import (
	"context"
	"github.com/aabbuukkaarr8/internal/model"
	"github.com/google/uuid"
)

func (s *Service) MachineForImg(ctx context.Context, image model.Image) (uuid.UUID, error)
