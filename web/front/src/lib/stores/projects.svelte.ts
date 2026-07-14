import type { MCProject, MCProjectProgress } from "$lib/types/pawkit";
import { connectProjectWebSocket, type PocketAgentClientOptions } from "@pocketagent/client";
import { connectionStore } from "./connection.svelte";

export type PlanningPhase = "goal_analysis" | "research" | "prd" | "tasks" | "team";

const PLAN_POLL_MS = 2000;

class ProjectStore {
  projects = $state<MCProject[]>([]);
  isLoading = $state(false);

  planningProjectId = $state<string | null>(null);
  planningPhase = $state<PlanningPhase | null>(null);
  planningMessage = $state("");
  planningDone = $state(false);
  planningError = $state<string | null>(null);

  private pollTimer: ReturnType<typeof setInterval> | null = null;
  private projectWs: WebSocket | null = null;
  private wsActive = false;

  private clientOpts(): PocketAgentClientOptions {
    return {
      baseUrl: connectionStore.backendUrl,
      token: connectionStore.token,
      spaceId: connectionStore.spaceId,
    };
  }

  async loadProjects(status?: string): Promise<void> {
    this.isLoading = true;
    try {
      const client = connectionStore.getAgentClient();
      const res = await client.projects.listMC(status ? { status } : {});
      this.projects = res.projects as unknown as MCProject[];
    } catch (err) {
      console.error("[ProjectStore] Failed to load projects:", err);
    } finally {
      this.isLoading = false;
    }
  }

  async getProjectDetail(
    projectId: string,
  ): Promise<{ project: MCProject; tasks: Record<string, unknown>[]; progress: MCProjectProgress } | null> {
    try {
      const client = connectionStore.getAgentClient();
      const res = await client.projects.getDetail(projectId);
      return {
        project: res.project as unknown as MCProject,
        tasks: res.tasks,
        progress: res.progress,
      };
    } catch (err) {
      console.error("[ProjectStore] Failed to get project:", err);
      return null;
    }
  }

  async parseGoal(goal: string): Promise<Record<string, unknown> | null> {
    try {
      const client = connectionStore.getAgentClient();
      const res = await client.projects.parseGoal(goal);
      return res.goal_analysis ?? null;
    } catch (err) {
      console.error("[ProjectStore] Failed to parse goal:", err);
      return null;
    }
  }

  async startProject(goal: string, description?: string): Promise<string | null> {
    this.planningPhase = null;
    this.planningMessage = "";
    this.planningDone = false;
    this.planningError = null;
    try {
      const client = connectionStore.getAgentClient();
      const res = await client.projects.start(goal, description ? { title: description } : {});
      const projectId = res.project_id ?? res.project?.id ?? null;
      this.planningProjectId = projectId;
      if (projectId) {
        this.startPlanningTrack(projectId);
      }
      return projectId;
    } catch (err) {
      console.error("[ProjectStore] Failed to start project:", err);
      this.planningError = String(err);
      return null;
    }
  }

  async getPlan(projectId: string) {
    try {
      const client = connectionStore.getAgentClient();
      return await client.projects.getPlan(projectId);
    } catch (err) {
      console.error("[ProjectStore] Failed to get plan:", err);
      return null;
    }
  }

  private startPlanningTrack(projectId: string) {
    this.stopPlanningTrack();
    this.startPlanningWs(projectId);
    this.startPlanningPoll(projectId);
  }

  private startPlanningWs(projectId: string) {
    this.wsActive = false;
    try {
      const ws = connectProjectWebSocket(this.clientOpts(), projectId, {
        onPhase: (e) => {
          if (this.planningProjectId !== projectId) return;
          this.wsActive = true;
          this.stopPlanningPoll();
          this.planningPhase = e.phase as PlanningPhase;
          this.planningMessage = e.message;
        },
        onComplete: (e) => {
          if (this.planningProjectId !== projectId) return;
          this.planningDone = true;
          if (e.error) this.planningError = e.error;
          this.stopPlanningTrack();
          void this.loadProjects();
        },
        onError: () => {
          if (!this.wsActive && !this.pollTimer) {
            this.startPlanningPoll(projectId);
          }
        },
        onClose: () => {
          this.projectWs = null;
          if (!this.planningDone && this.planningProjectId === projectId && !this.pollTimer) {
            this.startPlanningPoll(projectId);
          }
        },
      });
      this.projectWs = ws;
      ws.onopen = () => {
        this.wsActive = true;
      };
    } catch (err) {
      console.warn("[ProjectStore] Project WS unavailable, using poll fallback:", err);
    }
  }

  private startPlanningPoll(projectId: string) {
    if (this.pollTimer) return;
    const poll = async () => {
      const plan = await this.getPlan(projectId);
      if (!plan || this.planningProjectId !== projectId) return;

      const phase = plan.current_phase as PlanningPhase | undefined;
      if (phase) this.planningPhase = phase;
      if (typeof plan.phase_message === "string") {
        this.planningMessage = plan.phase_message;
      }
      if (plan.planning_error) {
        this.planningError = String(plan.planning_error);
      }
      if (plan.planning_done || plan.status === "awaiting_approval") {
        this.planningDone = true;
        this.stopPlanningTrack();
        await this.loadProjects();
      }
      if (plan.status === "failed") {
        this.planningError = this.planningError ?? "Planning failed";
        this.stopPlanningTrack();
      }
    };
    void poll();
    this.pollTimer = setInterval(() => void poll(), PLAN_POLL_MS);
  }

  private stopPlanningPoll() {
    if (this.pollTimer) {
      clearInterval(this.pollTimer);
      this.pollTimer = null;
    }
  }

  private stopProjectWs() {
    if (this.projectWs) {
      this.projectWs.close();
      this.projectWs = null;
    }
    this.wsActive = false;
  }

  private stopPlanningTrack() {
    this.stopPlanningPoll();
    this.stopProjectWs();
  }

  async approveProject(projectId: string): Promise<boolean> {
    try {
      const client = connectionStore.getAgentClient();
      await client.projects.approve(projectId);
      await this.loadProjects();
      return true;
    } catch (err) {
      console.error("[ProjectStore] Failed to approve:", err);
      return false;
    }
  }

  async pauseProject(projectId: string): Promise<boolean> {
    try {
      const client = connectionStore.getAgentClient();
      await client.projects.pause(projectId);
      await this.loadProjects();
      return true;
    } catch (err) {
      console.error("[ProjectStore] Failed to pause:", err);
      return false;
    }
  }

  async resumeProject(projectId: string): Promise<boolean> {
    try {
      const client = connectionStore.getAgentClient();
      await client.projects.resume(projectId);
      await this.loadProjects();
      return true;
    } catch (err) {
      console.error("[ProjectStore] Failed to resume:", err);
      return false;
    }
  }

  async cancelProject(projectId: string): Promise<boolean> {
    try {
      const client = connectionStore.getAgentClient();
      await client.projects.cancel(projectId);
      await this.loadProjects();
      return true;
    } catch (err) {
      console.error("[ProjectStore] Failed to cancel:", err);
      return false;
    }
  }

  async deleteProject(projectId: string): Promise<boolean> {
    try {
      const client = connectionStore.getAgentClient();
      await client.projects.delete(projectId);
      this.projects = this.projects.filter((p) => p.id !== projectId);
      return true;
    } catch (err) {
      console.error("[ProjectStore] Failed to delete:", err);
      return false;
    }
  }

  async skipTask(projectId: string, taskId: string): Promise<boolean> {
    try {
      const client = connectionStore.getAgentClient();
      await client.projects.skipTask(projectId, taskId);
      return true;
    } catch (err) {
      console.error("[ProjectStore] Failed to skip task:", err);
      return false;
    }
  }

  async retryTask(projectId: string, taskId: string): Promise<boolean> {
    try {
      const client = connectionStore.getAgentClient();
      await client.projects.retryTask(projectId, taskId);
      return true;
    } catch (err) {
      console.error("[ProjectStore] Failed to retry task:", err);
      return false;
    }
  }

  /** @deprecated Global WS removed — per-project stream used in startPlanningTrack. */
  bindEvents(): void {}

  unbindEvents(): void {}

  initialize(): void {}

  dispose(): void {
    this.stopPlanningTrack();
  }
}

export const projectStore = new ProjectStore();