import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState, useMemo } from 'react'
import { Receipt, Undo2, AlertTriangle, Printer, Download } from 'lucide-react'
import toast from 'react-hot-toast'
import { useSales, useSale, useSaleReturns, useCreateReturn } from '@/hooks/useSales'
import {
  TwoPane, ListPane, ListRow, DetailPane, DetailSection, Field, FieldGrid,
} from '@/components/two-pane'
import { Drawer } from '@/components/drawer'
import { FormActions } from '@/components/form'
import { StatusBadge } from '@/components/status-badge'
import { ExportButton } from '@/components/export-button'
import { exportToExcel } from '@/lib/export'
import { formatUGX } from '@/lib/utils'
import { ReceiptPDF } from '@/components/pos/ReceiptPDF'
import { useReceiptPDF } from '@/components/pos/useReceiptPDF'

type SSearch = { selected?: number; q?: string; from?: string; to?: string }

export const Route = createFileRoute('/_app/sales/')({
  component: SalesPage,
  validateSearch: (s: Record<string, unknown>): SSearch => ({
    selected: s.selected ? Number(s.selected) : undefined,
    q: typeof s.q === 'string' ? s.q : undefined,
    from: typeof s.from === 'string' ? s.from : undefined,
    to: typeof s.to === 'string' ? s.to : undefined,
  }),
})

function SalesPage() {
  const search = Route.useSearch()
  const navigate = useNavigate()

  // Server filters by date when from/to are present (backend already accepts
  // from_date / to_date params on /sales).
  const list = useSales({
    from_date: search.from,
    to_date: search.to,
  } as any)
  const selected = useSale(search.selected || 0)

  const update = (patch: Partial<SSearch>) =>
    navigate({ to: '/sales', search: { ...search, ...patch }, replace: true })

  const items = (list.data?.data as any[] | undefined) || []
  const filtered = items.filter((s) => {
    if (!search.q) return true
    const q = search.q.toLowerCase()
    return (
      String(s.id).includes(q) ||
      s.customer_phone?.toLowerCase().includes(q) ||
      s.cashier_name?.toLowerCase().includes(q)
    )
  })

  return (
    <TwoPane>
      <ListPane
        title="Spares Sales"
        count={list.data?.total}
        search={search.q || ''}
        onSearch={(q) => update({ q: q || undefined })}
        searchPlaceholder="Sale #, phone, cashier..."
        newLabel="New"
        onNew={() => navigate({ to: '/pos' })}
        refreshKey={['sales']}
        filters={
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
        }
        extraAction={
          <ExportButton
            onClick={() => exportToExcel(
              filtered.map((s: any) => ({
                'Sale #': s.id,
                Date: new Date(s.created_at).toLocaleString(),
                Cashier: s.cashier_name,
                Branch: s.branch_name,
                'Customer Phone': s.customer_phone,
                'Payment Method': s.payment_method,
                Subtotal: s.subtotal,
                Discount: s.discount_amount,
                Total: s.total,
              })),
              'sales-history',
              'Sales',
            )}
            disabled={filtered.length === 0}
          />
        }
      >
        {list.isLoading && <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">Loading…</div>}
        {!list.isLoading && filtered.length === 0 && (
          <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">No sales yet.</div>
        )}
        {filtered.map((s: any) => (
          <ListRow
            key={s.id}
            selected={search.selected === s.id}
            onClick={() => update({ selected: s.id })}
            icon={<Receipt size={14} className="text-foreground-secondary" />}
            title={`Sale #${s.id}`}
            subtitle={
              <span>
                {s.cashier_name || `Cashier #${s.cashier_id}`}
                {s.branch_name && <> · {s.branch_name}</>}
              </span>
            }
            rightTop={new Date(s.created_at).toLocaleDateString()}
            rightBottom={<span className="text-[13px] font-semibold tabular-nums text-foreground">{formatUGX(s.total)}</span>}
          />
        ))}
      </ListPane>

      {!search.selected ? (
        <DetailPane empty emptyTitle="No sale selected" emptyHint="Pick a sale from the list to view its receipt." />
      ) : selected.isLoading || !selected.data ? (
        <DetailPane><div className="text-foreground-muted">Loading…</div></DetailPane>
      ) : (
        <SaleDetail s={selected.data} />
      )}
    </TwoPane>
  )
}

function SaleDetail({ s }: { s: any }) {
  const returns = useSaleReturns(s.id)
  const [returnOpen, setReturnOpen] = useState(false)

  // Compute already-returned quantity per sale_item so we can display
  // remaining-returnable + restrict the input.
  const returnedByItem = useMemo(() => {
    const m = new Map<number, number>()
    for (const r of returns.data || []) {
      for (const ri of r.items) {
        m.set(ri.sale_item_id, (m.get(ri.sale_item_id) || 0) + ri.quantity)
      }
    }
    return m
  }, [returns.data])

  const totalRefunded = (returns.data || []).reduce((sum, r) => sum + r.refunded_total, 0)
  const allReturned = (s.items || []).every((it: any) =>
    (returnedByItem.get(it.id) || 0) >= it.quantity,
  )

  return (
    <DetailPane
      header={
        <div className="flex items-start justify-between gap-4">
          <div>
            <p className="text-[10.5px] uppercase tracking-wider text-foreground-muted font-mono">Sale #{s.id}</p>
            <h1 className="text-[20px] font-semibold text-foreground tabular-nums mt-0.5">{formatUGX(s.total)}</h1>
            <p className="text-[12.5px] text-foreground-muted mt-0.5">{new Date(s.created_at).toLocaleString()} · {s.branch_name || '—'}</p>
            {totalRefunded > 0 && (
              <p className="text-[12px] text-danger-dark mt-1 font-medium">
                {formatUGX(totalRefunded)} refunded across {(returns.data || []).length} return{(returns.data || []).length === 1 ? '' : 's'}
              </p>
            )}
          </div>
          <div className="flex flex-col items-end gap-2">
            <StatusBadge tone={s.payment_method === 'cash' ? 'success' : 'info'}>
              {s.payment_method === 'cash' ? 'Cash' : 'Mobile Money'}
            </StatusBadge>
            <div className="flex items-center gap-1.5 flex-wrap justify-end">
              <ReceiptActions sale={s} />
              {!allReturned && (
                <button
                  type="button"
                  onClick={() => setReturnOpen(true)}
                  className="h-8 px-2.5 rounded-lg border border-border bg-surface text-[12px] font-medium text-foreground-secondary hover:bg-surface-hover transition inline-flex items-center gap-1"
                >
                  <Undo2 size={12} /> Return items
                </button>
              )}
            </div>
          </div>
        </div>
      }
    >
      <DetailSection title="Sale">
        <FieldGrid columns={2}>
          <Field label="Cashier" value={s.cashier_name || `User #${s.cashier_id}`} />
          <Field label="Branch" value={s.branch_name} />
          <Field label="Customer phone" value={s.customer_phone ? <span className="font-mono">{s.customer_phone}</span> : null} />
        </FieldGrid>
      </DetailSection>

      <DetailSection title={`Items (${s.items?.length || 0})`}>
        <div className="border border-border rounded-lg overflow-hidden">
          <table className="w-full text-[13px] tabular-nums">
            <thead className="bg-surface-2 text-[11px] uppercase tracking-wider text-foreground-muted">
              <tr>
                <th className="text-left py-2 px-3 font-semibold">Product</th>
                <th className="text-right py-2 px-3 font-semibold">Qty</th>
                <th className="text-right py-2 px-3 font-semibold">Unit</th>
                <th className="text-right py-2 px-3 font-semibold">Total</th>
              </tr>
            </thead>
            <tbody>
              {s.items?.map((it: any) => (
                <tr key={it.id} className="border-t border-border-subtle">
                  <td className="py-2 px-3">{it.product_name || it.product?.title || `Product #${it.product_id}`}</td>
                  <td className="py-2 px-3 text-right">{it.quantity}</td>
                  <td className="py-2 px-3 text-right">{formatUGX(it.unit_price)}</td>
                  <td className="py-2 px-3 text-right font-medium">{formatUGX(it.line_total)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </DetailSection>

      <DetailSection title="Totals">
        <div className="border border-border rounded-lg overflow-hidden bg-surface">
          <div className="px-4 py-2.5 flex justify-between text-[13px] border-b border-border-subtle">
            <span className="text-foreground-muted">Subtotal</span>
            <span className="font-medium tabular-nums">{formatUGX(s.subtotal)}</span>
          </div>
          {s.discount_amount > 0 && (
            <div className="px-4 py-2.5 flex justify-between text-[13px] border-b border-border-subtle">
              <span className="text-foreground-muted">Discount</span>
              <span className="font-medium tabular-nums text-danger-dark">−{formatUGX(s.discount_amount)}</span>
            </div>
          )}
          <div className="px-4 py-3 flex justify-between text-[14px] font-semibold bg-accent-light">
            <span>Total</span>
            <span className="tabular-nums text-accent-hover">{formatUGX(s.total)}</span>
          </div>
          {totalRefunded > 0 && (
            <>
              <div className="px-4 py-2.5 flex justify-between text-[13px] border-t border-border-subtle">
                <span className="text-foreground-muted">Refunded</span>
                <span className="font-medium tabular-nums text-danger-dark">−{formatUGX(totalRefunded)}</span>
              </div>
              <div className="px-4 py-2.5 flex justify-between text-[13px] border-t border-border-subtle bg-surface-2">
                <span className="font-medium text-foreground">Net Total</span>
                <span className="font-semibold tabular-nums">{formatUGX(s.total - totalRefunded)}</span>
              </div>
            </>
          )}
        </div>
      </DetailSection>

      {/* Returns history */}
      {(returns.data?.length || 0) > 0 && (
        <DetailSection title={`Returns (${returns.data?.length})`}>
          <div className="space-y-2">
            {returns.data!.map((r) => (
              <div key={r.id} className="border border-border rounded-lg p-3">
                <div className="flex items-center justify-between text-[12.5px] mb-2">
                  <span className="text-foreground-muted">
                    {new Date(r.created_at).toLocaleString()} · {r.processor?.name || 'Staff'}
                  </span>
                  <span className="font-semibold tabular-nums text-danger-dark">−{formatUGX(r.refunded_total)}</span>
                </div>
                <ul className="space-y-1 text-[12.5px]">
                  {r.items.map((ri) => (
                    <li key={ri.id} className="flex justify-between text-foreground-secondary">
                      <span>{ri.product?.title || `Product #${ri.product_id}`} × {ri.quantity}</span>
                      <span className="tabular-nums">{formatUGX(ri.line_total)}</span>
                    </li>
                  ))}
                </ul>
                {r.reason && <p className="text-[11.5px] text-foreground-muted mt-2 italic">"{r.reason}"</p>}
              </div>
            ))}
          </div>
        </DetailSection>
      )}

      <ReturnDrawer
        open={returnOpen}
        onClose={() => setReturnOpen(false)}
        sale={s}
        returnedByItem={returnedByItem}
      />
    </DetailPane>
  )
}

/* ─────────────── Return drawer ─────────────── */

function ReturnDrawer({
  open, onClose, sale, returnedByItem,
}: {
  open: boolean
  onClose: () => void
  sale: any
  returnedByItem: Map<number, number>
}) {
  const create = useCreateReturn()
  const [quantities, setQuantities] = useState<Record<number, number>>({})
  const [reason, setReason] = useState('')

  const close = () => {
    if (create.isPending) return
    setQuantities({})
    setReason('')
    onClose()
  }

  // Lines the customer is bringing back (qty > 0).
  const planned = (sale.items || []).map((it: any) => {
    const remaining = it.quantity - (returnedByItem.get(it.id) || 0)
    const qty = quantities[it.id] || 0
    return { item: it, remaining, qty }
  })
  const totalRefund = planned.reduce((sum, p) => sum + p.qty * p.item.unit_price, 0)
  const anyPicked = planned.some((p) => p.qty > 0)

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    const items = planned
      .filter((p) => p.qty > 0)
      .map((p) => ({ sale_item_id: p.item.id, quantity: p.qty }))
    if (items.length === 0) {
      toast.error('Pick at least one item to return')
      return
    }
    try {
      await create.mutateAsync({ saleId: sale.id, items, reason })
      toast.success('Return recorded — stock restored')
      setQuantities({})
      setReason('')
      onClose()
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed to record return')
    }
  }

  return (
    <Drawer open={open} onOpenChange={(o) => !o && close()} title="Return items" description={`Sale #${sale.id} · pick which items the customer is bringing back`} width={520}>
      <form onSubmit={submit} className="space-y-5">
        <div className="border border-border rounded-lg overflow-hidden">
          <table className="w-full text-[13px]">
            <thead className="bg-surface-2 text-[11px] uppercase tracking-wider text-foreground-muted">
              <tr>
                <th className="text-left py-2 px-3 font-semibold">Product</th>
                <th className="text-right py-2 px-3 font-semibold">Sold</th>
                <th className="text-right py-2 px-3 font-semibold">Returnable</th>
                <th className="text-right py-2 px-3 font-semibold w-24">Return qty</th>
              </tr>
            </thead>
            <tbody>
              {planned.map(({ item, remaining }) => {
                const fullyReturned = remaining <= 0
                return (
                  <tr key={item.id} className="border-t border-border-subtle">
                    <td className="py-2 px-3">
                      <p className="font-medium">{item.product_name || item.product?.title || `Product #${item.product_id}`}</p>
                      <p className="text-[11px] text-foreground-muted tabular-nums">{formatUGX(item.unit_price)} each</p>
                    </td>
                    <td className="py-2 px-3 text-right tabular-nums">{item.quantity}</td>
                    <td className="py-2 px-3 text-right tabular-nums">
                      {fullyReturned
                        ? <span className="text-foreground-muted">0</span>
                        : <span className="text-foreground">{remaining}</span>}
                    </td>
                    <td className="py-1 px-2">
                      <input
                        type="number"
                        min={0}
                        max={remaining}
                        disabled={fullyReturned}
                        value={quantities[item.id] || ''}
                        onChange={(e) => {
                          const v = Math.min(remaining, Math.max(0, Number(e.target.value) || 0))
                          setQuantities((q) => ({ ...q, [item.id]: v }))
                        }}
                        className="w-full h-8 px-2 rounded-md border border-border bg-surface text-[13px] text-right tabular-nums focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/15 disabled:bg-surface-2 disabled:text-foreground-muted"
                      />
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>

        <div>
          <label className="block text-[12.5px] font-medium text-foreground-secondary mb-1.5">Reason (optional)</label>
          <textarea
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            rows={3}
            placeholder="e.g. defective, wrong size, customer changed mind..."
            className="w-full px-3 py-2 rounded-lg border border-border bg-surface text-[13.5px] focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/15"
          />
        </div>

        {anyPicked && (
          <div className="flex items-center justify-between bg-accent-light px-3 py-2.5 rounded-lg">
            <span className="text-[12.5px] font-medium text-accent-hover">Refund total</span>
            <span className="text-[15px] font-semibold tabular-nums text-accent-hover">{formatUGX(totalRefund)}</span>
          </div>
        )}

        <div className="flex items-start gap-2 px-3 py-2.5 rounded-lg bg-warning-light text-warning-dark text-[12px]">
          <AlertTriangle size={14} className="shrink-0 mt-0.5" />
          <span>Stock will be restored to the sale's branch. Reports will reflect the refund as a deduction from this sale's revenue.</span>
        </div>

        <FormActions onCancel={close} submitLabel="Record return" isPending={create.isPending} />
      </form>
    </Drawer>
  )
}

function ReceiptActions({ sale }: { sale: any }) {
  const { print, download, busy } = useReceiptPDF()
  const filename = `grit-motors-receipt-${String(sale.id).padStart(6, '0')}.pdf`
  return (
    <>
      <button
        type="button"
        onClick={() => print(<ReceiptPDF sale={sale} />)}
        disabled={busy === 'print'}
        className="h-8 px-2.5 rounded-lg border border-border bg-surface text-[12px] font-medium text-foreground-secondary hover:bg-surface-hover transition inline-flex items-center gap-1 disabled:opacity-50"
      >
        <Printer size={12} /> {busy === 'print' ? 'Printing…' : 'Print'}
      </button>
      <button
        type="button"
        onClick={() => download(<ReceiptPDF sale={sale} />, filename)}
        disabled={busy === 'download'}
        className="h-8 px-2.5 rounded-lg border border-border bg-surface text-[12px] font-medium text-foreground-secondary hover:bg-surface-hover transition inline-flex items-center gap-1 disabled:opacity-50"
      >
        <Download size={12} /> {busy === 'download' ? 'Saving…' : 'PDF'}
      </button>
    </>
  )
}
