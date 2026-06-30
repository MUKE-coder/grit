'use client'

import { motion, useReducedMotion } from 'framer-motion'
import type { ReactNode } from 'react'

// Subtle scroll-into-view reveal: fade + small lift, once. Use to give each
// section a gentle entrance without being exaggerated. Honours
// prefers-reduced-motion (renders a plain div). `delay` enables staggering
// sibling blocks; pass the same className you'd put on the wrapped element so
// Reveal can BE that element (e.g. a grid container) rather than adding a layer.
export function Reveal({
  children,
  className,
  delay = 0,
  y = 18,
  as = 'div',
}: {
  children: ReactNode
  className?: string
  delay?: number
  y?: number
  as?: 'div' | 'section'
}) {
  const reduced = useReducedMotion()

  if (reduced) {
    const Tag = as
    return <Tag className={className}>{children}</Tag>
  }

  const MotionTag = as === 'section' ? motion.section : motion.div
  return (
    <MotionTag
      className={className}
      initial={{ opacity: 0, y }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true, margin: '-80px' }}
      transition={{ duration: 0.6, ease: [0.2, 0.7, 0.2, 1], delay }}
    >
      {children}
    </MotionTag>
  )
}
