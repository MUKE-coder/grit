import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/upgrades')

export default function UpgradesPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Stability</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Upgrading a project</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                A scaffolder gives you code once and then watches it rot. Grit keeps fixing the
                project it generated: <code>grit upgrade</code> brings an existing app up to the
                current templates, and it knows which files you have edited, so it never overwrites
                your work to do it.
              </p>
            </div>

            <div className="prose-grit">
              <h2>The two commands people confuse</h2>
              <CodeBlock
                language="bash"
                code={`grit update     # updates the CLI binary on your machine
grit upgrade    # run inside a project: brings THAT project up to the CLI's version`}
              />
              <p>
                They are different commands and the error from using the wrong one is confusing, so
                it is worth reading twice. <code>grit upgrade</code> is the one that touches your
                code.
              </p>

              <h2>It knows what you have edited</h2>
              <p>
                Grit records a hash of every file it writes, in <code>.grit/manifest.json</code>.
                That gives three answers for any generated file, and they lead to three different
                behaviours:
              </p>
              <ul>
                <li>
                  <strong>unchanged</strong>: the bytes on disk are the bytes Grit wrote, so it is
                  safe to replace. It gets the new version.
                </li>
                <li>
                  <strong>modified</strong>: you have edited it. The upgrade reports a conflict and
                  leaves your file alone. Nothing you wrote is lost by running this command.
                </li>
                <li>
                  <strong>untracked</strong>: Grit never wrote it. It is yours and is not touched.
                </li>
              </ul>
              <p>
                That is why the manifest exists, and it is also why an agent working in your project
                should read it before proposing an edit: changing a pristine file is the moment that
                file stops receiving framework fixes. <code>grit mcp</code> exposes it as{' '}
                <a href="/docs/ai-workflows/mcp">grit_file_ownership</a>.
              </p>
              <CodeBlock
                language="bash"
                code={`grit upgrade --diff     # show what it would change in the files it leaves alone
grit upgrade            # do it
grit upgrade --plan     # write a runbook for the files it left alone
grit upgrade --force    # overwrite everything, including your edits`}
              />

              <h2>The part that needs judgment</h2>
              <p>
                Leaving an edited file alone is a protection, and it is also a limit. What
                the new version would have changed in that file goes unapplied, and before
                v3.359.0 it was reported as one line of terminal output that scrolls away.
              </p>
              <p>
                There is no mechanical answer: the reason the file was skipped is that
                applying the change needs judgment about code you wrote. But a change that
                needs judgment is one a coding agent can carry out when it is told what the
                change is, where to look and how to check the result.
              </p>
              <p>
                <code>grit upgrade --plan</code> writes <code>UPGRADE-PLAN.md</code> at your
                project root: which version you came from, every file left alone, the diff
                of what the new version does to it, and a command whose success means it
                worked. An instruction with no verification step is one an agent will report
                as done whether or not it is, so every entry has all three.
              </p>
              <CodeBlock
                language="markdown"
                code={`### 1. \`docker-compose.yml\`

**Detect.** Open it, and see what you changed:

    git diff HEAD -- docker-compose.yml

**Change.** This is what the new version does to it:

    -    restart: unless-stopped
    +    container_name: shop-postgres
    +    restart: unless-stopped

**Verify.**

    docker compose -f docker-compose.yml config > /dev/null`}
              />
              <p>
                The diff is what Grit would have written against what you have. It is not a
                patch to apply blindly: your edits are in there for a reason, and the job is
                to carry the new behaviour into them rather than to replace one with the
                other.
              </p>

              <h2>The files you own, that Grit also writes into</h2>
              <p>
                Some files are genuinely shared. <code>routes.go</code> carries your routes and
                Grit&apos;s, <code>main.go</code> wires your services and Grit&apos;s,{' '}
                <code>config.go</code> holds your settings and Grit&apos;s. Replacing those would
                throw away your application; leaving them alone would mean a framework fix never
                reaches the file where it matters.
              </p>
              <p>
                So the upgrade carries 53 <em>repairs</em>. A repair finds the exact text Grit wrote,
                replaces it with the exact text Grit writes today, and when the text is not what Grit
                wrote, says so and changes nothing:
              </p>
              <CodeBlock
                language="text"
                code={`✓ apps/api/internal/routes/routes.go: /api/health answers ok, degraded, off or
  unknown per component, reports storage, and merges anything health.Register was given

⚠ the verification email is not sent the way Grit wrote it: send it with dispatchMail
  rather than from a \`go\` statement, so it is queued, retried and survives a restart`}
              />
              <p>
                The warning is the interesting half. Grit found code it did not write, so it did not
                touch it, and it told you what the fix would have been and why. You can take it or
                leave it.
              </p>

              <Callout type="note" title="An upgraded project is a fresh one">
                Each repair is tested against the file the previous release actually wrote, and the
                assertion is that the result is byte-for-byte what a fresh scaffold produces today.
                An upgraded project that differs from a new one by a blank line is a diff somebody
                has to read on every later upgrade.
              </Callout>

              <h2>What travels, and why that is a decision</h2>
              <p>
                Two sets of files are rewritten on every upgrade rather than only at scaffold time:
              </p>
              <ul>
                <li>
                  <strong>Framework-owned code</strong>: the webhook cluster, feature flags, the
                  permission engine, the audit services, the health package. Nothing describes them
                  in a resource definition and they only make sense as a set. A fix to one that did
                  not reach existing projects would be a fix for new projects only.
                </li>
                <li>
                  <strong>The codegen runtime</strong>: the packages generated code imports, such as{' '}
                  <code>internal/paginate</code> and <code>internal/export</code>. The generator
                  always emits against the CLI&apos;s current templates, so a project whose paginate
                  package is six versions old would get a handler referencing a field its own copy of
                  the struct does not declare, and the API would stop compiling with nothing in the
                  error to point at the cause.
                </li>
              </ul>
              <p>
                Your models, your services, your handlers and your frontend code are not in either
                set and are never rewritten.
              </p>

              <h2>What to do after</h2>
              <CodeBlock
                language="bash"
                code={`cd apps/api && go mod tidy    # if the upgrade added a dependency
pnpm install                  # if it added a frontend one
grit doctor                   # the audit, which is where a new check shows up
grit test                     # the project's own suites`}
              />
              <p>
                The upgrade tells you which of these it needs. <code>grit doctor</code> is worth
                running either way: a release often adds a check before it adds a fix, and the check
                is how you find out the fix applies to you.
              </p>

              <h2>How often</h2>
              <p>
                Grit releases often, and the version number is deliberately boring about it:{' '}
                <a href="/docs/versioning">a minor never breaks a running application</a>. Upgrading
                is a thing to do when you want something from the changelog, not a thing to keep up
                with. A project three months behind upgrades in one command, because every repair
                checks its own preconditions and the ones that do not apply do nothing.
              </p>

              <div className="mt-16 flex items-center justify-between border-t border-border/30 pt-8">
                <Link href="/docs/stability">
                  <Button variant="outline" size="sm" className="gap-2">
                    <ArrowLeft className="h-4 w-4" />
                    Stability and hardening
                  </Button>
                </Link>
                <Link href="/docs/versioning">
                  <Button variant="outline" size="sm" className="gap-2">
                    Versioning
                    <ArrowRight className="h-4 w-4" />
                  </Button>
                </Link>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
