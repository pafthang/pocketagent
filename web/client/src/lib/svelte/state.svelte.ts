import { createClient, type PocketAgentClient } from '../client.js';
import { isPocketAgentError, isUnauthorized } from '../errors.js';
import type {
  AuthSession,
  AuthUser,
  CreateAgentInput,
  CreateTaskInput,
  PocketAgentClientOptions,
} from '../types.js';
import { createResourceState } from './resource.svelte.js';

export interface PocketAgentStateOptions extends PocketAgentClientOptions {
  /** localStorage key for token, user, spaceId */
  persistKey?: string;
  /** Try refresh on init when persisted token exists */
  restoreSession?: boolean;
}

interface PersistedState {
  v: 1;
  token?: string;
  user?: AuthUser;
  spaceId?: string;
}

function loadPersisted(key?: string): PersistedState | null {
  if (!key || typeof localStorage === 'undefined') return null;
  try {
    const raw = localStorage.getItem(key);
    if (!raw) return null;
    const parsed = JSON.parse(raw) as PersistedState;
    return parsed?.v === 1 ? parsed : null;
  } catch {
    return null;
  }
}

function savePersisted(key: string | undefined, data: Omit<PersistedState, 'v'>): void {
  if (!key || typeof localStorage === 'undefined') return;
  if (!data.token) {
    localStorage.removeItem(key);
    return;
  }
  localStorage.setItem(key, JSON.stringify({ v: 1, ...data }));
}

export function createPocketAgentState(opts: PocketAgentStateOptions) {
  const persisted = loadPersisted(opts.persistKey);

  let token = $state<string | undefined>(opts.token ?? persisted?.token);
  let user = $state<AuthUser | null>(persisted?.user ?? null);
  let spaceId = $state<string | undefined>(opts.spaceId ?? persisted?.spaceId);
  let loading = $state(false);
  let error = $state<string | null>(null);
  let ready = $state(!opts.restoreSession);

  const client = $derived(
    createClient({
      baseUrl: opts.baseUrl,
      token,
      spaceId,
      fetch: opts.fetch,
      timeoutMs: opts.timeoutMs,
      onSession: (session) => setSession(session),
    }),
  );

  const isAuthenticated = $derived(Boolean(token && user));

  function persist(): void {
    savePersisted(opts.persistKey, { token, user: user ?? undefined, spaceId });
  }

  function setSession(session: AuthSession): void {
    token = session.token;
    user = session.user;
    persist();
  }

  function clearSession(): void {
    token = undefined;
    user = null;
    persist();
  }

  function handleError(err: unknown): void {
    error = isPocketAgentError(err) ? err.message : String(err);
    if (isUnauthorized(err)) clearSession();
  }

  async function run<T>(fn: () => Promise<T>): Promise<T | undefined> {
    loading = true;
    error = null;
    try {
      return await fn();
    } catch (err) {
      handleError(err);
      return undefined;
    } finally {
      loading = false;
    }
  }

  async function init(): Promise<void> {
    if (!opts.restoreSession || !token) {
      ready = true;
      return;
    }
    loading = true;
    try {
      const session = await client.auth.refresh();
      setSession(session);
    } catch {
      clearSession();
    } finally {
      loading = false;
      ready = true;
    }
  }

  if (opts.restoreSession) {
    init();
  }

  const agents = createResourceState({
    fetch: async (page) => {
      const res = await run(() => client.agents.list({ page }));
      return res ? { items: res.agents, total: res.total } : undefined;
    },
    onUnauthorized: clearSession,
  });

  const tasks = createResourceState({
    fetch: async (page) => {
      const res = await run(() => client.tasks.list({ page }));
      return res ? { items: res.tasks, total: res.total } : undefined;
    },
    onUnauthorized: clearSession,
  });

  const spaces = createResourceState({
    fetch: async () => {
      const res = await run(() => client.spaces.list());
      if (res && !spaceId && res.spaces.length === 1) {
        spaceId = res.spaces[0].id;
        persist();
      }
      return res ? { items: res.spaces, total: res.total } : undefined;
    },
    onUnauthorized: clearSession,
  });

  const schedules = createResourceState({
    fetch: async (page) => {
      const res = await run(() => client.schedules.list({ page }));
      return res ? { items: res.schedules, total: res.total } : undefined;
    },
    onUnauthorized: clearSession,
  });

  return {
    get token() {
      return token;
    },
    get user() {
      return user;
    },
    get spaceId() {
      return spaceId;
    },
    set spaceId(id: string | undefined) {
      spaceId = id;
      persist();
    },
    get loading() {
      return loading;
    },
    get error() {
      return error;
    },
    get ready() {
      return ready;
    },
    get client(): PocketAgentClient {
      return client;
    },
    get isAuthenticated() {
      return isAuthenticated;
    },
    get agents() {
      return agents;
    },
    get tasks() {
      return tasks;
    },
    get spaces() {
      return spaces;
    },
    get schedules() {
      return schedules;
    },

    clearError() {
      error = null;
    },

    async login(email: string, password: string) {
      return run(async () => {
        const session = await client.auth.login(email, password);
        setSession(session);
        return session;
      });
    },

    async register(email: string, password: string) {
      return run(() => client.auth.register(email, password));
    },

    async refresh() {
      return run(async () => {
        const session = await client.auth.refresh();
        setSession(session);
        return session;
      });
    },

    logout() {
      clearSession();
    },

    selectSpace(id: string) {
      spaceId = id;
      persist();
    },

    async createAgent(input: CreateAgentInput) {
      return run(() => client.agents.create(input));
    },

    async createTask(input: CreateTaskInput) {
      return run(() => client.tasks.create(input));
    },

    async cancelTask(id: string) {
      return run(() => client.tasks.cancel(id));
    },

    async getTask(id: string, includeSubtasks = false) {
      if (includeSubtasks) {
        return run(() => client.tasks.getWithSubtasks(id));
      }
      return run(() => client.tasks.get(id));
    },

    init,
  };
}

export type PocketAgentState = ReturnType<typeof createPocketAgentState>;