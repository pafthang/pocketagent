// See https://svelte.dev/docs/kit/types#app.d.ts
declare global {
  interface Window {
    showDirectoryPicker?(options?: {
      mode?: "read" | "readwrite";
    }): Promise<FileSystemDirectoryHandle>;
  }

  interface FileSystemDirectoryHandle {
    entries(): AsyncIterableIterator<[string, FileSystemHandle]>;
  }
}

export {};