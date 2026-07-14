package decompose

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
	"github.com/pafthang/pocketagent/pkgs/ollama"
)

type supervisorAssignment struct {
	Prompt  string `json:"prompt"`
	AgentID string `json:"agent_id"`
}

// Supervisor assigns decomposed subtasks to worker agents.
type Supervisor struct {
	Ollama      *ollama.Client
	PocketBase  *pbclient.Client
	Model       string
	MaxSubtasks int
	Log         *slog.Logger
}

func (s *Supervisor) Plan(ctx context.Context, task models.Task) []SubtaskPlan {
	workers := resolveWorkerAgents(s.PocketBase, task)
	if len(workers) == 0 {
		if s.Log != nil {
			s.Log.Warn("supervisor workflow without workers, falling back to default split")
		}
		return nil
	}

	model := s.resolveModel(task.AgentID)
	assignments, err := s.planWithLLM(ctx, model, task.Prompt, workers)
	if err != nil {
		if s.Log != nil {
			s.Log.Warn("supervisor planning failed, using fallback", "error", err)
		}
		return fallbackSupervisorPlan(task.Prompt, workers, s.MaxSubtasks)
	}
	if len(assignments) == 0 {
		return fallbackSupervisorPlan(task.Prompt, workers, s.MaxSubtasks)
	}

	allowed := make(map[string]struct{}, len(workers))
	for _, id := range workers {
		allowed[id] = struct{}{}
	}

	plans := make([]SubtaskPlan, 0, len(assignments))
	for i, item := range assignments {
		prompt := strings.TrimSpace(item.Prompt)
		if prompt == "" {
			continue
		}
		agentID := strings.TrimSpace(item.AgentID)
		if _, ok := allowed[agentID]; !ok {
			agentID = workers[i%len(workers)]
		}
		plans = append(plans, SubtaskPlan{Prompt: prompt, AgentID: agentID})
		if s.MaxSubtasks > 0 && len(plans) >= s.MaxSubtasks {
			break
		}
	}
	if len(plans) == 0 {
		return fallbackSupervisorPlan(task.Prompt, workers, s.MaxSubtasks)
	}
	return plans
}

func (s *Supervisor) resolveModel(agentID string) string {
	if agentID == "" || s.PocketBase == nil {
		return s.Model
	}
	agent, err := s.PocketBase.GetAgent(agentID)
	if err != nil || agent.Model == "" {
		return s.Model
	}
	return agent.Model
}

func (s *Supervisor) planWithLLM(ctx context.Context, model, prompt string, workers []string) ([]supervisorAssignment, error) {
	if s.Ollama == nil {
		return nil, fmt.Errorf("ollama client is nil")
	}

	workerLines := strings.Join(workers, ", ")
	llmPrompt := fmt.Sprintf(`You are a multi-agent supervisor.
Break the user task into %d or fewer parallel subtasks and assign each to one worker agent.

Available worker agent IDs (use only these exact values):
%s

Return ONLY valid JSON array:
[{"prompt":"subtask description","agent_id":"worker-id"}]

User task:
%s`, s.MaxSubtasks, workerLines, prompt)

	_ = ctx
	resp, err := s.Ollama.Generate(ollama.GenerateRequest{
		Model:  model,
		Prompt: llmPrompt,
		Format: "json",
		Stream: false,
	})
	if err != nil {
		return nil, err
	}

	resp = strings.TrimSpace(resp)
	if match := jsonArrayPattern.FindString(resp); match != "" {
		resp = match
	}

	var assignments []supervisorAssignment
	if err := json.Unmarshal([]byte(resp), &assignments); err != nil {
		return nil, err
	}
	return assignments, nil
}

func resolveWorkerAgents(pb *pbclient.Client, task models.Task) []string {
	if len(task.WorkerAgentIDs) > 0 {
		return task.WorkerAgentIDs
	}
	if task.AgentID == "" || pb == nil {
		return nil
	}

	agent, err := pb.GetAgent(task.AgentID)
	if err != nil {
		return nil
	}

	raw, ok := agent.Config["worker_agent_ids"]
	if !ok {
		return nil
	}
	switch ids := raw.(type) {
	case []string:
		return ids
	case []interface{}:
		out := make([]string, 0, len(ids))
		for _, item := range ids {
			if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
				out = append(out, strings.TrimSpace(s))
			}
		}
		return out
	default:
		return nil
	}
}

func fallbackSupervisorPlan(prompt string, workers []string, max int) []SubtaskPlan {
	prompts := fallbackSplit(prompt)
	plans := make([]SubtaskPlan, 0, len(prompts))
	for i, p := range prompts {
		plans = append(plans, SubtaskPlan{
			Prompt:  p,
			AgentID: workers[i%len(workers)],
		})
		if max > 0 && len(plans) >= max {
			break
		}
	}
	return plans
}

func isSupervisorWorkflow(task models.Task) bool {
	return strings.EqualFold(task.Workflow, models.WorkflowSupervisor) || len(task.WorkerAgentIDs) > 0
}
