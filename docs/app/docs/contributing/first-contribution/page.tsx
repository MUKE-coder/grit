import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/contributing/first-contribution')

export default function FirstContributionPage() {
  return (
    <DocsShell>
            <div className="mb-10">
              <span className="tag-mono text-accent/80 mb-3 block">Contributing</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Your first contribution</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                You learned Grit and you want to give back. Docs fixes are the usual first
                step: they are small, they are easy to check, and they help the next person
                who reads the page. The full guide is{' '}
                <a href="https://github.com/MUKE-coder/grit/blob/main/CONTRIBUTING.md" target="_blank" rel="noreferrer" className="text-accent hover:underline">CONTRIBUTING.md</a>{' '}
                in the repo root, and it is aimed at people changing templates. This page is
                the shorter path for docs fixes.
              </p>
            </div>

            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Fix a docs page
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Docs pages live in <code>docs/app/docs</code>, one <code>page.tsx</code> per
                page. Course lessons live in <code>docs/content/lessons</code>. Shared text
                such as versions and page metadata lives in <code>docs/config</code>.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Run the site with <code>pnpm dev</code> in <code>docs/</code>, then open the
                page you changed and check it reads right.
              </p>
              <CodeBlock language="bash" code={`cd docs\npnpm dev`} />
              <p className="text-muted-foreground leading-relaxed mb-4">
                A good first fix is something you noticed while reading: a stale version, a
                broken link, a command that does not do what the page says. Fix the page,
                check it renders, send the pull request.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Three rules for docs changes. First, never use an em dash in anything a
                reader sees. Use a colon, a comma or a full stop instead, and run{' '}
                <code>python scripts/check-em-dashes.py</code> in <code>docs/</code> before
                you send the pull request: it fails on a new em dash. Second, if you add or
                rename a page, regenerate the search index with{' '}
                <code>python scripts/search-index.py</code> in <code>docs/</code>, or the
                page is unfindable. Third, run <code>npx next build</code> in{' '}
                <code>docs/</code> before you send the pull request.
              </p>
            </div>

            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Find a good first issue
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Open issues labelled{' '}
                <a href="https://github.com/MUKE-coder/grit/labels/good%20first%20issue" target="_blank" rel="noreferrer" className="text-accent hover:underline">good first issue</a>{' '}
                are the starting list. Pick one that is small and well described: a docs
                fix, a missing test, a small bug.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Before you start, check the issue is still open and that no open pull
                request already covers it. If something is wrong and there is no issue,
                send the fix: fixes need no issue. Anything bigger than a fix starts with
                an issue first. Questions and ideas go to{' '}
                <a href="https://github.com/MUKE-coder/grit/discussions" target="_blank" rel="noreferrer" className="text-accent hover:underline">Discussions</a>, not the
                issue tracker.
              </p>
            </div>

            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                What a good issue report contains
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                A bug report needs the Grit version (<code>grit version</code>), where it
                happens, steps to reproduce starting from <code>grit new</code> if you can,
                what you expected, the actual output or error text, your operating system,
                and your Go and Node versions. Say that you searched the existing issues
                and that you are on the latest version.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-4">
                A feature request needs what you are trying to do, stated as the problem
                rather than the solution, what you do today instead, however ugly the
                workaround is, what it would look like, and where it would live. A feature
                that replaces forty lines everybody writes is an easier decision than one
                that replaces nothing.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Security problems do not go in an issue. Report them privately, as{' '}
                <a href="https://github.com/MUKE-coder/grit/blob/main/SECURITY.md" target="_blank" rel="noreferrer" className="text-accent hover:underline">SECURITY.md</a> says.
              </p>
            </div>

            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Send the fix
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                One change per pull request. Write the commit message for somebody reading
                it in two years with no context: what was wrong, and why the fix is the
                fix. Conventional prefixes apply: <code>feat:</code>, <code>fix:</code>,{' '}
                <code>docs:</code>, <code>refactor:</code>, <code>test:</code>,{' '}
                <code>chore:</code>.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-4">
                For code changes, <code>go test ./...</code> passes and{' '}
                <code>gofmt -l .</code> is quiet, and every behavioural change needs a test.
                Do not bump the version or edit the changelog: releases are cut from{' '}
                <code>main</code> by the maintainer.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-4">
                What merges quickly: docs that correct something wrong, a fix with a test
                that fails without it, and a bug reproduced in a scratch project.
              </p>
            </div>

            <div className="flex items-center justify-between border-t border-border pt-8 mt-12 max-w-3xl">
              <Button variant="ghost" asChild>
                <Link href="/docs" className="gap-2">
                  <ArrowLeft className="h-4 w-4" />
                  Introduction
                </Link>
              </Button>
              <Button variant="ghost" asChild>
                <Link href="https://github.com/MUKE-coder/grit/labels/good%20first%20issue" target="_blank" rel="noreferrer" className="gap-2">
                  Good first issues
                  <ArrowRight className="h-4 w-4" />
                </Link>
              </Button>
            </div>
    </DocsShell>
  )
}
