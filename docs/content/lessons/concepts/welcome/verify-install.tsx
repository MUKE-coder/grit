import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Three quick commands and you&apos;re ready for chapter 2. This lesson is
        the &quot;does my install actually work&quot; sanity check, plus how to keep
        grit up to date once you start using it daily.
      </p>

      <h2>1. Version check</h2>
      <CodeBlock
        terminal
        code={`grit version`}
      />
      <p>
        Should print <code>grit version 3.25.x</code> or later. If it says{' '}
        <code>3.24.x</code> or older, run <code>grit update</code> first —
        SQLite support, scaffolded random secrets, and the one-line install
        script all landed in 3.25.
      </p>

      <h2>2. Help — see what grit can do</h2>
      <CodeBlock
        terminal
        code={`grit --help`}
      />
      <p>You should see a printed list including (at minimum):</p>

      <CodeBlock
        language="text"
        code={`Available Commands:
  new            Scaffold a new Grit project
  new-desktop    Scaffold a Wails desktop app
  generate       Generate a resource (model + handler + service + types + ...)
  sync           Sync Go types → TypeScript
  migrate        Run database migrations
  seed           Seed the database
  start          Start dev servers
  update         Update the Grit CLI to the latest version
  deploy         Deploy to a VPS
  version        Print the grit CLI version`}
      />

      <TipBox tone="info">
        Every command also takes <code>--help</code> for its own flags. Try{' '}
        <code>grit new --help</code> — you&apos;ll see the architecture choices
        we&apos;ll cover in chapter 5.
      </TipBox>

      <h2>3. Update check — the command you&apos;ll use most</h2>
      <CodeBlock
        terminal
        code={`grit update`}
      />
      <p>
        Hits GitHub&apos;s release API, compares the latest tag to the version of
        the binary running this command, and either updates in place or
        exits with &quot;Already on the latest version. Nothing to do.&quot;
      </p>

      <p>
        Make a habit of running this once a week. New features ship
        regularly; <code>grit update</code> is a single HTTP round-trip when
        there&apos;s nothing to do, so the cost is essentially zero.
      </p>

      <h2>What if something goes wrong?</h2>
      <p>The two most common issues:</p>
      <ul>
        <li>
          <strong>&quot;command not found&quot;</strong> — your PATH didn&apos;t pick up the
          install dir. Open a new terminal first. If that doesn&apos;t fix it,
          re-read the install lesson&apos;s exercise solution.
        </li>
        <li>
          <strong>&quot;permission denied&quot;</strong> on macOS / Linux when{' '}
          <code>grit update</code> runs — the install dir requires sudo. The
          install script&apos;s default location (<code>$HOME/.local/bin</code>) is
          user-owned, so you should never hit this. If you do, the script
          probably installed to <code>/usr/local/bin</code>; re-run with{' '}
          <code>sudo</code> or reinstall to <code>$HOME/.local/bin</code>.
        </li>
      </ul>

      <KnowledgeCheck
        question="You ran `grit update` and it says 'Already on the latest version. Nothing to do.' What just happened?"
        choices={[
          {
            label: 'grit did nothing — it failed silently.',
            feedback:
              "Wrong. The 'Nothing to do' is the success path, not silent failure. grit hit GitHub, compared versions, and exited cleanly.",
          },
          {
            label: 'grit downloaded the binary anyway and overwrote your install.',
            feedback:
              "Wrong — the whole point of the version check is to skip the download when you're current. v3.26.1 added this short-circuit specifically to avoid wasted go install runs.",
          },
          {
            label: 'grit made a single HTTP request to GitHub, found you were current, and exited.',
            correct: true,
            feedback:
              "Right — one round-trip to api.github.com/repos/MUKE-coder/grit/releases/latest, version compare, exit. The whole thing takes <500ms typically.",
          },
          {
            label: 'grit cached the result locally and will skip the check next time.',
            feedback:
              "Wrong — every `grit update` hits GitHub fresh. There's no local cache.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Run all three sanity checks and capture proof:</p>
            <ol>
              <li>
                Run <code>grit version</code> and paste output in{' '}
                <code>notes.md</code>.
              </li>
              <li>
                Run <code>grit --help</code> and copy the list of commands
                into <code>notes.md</code>.
              </li>
              <li>
                Run <code>grit update</code> and paste its output too — if it
                says &quot;Already on the latest version&quot;, that&apos;s the success
                state.
              </li>
            </ol>
            <p>
              You&apos;ve now completed chapter 1&apos;s assignment. Look at the
              sidebar — <strong>ch.1 Assignment: Install Grit and capture
              proof</strong> is the next link. Click it to see the success
              criteria and tick them off.
            </p>
          </>
        }
        hint={
          <>
            If <code>grit update</code> tries to actually pull a new version,
            let it finish — then re-run <code>grit version</code> to confirm
            the new version is in place.
          </>
        }
        solution={
          <>
            <p>Your <code>notes.md</code> should now look something like:</p>
            <CodeBlock
              language="markdown"
              filename="notes.md"
              code={`# Grit Concepts notes

## Chapter 1 — install

\`\`\`
$ grit version
grit version 3.26.1
\`\`\`

\`\`\`
$ grit --help
Available Commands:
  new            Scaffold a new Grit project
  new-desktop    Scaffold a Wails desktop app
  ...
\`\`\`

\`\`\`
$ grit update
Checking GitHub for the latest release...
✓ Already on the latest version (v3.26.1). Nothing to do.
\`\`\``}
            />
            <p>
              That&apos;s the assignment done. Open the chapter assignment page
              and tick off the criteria.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 2 — <strong>your first real Grit project</strong>. We&apos;ll
        scaffold a project, tour every folder it produces, start the dev
        servers, and log into the admin panel. By the end of the chapter
        you&apos;ll have a working full-stack app in your browser.
      </p>
    </>
  )
}
