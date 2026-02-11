import { BehaviorSubject, Observable, Subscription } from 'rxjs';
import { distinctUntilChanged, map } from 'rxjs/operators';
import {
  MutationConfig,
  MutationContext,
  MutationPlugin,
  MutationState,
  initialMutationState,
} from '../types';

/**
 * A single mutation execution unit.
 *
 * Mirrors QueryInstance but for write operations:
 * beforeMutate → network → afterMutate / onError → onSettled → auto-invalidate.
 */
export class MutationInstance<TData, TVariables> {
  private readonly state$ = new BehaviorSubject<MutationState<TData>>(initialMutationState());
  private readonly plugins: MutationPlugin<TData, TVariables>[];
  private readonly ctx: MutationContext<TData, TVariables>;
  private activeSub: Subscription | null = null;
  private readonly onInvalidate?: (keys: string[]) => void;

  constructor(
    private config: MutationConfig<TData, TVariables>,
    invalidateCallback?: (keys: string[]) => void,
  ) {
    this.plugins = config.plugins ?? [];
    this.onInvalidate = invalidateCallback;

    this.ctx = {
      getState: () => this.state$.value,
      setState: (p) => this.state$.next({ ...this.state$.value, ...p }),
      meta: {},
    };
  }

  // ─── Reactive Selectors ───────────────────────────────────

  get loading$(): Observable<boolean> {
    return this.state$.pipe(map((s) => s.isLoading), distinctUntilChanged());
  }

  get error$(): Observable<any> {
    return this.state$.pipe(map((s) => s.error), distinctUntilChanged());
  }

  get data$(): Observable<TData | null> {
    return this.state$.pipe(map((s) => s.data), distinctUntilChanged());
  }

  // ─── Execute ──────────────────────────────────────────────

  mutate(variables: TVariables): Observable<TData> {
    return new Observable<TData>((subscriber) => {
      // beforeMutate hooks
      this.plugins.forEach((p) => p.beforeMutate?.(variables, this.ctx));

      this.state$.next({ ...this.state$.value, status: 'loading', isLoading: true, error: null });

      this.activeSub?.unsubscribe();
      this.activeSub = this.config.mutationFn(variables).subscribe({
        next: (data) => {
          this.state$.next({ data, error: null, status: 'success', isLoading: false });

          // afterMutate hooks
          this.plugins.forEach((p) => p.afterMutate?.(data, variables, this.ctx));

          // Config callbacks
          this.config.onSuccess?.(data, variables);
          this.config.onSettled?.(data, null, variables);

          // Auto-invalidate related queries
          if (this.config.invalidateKeys?.length && this.onInvalidate) {
            this.onInvalidate(this.config.invalidateKeys);
          }

          // onSettled hooks
          this.plugins.forEach((p) => p.onSettled?.(this.ctx));

          subscriber.next(data);
          subscriber.complete();
        },
        error: (err) => {
          this.state$.next({ ...this.state$.value, error: err, status: 'error', isLoading: false });

          // onError hooks
          this.plugins.forEach((p) => p.onError?.(err, variables, this.ctx));

          // Config callbacks
          this.config.onError?.(err, variables);
          this.config.onSettled?.(null, err, variables);

          // onSettled hooks
          this.plugins.forEach((p) => p.onSettled?.(this.ctx));

          subscriber.error(err);
        },
      });
    });
  }

  /** Reset to idle state */
  reset(): void {
    this.activeSub?.unsubscribe();
    this.state$.next(initialMutationState());
  }
}
