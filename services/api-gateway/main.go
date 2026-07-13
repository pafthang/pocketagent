package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pafthang/pocketagent/internal/models"
	"github.com/pafthang/pocketagent/internal/nats"
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

	// Routes
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})

	e.POST("/agents", createAgent)
	e.POST("/tasks", createTask)

	e.Logger.Fatal(e.Start(":8080"))
}

func createAgent(c echo.Context) error {
	var agent models.Agent
	if err := c.Bind(&agent); err != nil {
		return err
	}
	// TODO: save to PocketBase
	return c.JSON(http.StatusCreated, agent)
}

func createTask(c echo.Context) error {
	var task models.Task
	if err := c.Bind(&task); err != nil {
		return err
	}
	// TODO: publish to NATS
	return c.JSON(http.StatusCreated, task)
}
