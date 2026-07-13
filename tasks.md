# pocketagent — План разработки и статус

## Общий статус проекта (обновлено 14 июля 2026)

**Проект:** `pocketagent` — Платформа управления AI-агентами

### ✅ Что полностью закрыто

- Полноценный **ReAct с реальными tools** (web_search + scrape_page)
- Correlation ID + Structured Logging
- Retry + Circuit Breaker
- Prometheus Metrics
- Embedded NATS + PocketBase
- `internal/common` и `internal/service` пакеты
- `agent-service` (скелет)

### Текущий статус по приоритетам

| Приоритет | Задача                              | Статус          | Комментарий                     |
|-----------|-------------------------------------|-----------------|---------------------------------|
| Высокий   | Полноценный ReAct с реальными tools | **Закрыто**     | Реальное выполнение web_search и scrape_page |
| Высокий   | Полноценный CRUD в agent-service    | В процессе      | Create готов, остальные — заглушки |
| Высокий   | Project Manager в task-orchestrator | Базовый         | Нужно развивать                 |
| Средний   | Memory / RAG                        | Нет             | Следующий большой этап          |
| Средний   | Тесты                               | Нет             | —                               |

---

## Статус по сервисам

### `api-gateway` — **Готов**
- Полностью использует `internal/service`
- Correlation ID, slog, WebSocket

### `agent-service` — **Скелет**
- Создан отдельный сервис
- Реализован только Create
- Остальные методы — заглушки

### `execution-service` — **Хорошо**
- Полноценный ReAct Executor с реальными инструментами
- BaseConsumer + Correlation ID

### `task-orchestrator` — **Базовый**
- Есть Consumer
- Логика делегирования минимальна

### `ollama-client`, `nats-server`, `pocketbase-server` — **Готовы**

---

## План дальнейшей работы

1. Доработать `agent-service` (все CRUD методы)
2. Развить `task-orchestrator` (настоящий Project Manager)
3. Добавить Memory/RAG сервис
4. Написать тесты
5. UI (позже)
