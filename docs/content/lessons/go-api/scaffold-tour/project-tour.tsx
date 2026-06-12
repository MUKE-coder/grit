import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Let&apos;s walk every folder inside <code>apps/api/</code>. By the end
        you&apos;ll know what each file does and where to put new code without
        thinking.
      </p>

      <h2>The full tree</h2>
      <CodeBlock
        language="text"
        filename="apps/api/"
        code={`apps/api/
├── cmd/
│   └── server/
│       └── main.go             ← entry point, wires services + starts Gin
├── internal/
│   ├── config/                 loads .env into a typed Config struct
│   ├── database/               Postgres / SQLite connection + AutoMigrate
│   ├── handlers/               HTTP handlers — thin
│   ├── middleware/             auth, CORS, security headers, request ID
│   ├── models/                 GORM structs (User, Upload, …)
│   ├── routes/
│   │   └── routes.go           mounts every handler under /api/*
│   ├── services/               business logic — thick
│   ├── ai/                     AI Gateway client
│   ├── cache/                  Redis cache
│   ├── jobs/                   Asynq queue + workers
│   ├── mail/                   Resend client + templates
│   ├── storage/                S3 / MinIO client
│   └── totp/                   TOTP 2FA (encrypted seed storage)
├── go.mod
├── .air.toml                   hot-reload config
└── Dockerfile`}
      />

      <h2>cmd/server/main.go — the entry point</h2>
      <p>This is where Gin starts. The full flow:</p>
      <CodeBlock
        language="go"
        filename="apps/api/cmd/server/main.go (abridged)"
        code={`func main() {
    cfg := config.Load()
    db := database.MustConnect(cfg)

    // Instantiate every service the handlers need
    svc := &routes.Services{
        Cache:   cache.New(cfg),
        Storage: storage.New(cfg),
        Mailer:  mail.New(cfg),
    }

    r := routes.Setup(db, cfg, svc)  // builds the Gin router
    r.Run(":" + cfg.Port)
}`}
      />
      <p>
        Roughly 30 lines. <code>main.go</code> is intentionally tiny — its
        job is to wire dependencies, not to contain logic.
      </p>

      <h2>The Services pattern</h2>
      <p>
        Every external dependency (Redis, S3, Resend, AI Gateway) is wrapped
        in a tiny service. Those services compose into a single{' '}
        <code>routes.Services</code> struct that&apos;s passed to{' '}
        <code>routes.Setup()</code>. Handlers get the services they need by
        reference; nothing is global.
      </p>
      <TipBox tone="success">
        <strong>Why this matters:</strong> tests can pass a mock{' '}
        <code>Services</code> struct and exercise handlers without booting
        Redis / S3 / Resend. The handler doesn&apos;t care if the storage
        backend is real S3 or a fake in-memory one.
      </TipBox>

      <h2>internal/ — recap from Concepts</h2>
      <p>
        You met this folder in Grit Concepts ch.3. Quick refresher of the
        hot path:
      </p>
      <CodeBlock
        language="text"
        code={`request → middleware → routes.go picks handler → handler calls service
                                                  → service calls model
                                                  → response shaped by handler`}
      />

      <h2>The new folders for batteries</h2>
      <ul>
        <li>
          <code>jobs/</code> — Asynq integration. Enqueue from anywhere; the
          worker (separate process) consumes them.
        </li>
        <li>
          <code>mail/</code> — Resend client + HTML templates. In dev, all
          mail goes to Mailhog.
        </li>
        <li>
          <code>storage/</code> — S3-compatible (R2, MinIO, B2). One interface,
          swap the backend via env.
        </li>
        <li>
          <code>ai/</code> — AI Gateway client. Stream from Claude, OpenAI,
          and ~98 other models.
        </li>
        <li>
          <code>totp/</code> — 2FA: TOTP seeds encrypted at rest via AES-GCM.
        </li>
      </ul>

      <KnowledgeCheck
        question="You want to add a Stripe webhook handler. Which two files do you touch first?"
        choices={[
          {
            label: 'apps/api/internal/handlers/stripe.go + apps/api/internal/routes/routes.go',
            correct: true,
            feedback:
              "Right — the handler in handlers/ (logic split into a service if it grows), and a route registration in routes.go. That's the canonical pair for any new HTTP endpoint.",
          },
          {
            label: 'apps/api/cmd/server/main.go + apps/api/internal/middleware/stripe.go',
            feedback:
              "Wrong — main.go is for dependency wiring, not endpoint registration. And middleware is for cross-cutting concerns (auth, CORS), not webhook handling.",
          },
          {
            label: 'A new internal/stripe/ package + main.go',
            feedback:
              "Half right — yes you might add internal/stripe/ for the SDK wrapper. But the HTTP-facing piece is a handler, and the route still has to be registered in routes.go.",
          },
          {
            label: 'apps/api/internal/models/stripe.go and that\'s it',
            feedback:
              "Wrong — a model alone doesn't make an endpoint. You need a handler to receive the webhook + a route to mount it.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Open <code>apps/api/internal/routes/routes.go</code>. Count
              how many route groups are registered (look for{' '}
              <code>r.Group(...)</code> or <code>api.Group(...)</code>). For
              each, write one line in <code>notes.md</code> describing what
              it serves.
            </p>
          </>
        }
        hint={
          <>
            You&apos;ll see a public group (no auth), an auth group{' '}
            (<code>/api/auth</code>), and a protected group (auth middleware
            applied). Some kits also have an admin group.
          </>
        }
        solution={
          <>
            <p>You should find something like:</p>
            <ul>
              <li><code>/api/health</code> — public, returns the version envelope</li>
              <li><code>/api/auth/*</code> — login / register / refresh — public</li>
              <li><code>/api/users/*</code> — protected, requires JWT</li>
              <li><code>/api/uploads/*</code> — protected, file upload + serve</li>
            </ul>
            <p>
              Each group is registered with its own middleware stack — that&apos;s
              how Grit applies auth selectively without per-handler boilerplate.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        You know the layout. Last lesson of this chapter — start everything
        up, make a request, prove the API works.
      </p>
    </>
  )
}
