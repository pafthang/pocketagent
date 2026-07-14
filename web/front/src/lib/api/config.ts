/**
 * Gate API base URL.
 * - Set `VITE_GATE_URL=http://127.0.0.1:8080` for direct gate access.
 * - Set `VITE_GATE_URL=` (empty) in dev to use Vite proxy (same-origin relative paths).
 */
function resolveGateUrl(): string {
  const env = import.meta.env.VITE_GATE_URL as string | undefined;
  if (env === "") return "";
  if (env) return env.replace(/\/+$/, "");
  return "http://127.0.0.1:8080";
}

export const GATE_URL = resolveGateUrl();

/** @deprecated Legacy PocketPaw — use GATE_URL. */
export const BACKEND_URL = GATE_URL || "http://127.0.0.1:8080";

/** @deprecated Legacy PocketPaw API prefix — not used by gate. */
export const API_PREFIX = "";

/** @deprecated Legacy PocketPaw API base. */
export const API_BASE = BACKEND_URL;