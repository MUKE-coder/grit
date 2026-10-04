"use client";

import Link from "next/link";
import dynamic from "next/dynamic";
import { Suspense } from "react";
import { useResourceDashboardStats, type ResourceDashboardStats } from "@/hooks/use-resource";
import { dateRangeToQueryParams, type DateRange } from "@/components/tables/date-filter";
import { getIcon, ArrowUpRight } from "@/lib/icons";
import type { ResourceDefinition } from "@/lib/resource";

// The sparkline is recharts, about 400 KB, so it loads after the card renders
// instead of riding in the first load of every page that shows a card.
const ResourceSparkline = dynamic(
  () => import("@/components/dashboard/ResourceSparkline").then((m) => m.ResourceSparkline),
  { ssr: false },
);

interface Props {
  resource: ResourceDefinition;
  dateRange: DateRange;
}

// The card's part of the response the latest table shares. The sparkline is
// always 30 daily buckets so its shape stays stable; the date range only
// changes the total.
function pickTotals(stats: ResourceDashboardStats) {
  return { total: stats.total, series: stats.series };
}

export function ResourceStatCard({ resource, dateRange }: Props) {
  const query = useResourceDashboardStats(resource.endpoint, dateRangeToQueryParams(dateRange), pickTotals);

  const Icon = getIcon(resource.icon);
  const label = resource.label?.plural ?? resource.slug;
  const total = query.data?.total ?? 0;
  const series = query.data?.series ?? [];
  // Treat empty/error as zero so the layout doesn't shift.
  const sparkData = series.length
    ? series
    : Array.from({ length: 30 }).map((_, i) => ({
        date: String(i),
        count: 0,
      }));

  return (
    <Link
      href={"/resources/" + resource.slug}
      className="group block rounded-xl border border-border bg-bg-elevated p-4 transition-colors hover:bg-bg-hover"
    >
      <div className="flex items-start justify-between">
        <span className="inline-flex h-9 w-9 items-center justify-center rounded-lg bg-accent/10 text-accent">
          <Icon className="h-4 w-4" />
        </span>
        <ArrowUpRight className="h-4 w-4 text-text-muted opacity-0 transition-opacity group-hover:opacity-100" />
      </div>
      <p className="mt-3 text-xs font-medium uppercase tracking-wide text-text-muted">
        Total {label}
      </p>
      <p className="text-2xl font-bold text-foreground">
        {query.isLoading ? (
          <span className="text-text-muted">—</span>
        ) : (
          total.toLocaleString()
        )}
      </p>
      {/* Always render the sparkline space so the card height is
          stable, even before data lands. */}
      <div className="mt-2 h-12 w-full">
        <Suspense fallback={null}>
          <ResourceSparkline id={"spark-" + resource.slug} data={sparkData} />
        </Suspense>
      </div>
      <p className="mt-1 text-[11px] text-text-muted">Last 30 days</p>
    </Link>
  );
}
