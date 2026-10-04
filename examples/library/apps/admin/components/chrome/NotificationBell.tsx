"use client";

import Link from "next/link";
import { useEffect, useRef, useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Bell, Check, AlertCircle, AlertTriangle, Activity } from "@/lib/icons";
import { safeHref } from "@/lib/safe-href";
import { apiClient } from "@/lib/api-client";
import { notificationKeys, useNotificationList, type Notification } from "@/hooks/use-notifications";

const severityColor: Record<Notification["severity"], string> = {
  critical: "text-danger",
  high: "text-warning",
  medium: "text-info",
  low: "text-text-secondary",
  info: "text-text-secondary",
};

const sourceIcon: Record<Notification["source"], typeof AlertCircle> = {
  sentinel: AlertTriangle,
  pulse: Activity,
  system: AlertCircle,
};

/**
 * Bell + dropdown. The bell badge polls /api/notifications every 60s for
 * fresh data. Opening the dropdown reveals the most recent items with a
 * mark-read button per row + a mark-all-read footer. Clicking a row
 * navigates to the linked surface (Sentinel finding, Pulse trace, etc.).
 */
export function NotificationBell() {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);
  const queryClient = useQueryClient();

  const { data } = useNotificationList();

  const markRead = useMutation({
    mutationFn: async (id: string) => apiClient.post("/api/notifications/" + id + "/read"),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: notificationKeys.all }),
  });

  const markAllRead = useMutation({
    mutationFn: async () => apiClient.post("/api/notifications/read-all"),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: notificationKeys.all }),
  });

  useEffect(() => {
    if (!open) return;
    const onClick = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", onClick);
    return () => document.removeEventListener("mousedown", onClick);
  }, [open]);

  const unread = data?.unread || 0;
  const items = data?.data || [];

  return (
    <div ref={ref} className="relative">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        aria-label={"Notifications" + (unread > 0 ? " (" + unread + " unread)" : "")}
        className="relative inline-flex h-9 w-9 items-center justify-center rounded-lg border border-border bg-bg-elevated text-text-secondary hover:bg-bg-hover transition-colors"
      >
        <Bell className="h-4 w-4" />
        {unread > 0 && (
          <span className="absolute -top-1 -right-1 flex h-4 min-w-4 items-center justify-center rounded-full bg-danger px-1 text-[10px] font-semibold text-white">
            {unread > 99 ? "99+" : unread}
          </span>
        )}
      </button>

      {open && (
        <div className="absolute right-0 z-50 mt-2 w-96 max-w-[calc(100vw-2rem)] overflow-hidden rounded-xl border border-border bg-bg-elevated shadow-xl">
          <header className="flex items-center justify-between border-b border-border px-4 py-3">
            <h3 className="text-sm font-semibold text-foreground">Notifications</h3>
            {unread > 0 && (
              <button
                type="button"
                onClick={() => markAllRead.mutate()}
                disabled={markAllRead.isPending}
                className="text-xs font-medium text-accent hover:text-accent-hover disabled:opacity-50"
              >
                Mark all read
              </button>
            )}
          </header>

          <div className="max-h-96 overflow-y-auto">
            {items.length === 0 ? (
              <p className="px-4 py-8 text-center text-sm text-text-muted">
                You're all caught up
              </p>
            ) : (
              <ul className="divide-y divide-border">
                {items.map((n) => {
                  const Icon = sourceIcon[n.source] || AlertCircle;
                  const unreadRow = !n.read_at;
                  return (
                    <li key={n.id} className={unreadRow ? "bg-bg-secondary/50" : ""}>
                      <div className="flex gap-3 px-4 py-3">
                        <Icon className={"mt-0.5 h-4 w-4 shrink-0 " + severityColor[n.severity]} />
                        <div className="min-w-0 flex-1">
                          <Link
                            href={safeHref(n.link)}
                            onClick={() => { setOpen(false); if (unreadRow) markRead.mutate(n.id); }}
                            className="block"
                          >
                            <p className="text-sm font-medium text-foreground truncate">
                              {n.title}
                              {n.count > 1 && (
                                <span className="ml-1 text-xs text-text-muted">×{n.count}</span>
                              )}
                            </p>
                            {n.body && (
                              <p className="mt-0.5 text-xs text-text-secondary line-clamp-2">{n.body}</p>
                            )}
                            <p className="mt-1 text-[10px] uppercase tracking-wide text-text-muted">
                              {n.source} · {new Date(n.created_at).toLocaleString()}
                            </p>
                          </Link>
                        </div>
                        {unreadRow && (
                          <button
                            type="button"
                            onClick={() => markRead.mutate(n.id)}
                            aria-label="Mark read"
                            className="shrink-0 rounded p-1 text-text-muted hover:bg-bg-hover hover:text-foreground"
                          >
                            <Check className="h-3.5 w-3.5" />
                          </button>
                        )}
                      </div>
                    </li>
                  );
                })}
              </ul>
            )}
          </div>

          <footer className="border-t border-border px-4 py-2 text-center">
            <Link
              href="/system/notifications"
              onClick={() => setOpen(false)}
              className="inline-flex min-h-[24px] items-center text-xs font-medium text-accent hover:text-accent-hover"
            >
              View all notifications
            </Link>
          </footer>
        </div>
      )}
    </div>
  );
}
