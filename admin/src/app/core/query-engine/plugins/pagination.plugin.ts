import { QueryPlugin } from '../types';

// ─── Public types ─────────────────────────────────────────────

export interface PaginationMeta {
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
  hasNextPage: boolean;
  hasPreviousPage: boolean;
}

export interface PaginationOptions {
  /** Items per page. Default: 20 */
  pageSize?: number;
  /** Extract pagination metadata from the raw API response */
  extractMeta?: (response: any) => PaginationMeta;
  /** Extract the items array from the raw API response */
  extractItems?: (response: any) => any[];
}

// ─── Plugin ───────────────────────────────────────────────────

/**
 * Adds server-side pagination awareness to a query.
 *
 * After each fetch the plugin writes `PaginationMeta` into `ctx.meta['pagination']`
 * and optionally extracts just the items array from a wrapped response.
 */
export function paginationPlugin<T>(options: PaginationOptions = {}): QueryPlugin<T> {
  const { pageSize = 20, extractMeta, extractItems } = options;

  return {
    name: 'pagination',

    onInit(ctx) {
      ctx.meta['pagination'] = {
        page: 1,
        pageSize,
        total: 0,
        totalPages: 0,
        hasNextPage: false,
        hasPreviousPage: false,
      } as PaginationMeta;
    },

    afterFetch(data, ctx) {
      if (extractMeta) {
        ctx.meta['pagination'] = extractMeta(data);
      }
      if (extractItems) {
        return extractItems(data) as unknown as T;
      }
      return;
    },
  };
}

// ─── Helper ───────────────────────────────────────────────────

/** Read the current pagination metadata from a query's context */
export function getPaginationMeta(meta: Record<string, any>): PaginationMeta | undefined {
  return meta['pagination'];
}
