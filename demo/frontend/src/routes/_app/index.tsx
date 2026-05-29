import { createFileRoute, Link } from '@tanstack/react-router'
import { useAuthStore } from '@/stores/auth.store'
import { useMainReport } from '@/hooks/useReports'
import { useDashboard } from '@/hooks/useReports'
import { formatUGX } from '@/lib/utils'
import {
  ArrowRight, Bike, BarChart3, CircleDollarSign, ShoppingCart, HandCoins,
  AlertTriangle, Wallet, CalendarClock, Boxes, TrendingUp, Users,
} from 'lucide-react'

export const Route = createFileRoute('/_app/')({
  component: HomePage,
})

/**
 * Home dashboard — shows the dashboard for whichever workspace is active.
 * The user picks their workspace at the top of the sidebar; this component
 * reads `currentWorkspace` and dispatches to the matching dashboard.
 */
function HomePage() {
  const workspace = useAuthStore((s) => s.currentWorkspace)
  return workspace === 'loans' ? <LoansDashboard /> : <SparesDashboard />
}

/* ─────────────── Loans & Motorcycles dashboard ─────────────── */

function LoansDashboard() {
  const user = useAuthStore((s) => s.user)
  const { data, isLoading } = useMainReport()

  const loans = data?.loans
  const motos = data?.motorcycles
  const repayments = data?.repayments
  const dailyBoda = data?.daily_boda

  return (
    <div className="h-full overflow-y-auto">
      <div className="p-6 lg:p-8 max-w-7xl mx-auto">
        <Header
          eyebrow="Loans & Motorcycles"
          title={`Welcome back, ${user?.name?.split(' ')[0] || 'there'}.`}
          subtitle="Active loan portfolio, collections, and motorcycle inventory at a glance."
          rightStat={isLoading ? '—' : formatUGX(repayments?.collected_today || 0)}
          rightLabel="Collected today"
        />

        {/* Stat row */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-3 mb-6">
          <StatCard icon={<HandCoins size={16} />} label="Active loans" value={String(loans?.active || 0)} accent="bg-accent-light text-accent" />
          <StatCard icon={<CircleDollarSign size={16} />} label="Outstanding" value={formatUGX(loans?.total_outstanding || 0)} accent="bg-accent-light text-accent" />
          <StatCard icon={<Bike size={16} />} label="Bikes available" value={String(motos?.available || 0)} accent="bg-warning-light text-warning-dark" />
          <StatCard icon={<Wallet size={16} />} label="Bikes on loan" value={String(motos?.on_loan || 0)} accent="bg-success-light text-success-dark" />
        </div>

        {/* Action queues */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 mb-6">
          <QueueCard
            title="Pending loans"
            count={loans?.pending || 0}
            link="/loans?status=pending"
            empty="Nothing waiting on approval"
            icon={<HandCoins size={14} />}
          />
          <QueueCard
            title="Repayments to verify"
            count={repayments?.pending_verification || 0}
            link="/repayments?status=pending"
            empty="All repayments verified"
            icon={<AlertTriangle size={14} />}
            urgent={(repayments?.pending_verification || 0) > 0}
          />
          <QueueCard
            title="Overdue installments"
            count={repayments?.overdue_installments || 0}
            link="/loans?status=active"
            empty="No overdue installments"
            icon={<AlertTriangle size={14} />}
            urgent={(repayments?.overdue_installments || 0) > 0}
          />
        </div>

        {/* Loan portfolio breakdown */}
        <SectionTitle>Loan portfolio</SectionTitle>
        <div className="grid grid-cols-2 md:grid-cols-5 gap-3 mb-6">
          <PortfolioCell label="Pending" value={loans?.pending || 0} link="/loans?status=pending" />
          <PortfolioCell label="Approved" value={loans?.approved || 0} link="/loans?status=approved" />
          <PortfolioCell label="Active" value={loans?.active || 0} link="/loans?status=active" />
          <PortfolioCell label="Completed" value={loans?.completed || 0} link="/loans?status=completed" />
          <PortfolioCell label="Defaulted" value={loans?.defaulted || 0} link="/loans?status=defaulted" tone="danger" />
        </div>

        {/* Daily Boda strip */}
        {dailyBoda && (dailyBoda.today_count || 0) > 0 && (
          <>
            <SectionTitle>Daily Boda</SectionTitle>
            <Link
              to="/daily-boda/payments"
              className="block bg-surface rounded-xl border border-border p-4 hover:bg-surface-hover transition group"
            >
              <div className="flex items-center justify-between gap-4">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-lg bg-accent-light text-accent flex items-center justify-center">
                    <CalendarClock size={18} />
                  </div>
                  <div>
                    <p className="text-[13px] font-semibold text-foreground">{formatUGX(dailyBoda.today_total)}</p>
                    <p className="text-[12px] text-foreground-muted">{dailyBoda.today_count} collections today</p>
                  </div>
                </div>
                <ArrowRight size={14} className="text-foreground-muted group-hover:text-accent transition" />
              </div>
            </Link>
          </>
        )}
      </div>
    </div>
  )
}

/* ─────────────── Spares POS dashboard ─────────────── */

function SparesDashboard() {
  const user = useAuthStore((s) => s.user)
  const { data, isLoading } = useDashboard()

  return (
    <div className="h-full overflow-y-auto">
      <div className="p-6 lg:p-8 max-w-7xl mx-auto">
        <Header
          eyebrow="Spares POS"
          title={`Welcome back, ${user?.name?.split(' ')[0] || 'there'}.`}
          subtitle="Today's spare-parts inventory and POS performance."
          rightStat={isLoading ? '—' : formatUGX(data?.today_sales_total || 0)}
          rightLabel="Today's sales"
          rightAction={
            <Link to="/pos" className="inline-flex items-center gap-1.5 h-9 px-3.5 rounded-lg bg-accent text-white text-[13px] font-semibold hover:bg-accent-hover transition shadow-xs">
              Open POS <ArrowRight size={12} />
            </Link>
          }
        />

        <div className="grid grid-cols-2 md:grid-cols-4 gap-3 mb-6">
          <StatCard icon={<ShoppingCart size={16} />} label="Today's sales" value={formatUGX(data?.today_sales_total || 0)} accent="bg-accent-light text-accent" />
          <StatCard icon={<BarChart3 size={16} />} label="Transactions" value={String(data?.today_transaction_count || 0)} accent="bg-success-light text-success-dark" />
          <StatCard icon={<Boxes size={16} />} label="Stock value" value={formatUGX(data?.total_stock_value || 0)} accent="bg-warning-light text-warning-dark" />
          <StatCard icon={<TrendingUp size={16} />} label="Est. profit today" value={formatUGX(data?.estimated_profit_today || 0)} accent="bg-success-light text-success-dark" />
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <Card>
            <CardHeader title="Low stock" link="/stock" />
            {!data?.low_stock_items || data.low_stock_items.length === 0 ? (
              <p className="text-[13px] text-foreground-muted py-4 px-4">All stock levels are healthy.</p>
            ) : (
              <ul className="divide-y divide-border-subtle">
                {data.low_stock_items.slice(0, 6).map((item: any) => (
                  <li key={`${item.product_id}-${item.branch_id}`} className="px-4 py-2.5 flex items-center justify-between text-[13px]">
                    <div>
                      <p className="font-medium text-foreground">{item.product_title}</p>
                      <p className="text-[11.5px] text-foreground-muted">{item.branch_name}</p>
                    </div>
                    <span className={`px-2 py-0.5 rounded-full text-[11px] font-semibold ${
                      item.quantity === 0 ? 'bg-danger-light text-danger-dark' : 'bg-warning-light text-warning-dark'
                    }`}>
                      {item.quantity} left
                    </span>
                  </li>
                ))}
              </ul>
            )}
          </Card>

          <Card>
            <CardHeader title="Recent sales" link="/sales" />
            {!data?.recent_sales || data.recent_sales.length === 0 ? (
              <p className="text-[13px] text-foreground-muted py-4 px-4">No sales today yet.</p>
            ) : (
              <ul className="divide-y divide-border-subtle">
                {data.recent_sales.slice(0, 6).map((sale: any) => {
                  // Render up to 3 product names then a "+ N more" pill so
                  // the row stays compact even on big sales.
                  const items: { product_title: string; quantity: number }[] = sale.items || []
                  const shown = items.slice(0, 3)
                  const overflow = items.length - shown.length
                  const productLine = shown
                    .map((i) => `${i.quantity}× ${i.product_title || 'Product'}`)
                    .join(' · ')
                  return (
                    <li key={sale.id} className="px-4 py-2.5 text-[13px]">
                      <div className="flex items-start justify-between gap-3">
                        <div className="min-w-0 flex-1">
                          <p className="font-medium text-foreground truncate">
                            {productLine || 'No items'}
                            {overflow > 0 && (
                              <span className="ml-1 text-[11px] text-foreground-muted font-normal">
                                +{overflow} more
                              </span>
                            )}
                          </p>
                          <p className="text-[11.5px] text-foreground-muted mt-0.5">
                            {sale.cashier_name} · {sale.branch_name} · {new Date(sale.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                          </p>
                        </div>
                        <p className="font-semibold text-foreground tabular-nums shrink-0">
                          {formatUGX(sale.total)}
                        </p>
                      </div>
                    </li>
                  )
                })}
              </ul>
            )}
          </Card>
        </div>
      </div>
    </div>
  )
}

/* ─────────────── Shared dashboard primitives ─────────────── */

function Header({
  eyebrow, title, subtitle, rightStat, rightLabel, rightAction,
}: {
  eyebrow: string; title: string; subtitle?: string
  rightStat?: string; rightLabel?: string; rightAction?: React.ReactNode
}) {
  return (
    <div className="flex flex-col md:flex-row md:items-start md:justify-between gap-4 mb-6">
      <div>
        <p className="text-[10.5px] font-semibold uppercase tracking-wider text-foreground-muted">{eyebrow}</p>
        <h1 className="text-[22px] font-semibold text-foreground mt-1">{title}</h1>
        {subtitle && <p className="text-[13px] text-foreground-muted mt-0.5">{subtitle}</p>}
      </div>
      <div className="flex items-center gap-3">
        {rightStat && (
          <div className="text-right">
            {rightLabel && <p className="text-[10.5px] uppercase tracking-wider text-foreground-muted font-semibold">{rightLabel}</p>}
            <p className="text-[20px] font-semibold text-foreground tabular-nums leading-tight">{rightStat}</p>
          </div>
        )}
        {rightAction}
      </div>
    </div>
  )
}

function StatCard({
  icon, label, value, accent,
}: {
  icon: React.ReactNode; label: string; value: string; accent: string
}) {
  return (
    <div className="bg-surface rounded-xl border border-border p-4 shadow-xs">
      <div className="flex items-center gap-2.5 mb-2">
        <div className={`w-8 h-8 rounded-lg flex items-center justify-center ${accent}`}>{icon}</div>
        <p className="text-[11.5px] text-foreground-muted">{label}</p>
      </div>
      <p className="text-[18px] font-semibold text-foreground tabular-nums">{value}</p>
    </div>
  )
}

function SectionTitle({ children }: { children: React.ReactNode }) {
  return <h2 className="text-[11px] font-semibold uppercase tracking-wider text-foreground-muted mb-3">{children}</h2>
}

function QueueCard({
  title, count, link, empty, icon, urgent,
}: {
  title: string; count: number; link: string; empty: string; icon: React.ReactNode; urgent?: boolean
}) {
  return (
    <Link
      to={link}
      className="block bg-surface rounded-xl border border-border p-4 hover:bg-surface-hover transition group shadow-xs"
    >
      <div className="flex items-center justify-between mb-1.5">
        <div className="flex items-center gap-2 text-foreground-muted">
          <span className={urgent ? 'text-warning-dark' : 'text-foreground-muted'}>{icon}</span>
          <p className="text-[12.5px] font-medium">{title}</p>
        </div>
        <ArrowRight size={12} className="text-foreground-muted group-hover:text-accent transition" />
      </div>
      {count === 0 ? (
        <p className="text-[12.5px] text-foreground-muted">{empty}</p>
      ) : (
        <p className={`text-[22px] font-semibold tabular-nums leading-tight ${urgent ? 'text-warning-dark' : 'text-foreground'}`}>
          {count} <span className="text-[12.5px] font-normal text-foreground-muted">items</span>
        </p>
      )}
    </Link>
  )
}

function PortfolioCell({
  label, value, link, tone,
}: {
  label: string; value: number; link: string; tone?: 'danger'
}) {
  return (
    <Link
      to={link}
      className={`block bg-surface rounded-xl border p-4 hover:bg-surface-hover transition shadow-xs ${tone === 'danger' ? 'border-danger/30' : 'border-border'}`}
    >
      <p className="text-[11.5px] text-foreground-muted">{label}</p>
      <p className={`text-[22px] font-semibold tabular-nums leading-tight mt-1 ${tone === 'danger' ? 'text-danger-dark' : 'text-foreground'}`}>
        {value}
      </p>
    </Link>
  )
}

function Card({ children }: { children: React.ReactNode }) {
  return <div className="bg-surface rounded-xl border border-border shadow-xs overflow-hidden">{children}</div>
}

function CardHeader({ title, link }: { title: string; link?: string }) {
  return (
    <div className="px-4 py-3 border-b border-border-subtle flex items-center justify-between">
      <h3 className="text-[12.5px] font-semibold text-foreground">{title}</h3>
      {link && (
        <Link to={link} className="text-[11.5px] text-accent font-medium hover:underline inline-flex items-center gap-0.5">
          View all <ArrowRight size={11} />
        </Link>
      )}
    </div>
  )
}
