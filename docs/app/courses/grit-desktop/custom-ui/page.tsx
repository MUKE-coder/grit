import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import { LaneFlow } from '@/components/lane-flow'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Custom UI & Theming — Grit Desktop Course',
  description: 'Frameless window, custom title bar, sidebar navigation, dark theme, shadcn/ui components, and resizable panels in Grit desktop apps.',
}

export default function CustomUICourse() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <Link href="/courses/grit-desktop" className="hover:text-foreground transition-colors">Grit Desktop</Link>
          <span>/</span>
          <span className="text-foreground">Custom UI & Theming</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Course 3 of 5</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">10 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Custom UI & Theming
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            In this course, you will learn how to customize the look and feel of your desktop app.
            Frameless windows, custom title bars, sidebar navigation, dark theme CSS variables,
            shadcn/ui components, and resizable panels.
          </p>
          <LaneFlow id="c-gdui" lanes={['Desktop shell', 'Customize every part']}
            nodes={[
              { id: 'app', lane: 0, row: 2, title: 'Your app', sub: 'native shell', tone: 'primary' },
              { id: 'frame', lane: 1, row: 0, title: 'Frameless window', sub: 'custom title bar', tone: 'cyan' },
              { id: 'nav', lane: 1, row: 1, title: 'Sidebar nav', sub: 'layout', tone: 'blue' },
              { id: 'theme', lane: 1, row: 2, title: 'Dark theme', sub: 'CSS vars', tone: 'violet' },
              { id: 'ui', lane: 1, row: 3, title: 'shadcn/ui', sub: 'components', tone: 'green' },
            ]}
            edges={[{ from: 'app', to: 'frame', tone: 'cyan' }, { from: 'app', to: 'nav', label: 'style', tone: 'blue' }, { from: 'app', to: 'theme', tone: 'violet' }, { from: 'app', to: 'ui', tone: 'green' }]}
            legend={[{ tone: 'primary', label: 'Native shell' }, { tone: 'green', label: 'Fully yours' }]}
            caption="Frameless windows, custom title bars, sidebar navigation, dark theme, and shadcn/ui — fully yours" />
        </div>

        <div className="my-4 rounded-lg border border-primary/20 bg-primary/5 px-4 py-3">
          <p className="text-sm text-muted-foreground">
            <strong className="text-foreground">Reference docs:</strong>{' '}
            <a href="/docs/desktop" className="text-primary hover:underline">Desktop Overview →</a>
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ The Frameless Window ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The Frameless Window</h2>

          <Definition term="Frameless Window">
            A desktop window without the operating system{"'"}s default title bar and borders. Instead
            of the standard Windows/macOS/Linux title bar with its minimize, maximize, and close buttons,
            you build your own with React. This gives your app a modern, custom look — like VS Code,
            Spotify, or Discord.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit desktop apps use frameless windows by default. The setting lives in <Code>wails.json</Code>:
          </p>

          <CodeBlock filename="wails.json (partial)">
{`{
  "name": "myapp",
  "frontend:dir": "./frontend",
  "frontend:install": "pnpm install",
  "frontend:build": "pnpm build",
  "frontend:dev:watcher": "pnpm dev",
  "frontend:dev:serverUrl": "auto",
  "author": {
    "name": ""
  },
  "info": {},
  "nsisType": "exe",
  "obfuscated": false,
  "garbled": false,
  "frameless": true,
  "width": 1024,
  "height": 768,
  "minWidth": 800,
  "minHeight": 600
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When <Code>frameless</Code> is <Code>true</Code>, the OS title bar disappears entirely. Your React app
            fills the entire window. This means you need to build your own title bar with window controls.
          </p>

          <Challenge number={1} title="Find the Frameless Setting">
            <p>Open <Code>wails.json</Code> in your desktop project. Find the <Code>frameless</Code> setting. What are the <Code>width</Code> and <Code>height</Code> values? Try changing <Code>frameless</Code> to <Code>false</Code>, restart the app, and see what happens — you{"'"}ll get the default OS title bar.</p>
          </Challenge>
        </section>

        {/* ═══ Custom Title Bar ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Custom Title Bar</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Since there{"'"}s no OS title bar, Grit scaffolds a custom one with React. It includes
            the app name and window control buttons (minimize, maximize, close). These buttons
            call Wails runtime functions:
          </p>

          <CodeBlock filename="frontend/src/components/title-bar.tsx (simplified)">
{`import { WindowMinimise, WindowToggleMaximise, WindowClose } from '../../wailsjs/runtime'

export function TitleBar() {
  return (
    <div className="flex items-center justify-between h-10 bg-background border-b border-border px-4"
         style={{ "--wails-draggable": "drag" } as React.CSSProperties}>
      <span className="text-sm font-medium text-foreground">My App</span>

      <div className="flex items-center gap-1" style={{ "--wails-draggable": "no-drag" } as React.CSSProperties}>
        <button onClick={() => WindowMinimise()} className="p-1.5 hover:bg-bg-hover rounded">
          {/* Minimize icon */}
        </button>
        <button onClick={() => WindowToggleMaximise()} className="p-1.5 hover:bg-bg-hover rounded">
          {/* Maximize icon */}
        </button>
        <button onClick={() => WindowClose()} className="p-1.5 hover:bg-red-500/20 text-red-400 rounded">
          {/* Close icon */}
        </button>
      </div>
    </div>
  )
}`}
          </CodeBlock>

          <Definition term="Window Controls">
            Buttons that control the desktop window — minimize (hide to taskbar), maximize (fill screen),
            and close (exit the app). In a frameless window, you build these yourself with React instead
            of relying on the operating system{"'"}s built-in controls.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The Wails runtime provides these functions:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>WindowMinimise()</Code> — minimize the window to the taskbar</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>WindowToggleMaximise()</Code> — toggle between maximized and normal size</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>WindowClose()</Code> — close the application</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>WindowSetTitle(title)</Code> — change the window title programmatically</li>
          </ul>

          <Challenge number={2} title="Explore the Title Bar">
            <p>Find the title bar component in your project (look in <Code>frontend/src/components/</Code>). What Wails runtime functions does it import? Click each window control button — do minimize, maximize, and close all work?</p>
          </Challenge>
        </section>

        {/* ═══ Draggable Window ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Draggable Window</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Without the OS title bar, users cannot drag the window to move it. Wails solves this
            with a CSS property: <Code>--wails-draggable</Code>. Any element with this property set
            to <Code>drag</Code> becomes a drag handle:
          </p>

          <CodeBlock filename="Draggable region">
{`{/* This div can be dragged to move the window */}
<div style={{ "--wails-draggable": "drag" } as React.CSSProperties}>
  <span>Drag me to move the window</span>

  {/* Buttons inside should NOT be draggable */}
  <div style={{ "--wails-draggable": "no-drag" } as React.CSSProperties}>
    <button>Click me (not draggable)</button>
  </div>
</div>`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The title bar uses <Code>--wails-draggable: drag</Code> so you can grab it and move the window.
            The window control buttons use <Code>--wails-draggable: no-drag</Code> so clicking them
            triggers the button action instead of starting a drag.
          </p>

          <Challenge number={3} title="Test Dragging">
            <p>Try dragging the window by the title bar area — it should move. Now try dragging by the main content area — it should not move. Can you identify the exact CSS property that makes this work?</p>
          </Challenge>
        </section>

        {/* ═══ Sidebar Navigation ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Sidebar Navigation</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The sidebar is auto-generated from your resources. When you run <Code>grit generate resource</Code>,
            a new sidebar entry is injected automatically with an appropriate icon:
          </p>

          <CodeBlock filename="Sidebar entries (auto-generated)">
{`const sidebarItems = [
  { label: "Dashboard", icon: LayoutDashboard, href: "/" },
  { label: "Blogs", icon: FileText, href: "/blogs" },
  { label: "Contacts", icon: Users, href: "/contacts" },
  // Automatically added when you generate resources:
  { label: "Tasks", icon: CheckSquare, href: "/tasks" },
  { label: "Notes", icon: StickyNote, href: "/notes" },
]`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The sidebar is collapsible — it can shrink to just icons, giving more space to the main content.
            Grit picks Lucide icons based on the resource name (e.g., {'"'}Task{'"'} gets a check icon,
            {'"'}Note{'"'} gets a sticky note icon).
          </p>

          <Challenge number={4} title="Sidebar Auto-Updates">
            <p>Generate 3 new resources (e.g., Project, Invoice, Report). Restart the app. Does the sidebar show 3 new entries with icons? Check the sidebar component to see how the entries are defined.</p>
          </Challenge>
        </section>

        {/* ═══ Dark Theme ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Dark Theme</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit desktop apps use a premium dark theme defined with CSS custom properties (variables).
            These variables are used by all shadcn/ui components and your custom components:
          </p>

          <CodeBlock filename="frontend/src/index.css (theme variables)">
{`:root {
  --background: #0a0a0f;     /* Darkest — main background */
  --bg-elevated: #22222e;    /* Card surfaces */
  --bg-hover: #2a2a38;       /* Hover states */
  --border: #2a2a3a;         /* Borders and dividers */
  --foreground: #e8e8f0;     /* Primary text */
  --text-secondary: #9090a8; /* Secondary text */
  --text-muted: #606078;     /* Muted/disabled text */
  --primary: #6c5ce7;        /* Accent purple */
  --primary-hover: #7c6cf7;  /* Accent hover */
  --success: #00b894;        /* Green for success states */
  --danger: #ff6b6b;         /* Red for errors and delete */
  --warning: #fdcb6e;        /* Yellow for warnings */
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You use these variables through Tailwind CSS classes:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>bg-background</Code> — the darkest background (#0a0a0f)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>bg-bg-elevated</Code> — card and panel surfaces (#22222e)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>text-foreground</Code> — primary white text (#e8e8f0)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>text-primary</Code> — accent purple (#6c5ce7)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>border-border</Code> — subtle borders (#2a2a3a)</li>
          </ul>

          <Tip>
            To change the entire color scheme of your app, just modify the CSS variables. Change
            <Code>--primary</Code> from purple to blue and every button, link, and accent in the entire app
            updates instantly. No need to find and replace colors across dozens of files.
          </Tip>

          <Challenge number={5} title="Customize the Theme">
            <p>Find the CSS file with theme variables (look in <Code>frontend/src/</Code>). Change the <Code>--primary</Code> color from <Code>#6c5ce7</Code> to <Code>#3b82f6</Code> (blue). Save and check — does the entire app{"'"}s accent color update?</p>
          </Challenge>
        </section>

        {/* ═══ shadcn/ui in Desktop ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">shadcn/ui in Desktop</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit desktop apps use the same <a href="https://ui.shadcn.com" target="_blank" rel="noreferrer" className="text-primary hover:underline">shadcn/ui</a> components
            as web projects. These are unstyled, accessible React components that you can customize:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Component</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Used for</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Button</td>
                  <td className="px-4 py-3">Primary actions, form submissions, navigation</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Input</td>
                  <td className="px-4 py-3">Text fields in forms</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Dialog</td>
                  <td className="px-4 py-3">Modal dialogs for confirmations and forms</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">DataTable</td>
                  <td className="px-4 py-3">Sortable, filterable data tables for listing resources</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Select</td>
                  <td className="px-4 py-3">Dropdown selection (e.g., category picker)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Checkbox</td>
                  <td className="px-4 py-3">Boolean toggles (e.g., task done/not done)</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Textarea</td>
                  <td className="px-4 py-3">Multi-line text input for descriptions</td>
                </tr>
              </tbody>
            </table>
          </div>

          <p className="text-muted-foreground leading-relaxed mb-4">
            All components live in <Code>frontend/src/components/ui/</Code> and are imported like:
          </p>

          <CodeBlock filename="Import example">
{`import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Dialog, DialogContent, DialogHeader } from '@/components/ui/dialog'`}
          </CodeBlock>

          <Challenge number={6} title="Find shadcn/ui Components">
            <p>Search your project for imports from <Code>@/components/ui/</Code>. Find at least 5 different shadcn/ui components used in the app. Which component is used most frequently?</p>
          </Challenge>
        </section>

        {/* ═══ Resizable Panels ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Resizable Panels</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit desktop apps can include resizable split panels, similar to IDE-style interfaces where
            you can drag a divider between the sidebar and main content to adjust their widths:
          </p>

          <CodeBlock filename="Resizable layout pattern">
{`import { ResizablePanel, ResizablePanelGroup, ResizableHandle } from '@/components/ui/resizable'

export function AppLayout({ children }) {
  return (
    <ResizablePanelGroup direction="horizontal">
      <ResizablePanel defaultSize={20} minSize={15} maxSize={30}>
        <Sidebar />
      </ResizablePanel>
      <ResizableHandle />
      <ResizablePanel defaultSize={80}>
        {children}
      </ResizablePanel>
    </ResizablePanelGroup>
  )
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>ResizableHandle</Code> is the draggable divider between panels. Users can grab it
            and drag to resize. The <Code>minSize</Code> and <Code>maxSize</Code> props prevent panels from
            becoming too small or too large.
          </p>

          <Challenge number={7} title="Test Resizable Panels">
            <p>If your app has resizable panels, try dragging the divider between the sidebar and main content. The sidebar should resize smoothly. What are the minimum and maximum sizes?</p>
          </Challenge>
        </section>

        {/* ═══ Window Configuration ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Window Configuration</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>wails.json</Code> file controls window appearance and behavior:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Setting</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Default</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">What it does</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">width</td>
                  <td className="px-4 py-3">1024</td>
                  <td className="px-4 py-3">Default window width in pixels</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">height</td>
                  <td className="px-4 py-3">768</td>
                  <td className="px-4 py-3">Default window height in pixels</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">minWidth</td>
                  <td className="px-4 py-3">800</td>
                  <td className="px-4 py-3">Minimum window width (cannot shrink smaller)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">minHeight</td>
                  <td className="px-4 py-3">600</td>
                  <td className="px-4 py-3">Minimum window height</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">frameless</td>
                  <td className="px-4 py-3">true</td>
                  <td className="px-4 py-3">Hide OS title bar for custom UI</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-mono text-foreground text-xs">fullscreen</td>
                  <td className="px-4 py-3">false</td>
                  <td className="px-4 py-3">Start in fullscreen mode</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Challenge number={8} title="Configure the Window">
            <p>Change the window size to 1200x800 in <Code>wails.json</Code>. Set the minimum size to 900x650. Restart the app — is the window larger? Try resizing it below the minimum — does it stop?</p>
          </Challenge>
        </section>

        {/* ═══ Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What You Learned</h2>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How frameless windows work — hiding the OS title bar for a custom look</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Building a custom title bar with Wails runtime window controls</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Making windows draggable with <Code>--wails-draggable</Code> CSS property</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Auto-generated sidebar navigation from resources</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> The dark theme system with CSS custom properties</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Using shadcn/ui components in desktop apps</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Resizable panels for IDE-style layouts</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Window configuration in wails.json</li>
          </ul>

          <Challenge number={9} title="Full Customization Challenge">
            <p>Customize your app in four ways: (1) Change the accent color from purple to a color of your choice. (2) Modify the title bar text to show your app{"'"}s name. (3) Add a new icon to the sidebar for a section called {'"'}Settings{'"'}. (4) Adjust the default window size to 1280x720 in wails.json.</p>
          </Challenge>

          <Challenge number={10} title="Theme Experiment">
            <p>Create a second theme by duplicating the CSS variables with different colors (e.g., a light theme or a blue-tinted dark theme). Add a button in the title bar that toggles between themes by swapping CSS classes on the root element. Does the entire app update?</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses/grit-desktop/crud-data', label: 'Prev: Desktop CRUD & Data' }}
            next={{ href: '/courses/grit-desktop/export', label: 'Next: PDF & Excel Export' }}
          />
        </div>
      </main>
    </div>
  )
}
