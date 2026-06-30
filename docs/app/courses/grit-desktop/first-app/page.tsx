import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Your First Desktop App — Grit Desktop Course',
  description: 'Install Wails, scaffold a desktop app with Grit, understand the project structure, Wails bindings, and run in dev mode.',
}

export default function FirstDesktopAppCourse() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <Link href="/courses/grit-desktop" className="hover:text-foreground transition-colors">Grit Desktop</Link>
          <span>/</span>
          <span className="text-foreground">Your First Desktop App</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Course 1 of 5</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">12 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Your First Desktop App
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            In this course, you will install Wails, scaffold a desktop application with Grit,
            understand how Go functions are called directly from React, and run your app in dev mode.
            By the end, you will have a working desktop app with authentication, blog CRUD, and a SQLite database.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ What is a Desktop App? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is a Desktop App?</h2>

          <Definition term="Desktop Application">
            A program that runs directly on your computer{"'"}s operating system (Windows, macOS, or Linux).
            Unlike a web app that runs in a browser, a desktop app is a standalone window on your computer.
            Think VS Code, Spotify, or Discord — they run without needing a browser.
          </Definition>

          <Definition term="Native App">
            An application compiled specifically for one operating system. It can access local files,
            system notifications, the clipboard, and other OS features that web apps cannot.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Desktop apps are different from web apps in several important ways:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> They run <strong className="text-foreground">locally</strong> — no browser needed, no internet required (unless your app needs it)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> They can <strong className="text-foreground">access the filesystem</strong> — read and write files on the user{"'"}s computer</li>
            <li className="flex gap-2"><span className="text-primary">•</span> They feel <strong className="text-foreground">faster</strong> — no network latency between the UI and backend</li>
            <li className="flex gap-2"><span className="text-primary">•</span> They are <strong className="text-foreground">distributed as a single file</strong> — users download one executable and run it</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit uses <strong className="text-foreground">Wails</strong> to build desktop apps. You write Go for
            the backend and React for the UI — the same languages you already know from Grit web.
          </p>
        </section>

        {/* ═══ What is Wails? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is Wails?</h2>

          <Definition term="Wails">
            A Go framework for building desktop applications using web technologies.
            Go handles the backend logic, React handles the UI. The key difference from web apps:
            Go functions are <strong className="text-foreground">directly callable from JavaScript</strong> — no HTTP requests,
            no REST APIs, no fetch(). Just function calls.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Wails works by embedding a webview (a lightweight browser engine) inside a native window.
            Your React app runs in this webview, and Wails creates a bridge so JavaScript can call Go
            functions directly:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Go side:</strong> You write normal Go methods like <Code>GetUsers()</Code> or <Code>CreateBlog(input)</Code></li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">React side:</strong> You import and call them as if they were local functions — <Code>{'const users = await GetUsers()'}</Code></li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">No HTTP layer:</strong> No REST endpoints, no fetch(), no axios — direct function calls over an internal bridge</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Learn more at <a href="https://wails.io" target="_blank" rel="noreferrer" className="text-primary hover:underline">wails.io</a>.
            Wails supports Windows, macOS, and Linux from a single codebase.
          </p>

          <Challenge number={1} title="Visit Wails.io">
            <p>Open <a href="https://wails.io" target="_blank" rel="noreferrer" className="text-primary hover:underline">wails.io</a> in your browser. What platforms does Wails support? Find the {'"'}Getting Started{'"'} section — does it mention Go as a requirement?</p>
          </Challenge>
        </section>

        {/* ═══ Prerequisites ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Prerequisites</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You need the same tools as Grit web (Go, Node.js, pnpm) plus one more — the Wails CLI:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Tool</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Version</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Check</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Go</td>
                  <td className="px-4 py-3">1.21+</td>
                  <td className="px-4 py-3"><Code>go version</Code></td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Node.js</td>
                  <td className="px-4 py-3">18+</td>
                  <td className="px-4 py-3"><Code>node --version</Code></td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">pnpm</td>
                  <td className="px-4 py-3">8+</td>
                  <td className="px-4 py-3"><Code>pnpm --version</Code></td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Wails CLI</td>
                  <td className="px-4 py-3">v2</td>
                  <td className="px-4 py-3"><Code>wails version</Code></td>
                </tr>
              </tbody>
            </table>
          </div>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Install the Wails CLI with:
          </p>

          <CodeBlock filename="Terminal">
{`go install github.com/wailsapp/wails/v2/cmd/wails@latest`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            After installation, verify everything is working:
          </p>

          <CodeBlock filename="Terminal">
{`wails doctor`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>wails doctor</Code> command checks your system for all required dependencies
            (Go, Node, npm/pnpm, platform-specific build tools) and reports any issues.
          </p>

          <Challenge number={2} title="Install Wails CLI">
            <p>Run <Code>go install github.com/wailsapp/wails/v2/cmd/wails@latest</Code> to install the Wails CLI. Then run <Code>wails doctor</Code>. Does it pass all checks? If any checks fail, follow the instructions to fix them.</p>
          </Challenge>
        </section>

        {/* ═══ Scaffold a Desktop App ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Scaffold a Desktop App</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The command for desktop apps is <Code>grit new-desktop</Code> — this is different
            from <Code>grit new</Code> which creates a web monorepo:
          </p>

          <CodeBlock filename="Terminal">
{`grit new-desktop myapp`}
          </CodeBlock>

          <Note>
            <Code>grit new-desktop</Code> creates a <strong className="text-foreground">standalone app</strong>, not a monorepo.
            There is no <Code>apps/</Code> folder, no <Code>packages/</Code> folder, no Turborepo. Everything lives
            in a single flat directory. This is intentional — desktop apps are self-contained.
          </Note>

          <p className="text-muted-foreground leading-relaxed mb-4">
            After running the command, Grit creates a complete desktop project with authentication,
            blog CRUD, contact CRUD, a dark theme, and a SQLite database — all wired up and ready to run.
          </p>

          <Challenge number={3} title="Scaffold Your First Desktop App">
            <p>Run <Code>grit new-desktop notes-app</Code>. Look at the folder structure that was created. How many Go files are in the root directory? How many folders are there?</p>
          </Challenge>
        </section>

        {/* ═══ Project Structure ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Project Structure</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            A Grit desktop project has a flat, simple structure:
          </p>

          <CodeBlock filename="Project Structure">
{`notes-app/
├── main.go          <- Wails bootstrap + Go embed
├── app.go           <- Bound methods (Go functions callable from React)
├── internal/        <- Go backend (models, services, config, db)
├── frontend/        <- React app (Vite + TanStack Router)
├── wails.json       <- Wails configuration
└── .env             <- Environment variables`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Compare this to a Grit web project which has <Code>apps/api</Code>, <Code>apps/web</Code>, <Code>apps/admin</Code>, and <Code>packages/shared</Code>.
            A desktop project is much simpler — everything is in one place.
          </p>

          <Definition term="Wails Bindings">
            Go functions exposed to JavaScript. You write a Go method like <Code>GetUsers()</Code> and
            call it from React as if it were a local function. No HTTP, no REST — direct function calls.
            Wails automatically generates TypeScript wrappers for every bound Go method.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The key file is <Code>app.go</Code> — this is where all your Go methods that React can call are defined.
            Every public method on the <Code>App</Code> struct becomes available in the frontend:
          </p>

          <CodeBlock filename="app.go (simplified)">
{`type App struct {
    ctx         context.Context
    authService *services.AuthService
    blogService *services.BlogService
}

// React can call: Login({ email, password })
func (a *App) Login(input types.LoginInput) (*types.AuthResponse, error) {
    return a.authService.Login(input)
}

// React can call: GetBlogs()
func (a *App) GetBlogs() ([]models.Blog, error) {
    return a.blogService.List()
}`}
          </CodeBlock>

          <Challenge number={4} title="Explore app.go">
            <p>Open <Code>app.go</Code> in your notes-app project. What methods are available? These are the functions React can call. Count how many public methods the App struct has.</p>
          </Challenge>
        </section>

        {/* ═══ Running in Dev Mode ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Running in Dev Mode</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Start your desktop app with either command:
          </p>

          <CodeBlock filename="Terminal">
{`cd notes-app
grit start`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Or use the Wails command directly:
          </p>

          <CodeBlock filename="Terminal">
{`wails dev`}
          </CodeBlock>

          <Definition term="Hot Reload">
            The app automatically updates when you save a file. Change Go code and the backend restarts.
            Change React code and the UI updates instantly. No manual restart needed. This makes
            development fast — save a file and see the result immediately.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When you run <Code>grit start</Code>, a native desktop window opens with your app. The frontend
            dev server runs at <Code>localhost:34115</Code> (you can also open it in a browser for debugging).
          </p>

          <Tip>
            If you need to debug the frontend, right-click inside the desktop window and select {'"'}Inspect Element{'"'}
            to open the browser DevTools — just like in Chrome.
          </Tip>

          <Challenge number={5} title="Run Your Desktop App">
            <p>Run <Code>grit start</Code> in your notes-app directory. The app should open as a desktop window. Try the login page — can you see the registration form?</p>
          </Challenge>
        </section>

        {/* ═══ SQLite Database ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">SQLite Database</h2>

          <Definition term="SQLite">
            A file-based database that requires no server. Your entire database lives in a single <Code>.db</Code> file
            in your project folder. No Docker, no PostgreSQL installation, no configuration. SQLite is the
            most widely deployed database engine in the world — it{"'"}s inside every smartphone, every browser,
            and now inside your desktop app.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Desktop apps use SQLite instead of PostgreSQL because users should not need to install
            a database server. The database file travels with the app — it{"'"}s just a file on disk:
          </p>

          <CodeBlock filename=".env">
{`DB_PATH=./data.db`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When the app starts, GORM automatically creates the <Code>.db</Code> file and runs migrations.
            You can browse the database with GORM Studio:
          </p>

          <CodeBlock filename="Terminal">
{`grit studio`}
          </CodeBlock>

          <Challenge number={6} title="Find the Database File">
            <p>After running the app and registering an account, find the <Code>.db</Code> file in your project folder. Open GORM Studio with <Code>grit studio</Code> and browse the tables. Can you see the user you just registered?</p>
          </Challenge>
        </section>

        {/* ═══ Tour the App ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Tour the App</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Your scaffolded desktop app comes with everything you need to start building:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Login and Register pages</strong> — local authentication with bcrypt password hashing</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Dashboard</strong> — stats cards showing total blogs, contacts, and recent activity</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Blog CRUD</strong> — create, read, update, and delete blog posts</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Contact CRUD</strong> — manage a contacts list with full CRUD operations</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Sidebar navigation</strong> — collapsible sidebar with icons for every section</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Dark theme</strong> — premium dark UI using Grit{"'"}s design system</li>
          </ul>

          <Challenge number={7} title="Create Some Data">
            <p>Register an account in your app. Create 3 blog posts with different titles and content. Create 2 contacts with names and emails. Visit the dashboard — do the stats update to show your new data?</p>
          </Challenge>
        </section>

        {/* ═══ The .env File ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The .env File</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Desktop apps have a simpler <Code>.env</Code> file compared to web projects — no database URLs,
            no Redis, no S3 storage:
          </p>

          <CodeBlock filename=".env">
{`APP_NAME=notes-app
JWT_SECRET=your-secret-key-change-in-production
DB_PATH=./data.db
TOTP_ISSUER=notes-app`}
          </CodeBlock>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>APP_NAME</Code> — your application name, shown in the title bar</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>JWT_SECRET</Code> — secret key for signing authentication tokens</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>DB_PATH</Code> — path to the SQLite database file</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>TOTP_ISSUER</Code> — name shown in authenticator apps for two-factor auth</li>
          </ul>

          <Challenge number={8} title="Read the .env File">
            <p>Open the <Code>.env</Code> file in your project. What is the database file name? What is the JWT_SECRET set to?</p>
          </Challenge>
        </section>

        {/* ═══ Essential Commands ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Essential Commands</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here are the commands you{"'"}ll use every day when building desktop apps:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Command</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">What it does</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">grit start</td>
                  <td className="px-4 py-3">Start the app in dev mode with hot reload</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">grit studio</td>
                  <td className="px-4 py-3">Open the database browser for your SQLite database</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">grit generate resource</td>
                  <td className="px-4 py-3">Generate a new CRUD resource (model, service, routes, UI)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">grit routes</td>
                  <td className="px-4 py-3">List all bound methods (desktop equivalent of API routes)</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-mono text-foreground text-xs">grit compile</td>
                  <td className="px-4 py-3">Build a production-ready executable</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Challenge number={9} title="Try the Commands">
            <p>Run <Code>grit routes</Code> in your project folder. How many routes does the desktop app have? Can you identify which ones handle blogs and which handle contacts?</p>
          </Challenge>

          <Challenge number={10} title="Explore GORM Studio">
            <p>Run <Code>grit studio</Code>. Browse the tables in your SQLite database. Can you see the blogs and contacts you created earlier? Try editing a record directly in Studio.</p>
          </Challenge>
        </section>

        {/* ═══ Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What You Learned</h2>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> What desktop apps are and how they differ from web apps</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> What Wails is — Go + React desktop framework with direct function calls</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How to install the Wails CLI and verify with <Code>wails doctor</Code></li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How to scaffold a desktop app with <Code>grit new-desktop</Code></li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> The flat project structure (main.go, app.go, internal/, frontend/)</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How Wails bindings work — Go methods callable from React</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How to run in dev mode with <Code>grit start</Code> and hot reload</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> SQLite as the database — file-based, no server required</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> The essential CLI commands for desktop development</li>
          </ul>

          <Challenge number={11} title="Final Challenge: Build a Todo App">
            <p>Create a brand new desktop app called <Code>todo-app</Code> with <Code>grit new-desktop todo-app</Code>. Start it, register an account, create 5 blog posts, then open GORM Studio to see the data in SQLite. Verify the <Code>.db</Code> file exists in the project folder.</p>
          </Challenge>

          <Challenge number={12} title="Bonus: Compare Web vs Desktop">
            <p>If you{"'"}ve completed the Grit Web course, create a web project with <Code>grit new webtest</Code> and a desktop project with <Code>grit new-desktop desktest</Code> side by side. List 3 differences you notice in the project structure, the way data is fetched, and the dev experience.</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses/grit-desktop', label: 'Back to Grit Desktop' }}
            next={{ href: '/courses/grit-desktop/crud-data', label: 'Next: Desktop CRUD & Data' }}
          />
        </div>
      </main>
    </div>
  )
}
