package tools

import "testing"

func TestParseArgsJSON(t *testing.T) {
	args := ParseArgs(`{"query":"hello"}`)
	if ArgString(args, "query") != "hello" {
		t.Fatalf("unexpected args: %v", args)
	}
}

func TestParseArgsInvalidJSON(t *testing.T) {
	args := ParseArgs("plain query")
	if len(args) != 0 {
		t.Fatalf("expected empty args, got %v", args)
	}
}