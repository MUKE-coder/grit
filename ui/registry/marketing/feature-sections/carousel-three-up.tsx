'use client'

import { useEffect, useRef, useState } from 'react'
import { ArrowDown, ArrowLeft, ArrowRight, ArrowUp, CornerDownLeft, Sparkles } from 'lucide-react'

/**
 * A three-up feature carousel on a ruled grid.
 *
 * Same scroll-snap mechanics as the two-up variant — see the note there for why
 * it is not a transform track. The difference is the framing: dashed vertical
 * rules run the full height of the row, so the cells read as panels cut out of
 * a sheet rather than as cards floating on it.
 */

function useReducedMotion() {
  const [reduce, setReduce] = useState(false)
  useEffect(() => {
    const mq = window.matchMedia('(prefers-reduced-motion: reduce)')
    setReduce(mq.matches)
    const onChange = () => setReduce(mq.matches)
    mq.addEventListener('change', onChange)
    return () => mq.removeEventListener('change', onChange)
  }, [])
  return reduce
}

const SUGGESTIONS = [
  { text: '…authentication with OAuth 2.0?', hint: true },
  { text: '…a dark mode toggle in React?' },
  { text: '…caching for API responses?' },
]

function CompletionMock() {
  return (
    <div className="w-full max-w-[19rem]">
      <div className="rounded-xl border border-gray-200 bg-white px-3.5 py-3 shadow-[0_1px_2px_rgb(15_23_42_/_0.05)] dark:border-white/10 dark:bg-gray-900">
        <p className="text-[13.5px] text-gray-900 dark:text-white">
          How do I implement
          {/* A caret drawn as a box, so there is no animation to distract and no
              blinking element for anyone sensitive to motion. */}
          <span
            aria-hidden
            className="ml-0.5 inline-block h-[1.05em] w-px translate-y-[0.15em] bg-gray-900 dark:bg-white"
          />
        </p>
      </div>

      <div className="mt-2 overflow-hidden rounded-xl border border-gray-200 bg-white shadow-[0_12px_32px_-12px_rgb(15_23_42_/_0.18)] dark:border-white/10 dark:bg-gray-900">
        <p className="flex items-center gap-2 bg-indigo-50/70 px-3.5 py-2.5 text-[12.5px] font-semibold text-indigo-700 dark:bg-indigo-500/10 dark:text-indigo-300">
          <Sparkles aria-hidden="true" className="size-3.5" />
          AI Suggestions
        </p>
        <ul className="divide-y divide-gray-100 dark:divide-white/[0.06]">
          {SUGGESTIONS.map((s) => (
            <li
              key={s.text}
              className="flex items-center justify-between gap-3 px-3.5 py-2.5 text-[13px] text-gray-700 dark:text-gray-300"
            >
              {s.text}
              {s.hint && (
                <span className="inline-flex shrink-0 items-center gap-1 rounded-md border border-gray-200 px-1.5 py-0.5 font-mono text-[10.5px] text-gray-500 dark:border-white/15 dark:text-gray-400">
                  <CornerDownLeft aria-hidden="true" className="size-3" />
                  Tab
                </span>
              )}
            </li>
          ))}
        </ul>
      </div>

      <div className="mt-2.5 flex items-center justify-between text-[12px] text-gray-500 dark:text-gray-400">
        <span>3 suggestions</span>
        <span className="flex items-center gap-1.5">
          <Key>
            <ArrowUp aria-hidden="true" className="size-3" />
          </Key>
          <Key>
            <ArrowDown aria-hidden="true" className="size-3" />
          </Key>
          to navigate
        </span>
      </div>
    </div>
  )
}

function Key({ children }: { children: React.ReactNode }) {
  return (
    <span className="inline-flex size-[18px] items-center justify-center rounded border border-gray-200 text-gray-500 dark:border-white/15 dark:text-gray-400">
      {children}
    </span>
  )
}

function UsageMock() {
  return (
    <div className="w-full max-w-[19rem] rounded-2xl border border-gray-200 bg-white p-4 shadow-[0_12px_32px_-12px_rgb(15_23_42_/_0.18)] dark:border-white/10 dark:bg-gray-900">
      <p className="text-[14px] font-semibold text-gray-900 dark:text-white">Usage</p>

      <div className="mt-3 rounded-xl border border-gray-200 p-4 dark:border-white/10">
        <div className="flex items-baseline justify-between gap-3">
          <p className="text-[14px] font-semibold text-gray-900 dark:text-white">
            User prompt tokens
          </p>
          <p className="text-[13px] font-semibold text-gray-900 tabular-nums dark:text-white">
            43%
          </p>
        </div>
        <p className="mt-1.5 text-[12.5px]/[1.5] text-gray-500 dark:text-gray-400">
          Using a premium model costs one prompt credit per use.
        </p>

        <div className="mt-3 h-1.5 overflow-hidden rounded-full bg-gray-100 dark:bg-white/10">
          <div className="h-full w-[43%] rounded-full bg-gradient-to-r from-emerald-400 to-sky-400" />
        </div>
        <div className="mt-2 flex items-baseline justify-between text-[12px]">
          <span className="text-gray-500 tabular-nums dark:text-gray-400">
            550 / 1,500 tokens
          </span>
          <span className="font-semibold text-emerald-600 tabular-nums dark:text-emerald-400">
            950 left
          </span>
        </div>

        {/* Clipped by the panel, which implies the list continues. A fade would
            need a background colour that breaks in dark mode; overflow does not. */}
        <div className="mt-3 grid grid-cols-2 gap-2.5">
          {[
            ['Est. Cost', '$0.25'],
            ['Requests', '48'],
          ].map(([label, value]) => (
            <div
              key={label}
              className="rounded-lg border border-gray-200 px-3 py-2 dark:border-white/10"
            >
              <p className="text-[11.5px] text-gray-400 dark:text-gray-500">{label}</p>
              <p className="text-[13px] font-semibold text-gray-900 tabular-nums dark:text-white">
                {value}
              </p>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

function TranslationMock() {
  return (
    <div className="w-full max-w-[19rem]">
      <p className="text-[13.5px]/[1.6] text-gray-300 dark:text-gray-600">
        Corporis voluptates voluptatem atque excepturi, tempore dolor distinctio libero
        dicta vel, nihil rem consequatur esse aspernatur nostrum.
      </p>
      <p className="mt-3 text-[13.5px] font-medium text-indigo-600 dark:text-indigo-400">
        Auto translated from English
      </p>
      <div className="mt-2.5 rounded-xl border border-indigo-200 bg-white p-4 shadow-[0_12px_32px_-12px_rgb(79_70_229_/_0.22)] dark:border-indigo-500/30 dark:bg-gray-900">
        <p className="text-[12px] text-gray-500 dark:text-gray-400">Spanish</p>
        <p className="mt-1.5 text-[13.5px]/[1.6] text-gray-900 dark:text-white">
          Hola, ¿cómo puedo ayudarte hoy? Estoy aquí para responder cualquier pregunta que
          tengas sobre nuestros servicios y productos.
        </p>
      </div>
    </div>
  )
}

const SLIDES = [
  {
    mock: <CompletionMock />,
    lead: 'Intelligent code completion',
    body: 'that understands your codebase and suggests context-aware snippets.',
  },
  {
    mock: <UsageMock />,
    lead: 'Usage analytics',
    body: 'with detailed token tracking, cost estimation, and budget alerts.',
  },
  {
    mock: <TranslationMock />,
    lead: 'Real-time translation',
    body: 'across 50+ languages with natural-sounding output and dialect support.',
  },
  {
    mock: <CompletionMock />,
    lead: 'Inline refactors',
    body: 'proposed as a diff you accept or reject, never applied behind your back.',
  },
]

export default function CarouselThreeUp({
  title = 'Build modern AI\ndevelopment tools',
}: {
  title?: string
}) {
  const trackRef = useRef<HTMLDivElement>(null)
  const [atStart, setAtStart] = useState(true)
  const [atEnd, setAtEnd] = useState(false)
  const reduce = useReducedMotion()

  function sync() {
    const el = trackRef.current
    if (!el) return
    setAtStart(el.scrollLeft <= 2)
    setAtEnd(el.scrollLeft + el.clientWidth >= el.scrollWidth - 2)
  }

  useEffect(() => {
    sync()
    const onResize = () => sync()
    window.addEventListener('resize', onResize)
    return () => window.removeEventListener('resize', onResize)
  }, [])

  function page(direction: 1 | -1) {
    const el = trackRef.current
    if (!el) return
    const card = el.firstElementChild as HTMLElement | null
    const step = card ? card.offsetWidth : el.clientWidth
    el.scrollBy({ left: direction * step, behavior: reduce ? 'auto' : 'smooth' })
  }

  return (
    <section className="bg-white py-20 sm:py-28 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6">
        <div className="flex items-start justify-between gap-8">
          <h2 className="max-w-xl text-[2rem] leading-[1.1] font-semibold tracking-[-0.025em] whitespace-pre-line text-gray-900 sm:text-[3rem] dark:text-white">
            {title}
          </h2>
          <div className="hidden shrink-0 items-center gap-2 pt-3 sm:flex">
            <ArrowButton label="Previous features" disabled={atStart} onClick={() => page(-1)}>
              <ArrowLeft aria-hidden="true" className="size-4" />
            </ArrowButton>
            <ArrowButton label="Next features" disabled={atEnd} onClick={() => page(1)}>
              <ArrowRight aria-hidden="true" className="size-4" />
            </ArrowButton>
          </div>
        </div>
      </div>

      {/* The ruled band is full-bleed and the track sits inside it, so the dashed
          rules read as part of the page rather than as cell borders. */}
      <div className="mt-14 border-y border-dashed border-gray-200 dark:border-white/10">
        <div className="mx-auto max-w-7xl">
          <div
            ref={trackRef}
            onScroll={sync}
            className="flex snap-x snap-mandatory overflow-x-auto [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
          >
            {SLIDES.map((slide, i) => (
              <article
                key={i}
                className="flex w-[88%] shrink-0 snap-start flex-col border-r border-dashed border-gray-200 px-6 py-10 sm:w-1/2 lg:w-1/3 dark:border-white/10"
              >
                {/* The product mock is a drawing of an app, not an app: fake rows,
                invented dates, no behaviour. aria-hidden so a screen reader does
                not read the chrome out as page content, and so its deliberately
                low-contrast styling stops being a text-contrast failure. Nothing
                inside is focusable, which is what makes hiding it safe. */}
                <div aria-hidden="true" className="flex min-h-[15rem] flex-1 items-start">{slide.mock}</div>
                <p className="mt-10 text-[15px]/[1.6] text-gray-500 dark:text-gray-400">
                  <span className="font-semibold text-gray-900 dark:text-white">
                    {slide.lead}
                  </span>{' '}
                  {slide.body}
                </p>
              </article>
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}

function ArrowButton({
  label,
  disabled,
  onClick,
  children,
}: {
  label: string
  disabled: boolean
  onClick: () => void
  children: React.ReactNode
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled}
      aria-label={label}
      className="inline-flex size-11 items-center justify-center rounded-full border border-gray-200 bg-white text-gray-700 shadow-[0_1px_2px_rgb(15_23_42_/_0.06)] transition-all duration-200 hover:bg-gray-50 active:scale-[0.96] disabled:pointer-events-none disabled:opacity-35 dark:border-white/15 dark:bg-white/5 dark:text-gray-200 dark:hover:bg-white/10"
    >
      {children}
    </button>
  )
}
