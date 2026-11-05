import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { getDocMetadata } from '@/config/docs-metadata'
import { GRIT_VERSION } from '@/config/site'
import { STABILITY_AREAS, statusStyle } from '@/config/stability'

export const metadata = getDocMetadata('/docs/stability')

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
    label: 'The live suite on five databases',
    detail: 'the same 68 checks run on Postgres 15, 16 and 17, on MySQL 8, and on SQLite, with a smaller set on an in-memory database. Until v3.234.0 they ran on Postgres alone, so MySQL support rested on the driver compiling: the first project built on it hit a table AutoMigrate refused.',
  },
  {
    label: 'A tenancy fixture, run on every push',
    detail: 'CI scaffolds a project with the multitenant and impersonate plugins, a tenant-owned resource and one that is tenant-owned and owned by a user, then drives 26 checks across the seams: cross-tenant reads, an admin with no organization, a role held through a membership, and an impersonated session inside a tenant.',
  },
  {
    label: '106 error codes, one status each, checked',
    detail: 'every code a handler returns is resolved against the catalogue, and a code returned with a status other than its own fails the build. That test is why the statuses are consistent: VALIDATION_ERROR was 422 in thirty-eight handlers and 400 in twenty-five.',
  },
  {
    label: 'Every documented command, checked against the CLI',
    detail: 'all 586 of them are resolved against the real command tree, and every --fields spec is run through the generator\u2019s own parser. A renamed command or a field type that does not exist fails the build. Blocks marked verify= go further: CI runs them against a scaffolded project on a real Postgres.',
  },
  {
    label: 'Every architecture generated, run and looked at',
    detail:
      'the eight shapes the stack selector offers, each scaffolded fresh, built, type-checked, migrated, started and driven through a browser as a contact app. It found thirty faults in one pass, every one of them now a test: a Content-Security-Policy that authorised a different host from the one the client called, an owner picker whose answer the server discarded, a mobile app shipping the framework\u2019s logo as its own, a desktop app that did not build, and a Next single-binary project whose admin panel could not reach its API in development. Five of the thirty were the same fault in five places, which is why there is now a test for that family across every shape.',
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
    <DocsShell>
            <div className="mb-10">
              <span className="tag-mono text-accent/80 mb-3 block">Security &amp; Testing</span>
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
              {STABILITY_AREAS.map((row) => (
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

            {/* Release channels */}
            <div className="prose-grit mt-12">
              <h2 id="release-channels">Stable, and how to tell</h2>
              <p>
                The table above says which <em>features</em> are settled. This says which{' '}
                <em>releases</em> are, which is a different question and is answered entirely by
                the version number.
              </p>
              <p>
                A version with nothing after the patch number is stable:{' '}
                <code>{GRIT_VERSION}</code> is the current one. A version carrying a suffix is
                not, and the suffix says which kind:
              </p>
              <ul>
                <li>
                  <code>3.398.0-rc.1</code> is a <strong>release candidate</strong>: feature
                  complete, published so it can be tried before the real thing.
                </li>
                <li>
                  <code>3.398.0-beta.1</code> is a <strong>beta</strong>: the shape is settled,
                  the details are not.
                </li>
                <li>
                  <code>3.398.0-alpha.1</code> is an <strong>alpha</strong>: published to get
                  feedback, and expected to change.
                </li>
              </ul>
              <p>
                The suffix is not only a naming convention. The release workflow marks any tag
                containing <code>-rc</code>, <code>-beta</code> or <code>-alpha</code> as a GitHub
                prerelease, and <code>grit update</code> asks GitHub for the <em>latest release</em>,
                which by definition excludes prereleases. So a prerelease is something you install
                deliberately, and never something an update moves you onto.
              </p>
              <p>
                Everywhere this site shows a version, a prerelease is labelled next to it. No
                label means stable, which is why there is no badge saying so on every page.
              </p>
              <p>
                To install one on purpose, name it:{' '}
                <code>go install github.com/MUKE-coder/grit/v3/cmd/grit@v3.398.0-rc.1</code>, or
                take the binary from that release on GitHub. To go back,{' '}
                <code>grit update</code> returns you to the newest stable.
              </p>
              <p>
                Every release is cut when something ships rather than on a schedule, so the
                version number tells you what kind of release it is and nothing about how much
                changed. The{' '}
                <Link href="/docs/changelog">changelog</Link> is where that lives.
              </p>
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
                  <strong>No signed desktop release has been verified</strong> end to end. The
                  desktop app itself is built and run on every sweep; what is unverified is the
                  Authenticode and Apple Developer signing, which needs certificates this project
                  does not own.
                </li>
                <li>
                  <strong>The architecture sweep is a person, not a job.</strong> CI builds a
                  generated app and type-checks it, and the live suite drives the API over HTTP,
                  but nothing in CI opens a browser and looks at eight shapes. The thirty faults
                  that found were all invisible to a green build, and five of them were in work
                  that had shipped hours earlier the same day.
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
    </DocsShell>
  )
}
