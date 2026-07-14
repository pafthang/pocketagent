package blob

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3Backend persists file bytes in S3-compatible object storage (MinIO, AWS S3, …).
type S3Backend struct {
	client *minio.Client
	bucket string
}

// NewS3Backend creates an object-store backend.
func NewS3Backend(cfg StoreConfig) (*S3Backend, error) {
	endpoint := strings.TrimSpace(cfg.S3Endpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("files_s3_endpoint is required for s3 backend")
	}
	bucket := strings.TrimSpace(cfg.S3Bucket)
	if bucket == "" {
		return nil, fmt.Errorf("files_s3_bucket is required for s3 backend")
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		Secure: cfg.S3UseSSL,
		Region: strings.TrimSpace(cfg.S3Region),
	})
	if err != nil {
		return nil, fmt.Errorf("create s3 client: %w", err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("check s3 bucket: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("s3 bucket %q does not exist", bucket)
	}

	return &S3Backend{client: client, bucket: bucket}, nil
}

// Save writes content for a space file and returns storage key and checksum.
func (s *S3Backend) Save(spaceID, fileID string, reader io.Reader) (storageKey string, checksum string, size int64, err error) {
	key, err := blobStorageKey(spaceID, fileID)
	if err != nil {
		return "", "", 0, err
	}

	hasher := sha256.New()
	tee := io.TeeReader(reader, hasher)
	info, err := s.client.PutObject(context.Background(), s.bucket, key, tee, -1, minio.PutObjectOptions{})
	if err != nil {
		return "", "", 0, err
	}

	return key, hex.EncodeToString(hasher.Sum(nil)), info.Size, nil
}

// Open returns a reader for a stored blob.
func (s *S3Backend) Open(storageKey string) (io.ReadCloser, error) {
	if err := validateStorageKey(storageKey); err != nil {
		return nil, err
	}
	obj, err := s.client.GetObject(context.Background(), s.bucket, storageKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	if _, err := obj.Stat(); err != nil {
		_ = obj.Close()
		var resp minio.ErrorResponse
		if errors.As(err, &resp) && resp.Code == "NoSuchKey" {
			return nil, fmt.Errorf("blob not found")
		}
		return nil, err
	}
	return obj, nil
}

// Delete removes a blob if present.
func (s *S3Backend) Delete(storageKey string) error {
	if strings.TrimSpace(storageKey) == "" {
		return nil
	}
	if err := validateStorageKey(storageKey); err != nil {
		return err
	}
	err := s.client.RemoveObject(context.Background(), s.bucket, storageKey, minio.RemoveObjectOptions{})
	if err != nil {
		var resp minio.ErrorResponse
		if errors.As(err, &resp) && resp.Code == "NoSuchKey" {
			return nil
		}
		return err
	}
	return nil
}