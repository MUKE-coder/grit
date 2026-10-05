"use client";

import Link from "next/link";
import {
  useResourceDashboardStats,
  type ResourceDashboardStats,
} from "@/hooks/use-resource";
import { dateRangeToQueryParams, type DateRange } from "@/components/tables/date-filter";
import { renderCell } from "@/components/tables/cell-renderers";
import { FolderOpen } from "@/lib/icons";
import type { ResourceDefinition, ColumnDefinition } from "@/lib/resource";

interface Props {
  resource: ResourceDefinition;
  dateRange: DateRange;
  /** Rows shown, up to DASHBOARD_LATEST_LIMIT. */
  limit?: number;
}

// The table's part of the response the stat card shares.
function pickLatest(stats: ResourceDashboardStats) {
  return stats.latest;
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
  // The same query as the stat card beside it, so the two cost one request.
  const query = useResourceDashboardStats(resource.endpoint, dateRangeToQueryParams(dateRange), pickLatest);

  const cols = pickPreviewColumns(resource);
  const rows = query.data ?? [];
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
          className="inline-flex min-h-[24px] items-center text-xs font-medium text-accent hover:text-accent-hover"
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
