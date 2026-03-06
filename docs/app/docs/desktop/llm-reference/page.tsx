import Link from "next/link";
import { ArrowLeft, Download } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { CodeBlock } from "@/components/code-block";
import { getDocMetadata } from "@/config/docs-metadata";

export const metadata = getDocMetadata("/docs/desktop/llm-reference");

// ── reusable mini-components ────────────────────────────────────
function SectionNum({ n }: { n: number }) {
  return (
    <span className="inline-flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/20 text-xs font-mono font-semibold text-primary mr-2">
      {n}
    </span>
  );
}

function Warn({ children }: { children: React.ReactNode }) {
  return (
    <div className="p-4 rounded-lg border border-red-500/20 bg-red-500/5 mb-4">
      <p className="text-sm text-red-400/80 leading-relaxed font-medium">
        {children}
      </p>
    </div>
  );
}

function Tip({ children }: { children: React.ReactNode }) {
  return (
    <div className="p-4 rounded-lg border border-primary/20 bg-primary/5 mb-4">
      <p className="text-sm text-muted-foreground leading-relaxed">
        {children}
      </p>
    </div>
  );
}

function Note({ children }: { children: React.ReactNode }) {
  return (
    <div className="p-4 rounded-lg border border-amber-500/20 bg-amber-500/5 mb-4">
      <p className="text-sm text-amber-400/80 leading-relaxed">{children}</p>
    </div>
  );
}

export default function DesktopLLMReferencePage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* ── Header ── */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">
                AI Reference
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Grit Desktop — Complete LLM Reference
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                The canonical guide for AI assistants building native desktop
                applications with Grit Desktop. Covers architecture, every CLI
                command, resource generation, field types, Wails bindings,
                DataTable, FormBuilder, code markers, building executables, and
                the rules that must never be broken.
              </p>
              <Tip>
                <strong className="text-primary/90">For AI tools:</strong> Read
                this entire page before generating any code for Grit Desktop
                projects. Every convention here is enforced by the CLI —
                deviating breaks{" "}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                  grit generate
                </code>{" "}
                and{" "}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                  grit remove
                </code>
                . Desktop uses Wails bindings, not HTTP — there is no REST API.
              </Tip>
            </div>

            <div className="prose-grit">
              {/* ════════════════════════════════════════════════════
                  1. What is Grit Desktop
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={1} />
                  What is Grit Desktop?
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit Desktop builds{" "}
                  <strong className="text-foreground/80">
                    native desktop applications
                  </strong>{" "}
                  using Go + Wails + React. One command (
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    grit new-desktop myapp
                  </code>
                  ) scaffolds a complete project that compiles to a{" "}
                  <strong className="text-foreground/80">
                    single native binary
                  </strong>
                  . No browser, no Electron, no separate runtime.
                </p>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 mb-6">
                  {[
                    {
                      label: "Backend",
                      value:
                        "Go (Wails v2 bindings) — direct function calls, no HTTP server",
                    },
                    {
                      label: "Frontend",
                      value:
                        "Vite + React + TanStack Router (file-based) + TanStack Query + Tailwind CSS",
                    },
                    {
                      label: "Database",
                      value:
                        "SQLite (default, local) via GORM — fully offline-first",
                    },
                    {
                      label: "Distribution",
                      value:
                        "Single binary (~10-15MB) for Windows, macOS, and Linux",
                    },
                    {
                      label: "Code Generation",
                      value:
                        "grit generate resource — same CLI, auto-detects desktop",
                    },
                    {
                      label: "Dev Tools",
                      value:
                        "GORM Studio (localhost:8080/studio), hot-reload, PDF/Excel export",
                    },
                  ].map((item) => (
                    <div
                      key={item.label}
                      className="p-3 rounded-lg border border-border/30 bg-card/30"
                    >
                      <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-1">
                        {item.label}
                      </p>
                      <p className="text-xs font-mono text-foreground/70">
                        {item.value}
                      </p>
                    </div>
                  ))}
                </div>
                <p className="text-muted-foreground leading-relaxed">
                  The scaffolded project includes JWT authentication, blog and
                  contact CRUD resources, PDF and Excel data export, a custom
                  frameless title bar, dark theme with shadcn/ui components,
                  toast notifications, and GORM Studio for database browsing —
                  all out of the box.
                </p>
              </div>

              {/* ════════════════════════════════════════════════════
                  2. Why Wails
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={2} />
                  Why Wails (not Electron or Tauri)?
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit chose{" "}
                  <a
                    href="https://wails.io"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-primary hover:underline"
                  >
                    Wails
                  </a>{" "}
                  because Go is already the backend language for Grit web
                  projects. Using Wails means developers write Go for both web
                  and desktop — no Rust learning curve (Tauri), and none of
                  Electron&apos;s bloat.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/50 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Aspect
                        </th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-primary/80 uppercase">
                          Wails
                        </th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Electron
                        </th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Tauri
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/20 text-xs">
                      {[
                        [
                          "Backend language",
                          "Go",
                          "JavaScript/C++",
                          "Rust",
                        ],
                        [
                          "Binary size",
                          "~10-15MB",
                          "~150MB+",
                          "~5-10MB",
                        ],
                        [
                          "RAM usage",
                          "~30-50MB",
                          "~200MB+",
                          "~30-50MB",
                        ],
                        [
                          "Frontend",
                          "Any web framework",
                          "Any web framework",
                          "Any web framework",
                        ],
                        [
                          "Backend calls",
                          "Direct Go bindings",
                          "IPC (serialization overhead)",
                          "Rust commands (compile-time checked)",
                        ],
                        [
                          "Learning curve",
                          "Go + React",
                          "JavaScript only",
                          "Rust + JavaScript",
                        ],
                        ["Startup time", "<1s", "2-5s", "<1s"],
                        [
                          "Cross-platform",
                          "Win / macOS / Linux",
                          "Win / macOS / Linux",
                          "Win / macOS / Linux",
                        ],
                      ].map(([aspect, wails, electron, tauri]) => (
                        <tr key={aspect}>
                          <td className="px-4 py-2 font-medium text-foreground/70">
                            {aspect}
                          </td>
                          <td className="px-4 py-2 text-primary/80 font-mono">
                            {wails}
                          </td>
                          <td className="px-4 py-2 text-muted-foreground/60 font-mono">
                            {electron}
                          </td>
                          <td className="px-4 py-2 text-muted-foreground/60 font-mono">
                            {tauri}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                <div className="space-y-2 mb-4">
                  {[
                    "Go is the backend language — Wails is the natural choice (no Rust learning curve)",
                    "Direct function calls from React to Go — no serialization overhead, no IPC",
                    "Small binary size (~10-15MB vs Electron's 150MB+)",
                    "Memory efficient (~30-50MB vs Electron's 200MB+)",
                    "Same language for CLI, web API, and desktop — unified developer experience",
                  ].map((reason) => (
                    <div
                      key={reason}
                      className="flex items-start gap-2 text-sm text-muted-foreground/70"
                    >
                      <span className="text-primary/60 mt-0.5 shrink-0">
                        --
                      </span>
                      <span>{reason}</span>
                    </div>
                  ))}
                </div>
              </div>

              {/* ════════════════════════════════════════════════════
                  3. Architecture
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={3} />
                  Architecture
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Unlike Grit web projects (which use Gin HTTP server), desktop
                  projects have{" "}
                  <strong className="text-foreground/80">no HTTP server</strong>.
                  Wails bridges Go and JavaScript directly. The React frontend
                  calls Go functions through generated TypeScript bindings.
                </p>

                <Warn>
                  There is no REST API in desktop projects. React calls Go
                  functions directly via{" "}
                  <code className="font-mono bg-red-500/10 px-1 rounded">
                    window.go.main.App.MethodName()
                  </code>
                  . Never suggest HTTP endpoints, fetch calls, or Axios for
                  desktop projects.
                </Warn>

                <h3 className="text-base font-semibold text-foreground/80 mb-3">
                  Call Flow
                </h3>
                <CodeBlock
                  language="bash"
                  filename="Desktop call flow"
                  code={`React Component
  --> Wails TypeScript Binding (auto-generated)
    --> Go App Method (on the App struct)
      --> GORM Service (business logic + queries)
        --> SQLite Database`}
                />

                <div className="space-y-3 mt-6 mb-4">
                  {[
                    {
                      step: "1",
                      title: "App Struct (app.go)",
                      desc: "A Go struct with exported methods. Any method on this struct with exported name and supported types is automatically available to the React frontend. Wails generates TypeScript bindings during build.",
                    },
                    {
                      step: "2",
                      title: "React Frontend (frontend/)",
                      desc: "A Vite-powered React app with TanStack Router (file-based routing) and TanStack Query for state management. Calls Go functions via the generated Wails bindings — no fetch, no Axios, no HTTP.",
                    },
                    {
                      step: "3",
                      title: "GORM Services (internal/service/)",
                      desc: "Business logic and database queries. Each resource has a service with List, ListAll, GetByID, Create, Update, Delete methods. Services operate on GORM models backed by SQLite.",
                    },
                    {
                      step: "4",
                      title: "SQLite Database",
                      desc: "Local file-based database (app.db). All data lives on the user's machine. Fully offline — no network connection required.",
                    },
                  ].map((item) => (
                    <div
                      key={item.step}
                      className="flex gap-3 p-3 rounded-lg border border-border/30 bg-card/30"
                    >
                      <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                        {item.step}
                      </div>
                      <div>
                        <p className="text-sm font-medium">{item.title}</p>
                        <p className="text-xs text-muted-foreground/70 mt-0.5">
                          {item.desc}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* ════════════════════════════════════════════════════
                  4. Project Structure
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={4} />
                  Project Structure
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Desktop projects are a{" "}
                  <strong className="text-foreground/80">
                    single directory
                  </strong>{" "}
                  (not a monorepo). Go code lives at the root and in{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    internal/
                  </code>
                  , while the React frontend lives in{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    frontend/
                  </code>
                  .
                </p>
                <CodeBlock
                  language="bash"
                  filename="myapp/ (scaffolded desktop project)"
                  code={`myapp/
├── main.go                  # Wails entry point — creates app, binds methods
├── app.go                   # App struct with all bound methods + constructor
├── types.go                 # Input structs for bound methods
├── wails.json               # Wails project configuration
├── go.mod                   # Go module definition
├── go.sum                   # Go dependency checksums
│
├── internal/
│   ├── config/
│   │   └── config.go        # App configuration (reads .env)
│   ├── db/
│   │   └── db.go            # GORM database setup (SQLite) + AutoMigrate
│   ├── models/
│   │   ├── user.go          # User model (auth)
│   │   ├── blog.go          # Blog post model (scaffolded example)
│   │   └── contact.go       # Contact model (scaffolded example)
│   └── service/
│       ├── auth.go          # Authentication service (register, login, JWT)
│       ├── blog.go          # Blog CRUD service
│       └── contact.go       # Contact CRUD service
│
├── frontend/
│   ├── src/
│   │   ├── main.tsx          # React entry point (TanStack Router)
│   │   ├── routes/           # File-based routes (TanStack Router)
│   │   │   ├── __root.tsx    # Root route
│   │   │   ├── _layout.tsx   # Auth guard + sidebar layout
│   │   │   └── _layout/      # Protected page routes
│   │   ├── components/       # Shared UI: sidebar, title-bar, data-table, etc.
│   │   │   └── sidebar.tsx   # App sidebar with navigation
│   │   ├── hooks/            # TanStack Query hooks for Wails bindings
│   │   └── lib/              # Utilities
│   ├── index.html
│   ├── package.json
│   ├── vite.config.ts
│   └── tailwind.config.js
│
├── build/
│   ├── appicon.png           # App icon (1024x1024 — auto-converted per platform)
│   └── bin/                  # Build output directory
│
└── cmd/
    └── studio/
        └── main.go           # GORM Studio server (port 8080)`}
                />

                <Note>
                  Key difference from web projects: Desktop has{" "}
                  <code className="font-mono bg-amber-500/10 px-1 rounded">
                    app.go
                  </code>{" "}
                  +{" "}
                  <code className="font-mono bg-amber-500/10 px-1 rounded">
                    types.go
                  </code>{" "}
                  at the root instead of handlers/ and routes/. There is no
                  apps/ directory, no Turborepo, no monorepo structure.
                </Note>
              </div>

              {/* ════════════════════════════════════════════════════
                  5. CLI Commands
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={5} />
                  All CLI Commands
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every desktop-relevant command. The CLI auto-detects desktop
                  projects by checking for{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    wails.json
                  </code>
                  .
                </p>

                <div className="rounded-xl border border-border/40 bg-card/50 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Command
                        </th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Description
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/20 text-xs">
                      {[
                        [
                          "grit new-desktop <name>",
                          "Scaffold a complete desktop project with Go backend, React frontend, auth, CRUD, GORM Studio, and dark theme",
                        ],
                        [
                          'grit generate resource <Name> --fields "..."',
                          "Generate a full-stack CRUD resource: Go model, service, React list page, React form page, plus 10 code injections",
                        ],
                        [
                          "grit remove resource <Name>",
                          "Remove a generated resource: deletes files and reverses all 10 injections. Markers stay intact.",
                        ],
                        [
                          "grit start",
                          "Start development mode (runs wails dev). Opens native window with hot-reload for Go and React.",
                        ],
                        [
                          "grit compile",
                          "Build native executable (runs wails build). Output in build/bin/.",
                        ],
                        [
                          "grit studio",
                          "Open GORM Studio at localhost:8080/studio — visual database browser.",
                        ],
                        [
                          "grit version",
                          "Show installed CLI version.",
                        ],
                      ].map(([cmd, desc]) => (
                        <tr key={cmd}>
                          <td className="px-4 py-2.5 font-mono text-primary/80 whitespace-nowrap">
                            {cmd}
                          </td>
                          <td className="px-4 py-2.5 text-muted-foreground/70">
                            {desc}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                <Warn>
                  Desktop projects do not have{" "}
                  <code className="font-mono bg-red-500/10 px-1 rounded">
                    grit start server
                  </code>{" "}
                  or{" "}
                  <code className="font-mono bg-red-500/10 px-1 rounded">
                    grit start client
                  </code>
                  . Those are web-only. For desktop, use{" "}
                  <code className="font-mono bg-red-500/10 px-1 rounded">
                    grit start
                  </code>{" "}
                  or{" "}
                  <code className="font-mono bg-red-500/10 px-1 rounded">
                    wails dev
                  </code>{" "}
                  directly. There is no{" "}
                  <code className="font-mono bg-red-500/10 px-1 rounded">
                    grit migrate
                  </code>{" "}
                  or{" "}
                  <code className="font-mono bg-red-500/10 px-1 rounded">
                    grit seed
                  </code>{" "}
                  for desktop — GORM AutoMigrate runs on startup automatically.
                </Warn>
              </div>

              {/* ════════════════════════════════════════════════════
                  6. Resource Generation
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={6} />
                  Resource Generation — Complete Reference
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The same{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    grit generate resource
                  </code>{" "}
                  command works for both web and desktop projects. The CLI
                  auto-detects the project type by checking for{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    wails.json
                  </code>
                  .
                </p>

                <h3 className="text-base font-semibold text-foreground/80 mb-3">
                  Syntax
                </h3>
                <CodeBlock
                  terminal
                  code={`grit generate resource <Name> --fields "field1:type,field2:type,field3:type"`}
                />

                <h3 className="text-base font-semibold text-foreground/80 mb-3 mt-8">
                  Supported Field Types
                </h3>
                <div className="rounded-xl border border-border/40 bg-card/50 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Type
                        </th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Go Type
                        </th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Form Input
                        </th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Notes
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/20 text-xs font-mono">
                      {[
                        {
                          type: "string",
                          go: "string",
                          form: "Text input",
                          note: "Short text — titles, names, emails",
                        },
                        {
                          type: "text",
                          go: "string",
                          form: "Textarea",
                          note: "Long text — descriptions, notes",
                        },
                        {
                          type: "richtext",
                          go: "string",
                          form: "Textarea",
                          note: "Rich content (rendered as textarea in desktop)",
                        },
                        {
                          type: "int",
                          go: "int",
                          form: "Number input",
                          note: "Signed integer — stock, quantity",
                        },
                        {
                          type: "uint",
                          go: "uint",
                          form: "Number input",
                          note: "Unsigned integer — IDs, counts",
                        },
                        {
                          type: "float",
                          go: "float64",
                          form: "Number input (decimal)",
                          note: "Prices, coordinates, ratings",
                        },
                        {
                          type: "bool",
                          go: "bool",
                          form: "Toggle switch",
                          note: "Active, published, featured flags",
                        },
                        {
                          type: "date",
                          go: "time.Time",
                          form: "Date picker",
                          note: "Date only — due dates, birthdays",
                        },
                        {
                          type: "datetime",
                          go: "time.Time",
                          form: "DateTime picker",
                          note: "Date + time — event timestamps",
                        },
                        {
                          type: "slug",
                          go: "string (uniqueIndex)",
                          form: "Auto-generated (excluded from form)",
                          note: "URL-friendly slug from source field",
                        },
                        {
                          type: "belongs_to",
                          go: "uint (foreign key)",
                          form: "Number input (FK ID)",
                          note: "Foreign key to parent model",
                        },
                      ].map((row) => (
                        <tr key={row.type}>
                          <td className="px-4 py-2 text-primary/80">
                            {row.type}
                          </td>
                          <td className="px-4 py-2 text-foreground/60">
                            {row.go}
                          </td>
                          <td className="px-4 py-2 text-muted-foreground/60 font-sans">
                            {row.form}
                          </td>
                          <td className="px-4 py-2 text-muted-foreground/40 font-sans">
                            {row.note}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                <h3 className="text-base font-semibold text-foreground/80 mb-3">
                  Slug Field Syntax
                </h3>
                <p className="text-sm text-muted-foreground/70 mb-3">
                  The{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    slug
                  </code>{" "}
                  type uses a special syntax to specify the source field:
                </p>
                <CodeBlock
                  terminal
                  code={`grit generate resource Article --fields "title:string,slug:slug:source=title,content:richtext"`}
                />
                <p className="text-sm text-muted-foreground/60 mt-2 mb-6">
                  This adds a{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    uniqueIndex
                  </code>{" "}
                  and a GORM{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    BeforeCreate
                  </code>{" "}
                  hook that auto-generates the slug from the title. Slug fields
                  are excluded from forms and input structs.
                </p>

                <h3 className="text-base font-semibold text-foreground/80 mb-3">
                  Files Created Per Resource
                </h3>
                <div className="space-y-3 mb-6">
                  {[
                    {
                      file: "internal/models/<snake>.go",
                      desc: "GORM model struct with ID, fields, timestamps, and soft delete. Includes slug hook if applicable.",
                    },
                    {
                      file: "internal/service/<snake>.go",
                      desc: "Service with List (paginated), ListAll, GetByID, Create, Update, Delete methods. Returns PaginatedResult for list queries.",
                    },
                    {
                      file: "frontend/src/routes/_layout/<plural>.index.tsx",
                      desc: "List route — DataTable, search, pagination, PDF/Excel export",
                    },
                    {
                      file: "frontend/src/routes/_layout/<plural>.new.tsx",
                      desc: "Create form route — field-type-based inputs, validation, toast",
                    },
                    {
                      file: "frontend/src/routes/_layout/<plural>.$id.edit.tsx",
                      desc: "Edit form route — pre-fills from GetByID, type-safe params",
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

                <h3 className="text-base font-semibold text-foreground/80 mb-3">
                  Example: Generate a Product Resource
                </h3>
                <CodeBlock
                  terminal
                  code={`grit generate resource Product --fields "name:string,price:float,stock:int,published:bool"`}
                />

                <div className="mt-6 mb-4">
                  <h3 className="text-base font-semibold text-foreground/80 mb-3">
                    Generated Model
                  </h3>
                  <CodeBlock
                    language="go"
                    filename="internal/models/product.go"
                    code={`package models

import (
    "gorm.io/gorm"
)

type Product struct {
    gorm.Model
    Name      string  \`gorm:"not null" json:"name"\`
    Price     float64 \`json:"price"\`
    Stock     int     \`json:"stock"\`
    Published bool    \`gorm:"default:false" json:"published"\`
}

type ProductInput struct {
    Name      string  \`json:"name"\`
    Price     float64 \`json:"price"\`
    Stock     int     \`json:"stock"\`
    Published bool    \`json:"published"\`
}`}
                  />
                </div>

                <div className="mt-6 mb-4">
                  <h3 className="text-base font-semibold text-foreground/80 mb-3">
                    Generated Service (Skeleton)
                  </h3>
                  <CodeBlock
                    language="go"
                    filename="internal/service/product.go"
                    code={`package service

import (
    "myapp/internal/models"
    "gorm.io/gorm"
)

type ProductService struct {
    DB *gorm.DB
}

type PaginatedResult struct {
    Data  interface{} \`json:"data"\`
    Total int64       \`json:"total"\`
    Page  int         \`json:"page"\`
    Pages int         \`json:"pages"\`
}

func (s *ProductService) List(page, pageSize int, search string) (*PaginatedResult, error) { ... }
func (s *ProductService) ListAll() ([]models.Product, error) { ... }
func (s *ProductService) GetByID(id uint) (*models.Product, error) { ... }
func (s *ProductService) Create(input models.ProductInput) (*models.Product, error) { ... }
func (s *ProductService) Update(id uint, input models.ProductInput) (*models.Product, error) { ... }
func (s *ProductService) Delete(id uint) error { ... }`}
                  />
                </div>
              </div>

              {/* ════════════════════════════════════════════════════
                  7. DataTable (List Page)
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={7} />
                  DataTable (List Page)
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every generated resource gets a list page with a full-featured
                  DataTable. Here is what it includes out of the box:
                </p>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 mb-6">
                  {[
                    {
                      name: "Search Bar",
                      desc: "Debounced text input that filters records. Searches across the primary text field.",
                    },
                    {
                      name: "Paginated Table",
                      desc: "Data table with configurable page size. Shows total records, page numbers, and navigation.",
                    },
                    {
                      name: "Sortable Columns",
                      desc: "Click column headers to sort ascending/descending. Columns are generated per field.",
                    },
                    {
                      name: "Edit Button",
                      desc: "Navigates to the form page with the record ID — pre-fills all fields for editing.",
                    },
                    {
                      name: "Delete Button",
                      desc: "Shows a confirmation dialog. On confirm, calls the Go delete method via Wails binding.",
                    },
                    {
                      name: "PDF Export",
                      desc: "Exports all records to a styled PDF document. Calls ExportProductsPDF() Go method.",
                    },
                    {
                      name: "Excel Export",
                      desc: "Exports all records to an .xlsx spreadsheet. Calls ExportProductsExcel() Go method.",
                    },
                    {
                      name: "New Button",
                      desc: "Navigates to the form page in create mode (no ID). Form fields are empty.",
                    },
                  ].map((feature) => (
                    <div
                      key={feature.name}
                      className="rounded-lg border border-border/30 bg-card/30 px-4 py-3"
                    >
                      <span className="text-[15px] font-semibold block mb-1">
                        {feature.name}
                      </span>
                      <span className="text-sm text-muted-foreground/60">
                        {feature.desc}
                      </span>
                    </div>
                  ))}
                </div>

                <h3 className="text-base font-semibold text-foreground/80 mb-3">
                  Customizing Columns
                </h3>
                <p className="text-sm text-muted-foreground/70 mb-3">
                  Columns are defined in the list page{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    index.tsx
                  </code>
                  . You can reorder, rename, add custom formatting, or add
                  computed columns:
                </p>
                <CodeBlock
                  language="tsx"
                  filename="frontend/src/routes/_layout/products.index.tsx (column definition)"
                  code={`const columns = [
  { accessorKey: "name", header: "Product Name" },
  {
    accessorKey: "price",
    header: "Price",
    cell: ({ row }: any) => "$" + Number(row.getValue("price")).toFixed(2),
  },
  { accessorKey: "stock", header: "Stock" },
  {
    accessorKey: "published",
    header: "Status",
    cell: ({ row }: any) =>
      row.getValue("published") ? "Published" : "Draft",
  },
];`}
                />
              </div>

              {/* ════════════════════════════════════════════════════
                  8. FormBuilder (Form Page)
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={8} />
                  FormBuilder (Form Page)
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every generated resource gets a form page that handles both
                  creating and editing records. The form fields are determined by
                  the field types declared during generation.
                </p>

                <h3 className="text-base font-semibold text-foreground/80 mb-3">
                  Field Type to Form Input Mapping
                </h3>
                <div className="rounded-xl border border-border/40 bg-card/50 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Field Type
                        </th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Form Input
                        </th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Behavior
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/20 text-xs">
                      {[
                        [
                          "string",
                          'Text input (<input type="text">)',
                          "Standard text field with label",
                        ],
                        [
                          "text / richtext",
                          "Textarea (<textarea>)",
                          "Multi-line text area, resizable",
                        ],
                        [
                          "int / uint / float",
                          'Number input (<input type="number">)',
                          "Numeric input with step control",
                        ],
                        [
                          "bool",
                          "Toggle switch",
                          "On/off toggle, defaults to false",
                        ],
                        [
                          "date",
                          "Date picker",
                          "Calendar date selector",
                        ],
                        [
                          "datetime",
                          "DateTime picker",
                          "Calendar + time selector",
                        ],
                        [
                          "slug",
                          "(Excluded from form)",
                          "Auto-generated in Go service from source field",
                        ],
                        [
                          "belongs_to",
                          "Number input (FK ID)",
                          "Enter the related record ID directly",
                        ],
                      ].map(([type_, input, behavior]) => (
                        <tr key={type_}>
                          <td className="px-4 py-2 text-primary/80 font-mono">
                            {type_}
                          </td>
                          <td className="px-4 py-2 text-foreground/60">
                            {input}
                          </td>
                          <td className="px-4 py-2 text-muted-foreground/60">
                            {behavior}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                <h3 className="text-base font-semibold text-foreground/80 mb-3">
                  Create vs Edit Mode
                </h3>
                <div className="space-y-3 mb-6">
                  {[
                    {
                      mode: "Create",
                      desc: "URL has no ID parameter. Form fields are empty. On submit, calls the Go Create method via Wails binding. Shows success toast and navigates to list page.",
                    },
                    {
                      mode: "Edit",
                      desc: "URL includes the record ID. On load, calls the Go Get method to fetch the record and pre-fill all fields. On submit, calls the Go Update method. Shows success toast and navigates back.",
                    },
                  ].map((item) => (
                    <div
                      key={item.mode}
                      className="rounded-lg border border-border/30 bg-card/30 px-4 py-3"
                    >
                      <span className="text-[15px] font-semibold block mb-1">
                        {item.mode} Mode
                      </span>
                      <span className="text-sm text-muted-foreground/60">
                        {item.desc}
                      </span>
                    </div>
                  ))}
                </div>

                <Tip>
                  Form validation happens on the Go side. If a GORM{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    not null
                  </code>{" "}
                  constraint fails or a service returns an error, the React form
                  displays a toast notification with the error message.
                </Tip>
              </div>

              {/* ════════════════════════════════════════════════════
                  9. Wails Bound Methods
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={9} />
                  Wails Bound Methods
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Go methods on the{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    App
                  </code>{" "}
                  struct are automatically available in React. Wails generates
                  TypeScript bindings during{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    wails dev
                  </code>{" "}
                  and{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    wails build
                  </code>
                  . React calls them as async functions.
                </p>

                <h3 className="text-base font-semibold text-foreground/80 mb-3">
                  7 Methods Generated Per Resource
                </h3>
                <CodeBlock
                  language="go"
                  filename="app.go — methods generated for Product resource"
                  code={`// List with pagination and search
func (a *App) GetProducts(page, pageSize int, search string) (*service.PaginatedResult, error)

// Get single record by ID
func (a *App) GetProduct(id uint) (*models.Product, error)

// Create new record
func (a *App) CreateProduct(input models.ProductInput) (*models.Product, error)

// Update existing record
func (a *App) UpdateProduct(id uint, input models.ProductInput) (*models.Product, error)

// Delete record (soft delete)
func (a *App) DeleteProduct(id uint) error

// Export all records as PDF (returns raw bytes)
func (a *App) ExportProductsPDF() ([]byte, error)

// Export all records as Excel (returns raw bytes)
func (a *App) ExportProductsExcel() ([]byte, error)`}
                />

                <h3 className="text-base font-semibold text-foreground/80 mb-3 mt-8">
                  Calling from React
                </h3>
                <p className="text-sm text-muted-foreground/70 mb-3">
                  Wails generates TypeScript bindings in{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    frontend/wailsjs/
                  </code>
                  . Import and call them as async functions:
                </p>
                <CodeBlock
                  language="tsx"
                  filename="frontend/src/routes/_layout/products.index.tsx (calling Go from React)"
                  code={`import { GetProducts, DeleteProduct, ExportProductsPDF, ExportProductsExcel } from "../../wailsjs/go/main/App";

// In a component:
const result = await GetProducts(page, pageSize, searchQuery);
// result = { data: Product[], total: number, page: number, pages: number }

await DeleteProduct(productId);

const pdfBytes = await ExportProductsPDF();
// Save pdfBytes as a file using browser APIs

const excelBytes = await ExportProductsExcel();
// Save excelBytes as .xlsx file`}
                />

                <Note>
                  The React frontend calls{" "}
                  <code className="font-mono bg-amber-500/10 px-1 rounded">
                    GetProducts(...)
                  </code>{" "}
                  — NOT{" "}
                  <code className="font-mono bg-amber-500/10 px-1 rounded">
                    fetch(&quot;/api/products&quot;)
                  </code>
                  . There is no HTTP involved. The call goes directly from
                  JavaScript to Go through the Wails bridge.
                </Note>
              </div>

              {/* ════════════════════════════════════════════════════
                  10. Building into Executable
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={10} />
                  Building into Executable
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Compile the entire application — Go runtime, React UI, static
                  assets, and SQLite — into a single native binary.
                </p>

                <CodeBlock
                  terminal
                  code="grit compile"
                  className="mb-4 glow-purple-sm"
                />

                <p className="text-sm text-muted-foreground/70 mb-6">
                  This runs{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    wails build
                  </code>{" "}
                  under the hood. The React frontend is embedded via{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    //go:embed all:frontend/dist
                  </code>{" "}
                  — no separate files need to be distributed.
                </p>

                <h3 className="text-base font-semibold text-foreground/80 mb-3">
                  Build Output
                </h3>
                <div className="rounded-xl border border-border/40 bg-card/50 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Platform
                        </th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Output Path
                        </th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Size
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/20 text-xs">
                      {[
                        [
                          "Windows",
                          "build/bin/myapp.exe",
                          "~10-15MB",
                        ],
                        [
                          "macOS",
                          "build/bin/myapp.app",
                          "~10-15MB",
                        ],
                        [
                          "Linux",
                          "build/bin/myapp",
                          "~10-15MB",
                        ],
                      ].map(([platform, output, size]) => (
                        <tr key={platform}>
                          <td className="px-4 py-2.5 font-medium text-foreground/70">
                            {platform}
                          </td>
                          <td className="px-4 py-2.5 font-mono text-primary/80">
                            {output}
                          </td>
                          <td className="px-4 py-2.5 text-muted-foreground/60">
                            {size}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                <h3 className="text-base font-semibold text-foreground/80 mb-3">
                  Cross-Platform Builds
                </h3>
                <CodeBlock
                  terminal
                  code={`# Windows
wails build -platform windows/amd64

# macOS (Intel)
wails build -platform darwin/amd64

# macOS (Apple Silicon)
wails build -platform darwin/arm64

# Linux
wails build -platform linux/amd64`}
                />

                <h3 className="text-base font-semibold text-foreground/80 mb-3 mt-6">
                  Windows Installer (NSIS)
                </h3>
                <p className="text-sm text-muted-foreground/70 mb-3">
                  Create a Windows installer with Start Menu shortcuts and
                  uninstallation support. Requires{" "}
                  <a
                    href="https://nsis.sourceforge.io/"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-primary hover:underline"
                  >
                    NSIS
                  </a>{" "}
                  installed:
                </p>
                <CodeBlock terminal code="wails build -nsis" />

                <h3 className="text-base font-semibold text-foreground/80 mb-3 mt-6">
                  What Gets Embedded
                </h3>
                <div className="space-y-2 mb-4">
                  {[
                    "Go runtime and compiled backend code",
                    "All Wails bindings",
                    "Complete React frontend (Vite build output)",
                    "All static assets (CSS, fonts, images)",
                    "SQLite driver (database file created at runtime)",
                  ].map((item) => (
                    <div
                      key={item}
                      className="flex items-start gap-2 text-sm text-muted-foreground/70"
                    >
                      <span className="text-primary/60 mt-0.5 shrink-0">
                        --
                      </span>
                      <span>{item}</span>
                    </div>
                  ))}
                </div>
              </div>

              {/* ════════════════════════════════════════════════════
                  11. Desktop vs Web Comparison
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={11} />
                  Desktop vs Web — Full Comparison
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Both modes share Grit conventions (GORM models, code
                  generation, Studio), but the runtime, communication model, and
                  distribution are fundamentally different.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/50 overflow-hidden mb-4">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Aspect
                        </th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Web (grit new)
                        </th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-primary/80 uppercase">
                          Desktop (grit new-desktop)
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/20 text-xs">
                      {[
                        [
                          "Backend",
                          "Gin HTTP server",
                          "Wails bindings (direct Go calls)",
                        ],
                        [
                          "Frontend",
                          "Next.js (App Router)",
                          "Vite + React + TanStack Router",
                        ],
                        [
                          "Database",
                          "PostgreSQL",
                          "SQLite (default, local)",
                        ],
                        [
                          "Communication",
                          "HTTP/REST + React Query + fetch",
                          "Direct Go calls + React Query + Wails bindings",
                        ],
                        [
                          "State management",
                          "React Query + fetch",
                          "React Query + Wails bindings",
                        ],
                        [
                          "Auth storage",
                          "JWT + HTTP cookies",
                          "JWT + local storage",
                        ],
                        [
                          "Project structure",
                          "Turborepo monorepo (apps/api, apps/web, apps/admin)",
                          "Single directory (root Go + frontend/)",
                        ],
                        [
                          "Distribution",
                          "Deploy to cloud / VPS",
                          "Distribute .exe / .app / binary",
                        ],
                        [
                          "File size",
                          "N/A (web-hosted)",
                          "~10-15MB single binary",
                        ],
                        [
                          "Offline support",
                          "Requires server",
                          "Fully offline — all data local",
                        ],
                        [
                          "Code generation markers",
                          "GRIT:MODELS, GRIT:ROUTES, etc.",
                          "grit:models, grit:methods, grit:nav, etc.",
                        ],
                        [
                          "Admin panel",
                          "Full admin panel (Next.js app)",
                          "None — uses in-app pages directly",
                        ],
                        [
                          "File storage",
                          "S3 / R2 / MinIO (presigned uploads)",
                          "Local filesystem (no S3)",
                        ],
                        [
                          "Background jobs",
                          "asynq + Redis",
                          "None — desktop apps use goroutines if needed",
                        ],
                      ].map(([aspect, web, desktop]) => (
                        <tr key={aspect}>
                          <td className="px-4 py-2 font-medium text-foreground/70">
                            {aspect}
                          </td>
                          <td className="px-4 py-2 text-muted-foreground/60">
                            {web}
                          </td>
                          <td className="px-4 py-2 text-primary/70">
                            {desktop}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* ════════════════════════════════════════════════════
                  12. Code Markers Reference
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={12} />
                  Code Markers — NEVER Delete
                </h2>
                <Warn>
                  These comments are injection points for the CLI. Removing them
                  permanently breaks{" "}
                  <code className="font-mono bg-red-500/10 px-1 rounded">
                    grit generate
                  </code>{" "}
                  and{" "}
                  <code className="font-mono bg-red-500/10 px-1 rounded">
                    grit remove
                  </code>
                  . Never remove, rename, or move them.
                </Warn>

                <p className="text-muted-foreground leading-relaxed mb-4">
                  Desktop projects use 13{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    grit:
                  </code>{" "}
                  markers across 6 files. The CLI injects code at these
                  locations during{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    grit generate resource
                  </code>{" "}
                  and removes it during{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    grit remove resource
                  </code>
                  .
                </p>

                <div className="rounded-xl border border-border/40 bg-card/50 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          File
                        </th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Marker
                        </th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Type
                        </th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">
                          Purpose
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/20 text-xs">
                      {[
                        [
                          "db.go",
                          "// grit:models",
                          "line-before",
                          "Add model to AutoMigrate list",
                        ],
                        [
                          "main.go",
                          "// grit:service-init",
                          "line-before",
                          "Initialize service with DB connection",
                        ],
                        [
                          "main.go",
                          "/* grit:app-args */",
                          "inline",
                          "Pass service to NewApp constructor",
                        ],
                        [
                          "app.go",
                          "// grit:imports",
                          "line-before",
                          "Import service package",
                        ],
                        [
                          "app.go",
                          "// grit:fields",
                          "line-before",
                          "Add service field to App struct",
                        ],
                        [
                          "app.go",
                          "/* grit:constructor-params */",
                          "inline",
                          "Add constructor parameter",
                        ],
                        [
                          "app.go",
                          "/* grit:constructor-assign */",
                          "inline",
                          "Assign service to App field",
                        ],
                        [
                          "app.go",
                          "// grit:methods",
                          "line-before",
                          "Add 7 bound methods (CRUD + export)",
                        ],
                        [
                          "types.go",
                          "// grit:input-types",
                          "line-before",
                          "Add Input struct for resource",
                        ],
                        [
                          "cmd/studio/main.go",
                          "// grit:studio-models",
                          "line-before",
                          "Register model in GORM Studio",
                        ],
                        [
                          "sidebar.tsx",
                          "// grit:nav-icons + // grit:nav",
                          "line-before",
                          "Add navigation icon import + sidebar link",
                        ],
                      ].map(([file, marker, type_, purpose]) => (
                        <tr key={marker}>
                          <td className="px-4 py-2 font-mono text-foreground/70">
                            {file}
                          </td>
                          <td className="px-4 py-2 font-mono text-primary/80">
                            {marker}
                          </td>
                          <td className="px-4 py-2 text-muted-foreground/60">
                            {type_}
                          </td>
                          <td className="px-4 py-2 text-muted-foreground/60 font-sans">
                            {purpose}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                <h3 className="text-base font-semibold text-foreground/80 mb-3">
                  What line-before vs inline Means
                </h3>
                <div className="space-y-3 mb-4">
                  <div className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                    <span className="text-[15px] font-semibold block mb-1">
                      line-before
                    </span>
                    <span className="text-sm text-muted-foreground/60">
                      New code is inserted on the line BEFORE the marker
                      comment. The marker stays on its own line. Used for models,
                      fields, methods, imports, routes.
                    </span>
                  </div>
                  <div className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                    <span className="text-[15px] font-semibold block mb-1">
                      inline
                    </span>
                    <span className="text-sm text-muted-foreground/60">
                      New code is inserted immediately before the{" "}
                      <code className="text-xs font-mono bg-accent/50 px-1 rounded">
                        /* grit:marker */
                      </code>{" "}
                      comment on the SAME line. Used for constructor parameters
                      and assignment expressions.
                    </span>
                  </div>
                </div>

                <h3 className="text-base font-semibold text-foreground/80 mb-3 mt-6">
                  Example: What Injection Looks Like
                </h3>
                <CodeBlock
                  language="go"
                  filename="app.go — before and after generating Product"
                  code={`// BEFORE generation:
type App struct {
    ctx         context.Context
    authService *service.AuthService
    // grit:fields
}

// AFTER generation:
type App struct {
    ctx            context.Context
    authService    *service.AuthService
    productService *service.ProductService
    // grit:fields
}`}
                />
              </div>

              {/* ════════════════════════════════════════════════════
                  13. 12 Injection Points — Visual
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={13} />
                  All 12 Injection Points — What Gets Injected
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  For a resource named{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    Product
                  </code>
                  , here is exactly what the CLI injects at each marker:
                </p>

                <div className="space-y-4 mb-4">
                  <CodeBlock
                    language="go"
                    filename="internal/db/db.go — // grit:models"
                    code={`&models.Product{},
// grit:models`}
                  />
                  <CodeBlock
                    language="go"
                    filename="main.go — // grit:service-init"
                    code={`productService := &service.ProductService{DB: database}
// grit:service-init`}
                  />
                  <CodeBlock
                    language="go"
                    filename="main.go — /* grit:app-args */"
                    code={`app := NewApp(authService, productService, /* grit:app-args */)`}
                  />
                  <CodeBlock
                    language="go"
                    filename="app.go — struct fields and constructor"
                    code={`type App struct {
    ctx            context.Context
    authService    *service.AuthService
    productService *service.ProductService
    // grit:fields
}

func NewApp(authService *service.AuthService, productService *service.ProductService, /* grit:constructor-params */) *App {
    return &App{authService: authService, productService: productService, /* grit:constructor-assign */}
}`}
                  />
                  <CodeBlock
                    language="tsx"
                    filename="frontend/src/routes/_layout/products.index.tsx — route file"
                    code={`import { createFileRoute } from "@tanstack/react-router";
// ... component imports

export const Route = createFileRoute("/_layout/products/")({
  component: ProductsListPage,
});

function ProductsListPage() {
  // ... list page with DataTable
}`}
                  />
                  <CodeBlock
                    language="tsx"
                    filename="frontend/src/components/sidebar.tsx — navigation"
                    code={`import { Package } from "lucide-react";
// grit:nav-icons

{ name: "Products", href: "/products", icon: Package },
// grit:nav`}
                  />
                </div>
              </div>

              {/* ════════════════════════════════════════════════════
                  14. LLM Examples
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={14} />
                  LLM Examples — Copy and Adapt
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-6">
                  Concrete examples an AI assistant can learn from. These cover
                  the most common workflows.
                </p>

                {/* Example 1 */}
                <div className="mb-8">
                  <h3 className="text-base font-semibold text-foreground/80 mb-2">
                    Example 1: Scaffold and Generate Multiple Resources
                  </h3>
                  <p className="text-sm text-muted-foreground/70 mb-3">
                    Create an inventory management desktop app with products and
                    suppliers.
                  </p>
                  <CodeBlock
                    language="bash"
                    code={`# Create the project
grit new-desktop inventory

# Navigate into project
cd inventory

# Generate resources
grit generate resource Product --fields "name:string,sku:string,price:float,stock:int,reorder_level:int,active:bool"
grit generate resource Supplier --fields "name:string,email:string,phone:string,address:text,notes:text"

# Start development
wails dev`}
                  />
                </div>

                {/* Example 2 */}
                <div className="mb-8">
                  <h3 className="text-base font-semibold text-foreground/80 mb-2">
                    Example 2: Add a Custom Service Method
                  </h3>
                  <p className="text-sm text-muted-foreground/70 mb-3">
                    Add a{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      GetBySKU
                    </code>{" "}
                    method to the product service for barcode-based lookup.
                  </p>
                  <CodeBlock
                    language="go"
                    filename="internal/service/product.go — add this method"
                    code={`func (s *ProductService) GetBySKU(sku string) (*models.Product, error) {
    var product models.Product
    if err := s.DB.Where("sku = ?", sku).First(&product).Error; err != nil {
        return nil, fmt.Errorf("product not found with SKU %s: %w", sku, err)
    }
    return &product, nil
}`}
                  />
                </div>

                {/* Example 3 */}
                <div className="mb-8">
                  <h3 className="text-base font-semibold text-foreground/80 mb-2">
                    Example 3: Expose a Custom Method via Wails
                  </h3>
                  <p className="text-sm text-muted-foreground/70 mb-3">
                    Make the custom method available in React by adding it to
                    the App struct. Then call from React.
                  </p>
                  <CodeBlock
                    language="go"
                    filename="app.go — add method to App struct (below grit:methods marker)"
                    code={`// Custom method — add ABOVE the // grit:methods marker (not between generated code)
func (a *App) GetProductBySKU(sku string) (*models.Product, error) {
    return a.productService.GetBySKU(sku)
}`}
                  />
                  <CodeBlock
                    language="tsx"
                    filename="frontend/src/routes/_layout/products.lookup.tsx — call from React"
                    code={`import { GetProductBySKU } from "../../wailsjs/go/main/App";

export default function ProductLookup() {
  const [sku, setSku] = useState("");
  const [product, setProduct] = useState(null);

  const handleLookup = async () => {
    try {
      const result = await GetProductBySKU(sku);
      setProduct(result);
    } catch (err) {
      toast.error("Product not found");
    }
  };

  return (
    <div>
      <input value={sku} onChange={(e) => setSku(e.target.value)} placeholder="Scan barcode..." />
      <button onClick={handleLookup}>Lookup</button>
      {product && <div>{product.name} — ${"{product.price}"}</div>}
    </div>
  );
}`}
                  />
                  <Tip>
                    After adding a new Go method, restart{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      wails dev
                    </code>{" "}
                    so Wails regenerates the TypeScript bindings. Then import
                    the new function from{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      wailsjs/go/main/App
                    </code>
                    .
                  </Tip>
                </div>

                {/* Example 4 */}
                <div className="mb-8">
                  <h3 className="text-base font-semibold text-foreground/80 mb-2">
                    Example 4: Generate with Slug Field
                  </h3>
                  <p className="text-sm text-muted-foreground/70 mb-3">
                    Generate an Article resource with an auto-generated slug
                    from the title field.
                  </p>
                  <CodeBlock
                    terminal
                    code={`grit generate resource Article --fields "title:string,slug:slug:source=title,content:richtext,published:bool"`}
                  />
                  <p className="text-sm text-muted-foreground/60 mt-2">
                    The slug field is excluded from the form — it is
                    auto-generated in the Go model&apos;s{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      BeforeCreate
                    </code>{" "}
                    hook. The model gets a{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      uniqueIndex
                    </code>{" "}
                    on the slug column.
                  </p>
                </div>

                {/* Example 5 */}
                <div className="mb-8">
                  <h3 className="text-base font-semibold text-foreground/80 mb-2">
                    Example 5: Build for Distribution
                  </h3>
                  <CodeBlock
                    language="bash"
                    code={`# Build for current platform
grit compile

# Output:
#   Windows: build/bin/inventory.exe
#   macOS:   build/bin/inventory.app
#   Linux:   build/bin/inventory

# Cross-compile for a specific platform:
wails build -platform windows/amd64
wails build -platform darwin/arm64

# Create a Windows installer:
wails build -nsis`}
                  />
                </div>

                {/* Example 6 */}
                <div className="mb-4">
                  <h3 className="text-base font-semibold text-foreground/80 mb-2">
                    Example 6: Remove a Resource
                  </h3>
                  <p className="text-sm text-muted-foreground/70 mb-3">
                    Cleanly remove a resource — deletes all generated files
                    and reverses all 10 injections. Markers stay intact for
                    future generation.
                  </p>
                  <CodeBlock
                    terminal
                    code={`grit remove resource Supplier

# Deletes:
#   internal/models/supplier.go
#   internal/service/supplier.go
#   frontend/src/routes/_layout/suppliers.index.tsx
#   frontend/src/routes/_layout/suppliers.new.tsx
#   frontend/src/routes/_layout/suppliers.$id.edit.tsx
#
# Removes injected code from:
#   db.go, main.go, app.go, types.go, studio/main.go, sidebar.tsx`}
                  />
                </div>
              </div>

              {/* ════════════════════════════════════════════════════
                  15. Golden Rules
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={15} />
                  Golden Rules for AI Assistants — Never Break These
                </h2>
                <Note>
                  Non-negotiable. Violating them causes silent failures, broken
                  code generation, or corrupted project state.
                </Note>

                <div className="space-y-3 mb-4">
                  {[
                    {
                      rule: "Always use grit generate resource for new resources",
                      detail:
                        "Never create model, service, or page files manually. The CLI handles 5 file creations + 10 injections in one atomic operation. Manual creation will miss injection points and break future generation/removal.",
                    },
                    {
                      rule: "Never modify grit: markers",
                      detail:
                        "The comments // grit:models, // grit:fields, // grit:methods, // grit:nav, etc. are permanent injection points. Do not remove, rename, reformat, or move them. They must stay exactly as scaffolded.",
                    },
                    {
                      rule: "Use grit remove resource to undo — never delete files manually",
                      detail:
                        "Manual deletion leaves injected code orphaned in app.go, main.go, db.go, and sidebar.tsx. Always use grit remove resource <Name> which cleans up all 10 injection points.",
                    },
                    {
                      rule: "Desktop uses SQLite, not PostgreSQL",
                      detail:
                        "Never suggest PostgreSQL-specific features (arrays, JSONB, full-text search) for desktop projects. SQLite is the default and only supported database. No Docker needed.",
                    },
                    {
                      rule: "Methods must be on the App struct to be exposed to React",
                      detail:
                        "Only exported methods on the App struct (in app.go) are available as Wails bindings. Methods on other structs or in other files are not accessible from the frontend.",
                    },
                    {
                      rule: "The frontend calls Go functions directly — there is no REST API",
                      detail:
                        "Never suggest fetch(), Axios, HTTP endpoints, or REST patterns for desktop projects. React calls Go via GetProducts(), CreateProduct(), etc. — generated Wails TypeScript bindings.",
                    },
                    {
                      rule: "Restart wails dev after Go file changes",
                      detail:
                        "Wails auto-detects Go changes and rebuilds, but new methods require a full restart so TypeScript bindings are regenerated. If a new method is not available in React, restart wails dev.",
                    },
                    {
                      rule: "Add custom methods ABOVE the grit:methods marker, not below",
                      detail:
                        "Code injected by grit generate goes directly above the // grit:methods marker. Place custom methods above the generated block to avoid them being accidentally removed by grit remove.",
                    },
                    {
                      rule: "Do not create handlers or routes files for desktop",
                      detail:
                        "Desktop projects do not have handlers/, routes/, or middleware/ directories. All request handling is done through methods on the App struct. If a user asks for an API endpoint, explain that desktop uses Wails bindings instead.",
                    },
                    {
                      rule: "Generate parent models before child models (belongs_to)",
                      detail:
                        "When using belongs_to relationships, generate the parent resource first so the referenced model exists. For example, generate Category before generating Product with a category_id:belongs_to field.",
                    },
                  ].map((item, i) => (
                    <div
                      key={item.rule}
                      className="flex gap-3 p-4 rounded-lg border border-border/30 bg-card/30"
                    >
                      <div className="flex h-6 w-6 shrink-0 items-center justify-center rounded-md bg-amber-500/10 border border-amber-500/20 text-xs font-mono font-semibold text-amber-500/80 mt-0.5">
                        {i + 1}
                      </div>
                      <div>
                        <p className="text-sm font-semibold text-foreground/80 mb-1">
                          {item.rule}
                        </p>
                        <p className="text-xs text-muted-foreground/60 leading-relaxed">
                          {item.detail}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* ════════════════════════════════════════════════════
                  16. Quick Reference
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={16} />
                  Quick Build Reference
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Copy-paste recipes for the most common desktop workflows.
                </p>

                <div className="space-y-4">
                  <CodeBlock
                    language="bash"
                    filename="Scaffold a new desktop project"
                    code={`grit new-desktop myapp
cd myapp
wails dev`}
                  />

                  <CodeBlock
                    language="bash"
                    filename="Generate resources"
                    code={`grit generate resource Product --fields "name:string,price:float,stock:int,active:bool"
grit generate resource Category --fields "name:string,description:text"
grit generate resource Order --fields "customer_name:string,total:float,status:string,notes:text,completed:bool"`}
                  />

                  <CodeBlock
                    language="bash"
                    filename="Remove a resource"
                    code="grit remove resource Order"
                  />

                  <CodeBlock
                    language="bash"
                    filename="Open GORM Studio"
                    code={`grit studio
# Opens browser at http://localhost:8080/studio`}
                  />

                  <CodeBlock
                    language="bash"
                    filename="Build for production"
                    code={`grit compile
# Binary at: build/bin/myapp.exe (Windows) or build/bin/myapp (macOS/Linux)`}
                  />

                  <CodeBlock
                    language="bash"
                    filename="Cross-compile + installer"
                    code={`wails build -platform windows/amd64
wails build -platform darwin/arm64
wails build -nsis  # Windows installer`}
                  />
                </div>
              </div>

              {/* ════════════════════════════════════════════════════
                  17. Common Mistakes LLMs Make
              ════════════════════════════════════════════════════ */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={17} />
                  Common Mistakes LLMs Make with Desktop Projects
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  These are the most frequent errors AI assistants make when
                  helping with Grit Desktop. Avoid all of them.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/50 overflow-hidden mb-4">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-red-400/80 uppercase">
                          Wrong
                        </th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-emerald-400/80 uppercase">
                          Correct
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/20 text-xs">
                      {[
                        [
                          'Suggesting fetch("/api/products") in React',
                          "Use GetProducts() from Wails bindings",
                        ],
                        [
                          "Creating a handlers/ or routes/ directory",
                          "Add methods to the App struct in app.go",
                        ],
                        [
                          "Suggesting Docker or PostgreSQL setup",
                          "SQLite works out of the box — no Docker needed",
                        ],
                        [
                          "Running grit start server + grit start client",
                          "Use grit start or wails dev (single command)",
                        ],
                        [
                          "Manually creating model/service/page files",
                          "Use grit generate resource for everything",
                        ],
                        [
                          "Deleting generated files to remove a resource",
                          "Use grit remove resource <Name>",
                        ],
                        [
                          "Suggesting middleware for authentication",
                          "Auth is handled in Go methods — check JWT in service layer",
                        ],
                        [
                          "Adding icon + nav item to sidebar.tsx",
                          "Add imports ABOVE the marker — inject goes before",
                        ],
                        [
                          "Using pnpm dev or turbo dev for desktop",
                          "Use wails dev — no Turborepo in desktop projects",
                        ],
                        [
                          "Creating apps/ or packages/ directories",
                          "Desktop is a single directory — no monorepo structure",
                        ],
                      ].map(([wrong, correct]) => (
                        <tr key={wrong}>
                          <td className="px-4 py-2 text-red-400/70">
                            {wrong}
                          </td>
                          <td className="px-4 py-2 text-emerald-400/70">
                            {correct}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

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

              {/* Nav */}
              <div className="flex items-center justify-between pt-6 border-t border-border/30">
                <Button
                  variant="ghost"
                  size="sm"
                  asChild
                  className="text-muted-foreground/60 hover:text-foreground"
                >
                  <Link
                    href="/docs/desktop/pos-app"
                    className="gap-1.5"
                  >
                    <ArrowLeft className="h-3.5 w-3.5" />
                    Build a POS App
                  </Link>
                </Button>
                <div />
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
