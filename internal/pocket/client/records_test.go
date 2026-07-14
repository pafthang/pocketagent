package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateRecord(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method", http.StatusMethodNotAllowed)
			return
		}
		if !strings.HasSuffix(r.URL.Path, "/api/collections/agents/records") {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":   "rec-1",
			"name": "Test Agent",
		})
	}))
	defer srv.Close()

	c := New(srv.URL)
	rec, err := c.CreateRecord("agents", map[string]interface{}{"name": "Test Agent"})
	if err != nil {
		t.Fatalf("CreateRecord: %v", err)
	}
	if rec["id"] != "rec-1" {
		t.Fatalf("id = %v", rec["id"])
	}
}

func TestGetRecord(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || !strings.HasSuffix(r.URL.Path, "/records/rec-1") {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"id": "rec-1", "name": "Agent"})
	}))
	defer srv.Close()

	c := New(srv.URL)
	rec, err := c.GetRecord("agents", "rec-1")
	if err != nil {
		t.Fatalf("GetRecord: %v", err)
	}
	if rec["name"] != "Agent" {
		t.Fatalf("name = %v", rec["name"])
	}
}

func TestListRecordsPagination(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "2" || r.URL.Query().Get("perPage") != "10" {
			http.Error(w, "bad query", http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(listResponse{
			Page:       2,
			PerPage:    10,
			TotalItems: 15,
			Items:      []map[string]interface{}{{"id": "rec-2"}},
		})
	}))
	defer srv.Close()

	c := New(srv.URL)
	items, total, err := c.ListRecords("agents", 2, 10)
	if err != nil {
		t.Fatalf("ListRecords: %v", err)
	}
	if total != 15 || len(items) != 1 {
		t.Fatalf("total=%d len=%d", total, len(items))
	}
}

func TestDeleteRecord(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := New(srv.URL)
	if err := c.DeleteRecord("agents", "rec-1"); err != nil {
		t.Fatalf("DeleteRecord: %v", err)
	}
}

func TestCreateRecordAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"message":"validation failed"}`, http.StatusBadRequest)
	}))
	defer srv.Close()

	c := New(srv.URL)
	_, err := c.CreateRecord("agents", map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok || apiErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected error: %T %v", err, err)
	}
}