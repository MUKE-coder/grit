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

export type BackupFrequency = "daily" | "weekly" | "monthly" | "yearly";

export interface BackupSchedule {
  frequency: BackupFrequency;
  time: string; // "HH:MM"
  enabled: boolean;
}

// The automatic-backup schedule (period + time-of-day). Changing it takes effect
// on the scheduler's next tick — no restart needed.
export function useBackupSchedule() {
  return useQuery<BackupSchedule>({
    queryKey: ["backup-schedule"],
    queryFn: async () => (await apiClient.get("/api/backup-settings")).data.data,
    staleTime: 15_000,
  });
}

export function useUpdateBackupSchedule() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (s: BackupSchedule) => (await apiClient.put("/api/backup-settings", s)).data.data as BackupSchedule,
    onSuccess: (data) => qc.setQueryData(["backup-schedule"], data),
  });
}
`
}

// adminBackupsPage is app/(dashboard)/system/backups/page.tsx.
func adminBackupsPage() string {
	return `"use client";

import { useEffect, useState } from "react";
import {
  useBackups, useGenerateBackup, useDownloadBackup, useBackupSchedule, useUpdateBackupSchedule,
  type Backup, type BackupFrequency,
} from "@/hooks/use-backups";
import { Database, Download, RefreshCw, Loader2, AlertCircle, Clock, Save } from "@/lib/icons";
import { buttonClasses } from "@/components/ui/button";
import { inputClasses } from "@/components/ui/input";

const FREQUENCIES: { key: BackupFrequency; label: string; hint: string }[] = [
  { key: "daily", label: "Daily", hint: "every day" },
  { key: "weekly", label: "Weekly", hint: "every Sunday" },
  { key: "monthly", label: "Monthly", hint: "1st of month" },
  { key: "yearly", label: "Yearly", hint: "Jan 1" },
];

function ScheduleCard() {
  const { data: schedule, isLoading } = useBackupSchedule();
  const update = useUpdateBackupSchedule();
  const [freq, setFreq] = useState<BackupFrequency>("weekly");
  const [time, setTime] = useState("02:00");
  const [enabled, setEnabled] = useState(true);

  useEffect(() => {
    if (schedule) { setFreq(schedule.frequency); setTime(schedule.time); setEnabled(schedule.enabled); }
  }, [schedule]);

  const dirty = schedule && (schedule.frequency !== freq || schedule.time !== time || schedule.enabled !== enabled);

  return (
    <div className="rounded-xl border border-border bg-bg-secondary p-5">
      <div className="flex items-start justify-between gap-4">
        <div className="flex items-center gap-3">
          <span className="inline-flex h-9 w-9 items-center justify-center rounded-lg bg-accent/10 text-accent"><Clock className="h-5 w-5" /></span>
          <div>
            <h3 className="text-[15px] font-semibold text-foreground">Automatic backups</h3>
            <p className="text-sm text-text-secondary">A full backup runs on this schedule. Default is weekly.</p>
          </div>
        </div>
        <button
          type="button"
          onClick={() => setEnabled((v) => !v)}
          role="switch"
          aria-checked={enabled}
          className={"relative mt-1 inline-flex h-6 w-11 shrink-0 items-center rounded-full transition-colors " + (enabled ? "bg-accent" : "bg-text-muted/40")}
        >
          <span className={"inline-block h-4 w-4 rounded-full bg-white transition-transform " + (enabled ? "translate-x-6" : "translate-x-1")} />
        </button>
      </div>

      <div className={"mt-4 transition-opacity " + (enabled ? "" : "pointer-events-none opacity-50")}>
        <p className="mb-2 text-xs font-semibold uppercase tracking-wider text-text-muted">Frequency</p>
        <div className="grid grid-cols-2 gap-2 sm:grid-cols-4">
          {FREQUENCIES.map((f) => (
            <button
              key={f.key}
              onClick={() => setFreq(f.key)}
              className={
                "rounded-lg border px-3 py-2.5 text-center transition-colors " +
                (freq === f.key ? "border-accent bg-accent/5 text-accent" : "border-border text-text-secondary hover:bg-bg-hover")
              }
            >
              <div className="text-sm font-medium text-foreground">{f.label}</div>
              <div className="text-xs text-text-muted">{f.hint}</div>
            </button>
          ))}
        </div>
        <div className="mt-4 flex items-end gap-3">
          <div>
            <label className="mb-1.5 block text-xs font-semibold uppercase tracking-wider text-text-muted">Time of day</label>
            <input type="time" value={time} onChange={(e) => setTime(e.target.value)} className={inputClasses()} />
          </div>
          <span className="pb-2 text-xs text-text-muted">server local time</span>
        </div>
      </div>

      <div className="mt-5 flex items-center justify-end gap-3">
        {update.isSuccess && !dirty ? <span className="text-xs text-success">Saved</span> : null}
        <button
          onClick={() => update.mutate({ frequency: freq, time, enabled })}
          disabled={isLoading || update.isPending || !dirty}
          className={buttonClasses()}
        >
          {update.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />} Save schedule
        </button>
      </div>
    </div>
  );
}

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
          <h1 className="text-2xl font-bold text-foreground">Data &amp; Backup</h1>
          <p className="text-sm text-text-secondary mt-1">
            Full-database snapshots on a schedule you control, or on demand. The four most recent are kept.
          </p>
        </div>
        <button
          onClick={() => generate.mutate()}
          disabled={generate.isPending || running}
          className={buttonClasses()}
        >
          {generate.isPending || running ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <RefreshCw className="h-4 w-4" />
          )}
          Back up now
        </button>
      </div>

      <ScheduleCard />

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
