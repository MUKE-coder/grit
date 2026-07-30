'use client'

import { useEffect, useState } from 'react'
import { Menu, X } from 'lucide-react'

/** The Grit UI mark, inlined so the block stays self-contained. */
function GritMark({ className = 'size-7' }: { className?: string }) {
  return (
    <svg viewBox="0 0 32 32" fill="none" aria-hidden="true" className={className}>
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M8 0H24A8 8 0 0 1 32 8V24A8 8 0 0 1 24 32H8A8 8 0 0 1 0 24V8A8 8 0 0 1 8 0ZM16 9.4A6.6 6.6 0 1 0 21.4 19.9V17.4H17.2A1.7 1.7 0 0 1 17.2 14H23.1A1.7 1.7 0 0 1 24.8 15.7V20.6A1.7 1.7 0 0 1 24.4 21.7A10 10 0 1 1 22.6 8.2A1.7 1.7 0 0 1 20.4 10.8A6.6 6.6 0 0 0 16 9.4Z"
        fill="currentColor"
      />
    </svg>
  )
}

/**
 * The plainest useful header: mark, links, two actions.
 *
 * No dropdowns on purpose. Most sites have five destinations and do not need a
 * mega menu, and a flat bar is faster to scan and impossible to get wrong. Reach
 * for one of the mega-menu headers when the navigation genuinely outgrows this.
 *
 * Sticky with a translucent backdrop, so content scrolls under it rather than
 * disappearing behind an opaque slab.
 */
export default function SimpleWithActions({
  brand = 'tailark',
  links = ['Product', 'Solutions', 'Pricing', 'Company'],
  loginLabel = 'Login',
  ctaLabel = 'Get Started',
}: {
  brand?: string
  links?: string[]
  loginLabel?: string
  ctaLabel?: string
}) {
  const [mobileOpen, setMobileOpen] = useState(false)

  useEffect(() => {
    function onKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') setMobileOpen(false)
    }
    document.addEventListener('keydown', onKeyDown)
    return () => document.removeEventListener('keydown', onKeyDown)
  }, [])

  return (
    <header className="sticky top-0 z-50 border-b border-gray-200/80 bg-white/80 backdrop-blur-xl dark:border-white/10 dark:bg-gray-950/80">
      <nav
        aria-label="Global"
        className="mx-auto flex max-w-7xl items-center gap-8 px-6 py-4"
      >
        <a
          href="#"
          className="flex shrink-0 items-center gap-2 text-gray-900 transition-transform duration-200 active:scale-[0.97] dark:text-white"
        >
          <GritMark className="size-7" />
          <span className="text-lg font-semibold tracking-tight">{brand}</span>
        </a>

        <div className="hidden items-center gap-1 md:flex">
          {links.map((label) => (
            <a
              key={label}
              href="#"
              className="rounded-xl px-3.5 py-2 text-sm font-medium text-gray-700 transition-colors duration-200 hover:bg-gray-100/70 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-white/5 dark:hover:text-white"
            >
              {label}
            </a>
          ))}
        </div>

        <div className="ml-auto flex items-center gap-2">
          <a
            href="#"
            className="hidden rounded-xl px-3 py-2 text-sm font-medium text-gray-700 transition-colors duration-200 hover:text-gray-900 sm:block dark:text-gray-300 dark:hover:text-white"
          >
            {loginLabel}
          </a>
          <a
            href="#"
            className="rounded-xl bg-gray-900 px-4 py-2 text-sm font-semibold text-white shadow-[0_1px_2px_rgb(15_23_42_/_0.2)] transition-all duration-200 hover:bg-gray-800 active:scale-[0.98] dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100"
          >
            {ctaLabel}
          </a>
          <button
            type="button"
            onClick={() => setMobileOpen(!mobileOpen)}
            aria-expanded={mobileOpen}
            aria-label="Toggle navigation"
            className="-mr-1 inline-flex size-10 items-center justify-center rounded-xl text-gray-700 transition-colors hover:bg-gray-100 md:hidden dark:text-gray-300 dark:hover:bg-white/10"
          >
            {mobileOpen ? <X className="size-5" /> : <Menu className="size-5" />}
          </button>
        </div>
      </nav>

      {mobileOpen && (
        <div className="border-t border-gray-200 px-6 py-3 md:hidden dark:border-white/10">
          {links.map((label) => (
            <a
              key={label}
              href="#"
              className="block rounded-xl px-2 py-2.5 text-sm font-medium text-gray-800 hover:bg-gray-50 dark:text-gray-200 dark:hover:bg-white/5"
            >
              {label}
            </a>
          ))}
          <a
            href="#"
            className="mt-1 block rounded-xl border-t border-gray-200 px-2 pt-3.5 pb-2.5 text-sm font-medium text-gray-800 dark:border-white/10 dark:text-gray-200"
          >
            {loginLabel}
          </a>
        </div>
      )}
    </header>
  )
}
