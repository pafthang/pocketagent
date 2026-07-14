export const IMAGE_EXTENSIONS = new Set([
  "png",
  "jpg",
  "jpeg",
  "gif",
  "webp",
  "svg",
  "bmp",
]);

export function isImageFile(ext: string): boolean {
  return IMAGE_EXTENSIONS.has(ext.toLowerCase());
}

export function isPdfFile(ext: string): boolean {
  return ext.toLowerCase() === "pdf";
}

const MAX_CACHE_SIZE = 200;
const cache = new Map<string, string>();

function evictOldest() {
  if (cache.size >= MAX_CACHE_SIZE) {
    const oldest = cache.keys().next().value;
    if (oldest !== undefined) {
      cache.delete(oldest);
    }
  }
}

/** Thumbnails require local FS or blob URLs — not available in web stub. */
export async function getThumbnail(path: string): Promise<string | null> {
  const cached = cache.get(path);
  if (cached) {
    cache.delete(path);
    cache.set(path, cached);
    return cached;
  }
  return null;
}

export function setThumbnail(path: string, dataUrl: string): void {
  evictOldest();
  cache.set(path, dataUrl);
}

export function invalidateThumbnail(path: string): void {
  cache.delete(path);
}

export function clearThumbnailCache(): void {
  cache.clear();
}