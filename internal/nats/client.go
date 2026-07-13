package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/pafthang/pocketagent/internal/common"
	"github.com/pafthang/pocketagent/internal/models"
)

// Client wraps NATS + JetStream with helpers
type Client struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

func NewClient(url string) (*Client, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return &Client{nc: nc, js: js}, nil
}

// PublishTask publishes task with correlation ID
func (c *Client) PublishTask(ctx context.Context, task models.Task) error {
	msg := &nats.Msg{
		Subject: "agents.tasks." + task.AgentID,
		Data:    []byte(`{"prompt":"` + task.Prompt + `"}`),
		Header:  make(nats.Header),
	}

	if corrID := common.GetCorrelationID(ctx); corrID != "" {
		msg.Header.Set("X-Correlation-ID", corrID)
	}

	_, err := c.js.PublishMsg(msg)
	return err
}

// SubscribeWithCorrelation subscribes with correlation ID extraction
func (c *Client) SubscribeWithCorrelation(subject string, handler func(ctx context.Context, msg *nats.Msg)) (*nats.Subscription, error) {
	return c.js.Subscribe(subject, func(msg *nats.Msg) {
		ctx := context.Background()
		if corrID := msg.Header.Get("X-Correlation-ID"); corrID != "" {
			ctx = common.WithCorrelationID(ctx, corrID)
		}
		handler(ctx, msg)
	})
}

func (c *Client) Close() {
	c.nc.Close()
}
