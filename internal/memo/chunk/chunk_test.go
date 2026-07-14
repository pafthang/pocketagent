package chunk

import "testing"

func TestTextSingle(t *testing.T) {
	chunks := Text("short text", 100, 10)
	if len(chunks) != 1 || chunks[0] != "short text" {
		t.Fatalf("unexpected chunks: %#v", chunks)
	}
}

func TestTextMultiple(t *testing.T) {
	text := stringsRepeat("word ", 300)
	chunks := Text(text, 200, 20)
	if len(chunks) < 2 {
		t.Fatalf("expected multiple chunks, got %d", len(chunks))
	}
}

func stringsRepeat(s string, n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}