package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pafthang/pocketagent/pkgs/common"
)

// GenerateRequest configures a /api/generate call.
type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format,omitempty"`
	Tools  []Tool `json:"tools,omitempty"`
}

// Generate calls Ollama /api/generate.
func (c *Client) Generate(req GenerateRequest) (string, error) {
	var response string

	err := c.callWithResilience(context.Background(), func() error {
		text, err := c.generateOnce(req)
		if err != nil {
			return err
		}
		response = text
		return nil
	})

	return response, err
}

func (c *Client) generateOnce(req GenerateRequest) (string, error) {
	url := fmt.Sprintf("%s/api/generate", c.BaseURL)
	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := c.http().Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", common.NewHTTPStatusError(resp.StatusCode, string(respBody))
	}

	var result struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return string(respBody), nil
	}
	if result.Response != "" {
		return result.Response, nil
	}

	return string(respBody), nil
}