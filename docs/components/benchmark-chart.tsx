'use client'

import { useState } from 'react'
import Link from 'next/link'
import { ArrowRight } from 'lucide-react'

/*
 * The homepage benchmark chart.
 *
 * Plots the RATIO within each head-to-head pair, not raw requests per second,
 * and that is a correctness decision rather than a stylistic one.
 *
 * Every framework was measured against Grit in its own back-to-back run. Across
 * those runs the machine drifted a long way — Grit's single-row read measured
 * 6,600 req/s in the Bun pair, 4,392 in the Encore pair and 1,635 in the Express
 * pair, from an identical binary. Putting Bun's 3,196 next to Encore's 438 on
 * one axis would imply they faced the same Grit. They did not, and Encore would
 * look far worse than it is.
 *
 * The ratio is the part that survives: within a pair, both sides ran minutes
 * apart under the same conditions. Absolute numbers and their Grit baseline live
 * on each framework's own page, where the pairing is stated.
 */

export type BenchRow = {
  framework: string
  slug: string
  version: string
  /** how many times faster Grit was in this pair; below 1 means Grit lost */
  ratio: number
  /** what Grit measured in this pair, for the tooltip */
  gritRps: number
  /** what the opponent measured in the same pair */
  rps: number
  /** true when the app container saturated — a real ceiling, not a DB limit */
  appBound: boolean
}

export type BenchScenario = {
  id: string
  label: string
  title: string
  subtitle: string
  rows: BenchRow[]
}

export function BenchmarkChart({ scenarios }: { scenarios: BenchScenario[] }) {
  const [active, setActive] = useState(scenarios[0]?.id)
  const scenario = scenarios.find((s) => s.id === active) ?? scenarios[0]
  if (!scenario) return null

  const max = Math.max(...scenario.rows.map((r) => r.ratio), 1)

  /*
   * Square-root heights. Grit is 41.8x Laravel on a single-row read and 2.07x
   * Bun; on a linear axis the Bun column would be a 5% stub you cannot read, and
   * the chart would say "Bun is roughly as slow as Laravel", which is the
   * opposite of the truth. Sqrt keeps every column legible and the ordering
   * exact. The printed figure on each bar is the real ratio, and the footnote
   * says the axis is compressed — a reader is never left to infer scale from
   * bar height alone.
   */
  const height = (ratio: number) => (Math.sqrt(ratio) / Math.sqrt(max)) * 100

  return (
    <div className="rounded-2xl border border-border/40 bg-[#0d0d12] p-5 sm:p-8">
      <div
        role="tablist"
        aria-label="Benchmark scenario"
        className="flex flex-wrap items-center justify-center gap-1 mb-7"
      >
        {scenarios.map((s) => (
          <button
            key={s.id}
            role="tab"
            aria-selected={s.id === scenario.id}
            onClick={() => setActive(s.id)}
            className={
              'rounded-lg px-3.5 py-1.5 text-sm font-medium transition-colors ' +
              (s.id === scenario.id
                ? 'bg-primary/15 text-primary'
                : 'text-muted-foreground hover:text-foreground hover:bg-white/5')
            }
          >
            {s.label}
          </button>
        ))}
      </div>

      <div className="text-center mb-2">
        <h3 className="text-xl sm:text-2xl font-semibold text-foreground">
          How many times faster Grit is
        </h3>
        <p className="text-sm text-muted-foreground mt-1">{scenario.subtitle}</p>
      </div>

      <div className="flex items-center justify-center gap-5 mb-8 text-xs text-muted-foreground">
        <span className="inline-flex items-center gap-1.5">
          <span className="h-2.5 w-2.5 rounded-sm bg-primary" />
          Grit ahead
        </span>
        <span className="inline-flex items-center gap-1.5">
          <span className="h-2.5 w-2.5 rounded-sm bg-[#f5a623]" />
          Grit behind
        </span>
      </div>

      {/* Vertical columns, value inside, framework and version beneath. Scrolls
          sideways on a phone rather than crushing six columns into 320px. */}
      <div className="overflow-x-auto -mx-2 px-2">
        <div
          className="flex items-end justify-center gap-2 sm:gap-3"
          style={{ minHeight: 300, minWidth: scenario.rows.length * 82 }}
        >
          {scenario.rows.map((row) => {
            const pct = Math.max(4, height(row.ratio))
            const behind = row.ratio < 1

            return (
              <Link
                key={row.slug}
                href={`/docs/benchmarks/${row.slug}`}
                className="group flex flex-1 flex-col items-center min-w-[74px] max-w-[130px]"
                title={`Grit ${row.gritRps.toLocaleString()} req/s vs ${row.framework} ${row.rps.toLocaleString()} req/s, measured back to back`}
              >
                <div className="flex h-[240px] w-full flex-col justify-end items-center">
                  {/* Label above the bar, not inside it — a short column has no
                      room for text, and that is exactly where the number matters. */}
                  <span
                    className={
                      'mb-1.5 text-[13px] font-bold tabular-nums ' +
                      (behind ? 'text-[#f5a623]' : 'text-primary')
                    }
                  >
                    {row.ratio.toFixed(2)}&times;
                  </span>
                  <div
                    className={
                      'w-full rounded-t-md transition-all duration-500 group-hover:brightness-110 ' +
                      (behind
                        ? 'bg-[#f5a623]'
                        : 'bg-gradient-to-t from-primary/70 to-primary')
                    }
                    style={{ height: `${pct}%` }}
                  />
                </div>

                <div className="mt-3 text-center px-0.5">
                  <div className="text-[13px] font-medium leading-tight text-muted-foreground group-hover:text-primary transition-colors">
                    vs {row.framework}
                  </div>
                  <div className="text-[10px] font-mono text-muted-foreground/60 leading-tight mt-0.5">
                    {row.version}
                  </div>
                  <div className="text-[10px] text-muted-foreground/50 leading-tight mt-0.5 tabular-nums">
                    {row.gritRps.toLocaleString()} v {row.rps.toLocaleString()}
                  </div>
                </div>
              </Link>
            )
          })}
        </div>
      </div>

      <div className="mt-7 pt-4 border-t border-border/40 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
        <p className="text-xs text-muted-foreground leading-relaxed max-w-2xl">
          Each framework was run head to head against Grit, back to back, three times, medians
          reported, zero failed requests. 50 concurrent users, 4 CPUs and 2&nbsp;GB per container,
          one shared Postgres, the same 10,000 rows, every framework on its own ORM.{' '}
          <span className="text-muted-foreground/70">
            Ratios are only comparable <em>within</em> a pair &mdash; the two figures under each
            bar are that pair&rsquo;s own measurements. Bar heights use a square-root scale so a
            42&times; result does not flatten a 2&times; one into an unreadable stub; the printed
            figures are the real ratios.
          </span>
        </p>
        <Link
          href="/docs/benchmarks"
          className="shrink-0 inline-flex items-center gap-1 text-sm text-primary hover:underline"
        >
          Methodology
          <ArrowRight className="h-3.5 w-3.5" />
        </Link>
      </div>
    </div>
  )
}
