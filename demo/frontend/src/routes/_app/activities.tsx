import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useMemo, useRef } from 'react'
import { useQuery } from '@tanstack/react-query'
import {
  Activity, LogIn, LogOut, Plus, Pencil, Trash2, CheckCircle2, XCircle,
  ShieldCheck, Banknote, ArrowRightLeft, Upload, Undo2, Printer, User as UserIcon,
} from 'lucide-react'

import api from '@/lib/api'
import { useStaff } from '@/hooks/useStaff'
import { ExportButton } from '@/components/export-button'
import { exportToExcel } from '@/lib/export'

type ASearch = { user_id?: string; from?: string; to?: string; q?: string }

export const Route = createFileRoute('/_app/activities')({
  component: ActivitiesPage,
  validateSearch: (s: Record<string, unknown>): ASearch => ({
    user_id: typeof s.user_id === 'string' ? s.user_id : undefined,
    from: typeof s.from === 'string' ? s.from : undefined,
    to: typeof s.to === 'string' ? s.to : undefined,
    q: typeof s.q === 'string' ? s.q : undefined,
  }),
})

interface ActivityRow {
  id: number
  user_id: number
  user_name: string
  user_email: string
  action: string
  resource: string
  resource_id?: number | null
  description: string
  ip_address?: string
  created_at: string
}

function ActivitiesPage() {
  const search = Route.useSearch()
  const navigate = useNavigate()
  const printRef = useRef<HTMLDivElement>(null)

  const update = (patch: Partial<ASearch>) =>
    navigate({ to: '/activities', search: { ...search, ...patch }, replace: true })

  const list = useQuery({
    queryKey: ['activities', search],
    queryFn: async () => {
      const r = await api.get('/activities', {
        params: {
          user_id: search.user_id,
          from: search.from,
          to: search.to,
          q: search.q,
          per_page: 500,
        },
      })
      return r.data as { data: ActivityRow[]; total: number }
    },
  })

  const staff = useStaff()

  // Group activities by day (YYYY-MM-DD), preserving descending time order.
  const grouped = useMemo(() => {
    const groups = new Map<string, ActivityRow[]>()
    for (const row of list.data?.data || []) {
      const day = row.created_at.slice(0, 10)
      if (!groups.has(day)) groups.set(day, [])
      groups.get(day)!.push(row)
    }
    return Array.from(groups.entries()) // already sorted by API DESC
  }, [list.data])

  const handlePrint = () => {
    // Use a transient print-only stylesheet so only the activities pane appears
    const win = window.open('', '_blank', 'width=900,height=1200')
    if (!win || !printRef.current) return
    const html = printRef.current.innerHTML
    const filterLabel = [
      search.user_id ? `User: ${staff.data?.find((u) => u.user_id === Number(search.user_id))?.name || search.user_id}` : null,
      search.from ? `From: ${search.from}` : null,
      search.to ? `To: ${search.to}` : null,
    ].filter(Boolean).join(' · ')

    win.document.write(`<!doctype html><html><head><title>Activities</title>
      <style>
        * { box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", system-ui, sans-serif; color: #0F172A; padding: 24px; line-height: 1.5; }
        h1 { font-size: 18px; margin: 0 0 4px 0; }
        .sub { font-size: 12px; color: #64748B; margin-bottom: 18px; }
        .day { margin-top: 18px; font-size: 11px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.6px; color: #1565C0; padding-bottom: 4px; border-bottom: 2px solid #1565C0; }
        .row { display: flex; padding: 8px 0; border-bottom: 1px solid #F1F5F9; gap: 12px; align-items: flex-start; font-size: 12px; }
        .time { color: #64748B; font-variant-numeric: tabular-nums; min-width: 70px; }
        .who { color: #475569; font-weight: 600; min-width: 140px; }
        .what { color: #0F172A; flex: 1; }
        .badge { display: inline-block; padding: 1px 6px; border-radius: 4px; font-size: 10px; font-weight: 600; text-transform: uppercase; background: #F1F5F9; color: #64748B; margin-right: 6px; }
      </style></head><body>
      <h1>Grit Motors — Activity Log</h1>
      <div class="sub">${filterLabel || 'All activities'} · Printed ${new Date().toLocaleString()}</div>
      ${html}
    </body></html>`)
    win.document.close()
    setTimeout(() => { win.focus(); win.print() }, 200)
  }

  return (
    <div className="h-full flex flex-col bg-background">
      <header className="px-4 sm:px-6 py-4 border-b border-border-subtle bg-surface">
        <div className="flex items-center justify-between gap-3 flex-wrap">
          <div>
            <h1 className="text-[18px] sm:text-[20px] font-semibold text-foreground inline-flex items-center gap-2">
              <Activity size={20} className="text-accent" /> Activity Log
            </h1>
            <p className="text-[12.5px] text-foreground-muted mt-0.5">
              Every meaningful action — useful for accountability and reviewing staff work.
            </p>
          </div>
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={handlePrint}
              disabled={!list.data?.data?.length}
              className="h-9 px-3 rounded-lg border border-border bg-surface hover:bg-surface-hover text-[12.5px] font-medium text-foreground-secondary transition inline-flex items-center gap-1.5 disabled:opacity-50"
            >
              <Printer size={13} /> Print
            </button>
            <ExportButton
              onClick={() => exportToExcel(
                (list.data?.data || []).map((r) => ({
                  When: new Date(r.created_at).toLocaleString(),
                  User: r.user_name,
                  Email: r.user_email,
                  Action: r.action,
                  Resource: r.resource,
                  'Resource ID': r.resource_id ?? '',
                  Description: r.description,
                  IP: r.ip_address ?? '',
                })),
                'activities',
                'Activities',
              )}
              disabled={!list.data?.data?.length}
            />
          </div>
        </div>

        {/* Filters */}
        <div className="mt-3 flex flex-wrap items-center gap-2">
          <select
            value={search.user_id || ''}
            onChange={(e) => update({ user_id: e.target.value || undefined })}
            className="h-8 px-2 rounded-md border border-border bg-surface text-[12px] text-foreground"
          >
            <option value="">All users</option>
            {staff.data?.map((u) => (
              <option key={u.user_id} value={u.user_id}>{u.name}</option>
            ))}
          </select>
          <input
            type="date"
            value={search.from || ''}
            onChange={(e) => update({ from: e.target.value || undefined })}
            className="h-8 px-2 rounded-md border border-border bg-surface text-[12px] text-foreground"
            title="From"
          />
          <span className="text-foreground-muted text-[11px]">to</span>
          <input
            type="date"
            value={search.to || ''}
            onChange={(e) => update({ to: e.target.value || undefined })}
            className="h-8 px-2 rounded-md border border-border bg-surface text-[12px] text-foreground"
            title="To"
          />
          <input
            type="text"
            value={search.q || ''}
            onChange={(e) => update({ q: e.target.value || undefined })}
            placeholder="Search description…"
            className="h-8 px-2 rounded-md border border-border bg-surface text-[12px] text-foreground flex-1 min-w-[180px] max-w-[260px]"
          />
          {(search.user_id || search.from || search.to || search.q) && (
            <button
              type="button"
              onClick={() => navigate({ to: '/activities', search: {}, replace: true })}
              className="text-[11.5px] text-foreground-muted hover:text-foreground underline"
            >
              Clear
            </button>
          )}
          <span className="ml-auto text-[11.5px] text-foreground-muted">
            {list.data?.total ?? 0} {list.data?.total === 1 ? 'activity' : 'activities'}
          </span>
        </div>
      </header>

      {/* Body */}
      <div className="flex-1 overflow-y-auto px-4 sm:px-6 py-5">
        {list.isLoading ? (
          <div className="text-center py-16 text-foreground-muted text-[13px]">Loading…</div>
        ) : grouped.length === 0 ? (
          <div className="text-center py-16 text-foreground-muted text-[13px]">
            <Activity size={36} className="mx-auto mb-3 opacity-30" />
            <p>No activities match these filters.</p>
          </div>
        ) : (
          <div ref={printRef} className="max-w-5xl mx-auto space-y-6">
            {grouped.map(([day, rows]) => (
              <DayGroup key={day} day={day} rows={rows} />
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

function DayGroup({ day, rows }: { day: string; rows: ActivityRow[] }) {
  const date = new Date(day + 'T12:00:00')
  const today = new Date().toISOString().slice(0, 10)
  const yesterdayStr = (() => { const d = new Date(); d.setDate(d.getDate() - 1); return d.toISOString().slice(0, 10) })()
  const label = day === today ? 'Today' : day === yesterdayStr ? 'Yesterday' : date.toLocaleDateString(undefined, { weekday: 'long', month: 'long', day: 'numeric', year: 'numeric' })

  return (
    <section className="day">
      <div className="sticky top-0 -mx-4 sm:-mx-6 px-4 sm:px-6 py-1.5 bg-background/95 backdrop-blur-sm border-b border-border-subtle z-10 mb-2">
        <h3 className="text-[11px] font-bold uppercase tracking-wider text-accent">
          {label}
          <span className="ml-2 text-foreground-muted normal-case font-medium">
            {rows.length} {rows.length === 1 ? 'event' : 'events'}
          </span>
        </h3>
      </div>
      <ul className="divide-y divide-border-subtle">
        {rows.map((r) => (
          <ActivityItem key={r.id} row={r} />
        ))}
      </ul>
    </section>
  )
}

function ActivityItem({ row }: { row: ActivityRow }) {
  const time = new Date(row.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  const { Icon, tone } = iconFor(row.action)
  return (
    <li className="row flex items-start gap-3 py-2.5">
      <span className="time text-[11.5px] text-foreground-muted tabular-nums shrink-0 w-[58px] pt-0.5">
        {time}
      </span>
      <span className={`w-7 h-7 rounded-full flex items-center justify-center shrink-0 ${tone}`}>
        <Icon size={13} />
      </span>
      <div className="flex-1 min-w-0">
        <p className="text-[13px] text-foreground leading-snug">
          <span className="badge inline-block mr-1.5 px-1.5 py-0.5 rounded text-[9.5px] font-semibold uppercase tracking-wider bg-surface-2 text-foreground-muted">
            {row.resource.replace('_', ' ')}
          </span>
          {row.description}
        </p>
        <p className="who text-[11.5px] text-foreground-muted mt-0.5 inline-flex items-center gap-1">
          <UserIcon size={10} /> {row.user_name}
          {row.ip_address && <span className="text-foreground-muted/70 font-mono ml-1">· {row.ip_address}</span>}
        </p>
      </div>
    </li>
  )
}

function iconFor(action: string) {
  switch (action) {
    case 'login':    return { Icon: LogIn,        tone: 'bg-success-light text-success-dark' }
    case 'logout':   return { Icon: LogOut,       tone: 'bg-surface-2 text-foreground-muted' }
    case 'create':   return { Icon: Plus,         tone: 'bg-accent-light text-accent-hover' }
    case 'update':   return { Icon: Pencil,       tone: 'bg-info-light text-info-dark' }
    case 'delete':   return { Icon: Trash2,       tone: 'bg-danger-light text-danger-dark' }
    case 'approve':  return { Icon: CheckCircle2, tone: 'bg-success-light text-success-dark' }
    case 'reject':   return { Icon: XCircle,      tone: 'bg-danger-light text-danger-dark' }
    case 'verify':   return { Icon: ShieldCheck,  tone: 'bg-success-light text-success-dark' }
    case 'disburse': return { Icon: Banknote,     tone: 'bg-warning-light text-warning-dark' }
    case 'transfer': return { Icon: ArrowRightLeft, tone: 'bg-info-light text-info-dark' }
    case 'import':   return { Icon: Upload,       tone: 'bg-accent-light text-accent-hover' }
    case 'return':   return { Icon: Undo2,        tone: 'bg-warning-light text-warning-dark' }
    default:         return { Icon: Activity,     tone: 'bg-surface-2 text-foreground-muted' }
  }
}
