package path

import (
	"bytes"
	"io"
	"testing"
)

func TestReadTextBlobRejectsBinary(t *testing.T) {
	_, err := ReadTextBlob(func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader([]byte("data"))), nil
	}, "application/pdf", MaxIngestBytes)
	if err == nil {
		t.Fatal("expected error for binary mime")
	}
}

func TestReadTextBlobReadsMarkdown(t *testing.T) {
	content, err := ReadTextBlob(func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader([]byte("# hello"))), nil
	}, "text/markdown", MaxIngestBytes)
	if err != nil || content != "# hello" {
		t.Fatalf("unexpected: %q %v", content, err)
	}
}