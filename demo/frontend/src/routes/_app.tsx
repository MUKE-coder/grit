import { createFileRoute, Outlet, Navigate, useLocation } from '@tanstack/react-router'
import { useEffect, useState } from 'react'
import { useAuthStore, useAuthHydrated } from '@/stores/auth.store'
import { AppSidebar } from '@/components/layout/AppSidebar'
import { TopBar } from '@/components/layout/TopBar'
import { MobileBottomNav, MobileSheetMenu } from '@/components/layout/MobileBottomNav'
import { CommandPalette, useCommandPalette } from '@/components/command-palette'

export const Route = createFileRoute('/_app')({
  component: AppLayout,
})

const SIDEBAR_KEY = 'grit-sidebar-collapsed'

function AppLayout() {
  const hydrated = useAuthHydrated()
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated())
  const location = useLocation()
  const palette = useCommandPalette()
  const [sheetOpen, setSheetOpen] = useState(false)

  const [collapsed, setCollapsed] = useState<boolean>(() => {
    if (typeof window === 'undefined') return false
    return localStorage.getItem(SIDEBAR_KEY) === '1'
  })
  useEffect(() => {
    localStorage.setItem(SIDEBAR_KEY, collapsed ? '1' : '0')
  }, [collapsed])

  if (!hydrated) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-accent flex items-center justify-center text-white font-bold text-sm animate-pulse">
            KM
          </div>
          <span className="text-foreground-muted text-sm">Loading...</span>
        </div>
      </div>
    )
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" />
  }

  // POS owns its own full-screen layout (no sidebar / topbar / bottom nav).
  const isPOS = location.pathname === '/pos'
  if (isPOS) return <Outlet />

  return (
    <div className="h-screen bg-background flex overflow-hidden">
      {/* Desktop sidebar — hidden on mobile (< lg). Bottom nav takes its place. */}
      <AppSidebar
        collapsed={collapsed}
        onToggleCollapsed={() => setCollapsed(!collapsed)}
        onOpenCommandPalette={() => palette.setOpen(true)}
      />

      <div className="flex-1 flex flex-col min-w-0">
        <TopBar />
        {/* pb-16 reserves space for the bottom nav on mobile so the last row
            of any scrolling content isn't hidden behind it. */}
        <main className="flex-1 overflow-y-auto pb-16 lg:pb-0">
          <Outlet />
        </main>
      </div>

      {/* Mobile-only navigation surfaces */}
      <MobileBottomNav onOpenSheet={() => setSheetOpen(true)} />
      <MobileSheetMenu
        open={sheetOpen}
        onOpenChange={setSheetOpen}
        onOpenCommandPalette={() => palette.setOpen(true)}
      />

      <CommandPalette open={palette.open} onOpenChange={palette.setOpen} />
    </div>
  )
}
