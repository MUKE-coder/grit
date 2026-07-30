import { ArrowRight } from 'lucide-react'

/** The Grit UI mark, inlined so the block stays self-contained. */
function GritMark({ className = 'size-8' }: { className?: string }) {
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

const NAV = ['Product', 'Features', 'Marketplace', 'Company']

const PROJECTS = [
  { team: 'Planetaria', name: 'ios-app', state: 'Preview', tone: 'gray', meta: 'Initiated 1m 32s ago' },
  { team: 'Planetaria', name: 'mobile-api', state: 'Production', tone: 'indigo', meta: 'Deployed 3m ago · 23s' },
  { team: 'Grit Labs', name: 'gritframework.dev', state: 'Preview', tone: 'gray', meta: 'Initiated 5m 45s ago · 3m 4s' },
  { team: 'Grit Labs', name: 'ui.gritframework.dev', state: 'Preview', tone: 'gray', meta: 'Initiated 8m ago · 1m 30s' },
  { team: 'Protocol', name: 'relay-service', state: 'Production', tone: 'indigo', meta: 'Deployed 3h ago · 8s' },
]

const FEED = [
  { who: 'Cosetta Dusett', what: 'Pushed to ios-app (27c83j on main)', when: '30s' },
  { who: 'Pammi Kakani', what: 'Pushed to mobile-api (29jsd on main)', when: '3m' },
  { who: 'Kora Grisostomo', what: 'Pushed to ios-app (cdd2d on main)', when: '4m' },
  { who: 'Jean-Francois Tippy', what: 'Pushed to relay-service (9da3c on main)', when: '8m' },
]

export default function DarkWithDashboard({
  title = 'Data to enrich your online business',
  subtitle = 'Generate the API, the admin panel and the typed client from one definition. Deploy the whole thing as a single binary.',
  primaryLabel = 'Get started',
  secondaryLabel = 'Learn more',
}: {
  title?: string
  subtitle?: string
  primaryLabel?: string
  secondaryLabel?: string
}) {
  return (
    // Dark in both themes by design — a dark hero, not a light one inverted.
    <div className="relative isolate overflow-hidden bg-gray-950">
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-x-0 top-0 -z-10 h-[48rem]"
        style={{
          background:
            'radial-gradient(48rem 32rem at 50% -6rem, rgba(99,102,241,0.28), transparent 65%)',
        }}
      />

      {/* Nav */}
      <nav
        aria-label="Global"
        className="mx-auto flex max-w-7xl items-center justify-between px-6 py-6"
      >
        <a href="#" className="flex items-center gap-2 text-white">
          <GritMark className="size-8" />
          <span className="sr-only">Grit</span>
        </a>
        <div className="hidden gap-10 lg:flex">
          {NAV.map((item) => (
            <a key={item} href="#" className="text-sm font-semibold text-white">
              {item}
            </a>
          ))}
        </div>
        <a href="#" className="text-sm font-semibold text-white">
          Log in <span aria-hidden="true">&rarr;</span>
        </a>
      </nav>

      {/* Hero copy */}
      <div className="mx-auto max-w-3xl px-6 pt-16 pb-20 text-center sm:pt-24">
        <h1 className="text-5xl font-bold tracking-tight text-balance text-white sm:text-7xl">
          {title}
        </h1>
        <p className="mx-auto mt-8 max-w-2xl text-lg/8 text-gray-400">{subtitle}</p>
        <div className="mt-10 flex items-center justify-center gap-x-6">
          <a
            href="#"
            className="rounded-lg bg-indigo-500 px-5 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-indigo-400"
          >
            {primaryLabel}
          </a>
          <a
            href="#"
            className="group inline-flex items-center gap-1.5 text-sm font-semibold text-white"
          >
            {secondaryLabel}
            <ArrowRight
              aria-hidden="true"
              className="size-4 transition-transform group-hover:translate-x-0.5"
            />
          </a>
        </div>
      </div>

      {/* Dashboard */}
      <div className="mx-auto max-w-7xl px-6 pb-24">
        <div className="overflow-hidden rounded-2xl border border-white/10 bg-gray-900/70 shadow-2xl backdrop-blur">
          <div className="flex">
            {/* Sidebar */}
            <div className="hidden w-56 shrink-0 border-r border-white/5 p-5 lg:block">
              <div className="mb-8 flex items-center gap-2 text-white">
                <GritMark className="size-6" />
              </div>
              <p className="mb-3 font-mono text-[10px] tracking-wider text-gray-600 uppercase">
                Navigation
              </p>
              {['Projects', 'Deployments', 'Activity', 'Domains', 'Usage', 'Settings'].map(
                (item, i) => (
                  <div
                    key={item}
                    className={`mb-1 rounded-md px-2.5 py-1.5 text-sm ${
                      i === 1 ? 'bg-white/10 text-white' : 'text-gray-400'
                    }`}
                  >
                    {item}
                  </div>
                ),
              )}
            </div>

            {/* Projects */}
            <div className="min-w-0 flex-1 p-5">
              <div className="mb-5 flex items-center gap-4">
                <div className="flex-1 rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-sm text-gray-500">
                  Search projects…
                </div>
                <span className="hidden rounded-lg bg-indigo-500 px-3.5 py-2 text-sm font-medium text-white sm:block">
                  + New project
                </span>
              </div>

              <p className="mb-3 text-sm font-medium text-white">All projects</p>
              <div className="space-y-1">
                {PROJECTS.map((p) => (
                  <div
                    key={p.name}
                    className="flex items-center gap-3 rounded-lg px-3 py-2.5 hover:bg-white/5"
                  >
                    <span
                      className={`size-2 shrink-0 rounded-full ${
                        p.tone === 'indigo' ? 'bg-green-400' : 'bg-gray-600'
                      }`}
                    />
                    <div className="min-w-0 flex-1">
                      <p className="truncate text-sm text-white">
                        <span className="text-gray-400">{p.team}</span> / {p.name}
                      </p>
                      <p className="truncate text-xs text-gray-500">{p.meta}</p>
                    </div>
                    <span
                      className={`shrink-0 rounded-full px-2.5 py-0.5 text-xs font-medium ${
                        p.tone === 'indigo'
                          ? 'bg-indigo-500 text-white'
                          : 'bg-white/10 text-gray-300'
                      }`}
                    >
                      {p.state}
                    </span>
                  </div>
                ))}
              </div>
            </div>

            {/* Activity */}
            <div className="hidden w-80 shrink-0 border-l border-white/5 p-5 xl:block">
              <p className="mb-4 text-sm font-medium text-white">Activity feed</p>
              <div className="space-y-4">
                {FEED.map((f) => (
                  <div key={f.who + f.when} className="flex gap-3">
                    <span className="size-7 shrink-0 rounded-full bg-gradient-to-br from-indigo-400 to-fuchsia-500" />
                    <div className="min-w-0 flex-1">
                      <p className="flex items-baseline justify-between gap-2 text-sm text-white">
                        <span className="truncate">{f.who}</span>
                        <span className="shrink-0 text-xs text-gray-600">{f.when}</span>
                      </p>
                      <p className="truncate text-xs text-gray-500">{f.what}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
