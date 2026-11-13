package scaffold

// adminUseSystemHealth is hooks/use-system-health.ts: one query, read by
// everything that reports how the system is doing.
//
// The same facts appear in five places: the pill in the top bar, the badge on
// the sidebar's System health item, the dashboard's quick-access sub-line, the
// tiles on the system hub, and the health page itself. Five components each
// fetching for themselves is five requests on every page load and, worse, five
// answers that can differ: the pill says all clear while the tile two inches
// below it says two components are off, because one of them asked a second
// earlier.
//
// React Query dedupes by key, so every caller of this hook shares one request
// and one cache entry. The facts are then the same everywhere by construction
// rather than by everyone remembering to use the same interval.
func adminUseSystemHealth() string {
	return `"use client";

import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

/**
 * The four states /api/health answers in.
 *
 * "off" is a component this deployment never configured, which is not a
 * problem: a project running without Redis is a supported way to run.
 * "unknown" is a probe that could not find out, which is not an accusation
 * either. Only "degraded" is bad news, and only it moves the overall status.
 */
export type HealthState = "ok" | "degraded" | "off" | "unknown";

export interface HealthComponent {
  state?: HealthState;
  /**
   * The state as a boolean. Kept because an API that has not been upgraded
   * yet sends only this, and because a load balancer reads it.
   */
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

export interface HealthResponse {
  status: HealthState;
  database?: HealthComponent;
  redis?: HealthComponent;
  api?: HealthComponent;
  jobs?: HealthComponent;
  email?: HealthComponent;
  storage?: HealthComponent;
  events?: HealthComponent;
  /**
   * This replica's hub. With several replicas each reports its own share, so
   * the numbers are per process, not per deployment.
   */
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

/** The components the health page lays out, in the order it lays them out. */
export const HEALTH_COMPONENTS = [
  "database",
  "api",
  "storage",
  "email",
  "events",
  "realtime",
  "redis",
  "jobs",
] as const;

/**
 * Reads a component's state, with a fallback for an API older than it.
 *
 * A component that is not in the response was not reported at all, so
 * unknown. An older API sends only ok, where false covered both "down" and
 * "never configured": unknown is the honest reading, and reading it as down
 * is what put a red card on every machine that runs without Redis.
 */
export function stateOf(component?: HealthComponent): HealthState {
  if (!component) return "unknown";
  if (component.state) return component.state;
  return component.ok ? "ok" : "unknown";
}

export interface HealthSummary {
  /** The one word for the whole system. */
  overall: HealthState;
  healthy: number;
  degraded: number;
  off: number;
  unknown: number;
  /** "All systems operational", "2 components degraded". */
  label: string;
  /** Short enough for a sidebar badge: "2 off", "1 degraded", or "". */
  badge: string;
  /** What the badge is reporting, for its colour. */
  badgeTone: "neutral" | "warn" | "bad";
}

/**
 * Counts the components by state.
 *
 * Only degraded moves the overall status. A deployment with no Redis is a
 * supported deployment, and reporting it as a problem trains people to ignore
 * the one place that tells them about real problems.
 */
export function summarise(data: HealthResponse | undefined, failed: boolean): HealthSummary {
  if (failed) {
    return {
      overall: "degraded",
      healthy: 0,
      degraded: 0,
      off: 0,
      unknown: 0,
      label: "Cannot reach the API",
      badge: "down",
      badgeTone: "bad",
    };
  }

  let healthy = 0;
  let degraded = 0;
  let off = 0;
  let unknown = 0;

  if (data) {
    for (const key of HEALTH_COMPONENTS) {
      // realtime is reported as a hub rather than a component; it is healthy
      // whenever the API answered at all, because it lives in the API.
      const state =
        key === "realtime"
          ? data.realtime
            ? "ok"
            : "unknown"
          : stateOf(data[key as keyof HealthResponse] as HealthComponent | undefined);
      if (state === "ok") healthy += 1;
      else if (state === "degraded") degraded += 1;
      else if (state === "off") off += 1;
      else unknown += 1;
    }
  }

  const overall: HealthState = data?.status ?? (degraded > 0 ? "degraded" : "unknown");

  let label = "Checking…";
  if (data) {
    label =
      degraded > 0
        ? degraded + (degraded === 1 ? " component degraded" : " components degraded")
        : "All systems operational";
  }

  let badge = "";
  let badgeTone: HealthSummary["badgeTone"] = "neutral";
  if (degraded > 0) {
    badge = degraded + " degraded";
    badgeTone = "bad";
  } else if (off > 0) {
    badge = off + " off";
    badgeTone = "neutral";
  }

  return { overall, healthy, degraded, off, unknown, label, badge, badgeTone };
}

/**
 * The shared health query.
 *
 * No catch around the request. An earlier version answered a failed call with
 * a benign { status: "ok" } object, so an API that was down painted "All
 * systems operational", which is the one thing a status page must never do.
 */
export function useSystemHealth() {
  const query = useQuery<HealthResponse>({
    queryKey: ["system-health"],
    queryFn: async () => {
      const { data } = await apiClient.get<HealthResponse>("/api/health");
      return data;
    },
    refetchInterval: 30_000,
    staleTime: 25_000,
  });

  return { ...query, summary: summarise(query.data, query.isError) };
}
`
}
