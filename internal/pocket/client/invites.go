package client

import (
	"fmt"

	"github.com/pafthang/pocketagent/pkgs/models"
)



// CreateSpaceInvite stores a pending invite.
func (c *Client) CreateSpaceInvite(invite models.SpaceInvite, tokenHash string) (models.SpaceInvite, error) {
	record, err := c.CreateRecord(SpaceInvitesCollection, map[string]interface{}{
		"space_id":   invite.SpaceID,
		"email":      invite.Email,
		"role":       invite.Role,
		"token_hash": tokenHash,
		"invited_by": invite.InvitedBy,
		"status":     models.InvitePending,
		"expires_at": invite.ExpiresAt,
	})
	if err != nil {
		return models.SpaceInvite{}, err
	}
	return inviteFromRecord(record), nil
}

// GetSpaceInvite returns an invite by ID.
func (c *Client) GetSpaceInvite(id string) (models.SpaceInvite, error) {
	record, err := c.GetRecord(SpaceInvitesCollection, id)
	if err != nil {
		return models.SpaceInvite{}, err
	}
	return inviteFromRecord(record), nil
}

// FindSpaceInviteByTokenHash returns a pending invite matching the token hash.
func (c *Client) FindSpaceInviteByTokenHash(tokenHash string) (models.SpaceInvite, error) {
	filter := fmt.Sprintf("token_hash = %q && status = %q", tokenHash, models.InvitePending)
	records, _, err := c.ListRecordsOpts(SpaceInvitesCollection, ListOptions{Page: 1, PerPage: 1, Filter: filter})
	if err != nil {
		return models.SpaceInvite{}, err
	}
	if len(records) == 0 {
		return models.SpaceInvite{}, &APIError{StatusCode: 404, Message: "invite not found"}
	}
	return inviteFromRecord(records[0]), nil
}

// ListSpaceInvites lists invites for a space.
func (c *Client) ListSpaceInvites(spaceID string, opts ListOptions) ([]models.SpaceInvite, int, error) {
	filter := fmt.Sprintf("space_id = %q", spaceID)
	if opts.Filter != "" {
		filter = filter + " && (" + opts.Filter + ")"
	}
	opts.Filter = filter
	records, total, err := c.ListRecordsOpts(SpaceInvitesCollection, opts)
	if err != nil {
		return nil, 0, err
	}
	invites := make([]models.SpaceInvite, 0, len(records))
	for _, record := range records {
		invites = append(invites, inviteFromRecord(record))
	}
	return invites, total, nil
}

// UpdateSpaceInvite patches an invite.
func (c *Client) UpdateSpaceInvite(id string, data map[string]interface{}) (models.SpaceInvite, error) {
	record, err := c.UpdateRecord(SpaceInvitesCollection, id, data)
	if err != nil {
		return models.SpaceInvite{}, err
	}
	return inviteFromRecord(record), nil
}

// DeleteSpaceInvite removes an invite.
func (c *Client) DeleteSpaceInvite(id string) error {
	return c.DeleteRecord(SpaceInvitesCollection, id)
}

func inviteFromRecord(record map[string]interface{}) models.SpaceInvite {
	return models.SpaceInvite{
		ID:        stringField(record, "id"),
		SpaceID:   stringField(record, "space_id"),
		Email:     stringField(record, "email"),
		Role:      stringField(record, "role"),
		Status:    stringField(record, "status"),
		InvitedBy: stringField(record, "invited_by"),
		ExpiresAt: stringField(record, "expires_at"),
		CreatedAt: stringField(record, "created"),
		UpdatedAt: stringField(record, "updated"),
	}
}
