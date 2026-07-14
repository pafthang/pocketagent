package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pafthang/pocketagent/pkgs/models"
)

type authResponse struct {
	Token  string                 `json:"token"`
	Record map[string]interface{} `json:"record"`
}

// AuthWithPassword logs in a user and returns a session.
func (c *Client) AuthWithPassword(identity, password string) (models.AuthSession, error) {
	url := fmt.Sprintf("%s/api/collections/%s/auth-with-password", c.BaseURL, UsersCollection)
	body, err := json.Marshal(map[string]string{
		"identity": identity,
		"password": password,
	})
	if err != nil {
		return models.AuthSession{}, err
	}

	resp, err := c.doPost(url, bytes.NewReader(body))
	if err != nil {
		return models.AuthSession{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.AuthSession{}, readAPIError(resp)
	}

	var result authResponse
	if err := decodeJSON(resp.Body, &result); err != nil {
		return models.AuthSession{}, err
	}

	return models.AuthSession{
		Token: result.Token,
		User:  userFromRecord(result.Record),
	}, nil
}

// AuthRefresh validates a token and returns a refreshed session.
func (c *Client) AuthRefresh(token string) (models.AuthSession, error) {
	url := fmt.Sprintf("%s/api/collections/%s/auth-refresh", c.BaseURL, UsersCollection)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return models.AuthSession{}, err
	}
	req.Header.Set("Authorization", normalizeAuthHeader(token))

	resp, err := c.http().Do(req)
	if err != nil {
		return models.AuthSession{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.AuthSession{}, readAPIError(resp)
	}

	var result authResponse
	if err := decodeJSON(resp.Body, &result); err != nil {
		return models.AuthSession{}, err
	}

	return models.AuthSession{
		Token: result.Token,
		User:  userFromRecord(result.Record),
	}, nil
}

// AuthSuperuser logs in as PocketBase superuser for server-side operations.
func (c *Client) AuthSuperuser(identity, password string) (string, error) {
	url := fmt.Sprintf("%s/api/collections/_superusers/auth-with-password", c.BaseURL)
	body, err := json.Marshal(map[string]string{
		"identity": identity,
		"password": password,
	})
	if err != nil {
		return "", err
	}

	resp, err := c.doPost(url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", readAPIError(resp)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result authResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}
	if result.Token == "" {
		return "", fmt.Errorf("empty superuser token")
	}
	return result.Token, nil
}

// RegisterUser creates a new auth user.
func (c *Client) RegisterUser(email, password string) (models.AuthUser, error) {
	record, err := c.CreateRecord(UsersCollection, map[string]interface{}{
		"email":           email,
		"password":        password,
		"passwordConfirm": password,
	})
	if err != nil {
		return models.AuthUser{}, err
	}
	return userFromRecord(record), nil
}

func userFromRecord(record map[string]interface{}) models.AuthUser {
	return models.AuthUser{
		ID:       stringField(record, "id"),
		Email:    stringField(record, "email"),
		Verified: boolField(record, "verified"),
	}
}

func normalizeAuthHeader(token string) string {
	if token == "" {
		return ""
	}
	lower := token
	if len(token) > 7 && (token[:7] == "Bearer " || token[:7] == "bearer ") {
		return token[7:]
	}
	return lower
}
