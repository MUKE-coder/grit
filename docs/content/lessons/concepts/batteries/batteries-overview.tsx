import { TipBox } from '@/components/course/tip-box'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        Grit ships with five &quot;batteries&quot; — pre-wired services
        that solve real production problems. Cache, file storage,
        email, background jobs, AI. You don&apos;t have to choose,
        install, or configure them from scratch — they&apos;re there,
        ready, in the right shape. This chapter walks through each one
        in turn.
      </p>

      <h2>What is a &quot;battery&quot;?</h2>
      <p>
        A battery is an included service that:
      </p>
      <ul>
        <li>
          Is enabled in <code>docker-compose.yml</code> for local dev.
        </li>
        <li>
          Has a Go client / service initialized in{' '}
          <code>main.go</code>.
        </li>
        <li>
          Is exposed to your code via a clean interface (e.g.,{' '}
          <code>s.cache.Get(key)</code>, <code>s.mail.Send(to, ...)</code>).
        </li>
        <li>
          Has an admin page in the dashboard so non-devs can poke at it
          (Jobs dashboard, Mail preview, etc.).
        </li>
      </ul>

      <h2>The five batteries</h2>
      <Diagram label="Where the batteries plug in" caption="Each battery has a dedicated service in the API. Handlers don't talk to Redis or Resend directly — they go through the service interface.">
{`   Your handler
        │
        ▼
   Your service ──┬────► CACHE      (Redis)             — speed up reads, throttle
                  │
                  ├────► STORAGE    (S3 / R2 / MinIO)   — file uploads
                  │
                  ├────► MAIL       (Resend)            — transactional + marketing email
                  │
                  ├────► JOBS       (asynq + Redis)     — background work, retries, cron
                  │
                  └────► AI         (Claude / OpenAI)   — chat, embeddings, generation`}
      </Diagram>

      <h2>Why batteries-included matters</h2>
      <p>
        Without batteries, the &quot;day-1&quot; experience of a new
        API looks like this:
      </p>
      <ol>
        <li>Pick a Redis client library (3 choices, vibes-based decision).</li>
        <li>Wire connection pool, retry, env config.</li>
        <li>Decide on cache key naming convention.</li>
        <li>Wrap GET/SET in a service so you can test it.</li>
        <li>… and you haven&apos;t written any product code yet.</li>
      </ol>
      <p>
        Multiply that by five (S3, mail, jobs, AI). Grit makes the
        decisions for you — opinionated, but the opinions are
        defensible. You can replace any battery later, but you
        don&apos;t have to choose to get started.
      </p>

      <h2>Where the batteries live in the code</h2>
      <pre className="not-prose text-xs leading-relaxed bg-bg-elevated border border-border rounded-lg p-4 overflow-x-auto">{`apps/api/internal/
├── cache/         ← Redis cache service + middleware
├── storage/       ← S3/MinIO/R2 client + upload handler
├── mail/          ← Resend client + 4 HTML templates
├── jobs/          ← asynq client + worker definitions
└── ai/            ← Claude + OpenAI clients`}</pre>
      <p>
        Same shape: one package per battery, exposing a thin Go
        interface. Tests can swap in a fake; production wires the real
        thing.
      </p>

      <h2>Enabling / disabling at scaffold time</h2>
      <p>
        When you run <code>grit new</code>, the CLI asks which
        batteries you want. You can include everything (default for{' '}
        <code>triple</code> kit) or skip some for a minimal API. You
        can also add a battery later by running{' '}
        <code>grit add battery &lt;name&gt;</code>.
      </p>

      <h2>The admin dashboard for each battery</h2>
      <p>
        Every battery gets a page in the admin panel under{' '}
        <code>/admin/system</code>:
      </p>
      <ul>
        <li>
          <strong>Jobs</strong> — live queue depth, running jobs, retry
          attempts, failed jobs you can re-queue.
        </li>
        <li>
          <strong>Mail Preview</strong> — render templates with sample
          data without sending.
        </li>
        <li>
          <strong>Files</strong> — browse + delete uploads.
        </li>
        <li>
          <strong>Cron</strong> — see scheduled jobs, last/next run.
        </li>
        <li>
          <strong>AI</strong> — usage stats, token spend.
        </li>
      </ul>

      <TipBox tone="info">
        <strong>Pay-as-you-grow.</strong> Local dev uses MinIO (S3
        clone), a local Redis, and Mailhog (no real emails sent).
        Production swaps in real providers via env vars. Same code,
        different config. We&apos;ll show this for each battery in
        the next five lessons.
      </TipBox>

      <h2>What this chapter covers, in order</h2>
      <ol>
        <li>
          <strong>Redis cache</strong> — speed up reads, throttle, store
          ephemeral state.
        </li>
        <li>
          <strong>S3 storage</strong> — uploads, signed URLs, image
          processing.
        </li>
        <li>
          <strong>Mail (Resend)</strong> — transactional + marketing
          email with editable templates.
        </li>
        <li>
          <strong>Jobs (asynq)</strong> — background work, retries,
          cron.
        </li>
        <li>
          <strong>AI (Claude + OpenAI)</strong> — chat, embeddings,
          provider swap.
        </li>
      </ol>
      <p>
        Each lesson follows the same template: what it does, how it&apos;s
        implemented, where the file lives, how to modify it, and a
        small exercise.
      </p>

      <KnowledgeCheck
        question="A teammate asks &quot;why not just install Redis ourselves and write the client?&quot;. What's the best one-line response?"
        choices={[
          {
            label: 'Redis is too complex to install',
            feedback:
              "Redis install is one apt-get away. The reason is consistency, not difficulty.",
          },
          {
            label: 'The shape (interface, env config, admin page) being identical across projects means anyone joining can navigate from day one — that compounds over a team',
            correct: true,
            feedback:
              "Right — the value isn't &quot;we couldn't install Redis&quot;; it's &quot;every Grit project has the cache wired the same way, so any dev who knows one knows them all&quot;.",
          },
          {
            label: "Grit's Redis is more secure",
            feedback:
              "It's the same Redis. The wrapper isn't security-related.",
          },
          {
            label: "Grit forces you to use Redis",
            feedback:
              "It doesn't — you can skip the cache battery and write to whatever you want. The default is just opinionated.",
          },
        ]}
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — <strong>Redis cache</strong>. Speed up reads,
        cache list endpoints, throttle expensive calls. The most
        common production performance win.
      </p>
    </>
  )
}
