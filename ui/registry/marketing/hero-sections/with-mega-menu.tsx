'use client'

import { useEffect, useRef, useState } from 'react'
import {
  ChevronDown,
  Cloud,
  Cpu,
  Gauge,
  Gem,
  Rocket,
  Shield,
  Smartphone,
  Sparkles,
  Workflow,
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

const FEATURES = [
  { icon: Sparkles, name: 'AI', body: 'Generate insights and summaries' },
  { icon: Gauge, name: 'Performance', body: 'Lightning-fast load times' },
  { icon: Shield, name: 'Security', body: 'Keep your data safe and secure' },
]

const MORE_FEATURES = [
  { icon: Workflow, name: 'Automation', body: 'Automate your workflow' },
  { icon: Rocket, name: 'Scalability', body: 'Scale your application' },
  { icon: Cloud, name: 'Backup', body: 'Keep your data backed up' },
  { icon: Cpu, name: 'Compute', body: 'Run jobs at native speed' },
  { icon: Gem, name: 'Partnerships', body: 'Get help when you need it' },
  { icon: Smartphone, name: 'Mobile App', body: 'Ship to iOS and Android' },
]

export default function WithMegaMenu({
  announcement = 'Grit raises $12M series B',
  announcementCta = 'Read',
  title = 'Transform your sales with data-driven insights',
  subtitle = 'Efficiently manage your team with tools that generate the API, the admin panel and the typed client from one definition.',
  primaryLabel = 'Start building',
  secondaryLabel = 'Book a demo',
}: {
  announcement?: string
  announcementCta?: string
  title?: string
  subtitle?: string
  primaryLabel?: string
  secondaryLabel?: string
}) {
  const [openMenu, setOpenMenu] = useState<string | null>(null)
  const navRef = useRef<HTMLDivElement>(null)

  // A mega menu that only closes by clicking the trigger again is a trap on
  // touch devices and unusable by keyboard. Close on outside click and Escape.
  useEffect(() => {
    function onPointerDown(e: PointerEvent) {
      if (navRef.current && !navRef.current.contains(e.target as Node)) {
        setOpenMenu(null)
      }
    }
    function onKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') setOpenMenu(null)
    }
    document.addEventListener('pointerdown', onPointerDown)
    document.addEventListener('keydown', onKeyDown)
    return () => {
      document.removeEventListener('pointerdown', onPointerDown)
      document.removeEventListener('keydown', onKeyDown)
    }
  }, [])

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-950">
      {/* Nav */}
      <div ref={navRef} className="relative z-50">
        <nav
          aria-label="Global"
          className="mx-auto flex max-w-7xl items-center gap-6 px-6 py-4"
        >
          <a href="#" className="flex items-center gap-2 text-gray-900 dark:text-white">
            <GritMark className="size-7" />
            <span className="text-lg font-semibold tracking-tight">grit</span>
          </a>

          <div className="mx-auto hidden items-center gap-1 lg:flex">
            {['Product', 'Solutions'].map((label) => (
              <button
                key={label}
                type="button"
                onClick={() => setOpenMenu(openMenu === label ? null : label)}
                aria-expanded={openMenu === label}
                className="inline-flex items-center gap-1 rounded-lg px-3.5 py-2 text-sm font-medium text-gray-700 transition-colors hover:text-gray-900 dark:text-gray-300 dark:hover:text-white"
              >
                {label}
                <ChevronDown
                  aria-hidden="true"
                  className={`size-4 text-gray-400 transition-transform duration-200 ${
                    openMenu === label ? 'rotate-180' : ''
                  }`}
                />
              </button>
            ))}
            {['Pricing', 'Company'].map((label) => (
              <a
                key={label}
                href="#"
                className="rounded-lg px-3.5 py-2 text-sm font-medium text-gray-700 transition-colors hover:text-gray-900 dark:text-gray-300 dark:hover:text-white"
              >
                {label}
              </a>
            ))}
          </div>

          <div className="ml-auto flex items-center gap-2 lg:ml-0">
            <a
              href="#"
              className="hidden rounded-lg border border-gray-300 px-4 py-1.5 text-sm font-medium text-gray-900 transition-colors hover:bg-white sm:block dark:border-white/15 dark:text-white dark:hover:bg-white/5"
            >
              Login
            </a>
            <a
              href="#"
              className="rounded-lg bg-indigo-600 px-4 py-1.5 text-sm font-semibold text-white transition-colors hover:bg-indigo-500"
            >
              Get Started
            </a>
          </div>
        </nav>

        {/* Mega menu. Absolutely positioned so opening it never reflows the
            hero underneath — a menu that pushes the page down feels broken. */}
        {openMenu === 'Product' && (
          <div className="absolute inset-x-0 top-full px-6">
            <div className="mx-auto grid max-w-6xl grid-cols-1 gap-6 rounded-2xl border border-gray-200 bg-white p-6 shadow-xl lg:grid-cols-[1fr_1.4fr_auto] dark:border-white/10 dark:bg-gray-900">
              <div>
                <p className="mb-4 font-mono text-[11px] tracking-wider text-gray-500 uppercase">
                  Features
                </p>
                <div className="space-y-1">
                  {FEATURES.map((f) => (
                    <MenuItem key={f.name} {...f} />
                  ))}
                </div>
              </div>

              <div className="lg:border-l lg:border-gray-200 lg:pl-6 dark:lg:border-white/10">
                <p className="mb-4 font-mono text-[11px] tracking-wider text-gray-500 uppercase">
                  More features
                </p>
                <div className="grid gap-1 sm:grid-cols-2">
                  {MORE_FEATURES.map((f) => (
                    <MenuItem key={f.name} {...f} />
                  ))}
                </div>
              </div>

              {/* Promo card */}
              <a
                href="#"
                className="group hidden w-64 flex-col overflow-hidden rounded-xl border border-gray-200 bg-gray-50 transition-colors hover:border-indigo-300 lg:flex dark:border-white/10 dark:bg-white/5 dark:hover:border-indigo-500/40"
              >
                <div className="relative h-28 overflow-hidden bg-gradient-to-br from-indigo-500/15 to-sky-400/10">
                  <div className="absolute inset-x-6 top-6 rounded-t-lg border border-gray-300/60 bg-white/70 p-2 dark:border-white/10 dark:bg-white/5">
                    <div className="mb-1.5 h-1.5 w-10 rounded bg-indigo-500/60" />
                    <div className="mb-1 h-1.5 w-full rounded bg-gray-300 dark:bg-gray-700" />
                    <div className="h-1.5 w-2/3 rounded bg-gray-300 dark:bg-gray-700" />
                  </div>
                </div>
                <div className="p-4">
                  <p className="text-sm font-semibold text-gray-900 dark:text-white">
                    Multimodal learning
                  </p>
                  <p className="mt-1 text-xs/5 text-gray-600 dark:text-gray-400">
                    How the platform handles text, images and structured data together.
                  </p>
                </div>
              </a>
            </div>
          </div>
        )}

        {openMenu === 'Solutions' && (
          <div className="absolute inset-x-0 top-full px-6">
            <div className="mx-auto grid max-w-2xl gap-1 rounded-2xl border border-gray-200 bg-white p-6 shadow-xl sm:grid-cols-2 dark:border-white/10 dark:bg-gray-900">
              {[
                { icon: Rocket, name: 'Startups', body: 'Ship your first version fast' },
                { icon: Shield, name: 'Enterprise', body: 'SSO, audit logs and SLAs' },
                { icon: Workflow, name: 'Agencies', body: 'Spin up client projects' },
                { icon: Cpu, name: 'Internal tools', body: 'Admin panels in minutes' },
              ].map((f) => (
                <MenuItem key={f.name} {...f} />
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Hero */}
      <div className="mx-auto max-w-4xl px-6 pt-16 pb-14 text-center sm:pt-24">
        <a
          href="#"
          className="inline-flex items-center gap-3 rounded-full border border-gray-200 bg-white py-1 pr-2 pl-4 text-sm shadow-sm transition-colors hover:border-gray-300 dark:border-white/10 dark:bg-white/5 dark:hover:border-white/20"
        >
          <span className="text-gray-700 dark:text-gray-300">{announcement}</span>
          <span className="h-4 w-px bg-gray-200 dark:bg-white/10" />
          <span className="pr-1 font-medium text-indigo-600 dark:text-indigo-400">
            {announcementCta}
          </span>
        </a>

        <h1 className="mt-8 text-5xl font-bold tracking-tight text-balance text-gray-900 sm:text-7xl dark:text-white">
          {title}
        </h1>
        <p className="mx-auto mt-6 max-w-2xl text-lg text-gray-600 dark:text-gray-400">
          {subtitle}
        </p>

        <div className="mt-9 flex flex-wrap items-center justify-center gap-3">
          <a
            href="#"
            className="rounded-xl bg-indigo-600 px-6 py-3 text-sm font-semibold text-white transition-colors hover:bg-indigo-500"
          >
            {primaryLabel}
          </a>
          <a
            href="#"
            className="rounded-xl border border-gray-300 bg-white px-6 py-3 text-sm font-semibold text-gray-900 transition-colors hover:bg-gray-50 dark:border-white/15 dark:bg-white/5 dark:text-white dark:hover:bg-white/10"
          >
            {secondaryLabel}
          </a>
        </div>
      </div>

      {/* Product frame */}
      <div className="mx-auto max-w-6xl px-6 pb-20">
        <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-2xl dark:border-white/10 dark:bg-gray-900">
          <div className="flex">
            <div className="hidden w-52 shrink-0 border-r border-gray-200 p-4 sm:block dark:border-white/10">
              <div className="mb-5 flex items-center gap-2 rounded-lg border border-gray-200 px-2.5 py-2 dark:border-white/10">
                <GritMark className="size-5 text-gray-900 dark:text-white" />
                <span className="text-sm font-medium text-gray-900 dark:text-white">
                  Grit Pro
                </span>
              </div>
              {['Dashboard', 'Analytics', 'AI Insights', 'Predictions', 'Reports'].map(
                (item, i) => (
                  <div
                    key={item}
                    className={`mb-0.5 rounded-lg px-2.5 py-2 text-sm ${
                      i === 0
                        ? 'bg-gray-100 font-medium text-gray-900 dark:bg-white/10 dark:text-white'
                        : 'text-gray-600 dark:text-gray-400'
                    }`}
                  >
                    {item}
                  </div>
                ),
              )}
            </div>

            <div className="min-w-0 flex-1 p-5">
              <div className="mb-5 grid grid-cols-2 gap-3 lg:grid-cols-4">
                {[
                  ['Total balance', '$23,056.56', '+65%', true],
                  ['Orders', '562', '+12%', true],
                  ['Customers', '456', '-5%', false],
                  ['Recurring', '$8,420.00', '+65%', true],
                ].map(([label, value, delta, up]) => (
                  <div
                    key={label as string}
                    className="rounded-xl border border-gray-200 p-3.5 dark:border-white/10"
                  >
                    <div className="mb-1.5 flex items-center gap-1.5">
                      <span className="text-xs text-gray-500 dark:text-gray-500">
                        {label}
                      </span>
                      <span
                        className={`rounded px-1 text-[10px] font-medium ${
                          up
                            ? 'bg-green-100 text-green-700 dark:bg-green-500/15 dark:text-green-400'
                            : 'bg-rose-100 text-rose-700 dark:bg-rose-500/15 dark:text-rose-400'
                        }`}
                      >
                        {delta as string}
                      </span>
                    </div>
                    <p className="text-lg font-semibold text-gray-900 dark:text-white">
                      {value as string}
                    </p>
                  </div>
                ))}
              </div>
              <Sparkline />
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

function MenuItem({
  icon: Icon,
  name,
  body,
}: {
  icon: typeof Sparkles
  name: string
  body: string
}) {
  return (
    <a
      href="#"
      className="group flex items-start gap-3 rounded-xl p-2.5 transition-colors hover:bg-gray-50 dark:hover:bg-white/5"
    >
      <span className="mt-0.5 inline-flex size-8 shrink-0 items-center justify-center rounded-lg border border-gray-200 bg-white text-gray-700 transition-colors group-hover:border-indigo-300 group-hover:text-indigo-600 dark:border-white/10 dark:bg-white/5 dark:text-gray-300 dark:group-hover:border-indigo-500/40 dark:group-hover:text-indigo-400">
        <Icon aria-hidden="true" className="size-4" />
      </span>
      <span className="min-w-0">
        <span className="block text-sm font-medium text-gray-900 dark:text-white">
          {name}
        </span>
        <span className="block truncate text-xs text-gray-500 dark:text-gray-500">
          {body}
        </span>
      </span>
    </a>
  )
}

/** Two smoothed series. Static paths keep the block dependency-free — a chart
 *  library would be a heavier install than the entire rest of the hero. */
function Sparkline() {
  return (
    <div className="rounded-xl border border-gray-200 p-4 dark:border-white/10">
      <p className="text-sm font-medium text-gray-900 dark:text-white">Activity</p>
      <p className="mb-3 text-xs text-gray-500 dark:text-gray-500">
        Visitors and page views
      </p>
      <svg viewBox="0 0 600 140" className="h-36 w-full" preserveAspectRatio="none">
        <defs>
          <linearGradient id="grit-hero-fill" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="rgb(99 102 241)" stopOpacity="0.18" />
            <stop offset="100%" stopColor="rgb(99 102 241)" stopOpacity="0" />
          </linearGradient>
        </defs>
        <path
          d="M0 104 C 40 100, 60 84, 100 86 S 160 68, 200 74 S 260 52, 300 60 S 360 36, 400 42 S 460 22, 500 28 S 560 10, 600 14 L600 140 L0 140 Z"
          fill="url(#grit-hero-fill)"
        />
        <path
          d="M0 104 C 40 100, 60 84, 100 86 S 160 68, 200 74 S 260 52, 300 60 S 360 36, 400 42 S 460 22, 500 28 S 560 10, 600 14"
          fill="none"
          stroke="rgb(99 102 241)"
          strokeWidth="2.5"
          vectorEffect="non-scaling-stroke"
        />
        <path
          d="M0 122 C 40 120, 60 112, 100 114 S 160 104, 200 108 S 260 96, 300 100 S 360 88, 400 92 S 460 80, 500 84 S 560 72, 600 76"
          fill="none"
          className="stroke-gray-300 dark:stroke-gray-700"
          strokeWidth="2"
          vectorEffect="non-scaling-stroke"
        />
      </svg>
      <div className="mt-2 flex justify-between font-mono text-[10px] text-gray-400 dark:text-gray-600">
        {['Dec 1', 'Dec 9', 'Dec 17', 'Dec 25', 'Dec 31'].map((d) => (
          <span key={d}>{d}</span>
        ))}
      </div>
    </div>
  )
}
