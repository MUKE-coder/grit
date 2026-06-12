import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Polishing the window controls — OS-conventional layouts, the
        unmaximize state, double-click-to-maximize, the edge cases users
        will report.
      </p>

      <h2>Windows vs. macOS controls — different conventions</h2>
      <ul>
        <li>
          <strong>Windows:</strong> Minimise / Maximise / Close on the
          RIGHT, in that order. Square boxes, ~12-14px icons.
        </li>
        <li>
          <strong>macOS:</strong> Close / Minimise / Zoom on the LEFT,
          colored traffic lights (red / yellow / green).
        </li>
      </ul>
      <p>
        Wails ships the native OS traffic lights on macOS if{' '}
        <code>options.Mac</code> isn&apos;t set to hide them. The simplest
        approach: use Wails&apos;s defaults on Mac, draw your own on Windows.
      </p>

      <h2>OS detection in React</h2>
      <CodeBlock
        language="ts"
        filename="src/lib/os.ts"
        code={`import { Environment } from '../../wailsjs/runtime/runtime'

let cachedOS: 'windows' | 'darwin' | 'linux' | null = null

export async function detectOS() {
  if (cachedOS) return cachedOS
  const env = await Environment()
  cachedOS = env.platform as any
  return cachedOS
}`}
      />
      <p>
        Now your titlebar can render OS-appropriate controls or hide them
        entirely on Mac (Wails handles the traffic lights).
      </p>

      <h2>Double-click to maximise</h2>
      <p>
        Standard convention — double-click the titlebar maximises /
        unmaximises. Wire it:
      </p>
      <CodeBlock
        language="tsx"
        code={`<div
  style={{ '--wails-draggable': 'drag' } as any}
  onDoubleClick={WindowToggleMaximise}
  className="h-10 flex items-center justify-between px-4 ..."
>
  ...`}
      />
      <p>
        One line. Now power users have the shortcut they expect.
      </p>

      <h2>Tracking maximised state for the icon swap</h2>
      <p>
        Convention: the maximise button shows a single square when the
        window is normal, two stacked squares (or a restore icon) when
        maximised. Listen to the Wails event:
      </p>
      <CodeBlock
        language="tsx"
        code={`import { EventsOn } from '../../wailsjs/runtime/runtime'

export function useIsMaximised() {
  const [max, setMax] = useState(false)
  useEffect(() => {
    const cleanup = EventsOn('wails:window:maximise', () => setMax(true))
    const cleanup2 = EventsOn('wails:window:unmaximise', () => setMax(false))
    return () => { cleanup(); cleanup2() }
  }, [])
  return max
}`}
      />
      <CodeBlock
        language="tsx"
        code={`const isMax = useIsMaximised()

<WinBtn onClick={WindowToggleMaximise}>
  {isMax ? <Restore /> : <Square />}
</WinBtn>`}
      />

      <TipBox tone="info">
        <strong>Quit vs. Close.</strong> Wails&apos; <code>Quit()</code> ends
        the process. Some apps want close-to-tray instead (Slack, Discord
        — clicking X hides the window but the process stays). Wire that
        with the <code>OnBeforeClose</code> hook in main.go.
      </TipBox>

      <h2>Keyboard shortcuts</h2>
      <ul>
        <li>
          <strong>Cmd/Ctrl+W</strong> — close window. Wire via Wails
          accelerators or with an HTML keyboard listener.
        </li>
        <li>
          <strong>Cmd/Ctrl+M</strong> — minimise. macOS does this
          automatically; Wails has it.
        </li>
        <li>
          <strong>F11</strong> — fullscreen toggle. Wire to{' '}
          <code>WindowFullscreen()</code>.
        </li>
      </ul>

      <h2>Edge cases users WILL report</h2>
      <ul>
        <li>
          <strong>Window remembers wrong size on relaunch.</strong> Save
          the geometry on close, restore on startup. Wails has{' '}
          <code>WindowGetSize / WindowSetPosition</code> for this.
        </li>
        <li>
          <strong>App restored on a disconnected monitor.</strong> Detect
          off-screen position; reset to centre.
        </li>
        <li>
          <strong>HiDPI / scale 200%.</strong> Test on a high-DPI display.
          Pixel measurements that look fine at 100% can blur at 200%.
        </li>
      </ul>

      <KnowledgeCheck
        question="A user wants the app to close to the system tray (like Slack), not quit. Where do you wire that?"
        choices={[
          {
            label: 'Override Quit in your Go App struct',
            feedback:
              "Wrong — Quit is the runtime call. Overriding it doesn\'t intercept the user clicking X.",
          },
          {
            label: 'OnBeforeClose hook in main.go — return false to cancel the close',
            correct: true,
            feedback:
              "Right — OnBeforeClose fires when the user clicks X. Return false to cancel the close + hide the window instead. Implement a system-tray icon to bring it back.",
          },
          {
            label: 'In React, listen for the close event and call preventDefault',
            feedback:
              "Wrong — by the time React would hear it, the window is closing. The hook is in Go, in main.go.",
          },
          {
            label: 'Disable the close button',
            feedback:
              "Disables the button but doesn\'t handle Alt+F4 or the OS close-window shortcut. OnBeforeClose handles all paths.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>For chapter 3&apos;s assignment, polish the titlebar to production quality:</p>
            <ol>
              <li>
                Add the OS detection — show OS-appropriate controls.
              </li>
              <li>Add double-click-to-maximize.</li>
              <li>
                Swap the maximize icon between single-square and
                restore-icon based on state.
              </li>
              <li>
                Save + restore the window size on close + relaunch.
              </li>
              <li>
                Test on at least one OS other than your dev machine
                (Hyper-V Windows VM if you&apos;re on Mac, OrbStack macOS
                if you&apos;re on Linux).
              </li>
            </ol>
          </>
        }
        hint={
          <>
            For size persistence, save to{' '}
            <code>localStorage</code> on a <code>beforeunload</code> event
            and restore on app startup via the Wails{' '}
            <code>WindowSetSize</code> call.
          </>
        }
        solution={
          <>
            <p>
              You should have a titlebar that feels native on every OS —
              traffic lights on Mac, square controls on Windows,
              double-click maximises, the maximise icon reflects state,
              and the window remembers its size across launches.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 4 — <strong>In-app auto-update</strong>. The most
        impactful feature for a shipped desktop app. Bug fixes land
        without users having to re-download installers.
      </p>
    </>
  )
}
