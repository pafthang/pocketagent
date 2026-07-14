package tools

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pafthang/pocketagent/pkgs/ollama"
)

type mcpConnector interface {
	listTools() ([]mcpTool, error)
	callTool(toolName string, args map[string]interface{}) (string, error)
	close() error
}

type mcpClient struct {
	server MCPServerConfig
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	nextID atomic.Uint64
	mu     sync.Mutex
	tools  []mcpTool
	ready  bool
}

type mcpHTTPClient struct {
	server     MCPServerConfig
	httpClient *http.Client
	nextID     atomic.Uint64
	mu         sync.Mutex
}

type mcpTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

type mcpRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      uint64      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type mcpRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      uint64          `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *mcpRPCError    `json:"error,omitempty"`
}

type mcpRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func connectMCPServers(servers []MCPServerConfig) ([]mcpToolBinding, []ollama.Tool, []mcpConnector) {
	bindings := make([]mcpToolBinding, 0)
	catalog := make([]ollama.Tool, 0)
	clients := make([]mcpConnector, 0)

	for _, server := range servers {
		if !serverEnabled(server) {
			continue
		}
		transport := strings.ToLower(strings.TrimSpace(server.Transport))
		if transport == "" {
			transport = "stdio"
		}

		var client mcpConnector
		var err error
		switch transport {
		case "http":
			client, err = newMCPHTTPClient(server)
		case "stdio":
			if strings.TrimSpace(server.Command) == "" {
				continue
			}
			client, err = newMCPStdioClient(server)
		default:
			continue
		}
		if strings.TrimSpace(server.Name) == "" {
			continue
		}
		if err != nil {
			continue
		}
		tools, err := client.listTools()
		if err != nil {
			_ = client.close()
			continue
		}
		clients = append(clients, client)
		for _, tool := range tools {
			name := mcpToolName(server.Name, tool.Name)
			bindings = append(bindings, mcpToolBinding{
				name:   name,
				client: client,
				tool:   tool.Name,
			})
			catalog = append(catalog, ollama.Tool{
				Type: "function",
				Function: ollama.ToolFunction{
					Name:        name,
					Description: fmt.Sprintf("MCP/%s: %s", server.Name, tool.Description),
					Parameters:  normalizeMCPSchema(tool.InputSchema),
				},
			})
		}
	}

	return bindings, catalog, clients
}

type mcpToolBinding struct {
	name   string
	client mcpConnector
	tool   string
}

func serverEnabled(server MCPServerConfig) bool {
	if server.Name == "" {
		return false
	}
	// Configs without explicit enabled field default to enabled.
	if !server.Enabled && server.Transport != "" {
		return false
	}
	return true
}

func newMCPStdioClient(server MCPServerConfig) (*mcpClient, error) {
	cmd := exec.Command(server.Command, server.Args...)
	if len(server.Env) > 0 {
		cmd.Env = append(os.Environ(), envPairs(server.Env)...)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	go io.Copy(io.Discard, stderr)

	client := &mcpClient{
		server: server,
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewReader(stdoutPipe),
	}

	if err := client.initialize(); err != nil {
		_ = client.close()
		return nil, err
	}
	return client, nil
}

func (c *mcpClient) initialize() error {
	_, err := c.call("initialize", map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo": map[string]string{
			"name":    "pocketagent",
			"version": "1.0.0",
		},
	})
	if err != nil {
		return err
	}
	_ = c.notify("notifications/initialized", nil)
	c.ready = true
	return nil
}

func newMCPHTTPClient(server MCPServerConfig) (*mcpHTTPClient, error) {
	url := strings.TrimSpace(server.URL)
	if url == "" {
		return nil, fmt.Errorf("url is required for http transport")
	}
	client := &mcpHTTPClient{
		server: server,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	if _, err := client.initialize(); err != nil {
		return nil, err
	}
	return client, nil
}

func (c *mcpHTTPClient) initialize() (json.RawMessage, error) {
	return c.call("initialize", map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo": map[string]string{
			"name":    "pocketagent",
			"version": "1.0.0",
		},
	})
}

func (c *mcpHTTPClient) listTools() ([]mcpTool, error) {
	raw, err := c.call("tools/list", map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	var result struct {
		Tools []mcpTool `json:"tools"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}
	return result.Tools, nil
}

func (c *mcpHTTPClient) callTool(toolName string, args map[string]interface{}) (string, error) {
	raw, err := c.call("tools/call", map[string]interface{}{
		"name":      toolName,
		"arguments": args,
	})
	if err != nil {
		return "", err
	}
	return formatToolResult(raw)
}

func (c *mcpHTTPClient) call(method string, params interface{}) (json.RawMessage, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := c.nextID.Add(1)
	req := mcpRPCRequest{JSONRPC: "2.0", ID: id, Method: method, Params: params}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, c.server.URL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	for key, value := range c.server.Env {
		if strings.HasPrefix(strings.ToLower(key), "header:") {
			httpReq.Header.Set(strings.TrimPrefix(key, "header:"), value)
		}
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("mcp http %d: %s", resp.StatusCode, string(respBody))
	}

	var rpcResp mcpRPCResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return nil, err
	}
	if rpcResp.Error != nil {
		return nil, fmt.Errorf("mcp rpc %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}
	return rpcResp.Result, nil
}

func (c *mcpHTTPClient) close() error {
	return nil
}

func (c *mcpClient) listTools() ([]mcpTool, error) {
	raw, err := c.call("tools/list", map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	var result struct {
		Tools []mcpTool `json:"tools"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}
	c.tools = result.Tools
	return result.Tools, nil
}

func (c *mcpClient) callTool(toolName string, args map[string]interface{}) (string, error) {
	raw, err := c.call("tools/call", map[string]interface{}{
		"name":      toolName,
		"arguments": args,
	})
	if err != nil {
		return "", err
	}

	return formatToolResult(raw)
}

func formatToolResult(raw json.RawMessage) (string, error) {
	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		IsError bool `json:"isError"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return string(raw), nil
	}

	var parts []string
	for _, item := range result.Content {
		if item.Text != "" {
			parts = append(parts, item.Text)
		}
	}
	out := strings.Join(parts, "\n")
	if out == "" {
		out = string(raw)
	}
	if result.IsError {
		return out, fmt.Errorf("mcp tool error")
	}
	return out, nil
}

func envPairs(env map[string]string) []string {
	out := make([]string, 0, len(env))
	for key, value := range env {
		out = append(out, key+"="+value)
	}
	return out
}

func (c *mcpClient) call(method string, params interface{}) (json.RawMessage, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := c.nextID.Add(1)
	req := mcpRPCRequest{JSONRPC: "2.0", ID: id, Method: method, Params: params}
	if err := c.write(req); err != nil {
		return nil, err
	}
	return c.readResponse(id)
}

func (c *mcpClient) notify(method string, params interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
	}
	return c.write(req)
}

func (c *mcpClient) write(payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = c.stdin.Write(data)
	return err
}

func (c *mcpClient) readResponse(expectedID uint64) (json.RawMessage, error) {
	for {
		line, err := c.stdout.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		line = bytesTrimSpace(line)
		if len(line) == 0 {
			continue
		}

		var resp mcpRPCResponse
		if err := json.Unmarshal(line, &resp); err != nil {
			continue
		}
		if resp.ID == 0 {
			continue
		}
		if resp.ID != expectedID {
			continue
		}
		if resp.Error != nil {
			return nil, fmt.Errorf("mcp rpc %d: %s", resp.Error.Code, resp.Error.Message)
		}
		return resp.Result, nil
	}
}

func (c *mcpClient) close() error {
	if c.stdin != nil {
		_ = c.stdin.Close()
	}
	if c.cmd != nil && c.cmd.Process != nil {
		_ = c.cmd.Process.Kill()
	}
	if c.cmd != nil {
		_ = c.cmd.Wait()
	}
	return nil
}

// PublicToolName returns the agent registry name for an MCP tool.
func PublicToolName(serverName, toolName string) string {
	return mcpToolName(serverName, toolName)
}

func mcpToolName(serverName, toolName string) string {
	safeServer := sanitizeName(serverName)
	safeTool := sanitizeName(toolName)
	return "mcp__" + safeServer + "__" + safeTool
}

func sanitizeName(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-', r == '_':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	out := b.String()
	if out == "" {
		return "tool"
	}
	return out
}

func normalizeMCPSchema(schema map[string]interface{}) map[string]interface{} {
	if schema == nil {
		return map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		}
	}
	if _, ok := schema["type"]; !ok {
		schema["type"] = "object"
	}
	return schema
}

func bytesTrimSpace(b []byte) []byte {
	return []byte(strings.TrimSpace(string(b)))
}
