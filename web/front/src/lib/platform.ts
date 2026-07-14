// Web-only platform helpers (replaces former Tauri bridge).

import type { ChatMessage } from "$lib/api";
import type { Settings } from "$lib/api/types";

export function openExternalUrl(url: string): void {
  if (typeof window !== "undefined") {
    window.open(url, "_blank", "noopener,noreferrer");
  }
}

/** Local file paths cannot be opened from the browser without a Files API backend. */
export function openExternalPath(_path: string): void {
  // no-op until Files API is implemented
}

export function registerHotkeys(_handlers: {
  onQuickAsk?: () => void | Promise<void>;
  onToggleSidePanel?: () => void | Promise<void>;
}): void {}

export function unregisterHotkeys(): void {}

export function setupTrayListeners(_handlers: { onNavigate?: (path: string) => void }): void {}

export function cleanupTrayListeners(): void {}

export function requestNotificationPermission(): void {
  if (typeof Notification === "undefined") return;
  if (Notification.permission === "default") {
    Notification.requestPermission().catch(() => {});
  }
}

export function notifyAgentComplete(message: string): void {
  if (typeof document === "undefined" || !document.hidden) return;
  if (typeof Notification === "undefined" || Notification.permission !== "granted") return;
  try {
    new Notification("PocketAgent", { body: message });
  } catch {
    // ignore
  }
}

export function emitSessionSwitch(_sessionId: string): void {}

export function emitChatSync(_payload: {
  messages: ChatMessage[];
  streaming: boolean;
  streamingContent: string;
}): void {}

export function emitSettingsUpdate(_settings: Settings): void {}

export function onSidePanelReady(_handler: () => void): void {}

export function onChatSync(
  _handler: (payload: {
    messages: ChatMessage[];
    streaming: boolean;
    streamingContent: string;
  }) => void,
): void {}

export function disposeAllBridgeListeners(): void {}