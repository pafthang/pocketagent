package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/pafthang/pocketagent/services/execution-service/tools"
	"github.com/pafthang/pocketagent/services/ollama-client/ollama"
)

// ReActExecutor with Memory support (RAG)
type ReActExecutor struct {
	Ollama       *ollama.Client
	Tools        []ollama.Tool
	MaxSteps     int
	MemoryClient MemoryClient // optional
}

// MemoryClient interface for RAG
type MemoryClient interface {
	Search(ctx context.Context, queryEmbedding []float32, limit int) ([]string, error)
}

func NewReActExecutor(ollamaClient *ollama.Client, tools []ollama.Tool) *ReActExecutor {
	return &ReActExecutor{
		Ollama:   ollamaClient,
		Tools:    tools,
		MaxSteps: 8,
	}
}

// WithMemory adds RAG capability
func (e *ReActExecutor) WithMemory(mc MemoryClient) *ReActExecutor {
	e.MemoryClient = mc
	return e
}

// ReActResult ...
type ReActResult struct {
	FinalAnswer  string
	Steps        []string
	ToolCalls    []string
	Observations []string
}

func (e *ReActExecutor) Execute(ctx context.Context, prompt string) (ReActResult, error) {
	var result ReActResult
	var history strings.Builder

	// === RAG: Retrieve relevant context from Memory ===
	if e.MemoryClient != nil {
		// TODO: Generate embedding for prompt using Ollama
		// For now we skip real embedding generation
		relevantDocs, _ := e.MemoryClient.Search(ctx, nil, 3)
		if len(relevantDocs) > 0 {
			history.WriteString("Relevant context from memory:\n")
			for _, doc := range relevantDocs {
				history.WriteString("- " + doc + "\n")
			}
		}
	}

	history.WriteString("Thought: I will solve this task using tools and memory context.\n")

	for step := 0; step < e.MaxSteps; step++ {
		fullPrompt := fmt.Sprintf(`%s

%s`, prompt, history.String())

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

		if strings.Contains(resp, "Action:") {
			toolName, args := parseToolCall(resp)
			result.ToolCalls = append(result.ToolCalls, toolName)

			observation := executeRealTool(toolName, args)
			result.Observations = append(result.Observations, observation)
			history.WriteString("Observation: " + observation + "\n")
		}

		if strings.Contains(resp, "Final Answer:") {
			result.FinalAnswer = extractFinalAnswer(resp)
			break
		}
	}

	return result, nil
}

// ... (остальные функции parseToolCall, executeRealTool, extractFinalAnswer остаются без изменений)
