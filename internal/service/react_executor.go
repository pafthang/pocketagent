package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pafthang/pocketagent/services/ollama-client/ollama"
)

// ReActExecutor provides a more realistic ReAct implementation
type ReActExecutor struct {
	Ollama   *ollama.Client
	Tools    []ollama.Tool
	MaxSteps int
}

func NewReActExecutor(ollamaClient *ollama.Client, tools []ollama.Tool) *ReActExecutor {
	return &ReActExecutor{
		Ollama:   ollamaClient,
		Tools:    tools,
		MaxSteps: 6,
	}
}

// Execute runs a more realistic ReAct loop
type ReActResult struct {
	FinalAnswer string
	Steps       []string
	ToolCalls   []string
}

func (e *ReActExecutor) Execute(ctx context.Context, prompt string) (ReActResult, error) {
	var result ReActResult
	var history strings.Builder

	history.WriteString("Thought: I need to solve this task step by step.\n")

	for step := 0; step < e.MaxSteps; step++ {
		fullPrompt := fmt.Sprintf(`%s

%s

Respond with either:
- Thought: ...
- Action: tool_name with args
- Final Answer: ...`, prompt, history.String())

		resp, err := e.Ollama.Generate(ollama.GenerateRequest{
			Model:  "llama3.1",
			Prompt: fullPrompt,
			Tools:  e.Tools,
		})
		if err != nil {
			return result, err
		}

		stepLog := fmt.Sprintf("Step %d: %s", step, resp)
		result.Steps = append(result.Steps, stepLog)
		history.WriteString(stepLog + "\n")

		// Check for tool call
		if strings.Contains(resp, "Action:") {
			toolCall := extractToolCall(resp)
			result.ToolCalls = append(result.ToolCalls, toolCall)
			// TODO: actually execute the tool here
			// For now we simulate observation
			observation := fmt.Sprintf("Observation: Tool '%s' returned some result.", toolCall)
			history.WriteString(observation + "\n")
			result.Steps = append(result.Steps, observation)
		}

		if strings.Contains(resp, "Final Answer:") {
			result.FinalAnswer = extractFinalAnswer(resp)
			break
		}
	}

	return result, nil
}

func extractToolCall(text string) string {
	// Simple extraction
	if idx := strings.Index(text, "Action:"); idx != -1 {
		end := strings.Index(text[idx:], "\n")
		if end == -1 {
			end = len(text)
		}
		return strings.TrimSpace(text[idx+7 : idx+end])
	}
	return "unknown_tool"
}

func extractFinalAnswer(text string) string {
	if idx := strings.Index(text, "Final Answer:"); idx != -1 {
		return strings.TrimSpace(text[idx+13:])
	}
	return text
}
