package blob

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalBlobStoreRoundTrip(t *testing.T) {
	dir := t.TempDir()
	store, err := NewBlobStore(dir)
	if err != nil {
		t.Fatalf("NewBlobStore: %v", err)
	}

	payload := []byte("hello pocketagent files")
	key, checksum, size, err := store.Save("space-1", "file-1", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if size != int64(len(payload)) || key == "" || checksum == "" {
		t.Fatalf("unexpected save result: key=%q checksum=%q size=%d", key, checksum, size)
	}

	reader, err := store.Open(key)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer reader.Close()

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("content mismatch: %q", got)
	}

	abs := filepath.Join(dir, filepath.FromSlash(key))
	if _, err := os.Stat(abs); err != nil {
		t.Fatalf("blob file missing: %v", err)
	}

	if err := store.Delete(key); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := os.Stat(abs); !os.IsNotExist(err) {
		t.Fatalf("expected blob removed, stat err=%v", err)
	}
}

func TestNewBackendDefaultsToLocal(t *testing.T) {
	backend, err := NewBackend(StoreConfig{DataDir: t.TempDir()})
	if err != nil {
		t.Fatalf("NewBackend: %v", err)
	}
	if _, ok := backend.(*BlobStore); !ok {
		t.Fatalf("expected *BlobStore, got %T", backend)
	}
}

func TestValidateStorageKeyRejectsTraversal(t *testing.T) {
	store, err := NewBlobStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewBlobStore: %v", err)
	}
	if _, err := store.Open("../secret"); err == nil {
		t.Fatal("expected traversal rejection")
	}
}