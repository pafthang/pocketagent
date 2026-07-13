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
		subtasks := smartSplit(taskText)

		logger.Info("Task split into subtasks", "count", len(subtasks))

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

		// === ПАРАЛЛЕЛЬНАЯ публикация подзадач ===
		for i, subtask := range subtasks {
			wg.Add(1)
			go func(idx int, task string) {
				subject := fmt.Sprintf("agents.tasks.subtask-%d", idx)
				js.Publish(subject, []byte(task))
				logger.Info("Subtask published in parallel", "subject", subject)
			}(i, subtask)
		}

		// Реальное ожидание результатов
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			logger.Info("All parallel subtasks completed")
		case <-time.After(30 * time.Second):
			logger.Warn("Timeout waiting for parallel results")
		}

		final := buildFinalAnswer(results)
		logger.Info("Final parallel answer ready", "answer", final)
	})

	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Task Orchestrator with parallel execution started")
	select {}
}

func smartSplit(task string) []string {
	if strings.Contains(strings.ToLower(task), "and") {
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
	sb.WriteString("Parallel execution result:\n")
	for i, r := range results {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, r))
	}
	return sb.String()
}
