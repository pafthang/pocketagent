package svc

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/internal/memo/store"
)

func collectionStats(mgr *store.Manager) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID := strings.TrimSpace(c.QueryParam("space_id"))
		stats, err := mgr.Stats(c.Request().Context(), spaceID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, stats)
	}
}