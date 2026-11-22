import Link from 'next/link'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/contributing')

/**
 * Your first contribution, for somebody who used Grit rather than somebody
 * already inside the repository.
 *
 * Deliberately not a copy of CONTRIBUTING.md. That file is for a person who
 * has already decided to send a patch and needs the conventions. This page is
 * for the earlier moment: somebody who hit a rough edge, is mildly annoyed,
 * and does not know that the rough edge IS the contribution.
 *
 * So the order is: the things that are worth more than code, then the small
 * code changes, then the mechanics, then the honest bit about what gets
 * refused.
 */
export default function ContributingPage() {
  return (
    <DocsShell>
      <div className="mb-10">
        <span className="tag-mono text-accent/80 mb-3 block">Community</span>
        <h1 className="mb-4 text-4xl font-bold tracking-tight">Your first contribution</h1>
        <p className="text-lg leading-relaxed text-muted-foreground">
          The most useful thing you can send is not code. It is the thing that confused you,
          written down while you still remember being confused, because nobody who already knows
          the answer can see it any more.
        </p>
      </div>

      <div className="prose-grit">
        <h2>Start with what went wrong for you</h2>
        <p>
          If something took you an hour that should have taken five minutes, that hour is
          information the project cannot get any other way. The people who maintain it have the
          whole thing in their heads, which makes them the worst possible judges of what a
          newcomer needs.
        </p>
        <p>Four things in that category, in rough order of how much they help:</p>
        <ul>
          <li>
            <strong>An error message that did not explain itself.</strong> Open an issue with
            the exact text and what you were doing. If a page at{' '}
            <Link href="/e">/e</Link> would have saved you, say so: those pages exist to be
            linked from the CLI, and adding one is a small change.
          </li>
          <li>
            <strong>A documentation page that assumed something you did not know.</strong> Name
            the sentence. &quot;This page assumes I know what a resource definition is&quot; is a
            perfect issue.
          </li>
          <li>
            <strong>A command that did something surprising.</strong> Even if it was your
            mistake: a command that is easy to misuse is a command with a design problem.
          </li>
          <li>
            <strong>A bug, with the steps to reproduce it.</strong> A reproduction is worth more
            than a diagnosis, and far more than a guess at the cause.
          </li>
        </ul>

        <h2>Small code changes that are genuinely wanted</h2>
        <p>
          These are the ones labelled{' '}
          <a
            href="https://github.com/MUKE-coder/grit/labels/good%20first%20issue"
            target="_blank"
            rel="noreferrer"
          >
            good first issue
          </a>
          , and the kinds of thing they usually are:
        </p>
        <ul>
          <li>
            <strong>An explanation page.</strong> One file under{' '}
            <code>docs/app/e/&lt;slug&gt;/</code> plus an entry in{' '}
            <code>internal/clierr</code>. A test pairs the two, so you cannot ship one without
            the other.
          </li>
          <li>
            <strong>A glossary entry.</strong> One object in{' '}
            <code>docs/app/docs/glossary/page.tsx</code>. If a word confused you, it will confuse
            the next person.
          </li>
          <li>
            <strong>A guide.</strong> The hardest part of a guide is having been stuck recently,
            and you have. See the thirty at <Link href="/guides">/guides</Link> for the shape:
            answer it in plain Go and React first, and leave the one-command version until the
            end.
          </li>
          <li>
            <strong>A failing test for a bug you found.</strong> Even without the fix. A test
            that reproduces a bug is most of the work.
          </li>
          <li>
            <strong>A community plugin.</strong> Yours, in your repository, listed in the index.
            See <Link href="/plugins">the marketplace</Link> for how to submit one.
          </li>
        </ul>

        <h2>Getting set up</h2>
        <p>
          The CLI is one Go binary with no code generation step and no vendored dependencies, so
          this really is the whole setup:
        </p>
        <CodeBlock
          terminal
          code={`git clone https://github.com/MUKE-coder/grit.git
cd grit
go build ./...
go test ./...`}
        />
        <p>To try a change against a real project:</p>
        <CodeBlock
          terminal
          code={`go build -o /tmp/grit ./cmd/grit
cd /tmp && /tmp/grit new scratch --triple --next --db sqlite
cd scratch/apps/api && go build ./...`}
        />
        <p>
          For a documentation change you need none of that. The site is in{' '}
          <code>docs/</code>:
        </p>
        <CodeBlock terminal code={`cd docs && npm run dev`} />

        <h2>The one thing that catches everybody</h2>
        <p>
          <strong>
            A template is a Go function returning a string, not a file in a template directory.
          </strong>{' '}
          Everything a generated project contains lives in <code>internal/scaffold</code> as Go
          source. That is deliberate: the templates compile, so a typo is a build failure here
          rather than a broken project for somebody else.
        </p>
        <p>
          The consequence for you is that changing what a generated project looks like means
          editing Go, and that a change to a template needs to reach two places: the map that
          writes a new project, and the writer that an upgrade uses. There is a short note about
          that in <code>internal/scaffold/AGENTS.md</code>, which is worth reading before your
          first scaffold change.
        </p>

        <h2>Sending it</h2>
        <ul>
          <li>
            <strong>Branch, commit, open a pull request.</strong> Conventional commit messages:{' '}
            <code>fix(scaffold): ...</code>, <code>docs: ...</code>, <code>feat(generate): ...</code>
          </li>
          <li>
            <strong>Run <code>go test ./...</code> first.</strong> The suite is fast and it
            covers the generated output, not only the generator.
          </li>
          <li>
            <strong>Say what you verified.</strong> &quot;I generated a project and the admin
            screen renders&quot; is worth more in a pull request than a description of the code,
            which the reviewer can read.
          </li>
          <li>
            <strong>Small is better.</strong> A thirty-line change gets read today. A
            thousand-line one waits for an afternoon that may not come this week.
          </li>
        </ul>

        <h2>What gets refused, so you do not waste an evening</h2>
        <ul>
          <li>
            <strong>A new dependency without a reason that is not convenience.</strong> A
            generated project&apos;s dependency list is a thing its users have to audit.
          </li>
          <li>
            <strong>A rewrite of something that works</strong> to a different style. Known
            behaviour traded for unknown behaviour is a bad trade even when the new code is
            nicer.
          </li>
          <li>
            <strong>A second way to do something that already has one.</strong> Two commands
            that overlap is a decision every user then has to make.
          </li>
          <li>
            <strong>Em dashes in anything a reader sees.</strong> There is a checker, and it
            will fail your build. Use a colon, a comma or a full stop.
          </li>
        </ul>

        <h2>If you are not sure</h2>
        <p>
          Open an issue before writing the code. &quot;I was thinking of doing X, is that
          wanted?&quot; costs you a sentence and can save you a weekend. For anything larger than
          a page or a bug fix, that is the recommended order rather than a formality.
        </p>
        <p>
          The full conventions, the repository layout and the test guidance are in{' '}
          <a
            href="https://github.com/MUKE-coder/grit/blob/main/CONTRIBUTING.md"
            target="_blank"
            rel="noreferrer"
          >
            CONTRIBUTING.md
          </a>
          . The people who have already contributed are on{' '}
          <Link href="/contributors">the contributors page</Link>.
        </p>
      </div>
    </DocsShell>
  )
}
