package client

import (
	"fmt"

	"github.com/pafthang/pocketagent/pkgs/models"
)



// FindUserByEmail returns a user by email if it exists.
func (c *Client) FindUserByEmail(email string) (models.AuthUser, bool, error) {
	filter := fmt.Sprintf("email = %q", email)
	records, _, err := c.ListRecordsOpts(UsersCollection, ListOptions{Page: 1, PerPage: 1, Filter: filter})
	if err != nil {
		return models.AuthUser{}, false, err
	}
	if len(records) == 0 {
		return models.AuthUser{}, false, nil
	}
	return userFromRecord(records[0]), true, nil
}

// GetUser returns a user by ID.
func (c *Client) GetUser(id string) (models.AuthUser, error) {
	record, err := c.GetRecord(UsersCollection, id)
	if err != nil {
		return models.AuthUser{}, err
	}
	return userFromRecord(record), nil
}

// SetUserVerified marks a user's email as verified.
func (c *Client) SetUserVerified(userID string, verified bool) error {
	_, err := c.UpdateRecord(UsersCollection, userID, map[string]interface{}{"verified": verified})
	return err
}

// CreateEmailVerification stores a pending verification token.
func (c *Client) CreateEmailVerification(userID, email, tokenHash, expiresAt string) error {
	_, err := c.CreateRecord(EmailVerificationsCollection, map[string]interface{}{
		"user_id":    userID,
		"email":      email,
		"token_hash": tokenHash,
		"status":     models.VerificationPending,
		"expires_at": expiresAt,
	})
	return err
}

// FindEmailVerificationByTokenHash returns a pending verification by token hash.
func (c *Client) FindEmailVerificationByTokenHash(tokenHash string) (map[string]interface{}, error) {
	filter := fmt.Sprintf("token_hash = %q && status = %q", tokenHash, models.VerificationPending)
	records, _, err := c.ListRecordsOpts(EmailVerificationsCollection, ListOptions{Page: 1, PerPage: 1, Filter: filter})
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, &APIError{StatusCode: 404, Message: "verification not found"}
	}
	return records[0], nil
}

// MarkEmailVerificationDone marks a verification record as used.
func (c *Client) MarkEmailVerificationDone(id string) error {
	_, err := c.UpdateRecord(EmailVerificationsCollection, id, map[string]interface{}{
		"status": models.VerificationDone,
	})
	return err
}
