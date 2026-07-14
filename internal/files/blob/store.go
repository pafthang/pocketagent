package blob

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// BlobStore persists file bytes on local disk.
type BlobStore struct {
	root string
}

// NewBlobStore creates a store rooted at dataDir.
func NewBlobStore(dataDir string) (*BlobStore, error) {
	if dataDir == "" {
		dataDir = "data/files"
	}
	root := filepath.Clean(dataDir)
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("create files data dir: %w", err)
	}
	return &BlobStore{root: root}, nil
}

// Save writes content for a space file and returns storage key and checksum.
func (s *BlobStore) Save(spaceID, fileID string, reader io.Reader) (storageKey string, checksum string, size int64, err error) {
	rel, err := blobStorageKey(spaceID, fileID)
	if err != nil {
		return "", "", 0, err
	}
	abs := filepath.Join(s.root, rel)
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return "", "", 0, err
	}

	f, err := os.Create(abs)
	if err != nil {
		return "", "", 0, err
	}
	defer f.Close()

	hasher := sha256.New()
	written, err := io.Copy(io.MultiWriter(f, hasher), reader)
	if err != nil {
		_ = os.Remove(abs)
		return "", "", 0, err
	}

	return rel, hex.EncodeToString(hasher.Sum(nil)), written, nil
}

// Open returns a reader for a stored blob.
func (s *BlobStore) Open(storageKey string) (io.ReadCloser, error) {
	abs, err := s.resolve(storageKey)
	if err != nil {
		return nil, err
	}
	return os.Open(abs)
}

// Delete removes a blob if present.
func (s *BlobStore) Delete(storageKey string) error {
	if strings.TrimSpace(storageKey) == "" {
		return nil
	}
	abs, err := s.resolve(storageKey)
	if err != nil {
		return err
	}
	if err := os.Remove(abs); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *BlobStore) resolve(storageKey string) (string, error) {
	if err := validateStorageKey(storageKey); err != nil {
		return "", err
	}
	storageKey = filepath.Clean(filepath.FromSlash(storageKey))
	abs := filepath.Join(s.root, storageKey)
	if !strings.HasPrefix(abs, s.root+string(os.PathSeparator)) && abs != s.root {
		return "", fmt.Errorf("invalid storage key")
	}
	return abs, nil
}

func safeSegment(v string) string {
	v = strings.TrimSpace(v)
	v = strings.ReplaceAll(v, "/", "_")
	v = strings.ReplaceAll(v, "\\", "_")
	return v
}