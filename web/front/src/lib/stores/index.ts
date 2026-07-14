import { connectionStore } from "./connection.svelte";
import { chatStore } from "./chat.svelte";
import { sessionStore } from "./sessions.svelte";
import { agentsStore } from "./agents.svelte";
import { settingsStore } from "./settings.svelte";
import { activityStore } from "./activity.svelte";
import { skillStore } from "./skills.svelte";
import { uiStore } from "./ui.svelte";
import { platformStore } from "./platform.svelte";
import { explorerStore } from "./explorer.svelte";
import { kitStore } from "./kits.svelte";
import { mcStore } from "./mission-control.svelte";
import { projectStore } from "./projects.svelte";
import { metricsStore } from "./metrics.svelte";

export {
  connectionStore,
  chatStore,
  sessionStore,
  agentsStore,
  settingsStore,
  activityStore,
  skillStore,
  uiStore,
  platformStore,
  explorerStore,
  kitStore,
  mcStore,
  projectStore,
  metricsStore,
};
export type { FileTypeCategory, ExplorerTab } from "./explorer.svelte";
export type { ActivityEntry } from "./activity.svelte";
export type { ActiveExecution, ExecutionLogEntry } from "./mission-control.svelte";
export type { PlanningPhase } from "./projects.svelte";

/** Post-auth initialization — Phase 2 core (agents, tasks/sessions, chat). */
export async function initializeStores(): Promise<void> {
  if (!connectionStore.isAuthenticated) return;

  await connectionStore.pingHealth();

  await Promise.all([agentsStore.loadAgents(), sessionStore.loadSessions()]);

  explorerStore.initialize();
  explorerStore.bindEvents();
}