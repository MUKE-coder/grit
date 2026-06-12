import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        The installer&apos;s Welcome page and header strip are bitmaps —
        small PNG-as-BMP images. By default they&apos;re NSIS stock. With
        ~5 lines of script, they become your product&apos;s look.
      </p>

      <h2>The two bitmaps</h2>
      <ul>
        <li>
          <strong><code>header.bmp</code></strong> — 150 × 57 px. Sits in
          the upper-right of every install page after Welcome. Usually
          contains your logo on a brand background.
        </li>
        <li>
          <strong><code>welcome.bmp</code></strong> — 164 × 314 px. Left-
          side banner on the Welcome and Finish pages. Tall, narrow,
          eye-catching.
        </li>
      </ul>

      <h2>Wiring them into the NSI</h2>
      <CodeBlock
        language="nsi"
        filename="build/windows/installer/project.nsi"
        code={`!define MUI_HEADERIMAGE
!define MUI_HEADERIMAGE_BITMAP        "header.bmp"
!define MUI_HEADERIMAGE_UNBITMAP      "header.bmp"
!define MUI_HEADERIMAGE_RIGHT
!define MUI_WELCOMEFINISHPAGE_BITMAP  "welcome.bmp"
!define MUI_UNWELCOMEFINISHPAGE_BITMAP "welcome.bmp"
!define MUI_BGCOLOR                   "F5F1EA"
!define MUI_TEXTCOLOR                 "0E0E0E"`}
      />
      <p>
        Background colour + text colour also customisable. Match your
        brand palette.
      </p>

      <h2>Generate from icon.ico — no design step</h2>
      <p>
        The Grit release script regenerates both bitmaps from your{' '}
        <code>icon.ico</code> automatically. PowerShell does the work:
      </p>
      <CodeBlock
        language="powershell"
        filename="scripts/release-desktop.sh (excerpt)"
        code={`powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "
Add-Type -AssemblyName System.Drawing
function Save-Bmp(\\$icon, \\$w, \\$h, \\$out) {
    \\$bmp = New-Object System.Drawing.Bitmap \\$w, \\$h
    \\$g = [System.Drawing.Graphics]::FromImage(\\$bmp)
    \\$g.Clear([System.Drawing.Color]::White)
    \\$src = [System.Drawing.Icon]::ExtractAssociatedIcon('build/windows/icon.ico')
    \\$g.DrawIcon(\\$src, 8, 8)
    \\$bmp.Save(\\$out, [System.Drawing.Imaging.ImageFormat]::Bmp)
}
Save-Bmp 'build/windows/icon.ico' 150 57 'build/windows/installer/header.bmp'
Save-Bmp 'build/windows/icon.ico' 164 314 'build/windows/installer/welcome.bmp'"`}
      />
      <p>
        Net effect: change your <code>icon.ico</code>, run the release
        script, both bitmaps reflect the new icon. No Photoshop trips.
      </p>

      <h2>Going further — custom bitmaps</h2>
      <p>
        For a polished look, replace the auto-generated bitmaps with
        designer-made ones. Drop them at:
      </p>
      <ul>
        <li><code>build/windows/installer/header.bmp</code></li>
        <li><code>build/windows/installer/welcome.bmp</code></li>
      </ul>
      <p>
        The release script preserves them if they already exist (or you
        comment out the regeneration step).
      </p>

      <TipBox tone="warning">
        <strong>Real BMP, not PNG renamed.</strong> NSIS requires actual
        BMP format. Export from Photoshop / GIMP / Figma as BMP, not PNG
        with the file extension changed. Bad BMPs render as black
        rectangles in the wizard.
      </TipBox>

      <h2>The icon itself</h2>
      <p>
        Your <code>icon.ico</code> is a multi-resolution file — typically
        16, 32, 48, 256 px embedded. Generate from a square PNG with
        an online icon generator (e.g.,{' '}
        <a href="https://icoconvert.com" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">icoconvert.com</a>).
      </p>
      <p>This icon ends up in:</p>
      <ul>
        <li>The installer&apos;s wizard window</li>
        <li>The shortcut on the Start menu / Desktop</li>
        <li>Add/Remove Programs entry</li>
        <li>The Wails app window in Task Manager</li>
      </ul>

      <h2>Text strings</h2>
      <p>
        The Welcome page text comes from MUI defines:
      </p>
      <CodeBlock
        language="nsi"
        code={`!define MUI_WELCOMEPAGE_TITLE "Welcome to the FieldPOS Setup Wizard"
!define MUI_WELCOMEPAGE_TEXT  "This will install FieldPOS v\${INFO_PRODUCTVERSION}.$\\r$\\n$\\r$\\nClick Next to continue."

!define MUI_DIRECTORYPAGE_TEXT_TOP "Setup will install FieldPOS in the folder below."
!define MUI_DIRECTORYPAGE_TEXT_DESTINATION "Install location"

!define MUI_FINISHPAGE_RUN       "$INSTDIR\\FieldPOS.exe"
!define MUI_FINISHPAGE_RUN_TEXT  "Launch FieldPOS now"`}
      />

      <KnowledgeCheck
        question="You replace icon.ico with a new design. You build the installer. The header bitmap still shows the OLD icon. What did you skip?"
        choices={[
          {
            label: 'Re-running the release script (it regenerates the bitmaps from icon.ico)',
            correct: true,
            feedback:
              "Right — bitmaps are generated, not live-loaded. The script reads icon.ico and writes header.bmp + welcome.bmp. Skip that step and the OLD bitmaps stick around.",
          },
          {
            label: 'NSIS rebuild',
            feedback:
              "Half right — you do need to rebuild the installer, but the underlying issue is the bitmaps weren\'t regenerated. The release script handles both.",
          },
          {
            label: 'Restarting your editor',
            feedback:
              "Irrelevant — NSIS reads files at build time, not editor cache.",
          },
          {
            label: 'NSIS doesn\'t cache bitmaps',
            feedback:
              "Correct that there\'s no cache, but the old bitmap file still exists. Regeneration is the fix.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>For chapter 5&apos;s assignment, ship a fully-branded installer:</p>
            <ol>
              <li>Replace <code>build/windows/icon.ico</code> with your real product icon</li>
              <li>
                Run the release script — confirm <code>header.bmp</code>{' '}
                and <code>welcome.bmp</code> were regenerated
              </li>
              <li>Build both installers (full + slim)</li>
              <li>Run the full installer on a test machine</li>
              <li>
                Confirm the Welcome wizard shows YOUR icon, not a default
              </li>
              <li>
                After install, confirm the Start menu shortcut and the
                running app both use your icon
              </li>
            </ol>
            <p>Paste a screenshot of the branded Welcome page in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            If you don&apos;t have a designed icon yet, use any 256 × 256 PNG
            and convert to ICO at icoconvert.com. The point of the
            exercise is the round-trip, not pixel-perfect art.
          </>
        }
        solution={
          <>
            <p>
              You should see your icon in the upper-right of every install
              page (header.bmp) and on the left of the Welcome + Finish
              pages (welcome.bmp). Plus matching icons in the Start menu
              shortcut and the running window.
            </p>
            <p>
              From a customer&apos;s perspective: this is now a real
              branded product, not a default NSIS template.
            </p>
          </>
        }
      />

      <h2>You finished Building Desktop with Go API 🎉</h2>
      <p>
        Five chapters, 14 lessons. You can now:
      </p>
      <ul>
        <li>Scaffold + run + build a Wails desktop app</li>
        <li>
          Offline-first with SQLite, the outbox pattern, and a sync
          engine
        </li>
        <li>Frameless window with custom titlebar and OS-aware controls</li>
        <li>
          In-app auto-updater that swaps the binary from GitHub releases
        </li>
        <li>Branded NSIS installers (full + slim) with custom bitmaps</li>
      </ul>
      <p>
        From here: try{' '}
        <a href="/courses/multiplatform" className="text-primary hover:underline">
          Building Web + Desktop + Mobile
        </a>{' '}
        to add a web + admin + mobile companion to the same Go API.
      </p>
    </>
  )
}
