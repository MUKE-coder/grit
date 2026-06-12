import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Time to ship. <code>grit deploy</code> takes your local project and
        runs it on a VPS — Docker, Traefik for HTTPS, your domain, a
        functioning Postgres. One command. This lesson covers what it does
        and the prerequisites.
      </p>

      <h2>What the command does</h2>
      <CodeBlock
        terminal
        code={`grit deploy --host root@1.2.3.4 --domain api.acme.com`}
      />
      <p>It runs (in this order):</p>
      <ol>
        <li>SSH&apos;s to the host, installs Docker + Docker Compose if missing</li>
        <li>Builds your project locally + ships the image to the host</li>
        <li>Copies <code>docker-compose.prod.yml</code> + your <code>.env.production</code></li>
        <li>Provisions Traefik with Let&apos;s Encrypt for HTTPS</li>
        <li>
          Starts Postgres + Redis + your API behind Traefik with the domain
          you specified
        </li>
        <li>Health-checks the deployment</li>
        <li>Prints the URL — typically <code>https://api.acme.com/api/health</code></li>
      </ol>

      <h2>Prerequisites</h2>
      <ul>
        <li>
          A VPS (Hetzner / DigitalOcean / Linode / Vultr — anywhere with
          SSH + a public IP). $4/month tier works.
        </li>
        <li>SSH key access (no password auth)</li>
        <li>A domain you control + DNS access</li>
        <li>An <code>.env.production</code> in your project with real secrets</li>
      </ul>

      <h2>Setting up DNS</h2>
      <p>
        Point an A record at your VPS IP <em>before</em> running deploy.
        Otherwise Let&apos;s Encrypt can&apos;t verify domain ownership and HTTPS
        provisioning fails.
      </p>
      <CodeBlock
        language="text"
        code={`api.acme.com  →  A  →  1.2.3.4   (your VPS IP)`}
      />

      <h2>.env.production — what changes from .env</h2>
      <p>
        Your dev <code>.env</code> uses localhost services + Mailhog +
        MinIO. Production uses real services:
      </p>
      <CodeBlock
        language="dotenv"
        filename=".env.production"
        code={`APP_ENV=production
APP_PORT=8080
APP_URL=https://api.acme.com

DATABASE_URL=postgres://grit:STRONG_PASSWORD@postgres:5432/myapp?sslmode=disable
REDIS_URL=redis://redis:6379

JWT_SECRET=<openssl rand -hex 32>

# Real email
RESEND_API_KEY=re_live_...
MAIL_FROM=noreply@acme.com

# Real storage
STORAGE_DRIVER=r2
R2_ACCOUNT_ID=...
R2_ACCESS_KEY=...
R2_SECRET_KEY=...
R2_BUCKET=acme-prod

# Sentinel + Pulse — strong passwords required in production
SENTINEL_PASSWORD=<openssl rand -hex 16>
SENTINEL_SECRET_KEY=<openssl rand -hex 32>
PULSE_PASSWORD=<openssl rand -hex 16>

# CORS
CORS_ORIGINS=https://app.acme.com,https://acme.com`}
      />

      <TipBox tone="warning">
        <strong>Never commit <code>.env.production</code> to git.</strong>{' '}
        It contains every secret. Add it to <code>.gitignore</code> (Grit&apos;s
        scaffold already does). Store the file somewhere safe (1Password,
        Bitwarden, vault).
      </TipBox>

      <h2>What gets deployed</h2>
      <CodeBlock
        language="text"
        code={`on the VPS:

/opt/myapp/
├── docker-compose.prod.yml   ← Traefik + Postgres + Redis + API
├── .env                       ← your .env.production renamed
├── traefik/                   ← cert storage, config
└── postgres-data/             ← Postgres volume`}
      />
      <p>
        Docker Compose handles all process management. <code>systemd</code>{' '}
        is optional; the compose stack starts on boot via{' '}
        <code>restart: unless-stopped</code>.
      </p>

      <h2>Subsequent deploys</h2>
      <CodeBlock
        terminal
        code={`# After the first deploy, subsequent ones are this:
git push                            # commit your changes
grit deploy --host root@1.2.3.4    # ship them`}
      />
      <p>
        Grit rebuilds the image, ships it, restarts the API container with
        a rolling update. Zero downtime if you wire health checks
        (see next lesson).
      </p>

      <h2>Alternatives</h2>
      <ul>
        <li>
          <strong>Dokploy</strong> — UI for the same docker-compose deploys.
          Worth a look if you want a dashboard.
        </li>
        <li>
          <strong>Fly.io / Railway / Render</strong> — managed PaaS. More
          expensive than a VPS but zero config.
        </li>
        <li>
          <strong>Kubernetes</strong> — for products at scale. Grit ships
          example manifests; for &lt;10K users you don&apos;t need it.
        </li>
      </ul>

      <KnowledgeCheck
        question="You ran `grit deploy` but the API container restart-loops with `DATABASE_URL is required`. What's the most likely cause?"
        choices={[
          {
            label: 'Postgres failed to start',
            feedback:
              "Possible but wouldn't manifest as 'DATABASE_URL is required'. That error is from the config loader — DATABASE_URL is empty.",
          },
          {
            label: 'You forgot to put DATABASE_URL in .env.production',
            correct: true,
            feedback:
              "Right — the deploy script ships your .env.production to the VPS as the API's .env. If DATABASE_URL isn't set there, the config loader rejects boot.",
          },
          {
            label: 'Traefik is misconfigured',
            feedback:
              "Wrong — Traefik handles routing, not env vars. The container would start and Traefik would route to it; only then would Traefik see a 502.",
          },
          {
            label: 'SSL cert provisioning failed',
            feedback:
              "Cert failures cause TLS handshake errors, not DATABASE_URL complaints. The error tells you which env var is missing.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Spin up a $4 VPS (Hetzner / DigitalOcean / etc) and deploy
              your bench-api to it.
            </p>
            <ol>
              <li>Create the VPS, note its IP.</li>
              <li>
                Point an A record at it. Wait ~5 minutes for DNS to
                propagate.
              </li>
              <li>
                Create <code>.env.production</code> with your real secrets.
              </li>
              <li>Run <code>grit deploy --host root@&lt;IP&gt; --domain api.&lt;your-domain&gt;</code></li>
              <li>
                Hit <code>https://api.&lt;your-domain&gt;/api/health</code>{' '}
                from your phone (not localhost!) to confirm.
              </li>
            </ol>
            <p>Paste the response (with the real domain) in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            Most failures are DNS propagation — check with{' '}
            <code>dig api.&lt;your-domain&gt;</code>. If it doesn&apos;t return
            your VPS IP, wait 10 more minutes and retry.
          </>
        }
        solution={
          <>
            <p>You should see the same envelope you saw locally — now over HTTPS:</p>
            <CodeBlock
              language="json"
              code={`{ "status": "ok", "version": "0.1.0" }`}
            />
            <p>
              With a green padlock in the address bar. You have a real
              production API. That&apos;s chapter 6&apos;s assignment done.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Last lesson — the <strong>env vars checklist</strong>. What MUST be
        set in production before you ship to real users.
      </p>
    </>
  )
}
