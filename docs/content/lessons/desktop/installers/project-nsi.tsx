import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        The full installer (<code>project.nsi</code>) bundles the
        WebView2 runtime — ~125 MB Microsoft Edge component your app needs
        to render HTML. Total installer size: ~150 MB. Works offline.
      </p>

      <h2>Why bundle WebView2?</h2>
      <ul>
        <li>
          Windows 11 ships it. Older Windows installs do NOT have it.
        </li>
        <li>
          Without WebView2, the app launches and shows a blank window —
          worst kind of failure for a non-technical user.
        </li>
        <li>
          Bundled WebView2 = works on a fresh machine with no internet
          during install.
        </li>
        <li>
          If the customer is in a high-cost-data region, downloading 125
          MB during install is painful — bundling means it&apos;s already
          there.
        </li>
      </ul>

      <h2>The NSI structure</h2>
      <CodeBlock
        language="nsi"
        filename="build/windows/installer/project.nsi (key parts)"
        code={`Unicode true

!define PRODUCT_EXECUTABLE "FieldPOS.exe"
!define UNINST_KEY_NAME    "FieldPOS"

!include "wails_tools.nsh"   ; Wails' generated paths + macros
!include "MUI.nsh"           ; Modern UI

# Branded MUI
!define MUI_HEADERIMAGE
!define MUI_HEADERIMAGE_BITMAP   "header.bmp"
!define MUI_WELCOMEFINISHPAGE_BITMAP "welcome.bmp"

# Pages: Welcome → Directory → Install → Finish
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

# Per-user install — no UAC
RequestExecutionLevel user
InstallDir "$LOCALAPPDATA\\Programs\\FieldPOS"

Section
    SetOutPath "$INSTDIR"

    # Bundle WebView2
    File "MicrosoftEdgeWebView2RuntimeInstallerX64.exe"
    ExecWait '"$INSTDIR\\MicrosoftEdgeWebView2RuntimeInstallerX64.exe" /silent /install'

    # Main binary
    File "/oname=FieldPOS.exe" "..\\..\\bin\\FieldPOS-v\${INFO_PRODUCTVERSION}.exe"

    # Shortcuts
    CreateShortcut "$SMPROGRAMS\\FieldPOS.lnk" "$INSTDIR\\FieldPOS.exe"
    CreateShortcut "$DESKTOP\\FieldPOS.lnk" "$INSTDIR\\FieldPOS.exe"

    # Register uninstaller
    WriteUninstaller "$INSTDIR\\Uninstall.exe"
SectionEnd`}
      />

      <h2>Why per-user install (no UAC)</h2>
      <p>
        Default install location: <code>%LOCALAPPDATA%\Programs\FieldPOS\</code>.
        The user&apos;s own AppData — no admin needed.
      </p>
      <p>The trade-offs:</p>
      <ul>
        <li>
          <strong>Pro:</strong> no UAC prompt; the auto-updater can swap
          the .exe without elevation
        </li>
        <li>
          <strong>Pro:</strong> reverse-uninstalled per user when they log
          off — clean for shared workstations
        </li>
        <li>
          <strong>Con:</strong> if 3 users on one PC all install it,
          there are 3 copies of the .exe (~45 MB total)
        </li>
      </ul>
      <p>
        For line-of-business apps with 1 user per PC, per-user wins
        every time.
      </p>

      <TipBox tone="success">
        <strong>Same install path as Slack / Discord / VS Code /
        GitHub Desktop.</strong> If users have other modern apps, the
        location feels familiar.
      </TipBox>

      <h2>WebView2 detection — skip if already installed</h2>
      <CodeBlock
        language="nsi"
        code={`# Check 3 registry locations; skip the install if WebView2 is already present
ReadRegStr $0 HKLM "SOFTWARE\\Wow6432Node\\Microsoft\\EdgeUpdate\\Clients\\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}" "pv"
StrCmp $0 "" 0 webview2_installed

ReadRegStr $0 HKLM "SOFTWARE\\Microsoft\\EdgeUpdate\\Clients\\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}" "pv"
StrCmp $0 "" 0 webview2_installed

# … etc
ExecWait '"$INSTDIR\\MicrosoftEdgeWebView2RuntimeInstallerX64.exe" /silent /install'

webview2_installed:
  # Skip — already there`}
      />

      <h2>Single-instance lock</h2>
      <p>
        Prevent two installers running simultaneously:
      </p>
      <CodeBlock
        language="nsi"
        code={`# At the top of every section
System::Call 'kernel32::CreateMutex(p 0, i 0, t "FieldPOS_Installer") p .r0'
System::Call 'kernel32::GetLastError() i .r1'
IntCmp $1 183 mutex_busy

# 183 = ERROR_ALREADY_EXISTS
mutex_busy:
  MessageBox MB_OK|MB_ICONEXCLAMATION "Setup is already running."
  Abort`}
      />

      <KnowledgeCheck
        question="A customer with no internet runs your full installer on a fresh Windows 10 machine. Does it work?"
        choices={[
          {
            label: 'No — installers always need internet',
            feedback:
              "Wrong — that\'s the whole point of bundling WebView2. The full installer ships every dependency.",
          },
          {
            label: 'Yes — WebView2 runtime is bundled inside the installer, so install + first launch work offline',
            correct: true,
            feedback:
              "Right — that\'s exactly why the full installer is ~150 MB. Customer can install on an air-gapped box. The slim installer would fail in this case.",
          },
          {
            label: 'Only if Windows Update is configured',
            feedback:
              "Wrong — Windows Update isn\'t needed when WebView2 is bundled.",
          },
          {
            label: 'Yes — Wails apps don\'t use WebView2',
            feedback:
              "They do. WebView2 is the rendering engine. Without it, the window stays blank.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              On a test PC (or a Hyper-V VM with no internet):
            </p>
            <ol>
              <li>
                Run the full installer (<code>FieldPOS-Setup-v0.1.0.exe</code>)
              </li>
              <li>Walk through the wizard</li>
              <li>Launch the app from the Start menu</li>
              <li>
                Confirm it renders (not a blank window) — proves WebView2
                got installed
              </li>
              <li>Uninstall via Settings → Apps</li>
            </ol>
            <p>
              Paste a screenshot of the Welcome wizard page in{' '}
              <code>notes.md</code>.
            </p>
          </>
        }
        hint={
          <>
            If you don&apos;t have a Hyper-V VM, just disable the network
            interface in Network Settings before launching the installer.
            Same effect.
          </>
        }
        solution={
          <>
            <p>
              You should see the full installer flow: branded Welcome →
              Directory (offers a default with no UAC) → Install (with
              WebView2 install progress visible) → Finish (with &quot;Launch
              FieldPOS&quot; checkbox).
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next — the <strong>slim installer</strong>. ~22 MB by downloading
        WebView2 at install time. Right tradeoff when bandwidth matters.
      </p>
    </>
  )
}
