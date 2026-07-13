package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/philippgille/chromem-go"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	// Create Chromem collection (embeddings dimension 1536 - common for many models)
	collection, err := chromem.NewCollection("memory", 1536)
	if err != nil {
		log.Fatal(err)
	}

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":  "healthy",
			"service": "memory-service",
		})
	})

	// Add document
	e.POST("/documents", func(c echo.Context) error {
		var req struct {
			ID      string `json:"id"`
			Content string `json:"content"`
			Embedding []float32 `json:"embedding"`
		}
			if err := c.Bind(&req); err != nil {
				return err
			}

			err := collection.AddDocument(req.ID, req.Content, req.Embedding)
			if err != nil {
				return c.JSON(500, map[string]string{"error": err.Error()})
			}

			return c.JSON(201, map[string]string{"status": "added", "id": req.ID})
		})

	// Semantic search
	e.POST("/search", func(c echo.Context) error {
		var req struct {
			QueryEmbedding []float32 `json:"query_embedding"`
			Limit          int       `json:"limit"`
		}
			if err := c.Bind(&req); err != nil {
				return err
			}

			if req.Limit == 0 {
				req.Limit = 5
			}

			results, err := collection.Search(req.QueryEmbedding, req.Limit)
			if err != nil {
				return c.JSON(500, map[string]string{"error": err.Error()})
			}

			return c.JSON(200, results)
		})

	e.Logger.Fatal(e.Start(":8082"))
}
