import { Observable } from 'rxjs';
import { MutationConfig, MutationPlugin } from '../types';
import { MutationInstance } from '../core/mutation-instance';

/**
 * Fluent builder for composing a MutationInstance.
 *
 * ```ts
 * const mutation = MutationBuilder.for<Product, ProductCreate>()
 *   .mutateWith(p => api.create(p))
 *   .invalidates('products')
 *   .onSuccess(data => toast.success('Created'))
 *   .build(invalidateCb);
 * ```
 */
export class MutationBuilder<TData, TVariables> {
  private config: Partial<MutationConfig<TData, TVariables>> = { plugins: [] };

  static for<TData, TVariables>(): MutationBuilder<TData, TVariables> {
    return new MutationBuilder<TData, TVariables>();
  }

  mutateWith(fn: (variables: TVariables) => Observable<TData>): this {
    this.config.mutationFn = fn;
    return this;
  }

  use(...plugins: MutationPlugin<TData, TVariables>[]): this {
    this.config.plugins!.push(...plugins);
    return this;
  }

  onSuccess(cb: (data: TData, variables: TVariables) => void): this {
    this.config.onSuccess = cb;
    return this;
  }

  onError(cb: (error: any, variables: TVariables) => void): this {
    this.config.onError = cb;
    return this;
  }

  onSettled(cb: (data: TData | null, error: any | null, variables: TVariables) => void): this {
    this.config.onSettled = cb;
    return this;
  }

  invalidates(...keys: string[]): this {
    this.config.invalidateKeys = keys;
    return this;
  }

  /** Build a standalone MutationInstance */
  build(invalidateCallback?: (keys: string[]) => void): MutationInstance<TData, TVariables> {
    if (!this.config.mutationFn) throw new Error('[MutationBuilder] mutationFn is required');
    return new MutationInstance(this.config as MutationConfig<TData, TVariables>, invalidateCallback);
  }

  /** Export as plain config — pass to QueryEngine.mutation() */
  toConfig(): MutationConfig<TData, TVariables> {
    if (!this.config.mutationFn) throw new Error('[MutationBuilder] mutationFn is required');
    return this.config as MutationConfig<TData, TVariables>;
  }
}
