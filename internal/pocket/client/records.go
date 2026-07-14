package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type listResponse struct {
	Page       int                      `json:"page"`
	PerPage    int                      `json:"perPage"`
	TotalItems int                      `json:"totalItems"`
	TotalPages int                      `json:"totalPages"`
	Items      []map[string]interface{} `json:"items"`
}

// CreateRecord creates a record in a collection.
func (c *Client) CreateRecord(collection string, data map[string]interface{}) (map[string]interface{}, error) {
	return c.writeRecord(http.MethodPost, c.collectionURL(collection), data)
}

// GetRecord returns a single record by ID.
func (c *Client) GetRecord(collection, id string) (map[string]interface{}, error) {
	resp, err := c.doGet(recordURL(c.BaseURL, collection, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, readAPIError(resp)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// ListRecords returns records with optional pagination.
func (c *Client) ListRecords(collection string, page, perPage int) ([]map[string]interface{}, int, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 50
	}

	u, err := url.Parse(c.collectionURL(collection))
	if err != nil {
		return nil, 0, err
	}
	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	q.Set("perPage", strconv.Itoa(perPage))
	u.RawQuery = q.Encode()

	resp, err := c.doGet(u.String())
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, readAPIError(resp)
	}

	var result listResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	return result.Items, result.TotalItems, nil
}

// UpdateRecord patches a record by ID.
func (c *Client) UpdateRecord(collection, id string, data map[string]interface{}) (map[string]interface{}, error) {
	return c.writeRecord(http.MethodPatch, recordURL(c.BaseURL, collection, id), data)
}

// DeleteRecord removes a record by ID.
func (c *Client) DeleteRecord(collection, id string) error {
	req, err := http.NewRequest(http.MethodDelete, recordURL(c.BaseURL, collection, id), nil)
	if err != nil {
		return err
	}
	c.applyAuth(req)

	resp, err := c.http().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return readAPIError(resp)
	}
	return nil
}

func (c *Client) writeRecord(method, targetURL string, data map[string]interface{}) (map[string]interface{}, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, targetURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.applyAuth(req)

	resp, err := c.http().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, readAPIError(resp)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

func (c *Client) collectionURL(collection string) string {
	return fmt.Sprintf("%s/api/collections/%s/records", c.BaseURL, collection)
}

func recordURL(baseURL, collection, id string) string {
	return fmt.Sprintf("%s/api/collections/%s/records/%s", baseURL, collection, id)
}