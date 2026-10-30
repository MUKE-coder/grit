import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock, StepWithCode } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/deployment/deploy-command')

export default function DeployCommandPage() {
  return (
    <DocsShell>
          <div className="mb-14">
            <p className="text-sm font-mono font-medium text-accent mb-3 tracking-wide uppercase">
              Infrastructure
            </p>
            <h1 className="text-3xl lg:text-4xl font-bold tracking-tight text-foreground mb-4">
              Deploy Command
            </h1>
            <p className="text-lg text-muted-foreground leading-relaxed max-w-2xl">
              One-command production deployment. Build, upload, and configure your server with{' '}
              <code className="text-accent bg-accent/30 px-1.5 py-0.5 rounded text-[15px]">grit deploy</code>.
            </p>
            <LaneFlow
              id="infra-deploycmd"
              lanes={['grit deploy', 'On your server']}
              nodes={[
                { id: 'cmd', lane: 0, row: 1, title: 'grit deploy', sub: 'one command', tone: 'primary' },
                { id: 'build', lane: 1, row: 0, title: 'Build', sub: 'binary + assets', tone: 'cyan', badge: 1 },
                { id: 'upload', lane: 1, row: 1, title: 'Upload', sub: 'to the server', tone: 'blue', badge: 2 },
                { id: 'configure', lane: 1, row: 2, title: 'Configure', sub: 'systemd · Caddy', tone: 'amber', badge: 3 },
                { id: 'switch', lane: 1, row: 3, title: 'Switch', sub: 'current -> release', tone: 'violet', badge: 4 },
                { id: 'check', lane: 1, row: 4, title: 'Health check', sub: 'or roll back', tone: 'green', badge: 5 },
              ]}
              edges={[
                { from: 'cmd', to: 'build', label: 'runs', tone: 'cyan' },
                { from: 'build', to: 'upload', tone: 'blue' },
                { from: 'upload', to: 'configure', tone: 'amber' },
                { from: 'configure', to: 'switch', tone: 'violet' },
                { from: 'switch', to: 'check', tone: 'green' },
                { from: 'check', to: 'switch', label: 'fails: back one release', tone: 'amber' },
              ]}
              legend={[
                { tone: 'primary', label: 'One command' },
                { tone: 'green', label: 'Live server' },
              ]}
              caption="grit deploy builds, uploads, switches and checks: a release that does not answer is put back"
            />
          </div>

          <div className="mb-14">
            <h2 className="text-xl font-semibold text-foreground mb-8">How it works</h2>

            <StepWithCode
              number="01"
              title="Build"
              description={
                <p>
                  Cross-compiles your Go binary for Linux (<code className="text-accent bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">GOOS=linux GOARCH=amd64 CGO_ENABLED=0</code>).
                  If a <code className="text-accent bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">frontend/</code> directory
                  exists, builds the frontend first.
                </p>
              }
              code={`grit deploy --host user@server.com --domain myapp.com`}
              filename="Terminal"
              language="bash"
            />

            <StepWithCode
              number="02"
              title="Upload into a new release"
              description={
                <p>
                  Each deploy gets a timestamped directory of its own, and{' '}
                  <code className="text-accent bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">current</code>{' '}
                  is a symlink to one of them. Nothing is ever written over the binary
                  that is running, which is what used to make the second deploy to a
                  Linux box fail with <code className="text-accent bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">ETXTBSY</code>.
                  The <code className="text-accent bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">.env</code> goes
                  up with it, at 0600, owned by the service user.
                </p>
              }
              code={`/opt/my-app/releases/20261030-150405/
  my-app
  grit-migrate
  .env                      (0600, www-data)

/opt/my-app/current -> releases/20261030-150405`}
              filename="On the server"
              language="text"
            />

            <StepWithCode
              number="03"
              title="Migrate, then switch"
              description={
                <p>
                  If the project has a{' '}
                  <code className="text-accent bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">cmd/migrate</code>,
                  it is built and uploaded alongside the server and run{' '}
                  <strong>before</strong> the symlink moves: the old binary is still
                  serving, and a migration that fails must not leave the new one running
                  against a schema it does not expect. The switch itself is one symlink,
                  so a request in flight sees the old release or the new one and never
                  half of a copy.
                </p>
              }
              code={`cd /opt/my-app/releases/20261030-150405 && ./grit-migrate
ln -sfn /opt/my-app/releases/20261030-150405 /opt/my-app/current
systemctl restart my-app`}
              filename="On the server"
              language="bash"
            />

            <StepWithCode
              number="04"
              title="Systemd service"
              description={
                <p>
                  Written once, and it names{' '}
                  <code className="text-accent bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">current</code>{' '}
                  rather than a release, so a deploy moves a symlink instead of rewriting
                  and reloading the unit. The service restarts on failure and reads its
                  configuration from the{' '}
                  <code className="text-accent bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">.env</code>{' '}
                  inside the release.
                </p>
              }
              code={`[Unit]
Description=my-app
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/my-app/current
ExecStart=/opt/my-app/current/my-app
Restart=on-failure
EnvironmentFile=/opt/my-app/current/.env

# Security hardening
NoNewPrivileges=yes
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=/opt/my-app

[Install]
WantedBy=multi-user.target`}
              filename="my-app.service"
              language="ini"
            />

            <StepWithCode
              number="05"
              title="Caddy reverse proxy"
              description={
                <p>
                  With <code className="text-accent bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">--domain</code>,
                  this app{"'"}s site block is written to{' '}
                  <code className="text-accent bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">/etc/caddy/sites/my-app.caddy</code>{' '}
                  and the main Caddyfile gets one{' '}
                  <code className="text-accent bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">import sites/*.caddy</code>{' '}
                  line appended if it does not have one. Appended, never rewritten: that
                  file is yours, and every other site in it stays. The configuration is
                  validated before the reload, because a Caddyfile that does not parse
                  takes every site on the box down. HTTPS is automatic, through Let{"'"}s
                  Encrypt.
                </p>
              }
              code={`myapp.com {
    reverse_proxy localhost:8080
    encode gzip
    header {
        X-Frame-Options "DENY"
        X-Content-Type-Options "nosniff"
        -Server
    }
}`}
              filename="/etc/caddy/sites/my-app.caddy"
              language="nginx"
            />

            <StepWithCode
              number="06"
              title="Health check, or roll back"
              description={
                <p>
                  The health endpoint is fetched on the server itself, through the app{"'"}s
                  own port, for up to a minute. From the server deliberately: a deploy
                  that works and a certificate that has not issued yet are different
                  problems and should not look the same. A release that never answers is
                  put back, and the failed release is left on disk so you can read its
                  logs. Old releases are pruned only after the new one has answered.
                </p>
              }
              code={`  -> Checking the app answers on /api/health...
  !  the app answered 500 on /api/health for 60 seconds
  -> Rolling back to /opt/my-app/releases/20261029-094512...

the new release did not come up, so the server was rolled back to the
one before it. The failed release is still at
/opt/my-app/releases/20261030-150405. Read why with:
  ssh user@server.com 'journalctl -u my-app -n 50'`}
              filename="Terminal"
              language="text"
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
                    { flag: '--env-file', env: '-', desc: "Local .env uploaded into the release (default: the API's own)" },
                    { flag: '--health-path', env: '-', desc: 'Probed after the restart (default: /api/health)' },
                    { flag: '--keep', env: '-', desc: 'How many old releases to keep (default: 5)' },
                    { flag: '--rollback', env: '-', desc: 'Put the previous release back, building and uploading nothing' },
                    { flag: '--dry-run', env: '-', desc: 'Print every command that would run, and run none of them' },
                  ].map((row) => (
                    <tr key={row.flag} className="hover:bg-surface-2 transition-colors">
                      <td className="px-4 py-3 font-mono text-sm text-accent">{row.flag}</td>
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
              <Link href="/docs/deployment/vps" className="gap-1.5">
                <ArrowLeft className="h-3.5 w-3.5" />
                Deployment
              </Link>
            </Button>
            <Button variant="ghost" size="sm" asChild className="text-muted-foreground/70 hover:text-foreground">
              <Link href="/docs/deployment/dokploy" className="gap-1.5">
                Deploy with Dokploy
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </Button>
          </div>
    </DocsShell>
  )
}
