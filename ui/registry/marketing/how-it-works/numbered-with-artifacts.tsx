/*
 * Three numbered steps, each above a piece of the thing you actually do.
 *
 * The number lives in the heading text rather than in a badge beside it, so a
 * screen reader reads "1. Add payment information" as one heading instead of
 * announcing a decorative circle and then the title. The list is an <ol>: the
 * order is the point, and an ordered list is how you say that in markup rather
 * than in CSS.
 *
 * Each artifact below the copy is markup, marked `aria-hidden`, and faded out
 * at the bottom with a mask rather than a gradient overlay. A gradient overlay
 * is an opaque box that has to match whatever is behind it, so it breaks the
 * moment the section sits on a different background; a mask removes the pixels
 * and works on any backdrop, in either theme.
 */

const CARD =
  'relative overflow-hidden rounded-xl border border-gray-200 bg-white p-5 dark:border-white/10 dark:bg-gray-900'

/* Fades the artifact out instead of cutting it off. */
const FADE = {
  maskImage: 'linear-gradient(to bottom, black 55%, transparent 100%)',
  WebkitMaskImage: 'linear-gradient(to bottom, black 55%, transparent 100%)',
}

function Field({ label, value, muted }: { label: string; value: string; muted?: boolean }) {
  return (
    <div>
      <p className="text-xs font-medium text-gray-700 dark:text-gray-300">{label}</p>
      <div className="mt-1 rounded-md border border-gray-200 px-3 py-2 dark:border-white/10">
        <span className={`text-sm ${muted ? 'text-gray-400 dark:text-gray-600' : 'text-gray-700 dark:text-gray-300'}`}>
          {value}
        </span>
      </div>
    </div>
  )
}

function PaymentArtifact() {
  return (
    <div className={CARD} style={FADE}>
      <div className="space-y-4">
        <Field label="Email" value="irung@example.com" />
        <div>
          <p className="text-xs font-medium text-gray-700 dark:text-gray-300">Card information</p>
          <div className="mt-1 rounded-md border border-gray-200 dark:border-white/10">
            <div className="border-b border-gray-200 px-3 py-2 dark:border-white/10">
              <span className="text-sm text-gray-400 dark:text-gray-600">1234 1234 1234 1234</span>
            </div>
            <div className="grid grid-cols-2 divide-x divide-gray-200 dark:divide-white/10">
              <span className="px-3 py-2 text-sm text-gray-400 dark:text-gray-600">MM/YY</span>
              <span className="px-3 py-2 text-sm text-gray-400 dark:text-gray-600">CVV</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

function InvoiceArtifact() {
  return (
    <div className={CARD} style={FADE}>
      <p className="font-mono text-xs text-gray-500 dark:text-gray-400">INV-456789</p>
      <p className="mt-1 font-mono text-2xl font-semibold text-gray-900 tabular-nums dark:text-white">
        $284,342.57
      </p>
      <p className="mt-1 text-xs font-medium text-gray-600 dark:text-gray-400">Due in 15 days</p>
      <div className="mt-5 flex h-20 items-center justify-center rounded-md border border-dashed border-gray-300 dark:border-white/15">
        <span className="text-sm text-gray-400 italic dark:text-gray-600">Sign here</span>
      </div>
    </div>
  )
}

function ReceiptArtifact() {
  return (
    <div className={CARD} style={FADE}>
      <p className="font-mono text-xs text-gray-500 dark:text-gray-400">INV-456789</p>
      <p className="mt-1 font-mono text-2xl font-semibold text-gray-900 tabular-nums dark:text-white">
        $284,342.57
      </p>
      <p className="mt-1 text-xs font-medium text-gray-600 dark:text-gray-400">Paid in full</p>
      <dl className="mt-5 space-y-3">
        {['To', 'From', 'Address'].map((label, i) => (
          <div key={label} className="flex items-center gap-4">
            <dt className="w-16 flex-none text-sm text-gray-500 dark:text-gray-500">{label}</dt>
            <dd
              className="h-2 rounded-full bg-gray-200 dark:bg-white/10"
              style={{ width: `${60 - i * 14}%` }}
            />
          </div>
        ))}
      </dl>
    </div>
  )
}

const STEPS = [
  {
    title: 'Add payment information',
    body: 'Securely add your payment details to get started with our services.',
    Artifact: PaymentArtifact,
  },
  {
    title: 'Sign documents',
    body: 'Digitally sign and authorise transactions, with a full audit trail behind every signature.',
    Artifact: InvoiceArtifact,
  },
  {
    title: 'Receive confirmation',
    body: 'Get instant confirmation receipts for every completed transaction.',
    Artifact: ReceiptArtifact,
  },
]

export default function HowItWorksNumberedWithArtifacts({
  steps = STEPS,
}: {
  steps?: { title: string; body: string; Artifact: () => React.JSX.Element }[]
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <ol role="list" className="grid grid-cols-1 gap-x-8 gap-y-16 lg:grid-cols-3">
          {steps.map((step, i) => (
            <li key={step.title}>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                <span className="text-gray-400 tabular-nums dark:text-gray-500">{i + 1}.</span>{' '}
                {step.title}
              </h3>
              <p className="mt-2 text-sm/6 text-gray-600 dark:text-gray-400">{step.body}</p>
              <div aria-hidden="true" className="mt-8 select-none">
                <step.Artifact />
              </div>
            </li>
          ))}
        </ol>
      </div>
    </section>
  )
}
