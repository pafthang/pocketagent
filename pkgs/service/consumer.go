package service

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/pafthang/pocketagent/pkgs/service/consumer"
)

// Consumer is a base NATS JetStream consumer.
type Consumer = consumer.Consumer

// NewConsumer creates a JetStream consumer helper.
func NewConsumer(name string, js nats.JetStreamContext) *Consumer {
	return consumer.New(name, js)
}

// SubscribeJSON unmarshals the message and ack/nak on success/failure.
func SubscribeJSON[T any](c *Consumer, subject string, handler func(ctx context.Context, payload T) error) (*nats.Subscription, error) {
	return consumer.SubscribeJSON(c, subject, handler)
}