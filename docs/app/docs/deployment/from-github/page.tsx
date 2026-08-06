import Link from 'next/link'
import { ArrowLeft, ArrowRight, GitBranch, FileCode, Rocket, AlertTriangle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'
import { DEPLOYMENT_PROVIDERS } from '@/config/deployment-providers'

export const metadata = getDocMetadata('/docs/deployment/from-github')

/* The shared prerequisite for every provider guide.
   Each provider page assumes the repository exists and the compose file is
   understood; this is where both of those actually happen, once, rather than
   being half-explained six times. */

export default function FromGitHubPage() {
  const gitProviders = DEPLOYMENT_PROVIDERS.filter((p) => p.kind !== 'vps')

  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Deployment</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Deploy from GitHub</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Every platform in these guides deploys the same way: it clones your repository,
                reads your compose file, builds the images and runs them. This page covers the two
                things all of them assume you have already done, and explains the compose file
                line by line so the provider pages can get on with the provider-specific part.
              </p>
            </div>

            {/* ── 1. The compose file ─────────────────────────────── */}
            <section className="mb-14">
              <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center gap-2">
                <FileCode className="h-5 w-5 text-primary" />
                The production compose file, line by line
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Every Grit project ships two compose files.{' '}
                <code>docker-compose.yml</code> is for your laptop: it bind-mounts source so hot
                reload works and publishes Postgres on a host port so you can attach a GUI.{' '}
                <code>docker-compose.prod.yml</code> is the one you deploy, and the differences are
                the whole point.
              </p>

              <CodeBlock
                language="yaml"
                filename="docker-compose.prod.yml"
                code={`services:
  api:
    build:
      context: ./apps/api
      dockerfile: Dockerfile
    restart: unless-stopped
    expose:
      - "8080"
    environment:
      DATABASE_URL: postgres://\${POSTGRES_USER}:\${POSTGRES_PASSWORD}@postgres:5432/\${POSTGRES_DB}?sslmode=disable
      REDIS_URL: redis://redis:6379
      JWT_SECRET: \${JWT_SECRET}
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started
    networks: [internal]

  web:
    build:
      context: ./apps/web
      dockerfile: Dockerfile
    restart: unless-stopped
    expose:
      - "3000"
    environment:
      NEXT_PUBLIC_API_URL: \${API_URL}
    depends_on: [api]
    networks: [internal]

  postgres:
    image: postgres:16-alpine
    restart: unless-stopped
    environment:
      POSTGRES_USER: \${POSTGRES_USER}
      POSTGRES_PASSWORD: \${POSTGRES_PASSWORD}
      POSTGRES_DB: \${POSTGRES_DB}
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U \${POSTGRES_USER}"]
      interval: 5s
      timeout: 5s
      retries: 10
    networks: [internal]

  redis:
    image: redis:7-alpine
    restart: unless-stopped
    volumes:
      - redis-data:/data
    networks: [internal]

volumes:
  postgres-data:
  redis-data:

networks:
  internal:
    driver: bridge`}
              />

              <div className="mt-6 space-y-5">
                <div>
                  <h3 className="font-semibold text-foreground mb-1.5">
                    <code>expose</code>, never <code>ports</code>
                  </h3>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <code>expose</code> opens a port to the other containers on the network.{' '}
                    <code>ports</code> publishes it on the host. Every platform here puts a reverse
                    proxy in front and terminates TLS there, so publishing a host port bypasses
                    both. On a plain VPS with no firewall it puts your database on the public
                    internet. The development compose file publishes Postgres deliberately, which
                    is exactly why you do not deploy that one.
                  </p>
                </div>

                <div>
                  <h3 className="font-semibold text-foreground mb-1.5">
                    Service names are hostnames
                  </h3>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <code>postgres:5432</code> works because Docker gives every service a DNS entry
                    on the shared network. This is why the connection string says{' '}
                    <code>@postgres</code> and not <code>@localhost</code>: inside a container,
                    localhost is that container. Getting this wrong is the single most common
                    self-hosted deployment question.
                  </p>
                </div>

                <div>
                  <h3 className="font-semibold text-foreground mb-1.5">
                    The healthcheck is what makes <code>depends_on</code> mean anything
                  </h3>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <code>depends_on</code> on its own waits for the container to <em>start</em>,
                    not for Postgres to accept connections, and those are several seconds apart.{' '}
                    <code>condition: service_healthy</code> plus the <code>pg_isready</code>{' '}
                    healthcheck is what actually holds the API back until the database answers.
                    Platforms that run real Compose honour this. Platforms that translate your file
                    into their own model do not, which is covered per provider below.
                  </p>
                </div>

                <div>
                  <h3 className="font-semibold text-foreground mb-1.5">
                    Named volumes are your database
                  </h3>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <code>postgres-data</code> survives <code>docker compose down</code> and a
                    server reboot. It does not survive <code>docker compose down -v</code>. That
                    one flag is the difference between a restart and losing everything, and it is
                    worth knowing before you type it at 2am.
                  </p>
                </div>

                <div>
                  <h3 className="font-semibold text-foreground mb-1.5">
                    <code>${'{VAR}'}</code> substitution, not hardcoded secrets
                  </h3>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    Compose substitutes <code>${'{POSTGRES_PASSWORD}'}</code> from a{' '}
                    <code>.env</code> sitting beside the compose file. Every platform below fills
                    that file from its own dashboard, which is how your secrets reach the
                    containers without ever being committed. Keep the <code>${'{VAR}'}</code>{' '}
                    syntax: it is the mechanism, not a formality.
                  </p>
                </div>
              </div>
            </section>

            {/* ── 2. GitHub ───────────────────────────────────────── */}
            <section className="mb-14">
              <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center gap-2">
                <GitBranch className="h-5 w-5 text-primary" />
                Getting the code on GitHub
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Every platform on this page except a plain VPS deploys by cloning a repository. It
                needs to exist before you start clicking around in a dashboard.
              </p>

              <CodeBlock
                language="bash"
                code={`# from your project root
git init
git add .
git commit -m "Initial commit"

# create the repo and push, using the GitHub CLI
gh repo create my-app --private --source=. --push

# or, if you made the repo in the browser first
git remote add origin git@github.com:you/my-app.git
git branch -M main
git push -u origin main`}
              />

              <div className="mt-5 flex gap-3 rounded-xl border border-warning/30 bg-warning/[0.05] px-5 py-4">
                <AlertTriangle className="h-5 w-5 shrink-0 text-warning mt-0.5" />
                <div>
                  <h3 className="font-semibold text-foreground mb-1">
                    Check what you just committed
                  </h3>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    Grit&apos;s <code>.gitignore</code> excludes <code>.env</code>, but only if you
                    have not renamed it. Run <code>git ls-files | grep -i env</code> before pushing.
                    A committed <code>.env</code> means rotating every secret in it, and a private
                    repository does not save you: it is in the history the moment a collaborator
                    clones it.
                  </p>
                </div>
              </div>

              <p className="text-muted-foreground leading-relaxed mt-5">
                What every platform needs from the repository is the same three things: a{' '}
                <code>docker-compose.prod.yml</code> at a path you can name, a{' '}
                <code>Dockerfile</code> per service that builds, and no secrets in the tree. Grit
                scaffolds the first two.
              </p>
            </section>

            {/* ── 3. Then pick a platform ─────────────────────────── */}
            <section className="mb-14">
              <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center gap-2">
                <Rocket className="h-5 w-5 text-primary" />
                Then pick a platform
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-5">
                They divide into two groups, and the division matters more than the branding. Some
                run the real Docker Compose engine, so your file behaves exactly as it does on your
                laptop. Others read your file and translate it into their own model, which changes
                hostnames and drops <code>depends_on</code> ordering. Neither is wrong; knowing
                which you picked saves an afternoon.
              </p>

              <div className="grid gap-2 sm:grid-cols-2">
                {gitProviders.map((p) => (
                  <Link
                    key={p.slug}
                    href={`/docs/deployment/${p.slug}`}
                    className="rounded-lg border border-border/40 px-4 py-3 hover:border-primary/40 hover:bg-muted/30 transition-colors"
                  >
                    <div className="font-medium text-foreground">{p.name}</div>
                    <div className="text-xs text-muted-foreground mt-0.5 line-clamp-2">
                      {p.tagline}
                    </div>
                  </Link>
                ))}
              </div>
            </section>

            <div className="flex items-center justify-between pt-6 border-t border-border/40">
              <Button variant="ghost" asChild>
                <Link href="/docs/deployment">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  All platforms
                </Link>
              </Button>
              <Button variant="ghost" asChild>
                <Link href="/docs/deployment/checklist">
                  Pre-launch checklist
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
