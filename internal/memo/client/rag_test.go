package client

import (
	"testing"

	"github.com/pafthang/pocketagent/internal/memo/chunk"
)

func TestFormatRAGLinesDedup(t *testing.T) {
	lines := formatRAGLines([]Document{
		{Content: "same", Similarity: 0.9},
		{Content: "same", Similarity: 0.8},
		{Content: "other", Similarity: 0.7},
	})
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestChunkTextClient(t *testing.T) {
	chunks := chunk.Text("alpha. beta. "+stringsRepeatClient("gamma ", 200), 120, 10)
	if len(chunks) < 2 {
		t.Fatalf("expected chunks, got %d", len(chunks))
	}
}

func stringsRepeatClient(s string, n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}
