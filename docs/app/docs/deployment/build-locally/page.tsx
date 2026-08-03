import Link from 'next/link'
import { AlertCircle, ArrowLeft, ArrowRight, Check } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/deployment/build-locally')

export default function BuildLocallyPage() {
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
              <h1 className="text-4xl font-bold tracking-tight mb-4">Test the production build locally</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">Reproduce the deploy on your own machine first. Most deployment failures are visible here, two minutes in, instead of ten minutes into someone else&apos;s build log.</p>
            </div>

            <p className="text-muted-foreground leading-relaxed mb-6">
              Almost every failed deploy is reproducible on your own machine in two minutes.
              The production build differs from <code className="text-xs">dev</code> in ways
              that matter: it type-checks strictly, it inlines build-time environment, it
              tree-shakes, and it runs without the dev server&apos;s forgiving module
              resolution. A change that works under <code className="text-xs">dev</code> and
              fails in <code className="text-xs">build</code> is common, and finding that out
              from a platform&apos;s build log is a slow way to learn it.
            </p>

            <h2 className="text-2xl font-bold tracking-tight mb-4">1. Build the API image</h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              This is the same Dockerfile every platform here uses, so if it builds locally it
              builds there.
            </p>
            <CodeBlock
              language="bash"
              code={`docker build -f apps/api/Dockerfile -t my-app-api .

# Run it against your local Postgres to confirm it boots:
docker run --rm -p 8080:8080 \
  -e DATABASE_URL="postgres://user:pass@host.docker.internal:5432/mydb" \
  -e JWT_SECRET="$(openssl rand -base64 32)" \
  -e APP_ENV=production \
  my-app-api`}
            />

            <h2 className="text-2xl font-bold tracking-tight mt-12 mb-4">2. Build the frontend the way the platform will</h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Run the real build, not the dev server. Pass the build-time variables exactly as
              the platform will — this is the step that catches a missing{' '}
              <code className="text-xs">NEXT_PUBLIC_*</code> before it ships as an empty
              string.
            </p>
            <CodeBlock
              language="bash"
              code={`cd apps/web
NEXT_PUBLIC_API_URL=https://api.your-domain.com pnpm build
pnpm start`}
            />

            <div className="rounded-xl border border-amber-500/30 bg-amber-500/[0.05] p-5 my-8">
              <div className="flex gap-3">
                <AlertCircle className="h-4 w-4 shrink-0 text-amber-500 mt-0.5" />
                <p className="text-sm text-muted-foreground leading-relaxed">
                  Do this with your <code className="text-xs">.env.local</code> temporarily
                  renamed. Otherwise you are testing with variables the deployment does not
                  have, and the run proves nothing about production.
                </p>
              </div>
            </div>

            <h2 className="text-2xl font-bold tracking-tight mt-12 mb-4">3. Bring up the whole stack</h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              The closest local approximation of a server: production images, no bind mounts,
              no exposed database port.
            </p>
            <CodeBlock
              language="bash"
              code={`docker compose -f docker-compose.prod.yml up --build

# In another shell:
curl -s localhost:8080/api/health`}
            />

            <h2 className="text-2xl font-bold tracking-tight mt-12 mb-4">What this catches</h2>
            <ul className="space-y-2.5 mb-8">
              {[
                'Type errors that dev mode never surfaced, because dev does not type-check on every save.',
                'A build-time variable read at runtime — the empty-string failure that produces a green deploy and a broken app.',
                'A cgo dependency that breaks the static binary. It compiles on your machine and refuses to run in the container.',
                'Missing files in the image. A .dockerignore that excludes something the build needs fails only inside Docker.',
                'A health endpoint that is not where the platform is looking, which reads as "the deploy never finishes".',
              ].map((t) => (
                <li key={t} className="flex gap-3 text-sm text-muted-foreground leading-relaxed">
                  <Check className="h-4 w-4 shrink-0 text-emerald-500 mt-0.5" />
                  {t}
                </li>
              ))}
            </ul>

            <div className="flex items-center justify-between gap-4 border-t border-border/40 pt-8 mt-12">
              <Link href="/docs/deployment/environment" className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors">
                <ArrowLeft className="h-3.5 w-3.5" />
                Environment variables
              </Link>
              <Link href="/docs/deployment/checklist" className="inline-flex items-center gap-1.5 text-sm font-medium text-primary hover:underline">
                Go-live checklist
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </div>

          </div>
        </div>
      </main>
    </div>
  )
}
