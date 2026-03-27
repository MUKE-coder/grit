import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Batteries Included: Every Feature That Ships with Grit',
  description: 'Learn every battery that ships with a Grit project — Redis caching, file storage, email, background jobs, cron, AI, 2FA, GORM Studio, API docs, security, observability, UI components, and Docker.',
}

export default function BatteriesIncludedCourse() {
  return (
    <div className="min-h-screen bg-[#0b1120]">
      <SiteHeader />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Batteries Included</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">15 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Batteries Included: Every Feature That Ships with Grit
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Every Grit project ships with 18 production-ready features pre-configured and wired together.
            This course walks through every single one — what it is, why it matters, how it{"'"}s configured,
            and how to use it. By the end, you{"'"}ll know your project inside and out.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What "Batteries Included" Means ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What {'"'}Batteries Included{'"'} Means</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Most frameworks give you a foundation and tell you to install plugins for everything else.
            Need caching? Find a Redis library, wire it up, write middleware. Need email? Pick a service,
            install an SDK, build templates. Need file uploads? Choose between S3, Cloudflare, or local
            storage, configure presigned URLs, handle multipart forms. Every feature is a research project.
          </p>

          <Definition term="Batteries Included">
            A design philosophy where the framework ships with everything you need to build a production
            application out of the box. No hunting for plugins, no configuration marathons, no compatibility
            headaches. Every feature is pre-configured, tested, and wired together. Python{"'"}s standard
            library and Ruby on Rails popularized this approach.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit ships <strong className="text-foreground">18 features</strong> pre-configured and wired
            together in every project. They work on day one — no setup required. And because they{"'"}re
            modular, you can remove any feature you don{"'"}t need without breaking the others.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s the full list of what ships with every Grit project:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">#</th>
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Battery</th>
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">What It Does</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                {[
                  ['1', 'Redis Caching', 'GET response caching with TTL'],
                  ['2', 'File Storage', 'S3-compatible uploads (MinIO/R2/S3)'],
                  ['3', 'Email Service', 'Transactional email with 4 HTML templates'],
                  ['4', 'Background Jobs', 'Redis-backed async task queue'],
                  ['5', 'Cron Scheduler', 'Recurring tasks with cron expressions'],
                  ['6', 'AI Integration', 'Vercel AI Gateway (hundreds of models)'],
                  ['7', 'Two-Factor Auth', 'TOTP + backup codes + trusted devices'],
                  ['8', 'GORM Studio', 'Visual database browser at /studio'],
                  ['9', 'API Docs', 'Auto-generated OpenAPI 3.1 at /docs'],
                  ['10', 'Sentinel', 'WAF, rate limiting, security dashboard'],
                  ['11', 'Pulse', 'Request tracing, DB monitoring, metrics'],
                  ['12', 'Grit UI', '100 shadcn-compatible components'],
                  ['13', 'Docker Setup', 'Dev + production compose files'],
                  ['14', 'JWT Auth', 'Access + refresh tokens, middleware'],
                  ['15', 'CORS Middleware', 'Pre-configured cross-origin requests'],
                  ['16', 'Gzip Compression', 'Automatic response compression'],
                  ['17', 'Request Logging', 'Structured request/response logs'],
                  ['18', 'Connection Pooling', 'Optimized DB connection management'],
                ].map(([num, name, desc]) => (
                  <tr key={num} className="border-b border-border/20">
                    <td className="px-3 py-2 text-muted-foreground">{num}</td>
                    <td className="px-3 py-2 font-medium text-foreground">{name}</td>
                    <td className="px-3 py-2">{desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <Tip>
            You don{"'"}t need to memorize all 18 features right now. This course walks through each one
            with examples and hands-on challenges. By the end, you{"'"}ll have used every single battery
            in your project.
          </Tip>

          <Challenge number={1} title="Count the Setup Time">
            <p>Look at the table above. Pick any 5 features and estimate how long it would take you to
            set up each one from scratch in a new Go project (finding the library, reading docs, writing
            code, testing). Add up the total hours. That{"'"}s the {'"'}Setup Tax{'"'} Grit eliminates.</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: Redis Caching ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Redis Caching</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Caching is one of the most impactful performance optimizations you can make. Instead of
            querying the database every time a user requests the same data, you store the result in
            memory and serve it instantly. Grit{"'"}s cache service wraps Redis with a clean Go API.
          </p>

          <Definition term="Redis">
            An in-memory data store that operates at sub-millisecond speed. It stores data as key-value
            pairs in RAM, making reads hundreds of times faster than a database query. Redis is the
            industry standard for caching, session storage, and job queues.
          </Definition>

          <Definition term="Cache">
            A temporary storage layer that holds frequently accessed data closer to the consumer. Instead
            of computing a result or querying a database every time, you store the result once and serve
            it from the cache until it expires.
          </Definition>

          <Definition term="TTL (Time To Live)">
            How long a cached value remains valid before it expires and gets deleted automatically. A TTL
            of 5 minutes means the cached data will be refreshed every 5 minutes. Shorter TTLs mean
            fresher data. Longer TTLs mean better performance.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit{"'"}s cache service exposes three methods — <Code>Set</Code>, <Code>Get</Code>,
            and <Code>Delete</Code>:
          </p>

          <CodeBlock filename="internal/service/cache_service.go">
{`// Store a value with a 5-minute TTL
cache.Set("user:1", userData, 5*time.Minute)

// Retrieve it — returns the cached value or an error if expired/missing
cached, err := cache.Get("user:1")

// Explicitly remove a cached value (e.g., after an update)
cache.Delete("user:1")`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            On top of the service, Grit includes <strong className="text-foreground">cache middleware</strong> that
            automatically caches GET responses. When a GET request comes in, the middleware checks Redis
            first. If the response is cached, it returns it instantly without hitting your handler. If
            not, it calls your handler, caches the response, and returns it.
          </p>

          <CodeBlock filename=".env">
{`# Redis configuration
REDIS_URL=redis://localhost:6379`}
          </CodeBlock>

          <Note>
            The cache middleware only caches GET requests. POST, PUT, PATCH, and DELETE requests always
            go through to your handler. When a write happens, the relevant cache keys are automatically
            invalidated so users always see fresh data.
          </Note>

          <Challenge number={2} title="Explore the Cache Service">
            <p>Open the file <Code>internal/service/cache_service.go</Code> in your Grit project. What
            methods does it expose? What parameters does <Code>Set</Code> accept? What does <Code>Get</Code> return
            when the key doesn{"'"}t exist?</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: File Storage ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">File Storage (S3-Compatible)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every real application needs file uploads — profile pictures, documents, images, attachments.
            Grit{"'"}s file storage service works with any S3-compatible provider: <strong className="text-foreground">MinIO</strong> for
            local development, <strong className="text-foreground">AWS S3</strong> for production,
            or <strong className="text-foreground">Cloudflare R2</strong> for edge storage.
          </p>

          <Definition term="S3-Compatible">
            Amazon S3 (Simple Storage Service) defined the standard API for object storage. Any service
            that implements this same API is called {'"'}S3-compatible.{'"'} This means you can switch
            between providers (MinIO, R2, Backblaze B2) without changing your code — just update the
            endpoint URL.
          </Definition>

          <Definition term="Presigned URL">
            A temporary URL that grants time-limited access to upload or download a private file. Instead
            of routing files through your server (which uses bandwidth and memory), the client uploads
            directly to the storage provider using a presigned URL. The URL expires after a set time,
            so files stay secure.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The upload flow in Grit works like this:
          </p>

          <CodeBlock filename="Upload Flow">
{`1. Client requests a presigned URL from your API
   POST /api/uploads/presign { filename: "photo.jpg" }

2. API generates a presigned URL (valid for 15 minutes)
   Response: { url: "https://storage.../photo.jpg?signature=...", key: "uploads/abc123.jpg" }

3. Client uploads the file directly to storage using the presigned URL
   PUT https://storage.../photo.jpg?signature=... (file data)

4. Client confirms the upload to your API
   POST /api/uploads/confirm { key: "uploads/abc123.jpg" }

5. API stores the file reference in the database`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Configuration uses environment variables:
          </p>

          <CodeBlock filename=".env">
{`# File storage (S3-compatible)
STORAGE_ENDPOINT=localhost:9000
STORAGE_ACCESS_KEY=minioadmin
STORAGE_SECRET_KEY=minioadmin
STORAGE_BUCKET=uploads
STORAGE_USE_SSL=false
STORAGE_REGION=us-east-1`}
          </CodeBlock>

          <Tip>
            In development, MinIO runs via Docker Compose. You can browse uploaded files at
            <Code>localhost:9001</Code> using the MinIO Console. The default credentials
            are <Code>minioadmin</Code> / <Code>minioadmin</Code>.
          </Tip>

          <Challenge number={3} title="Upload a File">
            <p>Start your Grit project with Docker Compose. Open the admin panel and upload a file
            (any image will do). Then open the MinIO Console at <Code>localhost:9001</Code>. Can
            you find the file in the <Code>uploads</Code> bucket? What is the full path of the file?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Email Service ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Email Service (Resend)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Transactional email is essential for any application — welcome emails, password resets,
            verification codes, notifications. Grit includes a complete email service powered
            by <strong className="text-foreground">Resend</strong> with 4 pre-built HTML templates.
          </p>

          <Definition term="Transactional Email">
            Emails triggered by user actions — not marketing blasts. When someone registers, they get a
            welcome email. When they reset their password, they get a reset link. These are transactional
            emails. They{"'"}re sent one-at-a-time in response to specific events.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit ships 4 HTML email templates:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Welcome</strong> — sent when a user registers</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Password Reset</strong> — sent when a user requests a password reset</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Email Verification</strong> — sent to verify a user{"'"}s email address</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Notification</strong> — a generic template for system notifications</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            In development, emails are captured by <strong className="text-foreground">Mailhog</strong> instead
            of being sent to real inboxes. This means you can test every email flow without configuring
            a real email provider.
          </p>

          <CodeBlock filename="internal/service/email_service.go">
{`// Send a welcome email to a new user
emailService.SendWelcome(user.Email, user.Name)

// Send a password reset link
emailService.SendPasswordReset(user.Email, resetToken)

// Send an email verification code
emailService.SendVerification(user.Email, verificationCode)

// Send a generic notification
emailService.SendNotification(user.Email, "Order Shipped", "Your order #1234 has shipped.")`}
          </CodeBlock>

          <CodeBlock filename=".env">
{`# Email (Resend)
RESEND_API_KEY=re_xxxxxxxxxxxx
EMAIL_FROM=noreply@myapp.com

# Local testing (Mailhog captures all emails)
SMTP_HOST=localhost
SMTP_PORT=1025`}
          </CodeBlock>

          <Tip>
            Mailhog runs automatically with Docker Compose. Open <Code>localhost:8025</Code> in your
            browser to see every email your application sends during development.
          </Tip>

          <Challenge number={4} title="Catch a Welcome Email">
            <p>Start your Grit project. Open Mailhog at <Code>localhost:8025</Code>. Now register a new
            user through the web app or API. Go back to Mailhog — do you see the welcome email? Open it
            and look at the HTML template. What information does it include?</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Background Jobs ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Background Jobs (asynq)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Some tasks are too slow or too resource-intensive to run inside an HTTP request. Sending an
            email takes 200-500ms. Processing an image takes 1-2 seconds. Generating a report takes 5-10
            seconds. You don{"'"}t want users waiting for these. Background jobs let you offload work to
            a separate process that runs asynchronously.
          </p>

          <Definition term="Background Job">
            A task that runs outside the request/response cycle. Your API handler enqueues the job
            (pushes it onto a queue) and immediately returns a response to the user. A separate worker
            process picks up the job and executes it in the background. The user doesn{"'"}t wait.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit uses <strong className="text-foreground">asynq</strong>, a Redis-backed job queue for Go.
            It ships with 3 built-in workers:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Email Worker</strong> — sends transactional emails asynchronously</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Image Worker</strong> — processes uploaded images (resize, thumbnails)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Cleanup Worker</strong> — removes expired tokens, old sessions, temp files</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Dispatching a job is a single function call:
          </p>

          <CodeBlock filename="Dispatching a Background Job">
{`// Enqueue an email job — returns immediately
err := jobs.Dispatch("email:welcome", map[string]interface{}{
    "email": user.Email,
    "name":  user.Name,
})

// Enqueue an image processing job
err := jobs.Dispatch("image:resize", map[string]interface{}{
    "file_key": "uploads/abc123.jpg",
    "width":    800,
    "height":   600,
})

// Enqueue with options: delay, retry, deadline
err := jobs.DispatchWithOptions("cleanup:tokens", nil, asynq.ProcessIn(1*time.Hour))`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The admin panel includes a <strong className="text-foreground">Jobs dashboard</strong> where
            you can monitor queued, active, completed, and failed jobs in real time.
          </p>

          <Challenge number={5} title="Inspect the Jobs Dashboard">
            <p>Open the admin panel and navigate to the Jobs page. Are there any queued or completed
            jobs? Try registering a new user — the welcome email is sent as a background job. Go back
            to the Jobs page. Do you see the email job in the completed list?</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Cron Scheduler ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Cron Scheduler</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Background jobs handle one-off tasks. Cron jobs handle <strong className="text-foreground">recurring
            tasks</strong> — things that need to happen on a schedule. Clean up expired tokens every hour.
            Send a weekly digest every Monday. Generate reports every night at midnight.
          </p>

          <Definition term="Cron Expression">
            A string that defines a schedule using 5 fields: minute, hour, day of month, month, day of
            week. For example, <Code>0 9 * * 1</Code> means {'"'}at 9:00 AM every Monday.{'"'} The
            asterisk (<Code>*</Code>) means {'"'}every{'"'} — so <Code>* * * * *</Code> means every
            minute.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit{"'"}s cron scheduler uses the same worker pool as background jobs, so you don{"'"}t need
            a separate process. Here{"'"}s the cron expression format:
          </p>

          <CodeBlock filename="Cron Expression Format">
{`┌───────────── minute (0-59)
│ ┌───────────── hour (0-23)
│ │ ┌───────────── day of month (1-31)
│ │ │ ┌───────────── month (1-12)
│ │ │ │ ┌───────────── day of week (0-6, Sunday=0)
│ │ │ │ │
* * * * *`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Common examples:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Expression</th>
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Schedule</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-mono text-foreground">* * * * *</td><td className="px-3 py-2">Every minute</td></tr>
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-mono text-foreground">0 * * * *</td><td className="px-3 py-2">Every hour (at minute 0)</td></tr>
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-mono text-foreground">0 0 * * *</td><td className="px-3 py-2">Every day at midnight</td></tr>
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-mono text-foreground">0 9 * * 1</td><td className="px-3 py-2">Every Monday at 9:00 AM</td></tr>
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-mono text-foreground">0 0 1 * *</td><td className="px-3 py-2">First day of every month at midnight</td></tr>
                <tr><td className="px-3 py-2 font-mono text-foreground">30 2 * * 0</td><td className="px-3 py-2">Every Sunday at 2:30 AM</td></tr>
              </tbody>
            </table>
          </div>

          <Challenge number={6} title="Write Cron Expressions">
            <p>Write the cron expression for each schedule:</p>
            <ul className="list-disc list-inside mt-2 space-y-1">
              <li>Every hour on the hour</li>
              <li>Every day at midnight</li>
              <li>Every Monday at 9:00 AM</li>
              <li>Every 15 minutes</li>
              <li>The first of every month at 3:00 AM</li>
            </ul>
            <p className="mt-2">Check your answers against the table above. The last two are not in the
            table — you{"'"}ll need to figure them out yourself.</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: AI Integration ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">AI Integration (Vercel AI Gateway)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every Grit project ships with AI integration powered by <strong className="text-foreground">Vercel
            AI Gateway</strong>. One API key gives you access to hundreds of models from dozens of
            providers — Claude, GPT, Gemini, Llama, Mistral, and more. No provider-specific code required.
          </p>

          <Definition term="AI Gateway">
            A proxy service that sits between your application and AI providers. You send a standardized
            request to the gateway, and it routes it to whichever provider and model you choose. Switch
            models by changing one environment variable — no code changes.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit exposes 3 AI endpoints:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">POST /api/ai/complete</strong> — single prompt in, single response out</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">POST /api/ai/chat</strong> — multi-turn conversation with message history</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">POST /api/ai/stream</strong> — streaming response (tokens arrive in real-time)</li>
          </ul>

          <CodeBlock filename=".env">
{`# AI Gateway configuration
AI_GATEWAY_URL=https://gateway.ai.vercel.app/v1
AI_GATEWAY_API_KEY=your-api-key-here
AI_GATEWAY_MODEL=anthropic/claude-sonnet-4-6`}
          </CodeBlock>

          <CodeBlock filename="Testing the Complete Endpoint">
{`curl -X POST localhost:8080/api/ai/complete \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \\
  -d '{"prompt": "Explain REST APIs in 3 sentences."}'`}
          </CodeBlock>

          <Tip>
            The AI service is protected by authentication. You need a valid JWT token to call these
            endpoints. This prevents unauthorized users from consuming your AI credits.
          </Tip>

          <Challenge number={7} title="Test AI Integration">
            <p>If you have an AI Gateway API key, add it to your <Code>.env</Code> file. Start the API,
            get a JWT token by logging in, then test the <Code>/api/ai/complete</Code> endpoint with
            a prompt of your choice. What model did the response come from? (Check the response headers
            or body for provider information.)</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Two-Factor Auth ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Two-Factor Authentication (TOTP)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Passwords alone are not enough. If someone steals a user{"'"}s password (phishing, data breach,
            reuse), they can access the account. Two-Factor Authentication (2FA) adds a second layer —
            even with the password, an attacker needs the user{"'"}s physical device.
          </p>

          <Definition term="TOTP (Time-based One-Time Password)">
            An algorithm that generates a 6-digit code that changes every 30 seconds. The code is
            derived from a shared secret key and the current time. Both the server and the authenticator
            app know the secret, so they generate the same code independently. No internet connection
            needed — it{"'"}s purely math-based.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit{"'"}s 2FA system includes:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Authenticator app support</strong> — works with Google Authenticator, Authy, 1Password, Bitwarden</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">10 backup codes</strong> — one-time-use codes in case the user loses their device</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Trusted devices</strong> — skip 2FA on recognized devices for 30 days</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Zero external dependencies</strong> — the TOTP algorithm runs in pure Go, no third-party services</li>
          </ul>

          <CodeBlock filename="2FA Setup Flow">
{`1. User enables 2FA
   POST /api/auth/2fa/enable
   → Returns: QR code (base64) + secret key + 10 backup codes

2. User scans QR code with authenticator app

3. User verifies with the 6-digit code from the app
   POST /api/auth/2fa/verify { "code": "123456" }
   → 2FA is now active

4. On next login, after password is verified:
   POST /api/auth/2fa/validate { "code": "789012" }
   → Returns JWT token`}
          </CodeBlock>

          <Note>
            Backup codes are hashed before storage (like passwords). When a user uses a backup code,
            it{"'"}s marked as used and can{"'"}t be used again. If all 10 codes are used, the user
            must generate new ones.
          </Note>

          <Challenge number={8} title="Enable 2FA">
            <p>Enable 2FA on your account. Scan the QR code with an authenticator app (Google Authenticator
            is free). Verify with a code. Then log out and log back in — you{"'"}ll be prompted for the
            6-digit code. Finally, test one of your backup codes: log out, and instead of the app code,
            enter a backup code. Does it work? Can you use the same backup code twice?</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: GORM Studio ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">GORM Studio</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            GORM Studio is a visual database browser built directly into your Grit project. It runs
            at <Code>/studio</Code> and lets you browse tables, view and edit records, run raw SQL
            queries, and export schemas — all from your browser.
          </p>

          <Definition term="Database Browser">
            A graphical interface for inspecting and manipulating database contents. Instead of writing
            SQL in a terminal, you see tables, rows, and columns in a visual interface. Think of it
            like phpMyAdmin for PHP, but built specifically for GORM and embedded in your Go API.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            What you can do in GORM Studio:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Browse tables</strong> — see all tables, columns, types, and relationships</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">View and edit data</strong> — click any row to view details, edit fields inline</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Run SQL queries</strong> — execute raw SQL and see results in a table</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Export schemas</strong> — download your database schema as SQL</li>
          </ul>

          <CodeBlock filename="Accessing GORM Studio">
{`# Start your API
cd apps/api && go run cmd/server/main.go

# Open GORM Studio in your browser
http://localhost:8080/studio`}
          </CodeBlock>

          <Tip>
            GORM Studio is only available in development mode. In production, the <Code>/studio</Code> route
            is disabled for security. You don{"'"}t want random users browsing your database.
          </Tip>

          <Challenge number={9} title="Explore GORM Studio">
            <p>Open GORM Studio at <Code>localhost:8080/studio</Code>. Find at least 3 tables and
            describe what each one stores. How many columns does the <Code>users</Code> table have?
            Try running a SQL query: <Code>SELECT COUNT(*) FROM users;</Code></p>
          </Challenge>
        </section>

        {/* ═══ Section 10: API Documentation ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">API Documentation (Auto-Generated)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit automatically generates an <strong className="text-foreground">OpenAPI 3.1 specification</strong> for
            your API and serves an interactive documentation UI at <Code>/docs</Code>. Every endpoint,
            every request body, every response shape is documented — and it stays in sync with your
            code because it{"'"}s generated from your route definitions.
          </p>

          <Definition term="OpenAPI Specification">
            A standard format for describing REST APIs. It defines endpoints, request/response schemas,
            authentication methods, and more in a machine-readable JSON or YAML format. Tools like
            Swagger UI, Redoc, and Postman can read OpenAPI specs to generate interactive documentation.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            No annotations or comments needed. Grit reads your Gin route definitions and GORM models
            to generate the spec automatically. When you add a new resource with <Code>grit generate</Code>,
            the docs update automatically.
          </p>

          <CodeBlock filename="Accessing API Docs">
{`# Start your API
cd apps/api && go run cmd/server/main.go

# Interactive docs
http://localhost:8080/docs

# Raw OpenAPI 3.1 JSON spec
http://localhost:8080/docs/openapi.json`}
          </CodeBlock>

          <Challenge number={10} title="Test an Endpoint">
            <p>Open <Code>localhost:8080/docs</Code>. Find the login endpoint. Use the interactive UI
            to send a login request with the test credentials. Did you get a JWT token back? Now find
            the {'"'}list users{'"'} endpoint — does it require authentication? Try calling it without
            a token and see what error you get.</p>
          </Challenge>
        </section>

        {/* ═══ Section 11: Security (Sentinel) ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Security (Sentinel)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Sentinel is Grit{"'"}s built-in security layer. It protects your API from common attacks
            without you writing a single line of security code. It includes a WAF (Web Application
            Firewall), rate limiting, brute-force protection, security headers, and a threat dashboard.
          </p>

          <Definition term="WAF (Web Application Firewall)">
            A firewall that inspects HTTP requests for malicious patterns — SQL injection, cross-site
            scripting (XSS), path traversal, and other attack vectors. It sits between the internet
            and your application, blocking bad requests before they reach your handlers.
          </Definition>

          <Definition term="Rate Limiting">
            Restricting how many requests a client can make in a given time period. For example, 100
            requests per minute per IP address. This prevents abuse, brute-force attacks, and
            denial-of-service attempts. Exceeding the limit returns a <Code>429 Too Many Requests</Code> response.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            What Sentinel protects against:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">SQL Injection</strong> — malicious SQL in query parameters or request bodies</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">XSS (Cross-Site Scripting)</strong> — injecting JavaScript into your pages</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Brute-force login</strong> — repeated login attempts are throttled and blocked</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Path traversal</strong> — attempts to access files outside allowed directories</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Missing security headers</strong> — auto-adds CSP, HSTS, X-Frame-Options, and more</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The Sentinel dashboard at <Code>/sentinel/ui</Code> shows real-time threat data: blocked
            requests, top attacking IPs, attack types, and rate limit violations.
          </p>

          <Challenge number={11} title="Visit the Security Dashboard">
            <p>Open <Code>localhost:8080/sentinel/ui</Code> in your browser. What information does the
            dashboard show? Make 50 rapid requests to any endpoint (use a loop in your terminal). Does
            the rate limiter kick in? Check the dashboard again for the rate limit data.</p>
          </Challenge>
        </section>

        {/* ═══ Section 12: Observability (Pulse) ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Observability (Pulse)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You can{"'"}t fix what you can{"'"}t see. Pulse is Grit{"'"}s built-in observability layer. It
            tracks every request, monitors database performance, collects metrics, and runs health
            checks — all viewable from a dashboard at <Code>/pulse/ui</Code>.
          </p>

          <Definition term="Observability">
            The ability to understand what{"'"}s happening inside your application by examining its
            outputs — logs, metrics, and traces. Observability answers questions like: {'"'}Why is this
            endpoint slow?{'"'} {'"'}How many requests per second am I handling?{'"'} {'"'}Is the database
            the bottleneck?{'"'}
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            What Pulse tracks:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Request tracing</strong> — every request is logged with timing, status code, and unique request ID</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Database monitoring</strong> — query count, average duration, slowest queries</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Metrics</strong> — requests per second, error rate, response time percentiles (p50, p95, p99)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Health checks</strong> — database connectivity, Redis connectivity, disk space, memory usage</li>
          </ul>

          <CodeBlock filename="Accessing Pulse">
{`# Pulse dashboard
http://localhost:8080/pulse/ui

# Health check endpoint (JSON)
http://localhost:8080/pulse/health

# Metrics endpoint
http://localhost:8080/pulse/metrics`}
          </CodeBlock>

          <Challenge number={12} title="Generate Some Metrics">
            <p>Make at least 10 API requests — a mix of GET, POST, and invalid requests (wrong URLs,
            missing auth). Then open Pulse at <Code>localhost:8080/pulse/ui</Code>. What is the average
            response time? Which endpoint is the slowest? How many errors did you generate?</p>
          </Challenge>
        </section>

        {/* ═══ Section 13: Grit UI Components ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Grit UI Components</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every Grit project ships with access to <strong className="text-foreground">100 pre-built
            UI components</strong> in the shadcn/ui format. These aren{"'"}t npm packages — they{"'"}re
            source code files you install directly into your project, giving you full control to
            customize them.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The components are organized into 5 categories:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Category</th>
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Count</th>
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Examples</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-medium text-foreground">Marketing</td><td className="px-3 py-2">21</td><td className="px-3 py-2">Hero sections, feature grids, pricing tables, testimonials</td></tr>
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-medium text-foreground">Auth</td><td className="px-3 py-2">10</td><td className="px-3 py-2">Login forms, register forms, forgot password, 2FA pages</td></tr>
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-medium text-foreground">SaaS</td><td className="px-3 py-2">30</td><td className="px-3 py-2">Dashboard layouts, billing cards, settings pages, onboarding</td></tr>
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-medium text-foreground">E-commerce</td><td className="px-3 py-2">20</td><td className="px-3 py-2">Product cards, cart, checkout, order history, wishlists</td></tr>
                <tr><td className="px-3 py-2 font-medium text-foreground">Layout</td><td className="px-3 py-2">20</td><td className="px-3 py-2">Navbars, sidebars, footers, app shells, mobile menus</td></tr>
              </tbody>
            </table>
          </div>

          <CodeBlock filename="Browsing & Installing Components">
{`# Browse all components in your browser
http://localhost:3000/components

# The registry API serves component metadata
GET /api/r.json          → lists all 100 components
GET /api/r/hero-split.json  → metadata + source for a specific component

# Install a component (shadcn CLI compatible)
npx shadcn@latest add "http://localhost:8080/r/hero-split.json"`}
          </CodeBlock>

          <Challenge number={13} title="Browse the Registry">
            <p>Open the component browser at <Code>localhost:3000/components</Code> (or call the API
            at <Code>localhost:8080/api/r.json</Code>). How many components are in the {'"'}saas{'"'} category?
            Pick one component and explain how you would install it into your project.</p>
          </Challenge>
        </section>

        {/* ═══ Section 14: Docker Setup ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Docker Setup</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit ships with two Docker Compose files: one for development and one for production. The
            development compose file starts all the infrastructure services your project needs with a
            single command.
          </p>

          <Definition term="Docker Compose">
            A tool for defining and running multi-container applications. Instead of starting each
            service manually (database, cache, storage, email), you define them all in a YAML file and
            start everything with <Code>docker compose up</Code>. Each service runs in its own isolated
            container.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The development compose file includes:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">PostgreSQL</strong> — the primary database (port 5432)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Redis</strong> — caching and job queue (port 6379)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">MinIO</strong> — S3-compatible file storage (port 9000, console on 9001)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Mailhog</strong> — email capture for testing (SMTP on 1025, UI on 8025)</li>
          </ul>

          <CodeBlock filename="Docker Commands">
{`# Start all development services
docker compose up -d

# Check running services
docker compose ps

# View logs for a specific service
docker compose logs -f postgres

# Stop everything
docker compose down

# Stop and remove all data (fresh start)
docker compose down -v`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The production compose file (<Code>docker-compose.prod.yml</Code>) is optimized for
            deployment — it builds your Go API into a minimal Docker image, uses environment variables
            for configuration, and includes health checks.
          </p>

          <Challenge number={14} title="List Running Services">
            <p>Run <Code>docker compose up -d</Code> to start all services. Then
            run <Code>docker compose ps</Code>. List every running service, its port mapping, and
            its status. Visit each service{"'"}s UI: PostgreSQL has no UI, but check MinIO
            at <Code>localhost:9001</Code> and Mailhog at <Code>localhost:8025</Code>.</p>
          </Challenge>
        </section>

        {/* ═══ Section 15: Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You{"'"}ve now toured every battery that ships with a Grit project. Here{"'"}s the complete list:
          </p>

          <ul className="space-y-3 mb-6">
            {[
              'Redis Caching — sub-millisecond response caching with automatic invalidation',
              'File Storage — S3-compatible uploads with presigned URLs (MinIO, R2, S3)',
              'Email Service — 4 HTML templates with Resend (Mailhog for local testing)',
              'Background Jobs — Redis-backed async task queue with 3 built-in workers',
              'Cron Scheduler — recurring tasks with cron expressions, same worker pool',
              'AI Integration — Vercel AI Gateway with hundreds of models, 3 endpoints',
              'Two-Factor Auth — TOTP with authenticator apps, backup codes, trusted devices',
              'GORM Studio — visual database browser at /studio',
              'API Documentation — auto-generated OpenAPI 3.1 spec at /docs',
              'Sentinel — WAF, rate limiting, brute-force protection, security dashboard',
              'Pulse — request tracing, DB monitoring, metrics, health checks',
              'Grit UI — 100 shadcn-compatible components across 5 categories',
              'Docker Setup — dev compose (4 services) + production compose',
              'JWT Authentication — access + refresh tokens with middleware',
              'CORS Middleware — pre-configured cross-origin request handling',
              'Gzip Compression — automatic response compression',
              'Request Logging — structured logs with request IDs',
              'Connection Pooling — optimized database connection management',
            ].map((item, i) => (
              <li key={i} className="flex gap-3 text-muted-foreground">
                <CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" />
                <span>{item}</span>
              </li>
            ))}
          </ul>

          <Challenge number={15} title="Tour Every Tool">
            <p>Visit each of these URLs in your running Grit project and confirm they work:</p>
            <ul className="list-disc list-inside mt-2 space-y-1">
              <li><Code>localhost:8080/studio</Code> — GORM Studio</li>
              <li><Code>localhost:8080/docs</Code> — API Documentation</li>
              <li><Code>localhost:8080/pulse/ui</Code> — Pulse Observability</li>
              <li><Code>localhost:8080/sentinel/ui</Code> — Sentinel Security</li>
              <li><Code>localhost:8025</Code> — Mailhog (email capture)</li>
              <li><Code>localhost:9001</Code> — MinIO Console (file storage)</li>
            </ul>
            <p className="mt-2">Take a screenshot of each. You now know every tool in your Grit project.
            These aren{"'"}t plugins you installed — they shipped with your project from the very first
            command.</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses', label: 'All Courses' }}
            next={{ href: '/courses', label: 'More Courses' }}
          />
        </div>
      </main>
    </div>
  )
}
