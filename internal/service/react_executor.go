package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/pafthang/pocketagent/services/execution-service/tools"
	"github.com/pafthang/pocketagent/services/ollama-client/ollama"
)

// ReActExecutor - полноценная реализация с реальным выполнением инструментов
type ReActExecutor struct {
	Ollama   *ollama.Client
	Tools    []ollama.Tool
	MaxSteps int
}

func NewReActExecutor(ollamaClient *ollama.Client, tools []ollama.Tool) *ReActExecutor {
	return &ReActExecutor{
		Ollama:   ollamaClient,
		Tools:    tools,
		MaxSteps: 8,
	}
}

// ReActResult содержит результат выполнения

type ReActResult struct {
	FinalAnswer string
	Steps       []string
	ToolCalls   []string
	Observations []string
}

// Execute выполняет полноценный ReAct цикл с реальными инструментами
func (e *ReActExecutor) Execute(ctx context.Context, prompt string) (ReActResult, error) {
	var result ReActResult
	var history strings.Builder

	history.WriteString("Thought: Я буду использовать инструменты для решения задачи.\n")

	for step := 0; step < e.MaxSteps; step++ {
		fullPrompt := fmt.Sprintf(`
Задача: %s

История:
%s

Инструкция: Отвечай только в формате:
- Thought: ...
- Action: tool_name аргументы
- Final Answer: ...
`, prompt, history.String())

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

		// Если модель хочет вызвать инструмент
		if strings.Contains(resp, "Action:") {
			toolName, args := parseToolCall(resp)
			result.ToolCalls = append(result.ToolCalls, fmt.Sprintf("%s(%s)", toolName, args))

			// === РЕАЛЬНОЕ ВЫПОЛНЕНИЕ ИНСТРУМЕНТА ===
			observation := executeRealTool(toolName, args)
			result.Observations = append(result.Observations, observation)

			history.WriteString("Observation: " + observation + "\n")
			result.Steps = append(result.Steps, "Observation: "+observation)
		}

		// Если модель дала финальный ответ
		if strings.Contains(resp, "Final Answer:") {
			result.FinalAnswer = extractFinalAnswer(resp)
			break
		}
	}

	return result, nil
}

// Вспомогательные функции

func parseToolCall(text string) (string, string) {
	if idx := strings.Index(text, "Action:"); idx != -1 {
		remaining := strings.TrimSpace(text[idx+7:])
		parts := strings.SplitN(remaining, " ", 2)
		if len(parts) >= 1 {
			tool := strings.TrimSpace(parts[0])
			args := ""
			if len(parts) == 2 {
				args = strings.TrimSpace(parts[1])
			}
			return tool, args
		}
	}
	return "unknown_tool", ""
}

func executeRealTool(toolName, args string) string {
	switch strings.ToLower(toolName) {
	case "web_search", "search_web", "search":
		res, err := tools.WebSearch(args)
		if err != nil {
			return "Error during web search: " + err.Error()
		}
		return res

	case "scrape_page", "scrape", "browse":
		res, err := tools.ScrapePage(args)
		if err != nil {
			return "Error scraping page: " + err.Error()
		}
		return res

	default:
		return fmt.Sprintf("Tool '%s' called with args '%s' (no real implementation yet)", toolName, args)
	}
}

func extractFinalAnswer(text string) string {
	if idx := strings.Index(text, "Final Answer:"); idx != -1 {
		return strings.TrimSpace(text[idx+13:])
	}
	return text
}
