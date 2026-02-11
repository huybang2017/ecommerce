import { Injectable } from '@angular/core';
import { QueryConfig, MutationConfig } from '../types';
import { QueryInstance } from './query-instance';
import { MutationInstance } from './mutation-instance';
import { QueryRegistry } from './query-registry.service';

/**
 * Main Angular DI entry-point for the query engine.
 *
 * Usage:
 * ```ts
 * private listQuery = this.engine.query<Product[]>({
 *   key: 'products',
 *   fetchFn: () => this.api.list(),
 *   plugins: [ cachePlugin({ ttl: 60_000 }), retryPlugin() ],
 * });
 * ```
 */
@Injectable({ providedIn: 'root' })
export class QueryEngine {
  constructor(private registry: QueryRegistry) {}

  /** Create (or retrieve an existing) query by key */
  query<T>(config: QueryConfig<T>): QueryInstance<T> {
    return this.registry.getOrCreate(config);
  }

  /** Create a one-shot mutation with optional auto-invalidation */
  mutation<TData, TVariables>(
    config: MutationConfig<TData, TVariables>,
  ): MutationInstance<TData, TVariables> {
    return new MutationInstance(config, (keys) => this.registry.invalidateKeys(keys));
  }

  /** Invalidate a single key or all keys matching a pattern */
  invalidate(key: string | RegExp): void {
    this.registry.invalidate(key);
  }

  /** Invalidate every registered query */
  invalidateAll(): void {
    this.registry.keys().forEach((k) => this.registry.invalidate(k));
  }

  /** Destroy and de-register a query */
  removeQuery(key: string): void {
    this.registry.remove(key);
  }

  /** Look up a live query by key (returns undefined if not registered) */
  getQuery<T>(key: string): QueryInstance<T> | undefined {
    return this.registry.get<T>(key);
  }

  /** Destroy all queries (e.g. on logout) */
  clear(): void {
    this.registry.clear();
  }
}
