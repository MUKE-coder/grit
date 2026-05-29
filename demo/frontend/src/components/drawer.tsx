import { useEffect } from 'react'
import { X } from 'lucide-react'
import { cn } from '@/lib/utils'

/**
 * Responsive overlay container used for create + edit forms.
 *
 * - Mobile (< lg): slides up from the bottom, max-height ~90vh, rounded
 *   top corners. Native bottom-sheet feel — the user can't lose context
 *   above the sheet.
 * - Desktop (lg+): slides in from the right, fixed width set by `width`.
 *
 * Esc closes. Click outside (backdrop) closes. Background scroll is
 * locked while the sheet is open so iOS doesn't bounce the page behind.
 */
export function Drawer({
  open,
  onOpenChange,
  title,
  children,
  width = 480,
  description,
}: {
  open: boolean
  onOpenChange: (o: boolean) => void
  title: string
  description?: string
  children: React.ReactNode
  width?: number
}) {
  useEffect(() => {
    if (!open) return
    const handler = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onOpenChange(false)
    }
    window.addEventListener('keydown', handler)
    return () => window.removeEventListener('keydown', handler)
  }, [open, onOpenChange])

  useEffect(() => {
    if (!open) return
    const prev = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    return () => { document.body.style.overflow = prev }
  }, [open])

  if (!open) return null

  return (
    <div className="fixed inset-0 z-50 flex flex-col-reverse lg:flex-row">
      {/* Backdrop */}
      <div
        className="flex-1 bg-black/40 lg:backdrop-blur-sm"
        onClick={() => onOpenChange(false)}
        aria-label="Close drawer"
      />

      {/* Sheet card. The CSS variable lets the desktop width be passed in
          via inline style without breaking the mobile w-full layout. */}
      <div
        className={cn(
          'bg-surface shadow-2xl flex flex-col shrink-0',
          // Mobile: bottom sheet
          'w-full max-h-[90vh] rounded-t-2xl border-t border-border animate-slide-up pb-[env(safe-area-inset-bottom)]',
          // Desktop: right slide-over
          'lg:w-[var(--drawer-w)] lg:max-h-none lg:h-full lg:rounded-none lg:border-t-0 lg:border-l lg:animate-slide-in lg:pb-0',
        )}
        style={{ ['--drawer-w' as any]: `${width}px` }}
      >
        {/* Drag-handle pip — mobile only */}
        <div className="lg:hidden h-1.5 w-10 rounded-full bg-border mx-auto mt-2.5 mb-1.5" aria-hidden />

        <div className="flex items-start justify-between px-5 py-3 lg:py-4 border-b border-border shrink-0">
          <div className="min-w-0">
            <h2 className="text-[15px] font-semibold text-foreground truncate">{title}</h2>
            {description && (
              <p className="text-[12.5px] text-foreground-muted mt-0.5">{description}</p>
            )}
          </div>
          <button
            type="button"
            onClick={() => onOpenChange(false)}
            className="h-8 w-8 rounded-lg hover:bg-surface-hover flex items-center justify-center text-foreground-muted shrink-0"
            aria-label="Close"
          >
            <X className="h-4 w-4" />
          </button>
        </div>

        <div className="flex-1 overflow-y-auto p-5">{children}</div>
      </div>
    </div>
  )
}
