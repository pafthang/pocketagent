import { errorFromResponse } from './errors.js';
import type { AuthSession, PocketAgentClientOptions, RequestConfig } from './types.js';

export const HEADER_SPACE_ID = 'X-Space-Id';
const DEFAULT_TIMEOUT_MS = 30_000;

export interface RequestContext {
  baseUrl: string;
  token?: string;
  spaceId?: string;
  fetchFn: typeof fetch;
  timeoutMs: number;
  onSession?: (session: AuthSession) => void;
  refreshSession?: () => Promise<AuthSession | undefined>;
}

export function createRequestContext(opts: PocketAgentClientOptions): RequestContext {
  const baseUrl = opts.baseUrl.replace(/\/$/, '');
  return {
    baseUrl,
    token: opts.token,
    spaceId: opts.spaceId,
    fetchFn: opts.fetch ?? globalThis.fetch.bind(globalThis),
    timeoutMs: opts.timeoutMs ?? DEFAULT_TIMEOUT_MS,
    onSession: opts.onSession,
  };
}

export function applySession(ctx: RequestContext, session: AuthSession): void {
  ctx.token = session.token;
  ctx.onSession?.(session);
}

export function buildHeaders(
  ctx: RequestContext,
  init?: HeadersInit,
  opts?: { requireAuth?: boolean; requireSpace?: boolean },
): Headers {
  const headers = new Headers(init);
  if (!headers.has('Accept')) {
    headers.set('Accept', 'application/json');
  }

  if (opts?.requireAuth !== false && ctx.token) {
    headers.set('Authorization', `Bearer ${ctx.token}`);
  }

  if (opts?.requireSpace !== false && ctx.spaceId) {
    headers.set(HEADER_SPACE_ID, ctx.spaceId);
  }

  return headers;
}

export async function requestJSON<T>(
  ctx: RequestContext,
  method: string,
  path: string,
  body?: unknown,
  opts?: {
    requireAuth?: boolean;
    requireSpace?: boolean;
    config?: RequestConfig;
  },
): Promise<T> {
  return requestJSONWithRetry(ctx, method, path, body, opts, false);
}

async function requestJSONWithRetry<T>(
  ctx: RequestContext,
  method: string,
  path: string,
  body: unknown,
  opts: { requireAuth?: boolean; requireSpace?: boolean; config?: RequestConfig } | undefined,
  retried: boolean,
): Promise<T> {
  try {
    return await executeJSON<T>(ctx, method, path, body, opts);
  } catch (err) {
    const shouldRefresh =
      !retried &&
      !opts?.config?.skipAuthRefresh &&
      opts?.requireAuth !== false &&
      ctx.token &&
      ctx.refreshSession &&
      typeof err === 'object' &&
      err !== null &&
      'status' in err &&
      (err as { status: number }).status === 401;

    if (!shouldRefresh) throw err;

    const session = await ctx.refreshSession!();
    if (!session) throw err;
    return requestJSONWithRetry(ctx, method, path, body, opts, true);
  }
}

async function executeJSON<T>(
  ctx: RequestContext,
  method: string,
  path: string,
  body: unknown,
  opts?: { requireAuth?: boolean; requireSpace?: boolean; config?: RequestConfig },
): Promise<T> {
  const res = await requestRaw(ctx, method, path, body, opts);
  if (res.status === 204) {
    return undefined as T;
  }
  return (await res.json()) as T;
}

export async function requestRaw(
  ctx: RequestContext,
  method: string,
  path: string,
  body?: unknown,
  opts?: {
    requireAuth?: boolean;
    requireSpace?: boolean;
    config?: RequestConfig;
    headers?: HeadersInit;
  },
): Promise<Response> {
  const headers = buildHeaders(
    ctx,
    body !== undefined
      ? { 'Content-Type': 'application/json', ...(opts?.headers as Record<string, string>) }
      : opts?.headers,
    opts,
  );

  const timeoutMs = opts?.config?.timeoutMs ?? ctx.timeoutMs;
  const signals: AbortSignal[] = [];
  if (opts?.config?.signal) signals.push(opts.config.signal);

  const timeoutController = new AbortController();
  const timer = setTimeout(() => timeoutController.abort(new Error('request timeout')), timeoutMs);
  signals.push(timeoutController.signal);

  const controller = new AbortController();
  for (const signal of signals) {
    if (signal.aborted) {
      controller.abort(signal.reason);
      break;
    }
    signal.addEventListener('abort', () => controller.abort(signal.reason), { once: true });
  }

  try {
    const res = await ctx.fetchFn(`${ctx.baseUrl}${path}`, {
      method,
      headers,
      body: body !== undefined ? JSON.stringify(body) : undefined,
      signal: controller.signal,
    });

    if (!res.ok) {
      throw await errorFromResponse(res);
    }
    return res;
  } finally {
    clearTimeout(timer);
  }
}

export function withContext(
  ctx: RequestContext,
  patch: Partial<Pick<RequestContext, 'token' | 'spaceId'>>,
): RequestContext {
  return { ...ctx, ...patch };
}