// OAuth PKCE was used by the Tauri desktop app. Web UI uses JWT auth (Phase 1).

import type { OAuthTokens } from "./token-store";

export interface OAuthResult {
  success: boolean;
  tokens?: OAuthTokens;
  error?: string;
}

export async function startOAuthFlow(): Promise<OAuthResult> {
  return {
    success: false,
    error: "OAuth is not available in the web app. Use email/password login (Phase 1).",
  };
}

export async function revokeTokens(): Promise<void> {
  // no-op
}