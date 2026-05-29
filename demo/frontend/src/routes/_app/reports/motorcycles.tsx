import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useMotorcyclesReport } from '@/hooks/useReports'
import { ExportButton } from '@/components/export-button'
import { exportSheetsToExcel } from '@/lib/export'
import { formatUGX } from '@/lib/utils'
import { Bike, Wallet, HandCoins } from 'lucide-react'

export const Route = createFileRoute('/_app/reports/motorcycles')({
  component: MotorcyclesReportPage,
})

function defaultRange() {
  const now = new Date()
  const from = new Date(now); from.setDate(from.getDate() - 30)
  return { from: from.toISOString().slice(0, 10), to: now.toISOString().slice(0, 10) }
}

function MotorcyclesReportPage() {
  const init = defaultRange()
  const [from, setFrom] = useState(init.from)
  const [to, setTo] = useState(init.to)
  const { data, isLoading } = useMotorcyclesReport({ from_date: from, to_date: to })

  if (isLoading) return <div className="p-6 text-text-muted">Loading...</div>

  const inventory: Array<{ branch_id: number; branch_name: string; status: string; count: number }> = data?.inventory || []
  const cashByBranch: Array<{ branch_id: number; branch_name: string; total: number; count: number }> = data?.cash_sales_by_branch || []

  // Pivot inventory into branch -> { status: count }
  const branches: Record<string, Record<string, number>> = {}
  for (const r of inventory) {
    if (!branches[r.branch_name]) branches[r.branch_name] = {}
    branches[r.branch_name][r.status] = r.count
  }
  const statuses = ['available', 'reserved', 'sold', 'on_loan', 'repossessed']

  const handleExport = () => {
    exportSheetsToExcel([
      { name: 'Summary', rows: [
        { Metric: 'Cash Sales Total', Value: data?.cash_sales_total || 0 },
        { Metric: 'Cash Sales Count', Value: data?.cash_sales_count || 0 },
        { Metric: 'Loan-Financed Total', Value: data?.loan_sales_total || 0 },
        { Metric: 'Bikes Moved', Value: data?.motorcycles_moved || 0 },
        { Metric: 'From', Value: from }, { Metric: 'To', Value: to },
      ]},
      { name: 'Inventory by Branch', rows: Object.entries(branches).map(([name, counts]) => ({
        Branch: name,
        ...Object.fromEntries(statuses.map((s) => [s.replace('_', ' '), counts[s] || 0])),
      })) },
      { name: 'Cash Sales by Branch', rows: cashByBranch.map((b) => ({ Branch: b.branch_name, Sales: b.count, Total: b.total })) },
    ], `motorcycle-sales-${from}-to-${to}`)
  }

  return (
    <div className="p-4 lg:p-8">
      <div className="mb-6 flex flex-col sm:flex-row sm:items-start sm:justify-between gap-3">
        <div>
          <h1 className="text-[20px] lg:text-2xl font-bold text-text-primary">Motorcycle Sales</h1>
          <p className="text-text-muted mt-1 text-[13px]">Inventory by branch and sales activity in the selected period.</p>
        </div>
        <div className="flex gap-2 flex-wrap">
          <input type="date" value={from} onChange={(e) => setFrom(e.target.value)} className="px-3 h-8 rounded-lg border border-border bg-white text-[12.5px]" />
          <input type="date" value={to} onChange={(e) => setTo(e.target.value)} className="px-3 h-8 rounded-lg border border-border bg-white text-[12.5px]" />
          <ExportButton onClick={handleExport} />
        </div>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
        <BigStat icon={<Wallet size={20} />} label="Cash Sales Total" value={formatUGX(data?.cash_sales_total || 0)} color="bg-grit-blue" />
        <BigStat icon={<HandCoins size={20} />} label="Loan-Financed" value={formatUGX(data?.loan_sales_total || 0)} color="bg-role-manager" />
        <BigStat icon={<Bike size={20} />} label="Bikes Moved" value={String(data?.motorcycles_moved || 0)} color="bg-role-admin" />
        <BigStat icon={<Bike size={20} />} label="Cash Sales Count" value={String(data?.cash_sales_count || 0)} color="bg-warning" />
      </div>

      <Card title="Inventory by Branch" className="mb-6">
        {Object.keys(branches).length === 0 ? (
          <p className="text-text-muted text-sm">No motorcycles in inventory</p>
        ) : (
          <table className="w-full text-sm">
            <thead className="text-xs text-text-muted uppercase tracking-wider">
              <tr>
                <th className="text-left pb-2">Branch</th>
                {statuses.map((s) => (
                  <th key={s} className="text-right pb-2 capitalize">{s.replace('_', ' ')}</th>
                ))}
                <th className="text-right pb-2">Total</th>
              </tr>
            </thead>
            <tbody>
              {Object.entries(branches).map(([name, counts]) => {
                const total = statuses.reduce((s, st) => s + (counts[st] || 0), 0)
                return (
                  <tr key={name} className="border-t border-border">
                    <td className="py-2 font-medium">{name}</td>
                    {statuses.map((s) => (
                      <td key={s} className="py-2 text-right">{counts[s] || 0}</td>
                    ))}
                    <td className="py-2 text-right font-bold">{total}</td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        )}
      </Card>

      <Card title="Cash Sales by Branch">
        {cashByBranch.length === 0 ? <p className="text-text-muted text-sm">No cash sales in range</p> : (
          <table className="w-full text-sm">
            <thead className="text-xs text-text-muted uppercase tracking-wider">
              <tr><th className="text-left pb-2">Branch</th><th className="text-right pb-2">Sales</th><th className="text-right pb-2">Total</th></tr>
            </thead>
            <tbody>
              {cashByBranch.map((b) => (
                <tr key={b.branch_id} className="border-t border-border">
                  <td className="py-2 font-medium">{b.branch_name}</td>
                  <td className="py-2 text-right">{b.count}</td>
                  <td className="py-2 text-right font-bold">{formatUGX(b.total)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </Card>
    </div>
  )
}

function BigStat({ icon, label, value, color }: { icon: React.ReactNode; label: string; value: string; color: string }) {
  return (
    <div className="bg-surface rounded-xl border border-border p-5 shadow-sm">
      <div className="flex items-center gap-3 mb-3">
        <div className={`w-9 h-9 rounded-lg ${color} text-white flex items-center justify-center`}>{icon}</div>
        <p className="text-sm text-text-muted">{label}</p>
      </div>
      <p className="text-xl font-bold text-text-primary">{value}</p>
    </div>
  )
}

function Card({ title, children, className }: { title: string; children: React.ReactNode; className?: string }) {
  return (
    <div className={`bg-surface rounded-xl border border-border p-5 shadow-sm ${className || ''}`}>
      <h3 className="text-sm font-semibold text-text-primary mb-4">{title}</h3>
      {children}
    </div>
  )
}
