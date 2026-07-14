export { PocketAgentClient, createClient } from './client.js';
export { ClientCore } from './core.js';
export {
  PocketAgentError,
  errorFromResponse,
  isForbidden,
  isNotFound,
  isPocketAgentError,
  isUnauthorized,
} from './errors.js';
export { HEADER_SPACE_ID } from './http.js';
export {
  TaskStream,
  connectProjectWebSocket,
  connectTaskWebSocket,
  isTerminalEvent,
  parseSSEFrame,
  projectWebSocketUrl,
  streamTaskSSE,
  taskStreamUrl,
  taskWebSocketUrl,
  type ProjectStreamHandlers,
  type TaskStreamHandlers,
  type TaskStreamOptions,
} from './streams.js';
export { SkillsAPI } from './api/skills.js';
export { ToolsAPI } from './api/tools.js';
export { McpAPI } from './api/mcp.js';
export { MemoryAPI } from './api/memory.js';
export { IdentityAPI } from './api/identity.js';
export { ActivityAPI } from './api/activity.js';
export { DashboardAPI } from './api/dashboard.js';
export { ProjectsAPI } from './api/projects.js';
export { FilesAPI } from './api/files.js';
export { buildQuery, isTaskTerminal, sleep, taskStreamId } from './utils.js';
export type {
  CreateSkillInput,
  PatchSkillInput,
  RunSkillInput,
  RunSkillResult,
  Skill,
  SkillCatalogEntry,
  SkillListResult,
} from './api/skills.js';
export type { ToolInfo } from './api/tools.js';
export type {
  CreateMCPServerInput,
  MCPPreset,
  MCPProbeResult,
  MCPServer,
  MCPServerStatus,
  PatchMCPServerInput,
} from './api/mcp.js';
export type {
  IngestMemoryInput,
  MemoryDocument,
  MemoryListResult,
  MemorySearchResult,
  MemoryStats,
  SearchMemoryInput,
} from './api/memory.js';
export type {
  IdentityFiles,
  IdentitySaveResponse,
  UpdateIdentityInput,
} from './api/identity.js';
export type {
  ActivityEntry,
  ActivityEntryType,
  ActivityListResult,
} from './api/activity.js';
export type {
  BuiltinKit,
  DashboardAgent,
  DashboardFeedItem,
  DashboardKanbanCard,
  DashboardListParams,
  DashboardMetrics,
  DashboardSummary,
  DashboardTaskRow,
} from './api/dashboard.js';
export type {
  CreateProjectInput,
  CreateProjectItemInput,
  PatchProjectInput,
  PatchProjectItemInput,
  Project,
  ProjectItem,
  ProjectItemStatus,
  ProjectProgress,
  ProjectStatus,
} from './api/projects.js';
export type {
  BrowseFileEntry,
  BrowseFilesParams,
  BrowseFilesResult,
  CreateFolderInput,
  FileContentPreview,
  IngestFileInput,
  IngestFileResult,
  RecentFileEntry,
  StoredFile,
  UploadFileInput,
} from './api/files.js';
export { projectFilesPath } from './api/files.js';
export type {
  Agent,
  ApiErrorBody,
  AuditLog,
  AuthSession,
  AuthUser,
  AuthorizeRequest,
  AuthorizeResponse,
  CreateAgentInput,
  CreateInviteResponse,
  CreateScheduleInput,
  CreateTaskInput,
  InvitePreview,
  ListParams,
  PaginatedResult,
  PocketAgentClientOptions,
  RegisterResponse,
  RequestConfig,
  Schedule,
  Space,
  SpaceInvite,
  SpaceMember,
  SpaceRole,
  Task,
  TaskEvent,
  TaskEventType,
  TaskStatus,
  TaskTokenEvent,
  TaskWithSubtasks,
  TaskWorkflow,
  Team,
  TeamMember,
} from './types.js';