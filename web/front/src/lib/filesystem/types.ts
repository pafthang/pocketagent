export type FileSource = "local" | "remote" | "cloud";

export interface FileEntry {
  name: string;
  path: string;
  isDir: boolean;
  size: number;
  modified: number;
  extension: string;
  source: FileSource;
  /** PocketAgent Files API record id (remote). */
  id?: string;
  mimeType?: string;
  projectId?: string;
}

/** Optional context for file operations (remote id lookup). */
export type FileOpContext = Pick<
  FileEntry,
  "id" | "path" | "isDir" | "name" | "projectId" | "mimeType" | "source"
>;

export interface DefaultDirs {
  home: string;
  documents: string;
  downloads: string;
  desktop: string;
}

export interface FileChangeEvent {
  path: string;
  kind: "create" | "modify" | "delete";
  is_dir: boolean;
}

export interface RecursiveSearchResult {
  entries: FileEntry[];
  totalScanned: number;
  truncated: boolean;
}

export interface FileStatExtended {
  name: string;
  path: string;
  isDir: boolean;
  size: number;
  modified: number;
  created: number;
  extension: string;
  readonly: boolean;
  isSymlink: boolean;
}

export interface FileSystemProvider {
  scheme: FileSource;
  readDir(path: string): Promise<FileEntry[]>;
  readFileText(path: string, ctx?: FileOpContext): Promise<string>;
  writeFile(path: string, content: string, ctx?: FileOpContext): Promise<void>;
  deleteFile(path: string, recursive?: boolean, ctx?: FileOpContext): Promise<void>;
  rename(oldPath: string, newPath: string, ctx?: FileOpContext): Promise<void>;
  stat(path: string, ctx?: FileOpContext): Promise<FileEntry>;
  createDir(path: string): Promise<void>;
  exists(path: string): Promise<boolean>;
  watch(path: string, callback: (event: FileChangeEvent) => void): Promise<() => void>;
  getDefaultDirs(): Promise<DefaultDirs>;
  readFileHead(path: string, maxBytes?: number, ctx?: FileOpContext): Promise<string>;
  getFileSrc(path: string, ctx?: FileOpContext): Promise<string>;
  readFileBase64(path: string, ctx?: FileOpContext): Promise<string>;
  copyFile(src: string, dest: string): Promise<void>;
  copyDir(src: string, dest: string): Promise<void>;
  moveFile(src: string, dest: string): Promise<void>;
  statExtended(path: string, ctx?: FileOpContext): Promise<FileStatExtended>;
  openInTerminal(path: string): Promise<void>;
  searchRecursive(
    rootPath: string,
    query: string,
    maxResults?: number,
    maxDepth?: number,
  ): Promise<RecursiveSearchResult>;
  uploadFiles?(dirPath: string, files: File[], projectId?: string): Promise<void>;
  supportsPaste?: boolean;
  supportsRename?: boolean;
}
