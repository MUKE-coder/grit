/*
 * FAQs with a category rail beside the questions.
 *
 * The rail is anchor links, not tabs, and that is deliberate. Tabs hide every
 * category you are not looking at, which means Ctrl+F finds nothing and a
 * search engine indexes a third of the page. Anchors keep all of it in the
 * document and just move you to the right part of it.
 *
 * `scroll-mt-24` on each heading stops the sticky rail landing on top of the
 * heading it just jumped to, which is the usual reason anchor navigation feels
 * broken. `scroll-smooth` is applied on the container rather than globally, and
 * a reduced-motion preference turns it off: for someone who gets motion sick,
 * an unrequested smooth scroll is the problem, not the polish.
 */

interface QA {
  q: string
  a: string
}

interface Group {
  id: string
  title: string
  items: QA[]
}

const GROUPS: Group[] = [
  {
    id: 'general',
    title: 'General',
    items: [
      {
        q: 'How long does shipping take?',
        a: 'Orders leave the warehouse within one working day. Delivery is two to five working days domestically, and seven to fourteen internationally.',
      },
      {
        q: 'What payment methods do you accept?',
        a: 'All major cards, Apple Pay, Google Pay and bank transfer for invoiced accounts.',
      },
      {
        q: 'Can I change or cancel my order?',
        a: 'Any time before it ships. Once it has a tracking number it is with the courier, and the returns process is the faster route.',
      },
      {
        q: 'Do you offer gift wrapping?',
        a: 'On any order, at checkout. It adds a hand-written card and recycled paper, and it does not add a day to dispatch.',
      },
    ],
  },
  {
    id: 'shipping',
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
        q: 'What are your shipping rates?',
        a: 'Flat rate domestically and free above a threshold shown in the basket. International is calculated by weight and destination.',
      },
    ],
  },
  {
    id: 'payment',
    title: 'Payment',
    items: [
      {
        q: 'What currencies do you accept?',
        a: 'Around thirty, chosen from your location and changeable in the footer. You are charged in the currency you see.',
      },
      {
        q: 'Is my payment information secure?',
        a: 'Card details go straight to our PCI-compliant provider and never reach our servers, so there is nothing on our side to leak.',
      },
      {
        q: 'Can I get a refund if I change my mind?',
        a: 'Within the thirty-day window, back to the original payment method. It usually clears in three to five working days.',
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

export default function FaqsCategoryNav({
  title = 'FAQs',
  subtitle = 'Quick answers to what people ask most about the platform, the service and how it all fits together.',
  groups = GROUPS,
}: {
  title?: string
  subtitle?: string
  groups?: Group[]
}) {
  return (
    <section className="scroll-smooth bg-white py-24 motion-reduce:scroll-auto sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="max-w-2xl">
          <h2 className="text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
            {title}
          </h2>
          <p className="mt-6 text-lg/8 text-gray-600 dark:text-gray-400">{subtitle}</p>
        </div>

        <div className="mt-16 grid grid-cols-1 gap-x-16 gap-y-10 lg:grid-cols-[minmax(0,14rem)_minmax(0,1fr)]">
          <nav aria-label="FAQ categories" className="lg:sticky lg:top-24 lg:self-start">
            <ul role="list" className="flex flex-wrap gap-2 lg:flex-col lg:gap-1">
              {groups.map((group) => (
                <li key={group.id}>
                  <a
                    href={`#${group.id}`}
                    className="flex min-h-11 items-center rounded-lg px-3 text-sm font-medium text-gray-600 hover:bg-gray-100 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-gray-400 dark:hover:bg-white/5 dark:hover:text-white"
                  >
                    {group.title}
                  </a>
                </li>
              ))}
            </ul>
          </nav>

          <div className="space-y-14">
            {groups.map((group) => (
              <div key={group.id}>
                {/* scroll-mt clears the sticky rail when jumped to. */}
                <h3
                  id={group.id}
                  className="scroll-mt-24 text-lg font-semibold text-gray-900 dark:text-white"
                >
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
        </div>
      </div>
    </section>
  )
}
