package scaffold

// syncReactTS emits packages/sync/src/react.tsx: the provider and hooks that
// let a screen read and write a resource without knowing whether the network
// is there.
func syncReactTS() string {
	return `import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import type { ReactNode } from "react";

import { SyncEngine } from "./engine";
import type { OutboxEntry, SyncStatus } from "./types";

const SyncContext = createContext<SyncEngine | null>(null);

export interface SyncProviderProps {
  engine: SyncEngine;
  children: ReactNode;
  /** Start the background loop on mount. Defaults to true. */
  autoSync?: boolean;
}

/**
 * Opens the engine, starts the loop, and puts it on context.
 *
 * The engine is constructed by the caller rather than here, because it needs
 * a storage adapter and the right one is a property of the platform:
 * IndexedDB in a browser, SQLite on a phone.
 */
export function SyncProvider({
  engine,
  children,
  autoSync = true,
}: SyncProviderProps) {
  const [ready, setReady] = useState(false);

  useEffect(() => {
    let cancelled = false;
    void engine.open().then(() => {
      if (cancelled) return;
      setReady(true);
      if (autoSync) engine.startAutoSync();
    });
    return () => {
      cancelled = true;
      engine.stopAutoSync();
    };
  }, [engine, autoSync]);

  // Children render before the mirror is open. Every read hook returns its
  // loading state until then, so the alternative is a blank screen while
  // IndexedDB opens rather than a list that fills in.
  void ready;

  return (
    <SyncContext.Provider value={engine}>{children}</SyncContext.Provider>
  );
}

export function useSyncEngine(): SyncEngine {
  const engine = useContext(SyncContext);
  if (!engine) {
    throw new Error("useSyncEngine must be used inside a <SyncProvider>");
  }
  return engine;
}

/**
 * The status the badge in your chrome renders: synced, syncing, offline, or
 * conflict, with the pending count.
 */
export function useSyncStatus(): SyncStatus & {
  syncNow: () => Promise<void>;
  setForceOffline: (value: boolean) => Promise<void>;
} {
  const engine = useSyncEngine();
  const [status, setStatus] = useState<SyncStatus>(() => engine.status());

  useEffect(() => engine.subscribe(setStatus), [engine]);

  const syncNow = useCallback(async () => {
    await engine.sync();
  }, [engine]);

  const setForceOffline = useCallback(
    async (value: boolean) => {
      await engine.setForceOffline(value);
    },
    [engine],
  );

  return { ...status, syncNow, setForceOffline };
}

export interface OfflineResource<T> {
  data: T[];
  loading: boolean;
  error: Error | null;
  create: (values: Partial<T>) => Promise<string>;
  update: (id: string, values: Partial<T>) => Promise<void>;
  remove: (id: string) => Promise<void>;
  refresh: () => Promise<void>;
}

/**
 * One resource, read from the mirror and written through the outbox.
 *
 * This is the hook that makes "works offline" a property of a resource rather
 * than of a client. The screen calls create and gets an id back immediately,
 * whether the row reached the server or is sitting in the outbox waiting for
 * a network.
 *
 * The list re-reads on every status change. The engine emits one after each
 * local write and after each sync, which covers both the row you just added
 * and the rows somebody else added while you were away.
 */
export function useOfflineResource<T extends { id: string }>(
  model: string,
): OfflineResource<T> {
  const engine = useSyncEngine();
  const [data, setData] = useState<T[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const mounted = useRef(true);

  useEffect(() => {
    mounted.current = true;
    return () => {
      mounted.current = false;
    };
  }, []);

  const refresh = useCallback(async () => {
    try {
      const rows = await engine.list(model);
      if (mounted.current) {
        setData(rows as unknown as T[]);
        setError(null);
      }
    } catch (caught) {
      if (mounted.current) {
        setError(caught instanceof Error ? caught : new Error(String(caught)));
      }
    } finally {
      if (mounted.current) setLoading(false);
    }
  }, [engine, model]);

  useEffect(() => {
    return engine.subscribe(() => {
      void refresh();
    });
  }, [engine, refresh]);

  const create = useCallback(
    (values: Partial<T>) =>
      engine.create(model, values as Record<string, unknown>),
    [engine, model],
  );

  const update = useCallback(
    (id: string, values: Partial<T>) =>
      engine.update(model, id, values as Record<string, unknown>),
    [engine, model],
  );

  const remove = useCallback(
    (id: string) => engine.remove(model, id),
    [engine, model],
  );

  return useMemo(
    () => ({ data, loading, error, create, update, remove, refresh }),
    [data, loading, error, create, update, remove, refresh],
  );
}

/** One row from the mirror. */
export function useOfflineRecord<T extends { id: string }>(
  model: string,
  id: string | null | undefined,
): { data: T | null; loading: boolean } {
  const engine = useSyncEngine();
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) {
      setData(null);
      setLoading(false);
      return;
    }
    return engine.subscribe(() => {
      void engine.get(model, id).then((row) => {
        setData((row as unknown as T) ?? null);
        setLoading(false);
      });
    });
  }, [engine, model, id]);

  return { data, loading };
}

/**
 * The conflicts waiting for a decision, and the two ways to end one.
 *
 * resolve takes the merged row the user assembled from both sides. revert
 * throws the local change away and keeps the server's.
 */
export function useSyncConflicts(): {
  conflicts: OutboxEntry[];
  resolve: (
    model: string,
    id: string,
    merged: Record<string, unknown>,
    serverVersion: number,
  ) => Promise<void>;
  revert: (model: string, id: string) => Promise<void>;
} {
  const engine = useSyncEngine();
  const [conflicts, setConflicts] = useState<OutboxEntry[]>([]);

  useEffect(() => {
    return engine.subscribe(() => {
      void engine.listConflicts().then(setConflicts);
    });
  }, [engine]);

  const resolve = useCallback(
    (
      model: string,
      id: string,
      merged: Record<string, unknown>,
      serverVersion: number,
    ) => engine.resolveConflict(model, id, merged, serverVersion),
    [engine],
  );

  const revert = useCallback(
    (model: string, id: string) => engine.revert(model, id),
    [engine],
  );

  return { conflicts, resolve, revert };
}
`
}
