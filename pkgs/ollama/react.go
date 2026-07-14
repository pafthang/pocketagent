package ollama

import "github.com/pafthang/pocketagent/pkgs/ollama/agent"

type ToolRunner = agent.ToolRunner
type AgentLoopOptions = agent.LoopOptions
type AgentLoopResult = agent.LoopResult

// RunAgentLoop drives Ollama /api/chat with native tool calling.
func RunAgentLoop(client *Client, opts AgentLoopOptions) (AgentLoopResult, error) {
	return agent.RunLoop(client, opts)
}