import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { Shield, ExternalLink, AlertTriangle, Zap, Activity as ActivityIcon, AlertCircle, Globe } from 'lucide-react'
import api from '@/lib/api'

const API_URL = import.meta.env.VITE_API_URL || ''

export const Route = createFileRoute('/_app/system/security')({
  component: SecurityPage,
})

interface ThreatItem {
  id: string
  type: string
  severity?: string
  cvss?: number
  source_ip?: string
  route?: string
  created_at?: string
}

interface Summary {
  summary?: { threats_total?: number; threats_24h?: number; blocked?: number; actors?: number }
  score?: { score?: number; grade?: string }
  threats?: { data?: ThreatItem[] }
  auth_shield?: { enabled?: boolean; failures?: number; locked?: Array<{ ip: string; until: string }> }
  csp_top?: { data?: Array<{ violated_directive: string; blocked_uri: string; count: number }> }
  _errors?: Record<string, string>
}

function SecurityPage() {
  const { data, error, isLoading } = useQuery<Summary>({
    queryKey: ['admin', 'security', 'summary'],
    queryFn: async () => {
      const res = await api.get('/api/admin/security/summary')
      return (res.data?.data ?? res.data) as Summary
    },
    refetchInterval: 20_000,
  })

  return (
    <div className="space-y-6 p-6 max-w-6xl mx-auto">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground flex items-center gap-2.5">
            <Shield className="h-6 w-6 text-grit-blue" /> Security
          </h1>
          <p className="text-sm text-foreground-muted mt-1">
            Live summary from Sentinel — score, threats, AuthShield, CSP, performance.
          </p>
        </div>
        <a
          href={`${API_URL}/sentinel/ui`}
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
            <strong>Couldn&apos;t reach Sentinel.</strong> Make sure <code className="px-1 py-0.5 rounded bg-amber-100 text-xs font-mono">SENTINEL_ENABLED=true</code> and the API is running.
          </div>
        </div>
      )}

      {/* Top-line scorecard */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        <Kpi title="Security score" value={data?.score?.score ?? '—'} sub={data?.score?.grade} icon={Shield} tone="emerald" />
        <Kpi title="Threats (24h)" value={data?.summary?.threats_24h ?? '—'} sub={`${data?.summary?.threats_total ?? 0} total`} icon={AlertCircle} tone="amber" />
        <Kpi title="Blocked" value={data?.summary?.blocked ?? '—'} sub="auto by WAF" icon={Zap} tone="sky" />
        <Kpi title="Actors" value={data?.summary?.actors ?? '—'} sub="unique sources" icon={Globe} tone="slate" />
      </div>

      {/* Recent threats */}
      <Panel
        title="Recent threats"
        right={<a href={`${API_URL}/sentinel/ui/threats`} target="_blank" rel="noopener noreferrer" className="text-xs text-grit-blue hover:underline">View all</a>}
      >
        <div className="divide-y divide-border">
          {(data?.threats?.data ?? []).slice(0, 8).map((t) => (
            <div key={t.id} className="flex items-center gap-3 px-1 py-2.5 text-sm">
              <SeverityChip severity={t.severity} cvss={t.cvss} />
              <span className="font-mono text-xs truncate flex-1">{t.type}</span>
              <span className="text-foreground-muted text-xs hidden md:inline truncate max-w-[180px]">{t.route}</span>
              <span className="text-foreground-muted text-xs font-mono">{t.source_ip}</span>
            </div>
          ))}
          {(!data?.threats?.data || data.threats.data.length === 0) && (
            <div className="py-6 text-center text-sm text-foreground-muted">No threats in window.</div>
          )}
        </div>
      </Panel>

      {/* AuthShield + CSP side-by-side */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <Panel title="AuthShield">
          {data?.auth_shield?.enabled === false ? (
            <p className="text-sm text-foreground-muted">AuthShield is disabled in this environment.</p>
          ) : (
            <div className="space-y-3">
              <div className="grid grid-cols-2 gap-2 text-sm">
                <Stat label="Failures (rolling)" value={data?.auth_shield?.failures ?? 0} />
                <Stat label="Locked IPs" value={data?.auth_shield?.locked?.length ?? 0} />
              </div>
              <ul className="space-y-1.5 text-xs font-mono">
                {(data?.auth_shield?.locked ?? []).slice(0, 5).map((l) => (
                  <li key={l.ip} className="flex justify-between gap-2 text-foreground-muted">
                    <span>{l.ip}</span>
                    <span>until {l.until}</span>
                  </li>
                ))}
              </ul>
            </div>
          )}
        </Panel>
        <Panel title="Top CSP violations">
          <ul className="space-y-1.5 text-xs font-mono">
            {(data?.csp_top?.data ?? []).slice(0, 5).map((v, i) => (
              <li key={i} className="flex justify-between gap-2">
                <span className="truncate">{v.violated_directive}</span>
                <span className="text-foreground-muted shrink-0">{v.count}×</span>
              </li>
            ))}
            {(!data?.csp_top?.data || data.csp_top.data.length === 0) && (
              <li className="text-foreground-muted">No violations.</li>
            )}
          </ul>
        </Panel>
      </div>
    </div>
  )
}

const TONE: Record<string, { bg: string; fg: string }> = {
  emerald: { bg: 'bg-emerald-50', fg: 'text-emerald-700' },
  amber: { bg: 'bg-amber-50', fg: 'text-amber-700' },
  sky: { bg: 'bg-sky-50', fg: 'text-sky-700' },
  slate: { bg: 'bg-slate-50', fg: 'text-slate-700' },
  rose: { bg: 'bg-rose-50', fg: 'text-rose-700' },
}

function Kpi({ title, value, sub, icon: Icon, tone }: { title: string; value: React.ReactNode; sub?: React.ReactNode; icon: React.ComponentType<{ className?: string }>; tone: keyof typeof TONE }) {
  const t = TONE[tone]
  return (
    <div className="rounded-xl border border-border bg-surface p-4">
      <div className="flex items-center justify-between mb-1">
        <span className="text-[10px] font-mono uppercase tracking-wider text-foreground-muted">{title}</span>
        <span className={`h-7 w-7 rounded-md ${t.bg} flex items-center justify-center`}><Icon className={`h-4 w-4 ${t.fg}`} /></span>
      </div>
      <p className={`text-2xl font-semibold tabular-nums ${t.fg}`}>{value}</p>
      {sub != null && <p className="text-[11px] text-foreground-muted mt-0.5">{sub}</p>}
    </div>
  )
}

function Panel({ title, children, right }: { title: string; children: React.ReactNode; right?: React.ReactNode }) {
  return (
    <div className="rounded-xl border border-border bg-surface overflow-hidden">
      <div className="px-4 py-3 border-b border-border flex items-center justify-between">
        <h2 className="text-sm font-semibold">{title}</h2>
        {right}
      </div>
      <div className="p-4">{children}</div>
    </div>
  )
}

function Stat({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className="rounded-lg border border-border bg-bg px-3 py-2">
      <p className="text-[10px] uppercase tracking-wider text-foreground-muted">{label}</p>
      <p className="text-lg font-semibold tabular-nums">{value}</p>
    </div>
  )
}

function SeverityChip({ severity, cvss }: { severity?: string; cvss?: number }) {
  const tone =
    severity === 'critical' ? 'bg-rose-50 text-rose-700 border-rose-200' :
    severity === 'high'     ? 'bg-amber-50 text-amber-700 border-amber-200' :
    severity === 'medium'   ? 'bg-sky-50 text-sky-700 border-sky-200' :
                              'bg-slate-50 text-slate-600 border-slate-200'
  return (
    <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded border text-[10px] font-mono ${tone}`}>
      {severity ?? 'info'}
      {cvss != null && cvss > 0 && <span>· {cvss.toFixed(1)}</span>}
    </span>
  )
}
