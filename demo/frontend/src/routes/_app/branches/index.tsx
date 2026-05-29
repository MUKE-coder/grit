import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Building2 } from 'lucide-react'
import toast from 'react-hot-toast'
import { useQueryClient, useMutation } from '@tanstack/react-query'

import { useBranches } from '@/hooks/useBusiness'
import {
  TwoPane, ListPane, ListRow, DetailPane, DetailSection, Field, FieldGrid,
} from '@/components/two-pane'
import { Drawer } from '@/components/drawer'
import { TextField, FormActions } from '@/components/form'
import { StatusBadge } from '@/components/status-badge'
import { ExportButton } from '@/components/export-button'
import { exportToExcel } from '@/lib/export'
import api from '@/lib/api'

type BSearch = { selected?: number; q?: string }

export const Route = createFileRoute('/_app/branches/')({
  component: BranchesPage,
  validateSearch: (s: Record<string, unknown>): BSearch => ({
    selected: s.selected ? Number(s.selected) : undefined,
    q: typeof s.q === 'string' ? s.q : undefined,
  }),
})

function BranchesPage() {
  const search = Route.useSearch()
  const navigate = useNavigate()
  const [createOpen, setCreateOpen] = useState(false)
  const list = useBranches()

  const update = (patch: Partial<BSearch>) =>
    navigate({ to: '/branches', search: { ...search, ...patch }, replace: true })

  const branches = (list.data || []).filter((b) =>
    !search.q ? true : b.name.toLowerCase().includes(search.q.toLowerCase())
  )
  const selected = branches.find((b) => b.id === search.selected) || (list.data || []).find((b) => b.id === search.selected)

  return (
    <TwoPane>
      <ListPane
        title="Branches"
        count={branches.length}
        search={search.q || ''}
        onSearch={(q) => update({ q: q || undefined })}
        searchPlaceholder="Search branches..."
        onNew={() => setCreateOpen(true)}
        refreshKey={['branches']}
        extraAction={
          <ExportButton
            onClick={() => exportToExcel(
              branches.map((b: any) => ({
                Name: b.name,
                Address: b.address,
                Default: b.is_default ? 'Yes' : 'No',
                'Created At': b.created_at ? new Date(b.created_at).toLocaleDateString() : '',
              })),
              'branches',
              'Branches',
            )}
            disabled={branches.length === 0}
          />
        }
      >
        {list.isLoading && <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">Loading…</div>}
        {!list.isLoading && branches.length === 0 && (
          <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">No branches.</div>
        )}
        {branches.map((b) => (
          <ListRow
            key={b.id}
            selected={search.selected === b.id}
            onClick={() => update({ selected: b.id })}
            icon={<Building2 size={14} className="text-foreground-secondary" />}
            title={b.name}
            subtitle={b.address || '—'}
            rightBottom={b.is_default ? <StatusBadge tone="success" size="sm">Default</StatusBadge> : null}
          />
        ))}
      </ListPane>

      {!selected ? (
        <DetailPane empty emptyTitle="No branch selected" />
      ) : (
        <BranchDetail branch={selected} />
      )}

      <CreateBranchDrawer open={createOpen} onClose={() => setCreateOpen(false)} onCreated={(id) => update({ selected: id })} />
    </TwoPane>
  )
}

function BranchDetail({ branch }: { branch: any }) {
  const [editOpen, setEditOpen] = useState(false)
  const qc = useQueryClient()
  const navigate = useNavigate()
  const remove = useMutation({
    mutationFn: async () => api.delete(`/branches/${branch.id}`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['branches'] })
      toast.success('Branch deleted')
      navigate({ to: '/branches', search: { selected: undefined } })
    },
    onError: (e: any) => toast.error(e.response?.data?.error || 'Failed to delete'),
  })

  return (
    <DetailPane
      header={
        <div className="flex items-center justify-between gap-4">
          <div>
            <h1 className="text-[17px] font-semibold text-foreground">{branch.name}</h1>
            <p className="text-[12.5px] text-foreground-muted mt-0.5">{branch.address || 'No address set'}</p>
          </div>
          <div className="flex items-center gap-2">
            {branch.is_default && <StatusBadge tone="success">Default — Permanent</StatusBadge>}
            {/* Edit + Delete only on non-default branches. The default Main Branch
                is the system's permanent fallback and can't be removed. */}
            {!branch.is_default && (
              <>
                <button onClick={() => setEditOpen(true)} className="h-8 px-2.5 rounded-lg border border-border bg-surface text-[12px] font-medium text-foreground-secondary hover:bg-surface-hover transition">
                  Edit
                </button>
                <button
                  onClick={() => {
                    if (window.confirm(`Delete branch "${branch.name}"? Any inventory must be transferred first.`)) {
                      remove.mutate()
                    }
                  }}
                  disabled={remove.isPending}
                  className="h-8 px-2.5 rounded-lg bg-danger-light text-danger-dark text-[12px] font-medium hover:bg-danger/20 transition disabled:opacity-50"
                >
                  Delete
                </button>
              </>
            )}
          </div>
        </div>
      }
    >
      <DetailSection title="Branch">
        <FieldGrid columns={2}>
          <Field label="Name" value={branch.name} />
          <Field label="Address" value={branch.address} />
          <Field label="Default" value={branch.is_default ? 'Yes' : 'No'} />
          <Field label="Created" value={branch.created_at ? new Date(branch.created_at).toLocaleDateString() : null} />
        </FieldGrid>
      </DetailSection>

      <EditBranchDrawer open={editOpen} onClose={() => setEditOpen(false)} branch={branch} />
    </DetailPane>
  )
}

const branchSchema = z.object({
  name: z.string().min(1, 'Required'),
  address: z.string().optional(),
})
type BranchForm = z.infer<typeof branchSchema>

function CreateBranchDrawer({ open, onClose, onCreated }: { open: boolean; onClose: () => void; onCreated: (id: number) => void }) {
  const qc = useQueryClient()
  const { register, handleSubmit, formState: { errors }, reset } = useForm<BranchForm>({
    resolver: zodResolver(branchSchema),
    defaultValues: { name: '', address: '' },
  })

  const onSubmit = handleSubmit(async (data) => {
    try {
      const res = await api.post('/branches', data)
      qc.invalidateQueries({ queryKey: ['branches'] })
      toast.success('Branch created')
      reset()
      onClose()
      onCreated(res.data.data.id)
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  })

  return (
    <Drawer open={open} onOpenChange={(o) => { if (!o) reset(); onClose() }} title="New Branch">
      <form onSubmit={onSubmit} className="space-y-5">
        <TextField label="Name" required placeholder="e.g. Main Branch" error={errors.name?.message} {...register('name')} />
        <TextField label="Address" placeholder="e.g. Garden City Mall, Kampala" {...register('address')} />
        <FormActions onCancel={onClose} submitLabel="Create Branch" />
      </form>
    </Drawer>
  )
}

function EditBranchDrawer({ open, onClose, branch }: { open: boolean; onClose: () => void; branch: any }) {
  const qc = useQueryClient()
  const { register, handleSubmit, formState: { errors } } = useForm<BranchForm>({
    resolver: zodResolver(branchSchema),
    defaultValues: { name: branch.name, address: branch.address || '' },
  })

  const onSubmit = handleSubmit(async (data) => {
    try {
      await api.put(`/branches/${branch.id}`, data)
      qc.invalidateQueries({ queryKey: ['branches'] })
      toast.success('Updated')
      onClose()
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  })

  return (
    <Drawer open={open} onOpenChange={(o) => !o && onClose()} title="Edit Branch">
      <form onSubmit={onSubmit} className="space-y-5">
        <TextField label="Name" required error={errors.name?.message} {...register('name')} />
        <TextField label="Address" {...register('address')} />
        <FormActions onCancel={onClose} submitLabel="Save changes" />
      </form>
    </Drawer>
  )
}
