import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        A frameless Wails window has no OS titlebar — you draw your own.
        Pros: it looks like Linear / Discord / VS Code, not Win32. Cons:
        you have to wire drag-to-move yourself. This lesson covers both.
      </p>

      <h2>Enable frameless in main.go</h2>
      <CodeBlock
        language="go"
        filename="main.go (excerpt)"
        code={`wails.Run(&options.App{
    Title:      "Field POS",
    Width:      1280,
    Height:     800,
    Frameless:  true,    // ← removes OS titlebar
    CSSDragProperty: "--wails-draggable",
    CSSDragValue:    "drag",
    // ...
})`}
      />
      <p>
        Frameless + the CSS drag declarations let you mark which HTML
        elements can drag the window. Anything with{' '}
        <code>style=&quot;--wails-draggable: drag&quot;</code> becomes a drag handle.
      </p>

      <h2>The titlebar component</h2>
      <CodeBlock
        language="tsx"
        filename="src/components/titlebar.tsx"
        code={`import { Quit, WindowMinimise, WindowToggleMaximise } from '../../wailsjs/runtime/runtime'
import { Minus, Square, X } from 'lucide-react'

export function TitleBar() {
  return (
    <div
      style={{ '--wails-draggable': 'drag' } as any}
      className="h-10 flex items-center justify-between px-4 bg-card border-b border-border select-none"
    >
      {/* Left: app icon + title */}
      <div className="flex items-center gap-2">
        <img src="/icon.svg" className="h-5 w-5" alt="" />
        <span className="text-sm font-medium">Field POS</span>
      </div>

      {/* Right: window controls — NOT draggable */}
      <div
        style={{ '--wails-draggable': 'no-drag' } as any}
        className="flex items-center"
      >
        <WinBtn onClick={WindowMinimise}><Minus className="h-4 w-4" /></WinBtn>
        <WinBtn onClick={WindowToggleMaximise}><Square className="h-3 w-3" /></WinBtn>
        <WinBtn onClick={Quit} danger><X className="h-4 w-4" /></WinBtn>
      </div>
    </div>
  )
}

function WinBtn({ onClick, danger, children }) {
  return (
    <button
      onClick={onClick}
      className={\`h-10 w-12 flex items-center justify-center hover:bg-muted \${danger ? 'hover:bg-red-500 hover:text-white' : ''}\`}
    >
      {children}
    </button>
  )
}`}
      />
      <p>Two key bits:</p>
      <ul>
        <li>
          The outer div has <code>--wails-draggable: drag</code> — the
          whole titlebar can drag the window.
        </li>
        <li>
          The window-controls cluster has{' '}
          <code>--wails-draggable: no-drag</code> — buttons are clickable,
          not drag handles.
        </li>
      </ul>

      <TipBox tone="warning">
        <strong>Always mark buttons as no-drag.</strong> Without it,
        clicking the X starts a drag-then-click, the user&apos;s click
        registers on whatever happens to be under the cursor after the
        drag, and the X doesn&apos;t fire. Annoying for the user; subtle for
        you.
      </TipBox>

      <h2>The Wails runtime API</h2>
      <ul>
        <li><code>WindowMinimise()</code> — minimise to taskbar</li>
        <li><code>WindowToggleMaximise()</code> — toggle full vs. windowed</li>
        <li><code>WindowFullscreen()</code> — true OS fullscreen (different from maximise)</li>
        <li><code>Quit()</code> — exit the app</li>
        <li><code>WindowSetTitle(s)</code> — change the title at runtime</li>
        <li><code>WindowCenter()</code> — recenter on screen</li>
      </ul>
      <p>
        All exported from <code>../wailsjs/runtime/runtime</code> and ready
        to import.
      </p>

      <h2>Why frameless is the right default for line-of-business apps</h2>
      <ul>
        <li>
          <strong>Brand consistency.</strong> Your app looks the same on
          every OS. No mismatched native-vs-web fonts in the titlebar.
        </li>
        <li>
          <strong>More vertical pixels.</strong> No 30px OS titlebar
          eating space. For dashboards, every pixel counts.
        </li>
        <li>
          <strong>Custom controls.</strong> Add a tab bar, search box, or
          status indicator IN the titlebar. Common pattern in
          spreadsheet-shaped apps.
        </li>
      </ul>

      <KnowledgeCheck
        question="A user reports 'the close button doesn't work when I click it'. You inspect — `Quit` is wired correctly. What's the bug?"
        choices={[
          {
            label: 'The window controls container isn\'t marked `--wails-draggable: no-drag`',
            correct: true,
            feedback:
              "Right — without the no-drag marker, every click starts a drag gesture. The mouseup happens elsewhere and the click never registers on the button. Classic frameless-window gotcha.",
          },
          {
            label: 'The Quit import is wrong',
            feedback:
              "Wrong — Quit DOES work; the click just isn\'t reaching it.",
          },
          {
            label: 'Wails 2 doesn\'t support Quit on frameless windows',
            feedback:
              "It does — Quit works regardless of frameless or not.",
          },
          {
            label: 'The button needs a higher z-index',
            feedback:
              "Wrong — z-index doesn\'t affect drag-vs-click. The CSS drag declaration does.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Build a working titlebar:</p>
            <ol>
              <li>
                Set <code>Frameless: true</code> + the CSS drag declarations
                in <code>main.go</code>.
              </li>
              <li>
                Add the <code>TitleBar</code> component above your main
                layout.
              </li>
              <li>
                Drag the window by the empty area — confirm it moves.
              </li>
              <li>
                Click each window button — Minimise, Maximise, Close — all
                should work.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            If you can&apos;t move the window, double-check{' '}
            <code>main.go</code> sets both <code>CSSDragProperty</code> and{' '}
            <code>CSSDragValue</code>. They&apos;re both required.
          </>
        }
        solution={
          <>
            <p>
              You should have a draggable titlebar with three working
              window buttons. That&apos;s ~40 lines of React + 3 Go config
              lines for what feels like a native Win32 window.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Last lesson of the chapter — going deeper on the window controls:
        traffic-light styles on Mac, double-click-to-maximize, and the
        edge cases.
      </p>
    </>
  )
}
