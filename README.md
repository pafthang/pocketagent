# pocketagent

**AI Agent Forge** — Платформа управления AI-агентами на Ollama + **PocketBase** + NATS JetStream

## ✅ Что реализовано (финальная версия)

### `internal/common`
- Structured Logging (slog) + уровни через `LOG_LEVEL`
- Correlation ID / Tracing (через NATS headers)
- Retry + Backoff
- Circuit Breaker
- Prometheus Metrics
- .env + environment support
- Dependency Health checks
- Graceful NATS shutdown

### `internal/service`
- `BaseServer` (Echo)
- `BaseConsumer` (NATS JetStream)
- `ReActExecutor`

### NATS
- Улучшенный клиент с Correlation ID
- Publish / Subscribe helpers

## Готовность

Проект имеет:
- Полную traceability задач
- Защиту от ошибок (retry, circuit breaker)
- Мониторинг (metrics)
- Graceful shutdown
- Переиспользуемые компоненты

**Следующий шаг:** `make up` или применение ко всем сервисам.
