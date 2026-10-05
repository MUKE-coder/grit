"use client";

import Link from "next/link";
import { useQueryClient } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { IconButton } from "@/components/ui/IconButton";
import { SkeletonTable } from "@/components/ui/Skeleton";
import { useToastedMutation } from "@/hooks/use-toasted-mutation";
import { notificationKeys, useNotificationList, type Notification } from "@/hooks/use-notifications";
import { Check, AlertCircle, AlertTriangle, Activity as ActivityIcon, Bell } from "@/lib/icons";
import { safeHref } from "@/lib/safe-href";
import { apiClient } from "@/lib/api-client";

const severityClass: Record<Notification["severity"], string> = {
  critical: "bg-danger/10 text-danger",
  high:     "bg-warning/10 text-warning",
  medium:   "bg-info/10 text-info",
  low:      "bg-bg-hover text-text-secondary",
  info:     "bg-bg-hover text-text-secondary",
};

const sourceIcon: Record<Notification["source"], React.ReactNode> = {
  sentinel: <AlertTriangle className="h-4 w-4" />,
  pulse:    <ActivityIcon className="h-4 w-4" />,
  system:   <AlertCircle className="h-4 w-4" />,
};

export default function NotificationsPage() {
  const queryClient = useQueryClient();
  // The list the bell in the header reads, so the page adds no request of its own.
  const { data, isLoading } = useNotificationList();

  const markRead = useToastedMutation({
    mutationFn: async (id: string) => apiClient.post("/api/notifications/" + id + "/read"),
    successMessage: "Marked read",
    silentSuccess: true,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: notificationKeys.all }),
  });

  const markAllRead = useToastedMutation({
    mutationFn: async () => apiClient.post("/api/notifications/read-all"),
    successMessage: "All notifications marked read",
    onSuccess: () => queryClient.invalidateQueries({ queryKey: notificationKeys.all }),
  });

  const items = data?.data || [];
  const unread = data?.unread || 0;

  return (
    <div>
      <PageHeader
        title="Notifications"
        subtitle="Recent system, security, and performance events"
        actions={
          unread > 0 ? (
            <IconButton
              variant="secondary"
              icon={<Check className="h-4 w-4" />}
              label="Mark all read"
              onClick={() => markAllRead.mutate()}
              disabled={markAllRead.isPending}
            />
          ) : null
        }
      />

      {isLoading ? (
        <SkeletonTable rows={6} columns={3} />
      ) : items.length === 0 ? (
        <div className="rounded-xl border border-border bg-bg-elevated p-12 text-center">
          <Bell className="mx-auto h-10 w-10 text-text-muted" />
          <p className="mt-3 text-base font-medium text-foreground">You're all caught up</p>
          <p className="mt-1 text-sm text-text-muted">Nothing needs your attention right now.</p>
        </div>
      ) : (
        <ul className="space-y-2">
          {items.map((n) => {
            const unreadRow = !n.read_at;
            return (
              <li
                key={n.id}
                className={
                  "rounded-xl border p-4 transition-colors " +
                  (unreadRow ? "border-accent/30 bg-accent/5" : "border-border bg-bg-elevated")
                }
              >
                <div className="flex items-start gap-3">
                  <span className={"mt-0.5 inline-flex h-7 w-7 shrink-0 items-center justify-center rounded-md " + severityClass[n.severity]}>
                    {sourceIcon[n.source]}
                  </span>
                  <div className="min-w-0 flex-1">
                    <Link
                      href={safeHref(n.link)}
                      onClick={() => { if (unreadRow) markRead.mutate(n.id); }}
                      className="block"
                    >
                      <div className="flex items-center gap-2">
                        <p className="truncate text-sm font-semibold text-foreground">{n.title}</p>
                        {n.count > 1 && (
                          <span className="rounded bg-bg-hover px-1.5 text-xs font-medium text-text-muted">×{n.count}</span>
                        )}
                        <span className={"rounded px-1.5 py-0.5 text-[10px] font-semibold uppercase " + severityClass[n.severity]}>
                          {n.severity}
                        </span>
                      </div>
                      {n.body && <p className="mt-0.5 text-sm text-text-secondary">{n.body}</p>}
                      <p className="mt-1 text-xs text-text-muted">
                        {n.source} · {new Date(n.created_at).toLocaleString()}
                      </p>
                    </Link>
                  </div>
                  {unreadRow && (
                    <button
                      type="button"
                      onClick={() => markRead.mutate(n.id)}
                      aria-label="Mark read"
                      className="shrink-0 rounded p-1.5 text-text-muted hover:bg-bg-hover hover:text-foreground"
                    >
                      <Check className="h-4 w-4" />
                    </button>
                  )}
                </div>
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
}
