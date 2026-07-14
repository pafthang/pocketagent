package taskapis

import (
	"context"
	"net/http"

	"github.com/pafthang/pocketagent/pkgs/httpx"

	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
)

func verifyTaskInSpace(c echo.Context, pb *pbclient.Client, taskID, spaceID string) error {
	if taskID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "task id required")
	}
	task, err := pb.GetTaskByCorrelationID(taskID)
	if err != nil {
		return httpx.MapPocketError(c, err)
	}
	if task.SpaceID != spaceID {
		return echo.NewHTTPError(http.StatusNotFound, "task not found")
	}
	return nil
}

func subscribeTaskEvents(nc *natsclient.Client, taskID string) (chan []byte, *nats.Subscription, error) {
	events := make(chan []byte, 32)
	subject := natsclient.EventSubject(taskID)

	sub, err := nc.Subscribe(subject, func(ctx context.Context, msg *nats.Msg) {
		select {
		case events <- msg.Data:
		default:
		}
	})
	if err != nil {
		return nil, nil, err
	}

	return events, sub, nil
}

func isTerminalEvent(eventType string) bool {
	switch eventType {
	case models.EventCompleted, models.EventFailed, models.EventCancelled, models.EventTimeout:
		return true
	default:
		return false
	}
}
