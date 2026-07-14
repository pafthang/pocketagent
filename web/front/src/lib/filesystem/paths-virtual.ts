/** Virtual path helpers for PocketAgent Files API (/projects/{id}/…). */

export function normalizeVirtualPath(raw: string): string {
  const p = raw.trim();
  if (!p || p === ".") return "/";
  const normalized = p.startsWith("/") ? p : `/${p}`;
  return normalized.replace(/\/+/g, "/").replace(/\/$/, "") || "/";
}

export function parseBrowseScope(raw: string): { path: string; projectId?: string } {
  const path = normalizeVirtualPath(raw);
  const match = path.match(/^\/projects\/([^/]+)(\/.*)?$/);
  if (match) {
    return { projectId: match[1], path };
  }
  return { path };
}

export function parentVirtualPath(raw: string): string {
  const path = normalizeVirtualPath(raw);
  if (path === "/") return "/";
  const idx = path.lastIndexOf("/");
  return idx <= 0 ? "/" : path.slice(0, idx);
}

export function isVirtualRoot(raw: string): boolean {
  const path = normalizeVirtualPath(raw);
  return path === "/";
}