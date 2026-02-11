import { QueryPlugin } from "../types";

export interface CacheOptions {
  /** Time-to-live in ms. Data older than TTL triggers a fresh fetch. Default: 5 min */
  ttl?: number;
  /** Data younger than staleTime is considered fresh — fetch is skipped. Default: 0 */
  staleTime?: number;
}

/**
 * Prevents redundant network calls by honoring TTL and staleTime.
 *
 * - `staleTime = 0`  → data is immediately stale, background refetch allowed
 * - `staleTime = ttl` → data is fresh for the whole TTL window
 */
export function cachePlugin<T>(options: CacheOptions = {}): QueryPlugin<T> {
  const { ttl = 5 * 60 * 1000, staleTime = 0 } = options;

  return {
    name: "cache",

    beforeFetch(ctx) {
      const { data, fetchedAt } = ctx.getState();
      if (data !== null && fetchedAt !== null) {
        const age = Date.now() - fetchedAt;

        // Expired — allow fetch (data will be replaced)
        if (age > ttl) return;

        // Fresh — skip fetch entirely
        if (age <= staleTime) {
          ctx.setState({ isStale: false });
          return false;
        }

        // Stale but within TTL — mark stale, allow background refetch
        ctx.setState({ isStale: true });
        return;
      }
      return;
    },

    afterFetch(_data, ctx) {
      ctx.setState({ fetchedAt: Date.now(), isStale: false });
    },
  };
}
