package projectapis

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
)

const wsProjectTimeout = 30 * time.Minute

func wsProjectStream(c echo.Context, nc *natsclient.Client, pb *pbclient.Client) error {
	projectID := c.Param("projectId")
	spaceID, ok := apimw.SpaceIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(400, apimw.HeaderSpaceID+" header is required")
	}
	if err := verifyProjectInSpace(c, pb, projectID, spaceID); err != nil {
		return err
	}

	ws, err := wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	events, sub, err := subscribeProjectEvents(nc, projectID)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	ctx, cancel := context.WithTimeout(c.Request().Context(), wsProjectTimeout)
	defer cancel()

	go func() {
		for {
			if _, _, err := ws.ReadMessage(); err != nil {
				cancel()
				return
			}
		}
	}()

	if err := ws.WriteJSON(map[string]interface{}{
		"type":       "connected",
		"project_id": projectID,
		"message":    "subscribed to project planning events",
	}); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case data := <-events:
			if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
				return err
			}
			if isProjectTerminalEvent(data) {
				return nil
			}
		}
	}
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
