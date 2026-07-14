import {
  createClient,
  isPocketAgentError,
  type AuthSession,
  type AuthUser,
  type PocketAgentClient,
  type RegisterResponse,
  type Space,
} from "@pocketagent/client";
import { GATE_URL } from "$lib/api/config";
import type { PocketPawClient } from "$lib/api/client";
import type { PocketPawWebSocket } from "$lib/api/websocket";
import { createLegacyFacade } from "$lib/api/legacy-facade";
import { logger } from "$lib/utils/logger";

export type ConnectionState = "connecting" | "connected" | "disconnected";

const PERSIST_KEY = "pocketagent";

interface PersistedState {
  v: 1;
  token?: string;
  user?: AuthUser;
  spaceId?: string;
}

function loadPersisted(): PersistedState | null {
  if (typeof localStorage === "undefined") return null;
  try {
    const raw = localStorage.getItem(PERSIST_KEY);
    if (!raw) return null;
    const parsed = JSON.parse(raw) as PersistedState;
    return parsed?.v === 1 ? parsed : null;
  } catch {
    return null;
  }
}

function savePersisted(data: Omit<PersistedState, "v">): void {
  if (typeof localStorage === "undefined") return;
  if (!data.token) {
    localStorage.removeItem(PERSIST_KEY);
    return;
  }
  localStorage.setItem(PERSIST_KEY, JSON.stringify({ v: 1, ...data }));
}

class ConnectionStore {
  token = $state<string | undefined>(undefined);
  user = $state<AuthUser | null>(null);
  spaceId = $state<string | undefined>(undefined);
  spaces = $state<Space[]>([]);

  ready = $state(false);
  loading = $state(false);
  error = $state<string | null>(null);

  status = $state<ConnectionState>("disconnected");
  backendUrl = $state(GATE_URL);
  isOffline = $state(false);

  isAuthenticated = $derived(Boolean(this.token && this.user));
  isConnected = $derived(this.status === "connected");

  private agentClient: PocketAgentClient | null = null;
  private legacyClient: PocketPawClient = createLegacyFacade();
  private onlineHandler: (() => void) | null = null;
  private offlineHandler: (() => void) | null = null;

  constructor() {
    if (typeof window === "undefined") return;
    this.isOffline = !navigator.onLine;
    this.onlineHandler = () => {
      this.isOffline = false;
      if (this.isAuthenticated) void this.pingHealth();
    };
    this.offlineHandler = () => {
      this.isOffline = true;
      this.status = "disconnected";
    };
    window.addEventListener("online", this.onlineHandler);
    window.addEventListener("offline", this.offlineHandler);
  }

  private buildClient(token?: string, spaceId?: string): PocketAgentClient {
    const client = createClient({
      baseUrl: this.backendUrl,
      token: token ?? this.token,
      spaceId: spaceId ?? this.spaceId,
      onSession: (session) => this.applySession(session),
    });
    client.auth.wireRefresh();
    return client;
  }

  private applySession(session: AuthSession): void {
    this.token = session.token;
    this.user = session.user;
    this.agentClient = this.buildClient(session.token, this.spaceId);
    savePersisted({ token: session.token, user: session.user, spaceId: this.spaceId });
  }

  private clearSession(): void {
    this.token = undefined;
    this.user = null;
    this.agentClient = null;
    savePersisted({});
  }

  async bootstrap(): Promise<boolean> {
    this.loading = true;
    this.error = null;
    this.backendUrl = GATE_URL;

    const persisted = loadPersisted();
    if (persisted?.token) {
      this.token = persisted.token;
      this.user = persisted.user ?? null;
      this.spaceId = persisted.spaceId;
      this.agentClient = this.buildClient(persisted.token, persisted.spaceId);
      try {
        const session = await this.agentClient.auth.refresh();
        this.applySession(session);
      } catch {
        this.clearSession();
      }
    }

    if (!this.isAuthenticated) {
      this.status = "disconnected";
      this.ready = true;
      this.loading = false;
      return false;
    }

    try {
      await this.loadSpaces();
      await this.pingHealth();
    } catch (err) {
      this.error = isPocketAgentError(err) ? err.message : String(err);
      this.status = "disconnected";
    }

    this.ready = true;
    this.loading = false;
    return true;
  }

  async login(email: string, password: string): Promise<void> {
    this.loading = true;
    this.error = null;
    try {
      const client = this.buildClient(undefined, this.spaceId);
      const session = await client.auth.login(email, password);
      this.applySession(session);
      await this.loadSpaces();
      await this.pingHealth();
    } catch (err) {
      this.error = isPocketAgentError(err) ? err.message : String(err);
      throw err;
    } finally {
      this.loading = false;
      this.ready = true;
    }
  }

  async register(email: string, password: string): Promise<RegisterResponse> {
    this.loading = true;
    this.error = null;
    try {
      const client = createClient({ baseUrl: this.backendUrl });
      return await client.auth.register(email, password);
    } catch (err) {
      this.error = isPocketAgentError(err) ? err.message : String(err);
      throw err;
    } finally {
      this.loading = false;
    }
  }

  logout(): void {
    this.clearSession();
    this.spaces = [];
    this.spaceId = undefined;
    this.status = "disconnected";
  }

  selectSpace(id: string): void {
    this.spaceId = id;
    if (this.token) {
      savePersisted({ token: this.token, user: this.user ?? undefined, spaceId: id });
    }
    if (this.agentClient) {
      this.agentClient.spaceId = id;
    }
  }

  async loadSpaces(): Promise<void> {
    const client = this.requireAgentClient();
    const res = await client.spaces.list();
    this.spaces = res.spaces;
    if (!this.spaceId && res.spaces.length === 1) {
      this.selectSpace(res.spaces[0].id);
    } else if (this.spaceId && !res.spaces.some((s) => s.id === this.spaceId)) {
      this.selectSpace(res.spaces[0]?.id ?? "");
    }
  }

  async pingHealth(): Promise<void> {
    this.status = "connecting";
    try {
      const healthUrl = this.backendUrl ? `${this.backendUrl}/health` : "/health";
      const res = await fetch(healthUrl, {
        headers: this.token ? { Authorization: `Bearer ${this.token}` } : undefined,
      });
      this.status = res.ok ? "connected" : "disconnected";
      if (!res.ok) {
        logger.warn(`[Connection] Health check failed: ${res.status}`);
      }
    } catch (err) {
      this.status = "disconnected";
      logger.warn("[Connection] Health check error:", err);
    }
  }

  /** PocketAgent SDK client — use for gate API (auth, spaces, agents, tasks). */
  requireAgentClient(): PocketAgentClient {
    if (!this.agentClient) {
      throw new Error("Not authenticated. Log in first.");
    }
    return this.agentClient;
  }

  getAgentClient(): PocketAgentClient {
    return this.requireAgentClient();
  }

  /**
   * Legacy PocketPaw client stub for unmigrated pages.
   * @deprecated Use getAgentClient() — Phase 2 will remove this.
   */
  getClient(): PocketPawClient {
    if (!this.isAuthenticated) {
      throw new Error("Not authenticated. Log in first.");
    }
    return this.legacyClient;
  }

  /** @deprecated Global WS removed — use per-task streams in Phase 2. */
  getWebSocket(): PocketPawWebSocket | null {
    return null;
  }

  disconnect(): void {
    this.status = "disconnected";
  }
}

export const connectionStore = new ConnectionStore();