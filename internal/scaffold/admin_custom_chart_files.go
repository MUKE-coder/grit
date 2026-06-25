package scaffold

// admin_custom_chart_files.go scaffolds the v3.31.47 Preset Chart
// builder client-side pieces:
//
//   - CustomChartCard.tsx       -- renders one saved chart with Recharts
//   - ChartBuilderForm.tsx      -- the inline form used in Settings
//
// Both files are written into apps/admin/components/dashboard/ next
// to the v3.31.44 ResourceStatCard / ResourceLatestTable.

import (
	"path/filepath"
)

func writeAdminCustomChartFiles(root string, opts Options) error {
	_ = opts
	adminRoot := filepath.Join(root, "apps", "admin")

	files := map[string]string{
		filepath.Join(adminRoot, "components", "dashboard", "CustomChartCard.tsx"):  adminCustomChartCardTSX(),
		filepath.Join(adminRoot, "components", "dashboard", "ChartBuilderForm.tsx"): adminChartBuilderFormTSX(),
	}
	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return err
		}
	}
	return nil
}

func adminCustomChartCardTSX() string {
	return `"use client";

// v3.31.47 -- one card per saved CustomChart. Pulls data from
// /api/admin/dashboard/chart/:resource on render and dispatches to
// the right Recharts component based on the saved viz. Loading +
// error states are rendered inline so a broken chart never blanks
// the whole dashboard section.

import { useQuery } from "@tanstack/react-query";
import {
  BarChart, Bar,
  LineChart, Line,
  AreaChart, Area,
  PieChart, Pie, Cell,
  XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
} from "recharts";
import { apiClient } from "@/lib/api-client";
import { dateRangeToQueryParams, type DateRange } from "@/components/tables/date-filter";
import { AlertCircle, TrendingUp } from "@/lib/icons";
import type { CustomChart, ChartViz } from "@/lib/dashboard-catalog";

interface ChartRow {
  x: string | number;
  y: number;
}

interface ChartResponse {
  data: {
    preset: string;
    rows: ChartRow[];
    meta: Record<string, unknown>;
  };
}

interface Props {
  chart: CustomChart;
  dateRange: DateRange;
}

const PIE_COLORS = [
  "var(--accent)",
  "var(--info)",
  "var(--success)",
  "var(--warning)",
  "var(--danger)",
  "#a78bfa",
  "#f472b6",
  "#34d399",
];

export function CustomChartCard({ chart, dateRange }: Props) {
  const params = {
    preset: chart.preset,
    ...(chart.field ? { field: chart.field } : {}),
    ...(chart.limit ? { limit: String(chart.limit) } : {}),
    ...(chart.grain ? { grain: chart.grain } : {}),
    ...dateRangeToQueryParams(dateRange),
  };

  const query = useQuery<ChartResponse["data"]>({
    queryKey: ["dashboard", "chart", chart.id, params],
    queryFn: async () => {
      const search = new URLSearchParams(params).toString();
      const url =
        "/api/admin/dashboard/chart/" +
        chart.resource +
        (search ? "?" + search : "");
      const { data } = await apiClient.get<ChartResponse>(url);
      return data.data;
    },
    staleTime: 30_000,
    refetchInterval: 60_000,
  });

  return (
    <div className="rounded-xl border border-border bg-bg-elevated p-5">
      <div className="mb-4 flex items-center justify-between">
        <div className="min-w-0">
          <p className="truncate text-sm font-semibold text-foreground">
            {chart.title || chart.resource}
          </p>
          <p className="text-xs text-text-muted">{describeChart(chart)}</p>
        </div>
        <TrendingUp className="h-4 w-4 text-text-muted shrink-0" />
      </div>
      <div className="h-56 w-full">
        {query.isLoading ? (
          <div className="flex h-full items-center justify-center text-sm text-text-muted">
            Loading…
          </div>
        ) : query.isError ? (
          <div className="flex h-full flex-col items-center justify-center gap-1.5 text-center">
            <AlertCircle className="h-5 w-5 text-danger" />
            <p className="text-xs text-text-muted">Couldn&apos;t load this chart.</p>
            <p className="max-w-[260px] truncate text-[11px] text-text-muted">
              {(query.error as Error)?.message ?? "Unknown error"}
            </p>
          </div>
        ) : (query.data?.rows ?? []).length === 0 ? (
          <div className="flex h-full items-center justify-center text-sm text-text-muted">
            No data yet.
          </div>
        ) : (
          <ResponsiveContainer width="100%" height="100%">
            {renderChart(chart.viz, query.data!.rows, chart.id)}
          </ResponsiveContainer>
        )}
      </div>
    </div>
  );
}

function describeChart(c: CustomChart): string {
  switch (c.preset) {
    case "count_over_time":
      return "New records per day, last 30 days";
    case "group_by":
      return "Top " + (c.limit ?? 10) + " by " + (c.field ?? "—");
    case "sum_over_time":
      return "Sum of " + (c.field ?? "—") + " per day";
    case "avg_over_time":
      return "Average " + (c.field ?? "—") + " per day";
  }
}

function renderChart(viz: ChartViz, rows: ChartRow[], id: string) {
  switch (viz) {
    case "bar":
      return (
        <BarChart data={rows}>
          <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" />
          <XAxis dataKey="x" stroke="var(--text-muted)" fontSize={11} tickLine={false} axisLine={false} />
          <YAxis stroke="var(--text-muted)" fontSize={11} tickLine={false} axisLine={false} />
          <Tooltip contentStyle={chartTooltipStyle} cursor={{ fill: "var(--bg-hover)" }} />
          <Bar dataKey="y" fill="var(--accent)" radius={[4, 4, 0, 0]} />
        </BarChart>
      );
    case "line":
      return (
        <LineChart data={rows}>
          <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" />
          <XAxis dataKey="x" stroke="var(--text-muted)" fontSize={11} tickLine={false} axisLine={false} />
          <YAxis stroke="var(--text-muted)" fontSize={11} tickLine={false} axisLine={false} />
          <Tooltip contentStyle={chartTooltipStyle} />
          <Line type="monotone" dataKey="y" stroke="var(--accent)" strokeWidth={2} dot={false} />
        </LineChart>
      );
    case "area":
      return (
        <AreaChart data={rows}>
          <defs>
            <linearGradient id={"area-" + id} x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stopColor="var(--accent)" stopOpacity={0.35} />
              <stop offset="100%" stopColor="var(--accent)" stopOpacity={0} />
            </linearGradient>
          </defs>
          <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" />
          <XAxis dataKey="x" stroke="var(--text-muted)" fontSize={11} tickLine={false} axisLine={false} />
          <YAxis stroke="var(--text-muted)" fontSize={11} tickLine={false} axisLine={false} />
          <Tooltip contentStyle={chartTooltipStyle} />
          <Area type="monotone" dataKey="y" stroke="var(--accent)" strokeWidth={2} fill={"url(#area-" + id + ")"} />
        </AreaChart>
      );
    case "pie":
    case "donut":
      return (
        <PieChart>
          <Pie
            data={rows}
            dataKey="y"
            nameKey="x"
            cx="50%"
            cy="50%"
            innerRadius={viz === "donut" ? 50 : 0}
            outerRadius={80}
            paddingAngle={2}
            label={(entry: { x: string }) => entry.x}
            labelLine={false}
          >
            {rows.map((_, i) => (
              <Cell key={i} fill={PIE_COLORS[i % PIE_COLORS.length]} />
            ))}
          </Pie>
          <Tooltip contentStyle={chartTooltipStyle} />
        </PieChart>
      );
  }
}

const chartTooltipStyle = {
  background: "var(--bg-elevated)",
  border: "1px solid var(--border)",
  borderRadius: 8,
  fontSize: 12,
} as const;
`
}

func adminChartBuilderFormTSX() string {
	return `"use client";

// v3.31.47 -- the Preset Chart builder form. Renders as an inline
// panel inside Dashboard Settings; the parent owns the saved chart
// list and decides whether the form is in "create new" mode (initial
// = undefined) or "edit existing" mode (initial = an existing chart).

import { useEffect, useMemo, useState } from "react";
import {
  CHART_PRESET_LABELS,
  CHART_VIZ_LABELS,
  vizesForPreset,
  type CustomChart,
  type ChartPreset,
  type ChartViz,
} from "@/lib/dashboard-catalog";
import type { ResourceDefinition, ColumnDefinition } from "@/lib/resource";

interface Props {
  resources: ResourceDefinition[];
  initial?: CustomChart;
  onSubmit: (chart: CustomChart) => void;
  onCancel: () => void;
}

// classifyFields buckets a resource's table.columns into the two
// kinds the chart builder cares about. Only "currency" reliably
// tells us numeric; everything else might be either, so the API's
// reflection-based validator is authoritative.
function classifyFields(resource: ResourceDefinition | undefined): {
  categorical: ColumnDefinition[];
  numeric: ColumnDefinition[];
} {
  if (!resource) return { categorical: [], numeric: [] };
  const cols = (resource.table?.columns ?? []) as ColumnDefinition[];
  const categorical: ColumnDefinition[] = [];
  const numeric: ColumnDefinition[] = [];
  for (const c of cols) {
    if (c.hidden) continue;
    if (c.key === "id" || c.key.endsWith("_id")) continue;
    if (
      c.format === "image" ||
      c.format === "file" ||
      c.format === "files" ||
      c.format === "richtext" ||
      c.format === "user"
    ) {
      continue;
    }
    if (c.format === "currency") {
      numeric.push(c);
      continue;
    }
    categorical.push(c);
    if (c.format === undefined || c.format === "text") {
      numeric.push(c);
    }
  }
  return { categorical, numeric };
}

function genId(): string {
  if (typeof crypto !== "undefined" && typeof crypto.randomUUID === "function") {
    return crypto.randomUUID();
  }
  return "chart-" + Math.random().toString(36).slice(2, 10) + "-" + Date.now().toString(36);
}

export function ChartBuilderForm({ resources, initial, onSubmit, onCancel }: Props) {
  const [resourceSlug, setResourceSlug] = useState<string>(initial?.resource ?? resources[0]?.slug ?? "");
  const [preset, setPreset] = useState<ChartPreset>(initial?.preset ?? "count_over_time");
  const [field, setField] = useState<string>(initial?.field ?? "");
  const [viz, setViz] = useState<ChartViz>(initial?.viz ?? "line");
  const [title, setTitle] = useState<string>(initial?.title ?? "");
  const [limit, setLimit] = useState<number>(initial?.limit ?? 10);

  const resource = resources.find((r) => r.slug === resourceSlug);
  const { categorical, numeric } = useMemo(() => classifyFields(resource), [resource]);

  useEffect(() => {
    if (preset === "count_over_time") {
      setField("");
      return;
    }
    const valid = preset === "group_by" ? categorical : numeric;
    if (!valid.find((c) => c.key === field)) {
      setField(valid[0]?.key ?? "");
    }
  }, [preset, resourceSlug]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    const allowed = vizesForPreset(preset);
    if (!allowed.includes(viz)) {
      setViz(allowed[0]);
    }
  }, [preset]); // eslint-disable-line react-hooks/exhaustive-deps

  const fieldRequired = preset !== "count_over_time";
  const fieldList = preset === "group_by" ? categorical : numeric;
  const canSubmit =
    resourceSlug !== "" &&
    (!fieldRequired || field !== "") &&
    title.trim() !== "";

  const handleSubmit = () => {
    if (!canSubmit) return;
    const chart: CustomChart = {
      id: initial?.id ?? genId(),
      title: title.trim(),
      resource: resourceSlug,
      preset,
      viz,
      ...(field ? { field } : {}),
      ...(preset === "group_by" ? { limit } : {}),
    };
    onSubmit(chart);
  };

  return (
    <div className="space-y-4 rounded-xl border border-accent/30 bg-bg-elevated p-5">
      <header>
        <h3 className="text-sm font-semibold text-foreground">
          {initial ? "Edit chart" : "New chart"}
        </h3>
        <p className="text-xs text-text-muted">
          Pick a resource, a preset, and a visualization. The chart
          will appear in the Charts section of your dashboard.
        </p>
      </header>

      <Field label="Title" hint="Shown as the chart's heading.">
        <input
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Orders by status"
          className="w-full rounded-lg border border-border bg-bg-secondary px-3 py-2 text-sm text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
        />
      </Field>

      <Field label="Resource" hint="Which model the chart aggregates.">
        <select
          value={resourceSlug}
          onChange={(e) => setResourceSlug(e.target.value)}
          className="w-full rounded-lg border border-border bg-bg-secondary px-3 py-2 text-sm text-foreground focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
        >
          {resources.map((r) => (
            <option key={r.slug} value={r.slug}>
              {r.label?.plural ?? r.name}
            </option>
          ))}
        </select>
      </Field>

      <Field label="Preset" hint="The aggregation the chart computes.">
        <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
          {(Object.keys(CHART_PRESET_LABELS) as ChartPreset[]).map((p) => {
            const meta = CHART_PRESET_LABELS[p];
            const active = preset === p;
            return (
              <button
                key={p}
                type="button"
                onClick={() => setPreset(p)}
                className={
                  "rounded-lg border p-3 text-left transition-colors " +
                  (active
                    ? "border-accent bg-accent/10"
                    : "border-border bg-bg-tertiary hover:bg-bg-hover")
                }
              >
                <p className="text-sm font-medium text-foreground">{meta.title}</p>
                <p className="text-xs text-text-muted">{meta.hint}</p>
              </button>
            );
          })}
        </div>
      </Field>

      {fieldRequired && (
        <Field
          label="Field"
          hint={
            preset === "group_by"
              ? "Categorical column to group rows by."
              : "Numeric column to aggregate."
          }
        >
          {fieldList.length === 0 ? (
            <p className="rounded-lg border border-warning/40 bg-warning/5 px-3 py-2 text-xs text-warning">
              No suitable columns on this resource. Try a different preset or resource.
            </p>
          ) : (
            <select
              value={field}
              onChange={(e) => setField(e.target.value)}
              className="w-full rounded-lg border border-border bg-bg-secondary px-3 py-2 text-sm text-foreground focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
            >
              {fieldList.map((c) => (
                <option key={c.key} value={c.key}>
                  {c.label}
                </option>
              ))}
            </select>
          )}
        </Field>
      )}

      {preset === "group_by" && (
        <Field label="Top N" hint="How many groups to show. Max 100.">
          <input
            type="number"
            min={1}
            max={100}
            value={limit}
            onChange={(e) => setLimit(Math.max(1, Math.min(100, parseInt(e.target.value, 10) || 1)))}
            className="w-32 rounded-lg border border-border bg-bg-secondary px-3 py-2 text-sm text-foreground focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
          />
        </Field>
      )}

      <Field label="Visualization" hint="How the data renders. Time presets prefer line/area; group_by prefers bar/pie.">
        <div className="flex flex-wrap gap-1.5">
          {(Object.keys(CHART_VIZ_LABELS) as ChartViz[]).map((v) => {
            const allowed = vizesForPreset(preset).includes(v);
            const active = viz === v;
            return (
              <button
                key={v}
                type="button"
                onClick={() => allowed && setViz(v)}
                disabled={!allowed}
                className={
                  "rounded-md border px-3 py-1.5 text-xs font-medium transition-colors " +
                  (active
                    ? "border-accent bg-accent text-white"
                    : allowed
                      ? "border-border bg-bg-tertiary text-text-secondary hover:bg-bg-hover"
                      : "border-border bg-bg-tertiary text-text-muted opacity-40 cursor-not-allowed")
                }
              >
                {CHART_VIZ_LABELS[v]}
              </button>
            );
          })}
        </div>
      </Field>

      <div className="flex items-center justify-end gap-2 border-t border-border pt-4">
        <button
          type="button"
          onClick={onCancel}
          className="rounded-lg border border-border px-3 py-1.5 text-sm text-text-secondary hover:bg-bg-hover"
        >
          Cancel
        </button>
        <button
          type="button"
          onClick={handleSubmit}
          disabled={!canSubmit}
          className="inline-flex items-center gap-1.5 rounded-lg bg-accent px-4 py-1.5 text-sm font-medium text-white hover:bg-accent-hover disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {initial ? "Save chart" : "Add chart"}
        </button>
      </div>
    </div>
  );
}

function Field({
  label,
  hint,
  children,
}: {
  label: string;
  hint?: string;
  children: React.ReactNode;
}) {
  return (
    <div className="space-y-1.5">
      <label className="block text-xs font-medium uppercase tracking-wider text-text-muted">
        {label}
      </label>
      {hint && <p className="text-[11px] text-text-muted">{hint}</p>}
      {children}
    </div>
  );
}
`
}
