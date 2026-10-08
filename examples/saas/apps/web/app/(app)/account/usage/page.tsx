"use client";

import { useQuery } from "@tanstack/react-query";
import { Activity } from "lucide-react";
import { api } from "@/lib/api";
import type { UsageRecord } from "@repo/shared/types";
import { API_ROUTES } from "@repo/shared/constants";

/**
 * Metered usage, from the customer's side.
 *
 * Same story as billing: one owner-scoped endpoint, no portal-only API. What is
 * worth copying from this page is the breakdown, because it is answered by the
 * server rather than by adding up rows here.
 *
 * `?breakdown=metric` returns a count per value of a column on the same query
 * the list is running. It is the parameter the admin's Insights panel uses, it
 * works on any list endpoint built on paginate.List, and because the query it
 * rides on is already scoped to the signed-in customer, the totals are theirs.
 * Summing on the client would mean fetching every row to add them up, which is
 * fine at twenty records and a slow page at twenty thousand.
 */

type Slice = { value: string; count: number };

export default function UsagePage() {
  // page_size 1 because the rows are not what this page wants: the breakdown
  // in meta is. One row comes back so the response shape stays the same.
  const breakdown = useQuery({
    queryKey: ["portal", "usage", "breakdown"],
    queryFn: async () => {
      const { data } = await api.get(`${API_ROUTES.USAGE_RECORDS.LIST}?page_size=1&breakdown=metric`);
      return (data?.meta?.breakdown?.metric ?? []) as Slice[];
    },
  });

  const recent = useQuery({
    queryKey: ["portal", "usage", "recent"],
    queryFn: async () => {
      const { data } = await api.get(`${API_ROUTES.USAGE_RECORDS.LIST}?page_size=10&sort=recorded_on&order=desc`);
      return (data?.data ?? []) as UsageRecord[];
    },
  });

  const slices = breakdown.data ?? [];

  return (
    <div className="space-y-8">
      <div className="flex items-center gap-4">
        <span className="flex h-11 w-11 items-center justify-center rounded-xl bg-accent/10 text-accent">
          <Activity className="h-5 w-5" />
        </span>
        <div>
          <h1 className="text-xl font-bold tracking-tight text-foreground">Usage</h1>
          <p className="text-sm text-text-secondary">What you have used, by metric</p>
        </div>
      </div>

      {breakdown.isLoading ? (
        <p className="text-sm text-text-secondary">Loading your usage...</p>
      ) : slices.length > 0 ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {slices.map((slice) => (
            <div key={slice.value} className="rounded-xl border border-border bg-background p-5">
              <p className="text-[11px] font-semibold uppercase tracking-wider text-text-muted">
                {slice.value || "Unlabelled"}
              </p>
              <p className="mt-2 text-2xl font-bold text-foreground tabular-nums">
                {slice.count.toLocaleString()}
              </p>
              <p className="mt-1 text-xs text-text-secondary">records this period</p>
            </div>
          ))}
        </div>
      ) : (
        <p className="rounded-xl border border-border bg-background p-6 text-sm text-text-secondary">
          Nothing recorded yet. Usage shows here as it is metered.
        </p>
      )}

      <section className="space-y-4">
        <h2 className="text-sm font-semibold uppercase tracking-wider text-text-secondary">Most recent</h2>
        {recent.isLoading ? (
          <p className="text-sm text-text-secondary">Loading...</p>
        ) : recent.data && recent.data.length > 0 ? (
          <div className="overflow-hidden rounded-xl border border-border">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border bg-bg-secondary text-left">
                  <th className="px-4 py-3 font-medium text-text-secondary">Metric</th>
                  <th className="px-4 py-3 font-medium text-text-secondary">Recorded</th>
                  <th className="px-4 py-3 text-right font-medium text-text-secondary">Quantity</th>
                </tr>
              </thead>
              <tbody>
                {recent.data.map((row) => (
                  <tr key={row.id} className="border-b border-border last:border-0">
                    <td className="px-4 py-3 font-medium text-foreground">{row.metric}</td>
                    <td className="px-4 py-3 text-text-secondary">
                      {row.recorded_on
                        ? new Date(row.recorded_on).toLocaleDateString(undefined, {
                            year: "numeric",
                            month: "short",
                            day: "numeric",
                          })
                        : "--"}
                    </td>
                    <td className="px-4 py-3 text-right text-foreground tabular-nums">
                      {row.quantity.toLocaleString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <p className="rounded-xl border border-border bg-background p-6 text-sm text-text-secondary">
            No records yet.
          </p>
        )}
      </section>
    </div>
  );
}
