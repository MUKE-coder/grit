import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Bike, UserPlus, RotateCcw } from 'lucide-react'
import toast from 'react-hot-toast'

import {
  useDailyBodaMotorcycles, useCreateDailyBodaMotorcycle,
  useDailyBodaDrivers, useAssignDriver, useReturnMotorcycle,
} from '@/hooks/useDailyBoda'
import { useBranches } from '@/hooks/useBusiness'
import {
  TwoPane, ListPane, ListRow, DetailPane, DetailSection, Field, FieldGrid,
} from '@/components/two-pane'
import { Drawer } from '@/components/drawer'
import { TextField, SelectField, FormActions } from '@/components/form'
import { StatusPill } from '@/components/status-badge'
import { ExportButton } from '@/components/export-button'
import { exportToExcel } from '@/lib/export'

type DBMSearch = { selected?: number; q?: string; status?: string }

export const Route = createFileRoute('/_app/daily-boda/motorcycles')({
  component: DBMPage,
  validateSearch: (s: Record<string, unknown>): DBMSearch => ({
    selected: s.selected ? Number(s.selected) : undefined,
    q: typeof s.q === 'string' ? s.q : undefined,
    status: typeof s.status === 'string' ? s.status : undefined,
  }),
})

const TABS = [
  { val: 'all', label: 'All' },
  { val: 'available', label: 'Available' },
  { val: 'occupied', label: 'Occupied' },
  { val: 'returned', label: 'Returned' },
  { val: 'in_service', label: 'In Service' },
]

function DBMPage() {
  const search = Route.useSearch()
  const navigate = useNavigate()
  const [createOpen, setCreateOpen] = useState(false)
  const [assignOpen, setAssignOpen] = useState(false)

  const list = useDailyBodaMotorcycles({ status: search.status === 'all' ? undefined : search.status })
  const drivers = useDailyBodaDrivers({ active: 'true' })
  const returnMoto = useReturnMotorcycle()

  const update = (patch: Partial<DBMSearch>) =>
    navigate({ to: '/daily-boda/motorcycles', search: { ...search, ...patch }, replace: true })

  const motos = (list.data || []).filter((m) => {
    if (!search.q) return true
    const q = search.q.toLowerCase()
    return m.name.toLowerCase().includes(q) || m.number_plate.toLowerCase().includes(q)
  })
  const selected = motos.find((m) => m.id === search.selected) || (list.data || []).find((m) => m.id === search.selected)
  const driverFor = (id: number | null) => drivers.data?.find((d) => d.id === id)

  const handleReturn = async () => {
    if (!selected) return
    if (!confirm('Return this motorcycle? It will be available for re-assignment.')) return
    try {
      await returnMoto.mutateAsync(selected.id)
      toast.success('Returned')
    } catch (e: any) {
      toast.error(e.response?.data?.error || 'Failed')
    }
  }

  return (
    <TwoPane>
      <ListPane
        title="Daily Boda Motorcycles"
        count={motos.length}
        search={search.q || ''}
        onSearch={(q) => update({ q: q || undefined })}
        searchPlaceholder="Name or number plate..."
        onNew={() => setCreateOpen(true)}
        refreshKey={['daily-boda-motorcycles']}
        extraAction={
          <ExportButton
            onClick={() => exportToExcel(
              motos.map((m: any) => ({
                Name: m.name,
                'Number Plate': m.number_plate,
                Branch: m.branch?.name,
                Status: m.status,
                'Assigned Driver': driverFor(m.assigned_driver_id)?.full_name || '',
                'Driver Phone': driverFor(m.assigned_driver_id)?.phone || '',
              })),
              'daily-boda-fleet',
              'Fleet',
            )}
            disabled={motos.length === 0}
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
        {!list.isLoading && motos.length === 0 && (
          <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">No motorcycles in fleet.</div>
        )}
        {motos.map((m) => {
          const driver = driverFor(m.assigned_driver_id)
          return (
            <ListRow
              key={m.id}
              selected={search.selected === m.id}
              onClick={() => update({ selected: m.id })}
              icon={<Bike size={14} className="text-foreground-secondary" />}
              title={m.name}
              subtitle={
                <span>
                  <span className="font-mono">{m.number_plate}</span>
                  {driver && <> · {driver.full_name}</>}
                </span>
              }
              rightTop={m.branch?.name}
              rightBottom={<StatusPill status={m.status} size="sm" />}
            />
          )
        })}
      </ListPane>

      {!selected ? (
        <DetailPane empty emptyTitle="No motorcycle selected" emptyHint="Pick one from the list." />
      ) : (
        <DetailPane
          header={
            <div className="flex items-center justify-between gap-4">
              <div>
                <h1 className="text-[17px] font-semibold text-foreground">{selected.name}</h1>
                <p className="text-[12.5px] text-foreground-muted font-mono">{selected.number_plate}</p>
              </div>
              <div className="flex gap-2 items-center">
                <StatusPill status={selected.status} />
                {selected.status === 'available' && (
                  <button onClick={() => setAssignOpen(true)} className="h-8 px-3 rounded-lg bg-accent text-white text-[12px] font-medium hover:bg-accent-hover transition inline-flex items-center gap-1">
                    <UserPlus size={12} /> Assign
                  </button>
                )}
                {selected.status === 'occupied' && (
                  <button onClick={handleReturn} className="h-8 px-3 rounded-lg bg-warning-light text-warning-dark text-[12px] font-medium hover:bg-warning/20 transition inline-flex items-center gap-1">
                    <RotateCcw size={12} /> Return
                  </button>
                )}
              </div>
            </div>
          }
        >
          <DetailSection title="Motorcycle">
            <FieldGrid columns={2}>
              <Field label="Name" value={selected.name} />
              <Field label="Number plate" value={<span className="font-mono">{selected.number_plate}</span>} />
              <Field label="Branch" value={selected.branch?.name} />
              <Field label="Assigned driver" value={driverFor(selected.assigned_driver_id)?.full_name} />
            </FieldGrid>
          </DetailSection>

          <AssignDriverDrawer
            open={assignOpen}
            onClose={() => setAssignOpen(false)}
            motorcycleId={selected.id}
            drivers={drivers.data || []}
          />
        </DetailPane>
      )}

      <CreateMotorcycleDrawer open={createOpen} onClose={() => setCreateOpen(false)} onCreated={(id) => update({ selected: id })} />
    </TwoPane>
  )
}

const motoSchema = z.object({
  branch_id: z.coerce.number().int().positive('Required'),
  name: z.string().min(1, 'Required'),
  number_plate: z.string().min(1, 'Required'),
})
type MotoForm = z.infer<typeof motoSchema>

function CreateMotorcycleDrawer({ open, onClose, onCreated }: { open: boolean; onClose: () => void; onCreated: (id: number) => void }) {
  const { data: branches } = useBranches()
  const create = useCreateDailyBodaMotorcycle()
  const { register, handleSubmit, formState: { errors }, reset } = useForm<MotoForm>({
    resolver: zodResolver(motoSchema),
    defaultValues: { branch_id: undefined as any, name: '', number_plate: '' },
  })

  const onSubmit = handleSubmit(async (data) => {
    try {
      const created = await create.mutateAsync({ ...data, number_plate: data.number_plate.toUpperCase().replace(/\s/g, '') })
      toast.success('Motorcycle added to fleet')
      reset()
      onClose()
      onCreated(created.id)
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  })

  return (
    <Drawer open={open} onOpenChange={(o) => { if (!o) reset(); onClose() }} title="New Daily Boda Motorcycle">
      <form onSubmit={onSubmit} className="space-y-5">
        <SelectField label="Branch" required options={branches?.map((b) => ({ value: String(b.id), label: b.name })) || []} error={errors.branch_id?.message} {...register('branch_id')} />
        <TextField label="Name" required placeholder="e.g. Bajaj Boxer" error={errors.name?.message} {...register('name')} />
        <TextField label="Number plate" required error={errors.number_plate?.message} {...register('number_plate', { setValueAs: (v) => v?.toString().toUpperCase() })} />
        <FormActions onCancel={onClose} submitLabel="Add to Fleet" isPending={create.isPending} />
      </form>
    </Drawer>
  )
}

function AssignDriverDrawer({ open, onClose, motorcycleId, drivers }: { open: boolean; onClose: () => void; motorcycleId: number; drivers: any[] }) {
  const assign = useAssignDriver()
  const { register, handleSubmit, formState: { errors }, reset } = useForm<{ driver_id: number }>()

  const onSubmit = handleSubmit(async (data) => {
    try {
      await assign.mutateAsync({ motorcycleId, driverId: Number(data.driver_id) })
      toast.success('Driver assigned')
      reset()
      onClose()
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  })

  return (
    <Drawer open={open} onOpenChange={(o) => !o && onClose()} title="Assign Driver">
      <form onSubmit={onSubmit} className="space-y-5">
        <SelectField
          label="Driver"
          required
          placeholder="Select driver"
          options={drivers.filter((d) => d.is_active).map((d) => ({ value: String(d.id), label: `${d.full_name} · ${d.phone}` }))}
          error={errors.driver_id?.message as any}
          {...register('driver_id', { required: true })}
        />
        <FormActions onCancel={onClose} submitLabel="Assign" isPending={assign.isPending} />
      </form>
    </Drawer>
  )
}
