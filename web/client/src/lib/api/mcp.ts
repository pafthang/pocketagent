import type { ClientCore } from '../core.js';
import type { ListParams, RequestConfig } from '../types.js';
import { buildQuery } from '../utils.js';

export interface MCPServer {
  id: string;
  space_id: string;
  name: string;
  transport: 'stdio' | 'http' | string;
  command?: string;
  args?: string[];
  url?: string;
  env?: Record<string, string>;
  enabled: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface MCPServerStatus {
  connected: boolean;
  connecting?: boolean;
  tool_count: number;
  error: string;
  transport: string;
  enabled: boolean;
}

export interface MCPProbeTool {
  name: string;
  description: string;
}

export interface MCPProbeResult {
  connected: boolean;
  error?: string;
  tools?: MCPProbeTool[];
}

export interface MCPPreset {
  id: string;
  name: string;
  description: string;
  icon: string;
  category: string;
  package: string;
  transport: string;
  url?: string;
  docs_url: string;
  needs_args: boolean;
  oauth: boolean;
  installed: boolean;
  env_keys: {
    key: string;
    label: string;
    required: boolean;
    placeholder: string;
    secret: boolean;
  }[];
}

export interface CreateMCPServerInput {
  name: string;
  transport: string;
  command?: string;
  args?: string[];
  url?: string;
  env?: Record<string, string>;
  enabled?: boolean;
}

export interface PatchMCPServerInput {
  name?: string;
  transport?: string;
  command?: string;
  args?: string[];
  url?: string;
  env?: Record<string, string>;
  enabled?: boolean;
}

export class McpAPI {
  constructor(private readonly core: ClientCore) {}

  list(params: ListParams = {}, config?: RequestConfig): Promise<{ servers: MCPServer[]; total: number }> {
    return this.core.request('GET', buildQuery('/mcp/servers', params), undefined, { config });
  }

  get(id: string, config?: RequestConfig): Promise<MCPServer> {
    return this.core.request('GET', `/mcp/servers/${encodeURIComponent(id)}`, undefined, { config });
  }

  create(input: CreateMCPServerInput, config?: RequestConfig): Promise<MCPServer> {
    return this.core.request('POST', '/mcp/servers', input, { config });
  }

  patch(id: string, input: PatchMCPServerInput, config?: RequestConfig): Promise<MCPServer> {
    return this.core.request('PATCH', `/mcp/servers/${encodeURIComponent(id)}`, input, { config });
  }

  delete(id: string, config?: RequestConfig): Promise<void> {
    return this.core.request('DELETE', `/mcp/servers/${encodeURIComponent(id)}`, undefined, { config });
  }

  test(id: string, config?: RequestConfig): Promise<MCPProbeResult> {
    return this.core.request('POST', `/mcp/servers/${encodeURIComponent(id)}/test`, undefined, { config });
  }

  status(config?: RequestConfig): Promise<Record<string, MCPServerStatus>> {
    return this.core.request('GET', '/mcp/status', undefined, { config });
  }

  presets(config?: RequestConfig): Promise<MCPPreset[]> {
    return this.core.request('GET', '/mcp/presets', undefined, { config });
  }

  installPreset(
    input: { preset_id: string; env?: Record<string, string>; extra_args?: string[] },
    config?: RequestConfig,
  ): Promise<{ status: string; connected?: boolean; error?: string }> {
    return this.core.request('POST', '/mcp/presets/install', input, { config });
  }
}