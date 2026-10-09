import Link from 'next/link'
import { notFound } from 'next/navigation'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { SystemDesignBody } from '@/components/system-design'
import { SYSTEMS, getSystem } from '@/config/systems'

/**
 * One page per system, rendered from config/systems.ts.
 *
 * A shared template rather than thirty files, for the reason the deployment
 * guides give: the value of a set like this is that every entry answers the
 * same questions in the same order. Hand-written pages drift, and then a
 * missing section is ambiguous between "does not apply" and "unfinished".
 */

export function generateStaticParams() {
  return SYSTEMS.map((s) => ({ system: s.slug }))
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ system: string }>
}) {
  const { system: slug } = await params
  const s = getSystem(slug)
  if (!s) return { title: 'Not found' }
  return {
    title: `${s.name} System Design`,
    description: `${s.tagline} Requirements, capacity estimation, high and low level design, API surface, scaling and bottlenecks for Grit's ${s.name.toLowerCase()} system.`,
  }
}

export default async function SystemPage({
  params,
}: {
  params: Promise<{ system: string }>
}) {
  const { system: slug } = await params
  const s = getSystem(slug)
  if (!s) notFound()

  const index = SYSTEMS.findIndex((x) => x.slug === slug)
  const prev = SYSTEMS[index - 1]
  const next = SYSTEMS[index + 1]

  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl px-6 py-10">
          <div className="max-w-3xl">
            <Link
              href="/docs/systems"
              className="mb-6 inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground"
            >
              <ArrowLeft className="h-3.5 w-3.5" />
              All systems
            </Link>

            <div className="mb-8">
              <span className="tag-mono mb-3 block text-primary/80">{s.group}</span>
              <h1 className="mb-4 text-4xl font-bold tracking-tight">
                {s.name} <span className="text-muted-foreground">System Design</span>
              </h1>
              <p className="text-lg leading-relaxed text-muted-foreground">{s.tagline}</p>
              <div className="mt-5 flex flex-wrap gap-1.5">
                {s.packages.map((p) => (
                  <code
                    key={p}
                    className="rounded bg-accent/50 px-2 py-1 font-mono text-xs text-muted-foreground"
                  >
                    {p}
                  </code>
                ))}
              </div>
            </div>

            <SystemDesignBody s={s} />

            {s.seeAlso && s.seeAlso.length > 0 && (
              <section className="mt-12">
                <h2 className="text-2xl font-bold tracking-tight">Read next</h2>
                <ul className="mt-4 space-y-2">
                  {s.seeAlso.map((l) => (
                    <li key={l.href}>
                      <Link
                        href={l.href}
                        className="inline-flex items-center gap-1.5 text-[15px] text-primary hover:underline"
                      >
                        {l.title}
                        <ArrowRight className="h-3.5 w-3.5" />
                      </Link>
                    </li>
                  ))}
                </ul>
              </section>
            )}

            <nav className="mt-14 flex flex-col gap-3 border-t border-border/50 pt-6 sm:flex-row sm:justify-between">
              {prev ? (
                <Link
                  href={`/docs/systems/${prev.slug}`}
                  className="group flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground"
                >
                  <ArrowLeft className="h-4 w-4" />
                  <span>
                    <span className="block text-xs">Previous</span>
                    <span className="font-medium text-foreground">{prev.name}</span>
                  </span>
                </Link>
              ) : (
                <span />
              )}
              {next && (
                <Link
                  href={`/docs/systems/${next.slug}`}
                  className="group flex items-center gap-2 text-right text-sm text-muted-foreground hover:text-foreground sm:ml-auto"
                >
                  <span>
                    <span className="block text-xs">Next</span>
                    <span className="font-medium text-foreground">{next.name}</span>
                  </span>
                  <ArrowRight className="h-4 w-4" />
                </Link>
              )}
            </nav>
          </div>
        </div>
      </main>
    </div>
  )
}
