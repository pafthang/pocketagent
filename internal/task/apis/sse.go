package taskapis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	apimw "github.com/pafthang/pocketagent/pkgs/middle"

	"github.com/labstack/echo/v4"
	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
)

const sseHeartbeatInterval = 25 * time.Second

type sseTokenPayload struct {
	TaskID string `json:"task_id"`
	Step   int    `json:"step,omitempty"`
	Delta  string `json:"delta"`
}

func sseTaskStream(c echo.Context, nc *natsclient.Client, pb *pbclient.Client) error {
	taskID := c.Param("id")
	spaceID, ok := apimw.SpaceIDFromContext(c)
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": apimw.HeaderSpaceID + " header is required"})
	}
	if err := verifyTaskInSpace(c, pb, taskID, spaceID); err != nil {
		return err
	}

	includeAll := c.QueryParam("events") == "all"

	events, sub, err := subscribeTaskEvents(nc, taskID)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	res := c.Response()
	res.Header().Set(echo.HeaderContentType, "text/event-stream; charset=utf-8")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")
	res.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := res.Writer.(http.Flusher)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "streaming not supported")
	}

	if err := writeSSE(res, flusher, "connected", map[string]string{
		"task_id": taskID,
		"status":  "listening",
	}); err != nil {
		return err
	}

	ctx, cancel := contextWithCancel(c)
	defer cancel()

	go sseHeartbeat(ctx, res, flusher)

	for {
		select {
		case <-ctx.Done():
			return nil
		case data := <-events:
			var event models.TaskEvent
			if err := json.Unmarshal(data, &event); err != nil {
				continue
			}

			switch event.Type {
			case models.EventLLMToken:
				if err := writeSSE(res, flusher, "token", sseTokenPayload{
					TaskID: taskID,
					Step:   event.Step,
					Delta:  event.Message,
				}); err != nil {
					return nil
				}
			default:
				if includeAll {
					if err := writeSSE(res, flusher, "event", event); err != nil {
						return nil
					}
				}
			}

			if isTerminalEvent(event.Type) {
				_ = writeSSE(res, flusher, event.Type, event)
				return nil
			}
		}
	}
}

func writeSSE(res *echo.Response, flusher http.Flusher, event string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(res.Writer, "event: %s\ndata: %s\n\n", event, data); err != nil {
		return err
	}
	flusher.Flush()
	return nil
}

func sseHeartbeat(ctx context.Context, res *echo.Response, flusher http.Flusher) {
	ticker := time.NewTicker(sseHeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if _, err := fmt.Fprint(res.Writer, ": ping\n\n"); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func contextWithCancel(c echo.Context) (context.Context, context.CancelFunc) {
	ctx := c.Request().Context()
	ctx, cancel := context.WithTimeout(ctx, wsStreamTimeout)
	return ctx, cancel
}
