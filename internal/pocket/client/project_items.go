package client

import (
	"fmt"

	"github.com/pafthang/pocketagent/pkgs/models"
)



func (c *Client) CreateProjectItem(item models.ProjectItem) (models.ProjectItem, error) {
	record, err := c.CreateRecord(ProjectItemsCollection, projectItemRecordData(item))
	if err != nil {
		return models.ProjectItem{}, err
	}
	return projectItemFromRecord(record), nil
}

func (c *Client) GetProjectItem(id string) (models.ProjectItem, error) {
	record, err := c.GetRecord(ProjectItemsCollection, id)
	if err != nil {
		return models.ProjectItem{}, err
	}
	return projectItemFromRecord(record), nil
}

func (c *Client) ListProjectItems(opts ListOptions) ([]models.ProjectItem, int, error) {
	records, total, err := c.ListRecordsOpts(ProjectItemsCollection, opts)
	if err != nil {
		return nil, 0, err
	}
	out := make([]models.ProjectItem, 0, len(records))
	for _, record := range records {
		out = append(out, projectItemFromRecord(record))
	}
	return out, total, nil
}

func (c *Client) UpdateProjectItem(id string, item models.ProjectItem) (models.ProjectItem, error) {
	record, err := c.UpdateRecord(ProjectItemsCollection, id, projectItemRecordData(item))
	if err != nil {
		return models.ProjectItem{}, err
	}
	return projectItemFromRecord(record), nil
}

func (c *Client) DeleteProjectItem(id string) error {
	return c.DeleteRecord(ProjectItemsCollection, id)
}

func (c *Client) DeleteProjectItemsByProject(projectID string) error {
	filter := fmt.Sprintf("project_id = %q", projectID)
	for page := 1; ; page++ {
		items, _, err := c.ListProjectItems(ListOptions{Page: page, PerPage: 100, Filter: filter})
		if err != nil {
			return err
		}
		if len(items) == 0 {
			return nil
		}
		for _, item := range items {
			if err := c.DeleteProjectItem(item.ID); err != nil {
				return err
			}
		}
	}
}

func projectItemRecordData(item models.ProjectItem) map[string]interface{} {
	data := map[string]interface{}{
		"space_id":   item.SpaceID,
		"project_id": item.ProjectID,
		"title":      item.Title,
	}
	if item.Description != "" {
		data["description"] = item.Description
	}
	if item.Status != "" {
		data["status"] = item.Status
	}
	if item.Priority != "" {
		data["priority"] = item.Priority
	}
	if item.AssigneeIDs != nil {
		data["assignee_ids"] = item.AssigneeIDs
	}
	if item.ExecutionTaskID != "" {
		data["execution_task_id"] = item.ExecutionTaskID
	}
	if item.SortOrder != 0 {
		data["sort_order"] = item.SortOrder
	}
	if item.Tags != nil {
		data["tags"] = item.Tags
	}
	return data
}

func projectItemFromRecord(record map[string]interface{}) models.ProjectItem {
	item := models.ProjectItem{
		ID:              stringField(record, "id"),
		SpaceID:         stringField(record, "space_id"),
		ProjectID:       stringField(record, "project_id"),
		Title:           stringField(record, "title"),
		Description:     stringField(record, "description"),
		Status:          stringField(record, "status"),
		Priority:        stringField(record, "priority"),
		ExecutionTaskID: stringField(record, "execution_task_id"),
		CreatedAt:       stringField(record, "created"),
		UpdatedAt:       stringField(record, "updated"),
	}
	if v, ok := record["sort_order"]; ok {
		switch n := v.(type) {
		case float64:
			item.SortOrder = int(n)
		case int:
			item.SortOrder = n
		}
	}
	item.AssigneeIDs = stringSliceField(record, "assignee_ids")
	item.Tags = stringSliceField(record, "tags")
	return item
}

// ProjectItemsFilter returns a PocketBase filter for items in a project.
func ProjectItemsFilter(spaceID, projectID string) string {
	return fmt.Sprintf("space_id = %q && project_id = %q", spaceID, projectID)
}
