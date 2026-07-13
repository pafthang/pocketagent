package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/pafthang/pocketagent/internal/models"
	"github.com/pafthang/pocketagent/services/ollama-client/ollama"
)

func main() {
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, _ := nc.JetStream()

	ollama := ollama.NewClient("http://ollama:11434")

	// Subscribe to tasks
	_, err = js.Subscribe("agents.tasks.*", func(msg *nats.Msg) {
		var task models.Task
		// unmarshal...
		fmt.Printf("Executing task: %s\n", task.Prompt)

		// Simple ReAct simulation
		response, _ := ollama.Generate(ollama.GenerateRequest{
			Model:  "llama3",
			Prompt: task.Prompt,
			Stream: false,
		})

		fmt.Println("Ollama response:", response)
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Execution service with ReAct started...")
	select {}
}
