package service

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/pafthang/pocketagent/internal/common"
)

// Consumer is a base NATS JetStream consumer
type Consumer struct {
	JS     nats.JetStreamContext
	Logger *log.Logger
	Name   string
}

func NewConsumer(name string, js nats.JetStreamContext) *Consumer {
	return &Consumer{
		JS:     js,
		Logger: common.NewLogger(name),
		Name:   name,
	}
}

// Subscribe subscribes to a subject with correlation ID support
func (c *Consumer) Subscribe(subject string, handler func(ctx context.Context, msg *nats.Msg)) (*nats.Subscription, error) {
	return c.JS.Subscribe(subject, func(msg *nats.Msg) {
		ctx := context.Background()

		// Extract correlation ID from headers
		if correlationID := msg.Header.Get("X-Correlation-ID"); correlationID != "" {
			ctx = common.WithCorrelationID(ctx, correlationID)
		}

		handler(ctx, msg)
	})
}
