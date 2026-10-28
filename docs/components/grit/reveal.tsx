'use client'

import { useEffect, useRef, useState } from 'react'

/**
 * Reveal — fade and rise, once.
 *
 * Once is the whole specification. A section that re-animates every time it
 * scrolls back into view is a page that will not sit still to be read, and
 * it is the single most common way a scroll animation turns from a nice
 * touch into an irritation. The observer disconnects after it fires.
 *
 * The CSS carries the reduced-motion override, but this also checks the
 * query: an element that starts at opacity 0 and waits for an observer is
 * invisible content, not a disabled animation, so under reduced motion it
 * renders in place with no observer at all.
 */
export function Reveal({
  children,
  delay = 0,
  className = '',
  as: Tag = 'div',
}: {
  children: React.ReactNode
  /** Milliseconds, for staggering a row. 60 per item is the system's step. */
  delay?: number
  className?: string
  as?: 'div' | 'section' | 'li' | 'article'
}) {
  const ref = useRef<HTMLElement>(null)
  const [shown, setShown] = useState(false)

  useEffect(() => {
    const node = ref.current
    if (!node) return

    const reduced =
      typeof window.matchMedia === 'function' &&
      window.matchMedia('(prefers-reduced-motion: reduce)').matches
    if (reduced || typeof IntersectionObserver === 'undefined') {
      setShown(true)
      return
    }

    const observer = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (!entry.isIntersecting) continue
          // Disconnect before the state change: once means once.
          observer.disconnect()
          if (delay) window.setTimeout(() => setShown(true), delay)
          else setShown(true)
        }
      },
      // Fire a little before the element arrives, so the motion reads as the
      // section settling rather than as something starting late.
      { rootMargin: '0px 0px -10% 0px', threshold: 0.05 }
    )
    observer.observe(node)
    return () => observer.disconnect()
  }, [delay])

  // A polymorphic `as` and a ref do not unify: TypeScript intersects the ref
  // types of every allowed tag, and nothing satisfies all of them at once.
  // The cast is confined to this one line, and `as` is worth keeping because
  // a staggered list needs <li> rather than a div inside a <ul>.
  const Element = Tag as React.ElementType

  return (
    <Element
      ref={ref}
      className={`g-reveal ${shown ? 'is-in' : ''} ${className}`.trim()}
    >
      {children}
    </Element>
  )
}
