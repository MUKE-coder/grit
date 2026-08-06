/*
 * Three tiles above a wide pair: copy first, artifact underneath.
 *
 * A bento grid works because not every feature deserves the same box. The
 * spans here are asymmetric on purpose — the wide tile at the bottom carries
 * the feature you actually want read, and the three above it are supporting.
 * A grid of six equal boxes is a list with extra borders.
 *
 * The spans only exist above `lg`. Below that everything is one column in
 * source order, so put the tiles in the order you want them read: on a phone
 * that order is the design.
 *
 * Every artifact is markup and `aria-hidden`. They illustrate the sentence
 * above them and contain invented data — a screen reader reading out a fake
 * spending limit between two feature descriptions is noise, not information.
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
      <div className="mt-3 flex gap-2 text-gray-400 dark:text-gray-600">
        {['@', '☺', '🖇'].map((glyph) => (
          <span key={glyph} className="text-xs">
            {glyph}
          </span>
        ))}
      </div>
    </Card>
  )
}

function TimelineArtifact() {
  const events = [
    { time: '06 AM', label: 'Poll created', highlight: false },
    { time: '12 PM', label: '+50 users voted', highlight: true },
    { time: '12:30 PM', label: 'Poll closed', highlight: false },
  ]
  return (
    <div className="space-y-2">
      {events.map((event) => (
        <div
          key={event.label}
          className={`flex items-center gap-3 rounded-lg px-3 py-2 ${
            event.highlight
              ? 'border border-gray-200 bg-white shadow-sm dark:border-white/10 dark:bg-gray-900'
              : ''
          }`}
        >
          <span className="size-1.5 flex-none rounded-full border border-gray-400 dark:border-gray-600" />
          <span className="w-16 flex-none text-[10px] text-gray-500 tabular-nums dark:text-gray-400">
            {event.time}
          </span>
          <span className="text-xs font-semibold text-gray-900 dark:text-white">{event.label}</span>
        </div>
      ))}
    </div>
  )
}

function SpendingArtifact() {
  return (
    <Card>
      <p className="text-xs font-semibold text-gray-900 dark:text-white">
        <span className="rounded bg-amber-100 px-1 dark:bg-amber-500/20">Spending</span> limit
      </p>
      <p className="mt-1 text-[11px] text-gray-500 dark:text-gray-400">
        Usage by primary channel group
      </p>
      <div className="mt-4 flex h-2 overflow-hidden rounded-full bg-gray-100 dark:bg-white/10">
        <span className="block w-[22%] bg-gray-900 dark:bg-white" />
        <span className="block w-[18%] bg-indigo-500" />
      </div>
      <div className="mt-3 flex gap-10">
        {[
          { value: '40%', label: 'Used' },
          { value: '60%', label: 'Free' },
        ].map((item) => (
          <div key={item.label}>
            <p className="text-lg font-semibold text-gray-900 tabular-nums dark:text-white">
              {item.value}
            </p>
            <p className="text-[10px] text-gray-500 dark:text-gray-400">{item.label}</p>
          </div>
        ))}
      </div>
    </Card>
  )
}

export interface Tile {
  title: string
  body: string
  Artifact: () => React.JSX.Element
  /** Column span above lg. */
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
    title: 'Activity timeline',
    body: 'Every state change recorded against the record that changed.',
    Artifact: TimelineArtifact,
    span: 'lg:col-span-1',
  },
  {
    title: 'Budgets and limits',
    body: 'Set a ceiling per team and watch it in real time rather than at the end of the month.',
    Artifact: SpendingArtifact,
    span: 'lg:col-span-2',
  },
]

export default function BentoThreeUpWithFeatureRow({ tiles = TILES }: { tiles?: Tile[] }) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        {/* One column below lg, in source order. On a phone that order is the
            design, so the array order is the reading order. */}
        <div className="grid grid-cols-1 gap-4 lg:grid-cols-3">
          {tiles.map((tile) => (
            <div
              key={tile.title}
              className={`flex flex-col rounded-2xl border border-gray-200 p-6 dark:border-white/10 ${
                tile.span ?? ''
              }`}
            >
              <h3 className="text-sm font-semibold text-gray-900 dark:text-white">{tile.title}</h3>
              <p className="mt-2 text-sm/6 text-gray-600 dark:text-gray-400">{tile.body}</p>
              <div aria-hidden="true" className="mt-8 select-none">
                <tile.Artifact />
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
