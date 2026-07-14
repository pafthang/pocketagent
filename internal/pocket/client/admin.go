package client

import "fmt"

// NewServiceClient creates a PocketBase client for backend services.
// When admin credentials are set, authenticates as superuser (required for locked collections).
func NewServiceClient(baseURL, adminEmail, adminPassword string) (*Client, error) {
	c := New(baseURL)
	if adminEmail == "" || adminPassword == "" {
		return nil, fmt.Errorf("pocketbase admin credentials required (POCKETBASE_ADMIN_EMAIL/PASSWORD)")
	}

	token, err := c.AuthSuperuser(adminEmail, adminPassword)
	if err != nil {
		return nil, fmt.Errorf("pocketbase superuser auth: %w", err)
	}
	c.Token = token
	return c, nil
}