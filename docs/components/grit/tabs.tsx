'use client'

import { useCallback, useEffect, useRef, useState } from 'react'

/**
 * Tabs — the design system's segmented control, with the sliding indicator.
 *
 * The indicator is measured from the selected tab rather than animated by
 * index, because the tabs are not equal width: an index-driven transform
 * lands in the wrong place the moment one label is longer than another.
 *
 * Keyboard behaviour is the reason this is a component rather than markup.
 * `role="tablist"` promises arrow-key navigation, and a tablist that only
 * responds to clicks is worse than a row of buttons: it has told a screen
 * reader to expect something it does not do.
 */
export interface TabItem {
  id: string
  label: React.ReactNode
}

export function Tabs({
  items,
  value,
  onChange,
  variant = 'segmented',
  className = '',
  label,
}: {
  items: TabItem[]
  value: string
  onChange: (id: string) => void
  variant?: 'segmented' | 'line'
  className?: string
  /** Names the group for assistive tech. */
  label: string
}) {
  const listRef = useRef<HTMLDivElement>(null)
  const [indicator, setIndicator] = useState<{ left: number; width: number } | null>(null)

  const measure = useCallback(() => {
    const list = listRef.current
    if (!list) return
    const active = list.querySelector<HTMLElement>('[aria-selected="true"]')
    if (!active) return
    setIndicator({ left: active.offsetLeft, width: active.offsetWidth })
  }, [])

  useEffect(() => {
    measure()
  }, [measure, value, items])

  // A font swap or a container resize moves the tabs under the indicator.
  useEffect(() => {
    const list = listRef.current
    if (!list || typeof ResizeObserver === 'undefined') return
    const ro = new ResizeObserver(measure)
    ro.observe(list)
    return () => ro.disconnect()
  }, [measure])

  function onKeyDown(event: React.KeyboardEvent) {
    const i = items.findIndex((t) => t.id === value)
    if (i < 0) return
    let next = i
    if (event.key === 'ArrowRight') next = (i + 1) % items.length
    else if (event.key === 'ArrowLeft') next = (i - 1 + items.length) % items.length
    else if (event.key === 'Home') next = 0
    else if (event.key === 'End') next = items.length - 1
    else return

    event.preventDefault()
    onChange(items[next].id)
    // Move focus with selection, which is what a tablist is expected to do.
    const list = listRef.current
    list?.querySelectorAll<HTMLElement>('[role="tab"]')[next]?.focus()
  }

  return (
    <div
      ref={listRef}
      role="tablist"
      aria-label={label}
      onKeyDown={onKeyDown}
      className={`g-tabs ${variant === 'line' ? 'g-tabs--line' : ''} ${className}`.trim()}
    >
      {indicator && (
        <span
          aria-hidden="true"
          className="g-tabs__indicator"
          style={{ transform: `translateX(${indicator.left}px)`, width: indicator.width }}
        />
      )}
      {items.map((tab) => (
        <button
          key={tab.id}
          type="button"
          role="tab"
          id={`tab-${tab.id}`}
          aria-selected={tab.id === value}
          aria-controls={`panel-${tab.id}`}
          // Only the selected tab is in the tab order; the arrows reach the
          // rest. Tabbing through eight tabs to leave a tablist is the thing
          // roving focus exists to stop.
          tabIndex={tab.id === value ? 0 : -1}
          className="g-tab"
          onClick={() => onChange(tab.id)}
        >
          {tab.label}
        </button>
      ))}
    </div>
  )
}

/** The panel a tab controls. */
export function TabPanel({
  id,
  active,
  children,
}: {
  id: string
  active: boolean
  children: React.ReactNode
}) {
  return (
    <div
      role="tabpanel"
      id={`panel-${id}`}
      aria-labelledby={`tab-${id}`}
      hidden={!active}
      tabIndex={0}
    >
      {children}
    </div>
  )
}
