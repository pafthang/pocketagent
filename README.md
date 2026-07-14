# pocketagent

**AI Agent Forge** — multi-tenant platform for running Ollama-powered agents with PocketBase persistence, NATS JetStream orchestration, vector memory, files, skills, and MCP integrations.

The public HTTP API lives behind **gate** (`:8080`). Internal services communicate over NATS and HTTP; PocketBase is the system of record for agents, tasks, spaces, and RBAC metadata.

## Architecture

```text
Browser / API client
        │
        ▼
   gate :8080          ← auth, RBAC, tenant routing, SSE/WS
        │
        ├──► space :8083     tenants, members, invites, /authorize
        ├──► agent :8081     agents CRUD, identity, runtime-config
        ├──► files :8086     uploads, browse, project files
        ├──► memo  :8082     vector memory (chromem-go, internal)
        ├──► task  (worker)  LLM decomposition, schedules, orchestration
        ├──► exec  (worker)  ReAct tool loop, MCP, streaming events
        ├──► pocket :8090    PocketBase (auth + collections)
        ├──► nats  :4222     JetStream
        └──► ollama :11434    LLM + embeddings
```

### Services

| Service | Port (compose) | Role |
|---------|----------------|------|
| `gate` | 8080 | HTTP edge: JWT refresh, `X-Space-Id`, RBAC, proxies, tenant APIs |
| `agent` | 8081 | Agent domain: CRUD, identity files, runtime config |
| `space` | 8083 | Multi-tenancy: spaces, members, teams, invites, audit |
| `memo` | 8082 (internal) | Space-scoped vector memory + RAG |
| `files` | 8086 | Blob storage, ingest into memory, project attachments |
| `task` | health 9085 | Task orchestration, cron schedules, LLM split |
| `exec` | health 9084 | Subtask execution (native Ollama `tool_calls`) |
| `pocket` | 8090 | Embedded PocketBase |
| `nats` | 4222 / 8222 | NATS + monitoring |
| `ollama` | 11434 | Models (`llama3.1`, `nomic-embed-text`, …) |
| `ctrl` | — | Local supervisor for all of the above |

Each service follows the same layout: `cmd/<name>`, `configs/<name>.yaml`, `internal/<name>/` with `service.go`, `config.go`, `deps.go`, `routes.go`, and domain subpackages.

### NATS subjects

| Subject | Producer → Consumer |
|---------|-------------------|
| `agents.orchestrator.commands` | gate → task (decompose prompt) |
| `agents.tasks.{corr}-{n}` | task → exec (run subtask) |
| `agents.results.{corr}-{n}` | exec → task (subtask result) |
| `agents.events.{task_id}` | * → gate SSE/WebSocket |

Stream `AGENTS` is defined in `internal/nats/client/streams.go`.

### Shared packages (`pkgs/`)

| Package | Purpose |
|---------|---------|
| `pkgs/common` | Config, logging, metrics, health, retry, circuit breaker, rate limits, egress guard |
| `pkgs/service` | Echo HTTP server, NATS worker/consumer scaffolding |
| `pkgs/middle` | Auth context, `RequireSpace`, in-process `PocketRBAC` |
| `pkgs/ollama` | Ollama client, chat loop with native tool calling |
| `pkgs/models` | Shared DTOs (Agent, Task, Space, Skill, MCP, …) |
| `pkgs/httpx` | Echo error mapping helpers |

## Quick start

### Prerequisites

- Docker + Docker Compose
- [Ollama](https://ollama.com) models pulled inside the compose stack (first run may take a few minutes)

### Run the stack

```bash
cp .env.example .env   # adjust credentials if needed
make up
```

Pull models after Ollama is healthy:

```bash
docker compose exec ollama ollama pull llama3.1
docker compose exec ollama ollama pull nomic-embed-text
```

Smoke-test the full flow:

```bash
./scripts/e2e-smoke.sh
```

### Local development (no Docker)

Start every service from `configs/` via the supervisor:

```bash
make run          # or: go run ./cmd/ctrl -config-dir=configs
```

Individual binaries: `go run ./cmd/gate`, `./cmd/agent`, `./cmd/exec`, etc.

Build and test:

```bash
make build
make test
```

## Authentication and tenancy

1. Register / login through gate (`POST /auth/register`, `POST /auth/login`).
2. Send `Authorization: Bearer <token>` on every request.
3. For tenant-scoped routes, set `X-Space-Id: <space_uuid>`.
4. RBAC actions (`agent:read`, `task:write`, `memory:read`, …) are enforced in-process via PocketBase (`pkgs/middle/rbac`).

Gate proxies domain services where needed; tenant APIs for tasks, memory, skills, MCP, projects, and files are registered directly on gate with shared middleware.

## Tenant API surface (gate)

All routes below require auth + `X-Space-Id` unless noted.

| Area | Endpoints |
|------|-----------|
| Agents | `GET/POST /agents`, `GET/PUT/DELETE /agents/:id`, `GET/PUT /agents/:id/identity`, `GET /agents/:id/runtime-config` |
| Tasks | `GET/POST /tasks`, `GET/DELETE /tasks/:id`, SSE + WebSocket streams |
| Schedules | CRUD `/schedules` |
| Memory | `GET/POST /memory`, `GET/DELETE /memory/:id`, `POST /memory/search`, `GET /memory/stats` |
| Skills | CRUD `/skills`, `GET /skills/catalog`, `POST /skills/:id/run` |
| MCP | CRUD `/mcp/servers`, `GET /mcp/status`, `GET /mcp/presets`, `POST /mcp/presets/install` |
| Projects | CRUD `/projects`, items, planning, WebSocket |
| Files | `/files/*`, `/projects/:id/files/*` |
| Dashboard | `/dashboard/*` |
| Activity | `/activity` |

Spaces, members, invites, and `/authorize` are proxied to the **space** service.

## Configuration

- Per-service YAML in `configs/` (loaded via `CONFIG_DIR` or `-config-dir` for ctrl).
- Environment overrides — see [`.env.example`](.env.example).
- Catalogs: `configs/skills-catalog.json`, `configs/mcp-presets.json`.

Key variables:

| Variable | Description |
|----------|-------------|
| `POCKETBASE_SUPERUSER_EMAIL/PASSWORD` | Bootstrap admin on first `pocket` start |
| `POCKETBASE_ADMIN_EMAIL/PASSWORD` | Admin API client for workers |
| `MEMO_SERVICE_TOKEN` | Service-to-service auth for memo |
| `OLLAMA_URL`, `EMBED_MODEL`, `LLM_MODEL` | Model endpoints |
| `APP_ENV=production` | Enables rate limits and secret hardening |

## Web

| Path | Description |
|------|-------------|
| [`web/client`](web/client) | TypeScript SDK `@pocketagent/client` — namespaced API (`client.agents`, `client.tasks`, …) |
| [`web/front`](web/front) | SvelteKit UI (`@pocketagent/front`) |

The SDK targets gate only. Example:

```ts
import { createClient } from '@pocketagent/client';

const client = createClient({ baseUrl: 'http://127.0.0.1:8080', token, spaceId });
const { agents } = await client.agents.list();
const identity = await client.identity.get(agentId);
```

## Project layout

```text
cmd/           service entrypoints (gate, agent, exec, task, …)
configs/       typed YAML per service + JSON catalogs
internal/      domain services (gate, space, agent, exec, task, memo, files, pocket, nats, ctrl)
pkgs/          shared libraries (common, service, middle, ollama, models, httpx)
web/client     TypeScript API client
web/front      SvelteKit frontend
scripts/       e2e-smoke.sh
```

## PocketBase collections

Bootstrapped on first start: `users`, `spaces`, `space_members`, `space_invites`, `teams`, `agents`, `tasks`, `skills`, `mcp_servers`, `projects`, `audit_logs`, and related records. System space `admin` is created for the bootstrap superuser.

## Tool execution

`exec` runs a ReAct loop over Ollama `/api/chat` with **native `tool_calls`** (no text-based Action:/Final Answer fallback). Built-in tools include web search, scraper, optional code exec, and space-scoped MCP servers loaded from PocketBase.

## License

See repository license file if present; otherwise treat as private / all rights reserved by the maintainer.