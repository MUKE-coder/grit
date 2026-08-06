/*
 * FAQs: a sticky title column beside grouped questions.
 *
 * Built on <details> and <summary> rather than React state, and that is the
 * interesting decision here.
 *
 * A hand-rolled accordion has to re-implement things the browser already does:
 * focus, Enter and Space, the expanded state announced to a screen reader, and
 * find-in-page. Chrome and Firefox will open a closed <details> when the text
 * inside it matches a Ctrl+F search. A div with onClick will not, so the answer
 * someone searched for stays invisible and they conclude you do not cover it.
 *
 * It also works before hydration, which matters most for exactly this content:
 * the FAQ is what someone reads while deciding whether to trust you.
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
      {/* min-h-14 keeps the tap target above 44px without padding the text. */}
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

export default function FaqsSplitWithContact({
  title = 'FAQs',
  subtitle = 'Your questions answered',
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
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="grid grid-cols-1 gap-x-16 gap-y-12 lg:grid-cols-[minmax(0,20rem)_minmax(0,1fr)]">
          <div className="lg:sticky lg:top-24 lg:self-start">
            <h2 className="text-4xl font-semibold tracking-tight text-gray-900 sm:text-5xl dark:text-white">
              {title}
            </h2>
            <p className="mt-4 text-lg text-gray-600 dark:text-gray-400">{subtitle}</p>
            <p className="mt-8 text-base/7 text-gray-600 dark:text-gray-400">
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

          <div className="space-y-14">
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
        </div>
      </div>
    </section>
  )
}
