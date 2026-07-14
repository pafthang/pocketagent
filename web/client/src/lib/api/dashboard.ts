import type { ClientCore } from '../core.js';
import type { RequestConfig, Task } from '../types.js';
import { buildQuery } from '../utils.js';

export interface DashboardMetrics {
  agents_total: number;
  tasks_total: number;
  tasks_queued: number;
  tasks_running: number;
  tasks_completed: number;
  tasks_failed: number;
  tasks_cancelled: number;
}

export interface DashboardAgent {
  id: string;
  name: string;
  role: string;
  description: string;
  model?: string;
  status: string;
  level: string;
  specialties: string[];
  current_task_id?: string;
  created_at?: string;
  updated_at?: string;
}

export interface DashboardKanbanCard {
  id: string;
  title: string;
  status: string;
  priority?: string;
  agent_id?: string;
  created_at?: string;
  updated_at?: string;
}

export interface DashboardTaskRow {
  id: string;
  title: string;
  status: string;
  agent_id?: string;
  created_at?: string;
  updated_at?: string;
}

export interface DashboardFeedItem {
  message: string;
  type?: string;
  created_at?: string;
  agent_id?: string;
  task_id?: string;
}

export interface DashboardSummary {
  metrics: DashboardMetrics;
  agents: DashboardAgent[];
  running_tasks: Task[];
  recent_tasks: DashboardTaskRow[];
  kanban: Record<string, DashboardKanbanCard[]>;
  activity: DashboardFeedItem[];
  generated_at: string;
}

export interface BuiltinKit {
  id: string;
  config: Record<string, unknown>;
  user_values: Record<string, string>;
  installed_at: string;
  active: boolean;
}

export interface DashboardListParams {
  limit?: number;
}

export class DashboardAPI {
  constructor(private readonly core: ClientCore) {}

  summary(params: DashboardListParams = {}, config?: RequestConfig): Promise<DashboardSummary> {
    return this.core.request('GET', buildQuery('/dashboard', params), undefined, { config });
  }

  listKits(config?: RequestConfig): Promise<{ kits: BuiltinKit[] }> {
    return this.core.request('GET', '/kits', undefined, { config });
  }

  kitData(kitId: string, params: DashboardListParams = {}, config?: RequestConfig): Promise<Record<string, unknown>> {
    return this.core.request(
      'GET',
      buildQuery(`/kits/${encodeURIComponent(kitId)}/data`, params),
      undefined,
      { config },
    );
  }

  kitCatalog(config?: RequestConfig): Promise<{ catalog: Record<string, unknown>[] }> {
    return this.core.request('GET', '/kits/catalog', undefined, { config });
  }

  activateKit(kitId: string, config?: RequestConfig): Promise<{ ok: boolean; id: string; active: boolean }> {
    return this.core.request('POST', `/kits/${encodeURIComponent(kitId)}/activate`, undefined, { config });
  }
}