package scaffold

// Admin panel surface for full-database backups: a React Query hook that polls
// only while a backup is RUNNING, and a page to list / generate / download.
//
// Download never proxies the archive through the API — the handler mints a
// short-lived pre-signed URL and the browser pulls straight from object storage,
// so a multi-hundred-MB file can't time out a request.

// adminUseBackups is hooks/use-backups.ts.
func adminUseBackups() string {
	return `import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

export interface Backup {
  id: string;
  kind: "WEEKLY" | "MANUAL" | "CLI";
  status: "RUNNING" | "READY" | "FAILED" | "PURGED";
  size_bytes: number;
  table_count: number;
  row_count: number;
  error?: string;
  created_at: string;
  completed_at?: string;
}

// Poll every 3s while any backup is RUNNING, then go idle. staleTime keeps the
// list from refetching on every window focus.
export function useBackups() {
  return useQuery<Backup[]>({
    queryKey: ["backups"],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/backups");
      return data.data ?? [];
    },
    staleTime: 15_000,
    refetchInterval: (query) =>
      (query.state.data ?? []).some((b) => b.status === "RUNNING") ? 3000 : false,
  });
}

// Manual backups are rate-limited to one per 24h server-side; a 429 surfaces here.
export function useGenerateBackup() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async () => {
      const { data } = await apiClient.post("/api/backups/generate");
      return data.data as Backup;
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ["backups"] }),
  });
}

// Mints a 15-minute pre-signed URL, then lets the browser download from storage.
export function useDownloadBackup() {
  return useMutation({
    mutationFn: async (id: string) => {
      const { data } = await apiClient.get(` + "`" + `/api/backups/${id}/download` + "`" + `);
      return data.data.url as string;
    },
    onSuccess: (url) => {
      window.location.assign(url);
    },
  });
}
`
}

// adminBackupsPage is app/(dashboard)/system/backups/page.tsx.
func adminBackupsPage() string {
	return `"use client";

import { useBackups, useGenerateBackup, useDownloadBackup, type Backup } from "@/hooks/use-backups";
import { Database, Download, RefreshCw, Loader2, AlertCircle } from "@/lib/icons";

function formatBytes(n: number): string {
  if (!n) return "—";
  if (n < 1024) return n + " B";
  if (n < 1024 * 1024) return (n / 1024).toFixed(1) + " KB";
  return (n / (1024 * 1024)).toFixed(1) + " MB";
}

function StatusBadge({ status }: { status: Backup["status"] }) {
  const styles: Record<Backup["status"], string> = {
    RUNNING: "bg-info/15 text-info",
    READY: "bg-success/15 text-success",
    FAILED: "bg-danger/15 text-danger",
    PURGED: "bg-text-muted/15 text-text-muted",
  };
  return (
    <span className={"inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium " + styles[status]}>
      {status === "RUNNING" ? <Loader2 className="h-3 w-3 animate-spin" /> : null}
      {status}
    </span>
  );
}

export default function BackupsPage() {
  const { data: backups, isLoading } = useBackups();
  const generate = useGenerateBackup();
  const download = useDownloadBackup();

  const running = (backups ?? []).some((b) => b.status === "RUNNING");

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Database Backups</h1>
          <p className="text-sm text-text-secondary mt-1">
            A full backup runs automatically every Sunday at 02:00 UTC. The four most recent are kept.
          </p>
        </div>
        <button
          onClick={() => generate.mutate()}
          disabled={generate.isPending || running}
          className="inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white disabled:opacity-50"
        >
          {generate.isPending || running ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <RefreshCw className="h-4 w-4" />
          )}
          Back up now
        </button>
      </div>

      {generate.isError ? (
        <div className="flex items-center gap-2 rounded-lg border border-danger/25 bg-danger/10 px-4 py-3 text-sm text-danger">
          <AlertCircle className="h-4 w-4 shrink-0" />
          {(generate.error as any)?.response?.data?.error?.message ?? "Failed to start the backup"}
        </div>
      ) : null}

      <div className="rounded-xl border border-border bg-bg-secondary overflow-hidden">
        {isLoading ? (
          <div className="flex items-center justify-center py-16">
            <Loader2 className="h-5 w-5 animate-spin text-text-muted" />
          </div>
        ) : !backups?.length ? (
          <div className="flex flex-col items-center justify-center py-16 text-center">
            <Database className="h-8 w-8 text-text-muted" />
            <p className="mt-3 text-sm text-text-secondary">No backups yet</p>
            <p className="mt-1 text-xs text-text-muted">
              The first one lands on Sunday, or take one now.
            </p>
          </div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border text-left text-xs font-medium text-text-secondary">
                <th className="px-4 py-3">Date</th>
                <th className="px-4 py-3">Kind</th>
                <th className="px-4 py-3">Status</th>
                <th className="px-4 py-3">Tables</th>
                <th className="px-4 py-3">Rows</th>
                <th className="px-4 py-3">Size</th>
                <th className="px-4 py-3"></th>
              </tr>
            </thead>
            <tbody>
              {backups.map((b) => (
                <tr key={b.id} className="border-b border-border last:border-0">
                  <td className="px-4 py-3 text-foreground">
                    {new Date(b.created_at).toLocaleString()}
                  </td>
                  <td className="px-4 py-3 text-text-secondary">{b.kind}</td>
                  <td className="px-4 py-3">
                    <StatusBadge status={b.status} />
                    {b.status === "FAILED" && b.error ? (
                      <p className="mt-1 max-w-md truncate text-xs text-danger" title={b.error}>
                        {b.error}
                      </p>
                    ) : null}
                  </td>
                  <td className="px-4 py-3 text-text-secondary">{b.table_count || "—"}</td>
                  <td className="px-4 py-3 text-text-secondary">
                    {b.row_count ? b.row_count.toLocaleString() : "—"}
                  </td>
                  <td className="px-4 py-3 text-text-secondary">{formatBytes(b.size_bytes)}</td>
                  <td className="px-4 py-3 text-right">
                    {b.status === "READY" ? (
                      <button
                        onClick={() => download.mutate(b.id)}
                        disabled={download.isPending}
                        className="inline-flex items-center gap-1.5 rounded-lg border border-border px-3 py-1.5 text-xs font-medium text-foreground hover:bg-bg-hover disabled:opacity-50"
                      >
                        <Download className="h-3.5 w-3.5" />
                        Download
                      </button>
                    ) : null}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      <p className="text-xs text-text-muted">
        Each archive is a ZIP: one CSV per table, a <code>dump.sql</code> of INSERTs, and a{" "}
        <code>metadata.json</code> manifest. Restore with{" "}
        <code>grit restore backup.zip</code> — test it before you need it.
      </p>
    </div>
  );
}
`
}
