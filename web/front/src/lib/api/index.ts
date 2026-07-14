export * from "./types";
export { friendlyErrorMessage } from "./errors";
export { GATE_URL, BACKEND_URL, API_PREFIX, API_BASE } from "./config";
export { PocketPawClient } from "./client";
export { PocketPawWebSocket } from "./websocket";
export type { ConnectionState } from "./websocket";
export {
  createClient,
  PocketAgentClient,
  PocketAgentError,
  isPocketAgentError,
  isUnauthorized,
  isForbidden,
  isNotFound,
} from "@pocketagent/client";
export type {
  Agent,
  AuthSession,
  AuthUser,
  RegisterResponse,
  Space,
  SpaceRole,
  Task,
  TaskEvent,
  TaskStatus,
} from "@pocketagent/client";