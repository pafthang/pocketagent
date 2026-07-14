package projectapis

import (
	"context"
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
)

func verifyProjectInSpace(c echo.Context, pb *pbclient.Client, projectID, spaceID string) error {
	if projectID == "" {
		return echo.NewHTTPError(400, "project id required")
	}
	project, err := pb.GetProject(projectID)
	if err != nil {
		return httpx.MapPocketError(c, err)
	}
	if project.SpaceID != spaceID {
		return echo.NewHTTPError(404, "project not found")
	}
	return nil
}

func subscribeProjectEvents(nc *natsclient.Client, projectID string) (chan []byte, *nats.Subscription, error) {
	events := make(chan []byte, 32)
	subject := natsclient.EventSubject(projectID)

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

func isProjectTerminalEvent(data []byte) bool {
	var payload struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return false
	}
	return payload.Type == "dw_planning_complete"
}