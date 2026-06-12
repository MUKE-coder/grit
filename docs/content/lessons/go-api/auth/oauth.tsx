import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        OAuth2 means &quot;sign in with Google / GitHub&quot; — the user clicks a
        button, gets redirected to the provider, comes back authenticated.
        Grit ships handlers for both. This lesson covers the wiring; the
        provider setup is a one-time copy-paste.
      </p>

      <h2>The flow</h2>
      <CodeBlock
        language="text"
        code={`User → Click "Sign in with Google"
       → GET /api/auth/google
           → Grit redirects to accounts.google.com
       → Google authenticates user
       → Google redirects to /api/auth/google/callback?code=...
           → Grit exchanges code for user info
           → Grit creates or finds the user
           → Grit issues JWT (access + refresh tokens)
           → Frontend stores tokens, user is signed in`}
      />

      <h2>The 3-line provider setup</h2>
      <p>
        Get OAuth credentials from the provider once. For Google: create a
        project in <a href="https://console.cloud.google.com/apis/credentials" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Google Cloud Console</a>,
        configure the consent screen, get client ID + secret. For GitHub:
        <a href="https://github.com/settings/developers" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline"> Developer Settings → OAuth Apps</a> → New OAuth App.
      </p>
      <p>
        Then add to <code>.env</code>:
      </p>
      <CodeBlock
        language="dotenv"
        filename=".env"
        code={`# Google OAuth
GOOGLE_CLIENT_ID=1234-...apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-...

# GitHub OAuth
GITHUB_CLIENT_ID=Iv1....
GITHUB_CLIENT_SECRET=...

# Where to redirect after success (your web frontend)
OAUTH_FRONTEND_URL=http://localhost:3001`}
      />
      <p>
        That&apos;s it. Restart the API. Both providers are wired.
      </p>

      <h2>The callback URLs you give the provider</h2>
      <ul>
        <li>Google: <code>http://localhost:8080/api/auth/google/callback</code></li>
        <li>GitHub: <code>http://localhost:8080/api/auth/github/callback</code></li>
      </ul>
      <p>
        For production, swap <code>localhost:8080</code> for your real domain.
        Each provider lets you list multiple callback URLs — typically you&apos;d
        list both your localhost dev URL and your production URL.
      </p>

      <TipBox tone="info">
        <strong>State parameter:</strong> Grit&apos;s OAuth handler uses a
        CSRF-resistant <code>state</code> token signed with{' '}
        <code>JWT_SECRET</code>. The callback verifies the state matches,
        preventing CSRF attacks during the OAuth handshake.
      </TipBox>

      <h2>Account linking — what happens if email already exists?</h2>
      <p>
        User signs up with <code>alex@example.com</code> + password. Later,
        they hit &quot;Sign in with Google&quot; and Google says &quot;this account is{' '}
        <code>alex@example.com</code>&quot;. Grit&apos;s default: match the
        existing user, link the Google ID, sign them in. The user has both
        a password AND a Google login from then on.
      </p>
      <p>
        If you want strict separation, edit{' '}
        <code>authService.HandleOAuthCallback</code> to reject the merge or
        prompt the user.
      </p>

      <h2>The user&apos;s profile data</h2>
      <p>What you get back from each provider:</p>
      <CodeBlock
        language="text"
        code={`Google:  email, email_verified, name, picture, locale
GitHub:  email (sometimes private!), login (username), name, avatar_url`}
      />
      <p>
        GitHub&apos;s &quot;email&quot; may be private. Grit&apos;s callback handler grabs the
        primary verified email from{' '}
        <code>GET /user/emails</code> if the primary endpoint hides it.
      </p>

      <KnowledgeCheck
        question="A user signs up with email/password, then later clicks 'Sign in with Google'. What's the default behaviour?"
        choices={[
          {
            label: 'Grit rejects the OAuth login because email already exists',
            feedback:
              "Wrong — that'd be a terrible UX. The default is to link both methods to the same account.",
          },
          {
            label: 'Grit creates a duplicate user account',
            feedback:
              "Wrong — that'd give the user two separate identities. Email uniqueness in the DB would also reject it.",
          },
          {
            label: 'Grit matches the email, links the Google ID, signs the user in. They can now use either method.',
            correct: true,
            feedback:
              "Right — account linking by email is the friendly default. You can override in authService.HandleOAuthCallback if you want stricter rules.",
          },
          {
            label: 'Grit issues a verification email asking the user to confirm the link',
            feedback:
              "Sensible but not the default. You can add this in HandleOAuthCallback if your security model requires it.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Wire one OAuth provider (Google or GitHub — your pick) on your
              bench-api. Then trigger the flow:
            </p>
            <ol>
              <li>Open <code>http://localhost:8080/api/auth/google</code> in your browser</li>
              <li>Sign in with your Google account</li>
              <li>You should land on the OAUTH_FRONTEND_URL with tokens in the query string</li>
              <li>
                Use that access token to hit <code>/api/auth/me</code> with
                curl
              </li>
            </ol>
            <p>Paste the /me response in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            If you don&apos;t have a frontend running on OAUTH_FRONTEND_URL,
            point it at <code>http://localhost:8080/api/health</code> for the
            exercise — you&apos;ll see the tokens in the redirect URL bar even
            though the page is just JSON.
          </>
        }
        solution={
          <>
            <p>
              You should see a redirect that looks like:
            </p>
            <CodeBlock
              language="text"
              code={`http://localhost:3001/auth/callback?access_token=eyJhb...&refresh_token=eyJhb...`}
            />
            <p>
              Pluck the access_token and hit <code>/api/auth/me</code> — you
              should see your real Google profile info populated in the user
              record.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Password and OAuth handle the &quot;something you know&quot; factor. Next
        lesson — <strong>TOTP 2FA</strong> for the &quot;something you have&quot;
        factor.
      </p>
    </>
  )
}
