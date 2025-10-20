package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
)

func (s *StorageMinio) Delete(ctx context.Context, path string) error {
	return s.Client.RemoveObject(ctx, s.Bucket, path, minio.RemoveObjectOptions{})
}
