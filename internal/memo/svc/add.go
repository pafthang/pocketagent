package svc

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/internal/memo/store"
	"github.com/philippgille/chromem-go"
)

func addDocument(mgr *store.Manager) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req AddDocumentRequest
		if err := c.Bind(&req); err != nil {
			return err
		}
		if req.ID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
		}

		collection, err := mgr.Collection(req.SpaceID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		metadata := req.Metadata
		if metadata == nil {
			metadata = map[string]string{}
		}
		if req.SpaceID != "" {
			metadata["space_id"] = req.SpaceID
		}

		doc := chromem.Document{
			ID:        req.ID,
			Content:   req.Content,
			Embedding: req.Embedding,
			Metadata:  metadata,
		}
		if err := collection.AddDocument(c.Request().Context(), doc); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusCreated, map[string]string{
			"status":     "added",
			"id":         req.ID,
			"collection": collection.Name,
		})
	}
}