package client

import (
	"fmt"

	"github.com/pafthang/pocketagent/pkgs/models"
)



func (c *Client) CreateFile(record models.StoredFile) (models.StoredFile, error) {
	out, err := c.CreateRecord(FilesCollection, fileRecordData(record))
	if err != nil {
		return models.StoredFile{}, err
	}
	return fileFromRecord(out), nil
}

func (c *Client) GetFile(id string) (models.StoredFile, error) {
	record, err := c.GetRecord(FilesCollection, id)
	if err != nil {
		return models.StoredFile{}, err
	}
	return fileFromRecord(record), nil
}

func (c *Client) ListFiles(opts ListOptions) ([]models.StoredFile, int, error) {
	records, total, err := c.ListRecordsOpts(FilesCollection, opts)
	if err != nil {
		return nil, 0, err
	}
	out := make([]models.StoredFile, 0, len(records))
	for _, record := range records {
		out = append(out, fileFromRecord(record))
	}
	return out, total, nil
}

func (c *Client) UpdateFile(id string, record models.StoredFile) (models.StoredFile, error) {
	out, err := c.UpdateRecord(FilesCollection, id, fileRecordData(record))
	if err != nil {
		return models.StoredFile{}, err
	}
	return fileFromRecord(out), nil
}

func (c *Client) DeleteFile(id string) error {
	return c.DeleteRecord(FilesCollection, id)
}

func (c *Client) FindFileByPath(spaceID, virtualPath string) (models.StoredFile, error) {
	filter := fmt.Sprintf("space_id = %q && virtual_path = %q", spaceID, virtualPath)
	records, _, err := c.ListRecordsOpts(FilesCollection, ListOptions{Page: 1, PerPage: 1, Filter: filter})
	if err != nil {
		return models.StoredFile{}, err
	}
	if len(records) == 0 {
		return models.StoredFile{}, &APIError{StatusCode: 404, Message: "file not found"}
	}
	return fileFromRecord(records[0]), nil
}

func (c *Client) ListChildren(spaceID, parentID, projectID string, page, perPage int) ([]models.StoredFile, int, error) {
	filter := fmt.Sprintf("space_id = %q && parent_id = %q", spaceID, parentID)
	if projectID != "" {
		filter += fmt.Sprintf(" && project_id = %q", projectID)
	} else {
		filter += ` && project_id = ""`
	}
	return c.ListFiles(ListOptions{Page: page, PerPage: perPage, Filter: filter})
}

func FilesFilter(spaceID string) string {
	return fmt.Sprintf("space_id = %q", spaceID)
}

func fileRecordData(record models.StoredFile) map[string]interface{} {
	data := map[string]interface{}{
		"space_id":     record.SpaceID,
		"name":         record.Name,
		"virtual_path": record.VirtualPath,
		"is_dir":       record.IsDir,
	}
	if record.ProjectID != "" {
		data["project_id"] = record.ProjectID
	}
	if record.ParentID != "" {
		data["parent_id"] = record.ParentID
	}
	if record.MimeType != "" {
		data["mime_type"] = record.MimeType
	}
	if record.Size > 0 {
		data["size"] = record.Size
	}
	if record.StorageKey != "" {
		data["storage_key"] = record.StorageKey
	}
	if record.Checksum != "" {
		data["checksum"] = record.Checksum
	}
	if record.MemoIngested {
		data["memo_ingested"] = true
	}
	if record.UploadedBy != "" {
		data["uploaded_by"] = record.UploadedBy
	}
	if record.Tags != nil {
		data["tags"] = record.Tags
	}
	return data
}

func fileFromRecord(record map[string]interface{}) models.StoredFile {
	file := models.StoredFile{
		ID:          stringField(record, "id"),
		SpaceID:     stringField(record, "space_id"),
		ProjectID:   stringField(record, "project_id"),
		ParentID:    stringField(record, "parent_id"),
		Name:        stringField(record, "name"),
		VirtualPath: stringField(record, "virtual_path"),
		MimeType:    stringField(record, "mime_type"),
		StorageKey:  stringField(record, "storage_key"),
		Checksum:    stringField(record, "checksum"),
		UploadedBy:  stringField(record, "uploaded_by"),
		CreatedAt:   stringField(record, "created"),
		UpdatedAt:   stringField(record, "updated"),
	}
	if v, ok := record["is_dir"].(bool); ok {
		file.IsDir = v
	}
	if v, ok := record["memo_ingested"].(bool); ok {
		file.MemoIngested = v
	}
	if v, ok := record["size"]; ok {
		switch n := v.(type) {
		case float64:
			file.Size = int64(n)
		case int64:
			file.Size = n
		}
	}
	file.Tags = stringSliceField(record, "tags")
	return file
}
