import Link from "next/link";
import { ArrowRight, ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/tutorials/learn')

export default function TutorialLearnPage() {
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
                Learn Grit Step by Step
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                A progressive, hands-on tutorial that takes you from zero to a
                full-stack task manager. You&apos;ll create 4 resources, add
                relationships, customize Go handlers, and build a working
                Next.js frontend &mdash; learning Go and React patterns along
                the way. Estimated time: ~1 hour.
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

            {/* Curriculum */}
            <div className="mb-12 rounded-xl border border-border/40 bg-card/30 p-6">
              <h2 className="text-xl font-semibold tracking-tight mb-4">
                Curriculum
              </h2>
              <div className="space-y-4">
                {[
                  {
                    milestone: "Milestone 1",
                    title: "Your First Resource",
                    desc: "Project setup, code generation, Go structs, admin CRUD",
                    models: "Task",
                  },
                  {
                    milestone: "Milestone 2",
                    title: "Adding Categories",
                    desc: "belongs_to relationships, foreign keys, eager loading",
                    models: "Task + Category",
                  },
                  {
                    milestone: "Milestone 3",
                    title: "Tags & Comments",
                    desc: "many_to_many, multi-select, complex relationships",
                    models: "+ Tag + Comment",
                  },
                  {
                    milestone: "Milestone 4",
                    title: "The Web App",
                    desc: "React data fetching, creating records, building pages",
                    models: "All 4 models",
                  },
                ].map((m) => (
                  <div key={m.milestone} className="flex items-start gap-4">
                    <span className="shrink-0 rounded-md bg-primary/10 border border-primary/15 px-2 py-1 text-xs font-mono font-semibold text-primary">
                      {m.milestone}
                    </span>
                    <div>
                      <p className="text-sm font-medium text-foreground/90">
                        {m.title}{" "}
                        <span className="text-xs text-muted-foreground/50 font-normal">
                          ({m.models})
                        </span>
                      </p>
                      <p className="text-xs text-muted-foreground/60">
                        {m.desc}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* ============================================================ */}
            {/* MILESTONE 1 */}
            {/* ============================================================ */}
            <div className="mb-8">
              <h2 className="text-2xl font-semibold tracking-tight mb-2">
                Milestone 1: Your First Resource
              </h2>
              <p className="text-muted-foreground text-sm">
                One model, one command, full CRUD &mdash; let&apos;s get
                familiar with how Grit works.
              </p>
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
                  Create a New Project
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Scaffold a new Grit monorepo called{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    task-manager
                  </code>
                  . This creates a monorepo with a Go API, two Next.js apps (web
                  and admin), a shared package, and Docker config. Here&apos;s
                  the command:
                </p>

                <CodeBlock terminal code="grit new task-manager" className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  This creates a monorepo with a Go API, two Next.js apps (web
                  and admin), a shared package, and Docker config. Here&apos;s
                  the structure:
                </p>

                <CodeBlock filename="project structure" code={`task-manager/
├── apps/
│   ├── api/          # Go backend (Gin + GORM)
│   │   ├── cmd/server/main.go
│   │   └── internal/ # models, handlers, services, middleware
│   ├── web/          # Next.js public frontend
│   └── admin/        # Next.js admin panel
├── packages/
│   └── shared/       # Zod schemas, TypeScript types
├── docker-compose.yml
└── pnpm-workspace.yaml`} />
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
                  Start the Infrastructure
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Spin up PostgreSQL, Redis, MinIO (S3-compatible file storage),
                  and Mailhog. These run in the background and persist data
                  across restarts.
                </p>

                <CodeBlock terminal code={`cd task-manager\ndocker compose up -d`} className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Now install dependencies and start the servers:
                </p>

                <CodeBlock terminal code={`pnpm install\ncd apps/api && go run cmd/server/main.go &\ncd apps/admin && pnpm dev &`} className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  Your API is running at{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    http://localhost:8080
                  </code>{" "}
                  and the admin panel at{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    http://localhost:3001
                  </code>
                  . Let&apos;s generate our first resource.
                </p>
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
                  Generate Your First Resource
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Use the code generator to create a full-stack Task resource
                  with a single command. This generates the Go model, handler,
                  service, Zod schema, TypeScript types, React Query hooks, and
                  admin page &mdash; all wired together.
                </p>

                <CodeBlock terminal code={`grit generate resource Task --fields "title:string,description:text,status:string,priority:int,due_date:date,completed:bool"`} className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  One command just generated 7 files across the entire stack:
                </p>

                <CodeBlock language="bash" filename="generated files" code={`✓ apps/api/internal/models/task.go        # Go model
✓ apps/api/internal/services/task.go      # Service layer
✓ apps/api/internal/handlers/task.go      # API endpoints
✓ packages/shared/schemas/task.ts         # Zod validation
✓ packages/shared/types/task.ts           # TypeScript types
✓ apps/admin/resources/tasks.ts           # Admin resource definition
✓ apps/admin/app/.../tasks/page.tsx       # Admin page`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Here is the generated Go model:
                </p>

                <CodeBlock filename="apps/api/internal/models/task.go" code={`package models

import (
    "time"
    "gorm.io/gorm"
)

type Task struct {
    ID          uint           \`gorm:"primarykey" json:"id"\`
    Title       string         \`gorm:"size:255" json:"title" binding:"required"\`
    Description string         \`gorm:"type:text" json:"description"\`
    Status      string         \`gorm:"size:255" json:"status" binding:"required"\`
    Priority    int            \`json:"priority"\`
    DueDate     *time.Time     \`gorm:"type:date" json:"due_date"\`
    Completed   bool           \`json:"completed"\`
    CreatedAt   time.Time      \`json:"created_at"\`
    UpdatedAt   time.Time      \`json:"updated_at"\`
    DeletedAt   gorm.DeletedAt \`gorm:"index" json:"-"\`
}`} />

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-4 mb-4">
                  <p className="text-xs font-semibold text-primary/80 uppercase tracking-wider mb-2">
                    Go Basics
                  </p>
                  <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                    In Go, a <strong>struct</strong> is a typed collection of
                    fields &mdash; similar to a class in other languages, but
                    without methods baked in. Each field has three parts:
                  </p>
                  <ul className="mt-2 space-y-1.5 text-[13px] text-muted-foreground/70 leading-relaxed">
                    <li>
                      <strong>Name and type</strong>:{" "}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        Title string
                      </code>{" "}
                      &mdash; a field called Title of type string.
                    </li>
                    <li>
                      <strong>Struct tags</strong>: the backtick section after
                      the type. These are metadata used by libraries:
                      <ul className="mt-1 ml-4 space-y-1">
                        <li>
                          <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                            gorm:&quot;size:255&quot;
                          </code>{" "}
                          tells the ORM to create a VARCHAR(255) column.
                        </li>
                        <li>
                          <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                            json:&quot;title&quot;
                          </code>{" "}
                          controls the JSON key name when serializing.
                        </li>
                        <li>
                          <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                            binding:&quot;required&quot;
                          </code>{" "}
                          tells the Gin framework this field must be provided.
                        </li>
                      </ul>
                    </li>
                    <li>
                      <strong>Special types</strong>:{" "}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        *time.Time
                      </code>{" "}
                      (pointer) means the date can be null.{" "}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        gorm.DeletedAt
                      </code>{" "}
                      enables soft delete &mdash; records aren&apos;t actually
                      deleted, just hidden.
                    </li>
                  </ul>
                </div>
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
                  Understanding the Go Handler
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Grit also generated a handler with full CRUD endpoints.
                  Let&apos;s look at the Create handler:
                </p>

                <CodeBlock filename="apps/api/internal/handlers/task.go &mdash; Create" code={`func (h *TaskHandler) Create(c *gin.Context) {
    var req struct {
        Title       string     \`json:"title" binding:"required"\`
        Description string     \`json:"description"\`
        Status      string     \`json:"status" binding:"required"\`
        Priority    int        \`json:"priority"\`
        DueDate     *time.Time \`json:"due_date"\`
        Completed   bool       \`json:"completed"\`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": err.Error(),
            },
        })
        return
    }

    item := models.Task{
        Title:       req.Title,
        Description: req.Description,
        Status:      req.Status,
        Priority:    req.Priority,
        DueDate:     req.DueDate,
        Completed:   req.Completed,
    }

    if err := h.DB.Create(&item).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "INTERNAL_ERROR",
                "message": "Failed to create task",
            },
        })
        return
    }

    h.DB.First(&item, item.ID)

    c.JSON(http.StatusCreated, gin.H{
        "data":    item,
        "message": "Task created successfully",
    })
}`} />

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-4 mb-4">
                  <p className="text-xs font-semibold text-primary/80 uppercase tracking-wider mb-2">
                    Go Basics
                  </p>
                  <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                    Let&apos;s break down the Go patterns:
                  </p>
                  <ul className="mt-2 space-y-1.5 text-[13px] text-muted-foreground/70 leading-relaxed">
                    <li>
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        func (h *TaskHandler) Create(c *gin.Context)
                      </code>{" "}
                      &mdash; this is a <strong>method</strong> on the
                      TaskHandler struct. The{" "}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        (h *TaskHandler)
                      </code>{" "}
                      part is called a <strong>receiver</strong> &mdash;
                      it&apos;s how Go does methods. The{" "}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        *
                      </code>{" "}
                      means it&apos;s a pointer receiver, so we can access{" "}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        h.DB
                      </code>
                      .
                    </li>
                    <li>
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        c.ShouldBindJSON(&amp;req)
                      </code>{" "}
                      &mdash; Gin reads the request body and fills the{" "}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        req
                      </code>{" "}
                      struct. The{" "}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        &amp;
                      </code>{" "}
                      passes a pointer so Gin can modify it.
                    </li>
                    <li>
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        if err := ...; err != nil
                      </code>{" "}
                      &mdash; Go&apos;s standard error handling. Functions
                      return errors instead of throwing exceptions. You always
                      check if{" "}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        err
                      </code>{" "}
                      is not nil.
                    </li>
                    <li>
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        {"gin.H{}"}
                      </code>{" "}
                      &mdash; shorthand for{" "}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        {"map[string]interface{}{}"}
                      </code>
                      , used to build JSON responses.
                    </li>
                  </ul>
                </div>
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
                <h3 className="text-xl font-semibold mb-2">
                  See It in Action
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Restart the API server to pick up the new model, then open the
                  admin panel at{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    http://localhost:3001
                  </code>
                  . You&apos;ll see Tasks in the sidebar. Click it to see the
                  data table, then click &quot;Create Task&quot; to add your
                  first task.
                </p>

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  You can also test the API directly from the terminal:
                </p>

                <CodeBlock terminal code="curl http://localhost:8080/api/tasks | jq" className="glow-purple-sm mb-4" />

                <CodeBlock language="json" filename="API response" code={`{
  "data": [
    {
      "id": 1,
      "title": "Build the landing page",
      "description": "Design and implement the hero section",
      "status": "in-progress",
      "priority": 2,
      "due_date": "2026-03-01T00:00:00Z",
      "completed": false,
      "created_at": "2026-02-16T10:30:00Z",
      "updated_at": "2026-02-16T10:30:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "page_size": 20,
    "pages": 1
  }
}`} />

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-4 mb-4">
                  <p className="text-xs font-semibold text-primary/80 uppercase tracking-wider mb-2">
                    Go Basics
                  </p>
                  <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                    Notice how{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      json:&quot;title&quot;
                    </code>{" "}
                    in the Go struct maps directly to{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      &quot;title&quot;
                    </code>{" "}
                    in the JSON output. The{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      json:&quot;-&quot;
                    </code>{" "}
                    tag on DeletedAt means that field is hidden from API
                    responses entirely. This is how Go controls serialization
                    &mdash; through struct tags, not decorators or annotations.
                  </p>
                </div>
              </div>
            </div>

            {/* ============================================================ */}
            {/* MILESTONE 2 */}
            {/* ============================================================ */}
            <div className="mb-8">
              <h2 className="text-2xl font-semibold tracking-tight mb-2">
                Milestone 2: Adding Categories
              </h2>
              <p className="text-muted-foreground text-sm">
                Two models with a belongs_to relationship &mdash; foreign keys,
                eager loading, and relationship dropdowns.
              </p>
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
                  Generate the Category Resource
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Categories help organize tasks into groups. Generate a simple
                  Category resource with a name and a color:
                </p>

                <CodeBlock terminal code={`grit generate resource Category --fields "name:string,color:string"`} className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  Now we have a Category resource with just a name and a color.
                  Next, we&apos;ll connect Tasks to Categories.
                </p>
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
                  Add a belongs_to Relationship
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  In a real project, you&apos;d plan your data model upfront.
                  But Grit makes it easy to regenerate. Let&apos;s regenerate
                  the Task resource with a category relationship:
                </p>

                <CodeBlock terminal code={`grit generate resource Task --fields "title:string,description:text,category:belongs_to,status:string,priority:int,due_date:date,completed:bool"`} className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  The{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    category:belongs_to
                  </code>{" "}
                  syntax tells Grit to create a foreign key relationship.
                  Here&apos;s what changed in the Go model:
                </p>

                <CodeBlock filename="apps/api/internal/models/task.go &mdash; with relationship" code={`type Task struct {
    ID          uint           \`gorm:"primarykey" json:"id"\`
    Title       string         \`gorm:"size:255" json:"title" binding:"required"\`
    Description string         \`gorm:"type:text" json:"description"\`
    CategoryID  uint           \`gorm:"index" json:"category_id" binding:"required"\`
    Category    Category       \`gorm:"foreignKey:CategoryID" json:"category"\`
    Status      string         \`gorm:"size:255" json:"status" binding:"required"\`
    Priority    int            \`json:"priority"\`
    DueDate     *time.Time     \`gorm:"type:date" json:"due_date"\`
    Completed   bool           \`json:"completed"\`
    CreatedAt   time.Time      \`json:"created_at"\`
    UpdatedAt   time.Time      \`json:"updated_at"\`
    DeletedAt   gorm.DeletedAt \`gorm:"index" json:"-"\`
}`} />

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-4 mb-4">
                  <p className="text-xs font-semibold text-primary/80 uppercase tracking-wider mb-2">
                    Go Basics
                  </p>
                  <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                    Two new fields were added:
                  </p>
                  <ul className="mt-2 space-y-1.5 text-[13px] text-muted-foreground/70 leading-relaxed">
                    <li>
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        CategoryID uint
                      </code>{" "}
                      &mdash; the <strong>foreign key</strong> column. This
                      stores the numeric ID of the related category. The{" "}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        gorm:&quot;index&quot;
                      </code>{" "}
                      tag creates a database index for faster lookups.
                    </li>
                    <li>
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        Category Category
                      </code>{" "}
                      &mdash; the <strong>association</strong> struct. This
                      isn&apos;t a database column &mdash; it&apos;s a Go-only
                      field that GORM populates when you use{" "}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        Preload
                      </code>
                      . The{" "}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        gorm:&quot;foreignKey:CategoryID&quot;
                      </code>{" "}
                      tag tells GORM which column links the two models.
                    </li>
                  </ul>
                  <p className="mt-2 text-[13px] text-muted-foreground/70 leading-relaxed">
                    Think of{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      CategoryID
                    </code>{" "}
                    as the actual data (stored in the database) and{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      Category
                    </code>{" "}
                    as the convenient Go object (loaded on demand).
                  </p>
                </div>

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  The handler now includes{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    .Preload(&quot;Category&quot;)
                  </code>{" "}
                  to eager-load the related category data:
                </p>

                <CodeBlock filename="handler snippet" code={`// List — eager loading with Preload
query := h.DB.Model(&models.Task{}).Preload("Category")

// GetByID — also preloads
h.DB.Preload("Category").First(&item, id)`} />
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
                  Test the Relationship
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  First, create a few categories in the admin panel (e.g.,
                  &quot;Work&quot;, &quot;Personal&quot;, &quot;Urgent&quot;).
                  Then create or edit a task &mdash; you&apos;ll see a Category
                  dropdown that fetches options from the API automatically.
                </p>

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  The API response now includes the nested category:
                </p>

                <CodeBlock terminal code="curl http://localhost:8080/api/tasks/1 | jq" className="glow-purple-sm mb-4" />

                <CodeBlock language="json" filename="API response" code={`{
  "data": {
    "id": 1,
    "title": "Build the landing page",
    "description": "Design and implement the hero section",
    "category_id": 1,
    "category": {
      "id": 1,
      "name": "Work",
      "color": "#6c5ce7",
      "created_at": "2026-02-16T10:00:00Z",
      "updated_at": "2026-02-16T10:00:00Z"
    },
    "status": "in-progress",
    "priority": 2,
    "completed": false
  }
}`} />

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-4 mb-4">
                  <p className="text-xs font-semibold text-primary/80 uppercase tracking-wider mb-2">
                    Go Basics
                  </p>
                  <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                    GORM&apos;s{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      Preload
                    </code>{" "}
                    works by running a second query behind the scenes: first it
                    fetches the tasks, then it fetches all related categories in
                    one batch query (
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      WHERE id IN (1, 2, 3)
                    </code>
                    ) and stitches them together in Go. This is much more
                    efficient than N+1 queries.
                  </p>
                </div>
              </div>
            </div>

            {/* ============================================================ */}
            {/* MILESTONE 3 */}
            {/* ============================================================ */}
            <div className="mb-8">
              <h2 className="text-2xl font-semibold tracking-tight mb-2">
                Milestone 3: Tags &amp; Comments
              </h2>
              <p className="text-muted-foreground text-sm">
                Adding many_to_many relationships, multiple belongs_to, and
                seeing the full data model come together.
              </p>
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
                  Generate Tags with many_to_many
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  First, generate a simple Tag resource:
                </p>

                <CodeBlock terminal code={`grit generate resource Tag --fields "name:string:unique"`} className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Tags and tasks have a many-to-many relationship &mdash; a task
                  can have multiple tags, and a tag can be on multiple tasks.
                  Let&apos;s regenerate Task to add this:
                </p>

                <CodeBlock terminal code={`grit generate resource Task --fields "title:string,description:text,category:belongs_to,tags:many_to_many:Tag,status:string,priority:int,due_date:date,completed:bool"`} className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  The{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    tags:many_to_many:Tag
                  </code>{" "}
                  syntax creates a junction table. Here&apos;s the updated
                  model:
                </p>

                <CodeBlock filename="apps/api/internal/models/task.go &mdash; with m2m" code={`type Task struct {
    ID          uint           \`gorm:"primarykey" json:"id"\`
    Title       string         \`gorm:"size:255" json:"title" binding:"required"\`
    Description string         \`gorm:"type:text" json:"description"\`
    CategoryID  uint           \`gorm:"index" json:"category_id" binding:"required"\`
    Category    Category       \`gorm:"foreignKey:CategoryID" json:"category"\`
    Tags        []Tag          \`gorm:"many2many:task_tags" json:"tags"\`
    Status      string         \`gorm:"size:255" json:"status" binding:"required"\`
    Priority    int            \`json:"priority"\`
    DueDate     *time.Time     \`gorm:"type:date" json:"due_date"\`
    Completed   bool           \`json:"completed"\`
    CreatedAt   time.Time      \`json:"created_at"\`
    UpdatedAt   time.Time      \`json:"updated_at"\`
    DeletedAt   gorm.DeletedAt \`gorm:"index" json:"-"\`
}`} />

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-4 mb-4">
                  <p className="text-xs font-semibold text-primary/80 uppercase tracking-wider mb-2">
                    Go Basics
                  </p>
                  <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                    The{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      Tags []Tag
                    </code>{" "}
                    field with{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      gorm:&quot;many2many:task_tags&quot;
                    </code>{" "}
                    tells GORM to create a <strong>junction table</strong>{" "}
                    called{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      task_tags
                    </code>{" "}
                    with columns{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      task_id
                    </code>{" "}
                    and{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      tag_id
                    </code>
                    . Unlike{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      belongs_to
                    </code>
                    , there&apos;s no foreign key column on the Task table itself
                    &mdash; the relationship lives entirely in the junction
                    table.
                  </p>
                  <p className="mt-2 text-[13px] text-muted-foreground/70 leading-relaxed">
                    In the handler, GORM uses{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      Association(&quot;Tags&quot;).Replace(tags)
                    </code>{" "}
                    to manage the junction table entries. The{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      Replace
                    </code>{" "}
                    method handles both adding and removing associations.
                  </p>
                </div>
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
                <h3 className="text-xl font-semibold mb-2">Add Comments</h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  A Comment belongs to both a Task and a User (the built-in User
                  model). The{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    task:belongs_to:Task
                  </code>{" "}
                  syntax uses the explicit form &mdash; the field name is{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    task
                  </code>{" "}
                  and the related model is{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    Task
                  </code>
                  .
                </p>

                <CodeBlock terminal code={`grit generate resource Comment --fields "content:text,task:belongs_to:Task,author:belongs_to:User"`} className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Here is the generated Comment model:
                </p>

                <CodeBlock filename="apps/api/internal/models/comment.go" code={`type Comment struct {
    ID        uint           \`gorm:"primarykey" json:"id"\`
    Content   string         \`gorm:"type:text" json:"content"\`
    TaskID    uint           \`gorm:"index" json:"task_id" binding:"required"\`
    Task      Task           \`gorm:"foreignKey:TaskID" json:"task"\`
    AuthorID  uint           \`gorm:"index" json:"author_id" binding:"required"\`
    Author    User           \`gorm:"foreignKey:AuthorID" json:"author"\`
    CreatedAt time.Time      \`json:"created_at"\`
    UpdatedAt time.Time      \`json:"updated_at"\`
    DeletedAt gorm.DeletedAt \`gorm:"index" json:"-"\`
}`} />

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-4 mb-4">
                  <p className="text-xs font-semibold text-primary/80 uppercase tracking-wider mb-2">
                    Go Basics
                  </p>
                  <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                    A model can have multiple{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      belongs_to
                    </code>{" "}
                    relationships. Each one creates its own foreign key (
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      TaskID
                    </code>
                    ,{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      AuthorID
                    </code>
                    ) and association struct (
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      Task
                    </code>
                    ,{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      Author
                    </code>
                    ). The{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      author:belongs_to:User
                    </code>{" "}
                    syntax is the explicit form &mdash; you need it when the
                    field name (
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      author
                    </code>
                    ) differs from the model name (
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      User
                    </code>
                    ).
                  </p>
                </div>
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
                  The Full Data Model
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  You now have 4 models with 3 types of relationships:
                </p>

                <CodeBlock filename="data model diagram" code={`Category ──< Task >── Tag
               │        (many-to-many via task_tags)
               │
            Comment
               │
              User

── belongs_to (foreign key on child)
>── many_to_many (junction table)`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  Create some tags (&quot;frontend&quot;, &quot;backend&quot;,
                  &quot;urgent&quot;), add them to tasks, and create comments.
                  The admin panel handles everything through the generated forms
                  &mdash; no custom code required.
                </p>
              </div>
            </div>

            {/* ============================================================ */}
            {/* MILESTONE 4 */}
            {/* ============================================================ */}
            <div className="mb-8">
              <h2 className="text-2xl font-semibold tracking-tight mb-2">
                Milestone 4: Building the Web App
              </h2>
              <p className="text-muted-foreground text-sm">
                Fetch and display data from the Next.js frontend &mdash; list
                tasks, view details, and create new ones.
              </p>
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
                  Set Up the API Client
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  The web app at{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    apps/web
                  </code>{" "}
                  is a Next.js project. Let&apos;s create a simple API helper:
                </p>

                <CodeBlock language="typescript" filename="apps/web/lib/api.ts" code={`const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function fetchAPI<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(\`\${API_URL}\${path}\`, {
    headers: { "Content-Type": "application/json", ...options?.headers },
    ...options,
  });
  if (!res.ok) throw new Error(\`API error: \${res.status}\`);
  return res.json();
}`} />

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-4 mb-4">
                  <p className="text-xs font-semibold text-primary/80 uppercase tracking-wider mb-2">
                    React Basics
                  </p>
                  <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                    We&apos;re using a simple{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      fetch
                    </code>{" "}
                    wrapper instead of a full library. The{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      NEXT_PUBLIC_
                    </code>{" "}
                    prefix makes the env variable available in the browser. For a
                    production app, you&apos;d use React Query (which Grit
                    generates for the admin panel) &mdash; but plain fetch is a
                    good way to learn the fundamentals.
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
                <h3 className="text-xl font-semibold mb-2">
                  Create a Task List Page
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Let&apos;s build a page that shows all tasks with their
                  categories and tags:
                </p>

                <CodeBlock language="tsx" filename="apps/web/app/tasks/page.tsx" code={`"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { fetchAPI } from "@/lib/api";

interface Task {
  id: number;
  title: string;
  description: string;
  status: string;
  priority: number;
  completed: boolean;
  category?: { id: number; name: string; color: string };
  tags?: { id: number; name: string }[];
  created_at: string;
}

export default function TasksPage() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchAPI<{ data: Task[] }>("/api/tasks?page_size=50")
      .then((res) => setTasks(res.data))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="p-8 text-center">Loading tasks...</div>;

  return (
    <div className="max-w-4xl mx-auto p-8">
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-3xl font-bold">Tasks</h1>
        <Link
          href="/tasks/new"
          className="rounded-lg bg-purple-600 px-4 py-2 text-sm font-medium text-white hover:bg-purple-700"
        >
          New Task
        </Link>
      </div>

      <div className="space-y-4">
        {tasks.map((task) => (
          <Link
            key={task.id}
            href={\`/tasks/\${task.id}\`}
            className="block rounded-xl border border-gray-800 bg-gray-900 p-5 hover:border-purple-500/50 transition-colors"
          >
            <div className="flex items-start justify-between">
              <div>
                <h2 className="text-lg font-semibold">{task.title}</h2>
                {task.description && (
                  <p className="mt-1 text-sm text-gray-400 line-clamp-2">{task.description}</p>
                )}
                <div className="mt-3 flex items-center gap-2">
                  {task.category && (
                    <span className="rounded-full bg-purple-500/10 border border-purple-500/20 px-2.5 py-0.5 text-xs text-purple-400">
                      {task.category.name}
                    </span>
                  )}
                  {task.tags?.map((tag) => (
                    <span key={tag.id} className="rounded-full bg-gray-800 px-2.5 py-0.5 text-xs text-gray-400">
                      {tag.name}
                    </span>
                  ))}
                </div>
              </div>
              <span className={\`rounded-full px-2.5 py-0.5 text-xs font-medium \${
                task.completed ? "bg-green-500/10 text-green-400" : "bg-yellow-500/10 text-yellow-400"
              }\`}>
                {task.completed ? "Done" : task.status}
              </span>
            </div>
          </Link>
        ))}
      </div>
    </div>
  );
}`} />

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-4 mb-4">
                  <p className="text-xs font-semibold text-primary/80 uppercase tracking-wider mb-2">
                    React Basics
                  </p>
                  <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      &quot;use client&quot;
                    </code>{" "}
                    tells Next.js this is a client component (it runs in the
                    browser). We use{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      useState
                    </code>{" "}
                    to store tasks and loading state, and{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      useEffect
                    </code>{" "}
                    to fetch data when the page loads. The{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      fetchAPI
                    </code>{" "}
                    call hits our Go backend, which returns the paginated task
                    list with preloaded categories and tags.
                  </p>
                  <p className="mt-2 text-[13px] text-muted-foreground/70 leading-relaxed">
                    Note: In a production app, you&apos;d use React Query for
                    caching, background refetching, and optimistic updates. The
                    admin panel already uses it &mdash; this vanilla approach is
                    just to show the fundamentals.
                  </p>
                </div>
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
                  Create a Task Detail Page
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Now let&apos;s add a detail page that shows a single task with
                  its comments:
                </p>

                <CodeBlock language="tsx" filename="apps/web/app/tasks/[id]/page.tsx" code={`"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { fetchAPI } from "@/lib/api";

interface Comment {
  id: number;
  content: string;
  author?: { first_name: string; last_name: string };
  created_at: string;
}

interface Task {
  id: number;
  title: string;
  description: string;
  status: string;
  priority: number;
  completed: boolean;
  due_date: string | null;
  category?: { id: number; name: string; color: string };
  tags?: { id: number; name: string }[];
  created_at: string;
}

export default function TaskDetailPage() {
  const { id } = useParams();
  const [task, setTask] = useState<Task | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);

  useEffect(() => {
    fetchAPI<{ data: Task }>(\`/api/tasks/\${id}\`).then((res) => setTask(res.data));
    fetchAPI<{ data: Comment[] }>(\`/api/comments?search=\${id}\`).then((res) => setComments(res.data));
  }, [id]);

  if (!task) return <div className="p-8 text-center">Loading...</div>;

  return (
    <div className="max-w-3xl mx-auto p-8">
      <Link href="/tasks" className="text-sm text-purple-400 hover:text-purple-300 mb-6 inline-block">
        &larr; Back to tasks
      </Link>

      <div className="rounded-xl border border-gray-800 bg-gray-900 p-6">
        <div className="flex items-start justify-between">
          <h1 className="text-2xl font-bold">{task.title}</h1>
          <span className={\`rounded-full px-3 py-1 text-xs font-medium \${
            task.completed ? "bg-green-500/10 text-green-400" : "bg-yellow-500/10 text-yellow-400"
          }\`}>
            {task.completed ? "Completed" : task.status}
          </span>
        </div>

        {task.description && (
          <p className="mt-4 text-gray-400 leading-relaxed">{task.description}</p>
        )}

        <div className="mt-6 flex flex-wrap items-center gap-3">
          {task.category && (
            <span className="rounded-full bg-purple-500/10 border border-purple-500/20 px-3 py-1 text-sm text-purple-400">
              {task.category.name}
            </span>
          )}
          {task.tags?.map((tag) => (
            <span key={tag.id} className="rounded-full bg-gray-800 px-3 py-1 text-sm text-gray-400">
              {tag.name}
            </span>
          ))}
        </div>

        <div className="mt-6 flex items-center gap-6 text-sm text-gray-500">
          <span>Priority: {task.priority}</span>
          {task.due_date && <span>Due: {new Date(task.due_date).toLocaleDateString()}</span>}
          <span>Created: {new Date(task.created_at).toLocaleDateString()}</span>
        </div>
      </div>

      {/* Comments */}
      <div className="mt-8">
        <h2 className="text-xl font-semibold mb-4">Comments ({comments.length})</h2>
        <div className="space-y-3">
          {comments.map((comment) => (
            <div key={comment.id} className="rounded-lg border border-gray-800 bg-gray-900/50 p-4">
              <p className="text-sm text-gray-300">{comment.content}</p>
              <p className="mt-2 text-xs text-gray-500">
                {comment.author ? \`\${comment.author.first_name} \${comment.author.last_name}\` : "Unknown"} &middot;{" "}
                {new Date(comment.created_at).toLocaleDateString()}
              </p>
            </div>
          ))}
          {comments.length === 0 && (
            <p className="text-sm text-gray-500">No comments yet. Add one from the admin panel.</p>
          )}
        </div>
      </div>
    </div>
  );
}`} />
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
                  Create a New Task Form
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Finally, let&apos;s build a form that creates tasks by POSTing
                  to the Go API:
                </p>

                <CodeBlock language="tsx" filename="apps/web/app/tasks/new/page.tsx" code={`"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { fetchAPI } from "@/lib/api";

interface Category {
  id: number;
  name: string;
}

export default function NewTaskPage() {
  const router = useRouter();
  const [categories, setCategories] = useState<Category[]>([]);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    fetchAPI<{ data: Category[] }>("/api/categories?page_size=100")
      .then((res) => setCategories(res.data));
  }, []);

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setSubmitting(true);

    const form = new FormData(e.currentTarget);
    const body = {
      title: form.get("title") as string,
      description: form.get("description") as string,
      category_id: Number(form.get("category_id")),
      status: form.get("status") as string || "todo",
      priority: Number(form.get("priority") || 0),
    };

    try {
      await fetchAPI("/api/tasks", {
        method: "POST",
        body: JSON.stringify(body),
      });
      router.push("/tasks");
    } catch (err) {
      alert("Failed to create task");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="max-w-2xl mx-auto p-8">
      <h1 className="text-3xl font-bold mb-8">New Task</h1>

      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label className="block text-sm font-medium mb-2">Title *</label>
          <input
            name="title"
            required
            className="w-full rounded-lg border border-gray-700 bg-gray-900 px-4 py-2.5 text-sm focus:border-purple-500 focus:outline-none"
            placeholder="What needs to be done?"
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-2">Description</label>
          <textarea
            name="description"
            rows={4}
            className="w-full rounded-lg border border-gray-700 bg-gray-900 px-4 py-2.5 text-sm focus:border-purple-500 focus:outline-none"
            placeholder="Add some details..."
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium mb-2">Category</label>
            <select
              name="category_id"
              className="w-full rounded-lg border border-gray-700 bg-gray-900 px-4 py-2.5 text-sm focus:border-purple-500 focus:outline-none"
            >
              <option value="">Select category</option>
              {categories.map((cat) => (
                <option key={cat.id} value={cat.id}>{cat.name}</option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">Priority</label>
            <input
              name="priority"
              type="number"
              min={0}
              max={5}
              defaultValue={0}
              className="w-full rounded-lg border border-gray-700 bg-gray-900 px-4 py-2.5 text-sm focus:border-purple-500 focus:outline-none"
            />
          </div>
        </div>

        <button
          type="submit"
          disabled={submitting}
          className="rounded-lg bg-purple-600 px-6 py-2.5 text-sm font-medium text-white hover:bg-purple-700 disabled:opacity-50"
        >
          {submitting ? "Creating..." : "Create Task"}
        </button>
      </form>
    </div>
  );
}`} />

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-4 mb-4">
                  <p className="text-xs font-semibold text-primary/80 uppercase tracking-wider mb-2">
                    React Basics
                  </p>
                  <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                    This form demonstrates the fundamentals of creating data
                    from React:
                  </p>
                  <ul className="mt-2 space-y-1.5 text-[13px] text-muted-foreground/70 leading-relaxed">
                    <li>
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        e.preventDefault()
                      </code>{" "}
                      stops the browser from doing a full page reload.
                    </li>
                    <li>
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        FormData
                      </code>{" "}
                      extracts values from form inputs by their{" "}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        name
                      </code>{" "}
                      attributes.
                    </li>
                    <li>
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        JSON.stringify(body)
                      </code>{" "}
                      converts the JavaScript object to a JSON string for the
                      API.
                    </li>
                    <li>
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                        router.push(&quot;/tasks&quot;)
                      </code>{" "}
                      navigates to the task list after successful creation.
                    </li>
                  </ul>
                  <p className="mt-2 text-[13px] text-muted-foreground/70 leading-relaxed">
                    The Go API validates the request body using the{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      binding:&quot;required&quot;
                    </code>{" "}
                    tags. If{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      title
                    </code>{" "}
                    is missing, it returns a 422 error with a validation
                    message.
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
                <h3 className="text-xl font-semibold mb-2">
                  Run Everything Together
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Start the web frontend and test the full flow:
                </p>

                <CodeBlock terminal code="cd apps/web && pnpm dev" className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Now open{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    http://localhost:3000/tasks
                  </code>{" "}
                  in your browser. You should see:
                </p>

                <ul className="space-y-2 mb-4">
                  {[
                    "A list of all tasks with category badges and tags",
                    "Click a task to see its details and comments",
                    "Click \"New Task\" to create one from the web app",
                    "Check the admin panel — your new task appears there too",
                  ].map((item, i) => (
                    <li
                      key={i}
                      className="flex items-start gap-2.5 text-[13px] text-muted-foreground/70"
                    >
                      <span className="text-primary/70 font-mono text-xs mt-0.5">
                        {i + 1}.
                      </span>
                      {item}
                    </li>
                  ))}
                </ul>

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  Congratulations! You&apos;ve built a full-stack application
                  with 4 related models, a Go API with eager loading, an admin
                  panel with relationship forms, and a Next.js frontend that
                  fetches and creates data.
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
                  "A full-stack task manager with Go API and Next.js frontend",
                  "4 resources (Task, Category, Tag, Comment) generated with CLI commands",
                  "A belongs_to relationship between Tasks and Categories",
                  "A many_to_many relationship between Tasks and Tags via a junction table",
                  "Comments with multiple belongs_to relationships (Task + User)",
                  "An admin panel with data tables, forms, relationship selects, and multi-selects",
                  "A web app with task list, detail page, and create form",
                  "Type-safe API responses with nested related data via GORM Preload",
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
                Now that you have a working task manager, here are some ways to
                extend it:
              </p>
              <ul>
                <li>
                  <strong>Add authentication to the web app</strong> &mdash;
                  protect the /tasks routes with JWT tokens from the Go API.
                </li>
                <li>
                  <strong>Use grit sync</strong> &mdash; run{" "}
                  <code>grit sync</code> to automatically keep Go types and
                  TypeScript types in sync as you change models.
                </li>
                <li>
                  <strong>File uploads</strong> &mdash; add an attachment field
                  to tasks using Grit&apos;s S3 storage service.
                </li>
                <li>
                  <strong>Email notifications</strong> &mdash; send an email
                  when a task is completed using the background jobs system.
                </li>
                <li>
                  <strong>Try the Blog tutorial</strong> &mdash; build a
                  complete blogging platform with custom public endpoints and
                  frontend pages.
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
                <Link href="/docs/design/theme" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Theme &amp; Colors
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
