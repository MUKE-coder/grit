import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { Search } from 'lucide-react'
import { useStockReport } from '@/hooks/useReports'
import { useBranches } from '@/hooks/useBusiness'
import { useCategories } from '@/hooks/useProducts'
import { ExportButton } from '@/components/export-button'
import { exportToExcel } from '@/lib/export'
import { formatUGX } from '@/lib/utils'

export const Route = createFileRoute('/_app/reports/stock')({
  component: StockReportPage,
})

function StockReportPage() {
  const [search, setSearch] = useState('')
  const [branchId, setBranchId] = useState('')
  const [categoryId, setCategoryId] = useState('')
  const [status, setStatus] = useState('')

  const { data: branches } = useBranches()
  const { data: categories } = useCategories()
  const { data, isLoading } = useStockReport({
    branch_id: branchId || undefined,
    category_id: categoryId || undefined,
    status: status || undefined,
  })

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
          <h1 className="text-2xl font-bold text-text-primary">Stock Report</h1>
          <p className="text-text-muted mt-1">Inventory levels across all branches</p>
        </div>
        <div className="flex items-center gap-2 flex-wrap">
          <div className="relative">
            <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" />
            <input type="text" value={search} onChange={(e) => setSearch(e.target.value)}
              placeholder="Search products..."
              className="pl-8 pr-3 py-1.5 rounded-lg border border-border text-sm bg-white w-48 focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue" />
          </div>
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
            {['', 'low', 'out'].map((s) => (
              <button key={s} onClick={() => setStatus(s)}
                className={`px-3 py-1.5 text-sm font-medium transition ${status === s ? 'bg-grit-blue text-white' : 'bg-white text-text-muted hover:bg-background'}`}>
                {s === '' ? 'All' : s === 'low' ? 'Low' : 'Out'}
              </button>
            ))}
          </div>
          <ExportButton
            onClick={() => exportToExcel(
              (data || []).filter((item: any) => !search || item.product_title.toLowerCase().includes(search.toLowerCase()))
                .map((item: any) => ({
                  Product: item.product_title,
                  Category: item.category_name,
                  Branch: item.branch_name,
                  Quantity: item.quantity,
                  Threshold: item.threshold,
                  'Stock Value': item.stock_value,
                  Status: item.status,
                })),
              'stock-report',
              'Stock',
            )}
            disabled={!data || data.length === 0}
          />
        </div>
      </div>

      <div className="bg-surface rounded-xl border border-border shadow-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border bg-background/50">
                <th className="text-left px-5 py-3 font-semibold text-text-muted">Product</th>
                <th className="text-left px-5 py-3 font-semibold text-text-muted">Category</th>
                <th className="text-left px-5 py-3 font-semibold text-text-muted">Branch</th>
                <th className="text-right px-5 py-3 font-semibold text-text-muted">Qty</th>
                <th className="text-right px-5 py-3 font-semibold text-text-muted">Threshold</th>
                <th className="text-right px-5 py-3 font-semibold text-text-muted">Value</th>
                <th className="text-center px-5 py-3 font-semibold text-text-muted">Status</th>
              </tr>
            </thead>
            <tbody>
              {isLoading ? (
                <tr><td colSpan={7} className="px-5 py-10 text-center text-text-muted">Loading...</td></tr>
              ) : !data || data.length === 0 ? (
                <tr><td colSpan={7} className="px-5 py-10 text-center text-text-muted">No stock data</td></tr>
              ) : (
                data
                  .filter((item: any) => !search || item.product_title.toLowerCase().includes(search.toLowerCase()))
                  .map((item: any, i: number) => (
                  <tr key={i} className="border-b border-border last:border-0 hover:bg-background/30">
                    <td className="px-5 py-3 font-medium text-text-primary">{item.product_title}</td>
                    <td className="px-5 py-3 text-text-muted">{item.category_name}</td>
                    <td className="px-5 py-3 text-text-muted">{item.branch_name}</td>
                    <td className="px-5 py-3 text-right font-semibold text-text-primary">{item.quantity}</td>
                    <td className="px-5 py-3 text-right text-text-muted">{item.threshold}</td>
                    <td className="px-5 py-3 text-right text-text-primary">{formatUGX(item.stock_value)}</td>
                    <td className="px-5 py-3 text-center">
                      <span className={`px-2 py-0.5 rounded-full text-xs font-semibold capitalize ${statusBadge(item.status)}`}>
                        {item.status === 'out' ? 'Out of Stock' : item.status}
                      </span>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
