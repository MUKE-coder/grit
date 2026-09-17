import Link from 'next/link'
import { AlertCircle, ArrowLeft } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/deployment/checklist')

interface Item {
  label: string
  detail: string
}

const GROUPS: { title: string; note?: string; items: Item[] }[] = [
  {
    title: 'Secrets',
    note: 'The single most common way a hobby deploy becomes an incident.',
    items: [
      {
        label: 'JWT_SECRET is unique to this environment',
        detail:
          'Not the value from .env.example — that string is identical in every Grit project ever generated. Staging and production must differ, or a token minted in staging is valid in production.',
      },
      {
        label: 'FIELD_ENCRYPTION_KEY is in the secret store and backed up separately',
        detail:
          'It encrypts two-factor secrets and encrypted columns. Lose it and every two-factor user is locked out. See the section on backing it up below.',
      },
      {
        label: 'No secret is a build argument',
        detail:
          'Build logs are retained, often readable by more people than the container environment, and frequently shipped to a third-party CI provider.',
      },
      {
        label: 'Nothing sensitive is in a NEXT_PUBLIC_ variable',
        detail:
          'Those are compiled into JavaScript that every visitor downloads. Grep the built output before launch: a leaked key here is public the moment the page loads.',
      },
      {
        label: '.env is not committed',
        detail: 'Check the history, not just the working tree. A rotated key is only rotated if the old one is dead.',
      },
    ],
  },
  {
    title: 'Access',
    items: [
      {
        label: 'CORS_ORIGINS lists exact origins',
        detail: 'No wildcard. With credentialed requests a wildcard is both a security hole and rejected by browsers.',
      },
      {
        label: 'The database is not reachable from the internet',
        detail:
          'The development compose file publishes Postgres on a host port so you can attach a GUI. The production one must not. Confirm from outside the network, not from the server.',
      },
      {
        label: 'Studio, Pulse and Sentinel dashboards are protected',
        detail:
          'GORM Studio at /studio browses every table. Pulse and Sentinel expose traffic and configuration. All three must be behind auth, an IP allowlist, or switched off in production.',
      },
      {
        label: 'The default admin account is gone or has a new password',
        detail: 'admin@example.com / admin123 is in the seeder of every generated project, and in these docs.',
      },
    ],
  },
  {
    title: 'Data',
    items: [
      {
        label: 'Migrations ran as an explicit step',
        detail:
          'Grit does not auto-migrate in production. A process that rewrites the schema on every boot is a bad idea when the platform restarts it for its own reasons.',
      },
      {
        label: 'Backups are configured AND a restore has been tested',
        detail:
          'An untested backup is a hypothesis. Restore into a scratch database and count the rows — that is the only thing that turns it into a fact.',
      },
      {
        label: 'Uploads go to object storage, not local disk',
        detail:
          'Any platform that can move or rebuild your container will lose local files. This includes every PaaS on the previous pages.',
      },
    ],
  },
  {
    title: 'Operations',
    items: [
      {
        label: 'The health check points at /api/health',
        detail:
          'Not /. A platform checking the wrong path either cycles a healthy machine forever or reports a deploy that never finishes.',
      },
      {
        label: 'You know how to roll back',
        detail:
          'Before you need to, not during. On a plain VPS that means keeping the previous binary; on a PaaS it means having actually clicked the button once.',
      },
      {
        label: 'Logs go somewhere you can search',
        detail: 'Container stdout that nobody collects is not logging. Confirm you can find a specific request from an hour ago.',
      },
      {
        label: 'Someone is told when it breaks',
        detail:
          'An uptime check hitting /api/health from outside your infrastructure. Finding out from a customer is the expensive way.',
      },
    ],
  },
]

export default function ChecklistPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <Link
              href="/docs/deployment"
              className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors mb-8"
            >
              <ArrowLeft className="h-3.5 w-3.5" />
              Deployment
            </Link>

            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Deployment</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Go-live checklist</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                The things that are cheap to check now and expensive to discover later.
                Everything here has bitten a real deployment.
              </p>
            </div>

            <div className="rounded-xl border border-border/50 bg-card/50 p-5 mb-12">
              <p className="text-sm text-muted-foreground leading-relaxed">
                Most of this is machine-checkable. Run{' '}
                <code className="text-xs">grit doctor</code> first and let it find what it
                can — then work through the rest, which needs a human who knows what the app
                is for.
              </p>
              <div className="mt-4">
                <CodeBlock language="bash" code="grit doctor" />
              </div>
            </div>

            {GROUPS.map((group) => (
              <div key={group.title} className="mb-12">
                <h2 className="text-2xl font-bold tracking-tight mb-2">{group.title}</h2>
                {group.note && (
                  <p className="text-sm text-muted-foreground mb-5 leading-relaxed">{group.note}</p>
                )}
                <div className="space-y-3">
                  {group.items.map((item) => (
                    <div
                      key={item.label}
                      className="rounded-xl border border-border/50 bg-card/40 p-5"
                    >
                      <div className="flex gap-3">
                        {/* A real checkbox, not an icon. This page is meant to be
                            worked through, and an unticked box is a different
                            thing from a bullet point. */}
                        <span className="mt-0.5 flex h-4 w-4 shrink-0 items-center justify-center rounded border-2 border-border" />
                        <div>
                          <p className="font-medium text-sm leading-snug">{item.label}</p>
                          <p className="mt-1.5 text-[13px] text-muted-foreground leading-relaxed">
                            {item.detail}
                          </p>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            ))}

            <div id="field-encryption-key" className="mb-12 scroll-mt-24">
              <h2 className="text-2xl font-bold tracking-tight mb-4">Back up FIELD_ENCRYPTION_KEY</h2>
              <div className="space-y-4 text-muted-foreground leading-relaxed">
                <p>
                  <code className="text-xs">grit new</code> writes a random key to{' '}
                  <code className="text-xs">.env</code>. It encrypts every user&rsquo;s two-factor
                  secret and every <code className="text-xs">crypto.EncryptedString</code> column, so
                  the database on its own is not enough to read them. That is the point, and it
                  is also why the key needs as much care as the database.
                </p>
                <ul className="list-disc pl-5 space-y-2">
                  <li>
                    <strong className="text-foreground">Where it lives.</strong> Put it in the
                    deployment&rsquo;s secret store, not in the image, the repository or a build
                    argument. Every instance of the API needs the same value.
                  </li>
                  <li>
                    <strong className="text-foreground">Where the backup lives.</strong> Keep a copy
                    somewhere separate from the database backups, such as a password manager or a
                    secrets vault that a second person can reach. A backup stored next to the key
                    it protects is one leak away from being readable.
                  </li>
                  <li>
                    <strong className="text-foreground">If it is lost.</strong> Nobody with two-factor
                    sign-in can finish signing in, and each of them has to have two-factor reset by
                    an administrator. Encrypted columns cannot be read at all. Restoring a database
                    backup does not help, because the backup is encrypted with the same key.
                  </li>
                  <li>
                    <strong className="text-foreground">Rotating it.</strong> Changing the value is
                    not a rotation. Everything already encrypted stays encrypted under the old key,
                    so a new value has exactly the same effect as losing the old one. Treat the key
                    as permanent for the life of the data.
                  </li>
                </ul>
                <p>
                  A project created before keys were generated has none. Generate one, store it
                  as above, then run <code className="text-xs">grit migrate</code>, which encrypts the
                  two-factor secrets and other values that are already stored:
                </p>
              </div>
              <div className="mt-4">
                <CodeBlock language="bash" code={'openssl rand -base64 32\ngrit migrate'} />
              </div>
            </div>

            <div className="rounded-xl border border-amber-500/30 bg-amber-500/[0.05] p-6">
              <div className="flex gap-3">
                <AlertCircle className="h-5 w-5 shrink-0 text-amber-500 mt-0.5" />
                <div>
                  <h2 className="font-semibold mb-2">One thing this list cannot check</h2>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    Whether anyone other than you can deploy it. Write down the steps —
                    where it runs, how to get a shell, how to restore the database, who pays
                    the bill. A deployment only one person understands is an outage waiting
                    for that person to be unavailable.
                  </p>
                </div>
              </div>
            </div>

            <div className="border-t border-border/40 pt-8 mt-12">
              <Link
                href="/docs/deployment"
                className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
              >
                <ArrowLeft className="h-3.5 w-3.5" />
                All deployment targets
              </Link>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
