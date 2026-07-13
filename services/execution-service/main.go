package main

import (
	"context"
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
	executor := service.NewReActExecutor(ollamaClient, ollama.GetExampleTools())

	_, err = consumer.Subscribe("agents.tasks.*", func(ctx context.Context, msg *nats.Msg) {
		var task models.Task
		logger := common.LogWithCorrelation(common.NewSlogLogger("execution"), ctx)

		logger.Info("Starting ReAct task", "prompt", task.Prompt)

		result, err := executor.Execute(ctx, task.Prompt)
		if err != nil {
			logger.Error("ReAct execution failed", "error", err)
			return
		}

		logger.Info("ReAct completed", 
			"final_answer", result.FinalAnswer,
			"steps", len(result.Steps),
			"tool_calls", len(result.ToolCalls),
		)
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Execution service with improved ReAct ready...")
	select {}
}
