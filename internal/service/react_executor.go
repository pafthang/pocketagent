package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/pafthang/pocketagent/services/execution-service/tools"
	"github.com/pafthang/pocketagent/services/ollama-client/ollama"
)

// ReActExecutor with real tool execution
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

// Execute runs ReAct with real tool execution
type ReActResult struct {
	FinalAnswer string
	Steps       []string
	ToolCalls   []string
}

func (e *ReActExecutor) Execute(ctx context.Context, prompt string) (ReActResult, error) {
	var result ReActResult
	var history strings.Builder

	history.WriteString("Thought: I will solve this using available tools when needed.\n")

	for step := 0; step < e.MaxSteps; step++ {
		fullPrompt := fmt.Sprintf(`%s

Current history:
%s

Respond with Thought, Action (tool_name args), or Final Answer.`, prompt, history.String())

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

		// Check if model wants to call a tool
		if strings.Contains(resp, "Action:") {
			toolName, args := parseToolCall(resp)
			result.ToolCalls = append(result.ToolCalls, toolName)

			// Execute real tool
			observation := executeRealTool(toolName, args)
			history.WriteString("Observation: " + observation + "\n")
			result.Steps = append(result.Steps, "Observation: "+observation)
		}

		if strings.Contains(resp, "Final Answer:") {
			result.FinalAnswer = extractFinalAnswer(resp)
			break
		}
	}

	return result, nil
}

func parseToolCall(text string) (string, string) {
	// Very simple parser
	if idx := strings.Index(text, "Action:"); idx != -1 {
		remaining := text[idx+7:]
		parts := strings.SplitN(remaining, " ", 2)
		if len(parts) == 2 {
			return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		}
		return strings.TrimSpace(parts[0]), ""
	}
	return "unknown", ""
}

func executeRealTool(toolName, args string) string {
	switch toolName {
	case "web_search", "search_web":
		result, err := tools.WebSearch(args)
		if err != nil {
			return "Error searching web: " + err.Error()
		}
		return result
	case "scrape_page", "scrape":
		result, err := tools.ScrapePage(args)
		if err != nil {
			return "Error scraping page: " + err.Error()
		}
		return result
	default:
		return fmt.Sprintf("Tool '%s' executed with args: %s (simulated)", toolName, args)
	}
}

func extractFinalAnswer(text string) string {
	if idx := strings.Index(text, "Final Answer:"); idx != -1 {
		return strings.TrimSpace(text[idx+13:])
	}
	return text
}
