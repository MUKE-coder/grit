'use client'

import { useEffect, useRef, useState } from 'react'
import {
  BookOpen,
  ChevronDown,
  Cloud,
  Cpu,
  FlaskConical,
  Gauge,
  Gem,
  GraduationCap,
  Menu,
  Newspaper,
  Rocket,
  Shield,
  Smartphone,
  Sparkles,
  Store,
  Workflow,
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

const FEATURES: Item[] = [
  { icon: Sparkles, name: 'AI', body: 'Generate insights and summaries' },
  { icon: Gauge, name: 'Performance', body: 'Lightning-fast load times' },
  { icon: Shield, name: 'Security', body: 'Keep your data safe and secure' },
]

const MORE_FEATURES: Item[] = [
  { icon: Workflow, name: 'Automation', body: 'Automate your workflow' },
  { icon: Rocket, name: 'Scalability', body: 'Scale your application' },
  { icon: Cloud, name: 'Backup', body: 'Keep your data backed up' },
  { icon: Gem, name: 'Partnerships', body: 'Get help when you need it' },
  { icon: Smartphone, name: 'Mobile App', body: 'Ship to iOS and Android' },
  { icon: Cpu, name: 'Compute', body: 'Run jobs at native speed' },
]

const USE_CASES: Item[] = [
  { icon: Store, name: 'Marketplace', body: 'Find and buy AI tools' },
  { icon: Cpu, name: 'API Integration', body: 'Integrate AI tools into your app' },
  { icon: GraduationCap, name: 'Guides', body: 'Learn how to use AI tools' },
  { icon: Gem, name: 'Partnerships', body: 'Get help when you need it' },
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
      className="group flex items-start gap-3.5 rounded-xl p-2.5 transition-colors duration-200 hover:bg-gray-50 dark:hover:bg-white/5"
    >
      <span className="flex size-10 shrink-0 items-center justify-center rounded-xl border border-gray-200 bg-gray-50 text-gray-700 transition-colors duration-200 group-hover:border-indigo-200 group-hover:bg-white group-hover:text-indigo-600 dark:border-white/10 dark:bg-white/5 dark:text-gray-300 dark:group-hover:border-indigo-500/40 dark:group-hover:text-indigo-400">
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

export default function FullWidthMegaMenu({
  brand = 'tailark',
  loginLabel = 'Login',
  ctaLabel = 'Get Started',
  changelogTitle = 'Multimodal learning',
  changelogBody = 'How the platform handles text, images and structured data together.',
}: {
  brand?: string
  loginLabel?: string
  ctaLabel?: string
  changelogTitle?: string
  changelogBody?: string
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

      {/* Full-bleed panel. It spans the viewport rather than floating, so the
          separator is a dashed rule instead of a card edge — a hard border across
          the full width would read as a second navbar. */}
      {openMenu && (
        <div
          className={`absolute inset-x-0 top-full hidden border-t border-dashed border-gray-300 bg-white shadow-[0_16px_48px_-16px_rgb(15_23_42_/_0.16)] transition-all duration-200 ease-out lg:block dark:border-white/15 dark:bg-gray-900 ${
            shown ? 'translate-y-0 opacity-100' : '-translate-y-1 opacity-0'
          }`}
        >
          <div className="mx-auto max-w-7xl px-6 py-8">
            {openMenu === 'Product' ? (
              <div className="grid grid-cols-[1fr_1.6fr_auto] gap-8">
                <div>
                  <p className="mb-4 font-mono text-[11px] tracking-wider text-gray-500 uppercase">
                    Features
                  </p>
                  <div className="space-y-1">
                    {FEATURES.map((item) => (
                      <MenuItem key={item.name} {...item} />
                    ))}
                  </div>
                </div>

                <div className="border-l border-gray-200 pl-8 dark:border-white/10">
                  <p className="mb-4 font-mono text-[11px] tracking-wider text-gray-500 uppercase">
                    More features
                  </p>
                  <div className="grid grid-cols-2 gap-x-4 gap-y-1">
                    {MORE_FEATURES.map((item) => (
                      <MenuItem key={item.name} {...item} />
                    ))}
                  </div>
                </div>

                <div className="w-72 border-l border-gray-200 pl-8 dark:border-white/10">
                  <p className="mb-4 font-mono text-[11px] tracking-wider text-gray-500 uppercase">
                    Changelog
                  </p>
                  <a
                    href="#"
                    className="group block overflow-hidden rounded-2xl border border-gray-200 bg-white transition-all duration-200 hover:-translate-y-0.5 hover:border-indigo-200 hover:shadow-[0_8px_24px_-8px_rgb(15_23_42_/_0.14)] dark:border-white/10 dark:bg-white/5 dark:hover:border-indigo-500/40"
                  >
                    {/* A hinted UI, not a screenshot. A real image here means the
                        block cannot be installed without shipping an asset. */}
                    <div className="relative h-28 overflow-hidden bg-gradient-to-br from-indigo-500/20 via-sky-400/15 to-emerald-300/15">
                      <div className="absolute inset-x-8 top-6 rounded-t-xl border border-white/60 bg-white/70 p-2.5 shadow-sm backdrop-blur-sm dark:border-white/10 dark:bg-white/10">
                        <div className="mb-2 h-1.5 w-10 rounded-full bg-indigo-500/70" />
                        <div className="mb-1.5 h-1.5 w-full rounded-full bg-gray-900/10 dark:bg-white/20" />
                        <div className="h-1.5 w-2/3 rounded-full bg-gray-900/10 dark:bg-white/20" />
                      </div>
                    </div>
                    <div className="p-4">
                      <p className="text-sm font-semibold text-gray-900 dark:text-white">
                        {changelogTitle}
                      </p>
                      <p className="mt-1 line-clamp-2 text-[13px]/5 text-gray-500 dark:text-gray-400">
                        {changelogBody}
                      </p>
                    </div>
                  </a>
                </div>
              </div>
            ) : (
              <div className="grid grid-cols-[1.8fr_auto] gap-8">
                <div>
                  <p className="mb-4 font-mono text-[11px] tracking-wider text-gray-500 uppercase">
                    Use cases
                  </p>
                  <div className="grid grid-cols-2 gap-x-4 gap-y-1">
                    {USE_CASES.map((item) => (
                      <MenuItem key={item.name} {...item} />
                    ))}
                  </div>
                </div>

                <div className="w-64 border-l border-gray-200 pl-8 dark:border-white/10">
                  <p className="mb-4 font-mono text-[11px] tracking-wider text-gray-500 uppercase">
                    Content
                  </p>
                  <div className="space-y-1">
                    {CONTENT.map((item) => (
                      <a
                        key={item.name}
                        href="#"
                        className="group flex items-center gap-3 rounded-xl p-2.5 text-sm font-semibold text-gray-900 transition-colors duration-200 hover:bg-gray-50 dark:text-white dark:hover:bg-white/5"
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
            )}
          </div>
        </div>
      )}

      {mobileOpen && (
        <div className="border-t border-gray-200 px-6 py-4 lg:hidden dark:border-white/10">
          <p className="mb-2 font-mono text-[11px] tracking-wider text-gray-500 uppercase">
            Menu
          </p>
          {[...FEATURES, ...MORE_FEATURES, ...USE_CASES].map((item, i) => (
            <a
              // Index in the key: these three lists are concatenated and
              // "Partnerships" appears in two of them, so the name alone is not
              // unique here even though it is unique within each list.
              key={`${item.name}-${i}`}
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
