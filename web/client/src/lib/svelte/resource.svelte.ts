import { isPocketAgentError, isUnauthorized } from '../errors.js';

export interface ResourceState<T> {
  readonly items: T[];
  readonly total: number;
  readonly page: number;
  readonly loading: boolean;
  readonly error: string | null;
  refresh: (page?: number) => Promise<T[] | undefined>;
  clearError: () => void;
}

export interface ResourceOptions<T> {
  page?: number;
  fetch: (page: number) => Promise<{ items: T[]; total: number } | undefined>;
  onUnauthorized?: () => void;
}

export function createResourceState<T>(opts: ResourceOptions<T>): ResourceState<T> {
  let items = $state<T[]>([]);
  let total = $state(0);
  let page = $state(opts.page ?? 1);
  let loading = $state(false);
  let error = $state<string | null>(null);

  async function refresh(nextPage = page): Promise<T[] | undefined> {
    loading = true;
    error = null;
    page = nextPage;

    try {
      const res = await opts.fetch(page);
      if (!res) return undefined;
      items = res.items;
      total = res.total;
      return items;
    } catch (err) {
      error = isPocketAgentError(err) ? err.message : String(err);
      if (isUnauthorized(err)) opts.onUnauthorized?.();
      return undefined;
    } finally {
      loading = false;
    }
  }

  return {
    get items() {
      return items;
    },
    get total() {
      return total;
    },
    get page() {
      return page;
    },
    get loading() {
      return loading;
    },
    get error() {
      return error;
    },
    refresh,
    clearError() {
      error = null;
    },
  };
}