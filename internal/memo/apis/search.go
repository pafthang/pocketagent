package memoapis

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
)

func searchMemoryHandler(deps Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := apimw.SpaceIDFromContext(c)
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": apimw.HeaderSpaceID + " header is required"})
		}

		var req SearchMemoryRequest
		if err := c.Bind(&req); err != nil {
			return err
		}
		req.Query = strings.TrimSpace(req.Query)
		if req.Query == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "query is required"})
		}

		embedding, err := deps.Ollama.Embed(c.Request().Context(), req.Query)
		if err != nil {
			return c.JSON(http.StatusBadGateway, map[string]string{"error": err.Error()})
		}

		limit := req.Limit
		if limit <= 0 {
			limit = 5
		}
		if req.MinSimilarity > 0 {
			deps.Memo.MinSimilarity = req.MinSimilarity
		}

		docs, err := deps.Memo.SearchDocumentsScoped(c.Request().Context(), embedding, spaceID, limit)
		if err != nil {
			return httpx.MapMemoError(c, err)
		}

		results := make([]map[string]interface{}, 0, len(docs))
		for _, doc := range docs {
			results = append(results, map[string]interface{}{
				"id":         doc.ID,
				"content":    doc.Content,
				"similarity": doc.Similarity,
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"results": results})
	}
}