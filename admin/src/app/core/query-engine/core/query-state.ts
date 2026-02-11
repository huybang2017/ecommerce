import { BehaviorSubject, Observable } from 'rxjs';
import { distinctUntilChanged, map } from 'rxjs/operators';
import { QueryState, initialQueryState } from '../types';

/**
 * Thin reactive wrapper around a BehaviorSubject holding QueryState<T>.
 * Provides immutable-update semantics and selector-based subscriptions.
 */
export class ReactiveQueryState<T> {
  private readonly state$: BehaviorSubject<QueryState<T>>;

  constructor(initialData?: T) {
    this.state$ = new BehaviorSubject<QueryState<T>>(initialQueryState(initialData));
  }

  /** Synchronous snapshot */
  get(): QueryState<T> {
    return this.state$.value;
  }

  /** Immutable partial update */
  set(partial: Partial<QueryState<T>>): void {
    this.state$.next({ ...this.state$.value, ...partial });
  }

  /** Derived observable — only emits when the selected slice changes */
  select<R>(selector: (state: QueryState<T>) => R): Observable<R> {
    return this.state$.pipe(map(selector), distinctUntilChanged());
  }

  /** Full state stream */
  asObservable(): Observable<QueryState<T>> {
    return this.state$.asObservable();
  }

  /** Complete the underlying subject */
  destroy(): void {
    this.state$.complete();
  }
}
