import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
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
      const { data } = await apiClient.get(`/api/admin/jobs/${status}?queue=${queue}`);
      return data.data;
    },
    refetchInterval: 5000,
  });
}

export function useRetryJob() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ id, queue }: { id: string; queue?: string }) => {
      await apiClient.post(`/api/admin/jobs/${id}/retry?queue=${queue || "default"}`);
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
      await apiClient.delete(`/api/admin/jobs/queue/${queue}`);
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
      const { data } = await apiClient.get(`/api/uploads?page=${page}&page_size=${pageSize}`);
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
      await apiClient.delete(`/api/uploads/${id}`);
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
