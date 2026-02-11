import { QueryPlugin } from '../types';
import { cachePlugin, CacheOptions } from '../plugins/cache.plugin';
import { retryPlugin, RetryOptions } from '../plugins/retry.plugin';
import { refetchPlugin, RefetchOptions } from '../plugins/refetch.plugin';
import { loggingPlugin, LoggingOptions } from '../plugins/logging.plugin';

export interface DefaultPresetOptions {
  /** Pass `false` to disable; pass options to customise */
  cache?: CacheOptions | false;
  retry?: RetryOptions | false;
  refetch?: RefetchOptions | false;
  logging?: LoggingOptions | false;
}

/**
 * Convenience factory — returns a sensible plugin stack for most queries.
 *
 * ```ts
 * engine.query({
 *   key: 'products',
 *   fetchFn: () => api.list(),
 *   plugins: defaultPreset({ cache: { ttl: 60_000 } }),
 * });
 * ```
 */
export function defaultPreset<T>(options: DefaultPresetOptions = {}): QueryPlugin<T>[] {
  const plugins: QueryPlugin<T>[] = [];

  if (options.cache !== false) plugins.push(cachePlugin(options.cache));
  if (options.retry !== false) plugins.push(retryPlugin(options.retry));
  if (options.refetch !== false) plugins.push(refetchPlugin(options.refetch));
  if (options.logging !== false) plugins.push(loggingPlugin(options.logging));

  return plugins;
}
