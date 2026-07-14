export { LocalFileSystem, localFs } from "./local";
export { RemoteFileSystem, remoteFs } from "./remote";
export {
  getFileSystem,
  defaultExplorerSource,
  entryContext,
  browserLocalPickerSupported,
  pickBrowserLocalDirectory,
  clearBrowserLocalDirectory,
  hasBrowserLocalDirectory,
  getBrowserLocalRootName,
} from "./adapter";
export type { FileStatExtended } from "./local";
export type {
  FileEntry,
  FileOpContext,
  FileSource,
  DefaultDirs,
  FileChangeEvent,
  FileSystemProvider,
  RecursiveSearchResult,
} from "./types";
export {
  getThumbnail,
  invalidateThumbnail,
  clearThumbnailCache,
  isImageFile,
  isPdfFile,
  IMAGE_EXTENSIONS,
} from "./thumbnail-cache";
export {
  resolvePath,
  getParentDir,
  parentDir,
  joinPath,
  isAbsolute,
  normalizeSeparators,
  getExtension,
  getFileName,
} from "./paths";
export {
  normalizeVirtualPath,
  parseBrowseScope,
  parentVirtualPath,
  isVirtualRoot,
} from "./paths-virtual";

import { base64DataUrlToArrayBuffer } from "./binary-utils";
export { base64DataUrlToArrayBuffer };