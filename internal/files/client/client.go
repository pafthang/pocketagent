package client

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	apimw "github.com/pafthang/pocketagent/pkgs/middle"
)

// Client is an HTTP client for the files service API.
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

// New creates a files service client.
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

// Proxy forwards a request to the files service and returns the upstream response.
func (c *Client) Proxy(method, path, token, spaceID string, body []byte, contentType string) (*http.Response, error) {
	var reader io.Reader
	if len(body) > 0 {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, c.BaseURL+path, reader)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	} else if len(body) > 0 {
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