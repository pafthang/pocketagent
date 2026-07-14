package agentapis

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
	"github.com/pafthang/pocketagent/pkgs/models"
)

func agentsFilter(spaceID string) string {
	return fmt.Sprintf("space_id = %q", spaceID)
}

func requireSpaceID(c echo.Context) (string, error) {
	spaceID, ok := apimw.SpaceIDFromContext(c)
	if !ok {
		return "", echo.NewHTTPError(http.StatusBadRequest, apimw.HeaderSpaceID+" header is required")
	}
	return spaceID, nil
}

func loadAgentInSpace(c echo.Context, pb *pbclient.Client, id string) (models.Agent, error) {
	spaceID, err := requireSpaceID(c)
	if err != nil {
		return models.Agent{}, err
	}

	agent, err := pb.GetAgent(id)
	if err != nil {
		return models.Agent{}, err
	}
	if agent.SpaceID != "" && agent.SpaceID != spaceID {
		return models.Agent{}, &pbclient.APIError{StatusCode: http.StatusNotFound, Message: "agent not found"}
	}
	return agent, nil
}
