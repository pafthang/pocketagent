import type { ClientCore } from '../core.js';
import type { ListParams, RequestConfig } from '../types.js';
import { buildQuery } from '../utils.js';

export type ProjectStatus =
  | 'draft'
  | 'planning'
  | 'awaiting_approval'
  | 'approved'
  | 'executing'
  | 'paused'
  | 'completed'
  | 'failed'
  | 'cancelled';

export type ProjectItemStatus =
  | 'inbox'
  | 'assigned'
  | 'in_progress'
  | 'review'
  | 'done'
  | 'blocked'
  | 'skipped';

export interface Project {
  id: string;
  space_id: string;
  title: string;
  goal?: string;
  description?: string;
  status: ProjectStatus;
  plan_json?: Record<string, unknown>;
  parent_task_id?: string;
  creator_id?: string;
  planner_agent_id?: string;
  team_agent_ids?: string[];
  tags?: string[];
  started_at?: string;
  completed_at?: string;
  created_at?: string;
  updated_at?: string;
  metadata?: Record<string, unknown>;
}

export interface ProjectItem {
  id: string;
  space_id: string;
  project_id: string;
  title: string;
  description?: string;
  status: ProjectItemStatus;
  priority?: string;
  assignee_ids?: string[];
  execution_task_id?: string;
  sort_order?: number;
  tags?: string[];
  created_at?: string;
  updated_at?: string;
}

export interface ProjectProgress {
  total: number;
  completed: number;
  skipped: number;
  in_progress: number;
  blocked: number;
  human_pending: number;
  percent: number;
}

export interface CreateProjectInput {
  title?: string;
  goal?: string;
  description?: string;
  planner_agent_id?: string;
  team_agent_ids?: string[];
  tags?: string[];
  status?: ProjectStatus;
}

export interface PatchProjectInput {
  title?: string;
  goal?: string;
  description?: string;
  status?: ProjectStatus;
  plan_json?: Record<string, unknown>;
  parent_task_id?: string;
  planner_agent_id?: string;
  team_agent_ids?: string[];
  tags?: string[];
  metadata?: Record<string, unknown>;
}

export interface CreateProjectItemInput {
  title: string;
  description?: string;
  status?: ProjectItemStatus;
  priority?: string;
  assignee_ids?: string[];
  execution_task_id?: string;
  sort_order?: number;
  tags?: string[];
}

export interface PatchProjectItemInput {
  title?: string;
  description?: string;
  status?: ProjectItemStatus;
  priority?: string;
  assignee_ids?: string[];
  execution_task_id?: string;
  sort_order?: number;
  tags?: string[];
}

export class ProjectsAPI {
  constructor(private readonly core: ClientCore) {}

  list(params: ListParams & { status?: string; format?: 'mc' } = {}, config?: RequestConfig) {
    return this.core.request<{ projects: Project[]; total: number }>(
      'GET',
      buildQuery('/projects', params),
      undefined,
      { config },
    );
  }

  listMC(params: ListParams & { status?: string } = {}, config?: RequestConfig) {
    return this.core.request<{ projects: Record<string, unknown>[]; count: number }>(
      'GET',
      buildQuery('/projects', { ...params, format: 'mc' }),
      undefined,
      { config },
    );
  }

  create(input: CreateProjectInput, config?: RequestConfig): Promise<Project> {
    return this.core.request('POST', '/projects', input, { config });
  }

  createMC(input: CreateProjectInput, config?: RequestConfig) {
    return this.core.request<{ project: Record<string, unknown> }>(
      'POST',
      '/projects?format=mc',
      input,
      { config },
    );
  }

  get(id: string, config?: RequestConfig) {
    return this.core.request<{ project: Project; items: ProjectItem[] }>(
      'GET',
      `/projects/${encodeURIComponent(id)}`,
      undefined,
      { config },
    );
  }

  getDetail(id: string, config?: RequestConfig) {
    return this.core.request<{
      project: Record<string, unknown>;
      tasks: Record<string, unknown>[];
      items: ProjectItem[];
      progress: ProjectProgress;
    }>('GET', `/projects/${encodeURIComponent(id)}?format=mc`, undefined, { config });
  }

  patch(id: string, input: PatchProjectInput, config?: RequestConfig): Promise<Project> {
    return this.core.request('PATCH', `/projects/${encodeURIComponent(id)}`, input, { config });
  }

  delete(id: string, config?: RequestConfig): Promise<void> {
    return this.core.request('DELETE', `/projects/${encodeURIComponent(id)}`, undefined, { config });
  }

  listItems(projectId: string, params: ListParams & { status?: string } = {}, config?: RequestConfig) {
    return this.core.request<{ items: ProjectItem[]; total: number }>(
      'GET',
      buildQuery(`/projects/${encodeURIComponent(projectId)}/items`, params),
      undefined,
      { config },
    );
  }

  createItem(projectId: string, input: CreateProjectItemInput, config?: RequestConfig): Promise<ProjectItem> {
    return this.core.request(
      'POST',
      `/projects/${encodeURIComponent(projectId)}/items`,
      input,
      { config },
    );
  }

  patchItem(
    projectId: string,
    itemId: string,
    input: PatchProjectItemInput,
    config?: RequestConfig,
  ): Promise<ProjectItem> {
    return this.core.request(
      'PATCH',
      `/projects/${encodeURIComponent(projectId)}/items/${encodeURIComponent(itemId)}`,
      input,
      { config },
    );
  }

  deleteItem(projectId: string, itemId: string, config?: RequestConfig): Promise<void> {
    return this.core.request(
      'DELETE',
      `/projects/${encodeURIComponent(projectId)}/items/${encodeURIComponent(itemId)}`,
      undefined,
      { config },
    );
  }

  parseGoal(description: string, config?: RequestConfig) {
    return this.core.request<{ success: boolean; goal_analysis: Record<string, unknown> }>(
      'POST',
      '/projects/parse-goal',
      { description },
      { config },
    );
  }

  start(description: string, input: Partial<CreateProjectInput> = {}, config?: RequestConfig) {
    return this.core.request<{ success: boolean; project_id: string; project: Project }>(
      'POST',
      '/projects/start',
      { description, ...input },
      { config },
    );
  }

  getPlan(id: string, config?: RequestConfig) {
    return this.core.request<Record<string, unknown>>(
      'GET',
      `/projects/${encodeURIComponent(id)}/plan`,
      undefined,
      { config },
    );
  }

  plan(id: string, config?: RequestConfig) {
    return this.core.request('POST', `/projects/${encodeURIComponent(id)}/plan`, undefined, { config });
  }

  approve(id: string, config?: RequestConfig) {
    return this.core.request('POST', `/projects/${encodeURIComponent(id)}/approve`, undefined, { config });
  }

  pause(id: string, config?: RequestConfig) {
    return this.core.request('POST', `/projects/${encodeURIComponent(id)}/pause`, undefined, { config });
  }

  resume(id: string, config?: RequestConfig) {
    return this.core.request('POST', `/projects/${encodeURIComponent(id)}/resume`, undefined, { config });
  }

  cancel(id: string, config?: RequestConfig) {
    return this.core.request('POST', `/projects/${encodeURIComponent(id)}/cancel`, undefined, { config });
  }

  skipTask(projectId: string, taskId: string, config?: RequestConfig) {
    return this.core.request(
      'POST',
      `/projects/${encodeURIComponent(projectId)}/tasks/${encodeURIComponent(taskId)}/skip`,
      undefined,
      { config },
    );
  }

  retryTask(projectId: string, taskId: string, config?: RequestConfig) {
    return this.core.request(
      'POST',
      `/projects/${encodeURIComponent(projectId)}/tasks/${encodeURIComponent(taskId)}/retry`,
      undefined,
      { config },
    );
  }
}