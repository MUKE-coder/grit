package scaffold

// syncEngineTS emits packages/sync/src/engine.ts.
//
// This is a port of apps/desktop/internal/sync/engine.go, not a second
// design. Same mirror, same outbox with the same squash rules, same
// version-conflict handling, over a storage interface so one engine serves
// IndexedDB in a browser and SQLite on a phone.
func syncEngineTS() string {
	return `import type {
  OutboxEntry,
  PullResponse,
  PushChange,
  PushResult,
  StorageAdapter,
  SyncEngineOptions,
  SyncOp,
  SyncResult,
  SyncStatus,
} from "./types";

const FORCE_OFFLINE_KEY = "force_offline";
const LAST_SYNCED_KEY = "last_synced_at";

/**
 * The offline-first engine.
 *
 * It owns three things: a mirror of server rows, an outbox of local changes
 * that have not reached the server yet, and a cursor per model so a pull asks
 * only for what changed. Reads come from the mirror, so a screen renders the
 * same way and at the same speed whether or not there is a network.
 */
export class SyncEngine {
  private readonly apiUrl: string;
  private readonly adapter: StorageAdapter;
  private readonly models: string[];
  private readonly getToken: () => string | null | Promise<string | null>;
  private readonly autoSyncIntervalMs: number;
  private readonly pullLimit: number;
  private readonly fetchImpl: typeof fetch;

  private listeners = new Set<(status: SyncStatus) => void>();
  private timer: ReturnType<typeof setInterval> | null = null;
  private opened = false;
  private syncing = false;
  private inFlight: Promise<SyncResult> | null = null;
  private online = true;
  private forceOffline = false;
  private pending = 0;
  private conflicts = 0;
  private lastSyncedAt: number | null = null;
  private lastError: string | null = null;

  constructor(options: SyncEngineOptions) {
    this.apiUrl = options.apiUrl.replace(/\/+$/, "");
    this.adapter = options.adapter;
    this.models = options.models;
    this.getToken = options.getToken ?? (() => null);
    this.autoSyncIntervalMs = options.autoSyncIntervalMs ?? 30000;
    this.pullLimit = options.pullLimit ?? 500;
    this.fetchImpl = options.fetchImpl ?? globalThis.fetch.bind(globalThis);
  }

  // --------------------------------------------------------------- lifecycle

  async open(): Promise<void> {
    if (this.opened) return;
    await this.adapter.open();
    this.opened = true;
    this.forceOffline = (await this.adapter.getMeta(FORCE_OFFLINE_KEY)) === "1";
    const last = await this.adapter.getMeta(LAST_SYNCED_KEY);
    this.lastSyncedAt = last ? Number(last) : null;
    await this.refreshCounts();
    this.emit();
  }

  async close(): Promise<void> {
    this.stopAutoSync();
    if (this.opened) {
      await this.adapter.close();
      this.opened = false;
    }
  }

  /**
   * Starts the background loop. Idempotent, because a component that mounts
   * twice must not leave two timers racing into the same outbox.
   */
  startAutoSync(): void {
    if (this.timer !== null || this.autoSyncIntervalMs <= 0) return;
    this.timer = setInterval(() => {
      void this.tick();
    }, this.autoSyncIntervalMs);
    void this.tick();
  }

  stopAutoSync(): void {
    if (this.timer !== null) {
      clearInterval(this.timer);
      this.timer = null;
    }
  }

  private async tick(): Promise<void> {
    if (this.forceOffline || this.syncing) return;
    if (!(await this.reachable())) return;
    try {
      await this.sync();
    } catch {
      // A background tick nobody asked for should not throw at the app. The
      // reason is on the status, and the next interval retries.
    }
  }

  // ------------------------------------------------------------------ status

  subscribe(listener: (status: SyncStatus) => void): () => void {
    this.listeners.add(listener);
    listener(this.status());
    return () => {
      this.listeners.delete(listener);
    };
  }

  status(): SyncStatus {
    let state: SyncStatus["state"] = "synced";
    if (this.conflicts > 0) state = "conflict";
    else if (this.syncing) state = "syncing";
    else if (this.forceOffline || !this.online) state = "offline";

    return {
      state,
      online: this.online && !this.forceOffline,
      syncing: this.syncing,
      forceOffline: this.forceOffline,
      pending: this.pending,
      conflicts: this.conflicts,
      lastSyncedAt: this.lastSyncedAt,
      lastError: this.lastError,
    };
  }

  private emit(): void {
    const snapshot = this.status();
    for (const listener of this.listeners) listener(snapshot);
  }

  private async refreshCounts(): Promise<void> {
    const entries = await this.adapter.listOutbox();
    this.pending = entries.length;
    this.conflicts = entries.filter((e) => e.hasConflict).length;
  }

  /** Work offline deliberately, even with a network available. */
  async setForceOffline(value: boolean): Promise<void> {
    this.forceOffline = value;
    await this.adapter.setMeta(FORCE_OFFLINE_KEY, value ? "1" : "0");
    this.emit();
  }

  async reachable(): Promise<boolean> {
    try {
      const response = await this.fetchImpl(this.apiUrl + "/health", {
        method: "GET",
      });
      // Anything short of a server error means we got there. A 401 still
      // answers the question being asked, which is whether the network is up.
      this.online = response.status < 500;
    } catch {
      this.online = false;
    }
    this.emit();
    return this.online;
  }

  // ------------------------------------------------------------- local reads

  async list(model: string): Promise<Array<Record<string, unknown>>> {
    const records = await this.adapter.listRecords(model);
    return records.filter((r) => !r.deleted).map((r) => r.data);
  }

  async get(model: string, id: string): Promise<Record<string, unknown> | null> {
    const record = await this.adapter.getRecord(model, id);
    if (!record || record.deleted) return null;
    return record.data;
  }

  // ------------------------------------------------------------ local writes

  /**
   * Creates a row locally and queues it.
   *
   * created_at and updated_at are filled in optimistically so the row does
   * not render with blanks until the first pull. The server stays
   * authoritative and overwrites both on the next sync.
   */
  async create(
    model: string,
    data: Record<string, unknown>,
    id?: string,
  ): Promise<string> {
    const rowId = id ?? newId();
    const now = new Date().toISOString();
    const payload: Record<string, unknown> = { ...data, id: rowId };
    if (payload.version === undefined || payload.version === null) {
      payload.version = 0;
    }
    if (typeof payload.created_at !== "string" || !payload.created_at) {
      payload.created_at = now;
    }
    if (typeof payload.updated_at !== "string" || !payload.updated_at) {
      payload.updated_at = now;
    }

    await this.adapter.putRecord({
      model,
      id: rowId,
      data: payload,
      version: 0,
      updatedAt: nowSeconds(),
      deleted: false,
    });
    await this.enqueue(model, rowId, "create", payload, 0);
    return rowId;
  }

  /** Merges a patch into the local row and queues an update. */
  async update(
    model: string,
    id: string,
    data: Record<string, unknown>,
  ): Promise<void> {
    const record = await this.adapter.getRecord(model, id);
    if (!record) {
      throw new Error("local update: " + model + "/" + id + " not found");
    }
    const merged: Record<string, unknown> = { ...record.data, ...data, id };
    await this.adapter.putRecord({
      ...record,
      data: merged,
      updatedAt: nowSeconds(),
    });
    await this.enqueue(model, id, "update", merged, record.version);
  }

  /** Removes the local row and queues a delete. */
  async remove(model: string, id: string): Promise<void> {
    const record = await this.adapter.getRecord(model, id);
    const knownVersion = record ? record.version : 0;
    if (record) await this.adapter.deleteRecord(model, id);
    await this.enqueue(model, id, "delete", null, knownVersion);
  }

  /**
   * Queues one change, squashing into whatever is already waiting for the
   * same row. At most one entry per (model, id), so the outbox stays
   * proportional to the rows you touched rather than the edits you made.
   *
   * Create-then-delete cancels both ends instead of sending a delete for a
   * row the server has never seen, which would come back as a NOT_FOUND the
   * user has no way to act on.
   */
  private async enqueue(
    model: string,
    id: string,
    op: SyncOp,
    data: Record<string, unknown> | null,
    version: number,
  ): Promise<void> {
    const existing = await this.adapter.getOutbox(model, id);

    if (existing) {
      if (op === "delete" && existing.op === "create") {
        await this.adapter.deleteOutbox(model, id);
      } else if (op === "delete") {
        await this.adapter.putOutbox({
          ...existing,
          op: "delete",
          data: null,
          hasConflict: false,
          serverData: null,
          serverVersion: 0,
          conflictMessage: "",
        });
      } else {
        // Keep the original op: an update on top of a queued create is still
        // a create as far as the server is concerned.
        await this.adapter.putOutbox({
          ...existing,
          data,
          hasConflict: false,
          serverData: null,
          serverVersion: 0,
          conflictMessage: "",
        });
      }
    } else {
      await this.adapter.putOutbox({
        model,
        entityId: id,
        op,
        data,
        version,
        createdAt: nowSeconds(),
        hasConflict: false,
        serverData: null,
        serverVersion: 0,
        conflictMessage: "",
      });
    }

    await this.refreshCounts();
    this.emit();
  }

  // -------------------------------------------------------------------- sync

  /**
   * Pull then push, for every registered model.
   *
   * Concurrent callers share one run. Two syncs draining the same outbox
   * would send every change twice, and the second copy comes back as a
   * version conflict against the first.
   */
  async sync(): Promise<SyncResult> {
    if (this.inFlight) return this.inFlight;
    this.inFlight = this.runSync().finally(() => {
      this.inFlight = null;
    });
    return this.inFlight;
  }

  private async runSync(): Promise<SyncResult> {
    this.syncing = true;
    this.lastError = null;
    this.emit();
    try {
      let pulled = 0;
      for (const model of this.models) {
        pulled += await this.pull(model);
      }
      const pushResult = await this.push();
      this.lastSyncedAt = Date.now();
      await this.adapter.setMeta(LAST_SYNCED_KEY, String(this.lastSyncedAt));
      return {
        pulled,
        pushed: pushResult.pushed,
        conflicts: pushResult.conflicts,
      };
    } catch (error) {
      this.lastError = error instanceof Error ? error.message : String(error);
      throw error;
    } finally {
      this.syncing = false;
      await this.refreshCounts();
      this.emit();
    }
  }

  /**
   * Pulls one model's changes since the stored cursor.
   *
   * Loops while the server fills a page. A cursor pull returns at most limit
   * rows, so stopping after one page would leave the mirror quietly behind
   * after any bulk change, and quietly is the worst way for a mirror to be
   * wrong.
   */
  async pull(model: string): Promise<number> {
    let total = 0;
    for (;;) {
      const since = await this.adapter.getMeta("cursor:" + model);
      let url =
        this.apiUrl +
        "/sync/pull?model=" +
        encodeURIComponent(model) +
        "&limit=" +
        String(this.pullLimit);
      if (since) url += "&since=" + encodeURIComponent(since);

      const response = await this.request(url, { method: "GET" });
      const body = (await response.json()) as PullResponse;
      const rows = body.data ?? [];

      for (const row of rows) {
        const id = typeof row.id === "string" ? row.id : "";
        if (!id) continue;
        await this.adapter.putRecord({
          model,
          id,
          data: row,
          version: toInt(row.version),
          updatedAt: parseUpdatedAt(row.updated_at),
          // A server-side delete arrives as a tombstone. The row is kept
          // rather than removed so a later re-create upserts over it.
          deleted: row._deleted === true,
        });
      }

      total += rows.length;
      const previous = since;
      if (body.cursor) {
        await this.adapter.setMeta("cursor:" + model, body.cursor);
      }

      // Stop on a short page. Also stop if the cursor did not move, which is
      // the only thing standing between a full page of same-timestamp rows
      // and an infinite loop.
      if (rows.length < this.pullLimit) break;
      if (!body.cursor || body.cursor === previous) break;
    }
    return total;
  }

  /**
   * Drains the outbox into one push and applies each result.
   *
   * Rows carrying an unresolved conflict are held back. Replaying them would
   * produce the same conflict and overwrite the server state the merge UI is
   * currently showing the user.
   */
  async push(): Promise<{ pushed: number; conflicts: number }> {
    const all = await this.adapter.listOutbox();
    const entries = all
      .filter((e) => !e.hasConflict)
      .sort((a, b) => a.createdAt - b.createdAt);
    if (entries.length === 0) return { pushed: 0, conflicts: 0 };

    const changes: PushChange[] = entries.map((e) => ({
      op: e.op,
      model: e.model,
      id: e.entityId,
      version: e.version,
      data: e.data,
    }));

    const response = await this.request(this.apiUrl + "/sync/push", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ changes }),
    });
    const body = (await response.json()) as { results: PushResult[] };
    const results = body.results ?? [];
    if (results.length !== entries.length) {
      throw new Error(
        "push: server returned " +
          String(results.length) +
          " results for " +
          String(entries.length) +
          " changes",
      );
    }

    let pushed = 0;
    let conflicts = 0;
    for (let i = 0; i < results.length; i += 1) {
      const result = results[i];
      const entry = entries[i];

      if (result.ok) {
        await this.adapter.deleteOutbox(entry.model, entry.entityId);
        if (entry.op === "delete") {
          await this.adapter.deleteRecord(entry.model, entry.entityId);
        } else {
          const record = await this.adapter.getRecord(entry.model, entry.entityId);
          if (record) {
            await this.adapter.putRecord({
              ...record,
              version: result.new_version ?? record.version,
            });
          }
        }
        pushed += 1;
        continue;
      }

      if (result.code === "VERSION_CONFLICT") {
        // Park the server's copy on the entry, so a merge UI has both sides
        // without a second round trip.
        await this.adapter.putOutbox({
          ...entry,
          hasConflict: true,
          serverData: result.server_data ?? null,
          serverVersion: result.server_version ?? 0,
          conflictMessage:
            result.message ?? "This record changed on the server",
        });
        conflicts += 1;
        continue;
      }

      // Anything else stays queued with the reason attached, so a transient
      // server error retries instead of dropping the change on the floor.
      await this.adapter.putOutbox({
        ...entry,
        conflictMessage:
          (result.code ?? "ERROR") + ": " + (result.message ?? ""),
      });
    }

    await this.refreshCounts();
    return { pushed, conflicts };
  }

  // --------------------------------------------------------------- conflicts

  async listConflicts(): Promise<OutboxEntry[]> {
    const entries = await this.adapter.listOutbox();
    return entries.filter((e) => e.hasConflict);
  }

  /**
   * Accepts the user's merge for one conflicted row.
   *
   * serverVersion becomes the version the next push claims to have seen,
   * which is what turns the retry into a legitimate update rather than the
   * same losing race a second time.
   */
  async resolveConflict(
    model: string,
    id: string,
    merged: Record<string, unknown>,
    serverVersion: number,
  ): Promise<void> {
    const entry = await this.adapter.getOutbox(model, id);
    if (!entry) throw new Error("no queued change for " + model + "/" + id);

    await this.adapter.putOutbox({
      ...entry,
      data: merged,
      version: serverVersion,
      hasConflict: false,
      serverData: null,
      serverVersion: 0,
      conflictMessage: "",
    });

    const record = await this.adapter.getRecord(model, id);
    if (record) {
      await this.adapter.putRecord({
        ...record,
        data: merged,
        version: serverVersion,
      });
    }
    await this.refreshCounts();
    this.emit();
  }

  /** Discards one queued change and puts the server's version back. */
  async revert(model: string, id: string): Promise<void> {
    const entry = await this.adapter.getOutbox(model, id);
    await this.adapter.deleteOutbox(model, id);

    if (entry && entry.op === "create") {
      // It never reached the server, so there is nothing to restore.
      await this.adapter.deleteRecord(model, id);
    } else if (entry && entry.serverData) {
      await this.adapter.putRecord({
        model,
        id,
        data: entry.serverData,
        version: entry.serverVersion,
        updatedAt: nowSeconds(),
        deleted: false,
      });
    } else {
      // No server copy on hand: clear the cursor so the next pull fetches the
      // row again, rather than leaving the edited version sitting in the
      // mirror looking authoritative.
      await this.adapter.setMeta("cursor:" + model, "");
    }

    await this.refreshCounts();
    this.emit();
  }

  async revertAll(): Promise<void> {
    const entries = await this.adapter.listOutbox();
    for (const entry of entries) {
      await this.revert(entry.model, entry.entityId);
    }
  }

  async pendingChanges(): Promise<OutboxEntry[]> {
    return this.adapter.listOutbox();
  }

  // ----------------------------------------------------------------- private

  private async request(url: string, init: RequestInit): Promise<Response> {
    const token = await this.getToken();
    const headers = new Headers(init.headers);
    if (token) headers.set("Authorization", "Bearer " + token);

    const response = await this.fetchImpl(url, { ...init, headers });
    if (!response.ok) {
      const text = await response.text().catch(() => "");
      throw new Error(
        "sync " + url + " failed: " + String(response.status) + " " + text,
      );
    }
    return response;
  }
}

function nowSeconds(): number {
  return Math.floor(Date.now() / 1000);
}

function toInt(value: unknown): number {
  if (typeof value === "number") return Math.trunc(value);
  if (typeof value === "string") {
    const parsed = Number.parseInt(value, 10);
    return Number.isNaN(parsed) ? 0 : parsed;
  }
  return 0;
}

function parseUpdatedAt(value: unknown): number {
  if (typeof value === "string") {
    const parsed = Date.parse(value);
    if (!Number.isNaN(parsed)) return Math.floor(parsed / 1000);
  }
  return nowSeconds();
}

/**
 * A UUID v4 for rows created offline.
 *
 * The client picks the id so the outbox can keep referring to the same row
 * after the server insert, which is also why the server accepts it. The
 * fallback exists because React Native has getRandomValues but, depending on
 * the runtime, not randomUUID.
 */
function newId(): string {
  const cryptoRef = globalThis.crypto as Crypto | undefined;
  if (cryptoRef && typeof cryptoRef.randomUUID === "function") {
    return cryptoRef.randomUUID();
  }

  const bytes = new Uint8Array(16);
  if (cryptoRef && typeof cryptoRef.getRandomValues === "function") {
    cryptoRef.getRandomValues(bytes);
  } else {
    for (let i = 0; i < 16; i += 1) {
      bytes[i] = Math.floor(Math.random() * 256);
    }
  }
  bytes[6] = (bytes[6] & 0x0f) | 0x40;
  bytes[8] = (bytes[8] & 0x3f) | 0x80;

  const hex: string[] = [];
  for (let i = 0; i < 16; i += 1) {
    hex.push(bytes[i].toString(16).padStart(2, "0"));
  }
  return (
    hex.slice(0, 4).join("") +
    "-" +
    hex.slice(4, 6).join("") +
    "-" +
    hex.slice(6, 8).join("") +
    "-" +
    hex.slice(8, 10).join("") +
    "-" +
    hex.slice(10, 16).join("")
  );
}
`
}
