// ─── Core ─────────────────────────────────────────────────────
export { QueryEngine } from './core/query-engine.service';
export { QueryRegistry } from './core/query-registry.service';
export { QueryInstance } from './core/query-instance';
export { MutationInstance } from './core/mutation-instance';
export { ReactiveQueryState } from './core/query-state';

// ─── Builders ─────────────────────────────────────────────────
export { QueryBuilder } from './builder/query.builder';
export { MutationBuilder } from './builder/mutation.builder';

// ─── Types ────────────────────────────────────────────────────
export type {
  QueryState,
  QueryStatus,
  QueryConfig,
  QueryContext,
  QueryPlugin,
  MutationState,
  MutationStatus,
  MutationConfig,
  MutationContext,
  MutationPlugin,
} from './types';
export { initialQueryState, initialMutationState } from './types';

// ─── Plugins ──────────────────────────────────────────────────
export { cachePlugin } from './plugins/cache.plugin';
export type { CacheOptions } from './plugins/cache.plugin';

export { retryPlugin } from './plugins/retry.plugin';
export type { RetryOptions } from './plugins/retry.plugin';

export { refetchPlugin } from './plugins/refetch.plugin';
export type { RefetchOptions } from './plugins/refetch.plugin';

export { errorPlugin } from './plugins/error.plugin';
export type { ErrorOptions } from './plugins/error.plugin';

export { pollingPlugin } from './plugins/polling.plugin';
export type { PollingOptions } from './plugins/polling.plugin';

export { paginationPlugin, getPaginationMeta } from './plugins/pagination.plugin';
export type { PaginationOptions, PaginationMeta } from './plugins/pagination.plugin';

export { loggingPlugin, mutationLoggingPlugin } from './plugins/logging.plugin';
export type { LoggingOptions } from './plugins/logging.plugin';

export { metricsPlugin, getMetrics } from './plugins/metrics.plugin';
export type { MetricsOptions, MetricsData } from './plugins/metrics.plugin';

export { syncPlugin } from './plugins/sync.plugin';
export type { SyncOptions } from './plugins/sync.plugin';

// ─── Presets ──────────────────────────────────────────────────
export { defaultPreset } from './presets/default.preset';
export type { DefaultPresetOptions } from './presets/default.preset';
