import type {
  DefaultDirs,
  FileChangeEvent,
  FileEntry,
  FileStatExtended,
  FileSystemProvider,
  RecursiveSearchResult,
} from "./types";
import { getExtension, getFileName } from "./paths";

const EMPTY_STAT: FileStatExtended = {
  name: "",
  path: "",
  isDir: false,
  size: 0,
  modified: 0,
  created: 0,
  extension: "",
  readonly: false,
  isSymlink: false,
};

let rootHandle: FileSystemDirectoryHandle | null = null;
let rootName = "local";

export function hasBrowserLocalDirectory(): boolean {
  return rootHandle !== null;
}

export function getBrowserLocalRootName(): string {
  return rootName;
}

export function browserLocalPickerSupported(): boolean {
  return typeof window !== "undefined" && "showDirectoryPicker" in window;
}

export async function pickBrowserLocalDirectory(): Promise<boolean> {
  if (!browserLocalPickerSupported() || !window.showDirectoryPicker) return false;
  const handle = await window.showDirectoryPicker({ mode: "readwrite" });
  rootHandle = handle;
  rootName = handle.name;
  return true;
}

export function clearBrowserLocalDirectory(): void {
  rootHandle = null;
  rootName = "local";
}

async function getDirHandle(
  segments: string[],
  create = false,
): Promise<FileSystemDirectoryHandle> {
  if (!rootHandle) throw new Error("No local folder selected");
  let dir = rootHandle;
  for (const segment of segments) {
    if (!segment) continue;
    dir = await dir.getDirectoryHandle(segment, { create });
  }
  return dir;
}

function virtualPathToSegments(path: string): string[] {
  return path.replace(/\\/g, "/").split("/").filter(Boolean);
}

function joinVirtualPath(segments: string[]): string {
  return segments.length ? `/${segments.join("/")}` : "/";
}

async function entryFromHandle(
  name: string,
  handle: FileSystemHandle,
  parentSegments: string[],
): Promise<FileEntry> {
  const path = joinVirtualPath([...parentSegments, name]);
  if (handle.kind === "directory") {
    return {
      name,
      path,
      isDir: true,
      size: 0,
      modified: Date.now(),
      extension: "",
      source: "local",
    };
  }

  const file = await (handle as FileSystemFileHandle).getFile();
  return {
    name,
    path,
    isDir: false,
    size: file.size,
    modified: file.lastModified,
    extension: getExtension(name),
    source: "local",
    mimeType: file.type,
  };
}

export class BrowserLocalFileSystem implements FileSystemProvider {
  scheme = "local" as const;
  supportsPaste = true;
  supportsRename = true;

  async readDir(path: string): Promise<FileEntry[]> {
    const segments = virtualPathToSegments(path);
    const dir = await getDirHandle(segments);
    const entries: FileEntry[] = [];

    for await (const [name, handle] of (dir as FileSystemDirectoryHandle).entries()) {
      entries.push(await entryFromHandle(name, handle, segments));
    }
    return entries.sort((a, b) => {
      if (a.isDir !== b.isDir) return a.isDir ? -1 : 1;
      return a.name.localeCompare(b.name);
    });
  }

  private async getFileHandle(path: string, create = false): Promise<FileSystemFileHandle> {
    const segments = virtualPathToSegments(path);
    const fileName = segments.pop();
    if (!fileName) throw new Error("Invalid file path");
    const dir = await getDirHandle(segments, create);
    return dir.getFileHandle(fileName, { create });
  }

  async readFileText(path: string): Promise<string> {
    const handle = await this.getFileHandle(path);
    const file = await handle.getFile();
    return file.text();
  }

  async readFileHead(path: string, maxBytes = 2048): Promise<string> {
    const text = await this.readFileText(path);
    return text.slice(0, maxBytes);
  }

  async readFileBase64(path: string): Promise<string> {
    const handle = await this.getFileHandle(path);
    const file = await handle.getFile();
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => resolve(String(reader.result));
      reader.onerror = () => reject(reader.error);
      reader.readAsDataURL(file);
    });
  }

  async getFileSrc(path: string): Promise<string> {
    const handle = await this.getFileHandle(path);
    const file = await handle.getFile();
    return URL.createObjectURL(file);
  }

  async writeFile(path: string, content: string): Promise<void> {
    const handle = await this.getFileHandle(path, true);
    const writable = await handle.createWritable();
    await writable.write(content);
    await writable.close();
  }

  async deleteFile(path: string, recursive = false): Promise<void> {
    const segments = virtualPathToSegments(path);
    const name = segments.pop();
    if (!name) return;
    const dir = await getDirHandle(segments);
    await dir.removeEntry(name, { recursive });
  }

  async rename(oldPath: string, newPath: string): Promise<void> {
    const content = await this.readFileText(oldPath);
    await this.writeFile(newPath, content);
    await this.deleteFile(oldPath);
  }

  async stat(path: string): Promise<FileEntry> {
    const segments = virtualPathToSegments(path);
    const name = segments.pop();
    if (!name) throw new Error("Invalid path");
    const dir = await getDirHandle(segments);
    const handle = await dir.getFileHandle(name).catch(async () => dir.getDirectoryHandle(name));
    return entryFromHandle(name, handle, segments);
  }

  async createDir(path: string): Promise<void> {
    const segments = virtualPathToSegments(path);
    await getDirHandle(segments, true);
  }

  async exists(path: string): Promise<boolean> {
    try {
      await this.stat(path);
      return true;
    } catch {
      return false;
    }
  }

  async watch(): Promise<() => void> {
    return () => {};
  }

  async getDefaultDirs(): Promise<DefaultDirs> {
    const root = rootHandle ? `/${rootName}` : "";
    return { home: root, documents: root, downloads: root, desktop: root };
  }

  async copyFile(): Promise<void> {
    throw new Error("Copy not implemented for browser local");
  }

  async copyDir(): Promise<void> {
    throw new Error("Copy not implemented for browser local");
  }

  async moveFile(oldPath: string, newPath: string): Promise<void> {
    await this.rename(oldPath, newPath);
  }

  async statExtended(path: string): Promise<FileStatExtended> {
    const entry = await this.stat(path);
    return { ...EMPTY_STAT, ...entry, created: entry.modified };
  }

  async openInTerminal(): Promise<void> {}

  async searchRecursive(
    rootPath: string,
    query: string,
    maxResults = 500,
    maxDepth = 8,
  ): Promise<RecursiveSearchResult> {
    const q = query.trim().toLowerCase();
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
        if (entry.name.toLowerCase().includes(q)) results.push(entry);
        if (results.length >= maxResults) {
          truncated = true;
          return;
        }
        if (entry.isDir) await walk(entry.path, depth + 1);
      }
    };

    await walk(rootPath || "/", 0);
    return { entries: results, totalScanned: scanned, truncated };
  }

  async uploadFiles(): Promise<void> {}
}

export const browserLocalFs = new BrowserLocalFileSystem();