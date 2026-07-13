// StatBars — a compact horizontal bar chart for benchmarks and before/after stats.
// Server component, theme-aware. Bars are scaled to the largest value in the set.
//
//   <StatBars
//     title="Typical JSON list payload"
//     items={[
//       { label: 'Uncompressed', value: 100, display: '100 KB', tone: 'default' },
//       { label: 'Gzip (BestSpeed)', value: 28, display: '~28 KB', tone: 'primary' },
//     ]}
//     caption="60–80% smaller responses, applied automatically to every route"
//   />

type BarTone = 'primary' | 'blue' | 'green' | 'amber' | 'rose' | 'cyan' | 'violet' | 'default'

const TONES: Record<BarTone, string> = {
  primary: 'var(--primary)',
  blue: '217 91% 60%',
  green: '160 84% 39%',
  amber: '38 92% 50%',
  rose: '350 89% 60%',
  cyan: '189 94% 43%',
  violet: '258 90% 66%',
  default: '215 16% 55%',
}

const hsl = (tone: BarTone, alpha = 1) =>
  `hsl(${TONES[tone] ?? TONES.default}${alpha === 1 ? '' : ` / ${alpha}`})`

export interface StatBarItem {
  label: string
  /** Numeric magnitude — bars are scaled to the largest in the set. */
  value: number
  /** Text shown at the end of the bar (defaults to the value). */
  display?: string
  tone?: BarTone
}

export function StatBars({
  title,
  items,
  caption,
}: {
  title?: string
  items: StatBarItem[]
  caption?: string
}) {
  const max = Math.max(1, ...items.map((i) => i.value))
  return (
    <figure className="my-8 not-prose overflow-hidden rounded-xl border border-border/60 bg-card/40 p-5">
      {title && (
        <div className="mb-4 text-sm font-semibold text-foreground/90">{title}</div>
      )}
      <div className="space-y-3">
        {items.map((item, i) => {
          const tone = item.tone ?? 'default'
          const pct = Math.max(2, (item.value / max) * 100)
          return (
            <div key={i} className="flex items-center gap-3">
              <div className="w-32 shrink-0 text-right text-[12px] text-muted-foreground">
                {item.label}
              </div>
              <div className="relative h-6 flex-1 overflow-hidden rounded-md bg-muted/40">
                <div
                  className="flex h-full items-center justify-end rounded-md pr-2"
                  style={{
                    width: `${pct}%`,
                    background: `linear-gradient(90deg, ${hsl(tone, 0.35)}, ${hsl(tone, 0.85)})`,
                  }}
                >
                  <span className="text-[11px] font-semibold text-foreground/90">
                    {item.display ?? item.value}
                  </span>
                </div>
              </div>
            </div>
          )
        })}
      </div>
      {caption && (
        <figcaption className="mt-4 border-t border-border/30 pt-3 text-xs italic text-muted-foreground/80">
          {caption}
        </figcaption>
      )}
    </figure>
  )
}
