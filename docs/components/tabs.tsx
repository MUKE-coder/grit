'use client'

import { useState, type ReactNode } from 'react'
import { cn } from '@/lib/utils'

export interface TabItem {
  id: string
  label: string
  /** Optional lucide icon element or emoji shown before the label. */
  icon?: ReactNode
  content: ReactNode
}

interface TabsProps {
  items: TabItem[]
  /** id of the tab open by default (defaults to the first). */
  defaultId?: string
  className?: string
}

/**
 * A small, reusable tab switcher for docs pages — e.g. choosing which app
 * to scaffold (API / Mobile / Desktop / Everything) or per-OS instructions.
 */
export function Tabs({ items, defaultId, className }: TabsProps) {
  const [active, setActive] = useState(defaultId ?? items[0]?.id)
  const current = items.find((t) => t.id === active) ?? items[0]
  if (!current) return null

  return (
    <div className={cn('my-6', className)}>
      <div className="flex flex-wrap gap-1 border-b border-border/60 mb-5">
        {items.map((t) => (
          <button
            key={t.id}
            type="button"
            onClick={() => setActive(t.id)}
            className={cn(
              'inline-flex items-center gap-1.5 px-3.5 py-2 text-sm font-medium transition-colors border-b-2 -mb-px',
              t.id === active
                ? 'border-primary text-foreground'
                : 'border-transparent text-muted-foreground hover:text-foreground'
            )}
          >
            {t.icon}
            {t.label}
          </button>
        ))}
      </div>
      <div>{current.content}</div>
    </div>
  )
}
