package tools

import "testing"

func TestMCPToolName(t *testing.T) {
	name := mcpToolName("File System", "read_file")
	if name != "mcp__file_system__read_file" {
		t.Fatalf("unexpected name: %s", name)
	}
}

func TestNormalizeMCPSchema(t *testing.T) {
	schema := normalizeMCPSchema(nil)
	if schema["type"] != "object" {
		t.Fatalf("expected object schema: %v", schema)
	}
}