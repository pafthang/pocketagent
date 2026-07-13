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

	// CRUD endpoints
	s.Echo.POST("/agents", createAgent(pb))
	s.Echo.GET("/agents/:id", getAgent(pb))
	s.Echo.PUT("/agents/:id", updateAgent(pb))
	s.Echo.DELETE("/agents/:id", deleteAgent(pb))
	s.Echo.GET("/agents", listAgents(pb))

	s.Start()
}

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

// TODO: implement getAgent, updateAgent, deleteAgent, listAgents

func getAgent(pb *pocketbase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		// TODO: fetch from PocketBase
		return c.JSON(200, map[string]string{"id": id, "status": "not_implemented_yet"})
	}
}

func updateAgent(pb *pocketbase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(501, map[string]string{"error": "not implemented"})
	}
}

func deleteAgent(pb *pocketbase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(501, map[string]string{"error": "not implemented"})
	}
}

func listAgents(pb *pocketbase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(501, map[string]string{"error": "not implemented"})
	}
}
