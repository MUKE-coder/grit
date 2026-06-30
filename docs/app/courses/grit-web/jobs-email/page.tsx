import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Background Jobs & Email — Grit Web Course',
  description: 'Learn how to use Redis job queues with asynq for background processing and send transactional emails with Resend. Covers workers, retry logic, cron scheduling, email templates, and the admin jobs dashboard.',
}

export default function JobsEmailCourse() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <Link href="/courses/grit-web" className="hover:text-foreground transition-colors">Grit Web</Link>
          <span>/</span>
          <span className="text-foreground">Background Jobs & Email</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Course 6 of 8</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">14 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Background Jobs & Email
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Some tasks are too slow or too unreliable to run while a user waits for a response. In this course,
            you will learn how Grit uses Redis-backed job queues to process work in the background — sending emails,
            processing images, generating reports — and how to send beautiful transactional emails with templates.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What are Background Jobs? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What are Background Jobs?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When a user clicks {'"'}Send Email{'"'} in your app, what should happen? If you send the email
            right there in the API handler, the user has to wait — maybe 2-5 seconds — while the email service
            processes the request. That{"'"}s a terrible experience.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Instead, the smart approach is: the API handler puts a message on a <strong className="text-foreground">queue</strong> saying
            {'"'}send this email{'"'}, responds to the user immediately with {'"'}Email queued!{'"'},
            and a separate <strong className="text-foreground">worker</strong> picks up that message and sends the email
            in the background. The user never waits.
          </p>

          <Definition term="Background Job">
            A task that runs asynchronously — outside the user{"'"}s HTTP request cycle. The API
            responds immediately, and the actual work happens later in a separate process. This keeps
            the user interface fast and responsive.
          </Definition>

          <Definition term="Job Queue">
            A data structure (usually backed by Redis) that holds a list of tasks waiting to be processed.
            Producers add tasks to the queue, and workers consume them one at a time. If a task fails,
            it can be retried automatically.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here are common tasks that should be background jobs instead of blocking the user{"'"}s request:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Sending emails</strong> — welcome emails, password resets, order confirmations</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Processing images</strong> — generating thumbnails, resizing uploads, optimizing file sizes</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Generating reports</strong> — PDF exports, CSV downloads, analytics calculations</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Cleanup tasks</strong> — removing orphaned files, purging expired sessions, archiving old data</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Third-party API calls</strong> — payment webhooks, push notifications, syncing with external services</li>
          </ul>

          <Note>
            A good rule of thumb: if a task takes more than 200ms or depends on an external service that might
            be slow or unreliable, it should be a background job.
          </Note>

          <Challenge number={1} title="Identify Background Job Candidates">
            <p>
              Name 3 things in a web app that should be background jobs instead of blocking the user{"'"}s
              request. For each one, explain why it would be bad to run it inline during the HTTP request.
            </p>
          </Challenge>
        </section>

        {/* ═══ Section 2: How asynq Works ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">How asynq Works</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit uses <Code>asynq</Code> — a Go library for distributed task processing backed by Redis.
            It{"'"}s simple, fast, and battle-tested. Here{"'"}s how the pieces fit together:
          </p>

          <Definition term="asynq">
            A Go library for background task processing. It uses Redis as a message broker to coordinate
            between producers (your API handlers) and consumers (worker processes). Think of it as
            Sidekiq for Go.
          </Definition>

          <Definition term="Worker">
            A long-running process that listens to the Redis queue, picks up tasks, and processes them.
            Workers run separately from your API server — they{"'"}re their own process.
          </Definition>

          <Definition term="Redis (as Job Broker)">
            Redis acts as the middleman between your API and your workers. When a handler enqueues a task,
            it gets stored in Redis. Workers poll Redis for new tasks. Redis ensures tasks aren{"'"}t lost
            and tracks their state (pending, active, completed, failed).
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The flow works in 4 steps:
          </p>

          <ol className="space-y-3 text-muted-foreground mb-6">
            <li className="flex gap-3">
              <span className="text-primary font-bold shrink-0">1.</span>
              <span><strong className="text-foreground">Handler dispatches a task</strong> — Your API handler creates a task with a type and payload, then pushes it to Redis via the asynq client.</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-bold shrink-0">2.</span>
              <span><strong className="text-foreground">Worker picks it up</strong> — The asynq worker process polls Redis, sees the new task, and dequeues it.</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-bold shrink-0">3.</span>
              <span><strong className="text-foreground">Worker processes it</strong> — The worker matches the task type to a registered handler function and executes it.</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-bold shrink-0">4.</span>
              <span><strong className="text-foreground">Task marked done/failed</strong> — If the handler returns <Code>nil</Code>, the task is marked as completed. If it returns an error, it{"'"}s marked as failed and may be retried.</span>
            </li>
          </ol>

          <CodeBlock filename="Architecture Diagram">
{`┌──────────────┐     enqueue     ┌─────────┐     dequeue     ┌──────────────┐
│  API Handler  │ ──────────────► │  Redis  │ ◄────────────── │    Worker    │
│  (producer)   │                 │ (broker) │                 │  (consumer)  │
└──────────────┘                 └─────────┘                 └──────────────┘
                                                                    │
                                                                    ▼
                                                            ┌──────────────┐
                                                            │  Task Handler │
                                                            │  (your code)  │
                                                            └──────────────┘`}
          </CodeBlock>

          <Challenge number={2} title="Check Redis">
            <p>
              Make sure Redis is running by checking your Docker containers:
            </p>
            <CodeBlock filename="Terminal">
{`docker compose ps`}
            </CodeBlock>
            <p className="mt-3">
              What port does Redis use? (Hint: it{"'"}s the default Redis port.)
            </p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Built-in Workers ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Built-in Workers</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit comes with 3 workers out of the box. You don{"'"}t need to build them — they{"'"}re already
            wired up and ready to process tasks the moment your app starts.
          </p>

          <ul className="space-y-4 text-muted-foreground mb-6">
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Email worker</strong> — Sends transactional emails via Resend. Handles welcome emails, password resets, email verification, and custom notifications. If Resend is down, the task retries automatically.</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Image worker</strong> — Generates thumbnails after a file is uploaded. Resizes images to multiple dimensions, optimizes file size, and stores the results back in S3/MinIO.</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Cleanup worker</strong> — Removes old or orphaned files. Runs on a schedule to purge expired uploads, temporary files, and data that{"'"}s no longer referenced by any record.</span>
            </li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            All worker code lives in a single directory:
          </p>

          <CodeBlock filename="Project Structure">
{`apps/api/internal/jobs/
├── email.go        # Email worker — sends emails via Resend
├── image.go        # Image worker — generates thumbnails
├── cleanup.go      # Cleanup worker — removes orphaned files
├── scheduler.go    # Cron scheduler — recurring tasks
└── worker.go       # Worker setup — registers all handlers`}
          </CodeBlock>

          <Tip>
            The <Code>worker.go</Code> file is where all task handlers are registered with the asynq mux.
            When you create a new custom job, this is where you{"'"}ll add its handler registration.
          </Tip>

          <Challenge number={3} title="Explore the Workers">
            <p>
              Open the <Code>apps/api/internal/jobs/</Code> directory. List the files you see. For each file,
              describe in one sentence what that worker does based on the filename and any comments at the top.
            </p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Creating a Custom Job ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Creating a Custom Job</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Let{"'"}s walk through creating a custom background job step by step. We{"'"}ll build a
            report generation job — when a user requests a report, we{"'"}ll queue it and generate it
            in the background.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every asynq job has 4 parts: a <strong className="text-foreground">task type</strong> (a string identifier),
            a <strong className="text-foreground">payload</strong> (data the worker needs),
            a <strong className="text-foreground">handler</strong> (the function that does the work),
            and a <strong className="text-foreground">registration</strong> (telling the worker about the handler).
          </p>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">Step 1: Define the Task Type and Payload</h3>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The task type is a unique string that identifies this kind of job. The payload is a struct
            containing all the data the worker needs to process the task.
          </p>

          <CodeBlock filename="apps/api/internal/jobs/report.go">
{`package jobs

const TypeReportGenerate = "report:generate"

type ReportPayload struct {
    UserID uint   \`json:"user_id"\`
    Type   string \`json:"type"\`
}`}
          </CodeBlock>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">Step 2: Create the Handler</h3>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The handler is a function that receives the task, unmarshals the payload, and does the actual work.
            If it returns <Code>nil</Code>, the task is marked as completed. If it returns an error, the task
            is marked as failed and may be retried.
          </p>

          <CodeBlock filename="apps/api/internal/jobs/report.go">
{`func HandleReportGenerate(ctx context.Context, t *asynq.Task) error {
    var p ReportPayload
    if err := json.Unmarshal(t.Payload(), &p); err != nil {
        return err
    }

    // Generate the report...
    log.Printf("Generating %s report for user %d", p.Type, p.UserID)

    // In a real app, you would:
    // 1. Query the database for the relevant data
    // 2. Format it into a PDF or CSV
    // 3. Upload the file to S3
    // 4. Notify the user that the report is ready

    return nil
}`}
          </CodeBlock>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">Step 3: Register the Handler</h3>

          <p className="text-muted-foreground leading-relaxed mb-4">
            In <Code>worker.go</Code>, register the handler so the worker knows which function to call
            when it sees a task with this type:
          </p>

          <CodeBlock filename="apps/api/internal/jobs/worker.go">
{`mux.HandleFunc(TypeReportGenerate, HandleReportGenerate)`}
          </CodeBlock>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">Step 4: Dispatch from a Handler or Service</h3>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Finally, enqueue the task from your API handler or service. Create the payload, wrap it in
            an asynq task, and enqueue it:
          </p>

          <CodeBlock filename="apps/api/internal/handlers/report_handler.go">
{`payload, _ := json.Marshal(jobs.ReportPayload{
    UserID: 1,
    Type:   "monthly",
})
task := asynq.NewTask(jobs.TypeReportGenerate, payload)
info, err := client.Enqueue(task)
if err != nil {
    // Handle error
}
log.Printf("Task enqueued: id=%s queue=%s", info.ID, info.Queue)`}
          </CodeBlock>

          <Note>
            Notice how the API handler doesn{"'"}t generate the report itself. It just puts a message on the
            queue and responds to the user immediately. The worker does the heavy lifting later.
          </Note>

          <Challenge number={4} title="Read the Email Worker">
            <p>
              Read through the email worker code in <Code>apps/api/internal/jobs/email.go</Code>. Can you
              identify all 4 parts?
            </p>
            <ol className="mt-2 space-y-1 list-decimal list-inside">
              <li>What is the task type string?</li>
              <li>What fields are in the payload struct?</li>
              <li>What does the handler function do?</li>
              <li>Where is it registered?</li>
            </ol>
          </Challenge>
        </section>

        {/* ═══ Section 5: Retry Logic ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Retry Logic</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            What happens when a background job fails? Maybe the email API is down, or the database is
            temporarily unavailable. You don{"'"}t want to lose the task — you want to retry it.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            asynq handles this automatically. By default, failed tasks are retried <strong className="text-foreground">3 times</strong> with
            exponential backoff. Each retry waits longer than the last, giving the failing service time to recover.
          </p>

          <Definition term="Exponential Backoff">
            A retry strategy where each attempt waits longer than the last: 1s, 2s, 4s, 8s, 16s, and so on.
            This prevents overwhelming a failing service with rapid retries. Instead of hammering a down
            server with requests every second, you gradually back off and give it time to recover.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You can configure retry behavior per task when you enqueue it:
          </p>

          <CodeBlock filename="Retry Configuration">
{`// Retry up to 5 times (default is 3)
client.Enqueue(task, asynq.MaxRetry(5))

// Delay processing by 10 seconds
client.Enqueue(task, asynq.ProcessIn(10 * time.Second))

// Set a deadline — fail if not completed by this time
client.Enqueue(task, asynq.Deadline(time.Now().Add(1 * time.Hour)))

// Combine options
client.Enqueue(task,
    asynq.MaxRetry(5),
    asynq.ProcessIn(30 * time.Second),
    asynq.Queue("critical"),
)`}
          </CodeBlock>

          <Tip>
            Use <Code>asynq.Queue({'"'}critical{'"'})</Code> to put important tasks (like payment processing)
            on a separate queue with higher priority. Less important tasks (like cleanup) can go on a
            {'"'}low{'"'} priority queue.
          </Tip>

          <Challenge number={5} title="Understand Retry Behavior">
            <p>
              What happens if a background job fails 3 times (the default max retry)? Where can you see
              failed jobs? (Hint: check the admin panel{"'"}s System pages.)
            </p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Cron Scheduler ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Cron Scheduler</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Some tasks need to run on a schedule — not triggered by a user action, but by the clock.
            Clean up expired sessions every night. Generate a daily report at 9am. Check for abandoned
            carts every hour. These are <strong className="text-foreground">scheduled tasks</strong>.
          </p>

          <Definition term="Cron Expression">
            A string that defines a recurring schedule using 5 fields: minute, hour, day of month, month,
            and day of week. Originally from Unix, cron expressions are the universal standard for
            scheduling recurring tasks.
          </Definition>

          <Definition term="Scheduled Task">
            A background job that runs automatically on a recurring schedule, defined by a cron expression.
            Unlike regular jobs that are triggered by user actions, scheduled tasks are triggered by time.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit uses asynq{"'"}s built-in scheduler. Here{"'"}s the cron expression format:
          </p>

          <CodeBlock filename="Cron Expression Format">
{`┌───────── minute (0-59)
│ ┌─────── hour (0-23)
│ │ ┌───── day of month (1-31)
│ │ │ ┌─── month (1-12)
│ │ │ │ ┌─ day of week (0-6, 0=Sunday)
│ │ │ │ │
* * * * *`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>*</Code> means {'"'}every{'"'}. Here are common patterns:
          </p>

          <div className="bg-[#0d1526] border border-border/40 rounded-lg overflow-hidden mb-6">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border/40">
                  <th className="text-left px-4 py-2 text-foreground font-medium">Expression</th>
                  <th className="text-left px-4 py-2 text-foreground font-medium">Meaning</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2"><Code>0 0 * * *</Code></td>
                  <td className="px-4 py-2">Every day at midnight</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2"><Code>*/5 * * * *</Code></td>
                  <td className="px-4 py-2">Every 5 minutes</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2"><Code>0 9 * * 1</Code></td>
                  <td className="px-4 py-2">Every Monday at 9:00 AM</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2"><Code>0 */2 * * *</Code></td>
                  <td className="px-4 py-2">Every 2 hours</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2"><Code>30 3 * * *</Code></td>
                  <td className="px-4 py-2">Every day at 3:30 AM</td>
                </tr>
                <tr>
                  <td className="px-4 py-2"><Code>0 0 1 * *</Code></td>
                  <td className="px-4 py-2">First day of every month at midnight</td>
                </tr>
              </tbody>
            </table>
          </div>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s how Grit registers scheduled tasks in the scheduler:
          </p>

          <CodeBlock filename="apps/api/internal/jobs/scheduler.go">
{`scheduler := asynq.NewScheduler(redisOpt, nil)

// Clean up orphaned files every day at 2am
scheduler.Register("0 2 * * *", asynq.NewTask(TypeCleanupFiles, nil))

// Generate daily stats at midnight
scheduler.Register("0 0 * * *", asynq.NewTask(TypeStatsDaily, nil))

// Check for expired sessions every hour
scheduler.Register("0 * * * *", asynq.NewTask(TypeSessionCleanup, nil))`}
          </CodeBlock>

          <Challenge number={6} title="Write Cron Expressions">
            <p>
              What cron expression would you use for each of these schedules?
            </p>
            <ol className="mt-2 space-y-1 list-decimal list-inside">
              <li>Every hour on the hour</li>
              <li>Every day at 3:30 PM</li>
              <li>Every Sunday at midnight</li>
            </ol>
          </Challenge>
        </section>

        {/* ═══ Section 7: Admin Jobs Dashboard ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Admin Jobs Dashboard</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit{"'"}s admin panel includes a built-in Jobs dashboard where you can monitor all background
            tasks. No need for a separate monitoring tool — it{"'"}s right there in your admin panel.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The Jobs page shows tasks grouped by status:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Active</strong> — tasks currently being processed by a worker</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Pending</strong> — tasks waiting in the queue to be picked up</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Completed</strong> — tasks that finished successfully</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Failed</strong> — tasks that failed after all retries were exhausted</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Scheduled</strong> — tasks delayed with <Code>ProcessIn</Code>, waiting for their time</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For each failed task, you can see the error message, how many times it was retried, and when it
            last failed. You can also <strong className="text-foreground">retry a failed task</strong> directly
            from the dashboard — no code changes needed.
          </p>

          <Tip>
            The Jobs dashboard is one of the admin panel{"'"}s System pages. You can access it from the
            sidebar under {'"'}System{'"'}. It{"'"}s useful for debugging issues in production — if emails
            aren{"'"}t going out, check the Jobs page first.
          </Tip>

          <Challenge number={7} title="Explore the Jobs Dashboard">
            <p>
              Visit the Jobs page in the admin panel at <Code>localhost:3001</Code>. Navigate to
              System &gt; Jobs. Is there a queue? Are there any completed or failed jobs?
              Try triggering an action that creates a background job (like registering a new user)
              and watch the Jobs page update.
            </p>
          </Challenge>
        </section>

        {/* ═══ Section 8: What is Transactional Email? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is Transactional Email?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            There are two kinds of email in web applications: <strong className="text-foreground">marketing email</strong> and
            {' '}<strong className="text-foreground">transactional email</strong>. They serve very different purposes.
          </p>

          <Definition term="Transactional Email">
            An email triggered by a specific user action or system event. Examples: password reset links,
            order confirmations, welcome emails after registration, account verification codes. These are
            one-to-one emails sent in response to something the user did — not bulk campaigns.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s how they differ:
          </p>

          <div className="bg-[#0d1526] border border-border/40 rounded-lg overflow-hidden mb-6">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border/40">
                  <th className="text-left px-4 py-2 text-foreground font-medium">Transactional</th>
                  <th className="text-left px-4 py-2 text-foreground font-medium">Marketing</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2">Triggered by user action</td>
                  <td className="px-4 py-2">Sent in bulk to a list</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2">One recipient at a time</td>
                  <td className="px-4 py-2">Hundreds or thousands at once</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2">Expected by the user</td>
                  <td className="px-4 py-2">May be unsolicited</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2">Password resets, receipts</td>
                  <td className="px-4 py-2">Newsletters, promotions</td>
                </tr>
                <tr>
                  <td className="px-4 py-2">High priority, must arrive fast</td>
                  <td className="px-4 py-2">Can be delayed or batched</td>
                </tr>
              </tbody>
            </table>
          </div>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit uses <strong className="text-foreground">Resend</strong> for transactional email. Resend is a
            modern email API designed for developers — simple, reliable, and with excellent deliverability.
            In local development, Grit uses <strong className="text-foreground">Mailhog</strong> to catch all
            outgoing emails so you can preview them without sending real emails.
          </p>
        </section>

        {/* ═══ Section 9: Email Configuration ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Email Configuration</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Email configuration lives in your <Code>.env</Code> file. There are two variables you need:
          </p>

          <CodeBlock filename=".env">
{`# Email — Resend
RESEND_API_KEY=re_xxxxxxxxxxxxxxxxxxxx
MAIL_FROM=noreply@yourapp.com`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <Code>RESEND_API_KEY</Code> is your API key from{' '}
            <Link href="https://resend.com" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">resend.com</Link>.
            <Code>MAIL_FROM</Code> is the {'"'}from{'"'} address that appears on all outgoing emails.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For local development, you don{"'"}t need a real Resend API key. Grit{"'"}s Docker Compose
            includes <strong className="text-foreground">Mailhog</strong> — a fake SMTP server that catches
            all outgoing emails and displays them in a web UI at <Code>localhost:8025</Code>.
          </p>

          <CodeBlock filename="docker-compose.yml (Mailhog section)">
{`mailhog:
  image: mailhog/mailhog
  ports:
    - "1025:1025"   # SMTP port — your app sends here
    - "8025:8025"   # Web UI — view caught emails here`}
          </CodeBlock>

          <Note>
            Mailhog catches <em>every</em> email your app sends during development. This means you can test
            password reset flows, welcome emails, and notifications without worrying about accidentally
            emailing real users.
          </Note>

          <Challenge number={8} title="Test Email with Mailhog">
            <p>
              Open <Code>localhost:8025</Code> in your browser (Mailhog). Then register a new user in the
              web app at <Code>localhost:3000</Code>. Go back to Mailhog — did a welcome email arrive?
              Open it and inspect the HTML content.
            </p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Email Templates ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Email Templates</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit includes 4 pre-built HTML email templates that cover the most common transactional email
            needs. They{"'"}re clean, responsive, and ready to customize.
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">welcome</strong> — sent when a user registers a new account</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">password-reset</strong> — sent when a user requests a password reset link</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">verify-email</strong> — sent to confirm a user{"'"}s email address</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">notification</strong> — a generic template for custom notifications</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            All templates live in the <Code>apps/api/internal/mail/</Code> directory:
          </p>

          <CodeBlock filename="Project Structure">
{`apps/api/internal/mail/
├── mailer.go             # Send function, SMTP config
├── templates/
│   ├── welcome.html      # Welcome email template
│   ├── password-reset.html
│   ├── verify-email.html
│   └── notification.html`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Templates use Go{"'"}s built-in <Code>html/template</Code> package. Dynamic data is injected
            using double curly braces:
          </p>

          <CodeBlock filename="apps/api/internal/mail/templates/welcome.html">
{`<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: sans-serif; background: #f4f4f4; padding: 20px; }
        .container { max-width: 600px; margin: 0 auto; background: #fff; padding: 32px; border-radius: 8px; }
        .button { background: #6c5ce7; color: #fff; padding: 12px 24px; text-decoration: none; border-radius: 6px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Welcome, {{.Name}}!</h1>
        <p>Thanks for joining {{.AppName}}. We're excited to have you.</p>
        <p><a href="{{.DashboardURL}}" class="button">Go to Dashboard</a></p>
    </div>
</body>
</html>`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>{"{{.Name}}"}</Code> syntax is Go{"'"}s template language. When you send the email,
            you pass a map of data, and Go replaces each <Code>{"{{.Key}}"}</Code> with the corresponding value.
          </p>

          <Challenge number={9} title="Inspect a Template">
            <p>
              Open the welcome email template file. Find the placeholder for the user{"'"}s name.
              How does Go inject the user{"'"}s actual name? What other placeholders does the template use?
            </p>
          </Challenge>
        </section>

        {/* ═══ Section 11: Sending Email from Code ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Sending Email from Code</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Sending an email from your Go code is straightforward. Grit provides a <Code>mailer</Code> package
            with a <Code>Send</Code> function that handles template rendering and delivery.
          </p>

          <CodeBlock filename="Sending an Email">
{`err := mailer.Send(mail.Email{
    To:       "user@example.com",
    Subject:  "Welcome!",
    Template: "welcome",
    Data: map[string]string{
        "Name":         "John",
        "AppName":      "My App",
        "DashboardURL": "https://myapp.com/dashboard",
    },
})`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>Template</Code> field matches one of the HTML files in the <Code>templates/</Code> directory
            (without the <Code>.html</Code> extension). The <Code>Data</Code> map provides the values for
            all the <Code>{"{{.Key}}"}</Code> placeholders in the template.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            But remember — you should <strong className="text-foreground">never send emails directly from an API handler</strong>.
            Instead, dispatch a background job that sends the email. This way the user{"'"}s request isn{"'"}t
            blocked by the email service:
          </p>

          <CodeBlock filename="The Right Way — Dispatch as a Background Job">
{`// In your handler or service
payload, _ := json.Marshal(jobs.EmailPayload{
    To:       "user@example.com",
    Subject:  "Welcome!",
    Template: "welcome",
    Data:     map[string]string{"Name": "John"},
})
task := asynq.NewTask(jobs.TypeEmailSend, payload)
client.Enqueue(task)

// The email worker will pick this up and call mailer.Send()`}
          </CodeBlock>

          <Tip>
            Always send emails through the background job queue, not inline. This gives you automatic
            retries if the email service is down, and keeps your API responses fast.
          </Tip>

          <Challenge number={10} title="Find the Welcome Email Dispatch">
            <p>
              Look at the auth handler registration code in <Code>apps/api/internal/handlers/</Code>.
              Can you find where the welcome email is dispatched as a background job after a user registers?
              What task type does it use?
            </p>
          </Challenge>

          <Challenge number={11} title="Send a Test Email">
            <p>
              Using the mailer package, write the Go code to send a password reset email to
              {' '}<Code>test@example.com</Code> with the user{"'"}s name set to {'"'}Alice{'"'} and a reset
              link of <Code>https://myapp.com/reset?token=abc123</Code>. Which template would you use?
            </p>
          </Challenge>
        </section>

        {/* ═══ Section 12: Summary + Final Challenge ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-6">
            In this course, you learned how to offload slow or unreliable tasks to background workers using
            asynq and Redis, and how to send transactional emails with Go templates and Resend. Here{"'"}s
            what you covered:
          </p>

          <ul className="space-y-3 text-muted-foreground mb-8">
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Background jobs</strong> — tasks that run asynchronously outside the user{"'"}s request, keeping the API fast</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">asynq architecture</strong> — producers enqueue tasks to Redis, workers dequeue and process them</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Built-in workers</strong> — email, image processing, and cleanup workers come pre-configured</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Custom jobs</strong> — 4 steps: define task type, create handler, register it, dispatch from your code</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Retry logic</strong> — exponential backoff with configurable max retries and delayed processing</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Cron scheduler</strong> — recurring tasks on a schedule using cron expressions</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Admin Jobs dashboard</strong> — monitor active, pending, completed, and failed tasks from the admin panel</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Transactional email</strong> — user-triggered emails using Resend, with Go HTML templates and Mailhog for local dev</span>
            </li>
          </ul>

          <Challenge number={12} title="Final Challenge: Create a Weekly Digest">
            <p>
              Design a background job system that sends a weekly digest email to all users. Think through
              every piece you would need:
            </p>
            <ol className="mt-3 space-y-2 list-decimal list-inside">
              <li><strong className="text-foreground">Task type</strong> — What would you name it? (e.g., <Code>digest:weekly</Code>)</li>
              <li><strong className="text-foreground">Payload struct</strong> — What data does the worker need? (Think: does it need a user ID, or does it query all users?)</li>
              <li><strong className="text-foreground">Handler function</strong> — Write pseudocode: query new posts from the past 7 days, loop through all users, send each one a digest email</li>
              <li><strong className="text-foreground">Cron expression</strong> — What expression runs every Monday at 8:00 AM?</li>
              <li><strong className="text-foreground">Email template</strong> — What placeholders would it need? (Name, post titles, post links, date range)</li>
            </ol>
            <p className="mt-3">
              Write out the full task type constant, payload struct, handler skeleton, cron registration line,
              and a list of template placeholders.
            </p>
          </Challenge>

          <Challenge number={13} title="Bonus: Error Handling Strategy">
            <p>
              Think about what could go wrong with the weekly digest job. What if the database query fails
              halfway through? What if Resend is down for one user but works for others? Should you send
              all emails in one task, or enqueue a separate {'"'}send digest{'"'} task per user? What are the
              tradeoffs?
            </p>
          </Challenge>

          <Challenge number={14} title="Bonus: Monitoring Checklist">
            <p>
              Create a monitoring checklist for a production app using background jobs. What would you check
              daily? Consider: queue depth (are jobs piling up?), failure rate, retry counts, worker
              health, Redis memory usage, and email delivery rates. How would the admin Jobs dashboard
              help you catch problems early?
            </p>
          </Challenge>
        </section>

        {/* ═══ Footer ═══ */}
        <CourseFooter />

        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses/grit-web/file-storage', label: 'File Storage' }}
            next={{ href: '/courses/grit-web/ai-features', label: 'AI Features' }}
          />
        </div>
      </main>
    </div>
  )
}
