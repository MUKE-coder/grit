import { createFileRoute, useNavigate, Link } from '@tanstack/react-router'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Bike, Pencil, ArrowRight, Wallet, FileText, MapPin } from 'lucide-react'
import toast from 'react-hot-toast'

import { useMotorcycles, useMotorcycle, useCreateMotorcycle, useUpdateMotorcycle, useTransferMotorcycle } from '@/hooks/useMotorcycles'
import { useBranches } from '@/hooks/useBusiness'
import {
  TwoPane, ListPane, ListRow, DetailPane, DetailSection, Field, FieldGrid,
} from '@/components/two-pane'
import { Drawer } from '@/components/drawer'
import { TextField, SelectField, FormGrid, FormActions, FormSection, TextAreaField } from '@/components/form'
import { CurrencyField } from '@/components/currency-field'
import { StatusPill } from '@/components/status-badge'
import { ExportButton } from '@/components/export-button'
import { exportToExcel } from '@/lib/export'
import { formatUGX } from '@/lib/utils'
import type { MotorcycleStatus } from '@/types'

type MotoSearch = { selected?: number; q?: string; status?: MotorcycleStatus | 'all'; branch?: string }

export const Route = createFileRoute('/_app/motorcycles/')({
  component: MotorcyclesPage,
  validateSearch: (s: Record<string, unknown>): MotoSearch => ({
    selected: s.selected ? Number(s.selected) : undefined,
    q: typeof s.q === 'string' ? s.q : undefined,
    status: typeof s.status === 'string' ? (s.status as any) : undefined,
    branch: typeof s.branch === 'string' ? s.branch : undefined,
  }),
})

const STATUS_TABS: { val: MotorcycleStatus | 'all'; label: string }[] = [
  { val: 'all', label: 'All' },
  { val: 'available', label: 'Available' },
  { val: 'reserved', label: 'Reserved' },
  { val: 'on_loan', label: 'On Loan' },
  { val: 'sold', label: 'Sold' },
  { val: 'repossessed', label: 'Repossessed' },
]

function MotorcyclesPage() {
  const search = Route.useSearch()
  const navigate = useNavigate()
  const [createOpen, setCreateOpen] = useState(false)

  const list = useMotorcycles({
    search: search.q,
    status: search.status === 'all' ? undefined : search.status,
    branch_id: search.branch,
  })
  const selected = useMotorcycle(search.selected)

  const update = (patch: Partial<MotoSearch>) =>
    navigate({ to: '/motorcycles', search: { ...search, ...patch }, replace: true })

  return (
    <TwoPane>
      <ListPane
        title="Motorcycles"
        count={list.data?.total}
        search={search.q || ''}
        onSearch={(q) => update({ q: q || undefined })}
        searchPlaceholder="Name or number plate..."
        onNew={() => setCreateOpen(true)}
        refreshKey={['motorcycles']}
        extraAction={
          <ExportButton
            onClick={() => exportToExcel(
              (list.data?.data || []).map((m: any) => ({
                Name: m.name,
                'Number Plate': m.number_plate,
                Branch: m.branch?.name,
                Status: m.status,
                Color: m.color,
                'Year of Make': m.year_of_make,
                'Chassis No': m.chassis_no,
                'Engine No': m.engine_no,
                'Cost Price': m.cost_price,
                'Selling Price': m.selling_price,
                'Created At': m.created_at ? new Date(m.created_at).toLocaleDateString() : '',
              })),
              'motorcycles',
              'Motorcycles',
            )}
            disabled={!list.data?.data.length}
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
        {!list.isLoading && list.data?.data.length === 0 && (
          <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">No motorcycles match.</div>
        )}
        {list.data?.data.map((m) => (
          <ListRow
            key={m.id}
            selected={search.selected === m.id}
            onClick={() => update({ selected: m.id })}
            icon={<Bike size={14} className="text-foreground-secondary" />}
            title={m.name}
            subtitle={<span className="font-mono">{m.number_plate}</span>}
            rightTop={formatUGX(m.selling_price)}
            rightBottom={<StatusPill status={m.status} size="sm" />}
          />
        ))}
      </ListPane>

      {!search.selected ? (
        <DetailPane empty emptyTitle="No motorcycle selected" emptyHint="Pick one from the list, or click + New." />
      ) : selected.isLoading || !selected.data ? (
        <DetailPane><div className="text-foreground-muted">Loading…</div></DetailPane>
      ) : (
        <MotorcycleDetail m={selected.data} />
      )}

      <CreateMotorcycleDrawer
        open={createOpen}
        onClose={() => setCreateOpen(false)}
        onCreated={(id) => update({ selected: id })}
      />
    </TwoPane>
  )
}

function MotorcycleDetail({ m }: { m: NonNullable<ReturnType<typeof useMotorcycle>['data']> }) {
  const [editOpen, setEditOpen] = useState(false)
  const [transferOpen, setTransferOpen] = useState(false)
  const isAvailable = m.status === 'available'

  return (
    <DetailPane
      header={
        <div className="flex items-center justify-between gap-4">
          <div className="flex items-center gap-3 min-w-0">
            <div className="h-12 w-12 rounded-xl bg-surface-2 flex items-center justify-center shrink-0 overflow-hidden">
              {m.image_url ? <img src={m.image_url} alt={m.name} className="w-full h-full object-cover" /> : <Bike size={20} className="text-foreground-muted" />}
            </div>
            <div className="min-w-0">
              <h1 className="text-[17px] font-semibold text-foreground truncate">{m.name}</h1>
              <p className="text-[12.5px] text-foreground-muted font-mono">{m.number_plate} · {m.branch?.name || '—'}</p>
            </div>
          </div>
          <div className="flex items-center gap-2 shrink-0">
            <StatusPill status={m.status} />
            <button
              type="button"
              onClick={() => setEditOpen(true)}
              className="h-8 px-2.5 rounded-lg border border-border bg-surface text-[12px] font-medium text-foreground-secondary hover:bg-surface-hover transition inline-flex items-center gap-1"
            >
              <Pencil size={12} /> Edit
            </button>
          </div>
        </div>
      }
    >
      {/* Quick actions */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-3 mb-6">
        <ActionCard
          icon={<FileText size={14} />}
          label="Sell on Loan"
          to="/loans/new"
          search={{ motorcycle_id: String(m.id) }}
          disabled={!isAvailable}
        />
        <ActionCard
          icon={<Wallet size={14} />}
          label="Sell for Cash"
          to="/cash-sales/new"
          search={{ motorcycle_id: String(m.id) }}
          disabled={!isAvailable}
        />
        <button
          type="button"
          onClick={() => setTransferOpen(true)}
          disabled={!isAvailable}
          className="flex items-center gap-2 px-4 py-2.5 rounded-lg border border-border bg-surface text-[13px] font-medium text-foreground-secondary hover:bg-surface-hover transition disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <MapPin size={14} /> Transfer Branch
        </button>
      </div>

      <DetailSection title="Vehicle">
        <FieldGrid columns={2}>
          <Field label="Brand / Name" value={m.name} />
          <Field label="Number plate" value={<span className="font-mono">{m.number_plate}</span>} />
          <Field label="Color" value={m.color} />
          <Field label="Year of make" value={m.year_of_make ? String(m.year_of_make) : null} />
          <Field label="Chassis no." value={m.chassis_no ? <span className="font-mono">{m.chassis_no}</span> : null} />
          <Field label="Engine no." value={m.engine_no ? <span className="font-mono">{m.engine_no}</span> : null} />
        </FieldGrid>
      </DetailSection>

      <DetailSection title="Pricing & Branch">
        <FieldGrid columns={2}>
          <Field label="Branch" value={m.branch?.name} />
          <Field label="Cost price" value={m.cost_price > 0 ? formatUGX(m.cost_price) : null} />
          <Field label="Selling price" value={<span className="font-semibold text-foreground">{formatUGX(m.selling_price)}</span>} />
        </FieldGrid>
      </DetailSection>

      {m.notes && (
        <DetailSection title="Notes">
          <p className="text-[13.5px] text-foreground whitespace-pre-line">{m.notes}</p>
        </DetailSection>
      )}

      <EditMotorcycleDrawer open={editOpen} onClose={() => setEditOpen(false)} motorcycle={m} />
      <TransferDrawer open={transferOpen} onClose={() => setTransferOpen(false)} motorcycle={m} />
    </DetailPane>
  )
}

function ActionCard({ icon, label, to, search, disabled }: { icon: React.ReactNode; label: string; to: string; search?: Record<string, string>; disabled?: boolean }) {
  if (disabled) {
    return (
      <div className="flex items-center gap-2 px-4 py-2.5 rounded-lg border border-dashed border-border text-[13px] font-medium text-foreground-muted opacity-60">
        {icon} {label}
      </div>
    )
  }
  return (
    <Link to={to} search={search as any} className="flex items-center gap-2 px-4 py-2.5 rounded-lg bg-accent text-white text-[13px] font-medium hover:bg-accent-hover transition">
      {icon} {label} <ArrowRight size={12} className="ml-auto" />
    </Link>
  )
}

/* ─────────────── Create drawer ─────────────── */

const createSchema = z.object({
  branch_id: z.coerce.number().int().positive('Branch is required'),
  name: z.string().min(1, 'Required'),
  number_plate: z.string().min(1, 'Required'),
  selling_price: z.number({ message: 'Required' }).positive('Must be positive'),
})
type CreateMotoForm = z.infer<typeof createSchema>

function CreateMotorcycleDrawer({ open, onClose, onCreated }: { open: boolean; onClose: () => void; onCreated: (id: number) => void }) {
  const { data: branches } = useBranches()
  const create = useCreateMotorcycle()
  const { register, handleSubmit, formState: { errors }, reset, control } = useForm<CreateMotoForm>({
    resolver: zodResolver(createSchema),
    defaultValues: { branch_id: undefined as any, name: '', number_plate: '', selling_price: undefined },
  })

  const onSubmit = handleSubmit(async (data) => {
    try {
      const created = await create.mutateAsync({
        ...data,
        number_plate: data.number_plate.toUpperCase().replace(/\s/g, ''),
      })
      toast.success('Motorcycle added')
      reset()
      onClose()
      onCreated(created.id)
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  })

  return (
    <Drawer
      open={open}
      onOpenChange={(o) => { if (!o) reset(); onClose() }}
      title="New Motorcycle"
      description="Just the basics — fill in chassis/engine numbers and notes after creating."
    >
      <form onSubmit={onSubmit} className="space-y-5">
        <SelectField
          label="Branch"
          required
          placeholder="Select branch"
          options={branches?.map((b) => ({ value: String(b.id), label: b.name })) || []}
          error={errors.branch_id?.message}
          {...register('branch_id')}
        />
        <TextField label="Name / brand" required placeholder="e.g. KEVLA" error={errors.name?.message} {...register('name')} />
        <TextField
          label="Number plate"
          required
          placeholder="e.g. UMA842GV"
          error={errors.number_plate?.message}
          {...register('number_plate', { setValueAs: (v) => v?.toString().toUpperCase() })}
        />
        <CurrencyField control={control} name="selling_price" label="Selling price" required />

        <FormActions onCancel={onClose} submitLabel="Add Motorcycle" isPending={create.isPending} />
      </form>
    </Drawer>
  )
}

/* ─────────────── Edit drawer ─────────────── */

const editSchema = z.object({
  branch_id: z.coerce.number().int().positive(),
  name: z.string().min(1),
  chassis_no: z.string().optional(),
  engine_no: z.string().optional(),
  color: z.string().optional(),
  year_of_make: z.coerce.number().int().min(0).optional(),
  cost_price: z.number().optional(),
  selling_price: z.number().positive(),
  notes: z.string().optional(),
})
type EditMotoForm = z.infer<typeof editSchema>

function EditMotorcycleDrawer({ open, onClose, motorcycle }: { open: boolean; onClose: () => void; motorcycle: NonNullable<ReturnType<typeof useMotorcycle>['data']> }) {
  const { data: branches } = useBranches()
  const update = useUpdateMotorcycle(motorcycle.id)
  const { register, handleSubmit, formState: { errors }, control } = useForm<EditMotoForm>({
    resolver: zodResolver(editSchema),
    defaultValues: {
      branch_id: motorcycle.branch_id,
      name: motorcycle.name,
      chassis_no: motorcycle.chassis_no || '',
      engine_no: motorcycle.engine_no || '',
      color: motorcycle.color || '',
      year_of_make: motorcycle.year_of_make || undefined,
      cost_price: motorcycle.cost_price || undefined,
      selling_price: motorcycle.selling_price,
      notes: motorcycle.notes || '',
    },
  })

  const onSubmit = handleSubmit(async (data) => {
    try {
      await update.mutateAsync(data as any)
      toast.success('Updated')
      onClose()
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  })

  return (
    <Drawer open={open} onOpenChange={(o) => !o && onClose()} title="Edit Motorcycle" width={520}>
      <form onSubmit={onSubmit} className="space-y-6">
        <FormSection title="Vehicle">
          <SelectField
            label="Branch"
            required
            options={branches?.map((b) => ({ value: String(b.id), label: b.name })) || []}
            error={errors.branch_id?.message}
            {...register('branch_id')}
          />
          <TextField label="Name / brand" required error={errors.name?.message} {...register('name')} />
          <FormGrid columns={2}>
            <TextField label="Color" {...register('color')} />
            <TextField label="Year" type="number" {...register('year_of_make')} />
          </FormGrid>
          <FormGrid columns={2}>
            <TextField label="Chassis no." {...register('chassis_no')} />
            <TextField label="Engine no." {...register('engine_no')} />
          </FormGrid>
        </FormSection>

        <FormSection title="Pricing">
          <FormGrid columns={2}>
            <CurrencyField control={control} name="cost_price" label="Cost price" />
            <CurrencyField control={control} name="selling_price" label="Selling price" required />
          </FormGrid>
        </FormSection>

        <TextAreaField label="Notes" {...register('notes')} />

        <FormActions onCancel={onClose} submitLabel="Save changes" isPending={update.isPending} />
      </form>
    </Drawer>
  )
}

/* ─────────────── Transfer drawer ─────────────── */

function TransferDrawer({ open, onClose, motorcycle }: { open: boolean; onClose: () => void; motorcycle: NonNullable<ReturnType<typeof useMotorcycle>['data']> }) {
  const { data: branches } = useBranches()
  const transfer = useTransferMotorcycle()
  const { register, handleSubmit, formState: { errors } } = useForm<{ branch_id: number; note: string }>({
    defaultValues: { branch_id: motorcycle.branch_id, note: '' },
  })

  const onSubmit = handleSubmit(async (data) => {
    try {
      await transfer.mutateAsync({ id: motorcycle.id, branch_id: Number(data.branch_id), note: data.note })
      toast.success('Transferred')
      onClose()
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  })

  return (
    <Drawer open={open} onOpenChange={(o) => !o && onClose()} title="Transfer Motorcycle">
      <form onSubmit={onSubmit} className="space-y-5">
        <p className="text-[12.5px] text-foreground-muted">Currently at <span className="font-medium text-foreground">{motorcycle.branch?.name}</span>. Pick the destination branch.</p>
        <SelectField
          label="To branch"
          required
          options={branches?.filter((b) => b.id !== motorcycle.branch_id).map((b) => ({ value: String(b.id), label: b.name })) || []}
          error={errors.branch_id?.message as any}
          {...register('branch_id', { required: true })}
        />
        <TextAreaField label="Note (optional)" placeholder="Reason for transfer..." {...register('note')} />

        <FormActions onCancel={onClose} submitLabel="Transfer" isPending={transfer.isPending} />
      </form>
    </Drawer>
  )
}
