package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/pafthang/pocketagent/pkgs/models"
)

// Client wraps NATS + JetStream with helpers.
type Client struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

// NewClient connects to NATS and prepares JetStream.
func NewClient(url string) (*Client, error) {
	nc, js, err := Connect(url)
	if err != nil {
		return nil, err
	}

	return &Client{nc: nc, js: js}, nil
}

// Conn exposes the underlying NATS connection (e.g. for health checks).
func (c *Client) Conn() *nats.Conn {
	return c.nc
}

// JS exposes the JetStream context for advanced subscriptions.
func (c *Client) JS() nats.JetStreamContext {
	return c.js
}

// PublishOrchestrator enqueues a high-level task for the task service.
func (c *Client) PublishOrchestrator(ctx context.Context, task models.Task) error {
	return c.publishJSON(ctx, SubjectOrchestrator, task)
}

// PublishProjectPlan enqueues async planning for a project.
func (c *Client) PublishProjectPlan(ctx context.Context, cmd models.ProjectPlanCommand) error {
	return c.publishJSON(ctx, SubjectProjectsPlan, cmd)
}

// PublishRawEvent emits a JSON event on agents.events.{id} (tasks, projects, etc.).
func (c *Client) PublishRawEvent(ctx context.Context, id string, payload any) error {
	return c.publishJSON(ctx, EventSubject(id), payload)
}

// PublishSubtask dispatches a subtask to exec via agents.tasks.{id}.
func (c *Client) PublishSubtask(ctx context.Context, subtaskID string, task models.Task) error {
	return c.publishJSON(ctx, SubjectTasksPrefix+subtaskID, task)
}

// PublishResult publishes an exec result to agents.results.{key}.
func (c *Client) PublishResult(ctx context.Context, key, result string) error {
	msg := &nats.Msg{
		Subject: SubjectResultsPrefix + key,
		Data:    []byte(result),
		Header:  make(nats.Header),
	}
	common.InjectContextHeaders(ctx, msg.Header)

	_, err := c.js.PublishMsg(msg)
	return err
}

// Subscribe subscribes with correlation ID extraction.
func (c *Client) Subscribe(subject string, handler func(ctx context.Context, msg *nats.Msg)) (*nats.Subscription, error) {
	return c.js.Subscribe(subject, func(msg *nats.Msg) {
		handler(common.ContextFromNATSMsg(msg), msg)
	})
}

// Close closes the underlying NATS connection.
func (c *Client) Close() {
	c.nc.Close()
}

func (c *Client) publishJSON(ctx context.Context, subject string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	msg := &nats.Msg{
		Subject: subject,
		Data:    data,
		Header:  make(nats.Header),
	}
	common.InjectContextHeaders(ctx, msg.Header)

	_, err = c.js.PublishMsg(msg)
	return err
}
