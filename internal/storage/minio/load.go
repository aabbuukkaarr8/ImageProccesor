package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
)

func (s *StorageMinio) Load(ctx context.Context, path string) (io.ReadCloser, error) {
	obj, err := s.Client.GetObject(ctx, s.Bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to load file: %w", err)
	}

	return obj, nil
}
