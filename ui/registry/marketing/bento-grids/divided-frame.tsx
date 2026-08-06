/*
 * Six features in one frame, separated by hairlines rather than gaps.
 *
 * The gapless version reads as a single surface with regions, where the
 * card version reads as six objects. Use this for capabilities that belong to
 * one product and the card version for things a visitor picks between.
 *
 * The dividers come from a background colour on the grid showing through 1px
 * gaps, which is the one reliable way to get a single hairline between cells
 * whose spans differ per breakpoint. `divide-x` cannot do it once a cell spans
 * two columns: the divider follows the DOM order, not the visual grid, so it
 * lands in the middle of the wide cell.
 *
 * The spans have to tile the grid exactly. Here that is 1+1+1, then 2+1, then
 * a full-width 3 — nine units across three columns. Get it wrong and the last
 * row leaves a hole, and because the hairlines are the frame's background
 * showing through, that hole is a grey rectangle rather than nothing. Change
 * the number of tiles and re-check the arithmetic.
 *
 * All artifacts are markup and `aria-hidden` — they illustrate the sentence
 * above them and contain invented data.
 */

import type { ReactNode } from 'react'

function Card({ children }: { children: ReactNode }) {
  return (
    <div className="rounded-lg border border-gray-200 bg-white p-3 shadow-sm dark:border-white/10 dark:bg-gray-900">
      {children}
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

function DocumentsArtifact() {
  return (
    <div className="flex gap-2">
      {[0, 1, 2].map((i) => (
        <div
          key={i}
          className="w-14 rounded-md border border-gray-200 bg-white p-2 shadow-sm dark:border-white/10 dark:bg-gray-900"
        >
          <span className="block size-1.5 rounded-full bg-gray-300 dark:bg-white/20" />
          <div className="mt-2 space-y-1">
            <span className="block h-1 w-full rounded-full bg-gray-200 dark:bg-white/10" />
            <span className="block h-1 w-full rounded-full bg-gray-200 dark:bg-white/10" />
            <span className="block h-1 w-2/3 rounded-full bg-gray-200 dark:bg-white/10" />
          </div>
        </div>
      ))}
    </div>
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
          className={`${card.rotate} w-14 rounded-lg border border-gray-200 p-2 shadow-sm dark:border-white/10 ${card.tone}`}
        >
          <p className="font-mono text-[9px] font-semibold">{card.label}</p>
          <div className="mt-2 space-y-1">
            <span className="block h-1 w-full rounded-full bg-current opacity-25" />
          </div>
        </div>
      ))}
    </div>
  )
}

function ChatArtifact() {
  return (
    <div className="space-y-2">
      <p className="text-[10px] text-gray-400 dark:text-gray-600">Sat 22 Feb</p>
      <div className="max-w-[80%] rounded-lg border border-gray-200 bg-white px-3 py-2 text-xs text-gray-700 dark:border-white/10 dark:bg-gray-900 dark:text-gray-300">
        I am having trouble with my account.
      </div>
      <div className="ml-auto max-w-[80%] rounded-lg bg-indigo-600 px-3 py-2 text-xs text-white">
        Let me pull that up, one moment.
      </div>
    </div>
  )
}

function ScheduleArtifact() {
  return (
    <Card>
      <div className="flex items-center gap-2">
        <span className="rounded-md bg-indigo-600 px-2 py-1 text-[10px] font-semibold text-white">
          Schedule
        </span>
        <div className="flex gap-1.5 text-[11px] font-medium text-gray-400 dark:text-gray-600">
          <span className="font-bold">B</span>
          <span className="italic">I</span>
          <span className="underline">U</span>
        </div>
      </div>
      <p className="mt-3 text-xs text-gray-600 dark:text-gray-400">
        <span className="rounded bg-indigo-50 px-1 text-indigo-700 dark:bg-indigo-500/15 dark:text-indigo-300">
          Tomorrow 8:30 pm
        </span>{' '}
        is our priority.
      </p>
    </Card>
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
        </div>
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
    title: 'Real-time collaboration',
    body: 'Work together on the same record. Edits, comments and presence arrive as they happen rather than on refresh.',
    Artifact: MentionArtifact,
  },
  {
    title: 'Document management',
    body: 'Organise and find files with categorisation that happens on upload instead of never.',
    Artifact: DocumentsArtifact,
  },
  {
    title: 'Financial analytics',
    body: 'Reports and conversions across every currency you trade in.',
    Artifact: CurrencyArtifact,
  },
  {
    title: 'Assisted support',
    body: 'Answer the common questions automatically and hand the rest to a person with the context attached.',
    Artifact: ChatArtifact,
    span: 'lg:col-span-2',
  },
  {
    title: 'Automated scheduling',
    body: 'Find the meeting time that works without the six-message thread.',
    Artifact: ScheduleArtifact,
  },
  {
    title: 'Secure file sharing',
    body: 'Signed links with an expiry and a permission model you can explain to somebody in one sentence.',
    Artifact: FileArtifact,
    span: 'lg:col-span-3',
  },
]

export default function BentoDividedFrame({ tiles = TILES }: { tiles?: Tile[] }) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        {/* The 1px gap plus a background colour is what draws the hairlines.
            divide-x cannot do this once a cell spans two columns: it follows
            the DOM order, so the rule lands inside the wide cell. */}
        <div className="overflow-hidden rounded-2xl border border-gray-200 bg-gray-200 dark:border-white/10 dark:bg-white/10">
          <div className="grid grid-cols-1 gap-px lg:grid-cols-3">
            {tiles.map((tile) => (
              <div
                key={tile.title}
                className={`flex flex-col bg-white p-8 dark:bg-gray-950 ${tile.span ?? ''}`}
              >
                <h3 className="text-base font-semibold text-gray-900 dark:text-white">
                  {tile.title}
                </h3>
                <p className="mt-2 text-sm/6 text-gray-600 dark:text-gray-400">{tile.body}</p>
                <div aria-hidden="true" className="mt-auto pt-10 select-none">
                  <tile.Artifact />
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}
