package path

import (
	"mime"
	"path/filepath"
	"strings"
)

func init() {
	_ = mime.AddExtensionType(".md", "text/markdown")
}

// DetectMimeType guesses a MIME type from filename and optional content sniff.
func DetectMimeType(name string, sniff []byte) string {
	if len(sniff) > 0 {
		return mime.TypeByExtension(filepath.Ext(name))
	}
	ext := strings.ToLower(filepath.Ext(name))
	if ext == "" {
		return "application/octet-stream"
	}
	if t := mime.TypeByExtension(ext); t != "" {
		return t
	}
	return "application/octet-stream"
}

// IsTextMime reports whether content can be returned as text preview.
func IsTextMime(mimeType string) bool {
	return strings.HasPrefix(mimeType, "text/") ||
		mimeType == "application/json" ||
		mimeType == "application/markdown" ||
		strings.Contains(mimeType, "markdown")
}