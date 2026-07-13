package main

import (
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

	log.Println("Task Orchestrator (Project Manager) started")

	// TODO: Subscribe to high-level tasks and delegate
	select {}
}
