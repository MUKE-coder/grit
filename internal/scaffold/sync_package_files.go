package scaffold

import (
	"fmt"
	"path/filepath"
)

// writeSyncPackageFiles emits packages/sync: the offline-first client that
// apps/web and apps/expo share.
//
// The desktop app has had one of these since v3.60, written in Go and living
// inside apps/desktop. The wire protocol it speaks (/api/sync/pull with a
// cursor, /api/sync/push with per-change version checks) is not
// desktop-specific and neither is the local outbox, so the second
// implementation is a port rather than a new design: same record shape, same
// squash rules, same conflict semantics, in TypeScript, over a storage
// interface so one engine serves IndexedDB in a browser and SQLite on a
// phone.
func writeSyncPackageFiles(root string, opts Options) error {
	pkg := filepath.Join(root, "packages", "sync")
	files := map[string]string{
		filepath.Join(pkg, "package.json"):                  syncPackageJSON(opts),
		filepath.Join(pkg, "tsconfig.json"):                 syncTSConfig(),
		filepath.Join(pkg, "src", "index.ts"):               syncIndexTS(),
		filepath.Join(pkg, "src", "types.ts"):               syncTypesTS(),
		filepath.Join(pkg, "src", "engine.ts"):              syncEngineTS(),
		filepath.Join(pkg, "src", "adapters", "memory.ts"):  syncMemoryAdapterTS(),
		filepath.Join(pkg, "src", "adapters", "indexed.ts"): syncIndexedAdapterTS(),
		filepath.Join(pkg, "src", "adapters", "sqlite.ts"):  syncSQLiteAdapterTS(),
		filepath.Join(pkg, "src", "react.tsx"):              syncReactTS(),
	}
	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

func syncPackageJSON(opts Options) string {
	return `{
  "name": "@` + opts.ProjectName + `/sync",
  "version": "0.0.0",
  "private": true,
  "main": "./src/index.ts",
  "types": "./src/index.ts",
  "exports": {
    ".": "./src/index.ts",
    "./react": "./src/react.tsx",
    "./adapters/indexed": "./src/adapters/indexed.ts",
    "./adapters/sqlite": "./src/adapters/sqlite.ts",
    "./adapters/memory": "./src/adapters/memory.ts"
  },
  "peerDependencies": {
    "react": ">=18"
  },
  "peerDependenciesMeta": {
    "react": { "optional": true }
  }
}
`
}

func syncTSConfig() string {
	return `{
  "compilerOptions": {
    "target": "ES2020",
    "lib": ["ES2020", "DOM"],
    "module": "ESNext",
    "moduleResolution": "Bundler",
    "jsx": "react-jsx",
    "strict": true,
    "declaration": true,
    "skipLibCheck": true,
    "noEmit": true,
    "esModuleInterop": true,
    "forceConsistentCasingInFileNames": true
  },
  "include": ["src/**/*"]
}
`
}

func syncIndexTS() string {
	return `export * from "./types";
export { SyncEngine } from "./engine";
export { MemoryAdapter } from "./adapters/memory";
`
}

func syncTypesTS() string {
	return `// The offline-first contract, shared by every Grit client.
//
// These types mirror the Go structs the API already serves: PushChange,
// PushResult and the pull envelope in apps/api/internal/handlers/sync.go.
// If you change one side, change the other.

/** A local change waiting to reach the server. */
export type SyncOp = "create" | "update" | "delete";

/** One entry in a /api/sync/push batch. */
export interface PushChange {
  op: SyncOp;
  model: string;
  id: string;
  /** The version the client believes the server holds. */
  version: number;
  data?: Record<string, unknown> | null;
}

/**
 * The per-change result, returned in the same order as the batch.
 *
 * Field names are snake_case because they come off the wire that way. On
 * VERSION_CONFLICT, server_version and server_data carry the current server
 * state so the client can build a merge UI without a second round trip.
 */
export interface PushResult {
  ok: boolean;
  code?: string;
  message?: string;
  server_version?: number;
  server_data?: Record<string, unknown>;
  new_version?: number;
}

/** The /api/sync/pull envelope. */
export interface PullResponse {
  data: Array<Record<string, unknown>>;
  cursor: string;
  count: number;
}

/**
 * One row in the local mirror.
 *
 * deleted is a tombstone rather than a removal: a pulled server delete has to
 * hide the row from reads while still letting a later re-create upsert over
 * it.
 */
export interface LocalRecord {
  model: string;
  id: string;
  data: Record<string, unknown>;
  version: number;
  updatedAt: number;
  deleted: boolean;
}

/**
 * One pending local change. At most one per (model, id): a second edit to the
 * same row squashes into the existing entry rather than queueing behind it,
 * which is what keeps the outbox proportional to the number of touched rows
 * instead of the number of keystrokes.
 */
export interface OutboxEntry {
  model: string;
  entityId: string;
  op: SyncOp;
  data: Record<string, unknown> | null;
  version: number;
  createdAt: number;
  hasConflict: boolean;
  serverData: Record<string, unknown> | null;
  serverVersion: number;
  conflictMessage: string;
}

/**
 * How one model behaves offline, as declared on the server.
 *
 * Fetched from GET /api/sync/policy rather than compiled in, because a copy
 * is a second thing to keep in sync and an offline client running last
 * month's conflict rules is the exact silent failure this prevents.
 */
export interface SyncPolicy {
  mode: "offline_first" | "online_only";
  conflict: "manual" | "server_wins" | "client_wins";
  fields?: string[];
  local_only?: string[];
  max_offline_age_seconds?: number;
}

/**
 * What the UI shows.
 *
 * "stale" is its own state and not a flavour of offline: an app that has been
 * offline for ten seconds and one that has been offline past its declared
 * max_offline_age are telling the user different things. The second is
 * showing data nobody should act on.
 */
export type SyncState =
  | "synced"
  | "syncing"
  | "offline"
  | "conflict"
  | "stale";

export interface SyncStatus {
  state: SyncState;
  online: boolean;
  syncing: boolean;
  forceOffline: boolean;
  pending: number;
  conflicts: number;
  lastSyncedAt: number | null;
  lastError: string | null;
  /** Past the strictest max_offline_age of any mirrored model. */
  stale: boolean;
}

/**
 * The numbers that tell you an offline app is in trouble.
 *
 * Offline products fail quietly. An outbox that stopped draining three days
 * ago looks exactly like an outbox with nothing in it, and the only
 * difference visible from inside the app is these figures.
 */
export interface SyncHealth {
  pending: number;
  conflicts: number;
  /** Seconds since the oldest queued change was made. Zero when empty. */
  oldestPendingAgeSeconds: number;
  /** Seconds since the last successful sync, or null if there has never been one. */
  lastSyncAgeSeconds: number | null;
  stale: boolean;
  maxOfflineAgeSeconds: number | null;
  lastError: string | null;
  /** Per model: how many rows are mirrored and how many are queued. */
  models: Array<{ model: string; mirrored: number; pending: number }>;
}

export interface SyncResult {
  pulled: number;
  pushed: number;
  conflicts: number;
  /** Changes the server discarded under a server_wins policy. */
  overridden: number;
}

/**
 * Where the mirror lives.
 *
 * Three implementations ship: IndexedDB for the browser, expo-sqlite for
 * React Native, and an in-memory one for tests and server rendering. The
 * engine holds no storage-specific code, so a fourth is a file rather than a
 * fork.
 */
export interface StorageAdapter {
  open(): Promise<void>;
  close(): Promise<void>;

  putRecord(record: LocalRecord): Promise<void>;
  getRecord(model: string, id: string): Promise<LocalRecord | null>;
  listRecords(model: string): Promise<LocalRecord[]>;
  deleteRecord(model: string, id: string): Promise<void>;
  clearModel(model: string): Promise<void>;

  putOutbox(entry: OutboxEntry): Promise<void>;
  getOutbox(model: string, id: string): Promise<OutboxEntry | null>;
  listOutbox(): Promise<OutboxEntry[]>;
  deleteOutbox(model: string, id: string): Promise<void>;

  getMeta(key: string): Promise<string>;
  setMeta(key: string, value: string): Promise<void>;
}

export interface SyncEngineOptions {
  /** Base API URL, e.g. "http://localhost:8080/api". */
  apiUrl: string;
  adapter: StorageAdapter;
  /** The models to mirror, in plural snake_case as the server registered them. */
  models: string[];
  /** Returns the current access token, or null when signed out. */
  getToken?: () => string | null | Promise<string | null>;
  /** How often to sync while online. Set 0 to disable. Defaults to 30s. */
  autoSyncIntervalMs?: number;
  /** Rows per pull page. The server caps this at 5000. */
  pullLimit?: number;
  fetchImpl?: typeof fetch;
  /**
   * Skip fetching policies from the server and use these instead. For tests,
   * and for a client that must work before it has ever reached the API.
   */
  policies?: Record<string, SyncPolicy>;
}
`
}
