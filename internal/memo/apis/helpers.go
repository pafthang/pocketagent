package memoapis

import (
	"strings"

	memoclient "github.com/pafthang/pocketagent/internal/memo/client"
)

func toMemoryDocument(doc memoclient.DocumentRecord) memoryDocument {
	meta := doc.Metadata
	if meta == nil {
		meta = map[string]string{}
	}
	return memoryDocument{
		ID:        doc.ID,
		Content:   doc.Content,
		Metadata:  meta,
		CreatedAt: meta["created_at"],
		Tags:      parseTags(meta["tags"]),
	}
}

func parseTags(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func copyMetadata(in map[string]string) map[string]string {
	if len(in) == 0 {
		return map[string]string{}
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}