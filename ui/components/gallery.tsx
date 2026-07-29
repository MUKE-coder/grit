'use client'

import { useMemo, useState } from 'react'
import { Search } from 'lucide-react'
import { ComponentCard } from './component-card'

export interface GalleryItem {
  name: string
  title: string
  description: string
  category: string
}

/**
 * Client-side filtering rather than routed category pages.
 *
 * The whole catalogue is a hundred rows of metadata — a few KB — so filtering
 * in the browser is instant and keeps the iframes already mounted instead of
 * tearing them down on every category change.
 */
export function Gallery({
  items,
  categories,
  baseUrl,
}: {
  items: GalleryItem[]
  categories: { name: string; count: number }[]
  baseUrl: string
}) {
  const [category, setCategory] = useState<string>('all')
  const [query, setQuery] = useState('')

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase()
    return items.filter((item) => {
      if (category !== 'all' && item.category !== category) return false
      if (!q) return true
      return (
        item.name.includes(q) ||
        item.title.toLowerCase().includes(q) ||
        item.description.toLowerCase().includes(q)
      )
    })
  }, [items, category, query])

  return (
    <>
      {/* Controls */}
      <div className="mb-8 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div className="flex flex-wrap gap-2">
          <FilterChip
            label="All"
            count={items.length}
            active={category === 'all'}
            onClick={() => setCategory('all')}
          />
          {categories.map((c) => (
            <FilterChip
              key={c.name}
              label={c.name}
              count={c.count}
              active={category === c.name}
              onClick={() => setCategory(c.name)}
            />
          ))}
        </div>

        <div className="relative md:w-72">
          <Search
            size={14}
            className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-text-muted"
          />
          <input
            type="search"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search components"
            aria-label="Search components"
            className="h-10 w-full rounded-full border border-border bg-bg-secondary pl-9 pr-4 text-sm text-foreground placeholder:text-text-muted focus:border-accent/50 focus:outline-none"
          />
        </div>
      </div>

      {filtered.length === 0 ? (
        <p className="py-20 text-center text-sm text-text-muted">
          Nothing matches{query ? ` “${query}”` : ''}. Try a different search or
          category.
        </p>
      ) : (
        <>
          <p className="mb-5 font-mono text-xs text-text-muted">
            {filtered.length} component{filtered.length === 1 ? '' : 's'}
          </p>
          <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
            {filtered.map((item) => (
              <ComponentCard
                key={item.name}
                {...item}
                installCommand={`npx shadcn@latest add ${baseUrl}/r/${item.name}.json`}
              />
            ))}
          </div>
        </>
      )}
    </>
  )
}

function FilterChip({
  label,
  count,
  active,
  onClick,
}: {
  label: string
  count: number
  active: boolean
  onClick: () => void
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      aria-pressed={active}
      className={`inline-flex items-center gap-1.5 rounded-full border px-3.5 py-1.5 text-xs font-medium capitalize transition-colors ${
        active
          ? 'border-accent/50 bg-accent/15 text-foreground'
          : 'border-border bg-bg-secondary text-text-secondary hover:border-border hover:bg-bg-elevated hover:text-foreground'
      }`}
    >
      {label}
      <span className="font-mono text-[10px] text-text-muted">{count}</span>
    </button>
  )
}
