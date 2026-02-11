import { Injectable } from '@angular/core';
import { QueryConfig } from '../types';
import { QueryInstance } from './query-instance';

/**
 * Global singleton registry that de-duplicates QueryInstances by key.
 *
 * If two components request the same key they get the SAME instance,
 * sharing cache, loading state, and plugin side-effects.
 */
@Injectable({ providedIn: 'root' })
export class QueryRegistry {
  private readonly queries = new Map<string, QueryInstance<any>>();

  /** Return an existing instance or create (and register) a new one */
  getOrCreate<T>(config: QueryConfig<T>): QueryInstance<T> {
    if (this.queries.has(config.key)) {
      return this.queries.get(config.key)! as QueryInstance<T>;
    }
    const instance = new QueryInstance(config);
    this.queries.set(config.key, instance);
    return instance;
  }

  get<T>(key: string): QueryInstance<T> | undefined {
    return this.queries.get(key) as QueryInstance<T> | undefined;
  }

  /** Invalidate one key or every key matching a RegExp */
  invalidate(key: string | RegExp): void {
    if (typeof key === 'string') {
      this.queries.get(key)?.invalidate();
    } else {
      this.queries.forEach((instance, k) => {
        if (key.test(k)) instance.invalidate();
      });
    }
  }

  /** Batch-invalidate multiple keys */
  invalidateKeys(keys: string[]): void {
    keys.forEach((k) => this.invalidate(k));
  }

  /** Destroy and remove a single query */
  remove(key: string): void {
    const instance = this.queries.get(key);
    if (instance) {
      instance.destroy();
      this.queries.delete(key);
    }
  }

  /** Destroy everything (e.g. on logout) */
  clear(): void {
    this.queries.forEach((i) => i.destroy());
    this.queries.clear();
  }

  has(key: string): boolean {
    return this.queries.has(key);
  }

  keys(): string[] {
    return Array.from(this.queries.keys());
  }
}
