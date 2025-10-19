package storage

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageMinio struct {
	Client *minio.Client
	Bucket string
}

func NewMinio(endpoint, accessKey, secretKey, bucket string) (*StorageMinio, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("cannot create bucket: %w", err)
		}
	}

	return &StorageMinio{Client: client, Bucket: bucket}, nil
}

func (s *StorageMinio) Upload(ctx context.Context, objectName string, filePath string) error {
	_, err := s.Client.FPutObject(ctx, s.Bucket, objectName, filePath, minio.PutObjectOptions{})
	return err
}

func (s *StorageMinio) Download(ctx context.Context, objectName, destPath string) error {
	return s.Client.FGetObject(ctx, s.Bucket, objectName, destPath, minio.GetObjectOptions{})
}

func (s *StorageMinio) Delete(ctx context.Context, objectName string) error {
	return s.Client.RemoveObject(ctx, s.Bucket, objectName, minio.RemoveObjectOptions{})
}
