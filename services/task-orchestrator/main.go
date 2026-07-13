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

	_, err = consumer.Subscribe("agents.orchestrator.commands", func(ctx context.Context, msg *nats.Msg) {
		corrID := common.GetCorrelationID(ctx)
		logger.Info("High-level task received", "correlation_id", corrID)

		// ============================================
		// TODO: Реальная логика Project Manager
		// ============================================
		// 1. Разбить высокоуровневую задачу на подзадачи
		// 2. Определить, каким агентам делегировать
		// 3. Опубликовать подзадачи в agents.tasks.*
		// 4. Подписаться на результаты
		// 5. Собрать результаты и сформировать финальный ответ
		// ============================================

		fmt.Printf("[Orchestrator] Would now break task into subtasks and delegate (corr=%s)\n", corrID)
	})

	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Task Orchestrator (Project Manager) is running")
	select {}
}
