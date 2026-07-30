'use client'

import { useState } from 'react'
import {
  Bot,
  Brain,
  FileText,
  Globe,
  Hourglass,
  Play,
  ShieldCheck,
  Sparkles,
  Target,
} from 'lucide-react'

/**
 * Split feature section: an accordion on the left, a product card on the right.
 *
 * Two independent interactions on one section — the accordion and the card's own
 * tabs — so each keeps its own state. Tying them together would mean opening
 * "Global Reach" silently changed what the card was showing, which is a
 * surprise, not a feature.
 */

const FEATURES = [
  {
    icon: Brain,
    name: 'AI Models',
    body: 'Switch between GPT, Claude, and Gemini from a unified interface.',
  },
  {
    icon: Globe,
    name: 'Global Reach',
    body: 'Deploy to fourteen regions with automatic failover and no config change.',
  },
  {
    icon: Bot,
    name: 'Smart Agent',
    body: 'Delegate multi-step work and watch each tool call as it happens.',
  },
]

const COMPLIANCE = [
  { icon: ShieldCheck, label: 'SOC 2' },
  { icon: ShieldCheck, label: 'ISO 27001' },
  { icon: ShieldCheck, label: 'GDPR' },
  { icon: Hourglass, label: '99.9%', muted: 'uptime' },
]

const TABS = [
  { key: 'summary', label: 'Summary', icon: Target },
  { key: 'transcript', label: 'Transcript', icon: FileText },
  { key: 'creations', label: 'Creations', icon: Sparkles },
] as const

const PANELS: Record<string, { intro: string; points: [string, string][]; tail: string }> = {
  summary: {
    intro: "Key decisions from today's marketing sync:",
    points: [
      ['Launch date confirmed:', 'Monday, March 18th at 9 AM EST'],
      ['Hero section:', 'New animated product showcase with A/B testing enabled'],
      ['SEO improvements:', 'Meta tags updated, sitemap regenerated'],
    ],
    tail: 'Next steps: QA review by Friday, staging deployment Thursday EOD.',
  },
  transcript: {
    intro: 'Transcript, speaker-separated:',
    points: [
      ['Priya:', 'Can we hold the launch until the analytics dashboard is in?'],
      ['Marcus:', 'The dashboard is decoupled — it can ship the week after.'],
      ['Priya:', 'Then Monday works. I will confirm with the content team.'],
    ],
    tail: 'Recording retained for 90 days, then deleted automatically.',
  },
  creations: {
    intro: 'Generated from this meeting:',
    points: [
      ['Launch checklist:', '14 items, 3 assigned to you'],
      ['Follow-up email:', 'Drafted for the content team, awaiting review'],
      ['Changelog entry:', 'Written from the decisions above'],
    ],
    tail: 'Every artefact links back to the moment it came from.',
  },
}

export default function ExpandableWithTabbedCard({
  title = 'Ship with\nconfidence using\nour unified platform',
  subtitle = 'Streamline your workflow with tools designed to enhance productivity at every step.',
  ctaLabel = 'Learn more',
  quote = 'Looks really good. Did you design in code or Figma first?',
  quoteName = 'Shadcn',
  quoteRole = 'Creator of Shadcn UI',
}: {
  title?: string
  subtitle?: string
  ctaLabel?: string
  quote?: string
  quoteName?: string
  quoteRole?: string
}) {
  const [open, setOpen] = useState(FEATURES[0].name)
  const [tab, setTab] = useState<string>(TABS[0].key)
  const panel = PANELS[tab]

  return (
    <section className="bg-white py-20 sm:py-28 dark:bg-gray-950">
      <div className="mx-auto grid max-w-7xl gap-16 px-6 lg:grid-cols-2 lg:gap-20">
        {/* ── Left column ───────────────────────────────────────────────── */}
        <div>
          <h2 className="text-[2rem] leading-[1.12] font-semibold tracking-[-0.025em] whitespace-pre-line text-gray-900 sm:text-[2.75rem] dark:text-white">
            {title}
          </h2>
          <p className="mt-5 max-w-md text-[16.5px]/[1.6] text-gray-500 dark:text-gray-400">
            {subtitle}
          </p>

          <a
            href="#"
            className="mt-8 inline-flex h-11 items-center rounded-xl border border-gray-200 bg-white px-5 text-[14.5px] font-semibold text-gray-900 shadow-[0_1px_2px_rgb(15_23_42_/_0.06)] transition-all duration-200 hover:bg-gray-50 active:scale-[0.98] dark:border-white/15 dark:bg-white/5 dark:text-white dark:hover:bg-white/10"
          >
            {ctaLabel}
          </a>

          {/* Accordion. One panel at a time: the point of the pattern is that the
              section stays a fixed, scannable height. */}
          <div className="mt-14 divide-y divide-gray-200 border-t border-gray-200 dark:divide-white/10 dark:border-white/10">
            {FEATURES.map((feature) => {
              const isOpen = open === feature.name
              return (
                <div key={feature.name}>
                  <h3>
                    <button
                      type="button"
                      onClick={() => setOpen(feature.name)}
                      aria-expanded={isOpen}
                      className="group flex w-full items-center gap-3 py-5 text-left"
                    >
                      <feature.icon
                        aria-hidden="true"
                        className={`size-[19px] shrink-0 transition-colors duration-200 ${
                          isOpen
                            ? 'text-gray-900 dark:text-white'
                            : 'text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300'
                        }`}
                      />
                      <span
                        className={`text-[15.5px] font-semibold transition-colors duration-200 ${
                          isOpen
                            ? 'text-gray-900 dark:text-white'
                            : 'text-gray-600 group-hover:text-gray-900 dark:text-gray-400 dark:group-hover:text-white'
                        }`}
                      >
                        {feature.name}
                      </span>
                    </button>
                  </h3>
                  {/* grid-rows trick: animates to the content's real height without
                      hard-coding a max-height that clips longer copy. */}
                  <div
                    className={`grid transition-all duration-300 ease-out ${
                      isOpen
                        ? 'grid-rows-[1fr] opacity-100'
                        : 'grid-rows-[0fr] opacity-0'
                    }`}
                  >
                    <div className="overflow-hidden">
                      <p className="pb-6 pl-[31px] text-[15px]/[1.6] text-gray-500 dark:text-gray-400">
                        {feature.body}
                      </p>
                    </div>
                  </div>
                </div>
              )
            })}
          </div>

          <ul className="mt-12 space-y-3">
            {COMPLIANCE.map((item) => (
              <li key={item.label} className="flex items-center gap-3 text-[15px]">
                <item.icon aria-hidden="true" className="size-[18px] text-gray-400" />
                <span className="font-semibold text-gray-900 dark:text-white">
                  {item.label}
                </span>
                {item.muted && (
                  <span className="text-gray-400 dark:text-gray-500">{item.muted}</span>
                )}
              </li>
            ))}
          </ul>
        </div>

        {/* ── Right column ──────────────────────────────────────────────── */}
        <div>
          {/* A gradient canvas, not a photograph — see the note in carousel-two-up
              for why this block ships no remote images. */}
          <div className="relative aspect-[4/4.4] overflow-hidden rounded-3xl bg-[linear-gradient(180deg,#93b7d8_0%,#a9c4dd_28%,#c9d3d0_58%,#8fa584_78%,#5f7a55_100%)] dark:bg-[linear-gradient(180deg,#16283a_0%,#17293a_28%,#1b2b31_58%,#1e2c22_78%,#16210f_100%)]">
            <div
              aria-hidden
              className="pointer-events-none absolute inset-0"
              style={{
                background:
                  'radial-gradient(52% 30% at 18% 96%, rgb(28 56 26 / 0.5), transparent 70%), radial-gradient(46% 26% at 84% 98%, rgb(24 48 22 / 0.45), transparent 70%)',
              }}
            />

            <div className="absolute inset-0 flex items-center justify-center p-6">
              <div className="w-full max-w-[22rem] rounded-2xl bg-white/95 p-5 shadow-[0_24px_60px_-20px_rgb(15_23_42_/_0.45)] backdrop-blur-sm dark:bg-gray-900/95">
                <p className="text-[12.5px] text-gray-500 dark:text-gray-400">
                  Today 09:15 AM
                </p>
                <p className="mt-1 text-[15.5px] font-semibold text-gray-900 dark:text-white">
                  Marketing Website Launch
                </p>

                <span className="mt-3 inline-flex items-center gap-2 rounded-full bg-sky-50 px-3 py-1.5 text-[12.5px] font-semibold text-gray-800 ring-1 ring-sky-200 dark:bg-sky-500/10 dark:text-sky-100 dark:ring-sky-500/30">
                  <Play aria-hidden="true" className="size-3 fill-current" />
                  03:47
                </span>

                <div
                  role="tablist"
                  aria-label="Meeting output"
                  className="mt-4 flex gap-4 border-b border-gray-200 dark:border-white/10"
                >
                  {TABS.map((t) => {
                    const active = tab === t.key
                    return (
                      <button
                        key={t.key}
                        type="button"
                        role="tab"
                        aria-selected={active}
                        onClick={() => setTab(t.key)}
                        className={`-mb-px inline-flex items-center gap-1.5 border-b-2 pb-2.5 text-[13px] font-medium transition-colors duration-200 ${
                          active
                            ? 'border-gray-900 text-gray-900 dark:border-white dark:text-white'
                            : 'border-transparent text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'
                        }`}
                      >
                        <t.icon aria-hidden="true" className="size-3.5" />
                        {t.label}
                      </button>
                    )
                  })}
                </div>

                <div role="tabpanel" className="mt-3.5">
                  <p className="text-[13px] text-gray-600 dark:text-gray-400">
                    {panel.intro}
                  </p>
                  <ul className="mt-2.5 space-y-2">
                    {panel.points.map(([lead, rest]) => (
                      <li
                        key={lead}
                        className="flex gap-2 text-[13px]/[1.5] text-gray-500 dark:text-gray-400"
                      >
                        <span aria-hidden className="text-gray-300 dark:text-gray-600">
                          •
                        </span>
                        <span>
                          <span className="font-semibold text-gray-900 dark:text-white">
                            {lead}
                          </span>{' '}
                          {rest}
                        </span>
                      </li>
                    ))}
                  </ul>
                  {/* Faded rather than cut: the card is a glimpse of a longer
                      document, and a hard clip looks like a layout bug. */}
                  <p className="mt-3 text-[13px]/[1.5] text-gray-400/70 dark:text-gray-500/70">
                    {panel.tail}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <figure className="mt-12">
            <blockquote className="text-[17px]/[1.55] text-gray-900 dark:text-white">
              &ldquo;{quote}&rdquo;
            </blockquote>
            <figcaption className="mt-4 flex items-center gap-3">
              {/* Initials, not an avatar image. Same reason as the canvas. */}
              <span
                aria-hidden
                className="flex size-10 items-center justify-center rounded-xl bg-gradient-to-br from-violet-400 to-sky-400 text-[13px] font-bold text-white"
              >
                {quoteName.slice(0, 2).toUpperCase()}
              </span>
              <span>
                <span className="block text-[14.5px] font-semibold text-gray-900 dark:text-white">
                  {quoteName}
                </span>
                <span className="block text-[13.5px] text-gray-500 dark:text-gray-400">
                  {quoteRole}
                </span>
              </span>
            </figcaption>
          </figure>
        </div>
      </div>
    </section>
  )
}
