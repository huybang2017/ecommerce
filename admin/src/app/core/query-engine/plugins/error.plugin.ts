import { QueryPlugin } from '../types';

export interface ErrorOptions {
  /** Custom error handler — receives the raw error and the query key */
  handler?: (error: any, key: string) => void;
  /** Log errors to the console. Default: false */
  notify?: boolean;
}

/**
 * Centralised error handling / notification for queries.
 */
export function errorPlugin<T>(options: ErrorOptions = {}): QueryPlugin<T> {
  const { handler, notify = false } = options;

  return {
    name: 'error',

    onError(error, ctx) {
      if (handler) {
        handler(error, ctx.key);
        return;
      }
      if (notify) {
        console.error(`[QueryEngine] Error in "${ctx.key}":`, error);
      }
    },
  };
}
