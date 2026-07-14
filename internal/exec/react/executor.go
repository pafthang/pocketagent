package react

import (
	"context"
	"fmt"

	memoclient "github.com/pafthang/pocketagent/internal/memo/client"
	"github.com/pafthang/pocketagent/pkgs/models"
	"github.com/pafthang/pocketagent/pkgs/ollama"
)

// Executor runs a chat-based tool loop with optional memory and streaming.
type Executor struct {
	Ollama       *ollama.Client
	Tools        []ollama.Tool
	LLMModel     string
	MaxSteps     int
	MemoryClient *memoclient.Client
	ExecuteTool  ToolExecutor
	StreamLLM    bool
	OnStream     ollama.StreamChunkHandler
}

// New creates a ReAct executor.
func New(ollamaClient *ollama.Client, tools []ollama.Tool, llmModel string) *Executor {
	if llmModel == "" {
		llmModel = "llama3.1"
	}
	return &Executor{
		Ollama:   ollamaClient,
		Tools:    tools,
		LLMModel: llmModel,
		MaxSteps: 8,
	}
}

// WithMemory adds RAG capability via memo client.
func (e *Executor) WithMemory(mc *memoclient.Client) *Executor {
	e.MemoryClient = mc
	return e
}

// WithStreaming enables Ollama chat streaming and chunk callbacks.
func (e *Executor) WithStreaming(enabled bool, handler ollama.StreamChunkHandler) *Executor {
	e.StreamLLM = enabled
	e.OnStream = handler
	return e
}

// Result is the outcome of a ReAct run.
type Result struct {
	FinalAnswer  string
	Steps        int
	ToolCalls    []string
	Observations []string
}

// Execute runs the chat agent loop for a task using optional agent configuration.
func (e *Executor) Execute(ctx context.Context, task models.Task, agent models.Agent, userProfile string) (Result, error) {
	return e.ExecuteWithOverlay(ctx, task, agent, nil, userProfile)
}

// ExecuteWithOverlay runs ReAct with optional per-task MCP tools merged into the base registry.
func (e *Executor) ExecuteWithOverlay(ctx context.Context, task models.Task, agent models.Agent, overlay *ToolOverlay, userProfile string) (Result, error) {
	var result Result
	_ = ctx

	model := e.LLMModel
	if agent.Model != "" {
		model = agent.Model
	}

	catalog := e.Tools
	runTool := e.ExecuteTool
	if overlay != nil {
		catalog = append(append([]ollama.Tool{}, e.Tools...), overlay.Catalog...)
		runTool = mergeToolExecutor(e.ExecuteTool, overlay.Run)
	}

	allowed := EffectiveAllowedTools(task, agent)
	tools := ollama.ToolsForAgent(allowed, catalog)
	runTool = filteredToolExecutor(runTool, allowed)

	messages := buildChatMessages(task, agent, userProfile, e.MemoryClient, e.Ollama)

	loopResult, err := ollama.RunAgentLoop(e.Ollama, ollama.AgentLoopOptions{
		Model:    model,
		Messages: messages,
		Tools:    tools,
		MaxSteps: e.MaxSteps,
		Stream:   e.StreamLLM,
		OnStream: e.OnStream,
		RunTool: func(name string, args map[string]interface{}) string {
			observation := runTool(name, ollama.FormatToolArguments(args))
			result.Observations = append(result.Observations, observation)
			return observation
		},
	})
	if err != nil {
		return result, err
	}

	result.FinalAnswer = loopResult.FinalAnswer
	result.Steps = loopResult.Steps
	result.ToolCalls = loopResult.ToolCalls
	if result.FinalAnswer == "" {
		return result, fmt.Errorf("agent loop finished without an answer")
	}
	return result, nil
}