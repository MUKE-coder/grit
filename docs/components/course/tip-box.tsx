import { Info, AlertTriangle, CheckCircle2, AlertCircle } from 'lucide-react'
import { cn } from '@/lib/utils'

interface Props {
  tone?: 'info' | 'warning' | 'success' | 'danger'
  children: React.ReactNode
}

/**
 * TipBox — inline callout used inside lesson content. Four tones for
 * the situation:
 *   info     — neutral guidance, "good to know"
 *   warning  — gotcha, "if you don't do this you'll regret it"
 *   success  — pro tip, "this is the elegant way"
 *   danger   — destructive / irreversible action
 */
export function TipBox({ tone = 'info', children }: Props) {
  const Icon = {
    info: Info,
    warning: AlertTriangle,
    success: CheckCircle2,
    danger: AlertCircle,
  }[tone]

  const styles = {
    info: 'border-sky-500/30 bg-sky-500/[0.05]',
    warning: 'border-amber-500/30 bg-amber-500/[0.05]',
    success: 'border-emerald-500/30 bg-emerald-500/[0.05]',
    danger: 'border-rose-500/30 bg-rose-500/[0.05]',
  }[tone]

  const iconColor = {
    info: 'text-sky-500',
    warning: 'text-amber-500',
    success: 'text-emerald-500',
    danger: 'text-rose-500',
  }[tone]

  return (
    <div className={cn('not-prose my-5 rounded-xl border p-4 flex gap-3 text-sm leading-relaxed', styles)}>
      <Icon className={cn('h-4 w-4 mt-0.5 shrink-0', iconColor)} />
      <div className="text-foreground/90 flex-1 prose-grit max-w-none">{children}</div>
    </div>
  )
}
