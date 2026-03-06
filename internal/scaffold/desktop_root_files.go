package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeDesktopRootFiles(root string, opts DesktopOptions) error {
	files := map[string]string{
		filepath.Join(root, "wails.json"):     desktopWailsJSON(opts),
		filepath.Join(root, "go.mod"):         desktopGoMod(opts),
		filepath.Join(root, ".gitignore"):     desktopGitignore(),
		filepath.Join(root, ".env"):           desktopEnvFile(opts),
		filepath.Join(root, ".env.example"):   desktopEnvExample(opts),
		filepath.Join(root, "README.md"):      desktopReadme(opts),
		filepath.Join(root, "GRIT_SKILL.md"):  desktopGritSkill(opts),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func desktopWailsJSON(opts DesktopOptions) string {
	return fmt.Sprintf(`{
  "$schema": "https://wails.io/schemas/config.v2.json",
  "name": "%s",
  "outputfilename": "%s",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto",
  "author": {
    "name": ""
  }
}
`, opts.ProjectName, opts.ProjectName)
}

func desktopGoMod(opts DesktopOptions) string {
	return fmt.Sprintf(`module %s

go 1.21

require (
	github.com/wailsapp/wails/v2 v2.9.1
	gorm.io/gorm v1.25.12
	github.com/glebarez/sqlite v1.11.0
	gorm.io/driver/postgres v1.5.11
	golang.org/x/crypto v0.31.0
	github.com/jung-kurt/gofpdf v1.16.2
	github.com/xuri/excelize/v2 v2.8.1
	github.com/joho/godotenv v1.5.1
)
`, opts.ProjectName)
}

func desktopGitignore() string {
	return `# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
vendor/

# Node
node_modules/
frontend/dist/
frontend/src/routeTree.gen.ts

# Wails
build/bin/

# IDE
.idea/
.vscode/
*.swp
*.swo

# Env
.env

# OS
.DS_Store
Thumbs.db
`
}

func desktopEnvFile(opts DesktopOptions) string {
	return fmt.Sprintf(`# Database: sqlite (default) or postgres
DB_DRIVER=sqlite
DB_DSN=%s.db

# For Postgres (uncomment and set):
# DB_DRIVER=postgres
# DB_DSN=host=localhost user=postgres password=postgres dbname=%s port=5432 sslmode=disable

APP_NAME=%s
`, opts.ProjectName, opts.ProjectName, desktopTitleCase(opts.ProjectName))
}

func desktopEnvExample(opts DesktopOptions) string {
	return fmt.Sprintf(`# Database: sqlite (default) or postgres
DB_DRIVER=sqlite
DB_DSN=%s.db

# For Postgres (uncomment and set):
DB_DRIVER=postgres
DB_DSN=host=localhost user=postgres password=postgres dbname=%s port=5432 sslmode=disable

APP_NAME=%s
`, opts.ProjectName, opts.ProjectName, desktopTitleCase(opts.ProjectName))
}

func desktopReadme(opts DesktopOptions) string {
	title := desktopTitleCase(opts.ProjectName)
	r := strings.NewReplacer(
		"{{NAME}}", opts.ProjectName,
		"{{TITLE}}", title,
		"{{BT}}", "```",
		"{{B}}", "`",
	)
	return r.Replace(`# {{TITLE}}

A native desktop application built with **Go + Wails v2 + React + TypeScript + Tailwind CSS**.
Compiles to a single binary (~10-15 MB) for Windows, macOS, and Linux.

---

## Prerequisites

| Tool | Version | Install |
|------|---------|---------|
| Go | 1.21+ | [golang.org/dl](https://golang.org/dl) |
| Node.js | 18+ | [nodejs.org](https://nodejs.org) |
| Wails CLI | v2 | {{B}}go install github.com/wailsapp/wails/v2/cmd/wails@latest{{B}} |

---

## Quick Start

{{BT}}bash
# Start development (hot-reload for Go + React)
wails dev

# Or use the Grit CLI wrapper
grit start
{{BT}}

The app opens in a native window. Go changes trigger a rebuild; React changes hot-reload instantly.

---

## Project Structure

{{BT}}
{{NAME}}/
├── main.go                  # Wails entry point — creates app, binds methods
├── app.go                   # App struct with all bound methods + constructor
├── types.go                 # Input structs for bound methods
├── wails.json               # Wails project configuration
├── go.mod / go.sum          # Go module
├── .env                     # Environment config (DB_DRIVER, DB_DSN, APP_NAME)
│
├── internal/
│   ├── config/config.go     # Reads .env configuration
│   ├── db/db.go             # GORM database setup + AutoMigrate
│   ├── models/              # GORM models (user.go, blog.go, contact.go, ...)
│   └── service/             # Business logic (auth.go, blog.go, contact.go, export.go)
│
├── frontend/
│   ├── src/
│   │   ├── main.tsx          # React entry (TanStack Router + hash history)
│   │   ├── routes/           # File-based routes (TanStack Router)
│   │   │   ├── __root.tsx    # Root route (QueryClientProvider, Toaster)
│   │   │   ├── login.tsx     # Login page
│   │   │   ├── register.tsx  # Register page
│   │   │   ├── _layout.tsx   # Auth guard + sidebar layout
│   │   │   └── _layout/      # Protected pages
│   │   │       ├── index.tsx             # Dashboard
│   │   │       ├── blogs.index.tsx       # Blog list (DataTable)
│   │   │       ├── blogs.new.tsx         # Blog create form
│   │   │       ├── blogs.$id.edit.tsx    # Blog edit form
│   │   │       ├── contacts.index.tsx    # Contact list (DataTable)
│   │   │       ├── contacts.new.tsx      # Contact create form
│   │   │       └── contacts.$id.edit.tsx # Contact edit form
│   │   ├── components/
│   │   │   ├── layout/       # title-bar, sidebar, draggable-window
│   │   │   ├── ui/           # button, input, card, dialog, toast, etc.
│   │   │   ├── data-table.tsx  # Reusable DataTable component
│   │   │   └── form-builder.tsx # Reusable FormBuilder component
│   │   ├── hooks/            # use-auth.tsx, use-theme.ts
│   │   └── lib/              # utils.ts, query-client.ts
│   ├── index.html
│   ├── package.json
│   ├── vite.config.ts
│   └── tailwind.config.ts
│
├── build/
│   ├── appicon.png           # App icon (1024x1024)
│   └── bin/                  # Build output (.exe / .app / binary)
│
└── cmd/studio/
    └── main.go               # GORM Studio server (port 4000)
{{BT}}

---

## How It Works

This is NOT a web app. There is **no HTTP server**.

{{BT}}
React Component
  → Wails TypeScript Binding (auto-generated)
    → Go App Method (on App struct in app.go)
      → GORM Service (internal/service/)
        → SQLite Database (local .db file)
{{BT}}

React calls Go functions **directly** via Wails bindings:
{{BT}}tsx
import { GetBlogs, CreateBlog } from "../../wailsjs/go/main/App";

const result = await GetBlogs(1, 10, "search term");
await CreateBlog({ title: "Hello", content: "World" });
{{BT}}

---

## Adding a New Module (Resource)

Use the Grit CLI to generate a full-stack CRUD resource:

{{BT}}bash
grit generate resource Product --fields "name:string,price:float,stock:int,active:bool"
{{BT}}

This creates **5 files** and makes **10 code injections** in a single atomic operation:

**Files created:**
- {{B}}internal/models/product.go{{B}} — GORM model struct
- {{B}}internal/service/product.go{{B}} — Service with List, GetByID, Create, Update, Delete
- {{B}}frontend/src/routes/_layout/products.index.tsx{{B}} — List page (DataTable + search + pagination + export)
- {{B}}frontend/src/routes/_layout/products.new.tsx{{B}} — Create form
- {{B}}frontend/src/routes/_layout/products.$id.edit.tsx{{B}} — Edit form

**Code injected into:**
- {{B}}db.go{{B}} — AutoMigrate
- {{B}}main.go{{B}} — Service init + App constructor args
- {{B}}app.go{{B}} — Struct field, constructor, 7 bound methods (CRUD + PDF/Excel export)
- {{B}}types.go{{B}} — Input struct
- {{B}}sidebar.tsx{{B}} — Nav icon import + nav link
- {{B}}cmd/studio/main.go{{B}} — GORM Studio model registration

### Supported field types

| Type | Go Type | Form Input | Notes |
|------|---------|------------|-------|
| string | string | Text input | Short text |
| text | string | Textarea | Long text |
| richtext | string | Textarea | Rich content |
| int | int | Number input | Signed integer |
| uint | uint | Number input | Unsigned integer |
| float | float64 | Number input | Decimal values |
| bool | bool | Toggle switch | On/off flags |
| date | time.Time | Date picker | Date only |
| datetime | time.Time | DateTime picker | Date + time |
| slug | string | Auto-generated | Uses {{B}}slug:slug:source=title{{B}} syntax |
| belongs_to | uint | Number input | Foreign key |

### Generate more examples

{{BT}}bash
# With slug field
grit generate resource Article --fields "title:string,slug:slug:source=title,content:richtext,published:bool"

# With relationship
grit generate resource Category --fields "name:string,description:text"
grit generate resource Product --fields "name:string,price:float,category_id:belongs_to"
{{BT}}

### Remove a resource

{{BT}}bash
grit remove resource Product
# Deletes all 5 files + reverses all 10 injections. Markers stay intact.
{{BT}}

---

## Customization

### Window Size and Title Bar

The window dimensions are set in {{B}}main.go{{B}}:

{{BT}}go
// main.go
err := wails.Run(&options.App{
    Title:     cfg.AppName,
    Width:     1280,
    Height:    800,
    MinWidth:  900,
    MinHeight: 600,
    Frameless: true,
    // ...
})
{{BT}}

**Common customizations:**

| Property | Default | Description |
|----------|---------|-------------|
| {{B}}Width{{B}} | 1280 | Initial window width in pixels |
| {{B}}Height{{B}} | 800 | Initial window height in pixels |
| {{B}}MinWidth{{B}} | 900 | Minimum window width |
| {{B}}MinHeight{{B}} | 600 | Minimum window height |
| {{B}}Frameless{{B}} | true | Hides native title bar (uses custom title bar) |
| {{B}}Fullscreen{{B}} | false | Start in fullscreen mode |
| {{B}}StartHidden{{B}} | false | Start the window hidden |

If window controls (minimize, maximize, close) are not visible on your screen, increase the {{B}}Width{{B}} / {{B}}Height{{B}}, or set {{B}}MinWidth{{B}} / {{B}}MinHeight{{B}} to smaller values:

{{BT}}go
Width:     1400,  // wider window
Height:    900,   // taller window
MinWidth:  800,   // allow shrinking smaller
MinHeight: 500,
{{BT}}

To use the native OS title bar instead of the custom frameless one:

{{BT}}go
Frameless: false,  // shows native minimize/maximize/close controls
{{BT}}

### Database

Switch between SQLite and PostgreSQL in {{B}}.env{{B}}:

{{BT}}bash
# SQLite (default — fully offline)
DB_DRIVER=sqlite
DB_DSN={{NAME}}.db

# PostgreSQL (uncomment to use)
# DB_DRIVER=postgres
# DB_DSN=host=localhost user=postgres password=postgres dbname={{NAME}} port=5432 sslmode=disable
{{BT}}

### App Name

Change the display name in {{B}}.env{{B}}:

{{BT}}bash
APP_NAME={{TITLE}}
{{BT}}

---

## Build for Distribution

{{BT}}bash
# Build for current platform
grit compile
# Output: build/bin/{{NAME}}.exe (Windows) or build/bin/{{NAME}} (macOS/Linux)

# Cross-compile
wails build -platform windows/amd64
wails build -platform darwin/arm64
wails build -platform linux/amd64

# Windows installer (requires NSIS)
wails build -nsis
{{BT}}

The binary embeds everything: Go runtime, React frontend, static assets, and SQLite driver. Distribute a single file.

---

## GORM Studio

Browse your SQLite database visually:

{{BT}}bash
grit studio
# Opens browser at http://localhost:4000
{{BT}}

---

## Code Markers — Do Not Delete

These comments in your source code are injection points for {{B}}grit generate{{B}} and {{B}}grit remove{{B}}:

{{BT}}
// grit:models          (db.go)
// grit:service-init    (main.go)
/* grit:app-args */     (main.go)
// grit:imports         (app.go)
// grit:fields          (app.go)
/* grit:constructor-params */  (app.go)
/* grit:constructor-assign */  (app.go)
// grit:methods         (app.go)
// grit:input-types     (types.go)
// grit:studio-models   (cmd/studio/main.go)
// grit:nav-icons       (sidebar.tsx)
// grit:nav             (sidebar.tsx)
{{BT}}

Removing or renaming any marker permanently breaks code generation.

---

## AI Prompt

Copy this prompt to ask an AI assistant to build a simple app with this project:

> I have a Grit Desktop project (Go + Wails v2 + React + TanStack Router + TanStack Query + Tailwind CSS + GORM + SQLite). The project already has auth (login/register), blog and contact CRUD resources, a dashboard, PDF/Excel export, and GORM Studio.
>
> Please help me build a **Task Manager** app. Generate these resources using the Grit CLI:
>
> 1. {{B}}grit generate resource Category --fields "name:string,color:string,description:text"{{B}}
> 2. {{B}}grit generate resource Task --fields "title:string,description:text,priority:string,status:string,due_date:date,completed:bool,category_id:belongs_to"{{B}}
>
> After generating, customize the task list page to show priority as colored badges and add a filter dropdown for status (todo/in-progress/done). Add a custom Go method {{B}}GetTasksByStatus(status string){{B}} on the App struct that filters tasks by status.
>
> Read the GRIT_SKILL.md file in the project root for the full reference on how this project works.

---

## Stack

- **Backend**: Go + GORM + SQLite/PostgreSQL
- **Frontend**: React + TypeScript + Tailwind CSS + TanStack Router + TanStack Query
- **Framework**: Wails v2
- **Dev Tools**: GORM Studio, hot-reload, PDF/Excel export
- **Docs**: [gritframework.com/docs/desktop](https://gritframework.com/docs/desktop)
`)
}

func desktopGritSkill(opts DesktopOptions) string {
	r := strings.NewReplacer(
		"{{NAME}}", opts.ProjectName,
		"{{BT}}", "```",
		"{{B}}", "`",
	)
	return r.Replace(`# GRIT_SKILL.md — AI Reference for Grit Desktop

> **Read this file before generating any code.** Every convention here is enforced by the CLI.

---

## 1. What Is This Project?

A **native desktop application** built with:

- **Backend**: Go (Wails v2 bindings) — direct function calls, no HTTP server
- **Frontend**: Vite + React + TanStack Router (file-based) + TanStack Query + Tailwind CSS
- **Database**: SQLite (default, local) via GORM — fully offline-first
- **Distribution**: Single binary (~10-15 MB) for Windows, macOS, and Linux
- **Code Generation**: {{B}}grit generate resource{{B}} — same CLI, auto-detects desktop

Scaffolded with {{B}}grit new-desktop {{NAME}}{{B}}.

---

## 2. Architecture — No HTTP Server

There is **no REST API**. React calls Go functions directly via Wails bindings.

{{BT}}
React Component
  → Wails TypeScript Binding (auto-generated at frontend/wailsjs/)
    → Go App Method (on App struct in app.go)
      → GORM Service (internal/service/)
        → SQLite Database (local .db file)
{{BT}}

**Never** suggest fetch(), Axios, HTTP endpoints, or REST patterns.

---

## 3. Project Structure

{{BT}}
{{NAME}}/
├── main.go                  # Wails entry point — creates app, binds methods
├── app.go                   # App struct with all bound methods + constructor
├── types.go                 # Input structs for bound methods
├── wails.json               # Wails project configuration
├── internal/
│   ├── config/config.go     # Reads .env configuration
│   ├── db/db.go             # GORM database setup + AutoMigrate
│   ├── models/              # GORM models (user.go, blog.go, contact.go, ...)
│   └── service/             # Business logic (auth.go, blog.go, contact.go, export.go)
├── frontend/
│   ├── src/
│   │   ├── main.tsx          # React entry (TanStack Router + hash history)
│   │   ├── routes/           # File-based routes (TanStack Router)
│   │   │   ├── __root.tsx    # Root route
│   │   │   ├── login.tsx / register.tsx
│   │   │   ├── _layout.tsx   # Auth guard + sidebar layout
│   │   │   └── _layout/      # Protected pages (dashboard, blogs, contacts, ...)
│   │   ├── components/       # layout/ (title-bar, sidebar), ui/ (button, input, ...)
│   │   ├── hooks/            # use-auth.tsx, use-theme.ts
│   │   └── lib/              # utils.ts, query-client.ts
│   └── wailsjs/              # Auto-generated Wails TypeScript bindings
├── build/                    # appicon.png + build output (bin/)
└── cmd/studio/main.go        # GORM Studio (port 4000)
{{BT}}

---

## 4. CLI Commands

| Command | Description |
|---------|-------------|
| {{B}}grit new-desktop <name>{{B}} | Scaffold complete desktop project |
| {{B}}grit generate resource <Name> --fields "..."{{B}} | Generate full-stack CRUD resource (5 files + 10 injections) |
| {{B}}grit remove resource <Name>{{B}} | Remove resource (delete files + reverse injections) |
| {{B}}grit start{{B}} | Start development mode (runs {{B}}wails dev{{B}}) |
| {{B}}grit compile{{B}} | Build native executable (runs {{B}}wails build{{B}}) |
| {{B}}grit studio{{B}} | Open GORM Studio on port 4000 |

**Desktop does NOT have**: {{B}}grit start server{{B}}, {{B}}grit start client{{B}}, {{B}}grit migrate{{B}}, or {{B}}grit seed{{B}}.
GORM AutoMigrate runs on startup automatically.

---

## 5. Resource Generation

{{BT}}bash
grit generate resource <Name> --fields "field1:type,field2:type"
{{BT}}

### Supported Field Types

| Type | Go Type | Form Input |
|------|---------|------------|
| string | string | Text input |
| text | string | Textarea |
| richtext | string | Textarea |
| int | int | Number input |
| uint | uint | Number input |
| float | float64 | Number input |
| bool | bool | Toggle switch |
| date | time.Time | Date picker |
| datetime | time.Time | DateTime picker |
| slug | string (uniqueIndex) | Auto-generated (excluded from form) |
| belongs_to | uint (foreign key) | Number input |

### Slug syntax

{{BT}}bash
grit generate resource Article --fields "title:string,slug:slug:source=title,content:richtext"
{{BT}}

### Files Created Per Resource

1. {{B}}internal/models/<snake>.go{{B}} — GORM model
2. {{B}}internal/service/<snake>.go{{B}} — CRUD service (List, ListAll, GetByID, Create, Update, Delete)
3. {{B}}frontend/src/routes/_layout/<plural>.index.tsx{{B}} — List page (DataTable)
4. {{B}}frontend/src/routes/_layout/<plural>.new.tsx{{B}} — Create form
5. {{B}}frontend/src/routes/_layout/<plural>.$id.edit.tsx{{B}} — Edit form

### 7 Wails Bound Methods Generated Per Resource

{{BT}}go
func (a *App) GetProducts(page, pageSize int, search string) (*service.PaginatedResult, error)
func (a *App) GetProduct(id uint) (*models.Product, error)
func (a *App) CreateProduct(input models.ProductInput) (*models.Product, error)
func (a *App) UpdateProduct(id uint, input models.ProductInput) (*models.Product, error)
func (a *App) DeleteProduct(id uint) error
func (a *App) ExportProductsPDF() ([]byte, error)
func (a *App) ExportProductsExcel() ([]byte, error)
{{BT}}

### Calling from React

{{BT}}tsx
import { GetProducts, CreateProduct } from "../../wailsjs/go/main/App";

const result = await GetProducts(1, 10, "search");
// result = { data: Product[], total: number, page: number, pages: number }

await CreateProduct({ name: "Widget", price: 9.99, stock: 100, active: true });
{{BT}}

---

## 6. Code Markers — NEVER Delete

These comments are injection points for {{B}}grit generate{{B}} and {{B}}grit remove{{B}}. Removing them **permanently breaks** code generation.

| File | Marker | Type | Purpose |
|------|--------|------|---------|
| db.go | {{B}}// grit:models{{B}} | line-before | Add model to AutoMigrate |
| main.go | {{B}}// grit:service-init{{B}} | line-before | Initialize service |
| main.go | {{B}}/* grit:app-args */{{B}} | inline | Pass service to constructor |
| app.go | {{B}}// grit:imports{{B}} | line-before | Import service package |
| app.go | {{B}}// grit:fields{{B}} | line-before | Add service field to App struct |
| app.go | {{B}}/* grit:constructor-params */{{B}} | inline | Add constructor parameter |
| app.go | {{B}}/* grit:constructor-assign */{{B}} | inline | Assign service to field |
| app.go | {{B}}// grit:methods{{B}} | line-before | Add 7 bound methods |
| types.go | {{B}}// grit:input-types{{B}} | line-before | Add Input struct |
| studio/main.go | {{B}}// grit:studio-models{{B}} | line-before | Register in GORM Studio |
| sidebar.tsx | {{B}}// grit:nav-icons{{B}} | line-before | Add icon import |
| sidebar.tsx | {{B}}// grit:nav{{B}} | line-before | Add nav link |

**line-before**: Code inserted on the line BEFORE the marker.
**inline**: Code inserted immediately before the marker on the SAME line.

---

## 7. Adding Custom Methods

Add custom Go methods to the App struct **ABOVE** the {{B}}// grit:methods{{B}} marker:

{{BT}}go
// Custom method — add ABOVE the // grit:methods marker
func (a *App) GetProductBySKU(sku string) (*models.Product, error) {
    return a.productService.GetBySKU(sku)
}
// grit:methods
{{BT}}

Then call from React (restart {{B}}wails dev{{B}} first to regenerate bindings):

{{BT}}tsx
import { GetProductBySKU } from "../../wailsjs/go/main/App";
const product = await GetProductBySKU("SKU-123");
{{BT}}

---

## 8. DataTable Features

Every list page includes: search bar (debounced), paginated table, sortable columns, edit/delete buttons, PDF export, Excel export, and new button.

Customize columns in the list page {{B}}index.tsx{{B}}:

{{BT}}tsx
const columns = [
  { accessorKey: "name", header: "Product Name" },
  { accessorKey: "price", header: "Price",
    cell: ({ row }: any) => "$" + Number(row.getValue("price")).toFixed(2) },
  { accessorKey: "stock", header: "Stock" },
  { accessorKey: "active", header: "Status",
    cell: ({ row }: any) => row.getValue("active") ? "Active" : "Inactive" },
];
{{BT}}

---

## 9. Building for Distribution

{{BT}}bash
grit compile              # Current platform → build/bin/
wails build -platform windows/amd64
wails build -platform darwin/arm64
wails build -nsis         # Windows installer (requires NSIS)
{{BT}}

Binary embeds: Go runtime, React frontend, static assets, SQLite driver. Single file distribution.

---

## 10. Golden Rules

1. **Always use {{B}}grit generate resource{{B}}** for new resources — never create files manually
2. **Never modify grit: markers** — they are permanent injection points
3. **Use {{B}}grit remove resource{{B}} to undo** — never delete files manually
4. **Desktop uses SQLite**, not PostgreSQL — no Docker needed
5. **Methods must be on App struct** to be exposed to React
6. **No REST API** — React calls Go directly via Wails bindings
7. **Restart {{B}}wails dev{{B}}** after adding new Go methods
8. **Add custom methods ABOVE** the {{B}}// grit:methods{{B}} marker
9. **No handlers/ or routes/ directories** — all logic in App struct methods
10. **Generate parent models before child models** (belongs_to)

---

## 11. Common LLM Mistakes

| Wrong | Correct |
|-------|---------|
| {{B}}fetch("/api/products"){{B}} | {{B}}GetProducts(){{B}} from Wails bindings |
| Creating handlers/ or routes/ | Add methods to App struct in app.go |
| Suggesting Docker / PostgreSQL | SQLite works out of the box |
| {{B}}grit start server{{B}} + {{B}}grit start client{{B}} | {{B}}grit start{{B}} or {{B}}wails dev{{B}} |
| Manually creating model/service files | Use {{B}}grit generate resource{{B}} |
| Deleting files to remove a resource | Use {{B}}grit remove resource <Name>{{B}} |
| Adding middleware for auth | Check JWT in service layer |
| Using pnpm dev or turbo dev | Use {{B}}wails dev{{B}} |
| Creating apps/ or packages/ | Single directory — no monorepo |

---

## 12. Quick Reference

{{BT}}bash
# Generate resources
grit generate resource Product --fields "name:string,price:float,stock:int,active:bool"
grit generate resource Category --fields "name:string,description:text"

# Remove a resource
grit remove resource Product

# Open database browser
grit studio

# Build for production
grit compile

# Cross-compile + installer
wails build -platform windows/amd64
wails build -platform darwin/arm64
wails build -nsis
{{BT}}
`)
}

// desktopTitleCase converts a kebab-case or snake_case project name to Title Case.
// For example, "my-cool-app" becomes "My Cool App".
func desktopTitleCase(name string) string {
	// Replace hyphens and underscores with spaces
	s := strings.ReplaceAll(name, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")

	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}
