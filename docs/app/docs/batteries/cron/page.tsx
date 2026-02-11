import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'

export default function CronPage() {
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
                Cron Scheduler
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit uses the asynq Scheduler for cron-like recurring task execution. Schedule periodic
                cleanup jobs, report generation, or any recurring work using standard cron expressions.
                Tasks are enqueued into the same Redis-backed job queue used by background jobs.
              </p>
            </div>

            <div className="prose-grit">
              {/* How It Works */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  How It Works
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The cron scheduler is a separate component from the job worker. While the worker
                  processes jobs from the queue, the scheduler periodically enqueues new tasks according
                  to cron expressions. Both share the same Redis connection and task types.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">cron-flow.txt</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`┌────────────────────┐     ┌─────────┐     ┌──────────────────┐
│  Cron Scheduler    │     │  Redis  │     │  Worker          │
│                    │     │  Queue  │     │                  │
│  Every hour:       │     │         │     │                  │
│  enqueue           │────>│  task   │────>│  handleTokens    │
│  tokens:cleanup    │     │         │     │  Cleanup()       │
│                    │     │         │     │                  │
│  Custom schedule:  │     │         │     │                  │
│  enqueue           │────>│  task   │────>│  handleCustom()  │
│  report:generate   │     │         │     │                  │
└────────────────────┘     └─────────┘     └──────────────────┘`}</pre>
                </div>
              </div>

              {/* Scheduler Service */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Scheduler Service
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The scheduler lives at <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">internal/cron/cron.go</code>. It
                  wraps the asynq Scheduler and registers cron tasks at creation time.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">internal/cron/cron.go</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`// Task represents a registered cron task for display purposes.
type Task struct {
    Name     string \`json:"name"\`
    Schedule string \`json:"schedule"\`
    Type     string \`json:"type"\`
}

// RegisteredTasks holds the list of cron tasks for the admin API.
var RegisteredTasks []Task

// Scheduler wraps asynq.Scheduler for cron-like job scheduling.
type Scheduler struct {
    scheduler *asynq.Scheduler
}

// New creates a new cron Scheduler connected to Redis.
func New(redisURL string) (*Scheduler, error)

// Start begins executing scheduled tasks (runs in a goroutine).
func (s *Scheduler) Start() error

// Stop shuts down the scheduler gracefully.
func (s *Scheduler) Stop()`}</pre>
                </div>
              </div>

              {/* Built-in Tasks */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Built-in Tasks
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit ships with one built-in cron task: cleanup expired tokens, which runs every hour.
                  This removes soft-deleted user records older than 30 days.
                </p>

                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Task</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Schedule</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Type</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Description</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr>
                        <td className="px-4 py-2.5">Cleanup expired tokens</td>
                        <td className="px-4 py-2.5 font-mono text-xs">0 * * * *</td>
                        <td className="px-4 py-2.5 font-mono text-xs">tokens:cleanup</td>
                        <td className="px-4 py-2.5">Every hour, on the hour</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">internal/cron/cron.go (registration)</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`func New(redisURL string) (*Scheduler, error) {
    redisOpt, err := asynq.ParseRedisURI(redisURL)
    if err != nil {
        return nil, fmt.Errorf("parsing redis URL for cron: %w", err)
    }

    scheduler := asynq.NewScheduler(redisOpt, nil)

    // Register built-in cron tasks
    RegisteredTasks = []Task{}

    // Cleanup expired tokens -- every hour
    _, err = scheduler.Register("0 * * * *", asynq.NewTask("tokens:cleanup", nil))
    if err != nil {
        return nil, fmt.Errorf("registering tokens cleanup: %w", err)
    }
    RegisteredTasks = append(RegisteredTasks, Task{
        Name:     "Cleanup expired tokens",
        Schedule: "0 * * * *",
        Type:     "tokens:cleanup",
    })

    // grit:cron-tasks

    return &Scheduler{scheduler: scheduler}, nil
}`}</pre>
                </div>
              </div>

              {/* Adding Custom Cron Tasks */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Adding Custom Cron Tasks
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Add custom cron tasks by registering them in the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">New()</code> function.
                  Use standard cron expressions (5-field format). The task must have a corresponding handler
                  registered in the job worker.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">custom-cron-task.go</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`// Add after the "// grit:cron-tasks" marker:

// Generate daily reports -- every day at midnight
_, err = scheduler.Register("0 0 * * *", asynq.NewTask("report:daily", nil))
if err != nil {
    return nil, fmt.Errorf("registering daily report: %w", err)
}
RegisteredTasks = append(RegisteredTasks, Task{
    Name:     "Generate daily reports",
    Schedule: "0 0 * * *",
    Type:     "report:daily",
})

// Send weekly digest -- every Monday at 9 AM
payload, _ := json.Marshal(map[string]string{"type": "weekly"})
_, err = scheduler.Register("0 9 * * 1", asynq.NewTask("email:digest", payload))
if err != nil {
    return nil, fmt.Errorf("registering weekly digest: %w", err)
}
RegisteredTasks = append(RegisteredTasks, Task{
    Name:     "Send weekly digest",
    Schedule: "0 9 * * 1",
    Type:     "email:digest",
})`}</pre>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3">Cron Expression Reference</h3>

                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Expression</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Schedule</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">* * * * *</td>
                        <td className="px-4 py-2.5">Every minute</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">0 * * * *</td>
                        <td className="px-4 py-2.5">Every hour</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">0 0 * * *</td>
                        <td className="px-4 py-2.5">Every day at midnight</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">0 9 * * 1</td>
                        <td className="px-4 py-2.5">Every Monday at 9 AM</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">*/5 * * * *</td>
                        <td className="px-4 py-2.5">Every 5 minutes</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">0 0 1 * *</td>
                        <td className="px-4 py-2.5">First day of every month</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* The grit:cron-tasks Marker */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  The grit:cron-tasks Marker
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">// grit:cron-tasks</code> comment in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">cron.go</code> is
                  a marker used by the code generator. When you use <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit add cron</code> or
                  future CLI extensions, new cron tasks are injected at this marker position. You can
                  add tasks manually either above or below it -- just do not remove the marker.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">marker-location.go</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`    // ... built-in tasks above ...

    // grit:cron-tasks    <-- CLI injects new tasks here

    return &Scheduler{scheduler: scheduler}, nil
}`}</pre>
                </div>
              </div>

              {/* Admin Cron Viewer */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Admin Cron Viewer
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The admin panel shows all registered cron tasks via the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">GET /api/admin/cron/tasks</code> endpoint.
                  This endpoint reads from the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">RegisteredTasks</code> slice,
                  which is populated during scheduler initialization.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">internal/handlers/cron.go</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`// CronHandler handles admin cron task endpoints.
type CronHandler struct{}

// ListTasks returns all registered cron tasks.
func (h *CronHandler) ListTasks(c *gin.Context) {
    tasks := cron.RegisteredTasks
    if tasks == nil {
        tasks = []cron.Task{}
    }

    c.JSON(http.StatusOK, gin.H{
        "data": tasks,
    })
}

// Response example:
// {
//   "data": [
//     {
//       "name": "Cleanup expired tokens",
//       "schedule": "0 * * * *",
//       "type": "tokens:cleanup"
//     }
//   ]
// }`}</pre>
                </div>
              </div>

              {/* Lifecycle */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Lifecycle
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The scheduler is started in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">main.go</code> after the worker,
                  and stopped during graceful shutdown. It runs as a goroutine alongside the HTTP server.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">cmd/server/main.go (excerpt)</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`// Create and start the cron scheduler
cronScheduler, err := cron.New(cfg.RedisURL)
if err != nil {
    log.Printf("Warning: Cron scheduler failed to start: %v", err)
} else {
    cronScheduler.Start()
    defer cronScheduler.Stop()
    log.Println("Cron scheduler started")
}`}</pre>
                </div>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/batteries/jobs" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Background Jobs
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/batteries/caching" className="gap-1.5">
                  Redis Caching
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
