'use client'

import { useEffect, useRef, useState } from 'react'
import {
  CheckSquare,
  ChevronDown,
  GraduationCap,
  HeartPulse,
  MessagesSquare,
  Star,
} from 'lucide-react'

/** The Grit UI mark, inlined so the block stays self-contained. */
function GritMark({ className = 'size-7' }: { className?: string }) {
  return (
    <svg viewBox="0 0 32 32" fill="none" aria-hidden="true" className={className}>
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M8 0H24A8 8 0 0 1 32 8V24A8 8 0 0 1 24 32H8A8 8 0 0 1 0 24V8A8 8 0 0 1 8 0ZM16 9.4A6.6 6.6 0 1 0 21.4 19.9V17.4H17.2A1.7 1.7 0 0 1 17.2 14H23.1A1.7 1.7 0 0 1 24.8 15.7V20.6A1.7 1.7 0 0 1 24.4 21.7A10 10 0 1 1 22.6 8.2A1.7 1.7 0 0 1 20.4 10.8A6.6 6.6 0 0 0 16 9.4Z"
        fill="currentColor"
      />
    </svg>
  )
}

const MENU_COLUMNS = [
  {
    icon: MessagesSquare,
    title: 'Communication',
    links: ['Chat', 'Social walls', 'Events'],
  },
  {
    icon: GraduationCap,
    title: 'Information & Training',
    links: ['Onboarding', 'Courses', 'Handbook', 'News'],
  },
  {
    icon: HeartPulse,
    title: 'Well-being & Retention',
    links: ['Satisfaction surveys', 'Development surveys'],
  },
  {
    icon: CheckSquare,
    title: 'Daily Operations',
    links: ['To-do lists', 'Forms'],
  },
]

export default function SplitWithStats({
  title = 'Put people first',
  underlinedWord = 'people',
  subtitle = 'Fast, user-friendly and engaging — turn HR into people and culture, and streamline your daily operations with your own branded app.',
  ctaLabel = 'Book a demo',
  placeholder = 'Enter work email',
  stats = [
    { value: '75.2%', label: 'Average daily activity' },
    { value: '~20k', label: 'Average daily users' },
  ],
  rating = 4.5,
  ratingLabel = 'Average user rating',
}: {
  title?: string
  underlinedWord?: string
  subtitle?: string
  ctaLabel?: string
  placeholder?: string
  stats?: { value: string; label: string }[]
  rating?: number
  ratingLabel?: string
}) {
  const [menuOpen, setMenuOpen] = useState(false)
  const [email, setEmail] = useState('')
  const navRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    function onPointerDown(e: PointerEvent) {
      if (navRef.current && !navRef.current.contains(e.target as Node)) setMenuOpen(false)
    }
    function onKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') setMenuOpen(false)
    }
    document.addEventListener('pointerdown', onPointerDown)
    document.addEventListener('keydown', onKeyDown)
    return () => {
      document.removeEventListener('pointerdown', onPointerDown)
      document.removeEventListener('keydown', onKeyDown)
    }
  }, [])

  // Split the headline so one word can carry the hand-drawn underline. Falls
  // back to the plain title when the word is not present, rather than dropping
  // half the headline on the floor.
  const parts = title.split(underlinedWord)
  const hasAccent = parts.length === 2

  return (
    <div className="min-h-screen bg-white dark:bg-gray-950">
      {/* Nav */}
      <div ref={navRef} className="relative z-50">
        <nav
          aria-label="Global"
          className="mx-auto flex max-w-7xl items-center gap-8 px-6 py-5"
        >
          <a href="#" className="flex items-center gap-2 text-gray-900 dark:text-white">
            <GritMark className="size-7" />
            <span className="text-lg font-semibold tracking-tight">Grit</span>
          </a>

          <div className="hidden items-center gap-6 lg:flex">
            <button
              type="button"
              onClick={() => setMenuOpen((o) => !o)}
              aria-expanded={menuOpen}
              className="inline-flex items-center gap-1 text-sm font-medium text-gray-900 dark:text-white"
            >
              Product
              <ChevronDown
                aria-hidden="true"
                className={`size-4 text-gray-400 transition-transform duration-200 ${
                  menuOpen ? 'rotate-180' : ''
                }`}
              />
            </button>
            {['Why us', 'About us', 'Cases', 'Blog'].map((item) => (
              <a
                key={item}
                href="#"
                className="text-sm text-gray-600 transition-colors hover:text-gray-900 dark:text-gray-400 dark:hover:text-white"
              >
                {item}
              </a>
            ))}
          </div>

          <div className="ml-auto flex items-center gap-4">
            <a
              href="#"
              className="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-900 transition-colors hover:bg-gray-50 dark:border-white/15 dark:text-white dark:hover:bg-white/5"
            >
              {ctaLabel}
            </a>
            <span className="hidden items-center gap-1 text-sm text-gray-600 sm:inline-flex dark:text-gray-400">
              English
              <ChevronDown aria-hidden="true" className="size-4 text-gray-400" />
            </span>
          </div>
        </nav>

        {/* Four-column mega menu, flush to the nav as a full-width sheet */}
        {menuOpen && (
          <div className="absolute inset-x-0 top-full border-b border-gray-200 bg-white shadow-lg dark:border-white/10 dark:bg-gray-900">
            <div className="mx-auto max-w-7xl px-6 py-8">
              <div className="mb-6 flex items-baseline gap-4">
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                  Product
                </h2>
                <a
                  href="#"
                  className="text-sm text-gray-500 transition-colors hover:text-gray-900 dark:text-gray-400 dark:hover:text-white"
                >
                  All features <span aria-hidden="true">&rarr;</span>
                </a>
              </div>

              <div className="grid grid-cols-2 gap-8 md:grid-cols-4 md:divide-x md:divide-gray-200 dark:md:divide-white/10">
                {MENU_COLUMNS.map(({ icon: Icon, title: col, links }, i) => (
                  <div key={col} className={i > 0 ? 'md:pl-8' : undefined}>
                    <div className="mb-4 flex items-center gap-2">
                      <Icon
                        aria-hidden="true"
                        className="size-4 text-gray-900 dark:text-white"
                      />
                      <span className="text-sm font-semibold text-gray-900 dark:text-white">
                        {col}
                      </span>
                    </div>
                    <ul className="space-y-2.5">
                      {links.map((link) => (
                        <li key={link}>
                          <a
                            href="#"
                            className="text-sm text-gray-500 transition-colors hover:text-gray-900 dark:text-gray-400 dark:hover:text-white"
                          >
                            {link}
                          </a>
                        </li>
                      ))}
                    </ul>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Hero */}
      <div className="mx-auto grid max-w-7xl grid-cols-1 items-center gap-14 px-6 py-16 lg:grid-cols-2 lg:py-24">
        <div>
          <h1 className="text-5xl font-bold tracking-tight text-gray-900 sm:text-6xl dark:text-white">
            {hasAccent ? (
              <>
                {parts[0]}
                <span className="relative inline-block">
                  {underlinedWord}
                  {/* Hand-drawn underline: a double stroke reads as pen, a
                      single straight rule reads as a text-decoration. */}
                  <svg
                    aria-hidden="true"
                    viewBox="0 0 200 16"
                    preserveAspectRatio="none"
                    className="absolute -bottom-1 left-0 h-2.5 w-full text-gray-900 dark:text-white"
                  >
                    <path
                      d="M2 9C40 4 120 3 198 7"
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="2.5"
                      strokeLinecap="round"
                    />
                    <path
                      d="M14 13C60 9 130 9 186 12"
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="1.5"
                      strokeLinecap="round"
                      opacity="0.5"
                    />
                  </svg>
                </span>
                {parts[1]}
              </>
            ) : (
              title
            )}
          </h1>

          <p className="mt-7 max-w-md text-lg/8 text-gray-600 dark:text-gray-400">
            {subtitle}
          </p>

          {/* Inline capture: the button lives inside the field border, which is
              what makes this read as one control rather than two. */}
          <form
            onSubmit={(e) => e.preventDefault()}
            className="mt-8 flex max-w-md items-center gap-2 rounded-xl border border-gray-300 p-1.5 focus-within:border-gray-400 dark:border-white/15"
          >
            <label htmlFor="grit-work-email" className="sr-only">
              Work email
            </label>
            <input
              id="grit-work-email"
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder={placeholder}
              className="min-w-0 flex-1 bg-transparent px-3 py-2 text-sm text-gray-900 placeholder:text-gray-400 focus:outline-none dark:text-white dark:placeholder:text-gray-500"
            />
            <button
              type="submit"
              className="shrink-0 rounded-lg bg-green-400 px-5 py-2.5 text-sm font-semibold text-gray-900 transition-colors hover:bg-green-300"
            >
              {ctaLabel}
            </button>
          </form>

          {/* Stats */}
          <div className="mt-12 flex max-w-md divide-x divide-gray-200 border-b border-gray-200 pb-8 dark:divide-white/10 dark:border-white/10">
            {stats.map((s, i) => (
              <div key={s.label} className={i > 0 ? 'pl-8' : 'pr-8'}>
                <p className="text-3xl font-bold text-gray-900 dark:text-white">
                  {s.value}
                </p>
                <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">{s.label}</p>
              </div>
            ))}
          </div>

          {/* Rating */}
          <div className="mt-6 flex items-center gap-3">
            <div className="flex items-center gap-0.5">
              {[0, 1, 2, 3, 4].map((i) => (
                <Star
                  key={i}
                  aria-hidden="true"
                  className={`size-5 ${
                    i < Math.floor(rating)
                      ? 'fill-gray-900 text-gray-900 dark:fill-white dark:text-white'
                      : i < rating
                        ? 'fill-gray-900/50 text-gray-900 dark:fill-white/50 dark:text-white'
                        : 'text-gray-300 dark:text-gray-700'
                  }`}
                />
              ))}
            </div>
            <span className="text-sm font-semibold text-gray-900 dark:text-white">
              {rating}
            </span>
            <span className="text-sm text-gray-500 dark:text-gray-400">{ratingLabel}</span>
          </div>
        </div>

        {/* Isometric card stack, built with CSS transforms rather than a
            flat image so it stays crisp at any size and themes correctly. */}
        <div className="relative hidden h-[30rem] lg:block" aria-hidden="true">
          <div
            className="absolute inset-0"
            style={{ perspective: '1400px', perspectiveOrigin: '60% 40%' }}
          >
            <div
              className="absolute inset-0"
              style={{ transform: 'rotateX(52deg) rotateZ(-42deg)', transformStyle: 'preserve-3d' }}
            >
              {/* Phone slab */}
              <div className="absolute top-16 left-40 h-72 w-44 rounded-[2rem] border-2 border-gray-300 dark:border-white/20">
                <div className="mx-auto mt-3 h-4 w-16 rounded-full bg-gray-200 dark:bg-white/10" />
              </div>

              {/* Chat card */}
              <div className="absolute top-8 left-16 w-60 rounded-xl border-2 border-gray-900 bg-white p-3.5 dark:border-white dark:bg-gray-900">
                <p className="text-sm font-semibold text-gray-900 dark:text-white">
                  Team Chat
                </p>
                <div className="mt-3 flex justify-end">
                  <div className="h-6 w-24 rounded-lg bg-green-300" />
                </div>
                <div className="mt-2 h-5 w-20 rounded-lg bg-gray-200 dark:bg-white/10" />
                <div className="mt-3 h-7 rounded-lg border border-gray-200 dark:border-white/10" />
              </div>

              {/* Schedule card */}
              <div className="absolute top-24 -left-4 w-48 rounded-xl border border-gray-300 bg-white p-3 dark:border-white/20 dark:bg-gray-900">
                <div className="mb-2 inline-block rounded-md bg-gray-900 px-2 py-0.5 text-[10px] font-semibold text-white dark:bg-white dark:text-gray-900">
                  TUE 12
                </div>
                {['9:00', '10:00', '11:00'].map((t) => (
                  <div key={t} className="mb-1.5 flex items-center gap-2">
                    <span className="w-8 text-[9px] text-gray-400">{t}</span>
                    <div className="h-3 flex-1 rounded bg-gray-100 dark:bg-white/5" />
                  </div>
                ))}
              </div>

              {/* To-do card */}
              <div className="absolute top-64 left-8 w-52 rounded-xl border-2 border-gray-900 bg-white p-3.5 dark:border-white dark:bg-gray-900">
                <p className="text-sm font-semibold text-gray-900 dark:text-white">
                  To-do list
                </p>
                <div className="mt-3 flex items-center gap-2">
                  <span className="flex size-5 items-center justify-center rounded bg-green-400 text-[10px] font-bold text-gray-900">
                    ✓
                  </span>
                  <div className="h-2.5 w-16 rounded bg-gray-900 dark:bg-white" />
                </div>
                {[0, 1, 2].map((i) => (
                  <div key={i} className="mt-2 flex items-center gap-2">
                    <span className="size-5 rounded border border-gray-300 dark:border-white/20" />
                    <div className="h-2 w-24 rounded bg-gray-200 dark:bg-white/10" />
                  </div>
                ))}
              </div>

              {/* Performance card */}
              <div className="absolute top-72 left-56 w-52 rounded-xl border border-gray-300 bg-white p-3.5 dark:border-white/20 dark:bg-gray-900">
                <div className="mb-2 flex items-end gap-1">
                  {[10, 16, 8, 20, 12, 18].map((h, i) => (
                    <span
                      key={i}
                      className="w-1.5 rounded-sm bg-gray-900 dark:bg-white"
                      style={{ height: h }}
                    />
                  ))}
                </div>
                <p className="text-[10px] text-gray-500 dark:text-gray-400">
                  Overall task performance
                </p>
                <p className="text-lg font-bold text-gray-900 dark:text-white">85.3%</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
