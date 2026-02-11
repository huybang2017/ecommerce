import { Observable } from 'rxjs';

// ─── Query Status ─────────────────────────────────────────────
export type QueryStatus = 'idle' | 'loading' | 'success' | 'error';

// ─── Query State ──────────────────────────────────────────────
export interface QueryState<T> {
  data: T | null;
  error: any | null;
  status: QueryStatus;
  isFetching: boolean;
  isStale: boolean;
  fetchedAt: number | null;
  retryCount: number;
  /** Extensible metadata bag — plugins read/write freely */
  meta: Record<string, any>;
}

export function initialQueryState<T>(initialData?: T): QueryState<T> {
  return {
    data: initialData ?? null,
    error: null,
    status: initialData ? 'success' : 'idle',
    isFetching: false,
    isStale: false,
    fetchedAt: initialData ? Date.now() : null,
    retryCount: 0,
    meta: {},
  };
}

// ─── Query Config ─────────────────────────────────────────────
export interface QueryConfig<T> {
  /** Unique cache key — also used for cross-query invalidation */
  key: string;
  /** The async data-fetching function */
  fetchFn: (params?: any) => Observable<T>;
  /** Plugins to compose into this query's lifecycle */
  plugins?: QueryPlugin<T>[];
  /** If false the query will NOT auto-fetch on creation. Default: true */
  enabled?: boolean;
  /** Seed data used before the first fetch completes */
  initialData?: T;
}

// ─── Query Context (passed to every plugin hook) ──────────────
export interface QueryContext<T> {
  readonly key: string;
  getState(): QueryState<T>;
  setState(partial: Partial<QueryState<T>>): void;
  select<R>(selector: (state: QueryState<T>) => R): Observable<R>;
  refetch(): void;
  invalidate(): void;
  destroy(): void;
  /** Shared bag for inter-plugin communication */
  meta: Record<string, any>;
}

// ─── Query Plugin ─────────────────────────────────────────────
export interface QueryPlugin<T = any> {
  readonly name: string;
  /** Fires once when the QueryInstance is created */
  onInit?(ctx: QueryContext<T>): void;
  /** Fires before every fetch — return `false` to skip the network call */
  beforeFetch?(ctx: QueryContext<T>): boolean | void;
  /** Fires after a successful fetch — may transform data in-place */
  afterFetch?(data: T, ctx: QueryContext<T>): T | void;
  /** Fires on fetch error */
  onError?(error: any, ctx: QueryContext<T>): void;
  /** Fires when the query is explicitly invalidated */
  onInvalidate?(ctx: QueryContext<T>): void;
  /** Fires when the QueryInstance is destroyed (cleanup timers, listeners) */
  onDestroy?(ctx: QueryContext<T>): void;
}

// ─── Mutation Status ──────────────────────────────────────────
export type MutationStatus = 'idle' | 'loading' | 'success' | 'error';

// ─── Mutation State ───────────────────────────────────────────
export interface MutationState<TData> {
  data: TData | null;
  error: any | null;
  status: MutationStatus;
  isLoading: boolean;
}

export function initialMutationState<TData>(): MutationState<TData> {
  return { data: null, error: null, status: 'idle', isLoading: false };
}

// ─── Mutation Config ──────────────────────────────────────────
export interface MutationConfig<TData, TVariables> {
  mutationFn: (variables: TVariables) => Observable<TData>;
  plugins?: MutationPlugin<TData, TVariables>[];
  onSuccess?: (data: TData, variables: TVariables) => void;
  onError?: (error: any, variables: TVariables) => void;
  onSettled?: (data: TData | null, error: any | null, variables: TVariables) => void;
  /** Query keys to auto-invalidate after a successful mutation */
  invalidateKeys?: string[];
}

// ─── Mutation Context ─────────────────────────────────────────
export interface MutationContext<TData, TVariables> {
  getState(): MutationState<TData>;
  setState(partial: Partial<MutationState<TData>>): void;
  meta: Record<string, any>;
}

// ─── Mutation Plugin ──────────────────────────────────────────
export interface MutationPlugin<TData = any, TVariables = any> {
  readonly name: string;
  beforeMutate?(variables: TVariables, ctx: MutationContext<TData, TVariables>): void;
  afterMutate?(data: TData, variables: TVariables, ctx: MutationContext<TData, TVariables>): void;
  onError?(error: any, variables: TVariables, ctx: MutationContext<TData, TVariables>): void;
  onSettled?(ctx: MutationContext<TData, TVariables>): void;
}
