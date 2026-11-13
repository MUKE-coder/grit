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
import { HubTile } from "@/components/ui/hub-tile";
import { UnderlineTabs } from "@/components/ui/tabs-underline";
import { EmptyState } from "@/components/ui/empty-state";
import { type PillTone } from "@/components/ui/pill";
import { inputClasses } from "@/components/ui/input";
import { useSystemHealth, stateOf } from "@/hooks/use-system-health";
import {
  Activity, Bell, Calendar, Database, FileText, Mail,
  MessageSquare, Shield, ShieldCheck, TrendingUp, Upload, Link as LinkIcon, Search,
  UserCheck, Settings, LayoutGrid, KeyRound, Workflow, Server, Trash2,
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
  { href: "/system/deleted-accounts", category: "Security & Access", title: "Deleted accounts",  description: "Accounts somebody closed. Restore one, or remove it for good.", icon: <UserCheck className="h-5 w-5" /> },
  { href: "/system/trash",         category: "Data & Files", title: "Trash",                description: "Deleted records, restorable for 30 days before they go for good.", icon: <Trash2 className="h-5 w-5" /> },
  { href: "/system/performance",   category: "Operations", title: "Operations",      description: "Latency, traffic, errors, saturation, slow routes and scaling readiness.", icon: <TrendingUp className="h-5 w-5" /> },
  { href: "/system/jobs",          category: "Operations", title: "Background Jobs",  description: "Queue depth, in-flight workers, dead-letter queue.",                      icon: <Database className="h-5 w-5" /> , module: "jobs" },
  { href: "/system/sagas",         category: "Operations", title: "Sagas",           description: "Multi-step processes, and the ones whose undo could not be completed.", icon: <Workflow className="h-5 w-5" /> },
  { href: "/system/routes",        category: "Operations", title: "Routes",          description: "Every endpoint this API serves, searchable, with a curl for each.",     icon: <Server className="h-5 w-5" /> },
  { href: "/system/cron",          category: "Operations", title: "Cron Schedules",  description: "Recurring jobs, next-run times, run history.",                            icon: <Calendar className="h-5 w-5" /> , module: "cron" },
  { href: "/system/activity",      category: "Operations", title: "User Activity",   description: "Auth events, writes, operator actions with IP + severity.",               icon: <Activity className="h-5 w-5" /> , module: "audit" },
  // ── Security & Access ───────────────────────────────────────────────────
  { href: "/account",               category: "Security & Access", title: "Account",             description: "Your password, two-factor, passkeys and the devices you are signed in on.",  icon: <UserCheck className="h-5 w-5" /> },
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
  const [filter, setFilter] = useState("");

  const health = useSystemHealth();

  // Everything on the hub, as one list, so the filter can search across tabs.
  //
  // A filter that only searches the tab you are on is a filter you have to
  // use five times, and the whole reason somebody types into this box is that
  // they do not know which tab the thing is under.
  const all: Array<{
    href: string;
    icon: React.ReactNode;
    title: string;
    description: string;
    category: string;
  }> = [
    ...tiles.map((t) => ({
      href: t.href,
      icon: t.icon,
      title: t.title,
      description: t.description,
      category: t.category as string,
    })),
    ...extensions.map((n) => ({
      href: n.href,
      icon: INTERNAL_ICON[n.iconKey] ?? <LayoutGrid className="h-5 w-5" />,
      title: n.label,
      description: "Added by a Grit plugin.",
      category: EXTENSIONS,
    })),
  ];

  const query = filter.trim().toLowerCase();
  const matches = query
    ? all.filter((t) =>
        (t.title + " " + t.description + " " + t.category).toLowerCase().includes(query)
      )
    : all.filter((t) => t.category === current);

  const countFor = (tab: string) => all.filter((t) => t.category === tab).length;

  // The status line each tile shows, from facts the panel already has.
  //
  // A grid of tiles that all say the same static sentence is a menu wearing a
  // dashboard's clothes. These say what is true now, and the one that matters
  // most is that a deployment without Redis has two components off on purpose
  // rather than two things broken.
  const statusFor = (href: string): { text: string; tone: PillTone } | undefined => {
    if (href === "/system/health") {
      if (health.summary.degraded > 0) {
        return { text: health.summary.degraded + " degraded", tone: "bad" };
      }
      if (health.summary.off > 0) {
        return {
          text: health.summary.healthy + " healthy · " + health.summary.off + " off",
          tone: "neutral",
        };
      }
      if (health.summary.healthy > 0) {
        return { text: health.summary.healthy + " healthy", tone: "good" };
      }
      return undefined;
    }
    if (href === "/system/jobs" || href === "/system/cron") {
      const redis = stateOf(health.data?.redis);
      if (redis === "off") return { text: "Off · needs Redis", tone: "neutral" };
      if (redis === "degraded") return { text: "Redis degraded", tone: "bad" };
      return undefined;
    }
    return undefined;
  };

  return (
    <div>
      <PageHeader
        title="System"
        subtitle="Every operational surface for this app, grouped. Pick a tab, then a tile."
        actions={
          <div className="relative w-full sm:w-72">
            <Search
              className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-text-muted"
              aria-hidden="true"
            />
            <input
              type="search"
              value={filter}
              onChange={(event) => setFilter(event.target.value)}
              placeholder="Filter every tab…"
              aria-label="Filter system pages"
              className={inputClasses({ className: "pl-9 pr-3" })}
            />
          </div>
        }
      />

      {/* Hidden while filtering: the results span every tab, so highlighting
          one of them would be a lie about where they came from. */}
      {!query && (
        <UnderlineTabs
          className="mb-6"
          label="System areas"
          value={current}
          onChange={setActive}
          tabs={tabs.map((t) => ({ value: t, label: t, count: countFor(t) }))}
        />
      )}

      {query && (
        <p className="mb-4 text-sm text-text-secondary">
          {matches.length === 0
            ? "Nothing matches"
            : matches.length + (matches.length === 1 ? " page" : " pages") + " across every tab"}
        </p>
      )}

      {matches.length === 0 ? (
        <div className="rounded-xl border border-border bg-bg-elevated">
          <EmptyState
            icon={<LayoutGrid className="h-5 w-5" />}
            title={"Nothing matches “" + filter.trim() + "”"}
            description="Every system page is searched, not only the tab you were on. Try a shorter word."
          />
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
          {matches.map((t) => {
            const status = statusFor(t.href);
            return (
              <HubTile
                key={t.href}
                href={t.href}
                icon={t.icon}
                title={t.title}
                description={t.description}
                status={status?.text}
                statusTone={status?.tone}
              />
            );
          })}
        </div>
      )}

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
  summarise,
  stateOf,
  type HealthState,
  type HealthResponse,
} from "@/hooks/use-system-health";
import { Note } from "@/components/ui/note";
import { Pill } from "@/components/ui/pill";
import {
  CheckCircle, AlertCircle, AlertTriangle, RefreshCw, Database, Mail, Server,
  Activity as ActivityIcon, HardDrive, Zap, FolderOpen, Webhook,
} from "@/lib/icons";

// The shapes and the state reading live with the shared query, in
// hooks/use-system-health. They were declared here as well, which is two
// copies of a type that describes one endpoint's response: the pill in the
// top bar and this page would have disagreed about what "off" means the first
// time one of them was edited.

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
  const summary = summarise(data, isError);
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
        // One action. The status pill that used to sit beside it said the
        // same thing as the pill in the top bar and the same thing as the
        // overall card below, which is the same sentence three times on one
        // screen.
        actions={
          <IconButton
            icon={<RefreshCw className={"h-4 w-4 " + (isFetching ? "animate-spin" : "")} />}
            label="Run health check"
            onClick={() => { refetch(); queryClient.invalidateQueries({ queryKey: ["system-health"] }); }}
          />
        }
      />

      {isError && (
        <div className="mb-4 flex items-start gap-2 rounded-lg border border-danger-line bg-danger-soft px-4 py-3 text-sm text-danger">
          <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0" />
          <span>
            The API did not answer /api/health. Nothing can be reported about the
            components behind it, which is itself the finding.
          </span>
        </div>
      )}

      {/* The one sentence, and the proportions behind it.
          The page used to open on a grid of eight cards, which is eight
          things to read before knowing whether anything is wrong. The bar is
          the same counts as a shape, so the answer arrives before any of the
          words do. */}
      {!isLoading && (
        <div className="mb-6 rounded-xl border border-border bg-bg-elevated p-5">
          <div className="flex flex-wrap items-center justify-between gap-4">
            <div className="flex items-center gap-3">
              <span
                className={
                  "inline-flex h-9 w-9 shrink-0 items-center justify-center rounded-lg " +
                  (overall === "ok"
                    ? "bg-success-soft text-success"
                    : overall === "degraded"
                      ? "bg-danger-soft text-danger"
                      : "bg-bg-hover text-text-muted")
                }
              >
                {overall === "ok" ? (
                  <CheckCircle className="h-5 w-5" />
                ) : (
                  <AlertTriangle className="h-5 w-5" />
                )}
              </span>
              <div>
                <p className="text-[15px] font-semibold text-foreground">{summary.label}</p>
                <p className="text-xs text-text-muted">
                  Re-checked every 30 seconds
                </p>
              </div>
            </div>

            <div className="flex items-center gap-4 text-sm">
              <span className="flex items-center gap-1.5 text-text-secondary">
                <span className="h-2 w-2 rounded-full bg-success" aria-hidden="true" />
                <span className="font-semibold text-foreground tabular-nums">{summary.healthy}</span>
                healthy
              </span>
              <span className="flex items-center gap-1.5 text-text-secondary">
                <span className="h-2 w-2 rounded-full bg-danger" aria-hidden="true" />
                <span className="font-semibold text-foreground tabular-nums">{summary.degraded}</span>
                degraded
              </span>
              <span className="flex items-center gap-1.5 text-text-secondary">
                <span className="h-2 w-2 rounded-full bg-border-strong" aria-hidden="true" />
                <span className="font-semibold text-foreground tabular-nums">
                  {summary.off + summary.unknown}
                </span>
                off
              </span>
            </div>
          </div>

          {/* The same counts, as proportions. */}
          <div
            className="mt-4 flex h-1.5 overflow-hidden rounded-full bg-bg-hover"
            role="img"
            aria-label={
              summary.healthy + " healthy, " + summary.degraded + " degraded, " +
              (summary.off + summary.unknown) + " off"
            }
          >
            {summary.healthy > 0 && (
              <span className="bg-success" style={{ flexGrow: summary.healthy }} />
            )}
            {summary.degraded > 0 && (
              <span className="bg-danger" style={{ flexGrow: summary.degraded }} />
            )}
            {summary.off + summary.unknown > 0 && (
              <span
                className="bg-border-strong"
                style={{ flexGrow: summary.off + summary.unknown }}
              />
            )}
          </div>
        </div>
      )}

      {isLoading ? (
        <SkeletonCards count={8} />
      ) : (
        <>
          {/* Running first, then the ones nobody configured.
              A single grid put Redis, which this deployment runs without on
              purpose, between two components that are working, and the only
              thing distinguishing it was a word in small type. Grouped, the
              bottom section is read once and then ignored, which is what it
              deserves. */}
          <HealthGroup
            title="Running"
            cards={cards.filter((c) => c.state === "ok" || c.state === "degraded")}
          />
          <HealthGroup
            title="Not configured"
            cards={cards.filter((c) => c.state === "off" || c.state === "unknown")}
          />
        </>
      )}

      <Note className="mt-6 max-w-2xl">
        <strong className="font-medium text-foreground">Off</strong> means this
        deployment never configured the component.{" "}
        <strong className="font-medium text-foreground">Unknown</strong> means a
        probe could not measure it. Neither counts against the status above: only{" "}
        <strong className="font-medium text-foreground">degraded</strong> does.
      </Note>
    </div>
  );
}

// One group of component cards, with its count in the heading.
//
// Nothing at all is not an empty state here: a deployment with every
// component configured should simply not show a "Not configured" heading.
function HealthGroup({ title, cards }: { title: string; cards: Card[] }) {
  if (cards.length === 0) return null;
  return (
    <section className="mb-6">
      <h2 className="mb-3 font-mono text-[11px] uppercase tracking-wider text-text-muted">
        {title} · {cards.length}
      </h2>
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
        {cards.map((c) => (
          <HealthCard key={c.key} card={c} />
        ))}
      </div>
    </section>
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

// A state is a word, in the same pill the rest of the panel uses.
//
// It was a tick for healthy, a warning triangle for degraded, and small
// uppercase text for the other two: three different objects for one
// dimension, and two of them colour-only, which is nothing at all to a reader
// who cannot separate the green from the red.
function StateBadge({ state }: { state: HealthState }) {
  if (state === "ok") return <Pill tone="good" dot>Healthy</Pill>;
  if (state === "degraded") return <Pill tone="bad" dot>Degraded</Pill>;
  if (state === "off") return <Pill tone="neutral" dot>Off</Pill>;
  return <Pill tone="neutral" dot>Unknown</Pill>;
}

// What to do about a component nobody configured.
//
// "Off" answers what, and leaves why and how. These two are off on every
// machine that runs Grit without Redis, which is a supported way to run, so
// the card says what turning it on would take rather than leaving somebody to
// search the docs for it.
const TURN_ON: Record<string, string> = {
  redis: "Set REDIS_URL in .env to turn this on",
  jobs: "Turns on when Redis is configured",
};

function HealthCard({ card }: { card: Card }) {
  const tone = TONES[card.state];
  const turnOn = TURN_ON[card.key];

  return (
    <div className={"flex flex-col rounded-xl border p-4 " + tone.card}>
      <div className="mb-3 flex items-center justify-between gap-2">
        <span
          className={"inline-flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-bg-elevated " + tone.icon}
          aria-hidden="true"
        >
          {card.icon}
        </span>
        <StateBadge state={card.state} />
      </div>
      <p className="text-sm font-semibold text-foreground">{card.label}</p>
      <p className="mt-1 flex-1 text-xs text-text-secondary">{card.detail}</p>
      {card.meta && (
        <p className="mt-2 inline-flex w-fit items-center gap-1 rounded bg-bg-elevated px-1.5 py-0.5 font-mono text-[10px] text-text-muted">
          {card.meta}
        </p>
      )}
      {/* Dashed, because it is not a state: it is the next step.
          A solid well here reads as another fact about the component. */}
      {turnOn && (
        <p className="mt-3 rounded-md border border-dashed border-border-strong px-3 py-2 font-mono text-[11.5px] text-text-secondary">
          {turnOn}
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
func adminSecurityPageV2() string { return tmpl("admin/pages/security.tsx") }

// adminPerformancePageV2 — DGateway Pulse-inspired four-golden-signals
// layout. Latency p50/p95/p99/avg, traffic throughput, error rate +
// active errors, saturation by goroutines/heap/gc/cpu, slowest routes
// table, N+1 query detections, recent errors. "Open full Pulse" deep-link.
// adminTrashPage emits app/(dashboard)/system/trash/page.tsx: what a soft
// delete left behind, and the two buttons that resolve it.
func adminTrashPage() string { return tmpl("admin/app/system/trash-page.tsx") }

// adminDeletedAccountsPage emits app/(dashboard)/system/deleted-accounts/page.tsx.
// Separate from the bin because users are not resources: they are kept out of
// the sync registry the bin reads, deliberately.
func adminDeletedAccountsPage() string { return tmpl("admin/app/system/deleted-accounts-page.tsx") }

// adminPerformancePageV2 emits app/(dashboard)/system/performance/page.tsx:
// the operations page, which /system/observability now redirects to.
func adminPerformancePageV2() string { return tmpl("admin/app/system/operations-page.tsx") }
