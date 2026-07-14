import type { ClientCore } from '../core.js';
import type {
  CreateTaskInput,
  ListParams,
  RequestConfig,
  Task,
  TaskWithSubtasks,
} from '../types.js';
import { buildQuery } from '../utils.js';
import { TaskStream, type TaskStreamOptions } from '../streams.js';

export class TasksAPI {
  constructor(private readonly core: ClientCore) {}

  list(params: ListParams = {}, config?: RequestConfig): Promise<{ tasks: Task[]; total: number }> {
    return this.core.request('GET', buildQuery('/tasks', params), undefined, { config });
  }

  get(id: string, config?: RequestConfig): Promise<Task> {
    return this.core.request('GET', `/tasks/${id}`, undefined, { config });
  }

  getWithSubtasks(id: string, config?: RequestConfig): Promise<TaskWithSubtasks> {
    return this.core.request('GET', `/tasks/${id}?include=subtasks`, undefined, { config });
  }

  create(input: CreateTaskInput, config?: RequestConfig): Promise<Task> {
    return this.core.request('POST', '/tasks', input, { config });
  }

  cancel(id: string, config?: RequestConfig): Promise<Task> {
    return this.core.request('DELETE', `/tasks/${id}`, undefined, { config });
  }

  openStream(taskId: string, opts?: TaskStreamOptions): TaskStream {
    return new TaskStream(
      {
        baseUrl: this.core.baseUrl,
        token: this.core.token,
        spaceId: this.core.spaceId,
        fetch: this.core.ctx.fetchFn,
        timeoutMs: this.core.ctx.timeoutMs,
      },
      taskId,
      opts,
    );
  }
}