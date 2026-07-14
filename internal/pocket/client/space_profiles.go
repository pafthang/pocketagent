package client

import (
	"fmt"

	"github.com/pafthang/pocketagent/pkgs/models"
)



// GetSpaceProfile returns the profile for a user in a space, or empty content if missing.
func (c *Client) GetSpaceProfile(spaceID, userID string) (models.SpaceProfile, error) {
	filter := fmt.Sprintf("space_id = %q && user_id = %q", spaceID, userID)
	records, _, err := c.ListRecordsOpts(SpaceProfilesCollection, ListOptions{Page: 1, PerPage: 1, Filter: filter})
	if err != nil {
		return models.SpaceProfile{}, err
	}
	if len(records) == 0 {
		return models.SpaceProfile{SpaceID: spaceID, UserID: userID}, nil
	}
	return spaceProfileFromRecord(records[0]), nil
}

// UpsertSpaceProfile creates or updates a user profile in a space.
func (c *Client) UpsertSpaceProfile(profile models.SpaceProfile) (models.SpaceProfile, error) {
	existing, err := c.GetSpaceProfile(profile.SpaceID, profile.UserID)
	if err != nil {
		return models.SpaceProfile{}, err
	}
	data := map[string]interface{}{
		"space_id": profile.SpaceID,
		"user_id":  profile.UserID,
		"content":  profile.Content,
	}
	if existing.ID == "" {
		record, err := c.CreateRecord(SpaceProfilesCollection, data)
		if err != nil {
			return models.SpaceProfile{}, err
		}
		return spaceProfileFromRecord(record), nil
	}
	record, err := c.UpdateRecord(SpaceProfilesCollection, existing.ID, data)
	if err != nil {
		return models.SpaceProfile{}, err
	}
	return spaceProfileFromRecord(record), nil
}

func spaceProfileFromRecord(record map[string]interface{}) models.SpaceProfile {
	return models.SpaceProfile{
		ID:        stringField(record, "id"),
		SpaceID:   stringField(record, "space_id"),
		UserID:    stringField(record, "user_id"),
		Content:   stringField(record, "content"),
		CreatedAt: stringField(record, "created"),
		UpdatedAt: stringField(record, "updated"),
	}
}
