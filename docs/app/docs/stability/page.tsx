import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/stability')

type Status = 'stable' | 'beta' | 'new'

const statusStyle: Record<Status, string> = {
  stable: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
  beta: 'bg-amber-500/10 text-amber-400 border-amber-500/20',
  new: 'bg-sky-500/10 text-sky-400 border-sky-500/20',
}

type Row = {
  area: string
  status: Status
  guaranteed: string
  proof: string
  yours: string
  history?: string
}

// One row per subsystem. The rule for this table: a claim in "Grit guarantees"
// has to name the test that would fail if it stopped being true. Where there is
// no such test, the claim belongs in "Your responsibility" instead, however
// much the code looks right.
const rows: Row[] = [
  {
    area: 'Auth: JWT, sessions, refresh rotation',
    status: 'stable',
    guaranteed:
      'Every refresh token is backed by a session row, rotation detects replay, and a password change signs out every device.',
    proof: 'handlers/auth_test.go, services/session_test.go, services/password_reset_test.go (41 handler tests, 75 service tests ship in your project)',
    yours: 'Set JWT_SECRET to 32+ random bytes and keep it out of the repository. grit doctor fails if it is short, empty or a placeholder.',
    history: 'Two logins in the same second once produced identical refresh tokens (fixed in v3.86.0: every token carries a jti).',
  },
  {
    area: 'Two-factor, passkeys, recovery codes',
    status: 'beta',
    guaranteed: 'TOTP is RFC 6238 (HMAC-SHA1 by the spec), recovery codes are single-use, and passkey challenges are bound to a session.',
    proof: 'services/passkey_test.go, services/recovery_test.go, crypto/field_test.go',
    yours: 'Decide your own lockout and recovery policy, and test the flow with the authenticator app your users actually have.',
  },
  {
    area: 'RBAC: roles, permissions, the staff group',
    status: 'beta',
    guaranteed:
      'A custom role reaches only the endpoints its permissions name, a revoked permission stops working everywhere, and a route that names no permission stays ADMIN-only, so one that forgets fails closed.',
    proof: 'authz/grants_test.go, authz/permissions_test.go, handlers/role_test.go, handlers/user_role_sync_test.go, plus a live check that a staff account holding notes.delete still cannot delete another user’s row',
    yours: 'Name a permission on every route you add yourself. The generator does it for generated resources; a hand-written route is yours.',
    history: 'A custom role could not reach a single admin endpoint until v3.220.0; a revoked permission kept working on other replicas until v3.218.0.',
  },
  {
    area: 'Generated CRUD: services, handlers, pagination',
    status: 'stable',
    guaranteed:
      'Handlers run no queries; the service owns every read and write, takes a context, and whitelists what a client may search, sort, filter and patch, because those names reach SQL.',
    proof: 'internal/generate tests in the CLI (492 test functions across 99 files), an import-versus-use check on every generated file, and 57 live checks against a running app on Postgres 15, 16 and 17',
    yours: 'Put your own queries in the service, not the handler, so a job and a route cannot disagree.',
    history: 'Handlers ran all of their own queries until v3.224.0, and the service written beside them was never called.',
  },
  {
    area: 'Owned resources (--owned-by)',
    status: 'stable',
    guaranteed:
      'List, export, get, PDF, update, patch, delete and bulk are all scoped to the caller; somebody else’s row is 404 rather than 403; the owner comes from the session and never from the request body. ADMIN is exempt.',
    proof: 'Live checks cover every one of those paths, including the staff-permission case, plus generator tests that fail if a method loses its scoping. grit doctor reports the same thing in your project.',
    yours: 'Add --owned-by when rows belong to a user. A resource that merely references a user is not scoped, and grit doctor asks about it.',
    history: 'The export, PDF, patch and bulk routes were unscoped once: a second ordinary account printed another user’s record and rewrote it.',
  },
  {
    area: 'Optimistic locking (If-Match)',
    status: 'stable',
    guaranteed:
      'A read returns the version as an ETag, a write against a stale version is 409 naming the current one, and of twenty simultaneous writes on one version exactly one lands.',
    proof: 'concurrency/concurrency_test.go in your project, generator tests, and a 20-way race in the live suite on all three Postgres versions',
    yours: 'Send If-Match from your own clients. Without it the last write still wins, by design.',
    history: 'Added in v3.219.0, after two people saving the same record meant the second silently won.',
  },
  {
    area: 'Encryption at rest (:encrypted fields)',
    status: 'beta',
    guaranteed:
      'An encrypted column is ciphertext in the database and plaintext through the API, and it is kept out of search, sort and filter lists, where ciphertext would match nothing.',
    proof: 'crypto/field_test.go, crypto/map_update_test.go, and a live check that reads the column straight out of Postgres',
    yours:
      'Set FIELD_ENCRYPTION_KEY. Without it the value is stored as it is, and nothing fails: grit doctor makes this an error, and it is the first thing it found in a real project.',
    history: 'A plaintext leak was fixed in v3.212.0. Re-audit before storing real PII.',
  },
  {
    area: 'Multitenancy (plugin, tenant.Owned)',
    status: 'beta',
    guaranteed: 'A tenant-owned model is scoped by the organization on the request context, and a query with no organization is refused rather than answered.',
    proof: 'The plugin’s own tests, and grit doctor reports a resource with no tenant.Owned in a project that has the plugin.',
    yours:
      'Add tenant.Owned to every model that belongs to an organization, and verify isolation with two organizations before you trust it. Combining it with --tree or --public is not covered by the live suite yet.',
    history: 'Isolation bugs were found and fixed in v3.197 and v3.198. This needs more adversarial testing before trusting it blind.',
  },
  {
    area: 'Money (exact, multi-currency)',
    status: 'stable',
    guaranteed: 'Amounts are integer minor units with their currency, two money columns per row do not collide, and sums stay exact.',
    proof: 'money/money_test.go (13 tests), generator tests for the embedded columns and whitelists, and a live round-trip check',
    yours: 'Pick the currency per row deliberately; Grit does not convert between currencies.',
  },
  {
    area: 'CSV/XLSX import and export',
    status: 'beta',
    guaranteed:
      'An export streams every matching row; an import runs in the background with a job to poll, resolves relations, skips duplicates, and an owned resource’s rows belong to whoever imported them unless an ADMIN names another owner.',
    proof: 'Live checks for both, including the owned cases, and generator tests for the importer’s columns',
    yours: 'Validate the files your users upload. The parser is hardened and bounded, and it is still parsing a stranger’s spreadsheet.',
    history: 'Every export was an empty 200 once. The importer created a user for every unrecognised email until that was removed.',
  },
  {
    area: 'Trees (--tree)',
    status: 'beta',
    guaranteed:
      'A move carries its subtree, a reorder keeps the parent the node has, and a node cannot be moved under its own descendant.',
    proof: 'services/category_tree_test.go ships with the resource (10 tests), plus live checks for the move, the reorder and the refusal',
    yours: 'Run the rebuild endpoint once after adding --tree to a table that already had rows.',
    history: 'Dragging a node that predated --tree once added one to the depth of every row in the table.',
  },
  {
    area: 'Public API (--public)',
    status: 'beta',
    guaranteed:
      'Public reads are an allowlist struct rather than the model, they are guarded by an API key, archived rows are excluded, and cost-shaped columns are held back by default.',
    proof: 'Live checks for the key, the allowlist, the filters and the archived case; grit doctor reports a held-back column that was published by hand',
    yours:
      'Review the allowlist when you add fields; it is your file and the generator will not rewrite it. Public responses are cached by URL for cache.public_ttl_seconds (60 by default), so an archived row can linger in a cached list for up to a minute.',
  },
  {
    area: 'Append-only resources',
    status: 'beta',
    guaranteed: 'Rows are created and read, never changed or deleted: a GORM guard, no mutating routes, and a database trigger that refuses the write.',
    proof: 'appendonly/appendonly_test.go, and grit doctor reports an append-only resource that still mounts a write',
    yours: 'Run grit migrate so the trigger exists. The guard alone is not the database refusing.',
  },
  {
    area: 'Realtime (WebSockets, Redis backplane)',
    status: 'beta',
    guaranteed: 'Events reach a user on every replica through the backplane, and a closed socket is cleaned up.',
    proof: 'realtime/backplane_test.go, cluster/cluster_test.go',
    yours: 'Run Redis if you run more than one replica. Without it, events stay in the process that emitted them.',
    history: 'Cross-replica and revocation bugs were fixed in v3.193 and v3.196.',
  },
  {
    area: 'Durable events (outbox)',
    status: 'beta',
    guaranteed: 'A durable subscriber’s work is committed with the row that triggered it, and the relay delivers it once.',
    proof: 'outbox/outbox_test.go (12 tests)',
    yours: 'Start the relay. grit upgrade wires it for projects that predate it.',
  },
  {
    area: 'Backups and restore',
    status: 'beta',
    guaranteed: 'A backup streams row by row, and a restore writes parents before children.',
    proof: 'The CLI’s backup tests and the generated project’s backup command',
    yours: 'Restore into a scratch database and look at the rows. A backup nobody has restored is a hope.',
    history: 'Restore correctness bugs were fixed in v3.207.0 and v3.213.0.',
  },
  {
    area: 'Offline sync (desktop)',
    status: 'new',
    guaranteed: 'Registered models mirror to the client, and conflicts resolve by the policy you set.',
    proof: 'The sync engine’s tests in the CLI, and grit sync doctor reports a model that cannot work (no version column, an allowlist naming a column that does not exist)',
    yours: 'Run grit sync doctor, and test the flows your users will have offline.',
  },
  {
    area: 'Workflows (status state machines)',
    status: 'new',
    guaranteed: 'Only declared transitions are allowed, they are guarded by role, and a transition hook commits with the move or not at all.',
    proof: 'The generated workflow service’s tests, plus generator tests for the transition route',
    yours: 'Decide who may make each transition. The generator gives you the shape, not your policy.',
    history: 'Shipped in v3.221.0: before it, a workflow could change a status and nothing else.',
  },
  {
    area: 'Feature flags',
    status: 'new',
    guaranteed: 'A flag can target attributes, a percentage and a date window, and flags.IsEnabled works from anywhere, answering false before the engine starts.',
    proof: 'flags/flags_test.go ships in every project and runs against a real engine',
    yours: 'Manage flags through the API: there is no admin screen for them yet.',
    history: 'Shipped in v3.223.0, when the documented package-level API did not exist at all.',
  },
  {
    area: 'Desktop (Wails) and mobile (Expo)',
    status: 'new',
    guaranteed: 'Both share the monorepo API and its generated types.',
    proof: 'The desktop module builds in CI; the Expo app is built by the nightly canary',
    yours:
      'Signing and store submission. A signed desktop release has never been verified end to end, because that needs Authenticode and Apple Developer certificates this project does not own.',
  },
]

const testing: { label: string; detail: string }[] = [
  {
    label: '492 test functions in the CLI',
    detail: 'across 99 files, over the generator, the scaffolder and every command. They run with -race on every push.',
  },
  {
    label: '256 tests ship into your project',
    detail: 'across 46 files: authz, concurrency, crypto, money, files, media, outbox, stock, paginate, realtime, webhooks, appendonly and the services. They are yours to run and to extend.',
  },
  {
    label: '57 live checks on Postgres 15, 16 and 17',
    detail: 'CI scaffolds a project, generates the resource shapes past bugs lived in, migrates, starts the server and drives it over HTTP, reading Postgres directly where the database is the only witness.',
  },
  {
    label: '13 project checks in grit doctor',
    detail: 'the mistakes that fail silently, reported in your own project: an encrypted field with no key, a resource nothing scopes, a table shared across organizations, a database browser with no login.',
  },
  {
    label: 'gosec, govulncheck, Trivy, CodeQL, Scorecard',
    detail: 'gosec gates the CLI and a generated app, govulncheck reports what is reachable, and Trivy reads the dependency graph. That last one found five HIGH CVEs in a freshly scaffolded project that reachability analysis does not see.',
  },
  {
    label: 'A nightly dependency canary',
    detail: 'scaffolds with no lockfile and builds every frontend, because a generated project resolves its dependencies at your install time, not at our release time.',
  },
]

export default function StabilityPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-4xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Security &amp; Testing</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Stability and hardening</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                What is proven, what is young, and where the line runs between what Grit guarantees
                and what stays your job. Per subsystem, so you can make your own risk call on the
                part you are about to depend on instead of on the whole framework at once.
              </p>
            </div>

            <div className="prose-grit">
              <p>
                The rule for the table below: a claim in <strong>Grit guarantees</strong> names the
                test that would fail if it stopped being true. Where no such test exists, the claim
                sits in <strong>your responsibility</strong> instead, however right the code looks.
                Three of the worst entries in the{' '}
                <Link href="/docs/changelog">changelog</Link> were docs asserting a guarantee the
                code did not provide, so the discipline is the point.
              </p>
              <ul>
                <li>
                  <span className={`inline-flex items-center rounded border px-2 py-0.5 text-xs font-medium ${statusStyle.stable}`}>stable</span>{' '}
                  long-lived, covered by tests and live checks, and no correctness bug found in it recently.
                </li>
                <li>
                  <span className={`inline-flex items-center rounded border px-2 py-0.5 text-xs font-medium ${statusStyle.beta}`}>beta</span>{' '}
                  works and is tested, and has had a real bug found in it recently enough that you should verify your own case.
                </li>
                <li>
                  <span className={`inline-flex items-center rounded border px-2 py-0.5 text-xs font-medium ${statusStyle.new}`}>new</span>{' '}
                  shipped in the last few weeks. Integration bugs are still plausible.
                </li>
              </ul>
            </div>

            {/* The matrix */}
            <div className="mt-8 space-y-4">
              {rows.map((row) => (
                <div key={row.area} className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <div className="flex flex-wrap items-center gap-3 mb-3">
                    <h2 className="text-base font-semibold tracking-tight">{row.area}</h2>
                    <span className={`inline-flex items-center rounded border px-2 py-0.5 text-xs font-medium ${statusStyle[row.status]}`}>
                      {row.status}
                    </span>
                  </div>
                  <dl className="grid gap-3 text-sm sm:grid-cols-2">
                    <div>
                      <dt className="text-xs uppercase tracking-wide text-emerald-400/80 mb-1">Grit guarantees</dt>
                      <dd className="text-muted-foreground leading-relaxed">{row.guaranteed}</dd>
                      <dt className="text-xs uppercase tracking-wide text-muted-foreground/60 mt-2 mb-1">Proven by</dt>
                      <dd className="text-muted-foreground/80 leading-relaxed text-xs font-mono">{row.proof}</dd>
                    </div>
                    <div>
                      <dt className="text-xs uppercase tracking-wide text-amber-400/80 mb-1">Your responsibility</dt>
                      <dd className="text-muted-foreground leading-relaxed">{row.yours}</dd>
                      {row.history && (
                        <>
                          <dt className="text-xs uppercase tracking-wide text-muted-foreground/60 mt-2 mb-1">What went wrong once</dt>
                          <dd className="text-muted-foreground/80 leading-relaxed">{row.history}</dd>
                        </>
                      )}
                    </div>
                  </dl>
                </div>
              ))}
            </div>

            {/* How Grit is tested */}
            <div className="prose-grit mt-12">
              <h2 id="how-we-test">How Grit is tested</h2>
              <p>
                Nearly every serious bug in Grit&apos;s history was found the same way: by building a
                real application on it and using it until something was wrong. That method works, and
                it only protects the next release if what it found is kept as a test. So it is.
              </p>
            </div>
            <div className="mt-4 space-y-3">
              {testing.map((item) => (
                <div key={item.label} className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                  <div className="text-sm font-medium text-foreground/90">{item.label}</div>
                  <div className="text-sm text-muted-foreground leading-relaxed mt-1">{item.detail}</div>
                </div>
              ))}
            </div>

            <div className="prose-grit mt-10">
              <p>
                The detail is in{' '}
                <Link href="/docs/testing#live-checks">how Grit itself is tested</Link>, and the
                project-side audit is{' '}
                <Link href="/docs/security/doctor">grit doctor</Link>.
              </p>

              <h2 id="not-yet">What is not covered yet</h2>
              <p>Stated plainly, because a matrix that only lists strengths is marketing:</p>
              <ul>
                <li>
                  <strong>No independent security audit.</strong> Every bug in auth, tenancy and
                  encryption listed above was found by the person who wrote it. An adversarial review
                  looks for what should not be possible, which is a different search.
                </li>
                <li>
                  <strong>The framework&apos;s own handlers still query directly.</strong> Auth,
                  two-factor, uploads, form shares and the dashboard read and write from their
                  handlers. Generated resources do not, since v3.225.0.
                </li>
                <li>
                  <strong>Multitenancy is not covered by the live suite</strong>, and not at all in
                  combination with <code>--tree</code> or <code>--public</code>.
                </li>
                <li>
                  <strong>No LTS channel.</strong> Releases land daily on one track. If you need a
                  slower train, pin a version and read the changelog before moving.
                </li>
                <li>
                  <strong>No signed desktop release has been verified</strong> end to end.
                </li>
                <li>
                  <strong>No production case studies yet.</strong> Benchmarks and tests are not the
                  same evidence as somebody else&apos;s traffic.
                </li>
              </ul>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 mt-10 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/security/doctor" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Project audit (grit doctor)
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
