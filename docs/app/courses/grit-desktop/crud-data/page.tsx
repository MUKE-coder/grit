import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import { LaneFlow } from '@/components/lane-flow'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Desktop CRUD & Data — Grit Desktop Course',
  description: 'Generate resources, understand Wails bindings, GORM models with SQLite, TanStack Router, and the service layer in Grit desktop apps.',
}

export default function CrudDataCourse() {
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
          <span className="text-foreground">Desktop CRUD & Data</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Course 2 of 5</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">12 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Desktop CRUD & Data
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            In this course, you will learn how CRUD operations work in desktop apps, generate resources
            with the Grit code generator, understand Wails bindings in depth, and work with SQLite through GORM.
          </p>
          <LaneFlow id="c-gdcrud" lanes={['grit generate', 'Wails bindings', 'SQLite']}
            nodes={[
              { id: 'gen', lane: 0, row: 1, title: 'generate resource', sub: 'CRUD', tone: 'primary' },
              { id: 'bind', lane: 1, row: 1, title: 'Wails bindings', sub: 'typed Go calls', tone: 'cyan' },
              { id: 'gorm', lane: 2, row: 0, title: 'GORM', sub: 'queries', tone: 'blue' },
              { id: 'sqlite', lane: 2, row: 1, title: 'SQLite', sub: 'local DB', tone: 'green' },
            ]}
            edges={[{ from: 'gen', to: 'bind', label: 'wires', tone: 'cyan' }, { from: 'bind', to: 'gorm', tone: 'blue' }, { from: 'gorm', to: 'sqlite', tone: 'green' }]}
            legend={[{ tone: 'cyan', label: 'Wails bindings' }, { tone: 'green', label: 'Local SQLite' }]}
            caption="Generated CRUD calls Go through Wails bindings, straight to SQLite via GORM" />
        </div>

        <div className="my-4 rounded-lg border border-primary/20 bg-primary/5 px-4 py-3">
          <p className="text-sm text-muted-foreground">
            <strong className="text-foreground">Reference docs:</strong>{' '}
            <a href="/docs/desktop/resource-generation" className="text-primary hover:underline">Desktop Resource Generation →</a>
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ How Desktop CRUD Differs ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">How Desktop CRUD Differs</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            In <strong className="text-foreground">web apps</strong>, React calls REST APIs over HTTP. You
            use <Code>fetch()</Code> or axios to send requests to endpoints like <Code>GET /api/blogs</Code>.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            In <strong className="text-foreground">desktop apps</strong>, React calls Go functions directly
            via Wails bindings — no REST endpoints to define and no <Code>fetch()</Code> calls to write by hand.
            (Under the hood, Grit runs an embedded Gin API bound to <Code>127.0.0.1:34999</Code> that the bindings
            talk to, but you never touch the HTTP layer yourself.) From your React code it{"'"}s just function calls:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Operation</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Web App (HTTP)</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Desktop App (Wails)</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">List</td>
                  <td className="px-4 py-3"><Code>{'fetch("/api/tasks")'}</Code></td>
                  <td className="px-4 py-3"><Code>GetTasks()</Code></td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Create</td>
                  <td className="px-4 py-3"><Code>{'fetch("/api/tasks", { method: "POST" })'}</Code></td>
                  <td className="px-4 py-3"><Code>{'CreateTask({ title: "..." })'}</Code></td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Update</td>
                  <td className="px-4 py-3"><Code>{'fetch("/api/tasks/1", { method: "PUT" })'}</Code></td>
                  <td className="px-4 py-3"><Code>{'UpdateTask(1, { title: "..." })'}</Code></td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Delete</td>
                  <td className="px-4 py-3"><Code>{'fetch("/api/tasks/1", { method: "DELETE" })'}</Code></td>
                  <td className="px-4 py-3"><Code>DeleteTask(1)</Code></td>
                </tr>
              </tbody>
            </table>
          </div>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This is simpler and faster — no hand-written serialization, no CORS configuration to manage.
            Your React code just calls Go functions and gets results back; the Wails bindings and the
            embedded local API handle the plumbing.
          </p>
        </section>

        {/* ═══ Generating Resources ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Generating Resources</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>grit generate resource</Code> command works for desktop projects too — same syntax, same workflow,
            but different output. For desktop, it generates Wails bindings instead of HTTP handlers, and TanStack Router
            routes instead of Next.js pages:
          </p>

          <CodeBlock filename="Terminal">
{`grit generate resource Task --fields "title:string,description:text,done:bool,priority:int"`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This creates everything you need for a full Task CRUD:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">File</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Purpose</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">internal/models/task.go</td>
                  <td className="px-4 py-3">GORM model with struct tags</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">internal/service/task.go</td>
                  <td className="px-4 py-3">Business logic (List, Create, GetByID, Update, Delete)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">internal/types/task_types.go</td>
                  <td className="px-4 py-3">Input/output types for Wails bindings</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">app.go (injected)</td>
                  <td className="px-4 py-3">New bound methods added to the App struct</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">frontend/src/routes/tasks/</td>
                  <td className="px-4 py-3">TanStack Router pages (list, create, edit)</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-mono text-foreground text-xs">frontend/src/components/sidebar.tsx</td>
                  <td className="px-4 py-3">Sidebar entry injected automatically</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Challenge number={1} title="Generate a Task Resource">
            <p>Run <Code>grit generate resource Task --fields {'"'}title:string,description:text,done:bool,priority:int{'"'}</Code> in your desktop project. List the files that were created. Did the sidebar update with a new {'"'}Tasks{'"'} entry?</p>
          </Challenge>
        </section>

        {/* ═══ Understanding Wails Bindings ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Understanding Wails Bindings</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When you generate a resource, Grit adds methods to <Code>app.go</Code> that become the bridge
            between Go and React:
          </p>

          <CodeBlock filename="app.go (Task methods)">
{`// List all tasks
func (a *App) GetTasks() ([]models.Task, error) {
    return a.taskService.List()
}

// Create a new task
func (a *App) CreateTask(input types.CreateTaskInput) (*models.Task, error) {
    return a.taskService.Create(input)
}

// Get a single task by ID
func (a *App) GetTask(id string) (*models.Task, error) {
    return a.taskService.GetByID(id)
}

// Update an existing task
func (a *App) UpdateTask(id string, input types.UpdateTaskInput) (*models.Task, error) {
    return a.taskService.Update(id, input)
}

// Delete a task
func (a *App) DeleteTask(id string) error {
    return a.taskService.Delete(id)
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            On the React side, Wails automatically generates TypeScript wrappers. You import and call
            them like any other function:
          </p>

          <CodeBlock filename="frontend/src/routes/tasks/index.tsx">
{`import { GetTasks, CreateTask, DeleteTask } from '../../../wailsjs/go/main/App'

// Fetch all tasks
const tasks = await GetTasks()

// Create a new task
const newTask = await CreateTask({
  title: "Buy groceries",
  description: "Milk, eggs, bread",
  done: false,
  priority: 1
})

// Delete a task
await DeleteTask(newTask.ID)`}
          </CodeBlock>

          <Definition term="Wails Bindings">
            Auto-generated TypeScript wrappers for your Go methods. Wails reads your Go code, finds
            all public methods on bound structs, and creates matching TypeScript functions in the
            <Code>wailsjs/go/</Code> directory. These functions handle serialization and deserialization
            automatically — you pass JavaScript objects and get JavaScript objects back.
          </Definition>

          <Challenge number={2} title="Explore the Generated Bindings">
            <p>After generating the Task resource, look in <Code>wailsjs/go/main/App.ts</Code> (this file is generated when you run <Code>wails dev</Code>). What functions are available? Do they match the Go methods in <Code>app.go</Code>?</p>
          </Challenge>
        </section>

        {/* ═══ GORM Models for Desktop ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">GORM Models for Desktop</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Desktop models use the same GORM patterns as web projects, but with the SQLite driver
            instead of PostgreSQL:
          </p>

          <CodeBlock filename="internal/models/task.go">
{`package models

import "gorm.io/gorm"

type Task struct {
    gorm.Model
    Title       string ` + '`' + `gorm:"not null" json:"title"` + '`' + `
    Description string ` + '`' + `gorm:"type:text" json:"description"` + '`' + `
    Done        bool   ` + '`' + `gorm:"default:false" json:"done"` + '`' + `
    Priority    int    ` + '`' + `gorm:"default:0" json:"priority"` + '`' + `
}`}
          </CodeBlock>

          <Note>
            SQLite has some differences from PostgreSQL. There is no <Code>JSONB</Code> type — use <Code>TEXT</Code> for
            storing JSON strings. There are no arrays — use <Code>TEXT</Code> with comma-separated values or JSON.
            Most standard types (string, int, bool, float, time) work identically.
          </Note>

          <Challenge number={3} title="Compare Models">
            <p>Open the Task model at <Code>internal/models/task.go</Code>. Compare it to a web project model if you have one. What{"'"}s the same? The struct tags, the GORM model embedding, and the JSON tags should all look familiar.</p>
          </Challenge>
        </section>

        {/* ═══ TanStack Router for Desktop ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">TanStack Router for Desktop</h2>

          <Definition term="TanStack Router">
            A type-safe router for React applications. It uses file-based routing — each file in the
            routes directory becomes a URL. TanStack Router is the router used in Grit desktop apps
            (Next.js App Router is used for web apps).
          </Definition>

          <Definition term="File-based Routing">
            A convention where the file structure in your routes folder directly maps to URL paths.
            A file at <Code>routes/tasks/index.tsx</Code> becomes the <Code>/tasks</Code> page.
            A file at <Code>routes/tasks/$id.edit.tsx</Code> becomes <Code>/tasks/123/edit</Code>.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Desktop apps use TanStack Router with hash history for Wails compatibility. The route files
            live in <Code>frontend/src/routes/</Code>:
          </p>

          <CodeBlock filename="Route Structure">
{`frontend/src/routes/
├── __root.tsx           <- Root layout (sidebar + title bar)
├── index.tsx            <- Dashboard (/)
├── tasks/
│   ├── index.tsx        <- Task list (/tasks)
│   ├── new.tsx          <- Create task (/tasks/new)
│   └── $id.edit.tsx     <- Edit task (/tasks/123/edit)
├── blogs/
│   ├── index.tsx        <- Blog list (/blogs)
│   ├── new.tsx          <- Create blog (/blogs/new)
│   └── $id.edit.tsx     <- Edit blog (/blogs/123/edit)
└── contacts/
    ├── index.tsx
    ├── new.tsx
    └── $id.edit.tsx`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>$id</Code> prefix means a dynamic parameter — TanStack Router extracts the value from
            the URL and makes it available in your component. So <Code>/tasks/42/edit</Code> gives you <Code>id = 42</Code>.
          </p>

          <Challenge number={4} title="Explore Route Files">
            <p>Find the route files for your generated Task resource in <Code>frontend/src/routes/tasks/</Code>. What URL patterns do they follow? Open <Code>index.tsx</Code> — can you see where it calls <Code>GetTasks()</Code>?</p>
          </Challenge>
        </section>

        {/* ═══ The Service Layer ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The Service Layer</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Services contain all business logic. The <Code>App</Code> methods in <Code>app.go</Code> call services,
            and services call GORM. This is the same pattern used in web projects — keeping logic
            separate from the transport layer:
          </p>

          <CodeBlock filename="internal/service/task.go (simplified)">
{`type TaskService struct {
    db *gorm.DB
}

func NewTaskService(db *gorm.DB) *TaskService {
    return &TaskService{db: db}
}

func (s *TaskService) List() ([]models.Task, error) {
    var tasks []models.Task
    result := s.db.Order("created_at DESC").Find(&tasks)
    return tasks, result.Error
}

func (s *TaskService) Create(input types.CreateTaskInput) (*models.Task, error) {
    task := models.Task{
        Title:       input.Title,
        Description: input.Description,
        Done:        input.Done,
        Priority:    input.Priority,
    }
    result := s.db.Create(&task)
    return &task, result.Error
}

func (s *TaskService) GetByID(id string) (*models.Task, error) {
    var task models.Task
    result := s.db.First(&task, id)
    return &task, result.Error
}

func (s *TaskService) Update(id string, input types.UpdateTaskInput) (*models.Task, error) {
    var task models.Task
    if err := s.db.First(&task, id).Error; err != nil {
        return nil, err
    }
    result := s.db.Model(&task).Updates(input)
    return &task, result.Error
}

func (s *TaskService) Delete(id string) error {
    return s.db.Delete(&models.Task{}, id).Error
}`}
          </CodeBlock>

          <Tip>
            The service layer is where you add custom business rules. Need to validate that priority
            is between 1 and 5? Add it in the service{"'"}s <Code>Create()</Code> method. Need to send a notification
            when a task is completed? Add it in <Code>Update()</Code>. Keep <Code>app.go</Code> thin — it just delegates.
          </Tip>

          <Challenge number={5} title="Read the Service">
            <p>Open the task service at <Code>internal/service/task.go</Code>. Read the <Code>List()</Code> method. What GORM query does it use? What order does it sort by?</p>
          </Challenge>
        </section>

        {/* ═══ Data Types and Field Mapping ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Data Types and Field Mapping</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When you define fields in <Code>grit generate resource</Code>, each type maps through three layers:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Grit Type</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Go Type</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">SQLite Type</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">TypeScript Type</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">string</td>
                  <td className="px-4 py-3"><Code>string</Code></td>
                  <td className="px-4 py-3">TEXT</td>
                  <td className="px-4 py-3"><Code>string</Code></td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">text</td>
                  <td className="px-4 py-3"><Code>string</Code></td>
                  <td className="px-4 py-3">TEXT</td>
                  <td className="px-4 py-3"><Code>string</Code></td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">int</td>
                  <td className="px-4 py-3"><Code>int</Code></td>
                  <td className="px-4 py-3">INTEGER</td>
                  <td className="px-4 py-3"><Code>number</Code></td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">float</td>
                  <td className="px-4 py-3"><Code>float64</Code></td>
                  <td className="px-4 py-3">REAL</td>
                  <td className="px-4 py-3"><Code>number</Code></td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">bool</td>
                  <td className="px-4 py-3"><Code>bool</Code></td>
                  <td className="px-4 py-3">INTEGER (0/1)</td>
                  <td className="px-4 py-3"><Code>boolean</Code></td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">date</td>
                  <td className="px-4 py-3"><Code>*time.Time</Code></td>
                  <td className="px-4 py-3">DATETIME</td>
                  <td className="px-4 py-3"><Code>string</Code> (ISO 8601)</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-mono text-foreground text-xs">richtext</td>
                  <td className="px-4 py-3"><Code>string</Code></td>
                  <td className="px-4 py-3">TEXT</td>
                  <td className="px-4 py-3"><Code>string</Code> (HTML)</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Challenge number={6} title="Type Mapping Quiz">
            <p>Without looking at the table: What SQLite type does a <Code>float64</Code> become? What about <Code>*time.Time</Code>? What Go type does <Code>bool</Code> map to? Check your answers against the table above.</p>
          </Challenge>
        </section>

        {/* ═══ Working with Relationships ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Working with Relationships</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You can create related resources using the <Code>belongs_to</Code> field type. For example,
            a Note that belongs to a Category:
          </p>

          <CodeBlock filename="Terminal">
{`grit generate resource Category --fields "name:string:unique"
grit generate resource Note --fields "title:string,content:richtext,category:belongs_to:Category"`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This creates a <Code>CategoryID</Code> foreign key on the Note model and sets up the GORM relationship.
            The generated UI includes a dropdown to select the category when creating or editing a note.
          </p>

          <Challenge number={7} title="Create Related Resources">
            <p>Generate a Category resource with a <Code>name</Code> field, then generate a Note resource with <Code>title:string,content:text,category:belongs_to:Category</Code>. Create 3 categories and 5 notes. Does the note form show a category dropdown?</p>
          </Challenge>
        </section>

        {/* ═══ Removing Resources ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Removing Resources</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Made a mistake? The <Code>grit remove resource</Code> command cleanly removes everything that was generated:
          </p>

          <CodeBlock filename="Terminal">
{`grit remove resource Task`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This removes the model file, service file, types file, route files, the sidebar entry,
            and the bindings in <Code>app.go</Code>. It{"'"}s a clean undo of <Code>grit generate resource</Code>.
          </p>

          <Challenge number={8} title="Remove and Recreate">
            <p>Remove the Task resource with <Code>grit remove resource Task</Code>. Verify: is the model file gone? Are the route files gone? Is the sidebar entry removed? Now generate it again with different fields: <Code>grit generate resource Task --fields {'"'}title:string,due_date:date,status:string,done:bool{'"'}</Code>.</p>
          </Challenge>
        </section>

        {/* ═══ Database Migrations ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Database Migrations</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            GORM auto-migrates your models when the app starts. When you add a new field to a model,
            restart the app and GORM adds the column to the SQLite database automatically.
          </p>

          <Note>
            SQLite does not support dropping columns or changing column types through GORM auto-migrate.
            If you need to make destructive changes, delete the <Code>.db</Code> file and let GORM recreate it
            from scratch. In development, this is fine — your seed data will repopulate.
          </Note>

          <Challenge number={9} title="Watch Auto-Migration">
            <p>Add a new field to your Task model manually: <Code>{'Notes string `gorm:"type:text" json:"notes"`'}</Code>. Restart the app. Open GORM Studio — is the new column there? Create a task with notes to confirm it works.</p>
          </Challenge>
        </section>

        {/* ═══ Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What You Learned</h2>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How desktop CRUD differs from web CRUD — direct function calls vs HTTP</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How to generate resources with <Code>grit generate resource</Code> for desktop</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How Wails bindings work — Go methods become TypeScript functions</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> GORM models with SQLite — same patterns, minor differences</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> TanStack Router file-based routing for desktop apps</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> The service layer pattern — business logic separated from bindings</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Data type mappings from Grit types to Go, SQLite, and TypeScript</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How to remove resources cleanly with <Code>grit remove resource</Code></li>
          </ul>

          <Challenge number={10} title="Build a Notes App">
            <p>Starting from a fresh desktop project, generate two resources: <Code>Category --fields {'"'}name:string:unique{'"'}</Code> and <Code>Note --fields {'"'}title:string,content:richtext,category:belongs_to:Category,pinned:bool,color:string:optional{'"'}</Code>. Create 3 categories and 10 notes across those categories.</p>
          </Challenge>

          <Challenge number={11} title="Verify in GORM Studio">
            <p>Open GORM Studio with <Code>grit studio</Code>. Browse the notes table. Can you see the <Code>category_id</Code> foreign key? Filter notes by category. Export the data if Studio supports it.</p>
          </Challenge>

          <Challenge number={12} title="Test the Full Cycle">
            <p>Create a note, edit it (change the title and category), then delete it. Verify each operation in GORM Studio. Does the note count on the dashboard update after each operation?</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses/grit-desktop/first-app', label: 'Prev: Your First Desktop App' }}
            next={{ href: '/courses/grit-desktop/custom-ui', label: 'Next: Custom UI & Theming' }}
          />
        </div>
      </main>
    </div>
  )
}
