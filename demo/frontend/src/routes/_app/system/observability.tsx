import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { Activity, ExternalLink, AlertTriangle, Cpu, Zap } from 'lucide-react'
import api from '@/lib/api'

const API_URL = import.meta.env.VITE_API_URL || ''

export const Route = createFileRoute('/_app/system/observability')({
  component: ObservabilityPage,
})

interface Summary {
  overview?: { p50_ms?: number; p95_ms?: number; p99_ms?: number; rps?: number; error_rate?: number; total_requests?: number }
  slos?: { data?: Array<{ name: string; target: number; current: number; budget_remaining: number; status: string }> }
  use?: { resources?: Array<{ name: string; utilization: { value: number; band: string }; saturation: { value: number; band: string }; errors: { value: number; band: string } }> }
  n1_ranked?: { data?: Array<{ route: string; pattern: string; occurrences: number; avg_queries_per_request: number; impact_score: number }> }
  errors?: { data?: Array<{ id: string; type: string; message: string; route: string; count: number }> }
  runtime?: { heap_alloc_mb?: number; goroutines?: number; gc_pause_ms?: number }
  health_checks?: { data?: Array<{ name: string; status: string; latency_ms: number }> }
  _errors?: Record<string, string>
}

const BAND_CLS: Record<string, string> = {
  green:   'bg-emerald-50 text-emerald-700 border-emerald-200',
  amber:   'bg-amber-50 text-amber-700 border-amber-200',
  red:     'bg-rose-50 text-rose-700 border-rose-200',
  unknown: 'bg-slate-50 text-slate-500 border-slate-200',
}

function ObservabilityPage() {
  const { data, error, isLoading } = useQuery<Summary>({
    queryKey: ['admin', 'observability', 'summary'],
    queryFn: async () => {
      const res = await api.get('/api/admin/observability/summary')
      return (res.data?.data ?? res.data) as Summary
    },
    refetchInterval: 10_000,
  })

  return (
    <div className="space-y-6 p-6 max-w-6xl mx-auto">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground flex items-center gap-2.5">
            <Activity className="h-6 w-6 text-grit-blue" /> Observability
          </h1>
          <p className="text-sm text-foreground-muted mt-1">
            Live summary from Pulse — percentile latency, SLOs, USE grid, top N+1, errors, runtime.
          </p>
        </div>
        <a
          href={`${API_URL}/pulse/ui`}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-2 rounded-lg border border-border bg-surface px-4 py-2.5 text-sm font-medium hover:bg-surface-hover transition-colors"
        >
          Open full dashboard <ExternalLink className="h-4 w-4" />
        </a>
      </div>

      {error && !isLoading && (
        <div className="rounded-xl border border-amber-300 bg-amber-50 p-4 text-sm text-amber-900 flex items-start gap-2">
          <AlertTriangle className="h-4 w-4 mt-0.5 shrink-0" />
          <div>
            <strong>Couldn&apos;t reach Pulse.</strong> Make sure <code className="px-1 py-0.5 rounded bg-amber-100 text-xs font-mono">PULSE_ENABLED=true</code> and the API is running.
          </div>
        </div>
      )}

      {/* Top-line KPIs */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        <Kpi label="p95 latency" value={data?.overview?.p95_ms != null ? `${data.overview.p95_ms.toFixed(0)} ms` : '—'} tone="sky" icon={Zap} />
        <Kpi label="p99 latency" value={data?.overview?.p99_ms != null ? `${data.overview.p99_ms.toFixed(0)} ms` : '—'} tone="amber" icon={Zap} />
        <Kpi label="Error rate" value={data?.overview?.error_rate != null ? `${(data.overview.error_rate * 100).toFixed(2)}%` : '—'} tone="rose" icon={AlertTriangle} />
        <Kpi label="RPS" value={data?.overview?.rps?.toFixed?.(1) ?? '—'} tone="emerald" icon={Activity} />
      </div>

      <Panel title="SLOs">
        {(data?.slos?.data ?? []).length === 0 ? (
          <p className="text-sm text-foreground-muted">No SLOs configured.</p>
        ) : (
          <div className="space-y-3">
            {(data?.slos?.data ?? []).map((s) => (
              <div key={s.name} className="space-y-1">
                <div className="flex items-center justify-between text-xs">
                  <span className="font-medium">{s.name}</span>
                  <span className={`font-mono ${s.status === 'firing' ? 'text-rose-700' : 'text-emerald-700'}`}>
                    {(s.current * 100).toFixed(2)}% / {(s.target * 100).toFixed(2)}%
                  </span>
                </div>
                <div className="h-2 rounded-full bg-bg overflow-hidden">
                  <div className={`h-full ${s.status === 'firing' ? 'bg-rose-500' : 'bg-emerald-500'}`} style={{ width: `${Math.min(s.current * 100, 100)}%` }} />
                </div>
                <p className="text-[10px] text-foreground-muted">Budget remaining: {(s.budget_remaining * 100).toFixed(1)}%</p>
              </div>
            ))}
          </div>
        )}
      </Panel>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <Panel title="USE method">
          <table className="w-full text-xs">
            <thead>
              <tr className="text-foreground-muted">
                <th className="text-left font-medium py-1">Resource</th>
                <th className="text-center font-medium">U</th>
                <th className="text-center font-medium">S</th>
                <th className="text-center font-medium">E</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {(data?.use?.resources ?? []).map((r) => (
                <tr key={r.name}>
                  <td className="py-1.5 font-medium capitalize">{r.name}</td>
                  <td className="py-1.5 text-center"><Cell band={r.utilization?.band} value={r.utilization?.value} /></td>
                  <td className="py-1.5 text-center"><Cell band={r.saturation?.band} value={r.saturation?.value} /></td>
                  <td className="py-1.5 text-center"><Cell band={r.errors?.band} value={r.errors?.value} /></td>
                </tr>
              ))}
              {(!data?.use?.resources || data.use.resources.length === 0) && (
                <tr><td colSpan={4} className="text-center text-foreground-muted py-3">No USE samples yet.</td></tr>
              )}
            </tbody>
          </table>
        </Panel>

        <Panel title="Top N+1 by impact">
          <ul className="space-y-2 text-xs">
            {(data?.n1_ranked?.data ?? []).slice(0, 6).map((n, i) => (
              <li key={i} className="rounded-lg border border-border bg-bg p-2.5">
                <div className="flex items-center justify-between mb-0.5">
                  <span className="font-mono truncate">{n.route}</span>
                  <span className="text-foreground-muted shrink-0 ml-2">impact {n.impact_score.toFixed(0)}</span>
                </div>
                <p className="font-mono text-[10px] text-foreground-muted truncate">{n.pattern}</p>
                <p className="text-[10px] text-foreground-muted mt-0.5">{n.occurrences} occurrences · ~{n.avg_queries_per_request} queries/req</p>
              </li>
            ))}
            {(!data?.n1_ranked?.data || data.n1_ranked.data.length === 0) && (
              <li className="text-foreground-muted text-center py-3">No N+1 detections.</li>
            )}
          </ul>
        </Panel>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <Panel title="Recent unresolved errors">
          <ul className="space-y-1.5 text-xs">
            {(data?.errors?.data ?? []).slice(0, 6).map((e) => (
              <li key={e.id} className="flex items-start justify-between gap-2 py-1 border-b border-border last:border-0">
                <div className="min-w-0 flex-1">
                  <p className="truncate"><span className="font-mono text-[10px] text-foreground-muted">{e.type}</span> {e.message}</p>
                  <p className="font-mono text-[10px] text-foreground-muted truncate">{e.route}</p>
                </div>
                <span className="text-foreground-muted shrink-0">{e.count}×</span>
              </li>
            ))}
            {(!data?.errors?.data || data.errors.data.length === 0) && <li className="text-foreground-muted text-center py-3">No unresolved errors.</li>}
          </ul>
        </Panel>

        <Panel title="Go runtime">
          <div className="grid grid-cols-3 gap-2 text-xs">
            <Stat label="Heap alloc" value={data?.runtime?.heap_alloc_mb != null ? `${data.runtime.heap_alloc_mb.toFixed(0)} MB` : '—'} />
            <Stat label="Goroutines" value={data?.runtime?.goroutines ?? '—'} />
            <Stat label="GC pause" value={data?.runtime?.gc_pause_ms != null ? `${data.runtime.gc_pause_ms.toFixed(1)} ms` : '—'} />
          </div>
          {data?.health_checks?.data && data.health_checks.data.length > 0 && (
            <div className="mt-4 pt-3 border-t border-border space-y-1">
              <p className="text-[10px] uppercase font-mono tracking-wider text-foreground-muted">Health checks</p>
              {data.health_checks.data.map((h) => (
                <div key={h.name} className="flex items-center justify-between text-xs">
                  <span className="font-mono">{h.name}</span>
                  <span className={`text-[10px] ${h.status === 'healthy' ? 'text-emerald-700' : 'text-rose-700'}`}>
                    {h.status} · {h.latency_ms}ms
                  </span>
                </div>
              ))}
            </div>
          )}
        </Panel>
      </div>
    </div>
  )
}

const TONE: Record<string, { bg: string; fg: string }> = {
  emerald: { bg: 'bg-emerald-50', fg: 'text-emerald-700' },
  amber: { bg: 'bg-amber-50', fg: 'text-amber-700' },
  sky: { bg: 'bg-sky-50', fg: 'text-sky-700' },
  rose: { bg: 'bg-rose-50', fg: 'text-rose-700' },
}

function Kpi({ label, value, tone, icon: Icon }: { label: string; value: React.ReactNode; tone: keyof typeof TONE; icon: React.ComponentType<{ className?: string }> }) {
  const t = TONE[tone]
  return (
    <div className="rounded-xl border border-border bg-surface p-4">
      <div className="flex items-center justify-between mb-1">
        <span className="text-[10px] font-mono uppercase tracking-wider text-foreground-muted">{label}</span>
        <span className={`h-7 w-7 rounded-md ${t.bg} flex items-center justify-center`}><Icon className={`h-4 w-4 ${t.fg}`} /></span>
      </div>
      <p className={`text-2xl font-semibold tabular-nums ${t.fg}`}>{value}</p>
    </div>
  )
}

function Panel({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="rounded-xl border border-border bg-surface overflow-hidden">
      <div className="px-4 py-3 border-b border-border"><h2 className="text-sm font-semibold">{title}</h2></div>
      <div className="p-4">{children}</div>
    </div>
  )
}

function Stat({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className="rounded-lg border border-border bg-bg px-3 py-2">
      <p className="text-[10px] uppercase tracking-wider text-foreground-muted">{label}</p>
      <p className="text-base font-semibold tabular-nums">{value}</p>
    </div>
  )
}

function Cell({ band, value }: { band?: string; value?: number }) {
  const cls = BAND_CLS[band || 'unknown']
  return (
    <span className={`inline-flex items-center justify-center min-w-[44px] px-2 py-0.5 rounded border text-[10px] font-mono ${cls}`}>
      {value != null ? value.toFixed(0) : '—'}
    </span>
  )
}

// suppress unused-import warning for Cpu (kept for potential future hardware-icon use)
void Cpu
