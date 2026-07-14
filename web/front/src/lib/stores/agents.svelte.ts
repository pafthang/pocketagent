import type { Agent, CreateAgentInput } from "@pocketagent/client";
import { isPocketAgentError } from "@pocketagent/client";
import { toast } from "svelte-sonner";
import { connectionStore } from "./connection.svelte";
import { logger } from "$lib/utils/logger";

const SELECTED_KEY = "pocketagent_selected_agent";

class AgentsStore {
  agents = $state<Agent[]>([]);
  selectedAgentId = $state<string | null>(null);
  isLoading = $state(false);
  error = $state<string | null>(null);

  selectedAgent = $derived(
    this.agents.find((a) => a.id === this.selectedAgentId) ?? null,
  );

  constructor() {
    try {
      this.selectedAgentId = localStorage.getItem(SELECTED_KEY);
    } catch {
      // ignore
    }
  }

  selectAgent(id: string | null): void {
    this.selectedAgentId = id;
    try {
      if (id) localStorage.setItem(SELECTED_KEY, id);
      else localStorage.removeItem(SELECTED_KEY);
    } catch {
      // ignore
    }
  }

  async loadAgents(): Promise<void> {
    this.isLoading = true;
    this.error = null;
    try {
      const client = connectionStore.getAgentClient();
      const res = await client.agents.list({ per_page: 100 });
      this.agents = res.agents;

      if (this.selectedAgentId && !res.agents.some((a) => a.id === this.selectedAgentId)) {
        this.selectAgent(res.agents[0]?.id ?? null);
      } else if (!this.selectedAgentId && res.agents.length === 1) {
        this.selectAgent(res.agents[0].id);
      }
    } catch (err) {
      this.error = isPocketAgentError(err) ? err.message : String(err);
      logger.error("[AgentsStore] Failed to load agents:", err);
      toast.error("Failed to load agents");
    } finally {
      this.isLoading = false;
    }
  }

  async createAgent(input: CreateAgentInput): Promise<Agent | null> {
    try {
      const client = connectionStore.getAgentClient();
      const agent = await client.agents.create(input);
      this.agents = [agent, ...this.agents];
      if (!this.selectedAgentId) this.selectAgent(agent.id);
      toast.success(`Agent "${agent.name}" created`);
      return agent;
    } catch (err) {
      const msg = isPocketAgentError(err) ? err.message : String(err);
      toast.error(msg);
      return null;
    }
  }

  async updateAgent(id: string, input: Partial<CreateAgentInput>): Promise<Agent | null> {
    try {
      const client = connectionStore.getAgentClient();
      const agent = await client.agents.update(id, input);
      this.agents = this.agents.map((a) => (a.id === id ? agent : a));
      toast.success("Agent updated");
      return agent;
    } catch (err) {
      const msg = isPocketAgentError(err) ? err.message : String(err);
      toast.error(msg);
      return null;
    }
  }

  async deleteAgent(id: string): Promise<boolean> {
    try {
      const client = connectionStore.getAgentClient();
      await client.agents.delete(id);
      this.agents = this.agents.filter((a) => a.id !== id);
      if (this.selectedAgentId === id) {
        this.selectAgent(this.agents[0]?.id ?? null);
      }
      toast.success("Agent deleted");
      return true;
    } catch (err) {
      const msg = isPocketAgentError(err) ? err.message : String(err);
      toast.error(msg);
      return false;
    }
  }

  reset(): void {
    this.agents = [];
    this.error = null;
    this.isLoading = false;
  }
}

export const agentsStore = new AgentsStore();