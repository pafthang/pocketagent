# pocketagent

**AI Agent Forge** — Платформа управления AI-агентами на Ollama + PocketBase + NATS

## Обновлённое ТЗ

**Ключевые изменения:**
- Вместо SQLite/Postgres — **PocketBase** как основной бэкенд (авторизация, хранилище, realtime, SDK для Go).
- Новый сервис: `pocketbase` (или встроенный запуск).
- Все сервисы работают с PocketBase через официальный Go SDK или HTTP.

Остальная архитектура остаётся микросервисной на Go + NATS JetStream.

Готовы приступать к структуре проекта.