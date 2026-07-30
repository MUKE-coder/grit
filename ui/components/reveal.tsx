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
 */
export function Reveal({
  children,
  delay = 0,
  className = '',
}: {
  children: React.ReactNode
  delay?: number
  className?: string
}) {
  const reduce = useReducedMotion()

  if (reduce) return <div className={className}>{children}</div>

  return (
    <motion.div
      initial={{ opacity: 0, y: 16 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true, margin: '-80px' }}
      transition={{ duration: 0.55, delay, ease: [0.22, 1, 0.36, 1] }}
      className={className}
    >
      {children}
    </motion.div>
  )
}
