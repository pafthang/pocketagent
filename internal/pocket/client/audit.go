package client

import (
	"fmt"

	"github.com/pafthang/pocketagent/pkgs/models"
)



// CreateAuditLog stores an audit event.
func (c *Client) CreateAuditLog(entry models.AuditLog) (models.AuditLog, error) {
	data := map[string]interface{}{
		"space_id": entry.SpaceID,
		"action":   entry.Action,
	}
	if entry.ActorID != "" {
		data["actor_id"] = entry.ActorID
	}
	if entry.ActorEmail != "" {
		data["actor_email"] = entry.ActorEmail
	}
	if entry.ResourceType != "" {
		data["resource_type"] = entry.ResourceType
	}
	if entry.ResourceID != "" {
		data["resource_id"] = entry.ResourceID
	}
	if entry.Metadata != nil {
		data["metadata"] = entry.Metadata
	}
	if entry.IPAddress != "" {
		data["ip_address"] = entry.IPAddress
	}

	record, err := c.CreateRecord(AuditLogsCollection, data)
	if err != nil {
		return models.AuditLog{}, err
	}
	return auditFromRecord(record), nil
}

// ListAuditLogs returns audit events for a space.
func (c *Client) ListAuditLogs(spaceID string, opts ListOptions) ([]models.AuditLog, int, error) {
	filter := fmt.Sprintf("space_id = %q", spaceID)
	if opts.Filter != "" {
		filter = filter + " && (" + opts.Filter + ")"
	}
	opts.Filter = filter
	records, total, err := c.ListRecordsOpts(AuditLogsCollection, opts)
	if err != nil {
		return nil, 0, err
	}
	logs := make([]models.AuditLog, 0, len(records))
	for _, record := range records {
		logs = append(logs, auditFromRecord(record))
	}
	return logs, total, nil
}

func auditFromRecord(record map[string]interface{}) models.AuditLog {
	var metadata map[string]interface{}
	if raw, ok := record["metadata"]; ok {
		if m, ok := raw.(map[string]interface{}); ok {
			metadata = m
		}
	}
	return models.AuditLog{
		ID:           stringField(record, "id"),
		SpaceID:      stringField(record, "space_id"),
		ActorID:      stringField(record, "actor_id"),
		ActorEmail:   stringField(record, "actor_email"),
		Action:       stringField(record, "action"),
		ResourceType: stringField(record, "resource_type"),
		ResourceID:   stringField(record, "resource_id"),
		Metadata:     metadata,
		IPAddress:    stringField(record, "ip_address"),
		CreatedAt:    stringField(record, "created"),
	}
}
