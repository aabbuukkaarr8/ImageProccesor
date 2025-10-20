package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
	"path/filepath"
)

func (s *StorageMinio) Save(ctx context.Context, subdir, filename string, src io.Reader) (string, error) {
	objectName := filepath.Join(subdir, filename)

	_, err := s.Client.PutObject(ctx, s.Bucket, objectName, src, -1, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return objectName, nil
}
