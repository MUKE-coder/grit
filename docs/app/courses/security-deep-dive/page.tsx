import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Security Deep Dive: Auth, 2FA, WAF & Rate Limiting',
  description: 'Comprehensive security course covering bcrypt passwords, JWT authentication, RBAC, two-factor authentication, backup codes, trusted devices, Sentinel WAF, rate limiting, security headers, and OAuth2.',
}

export default function SecurityDeepDiveCourse() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Security Deep Dive</span>
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
            Security Deep Dive: Auth, 2FA, WAF & Rate Limiting
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Most applications are vulnerable because security is added as an afterthought — a middleware
            bolted on after launch, a rate limiter added after an attack. Grit builds security in from
            day one. In this course, you will explore every security layer in a Grit application: password
            hashing, JWT tokens, role-based access, two-factor authentication, backup codes, trusted
            devices, a web application firewall, rate limiting, security headers, and OAuth2.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: Why Security Matters ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Why Security Matters</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Before we look at Grit{"'"}s security features, let{"'"}s understand what we{"'"}re defending against.
            Here are the most common attacks against web applications:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Brute Force:</strong> An attacker tries thousands of password combinations until one works. Without rate limiting, they can try millions of passwords per hour</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">SQL Injection:</strong> An attacker sends malicious SQL in form fields or URL parameters, tricking the database into executing unintended commands. Can dump your entire database</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">XSS (Cross-Site Scripting):</strong> An attacker injects JavaScript into your pages that runs in other users{"'"} browsers, stealing cookies, tokens, or personal data</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">CSRF (Cross-Site Request Forgery):</strong> An attacker tricks a logged-in user{"'"}s browser into making requests to your API — transferring money, changing passwords, deleting data</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Credential Stuffing:</strong> Attackers use leaked email/password combinations from other breaches to try logging into your app. Most people reuse passwords across sites</li>
          </ul>

          <Note>
            Security is not optional. A single breach can destroy user trust, violate regulations
            (GDPR, HIPAA), and cost millions in damages. The good news: Grit{"'"}s 10-layer security
            stack defends against all of these attacks out of the box.
          </Note>

          <Challenge number={1} title="Know Your Enemy">
            <p>Name 3 common web attacks and explain how each one works in your own words.
            For each attack, describe: (1) what the attacker does, (2) what they gain if successful,
            and (3) how you think it can be prevented.</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: Layer 1 — Password Security ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Layer 1: Password Security (bcrypt)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The most fundamental security rule: <strong className="text-foreground">never store passwords
            in plain text</strong>. If an attacker gains access to your database (through SQL injection,
            a backup leak, or an insider threat), plain text passwords give them immediate access to
            every account. Instead, Grit stores passwords as one-way hashes.
          </p>

          <Definition term="Hash">
            A one-way mathematical function that converts input (like a password) into a fixed-length
            string of characters. The key property: you can hash a password to get the hash, but you
            cannot reverse the hash to get the password. To verify a login, you hash the submitted
            password and compare hashes.
          </Definition>

          <Definition term="Salt">
            A random string added to the password before hashing. Without a salt, two users with the
            same password would have the same hash — making it easy to crack with precomputed tables
            (rainbow tables). With a salt, every hash is unique even for identical passwords.
          </Definition>

          <Definition term="bcrypt">
            A password hashing algorithm specifically designed to be slow. It includes a built-in salt
            and a configurable {'"'}cost factor{'"'} that controls how many iterations of hashing are
            performed. Higher cost = slower hashing = harder to brute force. The default cost factor
            is 10, meaning 2^10 (1,024) iterations.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s the process when a user registers:
          </p>

          <ol className="space-y-3 text-muted-foreground mb-6">
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">1.</span>
              <span>User submits password: <Code>myPassword123</Code></span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">2.</span>
              <span>bcrypt generates a random salt: <Code>$2a$10$N9qo8uLOickgx2ZMRZoMye</Code></span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">3.</span>
              <span>bcrypt hashes password + salt through 1,024 iterations</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">4.</span>
              <span>Result stored in database: <Code>$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy</Code></span>
            </li>
          </ol>

          <CodeBlock filename="apps/api/internal/service/auth_service.go">
{`// Hashing a password during registration
hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(input.Password),
    bcrypt.DefaultCost, // Cost factor 10 = 2^10 iterations
)

// Verifying a password during login
err := bcrypt.CompareHashAndPassword(
    []byte(user.Password), // stored hash
    []byte(input.Password), // submitted password
)`}
          </CodeBlock>

          <Challenge number={2} title="Find the Hash">
            <p>Open the auth service code in your project. Find where bcrypt hashes the password
            during registration and where it compares during login. What cost factor is used? What
            would happen if you changed it to 14? (Hint: each increment doubles the time.)</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Layer 2 — JWT Authentication ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Layer 2: JWT Authentication</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            After verifying a password, Grit issues two tokens: an <strong className="text-foreground">access
            token</strong> (short-lived, 15 minutes) and a <strong className="text-foreground">refresh
            token</strong> (long-lived, 7 days). This dual-token approach limits the damage if a token
            is stolen.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Why two tokens? If an access token is intercepted by an attacker, it only works for 15
            minutes. The refresh token is stored in an HTTP-only cookie (inaccessible to JavaScript),
            so XSS attacks cannot steal it. When the access token expires, the frontend silently
            uses the refresh token to get a new one.
          </p>

          <CodeBlock filename="Token Lifecycle">
{`Login successful
  ├── Access Token (JWT, 15 min expiry)
  │   └── Sent in Authorization header: Bearer <token>
  │   └── Contains: user ID, email, role
  │   └── Stateless — no database lookup needed
  │
  └── Refresh Token (JWT, 7 day expiry)
      └── Stored in HTTP-only cookie
      └── Used only to get new access tokens
      └── Cannot be read by JavaScript (XSS protection)`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The auth middleware runs on every protected route. It extracts the JWT from the
            Authorization header, verifies the signature, checks the expiration, and attaches the
            user information to the request context.
          </p>

          <Challenge number={3} title="Decode a JWT">
            <p>Login to your Grit API and copy the access token. Go to <a href="https://jwt.io" target="_blank" rel="noreferrer" className="text-primary hover:underline">jwt.io</a> and
            paste it. What information is in the payload? What algorithm is used for signing?
            What is the expiration time (exp claim)?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Layer 3 — RBAC ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Layer 3: Role-Based Access Control</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Authentication answers {'"'}who are you?{'"'} Authorization answers {'"'}what are you allowed
            to do?{'"'} Grit uses role-based access control (RBAC) with three built-in roles:
          </p>

          <div className="overflow-x-auto mb-6">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Role</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Permissions</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">ADMIN</td>
                  <td className="px-4 py-3">Full access — manage users, view all data, access admin panel, system settings</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">EDITOR</td>
                  <td className="px-4 py-3">Create and edit content — manage resources, upload files, but cannot manage users</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">USER</td>
                  <td className="px-4 py-3">Basic access — view content, manage own profile, limited write access</td>
                </tr>
              </tbody>
            </table>
          </div>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>RequireRole</Code> middleware enforces access control at the route level:
          </p>

          <CodeBlock filename="apps/api/internal/routes/routes.go">
{`// Public routes — no auth required
auth := api.Group("/auth")
auth.POST("/login", h.Login)
auth.POST("/register", h.Register)

// Protected routes — any authenticated user
users := api.Group("/users")
users.Use(middleware.AuthMiddleware(cfg))
users.GET("/me", h.GetMe)

// Admin-only routes
admin := api.Group("/admin")
admin.Use(middleware.AuthMiddleware(cfg))
admin.Use(middleware.RequireRole("ADMIN"))
admin.GET("/users", h.ListUsers)
admin.DELETE("/users/:id", h.DeleteUser)`}
          </CodeBlock>

          <Challenge number={4} title="Test Role Enforcement">
            <p>Log in as a user with the USER role. Try to access an ADMIN-only endpoint like
            <Code> GET /api/admin/users</Code>. What HTTP status code do you get? What does the
            error message say? Now log in as an ADMIN and try the same endpoint.</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Layer 4 — 2FA ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Layer 4: Two-Factor Authentication (TOTP)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Passwords can be stolen through phishing, keyloggers, or data breaches. Two-factor
            authentication (2FA) adds a second layer: something you <strong className="text-foreground">know</strong> (password)
            plus something you <strong className="text-foreground">have</strong> (your phone). Even if
            an attacker has your password, they cannot log in without your phone.
          </p>

          <Definition term="TOTP (Time-based One-Time Password)">
            An algorithm defined in RFC 6238 that generates a 6-digit code that changes every 30
            seconds. Both your authenticator app and the server share a secret key. They independently
            calculate the same code based on the current time. If the codes match, authentication
            succeeds.
          </Definition>

          <Definition term="Authenticator App">
            A mobile application (Google Authenticator, Authy, 1Password, Bitwarden) that stores
            TOTP secrets and displays the current 6-digit code. You scan a QR code to add an account,
            and the app generates codes offline — no internet connection needed.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The 2FA setup flow in Grit:
          </p>

          <ol className="space-y-3 text-muted-foreground mb-6">
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">1.</span>
              <span><strong className="text-foreground">Generate secret.</strong> The API creates a random TOTP secret and returns it as a QR code URL.</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">2.</span>
              <span><strong className="text-foreground">Scan QR code.</strong> The user scans the QR code with their authenticator app, which stores the secret.</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">3.</span>
              <span><strong className="text-foreground">Verify code.</strong> The user enters the current 6-digit code from their app to prove it{"'"}s working.</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">4.</span>
              <span><strong className="text-foreground">2FA enabled.</strong> The secret is saved to the database. Future logins require both password and TOTP code.</span>
            </li>
          </ol>

          <CodeBlock filename="2FA Login Flow">
{`Step 1: POST /api/auth/login
  Body: { "email": "user@example.com", "password": "..." }
  Response: { "requires_2fa": true, "temp_token": "..." }

Step 2: POST /api/auth/verify-2fa
  Body: { "temp_token": "...", "code": "482917" }
  Response: { "access_token": "...", "refresh_token": "..." }`}
          </CodeBlock>

          <Tip>
            The 30-second time window actually has a small grace period. Most TOTP implementations
            accept the previous code and the next code in addition to the current one. This accounts
            for slight clock drift between the server and the user{"'"}s phone.
          </Tip>

          <Challenge number={5} title="Enable 2FA">
            <p>Enable two-factor authentication on your account using Google Authenticator, Authy,
            or any TOTP-compatible app. Follow the setup flow: generate the secret, scan the QR code,
            enter the verification code. Then log out and log back in — you should be prompted for
            a 6-digit code.</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Layer 5 — Backup Codes ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Layer 5: Backup Codes</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            What happens if you lose your phone? Without backup codes, you{"'"}re locked out of your
            account forever. Grit generates 10 one-time backup codes when 2FA is enabled. Each code
            can be used exactly once as a substitute for the TOTP code.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Backup codes are:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">10 codes</strong> generated when 2FA is enabled</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">One-time use</strong> — each code is invalidated after use</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">bcrypt-hashed</strong> in the database — if the database is compromised, the raw codes are not exposed</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Shown once</strong> — the user sees the plain text codes only at generation time and must save them</li>
          </ul>

          <CodeBlock filename="Backup Code Format">
{`Your backup codes (save these somewhere safe):

1.  a8f2-k9m3-p4x7
2.  b3n6-j7w2-q8r1
3.  c5t9-h1y4-s6v3
4.  d2p8-g3z7-u9e5
5.  e7k1-f6x4-w2m8
6.  f4j3-e9n6-y7t2
7.  g1w7-d2p5-z3k9
8.  h8y4-c7m1-a6j3
9.  i6z2-b4t8-e1w5
10. j3x9-a1s7-f4n6

Each code can only be used once.
Store them in a password manager or print them.`}
          </CodeBlock>

          <Challenge number={6} title="Test Backup Codes">
            <p>With 2FA enabled, use a backup code to log in instead of your authenticator app.
            Then try the same backup code again. What happens? How many backup codes do you have
            remaining?</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Layer 6 — Trusted Devices ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Layer 6: Trusted Devices</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Entering a 2FA code every time you log in on the same computer gets annoying fast.
            Trusted devices solve this: after you verify 2FA once, the server sets a 30-day cookie
            on your browser. On subsequent logins from that device, the 2FA step is skipped.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The trusted device cookie is:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Secure</strong> — only sent over HTTPS (in production)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">HTTP-only</strong> — JavaScript cannot read it (XSS protection)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">30-day expiry</strong> — after 30 days, you{"'"}ll need to enter a 2FA code again</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Per-device</strong> — trusting your laptop does not trust your phone{"'"}s browser</li>
          </ul>

          <Challenge number={7} title="Inspect the Cookie">
            <p>After completing a 2FA verification, open your browser{"'"}s DevTools (F12) and go to
            the Application tab → Cookies. Find the trusted device cookie. What is its name? What is
            its expiry date? Is it marked as HTTP-only and Secure?</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Layer 7 — Sentinel WAF ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Layer 7: Sentinel WAF</h2>

          <Definition term="WAF (Web Application Firewall)">
            A security layer that sits between the internet and your application, inspecting every
            incoming request. It blocks malicious requests — SQL injection attempts, XSS payloads,
            suspicious user agents, and known attack patterns — before they reach your handlers.
            Think of it as a bouncer at the door of your API.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit includes <strong className="text-foreground">Sentinel</strong>, a built-in WAF that
            runs as middleware. It inspects every request and blocks:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">SQL injection patterns</strong> — detects <Code>{"'"} OR 1=1 --</Code>, <Code>UNION SELECT</Code>, and other injection attempts in query parameters and request bodies</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">XSS attempts</strong> — blocks <Code>{'<script>'}</Code> tags, event handlers like <Code>onerror=</Code>, and other JavaScript injection vectors</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Suspicious user agents</strong> — blocks known attack tools like sqlmap, nikto, and dirbuster</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Path traversal</strong> — blocks attempts to access files outside the web root using <Code>../../</Code></li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Sentinel also includes a dashboard at <Code>/sentinel/ui</Code> where you can see:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> Total requests processed</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Blocked requests and the reasons they were blocked</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Top attacking IPs</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Attack type distribution</li>
          </ul>

          <Challenge number={8} title="Explore the WAF Dashboard">
            <p>Visit <Code>/sentinel/ui</Code> in your browser. What metrics does the dashboard show?
            How many requests have been processed? Have any been blocked? Try sending a request
            with a SQL injection pattern in a query parameter and check if Sentinel catches it.</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Layer 8 — Rate Limiting ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Layer 8: Rate Limiting</h2>

          <Definition term="Rate Limiting">
            Restricting how many requests a client can make within a time window. For example,
            100 requests per minute per IP address. This prevents brute force attacks (trying
            thousands of passwords), DDoS attacks (overwhelming your server with traffic), and
            API abuse (scraping all your data).
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit{"'"}s Sentinel middleware includes configurable rate limiting. Limits are applied
            per IP address and can be configured per route or globally:
          </p>

          <CodeBlock filename="Rate Limiting Configuration">
{`// Default rate limits in Grit
Global:     100 requests per minute per IP
Auth:       10 requests per minute per IP (login, register)
Upload:     20 requests per minute per IP
API:        60 requests per minute per IP

// When a client exceeds the limit:
// HTTP 429 Too Many Requests
{
  "error": {
    "code": "RATE_LIMITED",
    "message": "Too many requests. Please try again later.",
    "details": {
      "retry_after": 45
    }
  }
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The auth endpoints have a stricter limit (10/minute) because brute force attacks
            target login endpoints. An attacker trying 10,000 passwords per minute would be
            blocked after the first 10 attempts.
          </p>

          <Tip>
            Rate limiting is enforced using Redis for distributed counting. If your API runs on
            multiple servers behind a load balancer, the rate limit is shared across all instances.
            An attacker cannot bypass the limit by hitting different servers.
          </Tip>

          <Challenge number={9} title="Find the Rate Limits">
            <p>What rate limits are configured in your Grit project by default? Look at the
            Sentinel configuration in your code. What happens when you exceed the limit? What
            HTTP status code is returned? What does the <Code>retry_after</Code> field tell the client?</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Layer 9 — Security Headers ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Layer 9: Security Headers</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            HTTP response headers can instruct browsers to enable security features. Grit sets
            these headers on every response:
          </p>

          <div className="overflow-x-auto mb-6">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Header</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Value</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Prevents</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Strict-Transport-Security</td>
                  <td className="px-4 py-3"><Code>max-age=31536000</Code></td>
                  <td className="px-4 py-3">Forces HTTPS for 1 year, prevents downgrade attacks</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">X-Frame-Options</td>
                  <td className="px-4 py-3"><Code>DENY</Code></td>
                  <td className="px-4 py-3">Prevents your site from being embedded in iframes (clickjacking)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">X-Content-Type-Options</td>
                  <td className="px-4 py-3"><Code>nosniff</Code></td>
                  <td className="px-4 py-3">Prevents browsers from guessing content types (MIME sniffing)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Content-Security-Policy</td>
                  <td className="px-4 py-3"><Code>default-src {"'"}self{"'"}</Code></td>
                  <td className="px-4 py-3">Controls which resources can load, blocks inline scripts (XSS)</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Referrer-Policy</td>
                  <td className="px-4 py-3"><Code>strict-origin-when-cross-origin</Code></td>
                  <td className="px-4 py-3">Limits referrer information sent to external sites (privacy)</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Challenge number={10} title="Check Your Headers">
            <p>Open DevTools (F12) in your browser. Go to the Network tab. Make any request to
            your API. Click on the request and look at the Response Headers section. List all
            the security headers you see. Are all five headers from the table present?</p>
          </Challenge>
        </section>

        {/* ═══ Section 11: Layer 10 — JWT Secret Management ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Layer 10: JWT Secret Management</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The JWT secret is the key used to sign and verify tokens. If an attacker knows your
            secret, they can forge tokens for any user — including admin accounts. The secret must be:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Long</strong> — at least 32 characters (256 bits)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Random</strong> — generated by a cryptographic random number generator, not a human-chosen phrase</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Secret</strong> — stored in environment variables, never committed to Git</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Unique per environment</strong> — development, staging, and production should each have different secrets</li>
          </ul>

          <CodeBlock filename="Generate a Strong JWT Secret">
{`# Generate a 64-character hex string (256 bits of entropy)
openssl rand -hex 32

# Example output:
# a3f8b2c1e9d7f4a6b8c2e1d3f5a7b9c1e3d5f7a9b1c3e5d7f9a1b3c5e7d9f1a3

# Add to your .env file (NEVER commit this file)
JWT_SECRET=a3f8b2c1e9d7f4a6b8c2e1d3f5a7b9c1e3d5f7a9b1c3e5d7f9a1b3c5e7d9f1a3`}
          </CodeBlock>

          <Note>
            Changing the JWT secret invalidates ALL existing tokens instantly. Every logged-in user
            will be forced to log in again. This is actually useful — if you suspect a token has
            been compromised, rotating the secret is the nuclear option that immediately revokes
            all sessions.
          </Note>

          <Challenge number={11} title="Secure Your Secret">
            <p>Generate a strong JWT secret using <Code>openssl rand -hex 32</Code>. Replace the
            default JWT_SECRET in your <Code>.env</Code> file. Restart the API. Are you still logged
            in? (You should not be — the old tokens are now invalid.)</p>
          </Challenge>
        </section>

        {/* ═══ Section 12: OAuth2 Security ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">OAuth2 Security</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Social login (Google, GitHub, etc.) is not just convenient — it{"'"}s a security feature.
            When users log in with Google, your application never sees their password. This prevents:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Password reuse attacks</strong> — if a user{"'"}s password is leaked from another site, your app is unaffected because you never stored a password</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Phishing</strong> — users authenticate on Google{"'"}s domain, not on a page you control, so fake login pages are harder to create</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Weak passwords</strong> — Google enforces its own password policies and 2FA, which are typically stronger than what most apps enforce</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit links OAuth accounts by email. If a user registers with email/password and later
            logs in with Google using the same email, the accounts are automatically linked. One
            user, multiple login methods.
          </p>

          <Challenge number={12} title="Explore OAuth Config">
            <p>Look at the OAuth configuration in your project. What scopes does Google OAuth
            request? (Scopes define what information your app can access from Google.) What
            happens if a user logs in with Google using an email that already has a password account?</p>
          </Challenge>
        </section>

        {/* ═══ Section 13: Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary: The 10-Layer Security Stack</h2>

          <div className="overflow-x-auto mb-6">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Layer</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Feature</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Protects Against</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-primary">1</td>
                  <td className="px-4 py-3 font-medium text-foreground">bcrypt Passwords</td>
                  <td className="px-4 py-3">Database leaks, rainbow tables</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-primary">2</td>
                  <td className="px-4 py-3 font-medium text-foreground">JWT Tokens</td>
                  <td className="px-4 py-3">Session hijacking, token theft (short expiry)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-primary">3</td>
                  <td className="px-4 py-3 font-medium text-foreground">RBAC</td>
                  <td className="px-4 py-3">Privilege escalation, unauthorized access</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-primary">4</td>
                  <td className="px-4 py-3 font-medium text-foreground">TOTP 2FA</td>
                  <td className="px-4 py-3">Stolen passwords, phishing</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-primary">5</td>
                  <td className="px-4 py-3 font-medium text-foreground">Backup Codes</td>
                  <td className="px-4 py-3">Lost device lockout</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-primary">6</td>
                  <td className="px-4 py-3 font-medium text-foreground">Trusted Devices</td>
                  <td className="px-4 py-3">2FA fatigue (usability without sacrificing security)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-primary">7</td>
                  <td className="px-4 py-3 font-medium text-foreground">Sentinel WAF</td>
                  <td className="px-4 py-3">SQL injection, XSS, path traversal</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-primary">8</td>
                  <td className="px-4 py-3 font-medium text-foreground">Rate Limiting</td>
                  <td className="px-4 py-3">Brute force, DDoS, API abuse</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-primary">9</td>
                  <td className="px-4 py-3 font-medium text-foreground">Security Headers</td>
                  <td className="px-4 py-3">Clickjacking, MIME sniffing, protocol downgrade</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-mono text-primary">10</td>
                  <td className="px-4 py-3 font-medium text-foreground">JWT Secret Management</td>
                  <td className="px-4 py-3">Token forgery, session compromise</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Challenge number={13} title="Security Audit — Part 1">
            <p>Perform the first half of a security audit on your Grit application:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Check your <Code>JWT_SECRET</Code> — is it at least 32 characters? Is it random (not a dictionary word)?</li>
              <li>Verify password hashing — register a user, then look at the users table in GORM Studio. Is the password a bcrypt hash (starts with <Code>$2a$</Code>)?</li>
              <li>Test 2FA — enable it, log out, log back in. Does it require the 6-digit code?</li>
              <li>Test a backup code — use one, then try it again. Is it rejected the second time?</li>
            </ol>
          </Challenge>

          <Challenge number={14} title="Security Audit — Part 2">
            <p>Complete the second half of the security audit:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Test rate limiting — send 15 rapid login requests. Are you rate limited?</li>
              <li>Check security headers — inspect response headers in DevTools. Are all 5 headers present?</li>
              <li>Test RBAC — log in as USER, attempt an ADMIN endpoint. Is it blocked?</li>
              <li>Check the Sentinel dashboard — are any attacks logged?</li>
            </ol>
          </Challenge>

          <Challenge number={15} title="Final Challenge: Write a Security Report">
            <p>Write a one-page security report for your Grit application. For each of the 10
            security layers, document: (1) whether it{"'"}s enabled, (2) how it{"'"}s configured,
            and (3) any improvements you would make. This is the kind of document you would give
            to a security auditor or include in your project{"'"}s documentation.</p>
          </Challenge>
        </section>

        {/* ═══ Footer ═══ */}
        <CourseFooter />

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
