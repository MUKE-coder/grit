import { createFileRoute, Link } from '@tanstack/react-router'
import { useState } from 'react'
import { ArrowLeft, Check } from 'lucide-react'
import toast from 'react-hot-toast'
import { useBranches } from '@/hooks/useBusiness'
import { useProducts } from '@/hooks/useProducts'
import { useStockIn } from '@/hooks/useStock'
import { formatUGX } from '@/lib/utils'

export const Route = createFileRoute('/_app/stock/in')({
  component: StockInPage,
})

function StockInPage() {
  const { data: branches } = useBranches()
  const { data: productsData } = useProducts({ per_page: 200 } as any)
  const stockIn = useStockIn()

  const [form, setForm] = useState({
    branch_id: '',
    product_id: '',
    quantity: '',
    cost_price_override: '',
    note: '',
  })
  const [lastResult, setLastResult] = useState<any>(null)

  const products = productsData?.data || []
  const selectedProduct = products.find((p: any) => String(p.id) === form.product_id)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const result = await stockIn.mutateAsync({
        product_id: Number(form.product_id),
        branch_id: Number(form.branch_id),
        quantity: Number(form.quantity),
        cost_price_override: form.cost_price_override ? Number(form.cost_price_override) : undefined,
        note: form.note || undefined,
      })
      setLastResult(result)
      toast.success(`Added ${form.quantity} units. Current stock: ${result.current_stock}`)
      setForm({ ...form, product_id: '', quantity: '', cost_price_override: '', note: '' })
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed to record stock-in')
    }
  }

  const set = (field: string) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) =>
    setForm({ ...form, [field]: e.target.value })

  return (
    <div className="p-6 lg:p-8 max-w-2xl">
      <Link to="/stock" className="inline-flex items-center gap-1.5 text-sm text-text-muted hover:text-text-primary mb-4">
        <ArrowLeft size={16} /> Back to Stock
      </Link>
      <h1 className="text-2xl font-bold text-text-primary mb-6">Record Stock In</h1>

      {lastResult && (
        <div className="mb-6 p-4 rounded-xl bg-grit-blue-50 border border-grit-blue/20 flex items-center gap-3">
          <Check size={20} className="text-grit-blue" />
          <div>
            <p className="text-sm font-semibold text-grit-blue">Stock recorded successfully</p>
            <p className="text-xs text-grit-blue-dark">
              Added {lastResult.quantity_added} units. Current stock: {lastResult.current_stock}
            </p>
          </div>
        </div>
      )}

      <form onSubmit={handleSubmit} className="bg-surface rounded-xl border border-border p-5 space-y-4">
        <div>
          <label className="block text-sm font-medium text-text-primary mb-1.5">Branch</label>
          <select required value={form.branch_id} onChange={set('branch_id')}
            className="w-full px-3 py-2.5 rounded-lg border border-border bg-white text-sm focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue">
            <option value="">Select branch</option>
            {branches?.map((b) => <option key={b.id} value={b.id}>{b.name}</option>)}
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-text-primary mb-1.5">Product</label>
          <select required value={form.product_id} onChange={set('product_id')}
            className="w-full px-3 py-2.5 rounded-lg border border-border bg-white text-sm focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue">
            <option value="">Select product</option>
            {products.map((p: any) => (
              <option key={p.id} value={p.id}>
                {p.title} — {formatUGX(p.selling_price)}
              </option>
            ))}
          </select>
          {selectedProduct && (
            <p className="text-xs text-text-muted mt-1">
              Cost: {selectedProduct.cost_price ? formatUGX(selectedProduct.cost_price) : 'Not set'} |
              Current stock: {selectedProduct.stock?.map((s: any) => `${s.branch_name}: ${s.quantity}`).join(', ') || 'None'}
            </p>
          )}
        </div>

        <div>
          <label className="block text-sm font-medium text-text-primary mb-1.5">Quantity</label>
          <input type="number" required min="1" value={form.quantity} onChange={set('quantity')}
            className="w-full px-3 py-2.5 rounded-lg border border-border bg-white text-sm focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue max-w-48"
            placeholder="50" />
        </div>

        <div>
          <label className="block text-sm font-medium text-text-primary mb-1.5">
            Cost Price Override (UGX) <span className="text-text-light font-normal">— optional</span>
          </label>
          <input type="number" min="0" step="100" value={form.cost_price_override} onChange={set('cost_price_override')}
            className="w-full px-3 py-2.5 rounded-lg border border-border bg-white text-sm focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue max-w-48"
            placeholder={selectedProduct?.cost_price ? String(selectedProduct.cost_price) : 'e.g. 180000'} />
        </div>

        <div>
          <label className="block text-sm font-medium text-text-primary mb-1.5">
            Note <span className="text-text-light font-normal">— optional</span>
          </label>
          <textarea rows={2} value={form.note} onChange={set('note')}
            className="w-full px-3 py-2.5 rounded-lg border border-border bg-white text-sm resize-none focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue"
            placeholder="e.g. China shipment batch #42" />
        </div>

        <button type="submit" disabled={stockIn.isPending}
          className="w-full py-3 rounded-xl bg-accent text-white font-semibold text-sm hover:bg-accent-hover transition disabled:opacity-50">
          {stockIn.isPending ? 'Recording...' : 'Record Stock In'}
        </button>
      </form>
    </div>
  )
}
