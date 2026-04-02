# Build Guide: Job Portal (Single + Vite)

Step-by-step guide to building a job portal with Grit's Single architecture
and TanStack Router. One Go binary serves the API and the React frontend.
No Turborepo, no monorepo -- flat project structure like Laravel or Next.js.

---

## Step 1: Scaffold the Project

```bash
grit new job-portal --single --vite
cd job-portal
```

This creates a flat project with `main.go` at the root and `frontend/`
for the React SPA. No `apps/` directory, no `packages/`, no `turbo.json`.

Look at what was generated:

```
job-portal/
├── main.go            # go:embed + Gin server
├── go.mod             # module job-portal
├── internal/          # Go backend code
├── frontend/          # React + TanStack Router + Vite
├── .env
├── docker-compose.yml
└── Makefile
```

Start infrastructure:

```bash
docker compose up -d   # PostgreSQL, Redis, MinIO, Mailhog
```

## Step 2: Understand the go:embed Pattern

Open `main.go`. The key lines:

```go
//go:embed frontend/dist/*
var frontendFS embed.FS
```

This directive tells the Go compiler to embed the contents of `frontend/dist/`
into the binary at compile time. The variable `frontendFS` becomes a virtual
file system containing the built frontend assets.

Later in `main.go`, all API routes are registered under `/api/*`, and a
catch-all route serves the embedded frontend for everything else:

```go
// API routes
routes.Register(r, db, cfg)

// SPA fallback — serve embedded frontend
distFS, _ := fs.Sub(frontendFS, "frontend/dist")
r.NoRoute(gin.WrapH(http.FileServer(http.FS(distFS))))
```

In development, `frontend/dist/` does not exist yet, so the Go server only
handles API requests. You run Vite separately for the frontend.

## Step 3: Development Workflow

You need two terminals during development:

```bash
# Terminal 1: Go API
go run main.go
# Runs on http://localhost:8080
```

```bash
# Terminal 2: Vite frontend
cd frontend && pnpm install && pnpm dev
# Runs on http://localhost:5173
```

Open http://localhost:5173 in your browser. Vite serves the React app with
hot module replacement. The Vite config proxies `/api` requests to Go:

```typescript
// frontend/vite.config.ts
export default defineConfig({
  plugins: [react(), TanStackRouterVite()],
  server: {
    port: 5173,
    proxy: {
      "/api": "http://localhost:8080",
    },
  },
});
```

This means `fetch("/api/jobs")` in the browser goes through Vite to Go
on port 8080. No CORS issues during development.

## Step 4: Generate Resources

```bash
grit generate resource Category --fields "name:string:unique,slug:slug:name"

grit generate resource Company --fields \
  "name:string,slug:slug:name,description:richtext,logo:string,website:string:optional,email:string"

grit generate resource Job --fields \
  "title:string,description:richtext,company_id:belongs_to:Company,location:string,salary_min:float,salary_max:float,type:string,experience:string,active:bool,deadline:date"

grit generate resource Application --fields \
  "job_id:belongs_to:Job,applicant_name:string,applicant_email:string,resume:string,cover_letter:text:optional,status:string"
```

Each command generates:
- Go model in `internal/models/`
- Go service in `internal/services/`
- Go handler in `internal/handlers/`
- Route injection in `internal/routes/routes.go`
- TypeScript types in `frontend/src/lib/types.ts`
- Zod schemas in `frontend/src/lib/schemas.ts`
- React Query hooks in `frontend/src/hooks/`

Note the import paths use `job-portal/internal/...` (not `job-portal/apps/api/internal/...`).

## Step 5: Run Migrations

```bash
grit migrate
grit seed
```

Restart the Go server to pick up the new routes.

## Step 6: Build the Landing Page

Edit `frontend/src/routes/index.tsx`:

```tsx
import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { api } from "../lib/api-client";

export const Route = createFileRoute("/")({
  component: HomePage,
});

function HomePage() {
  const { data: featuredJobs } = useQuery({
    queryKey: ["jobs", "featured"],
    queryFn: () => api.get("/api/jobs?active=true&page_size=6"),
  });

  const { data: categories } = useQuery({
    queryKey: ["categories"],
    queryFn: () => api.get("/api/categories"),
  });

  return (
    <div className="min-h-screen bg-background">
      <section className="py-20 text-center">
        <h1 className="text-5xl font-bold text-foreground">
          Find Your Next Opportunity
        </h1>
        <p className="mt-4 text-text-secondary">
          Browse thousands of jobs from top companies
        </p>
        <JobSearchBar />
      </section>
      <section className="mx-auto max-w-6xl px-4 py-12">
        <h2 className="text-2xl font-semibold text-foreground">Categories</h2>
        <CategoryGrid categories={categories?.data} />
      </section>
      <section className="mx-auto max-w-6xl px-4 py-12">
        <h2 className="text-2xl font-semibold text-foreground">Featured Jobs</h2>
        <JobGrid jobs={featuredJobs?.data} />
      </section>
    </div>
  );
}
```

## Step 7: Build Job Listings and Detail Pages

`frontend/src/routes/jobs/index.tsx`:

```tsx
export const Route = createFileRoute("/jobs/")({
  component: JobListings,
  validateSearch: (search) => ({
    q: (search.q as string) || "",
    location: (search.location as string) || "",
    type: (search.type as string) || "",
    page: Number(search.page) || 1,
  }),
});

function JobListings() {
  const { q, location, type, page } = Route.useSearch();
  const { data, isLoading } = useQuery({
    queryKey: ["jobs", { q, location, type, page }],
    queryFn: () => api.get("/api/jobs", { params: { q, location, type, page } }),
  });
  // Search bar, filters sidebar, job cards grid, pagination
}
```

`frontend/src/routes/jobs/$id.tsx`:

```tsx
export const Route = createFileRoute("/jobs/$id")({
  component: JobDetail,
});

function JobDetail() {
  const { id } = Route.useParams();
  const { data: job } = useQuery({
    queryKey: ["jobs", id],
    queryFn: () => api.get("/api/jobs/" + id),
  });
  // Job title, company info, description, salary range, apply button
}
```

## Step 8: Build Company Profiles

`frontend/src/routes/companies/$slug.tsx`:

```tsx
export const Route = createFileRoute("/companies/$slug")({
  component: CompanyProfile,
});

function CompanyProfile() {
  const { slug } = Route.useParams();
  const { data: company } = useQuery({
    queryKey: ["companies", slug],
    queryFn: () => api.get("/api/companies/slug/" + slug),
  });
  // Logo, description, website link, open positions list
}
```

## Step 9: Build the Application Flow

`frontend/src/routes/apply/$jobId.tsx`:

```tsx
export const Route = createFileRoute("/apply/$jobId")({
  beforeLoad: ({ context }) => {
    if (!context.auth.user) {
      throw redirect({ to: "/auth/login" });
    }
  },
  component: ApplyPage,
});

function ApplyPage() {
  const { jobId } = Route.useParams();
  const mutation = useMutation({
    mutationFn: (formData: FormData) =>
      api.post("/api/applications", formData),
    onSuccess: () => navigate({ to: "/jobs/" + jobId }),
  });
  // Form: name, email, resume upload, cover letter, submit button
}
```

## Step 10: Protect Admin Routes

Create a layout route at `frontend/src/routes/admin.tsx`:

```tsx
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { AdminSidebar } from "../components/admin-sidebar";

export const Route = createFileRoute("/admin")({
  beforeLoad: ({ context }) => {
    if (!context.auth.user) {
      throw redirect({ to: "/auth/login" });
    }
    if (context.auth.user.role !== "ADMIN") {
      throw redirect({ to: "/" });
    }
  },
  component: AdminLayout,
});

function AdminLayout() {
  return (
    <div className="flex min-h-screen bg-background">
      <AdminSidebar />
      <main className="flex-1 overflow-auto p-6">
        <Outlet />
      </main>
    </div>
  );
}
```

All routes under `frontend/src/routes/admin/` inherit this guard.

## Step 11: Build Admin Management Pages

`frontend/src/routes/admin/jobs.tsx`:

```tsx
export const Route = createFileRoute("/admin/jobs")({
  component: AdminJobs,
});

function AdminJobs() {
  const { data } = useQuery({
    queryKey: ["admin", "jobs"],
    queryFn: () => api.get("/api/jobs"),
  });

  return (
    <div>
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-foreground">Manage Jobs</h1>
        <Button onClick={() => setShowForm(true)}>Add Job</Button>
      </div>
      <DataTable
        columns={jobColumns}
        data={data?.data || []}
        onEdit={handleEdit}
        onDelete={handleDelete}
      />
    </div>
  );
}
```

`frontend/src/routes/admin/applications.tsx` -- review with status updates:

```tsx
function AdminApplications() {
  const updateStatus = useMutation({
    mutationFn: ({ id, status }) =>
      api.patch("/api/applications/" + id, { status }),
    onSuccess: () => queryClient.invalidateQueries(["admin", "applications"]),
  });
  // DataTable with status dropdown: Pending, Reviewed, Interview, Accepted, Rejected
}
```

Repeat the pattern for `admin/companies.tsx` and `admin/categories.tsx`.

## Step 12: Root Layout with Conditional Admin Link

`frontend/src/routes/__root.tsx`:

```tsx
function RootLayout() {
  const { user } = useAuth();
  return (
    <div className="min-h-screen bg-background">
      <nav className="flex items-center justify-between border-b border-border px-6 py-3">
        <div className="flex items-center gap-6">
          <Link to="/" className="text-lg font-bold text-foreground">JobPortal</Link>
          <Link to="/jobs" className="text-text-secondary hover:text-foreground">Jobs</Link>
          <Link to="/companies" className="text-text-secondary hover:text-foreground">Companies</Link>
          {user?.role === "ADMIN" && (
            <Link to="/admin" className="text-accent hover:text-accent-hover">Admin</Link>
          )}
        </div>
        <AuthButtons />
      </nav>
      <Outlet />
    </div>
  );
}
```

## Step 13: File Uploads and Emails

File uploads work the same as other architectures -- the API stores files in
MinIO (dev) or S3/R2 (prod):

```tsx
async function uploadFile(file: File): Promise<string> {
  const formData = new FormData();
  formData.append("file", file);
  const res = await api.post("/api/uploads", formData);
  return res.data.url;
}
```

Emails are sent via Resend through background jobs (asynq):
- Welcome email on registration
- Application confirmation on submit
- Status update notification when admin changes application status

## Step 14: Production Build

This is the key advantage of Single architecture. Three commands:

```bash
# 1. Build the frontend into frontend/dist/
cd frontend && pnpm install && pnpm build

# 2. Build the Go binary (embeds frontend/dist/ into the binary)
cd .. && go build -o job-portal main.go

# 3. Done. One binary.
ls -lh job-portal
# job-portal  ~25MB (includes Go binary + all frontend assets)
```

Or use the Makefile:

```bash
make build
# Runs both steps and outputs ./job-portal
```

Test it locally:

```bash
./job-portal
# Open http://localhost:8080 — both API and frontend served from one port
```

## Step 15: Deploy the Single Binary

### Option A: Copy binary to server

```bash
scp job-portal deploy@server.com:/opt/job-portal/
scp .env deploy@server.com:/opt/job-portal/

ssh deploy@server.com
cd /opt/job-portal
./job-portal
```

Create a systemd service to keep it running:

```ini
[Unit]
Description=Job Portal
After=network.target

[Service]
Type=simple
User=deploy
WorkingDirectory=/opt/job-portal
ExecStart=/opt/job-portal/job-portal
EnvironmentFile=/opt/job-portal/.env
Restart=always

[Install]
WantedBy=multi-user.target
```

Use Caddy for automatic TLS:

```
jobs.example.com {
    reverse_proxy localhost:8080
}
```

### Option B: grit deploy

```bash
grit deploy --host deploy@server.com --domain jobs.example.com
```

### Option C: Docker Compose

```bash
docker compose -f docker-compose.prod.yml up -d
```

The Docker setup builds the Go binary with embedded frontend inside the
container. No separate frontend container needed.

---

## Summary

The Single architecture gives you the simplest deployment story:

- **One binary** -- `go build` produces a single executable with embedded frontend
- **One port** -- API and frontend served on the same port
- **Flat structure** -- `main.go` at root, `internal/` for Go, `frontend/` for React
- **Same features** -- auth, CRUD, file uploads, emails, admin pages

The tradeoff: no SSR (it is an SPA), and development requires two terminals.
For most internal tools, dashboards, and smaller apps, this is the right choice.

For apps that need SSR, SEO, or separate frontend deployments, use the
Double (`--double --vite`) or Triple (`--triple --next`) architecture.
