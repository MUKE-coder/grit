import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/security/doctor')

const checks: { name: string; level: string; what: string }[] = [
  {
    name: 'encryption-key-unset',
    level: 'error',
    what: 'A model has an encrypted field and FIELD_ENCRYPTION_KEY is not set, so those columns are stored in the clear. Nothing fails: the write succeeds and the value is readable in the database.',
  },
  {
    name: 'encrypted-column-in-whitelist',
    level: 'warning',
    what: 'An encrypted column can be searched, sorted or filtered. Ciphertext differs on every write, so the query matches nothing.',
  },
  {
    name: 'owned-resource-unscoped',
    level: 'error',
    what: 'A resource with an owner that nothing scopes, or one method that lost its scoping while the others kept it. The list is the path people forget, and an export is the list without pages.',
  },
  {
    name: 'owner-settable-from-body',
    level: 'error',
    what: 'The create or update request accepts the owner column, so a caller can file a row under somebody else’s account.',
  },
  {
    name: 'resource-could-be-owned',
    level: 'warning',
    what: 'A resource references a user and is scoped to nobody, so every signed-in caller sees every row. Sometimes that is right, which is why it is a question.',
  },
  {
    name: 'tenant-shared-resource',
    level: 'warning',
    what: 'The multitenant plugin is installed and a resource has no tenant.Owned, so its rows are shared across organizations.',
  },
  {
    name: 'pii-column-not-encrypted',
    level: 'warning',
    what: 'A column named like something that should not be readable in a database dump (an SSN, a card number, a diagnosis) is a plain column.',
  },
  {
    name: 'append-only-mutable-routes',
    level: 'error',
    what: 'An append-only resource still mounts PUT, PATCH or DELETE. The model and the database trigger refuse the write, so the endpoint can only fail.',
  },
  {
    name: 'studio-unprotected',
    level: 'error / warning',
    what: 'GORM Studio browses and edits every table. An error when it has no login at all, and when its password is still the default in production.',
  },
  {
    name: 'default-credentials',
    level: 'error / warning',
    what: 'Dashboard credentials still at their defaults, a JWT secret that is empty, a placeholder or under 32 characters, and an unkeyed security audit chain.',
  },
  {
    name: 'framework-library-behind',
    level: 'warning',
    what: 'Sentinel, GORM Studio or Pulse is below the version this CLI ships with, which for Sentinel means missing security fixes. grit upgrade raises them.',
  },
  {
    name: 'rate-limits-per-process',
    level: 'warning',
    what: 'Sentinel counts rate limits and lockouts in this process while Redis is configured, so each replica allows a client the full limit again.',
  },
  {
    name: 'public-allowlist-sensitive',
    level: 'warning',
    what: 'A --public allowlist publishes a column the generator holds back: a cost, a margin, a stock count, an internal note.',
  },
]

export default function DoctorPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Security &amp; Testing</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Project audit: grit doctor</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                <code>grit doctor</code> reads your project and reports the mistakes that do not
                announce themselves. Every check exists because the mistake it finds was made once,
                in a real project, and cost something to find.
              </p>
            </div>

            <div className="prose-grit">
              <p>
                A build catches a type error. A test catches a wrong answer. Neither catches an
                encrypted column with no key, a list that was never scoped to its owner, or a
                database browser whose password is still <code>studio</code>. Those work, which is
                the problem.
              </p>

              <CodeBlock terminal code={`$ grit doctor

  13 checks over 0 resource(s)

  ✓ Nothing to report.`} />

              <p>
                That is a project with nothing to say about it, which is the point: a linter that
                cries wolf is a linter people turn off. That claim is checked rather than asserted:
                the docs workflow scaffolds a project on every change to these pages and runs this,
                which fails if anything is reported.
              </p>

              <CodeBlock terminal verify="doctor-clean" code={`grit doctor`} />

              <p>
                Here is the same command on an older project with ten resources:
              </p>

              <CodeBlock terminal code={`$ grit doctor

  13 checks over 10 resource(s)

  ✗ Product has encrypted fields, and FIELD_ENCRYPTION_KEY is not set in .env, so they are stored in the clear
      set FIELD_ENCRYPTION_KEY to 32 random bytes (openssl rand -base64 32), or set it in the deployment's environment
      (encryption-key-unset)
  ⚠ /studio still has its default password, and it can browse and edit every table
      set GORM_STUDIO_PASSWORD in .env: openssl rand -hex 16
      (studio-unprotected)
  ⚠ SENTINEL_AUDIT_KEY is not set, so the security audit log's chain is unkeyed
      set SENTINEL_AUDIT_KEY in .env: openssl rand -hex 32
      (default-credentials)
  ⚠ github.com/MUKE-coder/sentinel/v2 is at v2.2.1, below v2.5.0: security fixes from v2.2.2
      run grit upgrade, which raises it
      (framework-library-behind)

  1 error(s), 5 warning(s)`} />

              <p>
                Every finding names the check it came from, so you can talk about it, and a fix you
                can act on. The first one there is a real bug from a live test: the column was
                readable in Postgres and the API had answered 201.
              </p>

              {/* ── What it checks ─────────────────────────────── */}
              <h2 id="checks">What it checks</h2>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Check</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Level</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">What it means</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    {checks.map((c) => (
                      <tr key={c.name} className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary/60 align-top whitespace-nowrap">{c.name}</td>
                        <td className="px-4 py-2.5 text-xs align-top whitespace-nowrap">{c.level}</td>
                        <td className="px-4 py-2.5">{c.what}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              {/* ── In CI ─────────────────────────────── */}
              <h2 id="ci">In CI</h2>
              <p>
                <code>grit doctor</code> exits non-zero when anything is an error, and{' '}
                <code>--json</code> prints the report for a machine to read:
              </p>
              <CodeBlock language="yaml" filename=".github/workflows/ci.yml" code={`- name: Audit the project
  run: grit doctor`} />
              <CodeBlock terminal code={`$ grit doctor --json
{
  "Findings": [
    {
      "Level": "error",
      "Check": "encryption-key-unset",
      "Resource": "Product",
      "Message": "has encrypted fields, and FIELD_ENCRYPTION_KEY is not set in .env, so they are stored in the clear",
      "Fix": "set FIELD_ENCRYPTION_KEY to 32 random bytes (openssl rand -base64 32), or set it in the deployment's environment"
    }
  ],
  "Resources": ["Invoice", "Product"],
  "Checks": 13
}`} />

              {/* ── What it reads ─────────────────────────────── */}
              <h2 id="what-it-reads">What it reads, and what it does not</h2>
              <p>
                It reads your project: the models, handlers, services, routes, <code>.env</code> and{' '}
                <code>go.mod</code>. It connects to nothing, so it is safe to run anywhere and it
                cannot check what only the database knows.
              </p>
              <p>
                Two things follow from that. A deployment that sets its secrets outside{' '}
                <code>.env</code> will still be told the keys are missing, because from here they
                are. And only generated resources are audited: the framework&apos;s own tables are
                skipped, because a session row referencing a user is not a user-owned resource.
                Asked the loose way, every project reported four framework tables as
                possibly-owned, including a project with no resources of its own.
              </p>

              {/* ── Its first finding ─────────────────────────────── */}
              <h2 id="first-finding">Its first finding</h2>
              <p>
                Run against Grit&apos;s own test projects, the first thing it found was in the
                scaffold: every project shipped <code>GORM_STUDIO_PASSWORD=studio</code>, a known
                password on a tool that browses and edits every table, while the Sentinel and Pulse
                passwords beside it were generated per project. From v3.227.0 the Studio password is
                generated too, so a new project has nothing to report here, and an older one is told.
              </p>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 mt-10 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/security/defenders-handbook" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Defender&apos;s Handbook
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/testing" className="gap-1.5">
                  Performance &amp; Pentest Testing
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
