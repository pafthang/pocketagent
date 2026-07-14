import type { Task } from "@pocketagent/client";
import { taskStreamId } from "@pocketagent/client";
import type { ChatMessage, Session } from "$lib/api";

const TITLE_MAX = 48;

export function truncateTitle(text: string, max = TITLE_MAX): string {
  const oneLine = text.replace(/\s+/g, " ").trim();
  if (oneLine.length <= max) return oneLine || "New chat";
  return oneLine.slice(0, max - 1) + "…";
}

/** Build chat messages from tasks sharing a correlation id (sorted oldest first). */
export function tasksToMessages(tasks: Task[]): ChatMessage[] {
  const sorted = [...tasks]
    .filter((t) => !t.parent_id)
    .sort((a, b) => String(a.created_at ?? "").localeCompare(String(b.created_at ?? "")));

  const messages: ChatMessage[] = [];
  for (const task of sorted) {
    if (task.prompt) {
      messages.push({
        id: `${task.id}-user`,
        role: "user",
        content: task.prompt,
        timestamp: task.created_at,
      });
    }
    const assistant = task.result || (task.error ? `Error: ${task.error}` : "");
    if (assistant) {
      messages.push({
        id: `${task.id}-assistant`,
        role: "assistant",
        content: assistant,
        timestamp: task.updated_at ?? task.created_at,
      });
    }
  }
  return messages;
}

export function taskToMessages(task: Task): ChatMessage[] {
  return tasksToMessages([task]);
}

export function taskToSession(
  task: Task,
  opts?: { title?: string; messageCount?: number },
): Session {
  const id = taskStreamId(task);
  return {
    id,
    title: opts?.title ?? truncateTitle(task.prompt),
    channel: "task",
    last_activity: task.updated_at ?? task.created_at ?? new Date().toISOString(),
    message_count: opts?.messageCount ?? (task.result || task.error ? 2 : task.prompt ? 1 : 0),
  };
}

/** Group root tasks by correlation id into sidebar sessions. */
export function groupTasksIntoSessions(
  tasks: Task[],
  opts?: {
    titleOverrides?: Record<string, string>;
    hiddenIds?: Set<string>;
  },
): Session[] {
  const byCorr = new Map<string, Task[]>();
  for (const task of tasks) {
    if (task.parent_id) continue;
    const key = taskStreamId(task);
    if (opts?.hiddenIds?.has(key)) continue;
    const list = byCorr.get(key) ?? [];
    list.push(task);
    byCorr.set(key, list);
  }

  const sessions: Session[] = [];
  for (const [id, group] of byCorr) {
    const sorted = [...group].sort((a, b) =>
      String(b.updated_at ?? b.created_at ?? "").localeCompare(
        String(a.updated_at ?? a.created_at ?? ""),
      ),
    );
    const latest = sorted[0];
    const first = [...group].sort((a, b) =>
      String(a.created_at ?? "").localeCompare(String(b.created_at ?? "")),
    )[0];
    const messageCount = group.reduce((n, t) => {
      let c = t.prompt ? 1 : 0;
      if (t.result || t.error) c += 1;
      return n + c;
    }, 0);
    sessions.push(
      taskToSession(latest, {
        title: opts?.titleOverrides?.[id] ?? truncateTitle(first.prompt),
        messageCount,
      }),
    );
  }

  return sessions.sort((a, b) => b.last_activity.localeCompare(a.last_activity));
}

export function appendFileContextToPrompt(
  prompt: string,
  ctx?: {
    current_dir?: string;
    open_file?: string;
    open_file_name?: string;
    selected_files?: string[];
    source?: string;
  },
): string {
  if (!ctx) return prompt;
  const parts: string[] = [];
  if (ctx.current_dir) parts.push(`Current directory: ${ctx.current_dir}`);
  if (ctx.open_file) {
    parts.push(
      `Open file: ${ctx.open_file_name ?? ctx.open_file}${ctx.source ? ` (${ctx.source})` : ""}`,
    );
  }
  if (ctx.selected_files?.length) {
    parts.push(`Selected files: ${ctx.selected_files.join(", ")}`);
  }
  if (parts.length === 0) return prompt;
  return `${prompt}\n\n---\nFile context:\n${parts.join("\n")}`;
}