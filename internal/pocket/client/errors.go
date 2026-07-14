package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// APIError is returned when PocketBase responds with a non-success status.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("pocketbase: HTTP %d", e.StatusCode)
}

type apiError struct {
	Message string         `json:"message"`
	Data    map[string]any `json:"data"`
}

func readAPIError(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))

	err := &APIError{StatusCode: resp.StatusCode}
	if len(body) == 0 {
		return err
	}

	var errResp apiError
	if json.Unmarshal(body, &errResp) == nil && errResp.Message != "" {
		err.Message = errResp.Message
		return err
	}

	err.Message = string(body)
	return err
}