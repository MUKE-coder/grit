"use client";

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
import { resources } from "@/resources";
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
  // A chart is saved with the admin slug (purchase-requests), and the API
  // registers charts under its own resource name (purchase_requests), the
  // last segment of the endpoint. Asking by slug failed every multi-word
  // resource.
  const apiName =
    resources.find((r) => r.slug === chart.resource)?.endpoint.split("/").filter(Boolean).pop() ??
    chart.resource;
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
        apiName +
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
