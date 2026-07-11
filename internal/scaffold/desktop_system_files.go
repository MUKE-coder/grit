package scaffold

// Desktop system/nav pages — full parity with the admin panel's Internal +
// System sections: Users, Blogs, Activity, Support, Notifications, Dashboard
// settings, System Health, Performance, Security, and the System hub. They use
// the desktop apiClient (React Query) and degrade gracefully when offline
// (every query falls back to an empty/zero shape, matching the admin's own
// try/catch behaviour). Purely-static pages (System hub) render with no
// network at all.

// desktopClientSystemUI holds the small shared building blocks the system
// pages reuse: a stat card, a titled section shell, and an empty state.
func desktopClientSystemUI() string {
	return `import type { LucideIcon } from "lucide-react";
import { cn } from "@/lib/utils";

type Tone = "default" | "success" | "warning" | "danger" | "info";

const toneChip: Record<Tone, string> = {
  default: "bg-accent/10 text-accent",
  success: "bg-success/10 text-success",
  warning: "bg-warning/10 text-warning",
  danger: "bg-danger/10 text-danger",
  info: "bg-info/10 text-info",
};

export function SystemStat({
  label, value, icon: Icon, tone = "default", sub,
}: { label: string; value: string | number; icon: LucideIcon; tone?: Tone; sub?: string }) {
  return (
    <div className="rounded-xl border border-border bg-surface p-5">
      <div className="flex items-start justify-between">
        <div className="min-w-0">
          <p className="text-[11px] font-semibold uppercase tracking-wider text-foreground-muted">{label}</p>
          <p className="mt-1 text-2xl font-bold tabular-nums text-foreground">{value}</p>
          {sub && <p className="mt-0.5 text-[12px] text-foreground-muted">{sub}</p>}
        </div>
        <span className={cn("inline-flex h-9 w-9 items-center justify-center rounded-lg", toneChip[tone])}>
          <Icon className="h-4 w-4" />
        </span>
      </div>
    </div>
  );
}

export function SectionCard({
  title, description, action, children,
}: { title: string; description?: string; action?: React.ReactNode; children: React.ReactNode }) {
  return (
    <div className="rounded-xl border border-border bg-surface">
      <header className="flex items-center justify-between border-b border-border px-5 py-3.5">
        <div>
          <p className="text-sm font-semibold text-foreground">{title}</p>
          {description && <p className="text-xs text-foreground-muted">{description}</p>}
        </div>
        {action}
      </header>
      {children}
    </div>
  );
}

export function EmptyState({ icon: Icon, title, hint }: { icon: LucideIcon; title: string; hint?: string }) {
  return (
    <div className="px-5 py-16 text-center">
      <div className="mx-auto mb-3 inline-flex rounded-full bg-surface-2 p-4">
        <Icon className="h-6 w-6 text-foreground-muted" />
      </div>
      <p className="text-[13px] text-foreground">{title}</p>
      {hint && <p className="mt-1 text-[12px] text-foreground-muted">{hint}</p>}
    </div>
  );
}

export function relTime(iso: unknown): string {
  if (!iso) return "";
  const t = new Date(String(iso)).getTime();
  if (Number.isNaN(t)) return "";
  const diff = Date.now() - t;
  const sec = Math.round(diff / 1000);
  if (sec < 60) return sec + "s ago";
  const min = Math.round(sec / 60);
  if (min < 60) return min + "m ago";
  const hr = Math.round(min / 60);
  if (hr < 24) return hr + "h ago";
  return Math.round(hr / 24) + "d ago";
}
`
}

// desktopClientSystemHubPage is the /app/system landing grid — static tiles,
// no network, renders fully offline.
func desktopClientSystemHubPage() string {
	return `import { createFileRoute, Link } from "@tanstack/react-router";
import {
  Activity, TrendingUp, Shield, Bell, MessageSquare, Settings,
  Users, FileText, type LucideIcon,
} from "lucide-react";
import { PageHeader } from "@/components/layout/page-header";

export const Route = createFileRoute("/app/system/")({
  component: SystemHubPage,
});

const TILES: { to: string; title: string; description: string; icon: LucideIcon }[] = [
  { to: "/app/system/health", title: "System Health", description: "Database, cache, jobs & email status", icon: Activity },
  { to: "/app/system/performance", title: "Performance", description: "Latency, traffic, errors & saturation", icon: TrendingUp },
  { to: "/app/system/security", title: "Security", description: "Bans, rate limits & recent threats", icon: Shield },
  { to: "/app/system/activity", title: "User Activity", description: "Audit log across the platform", icon: Activity },
  { to: "/app/system/support", title: "Support", description: "Tickets & conversations", icon: MessageSquare },
  { to: "/app/system/notifications", title: "Notifications", description: "System & security alerts", icon: Bell },
  { to: "/app/system/users", title: "Users", description: "Manage accounts & roles", icon: Users },
  { to: "/app/system/blogs", title: "Blogs", description: "Write & publish posts", icon: FileText },
  { to: "/app/system/dashboard-settings", title: "Dashboard settings", description: "Customize your dashboard widgets", icon: Settings },
];

function SystemHubPage() {
  return (
    <div>
      <PageHeader title="System" description="Operate, observe and administer your app." />
      <div className="mt-6 grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {TILES.map((t) => (
          <Link
            key={t.to}
            to={t.to}
            className="group rounded-xl border border-border bg-surface p-5 transition-colors hover:bg-surface-hover"
          >
            <div className="mb-3 inline-flex h-10 w-10 items-center justify-center rounded-lg bg-accent/10 text-accent">
              <t.icon className="h-5 w-5" />
            </div>
            <p className="text-sm font-semibold text-foreground group-hover:text-accent">{t.title}</p>
            <p className="mt-0.5 text-[12px] text-foreground-muted">{t.description}</p>
          </Link>
        ))}
      </div>
    </div>
  );
}
`
}

// desktopClientSystemUsersPage lists user accounts via the shared DataTable.
func desktopClientSystemUsersPage() string {
	return `import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { PageHeader } from "@/components/layout/page-header";
import { DataTable, type DataColumn } from "@/components/tables/data-table";
import { apiClient } from "@/lib/api-client";

export const Route = createFileRoute("/app/system/users")({
  component: SystemUsersPage,
});

type UserRow = Record<string, unknown> & { id: string };

const COLUMNS: DataColumn[] = [
  { key: "name", label: "Name", format: "text" },
  { key: "email", label: "Email", format: "email" },
  { key: "role", label: "Role", format: "badge" },
  { key: "active", label: "Active", format: "boolean" },
  { key: "created_at", label: "Created", format: "relative" },
];

function SystemUsersPage() {
  const { data = [], isLoading } = useQuery<UserRow[]>({
    queryKey: ["system", "users"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<{ data: Record<string, unknown>[] }>("/users?page_size=200");
        return (data.data ?? []).map((u) => ({
          ...u,
          id: String(u.id),
          name: [u.first_name, u.last_name].filter(Boolean).join(" ") || String(u.email ?? ""),
        })) as UserRow[];
      } catch {
        return [];
      }
    },
    refetchInterval: 60_000,
  });

  return (
    <div>
      <PageHeader title="Users" description="Accounts, roles and status" />
      <div className="mt-6">
        <DataTable<UserRow>
          title="Users"
          singular="User"
          columns={COLUMNS}
          rows={data}
          loading={isLoading}
          searchKeys={["name", "email", "role"]}
        />
      </div>
    </div>
  );
}
`
}

// desktopClientSystemActivityPage is the audit-log feed with summary cards.
func desktopClientSystemActivityPage() string {
	return `import { useMemo, useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { Activity, Flag, AlertTriangle, Users } from "lucide-react";
import { PageHeader } from "@/components/layout/page-header";
import { SystemStat, SectionCard, EmptyState, relTime } from "@/components/system-ui";
import { apiClient } from "@/lib/api-client";

export const Route = createFileRoute("/app/system/activity")({
  component: SystemActivityPage,
});

interface Row {
  id: string; action: string; severity: "info" | "warn" | "critical";
  summary: string; ip_address: string; user_id?: string; created_at: string;
}

const dot: Record<Row["severity"], string> = { info: "bg-info", warn: "bg-warning", critical: "bg-danger" };
const TABS = [
  { key: "all", label: "All" },
  { key: "flagged", label: "Flagged" },
  { key: "critical", label: "Critical" },
] as const;

function SystemActivityPage() {
  const [tab, setTab] = useState<(typeof TABS)[number]["key"]>("all");
  const [search, setSearch] = useState("");

  const stats = useQuery({
    queryKey: ["system", "activity-stats"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<{ data: { info: number; warn: number; critical: number; total: number } }>("/user-activity/stats");
        return data.data;
      } catch { return { info: 0, warn: 0, critical: 0, total: 0 }; }
    },
    refetchInterval: 60_000,
  });

  const feed = useQuery<Row[]>({
    queryKey: ["system", "activity-feed"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<{ data: Row[] }>("/user-activity?page_size=200");
        return data.data ?? [];
      } catch { return []; }
    },
    refetchInterval: 30_000,
  });

  const rows = useMemo(() => {
    const q = search.trim().toLowerCase();
    return (feed.data ?? []).filter((r) => {
      if (tab === "flagged" && r.severity === "info") return false;
      if (tab === "critical" && r.severity !== "critical") return false;
      if (q && !(r.summary + " " + r.action).toLowerCase().includes(q)) return false;
      return true;
    });
  }, [feed.data, tab, search]);

  const activeUsers = useMemo(
    () => new Set((feed.data ?? []).map((r) => r.user_id).filter(Boolean)).size,
    [feed.data],
  );

  return (
    <div>
      <PageHeader
        title="Activity"
        description="Audit log across the platform"
        actions={
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search…"
            className="w-56 rounded-lg border border-border bg-surface-2 px-3 py-2 text-[13px] text-foreground outline-none focus:border-accent"
          />
        }
      />

      <div className="mt-6 grid grid-cols-2 gap-4 lg:grid-cols-4">
        <SystemStat label="Events" value={stats.data?.total ?? 0} icon={Activity} />
        <SystemStat label="Flagged" value={(stats.data?.warn ?? 0) + (stats.data?.critical ?? 0)} icon={Flag} tone="warning" />
        <SystemStat label="Critical" value={stats.data?.critical ?? 0} icon={AlertTriangle} tone="danger" />
        <SystemStat label="Active users" value={activeUsers} icon={Users} tone="info" />
      </div>

      <div className="mt-4 flex gap-1">
        {TABS.map((t) => (
          <button
            key={t.key}
            onClick={() => setTab(t.key)}
            className={
              "rounded-lg px-3 py-1.5 text-[13px] font-medium transition-colors " +
              (tab === t.key ? "bg-accent/10 text-accent" : "text-foreground-secondary hover:bg-surface-hover")
            }
          >
            {t.label}
          </button>
        ))}
      </div>

      <div className="mt-4">
        <SectionCard title="Recent events" description={rows.length + " event" + (rows.length === 1 ? "" : "s")}>
          {feed.isLoading ? (
            <div className="px-5 py-12 text-center text-[13px] text-foreground-muted">Loading…</div>
          ) : rows.length === 0 ? (
            <EmptyState icon={Activity} title="No activity yet." hint="Events will appear here as they happen." />
          ) : (
            <ul className="divide-y divide-border">
              {rows.map((r) => (
                <li key={r.id} className="flex items-start gap-3 px-5 py-3 text-[13px]">
                  <span className="mt-1.5 inline-block h-2 w-2 shrink-0 rounded-full">
                    <span className={"block h-full w-full rounded-full " + dot[r.severity]} />
                  </span>
                  <div className="min-w-0 flex-1">
                    <p className="truncate text-foreground">{r.summary}</p>
                    <p className="text-xs text-foreground-muted">
                      <code className="font-mono">{r.action}</code>
                      {r.ip_address && <span>{" · "}{r.ip_address === "::1" || r.ip_address === "127.0.0.1" ? "localhost" : r.ip_address}</span>}
                    </p>
                  </div>
                  <span className="shrink-0 text-xs text-foreground-muted">{relTime(r.created_at)}</span>
                </li>
              ))}
            </ul>
          )}
        </SectionCard>
      </div>
    </div>
  );
}
`
}

// desktopClientSystemNotificationsPage lists system/security notifications.
func desktopClientSystemNotificationsPage() string {
	return `import { createFileRoute } from "@tanstack/react-router";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Bell, CheckCheck } from "lucide-react";
import { PageHeader } from "@/components/layout/page-header";
import { EmptyState, relTime } from "@/components/system-ui";
import { apiClient } from "@/lib/api-client";

export const Route = createFileRoute("/app/system/notifications")({
  component: SystemNotificationsPage,
});

interface Notification {
  id: string; source: "sentinel" | "pulse" | "system";
  severity: "critical" | "high" | "medium" | "low" | "info";
  title: string; body: string; count?: number; read_at?: string | null; created_at: string;
}

const sevTone: Record<Notification["severity"], string> = {
  critical: "bg-danger/10 text-danger",
  high: "bg-danger/10 text-danger",
  medium: "bg-warning/10 text-warning",
  low: "bg-info/10 text-info",
  info: "bg-accent/10 text-accent",
};

function SystemNotificationsPage() {
  const qc = useQueryClient();
  const { data, isLoading } = useQuery<{ data: Notification[]; unread: number }>({
    queryKey: ["system", "notifications"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<{ data: Notification[]; unread: number }>("/notifications");
        return { data: data.data ?? [], unread: data.unread ?? 0 };
      } catch { return { data: [], unread: 0 }; }
    },
    refetchInterval: 60_000,
  });

  const markAll = useMutation({
    mutationFn: () => apiClient.post("/notifications/read-all"),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["system", "notifications"] }),
  });
  const markOne = useMutation({
    mutationFn: (id: string) => apiClient.post("/notifications/" + id + "/read"),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["system", "notifications"] }),
  });

  const items = data?.data ?? [];

  return (
    <div>
      <PageHeader
        title="Notifications"
        description="System & security alerts"
        actions={
          (data?.unread ?? 0) > 0 ? (
            <button
              onClick={() => markAll.mutate()}
              className="flex items-center gap-2 rounded-lg border border-border px-3 py-2 text-[13px] text-foreground-secondary hover:bg-surface-hover"
            >
              <CheckCheck className="h-4 w-4" /> Mark all read
            </button>
          ) : undefined
        }
      />

      <div className="mt-6 rounded-xl border border-border bg-surface">
        {isLoading ? (
          <div className="px-5 py-12 text-center text-[13px] text-foreground-muted">Loading…</div>
        ) : items.length === 0 ? (
          <EmptyState icon={Bell} title="You're all caught up." hint="New alerts will show up here." />
        ) : (
          <ul className="divide-y divide-border">
            {items.map((n) => (
              <li
                key={n.id}
                onClick={() => !n.read_at && markOne.mutate(n.id)}
                className={
                  "flex cursor-pointer items-start gap-3 px-5 py-4 transition-colors hover:bg-surface-hover " +
                  (!n.read_at ? "bg-accent/5" : "")
                }
              >
                <span className={"mt-0.5 inline-flex h-8 w-8 items-center justify-center rounded-lg " + sevTone[n.severity]}>
                  <Bell className="h-4 w-4" />
                </span>
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <p className="truncate text-[13px] font-semibold text-foreground">{n.title}</p>
                    {(n.count ?? 0) > 1 && <span className="rounded-full bg-surface-2 px-1.5 text-[11px] text-foreground-muted">×{n.count}</span>}
                  </div>
                  <p className="truncate text-[12px] text-foreground-secondary">{n.body}</p>
                  <p className="mt-0.5 text-[11px] text-foreground-muted">{n.source} · {relTime(n.created_at)}</p>
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
`
}

// desktopClientSystemHealthPage renders the infrastructure health cards.
func desktopClientSystemHealthPage() string {
	return `import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import {
  Database, HardDrive, Server, Activity, Mail, RefreshCw,
  CheckCircle2, AlertCircle, type LucideIcon,
} from "lucide-react";
import { PageHeader } from "@/components/layout/page-header";
import { cn } from "@/lib/utils";
import { apiClient } from "@/lib/api-client";

export const Route = createFileRoute("/app/system/health")({
  component: SystemHealthPage,
});

interface Health {
  status: "ok" | "degraded" | "down";
  database?: { ok: boolean; latency_ms?: number };
  redis?: { ok: boolean; latency_ms?: number };
  api?: { ok: boolean };
  jobs?: { ok: boolean };
  email?: { ok: boolean; configured?: boolean };
}

function SystemHealthPage() {
  const { data, isFetching, refetch } = useQuery<Health>({
    queryKey: ["system", "health"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<Health>("/health");
        return data;
      } catch { return { status: "ok" }; }
    },
    refetchInterval: 30_000,
  });

  const ok = data?.status !== "down" && data?.status !== "degraded";
  const cards: { label: string; icon: LucideIcon; up?: boolean; detail: string }[] = [
    { label: "PostgreSQL", icon: Database, up: data?.database?.ok, detail: data?.database?.latency_ms != null ? data.database.latency_ms + "ms" : "Primary database" },
    { label: "Redis", icon: HardDrive, up: data?.redis?.ok, detail: data?.redis?.latency_ms != null ? data.redis.latency_ms + "ms" : "Cache & queue" },
    { label: "API Server", icon: Server, up: data?.api?.ok ?? true, detail: "HTTP gateway" },
    { label: "Background Jobs", icon: Activity, up: data?.jobs?.ok, detail: "Worker queue" },
    { label: "Email (Resend)", icon: Mail, up: data?.email?.ok, detail: data?.email?.configured === false ? "Not configured" : "Transactional mail" },
  ];

  return (
    <div>
      <PageHeader
        title="System Health"
        description="Live status of your infrastructure"
        actions={
          <div className="flex items-center gap-3">
            <span className={cn("rounded-full px-3 py-1 text-[12px] font-semibold", ok ? "bg-success/10 text-success" : "bg-warning/10 text-warning")}>
              {ok ? "All systems operational" : "Degraded"}
            </span>
            <button
              onClick={() => refetch()}
              className="flex items-center gap-2 rounded-lg border border-border px-3 py-2 text-[13px] text-foreground-secondary hover:bg-surface-hover"
            >
              <RefreshCw className={cn("h-4 w-4", isFetching && "animate-spin")} /> Run check
            </button>
          </div>
        }
      />

      <h2 className="mb-3 mt-6 text-sm font-semibold uppercase tracking-wide text-foreground-muted">Infrastructure</h2>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {cards.map((c) => {
          const up = c.up === undefined ? null : c.up;
          return (
            <div
              key={c.label}
              className={cn(
                "rounded-xl border p-5",
                up === true ? "border-success/30 bg-success/5" : up === false ? "border-danger/30 bg-danger/5" : "border-border bg-surface",
              )}
            >
              <div className="flex items-start justify-between">
                <span className="inline-flex h-9 w-9 items-center justify-center rounded-lg bg-surface-2 text-foreground-secondary">
                  <c.icon className="h-4 w-4" />
                </span>
                {up === true ? <CheckCircle2 className="h-5 w-5 text-success" /> : up === false ? <AlertCircle className="h-5 w-5 text-danger" /> : <span className="text-[11px] text-foreground-muted">N/A</span>}
              </div>
              <p className="mt-3 text-sm font-semibold text-foreground">{c.label}</p>
              <p className="text-[12px] text-foreground-muted">{c.detail}</p>
            </div>
          );
        })}
      </div>
    </div>
  );
}
`
}

// desktopClientSystemPerformancePage renders the golden-signal metrics.
func desktopClientSystemPerformancePage() string {
	return `import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { Gauge, TrendingUp, AlertCircle, Cpu } from "lucide-react";
import { PageHeader } from "@/components/layout/page-header";
import { SystemStat, SectionCard, EmptyState } from "@/components/system-ui";
import { apiClient } from "@/lib/api-client";

export const Route = createFileRoute("/app/system/performance")({
  component: SystemPerformancePage,
});

interface Summary {
  latency?: { p50: number; p95: number; p99: number; avg: number };
  traffic?: { throughput: number; total: number };
  errors?: { rate: number; active_open: number };
  saturation?: { goroutines: number; heap_mb: number; gc_cycles: number; cpu_cores: number };
  slowest_routes?: { route: string; method: string; requests: number; avg: number; p95: number; p99: number; error_rate: number }[];
}

function SystemPerformancePage() {
  const { data } = useQuery<Summary>({
    queryKey: ["system", "performance"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<{ data?: Summary } & Summary>("/admin/observability/summary");
        return (data.data ?? data) as Summary;
      } catch { return {}; }
    },
    refetchInterval: 30_000,
  });

  const lat = data?.latency, tr = data?.traffic, err = data?.errors, sat = data?.saturation;
  const routes = data?.slowest_routes ?? [];

  return (
    <div>
      <PageHeader title="Performance" description="Latency, traffic, errors & saturation" />

      <div className="mt-6 grid grid-cols-2 gap-4 lg:grid-cols-4">
        <SystemStat label="P95 latency" value={lat ? Math.round(lat.p95) + "ms" : "—"} icon={Gauge} sub={lat ? "avg " + Math.round(lat.avg) + "ms" : undefined} />
        <SystemStat label="Throughput" value={tr ? tr.throughput + "/s" : "—"} icon={TrendingUp} tone="info" sub={tr ? tr.total + " total" : undefined} />
        <SystemStat label="Error rate" value={err ? (err.rate * 100).toFixed(1) + "%" : "—"} icon={AlertCircle} tone={err && err.rate > 0 ? "danger" : "default"} sub={err ? err.active_open + " open" : undefined} />
        <SystemStat label="Goroutines" value={sat?.goroutines ?? "—"} icon={Cpu} sub={sat ? sat.heap_mb + "MB heap" : undefined} />
      </div>

      <div className="mt-6">
        <SectionCard title="Slowest routes" description="By P99 latency">
          {routes.length === 0 ? (
            <EmptyState icon={Gauge} title="No route metrics yet." hint="Metrics appear once the API serves traffic." />
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-[13px]">
                <thead>
                  <tr className="border-b border-border text-left text-[11px] uppercase tracking-wider text-foreground-muted">
                    <th className="px-5 py-2.5 font-medium">Route</th>
                    <th className="px-4 py-2.5 font-medium">Reqs</th>
                    <th className="px-4 py-2.5 font-medium">Avg</th>
                    <th className="px-4 py-2.5 font-medium">P95</th>
                    <th className="px-4 py-2.5 font-medium">P99</th>
                    <th className="px-4 py-2.5 font-medium">Err</th>
                  </tr>
                </thead>
                <tbody>
                  {routes.map((r, i) => (
                    <tr key={i} className="border-b border-border/50">
                      <td className="px-5 py-2.5 font-mono text-foreground">{r.method} {r.route}</td>
                      <td className="px-4 py-2.5 text-foreground-secondary">{r.requests}</td>
                      <td className="px-4 py-2.5 text-foreground-secondary">{Math.round(r.avg)}ms</td>
                      <td className="px-4 py-2.5 text-foreground-secondary">{Math.round(r.p95)}ms</td>
                      <td className="px-4 py-2.5 text-foreground-secondary">{Math.round(r.p99)}ms</td>
                      <td className={"px-4 py-2.5 " + (r.error_rate > 0 ? "text-danger" : "text-foreground-secondary")}>{(r.error_rate * 100).toFixed(1)}%</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </SectionCard>
      </div>
    </div>
  );
}
`
}

// desktopClientSystemSecurityPage renders bans / rate limits / threats.
func desktopClientSystemSecurityPage() string {
	return `import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { Shield, AlertCircle, AlertTriangle } from "lucide-react";
import { PageHeader } from "@/components/layout/page-header";
import { SystemStat, SectionCard, EmptyState, relTime } from "@/components/system-ui";
import { apiClient } from "@/lib/api-client";

export const Route = createFileRoute("/app/system/security")({
  component: SystemSecurityPage,
});

interface Summary {
  banned_ips_now: number; auto_bans_24h: number; rate_limited_last_hour: number;
  active_bans?: { ip: string; reason: string; level?: number; expires_at?: string }[];
  rate_limit_hits_5min?: { ip: string; hits: number; last_hit: string }[];
  recent_threats?: { id: string; type: string; ip: string; description: string; created_at: string }[];
}

const TIERS = ["1st offence · 5 hours", "2nd offence · 8 hours", "3rd offence · 24 hours", "4th+ offence · 7 days"];

function SystemSecurityPage() {
  const { data } = useQuery<Summary>({
    queryKey: ["system", "security"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<{ data?: Summary } & Summary>("/admin/security/summary");
        return (data.data ?? data) as Summary;
      } catch { return { banned_ips_now: 0, auto_bans_24h: 0, rate_limited_last_hour: 0 }; }
    },
    refetchInterval: 60_000,
  });

  const bans = data?.active_bans ?? [];
  const limits = data?.rate_limit_hits_5min ?? [];
  const threats = data?.recent_threats ?? [];

  return (
    <div>
      <PageHeader title="Security" description="Bans, rate limits & recent threats" />

      <div className="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-3">
        <SystemStat label="Banned IPs now" value={data?.banned_ips_now ?? 0} icon={Shield} tone={(data?.banned_ips_now ?? 0) > 0 ? "danger" : "default"} />
        <SystemStat label="Auto-bans (24h)" value={data?.auto_bans_24h ?? 0} icon={AlertCircle} tone={(data?.auto_bans_24h ?? 0) > 0 ? "warning" : "default"} />
        <SystemStat label="Rate-limited (1h)" value={data?.rate_limited_last_hour ?? 0} icon={AlertTriangle} tone={(data?.rate_limited_last_hour ?? 0) > 0 ? "warning" : "default"} />
      </div>

      <div className="mt-6 rounded-xl border border-border bg-surface p-5">
        <p className="text-sm font-semibold text-foreground">Escalating auto-ban policy</p>
        <p className="text-xs text-foreground-muted">Repeat offenders are banned for progressively longer.</p>
        <div className="mt-3 grid grid-cols-1 gap-2 sm:grid-cols-2 lg:grid-cols-4">
          {TIERS.map((t, i) => (
            <div key={i} className="rounded-lg border border-border bg-surface-2 px-3 py-2 text-[12px] text-foreground-secondary">{t}</div>
          ))}
        </div>
      </div>

      <div className="mt-6 grid grid-cols-1 gap-4 lg:grid-cols-2">
        <SectionCard title="Active IP bans">
          {bans.length === 0 ? (
            <EmptyState icon={Shield} title="No active bans." />
          ) : (
            <ul className="divide-y divide-border">
              {bans.map((b, i) => (
                <li key={i} className="flex items-center justify-between px-5 py-3 text-[13px]">
                  <div className="min-w-0">
                    <p className="font-mono text-foreground">{b.ip}</p>
                    <p className="truncate text-[12px] text-foreground-muted">{b.reason}</p>
                  </div>
                  <div className="text-right">
                    {b.level != null && <span className="rounded-full bg-danger/10 px-2 py-0.5 text-[11px] text-danger">L{b.level}</span>}
                    {b.expires_at && <p className="mt-0.5 text-[11px] text-foreground-muted">{relTime(b.expires_at)}</p>}
                  </div>
                </li>
              ))}
            </ul>
          )}
        </SectionCard>

        <SectionCard title="Rate-limit hits" description="Last 5 minutes">
          {limits.length === 0 ? (
            <EmptyState icon={AlertTriangle} title="No rate-limit hits." />
          ) : (
            <ul className="divide-y divide-border">
              {limits.map((l, i) => (
                <li key={i} className="flex items-center justify-between px-5 py-3 text-[13px]">
                  <span className="font-mono text-foreground">{l.ip}</span>
                  <div className="text-right">
                    <span className="rounded-full bg-warning/10 px-2 py-0.5 text-[11px] text-warning">{l.hits} hits</span>
                    <p className="mt-0.5 text-[11px] text-foreground-muted">{relTime(l.last_hit)}</p>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </SectionCard>
      </div>

      <div className="mt-4">
        <SectionCard title="Recent threats">
          {threats.length === 0 ? (
            <EmptyState icon={AlertTriangle} title="No threats detected." />
          ) : (
            <ul className="divide-y divide-border">
              {threats.map((t) => (
                <li key={t.id} className="px-5 py-3 text-[13px]">
                  <div className="flex items-center justify-between">
                    <span className="font-medium text-foreground">{t.type}</span>
                    <span className="text-[11px] text-foreground-muted">{relTime(t.created_at)}</span>
                  </div>
                  <p className="text-[12px] text-foreground-muted"><span className="font-mono">{t.ip}</span> · {t.description}</p>
                </li>
              ))}
            </ul>
          )}
        </SectionCard>
      </div>
    </div>
  );
}
`
}
