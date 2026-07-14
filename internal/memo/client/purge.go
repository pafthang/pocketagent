package client

import (
	"context"
	"strings"
)

// PurgeByParentID deletes a document and any chunked children sharing parent_id.
func (c *Client) PurgeByParentID(ctx context.Context, spaceID, parentID string) (int, error) {
	result, err := c.ListDocuments(ctx, spaceID, 1, 10000)
	if err != nil {
		return 0, err
	}
	ids := relatedDocumentIDs(parentID, result.Documents)
	for _, docID := range ids {
		if err := c.DeleteDocument(ctx, spaceID, docID); err != nil {
			return 0, err
		}
	}
	return len(ids), nil
}

func relatedDocumentIDs(parentID string, docs []DocumentRecord) []string {
	ids := make([]string, 0, 1)
	for _, doc := range docs {
		if doc.ID == parentID || strings.HasPrefix(doc.ID, parentID+"#") || doc.Metadata["parent_id"] == parentID {
			ids = append(ids, doc.ID)
		}
	}
	if len(ids) == 0 {
		ids = append(ids, parentID)
	}
	return ids
}