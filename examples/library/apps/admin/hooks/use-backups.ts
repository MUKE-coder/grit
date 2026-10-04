import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
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
      const { data } = await apiClient.get(`/api/backups/${id}/download`);
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
