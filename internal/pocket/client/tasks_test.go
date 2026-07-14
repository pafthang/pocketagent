package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pafthang/pocketagent/pkgs/models"
)

func TestCreateTaskRoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method", http.StatusMethodNotAllowed)
			return
		}
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		body["id"] = "task-1"
		_ = json.NewEncoder(w).Encode(body)
	}))
	defer srv.Close()

	c := New(srv.URL)
	task, err := c.CreateTask(models.Task{
		SpaceID:       "space-1",
		AgentID:       "agent-1",
		Prompt:        "do work",
		Status:        models.TaskQueued,
		CorrelationID: "corr-1",
	})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	if task.ID != "task-1" || task.CorrelationID != "corr-1" {
		t.Fatalf("task = %+v", task)
	}
}

func TestGetTaskByCorrelationID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.RawQuery, "correlation_id") {
			http.Error(w, "bad filter", http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"page":       1,
			"perPage":    1,
			"totalItems": 1,
			"items": []map[string]interface{}{
				{
					"id":             "task-9",
					"space_id":       "space-1",
					"correlation_id": "corr-9",
					"prompt":         "run",
					"status":         models.TaskRunning,
				},
			},
		})
	}))
	defer srv.Close()

	c := New(srv.URL)
	task, err := c.GetTaskByCorrelationID("corr-9")
	if err != nil {
		t.Fatalf("GetTaskByCorrelationID: %v", err)
	}
	if task.ID != "task-9" || task.Status != models.TaskRunning {
		t.Fatalf("task = %+v", task)
	}
}
