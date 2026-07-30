import { ArrowRight } from 'lucide-react'

const LOGOS = ['Northwind', 'Acme', 'Globex', 'Initech', 'Umbrella']

export default function WithAppScreenshot({
  badge = 'New',
  badgeText = 'Admin panel generator is live',
  title = 'The admin panel you never have to build',
  subtitle = 'Define a resource once and get a filterable table, a validated form and a detail view — generated into your repo as code you own and can edit.',
  primaryLabel = 'Get started',
  secondaryLabel = 'See a live demo',
}: {
  badge?: string
  badgeText?: string
  title?: string
  subtitle?: string
  primaryLabel?: string
  secondaryLabel?: string
}) {
  return (
    <div className="relative isolate overflow-hidden bg-white dark:bg-gray-900">
      {/* Faint grid backdrop */}
      <svg
        aria-hidden="true"
        className="absolute inset-0 -z-10 size-full stroke-gray-200 [mask-image:radial-gradient(64rem_64rem_at_top,white,transparent)] dark:stroke-white/10"
      >
        <defs>
          <pattern
            id="hero-grid"
            width={200}
            height={200}
            x="50%"
            y={-1}
            patternUnits="userSpaceOnUse"
          >
            <path d="M.5 200V.5H200" fill="none" />
          </pattern>
        </defs>
        <rect fill="url(#hero-grid)" width="100%" height="100%" strokeWidth={0} />
      </svg>

      <div className="mx-auto max-w-7xl px-6 pt-16 pb-24 sm:pt-24 lg:px-8 lg:pt-32">
        <div className="mx-auto max-w-2xl text-center">
          <div className="mb-8 flex justify-center">
            <div className="inline-flex items-center gap-x-2 rounded-full px-3 py-1 text-sm/6 text-gray-600 ring-1 ring-gray-900/10 dark:text-gray-400 dark:ring-white/10">
              <span className="rounded-full bg-indigo-600 px-2 py-0.5 text-xs font-semibold text-white">
                {badge}
              </span>
              {badgeText}
            </div>
          </div>

          <h1 className="text-5xl font-semibold tracking-tight text-balance text-gray-900 sm:text-6xl dark:text-white">
            {title}
          </h1>
          <p className="mt-6 text-lg/8 text-pretty text-gray-600 dark:text-gray-400">
            {subtitle}
          </p>

          <div className="mt-10 flex items-center justify-center gap-x-6">
            <a
              href="#"
              className="rounded-md bg-indigo-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
            >
              {primaryLabel}
            </a>
            <a
              href="#"
              className="group inline-flex items-center gap-1 text-sm/6 font-semibold text-gray-900 dark:text-white"
            >
              {secondaryLabel}
              <ArrowRight
                aria-hidden="true"
                className="size-4 transition-transform group-hover:translate-x-0.5"
              />
            </a>
          </div>
        </div>

        {/* Browser frame */}
        <div className="mt-16 flow-root sm:mt-24">
          <div className="-m-2 rounded-xl bg-gray-900/5 p-2 ring-1 ring-gray-900/10 ring-inset lg:-m-4 lg:rounded-2xl lg:p-4 dark:bg-white/5 dark:ring-white/10">
            <div className="overflow-hidden rounded-lg bg-white shadow-2xl ring-1 ring-gray-900/10 dark:bg-gray-950 dark:ring-white/10">
              <div className="flex items-center gap-2 border-b border-gray-200 bg-gray-50 px-4 py-2.5 dark:border-white/10 dark:bg-gray-900">
                <span className="size-2.5 rounded-full bg-[#ff5f57]" />
                <span className="size-2.5 rounded-full bg-[#febc2e]" />
                <span className="size-2.5 rounded-full bg-[#28c840]" />
                <div className="mx-auto rounded-md bg-white px-8 py-1 font-mono text-[11px] text-gray-400 ring-1 ring-gray-200 dark:bg-gray-800 dark:text-gray-500 dark:ring-white/10">
                  app.example.com/admin/products
                </div>
              </div>

              {/* Simplified admin table — structure only, no real data */}
              <div className="flex">
                <div className="hidden w-44 shrink-0 border-r border-gray-200 p-4 sm:block dark:border-white/10">
                  <div className="mb-4 h-2.5 w-20 rounded bg-indigo-600/70" />
                  {[...Array(6)].map((_, i) => (
                    <div
                      key={i}
                      className={`mb-3 h-2 rounded ${
                        i === 1
                          ? 'w-24 bg-indigo-600/30'
                          : 'w-full bg-gray-200 dark:bg-gray-700'
                      }`}
                    />
                  ))}
                </div>
                <div className="flex-1 p-4 sm:p-6">
                  <div className="mb-5 flex items-center justify-between">
                    <div className="h-3 w-28 rounded bg-gray-300 dark:bg-gray-600" />
                    <div className="h-7 w-20 rounded-md bg-indigo-600" />
                  </div>
                  <div className="overflow-hidden rounded-lg ring-1 ring-gray-200 dark:ring-white/10">
                    {[...Array(5)].map((_, row) => (
                      <div
                        key={row}
                        className={`grid grid-cols-4 gap-4 px-4 py-3 ${
                          row === 0
                            ? 'bg-gray-50 dark:bg-gray-900'
                            : 'border-t border-gray-200 dark:border-white/10'
                        }`}
                      >
                        {[...Array(4)].map((_, col) => (
                          <div
                            key={col}
                            className={`h-2 rounded ${
                              row === 0
                                ? 'w-12 bg-gray-300 dark:bg-gray-600'
                                : 'w-full bg-gray-200 dark:bg-gray-800'
                            }`}
                          />
                        ))}
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Logo cloud */}
        <div className="mx-auto mt-16 max-w-lg sm:mt-20">
          <p className="text-center font-mono text-xs uppercase tracking-wider text-gray-500 dark:text-gray-500">
            Trusted by teams at
          </p>
          <div className="mt-6 flex flex-wrap items-center justify-center gap-x-8 gap-y-4">
            {LOGOS.map((logo) => (
              <span
                key={logo}
                className="text-sm font-semibold text-gray-400 dark:text-gray-600"
              >
                {logo}
              </span>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
