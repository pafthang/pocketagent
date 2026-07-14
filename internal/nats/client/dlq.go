package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pafthang/pocketagent/pkgs/common"
)

// DLQStreamName is the JetStream stream holding failed agent messages.
const DLQStreamName = "AGENTS_DLQ"

const dlqStreamName = DLQStreamName

// SubjectDLQAll matches all dead-letter records.
const SubjectDLQAll = "agents.dlq.>"

// SubjectDLQPrefix is the per-service DLQ subject prefix.
const SubjectDLQPrefix = "agents.dlq."

// DLQMessage is a failed JetStream payload archived for inspection/replay.
type DLQMessage struct {
	Service       string          `json:"service"`
	Subject       string          `json:"subject"`
	Reason        string          `json:"reason"`
	Error         string          `json:"error,omitempty"`
	CorrelationID string          `json:"correlation_id,omitempty"`
	NumDelivered  uint64          `json:"num_delivered"`
	Payload       json.RawMessage `json:"payload"`
	FailedAt      string          `json:"failed_at"`
}

// EnsureDLQStream creates the JetStream DLQ stream if missing.
func EnsureDLQStream(js nats.JetStreamContext) error {
	cfg := &nats.StreamConfig{
		Name:      dlqStreamName,
		Subjects:  []string{SubjectDLQAll},
		Storage:   nats.FileStorage,
		Retention: nats.LimitsPolicy,
	}

	_, err := js.AddStream(cfg)
	if err == nil {
		return nil
	}

	if errors.Is(err, nats.ErrStreamNameAlreadyInUse) {
		_, err = js.UpdateStream(cfg)
		if err != nil {
			return fmt.Errorf("update stream %s: %w", dlqStreamName, err)
		}
		return nil
	}

	return fmt.Errorf("add stream %s: %w", dlqStreamName, err)
}

// DLQStreamStats reports current dead-letter queue depth.
type DLQStreamStats struct {
	Messages uint64 `json:"messages"`
	Bytes    uint64 `json:"bytes"`
}

// StreamStats returns message counts for the AGENTS_DLQ stream.
func StreamStats(js nats.JetStreamContext) (DLQStreamStats, error) {
	if js == nil {
		return DLQStreamStats{}, fmt.Errorf("jetstream unavailable")
	}
	info, err := js.StreamInfo(DLQStreamName)
	if err != nil {
		return DLQStreamStats{}, err
	}
	if info == nil {
		return DLQStreamStats{}, nil
	}
	return DLQStreamStats{
		Messages: info.State.Msgs,
		Bytes:    info.State.Bytes,
	}, nil
}

// PublishDLQ archives a failed message for later inspection.
func PublishDLQ(js nats.JetStreamContext, service, subject, reason, corrID string, numDelivered uint64, payload []byte, handlerErr error) error {
	if js == nil {
		return fmt.Errorf("jetstream unavailable")
	}

	errText := ""
	if handlerErr != nil {
		errText = handlerErr.Error()
	}

	envelope := DLQMessage{
		Service:       service,
		Subject:       subject,
		Reason:        reason,
		Error:         errText,
		CorrelationID: corrID,
		NumDelivered:  numDelivered,
		Payload:       json.RawMessage(payload),
		FailedAt:      time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	dlqSubject := SubjectDLQPrefix + service
	if _, err := js.Publish(dlqSubject, data); err != nil {
		return fmt.Errorf("publish dlq: %w", err)
	}
	common.DLQMessagesTotal.WithLabelValues(service, reason).Inc()
	return nil
}
