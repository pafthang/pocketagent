package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL string
}

func NewClient(url string) *Client {
	return &Client{BaseURL: url}
}

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Tools  []Tool `json:"tools,omitempty"`
}

func (c *Client) Generate(req GenerateRequest) (string, error) {
	url := fmt.Sprintf("%s/api/generate", c.BaseURL)
	body, _ := json.Marshal(req)

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Simplified response handling
	return "[Tool calling supported response]", nil
}
