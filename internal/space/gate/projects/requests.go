package projectapis

import (
	"strings"

	"github.com/pafthang/pocketagent/pkgs/models"
)

// CreateProjectRequest is the gate API body for POST /projects.
type CreateProjectRequest struct {
	Title          string   `json:"title"`
	Goal           string   `json:"goal"`
	Description    string   `json:"description"`
	PlannerAgentID string   `json:"planner_agent_id"`
	TeamAgentIDs   []string `json:"team_agent_ids"`
	Tags           []string `json:"tags"`
	Status         string   `json:"status"`
}

// ToModel builds a project record; title should already be normalized.
func (r CreateProjectRequest) ToModel(spaceID, creatorID, title string) models.Project {
	status := strings.TrimSpace(r.Status)
	if status == "" {
		status = models.ProjectDraft
	}
	return models.Project{
		SpaceID:        spaceID,
		Title:          title,
		Goal:           strings.TrimSpace(r.Goal),
		Description:    strings.TrimSpace(r.Description),
		Status:         status,
		PlannerAgentID: strings.TrimSpace(r.PlannerAgentID),
		TeamAgentIDs:   r.TeamAgentIDs,
		Tags:           r.Tags,
		CreatorID:      creatorID,
	}
}

// PatchProjectRequest is the gate API body for PATCH /projects/:id.
type PatchProjectRequest struct {
	Title          *string                `json:"title"`
	Goal           *string                `json:"goal"`
	Description    *string                `json:"description"`
	Status         *string                `json:"status"`
	PlanJSON       map[string]interface{} `json:"plan_json"`
	ParentTaskID   *string                `json:"parent_task_id"`
	PlannerAgentID *string                `json:"planner_agent_id"`
	TeamAgentIDs   *[]string              `json:"team_agent_ids"`
	Tags           *[]string              `json:"tags"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// ApplyPatch mutates project in place. titleNormalizer resolves title when Title is set.
func (r PatchProjectRequest) ApplyPatch(project *models.Project, titleNormalizer func(title, goal string) string) {
	if r.Title != nil && titleNormalizer != nil {
		project.Title = titleNormalizer(*r.Title, project.Goal)
	}
	if r.Goal != nil {
		project.Goal = strings.TrimSpace(*r.Goal)
	}
	if r.Description != nil {
		project.Description = strings.TrimSpace(*r.Description)
	}
	if r.Status != nil {
		project.Status = strings.TrimSpace(*r.Status)
	}
	if r.PlanJSON != nil {
		project.PlanJSON = r.PlanJSON
	}
	if r.ParentTaskID != nil {
		project.ParentTaskID = strings.TrimSpace(*r.ParentTaskID)
	}
	if r.PlannerAgentID != nil {
		project.PlannerAgentID = strings.TrimSpace(*r.PlannerAgentID)
	}
	if r.TeamAgentIDs != nil {
		project.TeamAgentIDs = *r.TeamAgentIDs
	}
	if r.Tags != nil {
		project.Tags = *r.Tags
	}
	if r.Metadata != nil {
		project.Metadata = r.Metadata
	}
}

// CreateProjectItemRequest is the gate API body for POST /projects/:id/items.
type CreateProjectItemRequest struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Status          string   `json:"status"`
	Priority        string   `json:"priority"`
	AssigneeIDs     []string `json:"assignee_ids"`
	ExecutionTaskID string   `json:"execution_task_id"`
	SortOrder       int      `json:"sort_order"`
	Tags            []string `json:"tags"`
}

// ToModel builds a project item record.
func (r CreateProjectItemRequest) ToModel(spaceID, projectID string) models.ProjectItem {
	status := strings.TrimSpace(r.Status)
	if status == "" {
		status = models.ItemInbox
	}
	return models.ProjectItem{
		SpaceID:         spaceID,
		ProjectID:       projectID,
		Title:           strings.TrimSpace(r.Title),
		Description:     strings.TrimSpace(r.Description),
		Status:          status,
		Priority:        strings.TrimSpace(r.Priority),
		AssigneeIDs:     r.AssigneeIDs,
		ExecutionTaskID: strings.TrimSpace(r.ExecutionTaskID),
		SortOrder:       r.SortOrder,
		Tags:            r.Tags,
	}
}

// PatchProjectItemRequest is the gate API body for PATCH /projects/:id/items/:itemId.
type PatchProjectItemRequest struct {
	Title           *string   `json:"title"`
	Description     *string   `json:"description"`
	Status          *string   `json:"status"`
	Priority        *string   `json:"priority"`
	AssigneeIDs     *[]string `json:"assignee_ids"`
	ExecutionTaskID *string   `json:"execution_task_id"`
	SortOrder       *int      `json:"sort_order"`
	Tags            *[]string `json:"tags"`
}

// ApplyPatch mutates a project item in place.
func (r PatchProjectItemRequest) ApplyPatch(item *models.ProjectItem) {
	if r.Title != nil {
		item.Title = strings.TrimSpace(*r.Title)
	}
	if r.Description != nil {
		item.Description = strings.TrimSpace(*r.Description)
	}
	if r.Status != nil {
		item.Status = strings.TrimSpace(*r.Status)
	}
	if r.Priority != nil {
		item.Priority = strings.TrimSpace(*r.Priority)
	}
	if r.AssigneeIDs != nil {
		item.AssigneeIDs = *r.AssigneeIDs
	}
	if r.ExecutionTaskID != nil {
		item.ExecutionTaskID = strings.TrimSpace(*r.ExecutionTaskID)
	}
	if r.SortOrder != nil {
		item.SortOrder = *r.SortOrder
	}
	if r.Tags != nil {
		item.Tags = *r.Tags
	}
}

// ParseGoalRequest is the gate API body for POST /projects/parse-goal.
type ParseGoalRequest struct {
	Description string `json:"description"`
}

// StartProjectRequest is the gate API body for POST /projects/start.
type StartProjectRequest struct {
	Description    string   `json:"description"`
	Title          string   `json:"title"`
	PlannerAgentID string   `json:"planner_agent_id"`
	TeamAgentIDs   []string `json:"team_agent_ids"`
}
