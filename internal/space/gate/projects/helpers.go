package projectapis

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
)

func loadProjectInSpace(pb *pbclient.Client, spaceID, id string) (models.Project, error) {
	project, err := pb.GetProject(id)
	if err != nil {
		return models.Project{}, err
	}
	if project.SpaceID != spaceID {
		return models.Project{}, &pbclient.APIError{StatusCode: http.StatusNotFound, Message: "project not found"}
	}
	return project, nil
}

func loadProjectItemInSpace(pb *pbclient.Client, spaceID, projectID, itemID string) (models.ProjectItem, error) {
	item, err := pb.GetProjectItem(itemID)
	if err != nil {
		return models.ProjectItem{}, err
	}
	if item.SpaceID != spaceID || item.ProjectID != projectID {
		return models.ProjectItem{}, &pbclient.APIError{StatusCode: http.StatusNotFound, Message: "project item not found"}
	}
	return item, nil
}

func projectItemIDs(pb *pbclient.Client, spaceID, projectID string) ([]string, error) {
	items, _, err := pb.ListProjectItems(pbclient.ListOptions{
		Page: 1, PerPage: 500, Filter: pbclient.ProjectItemsFilter(spaceID, projectID),
	})
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids, nil
}

func validatePlannerAgent(pb *pbclient.Client, spaceID, agentID string) error {
	if agentID == "" {
		return nil
	}
	agent, err := pb.GetAgent(agentID)
	if err != nil {
		return &pbclient.APIError{StatusCode: http.StatusBadRequest, Message: "invalid planner_agent_id"}
	}
	if agent.SpaceID != "" && agent.SpaceID != spaceID {
		return &pbclient.APIError{StatusCode: http.StatusBadRequest, Message: "planner agent does not belong to this space"}
	}
	return nil
}

func parsePageParams(c echo.Context) (int, int) {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 50
	}
	if perPage > 200 {
		perPage = 200
	}
	return page, perPage
}