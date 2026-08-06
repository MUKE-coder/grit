import { CircleCheck } from 'lucide-react'

/*
 * A recruiting CTA: photo on the left, what you get on the right.
 *
 * The benefits are a real <ul>, so a screen reader announces "list, six items"
 * before reading them. Six perks that arrive as an unannounced run of
 * paragraphs is the same information with the shape stripped out of it, and the
 * shape is most of what makes a list scannable.
 *
 * The tick icons are `aria-hidden` and carry no label. Every row has one, so a
 * tick means nothing beyond "this is a list item" — which the list already
 * says. Announcing "check, check, check" six times is noise, not access.
 *
 * The photo is `object-cover` in a fixed aspect ratio rather than free-height,
 * so swapping in a portrait shot cannot stretch the card and shove the link
 * below the fold.
 */

export interface Perk {
  label: string
}

/* Demo photography from Unsplash, free for commercial use with no attribution
   required, checked to return 200. Sized in the query string so the browser
   fetches roughly what it paints. Replace it: your own team is the entire
   point of this block, and a stock photo of strangers on a beach says the
   opposite of what a careers page is trying to say. */
const PHOTO =
  'https://images.unsplash.com/photo-1529156069898-49953e39b3ac?w=1000&h=1000&fit=crop&q=80'

const PERKS: Perk[] = [
  { label: 'Competitive salaries' },
  { label: 'Flexible work hours' },
  { label: '30 days of paid vacation' },
  { label: 'Annual team retreats' },
  { label: 'Benefits for you and your family' },
  { label: 'A great work environment' },
]

export default function CtaJoinTheTeam({
  title = 'Join our team',
  subtitle = 'We are a small group of engineers and designers who like shipping. If that sounds like your kind of week, we are hiring across the stack.',
  linkLabel = 'See our job postings',
  linkHref = '#',
  photo = PHOTO,
  perks = PERKS,
}: {
  title?: string
  subtitle?: string
  linkLabel?: string
  linkHref?: string
  photo?: string
  perks?: Perk[]
}) {
  return (
    <section className="bg-gray-950 py-24 sm:py-32">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div
          className="rounded-3xl border border-white/10 p-8 sm:p-12"
          style={{
            background:
              'radial-gradient(70rem 40rem at 85% 10%, rgba(59,90,150,0.35), transparent 70%)',
          }}
        >
          <div className="grid grid-cols-1 items-center gap-10 lg:grid-cols-[minmax(0,22rem)_minmax(0,1fr)] lg:gap-16">
            <img
              src={photo}
              alt=""
              className="aspect-square w-full rounded-2xl object-cover shadow-xl"
            />

            <div>
              <h2 className="text-4xl font-semibold tracking-tight text-balance text-white sm:text-5xl">
                {title}
              </h2>
              <p className="mt-4 max-w-xl text-lg/8 text-pretty text-gray-400">{subtitle}</p>

              <ul
                role="list"
                className="mt-10 grid grid-cols-1 gap-x-8 gap-y-4 sm:grid-cols-2"
              >
                {perks.map((perk) => (
                  <li key={perk.label} className="flex items-start gap-3">
                    <CircleCheck
                      aria-hidden="true"
                      className="mt-0.5 size-5 flex-none text-white"
                    />
                    <span className="text-base text-gray-200">{perk.label}</span>
                  </li>
                ))}
              </ul>

              <a
                href={linkHref}
                className="mt-10 inline-flex min-h-11 items-center rounded-md px-1 text-sm font-semibold text-indigo-400 hover:text-indigo-300 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-400"
              >
                {linkLabel} <span aria-hidden="true">&nbsp;&rarr;</span>
              </a>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
