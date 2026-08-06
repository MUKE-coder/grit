'use client'

import { useId, useState } from 'react'
import { Check } from 'lucide-react'

/*
 * Three tiers with a monthly/annual switch.
 *
 * The switch is two real radio inputs in a fieldset, not two buttons with
 * click handlers. That is the whole trick and it buys a lot: arrow keys move
 * between the options, the group is announced with its legend before the
 * options are read, the checked one is announced as checked, and the
 * relationship between the two is expressed in the markup rather than implied
 * by them sitting next to each other. A pair of divs with onClick gets none of
 * that, and a screen reader user hears two unrelated buttons.
 *
 * The inputs are hidden with `sr-only`, not `display: none`. Hidden that way
 * they are still focusable and still in the accessibility tree; `hidden` or
 * `display: none` would remove them from both and take the keyboard support
 * with them. The visible pill is the <label>, which is why clicking it works
 * without a handler.
 *
 * Prices are held as numbers and the annual figure is derived, so a discount
 * change is one edit and the two prices cannot drift apart. `tabular-nums`
 * stops the figure jittering horizontally as it changes.
 */

export interface Tier {
  name: string
  description: string
  monthly: number
  featured?: boolean
  /** Rendered above the list, e.g. "Everything in Free, plus:". */
  inherits?: string
  features: string[]
  cta?: string
  href?: string
}

const TIERS: Tier[] = [
  {
    name: 'Free',
    description: 'For developers trying Grit for the first time.',
    monthly: 0,
    features: ['Basic analytics dashboard', '5GB cloud storage', 'Email and chat support'],
  },
  {
    name: 'Pro',
    description: 'For developers who need more features and support.',
    monthly: 14,
    featured: true,
    inherits: 'Everything in Free, plus:',
    features: [
      '5GB cloud storage',
      'Email and chat support',
      'Access to community forum',
      'Single user access',
      'Access to basic templates',
      'Mobile app access',
      'One custom report per month',
      'Monthly product updates',
      'Standard security features',
    ],
  },
  {
    name: 'Startup',
    description: 'For startups that need advanced features and support.',
    monthly: 37,
    inherits: 'Everything in Pro, plus:',
    features: [
      '5GB cloud storage',
      'Email and chat support',
      'Multi-user access',
      'One custom report per month',
      'Monthly product updates',
      'Standard security features',
      'Access to advanced templates',
      'Access to community forum',
      'Mobile app access',
    ],
  },
]

export function BillingToggle({
  annual,
  onChange,
  discount,
}: {
  annual: boolean
  onChange: (annual: boolean) => void
  discount: number
}) {
  const name = useId()
  const options: { label: string; value: boolean }[] = [
    { label: 'Monthly', value: false },
    { label: 'Annually', value: true },
  ]

  return (
    <fieldset className="mt-10 flex flex-col items-center">
      <legend className="sr-only">Billing period</legend>
      <div className="inline-flex rounded-full bg-gray-100 p-1 dark:bg-white/5">
        {options.map((option) => (
          <label
            key={option.label}
            className={`flex min-h-11 cursor-pointer items-center rounded-full px-5 text-sm font-medium transition-colors has-[:focus-visible]:outline-2 has-[:focus-visible]:outline-offset-2 has-[:focus-visible]:outline-indigo-600 ${
              annual === option.value
                ? 'bg-white text-gray-900 shadow-sm dark:bg-gray-800 dark:text-white'
                : 'text-gray-600 dark:text-gray-400'
            }`}
          >
            {/* sr-only, not hidden: it keeps the keyboard and the a11y tree. */}
            <input
              type="radio"
              name={name}
              className="sr-only"
              checked={annual === option.value}
              onChange={() => onChange(option.value)}
            />
            {option.label}
          </label>
        ))}
      </div>
      <p className="mt-4 text-sm text-gray-600 dark:text-gray-400">
        <span className="font-semibold text-indigo-600 dark:text-indigo-400">
          Save {discount}%
        </span>{' '}
        on annual billing
      </p>
    </fieldset>
  )
}

export function TierCard({ tier, annual, discount }: { tier: Tier; annual: boolean; discount: number }) {
  const price = annual ? Math.round(tier.monthly * (1 - discount / 100)) : tier.monthly

  return (
    <div
      className={`flex flex-col rounded-2xl p-8 ${
        tier.featured
          ? 'bg-white shadow-xl ring-1 ring-gray-900/5 dark:bg-gray-900 dark:ring-white/10'
          : ''
      }`}
    >
      <h3 className="text-xl font-semibold text-gray-900 dark:text-white">{tier.name}</h3>
      <p className="mt-2 text-sm/6 text-gray-600 dark:text-gray-400">{tier.description}</p>

      <p className="mt-8 flex items-baseline gap-x-2">
        <span className="text-5xl font-semibold tracking-tight text-gray-900 tabular-nums dark:text-white">
          ${price}
        </span>
      </p>
      <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">Per month</p>

      <a
        href={tier.href ?? '#'}
        className={`mt-8 flex min-h-11 items-center justify-center rounded-lg px-4 text-sm font-semibold focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 ${
          tier.featured
            ? 'bg-indigo-600 text-white hover:bg-indigo-500'
            : 'border border-gray-300 text-gray-900 hover:bg-gray-50 dark:border-white/15 dark:text-white dark:hover:bg-white/5'
        }`}
      >
        {tier.cta ?? 'Get started'}
        <span className="sr-only"> with the {tier.name} plan</span>
      </a>

      {tier.inherits && (
        <p className="mt-8 text-sm font-semibold text-gray-900 dark:text-white">{tier.inherits}</p>
      )}
      <ul role="list" className={`space-y-3 ${tier.inherits ? 'mt-4' : 'mt-8'}`}>
        {tier.features.map((feature) => (
          <li key={feature} className="flex gap-x-3 text-sm/6 text-gray-600 dark:text-gray-400">
            <Check aria-hidden="true" className="mt-1 size-4 flex-none text-gray-900 dark:text-white" />
            {feature}
          </li>
        ))}
      </ul>
    </div>
  )
}

export default function PricingThreeTiersWithToggle({
  title = 'Pricing that scales\nwith your business',
  subtitle = 'Choose the plan that fits, and start shipping today.',
  discount = 25,
  tiers = TIERS,
}: {
  title?: string
  subtitle?: string
  discount?: number
  tiers?: Tier[]
}) {
  const [annual, setAnnual] = useState(true)

  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="text-4xl font-semibold tracking-tight text-balance whitespace-pre-line text-gray-900 sm:text-5xl dark:text-white">
            {title}
          </h2>
          <p className="mt-6 text-lg/8 text-pretty text-gray-600 dark:text-gray-400">{subtitle}</p>
          <BillingToggle annual={annual} onChange={setAnnual} discount={discount} />
        </div>

        <div className="mt-16 rounded-3xl border border-gray-200 p-2 dark:border-white/10">
          <div className="grid grid-cols-1 items-start gap-2 lg:grid-cols-3">
            {tiers.map((tier) => (
              <TierCard key={tier.name} tier={tier} annual={annual} discount={discount} />
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}
