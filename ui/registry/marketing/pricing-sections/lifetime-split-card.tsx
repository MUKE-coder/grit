import { Check } from 'lucide-react'

/*
 * A one-off purchase: what you get on the left, what it costs on the right.
 *
 * A lifetime price has no billing period to switch and no tier to compare
 * against, so a pricing grid has nothing to do here. The split says the same
 * thing the layout of a receipt says: this is the thing, this is the number.
 *
 * "What's included" is a heading with a rule running off it, drawn as a
 * flex-1 <span> rather than a border on the heading. A border would stop where
 * the text stops, which is the opposite of what the design asks for, and an
 * <hr> in the middle of a heading is a separator between nothing and nothing.
 *
 * Dark in both themes by design.
 */

const INCLUDED = [
  'Private forum access',
  'Member resources',
  'Entry to annual conference',
  'Official member t-shirt',
]

export default function PricingLifetimeSplitCard({
  title = 'Simple no-tricks pricing',
  subtitle = 'One payment, no renewal, no seat count to keep an eye on. Everything below is yours for good.',
  planName = 'Lifetime membership',
  planDescription = 'Access to every resource we publish, the private forum, and a seat at the annual conference. Buy it once.',
  priceLabel = 'Pay once, own it forever',
  price = 349,
  currency = 'USD',
  cta = 'Get access',
  href = '#',
  footnote = 'Invoices and receipts available for easy company reimbursement.',
  included = INCLUDED,
}: {
  title?: string
  subtitle?: string
  planName?: string
  planDescription?: string
  priceLabel?: string
  price?: number
  currency?: string
  cta?: string
  href?: string
  footnote?: string
  included?: string[]
}) {
  return (
    <section className="bg-gray-950 py-24 sm:py-32">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="mx-auto max-w-3xl text-center">
          <h2 className="text-4xl font-semibold tracking-tight text-balance text-white sm:text-5xl">
            {title}
          </h2>
          <p className="mt-6 text-lg/8 text-pretty text-gray-400">{subtitle}</p>
        </div>

        <div className="mx-auto mt-16 max-w-5xl rounded-3xl border border-white/10 bg-white/5 p-2">
          <div className="grid grid-cols-1 items-stretch gap-2 lg:grid-cols-[minmax(0,1.4fr)_minmax(0,1fr)]">
            <div className="p-8 sm:p-10">
              <h3 className="text-3xl font-semibold tracking-tight text-white">{planName}</h3>
              <p className="mt-4 text-base/7 text-gray-400">{planDescription}</p>

              <div className="mt-10 flex items-center gap-x-4">
                <h4 className="flex-none text-sm font-semibold text-indigo-400">
                  What&rsquo;s included
                </h4>
                {/* A rule that runs past the heading, which a border on the
                    heading could not do. */}
                <span aria-hidden="true" className="h-px flex-1 bg-white/10" />
              </div>

              <ul
                role="list"
                className="mt-8 grid grid-cols-1 gap-x-8 gap-y-3 sm:grid-cols-2"
              >
                {included.map((item) => (
                  <li key={item} className="flex gap-x-3 text-sm/6 text-gray-300">
                    <Check aria-hidden="true" className="mt-1 size-4 flex-none text-indigo-400" />
                    {item}
                  </li>
                ))}
              </ul>
            </div>

            <div className="flex flex-col items-center justify-center rounded-2xl border border-white/10 bg-gray-900/60 p-8 text-center sm:p-10">
              <p className="text-base font-semibold text-gray-300">{priceLabel}</p>
              <p className="mt-4 flex items-baseline justify-center gap-x-2">
                <span className="text-6xl font-semibold tracking-tight text-white tabular-nums">
                  ${price}
                </span>
                <span className="text-sm font-semibold text-gray-400">{currency}</span>
              </p>
              <a
                href={href}
                className="mt-8 flex min-h-11 w-full items-center justify-center rounded-lg bg-indigo-500 px-4 text-sm font-semibold text-white hover:bg-indigo-400 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-400"
              >
                {cta}
              </a>
              <p className="mt-6 text-xs/5 text-pretty text-gray-500">{footnote}</p>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
