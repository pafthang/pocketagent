package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApplyAuthSetsAuthorizationHeader(t *testing.T) {
	c := New("http://example.com")
	c.Token = "secret-token"

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	c.applyAuth(req)

	if got := req.Header.Get("Authorization"); got != "secret-token" {
		t.Fatalf("Authorization = %q", got)
	}
}

func TestApplyAuthSkipsWhenTokenEmpty(t *testing.T) {
	c := New("http://example.com")
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	c.applyAuth(req)
	if req.Header.Get("Authorization") != "" {
		t.Fatal("expected no Authorization header")
	}
}

func TestDoGetUsesAuth(t *testing.T) {
	var seenAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := New(srv.URL)
	c.Token = "pb-admin"
	if _, err := c.doGet(srv.URL); err != nil {
		t.Fatalf("doGet: %v", err)
	}
	if seenAuth != "pb-admin" {
		t.Fatalf("seenAuth = %q", seenAuth)
	}
}