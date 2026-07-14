package blob

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// Backend persists file bytes (local disk or object storage).
type Backend interface {
	Save(spaceID, fileID string, reader io.Reader) (storageKey string, checksum string, size int64, err error)
	Open(storageKey string) (io.ReadCloser, error)
	Delete(storageKey string) error
}

// StoreConfig selects a blob backend and its connection settings.
type StoreConfig struct {
	Backend string

	DataDir string

	S3Endpoint  string
	S3Bucket    string
	S3AccessKey string
	S3SecretKey string
	S3UseSSL    bool
	S3Region    string
}

// NewBackend creates a blob store from config.
func NewBackend(cfg StoreConfig) (Backend, error) {
	backend := strings.ToLower(strings.TrimSpace(cfg.Backend))
	if backend == "" {
		backend = "local"
	}
	switch backend {
	case "local":
		return NewBlobStore(cfg.DataDir)
	case "s3":
		return NewS3Backend(cfg)
	default:
		return nil, fmt.Errorf("unsupported files backend: %q", cfg.Backend)
	}
}

func blobStorageKey(spaceID, fileID string) (string, error) {
	spaceID = safeSegment(spaceID)
	fileID = safeSegment(fileID)
	if spaceID == "" || fileID == "" {
		return "", fmt.Errorf("space_id and file_id are required")
	}
	return filepath.ToSlash(filepath.Join("spaces", spaceID, "blobs", fileID)), nil
}

func validateStorageKey(storageKey string) error {
	storageKey = filepath.Clean(filepath.FromSlash(storageKey))
	if storageKey == "." || strings.HasPrefix(storageKey, "..") {
		return fmt.Errorf("invalid storage key")
	}
	return nil
}