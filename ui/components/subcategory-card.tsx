'use client'

import Link from 'next/link'
import { motion, useReducedMotion } from 'motion/react'
import { ArrowUpRight } from 'lucide-react'

/**
 * A catalogue tile.
 *
 * The hover behaviour is the micro-interaction that matters here: the card
 * lifts 2px, its shadow deepens, and the arrow slides. All three on the same
 * short soft ease, so they read as one gesture rather than three animations
 * that happen to fire together.
 */
export function SubcategoryCard({
  href,
  name,
  count,
}: {
  href: string
  name: string
  count: number
}) {
  const empty = count === 0
  const reduce = useReducedMotion()

  const inner = (
    <>
      <div
        className={`hairline relative flex aspect-[4/3] items-center justify-center overflow-hidden rounded-2xl border transition-colors duration-300 ${
          empty
            ? 'border-dashed bg-gray-50/50 dark:bg-white/[0.015]'
            : 'bg-gradient-to-b from-gray-50 to-white group-hover:border-gray-300/70 dark:from-white/[0.04] dark:to-transparent dark:group-hover:border-white/20'
        }`}
      >
        {/* Wireframe rather than a screenshot: a thumbnail goes stale the
            moment a block changes and nobody notices for months. */}
        <div className="w-1/2 space-y-2 opacity-70">
          <div className="h-1.5 w-1/3 rounded-full bg-gray-300 dark:bg-gray-700" />
          <div className="h-1.5 w-full rounded-full bg-gray-200 dark:bg-gray-800" />
          <div className="h-1.5 w-4/5 rounded-full bg-gray-200 dark:bg-gray-800" />
          {!empty && (
            <div className="mt-3.5 h-4 w-14 rounded-md bg-gradient-to-r from-indigo-500 to-violet-500 opacity-90" />
          )}
        </div>

        {!empty && (
          <span className="absolute top-3 right-3 opacity-0 transition-opacity duration-300 group-hover:opacity-100">
            <ArrowUpRight aria-hidden="true" className="size-4 text-gray-400" />
          </span>
        )}
      </div>

      <h3
        className={`display mt-3.5 text-[13.5px] transition-colors duration-200 ${
          empty
            ? 'text-gray-400 dark:text-gray-600'
            : 'text-gray-900 group-hover:text-indigo-600 dark:text-white dark:group-hover:text-indigo-400'
        }`}
      >
        {name}
      </h3>
      <p className="mt-0.5 font-mono text-[11px] text-gray-400 dark:text-gray-600">
        {empty ? 'coming soon' : `${count} block${count === 1 ? '' : 's'}`}
      </p>
    </>
  )

  if (empty) {
    return (
      <div aria-disabled className="group block cursor-default">
        {inner}
      </div>
    )
  }

  return (
    <Link href={href} className="group block">
      {reduce ? (
        <div>{inner}</div>
      ) : (
        <motion.div
          whileHover={{ y: -2 }}
          whileTap={{ scale: 0.985 }}
          transition={{ duration: 0.25, ease: [0.22, 1, 0.36, 1] }}
        >
          {inner}
        </motion.div>
      )}
    </Link>
  )
}
