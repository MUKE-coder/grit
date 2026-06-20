package scaffold

// v3.31 UI bundle: toast hook, skeleton primitive, table user-cell helper,
// system hub, and the notifications page that v3.30 forgot.

// adminToastHook — useToastedMutation wraps useMutation from react-query
// and emits a sonner toast on success / error with millisecond timing.
// Every mutation in the dashboard can opt in by replacing useMutation
// with useToastedMutation, no other code change needed.
func adminToastHook() string {
	return `"use client";

import { useMutation, type UseMutationOptions } from "@tanstack/react-query";
import { toast } from "sonner";

interface ToastedOptions<TData, TError, TVariables, TContext>
  extends UseMutationOptions<TData, TError, TVariables, TContext> {
  /** Success message. When a function, receives the mutation result. */
  successMessage?: string | ((data: TData) => string);
  /** Error message. When a function, receives the thrown error. */
  errorMessage?: string | ((err: TError) => string);
  /** Skip the success toast (e.g. when navigating away on success). */
  silentSuccess?: boolean;
}

/**
 * Drop-in replacement for useMutation. Times every mutation and emits a
 * toast on settle:
 *
 *   ✓ Ticket opened — 142ms
 *   ✗ Couldn't reach the API — 1203ms
 *
 * The timing is intentionally surfaced. Users get tactile feedback that
 * the request ran (vs hung) and developers spot regressions early.
 */
export function useToastedMutation<TData = unknown, TError = unknown, TVariables = void, TContext = unknown>(
  options: ToastedOptions<TData, TError, TVariables, TContext>
) {
  const { successMessage, errorMessage, silentSuccess, onSuccess, onError, ...rest } = options;

  return useMutation<TData, TError, TVariables, TContext>({
    ...rest,
    onSuccess: (data, vars, ctx) => {
      const elapsed = readElapsed(ctx as MutationContext | undefined);
      if (!silentSuccess) {
        const msg = typeof successMessage === "function" ? successMessage(data) : (successMessage || "Done");
        toast.success(msg + (elapsed != null ? " — " + elapsed + "ms" : ""));
      }
      onSuccess?.(data, vars, ctx);
    },
    onError: (err, vars, ctx) => {
      const elapsed = readElapsed(ctx as MutationContext | undefined);
      const fallback = pickErrorMessage(err);
      const msg = typeof errorMessage === "function" ? errorMessage(err) : (errorMessage || fallback);
      toast.error(msg + (elapsed != null ? " — " + elapsed + "ms" : ""));
      onError?.(err, vars, ctx);
    },
    onMutate: async (vars) => {
      // Stamp the start time on the mutation context so onSuccess/onError
      // can subtract. We piggy-back on the caller's context (if they
      // supplied an onMutate) by spreading the awaited result.
      const userCtx = rest.onMutate ? await rest.onMutate(vars) : undefined;
      return { ...(userCtx as object), __startedAt: performance.now() } as TContext;
    },
  });
}

type MutationContext = { __startedAt?: number } | undefined;

function readElapsed(ctx: MutationContext): number | null {
  const t = ctx?.__startedAt;
  if (typeof t !== "number") return null;
  return Math.round(performance.now() - t);
}

function pickErrorMessage(err: unknown): string {
  const m = (err as { response?: { data?: { error?: { message?: string } } } })
    ?.response?.data?.error?.message;
  if (m) return m;
  if (err instanceof Error) return err.message;
  return "Something went wrong";
}

// Re-export sonner's toast so pages can drop ad-hoc toasts (e.g. on a
// copy-to-clipboard) without a second import line.
export { toast };
`
}

// adminSkeletonComponent — composable loading-state primitive. Pages
// render a Skeleton tree while data is pending instead of a generic
// spinner; perceived performance improves and layout doesn't jump when
// real content arrives.
func adminSkeletonComponent() string {
	return `"use client";

import type { HTMLAttributes } from "react";

interface SkeletonProps extends HTMLAttributes<HTMLDivElement> {
  /** Convenience preset — leaves you free to override via className. */
  shape?: "rect" | "text" | "circle";
}

/**
 * Animated placeholder block. Compose larger skeletons by stacking these
 * with sizing classes:
 *
 *   <Skeleton className="h-8 w-64" />
 *   <Skeleton shape="circle" className="h-10 w-10" />
 *   <Skeleton shape="text" className="w-3/4" />
 *
 * Uses the page's --bg-hover token + a subtle pulse animation so it
 * adapts to every theme + light/dark mode out of the box.
 */
export function Skeleton({ shape = "rect", className = "", ...rest }: SkeletonProps) {
  const shapeClass =
    shape === "circle" ? "rounded-full" :
    shape === "text"   ? "rounded h-3.5" :
                         "rounded-lg";
  return (
    <div
      {...rest}
      className={"animate-pulse bg-bg-hover " + shapeClass + " " + className}
    />
  );
}

/**
 * SkeletonTable renders a placeholder for the ResponsiveTable primitive.
 * Pass the column count to roughly match the live table's geometry so
 * the layout doesn't jump when data arrives.
 */
export function SkeletonTable({ rows = 5, columns = 4 }: { rows?: number; columns?: number }) {
  return (
    <>
      <div className="hidden md:block overflow-hidden rounded-xl border border-border bg-bg-elevated">
        <div className="border-b border-border px-4 py-3 flex gap-4">
          {Array.from({ length: columns }).map((_, i) => (
            <Skeleton key={i} shape="text" className="flex-1 max-w-[120px]" />
          ))}
        </div>
        <div className="divide-y divide-border">
          {Array.from({ length: rows }).map((_, i) => (
            <div key={i} className="flex gap-4 px-4 py-3.5">
              {Array.from({ length: columns }).map((_, j) => (
                <Skeleton key={j} shape="text" className="flex-1" />
              ))}
            </div>
          ))}
        </div>
      </div>
      <ul className="md:hidden space-y-3">
        {Array.from({ length: rows }).map((_, i) => (
          <li key={i} className="rounded-xl border border-border bg-bg-elevated p-4 space-y-2">
            <Skeleton shape="text" className="w-1/2" />
            <Skeleton shape="text" className="w-3/4" />
            <Skeleton shape="text" className="w-1/3" />
          </li>
        ))}
      </ul>
    </>
  );
}

/**
 * SkeletonCards — placeholder for stats-card rows on dashboards.
 */
export function SkeletonCards({ count = 4 }: { count?: number }) {
  return (
    <div className={"grid grid-cols-2 gap-3 md:grid-cols-" + count}>
      {Array.from({ length: count }).map((_, i) => (
        <div key={i} className="rounded-xl border border-border bg-bg-elevated p-4 space-y-2">
          <Skeleton shape="text" className="w-20" />
          <Skeleton className="h-7 w-16" />
        </div>
      ))}
    </div>
  );
}
`
}

// adminUserCellComponent — packed presentation cell for "user" columns
// in tables. Renders avatar (or initials), name, and email in a single
// flex row so wide tables don't need a separate Email column.
func adminUserCellComponent() string {
	return `"use client";

interface UserCellProps {
  user?: {
    first_name?: string;
    last_name?: string;
    email?: string;
    avatar?: string;
  } | null;
  /** Override the displayed primary line. */
  name?: string;
  /** Fallback to use when user is null (e.g. system events). */
  fallback?: string;
  /** Set true to render initials only (no name + email). */
  compact?: boolean;
}

/**
 * Single-cell user display: avatar + name (top) + email (bottom). Falls
 * back to "—" or a custom fallback when no user is present. Pack this
 * into a table's user column to keep tables narrow on small screens:
 *
 *   cell: (row) => <UserCell user={row.user} />
 */
export function UserCell({ user, name, fallback, compact }: UserCellProps) {
  if (!user) {
    return <span className="text-sm text-text-muted">{fallback || "—"}</span>;
  }
  const fullName = name || [user.first_name, user.last_name].filter(Boolean).join(" ") || "User";
  const initials = ((user.first_name?.[0] || "") + (user.last_name?.[0] || "")).toUpperCase() || "U";

  return (
    <div className="flex items-center gap-2.5 min-w-0">
      <span className="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-full ring-1 ring-border bg-bg-elevated text-xs font-semibold text-foreground overflow-hidden">
        {user.avatar ? (
          <img src={user.avatar} alt={fullName} className="h-full w-full object-cover" />
        ) : initials}
      </span>
      {!compact && (
        <div className="min-w-0">
          <p className="truncate text-sm font-medium text-foreground">{fullName}</p>
          {user.email && (
            <p className="truncate text-xs text-text-muted">{user.email}</p>
          )}
        </div>
      )}
    </div>
  );
}
`
}

// adminSystemHubPage — single /system page listing every system surface
// (jobs, files, cron, mail, security, observability, activity, support).
// Replaces the six individual sidebar entries with one hub entry that
// users land on first.
func adminSystemHubPage() string {
	return `"use client";

import Link from "next/link";
import { PageHeader } from "@/components/chrome/PageHeader";
import {
  Activity, Bell, Calendar, Database, FileText, Mail,
  MessageSquare, Settings, Shield, TrendingUp, Upload,
} from "@/lib/icons";

interface SystemTile {
  href: string;
  title: string;
  description: string;
  icon: React.ReactNode;
  tone: "default" | "danger" | "warning" | "info";
}

const TILES: SystemTile[] = [
  { href: "/system/jobs",         title: "Background Jobs",  description: "Queue depth, in-flight workers, dead-letter queue.",       icon: <Database className="h-5 w-5" />,     tone: "default" },
  { href: "/system/files",        title: "File Storage",     description: "Browse uploads, manage retention, audit usage.",            icon: <Upload className="h-5 w-5" />,       tone: "default" },
  { href: "/system/cron",         title: "Cron Schedules",   description: "Recurring jobs, next-run times, run history.",              icon: <Calendar className="h-5 w-5" />,     tone: "default" },
  { href: "/system/mail",         title: "Mail Preview",     description: "Email template gallery + recent send log.",                 icon: <Mail className="h-5 w-5" />,         tone: "default" },
  { href: "/system/security",     title: "Security",         description: "Sentinel summary — score, threats, AuthShield, CSP.",       icon: <Shield className="h-5 w-5" />,       tone: "danger"  },
  { href: "/system/observability",title: "Observability",    description: "Pulse summary — latency, SLOs, top N+1, runtime.",          icon: <TrendingUp className="h-5 w-5" />,   tone: "info"    },
  { href: "/system/activity",     title: "User Activity",    description: "Auth events, writes, operator actions with IP + severity.", icon: <Activity className="h-5 w-5" />,     tone: "default" },
  { href: "/system/support",      title: "Support",          description: "Incoming tickets, threads, assignments, closures.",         icon: <MessageSquare className="h-5 w-5" />, tone: "warning" },
  { href: "/system/notifications",title: "Notifications",    description: "Recent system + Sentinel + Pulse notifications.",           icon: <Bell className="h-5 w-5" />,         tone: "default" },
];

const TONE_RING: Record<SystemTile["tone"], string> = {
  default: "ring-border bg-bg-elevated",
  danger:  "ring-danger/30 bg-danger/5",
  warning: "ring-warning/30 bg-warning/5",
  info:    "ring-info/30 bg-info/5",
};

export default function SystemHubPage() {
  return (
    <div>
      <PageHeader
        title="System"
        subtitle="Every operational surface for this app. Pick a tile to dive in."
      />

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-3">
        {TILES.map((t) => (
          <Link
            key={t.href}
            href={t.href}
            className={
              "group rounded-xl ring-1 p-5 transition-colors hover:bg-bg-hover " + TONE_RING[t.tone]
            }
          >
            <div className="mb-3 inline-flex h-10 w-10 items-center justify-center rounded-lg bg-bg-secondary text-foreground group-hover:text-accent">
              {t.icon}
            </div>
            <p className="text-base font-semibold text-foreground">{t.title}</p>
            <p className="mt-1 text-sm text-text-secondary">{t.description}</p>
          </Link>
        ))}
      </div>
    </div>
  );
}
`
}

// adminNotificationsPage — full notifications list. The bell dropdown
// only shows the most recent 50; this page is the authoritative view
// with filters, mark-all-read, and severity buckets.
func adminNotificationsPage() string {
	return `"use client";

import Link from "next/link";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { IconButton } from "@/components/ui/IconButton";
import { SkeletonTable } from "@/components/ui/Skeleton";
import { useToastedMutation } from "@/hooks/use-toasted-mutation";
import { Check, AlertCircle, AlertTriangle, Activity as ActivityIcon, Bell } from "@/lib/icons";
import { apiClient } from "@/lib/api-client";

interface Notification {
  id: string;
  source: "sentinel" | "pulse" | "system";
  severity: "critical" | "high" | "medium" | "low" | "info";
  title: string;
  body: string;
  link: string;
  count: number;
  read_at: string | null;
  created_at: string;
}

interface ListResponse {
  data: Notification[];
  unread: number;
}

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
  const { data, isLoading } = useQuery<ListResponse>({
    queryKey: ["notifications", "list"],
    queryFn: async () => {
      const { data } = await apiClient.get<ListResponse>("/api/notifications");
      return data;
    },
  });

  const markRead = useToastedMutation({
    mutationFn: async (id: string) => apiClient.post("/api/notifications/" + id + "/read"),
    successMessage: "Marked read",
    silentSuccess: true,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["notifications"] }),
  });

  const markAllRead = useToastedMutation({
    mutationFn: async () => apiClient.post("/api/notifications/read-all"),
    successMessage: "All notifications marked read",
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["notifications"] }),
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
                      href={n.link || "#"}
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
`
}
