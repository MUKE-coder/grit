import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { TrendingUp, DollarSign, Boxes, PieChart } from 'lucide-react'
import { usePnLReport } from '@/hooks/useReports'
import { StatCard } from '@/components/reports/StatCard'
import { ExportButton } from '@/components/export-button'
import { exportToExcel } from '@/lib/export'
import { formatUGX } from '@/lib/utils'

export const Route = createFileRoute('/_app/reports/pnl')({
  component: PnLReportPage,
})

function PnLReportPage() {
  const monthAgo = new Date(Date.now() - 30 * 86400000).toISOString().split('T')[0]
  const today = new Date().toISOString().split('T')[0]

  const [fromDate, setFromDate] = useState(monthAgo)
  const [toDate, setToDate] = useState(today)

  const { data, isLoading } = usePnLReport({ from_date: fromDate, to_date: toDate })

  const margin = data?.gross_margin_percent || 0

  const handleExport = () => {
    if (!data) return
    exportToExcel(
      [
        { Metric: 'Capital Invested', Value: data.total_capital_invested || 0 },
        { Metric: 'Revenue', Value: data.total_revenue || 0 },
        { Metric: 'Cost of Goods Sold', Value: data.total_cogs || 0 },
        { Metric: 'Gross Profit', Value: data.gross_profit || 0 },
        { Metric: 'Gross Margin %', Value: margin.toFixed(2) },
        { Metric: 'From', Value: fromDate },
        { Metric: 'To', Value: toDate },
      ],
      `pnl-${fromDate}-to-${toDate}`,
      'Profit & Loss',
    )
  }

  return (
    <div className="p-4 lg:p-8">
      <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-3 mb-6">
        <div>
          <h1 className="text-[20px] lg:text-2xl font-bold text-text-primary">Profit & Loss</h1>
          <p className="text-text-muted mt-1 text-[13px]">Capital investment and profitability analysis</p>
        </div>
        <div className="flex items-center gap-2 flex-wrap">
          <input type="date" value={fromDate} onChange={(e) => setFromDate(e.target.value)}
            className="px-3 h-8 rounded-lg border border-border text-[12.5px] bg-white" />
          <span className="text-text-muted text-[12.5px]">to</span>
          <input type="date" value={toDate} onChange={(e) => setToDate(e.target.value)}
            className="px-3 h-8 rounded-lg border border-border text-[12.5px] bg-white" />
          <ExportButton onClick={handleExport} disabled={!data || isLoading} />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <StatCard icon={<Boxes size={20} />} label="Capital Invested" value={isLoading ? '...' : formatUGX(data?.total_capital_invested || 0)} color="bg-role-admin" subtitle="All-time stock-in cost" />
        <StatCard icon={<DollarSign size={20} />} label="Revenue" value={isLoading ? '...' : formatUGX(data?.total_revenue || 0)} color="bg-grit-blue" subtitle="Sales in period" />
        <StatCard icon={<PieChart size={20} />} label="Cost of Goods Sold" value={isLoading ? '...' : formatUGX(data?.total_cogs || 0)} color="bg-warning" subtitle="Product cost for sales" />
        <StatCard icon={<TrendingUp size={20} />} label="Gross Profit" value={isLoading ? '...' : formatUGX(data?.gross_profit || 0)} color={data?.gross_profit >= 0 ? 'bg-grit-blue' : 'bg-error'} subtitle="Revenue minus COGS" />
      </div>

      <div className="bg-surface rounded-xl border border-border shadow-sm p-6">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-sm font-semibold text-text-primary">Gross Margin</h3>
          <span className={`text-3xl font-bold ${margin >= 20 ? 'text-grit-blue' : margin >= 10 ? 'text-warning' : 'text-error'}`}>
            {isLoading ? '...' : `${margin.toFixed(1)}%`}
          </span>
        </div>
        <div className="w-full h-4 bg-background rounded-full overflow-hidden">
          <div
            className={`h-full rounded-full transition-all duration-500 ${margin >= 20 ? 'bg-grit-blue' : margin >= 10 ? 'bg-warning' : 'bg-error'}`}
            style={{ width: `${Math.min(margin, 100)}%` }}
          />
        </div>
        <div className="flex justify-between mt-2 text-xs text-text-muted">
          <span>0%</span>
          <span>Revenue: {formatUGX(data?.total_revenue || 0)} | COGS: {formatUGX(data?.total_cogs || 0)}</span>
          <span>100%</span>
        </div>
      </div>
    </div>
  )
}
