package store

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/philippgille/chromem-go"
)

func TestManagerDocumentCRUD(t *testing.T) {
	dir := t.TempDir()
	mgr, err := Open(dir, "memory", false, 0.25)
	if err != nil {
		t.Fatalf("open manager: %v", err)
	}

	ctx := context.Background()
	spaceID := "space-a"
	collection, err := mgr.Collection(spaceID)
	if err != nil {
		t.Fatalf("collection: %v", err)
	}

	doc := chromem.Document{
		ID:        "doc-1",
		Content:   "hello memory",
		Embedding: []float32{1, 0, 0},
		Metadata:  map[string]string{"tags": "alpha,beta"},
	}
	if err := collection.AddDocument(ctx, doc); err != nil {
		t.Fatalf("add document: %v", err)
	}

	got, err := mgr.GetDocument(ctx, spaceID, "doc-1")
	if err != nil {
		t.Fatalf("get document: %v", err)
	}
	if got.Content != doc.Content {
		t.Fatalf("unexpected content: %q", got.Content)
	}

	items, total, err := mgr.ListDocuments(ctx, spaceID, 1, 10)
	if err != nil {
		t.Fatalf("list documents: %v", err)
	}
	if total != 1 || len(items) != 1 || items[0].ID != "doc-1" {
		t.Fatalf("unexpected list: total=%d items=%#v", total, items)
	}

	stats, err := mgr.Stats(ctx, spaceID)
	if err != nil {
		t.Fatalf("stats: %v", err)
	}
	if stats.DocumentCount != 1 || stats.ContentBytes != len(doc.Content) {
		t.Fatalf("unexpected stats: %#v", stats)
	}

	if err := mgr.DeleteDocument(ctx, spaceID, "doc-1"); err != nil {
		t.Fatalf("delete document: %v", err)
	}
	if collection.Count() != 0 {
		t.Fatalf("expected empty collection")
	}

	collectionDir := filepath.Join(dir, collectionDirHash(collection.Name))
	entries, err := os.ReadDir(collectionDir)
	if err != nil {
		t.Fatalf("read collection dir: %v", err)
	}
	for _, entry := range entries {
		if entry.Name() != metadataFileName+".gob" {
			t.Fatalf("expected metadata only, found %s", entry.Name())
		}
	}
}