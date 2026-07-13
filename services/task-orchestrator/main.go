package main

import (
	"context"
	"fmt"
	"log"
	"strings"

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

		// Простая эвристика разбиения задачи
		taskText := string(msg.Data)
		subtasks := splitIntoSubtasks(taskText)

		logger.Info("Task split into subtasks", "count", len(subtasks))

		for i, subtask := range subtasks {
			// Публикуем подзадачу
			subject := fmt.Sprintf("agents.tasks.subtask-%d", i)
			_, err := js.Publish(subject, []byte(subtask))
			if err != nil {
				logger.Error("Failed to publish subtask", "error", err)
				continue
			}
			logger.Info("Subtask published", "subject", subject, "task", subtask)
		}

		// TODO: В будущем — собирать результаты и формировать финальный ответ
	})

	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Task Orchestrator (Project Manager) started with basic delegation logic")
	select {}
}

// Простая функция разбиения задачи на подзадачи
func splitIntoSubtasks(task string) []string {
	// Очень простая эвристика
	if strings.Contains(task, "and") {
		parts := strings.Split(task, "and")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}

	// Если не получилось разбить — возвращаем как одну задачу
	return []string{task}
}
