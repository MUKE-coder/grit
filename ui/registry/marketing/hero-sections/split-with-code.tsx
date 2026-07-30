import { ArrowRight, Check } from 'lucide-react'

const TERMINAL_LINES = [
  { prompt: true, text: 'grit new storefront --triple' },
  { text: '✓ Go API, admin panel and web app scaffolded' },
  { prompt: true, text: 'grit generate resource Product \\' },
  { text: '    --fields "name:string,price:float,stock:int"' },
  { text: '✓ model · service · handler · routes' },
  { text: '✓ Zod schema · TS types · React Query hooks' },
  { text: '✓ admin page' },
]

export default function SplitWithCode({
  eyebrow = 'Full-stack, one command',
  title = 'From idea to a running app in five minutes',
  subtitle = 'Grit generates the Go backend and the React frontend from a single resource definition, so the two halves cannot drift apart.',
  bullets = [
    'Typed end to end, database to UI',
    'Admin panel generated with every resource',
    'Deploys as one container',
  ],
  primaryLabel = 'Start building',
  secondaryLabel = 'Live demo',
}: {
  eyebrow?: string
  title?: string
  subtitle?: string
  bullets?: string[]
  primaryLabel?: string
  secondaryLabel?: string
}) {
  return (
    <div className="overflow-hidden bg-white py-24 sm:py-32 dark:bg-gray-900">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="mx-auto grid max-w-2xl grid-cols-1 gap-x-12 gap-y-16 lg:mx-0 lg:max-w-none lg:grid-cols-2 lg:items-center">
          {/* Copy */}
          <div className="lg:pr-4">
            <p className="text-base/7 font-semibold text-indigo-600 dark:text-indigo-400">
              {eyebrow}
            </p>
            <h1 className="mt-2 text-4xl font-semibold tracking-tight text-pretty text-gray-900 sm:text-5xl dark:text-white">
              {title}
            </h1>
            <p className="mt-6 text-lg/8 text-gray-600 dark:text-gray-400">{subtitle}</p>

            <ul role="list" className="mt-8 space-y-3">
              {bullets.map((bullet) => (
                <li key={bullet} className="flex gap-x-3">
                  <Check
                    aria-hidden="true"
                    className="mt-1 size-5 flex-none text-indigo-600 dark:text-indigo-400"
                  />
                  <span className="text-base/7 text-gray-600 dark:text-gray-400">{bullet}</span>
                </li>
              ))}
            </ul>

            <div className="mt-10 flex items-center gap-x-6">
              <a
                href="#"
                className="rounded-md bg-indigo-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
              >
                {primaryLabel}
              </a>
              <a
                href="#"
                className="group inline-flex items-center gap-1 text-sm/6 font-semibold text-gray-900 dark:text-white"
              >
                {secondaryLabel}
                <ArrowRight
                  aria-hidden="true"
                  className="size-4 transition-transform group-hover:translate-x-0.5"
                />
              </a>
            </div>
          </div>

          {/* Terminal */}
          <div className="relative">
            <div
              aria-hidden="true"
              className="absolute -inset-x-4 -inset-y-8 -z-10 rounded-3xl bg-gradient-to-tr from-indigo-500/10 to-sky-400/10 blur-2xl"
            />
            <div className="overflow-hidden rounded-xl bg-gray-900 shadow-2xl ring-1 ring-white/10 dark:bg-gray-950">
              <div className="flex items-center gap-2 border-b border-white/10 px-4 py-3">
                <span className="size-3 rounded-full bg-[#ff5f57]" />
                <span className="size-3 rounded-full bg-[#febc2e]" />
                <span className="size-3 rounded-full bg-[#28c840]" />
                <span className="ml-2 font-mono text-xs text-gray-400">bash</span>
              </div>
              <pre className="overflow-x-auto p-5 font-mono text-[13px]/6">
                <code>
                  {TERMINAL_LINES.map((line, i) => (
                    <div key={i} className="whitespace-pre">
                      {line.prompt ? (
                        <>
                          <span className="text-indigo-400 select-none">$ </span>
                          <span className="text-gray-100">{line.text}</span>
                        </>
                      ) : (
                        <span className="text-gray-400">{line.text}</span>
                      )}
                    </div>
                  ))}
                </code>
              </pre>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
