package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/pafthang/pocketagent/internal/memo/chunk"
)

// Embedder generates vector embeddings for text (e.g. Ollama client).
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}

// Store embeds content and saves it to memo.
func (c *Client) Store(ctx context.Context, embedder Embedder, id, content string) error {
	return c.StoreScoped(ctx, embedder, "", id, content)
}

// StoreScoped stores content in a per-space chromem collection with optional chunking.
func (c *Client) StoreScoped(ctx context.Context, embedder Embedder, spaceID, id, content string) error {
	return c.StoreScopedWithMeta(ctx, embedder, spaceID, id, content, nil)
}

// StoreScopedWithMeta stores chunked content with metadata for better RAG recall.
func (c *Client) StoreScopedWithMeta(ctx context.Context, embedder Embedder, spaceID, id, content string, metadata map[string]string) error {
	chunks := chunk.Text(content, c.ChunkSize, c.ChunkOverlap)
	if len(chunks) == 0 {
		return nil
	}

	for i, chunk := range chunks {
		chunkID := id
		if len(chunks) > 1 {
			chunkID = fmt.Sprintf("%s#%d", id, i)
		}

		embedding, err := embedder.Embed(ctx, chunk)
		if err != nil {
			return err
		}

		meta := copyMetadata(metadata)
		meta["chunk_index"] = fmt.Sprintf("%d", i)
		meta["chunk_total"] = fmt.Sprintf("%d", len(chunks))
		if spaceID != "" {
			meta["space_id"] = spaceID
		}

		if err := c.AddScoped(ctx, spaceID, chunkID, chunk, embedding, meta); err != nil {
			return err
		}
	}
	return nil
}

// Query embeds the query and searches memo for similar documents.
func (c *Client) Query(ctx context.Context, embedder Embedder, query string, limit int) ([]string, error) {
	return c.QueryScoped(ctx, embedder, "", query, limit)
}

// QueryScoped searches a space collection and returns formatted RAG context lines.
func (c *Client) QueryScoped(ctx context.Context, embedder Embedder, spaceID, query string, limit int) ([]string, error) {
	embedding, err := embedder.Embed(ctx, query)
	if err != nil {
		return nil, err
	}

	docs, err := c.SearchDocumentsScoped(ctx, embedding, spaceID, limit)
	if err != nil {
		return nil, err
	}

	return formatRAGLines(docs), nil
}

func formatRAGLines(docs []Document) []string {
	lines := make([]string, 0, len(docs))
	seen := make(map[string]struct{}, len(docs))
	for i, doc := range docs {
		key := strings.TrimSpace(doc.Content)
		if key == "" {
			continue
		}
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		lines = append(lines, fmt.Sprintf("[%d] (%.2f) %s", i+1, doc.Similarity, doc.Content))
	}
	return lines
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
