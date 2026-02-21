'use client'

import { useState, useEffect, useCallback } from 'react'
import { cn } from '@/lib/utils'

interface TocItem {
  id: string
  label: string
}

interface TableOfContentsProps {
  items: TocItem[]
}

export function TableOfContents({ items }: TableOfContentsProps) {
  const [activeId, setActiveId] = useState<string>('')

  const handleObserver = useCallback((entries: IntersectionObserverEntry[]) => {
    // Find the first visible entry (closest to top of viewport)
    const visible = entries
      .filter((e) => e.isIntersecting)
      .sort((a, b) => a.boundingClientRect.top - b.boundingClientRect.top)

    if (visible.length > 0) {
      setActiveId(visible[0].target.id)
    }
  }, [])

  useEffect(() => {
    const observer = new IntersectionObserver(handleObserver, {
      rootMargin: '-80px 0px -60% 0px',
      threshold: 0,
    })

    items.forEach(({ id }) => {
      const el = document.getElementById(id)
      if (el) observer.observe(el)
    })

    return () => observer.disconnect()
  }, [items, handleObserver])

  return (
    <nav className="hidden xl:block fixed right-6 top-20 w-56 max-h-[calc(100vh-6rem)] overflow-y-auto">
      <p className="text-xs font-semibold uppercase tracking-wider text-muted-foreground/50 mb-3">
        On This Page
      </p>
      <ul className="space-y-0.5">
        {items.map(({ id, label }) => (
          <li key={id}>
            <a
              href={`#${id}`}
              onClick={(e) => {
                e.preventDefault()
                document.getElementById(id)?.scrollIntoView({ behavior: 'smooth' })
                setActiveId(id)
              }}
              className={cn(
                'block text-[12px] leading-relaxed py-0.5 border-l-2 pl-3 transition-colors',
                activeId === id
                  ? 'border-primary text-primary font-medium'
                  : 'border-transparent text-muted-foreground/60 hover:text-muted-foreground hover:border-border'
              )}
            >
              {label}
            </a>
          </li>
        ))}
      </ul>
    </nav>
  )
}
