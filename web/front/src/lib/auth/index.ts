export { generateCodeVerifier, generateCodeChallenge, generateState } from "./pkce";
export { readTokens, saveTokens, clearTokens, type OAuthTokens } from "./token-store";
export { startOAuthFlow, revokeTokens, type OAuthResult } from "./oauth-flow";
export {
  refreshAccessToken,
  scheduleTokenRefresh,
  cancelScheduledRefresh,
} from "./token-refresh";

/**
 * Get a valid access token from localStorage. Refreshes if expired.
 * Returns null if unavailable (caller should trigger login in Phase 1).
 */
export async function getValidToken(): Promise<string | null> {
  const { readTokens } = await import("./token-store");
  const tokens = await readTokens();
  if (!tokens) return null;

  const nowS = Math.floor(Date.now() / 1000);
  if (tokens.expires_at > nowS + 30) {
    return tokens.access_token;
  }

  if (tokens.refresh_token) {
    try {
      const { refreshAccessToken } = await import("./token-refresh");
      const newTokens = await refreshAccessToken(tokens);
      return newTokens.access_token;
    } catch {
      return null;
    }
  }

  return null;
}