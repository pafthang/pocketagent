package agent

import (
	"fmt"
	"strings"

	"github.com/pafthang/pocketagent/pkgs/ollama/api"
)

// ToolRunner executes a tool by name with native or string arguments.
type ToolRunner func(name string, args map[string]interface{}) string

// LoopOptions configures a multi-turn chat tool loop.
type LoopOptions struct {
	Model    string
	Messages []api.ChatMessage
	Tools    []api.Tool
	MaxSteps int
	RunTool  ToolRunner
	Stream   bool
	OnStream api.StreamChunkHandler
}

// LoopResult is the outcome of a chat-based agent loop.
type LoopResult struct {
	FinalAnswer string
	Messages    []api.ChatMessage
	ToolCalls   []string
	Steps       int
}

// RunLoop drives Ollama /api/chat with native tool calling.
func RunLoop(client *api.Client, opts LoopOptions) (LoopResult, error) {
	var result LoopResult
	if client == nil {
		return result, fmt.Errorf("ollama client is nil")
	}
	if opts.MaxSteps <= 0 {
		opts.MaxSteps = 8
	}

	messages := append([]api.ChatMessage{}, opts.Messages...)

	for step := 0; step < opts.MaxSteps; step++ {
		var assistant api.ChatMessage
		var err error

		req := api.ChatRequest{
			Model:    opts.Model,
			Messages: messages,
			Tools:    opts.Tools,
		}

		if opts.Stream {
			assistant, err = client.ChatStream(req, opts.OnStream)
		} else {
			var resp api.ChatResponse
			resp, err = client.Chat(req)
			assistant = resp.Message
		}
		if err != nil {
			return result, err
		}

		result.Steps++
		messages = append(messages, assistant)

		if len(assistant.ToolCalls) > 0 {
			for _, call := range assistant.ToolCalls {
				name := call.Function.Name
				result.ToolCalls = append(result.ToolCalls, name)
				observation := ""
				if opts.RunTool != nil {
					observation = opts.RunTool(name, call.Function.Arguments)
				}
				messages = append(messages, api.ChatMessage{
					Role:     "tool",
					ToolName: name,
					Content:  observation,
				})
			}
			continue
		}

		if strings.TrimSpace(assistant.Content) != "" {
			result.FinalAnswer = strings.TrimSpace(assistant.Content)
			break
		}
	}

	result.Messages = messages
	return result, nil
}