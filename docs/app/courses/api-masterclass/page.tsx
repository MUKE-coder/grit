import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Build & Deploy a REST API with Go (API-Only Masterclass)',
  description: 'Build a complete REST API from scratch using Grit API-only mode. No frontend — pure Go API with authentication, relationships, role-based access, and deployment.',
}

export default function APIMasterclassCourse() {
  return (
    <div className="min-h-screen bg-[#0b1120]">
      <SiteHeader />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">API-Only Masterclass</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">14 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Build & Deploy a REST API with Go (API-Only Masterclass)
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Build a complete REST API from scratch using <Code>grit new myapi --api</Code>. No frontend,
            no admin panel — pure Go API. You{"'"}ll scaffold, generate resources, test with API docs,
            add authentication and roles, and deploy to production.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: When to Go API-Only ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">When to Go API-Only</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Not every project needs a frontend. Sometimes you just need a clean, fast API that other
            services or apps consume. Grit{"'"}s <Code>--api</Code> flag gives you exactly that — a lean
            Go API with no frontend overhead.
          </p>

          <Definition term="Headless API">
            An API that has no built-in user interface. It only serves JSON data over HTTP. The frontend
            is built separately (or doesn{"'"}t exist at all). The API is the product — everything else
            consumes it. This is sometimes called {'"'}API-first{'"'} or {'"'}backend-only{'"'} architecture.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When does API-only make sense?
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Mobile app backend</strong> — your iOS/Android app needs an API, but the frontend is native code, not React</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Microservice</strong> — one service in a larger architecture that handles a specific domain (payments, notifications, search)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Third-party integrations</strong> — building a public API that other developers consume (like Stripe or Twilio)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Separate frontend team</strong> — the frontend is built by a different team using their own stack (Vue, Svelte, Flutter)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Data pipeline</strong> — an API that ingests, transforms, and serves data with no user-facing UI</li>
          </ul>

          <Tip>
            If you{"'"}re unsure whether to go API-only, ask yourself: {'"'}Will I build the frontend
            with React in this same repo?{'"'} If yes, use the default triple architecture. If no
            (or if the frontend doesn{"'"}t exist yet), go API-only.
          </Tip>

          <Challenge number={1} title="Name 3 API-Only Scenarios">
            <p>Think of 3 real-world projects where API-only is the right choice. For each one, explain
            why a frontend would be built separately (or not at all). Example: a weather data service
            that mobile apps and websites consume.</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: Scaffold the API ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Scaffold the API</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Let{"'"}s build a task management API. One command gives you the full project:
          </p>

          <CodeBlock filename="Terminal">
{`grit new taskapi --api`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>--api</Code> flag tells Grit to skip the frontend entirely. No <Code>apps/web</Code>,
            no <Code>apps/admin</Code>, no <Code>turbo.json</Code>, no <Code>pnpm-workspace.yaml</Code>.
            What you get is a clean Go API project:
          </p>

          <CodeBlock filename="Project Structure">
{`taskapi/
├── docker-compose.yml          # PostgreSQL, Redis, MinIO, Mailhog
├── .env                        # Environment variables
├── .gitignore
├── README.md
└── apps/
    └── api/
        ├── cmd/
        │   └── server/
        │       └── main.go     # Entry point
        ├── internal/
        │   ├── config/         # Environment config
        │   ├── database/       # DB connection + migrations
        │   ├── handler/        # HTTP handlers
        │   ├── middleware/      # Auth, CORS, logging, cache
        │   ├── model/          # GORM models
        │   ├── router/         # Route definitions
        │   └── service/        # Business logic
        ├── go.mod
        └── go.sum`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Compare this to a {'"'}triple{'"'} project (the default): no <Code>apps/web</Code>,
            no <Code>apps/admin</Code>, no <Code>packages/shared</Code>. The API stands alone.
          </p>

          <Note>
            Even in API-only mode, you still get GORM Studio (<Code>/studio</Code>), API docs
            (<Code>/docs</Code>), Sentinel (<Code>/sentinel/ui</Code>), and Pulse (<Code>/pulse/ui</Code>).
            These are part of the Go API, not the frontend.
          </Note>

          <Challenge number={2} title="Scaffold and Explore">
            <p>Run <Code>grit new taskapi --api</Code>. Open the project in your editor. How many
            folders are inside <Code>apps/api/internal/</Code>? List them. Compare the folder count
            to what a <Code>grit new myapp</Code> (triple) project would create.</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Project Structure Tour ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Project Structure Tour</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Let{"'"}s walk through each folder in the API and understand what it does:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Folder</th>
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Purpose</th>
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Key Files</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-mono text-foreground">cmd/server/</td><td className="px-3 py-2">Application entry point</td><td className="px-3 py-2">main.go</td></tr>
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-mono text-foreground">internal/config/</td><td className="px-3 py-2">Environment variables → Go struct</td><td className="px-3 py-2">config.go</td></tr>
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-mono text-foreground">internal/database/</td><td className="px-3 py-2">DB connection, migrations, seeding</td><td className="px-3 py-2">database.go</td></tr>
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-mono text-foreground">internal/model/</td><td className="px-3 py-2">GORM struct definitions</td><td className="px-3 py-2">user.go, upload.go</td></tr>
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-mono text-foreground">internal/handler/</td><td className="px-3 py-2">HTTP request handlers (thin layer)</td><td className="px-3 py-2">auth_handler.go, user_handler.go</td></tr>
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-mono text-foreground">internal/service/</td><td className="px-3 py-2">Business logic (called by handlers)</td><td className="px-3 py-2">auth_service.go, cache_service.go</td></tr>
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-mono text-foreground">internal/middleware/</td><td className="px-3 py-2">Request pipeline (auth, CORS, etc.)</td><td className="px-3 py-2">auth.go, cors.go, logger.go</td></tr>
                <tr><td className="px-3 py-2 font-mono text-foreground">internal/router/</td><td className="px-3 py-2">Route registration</td><td className="px-3 py-2">router.go</td></tr>
              </tbody>
            </table>
          </div>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The pattern is simple: <strong className="text-foreground">handlers</strong> receive HTTP
            requests, validate input, and call <strong className="text-foreground">services</strong>.
            Services contain business logic, call the database through <strong className="text-foreground">models</strong>,
            and return results. Handlers format the results as JSON responses. This separation keeps
            your code testable and organized.
          </p>

          <CodeBlock filename="Request Flow">
{`HTTP Request
  → Router (matches URL to handler)
    → Middleware (auth, logging, CORS)
      → Handler (validates input, calls service)
        → Service (business logic, DB queries)
          → Model (GORM struct, database table)
        ← Service returns result
      ← Handler formats JSON response
    ← Middleware adds headers
  ← Router sends response
HTTP Response`}
          </CodeBlock>

          <Challenge number={3} title="Count the Folders">
            <p>Open <Code>apps/api/internal/</Code>. How many folders are there? List each one and write
            a one-sentence description of what it does. Which folder contains the most files?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Start the API ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Start the API</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Starting an API-only project requires two steps: start the infrastructure services (database,
            Redis, etc.) and then start the Go API itself.
          </p>

          <CodeBlock filename="Terminal">
{`# Step 1: Start infrastructure (PostgreSQL, Redis, MinIO, Mailhog)
cd taskapi
docker compose up -d

# Step 2: Start the Go API
cd apps/api
go run cmd/server/main.go`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You should see output like:
          </p>

          <CodeBlock filename="Output">
{`[GIN-debug] Listening and serving HTTP on :8080
Connected to PostgreSQL
Connected to Redis
GORM Studio available at /studio
API Docs available at /docs`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Verify everything is working:
          </p>

          <CodeBlock filename="Terminal">
{`# Health check
curl localhost:8080/pulse/health

# Expected response:
# {"status":"ok","database":"connected","redis":"connected"}`}
          </CodeBlock>

          <Tip>
            If the API fails to start, check that Docker services are running
            with <Code>docker compose ps</Code>. The most common issue is PostgreSQL not being ready
            yet — wait a few seconds and try again.
          </Tip>

          <Challenge number={4} title="Start and Verify">
            <p>Start the infrastructure and API. Verify the health check
            at <Code>localhost:8080/pulse/health</Code>. Are both database and Redis connected? Open
            GORM Studio at <Code>localhost:8080/studio</Code> — how many tables exist in the fresh
            database?</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Generate Resources ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Generate Resources</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This is where the code generator shines. Instead of manually creating models, handlers,
            services, and routes for each entity, you describe the resource and Grit generates
            everything.
          </p>

          <Definition term="Resource">
            A data entity in your API — something you can Create, Read, Update, and Delete (CRUD).
            Examples: User, Task, Product, Order. Each resource gets a database table, a Go model,
            a service layer, HTTP handlers, and API routes.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Let{"'"}s generate two resources for our task API:
          </p>

          <CodeBlock filename="Terminal">
{`# Generate a Task resource
grit generate resource Task --fields "title:string,description:text,priority:int,done:bool,due_date:date:optional"

# Generate a Category resource
grit generate resource Category --fields "name:string:unique,color:string"`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For each resource, Grit generates:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Model</strong> — <Code>internal/model/task.go</Code> with GORM struct and field tags</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Service</strong> — <Code>internal/service/task_service.go</Code> with CRUD operations + pagination + filtering</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Handler</strong> — <Code>internal/handler/task_handler.go</Code> with HTTP endpoints</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Routes</strong> — injected into <Code>internal/router/router.go</Code></li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Migration</strong> — auto-migrated by GORM on startup</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The field type syntax is <Code>name:type</Code> with optional modifiers:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Field</th>
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Go Type</th>
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">DB Column</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-mono">title:string</td><td className="px-3 py-2">string</td><td className="px-3 py-2">VARCHAR NOT NULL</td></tr>
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-mono">description:text</td><td className="px-3 py-2">string</td><td className="px-3 py-2">TEXT NOT NULL</td></tr>
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-mono">priority:int</td><td className="px-3 py-2">int</td><td className="px-3 py-2">INTEGER NOT NULL</td></tr>
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-mono">done:bool</td><td className="px-3 py-2">bool</td><td className="px-3 py-2">BOOLEAN NOT NULL</td></tr>
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-mono">due_date:date:optional</td><td className="px-3 py-2">*time.Time</td><td className="px-3 py-2">DATE NULL</td></tr>
                <tr><td className="px-3 py-2 font-mono">name:string:unique</td><td className="px-3 py-2">string</td><td className="px-3 py-2">VARCHAR NOT NULL UNIQUE</td></tr>
              </tbody>
            </table>
          </div>

          <Challenge number={5} title="Generate and Inspect">
            <p>Generate both the Task and Category resources. Then run <Code>grit routes</Code> to see
            all registered API endpoints. How many new endpoints were created? List them. Restart the
            API and check GORM Studio — are the new tables there?</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Test with API Docs ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Test with API Docs</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit{"'"}s auto-generated API docs at <Code>/docs</Code> are not just documentation — they{"'"}re
            an interactive testing tool. You can send real requests, see real responses, and explore
            every endpoint without writing a single curl command.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Let{"'"}s create some test data through the API:
          </p>

          <CodeBlock filename="Create a Category">
{`POST /api/categories
Content-Type: application/json

{
  "name": "Work",
  "color": "#6c5ce7"
}`}
          </CodeBlock>

          <CodeBlock filename="Create a Task">
{`POST /api/tasks
Content-Type: application/json

{
  "title": "Finish API course",
  "description": "Complete all 14 challenges in the Grit API Masterclass",
  "priority": 1,
  "done": false,
  "due_date": "2026-04-01"
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Once you have data, test the list endpoint with pagination and filtering:
          </p>

          <CodeBlock filename="List with Pagination">
{`# Page 1, 5 items per page
GET /api/tasks?page=1&page_size=5

# Filter by priority
GET /api/tasks?priority=1

# Sort by due date (ascending)
GET /api/tasks?sort=due_date&order=asc`}
          </CodeBlock>

          <Note>
            Every list endpoint automatically supports pagination (<Code>page</Code>, <Code>page_size</Code>),
            sorting (<Code>sort</Code>, <Code>order</Code>), and filtering by any field. Grit generates
            all of this — you don{"'"}t need to write pagination logic.
          </Note>

          <Challenge number={6} title="Create Test Data">
            <p>Open <Code>/docs</Code> and create 3 categories (Work, Personal, Learning) and 10 tasks
            spread across them. Test pagination with <Code>?page=1&page_size=5</Code> — do you get
            5 results on page 1 and 5 on page 2? Sort by priority descending. Does the highest
            priority appear first?</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Add Relationships ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Add Relationships</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Real APIs have relationships between resources. A Task belongs to a Category. An Order
            has many OrderItems. Grit handles relationships with the <Code>belongs_to</Code> field type.
          </p>

          <Definition term="belongs_to Relationship">
            A database relationship where one record references another. A Task {'"'}belongs to{'"'} a
            Category — meaning the tasks table has a <Code>category_id</Code> column that references the
            categories table. When you query a Task, you can include its Category data in the response.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            To add a relationship, regenerate the Task resource with a <Code>category_id</Code> field:
          </p>

          <CodeBlock filename="Terminal">
{`grit generate resource Task --fields "title:string,description:text,priority:int,done:bool,due_date:date:optional,category_id:belongs_to:Category"`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>category_id:belongs_to:Category</Code> field tells Grit:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> Add a <Code>category_id</Code> column (foreign key) to the tasks table</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Add a <Code>Category</Code> field to the Task struct for preloading</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Include the Category data when fetching Tasks (GORM Preload)</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Now when you create a task, you include the <Code>category_id</Code>:
          </p>

          <CodeBlock filename="Create a Task with Category">
{`POST /api/tasks
Content-Type: application/json

{
  "title": "Review pull request",
  "description": "Review the auth middleware changes",
  "priority": 2,
  "done": false,
  "category_id": 1
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            And when you fetch the task, the category data comes along:
          </p>

          <CodeBlock filename="Response">
{`{
  "data": {
    "id": 1,
    "title": "Review pull request",
    "description": "Review the auth middleware changes",
    "priority": 2,
    "done": false,
    "category_id": 1,
    "category": {
      "id": 1,
      "name": "Work",
      "color": "#6c5ce7"
    }
  }
}`}
          </CodeBlock>

          <Challenge number={7} title="Test Relationships">
            <p>Regenerate the Task resource with the <Code>category_id:belongs_to:Category</Code> field.
            Create 3 categories and 5 tasks, each assigned to a category. Fetch a single task — does the
            response include the full category object? Fetch the task list — does every task show its
            category?</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Authentication ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Authentication</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every Grit API ships with JWT authentication. Users register, log in, receive a token,
            and send that token with every request to access protected endpoints.
          </p>

          <Definition term="JWT (JSON Web Token)">
            A self-contained token that encodes user identity. It{"'"}s a signed string that contains the
            user{"'"}s ID, email, and role. The server validates the signature without checking a database —
            making JWT auth stateless and fast. Tokens expire after a set time (e.g., 24 hours).
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The authentication flow:
          </p>

          <CodeBlock filename="1. Register a User">
{`curl -X POST localhost:8080/api/auth/register \\
  -H "Content-Type: application/json" \\
  -d '{"name":"John","email":"john@test.com","password":"password123"}'

# Response:
# {"data":{"id":1,"name":"John","email":"john@test.com","role":"USER"},"message":"Registration successful"}`}
          </CodeBlock>

          <CodeBlock filename="2. Log In">
{`curl -X POST localhost:8080/api/auth/login \\
  -H "Content-Type: application/json" \\
  -d '{"email":"john@test.com","password":"password123"}'

# Response:
# {"data":{"token":"eyJhbGciOiJIUzI1NiIs...","refresh_token":"...","user":{...}}}`}
          </CodeBlock>

          <CodeBlock filename="3. Use the Token on Protected Endpoints">
{`curl localhost:8080/api/tasks \\
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."

# Without the token, protected endpoints return 401 Unauthorized`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>Authorization</Code> header uses the <Code>Bearer</Code> scheme: the
            word {'"'}Bearer{'"'} followed by a space and then the token. This is the industry standard
            for JWT authentication.
          </p>

          <Tip>
            In the API docs UI at <Code>/docs</Code>, there{"'"}s an {'"'}Authorize{'"'} button where you
            can paste your JWT token. Once set, every request from the docs UI will include the token
            automatically.
          </Tip>

          <Challenge number={8} title="Full Auth Flow">
            <p>Using curl or the API docs:</p>
            <ul className="list-disc list-inside mt-2 space-y-1">
              <li>Register a user with your name and email</li>
              <li>Log in and copy the JWT token</li>
              <li>Call <Code>GET /api/tasks</Code> without a token — what error do you get?</li>
              <li>Call <Code>GET /api/tasks</Code> with the token — does it work?</li>
              <li>Try creating a task with the token — does it associate with your user?</li>
            </ul>
          </Challenge>
        </section>

        {/* ═══ Section 9: Role-Based Access ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Role-Based Access</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Authentication answers {'"'}who are you?{'"'} Authorization answers {'"'}what can you do?{'"'}
            Grit supports role-based access control (RBAC) with two built-in roles: <Code>USER</Code> and <Code>ADMIN</Code>.
          </p>

          <Definition term="RBAC (Role-Based Access Control)">
            A security model where permissions are assigned to roles, and roles are assigned to users.
            Instead of giving each user individual permissions, you assign them a role (like USER or
            ADMIN), and the role determines what they can access.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            By default, all generated resource routes require authentication but allow any role. You
            can restrict routes to admin-only by adding the <Code>--roles</Code> flag:
          </p>

          <CodeBlock filename="Admin-Only Resource">
{`# Generate a resource where write operations are admin-only
grit generate resource AuditLog --fields "action:string,user_id:int,details:text" --roles "ADMIN"`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When a route is restricted to ADMIN, the middleware checks the user{"'"}s role from the JWT
            token. If the user has the USER role, they get a <Code>403 Forbidden</Code> response:
          </p>

          <CodeBlock filename="403 Forbidden Response">
{`{
  "error": {
    "code": "FORBIDDEN",
    "message": "You do not have permission to access this resource"
  }
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You can promote a user to ADMIN through GORM Studio — just edit the user{"'"}s record
            and change the <Code>role</Code> field from {'"'}USER{'"'} to {'"'}ADMIN{'"'}.
          </p>

          <Challenge number={9} title="Test Role Restrictions">
            <p>Register two users: one regular user and one admin (promote via GORM Studio). Generate
            an admin-only resource. Log in as the regular user and try to access the admin endpoint.
            What response do you get? Now log in as the admin — does it work?</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Deploy the API ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Deploy the API</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit includes a one-command deployment tool. It connects to your server via SSH, builds
            your Go binary, sets up a systemd service, configures Caddy for HTTPS, and starts
            everything.
          </p>

          <CodeBlock filename="Deploy Command">
{`grit deploy --host deploy@server.com --domain api.myapp.com`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            What <Code>grit deploy</Code> does behind the scenes:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">1.</span> Connects to your server via SSH</li>
            <li className="flex gap-2"><span className="text-primary">2.</span> Uploads your project files</li>
            <li className="flex gap-2"><span className="text-primary">3.</span> Builds the Go binary on the server (<Code>go build</Code>)</li>
            <li className="flex gap-2"><span className="text-primary">4.</span> Creates a systemd service for auto-restart</li>
            <li className="flex gap-2"><span className="text-primary">5.</span> Configures Caddy reverse proxy with automatic HTTPS</li>
            <li className="flex gap-2"><span className="text-primary">6.</span> Starts PostgreSQL and Redis via Docker Compose</li>
            <li className="flex gap-2"><span className="text-primary">7.</span> Runs database migrations</li>
            <li className="flex gap-2"><span className="text-primary">8.</span> Starts the API</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Your server needs: Ubuntu 22.04+, Docker, Go 1.21+, and Caddy. The deploy command handles
            the rest.
          </p>

          <Note>
            For an API-only project, deployment is especially clean. There{"'"}s no frontend to build,
            no static files to serve. It{"'"}s just a single Go binary behind a reverse proxy. The
            entire deployment takes under 2 minutes.
          </Note>

          <Challenge number={10} title="Write Your Deploy Command">
            <p>Write the <Code>grit deploy</Code> command you would use to deploy your taskapi to a
            server. Include the SSH host and domain. What DNS record would you need to create before
            deploying? (Hint: an A record pointing your domain to your server{"'"}s IP address.)</p>
          </Challenge>
        </section>

        {/* ═══ Section 11: Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s everything you learned in this course:
          </p>

          <ul className="space-y-3 mb-6">
            {[
              'API-only mode (--api) gives you a lean Go API with no frontend overhead',
              'The project structure follows handler → service → model separation',
              'grit generate resource creates model, service, handler, and routes in one command',
              'Field types include string, text, int, bool, date, float, and belongs_to for relationships',
              'Every list endpoint gets automatic pagination, sorting, and filtering',
              'JWT authentication with register, login, and token-based access',
              'Role-based access control with USER and ADMIN roles',
              'API docs at /docs provide interactive testing without curl',
              'GORM Studio at /studio lets you browse and edit database records',
              'grit deploy ships your API to production in under 2 minutes',
            ].map((item, i) => (
              <li key={i} className="flex gap-3 text-muted-foreground">
                <CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" />
                <span>{item}</span>
              </li>
            ))}
          </ul>

          <Challenge number={11} title="Build a Bookmark API (Part 1)">
            <p>Time to build something on your own. Create a new API-only project
            called <Code>bookmarkapi</Code>. Generate these resources:</p>
            <ul className="list-disc list-inside mt-2 space-y-1">
              <li><strong>Category</strong> — name:string, color:string</li>
              <li><strong>Bookmark</strong> — title:string, url:string:unique, description:text:optional, category_id:belongs_to:Category, favorite:bool</li>
            </ul>
          </Challenge>

          <Challenge number={12} title="Build a Bookmark API (Part 2)">
            <p>Start the API and create test data: 5 categories and 20 bookmarks. Test these queries:</p>
            <ul className="list-disc list-inside mt-2 space-y-1">
              <li>List all bookmarks with pagination (page_size=5)</li>
              <li>Filter bookmarks by favorite=true</li>
              <li>Sort bookmarks by title ascending</li>
              <li>Get a single bookmark — does it include the category?</li>
            </ul>
          </Challenge>

          <Challenge number={13} title="Build a Bookmark API (Part 3)">
            <p>Add authentication. Register a user, log in, and use the token to create bookmarks.
            Try accessing bookmarks without a token. Does the API reject unauthenticated requests?</p>
          </Challenge>

          <Challenge number={14} title="Build a Bookmark API (Part 4)">
            <p>Write the deploy command for your bookmark API. What domain would you use? What
            infrastructure does the server need? If you have a server, deploy it for real. Your
            first production Go API is live.</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses', label: 'All Courses' }}
            next={{ href: '/courses', label: 'More Courses' }}
          />
        </div>
      </main>
    </div>
  )
}
