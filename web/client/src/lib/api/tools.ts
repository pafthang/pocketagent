import type { ClientCore } from '../core.js';
import type { RequestConfig } from '../types.js';

export interface ToolInfo {
  name: string;
  description: string;
  source: 'builtin' | 'mcp' | string;
  server?: string;
}

export class ToolsAPI {
  constructor(private readonly core: ClientCore) {}

  list(config?: RequestConfig): Promise<{ tools: ToolInfo[] }> {
    return this.core.request('GET', '/tools', undefined, { config });
  }
}