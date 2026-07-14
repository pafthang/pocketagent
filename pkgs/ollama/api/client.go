package api

import "net/http"

// Client talks to a local or remote Ollama HTTP API.
type Client struct {
	BaseURL    string
	EmbedModel string
	HTTP       *http.Client
	Breaker    *Breaker
}

// New creates a client with default settings.
func New(url string) *Client {
	return &Client{
		BaseURL:    url,
		EmbedModel: defaultEmbedModel,
		HTTP:       &http.Client{},
	}
}

// NewConfigured creates a client with an optional embedding model override.
func NewConfigured(url, embedModel string) *Client {
	c := New(url)
	if embedModel != "" {
		c.EmbedModel = embedModel
	}
	return c
}

func (c *Client) http() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return http.DefaultClient
}