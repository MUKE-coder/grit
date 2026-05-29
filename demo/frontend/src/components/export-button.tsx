import { Download } from 'lucide-react'
import { cn } from '@/lib/utils'

/**
 * Small "Export" button styled to match the secondary-action pattern across
 * the app (h-8, border, surface bg). Pass an `onClick` that calls
 * `exportToExcel` from `lib/export.ts`.
 */
export function ExportButton({
  onClick,
  disabled,
  label = 'Export',
  className,
}: {
  onClick: () => void
  disabled?: boolean
  label?: string
  className?: string
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled}
      className={cn(
        'h-8 px-2.5 rounded-lg border border-border bg-surface text-foreground-secondary text-[12px] font-medium',
        'hover:bg-surface-hover hover:border-foreground-muted',
        'disabled:opacity-50 disabled:cursor-not-allowed',
        'inline-flex items-center gap-1.5 transition',
        className,
      )}
      title="Export to Excel"
    >
      <Download size={12} />
      <span className="hidden sm:inline">{label}</span>
    </button>
  )
}
