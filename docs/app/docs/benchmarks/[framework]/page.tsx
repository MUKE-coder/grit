import { notFound } from 'next/navigation'
import Link from 'next/link'
import Image from 'next/image'
import type { Metadata } from 'next'
import { ArrowLeft, ArrowRight, AlertTriangle, Scale, Video, CheckCircle2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { FrameworkLogo } from '@/components/framework-logo'
import { GUIDES, guideFor } from '@/config/benchmark-guides'
import {
  bySlug,
  GRIT,
  SCENARIOS,
  UNMEASURED,
  isAppBound,
  isDbBound,
  ratio,
  gritCpu,
} from '@/config/benchmarks'
import { siteConfig } from '@/config/site'

export function generateStaticParams() {
  return GUIDES.map((g) => ({ framework: g.slug }))
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ framework: string }>
}): Promise<Metadata> {
  const { framework } = await params
  const guide = GUIDES.find((g) => g.slug === framework)
  const fw = bySlug(framework)
  if (!guide) return {}

  const title = fw ? `Grit vs ${fw.name} — benchmark methodology` : guide.videoTitle
  const description =
    `Step-by-step reproduction of the Grit vs ${fw?.name ?? framework} benchmark: the exact ` +
    `commands, the application code, how the comparison is kept fair, and what numbers to expect.`

  return {
    title,
    description,
    alternates: { canonical: `${siteConfig.url}/docs/benchmarks/${framework}` },
    openGraph: { title: `${title} | Grit`, description, type: 'article' },
  }
}

export default async function FrameworkBenchmarkPage({
  params,
}: {
  params: Promise<{ framework: string }>
}) {
  const { framework } = await params
  const guide = guideFor(framework)
  const fw = bySlug(framework)
  const unmeasured = UNMEASURED.find((u) => u.slug === framework)
  // The 'grit' guide has no opponent, so it gets no lockup.
  const logo = fw ?? unmeasured
  if (!guide) notFound()


  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Benchmark methodology</span>

              {/* The two marks, side by side — it is a head-to-head and the page
                  should look like one before you read a word of it. */}
              {logo && (
                <div className="flex items-center gap-4 mb-5">
                  <Image
                    src={GRIT.logo}
                    alt="Grit"
                    width={44}
                    height={44}
                    className="h-11 w-11 rounded-xl"
                  />
                  <span className="text-xl font-light text-muted-foreground/50">vs</span>
                  <FrameworkLogo
                    src={logo.logo}
                    alt={logo.name}
                    onLight={logo.invertOnDark}
                    className="h-10 max-w-[150px]"
                  />
                </div>
              )}

              <h1 className="text-4xl font-bold tracking-tight mb-4">
                {fw ? `Grit vs ${fw.name}` : guide.videoTitle}
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">{guide.intro}</p>

              {fw && (
                <p className="text-sm text-muted-foreground mt-4">
                  <strong className="text-foreground">{fw.stack}</strong> &middot; {fw.version}
                </p>
              )}
            </div>

            {/* ── The numbers being reproduced ─────────────────────── */}
            {fw && (
              <section className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  What you should end up with
                </h2>
                <div className="overflow-x-auto rounded-xl border border-border/40">
                  <table className="w-full text-sm border-collapse">
                    <thead>
                      <tr className="border-b border-border/40 bg-muted/30">
                        <th className="text-left px-4 py-2 font-medium">Scenario</th>
                        <th className="text-right px-4 py-2 font-medium">Grit</th>
                        <th className="text-right px-4 py-2 font-medium">{fw.name}</th>
                        <th className="text-right px-4 py-2 font-medium">Ratio</th>
                      </tr>
                    </thead>
                    <tbody>
                      {SCENARIOS.map((s) => {
                        const f = fw.results[s.id]
                        // Grit's figure comes from THIS pair, measured minutes
                        // apart from the opponent's. Its absolute value differs
                        // between pairs because the machine drifts; the ratio is
                        // what carries across.
                        const bound = (x: Parameters<typeof isAppBound>[0]) =>
                          isAppBound(x) ? 'its own ceiling'
                            : isDbBound(x) ? 'database-bound'
                            : 'neither saturated'
                        return (
                          <tr key={s.id} className="border-b border-border/20 last:border-0">
                            <td className="px-4 py-3">
                              <div className="font-mono text-xs text-primary">{s.id}</div>
                              <div className="text-xs text-muted-foreground mt-0.5">
                                {s.request}
                              </div>
                            </td>
                            <td className="px-4 py-3 text-right">
                              <div className="font-semibold">
                                {f.gritRps.toLocaleString()} req/s
                              </div>
                              <div className="text-xs text-muted-foreground">
                                {f.gritMedian} &middot; {bound(gritCpu(f))}
                              </div>
                            </td>
                            <td className="px-4 py-3 text-right">
                              <div className="font-semibold">
                                {f.rps.toLocaleString()} req/s
                              </div>
                              <div className="text-xs text-muted-foreground">
                                {f.median} &middot; {bound(f)}
                              </div>
                            </td>
                            <td
                              className={
                                'px-4 py-3 text-right font-semibold ' +
                                (ratio(f) < 1 ? 'text-warning' : 'text-primary')
                              }
                            >
                              {ratio(f).toFixed(2)}&times;
                            </td>
                          </tr>
                        )
                      })}
                    </tbody>
                  </table>
                </div>
                <p className="text-sm text-muted-foreground mt-3 leading-relaxed">
                  Both figures come from the same run, minutes apart. Your absolute numbers
                  will differ &mdash; different CPU, different disk, different background load
                  &mdash; and so will ours: Grit&apos;s single-row read measured 6,600 req/s in
                  the Bun pair and 1,635 in the Express pair from an identical binary, as hours
                  of write scenarios accumulated in Postgres. The <em>ratio</em> is what should
                  survive. If it does not, something in the setup differs and the steps below are
                  where to look.
                </p>
              </section>
            )}

            {/* Guides exist for frameworks that have no published pair yet. Say
                why, rather than quietly showing a page with no numbers on it. */}
            {!fw && unmeasured && (
              <section className="mb-12">
                <div className="flex gap-3 rounded-xl border border-warning/30 bg-warning/[0.05] px-5 py-4">
                  <AlertTriangle className="h-5 w-5 shrink-0 text-warning mt-0.5" />
                  <div>
                    <h2 className="font-semibold text-foreground mb-1">
                      No published numbers for {unmeasured.name} yet
                    </h2>
                    <p className="text-sm text-muted-foreground leading-relaxed">
                      {unmeasured.reason} Everything below is the full method, so you can run it
                      yourself before we do.
                    </p>
                  </div>
                </div>
              </section>
            )}

            {/* ── Fairness ─────────────────────────────────────────── */}
            <section className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center gap-2">
                <Scale className="h-5 w-5 text-primary" />
                How {fw?.name ?? 'this framework'} is kept from being handicapped
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                A benchmark is only worth publishing if the loser was given every reasonable
                advantage. These are the specific decisions made for this framework, and each one
                is worth saying out loud on camera.
              </p>
              <ul className="space-y-2.5 text-muted-foreground leading-relaxed list-disc pl-5">
                {guide.fairness.map((f, i) => (
                  <li key={i}>{f}</li>
                ))}
              </ul>
            </section>

            {/* ── Steps ────────────────────────────────────────────── */}
            <section className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-6 flex items-center gap-2">
                <Video className="h-5 w-5 text-primary" />
                Reproduce it, start to finish
              </h2>

              <ol className="space-y-8">
                {guide.steps.map((step, i) => (
                  <li key={i} className="relative pl-10">
                    <span className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-full bg-primary/10 text-xs font-semibold text-primary">
                      {i + 1}
                    </span>
                    <h3 className="font-semibold text-foreground mb-1.5">{step.title}</h3>
                    <p className="text-sm text-muted-foreground leading-relaxed mb-3">
                      {step.body}
                    </p>
                    {step.code && (
                      <CodeBlock language={step.language ?? 'bash'} code={step.code} />
                    )}
                    {step.warning && (
                      <div className="mt-3 flex gap-2.5 rounded-lg border border-warning/30 bg-warning/[0.05] px-4 py-3">
                        <AlertTriangle className="h-4 w-4 shrink-0 text-warning mt-0.5" />
                        <p className="text-sm text-muted-foreground leading-relaxed">
                          {step.warning}
                        </p>
                      </div>
                    )}
                  </li>
                ))}
              </ol>
            </section>

            {/* ── Expect ───────────────────────────────────────────── */}
            <section className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center gap-2">
                <CheckCircle2 className="h-5 w-5 text-success" />
                Reading your own results
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">{guide.expect}</p>
              <div className="rounded-xl border border-border/40 p-5">
                <p className="text-sm text-muted-foreground leading-relaxed">
                  <code>aggregate.py</code> prints the app container&apos;s CPU and Postgres&apos;s
                  CPU next to every row. That pair is what tells you whether a number is a real
                  ceiling. If the app is pinned near 400% then you are seeing its limit. If the app
                  is idling while Postgres is near 800%, the database gave out first &mdash; that
                  row is a floor, the framework would go faster on a bigger database, and quoting
                  it as &ldquo;X does N req/s&rdquo; overstates what you measured.
                </p>
              </div>
            </section>

            {/* ── Other frameworks ─────────────────────────────────── */}
            <section className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">The other guides</h2>
              <div className="grid gap-2 sm:grid-cols-2">
                {GUIDES.filter((g) => g.slug !== guide.slug).map((g) => {
                  const other = bySlug(g.slug)
                  const mark = other ?? UNMEASURED.find((u) => u.slug === g.slug)
                  return (
                    <Link
                      key={g.slug}
                      href={`/docs/benchmarks/${g.slug}`}
                      className="flex items-center gap-3 rounded-lg border border-border/40 px-4 py-3 hover:border-primary/40 hover:bg-muted/30 transition-colors"
                    >
                      {mark && (
                        <span className="w-11 shrink-0 flex justify-center">
                          <FrameworkLogo
                            src={mark.logo}
                            alt=""
                            onLight={mark.invertOnDark}
                            className="h-6 max-w-[44px]"
                          />
                        </span>
                      )}
                      <div className="min-w-0">
                        <div className="font-medium text-foreground">
                          {mark ? `Grit vs ${mark.name}` : g.slug}
                        </div>
                        <div className="text-xs text-muted-foreground mt-0.5 line-clamp-1">
                          {other?.tagline ?? 'Method published, numbers pending'}
                        </div>
                      </div>
                    </Link>
                  )
                })}
              </div>
            </section>

            <div className="flex items-center justify-between pt-6 border-t border-border/40">
              <Button variant="ghost" asChild>
                <Link href="/docs/benchmarks">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  All benchmarks
                </Link>
              </Button>
              <Button variant="ghost" asChild>
                <Link href="/docs/testing">
                  Performance testing
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
