import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        The Grit scaffold ships an updater that polls GitHub for new
        releases, downloads the new .exe to <code>%TEMP%</code>, then swaps
        it with the running binary on user confirmation. This lesson
        covers the Go side; next lesson covers the React modal.
      </p>

      <h2>The swap mechanism — the Windows .exe rename trick</h2>
      <Diagram label="Binary-swap sequence (Windows)" caption="You can't overwrite a running .exe — but you CAN rename it. Windows keeps the kernel object alive while the directory entry moves.">
{`   ┌─────────────────────────────────────────────────┐
   │                                                 │
   │  1.  Download new v1.2.1.exe to %TEMP%          │
   │                                                 │
   │      C:/Users/.../%TEMP%/FieldPOS-v1.2.1.exe    │
   │                                                 │
   │  2.  Rename running .exe out of the way         │
   │                                                 │
   │      FieldPOS.exe          ──►   FieldPOS.exe.old │
   │                                                 │
   │      (Windows allows this even while running)   │
   │                                                 │
   │  3.  Move new .exe into place                   │
   │                                                 │
   │      %TEMP%/FieldPOS-v1.2.1.exe                 │
   │                  ──►   FieldPOS.exe             │
   │                                                 │
   │  4.  Spawn the new binary, exit the old         │
   │                                                 │
   │      cmd.Start(FieldPOS.exe)                    │
   │      runtime.Quit(ctx)                          │
   │                                                 │
   │  5.  Next startup: delete FieldPOS.exe.old      │
   │                                                 │
   │      go updater.CleanupOldOnStartup()           │
   │                                                 │
   └─────────────────────────────────────────────────┘`}
      </Diagram>

      <h2>The Updater struct</h2>
      <CodeBlock
        language="go"
        filename="updater.go (scaffolded)"
        code={`type Updater struct {
    ctx           context.Context
    GithubOwner   string
    GithubRepo    string

    mu            sync.RWMutex
    status        string  // idle | downloading | downloaded | installing | error
    target        string  // path to the downloaded .exe
    version       string
    errMsg        string

    bytesDownloaded atomic.Int64
    bytesTotal      atomic.Int64
}

// Bound to React — kicks off the download
func (u *Updater) StartDownload(url, version string) error { ... }

// Polled by React for the progress bar
func (u *Updater) GetProgress() UpdaterState { ... }

// Final user action — swap + relaunch
func (u *Updater) InstallAndRestart() error { ... }`}
      />

      <h2>The check + download</h2>
      <CodeBlock
        language="go"
        code={`func (u *Updater) CheckForUpdates() (*UpdateRelease, error) {
    res, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest",
        u.GithubOwner, u.GithubRepo))
    // … decode JSON, extract version, URLs, release notes
    asset := pickAsset(release.Assets, runtime.GOOS, runtime.GOARCH)
    return &UpdateRelease{
        Version:      release.TagName,
        DownloadURL:  asset.URL,
        ReleaseNotes: release.Body,
    }, nil
}`}
      />
      <p>
        Polls every 6 hours, caches the result. On a hit, the React side
        renders the &quot;update available&quot; banner.
      </p>

      <h2>The swap call</h2>
      <CodeBlock
        language="go"
        code={`func (u *Updater) InstallAndRestart() error {
    currentExe, _ := os.Executable()
    oldPath := currentExe + ".old"
    _ = os.Remove(oldPath)

    // Move the running .exe aside
    if err := os.Rename(currentExe, oldPath); err != nil {
        return err
    }

    // Move the downloaded .exe into place
    if err := copyFile(u.target, currentExe); err != nil {
        _ = os.Rename(oldPath, currentExe)  // rollback
        return err
    }
    _ = os.Remove(u.target)

    // Launch the new version + exit the old
    exec.Command(currentExe).Start()
    go func() {
        time.Sleep(250 * time.Millisecond)
        runtime.Quit(u.ctx)
    }()
    return nil
}`}
      />

      <TipBox tone="warning">
        <strong>Why copy, not rename, for the install step?</strong>{' '}
        <code>%TEMP%</code> is often on a different drive (C:) than{' '}
        <code>%LOCALAPPDATA%</code>. <code>os.Rename</code> across drives
        fails with EXDEV. <code>copyFile</code> works regardless. The
        rename trick from step 2 works because both paths are on the
        same drive.
      </TipBox>

      <h2>Cleanup on startup</h2>
      <CodeBlock
        language="go"
        filename="main.go"
        code={`func main() {
    app := NewApp()
    updater := NewUpdater()
    go updater.CleanupOldOnStartup()  // delete leftover .exe.old, retry 25x

    wails.Run(...)
}`}
      />
      <p>
        Background goroutine, fires once at startup. If the old .exe is
        still locked (antivirus scanning it), retry 25 times at 200ms
        intervals. Best-effort.
      </p>

      <h2>POSIX (macOS/Linux) is simpler</h2>
      <p>
        On Mac/Linux, you CAN overwrite a running binary directly. The
        OS keeps the running process&apos;s inode alive while the file at
        the path is replaced. The Grit scaffold branches on OS:
      </p>
      <CodeBlock
        language="go"
        code={`if runtime.GOOS == "windows" {
    u.swapWindows(currentExe, src)
} else {
    copyFile(src, currentExe)  // straight overwrite — POSIX magic
}`}
      />

      <KnowledgeCheck
        question="You're testing the auto-updater. Step 2 (rename running .exe) succeeds, but step 3 (copy new .exe into place) fails with 'disk full'. What's the safest behaviour?"
        choices={[
          {
            label: 'Rename the .old back to the original name — leave the user with a working old binary',
            correct: true,
            feedback:
              "Right — without rollback, the user has no .exe at all. The rollback line in the scaffolded updater handles this exactly. Better to stay on v1 than have a broken install.",
          },
          {
            label: 'Show an error and exit',
            feedback:
              "Worse than a rollback — the user is stuck with no app until they manually rename .old. Always roll back automatic state changes on failure.",
          },
          {
            label: 'Re-attempt the copy 5 times',
            feedback:
              "Disk full won\'t magically resolve in 5 retries. Better to roll back and tell the user to free space, then retry the update.",
          },
          {
            label: 'Truncate the .old to free space',
            feedback:
              "Truncates the user\'s working binary. Catastrophic.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Inspect the scaffolded <code>updater.go</code> in your
              desktop project:
            </p>
            <ol>
              <li>
                Find the <code>swapWindows</code> function. List the 4
                steps it takes.
              </li>
              <li>
                Find the rollback line — what conditions trigger it?
              </li>
              <li>
                In <code>notes.md</code>, write one paragraph explaining
                why <code>os.Rename(currentExe, oldPath)</code> works
                even though the .exe is running.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            The trick is that Windows separates the file&apos;s data (kernel
            object) from its directory entry. The kernel object stays
            alive for the running process; the directory entry can move.
          </>
        }
        solution={
          <>
            <p>The four steps inside <code>swapWindows</code>:</p>
            <ol>
              <li>Remove any stale <code>.exe.old</code> from a prior run</li>
              <li>Rename the running .exe to <code>.exe.old</code></li>
              <li>Copy the new .exe into place</li>
              <li>If copy failed, rename .old back (rollback)</li>
            </ol>
            <p>
              Windows allows the rename because the kernel doesn&apos;t care
              what path leads to the running file; it only cares the
              inode is alive. The path can move freely.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Updater logic done. Next — the <strong>React modal + banner</strong>{' '}
        that drive it: progress bar, install button, release notes.
      </p>
    </>
  )
}
