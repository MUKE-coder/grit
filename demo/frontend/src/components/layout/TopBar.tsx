import { useEffect, useRef, useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import {
  ChevronDown, LogOut, Building2, GitBranch, User, Settings as SettingsIcon, Check, Shield,
} from 'lucide-react'
import { useAuthStore } from '@/stores/auth.store'
import api from '@/lib/api'
import { useQuery } from '@tanstack/react-query'
import type { Branch } from '@/types'
import { cn } from '@/lib/utils'

/**
 * TopBar — slim header above the main content area.
 *
 * - Business name + branch selector (left)
 * - User menu dropdown (right): name, email, role, switch business,
 *   settings, sign out (visible on lg+; the mobile sheet menu replaces
 *   the user-menu controls on small screens via the bottom nav).
 *
 * No hamburger — mobile uses the dedicated MobileBottomNav + sheet menu.
 */
export function TopBar() {
  const navigate = useNavigate()
  const user = useAuthStore((s) => s.user)
  const businesses = useAuthStore((s) => s.businesses)
  const currentBusinessID = useAuthStore((s) => s.currentBusinessID)
  const currentBranchID = useAuthStore((s) => s.currentBranchID)
  const currentRole = useAuthStore((s) => s.currentRole)
  const switchBusiness = useAuthStore((s) => s.switchBusiness)
  const setCurrentBranch = useAuthStore((s) => s.setCurrentBranch)
  const logout = useAuthStore((s) => s.logout)

  const [userOpen, setUserOpen] = useState(false)
  const [branchOpen, setBranchOpen] = useState(false)
  const userMenuRef = useRef<HTMLDivElement>(null)
  const branchMenuRef = useRef<HTMLDivElement>(null)

  const currentBiz = businesses.find((b) => b.id === currentBusinessID)

  const { data: branches } = useQuery({
    queryKey: ['branches', currentBusinessID],
    queryFn: async () => {
      const res = await api.get('/branches')
      return res.data.data as Branch[]
    },
    enabled: !!currentBusinessID,
  })

  const currentBranch = branches?.find((b) => b.id === currentBranchID)

  // Close on outside click.
  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (userOpen && userMenuRef.current && !userMenuRef.current.contains(e.target as Node)) setUserOpen(false)
      if (branchOpen && branchMenuRef.current && !branchMenuRef.current.contains(e.target as Node)) setBranchOpen(false)
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [userOpen, branchOpen])

  const handleLogout = () => {
    logout()
    navigate({ to: '/login' })
  }

  const initials = (user?.name || '?').split(' ').map((s) => s[0]).slice(0, 2).join('').toUpperCase()

  return (
    <header className="h-14 bg-surface border-b border-border flex items-center justify-between px-3 lg:px-6 sticky top-0 z-30">
      {/* Left: business name + branch selector */}
      <div className="flex items-center gap-2 lg:gap-3 min-w-0">
        <div className="flex items-center gap-2 min-w-0">
          <Building2 size={14} className="text-foreground-muted shrink-0" />
          <span className="text-[13px] font-semibold text-foreground truncate">
            {currentBiz?.name || 'Select Business'}
          </span>
        </div>

        {/* Branch selector */}
        {branches && branches.length > 1 && (
          <div ref={branchMenuRef} className="relative">
            <button
              type="button"
              onClick={() => setBranchOpen(!branchOpen)}
              className="ml-2 flex items-center gap-1.5 h-8 px-2.5 rounded-lg border border-border bg-surface hover:bg-surface-hover text-[12.5px] transition"
            >
              <GitBranch size={12} className="text-foreground-muted" />
              <span className="text-foreground">{currentBranch?.name || 'All Branches'}</span>
              <ChevronDown size={12} className={cn('text-foreground-muted transition-transform', branchOpen && 'rotate-180')} />
            </button>
            {branchOpen && (
              <div className="absolute left-0 top-full mt-1 w-56 bg-surface rounded-lg border border-border shadow-md py-1 z-50">
                {branches.map((branch) => (
                  <button
                    key={branch.id}
                    type="button"
                    onClick={() => {
                      setCurrentBranch(branch.id)
                      setBranchOpen(false)
                    }}
                    className={cn(
                      'w-full text-left px-3 py-2 text-[13px] hover:bg-surface-hover flex items-center gap-2',
                      branch.id === currentBranchID ? 'text-accent font-medium' : 'text-foreground',
                    )}
                  >
                    <span className="flex-1 truncate">{branch.name}</span>
                    {branch.id === currentBranchID && <Check size={14} className="text-accent shrink-0" />}
                  </button>
                ))}
              </div>
            )}
          </div>
        )}
      </div>

      {/* Right: user menu */}
      <div ref={userMenuRef} className="relative">
        <button
          type="button"
          onClick={() => setUserOpen(!userOpen)}
          className="flex items-center gap-2 pl-2 pr-2.5 h-9 rounded-lg hover:bg-surface-hover transition"
        >
          <div className="w-7 h-7 rounded-full bg-accent-light flex items-center justify-center text-accent text-[11px] font-semibold">
            {initials}
          </div>
          <div className="hidden sm:block text-left">
            <p className="text-[12.5px] font-medium text-foreground leading-tight">{user?.name}</p>
            <p className="text-[10.5px] text-foreground-muted capitalize leading-tight">{currentRole?.replace('_', ' ')}</p>
          </div>
          <ChevronDown size={13} className={cn('text-foreground-muted transition-transform', userOpen && 'rotate-180')} />
        </button>
        {userOpen && (
          <div className="absolute right-0 top-full mt-1 w-64 bg-surface rounded-lg border border-border shadow-md py-1 z-50">
            {/* Header */}
            <div className="px-3 py-2.5 border-b border-border-subtle">
              <p className="text-[13px] font-semibold text-foreground truncate">{user?.name}</p>
              <p className="text-[11.5px] text-foreground-muted truncate">{user?.email}</p>
              <span className="inline-flex mt-1.5 px-1.5 py-0.5 rounded-full text-[10.5px] font-medium bg-accent-light text-accent capitalize">
                <Shield size={10} className="inline mr-1 mt-px" />
                {currentRole?.replace('_', ' ')}
              </span>
            </div>

            {/* Switch business */}
            {businesses.length > 1 && (
              <div className="border-b border-border-subtle py-1">
                <p className="px-3 py-1 text-[10.5px] font-semibold uppercase tracking-wider text-foreground-muted">
                  Switch Business
                </p>
                {businesses.map((biz) => (
                  <button
                    key={biz.id}
                    type="button"
                    onClick={() => {
                      switchBusiness(biz.id)
                      setUserOpen(false)
                    }}
                    className={cn(
                      'w-full text-left px-3 py-1.5 text-[13px] hover:bg-surface-hover flex items-center gap-2',
                      biz.id === currentBusinessID ? 'text-accent font-medium' : 'text-foreground',
                    )}
                  >
                    <span className="flex-1 truncate">{biz.name}</span>
                    {biz.id === currentBusinessID && <Check size={14} className="text-accent shrink-0" />}
                  </button>
                ))}
              </div>
            )}

            {/* Actions */}
            <div className="py-1">
              <MenuItem
                icon={<User size={14} />}
                label="Profile"
                onClick={() => { setUserOpen(false); navigate({ to: '/settings' }) }}
              />
              <MenuItem
                icon={<SettingsIcon size={14} />}
                label="Settings"
                onClick={() => { setUserOpen(false); navigate({ to: '/settings' }) }}
              />
            </div>

            {/* Sign out */}
            <div className="border-t border-border-subtle py-1">
              <MenuItem
                icon={<LogOut size={14} />}
                label="Sign out"
                onClick={handleLogout}
                tone="danger"
              />
            </div>
          </div>
        )}
      </div>
    </header>
  )
}

function MenuItem({
  icon, label, onClick, tone = 'default',
}: { icon: React.ReactNode; label: string; onClick: () => void; tone?: 'default' | 'danger' }) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(
        'w-full text-left px-3 py-2 text-[13px] flex items-center gap-2.5 transition',
        tone === 'danger'
          ? 'text-danger-dark hover:bg-danger-light'
          : 'text-foreground hover:bg-surface-hover',
      )}
    >
      <span className={tone === 'danger' ? 'text-danger-dark' : 'text-foreground-muted'}>{icon}</span>
      {label}
    </button>
  )
}
