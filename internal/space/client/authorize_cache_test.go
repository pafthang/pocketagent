package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pafthang/pocketagent/pkgs/models"
)

func TestAuthorizeCachedHitMiss(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/authorize" {
			http.NotFound(w, r)
			return
		}
		calls.Add(1)
		_ = json.NewEncoder(w).Encode(models.AuthorizeResponse{Allowed: true})
	}))
	defer srv.Close()

	c := New(srv.URL)
	c.EnableAuthorizeCache(30 * time.Second)

	resp1, err := c.AuthorizeCached("user-1", "token", "space-1", "agent:read")
	if err != nil {
		t.Fatalf("first authorize: %v", err)
	}
	if !resp1.Allowed {
		t.Fatal("expected allowed")
	}

	resp2, err := c.AuthorizeCached("user-1", "token", "space-1", "agent:read")
	if err != nil {
		t.Fatalf("cached authorize: %v", err)
	}
	if !resp2.Allowed {
		t.Fatal("expected allowed from cache")
	}
	if calls.Load() != 1 {
		t.Fatalf("upstream calls = %d, want 1", calls.Load())
	}
}

func TestAuthorizeCachedDifferentKeys(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		_ = json.NewEncoder(w).Encode(models.AuthorizeResponse{Allowed: true})
	}))
	defer srv.Close()

	c := New(srv.URL)
	c.EnableAuthorizeCache(30 * time.Second)

	if _, err := c.AuthorizeCached("user-1", "token", "space-1", "agent:read"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.AuthorizeCached("user-1", "token", "space-1", "task:read"); err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 2 {
		t.Fatalf("upstream calls = %d, want 2", calls.Load())
	}
}

func TestAuthorizeCachedExpiry(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		_ = json.NewEncoder(w).Encode(models.AuthorizeResponse{Allowed: true})
	}))
	defer srv.Close()

	c := New(srv.URL)
	c.EnableAuthorizeCache(20 * time.Millisecond)

	if _, err := c.AuthorizeCached("user-1", "token", "space-1", "agent:read"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(30 * time.Millisecond)
	if _, err := c.AuthorizeCached("user-1", "token", "space-1", "agent:read"); err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 2 {
		t.Fatalf("upstream calls = %d, want 2 after expiry", calls.Load())
	}
}

func TestAuthorizeCachedDoesNotCacheErrors(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		http.Error(w, `{"error":"upstream"}`, http.StatusBadGateway)
	}))
	defer srv.Close()

	c := New(srv.URL)
	c.EnableAuthorizeCache(30 * time.Second)

	if _, err := c.AuthorizeCached("user-1", "token", "space-1", "agent:read"); err == nil {
		t.Fatal("expected error")
	}
	if _, err := c.AuthorizeCached("user-1", "token", "space-1", "agent:read"); err == nil {
		t.Fatal("expected error on retry")
	}
	if calls.Load() != 2 {
		t.Fatalf("upstream calls = %d, want 2 (no error cache)", calls.Load())
	}
}
