import type { BrowseFileEntry } from "@pocketagent/client";
import { connectionStore } from "$lib/stores/connection.svelte";
import type {
  DefaultDirs,
  FileChangeEvent,
  FileEntry,
  FileOpContext,
  FileStatExtended,
  FileSystemProvider,
  RecursiveSearchResult,
} from "./types";
import { getExtension, getFileName } from "./paths";
import {
  isVirtualRoot,
  normalizeVirtualPath,
  parentVirtualPath,
  parseBrowseScope,
} from "./paths-virtual";

function toModified(iso?: string): number {
  if (!iso) return 0;
  const ms = Date.parse(iso);
  return Number.isFinite(ms) ? ms : 0;
}

function browseToEntry(entry: BrowseFileEntry): FileEntry {
  return {
    id: entry.id,
    name: entry.name,
    path: entry.path,
    isDir: entry.is_dir,
    size: entry.size ?? 0,
    modified: toModified(entry.modified_at),
    extension: getExtension(entry.name),
    source: "remote",
    mimeType: entry.mime_type,
    projectId: entry.project_id,
  };
}

function ctxKey(ctx?: FileOpContext): string | undefined {
  return ctx?.id ?? (ctx?.path ? ctx.path : undefined);
}

export class RemoteFileSystem implements FileSystemProvider {
  scheme = "remote" as const;
  supportsPaste = false;
  supportsRename = false;

  private entryByPath = new Map<string, FileEntry>();
  private blobUrlCache = new Map<string, string>();

  private client() {
    return connectionStore.getAgentClient();
  }

  private resolveId(path: string, ctx?: FileOpContext): string {
    const id = ctx?.id ?? this.entryByPath.get(path)?.id;
    if (!id) throw new Error(`Remote file id not found for ${path}`);
    return id;
  }

  private scopeFor(path: string) {
    return parseBrowseScope(path);
  }

  async readDir(path: string): Promise<FileEntry[]> {
    const { path: browsePath, projectId } = this.scopeFor(path);
    const res = await this.client().files.browse({
      path: isVirtualRoot(browsePath) ? "/" : browsePath,
      projectId,
    });

    const entries = res.files.map(browseToEntry);
    for (const entry of entries) {
      if (entry.id) this.entryByPath.set(entry.path, entry);
    }
    return entries;
  }

  async readFileText(path: string, ctx?: FileOpContext): Promise<string> {
    const id = this.resolveId(path, ctx);
    const res = await this.client().files.content(id);
    return res.content;
  }

  async readFileHead(path: string, maxBytes = 2048, ctx?: FileOpContext): Promise<string> {
    const text = await this.readFileText(path, ctx);
    return text.slice(0, maxBytes);
  }

  private async downloadBlob(path: string, ctx?: FileOpContext): Promise<Blob> {
    const id = this.resolveId(path, ctx);
    return this.client().files.download(id);
  }

  async getFileSrc(path: string, ctx?: FileOpContext): Promise<string> {
    const key = ctxKey(ctx) ?? path;
    const cached = this.blobUrlCache.get(key);
    if (cached) return cached;

    const blob = await this.downloadBlob(path, ctx);
    const url = URL.createObjectURL(blob);
    this.blobUrlCache.set(key, url);
    return url;
  }

  async readFileBase64(path: string, ctx?: FileOpContext): Promise<string> {
    const blob = await this.downloadBlob(path, ctx);
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => resolve(String(reader.result));
      reader.onerror = () => reject(reader.error);
      reader.readAsDataURL(blob);
    });
  }

  async writeFile(path: string, content: string, ctx?: FileOpContext): Promise<void> {
    const name = getFileName(path);
    const parent = parentVirtualPath(path);
    const { projectId } = this.scopeFor(parent);
    const mime = ctx?.mimeType ?? "text/plain";

    if (ctx?.id) {
      try {
        await this.client().files.delete(ctx.id);
      } catch {
        // replace flow — continue upload
      }
    }

    const blob = new Blob([content], { type: mime });
    await this.client().files.upload({
      file: blob,
      filename: name,
      path: isVirtualRoot(parent) ? "/" : parent,
      projectId,
    });
  }

  async deleteFile(path: string, _recursive = false, ctx?: FileOpContext): Promise<void> {
    const id = this.resolveId(path, ctx);
    await this.client().files.delete(id);
    this.entryByPath.delete(path);
    const key = ctxKey(ctx) ?? path;
    const url = this.blobUrlCache.get(key);
    if (url) {
      URL.revokeObjectURL(url);
      this.blobUrlCache.delete(key);
    }
  }

  async rename(): Promise<void> {
    throw new Error("Rename is not supported for remote files");
  }

  async stat(path: string, ctx?: FileOpContext): Promise<FileEntry> {
    const cached = ctx?.id ? { ...this.entryByPath.get(path), ...ctx } : this.entryByPath.get(path);
    if (cached?.id) return cached as FileEntry;

    const parent = parentVirtualPath(path);
    const name = getFileName(path);
    const siblings = await this.readDir(parent);
    const found = siblings.find((f) => f.path === path || f.name === name);
    if (!found) throw new Error(`File not found: ${path}`);
    return found;
  }

  async createDir(path: string): Promise<void> {
    const name = getFileName(path);
    const parent = parentVirtualPath(path);
    const { projectId } = this.scopeFor(parent);
    await this.client().files.createFolder({
      name,
      path: isVirtualRoot(parent) ? "/" : parent,
      projectId,
    });
  }

  async exists(path: string): Promise<boolean> {
    try {
      await this.stat(path);
      return true;
    } catch {
      return false;
    }
  }

  async watch(_path: string, _callback: (event: FileChangeEvent) => void): Promise<() => void> {
    return () => {};
  }

  async getDefaultDirs(): Promise<DefaultDirs> {
    return { home: "/", documents: "/", downloads: "/", desktop: "/" };
  }

  async copyFile(): Promise<void> {
    throw new Error("Copy is not supported for remote files");
  }

  async copyDir(): Promise<void> {
    throw new Error("Copy is not supported for remote files");
  }

  async moveFile(): Promise<void> {
    throw new Error("Move is not supported for remote files");
  }

  async statExtended(path: string, ctx?: FileOpContext): Promise<FileStatExtended> {
    const entry = await this.stat(path, ctx);
    return {
      name: entry.name,
      path: entry.path,
      isDir: entry.isDir,
      size: entry.size,
      modified: entry.modified,
      created: entry.modified,
      extension: entry.extension,
      readonly: false,
      isSymlink: false,
    };
  }

  async openInTerminal(): Promise<void> {}

  async searchRecursive(
    rootPath: string,
    query: string,
    maxResults = 500,
    maxDepth = 6,
  ): Promise<RecursiveSearchResult> {
    const q = query.trim().toLowerCase();
    if (!q) return { entries: [], totalScanned: 0, truncated: false };

    const results: FileEntry[] = [];
    let scanned = 0;
    let truncated = false;

    const walk = async (dirPath: string, depth: number): Promise<void> => {
      if (depth > maxDepth || results.length >= maxResults) {
        truncated = true;
        return;
      }
      const entries = await this.readDir(dirPath);
      for (const entry of entries) {
        scanned++;
        if (entry.name.toLowerCase().includes(q)) {
          results.push(entry);
          if (results.length >= maxResults) {
            truncated = true;
            return;
          }
        }
        if (entry.isDir) {
          await walk(entry.path, depth + 1);
          if (truncated && results.length >= maxResults) return;
        }
      }
    };

    const start = normalizeVirtualPath(rootPath);
    await walk(start, 0);
    return { entries: results, totalScanned: scanned, truncated };
  }

  async uploadFiles(dirPath: string, files: File[], projectId?: string): Promise<void> {
    const { path: browsePath, projectId: scopedProject } = this.scopeFor(dirPath);
    const pid = projectId ?? scopedProject;
    const path = isVirtualRoot(browsePath) ? "/" : browsePath;

    for (const file of files) {
      await this.client().files.upload({
        file,
        filename: file.name,
        path,
        projectId: pid,
      });
    }
  }
}

export const remoteFs = new RemoteFileSystem();