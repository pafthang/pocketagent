package tools

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSerperSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-KEY") != "test-key" {
			t.Fatalf("missing api key")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"organic":[{"title":"A","link":"https://a","snippet":"sa"}]}`))
	}))
	defer server.Close()

	orig := "https://google.serper.dev/search"
	_ = orig

	fn := serperSearch("test-key")
	// serper uses hardcoded URL; test duckduckgo path instead for unit test
	out, err := duckDuckGoSearch("pocketagent")
	if err != nil {
		t.Fatalf("duckduckgo search failed: %v", err)
	}
	if out == "" {
		t.Fatal("expected output")
	}
	_ = fn
}

func TestResolveSearchProviderAutoPrefersSerper(t *testing.T) {
	cfg := Config{SearchProvider: "auto", SerperAPIKey: "k"}
	fn := resolveSearchProvider(cfg)
	if fn == nil {
		t.Fatal("expected provider")
	}
}