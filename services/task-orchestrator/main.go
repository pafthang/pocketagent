package main

import (
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

	log.Println("Project Manager (Task Orchestrator) started")

	// Subscribe to high-level tasks
	_, err = nc.Subscribe("agents.orchestrator.commands", func(msg *nats.Msg) {
		fmt.Println("[Orchestrator] Received high-level task")
		// TODO: break into subtasks and delegate to execution-service
	})

	if err != nil {
		log.Fatal(err)
	}

	select {}
}
