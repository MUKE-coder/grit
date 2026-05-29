import { createFileRoute, useNavigate, Link } from '@tanstack/react-router'
import { Wallet } from 'lucide-react'
import { useCashSales, useCashSale } from '@/hooks/useCashSales'
import {
  TwoPane, ListPane, ListRow, DetailPane, DetailSection, Field, FieldGrid,
} from '@/components/two-pane'
import { StatusBadge } from '@/components/status-badge'
import { ExportButton } from '@/components/export-button'
import { exportToExcel } from '@/lib/export'
import { formatUGX } from '@/lib/utils'

type CSSearch = { selected?: number; q?: string }

export const Route = createFileRoute('/_app/cash-sales/')({
  component: CashSalesPage,
  validateSearch: (s: Record<string, unknown>): CSSearch => ({
    selected: s.selected ? Number(s.selected) : undefined,
    q: typeof s.q === 'string' ? s.q : undefined,
  }),
})

function CashSalesPage() {
  const search = Route.useSearch()
  const navigate = useNavigate()

  const list = useCashSales({ search: search.q })
  const selected = useCashSale(search.selected)

  const update = (patch: Partial<CSSearch>) =>
    navigate({ to: '/cash-sales', search: { ...search, ...patch }, replace: true })

  return (
    <TwoPane>
      <ListPane
        title="Cash Sales"
        count={list.data?.total}
        search={search.q || ''}
        onSearch={(q) => update({ q: q || undefined })}
        searchPlaceholder="Customer, plate, sale #..."
        newLabel="New Cash Sale"
        onNew={() => navigate({ to: '/cash-sales/new' })}
        refreshKey={['cash-sales']}
        extraAction={
          <ExportButton
            onClick={() => exportToExcel(
              (list.data?.data || []).map((s: any) => ({
                'Sale #': s.sale_number,
                Date: new Date(s.created_at).toLocaleString(),
                Customer: s.customer_name,
                'Customer Phone': s.customer_phone,
                'Customer NIN': s.customer_nin,
                Motorcycle: s.motorcycle?.name,
                'Number Plate': s.motorcycle?.number_plate,
                Branch: s.branch?.name,
                'List Price': s.list_price,
                Discount: s.discount_amount,
                Total: s.total,
                'Payment Method': s.payment_method,
                'Transaction Ref': s.transaction_ref,
              })),
              'cash-sales',
              'Cash Sales',
            )}
            disabled={!list.data?.data.length}
          />
        }
      >
        {list.isLoading && <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">Loading…</div>}
        {!list.isLoading && (list.data?.data.length ?? 0) === 0 && (
          <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">No cash sales yet.</div>
        )}
        {list.data?.data.map((s) => (
          <ListRow
            key={s.id}
            selected={search.selected === s.id}
            onClick={() => update({ selected: s.id })}
            icon={<Wallet size={14} className="text-foreground-secondary" />}
            title={s.customer_name}
            subtitle={
              <span>
                <span className="font-mono text-accent">{s.sale_number}</span>
                {s.motorcycle?.number_plate && <> · <span className="font-mono">{s.motorcycle.number_plate}</span></>}
              </span>
            }
            rightTop={new Date(s.created_at).toLocaleDateString()}
            rightBottom={<span className="text-[13px] font-semibold tabular-nums text-foreground">{formatUGX(s.total)}</span>}
          />
        ))}
      </ListPane>

      {!search.selected ? (
        <DetailPane empty emptyTitle="No cash sale selected" emptyHint="Pick a sale from the list to view the receipt." />
      ) : selected.isLoading || !selected.data ? (
        <DetailPane><div className="text-foreground-muted">Loading…</div></DetailPane>
      ) : (
        <CashSaleDetail s={selected.data} />
      )}
    </TwoPane>
  )
}

function CashSaleDetail({ s }: { s: NonNullable<ReturnType<typeof useCashSale>['data']> }) {
  return (
    <DetailPane
      header={
        <div className="flex items-start justify-between gap-4">
          <div>
            <p className="text-[10.5px] uppercase tracking-wider text-foreground-muted font-mono">{s.sale_number}</p>
            <h1 className="text-[17px] font-semibold text-foreground mt-0.5">{s.customer_name}</h1>
            <p className="text-[12.5px] text-foreground-muted mt-0.5">{new Date(s.created_at).toLocaleString()}</p>
          </div>
          <StatusBadge tone="success">Cash Sale</StatusBadge>
        </div>
      }
    >
      <DetailSection title="Motorcycle">
        <FieldGrid columns={2}>
          <Field label="Brand" value={
            <Link to="/motorcycles" search={{ selected: s.motorcycle_id }} className="text-accent hover:underline">
              {s.motorcycle?.name || `#${s.motorcycle_id}`}
            </Link>
          } />
          <Field label="Number plate" value={<span className="font-mono">{s.motorcycle?.number_plate}</span>} />
          <Field label="Branch" value={s.branch?.name} />
          <Field label="Sold by" value={s.seller?.name || `User #${s.sold_by}`} />
        </FieldGrid>
      </DetailSection>

      <DetailSection title="Customer">
        <FieldGrid columns={2}>
          <Field label="Name" value={s.customer_name} />
          <Field label="Phone" value={s.customer_phone ? <span className="font-mono">{s.customer_phone}</span> : null} />
          <Field label="Email" value={s.customer_email} />
          <Field label="NIN" value={s.customer_nin ? <span className="font-mono">{s.customer_nin}</span> : null} />
        </FieldGrid>
      </DetailSection>

      <DetailSection title="Payment">
        <div className="border border-border rounded-lg overflow-hidden bg-surface">
          <div className="px-4 py-3 flex justify-between text-[13px] border-b border-border-subtle">
            <span className="text-foreground-muted">List price</span>
            <span className="font-medium tabular-nums">{formatUGX(s.list_price)}</span>
          </div>
          {s.discount_amount > 0 && (
            <div className="px-4 py-3 flex justify-between text-[13px] border-b border-border-subtle">
              <span className="text-foreground-muted">Discount</span>
              <span className="font-medium tabular-nums text-danger-dark">−{formatUGX(s.discount_amount)}</span>
            </div>
          )}
          <div className="px-4 py-3 flex justify-between text-[14px] font-semibold border-b border-border-subtle bg-accent-light">
            <span>Total</span>
            <span className="tabular-nums text-accent-hover">{formatUGX(s.total)}</span>
          </div>
          <div className="px-4 py-3 flex justify-between text-[13px]">
            <span className="text-foreground-muted">Method</span>
            <span className="capitalize font-medium">{s.payment_method.replace('_', ' ')}</span>
          </div>
          {s.transaction_ref && (
            <div className="px-4 py-3 flex justify-between text-[12.5px] border-t border-border-subtle">
              <span className="text-foreground-muted">Reference</span>
              <span className="font-mono text-foreground-secondary">{s.transaction_ref}</span>
            </div>
          )}
        </div>
      </DetailSection>

      {s.notes && (
        <DetailSection title="Notes">
          <p className="text-[13px] text-foreground whitespace-pre-line">{s.notes}</p>
        </DetailSection>
      )}
    </DetailPane>
  )
}
