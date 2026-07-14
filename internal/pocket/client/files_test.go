package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pafthang/pocketagent/pkgs/models"
)

func TestCreateFileRoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method", http.StatusMethodNotAllowed)
			return
		}
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		body["id"] = "file-1"
		_ = json.NewEncoder(w).Encode(body)
	}))
	defer srv.Close()

	c := New(srv.URL)
	file, err := c.CreateFile(models.StoredFile{
		SpaceID:     "space-1",
		VirtualPath: "/docs/readme.md",
		Name:        "readme.md",
		MimeType:    "text/markdown",
		Size:        42,
	})
	if err != nil {
		t.Fatalf("CreateFile: %v", err)
	}
	if file.ID != "file-1" || file.VirtualPath != "/docs/readme.md" {
		t.Fatalf("file = %+v", file)
	}
}

func TestGetFile(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/records/file-1") {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":           "file-1",
			"space_id":     "space-1",
			"virtual_path": "/a.txt",
			"name":         "a.txt",
			"mime_type":    "text/plain",
			"size":         10,
		})
	}))
	defer srv.Close()

	c := New(srv.URL)
	file, err := c.GetFile("file-1")
	if err != nil {
		t.Fatalf("GetFile: %v", err)
	}
	if file.Name != "a.txt" || file.Size != 10 {
		t.Fatalf("file = %+v", file)
	}
}
