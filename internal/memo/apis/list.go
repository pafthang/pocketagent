package memoapis

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
)

func listMemoryHandler(deps Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := apimw.SpaceIDFromContext(c)
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": apimw.HeaderSpaceID + " header is required"})
		}

		page, _ := strconv.Atoi(c.QueryParam("page"))
		perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
		if limit, err := strconv.Atoi(c.QueryParam("limit")); err == nil && limit > 0 && perPage <= 0 {
			perPage = limit
		}

		result, err := deps.Memo.ListDocuments(c.Request().Context(), spaceID, page, perPage)
		if err != nil {
			return httpx.MapMemoError(c, err)
		}

		items := make([]memoryDocument, 0, len(result.Documents))
		for _, doc := range result.Documents {
			items = append(items, toMemoryDocument(doc))
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"documents": items,
			"total":     result.Total,
			"page":      result.Page,
			"per_page":  result.PerPage,
		})
	}
}

func getMemoryHandler(deps Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := apimw.SpaceIDFromContext(c)
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": apimw.HeaderSpaceID + " header is required"})
		}

		id := strings.TrimSpace(c.Param("id"))
		if id == "" || id == "stats" || id == "settings" || id == "search" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
		}

		doc, err := deps.Memo.GetDocument(c.Request().Context(), spaceID, id)
		if err != nil {
			return httpx.MapMemoError(c, err)
		}

		return c.JSON(http.StatusOK, toMemoryDocument(doc))
	}
}