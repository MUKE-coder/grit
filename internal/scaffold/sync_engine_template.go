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
  SyncHealth,
  SyncOp,
  SyncPolicy,
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
  private policies: Record<string, SyncPolicy> = {};
  private policiesLoaded = false;

  constructor(options: SyncEngineOptions) {
    this.apiUrl = options.apiUrl.replace(/\/+$/, "");
    this.adapter = options.adapter;
    this.models = options.models;
    this.getToken = options.getToken ?? (() => null);
    this.autoSyncIntervalMs = options.autoSyncIntervalMs ?? 30000;
    this.pullLimit = options.pullLimit ?? 500;
    this.fetchImpl = options.fetchImpl ?? globalThis.fetch.bind(globalThis);
    if (options.policies) {
      this.policies = options.policies;
      this.policiesLoaded = true;
    }
  }

  // --------------------------------------------------------------- policies

  /**
   * The declared behaviour for a model, or the defaults for one with none.
   *
   * Defaults match what every project had before policies existed, so a
   * client that has never reached the server behaves exactly as it used to
   * rather than refusing to work.
   */
  policyFor(model: string): SyncPolicy {
    return this.policies[model] ?? { mode: "offline_first", conflict: "manual" };
  }

  /**
   * Loads the policies the server declares.
   *
   * Failure is not fatal. A client that cannot reach the server is the case
   * this whole feature exists for, and it should fall back to the defaults
   * rather than refuse to open.
   */
  async loadPolicies(): Promise<void> {
    try {
      const response = await this.request(this.apiUrl + "/sync/policy", {
        method: "GET",
      });
      const body = (await response.json()) as {
        data?: { models?: Record<string, SyncPolicy> };
      };
      if (body.data?.models) {
        this.policies = body.data.models;
        this.policiesLoaded = true;
        this.emit();
      }
    } catch {
      // Keep whatever we had. Stored policies would be better still; that is
      // worth adding once there is a reason to believe they change often.
    }
  }

  /**
   * The strictest max_offline_age across mirrored models, in seconds.
   *
   * Strictest rather than per-model on purpose: the badge answers one
   * question, and if any table on the screen is too old to act on, the honest
   * answer for the screen is stale.
   */
  private strictestMaxAge(): number | null {
    let strictest: number | null = null;
    for (const model of this.models) {
      const seconds = this.policyFor(model).max_offline_age_seconds;
      if (!seconds) continue;
      if (strictest === null || seconds < strictest) strictest = seconds;
    }
    return strictest;
  }

  /** Whether the mirror is past its declared age limit. */
  isStale(): boolean {
    const limit = this.strictestMaxAge();
    if (limit === null) return false;
    // Never synced and a limit declared is the worst case, not an exemption.
    if (this.lastSyncedAt === null) return true;
    return (Date.now() - this.lastSyncedAt) / 1000 > limit;
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
    if (!this.policiesLoaded) {
      // Not awaited: the app should render from the mirror immediately, and
      // the defaults are correct until this lands.
      void this.loadPolicies();
    }
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
    const stale = this.isStale();

    // Ordered by what the user most needs to know. Stale outranks offline
    // because "your data is too old to act on" is a different message from
    // "you are offline", and only the first one should stop somebody
    // shipping against a three-day-old stock level.
    let state: SyncStatus["state"] = "synced";
    if (this.conflicts > 0) state = "conflict";
    else if (stale) state = "stale";
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
      stale,
    };
  }

  /**
   * The figures that show whether the mirror is healthy.
   *
   * Worth having as its own call rather than folding into status: this reads
   * every mirrored model to count rows, which is fine for a diagnostics panel
   * and wasteful on every badge repaint.
   */
  async health(): Promise<SyncHealth> {
    const entries = await this.adapter.listOutbox();
    const now = nowSeconds();

    let oldest = 0;
    for (const entry of entries) {
      const age = now - entry.createdAt;
      if (age > oldest) oldest = age;
    }

    const models = [];
    for (const model of this.models) {
      const records = await this.adapter.listRecords(model);
      models.push({
        model,
        mirrored: records.filter((r) => !r.deleted).length,
        pending: entries.filter((e) => e.model === model).length,
      });
    }

    return {
      pending: entries.length,
      conflicts: entries.filter((e) => e.hasConflict).length,
      oldestPendingAgeSeconds: oldest,
      lastSyncAgeSeconds:
        this.lastSyncedAt === null
          ? null
          : Math.floor((Date.now() - this.lastSyncedAt) / 1000),
      stale: this.isStale(),
      maxOfflineAgeSeconds: this.strictestMaxAge(),
      lastError: this.lastError,
      models,
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
        overridden: pushResult.overridden,
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
    if (this.policyFor(model).mode === "online_only") return 0;
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
  async push(): Promise<{
    pushed: number;
    conflicts: number;
    overridden: number;
  }> {
    const all = await this.adapter.listOutbox();
    const entries = all
      .filter((e) => !e.hasConflict)
      .sort((a, b) => a.createdAt - b.createdAt);
    if (entries.length === 0) return { pushed: 0, conflicts: 0, overridden: 0 };

    const changes: PushChange[] = entries.map((e) => ({
      op: e.op,
      model: e.model,
      id: e.entityId,
      version: e.version,
      // Stripped here as well as on the server. The server enforces it,
      // because a promise kept only by well-behaved clients is not one; the
      // client strips it too so a local_only field is not sitting in a
      // request body on the wire waiting to be enforced away.
      data: this.stripLocalOnly(e.model, e.data),
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
    let overridden = 0;
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

      if (result.code === "SERVER_WINS") {
        // The resource declared server_wins, so the server discarded this
        // change and sent its row. Nobody is asked anything: the point of
        // declaring server_wins is that there is no decision to make.
        await this.adapter.deleteOutbox(entry.model, entry.entityId);
        if (result.server_data) {
          await this.adapter.putRecord({
            model: entry.model,
            id: entry.entityId,
            data: result.server_data,
            version: result.server_version ?? 0,
            updatedAt: nowSeconds(),
            deleted: false,
          });
        }
        overridden += 1;
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
    return { pushed, conflicts, overridden };
  }

  /**
   * Removes a model's local_only fields from an outgoing payload.
   *
   * The values stay in the mirror. That is the whole point: a draft note or a
   * scratch flag lives on the device and is readable by the screen that wrote
   * it, and the server never learns it exists.
   */
  private stripLocalOnly(
    model: string,
    data: Record<string, unknown> | null,
  ): Record<string, unknown> | null {
    const localOnly = this.policyFor(model).local_only;
    if (!localOnly || localOnly.length === 0 || data === null) return data;
    const out: Record<string, unknown> = {};
    for (const [key, value] of Object.entries(data)) {
      if (!localOnly.includes(key)) out[key] = value;
    }
    return out;
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
