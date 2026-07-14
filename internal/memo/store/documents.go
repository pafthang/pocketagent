package store

import (
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/philippgille/chromem-go"
)

// DocumentView is a memo document without embedding vectors.
type DocumentView struct {
	ID       string            `json:"id"`
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// CollectionStats summarizes a space collection.
type CollectionStats struct {
	Collection    string `json:"collection"`
	DocumentCount int    `json:"document_count"`
	ContentBytes  int    `json:"content_bytes"`
}

func (m *Manager) GetDocument(ctx context.Context, spaceID, id string) (DocumentView, error) {
	if strings.TrimSpace(id) == "" {
		return DocumentView{}, fmt.Errorf("id is required")
	}

	collection, err := m.Collection(spaceID)
	if err != nil {
		return DocumentView{}, err
	}

	doc, err := collection.GetByID(ctx, id)
	if err != nil {
		return DocumentView{}, err
	}

	return toDocumentView(doc), nil
}

func (m *Manager) DeleteDocument(ctx context.Context, spaceID, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("id is required")
	}

	collection, err := m.Collection(spaceID)
	if err != nil {
		return err
	}

	return collection.Delete(ctx, nil, nil, id)
}

func (m *Manager) DeleteDocuments(ctx context.Context, spaceID string, ids []string) (int, error) {
	if len(ids) == 0 {
		return 0, fmt.Errorf("ids are required")
	}

	collection, err := m.Collection(spaceID)
	if err != nil {
		return 0, err
	}

	deleted := 0
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if err := collection.Delete(ctx, nil, nil, id); err != nil {
			return deleted, err
		}
		deleted++
	}
	return deleted, nil
}

func (m *Manager) ListDocuments(ctx context.Context, spaceID string, page, perPage int) ([]DocumentView, int, error) {
	collection, err := m.Collection(spaceID)
	if err != nil {
		return nil, 0, err
	}

	total := collection.Count()
	if total == 0 {
		return nil, 0, nil
	}

	docs, err := m.loadCollectionDocuments(collection.Name)
	if err != nil {
		return nil, 0, err
	}
	if len(docs) == 0 {
		return nil, 0, nil
	}

	sort.Slice(docs, func(i, j int) bool {
		return docs[i].ID < docs[j].ID
	})

	if perPage <= 0 {
		perPage = 50
	}
	if page <= 0 {
		page = 1
	}

	total = len(docs)
	start := (page - 1) * perPage
	if start >= total {
		return []DocumentView{}, total, nil
	}
	end := start + perPage
	if end > total {
		end = total
	}

	out := make([]DocumentView, 0, end-start)
	for _, doc := range docs[start:end] {
		out = append(out, toDocumentView(doc))
	}
	return out, total, nil
}

func (m *Manager) Stats(ctx context.Context, spaceID string) (CollectionStats, error) {
	collection, err := m.Collection(spaceID)
	if err != nil {
		return CollectionStats{}, err
	}

	stats := CollectionStats{
		Collection:    collection.Name,
		DocumentCount: collection.Count(),
	}

	docs, err := m.loadCollectionDocuments(collection.Name)
	if err != nil {
		return stats, err
	}
	for _, doc := range docs {
		stats.ContentBytes += len(doc.Content)
	}
	if stats.DocumentCount == 0 {
		stats.DocumentCount = len(docs)
	}
	return stats, nil
}

func (m *Manager) loadCollectionDocuments(collectionName string) ([]chromem.Document, error) {
	dir := filepath.Join(m.dataDir, collectionDirHash(collectionName))
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read collection dir: %w", err)
	}

	ext := ".gob"
	if m.compress {
		ext += ".gz"
	}

	docs := make([]chromem.Document, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if name == metadataFileName+ext {
			continue
		}
		if !strings.HasSuffix(name, ext) {
			continue
		}

		doc, err := readDocumentFile(filepath.Join(dir, name), m.compress)
		if err != nil {
			return nil, fmt.Errorf("read document %q: %w", name, err)
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func readDocumentFile(path string, compress bool) (chromem.Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return chromem.Document{}, err
	}
	defer f.Close()

	var r io.Reader = f
	if compress {
		gz, err := gzip.NewReader(f)
		if err != nil {
			return chromem.Document{}, err
		}
		defer gz.Close()
		r = gz
	}

	var doc chromem.Document
	if err := gob.NewDecoder(r).Decode(&doc); err != nil {
		return chromem.Document{}, err
	}
	return doc, nil
}

func collectionDirHash(name string) string {
	hash := sha256.Sum256([]byte(name))
	return hex.EncodeToString(hash[:4])
}

func toDocumentView(doc chromem.Document) DocumentView {
	meta := doc.Metadata
	if meta == nil {
		meta = map[string]string{}
	}
	return DocumentView{
		ID:       doc.ID,
		Content:  doc.Content,
		Metadata: meta,
	}
}