package main

import (
	"net/http"

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

	natsClient, err := nats.NewClient("nats://nats:4222")
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer natsClient.Close()

	pbClient := pocketbase.NewClient("http://pocketbase:8090")

	// Routes
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})

	e.POST("/agents", func(c echo.Context) error {
		return createAgent(c, pbClient)
	})
	e.POST("/tasks", func(c echo.Context) error {
		return createTask(c, natsClient)
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func createAgent(c echo.Context, pb *pocketbase.Client) error {
	var agent models.Agent
	if err := c.Bind(&agent); err != nil {
		return err
	}

	// Save to PocketBase
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
