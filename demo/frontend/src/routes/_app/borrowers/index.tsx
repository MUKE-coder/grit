import { createFileRoute, useNavigate, Link } from '@tanstack/react-router'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Phone, Mail, IdCard, Briefcase, Pencil, Check, X, ArrowRight } from 'lucide-react'
import toast from 'react-hot-toast'

import { useBorrowers, useBorrower, useCreateBorrower, useUpdateBorrower } from '@/hooks/useBorrowers'
import { useBranches } from '@/hooks/useBusiness'
import { useLoans } from '@/hooks/useLoans'
import {
  TwoPane, ListPane, ListRow, DetailPane, DetailSection, Field, FieldGrid,
} from '@/components/two-pane'
import { Drawer } from '@/components/drawer'
import {
  TextField, SelectField, FormGrid, FormActions, FormSection, TextAreaField,
} from '@/components/form'
import { CurrencyField } from '@/components/currency-field'
import { StatusPill } from '@/components/status-badge'
import { ExportButton } from '@/components/export-button'
import { exportToExcel } from '@/lib/export'
import { formatUGX } from '@/lib/utils'

/** URL search shape: `?selected=123` keeps the active row in the URL so it
 *  survives refresh + browser back/forward. */
type BorrowerSearch = { selected?: number; q?: string; risk?: string }

export const Route = createFileRoute('/_app/borrowers/')({
  component: BorrowersPage,
  validateSearch: (s: Record<string, unknown>): BorrowerSearch => ({
    selected: typeof s.selected === 'number' ? s.selected : s.selected ? Number(s.selected) : undefined,
    q: typeof s.q === 'string' ? s.q : undefined,
    risk: typeof s.risk === 'string' ? s.risk : undefined,
  }),
})

function BorrowersPage() {
  const search = Route.useSearch()
  const navigate = useNavigate()
  const [createOpen, setCreateOpen] = useState(false)

  const list = useBorrowers({ search: search.q, risk_level: search.risk })
  const selected = useBorrower(search.selected)

  const setSelected = (id: number | undefined) =>
    navigate({ to: '/borrowers', search: { ...search, selected: id }, replace: true })

  const setQ = (q: string) =>
    navigate({ to: '/borrowers', search: { ...search, q: q || undefined }, replace: true })

  const setRisk = (risk: string) =>
    navigate({ to: '/borrowers', search: { ...search, risk: risk || undefined }, replace: true })

  return (
    <TwoPane>
      <ListPane
        title="Borrowers"
        count={list.data?.total}
        search={search.q || ''}
        onSearch={setQ}
        searchPlaceholder="Name, phone, NIN..."
        onNew={() => setCreateOpen(true)}
        refreshKey={['borrowers']}
        extraAction={
          <ExportButton
            onClick={() => exportToExcel(
              (list.data?.data || []).map((b: any) => ({
                'First Name': b.first_name,
                'Last Name': b.last_name,
                Phone: b.phone,
                'Alt Phone': b.alt_phone,
                Email: b.email,
                NIN: b.national_id,
                Branch: b.branch?.name,
                Gender: b.gender,
                'Date of Birth': b.date_of_birth ? new Date(b.date_of_birth).toLocaleDateString() : '',
                Address: b.address,
                'Employment Status': b.employment_status,
                Occupation: b.occupation,
                Employer: b.employer,
                'Monthly Income': b.monthly_income,
                'Next of Kin': b.next_of_kin_name,
                'Next of Kin Phone': b.next_of_kin_phone,
                'Next of Kin Relation': b.next_of_kin_relation,
                'Risk Level': b.risk_level,
                'Credit Score': b.credit_score,
                'Created At': b.created_at ? new Date(b.created_at).toLocaleDateString() : '',
              })),
              'borrowers',
              'Borrowers',
            )}
            disabled={!list.data?.data.length}
          />
        }
        filters={
          <div className="flex gap-1 flex-wrap">
            {[
              { val: '', label: 'All' },
              { val: 'low', label: 'Low' },
              { val: 'medium', label: 'Medium' },
              { val: 'high', label: 'High' },
              { val: 'critical', label: 'Critical' },
            ].map((r) => (
              <button
                key={r.val}
                type="button"
                onClick={() => setRisk(r.val)}
                className={`px-2.5 py-1 rounded-md text-[11.5px] font-medium transition ${
                  (search.risk || '') === r.val
                    ? 'bg-foreground text-white'
                    : 'bg-surface-2 text-foreground-secondary hover:bg-surface-hover'
                }`}
              >
                {r.label}
              </button>
            ))}
          </div>
        }
      >
        {list.isLoading && (
          <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">Loading…</div>
        )}
        {!list.isLoading && list.data?.data.length === 0 && (
          <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">
            No borrowers match.
          </div>
        )}
        {list.data?.data.map((b) => (
          <ListRow
            key={b.id}
            selected={search.selected === b.id}
            onClick={() => setSelected(b.id)}
            icon={<Initials first={b.first_name} last={b.last_name} />}
            title={b.full_name}
            subtitle={b.phone}
            rightTop={b.branch?.name}
            rightBottom={<StatusPill status={b.risk_level} size="sm" />}
          />
        ))}
      </ListPane>

      {!search.selected ? (
        <DetailPane
          empty
          emptyTitle="No borrower selected"
          emptyHint="Pick a borrower from the list, or click + New to register one."
        />
      ) : selected.isLoading || !selected.data ? (
        <DetailPane><div className="text-foreground-muted">Loading…</div></DetailPane>
      ) : (
        <BorrowerDetail b={selected.data} />
      )}

      <CreateBorrowerDrawer
        open={createOpen}
        onClose={() => setCreateOpen(false)}
        onCreated={(id) => setSelected(id)}
      />
    </TwoPane>
  )
}

function Initials({ first, last }: { first: string; last: string }) {
  return <span>{(first[0] || '') + (last[0] || '')}</span>
}

/* ─────────────── Detail pane ─────────────── */

function BorrowerDetail({ b }: { b: ReturnType<typeof useBorrower>['data'] & {} }) {
  const [editOpen, setEditOpen] = useState(false)
  const loans = useLoans({ borrower_id: String(b.id) })

  return (
    <DetailPane
      header={
        <div className="flex items-center justify-between gap-4">
          <div className="min-w-0 flex items-center gap-3">
            <div className="h-11 w-11 rounded-full bg-accent-light flex items-center justify-center text-accent font-semibold text-[14px] shrink-0">
              {(b.first_name[0] || '') + (b.last_name[0] || '')}
            </div>
            <div className="min-w-0">
              <h1 className="text-[17px] font-semibold text-foreground truncate">{b.full_name}</h1>
              <p className="text-[12.5px] text-foreground-muted">{b.phone} · {b.branch?.name || '—'}</p>
            </div>
          </div>
          <div className="flex items-center gap-2 shrink-0">
            <StatusPill status={b.risk_level} />
            <button
              type="button"
              onClick={() => setEditOpen(true)}
              className="h-8 px-2.5 rounded-lg border border-border bg-surface text-[12px] font-medium text-foreground-secondary hover:bg-surface-hover transition inline-flex items-center gap-1"
            >
              <Pencil size={12} /> Edit
            </button>
            <Link
              to="/loans/new"
              search={{ borrower_id: String(b.id) }}
              className="h-8 px-3 rounded-lg bg-accent text-white text-[12px] font-medium hover:bg-accent-hover transition inline-flex items-center gap-1"
            >
              New Loan <ArrowRight size={12} />
            </Link>
          </div>
        </div>
      }
    >
      <DetailSection title="Contact">
        <FieldGrid columns={2}>
          <Field label="Phone" value={<span className="font-mono">{b.phone}</span>} />
          {b.alt_phone && <Field label="Alt phone" value={<span className="font-mono">{b.alt_phone}</span>} />}
          {b.email && <Field label="Email" value={b.email} />}
          {b.national_id && <Field label="NIN" value={<span className="font-mono">{b.national_id}</span>} />}
          {b.address && <Field label="Address" value={b.address} />}
        </FieldGrid>
      </DetailSection>

      <DetailSection title="Employment & Income">
        <FieldGrid columns={2}>
          <Field label="Status" value={b.employment_status} />
          <Field label="Occupation" value={b.occupation} />
          <Field label="Employer" value={b.employer} />
          <Field
            label="Monthly income"
            value={b.monthly_income > 0 ? formatUGX(b.monthly_income) : null}
          />
          <Field label="Credit score" value={String(b.credit_score)} />
          <Field label="Risk level" value={<StatusPill status={b.risk_level} size="sm" />} />
        </FieldGrid>
      </DetailSection>

      {(b.next_of_kin_name || b.next_of_kin_phone) && (
        <DetailSection title="Next of Kin">
          <FieldGrid columns={3}>
            <Field label="Name" value={b.next_of_kin_name} />
            <Field label="Phone" value={<span className="font-mono">{b.next_of_kin_phone}</span>} />
            <Field label="Relation" value={b.next_of_kin_relation} />
          </FieldGrid>
        </DetailSection>
      )}

      <DetailSection
        title={`Loans (${loans.data?.total ?? 0})`}
        action={
          <Link
            to="/loans"
            search={{ borrower_id: String(b.id) }}
            className="text-[11.5px] text-accent font-medium hover:underline inline-flex items-center gap-0.5"
          >
            View all <ArrowRight size={11} />
          </Link>
        }
      >
        {!loans.data?.data.length ? (
          <p className="text-[13px] text-foreground-muted">No loans yet.</p>
        ) : (
          <div className="border border-border rounded-lg overflow-hidden">
            <table className="w-full text-[13px]">
              <thead className="bg-surface-2 text-[11px] uppercase tracking-wider text-foreground-muted">
                <tr>
                  <th className="text-left py-2 px-3 font-semibold">Loan #</th>
                  <th className="text-right py-2 px-3 font-semibold">Principal</th>
                  <th className="text-right py-2 px-3 font-semibold">Balance</th>
                  <th className="text-left py-2 px-3 font-semibold">Status</th>
                </tr>
              </thead>
              <tbody>
                {loans.data.data.slice(0, 5).map((l) => (
                  <tr key={l.id} className="border-t border-border-subtle hover:bg-surface-2">
                    <td className="py-2 px-3">
                      <Link to={`/loans/${l.id}`} className="text-accent font-mono hover:underline">
                        {l.loan_number}
                      </Link>
                    </td>
                    <td className="py-2 px-3 text-right tabular-nums">{formatUGX(l.principal_amount)}</td>
                    <td className="py-2 px-3 text-right tabular-nums font-medium">{formatUGX(l.balance_remaining)}</td>
                    <td className="py-2 px-3"><StatusPill status={l.status} size="sm" /></td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </DetailSection>

      <EditBorrowerDrawer open={editOpen} onClose={() => setEditOpen(false)} borrower={b} />
    </DetailPane>
  )
}

/* ─────────────── Create drawer (REQUIRED fields only) ─────────────── */

const createSchema = z.object({
  branch_id: z.coerce.number({ message: 'Branch is required' }).int().positive('Branch is required'),
  first_name: z.string().min(1, 'Required'),
  last_name: z.string().min(1, 'Required'),
  phone: z.string().min(7, 'Phone is required'),
})
type CreateForm = z.infer<typeof createSchema>

function CreateBorrowerDrawer({
  open,
  onClose,
  onCreated,
}: {
  open: boolean
  onClose: () => void
  onCreated: (id: number) => void
}) {
  const { data: branches } = useBranches()
  const create = useCreateBorrower()

  const { register, handleSubmit, formState: { errors }, reset } = useForm<CreateForm>({
    resolver: zodResolver(createSchema),
    defaultValues: { branch_id: undefined as any, first_name: '', last_name: '', phone: '' },
  })

  const onSubmit = handleSubmit(async (data) => {
    try {
      const created = await create.mutateAsync(data as any)
      toast.success('Borrower added')
      reset()
      onClose()
      onCreated(created.id)
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed to create borrower')
    }
  })

  return (
    <Drawer
      open={open}
      onOpenChange={(o) => {
        if (!o) reset()
        onClose()
      }}
      title="New Borrower"
      description="Just the basics — you can fill in the rest after creation."
    >
      <form onSubmit={onSubmit} className="space-y-5">
        <SelectField
          label="Branch"
          required
          placeholder="Select branch"
          error={errors.branch_id?.message}
          options={branches?.map((b) => ({ value: String(b.id), label: b.name })) || []}
          {...register('branch_id')}
        />

        <FormGrid columns={2}>
          <TextField label="First name" required error={errors.first_name?.message} {...register('first_name')} />
          <TextField label="Last name" required error={errors.last_name?.message} {...register('last_name')} />
        </FormGrid>

        <TextField
          label="Phone"
          required
          placeholder="0700000000"
          hint="Used to contact the borrower and to send mobile-money requests."
          error={errors.phone?.message}
          {...register('phone')}
        />

        <FormActions
          onCancel={onClose}
          submitLabel="Create Borrower"
          isPending={create.isPending}
        />
      </form>
    </Drawer>
  )
}

/* ─────────────── Edit drawer (FULL field set) ─────────────── */

const editSchema = z.object({
  branch_id: z.coerce.number().int().positive(),
  first_name: z.string().min(1),
  last_name: z.string().min(1),
  phone: z.string().min(7),
  alt_phone: z.string().optional(),
  email: z.string().email().or(z.literal('')).optional(),
  national_id: z.string().optional(),
  address: z.string().optional(),
  employment_status: z.string().optional(),
  occupation: z.string().optional(),
  employer: z.string().optional(),
  monthly_income: z.number().optional(),
  next_of_kin_name: z.string().optional(),
  next_of_kin_phone: z.string().optional(),
  next_of_kin_relation: z.string().optional(),
  risk_level: z.enum(['low', 'medium', 'high', 'critical']),
})
type EditForm = z.infer<typeof editSchema>

function EditBorrowerDrawer({
  open,
  onClose,
  borrower,
}: {
  open: boolean
  onClose: () => void
  borrower: NonNullable<ReturnType<typeof useBorrower>['data']>
}) {
  const { data: branches } = useBranches()
  const update = useUpdateBorrower(borrower.id)

  const { register, handleSubmit, formState: { errors }, control } = useForm<EditForm>({
    resolver: zodResolver(editSchema),
    defaultValues: {
      branch_id: borrower.branch_id,
      first_name: borrower.first_name,
      last_name: borrower.last_name,
      phone: borrower.phone,
      alt_phone: borrower.alt_phone || '',
      email: borrower.email || '',
      national_id: borrower.national_id || '',
      address: borrower.address || '',
      employment_status: borrower.employment_status || '',
      occupation: borrower.occupation || '',
      employer: borrower.employer || '',
      monthly_income: borrower.monthly_income || undefined,
      next_of_kin_name: borrower.next_of_kin_name || '',
      next_of_kin_phone: borrower.next_of_kin_phone || '',
      next_of_kin_relation: borrower.next_of_kin_relation || '',
      risk_level: (borrower.risk_level as any) || 'medium',
    },
  })

  const onSubmit = handleSubmit(async (data) => {
    try {
      await update.mutateAsync(data as any)
      toast.success('Borrower updated')
      onClose()
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed to update')
    }
  })

  return (
    <Drawer open={open} onOpenChange={(o) => !o && onClose()} title="Edit Borrower" width={560}>
      <form onSubmit={onSubmit} className="space-y-6">
        <FormSection title="Identity">
          <SelectField
            label="Branch"
            required
            options={branches?.map((b) => ({ value: String(b.id), label: b.name })) || []}
            error={errors.branch_id?.message}
            {...register('branch_id')}
          />
          <FormGrid columns={2}>
            <TextField label="First name" required error={errors.first_name?.message} {...register('first_name')} />
            <TextField label="Last name" required error={errors.last_name?.message} {...register('last_name')} />
          </FormGrid>
          <FormGrid columns={2}>
            <TextField label="Phone" required error={errors.phone?.message} {...register('phone')} />
            <TextField label="Alt phone" error={errors.alt_phone?.message} {...register('alt_phone')} />
          </FormGrid>
          <FormGrid columns={2}>
            <TextField label="Email" type="email" error={errors.email?.message} {...register('email')} />
            <TextField label="National ID (NIN)" error={errors.national_id?.message} {...register('national_id')} />
          </FormGrid>
          <TextAreaField label="Address" error={errors.address?.message} {...register('address')} />
        </FormSection>

        <FormSection title="Employment">
          <FormGrid columns={2}>
            <SelectField
              label="Status"
              placeholder="—"
              options={[
                { value: 'employed', label: 'Employed' },
                { value: 'self_employed', label: 'Self-employed' },
                { value: 'unemployed', label: 'Unemployed' },
                { value: 'student', label: 'Student' },
                { value: 'retired', label: 'Retired' },
              ]}
              {...register('employment_status')}
            />
            <TextField label="Occupation" {...register('occupation')} />
          </FormGrid>
          <FormGrid columns={2}>
            <TextField label="Employer" {...register('employer')} />
            <CurrencyField control={control} name="monthly_income" label="Monthly income" />
          </FormGrid>
        </FormSection>

        <FormSection title="Next of Kin">
          <FormGrid columns={2}>
            <TextField label="Name" {...register('next_of_kin_name')} />
            <TextField label="Phone" {...register('next_of_kin_phone')} />
          </FormGrid>
          <TextField label="Relation" placeholder="e.g. Spouse, Parent, Sibling" {...register('next_of_kin_relation')} />
        </FormSection>

        <FormSection title="Risk Assessment">
          <SelectField
            label="Risk level"
            required
            options={[
              { value: 'low', label: 'Low' },
              { value: 'medium', label: 'Medium' },
              { value: 'high', label: 'High' },
              { value: 'critical', label: 'Critical' },
            ]}
            {...register('risk_level')}
          />
        </FormSection>

        <FormActions onCancel={onClose} submitLabel="Save changes" isPending={update.isPending} />
      </form>
    </Drawer>
  )
}
