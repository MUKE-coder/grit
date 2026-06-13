import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        Some work shouldn&apos;t happen inside a request: sending
        emails, resizing images, generating PDFs, syncing third-party
        data. Background jobs handle it asynchronously — your handler
        enqueues, returns instantly, a worker processes later. Grit
        ships{' '}
        <a href="https://github.com/hibiken/asynq" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">asynq</a> as the job engine.
      </p>

      <h2>Why we need it</h2>
      <p>
        Three reasons to push work to a job:
      </p>
      <ul>
        <li>
          <strong>Speed up requests.</strong> Sending email = 200ms; the
          user shouldn&apos;t wait. Enqueue, return in 5ms.
        </li>
        <li>
          <strong>Retries on failure.</strong> External API hiccupped?
          The job retries automatically with backoff. The handler would
          have just 500&apos;d.
        </li>
        <li>
          <strong>Scheduled work.</strong> Cron-like — nightly digest,
          hourly cleanup, weekly report. Without a job system, you
          spawn a cron container that&apos;s hard to monitor.
        </li>
      </ul>

      <h2>The two halves</h2>
      <Diagram label="Jobs architecture" caption="Producer and consumer separated by Redis. Workers run in a separate process from the API.">
{`   ┌─────────────────┐                          ┌─────────────────┐
   │  API process    │                          │  Worker process │
   │                 │                          │                 │
   │  handler        │                          │  process(job)   │
   │    │            │                          │    │            │
   │    │ enqueue    │                          │    │ run handler│
   │    ▼            │                          │    ▼            │
   │  ┌──────────┐   │   ┌──────────────────┐   │  ┌──────────┐   │
   │  │ asynq    │───┼──►│   Redis  queue   │◄──┼──│ asynq    │   │
   │  │ Client   │   │   │   (jobs:default) │   │  │ Worker   │   │
   │  └──────────┘   │   └──────────────────┘   │  └──────────┘   │
   └─────────────────┘                          └─────────────────┘`}
      </Diagram>
      <p>
        The API enqueues; the worker dequeues. They&apos;re separate
        Go processes. Both connect to Redis. You can scale workers
        independently — 10x more worker pods if your queue grows
        without touching the API.
      </p>

      <h2>Where it lives</h2>
      <pre className="not-prose text-xs leading-relaxed bg-bg-elevated border border-border rounded-lg p-4 overflow-x-auto">{`apps/api/internal/jobs/
├── client.go        ← Enqueue(taskName, payload)
├── handlers.go      ← Register handlers for each task name
├── server.go        ← Worker process bootstrap
└── cron.go          ← Cron schedule definitions
cmd/worker/main.go   ← Worker entry point (separate binary)`}</pre>

      <h2>How it&apos;s implemented</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/jobs/client.go"
        code={`package jobs

import (
  "context"
  "encoding/json"
  "github.com/hibiken/asynq"
)

type Client struct {
  inner *asynq.Client
}

func NewClient(redisURL string) *Client {
  opt, _ := asynq.ParseRedisURI(redisURL)
  return &Client{inner: asynq.NewClient(opt)}
}

func (c *Client) Enqueue(ctx context.Context, taskName string, payload any) error {
  b, err := json.Marshal(payload)
  if err != nil { return err }
  _, err = c.inner.EnqueueContext(ctx, asynq.NewTask(taskName, b))
  return err
}

func (c *Client) EnqueueIn(ctx context.Context, taskName string, payload any, delay time.Duration) error {
  b, _ := json.Marshal(payload)
  _, err := c.inner.EnqueueContext(ctx, asynq.NewTask(taskName, b), asynq.ProcessIn(delay))
  return err
}`}
      />
      <p>
        Two operations: enqueue now, enqueue in <code>delay</code>{' '}
        (great for &quot;send reminder in 24h&quot;).
      </p>

      <h2>Defining a job handler</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/jobs/handlers.go"
        code={`type Handlers struct {
  mail  *mail.Service
  store *storage.Service
}

func (h *Handlers) Register(mux *asynq.ServeMux) {
  mux.HandleFunc("send_welcome",     h.SendWelcome)
  mux.HandleFunc("resize_image",     h.ResizeImage)
  mux.HandleFunc("cleanup_uploads",  h.CleanupUploads)
}

func (h *Handlers) SendWelcome(ctx context.Context, t *asynq.Task) error {
  var p struct{ To, Name string }
  if err := json.Unmarshal(t.Payload(), &p); err != nil {
    return fmt.Errorf("bad payload: %w", err)
  }
  return h.mail.Send(ctx, p.To, "welcome.html", "Welcome!", p)
}`}
      />
      <p>
        Each task name maps to a handler function. Return{' '}
        <code>nil</code> for success; return an error to trigger retry.
        asynq re-runs the job up to N times with exponential backoff,
        then dead-letters it.
      </p>

      <h2>The worker process</h2>
      <CodeBlock
        language="go"
        filename="cmd/worker/main.go"
        code={`func main() {
  cfg := config.Load()
  db := db.MustOpen(cfg)
  services := services.New(db, cfg)

  srv := asynq.NewServer(
    asynq.RedisClientOpt{Addr: cfg.RedisAddr},
    asynq.Config{
      Concurrency: 10,                                   // workers per process
      Queues: map[string]int{
        "critical": 6,                                  // priority weight
        "default":  3,
        "low":      1,
      },
    },
  )

  mux := asynq.NewServeMux()
  jobs.NewHandlers(services).Register(mux)

  log.Fatal(srv.Run(mux))
}`}
      />
      <p>
        Run this as a separate process:{' '}
        <code>go run cmd/worker/main.go</code> in dev, or a separate
        deploy unit in prod. Two processes, one DB, one Redis.
      </p>

      <h2>Cron — scheduled jobs</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/jobs/cron.go"
        code={`func RegisterCron(s *asynq.Scheduler) error {
  // Every day at 2am
  _, err := s.Register("0 2 * * *", asynq.NewTask("nightly_digest", nil))
  if err != nil { return err }

  // Every hour
  _, err = s.Register("0 * * * *", asynq.NewTask("cleanup_uploads", nil))
  if err != nil { return err }

  return nil
}`}
      />
      <p>
        Standard cron syntax. The scheduler runs as part of the worker;
        it enqueues the task at the right time, the same workers process
        it.
      </p>

      <TipBox tone="warning">
        <strong>Cron in two processes = double-fire.</strong> If you
        run the scheduler in 2 worker pods, both will enqueue the same
        nightly digest. Run the scheduler in EXACTLY ONE process (an
        env flag like <code>RUN_SCHEDULER=true</code> on one worker) or
        you&apos;ll send 2× emails on cron days.
      </TipBox>

      <h2>The Jobs admin page</h2>
      <p>
        Grit&apos;s admin panel at <code>/admin/system/jobs</code>{' '}
        shows:
      </p>
      <ul>
        <li>Queue depth per queue (active, pending, retry, dead).</li>
        <li>Recent runs with status + duration.</li>
        <li>Failed jobs — click to inspect payload + error, re-queue.</li>
        <li>Throughput chart (jobs/min).</li>
      </ul>
      <p>
        When something&apos;s wrong, this is the first place to look.
      </p>

      <h2>How you call it from a service</h2>
      <CodeBlock
        language="go"
        code={`func (s *OrderService) Create(ctx context.Context, in CreateOrderInput) (Order, error) {
  order := models.Order{...}
  if err := s.db.WithContext(ctx).Create(&order).Error; err != nil { return Order{}, err }

  // 3 jobs, fire-and-forget — never block the create.
  // The IdempotencyKey makes a same-order retry safe — the receipt won't
  // fire twice even if the caller retries, even if a load-balancer
  // double-delivers, even if the job worker crashed mid-process.
  _ = s.jobs.EnqueueSendEmail(ctx, in.Email, "Receipt", "receipt",
      map[string]any{"order_id": order.ID},
      jobs.EnqueueOption{IdempotencyKey: "receipt:" + order.ID},
  )
  _ = s.jobs.Enqueue(ctx, "update_inventory",
      map[string]any{"order_id": order.ID},
      jobs.EnqueueOption{IdempotencyKey: "inventory:" + order.ID},
  )
  _ = s.jobs.Enqueue(ctx, "ask_for_review",
      map[string]any{"order_id": order.ID},
      jobs.EnqueueOption{
          Delay:          7 * 24 * time.Hour,
          IdempotencyKey: "review_request:" + order.ID,
      },
  )

  return order, nil
}`}
      />
      <p>
        Receipt: now. Inventory update: now. Review request: in 7 days.
        Customer&apos;s checkout response still returns in 50ms — and
        none of these jobs fire twice for the same order.
      </p>

      <h2>Idempotency keys — the safety net for retries</h2>
      <p>
        Retries cause duplicates. The most common production incident:
        a payment job fires, the worker crashes mid-flight, the next
        worker picks it up and charges the customer twice. Idempotency
        keys fix this at the queue level.
      </p>
      <p>
        Pass an <code>IdempotencyKey</code> in <code>EnqueueOption</code>{' '}
        and the queue refuses to enqueue a second task with the same
        key during the window (default 24h). The first call wins; the
        second silently returns success (the operation is already on
        its way, which is what the caller wanted).
      </p>
      <CodeBlock
        language="go"
        code={`// First call enqueues
err := s.jobs.Enqueue(ctx, "charge_card", payload, jobs.EnqueueOption{
    IdempotencyKey: "charge:" + invoice.ID,
})
// err == nil

// A retry within the window — dedup'd
err = s.jobs.Enqueue(ctx, "charge_card", payload, jobs.EnqueueOption{
    IdempotencyKey: "charge:" + invoice.ID,
})
// err == nil — the EnqueueSendEmail/EnqueueProcessImage helpers swallow
// the ErrDuplicateTask sentinel so you don't have to check for it on
// fire-and-forget paths. Use the generic Enqueue if you want to detect it.`}
      />

      <h3>Picking a good key</h3>
      <p>
        The key should be a natural business identifier for the
        action — not the payload, not a timestamp, not random bytes:
      </p>
      <ul>
        <li><code>&quot;receipt:&quot; + order.ID</code> — one receipt per order.</li>
        <li><code>&quot;welcome:&quot; + user.ID</code> — one welcome email per user.</li>
        <li><code>&quot;charge:&quot; + invoice.ID</code> — one charge per invoice.</li>
        <li><code>&quot;digest:&quot; + user.ID + &quot;:&quot; + today</code> — one digest per user per day.</li>
      </ul>

      <h2>Retry backoff — what actually happens on failure</h2>
      <p>
        Return an error from a job handler and asynq re-runs it with
        exponential backoff. Grit&apos;s explicit schedule (set via{' '}
        <code>ExponentialBackoff</code> in <code>workers.go</code>):
      </p>
      <CodeBlock
        language="text"
        code={`Attempt 1:  immediate
Attempt 2:  +1s     (after first failure)
Attempt 3:  +2s
Attempt 4:  +4s
Attempt 5:  +8s
Attempt 6:  +16s
Attempt 7:  +32s
Attempt 8:  +64s
Attempt 9:  +128s
Attempt 10: +256s
…           capped at 5 minutes
After DefaultMaxRetries failures (5 by default): dead-letter queue`}
      />
      <p>
        The total budget at default settings is ~16 minutes — long
        enough to ride out a typical downstream outage, short enough
        that the dead queue doesn&apos;t fill with hours-old work. Tune{' '}
        <code>MaxRetries</code> per-task when the failure mode warrants
        it (e.g., a webhook delivery to a flaky third party might want{' '}
        <code>MaxRetries: 25</code>).
      </p>

      <h2>How to modify this battery</h2>
      <ul>
        <li>
          <strong>Add a new task</strong> — register the handler in{' '}
          <code>handlers.go</code>, enqueue it from a service. That&apos;s
          it.
        </li>
        <li>
          <strong>Add a queue priority</strong> — edit the{' '}
          <code>Queues</code> map in <code>server.go</code>. Higher
          weight = more worker attention.
        </li>
        <li>
          <strong>Increase concurrency</strong> — bump{' '}
          <code>Concurrency: 10</code> to 50. Watch DB connection pool
          — if you exhaust it, raise the pool too.
        </li>
        <li>
          <strong>Change retry policy</strong> — pass{' '}
          <code>asynq.MaxRetry(3)</code> at enqueue, or set a default
          in the server config.
        </li>
      </ul>

      <KnowledgeCheck
        question="A job to resize an uploaded image fails because the S3 bucket is temporarily unreachable. With Grit's default config, what happens?"
        choices={[
          {
            label: 'The job fails permanently and the image is never resized',
            feedback:
              "No — that's why we use a job system. Retries are the whole point.",
          },
          {
            label: 'asynq retries the job with exponential backoff (a few times, increasing intervals), then dead-letters if it still fails. The admin Jobs page shows the failure with the error so you can investigate or re-queue.',
            correct: true,
            feedback:
              "Right — that's the resilience benefit. Transient failures heal themselves; persistent ones land in the dead queue for human review.",
          },
          {
            label: 'The worker process crashes',
            feedback:
              "A failed job returns an error from the handler — it doesn't crash the worker.",
          },
          {
            label: 'The whole job queue freezes',
            feedback:
              "Workers are concurrent; one failing job doesn't block others. Only that specific job retries.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Wire your first custom job:</p>
            <ol>
              <li>
                Define a task: &quot;send_birthday_email&quot; that
                takes a <code>user_id</code> payload.
              </li>
              <li>
                Register the handler — it loads the user, sends an
                email via the mail service.
              </li>
              <li>
                Enqueue it from a test endpoint{' '}
                <code>POST /api/dev/send-birthday/:user_id</code>.
              </li>
              <li>
                Boot the worker:{' '}
                <code>go run cmd/worker/main.go</code> in a second
                terminal.
              </li>
              <li>
                Hit the endpoint. Watch the worker logs — should see
                the job processed. Mailhog should show the email.
              </li>
              <li>
                Bonus: make the handler return an error on purpose
                once. Confirm it retries (asynq logs retry attempts).
              </li>
            </ol>
          </>
        }
        hint={
          <>
            For the retry test, return a generic error from the handler
            on the first invocation only — use a static counter in the
            handler struct. asynq will retry per its policy.
          </>
        }
        solution={
          <>
            <p>
              You should have an end-to-end job flow: enqueue → Redis →
              worker → email. The admin page should show the job;
              re-queueing a failed one should work.
            </p>
            <p>
              This pattern unlocks every &quot;don&apos;t block the
              request&quot; feature you&apos;ll ever build.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Last battery — <strong>AI (Claude + OpenAI)</strong>. Streaming
        chat, embeddings, and how Grit lets you swap providers without
        touching feature code.
      </p>
    </>
  )
}
