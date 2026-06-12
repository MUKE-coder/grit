import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Some work shouldn&apos;t block the HTTP request — sending an email,
        resizing an image, calling an LLM. Grit uses{' '}
        <a href="https://github.com/hibiken/asynq" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Asynq</a>{' '}
        for background jobs: enqueue from anywhere, a worker process consumes
        them. This lesson covers the pattern + the four primitives you&apos;ll
        actually use.
      </p>

      <h2>How it&apos;s shaped</h2>
      <CodeBlock
        language="text"
        code={`Producer (your handler/service)
  → jobs.Enqueue(ctx, task)          puts JSON into Redis
                                     returns immediately

Worker (separate Go process)
  → polls Redis, dequeues tasks
  → calls the registered handler for the task type
  → marks done / retries on error / moves to dead-letter`}
      />
      <p>
        Two processes — the API and the worker — share Redis as the queue.
        Both run from the same binary; you start the worker with{' '}
        <code>./api worker</code> instead of <code>./api server</code>.
      </p>

      <h2>Defining a task</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/jobs/welcome_email.go"
        code={`const TaskTypeWelcomeEmail = "user:welcome_email"

type WelcomeEmailPayload struct {
    UserID string \`json:"user_id"\`
}

func (j *Jobs) EnqueueWelcomeEmail(ctx context.Context, userID string) error {
    payload, _ := json.Marshal(WelcomeEmailPayload{UserID: userID})
    _, err := j.client.EnqueueContext(ctx, asynq.NewTask(TaskTypeWelcomeEmail, payload),
        asynq.Queue("default"),
        asynq.MaxRetry(5),
    )
    return err
}`}
      />

      <h2>Handling the task in the worker</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/jobs/welcome_email_handler.go"
        code={`func (j *Jobs) HandleWelcomeEmail(ctx context.Context, task *asynq.Task) error {
    var p WelcomeEmailPayload
    if err := json.Unmarshal(task.Payload(), &p); err != nil {
        return fmt.Errorf("unmarshal: %w", err)
    }
    user, err := j.userService.GetByID(ctx, p.UserID)
    if err != nil {
        return fmt.Errorf("get user: %w", err)
    }
    return j.mailer.SendWelcome(ctx, user.Email, user.Name)
}`}
      />
      <p>
        Return <code>nil</code> = task completed. Return an error = Asynq
        retries with exponential backoff up to <code>MaxRetry</code>, then
        moves it to the dead-letter queue.
      </p>

      <h2>Enqueuing from a handler</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/auth.go (excerpt)"
        code={`// After register succeeds
if err := h.jobs.EnqueueWelcomeEmail(c.Request.Context(), user.ID); err != nil {
    // Log + continue — registration succeeded; the welcome email is best-effort
    log.Printf("welcome email enqueue: %v", err)
}
respond.Created(c, user, "Account created")`}
      />
      <p>
        Notice — enqueue failure doesn&apos;t fail the request. The user&apos;s
        account exists; the welcome email is best-effort. Different
        behaviour for critical paths (payments, audit) where you&apos;d return
        the error.
      </p>

      <h2>Scheduled tasks (cron)</h2>
      <p>For recurring work (daily digests, hourly metrics roll-up):</p>
      <CodeBlock
        language="go"
        code={`// In internal/jobs/scheduler.go
scheduler.Register("0 9 * * *", "daily:digest", nil)  // 9 AM daily
scheduler.Register("0 * * * *", "metrics:rollup", nil) // hourly`}
      />
      <p>
        Same handler signature as one-off tasks. The scheduler runs as part
        of the worker process — no separate cron daemon.
      </p>

      <TipBox tone="warning">
        <strong>Don&apos;t run the worker inside the API process.</strong> If a
        long-running job hangs, it blocks an HTTP request slot. Always boot
        the worker separately: <code>./api worker</code> on its own
        container/VM/systemd unit.
      </TipBox>

      <h2>The Asynq admin UI</h2>
      <p>
        Grit&apos;s scaffold includes <code>/admin/jobs</code> in the admin app
        — visual dashboard for active/pending/failed/dead jobs, with retry
        + delete buttons. For the API-only kit, expose Asynq&apos;s web UI on a
        dedicated port; password-protected via the Sentinel dashboard
        credentials.
      </p>

      <KnowledgeCheck
        question="A user signs up. You enqueue a welcome-email task. The worker is down for an hour. What happens to the task?"
        choices={[
          {
            label: 'Lost — Asynq dropped it',
            feedback:
              "Wrong — Asynq persists tasks in Redis. As long as Redis is alive (and the data survives), the task waits.",
          },
          {
            label: 'It sits in Redis until the worker comes back, then runs',
            correct: true,
            feedback:
              "Right — the queue persists; tasks survive worker outages as long as Redis stays alive. The user gets their welcome email an hour late.",
          },
          {
            label: 'The API blocks waiting for the worker',
            feedback:
              "Wrong — enqueue returns immediately. The worker is decoupled; that's the whole point of background jobs.",
          },
          {
            label: 'Asynq fails the API request',
            feedback:
              "Wrong — enqueue puts the JSON in Redis and returns. The user's request completes regardless of worker state.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Wire your first job. Register a new task type:</p>
            <ol>
              <li>
                Define <code>TaskTypeLogHello</code> + a handler that just{' '}
                <code>log.Println(&quot;hello from a job!&quot;)</code>.
              </li>
              <li>
                Add an enqueue call to your <code>/api/health</code> handler.
              </li>
              <li>
                Start the worker: <code>go run ./cmd/server worker</code>{' '}
                (or whatever the scaffolded subcommand is).
              </li>
              <li>
                Hit <code>/api/health</code> — you should see &quot;hello from a
                job!&quot; in the worker&apos;s output.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            Register the handler with <code>mux.HandleFunc(TaskTypeLogHello,
            j.HandleLogHello)</code> in the worker bootstrap.
          </>
        }
        solution={
          <>
            <p>
              You should see two log streams — your API logs the request,
              the worker logs the job:
            </p>
            <CodeBlock
              language="text"
              code={`[API] GIN-debug: /api/health 200
[WORKER] hello from a job!`}
            />
            <p>
              That&apos;s producer + consumer working. Now you can offload
              anything heavy from the request path.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Transactional email is the most common job. Next lesson — how to
        send it via Resend and watch it land in Mailhog during dev.
      </p>
    </>
  )
}
