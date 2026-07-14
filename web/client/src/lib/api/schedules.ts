import type { ClientCore } from '../core.js';
import type { CreateScheduleInput, ListParams, RequestConfig, Schedule } from '../types.js';
import { buildQuery } from '../utils.js';

export class SchedulesAPI {
  constructor(private readonly core: ClientCore) {}

  list(
    params: ListParams = {},
    config?: RequestConfig,
  ): Promise<{ schedules: Schedule[]; total: number }> {
    return this.core.request('GET', buildQuery('/schedules', params), undefined, { config });
  }

  get(id: string, config?: RequestConfig): Promise<Schedule> {
    return this.core.request('GET', `/schedules/${id}`, undefined, { config });
  }

  create(input: CreateScheduleInput, config?: RequestConfig): Promise<Schedule> {
    return this.core.request('POST', '/schedules', input, { config });
  }

  update(
    id: string,
    input: Partial<CreateScheduleInput>,
    config?: RequestConfig,
  ): Promise<Schedule> {
    return this.core.request('PATCH', `/schedules/${id}`, input, { config });
  }

  delete(id: string, config?: RequestConfig): Promise<void> {
    return this.core.request('DELETE', `/schedules/${id}`, undefined, { config });
  }
}