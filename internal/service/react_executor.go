package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/pafthang/pocketagent/services/ollama-client/ollama"
)

// ReActExecutor provides reusable ReAct logic
type ReActExecutor struct {
	Ollama *ollama.Client
	Tools  []ollama.Tool
	MaxSteps int
}

func NewReActExecutor(ollamaClient *ollama.Client, tools []ollama.Tool) *ReActExecutor {
	return &ReActExecutor{
		Ollama: ollamaClient,
		Tools:  tools,
		MaxSteps: 8,
	}
}

// Execute runs ReAct loop
func (e *ReActExecutor) Execute(ctx context.Context, prompt string) (string, error) {
	var history strings.Builder
	history.WriteString("Thought: Starting task\n")

	for step := 0; step < e.MaxSteps; step++ {
		resp, err := e.Ollama.Generate(ollama.GenerateRequest{
			Model:  "llama3.1",
			Prompt: prompt + "\n" + history.String(),
			Tools:  e.Tools,
		})
		if err != nil {
			return "", err
		}

		history.WriteString(fmt.Sprintf("Step %d: %s\n", step, resp))

		if strings.Contains(resp, "Final Answer") {
			break
		}
	}

	return history.String(), nil
}
