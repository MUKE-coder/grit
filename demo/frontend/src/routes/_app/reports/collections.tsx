import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useCollectionsReport } from '@/hooks/useReports'
import { ExportButton } from '@/components/export-button'
import { exportSheetsToExcel } from '@/lib/export'
import { formatUGX } from '@/lib/utils'
import { CircleDollarSign, Wallet, Smartphone } from 'lucide-react'

export const Route = createFileRoute('/_app/reports/collections')({
  component: CollectionsReportPage,
})

function defaultRange() {
  const now = new Date()
  const from = new Date(now); from.setDate(from.getDate() - 30)
  return { from: from.toISOString().slice(0, 10), to: now.toISOString().slice(0, 10) }
}

function CollectionsReportPage() {
  const init = defaultRange()
  const [from, setFrom] = useState(init.from)
  const [to, setTo] = useState(init.to)
  const { data, isLoading } = useCollectionsReport({ from_date: from, to_date: to })

  if (isLoading) return <div className="p-6 text-text-muted">Loading...</div>

  const byMethod: Array<{ payment_method: string; total: number; count: number }> = data?.by_method || []
  const byCollector: Array<{ user_id: number; name: string; total: number; count: number }> = data?.by_collector || []
  const daily: Array<{ date: string; total: number }> = data?.daily || []

  const handleExport = () => {
    exportSheetsToExcel([
      { name: 'Summary', rows: [
        { Metric: 'Total Collected', Value: data?.total_collected || 0 },
        { Metric: 'Transactions', Value: data?.transaction_count || 0 },
        { Metric: 'From', Value: from }, { Metric: 'To', Value: to },
      ]},
      { name: 'By Method', rows: byMethod.map((m) => ({ Method: m.payment_method, Count: m.count, Total: m.total })) },
      { name: 'By Collector', rows: byCollector.map((c) => ({ Collector: c.name, Count: c.count, Total: c.total })) },
      { name: 'Daily', rows: daily.map((d) => ({ Date: d.date, Total: d.total })) },
    ], `collections-${from}-to-${to}`)
  }

  return (
    <div className="p-4 lg:p-8">
      <div className="mb-6 flex flex-col sm:flex-row sm:items-start sm:justify-between gap-3">
        <div>
          <h1 className="text-[20px] lg:text-2xl font-bold text-text-primary">Loan Collections</h1>
          <p className="text-text-muted mt-1 text-[13px]">Repayments collected (approved) in the selected period.</p>
        </div>
        <div className="flex gap-2 flex-wrap">
          <input type="date" value={from} onChange={(e) => setFrom(e.target.value)} className="px-3 h-8 rounded-lg border border-border bg-white text-[12.5px]" />
          <input type="date" value={to} onChange={(e) => setTo(e.target.value)} className="px-3 h-8 rounded-lg border border-border bg-white text-[12.5px]" />
          <ExportButton onClick={handleExport} />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <BigStat icon={<CircleDollarSign size={20} />} label="Total Collected" value={formatUGX(data?.total_collected || 0)} color="bg-grit-blue" />
        <BigStat icon={<Wallet size={20} />} label="Transactions" value={String(data?.transaction_count || 0)} color="bg-role-manager" />
        <BigStat icon={<Smartphone size={20} />} label="Avg / Transaction" value={formatUGX(data?.transaction_count ? (data.total_collected / data.transaction_count) : 0)} color="bg-role-admin" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        <Card title="By Payment Method">
          {byMethod.length === 0 ? <p className="text-text-muted text-sm">No collections in range</p> : (
            <table className="w-full text-sm">
              <thead className="text-xs text-text-muted uppercase tracking-wider">
                <tr><th className="text-left pb-2">Method</th><th className="text-right pb-2">Count</th><th className="text-right pb-2">Total</th></tr>
              </thead>
              <tbody>
                {byMethod.map((m) => (
                  <tr key={m.payment_method} className="border-t border-border">
                    <td className="py-2 capitalize">{m.payment_method.replace('_', ' ')}</td>
                    <td className="py-2 text-right">{m.count}</td>
                    <td className="py-2 text-right font-medium">{formatUGX(m.total)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </Card>

        <Card title="By Collector">
          {byCollector.length === 0 ? <p className="text-text-muted text-sm">No collections in range</p> : (
            <table className="w-full text-sm">
              <thead className="text-xs text-text-muted uppercase tracking-wider">
                <tr><th className="text-left pb-2">User</th><th className="text-right pb-2">Count</th><th className="text-right pb-2">Total</th></tr>
              </thead>
              <tbody>
                {byCollector.map((c) => (
                  <tr key={c.user_id} className="border-t border-border">
                    <td className="py-2 font-medium">{c.name}</td>
                    <td className="py-2 text-right">{c.count}</td>
                    <td className="py-2 text-right font-medium">{formatUGX(c.total)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </Card>
      </div>

      <Card title="Daily">
        {daily.length === 0 ? <p className="text-text-muted text-sm">No data</p> : (
          <table className="w-full text-sm">
            <thead className="text-xs text-text-muted uppercase tracking-wider">
              <tr><th className="text-left pb-2">Date</th><th className="text-right pb-2">Collected</th></tr>
            </thead>
            <tbody>
              {daily.map((d) => (
                <tr key={d.date} className="border-t border-border">
                  <td className="py-2">{new Date(d.date).toLocaleDateString()}</td>
                  <td className="py-2 text-right font-medium">{formatUGX(d.total)}</td>
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
      <p className="text-2xl font-bold text-text-primary">{value}</p>
    </div>
  )
}

function Card({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="bg-surface rounded-xl border border-border p-5 shadow-sm">
      <h3 className="text-sm font-semibold text-text-primary mb-4">{title}</h3>
      {children}
    </div>
  )
}
