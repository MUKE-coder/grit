package scaffold

// v3.31.5 page redesigns + new System surfaces.
//
//   - System hub (/system) — redesigned tiles matching the Shoppleet Setup
//     look the user shared: clean rounded card, icon top-left, title +
//     description, no garish tone colours.
//   - System Health (/system/health) — Infrastructure status cards for
//     PostgreSQL / Redis / API server / Background Jobs / Email. No
//     Payment Lines — that's app-specific and lives in app code.
//   - Security (/system/security) — DGateway Sentinel-inspired KPIs,
//     auto-ban escalation policy explainer, active bans / recent threats.
//   - Performance (/system/performance) — DGateway Pulse-inspired four
//     golden signals layout (latency, traffic, errors, saturation).
//
// All four pages use the v3.29 PageHeader + Skeleton primitives.

// adminSystemHubPageV2 — redesigned system hub. Plain cards (no tone
// tinting), icon in a square at the top-left, title + description. Same
// landing tile set as v3.30 but presented more clearly.
func adminSystemHubPageV2() string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { useModules } from "@/hooks/use-modules";
import { PageHeader } from "@/components/chrome/PageHeader";
import { SYSTEM_NAV, INTERNAL_ICON } from "@/components/chrome/CollapsibleSidebar";
import {
  Activity, Bell, Calendar, Database, FileText, Mail,
  MessageSquare, Shield, ShieldCheck, TrendingUp, Upload, Link as LinkIcon,
  UserCheck, Settings, LayoutGrid, KeyRound, Workflow,
} from "@/lib/icons";

// Every operational surface is grouped under one of these tabs. The whole
// point of the hub is to keep the sidebar to just resources + one "System
// Hub" link, so ALL of these live here rather than in the rail.
type Category =
  | "Operations"
  | "Security & Access"
  | "Data & Files"
  | "Communication"
  | "Settings";

const CATEGORIES: Category[] = [
  "Operations",
  "Security & Access",
  "Data & Files",
  "Communication",
  "Settings",
];

// Plugin-injected system links (Webhooks, Impersonate, …) that aren't built-in
// tiles show under this extra tab, which only appears when something injected.
const EXTENSIONS = "Extensions";

interface SystemTile {
  href: string;
  /** Optional module; the tile is hidden when that module is disabled. */
  module?: string;
  category: Category;
  title: string;
  description: string;
  icon: React.ReactNode;
}

const TILES: SystemTile[] = [
  // ── Operations ──────────────────────────────────────────────────────────
  { href: "/system/health",        category: "Operations", title: "System Health",   description: "Real-time infrastructure status — Postgres, Redis, API, jobs, email.",    icon: <Activity className="h-5 w-5" /> },
  { href: "/system/performance",   category: "Operations", title: "Performance",     description: "Four Google SRE golden signals — latency, traffic, errors, saturation.", icon: <TrendingUp className="h-5 w-5" /> },
  { href: "/system/observability", category: "Operations", title: "Observability",   description: "Pulse summary — latency, SLOs, top N+1, runtime.",                        icon: <TrendingUp className="h-5 w-5" /> },
  { href: "/system/jobs",          category: "Operations", title: "Background Jobs",  description: "Queue depth, in-flight workers, dead-letter queue.",                      icon: <Database className="h-5 w-5" /> , module: "jobs" },
  { href: "/system/sagas",         category: "Operations", title: "Sagas",           description: "Multi-step processes, and the ones whose undo could not be completed.", icon: <Workflow className="h-5 w-5" /> },
  { href: "/system/cron",          category: "Operations", title: "Cron Schedules",  description: "Recurring jobs, next-run times, run history.",                            icon: <Calendar className="h-5 w-5" /> , module: "cron" },
  { href: "/system/activity",      category: "Operations", title: "User Activity",   description: "Auth events, writes, operator actions with IP + severity.",               icon: <Activity className="h-5 w-5" /> , module: "audit" },
  // ── Security & Access ───────────────────────────────────────────────────
  { href: "/system/account",        category: "Security & Access", title: "Account",             description: "Your password, two-factor, passkeys and the devices you are signed in on.",  icon: <UserCheck className="h-5 w-5" /> },
  { href: "/system/security",       category: "Security & Access", title: "Security",            description: "Sentinel summary — banned IPs, rate-limit pressure, recent threats.",       icon: <Shield className="h-5 w-5" /> },
  { href: "/system/roles",          category: "Security & Access", title: "Roles & permissions", description: "Define what each role can see and do, and assign roles to users.",           icon: <ShieldCheck className="h-5 w-5" /> },
  { href: "/system/access-reviews", category: "Security & Access", title: "Access Reviews",      description: "Point-in-time snapshots of who holds which role — for audits and attestations.", icon: <UserCheck className="h-5 w-5" /> },
  { href: "/system/gdpr",           category: "Security & Access", title: "GDPR",                description: "Export or erase a user's data, with a tamper-evident deletion journal.",     icon: <Shield className="h-5 w-5" /> },
  { href: "/system/api-keys",       category: "Security & Access", title: "API keys",            description: "Long-lived credentials for scripts and integrations, instead of a login.", icon: <KeyRound className="h-5 w-5" /> },
  { href: "/system/audit",          category: "Security & Access", title: "Audit log",           description: "Every authenticated write, hash-chained so a changed row can be detected.", icon: <ShieldCheck className="h-5 w-5" /> },
  { href: "/system/sso",            category: "Security & Access", title: "Single sign-on",      description: "Let a customer's team sign in with their own identity provider — OIDC or SAML.",     icon: <ShieldCheck className="h-5 w-5" /> },
  // ── Data & Files ────────────────────────────────────────────────────────
  { href: "/system/backups",     category: "Data & Files", title: "Data & Backup",       description: "Create, schedule and restore database backups.",                                     icon: <Database className="h-5 w-5" /> , module: "backup" },
  { href: "/system/files",       category: "Data & Files", title: "File Storage",        description: "Browse uploads, manage retention, audit usage.",                                      icon: <Upload className="h-5 w-5" /> , module: "files" },
  { href: "/system/form-shares", category: "Data & Files", title: "Public form sharing", description: "Token-gated public submission links. Generate, enable/disable, view submissions.",     icon: <LinkIcon className="h-5 w-5" /> },
  // ── Communication ───────────────────────────────────────────────────────
  { href: "/system/mail",          category: "Communication", title: "Mail Preview",   description: "Every email template, rendered by the API.",       icon: <Mail className="h-5 w-5" /> , module: "mail" },
  { href: "/system/support",       category: "Communication", title: "Support",        description: "Incoming tickets, threads, assignments, closures.", icon: <MessageSquare className="h-5 w-5" /> },
  { href: "/system/notifications", category: "Communication", title: "Notifications",  description: "Recent system + Sentinel + Pulse notifications.",   icon: <Bell className="h-5 w-5" /> },
  // ── Settings ────────────────────────────────────────────────────────────
  { href: "/settings/dashboard",   category: "Settings", title: "Dashboard settings", description: "Configure the admin dashboard — widgets, stats, and layout.", icon: <Settings className="h-5 w-5" /> },
];

function Tile({ href, icon, title, description }: { href: string; icon: React.ReactNode; title: string; description: string }) {
  return (
    <Link
      href={href}
      className="group rounded-xl border border-border bg-bg-elevated p-5 transition-colors hover:bg-bg-hover hover:border-accent/30"
    >
      <div className="mb-3 inline-flex h-10 w-10 items-center justify-center rounded-lg bg-accent/10 text-accent">
        {icon}
      </div>
      <p className="text-base font-semibold text-foreground group-hover:text-accent">{title}</p>
      <p className="mt-1 text-sm text-text-secondary">{description}</p>
    </Link>
  );
}

export default function SystemHubPage() {
  // Hide tiles for modules that are switched off — a tile linking to a route
  // the server no longer mounts just 404s.
  const { moduleEnabled } = useModules();
  const tiles = TILES.filter((t) => !t.module || moduleEnabled(t.module));

  // Anything a plugin injected into SYSTEM_NAV that isn't one of our built-in
  // tiles (matched by href) becomes an "Extensions" entry.
  const builtinHrefs = new Set(TILES.map((t) => t.href));
  const extensions = SYSTEM_NAV.filter(
    (n) => !builtinHrefs.has(n.href) && n.href !== "/system" && (!n.module || moduleEnabled(n.module))
  );

  const tabs: string[] = [...CATEGORIES, ...(extensions.length ? [EXTENSIONS] : [])];
  const [active, setActive] = useState<string>(CATEGORIES[0]);
  const current = tabs.includes(active) ? active : CATEGORIES[0];

  return (
    <div>
      <PageHeader
        title="System"
        subtitle="Every operational surface for this app, grouped. Pick a tab, then a tile."
      />

      {/* Category tabs */}
      <div className="mb-6 flex flex-wrap gap-1 border-b border-border">
        {tabs.map((t) => (
          <button
            key={t}
            type="button"
            onClick={() => setActive(t)}
            className={
              "relative px-4 py-2.5 text-sm font-medium transition-colors " +
              (current === t ? "text-accent" : "text-text-secondary hover:text-foreground")
            }
          >
            {t}
            {current === t && (
              <span className="absolute inset-x-2 -bottom-px h-0.5 rounded-full bg-accent" />
            )}
          </button>
        ))}
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
        {current === EXTENSIONS
          ? extensions.map((n) => (
              <Tile
                key={n.href}
                href={n.href}
                icon={INTERNAL_ICON[n.iconKey] ?? <LayoutGrid className="h-5 w-5" />}
                title={n.label}
                description="Added by a Grit plugin."
              />
            ))
          : tiles
              .filter((t) => t.category === current)
              .map((t) => (
                <Tile key={t.href} href={t.href} icon={t.icon} title={t.title} description={t.description} />
              ))}
      </div>

      <FileText className="hidden" />
    </div>
  );
}
`
}

// adminSystemHealthPage — Infrastructure status. Hits /api/health for the
// summary + uses fixed dummy/derived data for individual components when
// the API doesn't yet break them out. Each card is green when up, danger
// when down. Top-right "All Systems Operational" pill + "Run Health Check"
// button that re-runs the query.
func adminSystemHealthPage() string {
	return `"use client";

import { useQuery, useQueryClient } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { SkeletonCards } from "@/components/ui/Skeleton";
import { IconButton } from "@/components/ui/IconButton";
import { apiClient } from "@/lib/api-client";
import {
  CheckCircle, AlertCircle, AlertTriangle, RefreshCw, Database, Mail, Server,
  Activity as ActivityIcon, HardDrive, Zap, FolderOpen, Webhook,
} from "@/lib/icons";

// The four states /api/health answers in.
//
// "off" is a component this deployment never configured, which is not a
// problem. "unknown" is a probe that could not find out, which is not an
// accusation either. Only "degraded" is bad news.
//
// This page used to reconstruct the state from a boolean, and did it
// differently per component: Redis read ok === false ? down : ok ? ok :
// unknown, email read configured ? ok : unknown, and storage was not on the
// page at all. Nothing here infers a state any more.
type HealthState = "ok" | "degraded" | "off" | "unknown";

interface Component {
  state?: HealthState;
  // The state as a boolean. Kept because an API that has not been upgraded yet
  // sends only this, and because a load balancer reads it.
  ok?: boolean;
  detail?: string;
  latency_ms?: number;
  tables?: number;
  queued?: number;
  active?: number;
  configured?: boolean;
  driver?: string;
  subscribers?: number;
  capacity?: number;
  dropped?: number;
}

interface HealthResponse {
  status: HealthState;
  database?: Component;
  redis?: Component;
  api?: Component;
  jobs?: Component;
  email?: Component;
  storage?: Component;
  events?: Component;
  // This replica's hub. With several replicas each reports its own share, so
  // the numbers are per process, not per deployment.
  realtime?: {
    connections: number;
    users: number;
    channels: number;
    messages_sent: number;
    messages_dropped: number;
    backplane: boolean;
    backplane_publish_errors: number;
  };
}

// stateOf reads a component's state, with a fallback for an API older than it.
//
// A component that is not in the response was not reported at all, so unknown.
// An older API sends only ok, where false covered both "down" and "never
// configured": unknown is the honest reading, and reading it as down is what
// put a red card on every machine that runs without Redis.
function stateOf(component?: Component): HealthState {
  if (!component) return "unknown";
  if (component.state) return component.state;
  return component.ok ? "ok" : "unknown";
}

interface Card {
  key: string;
  label: string;
  icon: React.ReactNode;
  state: HealthState;
  detail: string;
  meta?: string;
}

export default function SystemHealthPage() {
  const queryClient = useQueryClient();
  const { data, isLoading, isError, refetch, isFetching } = useQuery<HealthResponse>({
    queryKey: ["system-health"],
    // No catch here. This page used to answer a failed request with a benign
    // { status: "ok" } object, so an API that was down painted "All Systems
    // Operational", which is the one thing a status page must never do.
    queryFn: async () => {
      const { data } = await apiClient.get<HealthResponse>("/api/health");
      return data;
    },
    refetchInterval: 30_000,
  });

  const overall: HealthState = isError ? "degraded" : (data?.status ?? "unknown");
  const cards: Card[] = data ? [
    {
      key: "database",
      label: "Database",
      icon: <Database className="h-5 w-5" />,
      state: stateOf(data.database),
      detail: data.database?.detail
        ?? (data.database?.tables ? data.database.tables + " tables, ping OK" : "Ping OK"),
      meta: data.database?.latency_ms != null ? data.database.latency_ms + "ms" : undefined,
    },
    {
      key: "redis",
      label: "Redis",
      icon: <HardDrive className="h-5 w-5" />,
      state: stateOf(data.redis),
      detail: data.redis?.detail ?? "Ping OK",
      meta: data.redis?.latency_ms != null ? data.redis.latency_ms + "ms" : undefined,
    },
    {
      key: "api",
      label: "API Server",
      icon: <Server className="h-5 w-5" />,
      state: stateOf(data.api),
      detail: data.api?.detail ?? "Responding to requests",
    },
    {
      key: "jobs",
      label: "Background Jobs",
      icon: <ActivityIcon className="h-5 w-5" />,
      state: stateOf(data.jobs),
      detail: data.jobs?.queued != null
        ? data.jobs.queued + " queued, " + (data.jobs.active ?? 0) + " running"
        : (data.jobs?.detail ?? "Worker pool healthy"),
    },
    {
      key: "storage",
      label: "File Storage",
      icon: <FolderOpen className="h-5 w-5" />,
      state: stateOf(data.storage),
      detail: data.storage?.detail
        ?? (data.storage?.driver ? "Storing files on " + data.storage.driver : "Configured"),
    },
    {
      key: "email",
      label: "Email",
      icon: <Mail className="h-5 w-5" />,
      state: stateOf(data.email),
      detail: data.email?.detail
        ?? (data.email?.driver ? "Sending with " + data.email.driver : "Configured"),
    },
    {
      key: "events",
      label: "Event Bus",
      icon: <Webhook className="h-5 w-5" />,
      state: stateOf(data.events),
      detail: data.events?.detail ?? ((data.events?.subscribers ?? 0) + " subscribers"),
      // Dropped is a counter for the life of the process, so it is the trend
      // that matters: a subscriber too slow, or a queue too small.
      meta: data.events?.dropped != null
        ? data.events.dropped + " dropped / " + (data.events.capacity ?? 0) + " capacity"
        : undefined,
    },
    {
      key: "realtime",
      label: "Realtime",
      icon: <Zap className="h-5 w-5" />,
      state: data.realtime ? "ok" : "unknown",
      detail: data.realtime
        ? data.realtime.connections + " sockets, " + data.realtime.users + " users, " + data.realtime.channels + " channels"
        : "Not reported",
      meta: data.realtime
        ? data.realtime.messages_dropped + " dropped / " + data.realtime.messages_sent + " sent"
        : undefined,
    },
  ] : [];

  return (
    <div>
      <PageHeader
        title="System Health"
        subtitle="Real-time status of every platform component."
        actions={
          <>
            <span
              className={
                "hidden md:inline-flex items-center gap-2 rounded-full border px-3 py-1.5 text-xs font-semibold " +
                (overall === "ok"
                  ? "border-success/30 bg-success/5 text-success"
                  : "border-warning/30 bg-warning/5 text-warning")
              }
            >
              {overall === "ok"
                ? <CheckCircle className="h-3.5 w-3.5" />
                : <AlertTriangle className="h-3.5 w-3.5" />}
              {overall === "ok" ? "All Systems Operational" : "Degraded, review components"}
            </span>
            <IconButton
              variant="secondary"
              icon={<RefreshCw className={"h-4 w-4 " + (isFetching ? "animate-spin" : "")} />}
              label="Run Health Check"
              onClick={() => { refetch(); queryClient.invalidateQueries({ queryKey: ["system-health"] }); }}
            />
          </>
        }
      />

      {isError && (
        <div className="mb-4 flex items-start gap-2 rounded-lg border border-danger/30 bg-danger/5 px-4 py-3 text-sm text-danger">
          <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0" />
          <span>
            The API did not answer /api/health. Nothing can be reported about the
            components behind it, which is itself the finding.
          </span>
        </div>
      )}

      <h2 className="mb-3 text-xl font-bold text-foreground">Infrastructure</h2>

      {isLoading ? (
        <SkeletonCards count={8} />
      ) : (
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
          {cards.map((c) => (
            <HealthCard key={c.key} card={c} />
          ))}
        </div>
      )}

      <p className="mt-4 max-w-2xl text-xs text-text-muted">
        Off is a component this deployment never configured, and unknown is a
        probe that could not measure one. Neither counts against the status
        above: only degraded does.
      </p>
    </div>
  );
}

// A state's colours. Off and unknown get none, because a deployment that runs
// without Redis should not show a red card: that is how an operator learns to
// ignore this page.
const TONES: Record<HealthState, { card: string; icon: string }> = {
  ok: { card: "border-success/30 bg-success/5", icon: "text-success" },
  degraded: { card: "border-danger/30 bg-danger/5", icon: "text-danger" },
  off: { card: "border-border bg-bg-elevated", icon: "text-text-muted" },
  unknown: { card: "border-border bg-bg-elevated", icon: "text-text-muted" },
};

function StateBadge({ state }: { state: HealthState }) {
  if (state === "ok") return <CheckCircle className="h-4 w-4 text-success" />;
  if (state === "degraded") return <AlertCircle className="h-4 w-4 text-danger" />;
  // Off and unknown are words, not symbols. A dash beside "off" reads as a
  // missing value, which is the opposite of what off means here.
  return (
    <span className="text-[10px] font-semibold uppercase text-text-muted">
      {state === "off" ? "Off" : "Unknown"}
    </span>
  );
}

function HealthCard({ card }: { card: Card }) {
  const tone = TONES[card.state];

  return (
    <div className={"rounded-xl border p-4 " + tone.card}>
      <div className="mb-3 flex items-center justify-between">
        <span className={"inline-flex h-9 w-9 items-center justify-center rounded-lg bg-bg-elevated " + tone.icon}>
          {card.icon}
        </span>
        <StateBadge state={card.state} />
      </div>
      <p className="text-sm font-semibold text-foreground">{card.label}</p>
      <p className="mt-1 text-xs text-text-secondary">{card.detail}</p>
      {card.meta && (
        <p className="mt-2 inline-flex items-center gap-1 rounded bg-bg-elevated px-1.5 py-0.5 text-[10px] font-mono text-text-muted">
          {card.meta}
        </p>
      )}
    </div>
  );
}
`
}

// adminSecurityPageV2 — DGateway Sentinel-inspired security summary.
//
// Three KPI cards (currently banned IPs, auto-bans last 24h, rate-limited
// IPs last hour) + a section explaining the escalating auto-ban policy
// (5h on first hit, doubles on repeat) + Active IP bans table + IPs hit
// rate limits in the last 5 minutes + recent threats. "Open full Sentinel"
// pill in the top-right that deep-links to /sentinel/ui.
func adminSecurityPageV2() string {
	return `"use client";

import { useQuery } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { SkeletonCards } from "@/components/ui/Skeleton";
import { apiClient } from "@/lib/api-client";
import {
  Shield, AlertTriangle, AlertCircle, ExternalLink, Activity as ActivityIcon, Clock,
} from "@/lib/icons";

// The Sentinel UI is mounted on the Go API, not on this admin host. Use
// the API base so "Open Sentinel" works whether the admin is on :3001
// in dev or a different origin in prod.
import { API_URL } from "@/lib/api-core";

// DegradedBanner says which panels have no data behind them.
//
// The dashboard used to answer a dead upstream with a page of zeros and a 200,
// which reads as good news. The API now names the calls that failed, and this
// puts that on the screen above the numbers it applies to.
function DegradedBanner({ source, degraded }: { source: string; degraded?: string[] }) {
  if (!degraded || degraded.length === 0) return null;
  return (
    <div className="mt-4 flex items-start gap-2 rounded-lg border border-warning/30 bg-warning/5 px-4 py-3 text-sm text-warning">
      <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0" />
      <span>
        {source} did not answer for {degraded.join(", ")}. Those panels are empty rather than zero.
      </span>
    </div>
  );
}

interface SecuritySummary {
  banned_ips_now: number;
  auto_bans_24h: number;
  rate_limited_last_hour: number;
  active_bans?: Array<{ ip: string; reason: string; expires_at: string; level: number }>;
  rate_limit_hits_5min?: Array<{ ip: string; hits: number; last_hit: string }>;
  recent_threats?: Array<{ id: string; type: string; ip: string; description: string; created_at: string }>;
  // The calls to Sentinel that did not answer. Empty is the healthy case.
  degraded?: string[];
}

export default function SecurityPage() {
  const { data, isLoading } = useQuery<SecuritySummary>({
    queryKey: ["security", "summary"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<SecuritySummary>("/api/admin/security/summary");
        return data;
      } catch {
        // Nothing came back. A page of zeros would read as "no IP is banned and
        // nothing has attacked you", so say what actually happened instead.
        return { banned_ips_now: 0, auto_bans_24h: 0, rate_limited_last_hour: 0, degraded: ["the whole summary"] };
      }
    },
    refetchInterval: 60_000,
  });

  return (
    <div>
      <PageHeader
        title="Security"
        subtitle="IP bans, rate-limit pressure, recent threats — powered by Sentinel."
        actions={
          <a
            href={` + "`${API_URL}/sentinel/ui`" + `}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex shrink-0 items-center gap-2 whitespace-nowrap rounded-lg border border-border bg-bg-elevated px-3 py-2 text-sm font-medium text-foreground hover:bg-bg-hover"
          >
            <ExternalLink className="h-4 w-4" />
            Open Sentinel
          </a>
        }
      />

      <DegradedBanner source="Sentinel" degraded={data?.degraded} />

      {/* KPI row */}
      {isLoading ? (
        <SkeletonCards count={3} />
      ) : (
        <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
          <KPI
            label="Currently Banned IPs"
            value={data?.banned_ips_now ?? 0}
            icon={<Shield className="h-4 w-4" />}
            tone={(data?.banned_ips_now ?? 0) > 0 ? "danger" : "default"}
          />
          <KPI
            label="Auto-bans (last 24h)"
            value={data?.auto_bans_24h ?? 0}
            icon={<AlertCircle className="h-4 w-4" />}
            tone={(data?.auto_bans_24h ?? 0) > 0 ? "warning" : "default"}
          />
          <KPI
            label="Rate-limited IPs (last hour)"
            value={data?.rate_limited_last_hour ?? 0}
            icon={<AlertTriangle className="h-4 w-4" />}
            tone={(data?.rate_limited_last_hour ?? 0) > 0 ? "warning" : "default"}
          />
        </div>
      )}

      {/* Auto-ban escalation policy — explainer card. The schedule is
          surfaced here so operators don't have to grep Sentinel config to
          understand what's about to happen to a re-offender. */}
      <section className="mt-6 rounded-xl border border-accent/20 bg-accent/5 p-5">
        <div className="flex items-start gap-3">
          <Clock className="mt-0.5 h-4 w-4 shrink-0 text-accent" />
          <div className="min-w-0">
            <p className="text-sm font-semibold text-foreground">Escalating auto-ban policy</p>
            <p className="mt-1 text-sm text-text-secondary">
              When an IP trips the brute-force rate limit, Sentinel auto-bans it.
              Re-offenders escalate quickly so a bot can't simply wait out the cooldown.
            </p>
            <ul className="mt-3 grid grid-cols-1 gap-2 text-xs text-text-secondary sm:grid-cols-4">
              <li className="rounded-lg border border-border bg-bg-elevated px-3 py-2">
                <span className="block text-[10px] font-semibold uppercase tracking-wide text-text-muted">1st offence</span>
                <span className="text-foreground font-mono">5 hours</span>
              </li>
              <li className="rounded-lg border border-border bg-bg-elevated px-3 py-2">
                <span className="block text-[10px] font-semibold uppercase tracking-wide text-text-muted">2nd offence</span>
                <span className="text-foreground font-mono">8 hours</span>
              </li>
              <li className="rounded-lg border border-border bg-bg-elevated px-3 py-2">
                <span className="block text-[10px] font-semibold uppercase tracking-wide text-text-muted">3rd offence</span>
                <span className="text-foreground font-mono">24 hours</span>
              </li>
              <li className="rounded-lg border border-border bg-bg-elevated px-3 py-2">
                <span className="block text-[10px] font-semibold uppercase tracking-wide text-text-muted">4th+ offence</span>
                <span className="text-foreground font-mono">7 days</span>
              </li>
            </ul>
          </div>
        </div>
      </section>

      {/* Active bans */}
      <Section title="Active IP bans" icon={<Shield className="h-4 w-4" />}>
        {(data?.active_bans?.length ?? 0) === 0 ? (
          <p className="px-5 py-8 text-center text-sm text-text-muted">No IPs are currently banned.</p>
        ) : (
          <ul className="divide-y divide-border">
            {data!.active_bans!.map((b) => (
              <li key={b.ip} className="flex items-center justify-between gap-3 px-5 py-3">
                <div className="min-w-0">
                  <p className="font-mono text-sm text-foreground">{b.ip}</p>
                  <p className="text-xs text-text-muted">{b.reason}</p>
                </div>
                <div className="text-right text-xs">
                  <span className="rounded bg-danger/10 px-1.5 py-0.5 font-semibold uppercase text-danger">
                    Level {b.level}
                  </span>
                  <p className="mt-1 text-text-muted">expires {new Date(b.expires_at).toLocaleString()}</p>
                </div>
              </li>
            ))}
          </ul>
        )}
      </Section>

      {/* Rate-limit pressure */}
      <Section title="IPs hitting rate limits (last 5 min)" icon={<AlertTriangle className="h-4 w-4" />}>
        {(data?.rate_limit_hits_5min?.length ?? 0) === 0 ? (
          <p className="px-5 py-8 text-center text-sm text-text-muted">No rate-limit caps in the last 5 minutes.</p>
        ) : (
          <ul className="divide-y divide-border">
            {data!.rate_limit_hits_5min!.map((r) => (
              <li key={r.ip} className="flex items-center justify-between gap-3 px-5 py-3">
                <p className="font-mono text-sm text-foreground">{r.ip}</p>
                <span className="text-xs text-text-muted">
                  {r.hits} hits · last {new Date(r.last_hit).toLocaleTimeString()}
                </span>
              </li>
            ))}
          </ul>
        )}
      </Section>

      {/* Recent threats */}
      <Section title="Recent threats" icon={<ActivityIcon className="h-4 w-4" />}>
        {(data?.recent_threats?.length ?? 0) === 0 ? (
          <p className="px-5 py-8 text-center text-sm text-text-muted">No threats detected recently.</p>
        ) : (
          <ul className="divide-y divide-border">
            {data!.recent_threats!.map((t) => (
              <li key={t.id} className="px-5 py-3">
                <div className="flex items-center justify-between">
                  <p className="text-sm font-semibold text-foreground">{t.type}</p>
                  <p className="text-xs text-text-muted">{new Date(t.created_at).toLocaleString()}</p>
                </div>
                <p className="mt-1 text-xs text-text-secondary">{t.description}</p>
                <p className="mt-1 font-mono text-xs text-text-muted">{t.ip}</p>
              </li>
            ))}
          </ul>
        )}
      </Section>
    </div>
  );
}

function KPI({ label, value, icon, tone }: { label: string; value: number; icon: React.ReactNode; tone: "default" | "warning" | "danger" }) {
  const toneClass = {
    default: "border-border bg-bg-elevated",
    warning: "border-warning/30 bg-warning/5",
    danger:  "border-danger/30 bg-danger/5",
  }[tone];
  const iconClass = { default: "text-text-secondary", warning: "text-warning", danger: "text-danger" }[tone];
  return (
    <div className={"rounded-xl border p-4 " + toneClass}>
      <div className="flex items-center justify-between">
        <p className="text-xs font-semibold uppercase tracking-wide text-text-muted">{label}</p>
        <span className={iconClass}>{icon}</span>
      </div>
      <p className="mt-2 text-3xl font-bold text-foreground">{value}</p>
    </div>
  );
}

function Section({ title, icon, children }: { title: string; icon: React.ReactNode; children: React.ReactNode }) {
  return (
    <section className="mt-6 overflow-hidden rounded-xl border border-border bg-bg-elevated">
      <header className="flex items-center gap-2 border-b border-border px-5 py-3">
        <span className="text-text-secondary">{icon}</span>
        <p className="text-xs font-semibold uppercase tracking-wider text-text-muted">{title}</p>
      </header>
      {children}
    </section>
  );
}
`
}

// adminPerformancePageV2 — DGateway Pulse-inspired four-golden-signals
// layout. Latency p50/p95/p99/avg, traffic throughput, error rate +
// active errors, saturation by goroutines/heap/gc/cpu, slowest routes
// table, N+1 query detections, recent errors. "Open full Pulse" deep-link.
func adminPerformancePageV2() string {
	return `"use client";

import { useQuery } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { SkeletonCards } from "@/components/ui/Skeleton";
import { apiClient } from "@/lib/api-client";
import {
  TrendingUp, AlertCircle, AlertTriangle, ExternalLink, Activity as ActivityIcon,
  Cpu, Database, Gauge,
} from "@/lib/icons";

// The Pulse UI is mounted on the Go API, not on this admin host. Use the
// API base so "Open Pulse" works whether the admin is on :3001 in dev or
// a different origin in prod.
import { API_URL } from "@/lib/api-core";

// DegradedBanner says which panels have no data behind them.
//
// The dashboard used to answer a dead upstream with a page of zeros and a 200,
// which reads as good news. The API now names the calls that failed, and this
// puts that on the screen above the numbers it applies to.
function DegradedBanner({ source, degraded }: { source: string; degraded?: string[] }) {
  if (!degraded || degraded.length === 0) return null;
  return (
    <div className="mt-4 flex items-start gap-2 rounded-lg border border-warning/30 bg-warning/5 px-4 py-3 text-sm text-warning">
      <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0" />
      <span>
        {source} did not answer for {degraded.join(", ")}. Those panels are empty rather than zero.
      </span>
    </div>
  );
}

interface PerformanceSummary {
  latency?: { p50: number; p95: number; p99: number; avg: number };
  traffic?: { throughput: number; total: number };
  errors?:  { rate: number; active_open: number };
  saturation?: { goroutines: number; heap_mb: number; gc_cycles: number; cpu_cores: number };
  slowest_routes?: Array<{ route: string; method: string; requests: number; avg: number; p95: number; p99: number; error_rate: number }>;
  n1_detections?: Array<{ route: string; query_count: number; first_seen: string }>;
  recent_errors?: Array<{ id: string; route: string; message: string; created_at: string }>;
  // The calls to Pulse that did not answer. Empty is the healthy case.
  degraded?: string[];
}

export default function PerformancePage() {
  const { data, isLoading } = useQuery<PerformanceSummary>({
    queryKey: ["performance", "summary"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<PerformanceSummary>("/api/admin/observability/summary");
        return data;
      } catch {
        // Nothing came back. Every signal would render as a dash, which looks
        // like an idle app rather than a dashboard that cannot see one.
        return { degraded: ["the whole summary"] };
      }
    },
    refetchInterval: 30_000,
  });

  return (
    <div>
      <PageHeader
        title="Performance"
        subtitle="Four SRE golden signals + route, query, and error detail — powered by Pulse."
        actions={
          <a
            href={` + "`${API_URL}/pulse/ui`" + `}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex shrink-0 items-center gap-2 whitespace-nowrap rounded-lg border border-border bg-bg-elevated px-3 py-2 text-sm font-medium text-foreground hover:bg-bg-hover"
          >
            <ExternalLink className="h-4 w-4" />
            Open Pulse
          </a>
        }
      />

      <DegradedBanner source="Pulse" degraded={data?.degraded} />

      {isLoading ? <SkeletonCards count={4} /> : (
        <>
          {/* Latency */}
          <SignalGroup
            title="Latency"
            tagline="How long requests take. Watch p95/p99 — the average hides the long tail."
          >
            <Signal label="P50" value={fmt(data?.latency?.p50, "ms")} icon={<Gauge className="h-4 w-4" />} />
            <Signal label="P95" value={fmt(data?.latency?.p95, "ms")} icon={<Gauge className="h-4 w-4" />} />
            <Signal label="P99" value={fmt(data?.latency?.p99, "ms")} icon={<Gauge className="h-4 w-4" />} tone="warning" />
            <Signal label="AVG" value={fmt(data?.latency?.avg, "ms")} icon={<Gauge className="h-4 w-4" />} />
          </SignalGroup>

          {/* Traffic */}
          <SignalGroup
            title="Traffic"
            tagline="How much demand the API is handling right now."
          >
            <Signal label="Throughput" value={fmtRate(data?.traffic?.throughput)} icon={<TrendingUp className="h-4 w-4 text-success" />} />
            <Signal label="Total requests" value={fmt(data?.traffic?.total)} icon={<TrendingUp className="h-4 w-4 text-success" />} />
          </SignalGroup>

          {/* Errors */}
          <SignalGroup
            title="Errors"
            tagline="Rate of failures. Spikes here usually correlate with latency spikes — check both."
          >
            <Signal label="Error rate" value={fmtPct(data?.errors?.rate)} icon={<AlertCircle className="h-4 w-4 text-danger" />} tone={(data?.errors?.rate ?? 0) > 0 ? "danger" : "default"} />
            <Signal label="Active errors (open)" value={fmt(data?.errors?.active_open)} icon={<AlertCircle className="h-4 w-4 text-danger" />} tone={(data?.errors?.active_open ?? 0) > 0 ? "danger" : "default"} />
          </SignalGroup>

          {/* Saturation */}
          <SignalGroup
            title="Saturation"
            tagline="How full your resources are. Red bands here mean a bottleneck is imminent or already firing — fix before users feel it."
          >
            <Signal label="Goroutines" value={fmt(data?.saturation?.goroutines)} icon={<ActivityIcon className="h-4 w-4" />} />
            <Signal label="Heap alloc" value={fmt(data?.saturation?.heap_mb, "MB")} icon={<Database className="h-4 w-4" />} />
            <Signal label="GC cycles" value={fmt(data?.saturation?.gc_cycles)} icon={<Database className="h-4 w-4" />} />
            <Signal label="CPU cores" value={fmt(data?.saturation?.cpu_cores)} icon={<Cpu className="h-4 w-4" />} />
          </SignalGroup>
        </>
      )}

      {/* Slowest routes */}
      <section className="mt-6 overflow-hidden rounded-xl border border-border bg-bg-elevated">
        <header className="border-b border-border px-5 py-3">
          <p className="text-xs font-semibold uppercase tracking-wider text-text-muted">Slowest routes</p>
        </header>
        {(data?.slowest_routes?.length ?? 0) === 0 ? (
          <p className="px-5 py-8 text-center text-sm text-text-muted">No route latency data yet.</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full table-fixed">
              <thead>
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-semibold uppercase text-text-muted">Route</th>
                  <th className="w-20 px-2 py-3 text-right text-xs font-semibold uppercase text-text-muted">Reqs</th>
                  <th className="w-20 px-2 py-3 text-right text-xs font-semibold uppercase text-text-muted">Avg</th>
                  <th className="w-20 px-2 py-3 text-right text-xs font-semibold uppercase text-text-muted">P95</th>
                  <th className="w-20 px-2 py-3 text-right text-xs font-semibold uppercase text-text-muted">P99</th>
                  <th className="w-24 px-2 py-3 text-right text-xs font-semibold uppercase text-text-muted">Err rate</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {data!.slowest_routes!.map((r, i) => (
                  <tr key={i}>
                    <td className="truncate px-4 py-3 font-mono text-xs text-foreground">{r.method} {r.route}</td>
                    <td className="px-2 py-3 text-right text-xs text-foreground">{r.requests}</td>
                    <td className="px-2 py-3 text-right text-xs text-foreground">{fmt(r.avg, "ms")}</td>
                    <td className="px-2 py-3 text-right text-xs text-foreground">{fmt(r.p95, "ms")}</td>
                    <td className="px-2 py-3 text-right text-xs text-foreground">{fmt(r.p99, "ms")}</td>
                    <td className="px-2 py-3 text-right text-xs text-foreground">{fmtPct(r.error_rate)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {/* N+1 detections */}
      <section className="mt-6 overflow-hidden rounded-xl border border-border bg-bg-elevated">
        <header className="border-b border-border px-5 py-3">
          <p className="text-xs font-semibold uppercase tracking-wider text-text-muted">N+1 query detections</p>
        </header>
        {(data?.n1_detections?.length ?? 0) === 0 ? (
          <p className="px-5 py-8 text-center text-sm text-text-muted">No N+1 queries detected.</p>
        ) : (
          <ul className="divide-y divide-border">
            {data!.n1_detections!.map((n, i) => (
              <li key={i} className="flex items-center justify-between gap-3 px-5 py-3">
                <p className="truncate font-mono text-xs text-foreground">{n.route}</p>
                <span className="shrink-0 text-xs text-text-muted">
                  {n.query_count} queries · first seen {new Date(n.first_seen).toLocaleString()}
                </span>
              </li>
            ))}
          </ul>
        )}
      </section>

      {/* Recent errors */}
      <section className="mt-6 overflow-hidden rounded-xl border border-border bg-bg-elevated">
        <header className="border-b border-border px-5 py-3">
          <p className="text-xs font-semibold uppercase tracking-wider text-text-muted">Recent errors</p>
        </header>
        {(data?.recent_errors?.length ?? 0) === 0 ? (
          <p className="px-5 py-8 text-center text-sm text-text-muted">No errors recorded.</p>
        ) : (
          <ul className="divide-y divide-border">
            {data!.recent_errors!.map((e) => (
              <li key={e.id} className="px-5 py-3">
                <div className="flex items-center justify-between">
                  <p className="truncate font-mono text-xs text-foreground">{e.route}</p>
                  <p className="shrink-0 text-xs text-text-muted">{new Date(e.created_at).toLocaleString()}</p>
                </div>
                <p className="mt-1 text-sm text-danger">{e.message}</p>
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}

interface SignalGroupProps { title: string; tagline: string; children: React.ReactNode }
function SignalGroup({ title, tagline, children }: SignalGroupProps) {
  return (
    <section className="mt-6">
      <header className="mb-3">
        <p className="text-sm font-semibold uppercase tracking-wider text-text-muted">{title}</p>
        <p className="text-xs text-text-secondary">{tagline}</p>
      </header>
      <div className="grid grid-cols-2 gap-3 md:grid-cols-4">{children}</div>
    </section>
  );
}

interface SignalProps { label: string; value: string; icon: React.ReactNode; tone?: "default" | "warning" | "danger" }
function Signal({ label, value, icon, tone = "default" }: SignalProps) {
  const toneClass = { default: "border-border bg-bg-elevated", warning: "border-warning/30 bg-warning/5", danger: "border-danger/30 bg-danger/5" }[tone];
  return (
    <div className={"rounded-xl border p-4 " + toneClass}>
      <div className="flex items-center justify-between">
        <p className="text-xs font-semibold uppercase tracking-wide text-text-muted">{label}</p>
        <span className="text-text-secondary">{icon}</span>
      </div>
      <p className="mt-2 text-2xl font-bold text-foreground">{value}</p>
    </div>
  );
}

function fmt(n: number | undefined, suffix: string = ""): string {
  if (n === undefined || n === null || Number.isNaN(n)) return "—";
  return Math.round(n).toLocaleString() + suffix;
}

// fmtRate formats a per-second rate. Rounding to an integer would collapse a
// low-but-real throughput (e.g. 0.14 req/s) to "0/s", so keep two decimals
// below 10 and switch to whole numbers once the rate is large.
function fmtRate(n: number | undefined): string {
  if (n === undefined || n === null || Number.isNaN(n)) return "—";
  if (n > 0 && n < 10) return n.toFixed(2) + "/s";
  return Math.round(n).toLocaleString() + "/s";
}

// Both callers pass a Pulse error rate, and Pulse already reports those as a
// percentage (errCount / total * 100 — see storage_memory.go). Multiplying by
// 100 here scaled every value by 100 again: 3 errors in 38 requests rendered as
// "789.47%". Takes a percentage, prints a percentage.
function fmtPct(n: number | undefined): string {
  if (n === undefined || n === null || Number.isNaN(n)) return "—";
  return n.toFixed(2) + "%";
}
`
}
