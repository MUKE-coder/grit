import Link from "next/link";
import { ArrowLeft, ArrowRight, Download } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { CodeBlock } from "@/components/code-block";
import { Files, Folder, File } from "@/components/files";
import { Steps, Step } from "@/components/steps";
import { getDocMetadata } from "@/config/docs-metadata";

export const metadata = getDocMetadata("/docs/desktop/first-app");

export default function DesktopFirstAppPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">
                Desktop (Wails)
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Your First Desktop App
              </h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                Build a Task Manager desktop application from scratch. This
                step-by-step tutorial covers scaffolding, development, resource
                generation, GORM Studio, and building for distribution.
              </p>
            </div>

            <div className="rounded-lg border border-primary/20 bg-primary/5 p-4 mb-10">
              <p className="text-sm text-muted-foreground leading-relaxed">
                <strong className="text-foreground">Prefer a narrative walkthrough?</strong> This build has a
                story-form companion on the blog &mdash; an offline-first POS:{' '}
                <Link href="/blog/build-desktop-app-with-grit" className="text-primary hover:underline">Build a desktop app with Grit &rarr;</Link>
              </p>
            </div>

            {/* What you'll build */}
            <div className="rounded-lg border border-border/30 bg-card/30 px-5 py-4 mb-10">
              <h3 className="font-semibold text-foreground mb-2">
                What you&apos;ll build
              </h3>
              <p className="text-sm text-muted-foreground mb-3">
                A native desktop Task Manager with:
              </p>
              <ul className="text-sm text-muted-foreground space-y-1.5">
                <li>
                  &bull; Full CRUD for tasks (title, description, priority,
                  completion status, due date)
                </li>
                <li>
                  &bull; A categories resource to organize tasks
                </li>
                <li>
                  &bull; Search, pagination, PDF and Excel export
                </li>
                <li>
                  &bull; Authentication with login and default admin user
                </li>
                <li>
                  &bull; Visual database browsing with GORM Studio
                </li>
                <li>
                  &bull; A single native executable for distribution
                </li>
              </ul>
            </div>

            <Steps>
            {/* Step 1: Prerequisites */}
            <Step title="Prerequisites">
              <div className="prose-grit mb-4">
                <p>
                  Make sure the following tools are installed on your system
                  before starting:
                </p>
              </div>

              <div className="grid gap-3 sm:grid-cols-2 mb-6">
                {[
                  { name: "Go", version: "1.21+", check: "go version" },
                  { name: "Node.js", version: "18+", check: "node --version" },
                  { name: "Wails CLI", version: "v2", check: "wails version" },
                  { name: "Grit CLI", version: "Latest", check: "grit --help" },
                ].map((tool) => (
                  <div
                    key={tool.name}
                    className="rounded-lg border border-border/30 bg-card/30 px-4 py-3"
                  >
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-[15px] font-semibold">
                        {tool.name}
                      </span>
                      <span className="text-sm font-mono text-primary/60">
                        {tool.version}
                      </span>
                    </div>
                    <code className="text-sm font-mono text-muted-foreground/50">
                      {tool.check}
                    </code>
                  </div>
                ))}
              </div>

              <div className="prose-grit mb-4">
                <p>Install the Wails CLI:</p>
              </div>
              <CodeBlock
                terminal
                code="go install github.com/wailsapp/wails/v2/cmd/wails@latest"
                className="mb-4"
              />

              <div className="prose-grit mb-4">
                <p>Install the Grit CLI:</p>
              </div>
              <CodeBlock
                terminal
                code={`# macOS / Linux
curl -fsSL https://gritframework.dev/install.sh | sh

# Windows (PowerShell)
iwr -useb https://gritframework.dev/install.ps1 | iex`}
                className="mb-4"
              />

              <div className="prose-grit mb-0">
                <p>
                  Run <code>wails doctor</code> to verify your environment. It
                  checks for Go, Node.js, npm/pnpm, and platform-specific build
                  tools (GCC on Linux, Xcode on macOS, or WebView2 on Windows).
                </p>
              </div>
              <CodeBlock terminal code="wails doctor" className="mb-0" />
            </Step>

            {/* Step 2: Scaffold the Project */}
            <Step title="Scaffold the Project">
              <div className="prose-grit mb-4">
                <p>
                  Create a new desktop project called{" "}
                  <code>taskmanager</code> with the Grit CLI. This generates a
                  complete Wails application with everything included.
                </p>
              </div>
              <CodeBlock
                terminal
                code="grit new-desktop taskmanager"
                className="mb-6 glow-purple-sm"
              />

              <div className="prose-grit mb-4">
                <p>
                  Here is the project structure that gets created:
                </p>
              </div>
              <Files title="taskmanager/">
                <File name="main.go" comment="Wails entry point" />
                <File name="app.go" comment="App struct with bound methods" />
                <File name="wails.json" comment="Wails project configuration" />
                <File name="go.mod" />
                <File name="go.sum" />
                <Folder name="internal" defaultOpen>
                  <Folder name="config" defaultOpen>
                    <File name="config.go" comment="App configuration" />
                  </Folder>
                  <Folder name="db" defaultOpen>
                    <File name="db.go" comment="GORM database setup (SQLite)" />
                  </Folder>
                  <Folder name="models" defaultOpen>
                    <File name="user.go" comment="User model + AutoMigrate" />
                    <File name="blog.go" comment="Blog post model" />
                    <File name="contact.go" comment="Contact model" />
                  </Folder>
                  <Folder name="service" defaultOpen>
                    <File name="auth.go" comment="Authentication service" />
                    <File name="blog.go" comment="Blog CRUD service" />
                    <File name="contact.go" comment="Contact CRUD service" />
                    <File name="export.go" comment="PDF / Excel export service" />
                  </Folder>
                  <Folder name="api" comment="Embedded Gin REST API (127.0.0.1:34999)" />
                  <Folder name="storage" comment="Upload storage in the OS app-data dir" />
                  <Folder name="files" comment="FileRef — stored upload JSON shape" />
                  <Folder name="types" defaultOpen>
                    <File name="types.go" comment="Shared request/response types" />
                  </Folder>
                </Folder>
                <Folder name="frontend" defaultOpen>
                  <Folder name="src" defaultOpen>
                    <File name="main.tsx" comment="React entry point (TanStack Router)" />
                    <Folder name="routes" comment="File-based routes (TanStack Router)" defaultOpen>
                      <File name="__root.tsx" comment="Root route" />
                      <File name="_layout.tsx" comment="Auth guard + sidebar layout" />
                      <Folder name="_layout" comment="Protected page routes" />
                    </Folder>
                    <Folder name="components" comment="Reusable UI components" />
                    <Folder name="hooks" comment="TanStack Query hooks" />
                    <Folder name="lib" comment="Utilities" />
                  </Folder>
                  <File name="index.html" />
                  <File name="package.json" />
                  <File name="vite.config.ts" />
                  <File name="tailwind.config.ts" />
                </Folder>
                <Folder name="cmd" defaultOpen>
                  <Folder name="studio" defaultOpen>
                    <File name="main.go" comment="GORM Studio standalone server" />
                  </Folder>
                </Folder>
              </Files>

              <div className="prose-grit mb-0">
                <p>What gets created out of the box:</p>
              </div>
              <div className="space-y-3 mt-4">
                {[
                  {
                    label: "Go backend with Wails bindings",
                    desc: "Go services are bound to the Wails runtime so React can call them directly for CRUD, while an embedded Gin REST API on 127.0.0.1:34999 handles native uploads, health, and external clients.",
                  },
                  {
                    label: "React frontend with Vite",
                    desc: "A full React app with Tailwind CSS, TanStack Router, TanStack Query, and pre-built pages.",
                  },
                  {
                    label: "SQLite database",
                    desc: "Zero-config local database. GORM AutoMigrate runs on startup.",
                  },
                  {
                    label: "Authentication",
                    desc: "Login and register pages with JWT-based auth. A default admin user is seeded automatically.",
                  },
                  {
                    label: "Blog + Contact CRUD",
                    desc: "Two fully working resources with list pages, forms, search, pagination, and PDF/Excel export.",
                  },
                  {
                    label: "GORM Studio",
                    desc: "A standalone visual database browser you can launch alongside the app.",
                  },
                ].map((item) => (
                  <div
                    key={item.label}
                    className="rounded-lg border border-border/30 bg-card/30 px-4 py-3"
                  >
                    <span className="text-[15px] font-semibold block mb-1">
                      {item.label}
                    </span>
                    <span className="text-sm text-muted-foreground/60">
                      {item.desc}
                    </span>
                  </div>
                ))}
              </div>
            </Step>

            {/* Step 3: Start Development */}
            <Step title="Start Development">
              <div className="prose-grit mb-4">
                <p>
                  Navigate into the project and start Wails in development mode:
                </p>
              </div>
              <CodeBlock
                terminal
                code={`cd taskmanager
wails dev`}
                className="mb-4 glow-purple-sm"
              />

              <div className="prose-grit mb-4">
                <p>
                  When you run <code>wails dev</code>, the following happens:
                </p>
              </div>

              <div className="space-y-2 mb-6">
                {[
                  "A native desktop window opens with your app",
                  "The Go backend compiles and starts with Wails bindings",
                  "The Vite dev server starts for the React frontend",
                  "Hot-reload is active for both Go and React changes",
                ].map((item, i) => (
                  <div
                    key={i}
                    className="flex items-start gap-3 rounded-lg border border-border/30 bg-card/30 px-4 py-2.5"
                  >
                    <span className="text-sm font-mono text-primary/60 shrink-0 mt-0.5">
                      {i + 1}.
                    </span>
                    <span className="text-sm text-muted-foreground">
                      {item}
                    </span>
                  </div>
                ))}
              </div>

              <div className="prose-grit mb-4">
                <p>
                  The frontend dev server also runs at{" "}
                  <code>http://localhost:34115</code>, which you can open in a
                  browser for debugging with browser DevTools.
                </p>
                <blockquote>
                  On first run, <code>wails dev</code> installs frontend npm
                  dependencies automatically. This takes a minute or two.
                  Subsequent starts are much faster.
                </blockquote>
              </div>
            </Step>

            {/* Step 4: Explore the Default App */}
            <Step title="Explore the Default App">
              <div className="prose-grit mb-4">
                <p>
                  Once the desktop window opens, log in with the default admin
                  credentials:
                </p>
              </div>

              <div className="rounded-lg border border-border/30 bg-card/30 px-4 py-3 mb-6">
                <div className="grid grid-cols-2 gap-2 text-sm">
                  <div>
                    <span className="text-muted-foreground/60">Email:</span>{" "}
                    <code className="text-primary/80">admin@example.com</code>
                  </div>
                  <div>
                    <span className="text-muted-foreground/60">Password:</span>{" "}
                    <code className="text-primary/80">password</code>
                  </div>
                </div>
              </div>

              <div className="prose-grit mb-4">
                <p>
                  After logging in, explore what the scaffold gives you out of
                  the box:
                </p>
              </div>

              <div className="space-y-3 mb-6">
                {[
                  {
                    section: "Dashboard",
                    desc: "A landing page with stats cards showing total blogs, contacts, and users.",
                  },
                  {
                    section: "Blog",
                    desc: "Full CRUD with list table, create/edit forms, search, pagination, and PDF/Excel export.",
                  },
                  {
                    section: "Contacts",
                    desc: "Another complete CRUD resource following the same pattern as Blog.",
                  },
                ].map((item) => (
                  <div
                    key={item.section}
                    className="rounded-lg border border-border/30 bg-card/30 px-4 py-3"
                  >
                    <div className="flex items-center gap-2 mb-1">
                      <span className="text-[15px] font-semibold">
                        {item.section}
                      </span>
                    </div>
                    <span className="text-sm text-muted-foreground/60">
                      {item.desc}
                    </span>
                  </div>
                ))}
              </div>

              <div className="prose-grit mb-4">
                <p>
                  These built-in resources demonstrate every feature that your
                  generated resources will also have: search, sorting,
                  pagination, inline editing, bulk operations, and export.
                </p>
              </div>
            </Step>

            {/* Step 5: Generate the Task Resource */}
            <Step title="Generate the Task Resource">
              <div className="prose-grit mb-4">
                <p>
                  Now for the main event. Open a new terminal in the{" "}
                  <code>taskmanager</code> directory and run:
                </p>
              </div>
              <CodeBlock
                terminal
                code={`grit generate resource Task --fields "title:string,description:text,priority:string,completed:bool,due_date:date"`}
                className="mb-6 glow-purple-sm"
              />

              <div className="prose-grit mb-4">
                <p>This single command creates 5 new files:</p>
              </div>

              <div className="space-y-3 mb-6">
                {[
                  {
                    file: "internal/models/task.go",
                    desc: "GORM model struct with all fields, timestamps, and soft delete",
                  },
                  {
                    file: "internal/service/task.go",
                    desc: "Service with List, ListAll, GetByID, Create, Update, Delete methods",
                  },
                  {
                    file: "frontend/src/routes/_layout/tasks.index.tsx",
                    desc: "List route with search, pagination, edit/delete, PDF and Excel export",
                  },
                  {
                    file: "frontend/src/routes/_layout/tasks.new.tsx",
                    desc: "Create form route with inputs mapped to each field type",
                  },
                  {
                    file: "frontend/src/routes/_layout/tasks.$id.edit.tsx",
                    desc: "Edit form route with pre-filled field values",
                  },
                ].map((item) => (
                  <div
                    key={item.file}
                    className="rounded-lg border border-border bg-card p-4"
                  >
                    <code className="text-sm font-semibold text-primary">
                      {item.file}
                    </code>
                    <p className="text-sm text-muted-foreground mt-1">
                      {item.desc}
                    </p>
                  </div>
                ))}
              </div>

              <div className="prose-grit mb-4">
                <p>
                  It also injects code into 10 locations in existing files using{" "}
                  <code>grit:</code> markers:
                </p>
              </div>

              <div className="overflow-x-auto mb-6">
                <table className="w-full text-sm border border-border rounded-lg">
                  <thead>
                    <tr className="border-b border-border bg-muted/50">
                      <th className="text-left p-3 font-medium">File</th>
                      <th className="text-left p-3 font-medium">Marker</th>
                      <th className="text-left p-3 font-medium">What</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border text-muted-foreground">
                    {[
                      ["db.go", "// grit:models", "Task model in AutoMigrate"],
                      ["main.go", "// grit:service-init", "TaskService initialization"],
                      ["main.go", "/* grit:app-args */", "Service passed to NewApp"],
                      ["app.go", "// grit:fields", "TaskService field on App struct"],
                      ["app.go", "/* grit:constructor-params */", "Constructor parameter"],
                      ["app.go", "/* grit:constructor-assign */", "Field assignment"],
                      ["app.go", "// grit:methods", "7 bound methods (CRUD + export)"],
                      ["types.go", "// grit:input-types", "TaskInput struct"],
                      ["cmd/studio/main.go", "// grit:studio-models", "Task model in Studio"],
                      ["sidebar.tsx", "// grit:nav-icons + nav", "Nav icon + sidebar item"],
                    ].map(([file, marker, what]) => (
                      <tr key={marker}>
                        <td className="p-3 font-mono text-xs">{file}</td>
                        <td className="p-3 font-mono text-xs">{marker}</td>
                        <td className="p-3">{what}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <div className="prose-grit mb-4">
                <p>
                  Here is the generated Go model:
                </p>
              </div>
              <CodeBlock
                language="go"
                filename="internal/models/task.go"
                code={`type Task struct {
    ID          uint           \`gorm:"primaryKey" json:"id"\`
    Title       string         \`json:"title"\`
    Description string         \`json:"description"\`
    Priority    string         \`json:"priority"\`
    Completed   bool           \`json:"completed"\`
    DueDate     time.Time      \`json:"due_date"\`
    CreatedAt   time.Time      \`json:"created_at"\`
    UpdatedAt   time.Time      \`json:"updated_at"\`
    DeletedAt   gorm.DeletedAt \`gorm:"index" json:"deleted_at"\`
}`}
                className="mb-0"
              />

              {/* Understanding the Generated Routes */}
              <div className="prose-grit mb-4 mt-8">
                <h3>Understanding the Generated Route Files</h3>
                <p>
                  The three frontend files are{" "}
                  <a href="https://tanstack.com/router" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">TanStack Router</a>{" "}
                  route files. Each file exports a <code>Route</code> constant
                  created with <code>createFileRoute()</code>. The TanStack Router
                  Vite plugin auto-discovers these files — no import or route
                  registry update is needed.
                </p>
              </div>

              <CodeBlock
                language="tsx"
                filename="frontend/src/routes/_layout/tasks.index.tsx"
                code={`import { createFileRoute } from "@tanstack/react-router";

// This line registers the route at /_layout/tasks/
export const Route = createFileRoute("/_layout/tasks/")({
  component: TasksPage,
});

function TasksPage() {
  // DataTable with search, pagination, PDF/Excel export
  // Calls ListTasks() via Wails bindings
}`}
                className="mb-4"
              />

              <div className="prose-grit mb-4">
                <p>
                  The edit route uses <code>Route.useParams()</code> for type-safe
                  parameter access. The <code>$id</code> in the filename becomes a
                  typed parameter:
                </p>
              </div>

              <CodeBlock
                language="tsx"
                filename="frontend/src/routes/_layout/tasks.$id.edit.tsx"
                code={`import { createFileRoute, useNavigate } from "@tanstack/react-router";

export const Route = createFileRoute("/_layout/tasks/$id/edit")({
  component: EditTaskPage,
});

function EditTaskPage() {
  const { id } = Route.useParams(); // typed string
  const navigate = useNavigate();

  // Fetch task by ID, populate form
  // On save: navigate({ to: "/tasks" })
}`}
                className="mb-4"
              />

              <div className="prose-grit mb-0">
                <p>
                  Navigation uses TanStack Router&apos;s object syntax:{" "}
                  <code>{"navigate({ to: \"/tasks\" })"}</code> instead of{" "}
                  <code>{"navigate(\"/tasks\")"}</code>. All routes, params, and
                  navigation calls are validated by TypeScript at compile time.
                </p>
              </div>
            </Step>

            {/* Step 6: Test the Task Manager */}
            <Step title="Test the Task Manager">
              <div className="prose-grit mb-4">
                <p>
                  Restart the development server to pick up the new Go code. If{" "}
                  <code>wails dev</code> is still running, it automatically
                  detects Go file changes and rebuilds. Otherwise, restart it:
                </p>
              </div>
              <CodeBlock terminal code="wails dev" className="mb-6" />

              <div className="prose-grit mb-4">
                <p>
                  Once the app reloads, you will see a new{" "}
                  <strong>Tasks</strong> item in the sidebar. Click it to open
                  the task list page. Try the following:
                </p>
              </div>

              <div className="space-y-2 mb-6">
                {[
                  {
                    action: "Create a task",
                    detail:
                      'Click the "New Task" button. Fill in the title, description, priority, due date, and toggle the completed switch. Hit save.',
                  },
                  {
                    action: "Search tasks",
                    detail:
                      "Type in the search box at the top of the list. It filters across all text fields in real time.",
                  },
                  {
                    action: "Paginate",
                    detail:
                      "Create a few more tasks. The table paginates automatically with page size controls.",
                  },
                  {
                    action: "Edit a task",
                    detail:
                      "Click the edit icon on any row. The form pre-fills with the existing values.",
                  },
                  {
                    action: "Delete a task",
                    detail:
                      "Click the delete icon. A confirmation dialog appears before the soft delete.",
                  },
                  {
                    action: "Export to PDF",
                    detail:
                      'Click the "PDF" export button to download a formatted PDF of all tasks.',
                  },
                  {
                    action: "Export to Excel",
                    detail:
                      'Click the "Excel" export button to download an .xlsx spreadsheet.',
                  },
                ].map((item) => (
                  <div
                    key={item.action}
                    className="flex items-start gap-3 rounded-lg border border-border/30 bg-card/30 px-4 py-2.5"
                  >
                    <span className="text-sm font-semibold text-foreground shrink-0 w-32">
                      {item.action}
                    </span>
                    <span className="text-sm text-muted-foreground/60">
                      {item.detail}
                    </span>
                  </div>
                ))}
              </div>

              <div className="prose-grit mb-0">
                <p>
                  All of this was generated from a single CLI command. The list
                  page, form, service methods, Wails bindings, routes, and
                  sidebar navigation are all wired up automatically.
                </p>
              </div>
            </Step>

            {/* Step 7: Generate Another Resource */}
            <Step title="Generate Another Resource — Categories">
              <div className="prose-grit mb-4">
                <p>
                  To demonstrate how easy it is to add more resources, generate a
                  Category resource for organizing tasks:
                </p>
              </div>
              <CodeBlock
                terminal
                code={`grit generate resource Category --fields "name:string,color:string"`}
                className="mb-6 glow-purple-sm"
              />

              <div className="prose-grit mb-4">
                <p>
                  This creates 5 more files and injects into the same 10
                  locations. Your app now has:
                </p>
              </div>

              <div className="overflow-x-auto mb-4">
                <table className="w-full text-sm border border-border rounded-lg">
                  <thead>
                    <tr className="border-b border-border bg-muted/50">
                      <th className="text-left p-3 font-medium">Resource</th>
                      <th className="text-left p-3 font-medium">Source</th>
                      <th className="text-left p-3 font-medium">Fields</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border text-muted-foreground">
                    <tr>
                      <td className="p-3 font-semibold text-foreground">Blog</td>
                      <td className="p-3">Built-in (scaffold)</td>
                      <td className="p-3 font-mono text-xs">
                        title, slug, content, image, published
                      </td>
                    </tr>
                    <tr>
                      <td className="p-3 font-semibold text-foreground">
                        Contact
                      </td>
                      <td className="p-3">Built-in (scaffold)</td>
                      <td className="p-3 font-mono text-xs">
                        name, email, phone, message
                      </td>
                    </tr>
                    <tr>
                      <td className="p-3 font-semibold text-foreground">Task</td>
                      <td className="p-3">Generated (Step 5)</td>
                      <td className="p-3 font-mono text-xs">
                        title, description, priority, completed, due_date
                      </td>
                    </tr>
                    <tr>
                      <td className="p-3 font-semibold text-foreground">
                        Category
                      </td>
                      <td className="p-3">Generated (Step 7)</td>
                      <td className="p-3 font-mono text-xs">name, color</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <div className="prose-grit mb-0">
                <p>
                  Each resource gets its own sidebar entry, list page, form page,
                  Go model, and service. You can keep generating as many
                  resources as your app needs.
                </p>
              </div>
            </Step>

            {/* Step 8: Open GORM Studio */}
            <Step title="Open GORM Studio">
              <div className="prose-grit mb-4">
                <p>
                  GORM Studio is a visual database browser bundled with every
                  Grit desktop project. Open a separate terminal in the{" "}
                  <code>taskmanager</code> directory and run:
                </p>
              </div>
              <CodeBlock
                terminal
                code="grit studio"
                className="mb-4 glow-purple-sm"
              />

              <div className="prose-grit mb-4">
                <p>
                  Your browser opens automatically at{" "}
                  <code>http://localhost:8080/studio</code>. You can:
                </p>
              </div>

              <div className="space-y-2 mb-6">
                {[
                  "Browse all tables: users, blogs, contacts, tasks, categories",
                  "View and edit records directly in the browser",
                  "Run raw SQL queries against the SQLite database",
                  "Inspect table schemas, column types, and indexes",
                  "Export query results for debugging",
                ].map((item, i) => (
                  <div
                    key={i}
                    className="flex items-start gap-3 rounded-lg border border-border/30 bg-card/30 px-4 py-2.5"
                  >
                    <span className="text-sm font-mono text-primary/60 shrink-0 mt-0.5">
                      {i + 1}.
                    </span>
                    <span className="text-sm text-muted-foreground">
                      {item}
                    </span>
                  </div>
                ))}
              </div>

              <div className="prose-grit mb-0">
                <p>
                  Studio runs as a standalone Go process that connects to the
                  same SQLite file your desktop app uses. This is useful for
                  verifying data, debugging issues, and understanding the
                  database schema during development.
                </p>
              </div>
            </Step>

            {/* Step 9: Build for Distribution */}
            <Step title="Build for Distribution">
              <div className="prose-grit mb-4">
                <p>
                  When you are ready to ship, compile the app into a native
                  executable:
                </p>
              </div>
              <CodeBlock
                terminal
                code="grit compile"
                className="mb-4 glow-purple-sm"
              />

              <div className="prose-grit mb-4">
                <p>
                  This runs <code>wails build</code> under the hood. The output
                  binary is placed in <code>build/bin/</code>:
                </p>
              </div>

              <div className="overflow-x-auto mb-6">
                <table className="w-full text-sm border border-border rounded-lg">
                  <thead>
                    <tr className="border-b border-border bg-muted/50">
                      <th className="text-left p-3 font-medium">Platform</th>
                      <th className="text-left p-3 font-medium">Output</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border text-muted-foreground">
                    <tr>
                      <td className="p-3 font-medium text-foreground">
                        Windows
                      </td>
                      <td className="p-3 font-mono text-xs">
                        build/bin/taskmanager.exe
                      </td>
                    </tr>
                    <tr>
                      <td className="p-3 font-medium text-foreground">macOS</td>
                      <td className="p-3 font-mono text-xs">
                        build/bin/taskmanager.app
                      </td>
                    </tr>
                    <tr>
                      <td className="p-3 font-medium text-foreground">Linux</td>
                      <td className="p-3 font-mono text-xs">
                        build/bin/taskmanager
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <div className="prose-grit mb-4">
                <p>
                  The resulting binary is a single executable with everything
                  embedded:
                </p>
              </div>

              <div className="space-y-3 mb-4">
                {[
                  {
                    label: "Embedded frontend",
                    desc: "The entire React app is compiled and embedded into the Go binary via go:embed. No separate files to distribute.",
                  },
                  {
                    label: "SQLite database",
                    desc: "The database file (app.db) is created at runtime in the working directory. Ship the binary alone.",
                  },
                  {
                    label: "No runtime dependencies",
                    desc: "Users do not need Go, Node.js, or any other tooling. Just double-click the executable.",
                  },
                ].map((item) => (
                  <div
                    key={item.label}
                    className="rounded-lg border border-border/30 bg-card/30 px-4 py-3"
                  >
                    <span className="text-[15px] font-semibold block mb-1">
                      {item.label}
                    </span>
                    <span className="text-sm text-muted-foreground/60">
                      {item.desc}
                    </span>
                  </div>
                ))}
              </div>

              <div className="prose-grit mb-0">
                <blockquote>
                  For a distributable Windows installer, run{" "}
                  <code>grit package</code> — it builds an NSIS installer (.exe)
                  and auto-detects <code>makensis</code> on your PATH. See the{" "}
                  <Link
                    href="/docs/desktop/building"
                    className="text-primary hover:underline"
                  >
                    Building & Distribution
                  </Link>{" "}
                  guide for cross-platform builds and distribution details.
                </blockquote>
              </div>
            </Step>

            {/* Step 10: What's Next */}
            <Step title="What's Next">
              <div className="prose-grit mb-6">
                <p>
                  You have built a fully functional Task Manager desktop
                  application with two custom resources, authentication, and
                  database browsing — all from a few CLI commands. Here is where
                  to go from here:
                </p>
              </div>

              <div className="grid gap-3 sm:grid-cols-2 mb-6">
                {[
                  {
                    title: "Resource Generation",
                    desc: "Deep dive into field types, slug generation, belongs_to relationships, and the remove command.",
                    href: "/docs/desktop/resource-generation",
                  },
                  {
                    title: "Building & Distribution",
                    desc: "Cross-platform builds, NSIS installers, app icons, environment config, and production tips.",
                    href: "/docs/desktop/building",
                  },
                  {
                    title: "Build a POS App",
                    desc: "Follow the next tutorial to build a Point of Sale desktop application with products, orders, and receipts.",
                    href: "/docs/desktop/pos-app",
                  },
                  {
                    title: "Desktop LLM Reference",
                    desc: "The complete Grit desktop reference for AI assistants — every convention, pattern, and code marker.",
                    href: "/docs/ai-skill/llm-guide",
                  },
                ].map((card) => (
                  <Link
                    key={card.href}
                    href={card.href}
                    className="group rounded-lg border border-border/30 bg-card/30 px-4 py-4 hover:border-primary/30 hover:bg-card/50 transition-colors"
                  >
                    <h3 className="text-[15px] font-semibold mb-1 group-hover:text-primary transition-colors">
                      {card.title}
                    </h3>
                    <p className="text-sm text-muted-foreground/60">
                      {card.desc}
                    </p>
                  </Link>
                ))}
              </div>
            </Step>
            </Steps>

            {/* Desktop Handbook */}
            <div className="mb-8">
              <a
                href="https://14j7oh8kso.ufs.sh/f/HLxTbDBCDLwfpiJDPD3QgNG9hYzVFo5iLR0yrDPTJedWnBH7"
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-3 rounded-lg border border-primary/25 bg-primary/5 px-5 py-4 hover:bg-primary/10 transition-colors group"
              >
                <Download className="h-5 w-5 text-primary/70 group-hover:text-primary shrink-0" />
                <div>
                  <span className="text-sm font-semibold text-foreground block">Download Desktop Handbook (PDF)</span>
                  <span className="text-xs text-muted-foreground/60">Complete offline reference for Grit Desktop development</span>
                </div>
              </a>
            </div>

            {/* Bottom Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="text-muted-foreground/60 hover:text-foreground"
              >
                <Link href="/docs/desktop/getting-started" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Getting Started
                </Link>
              </Button>
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="text-muted-foreground/60 hover:text-foreground"
              >
                <Link href="/docs/desktop/pos-app" className="gap-1.5">
                  Build a POS App
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
