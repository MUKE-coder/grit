import { cn } from '@/lib/utils'

/**
 * StatusBadge — semantic colored pill used for entity status across the app.
 *
 * Color palette mapped to meaning, NOT to value. Same green for paid /
 * occupied / active / available / completed; same red for overdue / failed /
 * defaulted / vacant / repossessed. Keeps cognitive load down — colors mean
 * the same thing everywhere.
 */

type Tone = 'success' | 'warning' | 'danger' | 'info' | 'neutral' | 'accent'

const toneClasses: Record<Tone, string> = {
  success: 'bg-success-light text-success-dark',
  warning: 'bg-warning-light text-warning-dark',
  danger: 'bg-danger-light text-danger-dark',
  info: 'bg-info-light text-info-dark',
  neutral: 'bg-surface-2 text-foreground-secondary',
  accent: 'bg-accent-light text-accent-hover',
}

export function StatusBadge({
  tone = 'neutral',
  children,
  className,
  size = 'md',
}: {
  tone?: Tone
  children: React.ReactNode
  className?: string
  size?: 'sm' | 'md'
}) {
  return (
    <span
      className={cn(
        'inline-flex items-center font-medium rounded-full',
        size === 'sm' ? 'px-2 py-0.5 text-[10.5px]' : 'px-2.5 py-1 text-[11.5px]',
        toneClasses[tone],
        className,
      )}
    >
      {children}
    </span>
  )
}

/** Map domain status strings to a tone + label. Keep all known statuses here
 *  so the rest of the app doesn't have to know which color belongs to which
 *  state. Add new statuses as the domain grows. */
const STATUS_MAP: Record<string, { tone: Tone; label?: string }> = {
  // Loans
  pending: { tone: 'warning' },
  approved: { tone: 'info' },
  active: { tone: 'success' },
  completed: { tone: 'neutral' },
  defaulted: { tone: 'danger' },
  rejected: { tone: 'danger' },
  cancelled: { tone: 'neutral' },

  // Repayments
  failed: { tone: 'danger' },

  // Motorcycles
  available: { tone: 'success' },
  reserved: { tone: 'warning' },
  sold: { tone: 'neutral' },
  on_loan: { tone: 'info', label: 'On Loan' },
  repossessed: { tone: 'danger' },

  // Daily Boda
  occupied: { tone: 'info' },
  returned: { tone: 'warning' },
  in_service: { tone: 'danger', label: 'In Service' },

  // Daily Boda payments
  paid: { tone: 'success' },
  partial: { tone: 'warning' },

  // Risk levels
  low: { tone: 'success' },
  medium: { tone: 'warning' },
  high: { tone: 'danger' },
  critical: { tone: 'danger' },

  // Spares stock
  out: { tone: 'danger', label: 'Out of stock' },
  low_stock: { tone: 'warning', label: 'Low stock' },
  ok: { tone: 'success', label: 'In stock' },
}

/** Convenience: pass a domain status string and get a styled badge automatically. */
export function StatusPill({ status, size = 'md' }: { status: string; size?: 'sm' | 'md' }) {
  const meta = STATUS_MAP[status] || { tone: 'neutral' as Tone }
  const label = meta.label || status.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase())
  return (
    <StatusBadge tone={meta.tone} size={size}>
      {label}
    </StatusBadge>
  )
}
