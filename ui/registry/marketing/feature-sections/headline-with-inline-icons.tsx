'use client'

import { Bot, Brain, ChevronRight, Plus } from 'lucide-react'

/**
 * An oversized statement headline with inline icon chips, over a ruled canvas.
 *
 * The two-tone headline is the whole idea: the named concepts are full-contrast
 * and everything else drops to a mid grey, so the sentence still reads as a
 * sentence while the two terms you are selling carry the weight.
 *
 * Not a client component by necessity — nothing here holds state — but marked
 * anyway so it drops into an App Router project unchanged whichever way the
 * consumer composes it.
 */

/** A named concept in the headline: icon chip plus emphasised label. */
function Concept({
  icon: Icon,
  children,
  tint = '',
}: {
  icon: typeof Brain
  children: React.ReactNode
  tint?: string
}) {
  return (
    <span className="inline-flex items-baseline gap-2.5 align-baseline">
      <span
        aria-hidden
        // inline-flex on a baseline-aligned span: the chip needs to sit on the
        // text baseline, not the line box, or it rides high on the taller display
        // sizes and looks pasted on.
        className="inline-flex size-[0.82em] translate-y-[0.08em] items-center justify-center rounded-[0.22em] border border-gray-200 bg-white text-gray-700 shadow-[0_1px_2px_rgb(15_23_42_/_0.06)] dark:border-white/15 dark:bg-white/10 dark:text-gray-200"
      >
        <Icon className="size-[0.5em]" />
      </span>
      <span className={tint || 'text-gray-900 dark:text-white'}>{children}</span>
    </span>
  )
}

/**
 * A drawing of a compose window. aria-hidden because it is scenery: the
 * addresses and field labels are invented, and read aloud they arrive as if
 * this page were an inbox. Hiding it also settles the contrast question — the
 * faint field labels are decoration, not text someone has to read.
 */
function ComposeMock() {
  return (
    <div
      aria-hidden="true"
      className="w-full max-w-[23rem] rounded-2xl bg-white/95 p-4 shadow-[0_24px_60px_-20px_rgb(15_23_42_/_0.4)] backdrop-blur-sm dark:bg-gray-900/95">
      <div className="flex items-center justify-between gap-3 border-b border-gray-200 pb-3 dark:border-white/10">
        <span className="flex items-center gap-1.5 text-[13px] text-gray-500 dark:text-gray-400">
          To:
          <span className="inline-flex items-center gap-1.5 rounded-full bg-gray-100 py-0.5 pr-2.5 pl-1 dark:bg-white/10">
            <span
              aria-hidden
              className="flex size-5 items-center justify-center rounded-full bg-gradient-to-br from-violet-400 to-sky-400 text-[9px] font-bold text-white"
            >
              SH
            </span>
            <span className="text-[12.5px] font-medium text-gray-900 dark:text-white">
              Shadcn
            </span>
          </span>
        </span>
        <span
          aria-hidden
          className="flex size-7 shrink-0 items-center justify-center rounded-full bg-gray-100 text-gray-600 dark:bg-white/10 dark:text-gray-300"
        >
          <Plus className="size-3.5" />
        </span>
      </div>

      {['Cc:', 'Subject:', 'From:'].map((label) => (
        <p
          key={label}
          className="border-b border-gray-100 py-2.5 text-[13px] text-gray-400 dark:border-white/[0.06] dark:text-gray-500"
        >
          {label}
        </p>
      ))}

      <p className="mt-4 text-[13.5px]/[1.6] text-gray-700 dark:text-gray-300">
        Web applications with{' '}
        <span className="text-sky-600 dark:text-sky-400">React and TypeScript</span> using
        best practices.
      </p>
      <p className="mt-3 text-[13.5px] text-gray-500 dark:text-gray-400">
        Sent from my iPhone
      </p>
    </div>
  )
}

export default function HeadlineWithInlineIcons({
  eyebrow = '[ 0.1 ]  FEATURES',
  subhead = 'Accelerate development\nwith AI-powered assistance',
  body = 'Switch between GPT, Claude, and Gemini from a unified interface with seamless provider switching.',
  ctaLabel = 'Learn more',
  quote = 'Looks really good. Did you design in code or Figma first?',
  quoteName = 'Shadcn',
  quoteRole = 'Creator of Shadcn UI',
}: {
  eyebrow?: string
  subhead?: string
  body?: string
  ctaLabel?: string
  quote?: string
  quoteName?: string
  quoteRole?: string
}) {
  return (
    <section className="relative overflow-hidden bg-gray-50/60 py-20 sm:py-24 dark:bg-gray-950">
      {/* A faint dashed grid, pinned to the content column so the rules line up
          with the copy instead of floating at arbitrary offsets. */}
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0 mx-auto max-w-7xl"
      >
        <div className="absolute inset-y-0 left-[30%] border-l border-dashed border-gray-200 dark:border-white/[0.07]" />
        <div className="absolute inset-y-0 right-[8%] border-l border-dashed border-gray-200 dark:border-white/[0.07]" />
        <div className="absolute inset-x-0 top-[38%] border-t border-dashed border-gray-200 dark:border-white/[0.07]" />
      </div>

      <div className="relative mx-auto max-w-7xl px-6">
        <p className="font-mono text-[11.5px] tracking-wider text-gray-600 uppercase dark:text-gray-400">
          {eyebrow}
        </p>

        <h2 className="mt-8 max-w-5xl text-[2.25rem] leading-[1.18] font-semibold tracking-[-0.025em] text-gray-500 sm:text-[3.25rem] sm:leading-[1.16] dark:text-gray-400">
          Build faster with{' '}
          <Concept icon={Brain} tint="text-pink-500 dark:text-pink-400">
            LLMs
          </Concept>{' '}
          that adapt to your workflow, and let{' '}
          <Concept icon={Bot}>Personal Agents</Concept> handle the rest.
        </h2>

        <div className="mt-20 grid items-start gap-14 lg:grid-cols-2 lg:gap-16">
          <div>
            <h3 className="text-[1.375rem] leading-[1.25] font-semibold tracking-[-0.015em] whitespace-pre-line text-gray-900 dark:text-white">
              {subhead}
            </h3>
            <p className="mt-4 max-w-md text-[15.5px]/[1.6] text-gray-500 dark:text-gray-400">
              {body}
            </p>

            <a
              href="#"
              className="group mt-7 inline-flex h-11 items-center gap-1.5 rounded-xl border border-gray-200 bg-white px-5 text-[14.5px] font-semibold text-gray-900 shadow-[0_1px_2px_rgb(15_23_42_/_0.06)] transition-all duration-200 hover:bg-gray-50 active:scale-[0.98] dark:border-white/15 dark:bg-white/5 dark:text-white dark:hover:bg-white/10"
            >
              {ctaLabel}
              <ChevronRight
                aria-hidden="true"
                className="size-4 transition-transform duration-200 group-hover:translate-x-0.5"
              />
            </a>

            <figure className="mt-20">
              <blockquote className="max-w-xs text-[16.5px]/[1.5] text-gray-900 dark:text-white">
                &ldquo;{quote}&rdquo;
              </blockquote>
              <figcaption className="mt-4 flex items-center gap-3">
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

          <div className="relative aspect-[4/4.2] overflow-hidden rounded-3xl bg-[linear-gradient(180deg,#cfd6dd_0%,#dcd8cf_30%,#e0cfc0_58%,#cbb2a0_80%,#a9c4c8_100%)] dark:bg-[linear-gradient(180deg,#1e2732_0%,#262319_30%,#2b201a_58%,#221812_80%,#16262a_100%)]">
            <div
              aria-hidden
              className="pointer-events-none absolute inset-0"
              style={{
                background:
                  'radial-gradient(56% 30% at 20% 96%, rgb(120 108 88 / 0.4), transparent 70%), radial-gradient(46% 26% at 82% 99%, rgb(96 118 122 / 0.4), transparent 70%)',
              }}
            />
            <div className="absolute inset-0 flex items-center justify-center p-6">
              <ComposeMock />
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
