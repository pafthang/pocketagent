package client

import "testing"

func TestRelatedDocumentIDs(t *testing.T) {
	ids := relatedDocumentIDs("mem-1", []DocumentRecord{
		{ID: "mem-1#0", Metadata: map[string]string{"parent_id": "mem-1"}},
		{ID: "mem-1#1", Metadata: map[string]string{"parent_id": "mem-1"}},
		{ID: "mem-2"},
	})
	if len(ids) != 2 {
		t.Fatalf("expected 2 ids, got %#v", ids)
	}
}