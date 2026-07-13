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

	// WebSocket streaming
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
		return err
	}
	defer ws.Close()

	for i := 0; i < 10; i++ {
		msg := map[string]string{"step": fmt.Sprintf("Step %d for task %s", i, taskID), "status": "thinking"}
		if err := ws.WriteJSON(msg); err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func createAgent(c echo.Context, pb *pocketbase.Client) error {
	var agent models.Agent
	if err := c.Bind(&agent); err != nil {
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, agent)
}

func createTask(c echo.Context, nc *nats.Client) error {
	var task models.Task
	if err := c.Bind(&task); err != nil {
		return err
	}
	// TODO: publish to NATS
	return c.JSON(http.StatusCreated, task)
}
