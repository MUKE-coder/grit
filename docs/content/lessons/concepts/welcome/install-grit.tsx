import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Time to install Grit. Pick one path — the install script (everyone)
        or <code>go install</code> (if you have the Go toolchain). Both
        produce the exact same binary at the same place on your PATH.
      </p>

      <h2>Path A — the install script (recommended)</h2>
      <p>
        One line. Works on macOS, Linux, and Windows. Detects whether you
        already have grit and either installs fresh or runs the self-update.
      </p>

      <h3>macOS / Linux</h3>
      <CodeBlock
        terminal
        code={`curl -fsSL https://gritframework.dev/install.sh | sh`}
      />

      <h3>Windows (PowerShell)</h3>
      <CodeBlock
        terminal
        code={`iwr -useb https://gritframework.dev/install.ps1 | iex`}
      />

      <p>
        The script does three things: detects your OS + architecture,
        downloads the matching release binary from GitHub, and puts it on
        your PATH (<code>/usr/local/bin/grit</code> on Unix,{' '}
        <code>%USERPROFILE%\.grit\bin\grit.exe</code> on Windows). If you
        already had grit, it runs <code>grit update</code> instead.
      </p>

      <TipBox tone="info">
        On Windows the install script appends <code>%USERPROFILE%\.grit\bin</code>{' '}
        to your user PATH. You&apos;ll need to <strong>open a new terminal</strong>{' '}
        for the change to take effect — the script tells you that, but
        it&apos;s easy to miss.
      </TipBox>

      <h2>Path B — `go install` (if you have Go)</h2>
      <p>
        If you already have Go 1.21+ on your machine, this is the most
        traditional path:
      </p>

      <CodeBlock
        terminal
        code={`go install github.com/MUKE-coder/grit/v3/cmd/grit@latest`}
      />

      <p>
        Compiles grit from source and puts it in{' '}
        <code>$(go env GOPATH)/bin/grit</code>. Make sure that directory is on
        your PATH — most Go setups add it automatically.
      </p>

      <h2>Want a pinned version?</h2>
      <p>
        For reproducible installs (CI scripts, team alignment), pin the exact
        version:
      </p>

      <CodeBlock
        terminal
        code={`# Pin via the install script
GRIT_VERSION=v3.26.5 curl -fsSL https://gritframework.dev/install.sh | sh

# Pin via go install
go install github.com/MUKE-coder/grit/v3/cmd/grit@v3.26.5`}
      />

      <TipBox tone="warning">
        <strong>Don&apos;t install grit globally via npm.</strong> There&apos;s no npm
        package. If you find one, it&apos;s not us — install via the script or
        <code> go install</code> only.
      </TipBox>

      <h2>What ships with the install</h2>
      <p>
        One binary (~20 MB). That&apos;s it. No background daemon, no telemetry,
        no licence server. You can move <code>grit</code> to a USB stick and
        run it anywhere with the matching OS.
      </p>

      <KnowledgeCheck
        question="You don't have Go installed and you want grit on a Windows machine. Which command works?"
        choices={[
          {
            label: 'go install github.com/MUKE-coder/grit/v3/cmd/grit@latest',
            feedback:
              "Won't work — go install needs the Go toolchain. The error will be 'go: command not found'.",
          },
          {
            label: 'npm install -g grit-framework',
            feedback:
              "Won't work — there's no npm package. Anything claiming to be Grit on npm isn't ours.",
          },
          {
            label: 'iwr -useb https://gritframework.dev/install.ps1 | iex',
            correct: true,
            feedback:
              "Right — the PowerShell install script downloads the release binary directly. No Go toolchain required.",
          },
          {
            label: 'brew install grit',
            feedback:
              "Won't work on Windows. Even on macOS there's no Homebrew formula yet — the install script is the canonical path.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Install grit on your machine using whichever path matches your
              setup. Then run <code>grit version</code> and paste the output
              into your <code>notes.md</code>.
            </p>
          </>
        }
        hint={
          <>
            If <code>grit version</code> says &quot;command not found&quot;, your PATH
            probably hasn&apos;t picked up the new directory yet. <strong>Open a
            new terminal</strong> first. On Windows specifically you may need
            to close all terminal windows and reopen one.
          </>
        }
        solution={
          <>
            <p>You should see something like:</p>
            <CodeBlock
              terminal
              code={`$ grit version
grit version 3.25.x`}
            />
            <p>
              If the version is <code>3.24.x</code> or older, run{' '}
              <code>grit update</code> to pick up the latest features.
            </p>
            <p>
              If <code>grit version</code> still doesn&apos;t work after opening a
              new terminal:
            </p>
            <ul>
              <li>
                <strong>Mac/Linux:</strong> run <code>which grit</code>. If
                nothing prints, the install dir isn&apos;t on PATH. Add{' '}
                <code>export PATH=&quot;$HOME/.local/bin:$PATH&quot;</code> to your{' '}
                <code>~/.zshrc</code> or <code>~/.bashrc</code> and re-source.
              </li>
              <li>
                <strong>Windows:</strong> run <code>where grit</code> in cmd
                or <code>(Get-Command grit).Source</code> in PowerShell. If
                nothing prints, manually add <code>%USERPROFILE%\.grit\bin</code>{' '}
                to your user PATH via System Properties → Environment
                Variables.
              </li>
            </ul>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        You have grit on your machine. Next lesson — a 3-minute sanity check
        — confirms everything&apos;s wired right and you know how to update
        when a new version drops.
      </p>
    </>
  )
}
