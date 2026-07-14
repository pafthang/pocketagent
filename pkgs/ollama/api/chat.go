package api

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pafthang/pocketagent/pkgs/common"
)

// ChatMessage is a single message in a chat conversation.
type ChatMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	Thinking  string     `json:"thinking,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
	ToolName  string     `json:"tool_name,omitempty"`
}

// ToolCall is a native tool invocation from the model.
type ToolCall struct {
	Type     string       `json:"type"`
	Function ToolCallFunc `json:"function"`
}

// ToolCallFunc holds the function name and JSON arguments.
type ToolCallFunc struct {
	Index     int                    `json:"index,omitempty"`
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// ChatRequest configures a /api/chat call.
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Tools    []Tool        `json:"tools,omitempty"`
	Stream   bool          `json:"stream"`
	Format   string        `json:"format,omitempty"`
	Think    bool          `json:"think,omitempty"`
}

// ChatResponse is a non-streaming chat result.
type ChatResponse struct {
	Message ChatMessage `json:"message"`
	Done    bool        `json:"done"`
}

// ChatStreamEvent is one NDJSON chunk from a streaming chat response.
type ChatStreamEvent struct {
	Message ChatMessage `json:"message"`
	Done    bool        `json:"done"`
}

// StreamChunkHandler receives incremental model output.
type StreamChunkHandler func(event ChatStreamEvent)

// Chat calls Ollama /api/chat (non-streaming).
func (c *Client) Chat(req ChatRequest) (ChatResponse, error) {
	var result ChatResponse

	err := c.callWithResilience(context.Background(), func() error {
		resp, err := c.chatOnce(req)
		if err != nil {
			return err
		}
		result = resp
		return nil
	})

	return result, err
}

func (c *Client) chatOnce(req ChatRequest) (ChatResponse, error) {
	req.Stream = false
	url := fmt.Sprintf("%s/api/chat", c.BaseURL)
	body, err := json.Marshal(req)
	if err != nil {
		return ChatResponse{}, err
	}

	resp, err := c.http().Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return ChatResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return ChatResponse{}, common.NewHTTPStatusError(resp.StatusCode, fmt.Sprintf("ollama chat: HTTP %d: %s", resp.StatusCode, string(respBody)))
	}

	var result ChatResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return ChatResponse{}, fmt.Errorf("decode chat response: %w", err)
	}
	return result, nil
}

// ChatStream calls Ollama /api/chat with streaming enabled and returns the accumulated assistant message.
func (c *Client) ChatStream(req ChatRequest, onChunk StreamChunkHandler) (ChatMessage, error) {
	var accum ChatMessage

	err := c.callWithResilience(context.Background(), func() error {
		msg, err := c.chatStreamOnce(req, onChunk)
		if err != nil {
			return err
		}
		accum = msg
		return nil
	})

	return accum, err
}

func (c *Client) chatStreamOnce(req ChatRequest, onChunk StreamChunkHandler) (ChatMessage, error) {
	req.Stream = true
	url := fmt.Sprintf("%s/api/chat", c.BaseURL)
	body, err := json.Marshal(req)
	if err != nil {
		return ChatMessage{}, err
	}

	resp, err := c.http().Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return ChatMessage{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return ChatMessage{}, common.NewHTTPStatusError(resp.StatusCode, fmt.Sprintf("ollama chat stream: HTTP %d: %s", resp.StatusCode, string(respBody)))
	}

	accum := ChatMessage{Role: "assistant"}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var event ChatStreamEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}
		if event.Message.Thinking != "" {
			accum.Thinking += event.Message.Thinking
		}
		if event.Message.Content != "" {
			accum.Content += event.Message.Content
		}
		if len(event.Message.ToolCalls) > 0 {
			accum.ToolCalls = mergeToolCalls(accum.ToolCalls, event.Message.ToolCalls)
		}
		if onChunk != nil {
			onChunk(event)
		}
		if event.Done {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return accum, err
	}
	return accum, nil
}

func mergeToolCalls(existing, delta []ToolCall) []ToolCall {
	if len(delta) == 0 {
		return existing
	}
	out := append([]ToolCall{}, existing...)
	for _, call := range delta {
		if call.Function.Name == "" && len(out) > 0 {
			last := &out[len(out)-1]
			if call.Function.Arguments != nil {
				if last.Function.Arguments == nil {
					last.Function.Arguments = map[string]interface{}{}
				}
				for k, v := range call.Function.Arguments {
					last.Function.Arguments[k] = v
				}
			}
			continue
		}
		out = append(out, call)
	}
	return out
}

// FormatToolArguments converts tool arguments to a string for tool executors.
func FormatToolArguments(arguments map[string]interface{}) string {
	if len(arguments) == 0 {
		return ""
	}
	for _, key := range []string{"query", "url", "input", "q", "code"} {
		if v, ok := arguments[key]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	b, err := json.Marshal(arguments)
	if err != nil {
		return ""
	}
	return string(b)
}