import { ArrowRight, Play } from 'lucide-react'

/*
 * A storefront hero laid out as tiles: the headline, a promise, the proof, and
 * three ways in.
 *
 * The thing to copy here is what is NOT in it. The source wrapped a <button>
 * inside a <Link>, which is invalid HTML — interactive content is not allowed
 * inside an anchor — and browsers resolve it inconsistently: some fire the
 * link, some the button, and a screen reader announces a link containing a
 * button with no way to tell which it is about to activate. Every tile here is
 * either a link or a button, never one inside the other, and the play glyph is
 * decoration inside the link rather than a control of its own.
 *
 * The headline is set with `text-balance` and a max width rather than the
 * source's hard `<br />` after every second word. A forced break is right at
 * exactly one viewport and wrong at all the others, and on a narrow screen it
 * produces a column of orphans.
 *
 * The avatars are `alt=""` and the count beside them carries the meaning. The
 * source labelled them "User 1", "User 2", "User 3", which is three
 * announcements that tell a screen reader user nothing at all — and two of the
 * three were the same photograph.
 *
 * Images are plain <img> with verified remote URLs. The source pointed at the
 * author's own UploadThing CDN, which breaks for everyone else and bills them
 * for the bandwidth until it does.
 */

export interface Tile {
  title: string
  body?: string
  href?: string
}

const HERO_IMAGE =
  'https://images.unsplash.com/photo-1570172619644-dfd03ed5d881?w=1000&h=800&fit=crop&q=80'

const PRODUCT_IMAGES = [
  'https://images.unsplash.com/photo-1620916566398-39f1143ab7be?w=200&h=200&fit=crop&q=80',
  'https://images.unsplash.com/photo-1608248543803-ba4f8c70ae0b?w=200&h=200&fit=crop&q=80',
]

const AVATARS = [
  'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=120&h=120&fit=crop&q=80',
  'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=120&h=120&fit=crop&q=80',
  'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=120&h=120&fit=crop&q=80',
]

export default function BentoHeroWithProof({
  title = 'Skincare that actually does something',
  ctaLabel = 'Shop the range',
  ctaHref = '#',
  reviewCount = '5.2K',
  reviewLabel = 'five-star reviews from our customers',
  claim = 'Formulated with dermatologists, tested on nobody.',
  bestSelling = { title: 'Best selling products', body: 'View the range', href: '#' },
  featured = { title: 'Take care of your skin', body: 'See the latest arrivals', href: '#' },
  tutorialsLabel = 'Watch the tutorials',
  tutorialsHref = '#',
  learnMoreHref = '#',
}: {
  title?: string
  ctaLabel?: string
  ctaHref?: string
  reviewCount?: string
  reviewLabel?: string
  claim?: string
  bestSelling?: Tile
  featured?: Tile
  tutorialsLabel?: string
  tutorialsHref?: string
  learnMoreHref?: string
}) {
  return (
    <section className="bg-white p-6 md:p-8 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl rounded-3xl bg-gradient-to-r from-pink-50 to-purple-50 p-6 md:p-8 lg:p-10 dark:from-pink-500/10 dark:to-purple-500/10">
        <div className="grid grid-cols-1 gap-8 lg:grid-cols-2">
          <div className="flex flex-col justify-center">
            {/* Balanced rather than hard-broken: a forced <br /> is correct at
                one width and wrong at every other. */}
            <h2 className="max-w-md text-4xl leading-tight font-bold text-balance text-gray-900 md:text-5xl lg:text-6xl dark:text-white">
              {title}
            </h2>

            <a
              href={ctaHref}
              className="mt-8 inline-flex min-h-11 w-fit items-center gap-2 rounded-full bg-gray-900 px-8 text-sm font-medium text-white hover:bg-gray-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-900 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100"
            >
              {ctaLabel}
              <ArrowRight aria-hidden="true" className="size-4" />
            </a>

            <div className="mt-10 grid grid-cols-1 gap-4 md:grid-cols-2">
              <a
                href={bestSelling.href ?? '#'}
                className="group relative flex flex-col overflow-hidden rounded-3xl bg-pink-400 p-6 text-white transition-shadow hover:shadow-lg focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-pink-600"
              >
                <h3 className="text-2xl font-bold">{bestSelling.title}</h3>
                {bestSelling.body && <p className="mt-1">{bestSelling.body}</p>}
                <span
                  aria-hidden="true"
                  className="mt-10 inline-flex size-10 items-center justify-center self-end rounded-full bg-white/20 transition group-hover:bg-white/30"
                >
                  <ArrowRight className="size-5" />
                </span>
              </a>

              <div className="flex flex-col justify-between rounded-3xl bg-gray-200 p-6 dark:bg-white/10">
                <ul role="list" className="flex justify-center gap-2">
                  {PRODUCT_IMAGES.map((src) => (
                    <li
                      key={src}
                      className="size-16 overflow-hidden rounded-full bg-white dark:bg-gray-900"
                    >
                      <img src={src} alt="" className="size-full object-cover" />
                    </li>
                  ))}
                </ul>
                <a
                  href={learnMoreHref}
                  className="mt-6 inline-flex min-h-11 items-center justify-center gap-2 rounded-full bg-white px-4 text-sm font-medium text-gray-900 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-900"
                >
                  Learn more
                  <ArrowRight aria-hidden="true" className="size-4" />
                </a>
              </div>
            </div>
          </div>

          <div className="space-y-6">
            <a
              href={featured.href ?? '#'}
              className="group relative block overflow-hidden rounded-3xl bg-gray-100 p-4 transition-shadow hover:shadow-lg focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-900 dark:bg-white/5"
            >
              <img
                src={HERO_IMAGE}
                alt=""
                className="h-64 w-full rounded-2xl object-cover md:h-80"
              />
              {/* A scrim, not a guess: the caption sits over a photograph
                  whose brightness you do not control. */}
              <span
                aria-hidden="true"
                className="absolute inset-4 rounded-2xl bg-gradient-to-t from-gray-950/70 via-gray-950/10 to-transparent"
              />
              <span className="absolute bottom-8 left-8 text-white">
                <span className="block text-2xl font-bold">{featured.title}</span>
                {featured.body && (
                  <span className="mt-1 hidden text-sm sm:block">{featured.body}</span>
                )}
              </span>
              <span
                aria-hidden="true"
                className="absolute top-1/2 left-8 inline-flex size-12 -translate-y-1/2 items-center justify-center rounded-full bg-white/30 backdrop-blur-sm transition group-hover:bg-white/40"
              >
                <ArrowRight className="size-5 text-white" />
              </span>
            </a>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div className="flex flex-col justify-center">
                {/* Decorative: the count below is what carries the meaning. */}
                <ul role="list" aria-hidden="true" className="flex -space-x-4">
                  {AVATARS.map((src) => (
                    <li
                      key={src}
                      className="size-12 overflow-hidden rounded-full border-2 border-white dark:border-gray-900"
                    >
                      <img src={src} alt="" className="size-full object-cover" />
                    </li>
                  ))}
                </ul>

                <p className="mt-4">
                  <span className="block text-5xl font-bold text-gray-900 tabular-nums dark:text-white">
                    {reviewCount}
                  </span>
                  <span className="mt-1 block text-sm text-gray-700 dark:text-gray-300">
                    {reviewLabel}
                  </span>
                </p>

                <p className="mt-6 text-xl font-medium text-balance text-gray-900 dark:text-white">
                  {claim}
                </p>
              </div>

              {/* One control, not a button inside a link. */}
              <a
                href={tutorialsHref}
                className="group flex items-center justify-between gap-3 rounded-3xl bg-purple-300 p-6 transition hover:bg-purple-400 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-purple-700 dark:bg-purple-500/30 dark:hover:bg-purple-500/40"
              >
                <span className="inline-flex min-h-11 items-center gap-2 rounded-full bg-white px-4 text-sm font-medium text-gray-900">
                  <Play aria-hidden="true" className="size-4 fill-current" />
                  {tutorialsLabel}
                </span>
                <span
                  aria-hidden="true"
                  className="inline-flex size-10 flex-none items-center justify-center rounded-full bg-gray-900 transition group-hover:scale-110"
                >
                  <ArrowRight className="size-5 text-white" />
                </span>
              </a>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
