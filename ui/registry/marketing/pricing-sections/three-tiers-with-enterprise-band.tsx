'use client'

import { useId, useState } from 'react'
import { Check } from 'lucide-react'

/*
 * The three-tier layout with a self-serve ceiling, plus a band underneath for
 * the plan that does not have a price.
 *
 * The band is a different shape on purpose. Enterprise is not a fourth column:
 * it has no number to compare, its features are a wide two-column list rather
 * than a stack, and its action is "contact sales" rather than "get started".
 * Forcing it into the grid invites a reader to compare it against tiers it
 * cannot be compared with, and leaves a card with a blank where the price
 * should be.
 *
 * Everything is duplicated in this file rather than imported from the sibling
 * block, because the registry installs one file. A shared import would land in
 * a project with nothing to resolve.
 *
 * The billing switch is two real radios in a fieldset; see the note on
 * BillingToggle below for why that matters.
 */

export interface Tier {
  name: string
  description: string
  monthly: number
  featured?: boolean
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

const ENTERPRISE_FEATURES = [
  'One custom report per month',
  'Standard security features',
  'Access to advanced templates',
  'Access to community forum',
  'Mobile app access',
  'Custom invoicing',
  'Custom user roles',
  'Enhanced reporting',
  'Priority support',
]

export default function PricingThreeTiersWithEnterpriseBand({
  title = 'Pricing that scales\nwith your business',
  subtitle = 'Choose the plan that fits, and start shipping today.',
  discount = 25,
  tiers = TIERS,
  enterpriseName = 'Enterprise custom plan',
  enterpriseDescription = 'For large organisations with complex workflows and advanced reporting requirements.',
  enterpriseFeatures = ENTERPRISE_FEATURES,
  enterpriseCta = 'Contact sales',
  enterpriseHref = '#',
}: {
  title?: string
  subtitle?: string
  discount?: number
  tiers?: Tier[]
  enterpriseName?: string
  enterpriseDescription?: string
  enterpriseFeatures?: string[]
  enterpriseCta?: string
  enterpriseHref?: string
}) {
  const [annual, setAnnual] = useState(true)
  const billingName = useId()

  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="text-4xl font-semibold tracking-tight text-balance whitespace-pre-line text-gray-900 sm:text-5xl dark:text-white">
            {title}
          </h2>
          <p className="mt-6 text-lg/8 text-pretty text-gray-600 dark:text-gray-400">{subtitle}</p>

          {/* Two real radios in a fieldset: arrow keys work, the group is
              announced with its legend, and the checked one says so. The inputs
              are sr-only rather than hidden, which keeps them focusable and in
              the accessibility tree. The visible pill is the label, so clicking
              it works without a handler. */}
          <fieldset className="mt-10 flex flex-col items-center">
            <legend className="sr-only">Billing period</legend>
            <div className="inline-flex rounded-full bg-gray-100 p-1 dark:bg-white/5">
              {[
                { label: 'Monthly', value: false },
                { label: 'Annually', value: true },
              ].map((option) => (
                <label
                  key={option.label}
                  className={`flex min-h-11 cursor-pointer items-center rounded-full px-5 text-sm font-medium transition-colors has-[:focus-visible]:outline-2 has-[:focus-visible]:outline-offset-2 has-[:focus-visible]:outline-indigo-600 ${
                    annual === option.value
                      ? 'bg-white text-gray-900 shadow-sm dark:bg-gray-800 dark:text-white'
                      : 'text-gray-600 dark:text-gray-400'
                  }`}
                >
                  <input
                    type="radio"
                    name={billingName}
                    className="sr-only"
                    checked={annual === option.value}
                    onChange={() => setAnnual(option.value)}
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
        </div>

        <div className="mt-16 rounded-3xl border border-gray-200 p-2 dark:border-white/10">
          <div className="grid grid-cols-1 items-start gap-2 lg:grid-cols-3">
            {tiers.map((tier) => {
              const price = annual
                ? Math.round(tier.monthly * (1 - discount / 100))
                : tier.monthly
              return (
                <div
                  key={tier.name}
                  className={`flex flex-col rounded-2xl p-8 ${
                    tier.featured
                      ? 'bg-white shadow-xl ring-1 ring-gray-900/5 dark:bg-gray-900 dark:ring-white/10'
                      : ''
                  }`}
                >
                  <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
                    {tier.name}
                  </h3>
                  <p className="mt-2 text-sm/6 text-gray-600 dark:text-gray-400">
                    {tier.description}
                  </p>
                  <p className="mt-8 text-5xl font-semibold tracking-tight text-gray-900 tabular-nums dark:text-white">
                    ${price}
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
                    <p className="mt-8 text-sm font-semibold text-gray-900 dark:text-white">
                      {tier.inherits}
                    </p>
                  )}
                  <ul role="list" className={`space-y-3 ${tier.inherits ? 'mt-4' : 'mt-8'}`}>
                    {tier.features.map((feature) => (
                      <li
                        key={feature}
                        className="flex gap-x-3 text-sm/6 text-gray-600 dark:text-gray-400"
                      >
                        <Check
                          aria-hidden="true"
                          className="mt-1 size-4 flex-none text-gray-900 dark:text-white"
                        />
                        {feature}
                      </li>
                    ))}
                  </ul>
                </div>
              )
            })}
          </div>
        </div>

        {/* No price, so no column. The divider is a border on the right half,
            drawn only where the two halves sit side by side. */}
        <div className="mt-8 rounded-3xl border border-gray-200 dark:border-white/10">
          <div className="grid grid-cols-1 lg:grid-cols-[minmax(0,22rem)_minmax(0,1fr)]">
            <div className="p-8 sm:p-10">
              <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
                {enterpriseName}
              </h3>
              <p className="mt-2 text-sm/6 text-gray-600 dark:text-gray-400">
                {enterpriseDescription}
              </p>
              <a
                href={enterpriseHref}
                className="mt-8 inline-flex min-h-11 items-center justify-center rounded-lg border border-gray-300 px-5 text-sm font-semibold text-gray-900 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:border-white/15 dark:text-white dark:hover:bg-white/5"
              >
                {enterpriseCta}
              </a>
            </div>
            <div className="p-8 sm:p-10 lg:border-l lg:border-gray-200 dark:lg:border-white/10">
              <ul role="list" className="grid grid-cols-1 gap-x-8 gap-y-3 sm:grid-cols-2">
                {enterpriseFeatures.map((feature) => (
                  <li
                    key={feature}
                    className="flex gap-x-3 text-sm/6 text-gray-600 dark:text-gray-400"
                  >
                    <Check
                      aria-hidden="true"
                      className="mt-1 size-4 flex-none text-gray-900 dark:text-white"
                    />
                    {feature}
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
