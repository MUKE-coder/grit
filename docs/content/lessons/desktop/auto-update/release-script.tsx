import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        The release script bumps the version, builds the binary +
        installer, tags git, and uploads everything to GitHub. One
        command, one new release that the auto-updater finds. Grit
        scaffolds this as <code>scripts/release-desktop.sh</code>.
      </p>

      <h2>The full pipeline</h2>
      <CodeBlock
        terminal
        code={`scripts/release-desktop.sh 1.2.3`}
      />
      <p>That runs:</p>
      <ol>
        <li>Bump <code>wails.json</code>&apos;s <code>productVersion</code> + <code>outputfilename</code></li>
        <li>Bump <code>version.go</code>&apos;s <code>AppVersion</code> constant</li>
        <li>Regenerate branded NSIS bitmaps from <code>icon.ico</code></li>
        <li>Download + cache the offline WebView2 runtime (first time only)</li>
        <li><code>wails build -nsis -platform windows/amd64 -trimpath -ldflags &apos;-s -w&apos;</code></li>
        <li>Run <code>makensis</code> on the slim variant</li>
        <li>Tag <code>v1.2.3</code> + git push</li>
        <li><code>gh release create</code> with three assets (raw .exe, full installer, slim installer)</li>
      </ol>

      <h2>Prerequisites</h2>
      <ul>
        <li>
          <code>bash</code>, <code>jq</code>, <code>python3</code> on PATH
          (Git Bash on Windows gives you all)
        </li>
        <li>
          <code>wails</code> CLI
        </li>
        <li>
          <code>makensis</code> — Windows NSIS install
        </li>
        <li>
          <code>gh</code> — GitHub CLI, signed in (<code>gh auth login</code>)
        </li>
        <li>
          <code>powershell.exe</code> on PATH (needed for the bitmap
          generation step)
        </li>
      </ul>

      <h2>Versioning</h2>
      <p>Semantic versioning:</p>
      <CodeBlock
        language="text"
        code={`v1.0.0    Initial release
v1.0.1    Bug fix only
v1.1.0    New feature, backward-compatible
v2.0.0    Breaking change`}
      />
      <p>
        For internal-only POS apps, you can be looser — many shops just
        use date-based versions like <code>2026.06.12</code>. The
        important rule is that versions sort lexically + the auto-updater
        recognises a higher value as newer.
      </p>

      <h2>Three assets per release</h2>
      <ul>
        <li>
          <strong>Raw .exe</strong>{' '}
          (<code>FieldPOS-v1.2.3.exe</code>) — what the auto-updater
          downloads and swaps. ~15 MB.
        </li>
        <li>
          <strong>Full installer</strong>{' '}
          (<code>FieldPOS-Setup-v1.2.3.exe</code>) — ~150 MB with bundled
          WebView2 runtime. Customer-facing.
        </li>
        <li>
          <strong>Slim installer</strong>{' '}
          (<code>FieldPOS-Setup-Slim-v1.2.3.exe</code>) — ~22 MB,
          downloads WebView2 online. Easier email attachment.
        </li>
      </ul>

      <TipBox tone="warning">
        <strong>Always upload the raw .exe.</strong> The auto-updater
        looks for it (not the installer). If you only upload installers,
        running v1.0 binaries can&apos;t find a swap target — the in-app
        update silently fails.
      </TipBox>

      <h2>GitHub Actions alternative</h2>
      <p>
        For teams, automate via GitHub Actions instead of running the
        script locally:
      </p>
      <CodeBlock
        language="yaml"
        filename=".github/workflows/release.yml"
        code={`on:
  push:
    tags: ['v*']

jobs:
  release:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - uses: pnpm/action-setup@v2
      - run: go install github.com/wailsapp/wails/v2/cmd/wails@latest
      - run: choco install nsis -y
      - run: scripts/release-desktop.sh \${{ github.ref_name }}
        env:
          GH_TOKEN: \${{ secrets.GITHUB_TOKEN }}`}
      />
      <p>
        Push a tag locally, GitHub Actions runs the same script in CI,
        the release appears on the Releases page, every running v1
        binary sees the new release on its next 6-hour check.
      </p>

      <h2>Code signing — optional but worth it</h2>
      <p>
        An unsigned .exe shows the Windows SmartScreen warning the first
        time a customer runs it. A code-signing cert (~$200/year) skips
        this. Add a signing step to the script:
      </p>
      <CodeBlock
        language="bash"
        code={`signtool sign /f cert.pfx /p "$CERT_PWD" /tr http://timestamp.digicert.com /td sha256 /fd sha256 \\
  "build/bin/FieldPOS-v$VERSION.exe"`}
      />

      <KnowledgeCheck
        question="You forgot to upload the raw .exe — only the installers are in the GitHub release. What happens for users running v1.0?"
        choices={[
          {
            label: 'Auto-updater downloads the installer and runs it',
            feedback:
              "Wrong — the updater downloads + swaps a raw binary. It can\'t run an installer mid-process. The installer would try to register itself in Add/Remove Programs, prompt for install location, etc.",
          },
          {
            label: 'Auto-updater fails to find a swap target — update silently fails',
            correct: true,
            feedback:
              "Right — the updater scans assets looking for FieldPOS-v1.2.3.exe (raw). If only installers are uploaded, no match. Users get no update banner; they\'re stuck on v1.0 until they manually download the installer.",
          },
          {
            label: 'Updater shows an error to the user',
            feedback:
              "It might log one in the console, but the user-visible behaviour is silence — no update banner, no error toast. Worse UX than an explicit error.",
          },
          {
            label: 'GitHub auto-converts the installer',
            feedback:
              "GitHub does nothing of the kind. The release contains exactly what you upload.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              For chapter 4&apos;s assignment, ship a release:
            </p>
            <ol>
              <li>Install v0.1.0 of your scaffolded app via the installer</li>
              <li>Make a small visible change (change a label, add a route)</li>
              <li>
                Run <code>scripts/release-desktop.sh 0.1.1</code>
              </li>
              <li>
                Confirm v0.1.1 appears on{' '}
                <code>github.com/&lt;you&gt;/&lt;repo&gt;/releases</code> with
                3 assets
              </li>
              <li>
                Open the installed v0.1.0 — the &quot;Update available&quot;
                banner should appear
              </li>
              <li>
                Click <strong>Install &amp; restart</strong>. App relaunches
                as v0.1.1.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            If the banner doesn&apos;t appear, check{' '}
            <code>UPDATER_GITHUB_OWNER</code> and{' '}
            <code>UPDATER_GITHUB_REPO</code> in your <code>.env</code> match
            where you published.
          </>
        }
        solution={
          <>
            <p>
              That cycle — code change, run script, watch in-app update —
              is your shipping rhythm forever. Customers get bug fixes
              within hours of you fixing them, with one click. Chapter 4
              done.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Final chapter — <strong>Branded installers</strong>. NSIS deep
        dive: full vs slim variants, branded MUI bitmaps from your icon.
      </p>
    </>
  )
}
