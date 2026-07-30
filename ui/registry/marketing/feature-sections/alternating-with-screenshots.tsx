import { ArrowRight, Check } from 'lucide-react'

/** The Grit UI mark, inlined so the block stays self-contained. */
function GritMark({ className = 'size-5' }: { className?: string }) {
  return (
    <svg viewBox="0 0 32 32" fill="none" aria-hidden="true" className={className}>
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M8 0H24A8 8 0 0 1 32 8V24A8 8 0 0 1 24 32H8A8 8 0 0 1 0 24V8A8 8 0 0 1 8 0ZM16 9.4A6.6 6.6 0 1 0 21.4 19.9V17.4H17.2A1.7 1.7 0 0 1 17.2 14H23.1A1.7 1.7 0 0 1 24.8 15.7V20.6A1.7 1.7 0 0 1 24.4 21.7A10 10 0 1 1 22.6 8.2A1.7 1.7 0 0 1 20.4 10.8A6.6 6.6 0 0 0 16 9.4Z"
        fill="currentColor"
      />
    </svg>
  )
}

const ROWS = [
  {
    eyebrow: 'Generate',
    title: 'A resource is one command, not one afternoon',
    body: 'Define the fields once. Grit writes the model, the service, the handler, the routes, the validation schema, the typed hooks and the admin screen — consistent every time.',
    bullets: ['Model, service and handler', 'Typed client and hooks', 'Admin table and form'],
    mock: 'terminal' as const,
  },
  {
    eyebrow: 'Administer',
    title: 'An admin panel you never have to build',
    body: 'Every generated resource arrives with a filterable table, a validated form and a detail view. It is ordinary code in your repository, so change anything you like.',
    bullets: ['Sort, filter and paginate', 'Validation from your Go tags', 'Role-aware actions'],
    mock: 'table' as const,
  },
]

export default function AlternatingWithScreenshots({
  eyebrow = 'How it works',
  title = 'From definition to running feature',
}: {
  eyebrow?: string
  title?: string
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="mx-auto max-w-2xl text-center">
          <p className="text-base/7 font-semibold text-indigo-600 dark:text-indigo-400">
            {eyebrow}
          </p>
          <h2 className="mt-2 text-4xl font-semibold tracking-tight text-pretty text-gray-900 sm:text-5xl dark:text-white">
            {title}
          </h2>
        </div>

        <div className="mt-20 space-y-24 lg:space-y-32">
          {ROWS.map((row, i) => (
            <div
              key={row.title}
              className="grid grid-cols-1 items-center gap-12 lg:grid-cols-2 lg:gap-16"
            >
              {/* Odd rows put the copy on the right. order-first on large
                  screens only, so the stacked mobile order stays copy-then-art
                  and never leads with a screenshot that has no context yet. */}
              <div className={i % 2 === 1 ? 'lg:order-last' : undefined}>
                <p className="text-base/7 font-semibold text-indigo-600 dark:text-indigo-400">
                  {row.eyebrow}
                </p>
                <h3 className="mt-2 text-3xl font-semibold tracking-tight text-pretty text-gray-900 dark:text-white">
                  {row.title}
                </h3>
                <p className="mt-5 text-lg/8 text-gray-600 dark:text-gray-400">{row.body}</p>

                <ul role="list" className="mt-8 space-y-3">
                  {row.bullets.map((bullet) => (
                    <li key={bullet} className="flex gap-x-3">
                      <Check
                        aria-hidden="true"
                        className="mt-1 size-5 flex-none text-indigo-600 dark:text-indigo-400"
                      />
                      <span className="text-base/7 text-gray-600 dark:text-gray-400">
                        {bullet}
                      </span>
                    </li>
                  ))}
                </ul>

                <a
                  href="#"
                  className="group mt-8 inline-flex items-center gap-1.5 text-sm font-semibold text-indigo-600 dark:text-indigo-400"
                >
                  Read the guide
                  <ArrowRight
                    aria-hidden="true"
                    className="size-4 transition-transform group-hover:translate-x-0.5"
                  />
                </a>
              </div>

              <div>{row.mock === 'terminal' ? <TerminalMock /> : <TableMock />}</div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}

function TerminalMock() {
  const lines = [
    { prompt: true, text: 'grit generate resource Invoice \\' },
    { text: '    --fields "number:string,total:float,paid:bool"' },
    { text: '✓ internal/models/invoice.go' },
    { text: '✓ internal/services/invoice.go' },
    { text: '✓ internal/handlers/invoice.go' },
    { text: '✓ packages/shared/schemas/invoice.ts' },
    { text: '✓ apps/admin/app/resources/invoices/page.tsx' },
    { ok: true, text: '7 files written in 240ms' },
  ]

  return (
    <div className="overflow-hidden rounded-xl bg-gray-900 shadow-2xl ring-1 ring-white/10 dark:bg-gray-900">
      <div className="flex items-center gap-2 border-b border-white/10 px-4 py-3">
        <span className="size-3 rounded-full bg-[#ff5f57]" />
        <span className="size-3 rounded-full bg-[#febc2e]" />
        <span className="size-3 rounded-full bg-[#28c840]" />
        <span className="ml-2 font-mono text-xs text-gray-400">bash</span>
      </div>
      <pre className="overflow-x-auto p-5 font-mono text-[13px]/6">
        <code>
          {lines.map((line, i) => (
            <div key={i} className="whitespace-pre">
              {line.prompt ? (
                <>
                  <span className="text-indigo-400 select-none">$ </span>
                  <span className="text-gray-100">{line.text}</span>
                </>
              ) : (
                <span className={line.ok ? 'text-green-400' : 'text-gray-400'}>
                  {line.text}
                </span>
              )}
            </div>
          ))}
        </code>
      </pre>
    </div>
  )
}

function TableMock() {
  const rows = [
    ['INV-1042', 'Northwind', '$4,200.00', 'Paid'],
    ['INV-1041', 'Acme', '$1,850.00', 'Pending'],
    ['INV-1040', 'Globex', '$12,400.00', 'Paid'],
    ['INV-1039', 'Initech', '$640.00', 'Overdue'],
  ]
  const tone: Record<string, string> = {
    Paid: 'bg-green-100 text-green-700 dark:bg-green-500/15 dark:text-green-400',
    Pending: 'bg-amber-100 text-amber-700 dark:bg-amber-500/15 dark:text-amber-400',
    Overdue: 'bg-rose-100 text-rose-700 dark:bg-rose-500/15 dark:text-rose-400',
  }

  return (
    <div className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-xl dark:border-white/10 dark:bg-gray-900">
      <div className="flex items-center gap-2 border-b border-gray-200 px-4 py-3 dark:border-white/10">
        <GritMark className="size-5 text-gray-900 dark:text-white" />
        <span className="text-sm font-medium text-gray-900 dark:text-white">Invoices</span>
        <span className="ml-auto rounded-md bg-indigo-600 px-2.5 py-1 text-xs font-medium text-white">
          New invoice
        </span>
      </div>

      <table className="w-full text-left text-sm">
        <thead className="bg-gray-50 text-xs text-gray-500 dark:bg-white/5 dark:text-gray-400">
          <tr>
            {['Number', 'Client', 'Total', 'Status'].map((h) => (
              <th key={h} scope="col" className="px-4 py-2.5 font-medium">
                {h}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map(([number, client, total, status]) => (
            <tr key={number} className="border-t border-gray-100 dark:border-white/5">
              <td className="px-4 py-3 font-mono text-xs text-gray-500 dark:text-gray-400">
                {number}
              </td>
              <td className="px-4 py-3 text-gray-900 dark:text-white">{client}</td>
              <td className="px-4 py-3 text-gray-900 tabular-nums dark:text-white">
                {total}
              </td>
              <td className="px-4 py-3">
                <span
                  className={`rounded-full px-2 py-0.5 text-xs font-medium ${tone[status]}`}
                >
                  {status}
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
