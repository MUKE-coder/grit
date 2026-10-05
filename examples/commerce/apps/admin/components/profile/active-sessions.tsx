"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { Monitor, Smartphone, Loader2, LogOut, ShieldCheck } from "@/lib/icons";

type Session = {
  id: string;
  user_agent: string;
  ip: string;
  created_at: string;
  last_seen_at: string;
  expires_at: string;
  current: boolean;
};

// Just enough user-agent parsing to label a row. Deliberately not a UA-parsing
// dependency — the goal is "is this the laptop or the phone I remember?", not
// analytics-grade device detection.
export function describeDevice(ua: string): { label: string; mobile: boolean } {
  if (!ua) return { label: "Unknown device", mobile: false };
  const mobile = /Mobile|Android|iPhone|iPad|iPod/i.test(ua);

  let browser = "Browser";
  if (/Edg\//.test(ua)) browser = "Edge";
  else if (/OPR\/|Opera/.test(ua)) browser = "Opera";
  else if (/Chrome\//.test(ua) && !/Chromium/.test(ua)) browser = "Chrome";
  else if (/Firefox\//.test(ua)) browser = "Firefox";
  else if (/Safari\//.test(ua)) browser = "Safari";

  // Order matters. An iPhone UA contains "like Mac OS X" and an Android UA
  // contains "Linux", so the specific platforms have to be tested first or
  // every phone reports itself as a desktop.
  let os = "";
  if (/Windows NT/.test(ua)) os = "Windows";
  else if (/iPhone|iPad|iPod/.test(ua)) os = "iOS";
  else if (/Android/.test(ua)) os = "Android";
  else if (/Mac OS X/.test(ua)) os = "macOS";
  else if (/Linux/.test(ua)) os = "Linux";

  return { label: os ? browser + " on " + os : browser, mobile };
}

function relativeTime(iso: string): string {
  const then = new Date(iso).getTime();
  if (Number.isNaN(then)) return "";
  const seconds = Math.floor((Date.now() - then) / 1000);
  if (seconds < 60) return "just now";
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return minutes + (minutes === 1 ? " minute ago" : " minutes ago");
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return hours + (hours === 1 ? " hour ago" : " hours ago");
  const days = Math.floor(hours / 24);
  if (days < 30) return days + (days === 1 ? " day ago" : " days ago");
  return new Date(iso).toLocaleDateString();
}

export function ActiveSessions() {
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ["auth", "sessions"],
    queryFn: async () => {
      const { data: res } = await apiClient.get("/api/auth/sessions");
      return (res.data ?? []) as Session[];
    },
  });

  const revoke = useMutation({
    mutationFn: async (id: string) => {
      await apiClient.delete("/api/auth/sessions/" + id);
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["auth", "sessions"] }),
  });

  const revokeAll = useMutation({
    mutationFn: async () => {
      await apiClient.post("/api/auth/sessions/revoke-all");
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["auth", "sessions"] }),
  });

  const sessions = data ?? [];
  const others = sessions.filter((s) => !s.current).length;

  return (
    <div className="space-y-3">
      {isLoading ? (
        <div className="flex items-center gap-2 py-6 text-sm text-text-secondary">
          <Loader2 className="h-4 w-4 animate-spin" />
          Loading sessions…
        </div>
      ) : sessions.length === 0 ? (
        <p className="py-6 text-sm text-text-secondary">No active sessions.</p>
      ) : (
        <ul className="space-y-2">
          {sessions.map((s) => {
            const device = describeDevice(s.user_agent);
            const Icon = device.mobile ? Smartphone : Monitor;
            return (
              <li
                key={s.id}
                className="flex flex-col gap-3 rounded-xl border border-border bg-bg-secondary/40 p-4 sm:flex-row sm:items-center sm:justify-between"
              >
                <div className="flex items-start gap-3">
                  <span className="mt-0.5 rounded-lg border border-border bg-bg-tertiary p-2 text-text-secondary">
                    <Icon className="h-4 w-4" />
                  </span>
                  <div>
                    <p className="flex items-center gap-2 text-sm font-medium text-foreground">
                      {device.label}
                      {s.current && (
                        <span className="inline-flex items-center gap-1 rounded-full bg-success/10 px-2 py-0.5 text-[11px] font-semibold text-success">
                          <ShieldCheck className="h-3 w-3" />
                          This device
                        </span>
                      )}
                    </p>
                    <p className="mt-0.5 text-xs text-text-secondary">
                      {s.ip || "unknown IP"} · last active {relativeTime(s.last_seen_at)}
                    </p>
                  </div>
                </div>

                {!s.current && (
                  <button
                    type="button"
                    onClick={() => revoke.mutate(s.id)}
                    disabled={revoke.isPending}
                    className="inline-flex shrink-0 items-center gap-2 rounded-lg border border-border px-3 py-1.5 text-xs font-semibold text-text-secondary hover:border-danger/40 hover:bg-danger/10 hover:text-danger disabled:opacity-50"
                  >
                    {revoke.isPending && revoke.variables === s.id ? (
                      <Loader2 className="h-3 w-3 animate-spin" />
                    ) : (
                      <LogOut className="h-3 w-3" />
                    )}
                    Sign out
                  </button>
                )}
              </li>
            );
          })}
        </ul>
      )}

      {others > 0 && (
        <div className="flex justify-end pt-1">
          <button
            type="button"
            onClick={() => revokeAll.mutate()}
            disabled={revokeAll.isPending}
            className="inline-flex items-center gap-2 rounded-lg border border-danger/40 bg-danger/10 px-4 py-2 text-sm font-semibold text-danger hover:bg-danger/20 disabled:opacity-50"
          >
            {revokeAll.isPending ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <LogOut className="h-4 w-4" />
            )}
            Sign out of all other devices
          </button>
        </div>
      )}
    </div>
  );
}
