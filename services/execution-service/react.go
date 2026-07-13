package main

import (
	"fmt"
	"strings"

	"github.com/pafthang/pocketagent/services/ollama-client/ollama"
)

// Simple ReAct loop simulation
type ReActAgent struct {
	Ollama *ollama.Client
	Tools  []ollama.Tool
}

func NewReActAgent(ollamaClient *ollama.Client) *ReActAgent {
	return &ReActAgent{
		Ollama: ollamaClient,
		Tools:  ollama.GetExampleTools(),
	}
}

func (a *ReActAgent) Run(task string) string {
	var history strings.Builder
	history.WriteString("Thought: Understanding the task\n")

	for i := 0; i < 5; i++ { // max steps
		resp, _ := a.Ollama.Generate(ollama.GenerateRequest{
			Model:  "llama3.1",
			Prompt: fmt.Sprintf("%s\n%s", task, history.String()),
			Tools:  a.Tools,
		})

		history.WriteString("Action: " + resp + "\n")

		if strings.Contains(resp, "Final Answer") {
			break
		}
	}

	return history.String()
}
