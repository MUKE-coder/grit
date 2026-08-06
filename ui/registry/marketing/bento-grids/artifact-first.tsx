/*
 * The same tiles with the artifact on top and the copy underneath.
 *
 * Worth having as a separate block rather than a prop, because the order
 * changes what the section is for. Copy first is for features that need
 * explaining; artifact first is for a product people already understand, where
 * the picture is the hook and the sentence is the caption.
 *
 * The artifacts sit in a fixed-height box so the copy below them starts on the
 * same line across the row. Let each tile size itself to its own artifact and
 * the headings step up and down like a bar chart.
 *
 * Everything decorative is `aria-hidden`. These illustrate the sentence
 * beneath them and hold invented data.
 */

import type { ReactNode } from 'react'

function Card({ children }: { children: ReactNode }) {
  return (
    <div className="rounded-lg border border-gray-200 bg-white p-3 shadow-sm dark:border-white/10 dark:bg-gray-900">
      {children}
    </div>
  )
}

function FileArtifact() {
  return (
    <Card>
      <div className="flex items-center gap-3">
        <span className="rounded bg-red-500 px-1.5 py-1 text-[9px] font-bold text-white">PDF</span>
        <div className="min-w-0 flex-1">
          <p className="truncate text-xs font-medium text-gray-900 dark:text-white">
            annual-report.pdf
          </p>
          <div className="mt-1.5 h-1 overflow-hidden rounded-full bg-gray-100 dark:bg-white/10">
            <span className="block h-full w-1/4 rounded-full bg-indigo-500" />
          </div>
          <p className="mt-1 text-[10px] text-gray-500 dark:text-gray-400">29 KB / 120 KB</p>
        </div>
      </div>
    </Card>
  )
}

function CurrencyArtifact() {
  const cards = [
    { label: '₿ BTC', tone: 'bg-blue-50 text-blue-700 dark:bg-blue-500/10 dark:text-blue-300', rotate: '-rotate-6' },
    { label: '$ USD', tone: 'bg-emerald-50 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-300', rotate: '' },
    { label: '€ EUR', tone: 'bg-rose-50 text-rose-700 dark:bg-rose-500/10 dark:text-rose-300', rotate: 'rotate-6' },
  ]
  return (
    <div className="flex items-end -space-x-3">
      {cards.map((card) => (
        <div
          key={card.label}
          className={`${card.rotate} w-16 rounded-lg border border-gray-200 p-2 shadow-sm dark:border-white/10 ${card.tone}`}
        >
          <p className="font-mono text-[10px] font-semibold">{card.label}</p>
          <div className="mt-2 space-y-1">
            <span className="block h-1 w-full rounded-full bg-current opacity-25" />
            <span className="block h-1 w-2/3 rounded-full bg-current opacity-25" />
          </div>
        </div>
      ))}
    </div>
  )
}

function MentionArtifact() {
  return (
    <Card>
      <p className="text-xs text-gray-700 dark:text-gray-300">
        <span className="font-medium text-indigo-600 dark:text-indigo-400">@bernard</span> shared 2
        invoices
      </p>
      <div className="mt-3 flex gap-2 text-[11px] text-gray-400 dark:text-gray-600">
        <span>@</span>
        <span>☺</span>
        <span>🖇</span>
      </div>
    </Card>
  )
}

function ChartArtifact() {
  const bars = [42, 68, 51, 84, 62, 92]
  return (
    <Card>
      <p className="text-xs font-semibold text-gray-900 dark:text-white">Monthly volume</p>
      <div className="mt-4 flex h-16 items-end gap-1.5">
        {bars.map((height, i) => (
          <span
            key={i}
            className="flex-1 rounded-t-sm bg-indigo-500/70 dark:bg-indigo-400/60"
            style={{ height: `${height}%` }}
          />
        ))}
      </div>
    </Card>
  )
}

function MapArtifact() {
  /* A handful of pins rather than a whole map: at this size a world map is a
     grey smudge, and three markers say "everywhere" just as well. */
  const pins = [
    { left: '18%', top: '32%' },
    { left: '46%', top: '54%' },
    { left: '72%', top: '28%' },
  ]
  return (
    <Card>
      <div className="relative h-24 overflow-hidden rounded-md bg-gray-50 dark:bg-white/5">
        <div
          className="absolute inset-0 opacity-60"
          style={{
            backgroundImage: 'radial-gradient(currentColor 1px, transparent 1px)',
            backgroundSize: '8px 8px',
            color: 'rgb(156 163 175 / 0.6)',
          }}
        />
        {pins.map((pin) => (
          <span
            key={pin.left}
            className="absolute size-3 -translate-x-1/2 -translate-y-1/2 rounded-full border-2 border-white bg-indigo-500 shadow dark:border-gray-900"
            style={{ left: pin.left, top: pin.top }}
          />
        ))}
      </div>
    </Card>
  )
}

export interface Tile {
  title: string
  body: string
  Artifact: () => React.JSX.Element
  span?: string
}

const TILES: Tile[] = [
  {
    title: 'Scheduled reports',
    body: 'Automate delivery to stakeholders on whatever schedule they actually read.',
    Artifact: FileArtifact,
  },
  {
    title: 'Multi-currency',
    body: 'Convert, store and report in the currency each customer pays in.',
    Artifact: CurrencyArtifact,
  },
  {
    title: 'Collaborative analysis',
    body: 'Comment, share and work through the numbers without leaving the dashboard.',
    Artifact: MentionArtifact,
  },
  {
    title: 'Interactive dashboards',
    body: 'Build a view per team by dragging panels around. The arrangement is data, so it survives a redeploy.',
    Artifact: ChartArtifact,
    span: 'lg:col-span-2',
  },
  {
    title: 'Regional reporting',
    body: 'Break every figure down by the region it came from.',
    Artifact: MapArtifact,
  },
]

export default function BentoArtifactFirst({ tiles = TILES }: { tiles?: Tile[] }) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="grid grid-cols-1 gap-4 lg:grid-cols-3">
          {tiles.map((tile) => (
            <div
              key={tile.title}
              className={`rounded-2xl border border-gray-200 p-6 dark:border-white/10 ${
                tile.span ?? ''
              }`}
            >
              {/* Fixed height, so the headings below line up across the row. */}
              <div
                aria-hidden="true"
                className="flex min-h-32 items-center select-none"
              >
                <div className="w-full">
                  <tile.Artifact />
                </div>
              </div>
              <h3 className="mt-6 text-sm font-semibold text-gray-900 dark:text-white">
                {tile.title}
              </h3>
              <p className="mt-2 text-sm/6 text-gray-600 dark:text-gray-400">{tile.body}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
