import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Time to bring the API up. Four commands. By the end you&apos;ll have a
        live JSON response in your terminal.
      </p>

      <h2>1. Start Postgres + friends</h2>
      <CodeBlock terminal code={`cd bench-api
docker compose up -d`} />
      <p>
        Brings up Postgres, Redis, MinIO, and Mailhog in the background. The{' '}
        <code>-d</code> means &quot;detached&quot;; check status with{' '}
        <code>docker compose ps</code>.
      </p>

      <h2>2. Run the API from the project root</h2>
      <p>
        <strong>Important:</strong> run from project root, not from inside{' '}
        <code>apps/api/cmd/server/</code>. The config loader looks for{' '}
        <code>./.env</code> relative to your working directory, and{' '}
        <code>.env</code> lives at project root.
      </p>
      <CodeBlock terminal code={`# from bench-api/
cd apps/api && go run ./cmd/server`} />
      <p>
        Wait — that&apos;s also a non-root cwd. The trick: run from{' '}
        <code>apps/api/</code> (because that&apos;s where <code>go.mod</code>{' '}
        lives), and the loader will find the project-root .env via the
        <code> ../../.env</code> path it tries.
      </p>

      <TipBox tone="warning">
        Don&apos;t <code>cd apps/api/cmd/server</code> first — from there, the
        loader can&apos;t find <code>.env</code> and the API exits with{' '}
        <code>DATABASE_URL is required</code>. This is the gotcha most new
        users hit. Always run from <code>apps/api/</code>.
      </TipBox>

      <h2>3. Watch the startup log</h2>
      <CodeBlock
        language="text"
        code={`2026/06/01 09:00:00 Database connected successfully
2026/06/01 09:00:00 Redis cache connected
2026/06/01 09:00:00 File storage connected
2026/06/01 09:00:00 Job queue connected
2026/06/01 09:00:00 [sentinel] Mounted at /sentinel
2026/06/01 09:00:00 Pulse observability mounted at /pulse
2026/06/01 09:00:00 [GIN-debug] Listening and serving HTTP on :8080`}
      />
      <p>
        Every battery wired in five lines. If any of these say &quot;failed&quot;,
        scroll up — the error message tells you exactly what&apos;s missing
        (usually a env value or Docker service that didn&apos;t come up).
      </p>

      <h2>4. Hit /api/health</h2>
      <CodeBlock terminal code={`curl http://localhost:8080/api/health`} />
      <p>You should see:</p>
      <CodeBlock language="json" code={`{ "status": "ok", "version": "0.1.0" }`} />
      <p>
        That&apos;s success. The API connected to every service it needs and is
        ready to handle requests.
      </p>

      <h2>The other URLs you can poke</h2>
      <ul>
        <li>
          <code>http://localhost:8080/docs</code> — auto-generated OpenAPI
          documentation
        </li>
        <li>
          <code>http://localhost:8080/studio</code> — GORM Studio (DB browser)
        </li>
        <li>
          <code>http://localhost:8080/sentinel/ui</code> — Sentinel WAF
          dashboard
        </li>
        <li>
          <code>http://localhost:8080/pulse/ui</code> — Pulse observability
        </li>
        <li>
          <code>http://localhost:8025</code> — Mailhog (captures all dev
          email)
        </li>
      </ul>

      <KnowledgeCheck
        question="You see `DATABASE_URL is required` on startup. What's the most likely cause?"
        choices={[
          {
            label: 'Postgres isn\'t running',
            feedback:
              "Wrong — `DATABASE_URL is required` fires from the config loader BEFORE we try to connect. The DB might be down, but that's not what this error says.",
          },
          {
            label: 'You ran `go run .` from `apps/api/cmd/server/` instead of `apps/api/`',
            correct: true,
            feedback:
              "Right — the godotenv loader looks for .env relative to cwd. From cmd/server it can't find it; from apps/api it follows ../../.env and finds the project-root file.",
          },
          {
            label: 'Your .env file is malformed',
            feedback:
              "Possible but rare — a malformed .env usually errors with a parse complaint, not 'required'.",
          },
          {
            label: 'You forgot to scaffold with --api',
            feedback:
              "Wrong — the API kit still scaffolds the .env at project root. Kit choice doesn't change where .env lives.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Get the API running end-to-end. From your scaffolded project:
            </p>
            <ol>
              <li><code>docker compose up -d</code></li>
              <li><code>cd apps/api &amp;&amp; go run ./cmd/server</code></li>
              <li>In another terminal: <code>curl http://localhost:8080/api/health</code></li>
            </ol>
            <p>
              Paste the curl response in <code>notes.md</code>. That completes
              chapter 1&apos;s assignment — the criteria are in the sidebar.
            </p>
          </>
        }
        hint={
          <>
            If you get a port conflict on 8080, check what else is running
            (<code>lsof -i :8080</code> on Mac/Linux) or change{' '}
            <code>APP_PORT</code> in <code>.env</code>.
          </>
        }
        solution={
          <>
            <p>Final response:</p>
            <CodeBlock
              language="json"
              code={`{ "status": "ok", "version": "0.1.0" }`}
            />
            <p>
              Congrats — you have a Grit Go API running locally. Chapter 2
              opens with modelling the data layer.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 2 — <strong>Modeling with GORM</strong>. You&apos;ll define your
        first model, see how AutoMigrate handles tables, and learn the GORM
        tags Grit assumes.
      </p>
    </>
  )
}
