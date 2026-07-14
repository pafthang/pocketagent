import type {
  FileEntry,
  FileOpContext,
  DefaultDirs,
  FileChangeEvent,
  FileSystemProvider,
  RecursiveSearchResult,
  FileStatExtended,
} from "./types";

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

/**
 * Fallback local provider when no browser directory is picked.
 * Use pickBrowserLocalDirectory() or remote Files API instead.
 */
export class LocalFileSystem implements FileSystemProvider {
  scheme = "local" as const;
  supportsPaste = false;
  supportsRename = false;

  async readDir(_path: string): Promise<FileEntry[]> {
    return [];
  }

  async readFileText(_path: string, _ctx?: FileOpContext): Promise<string> {
    return "";
  }

  async writeFile(_path: string, _content: string, _ctx?: FileOpContext): Promise<void> {}

  async deleteFile(_path: string, _recursive = false, _ctx?: FileOpContext): Promise<void> {}

  async rename(_oldPath: string, _newPath: string, _ctx?: FileOpContext): Promise<void> {}

  async stat(path: string, _ctx?: FileOpContext): Promise<FileEntry> {
    return { name: "", path, isDir: false, size: 0, modified: 0, extension: "", source: "local" };
  }

  async createDir(_path: string): Promise<void> {}

  async exists(_path: string): Promise<boolean> {
    return false;
  }

  async watch(_path: string, _callback: (event: FileChangeEvent) => void): Promise<() => void> {
    return () => {};
  }

  async getDefaultDirs(): Promise<DefaultDirs> {
    return { home: "", documents: "", downloads: "", desktop: "" };
  }

  async readFileHead(_path: string, _maxBytes = 2048, _ctx?: FileOpContext): Promise<string> {
    return "";
  }

  async getThumbnail(path: string): Promise<string | null> {
    const { getThumbnail: getCachedThumbnail } = await import("./thumbnail-cache");
    return getCachedThumbnail(path);
  }

  async getFileSrc(_path: string, _ctx?: FileOpContext): Promise<string> {
    return "";
  }

  async readFileBase64(_path: string, _ctx?: FileOpContext): Promise<string> {
    return "";
  }

  async copyFile(_src: string, _dest: string): Promise<void> {}

  async copyDir(_src: string, _dest: string): Promise<void> {}

  async moveFile(_src: string, _dest: string): Promise<void> {}

  async statExtended(path: string, _ctx?: FileOpContext): Promise<FileStatExtended> {
    return { ...EMPTY_STAT, path };
  }

  async openInTerminal(_path: string): Promise<void> {}

  async searchRecursive(
    _rootPath: string,
    _query: string,
    _maxResults = 500,
    _maxDepth = 10,
  ): Promise<RecursiveSearchResult> {
    return { entries: [], totalScanned: 0, truncated: false };
  }
}

export type { FileStatExtended };
export const localFs = new LocalFileSystem();