package react

import (
	"context"
	"strings"

	"github.com/pafthang/pocketagent/internal/agent/identity"
	memoclient "github.com/pafthang/pocketagent/internal/memo/client"
	"github.com/pafthang/pocketagent/pkgs/models"
	"github.com/pafthang/pocketagent/pkgs/ollama"
)

func buildChatMessages(task models.Task, agent models.Agent, userProfile string, memory *memoclient.Client, ollamaClient *ollama.Client) []ollama.ChatMessage {
	var messages []ollama.ChatMessage

	prompt := identity.CompileAgentPrompt(identity.FromAgent(agent))
	if prompt != "" {
		messages = append(messages, ollama.ChatMessage{
			Role:    "system",
			Content: prompt,
		})
	} else {
		messages = append(messages, ollama.ChatMessage{
			Role:    "system",
			Content: "You are a helpful agent with access to tools. Use native tool_calls when available.",
		})
	}

	if profile := identity.FormatUserProfile(userProfile); profile != "" {
		messages = append(messages, ollama.ChatMessage{
			Role:    "system",
			Content: profile,
		})
	}

	if memory != nil && ollamaClient != nil && task.SpaceID != "" {
		relevantDocs, _ := memory.QueryScoped(context.Background(), ollamaClient, task.SpaceID, task.Prompt, 3)
		if len(relevantDocs) > 0 {
			var b strings.Builder
			b.WriteString("Relevant context from memory:\n")
			for _, doc := range relevantDocs {
				b.WriteString("- ")
				b.WriteString(doc)
				b.WriteString("\n")
			}
			messages = append(messages, ollama.ChatMessage{Role: "system", Content: b.String()})
		}
	}

	messages = append(messages, ollama.ChatMessage{
		Role:    "user",
		Content: task.Prompt,
	})
	return messages
}