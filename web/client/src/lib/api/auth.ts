import type { ClientCore } from '../core.js';
import type {
  AuthSession,
  RegisterResponse,
  RequestConfig,
} from '../types.js';

export class AuthAPI {
  constructor(private readonly core: ClientCore) {}

  async register(
    email: string,
    password: string,
    config?: RequestConfig,
  ): Promise<RegisterResponse> {
    return this.core.request<RegisterResponse>(
      'POST',
      '/auth/register',
      { email, password },
      { requireAuth: false, requireSpace: false, config },
    );
  }

  async login(email: string, password: string, config?: RequestConfig): Promise<AuthSession> {
    const session = await this.core.request<AuthSession>(
      'POST',
      '/auth/login',
      { email, password },
      { requireAuth: false, requireSpace: false, config },
    );
    this.core.applySession(session);
    return session;
  }

  async refresh(config?: RequestConfig): Promise<AuthSession> {
    const session = await this.refreshRaw(config);
    this.core.applySession(session);
    return session;
  }

  /** Internal refresh without auto-retry recursion. */
  async refreshRaw(config?: RequestConfig): Promise<AuthSession> {
    return this.core.request<AuthSession>(
      'POST',
      '/auth/refresh',
      undefined,
      { requireAuth: true, requireSpace: false, config: { ...config, skipAuthRefresh: true } },
    );
  }

  async verifyEmail(token: string, config?: RequestConfig): Promise<void> {
    await this.core.request<void>(
      'POST',
      '/auth/verify-email',
      { token },
      { requireAuth: false, requireSpace: false, config },
    );
  }

  async requestVerification(config?: RequestConfig): Promise<void> {
    await this.core.request<void>(
      'POST',
      '/auth/request-verification',
      undefined,
      { requireAuth: true, requireSpace: false, config },
    );
  }

  async previewInvite(token: string, config?: RequestConfig) {
    return this.core.request<import('../types.js').InvitePreview>(
      'GET',
      `/invites/${encodeURIComponent(token)}`,
      undefined,
      { requireAuth: false, requireSpace: false, config },
    );
  }

  async acceptInvite(
    token: string,
    password?: string,
    config?: RequestConfig,
  ): Promise<{ space_id: string; member_id: string }> {
    return this.core.request(
      'POST',
      '/invites/accept',
      { token, password },
      { requireAuth: false, requireSpace: false, config },
    );
  }

  wireRefresh(): void {
    this.core.ctx.refreshSession = async () => {
      if (!this.core.ctx.token) return undefined;
      try {
        const session = await this.refreshRaw();
        this.core.applySession(session);
        return session;
      } catch {
        return undefined;
      }
    };
  }
}