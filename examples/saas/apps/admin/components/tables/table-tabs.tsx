"use client";

import { useMemo, useRef } from "react";
import { useQueries } from "@tanstack/react-query";
import type { TableTab } from "@/lib/resource";
import { apiClient } from "@/lib/api-client";
import { resourceKeys } from "@/hooks/use-resource";
import { getIcon } from "@/lib/icons";

/*
 * Filter presets as a tab strip.
 *
 * A real tablist, not a row of buttons that happen to filter. That means
 * arrow keys move between tabs and Tab leaves the group, which matters here
 * more than it looks: without roving focus a keyboard user walks through every
 * filter on the way to the table, and with six tabs that is six stops before
 * reaching the thing being filtered.
 *
 * The panel these control is the table, so the table carries the tabpanel role
 * and is labelled by the active tab.
 *
 * Counts are opt-in and arrive late. The badge is rendered only once its number
 * is known rather than showing a spinner or a zero, because a tab that says 0
 * and then says 47 is worse than a tab that said nothing for a moment. Each
 * count is one request with page_size=1, reading meta.total, cached under the
 * resource's endpoint so a create, edit or delete refreshes the badges too.
 */

export interface TableTabsProps {
  tabs: TableTab[];
  active: string;
  onChange: (key: string) => void;
  /** Endpoint the counts are fetched from, when a tab asks for one. */
  endpoint: string;
  /** Filters already applied outside the tabs, so a count matches the table. */
  baseFilters?: Record<string, string>;
}

export function TableTabs({ tabs, active, onChange, endpoint, baseFilters }: TableTabsProps) {
  const refs = useRef<(HTMLButtonElement | null)[]>([]);
  const counted = useMemo(() => tabs.filter((tab) => tab.count), [tabs]);

  // Queries rather than an effect, so the counts live in the cache under the
  // endpoint: saving a record refreshes them, and a tab strip mounted twice
  // shares its requests.
  const results = useQueries({
    queries: counted.map((tab) => {
      const filters: Record<string, string> = { ...(baseFilters ?? {}), ...(tab.filters ?? {}) };
      return {
        queryKey: resourceKeys.tabCount(endpoint, tab.key, filters),
        queryFn: async ({ signal }: { signal: AbortSignal }) => {
          const params = new URLSearchParams({ ...filters, page_size: "1" });
          const { data } = await apiClient.get(endpoint + "?" + params.toString(), { signal });
          return Number(data?.meta?.total ?? 0);
        },
        // A count that fails is a missing badge, not a broken page.
        retry: false,
      };
    }),
  });
  const counts: Record<string, number> = {};
  counted.forEach((tab, index) => {
    const total = results[index]?.data;
    if (typeof total === "number") counts[tab.key] = total;
  });

  if (tabs.length === 0) return null;

  // Roving focus: Left and Right move between tabs, Home and End jump to the
  // ends, and Tab leaves the strip entirely.
  function onKeyDown(event: React.KeyboardEvent, index: number) {
    const last = tabs.length - 1;
    let next = index;
    if (event.key === "ArrowRight") next = index === last ? 0 : index + 1;
    else if (event.key === "ArrowLeft") next = index === 0 ? last : index - 1;
    else if (event.key === "Home") next = 0;
    else if (event.key === "End") next = last;
    else return;
    event.preventDefault();
    onChange(tabs[next].key);
    refs.current[next]?.focus();
  }

  return (
    <div
      role="tablist"
      aria-label="Filter presets"
      className="flex flex-wrap gap-1 border-b border-border px-3 pt-3"
    >
      {tabs.map((tab, index) => {
        const selected = tab.key === active;
        const Icon = tab.icon ? getIcon(tab.icon) : null;
        const count = counts[tab.key];
        return (
          <button
            key={tab.key}
            ref={(node) => {
              refs.current[index] = node;
            }}
            type="button"
            role="tab"
            id={"table-tab-" + tab.key}
            aria-selected={selected}
            aria-controls="table-panel"
            tabIndex={selected ? 0 : -1}
            onClick={() => onChange(tab.key)}
            onKeyDown={(event) => onKeyDown(event, index)}
            className={
              "inline-flex min-h-10 items-center gap-2 rounded-t-lg border-b-2 px-3 text-sm transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent " +
              (selected
                ? "border-accent font-medium text-accent"
                : "border-transparent text-text-secondary hover:bg-bg-hover hover:text-foreground")
            }
          >
            {Icon && <Icon className="h-3.5 w-3.5" aria-hidden="true" />}
            {tab.label}
            {typeof count === "number" && (
              <span
                className={
                  "rounded-full px-1.5 py-0.5 text-xs tabular-nums " +
                  (selected ? "bg-accent/15 text-accent" : "bg-bg-tertiary text-text-secondary")
                }
              >
                {count}
                {/* The number alone is ambiguous next to a label. */}
                <span className="sr-only"> matching</span>
              </span>
            )}
          </button>
        );
      })}
    </div>
  );
}
