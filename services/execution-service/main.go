package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/pafthang/pocketagent/internal/common"
	"github.com/pafthang/pocketagent/internal/models"
	"github.com/pafthang/pocketagent/internal/service"
	"github.com/pafthang/pocketagent/services/ollama-client/ollama"
)

func main() {
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, _ := nc.JetStream()

	consumer := service.NewConsumer("execution-service", js)
	ollamaClient := ollama.NewClient("http://ollama:11434")

	_, err = consumer.Subscribe("agents.tasks.*", func(ctx context.Context, msg *nats.Msg) {
		var task models.Task
		fmt.Printf("[ReAct] Task received (corr=%s): %s\n", 
			common.GetCorrelationID(ctx), task.Prompt)

		// TODO: full ReAct
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Execution service with BaseConsumer ready...")
	select {}
}
