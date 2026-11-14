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

/**
 * How many archives rolling retention keeps.
 *
 * It is backup.KeepBackups in Go and it rides along with the schedule, so the
 * "N of 4" on the page is the server's number rather than a second copy of it
 * that nobody updates.
 */
export function useBackupRetention() {
  return useQuery<number>({
    queryKey: ["backup-retention"],
    queryFn: async () => (await apiClient.get("/api/backup-settings")).data.meta?.keep ?? 0,
    staleTime: 60_000,
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

import { useEffect, useMemo, useState } from "react";
import {
  useBackups, useGenerateBackup, useDownloadBackup, useBackupSchedule, useUpdateBackupSchedule,
  useBackupRetention,
  type Backup, type BackupFrequency,
} from "@/hooks/use-backups";
import { Database, Download, RefreshCw, Loader2, AlertCircle, Clock, Save, Package, Copy } from "@/lib/icons";
import { Button } from "@/components/ui/button";
import { inputClasses } from "@/components/ui/input";
import { Pill, type PillTone } from "@/components/ui/pill";
import { Note } from "@/components/ui/note";
import { EmptyState } from "@/components/ui/empty-state";
import { RadioCard } from "@/components/ui/radio-card";
import { PageHeader } from "@/components/chrome/PageHeader";
import { getApiErrorMessage } from "@/lib/api-core";
import { toast } from "sonner";

const FREQUENCIES: { key: BackupFrequency; label: string; hint: string }[] = [
  { key: "daily", label: "Daily", hint: "Every day" },
  { key: "weekly", label: "Weekly", hint: "Every Sunday" },
  { key: "monthly", label: "Monthly", hint: "The 1st of the month" },
  { key: "yearly", label: "Yearly", hint: "The 1st of January" },
];

/**
 * When the next automatic backup lands, for the choice on screen right now.
 *
 * The schedule was three controls and no answer: you picked Monthly, typed a
 * time, and the page told you nothing about what you had just arranged. The
 * sentence is the point of the card, and it updates before you save so you can
 * see you got it wrong while it is still free to change.
 */
function nextRun(frequency: BackupFrequency, time: string): string {
  const [h, m] = time.split(":").map((n) => parseInt(n, 10));
  if (Number.isNaN(h) || Number.isNaN(m)) return "";
  const at = new Date();
  at.setHours(h, m, 0, 0);

  const advance = () => {
    if (frequency === "daily") at.setDate(at.getDate() + 1);
    else if (frequency === "weekly") at.setDate(at.getDate() + ((7 - at.getDay()) % 7 || 7));
    else if (frequency === "monthly") at.setMonth(at.getMonth() + 1, 1);
    else at.setFullYear(at.getFullYear() + 1, 0, 1);
  };

  if (frequency === "weekly" && at.getDay() !== 0) at.setDate(at.getDate() + ((7 - at.getDay()) % 7));
  else if (frequency === "monthly" && at.getDate() !== 1) at.setMonth(at.getMonth() + 1, 1);
  else if (frequency === "yearly" && (at.getMonth() !== 0 || at.getDate() !== 1))
    at.setFullYear(at.getFullYear() + (at.getMonth() === 0 && at.getDate() === 1 ? 0 : 1), 0, 1);

  // Today at that time has already gone by, so it is the next one.
  if (at.getTime() <= Date.now()) advance();

  return at.toLocaleString(undefined, {
    weekday: "long",
    day: "numeric",
    month: "long",
    hour: "2-digit",
    minute: "2-digit",
  });
}

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

  const when = useMemo(() => nextRun(freq, time), [freq, time]);

  return (
    <section className="flex flex-col rounded-2xl border border-border bg-bg-elevated">
      <header className="flex items-start justify-between gap-4 border-b border-border px-5 py-4">
        <div className="flex items-center gap-3">
          <span
            className="inline-flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-accent-soft text-accent"
            aria-hidden="true"
          >
            <Clock className="h-4 w-4" />
          </span>
          <div>
            <h2 className="text-[15px] font-semibold text-foreground">Schedule</h2>
            <p className="text-[13px] text-text-secondary">A full backup runs without anybody asking.</p>
          </div>
        </div>
        <button
          type="button"
          onClick={() => setEnabled((v) => !v)}
          role="switch"
          aria-checked={enabled}
          aria-label="Automatic backups"
          className={
            "relative mt-0.5 inline-flex h-6 w-11 shrink-0 items-center rounded-full transition-grit " +
            (enabled ? "bg-accent" : "bg-border-strong")
          }
        >
          <span
            className={
              "inline-block h-4 w-4 rounded-full bg-white transition-grit " +
              (enabled ? "translate-x-6" : "translate-x-1")
            }
          />
        </button>
      </header>

      <div className="flex-1 px-5 py-5">
        {enabled ? (
          <>
            <fieldset>
              <legend className="mb-2 text-xs font-medium uppercase tracking-wider text-text-muted">
                How often
              </legend>
              <div className="grid gap-2 sm:grid-cols-2">
                {FREQUENCIES.map((f) => (
                  <RadioCard
                    key={f.key}
                    name="backup-frequency"
                    value={f.key}
                    checked={freq === f.key}
                    onChange={(v) => setFreq(v as BackupFrequency)}
                    title={f.label}
                    description={f.hint}
                  />
                ))}
              </div>
            </fieldset>

            <div className="mt-4 flex flex-wrap items-end gap-3">
              <div>
                <label
                  htmlFor="backup-time"
                  className="mb-1.5 block text-xs font-medium uppercase tracking-wider text-text-muted"
                >
                  Time of day
                </label>
                <input
                  id="backup-time"
                  type="time"
                  value={time}
                  onChange={(e) => setTime(e.target.value)}
                  className={inputClasses({ fullWidth: false }) + " w-36"}
                />
              </div>
              <p className="pb-2.5 text-[13px] text-text-muted">
                Server local time, not yours.
              </p>
            </div>

            {/* What you have just arranged, before you commit to it. */}
            {when && (
              <p className="mt-4 rounded-lg border border-border bg-bg-tertiary px-3 py-2 text-[13px] text-text-secondary">
                Next run <span className="font-medium text-foreground">{when}</span>
              </p>
            )}
          </>
        ) : (
          <p className="text-[13px] text-text-secondary">
            Nothing runs on its own. Every snapshot has to be taken by hand from the button above,
            which is a thing to remember on the day you need one.
          </p>
        )}
      </div>

      <footer className="flex items-center justify-end gap-3 border-t border-border px-5 py-4">
        {update.isSuccess && !dirty ? (
          <Pill tone="good" dot>
            Saved
          </Pill>
        ) : null}
        <Button
          onClick={() => update.mutate({ frequency: freq, time, enabled })}
          disabled={isLoading || !dirty}
          loading={update.isPending}
        >
          {!update.isPending && <Save className="h-3.5 w-3.5" />}
          Save schedule
        </Button>
      </footer>
    </section>
  );
}

function formatBytes(n: number): string {
  if (!n) return "—";
  if (n < 1024) return n + " B";
  if (n < 1024 * 1024) return (n / 1024).toFixed(1) + " KB";
  return (n / (1024 * 1024)).toFixed(1) + " MB";
}

// PURGED is the one that needs a word rather than a label.
//
// It means rolling retention deleted the archive and kept the row, so a backup
// did run that week. Shown as the bare enum name it read as a failure, and it
// is the opposite: it is retention working.
const STATUS: Record<Backup["status"], { label: string; tone: PillTone; title: string }> = {
  RUNNING: { label: "Running", tone: "info", title: "Being written now" },
  READY: { label: "Ready", tone: "good", title: "Downloadable" },
  FAILED: { label: "Failed", tone: "bad", title: "Nothing was stored" },
  PURGED: {
    label: "Aged out",
    tone: "neutral",
    title: "Rolling retention deleted the archive. The row stays, so you can still see a backup ran.",
  },
};

function StatusBadge({ status }: { status: Backup["status"] }) {
  const s = STATUS[status];
  return (
    <Pill tone={s.tone} dot={status !== "RUNNING"} title={s.title}>
      {status === "RUNNING" ? <Loader2 className="h-3 w-3 animate-spin" aria-hidden="true" /> : null}
      {s.label}
    </Pill>
  );
}

export default function BackupsPage() {
  const { data: backups, isLoading } = useBackups();
  const { data: keep } = useBackupRetention();
  const generate = useGenerateBackup();
  const download = useDownloadBackup();

  const rows = backups ?? [];
  const running = rows.some((b) => b.status === "RUNNING");
  // Only the archives that still exist count towards the retention window.
  // The aged-out rows are history, not snapshots you have.
  const held = rows.filter((b) => b.status === "READY").length;

  return (
    <div>
      <PageHeader
        title="Backups"
        subtitle="Full-database snapshots, on a schedule you set or on demand."
        actions={
          <Button onClick={() => generate.mutate()} disabled={running} loading={generate.isPending}>
            {!generate.isPending && <RefreshCw className="h-3.5 w-3.5" />}
            {running ? "Backing up\u2026" : "Back up now"}
          </Button>
        }
      />

      {generate.isError ? (
        <div className="mb-6 flex items-start gap-3 rounded-xl border border-danger-line bg-danger-soft px-4 py-3">
          <AlertCircle className="mt-0.5 h-4 w-4 shrink-0 text-danger" aria-hidden="true" />
          <p className="text-sm text-text-secondary">
            {getApiErrorMessage(generate.error, "The backup did not start.")}
          </p>
        </div>
      ) : null}

      {/* Two cards, side by side: what will happen, and what has.
          They used to be stacked full width, which put the snapshot list below
          the fold on the one screen where "do I have a backup?" is the whole
          question. */}
      <div className="grid gap-6 lg:grid-cols-2">
        <ScheduleCard />

        <section className="flex flex-col overflow-hidden rounded-2xl border border-border bg-bg-elevated">
          <header className="flex flex-wrap items-center justify-between gap-3 border-b border-border px-5 py-4">
            <div className="flex items-center gap-3">
              <span
                className="inline-flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-bg-hover text-text-secondary"
                aria-hidden="true"
              >
                <Database className="h-4 w-4" />
              </span>
              <div>
                <h2 className="text-[15px] font-semibold text-foreground">Snapshots</h2>
                <p className="text-[13px] text-text-secondary">
                  {/* The retention window, from the server rather than typed
                      here, so raising it does not leave the page lying. */}
                  {keep ? held + " of " + keep + " kept" : "Downloadable archives"}
                </p>
              </div>
            </div>
          </header>

          {isLoading ? (
            <div className="flex flex-1 items-center justify-center py-16">
              <Loader2 className="h-5 w-5 animate-spin text-text-muted" aria-hidden="true" />
            </div>
          ) : rows.length === 0 ? (
            <EmptyState
              className="flex-1 justify-center"
              icon={<Database className="h-5 w-5" />}
              title="No snapshots yet"
              description="The first scheduled one lands on the next run above, and the button at the top of the page takes one immediately."
            />
          ) : (
            <div className="flex-1 overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border text-left text-xs font-medium uppercase tracking-wider text-text-muted">
                    <th className="px-5 py-3 font-medium">Taken</th>
                    <th className="px-4 py-3 font-medium">Status</th>
                    <th className="px-4 py-3 font-medium">Size</th>
                    <th className="px-5 py-3" />
                  </tr>
                </thead>
                <tbody>
                  {rows.map((b) => (
                    <tr key={b.id} className="border-b border-border/60 last:border-0">
                      <td className="px-5 py-3 align-top">
                        <span className="block text-foreground">
                          {new Date(b.created_at).toLocaleString()}
                        </span>
                        <span className="mt-1 block text-[13px] text-text-muted">
                          {b.kind.toLowerCase()}
                          {b.table_count
                            ? " \u00b7 " + b.table_count + " tables, " + b.row_count.toLocaleString() + " rows"
                            : ""}
                        </span>
                      </td>
                      <td className="px-4 py-3 align-top">
                        <StatusBadge status={b.status} />
                        {b.status === "FAILED" && b.error ? (
                          <p className="mt-1.5 max-w-xs text-[13px] text-danger" title={b.error}>
                            {b.error}
                          </p>
                        ) : null}
                      </td>
                      <td className="px-4 py-3 align-top text-text-secondary tabular-nums">
                        {formatBytes(b.size_bytes)}
                      </td>
                      <td className="px-5 py-3 align-top text-right">
                        {b.status === "READY" ? (
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => download.mutate(b.id)}
                            disabled={download.isPending}
                          >
                            <Download className="h-3.5 w-3.5" />
                            Download
                          </Button>
                        ) : null}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>
      </div>

      {/* What is in one, and how to put it back.
          The restore command was a footnote in muted 12px under a table, which
          is where you put something nobody needs. It is the opposite: it is
          the only part of this page that matters on the day it matters, and
          the day it matters is not a day for reading small print. */}
      <section className="mt-6 rounded-2xl border border-border bg-bg-elevated">
        <header className="flex items-center gap-3 border-b border-border px-5 py-4">
          <span
            className="inline-flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-bg-hover text-text-secondary"
            aria-hidden="true"
          >
            <Package className="h-4 w-4" />
          </span>
          <div>
            <h2 className="text-[15px] font-semibold text-foreground">What is in a backup</h2>
            <p className="text-[13px] text-text-secondary">
              And the one command that puts it back.
            </p>
          </div>
        </header>
        <div className="px-5 py-5">
          <ul className="mb-4 space-y-1.5 text-[13px] text-text-secondary">
            <li>
              <code className="font-mono text-foreground">*.csv</code>, one per table, readable
              by anything.
            </li>
            <li>
              <code className="font-mono text-foreground">dump.sql</code>, every row as an
              INSERT, for putting it back as it was.
            </li>
            <li>
              <code className="font-mono text-foreground">metadata.json</code>, what was taken
              and when, so an archive can identify itself.
            </li>
          </ul>

          <div className="flex items-center gap-2 rounded-lg border border-border bg-bg-tertiary px-3 py-2.5">
            <code className="min-w-0 flex-1 truncate font-mono text-[13px] text-foreground">
              grit restore backup.zip
            </code>
            <button
              type="button"
              onClick={async () => {
                try {
                  await navigator.clipboard.writeText("grit restore backup.zip");
                  toast.success("Command copied.");
                } catch {
                  toast.error("Could not copy. Select it and copy by hand.");
                }
              }}
              className="shrink-0 rounded p-1 text-text-muted transition-grit hover:bg-bg-hover hover:text-foreground"
              aria-label="Copy the restore command"
            >
              <Copy className="h-3.5 w-3.5" />
            </button>
          </div>

          <Note className="mt-3">
            Run it against a scratch database before you need it. An untested backup is a belief,
            not a backup, and the first time anybody finds out which is on the worst day of the
            year.
          </Note>
        </div>
      </section>
    </div>
  );
}
`
}
