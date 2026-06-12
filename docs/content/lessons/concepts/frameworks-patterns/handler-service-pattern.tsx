import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        Every Grit endpoint follows the same two-layer pattern:{' '}
        <strong>handler</strong> talks to HTTP, <strong>service</strong>{' '}
        talks to the DB. Get this right and every future feature you
        build feels obvious. Get it wrong and your codebase rots.
      </p>

      <h2>The split — at a glance</h2>
      <Diagram label="Handler / Service responsibilities" caption="The handler is the airlock between HTTP and your business logic. The service is the business logic.">
{`   ┌─────────────────────────────────────────────────────────────┐
   │   HANDLER  (apps/api/internal/handlers/*.go)                │
   ├─────────────────────────────────────────────────────────────┤
   │   - parse path / query / body  → typed input                │
   │   - call ONE service method                                 │
   │   - translate result + errors  → JSON + HTTP status         │
   │                                                             │
   │   No DB calls.  No business rules.  Thin.                   │
   └─────────────────────────────────────────────────────────────┘
                              │
                              ▼  (typed input)
   ┌─────────────────────────────────────────────────────────────┐
   │   SERVICE  (apps/api/internal/services/*.go)                │
   ├─────────────────────────────────────────────────────────────┤
   │   - validate business rules                                 │
   │   - GORM calls (read/write)                                 │
   │   - orchestrate other services (email, cache, jobs)         │
   │   - return DOMAIN value or DOMAIN error                     │
   │                                                             │
   │   No HTTP types.  No JSON.  Just Go.                        │
   └─────────────────────────────────────────────────────────────┘`}
      </Diagram>

      <h2>Why split? Three concrete reasons</h2>

      <h3>1. Testability</h3>
      <p>
        Services are plain Go functions with plain Go inputs. They&apos;re
        trivial to unit-test:
      </p>
      <CodeBlock
        language="go"
        code={`func TestRegisterRejectsDuplicate(t *testing.T) {
  s := NewAuthService(testDB(t))
  _, err := s.Register(ctx, RegisterInput{Email: "a@b.com", Password: "secret"})
  require.NoError(t, err)
  _, err = s.Register(ctx, RegisterInput{Email: "a@b.com", Password: "other"})
  require.ErrorIs(t, err, ErrEmailTaken)  // domain error, not HTTP status
}`}
      />
      <p>
        If the handler held the DB calls, you&apos;d have to spin up an
        HTTP server, send a request, parse a response — three layers
        of ceremony for one assertion. Service tests are direct.
      </p>

      <h3>2. Reuse</h3>
      <p>
        Logic in a service can be called from:
      </p>
      <ul>
        <li>An HTTP handler (the normal case).</li>
        <li>A CLI command (<code>grit seed</code>, <code>grit admin make-user</code>).</li>
        <li>A background job (the welcome-email worker).</li>
        <li>Another service (Order service calls UserService).</li>
      </ul>
      <p>
        Logic in a handler can be called from… an HTTP request. Once.
        That&apos;s the whole point.
      </p>

      <h3>3. Single responsibility</h3>
      <p>
        When you read a handler, you should see the SHAPE of the
        endpoint: what it accepts, what HTTP status it returns. When
        you read a service, you should see the RULES of the business.
        Splitting them gives each file a clear job and a clear test.
      </p>

      <h2>Concrete example — register a user</h2>

      <h3>The handler</h3>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/auth_handler.go"
        code={`func (h *AuthHandler) Register(c *gin.Context) {
  // 1. parse input
  var in services.RegisterInput
  if err := c.ShouldBindJSON(&in); err != nil {
    c.JSON(400, gin.H{"error": gin.H{"code": "invalid_body", "message": err.Error()}})
    return
  }

  // 2. call the service — that's the only line of business logic in this handler
  user, token, err := h.auth.Register(c.Request.Context(), in)

  // 3. translate result + errors to HTTP
  switch {
  case errors.Is(err, services.ErrEmailTaken):
    c.JSON(409, gin.H{"error": gin.H{"code": "email_taken", "message": "Email already in use"}})
  case err != nil:
    c.JSON(500, gin.H{"error": gin.H{"code": "internal", "message": "Something went wrong"}})
  default:
    c.JSON(201, gin.H{"data": gin.H{"user": user, "token": token}, "message": "Account created"})
  }
}`}
      />
      <p>
        Three sections, in order: parse → call → translate. Notice
        what&apos;s NOT here: no GORM, no password hashing, no email
        check. The handler doesn&apos;t know HOW; it knows WHAT to
        return.
      </p>

      <h3>The service</h3>
      <CodeBlock
        language="go"
        filename="apps/api/internal/services/auth_service.go"
        code={`var ErrEmailTaken = errors.New("email already in use")

type RegisterInput struct {
  Email    string \`json:"email"    binding:"required,email"\`
  Password string \`json:"password" binding:"required,min=8"\`
  Name     string \`json:"name"     binding:"required"\`
}

func (s *AuthService) Register(ctx context.Context, in RegisterInput) (models.User, string, error) {
  // 1. business rule: email must be unique
  var existing models.User
  if err := s.db.WithContext(ctx).Where("email = ?", in.Email).First(&existing).Error; err == nil {
    return models.User{}, "", ErrEmailTaken
  }

  // 2. hash password
  hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
  if err != nil { return models.User{}, "", err }

  // 3. create
  user := models.User{Email: in.Email, PasswordHash: string(hash), Name: in.Name}
  if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
    return models.User{}, "", err
  }

  // 4. side effect: enqueue welcome email
  _ = s.jobs.Enqueue("send_welcome", map[string]any{"user_id": user.ID})

  // 5. issue token
  token, err := s.jwt.Issue(user.ID)
  if err != nil { return models.User{}, "", err }

  return user, token, nil
}`}
      />
      <p>
        Notice the difference in vocabulary:
      </p>
      <ul>
        <li>Handler talks about <em>400, 409, 201, JSON</em>.</li>
        <li>Service talks about <em>email taken, hash password, enqueue welcome</em>.</li>
      </ul>
      <p>
        Different layers, different concerns. That separation is the
        whole pattern.
      </p>

      <h2>Domain errors instead of HTTP codes</h2>
      <p>
        The service returns <code>ErrEmailTaken</code> — a sentinel
        error. The handler maps it to HTTP 409. This indirection costs
        almost nothing and buys you:
      </p>
      <ul>
        <li>
          Services that work in non-HTTP contexts (CLI, jobs) — they
          don&apos;t care about status codes.
        </li>
        <li>
          A single place to change the mapping. If you decide email
          collision should be 422 instead of 409, you edit one line in
          the handler.
        </li>
        <li>
          Tests that assert business intent (<code>ErrEmailTaken</code>),
          not transport details (<code>409</code>).
        </li>
      </ul>

      <h2>Where each piece lives</h2>
      <CodeBlock
        language="text"
        code={`apps/api/internal/
├── handlers/          ← thin, HTTP-aware
│   ├── auth_handler.go
│   ├── user_handler.go
│   └── product_handler.go
├── services/          ← business logic
│   ├── auth_service.go
│   ├── user_service.go
│   └── product_service.go
├── models/            ← GORM structs (the data shape)
│   ├── user.go
│   └── product.go
├── routes/            ← routes.go wires handlers to URLs
└── middleware/        ← auth, CORS, logger, rate limit`}
      />
      <p>
        Same triple — <code>model + service + handler</code> — for
        every resource. <code>grit generate resource</code> generates
        all three in lockstep so you stop typing the boilerplate.
      </p>

      <TipBox tone="warning">
        <strong>The temptation: shortcut from handler to DB.</strong>{' '}
        &quot;It&apos;s just one query, why bother with a service?&quot; —
        because the day after, you need to fire an email too. And cache
        the result. And add a unique check. Suddenly the handler is
        100 lines and you have to test it through HTTP. Save the
        future-you trip; service-out from the start.
      </TipBox>

      <h2>When the rule bends</h2>
      <p>
        Two rare exceptions:
      </p>
      <ul>
        <li>
          <strong>Pure echo / health endpoints.</strong>{' '}
          <code>GET /healthz</code> returning <code>200 ok</code>{' '}
          doesn&apos;t need a service. Skip the layer; it&apos;s
          ceremony.
        </li>
        <li>
          <strong>Static file serving.</strong> If Gin is just shipping
          bytes from disk, no business logic is involved. No service.
        </li>
      </ul>
      <p>
        Anything that touches the DB, calls another service, or
        applies a business rule — service. Always.
      </p>

      <h2>What about a Repository layer?</h2>
      <p>
        Some teams add a <code>repository</code> between service and
        GORM, so the service doesn&apos;t know about GORM. Pros:
        swappable DB layer. Cons: another file per resource, another
        interface, more clicks.
      </p>
      <p>
        Grit&apos;s default is service → GORM, no repository. Add a
        repository ONLY if you have a real reason (e.g., you&apos;re
        going to swap to a different ORM, or you need both Postgres
        and a NoSQL store for the same data). For 95% of projects, the
        extra layer is overhead.
      </p>

      <KnowledgeCheck
        question="A new feature: when a Product is created, send an email to all admins. Where do you put this code?"
        choices={[
          {
            label: 'In the Product handler — right after c.JSON(201, ...)',
            feedback:
              "Wrong layer. Now the email send is HTTP-coupled. CLI imports or background jobs that create products won't send the email. Also harder to test.",
          },
          {
            label: 'In the Product SERVICE — inside Create, after the DB insert succeeds',
            correct: true,
            feedback:
              "Right — that's the &quot;business rule&quot; layer. Any caller that creates a product (HTTP, CLI, job) will get the email side-effect for free. Test the service, you've tested the rule.",
          },
          {
            label: 'In a middleware so it runs for any POST',
            feedback:
              "Middleware is for cross-cutting concerns (auth, logging). Business rules belong with the entity they're about.",
          },
          {
            label: 'In a global goroutine that watches the DB',
            feedback:
              "Wildly overcomplicated. Service-level enqueue + a job worker handles this with normal code paths.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Refactor a handler to enforce the split:</p>
            <ol>
              <li>
                In a fresh project, find any handler that calls{' '}
                <code>h.db</code> directly (some old code, your earlier
                CRUD lesson, or hand-written).
              </li>
              <li>
                Move the DB calls into a new service method. The
                handler should call that method and only that method.
              </li>
              <li>
                If the handler had business validation (&quot;email
                must be unique&quot;), move that into the service too.
              </li>
              <li>
                If the handler had error mapping ifs, keep those in
                the handler but have the service return sentinel
                errors instead of HTTP statuses.
              </li>
              <li>
                Write ONE service-level unit test (no HTTP needed) for
                a business rule.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            For the unit test, you only need an in-memory SQLite DB —
            see the previous lesson for the snippet. No fixtures, no
            Postgres, no Gin.
          </>
        }
        solution={
          <>
            <p>
              You should have: a thin handler (parse, call, translate),
              a thick service (rules, DB, side-effects), and a unit
              test that exercises a rule without HTTP. That&apos;s the
              pattern in production form.
            </p>
            <p>
              Every other course in the kit assumes you can do this.
              The pattern repeats across thousands of lines without
              variation — that&apos;s the point.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — <strong>CRUD walkthrough</strong>. You&apos;ve
        seen the pieces; now we apply them to all four CRUD operations
        on a real resource, end to end.
      </p>
    </>
  )
}
