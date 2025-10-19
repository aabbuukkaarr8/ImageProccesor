package service

import (
	"github.com/aabbuukkaarr8/internal/kafka"
	"github.com/aabbuukkaarr8/internal/repository"
	"github.com/aabbuukkaarr8/internal/storage"
)

type Service struct {
	repo    *repository.Repository
	storage *storage.StorageMinio
	kafka   *kafka.Producer
}

func NewImageService(r *repository.Repository, s *storage.StorageMinio, k *kafka.Producer) *Service {
	return &Service{
		repo:    r,
		storage: s,
		kafka:   k,
	}
}
