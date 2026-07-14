import { applySession, createRequestContext, requestJSON, type RequestContext } from './http.js';
import type { AuthSession, PocketAgentClientOptions, RequestConfig } from './types.js';

export class ClientCore {
  readonly ctx: RequestContext;

  constructor(options: PocketAgentClientOptions) {
    this.ctx = createRequestContext(options);
  }

  get baseUrl(): string {
    return this.ctx.baseUrl;
  }

  get token(): string | undefined {
    return this.ctx.token;
  }

  get spaceId(): string | undefined {
    return this.ctx.spaceId;
  }

  setToken(token?: string): void {
    this.ctx.token = token;
  }

  setSpaceId(spaceId?: string): void {
    this.ctx.spaceId = spaceId;
  }

  applySession(session: AuthSession): void {
    applySession(this.ctx, session);
  }

  withToken(token?: string): ClientCore {
    return new ClientCore({
      baseUrl: this.ctx.baseUrl,
      token,
      spaceId: this.ctx.spaceId,
      fetch: this.ctx.fetchFn,
      timeoutMs: this.ctx.timeoutMs,
      onSession: this.ctx.onSession,
    });
  }

  withSpace(spaceId?: string): ClientCore {
    return new ClientCore({
      baseUrl: this.ctx.baseUrl,
      token: this.ctx.token,
      spaceId,
      fetch: this.ctx.fetchFn,
      timeoutMs: this.ctx.timeoutMs,
      onSession: this.ctx.onSession,
    });
  }

  requireSpaceId(spaceId?: string): string {
    const id = spaceId ?? this.ctx.spaceId;
    if (!id) {
      throw new Error('spaceId is required — set client.spaceId or pass spaceId explicitly');
    }
    return id;
  }

  request<T>(
    method: string,
    path: string,
    body?: unknown,
    opts?: {
      requireAuth?: boolean;
      requireSpace?: boolean;
      config?: RequestConfig;
    },
  ): Promise<T> {
    return requestJSON<T>(this.ctx, method, path, body, opts);
  }
}