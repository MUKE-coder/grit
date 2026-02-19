import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export default function JobsPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Batteries</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Background Jobs
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit uses <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">asynq</code> -- a Redis-backed job
                queue library for Go -- to handle background processing. Send emails, generate thumbnails,
                clean up expired tokens, or run any async task with automatic retries and priority queues.
              </p>
            </div>

            <div className="prose-grit">
              {/* Architecture */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Architecture
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The job system has two components: a <strong>client</strong> that enqueues jobs
                  and a <strong>worker</strong> that processes them. Both connect to the same Redis instance.
                  The worker runs as a goroutine inside the API server -- no separate process needed.
                </p>

                <CodeBlock language="bash" filename="job-architecture.txt" code={`┌──────────────────┐     ┌─────────┐     ┌──────────────────┐
│  API Handler     │     │  Redis  │     │  Worker          │
│                  │     │         │     │                  │
│  client.Enqueue  │────>│  Queue  │────>│  handleEmailSend │
│  SendEmail(...)  │     │         │     │  handleImage...  │
│                  │     │         │     │  handleCleanup   │
└──────────────────┘     └─────────┘     └──────────────────┘`} />
              </div>

              {/* Job Client */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Job Client
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The job client at <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">internal/jobs/client.go</code> provides
                  typed methods for enqueuing each built-in job. It handles JSON serialization of payloads
                  and configures retry policies per job type.
                </p>

                <CodeBlock language="go" filename="internal/jobs/client.go" code={`// Task type constants
const (
    TypeEmailSend     = "email:send"
    TypeImageProcess  = "image:process"
    TypeTokensCleanup = "tokens:cleanup"
)

// Client wraps asynq.Client for enqueuing background jobs.
type Client struct {
    client *asynq.Client
}

// NewClient creates a new job queue client connected to Redis.
func NewClient(redisURL string) (*Client, error)

// Close shuts down the client connection.
func (c *Client) Close() error`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3">Enqueue Methods</h3>

                <CodeBlock language="go" filename="internal/jobs/client.go (enqueue methods)" code={`// EnqueueSendEmail enqueues an email send job.
// Max retries: 3
func (c *Client) EnqueueSendEmail(
    to, subject, template string,
    data map[string]interface{},
) error

// EnqueueProcessImage enqueues an image processing job.
// Max retries: 2
func (c *Client) EnqueueProcessImage(
    uploadID uint,
    key, mimeType string,
) error

// EnqueueTokensCleanup enqueues a token cleanup job.
// Max retries: 1
func (c *Client) EnqueueTokensCleanup() error`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3">Job Payloads</h3>

                <CodeBlock language="go" filename="internal/jobs/client.go (payloads)" code={`// EmailPayload holds the data for an email send job.
type EmailPayload struct {
    To       string                 \`json:"to"\`
    Subject  string                 \`json:"subject"\`
    Template string                 \`json:"template"\`
    Data     map[string]interface{} \`json:"data"\`
}

// ImagePayload holds the data for an image processing job.
type ImagePayload struct {
    UploadID uint   \`json:"upload_id"\`
    Key      string \`json:"key"\`
    MimeType string \`json:"mime_type"\`
}`} />
              </div>

              {/* Worker Setup */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Worker Setup
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The worker is started in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">main.go</code> alongside
                  the HTTP server. It receives all service dependencies via the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">WorkerDeps</code> struct,
                  giving handler functions access to the database, mailer, storage, and cache.
                </p>

                <CodeBlock language="go" filename="internal/jobs/workers.go" code={`// WorkerDeps holds dependencies needed by job handlers.
type WorkerDeps struct {
    DB      *gorm.DB
    Mailer  *mail.Mailer
    Storage *storage.Storage
    Cache   *cache.Cache
}

// StartWorker starts the asynq worker server in a goroutine.
// Returns a stop function and any startup error.
func StartWorker(redisURL string, deps WorkerDeps) (func(), error) {
    redisOpt, _ := asynq.ParseRedisURI(redisURL)

    srv := asynq.NewServer(redisOpt, asynq.Config{
        Concurrency: 10,
        Queues: map[string]int{
            "default":  6,
            "critical": 3,
            "low":      1,
        },
    })

    mux := asynq.NewServeMux()
    mux.HandleFunc(TypeEmailSend, handleEmailSend(deps))
    mux.HandleFunc(TypeImageProcess, handleImageProcess(deps))
    mux.HandleFunc(TypeTokensCleanup, handleTokensCleanup(deps))

    go func() {
        srv.Run(mux)
    }()

    return func() { srv.Shutdown() }, nil
}`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3">Starting the Worker in main.go</h3>

                <CodeBlock language="go" filename="cmd/server/main.go (excerpt)" code={`// Start the background worker
stopWorker, err := jobs.StartWorker(cfg.RedisURL, jobs.WorkerDeps{
    DB:      db,
    Mailer:  mailer,
    Storage: store,
    Cache:   cacheService,
})
if err != nil {
    log.Fatalf("Failed to start worker: %v", err)
}
defer stopWorker()`} />
              </div>

              {/* Queue Priorities */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Queue Priorities
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Jobs are distributed across three priority queues. The worker allocates processing
                  capacity based on these weights: critical gets 30%, default gets 60%, and low gets 10%.
                </p>

                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Queue</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Weight</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Use Case</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">critical</td>
                        <td className="px-4 py-2.5">3 (30%)</td>
                        <td className="px-4 py-2.5">Password resets, payment webhooks</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">default</td>
                        <td className="px-4 py-2.5">6 (60%)</td>
                        <td className="px-4 py-2.5">Emails, image processing</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">low</td>
                        <td className="px-4 py-2.5">1 (10%)</td>
                        <td className="px-4 py-2.5">Cleanup, analytics, reports</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Retry Configuration */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Retry Configuration
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Each job type has a configured maximum retry count. When a handler returns an error,
                  asynq automatically retries with exponential backoff. After exhausting all retries,
                  the job is moved to the &quot;archived&quot; (failed) state.
                </p>

                <CodeBlock language="go" filename="retry-config.go" code={`// Email: 3 retries (important to deliver)
task := asynq.NewTask(TypeEmailSend, payload)
_, err = c.client.Enqueue(task, asynq.MaxRetry(3))

// Image: 2 retries (can be re-triggered)
task := asynq.NewTask(TypeImageProcess, payload)
_, err = c.client.Enqueue(task, asynq.MaxRetry(2))

// Cleanup: 1 retry (runs hourly anyway)
task := asynq.NewTask(TypeTokensCleanup, nil)
_, err = c.client.Enqueue(task, asynq.MaxRetry(1))

// Custom: enqueue to a specific queue with custom retry
task := asynq.NewTask("invoice:generate", payload)
_, err = c.client.Enqueue(task,
    asynq.MaxRetry(5),
    asynq.Queue("critical"),
    asynq.Timeout(30*time.Second),
)`} />
              </div>

              {/* Adding Custom Jobs */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Adding Custom Jobs
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  To add a new job type, define the task type constant, create a payload struct,
                  add an enqueue method to the client, write the handler function, and register it
                  in the worker mux.
                </p>

                <CodeBlock language="go" filename="custom-job-example.go" code={`// 1. Add task type constant
const TypeInvoiceGenerate = "invoice:generate"

// 2. Define payload
type InvoicePayload struct {
    OrderID   uint   \`json:"order_id"\`
    UserEmail string \`json:"user_email"\`
}

// 3. Add enqueue method to Client
func (c *Client) EnqueueGenerateInvoice(orderID uint, email string) error {
    payload, _ := json.Marshal(InvoicePayload{
        OrderID:   orderID,
        UserEmail: email,
    })
    task := asynq.NewTask(TypeInvoiceGenerate, payload)
    _, err := c.client.Enqueue(task, asynq.MaxRetry(3))
    return err
}

// 4. Write handler function
func handleInvoiceGenerate(deps WorkerDeps) func(ctx context.Context, task *asynq.Task) error {
    return func(ctx context.Context, task *asynq.Task) error {
        var payload InvoicePayload
        if err := json.Unmarshal(task.Payload(), &payload); err != nil {
            return err
        }
        // Generate PDF, send email, etc.
        return nil
    }
}

// 5. Register in worker mux (in StartWorker)
mux.HandleFunc(TypeInvoiceGenerate, handleInvoiceGenerate(deps))`} />
              </div>

              {/* Admin Jobs Dashboard */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Admin Jobs Dashboard
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The admin panel includes a jobs dashboard that shows queue statistics and allows
                  admins to view, retry, and clear jobs. The dashboard uses the asynq Inspector API
                  under the hood.
                </p>

                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Endpoint</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Method</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Description</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">/api/admin/jobs/stats</td>
                        <td className="px-4 py-2.5 font-mono text-xs">GET</td>
                        <td className="px-4 py-2.5">Queue stats (pending, active, completed, failed)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">/api/admin/jobs/:status</td>
                        <td className="px-4 py-2.5 font-mono text-xs">GET</td>
                        <td className="px-4 py-2.5">List jobs by status (active, pending, completed, failed, retry)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">/api/admin/jobs/:id/retry</td>
                        <td className="px-4 py-2.5 font-mono text-xs">POST</td>
                        <td className="px-4 py-2.5">Retry a failed job</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">/api/admin/jobs/queue/:queue</td>
                        <td className="px-4 py-2.5 font-mono text-xs">DELETE</td>
                        <td className="px-4 py-2.5">Clear all completed tasks in a queue</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5">
                  <p className="text-sm text-foreground/80">
                    <strong>Concurrency:</strong> The default worker concurrency is 10. This means up to 10
                    jobs can be processed simultaneously. Adjust the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Concurrency</code> setting
                    in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">StartWorker()</code> based on your server resources.
                  </p>
                </div>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/batteries/email" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Email System
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/batteries/cron" className="gap-1.5">
                  Cron Scheduler
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
