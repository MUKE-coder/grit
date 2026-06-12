import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'
import { PrereqLinks } from '@/components/course/prereq-links'
import { PlaygroundChallenge } from '@/components/course/playground-challenge'

export default function Lesson() {
  return (
    <>
      <p>
        Every HTTP request that hits a Grit API passes through{' '}
        <strong>Gin</strong>. Before you can write or change a route,
        you need a feel for what Gin is, what it does for you, and what
        you have to do yourself.
      </p>

      <PrereqLinks
        prereqs={['golang']}
        intro={
          <>
            If you&apos;re new to Go syntax (structs, methods, slices,
            error handling), skim the Go primer first — this lesson
            assumes you can read it.
          </>
        }
      />

      <h2>What is Gin?</h2>
      <p>
        Gin is a small, fast HTTP web framework for Go. It does three
        jobs:
      </p>
      <ul>
        <li>
          <strong>Routes</strong> URLs to functions you write
          (<code>GET /api/users → ListUsers</code>).
        </li>
        <li>
          <strong>Parses</strong> the request (path params, query, body,
          headers) into Go values.
        </li>
        <li>
          <strong>Serializes</strong> your response (struct → JSON,
          with the right status code).
        </li>
      </ul>
      <p>
        Without Gin you would write{' '}
        <code>net/http</code> by hand and parse JSON manually for every
        endpoint. Gin removes 90% of that ceremony while staying
        thin enough to read in an afternoon.
      </p>

      <h2>The request lifecycle</h2>
      <Diagram label="A Gin request" caption="Gin matches the URL to a handler, runs middleware in order, then your handler, then writes the response.">
{`   HTTP request
        │
        ▼
   ┌─────────────────────────────────────────────┐
   │  Gin router  →  finds handler for URL       │
   └─────────────────────────────────────────────┘
        │
        ▼
   ┌─────────────────────────────────────────────┐
   │  Global middleware  (logger, CORS, recovery)│
   └─────────────────────────────────────────────┘
        │
        ▼
   ┌─────────────────────────────────────────────┐
   │  Group middleware  (e.g., Auth on /api/*)   │
   └─────────────────────────────────────────────┘
        │
        ▼
   ┌─────────────────────────────────────────────┐
   │  YOUR HANDLER                               │
   │  - reads input via c.Param / c.Query / c... │
   │  - calls service                            │
   │  - writes response via c.JSON               │
   └─────────────────────────────────────────────┘
        │
        ▼
   HTTP response`}
      </Diagram>

      <h2>The smallest possible Gin server</h2>
      <CodeBlock
        language="go"
        code={`package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()                   // 1. new router with sensible defaults
    r.GET("/ping", func(c *gin.Context) { // 2. register a route
        c.JSON(200, gin.H{"message": "pong"})  // 3. write response
    })
    r.Run(":8080")                       // 4. start listening
}`}
      />
      <p>Four lines worth understanding:</p>
      <ul>
        <li>
          <code>gin.Default()</code> gives you a router with the logger
          and panic-recovery middleware pre-attached. Grit uses this.
        </li>
        <li>
          <code>r.GET</code> registers a handler for a method+path.
          There&apos;s also <code>POST</code>, <code>PATCH</code>,{' '}
          <code>DELETE</code>, etc.
        </li>
        <li>
          <code>c *gin.Context</code> is the most important type in Gin.
          It bundles request + response + middleware data. Every handler
          gets one.
        </li>
        <li>
          <code>c.JSON(status, body)</code> serializes the body as JSON
          and sets the status code. Done.
        </li>
      </ul>

      <h2>How a Grit project wires it up</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/routes/routes.go (simplified)"
        code={`func Register(r *gin.Engine, s *Services) {
  // Health (public)
  r.GET("/healthz", s.HealthHandler.Check)

  // All API endpoints under /api
  api := r.Group("/api")

  // Public auth endpoints
  api.POST("/auth/register", s.AuthHandler.Register)
  api.POST("/auth/login",    s.AuthHandler.Login)

  // Everything below requires an auth token
  authed := api.Group("")
  authed.Use(middleware.RequireAuth(s.JWT))

  authed.GET("/users",       s.UserHandler.List)
  authed.GET("/users/:id",   s.UserHandler.Get)
  authed.PATCH("/users/:id", s.UserHandler.Update)
}`}
      />
      <p>
        Three things to notice:
      </p>
      <ul>
        <li>
          <strong>Groups.</strong> <code>r.Group(&quot;/api&quot;)</code>{' '}
          factors out a common prefix. Cleaner than typing it on every line.
        </li>
        <li>
          <strong>Middleware on groups.</strong>{' '}
          <code>authed.Use(middleware.RequireAuth(...))</code> applies
          the auth gate to ONLY routes registered after that call. Public
          ones above stay open.
        </li>
        <li>
          <strong>Path params with <code>:id</code>.</strong> Gin parses
          these for you — read with <code>c.Param(&quot;id&quot;)</code>{' '}
          inside the handler.
        </li>
      </ul>

      <h2>Reading input — the c.* family</h2>
      <CodeBlock
        language="go"
        code={`func (h *UserHandler) Get(c *gin.Context) {
  id := c.Param("id")                           // /users/:id
  q  := c.Query("expand")                       // ?expand=profile
  page := c.DefaultQuery("page", "1")           // ?page= (default "1")

  var input UpdateUserInput
  if err := c.ShouldBindJSON(&input); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})    // body parse failed
    return
  }

  authedID := c.GetUint("user_id")              // set by RequireAuth middleware
}`}
      />
      <p>
        Five lookups, five different sources. Memorise these — you&apos;ll
        type them dozens of times:
      </p>
      <ul>
        <li><code>c.Param</code> — URL placeholder (<code>:id</code>).</li>
        <li><code>c.Query</code> — querystring (<code>?key=v</code>).</li>
        <li><code>c.ShouldBindJSON</code> — JSON request body → struct.</li>
        <li><code>c.GetUint / GetString</code> — value set by middleware.</li>
        <li><code>c.GetHeader</code> — HTTP header.</li>
      </ul>

      <h2>Writing output — c.JSON, c.AbortWithStatusJSON</h2>
      <CodeBlock
        language="go"
        code={`// Success — standard Grit envelope
c.JSON(200, gin.H{"data": user, "message": "ok"})

// Error — stops middleware chain AND writes
c.AbortWithStatusJSON(404, gin.H{"error": gin.H{
  "code": "not_found",
  "message": "User not found",
}})`}
      />
      <p>
        Use <code>Abort</code> when something has gone wrong — it stops
        any downstream middleware from running. <code>c.JSON</code>{' '}
        alone keeps the chain going (rarely what you want for errors).
      </p>

      <h2>Middleware in 10 lines</h2>
      <CodeBlock
        language="go"
        code={`func RequireAuth(jwtSvc *JWTService) gin.HandlerFunc {
  return func(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
      c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
      return
    }
    claims, err := jwtSvc.Verify(strings.TrimPrefix(token, "Bearer "))
    if err != nil {
      c.AbortWithStatusJSON(401, gin.H{"error": "bad token"})
      return
    }
    c.Set("user_id", claims.UserID)  // available downstream via c.GetUint("user_id")
    c.Next()                          // continue to the handler
  }
}`}
      />
      <p>
        A middleware IS just a handler that calls <code>c.Next()</code>{' '}
        when it wants the chain to continue. Authentication, rate
        limiting, request logging, request ID — all the same shape.
      </p>

      <TipBox tone="info">
        <strong>Where this file lives in a Grit project.</strong>{' '}
        Middleware sits at <code>apps/api/internal/middleware/</code>.
        Routes at <code>apps/api/internal/routes/routes.go</code>.
        Handlers at <code>apps/api/internal/handlers/</code>. We&apos;ll
        walk all of them in the next lesson.
      </TipBox>

      <PlaygroundChallenge title="Add a hello endpoint">
        <p>
          In the playground, scaffold a minimal Gin server. Add a{' '}
          <code>GET /hello/:name</code> that returns{' '}
          <code>&#123;&quot;greeting&quot;: &quot;hi NAME&quot;&#125;</code>{' '}
          where NAME is the path param. No DB, no auth — just routing
          and a JSON response. Should take 5 minutes.
        </p>
      </PlaygroundChallenge>

      <KnowledgeCheck
        question="Inside a handler, which Gin method reads the JSON request body into a struct?"
        choices={[
          {
            label: 'c.Param',
            feedback: 'No — c.Param reads URL placeholders like :id, not the body.',
          },
          {
            label: 'c.Query',
            feedback: 'No — c.Query reads the URL querystring (?key=v), not the body.',
          },
          {
            label: 'c.ShouldBindJSON(&input)',
            correct: true,
            feedback:
              "Right — it parses the request body as JSON into the struct you pass by pointer. If parsing fails (bad JSON, missing required fields with `binding` tags), it returns an error you should respond to with 400.",
          },
          {
            label: 'c.GetBody',
            feedback: 'Not a real Gin method. ShouldBindJSON is the way.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Wire your first real endpoint by hand:</p>
            <ol>
              <li>
                In a fresh Grit project, open{' '}
                <code>apps/api/internal/routes/routes.go</code>.
              </li>
              <li>
                Add a public route:{' '}
                <code>GET /api/ping → returns &#123;&quot;data&quot;: &quot;pong&quot;&#125;</code>.
              </li>
              <li>
                Add a public route:{' '}
                <code>GET /api/echo/:msg → returns &#123;&quot;data&quot;: msg&#125;</code>{' '}
                using <code>c.Param</code>.
              </li>
              <li>
                Add an AUTHED route:{' '}
                <code>GET /api/me/ping → returns &#123;&quot;user_id&quot;: c.GetUint(&quot;user_id&quot;)&#125;</code>{' '}
                — inside the <code>authed</code> group.
              </li>
              <li>
                Test all three with curl. The authed one should return
                401 without a token and 200 with one.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            For the authed test, register a user via{' '}
            <code>POST /api/auth/register</code> first; the response
            includes a token. Send it as{' '}
            <code>Authorization: Bearer &lt;token&gt;</code>.
          </>
        }
        solution={
          <>
            <p>
              You should be comfortable opening routes.go and adding
              routes without looking anything up. Every API course
              lesson assumes you can do this.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — <strong>GORM</strong>. The ORM that turns Go
        structs into SQL. Once Gin gets the request to a handler and
        the handler calls a service, the service uses GORM to actually
        read/write the database.
      </p>
    </>
  )
}
