'use client'

import { useState } from 'react'
import Image from 'next/image'
import Link from 'next/link'
import { ArrowRight } from 'lucide-react'
import { FrameworkLogo } from '@/components/framework-logo'

/*
 * The homepage benchmark chart: one group per head-to-head pair, Grit's bar
 * beside its opponent's, both labelled with the requests per second actually
 * measured.
 *
 * Bar heights are normalised WITHIN each pair, and that is a correctness
 * decision rather than a stylistic one. Every framework was measured against
 * Grit in its own back-to-back run, and across those runs the machine drifts a
 * long way — Grit's single-row read measured 4,536 req/s in the Bun pair, 4,392
 * in the Encore pair, 1,911 in the Next.js pair and 1,635 in the Express pair,
 * from an identical binary, as hours of write scenarios accumulated in Postgres.
 * Putting every bar on one shared axis would say those were the same Grit. They were not,
 * and Encore would look far worse than it is.
 *
 * So the taller bar in each pair is full height and the other is drawn in
 * proportion to it. Compare the two bars inside a group; do not compare heights
 * across groups. The printed req/s figures are the real measurements either way.
 */

export type BenchRow = {
  framework: string
  slug: string
  version: string
  logo: string
  color: string
  invertOnDark?: boolean
  /** how many times faster Grit was in this pair; below 1 means Grit lost */
  ratio: number
  /** what Grit measured in this pair */
  gritRps: number
  /** what the opponent measured in the same pair */
  rps: number
  /** true when the opponent's own container saturated — a real ceiling */
  appBound: boolean
}

export type BenchScenario = {
  id: string
  label: string
  title: string
  subtitle: string
  rows: BenchRow[]
}

const GRIT_COLOR = '#3BB4F5'

function Bar({
  value,
  pct,
  color,
  label,
  faded,
}: {
  value: number
  pct: number
  color: string
  label: string
  faded?: boolean
}) {
  return (
    <div className="flex h-full flex-1 flex-col justify-end items-center min-w-0">
      <span
        className="mb-1.5 text-[11px] sm:text-xs font-bold tabular-nums whitespace-nowrap"
        style={{ color }}
      >
        {value.toLocaleString()}
      </span>
      {/* Top-lit rather than bottom-lit: fading toward the base keeps the pale
          brand colours (Bun's and Encore's creams) from reading as grey. */}
      <div
        className="w-full rounded-t-[5px] transition-all duration-500"
        style={{
          height: `${pct}%`,
          background: `linear-gradient(to bottom, ${color}, ${color}bb)`,
          opacity: faded ? 0.95 : 1,
        }}
      />
      <span className="mt-1.5 text-[10px] text-muted-foreground/70 whitespace-nowrap">
        {label}
      </span>
    </div>
  )
}

export function BenchmarkChart({ scenarios }: { scenarios: BenchScenario[] }) {
  const [active, setActive] = useState(scenarios[0]?.id)
  const scenario = scenarios.find((s) => s.id === active) ?? scenarios[0]
  if (!scenario) return null

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
          Requests per second, {scenario.title}
        </h3>
        <p className="text-sm text-muted-foreground mt-1">{scenario.subtitle}</p>
      </div>

      <div className="flex items-center justify-center gap-4 mb-8 text-xs text-muted-foreground">
        <span className="inline-flex items-center gap-1.5">
          <Image
            src="/logos/grit.png"
            alt=""
            width={16}
            height={16}
            className="rounded-[4px]"
          />
          <span className="font-medium text-foreground">Grit</span>
        </span>
        <span className="text-muted-foreground/50">vs</span>
        <span>each framework in its own colour</span>
      </div>

      {/* Scrolls sideways on a phone rather than crushing four groups into 320px. */}
      <div className="overflow-x-auto -mx-2 px-2 pb-1">
        <div
          className="flex items-end justify-center gap-4 sm:gap-7"
          style={{ minWidth: scenario.rows.length * 126 }}
        >
          {scenario.rows.map((row) => {
            // Normalised inside the pair: the winner is full height, the other
            // is drawn against it. Never against another pair's numbers.
            const peak = Math.max(row.gritRps, row.rps)
            const behind = row.ratio < 1

            return (
              <Link
                key={row.slug}
                href={`/docs/benchmarks/${row.slug}`}
                className="group flex flex-1 flex-col items-center min-w-[112px] max-w-[190px]"
                title={`Grit ${row.gritRps.toLocaleString()} req/s vs ${row.framework} ${row.rps.toLocaleString()} req/s, measured back to back`}
              >
                {/* A fixed 44px box, so a square mark and a wide wordmark carry
                    the same visual weight instead of the wordmark shrinking. */}
                <div className="h-11 flex items-center justify-center mb-3.5 opacity-95 transition-opacity group-hover:opacity-100">
                  <FrameworkLogo
                    src={row.logo}
                    alt={row.framework}
                    onLight={row.invertOnDark}
                    className="h-11 max-w-[140px]"
                  />
                </div>

                <div className="flex h-[230px] w-full items-end justify-center gap-1.5 px-1">
                  <Bar
                    value={row.gritRps}
                    pct={(row.gritRps / peak) * 100}
                    color={GRIT_COLOR}
                    label="Grit"
                  />
                  <Bar
                    value={row.rps}
                    pct={Math.max(1.5, (row.rps / peak) * 100)}
                    color={row.color}
                    label={row.framework}
                    faded
                  />
                </div>

                <div className="mt-3 text-center px-0.5">
                  <div
                    className={
                      'text-sm font-bold tabular-nums ' +
                      (behind ? 'text-[#f5a623]' : 'text-primary')
                    }
                  >
                    {behind
                      ? `${(1 / row.ratio).toFixed(2)}× slower`
                      : `${row.ratio.toFixed(2)}× faster`}
                  </div>
                  <div className="text-[11px] font-medium leading-tight text-muted-foreground group-hover:text-primary transition-colors mt-1">
                    vs {row.framework}
                  </div>
                  <div className="text-[10px] font-mono text-muted-foreground/60 leading-tight mt-0.5">
                    {row.version}
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
            Bar heights are scaled <em>within</em> each pair, because each pair is a separate run.
            Compare the two bars inside a group, not heights across groups. The printed
            figures are the real measurements.
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
