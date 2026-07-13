package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/internal/common"
	"github.com/pafthang/pocketagent/internal/models"
	"github.com/pafthang/pocketagent/internal/nats"
	"github.com/pafthang/pocketagent/internal/pocketbase"
	"github.com/pafthang/pocketagent/internal/service"
)

func main() {
	s := service.New("api-gateway")
	s.AddHealth()

	natsClient, _ := nats.NewClient(s.Config.NatsURL)
	pbClient := pocketbase.NewClient("http://pocketbase:8090")

	s.Echo.POST("/agents", func(c echo.Context) error {
		return createAgent(c, pbClient)
	})
	s.Echo.POST("/tasks", func(c echo.Context) error {
		return createTask(c, natsClient)
	})
	s.Echo.GET("/ws/task/:taskId", wsTaskStream)

	s.Start()
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

	corrID := fmt.Sprintf("task-%d", time.Now().UnixNano())
	ctx := common.WithCorrelationID(context.Background(), corrID)

	if err := nc.PublishTask(ctx, task); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"status":        "task_queued",
		"correlation_id": corrID,
	})
}

func wsTaskStream(c echo.Context) error {
	taskID := c.Param("taskId")
	ws, err := c.WebSocketUpgrade()
	if err != nil {
		return err
	}
	defer ws.Close()

	for i := 0; i < 10; i++ {
		msg := map[string]interface{}{
			"task_id": taskID,
			"step":     i,
			"status":   "thinking",
		}
		if err := ws.WriteJSON(msg); err != nil {
			return err
		}
		time.Sleep(600 * time.Millisecond)
	}

	final := map[string]interface{}{
		"task_id": taskID,
		"status":  "completed",
		"result":  "Task completed successfully",
	}
	ws.WriteJSON(final)

	return nil
}
