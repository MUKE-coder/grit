import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useCashierReport } from '@/hooks/useReports'
import { useBranches } from '@/hooks/useBusiness'
import { ExportButton } from '@/components/export-button'
import { exportToExcel } from '@/lib/export'
import { formatUGX } from '@/lib/utils'

export const Route = createFileRoute('/_app/reports/cashiers')({
  component: CashierReportPage,
})

function CashierReportPage() {
  const today = new Date().toISOString().split('T')[0]
  const weekAgo = new Date(Date.now() - 7 * 86400000).toISOString().split('T')[0]

  const [fromDate, setFromDate] = useState(weekAgo)
  const [toDate, setToDate] = useState(today)
  const [branchId, setBranchId] = useState('')

  const { data: branches } = useBranches()
  const { data, isLoading } = useCashierReport({
    from_date: fromDate,
    to_date: toDate,
    branch_id: branchId || undefined,
  })

  return (
    <div className="p-6 lg:p-8">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
        <div>
          <h1 className="text-2xl font-bold text-text-primary">Sales by Cashier</h1>
          <p className="text-text-muted mt-1">Employee performance and sales breakdown</p>
        </div>
        <div className="flex items-center gap-2 flex-wrap">
          <input type="date" value={fromDate} onChange={(e) => setFromDate(e.target.value)}
            className="px-3 py-1.5 rounded-lg border border-border text-sm bg-white" />
          <span className="text-text-muted text-sm">to</span>
          <input type="date" value={toDate} onChange={(e) => setToDate(e.target.value)}
            className="px-3 py-1.5 rounded-lg border border-border text-sm bg-white" />
          {branches && branches.length > 1 && (
            <select value={branchId} onChange={(e) => setBranchId(e.target.value)}
              className="px-3 py-1.5 rounded-lg border border-border text-sm bg-white">
              <option value="">All Branches</option>
              {branches.map((b) => <option key={b.id} value={b.id}>{b.name}</option>)}
            </select>
          )}
          <ExportButton
            onClick={() => exportToExcel(
              (data || []).map((c: any) => ({
                Cashier: c.name,
                Branch: c.branch_name,
                Transactions: c.transaction_count,
                'Total Sales': c.total_sales,
              })),
              `cashier-report-${fromDate}-to-${toDate}`,
              'By Cashier',
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
                <th className="text-left px-5 py-3 font-semibold text-text-muted">Employee</th>
                <th className="text-left px-5 py-3 font-semibold text-text-muted">Branch</th>
                <th className="text-right px-5 py-3 font-semibold text-text-muted">Transactions</th>
                <th className="text-right px-5 py-3 font-semibold text-text-muted">Total Sales</th>
              </tr>
            </thead>
            <tbody>
              {isLoading ? (
                <tr><td colSpan={4} className="px-5 py-10 text-center text-text-muted">Loading...</td></tr>
              ) : !data || data.length === 0 ? (
                <tr><td colSpan={4} className="px-5 py-10 text-center text-text-muted">No sales data for this period</td></tr>
              ) : (
                data.map((cashier: any) => (
                  <tr key={cashier.user_id} className="border-b border-border last:border-0 hover:bg-background/30">
                    <td className="px-5 py-3">
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 rounded-full bg-grit-blue flex items-center justify-center text-white text-sm font-bold">
                          {cashier.name?.charAt(0)}
                        </div>
                        <span className="font-medium text-text-primary">{cashier.name}</span>
                      </div>
                    </td>
                    <td className="px-5 py-3 text-text-muted">{cashier.branch_name}</td>
                    <td className="px-5 py-3 text-right font-semibold text-text-primary">{cashier.transaction_count}</td>
                    <td className="px-5 py-3 text-right font-semibold text-grit-blue">{formatUGX(cashier.total_sales)}</td>
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
