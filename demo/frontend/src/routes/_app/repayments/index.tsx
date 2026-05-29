import { createFileRoute, useNavigate, Link } from '@tanstack/react-router'
import { CircleDollarSign, Check, X } from 'lucide-react'
import toast from 'react-hot-toast'

import { useRepayments, useRepayment, useApproveRepayment, useRejectRepayment } from '@/hooks/useRepayments'
import {
  TwoPane, ListPane, ListRow, DetailPane, DetailSection, Field, FieldGrid,
} from '@/components/two-pane'
import { StatusPill } from '@/components/status-badge'
import { ExportButton } from '@/components/export-button'
import { exportToExcel } from '@/lib/export'
import { formatUGX } from '@/lib/utils'
import type { RepaymentStatus } from '@/types'

type RSearch = { selected?: number; status?: RepaymentStatus | 'all'; loan_id?: string; q?: string; from?: string; to?: string }

export const Route = createFileRoute('/_app/repayments/')({
  component: RepaymentsPage,
  validateSearch: (s: Record<string, unknown>): RSearch => ({
    selected: s.selected ? Number(s.selected) : undefined,
    status: typeof s.status === 'string' ? (s.status as any) : undefined,
    loan_id: typeof s.loan_id === 'string' ? s.loan_id : undefined,
    q: typeof s.q === 'string' ? s.q : undefined,
    from: typeof s.from === 'string' ? s.from : undefined,
    to: typeof s.to === 'string' ? s.to : undefined,
  }),
})

const STATUS_TABS: { val: RepaymentStatus | 'all'; label: string }[] = [
  { val: 'all', label: 'All' },
  { val: 'pending', label: 'Pending' },
  { val: 'approved', label: 'Approved' },
  { val: 'rejected', label: 'Rejected' },
  { val: 'failed', label: 'Failed' },
]

function RepaymentsPage() {
  const search = Route.useSearch()
  const navigate = useNavigate()

  const list = useRepayments({
    status: search.status === 'all' ? undefined : search.status,
    loan_id: search.loan_id,
    from_date: search.from,
    to_date: search.to,
  })
  const selected = useRepayment(search.selected)

  const update = (patch: Partial<RSearch>) =>
    navigate({ to: '/repayments', search: { ...search, ...patch }, replace: true })

  const filtered = (list.data?.data || []).filter((r) => {
    if (!search.q) return true
    const q = search.q.toLowerCase()
    return (
      r.transaction_ref?.toLowerCase().includes(q) ||
      r.dgateway_reference?.toLowerCase().includes(q) ||
      r.loan?.loan_number?.toLowerCase().includes(q)
    )
  })

  return (
    <TwoPane>
      <ListPane
        title="Repayments"
        count={list.data?.total}
        search={search.q || ''}
        onSearch={(q) => update({ q: q || undefined })}
        searchPlaceholder="Loan #, txn ref..."
        refreshKey={['repayments']}
        extraAction={
          <ExportButton
            onClick={() => exportToExcel(
              filtered.map((r: any) => ({
                'Receipt': r.receipt,
                'Loan #': r.loan?.loan_number,
                Amount: r.amount,
                Method: r.payment_method,
                'Collection Date': new Date(r.collection_date).toLocaleString(),
                'Transaction Ref': r.transaction_ref,
                'DGateway Ref': r.dgateway_reference,
                'DGateway Fee': r.dgateway_fee,
                'Net Amount': r.dgateway_net_amount,
                Status: r.status,
                'Verified At': r.verified_at ? new Date(r.verified_at).toLocaleString() : '',
              })),
              'repayments',
              'Repayments',
            )}
            disabled={filtered.length === 0}
          />
        }
        filters={
          <div className="space-y-1.5">
            <div className="flex gap-1 flex-wrap">
              {STATUS_TABS.map((t) => (
                <button
                  key={t.val}
                  type="button"
                  onClick={() => update({ status: t.val === 'all' ? undefined : t.val })}
                  className={`px-2.5 py-1 rounded-md text-[11.5px] font-medium transition ${
                    (search.status || 'all') === t.val
                      ? 'bg-foreground text-white'
                      : 'bg-surface-2 text-foreground-secondary hover:bg-surface-hover'
                  }`}
                >
                  {t.label}
                </button>
              ))}
            </div>
            <div className="flex items-center gap-1.5 flex-wrap">
              <input
                type="date"
                value={search.from || ''}
                onChange={(e) => update({ from: e.target.value || undefined })}
                className="px-2 h-7 rounded-md border border-border bg-surface text-[11.5px] text-foreground"
                title="From date"
              />
              <span className="text-foreground-muted text-[11px]">to</span>
              <input
                type="date"
                value={search.to || ''}
                onChange={(e) => update({ to: e.target.value || undefined })}
                className="px-2 h-7 rounded-md border border-border bg-surface text-[11.5px] text-foreground"
                title="To date"
              />
              {(search.from || search.to) && (
                <button
                  type="button"
                  onClick={() => update({ from: undefined, to: undefined })}
                  className="text-[11px] text-foreground-muted hover:text-foreground underline"
                >
                  Clear
                </button>
              )}
            </div>
          </div>
        }
      >
        {list.isLoading && <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">Loading…</div>}
        {!list.isLoading && filtered.length === 0 && (
          <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">No repayments match.</div>
        )}
        {filtered.map((r) => (
          <ListRow
            key={r.id}
            selected={search.selected === r.id}
            onClick={() => update({ selected: r.id })}
            icon={<CircleDollarSign size={14} className="text-foreground-secondary" />}
            title={formatUGX(r.amount)}
            subtitle={
              <span>
                <span className="font-mono text-accent">{r.loan?.loan_number || `Loan #${r.loan_id}`}</span>
                <span className="capitalize"> · {r.payment_method.replace('_', ' ')}</span>
              </span>
            }
            rightTop={new Date(r.collection_date).toLocaleDateString()}
            rightBottom={<StatusPill status={r.status} size="sm" />}
          />
        ))}
      </ListPane>

      {!search.selected ? (
        <DetailPane empty emptyTitle="No repayment selected" emptyHint="Pick a repayment from the list to verify or inspect it." />
      ) : selected.isLoading || !selected.data ? (
        <DetailPane><div className="text-foreground-muted">Loading…</div></DetailPane>
      ) : (
        <RepaymentDetail r={selected.data} />
      )}
    </TwoPane>
  )
}

function RepaymentDetail({ r }: { r: NonNullable<ReturnType<typeof useRepayment>['data']> }) {
  const approve = useApproveRepayment()
  const reject = useRejectRepayment()

  const handleApprove = async () => {
    try { await approve.mutateAsync(r.id); toast.success('Approved & applied to loan') } catch (e: any) { toast.error(e.response?.data?.error || 'Failed') }
  }
  const handleReject = async () => {
    const reason = prompt('Reason for rejection?')
    if (reason === null) return
    try { await reject.mutateAsync({ id: r.id, reason }); toast.success('Rejected') } catch (e: any) { toast.error(e.response?.data?.error || 'Failed') }
  }

  return (
    <DetailPane
      header={
        <div className="flex items-start justify-between gap-4">
          <div>
            <p className="text-[10.5px] uppercase tracking-wider text-foreground-muted font-mono">{r.receipt}</p>
            <h1 className="text-[20px] font-semibold text-foreground tabular-nums mt-0.5">{formatUGX(r.amount)}</h1>
            <p className="text-[12.5px] text-foreground-muted mt-0.5">
              <Link to="/loans" search={{ selected: r.loan_id }} className="font-mono text-accent hover:underline">{r.loan?.loan_number || `Loan #${r.loan_id}`}</Link>
              {' · '}
              <span className="capitalize">{r.payment_method.replace('_', ' ')}</span>
            </p>
          </div>
          <StatusPill status={r.status} />
        </div>
      }
    >
      {r.status === 'pending' && (
        <div className="flex gap-2 mb-6">
          <button onClick={handleApprove} className="h-9 px-3.5 rounded-lg bg-accent text-white text-[13px] font-medium hover:bg-accent-hover transition inline-flex items-center gap-1.5">
            <Check size={14} /> Approve & Apply
          </button>
          <button onClick={handleReject} className="h-9 px-3.5 rounded-lg bg-danger-light text-danger-dark text-[13px] font-medium hover:bg-danger/20 transition inline-flex items-center gap-1.5">
            <X size={14} /> Reject
          </button>
        </div>
      )}

      <DetailSection title="Repayment">
        <FieldGrid columns={2}>
          <Field label="Amount" value={formatUGX(r.amount)} />
          <Field label="Method" value={<span className="capitalize">{r.payment_method.replace('_', ' ')}</span>} />
          <Field label="Collection date" value={new Date(r.collection_date).toLocaleString()} />
          <Field label="Receipt #" value={<span className="font-mono">{r.receipt}</span>} />
          {r.transaction_ref && <Field label="Transaction ref" value={<span className="font-mono">{r.transaction_ref}</span>} />}
          {r.dgateway_reference && <Field label="DGateway ref" value={<span className="font-mono">{r.dgateway_reference}</span>} />}
          {r.dgateway_provider && <Field label="Provider" value={r.dgateway_provider} />}
          {r.dgateway_fee > 0 && <Field label="DGateway fee" value={formatUGX(r.dgateway_fee)} />}
          {r.dgateway_net_amount > 0 && <Field label="Net to us" value={formatUGX(r.dgateway_net_amount)} />}
        </FieldGrid>
      </DetailSection>

      <DetailSection title="Audit">
        <FieldGrid columns={2}>
          <Field label="Collected by" value={r.collector?.name || `User #${r.collected_by}`} />
          <Field label="Verified at" value={r.verified_at ? new Date(r.verified_at).toLocaleString() : null} />
        </FieldGrid>
      </DetailSection>

      {r.notes && (
        <DetailSection title="Notes">
          <p className="text-[13px] text-foreground whitespace-pre-line">{r.notes}</p>
        </DetailSection>
      )}

      {r.dgateway_fail_reason && (
        <DetailSection title="Failure Reason">
          <p className="text-[13px] text-danger-dark">{r.dgateway_fail_reason}</p>
        </DetailSection>
      )}
    </DetailPane>
  )
}
