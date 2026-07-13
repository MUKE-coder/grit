import { Info, Lightbulb, TriangleAlert, OctagonAlert, DoorOpen } from 'lucide-react'
import type { ReactNode } from 'react'

// Typed admonition boxes shared across the docs. Colour signals meaning,
// never decoration. The "escape-hatch" variant (violet) is a Grit signature:
// use it wherever an opinionated default can be overridden.
type CalloutType = 'note' | 'tip' | 'warning' | 'danger' | 'escape'

const STYLES: Record<
  CalloutType,
  { border: string; bg: string; icon: string; label: string; Icon: typeof Info }
> = {
  note: {
    border: 'border-blue-500/25',
    bg: 'bg-blue-500/[0.06]',
    icon: 'text-blue-400',
    label: 'Note',
    Icon: Info,
  },
  tip: {
    border: 'border-emerald-500/25',
    bg: 'bg-emerald-500/[0.06]',
    icon: 'text-emerald-400',
    label: 'Tip',
    Icon: Lightbulb,
  },
  warning: {
    border: 'border-amber-500/25',
    bg: 'bg-amber-500/[0.06]',
    icon: 'text-amber-400',
    label: 'Warning',
    Icon: TriangleAlert,
  },
  danger: {
    border: 'border-red-500/25',
    bg: 'bg-red-500/[0.06]',
    icon: 'text-red-400',
    label: 'Danger',
    Icon: OctagonAlert,
  },
  escape: {
    border: 'border-violet-500/30',
    bg: 'bg-violet-500/[0.07]',
    icon: 'text-violet-400',
    label: 'Escape hatch',
    Icon: DoorOpen,
  },
}

interface CalloutProps {
  type?: CalloutType
  /** Overrides the default label ("Note", "Escape hatch", …). */
  title?: string
  children: ReactNode
  className?: string
}

export function Callout({ type = 'note', title, children, className = '' }: CalloutProps) {
  const s = STYLES[type]
  const Icon = s.Icon
  return (
    <div
      className={`my-6 flex gap-3 rounded-lg border ${s.border} ${s.bg} px-4 py-3.5 ${className}`}
    >
      <Icon className={`h-[18px] w-[18px] shrink-0 mt-0.5 ${s.icon}`} aria-hidden />
      <div className="min-w-0 text-sm leading-relaxed text-muted-foreground [&_a]:text-primary [&_a:hover]:underline [&_code]:text-foreground/90 [&_code]:text-[13px]">
        <span className={`mr-2 font-semibold ${s.icon}`}>{title ?? s.label}</span>
        {children}
      </div>
    </div>
  )
}

/** Convenience wrapper — the recurring "you can override this" box. */
export function EscapeHatch({ title, children, className }: Omit<CalloutProps, 'type'>) {
  return (
    <Callout type="escape" title={title} className={className}>
      {children}
    </Callout>
  )
}
