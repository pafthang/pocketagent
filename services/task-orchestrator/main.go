package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

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

		taskText := string(msg.Data)
		subtasks := splitIntoSubtasks(taskText)

		logger.Info("Task split", "subtasks", len(subtasks))

		// Карта для сбора результатов
		results := make(map[int]string)
		var mu sync.Mutex

		// Подписываемся на результаты подзадач
		_, err := js.Subscribe("agents.results.*", func(m *nats.Msg) {
			mu.Lock()
			defer mu.Unlock()

			// Простая логика: сохраняем результат
			results[len(results)] = string(m.Data)
			logger.Info("Subtask result received", "count", len(results))
		})
		if err != nil {
			logger.Error("Failed to subscribe to results", "error", err)
			return
		}

		// Публикуем подзадачи
		for i, subtask := range subtasks {
			subject := fmt.Sprintf("agents.tasks.subtask-%d", i)
			_, err := js.Publish(subject, []byte(subtask))
			if err != nil {
				logger.Error("Failed to publish subtask", "error", err)
			}
		}

		// Ждём результаты (простая задержка)
		time.Sleep(5 * time.Second)

		// Формируем финальный ответ
		finalAnswer := buildFinalAnswer(results)
		logger.Info("Final answer assembled", "answer", finalAnswer)

		// TODO: отправить финальный результат обратно клиенту
	})

	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Task Orchestrator with result collection started")
	select {}
}

func splitIntoSubtasks(task string) []string {
	if strings.Contains(task, "and") {
		parts := strings.Split(task, "and")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}
	return []string{task}
}

func buildFinalAnswer(results map[int]string) string {
	var sb strings.Builder
	sb.WriteString("Final combined answer:\n")
	for i, res := range results {
		sb.WriteString(fmt.Sprintf("- Subtask %d: %s\n", i, res))
	}
	return sb.String()
}
