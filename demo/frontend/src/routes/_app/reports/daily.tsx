import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { BarChart3, ShoppingCart, DollarSign, CreditCard } from 'lucide-react'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'
import { useDailyReport } from '@/hooks/useReports'
import { useBranches } from '@/hooks/useBusiness'
import { StatCard } from '@/components/reports/StatCard'
import { ExportButton } from '@/components/export-button'
import { exportSheetsToExcel } from '@/lib/export'
import { formatUGX } from '@/lib/utils'

export const Route = createFileRoute('/_app/reports/daily')({
  component: DailyReportPage,
})

function DailyReportPage() {
  const today = new Date().toISOString().split('T')[0]
  const weekAgo = new Date(Date.now() - 7 * 86400000).toISOString().split('T')[0]

  const [fromDate, setFromDate] = useState(weekAgo)
  const [toDate, setToDate] = useState(today)
  const [branchId, setBranchId] = useState('')

  const { data: branches } = useBranches()
  const { data, isLoading } = useDailyReport({ from_date: fromDate, to_date: toDate, branch_id: branchId || undefined })

  const cashTotal = data?.by_payment_method?.find((p: any) => p.payment_method === 'cash')?.total || 0
  const mmTotal = data?.by_payment_method?.find((p: any) => p.payment_method === 'mobile_money')?.total || 0

  const handleExport = () => {
    if (!data) return
    exportSheetsToExcel(
      [
        {
          name: 'Summary',
          rows: [
            { Metric: 'Total Sales', Value: data.total_sales || 0 },
            { Metric: 'Transactions', Value: data.transaction_count || 0 },
            { Metric: 'Cash', Value: cashTotal },
            { Metric: 'Mobile Money', Value: mmTotal },
            { Metric: 'From', Value: fromDate },
            { Metric: 'To', Value: toDate },
          ],
        },
        {
          name: 'Per-Day',
          rows: (data.sales_per_day || []).map((d: any) => ({ Date: d.date, Total: d.total })),
        },
        {
          name: 'Top Products',
          rows: (data.top_products || []).map((p: any) => ({
            Product: p.product_title,
            'Qty Sold': p.quantity_sold,
            Revenue: p.revenue,
          })),
        },
        {
          name: 'Per-Product Breakdown',
          rows: (data.products_breakdown || []).map((p: any) => ({
            Product: p.product_title,
            'Cash Qty': p.cash_qty || 0,
            'Cash Revenue': p.cash_revenue || 0,
            'MoMo Qty': p.momo_qty || 0,
            'MoMo Revenue': p.momo_revenue || 0,
            'Total Qty': p.total_qty || 0,
            'Total Revenue': p.total_revenue || 0,
          })),
        },
      ],
      `daily-sales-${fromDate}-to-${toDate}`,
    )
  }

  return (
    <div className="p-4 lg:p-8">
      <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-3 mb-6">
        <div>
          <h1 className="text-[20px] lg:text-2xl font-bold text-text-primary">Daily Sales Report</h1>
          <p className="text-text-muted mt-1 text-[13px]">Sales breakdown by date range</p>
        </div>
        <div className="flex items-center gap-2 flex-wrap">
          <input type="date" value={fromDate} onChange={(e) => setFromDate(e.target.value)}
            className="px-3 h-8 rounded-lg border border-border text-[12.5px] bg-white" />
          <span className="text-text-muted text-[12.5px]">to</span>
          <input type="date" value={toDate} onChange={(e) => setToDate(e.target.value)}
            className="px-3 h-8 rounded-lg border border-border text-[12.5px] bg-white" />
          {branches && branches.length > 1 && (
            <select value={branchId} onChange={(e) => setBranchId(e.target.value)}
              className="px-3 h-8 rounded-lg border border-border text-[12.5px] bg-white">
              <option value="">All Branches</option>
              {branches.map((b) => <option key={b.id} value={b.id}>{b.name}</option>)}
            </select>
          )}
          <ExportButton onClick={handleExport} disabled={!data || isLoading} />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <StatCard icon={<ShoppingCart size={20} />} label="Total Sales" value={isLoading ? '...' : formatUGX(data?.total_sales || 0)} />
        <StatCard icon={<BarChart3 size={20} />} label="Transactions" value={isLoading ? '...' : String(data?.transaction_count || 0)} color="bg-role-manager" />
        <StatCard icon={<DollarSign size={20} />} label="Cash" value={isLoading ? '...' : formatUGX(cashTotal)} color="bg-grit-blue" />
        <StatCard icon={<CreditCard size={20} />} label="Mobile Money" value={isLoading ? '...' : formatUGX(mmTotal)} color="bg-role-admin" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-surface rounded-xl border border-border shadow-sm p-5">
          <h3 className="text-sm font-semibold text-text-primary mb-4">Sales Per Day</h3>
          {data?.sales_per_day && data.sales_per_day.length > 0 ? (
            <ResponsiveContainer width="100%" height={280}>
              <BarChart data={data.sales_per_day}>
                <CartesianGrid strokeDasharray="3 3" stroke="#E2E8F0" />
                <XAxis dataKey="date" tick={{ fontSize: 11 }} tickFormatter={(v) => new Date(v).toLocaleDateString('en', { month: 'short', day: 'numeric' })} />
                <YAxis tick={{ fontSize: 11 }} tickFormatter={(v) => `${(v / 1000).toFixed(0)}K`} />
                <Tooltip formatter={(v: number) => formatUGX(v)} labelFormatter={(v) => new Date(v).toLocaleDateString()} />
                <Bar dataKey="total" fill="#16A34A" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          ) : (
            <p className="text-sm text-text-muted py-10 text-center">No sales data for this period</p>
          )}
        </div>

        <div className="bg-surface rounded-xl border border-border shadow-sm">
          <div className="px-5 py-4 border-b border-border">
            <h3 className="text-sm font-semibold text-text-primary">Top Selling Products</h3>
          </div>
          <div className="p-5">
            {!data?.top_products || data.top_products.length === 0 ? (
              <p className="text-sm text-text-muted text-center py-6">No sales in this period</p>
            ) : (
              <div className="space-y-3">
                {data.top_products.map((p: any, i: number) => (
                  <div key={p.product_id} className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <span className="w-6 h-6 rounded-full bg-grit-blue-50 text-grit-blue text-xs font-bold flex items-center justify-center">{i + 1}</span>
                      <span className="text-sm font-medium text-text-primary">{p.product_title}</span>
                    </div>
                    <div className="text-right">
                      <p className="text-sm font-semibold text-text-primary">{formatUGX(p.revenue)}</p>
                      <p className="text-xs text-text-muted">{p.quantity_sold} sold</p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Per-product breakdown — every product sold in the period, with cash vs
          mobile money totals side by side. Renders below the existing Top
          Selling Products card so users still get the at-a-glance ranking. */}
      <div className="bg-surface rounded-xl border border-border shadow-sm mt-6">
        <div className="px-5 py-4 border-b border-border flex items-center justify-between">
          <div>
            <h3 className="text-sm font-semibold text-text-primary">Products sold</h3>
            <p className="text-xs text-text-muted mt-0.5">Per-product split between cash and mobile money</p>
          </div>
          <span className="text-xs text-text-muted">
            {data?.products_breakdown?.length || 0} product{(data?.products_breakdown?.length || 0) === 1 ? '' : 's'}
          </span>
        </div>
        {!data?.products_breakdown || data.products_breakdown.length === 0 ? (
          <p className="px-5 py-10 text-sm text-text-muted text-center">No sales in this period</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm tabular-nums">
              <thead className="bg-background border-b border-border text-[11px] uppercase tracking-wider text-text-muted">
                <tr>
                  <th className="text-left px-5 py-2.5 font-semibold">Product</th>
                  <th className="text-right px-3 py-2.5 font-semibold">Cash qty</th>
                  <th className="text-right px-3 py-2.5 font-semibold text-grit-blue">Cash revenue</th>
                  <th className="text-right px-3 py-2.5 font-semibold">MoMo qty</th>
                  <th className="text-right px-3 py-2.5 font-semibold text-role-admin">MoMo revenue</th>
                  <th className="text-right px-3 py-2.5 font-semibold">Total qty</th>
                  <th className="text-right px-5 py-2.5 font-semibold">Total revenue</th>
                </tr>
              </thead>
              <tbody>
                {data.products_breakdown.map((p: any) => (
                  <tr key={p.product_id} className="border-b border-border last:border-0 hover:bg-background/50">
                    <td className="px-5 py-2.5 font-medium text-text-primary">{p.product_title}</td>
                    <td className="px-3 py-2.5 text-right text-text-muted">{p.cash_qty || 0}</td>
                    <td className="px-3 py-2.5 text-right text-grit-blue font-medium">{formatUGX(p.cash_revenue || 0)}</td>
                    <td className="px-3 py-2.5 text-right text-text-muted">{p.momo_qty || 0}</td>
                    <td className="px-3 py-2.5 text-right text-role-admin font-medium">{formatUGX(p.momo_revenue || 0)}</td>
                    <td className="px-3 py-2.5 text-right font-medium">{p.total_qty || 0}</td>
                    <td className="px-5 py-2.5 text-right font-bold text-text-primary">{formatUGX(p.total_revenue || 0)}</td>
                  </tr>
                ))}
              </tbody>
              <tfoot className="bg-background border-t border-border font-semibold">
                <tr>
                  <td className="px-5 py-3">Totals</td>
                  <td className="px-3 py-3 text-right">
                    {data.products_breakdown.reduce((s: number, p: any) => s + (p.cash_qty || 0), 0)}
                  </td>
                  <td className="px-3 py-3 text-right text-grit-blue">
                    {formatUGX(data.products_breakdown.reduce((s: number, p: any) => s + (p.cash_revenue || 0), 0))}
                  </td>
                  <td className="px-3 py-3 text-right">
                    {data.products_breakdown.reduce((s: number, p: any) => s + (p.momo_qty || 0), 0)}
                  </td>
                  <td className="px-3 py-3 text-right text-role-admin">
                    {formatUGX(data.products_breakdown.reduce((s: number, p: any) => s + (p.momo_revenue || 0), 0))}
                  </td>
                  <td className="px-3 py-3 text-right">
                    {data.products_breakdown.reduce((s: number, p: any) => s + (p.total_qty || 0), 0)}
                  </td>
                  <td className="px-5 py-3 text-right">
                    {formatUGX(data.products_breakdown.reduce((s: number, p: any) => s + (p.total_revenue || 0), 0))}
                  </td>
                </tr>
              </tfoot>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}
