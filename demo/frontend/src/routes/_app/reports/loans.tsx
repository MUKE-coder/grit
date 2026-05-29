import { createFileRoute, Link } from '@tanstack/react-router'
import { useLoansReport } from '@/hooks/useReports'
import { ExportButton } from '@/components/export-button'
import { exportSheetsToExcel } from '@/lib/export'
import { formatUGX } from '@/lib/utils'
import { HandCoins, Building2, Users, AlertTriangle } from 'lucide-react'

export const Route = createFileRoute('/_app/reports/loans')({
  component: LoansPortfolioReportPage,
})

function LoansPortfolioReportPage() {
  const { data, isLoading } = useLoansReport()

  if (isLoading) return <div className="p-6 text-text-muted">Loading...</div>

  const byStatus: Array<{ status: string; count: number; value: number }> = data?.by_status || []
  const byBranch: Array<{ branch_id: number; branch_name: string; active_loans: number; total_outstanding: number }> = data?.by_branch || []
  const topBorrowers: Array<{ borrower_id: number; full_name: string; phone: string; outstanding: number; loan_count: number }> = data?.top_borrowers || []
  const overdue: Array<{ loan_id: number; loan_number: string; borrower_name: string; days_past_due: number; amount_due: number }> = data?.overdue || []

  const totalActive = byStatus.find((s) => s.status === 'active')?.count || 0
  const totalDefaulted = byStatus.find((s) => s.status === 'defaulted')?.count || 0
  const portfolioOutstanding = byBranch.reduce((sum, b) => sum + b.total_outstanding, 0)

  const handleExport = () => {
    exportSheetsToExcel([
      { name: 'By Status', rows: byStatus.map((s) => ({ Status: s.status, Count: s.count, 'Principal Value': s.value })) },
      { name: 'By Branch', rows: byBranch.map((b) => ({ Branch: b.branch_name, 'Active Loans': b.active_loans, Outstanding: b.total_outstanding })) },
      { name: 'Top Borrowers', rows: topBorrowers.map((b) => ({ Name: b.full_name, Phone: b.phone, Loans: b.loan_count, Outstanding: b.outstanding })) },
      { name: 'Overdue', rows: overdue.map((o) => ({ 'Loan #': o.loan_number, Borrower: o.borrower_name, 'Days Late': o.days_past_due, 'Amount Due': o.amount_due })) },
    ], 'loan-portfolio')
  }

  return (
    <div className="p-4 lg:p-8">
      <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-3 mb-6">
        <div>
          <h1 className="text-[20px] lg:text-2xl font-bold text-text-primary">Loans Portfolio</h1>
          <p className="text-text-muted mt-1 text-[13px]">Active book, branch concentration, and overdue exposure.</p>
        </div>
        <ExportButton onClick={handleExport} />
      </div>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
        <Stat icon={<HandCoins size={18} />} label="Active loans" value={String(totalActive)} color="bg-grit-blue" />
        <Stat icon={<Building2 size={18} />} label="Total outstanding" value={formatUGX(portfolioOutstanding)} color="bg-role-manager" />
        <Stat icon={<Users size={18} />} label="Top borrowers shown" value={String(topBorrowers.length)} color="bg-role-admin" />
        <Stat icon={<AlertTriangle size={18} />} label="Defaulted" value={String(totalDefaulted)} color={totalDefaulted > 0 ? 'bg-error' : 'bg-text-muted'} />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        <Card title="By Status">
          <table className="w-full text-sm">
            <thead className="text-xs text-text-muted uppercase tracking-wider">
              <tr><th className="text-left pb-2">Status</th><th className="text-right pb-2">Count</th><th className="text-right pb-2">Principal Value</th></tr>
            </thead>
            <tbody>
              {byStatus.map((s) => (
                <tr key={s.status} className="border-t border-border">
                  <td className="py-2 capitalize">{s.status}</td>
                  <td className="py-2 text-right font-medium">{s.count}</td>
                  <td className="py-2 text-right">{formatUGX(s.value)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </Card>

        <Card title="By Branch (Active)">
          <table className="w-full text-sm">
            <thead className="text-xs text-text-muted uppercase tracking-wider">
              <tr><th className="text-left pb-2">Branch</th><th className="text-right pb-2">Loans</th><th className="text-right pb-2">Outstanding</th></tr>
            </thead>
            <tbody>
              {byBranch.length === 0 ? (
                <tr><td colSpan={3} className="py-4 text-center text-text-muted">No active loans</td></tr>
              ) : byBranch.map((b) => (
                <tr key={b.branch_id} className="border-t border-border">
                  <td className="py-2 font-medium">{b.branch_name}</td>
                  <td className="py-2 text-right">{b.active_loans}</td>
                  <td className="py-2 text-right font-medium">{formatUGX(b.total_outstanding)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </Card>
      </div>

      <Card title="Top Borrowers by Outstanding Balance" className="mb-6">
        {topBorrowers.length === 0 ? (
          <p className="text-text-muted text-sm">No active borrowers</p>
        ) : (
          <table className="w-full text-sm">
            <thead className="text-xs text-text-muted uppercase tracking-wider">
              <tr><th className="text-left pb-2">Borrower</th><th className="text-left pb-2">Phone</th><th className="text-right pb-2">Loans</th><th className="text-right pb-2">Outstanding</th></tr>
            </thead>
            <tbody>
              {topBorrowers.map((b) => (
                <tr key={b.borrower_id} className="border-t border-border">
                  <td className="py-2"><Link to={`/borrowers/${b.borrower_id}`} className="text-grit-blue font-medium hover:underline">{b.full_name}</Link></td>
                  <td className="py-2 font-mono text-text-muted">{b.phone}</td>
                  <td className="py-2 text-right">{b.loan_count}</td>
                  <td className="py-2 text-right font-bold">{formatUGX(b.outstanding)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </Card>

      <Card title="Overdue Installments" accent={overdue.length > 0}>
        {overdue.length === 0 ? (
          <p className="text-text-muted text-sm">No overdue installments — keep collecting.</p>
        ) : (
          <table className="w-full text-sm">
            <thead className="text-xs text-text-muted uppercase tracking-wider">
              <tr><th className="text-left pb-2">Loan</th><th className="text-left pb-2">Borrower</th><th className="text-right pb-2">Days Late</th><th className="text-right pb-2">Amount Due</th></tr>
            </thead>
            <tbody>
              {overdue.map((o) => (
                <tr key={o.loan_id} className="border-t border-border">
                  <td className="py-2"><Link to={`/loans/${o.loan_id}`} className="text-grit-blue font-mono hover:underline">{o.loan_number}</Link></td>
                  <td className="py-2 font-medium">{o.borrower_name}</td>
                  <td className="py-2 text-right"><span className="px-2 py-0.5 rounded-full text-xs font-semibold bg-error/10 text-error">{o.days_past_due}d</span></td>
                  <td className="py-2 text-right font-bold">{formatUGX(o.amount_due)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </Card>
    </div>
  )
}

function Stat({ icon, label, value, color }: { icon: React.ReactNode; label: string; value: string; color: string }) {
  return (
    <div className="bg-surface rounded-xl border border-border p-4 shadow-sm">
      <div className="flex items-center gap-3 mb-2">
        <div className={`w-9 h-9 rounded-lg ${color} text-white flex items-center justify-center`}>{icon}</div>
        <p className="text-xs text-text-muted">{label}</p>
      </div>
      <p className="text-xl font-bold text-text-primary">{value}</p>
    </div>
  )
}

function Card({ title, children, className, accent }: { title: string; children: React.ReactNode; className?: string; accent?: boolean }) {
  return (
    <div className={`bg-surface rounded-xl border ${accent ? 'border-error/30' : 'border-border'} p-5 shadow-sm ${className || ''}`}>
      <h3 className="text-sm font-semibold text-text-primary mb-4">{title}</h3>
      {children}
    </div>
  )
}
