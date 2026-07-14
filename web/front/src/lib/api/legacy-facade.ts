import type { PocketPawClient } from "./client";

const EMPTY_LIST = { sessions: [], kits: [], catalog: [], skills: [], results: [], files: [], reminders: [] };
const EMPTY_MAP = {};
function stubValue(method: string): unknown {
  if (method === "listSessions") return { sessions: [] };
  if (method === "getSessionHistory") return [];
  if (method === "createSession") return { id: "local", title: "New Chat" };
  if (method === "getSettings") return null;
  if (method === "getHealth") return { status: "unknown", dependencies: {} };
  if (method === "getVersion") return { version: "web", agent_backend: "gate" };
  if (method === "getMcpStatus" || method === "getChannelStatus") return EMPTY_MAP;
  if (method === "getMcpPresets") return [];
  if (method === "listKits") return [];
  if (method === "listKitCatalog") return [];
  if (method === "listSkills" || method === "searchSkills") return [];
  if (method === "getRecentFiles") return [];
  if (method === "getLongTermMemory") return [];
  if (method === "getMemorySettings") return {};
  if (method === "getMemoryStats") return {};
  if (method === "getIdentity") return { soul: "", identity: "", user: "" };
  if (method === "getAuditLog") return [];
  if (method === "listBackends") return [];
  if (method === "fetchOllamaModels") return [];
  if (method === "getHealthErrors" || method === "getRecentUsage") return [];
  if (method === "getUsageSummary") return {};
  if (method === "getSystemMetrics") return {};
  if (method === "mcListAgents") return { agents: [], count: 0 };
  if (method === "mcListProjects") return { projects: [], count: 0 };
  if (method === "mcGetRunningTasks") return { tasks: [], count: 0 };
  if (method === "mcListNotifications") return { notifications: [], count: 0 };
  if (method === "mcGetStandup") return { entries: [] };
  if (method === "mcGetTaskMessages") return { messages: [], count: 0 };
  if (method === "getChannelConfig") return null;
  if (method === "checkExtra") return { installed: false };
  if (method === "getWhatsAppQR") return { qr: null, connected: false };
  if (method === "chat") return "";
  if (method === "chatStream") return undefined;
  if (method === "exportSession") return "";
  if (method.startsWith("mc") || method.startsWith("dw")) return {};
  if (method.includes("toggle") || method.includes("install") || method.includes("remove")) {
    return { status: "unavailable", error: "Not available in web UI yet" };
  }
  if (method.startsWith("list") || method.startsWith("get")) return EMPTY_LIST;
  return undefined;
}

/**
 * Runtime stub for legacy PocketPawClient call sites.
 * Typed as PocketPawClient so existing pages compile until Phase 2 migration.
 */
export function createLegacyFacade(): PocketPawClient {
  return new Proxy({} as PocketPawClient, {
    get(_target, prop) {
      if (prop === "setToken" || prop === "getApiBase" || prop === "getWsUrl") {
        return () => "";
      }
      const method = String(prop);
      return async (..._args: unknown[]) => stubValue(method);
    },
  });
}