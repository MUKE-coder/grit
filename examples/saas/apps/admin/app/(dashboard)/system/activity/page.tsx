"use client";

import { useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { IconButton } from "@/components/ui/IconButton";
import { SkeletonCards, Skeleton } from "@/components/ui/Skeleton";
import { ResponsiveSheet } from "@/components/ui/ResponsiveSheet";
import { apiClient } from "@/lib/api-client";
import { exportToExcel } from "@/lib/export";
import { inputClasses } from "@/components/ui/input";
import {
  Activity as ActivityIcon, AlertCircle, AlertTriangle, Flag, Download,
  LogIn, LogOut, UserPlus, ShoppingCart, FileText, Settings as SettingsIcon, Shield,
  Users as UsersIcon,
} from "@/lib/icons";

interface ActivityRow {
  id: string;
  user_id: string;
  action: string;
  severity: "info" | "warn" | "critical";
  summary: string;
  resource_type: string;
  resource_id: string;
  ip_address: string;
  user_agent: string;
  metadata: string;
  created_at: string;
}
interface ListResponse { data: ActivityRow[] }
interface StatsResponse { data: { info: number; warn: number; critical: number; total: number } }

const TABS = ["All", "Flagged", "Critical", "After hours"] as const;
type Tab = (typeof TABS)[number];

// Chip palette + icon per action prefix. Falls back to "Event" with a
// generic activity icon when an action doesn't match any known prefix.
function actionChip(action: string): { label: string; icon: React.ReactNode; tone: string } {
  if (action.startsWith("auth.login_failed") || action === "auth.login_blocked")
    return { label: "Auth failed", icon: <Shield className="h-3.5 w-3.5" />, tone: "bg-danger/10 text-danger" };
  if (action.startsWith("auth.login"))
    return { label: "Sign-in", icon: <LogIn className="h-3.5 w-3.5" />, tone: "bg-success/10 text-success" };
  if (action.startsWith("auth.logout"))
    return { label: "Sign-out", icon: <LogOut className="h-3.5 w-3.5" />, tone: "bg-text-muted/10 text-text-secondary" };
  if (action.startsWith("auth.register"))
    return { label: "Sign-up", icon: <UserPlus className="h-3.5 w-3.5" />, tone: "bg-info/10 text-info" };
  if (action.startsWith("ticket"))
    return { label: "Ticket", icon: <FileText className="h-3.5 w-3.5" />, tone: "bg-warning/10 text-warning" };
  if (action.startsWith("sale") || action.startsWith("order"))
    return { label: "Sale", icon: <ShoppingCart className="h-3.5 w-3.5" />, tone: "bg-accent/10 text-accent" };
  if (action.startsWith("user"))
    return { label: "User", icon: <UsersIcon className="h-3.5 w-3.5" />, tone: "bg-info/10 text-info" };
  if (action.startsWith("settings"))
    return { label: "Settings", icon: <SettingsIcon className="h-3.5 w-3.5" />, tone: "bg-text-muted/10 text-text-secondary" };
  return { label: "Event", icon: <ActivityIcon className="h-3.5 w-3.5" />, tone: "bg-accent/10 text-accent" };
}

export default function ActivityPage() {
  const [tab, setTab] = useState<Tab>("All");
  const [search, setSearch] = useState("");
  const [from, setFrom] = useState("");
  const [to, setTo] = useState("");
  const [detail, setDetail] = useState<ActivityRow | null>(null);

  const { data: stats } = useQuery<StatsResponse["data"]>({
    queryKey: ["user-activity", "stats"],
    queryFn: async () => {
      const { data } = await apiClient.get<StatsResponse>("/api/user-activity/stats");
      return data.data;
    },
    refetchInterval: 60_000,
  });

  const { data: rows, isLoading } = useQuery<ActivityRow[]>({
    queryKey: ["user-activity", "feed", search, tab, from, to],
    queryFn: async () => {
      const params = new URLSearchParams({ page_size: "200" });
      if (search) params.set("q", search);
      if (tab === "Critical") params.set("severity", "critical");
      const { data } = await apiClient.get<ListResponse>("/api/user-activity?" + params.toString());
      let out = data.data;
      // Client-side filters for Flagged + After-hours tabs and date range
      // since the API doesn't currently expose those facets — cheap on
      // the page-size we request, and keeps the handler small.
      if (tab === "Flagged") out = out.filter((r) => r.severity !== "info");
      if (tab === "After hours") out = out.filter((r) => {
        const h = new Date(r.created_at).getHours();
        return h < 7 || h >= 19;
      });
      if (from) out = out.filter((r) => new Date(r.created_at) >= new Date(from));
      if (to) {
        const toDate = new Date(to);
        toDate.setHours(23, 59, 59, 999);
        out = out.filter((r) => new Date(r.created_at) <= toDate);
      }
      return out;
    },
  });

  // Group rows by their day (e.g. "TODAY · SAT, JUN 20, 2026 · 12 EVENTS").
  const grouped = useMemo(() => {
    const groups: Record<string, ActivityRow[]> = {};
    for (const r of rows ?? []) {
      const key = new Date(r.created_at).toDateString();
      if (!groups[key]) groups[key] = [];
      groups[key].push(r);
    }
    return Object.entries(groups);
  }, [rows]);

  // Active-users count is sourced from the visible rows so the header
  // stays in sync with whatever the user has filtered to.
  const activeUsers = useMemo(() => {
    const ids = new Set<string>();
    for (const r of rows ?? []) if (r.user_id) ids.add(r.user_id);
    return ids.size;
  }, [rows]);

  const onExport = async () => {
    const payload = (rows ?? []).map((r) => ({
      Time: new Date(r.created_at).toLocaleString(),
      Severity: r.severity,
      Action: r.action,
      Summary: r.summary,
      User: r.user_id,
      Resource: r.resource_type + ":" + r.resource_id,
      IP: r.ip_address,
      UserAgent: r.user_agent,
    }));
    await exportToExcel(payload, "user-activity-" + new Date().toISOString().slice(0, 10));
  };

  return (
    <div>
      <PageHeader
        title="Activity"
        subtitle="Descriptive feed of what each user did — sign-ins, writes, security events."
        searchPlaceholder="Search the feed..."
        searchValue={search}
        onSearchChange={setSearch}
      />

      {/* Tabs */}
      <div className="mb-6 flex flex-wrap gap-1 border-b border-border">
        {TABS.map((t) => {
          const active = t === tab;
          const Icon = t === "Flagged" ? Flag : t === "Critical" ? AlertCircle : t === "After hours" ? AlertTriangle : ActivityIcon;
          const count = t === "All" ? (rows?.length ?? 0) : undefined;
          return (
            <button
              key={t}
              type="button"
              onClick={() => setTab(t)}
              className={
                "inline-flex items-center gap-2 border-b-2 px-4 py-2 text-sm font-medium transition-colors -mb-px " +
                (active
                  ? "border-accent text-accent"
                  : "border-transparent text-text-secondary hover:text-foreground")
              }
            >
              <Icon className="h-4 w-4" />
              {t === "All" && count !== undefined ? "All " + count : t}
            </button>
          );
        })}
      </div>

      {/* Summary cards */}
      {!stats ? (
        <SkeletonCards count={4} />
      ) : (
        <div className="mb-6 grid grid-cols-2 gap-3 md:grid-cols-4">
          <SummaryCard label="Events" value={rows?.length ?? 0} icon={<ActivityIcon className="h-4 w-4" />} tone="default" />
          <SummaryCard label="Flagged" value={(stats.warn || 0) + (stats.critical || 0)} icon={<Flag className="h-4 w-4" />} tone="warning" />
          <SummaryCard label="Critical" value={stats.critical} icon={<AlertCircle className="h-4 w-4" />} tone="danger" />
          <SummaryCard label="Active users" value={activeUsers} icon={<UsersIcon className="h-4 w-4" />} tone="info" />
        </div>
      )}

      {/* Filter bar */}
      <div className="mb-6 grid grid-cols-1 gap-3 rounded-xl border border-border bg-bg-elevated p-4 md:grid-cols-[1fr_1fr_auto] md:items-end">
        <Field label="From">
          <input
            type="date"
            value={from}
            onChange={(e) => setFrom(e.target.value)}
            className={inputClasses()}
          />
        </Field>
        <Field label="To">
          <input
            type="date"
            value={to}
            onChange={(e) => setTo(e.target.value)}
            className={inputClasses()}
          />
        </Field>
        <IconButton
          variant="secondary"
          icon={<Download className="h-4 w-4" />}
          label="Export"
          onClick={onExport}
        />
      </div>

      {/* Feed */}
      {isLoading ? (
        <div className="space-y-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="rounded-xl border border-border bg-bg-elevated p-4">
              <Skeleton shape="text" className="w-1/2 mb-2" />
              <Skeleton shape="text" className="w-3/4" />
            </div>
          ))}
        </div>
      ) : grouped.length === 0 ? (
        <div className="rounded-xl border border-border bg-bg-elevated p-12 text-center">
          <ActivityIcon className="mx-auto h-10 w-10 text-text-muted" />
          <p className="mt-3 text-base font-medium text-foreground">No matching events</p>
          <p className="mt-1 text-sm text-text-muted">Try a different tab or widen the date range.</p>
        </div>
      ) : (
        <div className="space-y-6">
          {grouped.map(([day, items]) => {
            const date = new Date(day);
            const today = new Date(); today.setHours(0,0,0,0);
            const target = new Date(date); target.setHours(0,0,0,0);
            const isToday = today.getTime() === target.getTime();
            return (
              <section key={day}>
                <p className="mb-2 text-xs font-semibold uppercase tracking-wider text-text-muted">
                  {isToday ? "Today" : "On"} · {date.toLocaleDateString(undefined, { weekday: "short", month: "short", day: "numeric", year: "numeric" })} · {items.length} event{items.length === 1 ? "" : "s"}
                </p>
                <ul className="space-y-2">
                  {items.map((row) => {
                    const chip = actionChip(row.action);
                    return (
                      <li
                        key={row.id}
                        className="flex items-start gap-4 rounded-xl border border-border bg-bg-elevated px-4 py-3"
                      >
                        <span className="shrink-0 pt-0.5 text-xs font-mono text-text-muted">
                          {new Date(row.created_at).toLocaleTimeString(undefined, { hour: "2-digit", minute: "2-digit" })}
                        </span>
                        <span className={"inline-flex shrink-0 items-center gap-1 rounded-md px-2 py-0.5 text-xs font-semibold " + chip.tone}>
                          {chip.icon}
                          {chip.label}
                        </span>
                        <div className="min-w-0 flex-1">
                          <p className="truncate text-sm font-medium text-foreground">{row.summary}</p>
                          <p className="truncate text-xs text-text-muted">
                            <code className="font-mono">{row.action}</code>
                            {row.ip_address && <span> · {prettyIP(row.ip_address)}</span>}
                          </p>
                        </div>
                        <button
                          type="button"
                          onClick={() => setDetail(row)}
                          className="shrink-0 rounded-md border border-border bg-bg-secondary px-2.5 py-1 text-xs font-medium text-text-secondary hover:bg-bg-hover hover:text-foreground"
                        >
                          Info
                        </button>
                      </li>
                    );
                  })}
                </ul>
              </section>
            );
          })}
        </div>
      )}

      {/* Detail drawer */}
      <ResponsiveSheet
        open={detail !== null}
        onClose={() => setDetail(null)}
        title={detail?.summary || "Event"}
        description={detail ? new Date(detail.created_at).toLocaleString() : undefined}
        size="lg"
      >
        {detail && (
          <dl className="space-y-3 text-sm">
            <KV label="Action" value={<code className="font-mono">{detail.action}</code>} />
            <KV label="Severity" value={detail.severity.toUpperCase()} />
            <KV label="Resource" value={detail.resource_type ? detail.resource_type + " · " + detail.resource_id : "—"} />
            <KV label="User" value={detail.user_id || "system"} />
            <KV label="IP" value={prettyIP(detail.ip_address)} />
            <KV label="User agent" value={<span className="text-xs text-text-muted">{detail.user_agent || "—"}</span>} />
            {detail.metadata && (
              <div>
                <p className="mb-1 text-xs font-semibold uppercase tracking-wide text-text-muted">Metadata</p>
                <pre className="overflow-x-auto rounded-lg bg-bg-secondary p-3 text-xs">{prettyJSON(detail.metadata)}</pre>
              </div>
            )}
          </dl>
        )}
      </ResponsiveSheet>
    </div>
  );
}

// v3.31.49 -- IPv6 loopback (::1) and IPv4 loopback (127.0.0.1) are
// the actual IP for local-dev requests; gin.Context.ClientIP() does
// the right thing. But "::1" reads as cryptic to operators -- and
// every dev sees it on every event. Show "localhost" with the raw
// value tucked next to it so the origin is obvious and inspectable.
function prettyIP(ip: string | null | undefined): React.ReactNode {
  if (!ip) return "—";
  if (ip === "::1" || ip === "127.0.0.1" || ip === "0.0.0.0") {
    return (
      <span title={ip}>
        localhost <span className="text-xs text-text-muted">({ip})</span>
      </span>
    );
  }
  return ip;
}

function SummaryCard({ label, value, icon, tone }: { label: string; value: number; icon: React.ReactNode; tone: "default" | "warning" | "danger" | "info" }) {
  const toneClass = {
    default: "border-border bg-bg-elevated",
    warning: "border-warning/30 bg-warning/5",
    danger:  "border-danger/30 bg-danger/5",
    info:    "border-info/30 bg-info/5",
  }[tone];
  const iconClass = {
    default: "text-text-secondary",
    warning: "text-warning",
    danger:  "text-danger",
    info:    "text-info",
  }[tone];
  return (
    <div className={"rounded-xl border p-4 " + toneClass}>
      <div className="flex items-center justify-between">
        <p className="text-xs font-semibold uppercase tracking-wide text-text-muted">{label}</p>
        <span className={iconClass}>{icon}</span>
      </div>
      <p className="mt-2 text-2xl font-bold text-foreground">{value}</p>
    </div>
  );
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <label className="block">
      <span className="mb-1 block text-xs font-semibold uppercase tracking-wide text-text-muted">{label}</span>
      {children}
    </label>
  );
}

function KV({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className="grid grid-cols-3 gap-3 border-b border-border pb-2 last:border-b-0">
      <dt className="text-xs font-semibold uppercase tracking-wide text-text-muted">{label}</dt>
      <dd className="col-span-2 text-sm text-foreground">{value}</dd>
    </div>
  );
}

function prettyJSON(s: string): string {
  try { return JSON.stringify(JSON.parse(s), null, 2); } catch { return s; }
}
