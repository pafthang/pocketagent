package probe

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
)

type connector interface {
	listTools() ([]mcpTool, error)
	close() error
}

type mcpTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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

type stdioClient struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	nextID atomic.Uint64
	mu     sync.Mutex
}

type httpClient struct {
	server     ServerConfig
	httpClient *http.Client
	nextID     atomic.Uint64
	mu         sync.Mutex
}

func newStdioClient(cfg ServerConfig) (*stdioClient, error) {
	cmd := exec.Command(cfg.Command, cfg.Args...)
	if len(cfg.Env) > 0 {
		cmd.Env = append(os.Environ(), envPairs(cfg.Env)...)
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

	client := &stdioClient{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewReader(stdoutPipe),
	}
	if _, err := client.call("initialize", map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo":      map[string]string{"name": "pocketagent", "version": "1.0.0"},
	}); err != nil {
		_ = client.close()
		return nil, err
	}
	_ = client.notify("notifications/initialized", nil)
	return client, nil
}

func newHTTPClient(cfg ServerConfig) (*httpClient, error) {
	url := strings.TrimSpace(cfg.URL)
	if url == "" {
		return nil, fmt.Errorf("url is required for http transport")
	}
	client := &httpClient{
		server:     cfg,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
	if _, err := client.call("initialize", map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo":      map[string]string{"name": "pocketagent", "version": "1.0.0"},
	}); err != nil {
		return nil, err
	}
	return client, nil
}

func (c *stdioClient) listTools() ([]mcpTool, error) {
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

func (c *httpClient) listTools() ([]mcpTool, error) {
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

func (c *stdioClient) call(method string, params interface{}) (json.RawMessage, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := c.nextID.Add(1)
	req := mcpRPCRequest{JSONRPC: "2.0", ID: id, Method: method, Params: params}
	if err := c.write(req); err != nil {
		return nil, err
	}
	return c.readResponse(id)
}

func (c *stdioClient) notify(method string, params interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.write(map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": params})
}

func (c *stdioClient) write(payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = c.stdin.Write(data)
	return err
}

func (c *stdioClient) readResponse(expectedID uint64) (json.RawMessage, error) {
	for {
		line, err := c.stdout.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		var resp mcpRPCResponse
		if err := json.Unmarshal(line, &resp); err != nil {
			continue
		}
		if resp.ID == 0 || resp.ID != expectedID {
			continue
		}
		if resp.Error != nil {
			return nil, fmt.Errorf("mcp rpc %d: %s", resp.Error.Code, resp.Error.Message)
		}
		return resp.Result, nil
	}
}

func (c *stdioClient) close() error {
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

func (c *httpClient) call(method string, params interface{}) (json.RawMessage, error) {
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

func (c *httpClient) close() error {
	return nil
}

func envPairs(env map[string]string) []string {
	out := make([]string, 0, len(env))
	for key, value := range env {
		out = append(out, key+"="+value)
	}
	return out
}