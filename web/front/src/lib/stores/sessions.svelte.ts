import type { Session } from "$lib/api";
import type { Task } from "@pocketagent/client";
import { taskStreamId } from "@pocketagent/client";
import { toast } from "svelte-sonner";
import { connectionStore } from "./connection.svelte";
import { chatStore } from "./chat.svelte";
import { logger } from "$lib/utils/logger";
import { groupTasksIntoSessions, tasksToMessages, truncateTitle } from "$lib/types/chat";

const STORAGE_KEY = "pocketagent_active_session";
const PINNED_KEY = "pocketagent_pinned_sessions";
const TITLES_KEY = "pocketagent_session_titles";
const HIDDEN_KEY = "pocketagent_hidden_sessions";

class SessionStore {
  sessions = $state<Session[]>([]);
  activeSessionId = $state<string | null>(null);
  isLoading = $state(false);
  isLoadingHistory = $state(false);
  pinnedSessionIds = $state<Set<string>>(new Set());

  /** Raw tasks cache for history reconstruction. */
  private tasksByCorrelation = new Map<string, Task[]>();
  private titleOverrides: Record<string, string> = {};
  private hiddenIds = new Set<string>();

  activeSession = $derived(
    this.sessions.find((s) => s.id === this.activeSessionId) ?? null,
  );

  pinnedSessions = $derived(
    this.sessions.filter((s) => this.pinnedSessionIds.has(s.id)),
  );

  constructor() {
    try {
      const pinned = localStorage.getItem(PINNED_KEY);
      if (pinned) this.pinnedSessionIds = new Set(JSON.parse(pinned));

      const titles = localStorage.getItem(TITLES_KEY);
      if (titles) this.titleOverrides = JSON.parse(titles);

      const hidden = localStorage.getItem(HIDDEN_KEY);
      if (hidden) this.hiddenIds = new Set(JSON.parse(hidden));
    } catch {
      // ignore
    }
  }

  private persistTitles(): void {
    try {
      localStorage.setItem(TITLES_KEY, JSON.stringify(this.titleOverrides));
    } catch {
      // ignore
    }
  }

  private persistHidden(): void {
    try {
      localStorage.setItem(HIDDEN_KEY, JSON.stringify([...this.hiddenIds]));
    } catch {
      // ignore
    }
  }

  private rebuildSessions(): void {
    const allTasks = [...this.tasksByCorrelation.values()].flat();
    this.sessions = groupTasksIntoSessions(allTasks, {
      titleOverrides: this.titleOverrides,
      hiddenIds: this.hiddenIds,
    });
  }

  private indexTasks(tasks: Task[]): void {
    const map = new Map<string, Task[]>();
    for (const task of tasks) {
      if (task.parent_id) continue;
      const key = taskStreamId(task);
      const list = map.get(key) ?? [];
      list.push(task);
      map.set(key, list);
    }
    this.tasksByCorrelation = map;
    this.rebuildSessions();
  }

  upsertTask(task: Task): void {
    if (task.parent_id) return;
    const key = taskStreamId(task);
    const list = [...(this.tasksByCorrelation.get(key) ?? [])];
    const idx = list.findIndex((t) => t.id === task.id);
    if (idx >= 0) list[idx] = task;
    else list.push(task);
    const next = new Map(this.tasksByCorrelation);
    next.set(key, list);
    this.tasksByCorrelation = next;
    this.rebuildSessions();
  }

  togglePin(sessionId: string): void {
    const next = new Set(this.pinnedSessionIds);
    if (next.has(sessionId)) next.delete(sessionId);
    else next.add(sessionId);
    this.pinnedSessionIds = next;
    try {
      localStorage.setItem(PINNED_KEY, JSON.stringify([...next]));
    } catch {
      // ignore
    }
  }

  isSessionPinned(sessionId: string): boolean {
    return this.pinnedSessionIds.has(sessionId);
  }

  setActiveSession(id: string | null): void {
    this.activeSessionId = id;
    try {
      if (id) localStorage.setItem(STORAGE_KEY, id);
      else localStorage.removeItem(STORAGE_KEY);
    } catch {
      // ignore
    }
  }

  private getSavedSessionId(): string | null {
    try {
      return localStorage.getItem(STORAGE_KEY);
    } catch {
      return null;
    }
  }

  async loadSessions(limit = 50): Promise<void> {
    this.isLoading = true;
    try {
      const client = connectionStore.getAgentClient();
      const res = await client.tasks.list({ per_page: limit });
      this.indexTasks(res.tasks);

      if (!this.activeSessionId && this.sessions.length > 0) {
        const savedId = this.getSavedSessionId();
        const target =
          savedId && this.sessions.find((s) => s.id === savedId)
            ? savedId
            : this.sessions[0].id;
        this.setActiveSession(target);
        await this.loadSessionHistory(target);
      }
    } catch (err) {
      logger.error("[SessionStore] Failed to load sessions:", err);
      toast.error("Failed to load sessions");
    } finally {
      this.isLoading = false;
    }
  }

  private loadSessionHistory(sessionId: string): void {
    const tasks = this.tasksByCorrelation.get(sessionId) ?? [];
    chatStore.loadHistory(tasksToMessages(tasks));
  }

  async switchSession(sessionId: string): Promise<void> {
    if (sessionId === this.activeSessionId) return;

    this.setActiveSession(sessionId);
    this.isLoadingHistory = true;

    try {
      const cached = this.tasksByCorrelation.get(sessionId);
      if (cached?.length) {
        this.loadSessionHistory(sessionId);
      } else {
        const client = connectionStore.getAgentClient();
        const task = await client.tasks.get(sessionId);
        this.upsertTask(task);
        this.loadSessionHistory(sessionId);
      }
    } catch (err) {
      console.error("[SessionStore] Failed to load session history:", err);
      chatStore.clearMessages();
    } finally {
      this.isLoadingHistory = false;
    }
  }

  async createNewSession(): Promise<void> {
    chatStore.clearMessages();
    this.setActiveSession(null);
  }

  async deleteSession(sessionId: string): Promise<void> {
    this.hiddenIds = new Set([...this.hiddenIds, sessionId]);
    this.persistHidden();
    this.rebuildSessions();

    if (this.activeSessionId === sessionId) {
      const next = this.sessions[0];
      if (next) await this.switchSession(next.id);
      else {
        this.setActiveSession(null);
        chatStore.clearMessages();
      }
    }

    if (chatStore.streamingTaskId === sessionId) {
      chatStore.stopGeneration();
    }
  }

  async renameSession(sessionId: string, title: string): Promise<void> {
    const trimmed = title.trim();
    if (!trimmed) return;

    this.titleOverrides = { ...this.titleOverrides, [sessionId]: trimmed };
    this.persistTitles();
    this.rebuildSessions();

    const session = this.sessions.find((s) => s.id === sessionId);
    if (session) session.title = trimmed;
  }

  async searchSessions(query: string): Promise<Session[]> {
    const q = query.trim().toLowerCase();
    if (!q) return [];

    const matches: Session[] = [];
    for (const [corrId, tasks] of this.tasksByCorrelation) {
      if (this.hiddenIds.has(corrId)) continue;
      const hit = tasks.some(
        (t) =>
          t.prompt.toLowerCase().includes(q) ||
          (t.result?.toLowerCase().includes(q) ?? false),
      );
      if (hit) {
        const session = this.sessions.find((s) => s.id === corrId);
        if (session) matches.push(session);
        else {
          const first = tasks[0];
          matches.push({
            id: corrId,
            title: this.titleOverrides[corrId] ?? truncateTitle(first?.prompt ?? q),
            channel: "task",
            last_activity: first?.updated_at ?? first?.created_at ?? new Date().toISOString(),
            message_count: tasks.length * 2,
          });
        }
      }
    }
    return matches;
  }

  async exportSession(
    sessionId: string,
    format: "json" | "md" = "json",
  ): Promise<string> {
    const messages = tasksToMessages(this.tasksByCorrelation.get(sessionId) ?? []);

    if (format === "md") {
      return messages
        .map((m) => `### ${m.role}\n\n${m.content}\n`)
        .join("\n");
    }
    return JSON.stringify({ session_id: sessionId, messages }, null, 2);
  }

  onTaskCreated(task: Task): void {
    this.upsertTask(task);
    const sid = taskStreamId(task);
    if (!this.activeSessionId) {
      this.setActiveSession(sid);
    }
  }

  onTaskUpdated(task: Task): void {
    this.upsertTask(task);
  }
}

export const sessionStore = new SessionStore();