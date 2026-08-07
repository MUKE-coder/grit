'use client'

import { useCallback, useEffect, useId, useRef, useState } from 'react'
import { ArrowRight, ChevronLeft, ChevronRight, Pause, Play, ShoppingBag, Star } from 'lucide-react'

import { Button } from '@/components/ui/button'

/*
 * A full-bleed storefront hero carousel that advances on its own.
 *
 * Auto-advancing content is the single most reliably broken pattern on a
 * storefront, so most of this file is the parts that usually are not there.
 *
 * A visible pause button, not pause-on-hover. WCAG 2.2.2 wants a way to stop
 * anything that moves for more than five seconds, and hover is not one: it does
 * not exist on touch, and a keyboard user reading a slide has no way to hold it
 * still. The button is the mechanism; hover and focus pausing are conveniences
 * layered on top of it.
 *
 * It starts paused when the visitor has asked for reduced motion. A carousel
 * that swaps its entire viewport every eight seconds is exactly the kind of
 * motion that setting is about, and honouring it only in the CSS transition
 * while still swapping the content misses the point.
 *
 * Off-screen slides are `visibility: hidden`, not merely `opacity: 0`. A
 * transparent slide still has focusable links in it, so tabbing through an
 * opacity-only carousel walks you through three slides' worth of buttons you
 * cannot see. `visibility` removes them from the tab order and the
 * accessibility tree, and unlike `display: none` it still transitions — the
 * property steps at the end of a fade-out and at the start of a fade-in, which
 * is precisely the timing a cross-fade wants.
 *
 * The live region is `off` while playing and `polite` once paused. That is the
 * APG carousel rule and it is counter-intuitive: announcing every automatic
 * change would talk over whatever the person is actually doing, but once they
 * take control with the arrows or the dots they need to hear what they landed
 * on.
 *
 * Prices are integer cents. Percentages are computed from them rather than
 * typed in beside them, because a hardcoded "SAVE 25%" next to a price pair is
 * a claim that goes stale the first time someone edits one number.
 */

export interface Slide {
  id: string
  eyebrow: string
  title: string
  description: string
  image: string
  cta: { label: string; href: string }
  secondaryCta?: { label: string; href: string }
  badge?: string
  /** Integer cents. */
  price?: number
  /** Integer cents, struck through. */
  wasPrice?: number
  align?: 'left' | 'right' | 'center'
  /** Gradient stops for the scrim over the photo. Keep the text side opaque. */
  scrim?: string
}

const SLIDES: Slide[] = [
  {
    id: 'shoes',
    eyebrow: 'Discover your own',
    title: 'The newest run of everyday sneakers',
    description:
      'Premium comfort and contemporary shape, handmade from materials chosen to still look right in three years.',
    image: 'https://images.unsplash.com/photo-1600185365926-3a2ce3cdb9eb?w=1600&h=900&fit=crop&q=80',
    cta: { label: 'Shop collection', href: '#' },
    secondaryCta: { label: 'View lookbook', href: '#' },
    badge: 'New arrival',
    price: 29900,
    wasPrice: 39900,
    align: 'left',
    scrim: 'from-white/90 via-white/60 to-transparent',
  },
  {
    id: 'watches',
    eyebrow: 'Timeless elegance',
    title: 'Watches for every occasion',
    description:
      'Precision movement and a case finished by hand. Each one carries a decade of the same workshop getting it slightly better.',
    image: 'https://images.unsplash.com/photo-1523170335258-f5ed11844a49?w=1600&h=900&fit=crop&q=80',
    cta: { label: 'View collection', href: '#' },
    secondaryCta: { label: 'How they are made', href: '#' },
    badge: 'Premium',
    price: 129900,
    wasPrice: 149900,
    align: 'right',
    scrim: 'from-transparent via-gray-50/60 to-gray-100/90',
  },
  {
    id: 'bags',
    eyebrow: 'Carry your style',
    title: 'Bags and small leather goods',
    description:
      'Full-grain leather, edge-painted and stitched to take a decade of being thrown into the back of a car.',
    image: 'https://images.unsplash.com/photo-1590874103328-eac38a683ce7?w=1600&h=900&fit=crop&q=80',
    cta: { label: 'Explore collection', href: '#' },
    secondaryCta: { label: 'About the leather', href: '#' },
    badge: 'Exclusive',
    price: 89900,
    wasPrice: 119900,
    align: 'left',
    scrim: 'from-amber-50/90 via-amber-50/50 to-transparent',
  },
]

const money = new Intl.NumberFormat('en-US', {
  style: 'currency',
  currency: 'USD',
  maximumFractionDigits: 0,
})

export default function HeroCarouselWithControls({
  slides = SLIDES,
  label = 'Featured collections',
  /** Milliseconds per slide. Below 5000 there is no reading time. */
  interval = 8000,
}: {
  slides?: Slide[]
  label?: string
  interval?: number
}) {
  const [index, setIndex] = useState(0)
  const [playing, setPlaying] = useState(true)
  const [suspended, setSuspended] = useState(false)
  const [progress, setProgress] = useState(0)
  const slidesId = useId()
  const started = useRef(0)

  const running = playing && !suspended

  /* Reduced motion means start stopped, not "start moving and animate less".
     Checked in an effect because the server has no media queries and rendering
     a different initial state there is a hydration mismatch. */
  useEffect(() => {
    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) setPlaying(false)
  }, [])

  const go = useCallback(
    (next: number) => setIndex(((next % slides.length) + slides.length) % slides.length),
    [slides.length],
  )

  /* One rAF loop drives both the advance and the progress bar, so the bar can
     never disagree with the timer it is supposed to represent — which is what
     happens when a setInterval for the bar and a setTimeout for the slide drift
     apart, or when a background tab throttles one and not the other. */
  useEffect(() => {
    if (!running) return
    started.current = performance.now()
    let frame = 0

    const tick = (now: number) => {
      const elapsed = now - started.current
      if (elapsed >= interval) {
        started.current = now
        setProgress(0)
        setIndex((current) => (current + 1) % slides.length)
      } else {
        setProgress((elapsed / interval) * 100)
      }
      frame = requestAnimationFrame(tick)
    }

    frame = requestAnimationFrame(tick)
    return () => cancelAnimationFrame(frame)
  }, [running, interval, slides.length, index])

  useEffect(() => setProgress(0), [index])

  function handleKeyDown(event: React.KeyboardEvent) {
    if (event.key === 'ArrowLeft') {
      event.preventDefault()
      go(index - 1)
    } else if (event.key === 'ArrowRight') {
      event.preventDefault()
      go(index + 1)
    }
  }

  return (
    <section
      aria-roledescription="carousel"
      aria-label={label}
      onKeyDown={handleKeyDown}
      /* Hover and focus suspend it without touching the pause button's own
         state, so leaving does not override a deliberate pause. */
      onMouseEnter={() => setSuspended(true)}
      onMouseLeave={() => setSuspended(false)}
      onFocus={() => setSuspended(true)}
      onBlur={(event) => {
        if (!event.currentTarget.contains(event.relatedTarget)) setSuspended(false)
      }}
      className="relative h-[42rem] w-full overflow-hidden bg-gray-100 md:h-[44rem] dark:bg-gray-900"
    >
      <div
        id={slidesId}
        aria-live={playing ? 'off' : 'polite'}
        aria-atomic="false"
        className="absolute inset-0"
      >
        {slides.map((slide, position) => {
          const active = position === index
          return (
            <div
              key={slide.id}
              role="group"
              aria-roledescription="slide"
              aria-label={`${position + 1} of ${slides.length}: ${slide.title}`}
              /* invisible, not just transparent — see the note above. */
              className={`absolute inset-0 transition-[opacity,visibility] duration-700 ease-in-out ${
                active ? 'visible opacity-100' : 'invisible opacity-0'
              } motion-reduce:transition-none`}
            >
              {/* Empty alt: the slide's heading is the next thing read, and
                  the photo says the same thing less precisely. */}
              <img src={slide.image} alt="" className="absolute inset-0 size-full object-cover" />
              <div
                aria-hidden="true"
                className={`absolute inset-0 bg-gradient-to-r ${
                  slide.scrim ?? 'from-white/90 via-white/60 to-transparent'
                }`}
              />

              {/* The padding here is the controls' gutter, not decoration. At
                  px-4 the previous arrow lands on top of the copy — the first
                  letter of the description disappears behind it — and without
                  the bottom inset the dots pill covers the primary button on a
                  phone, which is the one thing on the slide that had to work. */}
              <div className="relative mx-auto flex h-full max-w-7xl items-center px-16 pt-6 pb-24 md:px-24 md:pb-10">
                <div
                  className={`w-full max-w-xl ${
                    slide.align === 'right'
                      ? 'ml-auto'
                      : slide.align === 'center'
                        ? 'mx-auto text-center'
                        : 'mr-auto'
                  }`}
                >
                  <div className="rounded-2xl border border-white/60 bg-white/40 p-5 shadow-xl backdrop-blur-sm md:p-10">
                    {slide.badge && (
                      <p className="mb-4 inline-flex items-center gap-1.5 rounded-full bg-gray-900 px-3 py-1">
                        <Star aria-hidden="true" className="size-3 fill-amber-400 text-amber-400" />
                        <span className="text-xs font-bold tracking-wider text-white uppercase">
                          {slide.badge}
                        </span>
                      </p>
                    )}

                    <p className="text-sm font-medium tracking-wider text-blue-700 uppercase">
                      {slide.eyebrow}
                    </p>
                    <h2 className="mt-2 text-2xl leading-tight font-bold text-balance text-gray-900 sm:text-3xl md:text-5xl">
                      {slide.title}
                    </h2>
                    <p className="mt-3 text-sm text-pretty text-gray-700 md:mt-4 md:text-lg">
                      {slide.description}
                    </p>

                    {slide.price !== undefined && (
                      <p className="mt-4 inline-flex flex-wrap items-center gap-3 rounded-lg bg-black/10 px-4 py-2 md:mt-6">
                        <span className="text-2xl font-bold text-gray-900">
                          {money.format(slide.price / 100)}
                        </span>
                        {slide.wasPrice !== undefined && (
                          <>
                            {/* The unit matters here: "was $399" reads
                                correctly, a bare struck-through number does
                                not survive being read aloud. */}
                            <span className="text-base text-gray-600">
                              <span className="sr-only">was </span>
                              <s>{money.format(slide.wasPrice / 100)}</s>
                            </span>
                            <span className="rounded bg-red-600 px-2 py-1 text-xs font-bold text-white">
                              Save {Math.round((1 - slide.price / slide.wasPrice) * 100)}%
                            </span>
                          </>
                        )}
                      </p>
                    )}

                    <div className="mt-4 flex flex-wrap gap-3 md:mt-6">
                      <Button
                        asChild
                        className="min-h-12 rounded-full bg-blue-700 px-8 hover:bg-blue-800"
                      >
                        <a href={slide.cta.href}>
                          <ShoppingBag aria-hidden="true" className="size-4" />
                          {slide.cta.label}
                        </a>
                      </Button>
                      {slide.secondaryCta && (
                        <Button
                          asChild
                          variant="outline"
                          className="min-h-12 rounded-full border-gray-400 bg-white/60 px-6 text-gray-800 hover:bg-white"
                        >
                          <a href={slide.secondaryCta.href}>
                            {slide.secondaryCta.label}
                            <ArrowRight aria-hidden="true" className="size-4" />
                          </a>
                        </Button>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          )
        })}
      </div>

      <button
        type="button"
        onClick={() => go(index - 1)}
        aria-controls={slidesId}
        className="absolute top-1/2 left-4 z-20 inline-flex size-12 -translate-y-1/2 items-center justify-center rounded-full border border-white/60 bg-white/40 text-gray-800 shadow-lg backdrop-blur-sm hover:bg-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 md:left-8"
      >
        <ChevronLeft aria-hidden="true" className="size-6" />
        <span className="sr-only">Previous slide</span>
      </button>

      <button
        type="button"
        onClick={() => go(index + 1)}
        aria-controls={slidesId}
        className="absolute top-1/2 right-4 z-20 inline-flex size-12 -translate-y-1/2 items-center justify-center rounded-full border border-white/60 bg-white/40 text-gray-800 shadow-lg backdrop-blur-sm hover:bg-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 md:right-8"
      >
        <ChevronRight aria-hidden="true" className="size-6" />
        <span className="sr-only">Next slide</span>
      </button>

      <div className="absolute inset-x-0 bottom-6 z-20 flex justify-center px-4">
        <div className="flex items-center gap-3 rounded-full border border-white/20 bg-black/40 px-3 py-2 shadow-lg backdrop-blur-md">
          {/* The mechanism WCAG 2.2.2 asks for. Everything else here is a
              convenience; this one is the requirement. */}
          <button
            type="button"
            onClick={() => setPlaying((current) => !current)}
            aria-controls={slidesId}
            className="inline-flex size-11 items-center justify-center rounded-full text-white hover:bg-white/20 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-white"
          >
            {playing ? (
              <Pause aria-hidden="true" className="size-4" />
            ) : (
              <Play aria-hidden="true" className="size-4" />
            )}
            <span className="sr-only">
              {playing ? 'Pause the carousel' : 'Play the carousel'}
            </span>
          </button>

          <p className="text-sm text-white tabular-nums">
            <span className="sr-only">Slide </span>
            <span className="font-bold">{index + 1}</span>
            <span aria-hidden="true"> / </span>
            <span className="sr-only">of </span>
            {slides.length}
          </p>

          <ul role="list" className="flex items-center gap-2">
            {slides.map((slide, position) => (
              <li key={slide.id}>
                <button
                  type="button"
                  onClick={() => go(position)}
                  aria-controls={slidesId}
                  /* aria-current, not aria-selected: these are buttons that
                     move a carousel, not tabs in a tablist. Omitted rather
                     than set to "false" on the others — an explicit false is
                     an extra thing to read for no information. */
                  aria-current={position === index ? 'true' : undefined}
                  className="inline-flex h-11 items-center px-0.5 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
                >
                  <span
                    aria-hidden="true"
                    className={`block h-2 rounded-full transition-all ${
                      position === index ? 'w-8 bg-blue-500' : 'w-2 bg-white/50'
                    }`}
                  />
                  <span className="sr-only">{slide.title}</span>
                </button>
              </li>
            ))}
          </ul>

          {/* Decorative: the counter beside it already says where you are, and
              a progress bar read as a percentage every frame would be noise. */}
          <div aria-hidden="true" className="hidden h-1 w-24 rounded-full bg-white/20 sm:block">
            <div
              className="h-full rounded-full bg-blue-500"
              style={{ width: `${running ? progress : 0}%` }}
            />
          </div>
        </div>
      </div>
    </section>
  )
}
