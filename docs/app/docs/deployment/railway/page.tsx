import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/deployment/railway')

const kept: { group: string; what: string }[] = [
  { group: 'Everything with a value', what: 'Your app name, URLs, JWT secret, encryption key, provider keys, module switches.' },
  { group: 'APP_ENV', what: 'Forced to production, whatever the file on your laptop says.' },
  { group: 'DATABASE_URL, REDIS_URL', what: 'With --provision, replaced by Railway references so they follow the database rather than a copy of its password.' },
]

const dropped: { group: string; what: string }[] = [
  { group: 'Empty values', what: 'A generated .env is mostly placeholders for providers you do not use. An empty string configures nothing and makes the dashboard unreadable.' },
  { group: 'PORT, APP_PORT', what: 'Railway assigns the port. Overriding it is how a deploy ends up listening on the wrong one.' },
  { group: 'MINIO_*, MAILHOG_*', what: 'Containers on your machine. On Railway, storage and mail are configured with their own provider keys.' },
  { group: 'DB_HOST, MYSQL_*, POSTGRES_*', what: 'The compose database. On Railway the app connects through DATABASE_URL, and these would contradict it.' },
]

export default function RailwayPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Deployment</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Railway</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                <code>grit deploy --railway</code> links your project, adds Postgres and Redis,
                pushes the variables a deploy actually needs, uploads the API and gives you a URL.
              </p>
            </div>

            <div className="prose-grit">
              <h2 id="once">Once, per machine</h2>
              <p>
                Railway&apos;s API has no endpoint that accepts local source: a service is built
                either from a connected GitHub repository or from an archive their CLI uploads. So
                the CLI is what deploys, and Grit drives it.
              </p>
              <CodeBlock
                terminal
                code={`npm i -g @railway/cli     # or: brew install railway
railway login`}
              />
              <p>
                In CI there is no browser to log in with, so set{' '}
                <code>RAILWAY_TOKEN</code> instead. Grit checks for either before it does anything.
              </p>

              <h2 id="first">The first deploy</h2>
              <p>
                From the root of your project. The first step links a Railway project, and that is
                the one thing Grit will not decide for you: a new project or an existing one is
                your call, not a default.
              </p>
              <CodeBlock
                terminal
                code={`cd apps/api
railway init          # a new project
# or: railway link    # an existing one

cd ../..
grit deploy --railway --provision`}
              />
              <p>
                <code>--provision</code> adds Postgres and Redis and points{' '}
                <code>DATABASE_URL</code> and <code>REDIS_URL</code> at them. Leave it off on every
                deploy after the first, or if you are bringing your own database.
              </p>

              <Callout type="tip" title="See the plan before it runs">
                <code>grit deploy --railway --dry-run</code> prints every command it would run, in
                order, and runs none of them. Variable values are masked, so the output is safe to
                paste into an issue.
              </Callout>

              <h2 id="variables">Which variables are sent</h2>
              <p>
                A generated <code>.env</code> has around a hundred entries. Sending all of them is
                minutes of CLI calls, a dashboard nobody can read, and an app carrying settings
                that point at a laptop. A real project sends roughly sixty.
              </p>

              <div className="not-prose my-6 overflow-hidden rounded-xl border border-border">
                <table className="w-full text-sm">
                  <thead className="bg-muted/40">
                    <tr>
                      <th className="px-4 py-2.5 text-left font-medium">Sent</th>
                      <th className="px-4 py-2.5 text-left font-medium">Why</th>
                    </tr>
                  </thead>
                  <tbody>
                    {kept.map((r) => (
                      <tr key={r.group} className="border-t border-border">
                        <td className="px-4 py-2.5 align-top font-mono text-xs">{r.group}</td>
                        <td className="px-4 py-2.5 align-top text-muted-foreground">{r.what}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <div className="not-prose my-6 overflow-hidden rounded-xl border border-border">
                <table className="w-full text-sm">
                  <thead className="bg-muted/40">
                    <tr>
                      <th className="px-4 py-2.5 text-left font-medium">Held back</th>
                      <th className="px-4 py-2.5 text-left font-medium">Why</th>
                    </tr>
                  </thead>
                  <tbody>
                    {dropped.map((r) => (
                      <tr key={r.group} className="border-t border-border">
                        <td className="px-4 py-2.5 align-top font-mono text-xs">{r.group}</td>
                        <td className="px-4 py-2.5 align-top text-muted-foreground">{r.what}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <p>
                Anything Grit held back you can still set yourself, and it will not be overwritten:
              </p>
              <CodeBlock terminal code={`railway variable set SENTRY_DSN=... --service api`} />

              <h2 id="build">How it is built</h2>
              <p>
                From <code>apps/api/Dockerfile</code>, the same one <code>docker compose</code>{' '}
                uses, so there is one definition of how the API is built rather than two that
                drift. <code>apps/api/railway.json</code> says so, and adds the health check and
                restart policy:
              </p>
              <CodeBlock
                language="json"
                code={`{
  "$schema": "https://railway.com/railway.schema.json",
  "build": { "builder": "DOCKERFILE", "dockerfilePath": "Dockerfile" },
  "deploy": {
    "healthcheckPath": "/api/v1/health",
    "healthcheckTimeout": 60,
    "restartPolicyType": "ON_FAILURE",
    "restartPolicyMaxRetries": 5
  }
}`}
              />
              <p>
                The upload happens from <code>apps/api</code>, not the repository root, because the
                Dockerfile copies <code>go.mod</code> from the context root and the Go module is{' '}
                <code>apps/api</code>. Without that, Railway uploads the whole monorepo and the
                build cannot find the module.
              </p>

              <h2 id="after">After the first deploy</h2>
              <CodeBlock
                terminal
                code={`grit deploy --railway            # variables, upload, done
grit deploy --railway --no-domain  # skip asking for a URL
grit deploy --railway --service worker --environment staging`}
              />

              <h2 id="frontends">The web and admin apps</h2>
              <p>
                This deploys the API. The Next.js apps are separate Railway services, or a Vercel
                deploy, or anything that serves static files: point them at the API with{' '}
                <code>NEXT_PUBLIC_API_URL</code> and add that URL to <code>CORS_ORIGINS</code> on
                the API service.
              </p>

              <Callout type="warning" title="CORS_ORIGINS is the one to remember">
                The value in your <code>.env</code> is <code>http://localhost:3000</code>. It is
                pushed as-is, which means a deployed frontend is refused until you set the real
                origin. Grit cannot guess it, because it does not know where you will put the
                frontend.
              </Callout>
            </div>

            <div className="mt-12 flex items-center justify-between border-t border-border pt-6">
              <Button variant="ghost" asChild>
                <Link href="/docs/deployment">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  Deployment
                </Link>
              </Button>
              <Button variant="ghost" asChild>
                <Link href="/docs/security/doctor">
                  Project audit (grit doctor)
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
