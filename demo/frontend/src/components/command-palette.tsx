import { useEffect, useState, useMemo, useRef } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { Search, ArrowRight, Bike, HandCoins, Users, Wallet, Package, ShoppingCart, BarChart3, Settings, Building2, FileText, Boxes, ClipboardList, ArrowLeftRight, PackagePlus, UserPlus, PieChart, TrendingUp, CircleDollarSign, CalendarClock, Receipt } from 'lucide-react'
import { useAuthStore } from '@/stores/auth.store'
import type { Role } from '@/types'
import { cn } from '@/lib/utils'

/**
 * Command palette — Cmd/Ctrl+K opens a fuzzy-search overlay with every nav
 * destination in the app. The fastest way to jump anywhere without taking
 * a hand off the keyboard. Heavily inspired by Linear / Slack / VS Code.
 *
 * The palette is mounted once at the app layout level and toggled by a
 * global key handler. Each command is just a navigate target; commands
 * are visible only if the user's role permits the destination.
 */

interface Command {
  id: string
  label: string
  segment: 'general' | 'spares' | 'motorcycles' | 'loans' | 'daily-boda' | 'reports' | 'admin'
  icon: React.ReactNode
  path: string
  roles?: Role[]
  keywords?: string[]
}

const COMMANDS: Command[] = [
  // General
  { id: 'main-dash', label: 'Main Dashboard', segment: 'general', icon: <BarChart3 size={14} />, path: '/', keywords: ['home', 'overview'] },
  { id: 'spares-dash', label: 'Spares Dashboard', segment: 'spares', icon: <ShoppingCart size={14} />, path: '/spares', roles: ['admin', 'manager', 'cashier', 'stock_clerk', 'accountant'] },
  { id: 'loans-dash', label: 'Loans Dashboard', segment: 'loans', icon: <CircleDollarSign size={14} />, path: '/loans-overview', roles: ['admin', 'manager', 'loan_officer', 'accountant'] },

  // Spares
  { id: 'pos', label: 'Open POS', segment: 'spares', icon: <ShoppingCart size={14} />, path: '/pos', roles: ['admin', 'manager', 'cashier'], keywords: ['sell', 'checkout'] },
  { id: 'products', label: 'Products', segment: 'spares', icon: <Package size={14} />, path: '/products', roles: ['admin', 'manager'] },
  { id: 'stock', label: 'Stock Levels', segment: 'spares', icon: <Boxes size={14} />, path: '/stock', roles: ['admin', 'manager', 'stock_clerk'] },
  { id: 'stock-in', label: 'Stock In', segment: 'spares', icon: <PackagePlus size={14} />, path: '/stock/in', roles: ['admin', 'manager', 'stock_clerk'] },
  { id: 'stock-transfer', label: 'Stock Transfer', segment: 'spares', icon: <ArrowLeftRight size={14} />, path: '/stock/transfer', roles: ['admin'] },
  { id: 'stock-movements', label: 'Stock Movements', segment: 'spares', icon: <ClipboardList size={14} />, path: '/stock/movements', roles: ['admin', 'manager'] },
  { id: 'sales-history', label: 'Sales History', segment: 'spares', icon: <Receipt size={14} />, path: '/sales', roles: ['admin', 'manager'] },

  // Motorcycles
  { id: 'motorcycles', label: 'Motorcycle Inventory', segment: 'motorcycles', icon: <Bike size={14} />, path: '/motorcycles', roles: ['admin', 'manager', 'cashier', 'loan_officer'] },
  { id: 'cash-sales', label: 'Cash Sales', segment: 'motorcycles', icon: <Wallet size={14} />, path: '/cash-sales', roles: ['admin', 'manager', 'cashier'] },

  // Loans
  { id: 'loan-products', label: 'Loan Products', segment: 'loans', icon: <FileText size={14} />, path: '/loan-products', roles: ['admin', 'manager'] },
  { id: 'borrowers', label: 'Borrowers', segment: 'loans', icon: <Users size={14} />, path: '/borrowers', roles: ['admin', 'manager', 'loan_officer'] },
  { id: 'loans', label: 'Loans', segment: 'loans', icon: <HandCoins size={14} />, path: '/loans', roles: ['admin', 'manager', 'loan_officer'] },
  { id: 'new-loan', label: 'New Loan', segment: 'loans', icon: <HandCoins size={14} />, path: '/loans/new', roles: ['admin', 'manager', 'loan_officer'], keywords: ['create', 'add'] },
  { id: 'repayments', label: 'Repayments', segment: 'loans', icon: <CircleDollarSign size={14} />, path: '/repayments', roles: ['admin', 'manager', 'loan_officer', 'cashier', 'accountant'] },

  // Daily Boda
  { id: 'db-motorcycles', label: 'Daily Boda — Motorcycles', segment: 'daily-boda', icon: <Bike size={14} />, path: '/daily-boda/motorcycles', roles: ['admin', 'manager'] },
  { id: 'db-drivers', label: 'Daily Boda — Drivers', segment: 'daily-boda', icon: <Users size={14} />, path: '/daily-boda/drivers', roles: ['admin', 'manager'] },
  { id: 'db-payments', label: 'Daily Boda — Payments', segment: 'daily-boda', icon: <CalendarClock size={14} />, path: '/daily-boda/payments', roles: ['admin', 'manager', 'cashier', 'loan_officer', 'accountant'] },

  // Reports
  { id: 'r-daily', label: 'Report — Daily Sales', segment: 'reports', icon: <BarChart3 size={14} />, path: '/reports/daily', roles: ['admin', 'manager', 'accountant'] },
  { id: 'r-stock', label: 'Report — Stock', segment: 'reports', icon: <Boxes size={14} />, path: '/reports/stock', roles: ['admin', 'manager'] },
  { id: 'r-pnl', label: 'Report — Profit & Loss', segment: 'reports', icon: <TrendingUp size={14} />, path: '/reports/pnl', roles: ['admin', 'accountant'] },
  { id: 'r-cashiers', label: 'Report — By Cashier', segment: 'reports', icon: <PieChart size={14} />, path: '/reports/cashiers', roles: ['admin', 'accountant'] },
  { id: 'r-loans', label: 'Report — Loan Portfolio', segment: 'reports', icon: <HandCoins size={14} />, path: '/reports/loans', roles: ['admin', 'manager', 'loan_officer', 'accountant'] },
  { id: 'r-collections', label: 'Report — Collections', segment: 'reports', icon: <CircleDollarSign size={14} />, path: '/reports/collections', roles: ['admin', 'manager', 'loan_officer', 'accountant'] },
  { id: 'r-mc-sales', label: 'Report — Motorcycle Sales', segment: 'reports', icon: <Bike size={14} />, path: '/reports/motorcycles', roles: ['admin', 'manager', 'accountant'] },
  { id: 'r-db', label: 'Report — Daily Boda', segment: 'reports', icon: <CalendarClock size={14} />, path: '/reports/daily-boda', roles: ['admin', 'manager', 'accountant'] },

  // Admin
  { id: 'staff', label: 'Staff', segment: 'admin', icon: <Users size={14} />, path: '/staff', roles: ['admin'] },
  { id: 'invite-staff', label: 'Invite Staff', segment: 'admin', icon: <UserPlus size={14} />, path: '/staff/invite', roles: ['admin'] },
  { id: 'branches', label: 'Branches', segment: 'admin', icon: <Building2 size={14} />, path: '/branches', roles: ['admin'] },
  { id: 'settings', label: 'Settings', segment: 'admin', icon: <Settings size={14} />, path: '/settings', roles: ['admin'] },
]

const SEGMENT_LABEL: Record<Command['segment'], string> = {
  general: 'General',
  spares: 'Spares',
  motorcycles: 'Motorcycles',
  loans: 'Loans',
  'daily-boda': 'Daily Boda',
  reports: 'Reports',
  admin: 'Admin',
}

/** Hook the layout uses to render the palette + listen for the open shortcut. */
export function useCommandPalette() {
  const [open, setOpen] = useState(false)

  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
        e.preventDefault()
        setOpen((o) => !o)
      }
    }
    window.addEventListener('keydown', handler)
    return () => window.removeEventListener('keydown', handler)
  }, [])

  return { open, setOpen }
}

export function CommandPalette({ open, onOpenChange }: { open: boolean; onOpenChange: (o: boolean) => void }) {
  const navigate = useNavigate()
  const role = useAuthStore((s) => s.currentRole) as Role | null
  const [query, setQuery] = useState('')
  const [highlight, setHighlight] = useState(0)
  const inputRef = useRef<HTMLInputElement>(null)

  const allowed = useMemo(() => COMMANDS.filter((c) => !c.roles || (role && c.roles.includes(role))), [role])

  const filtered = useMemo(() => {
    if (!query.trim()) return allowed
    const q = query.toLowerCase()
    return allowed.filter((c) => {
      if (c.label.toLowerCase().includes(q)) return true
      if (c.keywords?.some((k) => k.toLowerCase().includes(q))) return true
      if (SEGMENT_LABEL[c.segment].toLowerCase().includes(q)) return true
      return false
    })
  }, [allowed, query])

  // Group by segment for visual scanning
  const grouped = useMemo(() => {
    const out: Record<string, Command[]> = {}
    for (const c of filtered) {
      const key = SEGMENT_LABEL[c.segment]
      if (!out[key]) out[key] = []
      out[key].push(c)
    }
    return out
  }, [filtered])

  useEffect(() => {
    setHighlight(0)
  }, [query])

  useEffect(() => {
    if (open) {
      setQuery('')
      setHighlight(0)
      setTimeout(() => inputRef.current?.focus(), 50)
    }
  }, [open])

  const onKey = (e: React.KeyboardEvent) => {
    if (e.key === 'Escape') {
      e.preventDefault()
      onOpenChange(false)
    } else if (e.key === 'ArrowDown') {
      e.preventDefault()
      setHighlight((h) => Math.min(h + 1, filtered.length - 1))
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()
      setHighlight((h) => Math.max(h - 1, 0))
    } else if (e.key === 'Enter') {
      e.preventDefault()
      const cmd = filtered[highlight]
      if (cmd) {
        navigate({ to: cmd.path })
        onOpenChange(false)
      }
    }
  }

  if (!open) return null

  // Build a flat index so we can match `highlight` against grouped rendering.
  let runningIdx = -1

  return (
    <div className="fixed inset-0 z-50 bg-black/40 flex items-start justify-center pt-[15vh] px-4" onClick={() => onOpenChange(false)}>
      <div
        className="w-full max-w-xl bg-surface rounded-2xl border border-border shadow-2xl overflow-hidden"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center gap-3 px-4 py-3 border-b border-border">
          <Search size={16} className="text-foreground-muted" />
          <input
            ref={inputRef}
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onKeyDown={onKey}
            placeholder="Search pages, actions..."
            className="flex-1 bg-transparent text-[14px] text-foreground placeholder:text-foreground-muted focus:outline-none"
          />
          <kbd className="text-[10.5px] text-foreground-muted px-1.5 py-0.5 rounded border border-border">esc</kbd>
        </div>
        <div className="max-h-[50vh] overflow-y-auto py-1">
          {filtered.length === 0 && (
            <div className="px-4 py-8 text-center text-[13px] text-foreground-muted">No matches.</div>
          )}
          {Object.entries(grouped).map(([segment, cmds]) => (
            <div key={segment} className="py-1">
              <div className="px-4 pt-2 pb-1 text-[10.5px] font-semibold uppercase tracking-wider text-foreground-muted">
                {segment}
              </div>
              {cmds.map((cmd) => {
                runningIdx += 1
                const isActive = runningIdx === highlight
                return (
                  <button
                    key={cmd.id}
                    type="button"
                    onMouseEnter={() => setHighlight(runningIdx)}
                    onClick={() => {
                      navigate({ to: cmd.path })
                      onOpenChange(false)
                    }}
                    className={cn(
                      'w-full flex items-center gap-3 px-4 py-2 text-left',
                      isActive && 'bg-accent-tint',
                    )}
                  >
                    <span className={cn('flex h-6 w-6 items-center justify-center rounded-md', isActive ? 'text-accent' : 'text-foreground-muted')}>
                      {cmd.icon}
                    </span>
                    <span className={cn('flex-1 text-[13px]', isActive ? 'text-accent font-medium' : 'text-foreground')}>
                      {cmd.label}
                    </span>
                    {isActive && <ArrowRight size={14} className="text-accent" />}
                  </button>
                )
              })}
            </div>
          ))}
        </div>
        <div className="px-4 py-2 border-t border-border-subtle bg-surface-2 flex items-center justify-between text-[11px] text-foreground-muted">
          <div className="flex items-center gap-3">
            <span><kbd className="px-1 rounded border border-border">↑</kbd> <kbd className="px-1 rounded border border-border">↓</kbd> navigate</span>
            <span><kbd className="px-1 rounded border border-border">↵</kbd> open</span>
          </div>
          <span><kbd className="px-1 rounded border border-border">⌘K</kbd> toggle</span>
        </div>
      </div>
    </div>
  )
}
