package path

import (
	"fmt"
	"io"
)

// MaxIngestBytes is the largest blob ingested into memo RAG.
const MaxIngestBytes = 5 << 20

// ReadTextBlob reads a stored blob as UTF-8 text when MIME type is text-compatible.
func ReadTextBlob(open func() (io.ReadCloser, error), mimeType string, maxBytes int64) (string, error) {
	if maxBytes <= 0 {
		maxBytes = MaxIngestBytes
	}
	if !IsTextMime(mimeType) {
		return "", fmt.Errorf("binary file cannot be ingested")
	}

	rc, err := open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	body, err := io.ReadAll(io.LimitReader(rc, maxBytes+1))
	if err != nil {
		return "", err
	}
	if int64(len(body)) > maxBytes {
		return "", fmt.Errorf("file exceeds ingest size limit (%d bytes)", maxBytes)
	}
	return string(body), nil
}