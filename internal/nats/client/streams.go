package client

import (
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
)

const agentsStreamName = "AGENTS"

// EnsureStreams creates the JetStream stream for agent messaging if missing.
func EnsureStreams(js nats.JetStreamContext) error {
	cfg := &nats.StreamConfig{
		Name:      agentsStreamName,
		Subjects:  agentsStreamSubjects(),
		Storage:   nats.FileStorage,
		Retention: nats.LimitsPolicy,
	}

	_, err := js.AddStream(cfg)
	if err == nil {
		return EnsureDLQStream(js)
	}

	if errors.Is(err, nats.ErrStreamNameAlreadyInUse) {
		_, err = js.UpdateStream(cfg)
		if err != nil {
			return fmt.Errorf("update stream %s: %w", agentsStreamName, err)
		}
		return EnsureDLQStream(js)
	}

	return fmt.Errorf("add stream %s: %w", agentsStreamName, err)
}

func agentsStreamSubjects() []string {
	return []string{
		"agents.orchestrator.>",
		"agents.tasks.>",
		"agents.results.>",
		"agents.events.>",
		"agents.projects.plan.>",
	}
}
