package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pafthang/pocketagent/pkgs/models"
)

// Client is an HTTP client for the space service API.
type Client struct {
	BaseURL string
	HTTP    *http.Client

	authorizeCache *authorizeCache
}

// New creates a space service client.
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

// Refresh validates a token and returns session info.
func (c *Client) Refresh(token string) (models.AuthSession, error) {
	resp, err := c.doJSON(http.MethodPost, "/auth/refresh", token, nil)
	if err != nil {
		return models.AuthSession{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.AuthSession{}, readError(resp)
	}

	var session models.AuthSession
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return models.AuthSession{}, err
	}
	return session, nil
}

// Authorize checks whether the caller may perform an action in a space.
func (c *Client) Authorize(token, spaceID, action string) (models.AuthorizeResponse, error) {
	body, err := json.Marshal(models.AuthorizeRequest{
		SpaceID: spaceID,
		Action:  action,
	})
	if err != nil {
		return models.AuthorizeResponse{}, err
	}

	resp, err := c.doJSON(http.MethodPost, "/authorize", token, bytes.NewReader(body))
	if err != nil {
		return models.AuthorizeResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.AuthorizeResponse{}, readError(resp)
	}

	var result models.AuthorizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.AuthorizeResponse{}, err
	}
	return result, nil
}

// Proxy forwards a request to the space service and returns the upstream response.
func (c *Client) Proxy(method, path, token string, body []byte) (*http.Response, error) {
	var reader io.Reader
	if len(body) > 0 {
		reader = bytes.NewReader(body)
	}
	return c.doJSON(method, path, token, reader)
}

func (c *Client) doJSON(method, path, token string, body io.Reader) (*http.Response, error) {
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
		return fmt.Errorf("space: HTTP %d", resp.StatusCode)
	}
	var errResp struct {
		Error string `json:"error"`
	}
	if json.Unmarshal(body, &errResp) == nil && errResp.Error != "" {
		return fmt.Errorf("space: %s", errResp.Error)
	}
	return fmt.Errorf("space: HTTP %d: %s", resp.StatusCode, string(body))
}
