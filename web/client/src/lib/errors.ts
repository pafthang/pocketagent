export class PocketAgentError extends Error {
  readonly status: number;
  readonly body?: unknown;

  constructor(message: string, status: number, body?: unknown) {
    super(message);
    this.name = 'PocketAgentError';
    this.status = status;
    this.body = body;
  }

  get isUnauthorized(): boolean {
    return this.status === 401;
  }

  get isForbidden(): boolean {
    return this.status === 403;
  }

  get isNotFound(): boolean {
    return this.status === 404;
  }

  get isConflict(): boolean {
    return this.status === 409;
  }

  get isValidation(): boolean {
    return this.status === 400;
  }

  get isServer(): boolean {
    return this.status >= 500;
  }
}

export function isPocketAgentError(err: unknown): err is PocketAgentError {
  return err instanceof PocketAgentError;
}

export function isUnauthorized(err: unknown): boolean {
  return isPocketAgentError(err) && err.isUnauthorized;
}

export function isForbidden(err: unknown): boolean {
  return isPocketAgentError(err) && err.isForbidden;
}

export function isNotFound(err: unknown): boolean {
  return isPocketAgentError(err) && err.isNotFound;
}

export async function errorFromResponse(res: Response): Promise<PocketAgentError> {
  let body: unknown;
  const text = await res.text();
  if (text) {
    try {
      body = JSON.parse(text);
    } catch {
      body = text;
    }
  }

  const message =
    typeof body === 'object' &&
    body !== null &&
    'error' in body &&
    typeof (body as { error: unknown }).error === 'string'
      ? (body as { error: string }).error
      : res.statusText || `HTTP ${res.status}`;

  return new PocketAgentError(message, res.status, body);
}