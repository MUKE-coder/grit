import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/oauth')

export default function OAuthPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend (Go API)</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Social Login (OAuth2)
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit includes Google and GitHub OAuth2 authentication out of the box via{' '}
                <a href="https://github.com/markbates/goth" target="_blank" rel="noreferrer" className="text-primary hover:underline">
                  Gothic
                </a>. Social login buttons appear on all login and register pages across
                every admin style variant. This guide walks you through setting up both providers.
              </p>
            </div>

            <div className="prose-grit">

              {/* ── How It Works ─────────────────────────── */}
              <h2 id="how-it-works">How It Works</h2>
              <p>
                When a user clicks &ldquo;Sign in with Google&rdquo; or &ldquo;Sign in with GitHub&rdquo;,
                they are redirected to the provider&apos;s consent screen. After granting access, the provider
                redirects back to your API, which creates or links the user account and issues JWT tokens.
              </p>

              <CodeBlock filename="OAuth2 flow" code={`User clicks "Sign in with Google"
    |
    v
Browser navigates to /api/auth/oauth/google
    |
    v
Go API (Gothic) redirects to Google consent screen
    |
    v
User grants access on Google
    |
    v
Google redirects to /api/auth/oauth/google/callback
    |
    v
Go API receives profile (name, email, avatar)
  - Finds existing user by email → links GoogleID
  - OR creates new user (provider="google", no password)
  - Generates JWT access + refresh tokens
    |
    v
Redirects to OAUTH_FRONTEND_URL/auth/callback?access_token=...&refresh_token=...
    |
    v
Frontend callback page stores tokens in cookies
Fetches /api/auth/me → redirects to dashboard`} />

              <h2 id="routes">OAuth Routes</h2>
              <p>These routes are automatically registered in your API:</p>

              <CodeBlock filename="OAuth endpoints" code={`GET /api/auth/oauth/google            → Start Google login
GET /api/auth/oauth/google/callback   → Google callback (handles redirect)

GET /api/auth/oauth/github            → Start GitHub login
GET /api/auth/oauth/github/callback   → GitHub callback (handles redirect)`} />

              {/* ── Google Setup ─────────────────────────── */}
              <h2 id="google-setup">Step 1 &mdash; Google OAuth Setup</h2>
              <p>
                Create OAuth2 credentials in the Google Cloud Console to allow users to sign in with their Google account.
              </p>
            </div>

            {/* Step 1.1 */}
            <div className="relative pl-10 mb-8 mt-6">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                1
              </div>
              <div>
                <h3 className="text-lg font-semibold mb-2">Create a Google Cloud Project</h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Go to the{' '}
                  <a href="https://console.cloud.google.com/" target="_blank" rel="noreferrer" className="text-primary hover:underline">
                    Google Cloud Console
                  </a>{' '}
                  and create a new project (or select an existing one).
                </p>
              </div>
            </div>

            {/* Step 1.2 */}
            <div className="relative pl-10 mb-8">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                2
              </div>
              <div>
                <h3 className="text-lg font-semibold mb-2">Configure the OAuth Consent Screen</h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Navigate to <strong>APIs &amp; Services &rarr; OAuth consent screen</strong>:
                </p>
                <ul className="space-y-1.5 text-[13px] text-muted-foreground/70">
                  <li className="flex items-start gap-2"><span className="text-primary mt-0.5">&#8226;</span> Choose <strong>External</strong> user type</li>
                  <li className="flex items-start gap-2"><span className="text-primary mt-0.5">&#8226;</span> Fill in App name (e.g. &ldquo;My Grit App&rdquo;)</li>
                  <li className="flex items-start gap-2"><span className="text-primary mt-0.5">&#8226;</span> Add your support email</li>
                  <li className="flex items-start gap-2"><span className="text-primary mt-0.5">&#8226;</span> Under Scopes, add <code>email</code>, <code>profile</code>, and <code>openid</code></li>
                  <li className="flex items-start gap-2"><span className="text-primary mt-0.5">&#8226;</span> Add test users if the app is in &ldquo;Testing&rdquo; mode</li>
                </ul>
              </div>
            </div>

            {/* Step 1.3 */}
            <div className="relative pl-10 mb-8">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                3
              </div>
              <div>
                <h3 className="text-lg font-semibold mb-2">Create OAuth2 Credentials</h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Go to <strong>APIs &amp; Services &rarr; Credentials &rarr; Create Credentials &rarr; OAuth client ID</strong>:
                </p>
                <ul className="space-y-1.5 text-[13px] text-muted-foreground/70 mb-4">
                  <li className="flex items-start gap-2"><span className="text-primary mt-0.5">&#8226;</span> Application type: <strong>Web application</strong></li>
                  <li className="flex items-start gap-2"><span className="text-primary mt-0.5">&#8226;</span> Name: anything (e.g. &ldquo;Grit Web Client&rdquo;)</li>
                  <li className="flex items-start gap-2"><span className="text-primary mt-0.5">&#8226;</span> Authorized redirect URIs &mdash; add both:</li>
                </ul>
                <CodeBlock filename="Authorized redirect URIs" code={`# Development
http://localhost:8080/api/auth/oauth/google/callback

# Production (replace with your domain)
https://api.yourdomain.com/api/auth/oauth/google/callback`} />
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mt-3">
                  After creating, copy the <strong>Client ID</strong> and <strong>Client Secret</strong>.
                </p>
              </div>
            </div>

            {/* Step 1.4 */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                4
              </div>
              <div>
                <h3 className="text-lg font-semibold mb-2">Add to .env</h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Paste the credentials into your project&apos;s <code>.env</code> file:
                </p>
                <CodeBlock filename=".env" code={`GOOGLE_CLIENT_ID=123456789-abcdefg.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-your-secret-here`} />
              </div>
            </div>

            <div className="prose-grit">
              {/* ── GitHub Setup ─────────────────────────── */}
              <h2 id="github-setup">Step 2 &mdash; GitHub OAuth Setup</h2>
              <p>
                Create an OAuth App on GitHub to allow users to sign in with their GitHub account.
              </p>
            </div>

            {/* Step 2.1 */}
            <div className="relative pl-10 mb-8 mt-6">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                1
              </div>
              <div>
                <h3 className="text-lg font-semibold mb-2">Create a GitHub OAuth App</h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Go to{' '}
                  <a href="https://github.com/settings/developers" target="_blank" rel="noreferrer" className="text-primary hover:underline">
                    GitHub &rarr; Settings &rarr; Developer settings &rarr; OAuth Apps
                  </a>{' '}
                  and click <strong>New OAuth App</strong>.
                </p>
              </div>
            </div>

            {/* Step 2.2 */}
            <div className="relative pl-10 mb-8">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                2
              </div>
              <div>
                <h3 className="text-lg font-semibold mb-2">Fill in the App Details</h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Fill in the following fields:
                </p>
                <div className="rounded-lg border border-border/40 bg-card/50 overflow-hidden mb-4">
                  <table className="w-full text-[13px]">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2 font-medium text-foreground/80">Field</th>
                        <th className="text-left px-4 py-2 font-medium text-foreground/80">Value</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground/70">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2 font-medium">Application name</td>
                        <td className="px-4 py-2">My Grit App</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2 font-medium">Homepage URL</td>
                        <td className="px-4 py-2"><code>http://localhost:3001</code></td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2 font-medium">Authorization callback URL</td>
                        <td className="px-4 py-2"><code>http://localhost:8080/api/auth/oauth/github/callback</code></td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  For production, replace the callback URL with your live API domain:
                </p>
                <CodeBlock filename="Production callback URL" code={`https://api.yourdomain.com/api/auth/oauth/github/callback`} />
              </div>
            </div>

            {/* Step 2.3 */}
            <div className="relative pl-10 mb-8">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                3
              </div>
              <div>
                <h3 className="text-lg font-semibold mb-2">Generate a Client Secret</h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  After creating the app, click <strong>Generate a new client secret</strong>. Copy both the
                  Client ID and Client Secret immediately &mdash; the secret is only shown once.
                </p>
              </div>
            </div>

            {/* Step 2.4 */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                4
              </div>
              <div>
                <h3 className="text-lg font-semibold mb-2">Add to .env</h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Add the GitHub credentials to your <code>.env</code>:
                </p>
                <CodeBlock filename=".env" code={`GITHUB_CLIENT_ID=Ov23li...
GITHUB_CLIENT_SECRET=abc123...`} />
              </div>
            </div>

            <div className="prose-grit">
              {/* ── Frontend URL ─────────────────────────── */}
              <h2 id="frontend-url">Step 3 &mdash; Set the Frontend Callback URL</h2>
              <p>
                After a successful OAuth login, the API redirects the user to your admin frontend
                with JWT tokens in the URL. Set the <code>OAUTH_FRONTEND_URL</code> in your <code>.env</code>:
              </p>
            </div>

            <div className="mt-4 mb-12">
              <CodeBlock filename=".env" code={`# Where to redirect after OAuth callback
# Development — your admin panel URL
OAUTH_FRONTEND_URL=http://localhost:3001

# Production
OAUTH_FRONTEND_URL=https://admin.yourdomain.com`} />
            </div>

            <div className="prose-grit">
              {/* ── Full .env ────────────────────────────── */}
              <h2 id="full-env">Complete .env Example</h2>
              <p>
                Here is the full OAuth section for your <code>.env</code> file:
              </p>
            </div>

            <div className="mt-4 mb-12">
              <CodeBlock filename=".env" code={`# OAuth2 — Social Login
# Get Google credentials: https://console.cloud.google.com/apis/credentials
GOOGLE_CLIENT_ID=123456789-abcdefg.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-your-secret-here

# Get GitHub credentials: https://github.com/settings/developers
GITHUB_CLIENT_ID=Ov23li...
GITHUB_CLIENT_SECRET=abc123...

# Where to redirect after successful OAuth login
OAUTH_FRONTEND_URL=http://localhost:3001`} />
            </div>

            <div className="prose-grit">
              {/* ── Account Linking ──────────────────────── */}
              <h2 id="account-linking">Account Linking</h2>
              <p>
                Grit automatically handles account linking by email:
              </p>
              <ul>
                <li>
                  <strong>Existing user with same email:</strong> The OAuth provider ID (GoogleID or GithubID)
                  is linked to the existing account. The user can continue logging in with either email/password
                  or the social provider.
                </li>
                <li>
                  <strong>New user:</strong> A new account is created with the provider set to <code>&quot;google&quot;</code> or{' '}
                  <code>&quot;github&quot;</code> and an empty password. The user can set a password later via the profile page.
                </li>
                <li>
                  <strong>Password-less login attempt:</strong> If an OAuth-only user tries to log in with
                  email/password, they get a helpful error: <em>&quot;This account uses social login. Please sign in
                  with Google/GitHub.&quot;</em>
                </li>
              </ul>

              {/* ── User Model ───────────────────────────── */}
              <h2 id="user-model">User Model Fields</h2>
              <p>
                The following fields are added to the User model for OAuth support:
              </p>
            </div>

            <div className="mt-4 mb-12">
              <CodeBlock filename="apps/api/internal/models/user.go" language="go" code={`type User struct {
    // ... existing fields ...

    Provider string \`gorm:"size:50;default:'local'" json:"provider"\`
    GoogleID string \`gorm:"size:255" json:"-"\`
    GithubID string \`gorm:"size:255" json:"-"\`
}`} />
              <div className="mt-3 rounded-lg border border-border/40 bg-card/50 overflow-hidden">
                <table className="w-full text-[13px]">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2 font-medium text-foreground/80">Field</th>
                      <th className="text-left px-4 py-2 font-medium text-foreground/80">Description</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground/70">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2 font-mono text-xs">Provider</td>
                      <td className="px-4 py-2"><code>&quot;local&quot;</code>, <code>&quot;google&quot;</code>, or <code>&quot;github&quot;</code></td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2 font-mono text-xs">GoogleID</td>
                      <td className="px-4 py-2">Google account ID (hidden from API responses)</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2 font-mono text-xs">GithubID</td>
                      <td className="px-4 py-2">GitHub account ID (hidden from API responses)</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            <div className="prose-grit">
              {/* ── Conditional Providers ─────────────────── */}
              <h2 id="conditional-providers">Conditional Providers</h2>
              <p>
                Providers are initialized conditionally in <code>main.go</code>. If a provider&apos;s
                credentials are empty, it is simply skipped &mdash; the app starts normally without it.
                This means you can:
              </p>
              <ul>
                <li>Enable only Google (leave GitHub credentials empty)</li>
                <li>Enable only GitHub (leave Google credentials empty)</li>
                <li>Enable both by filling in all credentials</li>
                <li>Disable both by leaving all credentials empty (social buttons still show but won&apos;t work)</li>
              </ul>

              {/* ── Frontend ─────────────────────────────── */}
              <h2 id="frontend-callback">Frontend Callback Page</h2>
              <p>
                Grit scaffolds a callback page at <code>apps/admin/app/(auth)/callback/page.tsx</code>
                that handles the OAuth redirect. It:
              </p>
              <ol>
                <li>Extracts <code>access_token</code> and <code>refresh_token</code> from the URL</li>
                <li>Stores them in cookies via <code>js-cookie</code></li>
                <li>Fetches <code>/api/auth/me</code> to get the user&apos;s role</li>
                <li>Redirects to <code>/dashboard</code> (ADMIN/EDITOR) or <code>/profile</code> (USER)</li>
              </ol>
              <p>
                The social login buttons on the login and register pages are simple <code>&lt;a&gt;</code> tags
                that navigate to the API OAuth endpoints. No JavaScript is needed &mdash; it&apos;s a full-page
                redirect flow.
              </p>

              {/* ── Production Checklist ──────────────────── */}
              <h2 id="production">Production Checklist</h2>
            </div>

            <div className="mt-4 mb-12">
              <div className="space-y-2.5">
                {[
                  'Set APP_URL in .env to your production API URL (used for callback URLs)',
                  'Update Google redirect URI to your production callback URL',
                  'Update GitHub callback URL to your production callback URL',
                  'Set OAUTH_FRONTEND_URL to your production admin panel URL',
                  'Publish your Google OAuth consent screen (move from "Testing" to "Production")',
                  'Ensure your production domain uses HTTPS',
                  'Test both login and registration flows with each provider',
                ].map((item) => (
                  <div
                    key={item}
                    className="flex items-start gap-2.5 text-[13px] text-muted-foreground/70"
                  >
                    <span className="text-primary mt-0.5">&#9744;</span>
                    {item}
                  </div>
                ))}
              </div>
            </div>

            <div className="prose-grit">
              {/* ── Troubleshooting ───────────────────────── */}
              <h2 id="troubleshooting">Troubleshooting</h2>
            </div>

            <div className="mt-4 mb-12">
              <div className="space-y-4">
                {[
                  {
                    q: '"redirect_uri_mismatch" error from Google',
                    a: 'The callback URL in your Google Cloud Console must exactly match your API URL + /api/auth/oauth/google/callback. Check for trailing slashes, http vs https, and port numbers.',
                  },
                  {
                    q: 'OAuth works locally but not in production',
                    a: 'Make sure APP_URL is set to your production API URL (e.g. https://api.yourdomain.com). The callback URLs are built from APP_URL + the OAuth path.',
                  },
                  {
                    q: 'User created with empty name',
                    a: 'Some GitHub accounts have no display name set. The OAuth handler falls back to the email username. Users can update their name via the profile page.',
                  },
                  {
                    q: '"This account uses social login" error',
                    a: 'The user signed up via Google/GitHub and has no password set. They need to use the same social provider to log in, or set a password in their profile first.',
                  },
                  {
                    q: 'Callback page shows error or blank screen',
                    a: 'Check that OAUTH_FRONTEND_URL matches your admin panel URL exactly (including port). Verify the callback page exists at apps/admin/app/(auth)/callback/page.tsx.',
                  },
                ].map((item) => (
                  <div
                    key={item.q}
                    className="rounded-lg border border-border/40 bg-card/50 p-4"
                  >
                    <div className="text-[13px] font-medium text-foreground/90 mb-1">
                      {item.q}
                    </div>
                    <div className="text-[12px] text-muted-foreground/60 leading-relaxed">
                      {item.a}
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Bottom Navigation */}
            <div className="flex flex-wrap gap-3 mt-12 pt-6 border-t border-border/30">
              <Button variant="outline" asChild className="border-border/60 bg-transparent hover:bg-accent/50">
                <Link href="/docs/backend/authentication">
                  <ArrowLeft className="mr-1.5 h-4 w-4" />
                  Authentication
                </Link>
              </Button>
              <Button asChild className="glow-purple-sm ml-auto">
                <Link href="/docs/backend/response-format">
                  API Response Format
                  <ArrowRight className="ml-1.5 h-4 w-4" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
