/*
 * One number, as large as it will go.
 *
 * Only worth using when the number is genuinely the headline. A figure this
 * size makes an implicit promise that it is remarkable, and a reader who does
 * the arithmetic and finds it is not remembers that instead.
 *
 * Two details make the digits behave. `tabular-nums` gives every digit the
 * same width, so a counter that ticks upward does not shuffle the ones beside
 * it sideways. And the thousands separators are real characters in the string
 * rather than spacing tricks, so the number can be selected, copied and read
 * aloud as a number.
 *
 * The figure is a <dd> and the caption a <dt>, so a screen reader hears them
 * as one pair rather than an orphan number followed by an unrelated sentence.
 * `text-[clamp(...)]` scales it against the viewport with a floor and a
 * ceiling, which is the one place a fluid size beats a stack of breakpoints:
 * the whole design here is "as large as it fits".
 */

export default function StatsSingleGiantNumber({
  value = '67,904,370',
  label = 'Requests served last month',
  body = 'Across every region, every plan and every customer running on the platform.',
}: {
  value?: string
  label?: string
  body?: string
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 text-center lg:px-8">
        <dl>
          <dd
            className="font-semibold tracking-tighter text-gray-900 tabular-nums dark:text-white"
            style={{ fontSize: 'clamp(2.75rem, 11vw, 9rem)', lineHeight: 1 }}
          >
            {value}
          </dd>
          <dt className="sr-only">{label}</dt>
        </dl>
        <p className="mx-auto mt-8 max-w-xl text-base/7 text-pretty text-gray-600 dark:text-gray-400">
          <span className="font-medium text-gray-900 dark:text-white">{label}.</span> {body}
        </p>
      </div>
    </section>
  )
}
