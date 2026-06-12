import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Grit ships JWT auth out of the box: login, register, refresh, logout,
        and a middleware that validates tokens on protected routes. This
        lesson walks how it&apos;s wired and the two extension points you&apos;ll
        actually touch.
      </p>

      <h2>The flow</h2>
      <CodeBlock
        language="text"
        code={`POST /api/auth/register     {email, password}     → 201 + {access_token, refresh_token}
POST /api/auth/login        {email, password}     → 200 + {access_token, refresh_token}
POST /api/auth/refresh      {refresh_token}       → 200 + {access_token}
POST /api/auth/logout       Authorization header   → 204
GET  /api/auth/me           Authorization header   → 200 + {user}`}
      />
      <p>
        Two tokens: a <strong>15-minute access token</strong> sent on every
        request, and a <strong>7-day refresh token</strong> used to get new
        access tokens without re-prompting login.
      </p>

      <h2>The access token shape</h2>
      <CodeBlock
        language="json"
        code={`// Decoded JWT payload
{
  "sub": "9b4d-...",      // user UUID
  "role": "user",         // role for RBAC (covered in lesson 3.4)
  "exp": 1730000000,      // expiry timestamp
  "iat": 1729000000
}`}
      />
      <p>
        Signed with <code>JWT_SECRET</code> from <code>.env</code> using
        HS256. Grit&apos;s middleware <strong>pins the algorithm</strong> —{' '}
        <code>alg:none</code> attacks are impossible.
      </p>

      <h2>The Auth middleware</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/routes/routes.go (excerpt)"
        code={`api := r.Group("/api")

// Public — no auth required
auth := api.Group("/auth")
{
    auth.POST("/login", authHandler.Login)
    auth.POST("/register", authHandler.Register)
    auth.POST("/refresh", authHandler.Refresh)
}

// Protected — Auth middleware applied
protected := api.Group("")
protected.Use(middleware.Auth(authService))
{
    protected.GET("/auth/me", authHandler.Me)
    protected.POST("/auth/logout", authHandler.Logout)
    protected.GET("/users", userHandler.List)
    // …
}`}
      />
      <p>
        The middleware reads <code>Authorization: Bearer &lt;token&gt;</code>,
        verifies signature + expiry, and puts the user ID + role on the
        request context. Handlers grab it via:
      </p>
      <CodeBlock
        language="go"
        code={`userID, _ := c.Get("user_id")
role, _ := c.Get("role")`}
      />

      <h2>What to actually extend</h2>
      <p>Two places you&apos;ll touch:</p>
      <ol>
        <li>
          <strong>What gets returned with /auth/me</strong> — add fields to
          the response in <code>authHandler.Me</code>. Avatar URL, default
          tenant, feature flags.
        </li>
        <li>
          <strong>Custom claims</strong> — add fields to the JWT payload (e.g.,
          tenant_id for multi-tenancy). Edit <code>authService.GenerateToken</code>{' '}
          + the verifier.
        </li>
      </ol>

      <TipBox tone="warning">
        <strong>JWT_SECRET</strong> must be at least 32 characters. Grit&apos;s
        config loader refuses to boot if it&apos;s shorter, and v3.25 auto-generates
        a 64-char hex secret on scaffold. Don&apos;t reuse it across environments
        — staging and prod should have different secrets.
      </TipBox>

      <h2>Refresh-token rotation</h2>
      <p>
        Each successful refresh issues a NEW refresh token and invalidates
        the old one (stored in a small DB table). This protects against
        stolen refresh tokens — a thief who uses the token rotates it, then
        the legitimate client&apos;s next refresh fails, and you know there&apos;s an
        incident.
      </p>

      <KnowledgeCheck
        question="A user signs in, gets a 15-minute access token. After 14 minutes, what's the standard flow?"
        choices={[
          {
            label: 'They sign in again at minute 15',
            feedback:
              "Wrong — that's bad UX. The whole point of the refresh token is to skip re-login.",
          },
          {
            label: 'The frontend silently calls /api/auth/refresh with the refresh token and gets a new 15-min access token',
            correct: true,
            feedback:
              "Right — the access token has a short life (limits the blast radius if stolen), the refresh token has a longer life and is the recovery vehicle. The frontend handles this transparently with lib/auth.",
          },
          {
            label: 'The API extends the access token automatically on every request',
            feedback:
              "Wrong — JWTs are stateless. The API can't extend an existing token; it can only issue new ones. That's what refresh is for.",
          },
          {
            label: 'Their session ends at minute 15 and they get a 401',
            feedback:
              "Only if the frontend doesn't auto-refresh. Grit's lib/auth handles the refresh transparently — the user never sees a 401 from this.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Try the full auth flow with curl:</p>
            <CodeBlock
              terminal
              code={`# 1. Register
curl -X POST http://localhost:8080/api/auth/register \\
  -H 'Content-Type: application/json' \\
  -d '{"email":"alex@example.com","password":"alexsecret123","name":"Alex"}'

# 2. Login (or save the access_token from register response)
curl -X POST http://localhost:8080/api/auth/login \\
  -H 'Content-Type: application/json' \\
  -d '{"email":"alex@example.com","password":"alexsecret123"}'

# 3. Hit /me with the token
curl http://localhost:8080/api/auth/me \\
  -H 'Authorization: Bearer YOUR_ACCESS_TOKEN'`}
            />
            <p>Paste the /me response in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            If /me returns 401, check that you used{' '}
            <code>Bearer </code> (with the space) before the token. Also that
            the token hasn&apos;t expired — generate a fresh one.
          </>
        }
        solution={
          <>
            <p>The /me response shape:</p>
            <CodeBlock
              language="json"
              code={`{
  "data": {
    "id": "9b4d-...",
    "email": "alex@example.com",
    "name": "Alex",
    "role": "user"
  }
}`}
            />
            <p>That&apos;s end-to-end JWT auth working.</p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Email + password is one path. Most users prefer{' '}
        <strong>OAuth2</strong> — Sign in with Google / GitHub. Next lesson
        wires both in 10 minutes.
      </p>
    </>
  )
}
