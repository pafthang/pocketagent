import type { TaskEvent, TaskEventType } from "@pocketagent/client";

export interface ActivityEntry {
  id: string;
  type: "tool_start" | "tool_result" | "thinking" | "error" | "status" | "task_event";
  content: string;
  data?: Record<string, unknown>;
  timestamp: string;
}

let entryCounter = 0;
function nextId(): string {
  return `activity-${++entryCounter}`;
}

const STATUS_LABELS: Partial<Record<TaskEventType, string>> = {
  queued: "Task queued",
  orchestrating: "Orchestrating",
  subtask_dispatched: "Subtask dispatched",
  subtask_started: "Subtask started",
  subtask_completed: "Subtask completed",
  subtask_result: "Subtask result",
  supervisor_delegated: "Delegated to supervisor",
  completed: "Task completed",
  failed: "Task failed",
  cancelled: "Task cancelled",
  timeout: "Task timed out",
  llm_token: "Generating",
};

function formatTaskEvent(event: TaskEvent): string {
  const label = STATUS_LABELS[event.type] ?? event.type;
  if (event.message) {
    const msg = event.message.length > 120 ? event.message.slice(0, 120) + "…" : event.message;
    return `${label}: ${msg}`;
  }
  return label;
}

class ActivityStore {
  entries = $state<ActivityEntry[]>([]);
  isAgentWorking = $state(false);
  currentModel = $state<string | null>(null);
  tokenUsage = $state<{ input?: number; output?: number } | null>(null);
  sseActive = $state(false);

  recentEntries = $derived(this.entries.slice(-50));
  latestEntry = $derived(this.entries.at(-1) ?? null);

  clear(): void {
    this.entries = [];
    this.tokenUsage = null;
  }

  /** Push a task stream event into the activity log. */
  pushTaskEvent(event: TaskEvent): void {
    if (event.type === "llm_token") return;

    const now = new Date().toISOString();
    const isError = event.type === "failed" || event.type === "timeout";

    this.entries.push({
      id: nextId(),
      type: isError ? "error" : "task_event",
      content: formatTaskEvent(event),
      data: { ...event } as Record<string, unknown>,
      timestamp: now,
    });
  }

  /** @deprecated Legacy PocketPaw SSE shape — kept for unmigrated code paths. */
  pushSSEEvent(
    eventType: "tool_start" | "tool_result" | "thinking",
    data: Record<string, unknown>,
  ): void {
    const now = new Date().toISOString();
    this.entries.push({
      id: nextId(),
      type: eventType,
      content: String(data.tool ?? data.content ?? eventType),
      data,
      timestamp: now,
    });
  }
}

export const activityStore = new ActivityStore();