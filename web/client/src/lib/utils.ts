import type { ListParams, Task, TaskStatus } from './types.js';

const TERMINAL_STATUSES = new Set<TaskStatus>([
  'completed',
  'degraded',
  'failed',
  'cancelled',
]);

export function buildQuery(
  path: string,
  params?: ListParams | Record<string, string | number | boolean | undefined>,
): string {
  if (!params) return path;

  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === '') continue;
    search.set(key, String(value));
  }

  const qs = search.toString();
  return qs ? `${path}?${qs}` : path;
}

/** Preferred task identifier for API paths (correlation_id). */
export function taskStreamId(task: Pick<Task, 'correlation_id' | 'id'>): string {
  return task.correlation_id || task.id;
}

export function isTaskTerminal(status?: TaskStatus): boolean {
  return status ? TERMINAL_STATUSES.has(status) : false;
}

export function sleep(ms: number, signal?: AbortSignal): Promise<void> {
  return new Promise((resolve, reject) => {
    if (signal?.aborted) {
      reject(signal.reason);
      return;
    }
    const timer = setTimeout(resolve, ms);
    signal?.addEventListener(
      'abort',
      () => {
        clearTimeout(timer);
        reject(signal.reason);
      },
      { once: true },
    );
  });
}