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

	_, err = js.Subscribe("agents.tasks.*", func(msg *nats.Msg) {
		var task models.Task
		fmt.Printf("[Tool Calling] Executing: %s\n", task.Prompt)

		tools := ollama.GetExampleTools()

		_, err := ollama.Generate(ollama.GenerateRequest{
			Model:  "llama3.1",
			Prompt: task.Prompt,
			Tools:  tools,
		})
		if err != nil {
			log.Println("Error:", err)
		}

		fmt.Println("Tool calling cycle completed")
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Execution service with Tool Calling ready...")
	select {}
}
