package client

import (
	"context"

	"github.com/pafthang/pocketagent/pkgs/models"
)

// PublishEvent emits a task progress event for WebSocket subscribers.
func (c *Client) PublishEvent(ctx context.Context, taskID string, event models.TaskEvent) error {
	event.TaskID = taskID
	return c.publishJSON(ctx, EventSubject(taskID), event)
}
