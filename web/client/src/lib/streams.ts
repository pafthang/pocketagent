import { errorFromResponse, PocketAgentError } from './errors.js';
import { buildHeaders, createRequestContext } from './http.js';
import type {
  PocketAgentClientOptions,
  TaskEvent,
  TaskEventType,
  TaskTokenEvent,
} from './types.js';

const TERMINAL_EVENTS = new Set<TaskEventType>([
  'completed',
  'failed',
  'cancelled',
  'timeout',
]);

export interface TaskStreamOptions {
  events?: 'all';
  signal?: AbortSignal;
}

export interface TaskStreamHandlers {
  onConnected?: (payload: { task_id: string; status: string }) => void;
  onToken?: (payload: TaskTokenEvent) => void;
  onEvent?: (payload: TaskEvent) => void;
  onTerminal?: (type: TaskEventType, payload: TaskEvent) => void;
  onError?: (error: unknown) => void;
  onClose?: () => void;
}

export function taskStreamUrl(
  baseUrl: string,
  taskId: string,
  opts?: { token?: string; spaceId?: string; events?: 'all' },
): string {
  const url = new URL(`/tasks/${taskId}/stream`, baseUrl.replace(/\/$/, '') + '/');
  if (opts?.events === 'all') url.searchParams.set('events', 'all');
  if (opts?.token) url.searchParams.set('token', opts.token);
  if (opts?.spaceId) url.searchParams.set('space_id', opts.spaceId);
  return url.toString();
}

export function taskWebSocketUrl(
  baseUrl: string,
  taskId: string,
  opts?: { token?: string; spaceId?: string },
): string {
  const http = baseUrl.replace(/\/$/, '');
  const wsBase = http.replace(/^http/, 'ws');
  const url = new URL(`/ws/task/${taskId}`, `${wsBase}/`);
  if (opts?.token) url.searchParams.set('token', opts.token);
  if (opts?.spaceId) url.searchParams.set('space_id', opts.spaceId);
  return url.toString();
}

export class TaskStream {
  private abortController: AbortController | null = null;
  private running = false;

  constructor(
    private readonly clientOpts: PocketAgentClientOptions,
    private readonly taskId: string,
    private readonly opts?: TaskStreamOptions,
  ) {}

  get isRunning(): boolean {
    return this.running;
  }

  async start(handlers: TaskStreamHandlers): Promise<void> {
    this.stop();
    this.running = true;
    this.abortController = new AbortController();

    const signal = this.opts?.signal;
    if (signal) {
      signal.addEventListener('abort', () => this.abortController?.abort(signal.reason), {
        once: true,
      });
    }

    try {
      await streamTaskSSE(this.clientOpts, this.taskId, handlers, {
        ...this.opts,
        signal: this.abortController.signal,
      });
    } finally {
      this.running = false;
      handlers.onClose?.();
    }
  }

  stop(): void {
    this.abortController?.abort();
    this.abortController = null;
    this.running = false;
  }
}

export async function streamTaskSSE(
  clientOpts: PocketAgentClientOptions,
  taskId: string,
  handlers: TaskStreamHandlers,
  streamOpts?: TaskStreamOptions,
): Promise<void> {
  const ctx = createRequestContext(clientOpts);
  const url = taskStreamUrl(ctx.baseUrl, taskId, {
    token: ctx.token,
    spaceId: ctx.spaceId,
    events: streamOpts?.events,
  });

  const headers = buildHeaders(ctx, { Accept: 'text/event-stream' });
  let res: Response;
  try {
    res = await ctx.fetchFn(url, {
      headers,
      signal: streamOpts?.signal,
    });
  } catch (err) {
    if ((err as Error).name !== 'AbortError') handlers.onError?.(err);
    return;
  }

  if (!res.ok) {
    handlers.onError?.(await errorFromResponse(res));
    return;
  }
  if (!res.body) {
    handlers.onError?.(new PocketAgentError('empty SSE body', res.status));
    return;
  }

  const reader = res.body.getReader();
  const decoder = new TextDecoder();
  let buffer = '';

  try {
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      buffer += decoder.decode(value, { stream: true });
      const chunks = buffer.split('\n\n');
      buffer = chunks.pop() ?? '';

      for (const chunk of chunks) {
        dispatchSSEChunk(chunk, handlers);
      }
    }
  } catch (err) {
    if ((err as Error).name !== 'AbortError') {
      handlers.onError?.(err);
    }
  }
}

export interface ProjectStreamHandlers {
  onEvent?: (payload: Record<string, unknown>) => void;
  onPhase?: (payload: { project_id: string; phase: string; message: string }) => void;
  onComplete?: (payload: { project_id: string; error?: string }) => void;
  onError?: (error: unknown) => void;
  onClose?: () => void;
}

export function projectWebSocketUrl(
  baseUrl: string,
  projectId: string,
  opts?: { token?: string; spaceId?: string },
): string {
  const http = baseUrl.replace(/\/$/, '');
  const wsBase = http.replace(/^http/, 'ws');
  const url = new URL(`/ws/project/${projectId}`, `${wsBase}/`);
  if (opts?.token) url.searchParams.set('token', opts.token);
  if (opts?.spaceId) url.searchParams.set('space_id', opts.spaceId);
  return url.toString();
}

export function connectProjectWebSocket(
  clientOpts: PocketAgentClientOptions,
  projectId: string,
  handlers: ProjectStreamHandlers,
): WebSocket {
  const ctx = createRequestContext(clientOpts);
  const url = projectWebSocketUrl(ctx.baseUrl, projectId, {
    token: ctx.token,
    spaceId: ctx.spaceId,
  });

  const ws = new WebSocket(url);

  ws.onmessage = (msg) => {
    try {
      const event = JSON.parse(String(msg.data)) as Record<string, unknown>;
      handlers.onEvent?.(event);

      switch (event.type) {
        case 'dw_planning_phase':
          handlers.onPhase?.(event as { project_id: string; phase: string; message: string });
          break;
        case 'dw_planning_complete':
          handlers.onComplete?.(event as { project_id: string; error?: string });
          ws.close();
          break;
      }
    } catch (err) {
      handlers.onError?.(err);
    }
  };

  ws.onerror = () => handlers.onError?.(new Error('WebSocket error'));
  ws.onclose = () => handlers.onClose?.();
  return ws;
}

export function connectTaskWebSocket(
  clientOpts: PocketAgentClientOptions,
  taskId: string,
  handlers: TaskStreamHandlers,
): WebSocket {
  const ctx = createRequestContext(clientOpts);
  const url = taskWebSocketUrl(ctx.baseUrl, taskId, {
    token: ctx.token,
    spaceId: ctx.spaceId,
  });

  const ws = new WebSocket(url);

  ws.onmessage = (msg) => {
    try {
      const event = JSON.parse(String(msg.data)) as TaskEvent;
      handlers.onEvent?.(event);

      if (event.type === 'llm_token' && event.message) {
        handlers.onToken?.({
          task_id: event.task_id,
          step: event.step,
          delta: event.message,
        });
      }

      if (TERMINAL_EVENTS.has(event.type)) {
        handlers.onTerminal?.(event.type, event);
        ws.close();
      }
    } catch (err) {
      handlers.onError?.(err);
    }
  };

  ws.onerror = () => handlers.onError?.(new Error('WebSocket error'));
  ws.onclose = () => handlers.onClose?.();
  return ws;
}

/** Parse one SSE frame (testable). */
export function parseSSEFrame(chunk: string): { event: string; data: string } | null {
  if (!chunk.trim() || chunk.startsWith(':')) return null;

  let event = 'message';
  let data = '';

  for (const line of chunk.split('\n')) {
    if (line.startsWith('event:')) {
      event = line.slice(6).trim();
    } else if (line.startsWith('data:')) {
      data += line.slice(5).trim();
    }
  }

  return data ? { event, data } : null;
}

function dispatchSSEChunk(chunk: string, handlers: TaskStreamHandlers): void {
  const frame = parseSSEFrame(chunk);
  if (!frame) return;

  try {
    const payload = JSON.parse(frame.data) as Record<string, unknown>;

    switch (frame.event) {
      case 'connected':
        handlers.onConnected?.(payload as unknown as { task_id: string; status: string });
        break;
      case 'token':
        handlers.onToken?.(payload as unknown as TaskTokenEvent);
        break;
      default:
        handlers.onEvent?.(payload as unknown as TaskEvent);
        if (TERMINAL_EVENTS.has(frame.event as TaskEventType)) {
          handlers.onTerminal?.(frame.event as TaskEventType, payload as unknown as TaskEvent);
        }
    }
  } catch (err) {
    handlers.onError?.(err);
  }
}

export function isTerminalEvent(type: TaskEventType): boolean {
  return TERMINAL_EVENTS.has(type);
}