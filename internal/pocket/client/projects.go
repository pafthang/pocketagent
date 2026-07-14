package client

import (
	"fmt"

	"github.com/pafthang/pocketagent/pkgs/models"
)



func (c *Client) CreateProject(project models.Project) (models.Project, error) {
	record, err := c.CreateRecord(ProjectsCollection, projectRecordData(project))
	if err != nil {
		return models.Project{}, err
	}
	return projectFromRecord(record), nil
}

func (c *Client) GetProject(id string) (models.Project, error) {
	record, err := c.GetRecord(ProjectsCollection, id)
	if err != nil {
		return models.Project{}, err
	}
	return projectFromRecord(record), nil
}

func (c *Client) ListProjects(opts ListOptions) ([]models.Project, int, error) {
	records, total, err := c.ListRecordsOpts(ProjectsCollection, opts)
	if err != nil {
		return nil, 0, err
	}
	out := make([]models.Project, 0, len(records))
	for _, record := range records {
		out = append(out, projectFromRecord(record))
	}
	return out, total, nil
}

func (c *Client) UpdateProject(id string, project models.Project) (models.Project, error) {
	record, err := c.UpdateRecord(ProjectsCollection, id, projectRecordData(project))
	if err != nil {
		return models.Project{}, err
	}
	return projectFromRecord(record), nil
}

func (c *Client) DeleteProject(id string) error {
	return c.DeleteRecord(ProjectsCollection, id)
}

func projectRecordData(project models.Project) map[string]interface{} {
	data := map[string]interface{}{
		"space_id": project.SpaceID,
		"title":    project.Title,
	}
	if project.Goal != "" {
		data["goal"] = project.Goal
	}
	if project.Description != "" {
		data["description"] = project.Description
	}
	if project.Status != "" {
		data["status"] = project.Status
	}
	if project.PlanJSON != nil {
		data["plan_json"] = project.PlanJSON
	}
	if project.ParentTaskID != "" {
		data["parent_task_id"] = project.ParentTaskID
	}
	if project.CreatorID != "" {
		data["creator_id"] = project.CreatorID
	}
	if project.PlannerAgentID != "" {
		data["planner_agent_id"] = project.PlannerAgentID
	}
	if project.TeamAgentIDs != nil {
		data["team_agent_ids"] = project.TeamAgentIDs
	}
	if project.Tags != nil {
		data["tags"] = project.Tags
	}
	if project.StartedAt != "" {
		data["started_at"] = project.StartedAt
	}
	if project.CompletedAt != "" {
		data["completed_at"] = project.CompletedAt
	}
	if project.Metadata != nil {
		data["metadata"] = project.Metadata
	}
	return data
}

func projectFromRecord(record map[string]interface{}) models.Project {
	project := models.Project{
		ID:             stringField(record, "id"),
		SpaceID:        stringField(record, "space_id"),
		Title:          stringField(record, "title"),
		Goal:           stringField(record, "goal"),
		Description:    stringField(record, "description"),
		Status:         stringField(record, "status"),
		ParentTaskID:   stringField(record, "parent_task_id"),
		CreatorID:      stringField(record, "creator_id"),
		PlannerAgentID: stringField(record, "planner_agent_id"),
		StartedAt:      stringField(record, "started_at"),
		CompletedAt:    stringField(record, "completed_at"),
		CreatedAt:      stringField(record, "created"),
		UpdatedAt:      stringField(record, "updated"),
	}
	project.TeamAgentIDs = stringSliceField(record, "team_agent_ids")
	project.Tags = stringSliceField(record, "tags")
	if raw, ok := record["plan_json"].(map[string]interface{}); ok {
		project.PlanJSON = raw
	}
	if raw, ok := record["metadata"].(map[string]interface{}); ok {
		project.Metadata = raw
	}
	return project
}

// ProjectsFilter returns a PocketBase filter for space-scoped projects.
func ProjectsFilter(spaceID string) string {
	return fmt.Sprintf("space_id = %q", spaceID)
}
