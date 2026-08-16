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

/** What the UI shows: a badge, a pending count, a conflict list. */
export type SyncState = "synced" | "syncing" | "offline" | "conflict";

export interface SyncStatus {
  state: SyncState;
  online: boolean;
  syncing: boolean;
  forceOffline: boolean;
  pending: number;
  conflicts: number;
  lastSyncedAt: number | null;
  lastError: string | null;
}

export interface SyncResult {
  pulled: number;
  pushed: number;
  conflicts: number;
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
}
`
}
