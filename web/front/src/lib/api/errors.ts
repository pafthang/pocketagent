import { ApiError } from "./types";
import { isPocketAgentError } from "@pocketagent/client";

/** Map fetch/SDK errors to user-friendly messages. */
export function friendlyErrorMessage(err: unknown): string {
  if (isPocketAgentError(err)) {
    if (err.status === 0) return "Could not reach the backend. Is it running?";
    if (err.isUnauthorized) return "Session expired. Please sign in again.";
    if (err.status === 502 || err.status === 503)
      return "Backend is starting up or temporarily unavailable. Try again in a moment.";
    if (err.status === 504) return "Request timed out. The backend may be overloaded.";
    return err.message;
  }
  if (err instanceof ApiError) {
    if (err.status === 0) return "Could not reach the backend. Is it running?";
    if (err.status === 401) return "Session expired. Please sign in again.";
    if (err.status === 502 || err.status === 503)
      return "Backend is starting up or temporarily unavailable. Try again in a moment.";
    if (err.status === 504) return "Request timed out. The backend may be overloaded.";
    if (err.detail) return err.detail;
    return err.message;
  }
  if (err instanceof DOMException && err.name === "AbortError")
    return "Request timed out. The backend may be unresponsive.";
  if (err instanceof TypeError)
    return "Could not reach the backend. Check your connection and make sure it's running.";
  if (err instanceof Error) return err.message;
  return "An unexpected error occurred.";
}