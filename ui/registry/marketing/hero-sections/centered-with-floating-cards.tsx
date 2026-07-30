import { Brain, Cpu, ShieldCheck, Users } from 'lucide-react'

/** The Grit UI mark, inlined so the block stays self-contained. */
function GritMark({ className = 'size-6' }: { className?: string }) {
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

export default function CenteredWithFloatingCards({
  title = 'Trusted third-party request handling',
  subtitle = 'Grit automates third-party support requests, speeding up responses while keeping your processes secure.',
  primaryLabel = 'Get started',
  secondaryLabel = 'View live demo',
}: {
  title?: string
  subtitle?: string
  primaryLabel?: string
  secondaryLabel?: string
}) {
  return (
    <div className="relative isolate min-h-screen overflow-hidden bg-slate-50 dark:bg-gray-950">
      {/* Dotted connector rails. Decorative, so they sit behind everything and
          are hidden from assistive tech. */}
      <svg
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 -z-10 size-full text-slate-300 dark:text-white/10"
      >
        <line x1="12%" y1="0" x2="12%" y2="100%" stroke="currentColor" strokeDasharray="4 6" />
        <line x1="88%" y1="0" x2="88%" y2="100%" stroke="currentColor" strokeDasharray="4 6" />
        <line x1="0" y1="46%" x2="100%" y2="46%" stroke="currentColor" strokeDasharray="4 6" />
      </svg>

      {/* Pill nav */}
      <div className="px-6 pt-6">
        <nav
          aria-label="Global"
          className="mx-auto flex max-w-3xl items-center gap-6 rounded-full border border-slate-200 bg-white/90 px-4 py-2 shadow-sm backdrop-blur dark:border-white/10 dark:bg-white/5"
        >
          <a href="#" className="flex items-center gap-2 text-gray-900 dark:text-white">
            <GritMark className="size-6" />
            <span className="text-sm font-semibold tracking-tight">Grit</span>
          </a>
          <div className="hidden gap-5 sm:flex">
            {['Features', 'Pricing', 'About', 'FAQ'].map((item) => (
              <a
                key={item}
                href="#"
                className="text-sm text-gray-600 transition-colors hover:text-gray-900 dark:text-gray-400 dark:hover:text-white"
              >
                {item}
              </a>
            ))}
          </div>
          <div className="ml-auto flex items-center gap-3">
            <a
              href="#"
              className="text-sm text-gray-600 transition-colors hover:text-gray-900 dark:text-gray-400 dark:hover:text-white"
            >
              Sign in
            </a>
            <a
              href="#"
              className="rounded-full bg-blue-600 px-4 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-500"
            >
              Book a demo
            </a>
          </div>
        </nav>
      </div>

      {/* Hero */}
      <div className="relative mx-auto max-w-3xl px-6 pt-20 pb-10 text-center">
        <h1 className="text-4xl font-bold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
          {title}
        </h1>
        <p className="mx-auto mt-6 max-w-xl text-base/7 text-gray-600 dark:text-gray-400">
          {subtitle}
        </p>
        <div className="mt-8 flex flex-wrap items-center justify-center gap-3">
          <a
            href="#"
            className="rounded-full bg-blue-600 px-6 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-blue-500"
          >
            {primaryLabel}
          </a>
          <a
            href="#"
            className="rounded-full border border-slate-200 bg-white px-6 py-2.5 text-sm font-semibold text-gray-900 transition-colors hover:bg-slate-50 dark:border-white/10 dark:bg-white/5 dark:text-white dark:hover:bg-white/10"
          >
            {secondaryLabel}
          </a>
        </div>
      </div>

      {/* Side cards — absolute on large screens where there is room beside the
          headline, stacked into the normal flow below it otherwise. */}
      <div className="mx-auto max-w-7xl px-6 pb-20">
        <div className="grid gap-4 lg:hidden">
          <ComparisonCard />
          <ProxyCard />
        </div>

        <div className="pointer-events-none absolute inset-x-0 top-40 hidden lg:block">
          <div className="mx-auto flex max-w-[86rem] justify-between px-6">
            <div className="pointer-events-auto w-64">
              <ComparisonCard />
            </div>
            <div className="pointer-events-auto w-64">
              <ProxyCard />
            </div>
          </div>
        </div>

        {/* Flow diagram. Capped and centred: an aspect-square grid inside a
            free-growing column turns into enormous tiles on a wide screen. */}
        <div className="mx-auto mt-10 grid max-w-4xl justify-items-center gap-6 lg:mt-16 lg:grid-cols-[auto_auto_auto] lg:items-center lg:justify-between">
          <div className="w-64 rounded-2xl border border-slate-200 bg-white p-4 shadow-sm dark:border-white/10 dark:bg-gray-900">
            <p className="mb-3 text-sm font-medium text-gray-900 dark:text-white">
              Your customers
            </p>
            <div className="grid grid-cols-3 gap-2">
              {[
                'from-amber-200 to-orange-300',
                null,
                'from-sky-200 to-blue-300',
                'from-emerald-200 to-teal-300',
                'from-rose-200 to-pink-300',
                null,
                null,
                'from-violet-200 to-purple-300',
                null,
              ].map((tone, i) => (
                <span
                  key={i}
                  className={`size-14 rounded-lg ${
                    tone
                      ? `bg-gradient-to-br ${tone}`
                      : 'border border-slate-100 dark:border-white/5'
                  }`}
                />
              ))}
            </div>
          </div>

          <div className="flex items-center justify-center">
            <span className="inline-flex items-center gap-2 rounded-full bg-blue-600 px-5 py-2.5 text-sm font-semibold text-white shadow-lg shadow-blue-600/25">
              <GritMark className="size-4" />
              Grit
            </span>
          </div>

          <div className="w-64 rounded-2xl border border-slate-200 bg-white p-4 shadow-sm dark:border-white/10 dark:bg-gray-900">
            <p className="mb-3 flex items-center gap-1.5 text-sm font-medium text-gray-900 dark:text-white">
              <Users aria-hidden="true" className="size-3.5" />
              Your support team
            </p>
            <div className="grid grid-cols-3 gap-3">
              {[
                { icon: Cpu, label: 'Triage' },
                { icon: Brain, label: 'Assist' },
                { icon: ShieldCheck, label: 'Verify' },
              ].map(({ icon: Icon, label }) => (
                <div key={label} className="text-center">
                  <span className="mx-auto mb-1.5 flex size-10 items-center justify-center rounded-full bg-slate-100 text-gray-700 dark:bg-white/10 dark:text-gray-300">
                    <Icon aria-hidden="true" className="size-4" />
                  </span>
                  <p className="text-[11px] text-gray-500 dark:text-gray-400">{label}</p>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

function ComparisonCard() {
  return (
    <div className="rounded-xl border border-slate-200 bg-white p-3.5 shadow-sm dark:border-white/10 dark:bg-gray-900">
      {[
        { pct: '8', label: 'Without Grit', on: false },
        { pct: '75', label: 'With Grit', on: true },
      ].map((row) => (
        <div key={row.label} className="flex items-center gap-3 py-1.5">
          <span className="w-8 text-right text-sm font-semibold text-gray-900 dark:text-white">
            {row.pct}
            <span className="text-[10px] text-gray-400">%</span>
          </span>
          <span className="flex-1 text-xs text-gray-600 dark:text-gray-400">
            {row.label}
          </span>
          <span
            className={`flex h-4 w-7 shrink-0 items-center rounded-full px-0.5 ${
              row.on ? 'justify-end bg-green-500' : 'justify-start bg-slate-200 dark:bg-white/15'
            }`}
          >
            <span className="size-3 rounded-full bg-white" />
          </span>
        </div>
      ))}
    </div>
  )
}

function ProxyCard() {
  return (
    <div className="rounded-xl border border-slate-200 bg-white p-4 text-center shadow-sm dark:border-white/10 dark:bg-gray-900">
      <p className="mb-3 flex items-center justify-center gap-1.5 text-xs text-gray-600 dark:text-gray-400">
        <ShieldCheck aria-hidden="true" className="size-3.5 text-green-600" />
        Secured proxy protection
      </p>
      <div className="relative mx-auto flex size-20 items-center justify-center rounded-xl bg-gradient-to-br from-green-50 to-emerald-100 dark:from-green-500/10 dark:to-emerald-500/10">
        <span className="flex size-10 items-center justify-center rounded-full bg-green-600 text-white">
          <ShieldCheck aria-hidden="true" className="size-5" />
        </span>
      </div>
    </div>
  )
}
