/*
 * FAQs in one column: heading, grouped questions, a contact line at the end.
 *
 * The narrowest of these layouts, and the one to reach for when the FAQ is the
 * whole page rather than a section of one. Capped at a readable measure rather
 * than run to the full container width, because a question that wraps across
 * 1400px is harder to scan than one that wraps at 700.
 *
 * <details> and <summary> rather than React state: the browser already handles
 * focus, Enter and Space, the announced expanded state, and opening a closed
 * section when Ctrl+F matches text inside it. See split-with-contact.tsx.
 */

interface QA {
  q: string
  a: string
}

interface Group {
  title: string
  items: QA[]
}

const GROUPS: Group[] = [
  {
    title: 'General',
    items: [
      {
        q: 'How long does shipping take?',
        a: 'Orders leave the warehouse within one working day. Delivery is two to five working days domestically, and seven to fourteen internationally.',
      },
      {
        q: 'What payment methods do you accept?',
        a: 'All major cards, Apple Pay, Google Pay and bank transfer for invoiced accounts. Cards are handled by our payment provider and never touch our servers.',
      },
      {
        q: 'Can I change or cancel my order?',
        a: 'Any time before it ships. Once it has a tracking number it is with the courier, and the returns process is the faster route.',
      },
    ],
  },
  {
    title: 'Shipping',
    items: [
      {
        q: 'Do you ship internationally?',
        a: 'To most countries. Duties and import taxes are calculated at checkout so the price you pay is the price you pay.',
      },
      {
        q: 'What is your return policy?',
        a: 'Thirty days, unused and in the original packaging. Return shipping is on us if the item arrived damaged or was not what you ordered.',
      },
      {
        q: 'How do I track my order?',
        a: 'The dispatch email carries a tracking link. It also lives in your account under Orders, which is more reliable than an email you have to find again.',
      },
    ],
  },
]

function Item({ item }: { item: QA }) {
  return (
    <details className="group border-b border-gray-200 dark:border-white/10">
      <summary className="flex min-h-14 cursor-pointer list-none items-center justify-between gap-6 py-4 text-base font-semibold text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-white">
        {item.q}
        <svg
          aria-hidden="true"
          viewBox="0 0 20 20"
          fill="none"
          className="size-5 flex-none text-gray-400 transition-transform duration-200 group-open:rotate-180"
        >
          <path
            d="M6 8l4 4 4-4"
            stroke="currentColor"
            strokeWidth="1.5"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
        </svg>
      </summary>
      <p className="pb-5 pr-10 text-base/7 text-gray-600 dark:text-gray-400">{item.a}</p>
    </details>
  )
}

export default function FaqsStackedGrouped({
  title = 'Frequently asked questions',
  subtitle = 'Quick answers to what people ask most about the platform, the service and how it all fits together.',
  contactHref = '#',
  groups = GROUPS,
}: {
  title?: string
  subtitle?: string
  contactHref?: string
  groups?: Group[]
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-3xl px-6 lg:px-8">
        <h2 className="text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
          {title}
        </h2>
        <p className="mt-6 text-lg/8 text-gray-600 dark:text-gray-400">{subtitle}</p>

        <div className="mt-16 space-y-14">
          {groups.map((group) => (
            <div key={group.title}>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                {group.title}
              </h3>
              <div className="mt-4">
                {group.items.map((item) => (
                  <Item key={item.q} item={item} />
                ))}
              </div>
            </div>
          ))}
        </div>

        <p className="mt-16 text-base/7 text-gray-600 dark:text-gray-400">
          Cannot find what you are looking for? Contact our{' '}
          <a
            href={contactHref}
            className="font-medium text-indigo-600 underline-offset-4 hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-indigo-400"
          >
            customer support team
          </a>
          .
        </p>
      </div>
    </section>
  )
}
