import type { ClientCore } from '../core.js';
import type { ListParams, RequestConfig, Task } from '../types.js';
import { buildQuery } from '../utils.js';

export interface Skill {
  id: string;
  space_id: string;
  name: string;
  description: string;
  prompt: string;
  category?: string;
  tools?: string[];
  argument_hint?: string;
  catalog_id?: string;
  created_at?: string;
  updated_at?: string;
}

export interface SkillCatalogEntry {
  id: string;
  name: string;
  description: string;
  category: string;
  argument_hint: string;
  prompt?: string;
  tools?: string[];
  installed: boolean;
}

export interface CreateSkillInput {
  name: string;
  description?: string;
  prompt: string;
  category?: string;
  tools?: string[];
  argument_hint?: string;
  catalog_id?: string;
}

export interface PatchSkillInput {
  name?: string;
  description?: string;
  prompt?: string;
  category?: string;
  tools?: string[];
  argument_hint?: string;
}

export interface RunSkillInput {
  agent_id?: string;
  input?: string;
}

export interface RunSkillResult extends Task {
  tools?: string[];
  skill_id?: string;
}

export interface SkillListResult {
  skills: Skill[];
  total: number;
  page: number;
  per_page: number;
}

export class SkillsAPI {
  constructor(private readonly core: ClientCore) {}

  list(params: ListParams = {}, config?: RequestConfig): Promise<SkillListResult> {
    return this.core.request('GET', buildQuery('/skills', params), undefined, { config });
  }

  get(id: string, config?: RequestConfig): Promise<Skill> {
    return this.core.request('GET', `/skills/${encodeURIComponent(id)}`, undefined, { config });
  }

  create(input: CreateSkillInput, config?: RequestConfig): Promise<Skill> {
    return this.core.request('POST', '/skills', input, { config });
  }

  patch(id: string, input: PatchSkillInput, config?: RequestConfig): Promise<Skill> {
    return this.core.request('PATCH', `/skills/${encodeURIComponent(id)}`, input, { config });
  }

  delete(id: string, config?: RequestConfig): Promise<void> {
    return this.core.request('DELETE', `/skills/${encodeURIComponent(id)}`, undefined, { config });
  }

  catalog(config?: RequestConfig): Promise<SkillCatalogEntry[]> {
    return this.core.request('GET', '/skills/catalog', undefined, { config });
  }

  run(id: string, input: RunSkillInput = {}, config?: RequestConfig): Promise<RunSkillResult> {
    return this.core.request('POST', `/skills/${encodeURIComponent(id)}/run`, input, { config });
  }
}