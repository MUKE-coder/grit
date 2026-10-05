"use client";

import type { ReactNode } from "react";
import type { LucideIcon } from "lucide-react";
import { Activity, ExternalLink, AlertTriangle, Zap, CheckCircle, Circle, CircleDashed, Minus, Gauge } from "@/lib/icons";
import { useObservabilitySummary, useScaleReport } from "@/hooks/use-system";
import { getApiErrorMessage } from "@/lib/api-core";

import { API_URL } from "@/lib/api-core";

interface Summary {
  overview?: { p50_ms?: number; p95_ms?: number; p99_ms?: number; rps?: number; error_rate?: number; total_requests?: number }
  slos?: { data?: Array<{ name: string; target: number; current: number; budget_remaining: number; status: string }> }
  use?: { resources?: Array<{ name: string; utilization: { value: number; band: string }; saturation: { value: number; band: string }; errors: { value: number; band: string } }> }
  n1_ranked?: { data?: Array<{ route: string; pattern: string; occurrences: number; avg_queries_per_request: number; impact_score: number }> }
  errors?: { data?: Array<{ id: string; type: string; message: string; route: string; count: number }> }
  runtime?: { heap_alloc_mb?: number; goroutines?: number; gc_pause_ms?: number }
  health_checks?: { data?: Array<{ name: string; status: string; latency_ms: number; error?: string }> }
  alerts?: { data?: Array<{ id: string; name: string; severity: string; message: string }> }
  _errors?: Record<string, string>
}

type StageState = "on" | "ready" | "attention" | "yours";

interface ScaleReport {
  stage?: string
  next?: string
  why?: string
  notes?: string[]
  stages?: Array<{ n: string; name: string; state: StageState; what: string; evidence: string }>
  latency?: { samples?: number; p50_ms?: number; p95_ms?: number; p99_ms?: number }
  database?: { replicas?: number; open_connections?: number; max_connections?: number }
  cache?: { configured?: boolean; hit_rate?: number }
}

// How each state reads. "ready" is a tick too, hollow rather than filled: the
// stage is shipped and switched off because nothing yet needs it, which is a
// good answer and not a missing feature.
const STAGE_STATE: Record<StageState, { icon: typeof CheckCircle; cls: string; label: string }> = {
  on:        { icon: CheckCircle,  cls: "text-success",    label: "In force" },
  ready:     { icon: Circle,       cls: "text-success/60", label: "Ready, not needed yet" },
  attention: { icon: AlertTriangle, cls: "text-warning",   label: "Needs attention" },
  yours:     { icon: Minus,        cls: "text-text-muted", label: "Yours to decide" },
};

const BAND_CLS: Record<string, string> = {
  green:   "bg-success/15 text-success border-success/30",
  amber:   "bg-warning/15 text-warning border-warning/30",
  red:     "bg-danger/15 text-danger border-danger/30",
  unknown: "bg-bg-elevated text-text-muted border-border",
};

export default function ObservabilityPage() {
  // Polls every 10 seconds. The local helper this replaced was called "fetch",
  // which shadowed the global one for the rest of the effect: any code added
  // there that reached for fetch() would have called the poller instead.
  const query = useObservabilitySummary<Summary>();
  const data = query.data ?? null;
  const loading = query.isPending;
  const err = query.error ? getApiErrorMessage(query.error, "Failed to load") : null;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground flex items-center gap-2.5">
            <Activity className="h-6 w-6 text-accent" /> Observability
          </h1>
          <p className="text-sm text-text-secondary mt-1">
            Live summary from Pulse — percentile latency, SLOs, USE grid, top N+1, errors, runtime
          </p>
        </div>
        <a
          href={`${API_URL}/pulse/ui`}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-2 rounded-lg border border-border bg-bg-secondary px-4 py-2.5 text-sm font-medium text-foreground hover:bg-bg-hover transition-colors"
        >
          Open full dashboard <ExternalLink className="h-4 w-4" />
        </a>
      </div>

      {err && !loading && (
        <div className="rounded-xl border border-warning/30 bg-warning/5 p-4 text-sm text-warning flex items-start gap-2">
          <AlertTriangle className="h-4 w-4 mt-0.5 shrink-0" />
          <div>
            <strong>Couldn&apos;t reach Pulse.</strong> Make sure <code className="px-1 py-0.5 rounded bg-bg-secondary text-xs font-mono">PULSE_ENABLED=true</code> and the API is running.
            <div className="text-xs text-text-muted mt-1">{err}</div>
          </div>
        </div>
      )}

      <ScalingReadiness />

      {/* Top-line KPIs */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        <Kpi label="p95 latency" value={data?.overview?.p95_ms != null ? `${data.overview.p95_ms.toFixed(0)} ms` : "—"} tone="info" icon={Zap} />
        <Kpi label="p99 latency" value={data?.overview?.p99_ms != null ? `${data.overview.p99_ms.toFixed(0)} ms` : "—"} tone="warning" icon={Zap} />
        <Kpi label="Error rate" value={data?.overview?.error_rate != null ? `${(data.overview.error_rate * 100).toFixed(2)}%` : "—"} tone="danger" icon={AlertTriangle} />
        <Kpi label="RPS" value={data?.overview?.rps?.toFixed?.(1) ?? "—"} tone="success" icon={Activity} />
      </div>

      {/* SLO compliance bars */}
      <Panel title="SLOs">
        {(data?.slos?.data ?? []).length === 0 ? (
          <p className="text-sm text-text-muted">No SLOs configured. Add to <code className="px-1 py-0.5 rounded bg-bg-elevated text-xs font-mono">pulse.Config.SLOs</code>.</p>
        ) : (
          <div className="space-y-3">
            {(data?.slos?.data ?? []).map((s) => (
              <div key={s.name} className="space-y-1">
                <div className="flex items-center justify-between text-xs">
                  <span className="font-medium text-foreground">{s.name}</span>
                  <span className={`font-mono ${s.status === "firing" ? "text-danger" : "text-success"}`}>{(s.current * 100).toFixed(2)}% / {(s.target * 100).toFixed(2)}%</span>
                </div>
                <div className="h-2 rounded-full bg-bg-elevated overflow-hidden">
                  <div className={`h-full ${s.status === "firing" ? "bg-danger" : "bg-success"}`} style={{ width: `${Math.min(s.current * 100, 100)}%` }} />
                </div>
                <p className="text-[10px] text-text-muted">Budget remaining: {(s.budget_remaining * 100).toFixed(1)}%</p>
              </div>
            ))}
          </div>
        )}
      </Panel>

      {/* USE grid + N+1 side-by-side */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <Panel title="USE method">
          <table className="w-full text-xs">
            <thead>
              <tr className="text-text-muted">
                <th className="text-left font-medium py-1">Resource</th>
                <th className="text-center font-medium">U</th>
                <th className="text-center font-medium">S</th>
                <th className="text-center font-medium">E</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {(data?.use?.resources ?? []).map((r) => (
                <tr key={r.name}>
                  <td className="py-1.5 font-medium text-foreground capitalize">{r.name}</td>
                  <td className="py-1.5 text-center"><Cell band={r.utilization?.band} value={r.utilization?.value} /></td>
                  <td className="py-1.5 text-center"><Cell band={r.saturation?.band} value={r.saturation?.value} /></td>
                  <td className="py-1.5 text-center"><Cell band={r.errors?.band} value={r.errors?.value} /></td>
                </tr>
              ))}
              {(!data?.use?.resources || data.use.resources.length === 0) && (
                <tr><td colSpan={4} className="text-center text-text-muted py-3">No USE samples yet.</td></tr>
              )}
            </tbody>
          </table>
        </Panel>

        <Panel title="Top N+1 by impact">
          <ul className="space-y-2 text-xs">
            {(data?.n1_ranked?.data ?? []).slice(0, 6).map((n, i) => (
              <li key={i} className="rounded-lg border border-border bg-bg-elevated p-2.5">
                <div className="flex items-center justify-between mb-0.5">
                  <span className="font-mono text-foreground truncate">{n.route}</span>
                  <span className="text-text-muted shrink-0 ml-2">impact {n.impact_score.toFixed(0)}</span>
                </div>
                <p className="font-mono text-[10px] text-text-muted truncate">{n.pattern}</p>
                <p className="text-[10px] text-text-muted mt-0.5">{n.occurrences} occurrences · ~{n.avg_queries_per_request} queries/req</p>
              </li>
            ))}
            {(!data?.n1_ranked?.data || data.n1_ranked.data.length === 0) && (
              <li className="text-text-muted text-center py-3">No N+1 detections.</li>
            )}
          </ul>
        </Panel>
      </div>

      {/* Errors + Runtime side-by-side */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <Panel title="Recent unresolved errors">
          <ul className="space-y-1.5 text-xs">
            {(data?.errors?.data ?? []).slice(0, 6).map((e) => (
              <li key={e.id} className="flex items-start justify-between gap-2 py-1 border-b border-border last:border-0">
                <div className="min-w-0 flex-1">
                  <p className="text-foreground truncate"><span className="font-mono text-[10px] text-text-muted">{e.type}</span> {e.message}</p>
                  <p className="font-mono text-[10px] text-text-muted truncate">{e.route}</p>
                </div>
                <span className="text-text-muted shrink-0">{e.count}×</span>
              </li>
            ))}
            {(!data?.errors?.data || data.errors.data.length === 0) && <li className="text-text-muted text-center py-3">No unresolved errors.</li>}
          </ul>
        </Panel>

        <Panel title="Go runtime">
          <div className="grid grid-cols-3 gap-2 text-xs">
            <Stat label="Heap alloc" value={data?.runtime?.heap_alloc_mb != null ? `${data.runtime.heap_alloc_mb.toFixed(0)} MB` : "—"} />
            <Stat label="Goroutines" value={data?.runtime?.goroutines ?? "—"} />
            <Stat label="GC pause" value={data?.runtime?.gc_pause_ms != null ? `${data.runtime.gc_pause_ms.toFixed(1)} ms` : "—"} />
          </div>
          {data?.health_checks?.data && data.health_checks.data.length > 0 && (
            <div className="mt-4 pt-3 border-t border-border space-y-1">
              <p className="text-[10px] uppercase font-mono tracking-wider text-text-muted">Health checks</p>
              {data.health_checks.data.map((h) => (
                <div key={h.name} className="flex items-center justify-between text-xs">
                  <span className="font-mono text-foreground">{h.name}</span>
                  <span className={`text-[10px] ${h.status === "healthy" ? "text-success" : "text-danger"}`}>
                    {h.status} · {h.latency_ms}ms
                  </span>
                </div>
              ))}
            </div>
          )}
        </Panel>
      </div>
    </div>
  );
}

// Which of the ten scaling stages this deployment is past, measured rather than
// claimed. The verdict names the one thing to do next, and most of the time
// that is nothing; the list below it is the half people ask about first, which
// is what is already handled.
function ScalingReadiness() {
  const query = useScaleReport<ScaleReport>();
  const report = query.data ?? null;
  const stages = report?.stages ?? [];
  const attention = stages.filter((s) => s.state === "attention").length;
  const covered = stages.filter((s) => s.state === "on" || s.state === "ready").length;

  if (query.isError) {
    return (
      <Panel title="Scaling readiness">
        <p className="text-sm text-text-muted">
          {getApiErrorMessage(query.error, "Could not read the scaling report")}
        </p>
      </Panel>
    );
  }

  return (
    <div className="rounded-xl border border-border bg-bg-secondary overflow-hidden">
      <div className="px-4 py-3 border-b border-border flex items-center justify-between gap-3">
        <h2 className="text-sm font-semibold text-foreground flex items-center gap-2">
          <Gauge className="h-4 w-4 text-accent" /> Scaling readiness
        </h2>
        {stages.length > 0 && (
          <span className={`${attention > 0 ? "text-warning" : "text-success"} text-xs font-mono`}>
            {covered} of {stages.length} stages handled{attention > 0 ? `, ${attention} need${attention === 1 ? "s" : ""} attention` : ""}
          </span>
        )}
      </div>

      {/* The verdict, first, because it is the answer to "what now". */}
      <div className="px-4 py-3 border-b border-border bg-bg-elevated/40">
        {query.isPending ? (
          <p className="text-sm text-text-muted">Measuring…</p>
        ) : (
          <>
            <p className="text-sm font-medium text-foreground">{report?.stage ?? "—"}</p>
            {report?.why && <p className="text-xs text-text-secondary mt-0.5">{report.why}</p>}
            {report?.next && (
              <p className="text-xs text-text-muted mt-1.5">
                <span className="font-mono uppercase tracking-wider text-[10px] text-text-muted">Do this next</span>{" "}
                <span className="text-foreground">{report.next}</span>
              </p>
            )}
          </>
        )}
      </div>

      <ul className="divide-y divide-border">
        {stages.map((stage) => {
          const state = STAGE_STATE[stage.state] ?? STAGE_STATE.yours;
          const Icon = state.icon;
          return (
            <li key={stage.n} className="flex items-start gap-3 px-4 py-2.5">
              <Icon className={`h-4 w-4 mt-0.5 shrink-0 ${state.cls}`} aria-hidden />
              <div className="min-w-0 flex-1">
                <div className="flex items-baseline gap-2">
                  <span className="font-mono text-[10px] text-text-muted shrink-0">{stage.n}</span>
                  <span className="text-sm font-medium text-foreground">{stage.name}</span>
                </div>
                <p className="text-xs text-text-secondary mt-0.5">{stage.what}</p>
                {stage.evidence && (
                  <p className="text-[11px] text-text-muted mt-0.5 font-mono">{stage.evidence}</p>
                )}
              </div>
              <span className={`text-[10px] shrink-0 mt-1 ${state.cls}`}>{state.label}</span>
            </li>
          );
        })}
        {!query.isPending && stages.length === 0 && (
          <li className="px-4 py-6 text-center text-sm text-text-muted">
            No stage report. This needs an API on v3.336.0 or later.
          </li>
        )}
      </ul>

      {(report?.notes ?? []).length > 0 && (
        <div className="px-4 py-3 border-t border-border space-y-1">
          {(report?.notes ?? []).map((note, i) => (
            <p key={i} className="text-[11px] text-text-muted flex items-start gap-1.5">
              <CircleDashed className="h-3 w-3 mt-0.5 shrink-0" aria-hidden />
              {note}
            </p>
          ))}
        </div>
      )}
    </div>
  );
}

function Kpi({ label, value, tone, icon: Icon }: { label: string; value: ReactNode; tone: "default" | "success" | "warning" | "info" | "danger"; icon: LucideIcon }) {
  const toneCls = { default: "text-foreground", success: "text-success", warning: "text-warning", info: "text-info", danger: "text-danger" }[tone];
  return (
    <div className="rounded-xl border border-border bg-bg-secondary p-4">
      <div className="flex items-center justify-between mb-1">
        <span className="text-[10px] font-mono uppercase tracking-wider text-text-muted">{label}</span>
        <Icon className={`h-4 w-4 ${toneCls}`} />
      </div>
      <p className={`text-2xl font-semibold tabular-nums ${toneCls}`}>{value}</p>
    </div>
  );
}

function Panel({ title, children }: { title: string; children: ReactNode }) {
  return (
    <div className="rounded-xl border border-border bg-bg-secondary overflow-hidden">
      <div className="px-4 py-3 border-b border-border"><h2 className="text-sm font-semibold text-foreground">{title}</h2></div>
      <div className="p-4">{children}</div>
    </div>
  );
}

function Stat({ label, value }: { label: string; value: ReactNode }) {
  return (
    <div className="rounded-lg border border-border bg-bg-elevated px-3 py-2">
      <p className="text-[10px] uppercase tracking-wider text-text-muted">{label}</p>
      <p className="text-base font-semibold text-foreground tabular-nums">{value}</p>
    </div>
  );
}

function Cell({ band, value }: { band?: string; value?: number }) {
  const cls = BAND_CLS[band || "unknown"];
  return (
    <span className={`inline-flex items-center justify-center min-w-[44px] px-2 py-0.5 rounded border text-[10px] font-mono ${cls}`}>
      {value != null ? value.toFixed(0) : "—"}
    </span>
  );
}
