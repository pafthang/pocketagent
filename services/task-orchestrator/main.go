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

		// Более умное разбиение (простая эвристика по ключевым словам)
		subtasks := smartSplit(taskText)
		logger.Info("Task intelligently split", "count", len(subtasks))

		results := make(map[int]string)
		var mu sync.Mutex
		var wg sync.WaitGroup

		// Подписка на результаты
		sub, _ := js.Subscribe("agents.results.*", func(m *nats.Msg) {
			mu.Lock()
			results[len(results)] = string(m.Data)
			mu.Unlock()
			wg.Done()
		})
		defer sub.Unsubscribe()

		// Публикуем подзадачи и ждём результатов
		for i, subtask := range subtasks {
			wg.Add(1)
			subject := fmt.Sprintf("agents.tasks.subtask-%d", i)
			js.Publish(subject, []byte(subtask))
		}

		// Реальное ожидание (с таймаутом)
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			logger.Info("All subtasks completed")
		case <-time.After(30 * time.Second):
			logger.Warn("Timeout waiting for results")
		}

		final := buildFinalAnswer(results)
		logger.Info("Final answer ready", "answer", final)
	})

	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Advanced Task Orchestrator started")
	select {}
}

func smartSplit(task string) []string {
	// Улучшенная эвристика
	keywords := []string{"and", "then", "after that", "also"}
	for _, kw := range keywords {
		if strings.Contains(strings.ToLower(task), kw) {
			parts := strings.Split(task, kw)
			for i := range parts {
				parts[i] = strings.TrimSpace(parts[i])
			}
			return parts
		}
	}
	return []string{task}
}

func buildFinalAnswer(results map[int]string) string {
	var sb strings.Builder
	sb.WriteString("Combined result:\n")
	for i, r := range results {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, r))
	}
	return sb.String()
}
