"use client";

import Link from "next/link";
import dynamic from "next/dynamic";
import { Suspense, useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { useMe } from "@/hooks/use-auth";
import { useNotificationList, type NotificationList } from "@/hooks/use-notifications";
import { DashboardSecurityNudges } from "@/components/dashboard/security-nudges";
import { usePermissions } from "@/hooks/use-permissions";
import { resources } from "@/resources";
import { PageHeader } from "@/components/chrome/PageHeader";
import { LiveClock } from "@/components/chrome/LiveClock";
import { SkeletonCards } from "@/components/ui/Skeleton";
import { DateFilter, type DateRange } from "@/components/tables/date-filter";
import { ResourceWidgetsRow } from "@/components/dashboard/ResourceWidgetsRow";
import { apiClient } from "@/lib/api-client";
import {
  Activity as ActivityIcon, ArrowUpRight,
  Users, Bell, Database, Shield, getIcon,
} from "@/lib/icons";

// Recharts is about 400 KB. The chart loads after the page renders, so the
// greeting, stat tiles and activity feed do not wait for it. Suspense shows the
// placeholder in the Vite admin, where dynamic() is React.lazy.
const SeverityPieChart = dynamic(
  () => import("@/components/dashboard/DashboardCharts").then((m) => m.SeverityPieChart),
  { ssr: false, loading: ChartPlaceholder },
);

function ChartPlaceholder() {
  return <div className="h-full w-full animate-pulse rounded-lg bg-bg-hover/40" />;
}

interface MeStats {
  users: number;
  active_users: number;
}
interface ActivityRow {
  id: string;
  action: string;
  severity: "info" | "warn" | "critical";
  summary: string;
  ip_address: string;
  created_at: string;
}
interface ActivityListResponse { data: ActivityRow[] }
interface ActivityStatsResponse { data: { info: number; warn: number; critical: number; total: number } }

function unreadCount(list: NotificationList) {
  return list.unread || 0;
}

export default function DashboardPage() {
  const { data: user } = useMe();
  // Gate every resource-driven widget by the viewer's "<slug>.view" grant so
  // the dashboard shows the same surface the sidebar does — a user granted
  // only Categories + Products sees only those, not the whole catalog. Fails
  // open while permissions load (can() is closed until then); super-admins
  // short-circuit inside can(). Every underlying route is enforced server-side.
  const { can, isSuper, isLoading: permsLoading } = usePermissions();
  // hidden is an inline --items child: no sidebar entry, no page of its own,
  // and no case in the server's resource-stats dispatch, which the generator
  // writes only for top-level resources. Left in, each one costs two 400s on
  // every dashboard load and offers a tile to a table nobody browses.
  const visibleResources = resources.filter(
    (r) => !r.hidden && (permsLoading || can(r.slug + ".view")),
  );
  const canViewUsers = permsLoading || can("users.view");
  const canViewActivity = permsLoading || can("audit.view");
  // v3.31.44 -- shared DateFilter scopes the per-resource widgets
  // (Total + Latest N). Held in state here so the filter survives
  // hot-reload but doesn't bleed into the URL (would conflict with
  // the resource list pages' own URL-bound filter).
  const [dateRange, setDateRange] = useState<DateRange>({});
  const greeting = (() => {
    const h = new Date().getHours();
    if (h < 12) return "Good morning";
    if (h < 17) return "Good afternoon";
    return "Good evening";
  })();

  // Pull a small bundle of stats from endpoints that already exist —
  // each query is cheap, falls back to zero on error so partial outages
  // don't blank the dashboard.
  const userCount = useQuery<MeStats>({
    queryKey: ["dashboard", "users"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<{ meta?: { total: number } }>("/api/users?page_size=1");
        return { users: data.meta?.total || 0, active_users: data.meta?.total || 0 };
      } catch {
        return { users: 0, active_users: 0 };
      }
    },
    refetchInterval: 60_000,
  });

  const activityStats = useQuery<ActivityStatsResponse["data"]>({
    queryKey: ["dashboard", "activity-stats"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<ActivityStatsResponse>("/api/user-activity/stats");
        return data.data;
      } catch {
        return { info: 0, warn: 0, critical: 0, total: 0 };
      }
    },
    refetchInterval: 60_000,
  });

  // The bell in the header reads the same list, so the two share one request
  // and a notification marked read in the bell updates this tile as well.
  const notifications = useNotificationList(unreadCount);

  const recentActivity = useQuery<ActivityRow[]>({
    queryKey: ["dashboard", "recent-activity"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<ActivityListResponse>("/api/user-activity?page_size=8");
        return data.data;
      } catch {
        return [];
      }
    },
    refetchInterval: 60_000,
  });

  const severitySeries = useMemo(() => {
    const s = activityStats.data;
    if (!s || s.total === 0) {
      return [
        { name: "Info", value: 1, color: "var(--info)" },
      ];
    }
    return [
      { name: "Info",     value: s.info,     color: "var(--info)"    },
      { name: "Warn",     value: s.warn,     color: "var(--warning)" },
      { name: "Critical", value: s.critical, color: "var(--danger)"  },
    ].filter((d) => d.value > 0);
  }, [activityStats.data]);

  const statsLoading = userCount.isLoading || activityStats.isLoading;

  return (
    <div>
      <DashboardSecurityNudges />
      <PageHeader
        title={greeting + ", " + (user?.first_name || "Admin")}
        subtitle="Here's a snapshot of what's happening across your app right now."
        actions={<LiveClock />}
      />

      {/* v3.31.44 -- range filter scopes the per-resource widgets
          below. Sits on its own row so the existing "Resources"
          stat tile (count of registered modules) reads correctly
          against the dashboard's overall mode. */}
      <div className="mb-4 flex items-center justify-between">
        <p className="text-xs font-medium uppercase tracking-wide text-text-muted">
          Showing resource activity for
        </p>
        <DateFilter value={dateRange} onChange={setDateRange} label="Range" />
      </div>

      {/* Stat tiles */}
      {statsLoading ? (
        <SkeletonCards count={4} />
      ) : (
        <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
          {canViewUsers && (
            <StatTile
              label="Users"
              value={userCount.data?.users ?? 0}
              icon={<Users className="h-4 w-4" />}
              href="/resources/users"
              accent="info"
            />
          )}
          {canViewActivity && (
            <StatTile
              label="Events (24h)"
              value={activityStats.data?.total ?? 0}
              icon={<ActivityIcon className="h-4 w-4" />}
              href="/system/activity"
              accent="default"
              sublabel={(activityStats.data?.critical ?? 0) > 0 ? (activityStats.data?.critical + " critical") : "All clear"}
              sublabelTone={(activityStats.data?.critical ?? 0) > 0 ? "danger" : "success"}
            />
          )}
          <StatTile
            label="Notifications"
            value={notifications.data ?? 0}
            icon={<Bell className="h-4 w-4" />}
            href="/system/notifications"
            accent={notifications.data ? "warning" : "default"}
            sublabel={notifications.data ? "unread" : "you're caught up"}
          />
          <StatTile
            label="Resources"
            value={visibleResources.length}
            icon={<Database className="h-4 w-4" />}
            href="/dashboard"
            accent="default"
            sublabel="modules registered"
          />
        </div>
      )}

      {/* Activity row. There is no per-day event count endpoint, so the
          "past 7 days" chart that used to sit here plotted random numbers.
          It is gone rather than faked: the feed and the 24 hour severity mix
          below are both read from /api/user-activity. */}
      <div className="mt-6 grid grid-cols-1 gap-4 lg:grid-cols-3">
        <div className="rounded-xl border border-border bg-bg-elevated lg:col-span-2">
          <header className="flex items-center justify-between border-b border-border px-5 py-3.5">
            <div>
              <p className="text-sm font-semibold text-foreground">Recent activity</p>
              <p className="text-xs text-text-muted">Latest 8 events across the platform</p>
            </div>
            <Link href="/system/activity" className="inline-flex min-h-[24px] items-center text-xs font-medium text-accent hover:text-accent-hover">
              View all
            </Link>
          </header>
          {recentActivity.isLoading ? (
            <div className="px-5 py-12 text-center text-sm text-text-muted">Loading...</div>
          ) : (recentActivity.data ?? []).length === 0 ? (
            <div className="px-5 py-12 text-center text-sm text-text-muted">No activity yet.</div>
          ) : (
            <ul className="divide-y divide-border">
              {(recentActivity.data ?? []).map((row) => (
                <li key={row.id} className="flex items-start gap-3 px-5 py-3 text-sm">
                  <SeverityDot severity={row.severity} />
                  <div className="min-w-0 flex-1">
                    <p className="truncate text-foreground">{row.summary}</p>
                    <p className="text-xs text-text-muted">
                      <code className="font-mono">{row.action}</code>
                      {row.ip_address && (
                        <span title={row.ip_address}>
                          {" · "}
                          {row.ip_address === "::1" || row.ip_address === "127.0.0.1"
                            ? "localhost"
                            : row.ip_address}
                        </span>
                      )}
                    </p>
                  </div>
                  <span className="shrink-0 text-xs text-text-muted">
                    {timeAgo(row.created_at)}
                  </span>
                </li>
              ))}
            </ul>
          )}
        </div>

        <div className="rounded-xl border border-border bg-bg-elevated p-5">
          <div className="mb-4">
            <p className="text-sm font-semibold text-foreground">Severity mix</p>
            <p className="text-xs text-text-muted">Past 24 hours</p>
          </div>
          <div className="h-48 w-full">
            <Suspense fallback={<ChartPlaceholder />}>
              <SeverityPieChart data={severitySeries} />
            </Suspense>
          </div>
          <ul className="mt-2 space-y-1.5 text-xs">
            {severitySeries.map((s) => (
              <li key={s.name} className="flex items-center justify-between">
                <span className="flex items-center gap-2 text-text-secondary">
                  <span className="inline-block h-2 w-2 rounded-full" style={{ background: s.color }} />
                  {s.name}
                </span>
                <span className="font-semibold text-foreground">{s.value}</span>
              </li>
            ))}
          </ul>
        </div>
      </div>

      {/* Quick Access tiles — bottom section */}
      <div className="mt-6">
        <h2 className="mb-3 text-sm font-semibold uppercase tracking-wide text-text-muted">Quick access</h2>
        <div className="grid grid-cols-2 gap-3 md:grid-cols-3 lg:grid-cols-4">
          {visibleResources.slice(0, 8).map((r) => {
            const Icon = getIcon(r.icon);
            return (
              <Link
                key={r.slug}
                href={"/resources/" + r.slug}
                className="group rounded-xl border border-border bg-bg-elevated p-4 transition-colors hover:bg-bg-hover"
              >
                <div className="mb-2 inline-flex h-9 w-9 items-center justify-center rounded-lg bg-accent/10 text-accent">
                  <Icon className="h-4 w-4" />
                </div>
                <p className="text-sm font-semibold text-foreground group-hover:text-accent">
                  {r.label?.plural ?? r.name}
                </p>
                <p className="text-xs text-text-muted">Manage {(r.label?.plural ?? r.slug).toLowerCase()}</p>
              </Link>
            );
          })}
          {(permsLoading || isSuper) && (
            <Link
              href="/system"
              className="group rounded-xl border border-dashed border-border bg-bg-elevated p-4 transition-colors hover:bg-bg-hover"
            >
              <div className="mb-2 inline-flex h-9 w-9 items-center justify-center rounded-lg bg-bg-hover text-text-secondary">
                <Shield className="h-4 w-4" />
              </div>
              <p className="text-sm font-semibold text-foreground group-hover:text-accent">System hub</p>
              <p className="text-xs text-text-muted">Jobs, files, security, observability</p>
            </Link>
          )}
        </div>
      </div>

      {/* v3.31.44 -- per-resource widgets. One row per registered
          resource: Total + 30-day sparkline on the left, Latest N
          records on the right. Resources opting out (`dashboard: { enabled: false }`)
          are skipped here. Each row's queries are keyed on the
          shared DateFilter range so changing the filter refetches
          every widget in lockstep. */}
      {visibleResources.filter((r) => r.dashboard?.enabled !== false).length > 0 && (
        <div className="mt-8">
          <div className="mb-3 flex items-center justify-between">
            <h2 className="text-sm font-semibold uppercase tracking-wide text-text-muted">
              By resource
            </h2>
            <p className="text-[11px] text-text-muted">
              Scoped to the range above &middot; sparkline always last 30 days
            </p>
          </div>
          <div className="space-y-4">
            {visibleResources
              .filter((r) => r.dashboard?.enabled !== false)
              .map((r) => (
                <ResourceWidgetsRow key={r.slug} resource={r} dateRange={dateRange} />
              ))}
          </div>
        </div>
      )}
    </div>
  );
}

interface StatTileProps {
  label: string;
  value: number | string;
  icon: React.ReactNode;
  href: string;
  accent: "default" | "info" | "warning" | "danger";
  sublabel?: string;
  sublabelTone?: "success" | "danger" | "muted";
}

const accentClass: Record<StatTileProps["accent"], string> = {
  default: "bg-accent/10 text-accent",
  info: "bg-info/10 text-info",
  warning: "bg-warning/10 text-warning",
  danger: "bg-danger/10 text-danger",
};

const sublabelClass: Record<NonNullable<StatTileProps["sublabelTone"]>, string> = {
  success: "text-success",
  danger: "text-danger",
  muted: "text-text-muted",
};

function StatTile({ label, value, icon, href, accent, sublabel, sublabelTone = "muted" }: StatTileProps) {
  return (
    <Link
      href={href}
      className="group rounded-xl border border-border bg-bg-elevated p-4 transition-colors hover:bg-bg-hover"
    >
      <div className="flex items-center justify-between">
        <span className={"inline-flex h-9 w-9 items-center justify-center rounded-lg " + accentClass[accent]}>
          {icon}
        </span>
        <ArrowUpRight className="h-4 w-4 text-text-muted opacity-0 transition-opacity group-hover:opacity-100" />
      </div>
      <p className="mt-3 text-xs font-medium uppercase tracking-wide text-text-muted">{label}</p>
      <p className="text-2xl font-bold text-foreground">{value}</p>
      {sublabel && (
        <p className={"mt-1 text-xs " + sublabelClass[sublabelTone]}>{sublabel}</p>
      )}
    </Link>
  );
}

const severityDotClass: Record<ActivityRow["severity"], string> = {
  info: "bg-info",
  warn: "bg-warning",
  critical: "bg-danger",
};

function SeverityDot({ severity }: { severity: ActivityRow["severity"] }) {
  return (
    <span className="mt-1.5 inline-block h-2 w-2 shrink-0 rounded-full" >
      <span className={"block h-full w-full rounded-full " + severityDotClass[severity]} />
    </span>
  );
}

function timeAgo(iso: string): string {
  const diff = Date.now() - new Date(iso).getTime();
  const sec = Math.round(diff / 1000);
  if (sec < 60) return sec + "s ago";
  const min = Math.round(sec / 60);
  if (min < 60) return min + "m ago";
  const hr = Math.round(min / 60);
  if (hr < 24) return hr + "h ago";
  const days = Math.round(hr / 24);
  return days + "d ago";
}
