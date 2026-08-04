import Link from 'next/link'
import { ArrowRight } from 'lucide-react'
import { FrameworkLogo } from '@/components/framework-logo'
import { BenchmarkChart, type BenchScenario } from '@/components/benchmark-chart'
import { FRAMEWORKS, SCENARIOS, ranked, ratio, isAppBound } from '@/config/benchmarks'

/*
 * The homepage benchmark block, sat directly under the hero.
 *
 * Reads from config/benchmarks.ts so it cannot drift from the methodology pages
 * behind it — there is one set of numbers on this site, not two.
 */

export function HomepageBenchmarks() {
  if (FRAMEWORKS.length === 0) return null

  const scenarios: BenchScenario[] = SCENARIOS.map((s) => ({
    id: s.id,
    label: s.label,
    title: s.title,
    subtitle: s.subtitle,
    rows: ranked(s.id).map((f) => {
      const r = f.results[s.id]
      return {
        framework: f.name,
        slug: f.slug,
        version: f.version,
        logo: f.logo,
        color: f.color,
        invertOnDark: f.invertOnDark,
        ratio: ratio(r),
        gritRps: Math.round(r.gritRps),
        rps: Math.round(r.rps),
        appBound: isAppBound(r),
      }
    }),
  }))

  return (
    <section className="py-16 md:py-20 px-6 border-b border-border/40">
      <div className="max-w-5xl mx-auto">
        <div className="text-center mb-8">
          <span className="tag-mono text-primary/80 mb-3 block">Benchmarks</span>
          <h2 className="text-2xl md:text-3xl font-bold tracking-tight text-foreground mb-3">
            The same CRUD API, built in every framework
          </h2>
          <p className="text-sm md:text-base text-muted-foreground max-w-2xl mx-auto leading-relaxed">
            One Postgres, identical container limits, the same 10,000 rows, the same k6 script,
            every framework in production shape and on its own ORM. Each has a page showing
            exactly how to reproduce it &mdash; including the scenario Grit loses, and the two
            bugs this benchmark found in Grit itself.
          </p>
        </div>

        <BenchmarkChart scenarios={scenarios} />

        <div className="mt-6 flex flex-wrap items-center justify-center gap-2">
          {FRAMEWORKS.map((f) => (
            <Link
              key={f.slug}
              href={`/docs/benchmarks/${f.slug}`}
              className="inline-flex items-center gap-2 rounded-full border border-border/40 px-3.5 py-1.5 text-xs text-muted-foreground hover:border-primary/40 hover:text-foreground transition-colors"
            >
              <FrameworkLogo
                src={f.logo}
                alt=""
                onLight={f.invertOnDark}
                className="h-3.5 max-w-[36px]"
              />
              vs {f.name}
              <span className="text-muted-foreground/50">{f.orm}</span>
              <ArrowRight className="h-3 w-3" />
            </Link>
          ))}
        </div>
      </div>
    </section>
  )
}
