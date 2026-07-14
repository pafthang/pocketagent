package bootstrap

import "github.com/pocketbase/pocketbase/core"

func addFieldIfMissing(col *core.Collection, field core.Field) {
	if col.Fields.GetByName(field.GetName()) != nil {
		return
	}
	col.Fields.Add(field)
}