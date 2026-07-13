package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pafthang/pocketagent/internal/models"
	"github.com/pafthang/pocketagent/internal/nats"
	"github.com/pafthang/pocketagent/internal/pocketbase"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	natsClient, _ := nats.NewClient("nats://nats:4222")
	pbClient := pocketbase.NewClient("http://pocketbase:8090")

	e.GET("/health", healthHandler)
	e.POST("/agents", func(c echo.Context) error { return createAgent(c, pbClient) })
	e.POST("/tasks", func(c echo.Context) error { return createTask(c, natsClient) })

	e.GET("/ws/task/:taskId", wsTaskStream)

	e.Logger.Fatal(e.Start(":8080"))
}

func healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
}

func wsTaskStream(c echo.Context) error {
	taskID := c.Param("taskId")
	ws, err := c.WebSocketUpgrade()
	if err != nil {
		c.Logger().Errorf("WebSocket upgrade failed: %v", err)
		return err
	}
	defer ws.Close()

	c.Logger().Infof("WebSocket connected for task: %s", taskID)

	for i := 0; i < 10; i++ {
		msg := map[string]interface{}{
			"task_id": taskID,
			"step":     i,
			"status":   "thinking",
			"message":  fmt.Sprintf("Step %d in progress...", i),
		}
		if err := ws.WriteJSON(msg); err != nil {
			c.Logger().Errorf("WebSocket write error: %v", err)
			return err
		}
		time.Sleep(600 * time.Millisecond)
	}

	// Final message
	final := map[string]interface{}{
		"task_id": taskID,
		"status":  "completed",
		"result":  "Task completed successfully",
	}
	ws.WriteJSON(final)

	return nil
}

func createAgent(c echo.Context, pb *pocketbase.Client) error {
	var agent models.Agent
	if err := c.Bind(&agent); err != nil {
		c.Logger().Error(err)
		return err
	}

	data := map[string]interface{}{
		"name":        agent.Name,
		"description": agent.Description,
		"model":       agent.Model,
		"system_prompt": agent.SystemPrompt,
	}
	_, err := pb.CreateRecord("agents", data)
	if err != nil {
		c.Logger().Errorf("PocketBase error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, agent)
}

func createTask(c echo.Context, nc *nats.Client) error {
	var task models.Task
	if err := c.Bind(&task); err != nil {
		c.Logger().Error(err)
		return err
	}
	// TODO: publish to NATS
	return c.JSON(http.StatusCreated, task)
}
