import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock, StepWithCode } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/infrastructure/deploy-command')

export default function DeployCommandPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="max-w-4xl mx-auto py-12 px-6 lg:px-8">
          <div className="mb-14">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">
              Infrastructure
            </p>
            <h1 className="text-3xl lg:text-4xl font-bold tracking-tight text-foreground mb-4">
              Deploy Command
            </h1>
            <p className="text-lg text-muted-foreground leading-relaxed max-w-2xl">
              One-command production deployment. Build, upload, and configure your server with{' '}
              <code className="text-primary bg-accent/30 px-1.5 py-0.5 rounded text-[15px]">grit deploy</code>.
            </p>
          </div>

          <div className="mb-14">
            <h2 className="text-xl font-semibold text-foreground mb-8">How it works</h2>

            <StepWithCode
              number="01"
              title="Build"
              description={
                <p>
                  Cross-compiles your Go binary for Linux (<code className="text-primary bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">GOOS=linux GOARCH=amd64 CGO_ENABLED=0</code>).
                  If a <code className="text-primary bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">frontend/</code> directory
                  exists, builds the frontend first.
                </p>
              }
              code={`grit deploy --host user@server.com --domain myapp.com`}
              filename="Terminal"
              language="bash"
            />

            <StepWithCode
              number="02"
              title="Upload"
              description={<p>Uploads the compiled binary to the remote server via SCP. Creates the deployment directory if it does not exist.</p>}
              code={`# Binary uploaded to /opt/my-app/my-app
# via: scp -P 22 bin/my-app user@server.com:/opt/my-app/`}
              filename="Terminal"
              language="bash"
            />

            <StepWithCode
              number="03"
              title="Systemd service"
              description={<p>Creates and enables a systemd service unit. The service auto-restarts on failure and reads environment variables from your <code className="text-primary bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">.env</code> file.</p>}
              code={`[Unit]
Description=my-app
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/my-app
ExecStart=/opt/my-app/my-app
Restart=on-failure
EnvironmentFile=/opt/my-app/.env

[Install]
WantedBy=multi-user.target`}
              filename="my-app.service"
              language="ini"
            />

            <StepWithCode
              number="04"
              title="Caddy reverse proxy"
              description={<p>If <code className="text-primary bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">--domain</code> is provided, configures Caddy as a reverse proxy with automatic HTTPS via Let{"'"}s Encrypt.</p>}
              code={`myapp.com {
    reverse_proxy localhost:8080
    encode gzip
    header {
        X-Frame-Options "DENY"
        X-Content-Type-Options "nosniff"
        -Server
    }
}`}
              filename="Caddyfile"
              language="nginx"
            />
          </div>

          {/* Flags */}
          <div className="mb-14">
            <h2 className="text-xl font-semibold text-foreground mb-4">Flags</h2>
            <div className="overflow-x-auto rounded-lg border border-border/40">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border/40 bg-accent/20">
                    <th className="px-4 py-3 text-left font-medium text-foreground/80">Flag</th>
                    <th className="px-4 py-3 text-left font-medium text-foreground/80">Env Var</th>
                    <th className="px-4 py-3 text-left font-medium text-foreground/80">Description</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border/30">
                  {[
                    { flag: '--host', env: 'DEPLOY_HOST', desc: 'SSH host (e.g. user@server.com)' },
                    { flag: '--port', env: '-', desc: 'SSH port (default: 22)' },
                    { flag: '--key', env: 'DEPLOY_KEY_FILE', desc: 'Path to SSH private key' },
                    { flag: '--domain', env: 'DEPLOY_DOMAIN', desc: 'Domain for Caddy reverse proxy' },
                    { flag: '--app-port', env: '-', desc: 'Port the app runs on (default: 8080)' },
                  ].map((row) => (
                    <tr key={row.flag} className="hover:bg-accent/20 transition-colors">
                      <td className="px-4 py-3 font-mono text-sm text-primary">{row.flag}</td>
                      <td className="px-4 py-3 font-mono text-xs text-muted-foreground/70">{row.env}</td>
                      <td className="px-4 py-3 text-muted-foreground">{row.desc}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          <div className="flex items-center justify-between pt-8 border-t border-border/40">
            <Button variant="ghost" size="sm" asChild className="text-muted-foreground/70 hover:text-foreground">
              <Link href="/docs/infrastructure/deployment" className="gap-1.5">
                <ArrowLeft className="h-3.5 w-3.5" />
                Deployment
              </Link>
            </Button>
            <Button variant="ghost" size="sm" asChild className="text-muted-foreground/70 hover:text-foreground">
              <Link href="/docs/infrastructure/dokploy" className="gap-1.5">
                Deploy with Dokploy
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
