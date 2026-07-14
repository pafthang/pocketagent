package memoapis

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
)

func ingestMemoryHandler(deps Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := apimw.SpaceIDFromContext(c)
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": apimw.HeaderSpaceID + " header is required"})
		}

		var req IngestMemoryRequest
		if err := c.Bind(&req); err != nil {
			return err
		}
		req.Content = strings.TrimSpace(req.Content)
		if req.Content == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "content is required"})
		}

		id := strings.TrimSpace(req.ID)
		if id == "" {
			id = fmt.Sprintf("mem-%d", time.Now().UnixNano())
		}

		meta := copyMetadata(req.Metadata)
		meta["parent_id"] = id
		meta["created_at"] = time.Now().UTC().Format(time.RFC3339)
		if len(req.Tags) > 0 {
			meta["tags"] = strings.Join(req.Tags, ",")
		}

		if err := deps.Memo.StoreScopedWithMeta(c.Request().Context(), deps.Ollama, spaceID, id, req.Content, meta); err != nil {
			return httpx.MapMemoError(c, err)
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"id":      id,
			"status":  "ingested",
			"content": req.Content,
		})
	}
}