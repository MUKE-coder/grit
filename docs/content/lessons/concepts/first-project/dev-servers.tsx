import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { PrereqLinks } from '@/components/course/prereq-links'

export default function Lesson() {
  return (
    <>
      <p>
        Time to start everything. By the end of this lesson, Postgres, the
        API, the web app, and the admin panel will all be running together
        and you&apos;ll have URLs to open in your browser.
      </p>

      <PrereqLinks
        prereqs={['docker']}
        intro={
          <>
            We&apos;ll use Docker Compose to boot infra (Postgres, Redis,
            MinIO, Mailhog). If you&apos;ve never used containers or
            compose files, the Docker primer is a 10-minute read that
            makes the rest of this lesson click.
          </>
        }
      />

      <h2>1. Start the infrastructure</h2>
      <CodeBlock
        terminal
        code={`cd my-first-grit
docker compose up -d`}
      />
      <p>
        Brings up Postgres (host 5434), Redis (host 6380), MinIO (host
        9002 + 9003), and Mailhog (1025 + 8025) in the background — the
        host ports are intentionally offset from the canonical defaults
        so a native install of Postgres / Redis / MinIO on your machine
        doesn&apos;t collide. Wait ~5 seconds for them to be ready.
      </p>

      <TipBox tone="warning">
        If you don&apos;t have Docker installed, install Docker Desktop first.
        Grit&apos;s dev loop assumes docker compose works. Production deploys
        don&apos;t need Docker — but development does.
      </TipBox>

      <h2>2. Install JS deps</h2>
      <CodeBlock terminal code={`pnpm install`} />
      <p>
        Installs everything for web, admin, and the shared package — all in
        one pnpm workspace.
      </p>

      <h2>3. Migrate + seed</h2>
      <CodeBlock
        terminal
        code={`grit migrate    # creates the DB tables
grit seed       # inserts the admin user + sample data`}
      />
      <p>
        The seed creates an <code>admin@example.com</code> user with
        password <code>admin123</code>. We&apos;ll log in with that in the next
        lesson.
      </p>

      <h2>4. Start the dev servers</h2>
      <CodeBlock terminal code={`grit start`} />
      <p>One command, three apps:</p>
      <CodeBlock
        language="text"
        code={`✓ API     :8080    apps/api
✓ Web     :3000    apps/web
✓ Admin   :3001    apps/admin`}
      />
      <p>Hot reload is on for all three. Edit a file, the page reloads.</p>

      <KnowledgeCheck
        question="You forgot to run `grit migrate` before `grit start`. What happens?"
        choices={[
          {
            label: 'The API starts but every request errors with `relation "users" does not exist`.',
            correct: true,
            feedback:
              'Right — the API itself starts (it connects to Postgres successfully) but the tables don\'t exist yet. Run `grit migrate` and you\'re unblocked.',
          },
          {
            label: 'grit start refuses to run and tells you to migrate first.',
            feedback:
              'Wrong — there\'s no pre-flight check. The API will try to query a table that doesn\'t exist and fail at request time.',
          },
          {
            label: 'grit start auto-migrates for you.',
            feedback:
              'Wrong — migrations are explicit. This is intentional: in production you don\'t want the API auto-migrating on boot.',
          },
          {
            label: 'Docker Postgres refuses to accept connections.',
            feedback:
              'Wrong — Postgres is connecting just fine; the tables inside the DB simply don\'t exist yet.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Run the full boot sequence on your machine. When all four steps
              have run successfully, paste the <code>grit start</code> output
              into <code>notes.md</code>.
            </p>
          </>
        }
        hint={
          <>
            If any step fails, check that Docker is running (
            <code>docker ps</code>) and that no port (5434, 6380, 9002, 9003,
            8025, 8080, 3000, 3001) is already in use by another app.
          </>
        }
        solution={
          <>
            <p>You should see the three checkmarks:</p>
            <CodeBlock
              language="text"
              code={`✓ API     :8080
✓ Web     :3000
✓ Admin   :3001`}
            />
            <p>
              If one of them is red instead of green, scroll up — the
              error message tells you what&apos;s wrong (usually a port
              conflict or a missing dep).
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Everything&apos;s running. Next lesson — open the URLs, log into the
        admin panel, and click around so you understand what each surface
        looks like.
      </p>
    </>
  )
}
