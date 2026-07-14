package memoapis

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
)

func deleteMemoryHandler(deps Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := apimw.SpaceIDFromContext(c)
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": apimw.HeaderSpaceID + " header is required"})
		}

		id := strings.TrimSpace(c.Param("id"))
		if id == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
		}

		removed, err := deps.Memo.PurgeByParentID(c.Request().Context(), spaceID, id)
		if err != nil {
			return httpx.MapMemoError(c, err)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "deleted",
			"id":      id,
			"removed": removed,
		})
	}
}