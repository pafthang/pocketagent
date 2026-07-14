package memoapis

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
)

func memoryStatsHandler(deps Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := apimw.SpaceIDFromContext(c)
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": apimw.HeaderSpaceID + " header is required"})
		}

		stats, err := deps.Memo.Stats(c.Request().Context(), spaceID)
		if err != nil {
			return httpx.MapMemoError(c, err)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"backend":          "chromem",
			"total_memories":   stats.DocumentCount,
			"document_count":   stats.DocumentCount,
			"content_bytes":    stats.ContentBytes,
			"collection":       stats.Collection,
			"memories_by_type": map[string]int{"long_term": stats.DocumentCount},
		})
	}
}

func memorySettingsHandler(embedModel, ollamaURL string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, memorySettings{
			MemoryBackend:        "chromem",
			MemoryUseInference:   true,
			Mem0EmbedderProvider: "ollama",
			Mem0EmbedderModel:    embedModel,
			Mem0VectorStore:      "chromem",
			Mem0OllamaBaseURL:    ollamaURL,
			Mem0AutoLearn:        false,
		})
	}
}

func saveMemorySettingsHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "saved"})
	}
}