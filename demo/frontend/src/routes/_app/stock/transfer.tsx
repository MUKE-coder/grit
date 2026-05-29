import { createFileRoute, Link } from '@tanstack/react-router'
import { useState } from 'react'
import { ArrowLeft, ArrowRight, Check } from 'lucide-react'
import toast from 'react-hot-toast'
import { useBranches } from '@/hooks/useBusiness'
import { useProducts } from '@/hooks/useProducts'
import { useStockTransfer } from '@/hooks/useStock'
import { formatUGX } from '@/lib/utils'

export const Route = createFileRoute('/_app/stock/transfer')({
  component: StockTransferPage,
})

function StockTransferPage() {
  const { data: branches } = useBranches()
  const { data: productsData } = useProducts({ per_page: 200 } as any)
  const transfer = useStockTransfer()

  const [form, setForm] = useState({
    from_branch_id: '',
    to_branch_id: '',
    product_id: '',
    quantity: '',
    note: '',
  })
  const [lastResult, setLastResult] = useState<any>(null)

  const products = productsData?.data || []
  const selectedProduct = products.find((p: any) => String(p.id) === form.product_id)
  const fromBranchStock = selectedProduct?.stock?.find((s: any) => String(s.branch_id) === form.from_branch_id)?.quantity || 0

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (form.from_branch_id === form.to_branch_id) {
      toast.error('Cannot transfer to the same branch')
      return
    }
    try {
      const result = await transfer.mutateAsync({
        product_id: Number(form.product_id),
        from_branch_id: Number(form.from_branch_id),
        to_branch_id: Number(form.to_branch_id),
        quantity: Number(form.quantity),
        note: form.note || undefined,
      })
      setLastResult(result)
      const fromBranch = branches?.find((b) => String(b.id) === form.from_branch_id)?.name
      const toBranch = branches?.find((b) => String(b.id) === form.to_branch_id)?.name
      toast.success(`Transferred ${form.quantity} units from ${fromBranch} to ${toBranch}`)
      setForm({ ...form, quantity: '', note: '' })
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Transfer failed')
    }
  }

  const set = (field: string) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) =>
    setForm({ ...form, [field]: e.target.value })

  return (
    <div className="p-6 lg:p-8 max-w-2xl">
      <Link to="/stock" className="inline-flex items-center gap-1.5 text-sm text-text-muted hover:text-text-primary mb-4">
        <ArrowLeft size={16} /> Back to Stock
      </Link>
      <h1 className="text-2xl font-bold text-text-primary mb-6">Stock Transfer</h1>

      {lastResult && (
        <div className="mb-6 p-4 rounded-xl bg-grit-blue-50 border border-grit-blue/20 flex items-center gap-3">
          <Check size={20} className="text-grit-blue" />
          <div>
            <p className="text-sm font-semibold text-grit-blue">Transfer complete</p>
            <p className="text-xs text-grit-blue-dark">
              Source: {lastResult.from_branch_stock} units | Destination: {lastResult.to_branch_stock} units
            </p>
          </div>
        </div>
      )}

      <form onSubmit={handleSubmit} className="bg-surface rounded-xl border border-border p-5 space-y-4">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-text-primary mb-1.5">From Branch</label>
            <select required value={form.from_branch_id} onChange={set('from_branch_id')}
              className="w-full px-3 py-2.5 rounded-lg border border-border bg-white text-sm focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue">
              <option value="">Select source</option>
              {branches?.map((b) => <option key={b.id} value={b.id}>{b.name}</option>)}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-text-primary mb-1.5">To Branch</label>
            <select required value={form.to_branch_id} onChange={set('to_branch_id')}
              className="w-full px-3 py-2.5 rounded-lg border border-border bg-white text-sm focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue">
              <option value="">Select destination</option>
              {branches?.filter((b) => String(b.id) !== form.from_branch_id).map((b) => (
                <option key={b.id} value={b.id}>{b.name}</option>
              ))}
            </select>
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-text-primary mb-1.5">Product</label>
          <select required value={form.product_id} onChange={set('product_id')}
            className="w-full px-3 py-2.5 rounded-lg border border-border bg-white text-sm focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue">
            <option value="">Select product</option>
            {products.map((p: any) => <option key={p.id} value={p.id}>{p.title}</option>)}
          </select>
          {selectedProduct && form.from_branch_id && (
            <p className="text-xs text-text-muted mt-1">
              Available in source: <span className="font-semibold">{fromBranchStock} units</span>
            </p>
          )}
        </div>

        {/* Preview */}
        {form.product_id && form.from_branch_id && form.to_branch_id && form.quantity && (
          <div className="p-3 rounded-lg bg-background border border-border flex items-center justify-center gap-3 text-sm">
            <span className="font-medium">{branches?.find((b) => String(b.id) === form.from_branch_id)?.name}</span>
            <ArrowRight size={16} className="text-grit-blue" />
            <span className="font-medium">{branches?.find((b) => String(b.id) === form.to_branch_id)?.name}</span>
            <span className="text-text-muted">|</span>
            <span className="font-semibold text-grit-blue">{form.quantity} units</span>
          </div>
        )}

        <div>
          <label className="block text-sm font-medium text-text-primary mb-1.5">Quantity</label>
          <input type="number" required min="1" max={fromBranchStock || undefined} value={form.quantity} onChange={set('quantity')}
            className="w-full px-3 py-2.5 rounded-lg border border-border bg-white text-sm focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue max-w-48"
            placeholder="10" />
          {form.quantity && Number(form.quantity) > fromBranchStock && fromBranchStock > 0 && (
            <p className="text-xs text-error mt-1">Exceeds available stock ({fromBranchStock})</p>
          )}
        </div>

        <div>
          <label className="block text-sm font-medium text-text-primary mb-1.5">
            Note <span className="text-text-light font-normal">— optional</span>
          </label>
          <textarea rows={2} value={form.note} onChange={set('note')}
            className="w-full px-3 py-2.5 rounded-lg border border-border bg-white text-sm resize-none focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue"
            placeholder="e.g. Rebalance stock for weekend rush" />
        </div>

        <button type="submit" disabled={transfer.isPending}
          className="w-full py-3 rounded-xl bg-accent text-white font-semibold text-sm hover:bg-accent-hover transition disabled:opacity-50">
          {transfer.isPending ? 'Transferring...' : 'Confirm & Transfer'}
        </button>
      </form>
    </div>
  )
}
