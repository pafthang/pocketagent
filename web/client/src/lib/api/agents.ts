import type { ClientCore } from '../core.js';
import type { Agent, CreateAgentInput, ListParams, RequestConfig } from '../types.js';
import { buildQuery } from '../utils.js';

export class AgentsAPI {
  constructor(private readonly core: ClientCore) {}

  list(params: ListParams = {}, config?: RequestConfig): Promise<{ agents: Agent[]; total: number }> {
    return this.core.request('GET', buildQuery('/agents', params), undefined, { config });
  }

  get(id: string, config?: RequestConfig): Promise<Agent> {
    return this.core.request('GET', `/agents/${id}`, undefined, { config });
  }

  create(input: CreateAgentInput, config?: RequestConfig): Promise<Agent> {
    return this.core.request('POST', '/agents', input, { config });
  }

  update(id: string, input: Partial<CreateAgentInput>, config?: RequestConfig): Promise<Agent> {
    return this.core.request('PUT', `/agents/${id}`, input, { config });
  }

  delete(id: string, config?: RequestConfig): Promise<void> {
    return this.core.request('DELETE', `/agents/${id}`, undefined, { config });
  }
}