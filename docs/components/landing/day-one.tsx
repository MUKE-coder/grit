import Image from 'next/image'
import Link from 'next/link'
import { ArrowRight } from 'lucide-react'
import { COMPARED, readyCount, TOTAL_ROWS } from '@/config/capabilities'

/**
 * The half a benchmark cannot measure, as a summary rather than as the table.
 *
 * The full 33-row matrix used to sit on the home page. It is the most
 * defensible thing on the site and it was also 2,800 pixels of it, which meant
 * the two sections after it were read by almost nobody. The number and the
 * ranking are what a visitor needs here; the evidence lives one click away,
 * where somebody who doubts the number goes to check it.
 *
 * Every count comes from `readyCount`, so the figure in the heading and the
 * figure in the bars cannot disagree: they are the same function over the same
 * table that the full page renders.
 */
export function DayOne() {
  const grit = readyCount('grit')
  const ranked = [...COMPARED]
    .map((f) => ({ ...f, count: readyCount(f.slug) }))
    .sort((a, b) => b.count - a.count)

  return (
    <div className="grid gap-4 lg:grid-cols-12">
      <div className="g-bento lg:col-span-5">
        <p className="g-eyebrow">Running on day one</p>
        <p className="mt-6 flex items-end gap-2">
          <span className="g-hl text-[64px] font-medium leading-none tracking-[-0.04em]">
            {grit}
          </span>
          <span className="pb-2 text-[24px] leading-none text-fg-subtle">/ {TOTAL_ROWS}</span>
        </p>
        <p className="mt-6 text-[15px] leading-[1.7] text-fg-muted">
          Production capabilities running the first time you open the app: identity,
          security, compliance, operations and reach. The next closest framework has{' '}
          {ranked[1].count}.
        </p>
        <Link href="/docs/benchmarks#capabilities" className="g-link mt-auto pt-6 text-[14px]">
          See all {TOTAL_ROWS} rows
          <ArrowRight className="ml-1.5 inline h-3.5 w-3.5" aria-hidden="true" />
        </Link>
      </div>

      <div className="g-bento lg:col-span-7">
        <ol className="m-0 list-none space-y-4 p-0">
          {ranked.map((f) => (
            <li key={f.slug} className="flex items-center gap-4">
              <Image
                src={f.logo}
                alt=""
                width={20}
                height={20}
                className={`h-5 w-5 flex-none object-contain ${
                  'invertOnDark' in f && f.invertOnDark ? 'dark:invert' : ''
                }`}
              />
              <span
                className={`w-[86px] flex-none text-[14px] ${
                  f.slug === 'grit' ? 'font-semibold text-fg' : 'text-fg-muted'
                }`}
              >
                {f.name}
              </span>
              <span
                aria-hidden="true"
                className="h-2 flex-1 overflow-hidden rounded-full bg-surface-3"
              >
                <span
                  className={`block h-full rounded-full ${
                    f.slug === 'grit' ? 'bg-chart-grit' : 'bg-chart-rival'
                  }`}
                  style={{ width: `${Math.max((f.count / TOTAL_ROWS) * 100, 1)}%` }}
                />
              </span>
              <span
                className={`w-7 flex-none text-right font-mono text-[13px] ${
                  f.slug === 'grit' ? 'text-fg' : 'text-fg-subtle'
                }`}
              >
                {f.count}
              </span>
            </li>
          ))}
        </ol>
      </div>
    </div>
  )
}
