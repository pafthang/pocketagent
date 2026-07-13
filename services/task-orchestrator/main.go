package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/pafthang/pocketagent/internal/common"
	"github.com/pafthang/pocketagent/internal/service"
)

func main() {
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, _ := nc.JetStream()

	consumer := service.NewConsumer("task-orchestrator", js)
	logger := common.NewSlogLogger("orchestrator")

	// Подписка на высокоуровневые задачи
	_, err = consumer.Subscribe("agents.orchestrator.commands", func(ctx context.Context, msg *nats.Msg) {
		logger.Info("High-level task received", "corr", common.GetCorrelationID(ctx))

		// TODO: Настоящая логика Project Manager
		// 1. Разбить задачу на подзадачи
		// 2. Назначить агентам
		// 3. Следить за выполнением
		// 4. Собрать результаты

		fmt.Println("[Orchestrator] Task delegation logic should be implemented here")
	})

	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Task Orchestrator (Project Manager) started")
	select {}
}
