import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

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
                login, token refresh, logout, password reset, and role-based access control --
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
              <p>
                The Next.js frontend stores tokens and includes them in API requests using
                an Axios interceptor. The recommended pattern:
              </p>
              <CodeBlock language="typescript" filename="apps/web/lib/api-client.ts" code={`import axios from 'axios';

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
});

// Attach access token to every request
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = \`Bearer \${token}\`;
  }
  return config;
});

// Auto-refresh on 401
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = localStorage.getItem('refresh_token');
        const { data } = await axios.post(
          \`\${api.defaults.baseURL}/api/auth/refresh\`,
          { refresh_token: refreshToken },
        );

        const { access_token, refresh_token } = data.data.tokens;
        localStorage.setItem('access_token', access_token);
        localStorage.setItem('refresh_token', refresh_token);

        originalRequest.headers.Authorization = \`Bearer \${access_token}\`;
        return api(originalRequest);
      } catch {
        // Refresh failed -- redirect to login
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        window.location.href = '/login';
      }
    }

    return Promise.reject(error);
  },
);

export default api;`} />

              <div className="p-4 rounded-lg border border-primary/20 bg-primary/5 mb-6">
                <p className="text-sm text-foreground/80 mb-0">
                  <strong>Security note:</strong> Storing tokens in <code>localStorage</code> is
                  acceptable for most applications. For higher security, consider using
                  <code>httpOnly</code> cookies by modifying the login/refresh endpoints to set
                  cookies instead of returning tokens in the JSON body.
                </p>
              </div>

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

              {/* ── Configuration ─────────────────────────────── */}
              <h2 id="configuration">Auth Configuration</h2>
              <p>
                All authentication settings are configured via environment variables:
              </p>
              <CodeBlock language="bash" filename=".env" code={`# Required
JWT_SECRET=change-this-to-a-long-random-string

# Optional (defaults shown)
JWT_ACCESS_EXPIRY=15m      # Go duration format
JWT_REFRESH_EXPIRY=168h    # 7 days`} />
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
