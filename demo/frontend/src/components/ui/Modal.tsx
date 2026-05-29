import { useEffect } from 'react'
import { X } from 'lucide-react'
import { cn } from '@/lib/utils'

interface ModalProps {
  open: boolean
  onClose: () => void
  title: string
  children: React.ReactNode
}

/**
 * Responsive modal.
 * - Mobile: slides up from the bottom as a sheet (max-height 90vh).
 * - Desktop (lg+): centred dialog, max-width lg.
 *
 * Drawer.tsx uses the same pattern. Modal is for legacy callsites that
 * haven't migrated to Drawer yet (e.g. some quick confirmations).
 */
export function Modal({ open, onClose, title, children }: ModalProps) {
  useEffect(() => {
    if (open) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = ''
    }
    return () => {
      document.body.style.overflow = ''
    }
  }, [open])

  useEffect(() => {
    if (!open) return
    const handler = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose() }
    window.addEventListener('keydown', handler)
    return () => window.removeEventListener('keydown', handler)
  }, [open, onClose])

  if (!open) return null

  return (
    <div className="fixed inset-0 z-50 flex flex-col-reverse lg:flex-row lg:items-center lg:justify-center lg:p-4">
      {/* Backdrop */}
      <div className="flex-1 lg:absolute lg:inset-0 bg-black/40 lg:bg-black/50" onClick={onClose} />

      {/* Card */}
      <div
        className={cn(
          'relative bg-surface shadow-2xl flex flex-col',
          // Mobile: bottom sheet
          'w-full max-h-[90vh] rounded-t-2xl border-t border-border animate-slide-up pb-[env(safe-area-inset-bottom)]',
          // Desktop: centred dialog
          'lg:max-w-lg lg:w-full lg:max-h-[90vh] lg:rounded-2xl lg:border-0 lg:border-none lg:animate-none lg:pb-0',
        )}
      >
        {/* Drag-handle pip on mobile only */}
        <div className="lg:hidden h-1.5 w-10 rounded-full bg-border mx-auto mt-2.5 mb-1" aria-hidden />

        <div className="flex items-center justify-between px-5 lg:px-6 py-3 lg:py-4 border-b border-border shrink-0">
          <h2 className="text-[15px] lg:text-[17px] font-semibold text-foreground">{title}</h2>
          <button
            onClick={onClose}
            className="h-8 w-8 rounded-lg hover:bg-surface-hover text-foreground-muted inline-flex items-center justify-center"
            aria-label="Close"
          >
            <X size={16} />
          </button>
        </div>
        <div className="flex-1 overflow-y-auto px-5 lg:px-6 py-4 lg:py-5">
          {children}
        </div>
      </div>
    </div>
  )
}
