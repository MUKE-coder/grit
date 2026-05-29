import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useStockMovements } from '@/hooks/useStock'
import { useBranches } from '@/hooks/useBusiness'

export const Route = createFileRoute('/_app/stock/movements')({
  component: StockMovementsPage,
})

function StockMovementsPage() {
  const today = new Date().toISOString().split('T')[0]
  const weekAgo = new Date(Date.now() - 7 * 86400000).toISOString().split('T')[0]

  const [branchId, setBranchId] = useState('')
  const [type, setType] = useState('')
  const [fromDate, setFromDate] = useState(weekAgo)
  const [toDate, setToDate] = useState(today)

  const { data: branches } = useBranches()
  const { data, isLoading } = useStockMovements({
    branch_id: branchId || undefined,
    type: type || undefined,
    from_date: fromDate,
    to_date: toDate,
  } as any)

  const movements = data?.data || []

  const typeBadge = (t: string) => {
    switch (t) {
      case 'stock_in': return { label: 'Stock In', cls: 'bg-grit-blue-50 text-grit-blue' }
      case 'sale': return { label: 'Sale', cls: 'bg-role-manager/10 text-role-manager' }
      case 'transfer_out': return { label: 'Transfer Out', cls: 'bg-error/10 text-error' }
      case 'transfer_in': return { label: 'Transfer In', cls: 'bg-role-admin/10 text-role-admin' }
      default: return { label: t, cls: 'bg-background text-text-muted' }
    }
  }

  return (
    <div className="p-6 lg:p-8">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-text-primary">Stock Movements</h1>
        <p className="text-text-muted mt-1">Full audit trail of all stock changes</p>
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
        <select value={type} onChange={(e) => setType(e.target.value)}
          className="px-3 py-1.5 rounded-lg border border-border text-sm bg-white">
          <option value="">All Types</option>
          <option value="stock_in">Stock In</option>
          <option value="sale">Sale</option>
          <option value="transfer_out">Transfer Out</option>
          <option value="transfer_in">Transfer In</option>
        </select>
        <input type="date" value={fromDate} onChange={(e) => setFromDate(e.target.value)}
          className="px-3 py-1.5 rounded-lg border border-border text-sm bg-white" />
        <span className="text-text-muted text-sm self-center">to</span>
        <input type="date" value={toDate} onChange={(e) => setToDate(e.target.value)}
          className="px-3 py-1.5 rounded-lg border border-border text-sm bg-white" />
      </div>

      <div className="bg-surface rounded-xl border border-border shadow-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border bg-background/50">
                <th className="text-left px-5 py-3 font-semibold text-text-muted">Date</th>
                <th className="text-left px-5 py-3 font-semibold text-text-muted">Product</th>
                <th className="text-left px-5 py-3 font-semibold text-text-muted">Branch</th>
                <th className="text-center px-5 py-3 font-semibold text-text-muted">Type</th>
                <th className="text-right px-5 py-3 font-semibold text-text-muted">Quantity</th>
                <th className="text-left px-5 py-3 font-semibold text-text-muted">Note</th>
              </tr>
            </thead>
            <tbody>
              {isLoading ? (
                <tr><td colSpan={6} className="px-5 py-10 text-center text-text-muted">Loading...</td></tr>
              ) : movements.length === 0 ? (
                <tr><td colSpan={6} className="px-5 py-10 text-center text-text-muted">No movements found</td></tr>
              ) : (
                movements.map((m: any) => {
                  const badge = typeBadge(m.movement_type)
                  return (
                    <tr key={m.id} className="border-b border-border last:border-0 hover:bg-background/30">
                      <td className="px-5 py-3 text-text-muted whitespace-nowrap">
                        {new Date(m.created_at).toLocaleDateString()} {new Date(m.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                      </td>
                      <td className="px-5 py-3 font-medium text-text-primary">{m.product?.title}</td>
                      <td className="px-5 py-3 text-text-muted">{m.branch?.name}</td>
                      <td className="px-5 py-3 text-center">
                        <span className={`px-2 py-0.5 rounded-full text-xs font-semibold ${badge.cls}`}>{badge.label}</span>
                      </td>
                      <td className={`px-5 py-3 text-right font-semibold ${m.quantity > 0 ? 'text-grit-blue' : 'text-error'}`}>
                        {m.quantity > 0 ? '+' : ''}{m.quantity}
                      </td>
                      <td className="px-5 py-3 text-text-muted truncate max-w-48">{m.note || '—'}</td>
                    </tr>
                  )
                })
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
