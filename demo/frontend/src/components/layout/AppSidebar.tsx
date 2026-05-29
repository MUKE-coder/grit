import { Link, useLocation } from '@tanstack/react-router'
import { useEffect, useRef, useState } from 'react'
import {
  LayoutDashboard, ShoppingCart, Package, Boxes, Receipt, BarChart3, TrendingUp,
  Users, Building2, Settings, ArrowLeftRight, ClipboardList, PackagePlus, UserPlus,
  PieChart, Bike, HandCoins, FileText, CalendarClock, Wallet, CircleDollarSign,
  Search, ChevronsLeft, ChevronsRight, ChevronDown, Check, Briefcase, Wrench,
  Activity, Shield,
} from 'lucide-react'
import { useAuthStore, type Workspace } from '@/stores/auth.store'
import { useMainReport } from '@/hooks/useReports'
import type { Role } from '@/types'
import { cn } from '@/lib/utils'

interface NavItem {
  label: string
  path: string
  icon: React.ReactNode
  roles?: Role[]
  badge?: number
}

interface NavSection {
  title?: string
  /** Color of the dot next to a section label (segment identity). */
  dot?: string
  items: NavItem[]
}

/**
 * Sidebar with TWO LAYOUTS — Loans & Spares.
 *
 * The user picks a workspace at the top of the sidebar and ONLY that
 * workspace's nav items are visible. This solves the "Spares vs Loans
 * is confusing" problem by making them strictly separate experiences,
 * one click apart via the workspace switcher.
 *
 * The sidebar itself is also COLLAPSIBLE — toggle button at the top-left.
 * In collapsed mode, only icons show; tooltips on hover would be a follow-up.
 */

function buildLoansNav(pendingRepayments: number, overdueInstallments: number): NavSection[] {
  return [
    {
      items: [
        { label: 'Dashboard', path: '/', icon: <LayoutDashboard size={18} /> },
      ],
    },
    {
      title: 'Loan Management',
      dot: 'bg-accent',
      items: [
        { label: 'Borrowers', path: '/borrowers', icon: <Users size={18} />, roles: ['admin', 'manager', 'loan_officer'] },
        { label: 'Loans', path: '/loans', icon: <HandCoins size={18} />, roles: ['admin', 'manager', 'loan_officer'] },
        { label: 'Repayments', path: '/repayments', icon: <CircleDollarSign size={18} />, roles: ['admin', 'manager', 'loan_officer', 'cashier', 'accountant'], badge: pendingRepayments },
        { label: 'Loan Products', path: '/loan-products', icon: <FileText size={18} />, roles: ['admin', 'manager'] },
      ],
    },
    {
      title: 'Motorcycles',
      dot: 'bg-warning',
      items: [
        { label: 'Inventory', path: '/motorcycles', icon: <Bike size={18} />, roles: ['admin', 'manager', 'cashier', 'loan_officer'] },
        { label: 'Cash Sales', path: '/cash-sales', icon: <Wallet size={18} />, roles: ['admin', 'manager', 'cashier'] },
      ],
    },
    {
      title: 'Daily Boda',
      dot: 'bg-role-admin',
      items: [
        { label: 'Fleet', path: '/daily-boda/motorcycles', icon: <Bike size={18} />, roles: ['admin', 'manager'] },
        { label: 'Drivers', path: '/daily-boda/drivers', icon: <Users size={18} />, roles: ['admin', 'manager'] },
        { label: 'Payments', path: '/daily-boda/payments', icon: <CalendarClock size={18} />, roles: ['admin', 'manager', 'cashier', 'loan_officer', 'accountant'] },
      ],
    },
    {
      title: 'Reports',
      dot: 'bg-info',
      items: [
        { label: 'Performance', path: '/reports/performance', icon: <TrendingUp size={18} />, roles: ['admin', 'manager', 'accountant'] },
        { label: 'Loan Portfolio', path: '/reports/loans', icon: <HandCoins size={18} />, roles: ['admin', 'manager', 'loan_officer', 'accountant'], badge: overdueInstallments },
        { label: 'Collections', path: '/reports/collections', icon: <CircleDollarSign size={18} />, roles: ['admin', 'manager', 'loan_officer', 'accountant'] },
        { label: 'Motorcycle Sales', path: '/reports/motorcycles', icon: <Bike size={18} />, roles: ['admin', 'manager', 'accountant'] },
        { label: 'Daily Boda', path: '/reports/daily-boda', icon: <CalendarClock size={18} />, roles: ['admin', 'manager', 'accountant'] },
      ],
    },
  ]
}

function buildSparesNav(): NavSection[] {
  return [
    {
      items: [
        { label: 'Dashboard', path: '/', icon: <LayoutDashboard size={18} /> },
      ],
    },
    {
      title: 'Operations',
      dot: 'bg-success',
      items: [
        { label: 'POS', path: '/pos', icon: <ShoppingCart size={18} />, roles: ['admin', 'manager', 'cashier'] },
        { label: 'Sales History', path: '/sales', icon: <Receipt size={18} />, roles: ['admin', 'manager'] },
      ],
    },
    {
      title: 'Inventory',
      dot: 'bg-warning',
      items: [
        { label: 'Products', path: '/products', icon: <Package size={18} />, roles: ['admin', 'manager'] },
        { label: 'Stock Levels', path: '/stock', icon: <Boxes size={18} />, roles: ['admin', 'manager', 'stock_clerk'] },
        { label: 'Stock In', path: '/stock/in', icon: <PackagePlus size={18} />, roles: ['admin', 'manager', 'stock_clerk'] },
        { label: 'Transfer', path: '/stock/transfer', icon: <ArrowLeftRight size={18} />, roles: ['admin'] },
        { label: 'Movements', path: '/stock/movements', icon: <ClipboardList size={18} />, roles: ['admin', 'manager'] },
      ],
    },
    {
      title: 'Reports',
      dot: 'bg-info',
      items: [
        { label: 'Performance', path: '/reports/performance', icon: <TrendingUp size={18} />, roles: ['admin', 'manager', 'accountant'] },
        { label: 'Daily Sales', path: '/reports/daily', icon: <BarChart3 size={18} />, roles: ['admin', 'manager', 'accountant'] },
        { label: 'Stock Report', path: '/reports/stock', icon: <Boxes size={18} />, roles: ['admin', 'manager'] },
        { label: 'Profit & Loss', path: '/reports/pnl', icon: <TrendingUp size={18} />, roles: ['admin', 'accountant'] },
        { label: 'By Cashier', path: '/reports/cashiers', icon: <PieChart size={18} />, roles: ['admin', 'accountant'] },
      ],
    },
  ]
}

const SHARED_ADMIN: NavSection = {
  title: 'Manage',
  dot: 'bg-foreground-muted',
  items: [
    { label: 'Staff', path: '/staff', icon: <Users size={18} />, roles: ['admin'] },
    { label: 'Branches', path: '/branches', icon: <Building2 size={18} />, roles: ['admin'] },
    { label: 'Activities', path: '/activities', icon: <Activity size={18} />, roles: ['admin'] },
    { label: 'Security', path: '/system/security', icon: <Shield size={18} />, roles: ['admin'] },
    { label: 'Observability', path: '/system/observability', icon: <Activity size={18} />, roles: ['admin'] },
    { label: 'Settings', path: '/settings', icon: <Settings size={18} />, roles: ['admin'] },
  ],
}

const WORKSPACES: { id: Workspace; label: string; description: string; icon: React.ReactNode; accent: string }[] = [
  {
    id: 'loans',
    label: 'Loans & Motorcycles',
    description: 'Loan management, motorcycles, daily boda',
    icon: <HandCoins size={16} />,
    accent: 'text-accent',
  },
  {
    id: 'spares',
    label: 'Spares POS',
    description: 'Spare parts inventory and sales',
    icon: <Wrench size={16} />,
    accent: 'text-success-dark',
  },
]

export function AppSidebar({
  collapsed = false,
  onToggleCollapsed,
  onOpenCommandPalette,
}: {
  collapsed?: boolean
  onToggleCollapsed?: () => void
  onOpenCommandPalette?: () => void
}) {
  const location = useLocation()
  const role = useAuthStore((s) => s.currentRole) as Role | null
  const workspace = useAuthStore((s) => s.currentWorkspace)
  const setWorkspace = useAuthStore((s) => s.setWorkspace)
  const access = useAuthStore((s) => s.getWorkspaceAccess())
  const { data: report } = useMainReport()

  const sections = workspace === 'loans'
    ? buildLoansNav(
        report?.repayments?.pending_verification || 0,
        report?.repayments?.overdue_installments || 0,
      )
    : buildSparesNav()

  // Always show admin section at the bottom (it's cross-segment).
  const allSections = [...sections, SHARED_ADMIN]

  const isActive = (path: string) => {
    if (path === '/') return location.pathname === '/'
    return location.pathname === path || location.pathname.startsWith(path + '/')
  }

  return (
    <aside
      className={cn(
        'bg-sidebar text-foreground border-r border-border shrink-0 hidden lg:flex flex-col h-screen sticky top-0 transition-[width] duration-200',
        collapsed ? 'w-16' : 'w-64',
      )}
    >
      {/* Brand + collapse toggle */}
      <div className="h-14 flex items-center gap-2 px-3 border-b border-border-subtle">
        <button
          type="button"
          onClick={onToggleCollapsed}
          className="h-8 w-8 rounded-lg text-foreground-muted hover:bg-surface-hover hover:text-foreground inline-flex items-center justify-center shrink-0"
          aria-label={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
        >
          {collapsed ? <ChevronsRight size={16} /> : <ChevronsLeft size={16} />}
        </button>
        {!collapsed && (
          <Link to="/" className="flex items-center gap-2 min-w-0">
            <img src="/grit_logo.png" alt="" className="w-7 h-7 rounded-lg shrink-0" />
            <span className="text-[14px] font-semibold tracking-tight truncate">Grit Motors</span>
          </Link>
        )}
      </div>

      {/* Workspace switcher — hidden when user only has access to one workspace */}
      {!collapsed && access === 'both' && (
        <div className="px-3 pt-3">
          <WorkspaceSwitcher current={workspace} onPick={setWorkspace} />
        </div>
      )}

      {/* Cmd+K trigger */}
      {!collapsed && onOpenCommandPalette && (
        <button
          type="button"
          onClick={onOpenCommandPalette}
          className="mx-3 mt-3 flex items-center gap-2 px-3 h-8 rounded-lg bg-surface-2 border border-border text-[12.5px] text-foreground-muted hover:bg-surface-hover hover:text-foreground-secondary transition"
        >
          <Search size={13} />
          <span className="flex-1 text-left">Search anywhere</span>
          <kbd className="text-[10px] bg-surface px-1.5 py-0.5 rounded border border-border">⌘K</kbd>
        </button>
      )}

      {/* Navigation */}
      <nav className="flex-1 overflow-y-auto py-3 px-2 space-y-3">
        {allSections.map((section, idx) => {
          const items = section.items.filter((i) => !i.roles || (role && i.roles.includes(role)))
          if (items.length === 0) return null

          return (
            <div key={section.title || `section-${idx}`}>
              {section.title && !collapsed && (
                <div className="flex items-center gap-2 px-3 mb-1.5">
                  {section.dot && <span className={cn('h-1.5 w-1.5 rounded-full', section.dot)} />}
                  <p className="text-[10.5px] font-semibold uppercase tracking-wider text-foreground-muted">
                    {section.title}
                  </p>
                </div>
              )}
              <div className="space-y-0.5">
                {items.map((item) => {
                  const active = isActive(item.path)
                  return (
                    <Link
                      key={item.path}
                      to={item.path}
                      // Preload the route chunk as soon as the sidebar mounts
                      // — by the time the user clicks, everything is in
                      // memory and nav is instant. The default 'intent' only
                      // fetches on hover, which leaves a click-without-hover
                      // (touch / fast click) waiting for the chunk download.
                      preload="render"
                      className={cn(
                        'group flex items-center gap-2.5 h-9 rounded-md text-[13px] font-medium transition-colors relative',
                        collapsed ? 'justify-center px-0' : 'px-3',
                        active
                          ? 'bg-sidebar-active text-accent'
                          : 'text-foreground-secondary hover:bg-sidebar-hover hover:text-foreground',
                      )}
                      title={collapsed ? item.label : undefined}
                    >
                      {active && !collapsed && <span className="absolute left-0 inset-y-1.5 w-[2px] rounded-r-full bg-accent" />}
                      <span className={cn('shrink-0', active ? 'text-accent' : 'text-foreground-muted group-hover:text-foreground-secondary')}>
                        {item.icon}
                      </span>
                      {!collapsed && (
                        <>
                          <span className="flex-1 truncate">{item.label}</span>
                          {item.badge !== undefined && item.badge > 0 && (
                            <span className="px-1.5 py-0.5 rounded-full bg-warning-light text-warning-dark text-[10px] font-bold min-w-[18px] text-center">
                              {item.badge}
                            </span>
                          )}
                        </>
                      )}
                    </Link>
                  )
                })}
              </div>
            </div>
          )
        })}
      </nav>
    </aside>
  )
}

/* ─────────────── Workspace switcher (top of sidebar) ─────────────── */

function WorkspaceSwitcher({ current, onPick }: { current: Workspace; onPick: (w: Workspace) => void }) {
  const [open, setOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)
  const cur = WORKSPACES.find((w) => w.id === current)!

  useEffect(() => {
    if (!open) return
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false)
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [open])

  return (
    <div ref={ref} className="relative">
      <button
        type="button"
        onClick={() => setOpen(!open)}
        className="w-full flex items-center gap-2.5 px-3 py-2 rounded-lg bg-surface border border-border hover:bg-surface-hover transition text-left"
      >
        <span className={cn('shrink-0 w-7 h-7 rounded-md bg-surface-2 flex items-center justify-center', cur.accent)}>
          {cur.icon}
        </span>
        <span className="flex-1 min-w-0">
          <span className="block text-[10.5px] font-semibold uppercase tracking-wider text-foreground-muted">
            Workspace
          </span>
          <span className="block text-[12.5px] font-semibold text-foreground truncate">
            {cur.label}
          </span>
        </span>
        <ChevronDown size={14} className={cn('text-foreground-muted transition-transform shrink-0', open && 'rotate-180')} />
      </button>

      {open && (
        <div className="absolute left-0 right-0 top-full mt-1 z-50 bg-surface border border-border rounded-lg shadow-md overflow-hidden">
          {WORKSPACES.map((w) => (
            <button
              key={w.id}
              type="button"
              onClick={() => {
                onPick(w.id)
                setOpen(false)
              }}
              className={cn(
                'w-full flex items-center gap-2.5 px-3 py-2.5 text-left transition',
                w.id === current ? 'bg-sidebar-active' : 'hover:bg-surface-hover',
              )}
            >
              <span className={cn('shrink-0 w-7 h-7 rounded-md bg-surface-2 flex items-center justify-center', w.accent)}>
                {w.icon}
              </span>
              <span className="flex-1 min-w-0">
                <span className="block text-[12.5px] font-semibold text-foreground">{w.label}</span>
                <span className="block text-[11px] text-foreground-muted truncate">{w.description}</span>
              </span>
              {w.id === current && <Check size={14} className="text-accent shrink-0" />}
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
