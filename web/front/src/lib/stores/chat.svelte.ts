import type { ChatMessage, FileContext, MediaAttachment } from "$lib/api";
import type { TaskStream } from "@pocketagent/client";
import { isPocketAgentError, taskStreamId } from "@pocketagent/client";
import { friendlyErrorMessage } from "$lib/api/errors";
import { toast } from "svelte-sonner";
import { connectionStore } from "./connection.svelte";
import { sessionStore } from "./sessions.svelte";
import { agentsStore } from "./agents.svelte";
import { explorerStore } from "./explorer.svelte";
import { activityStore } from "./activity.svelte";
import { appendFileContextToPrompt } from "$lib/types/chat";

class ChatStore {
  messages = $state<ChatMessage[]>([]);
  isStreaming = $state(false);
  streamingContent = $state("");
  streamingStatus = $state<string | null>(null);
  error = $state<string | null>(null);
  streamingTaskId = $state<string | null>(null);

  /** Pending AskUserQuestion state (legacy UI — not emitted by task stream). */
  pendingAskQuestion = $state<string | null>(null);
  pendingAskOptions = $state<string[]>([]);

  isEmpty = $derived(this.messages.length === 0);
  lastMessage = $derived(this.messages.at(-1) ?? null);

  private taskStream: TaskStream | null = null;
  private abortController: AbortController | null = null;

  private getFileContext(): FileContext | undefined {
    const dir = explorerStore.currentPath;
    const file = explorerStore.openFile;
    const selected = explorerStore.selectedFiles;

    if (!dir && !file && selected.size === 0) return undefined;

    const ctx: FileContext = {};
    if (dir) ctx.current_dir = dir;
    if (file) {
      ctx.open_file = file.path;
      ctx.open_file_name = file.name;
      ctx.open_file_extension = file.extension;
      ctx.open_file_size = file.size;
    }
    if (selected.size > 0) ctx.selected_files = [...selected];
    const source = explorerStore.currentSource;
    if (source !== "local") ctx.source = source;
    return ctx;
  }

  sendMessage(content: string, media?: MediaAttachment[]): void {
    if (media?.length) {
      toast.info("File attachments are not yet supported with the task API");
    }

    const userMsg: ChatMessage = {
      role: "user",
      content,
      timestamp: new Date().toISOString(),
      media,
    };
    this.messages.push(userMsg);
    this.error = null;
    void this.runTask(content);
  }

  regenerateLastResponse(): void {
    if (this.isStreaming) return;

    let lastUserIdx = -1;
    for (let i = this.messages.length - 1; i >= 0; i--) {
      if (this.messages[i].role === "user") {
        lastUserIdx = i;
        break;
      }
    }
    if (lastUserIdx === -1) return;

    const userContent = this.messages[lastUserIdx].content;
    const userMedia = this.messages[lastUserIdx].media;

    this.messages = this.messages.slice(0, lastUserIdx + 1);
    this.error = null;
    void this.runTask(userContent, userMedia);
  }

  answerAskUser(answer: string): void {
    this.pendingAskQuestion = null;
    this.pendingAskOptions = [];
    this.sendMessage(answer);
  }

  stopGeneration(): void {
    this.taskStream?.stop();
    this.taskStream = null;
    this.abortController?.abort();
    this.abortController = null;

    const streamId = this.streamingTaskId;
    if (streamId) {
      try {
        const client = connectionStore.getAgentClient();
        client.tasks.cancel(streamId).catch(() => {
          // best-effort
        });
      } catch {
        // ignore
      }
    }

    this.finalizeStream();
  }

  loadHistory(messages: ChatMessage[]): void {
    this.taskStream?.stop();
    this.taskStream = null;
    this.abortController?.abort();
    this.abortController = null;

    this.messages = messages;
    this.isStreaming = false;
    this.streamingContent = "";
    this.streamingStatus = null;
    this.error = null;
    this.streamingTaskId = null;

    activityStore.isAgentWorking = false;
    activityStore.sseActive = false;
  }

  clearMessages(): void {
    this.taskStream?.stop();
    this.taskStream = null;
    this.abortController?.abort();
    this.abortController = null;

    this.messages = [];
    this.isStreaming = false;
    this.streamingContent = "";
    this.streamingStatus = null;
    this.error = null;
    this.streamingTaskId = null;

    activityStore.isAgentWorking = false;
    activityStore.sseActive = false;
  }

  private async runTask(content: string, _media?: MediaAttachment[]): Promise<void> {
    this.isStreaming = true;
    this.streamingContent = "";
    this.streamingStatus = "Starting task...";

    activityStore.clear();
    activityStore.isAgentWorking = true;
    activityStore.sseActive = true;

    this.abortController?.abort();
    this.abortController = new AbortController();

    try {
      const client = connectionStore.getAgentClient();
      const fileContext = this.getFileContext();
      const prompt = appendFileContextToPrompt(content, fileContext);

      const correlationId = sessionStore.activeSessionId ?? undefined;
      const agentId = agentsStore.selectedAgentId ?? undefined;

      const task = await client.tasks.create({
        prompt,
        agent_id: agentId,
        correlation_id: correlationId,
      });

      sessionStore.onTaskCreated(task);

      const streamId = taskStreamId(task);
      this.streamingTaskId = streamId;

      this.taskStream = client.tasks.openStream(streamId, {
        events: "all",
        signal: this.abortController.signal,
      });

      await this.taskStream.start({
        onConnected: () => {
          this.streamingStatus = "Thinking...";
        },
        onToken: (payload) => {
          if (this.streamingStatus) this.streamingStatus = null;
          this.streamingContent += payload.delta;
        },
        onEvent: (event) => {
          activityStore.pushTaskEvent(event);
          if (event.type === "subtask_started") {
            this.streamingStatus = "Running subtask...";
          } else if (event.type === "orchestrating") {
            this.streamingStatus = "Orchestrating...";
          }
        },
        onTerminal: (type, event) => {
          activityStore.pushTaskEvent(event);
          if (type === "failed" && event.message) {
            this.error = event.message;
            toast.error(event.message);
          }
          void this.refreshTaskAfterTerminal(streamId, event.result);
          this.finalizeStream();
        },
        onError: (err) => {
          const message = friendlyErrorMessage(err);
          this.error = message;
          toast.error(message);
          this.finalizeStream();
        },
        onClose: () => {
          if (this.isStreaming) this.finalizeStream();
        },
      });
    } catch (err: unknown) {
      if (err instanceof DOMException && err.name === "AbortError") {
        activityStore.sseActive = false;
        activityStore.isAgentWorking = false;
        return;
      }
      const message = friendlyErrorMessage(err);
      this.error = message;
      if (isPocketAgentError(err) && err.isUnauthorized) {
        connectionStore.logout();
      }
      toast.error(message);
      if (this.isStreaming) this.finalizeStream();
    } finally {
      this.abortController = null;
      this.taskStream = null;
    }
  }

  private async refreshTaskAfterTerminal(streamId: string, streamResult?: string): Promise<void> {
    try {
      const client = connectionStore.getAgentClient();
      const task = await client.tasks.get(streamId);
      if (!task.result && streamResult) {
        task.result = streamResult;
      }
      sessionStore.onTaskUpdated(task);
    } catch {
      // history refresh is best-effort
    }
  }

  private finalizeStream(): void {
    if (this.streamingContent) {
      const msg: ChatMessage = {
        role: "assistant",
        content: this.streamingContent,
        timestamp: new Date().toISOString(),
      };
      if (this.pendingAskOptions.length > 0) {
        msg.metadata = { askUser: true, options: [...this.pendingAskOptions] };
      }
      this.messages.push(msg);
    } else if (this.error && this.messages.at(-1)?.role === "user") {
      const msg: ChatMessage = {
        role: "assistant",
        content: `_${this.error}_`,
        timestamp: new Date().toISOString(),
      };
      this.messages.push(msg);
    }

    this.isStreaming = false;
    this.streamingContent = "";
    this.streamingStatus = null;
    this.streamingTaskId = null;
    this.pendingAskQuestion = null;
    this.pendingAskOptions = [];

    activityStore.isAgentWorking = false;
    activityStore.sseActive = false;
  }
}

export const chatStore = new ChatStore();