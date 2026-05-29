import { createFileRoute, Link } from '@tanstack/react-router'
import { useState } from 'react'
import { PackagePlus, ArrowLeftRight, Boxes, AlertTriangle, PackageX, Check } from 'lucide-react'
import toast from 'react-hot-toast'
import { useStockLevels, useStockIn, useStockTransfer } from '@/hooks/useStock'
import { useBranches } from '@/hooks/useBusiness'
import { useCategories, useProducts } from '@/hooks/useProducts'
import { useAuthStore } from '@/stores/auth.store'
import { StatCard } from '@/components/reports/StatCard'
import { Modal } from '@/components/ui/Modal'
import { SearchableSelect } from '@/components/ui/SearchableSelect'
import { formatUGX } from '@/lib/utils'

export const Route = createFileRoute('/_app/stock/')({
  component: StockOverviewPage,
})

function StockOverviewPage() {
  const role = useAuthStore((s) => s.currentRole)
  const [branchId, setBranchId] = useState('')
  const [categoryId, setCategoryId] = useState('')
  const [status, setStatus] = useState('')

  // Modals
  const [showStockIn, setShowStockIn] = useState(false)
  const [showTransfer, setShowTransfer] = useState(false)

  const { data: branches } = useBranches()
  const { data: categories } = useCategories()
  const { data, isLoading } = useStockLevels({
    branch_id: branchId || undefined,
    category_id: categoryId || undefined,
    status: status || undefined,
  })

  const allData = data || []
  const totalSKUs = new Set(allData.map((s: any) => s.product_id)).size
  const lowItems = allData.filter((s: any) => s.status === 'low').length
  const outItems = allData.filter((s: any) => s.status === 'out').length

  const statusBadge = (s: string) => {
    switch (s) {
      case 'ok': return 'bg-success/10 text-success'
      case 'low': return 'bg-warning/10 text-warning'
      case 'out': return 'bg-error/10 text-error'
      default: return 'bg-background text-text-muted'
    }
  }

  return (
    <div className="p-6 lg:p-8">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
        <div>
          <h1 className="text-2xl font-bold text-text-primary">Stock Overview</h1>
          <p className="text-text-muted mt-1">Inventory levels across branches</p>
        </div>
        <div className="flex gap-2">
          <button onClick={() => setShowStockIn(true)}
            className="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-accent text-white font-semibold text-sm hover:bg-accent-hover transition">
            <PackagePlus size={16} /> Stock In
          </button>
          {role === 'admin' && (
            <button onClick={() => setShowTransfer(true)}
              className="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl border border-border text-text-primary font-semibold text-sm hover:bg-background transition">
              <ArrowLeftRight size={16} /> Transfer
            </button>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <StatCard icon={<Boxes size={20} />} label="Total SKUs" value={String(totalSKUs)} color="bg-role-manager" />
        <StatCard icon={<AlertTriangle size={20} />} label="Low Stock" value={String(lowItems)} color="bg-warning" />
        <StatCard icon={<PackageX size={20} />} label="Out of Stock" value={String(outItems)} color="bg-error" />
      </div>

      {/* Filters */}
      <div className="flex flex-wrap gap-2 mb-4">
        {branches && branches.length > 1 && (
          <select value={branchId} onChange={(e) => setBranchId(e.target.value)}
            className="px-3 py-1.5 rounded-lg border border-border text-sm bg-white">
            <option value="">All Branches</option>
            {branches.map((b) => <option key={b.id} value={b.id}>{b.name}</option>)}
          </select>
        )}
        <select value={categoryId} onChange={(e) => setCategoryId(e.target.value)}
          className="px-3 py-1.5 rounded-lg border border-border text-sm bg-white">
          <option value="">All Categories</option>
          {categories?.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
        </select>
        <div className="flex rounded-lg border border-border overflow-hidden">
          {[{v:'',l:'All'},{v:'low',l:'Low'},{v:'out',l:'Out'}].map(({v,l}) => (
            <button key={v} onClick={() => setStatus(v)}
              className={`px-3 py-1.5 text-sm font-medium transition ${status === v ? 'bg-accent text-white' : 'bg-white text-text-muted hover:bg-background'}`}>
              {l}
            </button>
          ))}
        </div>
      </div>

      {/* Table */}
      <div className="bg-surface rounded-xl border border-border shadow-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border bg-background/50">
                <th className="text-left px-5 py-3 font-semibold text-text-muted">Product</th>
                <th className="text-left px-5 py-3 font-semibold text-text-muted">Category</th>
                <th className="text-left px-5 py-3 font-semibold text-text-muted">Branch</th>
                <th className="text-right px-5 py-3 font-semibold text-text-muted">Quantity</th>
                <th className="text-right px-5 py-3 font-semibold text-text-muted">Threshold</th>
                <th className="text-right px-5 py-3 font-semibold text-text-muted">Value</th>
                <th className="text-center px-5 py-3 font-semibold text-text-muted">Status</th>
              </tr>
            </thead>
            <tbody>
              {isLoading ? (
                <tr><td colSpan={7} className="px-5 py-10 text-center text-text-muted">Loading...</td></tr>
              ) : allData.length === 0 ? (
                <tr><td colSpan={7} className="px-5 py-10 text-center text-text-muted">No stock recorded. Click "Stock In" to start.</td></tr>
              ) : (
                allData.map((item: any, i: number) => (
                  <tr key={i} className="border-b border-border last:border-0 hover:bg-background/30">
                    <td className="px-5 py-3 font-medium text-text-primary">{item.product_title}</td>
                    <td className="px-5 py-3 text-text-muted">{item.category_name}</td>
                    <td className="px-5 py-3 text-text-muted">{item.branch_name}</td>
                    <td className="px-5 py-3 text-right font-semibold text-text-primary">{item.quantity}</td>
                    <td className="px-5 py-3 text-right text-text-muted">{item.threshold}</td>
                    <td className="px-5 py-3 text-right text-text-primary">{formatUGX(item.selling_price * item.quantity)}</td>
                    <td className="px-5 py-3 text-center">
                      <span className={`px-2 py-0.5 rounded-full text-xs font-semibold capitalize ${statusBadge(item.status)}`}>
                        {item.status === 'out' ? 'Out' : item.status}
                      </span>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Stock In Modal */}
      <StockInModal open={showStockIn} onClose={() => setShowStockIn(false)} />

      {/* Transfer Modal */}
      {role === 'admin' && (
        <TransferModal open={showTransfer} onClose={() => setShowTransfer(false)} />
      )}
    </div>
  )
}

function StockInModal({ open, onClose }: { open: boolean; onClose: () => void }) {
  const { data: branches } = useBranches()
  const { data: productsData } = useProducts({ per_page: 200 } as any)
  const stockIn = useStockIn()
  const [form, setForm] = useState({ branch_id: '', product_id: '', quantity: '', note: '' })

  const products = productsData?.data || []

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const result = await stockIn.mutateAsync({
        product_id: Number(form.product_id),
        branch_id: Number(form.branch_id),
        quantity: Number(form.quantity),
        note: form.note || undefined,
      })
      toast.success(`Added ${form.quantity} units. Current stock: ${result.current_stock}`)
      setForm({ branch_id: '', product_id: '', quantity: '', note: '' })
      onClose()
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed')
    }
  }

  return (
    <Modal open={open} onClose={onClose} title="Record Stock In">
      <form onSubmit={handleSubmit} className="space-y-4">
        <SearchableSelect
          label="Branch"
          required
          placeholder="Select branch"
          value={form.branch_id}
          onChange={(v) => setForm({ ...form, branch_id: v })}
          options={branches?.map((b) => ({ value: String(b.id), label: b.name })) || []}
        />
        <SearchableSelect
          label="Product"
          required
          placeholder="Search products..."
          value={form.product_id}
          onChange={(v) => setForm({ ...form, product_id: v })}
          options={products.map((p: any) => ({
            value: String(p.id),
            label: p.title,
            subtitle: formatUGX(p.selling_price),
          }))}
        />
        <div>
          <label className="block text-sm font-semibold text-text-primary mb-2">Quantity <span className="text-error">*</span></label>
          <input type="number" required min="1" value={form.quantity} onChange={(e) => setForm({ ...form, quantity: e.target.value })}
            className="w-full px-4 py-3 rounded-xl border border-border bg-white text-sm focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue" placeholder="50" />
        </div>
        <div>
          <label className="block text-sm font-semibold text-text-primary mb-2">Note <span className="text-text-light font-normal">— optional</span></label>
          <input type="text" value={form.note} onChange={(e) => setForm({ ...form, note: e.target.value })}
            className="w-full px-4 py-3 rounded-xl border border-border bg-white text-sm focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue" placeholder="e.g. China shipment" />
        </div>
        <button type="submit" disabled={stockIn.isPending}
          className="w-full py-3 rounded-xl bg-accent text-white font-semibold hover:bg-accent-hover transition disabled:opacity-50">
          {stockIn.isPending ? 'Recording...' : 'Record Stock In'}
        </button>
      </form>
    </Modal>
  )
}

function TransferModal({ open, onClose }: { open: boolean; onClose: () => void }) {
  const { data: branches } = useBranches()
  const { data: productsData } = useProducts({ per_page: 200 } as any)
  const transfer = useStockTransfer()
  const [form, setForm] = useState({ from_branch_id: '', to_branch_id: '', product_id: '', quantity: '', note: '' })

  const products = productsData?.data || []

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await transfer.mutateAsync({
        product_id: Number(form.product_id),
        from_branch_id: Number(form.from_branch_id),
        to_branch_id: Number(form.to_branch_id),
        quantity: Number(form.quantity),
        note: form.note || undefined,
      })
      toast.success('Transfer complete!')
      setForm({ from_branch_id: '', to_branch_id: '', product_id: '', quantity: '', note: '' })
      onClose()
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Transfer failed')
    }
  }

  return (
    <Modal open={open} onClose={onClose} title="Stock Transfer">
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid grid-cols-2 gap-3">
          <SearchableSelect
            label="From"
            required
            placeholder="Source"
            value={form.from_branch_id}
            onChange={(v) => setForm({ ...form, from_branch_id: v })}
            options={branches?.map((b) => ({ value: String(b.id), label: b.name })) || []}
          />
          <SearchableSelect
            label="To"
            required
            placeholder="Destination"
            value={form.to_branch_id}
            onChange={(v) => setForm({ ...form, to_branch_id: v })}
            options={branches?.filter((b) => String(b.id) !== form.from_branch_id).map((b) => ({ value: String(b.id), label: b.name })) || []}
          />
        </div>
        <SearchableSelect
          label="Product"
          required
          placeholder="Search products..."
          value={form.product_id}
          onChange={(v) => setForm({ ...form, product_id: v })}
          options={products.map((p: any) => ({ value: String(p.id), label: p.title }))}
        />
        <div>
          <label className="block text-sm font-semibold text-text-primary mb-2">Quantity <span className="text-error">*</span></label>
          <input type="number" required min="1" value={form.quantity} onChange={(e) => setForm({ ...form, quantity: e.target.value })}
            className="w-full px-4 py-3 rounded-xl border border-border bg-white text-sm focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue" placeholder="10" />
        </div>
        <button type="submit" disabled={transfer.isPending}
          className="w-full py-3 rounded-xl bg-accent text-white font-semibold hover:bg-accent-hover transition disabled:opacity-50">
          {transfer.isPending ? 'Transferring...' : 'Confirm & Transfer'}
        </button>
      </form>
    </Modal>
  )
}
