package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pafthang/pocketagent/pkgs/models"
)

func TestCreateAgentRoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method", http.StatusMethodNotAllowed)
			return
		}
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		body["id"] = "agent-1"
		_ = json.NewEncoder(w).Encode(body)
	}))
	defer srv.Close()

	c := New(srv.URL)
	agent, err := c.CreateAgent(models.Agent{
		SpaceID: "space-1",
		Name:    "Planner",
		Model:   "llama3.1",
	})
	if err != nil {
		t.Fatalf("CreateAgent: %v", err)
	}
	if agent.ID != "agent-1" || agent.Name != "Planner" {
		t.Fatalf("agent = %+v", agent)
	}
}

func TestGetAgent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/records/agent-1") {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":       "agent-1",
			"space_id": "space-1",
			"name":     "Worker",
			"model":    "llama3.1",
		})
	}))
	defer srv.Close()

	c := New(srv.URL)
	agent, err := c.GetAgent("agent-1")
	if err != nil {
		t.Fatalf("GetAgent: %v", err)
	}
	if agent.Name != "Worker" || agent.SpaceID != "space-1" {
		t.Fatalf("agent = %+v", agent)
	}
}
