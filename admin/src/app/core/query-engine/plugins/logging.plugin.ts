import { QueryPlugin, MutationPlugin } from '../types';

export interface LoggingOptions {
  /** Custom logger object. Default: console */
  logger?: {
    info: (...args: any[]) => void;
    warn: (...args: any[]) => void;
    error: (...args: any[]) => void;
  };
  /** Log line prefix. Default: [QueryEngine] */
  prefix?: string;
  /** Include data payloads in log output. Default: false */
  verbose?: boolean;
}

// ─── Query plugin ─────────────────────────────────────────────

export function loggingPlugin<T>(options: LoggingOptions = {}): QueryPlugin<T> {
  const { logger = console, prefix = '[QueryEngine]', verbose = false } = options;

  return {
    name: 'logging',

    onInit(ctx) {
      if (verbose) logger.info(`${prefix} Init query "${ctx.key}"`);
    },

    beforeFetch(ctx) {
      logger.info(`${prefix} Fetching "${ctx.key}"`);
    },

    afterFetch(data, ctx) {
      logger.info(`${prefix} Success "${ctx.key}"`, verbose ? data : '');
    },

    onError(error, ctx) {
      logger.error(`${prefix} Error "${ctx.key}"`, error);
    },

    onInvalidate(ctx) {
      if (verbose) logger.info(`${prefix} Invalidated "${ctx.key}"`);
    },

    onDestroy(ctx) {
      if (verbose) logger.info(`${prefix} Destroyed "${ctx.key}"`);
    },
  };
}

// ─── Mutation plugin ──────────────────────────────────────────

export function mutationLoggingPlugin<TData, TVariables>(
  options: LoggingOptions = {},
): MutationPlugin<TData, TVariables> {
  const { logger = console, prefix = '[QueryEngine]', verbose = false } = options;

  return {
    name: 'mutation-logging',

    beforeMutate(variables) {
      logger.info(`${prefix} Mutating…`, verbose ? variables : '');
    },

    afterMutate(data) {
      logger.info(`${prefix} Mutation success`, verbose ? data : '');
    },

    onError(error) {
      logger.error(`${prefix} Mutation error`, error);
    },
  };
}
