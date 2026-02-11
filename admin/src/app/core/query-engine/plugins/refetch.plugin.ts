import { QueryPlugin } from '../types';

export interface RefetchOptions {
  /** Refetch when the browser tab regains focus. Default: true */
  onWindowFocus?: boolean;
  /** Refetch when the network comes back online. Default: true */
  onReconnect?: boolean;
}

/**
 * Keeps data fresh by refetching on window-focus and/or network-reconnect.
 */
export function refetchPlugin<T>(options: RefetchOptions = {}): QueryPlugin<T> {
  const { onWindowFocus = true, onReconnect = true } = options;

  return {
    name: 'refetch',

    onInit(ctx) {
      if (typeof window === 'undefined') return;

      if (onWindowFocus) {
        const handler = () => {
          const { isStale, status } = ctx.getState();
          if (isStale || status === 'error') ctx.refetch();
        };
        window.addEventListener('focus', handler);
        ctx.meta['refetch:focusHandler'] = handler;
      }

      if (onReconnect) {
        const handler = () => ctx.refetch();
        window.addEventListener('online', handler);
        ctx.meta['refetch:onlineHandler'] = handler;
      }
    },

    onDestroy(ctx) {
      if (typeof window === 'undefined') return;
      const focus = ctx.meta['refetch:focusHandler'];
      if (focus) window.removeEventListener('focus', focus);
      const online = ctx.meta['refetch:onlineHandler'];
      if (online) window.removeEventListener('online', online);
    },
  };
}
