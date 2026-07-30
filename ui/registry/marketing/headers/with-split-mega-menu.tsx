'use client'

import { useEffect, useRef, useState } from 'react'
import {
  ArrowRight,
  BookOpen,
  ChevronDown,
  Cpu,
  FlaskConical,
  Gem,
  GraduationCap,
  Menu,
  Newspaper,
  Store,
  X,
} from 'lucide-react'

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

const USE_CASES = [
  { icon: Store, name: 'Marketplace', body: 'Find and buy AI tools' },
  { icon: GraduationCap, name: 'Guides', body: 'Learn how to use AI tools' },
  { icon: Cpu, name: 'API Integration', body: 'Integrate AI tools into your app' },
  { icon: Gem, name: 'Partnerships', body: 'Get help when you need it' },
]

const SOLUTIONS_CONTENT = [
  { icon: BookOpen, name: 'Announcements' },
  { icon: FlaskConical, name: 'Resources' },
  { icon: Newspaper, name: 'Blog' },
]

const PRODUCT_CONTENT = [
  { icon: Newspaper, name: 'Changelog' },
  { icon: BookOpen, name: 'Documentation' },
  { icon: FlaskConical, name: 'API reference' },
]

const MENUS = ['Product', 'Solutions'] as const
const LINKS = ['Pricing', 'Company']

/**
 * Entrance transition without a motion library.
 *
 * The panel is conditionally rendered, so a plain `transition` class has nothing
 * to transition from — the element arrives already in its final state. Flipping
 * the class on the next frame gives the browser one paint at the start values,
 * which is what makes the transition run.
 */
function useEntrance(open: boolean) {
  const [shown, setShown] = useState(false)
  useEffect(() => {
    if (!open) {
      setShown(false)
      return
    }
    const id = requestAnimationFrame(() => setShown(true))
    return () => cancelAnimationFrame(id)
  }, [open])
  return shown
}

export default function WithSplitMegaMenu({
  brand = 'tailark',
  loginLabel = 'Login',
  ctaLabel = 'Get Started',
}: {
  brand?: string
  loginLabel?: string
  ctaLabel?: string
}) {
  const [openMenu, setOpenMenu] = useState<string | null>(null)
  const [mobileOpen, setMobileOpen] = useState(false)
  const navRef = useRef<HTMLDivElement>(null)
  const shown = useEntrance(openMenu !== null)

  // A menu that only closes by clicking its trigger again is a trap on touch and
  // unusable by keyboard. Close on outside click and on Escape.
  useEffect(() => {
    function onPointerDown(e: PointerEvent) {
      if (navRef.current && !navRef.current.contains(e.target as Node)) {
        setOpenMenu(null)
      }
    }
    function onKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') {
        setOpenMenu(null)
        setMobileOpen(false)
      }
    }
    document.addEventListener('pointerdown', onPointerDown)
    document.addEventListener('keydown', onKeyDown)
    return () => {
      document.removeEventListener('pointerdown', onPointerDown)
      document.removeEventListener('keydown', onKeyDown)
    }
  }, [])

  const content = openMenu === 'Product' ? PRODUCT_CONTENT : SOLUTIONS_CONTENT

  return (
    <div ref={navRef} className="relative z-50 bg-white dark:bg-gray-950">
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

        <div className="hidden items-center gap-1 lg:flex">
          {MENUS.map((label) => {
            const active = openMenu === label
            return (
              <button
                key={label}
                type="button"
                onClick={() => setOpenMenu(active ? null : label)}
                aria-expanded={active}
                className={`inline-flex items-center gap-1.5 rounded-xl px-4 py-2 text-sm font-medium transition-colors duration-200 ${
                  active
                    ? 'bg-gray-100 text-gray-900 dark:bg-white/10 dark:text-white'
                    : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-white/5 dark:hover:text-white'
                }`}
              >
                {label}
                <ChevronDown
                  aria-hidden="true"
                  className={`size-4 text-gray-400 transition-transform duration-200 ${
                    active ? 'rotate-180' : ''
                  }`}
                />
              </button>
            )
          })}
          {LINKS.map((label) => (
            <a
              key={label}
              href="#"
              className="rounded-xl px-4 py-2 text-sm font-medium text-gray-700 transition-colors duration-200 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-white/5 dark:hover:text-white"
            >
              {label}
            </a>
          ))}
        </div>

        <div className="ml-auto flex items-center gap-2">
          <a
            href="#"
            className="hidden rounded-xl border border-gray-200 px-4 py-2 text-sm font-medium text-gray-900 transition-colors duration-200 hover:bg-gray-50 sm:block dark:border-white/15 dark:text-white dark:hover:bg-white/5"
          >
            {loginLabel}
          </a>
          <a
            href="#"
            className="rounded-xl bg-indigo-600 px-4 py-2 text-sm font-semibold text-white shadow-[0_1px_2px_rgb(15_23_42_/_0.18)] transition-all duration-200 hover:bg-indigo-500 active:scale-[0.98]"
          >
            {ctaLabel}
          </a>
          <button
            type="button"
            onClick={() => setMobileOpen(!mobileOpen)}
            aria-expanded={mobileOpen}
            aria-label="Toggle navigation"
            className="-mr-1 inline-flex size-10 items-center justify-center rounded-xl text-gray-700 transition-colors hover:bg-gray-100 lg:hidden dark:text-gray-300 dark:hover:bg-white/10"
          >
            {mobileOpen ? <X className="size-5" /> : <Menu className="size-5" />}
          </button>
        </div>
      </nav>

      {/* Desktop mega menu. Absolutely positioned so opening it never reflows the
          page underneath — a menu that pushes content down feels broken. */}
      {openMenu && (
        <div className="absolute inset-x-0 top-full hidden px-6 lg:block">
          {/* Aligned to the nav's own content column, not centred on the viewport.
              A centred panel under a left-aligned nav floats free of the trigger
              that opened it; sharing the container's left edge visually connects
              the two. */}
          <div className="mx-auto max-w-7xl">
            <div
              className={`grid w-[46rem] grid-cols-[1.25fr_1fr] overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-[0_16px_48px_-12px_rgb(15_23_42_/_0.18)] transition-all duration-200 ease-out dark:border-white/10 dark:bg-gray-900 ${
                shown ? 'translate-y-0 opacity-100' : '-translate-y-1 opacity-0'
              }`}
            >
              <div className="p-6">
                <p className="mb-4 font-mono text-[11px] tracking-wider text-gray-500 uppercase">
                  Use cases
                </p>
                <div className="space-y-1">
                  {USE_CASES.map((item) => (
                    <a
                      key={item.name}
                      href="#"
                      className="group flex items-start gap-3.5 rounded-xl p-2.5 transition-colors duration-200 hover:bg-gray-50 dark:hover:bg-white/5"
                    >
                      <span className="flex size-10 shrink-0 items-center justify-center rounded-xl border border-gray-200 bg-gray-50 text-gray-700 transition-colors duration-200 group-hover:border-indigo-200 group-hover:bg-white group-hover:text-indigo-600 dark:border-white/10 dark:bg-white/5 dark:text-gray-300 dark:group-hover:border-indigo-500/40 dark:group-hover:text-indigo-400">
                        <item.icon aria-hidden="true" className="size-[18px]" />
                      </span>
                      <span className="min-w-0">
                        <span className="block text-sm font-semibold text-gray-900 dark:text-white">
                          {item.name}
                        </span>
                        <span className="mt-0.5 block text-[13px]/[1.35] text-gray-500 dark:text-gray-400">
                          {item.body}
                        </span>
                      </span>
                    </a>
                  ))}
                </div>
              </div>

              {/* flex-col + mt-auto on the footer: the content column is shorter
                  than the use-cases column, and an unused gap at the bottom of a
                  panel reads as something failed to load. */}
              <div className="flex flex-col border-l border-gray-200 bg-gray-50/60 p-6 dark:border-white/10 dark:bg-white/[0.02]">
                <p className="mb-4 font-mono text-[11px] tracking-wider text-gray-500 uppercase">
                  Content
                </p>
                <div className="space-y-1">
                  {content.map((item) => (
                    <a
                      key={item.name}
                      href="#"
                      className="group flex items-center gap-3 rounded-xl p-2.5 text-sm font-semibold text-gray-900 transition-colors duration-200 hover:bg-white dark:text-white dark:hover:bg-white/5"
                    >
                      <item.icon
                        aria-hidden="true"
                        className="size-[18px] shrink-0 text-gray-500 transition-colors duration-200 group-hover:text-indigo-600 dark:text-gray-400 dark:group-hover:text-indigo-400"
                      />
                      {item.name}
                    </a>
                  ))}
                </div>

                <a
                  href="#"
                  className="group mt-auto inline-flex items-center gap-1.5 pt-6 pl-2.5 text-[13px] font-semibold text-indigo-600 dark:text-indigo-400"
                >
                  View everything
                  <ArrowRight
                    aria-hidden="true"
                    className="size-3.5 transition-transform duration-200 group-hover:translate-x-0.5"
                  />
                </a>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Mobile: one flat list. Nested accordions inside a small viewport is more
          navigation than anyone wants to operate with a thumb. */}
      {mobileOpen && (
        <div className="border-t border-gray-200 px-6 py-4 lg:hidden dark:border-white/10">
          <p className="mb-2 font-mono text-[11px] tracking-wider text-gray-500 uppercase">
            Menu
          </p>
          <div className="space-y-0.5">
            {[...USE_CASES, ...SOLUTIONS_CONTENT, ...PRODUCT_CONTENT].map((item) => (
              <a
                key={item.name}
                href="#"
                className="flex items-center gap-3 rounded-xl px-2 py-2.5 text-sm font-medium text-gray-800 hover:bg-gray-50 dark:text-gray-200 dark:hover:bg-white/5"
              >
                <item.icon aria-hidden="true" className="size-[18px] text-gray-500" />
                {item.name}
              </a>
            ))}
          </div>
          <div className="mt-4 space-y-0.5 border-t border-gray-200 pt-4 dark:border-white/10">
            {LINKS.map((label) => (
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
              className="block rounded-xl px-2 py-2.5 text-sm font-medium text-gray-800 hover:bg-gray-50 dark:text-gray-200 dark:hover:bg-white/5"
            >
              {loginLabel}
            </a>
          </div>
        </div>
      )}
    </div>
  )
}
