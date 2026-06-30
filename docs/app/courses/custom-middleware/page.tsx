import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Custom Middleware & Hooks: Extending Grit',
  description: 'Learn to write custom Gin middleware and GORM hooks for Grit projects. Request timing, API key auth, maintenance mode, error handling, and database lifecycle hooks.',
}

export default function CustomMiddlewareCourse() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Custom Middleware & Hooks</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">12 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Custom Middleware & Hooks: Extending Grit
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Learn to extend your Grit API with custom middleware and database hooks. You{"'"}ll build
            request timers, API key authentication, maintenance mode, error recovery, and GORM
            lifecycle hooks — the building blocks for production-grade APIs.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What is Middleware? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is Middleware?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every request to your API passes through a chain of middleware before reaching the handler.
            Think of middleware as security guards at a building entrance — each one checks something
            different (ID badge, temperature, visitor log) before letting you in.
          </p>

          <Definition term="Middleware">
            Code that runs BEFORE (and optionally AFTER) your route handler. Middleware can inspect
            requests, modify responses, short-circuit the chain, or pass control to the next
            middleware. It{"'"}s how you add cross-cutting concerns like auth, logging, and rate
            limiting without repeating code in every handler.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The request lifecycle with middleware:
          </p>

          <CodeBlock filename="Middleware Chain">
{`Request arrives
  -> Logger middleware (logs method, path, start time)
  -> CORS middleware (adds cross-origin headers)
  -> Auth middleware (verifies JWT token)
  -> RateLimit middleware (checks request count)
  -> Your Handler (processes the request)
  <- RateLimit middleware (nothing to do after)
  <- Auth middleware (nothing to do after)
  <- CORS middleware (nothing to do after)
  <- Logger middleware (logs response time, status code)
Response sent`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Notice that middleware runs in order going in, and in reverse order coming out. Code
            before <Code>c.Next()</Code> runs on the way in. Code after <Code>c.Next()</Code> runs
            on the way out.
          </p>

          <Challenge number={1} title="Name 3 Middleware Examples">
            <p>Name 3 types of middleware you use every day (even if you didn{"'"}t realize they were
            middleware). Hint: think about authentication, cross-origin requests, request logging,
            compression, rate limiting, and CSRF protection. What does each one check or do?</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: How Gin Middleware Works ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">How Gin Middleware Works</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            A Gin middleware is a function that returns a <Code>gin.HandlerFunc</Code>. Inside, you
            receive <Code>*gin.Context</Code> and call <Code>c.Next()</Code> to pass control to the
            next middleware (or the handler).
          </p>

          <CodeBlock filename="middleware/timer.go">
{`// A simple middleware that logs request timing
// func MyMiddleware() gin.HandlerFunc
//     return func(c *gin.Context)
//
//         // --- BEFORE the handler ---
//         fmt.Println("Request started:", c.Request.URL.Path)
//
//         c.Next() // Call next middleware or handler
//
//         // --- AFTER the handler ---
//         fmt.Println("Response status:", c.Writer.Status())`}
          </CodeBlock>

          <Definition term="gin.HandlerFunc">
            The function signature for all Gin middleware and route handlers. It takes a single
            parameter: <Code>*gin.Context</Code>, which contains the request, response writer,
            URL parameters, headers, and methods to read/write data.
          </Definition>

          <Definition term="c.Next()">
            Calls the next middleware in the chain. Code before <Code>c.Next()</Code> runs before
            the handler processes the request. Code after <Code>c.Next()</Code> runs after the
            handler has written the response. If you don{"'"}t call <Code>c.Next()</Code>, the chain stops.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            To use middleware, register it with <Code>r.Use()</Code>:
          </p>

          <CodeBlock filename="Registering Middleware">
{`// Apply to ALL routes
r := gin.Default()
r.Use(MyMiddleware())

// Apply to a specific group
api := r.Group("/api")
api.Use(MyMiddleware())`}
          </CodeBlock>

          <Challenge number={2} title="Write a Simple Middleware">
            <p>Write a middleware function that prints the HTTP method and URL path for every request.
            Example output: <Code>GET /api/users</Code>, <Code>POST /api/posts</Code>. Where would
            you put the print statement — before or after <Code>c.Next()</Code>?</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Built-in Grit Middleware ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Built-in Grit Middleware</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit scaffolds several middleware out of the box. Each one lives in its own file in
            <Code>internal/middleware/</Code>:
          </p>

          <ul className="space-y-3 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Auth</strong> — extracts JWT from the Authorization header, validates it, and sets the user in context. Protected routes use this.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">CORS</strong> — adds Access-Control-Allow-Origin headers so the frontend (different port) can call the API.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Logger</strong> — structured request logging with method, path, status, and duration.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Cache</strong> — caches GET responses in Redis. Subsequent identical requests return the cached response.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">RequestID</strong> — adds an <Code>X-Request-ID</Code> header to every response for tracing.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Gzip</strong> — compresses responses to reduce bandwidth.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">RateLimit</strong> — limits requests per IP using a sliding window. Prevents abuse.</li>
          </ul>

          <Tip>
            Read the Auth middleware source code. It{"'"}s a great example of a real middleware: extract
            the token from the header, validate it, set the user in context with <Code>c.Set()</Code>,
            and call <Code>c.Next()</Code>. If the token is missing or invalid, it aborts with 401.
          </Tip>

          <Challenge number={3} title="Explore Built-in Middleware">
            <p>Find each middleware file in <Code>internal/middleware/</Code>. Read the Auth middleware
            carefully. How does it extract the token? What header does it look for? What happens
            when the token is missing? What does it store in the context for handlers to use?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Creating a Request Timer ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Creating a Request Timer</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Your first custom middleware: measure how long each request takes. This is useful for
            performance monitoring — if a request takes 2 seconds, you know something{"'"}s wrong.
          </p>

          <CodeBlock filename="middleware/timer.go">
{`// RequestTimer measures request duration
// func RequestTimer() gin.HandlerFunc
//     return func(c *gin.Context)
//         start := time.Now()
//
//         c.Next()  // Wait for handler to finish
//
//         duration := time.Since(start)
//         c.Header("X-Response-Time", duration.String())
//         log.Printf("%s %s took %v",
//             c.Request.Method,
//             c.Request.URL.Path,
//             duration,
//         )`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This middleware does two things after the handler runs:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> Adds an <Code>X-Response-Time</Code> header to the response — clients can read this</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Logs the method, path, and duration — you can see slow requests in your logs</li>
          </ul>

          <Note>
            The key insight: <Code>time.Now()</Code> is called BEFORE <Code>c.Next()</Code>, and
            <Code>time.Since(start)</Code> is called AFTER. The difference is how long the handler
            (and all subsequent middleware) took to run.
          </Note>

          <Challenge number={4} title="Add Request Timer">
            <p>Create a <Code>middleware/timer.go</Code> file with the RequestTimer middleware. Register
            it in your router. Make a few API requests and check: (1) Does the <Code>X-Response-Time</Code>
            header appear in responses? (2) Are the timings logged to the console?</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Creating an API Key Middleware ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Creating an API Key Middleware</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            JWT auth works great for user-facing apps, but what about external API consumers?
            Third-party services, webhooks, and server-to-server communication often use API keys
            instead of JWT tokens. An API key is a simple string passed in a header.
          </p>

          <CodeBlock filename="middleware/apikey.go">
{`// RequireAPIKey checks the X-API-Key header
// func RequireAPIKey() gin.HandlerFunc
//     return func(c *gin.Context)
//         key := c.GetHeader("X-API-Key")
//
//         if key == "" or key is not valid
//             c.AbortWithStatusJSON(401, error response)
//             // "code": "INVALID_API_KEY"
//             // "message": "Valid API key required"
//             return
//
//         c.Next()`}
          </CodeBlock>

          <Definition term="c.Abort()">
            Stops the middleware chain immediately. No further middleware or handlers will run. The
            response is sent as-is. Use <Code>c.AbortWithStatusJSON()</Code> to stop the chain AND
            send an error response in one call. This is how middleware rejects bad requests.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The pattern: check the header, validate it, and either continue (<Code>c.Next()</Code>)
            or reject (<Code>c.Abort()</Code>). Every auth-style middleware follows this exact pattern.
          </p>

          <Tip>
            Store valid API keys in the database or environment variables. Never hardcode them.
            For a simple setup, use an environment variable: <Code>API_KEY=sk_live_abc123</Code>.
            For multi-tenant APIs, store keys in a database table with the associated user/org.
          </Tip>

          <Challenge number={5} title="Create API Key Middleware">
            <p>Create the RequireAPIKey middleware. Test it by making requests with and without the
            <Code>X-API-Key</Code> header. Without the header, you should get a 401 response.
            With a valid key, the request should proceed normally.</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Creating a Maintenance Middleware ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Creating a Maintenance Middleware</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Sometimes you need to take the API offline for maintenance — database migrations,
            infrastructure changes, or emergency fixes. The maintenance middleware checks for a
            <Code>.maintenance</Code> file and returns 503 if it exists.
          </p>

          <CodeBlock filename="middleware/maintenance.go">
{`// Maintenance checks if the API is in maintenance mode
// func Maintenance() gin.HandlerFunc
//     return func(c *gin.Context)
//         // Check if .maintenance file exists
//         if file ".maintenance" exists
//             c.AbortWithStatusJSON(503, error response)
//             // "code": "MAINTENANCE"
//             // "message": "Service temporarily unavailable"
//             return
//
//         c.Next()`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This is how <Code>grit down</Code> works internally — it creates the <Code>.maintenance</Code>
            file, and the middleware blocks all requests. <Code>grit up</Code> removes the file and
            the API resumes instantly.
          </p>

          <CodeBlock filename="Maintenance Mode Commands">
{`# Take the API offline
touch .maintenance
# All requests now return 503

# Bring it back online
rm .maintenance
# Requests resume immediately`}
          </CodeBlock>

          <Note>
            No restart required. The middleware checks the file on every request. Create the file
            and the API is down. Delete the file and it{"'"}s back. This is far better than stopping
            the server, because the process stays running and health checks still work.
          </Note>

          <Challenge number={6} title="Test Maintenance Mode">
            <p>Read the actual maintenance middleware in your project. Compare it to the simplified
            version above. Create the <Code>.maintenance</Code> file and make an API request — do you
            get a 503? Remove the file and try again — does the API respond normally?</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Middleware Groups ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Middleware Groups</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Not every route needs every middleware. Public routes don{"'"}t need auth. Admin routes
            need auth AND role checks. Gin{"'"}s route groups let you apply different middleware
            to different sets of routes.
          </p>

          <CodeBlock filename="Route Groups with Middleware">
{`// Public routes - no auth required
public := r.Group("/api/auth")
// POST /api/auth/register
// POST /api/auth/login
// POST /api/auth/forgot-password

// Protected routes - auth required
protected := r.Group("/api")
// Uses Auth middleware to verify JWT
// GET /api/posts, GET /api/users, etc.

// Admin routes - auth + admin role required
admin := protected.Group("/admin")
// Uses RequireRole("ADMIN") middleware
// GET /admin/users, DELETE /admin/posts, etc.`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The nesting is important: the admin group inherits the auth middleware from the protected
            group, then adds its own role check. A request to <Code>/api/admin/users</Code> goes through:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">1.</span> Global middleware (logger, CORS, gzip)</li>
            <li className="flex gap-2"><span className="text-primary">2.</span> Auth middleware (verify JWT, set user in context)</li>
            <li className="flex gap-2"><span className="text-primary">3.</span> RequireRole({'"'}ADMIN{'"'}) middleware (check user role)</li>
            <li className="flex gap-2"><span className="text-primary">4.</span> Handler (list users)</li>
          </ul>

          <Challenge number={7} title="Create a Custom Route Group">
            <p>Create a new route group at <Code>/api/external</Code> that uses the API key middleware
            instead of JWT auth. Add a test endpoint to this group. Verify: JWT-protected routes
            reject API keys, and API-key-protected routes reject JWT tokens. Different auth for
            different consumers.</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: GORM Hooks ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">GORM Hooks</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Middleware handles HTTP concerns. Hooks handle database concerns. GORM hooks are methods
            on your models that run automatically during database operations — before a record is
            created, after it{"'"}s updated, before it{"'"}s deleted.
          </p>

          <Definition term="Hook">
            Code that runs automatically at a specific lifecycle point. GORM hooks attach to database
            operations: BeforeCreate, AfterCreate, BeforeUpdate, AfterUpdate, BeforeDelete, AfterDelete.
            They{"'"}re defined as methods on your model structs.
          </Definition>

          <CodeBlock filename="model/user.go - Hooks">
{`// BeforeCreate runs before inserting a new user
// func (u *User) BeforeCreate(tx *gorm.DB) error
//     u.Role = "USER"  // Default role for new users
//     return nil

// BeforeCreate runs before inserting a new post
// func (p *Post) BeforeCreate(tx *gorm.DB) error
//     p.Slug = slugify(p.Title)  // Auto-generate slug from title
//     return nil

// BeforeUpdate runs before updating a post
// func (p *Post) BeforeUpdate(tx *gorm.DB) error
//     if title changed
//         p.Slug = slugify(p.Title)  // Regenerate slug
//     return nil`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Common uses for GORM hooks:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">BeforeCreate</strong> — set defaults, generate slugs, hash passwords, validate data</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">AfterCreate</strong> — send welcome email, create related records, log audit trail</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">BeforeUpdate</strong> — regenerate slugs, validate changes, track modifications</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">BeforeDelete</strong> — soft-delete instead of hard-delete, cascade to related records</li>
          </ul>

          <Tip>
            If a hook returns an error, the database operation is cancelled and the error is returned
            to the caller. This is how you enforce business rules at the data layer — even if the
            handler forgets to validate, the hook catches it.
          </Tip>

          <Challenge number={8} title="Add a BeforeCreate Hook">
            <p>Add a BeforeCreate hook to the User model that validates email format. If the email
            doesn{"'"}t contain an @ symbol, return an error. Test it: try creating a user with an
            invalid email via the API. Does the hook reject it?</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Error Handling Middleware ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Error Handling Middleware</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            What happens when a handler panics? Without error handling middleware, the entire server
            crashes. With it, the panic is caught, logged, and a clean 500 response is returned.
            The server stays running.
          </p>

          <CodeBlock filename="middleware/recovery.go">
{`// ErrorHandler recovers from panics gracefully
// func ErrorHandler() gin.HandlerFunc
//     return func(c *gin.Context)
//         defer func()
//             if err := recover(); err != nil
//                 log.Printf("Panic recovered: %v", err)
//                 // Log stack trace for debugging
//                 c.AbortWithStatusJSON(500, error response)
//                 // "code": "INTERNAL_ERROR"
//                 // "message": "Something went wrong"
//
//         c.Next()`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>defer</Code> + <Code>recover()</Code> pattern is Go{"'"}s way of catching panics.
            The deferred function runs after <Code>c.Next()</Code> returns (or panics), and
            <Code>recover()</Code> catches the panic value if one occurred.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This middleware should be registered FIRST in the chain — before any other middleware.
            That way, it catches panics from any middleware or handler in the entire chain.
          </p>

          <Note>
            Gin{"'"}s default engine (<Code>gin.Default()</Code>) includes a recovery middleware. But
            it returns HTML error pages, not JSON. For an API, you want JSON error responses — which
            is why Grit scaffolds a custom error handler.
          </Note>

          <Challenge number={9} title="Test Error Recovery">
            <p>Add the ErrorHandler middleware to your router (register it first, before all other
            middleware). Create a test endpoint that intentionally panics. Hit it with a request.
            Does the middleware catch the panic? Does the server stay running? Check: is the response
            a clean JSON error (not an HTML page)?</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s everything you learned in this course:
          </p>

          <ul className="space-y-3 mb-6">
            {[
              'Middleware runs before (and after) your handlers in a chain',
              'gin.HandlerFunc takes *gin.Context — call c.Next() to continue, c.Abort() to stop',
              'Grit scaffolds 7 middleware: Auth, CORS, Logger, Cache, RequestID, Gzip, RateLimit',
              'Code before c.Next() runs on the way in, code after runs on the way out',
              'c.Abort() stops the chain and sends the response immediately',
              'Route groups apply different middleware to different sets of routes',
              'GORM hooks (BeforeCreate, AfterCreate, etc.) run on database operations',
              'Hooks can set defaults, validate data, generate slugs, and enforce business rules',
              'Error handling middleware with recover() keeps the server running after panics',
              'Register error handling middleware first so it catches panics from the entire chain',
            ].map((item, i) => (
              <li key={i} className="flex gap-3 text-muted-foreground">
                <CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" />
                <span>{item}</span>
              </li>
            ))}
          </ul>

          <Challenge number={10} title="Build a Request Logger">
            <p>Build a middleware that writes request logs to a file (not just the console). Each line
            should include: timestamp, method, path, status code, response time, and client IP. Use
            Go{"'"}s <Code>os.OpenFile</Code> with append mode. After 10 requests, open the log file
            and verify the entries.</p>
          </Challenge>

          <Challenge number={11} title="Build an IP Whitelist">
            <p>Build a middleware that only allows requests from specific IP addresses. Store the
            whitelist in an environment variable: <Code>ALLOWED_IPS=127.0.0.1,10.0.0.1</Code>. Requests
            from other IPs get a 403 Forbidden response. Test with your local IP — does it pass?</p>
          </Challenge>

          <Challenge number={12} title="Build a Slow Request Tracker">
            <p>Build a middleware that logs slow requests as warnings. If a request takes longer than
            500ms, log it with a {'"'}SLOW{'"'} prefix: <Code>SLOW: GET /api/posts took 1.2s</Code>.
            Combine this with the request timer from earlier. Create an intentionally slow handler
            (add a sleep) to test it.</p>
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
