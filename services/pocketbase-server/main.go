package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	// Optional: add custom routes or migrations
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		log.Println("PocketBase embedded server starting...")
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
