package scaffold

// The three storage adapters. The engine holds no storage-specific code, so
// adding a fourth is a file rather than a fork.

// syncMemoryAdapterTS emits packages/sync/src/adapters/memory.ts.
func syncMemoryAdapterTS() string {
	return `import type { LocalRecord, OutboxEntry, StorageAdapter } from "../types";

/**
 * An in-memory mirror.
 *
 * For tests, and for server rendering, where IndexedDB does not exist and a
 * component that reaches for it throws during the render rather than
 * degrading. Nothing survives a reload, which is the whole point: it is the
 * adapter you use when persistence is not what you are testing.
 */
export class MemoryAdapter implements StorageAdapter {
  private records = new Map<string, LocalRecord>();
  private outbox = new Map<string, OutboxEntry>();
  private meta = new Map<string, string>();

  async open(): Promise<void> {}
  async close(): Promise<void> {}

  async putRecord(record: LocalRecord): Promise<void> {
    this.records.set(key(record.model, record.id), { ...record });
  }

  async getRecord(model: string, id: string): Promise<LocalRecord | null> {
    const found = this.records.get(key(model, id));
    return found ? { ...found } : null;
  }

  async listRecords(model: string): Promise<LocalRecord[]> {
    const out: LocalRecord[] = [];
    for (const record of this.records.values()) {
      if (record.model === model) out.push({ ...record });
    }
    return out;
  }

  async deleteRecord(model: string, id: string): Promise<void> {
    this.records.delete(key(model, id));
  }

  async clearModel(model: string): Promise<void> {
    for (const [k, record] of this.records) {
      if (record.model === model) this.records.delete(k);
    }
  }

  async putOutbox(entry: OutboxEntry): Promise<void> {
    this.outbox.set(key(entry.model, entry.entityId), { ...entry });
  }

  async getOutbox(model: string, id: string): Promise<OutboxEntry | null> {
    const found = this.outbox.get(key(model, id));
    return found ? { ...found } : null;
  }

  async listOutbox(): Promise<OutboxEntry[]> {
    return Array.from(this.outbox.values()).map((e) => ({ ...e }));
  }

  async deleteOutbox(model: string, id: string): Promise<void> {
    this.outbox.delete(key(model, id));
  }

  async getMeta(k: string): Promise<string> {
    return this.meta.get(k) ?? "";
  }

  async setMeta(k: string, value: string): Promise<void> {
    this.meta.set(k, value);
  }
}

function key(model: string, id: string): string {
  return model + "::" + id;
}
`
}

// syncIndexedAdapterTS emits packages/sync/src/adapters/indexed.ts.
func syncIndexedAdapterTS() string {
	return `import type { LocalRecord, OutboxEntry, StorageAdapter } from "../types";

const DB_VERSION = 1;
const RECORDS = "records";
const OUTBOX = "outbox";
const META = "meta";

/**
 * IndexedDB, for the browser and for a PWA build.
 *
 * Records and outbox entries both key on [model, id], which is what makes the
 * squash rule a put rather than a scan. Records also carry a "model" index so
 * listing one model does not walk every row of every other one.
 *
 * localStorage was the alternative and is the wrong shape: synchronous, string
 * only, and a few megabytes at best. A mirror of a real table outgrows it.
 */
export class IndexedDBAdapter implements StorageAdapter {
  private db: IDBDatabase | null = null;

  constructor(private readonly name = "grit-sync") {}

  async open(): Promise<void> {
    if (this.db) return;
    this.db = await new Promise<IDBDatabase>((resolve, reject) => {
      const request = indexedDB.open(this.name, DB_VERSION);
      request.onupgradeneeded = () => {
        const db = request.result;
        if (!db.objectStoreNames.contains(RECORDS)) {
          const store = db.createObjectStore(RECORDS, {
            keyPath: ["model", "id"],
          });
          store.createIndex("model", "model", { unique: false });
        }
        if (!db.objectStoreNames.contains(OUTBOX)) {
          const store = db.createObjectStore(OUTBOX, {
            keyPath: ["model", "entityId"],
          });
          store.createIndex("createdAt", "createdAt", { unique: false });
        }
        if (!db.objectStoreNames.contains(META)) {
          db.createObjectStore(META, { keyPath: "key" });
        }
      };
      request.onsuccess = () => resolve(request.result);
      request.onerror = () => reject(request.error);
    });
  }

  async close(): Promise<void> {
    this.db?.close();
    this.db = null;
  }

  async putRecord(record: LocalRecord): Promise<void> {
    await this.write(RECORDS, (store) => store.put(record));
  }

  async getRecord(model: string, id: string): Promise<LocalRecord | null> {
    const found = await this.read<LocalRecord>(RECORDS, (store) =>
      store.get([model, id]),
    );
    return found ?? null;
  }

  async listRecords(model: string): Promise<LocalRecord[]> {
    const found = await this.read<LocalRecord[]>(RECORDS, (store) =>
      store.index("model").getAll(model),
    );
    return found ?? [];
  }

  async deleteRecord(model: string, id: string): Promise<void> {
    await this.write(RECORDS, (store) => store.delete([model, id]));
  }

  async clearModel(model: string): Promise<void> {
    const records = await this.listRecords(model);
    for (const record of records) {
      await this.deleteRecord(model, record.id);
    }
  }

  async putOutbox(entry: OutboxEntry): Promise<void> {
    await this.write(OUTBOX, (store) => store.put(entry));
  }

  async getOutbox(model: string, id: string): Promise<OutboxEntry | null> {
    const found = await this.read<OutboxEntry>(OUTBOX, (store) =>
      store.get([model, id]),
    );
    return found ?? null;
  }

  async listOutbox(): Promise<OutboxEntry[]> {
    const found = await this.read<OutboxEntry[]>(OUTBOX, (store) =>
      store.getAll(),
    );
    return found ?? [];
  }

  async deleteOutbox(model: string, id: string): Promise<void> {
    await this.write(OUTBOX, (store) => store.delete([model, id]));
  }

  async getMeta(key: string): Promise<string> {
    const found = await this.read<{ key: string; value: string }>(
      META,
      (store) => store.get(key),
    );
    return found ? found.value : "";
  }

  async setMeta(key: string, value: string): Promise<void> {
    await this.write(META, (store) => store.put({ key, value }));
  }

  private store(name: string, mode: IDBTransactionMode): IDBObjectStore {
    if (!this.db) throw new Error("sync: adapter not open");
    return this.db.transaction(name, mode).objectStore(name);
  }

  private read<T>(
    name: string,
    run: (store: IDBObjectStore) => IDBRequest,
  ): Promise<T | undefined> {
    return new Promise((resolve, reject) => {
      const request = run(this.store(name, "readonly"));
      request.onsuccess = () => resolve(request.result as T);
      request.onerror = () => reject(request.error);
    });
  }

  private write(
    name: string,
    run: (store: IDBObjectStore) => IDBRequest,
  ): Promise<void> {
    return new Promise((resolve, reject) => {
      // Resolve on the transaction, not the request. A request that succeeds
      // inside a transaction that later aborts has not written anything, and
      // resolving early would report a durable write that is not.
      const store = this.store(name, "readwrite");
      const request = run(store);
      request.onerror = () => reject(request.error);
      store.transaction.oncomplete = () => resolve();
      store.transaction.onabort = () => reject(store.transaction.error);
    });
  }
}
`
}

// syncSQLiteAdapterTS emits packages/sync/src/adapters/sqlite.ts.
func syncSQLiteAdapterTS() string {
	return `import type { LocalRecord, OutboxEntry, StorageAdapter } from "../types";

/**
 * The slice of expo-sqlite this adapter uses.
 *
 * Declared structurally rather than imported, so packages/sync does not take
 * a hard dependency on expo-sqlite. A web build that never touches this file
 * should not have to install a native module to typecheck.
 */
export interface SQLiteDatabase {
  execAsync(source: string): Promise<void>;
  runAsync(source: string, params?: unknown[]): Promise<unknown>;
  getAllAsync<T>(source: string, params?: unknown[]): Promise<T[]>;
  getFirstAsync<T>(source: string, params?: unknown[]): Promise<T | null>;
  closeAsync(): Promise<void>;
}

interface RecordRow {
  model: string;
  id: string;
  data: string;
  version: number;
  updated_at: number;
  deleted: number;
}

interface OutboxRow {
  model: string;
  entity_id: string;
  op: string;
  data: string | null;
  version: number;
  created_at: number;
  has_conflict: number;
  server_data: string | null;
  server_version: number;
  conflict_message: string;
}

/**
 * SQLite, for the Expo app.
 *
 * The same three tables the desktop engine keeps in GORM, so a row means the
 * same thing on a phone as it does on a laptop. Pass the database in rather
 * than opening it here: the caller already has one, and Expo wants the open
 * to happen where the app can show a splash screen while it happens.
 *
 *   import * as SQLite from "expo-sqlite";
 *   const db = await SQLite.openDatabaseAsync("grit-sync.db");
 *   const adapter = new SQLiteAdapter(db);
 */
export class SQLiteAdapter implements StorageAdapter {
  constructor(private readonly db: SQLiteDatabase) {}

  async open(): Promise<void> {
    // WAL is not a tuning knob here. The default journal blocks readers during
    // a write, and a sync writing a few hundred pulled rows would freeze every
    // list on screen for the duration.
    await this.db.execAsync(
      "PRAGMA journal_mode = WAL;" +
        "CREATE TABLE IF NOT EXISTS sync_records (" +
        "  model TEXT NOT NULL," +
        "  id TEXT NOT NULL," +
        "  data TEXT NOT NULL," +
        "  version INTEGER NOT NULL DEFAULT 0," +
        "  updated_at INTEGER NOT NULL DEFAULT 0," +
        "  deleted INTEGER NOT NULL DEFAULT 0," +
        "  PRIMARY KEY (model, id)" +
        ");" +
        "CREATE INDEX IF NOT EXISTS idx_sync_records_model" +
        "  ON sync_records (model);" +
        "CREATE TABLE IF NOT EXISTS sync_outbox (" +
        "  model TEXT NOT NULL," +
        "  entity_id TEXT NOT NULL," +
        "  op TEXT NOT NULL," +
        "  data TEXT," +
        "  version INTEGER NOT NULL DEFAULT 0," +
        "  created_at INTEGER NOT NULL DEFAULT 0," +
        "  has_conflict INTEGER NOT NULL DEFAULT 0," +
        "  server_data TEXT," +
        "  server_version INTEGER NOT NULL DEFAULT 0," +
        "  conflict_message TEXT NOT NULL DEFAULT ''," +
        "  PRIMARY KEY (model, entity_id)" +
        ");" +
        "CREATE TABLE IF NOT EXISTS sync_meta (" +
        "  key TEXT PRIMARY KEY," +
        "  value TEXT NOT NULL" +
        ");",
    );
  }

  async close(): Promise<void> {
    await this.db.closeAsync();
  }

  async putRecord(record: LocalRecord): Promise<void> {
    await this.db.runAsync(
      "INSERT INTO sync_records (model, id, data, version, updated_at, deleted)" +
        " VALUES (?, ?, ?, ?, ?, ?)" +
        " ON CONFLICT(model, id) DO UPDATE SET" +
        "   data = excluded.data," +
        "   version = excluded.version," +
        "   updated_at = excluded.updated_at," +
        "   deleted = excluded.deleted",
      [
        record.model,
        record.id,
        JSON.stringify(record.data),
        record.version,
        record.updatedAt,
        record.deleted ? 1 : 0,
      ],
    );
  }

  async getRecord(model: string, id: string): Promise<LocalRecord | null> {
    const row = await this.db.getFirstAsync<RecordRow>(
      "SELECT * FROM sync_records WHERE model = ? AND id = ?",
      [model, id],
    );
    return row ? toRecord(row) : null;
  }

  async listRecords(model: string): Promise<LocalRecord[]> {
    const rows = await this.db.getAllAsync<RecordRow>(
      "SELECT * FROM sync_records WHERE model = ?",
      [model],
    );
    return rows.map(toRecord);
  }

  async deleteRecord(model: string, id: string): Promise<void> {
    await this.db.runAsync(
      "DELETE FROM sync_records WHERE model = ? AND id = ?",
      [model, id],
    );
  }

  async clearModel(model: string): Promise<void> {
    await this.db.runAsync("DELETE FROM sync_records WHERE model = ?", [model]);
  }

  async putOutbox(entry: OutboxEntry): Promise<void> {
    await this.db.runAsync(
      "INSERT INTO sync_outbox (model, entity_id, op, data, version," +
        " created_at, has_conflict, server_data, server_version, conflict_message)" +
        " VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)" +
        " ON CONFLICT(model, entity_id) DO UPDATE SET" +
        "   op = excluded.op," +
        "   data = excluded.data," +
        "   version = excluded.version," +
        "   has_conflict = excluded.has_conflict," +
        "   server_data = excluded.server_data," +
        "   server_version = excluded.server_version," +
        "   conflict_message = excluded.conflict_message",
      [
        entry.model,
        entry.entityId,
        entry.op,
        entry.data === null ? null : JSON.stringify(entry.data),
        entry.version,
        entry.createdAt,
        entry.hasConflict ? 1 : 0,
        entry.serverData === null ? null : JSON.stringify(entry.serverData),
        entry.serverVersion,
        entry.conflictMessage,
      ],
    );
  }

  async getOutbox(model: string, id: string): Promise<OutboxEntry | null> {
    const row = await this.db.getFirstAsync<OutboxRow>(
      "SELECT * FROM sync_outbox WHERE model = ? AND entity_id = ?",
      [model, id],
    );
    return row ? toOutbox(row) : null;
  }

  async listOutbox(): Promise<OutboxEntry[]> {
    const rows = await this.db.getAllAsync<OutboxRow>(
      "SELECT * FROM sync_outbox ORDER BY created_at ASC",
    );
    return rows.map(toOutbox);
  }

  async deleteOutbox(model: string, id: string): Promise<void> {
    await this.db.runAsync(
      "DELETE FROM sync_outbox WHERE model = ? AND entity_id = ?",
      [model, id],
    );
  }

  async getMeta(key: string): Promise<string> {
    const row = await this.db.getFirstAsync<{ value: string }>(
      "SELECT value FROM sync_meta WHERE key = ?",
      [key],
    );
    return row ? row.value : "";
  }

  async setMeta(key: string, value: string): Promise<void> {
    await this.db.runAsync(
      "INSERT INTO sync_meta (key, value) VALUES (?, ?)" +
        " ON CONFLICT(key) DO UPDATE SET value = excluded.value",
      [key, value],
    );
  }
}

function toRecord(row: RecordRow): LocalRecord {
  return {
    model: row.model,
    id: row.id,
    data: parseJSON(row.data) ?? {},
    version: row.version,
    updatedAt: row.updated_at,
    deleted: row.deleted === 1,
  };
}

function toOutbox(row: OutboxRow): OutboxEntry {
  return {
    model: row.model,
    entityId: row.entity_id,
    op: row.op as OutboxEntry["op"],
    data: row.data === null ? null : parseJSON(row.data),
    version: row.version,
    createdAt: row.created_at,
    hasConflict: row.has_conflict === 1,
    serverData: row.server_data === null ? null : parseJSON(row.server_data),
    serverVersion: row.server_version,
    conflictMessage: row.conflict_message,
  };
}

/**
 * Corrupt JSON in the mirror is recoverable: the next pull overwrites the row.
 * Throwing here would take down the whole list instead.
 */
function parseJSON(raw: string): Record<string, unknown> | null {
  try {
    return JSON.parse(raw) as Record<string, unknown>;
  } catch {
    return null;
  }
}
`
}
