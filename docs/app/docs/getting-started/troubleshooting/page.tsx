import Link from 'next/link'
import { ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'

export default function TroubleshootingPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Getting Started</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Troubleshooting
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Common errors you may encounter when setting up or running a Grit project,
                and how to fix them.
              </p>
            </div>

            {/* Port Conflicts */}
            <div className="prose-grit mb-12">
              <h2>Port Already in Use</h2>
              <p>
                This is the most common error when starting Docker services or dev servers.
                It means another process or container is already using the port Grit needs.
              </p>
            </div>

            {/* Port 5432 */}
            <div className="mb-8">
              <div className="rounded-xl border border-destructive/30 bg-destructive/5 overflow-hidden">
                <div className="px-5 py-3 border-b border-destructive/20 bg-destructive/10">
                  <code className="text-sm font-mono text-destructive">
                    Bind for 0.0.0.0:5432 failed: port is already allocated
                  </code>
                </div>
                <div className="p-5 space-y-4">
                  <div className="prose-grit">
                    <p>
                      <strong>Cause:</strong> Another PostgreSQL instance or Docker container is using port 5432.
                    </p>
                    <p><strong>Solutions:</strong></p>
                  </div>
                  <div className="space-y-3">
                    <div className="rounded-lg border border-border/30 bg-card/30 p-4">
                      <h4 className="text-sm font-semibold mb-2">Option 1: Stop the conflicting container</h4>
                      <div className="rounded-lg bg-card/80 border border-border/40 p-4 font-mono text-sm space-y-2">
                        <div className="text-muted-foreground/40 text-xs">{`# Find what's using the port`}</div>
                        <div>
                          <span className="text-primary/50 select-none">$ </span>
                          <span className="text-foreground/80">{`docker ps --format "table {{.Names}}\\t{{.Ports}}" | grep 5432`}</span>
                        </div>
                        <div className="text-muted-foreground/40 text-xs mt-3">{`# Stop the conflicting container`}</div>
                        <div>
                          <span className="text-primary/50 select-none">$ </span>
                          <span className="text-foreground/80">docker stop &lt;container-name&gt;</span>
                        </div>
                      </div>
                    </div>
                    <div className="rounded-lg border border-border/30 bg-card/30 p-4">
                      <h4 className="text-sm font-semibold mb-2">Option 2: Stop a local PostgreSQL service</h4>
                      <div className="rounded-lg bg-card/80 border border-border/40 p-4 font-mono text-sm space-y-2">
                        <div className="text-muted-foreground/40 text-xs">{`# macOS (Homebrew)`}</div>
                        <div>
                          <span className="text-primary/50 select-none">$ </span>
                          <span className="text-foreground/80">brew services stop postgresql@16</span>
                        </div>
                        <div className="text-muted-foreground/40 text-xs mt-3">{`# Linux (systemd)`}</div>
                        <div>
                          <span className="text-primary/50 select-none">$ </span>
                          <span className="text-foreground/80">sudo systemctl stop postgresql</span>
                        </div>
                        <div className="text-muted-foreground/40 text-xs mt-3">{`# Windows`}</div>
                        <div>
                          <span className="text-primary/50 select-none">&gt; </span>
                          <span className="text-foreground/80">net stop postgresql-x64-16</span>
                        </div>
                      </div>
                    </div>
                    <div className="rounded-lg border border-border/30 bg-card/30 p-4">
                      <h4 className="text-sm font-semibold mb-2">Option 3: Change the port in docker-compose.yml</h4>
                      <div className="prose-grit">
                        <p>
                          Edit <code>docker-compose.yml</code> and change the host port (left side):
                        </p>
                      </div>
                      <div className="rounded-lg bg-card/80 border border-border/40 p-4 font-mono text-sm mt-2">
                        <div className="text-muted-foreground/40 text-xs">{`# Change "5432:5432" to "5433:5432"`}</div>
                        <div className="mt-1">
                          <span className="text-foreground/80">ports:</span>
                        </div>
                        <div>
                          <span className="text-foreground/80">  - &quot;5433:5432&quot;</span>
                        </div>
                      </div>
                      <div className="prose-grit mt-2">
                        <p>
                          Then update <code>.env</code> to match: <code>DB_PORT=5433</code>
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Port 6379 - Redis */}
            <div className="mb-8">
              <div className="rounded-xl border border-destructive/30 bg-destructive/5 overflow-hidden">
                <div className="px-5 py-3 border-b border-destructive/20 bg-destructive/10">
                  <code className="text-sm font-mono text-destructive">
                    Bind for 0.0.0.0:6379 failed: port is already allocated
                  </code>
                </div>
                <div className="p-5 space-y-4">
                  <div className="prose-grit">
                    <p>
                      <strong>Cause:</strong> Another Redis instance or Docker container is using port 6379.
                    </p>
                    <p><strong>Fix:</strong> Same approach as PostgreSQL -- find and stop the conflicting
                    container, stop the local Redis service, or remap the port.</p>
                  </div>
                  <div className="rounded-lg bg-card/80 border border-border/40 p-4 font-mono text-sm space-y-2">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">{`docker ps --format "table {{.Names}}\\t{{.Ports}}" | grep 6379`}</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker stop &lt;container-name&gt;</span>
                    </div>
                  </div>
                  <div className="prose-grit">
                    <p>
                      Or to stop a locally installed Redis:
                    </p>
                  </div>
                  <div className="rounded-lg bg-card/80 border border-border/40 p-4 font-mono text-sm space-y-2">
                    <div className="text-muted-foreground/40 text-xs">{`# macOS`}</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">brew services stop redis</span>
                    </div>
                    <div className="text-muted-foreground/40 text-xs mt-3">{`# Linux`}</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">sudo systemctl stop redis-server</span>
                    </div>
                    <div className="text-muted-foreground/40 text-xs mt-3">{`# Windows`}</div>
                    <div>
                      <span className="text-primary/50 select-none">&gt; </span>
                      <span className="text-foreground/80">net stop Redis</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Port 3000/3001 */}
            <div className="mb-8">
              <div className="rounded-xl border border-destructive/30 bg-destructive/5 overflow-hidden">
                <div className="px-5 py-3 border-b border-destructive/20 bg-destructive/10">
                  <code className="text-sm font-mono text-destructive">
                    Port 3000 is already in use
                  </code>
                </div>
                <div className="p-5 space-y-4">
                  <div className="prose-grit">
                    <p>
                      <strong>Cause:</strong> Another Next.js app, React dev server, or other
                      process is running on port 3000 (web) or 3001 (admin).
                    </p>
                    <p><strong>Fix:</strong> Find and kill the process:</p>
                  </div>
                  <div className="rounded-lg bg-card/80 border border-border/40 p-4 font-mono text-sm space-y-2">
                    <div className="text-muted-foreground/40 text-xs">{`# macOS / Linux`}</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">lsof -i :3000</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">kill -9 &lt;PID&gt;</span>
                    </div>
                    <div className="text-muted-foreground/40 text-xs mt-3">{`# Windows`}</div>
                    <div>
                      <span className="text-primary/50 select-none">&gt; </span>
                      <span className="text-foreground/80">netstat -ano | findstr :3000</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">&gt; </span>
                      <span className="text-foreground/80">taskkill /PID &lt;PID&gt; /F</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Port 8080 */}
            <div className="mb-8">
              <div className="rounded-xl border border-destructive/30 bg-destructive/5 overflow-hidden">
                <div className="px-5 py-3 border-b border-destructive/20 bg-destructive/10">
                  <code className="text-sm font-mono text-destructive">
                    listen tcp :8080: bind: address already in use
                  </code>
                </div>
                <div className="p-5 space-y-4">
                  <div className="prose-grit">
                    <p>
                      <strong>Cause:</strong> Another Go server or process is already running on port 8080.
                    </p>
                    <p><strong>Fix:</strong> Kill the existing process or change the port in <code>.env</code>:</p>
                  </div>
                  <div className="rounded-lg bg-card/80 border border-border/40 p-4 font-mono text-sm space-y-2">
                    <div className="text-muted-foreground/40 text-xs">{`# Find and kill the process`}</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">lsof -i :8080</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">kill -9 &lt;PID&gt;</span>
                    </div>
                    <div className="text-muted-foreground/40 text-xs mt-3">{`# Or change the port in .env`}</div>
                    <div>
                      <span className="text-foreground/80">PORT=8081</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Nuclear option */}
            <div className="mb-12">
              <div className="rounded-xl border border-primary/30 bg-primary/5 overflow-hidden">
                <div className="px-5 py-3 border-b border-primary/20 bg-primary/10">
                  <span className="text-sm font-semibold text-primary">Nuclear option: stop ALL Docker containers</span>
                </div>
                <div className="p-5">
                  <div className="prose-grit mb-3">
                    <p>
                      If you have many stale containers from other projects hogging ports, you can
                      stop them all at once:
                    </p>
                  </div>
                  <div className="rounded-lg bg-card/80 border border-border/40 p-4 font-mono text-sm space-y-2">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">docker stop $(docker ps -q)</span>
                    </div>
                  </div>
                  <div className="prose-grit mt-3">
                    <p>
                      This stops every running container. Then retry <code>docker compose up -d</code>
                      in your Grit project.
                    </p>
                  </div>
                </div>
              </div>
            </div>

            {/* Go / API Errors */}
            <div className="prose-grit mb-6">
              <h2>Go API Errors</h2>
            </div>

            {/* go mod tidy errors */}
            <div className="mb-8">
              <div className="rounded-xl border border-destructive/30 bg-destructive/5 overflow-hidden">
                <div className="px-5 py-3 border-b border-destructive/20 bg-destructive/10">
                  <code className="text-sm font-mono text-destructive">
                    go: module not found / cannot find module providing package
                  </code>
                </div>
                <div className="p-5 space-y-4">
                  <div className="prose-grit">
                    <p>
                      <strong>Cause:</strong> Go module cache is stale or dependencies haven&apos;t been downloaded yet.
                    </p>
                    <p><strong>Fix:</strong></p>
                  </div>
                  <div className="rounded-lg bg-card/80 border border-border/40 p-4 font-mono text-sm space-y-2">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">cd apps/api</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">go mod tidy</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">go mod download</span>
                    </div>
                  </div>
                  <div className="prose-grit">
                    <p>
                      If that doesn&apos;t work, clear the module cache and try again:
                    </p>
                  </div>
                  <div className="rounded-lg bg-card/80 border border-border/40 p-4 font-mono text-sm">
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">go clean -modcache &amp;&amp; go mod tidy</span>
                  </div>
                </div>
              </div>
            </div>

            {/* Database connection */}
            <div className="mb-8">
              <div className="rounded-xl border border-destructive/30 bg-destructive/5 overflow-hidden">
                <div className="px-5 py-3 border-b border-destructive/20 bg-destructive/10">
                  <code className="text-sm font-mono text-destructive">
                    failed to connect to host=localhost: dial tcp connection refused
                  </code>
                </div>
                <div className="p-5 space-y-4">
                  <div className="prose-grit">
                    <p>
                      <strong>Cause:</strong> The Go API can&apos;t connect to PostgreSQL. Either Docker
                      isn&apos;t running or the database container hasn&apos;t started yet.
                    </p>
                    <p><strong>Fix:</strong></p>
                    <ol>
                      <li>Make sure Docker is running: <code>docker ps</code></li>
                      <li>Start infrastructure: <code>docker compose up -d</code></li>
                      <li>Wait a few seconds for PostgreSQL to initialize, then retry</li>
                      <li>Check the container logs: <code>docker compose logs postgres</code></li>
                    </ol>
                    <p>
                      Also verify your <code>.env</code> has the correct database URL:
                    </p>
                  </div>
                  <div className="rounded-lg bg-card/80 border border-border/40 p-4 font-mono text-sm">
                    <span className="text-foreground/80">DATABASE_URL=postgres://grit:grit@localhost:5432/grit_dev?sslmode=disable</span>
                  </div>
                </div>
              </div>
            </div>

            {/* Redis connection */}
            <div className="mb-8">
              <div className="rounded-xl border border-destructive/30 bg-destructive/5 overflow-hidden">
                <div className="px-5 py-3 border-b border-destructive/20 bg-destructive/10">
                  <code className="text-sm font-mono text-destructive">
                    redis: dial tcp localhost:6379: connection refused
                  </code>
                </div>
                <div className="p-5 space-y-4">
                  <div className="prose-grit">
                    <p>
                      <strong>Cause:</strong> Redis container isn&apos;t running. Grit&apos;s batteries
                      (cache, jobs, cron) need Redis.
                    </p>
                    <p><strong>Fix:</strong> Start the Redis container:</p>
                  </div>
                  <div className="rounded-lg bg-card/80 border border-border/40 p-4 font-mono text-sm">
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">docker compose up -d redis</span>
                  </div>
                  <div className="prose-grit">
                    <p>
                      The API will still work without Redis -- it logs a warning and disables
                      caching/jobs/cron features gracefully.
                    </p>
                  </div>
                </div>
              </div>
            </div>

            {/* Frontend Errors */}
            <div className="prose-grit mb-6">
              <h2>Frontend Errors</h2>
            </div>

            {/* pnpm install fails */}
            <div className="mb-8">
              <div className="rounded-xl border border-destructive/30 bg-destructive/5 overflow-hidden">
                <div className="px-5 py-3 border-b border-destructive/20 bg-destructive/10">
                  <code className="text-sm font-mono text-destructive">
                    ERR_PNPM_NO_MATCHING_VERSION / pnpm: command not found
                  </code>
                </div>
                <div className="p-5 space-y-4">
                  <div className="prose-grit">
                    <p>
                      <strong>Cause:</strong> pnpm is not installed, or an outdated version is
                      installed.
                    </p>
                    <p><strong>Fix:</strong></p>
                  </div>
                  <div className="rounded-lg bg-card/80 border border-border/40 p-4 font-mono text-sm space-y-2">
                    <div className="text-muted-foreground/40 text-xs">{`# Install pnpm globally`}</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">npm install -g pnpm</span>
                    </div>
                    <div className="text-muted-foreground/40 text-xs mt-3">{`# Or use corepack (built into Node.js 16+)`}</div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">corepack enable</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">corepack prepare pnpm@latest --activate</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Module not found */}
            <div className="mb-8">
              <div className="rounded-xl border border-destructive/30 bg-destructive/5 overflow-hidden">
                <div className="px-5 py-3 border-b border-destructive/20 bg-destructive/10">
                  <code className="text-sm font-mono text-destructive">
                    Module not found: Can&apos;t resolve &apos;@/components/...&apos;
                  </code>
                </div>
                <div className="p-5 space-y-4">
                  <div className="prose-grit">
                    <p>
                      <strong>Cause:</strong> Node modules haven&apos;t been installed yet, or the
                      install was interrupted.
                    </p>
                    <p><strong>Fix:</strong> Run install from the project root:</p>
                  </div>
                  <div className="rounded-lg bg-card/80 border border-border/40 p-4 font-mono text-sm space-y-2">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">cd myapp</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">pnpm install</span>
                    </div>
                  </div>
                  <div className="prose-grit">
                    <p>
                      If that doesn&apos;t work, delete <code>node_modules</code> and the lockfile, then
                      reinstall:
                    </p>
                  </div>
                  <div className="rounded-lg bg-card/80 border border-border/40 p-4 font-mono text-sm space-y-2">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">rm -rf node_modules apps/*/node_modules pnpm-lock.yaml</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">pnpm install</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* CORS errors */}
            <div className="mb-8">
              <div className="rounded-xl border border-destructive/30 bg-destructive/5 overflow-hidden">
                <div className="px-5 py-3 border-b border-destructive/20 bg-destructive/10">
                  <code className="text-sm font-mono text-destructive">
                    Access to XMLHttpRequest blocked by CORS policy
                  </code>
                </div>
                <div className="p-5 space-y-4">
                  <div className="prose-grit">
                    <p>
                      <strong>Cause:</strong> The browser is blocking requests from the frontend
                      (localhost:3000) to the API (localhost:8080) due to CORS.
                    </p>
                    <p><strong>Fix:</strong> Make sure your <code>.env</code> includes the frontend origins:</p>
                  </div>
                  <div className="rounded-lg bg-card/80 border border-border/40 p-4 font-mono text-sm">
                    <span className="text-foreground/80">CORS_ORIGINS=http://localhost:3000,http://localhost:3001</span>
                  </div>
                  <div className="prose-grit">
                    <p>
                      Then restart the Go API. The CORS middleware reads this value on startup.
                    </p>
                  </div>
                </div>
              </div>
            </div>

            {/* Docker errors */}
            <div className="prose-grit mb-6">
              <h2>Docker Errors</h2>
            </div>

            {/* Docker not running */}
            <div className="mb-8">
              <div className="rounded-xl border border-destructive/30 bg-destructive/5 overflow-hidden">
                <div className="px-5 py-3 border-b border-destructive/20 bg-destructive/10">
                  <code className="text-sm font-mono text-destructive">
                    Cannot connect to the Docker daemon. Is the docker daemon running?
                  </code>
                </div>
                <div className="p-5 space-y-4">
                  <div className="prose-grit">
                    <p>
                      <strong>Cause:</strong> Docker Desktop is not running.
                    </p>
                    <p><strong>Fix:</strong></p>
                    <ul>
                      <li><strong>macOS/Windows:</strong> Open Docker Desktop from the Applications menu</li>
                      <li><strong>Linux:</strong> <code>sudo systemctl start docker</code></li>
                    </ul>
                    <p>
                      Wait for Docker to fully start (you&apos;ll see the whale icon in the
                      system tray), then retry <code>docker compose up -d</code>.
                    </p>
                  </div>
                </div>
              </div>
            </div>

            {/* Volume permissions */}
            <div className="mb-8">
              <div className="rounded-xl border border-destructive/30 bg-destructive/5 overflow-hidden">
                <div className="px-5 py-3 border-b border-destructive/20 bg-destructive/10">
                  <code className="text-sm font-mono text-destructive">
                    permission denied while trying to connect to the Docker daemon socket
                  </code>
                </div>
                <div className="p-5 space-y-4">
                  <div className="prose-grit">
                    <p>
                      <strong>Cause:</strong> Your user doesn&apos;t have permission to use Docker (Linux only).
                    </p>
                    <p><strong>Fix:</strong> Add your user to the docker group:</p>
                  </div>
                  <div className="rounded-lg bg-card/80 border border-border/40 p-4 font-mono text-sm space-y-2">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">sudo usermod -aG docker $USER</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">newgrp docker</span>
                    </div>
                  </div>
                  <div className="prose-grit">
                    <p>
                      You may need to log out and log back in for the group change to take effect.
                    </p>
                  </div>
                </div>
              </div>
            </div>

            {/* Code Generation Errors */}
            <div className="prose-grit mb-6">
              <h2>Code Generation Errors</h2>
            </div>

            {/* Not in a Grit project */}
            <div className="mb-8">
              <div className="rounded-xl border border-destructive/30 bg-destructive/5 overflow-hidden">
                <div className="px-5 py-3 border-b border-destructive/20 bg-destructive/10">
                  <code className="text-sm font-mono text-destructive">
                    Error: not a Grit project (no apps/api directory found)
                  </code>
                </div>
                <div className="p-5 space-y-4">
                  <div className="prose-grit">
                    <p>
                      <strong>Cause:</strong> You ran <code>grit generate</code> from outside
                      the project root.
                    </p>
                    <p><strong>Fix:</strong> Make sure you&apos;re in the project root directory
                    (the folder that contains <code>docker-compose.yml</code>, <code>apps/</code>,
                    and <code>packages/</code>):</p>
                  </div>
                  <div className="rounded-lg bg-card/80 border border-border/40 p-4 font-mono text-sm space-y-2">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">cd myapp</span>
                    </div>
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">grit generate resource Post --fields &quot;title:string&quot;</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Resource already exists */}
            <div className="mb-8">
              <div className="rounded-xl border border-destructive/30 bg-destructive/5 overflow-hidden">
                <div className="px-5 py-3 border-b border-destructive/20 bg-destructive/10">
                  <code className="text-sm font-mono text-destructive">
                    Error: resource &quot;Post&quot; already exists
                  </code>
                </div>
                <div className="p-5 space-y-4">
                  <div className="prose-grit">
                    <p>
                      <strong>Cause:</strong> A model file for this resource already exists in
                      <code>apps/api/internal/models/</code>.
                    </p>
                    <p><strong>Fix:</strong> If you want to regenerate it, delete the existing files
                    first. The code generator won&apos;t overwrite files to protect your custom code.</p>
                  </div>
                </div>
              </div>
            </div>

            {/* General tips */}
            <div className="prose-grit mb-12">
              <h2>General Tips</h2>
              <ul>
                <li>
                  <strong>Always start Docker first.</strong> Run <code>docker compose up -d</code> before
                  starting the Go API. The API needs PostgreSQL and Redis to be available.
                </li>
                <li>
                  <strong>Run <code>go mod tidy</code> after scaffolding.</strong> The first time you
                  scaffold a project, run <code>go mod tidy</code> in <code>apps/api</code> to download
                  all Go dependencies.
                </li>
                <li>
                  <strong>Check the <code>.env</code> file.</strong> Most connection errors are caused by
                  incorrect values in <code>.env</code>. Compare it with <code>.env.example</code> to
                  make sure all values are set.
                </li>
                <li>
                  <strong>Restart after <code>.env</code> changes.</strong> The Go API reads environment
                  variables on startup. If you change <code>.env</code>, restart the server.
                </li>
                <li>
                  <strong>Use <code>docker compose logs</code></strong> to debug container issues.
                  For example: <code>docker compose logs postgres</code> or <code>docker compose logs redis</code>.
                </li>
                <li>
                  <strong>Windows users:</strong> Make sure Docker Desktop has WSL 2 integration enabled.
                  Some port binding issues on Windows are resolved by restarting Docker Desktop.
                </li>
              </ul>
            </div>

            {/* Still stuck */}
            <div className="mb-8">
              <div className="rounded-xl border border-primary/30 bg-primary/5 p-6">
                <h3 className="text-lg font-semibold mb-2">Still stuck?</h3>
                <div className="prose-grit">
                  <p>
                    If you&apos;re encountering an error not listed here, open an issue on GitHub
                    with your error message, OS, and the command you ran. The community and
                    maintainers are happy to help.
                  </p>
                </div>
              </div>
            </div>

            {/* Nav */}
            <div className="flex flex-wrap gap-3 mt-12 pt-6 border-t border-border/30">
              <Button variant="outline" asChild className="border-border/60 bg-transparent hover:bg-accent/50">
                <Link href="/docs/getting-started/configuration">
                  Configuration
                </Link>
              </Button>
              <Button asChild className="glow-purple-sm ml-auto">
                <Link href="/docs/concepts/architecture">
                  Architecture
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
