import Link from 'next/link'
import {
  ArrowLeft,
  ArrowRight,
  Download,
  RefreshCw,
  Package,
  Box,
  ShieldCheck,
  Info,
  AlertTriangle,
  CheckCircle2,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/desktop/auto-update')

export default function DesktopAutoUpdatePage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Hero */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 inline-flex items-center gap-1.5">
                <RefreshCw className="h-3.5 w-3.5" />
                Desktop · Auto-update + installers
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4 leading-tight">
                Auto-update, full + slim installers
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Every <code>grit new-desktop</code> project ships with an in-app
                auto-updater (binary-swap, Wails-bound, modal UI), two Windows
                installers (full with bundled WebView2 + slim with online
                bootstrapper), and a one-shot release script that builds and
                publishes both to a GitHub release. Inspired by the production
                pattern from{' '}
                <Link
                  href="https://jb.desishub.com/blog/wails-desktop-auto-updater-github-releases"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-primary hover:underline"
                >
                  jb.desishub.com&apos;s walkthrough
                </Link>
                .
              </p>
            </div>

            {/* What ships */}
            <div className="mb-12 rounded-2xl border border-primary/20 bg-primary/[0.04] p-6">
              <p className="text-xs font-mono uppercase tracking-wider text-primary mb-3 inline-flex items-center gap-2">
                <Box className="h-3.5 w-3.5" /> What ships in every scaffolded desktop app
              </p>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                {[
                  { p: 'updater.go', d: 'Wails-bound binary-swap auto-updater (~430 lines)' },
                  { p: 'version.go', d: 'Single AppName / AppVersion source of truth' },
                  { p: 'frontend/src/lib/version.ts', d: 'Vite-injected APP_VERSION + semver helper' },
                  { p: 'frontend/src/hooks/use-update-checker.ts', d: 'React Query polling hook (6h cadence)' },
                  { p: 'frontend/src/components/update-modal.tsx', d: 'Download + install modal with progress bar' },
                  { p: 'frontend/src/components/update-banner.tsx', d: 'Dismissible header banner' },
                  { p: 'build/windows/installer/project.nsi', d: 'Full installer (~150 MB, bundled WebView2)' },
                  { p: 'build/windows/installer/project-slim.nsi', d: 'Slim installer (~22 MB, online WebView2)' },
                  { p: 'scripts/release-desktop.sh', d: 'One-shot release pipeline' },
                ].map((row) => (
                  <div key={row.p} className="text-sm">
                    <code className="text-primary/90 text-[12px]">{row.p}</code>
                    <p className="text-xs text-muted-foreground mt-0.5">{row.d}</p>
                  </div>
                ))}
              </div>
            </div>

            <div className="prose-grit">
              {/* ─── Section 1: The model ─────────────────────────── */}
              <h2 id="how-it-works">How auto-update works</h2>
              <p>
                The frontend polls a Go method that hits the GitHub releases
                API. If a newer version is available the banner appears; the user
                clicks <em>Update now</em>, the modal opens, the new binary
                downloads to <code>%TEMP%</code>, and the running .exe gets
                atomically swapped with it. The whole cycle is ~20 lines of
                glue from the React side because all the hard parts live in
                <code> updater.go</code>.
              </p>

              <div className="overflow-hidden rounded-xl border border-border/50 my-5">
                <table className="w-full text-sm">
                  <thead className="bg-muted/40">
                    <tr>
                      <th className="text-left px-4 py-2.5 font-semibold">
                        Wails method
                      </th>
                      <th className="text-left px-4 py-2.5 font-semibold">
                        Frontend call
                      </th>
                      <th className="text-left px-4 py-2.5 font-semibold">
                        What it does
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border/40">
                    {[
                      {
                        m: 'CheckForUpdates',
                        f: 'useUpdateChecker (auto)',
                        d: 'Hits GitHub /releases/latest, picks the right asset for OS/arch.',
                      },
                      {
                        m: 'StartDownload',
                        f: 'UpdateModal onMount',
                        d: 'Goroutine-backed fetch to %TEMP%, atomic counter for progress.',
                      },
                      {
                        m: 'GetProgress',
                        f: 'UpdateModal 500ms poll',
                        d: 'Returns {status, bytes_downloaded, bytes_total, error}.',
                      },
                      {
                        m: 'CancelDownload',
                        f: 'UpdateModal "Cancel"',
                        d: 'Aborts the in-flight context, deletes the .partial.',
                      },
                      {
                        m: 'InstallAndRestart',
                        f: 'UpdateModal "Install & restart"',
                        d: 'Swap binary, spawn new process, runtime.Quit after 250ms.',
                      },
                    ].map((r) => (
                      <tr key={r.m} className="hover:bg-card/40 transition-colors">
                        <td className="px-4 py-2.5 align-top font-mono text-[12.5px] text-primary/90">
                          {r.m}
                        </td>
                        <td className="px-4 py-2.5 align-top text-[13px] text-foreground/80">
                          <code>{r.f}</code>
                        </td>
                        <td className="px-4 py-2.5 align-top text-[13px] text-muted-foreground">
                          {r.d}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <h3 id="binary-swap">The binary-swap trick (per OS)</h3>
              <p>
                Different OSes have different rules about overwriting a running
                executable. The updater branches accordingly:
              </p>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-3 not-prose my-5">
                <div className="rounded-xl border border-border/50 bg-card/30 p-4">
                  <p className="text-xs font-mono uppercase tracking-wider text-primary mb-1.5">
                    Linux / macOS
                  </p>
                  <p className="text-sm text-foreground/85 leading-relaxed">
                    POSIX inode semantics keep the running process&apos;s file
                    alive even when the path is overwritten. <code>copyFile</code>{' '}
                    straight onto the path — done. New invocations get the new
                    binary; the current process keeps using the old inode until
                    it exits.
                  </p>
                </div>
                <div className="rounded-xl border border-border/50 bg-card/30 p-4">
                  <p className="text-xs font-mono uppercase tracking-wider text-primary mb-1.5">
                    Windows
                  </p>
                  <p className="text-sm text-foreground/85 leading-relaxed">
                    .exe is locked while running, but Windows{' '}
                    <em>allows the directory entry to move</em> even with the
                    kernel object open. We rename <code>app.exe → app.exe.old</code>,
                    copy the new binary into place, spawn it, exit. Next startup
                    deletes the <code>.old</code> in a background goroutine
                    (retries 25× at 200ms intervals to ride out antivirus locks).
                  </p>
                </div>
              </div>

              <Callout tone="info">
                If anything goes wrong on Windows mid-swap (e.g.{' '}
                <code>copyFile</code> fails because <code>%TEMP%</code> is on a
                different volume than the install dir — different filesystem,{' '}
                <code>os.Rename</code> returns <code>EXDEV</code>), the updater
                renames <code>app.exe.old</code> back to <code>app.exe</code> so
                the user is never stranded without a working binary. Industry
                standard pattern (it&apos;s what{' '}
                <code>inconshreveable/go-update</code>, <code>gh</code>,{' '}
                <code>bun upgrade</code>, and <code>rustup</code> all do).
              </Callout>

              {/* ─── Section 2: Configuration ──────────────────── */}
              <h2 id="configure">Configure your release source</h2>
              <p>
                By default the updater points at the GitHub repo derived from
                your project name. Override with env vars if your repo lives
                somewhere else or you want to gate releases through a proxy:
              </p>

              <CodeBlock
                filename=".env"
                language="dotenv"
                code={`# Required if your repo name differs from the project name
UPDATER_GITHUB_OWNER=acme-co
UPDATER_GITHUB_REPO=my-desktop-app

# Optional: PAT for private repos / higher rate limits
UPDATER_GITHUB_TOKEN=

# Optional: point at your own proxy instead of api.github.com
# (Lets you gate by min_supported_version, hide a PAT,
#  flip an emergency override env var without a redeploy.)
UPDATER_PROXY_URL=https://my-api.com/api/desktop`}
              />

              <p>
                For most projects the GitHub defaults are enough. The proxy
                option matches the production pattern from the blog — useful
                when you have a paid customer base and want server-side gating.
              </p>

              {/* ─── Section 3: Two installers ─────────────────── */}
              <h2 id="installers">Two Windows installers: full vs slim</h2>
              <p>
                Wails ships one installer template. We add a second so you can
                pick the right trade-off for each distribution channel.
              </p>

              <div className="overflow-hidden rounded-xl border border-border/50 my-5">
                <table className="w-full text-sm">
                  <thead className="bg-muted/40">
                    <tr>
                      <th className="text-left px-4 py-2.5 font-semibold">Variant</th>
                      <th className="text-left px-4 py-2.5 font-semibold">Size</th>
                      <th className="text-left px-4 py-2.5 font-semibold">WebView2 strategy</th>
                      <th className="text-left px-4 py-2.5 font-semibold">When to use</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border/40">
                    <tr className="hover:bg-card/40 transition-colors">
                      <td className="px-4 py-2.5 align-top font-mono text-[12.5px] text-primary/90">
                        project.nsi (full)
                      </td>
                      <td className="px-4 py-2.5 align-top text-[13px] text-foreground/80">
                        ~150 MB
                      </td>
                      <td className="px-4 py-2.5 align-top text-[13px] text-muted-foreground">
                        Bundles the Microsoft Evergreen Standalone (~125 MB) and
                        installs it silently if missing.
                      </td>
                      <td className="px-4 py-2.5 align-top text-[13px] text-muted-foreground">
                        Field sites, offline machines, restrictive corporate
                        firewalls, one-shot USB-stick installs.
                      </td>
                    </tr>
                    <tr className="hover:bg-card/40 transition-colors">
                      <td className="px-4 py-2.5 align-top font-mono text-[12.5px] text-primary/90">
                        project-slim.nsi (slim)
                      </td>
                      <td className="px-4 py-2.5 align-top text-[13px] text-foreground/80">
                        ~22 MB
                      </td>
                      <td className="px-4 py-2.5 align-top text-[13px] text-muted-foreground">
                        Uses Wails&apos;s built-in online bootstrapper (~1.8 MB)
                        — downloads the actual runtime from Microsoft at install
                        time.
                      </td>
                      <td className="px-4 py-2.5 align-top text-[13px] text-muted-foreground">
                        Web downloads, email links, anywhere bandwidth matters.
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <p>
                Both installers default to{' '}
                <code>%LOCALAPPDATA%\Programs\&lt;ProjectName&gt;</code> — the
                Slack / Discord / VSCode convention. No admin required, so the
                in-app auto-updater can swap the binary without UAC.
              </p>

              {/* ─── Section 4: Release pipeline ───────────────── */}
              <h2 id="release-pipeline">Cutting a release</h2>
              <p>
                <code>scripts/release-desktop.sh</code> is the only command you
                need to know:
              </p>

              <CodeBlock
                terminal
                code={`scripts/release-desktop.sh 1.2.3`}
              />

              <p>That does, in order:</p>
              <ol>
                <li>Bumps <code>wails.json</code> productVersion + outputfilename</li>
                <li>Bumps <code>version.go</code>&apos;s <code>AppVersion</code> constant</li>
                <li>(Re)generates branded NSIS bitmaps from <code>icon.ico</code> via PowerShell</li>
                <li>Downloads + caches the offline WebView2 runtime (~125 MB, first time only)</li>
                <li>Runs <code>wails build -nsis -platform windows/amd64 -trimpath -ldflags &apos;-s -w&apos;</code></li>
                <li>Runs <code>makensis</code> on <code>project-slim.nsi</code> with the same raw binary</li>
                <li>Tags <code>v1.2.3</code>, pushes the tag</li>
                <li><code>gh release create</code> uploads three assets: raw .exe (auto-updater uses this), full installer, slim installer</li>
              </ol>

              <Callout tone="warning">
                <strong>Requires:</strong> <code>bash</code>, <code>jq</code>,{' '}
                <code>python3</code>, <code>wails</code>, <code>makensis</code>{' '}
                (NSIS), <code>gh</code> (GitHub CLI), and{' '}
                <code>powershell.exe</code> on PATH (Git Bash on Windows gives
                you all of these except NSIS — install with{' '}
                <code>winget install NSIS.NSIS</code>).
              </Callout>

              {/* ─── Section 5: Frontend wiring ─────────────────── */}
              <h2 id="frontend">Frontend wiring (auto-mounted)</h2>
              <p>
                The scaffold mounts <code>&lt;UpdateBanner /&gt;</code> in your
                root layout already, so nothing for you to wire up. The hook
                polls every 6 hours; on focus / dismiss patterns are handled
                inside the banner. If you want a manual&nbsp;
                &quot;Check for updates&quot; button in Settings, drop in:
              </p>

              <CodeBlock
                filename="frontend/src/routes/_layout/settings.tsx"
                language="tsx"
                code={`import { useUpdateChecker } from '@/hooks/use-update-checker'
import { APP_VERSION } from '@/lib/version'

export function SettingsPage() {
  const { isChecking, checkNow, update, isUpdateAvailable } = useUpdateChecker()

  return (
    <div className="space-y-4">
      <p>Current version: <strong>v{APP_VERSION}</strong></p>
      <button
        onClick={() => checkNow()}
        disabled={isChecking}
        className="rounded-md bg-primary text-primary-foreground px-3 py-1.5 text-sm disabled:opacity-50"
      >
        {isChecking ? 'Checking…' : 'Check for updates'}
      </button>
      {isUpdateAvailable && update && (
        <p className="text-sm text-success">v{update.version} is available.</p>
      )}
    </div>
  )
}`}
              />

              {/* ─── Section 6: Code signing ───────────────────── */}
              <h2 id="code-signing">Code signing (optional but recommended)</h2>
              <p>
                An unsigned .exe shows the Windows SmartScreen warning the first
                time a user runs it (<em>&quot;Windows protected your PC&quot;</em>) — you
                lose installs and gain support tickets. A code-signing cert is
                ~$200/year from any reputable CA (DigiCert, Sectigo,
                SSL.com). Once you have a <code>.pfx</code>, add a step to
                the release script before the final{' '}
                <code>gh release create</code>:
              </p>

              <CodeBlock
                filename="scripts/release-desktop.sh (extra step)"
                language="bash"
                code={`# Sign the raw .exe and both installers
SIGN_CERT="${'$'}{HOME}/.certs/code-signing.pfx"
SIGN_PASS="${'$'}{CODE_SIGN_PASSWORD}"   # from env, not committed

for exe in "${'$'}{RAW_EXE}" "${'$'}{TARGET_FULL}" "${'$'}{SLIM_TARGET}"; do
    signtool.exe sign \\
        /f "${'$'}{SIGN_CERT}" \\
        /p "${'$'}{SIGN_PASS}" \\
        /tr http://timestamp.digicert.com \\
        /td sha256 \\
        /fd sha256 \\
        "${'$'}exe"
done`}
              />

              <p>
                For macOS notarization the same pattern applies with{' '}
                <code>codesign</code> + <code>xcrun notarytool</code> — covered
                in the dedicated macOS release guide (coming soon).
              </p>

              {/* ─── Section 7: Bonus: emergency override ─────── */}
              <h2 id="emergency-rollback">Emergency rollback / override</h2>
              <p>
                If you push a bad release and need to roll back without yanking
                the GitHub release: set <code>UPDATER_PROXY_URL</code> to a
                custom endpoint that returns the previous version. Or — if
                you&apos;re already using a proxy — flip an{' '}
                <code>OVERRIDE_VERSION</code> env var on your server and every
                client picks up the pinned version on its next 6-hour poll.
              </p>

              {/* CTA strip */}
              <div className="mt-14 rounded-2xl border border-primary/20 bg-primary/[0.04] p-6 not-prose">
                <h3 className="text-lg font-semibold mb-2">Ready to ship</h3>
                <p className="text-sm text-muted-foreground mb-4">
                  <code>grit new-desktop my-app</code> gives you all of this
                  out of the box. Then{' '}
                  <code>scripts/release-desktop.sh 0.1.1</code> publishes your
                  first release and the auto-updater takes over from there.
                </p>
                <div className="flex flex-wrap gap-2.5">
                  <Button
                    asChild
                    className="rounded-full bg-primary hover:bg-primary/90 text-primary-foreground"
                  >
                    <Link href="/docs/desktop">
                      Desktop guide
                      <ArrowRight className="ml-1.5 h-3.5 w-3.5" />
                    </Link>
                  </Button>
                  <Button asChild variant="outline" className="rounded-full border-border/60">
                    <Link href="/docs/tech-kits/desktop">Desktop tech kit</Link>
                  </Button>
                </div>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 mt-10 border-t border-border/30">
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="text-muted-foreground/60 hover:text-foreground"
              >
                <Link href="/docs/desktop" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Desktop overview
                </Link>
              </Button>
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="text-muted-foreground/60 hover:text-foreground"
              >
                <Link href="/docs/tech-kits/desktop" className="gap-1.5">
                  Desktop tech kit
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}

function Callout({ tone, children }: { tone: 'info' | 'warning'; children: React.ReactNode }) {
  const Icon = tone === 'warning' ? AlertTriangle : Info
  const styles =
    tone === 'warning'
      ? 'border-amber-500/30 bg-amber-500/[0.05]'
      : 'border-sky-500/30 bg-sky-500/[0.05]'
  const iconColor = tone === 'warning' ? 'text-amber-500' : 'text-sky-500'
  return (
    <div className={`not-prose rounded-xl border ${styles} p-4 my-4 flex gap-3 text-sm leading-relaxed`}>
      <Icon className={`h-4 w-4 ${iconColor} mt-0.5 shrink-0`} />
      <div className="text-foreground/85">{children}</div>
    </div>
  )
}
