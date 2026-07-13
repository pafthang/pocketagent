package main

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/philippgille/chromem-go"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	// Create in-memory Chromem collection
	collection, err := chromem.NewCollection("memory", 1536) // 1536 for common embedding models
	if err != nil {
		log.Fatal(err)
	}

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "healthy", "service": "memory-service"})
	})

	// TODO: Add endpoints for:
	// - Add document + embedding
	// - Similarity search
	// - List documents

	e.Logger.Fatal(e.Start(":8082"))
}
