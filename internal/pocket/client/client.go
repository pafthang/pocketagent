package client

import "net/http"

// Client is an HTTP client for the PocketBase API.
type Client struct {
	BaseURL string
	HTTP    *http.Client
	Token   string
}

// New creates a PocketBase API client.
func New(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTP:    &http.Client{},
	}
}