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

func adminObservabilityPage() string {
	return `"use client";

import type { ReactNode } from "react";
import type { LucideIcon } from "lucide-react";
import { Activity, ExternalLink, AlertTriangle, Zap, CheckCircle, Circle, CircleDashed, Minus, Gauge } from "@/lib/icons";
import { useObservabilitySummary, useScaleReport } from "@/hooks/use-system";
import { getApiErrorMessage } from "@/lib/api-core";

import { API_URL } from "@/lib/api-core";

interface Summary {
  overview?: { p50_ms?: number; p95_ms?: number; p99_ms?: number; rps?: number; error_rate?: number; total_requests?: number }
  slos?: { data?: Array<{ name: string; target: number; current: number; budget_remaining: number; status: string }> }
  use?: { resources?: Array<{ name: string; utilization: { value: number; band: string }; saturation: { value: number; band: string }; errors: { value: number; band: string } }> }
  n1_ranked?: { data?: Array<{ route: string; pattern: string; occurrences: number; avg_queries_per_request: number; impact_score: number }> }
  errors?: { data?: Array<{ id: string; type: string; message: string; route: string; count: number }> }
  runtime?: { heap_alloc_mb?: number; goroutines?: number; gc_pause_ms?: number }
  health_checks?: { data?: Array<{ name: string; status: string; latency_ms: number; error?: string }> }
  alerts?: { data?: Array<{ id: string; name: string; severity: string; message: string }> }
  _errors?: Record<string, string>
}

type StageState = "on" | "ready" | "attention" | "yours";

interface ScaleReport {
  stage?: string
  next?: string
  why?: string
  notes?: string[]
  stages?: Array<{ n: string; name: string; state: StageState; what: string; evidence: string }>
  latency?: { samples?: number; p50_ms?: number; p95_ms?: number; p99_ms?: number }
  database?: { replicas?: number; open_connections?: number; max_connections?: number }
  cache?: { configured?: boolean; hit_rate?: number }
}

// How each state reads. "ready" is a tick too, hollow rather than filled: the
// stage is shipped and switched off because nothing yet needs it, which is a
// good answer and not a missing feature.
const STAGE_STATE: Record<StageState, { icon: typeof CheckCircle; cls: string; label: string }> = {
  on:        { icon: CheckCircle,  cls: "text-success",    label: "In force" },
  ready:     { icon: Circle,       cls: "text-success/60", label: "Ready, not needed yet" },
  attention: { icon: AlertTriangle, cls: "text-warning",   label: "Needs attention" },
  yours:     { icon: Minus,        cls: "text-text-muted", label: "Yours to decide" },
};

const BAND_CLS: Record<string, string> = {
  green:   "bg-success/15 text-success border-success/30",
  amber:   "bg-warning/15 text-warning border-warning/30",
  red:     "bg-danger/15 text-danger border-danger/30",
  unknown: "bg-bg-elevated text-text-muted border-border",
};

export default function ObservabilityPage() {
  // Polls every 10 seconds. The local helper this replaced was called "fetch",
  // which shadowed the global one for the rest of the effect: any code added
  // there that reached for fetch() would have called the poller instead.
  const query = useObservabilitySummary<Summary>();
  const data = query.data ?? null;
  const loading = query.isPending;
  const err = query.error ? getApiErrorMessage(query.error, "Failed to load") : null;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground flex items-center gap-2.5">
            <Activity className="h-6 w-6 text-accent" /> Observability
          </h1>
          <p className="text-sm text-text-secondary mt-1">
            Live summary from Pulse — percentile latency, SLOs, USE grid, top N+1, errors, runtime
          </p>
        </div>
        <a
          href={` + "`${API_URL}/pulse/ui`" + `}
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
            <strong>Couldn&apos;t reach Pulse.</strong> Make sure <code className="px-1 py-0.5 rounded bg-bg-secondary text-xs font-mono">PULSE_ENABLED=true</code> and the API is running.
            <div className="text-xs text-text-muted mt-1">{err}</div>
          </div>
        </div>
      )}

      <ScalingReadiness />

      {/* Top-line KPIs */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        <Kpi label="p95 latency" value={data?.overview?.p95_ms != null ? ` + "`${data.overview.p95_ms.toFixed(0)} ms`" + ` : "—"} tone="info" icon={Zap} />
        <Kpi label="p99 latency" value={data?.overview?.p99_ms != null ? ` + "`${data.overview.p99_ms.toFixed(0)} ms`" + ` : "—"} tone="warning" icon={Zap} />
        <Kpi label="Error rate" value={data?.overview?.error_rate != null ? ` + "`${(data.overview.error_rate * 100).toFixed(2)}%`" + ` : "—"} tone="danger" icon={AlertTriangle} />
        <Kpi label="RPS" value={data?.overview?.rps?.toFixed?.(1) ?? "—"} tone="success" icon={Activity} />
      </div>

      {/* SLO compliance bars */}
      <Panel title="SLOs">
        {(data?.slos?.data ?? []).length === 0 ? (
          <p className="text-sm text-text-muted">No SLOs configured. Add to <code className="px-1 py-0.5 rounded bg-bg-elevated text-xs font-mono">pulse.Config.SLOs</code>.</p>
        ) : (
          <div className="space-y-3">
            {(data?.slos?.data ?? []).map((s) => (
              <div key={s.name} className="space-y-1">
                <div className="flex items-center justify-between text-xs">
                  <span className="font-medium text-foreground">{s.name}</span>
                  <span className={` + "`font-mono ${s.status === \"firing\" ? \"text-danger\" : \"text-success\"}`" + `}>{(s.current * 100).toFixed(2)}% / {(s.target * 100).toFixed(2)}%</span>
                </div>
                <div className="h-2 rounded-full bg-bg-elevated overflow-hidden">
                  <div className={` + "`h-full ${s.status === \"firing\" ? \"bg-danger\" : \"bg-success\"}`" + `} style={{ width: ` + "`${Math.min(s.current * 100, 100)}%`" + ` }} />
                </div>
                <p className="text-[10px] text-text-muted">Budget remaining: {(s.budget_remaining * 100).toFixed(1)}%</p>
              </div>
            ))}
          </div>
        )}
      </Panel>

      {/* USE grid + N+1 side-by-side */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <Panel title="USE method">
          <table className="w-full text-xs">
            <thead>
              <tr className="text-text-muted">
                <th className="text-left font-medium py-1">Resource</th>
                <th className="text-center font-medium">U</th>
                <th className="text-center font-medium">S</th>
                <th className="text-center font-medium">E</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {(data?.use?.resources ?? []).map((r) => (
                <tr key={r.name}>
                  <td className="py-1.5 font-medium text-foreground capitalize">{r.name}</td>
                  <td className="py-1.5 text-center"><Cell band={r.utilization?.band} value={r.utilization?.value} /></td>
                  <td className="py-1.5 text-center"><Cell band={r.saturation?.band} value={r.saturation?.value} /></td>
                  <td className="py-1.5 text-center"><Cell band={r.errors?.band} value={r.errors?.value} /></td>
                </tr>
              ))}
              {(!data?.use?.resources || data.use.resources.length === 0) && (
                <tr><td colSpan={4} className="text-center text-text-muted py-3">No USE samples yet.</td></tr>
              )}
            </tbody>
          </table>
        </Panel>

        <Panel title="Top N+1 by impact">
          <ul className="space-y-2 text-xs">
            {(data?.n1_ranked?.data ?? []).slice(0, 6).map((n, i) => (
              <li key={i} className="rounded-lg border border-border bg-bg-elevated p-2.5">
                <div className="flex items-center justify-between mb-0.5">
                  <span className="font-mono text-foreground truncate">{n.route}</span>
                  <span className="text-text-muted shrink-0 ml-2">impact {n.impact_score.toFixed(0)}</span>
                </div>
                <p className="font-mono text-[10px] text-text-muted truncate">{n.pattern}</p>
                <p className="text-[10px] text-text-muted mt-0.5">{n.occurrences} occurrences · ~{n.avg_queries_per_request} queries/req</p>
              </li>
            ))}
            {(!data?.n1_ranked?.data || data.n1_ranked.data.length === 0) && (
              <li className="text-text-muted text-center py-3">No N+1 detections.</li>
            )}
          </ul>
        </Panel>
      </div>

      {/* Errors + Runtime side-by-side */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <Panel title="Recent unresolved errors">
          <ul className="space-y-1.5 text-xs">
            {(data?.errors?.data ?? []).slice(0, 6).map((e) => (
              <li key={e.id} className="flex items-start justify-between gap-2 py-1 border-b border-border last:border-0">
                <div className="min-w-0 flex-1">
                  <p className="text-foreground truncate"><span className="font-mono text-[10px] text-text-muted">{e.type}</span> {e.message}</p>
                  <p className="font-mono text-[10px] text-text-muted truncate">{e.route}</p>
                </div>
                <span className="text-text-muted shrink-0">{e.count}×</span>
              </li>
            ))}
            {(!data?.errors?.data || data.errors.data.length === 0) && <li className="text-text-muted text-center py-3">No unresolved errors.</li>}
          </ul>
        </Panel>

        <Panel title="Go runtime">
          <div className="grid grid-cols-3 gap-2 text-xs">
            <Stat label="Heap alloc" value={data?.runtime?.heap_alloc_mb != null ? ` + "`${data.runtime.heap_alloc_mb.toFixed(0)} MB`" + ` : "—"} />
            <Stat label="Goroutines" value={data?.runtime?.goroutines ?? "—"} />
            <Stat label="GC pause" value={data?.runtime?.gc_pause_ms != null ? ` + "`${data.runtime.gc_pause_ms.toFixed(1)} ms`" + ` : "—"} />
          </div>
          {data?.health_checks?.data && data.health_checks.data.length > 0 && (
            <div className="mt-4 pt-3 border-t border-border space-y-1">
              <p className="text-[10px] uppercase font-mono tracking-wider text-text-muted">Health checks</p>
              {data.health_checks.data.map((h) => (
                <div key={h.name} className="flex items-center justify-between text-xs">
                  <span className="font-mono text-foreground">{h.name}</span>
                  <span className={` + "`text-[10px] ${h.status === \"healthy\" ? \"text-success\" : \"text-danger\"}`" + `}>
                    {h.status} · {h.latency_ms}ms
                  </span>
                </div>
              ))}
            </div>
          )}
        </Panel>
      </div>
    </div>
  );
}

// Which of the ten scaling stages this deployment is past, measured rather than
// claimed. The verdict names the one thing to do next, and most of the time
// that is nothing; the list below it is the half people ask about first, which
// is what is already handled.
function ScalingReadiness() {
  const query = useScaleReport<ScaleReport>();
  const report = query.data ?? null;
  const stages = report?.stages ?? [];
  const attention = stages.filter((s) => s.state === "attention").length;
  const covered = stages.filter((s) => s.state === "on" || s.state === "ready").length;

  if (query.isError) {
    return (
      <Panel title="Scaling readiness">
        <p className="text-sm text-text-muted">
          {getApiErrorMessage(query.error, "Could not read the scaling report")}
        </p>
      </Panel>
    );
  }

  return (
    <div className="rounded-xl border border-border bg-bg-secondary overflow-hidden">
      <div className="px-4 py-3 border-b border-border flex items-center justify-between gap-3">
        <h2 className="text-sm font-semibold text-foreground flex items-center gap-2">
          <Gauge className="h-4 w-4 text-accent" /> Scaling readiness
        </h2>
        {stages.length > 0 && (
          <span className={` + "`${attention > 0 ? \"text-warning\" : \"text-success\"} text-xs font-mono`" + `}>
            {covered} of {stages.length} stages handled{attention > 0 ? ` + "`, ${attention} need${attention === 1 ? \"s\" : \"\"} attention`" + ` : ""}
          </span>
        )}
      </div>

      {/* The verdict, first, because it is the answer to "what now". */}
      <div className="px-4 py-3 border-b border-border bg-bg-elevated/40">
        {query.isPending ? (
          <p className="text-sm text-text-muted">Measuring…</p>
        ) : (
          <>
            <p className="text-sm font-medium text-foreground">{report?.stage ?? "—"}</p>
            {report?.why && <p className="text-xs text-text-secondary mt-0.5">{report.why}</p>}
            {report?.next && (
              <p className="text-xs text-text-muted mt-1.5">
                <span className="font-mono uppercase tracking-wider text-[10px] text-text-muted">Do this next</span>{" "}
                <span className="text-foreground">{report.next}</span>
              </p>
            )}
          </>
        )}
      </div>

      <ul className="divide-y divide-border">
        {stages.map((stage) => {
          const state = STAGE_STATE[stage.state] ?? STAGE_STATE.yours;
          const Icon = state.icon;
          return (
            <li key={stage.n} className="flex items-start gap-3 px-4 py-2.5">
              <Icon className={` + "`h-4 w-4 mt-0.5 shrink-0 ${state.cls}`" + `} aria-hidden />
              <div className="min-w-0 flex-1">
                <div className="flex items-baseline gap-2">
                  <span className="font-mono text-[10px] text-text-muted shrink-0">{stage.n}</span>
                  <span className="text-sm font-medium text-foreground">{stage.name}</span>
                </div>
                <p className="text-xs text-text-secondary mt-0.5">{stage.what}</p>
                {stage.evidence && (
                  <p className="text-[11px] text-text-muted mt-0.5 font-mono">{stage.evidence}</p>
                )}
              </div>
              <span className={` + "`text-[10px] shrink-0 mt-1 ${state.cls}`" + `}>{state.label}</span>
            </li>
          );
        })}
        {!query.isPending && stages.length === 0 && (
          <li className="px-4 py-6 text-center text-sm text-text-muted">
            No stage report. This needs an API on v3.336.0 or later.
          </li>
        )}
      </ul>

      {(report?.notes ?? []).length > 0 && (
        <div className="px-4 py-3 border-t border-border space-y-1">
          {(report?.notes ?? []).map((note, i) => (
            <p key={i} className="text-[11px] text-text-muted flex items-start gap-1.5">
              <CircleDashed className="h-3 w-3 mt-0.5 shrink-0" aria-hidden />
              {note}
            </p>
          ))}
        </div>
      )}
    </div>
  );
}

function Kpi({ label, value, tone, icon: Icon }: { label: string; value: ReactNode; tone: "default" | "success" | "warning" | "info" | "danger"; icon: LucideIcon }) {
  const toneCls = { default: "text-foreground", success: "text-success", warning: "text-warning", info: "text-info", danger: "text-danger" }[tone];
  return (
    <div className="rounded-xl border border-border bg-bg-secondary p-4">
      <div className="flex items-center justify-between mb-1">
        <span className="text-[10px] font-mono uppercase tracking-wider text-text-muted">{label}</span>
        <Icon className={` + "`h-4 w-4 ${toneCls}`" + `} />
      </div>
      <p className={` + "`text-2xl font-semibold tabular-nums ${toneCls}`" + `}>{value}</p>
    </div>
  );
}

function Panel({ title, children }: { title: string; children: ReactNode }) {
  return (
    <div className="rounded-xl border border-border bg-bg-secondary overflow-hidden">
      <div className="px-4 py-3 border-b border-border"><h2 className="text-sm font-semibold text-foreground">{title}</h2></div>
      <div className="p-4">{children}</div>
    </div>
  );
}

function Stat({ label, value }: { label: string; value: ReactNode }) {
  return (
    <div className="rounded-lg border border-border bg-bg-elevated px-3 py-2">
      <p className="text-[10px] uppercase tracking-wider text-text-muted">{label}</p>
      <p className="text-base font-semibold text-foreground tabular-nums">{value}</p>
    </div>
  );
}

function Cell({ band, value }: { band?: string; value?: number }) {
  const cls = BAND_CLS[band || "unknown"];
  return (
    <span className={` + "`inline-flex items-center justify-center min-w-[44px] px-2 py-0.5 rounded border text-[10px] font-mono ${cls}`" + `}>
      {value != null ? value.toFixed(0) : "—"}
    </span>
  );
}
`
}
