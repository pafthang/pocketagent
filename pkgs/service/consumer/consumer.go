package consumer

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/nats-io/nats.go"
	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	"github.com/pafthang/pocketagent/pkgs/common"
)

const defaultMaxDeliver = 5

// Consumer is a base NATS JetStream consumer.
type Consumer struct {
	JS         nats.JetStreamContext
	Log        *slog.Logger
	Name       string
	MaxDeliver int
}

// New creates a JetStream consumer helper.
func New(name string, js nats.JetStreamContext) *Consumer {
	return &Consumer{
		JS:         js,
		Log:        common.NewSlogLogger(name),
		Name:       name,
		MaxDeliver: defaultMaxDeliver,
	}
}

// Subscribe subscribes to a subject with correlation ID support.
func (c *Consumer) Subscribe(subject string, handler func(ctx context.Context, msg *nats.Msg)) (*nats.Subscription, error) {
	return c.JS.Subscribe(subject, func(msg *nats.Msg) {
		handler(common.ContextFromNATSMsg(msg), msg)
	})
}

// SubscribeJSON unmarshals the message and ack/nak on success/failure.
func SubscribeJSON[T any](c *Consumer, subject string, handler func(ctx context.Context, payload T) error) (*nats.Subscription, error) {
	return c.JS.Subscribe(subject, func(msg *nats.Msg) {
		ctx := common.ContextFromNATSMsg(msg)
		corrID := common.GetCorrelationID(ctx)

		var payload T
		if err := json.Unmarshal(msg.Data, &payload); err != nil {
			c.Log.Error("failed to unmarshal message",
				"subject", subject,
				"error", err,
				"correlation_id", corrID,
			)
			c.moveToDLQ(msg, subject, "unmarshal", corrID, err)
			_ = msg.Ack()
			return
		}

		if err := handler(ctx, payload); err != nil {
			deliveries := deliveryCount(msg)
			c.Log.Error("message handler failed",
				"subject", subject,
				"error", err,
				"correlation_id", corrID,
				"num_delivered", deliveries,
			)

			maxDeliver := c.maxDeliver()
			if deliveries >= uint64(maxDeliver) {
				c.moveToDLQ(msg, subject, "handler", corrID, err)
				_ = msg.Ack()
				return
			}

			_ = msg.Nak()
			return
		}

		_ = msg.Ack()
	})
}

func (c *Consumer) maxDeliver() int {
	if c == nil || c.MaxDeliver <= 0 {
		return defaultMaxDeliver
	}
	return c.MaxDeliver
}

func deliveryCount(msg *nats.Msg) uint64 {
	if msg == nil {
		return 1
	}
	meta, err := msg.Metadata()
	if err != nil || meta == nil {
		return 1
	}
	if meta.NumDelivered == 0 {
		return 1
	}
	return meta.NumDelivered
}

func (c *Consumer) moveToDLQ(msg *nats.Msg, subject, reason, corrID string, handlerErr error) {
	if c == nil || c.JS == nil {
		return
	}
	if err := natsclient.PublishDLQ(c.JS, c.Name, subject, reason, corrID, deliveryCount(msg), msg.Data, handlerErr); err != nil && c.Log != nil {
		c.Log.Error("failed to publish dlq message",
			"subject", subject,
			"reason", reason,
			"correlation_id", corrID,
			"error", err,
		)
	}
}