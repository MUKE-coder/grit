import Link from 'next/link'
import { ArrowRight, Check, Minus } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { getDocMetadata } from '@/config/docs-metadata'
import {
  DEPLOYMENT_PROVIDERS,
  PROVIDER_KIND_LABEL,
  type ProviderKind,
} from '@/config/deployment-providers'
import { GUIDE_COMPARISON } from '@/config/deployment-guides'

export const metadata = getDocMetadata('/docs/deployment')

const ORDER: ProviderKind[] = ['paas', 'self-hosted', 'vps', 'container']

/* The comparison table names providers in prose; this maps those names back to
   the page they belong to, so the first column links rather than being inert
   text next to four links. */
const GUIDE_SLUG_BY_NAME: Record<string, string> = {
  Dokploy: 'dokploy',
  Coolify: 'coolify',
  Render: 'render',
  'Fly.io': 'fly-io',
  Railway: 'railway',
}

export default function DeploymentIndexPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Deployment</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Where to deploy Grit</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit produces an ordinary Docker image and an ordinary Go binary. Anything
                that runs either will run your app — there is no platform it is coupled to.
                This section is about picking well, then the specifics of each.
              </p>
            </div>

            {/* The decision, before the options. A list of seven providers is not
                an answer to "where should I deploy this". */}
            <div className="rounded-xl border border-border/50 bg-card/50 p-6 mb-12">
              <h2 className="text-base font-semibold mb-4">If you have no opinion yet</h2>
              <ul className="space-y-3 text-sm text-muted-foreground leading-relaxed">
                <li>
                  <strong className="text-foreground">Shipping something today, alone or in a small team →</strong>{' '}
                  <Link href="/docs/deployment/railway" className="text-primary hover:underline">
                    Railway
                  </Link>{' '}
                  or{' '}
                  <Link href="/docs/deployment/render" className="text-primary hover:underline">
                    Render
                  </Link>
                  . Managed database, no server to patch, running in under an hour.
                </li>
                <li>
                  <strong className="text-foreground">You want a config file in the repo rather than a dashboard →</strong>{' '}
                  <Link href="/docs/deployment/fly-io" className="text-primary hover:underline">
                    Fly.io
                  </Link>
                  . One <code className="text-xs">fly.toml</code> per app, deployed with{' '}
                  <code className="text-xs">flyctl</code> and a GitHub Actions workflow.
                </li>
                <li>
                  <strong className="text-foreground">Several apps, one bill →</strong>{' '}
                  <Link href="/docs/deployment/dokploy" className="text-primary hover:underline">
                    Dokploy
                  </Link>{' '}
                  or{' '}
                  <Link href="/docs/deployment/coolify" className="text-primary hover:underline">
                    Coolify
                  </Link>
                  . A platform you own, on one box.
                </li>
                <li>
                  <strong className="text-foreground">On-premise, or the customer owns the hardware →</strong>{' '}
                  <Link href="/docs/deployment/dokploy" className="text-primary hover:underline">
                    Dokploy
                  </Link>{' '}
                  on their box. It runs your real{' '}
                  <code className="text-xs">docker-compose.prod.yml</code> almost unmodified.
                </li>
              </ul>
              <p className="mt-5 text-xs text-muted-foreground">
                Every one of these is a real production answer. The wrong choice here is
                recoverable — moving a Grit app between them is a Dockerfile and a set of
                environment variables.
              </p>
            </div>

            {/* Before you pick anything */}
            <h2 className="text-2xl font-bold tracking-tight mb-4">First, four things that apply everywhere</h2>
            <div className="grid gap-3 sm:grid-cols-2 mb-14">
              {[
                {
                  href: '/docs/deployment/environment',
                  title: 'Environment variables',
                  body: 'Which are required, which are build-time, and the ones that fail silently.',
                },
                {
                  href: '/docs/infrastructure/database',
                  title: 'Database & migrations',
                  body: 'Why production does not auto-migrate, and how to run them safely.',
                },
                {
                  href: '/docs/deployment/build-locally',
                  title: 'Test the production build locally',
                  body: 'Reproduce the deploy on your machine first. Most failures are visible here.',
                },
                {
                  href: '/docs/deployment/checklist',
                  title: 'Go-live checklist',
                  body: 'The things that are embarrassing to discover after launch.',
                },
              ].map((c) => (
                <Link
                  key={c.href}
                  href={c.href}
                  className="group rounded-xl border border-border/50 bg-card/40 p-5 transition-colors hover:border-border hover:bg-card/70"
                >
                  <div className="flex items-center gap-1.5 font-semibold text-sm">
                    {c.title}
                    <ArrowRight className="h-3.5 w-3.5 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
                  </div>
                  <p className="mt-1.5 text-[13px] text-muted-foreground leading-relaxed">{c.body}</p>
                </Link>
              ))}
            </div>

            {/* Comparison */}
            <h2 className="text-2xl font-bold tracking-tight mb-4">Side by side</h2>
            <div className="overflow-x-auto mb-14 rounded-xl border border-border/50">
              <table className="w-full text-sm">
                <thead className="bg-card/60">
                  <tr className="text-left">
                    <th className="px-4 py-3 font-semibold">Provider</th>
                    <th className="px-4 py-3 font-semibold">From</th>
                    <th className="px-4 py-3 font-semibold">Ops effort</th>
                    <th className="px-4 py-3 font-semibold whitespace-nowrap">Managed DB</th>
                    <th className="px-4 py-3 font-semibold whitespace-nowrap">Disk</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border/40">
                  {DEPLOYMENT_PROVIDERS.map((p) => (
                    <tr key={p.slug} className="hover:bg-card/40">
                      <td className="px-4 py-3">
                        <Link href={`/docs/deployment/${p.slug}`} className="font-medium text-primary hover:underline">
                          {p.name}
                        </Link>
                      </td>
                      <td className="px-4 py-3 text-muted-foreground whitespace-nowrap">{p.costFrom}</td>
                      <td className="px-4 py-3 text-muted-foreground">{p.effort}</td>
                      <td className="px-4 py-3">
                        {p.managedPostgres ? (
                          <Check className="h-4 w-4 text-emerald-500" />
                        ) : (
                          <Minus className="h-4 w-4 text-muted-foreground/50" />
                        )}
                      </td>
                      <td className="px-4 py-3">
                        {p.persistentDisk ? (
                          <Check className="h-4 w-4 text-emerald-500" />
                        ) : (
                          <Minus className="h-4 w-4 text-muted-foreground/50" />
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* The four questions that actually differ between these platforms.
                Everything else is detail you can absorb after choosing. */}
            <h2 className="text-2xl font-bold tracking-tight mb-2">How each one handles your Compose file</h2>
            <p className="text-sm text-muted-foreground leading-relaxed mb-4">
              Two of these run <code className="text-xs">docker-compose.prod.yml</code> nearly
              as written. The other three do not run Compose at all and need it translated into
              their own model, which changes how migrations are ordered and where domains are
              configured.
            </p>
            <div className="overflow-x-auto mb-14 rounded-xl border border-border/50">
              <table className="w-full text-left text-sm">
                <thead className="bg-card/60">
                  <tr>
                    {GUIDE_COMPARISON.headers.map((h, i) => (
                      <th key={i} scope="col" className="px-4 py-3 font-semibold whitespace-nowrap">
                        {h}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody className="divide-y divide-border/40">
                  {GUIDE_COMPARISON.rows.map((row) => (
                    <tr key={row[0]} className="hover:bg-card/40">
                      {row.map((cell, c) => (
                        <td
                          key={c}
                          className={`px-4 py-3 align-top ${
                            c === 0 ? 'font-medium whitespace-nowrap' : 'text-muted-foreground'
                          }`}
                        >
                          {c === 0 ? (
                            <Link
                              href={`/docs/deployment/${GUIDE_SLUG_BY_NAME[cell] ?? ''}`}
                              className="text-primary hover:underline"
                            >
                              {cell}
                            </Link>
                          ) : (
                            cell
                          )}
                        </td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Full list, grouped */}
            {ORDER.map((kind) => {
              const group = DEPLOYMENT_PROVIDERS.filter((p) => p.kind === kind)
              if (group.length === 0) return null
              return (
                <div key={kind} className="mb-10">
                  <h3 className="tag-mono text-primary/80 mb-4">{PROVIDER_KIND_LABEL[kind]}</h3>
                  <div className="space-y-3">
                    {group.map((p) => (
                      <Link
                        key={p.slug}
                        href={`/docs/deployment/${p.slug}`}
                        className="group block rounded-xl border border-border/50 bg-card/40 p-5 transition-colors hover:border-border hover:bg-card/70"
                      >
                        <div className="flex flex-wrap items-baseline gap-x-3 gap-y-1">
                          <span className="font-semibold">{p.name}</span>
                          <span className="text-xs text-muted-foreground">{p.costFrom}/mo</span>
                        </div>
                        <p className="mt-1.5 text-sm text-muted-foreground leading-relaxed">{p.tagline}</p>
                        <p className="mt-2 text-[13px] text-muted-foreground/80 leading-relaxed">
                          <strong className="text-foreground/70">Best for:</strong> {p.bestFor}
                        </p>
                      </Link>
                    ))}
                  </div>
                </div>
              )
            })}
          </div>
        </div>
      </main>
    </div>
  )
}
