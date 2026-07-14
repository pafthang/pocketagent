import { AgentsAPI } from './api/agents.js';
import { AuthAPI } from './api/auth.js';
import { SchedulesAPI } from './api/schedules.js';
import { SpacesAPI } from './api/spaces.js';
import { McpAPI } from './api/mcp.js';
import { MemoryAPI } from './api/memory.js';
import { SkillsAPI } from './api/skills.js';
import { TasksAPI } from './api/tasks.js';
import { ToolsAPI } from './api/tools.js';
import { IdentityAPI } from './api/identity.js';
import { ActivityAPI } from './api/activity.js';
import { DashboardAPI } from './api/dashboard.js';
import { ProjectsAPI } from './api/projects.js';
import { FilesAPI } from './api/files.js';
import { ClientCore } from './core.js';
import type { AuthSession, PocketAgentClientOptions } from './types.js';

export class PocketAgentClient {
  private readonly core: ClientCore;

  readonly auth: AuthAPI;
  readonly spaces: SpacesAPI;
  readonly agents: AgentsAPI;
  readonly tasks: TasksAPI;
  readonly schedules: SchedulesAPI;
  readonly memory: MemoryAPI;
  readonly mcp: McpAPI;
  readonly skills: SkillsAPI;
  readonly tools: ToolsAPI;
  readonly identity: IdentityAPI;
  readonly activity: ActivityAPI;
  readonly dashboard: DashboardAPI;
  readonly projects: ProjectsAPI;
  readonly files: FilesAPI;

  constructor(options: PocketAgentClientOptions) {
    this.core = new ClientCore(options);
    this.auth = new AuthAPI(this.core);
    this.spaces = new SpacesAPI(this.core);
    this.agents = new AgentsAPI(this.core);
    this.tasks = new TasksAPI(this.core);
    this.schedules = new SchedulesAPI(this.core);
    this.memory = new MemoryAPI(this.core);
    this.mcp = new McpAPI(this.core);
    this.skills = new SkillsAPI(this.core);
    this.tools = new ToolsAPI(this.core);
    this.identity = new IdentityAPI(this.core);
    this.activity = new ActivityAPI(this.core);
    this.dashboard = new DashboardAPI(this.core);
    this.projects = new ProjectsAPI(this.core);
    this.files = new FilesAPI(this.core);
    this.auth.wireRefresh();
  }

  get baseUrl(): string {
    return this.core.baseUrl;
  }

  get token(): string | undefined {
    return this.core.token;
  }

  get spaceId(): string | undefined {
    return this.core.spaceId;
  }

  set spaceId(id: string | undefined) {
    this.core.setSpaceId(id);
  }

  set token(value: string | undefined) {
    this.core.setToken(value);
  }

  withToken(token?: string): PocketAgentClient {
    return new PocketAgentClient({
      baseUrl: this.core.baseUrl,
      token,
      spaceId: this.core.spaceId,
      fetch: this.core.ctx.fetchFn,
      timeoutMs: this.core.ctx.timeoutMs,
      onSession: this.core.ctx.onSession,
    });
  }

  withSpace(spaceId?: string): PocketAgentClient {
    return new PocketAgentClient({
      baseUrl: this.core.baseUrl,
      token: this.core.token,
      spaceId,
      fetch: this.core.ctx.fetchFn,
      timeoutMs: this.core.ctx.timeoutMs,
      onSession: this.core.ctx.onSession,
    });
  }

  applySession(session: AuthSession): void {
    this.core.applySession(session);
  }
}

export function createClient(options: PocketAgentClientOptions): PocketAgentClient {
  return new PocketAgentClient(options);
}