export type TaskStatus =
  | 'queued'
  | 'running'
  | 'completed'
  | 'degraded'
  | 'failed'
  | 'cancelled';

export type TaskWorkflow = 'default' | 'supervisor';

export type SpaceRole = 'admin' | 'editor' | 'viewer';

export type TaskEventType =
  | 'connected'
  | 'queued'
  | 'orchestrating'
  | 'subtask_dispatched'
  | 'subtask_started'
  | 'subtask_completed'
  | 'subtask_result'
  | 'completed'
  | 'failed'
  | 'timeout'
  | 'llm_token'
  | 'cancelled'
  | 'supervisor_delegated';

export interface AuthUser {
  id: string;
  email: string;
  verified: boolean;
}

export interface AuthSession {
  token: string;
  user: AuthUser;
}

export interface RegisterResponse {
  user: AuthUser;
  email_verification_required: boolean;
}

export interface Space {
  id: string;
  name: string;
  slug: string;
  description?: string;
  is_system: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface SpaceMember {
  id: string;
  space_id: string;
  user_id: string;
  role: SpaceRole;
  created_at?: string;
  updated_at?: string;
}

export interface Team {
  id: string;
  space_id: string;
  name: string;
  description?: string;
  created_at?: string;
  updated_at?: string;
}

export interface TeamMember {
  id: string;
  team_id: string;
  member_type: 'user' | 'agent';
  member_id: string;
  created_at?: string;
  updated_at?: string;
}

export interface SpaceInvite {
  id: string;
  space_id: string;
  email: string;
  role: SpaceRole;
  status: string;
  invited_by?: string;
  expires_at?: string;
  created_at?: string;
  updated_at?: string;
}

export interface InvitePreview {
  space_id: string;
  space_name: string;
  email: string;
  role: SpaceRole;
  status: string;
  expires_at?: string;
}

export interface CreateInviteResponse {
  invite: SpaceInvite;
  token: string;
  invite_url?: string;
}

export interface AuditLog {
  id: string;
  space_id: string;
  actor_id?: string;
  actor_email?: string;
  action: string;
  resource_type?: string;
  resource_id?: string;
  metadata?: Record<string, unknown>;
  ip_address?: string;
  created_at?: string;
}

export interface AuthorizeRequest {
  space_id: string;
  action: string;
  resource_type?: string;
  resource_id?: string;
}

export interface AuthorizeResponse {
  allowed: boolean;
  role?: SpaceRole;
  reason?: string;
}

export interface Agent {
  id: string;
  space_id?: string;
  name: string;
  description: string;
  model: string;
  system_prompt: string;
  tools: string[];
  config: Record<string, unknown>;
  created_at: string;
  updated_at: string;
}

export interface Task {
  id: string;
  correlation_id?: string;
  space_id?: string;
  agent_id?: string;
  prompt: string;
  status?: TaskStatus;
  result?: string;
  error?: string;
  parent_id?: string | null;
  workflow?: TaskWorkflow;
  worker_agent_ids?: string[];
  tools?: string[];
  skill_id?: string;
  created_at?: string;
  updated_at?: string;
}

export interface TaskWithSubtasks {
  task: Task;
  subtasks: Task[];
}

export interface Schedule {
  id: string;
  space_id: string;
  name: string;
  agent_id?: string;
  prompt: string;
  cron_expr: string;
  workflow?: TaskWorkflow;
  worker_agent_ids?: string[];
  enabled: boolean;
  last_run_at?: string;
  next_run_at?: string;
  last_task_id?: string;
  created_at?: string;
  updated_at?: string;
}

export interface TaskEvent {
  task_id: string;
  type: TaskEventType;
  status: string;
  step?: number;
  message?: string;
  result?: string;
}

export interface TaskTokenEvent {
  task_id: string;
  step?: number;
  delta: string;
}

export interface ListParams {
  page?: number;
  per_page?: number;
}

export interface PaginatedResult<T> {
  items: T[];
  total: number;
}

export interface ApiErrorBody {
  error: string;
}

export interface CreateTaskInput {
  prompt: string;
  agent_id?: string;
  correlation_id?: string;
  workflow?: TaskWorkflow;
  worker_agent_ids?: string[];
  tools?: string[];
  skill_id?: string;
}

export interface CreateAgentInput {
  name: string;
  description?: string;
  model?: string;
  system_prompt?: string;
  tools?: string[];
  config?: Record<string, unknown>;
}

export interface CreateScheduleInput {
  name: string;
  prompt: string;
  cron_expr: string;
  agent_id?: string;
  workflow?: TaskWorkflow;
  worker_agent_ids?: string[];
  enabled?: boolean;
}

export interface RequestConfig {
  signal?: AbortSignal;
  timeoutMs?: number;
  skipAuthRefresh?: boolean;
}

export interface PocketAgentClientOptions {
  baseUrl: string;
  token?: string;
  spaceId?: string;
  fetch?: typeof fetch;
  /** Default request timeout (ms). */
  timeoutMs?: number;
  /** Called when login/refresh updates the session token. */
  onSession?: (session: AuthSession) => void;
}