package client

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func (c *Client) doGet(targetURL string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, err
	}
	c.applyAuth(req)
	return c.http().Do(req)
}

func (c *Client) doPost(targetURL string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, targetURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.applyAuth(req)
	return c.http().Do(req)
}

func (c *Client) applyAuth(req *http.Request) {
	if c.Token == "" {
		return
	}
	token := c.Token
	if !strings.HasPrefix(strings.ToLower(token), "bearer ") {
		token = token
	}
	req.Header.Set("Authorization", token)
}

func (c *Client) http() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return http.DefaultClient
}

func decodeJSON(r io.Reader, dest any) error {
	return json.NewDecoder(r).Decode(dest)
}