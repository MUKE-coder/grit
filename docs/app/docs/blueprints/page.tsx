import Link from 'next/link'
import { ArrowRight, Boxes, Github, Lock } from 'lucide-react'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'

export const metadata = {
  title: 'Blueprints: complete applications built with Grit — Grit',
  description:
    'A blueprint is a whole application in its own repository, built with Grit. WhatsApp, Instagram and a storefront, with what each one teaches and how it is kept current.',
  alternates: { canonical: 'https://gritframework.dev/docs/blueprints' },
}

// Mirrors internal/blueprint/registry.go, which `grit blueprint list` reads.
// TestBlueprintsPageListsEveryBlueprint compares the two.
const BLUEPRINTS = [
  {
    name: 'whatsapp',
    title: 'WhatsApp',
    summary: 'Realtime chat across web, mobile and desktop, on one Go API',
    access: 'open' as const,
    surfaces: 'web, admin, mobile, desktop',
    repo: 'https://github.com/MUKE-coder/grit-whatsapp',
    teaches: [
      'One conversation model reaching four clients without four implementations',
      'Message delivery and read state over the domain event bus',
      'Attachments through the same storage the rest of the app uses',
      'An Expo app and a Wails desktop app against the same shared types',
    ],
  },
  {
    name: 'instagram',
    title: 'Instagram',
    summary: 'Photo and video feed with an Expo app',
    access: 'paid' as const,
    surfaces: 'web, admin, mobile',
    repo: 'https://github.com/MUKE-coder/grit-instagram',
    teaches: [
      'Image and video processing on upload, and what to do off the request',
      'A feed that paginates without drifting as new posts arrive',
      'Mobile capture and upload from Expo',
    ],
  },
  {
    name: 'storefront',
    title: 'Kaya Goods',
    summary: 'A shop: products, variants, cart and Stripe checkout',
    access: 'paid' as const,
    surfaces: 'web, admin',
    repo: 'https://github.com/MUKE-coder/grit-storefront',
    teaches: [
      'Product variants, and the option matrix that generates them',
      'A cart that survives sign-in, and a checkout that survives a retry',
      'Stripe webhooks that are idempotent when Stripe redelivers',
    ],
  },
  {
    name: 'stripe-dashboard',
    title: 'Stripe dashboard',
    summary: 'Payments, subscriptions and revenue reporting',
    access: 'planned' as const,
    surfaces: 'admin',
    repo: '',
    teaches: [
      'Reporting over money without recomputing it on every page view',
      'Subscription lifecycle: trials, dunning, cancellation',
    ],
  },
  {
    name: 'uber',
    title: 'Uber',
    summary: 'Live tracking on a map, with PostGIS',
    access: 'planned' as const,
    surfaces: 'web, admin, mobile',
    repo: '',
    teaches: [
      'Geospatial queries in Postgres, and what they cost',
      'A position stream that does not melt the database',
    ],
  },
  {
    name: 'linear',
    title: 'Linear',
    summary: 'An issue tracker built for the keyboard',
    access: 'planned' as const,
    surfaces: 'web, admin',
    repo: '',
    teaches: [
      'Keyboard-first navigation across a whole application',
      'Optimistic updates that stay correct when the write is refused',
    ],
  },
]

const ACCESS_LABEL = {
  open: { text: 'Open source', className: 'bg-emerald-500/10 text-emerald-500' },
  paid: { text: 'Membership', className: 'bg-amber-500/10 text-amber-500' },
  planned: { text: 'Planned', className: 'bg-muted text-muted-foreground' },
}

export default function BlueprintsPage() {
  const available = BLUEPRINTS.filter((b) => b.access !== 'planned')
  const planned = BLUEPRINTS.filter((b) => b.access === 'planned')

  return (
    <DocsShell>
          <div className="mb-3 flex items-center gap-2">
            <Boxes className="h-5 w-5 text-accent" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Blueprints
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">Blueprints</h1>
          <p className="mb-8 text-lg leading-relaxed text-muted-foreground">
            A blueprint is a whole application, in its own repository, built with Grit. Not
            something you install: something you clone and work from. A{' '}
            <Link href="/docs/plugins/overview" className="text-accent hover:underline">
              plugin
            </Link>{' '}
            adds a feature to your project; a blueprint is a project.
          </p>

          <CodeBlock language="bash" code={`grit blueprint list\ngrit blueprint info whatsapp`} />

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Why they exist</h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            Two reasons, and the second one is the one that matters to you even if you never open
            one.
          </p>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            The first is a starting point. If you are building a storefront, starting from a
            storefront beats starting from an empty project and a tutorial.
          </p>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            The second is that each blueprint is a load test for the framework. Building a real
            application finds things a test suite does not, and every gap a blueprint finds is fixed
            in Grit first, released, and then the blueprint runs <code>grit upgrade</code> into it.
            Its code stays Grit&apos;s own rather than drifting into a fork. The storefront alone
            found five framework bugs in one sitting, and those fixes shipped to everybody.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Available now</h2>
          <div className="mb-10 space-y-4">
            {available.map((b) => {
              const label = ACCESS_LABEL[b.access]
              return (
                <div key={b.name} className="rounded-xl border border-border p-5">
                  <div className="mb-2 flex flex-wrap items-center gap-3">
                    <h3 className="text-lg font-semibold tracking-tight">{b.title}</h3>
                    <span className={`rounded-full px-2.5 py-0.5 text-xs font-medium ${label.className}`}>
                      {label.text}
                    </span>
                    <code className="ml-auto font-mono text-xs text-muted-foreground">{b.name}</code>
                  </div>
                  <p className="mb-3 text-sm text-muted-foreground">{b.summary}</p>
                  <p className="mb-3 font-mono text-xs text-muted-foreground">{b.surfaces}</p>
                  <p className="mb-2 text-xs font-medium uppercase tracking-wider text-muted-foreground">
                    What you open it for
                  </p>
                  <ul className="mb-4 ml-4 list-disc space-y-1 text-sm text-muted-foreground">
                    {b.teaches.map((t) => (
                      <li key={t}>{t}</li>
                    ))}
                  </ul>
                  {b.access === 'open' ? (
                    <Link
                      href={b.repo}
                      className="inline-flex items-center gap-1.5 text-sm text-accent hover:underline"
                    >
                      <Github className="h-3.5 w-3.5" />
                      Clone it
                    </Link>
                  ) : (
                    <span className="inline-flex items-center gap-1.5 text-sm text-muted-foreground">
                      <Lock className="h-3.5 w-3.5" />
                      A membership sends an invitation to the private repository
                    </span>
                  )}
                </div>
              )
            })}
          </div>

          <h2 className="mb-2 mt-12 text-2xl font-semibold tracking-tight">On the list</h2>
          <p className="mb-5 leading-relaxed text-muted-foreground">
            Announced, not built yet. Listed because &quot;is there a Linear blueprint?&quot; is
            better answered with &quot;not yet&quot; than with silence.
          </p>
          <div className="mb-10 grid gap-3 sm:grid-cols-3">
            {planned.map((b) => (
              <div key={b.name} className="rounded-xl border border-dashed border-border p-4">
                <p className="font-medium">{b.title}</p>
                <p className="mt-1 text-sm text-muted-foreground">{b.summary}</p>
              </div>
            ))}
          </div>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">
            Blueprints and examples are different things
          </h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            The Grit repository also carries <code>examples/</code>: complete applications checked
            into the repo and built on every CI run. Those are free and they are there to prove the
            framework works and to be read.
          </p>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            A blueprint is a bigger, designed product in a repository of its own, kept current
            against new releases. Both are real applications; the difference is scope and where they
            live.
          </p>

          <div className="mt-12 flex flex-col gap-3 sm:flex-row">
            <Link
              href="/docs/extending"
              className="group flex flex-1 items-center gap-3 rounded-xl border border-border p-4 transition-colors hover:border-accent"
            >
              <span className="flex-1">
                <span className="block font-medium">Blueprint, plugin, resource or block?</span>
                <span className="text-sm text-muted-foreground">Which one a job needs</span>
              </span>
              <ArrowRight className="h-4 w-4 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
            </Link>
            <Link
              href="https://gritframework.dev/blueprints"
              className="group flex flex-1 items-center gap-3 rounded-xl border border-border p-4 transition-colors hover:border-accent"
            >
              <span className="flex-1">
                <span className="block font-medium">Get access</span>
                <span className="text-sm text-muted-foreground">Memberships and repositories</span>
              </span>
              <ArrowRight className="h-4 w-4 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
            </Link>
          </div>
    </DocsShell>
  )
}
