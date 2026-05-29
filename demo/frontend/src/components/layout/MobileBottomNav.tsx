import { Link, useLocation, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import {
  LayoutDashboard, ShoppingCart, Package, Boxes, Receipt, BarChart3, TrendingUp,
  Users, Building2, Settings, ArrowLeftRight, ClipboardList, PackagePlus, UserPlus,
  PieChart, Bike, HandCoins, FileText, CalendarClock, Wallet, CircleDollarSign,
  Search, ChevronDown, Check, Wrench, Menu, LogOut, X, User as UserIcon, Shield,
} from 'lucide-react'
import { useAuthStore, type Workspace } from '@/stores/auth.store'
import { useMainReport } from '@/hooks/useReports'
import type { Role } from '@/types'
import { cn } from '@/lib/utils'

/**
 * Mobile bottom-tab navigation.
 *
 * Replaces the hamburger menu pattern. Shows 4 primary destinations for the
 * current workspace plus a "More" button that opens a full-height sheet
 * with the workspace switcher, every remaining nav item, and the user
 * profile / sign-out controls.
 *
 * Hidden on lg+ where the sidebar takes over.
 */

interface BottomTab {
  label: string
  path: string
  icon: React.ReactNode
  match?: string[]
}

function loansPrimary(): BottomTab[] {
  return [
    { label: 'Home', path: '/', icon: <LayoutDashboard size={20} /> },
    { label: 'Loans', path: '/loans', icon: <HandCoins size={20} />, match: ['/loans'] },
    { label: 'Borrowers', path: '/borrowers', icon: <Users size={20} /> },
    { label: 'Repayments', path: '/repayments', icon: <CircleDollarSign size={20} /> },
  ]
}

function sparesPrimary(): BottomTab[] {
  return [
    { label: 'Home', path: '/', icon: <LayoutDashboard size={20} /> },
    { label: 'POS', path: '/pos', icon: <ShoppingCart size={20} /> },
    { label: 'Products', path: '/products', icon: <Package size={20} /> },
    { label: 'Sales', path: '/sales', icon: <Receipt size={20} /> },
  ]
}

export function MobileBottomNav({ onOpenSheet }: { onOpenSheet: () => void }) {
  const location = useLocation()
  const workspace = useAuthStore((s) => s.currentWorkspace)
  const role = useAuthStore((s) => s.currentRole) as Role | null

  const tabs = workspace === 'loans' ? loansPrimary() : sparesPrimary()

  const isActive = (tab: BottomTab) => {
    const path = tab.path
    if (path === '/') return location.pathname === '/'
    if (tab.match) return tab.match.some((m) => location.pathname === m || location.pathname.startsWith(m + '/'))
    return location.pathname === path || location.pathname.startsWith(path + '/')
  }

  // POS owns its own full-screen layout — no bottom nav needed there.
  if (location.pathname === '/pos') return null

  return (
    <nav
      className="lg:hidden fixed bottom-0 inset-x-0 z-30 bg-surface border-t border-border flex items-stretch h-16 pb-[env(safe-area-inset-bottom)]"
      aria-label="Primary"
    >
      {tabs.map((tab) => {
        const active = isActive(tab)
        // Cashiers without permission to a target shouldn't see the tab —
        // but we let everyone tap each primary, and the page itself enforces
        // role gates. Keeps the bottom bar consistent across roles.
        return (
          <Link
            key={tab.path}
            to={tab.path}
            preload="render"
            className={cn(
              'flex-1 flex flex-col items-center justify-center gap-0.5 text-[10.5px] font-medium transition-colors',
              active ? 'text-accent' : 'text-foreground-muted',
            )}
          >
            <span className={active ? 'text-accent' : 'text-foreground-muted'}>{tab.icon}</span>
            {tab.label}
          </Link>
        )
      })}
      <button
        type="button"
        onClick={onOpenSheet}
        className="flex-1 flex flex-col items-center justify-center gap-0.5 text-[10.5px] font-medium text-foreground-muted hover:text-foreground"
        aria-label="More menu"
      >
        <Menu size={20} />
        More
      </button>
    </nav>
  )
}

/* ─────────────── Sheet menu (the "More" target) ─────────────── */

interface NavItem {
  label: string
  path: string
  icon: React.ReactNode
  roles?: Role[]
  badge?: number
}

interface NavSection {
  title?: string
  dot?: string
  items: NavItem[]
}

function buildLoansNav(pendingRepayments: number, overdueInstallments: number): NavSection[] {
  return [
    { items: [{ label: 'Dashboard', path: '/', icon: <LayoutDashboard size={18} /> }] },
    {
      title: 'Loan Management', dot: 'bg-accent',
      items: [
        { label: 'Borrowers', path: '/borrowers', icon: <Users size={18} />, roles: ['admin', 'manager', 'loan_officer'] },
        { label: 'Loans', path: '/loans', icon: <HandCoins size={18} />, roles: ['admin', 'manager', 'loan_officer'] },
        { label: 'Repayments', path: '/repayments', icon: <CircleDollarSign size={18} />, roles: ['admin', 'manager', 'loan_officer', 'cashier', 'accountant'], badge: pendingRepayments },
        { label: 'Loan Products', path: '/loan-products', icon: <FileText size={18} />, roles: ['admin', 'manager'] },
      ],
    },
    {
      title: 'Motorcycles', dot: 'bg-warning',
      items: [
        { label: 'Inventory', path: '/motorcycles', icon: <Bike size={18} />, roles: ['admin', 'manager', 'cashier', 'loan_officer'] },
        { label: 'Cash Sales', path: '/cash-sales', icon: <Wallet size={18} />, roles: ['admin', 'manager', 'cashier'] },
      ],
    },
    {
      title: 'Daily Boda', dot: 'bg-role-admin',
      items: [
        { label: 'Fleet', path: '/daily-boda/motorcycles', icon: <Bike size={18} />, roles: ['admin', 'manager'] },
        { label: 'Drivers', path: '/daily-boda/drivers', icon: <Users size={18} />, roles: ['admin', 'manager'] },
        { label: 'Payments', path: '/daily-boda/payments', icon: <CalendarClock size={18} />, roles: ['admin', 'manager', 'cashier', 'loan_officer', 'accountant'] },
      ],
    },
    {
      title: 'Reports', dot: 'bg-info',
      items: [
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
    { items: [{ label: 'Dashboard', path: '/', icon: <LayoutDashboard size={18} /> }] },
    {
      title: 'Operations', dot: 'bg-success',
      items: [
        { label: 'POS', path: '/pos', icon: <ShoppingCart size={18} />, roles: ['admin', 'manager', 'cashier'] },
        { label: 'Sales History', path: '/sales', icon: <Receipt size={18} />, roles: ['admin', 'manager'] },
      ],
    },
    {
      title: 'Inventory', dot: 'bg-warning',
      items: [
        { label: 'Products', path: '/products', icon: <Package size={18} />, roles: ['admin', 'manager'] },
        { label: 'Stock Levels', path: '/stock', icon: <Boxes size={18} />, roles: ['admin', 'manager', 'stock_clerk'] },
        { label: 'Stock In', path: '/stock/in', icon: <PackagePlus size={18} />, roles: ['admin', 'manager', 'stock_clerk'] },
        { label: 'Transfer', path: '/stock/transfer', icon: <ArrowLeftRight size={18} />, roles: ['admin'] },
        { label: 'Movements', path: '/stock/movements', icon: <ClipboardList size={18} />, roles: ['admin', 'manager'] },
      ],
    },
    {
      title: 'Reports', dot: 'bg-info',
      items: [
        { label: 'Daily Sales', path: '/reports/daily', icon: <BarChart3 size={18} />, roles: ['admin', 'manager', 'accountant'] },
        { label: 'Stock Report', path: '/reports/stock', icon: <Boxes size={18} />, roles: ['admin', 'manager'] },
        { label: 'Profit & Loss', path: '/reports/pnl', icon: <TrendingUp size={18} />, roles: ['admin', 'accountant'] },
        { label: 'By Cashier', path: '/reports/cashiers', icon: <PieChart size={18} />, roles: ['admin', 'accountant'] },
      ],
    },
  ]
}

const SHARED_ADMIN: NavSection = {
  title: 'Manage', dot: 'bg-foreground-muted',
  items: [
    { label: 'Staff', path: '/staff', icon: <Users size={18} />, roles: ['admin'] },
    { label: 'Branches', path: '/branches', icon: <Building2 size={18} />, roles: ['admin'] },
    { label: 'Settings', path: '/settings', icon: <Settings size={18} />, roles: ['admin'] },
  ],
}

const WORKSPACES: { id: Workspace; label: string; description: string; icon: React.ReactNode; accent: string }[] = [
  { id: 'loans', label: 'Loans & Motorcycles', description: 'Loans, motorcycles, daily boda', icon: <HandCoins size={16} />, accent: 'text-accent' },
  { id: 'spares', label: 'Spares POS', description: 'Spare parts inventory and sales', icon: <Wrench size={16} />, accent: 'text-success-dark' },
]

export function MobileSheetMenu({ open, onOpenChange, onOpenCommandPalette }: {
  open: boolean
  onOpenChange: (o: boolean) => void
  onOpenCommandPalette?: () => void
}) {
  const location = useLocation()
  const navigate = useNavigate()
  const role = useAuthStore((s) => s.currentRole) as Role | null
  const user = useAuthStore((s) => s.user)
  const businesses = useAuthStore((s) => s.businesses)
  const currentBusinessID = useAuthStore((s) => s.currentBusinessID)
  const switchBusiness = useAuthStore((s) => s.switchBusiness)
  const workspace = useAuthStore((s) => s.currentWorkspace)
  const setWorkspace = useAuthStore((s) => s.setWorkspace)
  const access = useAuthStore((s) => s.getWorkspaceAccess())
  const logout = useAuthStore((s) => s.logout)
  const { data: report } = useMainReport()
  const [businessOpen, setBusinessOpen] = useState(false)

  if (!open) return null

  const sections = workspace === 'loans'
    ? buildLoansNav(report?.repayments?.pending_verification || 0, report?.repayments?.overdue_installments || 0)
    : buildSparesNav()
  const all = [...sections, SHARED_ADMIN]

  const isActive = (path: string) => {
    if (path === '/') return location.pathname === '/'
    return location.pathname === path || location.pathname.startsWith(path + '/')
  }

  const close = () => onOpenChange(false)
  const go = (path: string) => {
    navigate({ to: path })
    close()
  }
  const handleLogout = () => {
    logout()
    navigate({ to: '/login' })
  }

  const currentBiz = businesses.find((b) => b.id === currentBusinessID)
  const initials = (user?.name || '?').split(' ').map((s) => s[0]).slice(0, 2).join('').toUpperCase()

  return (
    <div className="lg:hidden fixed inset-0 z-50 flex flex-col" role="dialog" aria-modal="true">
      {/* Backdrop */}
      <div className="flex-1 bg-black/40" onClick={close} />

      {/* Sheet */}
      <div className="bg-surface rounded-t-2xl border-t border-border shadow-2xl flex flex-col max-h-[90vh] animate-slide-up">
        <div className="px-4 py-3 border-b border-border-subtle flex items-center justify-between shrink-0">
          <h2 className="text-[15px] font-semibold text-foreground">Menu</h2>
          <button onClick={close} className="h-8 w-8 rounded-lg hover:bg-surface-hover flex items-center justify-center text-foreground-muted" aria-label="Close">
            <X size={16} />
          </button>
        </div>

        <div className="flex-1 overflow-y-auto">
          {/* Workspace switcher — hidden when user only has access to one workspace */}
          {access === 'both' && <div className="px-4 pt-4 pb-2">
            <p className="text-[10.5px] font-semibold uppercase tracking-wider text-foreground-muted mb-1.5">Workspace</p>
            <div className="grid grid-cols-2 gap-2">
              {WORKSPACES.map((w) => {
                const active = w.id === workspace
                return (
                  <button
                    key={w.id}
                    type="button"
                    onClick={() => {
                      setWorkspace(w.id)
                      navigate({ to: '/' })
                      close()
                    }}
                    className={cn(
                      'flex items-center gap-2 px-3 py-2 rounded-lg border text-left transition',
                      active ? 'border-accent bg-accent-light' : 'border-border bg-surface hover:bg-surface-hover',
                    )}
                  >
                    <span className={cn('w-7 h-7 rounded-md bg-surface-2 flex items-center justify-center', w.accent)}>
                      {w.icon}
                    </span>
                    <span className="flex-1 min-w-0">
                      <span className="block text-[12px] font-semibold text-foreground truncate">{w.label}</span>
                      <span className="block text-[10.5px] text-foreground-muted truncate">{w.description}</span>
                    </span>
                    {active && <Check size={14} className="text-accent shrink-0" />}
                  </button>
                )
              })}
            </div>
          </div>}

          {/* Business + branch row */}
          <div className="px-4 py-2">
            <button
              type="button"
              onClick={() => setBusinessOpen((o) => !o)}
              className="w-full flex items-center gap-2 px-3 py-2 rounded-lg bg-surface-2 border border-border-subtle"
            >
              <Building2 size={14} className="text-foreground-muted" />
              <span className="flex-1 text-left text-[12.5px] font-medium text-foreground truncate">
                {currentBiz?.name || 'No business'}
              </span>
              {businesses.length > 1 && (
                <ChevronDown size={14} className={cn('text-foreground-muted transition-transform', businessOpen && 'rotate-180')} />
              )}
            </button>
            {businessOpen && businesses.length > 1 && (
              <div className="mt-1 border border-border-subtle rounded-lg overflow-hidden">
                {businesses.map((biz) => (
                  <button
                    key={biz.id}
                    type="button"
                    onClick={() => {
                      switchBusiness(biz.id)
                      setBusinessOpen(false)
                    }}
                    className={cn(
                      'w-full flex items-center gap-2 px-3 py-2 text-left text-[13px] hover:bg-surface-hover',
                      biz.id === currentBusinessID && 'text-accent font-medium bg-accent-light/40',
                    )}
                  >
                    <span className="flex-1 truncate">{biz.name}</span>
                    {biz.id === currentBusinessID && <Check size={14} className="text-accent" />}
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* Cmd+K trigger (also useful on mobile via on-screen keyboard) */}
          {onOpenCommandPalette && (
            <div className="px-4 py-1">
              <button
                type="button"
                onClick={() => { onOpenCommandPalette(); close() }}
                className="w-full flex items-center gap-2 px-3 py-2 rounded-lg bg-surface-2 border border-border-subtle text-[12.5px] text-foreground-muted hover:bg-surface-hover transition"
              >
                <Search size={13} />
                <span className="flex-1 text-left">Search anywhere</span>
              </button>
            </div>
          )}

          {/* Nav sections */}
          {all.map((section, idx) => {
            const items = section.items.filter((i) => !i.roles || (role && i.roles.includes(role)))
            if (items.length === 0) return null
            return (
              <div key={section.title || `s-${idx}`} className="px-4 py-2">
                {section.title && (
                  <div className="flex items-center gap-2 mb-1.5 px-1">
                    {section.dot && <span className={cn('h-1.5 w-1.5 rounded-full', section.dot)} />}
                    <p className="text-[10.5px] font-semibold uppercase tracking-wider text-foreground-muted">{section.title}</p>
                  </div>
                )}
                <div className="grid grid-cols-1 gap-0.5">
                  {items.map((item) => {
                    const active = isActive(item.path)
                    return (
                      <button
                        key={item.path}
                        type="button"
                        onClick={() => go(item.path)}
                        className={cn(
                          'flex items-center gap-3 h-11 px-3 rounded-lg text-[13.5px] font-medium transition-colors',
                          active ? 'bg-accent-light text-accent' : 'text-foreground hover:bg-surface-hover',
                        )}
                      >
                        <span className={active ? 'text-accent' : 'text-foreground-muted'}>{item.icon}</span>
                        <span className="flex-1 text-left">{item.label}</span>
                        {item.badge !== undefined && item.badge > 0 && (
                          <span className="px-1.5 py-0.5 rounded-full bg-warning-light text-warning-dark text-[10px] font-bold min-w-4 text-center">
                            {item.badge}
                          </span>
                        )}
                      </button>
                    )
                  })}
                </div>
              </div>
            )
          })}
        </div>

        {/* User strip + logout — always visible at the bottom of the sheet */}
        <div className="border-t border-border-subtle px-4 py-3 flex items-center gap-3 shrink-0 pb-[env(safe-area-inset-bottom)]">
          <div className="w-9 h-9 rounded-full bg-accent-light flex items-center justify-center text-accent font-semibold text-[12px] shrink-0">
            {initials}
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-[13px] font-semibold text-foreground truncate">{user?.name}</p>
            <p className="text-[11px] text-foreground-muted truncate">
              <Shield size={9} className="inline mr-0.5" />
              <span className="capitalize">{role?.replace('_', ' ')}</span>
            </p>
          </div>
          <button
            onClick={handleLogout}
            className="h-9 px-3 rounded-lg bg-danger-light text-danger-dark text-[12.5px] font-semibold hover:bg-danger/20 inline-flex items-center gap-1.5"
          >
            <LogOut size={13} /> Sign out
          </button>
        </div>
      </div>
    </div>
  )
}
