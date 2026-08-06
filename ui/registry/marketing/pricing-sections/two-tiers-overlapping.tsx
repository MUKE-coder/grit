import { Check } from 'lucide-react'

/*
 * Two tiers, overlapping, with the featured one lifted out of the page.
 *
 * The overlap is done with a negative margin that only exists above `lg`, and
 * with `lg:rounded-r-none` on the card underneath so the seam reads as one
 * card sitting on another rather than two cards colliding. On a phone the
 * overlap is dropped entirely and the cards stack: two overlapping cards in a
 * single narrow column is just one card partly hiding the other.
 *
 * Set `featured` on whichever tier you are pushing. The featured card renders
 * dark with a solid button, the other renders light with an outline; swap the
 * flag and the emphasis moves with it, no class editing required.
 *
 * The card is a <div>, not a link, and only the button inside it is clickable.
 * A whole card wrapped in an anchor announces its entire contents, price and
 * every feature, as one enormous link name.
 */

export interface Tier {
  name: string
  monthly: number
  description: string
  features: string[]
  featured?: boolean
  cta?: string
  href?: string
}

const TIERS: Tier[] = [
  {
    name: 'Hobby',
    monthly: 29,
    description: "The perfect plan if you're just getting started with our product.",
    features: [
      '25 products',
      'Up to 10,000 subscribers',
      'Advanced analytics',
      '24-hour support response time',
    ],
  },
  {
    name: 'Enterprise',
    monthly: 99,
    description: 'Dedicated support and infrastructure for your company.',
    featured: true,
    features: [
      'Unlimited products',
      'Unlimited subscribers',
      'Advanced analytics',
      'Dedicated support representative',
      'Marketing automations',
      'Custom integrations',
    ],
  },
]

export default function PricingTwoTiersOverlapping({
  eyebrow = 'Pricing',
  title = 'Choose the right plan for you',
  subtitle = 'An affordable plan packed with the features you need to engage your audience, build loyalty and drive sales.',
  tiers = TIERS,
}: {
  eyebrow?: string
  title?: string
  subtitle?: string
  tiers?: Tier[]
}) {
  return (
    <section className="relative isolate overflow-hidden bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 -z-10"
        style={{
          background:
            'radial-gradient(48rem 28rem at 78% 12%, rgba(168,120,220,0.18), transparent 70%)',
        }}
      />

      <div className="mx-auto max-w-5xl px-6 lg:px-8">
        <div className="mx-auto max-w-3xl text-center">
          <p className="text-base font-semibold text-indigo-600 dark:text-indigo-400">{eyebrow}</p>
          <h2 className="mt-2 text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-6xl dark:text-white">
            {title}
          </h2>
          <p className="mx-auto mt-6 max-w-2xl text-lg/8 text-pretty text-gray-600 dark:text-gray-400">
            {subtitle}
          </p>
        </div>

        <div className="mt-20 grid grid-cols-1 items-center gap-8 lg:grid-cols-2 lg:gap-0">
          {tiers.map((tier) => (
            <div
              key={tier.name}
              className={
                tier.featured
                  ? 'relative z-10 rounded-3xl bg-gray-900 p-8 shadow-2xl sm:p-10 lg:-ml-8 dark:bg-gray-900 dark:ring-1 dark:ring-white/10'
                  : 'rounded-3xl bg-white/60 p-8 ring-1 ring-gray-900/10 sm:p-10 lg:mr-0 lg:rounded-r-none lg:py-14 dark:bg-white/5 dark:ring-white/10'
              }
            >
              <h3
                className={`text-base font-semibold ${
                  tier.featured ? 'text-indigo-400' : 'text-indigo-600 dark:text-indigo-400'
                }`}
              >
                {tier.name}
              </h3>

              <p className="mt-4 flex items-baseline gap-x-2">
                <span
                  className={`text-5xl font-semibold tracking-tight tabular-nums ${
                    tier.featured ? 'text-white' : 'text-gray-900 dark:text-white'
                  }`}
                >
                  ${tier.monthly}
                </span>
                <span
                  className={`text-base ${
                    tier.featured ? 'text-gray-400' : 'text-gray-600 dark:text-gray-400'
                  }`}
                >
                  /month
                </span>
              </p>

              <p
                className={`mt-6 text-base/7 ${
                  tier.featured ? 'text-gray-300' : 'text-gray-600 dark:text-gray-400'
                }`}
              >
                {tier.description}
              </p>

              <ul role="list" className="mt-8 space-y-3">
                {tier.features.map((feature) => (
                  <li
                    key={feature}
                    className={`flex gap-x-3 text-sm/6 ${
                      tier.featured ? 'text-gray-300' : 'text-gray-600 dark:text-gray-400'
                    }`}
                  >
                    <Check
                      aria-hidden="true"
                      className="mt-1 size-4 flex-none text-indigo-500 dark:text-indigo-400"
                    />
                    {feature}
                  </li>
                ))}
              </ul>

              <a
                href={tier.href ?? '#'}
                className={`mt-10 flex min-h-11 items-center justify-center rounded-lg px-4 text-sm font-semibold focus-visible:outline-2 focus-visible:outline-offset-2 ${
                  tier.featured
                    ? 'bg-indigo-500 text-white hover:bg-indigo-400 focus-visible:outline-indigo-400'
                    : 'border border-indigo-200 text-indigo-600 hover:bg-indigo-50 focus-visible:outline-indigo-600 dark:border-indigo-400/30 dark:text-indigo-400 dark:hover:bg-indigo-400/10'
                }`}
              >
                {tier.cta ?? 'Get started today'}
                <span className="sr-only"> with the {tier.name} plan</span>
              </a>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
