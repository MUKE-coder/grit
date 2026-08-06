'use client'

import { useId, useState } from 'react'
import { Check, X } from 'lucide-react'

/*
 * A dark hero with the plan cards straddling its lower edge, and the full
 * comparison below on light.
 *
 * The straddle is the whole idea: the cards sit half in the dark band and half
 * out of it, which ties the summary at the top to the detail underneath
 * instead of leaving two disconnected sections. It is built with a negative
 * bottom margin on the dark band rather than absolute positioning, so the
 * cards still take part in normal flow and the light section below moves down
 * when a card grows. Position them absolutely and a longer feature list
 * silently overlaps the table.
 *
 * The middle card is deliberately taller and lighter. On a dark band, white is
 * the loudest thing available, and it is doing the same job as the tallest bar
 * in a chart.
 *
 * The comparison table is a real table with `scope` on its headers, and every
 * boolean cell carries an `sr-only` "Included" / "Not included": a check icon
 * on its own leaves an empty cell for anyone not looking at it, which defeats
 * the one question the table answers.
 */

export interface Plan {
  name: string
  monthly: number
  blurb: string
  features: string[]
  featured?: boolean
  cta?: string
  href?: string
}

export interface FeatureRow {
  label: string
  values: (string | boolean)[]
}

export interface FeatureGroup {
  title: string
  rows: FeatureRow[]
}

/* Order is display order. The featured plan sits in the middle. */
const PLANS: Plan[] = [
  {
    name: 'Starter',
    monthly: 19,
    blurb: 'Everything you need to get started.',
    features: ['Custom domains', 'Edge content delivery', 'Advanced analytics'],
  },
  {
    name: 'Scale',
    monthly: 99,
    blurb: 'Added flexibility at scale.',
    featured: true,
    features: [
      'Custom domains',
      'Edge content delivery',
      'Advanced analytics',
      'Quarterly workshops',
      'Single sign-on (SSO)',
      'Priority phone support',
    ],
  },
  {
    name: 'Growth',
    monthly: 49,
    blurb: 'All the extras for your growing team.',
    features: [
      'Custom domains',
      'Edge content delivery',
      'Advanced analytics',
      'Quarterly workshops',
    ],
  },
]

/* Values are in the same order as PLANS. */
const GROUPS: FeatureGroup[] = [
  {
    title: 'Features',
    rows: [
      { label: 'Edge content delivery', values: [true, true, true] },
      { label: 'Custom domains', values: ['1', 'Unlimited', '3'] },
      { label: 'Team members', values: ['3', 'Unlimited', '20'] },
      { label: 'Single sign-on (SSO)', values: [false, true, false] },
    ],
  },
  {
    title: 'Reporting',
    rows: [
      { label: 'Advanced analytics', values: [true, true, true] },
      { label: 'Basic reports', values: [false, true, true] },
      { label: 'Professional reports', values: [false, true, false] },
      { label: 'Custom report builder', values: [false, true, false] },
    ],
  },
  {
    title: 'Support',
    rows: [
      { label: '24/7 online support', values: [true, true, true] },
      { label: 'Quarterly workshops', values: [false, true, true] },
      { label: 'Priority phone support', values: [false, true, false] },
      { label: '1:1 onboarding tour', values: [false, true, false] },
    ],
  },
]

function Value({ value, featured }: { value: string | boolean; featured?: boolean }) {
  if (value === true) {
    return (
      <>
        <Check aria-hidden="true" className="mx-auto size-5 text-indigo-600 dark:text-indigo-400" />
        <span className="sr-only">Included</span>
      </>
    )
  }
  if (value === false) {
    return (
      <>
        <X aria-hidden="true" className="mx-auto size-5 text-gray-400 dark:text-gray-600" />
        <span className="sr-only">Not included</span>
      </>
    )
  }
  return (
    <span
      className={`text-sm ${
        featured
          ? 'font-medium text-indigo-600 dark:text-indigo-400'
          : 'text-gray-600 dark:text-gray-400'
      }`}
    >
      {value}
    </span>
  )
}

export default function PricingDarkCardsOverComparison({
  title = 'Pricing that grows with you',
  subtitle = 'An affordable plan packed with the features you need to engage your audience, build loyalty and drive sales.',
  discount = 20,
  plans = PLANS,
  groups = GROUPS,
}: {
  title?: string
  subtitle?: string
  discount?: number
  plans?: Plan[]
  groups?: FeatureGroup[]
}) {
  const [annual, setAnnual] = useState(false)
  const billingName = useId()

  const priceOf = (plan: Plan) =>
    annual ? Math.round(plan.monthly * (1 - discount / 100)) : plan.monthly

  return (
    <section className="bg-white dark:bg-gray-950">
      {/* The band ends here; the negative margin is what the cards straddle. */}
      <div className="relative isolate overflow-hidden bg-gray-950 pt-24 pb-40 sm:pt-32">
        <div
          aria-hidden="true"
          className="pointer-events-none absolute inset-0 -z-10"
          style={{
            background:
              'radial-gradient(46rem 24rem at 50% 96%, rgba(139,92,246,0.4), transparent 68%)',
          }}
        />
        <div className="mx-auto max-w-3xl px-6 text-center lg:px-8">
          <h2 className="text-4xl font-semibold tracking-tight text-balance text-white sm:text-5xl">
            {title}
          </h2>
          <p className="mx-auto mt-6 max-w-2xl text-lg/8 text-pretty text-gray-400">{subtitle}</p>

          <fieldset className="mt-10 flex justify-center">
            <legend className="sr-only">Billing period</legend>
            <div className="inline-flex rounded-full bg-white/10 p-1">
              {[
                { label: 'Monthly', value: false },
                { label: 'Annually', value: true },
              ].map((option) => (
                <label
                  key={option.label}
                  className={`flex min-h-11 cursor-pointer items-center rounded-full px-5 text-sm font-medium transition-colors has-[:focus-visible]:outline-2 has-[:focus-visible]:outline-offset-2 has-[:focus-visible]:outline-white ${
                    annual === option.value ? 'bg-indigo-500 text-white' : 'text-gray-300'
                  }`}
                >
                  {/* sr-only, not hidden: keeps the keyboard and the a11y tree. */}
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
          </fieldset>
        </div>
      </div>

      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        {/* Pulled up into the band. Normal flow, so a longer card pushes the
            table down rather than landing on top of it.

            `relative z-10` on this container is load-bearing. The band above is
            positioned (it has `isolate`), and a positioned element paints above
            a static one regardless of document order, so without this the band
            covers the part of the cards that overlaps it. */}
        <div className="relative z-10 -mt-24 grid grid-cols-1 items-start gap-6 lg:grid-cols-3 lg:gap-0">
          {plans.map((plan) => (
            <div
              key={plan.name}
              className={
                plan.featured
                  ? 'relative z-10 rounded-2xl bg-white p-8 shadow-2xl lg:-mt-12 dark:bg-gray-900'
                  : 'rounded-2xl bg-gray-900 p-8 ring-1 ring-white/10 lg:my-6'
              }
            >
              <p
                className={`text-sm font-semibold ${
                  plan.featured ? 'text-gray-900 dark:text-white' : 'text-white'
                }`}
              >
                {plan.name}
              </p>
              <p className="mt-3 flex items-baseline gap-x-2">
                <span
                  className={`text-4xl font-semibold tracking-tight tabular-nums ${
                    plan.featured ? 'text-gray-900 dark:text-white' : 'text-white'
                  }`}
                >
                  ${priceOf(plan)}
                </span>
                <span
                  className={`text-left text-xs ${
                    plan.featured ? 'text-gray-600 dark:text-gray-400' : 'text-gray-400'
                  }`}
                >
                  USD
                  <br />
                  Billed {annual ? 'annually' : 'monthly'}
                </span>
              </p>

              <a
                href={plan.href ?? '#'}
                className={`mt-6 flex min-h-11 items-center justify-center rounded-lg px-4 text-sm font-semibold focus-visible:outline-2 focus-visible:outline-offset-2 ${
                  plan.featured
                    ? 'bg-indigo-600 text-white hover:bg-indigo-500 focus-visible:outline-indigo-600'
                    : 'bg-white/10 text-white hover:bg-white/15 focus-visible:outline-white'
                }`}
              >
                {plan.cta ?? 'Buy this plan'}
                <span className="sr-only"> — the {plan.name} plan</span>
              </a>

              <ul role="list" className="mt-6 space-y-3">
                {plan.features.map((feature) => (
                  <li
                    key={feature}
                    className={`flex gap-x-3 border-b pb-3 text-sm/6 last:border-0 ${
                      plan.featured
                        ? 'border-gray-200 text-gray-700 dark:border-white/10 dark:text-gray-300'
                        : 'border-white/10 text-gray-300'
                    }`}
                  >
                    <Check
                      aria-hidden="true"
                      className={`mt-1 size-4 flex-none ${
                        plan.featured ? 'text-indigo-600 dark:text-indigo-400' : 'text-gray-500'
                      }`}
                    />
                    {feature}
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        {/* Below lg the table becomes one labelled stack per plan: four
            columns do not fit on a phone, and scrolling a table sideways
            scrolls the row labels out of view. Only one version is in the
            document at a time, so nothing is announced twice. */}
        <div className="mt-20 space-y-14 pb-24 lg:hidden">
          {plans.map((plan, planIndex) => (
            <section key={plan.name} aria-labelledby={`compare-${planIndex}`}>
              <h3
                id={`compare-${planIndex}`}
                className="text-sm font-semibold text-indigo-600 dark:text-indigo-400"
              >
                {plan.name}
              </h3>
              <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">{plan.blurb}</p>
              {groups.map((group) => (
                <div key={group.title} className="mt-6">
                  <h4 className="text-sm font-semibold text-gray-900 dark:text-white">
                    {group.title}
                  </h4>
                  <dl className="mt-2 divide-y divide-gray-200 dark:divide-white/10">
                    {group.rows.map((row) => (
                      <div key={row.label} className="flex items-center justify-between gap-4 py-3">
                        <dt className="text-sm text-gray-600 dark:text-gray-400">{row.label}</dt>
                        <dd>
                          <Value value={row.values[planIndex]} featured={plan.featured} />
                        </dd>
                      </div>
                    ))}
                  </dl>
                </div>
              ))}
            </section>
          ))}
        </div>

        <div className="mt-24 hidden pb-24 lg:block">
          <table className="w-full table-fixed border-collapse text-left">
            <caption className="sr-only">Plan comparison</caption>
            <colgroup>
              <col className="w-1/4" />
              {plans.map((plan) => (
                <col key={plan.name} />
              ))}
            </colgroup>
            <thead>
              <tr>
                <td />
                {plans.map((plan) => (
                  <th key={plan.name} scope="col" className="px-6 pb-6 align-top">
                    <p
                      className={`text-sm font-semibold ${
                        plan.featured
                          ? 'text-indigo-600 dark:text-indigo-400'
                          : 'text-gray-900 dark:text-white'
                      }`}
                    >
                      {plan.name}
                    </p>
                    <p className="mt-1 text-sm font-normal text-gray-600 dark:text-gray-400">
                      {plan.blurb}
                    </p>
                  </th>
                ))}
              </tr>
            </thead>

            {groups.map((group) => (
              <tbody key={group.title}>
                <tr>
                  <th
                    scope="colgroup"
                    colSpan={plans.length + 1}
                    className="pt-10 pb-4 text-sm font-semibold text-gray-900 dark:text-white"
                  >
                    {group.title}
                  </th>
                </tr>
                {group.rows.map((row) => (
                  <tr key={row.label} className="border-t border-gray-200 dark:border-white/10">
                    <th
                      scope="row"
                      className="py-4 pr-6 text-sm font-normal text-gray-600 dark:text-gray-400"
                    >
                      {row.label}
                    </th>
                    {row.values.map((value, i) => (
                      <td
                        key={plans[i].name}
                        className={`px-6 py-4 text-center ${
                          plans[i].featured
                            ? 'bg-indigo-50/60 dark:bg-indigo-400/10'
                            : ''
                        }`}
                      >
                        <Value value={value} featured={plans[i].featured} />
                      </td>
                    ))}
                  </tr>
                ))}
              </tbody>
            ))}
          </table>
        </div>
      </div>
    </section>
  )
}
