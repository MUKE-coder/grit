package scaffold

// adminStatsCard returns the enhanced stats card component.
func adminStatsCard() string {
	return `"use client";

import { useEffect, useRef, useState } from "react";
import { getIcon } from "@/lib/icons";

interface StatsCardProps {
  label: string;
  value: string | number;
  change?: string;
  icon?: string;
  color?: string;
  format?: "number" | "currency" | "percentage";
  href?: string;
}

const colorMap: Record<string, string> = {
  accent: "from-accent/20 to-accent/5",
  success: "from-success/20 to-success/5",
  danger: "from-danger/20 to-danger/5",
  warning: "from-warning/20 to-warning/5",
  info: "from-info/20 to-info/5",
};

const iconColorMap: Record<string, string> = {
  accent: "text-accent",
  success: "text-success",
  danger: "text-danger",
  warning: "text-warning",
  info: "text-info",
};

export function StatsCard({ label, value, change, icon, color = "accent", format, href }: StatsCardProps) {
  const Icon = icon ? getIcon(icon) : null;
  const gradient = colorMap[color] ?? colorMap.accent;
  const iconColor = iconColorMap[color] ?? iconColorMap.accent;
  const animatedValue = useAnimatedCounter(typeof value === "number" ? value : 0);
  const displayValue = typeof value === "number" ? formatValue(animatedValue, format) : value;

  const Wrapper = href ? "a" : "div";
  const wrapperProps = href ? { href, className: "block" } : {};

  return (
    <Wrapper {...wrapperProps}>
      <div className={` + "`" + `rounded-xl border border-border bg-gradient-to-br ${gradient} p-6 transition-colors hover:border-border/80` + "`" + `}>
        <div className="flex items-center justify-between">
          {Icon && (
            <div className="rounded-lg bg-bg-secondary/50 p-2">
              <Icon className={` + "`" + `h-5 w-5 ${iconColor}` + "`" + `} />
            </div>
          )}
          {change && (
            <span className={` + "`" + `text-xs font-medium ${
              change.startsWith("+") ? "text-success" : change.startsWith("-") ? "text-danger" : "text-text-secondary"
            }` + "`" + `}>
              {change}
            </span>
          )}
        </div>
        <div className="mt-4">
          <p className="text-3xl font-bold text-foreground">{displayValue}</p>
          <p className="text-sm text-text-secondary mt-1">{label}</p>
        </div>
      </div>
    </Wrapper>
  );
}

function useAnimatedCounter(target: number, duration = 800): number {
  const [value, setValue] = useState(0);
  const startRef = useRef<number | null>(null);
  const frameRef = useRef<number>(0);

  useEffect(() => {
    if (target === 0) {
      setValue(0);
      return;
    }

    startRef.current = null;

    const animate = (timestamp: number) => {
      if (startRef.current === null) startRef.current = timestamp;
      const progress = Math.min((timestamp - startRef.current) / duration, 1);
      const eased = 1 - Math.pow(1 - progress, 3);
      setValue(Math.floor(eased * target));

      if (progress < 1) {
        frameRef.current = requestAnimationFrame(animate);
      }
    };

    frameRef.current = requestAnimationFrame(animate);
    return () => cancelAnimationFrame(frameRef.current);
  }, [target, duration]);

  return value;
}

function formatValue(value: number, format?: "number" | "currency" | "percentage"): string {
  switch (format) {
    case "currency":
      return "$" + value.toLocaleString();
    case "percentage":
      return value + "%";
    default:
      return value.toLocaleString();
  }
}
`
}

// adminChartWidget returns the Recharts chart widget component.
func adminChartWidget() string {
	return `"use client";

import dynamic from "next/dynamic";
import type { WidgetDefinition } from "@/lib/resource";

const LineChart = dynamic(
  () => import("recharts").then((mod) => {
    const { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } = mod;
    return function ChartLine({ data }: { data: ChartData[] }) {
      return (
        <ResponsiveContainer width="100%" height={300}>
          <LineChart data={data}>
            <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" />
            <XAxis dataKey="label" stroke="var(--text-muted)" fontSize={12} />
            <YAxis stroke="var(--text-muted)" fontSize={12} />
            <Tooltip
              contentStyle={{
                background: "var(--bg-elevated)",
                border: "1px solid var(--border)",
                borderRadius: "8px",
                color: "var(--text-primary)",
              }}
            />
            <Line type="monotone" dataKey="value" stroke="var(--accent)" strokeWidth={2} dot={false} />
          </LineChart>
        </ResponsiveContainer>
      );
    };
  }),
  { ssr: false, loading: () => <ChartSkeleton /> }
);

const BarChart = dynamic(
  () => import("recharts").then((mod) => {
    const { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } = mod;
    return function ChartBar({ data }: { data: ChartData[] }) {
      return (
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={data}>
            <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" />
            <XAxis dataKey="label" stroke="var(--text-muted)" fontSize={12} />
            <YAxis stroke="var(--text-muted)" fontSize={12} />
            <Tooltip
              contentStyle={{
                background: "var(--bg-elevated)",
                border: "1px solid var(--border)",
                borderRadius: "8px",
                color: "var(--text-primary)",
              }}
            />
            <Bar dataKey="value" fill="var(--accent)" radius={[4, 4, 0, 0]} />
          </BarChart>
        </ResponsiveContainer>
      );
    };
  }),
  { ssr: false, loading: () => <ChartSkeleton /> }
);

const PieChartComponent = dynamic(
  () => import("recharts").then((mod) => {
    const { PieChart, Pie, Cell, Tooltip, ResponsiveContainer } = mod;
    const COLORS = ["var(--accent)", "var(--success)", "var(--warning)", "var(--info)", "var(--danger)"];
    return function ChartPie({ data }: { data: ChartData[] }) {
      return (
        <ResponsiveContainer width="100%" height={300}>
          <PieChart>
            <Pie data={data} cx="50%" cy="50%" innerRadius={60} outerRadius={100} dataKey="value" nameKey="label">
              {data.map((_, i) => (
                <Cell key={i} fill={COLORS[i % COLORS.length]} />
              ))}
            </Pie>
            <Tooltip
              contentStyle={{
                background: "var(--bg-elevated)",
                border: "1px solid var(--border)",
                borderRadius: "8px",
                color: "var(--text-primary)",
              }}
            />
          </PieChart>
        </ResponsiveContainer>
      );
    };
  }),
  { ssr: false, loading: () => <ChartSkeleton /> }
);

interface ChartData {
  label: string;
  value: number;
}

interface ChartWidgetProps {
  config: WidgetDefinition;
  data?: ChartData[];
}

export function ChartWidget({ config, data = [] }: ChartWidgetProps) {
  return (
    <div className="rounded-xl border border-border bg-bg-secondary p-6">
      <h3 className="text-sm font-medium text-text-secondary mb-4">{config.label}</h3>
      {data.length === 0 ? (
        <div className="flex items-center justify-center h-[300px] text-text-muted text-sm">
          No chart data available
        </div>
      ) : config.chartType === "bar" ? (
        <BarChart data={data} />
      ) : config.chartType === "pie" ? (
        <PieChartComponent data={data} />
      ) : (
        <LineChart data={data} />
      )}
    </div>
  );
}

function ChartSkeleton() {
  return (
    <div className="flex items-center justify-center h-[300px]">
      <div className="h-6 w-6 animate-spin rounded-full border-2 border-accent border-t-transparent" />
    </div>
  );
}
`
}

// adminActivityWidget returns the activity feed widget.
func adminActivityWidget() string {
	return `interface ActivityItem {
  id: string;
  message: string;
  timestamp: string;
  user?: string;
}

interface ActivityWidgetProps {
  label?: string;
  items?: ActivityItem[];
}

export function ActivityWidget({ label = "Recent Activity", items = [] }: ActivityWidgetProps) {
  return (
    <div className="rounded-xl border border-border bg-bg-secondary p-6">
      <h3 className="text-sm font-medium text-text-secondary mb-4">{label}</h3>
      {items.length === 0 ? (
        <div className="space-y-4">
          {[1, 2, 3, 4, 5].map((i) => (
            <div key={i} className="flex items-center gap-3 text-sm">
              <div className="h-2 w-2 rounded-full bg-accent shrink-0" />
              <span className="text-text-secondary">Activity placeholder #{i}</span>
              <span className="ml-auto text-text-muted text-xs shrink-0">Just now</span>
            </div>
          ))}
        </div>
      ) : (
        <div className="space-y-4">
          {items.map((item) => (
            <div key={item.id} className="flex items-center gap-3 text-sm">
              <div className="h-2 w-2 rounded-full bg-accent shrink-0" />
              <span className="text-text-secondary flex-1 truncate">{item.message}</span>
              <span className="ml-auto text-text-muted text-xs shrink-0">{item.timestamp}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
`
}

// adminWidgetGrid returns the responsive widget grid component.
func adminWidgetGrid() string {
	return `"use client";

import type { WidgetDefinition } from "@/lib/resource";
import { StatsCard } from "./stats-card";
import { ChartWidget } from "./chart-widget";
import { ActivityWidget } from "./activity-widget";
import { useResource } from "@/hooks/use-resource";

interface WidgetGridProps {
  widgets: WidgetDefinition[];
}

export function WidgetGrid({ widgets }: WidgetGridProps) {
  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {widgets.map((widget, i) => (
        <div
          key={i}
          className={` + "`" + `${
            widget.colSpan === 2 ? "sm:col-span-2" :
            widget.colSpan === 3 ? "sm:col-span-2 lg:col-span-3" :
            widget.colSpan === 4 ? "sm:col-span-2 lg:col-span-4" :
            ""
          }` + "`" + `}
        >
          <WidgetRenderer widget={widget} />
        </div>
      ))}
    </div>
  );
}

function WidgetRenderer({ widget }: { widget: WidgetDefinition }) {
  switch (widget.type) {
    case "stat":
      return <StatWidgetLoader widget={widget} />;
    case "chart":
      return <ChartWidget config={widget} />;
    case "activity":
      return <ActivityWidget label={widget.label} />;
    default:
      return null;
  }
}

function StatWidgetLoader({ widget }: { widget: WidgetDefinition }) {
  const { data } = useResource(widget.endpoint ?? "", {
    page: 1,
    pageSize: 1,
  });

  const total = data?.meta?.total ?? 0;

  return (
    <StatsCard
      label={widget.label}
      value={total}
      icon={widget.icon}
      color={widget.color}
      format={widget.format}
    />
  );
}
`
}
