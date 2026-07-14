package agentapis

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
)

// CreateHandler handles POST /agents.
func CreateHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var agent models.Agent
		if err := c.Bind(&agent); err != nil {
			return err
		}
		if agent.Name == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "name is required"})
		}
		spaceID, err := requireSpaceID(c)
		if err != nil {
			return err
		}
		agent.SpaceID = spaceID

		created, err := pb.CreateAgent(agent)
		if err != nil {
			return mapPocketError(c, err)
		}

		return c.JSON(http.StatusCreated, created)
	}
}

// GetHandler handles GET /agents/:id.
func GetHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		agent, err := loadAgentInSpace(c, pb, c.Param("id"))
		if err != nil {
			return mapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, agent)
	}
}

// UpdateHandler handles PUT /agents/:id.
func UpdateHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		existing, err := loadAgentInSpace(c, pb, c.Param("id"))
		if err != nil {
			return mapPocketError(c, err)
		}

		var agent models.Agent
		if err := c.Bind(&agent); err != nil {
			return err
		}
		agent.SpaceID = existing.SpaceID

		updated, err := pb.UpdateAgent(c.Param("id"), agent)
		if err != nil {
			return mapPocketError(c, err)
		}

		return c.JSON(http.StatusOK, updated)
	}
}

// DeleteHandler handles DELETE /agents/:id.
func DeleteHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		if _, err := loadAgentInSpace(c, pb, c.Param("id")); err != nil {
			return mapPocketError(c, err)
		}
		if err := pb.DeleteAgent(c.Param("id")); err != nil {
			return mapPocketError(c, err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}

// ListHandler handles GET /agents.
func ListHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		perPage, _ := strconv.Atoi(c.QueryParam("per_page"))

		spaceID, err := requireSpaceID(c)
		if err != nil {
			return err
		}

		agents, total, err := pb.ListAgentsFilter(pbclient.ListOptions{
			Page:    page,
			PerPage: perPage,
			Filter:  agentsFilter(spaceID),
		})
		if err != nil {
			return mapPocketError(c, err)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"agents": agents,
			"total":  total,
		})
	}
}

func mapPocketError(c echo.Context, err error) error {
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		return httpErr
	}

	var apiErr *pbclient.APIError
	if errors.As(err, &apiErr) {
		status := apiErr.StatusCode
		if status < 400 {
			status = http.StatusInternalServerError
		}
		return c.JSON(status, map[string]string{"error": apiErr.Error()})
	}
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
}
