import { TaskStream } from '../streams.js';
import type {
  PocketAgentClientOptions,
  Task,
  TaskEvent,
  TaskEventType,
} from '../types.js';
import { taskStreamId } from '../utils.js';

export type TaskStreamStatus = 'idle' | 'streaming' | 'done' | 'error';

export interface TaskStreamStateOptions extends PocketAgentClientOptions {
  taskId: string;
  events?: 'all';
}

export function createTaskStreamState(opts: TaskStreamStateOptions) {
  let status = $state<TaskStreamStatus>('idle');
  let tokens = $state('');
  let events = $state<TaskEvent[]>([]);
  let lastEvent = $state<TaskEvent | null>(null);
  let error = $state<string | null>(null);

  let stream: TaskStream | null = null;

  const isStreaming = $derived(status === 'streaming');
  const isDone = $derived(status === 'done');

  async function start(): Promise<void> {
    stop();
    status = 'streaming';
    error = null;
    tokens = '';
    events = [];
    lastEvent = null;

    stream = new TaskStream(
      {
        baseUrl: opts.baseUrl,
        token: opts.token,
        spaceId: opts.spaceId,
        fetch: opts.fetch,
        timeoutMs: opts.timeoutMs,
      },
      opts.taskId,
      { events: opts.events },
    );

    await stream.start({
      onToken: (payload) => {
        tokens += payload.delta;
      },
      onEvent: (event) => {
        events = [...events, event];
        lastEvent = event;
      },
      onTerminal: (_type: TaskEventType, event) => {
        events = [...events, event];
        lastEvent = event;
        status = 'done';
      },
      onError: (err) => {
        error = String(err);
        status = 'error';
      },
      onClose: () => {
        if (status === 'streaming') status = 'done';
      },
    });
  }

  function stop(): void {
    stream?.stop();
    stream = null;
    if (status === 'streaming') status = 'idle';
  }

  function reset(): void {
    stop();
    tokens = '';
    events = [];
    lastEvent = null;
    error = null;
    status = 'idle';
  }

  return {
    get status() {
      return status;
    },
    get tokens() {
      return tokens;
    },
    get events() {
      return events;
    },
    get lastEvent() {
      return lastEvent;
    },
    get error() {
      return error;
    },
    get isStreaming() {
      return isStreaming;
    },
    get isDone() {
      return isDone;
    },
    start,
    stop,
    reset,
  };
}

export type TaskStreamState = ReturnType<typeof createTaskStreamState>;

export interface TaskRunStateOptions extends PocketAgentClientOptions {
  events?: 'all';
}

/** Create a task and stream its LLM output in one reactive helper. */
export function createTaskRunState(
  createTask: (prompt: string) => Promise<Task | undefined>,
  clientOpts: TaskRunStateOptions,
) {
  let task = $state<Task | null>(null);
  let status = $state<TaskStreamStatus>('idle');
  let tokens = $state('');
  let events = $state<TaskEvent[]>([]);
  let error = $state<string | null>(null);

  let stream: TaskStream | null = null;

  const isStreaming = $derived(status === 'streaming');
  const isDone = $derived(status === 'done');

  async function run(prompt: string): Promise<Task | undefined> {
    reset();
    status = 'streaming';
    error = null;

    const created = await createTask(prompt);
    if (!created) {
      status = 'error';
      return undefined;
    }

    task = created;
    const id = taskStreamId(created);

    stream = new TaskStream({ ...clientOpts }, id, { events: clientOpts.events });
    await stream.start({
      onToken: (payload) => {
        tokens += payload.delta;
      },
      onEvent: (event) => {
        events = [...events, event];
      },
      onTerminal: () => {
        status = 'done';
      },
      onError: (err) => {
        error = String(err);
        status = 'error';
      },
      onClose: () => {
        if (status === 'streaming') status = 'done';
      },
    });

    return created;
  }

  function stop(): void {
    stream?.stop();
    stream = null;
    if (status === 'streaming') status = 'idle';
  }

  function reset(): void {
    stop();
    task = null;
    tokens = '';
    events = [];
    error = null;
    status = 'idle';
  }

  return {
    get task() {
      return task;
    },
    get status() {
      return status;
    },
    get tokens() {
      return tokens;
    },
    get events() {
      return events;
    },
    get error() {
      return error;
    },
    get isStreaming() {
      return isStreaming;
    },
    get isDone() {
      return isDone;
    },
    run,
    stop,
    reset,
  };
}

export type TaskRunState = ReturnType<typeof createTaskRunState>;