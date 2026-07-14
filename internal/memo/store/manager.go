package store

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/philippgille/chromem-go"
)

var spaceIDPattern = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)

const metadataFileName = "00000000"

// Manager holds a persistent chromem DB with per-space collections.
type Manager struct {
	db                *chromem.DB
	dataDir           string
	compress          bool
	defaultCollection string
	embeddingFunc     chromem.EmbeddingFunc
	minSimilarity     float32
}

// Open creates a Manager backed by a persistent chromem store.
func Open(dataDir, defaultCollection string, compress bool, minSimilarity float32) (*Manager, error) {
	if dataDir == "" {
		dataDir = "data/memo"
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("create memo data dir: %w", err)
	}

	db, err := chromem.NewPersistentDB(filepath.Clean(dataDir), compress)
	if err != nil {
		return nil, fmt.Errorf("open persistent memo db: %w", err)
	}

	if defaultCollection == "" {
		defaultCollection = "memory"
	}
	if minSimilarity <= 0 {
		minSimilarity = 0.25
	}

	return &Manager{
		db:                db,
		dataDir:           filepath.Clean(dataDir),
		compress:          compress,
		defaultCollection: defaultCollection,
		embeddingFunc:     chromem.NewEmbeddingFuncDefault(),
		minSimilarity:     minSimilarity,
	}, nil
}

// Ping verifies the chromem persistent store is reachable.
func (m *Manager) Ping() error {
	if m == nil || m.db == nil {
		return fmt.Errorf("memo store not initialized")
	}
	_ = m.db.ListCollections()
	return nil
}

// MinSimilarity returns the default RAG similarity threshold.
func (m *Manager) MinSimilarity() float32 {
	if m == nil {
		return 0.25
	}
	return m.minSimilarity
}

func (m *Manager) Collection(spaceID string) (*chromem.Collection, error) {
	name := m.defaultCollection
	metadata := map[string]string{"scope": "global"}
	if spaceID != "" {
		name = spaceCollectionName(spaceID)
		metadata = map[string]string{"space_id": spaceID, "scope": "space"}
	}
	return m.db.GetOrCreateCollection(name, metadata, m.embeddingFunc)
}

func spaceCollectionName(spaceID string) string {
	safe := spaceIDPattern.ReplaceAllString(strings.TrimSpace(spaceID), "_")
	if safe == "" {
		safe = "unknown"
	}
	return "space_" + safe
}

// FilterBySimilarity keeps results above minSimilarity up to limit.
func FilterBySimilarity(docs []chromem.Result, minSimilarity float32, limit int) []chromem.Result {
	if minSimilarity <= 0 {
		minSimilarity = 0
	}
	filtered := make([]chromem.Result, 0, limit)
	for _, doc := range docs {
		if doc.Similarity < minSimilarity {
			continue
		}
		filtered = append(filtered, doc)
		if limit > 0 && len(filtered) >= limit {
			break
		}
	}
	return filtered
}