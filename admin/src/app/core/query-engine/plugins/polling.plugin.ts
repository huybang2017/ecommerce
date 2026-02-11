import { QueryPlugin } from '../types';

export interface PollingOptions {
  /** Polling interval in milliseconds */
  interval: number;
  /** Pause polling while the document tab is hidden. Default: true */
  pauseOnHidden?: boolean;
}

/**
 * Periodic background refetch on a fixed interval.
 */
export function pollingPlugin<T>(options: PollingOptions): QueryPlugin<T> {
  const { interval, pauseOnHidden = true } = options;

  return {
    name: 'polling',

    onInit(ctx) {
      const timer = setInterval(() => {
        if (pauseOnHidden && typeof document !== 'undefined' && document.hidden) return;
        ctx.refetch();
      }, interval);
      ctx.meta['polling:timer'] = timer;
    },

    onDestroy(ctx) {
      const timer = ctx.meta['polling:timer'];
      if (timer) clearInterval(timer);
    },
  };
}
