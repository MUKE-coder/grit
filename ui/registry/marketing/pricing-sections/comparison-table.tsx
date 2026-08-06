import { Check, Minus } from 'lucide-react'

/*
 * A full feature comparison, as a real table.
 *
 * Two things here are worth copying and are the reason most comparison tables
 * on the web are unusable.
 *
 * First, the ticks are not decorative. A check icon marked `aria-hidden` in an
 * otherwise empty cell leaves nothing to announce, so a screen reader reads
 * "Single sign-on, blank, blank, blank" and the user cannot tell which plans
 * include it — the exact question the table exists to answer. Every boolean
 * cell carries an `sr-only` "Included" or "Not included" beside the glyph.
 *
 * Second, a four-column table does not fit on a phone, and horizontally
 * scrolling a table means scrolling the row labels out of view. So below `lg`
 * the same data renders as one stacked list per plan, where every row carries
 * its own label. Only one of the two is in the document at a time — `hidden`
 * removes an element from the accessibility tree as well as from the page, so
 * nothing is announced twice.
 *
 * The plan names are `<th scope="col">` and the feature names are
 * `<th scope="row">`, which is what lets a screen reader say "Growth, custom
 * domains, 3" when you land on that cell instead of just "3".
 */

export interface Plan {
  name: string
  monthly: number
  featured?: boolean
  cta?: string
  href?: string
}

export interface FeatureRow {
  label: string
  /** One entry per plan, in the same order. true/false render a glyph. */
  values: (string | boolean)[]
}

export interface FeatureGroup {
  title: string
  rows: FeatureRow[]
}

const PLANS: Plan[] = [
  { name: 'Starter', monthly: 19 },
  { name: 'Growth', monthly: 49, featured: true },
  { name: 'Scale', monthly: 99 },
]

const GROUPS: FeatureGroup[] = [
  {
    title: 'Features',
    rows: [
      { label: 'Edge content delivery', values: [true, true, true] },
      { label: 'Custom domains', values: ['1', '3', 'Unlimited'] },
      { label: 'Team members', values: ['3', '20', 'Unlimited'] },
      { label: 'Single sign-on (SSO)', values: [false, false, true] },
    ],
  },
  {
    title: 'Reporting',
    rows: [
      { label: 'Advanced analytics', values: [true, true, true] },
      { label: 'Basic reports', values: [false, true, true] },
      { label: 'Professional reports', values: [false, false, true] },
      { label: 'Custom report builder', values: [false, false, true] },
    ],
  },
  {
    title: 'Support',
    rows: [
      { label: '24/7 online support', values: [true, true, true] },
      { label: 'Quarterly workshops', values: [false, true, true] },
      { label: 'Priority phone support', values: [false, false, true] },
      { label: '1:1 onboarding tour', values: [false, false, true] },
    ],
  },
]

/* Never just an icon: the sr-only text is the only thing that answers
   "is this included?" for anyone not looking at the glyph. */
function Value({ value }: { value: string | boolean }) {
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
        <Minus aria-hidden="true" className="mx-auto size-5 text-gray-400 dark:text-gray-600" />
        <span className="sr-only">Not included</span>
      </>
    )
  }
  return <span className="text-sm text-gray-600 dark:text-gray-400">{value}</span>
}

function PlanHeader({ plan }: { plan: Plan }) {
  return (
    <>
      <p className="text-sm font-semibold text-gray-900 dark:text-white">{plan.name}</p>
      <p className="mt-2 flex items-baseline gap-x-1">
        <span className="text-4xl font-semibold tracking-tight text-gray-900 tabular-nums dark:text-white">
          ${plan.monthly}
        </span>
        <span className="text-sm text-gray-600 dark:text-gray-400">/month</span>
      </p>
      <a
        href={plan.href ?? '#'}
        className={`mt-6 flex min-h-11 items-center justify-center rounded-lg px-4 text-sm font-semibold focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 ${
          plan.featured
            ? 'bg-indigo-600 text-white hover:bg-indigo-500'
            : 'border border-indigo-200 text-indigo-600 hover:bg-indigo-50 dark:border-indigo-400/30 dark:text-indigo-400 dark:hover:bg-indigo-400/10'
        }`}
      >
        {plan.cta ?? 'Buy plan'}
        <span className="sr-only">, the {plan.name} plan</span>
      </a>
    </>
  )
}

export default function PricingComparisonTable({
  eyebrow = 'Pricing',
  title = 'Pricing that grows with you',
  subtitle = 'An affordable plan packed with the features you need to engage your audience, build loyalty and drive sales.',
  plans = PLANS,
  groups = GROUPS,
}: {
  eyebrow?: string
  title?: string
  subtitle?: string
  plans?: Plan[]
  groups?: FeatureGroup[]
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="mx-auto max-w-3xl text-center">
          <p className="text-base font-semibold text-indigo-600 dark:text-indigo-400">{eyebrow}</p>
          <h2 className="mt-2 text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
            {title}
          </h2>
          <p className="mx-auto mt-6 max-w-2xl text-lg/8 text-pretty text-gray-600 dark:text-gray-400">
            {subtitle}
          </p>
        </div>

        {/* Phones and small tablets: one stack per plan, every row labelled. */}
        <div className="mt-16 space-y-16 lg:hidden">
          {plans.map((plan, planIndex) => (
            <section key={plan.name} aria-labelledby={`plan-${planIndex}`}>
              <div
                className={
                  plan.featured
                    ? 'rounded-2xl bg-gray-50 p-6 dark:bg-white/5'
                    : 'rounded-2xl border border-gray-200 p-6 dark:border-white/10'
                }
              >
                <div id={`plan-${planIndex}`}>
                  <PlanHeader plan={plan} />
                </div>
              </div>

              {groups.map((group) => (
                <div key={group.title} className="mt-8">
                  <h3 className="text-sm font-semibold text-gray-900 dark:text-white">
                    {group.title}
                  </h3>
                  <dl className="mt-3 divide-y divide-gray-200 dark:divide-white/10">
                    {group.rows.map((row) => (
                      <div key={row.label} className="flex items-center justify-between gap-4 py-3">
                        <dt className="text-sm text-gray-600 dark:text-gray-400">{row.label}</dt>
                        <dd className="text-sm font-medium text-gray-900 dark:text-white">
                          <Value value={row.values[planIndex]} />
                        </dd>
                      </div>
                    ))}
                  </dl>
                </div>
              ))}
            </section>
          ))}
        </div>

        {/* Large screens: the real thing. */}
        <div className="mt-16 hidden lg:block">
          <table className="w-full table-fixed border-collapse text-left">
            <caption className="sr-only">Plan comparison</caption>
            <colgroup>
              <col className="w-2/5" />
              {plans.map((plan) => (
                <col key={plan.name} />
              ))}
            </colgroup>

            <thead>
              <tr>
                <td />
                {plans.map((plan) => (
                  <th
                    key={plan.name}
                    scope="col"
                    className={`px-6 pb-8 align-top ${
                      plan.featured ? 'rounded-t-2xl bg-gray-50 dark:bg-white/5' : ''
                    }`}
                  >
                    <PlanHeader plan={plan} />
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
                          plans[i].featured ? 'bg-gray-50 dark:bg-white/5' : ''
                        }`}
                      >
                        <Value value={value} />
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
