'use client'

import { motion, useReducedMotion } from 'motion/react'

/**
 * Scroll-reveal wrapper.
 *
 * Deliberately restrained: 16px of travel and a short soft ease. Apple's
 * entrances are barely perceptible — you notice the content, not the animation.
 * A long slide from 40px reads as a template.
 *
 * `once` so content never re-animates when you scroll back up, and
 * useReducedMotion so anyone who has asked the OS for less motion gets none.
 *
 * Pass `immediate` on a page SHORT ENOUGH NOT TO SCROLL. The -80px viewport
 * margin shrinks the observer root, so an element sitting in that bottom band on
 * a page with no scroll never intersects and stays at opacity 0 forever. That is
 * not theoretical: it hid the whole category list on the 404 page at an 770px
 * viewport, and the DOM reported opacity 1 on the links themselves — the hidden
 * element was this wrapper.
 */
export function Reveal({
  children,
  delay = 0,
  className = '',
  immediate = false,
}: {
  children: React.ReactNode
  delay?: number
  className?: string
  /** Animate on mount instead of on scroll. For content always above the fold. */
  immediate?: boolean
}) {
  const reduce = useReducedMotion()

  if (reduce) return <div className={className}>{children}</div>

  const entrance = { opacity: 1, y: 0 }

  return (
    <motion.div
      initial={{ opacity: 0, y: 16 }}
      {...(immediate
        ? { animate: entrance }
        : { whileInView: entrance, viewport: { once: true, margin: '-80px' } })}
      transition={{ duration: 0.55, delay, ease: [0.22, 1, 0.36, 1] }}
      className={className}
    >
      {children}
    </motion.div>
  )
}
