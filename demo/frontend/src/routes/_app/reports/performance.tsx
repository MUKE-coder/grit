import { createFileRoute } from '@tanstack/react-router'
import { useState, useMemo } from 'react'
import {
  TrendingUp, ShoppingCart, HandCoins, Wallet, CalendarClock, Users, Bike, Award, BarChart3,
} from 'lucide-react'
import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend,
} from 'recharts'

import {
  useDailyReport, useCollectionsReport, useMotorcyclesReport, useDailyBodaReport, useLoansReport, useCashierReport,
} from '@/hooks/useReports'
import { ExportButton } from '@/components/export-button'
import { exportSheetsToExcel } from '@/lib/export'
import { formatUGX } from '@/lib/utils'

/**
 * Business Performance — admin's at-a-glance view of how the whole shop is
 * doing across all four segments (Spares, Motorcycle Cash Sales, Loan
 * Collections, Daily Boda) for any date range.
 *
 * Pulls four existing report endpoints in parallel and stitches them
 * together client-side. No new backend endpoint needed; we already have
 * each segment's report and this page is the unified roll-up.
 */
export const Route = createFileRoute('/_app/reports/performance')({
  component: PerformancePage,
})

function defaultRange() {
  const now = new Date()
  const from = new Date(now); from.setDate(from.getDate() - 30)
  return { from: from.toISOString().slice(0, 10), to: now.toISOString().slice(0, 10) }
}

function PerformancePage() {
  const init = defaultRange()
  const [from, setFrom] = useState(init.from)
  const [to, setTo] = useState(init.to)

  // Parallel queries — all four segment reports + portfolio + cashier ranking.
  const spares = useDailyReport({ from_date: from, to_date: to })
  const collections = useCollectionsReport({ from_date: from, to_date: to })
  const motorcycles = useMotorcyclesReport({ from_date: from, to_date: to })
  const dailyBoda = useDailyBodaReport({ from_date: from, to_date: to })
  const loanPortfolio = useLoansReport()
  const cashiers = useCashierReport({ from_date: from, to_date: to })

  const isLoading = spares.isLoading || collections.isLoading || motorcycles.isLoading || dailyBoda.isLoading

  // Per-segment totals — these are the "money in" lines.
  const sparesTotal = spares.data?.total_sales || 0
  const cashSalesTotal = motorcycles.data?.cash_sales_total || 0
  const collectionsTotal = collections.data?.total_collected || 0
  const dailyBodaTotal = dailyBoda.data?.total_collected || 0
  const grandTotal = sparesTotal + cashSalesTotal + collectionsTotal + dailyBodaTotal

  const segments = [
    { key: 'spares', label: 'Spares', value: sparesTotal, icon: <ShoppingCart size={16} />, accent: 'bg-success-light text-success-dark' },
    { key: 'cash', label: 'Motorcycle Cash', value: cashSalesTotal, icon: <Wallet size={16} />, accent: 'bg-warning-light text-warning-dark' },
    { key: 'loans', label: 'Loan Repayments', value: collectionsTotal, icon: <HandCoins size={16} />, accent: 'bg-accent-light text-accent' },
    { key: 'boda', label: 'Daily Boda', value: dailyBodaTotal, icon: <CalendarClock size={16} />, accent: 'bg-info-light text-info-dark' },
  ]

  // Combined daily series — merge the per-day arrays from each segment so a
  // single chart shows revenue per day broken down by segment.
  const daily = useMemo(() => {
    const map = new Map<string, { date: string; spares: number; cash: number; loans: number; boda: number; total: number }>()
    const ensure = (d: string) => {
      if (!map.has(d)) map.set(d, { date: d, spares: 0, cash: 0, loans: 0, boda: 0, total: 0 })
      return map.get(d)!
    }
    for (const r of spares.data?.sales_per_day || []) {
      const row = ensure(r.date.slice(0, 10))
      row.spares = r.total
      row.total += r.total
    }
    for (const r of collections.data?.daily || []) {
      const row = ensure(r.date.slice(0, 10))
      row.loans = r.total
      row.total += r.total
    }
    for (const r of dailyBoda.data?.daily || []) {
      const row = ensure(r.date.slice(0, 10))
      row.boda = r.total
      row.total += r.total
    }
    // Cash motorcycle sales report doesn't expose a per-day series — show
    // it as a flat band by dividing the total across the date range.
    return Array.from(map.values()).sort((a, b) => a.date.localeCompare(b.date))
  }, [spares.data, collections.data, dailyBoda.data])

  const topProducts = (spares.data?.top_products || []).slice(0, 8)
  const topBorrowers = (loanPortfolio.data?.top_borrowers || []).slice(0, 8)
  const topCashiers = (cashiers.data || []).slice(0, 8)
  const topDrivers = (dailyBoda.data?.top_drivers || []).slice(0, 8)

  const handleExport = () => {
    exportSheetsToExcel([
      {
        name: 'Performance Summary',
        rows: [
          { Metric: 'Period From', Value: from },
          { Metric: 'Period To', Value: to },
          { Metric: 'Total Revenue', Value: grandTotal },
          { Metric: 'Spares Revenue', Value: sparesTotal },
          { Metric: 'Motorcycle Cash Sales', Value: cashSalesTotal },
          { Metric: 'Loan Collections', Value: collectionsTotal },
          { Metric: 'Daily Boda Collections', Value: dailyBodaTotal },
          { Metric: 'Spares Transactions', Value: spares.data?.transaction_count || 0 },
          { Metric: 'Cash Sales Count', Value: motorcycles.data?.cash_sales_count || 0 },
          { Metric: 'Repayment Transactions', Value: collections.data?.transaction_count || 0 },
          { Metric: 'Daily Boda Transactions', Value: dailyBoda.data?.transaction_count || 0 },
          { Metric: 'Active Loans', Value: (loanPortfolio.data?.by_status || []).find((s: any) => s.status === 'active')?.count || 0 },
        ],
      },
      {
        name: 'Daily Revenue',
        rows: daily.map((d) => ({
          Date: d.date,
          Spares: d.spares,
          'Motorcycle Cash': d.cash,
          'Loan Repayments': d.loans,
          'Daily Boda': d.boda,
          Total: d.total,
        })),
      },
      {
        name: 'Top Products',
        rows: topProducts.map((p: any) => ({
          Product: p.product_title,
          'Qty Sold': p.quantity_sold,
          Revenue: p.revenue,
        })),
      },
      {
        name: 'Top Cashiers',
        rows: topCashiers.map((c: any) => ({
          Name: c.name,
          Branch: c.branch_name,
          Transactions: c.transaction_count,
          'Total Sales': c.total_sales,
        })),
      },
      {
        name: 'Top Borrowers',
        rows: topBorrowers.map((b: any) => ({
          Name: b.full_name,
          Phone: b.phone,
          Loans: b.loan_count,
          Outstanding: b.outstanding,
        })),
      },
      {
        name: 'Top Daily Boda Drivers',
        rows: topDrivers.map((d: any) => ({
          Driver: d.full_name,
          Phone: d.phone,
          Days: d.payment_count,
          'Total Paid': d.total_paid,
        })),
      },
    ], `business-performance-${from}-to-${to}`)
  }

  return (
    <div className="p-4 lg:p-8 max-w-7xl mx-auto">
      <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-3 mb-6">
        <div>
          <p className="text-[10.5px] font-semibold uppercase tracking-wider text-foreground-muted">Reports</p>
          <h1 className="text-[20px] lg:text-2xl font-bold text-foreground mt-1">Business Performance</h1>
          <p className="text-foreground-muted mt-1 text-[13px]">
            Combined revenue and activity across Spares, Motorcycles, Loans, and Daily Boda for the selected period.
          </p>
        </div>
        <div className="flex items-center gap-2 flex-wrap">
          <input type="date" value={from} onChange={(e) => setFrom(e.target.value)}
            className="px-3 h-8 rounded-lg border border-border bg-surface text-[12.5px]" />
          <span className="text-foreground-muted text-[12.5px]">to</span>
          <input type="date" value={to} onChange={(e) => setTo(e.target.value)}
            className="px-3 h-8 rounded-lg border border-border bg-surface text-[12.5px]" />
          <ExportButton onClick={handleExport} disabled={isLoading} />
        </div>
      </div>

      {/* Headline */}
      <div className="bg-surface border border-border rounded-2xl p-5 lg:p-6 shadow-xs mb-6">
        <div className="flex items-center justify-between gap-4">
          <div>
            <p className="text-[10.5px] uppercase tracking-wider text-foreground-muted font-semibold">Total revenue</p>
            <p className="text-[28px] lg:text-[34px] font-bold text-foreground tabular-nums leading-tight mt-1">
              {isLoading ? '—' : formatUGX(grandTotal)}
            </p>
            <p className="text-[12px] text-foreground-muted mt-1">{from} → {to}</p>
          </div>
          <div className="hidden sm:flex w-12 h-12 rounded-xl bg-accent text-white items-center justify-center">
            <TrendingUp size={22} />
          </div>
        </div>
      </div>

      {/* Per-segment cards */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3 mb-6">
        {segments.map((s) => {
          const pct = grandTotal > 0 ? (s.value / grandTotal) * 100 : 0
          return (
            <div key={s.key} className="bg-surface rounded-xl border border-border p-4 shadow-xs">
              <div className="flex items-center gap-2 mb-2">
                <div className={`w-8 h-8 rounded-lg flex items-center justify-center ${s.accent}`}>{s.icon}</div>
                <p className="text-[11.5px] text-foreground-muted">{s.label}</p>
              </div>
              <p className="text-[16px] lg:text-[18px] font-semibold text-foreground tabular-nums">{formatUGX(s.value)}</p>
              <div className="mt-2">
                <div className="h-1 rounded-full bg-surface-2 overflow-hidden">
                  <div className="h-full bg-accent rounded-full transition-all" style={{ width: `${Math.min(pct, 100)}%` }} />
                </div>
                <p className="text-[10.5px] text-foreground-muted mt-1 tabular-nums">{pct.toFixed(1)}% of total</p>
              </div>
            </div>
          )
        })}
      </div>

      {/* Daily revenue chart */}
      <div className="bg-surface rounded-xl border border-border shadow-xs p-4 lg:p-5 mb-6">
        <div className="flex items-center gap-2 mb-3">
          <BarChart3 size={14} className="text-foreground-muted" />
          <h3 className="text-[13px] font-semibold text-foreground">Daily revenue by segment</h3>
        </div>
        {daily.length === 0 ? (
          <p className="text-[13px] text-foreground-muted text-center py-12">No activity in this period.</p>
        ) : (
          <ResponsiveContainer width="100%" height={280}>
            <BarChart data={daily} margin={{ top: 5, right: 10, left: 0, bottom: 5 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="#E5E7EB" />
              <XAxis dataKey="date" tick={{ fontSize: 11 }} tickFormatter={(v) => new Date(v).toLocaleDateString('en', { month: 'short', day: 'numeric' })} />
              <YAxis tick={{ fontSize: 11 }} tickFormatter={(v) => `${(v / 1000).toFixed(0)}K`} />
              <Tooltip formatter={(v: number) => formatUGX(v)} labelFormatter={(v) => new Date(v).toLocaleDateString()} />
              <Legend wrapperStyle={{ fontSize: 12 }} />
              <Bar dataKey="spares" stackId="a" fill="#10B981" name="Spares" />
              <Bar dataKey="loans" stackId="a" fill="#1E7EF5" name="Loan Repayments" />
              <Bar dataKey="boda" stackId="a" fill="#9333EA" name="Daily Boda" />
            </BarChart>
          </ResponsiveContainer>
        )}
        <p className="text-[10.5px] text-foreground-muted mt-2">
          Note: motorcycle cash sales aren't shown per day (the report endpoint exposes only the period total).
        </p>
      </div>

      {/* Top performers grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <TopList icon={<Award size={14} />} title="Top products" rows={topProducts.map((p: any) => ({
          primary: p.product_title,
          secondary: `${p.quantity_sold} units`,
          value: formatUGX(p.revenue),
        }))} empty="No products sold in this period." />

        <TopList icon={<Users size={14} />} title="Top cashiers" rows={topCashiers.map((c: any) => ({
          primary: c.name,
          secondary: `${c.branch_name} · ${c.transaction_count} txns`,
          value: formatUGX(c.total_sales),
        }))} empty="No cashier activity." />

        <TopList icon={<HandCoins size={14} />} title="Top borrowers (outstanding)" rows={topBorrowers.map((b: any) => ({
          primary: b.full_name,
          secondary: `${b.phone} · ${b.loan_count} loan${b.loan_count === 1 ? '' : 's'}`,
          value: formatUGX(b.outstanding),
        }))} empty="No active borrowers." />

        <TopList icon={<Bike size={14} />} title="Top daily-boda drivers" rows={topDrivers.map((d: any) => ({
          primary: d.full_name,
          secondary: `${d.phone} · ${d.payment_count} day${d.payment_count === 1 ? '' : 's'}`,
          value: formatUGX(d.total_paid),
        }))} empty="No daily-boda activity." />
      </div>
    </div>
  )
}

function TopList({
  icon, title, rows, empty,
}: {
  icon: React.ReactNode
  title: string
  rows: Array<{ primary: string; secondary: string; value: string }>
  empty: string
}) {
  return (
    <div className="bg-surface rounded-xl border border-border shadow-xs">
      <div className="px-4 py-3 border-b border-border-subtle flex items-center gap-2">
        <span className="text-foreground-muted">{icon}</span>
        <h3 className="text-[12.5px] font-semibold text-foreground">{title}</h3>
      </div>
      {rows.length === 0 ? (
        <p className="px-4 py-6 text-[12.5px] text-foreground-muted text-center">{empty}</p>
      ) : (
        <ol className="divide-y divide-border-subtle">
          {rows.map((r, i) => (
            <li key={i} className="px-4 py-2.5 flex items-center justify-between gap-3 text-[13px]">
              <div className="flex items-center gap-3 min-w-0">
                <span className="w-6 h-6 rounded-full bg-surface-2 text-foreground-secondary text-[11.5px] font-bold flex items-center justify-center shrink-0">{i + 1}</span>
                <div className="min-w-0">
                  <p className="font-medium text-foreground truncate">{r.primary}</p>
                  <p className="text-[11.5px] text-foreground-muted truncate">{r.secondary}</p>
                </div>
              </div>
              <span className="font-semibold tabular-nums text-foreground">{r.value}</span>
            </li>
          ))}
        </ol>
      )}
    </div>
  )
}
