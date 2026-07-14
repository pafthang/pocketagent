package ollama

import "github.com/pafthang/pocketagent/pkgs/ollama/api"

type ChatMessage = api.ChatMessage
type ToolCall = api.ToolCall
type ToolCallFunc = api.ToolCallFunc
type ChatRequest = api.ChatRequest
type ChatResponse = api.ChatResponse
type ChatStreamEvent = api.ChatStreamEvent
type StreamChunkHandler = api.StreamChunkHandler
type GenerateRequest = api.GenerateRequest

func FormatToolArguments(arguments map[string]interface{}) string {
	return api.FormatToolArguments(arguments)
}