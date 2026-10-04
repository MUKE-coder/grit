"use client";

import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { resourceKeys } from "@/hooks/use-resource";
import { getIcon, TrendingUp, TrendingDown } from "@/lib/icons";

export interface StatCard {
  label: string;
  icon?: string;
  color?: "default" | "success" | "warning" | "danger" | "info";
  /** Either provide a static value... */
  value?: string | number;
  /** ...or an endpoint + field to fetch it from the API (e.g. endpoint: "/api/posts?page_size=1", field: "meta.total") */
  endpoint?: string;
  field?: string;
  /** Optional trend delta shown next to value */
  trend?: { value: number; direction: "up" | "down" };
  /** Shows the placeholder while a static value is still on its way. */
  loading?: boolean;
}

const colorClasses: Record<string, { bg: string; text: string }> = {
  default: { bg: "bg-accent/10", text: "text-accent" },
  success: { bg: "bg-success/10", text: "text-success" },
  warning: { bg: "bg-warning/10", text: "text-warning" },
  danger: { bg: "bg-danger/10", text: "text-danger" },
  info: { bg: "bg-info/10", text: "text-info" },
};

// Reads a dotted path from an object. e.g. getPath(data, "meta.total") → data.meta.total
function getPath(obj: unknown, path: string): unknown {
  return path.split(".").reduce<unknown>(
    (acc, key) => (acc == null ? acc : (acc as Record<string, unknown>)[key]),
    obj,
  );
}

function StatCardItem({ stat }: { stat: StatCard }) {
  const color = colorClasses[stat.color || "default"];
  const Icon = stat.icon ? getIcon(stat.icon) : null;

  // The key starts with the base endpoint (no query string) so a save to the
  // resource, which invalidates [endpoint], reaches this card too.
  const baseEndpoint = stat.endpoint ? stat.endpoint.split("?")[0] : "stat";

  const { data, isLoading } = useQuery({
    queryKey: [...resourceKeys.stats(baseEndpoint), stat.endpoint, stat.field],
    queryFn: async () => {
      if (!stat.endpoint) return null;
      const res = await apiClient.get(stat.endpoint);
      return res.data;
    },
    enabled: !!stat.endpoint,
  });

  const value =
    stat.value !== undefined
      ? stat.value
      : stat.endpoint && stat.field
      ? (getPath(data, stat.field) as string | number | undefined) ?? "—"
      : "—";

  return (
    <div className="rounded-xl border border-border bg-bg-secondary p-5 transition-colors hover:border-border/80">
      <div className="flex items-start justify-between">
        <div className="min-w-0 flex-1">
          <p className="text-[11px] font-semibold uppercase tracking-wider text-text-muted">
            {stat.label}
          </p>
          <div className="mt-2 flex items-baseline gap-2">
            {(isLoading && stat.endpoint) || stat.loading ? (
              <div className="h-7 w-16 rounded bg-bg-hover animate-pulse" />
            ) : (
              <p className="text-2xl font-bold text-foreground tabular-nums">
                {typeof value === "number" ? value.toLocaleString() : value}
              </p>
            )}
            {stat.trend && (
              <span
                className={`flex items-center gap-0.5 text-xs font-medium ${
                  stat.trend.direction === "up" ? "text-success" : "text-danger"
                }`}
              >
                {stat.trend.direction === "up" ? (
                  <TrendingUp className="h-3 w-3" />
                ) : (
                  <TrendingDown className="h-3 w-3" />
                )}
                {stat.trend.value}%
              </span>
            )}
          </div>
        </div>
        {Icon && (
          <div className={`flex h-9 w-9 items-center justify-center rounded-lg ${color.bg}`}>
            <Icon className={`h-4 w-4 ${color.text}`} />
          </div>
        )}
      </div>
    </div>
  );
}

/** The row of stat cards under a page header. */
export function StatCards({ stats }: { stats: StatCard[] }) {
  return (
    <div className="mb-8 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {stats.map((stat, i) => (
        <StatCardItem key={i} stat={stat} />
      ))}
    </div>
  );
}
