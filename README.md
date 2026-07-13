# pocketagent

**AI Agent Forge** — Платформа управления AI-агентами на Ollama + **PocketBase** + NATS JetStream

## Структура проекта

```
pocketagent/
├── go.work
├── Makefile
├── docker-compose.yml
├── README.md
├── services/
│   ├── api-gateway/         # HTTP + WS
│   ├── agent-service/       # CRUD агентов
│   ├── execution-service/   # Запуск ReAct
│   ├── task-orchestrator/   # Project Manager
│   ├── ollama-client/       # Клиент Ollama
│   └── pocketbase/          # (опционально)
├── internal/
│   └── models/              # Общие модели (Agent, Task и т.д.)
└── pkg/                     # Общие утилиты
```

## Что уже сделано:
- go workspace
- docker-compose с PocketBase + NATS
- Базовые модули

**Следующий шаг:** Создать общие модели и NATS-клиент.

Напиши «Продолжай» — продолжим генерацию кода.