package health

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

const agentsDLQStream = "AGENTS_DLQ"

type dlqDepth struct {
	Messages uint64
}

func dlqStreamDepth(js nats.JetStreamContext) (dlqDepth, error) {
	if js == nil {
		return dlqDepth{}, fmt.Errorf("jetstream unavailable")
	}
	info, err := js.StreamInfo(agentsDLQStream)
	if err != nil {
		return dlqDepth{}, fmt.Errorf("dlq stream: %w", err)
	}
	if info == nil || info.State.Msgs == 0 {
		return dlqDepth{}, nil
	}
	return dlqDepth{Messages: info.State.Msgs}, nil
}