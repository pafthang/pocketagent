package ollama

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Client for Ollama API
type Client struct {
	BaseURL string
}

func NewClient(url string) *Client {
	return &Client{BaseURL: url}
}

// Generate simple request to Ollama
type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

func (c *Client) Generate(req GenerateRequest) (string, error) {
	// TODO: real implementation
	return "Ollama response placeholder", nil
}
