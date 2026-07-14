import type { ClientCore } from '../core.js';
import type {
  AuditLog,
  AuthorizeRequest,
  AuthorizeResponse,
  CreateInviteResponse,
  ListParams,
  RequestConfig,
  Space,
  SpaceInvite,
  SpaceMember,
  SpaceRole,
  Team,
  TeamMember,
} from '../types.js';
import { buildQuery } from '../utils.js';

export class SpacesAPI {
  constructor(private readonly core: ClientCore) {}

  list(config?: RequestConfig): Promise<{ spaces: Space[]; total: number }> {
    return this.core.request('GET', '/spaces', undefined, { config });
  }

  create(
    input: { name: string; slug: string; description?: string },
    config?: RequestConfig,
  ): Promise<Space> {
    return this.core.request('POST', '/spaces', input, { config });
  }

  get(id: string, config?: RequestConfig): Promise<Space> {
    return this.core.request('GET', `/spaces/${id}`, undefined, { config });
  }

  update(
    id: string,
    input: Partial<Pick<Space, 'name' | 'slug' | 'description'>>,
    config?: RequestConfig,
  ): Promise<Space> {
    return this.core.request('PATCH', `/spaces/${id}`, input, { config });
  }

  delete(id: string, config?: RequestConfig): Promise<void> {
    return this.core.request('DELETE', `/spaces/${id}`, undefined, { config });
  }

  authorize(input: AuthorizeRequest, config?: RequestConfig): Promise<AuthorizeResponse> {
    return this.core.request('POST', '/authorize', input, { config });
  }

  members = {
    list: (spaceId?: string, params?: ListParams, config?: RequestConfig) => {
      const id = this.core.requireSpaceId(spaceId);
      return this.core.request<{ members: SpaceMember[]; total: number }>(
        'GET',
        buildQuery(`/spaces/${id}/members`, params),
        undefined,
        { config },
      );
    },
    add: (
      input: { user_id: string; role?: SpaceRole },
      spaceId?: string,
      config?: RequestConfig,
    ) => {
      const id = this.core.requireSpaceId(spaceId);
      return this.core.request<SpaceMember>('POST', `/spaces/${id}/members`, input, { config });
    },
    update: (memberId: string, role: SpaceRole, spaceId?: string, config?: RequestConfig) => {
      const id = this.core.requireSpaceId(spaceId);
      return this.core.request<SpaceMember>(
        'PATCH',
        `/spaces/${id}/members/${memberId}`,
        { role },
        { config },
      );
    },
    remove: (memberId: string, spaceId?: string, config?: RequestConfig) => {
      const id = this.core.requireSpaceId(spaceId);
      return this.core.request<void>(
        'DELETE',
        `/spaces/${id}/members/${memberId}`,
        undefined,
        { config },
      );
    },
  };

  invites = {
    list: (spaceId?: string, params?: ListParams, config?: RequestConfig) => {
      const id = this.core.requireSpaceId(spaceId);
      return this.core.request<{ invites: SpaceInvite[]; total: number }>(
        'GET',
        buildQuery(`/spaces/${id}/invites`, params),
        undefined,
        { config },
      );
    },
    create: (
      input: { email: string; role?: SpaceRole },
      spaceId?: string,
      config?: RequestConfig,
    ) => {
      const id = this.core.requireSpaceId(spaceId);
      return this.core.request<CreateInviteResponse>(
        'POST',
        `/spaces/${id}/invites`,
        input,
        { config },
      );
    },
    revoke: (inviteId: string, spaceId?: string, config?: RequestConfig) => {
      const id = this.core.requireSpaceId(spaceId);
      return this.core.request<void>(
        'DELETE',
        `/spaces/${id}/invites/${inviteId}`,
        undefined,
        { config },
      );
    },
  };

  teams = {
    list: (spaceId?: string, params?: ListParams, config?: RequestConfig) => {
      const id = this.core.requireSpaceId(spaceId);
      return this.core.request<{ teams: Team[]; total: number }>(
        'GET',
        buildQuery(`/spaces/${id}/teams`, params),
        undefined,
        { config },
      );
    },
    create: (
      input: { name: string; description?: string },
      spaceId?: string,
      config?: RequestConfig,
    ) => {
      const id = this.core.requireSpaceId(spaceId);
      return this.core.request<Team>('POST', `/spaces/${id}/teams`, input, { config });
    },
    get: (teamId: string, spaceId?: string, config?: RequestConfig) => {
      const id = this.core.requireSpaceId(spaceId);
      return this.core.request<Team>('GET', `/spaces/${id}/teams/${teamId}`, undefined, {
        config,
      });
    },
    update: (
      teamId: string,
      input: Partial<Pick<Team, 'name' | 'description'>>,
      spaceId?: string,
      config?: RequestConfig,
    ) => {
      const id = this.core.requireSpaceId(spaceId);
      return this.core.request<Team>('PATCH', `/spaces/${id}/teams/${teamId}`, input, {
        config,
      });
    },
    delete: (teamId: string, spaceId?: string, config?: RequestConfig) => {
      const id = this.core.requireSpaceId(spaceId);
      return this.core.request<void>('DELETE', `/spaces/${id}/teams/${teamId}`, undefined, {
        config,
      });
    },
    members: {
      list: (teamId: string, spaceId?: string, params?: ListParams, config?: RequestConfig) => {
        const id = this.core.requireSpaceId(spaceId);
        return this.core.request<{ members: TeamMember[]; total: number }>(
          'GET',
          buildQuery(`/spaces/${id}/teams/${teamId}/members`, params),
          undefined,
          { config },
        );
      },
      add: (
        teamId: string,
        input: { member_type: 'user' | 'agent'; member_id: string },
        spaceId?: string,
        config?: RequestConfig,
      ) => {
        const id = this.core.requireSpaceId(spaceId);
        return this.core.request<TeamMember>(
          'POST',
          `/spaces/${id}/teams/${teamId}/members`,
          input,
          { config },
        );
      },
      remove: (
        teamId: string,
        memberId: string,
        spaceId?: string,
        config?: RequestConfig,
      ) => {
        const id = this.core.requireSpaceId(spaceId);
        return this.core.request<void>(
          'DELETE',
          `/spaces/${id}/teams/${teamId}/members/${memberId}`,
          undefined,
          { config },
        );
      },
    },
  };

  auditLogs(spaceId?: string, params?: ListParams, config?: RequestConfig) {
    const id = this.core.requireSpaceId(spaceId);
    return this.core.request<{ logs: AuditLog[]; total: number }>(
      'GET',
      buildQuery(`/spaces/${id}/audit-logs`, params),
      undefined,
      { config },
    );
  }
}