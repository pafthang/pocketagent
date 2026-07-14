package taskapis

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	apimw "github.com/pafthang/pocketagent/pkgs/middle"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

const wsStreamTimeout = 10 * time.Minute

func wsTaskStream(c echo.Context, nc *natsclient.Client, pb *pbclient.Client) error {
	taskID := c.Param("taskId")
	spaceID, ok := apimw.SpaceIDFromContext(c)
	if !ok {
		return echo.NewHTTPError(400, apimw.HeaderSpaceID+" header is required")
	}
	if err := verifyTaskInSpace(c, pb, taskID, spaceID); err != nil {
		return err
	}

	ws, err := wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	events, sub, err := subscribeTaskEvents(nc, taskID)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	ctx, cancel := context.WithTimeout(c.Request().Context(), wsStreamTimeout)
	defer cancel()

	go func() {
		for {
			if _, _, err := ws.ReadMessage(); err != nil {
				cancel()
				return
			}
		}
	}()

	connected := models.NewTaskEvent(taskID, models.EventConnected, "listening", "subscribed to task events")
	if err := ws.WriteJSON(connected); err != nil {
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

			var event models.TaskEvent
			if err := json.Unmarshal(data, &event); err != nil {
				continue
			}
			if isTerminalEvent(event.Type) {
				return nil
			}
		}
	}
}
