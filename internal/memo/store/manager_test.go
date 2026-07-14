package store

import "testing"

func TestSpaceCollectionName(t *testing.T) {
	name := spaceCollectionName("abc-123")
	if name != "space_abc-123" {
		t.Fatalf("unexpected name: %s", name)
	}
	name = spaceCollectionName("bad id!")
	if name != "space_bad_id_" {
		t.Fatalf("unexpected sanitized name: %s", name)
	}
}