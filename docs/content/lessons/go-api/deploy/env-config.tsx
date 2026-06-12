import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Last lesson — the pre-launch env vars checklist. Miss one of these
        and prod either won&apos;t boot or will boot in an insecure state.
        Five minutes to read, save hours of debugging.
      </p>

      <h2>The must-set list</h2>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Var</th>
              <th className="text-left px-3 py-2 font-medium">Required</th>
              <th className="text-left px-3 py-2 font-medium">If missing / wrong</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr><td className="px-3 py-2 font-mono text-[12px]">APP_ENV</td><td>Must be <code>production</code></td><td className="text-[12px]">Gin in debug mode; verbose error responses; secrets leaked</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">DATABASE_URL</td><td>Yes</td><td className="text-[12px]">Boot fails at config load</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">JWT_SECRET</td><td>≥32 chars</td><td className="text-[12px]">Boot rejected; existing tokens become invalid</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">SENTINEL_PASSWORD</td><td>Random</td><td className="text-[12px]">Boot rejected with default password</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">SENTINEL_SECRET_KEY</td><td>≥32 chars</td><td className="text-[12px]">Same as above</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">PULSE_PASSWORD</td><td>Random</td><td className="text-[12px]">Boot rejected with default</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">CORS_ORIGINS</td><td>Your real domains</td><td className="text-[12px]">Frontend can&apos;t call API; CORS errors</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">APP_URL</td><td>https://your-domain</td><td className="text-[12px]">OAuth redirects break; absolute URLs in emails are wrong</td></tr>
          </tbody>
        </table>
      </div>

      <h2>The should-set list</h2>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Var</th>
              <th className="text-left px-3 py-2 font-medium">Why</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr><td className="px-3 py-2 font-mono text-[12px]">RESEND_API_KEY</td><td className="text-[12px]">Real email sends; without it, mail is silently dropped</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">STORAGE_DRIVER + creds</td><td className="text-[12px]">File uploads work and survive container restart</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">AI_GATEWAY_API_KEY</td><td className="text-[12px]">AI features stop working if you removed dev defaults</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">OAUTH_FRONTEND_URL</td><td className="text-[12px]">OAuth post-login redirect points to your real frontend</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">GORM_STUDIO_ENABLED=false</td><td className="text-[12px]">Don&apos;t expose the DB browser in production</td></tr>
          </tbody>
        </table>
      </div>

      <h2>Generating strong values</h2>
      <CodeBlock
        terminal
        code={`# JWT_SECRET + SENTINEL_SECRET_KEY — 32 bytes hex
openssl rand -hex 32

# SENTINEL_PASSWORD + PULSE_PASSWORD — 16 bytes hex (short but strong)
openssl rand -hex 16

# Strong DB password — 24+ chars, no shell-special characters
openssl rand -base64 24 | tr -d '/+='`}
      />
      <p>
        v3.25.1+ of Grit auto-generates these at scaffold time. If you
        scaffolded older, regenerate before going to prod.
      </p>

      <h2>OAuth — production callback URLs</h2>
      <p>
        Each OAuth provider needs your production callback URL added to
        its allow list:
      </p>
      <ul>
        <li>
          Google: <code>https://api.acme.com/api/auth/google/callback</code>
        </li>
        <li>
          GitHub: <code>https://api.acme.com/api/auth/github/callback</code>
        </li>
      </ul>
      <p>
        Both providers let you list multiple URLs — keep your localhost dev
        URLs too so dev still works.
      </p>

      <TipBox tone="danger">
        <strong>Don&apos;t reuse JWT_SECRET across environments.</strong> If
        staging&apos;s secret leaks, your prod tokens are exposed. Different
        random values per environment. Same for SENTINEL_SECRET_KEY.
      </TipBox>

      <h2>Database backups</h2>
      <p>
        Not strictly an env var, but: <strong>set up daily backups</strong>{' '}
        before going to prod. Five-minute setup, prevents catastrophic
        data loss. Options:
      </p>
      <ul>
        <li>
          <strong>pg_dump + cron</strong> — simplest. Dumps DB nightly, ships
          to S3 / R2.
        </li>
        <li>
          <strong>Backrest / pgBackRest</strong> — production-grade PITR
          (point-in-time recovery).
        </li>
        <li>
          <strong>Managed Postgres</strong> (Neon, Supabase) — backups are
          automatic. Skip the self-host.
        </li>
      </ul>

      <h2>The pre-flight checklist</h2>
      <p>Before <code>grit deploy --domain ...</code> in production:</p>
      <ul>
        <li>☐ <code>.env.production</code> exists and has every required var</li>
        <li>☐ All passwords are random (not <code>change-me</code>)</li>
        <li>☐ <code>APP_ENV=production</code></li>
        <li>☐ <code>CORS_ORIGINS</code> only lists real frontend domains</li>
        <li>☐ <code>GORM_STUDIO_ENABLED=false</code></li>
        <li>☐ OAuth callback URLs are added to provider allow lists</li>
        <li>☐ DNS A record points to VPS IP</li>
        <li>☐ DB backup plan exists</li>
        <li>☐ You can SSH to the VPS</li>
      </ul>

      <KnowledgeCheck
        question="You ran `grit deploy` to production. The Sentinel dashboard is reachable at /sentinel/ui with credentials `admin/sentinel`. What's the immediate priority?"
        choices={[
          {
            label: 'Nothing — Sentinel\'s defaults are fine for production',
            feedback:
              "Wrong — `sentinel` is the literal default password. Anyone who reads the docs can log into your WAF dashboard.",
          },
          {
            label: 'Rotate SENTINEL_PASSWORD + SENTINEL_SECRET_KEY immediately and redeploy',
            correct: true,
            feedback:
              "Right — change them to random values, redeploy. Better: Grit v3.25.1+ refuses to start in production with default credentials, so the API wouldn't have come up if these were defaults. Check your Grit version.",
          },
          {
            label: 'Disable the Sentinel dashboard',
            feedback:
              "Throwing out the baby with the bath water. The dashboard is useful; just secure it with strong creds.",
          },
          {
            label: 'Hide it behind a VPN',
            feedback:
              "Good defense-in-depth, but doesn't address the root issue: default credentials. Fix those first.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Build the pre-flight checklist for your bench-api&apos;s production
              env. In <code>notes.md</code>:
            </p>
            <ol>
              <li>List every env var your <code>.env.production</code> sets</li>
              <li>For each, write &quot;random&quot;, &quot;copied from dev&quot;, or &quot;from a service dashboard&quot;</li>
              <li>
                Confirm no value is <code>change-me</code>,{' '}
                <code>your-super-secret-jwt-key</code>, or similar template
                text
              </li>
            </ol>
            <p>
              That&apos;s the chapter 6 assignment partially done — the second
              half is the actual deploy + HTTPS verification.
            </p>
          </>
        }
        hint={
          <>
            Diff against <code>.env.example</code>. Anything that looks like
            a placeholder there — JWT_SECRET, passwords — should be a real
            random value in your production env.
          </>
        }
        solution={
          <>
            <p>
              A healthy <code>.env.production</code> looks like (abbreviated):
            </p>
            <CodeBlock
              language="dotenv"
              code={`APP_ENV=production                         ✓
DATABASE_URL=postgres://...                ✓ (random password)
JWT_SECRET=a3f8...64chars                  ✓ (openssl rand -hex 32)
SENTINEL_PASSWORD=9b4d...32chars           ✓ (openssl rand -hex 16)
SENTINEL_SECRET_KEY=ec47...64chars         ✓ (openssl rand -hex 32)
PULSE_PASSWORD=1f72...32chars              ✓ (openssl rand -hex 16)
RESEND_API_KEY=re_live_...                 ✓ (from Resend dashboard)
GORM_STUDIO_ENABLED=false                  ✓
CORS_ORIGINS=https://app.acme.com          ✓ (real frontend only)`}
            />
          </>
        }
      />

      <h2>You finished Building a Go API 🎉</h2>
      <p>
        Six chapters, 19 lessons. You can now scaffold a Go API, model the
        data layer, add auth with JWT + OAuth + 2FA + RBAC, integrate
        background jobs / mail / storage / AI, run a real WAF + observability
        + audit log, and deploy to a VPS with HTTPS.
      </p>
      <p>
        Next: pick the frontend that fits your product —{' '}
        <strong>Mobile (Expo)</strong>, <strong>Web (Next.js)</strong>,{' '}
        <strong>Web (TanStack)</strong>, <strong>Desktop (Wails)</strong>, or{' '}
        <strong>Multi-platform</strong>. All assume the API you just built.
      </p>
    </>
  )
}
