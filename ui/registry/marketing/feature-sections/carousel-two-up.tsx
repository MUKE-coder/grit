'use client'

import { useEffect, useRef, useState } from 'react'
import { ArrowLeft, ArrowRight, Check, Play } from 'lucide-react'

/**
 * A two-up feature carousel.
 *
 * Scroll-snap, not a transform track. It gets touch swipe, trackpad swipe,
 * keyboard scrolling and correct behaviour at every viewport width for free,
 * and the arrows are then a convenience on top of a control that already works
 * rather than the only way to move.
 *
 * The card surfaces are CSS gradients, not photographs. A block whose visual
 * depends on a remote image is one broken host away from shipping empty boxes
 * into someone's landing page, and it would force every consumer to add a
 * remote pattern to their Next config before it renders at all.
 */

/** True when the OS has been asked for less motion. */
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

/** The atmospheric card backdrop. Warm haze above, cool depth below. */
function Canvas({ children }: { children: React.ReactNode }) {
  return (
    <div className="relative aspect-[4/3.4] overflow-hidden rounded-3xl bg-[linear-gradient(180deg,#eae5db_0%,#e6e3dd_42%,#cdd6e0_74%,#b2c1cf_100%)] dark:bg-[linear-gradient(180deg,#1c2430_0%,#1a222c_42%,#16202b_74%,#111a24_100%)]">
      {/* Two soft blooms stand in for a horizon. Cheap, and it never 404s. */}
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0"
        style={{
          background:
            'radial-gradient(60% 40% at 30% 88%, rgb(120 145 175 / 0.35), transparent 70%), radial-gradient(50% 34% at 76% 94%, rgb(96 122 154 / 0.3), transparent 70%)',
        }}
      />
      <div className="absolute inset-0 flex items-center justify-center p-8">
        {children}
      </div>
    </div>
  )
}

function NoteMock() {
  return (
    <div className="w-full max-w-[19rem] rounded-2xl bg-white p-5 shadow-[0_18px_50px_-16px_rgb(15_23_42_/_0.35)] dark:bg-gray-900">
      <p className="text-[12.5px] text-gray-500 dark:text-gray-400">Today 09:15 AM</p>
      <p className="mt-1.5 text-[15px] font-semibold text-gray-900 dark:text-white">
        Marketing Website Launch
      </p>
      <p className="mt-3 text-[13.5px]/[1.55] text-gray-500 dark:text-gray-400">
        The new marketing website is scheduled to go live next Monday.{' '}
        <span className="font-semibold text-gray-900 dark:text-white">
          Key highlights include a redesigned hero section, improved SEO structure,
        </span>{' '}
        and an integrated analytics dashboard.
      </p>
      <span className="mt-4 inline-flex items-center gap-2 rounded-full bg-sky-50 px-3 py-1.5 text-[12.5px] font-semibold text-gray-800 ring-1 ring-sky-200 dark:bg-sky-500/10 dark:text-sky-100 dark:ring-sky-500/30">
        <Play aria-hidden="true" className="size-3 fill-current" />
        03:47
      </span>
    </div>
  )
}

const TASKS = [
  { label: 'Review pull requests', done: true },
  { label: 'Update documentation', done: true },
  { label: 'Deploy to staging', done: true },
  { label: 'Write unit tests', done: false },
  { label: 'Send weekly report', done: false },
]

function TasksMock() {
  const done = TASKS.filter((t) => t.done).length
  return (
    <div className="relative w-full max-w-[18rem]">
      {/* A second sheet peeking out below, so it reads as a pad of notes. */}
      <div
        aria-hidden
        className="absolute inset-x-2 -bottom-2 h-full rounded-xl bg-amber-200/70 dark:bg-amber-300/25"
      />
      <div className="relative rounded-xl bg-amber-100 p-5 shadow-[0_18px_50px_-16px_rgb(120_80_10_/_0.35)] dark:bg-amber-200/90">
        <div className="flex items-baseline justify-between">
          <p className="text-[14px] font-semibold text-amber-950">Quick Tasks</p>
          <p className="text-[12.5px] font-medium text-amber-800/70 tabular-nums">
            {done}/{TASKS.length}
          </p>
        </div>
        <ul className="mt-3 space-y-2">
          {TASKS.map((task) => (
            <li key={task.label} className="flex items-center gap-2.5">
              <span
                className={`flex size-4 shrink-0 items-center justify-center rounded ${
                  task.done
                    ? 'bg-emerald-400/90 text-emerald-950'
                    : 'border border-amber-700/35'
                }`}
              >
                {task.done && <Check aria-hidden="true" className="size-3" strokeWidth={3} />}
              </span>
              <span
                className={`text-[13.5px] ${
                  task.done
                    ? 'text-amber-800/55 line-through'
                    : 'font-medium text-amber-950'
                }`}
              >
                {task.label}
              </span>
            </li>
          ))}
        </ul>
      </div>
    </div>
  )
}

const SLIDES = [
  {
    mock: <NoteMock />,
    lead: 'Smart email composition',
    body: 'with AI-powered suggestions, templates, and seamless collaboration for faster communication.',
  },
  {
    mock: <TasksMock />,
    lead: 'Organized note-taking',
    body: 'with rich formatting, tagging, and instant search to capture and retrieve ideas effortlessly.',
  },
  {
    mock: <NoteMock />,
    lead: 'Shared team inbox',
    body: 'with assignment, internal notes, and a full history of every conversation in one place.',
  },
  {
    mock: <TasksMock />,
    lead: 'Lightweight task tracking',
    body: 'with checklists that live beside the work instead of in a separate tool nobody opens.',
  },
]

export default function CarouselTwoUp({
  title = 'Powerful features\nfor modern teams',
}: {
  /** Newlines are honoured, so the headline breaks where you intend it to. */
  title?: string
}) {
  const trackRef = useRef<HTMLDivElement>(null)
  const [atStart, setAtStart] = useState(true)
  const [atEnd, setAtEnd] = useState(false)
  const reduce = useReducedMotion()

  // Read the real scroll position rather than tracking an index. With snap
  // scrolling the user can move the track by swiping, and an index would drift
  // out of sync with what is actually on screen.
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
    // One card plus one gap, so a page always lands on a snap point.
    const step = card ? card.offsetWidth + 24 : el.clientWidth
    el.scrollBy({ left: direction * step, behavior: reduce ? 'auto' : 'smooth' })
  }

  return (
    <section className="bg-gray-50/70 py-20 sm:py-28 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6">
        <div className="flex items-end justify-between gap-8">
          <h2 className="max-w-xl text-[2rem] leading-[1.12] font-semibold tracking-[-0.02em] whitespace-pre-line text-gray-900 sm:text-[2.75rem] dark:text-white">
            {title}
          </h2>

          <div className="hidden shrink-0 items-center gap-2 sm:flex">
            <ArrowButton
              label="Previous features"
              disabled={atStart}
              onClick={() => page(-1)}
            >
              <ArrowLeft aria-hidden="true" className="size-4" />
            </ArrowButton>
            <ArrowButton label="Next features" disabled={atEnd} onClick={() => page(1)}>
              <ArrowRight aria-hidden="true" className="size-4" />
            </ArrowButton>
          </div>
        </div>

        <div
          ref={trackRef}
          onScroll={sync}
          // The scrollbar is hidden deliberately: the arrows and swipe are the
          // affordance, and a chrome scrollbar across a soft card row is the
          // loudest thing on the section.
          className="mt-12 flex snap-x snap-mandatory gap-6 overflow-x-auto pb-2 [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
        >
          {SLIDES.map((slide, i) => (
            <article
              key={i}
              className="w-[85%] shrink-0 snap-start sm:w-[calc((100%-1.5rem)/2)]"
            >
              <Canvas>{slide.mock}</Canvas>
              <p className="mt-5 text-[15px]/[1.6] text-gray-500 dark:text-gray-400">
                <span className="font-semibold text-gray-900 dark:text-white">
                  {slide.lead}
                </span>{' '}
                {slide.body}
              </p>
            </article>
          ))}
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
      // 44px, the smallest target that is comfortable to hit, and visibly
      // dimmed at the ends rather than silently doing nothing.
      className="inline-flex size-11 items-center justify-center rounded-full border border-gray-200 bg-white text-gray-700 shadow-[0_1px_2px_rgb(15_23_42_/_0.06)] transition-all duration-200 hover:bg-gray-50 active:scale-[0.96] disabled:pointer-events-none disabled:opacity-35 dark:border-white/15 dark:bg-white/5 dark:text-gray-200 dark:hover:bg-white/10"
    >
      {children}
    </button>
  )
}
