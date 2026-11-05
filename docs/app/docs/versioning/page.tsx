import Link from 'next/link'
import { AlertCircle, ArrowRight, Check } from 'lucide-react'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'
import { GRIT_VERSION } from '@/config/site'

export const metadata = getDocMetadata('/docs/versioning')

export default function VersioningPage() {
  return (
    <DocsShell>
            <div className="mb-10">
              <span className="tag-mono text-accent/80 mb-3 block">Project</span>
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
                className="text-accent hover:underline"
              >
                the issue tracker
              </Link>{' '}
              — it is worth knowing that a real team is blocked on it.
            </p>

            <h2 id="deprecation" className="text-2xl font-bold tracking-tight mb-4">
              Deprecation policy
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Nothing in the CLI is removed without notice. A command or a flag that is
              going away is <strong>deprecated first</strong>: it keeps working, it drops
              out of <code>--help</code>, and running it prints what to use instead.
            </p>
            <div className="rounded-xl border border-border/60 bg-card/40 p-6 mb-6">
              <p className="text-sm text-foreground font-medium mb-2">The window</p>
              <p className="text-sm text-muted-foreground leading-relaxed">
                A deprecated command or flag keeps working for at least{' '}
                <strong className="text-foreground">90 days and 30 releases</strong>,
                whichever is longer, and in practice until the next major. Both numbers,
                because this project ships several releases a week: thirty of them can pass
                in under a month, which is not notice. Ninety days alone would be a long
                time to hold something nobody uses.
              </p>
            </div>
            <p className="text-muted-foreground leading-relaxed mb-4">
              What that looks like. <code>grit new-desktop</code> does the same thing as{' '}
              <code>grit new --desktop</code>, and was deprecated in v3.467.0:
            </p>
            <CodeBlock
              language="bash"
              code={`$ grit new-desktop pos
Command "new-desktop" is deprecated, use \`grit new <project-name> --desktop\`
instead. This name will be removed no sooner than v4.0.0.`}
            />
            <p className="text-muted-foreground leading-relaxed mb-4">
              The message goes to stderr, so a script that pipes stdout still shows it to
              whoever runs the script. Three things are <em>not</em> covered by this, and
              it is worth saying which:
            </p>
            <ul className="space-y-2 mb-6 text-muted-foreground leading-relaxed">
              <li className="flex gap-3">
                <span className="text-accent mt-0.5">&bull;</span>
                <span>
                  <strong className="text-foreground">Generated code.</strong> It is yours,
                  so nothing is deprecated: the generator can change what it writes next
                  time without touching what it wrote last time. See above.
                </span>
              </li>
              <li className="flex gap-3">
                <span className="text-accent mt-0.5">&bull;</span>
                <span>
                  <strong className="text-foreground">A security fix.</strong> If something
                  has to go because leaving it is a risk, it goes, and the changelog says
                  why rather than hiding it in a notice period.
                </span>
              </li>
              <li className="flex gap-3">
                <span className="text-accent mt-0.5">&bull;</span>
                <span>
                  <strong className="text-foreground">Something nobody can be using.</strong>{' '}
                  A flag that has never worked is not a flag anybody depends on. The
                  changelog still says it went.
                </span>
              </li>
            </ul>

            <h2 className="text-2xl font-bold tracking-tight mb-4">Upgrading the CLI</h2>
            <CodeBlock
              language="bash"
              code={`grit update          # latest release
grit update --version 3.455.0   # a named release, not the latest
grit version         # what you are on

# Pin a version in CI so a release never changes your build under you:
go install github.com/MUKE-coder/grit/v3/cmd/grit@v${GRIT_VERSION}`}
            />
            <p className="text-muted-foreground leading-relaxed mb-4 mt-6">
              And a project can be moved to a named version without moving the CLI on the
              machine, which is how you step forward one release at a time rather than
              jumping from wherever you are to the latest:
            </p>
            <CodeBlock
              language="bash"
              code={`grit upgrade --to 3.455.0             # that version's framework files
grit upgrade --to 3.455.0 --dry-run   # what it would change, changing nothing`}
            />
            <p className="text-muted-foreground leading-relaxed mb-4 mt-4">
              The templates are compiled into the binary, so there is no archive of past
              templates to read: <code>--to</code> downloads that release, verifies it
              against the release&apos;s signed checksums, caches it under{' '}
              <code>~/.grit/versions</code> and runs it. Your installed CLI is not touched,
              and the download happens once per version.
            </p>

            <div className="flex flex-wrap items-center gap-6 border-t border-border/40 pt-8 mt-12">
              <Link
                href="/docs/changelog"
                className="inline-flex items-center gap-1.5 text-sm font-medium text-accent hover:underline"
              >
                Changelog
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
              <Link
                href="/docs/governance"
                className="inline-flex items-center gap-1.5 text-sm font-medium text-accent hover:underline"
              >
                Governance &amp; roadmap
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </div>
    </DocsShell>
  )
}
