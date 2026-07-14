package svc

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/internal/memo/store"
)

func deleteDocument(mgr *store.Manager) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID := strings.TrimSpace(c.QueryParam("space_id"))
		id := strings.TrimSpace(c.Param("id"))
		if id == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
		}

		if err := mgr.DeleteDocument(c.Request().Context(), spaceID, id); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "deleted", "id": id})
	}
}