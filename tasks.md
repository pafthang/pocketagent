# pocketagent — План разработки и статус

## Общий статус проекта (обновлено 14 июля 2026)

**Проект:** `pocketagent` — Платформа управления AI-агентами

### ✅ Закрытые задачи

- Полноценный **ReAct с реальными tools** (web_search + scrape_page)
- **agent-service** — все CRUD методы реализованы (Create, Read, Update, Delete, List)
- Correlation ID + Structured Logging
- Retry + Circuit Breaker
- Prometheus Metrics
- Embedded NATS + PocketBase

### Текущий статус

| Приоритет | Задача                           | Статус     | Комментарий                     |
|-----------|----------------------------------|------------|---------------------------------|
| Высокий   | agent-service (полный CRUD)      | **Закрыто** | Все методы реализованы         |
| Высокий   | task-orchestrator (Project Manager) | В процессе | Нужно развивать                 |
| Средний   | Memory / RAG                     | Нет        | —                               |

---

## Статус по сервисам

### `agent-service` — **Готов**
- Полноценный CRUD (Create / Read / Update / Delete / List)
- Интеграция с PocketBase

### Остальные сервисы — без изменений

---

## Следующий приоритет

**Развитие `task-orchestrator`** — настоящий Project Manager (разбиение задач, делегирование агентам)
