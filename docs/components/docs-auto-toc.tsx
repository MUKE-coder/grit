'use client'

import { useEffect, useState } from 'react'
import { usePathname } from 'next/navigation'
import { CircleAlert, FileText } from 'lucide-react'
import { cn } from '@/lib/utils'
import { siteConfig } from '@/config/site'

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

/**
 * @param variant `floating` pins the rail to the right of the viewport, which
 *   is what the 232 pages that render their own layout expect. `inline` makes
 *   it the third column of DocsShell, where it is part of the 1440px row
 *   rather than stuck to the window edge.
 */
export function DocsAutoToc({
  variant = 'floating',
}: {
  variant?: 'floating' | 'inline'
} = {}) {
  const pathname = usePathname()
  const [items, setItems] = useState<Item[]>([])
  const [activeId, setActiveId] = useState('')

  // Scan the document for content headings whenever the route changes.
  useEffect(() => {
    // Respect a page that curates its own table of contents.
    //
    // The floating instance also stands down inside DocsShell, which renders
    // a rail of its own in the third column. Both are mounted on a converted
    // page (one by app/docs/layout.tsx, one by the shell) and the marker is
    // what stops the reader seeing two.
    if (document.querySelector('[data-manual-toc]')) {
      setItems([])
      return
    }
    if (variant === 'floating' && document.querySelector('[data-docs-shell]')) {
      setItems([])
      return
    }

    const scan = () => {
      const main = document.querySelector('main')
      if (!main) return
      /*
       * The page's spine, which is not simply "every heading".
       *
       * Scoping to .prose-grit was right while a page was prose with
       * components dropped into it, and stopped being right the moment the
       * structure moved into components: on the quick start, six of the
       * eight sections a reader navigates by are <Step> titles, and the old
       * selector found one of them. A rail that lists one of eight sections
       * is worse than none, because it looks complete.
       *
       * Widening it to every h2 and h3 in <main> overshoots the other way:
       * it collects the heading of every card, callout and cross-link grid,
       * which on the same page is thirteen entries nobody navigates by.
       *
       * So: every h2, every h3 inside prose, and anything that marks itself
       * as a section with data-toc. [data-toc-ignore] opts a subtree out.
       */
      const heads = Array.from(
        main.querySelectorAll<HTMLElement>('h2, .prose-grit h3, [data-toc]'),
      ).filter((h) => !h.closest('[data-toc-ignore], header, aside, nav'))
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
        next.push({ id, label, level: h.tagName === 'H2' ? 2 : 3 })
      }
      setItems(next)
    }

    const raf = requestAnimationFrame(scan)
    return () => cancelAnimationFrame(raf)
  }, [pathname, variant])

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

  const inlineRail = variant === 'inline'
  const issueUrl =
    `${siteConfig.github}/issues/new?title=` +
    encodeURIComponent(`Docs: ${pathname}`) +
    '&body=' +
    encodeURIComponent(`Page: ${siteConfig.url}${pathname}

What is wrong:
`)

  return (
    <aside
      aria-label="On this page"
      className={cn(
        inlineRail
          ? 'sticky top-16 hidden w-[232px] flex-none py-12 pl-0 pr-6 min-[1180px]:block'
          : 'fixed right-8 top-24 hidden max-h-[calc(100vh-8rem)] w-52 overflow-y-auto min-[1340px]:block',
      )}
    >
      <p className="g-eyebrow mb-3.5">On this page</p>

      <ul className="m-0 list-none border-l border-border p-0">
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
                '-ml-px block border-l-2 py-1.5 text-[13.5px] leading-[1.4] transition-colors',
                level === 3 ? 'pl-7' : 'pl-3.5',
                activeId === id
                  ? 'border-accent font-medium text-fg'
                  : 'border-transparent text-fg-subtle hover:border-border-strong hover:text-fg-muted',
              )}
            >
              {label}
            </a>
          </li>
        ))}
      </ul>

      {/* Two actions, both of which do something. The design had a third,
          "Ask AI about this page", and the AI widget it pointed at is no
          longer on the site. */}
      {inlineRail && (
        <div className="mt-7 flex flex-col gap-2.5 border-t border-border pt-5 text-[13px]">
          <button
            type="button"
            onClick={() => window.print()}
            className="inline-flex items-center gap-2 text-left text-fg-subtle transition-colors hover:text-fg"
          >
            <FileText className="h-3.5 w-3.5 flex-none" aria-hidden="true" />
            Print or save as PDF
          </button>
          <a
            href={issueUrl}
            target="_blank"
            rel="noreferrer"
            className="inline-flex items-center gap-2 text-fg-subtle transition-colors hover:text-fg"
          >
            <CircleAlert className="h-3.5 w-3.5 flex-none" aria-hidden="true" />
            Report an issue
          </a>
        </div>
      )}
    </aside>
  )
}
