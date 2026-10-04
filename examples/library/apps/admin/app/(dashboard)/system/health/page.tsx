"use client";

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
