'use client'

import { useId, useState } from 'react'
import { Check } from 'lucide-react'

/*
 * One plan, one price, with what you get listed underneath.
 *
 * If you only sell one thing this beats three columns where two are decoys.
 * The tiering exists to make a middle option look reasonable; without that job
 * to do, a single card removes the comparison step entirely and the visitor
 * goes straight from price to button.
 *
 * The features sit below the card rather than inside it, which keeps the price
 * and the button in one glance instead of pushing the button below a list of
 * ten items. Someone who is sold does not have to scroll past the reasons.
 *
 * The switch is two real radios in a fieldset: arrow keys work, the group is
 * announced with its legend, and the checked option says so. They are sr-only
 * rather than hidden, which keeps them focusable and in the accessibility tree.
 */

const FEATURES = [
  'Team collaboration',
  'Custom templates',
  '24/7 customer support',
  'API access',
  'White labelling',
  'SSO integration',
  'Dedicated account manager',
  'Custom reporting',
]

export default function PricingSinglePlanWithFeatureGrid({
  title = 'One simple plan,\none price',
  subtitle = 'Everything you need to ship, in one package.',
  planName = 'All-in-one',
  planDescription = 'Everything you need in one simple plan',
  monthly = 349,
  discount = 25,
  cta = 'Get started now',
  href = '#',
  footnote = 'No hidden fees. Cancel anytime. Invoices available for easy reimbursement.',
  features = FEATURES,
}: {
  title?: string
  subtitle?: string
  planName?: string
  planDescription?: string
  monthly?: number
  discount?: number
  cta?: string
  href?: string
  footnote?: string
  features?: string[]
}) {
  const [annual, setAnnual] = useState(true)
  const billingName = useId()
  const price = annual ? Math.round(monthly * (1 - discount / 100)) : monthly

  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-2xl px-6 lg:px-8">
        <div className="text-center">
          <h2 className="text-4xl font-semibold tracking-tight text-balance whitespace-pre-line text-gray-900 sm:text-5xl dark:text-white">
            {title}
          </h2>
          <p className="mt-6 text-lg/8 text-pretty text-gray-600 dark:text-gray-400">{subtitle}</p>

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

        <div className="mt-14 rounded-3xl bg-white p-8 text-center shadow-xl ring-1 ring-gray-900/5 sm:p-10 dark:bg-gray-900 dark:ring-white/10">
          <h3 className="text-xl font-semibold text-gray-900 dark:text-white">{planName}</h3>
          <p className="mt-2 text-sm/6 text-gray-600 dark:text-gray-400">{planDescription}</p>

          <div className="mt-8 flex items-center justify-center gap-x-3">
            <span className="text-6xl font-semibold tracking-tight text-gray-900 tabular-nums dark:text-white">
              ${price}
            </span>
            <span className="text-left text-sm text-gray-600 dark:text-gray-400">
              Per month
              <br />
              {annual ? 'Billed annually' : 'Billed monthly'}
            </span>
          </div>

          <a
            href={href}
            className="mt-8 inline-flex min-h-11 items-center justify-center rounded-lg bg-indigo-600 px-6 text-sm font-semibold text-white hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
          >
            {cta}
          </a>

          <hr className="mt-8 border-gray-200 dark:border-white/10" />
          <p className="mt-6 text-sm/6 text-pretty text-gray-600 dark:text-gray-400">{footnote}</p>
        </div>

        <ul
          role="list"
          className="mx-auto mt-12 grid max-w-xl grid-cols-1 gap-x-8 gap-y-4 sm:grid-cols-2"
        >
          {features.map((feature) => (
            <li key={feature} className="flex items-start gap-x-3">
              <span
                aria-hidden="true"
                className="mt-0.5 flex size-5 flex-none items-center justify-center rounded-full bg-emerald-100 dark:bg-emerald-500/15"
              >
                <Check className="size-3 text-emerald-700 dark:text-emerald-400" />
              </span>
              <span className="text-base text-gray-700 dark:text-gray-300">{feature}</span>
            </li>
          ))}
        </ul>
      </div>
    </section>
  )
}
