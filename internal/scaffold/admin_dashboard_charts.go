package scaffold

// The dashboard's recharts code, in modules of their own. Recharts is about
// 400 KB; imported by the dashboard page and the resource stat card, it rode in
// the first load of the dashboard. The page and the card now load these with
// dynamic() behind Suspense, so the dashboard renders and becomes usable before
// the charts arrive.

// adminDashboardChartsTSX is the dashboard page's activity and severity charts.
func adminDashboardChartsTSX() string {
	return `"use client";

import {
  AreaChart, Area, CartesianGrid, ResponsiveContainer,
  Tooltip, XAxis, YAxis, PieChart, Pie, Cell,
} from "recharts";

// Loaded by the dashboard page with dynamic(), so recharts is not part of the
// page's first load.

export interface ActivityPoint {
  day: string;
  events: number;
}

export interface SeveritySlice {
  name: string;
  value: number;
  color: string;
}

export function ActivityAreaChart({ data }: { data: ActivityPoint[] }) {
  return (
    <ResponsiveContainer width="100%" height="100%">
      <AreaChart data={data}>
        <defs>
          <linearGradient id="activityFill" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="var(--accent)" stopOpacity={0.35} />
            <stop offset="100%" stopColor="var(--accent)" stopOpacity={0} />
          </linearGradient>
        </defs>
        <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" />
        <XAxis dataKey="day" stroke="var(--text-muted)" fontSize={11} tickLine={false} axisLine={false} />
        <YAxis stroke="var(--text-muted)" fontSize={11} tickLine={false} axisLine={false} />
        <Tooltip
          contentStyle={{
            background: "var(--bg-elevated)",
            border: "1px solid var(--border)",
            borderRadius: 8,
            fontSize: 12,
          }}
          labelStyle={{ color: "var(--text-secondary)" }}
          itemStyle={{ color: "var(--foreground)" }}
        />
        <Area type="monotone" dataKey="events" stroke="var(--accent)" strokeWidth={2} fill="url(#activityFill)" />
      </AreaChart>
    </ResponsiveContainer>
  );
}

export function SeverityPieChart({ data }: { data: SeveritySlice[] }) {
  return (
    <ResponsiveContainer width="100%" height="100%">
      <PieChart>
        <Pie data={data} dataKey="value" innerRadius={40} outerRadius={70} paddingAngle={2}>
          {data.map((s, i) => <Cell key={i} fill={s.color} />)}
        </Pie>
        <Tooltip
          contentStyle={{
            background: "var(--bg-elevated)",
            border: "1px solid var(--border)",
            borderRadius: 8,
            fontSize: 12,
          }}
        />
      </PieChart>
    </ResponsiveContainer>
  );
}
`
}

// adminResourceSparklineTSX is the 30-day sparkline of the resource stat card.
func adminResourceSparklineTSX() string {
	return `"use client";

import { ResponsiveContainer, AreaChart, Area, Tooltip } from "recharts";

// Loaded by ResourceStatCard with dynamic(), so recharts is not part of the
// first load of the pages that show the card.

export interface SparklinePoint {
  date: string;
  count: number;
}

// id names the gradient, and must be unique on the page.
export function ResourceSparkline({ id, data }: { id: string; data: SparklinePoint[] }) {
  return (
    <ResponsiveContainer width="100%" height="100%">
      <AreaChart data={data} margin={{ top: 4, right: 0, bottom: 0, left: 0 }}>
        <defs>
          <linearGradient id={id} x1="0" y1="0" x2="0" y2="1">
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
          fill={"url(#" + id + ")"}
          isAnimationActive={false}
        />
      </AreaChart>
    </ResponsiveContainer>
  );
}
`
}
