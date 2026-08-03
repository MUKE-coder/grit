import Link from 'next/link'
import { AlertCircle, ArrowLeft, ArrowRight } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/deployment/environment')

const REQUIRED = [
  { name: 'DATABASE_URL', note: 'Postgres connection string. Also accepts sqlite:./app.db for small single-binary deploys.' },
  { name: 'JWT_SECRET', note: 'Signs access tokens. 32+ random bytes. Changing it logs everyone out.' },
  { name: 'APP_ENV', note: 'Set to production. Controls error verbosity, cookie flags and whether debug routes mount.' },
  { name: 'APP_PORT', note: 'Defaults to 8080. Platforms that assign a port (Railway, Render, Heroku) inject their own — map it.' },
  { name: 'CORS_ORIGINS', note: 'Comma-separated exact origins. Never a wildcard in production with credentials on.' },
]

const CONDITIONAL = [
  { name: 'REDIS_URL', when: 'Background jobs, caching, or rate limiting are on' },
  { name: 'RESEND_API_KEY', when: 'The app sends email' },
  { name: 'S3_* / R2_*', when: 'File uploads are enabled' },
  { name: 'OAUTH_GOOGLE_* / OAUTH_GITHUB_*', when: 'Social login is enabled' },
  { name: 'AI_GATEWAY_API_KEY', when: 'The AI module is on' },
]

export default function EnvironmentPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <Link
              href="/docs/deployment"
              className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors mb-8"
            >
              <ArrowLeft className="h-3.5 w-3.5" />
              Deployment
            </Link>

            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Deployment</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Environment variables</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Which ones the app refuses to start without, which depend on what you have
                switched on, and the category that fails silently.
              </p>
            </div>

            {/* The build-time distinction first. It is the one that produces a
                deploy that looks successful and is broken. */}
            <div className="rounded-xl border border-amber-500/30 bg-amber-500/[0.05] p-6 mb-12">
              <div className="flex gap-3">
                <AlertCircle className="h-5 w-5 shrink-0 text-amber-500 mt-0.5" />
                <div>
                  <h2 className="font-semibold mb-2">Build-time vs runtime: read this one first</h2>
                  <p className="text-sm text-muted-foreground leading-relaxed mb-3">
                    Anything the frontend reads while it compiles is <strong>baked into the
                    output</strong>. Set it only at runtime and it is not missing — it is
                    the empty string, compiled in, with nothing logged and nothing failing.
                    The page renders, the API calls go to the wrong host, and the deploy
                    looks green.
                  </p>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    On every platform with a separate &ldquo;build arguments&rdquo; field,
                    that is where <code className="text-xs">NEXT_PUBLIC_*</code> and anything
                    read during prerender belongs. Secrets go in runtime environment, never
                    build args — build logs are a wider audience than a container.
                  </p>
                </div>
              </div>
            </div>

            <h2 className="text-2xl font-bold tracking-tight mb-4">Always required</h2>
            <div className="overflow-x-auto rounded-xl border border-border/50 mb-12">
              <table className="w-full text-sm">
                <thead className="bg-card/60">
                  <tr className="text-left">
                    <th className="px-4 py-3 font-semibold">Variable</th>
                    <th className="px-4 py-3 font-semibold">Notes</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border/40">
                  {REQUIRED.map((v) => (
                    <tr key={v.name}>
                      <td className="px-4 py-3 font-mono text-[13px] whitespace-nowrap align-top">{v.name}</td>
                      <td className="px-4 py-3 text-muted-foreground leading-relaxed">{v.note}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <h2 className="text-2xl font-bold tracking-tight mb-4">Required only if you use it</h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Each optional module can also be switched off entirely with{' '}
              <code className="text-xs">MODULE_&lt;NAME&gt;=false</code>, which is cleaner
              than supplying dummy credentials to satisfy a check.
            </p>
            <div className="overflow-x-auto rounded-xl border border-border/50 mb-12">
              <table className="w-full text-sm">
                <thead className="bg-card/60">
                  <tr className="text-left">
                    <th className="px-4 py-3 font-semibold">Variable</th>
                    <th className="px-4 py-3 font-semibold">Needed when</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border/40">
                  {CONDITIONAL.map((v) => (
                    <tr key={v.name}>
                      <td className="px-4 py-3 font-mono text-[13px] whitespace-nowrap align-top">{v.name}</td>
                      <td className="px-4 py-3 text-muted-foreground leading-relaxed">{v.when}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <h2 className="text-2xl font-bold tracking-tight mb-4">Generating secrets</h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Do not reuse the development values. The <code className="text-xs">.env.example</code>{' '}
              in a fresh project contains placeholders that are identical in every Grit
              project ever generated.
            </p>
            <CodeBlock
              language="bash"
              code={`# 32 random bytes, base64. One per environment.
openssl rand -base64 32

# Or without openssl:
head -c 32 /dev/urandom | base64`}
            />

            <h2 className="text-2xl font-bold tracking-tight mt-12 mb-4">Checking before you ship</h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              The API validates its configuration at startup and refuses to boot on a
              missing required value rather than failing later on the first request that
              needed it.
            </p>
            <CodeBlock
              language="bash"
              code={`# Same validation the server runs, without starting it.
grit doctor

# And confirm what the running app actually received:
curl -s https://your-domain.com/api/health`}
            />

            <div className="flex items-center justify-between gap-4 border-t border-border/40 pt-8 mt-12">
              <Link
                href="/docs/deployment"
                className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
              >
                <ArrowLeft className="h-3.5 w-3.5" />
                Deployment
              </Link>
              <Link
                href="/docs/deployment/build-locally"
                className="inline-flex items-center gap-1.5 text-sm font-medium text-primary hover:underline"
              >
                Test the build locally
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
