"use client";

import { useQuery } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { SkeletonCards } from "@/components/ui/Skeleton";
import { apiClient } from "@/lib/api-client";
import {
  TrendingUp, AlertCircle, AlertTriangle, ExternalLink, Activity as ActivityIcon,
  Cpu, Database, Gauge,
} from "@/lib/icons";

// The Pulse UI is mounted on the Go API, not on this admin host. Use the
// API base so "Open Pulse" works whether the admin is on :3001 in dev or
// a different origin in prod.
import { API_URL } from "@/lib/api-core";

// DegradedBanner says which panels have no data behind them.
//
// The dashboard used to answer a dead upstream with a page of zeros and a 200,
// which reads as good news. The API now names the calls that failed, and this
// puts that on the screen above the numbers it applies to.
function DegradedBanner({ source, degraded }: { source: string; degraded?: string[] }) {
  if (!degraded || degraded.length === 0) return null;
  return (
    <div className="mt-4 flex items-start gap-2 rounded-lg border border-warning/30 bg-warning/5 px-4 py-3 text-sm text-warning">
      <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0" />
      <span>
        {source} did not answer for {degraded.join(", ")}. Those panels are empty rather than zero.
      </span>
    </div>
  );
}

interface PerformanceSummary {
  latency?: { p50: number; p95: number; p99: number; avg: number };
  traffic?: { throughput: number; total: number };
  errors?:  { rate: number; active_open: number };
  saturation?: { goroutines: number; heap_mb: number; gc_cycles: number; cpu_cores: number };
  slowest_routes?: Array<{ route: string; method: string; requests: number; avg: number; p95: number; p99: number; error_rate: number }>;
  n1_detections?: Array<{ route: string; query_count: number; first_seen: string }>;
  recent_errors?: Array<{ id: string; route: string; message: string; created_at: string }>;
  // The calls to Pulse that did not answer. Empty is the healthy case.
  degraded?: string[];
}

export default function PerformancePage() {
  const { data, isLoading } = useQuery<PerformanceSummary>({
    queryKey: ["performance", "summary"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<PerformanceSummary>("/api/admin/observability/summary");
        return data;
      } catch {
        // Nothing came back. Every signal would render as a dash, which looks
        // like an idle app rather than a dashboard that cannot see one.
        return { degraded: ["the whole summary"] };
      }
    },
    refetchInterval: 30_000,
  });

  return (
    <div>
      <PageHeader
        title="Performance"
        subtitle="Four SRE golden signals + route, query, and error detail — powered by Pulse."
        actions={
          <a
            href={`${API_URL}/pulse/ui`}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex shrink-0 items-center gap-2 whitespace-nowrap rounded-lg border border-border bg-bg-elevated px-3 py-2 text-sm font-medium text-foreground hover:bg-bg-hover"
          >
            <ExternalLink className="h-4 w-4" />
            Open Pulse
          </a>
        }
      />

      <DegradedBanner source="Pulse" degraded={data?.degraded} />

      {isLoading ? <SkeletonCards count={4} /> : (
        <>
          {/* Latency */}
          <SignalGroup
            title="Latency"
            tagline="How long requests take. Watch p95/p99 — the average hides the long tail."
          >
            <Signal label="P50" value={fmt(data?.latency?.p50, "ms")} icon={<Gauge className="h-4 w-4" />} />
            <Signal label="P95" value={fmt(data?.latency?.p95, "ms")} icon={<Gauge className="h-4 w-4" />} />
            <Signal label="P99" value={fmt(data?.latency?.p99, "ms")} icon={<Gauge className="h-4 w-4" />} tone="warning" />
            <Signal label="AVG" value={fmt(data?.latency?.avg, "ms")} icon={<Gauge className="h-4 w-4" />} />
          </SignalGroup>

          {/* Traffic */}
          <SignalGroup
            title="Traffic"
            tagline="How much demand the API is handling right now."
          >
            <Signal label="Throughput" value={fmtRate(data?.traffic?.throughput)} icon={<TrendingUp className="h-4 w-4 text-success" />} />
            <Signal label="Total requests" value={fmt(data?.traffic?.total)} icon={<TrendingUp className="h-4 w-4 text-success" />} />
          </SignalGroup>

          {/* Errors */}
          <SignalGroup
            title="Errors"
            tagline="Rate of failures. Spikes here usually correlate with latency spikes — check both."
          >
            <Signal label="Error rate" value={fmtPct(data?.errors?.rate)} icon={<AlertCircle className="h-4 w-4 text-danger" />} tone={(data?.errors?.rate ?? 0) > 0 ? "danger" : "default"} />
            <Signal label="Active errors (open)" value={fmt(data?.errors?.active_open)} icon={<AlertCircle className="h-4 w-4 text-danger" />} tone={(data?.errors?.active_open ?? 0) > 0 ? "danger" : "default"} />
          </SignalGroup>

          {/* Saturation */}
          <SignalGroup
            title="Saturation"
            tagline="How full your resources are. Red bands here mean a bottleneck is imminent or already firing — fix before users feel it."
          >
            <Signal label="Goroutines" value={fmt(data?.saturation?.goroutines)} icon={<ActivityIcon className="h-4 w-4" />} />
            <Signal label="Heap alloc" value={fmt(data?.saturation?.heap_mb, "MB")} icon={<Database className="h-4 w-4" />} />
            <Signal label="GC cycles" value={fmt(data?.saturation?.gc_cycles)} icon={<Database className="h-4 w-4" />} />
            <Signal label="CPU cores" value={fmt(data?.saturation?.cpu_cores)} icon={<Cpu className="h-4 w-4" />} />
          </SignalGroup>
        </>
      )}

      {/* Slowest routes */}
      <section className="mt-6 overflow-hidden rounded-xl border border-border bg-bg-elevated">
        <header className="border-b border-border px-5 py-3">
          <p className="text-xs font-semibold uppercase tracking-wider text-text-muted">Slowest routes</p>
        </header>
        {(data?.slowest_routes?.length ?? 0) === 0 ? (
          <p className="px-5 py-8 text-center text-sm text-text-muted">No route latency data yet.</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full table-fixed">
              <thead>
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-semibold uppercase text-text-muted">Route</th>
                  <th className="w-20 px-2 py-3 text-right text-xs font-semibold uppercase text-text-muted">Reqs</th>
                  <th className="w-20 px-2 py-3 text-right text-xs font-semibold uppercase text-text-muted">Avg</th>
                  <th className="w-20 px-2 py-3 text-right text-xs font-semibold uppercase text-text-muted">P95</th>
                  <th className="w-20 px-2 py-3 text-right text-xs font-semibold uppercase text-text-muted">P99</th>
                  <th className="w-24 px-2 py-3 text-right text-xs font-semibold uppercase text-text-muted">Err rate</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {data!.slowest_routes!.map((r, i) => (
                  <tr key={i}>
                    <td className="truncate px-4 py-3 font-mono text-xs text-foreground">{r.method} {r.route}</td>
                    <td className="px-2 py-3 text-right text-xs text-foreground">{r.requests}</td>
                    <td className="px-2 py-3 text-right text-xs text-foreground">{fmt(r.avg, "ms")}</td>
                    <td className="px-2 py-3 text-right text-xs text-foreground">{fmt(r.p95, "ms")}</td>
                    <td className="px-2 py-3 text-right text-xs text-foreground">{fmt(r.p99, "ms")}</td>
                    <td className="px-2 py-3 text-right text-xs text-foreground">{fmtPct(r.error_rate)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {/* N+1 detections */}
      <section className="mt-6 overflow-hidden rounded-xl border border-border bg-bg-elevated">
        <header className="border-b border-border px-5 py-3">
          <p className="text-xs font-semibold uppercase tracking-wider text-text-muted">N+1 query detections</p>
        </header>
        {(data?.n1_detections?.length ?? 0) === 0 ? (
          <p className="px-5 py-8 text-center text-sm text-text-muted">No N+1 queries detected.</p>
        ) : (
          <ul className="divide-y divide-border">
            {data!.n1_detections!.map((n, i) => (
              <li key={i} className="flex items-center justify-between gap-3 px-5 py-3">
                <p className="truncate font-mono text-xs text-foreground">{n.route}</p>
                <span className="shrink-0 text-xs text-text-muted">
                  {n.query_count} queries · first seen {new Date(n.first_seen).toLocaleString()}
                </span>
              </li>
            ))}
          </ul>
        )}
      </section>

      {/* Recent errors */}
      <section className="mt-6 overflow-hidden rounded-xl border border-border bg-bg-elevated">
        <header className="border-b border-border px-5 py-3">
          <p className="text-xs font-semibold uppercase tracking-wider text-text-muted">Recent errors</p>
        </header>
        {(data?.recent_errors?.length ?? 0) === 0 ? (
          <p className="px-5 py-8 text-center text-sm text-text-muted">No errors recorded.</p>
        ) : (
          <ul className="divide-y divide-border">
            {data!.recent_errors!.map((e) => (
              <li key={e.id} className="px-5 py-3">
                <div className="flex items-center justify-between">
                  <p className="truncate font-mono text-xs text-foreground">{e.route}</p>
                  <p className="shrink-0 text-xs text-text-muted">{new Date(e.created_at).toLocaleString()}</p>
                </div>
                <p className="mt-1 text-sm text-danger">{e.message}</p>
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}

interface SignalGroupProps { title: string; tagline: string; children: React.ReactNode }
function SignalGroup({ title, tagline, children }: SignalGroupProps) {
  return (
    <section className="mt-6">
      <header className="mb-3">
        <p className="text-sm font-semibold uppercase tracking-wider text-text-muted">{title}</p>
        <p className="text-xs text-text-secondary">{tagline}</p>
      </header>
      <div className="grid grid-cols-2 gap-3 md:grid-cols-4">{children}</div>
    </section>
  );
}

interface SignalProps { label: string; value: string; icon: React.ReactNode; tone?: "default" | "warning" | "danger" }
function Signal({ label, value, icon, tone = "default" }: SignalProps) {
  const toneClass = { default: "border-border bg-bg-elevated", warning: "border-warning/30 bg-warning/5", danger: "border-danger/30 bg-danger/5" }[tone];
  return (
    <div className={"rounded-xl border p-4 " + toneClass}>
      <div className="flex items-center justify-between">
        <p className="text-xs font-semibold uppercase tracking-wide text-text-muted">{label}</p>
        <span className="text-text-secondary">{icon}</span>
      </div>
      <p className="mt-2 text-2xl font-bold text-foreground">{value}</p>
    </div>
  );
}

function fmt(n: number | undefined, suffix: string = ""): string {
  if (n === undefined || n === null || Number.isNaN(n)) return "—";
  return Math.round(n).toLocaleString() + suffix;
}

// fmtRate formats a per-second rate. Rounding to an integer would collapse a
// low-but-real throughput (e.g. 0.14 req/s) to "0/s", so keep two decimals
// below 10 and switch to whole numbers once the rate is large.
function fmtRate(n: number | undefined): string {
  if (n === undefined || n === null || Number.isNaN(n)) return "—";
  if (n > 0 && n < 10) return n.toFixed(2) + "/s";
  return Math.round(n).toLocaleString() + "/s";
}

// Both callers pass a Pulse error rate, and Pulse already reports those as a
// percentage (errCount / total * 100 — see storage_memory.go). Multiplying by
// 100 here scaled every value by 100 again: 3 errors in 38 requests rendered as
// "789.47%". Takes a percentage, prints a percentage.
function fmtPct(n: number | undefined): string {
  if (n === undefined || n === null || Number.isNaN(n)) return "—";
  return n.toFixed(2) + "%";
}
