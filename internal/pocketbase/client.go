package pocketbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client for PocketBase

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTP:    &http.Client{},
	}
}

// AuthAdmin authenticates as admin
func (c *Client) AuthAdmin(email, password string) (string, error) {
	// TODO: implement admin auth
	return "", nil
}

// CreateRecord creates a record in collection
type CreateRecordRequest struct {
	Collection string                 `json:"-"`
	Data       map[string]interface{} `json:"data"`
}

func (c *Client) CreateRecord(collection string, data map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/collections/%s/records", c.BaseURL, collection)
	body, _ := json.Marshal(data)

	resp, err := c.HTTP.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}
