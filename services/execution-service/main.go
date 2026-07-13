package main

import (
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
	agent := NewReActAgent(ollama)

	_, err = js.Subscribe("agents.tasks.*", func(msg *nats.Msg) {
		var task models.Task
		fmt.Printf("[ReAct] Starting task: %s\n", task.Prompt)

		result := agent.Run(task.Prompt)

		fmt.Println("ReAct Result:\n", result)
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Execution service with full ReAct loop ready...")
	select {}
}
