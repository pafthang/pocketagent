package bootstrap

import (
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pocketbase/pocketbase/core"
)

func ensureFilesCollection(app core.App) error {
	return ensureLockedCollection(app, pbclient.FilesCollection, func(col *core.Collection) {
		addFieldIfMissing(col, &core.TextField{Name: "space_id", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "project_id"})
		addFieldIfMissing(col, &core.TextField{Name: "parent_id"})
		addFieldIfMissing(col, &core.TextField{Name: "name", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "virtual_path", Required: true})
		addFieldIfMissing(col, &core.BoolField{Name: "is_dir"})
		addFieldIfMissing(col, &core.TextField{Name: "mime_type"})
		addFieldIfMissing(col, &core.NumberField{Name: "size"})
		addFieldIfMissing(col, &core.TextField{Name: "storage_key"})
		addFieldIfMissing(col, &core.TextField{Name: "checksum"})
		addFieldIfMissing(col, &core.BoolField{Name: "memo_ingested"})
		addFieldIfMissing(col, &core.TextField{Name: "uploaded_by"})
		addFieldIfMissing(col, &core.JSONField{Name: "tags"})
	})
}