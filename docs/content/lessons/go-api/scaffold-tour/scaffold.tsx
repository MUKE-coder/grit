import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Welcome to <strong>Building a Go API</strong>. Goal of this chapter:
        scaffold an API-only project, identify every file, get it serving
        requests. This first lesson is the scaffold itself — 5 minutes, one
        command, you&apos;re ready.
      </p>

      <h2>The command</h2>
      <CodeBlock
        terminal
        code={`grit new bench-api --api`}
      />
      <p>
        Pick the project name (anything kebab-case works). The{' '}
        <code>--api</code> flag tells Grit you want the headless Go backend
        kit — no <code>apps/web</code>, no <code>apps/admin</code>, just the
        Go API.
      </p>

      <h2>What you get</h2>
      <p>You&apos;ll see the standard checklist run:</p>
      <CodeBlock
        language="text"
        code={`✓ Created bench-api/
✓ Wrote 42 files
✓ Initialized go module
✓ Generated random secrets (JWT_SECRET, SENTINEL_*, PULSE_*)

Next steps:
    cd bench-api
    docker compose up -d
    cd apps/api && go run cmd/server/main.go`}
      />

      <TipBox tone="warning">
        <strong>Surprise:</strong> even the API-only kit is structured as a
        monorepo. The Go code lives at{' '}
        <code>apps/api/</code>, not the project root. This is so you can later
        add a frontend by dropping it next to <code>apps/api/</code> without
        moving files.
      </TipBox>

      <h2>The kit&apos;s shape at a glance</h2>
      <CodeBlock
        language="text"
        filename="bench-api/"
        code={`bench-api/
├── apps/
│   └── api/                    The Go API lives here
│       ├── cmd/server/main.go  Entry point
│       ├── internal/           All Go code
│       ├── .air.toml           Hot-reload config
│       ├── go.mod
│       └── Dockerfile
├── docker-compose.yml          Postgres + Redis + MinIO + Mailhog
├── docker-compose.prod.yml
├── .env                        Generated with random secrets
├── .env.example
├── grit.json
└── README.md`}
      />

      <h2>What didn&apos;t get scaffolded</h2>
      <ul>
        <li>No <code>apps/web/</code> — no public Next.js site</li>
        <li>No <code>apps/admin/</code> — no Filament-style admin panel</li>
        <li>No <code>packages/shared/</code> — no Zod schemas + TS types</li>
        <li>No <code>turbo.json</code>, <code>pnpm-workspace.yaml</code> — Go-only deps</li>
      </ul>
      <p>
        Everything else from the Concepts course — random JWT secret, Sentinel
        + Pulse pre-configured, SQLite support out of the box — still ships.
      </p>

      <KnowledgeCheck
        question="You scaffolded with `--api` but then decided you want an admin panel. What's the cleanest path?"
        choices={[
          {
            label: 'Re-scaffold with --triple and copy over your work',
            correct: true,
            feedback:
              "Right — the API-only kit doesn't ship the admin scaffolding, and bolting it on later is more work than re-scaffolding the right shape. This is exactly why ch.5 of Concepts emphasised picking your kit consciously.",
          },
          {
            label: 'Run `grit add admin` and let Grit add it',
            feedback:
              "Wrong — there's no incremental add-admin command. The kit choice is up-front.",
          },
          {
            label: 'Write your own Next.js admin in a sister folder',
            feedback:
              "You can — but you lose the defineResource() generator and have to wire everything yourself. The triple kit gives you a working admin in one command.",
          },
          {
            label: 'Add an apps/admin folder manually and import the API',
            feedback:
              "Way more work than re-scaffolding. The generated admin panel has 200+ lines of layout, sidebar, theme, resource registry — not worth hand-rolling.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Scaffold <code>bench-api</code> with the API-only kit. Then run{' '}
              <code>ls -la bench-api</code> (Mac/Linux) or <code>dir bench-api</code>{' '}
              (Windows) and paste the output in <code>notes.md</code>.
            </p>
            <p>
              Note: do <strong>not</strong> use <code>cd apps/api</code> yet
              — we&apos;ll tour the structure in the next lesson before running
              anything.
            </p>
          </>
        }
        hint={
          <>
            If you forgot the flag, just <code>grit new bench-api</code> and
            select &quot;API Only&quot; from the prompt.
          </>
        }
        solution={
          <>
            <p>You should see (Mac / Linux):</p>
            <CodeBlock
              terminal
              code={`$ ls -la bench-api
.env             .env.example     .gitignore       README.md
apps/            docker-compose.yml  docker-compose.prod.yml
grit.json`}
            />
            <p>
              No <code>package.json</code>, no <code>turbo.json</code> — that&apos;s
              the difference from triple. Tiny and focused.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Files are on disk. Next lesson — <strong>the tour</strong>. We open{' '}
        <code>apps/api/internal/</code> and walk every folder, naming the
        responsibility of each.
      </p>
    </>
  )
}
