import { createFileRoute, Link } from '@tanstack/react-router'
import { BarChart3, Boxes, TrendingUp, Users } from 'lucide-react'
import { useAuthStore } from '@/stores/auth.store'

export const Route = createFileRoute('/_app/reports/')({
  component: ReportsPage,
})

const reports = [
  { title: 'Daily Sales', description: 'Sales breakdown by date range with charts', path: '/reports/daily', icon: <BarChart3 size={24} />, color: 'bg-grit-blue', roles: ['admin', 'manager'] },
  { title: 'Stock Levels', description: 'Inventory levels across all branches', path: '/reports/stock', icon: <Boxes size={24} />, color: 'bg-role-manager', roles: ['admin', 'manager'] },
  { title: 'Profit & Loss', description: 'Capital, revenue, COGS, and margins', path: '/reports/pnl', icon: <TrendingUp size={24} />, color: 'bg-role-admin', roles: ['admin'] },
  { title: 'By Cashier', description: 'Sales performance per employee', path: '/reports/cashiers', icon: <Users size={24} />, color: 'bg-warning', roles: ['admin'] },
]

function ReportsPage() {
  const role = useAuthStore((s) => s.currentRole)

  return (
    <div className="p-6 lg:p-8">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-text-primary">Reports</h1>
        <p className="text-text-muted mt-1">Business analytics and insights</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {reports
          .filter((r) => role && r.roles.includes(role))
          .map((report) => (
            <Link key={report.path} to={report.path}
              className="bg-surface rounded-xl border border-border p-5 shadow-sm hover:shadow-md hover:border-grit-blue/30 transition group">
              <div className="flex items-start gap-4">
                <div className={`w-12 h-12 rounded-xl ${report.color} text-white flex items-center justify-center group-hover:scale-105 transition`}>
                  {report.icon}
                </div>
                <div>
                  <h3 className="text-base font-semibold text-text-primary group-hover:text-grit-blue transition">{report.title}</h3>
                  <p className="text-sm text-text-muted mt-1">{report.description}</p>
                </div>
              </div>
            </Link>
          ))}
      </div>
    </div>
  )
}
