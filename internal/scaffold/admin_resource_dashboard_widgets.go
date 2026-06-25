package scaffold

// admin_resource_dashboard_widgets.go scaffolds the v3.31.44
// per-resource dashboard widgets that the redesigned dashboard
// renders below the global stat row.
//
// Two widgets per resource:
//   - ResourceStatCard       — total + 30-day sparkline (small)
//   - ResourceLatestTable    — newest N rows, columns auto-picked
//
// Both widgets read from /api/admin/dashboard/resource-stats/:resource
// and accept a DateRange prop so the dashboard's top-level DateFilter
// scopes the count + the latest list. The sparkline is intentionally
// always-30-days (see services.buildDailySeries comment) so it never
// collapses to a single bar.
//
// The dashboard page renders each resource as one row (stat + latest
// side-by-side on desktop, stacked on mobile). Resources can opt out
// by setting `dashboard: { enabled: false }` in their definition.

import (
	"path/filepath"
)

func writeAdminResourceDashboardWidgets(root string, opts Options) error {
	_ = opts
	adminRoot := filepath.Join(root, "apps", "admin")

	files := map[string]string{
		filepath.Join(adminRoot, "components", "dashboard", "ResourceStatCard.tsx"):    adminResourceStatCardTSX(),
		filepath.Join(adminRoot, "components", "dashboard", "ResourceLatestTable.tsx"): adminResourceLatestTableTSX(),
		filepath.Join(adminRoot, "components", "dashboard", "ResourceWidgetsRow.tsx"):  adminResourceWidgetsRowTSX(),
	}
	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return err
		}
	}
	return nil
}

// adminResourceStatCardTSX emits the small stat card with sparkline.
// One React Query per render keyed on (resource, dateRange) — the
// dashboard's DateFilter triggers a refetch via the query key.
func adminResourceStatCardTSX() string {
	return `"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import {
  ResponsiveContainer, AreaChart, Area, Tooltip,
} from "recharts";
import { apiClient } from "@/lib/api-client";
import { dateRangeToQueryParams, type DateRange } from "@/components/tables/date-filter";
import { getIcon, ArrowUpRight } from "@/lib/icons";
import type { ResourceDefinition } from "@/lib/resource";

// One sparkline bucket = one calendar day. Always 30 buckets so the
// chart shape stays stable; counts inside the active date range only
// affect the "Total" number, not the sparkline window.
interface ResourceStatsBucket {
  date: string;
  count: number;
}

interface ResourceStatsResponse {
  data: {
    resource: string;
    total: number;
    series: ResourceStatsBucket[];
    latest: Record<string, unknown>[];
  };
}

interface Props {
  resource: ResourceDefinition;
  dateRange: DateRange;
}

export function ResourceStatCard({ resource, dateRange }: Props) {
  const params = dateRangeToQueryParams(dateRange);
  const query = useQuery<ResourceStatsResponse["data"]>({
    queryKey: ["dashboard", "resource-stats", resource.slug, params],
    queryFn: async () => {
      const search = new URLSearchParams(params).toString();
      const url =
        "/api/admin/dashboard/resource-stats/" +
        resource.slug +
        (search ? "?" + search : "");
      const { data } = await apiClient.get<ResourceStatsResponse>(url);
      return data.data;
    },
    // The dashboard cycles between resources quickly; keep stats
    // around so re-opening the page doesn't re-flash the skeleton.
    staleTime: 30_000,
    refetchInterval: 60_000,
  });

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
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={sparkData} margin={{ top: 4, right: 0, bottom: 0, left: 0 }}>
            <defs>
              <linearGradient id={"spark-" + resource.slug} x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor="var(--accent)" stopOpacity={0.45} />
                <stop offset="100%" stopColor="var(--accent)" stopOpacity={0} />
              </linearGradient>
            </defs>
            <Tooltip
              contentStyle={{
                background: "var(--bg-elevated)",
                border: "1px solid var(--border)",
                borderRadius: 6,
                fontSize: 11,
                padding: "4px 8px",
              }}
              labelStyle={{ color: "var(--text-secondary)" }}
              itemStyle={{ color: "var(--foreground)" }}
              cursor={{ stroke: "var(--accent)", strokeOpacity: 0.3 }}
              formatter={(value: number) => [value + " new", "Count"]}
              labelFormatter={(d: string) => d}
            />
            <Area
              type="monotone"
              dataKey="count"
              stroke="var(--accent)"
              strokeWidth={1.5}
              fill={"url(#spark-" + resource.slug + ")"}
              isAnimationActive={false}
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
      <p className="mt-1 text-[11px] text-text-muted">Last 30 days</p>
    </Link>
  );
}
`
}

// adminResourceLatestTableTSX emits the compact "Latest N" table used
// alongside the stat card. Columns are auto-picked from the resource
// definition's table.columns -- whichever first 2-3 are non-id keys.
func adminResourceLatestTableTSX() string {
	return `"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { dateRangeToQueryParams, type DateRange } from "@/components/tables/date-filter";
import { renderCell } from "@/components/tables/cell-renderers";
import { FolderOpen } from "@/lib/icons";
import type { ResourceDefinition, ColumnDefinition } from "@/lib/resource";

interface ResourceStatsResponse {
  data: {
    resource: string;
    total: number;
    series: { date: string; count: number }[];
    latest: Record<string, unknown>[];
  };
}

interface Props {
  resource: ResourceDefinition;
  dateRange: DateRange;
  limit?: number;
}

// v3.31.46 -- pickPreviewColumns now returns full ColumnDefinition
// objects so the table can use the same renderCell pipeline as the
// resource list page. That brings proper FileRef thumbnails, badge
// pills, date formatting, and currency rendering -- instead of the
// previous "stringify and truncate" fallback that turned image refs
// into a JSON blob.
//
// Heuristic: prefer recognisable columns (name/title + email/status/
// price), then fall back to the first non-id columns. Always include
// any image / file columns so users see thumbnails on the dashboard
// without having to customise.
function pickPreviewColumns(resource: ResourceDefinition): ColumnDefinition[] {
  const all = (resource.table?.columns ?? []) as ColumnDefinition[];
  if (all.length === 0) return [];

  // Image / file columns always make the cut -- thumbnails carry a
  // lot of signal in a small cell.
  const imageCols = all.filter(
    (c) => c.format === "image" || c.format === "file" || c.format === "files",
  );

  const priorityKeys = ["name", "title", "subject", "email", "status", "price"];
  const picked: ColumnDefinition[] = [...imageCols];
  for (const key of priorityKeys) {
    const hit = all.find((c) => c.key === key);
    if (hit && !picked.find((p) => p.key === hit.key)) {
      picked.push(hit);
    }
    if (picked.length >= 4) break;
  }
  if (picked.length < 4) {
    for (const c of all) {
      if (c.key === "id" || c.key.endsWith("_id")) continue;
      if (c.hidden) continue;
      if (picked.find((p) => p.key === c.key)) continue;
      picked.push(c);
      if (picked.length >= 4) break;
    }
  }
  return picked.slice(0, 4);
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

export function ResourceLatestTable({ resource, dateRange, limit = 5 }: Props) {
  const params = { ...dateRangeToQueryParams(dateRange), limit: String(limit) };
  const query = useQuery<ResourceStatsResponse["data"]>({
    queryKey: ["dashboard", "resource-latest", resource.slug, params],
    queryFn: async () => {
      const search = new URLSearchParams(params).toString();
      const url =
        "/api/admin/dashboard/resource-stats/" +
        resource.slug +
        (search ? "?" + search : "");
      const { data } = await apiClient.get<ResourceStatsResponse>(url);
      return data.data;
    },
    staleTime: 30_000,
    refetchInterval: 60_000,
  });

  const cols = pickPreviewColumns(resource);
  const rows = query.data?.latest ?? [];
  const label = resource.label?.plural ?? resource.slug;

  return (
    <div className="rounded-xl border border-border bg-bg-elevated">
      <header className="flex items-center justify-between border-b border-border px-4 py-3">
        <div>
          <p className="text-sm font-semibold text-foreground">Latest {label}</p>
          <p className="text-xs text-text-muted">
            Newest first within the selected range
          </p>
        </div>
        <Link
          href={"/resources/" + resource.slug}
          className="text-xs font-medium text-accent hover:text-accent-hover"
        >
          View all
        </Link>
      </header>
      {query.isLoading ? (
        <div className="px-4 py-10 text-center text-sm text-text-muted">Loading…</div>
      ) : rows.length === 0 ? (
        <div className="flex flex-col items-center justify-center gap-2 px-4 py-10 text-center">
          <FolderOpen className="h-5 w-5 text-text-muted" />
          <p className="text-sm text-text-muted">No {label.toLowerCase()} in this range yet.</p>
        </div>
      ) : (
        // v3.31.46 -- proper table layout. Reuses the same renderCell
        // dispatch the resource list page uses, so format hints
        // ("image", "badge", "currency", "date", "relative") get the
        // same treatment here as in the main resource view.
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border text-left">
                {cols.map((c) => (
                  <th
                    key={c.key}
                    className="px-4 py-2 text-[11px] font-semibold uppercase tracking-wider text-text-muted"
                  >
                    {c.label}
                  </th>
                ))}
                <th className="px-4 py-2 text-right text-[11px] font-semibold uppercase tracking-wider text-text-muted">
                  Created
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {rows.slice(0, limit).map((row, idx) => {
                const id = String(row.id ?? idx);
                const createdAt = String(row.created_at ?? "");
                return (
                  <tr key={id} className="transition-colors hover:bg-bg-hover">
                    {cols.map((c) => (
                      <td key={c.key} className="px-4 py-2.5 text-foreground">
                        {renderCell(c, row[c.key], row)}
                      </td>
                    ))}
                    <td className="px-4 py-2.5 text-right text-xs text-text-muted whitespace-nowrap">
                      {createdAt ? timeAgo(createdAt) : "—"}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
`
}

// adminResourceWidgetsRowTSX emits the wrapper that puts the stat
// card + latest table side-by-side. The dashboard page maps over
// resources and renders one of these per resource.
func adminResourceWidgetsRowTSX() string {
	return `"use client";

import { useState } from "react";
import { ResourceStatCard } from "@/components/dashboard/ResourceStatCard";
import { ResourceLatestTable } from "@/components/dashboard/ResourceLatestTable";
import type { DateRange } from "@/components/tables/date-filter";
import type { ResourceDefinition } from "@/lib/resource";

interface Props {
  resource: ResourceDefinition;
  dateRange: DateRange;
  // v3.31.45 -- show / hide each half independently. Both default
  // to true so existing call sites keep working without changes.
  // When only one is shown it stretches to fill the row.
  showStat?: boolean;
  showLatest?: boolean;
  // v3.31.46 -- layout mode. "split" = side-by-side (default);
  // "tabs" = each widget full-width inside its own tab.
  layout?: "split" | "tabs";
}

export function ResourceWidgetsRow({
  resource,
  dateRange,
  showStat = true,
  showLatest = true,
  layout = "split",
}: Props) {
  if (!showStat && !showLatest) return null;

  // Tabs mode only makes sense when both halves are enabled --
  // otherwise there's nothing to switch between, so fall through to
  // the single-pane render below.
  if (layout === "tabs" && showStat && showLatest) {
    return <ResourceTabs resource={resource} dateRange={dateRange} />;
  }

  if (showStat && !showLatest) {
    return (
      <div className="grid grid-cols-1">
        <ResourceStatCard resource={resource} dateRange={dateRange} />
      </div>
    );
  }
  if (!showStat && showLatest) {
    return (
      <div className="grid grid-cols-1">
        <ResourceLatestTable resource={resource} dateRange={dateRange} />
      </div>
    );
  }
  return (
    <div className="grid grid-cols-1 gap-4 lg:grid-cols-3">
      <div className="lg:col-span-1">
        <ResourceStatCard resource={resource} dateRange={dateRange} />
      </div>
      <div className="lg:col-span-2">
        <ResourceLatestTable resource={resource} dateRange={dateRange} />
      </div>
    </div>
  );
}

// v3.31.46 -- ResourceTabs renders the same Total + Latest widgets
// inside a tabbed container. Each tab body is full-width so the
// Latest table can use the entire dashboard row width.
function ResourceTabs({
  resource,
  dateRange,
}: {
  resource: ResourceDefinition;
  dateRange: DateRange;
}) {
  // Default to "latest" because that's the tab the user is most
  // likely opening tabs mode for in the first place -- the stat card
  // doesn't need full width.
  const [active, setActive] = useState<"total" | "latest">("latest");
  const label = resource.label?.plural ?? resource.slug;

  return (
    <div className="rounded-xl border border-border bg-bg-elevated">
      <div className="flex items-center gap-1 border-b border-border px-2 py-2">
        <TabButton
          active={active === "total"}
          onClick={() => setActive("total")}
          label={"Total " + label}
        />
        <TabButton
          active={active === "latest"}
          onClick={() => setActive("latest")}
          label={"Latest " + label}
        />
      </div>
      <div className="p-3">
        {active === "total" ? (
          <ResourceStatCard resource={resource} dateRange={dateRange} />
        ) : (
          <ResourceLatestTable resource={resource} dateRange={dateRange} limit={10} />
        )}
      </div>
    </div>
  );
}

function TabButton({
  active,
  onClick,
  label,
}: {
  active: boolean;
  onClick: () => void;
  label: string;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={
        "rounded-md px-3 py-1.5 text-sm font-medium transition-colors " +
        (active
          ? "bg-accent/15 text-accent"
          : "text-text-secondary hover:bg-bg-hover hover:text-foreground")
      }
    >
      {label}
    </button>
  );
}
`
}
