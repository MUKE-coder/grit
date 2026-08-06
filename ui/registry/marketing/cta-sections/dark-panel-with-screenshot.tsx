import { ChartNoAxesColumn, Folder, Globe, Radio, Search, Server, Settings } from 'lucide-react'

/*
 * Copy on the left, a slice of the product on the right, running off the edge
 * of a dark panel.
 *
 * The app on the right is built from markup rather than being an image, which
 * is worth explaining because the obvious move is to drop a PNG in. Markup
 * stays sharp on every display, weighs nothing, reflows instead of scaling, and
 * does not go stale the day someone renames a nav item. A screenshot is a
 * photograph of a build you shipped once.
 *
 * The whole mock carries `aria-hidden`, and that is the part people get wrong.
 * It is a picture of an app, not an app: none of it is reachable, half of it is
 * clipped by the panel edge, and a screen reader that walks into it starts
 * reading out invented project names in the middle of your call to action.
 * Decorative markup has to say so.
 *
 * It is also deliberately cut off. Showing the whole window shrinks it to
 * something unreadable; showing a slice at full size reads as a real interface
 * that continues past the frame.
 */

const PROJECTS = [
  { team: 'Planetaria', name: 'ios-app', live: false, meta: 'Initiated 1m 32s ago' },
  { team: 'Planetaria', name: 'mobile-api', live: true, meta: 'Deployed 3m ago · 23s' },
  { team: 'Grit Labs', name: 'gritframework.dev', live: false, meta: 'Initiated 5m 45s ago · 3m 4s' },
  { team: 'Grit Labs', name: 'ui.gritframework.dev', live: true, meta: 'Initiated 8m ago · 1m 30s' },
  { team: 'Protocol', name: 'relay-service', live: true, meta: 'Deployed 3h ago · 8s' },
  { team: 'Planetaria', name: 'android-app', live: true, meta: 'Deployed 12d ago · 5m 55s' },
]

const NAV = [
  { label: 'Projects', Icon: Folder },
  { label: 'Deployments', Icon: Server },
  { label: 'Activity', Icon: Radio },
  { label: 'Domains', Icon: Globe },
  { label: 'Usage', Icon: ChartNoAxesColumn },
  { label: 'Settings', Icon: Settings },
]

const TEAMS = ['Planetaria', 'Protocol', 'Grit Labs']

export default function CtaDarkPanelWithScreenshot({
  title = 'Boost your productivity.\nStart using our app today.',
  subtitle = 'Generate the API, the admin panel and the typed client from one definition. Deploy the whole thing as a single binary.',
  primaryLabel = 'Get started',
  primaryHref = '#',
  secondaryLabel = 'Learn more',
  secondaryHref = '#',
}: {
  title?: string
  subtitle?: string
  primaryLabel?: string
  primaryHref?: string
  secondaryLabel?: string
  secondaryHref?: string
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="relative isolate overflow-hidden rounded-3xl bg-gray-950 shadow-2xl">
          <div
            aria-hidden="true"
            className="pointer-events-none absolute inset-0 -z-10"
            style={{
              background:
                'radial-gradient(40rem 26rem at 28% 118%, rgba(139,92,246,0.4), transparent 66%)',
            }}
          />

          {/* The copy column is given the wider share: the mock is decorative
              and can be clipped further, the heading cannot reflow to four
              lines without looking like it overflowed. */}
          <div className="grid grid-cols-1 items-center gap-y-12 lg:grid-cols-[minmax(0,1.15fr)_minmax(0,1fr)]">
            <div className="px-6 pt-16 sm:px-14 lg:py-24 lg:pr-0">
              <h2 className="max-w-xl text-4xl font-semibold tracking-tight whitespace-pre-line text-white sm:text-5xl">
                {title}
              </h2>
              <p className="mt-6 max-w-md text-lg/8 text-pretty text-gray-300">{subtitle}</p>
              <div className="mt-10 flex flex-wrap items-center gap-x-6 gap-y-4">
                <a
                  href={primaryHref}
                  className="inline-flex min-h-11 items-center rounded-md bg-white px-4 text-sm font-semibold text-gray-900 shadow-sm hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
                >
                  {primaryLabel}
                </a>
                <a
                  href={secondaryHref}
                  className="inline-flex min-h-11 items-center rounded-md px-1 text-sm font-semibold text-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
                >
                  {secondaryLabel}<span aria-hidden="true">&nbsp;&rarr;</span>
                </a>
              </div>
            </div>

            {/* A picture of an app, not an app. Nothing here is focusable and
                nothing is announced. See the note at the top of the file. */}
            <div aria-hidden="true" className="relative h-[26rem] select-none lg:h-[34rem]">
              <div className="absolute inset-y-0 left-0 flex w-[46rem] rounded-tl-xl border-t border-l border-white/10 bg-gray-900/60 text-left backdrop-blur lg:top-12">
                <div className="w-56 flex-none border-r border-white/10 p-4">
                  <div className="flex items-center justify-between">
                    <div className="size-6 rounded bg-indigo-500/80" />
                    <div className="flex gap-1">
                      <span className="size-1 rounded-full bg-white/30" />
                      <span className="size-1 rounded-full bg-white/30" />
                      <span className="size-1 rounded-full bg-white/30" />
                    </div>
                  </div>
                  <p className="mt-6 text-[11px] font-medium tracking-wide text-gray-500 uppercase">
                    Navigation
                  </p>
                  <ul className="mt-3 space-y-2.5">
                    {NAV.map(({ label, Icon }) => (
                      <li key={label} className="flex items-center gap-2.5">
                        <Icon className="size-4 flex-none text-gray-400" />
                        <span className="text-sm text-gray-300">{label}</span>
                      </li>
                    ))}
                  </ul>
                  <p className="mt-7 text-[11px] font-medium tracking-wide text-gray-500 uppercase">
                    Your teams
                  </p>
                  <ul className="mt-3 space-y-2.5">
                    {TEAMS.map((team) => (
                      <li key={team} className="flex items-center gap-2.5">
                        <span className="flex size-4 items-center justify-center rounded border border-white/20 text-[9px] text-gray-400">
                          {team[0]}
                        </span>
                        <span className="text-sm text-gray-300">{team}</span>
                      </li>
                    ))}
                  </ul>
                </div>

                <div className="min-w-0 flex-1 bg-gray-950/60">
                  <div className="border-b border-white/10 p-4">
                    <div className="flex h-8 items-center gap-2 rounded-md border border-white/10 bg-white/5 px-3">
                      <Search className="size-3.5 text-gray-500" />
                      <span className="text-xs text-gray-500">Search projects...</span>
                    </div>
                  </div>
                  <div className="p-4">
                    <p className="text-sm font-semibold text-white">All projects</p>
                    <ul className="mt-3 divide-y divide-white/5">
                      {PROJECTS.map((project) => (
                        <li key={project.name} className="py-3">
                          <div className="flex items-center gap-2">
                            <span
                              className={`size-2 flex-none rounded-full ${
                                project.live ? 'bg-emerald-400' : 'bg-gray-600'
                              }`}
                            />
                            <span className="truncate text-sm font-semibold text-white">
                              {project.team}
                            </span>
                            <span className="text-sm text-gray-600">/</span>
                            <span className="truncate text-sm font-semibold text-white">
                              {project.name}
                            </span>
                          </div>
                          <p className="mt-1 pl-4 text-xs text-gray-500">
                            Deploys from GitHub · {project.meta}
                          </p>
                        </li>
                      ))}
                    </ul>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
