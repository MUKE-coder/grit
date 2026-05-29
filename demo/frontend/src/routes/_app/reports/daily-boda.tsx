import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useDailyBodaReport } from '@/hooks/useReports'
import { ExportButton } from '@/components/export-button'
import { exportSheetsToExcel } from '@/lib/export'
import { formatUGX } from '@/lib/utils'
import { CalendarClock, Wallet, AlertCircle } from 'lucide-react'

export const Route = createFileRoute('/_app/reports/daily-boda')({
  component: DailyBodaReportPage,
})

function defaultRange() {
  const now = new Date()
  const from = new Date(now); from.setDate(from.getDate() - 30)
  return { from: from.toISOString().slice(0, 10), to: now.toISOString().slice(0, 10) }
}

function DailyBodaReportPage() {
  const init = defaultRange()
  const [from, setFrom] = useState(init.from)
  const [to, setTo] = useState(init.to)
  const { data, isLoading } = useDailyBodaReport({ from_date: from, to_date: to })

  if (isLoading) return <div className="p-6 text-text-muted">Loading...</div>

  const topDrivers: Array<{ driver_id: number; full_name: string; phone: string; total_paid: number; payment_count: number }> = data?.top_drivers || []
  const daily: Array<{ date: string; total: number }> = data?.daily || []
  const fleet: Array<{ status: string; count: number }> = data?.fleet || []

  const handleExport = () => {
    exportSheetsToExcel([
      { name: 'Summary', rows: [
        { Metric: 'Total Collected', Value: data?.total_collected || 0 },
        { Metric: 'Transactions', Value: data?.transaction_count || 0 },
        { Metric: 'Unpaid Balance', Value: data?.unpaid_balance || 0 },
        { Metric: 'From', Value: from }, { Metric: 'To', Value: to },
      ]},
      { name: 'Top Drivers', rows: topDrivers.map((d) => ({ Driver: d.full_name, Phone: d.phone, Days: d.payment_count, 'Total Paid': d.total_paid })) },
      { name: 'Daily', rows: daily.map((d) => ({ Date: d.date, Total: d.total })) },
      { name: 'Fleet', rows: fleet.map((f) => ({ Status: f.status, Count: f.count })) },
    ], `daily-boda-${from}-to-${to}`)
  }

  return (
    <div className="p-4 lg:p-8">
      <div className="mb-6 flex flex-col sm:flex-row sm:items-start sm:justify-between gap-3">
        <div>
          <h1 className="text-[20px] lg:text-2xl font-bold text-text-primary">Daily Boda</h1>
          <p className="text-text-muted mt-1 text-[13px]">Rental fleet collections in the selected period.</p>
        </div>
        <div className="flex gap-2 flex-wrap">
          <input type="date" value={from} onChange={(e) => setFrom(e.target.value)} className="px-3 h-8 rounded-lg border border-border bg-white text-[12.5px]" />
          <input type="date" value={to} onChange={(e) => setTo(e.target.value)} className="px-3 h-8 rounded-lg border border-border bg-white text-[12.5px]" />
          <ExportButton onClick={handleExport} />
        </div>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
        <BigStat icon={<CalendarClock size={20} />} label="Total Collected" value={formatUGX(data?.total_collected || 0)} color="bg-grit-blue" />
        <BigStat icon={<Wallet size={20} />} label="Transactions" value={String(data?.transaction_count || 0)} color="bg-role-manager" />
        <BigStat icon={<AlertCircle size={20} />} label="Unpaid Balance" value={formatUGX(data?.unpaid_balance || 0)} color={data?.unpaid_balance > 0 ? 'bg-warning' : 'bg-text-muted'} />
        <BigStat icon={<CalendarClock size={20} />} label="Top Drivers Listed" value={String(topDrivers.length)} color="bg-role-admin" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        <Card title="Top Drivers">
          {topDrivers.length === 0 ? <p className="text-text-muted text-sm">No payments in range</p> : (
            <table className="w-full text-sm">
              <thead className="text-xs text-text-muted uppercase tracking-wider">
                <tr><th className="text-left pb-2">Driver</th><th className="text-right pb-2">Days Paid</th><th className="text-right pb-2">Total</th></tr>
              </thead>
              <tbody>
                {topDrivers.map((d) => (
                  <tr key={d.driver_id} className="border-t border-border">
                    <td className="py-2 font-medium">{d.full_name}<p className="text-xs text-text-muted font-mono">{d.phone}</p></td>
                    <td className="py-2 text-right">{d.payment_count}</td>
                    <td className="py-2 text-right font-bold">{formatUGX(d.total_paid)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </Card>

        <Card title="Fleet Status">
          {fleet.length === 0 ? <p className="text-text-muted text-sm">No motorcycles in fleet</p> : (
            <table className="w-full text-sm">
              <thead className="text-xs text-text-muted uppercase tracking-wider">
                <tr><th className="text-left pb-2">Status</th><th className="text-right pb-2">Count</th></tr>
              </thead>
              <tbody>
                {fleet.map((f) => (
                  <tr key={f.status} className="border-t border-border">
                    <td className="py-2 capitalize">{f.status.replace('_', ' ')}</td>
                    <td className="py-2 text-right font-bold">{f.count}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </Card>
      </div>

      <Card title="Daily Collections">
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
      <p className="text-xl font-bold text-text-primary">{value}</p>
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
