import type { ClientCore } from '../core.js';
import type { RequestConfig } from '../types.js';

export interface IdentityFiles {
  identity_file: string;
  soul_file: string;
  style_file: string;
  instructions_file: string;
  user_file: string;
}

export interface IdentitySaveResponse {
  ok: boolean;
  updated: string[];
  agent_id?: string;
}

export type UpdateIdentityInput = Partial<IdentityFiles>;

export class IdentityAPI {
  constructor(private readonly core: ClientCore) {}

  get(agentId: string, config?: RequestConfig): Promise<IdentityFiles> {
    return this.core.request(
      'GET',
      `/agents/${encodeURIComponent(agentId)}/identity`,
      undefined,
      { config },
    );
  }

  update(
    agentId: string,
    input: UpdateIdentityInput,
    config?: RequestConfig,
  ): Promise<IdentitySaveResponse> {
    return this.core.request(
      'PUT',
      `/agents/${encodeURIComponent(agentId)}/identity`,
      input,
      { config },
    );
  }
}