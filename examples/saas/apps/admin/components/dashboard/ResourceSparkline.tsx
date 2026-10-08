"use client";

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
