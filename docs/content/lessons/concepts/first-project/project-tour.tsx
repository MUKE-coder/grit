import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        You have a project. Now let&apos;s walk through every top-level folder so
        you know what goes where before you write a single line of code.
      </p>

      <h2>The 30-second map</h2>
      <CodeBlock
        language="text"
        filename="my-first-grit/"
        code={`my-first-grit/
├── apps/
│   ├── api/             # Go API (Gin + GORM)
│   ├── web/             # Next.js public site
│   └── admin/           # Next.js admin panel
├── packages/
│   └── shared/          # Zod schemas + TS types (used by web + admin)
├── docker-compose.yml   # Postgres, Redis, MinIO, Mailhog
├── docker-compose.prod.yml
├── grit.json            # Project metadata grit reads
├── package.json         # Root pnpm workspace
├── pnpm-workspace.yaml  # Workspace member list
├── turbo.json           # Turborepo task pipeline
├── .env                 # Secrets + URLs (gitignored)
├── .env.example         # Template for new collaborators
└── README.md            # Project-specific readme`}
      />

      <h2>apps/api — the Go backend</h2>
      <p>This is the engine. Open it up:</p>
      <CodeBlock
        language="text"
        filename="apps/api/"
        code={`apps/api/
├── cmd/server/main.go   # Entry point — where the API starts
├── internal/
│   ├── config/          # Loads .env into a typed Config struct
│   ├── database/        # Postgres connection + AutoMigrate
│   ├── handlers/        # HTTP handlers (Gin)
│   ├── middleware/      # Auth, CORS, security, etc.
│   ├── models/          # GORM structs
│   ├── routes/routes.go # Mounts every handler on a route
│   └── services/        # Business logic, called by handlers
├── .air.toml            # Hot-reload config (if you use air)
├── Dockerfile
└── go.mod`}
      />
      <p>
        The pattern is <strong>handler → service → model</strong>. Handlers
        are thin (parse input, return response). Services contain the
        business logic. Models are the data layer. We&apos;ll cover this in
        depth in chapter 3.
      </p>

      <h2>apps/web — the public Next.js app</h2>
      <p>
        Next.js 14 App Router. Standard layout: <code>app/</code> for routes,{' '}
        <code>components/</code> for reusable components,{' '}
        <code>lib/</code> for the API client. Nothing surprising if you&apos;ve
        used Next.js before.
      </p>

      <h2>apps/admin — the Filament-style admin panel</h2>
      <p>
        Also Next.js, but configured as an admin dashboard. The killer file
        is <code>app/resources/</code> — every model gets a <code>defineResource()</code>{' '}
        call that auto-builds a CRUD page. We&apos;ll see this in chapter 4.
      </p>

      <h2>packages/shared — the type bridge</h2>
      <p>
        Zod schemas + TypeScript types that web and admin both import. When
        you run <code>grit sync</code> (chapter 4), Go structs get translated
        to TS types here — so the types in your React code always match the
        Go reality.
      </p>

      <TipBox tone="info">
        The folder structure is the same for every Grit project. Once you
        learn it, you know your way around any Grit codebase — yours,
        someone else&apos;s, an open-source one.
      </TipBox>

      <h2>The root config files</h2>
      <ul>
        <li>
          <code>grit.json</code> — tells grit this is a Grit project. Includes
          the project name + architecture mode for diagnostics.
        </li>
        <li>
          <code>docker-compose.yml</code> — the dev infrastructure (Postgres,
          Redis, MinIO for S3-compatible storage, Mailhog for email testing).
        </li>
        <li>
          <code>turbo.json</code> — Turborepo task definitions for{' '}
          <code>build</code>, <code>dev</code>, <code>test</code>.
        </li>
        <li>
          <code>.env</code> — generated with crypto-random secrets (you&apos;ve
          seen these in past lessons).
        </li>
      </ul>

      <KnowledgeCheck
        question="A teammate asks: 'Where do I put the business logic for handling refunds?' What's your answer?"
        choices={[
          {
            label: 'apps/api/internal/handlers/refund_handler.go',
            feedback:
              "Wrong — handlers are intentionally thin. They just parse input and return responses. Business logic belongs in services.",
          },
          {
            label: 'apps/api/internal/services/refund_service.go',
            correct: true,
            feedback:
              "Right — services own the business logic. The handler calls the service, the service does the work, the handler shapes the response.",
          },
          {
            label: 'apps/api/internal/models/refund.go',
            feedback:
              "Wrong — models are data shape only (GORM struct + tags). No business logic.",
          },
          {
            label: 'apps/web/lib/refund.ts',
            feedback:
              "Wrong — apps/web is the public Next.js frontend. Business logic on the backend stays on the backend.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Open <code>apps/api/internal/handlers/</code>. List the files
              you see. Pick one, open it, and answer in <code>notes.md</code>:
            </p>
            <ul>
              <li>What model does it handle?</li>
              <li>How many functions does it expose?</li>
              <li>Does it import a service?</li>
            </ul>
          </>
        }
        hint={
          <>
            A fresh triple project ships with a <code>user.go</code> and an{' '}
            <code>upload.go</code> handler. Either one is a fine pick. Look
            for <code>&quot;&lt;MODULE&gt;/internal/services&quot;</code> in the imports.
          </>
        }
        solution={
          <>
            <p>For <code>user.go</code> you should find:</p>
            <ul>
              <li>Handles the <code>User</code> model</li>
              <li>5 functions: List, Get, Create, Update, Delete</li>
              <li>
                Yes — imports <code>internal/services</code>. The handler is
                tiny (~30 lines per function); the service is where the
                actual work lives.
              </li>
            </ul>
            <p>That&apos;s the handler/service split in action.</p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Tour done. Next lesson — fire up the dev servers. You&apos;ll have
        Postgres, Redis, the API, web, and admin all running together in
        under 5 minutes.
      </p>
    </>
  )
}
