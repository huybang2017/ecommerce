import { QueryPlugin } from '../types';

// ─── Public types ─────────────────────────────────────────────

export interface MetricsData {
  fetchCount: number;
  errorCount: number;
  lastFetchDuration: number | null;
  avgFetchDuration: number;
  totalDuration: number;
}

export interface MetricsOptions {
  /** Called after every fetch / error with the latest snapshot */
  onMetricsUpdate?: (key: string, metrics: MetricsData) => void;
}

// ─── Plugin ───────────────────────────────────────────────────

/**
 * Tracks fetch count, error count, and timing information.
 */
export function metricsPlugin<T>(options: MetricsOptions = {}): QueryPlugin<T> {
  const { onMetricsUpdate } = options;

  return {
    name: 'metrics',

    onInit(ctx) {
      ctx.meta['metrics'] = {
        fetchCount: 0,
        errorCount: 0,
        lastFetchDuration: null,
        avgFetchDuration: 0,
        totalDuration: 0,
      } as MetricsData;
    },

    beforeFetch(ctx) {
      ctx.meta['metrics:fetchStart'] = performance.now();
    },

    afterFetch(_data, ctx) {
      const start = ctx.meta['metrics:fetchStart'] as number;
      const duration = performance.now() - start;
      const m = ctx.meta['metrics'] as MetricsData;
      m.fetchCount++;
      m.lastFetchDuration = duration;
      m.totalDuration += duration;
      m.avgFetchDuration = m.totalDuration / m.fetchCount;
      onMetricsUpdate?.(ctx.key, { ...m });
    },

    onError(_error, ctx) {
      const m = ctx.meta['metrics'] as MetricsData;
      m.errorCount++;
      onMetricsUpdate?.(ctx.key, { ...m });
    },
  };
}

// ─── Helper ───────────────────────────────────────────────────

export function getMetrics(meta: Record<string, any>): MetricsData | undefined {
  return meta['metrics'];
}
