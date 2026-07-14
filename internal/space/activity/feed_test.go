package activity

import (
	"testing"

	"github.com/pafthang/pocketagent/pkgs/models"
)

func TestFromTaskEventThinking(t *testing.T) {
	entry := FromTaskEvent(models.StoredTaskEvent{
		ID:        "1",
		TaskID:    "task-1",
		EventType: models.EventLLMToken,
		Message:   "hello",
		CreatedAt: "2026-01-01T00:00:00Z",
	})
	if entry.Type != "thinking" || entry.Content != "hello" {
		t.Fatalf("unexpected entry: %+v", entry)
	}
}

func TestFromTaskEventError(t *testing.T) {
	entry := FromTaskEvent(models.StoredTaskEvent{
		ID:        "2",
		EventType: models.EventFailed,
		Message:   "boom",
		CreatedAt: "2026-01-02T00:00:00Z",
	})
	if entry.Type != "error" {
		t.Fatalf("expected error, got %q", entry.Type)
	}
}

func TestBuildFeedSortsNewestFirst(t *testing.T) {
	entries := BuildFeed(
		[]models.StoredTaskEvent{
			{ID: "a", EventType: models.EventQueued, Message: "old", CreatedAt: "2026-01-01T00:00:00Z"},
			{ID: "b", EventType: models.EventCompleted, Message: "new", CreatedAt: "2026-01-02T00:00:00Z"},
		},
		nil,
		false,
		0,
	)
	if len(entries) != 2 || entries[0].Content != "new" {
		t.Fatalf("unexpected order: %+v", entries)
	}
}

func TestBuildFeedIncludesAuditWhenEnabled(t *testing.T) {
	entries := BuildFeed(nil, []models.AuditLog{
		{ID: "x", Action: "member.add", CreatedAt: "2026-01-01T00:00:00Z"},
	}, true, 0)
	if len(entries) != 1 || entries[0].Source != SourceAudit {
		t.Fatalf("unexpected entries: %+v", entries)
	}
}
