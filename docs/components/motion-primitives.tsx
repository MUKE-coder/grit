'use client'

// Shared motion primitives used by the landing page.
//   - MagneticButton: a Framer Motion wrapper that subtly tracks the cursor
//     and adds a tactile spring on hover/tap.
//   - GSAPSection: registers a ScrollTrigger reveal on its children so each
//     section fades + lifts as it enters the viewport.
//   - GlowOrb: a slow-drifting blurred orb for hero / CTA backdrops.
// All three are dark-mode friendly and respect prefers-reduced-motion.

import React, { useEffect, useRef } from 'react'
import { motion, useMotionValue, useSpring, type HTMLMotionProps } from 'framer-motion'
import gsap from 'gsap'
import { ScrollTrigger } from 'gsap/ScrollTrigger'
import { cn } from '@/lib/utils'

if (typeof window !== 'undefined') {
  gsap.registerPlugin(ScrollTrigger)
}

const prefersReducedMotion = () =>
  typeof window !== 'undefined' && window.matchMedia('(prefers-reduced-motion: reduce)').matches

// ─── MagneticButton ──────────────────────────────────────────────────────
// Wraps children in a motion.div that springs toward the cursor. Use as
// the outer element of a Button — the inner Button doesn't need any
// changes. `strength` controls how far the element moves (px).
export interface MagneticButtonProps extends HTMLMotionProps<'div'> {
  strength?: number
  children: React.ReactNode
}

export function MagneticButton({
  children,
  strength = 18,
  className,
  ...rest
}: MagneticButtonProps) {
  const x = useMotionValue(0)
  const y = useMotionValue(0)
  const springX = useSpring(x, { stiffness: 200, damping: 14 })
  const springY = useSpring(y, { stiffness: 200, damping: 14 })

  function onMove(e: React.MouseEvent<HTMLDivElement>) {
    if (prefersReducedMotion()) return
    const rect = e.currentTarget.getBoundingClientRect()
    const cx = rect.left + rect.width / 2
    const cy = rect.top + rect.height / 2
    x.set(((e.clientX - cx) / rect.width) * strength * 2)
    y.set(((e.clientY - cy) / rect.height) * strength * 2)
  }
  function onLeave() {
    x.set(0)
    y.set(0)
  }

  return (
    <motion.div
      className={cn('inline-block', className)}
      style={{ x: springX, y: springY }}
      onMouseMove={onMove}
      onMouseLeave={onLeave}
      whileTap={{ scale: 0.96 }}
      {...rest}
    >
      {children}
    </motion.div>
  )
}

// ─── GSAPSection ─────────────────────────────────────────────────────────
// Reveals its children with a fade + translate on scroll. Honours
// prefers-reduced-motion by skipping the animation entirely.
export function GSAPSection({
  children,
  className,
  delay = 0,
  y = 24,
  start = 'top 85%',
}: {
  children: React.ReactNode
  className?: string
  delay?: number
  y?: number
  /** ScrollTrigger start. Defaults to "top 85%". */
  start?: string
}) {
  const ref = useRef<HTMLDivElement | null>(null)

  useEffect(() => {
    if (!ref.current || prefersReducedMotion()) return

    // Animate the section's direct children so each card / heading stagger
    // independently rather than the whole block sliding in as one.
    const targets = ref.current.querySelectorAll<HTMLElement>('[data-gsap-reveal]')
    const els = targets.length > 0 ? Array.from(targets) : [ref.current]

    const ctx = gsap.context(() => {
      gsap.from(els, {
        opacity: 0,
        y,
        duration: 0.7,
        ease: 'power3.out',
        stagger: 0.08,
        delay,
        scrollTrigger: {
          trigger: ref.current,
          start,
          toggleActions: 'play none none none',
        },
      })
    }, ref)

    return () => ctx.revert()
  }, [delay, y, start])

  return (
    <div ref={ref} className={className}>
      {children}
    </div>
  )
}

// ─── GlowOrb ────────────────────────────────────────────────────────────
// Decorative blurred orb that drifts slowly. Pure visual; no interaction.
export function GlowOrb({
  className,
  delay = 0,
  duration = 14,
}: {
  className?: string
  delay?: number
  duration?: number
}) {
  return (
    <motion.div
      aria-hidden
      className={cn('pointer-events-none absolute rounded-full blur-3xl', className)}
      animate={{
        x: [0, 20, -10, 0],
        y: [0, -15, 10, 0],
        scale: [1, 1.05, 0.98, 1],
      }}
      transition={{ duration, delay, repeat: Infinity, ease: 'easeInOut' }}
    />
  )
}

// ─── FadeIn ──────────────────────────────────────────────────────────────
// Small one-shot framer fade-in used for hero text + buttons. Triggers
// on mount, no scroll observer — saves an IntersectionObserver above
// the fold.
export function FadeIn({
  children,
  delay = 0,
  y = 14,
  className,
}: {
  children: React.ReactNode
  delay?: number
  y?: number
  className?: string
}) {
  return (
    <motion.div
      className={className}
      initial={{ opacity: 0, y }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.7, ease: [0.2, 0.7, 0.2, 1], delay }}
    >
      {children}
    </motion.div>
  )
}
