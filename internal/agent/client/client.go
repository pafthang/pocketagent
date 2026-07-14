package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	apimw "github.com/pafthang/pocketagent/pkgs/middle"
	"github.com/pafthang/pocketagent/pkgs/models"
)

// Client is an HTTP client for the agent service API.
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

// New creates an agent service client.
func New(baseURL string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTP:    &http.Client{},
	}
}

func (c *Client) http() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return http.DefaultClient
}

// Proxy forwards a request to the agent service and returns the upstream response.
func (c *Client) Proxy(method, path, token, spaceID string, body []byte) (*http.Response, error) {
	var reader io.Reader
	if len(body) > 0 {
		reader = bytes.NewReader(body)
	}
	return c.do(method, path, token, spaceID, reader)
}

// GetAgent fetches a single agent by ID.
func (c *Client) GetAgent(token, spaceID, id string) (models.Agent, error) {
	resp, err := c.do(http.MethodGet, "/agents/"+id, token, spaceID, nil)
	if err != nil {
		return models.Agent{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return models.Agent{}, readError(resp)
	}
	var agent models.Agent
	if err := json.NewDecoder(resp.Body).Decode(&agent); err != nil {
		return models.Agent{}, err
	}
	return agent, nil
}

// RuntimeConfig is the compiled agent configuration for workers and tooling.
type RuntimeConfig struct {
	ID           string               `json:"id"`
	SpaceID      string               `json:"space_id"`
	Name         string               `json:"name"`
	Model        string               `json:"model"`
	Tools        []string             `json:"tools"`
	SystemPrompt string               `json:"system_prompt"`
	Identity     models.IdentityFiles `json:"identity"`
}

// GetRuntimeConfig returns compiled prompt and tool configuration for an agent.
func (c *Client) GetRuntimeConfig(token, spaceID, id string) (RuntimeConfig, error) {
	resp, err := c.do(http.MethodGet, "/agents/"+id+"/runtime-config", token, spaceID, nil)
	if err != nil {
		return RuntimeConfig{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return RuntimeConfig{}, readError(resp)
	}
	var cfg RuntimeConfig
	if err := json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
		return RuntimeConfig{}, err
	}
	return cfg, nil
}

func (c *Client) do(method, path, token, spaceID string, body io.Reader) (*http.Response, error) {
	url := c.BaseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", normalizeBearer(token))
	}
	if spaceID != "" {
		req.Header.Set(apimw.HeaderSpaceID, spaceID)
	}
	return c.http().Do(req)
}

func normalizeBearer(token string) string {
	token = strings.TrimSpace(token)
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		return token
	}
	return "Bearer " + token
}

func readError(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if len(body) == 0 {
		return fmt.Errorf("agent: HTTP %d", resp.StatusCode)
	}
	var errResp struct {
		Error string `json:"error"`
	}
	if json.Unmarshal(body, &errResp) == nil && errResp.Error != "" {
		return fmt.Errorf("agent: %s", errResp.Error)
	}
	return fmt.Errorf("agent: HTTP %d: %s", resp.StatusCode, string(body))
}
