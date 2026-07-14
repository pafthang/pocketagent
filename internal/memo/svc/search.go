package svc

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/internal/memo/store"
)

func searchDocuments(mgr *store.Manager) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req SearchRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		if req.Limit <= 0 {
			req.Limit = 5
		}
		minSim := req.MinSimilarity
		if minSim <= 0 {
			minSim = mgr.MinSimilarity()
		}

		collection, err := mgr.Collection(req.SpaceID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		queryLimit := req.Limit * 3
		if queryLimit < req.Limit {
			queryLimit = req.Limit
		}

		results, err := collection.QueryEmbedding(context.Background(), req.QueryEmbedding, queryLimit, nil, nil)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		filtered := store.FilterBySimilarity(results, minSim, req.Limit)
		return c.JSON(http.StatusOK, toSearchResponse(filtered))
	}
}