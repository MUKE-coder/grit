import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { CalendarClock, Check } from 'lucide-react'
import toast from 'react-hot-toast'

import {
  useDailyBodaPayments, useCreateDailyBodaPayment, useVerifyDailyBodaPayment,
  useDailyBodaDrivers, useDailyBodaMotorcycles,
} from '@/hooks/useDailyBoda'
import { useBranches } from '@/hooks/useBusiness'
import {
  TwoPane, ListPane, ListRow, DetailPane, DetailSection, Field, FieldGrid,
} from '@/components/two-pane'
import { Drawer } from '@/components/drawer'
import { SelectField, FormActions, FormGrid } from '@/components/form'
import { CurrencyField } from '@/components/currency-field'
import { StatusPill } from '@/components/status-badge'
import { ExportButton } from '@/components/export-button'
import { exportToExcel } from '@/lib/export'
import { formatUGX } from '@/lib/utils'

type PSearch = { selected?: number; q?: string; status?: string }

export const Route = createFileRoute('/_app/daily-boda/payments')({
  component: PaymentsPage,
  validateSearch: (s: Record<string, unknown>): PSearch => ({
    selected: s.selected ? Number(s.selected) : undefined,
    q: typeof s.q === 'string' ? s.q : undefined,
    status: typeof s.status === 'string' ? s.status : undefined,
  }),
})

const TABS = [
  { val: 'all', label: 'All' },
  { val: 'paid', label: 'Paid' },
  { val: 'partial', label: 'Partial' },
  { val: 'pending', label: 'Pending' },
]

function PaymentsPage() {
  const search = Route.useSearch()
  const navigate = useNavigate()
  const [createOpen, setCreateOpen] = useState(false)

  const list = useDailyBodaPayments({ status: search.status === 'all' ? undefined : search.status })
  const drivers = useDailyBodaDrivers()
  const motos = useDailyBodaMotorcycles()
  const verify = useVerifyDailyBodaPayment()

  const update = (patch: Partial<PSearch>) =>
    navigate({ to: '/daily-boda/payments', search: { ...search, ...patch }, replace: true })

  const driverFor = (id: number) => drivers.data?.find((d) => d.id === id)
  const motoFor = (id: number) => motos.data?.find((m) => m.id === id)

  const items = (list.data?.data || []).filter((p) => {
    if (!search.q) return true
    const q = search.q.toLowerCase()
    return (
      driverFor(p.driver_id)?.full_name?.toLowerCase().includes(q) ||
      driverFor(p.driver_id)?.phone?.toLowerCase().includes(q) ||
      motoFor(p.motorcycle_id)?.number_plate?.toLowerCase().includes(q)
    )
  })
  const selected = items.find((p) => p.id === search.selected) || (list.data?.data || []).find((p) => p.id === search.selected)

  const handleVerify = async () => {
    if (!selected) return
    try {
      await verify.mutateAsync(selected.id)
      toast.success('Verified')
    } catch (e: any) {
      toast.error(e.response?.data?.error || 'Failed')
    }
  }

  return (
    <TwoPane>
      <ListPane
        title="Daily Boda Payments"
        count={list.data?.total}
        search={search.q || ''}
        onSearch={(q) => update({ q: q || undefined })}
        searchPlaceholder="Driver name, phone, plate..."
        newLabel="Record"
        onNew={() => setCreateOpen(true)}
        refreshKey={['daily-boda-payments']}
        extraAction={
          <ExportButton
            onClick={() => exportToExcel(
              items.map((p: any) => ({
                Date: new Date(p.payment_date).toLocaleDateString(),
                Driver: driverFor(p.driver_id)?.full_name || `Driver #${p.driver_id}`,
                'Driver Phone': driverFor(p.driver_id)?.phone || '',
                'Number Plate': motoFor(p.motorcycle_id)?.number_plate || '',
                'Daily Rate': p.daily_rate,
                Amount: p.amount,
                Balance: p.balance,
                Status: p.status,
                Method: p.payment_method,
                Verified: p.verified_at ? new Date(p.verified_at).toLocaleString() : '',
              })),
              'daily-boda-payments',
              'Payments',
            )}
            disabled={items.length === 0}
          />
        }
        filters={
          <div className="flex gap-1 flex-wrap">
            {TABS.map((t) => (
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
        {!list.isLoading && items.length === 0 && (
          <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">No payments yet.</div>
        )}
        {items.map((p) => {
          const driver = driverFor(p.driver_id)
          const moto = motoFor(p.motorcycle_id)
          return (
            <ListRow
              key={p.id}
              selected={search.selected === p.id}
              onClick={() => update({ selected: p.id })}
              icon={<CalendarClock size={14} className="text-foreground-secondary" />}
              title={driver?.full_name || `Driver #${p.driver_id}`}
              subtitle={
                <span>
                  <span className="font-mono">{moto?.number_plate}</span>
                  {' · '}
                  {new Date(p.payment_date).toLocaleDateString()}
                </span>
              }
              rightTop={<span className="tabular-nums">{formatUGX(p.amount)}</span>}
              rightBottom={<StatusPill status={p.status} size="sm" />}
            />
          )
        })}
      </ListPane>

      {!selected ? (
        <DetailPane empty emptyTitle="No payment selected" emptyHint="Pick a collection to verify or inspect." />
      ) : (
        <DetailPane
          header={
            <div className="flex items-start justify-between gap-4">
              <div>
                <h1 className="text-[20px] font-semibold text-foreground tabular-nums">{formatUGX(selected.amount)}</h1>
                <p className="text-[12.5px] text-foreground-muted mt-0.5">
                  {driverFor(selected.driver_id)?.full_name} · {new Date(selected.payment_date).toLocaleString()}
                </p>
              </div>
              <div className="flex items-center gap-2 shrink-0">
                <StatusPill status={selected.status} />
                {!selected.verified_at && (
                  <button onClick={handleVerify} className="h-8 px-3 rounded-lg bg-accent text-white text-[12px] font-medium hover:bg-accent-hover transition inline-flex items-center gap-1">
                    <Check size={12} /> Verify
                  </button>
                )}
              </div>
            </div>
          }
        >
          <DetailSection title="Collection">
            <FieldGrid columns={2}>
              <Field label="Amount paid" value={formatUGX(selected.amount)} />
              <Field label="Daily rate" value={formatUGX(selected.daily_rate)} />
              <Field label="Balance unpaid" value={selected.balance > 0 ? <span className="text-warning-dark font-medium">{formatUGX(selected.balance)}</span> : '—'} />
              <Field label="Method" value={<span className="capitalize">{selected.payment_method?.replace('_', ' ') || '—'}</span>} />
            </FieldGrid>
          </DetailSection>
          <DetailSection title="Linked">
            <FieldGrid columns={2}>
              <Field label="Driver" value={driverFor(selected.driver_id)?.full_name} />
              <Field label="Phone" value={<span className="font-mono">{driverFor(selected.driver_id)?.phone}</span>} />
              <Field label="Motorcycle" value={motoFor(selected.motorcycle_id)?.name} />
              <Field label="Plate" value={<span className="font-mono">{motoFor(selected.motorcycle_id)?.number_plate}</span>} />
            </FieldGrid>
          </DetailSection>
          <DetailSection title="Audit">
            <FieldGrid columns={2}>
              <Field label="Verified at" value={selected.verified_at ? new Date(selected.verified_at).toLocaleString() : null} />
            </FieldGrid>
          </DetailSection>
          {selected.notes && (
            <DetailSection title="Notes">
              <p className="text-[13px] text-foreground whitespace-pre-line">{selected.notes}</p>
            </DetailSection>
          )}
        </DetailPane>
      )}

      <CreatePaymentDrawer open={createOpen} onClose={() => setCreateOpen(false)} drivers={drivers.data || []} motos={motos.data || []} />
    </TwoPane>
  )
}

const paySchema = z.object({
  driver_id: z.coerce.number().int().positive('Required'),
  motorcycle_id: z.coerce.number().int().positive('Required'),
  branch_id: z.coerce.number().int().positive('Required'),
  amount: z.number({ message: 'Required' }).positive(),
  daily_rate: z.number({ message: 'Required' }).positive(),
})
type PayForm = z.infer<typeof paySchema>

function CreatePaymentDrawer({ open, onClose, drivers, motos }: { open: boolean; onClose: () => void; drivers: any[]; motos: any[] }) {
  const { data: branches } = useBranches()
  const create = useCreateDailyBodaPayment()
  const { register, handleSubmit, formState: { errors }, reset, setValue, control } = useForm<PayForm>({
    resolver: zodResolver(paySchema),
    defaultValues: {
      driver_id: undefined as any, motorcycle_id: undefined as any, branch_id: undefined as any,
      amount: undefined as any, daily_rate: 15000,
    },
  })

  const onSubmit = handleSubmit(async (data) => {
    try {
      await create.mutateAsync(data as any)
      toast.success('Recorded')
      reset()
      onClose()
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  })

  const onPickDriver = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const driver = drivers.find((d) => String(d.id) === e.target.value)
    if (driver) {
      setValue('daily_rate', driver.daily_rate)
      setValue('branch_id', driver.branch_id)
    }
  }

  return (
    <Drawer open={open} onOpenChange={(o) => { if (!o) reset(); onClose() }} title="Record Daily Payment">
      <form onSubmit={onSubmit} className="space-y-5">
        <SelectField
          label="Driver"
          required
          placeholder="Select driver"
          options={drivers.map((d) => ({ value: String(d.id), label: `${d.full_name} · ${d.phone}` }))}
          error={errors.driver_id?.message}
          {...register('driver_id', { onChange: onPickDriver })}
        />
        <SelectField
          label="Motorcycle"
          required
          placeholder="Select motorcycle"
          options={motos.map((m) => ({ value: String(m.id), label: `${m.name} · ${m.number_plate}` }))}
          error={errors.motorcycle_id?.message}
          {...register('motorcycle_id')}
        />
        <SelectField
          label="Branch"
          required
          options={branches?.map((b) => ({ value: String(b.id), label: b.name })) || []}
          error={errors.branch_id?.message}
          {...register('branch_id')}
        />
        <FormGrid columns={2}>
          <CurrencyField control={control} name="daily_rate" label="Daily rate" required />
          <CurrencyField control={control} name="amount" label="Amount paid" required />
        </FormGrid>
        <FormActions onCancel={onClose} submitLabel="Record Payment" isPending={create.isPending} />
      </form>
    </Drawer>
  )
}
