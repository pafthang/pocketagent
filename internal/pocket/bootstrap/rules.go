package bootstrap

import (
	"database/sql"
	"errors"

	"github.com/pocketbase/pocketbase/core"
)

// lockCollectionToSuperuser restricts API access to authenticated superusers only.
func lockCollectionToSuperuser(col *core.Collection) {
	col.ListRule = nil
	col.ViewRule = nil
	col.CreateRule = nil
	col.UpdateRule = nil
	col.DeleteRule = nil
}

func isPublicRule(rule *string) bool {
	return rule != nil && *rule == ""
}

func collectionIsSuperuserLocked(col *core.Collection) bool {
	return col.ListRule == nil &&
		col.ViewRule == nil &&
		col.CreateRule == nil &&
		col.UpdateRule == nil &&
		col.DeleteRule == nil
}

func collectionNeedsLock(col *core.Collection) bool {
	if collectionIsSuperuserLocked(col) {
		return false
	}
	return isPublicRule(col.ListRule) ||
		isPublicRule(col.ViewRule) ||
		isPublicRule(col.CreateRule) ||
		isPublicRule(col.UpdateRule) ||
		isPublicRule(col.DeleteRule) ||
		col.ListRule != nil ||
		col.ViewRule != nil ||
		col.CreateRule != nil ||
		col.UpdateRule != nil ||
		col.DeleteRule != nil
}

func ensureLockedCollection(app core.App, name string, extend func(col *core.Collection)) error {
	col, err := app.FindCollectionByNameOrId(name)
	if errors.Is(err, sql.ErrNoRows) {
		col = core.NewBaseCollection(name)
		if extend != nil {
			extend(col)
		}
		lockCollectionToSuperuser(col)
		return app.Save(col)
	}
	if err != nil {
		return err
	}

	beforeLen := len(col.Fields)
	if extend != nil {
		extend(col)
	}

	needsSave := collectionNeedsLock(col) || len(col.Fields) != beforeLen
	if !needsSave {
		return nil
	}

	lockCollectionToSuperuser(col)
	return app.Save(col)
}