package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pafthang/pocketagent/pkgs/models"
)



// ListOptions configures a filtered list query.
type ListOptions struct {
	Page    int
	PerPage int
	Filter  string
}

// CreateSpace stores a new space.
func (c *Client) CreateSpace(space models.Space) (models.Space, error) {
	record, err := c.CreateRecord(SpacesCollection, spaceRecordData(space))
	if err != nil {
		return models.Space{}, err
	}
	return spaceFromRecord(record), nil
}

// GetSpace returns a space by ID.
func (c *Client) GetSpace(id string) (models.Space, error) {
	record, err := c.GetRecord(SpacesCollection, id)
	if err != nil {
		return models.Space{}, err
	}
	return spaceFromRecord(record), nil
}

// GetSpaceBySlug returns a space by slug.
func (c *Client) GetSpaceBySlug(slug string) (models.Space, error) {
	filter := fmt.Sprintf("slug = %q", slug)
	records, _, err := c.ListRecordsOpts(SpacesCollection, ListOptions{Page: 1, PerPage: 1, Filter: filter})
	if err != nil {
		return models.Space{}, err
	}
	if len(records) == 0 {
		return models.Space{}, &APIError{StatusCode: 404, Message: "space not found"}
	}
	return spaceFromRecord(records[0]), nil
}

// ListSpaces returns spaces with optional filter.
func (c *Client) ListSpaces(opts ListOptions) ([]models.Space, int, error) {
	records, total, err := c.ListRecordsOpts(SpacesCollection, opts)
	if err != nil {
		return nil, 0, err
	}
	spaces := make([]models.Space, 0, len(records))
	for _, record := range records {
		spaces = append(spaces, spaceFromRecord(record))
	}
	return spaces, total, nil
}

// UpdateSpace patches a space.
func (c *Client) UpdateSpace(id string, space models.Space) (models.Space, error) {
	record, err := c.UpdateRecord(SpacesCollection, id, spaceRecordData(space))
	if err != nil {
		return models.Space{}, err
	}
	return spaceFromRecord(record), nil
}

// DeleteSpace removes a space.
func (c *Client) DeleteSpace(id string) error {
	return c.DeleteRecord(SpacesCollection, id)
}

// CreateSpaceMember adds a user to a space.
func (c *Client) CreateSpaceMember(member models.SpaceMember) (models.SpaceMember, error) {
	record, err := c.CreateRecord(SpaceMembersCollection, spaceMemberRecordData(member))
	if err != nil {
		return models.SpaceMember{}, err
	}
	return spaceMemberFromRecord(record), nil
}

// GetSpaceMember returns a membership record.
func (c *Client) GetSpaceMember(id string) (models.SpaceMember, error) {
	record, err := c.GetRecord(SpaceMembersCollection, id)
	if err != nil {
		return models.SpaceMember{}, err
	}
	return spaceMemberFromRecord(record), nil
}

// ListSpaceMembers lists memberships with optional filter.
func (c *Client) ListSpaceMembers(opts ListOptions) ([]models.SpaceMember, int, error) {
	records, total, err := c.ListRecordsOpts(SpaceMembersCollection, opts)
	if err != nil {
		return nil, 0, err
	}
	members := make([]models.SpaceMember, 0, len(records))
	for _, record := range records {
		members = append(members, spaceMemberFromRecord(record))
	}
	return members, total, nil
}

// UpdateSpaceMember patches a membership.
func (c *Client) UpdateSpaceMember(id string, member models.SpaceMember) (models.SpaceMember, error) {
	record, err := c.UpdateRecord(SpaceMembersCollection, id, spaceMemberRecordData(member))
	if err != nil {
		return models.SpaceMember{}, err
	}
	return spaceMemberFromRecord(record), nil
}

// DeleteSpaceMember removes a membership.
func (c *Client) DeleteSpaceMember(id string) error {
	return c.DeleteRecord(SpaceMembersCollection, id)
}

// CreateTeam stores a team in a space.
func (c *Client) CreateTeam(team models.Team) (models.Team, error) {
	record, err := c.CreateRecord(TeamsCollection, teamRecordData(team))
	if err != nil {
		return models.Team{}, err
	}
	return teamFromRecord(record), nil
}

// GetTeam returns a team by ID.
func (c *Client) GetTeam(id string) (models.Team, error) {
	record, err := c.GetRecord(TeamsCollection, id)
	if err != nil {
		return models.Team{}, err
	}
	return teamFromRecord(record), nil
}

// ListTeams lists teams with optional filter.
func (c *Client) ListTeams(opts ListOptions) ([]models.Team, int, error) {
	records, total, err := c.ListRecordsOpts(TeamsCollection, opts)
	if err != nil {
		return nil, 0, err
	}
	teams := make([]models.Team, 0, len(records))
	for _, record := range records {
		teams = append(teams, teamFromRecord(record))
	}
	return teams, total, nil
}

// UpdateTeam patches a team.
func (c *Client) UpdateTeam(id string, team models.Team) (models.Team, error) {
	record, err := c.UpdateRecord(TeamsCollection, id, teamRecordData(team))
	if err != nil {
		return models.Team{}, err
	}
	return teamFromRecord(record), nil
}

// DeleteTeam removes a team.
func (c *Client) DeleteTeam(id string) error {
	return c.DeleteRecord(TeamsCollection, id)
}

// CreateTeamMember links a user or agent to a team.
func (c *Client) CreateTeamMember(member models.TeamMember) (models.TeamMember, error) {
	record, err := c.CreateRecord(TeamMembersCollection, teamMemberRecordData(member))
	if err != nil {
		return models.TeamMember{}, err
	}
	return teamMemberFromRecord(record), nil
}

// ListTeamMembers lists team members with optional filter.
func (c *Client) ListTeamMembers(opts ListOptions) ([]models.TeamMember, int, error) {
	records, total, err := c.ListRecordsOpts(TeamMembersCollection, opts)
	if err != nil {
		return nil, 0, err
	}
	members := make([]models.TeamMember, 0, len(records))
	for _, record := range records {
		members = append(members, teamMemberFromRecord(record))
	}
	return members, total, nil
}

// DeleteTeamMember removes a team member.
func (c *Client) DeleteTeamMember(id string) error {
	return c.DeleteRecord(TeamMembersCollection, id)
}

// ListRecordsOpts returns records with pagination and filter.
func (c *Client) ListRecordsOpts(collection string, opts ListOptions) ([]map[string]interface{}, int, error) {
	page := opts.Page
	if page <= 0 {
		page = 1
	}
	perPage := opts.PerPage
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
	if opts.Filter != "" {
		q.Set("filter", opts.Filter)
	}
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
	if err := decodeJSON(resp.Body, &result); err != nil {
		return nil, 0, err
	}
	return result.Items, result.TotalItems, nil
}

func spaceRecordData(space models.Space) map[string]interface{} {
	data := map[string]interface{}{
		"name": space.Name,
		"slug": space.Slug,
	}
	if space.Description != "" {
		data["description"] = space.Description
	}
	if space.IsSystem {
		data["is_system"] = true
	}
	return data
}

func spaceFromRecord(record map[string]interface{}) models.Space {
	return models.Space{
		ID:          stringField(record, "id"),
		Name:        stringField(record, "name"),
		Slug:        stringField(record, "slug"),
		Description: stringField(record, "description"),
		IsSystem:    boolField(record, "is_system"),
		CreatedAt:   stringField(record, "created"),
		UpdatedAt:   stringField(record, "updated"),
	}
}

func spaceMemberRecordData(member models.SpaceMember) map[string]interface{} {
	return map[string]interface{}{
		"space_id": member.SpaceID,
		"user_id":  member.UserID,
		"role":     member.Role,
	}
}

func spaceMemberFromRecord(record map[string]interface{}) models.SpaceMember {
	return models.SpaceMember{
		ID:        stringField(record, "id"),
		SpaceID:   stringField(record, "space_id"),
		UserID:    stringField(record, "user_id"),
		Role:      stringField(record, "role"),
		CreatedAt: stringField(record, "created"),
		UpdatedAt: stringField(record, "updated"),
	}
}

func teamRecordData(team models.Team) map[string]interface{} {
	data := map[string]interface{}{
		"space_id": team.SpaceID,
		"name":     team.Name,
	}
	if team.Description != "" {
		data["description"] = team.Description
	}
	return data
}

func teamFromRecord(record map[string]interface{}) models.Team {
	return models.Team{
		ID:          stringField(record, "id"),
		SpaceID:     stringField(record, "space_id"),
		Name:        stringField(record, "name"),
		Description: stringField(record, "description"),
		CreatedAt:   stringField(record, "created"),
		UpdatedAt:   stringField(record, "updated"),
	}
}

func teamMemberRecordData(member models.TeamMember) map[string]interface{} {
	return map[string]interface{}{
		"team_id":     member.TeamID,
		"member_type": member.MemberType,
		"member_id":   member.MemberID,
	}
}

func teamMemberFromRecord(record map[string]interface{}) models.TeamMember {
	return models.TeamMember{
		ID:         stringField(record, "id"),
		TeamID:     stringField(record, "team_id"),
		MemberType: stringField(record, "member_type"),
		MemberID:   stringField(record, "member_id"),
		CreatedAt:  stringField(record, "created"),
		UpdatedAt:  stringField(record, "updated"),
	}
}

func boolField(record map[string]interface{}, key string) bool {
	if v, ok := record[key]; ok {
		switch t := v.(type) {
		case bool:
			return t
		case string:
			return t == "true" || t == "1"
		}
	}
	return false
}
