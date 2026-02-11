import { QueryPlugin } from '../types';

export interface SyncOptions {
  /** BroadcastChannel name. Default: 'query-engine-sync' */
  channel?: string;
}

/**
 * Cross-tab cache synchronisation via BroadcastChannel.
 *
 * When a query is invalidated in one tab, all other tabs
 * sharing the same channel automatically refetch.
 */
export function syncPlugin<T>(options: SyncOptions = {}): QueryPlugin<T> {
  const { channel = 'query-engine-sync' } = options;

  return {
    name: 'sync',

    onInit(ctx) {
      if (typeof BroadcastChannel === 'undefined') return;

      const bc = new BroadcastChannel(channel);
      bc.onmessage = (event: MessageEvent) => {
        const msg = event.data;
        if (msg?.key === ctx.key && msg?.type === 'invalidate') {
          ctx.refetch();
        }
      };
      ctx.meta['sync:channel'] = bc;
    },

    onInvalidate(ctx) {
      const bc = ctx.meta['sync:channel'] as BroadcastChannel | undefined;
      bc?.postMessage({ key: ctx.key, type: 'invalidate' });
    },

    onDestroy(ctx) {
      const bc = ctx.meta['sync:channel'] as BroadcastChannel | undefined;
      bc?.close();
    },
  };
}
