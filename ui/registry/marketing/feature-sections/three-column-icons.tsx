import {
  Boxes,
  GitBranch,
  Layers,
  Lock,
  Terminal,
  Zap,
} from 'lucide-react'

const FEATURES = [
  {
    icon: Terminal,
    title: 'One command',
    body: 'Scaffold the API, the admin panel and the frontend together — wired, typed and running before you have finished your coffee.',
  },
  {
    icon: Lock,
    title: 'Hardened by default',
    body: 'CSRF, a strict content policy, rate limiting and a tamper-evident audit log ship in the scaffold, not in a checklist for later.',
  },
  {
    icon: Layers,
    title: 'Five architectures',
    body: 'Embed a SPA in the Go binary, split into a monorepo, go API-only, or add mobile and desktop. Same generators throughout.',
  },
  {
    icon: Zap,
    title: 'Typed end to end',
    body: 'Go models generate TypeScript types and validation schemas, so the two halves of your app cannot drift apart.',
  },
  {
    icon: GitBranch,
    title: 'Code you own',
    body: 'Everything lands in your repository as ordinary source. No runtime magic to fight, nothing hidden behind a dependency.',
  },
  {
    icon: Boxes,
    title: 'Batteries included',
    body: 'Auth, storage, email, queues and cron are already there — each a separate package you can switch off or delete.',
  },
]

export default function ThreeColumnIcons({
  eyebrow = 'Everything included',
  title = 'The boring parts, already done',
  subtitle = 'The work that produces no product differentiation, finished before you start — so week one goes into the thing people actually pay for.',
}: {
  eyebrow?: string
  title?: string
  subtitle?: string
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="mx-auto max-w-2xl lg:text-center">
          <p className="text-base/7 font-semibold text-indigo-600 dark:text-indigo-400">
            {eyebrow}
          </p>
          <h2 className="mt-2 text-4xl font-semibold tracking-tight text-pretty text-gray-900 sm:text-5xl dark:text-white">
            {title}
          </h2>
          <p className="mt-6 text-lg/8 text-gray-600 dark:text-gray-400">{subtitle}</p>
        </div>

        <div className="mx-auto mt-16 max-w-2xl sm:mt-20 lg:mt-24 lg:max-w-none">
          <dl className="grid max-w-xl grid-cols-1 gap-x-8 gap-y-12 lg:max-w-none lg:grid-cols-3">
            {FEATURES.map(({ icon: Icon, title: name, body }) => (
              <div key={name} className="flex flex-col">
                <dt className="flex items-center gap-x-3 text-base/7 font-semibold text-gray-900 dark:text-white">
                  <span className="flex size-10 flex-none items-center justify-center rounded-lg bg-indigo-600">
                    <Icon aria-hidden="true" className="size-5 text-white" />
                  </span>
                  {name}
                </dt>
                <dd className="mt-4 flex flex-auto flex-col text-base/7 text-gray-600 dark:text-gray-400">
                  <p className="flex-auto">{body}</p>
                </dd>
              </div>
            ))}
          </dl>
        </div>
      </div>
    </section>
  )
}
