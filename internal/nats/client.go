package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/pafthang/pocketagent/internal/models"
)

type Client struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

// NewClient creates new NATS + JetStream client
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

// PublishTask publishes task to execution service
func (c *Client) PublishTask(task models.Task) error {
	// implementation
	return nil
}
