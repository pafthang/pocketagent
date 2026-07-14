package bootstrap

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
)

func TestCollectionIsSuperuserLocked(t *testing.T) {
	col := core.NewBaseCollection("test")
	lockCollectionToSuperuser(col)
	if !collectionIsSuperuserLocked(col) {
		t.Fatal("expected superuser-locked collection")
	}

	public := ""
	col2 := core.NewBaseCollection("test2")
	col2.ListRule = &public
	if collectionIsSuperuserLocked(col2) {
		t.Fatal("public list rule should not be locked")
	}
}

func TestCollectionNeedsLock(t *testing.T) {
	col := core.NewBaseCollection("test")
	lockCollectionToSuperuser(col)
	if collectionNeedsLock(col) {
		t.Fatal("already locked collection should not need lock")
	}

	public := ""
	col2 := core.NewBaseCollection("test2")
	col2.CreateRule = &public
	if !collectionNeedsLock(col2) {
		t.Fatal("public create rule should need lock")
	}
}

func TestIsPublicRule(t *testing.T) {
	empty := ""
	if !isPublicRule(&empty) {
		t.Fatal("empty string rule is public")
	}
	nilRule := (*string)(nil)
	if isPublicRule(nilRule) {
		t.Fatal("nil rule is not public")
	}
	locked := `@request.auth.id != ""`
	if isPublicRule(&locked) {
		t.Fatal("non-empty rule is not public")
	}
}