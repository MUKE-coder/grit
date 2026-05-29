import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Users } from 'lucide-react'
import toast from 'react-hot-toast'

import { useDailyBodaDrivers, useCreateDailyBodaDriver } from '@/hooks/useDailyBoda'
import { useBranches } from '@/hooks/useBusiness'
import {
  TwoPane, ListPane, ListRow, DetailPane, DetailSection, Field, FieldGrid,
} from '@/components/two-pane'
import { Drawer } from '@/components/drawer'
import { TextField, SelectField, FormActions } from '@/components/form'
import { ExportButton } from '@/components/export-button'
import { exportToExcel } from '@/lib/export'
import { CurrencyField } from '@/components/currency-field'
import { StatusBadge } from '@/components/status-badge'
import { formatUGX } from '@/lib/utils'

type DSearch = { selected?: number; q?: string }

export const Route = createFileRoute('/_app/daily-boda/drivers')({
  component: DriversPage,
  validateSearch: (s: Record<string, unknown>): DSearch => ({
    selected: s.selected ? Number(s.selected) : undefined,
    q: typeof s.q === 'string' ? s.q : undefined,
  }),
})

function DriversPage() {
  const search = Route.useSearch()
  const navigate = useNavigate()
  const [createOpen, setCreateOpen] = useState(false)

  const list = useDailyBodaDrivers({ search: search.q })

  const update = (patch: Partial<DSearch>) =>
    navigate({ to: '/daily-boda/drivers', search: { ...search, ...patch }, replace: true })

  const drivers = list.data || []
  const selected = drivers.find((d) => d.id === search.selected)

  return (
    <TwoPane>
      <ListPane
        title="Daily Boda Drivers"
        count={drivers.length}
        search={search.q || ''}
        onSearch={(q) => update({ q: q || undefined })}
        searchPlaceholder="Name, phone, NIN..."
        onNew={() => setCreateOpen(true)}
        refreshKey={['daily-boda-drivers']}
        extraAction={
          <ExportButton
            onClick={() => exportToExcel(
              drivers.map((d: any) => ({
                'Full Name': d.full_name,
                Phone: d.phone,
                'National ID': d.national_id,
                Branch: d.branch?.name,
                'Daily Rate': d.daily_rate,
                Active: d.is_active ? 'Yes' : 'No',
                Address: d.address,
              })),
              'daily-boda-drivers',
              'Drivers',
            )}
            disabled={drivers.length === 0}
          />
        }
      >
        {list.isLoading && <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">Loading…</div>}
        {!list.isLoading && drivers.length === 0 && (
          <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">No drivers yet.</div>
        )}
        {drivers.map((d) => (
          <ListRow
            key={d.id}
            selected={search.selected === d.id}
            onClick={() => update({ selected: d.id })}
            icon={<Users size={14} className="text-foreground-secondary" />}
            title={d.full_name}
            subtitle={<span className="font-mono">{d.phone}</span>}
            rightTop={formatUGX(d.daily_rate)}
            rightBottom={
              <StatusBadge tone={d.is_active ? 'success' : 'neutral'} size="sm">
                {d.is_active ? 'Active' : 'Inactive'}
              </StatusBadge>
            }
          />
        ))}
      </ListPane>

      {!selected ? (
        <DetailPane empty emptyTitle="No driver selected" emptyHint="Pick a driver, or click + New to register one." />
      ) : (
        <DriverDetail d={selected} />
      )}

      <CreateDriverDrawer open={createOpen} onClose={() => setCreateOpen(false)} onCreated={(id) => update({ selected: id })} />
    </TwoPane>
  )
}

function DriverDetail({ d }: { d: NonNullable<ReturnType<typeof useDailyBodaDrivers>['data']>[number] }) {
  return (
    <DetailPane
      header={
        <div className="flex items-center justify-between gap-4">
          <div className="min-w-0">
            <h1 className="text-[17px] font-semibold text-foreground truncate">{d.full_name}</h1>
            <p className="text-[12.5px] text-foreground-muted font-mono">{d.phone}</p>
          </div>
          <StatusBadge tone={d.is_active ? 'success' : 'neutral'}>
            {d.is_active ? 'Active' : 'Inactive'}
          </StatusBadge>
        </div>
      }
    >
      <DetailSection title="Profile">
        <FieldGrid columns={2}>
          <Field label="Phone" value={<span className="font-mono">{d.phone}</span>} />
          <Field label="National ID" value={d.national_id ? <span className="font-mono">{d.national_id}</span> : null} />
          <Field label="Branch" value={d.branch?.name} />
          <Field label="Daily rate" value={<span className="font-semibold">{formatUGX(d.daily_rate)}</span>} />
          {d.address && <Field label="Address" value={d.address} />}
        </FieldGrid>
      </DetailSection>
    </DetailPane>
  )
}

const driverSchema = z.object({
  branch_id: z.coerce.number().int().positive('Required'),
  full_name: z.string().min(1, 'Required'),
  phone: z.string().min(7, 'Required'),
  national_id: z.string().optional(),
  daily_rate: z.number({ message: 'Required' }).positive('Must be positive'),
})
type DriverForm = z.infer<typeof driverSchema>

function CreateDriverDrawer({ open, onClose, onCreated }: { open: boolean; onClose: () => void; onCreated: (id: number) => void }) {
  const { data: branches } = useBranches()
  const create = useCreateDailyBodaDriver()
  const { register, handleSubmit, formState: { errors }, reset, control } = useForm<DriverForm>({
    resolver: zodResolver(driverSchema),
    defaultValues: { branch_id: undefined as any, full_name: '', phone: '', national_id: '', daily_rate: 15000 },
  })

  const onSubmit = handleSubmit(async (data) => {
    try {
      const created = await create.mutateAsync(data as any)
      toast.success('Driver added')
      reset()
      onClose()
      onCreated(created.id)
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  })

  return (
    <Drawer open={open} onOpenChange={(o) => { if (!o) reset(); onClose() }} title="New Daily Boda Driver">
      <form onSubmit={onSubmit} className="space-y-5">
        <SelectField label="Branch" required placeholder="Select branch" options={branches?.map((b) => ({ value: String(b.id), label: b.name })) || []} error={errors.branch_id?.message} {...register('branch_id')} />
        <TextField label="Full name" required error={errors.full_name?.message} {...register('full_name')} />
        <TextField label="Phone" required placeholder="0700000000" error={errors.phone?.message} {...register('phone')} />
        <TextField label="National ID" {...register('national_id')} />
        <CurrencyField control={control} name="daily_rate" label="Daily rate" required hint="Default 15,000 UGX." />
        <FormActions onCancel={onClose} submitLabel="Add Driver" isPending={create.isPending} />
      </form>
    </Drawer>
  )
}
