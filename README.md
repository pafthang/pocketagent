# pocketagent

**AI Agent Forge** — Платформа управления AI-агентами на Ollama + **PocketBase** + NATS JetStream

## Интеграция PocketBase

✅ Добавлена базовая интеграция:
- `internal/pocketbase/client.go` — HTTP клиент
- api-gateway теперь сохраняет агентов в PocketBase
- docker-compose включает сервис pocketbase

## Структура

```
pocketagent/
├── internal/
│   ├── models/
│   ├── nats/
│   └── pocketbase/     ← новая интеграция
├── services/
│   ├── api-gateway/ (Echo + PocketBase)
│   ├── execution-service/
│   ├── task-orchestrator/
│   └── ollama-client/
```

**Следующий шаг:** Реализовать ReAct в execution-service или улучшить PocketBase клиент.