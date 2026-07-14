import type { ClientCore } from '../core.js';
import { buildHeaders } from '../http.js';
import { errorFromResponse } from '../errors.js';
import type { RequestConfig } from '../types.js';
import { buildQuery } from '../utils.js';

export interface StoredFile {
  id: string;
  space_id: string;
  project_id?: string;
  parent_id?: string;
  name: string;
  virtual_path: string;
  is_dir: boolean;
  mime_type?: string;
  size?: number;
  storage_key?: string;
  checksum?: string;
  memo_ingested?: boolean;
  uploaded_by?: string;
  tags?: string[];
  created_at?: string;
  updated_at?: string;
}

export interface BrowseFileEntry {
  id?: string;
  name: string;
  path: string;
  is_dir: boolean;
  size?: number;
  mime_type?: string;
  project_id?: string;
  modified_at?: string;
}

export interface BrowseFilesResult {
  path: string;
  project_id?: string;
  files: BrowseFileEntry[];
}

export interface RecentFileEntry {
  path: string;
  name: string;
  is_dir: boolean;
  extension: string;
  timestamp: number;
  tool: string;
}

export interface FileContentPreview {
  id: string;
  path: string;
  content: string;
}

export interface IngestFileInput {
  force?: boolean;
  tags?: string[];
}

export interface IngestFileResult {
  status: string;
  file_id: string;
  memo_id: string;
  path: string;
  content_bytes: number;
  replaced: boolean;
}

export interface UploadFileInput {
  file: Blob;
  filename: string;
  path?: string;
  projectId?: string;
}

export interface CreateFolderInput {
  name: string;
  path?: string;
  projectId?: string;
}

export interface BrowseFilesParams {
  path?: string;
  projectId?: string;
}

/** Build a virtual path under /projects/{id}/…. */
export function projectFilesPath(projectId: string, subpath = ''): string {
  const root = `/projects/${projectId}`;
  if (!subpath || subpath === '/' || subpath === '.') return root;
  const rel = subpath.startsWith('/') ? subpath.slice(1) : subpath;
  return `${root}/${rel}`.replace(/\/+/g, '/');
}

export class FilesAPI {
  constructor(private readonly core: ClientCore) {}

  browse(params: BrowseFilesParams = {}, config?: RequestConfig): Promise<BrowseFilesResult> {
    return this.core.request(
      'GET',
      buildQuery('/files/browse', {
        path: params.path,
        project_id: params.projectId,
      }),
      undefined,
      { config },
    );
  }

  recent(
    limit = 20,
    params: { projectId?: string } = {},
    config?: RequestConfig,
  ): Promise<{ files: RecentFileEntry[]; project_id?: string }> {
    return this.core.request(
      'GET',
      buildQuery('/files/recent', { limit, project_id: params.projectId }),
      undefined,
      { config },
    );
  }

  browseProject(
    projectId: string,
    subpath = '',
    config?: RequestConfig,
  ): Promise<BrowseFilesResult> {
    return this.core.request(
      'GET',
      buildQuery(`/projects/${encodeURIComponent(projectId)}/files/browse`, {
        path: subpath || undefined,
      }),
      undefined,
      { config },
    );
  }

  recentProject(
    projectId: string,
    limit = 20,
    config?: RequestConfig,
  ): Promise<{ files: RecentFileEntry[]; project_id?: string }> {
    return this.core.request(
      'GET',
      buildQuery(`/projects/${encodeURIComponent(projectId)}/files/recent`, { limit }),
      undefined,
      { config },
    );
  }

  get(id: string, config?: RequestConfig): Promise<StoredFile> {
    return this.core.request('GET', `/files/${encodeURIComponent(id)}`, undefined, { config });
  }

  content(id: string, config?: RequestConfig): Promise<FileContentPreview> {
    return this.core.request(
      'GET',
      `/files/${encodeURIComponent(id)}/content`,
      undefined,
      { config },
    );
  }

  async download(id: string, config?: RequestConfig): Promise<Blob> {
    const headers = buildHeaders(this.core.ctx);
    const res = await this.core.ctx.fetchFn(
      `${this.core.baseUrl}/files/${encodeURIComponent(id)}/download`,
      {
        method: 'GET',
        headers,
        signal: config?.signal,
      },
    );
    if (!res.ok) {
      throw await errorFromResponse(res);
    }
    return res.blob();
  }

  async uploadProject(
    projectId: string,
    input: Omit<UploadFileInput, 'projectId'>,
    config?: RequestConfig,
  ): Promise<StoredFile> {
    const form = new FormData();
    form.append('file', input.file, input.filename);
    if (input.path) form.append('path', input.path);

    const headers = buildHeaders(this.core.ctx);
    headers.delete('Content-Type');

    const res = await this.core.ctx.fetchFn(
      `${this.core.baseUrl}/projects/${encodeURIComponent(projectId)}/files/upload`,
      {
        method: 'POST',
        headers,
        body: form,
        signal: config?.signal,
      },
    );
    if (!res.ok) {
      throw await errorFromResponse(res);
    }
    return (await res.json()) as StoredFile;
  }

  async upload(input: UploadFileInput, config?: RequestConfig): Promise<StoredFile> {
    const form = new FormData();
    form.append('file', input.file, input.filename);
    if (input.path) form.append('path', input.path);
    if (input.projectId) form.append('project_id', input.projectId);

    const headers = buildHeaders(this.core.ctx);
    headers.delete('Content-Type');

    const res = await this.core.ctx.fetchFn(`${this.core.baseUrl}/files/upload`, {
      method: 'POST',
      headers,
      body: form,
      signal: config?.signal,
    });
    if (!res.ok) {
      throw await errorFromResponse(res);
    }
    return (await res.json()) as StoredFile;
  }

  createFolder(input: CreateFolderInput, config?: RequestConfig): Promise<StoredFile> {
    return this.core.request(
      'POST',
      '/files/folders',
      {
        name: input.name,
        path: input.path,
        project_id: input.projectId,
      },
      { config },
    );
  }

  createProjectFolder(
    projectId: string,
    input: Omit<CreateFolderInput, 'projectId'>,
    config?: RequestConfig,
  ): Promise<StoredFile> {
    return this.core.request(
      'POST',
      `/projects/${encodeURIComponent(projectId)}/files/folders`,
      {
        name: input.name,
        path: input.path,
      },
      { config },
    );
  }

  ingest(id: string, input: IngestFileInput = {}, config?: RequestConfig): Promise<IngestFileResult> {
    return this.core.request(
      'POST',
      `/files/${encodeURIComponent(id)}/ingest`,
      input,
      { config },
    );
  }

  delete(id: string, config?: RequestConfig): Promise<void> {
    return this.core.request('DELETE', `/files/${encodeURIComponent(id)}`, undefined, { config });
  }
}