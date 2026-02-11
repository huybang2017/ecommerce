import { Observable, Subscription } from 'rxjs';
import { QueryConfig, QueryContext, QueryPlugin, QueryState } from '../types';
import { ReactiveQueryState } from './query-state';

/**
 * A single query execution unit.
 *
 * Manages its own reactive state, runs the plugin pipeline on every
 * fetch cycle, and exposes fine-grained observable selectors.
 */
export class QueryInstance<T> {
  private readonly state: ReactiveQueryState<T>;
  private readonly plugins: QueryPlugin<T>[];
  private readonly ctx: QueryContext<T>;
  private fetchSub: Subscription | null = null;
  private destroyed = false;

  constructor(private config: QueryConfig<T>) {
    this.state = new ReactiveQueryState<T>(config.initialData);
    this.plugins = config.plugins ?? [];

    // Build context object shared with every plugin
    this.ctx = {
      key: config.key,
      getState: () => this.state.get(),
      setState: (p) => this.state.set(p),
      select: <R>(sel: (s: QueryState<T>) => R) => this.state.select(sel),
      refetch: () => this.fetch(),
      invalidate: () => this.invalidate(),
      destroy: () => this.destroy(),
      meta: {},
    };

    // Lifecycle: init plugins
    this.plugins.forEach((p) => p.onInit?.(this.ctx));

    // Auto-fetch when enabled (default) and no seed data
    if (config.enabled !== false && !config.initialData) {
      this.fetch();
    }
  }

  // ─── Reactive Selectors ───────────────────────────────────

  get data$(): Observable<T | null> {
    return this.state.select((s) => s.data);
  }

  get loading$(): Observable<boolean> {
    return this.state.select((s) => s.isFetching);
  }

  get error$(): Observable<any> {
    return this.state.select((s) => s.error);
  }

  get status$(): Observable<string> {
    return this.state.select((s) => s.status);
  }

  get state$(): Observable<QueryState<T>> {
    return this.state.asObservable();
  }

  // ─── Actions ──────────────────────────────────────────────

  /** Execute the fetch pipeline (beforeFetch → network → afterFetch) */
  fetch(): void {
    if (this.destroyed) return;

    // Run beforeFetch hooks — any plugin returning `false` skips the call
    for (const plugin of this.plugins) {
      if (plugin.beforeFetch?.(this.ctx) === false) return;
    }

    // Cancel any in-flight request
    this.fetchSub?.unsubscribe();

    this.state.set({ status: 'loading', isFetching: true, error: null });

    this.fetchSub = this.config.fetchFn().subscribe({
      next: (raw) => {
        let data = raw;

        // afterFetch hooks may transform data
        for (const plugin of this.plugins) {
          const transformed = plugin.afterFetch?.(data, this.ctx);
          if (transformed !== undefined && transformed !== null) {
            data = transformed as T;
          }
        }

        this.state.set({
          data,
          error: null,
          status: 'success',
          isFetching: false,
          retryCount: 0,
        });
      },
      error: (err) => {
        this.state.set({ error: err, status: 'error', isFetching: false });
        this.plugins.forEach((p) => p.onError?.(err, this.ctx));
      },
    });
  }

  /** Force a fresh fetch (resets retry count) */
  refetch(): void {
    this.state.set({ retryCount: 0 });
    this.fetch();
  }

  /** Mark stale → notify plugins → refetch */
  invalidate(): void {
    this.state.set({ isStale: true });
    this.plugins.forEach((p) => p.onInvalidate?.(this.ctx));
    this.fetch();
  }

  /**
   * Optimistic cache mutation — returns a rollback function.
   *
   * ```ts
   * const rollback = query.optimisticUpdate(cur => cur.filter(…));
   * mutation$.subscribe({ error: () => rollback() });
   * ```
   */
  optimisticUpdate(updater: (current: T | null) => T): () => void {
    const previous = this.state.get().data;
    this.state.set({ data: updater(previous) });
    return () => this.state.set({ data: previous });
  }

  /** Synchronous snapshot of current data */
  getSnapshot(): T | null {
    return this.state.get().data;
  }

  /** Tear down: cancel subscriptions, run onDestroy hooks, complete state */
  destroy(): void {
    if (this.destroyed) return;
    this.destroyed = true;
    this.fetchSub?.unsubscribe();
    this.plugins.forEach((p) => p.onDestroy?.(this.ctx));
    this.state.destroy();
  }
}
