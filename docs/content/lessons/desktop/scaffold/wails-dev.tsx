import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        <code>wails dev</code> is Wails&apos;s hot-reload loop — Vite for the
        frontend, recompiles + reloads Go on save. ~1 second iteration for
        small changes. This lesson covers the loop + the dev tools.
      </p>

      <h2>Start it</h2>
      <CodeBlock terminal code={`cd field-pos
wails dev`} />
      <p>
        First run takes ~30s — installs npm deps, compiles Go, opens a
        window. Save a React file, the window reloads. Save{' '}
        <code>app.go</code>, Wails rebuilds Go and re-launches.
      </p>

      <h2>The two-way bridge</h2>
      <Diagram label="Wails Go-to-React binding" caption="One bind call exposes a Go method to the React side as a typed Promise. The wailsjs/ folder is auto-generated at build.">
{`              ┌──────────────────────┐
              │   main.go            │
              │   wails.Run({        │
              │     Bind: []{        │
              │        app,          │
              │     }                │
              │   })                 │
              └──────────┬───────────┘
                         │
                         │  generates
                         ▼
              ┌──────────────────────┐
              │ frontend/wailsjs/    │
              │   go/main/App.d.ts   │
              │   go/main/App.js     │
              │                      │
              │   App.GetUser(id)    │
              │     returns User     │
              └──────────┬───────────┘
                         │
                         │  imported by
                         ▼
              ┌──────────────────────┐
              │  src/SomeComp.tsx    │
              │                      │
              │  GetUser(id)         │
              │    .then(user => …)  │
              └──────────────────────┘`}
      </Diagram>

      <p>
        Whatever struct + method you expose in Go shows up in TypeScript
        with full types. No REST plumbing. No JSON-encoding boilerplate.
      </p>

      <h2>Exposing your first Go method</h2>
      <CodeBlock
        language="go"
        filename="app.go"
        code={`type App struct {
    ctx context.Context
    db  *gorm.DB
}

// Exported = bound to React. Method name in TS = same.
func (a *App) Greet(name string) string {
    return fmt.Sprintf("Hello %s from Go!", name)
}`}
      />
      <p>And in React:</p>
      <CodeBlock
        language="tsx"
        filename="src/App.tsx"
        code={`import { Greet } from '../wailsjs/go/main/App'

function App() {
  const [greeting, setGreeting] = useState('')
  return (
    <button onClick={async () => setGreeting(await Greet('Alex'))}>
      {greeting || 'Click me'}
    </button>
  )
}`}
      />
      <p>
        Click the button — React calls <code>Greet</code>, Go runs, the
        return value lands in <code>greeting</code>. Zero plumbing.
      </p>

      <h2>The dev tools</h2>
      <ul>
        <li>
          <strong>Right-click in the app window → Inspect.</strong> Opens
          Chrome DevTools on the embedded WebView. Network tab, console,
          React DevTools — all the usual things.
        </li>
        <li>
          <strong>Wails terminal output.</strong> Anything you{' '}
          <code>log.Println</code> in Go shows here. Add prints liberally
          while developing.
        </li>
        <li>
          <strong>Hot reload paths.</strong> React: instant. Go: ~1-2s
          rebuild + window relaunch. Hot reload doesn&apos;t preserve state
          across Go changes — you re-login each time.
        </li>
      </ul>

      <TipBox tone="warning">
        <strong>State loss on Go change.</strong> Saving an{' '}
        <code>app.go</code> file rebuilds + relaunches the binary, so any
        in-memory React state resets. Persist what you need to
        SQLite or localStorage so you don&apos;t lose work.
      </TipBox>

      <h2>The first build is slow</h2>
      <p>
        First <code>wails dev</code> in a fresh project: ~30-60s. Subsequent
        starts: ~5s. The slow first start is npm install + Go module
        resolution. Be patient on day one.
      </p>

      <KnowledgeCheck
        question="You added a new method to your Go App struct. You save the file. wails dev rebuilds. You try to call it from React but TypeScript says it doesn't exist. What's the cause?"
        choices={[
          {
            label: 'TS server cache — restart your editor',
            feedback:
              "Sometimes helps, but Wails auto-regenerates wailsjs/. If the file isn't there, it's not a cache problem; it's a binding problem.",
          },
          {
            label: 'The method isn\'t exported (lowercase first letter) — Go won\'t bind it',
            correct: true,
            feedback:
              "Right — Wails only binds exported methods. `func (a *App) getUser(...)` is invisible to React. Capitalise: `GetUser`. Same rule as Go's package exports.",
          },
          {
            label: 'You need to manually run `wails generate`',
            feedback:
              "Wails auto-regenerates on save. Manual generation is for edge cases.",
          },
          {
            label: 'The method returns an unsupported type',
            feedback:
              "Possible — Wails has type-compatibility rules — but the export rule is the more common issue.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Build the call-Go-from-React loop end-to-end:</p>
            <ol>
              <li>
                Add a <code>Greet(name string) string</code> method to{' '}
                <code>app.go</code>.
              </li>
              <li>
                In <code>src/App.tsx</code>, import it and wire it to a
                button.
              </li>
              <li>Run <code>wails dev</code>.</li>
              <li>Click the button — confirm the greeting appears.</li>
            </ol>
            <p>Paste a screenshot of the working window in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            If TypeScript says <code>Greet doesn&apos;t exist</code>, check
            that the method name starts with a capital letter. Lowercase =
            unexported = unbound.
          </>
        }
        solution={
          <>
            <p>You should see &quot;Hello Alex from Go!&quot; after clicking.</p>
            <p>
              That round-trip — React → Wails → Go → response → React — is
              the entire interaction model. Every subsequent feature is a
              new bound method.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Last lesson of the chapter — <strong>wails build</strong>. Produce
        a real .exe / .app / .deb you can run without the CLI.
      </p>
    </>
  )
}
