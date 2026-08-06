/*
 * The same three steps, run down the page instead of across it.
 *
 * Worth having as its own block rather than a breakpoint of the horizontal
 * one: a vertical sequence reads as a process you go through, a horizontal one
 * reads as three things that exist. Use this where the order genuinely matters
 * and the steps take time.
 *
 * The connecting arrows point down at every width, because the layout is
 * vertical at every width. They are `aria-hidden` — the <ol> already says this
 * is a sequence, and "1 of 3" is announced without them.
 */

function SpreadsheetGlyph() {
  return (
    <div className="relative">
      <div className="grid w-24 grid-cols-3 gap-1 rounded-lg border border-gray-200 bg-white p-2 shadow-sm dark:border-white/10 dark:bg-gray-900">
        {Array.from({ length: 18 }).map((_, i) => (
          <span key={i} className="h-2 rounded-xs bg-gray-200 dark:bg-white/10" />
        ))}
      </div>
      <span className="absolute -right-2 -bottom-2 rounded-md bg-emerald-500 px-1.5 py-0.5 text-[10px] font-semibold text-white">
        CSV
      </span>
    </div>
  )
}

function CurrencyGlyph() {
  const cards = [
    { label: '₿ BTC', tone: 'bg-blue-50 text-blue-700 dark:bg-blue-500/10 dark:text-blue-300', rotate: '-rotate-6' },
    { label: '$ USD', tone: 'bg-emerald-50 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-300', rotate: '' },
    { label: '€ EUR', tone: 'bg-rose-50 text-rose-700 dark:bg-rose-500/10 dark:text-rose-300', rotate: 'rotate-6' },
  ]
  return (
    <div className="flex items-end -space-x-4">
      {cards.map((card) => (
        <div
          key={card.label}
          className={`${card.rotate} w-16 rounded-lg border border-gray-200 p-2 shadow-sm dark:border-white/10 ${card.tone}`}
        >
          <p className="font-mono text-[10px] font-semibold">{card.label}</p>
          <div className="mt-2 space-y-1">
            <span className="block h-1 w-full rounded-full bg-current opacity-25" />
            <span className="block h-1 w-2/3 rounded-full bg-current opacity-25" />
          </div>
        </div>
      ))}
    </div>
  )
}

function ReportsGlyph() {
  return (
    <div className="flex gap-2">
      {[0, 1].map((i) => (
        <div
          key={i}
          className="w-16 rounded-lg border border-gray-200 bg-white p-2 shadow-sm dark:border-white/10 dark:bg-gray-900"
        >
          <span className="block size-2 rounded-full bg-gray-300 dark:bg-white/20" />
          <div className="mt-2 space-y-1">
            <span className="block h-1 w-full rounded-full bg-gray-200 dark:bg-white/10" />
            <span className="block h-1 w-full rounded-full bg-gray-200 dark:bg-white/10" />
            <span className="block h-1 w-2/3 rounded-full bg-gray-200 dark:bg-white/10" />
          </div>
        </div>
      ))}
    </div>
  )
}

const STEPS = [
  {
    title: 'Data collection',
    body: 'Import data from any source or format with the integration tools.',
    Glyph: SpreadsheetGlyph,
  },
  {
    title: 'Automated analysis',
    body: 'The pipeline processes complex datasets and surfaces the patterns in them instantly.',
    Glyph: CurrencyGlyph,
  },
  {
    title: 'Actionable reports',
    body: 'Turn the result into visualisations and shareable reports that drive decisions.',
    Glyph: ReportsGlyph,
  },
]

export default function HowItWorksVerticalCenteredSteps({
  eyebrow = 'Our process',
  title = 'Simple three-step workflow',
  subtitle = 'A streamlined approach to data analysis that lets your team decide quickly and with confidence.',
  ctaLabel = 'Get started',
  ctaHref = '#',
  steps = STEPS,
}: {
  eyebrow?: string
  title?: string
  subtitle?: string
  ctaLabel?: string
  ctaHref?: string
  steps?: { title: string; body: string; Glyph: () => React.JSX.Element }[]
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-2xl px-6 text-center lg:px-8">
        <p className="text-base font-semibold text-indigo-600 dark:text-indigo-400">{eyebrow}</p>
        <h2 className="mt-2 text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
          {title}
        </h2>
        <p className="mt-6 text-lg/8 text-pretty text-gray-600 dark:text-gray-400">{subtitle}</p>

        <ol role="list" className="mt-20">
          {steps.map((step, i) => (
            <li key={step.title}>
              <span className="inline-flex size-8 items-center justify-center rounded-full bg-gray-100 text-sm font-semibold text-gray-700 tabular-nums dark:bg-white/10 dark:text-gray-300">
                {i + 1}
              </span>
              <div aria-hidden="true" className="mt-8 flex h-28 items-center justify-center">
                <step.Glyph />
              </div>
              <h3 className="mt-6 text-xl font-semibold text-gray-900 dark:text-white">
                {step.title}
              </h3>
              <p className="mx-auto mt-2 max-w-sm text-sm/6 text-gray-600 dark:text-gray-400">
                {step.body}
              </p>

              {i < steps.length - 1 && (
                <svg
                  aria-hidden="true"
                  viewBox="0 0 24 24"
                  fill="none"
                  className="mx-auto my-10 size-6 text-gray-300 dark:text-gray-700"
                >
                  <path
                    d="M12 4v15m0 0 5-5m-5 5-5-5"
                    stroke="currentColor"
                    strokeWidth="1.5"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                </svg>
              )}
            </li>
          ))}
        </ol>

        <div className="mt-16">
          <a
            href={ctaHref}
            className="inline-flex min-h-11 items-center rounded-lg border border-gray-300 px-6 text-sm font-semibold text-gray-900 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:border-white/15 dark:text-white dark:hover:bg-white/5"
          >
            {ctaLabel}
          </a>
        </div>
      </div>
    </section>
  )
}
