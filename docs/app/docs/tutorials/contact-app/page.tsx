import Link from "next/link";
import { ArrowRight, ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { CodeBlock } from "@/components/code-block";
import { getDocMetadata } from "@/config/docs-metadata";

export const metadata = getDocMetadata("/docs/tutorials/contact-app");

export default function TutorialContactAppPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">
                Tutorial
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Your First App
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Build a contact manager from scratch using Grit. You will create
                two resources &mdash; <strong>Group</strong> and{" "}
                <strong>Contact</strong> &mdash; generate them with a single CLI
                command each, and explore the admin panel, GORM Studio, and API
                docs that Grit gives you for free.
              </p>
            </div>

            {/* Prerequisites */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Prerequisites
              </h2>
              <ul className="space-y-2.5">
                {[
                  "Go 1.21+ installed",
                  "Node.js 18+ and pnpm installed",
                  "Docker and Docker Compose installed",
                  "Grit CLI installed globally (go install github.com/MUKE-coder/grit/cmd/grit@latest)",
                ].map((item) => (
                  <li
                    key={item}
                    className="flex items-start gap-2.5 text-[14px] text-muted-foreground"
                  >
                    <span className="text-primary mt-1">&#10003;</span>
                    {item}
                  </li>
                ))}
              </ul>
            </div>

            {/* ============================================================ */}
            {/* STEP 1 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                1
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Create the project
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Scaffold a new Grit monorepo called <code>contact-app</code>.
                  This creates the Go API, Next.js web app, admin panel, shared
                  package, and Docker configuration in one shot.
                </p>

                <CodeBlock terminal code={`grit new contact-app\ncd contact-app`} className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  Grit prints an ASCII art logo, creates the folder structure,
                  initializes
                  <code> go.mod</code>, and prints the next steps. Your project
                  is ready.
                </p>
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 2 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                2
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Install dependencies
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Install all Node.js packages for the web app, admin panel, and
                  shared package. Grit uses pnpm workspaces so everything
                  installs from the project root.
                </p>

                <CodeBlock terminal code="pnpm install" className="glow-purple-sm" />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 3 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                3
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Start Docker services
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Spin up PostgreSQL, Redis, MinIO (local S3), and Mailhog.
                  These run in the background and persist data across restarts.
                </p>

                <CodeBlock terminal code="docker compose up -d" className="glow-purple-sm" />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 4 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                4
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Check for port conflicts
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Docker needs port <code>5432</code> (PostgreSQL),{" "}
                  <code>6379</code> (Redis), <code>9000/9001</code> (MinIO), and{" "}
                  <code>1025/8025</code> (Mailhog). If any of these are already
                  in use, you will see an error.
                </p>

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  The most common conflict is port <code>5432</code> &mdash; if
                  you have PostgreSQL installed locally it&apos;s already
                  listening there. Check what&apos;s using the port:
                </p>

                <CodeBlock terminal filename="windows" code={`# Find the process using port 5432\nnetstat -ano | findstr :5432\n\n# Stop the local PostgreSQL service (Run as administrator)\nnet stop postgresql-x64-16`} className="glow-purple-sm mb-4" />

                <CodeBlock terminal filename="macOS / Linux" code={`# Find the process using port 5432\nlsof -i :5432\n\n# Stop local PostgreSQL (macOS with Homebrew)\nbrew services stop postgresql@16`} className="glow-purple-sm mb-4" />

                <div className="rounded-xl border border-yellow-500/20 bg-yellow-500/5 p-4 mb-4">
                  <p className="text-[13px] text-yellow-200/80 leading-relaxed mb-3">
                    <strong className="text-yellow-300">
                      You have three options:
                    </strong>
                  </p>
                  <ol className="space-y-2 text-[13px] text-yellow-200/80 leading-relaxed list-decimal list-inside">
                    <li>
                      <strong className="text-yellow-200/90">
                        Stop the local service
                      </strong>{" "}
                      &mdash; run <code>net stop postgresql-x64-16</code>{" "}
                      (Windows) or <code>brew services stop postgresql@16</code>{" "}
                      (macOS). This frees the port so Docker can use it.
                    </li>
                    <li>
                      <strong className="text-yellow-200/90">
                        Uninstall local PostgreSQL
                      </strong>{" "}
                      &mdash; if you only use Docker for databases, remove the
                      local installation entirely to avoid future conflicts.
                    </li>
                    <li>
                      <strong className="text-yellow-200/90">
                        Change the port
                      </strong>{" "}
                      &mdash; update <code>docker-compose.yml</code> to map a
                      different host port:
                      <div className="mt-1.5 font-mono text-xs text-yellow-200/60 bg-yellow-500/5 rounded px-2 py-1 inline-block">
                        ports: &quot;5433:5432&quot;
                      </div>
                      <br />
                      Then update <code>DATABASE_URL</code> in <code>.env</code>{" "}
                      to use port <code>5433</code>.
                    </li>
                  </ol>
                </div>

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  Run <code>docker compose ps</code> to verify all four services
                  are running and healthy.
                </p>
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 5 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                5
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">Run migrations</h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Create the database tables. Grit uses GORM AutoMigrate, which
                  reads your Go models and creates/updates the corresponding
                  PostgreSQL tables automatically.
                </p>

                <CodeBlock terminal code="grit migrate" className="glow-purple-sm" />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 6 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                6
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Seed the database
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Populate the database with an admin user and sample blog data.
                  The default admin credentials are{" "}
                  <code>admin@example.com</code> / <code>password</code>.
                </p>

                <CodeBlock terminal code="grit seed" className="glow-purple-sm" />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 7 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                7
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Start the API server
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Open a terminal and start the Go API. This compiles and runs
                  the server on <code>http://localhost:8080</code>.
                </p>

                <CodeBlock terminal filename="terminal 1" code="grit start server" className="glow-purple-sm" />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 8 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                8
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Start the frontend
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Open a <strong>second terminal</strong> and start both Next.js
                  apps (web + admin) using Turborepo. The web app runs on{" "}
                  <code>localhost:3000</code> and the admin panel on{" "}
                  <code>localhost:3001</code>.
                </p>

                <CodeBlock terminal filename="terminal 2" code="grit start client" className="glow-purple-sm" />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 9 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                9
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Create a user and log in
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Open <code>http://localhost:3000</code> in your browser. You
                  will see the login page. Click <strong>Register</strong> to
                  create a new account, then log in.
                </p>

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  For the admin panel at <code>http://localhost:3001</code>, use
                  the seeded admin credentials:
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-4">
                  <div className="p-4 font-mono text-sm space-y-1">
                    <div>
                      <span className="text-muted-foreground/50">Email: </span>
                      <span className="text-foreground/80">
                        admin@example.com
                      </span>
                    </div>
                    <div>
                      <span className="text-muted-foreground/50">
                        Password:{" "}
                      </span>
                      <span className="text-foreground/80">password</span>
                    </div>
                  </div>
                </div>

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  After logging in to the admin panel, you will see the
                  dashboard with stats cards, a chart, and the Users resource in
                  the sidebar. The blog example is also available &mdash; click{" "}
                  <strong>Posts</strong> in the sidebar to see the seeded blog
                  data.
                </p>
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 10 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                10
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Explore the built-in tools
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Grit ships with powerful development tools out of the box.
                  Open these URLs while your server is running:
                </p>

                <ul className="space-y-2">
                  {[
                    {
                      url: "http://localhost:8080/studio",
                      desc: "GORM Studio \u2014 visual database browser to inspect tables, run queries, and see relationships",
                    },
                    {
                      url: "http://localhost:8080/docs",
                      desc: "API Documentation \u2014 auto-generated interactive docs for every endpoint",
                    },
                    {
                      url: "http://localhost:8080/sentinel",
                      desc: "Sentinel \u2014 security dashboard with WAF, rate limiting, and threat monitoring",
                    },
                    {
                      url: "http://localhost:8025",
                      desc: "Mailhog \u2014 catches all emails sent during development",
                    },
                    {
                      url: "http://localhost:9001",
                      desc: "MinIO Console \u2014 browse uploaded files (login: minioadmin / minioadmin)",
                    },
                  ].map((item) => (
                    <li
                      key={item.url}
                      className="flex items-start gap-2.5 text-[13px] text-muted-foreground/70"
                    >
                      <code className="text-primary/70 text-xs whitespace-nowrap mt-0.5">
                        {item.url}
                      </code>
                      <span>&mdash; {item.desc}</span>
                    </li>
                  ))}
                </ul>
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 11 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                11
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Generate the Group resource
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  A contact belongs to a group (e.g. Friends, Family, Work).
                  Generate the Group resource first since Contact will reference
                  it. The <code>slug:slug:name</code> field automatically
                  generates a URL-friendly slug from the group name.
                </p>

                <CodeBlock terminal code={`grit generate resource Group --fields "name:string,slug:slug:name"`} className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  The generator creates these files:
                </p>

                <CodeBlock
                  language="bash"
                  filename="generated files"
                  code={`apps/api/internal/models/group.go        # GORM model
apps/api/internal/handlers/group.go      # CRUD handler
apps/api/internal/services/group.go      # Business logic
packages/shared/schemas/group.ts         # Zod validation
packages/shared/types/group.ts           # TypeScript types
apps/web/hooks/use-groups.ts             # React Query hooks (web)
apps/admin/hooks/use-groups.ts           # React Query hooks (admin)
apps/admin/app/resources/groups/page.tsx # Admin page
apps/admin/resources/groups.ts           # Resource definition`}
                />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Here is the generated Go model:
                </p>

                <CodeBlock
                  filename="apps/api/internal/models/group.go"
                  code={`package models

import (
    "time"
    "gorm.io/gorm"
)

type Group struct {
    ID        uint           \`gorm:"primarykey" json:"id"\`
    Name      string         \`gorm:"size:255;not null" json:"name" binding:"required"\`
    Slug      string         \`gorm:"size:255;uniqueIndex;not null" json:"slug"\`
    CreatedAt time.Time      \`json:"created_at"\`
    UpdatedAt time.Time      \`json:"updated_at"\`
    DeletedAt gorm.DeletedAt \`gorm:"index" json:"deleted_at,omitempty"\`
}`}
                />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 12 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                12
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Generate the Contact resource
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Now generate the Contact resource with all the fields: name,
                  email, phone, country, image (as a string array for multiple
                  photos), and a <code>belongs_to</code> relationship to Group.
                </p>

                <CodeBlock terminal code={`grit generate resource Contact --fields "name:string,email:string,phone:string,country:string,image:string_array,group:belongs_to:Group"`} className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  The generated Contact model includes a <code>GroupID</code>{" "}
                  foreign key and a <code>Group</code> relation. The{" "}
                  <code>Image</code> field is a JSON array that stores multiple
                  image URLs:
                </p>

                <CodeBlock
                  filename="apps/api/internal/models/contact.go"
                  code={`package models

import (
    "time"
    "gorm.io/datatypes"
    "gorm.io/gorm"
)

type Contact struct {
    ID        uint                        \`gorm:"primarykey" json:"id"\`
    Name      string                      \`gorm:"size:255;not null" json:"name" binding:"required"\`
    Email     string                      \`gorm:"size:255;not null" json:"email" binding:"required"\`
    Phone     string                      \`gorm:"size:255;not null" json:"phone" binding:"required"\`
    Country   string                      \`gorm:"size:255;not null" json:"country" binding:"required"\`
    Image     datatypes.JSONSlice[string] \`gorm:"type:json" json:"image"\`
    GroupID   uint                        \`gorm:"index;not null" json:"group_id" binding:"required"\`
    Group     Group                       \`gorm:"foreignKey:GroupID" json:"group,omitempty"\`
    CreatedAt time.Time                   \`json:"created_at"\`
    UpdatedAt time.Time                   \`json:"updated_at"\`
    DeletedAt gorm.DeletedAt              \`gorm:"index" json:"deleted_at,omitempty"\`
}`}
                />

                <div className="rounded-xl border border-blue-500/20 bg-blue-500/5 p-4 mt-4">
                  <p className="text-[13px] text-blue-200/80 leading-relaxed">
                    <strong className="text-blue-300">Note:</strong> The{" "}
                    <code>string_array</code> field type uses{" "}
                    <code>datatypes.JSONSlice[string]</code> from{" "}
                    <code>gorm.io/datatypes</code>. The generator automatically
                    runs <code>go get</code> to install this dependency for you.
                  </p>
                </div>
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 13 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                13
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">Rebuild the API</h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  The generator injected new code into your Go API &mdash;
                  models, handlers, routes, and GORM Studio registration. Build
                  the API to make sure everything compiles:
                </p>

                <CodeBlock terminal code="cd apps/api && go build ./..." className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  If the build succeeds with no output, everything is wired
                  correctly. If you see errors, check that both{" "}
                  <code>group.go</code> and <code>contact.go</code> models exist
                  in <code>apps/api/internal/models/</code>.
                </p>
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 14 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                14
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Restart and check the admin panel
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Stop the API server (Ctrl+C in terminal 1) and restart it.
                  GORM will automatically create the <code>groups</code> and{" "}
                  <code>contacts</code> tables in PostgreSQL.
                </p>

                <CodeBlock terminal filename="terminal 1" code="grit start server" className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  Open the admin panel at <code>http://localhost:3001</code>.
                  You will see <strong>Groups</strong> and{" "}
                  <strong>Contacts</strong> in the sidebar alongside Users and
                  Posts. The admin panel reads from the resource registry at
                  runtime &mdash; no manual sidebar configuration needed.
                </p>
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 15 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                15
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Create groups and contacts
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  In the admin panel, go to <strong>Groups</strong> and create a
                  few groups:
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-4">
                  <div className="p-4 text-sm space-y-2">
                    {[
                      { name: "Friends", slug: "friends" },
                      { name: "Family", slug: "family" },
                      { name: "Work", slug: "work" },
                    ].map((group) => (
                      <div key={group.slug} className="flex items-center gap-3">
                        <span className="text-primary/60 text-xs font-mono w-16">
                          {group.slug}
                        </span>
                        <span className="text-foreground/80">{group.name}</span>
                      </div>
                    ))}
                  </div>
                </div>

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Then go to <strong>Contacts</strong> and create some contacts.
                  When you click <strong>Create</strong>, the form will show a
                  dropdown to select the group &mdash; this is the{" "}
                  <code>belongs_to</code> relationship in action. Fill in the
                  name, email, phone, country, and optionally upload images.
                </p>

                <div className="rounded-xl border border-blue-500/20 bg-blue-500/5 p-4">
                  <p className="text-[13px] text-blue-200/80 leading-relaxed">
                    <strong className="text-blue-300">Tip:</strong> Open GORM
                    Studio at <code>http://localhost:8080/studio</code> to see
                    your groups and contacts in the database. You can browse the
                    tables, see the relationships, and even run raw SQL queries.
                  </p>
                </div>
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 16 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                16
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">Explore the API</h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Every resource you generate gets a full REST API. Try these
                  endpoints in your browser or with a tool like Postman:
                </p>

                <ul className="space-y-2">
                  {[
                    {
                      method: "GET",
                      url: "/api/groups",
                      desc: "List all groups with pagination",
                    },
                    {
                      method: "GET",
                      url: "/api/groups/:id",
                      desc: "Get a single group",
                    },
                    {
                      method: "POST",
                      url: "/api/groups",
                      desc: "Create a new group",
                    },
                    {
                      method: "GET",
                      url: "/api/contacts",
                      desc: "List all contacts with pagination",
                    },
                    {
                      method: "GET",
                      url: "/api/contacts/:id",
                      desc: "Get a single contact with its group",
                    },
                    {
                      method: "POST",
                      url: "/api/contacts",
                      desc: "Create a new contact",
                    },
                  ].map((item) => (
                    <li
                      key={`${item.method}-${item.url}`}
                      className="flex items-start gap-2.5 text-[13px] text-muted-foreground/70"
                    >
                      <code className="text-xs font-mono mt-0.5">
                        <span className="text-green-400/70">{item.method}</span>{" "}
                        <span className="text-primary/70">{item.url}</span>
                      </code>
                      <span>&mdash; {item.desc}</span>
                    </li>
                  ))}
                </ul>

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mt-4">
                  All endpoints support <code>?page=1&amp;page_size=20</code>{" "}
                  for pagination, <code>?sort=name&amp;order=asc</code> for
                  sorting, and <code>?search=john</code> for full-text search.
                  Check the interactive API docs at{" "}
                  <code>http://localhost:8080/docs</code> for the complete
                  reference.
                </p>
              </div>
            </div>

            {/* ============================================================ */}
            {/* SUMMARY */}
            {/* ============================================================ */}
            <div className="mb-12 rounded-xl border border-border/40 bg-card/30 p-6">
              <h2 className="text-xl font-semibold tracking-tight mb-4">
                What you&apos;ve built
              </h2>
              <ul className="space-y-2.5">
                {[
                  "A full-stack contact manager with Go API and Next.js frontend",
                  "Group and Contact resources generated with a single CLI command each",
                  "A BelongsTo relationship between Contacts and Groups",
                  "An admin panel with sortable tables, forms, and group selection",
                  "A JSON array field (image) for storing multiple images per contact",
                  "REST API with pagination, sorting, and search out of the box",
                  "GORM Studio for visual database browsing",
                  "Auto-generated API documentation",
                  "Docker-based PostgreSQL, Redis, MinIO, and Mailhog running locally",
                ].map((item) => (
                  <li
                    key={item}
                    className="flex items-start gap-2.5 text-[14px] text-muted-foreground"
                  >
                    <span className="text-primary mt-0.5">&#10003;</span>
                    {item}
                  </li>
                ))}
              </ul>
            </div>

            <div className="prose-grit mb-8">
              <h2>Next steps</h2>
              <p>
                Now that you have a working contact manager, here are some ideas
                to extend it:
              </p>
              <ul>
                <li>
                  <strong>Add tags</strong> &mdash; generate a <code>Tag</code>{" "}
                  resource and create a many-to-many relationship so contacts
                  can have multiple tags like &quot;VIP&quot; or
                  &quot;Newsletter&quot;.
                </li>
                <li>
                  <strong>Role-based access</strong> &mdash; use{" "}
                  <code>grit add role MANAGER</code> and the{" "}
                  <code>--roles</code> flag on generate to restrict who can
                  manage contacts.
                </li>
                <li>
                  <strong>Export contacts</strong> &mdash; add a custom Go
                  handler that exports contacts as a CSV file.
                </li>
                <li>
                  <strong>Send emails</strong> &mdash; use the built-in Resend
                  email service to send a welcome email when a new contact is
                  created.
                </li>
                <li>
                  <strong>Build the web frontend</strong> &mdash; create a
                  public contacts page in the web app using the generated React
                  Query hooks.
                </li>
              </ul>
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
                  href="/docs/getting-started/quick-start"
                  className="gap-1.5"
                >
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Quick Start
                </Link>
              </Button>
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="text-muted-foreground/60 hover:text-foreground"
              >
                <Link href="/docs/tutorials/blog" className="gap-1.5">
                  Build a Blog
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
