package client

import (
	"encoding/json"
	"testing"
)

func TestDLQMessageJSON(t *testing.T) {
	msg := DLQMessage{
		Service:       "exec",
		Subject:       SubjectOrchestrator,
		Reason:        "handler",
		Error:         "react failed",
		CorrelationID: "task-1",
		NumDelivered:  5,
		Payload:       json.RawMessage(`{"prompt":"hi"}`),
		FailedAt:      "2026-07-14T12:00:00Z",
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	var decoded DLQMessage
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Service != "exec" || decoded.Reason != "handler" {
		t.Fatalf("unexpected decode: %+v", decoded)
	}
}
