import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Last lesson of the chapter — produce a real binary you can run on
        any machine. Two commands cover it: <code>wails build</code> for
        the bare .exe / .app, <code>wails build -nsis</code> for a Windows
        installer.
      </p>

      <h2>The basic build</h2>
      <CodeBlock terminal code={`wails build`} />
      <p>Produces a single binary for your current OS / arch in <code>build/bin/</code>:</p>
      <ul>
        <li>
          <strong>Windows:</strong>{' '}
          <code>build/bin/field-pos.exe</code> (~15 MB)
        </li>
        <li>
          <strong>macOS:</strong>{' '}
          <code>build/bin/field-pos.app</code> (a bundle)
        </li>
        <li>
          <strong>Linux:</strong>{' '}
          <code>build/bin/field-pos</code> (ELF binary)
        </li>
      </ul>

      <h2>Cross-compiling for Windows from Mac/Linux</h2>
      <CodeBlock
        terminal
        code={`wails build -platform windows/amd64`}
      />
      <p>
        Wails uses Go&apos;s cross-compile. The frontend bundle is the same;
        only the Go binary differs. Result lands in{' '}
        <code>build/bin/field-pos.exe</code>.
      </p>

      <h2>Build with NSIS installer</h2>
      <CodeBlock terminal code={`wails build -nsis`} />
      <p>Produces:</p>
      <ul>
        <li><code>field-pos.exe</code> — raw binary (~15 MB)</li>
        <li>
          <code>field-pos-amd64-installer.exe</code> — NSIS installer (~22 MB)
        </li>
      </ul>
      <p>
        The installer is the customer-facing artifact. It installs to{' '}
        <code>%LOCALAPPDATA%\Programs\field-pos\</code> per-user — no UAC,
        no admin needed.
      </p>

      <TipBox tone="info">
        <strong>NSIS must be installed:</strong>{' '}
        <code>winget install NSIS.NSIS</code> on Windows. macOS/Linux can
        cross-build to Windows .exe but can&apos;t produce the NSIS installer
        unless NSIS is on the box (or via WINE).
      </TipBox>

      <h2>Build flags worth knowing</h2>
      <CodeBlock
        terminal
        code={`# Smaller binary — strips debug symbols
wails build -trimpath -ldflags "-s -w"

# Custom output name
wails build -o "FieldPOS-v1.0.0.exe"

# All platforms in one shot
wails build -platform windows/amd64,darwin/arm64,linux/amd64`}
      />
      <p>
        <code>-ldflags &quot;-s -w&quot;</code> can shave 30% off the binary
        size. Worth it for distribution.
      </p>

      <h2>Where the assets come from</h2>
      <p>
        Your app icon, splash image, and metadata come from{' '}
        <code>build/windows/</code> (and parallel folders for{' '}
        <code>darwin/</code>, <code>linux/</code>). The Grit scaffold ships
        placeholders. Swap them for your brand assets before shipping.
      </p>
      <ul>
        <li>
          <code>build/windows/icon.ico</code> — app icon (multi-resolution
          .ico)
        </li>
        <li>
          <code>build/windows/info.json</code> — version + company metadata
        </li>
        <li>
          <code>build/darwin/icon.icns</code> — macOS bundle icon
        </li>
      </ul>

      <KnowledgeCheck
        question="You ran `wails build -nsis` on Windows. Both `.exe` files appear in build/bin. Which one do you ship to customers?"
        choices={[
          {
            label: 'The raw .exe — smaller download',
            feedback:
              "Wrong for customers. The raw exe doesn't install anywhere, doesn't add Start menu entries, doesn't handle upgrades or uninstalls. Customers need the installer.",
          },
          {
            label: 'The Setup .exe (the installer) — handles install, Start menu, uninstall',
            correct: true,
            feedback:
              "Right — the installer is what customers expect. The raw .exe is for the auto-updater (chapter 4) which swaps just the binary, not the install bits.",
          },
          {
            label: 'Both, in a zip',
            feedback:
              "Possible but confusing. Customers don\'t need the raw exe; that\'s for the auto-updater\'s use only.",
          },
          {
            label: 'Neither — code-sign first, then ship',
            feedback:
              "Code-signing is recommended (skips SmartScreen) but optional. The installer is still the right artifact regardless of signing state.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>For chapter 1&apos;s assignment, build + install + uninstall:</p>
            <ol>
              <li>Run <code>wails build -nsis</code> in your project</li>
              <li>
                Find <code>build/bin/field-pos-amd64-installer.exe</code>
              </li>
              <li>Double-click → run the installer</li>
              <li>Launch the app from the Start menu</li>
              <li>
                Uninstall via{' '}
                <em>Settings → Apps → field-pos → Uninstall</em>
              </li>
            </ol>
            <p>Paste screenshots of the installer wizard + the uninstall confirmation in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            If the installer fails with &quot;NSIS not installed&quot;, run{' '}
            <code>winget install NSIS.NSIS</code> first. Re-run the build.
          </>
        }
        solution={
          <>
            <p>
              The installer wizard should look standard: Welcome page →
              Directory page → Install page → Finish page (with &quot;Launch
              field-pos&quot; checkbox). Uninstall removes the app cleanly.
            </p>
            <p>
              That&apos;s chapter 1 done. Your scaffolded desktop app is
              now a real, installable Windows binary.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 2 — <strong>Offline-first</strong>. Local SQLite, the
        outbox pattern for queued writes, and the sync engine that pushes
        changes when you reconnect.
      </p>
    </>
  )
}
