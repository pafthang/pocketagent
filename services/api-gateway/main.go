package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pafthang/pocketagent/internal/common"
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

func createTask(c echo.Context, nc *nats.Client) error {
	var task models.Task
	if err := c.Bind(&task); err != nil {
		return err
	}

	// Generate correlation ID
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

// ... other handlers (createAgent, wsTaskStream)
