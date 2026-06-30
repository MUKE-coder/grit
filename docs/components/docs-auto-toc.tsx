'use client'

import { useEffect, useState } from 'react'
import { usePathname } from 'next/navigation'
import { cn } from '@/lib/utils'

// Self-scanning "On this page" rail. Docs pages have no shared layout and
// rarely pass a heading list, so instead of wiring ~130 pages by hand this
// component reads the rendered DOM: it collects every <h2>/<h3> inside a
// .prose-grit block, assigns stable ids where missing, and tracks the
// active heading with an IntersectionObserver. Mounted once from
// app/docs/layout.tsx. Defers to any page that ships its own [data-manual-toc].

interface Item {
  id: string
  label: string
  level: number
}

function slugify(text: string): string {
  return text
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
}

export function DocsAutoToc() {
  const pathname = usePathname()
  const [items, setItems] = useState<Item[]>([])
  const [activeId, setActiveId] = useState('')

  // Scan the document for content headings whenever the route changes.
  useEffect(() => {
    // Respect a page that curates its own table of contents.
    if (document.querySelector('[data-manual-toc]')) {
      setItems([])
      return
    }

    const scan = () => {
      const main = document.querySelector('main')
      if (!main) return
      const heads = Array.from(
        main.querySelectorAll<HTMLElement>('.prose-grit h2, .prose-grit h3'),
      )
      const used = new Set<string>()
      const next: Item[] = []
      for (const h of heads) {
        const label = (h.textContent || '').trim()
        if (!label) continue
        let id = h.id || slugify(label)
        if (!id) continue
        while (used.has(id)) id = `${id}-1`
        used.add(id)
        if (!h.id) h.id = id
        next.push({ id, label, level: h.tagName === 'H3' ? 3 : 2 })
      }
      setItems(next)
    }

    const raf = requestAnimationFrame(scan)
    return () => cancelAnimationFrame(raf)
  }, [pathname])

  // Highlight the heading nearest the top of the viewport.
  useEffect(() => {
    if (items.length === 0) return
    const observer = new IntersectionObserver(
      (entries) => {
        const visible = entries
          .filter((e) => e.isIntersecting)
          .sort((a, b) => a.boundingClientRect.top - b.boundingClientRect.top)
        if (visible.length > 0) setActiveId(visible[0].target.id)
      },
      { rootMargin: '-80px 0px -65% 0px', threshold: 0 },
    )
    items.forEach(({ id }) => {
      const el = document.getElementById(id)
      if (el) observer.observe(el)
    })
    return () => observer.disconnect()
  }, [items])

  // Only show when there's room (content + sidebar + rail) and enough headings.
  if (items.length < 2) return null

  return (
    <nav className="hidden min-[1340px]:block fixed right-8 top-24 w-52 max-h-[calc(100vh-8rem)] overflow-y-auto">
      <p className="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground/50 mb-3">
        On this page
      </p>
      <ul className="space-y-0.5">
        {items.map(({ id, label, level }) => (
          <li key={id}>
            <a
              href={`#${id}`}
              onClick={(e) => {
                e.preventDefault()
                document.getElementById(id)?.scrollIntoView({ behavior: 'smooth' })
                setActiveId(id)
              }}
              className={cn(
                'block text-[12px] leading-relaxed py-0.5 border-l-2 transition-colors',
                level === 3 ? 'pl-6' : 'pl-3',
                activeId === id
                  ? 'border-primary text-primary font-medium'
                  : 'border-transparent text-muted-foreground/60 hover:text-muted-foreground hover:border-border',
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
