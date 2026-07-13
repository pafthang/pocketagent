package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/pafthang/pocketagent/internal/models"
)

func main() {
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, _ := nc.JetStream()

	// Subscribe to tasks
	_, err = js.Subscribe("agents.tasks.*", func(msg *nats.Msg) {
		var task models.Task
		// TODO: unmarshal and execute with Ollama
		fmt.Printf("Received task: %s\n", task.Prompt)
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Execution service started...")
	select {} // keep running
}
