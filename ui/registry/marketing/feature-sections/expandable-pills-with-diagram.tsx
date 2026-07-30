'use client'

import { useState } from 'react'
import { ChevronDown, ChevronUp, PlusCircle } from 'lucide-react'

/**
 * A vertical stack of pill triggers where the open one becomes a card, beside a
 * line-art diagram.
 *
 * The stepper buttons and the pills drive the same single index, so the two
 * controls can never disagree about what is open. They also wrap, which is the
 * point of a stepper on a short list — disabling it at both ends would leave a
 * control that is dead half the time it is looked at.
 */

const PILLARS = [
  {
    name: 'Server',
    body: 'Run your API on managed instances that scale with traffic and roll back on a bad deploy.',
  },
  {
    name: 'Router',
    body: 'Work with teammates across the globe with real-time presence indicators and automatic syncing.',
  },
  {
    name: 'Database',
    body: 'Postgres with automated backups, point-in-time recovery, and read replicas in one click.',
  },
  {
    name: 'Tab',
    body: 'Pick up on any device exactly where you left off, down to the cursor position.',
  },
  {
    name: 'Mobile',
    body: 'Ship to iOS and Android from the same codebase, with offline support included.',
  },
]

/** Isometric line art, drawn rather than photographed so nothing loads remotely. */
function Diagram() {
  return (
    <svg
      viewBox="0 0 520 400"
      fill="none"
      aria-hidden="true"
      className="h-auto w-full text-gray-300 dark:text-white/15"
    >
      <g stroke="currentColor" strokeWidth="1.2">
        {/* Upper platform */}
        <path
          d="M120 150 L300 60 L470 150 L290 240 Z"
          strokeDasharray="4 5"
          className="opacity-70"
        />
        {/* Lower platform */}
        <path d="M60 290 L240 200 L410 290 L230 380 Z" strokeDasharray="4 5" className="opacity-50" />

        {/* Server stack */}
        <path d="M196 118 L246 92 L286 114 L236 141 Z" />
        <path d="M196 118 L196 146 L236 169 L236 141" />
        <path d="M236 141 L236 169 L286 142 L286 114" />
        <path d="M206 138 L206 152 M214 142 L214 156 M222 146 L222 160" className="opacity-60" />
        <circle cx="228" cy="106" r="3" />
        <circle cx="252" cy="118" r="3" />

        {/* Router puck */}
        <path d="M268 152 L318 126 L368 152 L318 179 Z" />
        <path d="M268 152 L268 166 L318 193 L368 166 L368 152" />
        <path d="M300 148 L336 148" className="opacity-60" />
        <circle cx="304" cy="160" r="2.5" />
        <circle cx="318" cy="153" r="2.5" />

        {/* Database drum */}
        <path d="M352 196 L392 173 L432 196 L392 219 Z" />
        <path d="M352 196 L352 232 L392 255 L432 232 L432 196" />
        <path d="M352 214 L392 237 L432 214" className="opacity-60" />
        <circle cx="366" cy="212" r="2.5" />
        <circle cx="366" cy="230" r="2.5" />

        {/* Tab / laptop */}
        <path d="M150 300 L228 255 L300 297 L222 342 Z" />
        <path d="M166 300 L228 264 L284 297 L222 332 Z" className="opacity-60" />

        {/* Phone */}
        <path d="M256 322 L302 296 L330 312 L284 338 Z" />
        <path d="M256 322 L256 330 L284 346 L330 320 L330 312" />
        <circle cx="288" cy="322" r="3" className="opacity-70" />

        {/* Connector, platform to platform */}
        <path d="M330 250 L262 290" strokeDasharray="5 5" />
        <path d="M262 290 L274 282 M262 290 L272 296" />
      </g>
    </svg>
  )
}

export default function ExpandablePillsWithDiagram({
  title = 'Ship with confidence\nusing our unified platform',
  subtitle = 'Five pillars that power your development workflow from idea to deployment.',
}: {
  title?: string
  subtitle?: string
}) {
  const [index, setIndex] = useState(1)

  // Wrapping, so both stepper buttons always do something.
  const step = (direction: 1 | -1) =>
    setIndex((i) => (i + direction + PILLARS.length) % PILLARS.length)

  return (
    <section className="bg-gray-50/70 py-20 sm:py-28 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6">
        <div className="grid gap-8 lg:grid-cols-2 lg:gap-16">
          <h2 className="text-[2rem] leading-[1.12] font-semibold tracking-[-0.025em] whitespace-pre-line text-gray-900 sm:text-[2.75rem] dark:text-white">
            {title}
          </h2>
          <p className="max-w-sm text-[16.5px]/[1.6] text-gray-500 lg:pt-2 dark:text-gray-400">
            {subtitle}
          </p>
        </div>

        <div className="mt-14 grid items-center gap-12 lg:grid-cols-2 lg:gap-8">
          {/* ── Pills ─────────────────────────────────────────────────────── */}
          <div className="flex gap-5">
            <div className="hidden shrink-0 flex-col justify-center gap-2 sm:flex">
              <StepButton label="Previous pillar" onClick={() => step(-1)}>
                <ChevronUp aria-hidden="true" className="size-4" />
              </StepButton>
              <StepButton label="Next pillar" onClick={() => step(1)}>
                <ChevronDown aria-hidden="true" className="size-4" />
              </StepButton>
            </div>

            <ul className="min-w-0 flex-1 space-y-2.5">
              {PILLARS.map((pillar, i) => {
                const isOpen = i === index
                return (
                  <li key={pillar.name}>
                    {isOpen ? (
                      // The open item is a card, not a pill with text under it.
                      // Changing the container is what makes the open state
                      // obvious at a glance in a list this uniform.
                      <div className="rounded-2xl border border-gray-200 bg-white p-5 shadow-[0_10px_30px_-12px_rgb(15_23_42_/_0.14)] dark:border-white/10 dark:bg-gray-900">
                        <p className="text-[15px]/[1.6] text-gray-500 dark:text-gray-400">
                          <span className="font-semibold text-gray-900 dark:text-white">
                            {pillar.name}.
                          </span>{' '}
                          {pillar.body}
                        </p>
                      </div>
                    ) : (
                      <button
                        type="button"
                        onClick={() => setIndex(i)}
                        aria-expanded={false}
                        className="group inline-flex items-center gap-2 rounded-full border border-gray-200 bg-white px-4 py-2.5 text-[15px] font-medium text-gray-800 shadow-[0_1px_2px_rgb(15_23_42_/_0.05)] transition-all duration-200 hover:border-gray-300 hover:bg-gray-50 active:scale-[0.98] dark:border-white/10 dark:bg-white/5 dark:text-gray-200 dark:hover:bg-white/10"
                      >
                        <PlusCircle
                          aria-hidden="true"
                          className="size-[17px] text-gray-400 transition-colors duration-200 group-hover:text-gray-600 dark:group-hover:text-gray-300"
                        />
                        {pillar.name}
                      </button>
                    )}
                  </li>
                )
              })}
            </ul>
          </div>

          <div className="order-first lg:order-none">
            <Diagram />
          </div>
        </div>
      </div>
    </section>
  )
}

function StepButton({
  label,
  onClick,
  children,
}: {
  label: string
  onClick: () => void
  children: React.ReactNode
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      aria-label={label}
      className="inline-flex size-11 items-center justify-center rounded-full border border-gray-200 bg-white text-gray-600 shadow-[0_1px_2px_rgb(15_23_42_/_0.06)] transition-all duration-200 hover:bg-gray-50 hover:text-gray-900 active:scale-[0.96] dark:border-white/15 dark:bg-white/5 dark:text-gray-300 dark:hover:bg-white/10 dark:hover:text-white"
    >
      {children}
    </button>
  )
}
