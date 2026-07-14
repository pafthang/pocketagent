package fileapis

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	filepath "github.com/pafthang/pocketagent/internal/files/path"
	"github.com/pafthang/pocketagent/pkgs/httpx"
)

func ingestFileHandler(deps *Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		file, err := loadFileInSpace(deps.PB, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if file.IsDir || file.StorageKey == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "not a file"})
		}

		var req ingestFileRequest
		_ = c.Bind(&req)

		if file.MemoIngested && !req.Force {
			return c.JSON(http.StatusConflict, map[string]string{"error": "file already ingested (use force=true to re-ingest)"})
		}

		content, err := filepath.ReadTextBlob(func() (io.ReadCloser, error) {
			return deps.Store.Open(file.StorageKey)
		}, file.MimeType, filepath.MaxIngestBytes)
		if err != nil {
			if filepath.IsTextMime(file.MimeType) {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			}
			return c.JSON(http.StatusUnsupportedMediaType, map[string]string{"error": "only text files can be ingested into memo"})
		}
		content = strings.TrimSpace(content)
		if content == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "file has no text content"})
		}

		memoID := memoDocumentIDForFile(file.ID)
		replaced := file.MemoIngested || req.Force
		if replaced {
			if _, err := deps.Memo.PurgeByParentID(c.Request().Context(), spaceID, memoID); err != nil {
				return httpx.MapMemoError(c, err)
			}
		}

		meta := map[string]string{
			"parent_id":    memoID,
			"source":       "file",
			"file_id":      file.ID,
			"virtual_path": file.VirtualPath,
			"name":         file.Name,
			"mime_type":    file.MimeType,
			"created_at":   time.Now().UTC().Format(time.RFC3339),
		}
		if file.ProjectID != "" {
			meta["project_id"] = file.ProjectID
		}
		if len(req.Tags) > 0 {
			meta["tags"] = strings.Join(req.Tags, ",")
		} else if len(file.Tags) > 0 {
			meta["tags"] = strings.Join(file.Tags, ",")
		}

		if err := deps.Memo.StoreScopedWithMeta(c.Request().Context(), deps.Ollama, spaceID, memoID, content, meta); err != nil {
			return httpx.MapMemoError(c, err)
		}

		file.MemoIngested = true
		if _, err := deps.PB.UpdateFile(file.ID, file); err != nil {
			return httpx.MapPocketError(c, err)
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"status":        "ingested",
			"file_id":       file.ID,
			"memo_id":       memoID,
			"path":          file.VirtualPath,
			"content_bytes": len(content),
			"replaced":      replaced,
		})
	}
}