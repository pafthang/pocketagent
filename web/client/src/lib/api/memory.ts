import type { ClientCore } from '../core.js';
import type { ListParams, RequestConfig } from '../types.js';
import { buildQuery } from '../utils.js';

export interface MemoryDocument {
  id: string;
  content: string;
  metadata?: Record<string, string>;
  created_at?: string;
  tags?: string[];
}

export interface MemoryListResult {
  documents: MemoryDocument[];
  total: number;
  page: number;
  per_page: number;
}

export interface MemorySearchResult {
  id: string;
  content: string;
  similarity: number;
}

export interface MemoryStats {
  backend: string;
  total_memories: number;
  document_count: number;
  content_bytes: number;
  collection: string;
  memories_by_type?: Record<string, number>;
}

export interface IngestMemoryInput {
  content: string;
  id?: string;
  tags?: string[];
  metadata?: Record<string, string>;
}

export interface SearchMemoryInput {
  query: string;
  limit?: number;
  min_similarity?: number;
}

export class MemoryAPI {
  constructor(private readonly core: ClientCore) {}

  list(params: ListParams = {}, config?: RequestConfig): Promise<MemoryListResult> {
    return this.core.request('GET', buildQuery('/memory', params), undefined, { config });
  }

  get(id: string, config?: RequestConfig): Promise<MemoryDocument> {
    return this.core.request('GET', `/memory/${encodeURIComponent(id)}`, undefined, { config });
  }

  ingest(input: IngestMemoryInput, config?: RequestConfig): Promise<{ id: string; status: string; content: string }> {
    return this.core.request('POST', '/memory', input, { config });
  }

  search(input: SearchMemoryInput, config?: RequestConfig): Promise<{ results: MemorySearchResult[] }> {
    return this.core.request('POST', '/memory/search', input, { config });
  }

  delete(id: string, config?: RequestConfig): Promise<{ status: string; id: string; removed: number }> {
    return this.core.request('DELETE', `/memory/${encodeURIComponent(id)}`, undefined, { config });
  }

  stats(config?: RequestConfig): Promise<MemoryStats> {
    return this.core.request('GET', '/memory/stats', undefined, { config });
  }
}