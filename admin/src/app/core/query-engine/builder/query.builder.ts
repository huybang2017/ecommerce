import { Observable } from 'rxjs';
import { QueryConfig, QueryPlugin } from '../types';
import { QueryInstance } from '../core/query-instance';

/**
 * Fluent builder for composing a QueryInstance.
 *
 * ```ts
 * const query = QueryBuilder.for<Product[]>('products')
 *   .fetchWith(() => api.list())
 *   .use(cachePlugin(), retryPlugin())
 *   .build();
 * ```
 */
export class QueryBuilder<T> {
  private config: Partial<QueryConfig<T>> = { plugins: [] };

  static for<T>(key: string): QueryBuilder<T> {
    const b = new QueryBuilder<T>();
    b.config.key = key;
    return b;
  }

  fetchWith(fn: (params?: any) => Observable<T>): this {
    this.config.fetchFn = fn;
    return this;
  }

  use(...plugins: QueryPlugin<T>[]): this {
    this.config.plugins!.push(...plugins);
    return this;
  }

  initialData(data: T): this {
    this.config.initialData = data;
    return this;
  }

  enabled(flag: boolean): this {
    this.config.enabled = flag;
    return this;
  }

  /** Build a standalone QueryInstance (not registered in QueryRegistry) */
  build(): QueryInstance<T> {
    this.assertValid();
    return new QueryInstance(this.config as QueryConfig<T>);
  }

  /** Export as plain config — pass to QueryEngine.query() for registry-managed lifecycle */
  toConfig(): QueryConfig<T> {
    this.assertValid();
    return this.config as QueryConfig<T>;
  }

  private assertValid(): void {
    if (!this.config.key) throw new Error('[QueryBuilder] key is required');
    if (!this.config.fetchFn) throw new Error('[QueryBuilder] fetchFn is required');
  }
}
