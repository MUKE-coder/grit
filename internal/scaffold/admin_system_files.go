package scaffold

func adminUseSystem() string {
	return `import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type { Upload, PaginatedResponse } from "@repo/shared/types";
import { apiClient, uploadFile } from "@/lib/api-client";

// ── Jobs ────────────────────────────────────────────────────────

interface QueueStats {
  queue: string;
  size: number;
  active: number;
  pending: number;
  completed: number;
  failed: number;
  retry: number;
  scheduled: number;
  processed: number;
}

interface Job {
  id: string;
  type: string;
  queue: string;
  max_retry: number;
  retried: number;
  last_error: string;
}

// ── Security + observability summaries ──────────────────────────
// Both screens poll. They used to do it with useEffect and setInterval, which
// meant hand-written cancellation, a local fetch that shadowed the global one,
// and a second copy of the request in every page that wanted the same numbers.
// React Query already owns polling, deduplication and cancellation, so these
// are useQuery with a refetchInterval and nothing else.

export function useSecuritySummary<T = Record<string, unknown>>() {
  return useQuery<T>({
    queryKey: ["admin", "security", "summary"],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/admin/security/summary");
      return (data?.data ?? data) as T;
    },
    refetchInterval: 20_000,
  });
}

export function useObservabilitySummary<T = Record<string, unknown>>() {
  return useQuery<T>({
    queryKey: ["admin", "observability", "summary"],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/admin/observability/summary");
      return (data?.data ?? data) as T;
    },
    refetchInterval: 10_000,
  });
}

// The scaling report: percentiles, connection use, the slowest queries, the
// cache hit rate, one verdict, and the stage-by-stage state of this deployment.
// Admin only, because it reports connection counts and query shapes.
//
// Polled a good deal slower than the observability summary: it runs several
// statistics queries against Postgres, and nothing it reports changes in ten
// seconds.
export function useScaleReport<T = Record<string, unknown>>() {
  return useQuery<T>({
    queryKey: ["admin", "scale", "report"],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/scale");
      return (data?.data ?? data) as T;
    },
    refetchInterval: 60_000,
  });
}

export function useJobStats() {
  return useQuery<QueueStats[]>({
    queryKey: ["admin", "jobs", "stats"],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/admin/jobs/stats");
      return data.data;
    },
    refetchInterval: 5000,
  });
}

export function useJobsByStatus(status: string, queue = "default") {
  return useQuery<Job[]>({
    queryKey: ["admin", "jobs", status, queue],
    queryFn: async () => {
      const { data } = await apiClient.get(` + "`" + `/api/admin/jobs/${status}?queue=${queue}` + "`" + `);
      return data.data;
    },
    refetchInterval: 5000,
  });
}

export function useRetryJob() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ id, queue }: { id: string; queue?: string }) => {
      await apiClient.post(` + "`" + `/api/admin/jobs/${id}/retry?queue=${queue || "default"}` + "`" + `);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "jobs"] });
    },
  });
}

export function useClearQueue() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (queue: string) => {
      await apiClient.delete(` + "`" + `/api/admin/jobs/queue/${queue}` + "`" + `);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "jobs"] });
    },
  });
}

// ── Files ───────────────────────────────────────────────────────
// Upload + PaginatedResponse are imported from @repo/shared/types so
// the same shapes flow through web, admin, and the Go API (via
// grit sync) — no inline duplicates that silently drift.
type UploadListResponse = PaginatedResponse<Upload>;

export function useUploads(page = 1, pageSize = 20) {
  return useQuery<UploadListResponse>({
    queryKey: ["admin", "uploads", page, pageSize],
    queryFn: async () => {
      const { data } = await apiClient.get(` + "`" + `/api/uploads?page=${page}&page_size=${pageSize}` + "`" + `);
      return data;
    },
  });
}

// v3.31.32 — storage stats. Returns total file count, total bytes,
// and a per-kind breakdown (image / video / audio / pdf / document /
// spreadsheet / other) so the Files admin page can show usage at a
// glance.
export interface UploadStats {
  total_count: number;
  total_size: number;
  by_kind: { kind: string; count: number; size: number }[];
}

export function useUploadStats() {
  return useQuery<{ data: UploadStats }>({
    queryKey: ["admin", "uploads", "stats"],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/uploads/stats");
      return data;
    },
  });
}

export function useUploadFile() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (file: File) => {
      const result = await uploadFile(file);
      // uploadFile returns Record<string, unknown>; the server actually
      // sends the Upload shape so a two-step assertion through unknown
      // is sound and matches the shared type without a runtime cost.
      return result.data as unknown as Upload;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "uploads"] });
    },
  });
}

export function useDeleteUpload() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      await apiClient.delete(` + "`" + `/api/uploads/${id}` + "`" + `);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "uploads"] });
    },
  });
}

// ── Cron ────────────────────────────────────────────────────────

interface CronTask {
  name: string;
  schedule: string;
  type: string;
}

export function useCronTasks() {
  return useQuery<CronTask[]>({
    queryKey: ["admin", "cron", "tasks"],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/admin/cron/tasks");
      return data.data;
    },
  });
}
`
}

// adminSagasPage is the saga runs screen.
//
// A run that goes stuck needs a person: its compensation failed past its
// attempts, so there is a charge, a reservation or a booking still standing
// that was supposed to be taken back. Without this screen the only way to find
// one is SQL against saga_runs, which means nobody finds one until a customer
// complains.
// adminRoutesPage is the route table in a browser.
//
// `grit routes` has printed it in a terminal for a long time and remains the
// fuller answer, because the CLI parses routes.go and can report the middleware
// group and the permission each route wants. What nothing provided was a
// searchable version for somebody already in the admin, looking for the URL of
// an endpoint they are about to call.
func adminRoutesPage() string {
	return `"use client";

import { useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { apiClient, API_URL } from "@/lib/api-client";
import { inputClasses } from "@/components/ui/input";
import { Search, Copy, Check, Server, Loader2 } from "@/lib/icons";

interface Access {
  group: string;
  roles?: string[];
  permissions?: string[];
  authenticated: boolean;
  api_key_only?: boolean;
  staff?: boolean;
  summary: string;
}

interface Route {
  method: string;
  path: string;
  handler?: string;
  group: string;
  access_known: boolean;
  access?: Access;
}

interface RoutesResponse {
  data: Route[];
  count: number;
  covered?: { routes: number; note?: string };
}

const METHOD_TONE: Record<string, string> = {
  GET: "border-info/30 bg-info/10 text-info",
  POST: "border-success/30 bg-success/10 text-success",
  PUT: "border-warning/30 bg-warning/10 text-warning",
  PATCH: "border-warning/30 bg-warning/10 text-warning",
  DELETE: "border-danger/30 bg-danger/10 text-danger",
};

const METHODS = ["GET", "POST", "PUT", "PATCH", "DELETE"] as const;

// The four states a route's access can be in, which is three more than a
// boolean and one more than feels necessary.
//
// "unknown" is the one that earns its place. The reference, the profiler and
// the database browser mount their own routes, and this project's route files
// never saw them, so nothing can say what they require. Calling those public
// would be the one wrong answer worth avoiding: a reader takes this column as a
// security statement, and /auth/me is not public.
type AccessKind = "unknown" | "public" | "session" | "guarded";

function kindOf(route: Route): AccessKind {
  if (!route.access_known || !route.access) return "unknown";
  const { roles, permissions, authenticated } = route.access;
  if ((roles && roles.length > 0) || (permissions && permissions.length > 0)) return "guarded";
  if (authenticated) return "session";
  return "public";
}

const ACCESS_TONE: Record<AccessKind, string> = {
  unknown: "border-dashed border-border text-text-muted",
  public: "border-border text-text-secondary",
  session: "border-info/30 bg-info/10 text-info",
  guarded: "border-accent/30 bg-accent/10 text-accent",
};

const ACCESS_LABEL: Record<AccessKind, string> = {
  unknown: "not ours",
  public: "public",
  session: "signed in",
  guarded: "guarded",
};

function AccessBadge({ route }: { route: Route }) {
  const kind = kindOf(route);
  const text = kind === "unknown" ? "unknown" : (route.access?.summary ?? "public");
  return (
    <span
      className={
        "inline-flex whitespace-nowrap rounded border px-1.5 py-0.5 font-mono text-[11px] " +
        ACCESS_TONE[kind]
      }
      title={
        kind === "unknown"
          ? "This route was registered by a mounted dashboard, not by this project's route files, so what it requires is not knowable here."
          : undefined
      }
    >
      {text}
    </span>
  );
}

// A curl for the route, with the path parameters left as they are written. A
// command with :id still in it is obviously a template; one with a plausible
// fake id in it is the kind of thing people paste without reading.
function curlFor(route: Route): string {
  const url = API_URL.replace(/\/$/, "") + route.path;
  const parts = ["curl -X " + route.method, '"' + url + '"'];
  if (route.method !== "GET" && route.method !== "DELETE") {
    parts.push('-H "Content-Type: application/json"');
    parts.push("-d '{}'");
  }
  // Omitted for a route that needs nothing: a public endpoint's curl should be
  // the command that actually works, not one carrying a header the reader then
  // has to decide about.
  if (kindOf(route) !== "public") {
    parts.push('-H "Authorization: Bearer $TOKEN"');
  }
  return parts.join(" \\\n  ");
}

function CopyButton({ text }: { text: string }) {
  const [copied, setCopied] = useState(false);
  return (
    <button
      onClick={async () => {
        try {
          await navigator.clipboard.writeText(text);
          setCopied(true);
          setTimeout(() => setCopied(false), 1500);
        } catch {
          // A browser that refuses the clipboard, usually because the page is
          // not on https. Saying nothing is better than an error for something
          // the person can still select by hand.
        }
      }}
      title="Copy as curl"
      className="inline-flex items-center gap-1 rounded border border-border px-1.5 py-1 text-[11px] text-text-muted hover:bg-bg-hover hover:text-foreground"
    >
      {copied ? <Check className="h-3 w-3 text-success" /> : <Copy className="h-3 w-3" />}
      {copied ? "copied" : "curl"}
    </button>
  );
}

export default function RoutesPage() {
  const [query, setQuery] = useState("");
  const [method, setMethod] = useState<string>("");
  const [kind, setKind] = useState<AccessKind | "">("");

  const { data, isLoading } = useQuery<RoutesResponse>({
    queryKey: ["system-routes"],
    queryFn: async () => {
      const { data } = await apiClient.get<RoutesResponse>("/api/admin/routes");
      return data;
    },
    // The table changes when the binary does, which is not while somebody is
    // looking at it.
    staleTime: 5 * 60_000,
  });

  const all = data?.data ?? [];

  // Filtered in the browser rather than by refetching: it is a few hundred rows
  // of three short strings, and typing should not wait on a round trip.
  const routes = useMemo(() => {
    const q = query.trim().toLowerCase();
    return all.filter((r) => {
      if (method && r.method !== method) return false;
      if (kind && kindOf(r) !== kind) return false;
      if (!q) return true;
      return (
        r.path.toLowerCase().includes(q) ||
        (r.handler ?? "").toLowerCase().includes(q) ||
        (r.access?.summary ?? "").toLowerCase().includes(q) ||
        r.group.toLowerCase().includes(q)
      );
    });
  }, [all, query, method, kind]);

  // The counts are the reason this screen is worth opening rather than running
  // the command: "how many of my endpoints are reachable with no credentials"
  // is one number, and it was not available anywhere before.
  const counts = useMemo(() => {
    const out: Record<AccessKind, number> = { unknown: 0, public: 0, session: 0, guarded: 0 };
    for (const route of all) out[kindOf(route)] += 1;
    return out;
  }, [all]);

  const groups = useMemo(() => {
    const out = new Map<string, Route[]>();
    for (const route of routes) {
      const list = out.get(route.group);
      if (list) list.push(route);
      else out.set(route.group, [route]);
    }
    return [...out.entries()];
  }, [routes]);

  return (
    <div className="space-y-6">
      <PageHeader
        title="Routes"
        subtitle="Every endpoint this API serves, and what each one asks of a caller."
      />

      <div className="grid gap-3 sm:grid-cols-4">
        {(["public", "session", "guarded", "unknown"] as const).map((k) => (
          <button
            key={k}
            onClick={() => setKind(kind === k ? "" : k)}
            className={
              "rounded-xl border p-4 text-left transition-colors " +
              (kind === k ? "border-accent bg-accent/5" : "border-border bg-bg-secondary hover:bg-bg-hover")
            }
          >
            <div className="font-mono text-2xl font-semibold text-foreground">{counts[k]}</div>
            <div className="mt-1 text-xs text-text-secondary">{ACCESS_LABEL[k]}</div>
          </button>
        ))}
      </div>

      <div className="flex flex-wrap items-center gap-2">
        <div className="relative flex-1 min-w-[16rem]">
          <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-text-muted" />
          <input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Filter by path, handler, group or permission"
            className={inputClasses({ className: "pl-9" })}
          />
        </div>
        {METHODS.map((m) => (
          <button
            key={m}
            onClick={() => setMethod(method === m ? "" : m)}
            className={
              "rounded-full border px-2.5 py-1 text-xs font-semibold " +
              (method === m ? METHOD_TONE[m] : "border-border text-text-secondary hover:bg-bg-hover")
            }
          >
            {m}
          </button>
        ))}
        <span className="text-xs text-text-muted">
          {routes.length} of {data?.count ?? 0}
        </span>
      </div>

      {isLoading ? (
        <div className="flex items-center gap-2 text-sm text-text-muted">
          <Loader2 className="h-4 w-4 animate-spin" /> Reading the route table
        </div>
      ) : groups.length === 0 ? (
        <div className="rounded-xl border border-border bg-bg-secondary p-8 text-center">
          <Server className="mx-auto h-8 w-8 text-text-muted" />
          <p className="mt-3 text-sm font-semibold text-foreground">Nothing matches</p>
        </div>
      ) : (
        <div className="space-y-5">
          {groups.map(([group, rows]) => (
            <section key={group}>
              <h2 className="mb-2 text-xs font-semibold uppercase tracking-wider text-text-muted">
                {group} <span className="font-mono normal-case opacity-70">{rows.length}</span>
              </h2>
              <div className="overflow-x-auto rounded-xl border border-border">
                <table className="w-full text-sm">
                  <tbody className="divide-y divide-border">
                    {rows.map((route) => (
                      <tr key={route.method + route.path} className="hover:bg-bg-hover">
                        <td className="w-24 px-3 py-2">
                          <span
                            className={
                              "inline-flex rounded border px-1.5 py-0.5 font-mono text-[11px] font-semibold " +
                              (METHOD_TONE[route.method] ?? "border-border text-text-secondary")
                            }
                          >
                            {route.method}
                          </span>
                        </td>
                        <td className="whitespace-nowrap px-3 py-2 font-mono text-xs text-foreground">
                          {route.path}
                        </td>
                        <td className="w-px px-3 py-2">
                          <AccessBadge route={route} />
                        </td>
                        <td className="hidden px-3 py-2 xl:table-cell">
                          <div className="ml-auto max-w-[10rem] truncate text-right font-mono text-[11px] text-text-muted">
                            {route.handler}
                          </div>
                        </td>
                        <td className="w-20 px-3 py-2 text-right">
                          <CopyButton text={curlFor(route)} />
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </section>
          ))}
        </div>
      )}

      <p className="max-w-3xl text-xs text-text-muted">
        Method, path and handler come from the router. Access comes from a table{" "}
        <code className="font-mono">grit</code> builds out of{" "}
        <code className="font-mono">routes.go</code> and the per-resource route files, because
        the router does not carry it: Gin records a route&apos;s last handler and nothing about
        the middleware in front of it. The table covers{" "}
        {data?.covered?.routes ?? 0} routes, and anything outside it reads{" "}
        <span className="font-mono">unknown</span> rather than public, which is what the
        reference, the profiler and the database browser are. Run{" "}
        <code className="font-mono">grit routes</code> for the same data in a terminal, or{" "}
        <code className="font-mono">grit doctor</code> if you have added a route by hand and
        want to know whether the table is still current.
      </p>
    </div>
  );
}
`
}

func adminSagasPage() string {
	return `"use client";

import { useState } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { apiClient } from "@/lib/api-client";
import { useToastedMutation } from "@/hooks/use-toasted-mutation";
import {
  RefreshCw, AlertTriangle, CheckCircle, Loader2, ArrowRight, Workflow,
} from "@/lib/icons";

// The five statuses a run can be in. See internal/saga: only "stuck" needs a
// person, which is why it is the one with a button.
type RunStatus = "running" | "compensating" | "done" | "compensated" | "stuck";

interface SagaStep {
  id: string;
  idx: number;
  name: string;
  status: string;
  attempts: number;
  last_error?: string;
  started_at?: string;
  finished_at?: string;
}

interface SagaRun {
  id: string;
  name: string;
  key?: string;
  status: RunStatus;
  cursor: number;
  attempts: number;
  last_error?: string;
  created_at: string;
  finished_at?: string;
  steps?: SagaStep[];
}

interface RunsResponse {
  data: SagaRun[];
  meta: { total: number; page: number; page_size: number; pages: number };
  counts?: Record<string, number>;
}

const STATUS_TONE: Record<string, string> = {
  running: "border-info/30 bg-info/10 text-info",
  compensating: "border-warning/30 bg-warning/10 text-warning",
  done: "border-success/30 bg-success/10 text-success",
  compensated: "border-border bg-bg-elevated text-text-secondary",
  stuck: "border-danger/30 bg-danger/10 text-danger",
  // Step statuses, which share the palette.
  pending: "border-border bg-bg-elevated text-text-muted",
  failed: "border-danger/30 bg-danger/10 text-danger",
};

// What each status means, in the words an operator needs rather than the words
// the code uses. "compensated" reading as a failure is the misunderstanding
// worth preventing: it is the good outcome of a bad run.
const STATUS_HELP: Record<string, string> = {
  running: "Working through the steps.",
  compensating: "A step failed for good and the completed ones are being undone.",
  done: "Every step completed.",
  compensated: "A step failed and everything before it was undone. The world is back where it started.",
  stuck: "A compensation itself failed. Something happened that could not be taken back: this needs you.",
};

function StatusChip({ status, count }: { status: string; count?: number }) {
  return (
    <span
      title={STATUS_HELP[status]}
      className={
        "inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-xs font-semibold " +
        (STATUS_TONE[status] ?? "border-border bg-bg-elevated text-text-secondary")
      }
    >
      {status}
      {count != null && <span className="font-mono opacity-70">{count}</span>}
    </span>
  );
}

export default function SagasPage() {
  const queryClient = useQueryClient();
  const [status, setStatus] = useState<string>("");
  const [openRun, setOpenRun] = useState<string | null>(null);

  const { data, isLoading, refetch, isFetching } = useQuery<RunsResponse>({
    queryKey: ["system-sagas", status],
    queryFn: async () => {
      const { data } = await apiClient.get<RunsResponse>("/api/admin/sagas", {
        params: { page_size: 50, ...(status ? { "filter[status]": status } : {}) },
      });
      return data;
    },
    refetchInterval: 10_000,
  });

  const { data: detail } = useQuery<{ data: SagaRun }>({
    queryKey: ["system-saga", openRun],
    queryFn: async () => {
      const { data } = await apiClient.get<{ data: SagaRun }>("/api/admin/sagas/" + openRun);
      return data;
    },
    enabled: !!openRun,
  });

  const retry = useToastedMutation<unknown, unknown, string>({
    mutationFn: async (id: string) => {
      const { data } = await apiClient.post("/api/admin/sagas/" + id + "/retry");
      return data;
    },
    successMessage: "The run is compensating again",
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["system-sagas"] });
      queryClient.invalidateQueries({ queryKey: ["system-saga"] });
    },
  });

  const counts = data?.counts ?? {};
  const stuck = counts.stuck ?? 0;
  const runs = data?.data ?? [];
  const open = detail?.data;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Sagas"
        subtitle="Multi-step processes that have to either finish or be undone."
        actions={
          <button
            onClick={() => refetch()}
            className="inline-flex items-center gap-2 rounded-lg border border-border px-3 py-2 text-sm text-text-secondary hover:bg-bg-hover"
          >
            <RefreshCw className={"h-4 w-4 " + (isFetching ? "animate-spin" : "")} /> Refresh
          </button>
        }
      />

      {stuck > 0 && (
        <div className="flex items-start gap-2 rounded-lg border border-danger/30 bg-danger/5 px-4 py-3 text-sm text-danger">
          <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0" />
          <span>
            {stuck} run{stuck === 1 ? "" : "s"} could not be undone. Each one did something that was
            not taken back, so there is a charge, a reservation or a booking still standing. Fix
            what the error names, then retry the run.
          </span>
        </div>
      )}

      <div className="flex flex-wrap items-center gap-2">
        <button
          onClick={() => setStatus("")}
          className={
            "rounded-full border px-2.5 py-1 text-xs font-semibold " +
            (status === "" ? "border-accent/40 bg-accent/10 text-accent" : "border-border text-text-secondary hover:bg-bg-hover")
          }
        >
          all {data?.meta?.total ?? 0}
        </button>
        {(["running", "compensating", "stuck", "done", "compensated"] as const).map((s) => (
          <button key={s} onClick={() => setStatus(status === s ? "" : s)} className={status === s ? "ring-2 ring-accent/40 rounded-full" : ""}>
            <StatusChip status={s} count={counts[s] ?? 0} />
          </button>
        ))}
      </div>

      {isLoading ? (
        <div className="flex items-center gap-2 text-sm text-text-muted">
          <Loader2 className="h-4 w-4 animate-spin" /> Reading the runs
        </div>
      ) : runs.length === 0 ? (
        <div className="rounded-xl border border-border bg-bg-secondary p-8 text-center">
          <Workflow className="mx-auto h-8 w-8 text-text-muted" />
          <p className="mt-3 text-sm font-semibold text-foreground">No runs yet</p>
          <p className="mt-1 text-xs text-text-secondary">
            Generate one with <code className="font-mono">grit generate workflow Checkout --steps &quot;charge,ship&quot;</code>,
            then start it with <code className="font-mono">saga.Start</code>.
          </p>
        </div>
      ) : (
        <div className="overflow-hidden rounded-xl border border-border">
          <table className="w-full text-sm">
            <thead className="bg-bg-elevated text-left text-xs uppercase tracking-wider text-text-muted">
              <tr>
                <th className="px-4 py-3 font-medium">Saga</th>
                <th className="px-4 py-3 font-medium">Status</th>
                <th className="px-4 py-3 font-medium">Step</th>
                <th className="px-4 py-3 font-medium">Started</th>
                <th className="px-4 py-3 font-medium"></th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {runs.map((run) => (
                <tr key={run.id} className="hover:bg-bg-hover">
                  <td className="px-4 py-3">
                    <p className="font-medium text-foreground">{run.name}</p>
                    <p className="font-mono text-[11px] text-text-muted">{run.key ?? run.id}</p>
                  </td>
                  <td className="px-4 py-3">
                    <StatusChip status={run.status} />
                    {run.last_error && (
                      <p className="mt-1 max-w-md truncate text-[11px] text-text-muted" title={run.last_error}>
                        {run.last_error}
                      </p>
                    )}
                  </td>
                  <td className="px-4 py-3 font-mono text-xs text-text-secondary">{run.cursor}</td>
                  <td className="px-4 py-3 text-xs text-text-secondary">
                    {new Date(run.created_at).toLocaleString()}
                  </td>
                  <td className="px-4 py-3 text-right">
                    <div className="flex items-center justify-end gap-2">
                      {run.status === "stuck" && (
                        <button
                          onClick={() => retry.mutate(run.id)}
                          disabled={retry.isPending}
                          className="inline-flex items-center gap-1.5 rounded-lg border border-danger/30 bg-danger/10 px-2.5 py-1.5 text-xs font-semibold text-danger hover:bg-danger/20 disabled:opacity-50"
                        >
                          <RefreshCw className="h-3.5 w-3.5" /> Retry
                        </button>
                      )}
                      <button
                        onClick={() => setOpenRun(openRun === run.id ? null : run.id)}
                        className="inline-flex items-center gap-1 text-xs text-text-secondary hover:text-foreground"
                      >
                        Steps <ArrowRight className="h-3.5 w-3.5" />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {open && (
        <div className="rounded-xl border border-border bg-bg-secondary p-4">
          <div className="mb-3 flex items-center justify-between">
            <div>
              <p className="text-sm font-semibold text-foreground">{open.name}</p>
              <p className="font-mono text-[11px] text-text-muted">{open.id}</p>
            </div>
            <StatusChip status={open.status} />
          </div>
          <ol className="space-y-2">
            {(open.steps ?? []).map((step) => (
              <li key={step.id} className="flex items-start gap-3 rounded-lg border border-border bg-background px-3 py-2">
                <span className="mt-0.5 font-mono text-[11px] text-text-muted">{step.idx}</span>
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <span className="text-sm text-foreground">{step.name}</span>
                    <StatusChip status={step.status} />
                    {step.attempts > 1 && (
                      <span className="text-[11px] text-text-muted">{step.attempts} attempts</span>
                    )}
                  </div>
                  {step.last_error && (
                    <p className="mt-1 break-words text-[11px] text-danger">{step.last_error}</p>
                  )}
                </div>
                {step.status === "done" && <CheckCircle className="h-4 w-4 shrink-0 text-success" />}
              </li>
            ))}
          </ol>
          <p className="mt-3 text-[11px] text-text-muted">
            Steps run top to bottom and are undone bottom to top. A step with no Undo is skipped on
            the way back, which is why anything you cannot take back belongs at the end.
          </p>
        </div>
      )}
    </div>
  );
}
`
}

func adminJobsPage() string {
	return `"use client";

import { useState } from "react";
import { useJobStats, useJobsByStatus, useRetryJob, useClearQueue } from "@/hooks/use-system";
import { Briefcase, RefreshCw, Trash2, Loader2 } from "@/lib/icons";

const statuses = ["active", "pending", "completed", "failed", "retry"] as const;

export default function JobsPage() {
  const [activeTab, setActiveTab] = useState<string>("active");
  const { data: stats, isLoading: statsLoading } = useJobStats();
  const { data: jobs, isLoading: jobsLoading } = useJobsByStatus(activeTab);
  const retryJob = useRetryJob();
  const clearQueue = useClearQueue();

  const totalStats = stats?.[0];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Background Jobs</h1>
          <p className="text-sm text-text-secondary mt-1">Monitor and manage background job queues</p>
        </div>
      </div>

      {/* Stats cards */}
      <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
        {[
          { label: "Active", value: totalStats?.active ?? 0, color: "text-info" },
          { label: "Pending", value: totalStats?.pending ?? 0, color: "text-warning" },
          { label: "Completed", value: totalStats?.completed ?? 0, color: "text-success" },
          { label: "Failed", value: totalStats?.failed ?? 0, color: "text-danger" },
          { label: "Retry", value: totalStats?.retry ?? 0, color: "text-accent" },
          { label: "Processed", value: totalStats?.processed ?? 0, color: "text-foreground" },
        ].map((stat) => (
          <div key={stat.label} className="rounded-xl border border-border bg-bg-secondary p-4">
            <p className="text-xs text-text-muted uppercase tracking-wider">{stat.label}</p>
            <p className={` + "`" + `text-2xl font-bold mt-1 ${stat.color}` + "`" + `}>
              {statsLoading ? "—" : stat.value.toLocaleString()}
            </p>
          </div>
        ))}
      </div>

      {/* Tab navigation */}
      <div className="flex items-center gap-1 border-b border-border">
        {statuses.map((status) => (
          <button
            key={status}
            onClick={() => setActiveTab(status)}
            className={` + "`" + `px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${
              activeTab === status
                ? "border-accent text-accent"
                : "border-transparent text-text-secondary hover:text-foreground"
            }` + "`" + `}
          >
            {status.charAt(0).toUpperCase() + status.slice(1)}
          </button>
        ))}

        <div className="ml-auto">
          <button
            onClick={() => clearQueue.mutate("default")}
            className="flex items-center gap-2 px-3 py-1.5 text-sm text-text-secondary hover:text-danger transition-colors"
          >
            <Trash2 className="h-3.5 w-3.5" />
            Clear completed
          </button>
        </div>
      </div>

      {/* Job list */}
      <div className="rounded-xl border border-border bg-bg-secondary overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="border-b border-border">
              <th className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase">ID</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase">Type</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase">Queue</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase">Retries</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase">Error</th>
              <th className="px-4 py-3 text-right text-xs font-medium text-text-muted uppercase">Actions</th>
            </tr>
          </thead>
          <tbody>
            {jobsLoading ? (
              <tr>
                <td colSpan={6} className="px-4 py-12 text-center">
                  <Loader2 className="h-6 w-6 animate-spin text-accent mx-auto" />
                </td>
              </tr>
            ) : jobs && jobs.length > 0 ? (
              jobs.map((job) => (
                <tr key={job.id} className="border-b border-border last:border-0 hover:bg-bg-hover transition-colors">
                  <td className="px-4 py-3 text-sm font-mono text-text-secondary">{job.id.slice(0, 8)}...</td>
                  <td className="px-4 py-3">
                    <span className="rounded-md bg-accent/10 px-2 py-0.5 text-xs font-medium text-accent">
                      {job.type}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-sm text-text-secondary">{job.queue}</td>
                  <td className="px-4 py-3 text-sm text-text-secondary">{job.retried}/{job.max_retry}</td>
                  <td className="px-4 py-3 text-sm text-danger max-w-[200px] truncate">{job.last_error || "—"}</td>
                  <td className="px-4 py-3 text-right">
                    {(activeTab === "failed" || activeTab === "retry") && (
                      <button
                        onClick={() => retryJob.mutate({ id: job.id, queue: job.queue })}
                        className="inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-medium text-accent hover:bg-accent/10 rounded-md transition-colors"
                      >
                        <RefreshCw className="h-3 w-3" />
                        Retry
                      </button>
                    )}
                  </td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan={6} className="px-4 py-12 text-center">
                  <Briefcase className="h-8 w-8 text-text-muted mx-auto mb-2" />
                  <p className="text-sm text-text-secondary">No {activeTab} jobs</p>
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
`
}

func adminFilesPage() string {
	return `"use client";

import { useState, useRef } from "react";
import { useUploads, useUploadFile, useDeleteUpload, useUploadStats } from "@/hooks/use-system";
import { FolderOpen, Upload, Loader2, X, Image as ImageIcon, Play, Music, FileText, FileSpreadsheet, File as FileIcon } from "@/lib/icons";
import { buttonClasses } from "@/components/ui/button";

function formatFileSize(bytes: number): string {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + " " + sizes[i];
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

export default function FilesPage() {
  const [page, setPage] = useState(1);
  const { data, isLoading } = useUploads(page);
  const { data: statsResp } = useUploadStats();
  const uploadFile = useUploadFile();
  const deleteUpload = useDeleteUpload();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [dragOver, setDragOver] = useState(false);

  const uploads = data?.data ?? [];
  const meta = data?.meta;
  const stats = statsResp?.data;

  const handleFiles = (files: FileList | null) => {
    if (!files) return;
    for (const file of Array.from(files)) uploadFile.mutate(file);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setDragOver(false);
    handleFiles(e.dataTransfer.files);
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">File Storage</h1>
          <p className="text-sm text-text-secondary mt-1">Manage uploaded files and images</p>
        </div>
        <button
          onClick={() => fileInputRef.current?.click()}
          className={buttonClasses()}
        >
          <Upload className="h-4 w-4" />
          Upload File
        </button>
        <input
          ref={fileInputRef}
          type="file"
          className="hidden"
          multiple
          onChange={(e) => handleFiles(e.target.files)}
        />
      </div>

      {/* v3.31.32 — storage stats. Total + per-kind breakdown so
          you can see at a glance what's eating your bucket. */}
      {stats && <StorageStatsPanel stats={stats} />}

      {/* Drop zone */}
      <div
        onDragOver={(e) => { e.preventDefault(); setDragOver(true); }}
        onDragLeave={() => setDragOver(false)}
        onDrop={handleDrop}
        className={` + "`" + `rounded-xl border-2 border-dashed p-8 text-center transition-colors ${
          dragOver ? "border-accent bg-accent/5" : "border-border"
        }` + "`" + `}
      >
        <Upload className="h-8 w-8 text-text-muted mx-auto mb-2" />
        <p className="text-sm text-text-secondary">
          Drag & drop files here, or{" "}
          <button
            onClick={() => fileInputRef.current?.click()}
            className="text-accent hover:underline"
          >
            browse
          </button>
        </p>
        <p className="text-xs text-text-muted mt-1">Max 10 MB per file</p>
      </div>

      {/* Upload progress */}
      {uploadFile.isPending && (
        <div className="flex items-center gap-3 rounded-lg border border-border bg-bg-secondary p-4">
          <Loader2 className="h-5 w-5 animate-spin text-accent" />
          <span className="text-sm text-text-secondary">Uploading...</span>
        </div>
      )}

      {/* File grid */}
      {isLoading ? (
        <div className="flex justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-accent" />
        </div>
      ) : uploads.length === 0 ? (
        <div className="flex flex-col items-center py-16">
          <FolderOpen className="h-12 w-12 text-text-muted mb-3" />
          <p className="text-sm text-text-secondary">No files uploaded yet</p>
        </div>
      ) : (
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
          {uploads.map((upload) => (
            <div
              key={upload.id}
              className="group relative rounded-xl border border-border bg-bg-secondary overflow-hidden hover:border-accent/50 transition-colors"
            >
              {/* Preview */}
              <div className="aspect-square bg-bg-tertiary flex items-center justify-center">
                {upload.mime_type.startsWith("image/") ? (
                  <img
                    src={upload.thumbnail_url || upload.url}
                    alt={upload.original_name}
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <div className="text-center p-4">
                    <FolderOpen className="h-8 w-8 text-text-muted mx-auto mb-1" />
                    <p className="text-xs text-text-muted uppercase">
                      {upload.mime_type.split("/")[1]?.slice(0, 4)}
                    </p>
                  </div>
                )}
              </div>

              {/* Info */}
              <div className="p-3">
                <p className="text-sm font-medium text-foreground truncate">{upload.original_name}</p>
                <div className="flex items-center justify-between mt-1">
                  <span className="text-xs text-text-muted">{formatFileSize(upload.size)}</span>
                  <span className="text-xs text-text-muted">{formatDate(upload.created_at)}</span>
                </div>
              </div>

              {/* Delete overlay */}
              <button
                onClick={() => deleteUpload.mutate(upload.id)}
                className="absolute top-2 right-2 p-1.5 rounded-lg bg-black/50 text-white opacity-0 group-hover:opacity-100 hover:bg-danger transition-all"
              >
                <X className="h-3.5 w-3.5" />
              </button>
            </div>
          ))}
        </div>
      )}

      {/* Pagination */}
      {meta && meta.pages > 1 && (
        <div className="flex items-center justify-center gap-2">
          <button
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page === 1}
            className="px-3 py-1.5 text-sm rounded-lg border border-border text-text-secondary hover:bg-bg-hover disabled:opacity-50 transition-colors"
          >
            Previous
          </button>
          <span className="text-sm text-text-secondary">
            Page {page} of {meta.pages}
          </span>
          <button
            onClick={() => setPage((p) => Math.min(meta.pages, p + 1))}
            disabled={page === meta.pages}
            className="px-3 py-1.5 text-sm rounded-lg border border-border text-text-secondary hover:bg-bg-hover disabled:opacity-50 transition-colors"
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
}

// v3.31.32 — Storage stats panel. Shows total bytes / file count up
// top, then a per-kind breakdown so you can see what's filling the
// bucket. Per-kind rows compute their share via the count + size from
// the API; the bar fills proportional to the total size.
interface UploadStatsLite {
  total_count: number;
  total_size: number;
  by_kind: { kind: string; count: number; size: number }[];
}

const KIND_META: Record<string, { label: string; icon: React.ComponentType<{ className?: string }>; tint: string }> = {
  image:       { label: "Images",       icon: ImageIcon,       tint: "text-accent" },
  video:       { label: "Videos",       icon: Play,            tint: "text-info" },
  audio:       { label: "Audio",        icon: Music,           tint: "text-info" },
  pdf:         { label: "PDFs",         icon: FileText,        tint: "text-danger" },
  document:    { label: "Documents",    icon: FileText,        tint: "text-info" },
  spreadsheet: { label: "Spreadsheets", icon: FileSpreadsheet, tint: "text-success" },
  other:       { label: "Other",        icon: FileIcon,        tint: "text-text-muted" },
};

function formatStatsSize(bytes: number): string {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + " " + sizes[i];
}

function StorageStatsPanel({ stats }: { stats: UploadStatsLite }) {
  // Sort kinds by descending size so the biggest consumers are first.
  const rows = [...(stats.by_kind ?? [])].sort((a, b) => b.size - a.size);
  return (
    <div className="rounded-xl border border-border bg-bg-secondary">
      <div className="grid grid-cols-2 gap-px bg-border md:grid-cols-3">
        <div className="bg-bg-secondary p-4">
          <p className="text-[11px] font-semibold uppercase tracking-wider text-text-muted">Total files</p>
          <p className="mt-1 text-2xl font-bold text-foreground tabular-nums">
            {stats.total_count.toLocaleString()}
          </p>
        </div>
        <div className="bg-bg-secondary p-4">
          <p className="text-[11px] font-semibold uppercase tracking-wider text-text-muted">Total storage</p>
          <p className="mt-1 text-2xl font-bold text-foreground">
            {formatStatsSize(stats.total_size)}
          </p>
        </div>
        <div className="col-span-2 bg-bg-secondary p-4 md:col-span-1">
          <p className="text-[11px] font-semibold uppercase tracking-wider text-text-muted">Avg file size</p>
          <p className="mt-1 text-2xl font-bold text-foreground">
            {stats.total_count > 0 ? formatStatsSize(stats.total_size / stats.total_count) : "—"}
          </p>
        </div>
      </div>

      {rows.length > 0 && (
        <div className="border-t border-border px-4 py-3 space-y-2.5">
          {rows.map((row) => {
            const meta = KIND_META[row.kind] || KIND_META.other;
            const Icon = meta.icon;
            const share = stats.total_size > 0 ? (row.size / stats.total_size) * 100 : 0;
            return (
              <div key={row.kind} className="flex items-center gap-3">
                <div className={"flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-bg-tertiary " + meta.tint}>
                  <Icon className="h-3.5 w-3.5" />
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between text-xs">
                    <span className="font-medium text-foreground">{meta.label}</span>
                    <span className="tabular-nums text-text-muted">
                      {row.count.toLocaleString()} files · {formatStatsSize(row.size)}
                    </span>
                  </div>
                  <div className="mt-1.5 h-1 w-full overflow-hidden rounded-full bg-bg-tertiary">
                    <div
                      className="h-full rounded-full bg-accent transition-all"
                      style={{ width: share + "%" }}
                    />
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
`
}

func adminCronPage() string {
	return `"use client";

import { useCronTasks } from "@/hooks/use-system";
import { Calendar, Loader2 } from "@/lib/icons";

export default function CronPage() {
  const { data: tasks, isLoading } = useCronTasks();

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-foreground">Cron Scheduler</h1>
        <p className="text-sm text-text-secondary mt-1">View registered scheduled tasks</p>
      </div>

      <div className="rounded-xl border border-border bg-bg-secondary overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="border-b border-border">
              <th className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase">Task</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase">Schedule</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase">Type</th>
            </tr>
          </thead>
          <tbody>
            {isLoading ? (
              <tr>
                <td colSpan={3} className="px-4 py-12 text-center">
                  <Loader2 className="h-6 w-6 animate-spin text-accent mx-auto" />
                </td>
              </tr>
            ) : tasks && tasks.length > 0 ? (
              tasks.map((task, i) => (
                <tr key={i} className="border-b border-border last:border-0 hover:bg-bg-hover transition-colors">
                  <td className="px-4 py-3">
                    <span className="text-sm font-medium text-foreground">{task.name}</span>
                  </td>
                  <td className="px-4 py-3">
                    <code className="rounded-md bg-bg-tertiary px-2 py-1 text-xs font-mono text-accent">
                      {task.schedule}
                    </code>
                  </td>
                  <td className="px-4 py-3">
                    <span className="rounded-md bg-accent/10 px-2 py-0.5 text-xs font-medium text-accent">
                      {task.type}
                    </span>
                  </td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan={3} className="px-4 py-12 text-center">
                  <Calendar className="h-8 w-8 text-text-muted mx-auto mb-2" />
                  <p className="text-sm text-text-secondary">No cron tasks registered</p>
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
`
}

// adminMailPage is the Mail Preview. The templates and their rendering come
// from the API (GET /api/admin/mail/templates and /preview/:template), so what
// the page shows is what the mailer sends, including every template grit
// generate mail adds. It used to be JSX copies of four templates.
func adminMailPage() string {
	return `"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Mail, AlertTriangle } from "@/lib/icons";
import { apiClient } from "@/lib/api-client";

interface MailTemplate {
  name: string;
  description: string;
  subject: string;
  has_text: boolean;
}

interface MailTemplatesResponse {
  data: MailTemplate[];
  driver: string;
}

type Part = "html" | "text";

export default function MailPage() {
  const [selected, setSelected] = useState<string | null>(null);
  const [part, setPart] = useState<Part>("html");

  const templates = useQuery<MailTemplatesResponse>({
    queryKey: ["mail-templates"],
    queryFn: async () => {
      const { data } = await apiClient.get<MailTemplatesResponse>("/api/admin/mail/templates");
      return data;
    },
  });

  const list = templates.data?.data ?? [];
  const current = list.find((t) => t.name === selected) ?? list[0];
  const shownPart: Part = part === "text" && current?.has_text ? "text" : "html";

  // The HTML is rendered by the API with sample data and shown in a sandboxed
  // iframe, so the page shows exactly what the mailer sends.
  const preview = useQuery<string>({
    queryKey: ["mail-preview", current?.name ?? "", shownPart],
    enabled: Boolean(current),
    queryFn: async () => {
      const name = current ? current.name : "";
      const { data } = await apiClient.get<string>(
        "/api/admin/mail/preview/" + encodeURIComponent(name) + (shownPart === "text" ? "?part=text" : ""),
        { responseType: "text", transformResponse: [(body: string) => body] },
      );
      return data;
    },
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-foreground">Email Templates</h1>
        <p className="text-sm text-text-secondary mt-1">
          Rendered by the API with sample data.{" "}
          {templates.data && (templates.data.driver ? "Mail is sent with " + templates.data.driver + "." : "No mail driver is configured.")}
        </p>
      </div>

      {templates.isError && (
        <div className="rounded-xl border border-warning/30 bg-warning/5 p-4 text-sm text-warning flex items-start gap-2">
          <AlertTriangle className="h-4 w-4 mt-0.5 shrink-0" />
          <span>The templates could not be loaded from the API.</span>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        <div className="space-y-2">
          {list.map((t) => {
            const active = current?.name === t.name;
            return (
              <button
                key={t.name}
                type="button"
                onClick={() => setSelected(t.name)}
                aria-pressed={active}
                className={
                  "w-full text-left rounded-xl border p-4 transition-colors " +
                  (active ? "border-accent bg-accent/5" : "border-border bg-bg-secondary hover:border-accent/30")
                }
              >
                <div className="flex items-center gap-3">
                  <Mail className={"h-4 w-4 " + (active ? "text-accent" : "text-text-muted")} />
                  <div className="min-w-0">
                    <p className="text-sm font-medium text-foreground truncate">{t.name}</p>
                    <p className="text-xs text-text-muted mt-0.5">{t.description}</p>
                  </div>
                </div>
              </button>
            );
          })}
        </div>

        <div className="lg:col-span-3">
          {current && (
            <div className="rounded-xl border border-border bg-bg-secondary overflow-hidden">
              <div className="flex items-center justify-between gap-3 border-b border-border px-4 py-3">
                <div className="min-w-0">
                  <p className="text-sm font-medium text-foreground truncate">{current.subject || current.name}</p>
                  <p className="text-xs text-text-muted">Template: {current.name}</p>
                </div>
                <div className="flex items-center gap-1 rounded-lg border border-border p-0.5">
                  {(["html", "text"] as Part[]).map((p) => (
                    <button
                      key={p}
                      type="button"
                      disabled={p === "text" && !current.has_text}
                      onClick={() => setPart(p)}
                      aria-pressed={shownPart === p}
                      className={
                        "rounded-md px-2.5 py-1 text-xs font-medium disabled:opacity-40 " +
                        (shownPart === p ? "bg-accent/10 text-accent" : "text-text-secondary hover:text-foreground")
                      }
                    >
                      {p === "html" ? "HTML" : "Text"}
                    </button>
                  ))}
                </div>
              </div>
              <div className="p-4">
                {preview.isError ? (
                  <p className="text-sm text-warning">The preview could not be rendered.</p>
                ) : shownPart === "text" ? (
                  <pre className="text-xs text-text-secondary font-mono bg-bg-tertiary rounded-lg p-4 whitespace-pre-wrap">
                    {preview.data ?? ""}
                  </pre>
                ) : (
                  <iframe
                    title={"Preview of " + current.name}
                    sandbox=""
                    srcDoc={preview.data ?? ""}
                    className="w-full h-[640px] rounded-lg border border-border bg-white"
                  />
                )}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
`
}

func adminSecurityPage() string {
	return `"use client";

import { Shield, ExternalLink, AlertTriangle, Zap, Activity, AlertCircle, Globe } from "@/lib/icons";
import { useSecuritySummary } from "@/hooks/use-system";
import { getApiErrorMessage } from "@/lib/api-core";

import { API_URL } from "@/lib/api-core";

interface ScoreCard {
  score?: number
  grade?: string
  factors?: { name: string; score: number }[]
}

interface Summary {
  summary?: { threats_total?: number; threats_24h?: number; blocked?: number; actors?: number }
  score?: ScoreCard
  threats?: { data?: Array<{ id: string; type: string; severity: string; cvss?: number; source_ip?: string; route?: string; created_at?: string }> }
  auth_shield?: { enabled?: boolean; failures?: number; locked?: Array<{ ip: string; until: string }> }
  csp_top?: { data?: Array<{ violated_directive: string; blocked_uri: string; count: number }> }
  performance?: { p50_ms?: number; p95_ms?: number; p99_ms?: number; rps?: number; pipeline_drops?: number }
  trends?: unknown
  top_routes?: { data?: Array<{ route: string; count: number }> }
  _errors?: Record<string, string>
}

export default function SecurityPage() {
  // Polls every 20 seconds. The hook owns the interval, the in-flight request
  // and its cancellation, so this page has no effect of its own to get wrong.
  const query = useSecuritySummary<Summary>();
  const data = query.data ?? null;
  const loading = query.isPending;
  const err = query.error ? getApiErrorMessage(query.error, "Failed to load") : null;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground flex items-center gap-2.5">
            <Shield className="h-6 w-6 text-accent" /> Security
          </h1>
          <p className="text-sm text-text-secondary mt-1">
            Live summary from Sentinel — score, threats, AuthShield, CSP, performance
          </p>
        </div>
        <a
          href={` + "`${API_URL}/sentinel/ui`" + `}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-2 rounded-lg border border-border bg-bg-secondary px-4 py-2.5 text-sm font-medium text-foreground hover:bg-bg-hover transition-colors"
        >
          Open full dashboard <ExternalLink className="h-4 w-4" />
        </a>
      </div>

      {err && !loading && (
        <div className="rounded-xl border border-warning/30 bg-warning/5 p-4 text-sm text-warning flex items-start gap-2">
          <AlertTriangle className="h-4 w-4 mt-0.5 shrink-0" />
          <div>
            <strong>Couldn&apos;t reach Sentinel.</strong> Make sure <code className="px-1 py-0.5 rounded bg-bg-secondary text-xs font-mono">SENTINEL_ENABLED=true</code> and the API is running.
            <div className="text-xs text-text-muted mt-1">{err}</div>
          </div>
        </div>
      )}

      {/* Top-line scorecard */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        <KpiCard title="Security score" value={data?.score?.score ?? "—"} sub={data?.score?.grade} tone="success" icon={Shield} />
        <KpiCard title="Threats (24h)" value={data?.summary?.threats_24h ?? "—"} sub={` + "`${data?.summary?.threats_total ?? 0} total`" + `} tone="warning" icon={AlertCircle} />
        <KpiCard title="Blocked" value={data?.summary?.blocked ?? "—"} sub="auto by WAF" tone="info" icon={Zap} />
        <KpiCard title="Actors" value={data?.summary?.actors ?? "—"} sub="unique sources" tone="default" icon={Globe} />
      </div>

      {/* Recent threats */}
      <div className="rounded-xl border border-border bg-bg-secondary overflow-hidden">
        <div className="flex items-center justify-between px-4 py-3 border-b border-border">
          <h2 className="text-sm font-semibold text-foreground">Recent threats</h2>
          <a href={` + "`${API_URL}/sentinel/ui/threats`" + `} target="_blank" rel="noopener noreferrer" className="text-xs text-accent hover:underline">View all</a>
        </div>
        <div className="divide-y divide-border">
          {(data?.threats?.data ?? []).slice(0, 8).map((t) => (
            <div key={t.id} className="flex items-center gap-3 px-4 py-2.5 text-sm">
              <SeverityChip severity={t.severity} cvss={t.cvss} />
              <span className="text-foreground font-mono text-xs truncate flex-1">{t.type}</span>
              <span className="text-text-muted text-xs hidden md:inline">{t.route}</span>
              <span className="text-text-muted text-xs font-mono">{t.source_ip}</span>
            </div>
          ))}
          {(!data?.threats?.data || data.threats.data.length === 0) && (
            <div className="px-4 py-6 text-center text-sm text-text-muted">No threats in window.</div>
          )}
        </div>
      </div>

      {/* AuthShield + CSP side-by-side */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <Panel title="AuthShield">
          {data?.auth_shield?.enabled === false ? (
            <p className="text-sm text-text-muted">AuthShield is disabled in this environment.</p>
          ) : (
            <div className="space-y-3">
              <div className="grid grid-cols-2 gap-2 text-sm">
                <Stat label="Failures (rolling)" value={data?.auth_shield?.failures ?? 0} />
                <Stat label="Locked IPs" value={data?.auth_shield?.locked?.length ?? 0} />
              </div>
              <ul className="space-y-1.5 text-xs font-mono">
                {(data?.auth_shield?.locked ?? []).slice(0, 5).map((l) => (
                  <li key={l.ip} className="flex justify-between gap-2 text-text-secondary">
                    <span>{l.ip}</span>
                    <span className="text-text-muted">until {l.until}</span>
                  </li>
                ))}
              </ul>
            </div>
          )}
        </Panel>
        <Panel title="Top CSP violations">
          <ul className="space-y-1.5 text-xs font-mono">
            {(data?.csp_top?.data ?? []).slice(0, 5).map((v, i) => (
              <li key={i} className="flex justify-between gap-2">
                <span className="truncate text-text-secondary">{v.violated_directive}</span>
                <span className="text-text-muted shrink-0">{v.count}×</span>
              </li>
            ))}
            {(!data?.csp_top?.data || data.csp_top.data.length === 0) && (
              <li className="text-text-muted">No violations.</li>
            )}
          </ul>
        </Panel>
      </div>
    </div>
  );
}

function KpiCard({ title, value, sub, tone, icon: Icon }: { title: string; value: any; sub?: any; tone: "default" | "success" | "warning" | "info" | "danger"; icon: any }) {
  const toneCls = { default: "text-foreground", success: "text-success", warning: "text-warning", info: "text-info", danger: "text-danger" }[tone];
  return (
    <div className="rounded-xl border border-border bg-bg-secondary p-4">
      <div className="flex items-center justify-between mb-1">
        <span className="text-[10px] font-mono uppercase tracking-wider text-text-muted">{title}</span>
        <Icon className={` + "`h-4 w-4 ${toneCls}`" + `} />
      </div>
      <p className={` + "`text-2xl font-semibold tabular-nums ${toneCls}`" + `}>{value}</p>
      {sub && <p className="text-[11px] text-text-muted mt-0.5">{sub}</p>}
    </div>
  );
}

function Panel({ title, children }: { title: string; children: any }) {
  return (
    <div className="rounded-xl border border-border bg-bg-secondary overflow-hidden">
      <div className="px-4 py-3 border-b border-border"><h2 className="text-sm font-semibold text-foreground">{title}</h2></div>
      <div className="p-4">{children}</div>
    </div>
  );
}

function Stat({ label, value }: { label: string; value: any }) {
  return (
    <div className="rounded-lg border border-border bg-bg-elevated px-3 py-2">
      <p className="text-[10px] uppercase tracking-wider text-text-muted">{label}</p>
      <p className="text-lg font-semibold text-foreground tabular-nums">{value}</p>
    </div>
  );
}

function SeverityChip({ severity, cvss }: { severity?: string; cvss?: number }) {
  const tone =
    severity === "critical" ? "bg-danger/15 text-danger border-danger/30" :
    severity === "high"     ? "bg-warning/15 text-warning border-warning/30" :
    severity === "medium"   ? "bg-info/15 text-info border-info/30" :
                              "bg-bg-elevated text-text-muted border-border";
  return (
    <span className={` + "`inline-flex items-center gap-1 px-2 py-0.5 rounded border text-[10px] font-mono ${tone}`" + `}>
      {severity ?? "info"}
      {cvss !== undefined && cvss > 0 && <span>· {cvss.toFixed(1)}</span>}
    </span>
  );
}
`
}

// adminObservabilityPage emits app/(dashboard)/system/observability/page.tsx.
//
// The page it used to emit is merged into /system/performance. It stays as a
// redirect rather than a deleted route because it was linked from the System
// Hub, from release notes and from whatever anyone bookmarked, and a 404 is a
// worse answer than the page they were looking for.
func adminObservabilityPage() string {
	return `"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

// Merged into the operations page: everything this showed, and the latency,
// throughput and slow routes it never managed to.
export default function ObservabilityRedirect() {
	const router = useRouter();
	useEffect(() => {
		router.replace("/system/performance");
	}, [router]);

	return (
		<p className="p-8 text-sm text-text-secondary">
			Observability has moved to <a className="underline" href="/system/performance">Operations</a>.
		</p>
	);
}
`
}
