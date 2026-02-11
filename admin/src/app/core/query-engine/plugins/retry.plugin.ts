import { QueryPlugin } from '../types';

export interface RetryOptions {
  /** Max number of retries. Default: 3 */
  maxRetries?: number;
  /** Base delay in ms. Default: 1 000 */
  delay?: number;
  /** Use exponential back-off (delay × 2^attempt). Default: true */
  backoff?: boolean;
  /** Only retry when predicate returns true */
  retryIf?: (error: any) => boolean;
}

/**
 * Automatically retries failed fetches with configurable back-off.
 */
export function retryPlugin<T>(options: RetryOptions = {}): QueryPlugin<T> {
  const { maxRetries = 3, delay = 1000, backoff = true, retryIf } = options;

  return {
    name: 'retry',

    onError(error, ctx) {
      const { retryCount } = ctx.getState();
      if (retryCount >= maxRetries) return;
      if (retryIf && !retryIf(error)) return;

      const nextRetry = retryCount + 1;
      const retryDelay = backoff ? delay * Math.pow(2, retryCount) : delay;
      ctx.setState({ retryCount: nextRetry });
      setTimeout(() => ctx.refetch(), retryDelay);
    },

    onInvalidate(ctx) {
      // Reset counter so a manual invalidate gets full retry budget
      ctx.setState({ retryCount: 0 });
    },
  };
}
