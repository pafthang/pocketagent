import type { ClientCore } from '../core.js';
import type { ListParams, RequestConfig } from '../types.js';
import { buildQuery } from '../utils.js';

export type ActivityEntryType = 'tool_start' | 'tool_result' | 'thinking' | 'error' | 'status';

export interface ActivityEntry {
  id: string;
  type: ActivityEntryType;
  content: string;
  data?: Record<string, unknown>;
  timestamp: string;
  source?: string;
  task_id?: string;
}

export interface ActivityListResult {
  entries: ActivityEntry[];
  total: number;
}

export class ActivityAPI {
  constructor(private readonly core: ClientCore) {}

  list(params: ListParams = {}, config?: RequestConfig): Promise<ActivityListResult> {
    const spaceId = this.core.requireSpaceId();
    return this.core.request(
      'GET',
      buildQuery(`/spaces/${encodeURIComponent(spaceId)}/activity`, params),
      undefined,
      { config },
    );
  }
}