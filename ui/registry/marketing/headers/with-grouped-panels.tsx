'use client'

import { useEffect, useRef, useState } from 'react'
import {
  BookOpen,
  ChevronDown,
  Cpu,
  FlaskConical,
  Gauge,
  Gem,
  Menu,
  Newspaper,
  Smartphone,
  Sparkles,
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

type Item = { icon: typeof Store; name: string; body: string }

const USE_CASES: Item[] = [
  { icon: Store, name: 'Marketplace', body: 'Find and buy AI tools' },
  { icon: Cpu, name: 'API Integration', body: 'Integrate AI tools into your app' },
  { icon: Gem, name: 'Partnerships', body: 'Get help when you need it' },
  { icon: Smartphone, name: 'Mobile App', body: 'Ship to iOS and Android' },
]

const APPS: Item[] = [
  { icon: Sparkles, name: 'AI', body: 'Insights and summaries' },
  { icon: Gauge, name: 'Performance', body: 'Lightning-fast load times' },
]

const CONTENT = [
  { icon: BookOpen, name: 'Announcements' },
  { icon: FlaskConical, name: 'Resources' },
  { icon: Newspaper, name: 'Blog' },
]

const MENUS = ['Product', 'Solutions'] as const
const LINKS = ['Pricing', 'Company']

/** See the note in with-split-mega-menu: a rAF flip is what makes the transition run. */
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

function MenuItem({ icon: Icon, name, body }: Item) {
  return (
    <a
      href="#"
      className="group flex items-start gap-3.5 rounded-xl p-2.5 transition-colors duration-200 hover:bg-white dark:hover:bg-white/5"
    >
      <span className="flex size-10 shrink-0 items-center justify-center rounded-xl border border-gray-200 bg-white text-gray-700 transition-colors duration-200 group-hover:border-indigo-200 group-hover:text-indigo-600 dark:border-white/10 dark:bg-white/5 dark:text-gray-300 dark:group-hover:border-indigo-500/40 dark:group-hover:text-indigo-400">
        <Icon aria-hidden="true" className="size-[18px]" />
      </span>
      <span className="min-w-0">
        <span className="block text-sm font-semibold text-gray-900 dark:text-white">
          {name}
        </span>
        <span className="mt-0.5 block text-[13px]/[1.35] text-gray-500 dark:text-gray-400">
          {body}
        </span>
      </span>
    </a>
  )
}

export default function WithGroupedPanels({
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

  return (
    <div ref={navRef} className="relative z-50 bg-white dark:bg-gray-950">
      <nav
        aria-label="Global"
        className="mx-auto flex max-w-7xl items-center gap-8 px-6 py-4"
      >
        <div className="flex flex-1 items-center">
          <a
            href="#"
            className="flex shrink-0 items-center gap-2 text-gray-900 transition-transform duration-200 active:scale-[0.97] dark:text-white"
          >
            <GritMark className="size-7" />
            <span className="text-lg font-semibold tracking-tight">{brand}</span>
          </a>
        </div>

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

        <div className="flex flex-1 items-center justify-end gap-2">
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

      {/* Each group is its own card with its own tinted fill, separated by real
          gaps rather than divider rules. Dense navigation groups read faster when
          the grouping is a shape you can see at a glance than when it is a line. */}
      {openMenu && (
        <div className="absolute inset-x-0 top-full hidden px-6 lg:block">
          <div
            className={`mx-auto flex max-w-6xl items-stretch gap-3 transition-all duration-200 ease-out ${
              shown ? 'translate-y-0 opacity-100' : '-translate-y-1 opacity-0'
            }`}
          >
            <div className="flex-1 rounded-2xl border border-gray-200 bg-gray-50/80 p-6 shadow-[0_16px_48px_-12px_rgb(15_23_42_/_0.14)] backdrop-blur-xl dark:border-white/10 dark:bg-gray-900/90">
              <p className="mb-4 font-mono text-[11px] tracking-wider text-gray-500 uppercase">
                Use cases
              </p>
              <div className="grid grid-cols-2 gap-x-4 gap-y-1">
                {USE_CASES.map((item) => (
                  <MenuItem key={item.name} {...item} />
                ))}
              </div>
            </div>

            <div className="w-72 rounded-2xl border border-gray-200 bg-gray-50/80 p-6 shadow-[0_16px_48px_-12px_rgb(15_23_42_/_0.14)] backdrop-blur-xl dark:border-white/10 dark:bg-gray-900/90">
              <p className="mb-4 font-mono text-[11px] tracking-wider text-gray-500 uppercase">
                Apps
              </p>
              <div className="space-y-1">
                {APPS.map((item) => (
                  <MenuItem key={item.name} {...item} />
                ))}
              </div>
            </div>

            <div className="w-60 rounded-2xl border border-gray-200 bg-gray-50/80 p-6 shadow-[0_16px_48px_-12px_rgb(15_23_42_/_0.14)] backdrop-blur-xl dark:border-white/10 dark:bg-gray-900/90">
              <p className="mb-4 font-mono text-[11px] tracking-wider text-gray-500 uppercase">
                Content
              </p>
              <div className="space-y-1">
                {CONTENT.map((item) => (
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
            </div>
          </div>
        </div>
      )}

      {mobileOpen && (
        <div className="border-t border-gray-200 px-6 py-4 lg:hidden dark:border-white/10">
          {[...USE_CASES, ...APPS, ...CONTENT].map((item) => (
            <a
              key={item.name}
              href="#"
              className="flex items-center gap-3 rounded-xl px-2 py-2.5 text-sm font-medium text-gray-800 hover:bg-gray-50 dark:text-gray-200 dark:hover:bg-white/5"
            >
              <item.icon aria-hidden="true" className="size-[18px] text-gray-500" />
              {item.name}
            </a>
          ))}
          <div className="mt-3 border-t border-gray-200 pt-3 dark:border-white/10">
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
