import Link from 'next/link'
import { AlertCircle, ArrowRight, Check } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'
import { GRIT_VERSION } from '@/config/site'

export const metadata = getDocMetadata('/docs/versioning')

export default function VersioningPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Project</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Versioning &amp; breaking changes
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                What the version number means, what can change under you, and what happens
                to a project you generated a year ago. Currently on{' '}
                <code className="text-sm">v{GRIT_VERSION}</code>.
              </p>
            </div>

            {/* Address the actual question people arrive with. A high minor
                number with several releases a week reads as instability unless
                you explain what is moving. */}
            <div className="rounded-xl border border-border/50 bg-card/50 p-6 mb-12">
              <h2 className="font-semibold mb-3">
                &ldquo;Why is the version number so high?&rdquo;
              </h2>
              <p className="text-sm text-muted-foreground leading-relaxed mb-3">
                Because minor versions are cheap here and get used. Grit ships several times
                a week, and every shipped change gets a number rather than being batched into
                a quarterly release. A high minor is a record of release cadence, not of
                churn in your application.
              </p>
              <p className="text-sm text-muted-foreground leading-relaxed">
                The number that matters for stability is the <strong>major</strong>, and it
                has been <code className="text-xs">3</code> for a long time. That is the one
                that is allowed to break you.
              </p>
            </div>

            <h2 className="text-2xl font-bold tracking-tight mb-4">
              The thing most frameworks do not have to explain
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Grit is a code generator, not a runtime dependency. Your application does not
              import Grit and does not link against it. That changes what &ldquo;upgrading&rdquo;
              even means:
            </p>
            <div className="grid gap-3 sm:grid-cols-2 mb-6">
              <div className="rounded-xl border border-border/50 bg-card/40 p-5">
                <div className="text-sm font-semibold mb-2">A normal framework</div>
                <p className="text-[13px] text-muted-foreground leading-relaxed">
                  Bump the dependency, and every behaviour it controls changes at once.
                  Upgrading is a single risky step for the whole app.
                </p>
              </div>
              <div className="rounded-xl border border-emerald-500/25 bg-emerald-500/[0.04] p-5">
                <div className="text-sm font-semibold mb-2">Grit</div>
                <p className="text-[13px] text-muted-foreground leading-relaxed">
                  Upgrading the CLI changes what the <em>next</em> generated file looks like.
                  Code already in your repo is yours and does not move. There is no version
                  of Grit that can break a running application.
                </p>
              </div>
            </div>
            <p className="text-muted-foreground leading-relaxed mb-12">
              So the practical question is not &ldquo;is it safe to upgrade&rdquo; — it is
              always safe — but &ldquo;will the next thing I generate still fit alongside
              what I generated last year&rdquo;. That is what the guarantees below cover.
            </p>

            <h2 className="text-2xl font-bold tracking-tight mb-4">What each number means</h2>
            <div className="overflow-x-auto rounded-xl border border-border/50 mb-12">
              <table className="w-full text-sm">
                <thead className="bg-card/60">
                  <tr className="text-left">
                    <th className="px-4 py-3 font-semibold">Change</th>
                    <th className="px-4 py-3 font-semibold">Bumps</th>
                    <th className="px-4 py-3 font-semibold">Example</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border/40 align-top">
                  {[
                    ['New generator, new battery, new command', 'Minor', '`grit swap` arriving in v3.115.0'],
                    ['New field type or flag on an existing command', 'Minor', '`--items` on `generate resource`'],
                    ['Generated output changes shape for NEW resources', 'Minor', 'A generated form gaining a per-step save button'],
                    ['Bug fix in a template', 'Patch', 'A generated handler missing a nil check'],
                    ['A CLI command is renamed or removed', 'Major', '—'],
                    ['A `grit:` marker contract changes', 'Major', '—'],
                    ['The scaffolded project layout moves', 'Major', 'apps/api → services/api'],
                    ['A slot contract changes incompatibly', 'Major', '`button@1` → `button@2`'],
                  ].map(([change, bump, example]) => (
                    <tr key={change}>
                      <td className="px-4 py-3 text-muted-foreground">{change}</td>
                      <td className="px-4 py-3">
                        <span
                          className={
                            'inline-flex rounded-md px-2 py-0.5 text-xs font-semibold ' +
                            (bump === 'Major'
                              ? 'bg-amber-500/15 text-amber-600 dark:text-amber-400'
                              : 'bg-emerald-500/15 text-emerald-600 dark:text-emerald-400')
                          }
                        >
                          {bump}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-muted-foreground text-[13px]">{example}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <h2 className="text-2xl font-bold tracking-tight mb-4">What will not change in a minor</h2>
            <ul className="space-y-2.5 mb-12">
              {[
                'The meaning of a `grit:` marker. Regeneration finds the same anchors it always did.',
                'The scaffolded project layout. Where a model, handler or resource file lives is stable within a major.',
                'The name or argument shape of an existing CLI command.',
                'A slot contract (`button@1`). Adding a prop is additive; removing or narrowing one is a new major.',
                'Field-type syntax in `--fields`. New types get added; existing ones keep meaning what they meant.',
              ].map((t) => (
                <li key={t} className="flex gap-3 text-sm text-muted-foreground leading-relaxed">
                  <Check className="h-4 w-4 shrink-0 text-emerald-500 mt-0.5" />
                  {t}
                </li>
              ))}
            </ul>

            <h2 className="text-2xl font-bold tracking-tight mb-4">Generated code is yours</h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Nothing in a Grit release rewrites files you already have. Generators write new
              files and inject into marked regions of existing ones — never outside them.
            </p>
            <CodeBlock
              language="go"
              code={`// grit:models:auto-start
// Everything between these markers is regenerated.
// Edits here are lost on the next run.
&models.User{},
&models.Product{},
// grit:models:auto-end

// Anything outside the markers is yours, permanently.
// Generators do not read it and will not touch it.`}
            />
            <p className="text-muted-foreground leading-relaxed mt-4 mb-12">
              That boundary is the contract. If a release ever moves a marker, that is a major
              version and it will be in the changelog with a migration note.
            </p>

            <h2 className="text-2xl font-bold tracking-tight mb-4">Support policy</h2>
            <div className="rounded-xl border border-amber-500/30 bg-amber-500/[0.05] p-6 mb-8">
              <div className="flex gap-3">
                <AlertCircle className="h-5 w-5 shrink-0 text-amber-500 mt-0.5" />
                <div>
                  <p className="text-sm text-muted-foreground leading-relaxed mb-3">
                    <strong className="text-foreground">Being straight about this:</strong>{' '}
                    Grit is a young project with a small maintainer team. There is no
                    contractual LTS, no backport branch, and no support SLA. Claiming one
                    would be more useful to marketing than to you.
                  </p>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    What there is instead: fixes land on the current major, the CLI is
                    independent of your running application, and generated code has no
                    upstream to go stale. An app generated on v3.0 still runs today with no
                    Grit installed anywhere near it — which is a stronger guarantee than most
                    LTS promises, and it comes from the architecture rather than from a
                    commitment anyone has to keep.
                  </p>
                </div>
              </div>
            </div>
            <p className="text-muted-foreground leading-relaxed mb-12">
              If you need something firmer than that for a procurement process, say so on{' '}
              <Link
                href="https://github.com/MUKE-coder/grit/issues"
                target="_blank"
                rel="noreferrer"
                className="text-primary hover:underline"
              >
                the issue tracker
              </Link>{' '}
              — it is worth knowing that a real team is blocked on it.
            </p>

            <h2 className="text-2xl font-bold tracking-tight mb-4">Upgrading the CLI</h2>
            <CodeBlock
              language="bash"
              code={`grit update          # latest release
grit version         # what you are on

# Pin a version in CI so a release never changes your build under you:
go install github.com/MUKE-coder/grit/v3/cmd/grit@v${GRIT_VERSION}`}
            />

            <div className="flex flex-wrap items-center gap-6 border-t border-border/40 pt-8 mt-12">
              <Link
                href="/docs/changelog"
                className="inline-flex items-center gap-1.5 text-sm font-medium text-primary hover:underline"
              >
                Changelog
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
              <Link
                href="/docs/governance"
                className="inline-flex items-center gap-1.5 text-sm font-medium text-primary hover:underline"
              >
                Governance &amp; roadmap
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
