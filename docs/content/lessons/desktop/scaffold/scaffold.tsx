import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        Welcome to <strong>Building Desktop with Go API</strong>. Different
        kit, different command, different binary. This lesson covers the
        scaffold + what makes it different from the web/mobile kits.
      </p>

      <h2>The command — note it&apos;s different</h2>
      <CodeBlock terminal code={`grit new-desktop field-pos`} />
      <p>
        Yes — <code>new-desktop</code>, not <code>new --desktop</code>. The
        desktop kit produces a <strong>standalone Wails app</strong> (Go +
        React bundled into one binary), not a monorepo. That&apos;s why it has
        its own command.
      </p>

      <h2>What you get</h2>
      <Diagram label="Desktop project shape" caption="One binary. The frontend is bundled INTO the Go executable at build time.">
{`   ┌─────────────────────────────────────────────┐
   │     field-pos/  (one Wails project)         │
   ├─────────────────────────────────────────────┤
   │                                             │
   │   main.go              ◄─── entry point     │
   │   app.go               ◄─── Go methods      │
   │                            bound to React   │
   │   internal/                                 │
   │     ├── db/             local SQLite        │
   │     ├── models/         GORM structs        │
   │     ├── service/        business logic      │
   │     └── ...                                 │
   │                                             │
   │   frontend/             ◄─── React + Vite   │
   │     ├── src/                                │
   │     ├── package.json                        │
   │     └── wailsjs/        auto-generated Go-  │
   │                          to-TS bridge       │
   │                                             │
   │   build/                Wails build assets  │
   │   wails.json            Wails config        │
   │   go.mod                                    │
   └─────────────────────────────────────────────┘

           wails build  ──────────────►

   ┌─────────────────────────────────┐
   │       field-pos.exe             │  ◄── single binary,
   │   (frontend dist embedded)      │      no server,
   └─────────────────────────────────┘      no internet.`}
      </Diagram>

      <h2>What ships</h2>
      <ul>
        <li>
          <strong>Wails v2</strong> binding Go methods directly to React
          (no IPC plumbing)
        </li>
        <li>
          <strong>SQLite + GORM</strong> for local storage — survives
          across launches
        </li>
        <li>
          <strong>Frameless window</strong> with a custom titlebar (chapter
          3 covers the polish)
        </li>
        <li>
          <strong>Local auth with bcrypt</strong> — no server needed
        </li>
        <li>
          <strong>In-app auto-updater</strong> — binary swap from GitHub
          releases (chapter 4)
        </li>
        <li>
          <strong>NSIS installer templates</strong> (chapter 5)
        </li>
        <li>
          <strong>Tailwind + shadcn/ui</strong> on the React side
        </li>
      </ul>

      <h2>Prerequisites</h2>
      <ul>
        <li>Go 1.21+ installed</li>
        <li>Node 18+ and pnpm</li>
        <li>
          <strong>Wails v2 CLI</strong> —{' '}
          <code>go install github.com/wailsapp/wails/v2/cmd/wails@latest</code>
        </li>
        <li>
          <strong>Windows:</strong> WebView2 (ships with Windows 11; older
          Windows installs it on first build)
        </li>
        <li>
          <strong>macOS:</strong> Xcode CLI tools{' '}
          (<code>xcode-select --install</code>)
        </li>
        <li>
          <strong>Linux:</strong> WebKit2GTK + libgtk-3-dev (Wails docs lists
          the apt commands)
        </li>
      </ul>

      <TipBox tone="info">
        Confirm Wails CLI is healthy with{' '}
        <code>wails doctor</code> — it checks every dependency and tells
        you exactly what&apos;s missing.
      </TipBox>

      <KnowledgeCheck
        question="A teammate asks why we use `grit new-desktop` and not `grit new --desktop`. What's the cleanest answer?"
        choices={[
          {
            label: 'Historical reasons',
            feedback:
              'Wrong — there\'s a real structural reason.',
          },
          {
            label: 'The desktop kit produces a standalone Wails app, not a monorepo — different shape entirely',
            correct: true,
            feedback:
              "Right — the other kits all produce monorepos with apps/. Desktop is one project, structured for Wails. That difference earns its own command. `grit new --triple --desktop` IS a thing (covered in the Multi-Platform course) — that adds a Wails app to an existing triple monorepo.",
          },
          {
            label: 'Performance — new-desktop uses different scaffolder code',
            feedback:
              'Internal detail, not the user-facing reason.',
          },
          {
            label: 'Wails requires its own CLI tool',
            feedback:
              'Wails does — but that\'s a runtime requirement, not why the Grit command differs.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Scaffold + verify:</p>
            <ol>
              <li>
                Install Wails CLI:{' '}
                <code>go install github.com/wailsapp/wails/v2/cmd/wails@latest</code>
              </li>
              <li>Run <code>wails doctor</code> — fix anything red.</li>
              <li><code>grit new-desktop field-pos</code></li>
              <li>
                <code>cd field-pos && ls</code> — confirm{' '}
                <code>main.go</code>, <code>frontend/</code>, and{' '}
                <code>wails.json</code> are at the root (not under{' '}
                <code>apps/</code>).
              </li>
            </ol>
            <p>Paste the directory listing in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            If <code>wails doctor</code> complains on Windows about
            WebView2, install it from{' '}
            <a href="https://developer.microsoft.com/microsoft-edge/webview2/" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">developer.microsoft.com/microsoft-edge/webview2</a>.
          </>
        }
        solution={
          <>
            <p>You should see (one-level deep):</p>
            <CodeBlock
              language="text"
              code={`build/  frontend/  internal/  main.go  app.go  wails.json  go.mod  README.md`}
            />
            <p>
              No <code>apps/</code> folder. That&apos;s the structural
              difference from every other Grit kit.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — <code>wails dev</code>. Hot-reload loop for both Go
        AND React. Much faster than building every change.
      </p>
    </>
  )
}
