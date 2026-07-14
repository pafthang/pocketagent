import { connectionStore } from "$lib/stores/connection.svelte";
import type { FileSource, FileSystemProvider } from "./types";
import { localFs } from "./local";
import { remoteFs } from "./remote";
import {
  browserLocalFs,
  hasBrowserLocalDirectory,
  browserLocalPickerSupported,
  pickBrowserLocalDirectory,
  clearBrowserLocalDirectory,
  getBrowserLocalRootName,
} from "./browser-local";

export {
  browserLocalPickerSupported,
  pickBrowserLocalDirectory,
  clearBrowserLocalDirectory,
  hasBrowserLocalDirectory,
  getBrowserLocalRootName,
};

export function defaultExplorerSource(): FileSource {
  if (connectionStore.isAuthenticated && connectionStore.isConnected) {
    return "remote";
  }
  if (hasBrowserLocalDirectory()) {
    return "local";
  }
  return "remote";
}

export function getFileSystem(source: FileSource): FileSystemProvider {
  switch (source) {
    case "remote":
      return remoteFs;
    case "local":
      return hasBrowserLocalDirectory() ? browserLocalFs : localFs;
    case "cloud":
      return localFs;
    default:
      return remoteFs;
  }
}

export function entryContext(
  entry: Pick<import("./types").FileEntry, "id" | "path" | "isDir" | "name" | "projectId" | "mimeType" | "source">,
) {
  return entry;
}