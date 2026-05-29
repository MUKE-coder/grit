import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Package, Pencil, PackagePlus, Plus, Upload } from 'lucide-react'
import toast from 'react-hot-toast'
import { useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'

import {
  useProducts, useProduct, useCategories,
  useCreateProduct, useCreateCategory,
} from '@/hooks/useProducts'
import { useBranches } from '@/hooks/useBusiness'
import { useStockIn } from '@/hooks/useStock'
import {
  TwoPane, ListPane, ListRow, DetailPane, DetailSection, Field, FieldGrid,
} from '@/components/two-pane'
import { Drawer } from '@/components/drawer'
import { TextField, SelectField, FormActions, FormGrid, FormSection, TextAreaField } from '@/components/form'
import { CurrencyField } from '@/components/currency-field'
import { StatusBadge } from '@/components/status-badge'
import { ExportButton } from '@/components/export-button'
import { exportToExcel } from '@/lib/export'
import { formatUGX } from '@/lib/utils'

type PSearch = { selected?: number; q?: string; cat?: string }

export const Route = createFileRoute('/_app/products/')({
  component: ProductsPage,
  validateSearch: (s: Record<string, unknown>): PSearch => ({
    selected: s.selected ? Number(s.selected) : undefined,
    q: typeof s.q === 'string' ? s.q : undefined,
    cat: typeof s.cat === 'string' ? s.cat : undefined,
  }),
})

import { ImportProductsDrawer } from '@/components/products/ImportProductsDrawer'

function ProductsPage() {
  const search = Route.useSearch()
  const navigate = useNavigate()
  const [createOpen, setCreateOpen] = useState(false)
  const [importOpen, setImportOpen] = useState(false)

  const list = useProducts({ search: search.q, category_id: search.cat })
  const cats = useCategories()
  const selected = useProduct(search.selected || 0)

  const update = (patch: Partial<PSearch>) =>
    navigate({ to: '/products', search: { ...search, ...patch }, replace: true })

  const items = (list.data?.data as any[] | undefined) || []

  return (
    <TwoPane>
      <ListPane
        title="Spare Products"
        count={list.data?.total}
        search={search.q || ''}
        onSearch={(q) => update({ q: q || undefined })}
        searchPlaceholder="Title or barcode..."
        onNew={() => setCreateOpen(true)}
        refreshKey={['products']}
        extraAction={
          <>
            <button
              type="button"
              onClick={() => setImportOpen(true)}
              className="h-8 px-2.5 rounded-lg border border-border bg-surface text-[12px] font-medium text-foreground-secondary hover:bg-surface-hover transition inline-flex items-center gap-1"
              title="Import products from Excel"
            >
              <Upload size={12} /> <span className="hidden sm:inline">Import</span>
            </button>
            <ExportButton
              onClick={() => exportToExcel(
                (list.data?.data || []).map((p: any) => {
                  const totalStock = (p.stock || []).reduce((sum: number, s: any) => sum + s.quantity, 0)
                  return {
                    Title: p.title,
                    Barcode: p.barcode,
                    Category: p.category?.name,
                    'Cost Price': p.cost_price,
                    'Selling Price': p.selling_price,
                    'Total Stock': totalStock,
                    'Low Stock Threshold': p.low_stock_threshold,
                    Description: p.description,
                  }
                }),
                'products',
                'Products',
              )}
              disabled={!list.data?.data.length}
            />
          </>
        }
        filters={
          <div className="flex gap-1 flex-wrap">
            <button
              type="button"
              onClick={() => update({ cat: undefined })}
              className={`px-2.5 py-1 rounded-md text-[11.5px] font-medium transition ${
                !search.cat ? 'bg-foreground text-white' : 'bg-surface-2 text-foreground-secondary hover:bg-surface-hover'
              }`}
            >
              All
            </button>
            {cats.data?.slice(0, 8).map((c) => (
              <button
                key={c.id}
                type="button"
                onClick={() => update({ cat: String(c.id) })}
                className={`px-2.5 py-1 rounded-md text-[11.5px] font-medium transition ${
                  search.cat === String(c.id) ? 'bg-foreground text-white' : 'bg-surface-2 text-foreground-secondary hover:bg-surface-hover'
                }`}
              >
                {c.name}
              </button>
            ))}
          </div>
        }
      >
        {list.isLoading && <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">Loading…</div>}
        {!list.isLoading && items.length === 0 && (
          <div className="px-4 py-8 text-center text-[12.5px] text-foreground-muted">No products yet.</div>
        )}
        {items.map((p: any) => {
          const totalStock = p.stock?.reduce((s: number, r: any) => s + r.quantity, 0) || 0
          const status = totalStock === 0 ? 'out' : totalStock <= (p.low_stock_threshold || 5) ? 'low_stock' : 'ok'
          return (
            <ListRow
              key={p.id}
              selected={search.selected === p.id}
              onClick={() => update({ selected: p.id })}
              icon={p.image_url ? <img src={p.image_url} alt="" className="w-full h-full object-cover" /> : <Package size={14} className="text-foreground-secondary" />}
              title={p.title}
              subtitle={p.category?.name}
              rightTop={<span className="tabular-nums">{formatUGX(p.selling_price)}</span>}
              rightBottom={<StatusBadge tone={status === 'out' ? 'danger' : status === 'low_stock' ? 'warning' : 'success'} size="sm">{totalStock} in stock</StatusBadge>}
            />
          )
        })}
      </ListPane>

      {!search.selected ? (
        <DetailPane empty emptyTitle="No product selected" emptyHint="Pick a product from the list, or click + New to add one." />
      ) : selected.isLoading || !selected.data ? (
        <DetailPane><div className="text-foreground-muted">Loading…</div></DetailPane>
      ) : (
        <ProductDetail product={selected.data} />
      )}

      <CreateProductDrawer open={createOpen} onClose={() => setCreateOpen(false)} onCreated={(id) => update({ selected: id })} />
      <ImportProductsDrawer open={importOpen} onClose={() => setImportOpen(false)} />
    </TwoPane>
  )
}

function ProductDetail({ product }: { product: any }) {
  const [editOpen, setEditOpen] = useState(false)
  const [stockInOpen, setStockInOpen] = useState(false)
  const totalStock = product.stock?.reduce((s: number, r: any) => s + r.quantity, 0) || 0

  return (
    <DetailPane
      header={
        <div className="flex items-center justify-between gap-4">
          <div className="flex items-center gap-3 min-w-0">
            <div className="h-12 w-12 rounded-xl bg-surface-2 flex items-center justify-center shrink-0 overflow-hidden">
              {product.image_url ? <img src={product.image_url} alt={product.title} className="w-full h-full object-cover" /> : <Package size={20} className="text-foreground-muted" />}
            </div>
            <div className="min-w-0">
              <h1 className="text-[17px] font-semibold text-foreground truncate">{product.title}</h1>
              <p className="text-[12.5px] text-foreground-muted">{product.category?.name} {product.barcode && <>· <span className="font-mono">{product.barcode}</span></>}</p>
            </div>
          </div>
          <div className="flex items-center gap-2 shrink-0">
            <button onClick={() => setStockInOpen(true)} className="h-8 px-2.5 rounded-lg bg-accent text-white text-[12px] font-medium hover:bg-accent-hover transition inline-flex items-center gap-1">
              <PackagePlus size={12} /> Stock In
            </button>
            <button onClick={() => setEditOpen(true)} className="h-8 px-2.5 rounded-lg border border-border bg-surface text-[12px] font-medium text-foreground-secondary hover:bg-surface-hover transition inline-flex items-center gap-1">
              <Pencil size={12} /> Edit
            </button>
          </div>
        </div>
      }
    >
      <DetailSection title="Pricing">
        <FieldGrid columns={3}>
          <Field label="Selling price" value={<span className="font-semibold">{formatUGX(product.selling_price)}</span>} />
          <Field label="Cost price" value={product.cost_price > 0 ? formatUGX(product.cost_price) : null} />
          <Field label="Margin" value={product.cost_price > 0 ? `${Math.round(((product.selling_price - product.cost_price) / product.selling_price) * 100)}%` : null} />
        </FieldGrid>
      </DetailSection>

      <DetailSection title={`Stock by Branch · ${totalStock} total units`}>
        {!product.stock?.length ? (
          <p className="text-[13px] text-foreground-muted">No stock recorded yet. Click <span className="font-medium">Stock In</span> to receive the first units.</p>
        ) : (
          <div className="border border-border rounded-lg overflow-hidden">
            <table className="w-full text-[13px] tabular-nums">
              <thead className="bg-surface-2 text-[11px] uppercase tracking-wider text-foreground-muted">
                <tr>
                  <th className="text-left py-2 px-3 font-semibold">Branch</th>
                  <th className="text-right py-2 px-3 font-semibold">Quantity</th>
                  <th className="text-right py-2 px-3 font-semibold">Value</th>
                  <th className="text-left py-2 px-3 font-semibold">Status</th>
                </tr>
              </thead>
              <tbody>
                {product.stock.map((s: any) => {
                  const status = s.quantity === 0 ? 'out' : s.quantity <= (product.low_stock_threshold || 5) ? 'low_stock' : 'ok'
                  return (
                    <tr key={`${s.branch_id}-${s.product_id}`} className="border-t border-border-subtle">
                      <td className="py-2 px-3">{s.branch_name}</td>
                      <td className="py-2 px-3 text-right font-medium">{s.quantity}</td>
                      <td className="py-2 px-3 text-right">{formatUGX(s.quantity * product.selling_price)}</td>
                      <td className="py-2 px-3">
                        <StatusBadge tone={status === 'out' ? 'danger' : status === 'low_stock' ? 'warning' : 'success'} size="sm">
                          {status === 'out' ? 'Out' : status === 'low_stock' ? 'Low' : 'OK'}
                        </StatusBadge>
                      </td>
                    </tr>
                  )
                })}
              </tbody>
            </table>
          </div>
        )}
      </DetailSection>

      {product.description && (
        <DetailSection title="Description">
          <p className="text-[13px] text-foreground whitespace-pre-line">{product.description}</p>
        </DetailSection>
      )}

      <EditProductDrawer open={editOpen} onClose={() => setEditOpen(false)} product={product} />
      <StockInDrawer open={stockInOpen} onClose={() => setStockInOpen(false)} productId={product.id} />
    </DetailPane>
  )
}

/* ─────────────── Create drawer ─────────────── */

const createSchema = z.object({
  title: z.string().min(1, 'Required'),
  selling_price: z.number({ message: 'Required' }).positive(),
  category_id: z.coerce.number().int().positive('Required'),
})
type CreateForm = z.infer<typeof createSchema>

function CreateProductDrawer({ open, onClose, onCreated }: { open: boolean; onClose: () => void; onCreated: (id: number) => void }) {
  const cats = useCategories()
  const create = useCreateProduct()
  const createCat = useCreateCategory()
  const [newCat, setNewCat] = useState('')
  const [showNewCat, setShowNewCat] = useState(false)
  const { register, handleSubmit, formState: { errors }, reset, control, setValue } = useForm<CreateForm>({
    resolver: zodResolver(createSchema),
    defaultValues: { title: '', selling_price: undefined as any, category_id: undefined as any },
  })

  const onSubmit = handleSubmit(async (data) => {
    try {
      const created: any = await create.mutateAsync(data as any)
      toast.success('Product created')
      reset()
      onClose()
      onCreated(created.id)
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  })

  const handleNewCat = async () => {
    if (!newCat.trim()) return
    try {
      const c: any = await createCat.mutateAsync(newCat.trim())
      setValue('category_id', c.id)
      setNewCat('')
      setShowNewCat(false)
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  }

  return (
    <Drawer open={open} onOpenChange={(o) => { if (!o) reset(); onClose() }} title="New Product" description="Just title, price, and category. Add cost / barcode after.">
      <form onSubmit={onSubmit} className="space-y-5">
        <TextField label="Title" required placeholder="e.g. Brake pads" error={errors.title?.message} {...register('title')} />
        <CurrencyField control={control} name="selling_price" label="Selling price" required />

        {!showNewCat ? (
          <div className="flex gap-2 items-end">
            <div className="flex-1">
              <SelectField
                label="Category"
                required
                placeholder="Select category"
                options={cats.data?.map((c) => ({ value: String(c.id), label: c.name })) || []}
                error={errors.category_id?.message}
                {...register('category_id')}
              />
            </div>
            <button type="button" onClick={() => setShowNewCat(true)} className="h-10 px-3 rounded-lg border border-dashed border-accent text-accent text-[12px] font-medium hover:bg-accent-tint">
              <Plus size={14} />
            </button>
          </div>
        ) : (
          <div className="flex gap-2">
            <input
              autoFocus
              value={newCat}
              onChange={(e) => setNewCat(e.target.value)}
              placeholder="Category name"
              className="flex-1 h-10 px-3 rounded-lg border border-border bg-surface text-[13.5px] focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/15"
            />
            <button type="button" onClick={handleNewCat} className="h-10 px-3 rounded-lg bg-accent text-white text-[12px] font-medium">Add</button>
            <button type="button" onClick={() => { setShowNewCat(false); setNewCat('') }} className="h-10 px-3 rounded-lg border border-border bg-surface text-[12px] font-medium">Cancel</button>
          </div>
        )}

        <FormActions onCancel={onClose} submitLabel="Create Product" isPending={create.isPending} />
      </form>
    </Drawer>
  )
}

/* ─────────────── Edit drawer ─────────────── */

const editSchema = z.object({
  title: z.string().min(1, 'Required'),
  selling_price: z.number().positive(),
  cost_price: z.number().optional(),
  category_id: z.coerce.number().int().positive(),
  barcode: z.string().optional(),
  description: z.string().optional(),
  low_stock_threshold: z.coerce.number().int().min(0),
})
type EditForm = z.infer<typeof editSchema>

function EditProductDrawer({ open, onClose, product }: { open: boolean; onClose: () => void; product: any }) {
  const cats = useCategories()
  const queryClient = useQueryClient()
  const { register, handleSubmit, formState: { errors }, control } = useForm<EditForm>({
    resolver: zodResolver(editSchema),
    defaultValues: {
      title: product.title,
      selling_price: product.selling_price,
      cost_price: product.cost_price || undefined,
      category_id: product.category_id,
      barcode: product.barcode || '',
      description: product.description || '',
      low_stock_threshold: product.low_stock_threshold || 5,
    },
  })

  const onSubmit = handleSubmit(async (data) => {
    try {
      await api.put(`/products/${product.id}`, data)
      queryClient.invalidateQueries({ queryKey: ['products'] })
      queryClient.invalidateQueries({ queryKey: ['product', String(product.id)] })
      toast.success('Updated')
      onClose()
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  })

  return (
    <Drawer open={open} onOpenChange={(o) => !o && onClose()} title="Edit Product" width={520}>
      <form onSubmit={onSubmit} className="space-y-6">
        <FormSection title="Identity">
          <TextField label="Title" required error={errors.title?.message} {...register('title')} />
          <SelectField
            label="Category"
            required
            options={cats.data?.map((c) => ({ value: String(c.id), label: c.name })) || []}
            error={errors.category_id?.message}
            {...register('category_id')}
          />
          <TextField label="Barcode" {...register('barcode')} />
        </FormSection>

        <FormSection title="Pricing">
          <FormGrid columns={2}>
            <CurrencyField control={control} name="selling_price" label="Selling price" required />
            <CurrencyField control={control} name="cost_price" label="Cost price" />
          </FormGrid>
        </FormSection>

        <FormSection title="Inventory">
          <TextField label="Low stock threshold" type="number" hint="Alert when stock at any branch falls to this level." {...register('low_stock_threshold')} />
        </FormSection>

        <TextAreaField label="Description" {...register('description')} />

        <FormActions onCancel={onClose} submitLabel="Save changes" />
      </form>
    </Drawer>
  )
}

/* ─────────────── Stock In drawer ─────────────── */

const stockInSchema = z.object({
  branch_id: z.coerce.number().int().positive('Required'),
  quantity: z.coerce.number().int().positive('Must be positive'),
  note: z.string().optional(),
})
type StockInForm = z.infer<typeof stockInSchema>

function StockInDrawer({ open, onClose, productId }: { open: boolean; onClose: () => void; productId: number }) {
  const { data: branches } = useBranches()
  const stockIn = useStockIn()
  const { register, handleSubmit, formState: { errors }, reset } = useForm<StockInForm>({
    resolver: zodResolver(stockInSchema),
    defaultValues: { branch_id: undefined as any, quantity: undefined as any, note: '' },
  })

  const onSubmit = handleSubmit(async (data) => {
    try {
      await stockIn.mutateAsync({ product_id: productId, ...data })
      toast.success(`+${data.quantity} units added`)
      reset()
      onClose()
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  })

  return (
    <Drawer open={open} onOpenChange={(o) => { if (!o) reset(); onClose() }} title="Stock In">
      <form onSubmit={onSubmit} className="space-y-5">
        <SelectField
          label="Branch"
          required
          placeholder="Receiving at..."
          options={branches?.map((b) => ({ value: String(b.id), label: b.name })) || []}
          error={errors.branch_id?.message}
          {...register('branch_id')}
        />
        <TextField label="Quantity" required type="number" min="1" error={errors.quantity?.message} {...register('quantity')} />
        <TextAreaField label="Note (optional)" placeholder="Supplier, invoice ref..." {...register('note')} />
        <FormActions onCancel={onClose} submitLabel="Receive Stock" isPending={stockIn.isPending} />
      </form>
    </Drawer>
  )
}
