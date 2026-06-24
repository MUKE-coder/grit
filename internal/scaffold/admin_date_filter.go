package scaffold

// admin_date_filter.go emits the DateFilter component used by the
// TableToolbar to scope a resource's list (and stats) to a date
// window. Five quick presets + a custom range; selection is held by
// the parent (typically lifted to URL search params via ResourceListView).
//
// The component is deliberately presentation-only -- it doesn't read or
// write the URL itself. That keeps it reusable in non-resource contexts
// (e.g. a dashboard chart) and lets the resource page own URL plumbing
// without the DateFilter learning about Next.js routing.

func adminDateFilter() string {
	return `"use client";

import { useEffect, useRef, useState } from "react";
import { Calendar, ChevronDown, X } from "@/lib/icons";

// DateRange — the shape every consumer holds. A preset short-circuits
// the API call (?created_since=7d); custom passes explicit
// ?created_from=&created_to=. Empty = no filter (= "All time").
export type DateRange = {
  preset?: "today" | "7d" | "30d" | "month" | "custom";
  from?: string; // ISO date YYYY-MM-DD (custom mode only)
  to?: string;
};

interface DateFilterProps {
  value: DateRange;
  onChange: (next: DateRange) => void;
  label?: string;
}

const PRESETS: { key: DateRange["preset"]; label: string }[] = [
  { key: "today", label: "Today" },
  { key: "7d", label: "Last 7 days" },
  { key: "30d", label: "Last 30 days" },
  { key: "month", label: "This month" },
];

function describe(value: DateRange, fallback: string): string {
  if (!value.preset) return fallback;
  if (value.preset === "custom") {
    if (value.from && value.to) return value.from + " - " + value.to;
    if (value.from) return "From " + value.from;
    if (value.to) return "Until " + value.to;
    return "Custom range";
  }
  const hit = PRESETS.find((p) => p.key === value.preset);
  return hit?.label ?? fallback;
}

export function DateFilter({ value, onChange, label = "Date" }: DateFilterProps) {
  const [open, setOpen] = useState(false);
  const [customFrom, setCustomFrom] = useState(value.from ?? "");
  const [customTo, setCustomTo] = useState(value.to ?? "");
  const rootRef = useRef<HTMLDivElement>(null);

  // Close on outside-click. We attach to document so we catch presses
  // anywhere -- including other open popovers -- without an overlay
  // element that would block scroll-wheel pass-through.
  useEffect(() => {
    if (!open) return;
    function handle(e: MouseEvent) {
      if (rootRef.current && !rootRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handle);
    return () => document.removeEventListener("mousedown", handle);
  }, [open]);

  // Keep the custom inputs in sync if the parent state changes (e.g.
  // user clicks a preset elsewhere or URL navigation rehydrates).
  useEffect(() => {
    setCustomFrom(value.from ?? "");
    setCustomTo(value.to ?? "");
  }, [value.from, value.to]);

  const isActive = !!value.preset;
  const buttonText = describe(value, "All time");

  return (
    <div ref={rootRef} className="relative">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        className={
          "flex items-center gap-2 rounded-lg border px-3 py-2 text-sm transition-colors " +
          (isActive
            ? "border-accent bg-accent/10 text-accent"
            : "border-border bg-bg-tertiary text-text-secondary hover:text-foreground")
        }
      >
        <Calendar className="h-4 w-4" />
        <span>{buttonText}</span>
        {isActive ? (
          <span
            role="button"
            tabIndex={0}
            onClick={(e) => {
              e.stopPropagation();
              onChange({});
            }}
            onKeyDown={(e) => {
              if (e.key === "Enter" || e.key === " ") {
                e.preventDefault();
                e.stopPropagation();
                onChange({});
              }
            }}
            className="ml-1 rounded p-0.5 hover:bg-accent/20 cursor-pointer"
            aria-label="Clear date filter"
          >
            <X className="h-3 w-3" />
          </span>
        ) : (
          <ChevronDown className="h-3.5 w-3.5" />
        )}
      </button>

      {open && (
        <div className="absolute left-0 top-full z-30 mt-1 w-64 rounded-lg border border-border bg-bg-elevated shadow-lg">
          <div className="border-b border-border px-3 py-2 text-[11px] font-semibold uppercase tracking-wider text-text-muted">
            {label}
          </div>
          <div className="p-1">
            {PRESETS.map((p) => (
              <button
                key={p.key}
                type="button"
                onClick={() => {
                  onChange({ preset: p.key });
                  setOpen(false);
                }}
                className={
                  "flex w-full items-center justify-between rounded px-2 py-1.5 text-sm transition-colors " +
                  (value.preset === p.key
                    ? "bg-accent/10 text-accent"
                    : "text-foreground hover:bg-bg-hover")
                }
              >
                <span>{p.label}</span>
              </button>
            ))}
          </div>
          <div className="border-t border-border p-3 space-y-2">
            <p className="text-[11px] font-semibold uppercase tracking-wider text-text-muted">
              Custom range
            </p>
            <div className="grid grid-cols-2 gap-2">
              <div>
                <label className="block text-[10px] text-text-muted mb-0.5">From</label>
                <input
                  type="date"
                  value={customFrom}
                  onChange={(e) => setCustomFrom(e.target.value)}
                  className="w-full rounded border border-border bg-bg-secondary px-2 py-1 text-xs text-foreground focus:border-accent focus:outline-none"
                />
              </div>
              <div>
                <label className="block text-[10px] text-text-muted mb-0.5">To</label>
                <input
                  type="date"
                  value={customTo}
                  onChange={(e) => setCustomTo(e.target.value)}
                  className="w-full rounded border border-border bg-bg-secondary px-2 py-1 text-xs text-foreground focus:border-accent focus:outline-none"
                />
              </div>
            </div>
            <button
              type="button"
              disabled={!customFrom && !customTo}
              onClick={() => {
                onChange({
                  preset: "custom",
                  from: customFrom || undefined,
                  to: customTo || undefined,
                });
                setOpen(false);
              }}
              className="w-full rounded bg-accent px-3 py-1.5 text-xs font-medium text-white hover:bg-accent-hover transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
            >
              Apply
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

// dateRangeToQueryParams — maps a DateRange to the API's three
// supported query params. Returns a plain object the caller can spread
// into a URLSearchParams or fetch call. Empty range returns {}.
//
// Centralised so the resource list, stats hook, and any custom data
// hook agree on the wire shape and don't drift over time.
export function dateRangeToQueryParams(range: DateRange): Record<string, string> {
  if (!range.preset) return {};
  if (range.preset === "custom") {
    const out: Record<string, string> = {};
    if (range.from) out.created_from = range.from;
    if (range.to) out.created_to = range.to;
    return out;
  }
  const map: Record<string, string> = {
    today: "1d",
    "7d": "7d",
    "30d": "30d",
    month: "30d", // month preset is approximated server-side as last 30 days
  };
  const since = map[range.preset];
  return since ? { created_since: since } : {};
}
`
}
