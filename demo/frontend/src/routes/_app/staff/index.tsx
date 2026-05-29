import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Trash2, Mail, Check, Eye, EyeOff } from 'lucide-react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'

import api from '@/lib/api'
import { useBranches } from '@/hooks/useBusiness'
import {
  TwoPane, ListPane, ListRow, DetailPane, DetailSection, Field, FieldGrid,
} from '@/components/two-pane'
import { Drawer } from '@/components/drawer'
import { TextField, SelectField, FormActions } from '@/components/form'
import { StatusBadge } from '@/components/status-badge'
import { ExportButton } from '@/components/export-button'
import { exportToExcel } from '@/lib/export'

type SSearch = { selected?: number; q?: string }

export const Route = createFileRoute('/_app/staff/')({
  component: StaffPage,
  validateSearch: (s: Record<string, unknown>): SSearch => ({
    selected: s.selected ? Number(s.selected) : undefined,
    q: typeof s.q === 'string' ? s.q : undefined,
  }),
})

const roleTone = (role: string) => {
  if (role === 'admin') return 'accent' as const
  if (role === 'manager') return 'info' as const
  if (role === 'loan_officer') return 'info' as const
  if (role === 'accountant') return 'warning' as const
  if (role === 'cashier') return 'success' as const
  if (role === 'stock_clerk') return 'warning' as const
  return 'neutral' as const
}

function StaffPage() {
  const search = Route.useSearch()
  const navigate = useNavigate()
  const qc = useQueryClient()
  const [createOpen, setCreateOpen] = useState(false)

  const staff = useQuery({
    queryKey: ['staff'],
    queryFn: async () => { const r = await api.get('/staff'); return r.data.data as any[] },
  })
  const invitations = useQuery({
    queryKey: ['invitations'],
    queryFn: async () => { const r = await api.get('/invitations'); return r.data.data as any[] },
  })

  const update = (patch: Partial<SSearch>) =>
    navigate({ to: '/staff', search: { ...search, ...patch }, replace: true })

  const updateRole = useMutation({
    mutationFn: async ({ userId, role, branchId, workspaceAccess }: { userId: number; role: string; branchId?: number | null; workspaceAccess?: string }) => {
      await api.put(`/staff/${userId}/role`, { role, branch_id: branchId ?? null, workspace_access: workspaceAccess })
    },
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['staff'] }); toast.success('Updated') },
    onError: (err: any) => toast.error(err.response?.data?.error || 'Failed'),
  })

  const removeStaff = useMutation({
    mutationFn: async (userId: number) => { await api.delete(`/staff/${userId}`) },
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['staff'] }); toast.success('Staff removed') },
    onError: (err: any) => toast.error(err.response?.data?.error || 'Failed'),
  })

  const filtered = (staff.data || []).filter((s) =>
    !search.q ? true : s.name?.toLowerCase().includes(search.q.toLowerCase()) || s.email?.toLowerCase().includes(search.q.toLowerCase()),
  )
  const selected = filtered.find((s: any) => s.user_id === search.selected) || (staff.data || []).find((s: any) => s.user_id === search.selected)

  return (
    <TwoPane>
      <ListPane
        title="Staff"
        count={staff.data?.length}
        search={search.q || ''}
        onSearch={(q) => update({ q: q || undefined })}
        searchPlaceholder="Name or email..."
        newLabel="Add User"
        onNew={() => setCreateOpen(true)}
        refreshKey={['staff']}
        extraAction={
          <ExportButton
            onClick={() => exportToExcel(
              filtered.map((s: any) => ({
                Name: s.name,
                Email: s.email,
                Role: s.role,
                Branch: s.branch_name || 'All Branches',
              })),
              'staff',
              'Staff',
            )}
            disabled={filtered.length === 0}
          />
        }
        footer={
          invitations.data && invitations.data.length > 0 ? (
            <span><Mail className="inline h-3 w-3 mr-1" /> {invitations.data.length} pending invitation{invitations.data.length === 1 ? '' : 's'}</span>
          ) : null
        }
      >
        {staff.isLoading && <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">Loading…</div>}
        {!staff.isLoading && filtered.length === 0 && (
          <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">No team members yet.</div>
        )}
        {filtered.map((s: any) => (
          <ListRow
            key={s.user_id}
            selected={search.selected === s.user_id}
            onClick={() => update({ selected: s.user_id })}
            icon={<span className="text-[12px] font-semibold text-foreground-secondary">{(s.name || '').charAt(0).toUpperCase()}</span>}
            title={s.name}
            subtitle={s.email}
            rightTop={s.branch_name || 'All Branches'}
            rightBottom={<StatusBadge tone={roleTone(s.role)} size="sm">{s.role.replace('_', ' ')}</StatusBadge>}
          />
        ))}
      </ListPane>

      {!selected ? (
        <DetailPane empty emptyTitle="No staff selected" emptyHint="Pick a team member to view their profile, or click + Add User to create one." />
      ) : (
        <StaffDetail
          s={selected}
          onChange={(patch) => updateRole.mutate({
            userId: selected.user_id,
            role: patch.role ?? selected.role,
            branchId: 'branchId' in patch ? patch.branchId : selected.branch_id,
            workspaceAccess: patch.workspaceAccess ?? selected.workspace_access ?? 'both',
          })}
          onRemove={() => {
            if (confirm(`Remove ${selected.name}?`)) {
              removeStaff.mutate(selected.user_id)
              update({ selected: undefined })
            }
          }}
        />
      )}

      <CreateUserDrawer open={createOpen} onClose={() => setCreateOpen(false)} />
    </TwoPane>
  )
}

type StaffPatch = { role?: string; branchId?: number | null; workspaceAccess?: string }

function StaffDetail({ s, onChange, onRemove }: { s: any; onChange: (patch: StaffPatch) => void; onRemove: () => void }) {
  const { data: branches } = useBranches()
  const access = s.workspace_access || 'both'
  return (
    <DetailPane
      header={
        <div className="flex items-center justify-between gap-4">
          <div className="flex items-center gap-3 min-w-0">
            <div className="h-12 w-12 rounded-full bg-accent-light flex items-center justify-center text-accent font-semibold text-[16px] shrink-0">
              {(s.name || '').charAt(0).toUpperCase()}
            </div>
            <div className="min-w-0">
              <h1 className="text-[17px] font-semibold text-foreground truncate">{s.name}</h1>
              <p className="text-[12.5px] text-foreground-muted truncate">{s.email}</p>
            </div>
          </div>
          <button onClick={onRemove} className="h-8 px-2.5 rounded-lg bg-danger-light text-danger-dark text-[12px] font-medium hover:bg-danger/20 transition inline-flex items-center gap-1">
            <Trash2 size={12} /> Remove
          </button>
        </div>
      }
    >
      <DetailSection title="Access">
        <FieldGrid columns={2}>
          <div>
            <div className="text-[11px] font-medium uppercase tracking-wider text-foreground-muted mb-1">Role</div>
            <select
              value={s.role}
              onChange={(e) => onChange({ role: e.target.value })}
              className="w-full h-10 px-3 rounded-lg border border-border bg-surface text-[13.5px] text-foreground focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/15"
            >
              <option value="admin">Admin</option>
              <option value="manager">Manager</option>
              <option value="loan_officer">Loan Officer</option>
              <option value="accountant">Accountant</option>
              <option value="cashier">Cashier</option>
              <option value="stock_clerk">Stock Clerk</option>
            </select>
          </div>
          <div>
            <div className="text-[11px] font-medium uppercase tracking-wider text-foreground-muted mb-1">Branch scope</div>
            <select
              value={s.branch_id ?? ''}
              onChange={(e) => onChange({ branchId: e.target.value ? Number(e.target.value) : null })}
              className="w-full h-10 px-3 rounded-lg border border-border bg-surface text-[13.5px] text-foreground focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/15"
            >
              <option value="">All Branches</option>
              {branches?.map((b) => <option key={b.id} value={b.id}>{b.name}</option>)}
            </select>
          </div>
          <div className="col-span-2">
            <div className="text-[11px] font-medium uppercase tracking-wider text-foreground-muted mb-1">Workspace access</div>
            <div className="flex gap-2">
              {(['both', 'loans', 'spares'] as const).map((opt) => (
                <button
                  key={opt}
                  type="button"
                  onClick={() => onChange({ workspaceAccess: opt })}
                  className={`flex-1 h-10 px-3 rounded-lg border text-[13px] font-medium transition ${
                    access === opt
                      ? 'border-accent bg-accent-light text-accent-hover'
                      : 'border-border bg-surface text-foreground-secondary hover:bg-surface-hover'
                  }`}
                >
                  {opt === 'both' ? 'Both workspaces' : opt === 'loans' ? 'Loans only' : 'Spares only'}
                </button>
              ))}
            </div>
          </div>
        </FieldGrid>
      </DetailSection>

      <DetailSection title="What this role can do">
        <RolePermissions role={s.role} />
      </DetailSection>
    </DetailPane>
  )
}

const ROLE_PERMISSIONS: Record<string, string[]> = {
  admin: ['Full access to everything', 'Manage staff & branches', 'View all reports', 'Approve loans & repayments'],
  manager: ['POS, products, stock-in', 'Branch reports', 'Approve loans & repayments', 'Cannot manage staff'],
  loan_officer: ['Manage borrowers, loans, repayments', 'Read-only on motorcycles', 'No spares POS access'],
  accountant: ['Read-only finance reports across segments', 'Verify daily-boda payments', 'No write access'],
  cashier: ['POS sales', 'Cash sales of motorcycles', 'Record loan repayments (pending verification)'],
  stock_clerk: ['Receive stock only', 'No POS, no reports'],
}

function RolePermissions({ role }: { role: string }) {
  const perms = ROLE_PERMISSIONS[role] || []
  return (
    <ul className="space-y-1.5">
      {perms.map((p, i) => (
        <li key={i} className="flex items-center gap-2 text-[13px] text-foreground">
          <Check size={14} className="text-success-dark shrink-0" />
          {p}
        </li>
      ))}
    </ul>
  )
}

const createUserSchema = z.object({
  name: z.string().min(2, 'Name is required'),
  email: z.string().email('Valid email required'),
  password: z.string().min(6, 'At least 6 characters'),
  role: z.enum(['admin', 'manager', 'loan_officer', 'accountant', 'cashier', 'stock_clerk']),
  branch_id: z.string().optional(),
  workspace_access: z.enum(['both', 'loans', 'spares']).default('both'),
})
type CreateUserForm = z.infer<typeof createUserSchema>

function CreateUserDrawer({ open, onClose }: { open: boolean; onClose: () => void }) {
  const qc = useQueryClient()
  const { data: branches } = useBranches()
  const [showPwd, setShowPwd] = useState(false)
  const { register, handleSubmit, formState: { errors }, reset, watch, setValue } = useForm<CreateUserForm>({
    resolver: zodResolver(createUserSchema),
    defaultValues: { name: '', email: '', password: '', role: 'cashier', branch_id: '', workspace_access: 'both' },
  })

  const role = watch('role')
  const workspaceAccess = watch('workspace_access')

  const onSubmit = handleSubmit(async (data) => {
    try {
      await api.post('/staff', {
        name: data.name.trim(),
        email: data.email.trim().toLowerCase(),
        password: data.password,
        role: data.role,
        branch_id: data.branch_id ? Number(data.branch_id) : null,
        workspace_access: data.workspace_access,
      })
      qc.invalidateQueries({ queryKey: ['staff'] })
      toast.success(`${data.name} added`)
      reset()
      onClose()
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  })

  return (
    <Drawer open={open} onOpenChange={(o) => { if (!o) reset(); onClose() }} title="Add Staff User" width={520}>
      <form onSubmit={onSubmit} className="space-y-5">
        <TextField label="Full name" required placeholder="e.g. Jane Doe" error={errors.name?.message} {...register('name')} />
        <TextField label="Email" type="email" required placeholder="staff@example.com" error={errors.email?.message} {...register('email')} />

        <div>
          <label className="block text-[11.5px] font-medium uppercase tracking-wider text-foreground-muted mb-1">Password</label>
          <div className="relative">
            <input
              type={showPwd ? 'text' : 'password'}
              placeholder="At least 6 characters"
              {...register('password')}
              className="w-full h-10 pl-3 pr-10 rounded-lg border border-border bg-surface text-[13.5px] text-foreground focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/15"
            />
            <button
              type="button"
              onClick={() => setShowPwd((v) => !v)}
              className="absolute right-1 top-1/2 -translate-y-1/2 h-8 w-8 inline-flex items-center justify-center rounded-md text-foreground-muted hover:bg-surface-hover hover:text-foreground"
              aria-label={showPwd ? 'Hide password' : 'Show password'}
            >
              {showPwd ? <EyeOff size={14} /> : <Eye size={14} />}
            </button>
          </div>
          {errors.password && <p className="text-[11.5px] text-danger-dark mt-1">{errors.password.message}</p>}
        </div>

        <SelectField
          label="Role"
          required
          options={[
            { value: 'admin', label: 'Admin' },
            { value: 'manager', label: 'Manager' },
            { value: 'loan_officer', label: 'Loan Officer' },
            { value: 'accountant', label: 'Accountant' },
            { value: 'cashier', label: 'Cashier' },
            { value: 'stock_clerk', label: 'Stock Clerk' },
          ]}
          {...register('role')}
        />

        <SelectField
          label="Branch scope"
          placeholder="All branches"
          options={branches?.map((b) => ({ value: String(b.id), label: b.name })) || []}
          {...register('branch_id')}
        />

        <div>
          <label className="block text-[11.5px] font-medium uppercase tracking-wider text-foreground-muted mb-1.5">Workspace access</label>
          <div className="flex gap-2">
            {(['both', 'loans', 'spares'] as const).map((opt) => (
              <button
                key={opt}
                type="button"
                onClick={() => setValue('workspace_access', opt, { shouldDirty: true })}
                className={`flex-1 h-10 px-3 rounded-lg border text-[12.5px] font-medium transition ${
                  workspaceAccess === opt
                    ? 'border-accent bg-accent-light text-accent-hover'
                    : 'border-border bg-surface text-foreground-secondary hover:bg-surface-hover'
                }`}
              >
                {opt === 'both' ? 'Both' : opt === 'loans' ? 'Loans only' : 'Spares only'}
              </button>
            ))}
          </div>
          <p className="text-[11px] text-foreground-muted mt-1.5">
            Both = full sidebar. Loans only = motorcycle, loans, daily-boda. Spares only = POS + inventory.
          </p>
        </div>

        <div className="rounded-lg border border-border-subtle bg-surface-2 p-3">
          <p className="text-[10.5px] font-semibold uppercase tracking-wider text-foreground-muted mb-2">{role.replace('_', ' ')} can:</p>
          <RolePermissions role={role} />
        </div>
        <FormActions onCancel={onClose} submitLabel="Create User" />
      </form>
    </Drawer>
  )
}
