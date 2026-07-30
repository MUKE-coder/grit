'use client'

import { useState } from 'react'
import { BarChart3, Globe, Search, Sparkles } from 'lucide-react'

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

const TABS = [
  { key: 'tasks', label: 'Task Management', icon: Globe },
  { key: 'analytics', label: 'Analytics', icon: BarChart3 },
  { key: 'copilot', label: 'AI Copilot', icon: Sparkles },
] as const

type TabKey = (typeof TABS)[number]['key']

const ISSUES = [
  { id: 'GRIT-525', title: 'Add ownership guard to invoice handler', tag: 'Security', area: 'API', date: 'Apr 03' },
  { id: 'GRIT-513', title: 'Refactor resource generator import blocks', tag: 'Refactor', area: 'Codegen', date: 'Mar 22' },
  { id: 'GRIT-505', title: 'Improve admin table keyboard navigation', tag: 'A11y', area: 'Admin', date: 'Mar 14' },
  { id: 'GRIT-415', title: 'Design modal system with focus trapping', tag: 'Design', area: 'Grit UI', date: 'Mar 09' },
  { id: 'GRIT-501', title: 'Make the CLI report the path it wrote', tag: 'Bug', area: 'CLI', date: 'Mar 10' },
]

export default function SplitWithTabs({
  title = 'Full-stack tooling for the AI era',
  subtitle = 'Generate the API, the admin panel and the typed client from one definition. Open source, and it works.',
  primaryLabel = 'Get Started',
  secondaryLabel = 'Watch Demo',
}: {
  title?: string
  subtitle?: string
  primaryLabel?: string
  secondaryLabel?: string
}) {
  const [tab, setTab] = useState<TabKey>('tasks')

  return (
    <div className="min-h-screen bg-white dark:bg-gray-950">
      {/* Nav */}
      <nav
        aria-label="Global"
        className="mx-auto flex max-w-7xl items-center justify-between border-b border-gray-200 px-6 py-4 dark:border-white/10"
      >
        <a href="#" className="flex items-center gap-2 text-gray-900 dark:text-white">
          <GritMark className="size-7" />
          <span className="text-lg font-semibold tracking-tight">grit</span>
        </a>
        <div className="hidden gap-8 lg:flex">
          {['Product', 'Solutions', 'Pricing', 'Company'].map((item) => (
            <a
              key={item}
              href="#"
              className="text-sm font-medium text-gray-700 transition-colors hover:text-gray-900 dark:text-gray-300 dark:hover:text-white"
            >
              {item}
            </a>
          ))}
        </div>
        <div className="flex items-center gap-2">
          <a
            href="#"
            className="hidden rounded-lg border border-gray-300 px-4 py-1.5 text-sm font-medium text-gray-900 transition-colors hover:bg-gray-50 sm:block dark:border-white/15 dark:text-white dark:hover:bg-white/5"
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

      {/* Split hero */}
      <div className="mx-auto grid max-w-7xl grid-cols-1 items-center gap-10 px-6 py-20 lg:grid-cols-[1.15fr_1fr] lg:py-28">
        <h1 className="text-5xl font-bold tracking-tight text-balance text-gray-900 sm:text-7xl dark:text-white">
          {title}
        </h1>
        <div className="lg:pt-3">
          <p className="max-w-md text-lg/8 text-gray-600 dark:text-gray-400">{subtitle}</p>
          <div className="mt-8 flex flex-wrap gap-3">
            <a
              href="#"
              className="rounded-lg bg-indigo-600 px-5 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-indigo-500"
            >
              {primaryLabel}
            </a>
            <a
              href="#"
              className="rounded-lg border border-gray-300 bg-white px-5 py-2.5 text-sm font-semibold text-gray-900 transition-colors hover:bg-gray-50 dark:border-white/15 dark:bg-white/5 dark:text-white dark:hover:bg-white/10"
            >
              {secondaryLabel}
            </a>
          </div>
        </div>
      </div>

      {/* Tab strip, sitting on a hatched rail as in the reference */}
      <div className="relative border-y border-gray-200 dark:border-white/10">
        <div
          aria-hidden="true"
          className="pointer-events-none absolute inset-0"
          style={{
            backgroundImage:
              'repeating-linear-gradient(45deg, rgb(0 0 0 / 0.035) 0, rgb(0 0 0 / 0.035) 1px, transparent 0, transparent 50%)',
            backgroundSize: '10px 10px',
          }}
        />
        <div
          role="tablist"
          aria-label="Product areas"
          className="relative mx-auto flex max-w-3xl divide-x divide-gray-200 dark:divide-white/10"
        >
          {TABS.map(({ key, label, icon: Icon }) => {
            const active = tab === key
            return (
              <button
                key={key}
                role="tab"
                aria-selected={active}
                onClick={() => setTab(key)}
                className={`flex flex-1 items-center justify-center gap-2 px-4 py-5 text-sm font-medium transition-colors ${
                  active
                    ? 'bg-white text-gray-900 dark:bg-gray-950 dark:text-white'
                    : 'text-gray-500 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white'
                }`}
              >
                <Icon aria-hidden="true" className="size-4" />
                <span className="hidden sm:inline">{label}</span>
              </button>
            )
          })}
        </div>
      </div>

      {/* App mock */}
      <div className="mx-auto max-w-6xl px-6 py-14">
        <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-2xl dark:border-white/10 dark:bg-gray-900">
          <div className="flex">
            {/* Sidebar */}
            <div className="hidden w-52 shrink-0 border-r border-gray-200 p-4 md:block dark:border-white/10">
              <div className="mb-5 flex items-center gap-2">
                <span className="inline-flex size-6 items-center justify-center rounded-md bg-indigo-600 text-[10px] font-bold text-white">
                  GR
                </span>
                <span className="text-sm font-medium text-gray-900 dark:text-white">
                  grit-ui
                </span>
              </div>
              {['Inbox', 'My issues'].map((item) => (
                <div
                  key={item}
                  className="mb-0.5 rounded-md px-2 py-1.5 text-sm text-gray-600 dark:text-gray-400"
                >
                  {item}
                </div>
              ))}
              <p className="mt-5 mb-2 font-mono text-[10px] tracking-wider text-gray-400 uppercase dark:text-gray-600">
                Workspace
              </p>
              {['Teams', 'Projects', 'Members'].map((item) => (
                <div
                  key={item}
                  className="mb-0.5 rounded-md px-2 py-1.5 text-sm text-gray-600 dark:text-gray-400"
                >
                  {item}
                </div>
              ))}
            </div>

            {/* Body */}
            <div className="min-w-0 flex-1">
              <div className="flex items-center gap-3 border-b border-gray-200 px-4 py-3 dark:border-white/10">
                <Search aria-hidden="true" className="size-4 text-gray-400" />
                <span className="text-sm text-gray-400">Filter</span>
                <span className="ml-auto rounded-md border border-gray-200 px-2.5 py-1 text-xs text-gray-600 dark:border-white/10 dark:text-gray-400">
                  Display
                </span>
              </div>

              {tab === 'analytics' ? (
                <div className="p-6">
                  <div className="mb-4 grid grid-cols-3 gap-3">
                    {[
                      ['Open', '18'],
                      ['In review', '6'],
                      ['Shipped', '124'],
                    ].map(([l, v]) => (
                      <div
                        key={l}
                        className="rounded-xl border border-gray-200 p-3.5 dark:border-white/10"
                      >
                        <p className="text-xs text-gray-500 dark:text-gray-500">{l}</p>
                        <p className="mt-1 text-xl font-semibold text-gray-900 dark:text-white">
                          {v}
                        </p>
                      </div>
                    ))}
                  </div>
                  <svg viewBox="0 0 600 120" className="h-32 w-full" preserveAspectRatio="none">
                    <path
                      d="M0 96 C 50 92, 80 70, 130 76 S 210 52, 260 60 S 340 30, 400 38 S 500 14, 600 20"
                      fill="none"
                      stroke="rgb(99 102 241)"
                      strokeWidth="2.5"
                      vectorEffect="non-scaling-stroke"
                    />
                  </svg>
                </div>
              ) : tab === 'copilot' ? (
                <div className="space-y-4 p-6">
                  <div className="flex justify-end">
                    <div className="max-w-[70%] rounded-2xl rounded-br-md bg-gray-100 px-4 py-2.5 text-sm text-gray-800 dark:bg-white/10 dark:text-gray-200">
                      Group these issues by area and flag anything security-related.
                    </div>
                  </div>
                  <div className="flex gap-2.5">
                    <span className="mt-0.5 inline-flex size-6 shrink-0 items-center justify-center rounded-full bg-indigo-600/10 text-indigo-600 dark:bg-indigo-500/15 dark:text-indigo-400">
                      <Sparkles aria-hidden="true" className="size-3" />
                    </span>
                    <p className="text-sm/6 text-gray-700 dark:text-gray-300">
                      One security item: <strong>GRIT-525</strong> adds an ownership guard
                      to the invoice handler. The rest split across Codegen (1), Admin (1),
                      Grit UI (1) and CLI (1).
                    </p>
                  </div>
                </div>
              ) : (
                <div>
                  <div className="flex items-center gap-2 bg-amber-50/60 px-4 py-2 dark:bg-amber-500/5">
                    <span className="size-1.5 rounded-full bg-amber-500" />
                    <span className="text-xs font-medium text-gray-700 dark:text-gray-300">
                      In Progress
                    </span>
                    <span className="text-xs text-gray-400">{ISSUES.length}</span>
                  </div>
                  {ISSUES.map((issue) => (
                    <div
                      key={issue.id}
                      className="flex items-center gap-3 border-t border-gray-100 px-4 py-2.5 first:border-t-0 hover:bg-gray-50 dark:border-white/5 dark:hover:bg-white/5"
                    >
                      <span className="hidden font-mono text-[11px] text-gray-400 sm:block">
                        {issue.id}
                      </span>
                      <span className="min-w-0 flex-1 truncate text-sm text-gray-900 dark:text-white">
                        {issue.title}
                      </span>
                      <span className="hidden shrink-0 rounded-md border border-gray-200 px-2 py-0.5 text-[11px] text-gray-600 lg:block dark:border-white/10 dark:text-gray-400">
                        {issue.tag}
                      </span>
                      <span className="hidden shrink-0 rounded-md border border-gray-200 px-2 py-0.5 text-[11px] text-gray-600 lg:block dark:border-white/10 dark:text-gray-400">
                        {issue.area}
                      </span>
                      <span className="hidden shrink-0 font-mono text-[11px] text-gray-400 sm:block">
                        {issue.date}
                      </span>
                      <span className="size-5 shrink-0 rounded-full bg-gradient-to-br from-indigo-400 to-fuchsia-500" />
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
