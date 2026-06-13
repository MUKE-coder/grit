import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/authentication')

export default function AuthenticationPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend (Go API)</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Authentication
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit ships with a complete JWT-based authentication system. It includes register,
                login, token refresh, logout, password reset, role-based access control,
                and two-factor authentication (TOTP) with backup codes and trusted devices --
                all pre-configured and ready to use.
              </p>
            </div>

            <div className="prose-grit">
              {/* ── Auth Flow ─────────────────────────────── */}
              <h2 id="auth-flow">Authentication Flow</h2>
              <p>
                Grit uses a dual-token JWT strategy: a short-lived <strong>access token</strong> for
                API requests and a long-lived <strong>refresh token</strong> for obtaining new access
                tokens without re-authenticating.
              </p>

              <CodeBlock filename="authentication flow" code={`
  Client                           Grit API
    |                                  |
    |  POST /api/auth/register         |
    |  { name, email, password }       |
    | -------------------------------->|
    |                                  |  Hash password (bcrypt)
    |                                  |  Create user in DB
    |                                  |  Generate access + refresh tokens
    |  { user, tokens }                |
    | <--------------------------------|
    |                                  |
    |  GET /api/posts                  |
    |  Authorization: Bearer <access>  |
    | -------------------------------->|
    |                                  |  Validate JWT
    |                                  |  Load user from DB
    |  { data: [...] }                 |
    | <--------------------------------|
    |                                  |
    |  --- access token expires ---    |
    |                                  |
    |  POST /api/auth/refresh          |
    |  { refresh_token }               |
    | -------------------------------->|
    |                                  |  Validate refresh token
    |                                  |  Generate new token pair
    |  { tokens }                      |
    | <--------------------------------|
    |                                  |
`} />

              {/* ── JWT Tokens ─────────────────────────────── */}
              <h2 id="jwt-tokens">JWT Tokens</h2>
              <p>
                Grit generates two JWT tokens on login/register. Both are signed with
                HMAC-SHA256 using the <code>JWT_SECRET</code> environment variable.
              </p>

              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Token</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Default Expiry</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Purpose</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">access_token</td>
                      <td className="px-4 py-2.5 font-mono text-xs">15 minutes</td>
                      <td className="px-4 py-2.5">Sent with every API request in the Authorization header</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">refresh_token</td>
                      <td className="px-4 py-2.5 font-mono text-xs">7 days (168h)</td>
                      <td className="px-4 py-2.5">Used to get a new access token when it expires</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <p>Configure token expiry via environment variables:</p>
              <CodeBlock language="bash" filename=".env" code={`JWT_SECRET=your-super-secret-key-at-least-32-chars
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h`} />

              <h3 id="token-claims">Token Claims (JWT Payload)</h3>
              <p>Each token contains these claims:</p>
              <CodeBlock filename="services/auth.go" code={`type Claims struct {
    UserID uint   \`json:"user_id"\`
    Email  string \`json:"email"\`
    Role   string \`json:"role"\`
    jwt.RegisteredClaims  // exp, iat, etc.
}`} />

              {/* ── Auth Endpoints ─────────────────────────────── */}
              <h2 id="auth-endpoints">Auth Endpoints</h2>
              <p>
                All auth endpoints are mounted at <code>/api/auth</code>. Register, login,
                refresh, and forgot/reset-password are public. Me and logout require authentication.
              </p>

              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Method</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Endpoint</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Auth</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Description</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">POST</td>
                      <td className="px-4 py-2.5 font-mono text-xs">/api/auth/register</td>
                      <td className="px-4 py-2.5">No</td>
                      <td className="px-4 py-2.5">Create a new user account</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">POST</td>
                      <td className="px-4 py-2.5 font-mono text-xs">/api/auth/login</td>
                      <td className="px-4 py-2.5">No</td>
                      <td className="px-4 py-2.5">Authenticate and get tokens</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">POST</td>
                      <td className="px-4 py-2.5 font-mono text-xs">/api/auth/refresh</td>
                      <td className="px-4 py-2.5">No</td>
                      <td className="px-4 py-2.5">Get new tokens with refresh token</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">GET</td>
                      <td className="px-4 py-2.5 font-mono text-xs">/api/auth/me</td>
                      <td className="px-4 py-2.5">Yes</td>
                      <td className="px-4 py-2.5">Get current authenticated user</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">POST</td>
                      <td className="px-4 py-2.5 font-mono text-xs">/api/auth/logout</td>
                      <td className="px-4 py-2.5">Yes</td>
                      <td className="px-4 py-2.5">Invalidate user session</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">POST</td>
                      <td className="px-4 py-2.5 font-mono text-xs">/api/auth/forgot-password</td>
                      <td className="px-4 py-2.5">No</td>
                      <td className="px-4 py-2.5">Request a password reset link</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">POST</td>
                      <td className="px-4 py-2.5 font-mono text-xs">/api/auth/reset-password</td>
                      <td className="px-4 py-2.5">No</td>
                      <td className="px-4 py-2.5">Reset password with token</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              {/* ── Register ─────────────────────────────── */}
              <h3 id="register">Register</h3>
              <CodeBlock filename="POST /api/auth/register" code={`// Request
{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securepassword123"
}

// Response (201 Created)
{
    "data": {
        "user": {
            "id": 1,
            "name": "John Doe",
            "email": "john@example.com",
            "role": "user",
            "avatar": "",
            "active": true,
            "email_verified_at": null,
            "created_at": "2026-02-11T10:00:00Z",
            "updated_at": "2026-02-11T10:00:00Z"
        },
        "tokens": {
            "access_token": "eyJhbGciOiJIUzI1NiIs...",
            "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
            "expires_at": 1707649200
        }
    },
    "message": "User registered successfully"
}`} />

              <h3 id="login">Login</h3>
              <CodeBlock filename="POST /api/auth/login" code={`// Request
{
    "email": "john@example.com",
    "password": "securepassword123"
}

// Response (200 OK)
{
    "data": {
        "user": { ... },
        "tokens": {
            "access_token": "eyJhbGciOiJIUzI1NiIs...",
            "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
            "expires_at": 1707649200
        }
    },
    "message": "Logged in successfully"
}

// Error (401 Unauthorized)
{
    "error": {
        "code": "INVALID_CREDENTIALS",
        "message": "Invalid email or password"
    }
}`} />

              <h3 id="refresh">Refresh Token</h3>
              <CodeBlock filename="POST /api/auth/refresh" code={`// Request
{
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}

// Response (200 OK)
{
    "data": {
        "tokens": {
            "access_token": "eyJhbGciOiJIUzI1NiIs...",
            "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
            "expires_at": 1707650100
        }
    },
    "message": "Token refreshed successfully"
}`} />

              <h3 id="forgot-password">Forgot Password</h3>
              <CodeBlock filename="POST /api/auth/forgot-password" code={`// Request
{
    "email": "john@example.com"
}

// Response (200 OK) -- always returns success for security
{
    "message": "If an account with that email exists, a password reset link has been sent"
}`} />
              <p>
                The forgot-password endpoint always returns a success message regardless of whether
                the email exists. This prevents email enumeration attacks.
              </p>

              <h3 id="reset-password">Reset Password</h3>
              <CodeBlock filename="POST /api/auth/reset-password" code={`// Request
{
    "token": "abc123def456...",
    "password": "newSecurePassword456"
}

// Response (200 OK)
{
    "message": "Password reset successfully"
}`} />

              {/* ── Auth Middleware Usage ─────────────────────────────── */}
              <h2 id="auth-middleware">Auth Middleware Usage</h2>
              <p>
                Apply the <code>Auth</code> middleware to any route group that requires authentication.
                See the <Link href="/docs/backend/middleware" className="text-primary hover:underline">Middleware</Link> page
                for the full implementation.
              </p>
              <CodeBlock filename="routes.go" code={`// Protected routes -- any authenticated user
protected := r.Group("/api")
protected.Use(middleware.Auth(db, authService))
{
    protected.GET("/auth/me", authHandler.Me)
    protected.POST("/auth/logout", authHandler.Logout)
    protected.GET("/posts", postHandler.List)
}

// Admin routes -- admin role required
admin := r.Group("/api")
admin.Use(middleware.Auth(db, authService))
admin.Use(middleware.RequireRole("admin"))
{
    admin.GET("/users", userHandler.List)
    admin.DELETE("/users/:id", userHandler.Delete)
}`} />

              {/* ── Role-Based Access ─────────────────────────────── */}
              <h2 id="roles">Role-Based Access Control</h2>
              <p>
                Grit defines three built-in roles. You can extend these by adding new
                constants to the User model.
              </p>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Role</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Constant</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Access Level</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">admin</td>
                      <td className="px-4 py-2.5 font-mono text-xs">models.RoleAdmin</td>
                      <td className="px-4 py-2.5">Full access to all resources, user management, admin panel</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">editor</td>
                      <td className="px-4 py-2.5 font-mono text-xs">models.RoleEditor</td>
                      <td className="px-4 py-2.5">Can create and edit content, limited admin access</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">user</td>
                      <td className="px-4 py-2.5 font-mono text-xs">models.RoleUser</td>
                      <td className="px-4 py-2.5">Default role, can access own data only</td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <CodeBlock filename="models/user.go" code={`// Built-in roles
const (
    RoleAdmin  = "admin"
    RoleEditor = "editor"
    RoleUser   = "user"
)

// Add custom roles:
const (
    RoleManager   = "manager"
    RoleModerator = "moderator"
)`} />

              {/* ── Token Storage (Frontend) ─────────────────────────────── */}
              <h2 id="token-storage">Token Storage on the Frontend</h2>

              <div className="p-4 rounded-lg border border-destructive/30 bg-destructive/5 mb-6">
                <p className="text-sm text-foreground/85 mb-0">
                  <strong className="text-destructive">Do not store tokens in <code>localStorage</code> for the web client.</strong>{' '}
                  Anything readable from JavaScript is reachable by any XSS vector — a compromised
                  npm dependency, a stored XSS bug in a comment field, or a browser extension.
                  Tokens in <code>localStorage</code> are persistent, unscoped, and exfiltrate-able
                  with a single line of script. The OWASP guidance (and ours) is to put auth
                  cookies out of JavaScript&apos;s reach.
                </p>
              </div>

              <p>Pick the storage model that matches your client:</p>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Client</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Token storage</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Why</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-medium text-foreground">Web (Next.js)</td>
                      <td className="px-4 py-2.5"><code>httpOnly</code>, <code>Secure</code>, <code>SameSite=Lax</code> cookies set by the API</td>
                      <td className="px-4 py-2.5">XSS cannot read them; the browser attaches them automatically.</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-medium text-foreground">Mobile (Expo)</td>
                      <td className="px-4 py-2.5"><code>expo-secure-store</code> (iOS Keychain / Android Keystore)</td>
                      <td className="px-4 py-2.5">Hardware-backed, not readable by other apps or React Native bridges.</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-medium text-foreground">Desktop (Wails)</td>
                      <td className="px-4 py-2.5">OS keychain via Go binding (<code>keyring</code>)</td>
                      <td className="px-4 py-2.5">Same threat model as mobile; never the renderer.</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <h3 id="token-storage-cookies">Web: cookies set by the API</h3>

              <div className="p-4 rounded-lg border border-emerald-500/30 bg-emerald-500/5 mb-6">
                <p className="text-sm text-foreground/85 mb-0">
                  <strong className="text-emerald-400">This is on by default.</strong>{' '}
                  Every Grit project scaffolded with v3.25.3+ already sets{' '}
                  <code>grit_access</code> + <code>grit_refresh</code> on
                  Login / Register / Refresh / TOTP-verify and clears them on
                  Logout. The <code>middleware.Auth</code> chain reads cookies
                  first and falls back to the <code>Authorization</code> header
                  for native bearer clients. You don&apos;t have to wire any of
                  this yourself.
                </p>
              </div>

              <p>
                The API sets two cookies on login/register/refresh: <code>grit_access</code> (short-lived)
                and <code>grit_refresh</code> (long-lived, scoped to <code>/api/auth</code>). Both are{' '}
                <code>HttpOnly</code> so JavaScript cannot read them, <code>Secure</code> on HTTPS
                so they only travel over TLS, and <code>SameSite=Lax</code> so the CSRF surface
                is limited to top-level navigations.
              </p>
              <p>Helpers live on the auth service for use in your own handlers:</p>
              <CodeBlock language="go" filename="internal/services/auth.go (already in scaffold)" code={`// SetAuthCookies writes the token pair as HttpOnly cookies.
// Called from Register / Login / Refresh / TOTP verify.
func (s *AuthService) SetAuthCookies(c *gin.Context, pair *TokenPair) { ... }

// ClearAuthCookies expires both cookies. Called from Logout.
func (s *AuthService) ClearAuthCookies(c *gin.Context) { ... }`} />

              <p>
                The auth middleware reads cookies first, then the Authorization header — same
                <code>middleware.Auth(db, authService)</code> covers both flows:
              </p>
              <CodeBlock language="go" filename="internal/middleware/auth.go (already in scaffold)" code={`token := ""
if cookieValue, err := c.Cookie("grit_access"); err == nil && cookieValue != "" {
    token = cookieValue
} else if authHeader := c.GetHeader("Authorization"); authHeader != "" {
    // Bearer fallback for native mobile / desktop clients
    parts := strings.SplitN(authHeader, " ", 2)
    if len(parts) == 2 && parts[0] == "Bearer" {
        token = parts[1]
    }
}`} />

              <p>
                The Next.js client doesn&apos;t touch tokens at all — the browser handles them.
                Use <code>credentials: &apos;include&apos;</code> on every request:
              </p>
              <CodeBlock language="typescript" filename="apps/web/lib/api-client.ts" code={`import axios from 'axios'

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  withCredentials: true,   // <- the browser sends grit_access / grit_refresh automatically
})

// Auto-refresh on 401 — note: no token reading, no localStorage.
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      try {
        // The browser sends grit_refresh; the API sets a new grit_access cookie.
        await api.post('/api/auth/refresh')
        return api(originalRequest)
      } catch {
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  },
)

export default api`} />

              <div className="p-4 rounded-lg border border-emerald-500/30 bg-emerald-500/5 mb-6">
                <p className="text-sm text-foreground/85 mb-0">
                  <strong className="text-emerald-400">CSRF is on by default for the cookie flow.</strong>{' '}
                  Cookie-auth APIs are vulnerable to CSRF (a malicious site forging a state-changing
                  request from a logged-in user). The scaffolded{' '}
                  <code>middleware.AutoCSRF()</code> is wired globally and enforces a double-submit
                  CSRF token <em>only</em> when the request carries the <code>grit_access</code>{' '}
                  cookie. Bearer-token requests (mobile / desktop) pass through with no header
                  required — they aren&apos;t CSRF-vulnerable because browsers never auto-send a
                  Bearer header across origins.
                </p>
              </div>

              <p>
                The SPA reads the token from the <code>grit_csrf</code> cookie (which is{' '}
                <strong>not</strong> HttpOnly, so JS can read it) and echoes it back in the{' '}
                <code>X-CSRF-Token</code> header. With Axios:
              </p>
              <CodeBlock language="ts" filename="apps/web/lib/api-client.ts" code={`import axios from 'axios'

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  withCredentials: true,
})

// Echo the CSRF token from grit_csrf cookie into the X-CSRF-Token header.
// The middleware accepts this on every cookie-authenticated mutation.
api.interceptors.request.use((config) => {
  const m = document.cookie.match(/(?:^|; )grit_csrf=([^;]+)/)
  if (m) config.headers['X-CSRF-Token'] = decodeURIComponent(m[1])
  return config
})

export default api`} />
              <p>
                The first GET request (any GET) seeds the cookie; subsequent mutations carry the
                header. No bootstrap endpoint needed.
              </p>

              <h3 id="token-storage-native">Mobile + Desktop: bearer header from secure store</h3>
              <p>
                Native clients can&apos;t use HttpOnly cookies cleanly across all platforms. Use
                the secure OS-backed store and the <code>Authorization: Bearer</code> header path:
              </p>
              <CodeBlock language="typescript" filename="apps/mobile/lib/auth-store.ts" code={`import * as SecureStore from 'expo-secure-store'

const ACCESS = 'grit_access'
const REFRESH = 'grit_refresh'

export const saveTokens = async (access: string, refresh: string) => {
  await SecureStore.setItemAsync(ACCESS, access)
  await SecureStore.setItemAsync(REFRESH, refresh)
}

export const loadTokens = async () => ({
  access: await SecureStore.getItemAsync(ACCESS),
  refresh: await SecureStore.getItemAsync(REFRESH),
})

export const clearTokens = async () => {
  await SecureStore.deleteItemAsync(ACCESS)
  await SecureStore.deleteItemAsync(REFRESH)
}`} />
              <p>
                The <code>Authorization</code> header path on the API stays for these clients.
                Desktop (Wails) uses an equivalent OS-keychain binding from Go and exposes
                <code>saveTokens / loadTokens</code> to the React frontend via Wails bindings.
              </p>

              <h3 id="token-storage-summary">Summary — never use localStorage for tokens</h3>
              <ul>
                <li><strong>Web:</strong> HttpOnly cookies (default-on). No JS touches the access token. CSRF is auto-enforced via <code>AutoCSRF()</code>.</li>
                <li><strong>Mobile:</strong> <code>expo-secure-store</code> + bearer header.</li>
                <li><strong>Desktop:</strong> OS keychain + bearer header.</li>
                <li><strong>Never:</strong> <code>localStorage</code> / <code>sessionStorage</code> for auth tokens.</li>
              </ul>

              <h3 id="token-storage-defence-in-depth">Defence-in-depth — what else is on by default</h3>
              <ul>
                <li>
                  <strong>SecurityHeaders middleware</strong> — strict CSP (no inline script),
                  X-Frame-Options DENY, X-Content-Type-Options nosniff, Referrer-Policy
                  strict-origin-when-cross-origin, Permissions-Policy locking down camera /
                  mic / geolocation / payments / USB, COOP + CORP for Spectre isolation, and
                  HSTS on HTTPS. Globally applied.
                </li>
                <li>
                  <strong>Sentinel AuthShield</strong> — brute-force lockout on{' '}
                  <code>/api/auth/login</code> with progressive backoff. Default-on whenever
                  the Sentinel suite is enabled (which it is by default in fresh scaffolds).
                </li>
                <li>
                  <strong>Sentinel rate limiting</strong> — 5 requests / 15 minutes per IP on
                  <code>/api/auth/login</code> and 3 / 15 min on <code>/api/auth/register</code>{' '}
                  in production. Dev gets relaxed limits so testing doesn&apos;t lock you out.
                </li>
                <li>
                  <strong>WAF</strong> — Sentinel runs in block mode in production, log mode in
                  dev. Catches injection patterns at the edge before they reach handlers.
                </li>
                <li>
                  <strong>safefetch package</strong> — use{' '}
                  <code>safefetch.Client</code> for any URL the user supplies (webhooks, OG-image
                  preview, OEmbed expansion). Blocks private IP ranges + cloud metadata hostnames
                  and re-validates the resolved IP at TCP-connect time to defeat DNS rebinding.
                </li>
              </ul>

              {/* ── Password Hashing ─────────────────────────────── */}
              <h2 id="password-hashing">Password Hashing</h2>
              <p>
                Passwords are hashed using <strong>bcrypt</strong> with the default cost factor (10).
                Hashing happens automatically via the GORM <code>BeforeCreate</code> hook on the
                User model. Passwords are never stored in plain text and are never returned in
                API responses (the <code>Password</code> field uses <code>json:&quot;-&quot;</code>).
              </p>
              <CodeBlock filename="models/user.go" code={`// Password field -- never included in JSON responses
Password string \`gorm:"size:255;not null" json:"-"\`

// Automatically hash on create
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.Password != "" {
        hashedPassword, err := bcrypt.GenerateFromPassword(
            []byte(u.Password), bcrypt.DefaultCost,
        )
        if err != nil {
            return err
        }
        u.Password = string(hashedPassword)
    }
    return nil
}

// Verify password during login
func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword(
        []byte(u.Password), []byte(password),
    )
    return err == nil
}`} />

              {/* ── Token Generation ─────────────────────────────── */}
              <h2 id="token-generation">Token Generation</h2>
              <p>
                The <code>AuthService</code> handles all token operations. It uses the
                <code>golang-jwt/jwt/v5</code> library with HMAC-SHA256 signing.
              </p>
              <CodeBlock filename="services/auth.go" code={`// GenerateTokenPair creates access + refresh tokens.
func (s *AuthService) GenerateTokenPair(
    userID uint, email, role string,
) (*TokenPair, error) {
    accessToken, expiresAt, err := s.generateToken(
        userID, email, role, s.AccessExpiry,
    )
    if err != nil {
        return nil, fmt.Errorf("generating access token: %w", err)
    }

    refreshToken, _, err := s.generateToken(
        userID, email, role, s.RefreshExpiry,
    )
    if err != nil {
        return nil, fmt.Errorf("generating refresh token: %w", err)
    }

    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresAt:    expiresAt,
    }, nil
}

func (s *AuthService) generateToken(
    userID uint, email, role string, expiry time.Duration,
) (string, int64, error) {
    expiresAt := time.Now().Add(expiry)

    claims := &Claims{
        UserID: userID,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(s.Secret))
    if err != nil {
        return "", 0, err
    }

    return tokenString, expiresAt.Unix(), nil
}`} />

              {/* ── Password Reset Token ─────────────────────────────── */}
              <h2 id="reset-token">Password Reset Tokens</h2>
              <p>
                Password reset tokens are cryptographically random 32-byte hex strings.
                They are generated using Go&apos;s <code>crypto/rand</code> package, which is
                secure for this purpose.
              </p>
              <CodeBlock filename="services/auth.go" code={`// GenerateResetToken creates a random hex token for password resets.
func GenerateResetToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", fmt.Errorf("generating reset token: %w", err)
    }
    return hex.EncodeToString(bytes), nil
}

// Output example: "a3f4b2c1e5d6f7890123456789abcdef..."
// (64 hex characters = 32 bytes of randomness)`} />

              {/* ── Two-Factor Authentication ─────────────────────────────── */}
              <h2 id="two-factor">Two-Factor Authentication (TOTP)</h2>
              <p>
                Every Grit project includes a complete two-factor authentication system using
                TOTP (Time-based One-Time Passwords). It works with any authenticator app:
                Google Authenticator, Authy, 1Password, Bitwarden, etc.
              </p>

              <h3 id="totp-how-it-works">How It Works</h3>
              <CodeBlock filename="TOTP login flow" code={`
  Client                           Grit API
    |                                  |
    |  POST /api/auth/login            |
    |  { email, password }             |
    | -------------------------------->|
    |                                  |  Validate password ✓
    |                                  |  Check: TOTP enabled?
    |                                  |  Check: Trusted device cookie?
    |                                  |
    |  If TOTP required:               |
    |  { totp_required, pending_token }|
    | <--------------------------------|
    |                                  |
    |  POST /api/auth/totp/verify      |
    |  { pending_token, code, trust }  |
    | -------------------------------->|
    |                                  |  Validate TOTP code ✓
    |                                  |  (Optional) Set trusted device
    |  { user, tokens }                |
    | <--------------------------------|
`} />

              <p>
                If the user has 2FA enabled and no trusted device cookie, the login endpoint
                returns a short-lived <code>pending_token</code> (5 minutes) instead of JWT tokens.
                The client then redirects to a TOTP verification page.
              </p>

              <h3 id="totp-endpoints">TOTP Endpoints</h3>
              <div className="overflow-x-auto mb-6">
                <table>
                  <thead>
                    <tr>
                      <th>Method</th>
                      <th>Endpoint</th>
                      <th>Auth</th>
                      <th>Purpose</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr><td><code>POST</code></td><td><code>/api/auth/totp/setup</code></td><td>JWT</td><td>Generate secret + QR code URI</td></tr>
                    <tr><td><code>POST</code></td><td><code>/api/auth/totp/enable</code></td><td>JWT</td><td>Verify initial code, activate 2FA, get backup codes</td></tr>
                    <tr><td><code>POST</code></td><td><code>/api/auth/totp/verify</code></td><td>Public</td><td>Exchange pending token + TOTP code for JWT</td></tr>
                    <tr><td><code>POST</code></td><td><code>/api/auth/totp/backup-codes/verify</code></td><td>Public</td><td>Use backup code during login</td></tr>
                    <tr><td><code>POST</code></td><td><code>/api/auth/totp/disable</code></td><td>JWT</td><td>Turn off 2FA (requires password)</td></tr>
                    <tr><td><code>GET</code></td><td><code>/api/auth/totp/status</code></td><td>JWT</td><td>Check 2FA status, backup codes remaining</td></tr>
                    <tr><td><code>POST</code></td><td><code>/api/auth/totp/backup-codes</code></td><td>JWT</td><td>Regenerate backup codes</td></tr>
                    <tr><td><code>DELETE</code></td><td><code>/api/auth/totp/trusted-devices</code></td><td>JWT</td><td>Revoke all trusted devices</td></tr>
                  </tbody>
                </table>
              </div>

              <h3 id="totp-setup-flow">Enabling 2FA (User Flow)</h3>
              <CodeBlock language="typescript" filename="Enable TOTP" code={`// Step 1: Get the secret and QR code URI
const { data } = await api.post('/api/auth/totp/setup')
// data.secret = "JBSWY3DPEHPK3PXP..."
// data.uri    = "otpauth://totp/MyApp:user@email.com?secret=..."
// → Show QR code to user (use a QR library to render data.uri)

// Step 2: User scans QR code, enters the 6-digit code from their app
const { data: result } = await api.post('/api/auth/totp/enable', {
  secret: data.secret,
  code: '123456'  // from authenticator app
})
// result.enabled = true
// result.backup_codes = ["A1B2C3D4", "E5F6G7H8", ...]
// → Show backup codes to user (they must save these!)`} />

              <h3 id="totp-login-flow">Login with 2FA (Client Flow)</h3>
              <CodeBlock language="typescript" filename="Login with TOTP" code={`// Step 1: Normal login
const { data } = await api.post('/api/auth/login', { email, password })

if (data.totp_required) {
  // Step 2: Redirect to TOTP verification page
  // Store the pending token temporarily
  const pendingToken = data.pending_token

  // Step 3: User enters 6-digit code from authenticator app
  const { data: result } = await api.post('/api/auth/totp/verify', {
    pending_token: pendingToken,
    code: '123456',
    trust_device: true  // optional: remember this device for 30 days
  })
  // result.user = { ... }
  // result.tokens = { access_token, refresh_token }
} else {
  // No 2FA — normal login, tokens already returned
  // data.user = { ... }
  // data.tokens = { access_token, refresh_token }
}`} />

              <h3 id="backup-codes">Backup Codes</h3>
              <p>
                When 2FA is enabled, 10 one-time-use backup codes are generated. Each code is
                individually bcrypt-hashed before storage. When a user enters a backup code during
                login, the used code is permanently removed from the database.
              </p>
              <CodeBlock language="typescript" filename="Using a backup code" code={`// During login, if user lost their authenticator app:
const { data } = await api.post('/api/auth/totp/backup-codes/verify', {
  pending_token: pendingToken,
  code: 'A1B2C3D4',  // one of the saved backup codes
  trust_device: false
})
// data.backup_codes_remaining = 9  (one code was consumed)`} />

              <h3 id="trusted-devices">Trusted Devices</h3>
              <p>
                When <code>trust_device: true</code> is sent during TOTP verification, an HttpOnly
                cookie (<code>totp_trusted</code>) is set with a random token. The SHA-256 hash of this
                token is stored in the database with the user&apos;s IP and user agent. Trusted devices
                last 30 days with sliding expiry — each successful login refreshes the timer.
              </p>
              <p>
                Users can revoke all trusted devices:
              </p>
              <CodeBlock language="typescript" filename="Revoke trusted devices" code={`await api.delete('/api/auth/totp/trusted-devices')
// All trusted device cookies are now invalid`} />

              <h3 id="totp-implementation">Implementation Details</h3>
              <ul>
                <li><strong>Algorithm:</strong> HMAC-SHA1 (RFC 6238 / RFC 4226)</li>
                <li><strong>Code length:</strong> 6 digits</li>
                <li><strong>Period:</strong> 30 seconds</li>
                <li><strong>Clock skew:</strong> &plusmn;1 window tolerance (90 second total window)</li>
                <li><strong>Secret:</strong> 20 random bytes, base32-encoded (no padding)</li>
                <li><strong>Pending tokens:</strong> 32 random bytes, hex-encoded, SHA-256 hashed for DB, expires in 5 minutes</li>
                <li><strong>Backup codes:</strong> 8-character hex codes, individually bcrypt-hashed, one-time use</li>
                <li><strong>Trusted device tokens:</strong> 32 random bytes, SHA-256 hashed, 30-day sliding expiry</li>
                <li><strong>Dependencies:</strong> Zero external — uses only Go standard library + <code>golang.org/x/crypto/bcrypt</code></li>
              </ul>

              {/* ── Configuration ─────────────────────────────── */}
              <h2 id="configuration">Auth Configuration</h2>
              <p>
                All authentication settings are configured via environment variables:
              </p>
              <CodeBlock language="bash" filename=".env" code={`# Required
JWT_SECRET=change-this-to-a-long-random-string

# Optional (defaults shown)
JWT_ACCESS_EXPIRY=15m      # Go duration format
JWT_REFRESH_EXPIRY=168h    # 7 days

# TOTP (Two-Factor Authentication)
TOTP_ISSUER=MyApp          # App name shown in authenticator apps (defaults to APP_NAME)`} />
              <div className="p-4 rounded-lg border border-destructive/20 bg-destructive/5 mb-6">
                <p className="text-sm text-foreground/80 mb-0">
                  <strong>Important:</strong> The <code>JWT_SECRET</code> environment variable is required.
                  The server will not start without it. Use a random string of at least 32 characters
                  in production.
                </p>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 mt-10 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/backend/middleware" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Middleware
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/backend/response-format" className="gap-1.5">
                  API Response Format
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
