package main

import (
	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/internal/common"
	"github.com/pafthang/pocketagent/internal/models"
	"github.com/pafthang/pocketagent/internal/pocketbase"
	"github.com/pafthang/pocketagent/internal/service"
)

func main() {
	s := service.New("agent-service")
	s.AddHealth()

	pb := pocketbase.NewClient("http://pocketbase:8090")

	s.Echo.POST("/agents", createAgent(pb))
	s.Echo.GET("/agents/:id", getAgent(pb))
	s.Echo.PUT("/agents/:id", updateAgent(pb))
	s.Echo.DELETE("/agents/:id", deleteAgent(pb))
	s.Echo.GET("/agents", listAgents(pb))

	s.Start()
}

// Create
func createAgent(pb *pocketbase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var agent models.Agent
		if err := c.Bind(&agent); err != nil {
			return err
		}

		data := map[string]interface{}{
			"name":         agent.Name,
			"description":   agent.Description,
			"model":         agent.Model,
			"system_prompt": agent.SystemPrompt,
		}

		result, err := pb.CreateRecord("agents", data)
		if err != nil {
			return c.JSON(500, map[string]string{"error": err.Error()})
		}

		return c.JSON(201, result)
	}
}

// Read
func getAgent(pb *pocketbase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		// TODO: replace with real PocketBase get
		// For now return placeholder
		return c.JSON(200, map[string]interface{}{
			"id":   id,
			"note": "Implement real PocketBase get here",
		})
	}
}

// Update
func updateAgent(pb *pocketbase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		var agent models.Agent
		if err := c.Bind(&agent); err != nil {
			return err
		}

		// TODO: real update via PocketBase
		return c.JSON(200, map[string]interface{}{
			"id":      id,
			"updated": true,
			"note":    "Implement real update",
		})
	}
}

// Delete
func deleteAgent(pb *pocketbase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		// TODO: real delete
		return c.JSON(200, map[string]interface{}{
			"id":      id,
			"deleted": true,
		})
	}
}

// List
func listAgents(pb *pocketbase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: real list from PocketBase
		return c.JSON(200, map[string]interface{}{
			"agents": []interface{}{},
			"note":   "Implement real list from PocketBase",
		})
	}
}
