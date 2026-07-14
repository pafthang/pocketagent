package client

// JetStream subject layout and message flow.
//
// Streams (see streams.go):
//   - AGENTS      — orchestration, tasks, results, events, project planning
//   - AGENTS_DLQ  — dead letters after MaxDeliver (default 5)
//
// Subject tree:
//
//	agents.orchestrator.commands     task worker: decompose user prompt → subtasks
//	agents.projects.plan.commands    task worker: project planning pipeline
//	agents.tasks.{corr}-{index}      exec worker: ReAct subtask execution
//	agents.results.{corr}-{index}    task worker: subtask completion
//	agents.events.{root_task_id}     progress/SSE/WS (gate subscribes per task)
//	agents.dlq.{service}             failed payloads (exec | task | …)
//
// Flow (happy path):
//
//	┌──────┐  orchestrator.commands   ┌──────┐  tasks.{corr}-N   ┌──────┐
//	│ gate │ ────────────────────────► │ task │ ────────────────► │ exec │
//	└──────┘                           └──────┘                   └──────┘
//	    ▲                                  │  results.{corr}-N         │
//	    │                                  ◄──────────────────────────┘
//	    │  events.{task_id}  ◄── publish progress (task + exec)
//	    └── WebSocket / SSE / polling
//
// Project planning:
//
//	gate ──projects.plan.commands──► task (planner) ──events.{project_id}──► gate WS
//
// Retry / DLQ (SubscribeJSON in internal/service/consumer.go):
//
//	handler error ──► NAK (redeliver) until NumDelivered >= MaxDeliver (5)
//	             └──► PublishDLQ → agents.dlq.{consumer_name} + Ack
//	                  metric: pocketagent_dlq_messages_total{service,reason}
//	                  health: /health dependency "dlq" down when depth >= threshold
const (
	SubjectAgentsAll = "agents.>"

	SubjectOrchestrator  = "agents.orchestrator.commands"
	SubjectProjectsPlan  = "agents.projects.plan.commands"
	SubjectTasks         = "agents.tasks.*"
	SubjectResults       = "agents.results.*"
	SubjectEvents        = "agents.events.*"
	SubjectTasksPrefix   = "agents.tasks."
	SubjectResultsPrefix = "agents.results."
	SubjectEventsPrefix  = "agents.events."
)

// EventSubject returns the subject for task progress events.
func EventSubject(taskID string) string {
	return SubjectEventsPrefix + taskID
}