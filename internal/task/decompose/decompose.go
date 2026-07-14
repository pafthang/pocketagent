package decompose

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
	"github.com/pafthang/pocketagent/pkgs/ollama"
)

var jsonArrayPattern = regexp.MustCompile(`\[[\s\S]*\]`)

// Decomposer splits tasks into parallel subtasks using an LLM.
type Decomposer struct {
	Ollama      *ollama.Client
	PocketBase  *pbclient.Client
	Model       string
	MaxSubtasks int
	Log         *slog.Logger
}

func New(ollamaClient *ollama.Client, pb *pbclient.Client, model string, maxSubtasks int, log *slog.Logger) *Decomposer {
	if maxSubtasks <= 0 {
		maxSubtasks = 4
	}
	if model == "" {
		model = "llama3.1"
	}
	return &Decomposer{
		Ollama:      ollamaClient,
		PocketBase:  pb,
		Model:       model,
		MaxSubtasks: maxSubtasks,
		Log:         log,
	}
}

// Plan returns subtasks with optional per-agent assignment (supervisor workflow).
func (d *Decomposer) Plan(ctx context.Context, task models.Task) []SubtaskPlan {
	if isSupervisorWorkflow(task) {
		supervisor := &Supervisor{
			Ollama:      d.Ollama,
			PocketBase:  d.PocketBase,
			Model:       d.Model,
			MaxSubtasks: d.MaxSubtasks,
			Log:         d.Log,
		}
		if plans := supervisor.Plan(ctx, task); len(plans) > 0 {
			return plans
		}
	}

	prompts := d.Split(ctx, task)
	plans := make([]SubtaskPlan, len(prompts))
	for i, prompt := range prompts {
		plans[i] = SubtaskPlan{Prompt: prompt, AgentID: task.AgentID}
	}
	return plans
}

// Split returns subtask prompts for parallel execution.
func (d *Decomposer) Split(ctx context.Context, task models.Task) []string {
	model := d.resolveModel(task.AgentID)

	subtasks, err := d.decomposeWithLLM(ctx, model, task.Prompt)
	if err != nil {
		if d.Log != nil {
			d.Log.Warn("LLM decomposition failed, using fallback", "error", err)
		}
		return fallbackSplit(task.Prompt)
	}
	if len(subtasks) == 0 {
		return []string{task.Prompt}
	}

	if d.Log != nil {
		d.Log.Info("task decomposed", "subtasks", len(subtasks), "model", model)
	}
	return subtasks
}

func (d *Decomposer) resolveModel(agentID string) string {
	if agentID == "" || d.PocketBase == nil {
		return d.Model
	}

	agent, err := d.PocketBase.GetAgent(agentID)
	if err != nil || agent.Model == "" {
		return d.Model
	}
	return agent.Model
}

func (d *Decomposer) decomposeWithLLM(ctx context.Context, model, prompt string) ([]string, error) {
	if d.Ollama == nil {
		return nil, fmt.Errorf("ollama client is nil")
	}

	llmPrompt := fmt.Sprintf(`You are a task planner for an agent orchestrator.
Break the user task into %d or fewer independent parallel subtasks.
Each subtask must be self-contained and executable without the others.

Rules:
- If the task is already atomic, return a JSON array with one string (the original task).
- Return ONLY valid JSON: an array of strings.
- Do not include markdown or explanations.

User task:
%s`, d.MaxSubtasks, prompt)

	_ = ctx

	resp, err := d.Ollama.Generate(ollama.GenerateRequest{
		Model:  model,
		Prompt: llmPrompt,
		Format: "json",
		Stream: false,
	})
	if err != nil {
		return nil, err
	}

	return parseSubtasksResponse(resp, d.MaxSubtasks)
}

func parseSubtasksResponse(resp string, maxSubtasks int) ([]string, error) {
	resp = strings.TrimSpace(resp)
	if resp == "" {
		return nil, fmt.Errorf("empty LLM response")
	}

	candidates := []string{resp}
	if match := jsonArrayPattern.FindString(resp); match != "" && match != resp {
		candidates = append([]string{match}, candidates...)
	}

	var items []string
	var lastErr error
	for _, candidate := range candidates {
		if err := json.Unmarshal([]byte(candidate), &items); err != nil {
			lastErr = err
			continue
		}
		lastErr = nil
		break
	}
	if lastErr != nil {
		return nil, fmt.Errorf("parse subtasks JSON: %w", lastErr)
	}

	subtasks := make([]string, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		subtasks = append(subtasks, item)
		if maxSubtasks > 0 && len(subtasks) >= maxSubtasks {
			break
		}
	}

	if len(subtasks) == 0 {
		return nil, fmt.Errorf("no subtasks in LLM response")
	}
	return subtasks, nil
}
