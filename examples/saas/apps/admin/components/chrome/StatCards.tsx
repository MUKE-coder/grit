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
  /**
   * The rows this card counted. Give it one and the card becomes a button
   * that narrows the table to them.
   *
   * Only the cards that stand for a subset get it. "Total" and "This Week"
   * have nothing to narrow to and stay as they were.
   */
  filter?: { field: string; value: string };
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

function StatCardItem({
  stat,
  onFilter,
  active,
}: {
  stat: StatCard;
  onFilter?: (filter: StatFilter) => void;
  active?: boolean;
}) {
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

  const body = (
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
  );

  const base = "rounded-xl border p-5 text-left transition-colors ";

  // A card with nothing to narrow to stays a div, and loses the hover state
  // with it: a border that lights up under the pointer is a promise, and this
  // one had nothing behind it.
  if (!stat.filter || !onFilter) {
    return <div className={base + "border-border bg-bg-secondary"}>{body}</div>;
  }

  return (
    <button
      type="button"
      onClick={() => onFilter(stat.filter!)}
      aria-pressed={active}
      className={
        base +
        "w-full cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-accent/50 " +
        (active
          ? "border-accent bg-accent/10"
          : "border-border bg-bg-secondary hover:border-accent/50 hover:bg-bg-hover")
      }
    >
      {body}
      <span className="sr-only">
        {active ? ". Showing only these. Activate to show all." : ". Activate to show only these."}
      </span>
    </button>
  );
}

export interface StatFilter {
  field: string;
  value: string;
}

/** The row of stat cards under a page header. */
export function StatCards({
  stats,
  onFilter,
  isFilterActive,
}: {
  stats: StatCard[];
  /** Given, every card carrying a filter becomes a button. */
  onFilter?: (filter: StatFilter) => void;
  /** Whether the table is currently narrowed to that card. */
  isFilterActive?: (filter: StatFilter) => boolean;
}) {
  return (
    <div className="mb-8 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {stats.map((stat, i) => (
        <StatCardItem
          key={i}
          stat={stat}
          onFilter={onFilter}
          active={stat.filter ? (isFilterActive?.(stat.filter) ?? false) : false}
        />
      ))}
    </div>
  );
}
