# pocketagent — План разработки и статус

## Общий статус проекта (на 14 июля 2026)

**Проект:** `pocketagent` — Платформа управления AI-агентами на Go + Ollama + PocketBase + NATS

### Что уже реализовано (высокий уровень)

- ✅ Микросервисная архитектура
- ✅ Embedded NATS + JetStream
- ✅ Embedded PocketBase
- ✅ Общие пакеты `internal/common` и `internal/service`
- ✅ Полная traceability (Correlation ID)
- ✅ Защита от ошибок (Retry + Circuit Breaker)
- ✅ Мониторинг (Prometheus Metrics)
- ✅ Structured Logging (slog)
- ✅ ReAct Executor с реальными tool calls (web search + scraper)
- ✅ WebSocket streaming
- ✅ Отдельный `agent-service`

### Что ещё нужно сделать (приоритет)

| Приоритет | Задача                              | Сервис(ы)              | Статус    |
|-----------|-------------------------------------|------------------------|-----------|
| Высокий   | Полноценный ReAct с реальными tools | execution-service      | В процессе |
| Высокий   | Полноценный CRUD в agent-service    | agent-service          | Скелет    |
| Высокий   | Project Manager (оркестратор)       | task-orchestrator      | Базовый   |
| Средний   | Memory / RAG сервис                 | memory-service         | Нет       |
| Средний   | Тесты                               | Все                    | Нет       |
| Низкий    | UI (Svelte / React)                 | frontend               | Нет       |
| Низкий    | OpenTelemetry tracing               | Все                    | Нет       |

---

## Статус по сервисам

### 1. `api-gateway`
- ✅ Полностью на `internal/service` шаблоне
- ✅ Correlation ID + slog
- ✅ WebSocket streaming
- ✅ Интеграция с PocketBase и NATS
- Статус: **Готов к использованию**

### 2. `agent-service`
- ✅ Создан отдельный сервис
- ✅ Create агент
- ⚠️ Остальные CRUD методы — заглушки
- Статус: **Скелет готов**

### 3. `execution-service`
- ✅ ReAct Executor с реальными tools
- ✅ BaseConsumer + Correlation ID
- ✅ Structured logging
- Статус: **Хорошо работает (нужно доработать реальное выполнение tools)**

### 4. `task-orchestrator`
- ✅ BaseConsumer
- ⚠️ Логика Project Manager — минимальная
- Статус: **Нужно развивать**

### 5. `ollama-client`
- ✅ Базовый клиент + Tool support
- Статус: **Готов**

### 6. `nats-server` (embedded)
- ✅ Работает
- ✅ Structured logging
- Статус: **Готов**

### 7. `pocketbase-server` (embedded)
- ✅ Работает
- Статус: **Готов**

---

## План дальнейшей работы (рекомендуемый порядок)

1. **Стабилизация и запуск** (`make up`)
2. **Доработка ReAct** — реальное выполнение tool calls + обработка результатов
3. **Полноценный `agent-service`** — все CRUD методы + интеграция с PocketBase
4. **Развитие `task-orchestrator`** — настоящий Project Manager (разбиение задач, делегирование)
5. **Memory / RAG сервис** (опционально)
6. **Тесты** + CI
7. **UI** (позже)

---

**Текущий приоритет:** Доработка ReAct + полноценный agent-service
