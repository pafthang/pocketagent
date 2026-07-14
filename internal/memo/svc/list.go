package svc

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/internal/memo/store"
)

func listDocuments(mgr *store.Manager) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID := strings.TrimSpace(c.QueryParam("space_id"))
		page, _ := strconv.Atoi(c.QueryParam("page"))
		perPage, _ := strconv.Atoi(c.QueryParam("per_page"))

		items, total, err := mgr.ListDocuments(c.Request().Context(), spaceID, page, perPage)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"documents": items,
			"total":     total,
			"page":      maxInt(page, 1),
			"per_page":  defaultPerPage(perPage),
		})
	}
}

func getDocument(mgr *store.Manager) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID := strings.TrimSpace(c.QueryParam("space_id"))
		id := strings.TrimSpace(c.Param("id"))
		if id == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
		}

		doc, err := mgr.GetDocument(c.Request().Context(), spaceID, id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "document not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, doc)
	}
}