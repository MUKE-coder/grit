import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Authentication & Authorization — Grit Web Course',
  description: 'Learn the complete auth system in Grit: JWT tokens, bcrypt hashing, role-based access control, two-factor authentication with TOTP, OAuth2 social login, and frontend auth flows.',
}

export default function AuthenticationCourse() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <Link href="/courses/grit-web" className="hover:text-foreground transition-colors">Grit Web</Link>
          <span>/</span>
          <span className="text-foreground">Authentication & Authorization</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Course 3 of 8</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">16 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Authentication & Authorization
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            In this course, you will learn how Grit handles security from the ground up: JWT-based
            authentication, bcrypt password hashing, role-based access control, two-factor authentication
            with TOTP, OAuth2 social login, and how the frontend manages auth state. By the end, you will
            understand every layer of Grit{"'"}s auth system.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What is Authentication? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is Authentication?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Security in web applications comes down to two fundamental questions:
            {' '}<strong className="text-foreground">Who are you?</strong> and <strong className="text-foreground">What are you allowed to do?</strong>
            {' '}These are handled by two separate systems that work together.
          </p>

          <Definition term="Authentication">
            Verifying who a user is. When you log in with an email and password, the server checks your
            credentials against the database. If they match, you are authenticated — the server now knows
            who you are. Think of it like showing your ID at the door.
          </Definition>

          <Definition term="Authorization">
            Determining what an authenticated user is allowed to do. An ADMIN can delete users, but a
            regular USER cannot. A post owner can edit their own posts, but not someone else{"'"}s.
            Authorization happens <strong className="text-foreground">after</strong> authentication — you must
            know WHO someone is before you can decide WHAT they can do.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit handles both: <strong className="text-foreground">JWT tokens</strong> for authentication
            and <strong className="text-foreground">roles</strong> for authorization. Every API request goes
            through this two-step check: first, is the token valid? Second, does this user{"'"}s role allow
            this action?
          </p>

          <Challenge number={1} title="Authentication vs Authorization">
            <p>In your own words, explain the difference between authentication and authorization. Give an
            example from a real app — for instance, on Instagram, anyone with an account can view posts
            (authenticated), but only the post owner can delete it (authorized). Come up with your own
            example from an app you use daily.</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: How JWT Authentication Works ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">How JWT Authentication Works</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit uses <strong className="text-foreground">JSON Web Tokens (JWTs)</strong> for authentication.
            Here is the complete flow, step by step:
          </p>

          <ol className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">1.</span> User sends email + password to <Code>POST /api/auth/register</Code> or <Code>POST /api/auth/login</Code></li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">2.</span> Server verifies credentials against the database (password hashed with bcrypt)</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">3.</span> Server creates two tokens: an <strong className="text-foreground">access token</strong> (short-lived, 15 minutes) and a <strong className="text-foreground">refresh token</strong> (long-lived, 7 days)</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">4.</span> Frontend stores both tokens</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">5.</span> Every API request includes the access token in the <Code>Authorization</Code> header</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">6.</span> When the access token expires, the frontend uses the refresh token to get a new one</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">7.</span> When the refresh token expires, the user must log in again</li>
          </ol>

          <Definition term="JWT (JSON Web Token)">
            A compact, URL-safe token that contains encoded information about the user (like their ID and
            role). It is signed with a secret key so the server can verify it was not tampered with.
            The format is three parts separated by dots: <Code>header.payload.signature</Code>. You can
            decode a JWT at <a href="https://jwt.io" target="_blank" rel="noreferrer" className="text-primary hover:underline">jwt.io</a> to
            see its contents.
          </Definition>

          <Definition term="Access Token">
            A short-lived JWT (15 minutes in Grit) that is included in every API request. The short expiry
            limits the damage if the token is stolen — an attacker only has 15 minutes before it becomes
            useless.
          </Definition>

          <Definition term="Refresh Token">
            A long-lived token (7 days in Grit) used <strong className="text-foreground">only</strong> to get
            a new access token. It is stored securely and sent only to the <Code>/api/auth/refresh</Code> endpoint.
            This way, the user does not have to log in every 15 minutes.
          </Definition>

          <Definition term="bcrypt">
            A password hashing algorithm. Grit never stores plain text passwords. When you register,
            bcrypt creates a one-way hash — a scrambled version that cannot be reversed back into the
            original password. When you login, bcrypt compares your input against the stored hash. Even
            if the database is stolen, attackers cannot read the passwords.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here is what happens when you register a new user:
          </p>

          <CodeBlock filename="POST /api/auth/register">
{`// Request
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securePassword123"
}

// Response
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "role": "USER"
    }
  }
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Notice the response includes <strong className="text-foreground">both tokens</strong> and the
            user object. The password is never returned. Logging in works the same way:
          </p>

          <CodeBlock filename="POST /api/auth/login">
{`{
  "email": "john@example.com",
  "password": "securePassword123"
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Once you have the access token, every subsequent API request must include it in the
            {' '}<Code>Authorization</Code> header:
          </p>

          <CodeBlock filename="GET /api/users">
{`Authorization: Bearer eyJhbGciOiJIUzI1NiIs...`}
          </CodeBlock>

          <Definition term="Bearer Token">
            The format for sending JWTs in HTTP headers. {'"'}Bearer{'"'} means {'"'}the person bearing
            (carrying) this token is authorized.{'"'} The full header format is
            {' '}<Code>Authorization: Bearer &lt;token&gt;</Code>.
          </Definition>

          <Challenge number={2} title="Register a User">
            <p>Open the API docs at <Code>localhost:8080/docs</Code>. Find the register endpoint.
            Register a new user with your name and email. Copy the <Code>access_token</Code> from
            the response — you will need it for the next challenge.</p>
          </Challenge>

          <Challenge number={3} title="Use the Access Token">
            <p>Use the <Code>access_token</Code> you copied to call <Code>GET /api/users</Code>.
            In the API docs, paste it in the Authorization field. Does it return the user list?
            Now try calling the same endpoint <strong className="text-foreground">without</strong> a
            token. What error code and message do you get?</p>
          </Challenge>

          <Challenge number={4} title="Decode a JWT">
            <p>Copy your access token and paste it at <a href="https://jwt.io" target="_blank" rel="noreferrer" className="text-primary hover:underline">jwt.io</a>.
            Look at the payload section — can you see your user ID and role? What is the <Code>exp</Code> field?
            (Hint: it is a Unix timestamp for when the token expires.)</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Where Auth Lives in the Code ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Where Auth Lives in the Code</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit{"'"}s auth system is spread across four key files. Understanding where each piece lives
            helps you customize the auth behavior:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>apps/api/internal/handlers/auth.go</Code> — Register, Login, RefreshToken, and GetMe handlers</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>apps/api/internal/services/auth_service.go</Code> — Password hashing, token generation, user lookup</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>apps/api/internal/middleware/auth.go</Code> — JWT verification middleware</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>apps/api/internal/routes/routes.go</Code> — Public vs protected route groups</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The most important concept is the split between <strong className="text-foreground">public</strong> and
            {' '}<strong className="text-foreground">protected</strong> routes:
          </p>

          <CodeBlock filename="apps/api/internal/routes/routes.go">
{`// Public routes — no authentication required
public := r.Group("/api/auth")
public.POST("/register", authHandler.Register)
public.POST("/login", authHandler.Login)
public.POST("/refresh", authHandler.RefreshToken)

// Protected routes — JWT required
protected := r.Group("/api")
protected.Use(middleware.Auth(cfg.JWTSecret))
protected.GET("/users", userHandler.List)
protected.GET("/me", authHandler.GetMe)`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Routes inside the {'"'}protected{'"'} group require a valid JWT. The <Code>Auth</Code> middleware
            runs <strong className="text-foreground">before</strong> the handler and checks the token. If the
            token is missing, expired, or tampered with, the middleware returns <Code>401 Unauthorized</Code> immediately
            and the handler never executes.
          </p>

          <Definition term="Route Group">
            A collection of routes that share a common prefix and/or middleware. In Grit, {'"'}public{'"'} routes
            (like login and register) have no auth requirement, while {'"'}protected{'"'} routes require a valid
            JWT token. This keeps the auth logic in one place instead of checking it in every handler.
          </Definition>

          <Tip>
            The <Code>middleware.Auth()</Code> function extracts the user ID and role from the JWT and
            attaches them to the request context. Any handler in the protected group can then access
            {' '}<Code>c.GetUint({'"'}userID{'"'})</Code> and <Code>c.GetString({'"'}userRole{'"'})</Code> without
            re-parsing the token.
          </Tip>

          <Challenge number={5} title="Explore the Auth Code">
            <p>Open <Code>apps/api/internal/routes/routes.go</Code>. Find the public and protected
            groups. List 3 public routes and 3 protected routes you find in your project.</p>
          </Challenge>

          <Challenge number={6} title="Read the Middleware">
            <p>Open <Code>apps/api/internal/middleware/auth.go</Code>. Find the line where the JWT token
            is extracted from the request header. What happens if the token is missing? What happens if
            the token is expired? Trace the code path for each case.</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Roles and RBAC ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Roles and Role-Based Access Control</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Authentication tells you <em>who</em> the user is. But that is only half the story. You also
            need to control <em>what</em> each user can do. Grit solves this with three built-in roles:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><span className="text-primary font-bold">ADMIN</span> — Full access to everything: manage users, system settings, all resources</li>
            <li className="flex gap-2"><span className="text-primary font-bold">EDITOR</span> — Can create, edit, and publish content, but cannot manage users or system settings</li>
            <li className="flex gap-2"><span className="text-primary font-bold">USER</span> — Basic access: their own profile and public content</li>
          </ul>

          <Definition term="RBAC (Role-Based Access Control)">
            A security model where permissions are assigned to roles, and roles are assigned to users.
            Instead of checking {'"'}can user #42 delete posts?{'"'}, you check {'"'}does this user have
            the ADMIN role?{'"'} This is much simpler to manage — when a new admin joins, you assign
            the ADMIN role instead of configuring individual permissions.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here is how role-restricted routes look in Grit:
          </p>

          <CodeBlock filename="apps/api/internal/routes/routes.go">
{`// Only ADMIN can access these
admin := protected.Group("/admin")
admin.Use(middleware.RequireRole("ADMIN"))
admin.DELETE("/users/:id", userHandler.Delete)
admin.PUT("/users/:id/role", userHandler.UpdateRole)`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>RequireRole</Code> middleware checks the user{"'"}s role (which was set by the
            {' '}<Code>Auth</Code> middleware earlier). If the role does not match, the request is rejected
            with <Code>403 Forbidden</Code>.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here is the User model with the role field:
          </p>

          <CodeBlock filename="apps/api/internal/models/user.go">
{`type User struct {
    gorm.Model
    Name     string \`gorm:"not null" json:"name"\`
    Email    string \`gorm:"uniqueIndex;not null" json:"email"\`
    Password string \`gorm:"not null" json:"-"\`
    Role     string \`gorm:"default:USER" json:"role"\`
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Two important details: <Code>json:{'"'}-{'"'}</Code> means the Password field is
            {' '}<strong className="text-foreground">never</strong> included in API responses — it is invisible
            to the frontend. The Role field defaults to <Code>USER</Code> for all new registrations, so
            no one can register themselves as an ADMIN.
          </p>

          <Challenge number={7} title="Explore Your User Role">
            <p>Open GORM Studio at <Code>localhost:8080/studio</Code>. Find your user in the
            {' '}<Code>users</Code> table. What role do they have? Change it to <Code>ADMIN</Code> directly
            in the database. Refresh the admin panel — do you see more options now?</p>
          </Challenge>

          <Challenge number={8} title="Test Role Restrictions">
            <p>Create a second user through the web app (use a different email). This new user will have
            the default <Code>USER</Code> role. Try accessing the admin panel with this account. What
            happens? Can you access the user management page?</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Adding Custom Roles ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Adding Custom Roles</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Three roles might not be enough for your app. A forum might need a <Code>MODERATOR</Code>.
            A marketplace might need a <Code>VENDOR</Code>. Grit makes adding custom roles easy:
          </p>

          <CodeBlock filename="Terminal">
{`grit add role MODERATOR`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This single command updates <strong className="text-foreground">four places</strong> across your stack:
          </p>

          <ol className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">1.</span> <strong className="text-foreground">Go models</strong> — adds <Code>MODERATOR</Code> to the role constants</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">2.</span> <strong className="text-foreground">TypeScript types</strong> — adds {'"'}MODERATOR{'"'} to the role union type</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">3.</span> <strong className="text-foreground">Zod schemas</strong> — adds {'"'}MODERATOR{'"'} to the role validation</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">4.</span> <strong className="text-foreground">Admin UI</strong> — adds MODERATOR as an option in role dropdowns</li>
          </ol>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You can also restrict generated resources to specific roles:
          </p>

          <CodeBlock filename="Terminal">
{`grit generate resource Report --fields "title:string,content:text" --roles "ADMIN,MODERATOR"`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This generates routes that only ADMIN and MODERATOR can access. A regular USER trying to
            reach <Code>/api/reports</Code> will get <Code>403 Forbidden</Code>.
          </p>

          <Challenge number={9} title="Add a Custom Role">
            <p>Run <Code>grit add role MODERATOR</Code>. Then verify all four places were updated:
            (1) check the Go model constants, (2) check the TypeScript types file, (3) check the Zod
            schemas, and (4) open the admin user form — is MODERATOR an option in the role dropdown?</p>
          </Challenge>

          <Challenge number={10} title="Role-Restricted Resources">
            <p>Generate a <Code>Report</Code> resource restricted to ADMIN and MODERATOR:
            {' '}<Code>grit generate resource Report --fields {'"'}title:string,content:text{'"'} --roles {'"'}ADMIN,MODERATOR{'"'}</Code>.
            Log in as a regular USER and try to access <Code>/api/reports</Code>. What error do you get?
            Now change the user{"'"}s role to MODERATOR and try again.</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Two-Factor Authentication (TOTP) ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Two-Factor Authentication (TOTP)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Passwords alone are not enough. They can be guessed, phished, or leaked in data breaches.
            TOTP adds a <strong className="text-foreground">second layer of security</strong>: even if someone
            steals your password, they cannot log in without the code from your authenticator app.
          </p>

          <Definition term="Two-Factor Authentication (2FA)">
            Requiring two different things to log in: (1) something you <strong className="text-foreground">know</strong> (your
            password) and (2) something you <strong className="text-foreground">have</strong> (your authenticator
            app on your phone). Even if your password is compromised, the attacker needs your phone too.
          </Definition>

          <Definition term="TOTP (Time-Based One-Time Password)">
            A 6-digit code that changes every 30 seconds, generated by an authenticator app like Google
            Authenticator, Authy, or 1Password. It is based on{' '}
            <a href="https://datatracker.ietf.org/doc/html/rfc6238" target="_blank" rel="noreferrer" className="text-primary hover:underline">RFC 6238</a>
            {' '}— a shared secret combined with the current time produces a unique code. Both the server
            and your app generate the same code independently, so no network request is needed.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here is how Grit{"'"}s TOTP system works:
          </p>

          <ol className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">1.</span> User calls <Code>POST /api/totp/setup</Code> — server generates a secret key and returns a QR code URL</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">2.</span> User scans the QR code with their authenticator app (Google Authenticator, Authy, etc.)</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">3.</span> User enters the 6-digit code to verify: <Code>POST /api/totp/verify</Code></li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">4.</span> Server stores that 2FA is enabled for this user</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">5.</span> On next login, after the password check, the server asks for the TOTP code</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">6.</span> User enters the code from their app — if it matches, login succeeds</li>
          </ol>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here is the setup flow:
          </p>

          <CodeBlock filename="POST /api/totp/setup">
{`// Request (requires Authorization header)
POST /api/totp/setup
Authorization: Bearer <token>

// Response
{
  "data": {
    "secret": "JBSWY3DPEHPK3PXP",
    "qr_url": "otpauth://totp/myapp:john@example.com?secret=JBSWY3DPEHPK3PXP&issuer=myapp",
    "backup_codes": [
      "a1b2c3d4",
      "e5f6g7h8",
      "i9j0k1l2",
      "m3n4o5p6",
      "q7r8s9t0",
      "u1v2w3x4",
      "y5z6a7b8",
      "c9d0e1f2",
      "g3h4i5j6",
      "k7l8m9n0"
    ]
  }
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>qr_url</Code> follows the{' '}
            <a href="https://github.com/google/google-authenticator/wiki/Key-Uri-Format" target="_blank" rel="noreferrer" className="text-primary hover:underline">otpauth:// URI format</a>
            {' '}that all authenticator apps understand. The <Code>secret</Code> is the raw key — you
            can enter it manually if you cannot scan the QR code.
          </p>

          <Challenge number={11} title="Set Up TOTP">
            <p>Call the TOTP setup endpoint from the API docs (you need to be logged in). Copy the
            {' '}<Code>qr_url</Code> from the response. If you have Google Authenticator or Authy on your
            phone, scan the QR code. Otherwise, use an{' '}
            <a href="https://totp.danhersam.com/" target="_blank" rel="noreferrer" className="text-primary hover:underline">online TOTP generator</a>
            {' '}to test with the secret key.</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Backup Codes ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Backup Codes</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            What if you lose your phone? Or your authenticator app gets deleted? You would be locked out
            of your account forever. That is where <strong className="text-foreground">backup codes</strong> come in.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit generates <strong className="text-foreground">10 backup codes</strong> during TOTP setup.
            Each code is a one-time-use alternative to the TOTP code. They are hashed with bcrypt
            before being stored in the database — just like passwords, even if someone steals the
            database, they cannot read the backup codes.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            To use a backup code, send it to the verify endpoint instead of a TOTP code:
          </p>

          <CodeBlock filename="POST /api/totp/verify">
{`{
  "code": "a1b2c3d4"
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The server checks: is this a valid TOTP code? If not, is it a valid backup code? If it
            matches a backup code, the code is marked as used and can never be used again.
          </p>

          <Note>
            Store your backup codes somewhere safe — a password manager, a printed piece of paper in
            a secure location, or an encrypted note. If you lose both your phone AND your backup
            codes, the only way back into your account is through a database admin manually disabling 2FA.
          </Note>

          <Challenge number={12} title="Test Backup Codes">
            <p>During TOTP setup, save your 10 backup codes. Use one of them to verify instead of
            a TOTP code — does it work? Now try using the <strong className="text-foreground">same</strong> backup
            code again. It should fail because each code is one-time-use. How many backup codes do
            you have left?</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Trusted Devices ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Trusted Devices</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Entering a TOTP code every single time you log in from your own computer gets annoying fast.
            Grit solves this with <strong className="text-foreground">trusted devices</strong>.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When you verify a TOTP code, you can set a {'"'}trust this device{'"'} flag. The server creates
            a secure cookie that lasts <strong className="text-foreground">30 days</strong>. During those 30 days,
            that browser skips the TOTP step entirely — password only. After 30 days, or if you clear
            your cookies, you will need to enter a TOTP code again.
          </p>

          <Definition term="Cookie">
            A small piece of data stored in the browser by the server. The browser automatically sends
            it with every request to that domain. Grit uses a <strong className="text-foreground">secure, HTTP-only</strong> cookie
            for trusted devices — {'"'}secure{'"'} means it is only sent over HTTPS,
            and {'"'}HTTP-only{'"'} means JavaScript cannot read it (protecting against XSS attacks).
          </Definition>

          <Challenge number={13} title="Check Trusted Device Cookie">
            <p>After enabling 2FA and verifying with a TOTP code (with {'"'}trust this device{'"'}
            enabled), open your browser DevTools (F12) and go to the Application tab (Chrome) or
            Storage tab (Firefox). Look at Cookies for your domain. Can you find the trusted device
            cookie? What is its expiry date? Is it marked as HttpOnly?</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: OAuth2 Social Login ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">OAuth2 Social Login</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Not everyone wants to create yet another password. OAuth2 lets users sign in with their
            existing <strong className="text-foreground">Google</strong> or <strong className="text-foreground">GitHub</strong> account
            instead.
          </p>

          <Definition term="OAuth2">
            An authorization protocol that lets users grant your app limited access to their account on
            another service (Google, GitHub) without sharing their password. The user clicks {'"'}Sign in
            with Google{'"'}, Google confirms their identity, and sends your app a token with the
            user{"'"}s basic info (name and email). Read more in the{' '}
            <a href="https://oauth.net/2/" target="_blank" rel="noreferrer" className="text-primary hover:underline">OAuth 2.0 specification</a>.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit uses the <a href="https://github.com/markbates/goth" target="_blank" rel="noreferrer" className="text-primary hover:underline">goth</a> library
            for OAuth2. You configure it by adding client credentials to your <Code>.env</Code> file:
          </p>

          <CodeBlock filename=".env">
{`GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here is the complete OAuth2 flow:
          </p>

          <ol className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">1.</span> User clicks {'"'}Sign in with Google{'"'} on your login page</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">2.</span> Browser redirects to Google{"'"}s login page</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">3.</span> User logs in with their Google account</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">4.</span> Google redirects back to your app with a temporary code</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">5.</span> Your API exchanges the code for the user{"'"}s email and name</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">6.</span> If the email already exists in your database, the accounts are linked</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">7.</span> If it is a new email, a new account is created automatically</li>
          </ol>

          <Tip>
            To set up Google OAuth2, go to the{' '}
            <a href="https://console.cloud.google.com/apis/credentials" target="_blank" rel="noreferrer" className="text-primary hover:underline">Google Cloud Console</a>
            {' '}and create OAuth2 credentials. For GitHub, go to{' '}
            <a href="https://github.com/settings/developers" target="_blank" rel="noreferrer" className="text-primary hover:underline">GitHub Developer Settings</a>
            {' '}and create a new OAuth App. Both services are free for development.
          </Tip>

          <Challenge number={14} title="Explore OAuth2 Config">
            <p>Open the <Code>.env</Code> file in your project root. Find the <Code>GOOGLE_CLIENT_ID</Code> and
            {' '}<Code>GITHUB_CLIENT_ID</Code> variables. They are empty by default. Look at the login
            page of your web app — do you see the Google and GitHub sign-in buttons? (They will not
            work without client IDs configured, but the UI should already be there.)</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: The JWT_SECRET ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The JWT_SECRET</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>JWT_SECRET</Code> is arguably the most important configuration value in your app.
            It is the key used to <strong className="text-foreground">sign and verify</strong> every JWT token.
            Anyone who has this secret can create valid tokens for any user — including admin tokens.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This is why:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><span className="text-primary">•</span> It must be <strong className="text-foreground">long and random</strong> (at least 32 characters)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> It must <strong className="text-foreground">never</strong> be committed to Git (it is in <Code>.env</Code>, which is in <Code>.gitignore</Code>)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> It must be <strong className="text-foreground">different</strong> in development and production</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here is how to generate a cryptographically strong secret:
          </p>

          <CodeBlock filename="Terminal">
{`openssl rand -hex 32
# Output: a3f8b2c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1`}
          </CodeBlock>

          <Note>
            If you change the <Code>JWT_SECRET</Code>, <strong className="text-foreground">all existing tokens
            become invalid</strong> and every logged-in user will be forced to log in again. This is
            actually useful if you suspect tokens have been compromised — change the secret to
            immediately invalidate all sessions.
          </Note>

          <Challenge number={15} title="Secure Your JWT_SECRET">
            <p>Open your <Code>.env</Code> file. Look at the <Code>JWT_SECRET</Code> value. Is it a strong
            random string, or is it a weak default like {'"'}secret{'"'} or {'"'}changeme{'"'}? Generate a
            proper secret with <Code>openssl rand -hex 32</Code> and replace it. Restart the API server.
            What happens to your currently logged-in session? (You should be logged out.)</p>
          </Challenge>
        </section>

        {/* ═══ Section 11: Auth in the Frontend ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Auth in the Frontend</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The backend handles token creation and verification, but the frontend is responsible for
            {' '}<strong className="text-foreground">storing tokens</strong>, <strong className="text-foreground">attaching them to requests</strong>,
            and <strong className="text-foreground">managing auth state</strong>. Here is how Grit{"'"}s web app handles it:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">On login:</strong> stores the <Code>access_token</Code> and <Code>refresh_token</Code></li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">On every API call:</strong> attaches the <Code>Authorization: Bearer &lt;token&gt;</Code> header automatically</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">When access token expires:</strong> automatically calls <Code>/api/auth/refresh</Code> to get a new one</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">On logout:</strong> clears both tokens and redirects to the login page</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Auth context:</strong> a React context that provides user data and auth state to all components</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Protected pages use the auth hook to redirect unauthenticated users:
          </p>

          <CodeBlock filename="Protected Page Pattern">
{`// If not logged in, redirect to /login
const { user, isLoading } = useAuth()
if (!isLoading && !user) redirect("/login")`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This pattern runs on every protected page. While <Code>isLoading</Code> is true (the auth
            state is being checked), the page shows a loading spinner. Once loading completes, if
            there is no user, the visitor is immediately redirected to the login page. They never see
            the protected content.
          </p>

          <Tip>
            The token refresh happens transparently. When an API call returns <Code>401 Unauthorized</Code>,
            the API client automatically tries to refresh the token. If the refresh succeeds, it retries
            the original request. If the refresh fails (the refresh token is also expired), the user is
            logged out. This means users stay logged in for up to 7 days without any action.
          </Tip>

          <Challenge number={16} title="Test Frontend Auth Flow">
            <p>Do two things: (1) Log out of your web app. Try to visit <Code>/dashboard</Code> directly
            by typing the URL in the address bar. You should be redirected to the login page. (2) Open
            browser DevTools (F12) and go to the Network tab. Log in and watch the requests. Can you see
            the <Code>/api/auth/login</Code> request? Look at the response — do you see both tokens?</p>
          </Challenge>
        </section>

        {/* ═══ Section 12: Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-6">
            You now understand Grit{"'"}s complete auth system, from password hashing to social login.
            Here is everything you learned:
          </p>

          <ul className="space-y-3 text-muted-foreground mb-8">
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Authentication vs authorization</strong> — who you are vs what you can do</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">JWT flow</strong> — register, login, access token (15 min), refresh token (7 days)</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">bcrypt password hashing</strong> — one-way hashing, never store plain text</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Middleware-based route protection</strong> — Auth middleware checks tokens, RequireRole checks permissions</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">3 built-in roles</strong> — ADMIN, EDITOR, USER — plus custom roles with <Code>grit add role</Code></span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">TOTP two-factor authentication</strong> — 6-digit codes from authenticator apps, QR code setup</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Backup codes</strong> — 10 one-time-use recovery codes, bcrypt-hashed in the database</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Trusted devices</strong> — 30-day secure cookies to skip TOTP on known browsers</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">OAuth2 social login</strong> — Google and GitHub sign-in via the goth library</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">JWT_SECRET security</strong> — random, never committed, change to invalidate all sessions</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Frontend auth flow</strong> — token storage, auto-refresh, protected routes, auth context</span>
            </li>
          </ul>

          <Challenge number={17} title="Final Challenge: Complete Auth Setup">
            <p>Put everything together in one session:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li><strong className="text-foreground">Register a user</strong> through the API or web app</li>
              <li><strong className="text-foreground">Add a MODERATOR role</strong> with <Code>grit add role MODERATOR</Code></li>
              <li><strong className="text-foreground">Enable 2FA</strong> on your admin account using the TOTP setup endpoint</li>
              <li><strong className="text-foreground">Test backup codes</strong> — use one, verify it works, try reusing it</li>
              <li><strong className="text-foreground">Change the JWT_SECRET</strong> in your <Code>.env</Code> file and restart the API — observe what happens to your session</li>
              <li><strong className="text-foreground">Generate a role-restricted resource</strong> with <Code>--roles {'"'}ADMIN,MODERATOR{'"'}</Code> and verify a USER cannot access it</li>
            </ol>
          </Challenge>
        </section>

        {/* ═══ Footer ═══ */}
        <CourseFooter />

        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses/grit-web/code-generator', label: 'Code Generator Mastery' }}
            next={{ href: '/courses/grit-web/admin-panel', label: 'Admin Panel' }}
          />
        </div>
      </main>
    </div>
  )
}
