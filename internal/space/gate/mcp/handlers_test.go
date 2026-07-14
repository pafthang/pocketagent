package mcpapis

import "testing"

func TestBuildMCPServerFromRequestStdio(t *testing.T) {
	server, err := buildMCPServerFromRequest("space-1", "fs", "stdio", "npx", []string{"-y", "pkg"}, "", nil, boolPtr(true), true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if server.SpaceID != "space-1" || server.Command != "npx" || !server.Enabled {
		t.Fatalf("unexpected server: %#v", server)
	}
}

func TestBuildMCPServerFromRequestHTTP(t *testing.T) {
	server, err := buildMCPServerFromRequest("space-1", "remote", "http", "", nil, "https://example.com/mcp", nil, boolPtr(true), true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if server.URL != "https://example.com/mcp" {
		t.Fatalf("unexpected url: %s", server.URL)
	}
}

func TestBuildMCPServerFromRequestValidation(t *testing.T) {
	if _, err := buildMCPServerFromRequest("space-1", "", "stdio", "npx", nil, "", nil, nil, true); err == nil {
		t.Fatal("expected name validation error")
	}
	if _, err := buildMCPServerFromRequest("space-1", "x", "stdio", "", nil, "", nil, nil, true); err == nil {
		t.Fatal("expected command validation error")
	}
	if _, err := buildMCPServerFromRequest("space-1", "x", "http", "", nil, "", nil, nil, true); err == nil {
		t.Fatal("expected url validation error")
	}
}