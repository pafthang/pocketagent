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

const defaultEmbedModel = "nomic-embed-text"

type embedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type embedResponse struct {
	Embedding []float64 `json:"embedding"`
}

// Embed generates a vector embedding for text via Ollama /api/embeddings.
func (c *Client) Embed(ctx context.Context, text string) ([]float32, error) {
	var embedding []float32

	err := c.callWithResilience(ctx, func() error {
		vec, err := c.embedOnce(ctx, text)
		if err != nil {
			return err
		}
		embedding = vec
		return nil
	})

	return embedding, err
}

func (c *Client) embedOnce(ctx context.Context, text string) ([]float32, error) {
	model := c.EmbedModel
	if model == "" {
		model = defaultEmbedModel
	}

	body, err := json.Marshal(embedRequest{Model: model, Prompt: text})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, common.NewHTTPStatusError(resp.StatusCode, fmt.Sprintf("ollama embeddings: HTTP %d: %s", resp.StatusCode, string(respBody)))
	}

	var result embedResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}
	if len(result.Embedding) == 0 {
		return nil, fmt.Errorf("ollama embeddings: empty vector")
	}

	embedding := make([]float32, len(result.Embedding))
	for i, v := range result.Embedding {
		embedding[i] = float32(v)
	}

	return embedding, nil
}