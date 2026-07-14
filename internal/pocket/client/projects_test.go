package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pafthang/pocketagent/pkgs/models"
)

func TestCreateProjectRoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method", http.StatusMethodNotAllowed)
			return
		}
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		body["id"] = "proj-1"
		_ = json.NewEncoder(w).Encode(body)
	}))
	defer srv.Close()

	c := New(srv.URL)
	project, err := c.CreateProject(models.Project{
		SpaceID: "space-1",
		Title:   "Ship feature",
		Goal:    "Build it",
		Status:  models.ProjectDraft,
	})
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}
	if project.ID != "proj-1" || project.Title != "Ship feature" {
		t.Fatalf("project = %+v", project)
	}
}

func TestGetProject(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/records/proj-1") {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":       "proj-1",
			"space_id": "space-1",
			"title":    "Alpha",
			"status":   models.ProjectPlanning,
		})
	}))
	defer srv.Close()

	c := New(srv.URL)
	project, err := c.GetProject("proj-1")
	if err != nil {
		t.Fatalf("GetProject: %v", err)
	}
	if project.Title != "Alpha" || project.Status != models.ProjectPlanning {
		t.Fatalf("project = %+v", project)
	}
}
