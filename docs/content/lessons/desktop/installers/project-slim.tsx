import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        The slim installer (<code>project-slim.nsi</code>) is ~22 MB. It
        uses the WebView2 <em>online bootstrapper</em> — a 1.8 MB exe
        that downloads the full runtime from Microsoft at install time.
        Right choice when bandwidth at install time matters.
      </p>

      <h2>When slim wins</h2>
      <ul>
        <li>
          Email distribution — most providers cap attachments at 25 MB
        </li>
        <li>
          Web downloads where the conversion drops with file size (every
          extra MB hurts)
        </li>
        <li>
          Users on metered cellular at install — they&apos;d rather download
          22 MB once and let WebView2 install over Wi-Fi later
        </li>
        <li>
          GitHub release size limits (2 GB per asset is generous, but for
          older repos, smaller assets fit better)
        </li>
      </ul>

      <h2>The shape — minus the WebView2 bundle</h2>
      <CodeBlock
        language="nsi"
        filename="build/windows/installer/project-slim.nsi (excerpt)"
        code={`# Same Welcome/Directory/Install/Finish flow
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

# Wails' default bootstrapper macro — downloads WebView2 at install
!insertmacro wails.webview2runtime

Section
    !insertmacro wails.setShellContext
    !insertmacro wails.files                   # main binary
    !insertmacro wails.writeUninstaller
SectionEnd`}
      />
      <p>
        <code>wails.webview2runtime</code> drops a 1.8 MB
        <code>MicrosoftEdgeWebview2Setup.exe</code> bootstrapper. At
        install time, it runs that bootstrapper which contacts Microsoft
        and downloads the ~125 MB runtime if not present.
      </p>

      <h2>What happens on the customer&apos;s machine</h2>
      <ol>
        <li>Customer downloads the 22 MB slim installer.</li>
        <li>Runs it. NSIS prompts: Welcome → Directory → click Install.</li>
        <li>NSIS extracts the main binary + bootstrapper.</li>
        <li>
          Bootstrapper checks if WebView2 is installed. If yes, skip. If
          no, download from Microsoft (needs internet).
        </li>
        <li>Total install time: 30s if WebView2 already there; 2-3 min on a slow connection if not.</li>
      </ol>

      <TipBox tone="warning">
        <strong>Customer with no internet at install time runs the slim
        installer.</strong> The bootstrapper fails. The user sees a
        WebView2 install error mid-wizard. They&apos;re confused. Ship the
        full installer in that case.
      </TipBox>

      <h2>When to ship full vs slim</h2>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Channel</th>
              <th className="text-left px-3 py-2 font-medium">Use</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr><td className="px-3 py-2">Web download from your site</td><td>Slim — convert better; users typically online</td></tr>
            <tr><td className="px-3 py-2">Email attachment</td><td>Slim — fits in &lt; 25 MB caps</td></tr>
            <tr><td className="px-3 py-2">USB stick / on-prem install</td><td>Full — assume no internet</td></tr>
            <tr><td className="px-3 py-2">Air-gapped customer</td><td>Full</td></tr>
            <tr><td className="px-3 py-2">In-app auto-update</td><td>Neither — uses the raw .exe (chapter 4)</td></tr>
          </tbody>
        </table>
      </div>

      <h2>Ship both, let the customer pick</h2>
      <p>
        Grit&apos;s release script publishes both as GitHub release assets.
        Your download page can offer both with a brief description of
        when to use each. Most customers pick slim; the full is there for
        the edge case.
      </p>
      <CodeBlock
        language="text"
        filename="Your download page"
        code={`Download FieldPOS

Recommended:
[ Download FieldPOS-Setup-Slim-v1.2.3.exe ]  22 MB

Don't have internet during install?
[ Download FieldPOS-Setup-v1.2.3.exe (full) ]  150 MB`}
      />

      <KnowledgeCheck
        question="A customer-success ticket: 'Installer hangs at 60%.' Customer is on a satellite internet connection in a remote area. They downloaded the slim installer. What's most likely happening?"
        choices={[
          {
            label: 'WebView2 bootstrapper is downloading the runtime over a slow link — appears hung but is actually working',
            correct: true,
            feedback:
              "Right — 125 MB over satellite can take 30+ minutes. The bootstrapper doesn\'t show download progress in NSIS\'s UI; it just shows the install bar stuck. The fix: ship them the full installer (a one-time download that they can keep for re-installs).",
          },
          {
            label: 'NSIS bug — re-run the installer',
            feedback:
              "Wrong — re-running re-starts the same 125 MB download. The issue is the connection speed, not NSIS.",
          },
          {
            label: 'Corrupted download — verify the SHA',
            feedback:
              "Possible but unlikely. The hang at 60% pattern (right when WebView2 starts) points strongly to the runtime download.",
          },
          {
            label: 'Windows Update is interfering',
            feedback:
              "Unlikely. Windows Update doesn\'t typically block WebView2 install.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Compare the two installers side-by-side:</p>
            <ol>
              <li>Build both via your release script</li>
              <li>
                List the file sizes:{' '}
                <code>FieldPOS-Setup-v0.1.0.exe</code> +{' '}
                <code>FieldPOS-Setup-Slim-v0.1.0.exe</code>
              </li>
              <li>
                Install the slim version on a connected PC. Time how
                long it takes from double-click to first launch.
              </li>
              <li>Uninstall.</li>
              <li>
                Repeat with the full version on a disconnected PC.
              </li>
            </ol>
            <p>Paste your numbers in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            For a connected dev box where WebView2 is already installed,
            the slim install will be just a few seconds — that&apos;s
            because the bootstrapper skips when WebView2 is present.
            Test on a fresh VM for the realistic case.
          </>
        }
        solution={
          <>
            <p>Typical numbers:</p>
            <CodeBlock
              language="text"
              code={`Slim:  22 MB,  10s install (WebView2 already installed)
              45s install (WebView2 download needed)
Full: 150 MB,  20s install (always — extracts the bundled runtime)`}
            />
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Last lesson — the <strong>branded MUI bitmaps</strong>. Make the
        installer look like YOUR product, not a default NSIS template.
      </p>
    </>
  )
}
