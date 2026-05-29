import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { FileText, Pencil } from 'lucide-react'
import toast from 'react-hot-toast'

import { useLoanProducts, useLoanProduct, useCreateLoanProduct, useUpdateLoanProduct } from '@/hooks/useLoanProducts'
import {
  TwoPane, ListPane, ListRow, DetailPane, DetailSection, Field, FieldGrid,
} from '@/components/two-pane'
import { Drawer } from '@/components/drawer'
import { TextField, SelectField, FormGrid, FormActions, FormSection, TextAreaField } from '@/components/form'
import { CurrencyField } from '@/components/currency-field'
import { StatusBadge } from '@/components/status-badge'
import { ExportButton } from '@/components/export-button'
import { exportToExcel } from '@/lib/export'
import { formatUGX } from '@/lib/utils'

type LPSearch = { selected?: number; q?: string }

export const Route = createFileRoute('/_app/loan-products/')({
  component: LoanProductsPage,
  validateSearch: (s: Record<string, unknown>): LPSearch => ({
    selected: s.selected ? Number(s.selected) : undefined,
    q: typeof s.q === 'string' ? s.q : undefined,
  }),
})

function LoanProductsPage() {
  const search = Route.useSearch()
  const navigate = useNavigate()
  const [createOpen, setCreateOpen] = useState(false)

  const list = useLoanProducts()
  const selected = useLoanProduct(search.selected)

  const update = (patch: Partial<LPSearch>) =>
    navigate({ to: '/loan-products', search: { ...search, ...patch }, replace: true })

  const filtered = (list.data || []).filter((p) =>
    !search.q ? true : p.name.toLowerCase().includes(search.q.toLowerCase()),
  )

  return (
    <TwoPane>
      <ListPane
        title="Loan Products"
        count={filtered.length}
        search={search.q || ''}
        onSearch={(q) => update({ q: q || undefined })}
        searchPlaceholder="Search products..."
        onNew={() => setCreateOpen(true)}
        refreshKey={['loan-products']}
        extraAction={
          <ExportButton
            onClick={() => exportToExcel(
              filtered.map((p: any) => ({
                Name: p.name,
                'Min Amount': p.min_amount,
                'Max Amount': p.max_amount,
                'Min Duration': p.min_duration,
                'Max Duration': p.max_duration,
                'Repayment Cycle': p.repayment_cycle,
                'Interest Method': p.interest_method,
                'Interest Rate %': p.interest_rate,
                'Requires Collateral': p.requires_collateral ? 'Yes' : 'No',
                'Grace Period Days': p.grace_period_days,
                Active: p.is_active ? 'Yes' : 'No',
                Description: p.description,
              })),
              'loan-products',
              'Loan Products',
            )}
            disabled={filtered.length === 0}
          />
        }
      >
        {list.isLoading && <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">Loading…</div>}
        {!list.isLoading && filtered.length === 0 && (
          <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">No loan products yet.</div>
        )}
        {filtered.map((p) => (
          <ListRow
            key={p.id}
            selected={search.selected === p.id}
            onClick={() => update({ selected: p.id })}
            icon={<FileText size={14} className="text-foreground-secondary" />}
            title={p.name}
            subtitle={`${p.interest_rate}% ${p.interest_method.replace('_', ' ')} · ${p.repayment_cycle}`}
            rightTop={`${formatUGX(p.min_amount)} – ${formatUGX(p.max_amount)}`}
            rightBottom={!p.is_active ? <StatusBadge tone="neutral" size="sm">Inactive</StatusBadge> : null}
          />
        ))}
      </ListPane>

      {!search.selected ? (
        <DetailPane empty emptyTitle="No loan product selected" emptyHint="Pick a product, or click + New to create one." />
      ) : selected.isLoading || !selected.data ? (
        <DetailPane><div className="text-foreground-muted">Loading…</div></DetailPane>
      ) : (
        <LoanProductDetail product={selected.data} />
      )}

      <CreateLoanProductDrawer open={createOpen} onClose={() => setCreateOpen(false)} onCreated={(id) => update({ selected: id })} />
    </TwoPane>
  )
}

function LoanProductDetail({ product }: { product: NonNullable<ReturnType<typeof useLoanProduct>['data']> }) {
  const [editOpen, setEditOpen] = useState(false)
  const cycleUnit = product.repayment_cycle === 'monthly' ? 'months' : product.repayment_cycle === 'biweekly' ? 'biweeks' : 'weeks'

  return (
    <DetailPane
      header={
        <div className="flex items-center justify-between gap-4">
          <div className="min-w-0">
            <h1 className="text-[17px] font-semibold text-foreground">{product.name}</h1>
            {product.description && <p className="text-[12.5px] text-foreground-muted mt-0.5">{product.description}</p>}
          </div>
          <div className="flex items-center gap-2 shrink-0">
            <StatusBadge tone={product.is_active ? 'success' : 'neutral'}>
              {product.is_active ? 'Active' : 'Inactive'}
            </StatusBadge>
            <button type="button" onClick={() => setEditOpen(true)} className="h-8 px-2.5 rounded-lg border border-border bg-surface text-[12px] font-medium text-foreground-secondary hover:bg-surface-hover transition inline-flex items-center gap-1">
              <Pencil size={12} /> Edit
            </button>
          </div>
        </div>
      }
    >
      <DetailSection title="Pricing">
        <FieldGrid columns={2}>
          <Field label="Min amount" value={formatUGX(product.min_amount)} />
          <Field label="Max amount" value={formatUGX(product.max_amount)} />
          <Field label="Interest rate" value={`${product.interest_rate}% / yr`} />
          <Field label="Interest method" value={product.interest_method.replace('_', ' ')} />
        </FieldGrid>
      </DetailSection>

      <DetailSection title="Schedule">
        <FieldGrid columns={2}>
          <Field label="Repayment cycle" value={product.repayment_cycle} />
          <Field label="Duration range" value={`${product.min_duration} – ${product.max_duration} ${cycleUnit}`} />
          <Field label="Grace period" value={product.grace_period_days ? `${product.grace_period_days} days` : null} />
          <Field label="Requires collateral" value={product.requires_collateral ? 'Yes' : 'No'} />
        </FieldGrid>
      </DetailSection>

      <EditLoanProductDrawer open={editOpen} onClose={() => setEditOpen(false)} product={product} />
    </DetailPane>
  )
}

const productSchema = z.object({
  name: z.string().min(1, 'Required'),
  description: z.string().optional(),
  min_amount: z.number().positive(),
  max_amount: z.number().positive(),
  min_duration: z.coerce.number().int().positive(),
  max_duration: z.coerce.number().int().positive(),
  interest_method: z.enum(['flat', 'reducing_balance']),
  interest_rate: z.coerce.number().min(0),
  repayment_cycle: z.enum(['weekly', 'biweekly', 'monthly']),
  requires_collateral: z.coerce.boolean().optional(),
  grace_period_days: z.coerce.number().int().min(0).optional(),
}).refine((d) => d.max_amount >= d.min_amount, { message: 'Max must be ≥ min', path: ['max_amount'] })
  .refine((d) => d.max_duration >= d.min_duration, { message: 'Max must be ≥ min', path: ['max_duration'] })

type ProductForm = z.infer<typeof productSchema>

function CreateLoanProductDrawer({ open, onClose, onCreated }: { open: boolean; onClose: () => void; onCreated: (id: number) => void }) {
  const create = useCreateLoanProduct()
  const { register, handleSubmit, formState: { errors }, reset, control } = useForm<ProductForm>({
    resolver: zodResolver(productSchema),
    defaultValues: {
      name: '', description: '',
      min_amount: undefined as any, max_amount: undefined as any,
      min_duration: 1, max_duration: 12,
      interest_method: 'flat', interest_rate: undefined as any,
      repayment_cycle: 'monthly', requires_collateral: false, grace_period_days: 0,
    },
  })

  const onSubmit = handleSubmit(async (data) => {
    try {
      const created = await create.mutateAsync(data as any)
      toast.success('Loan product created')
      reset()
      onClose()
      onCreated(created.id)
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  })

  return (
    <Drawer open={open} onOpenChange={(o) => { if (!o) reset(); onClose() }} title="New Loan Product" width={520}>
      <form onSubmit={onSubmit} className="space-y-5">
        <TextField label="Name" required placeholder="e.g. Boda 12-month" error={errors.name?.message} {...register('name')} />
        <FormGrid columns={2}>
          <CurrencyField control={control} name="min_amount" label="Min amount" required />
          <CurrencyField control={control} name="max_amount" label="Max amount" required />
        </FormGrid>
        <FormGrid columns={2}>
          <SelectField label="Repayment cycle" required options={[
            { value: 'weekly', label: 'Weekly' },
            { value: 'biweekly', label: 'Biweekly' },
            { value: 'monthly', label: 'Monthly' },
          ]} {...register('repayment_cycle')} />
          <SelectField label="Interest method" required options={[
            { value: 'flat', label: 'Flat' },
            { value: 'reducing_balance', label: 'Reducing balance' },
          ]} {...register('interest_method')} />
        </FormGrid>
        <FormGrid columns={2}>
          <TextField label="Min duration (cycles)" required type="number" error={errors.min_duration?.message} {...register('min_duration')} />
          <TextField label="Max duration (cycles)" required type="number" error={errors.max_duration?.message} {...register('max_duration')} />
        </FormGrid>
        <TextField label="Annual interest rate (%)" required type="number" step="0.1" placeholder="24" error={errors.interest_rate?.message} {...register('interest_rate')} />
        <FormActions onCancel={onClose} submitLabel="Create Loan Product" isPending={create.isPending} />
      </form>
    </Drawer>
  )
}

function EditLoanProductDrawer({ open, onClose, product }: { open: boolean; onClose: () => void; product: NonNullable<ReturnType<typeof useLoanProduct>['data']> }) {
  const update = useUpdateLoanProduct(product.id)
  const { register, handleSubmit, formState: { errors }, control } = useForm<ProductForm>({
    resolver: zodResolver(productSchema),
    defaultValues: {
      name: product.name, description: product.description || '',
      min_amount: product.min_amount, max_amount: product.max_amount,
      min_duration: product.min_duration, max_duration: product.max_duration,
      interest_method: product.interest_method, interest_rate: product.interest_rate,
      repayment_cycle: product.repayment_cycle, requires_collateral: product.requires_collateral,
      grace_period_days: product.grace_period_days,
    },
  })

  const onSubmit = handleSubmit(async (data) => {
    try {
      await update.mutateAsync({ ...data, is_active: product.is_active } as any)
      toast.success('Updated')
      onClose()
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  })

  return (
    <Drawer open={open} onOpenChange={(o) => !o && onClose()} title="Edit Loan Product" width={560}>
      <form onSubmit={onSubmit} className="space-y-6">
        <FormSection title="Identity">
          <TextField label="Name" required error={errors.name?.message} {...register('name')} />
          <TextAreaField label="Description" {...register('description')} />
        </FormSection>
        <FormSection title="Pricing">
          <FormGrid columns={2}>
            <CurrencyField control={control} name="min_amount" label="Min amount" required />
            <CurrencyField control={control} name="max_amount" label="Max amount" required />
          </FormGrid>
          <TextField label="Annual interest rate (%)" required type="number" step="0.1" {...register('interest_rate')} />
        </FormSection>
        <FormSection title="Schedule">
          <FormGrid columns={2}>
            <SelectField label="Repayment cycle" required options={[
              { value: 'weekly', label: 'Weekly' },
              { value: 'biweekly', label: 'Biweekly' },
              { value: 'monthly', label: 'Monthly' },
            ]} {...register('repayment_cycle')} />
            <SelectField label="Interest method" required options={[
              { value: 'flat', label: 'Flat' },
              { value: 'reducing_balance', label: 'Reducing balance' },
            ]} {...register('interest_method')} />
          </FormGrid>
          <FormGrid columns={2}>
            <TextField label="Min duration" required type="number" {...register('min_duration')} />
            <TextField label="Max duration" required type="number" {...register('max_duration')} />
          </FormGrid>
          <TextField label="Grace period (days)" type="number" {...register('grace_period_days')} />
        </FormSection>
        <FormActions onCancel={onClose} submitLabel="Save changes" isPending={update.isPending} />
      </form>
    </Drawer>
  )
}
