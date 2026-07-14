/** Cross-platform path utilities (string-based; no native FS). */

export function normalizeSeparators(p: string): string {
  return p.replace(/\\/g, "/");
}

export function parentDir(filePath: string): string {
  const normalized = normalizeSeparators(filePath);
  const lastSlash = normalized.lastIndexOf("/");
  if (lastSlash <= 0) return normalized.slice(0, lastSlash + 1) || "/";
  return normalized.slice(0, lastSlash);
}

export function joinPath(...segments: string[]): string {
  if (segments.length === 0) return "";
  let result = segments[0];
  for (let i = 1; i < segments.length; i++) {
    const seg = segments[i];
    if (!seg) continue;
    if (seg.startsWith("/") || /^[A-Za-z]:[\\/]/.test(seg)) {
      result = seg;
      continue;
    }
    const sep = result.endsWith("/") || result.endsWith("\\") ? "" : "/";
    result = result + sep + seg;
  }
  return result;
}

export function isAbsolute(p: string): boolean {
  return p.startsWith("/") || /^[A-Za-z]:[\\/]/.test(p);
}

export async function resolvePath(path: string, baseDir?: string): Promise<string> {
  if (isAbsolute(path)) return path;
  if (baseDir) return joinPath(baseDir, path);
  return path;
}

export async function getParentDir(filePath: string): Promise<string> {
  return parentDir(filePath);
}

export function getExtension(filePath: string): string {
  const name = filePath.split(/[\\/]/).pop() ?? "";
  const dotIdx = name.lastIndexOf(".");
  if (dotIdx <= 0) return "";
  return name.slice(dotIdx + 1).toLowerCase();
}

export function getFileName(filePath: string): string {
  return filePath.split(/[\\/]/).pop() ?? "";
}