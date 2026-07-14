// Persist OAuth tokens in localStorage (web-only; legacy until Phase 1 JWT auth).

const STORAGE_KEY = "pocketagent_oauth_tokens";

export interface OAuthTokens {
  access_token: string;
  refresh_token: string | null;
  expires_at: number;
  scopes: string[];
}

export async function readTokens(): Promise<OAuthTokens | null> {
  if (typeof localStorage === "undefined") return null;
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return null;
    return JSON.parse(raw) as OAuthTokens;
  } catch {
    return null;
  }
}

export async function saveTokens(tokens: OAuthTokens): Promise<void> {
  if (typeof localStorage === "undefined") return;
  localStorage.setItem(STORAGE_KEY, JSON.stringify(tokens));
}

export async function clearTokens(): Promise<void> {
  if (typeof localStorage === "undefined") return;
  localStorage.removeItem(STORAGE_KEY);
}