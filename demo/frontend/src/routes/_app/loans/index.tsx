import { createFileRoute, useNavigate, Link } from '@tanstack/react-router'
import { useState } from 'react'
import { HandCoins, Check, X, ArrowRight, CircleDollarSign, Smartphone, Loader2, Pencil } from 'lucide-react'
import toast from 'react-hot-toast'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'

import { useLoans, useLoan, useApproveLoan, useDisburseLoan, useRejectLoan } from '@/hooks/useLoans'
import { useCreateRepayment } from '@/hooks/useRepayments'
import { useDGatewayCollect, pollDGatewayStatus } from '@/hooks/useDGateway'
import {
  TwoPane, ListPane, ListRow, DetailPane, DetailSection, Field, FieldGrid,
} from '@/components/two-pane'
import { Drawer } from '@/components/drawer'
import { TextField, SelectField, FormActions, TextAreaField } from '@/components/form'
import { CurrencyField } from '@/components/currency-field'
import { StatusPill } from '@/components/status-badge'
import { ExportButton } from '@/components/export-button'
import { exportToExcel } from '@/lib/export'
import { formatUGX } from '@/lib/utils'
import type { LoanStatus } from '@/types'

type LoansSearch = { selected?: number; q?: string; status?: LoanStatus | 'all'; borrower_id?: string }

export const Route = createFileRoute('/_app/loans/')({
  component: LoansPage,
  validateSearch: (s: Record<string, unknown>): LoansSearch => ({
    selected: s.selected ? Number(s.selected) : undefined,
    q: typeof s.q === 'string' ? s.q : undefined,
    status: typeof s.status === 'string' ? (s.status as any) : undefined,
    borrower_id: typeof s.borrower_id === 'string' ? s.borrower_id : undefined,
  }),
})

const STATUS_TABS: { val: LoanStatus | 'all'; label: string }[] = [
  { val: 'all', label: 'All' },
  { val: 'pending', label: 'Pending' },
  { val: 'active', label: 'Active' },
  { val: 'completed', label: 'Completed' },
  { val: 'defaulted', label: 'Defaulted' },
]

function LoansPage() {
  const search = Route.useSearch()
  const navigate = useNavigate()

  const list = useLoans({
    status: search.status === 'all' ? undefined : search.status,
    borrower_id: search.borrower_id,
  })
  const selected = useLoan(search.selected)

  const update = (patch: Partial<LoansSearch>) =>
    navigate({ to: '/loans', search: { ...search, ...patch }, replace: true })

  const filtered = (list.data?.data || []).filter((l) => {
    if (!search.q) return true
    const q = search.q.toLowerCase()
    return (
      l.loan_number?.toLowerCase().includes(q) ||
      l.borrower?.full_name?.toLowerCase().includes(q) ||
      l.motorcycle?.number_plate?.toLowerCase().includes(q)
    )
  })

  return (
    <TwoPane>
      <ListPane
        title="Loans"
        count={list.data?.total}
        search={search.q || ''}
        onSearch={(q) => update({ q: q || undefined })}
        searchPlaceholder="Loan #, borrower, plate..."
        newLabel="New Loan"
        onNew={() => navigate({ to: '/loans/new' })}
        refreshKey={['loans']}
        extraAction={
          <ExportButton
            onClick={() => exportToExcel(
              filtered.map((l: any) => ({
                'Loan #': l.loan_number,
                Borrower: l.borrower?.full_name,
                'Borrower Phone': l.borrower?.phone,
                Motorcycle: l.motorcycle?.number_plate,
                Branch: l.branch?.name,
                Principal: l.principal_amount,
                'Initial Deposit': l.initial_deposit,
                Disbursed: l.disbursed_amount,
                'Total Amount': l.total_amount,
                'Balance Remaining': l.balance_remaining,
                'Interest Rate %': l.interest_rate,
                'Interest Method': l.interest_method,
                'Repayment Cycle': l.repayment_cycle,
                Duration: l.duration,
                Status: l.status,
                'Disbursed Date': l.disbursement_date ? new Date(l.disbursement_date).toLocaleDateString() : '',
                'Maturity Date': l.maturity_date ? new Date(l.maturity_date).toLocaleDateString() : '',
                'Created At': l.created_at ? new Date(l.created_at).toLocaleDateString() : '',
              })),
              'loans',
              'Loans',
            )}
            disabled={filtered.length === 0}
          />
        }
        filters={
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
        }
      >
        {list.isLoading && <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">Loading…</div>}
        {!list.isLoading && filtered.length === 0 && (
          <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">No loans match.</div>
        )}
        {filtered.map((l) => (
          <ListRow
            key={l.id}
            selected={search.selected === l.id}
            onClick={() => update({ selected: l.id })}
            icon={<HandCoins size={14} className="text-foreground-secondary" />}
            title={l.borrower?.full_name || `Borrower #${l.borrower_id}`}
            subtitle={
              <span>
                <span className="font-mono text-accent">{l.loan_number}</span>
                {l.motorcycle?.number_plate && <> · <span className="font-mono">{l.motorcycle.number_plate}</span></>}
              </span>
            }
            rightTop={formatUGX(l.balance_remaining)}
            rightBottom={<StatusPill status={l.status} size="sm" />}
          />
        ))}
      </ListPane>

      {!search.selected ? (
        <DetailPane empty emptyTitle="No loan selected" emptyHint="Pick a loan from the list, or click + New Loan." />
      ) : selected.isLoading || !selected.data ? (
        <DetailPane><div className="text-foreground-muted">Loading…</div></DetailPane>
      ) : (
        <LoanDetail data={selected.data} />
      )}
    </TwoPane>
  )
}

/* ─────────────── Detail pane ─────────────── */

function LoanDetail({ data }: { data: NonNullable<ReturnType<typeof useLoan>['data']> }) {
  const { loan, schedule, repayments } = data
  const approve = useApproveLoan()
  const reject = useRejectLoan()
  const disburse = useDisburseLoan()
  const [collectOpen, setCollectOpen] = useState(false)
  const [momoOpen, setMomoOpen] = useState(false)

  const cycleUnit = loan.repayment_cycle === 'monthly' ? 'months' : loan.repayment_cycle === 'biweekly' ? 'biweeks' : 'weeks'

  const handleApprove = async () => {
    try { await approve.mutateAsync(loan.id); toast.success('Approved') } catch (e: any) { toast.error(e.response?.data?.error || 'Failed') }
  }
  const handleReject = async () => {
    if (!confirm('Reject this loan? The motorcycle will be released.')) return
    try { await reject.mutateAsync(loan.id); toast.success('Rejected') } catch (e: any) { toast.error(e.response?.data?.error || 'Failed') }
  }
  const handleDisburse = async () => {
    try { await disburse.mutateAsync(loan.id); toast.success('Disbursed') } catch (e: any) { toast.error(e.response?.data?.error || 'Failed') }
  }

  return (
    <DetailPane
      header={
        <div className="flex items-start justify-between gap-4">
          <div className="min-w-0">
            <p className="text-[10.5px] uppercase tracking-wider text-foreground-muted font-mono">{loan.loan_number}</p>
            <h1 className="text-[17px] font-semibold text-foreground mt-0.5">
              <Link to="/borrowers" search={{ selected: loan.borrower_id }} className="hover:text-accent">
                {loan.borrower?.full_name || `Borrower #${loan.borrower_id}`}
              </Link>
            </h1>
            <p className="text-[12.5px] text-foreground-muted mt-0.5">
              {loan.motorcycle?.number_plate ? (
                <>Motorcycle <Link to="/motorcycles" search={{ selected: loan.motorcycle_id! }} className="font-mono text-accent hover:underline">{loan.motorcycle.number_plate}</Link></>
              ) : (
                'Working capital loan (no asset)'
              )}
            </p>
          </div>
          <StatusPill status={loan.status} />
        </div>
      }
    >
      {/* Action bar */}
      {(loan.status === 'pending' || loan.status === 'approved' || loan.status === 'active') && (
        <div className="flex flex-wrap gap-2 mb-6">
          {loan.status === 'pending' && (
            <>
              <button onClick={handleApprove} className="h-9 px-3.5 rounded-lg bg-accent text-white text-[13px] font-medium hover:bg-accent-hover transition inline-flex items-center gap-1.5">
                <Check size={14} /> Approve
              </button>
              <button onClick={handleReject} className="h-9 px-3.5 rounded-lg bg-danger-light text-danger-dark text-[13px] font-medium hover:bg-danger/20 transition inline-flex items-center gap-1.5">
                <X size={14} /> Reject
              </button>
            </>
          )}
          {loan.status === 'approved' && (
            <button onClick={handleDisburse} className="h-9 px-3.5 rounded-lg bg-accent text-white text-[13px] font-medium hover:bg-accent-hover transition inline-flex items-center gap-1.5">
              <ArrowRight size={14} /> Disburse Loan
            </button>
          )}
          {loan.status === 'active' && (
            <>
              <button onClick={() => setMomoOpen(true)} className="h-9 px-3.5 rounded-lg bg-accent text-white text-[13px] font-medium hover:bg-accent-hover transition inline-flex items-center gap-1.5">
                <Smartphone size={14} /> Pay via Mobile Money
              </button>
              <button onClick={() => setCollectOpen(true)} className="h-9 px-3.5 rounded-lg border border-border bg-surface text-[13px] font-medium text-foreground-secondary hover:bg-surface-hover transition inline-flex items-center gap-1.5">
                <CircleDollarSign size={14} /> Record Cash
              </button>
            </>
          )}
        </div>
      )}

      <DetailSection title="Money">
        <FieldGrid columns={4}>
          <Field label="Principal" value={formatUGX(loan.principal_amount)} />
          <Field label="Disbursed" value={formatUGX(loan.disbursed_amount)} />
          <Field label="Total to repay" value={formatUGX(loan.total_amount)} />
          <Field label="Outstanding" value={<span className="font-semibold text-accent">{formatUGX(loan.balance_remaining)}</span>} />
        </FieldGrid>
      </DetailSection>

      <DetailSection title="Terms">
        <FieldGrid columns={4}>
          <Field label="Installment" value={formatUGX(loan.installment_amount)} />
          <Field label="Duration" value={`${loan.duration} ${cycleUnit}`} />
          <Field label="Interest" value={`${loan.interest_rate}% (${loan.interest_method.replace('_', ' ')})`} />
          <Field label="Maturity" value={loan.maturity_date ? new Date(loan.maturity_date).toLocaleDateString() : null} />
        </FieldGrid>
      </DetailSection>

      <DetailSection title={`Repayment Schedule (${schedule.filter((s) => s.is_paid).length}/${schedule.length} paid)`}>
        <div className="border border-border rounded-lg overflow-hidden">
          <table className="w-full text-[12.5px] tabular-nums">
            <thead className="bg-surface-2 text-[11px] uppercase tracking-wider text-foreground-muted">
              <tr>
                <th className="text-left py-2 px-3 font-semibold">#</th>
                <th className="text-left py-2 px-3 font-semibold">Due</th>
                <th className="text-right py-2 px-3 font-semibold">Principal</th>
                <th className="text-right py-2 px-3 font-semibold">Interest</th>
                <th className="text-right py-2 px-3 font-semibold">Total</th>
                <th className="text-right py-2 px-3 font-semibold">Paid</th>
                <th className="text-left py-2 px-3 font-semibold">Status</th>
              </tr>
            </thead>
            <tbody>
              {schedule.map((s) => (
                <tr key={s.id} className={`border-t border-border-subtle ${s.is_overdue && !s.is_paid ? 'bg-danger-light/40' : ''}`}>
                  <td className="py-2 px-3 font-mono text-foreground-muted">{s.installment_number}</td>
                  <td className="py-2 px-3">{new Date(s.due_date).toLocaleDateString()}</td>
                  <td className="py-2 px-3 text-right">{formatUGX(s.principal_amount)}</td>
                  <td className="py-2 px-3 text-right">{formatUGX(s.interest_amount)}</td>
                  <td className="py-2 px-3 text-right font-medium">{formatUGX(s.total_amount)}</td>
                  <td className="py-2 px-3 text-right">{formatUGX(s.paid_amount)}</td>
                  <td className="py-2 px-3">
                    {s.is_paid ? (
                      <StatusPill status="paid" size="sm" />
                    ) : s.is_overdue ? (
                      <StatusPill status="failed" size="sm" />
                    ) : (
                      <StatusPill status="pending" size="sm" />
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </DetailSection>

      <DetailSection
        title={`Repayment History (${repayments.length})`}
        action={<Link to="/repayments" search={{ loan_id: String(loan.id) }} className="text-[11.5px] text-accent font-medium hover:underline">View all</Link>}
      >
        {repayments.length === 0 ? (
          <p className="text-[13px] text-foreground-muted">No repayments yet.</p>
        ) : (
          <div className="border border-border rounded-lg overflow-hidden">
            <table className="w-full text-[12.5px] tabular-nums">
              <thead className="bg-surface-2 text-[11px] uppercase tracking-wider text-foreground-muted">
                <tr>
                  <th className="text-left py-2 px-3 font-semibold">Date</th>
                  <th className="text-left py-2 px-3 font-semibold">Method</th>
                  <th className="text-right py-2 px-3 font-semibold">Amount</th>
                  <th className="text-left py-2 px-3 font-semibold">Status</th>
                </tr>
              </thead>
              <tbody>
                {repayments.map((r) => (
                  <tr key={r.id} className="border-t border-border-subtle">
                    <td className="py-2 px-3">{new Date(r.collection_date).toLocaleDateString()}</td>
                    <td className="py-2 px-3 capitalize">{r.payment_method.replace('_', ' ')}</td>
                    <td className="py-2 px-3 text-right font-medium">{formatUGX(r.amount)}</td>
                    <td className="py-2 px-3"><StatusPill status={r.status} size="sm" /></td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </DetailSection>

      <CollectCashDrawer open={collectOpen} onClose={() => setCollectOpen(false)} loanId={loan.id} suggested={loan.installment_amount} />
      <MomoDrawer open={momoOpen} onClose={() => setMomoOpen(false)} loanId={loan.id} suggestedAmount={loan.installment_amount} suggestedPhone={loan.borrower?.phone || ''} />
    </DetailPane>
  )
}

/* ─────────────── Cash repayment drawer ─────────────── */

const cashSchema = z.object({
  amount: z.number().positive('Required'),
  payment_method: z.enum(['cash', 'bank_transfer', 'cheque']),
  transaction_ref: z.string().optional(),
  notes: z.string().optional(),
})
type CashForm = z.infer<typeof cashSchema>

function CollectCashDrawer({ open, onClose, loanId, suggested }: { open: boolean; onClose: () => void; loanId: number; suggested: number }) {
  const create = useCreateRepayment()
  const { register, handleSubmit, formState: { errors }, reset, control } = useForm<CashForm>({
    resolver: zodResolver(cashSchema),
    defaultValues: { amount: Math.round(suggested), payment_method: 'cash', transaction_ref: '', notes: '' },
  })

  const onSubmit = handleSubmit(async (data) => {
    try {
      await create.mutateAsync({ loan_id: loanId, ...data })
      toast.success('Repayment recorded — pending verification')
      reset()
      onClose()
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  })

  return (
    <Drawer open={open} onOpenChange={(o) => { if (!o) reset(); onClose() }} title="Record Cash Repayment">
      <form onSubmit={onSubmit} className="space-y-5">
        <CurrencyField control={control} name="amount" label="Amount" required />
        <SelectField label="Method" required options={[
          { value: 'cash', label: 'Cash' },
          { value: 'bank_transfer', label: 'Bank transfer' },
          { value: 'cheque', label: 'Cheque' },
        ]} {...register('payment_method')} />
        <TextField label="Transaction reference (optional)" {...register('transaction_ref')} />
        <TextAreaField label="Notes" {...register('notes')} />
        <p className="text-[11.5px] text-foreground-muted">A manager must verify this repayment before it's applied to the loan balance.</p>
        <FormActions onCancel={onClose} submitLabel="Record" isPending={create.isPending} />
      </form>
    </Drawer>
  )
}

/* ─────────────── Mobile money drawer ─────────────── */

function MomoDrawer({ open, onClose, loanId, suggestedAmount, suggestedPhone }: { open: boolean; onClose: () => void; loanId: number; suggestedAmount: number; suggestedPhone: string }) {
  const collect = useDGatewayCollect()
  const [stage, setStage] = useState<'idle' | 'pending' | 'success' | 'failed'>('idle')
  const [failReason, setFailReason] = useState('')

  const { register, handleSubmit, control, reset } = useForm<{ amount: number; phone: string; provider: 'iotec' | 'relworx' }>({
    defaultValues: { amount: Math.round(suggestedAmount), phone: suggestedPhone, provider: 'iotec' },
  })

  const close = () => {
    if (stage === 'pending') return
    setStage('idle')
    setFailReason('')
    reset()
    onClose()
  }

  const onSubmit = handleSubmit(async (data) => {
    setStage('pending')
    try {
      const initiated = await collect.mutateAsync({ loan_id: loanId, ...data })
      toast.success('Prompt sent — borrower must approve on phone')
      const final = await pollDGatewayStatus(initiated.repayment_id)
      if (final.gateway_status === 'completed') setStage('success')
      else if (final.gateway_status === 'failed') {
        setStage('failed')
        setFailReason(final.failure_reason || 'Transaction failed')
      } else {
        setStage('failed')
        setFailReason('Borrower did not respond in time. Try again or use cash.')
      }
    } catch (err: any) {
      setStage('failed')
      setFailReason(err.response?.data?.error || 'Failed to initiate')
    }
  })

  return (
    <Drawer open={open} onOpenChange={(o) => !o && close()} title="Pay via Mobile Money">
      {stage === 'idle' && (
        <form onSubmit={onSubmit} className="space-y-5">
          <CurrencyField control={control} name="amount" label="Amount" required />
          <TextField label="Phone" required placeholder="0700000000 or 256700000000" {...register('phone')} />
          <SelectField label="Provider" options={[
            { value: 'iotec', label: 'Iotec (MTN + Airtel)' },
            { value: 'relworx', label: 'Relworx' },
          ]} {...register('provider')} />
          <p className="text-[11.5px] text-foreground-muted">Borrower will receive a prompt on their phone to approve the deduction.</p>
          <FormActions onCancel={close} submitLabel="Send Prompt" />
        </form>
      )}
      {stage === 'pending' && (
        <div className="py-12 text-center space-y-4">
          <Loader2 size={48} className="mx-auto animate-spin text-accent" />
          <div>
            <p className="font-semibold text-foreground">Waiting for borrower to approve</p>
            <p className="text-[13px] text-foreground-muted mt-1">Tell them to check their phone now.</p>
          </div>
          <p className="text-[11.5px] text-foreground-muted">This will time out after about 90 seconds.</p>
        </div>
      )}
      {stage === 'success' && (
        <div className="py-12 text-center space-y-4">
          <div className="w-16 h-16 rounded-full bg-success-light mx-auto flex items-center justify-center">
            <Check size={32} className="text-success-dark" />
          </div>
          <div>
            <p className="font-semibold text-foreground">Payment received</p>
            <p className="text-[13px] text-foreground-muted mt-1">Credited to the loan.</p>
          </div>
          <button onClick={close} className="h-9 px-4 rounded-lg bg-accent text-white text-[13px] font-medium">Done</button>
        </div>
      )}
      {stage === 'failed' && (
        <div className="py-12 text-center space-y-4">
          <div className="w-16 h-16 rounded-full bg-danger-light mx-auto flex items-center justify-center">
            <X size={32} className="text-danger-dark" />
          </div>
          <div>
            <p className="font-semibold text-foreground">Payment did not complete</p>
            <p className="text-[13px] text-foreground-muted mt-1">{failReason}</p>
          </div>
          <div className="flex gap-2 justify-center">
            <button onClick={() => { setStage('idle'); setFailReason('') }} className="h-9 px-4 rounded-lg bg-accent text-white text-[13px] font-medium">Try again</button>
            <button onClick={close} className="h-9 px-4 rounded-lg border border-border bg-surface text-[13px] font-medium">Close</button>
          </div>
        </div>
      )}
    </Drawer>
  )
}
