package projects

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/pafthang/pocketagent/internal/exec/tools"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
	"github.com/pafthang/pocketagent/pkgs/ollama"
)

// EventPublisher emits planning progress for a project.
type EventPublisher func(ctx context.Context, projectID string, event map[string]any) error

// Planner runs multi-phase project planning and persists results.
type Planner struct {
	PB      *pbclient.Client
	Ollama  *ollama.Client
	Model   string
	Search  tools.Config
	Log     *slog.Logger
	Publish EventPublisher
}

var planPhases = []struct {
	key     string
	message string
}{
	{models.PlanPhaseGoalAnalysis, "Analyzing goal and scope"},
	{models.PlanPhaseResearch, "Gathering context"},
	{models.PlanPhasePRD, "Drafting requirements"},
	{models.PlanPhaseTasks, "Breaking down into tasks"},
	{models.PlanPhaseTeam, "Assigning team roles"},
}

// Run executes planning phases for a project.
func (p *Planner) Run(ctx context.Context, cmd models.ProjectPlanCommand) error {
	if p == nil || p.PB == nil {
		return fmt.Errorf("planner not configured")
	}
	project, err := p.PB.GetProject(cmd.ProjectID)
	if err != nil {
		return err
	}
	if project.SpaceID != cmd.SpaceID {
		return fmt.Errorf("project not in space")
	}
	if project.Status == models.ProjectCancelled {
		return nil
	}

	plan := project.PlanJSON
	if plan == nil {
		plan = map[string]interface{}{}
	}
	phases, _ := plan["phases"].(map[string]interface{})
	if phases == nil {
		phases = map[string]interface{}{}
		plan["phases"] = phases
	}

	project.Status = models.ProjectPlanning
	project.PlanJSON = plan
	if _, err := p.PB.UpdateProject(project.ID, project); err != nil {
		return err
	}

	goal := strings.TrimSpace(project.Goal)
	if goal == "" {
		goal = strings.TrimSpace(project.Description)
	}
	if goal == "" {
		goal = project.Title
	}

	for _, phase := range planPhases {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		project, err = p.PB.GetProject(cmd.ProjectID)
		if err != nil {
			return err
		}
		if project.Status == models.ProjectCancelled {
			return nil
		}

		plan = project.PlanJSON
		if plan == nil {
			plan = map[string]interface{}{}
		}
		plan["current_phase"] = phase.key
		plan["phase_message"] = phase.message
		project.PlanJSON = plan
		if _, err := p.PB.UpdateProject(project.ID, project); err != nil {
			return err
		}
		p.emitPhase(ctx, project.ID, phase.key, phase.message)

		switch phase.key {
		case models.PlanPhaseGoalAnalysis:
			analysis, err := ParseGoal(ctx, p.Ollama, p.Model, goal)
			if err != nil {
				return p.fail(ctx, project, err)
			}
			phases[phase.key] = analysis
		case models.PlanPhaseResearch:
			phases[phase.key] = p.researchPhase(goal, phases)
		case models.PlanPhasePRD:
			prd, err := p.generatePRD(ctx, goal, phases)
			if err != nil {
				return p.fail(ctx, project, err)
			}
			phases[phase.key] = map[string]interface{}{"content": prd}
		case models.PlanPhaseTasks:
			if err := p.createTasks(ctx, project, goal, phases); err != nil {
				return p.fail(ctx, project, err)
			}
		case models.PlanPhaseTeam:
			phases[phase.key] = map[string]interface{}{
				"team_agent_ids":   project.TeamAgentIDs,
				"planner_agent_id": project.PlannerAgentID,
			}
		}

		plan["phases"] = phases
		project.PlanJSON = plan
		if _, err := p.PB.UpdateProject(project.ID, project); err != nil {
			return err
		}
	}

	project, err = p.PB.GetProject(cmd.ProjectID)
	if err != nil {
		return err
	}
	plan = project.PlanJSON
	if plan == nil {
		plan = map[string]interface{}{}
	}
	plan["planning_done"] = true
	plan["current_phase"] = models.PlanPhaseTeam
	plan["phase_message"] = "Planning complete"
	project.PlanJSON = plan
	project.Status = models.ProjectAwaitingApproval
	if _, err := p.PB.UpdateProject(project.ID, project); err != nil {
		return err
	}
	p.emitComplete(ctx, project.ID, "")
	return nil
}

func (p *Planner) fail(ctx context.Context, project models.Project, err error) error {
	plan := project.PlanJSON
	if plan == nil {
		plan = map[string]interface{}{}
	}
	plan["planning_error"] = err.Error()
	project.PlanJSON = plan
	project.Status = models.ProjectFailed
	_, _ = p.PB.UpdateProject(project.ID, project)
	p.emitComplete(ctx, project.ID, err.Error())
	return err
}

func (p *Planner) emitPhase(ctx context.Context, projectID, phase, message string) {
	if p.Publish == nil {
		return
	}
	_ = p.Publish(ctx, projectID, map[string]any{
		"type":       "dw_planning_phase",
		"project_id": projectID,
		"phase":      phase,
		"message":    message,
	})
}

func (p *Planner) emitComplete(ctx context.Context, projectID, errMsg string) {
	if p.Publish == nil {
		return
	}
	payload := map[string]any{
		"type":       "dw_planning_complete",
		"project_id": projectID,
	}
	if errMsg != "" {
		payload["error"] = errMsg
	}
	_ = p.Publish(ctx, projectID, payload)
}

// ParseGoal analyzes a goal synchronously (gate handler).
func ParseGoal(ctx context.Context, oc *ollama.Client, model, goal string) (map[string]interface{}, error) {
	goal = strings.TrimSpace(goal)
	if goal == "" {
		return nil, fmt.Errorf("goal is required")
	}
	if oc == nil {
		return defaultGoalAnalysis(goal), nil
	}
	if model == "" {
		model = "llama3.1"
	}
	prompt := fmt.Sprintf(`Analyze this project goal and respond with JSON only:
{"domain":"","complexity":"low|medium|high","ai_roles":[],"human_roles":[],"questions":[]}

Goal:
%s`, goal)
	text, err := oc.Generate(ollama.GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
		Format: "json",
	})
	if err != nil {
		return defaultGoalAnalysis(goal), nil
	}
	var out map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(text)), &out); err != nil {
		return defaultGoalAnalysis(goal), nil
	}
	return out, nil
}

func defaultGoalAnalysis(goal string) map[string]interface{} {
	return map[string]interface{}{
		"domain":      "general",
		"complexity":  "medium",
		"ai_roles":    []string{"planner", "executor"},
		"human_roles": []string{"reviewer"},
		"questions":   []string{},
		"description": goal,
	}
}

func (p *Planner) generatePRD(ctx context.Context, goal string, phases map[string]interface{}) (string, error) {
	if p.Ollama == nil {
		return fmt.Sprintf("# Requirements\n\n%s", goal), nil
	}
	model := p.Model
	if model == "" {
		model = "llama3.1"
	}
	prompt := fmt.Sprintf(`Write a concise PRD (markdown) for this goal. Keep it under 400 words.

Goal:
%s`, goal)
	return p.Ollama.Generate(ollama.GenerateRequest{Model: model, Prompt: prompt, Stream: false})
}

type plannedTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (p *Planner) createTasks(ctx context.Context, project models.Project, goal string, phases map[string]interface{}) error {
	tasks, err := p.planTasks(ctx, goal)
	if err != nil {
		return err
	}
	for i, t := range tasks {
		title := strings.TrimSpace(t.Title)
		if title == "" {
			continue
		}
		item := models.ProjectItem{
			SpaceID:     project.SpaceID,
			ProjectID:   project.ID,
			Title:       title,
			Description: strings.TrimSpace(t.Description),
			Status:      models.ItemInbox,
			SortOrder:   i + 1,
		}
		if len(project.TeamAgentIDs) > 0 {
			item.AssigneeIDs = []string{project.TeamAgentIDs[i%len(project.TeamAgentIDs)]}
			item.Status = models.ItemAssigned
		}
		if _, err := p.PB.CreateProjectItem(item); err != nil {
			return err
		}
	}
	phases[models.PlanPhaseTasks] = map[string]interface{}{
		"count": len(tasks),
	}
	return nil
}

func (p *Planner) planTasks(ctx context.Context, goal string) ([]plannedTask, error) {
	if p.Ollama == nil {
		return []plannedTask{
			{Title: "Define scope", Description: goal},
			{Title: "Implement solution", Description: "Execute the main deliverable"},
			{Title: "Review and ship", Description: "Validate output and finalize"},
		}, nil
	}
	model := p.Model
	if model == "" {
		model = "llama3.1"
	}
	prompt := fmt.Sprintf(`Break this goal into 3-8 actionable tasks. Respond with JSON only:
{"tasks":[{"title":"","description":""}]}

Goal:
%s`, goal)
	text, err := p.Ollama.Generate(ollama.GenerateRequest{
		Model: model, Prompt: prompt, Stream: false, Format: "json",
	})
	if err != nil {
		return nil, err
	}
	var payload struct {
		Tasks []plannedTask `json:"tasks"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(text)), &payload); err != nil {
		return nil, err
	}
	if len(payload.Tasks) == 0 {
		return nil, fmt.Errorf("planner returned no tasks")
	}
	return payload.Tasks, nil
}

func (p *Planner) researchPhase(goal string, phases map[string]interface{}) map[string]interface{} {
	query := researchQuery(goal, phases)
	out := map[string]interface{}{
		"query":   query,
		"summary": fmt.Sprintf("Context prepared for: %s", truncate(goal, 200)),
	}
	if p == nil {
		return out
	}
	results, err := tools.SearchWeb(p.Search, query)
	if err != nil {
		if p.Log != nil {
			p.Log.Warn("project research search failed", "error", err, "query", query)
		}
		out["error"] = err.Error()
		return out
	}
	out["summary"] = truncate(results, 4000)
	out["source"] = "web"
	return out
}

func researchQuery(goal string, phases map[string]interface{}) string {
	query := strings.TrimSpace(goal)
	if analysis, ok := phases[models.PlanPhaseGoalAnalysis].(map[string]interface{}); ok {
		if domain, ok := analysis["domain"].(string); ok && strings.TrimSpace(domain) != "" {
			query = query + " " + strings.TrimSpace(domain)
		}
	}
	return query
}

func truncate(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
