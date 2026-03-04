package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeMigrateSeedFiles(root string, opts Options) error {
	apiRoot := filepath.Join(root, "apps", "api")
	module := opts.ProjectName + "/apps/api"

	files := map[string]string{
		filepath.Join(apiRoot, "cmd", "migrate", "main.go"):          apiMigrateMainGo(),
		filepath.Join(apiRoot, "cmd", "seed", "main.go"):             apiSeedMainGo(),
		filepath.Join(apiRoot, "internal", "database", "seed.go"):    apiSeedGo(),
		filepath.Join(apiRoot, "internal", "database", "migrate.go"): apiMigrateGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

// apiMigrateMainGo returns the cmd/migrate/main.go entrypoint.
func apiMigrateMainGo() string {
	return `package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"{{MODULE}}/internal/config"
	"{{MODULE}}/internal/database"
	"{{MODULE}}/internal/models"
)

func main() {
	fresh := flag.Bool("fresh", false, "Drop all tables before migrating")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if *fresh {
		fmt.Println("Dropping all tables...")
		if err := database.DropAll(db); err != nil {
			log.Fatalf("Failed to drop tables: %v", err)
		}
		fmt.Println("All tables dropped.")
	}

	fmt.Println("Running migrations...")
	if err := models.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migrations completed successfully.")
	os.Exit(0)
}
`
}

// apiSeedMainGo returns the cmd/seed/main.go entrypoint.
func apiSeedMainGo() string {
	return `package main

import (
	"fmt"
	"log"
	"os"

	"{{MODULE}}/internal/config"
	"{{MODULE}}/internal/database"
	"{{MODULE}}/internal/models"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Ensure tables exist before seeding
	fmt.Println("Running migrations...")
	if err := models.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Seeding database...")
	if err := database.Seed(db); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	fmt.Println("Database seeded successfully.")
	os.Exit(0)
}
`
}

// apiMigrateGo returns the database/migrate.go helper with DropAll.
func apiMigrateGo() string {
	return `package database

import (
	"fmt"

	"gorm.io/gorm"
)

// DropAll drops all tables in the database.
// Used by the migrate --fresh command.
func DropAll(db *gorm.DB) error {
	// Get all table names
	var tables []string
	if err := db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables).Error; err != nil {
		return fmt.Errorf("failed to list tables: %w", err)
	}

	if len(tables) == 0 {
		return nil
	}

	// Disable foreign key checks and drop all tables
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %q CASCADE", table)).Error; err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	return nil
}
`
}

// apiSeedGo returns the database/seed.go with the default seeder.
func apiSeedGo() string {
	return `package database

import (
	"fmt"
	"log"
	"time"

	"{{MODULE}}/internal/models"
	"gorm.io/gorm"
)

// Seed populates the database with initial data.
// Add your seeders to this function.
func Seed(db *gorm.DB) error {
	if err := seedAdminUser(db); err != nil {
		return fmt.Errorf("seeding admin user: %w", err)
	}

	if err := seedDemoUsers(db); err != nil {
		return fmt.Errorf("seeding demo users: %w", err)
	}

	if err := seedBlogs(db); err != nil {
		return fmt.Errorf("seeding blogs: %w", err)
	}

	// Seed Grit UI component registry (original 4 + 87 extended)
	models.SeedUIComponents(db)
	models.SeedUIComponentsExtended(db)

	// grit:seeders

	return nil
}

// seedAdminUser creates the default admin account.
func seedAdminUser(db *gorm.DB) error {
	var count int64
	db.Model(&models.User{}).Where("email = ?", "admin@example.com").Count(&count)
	if count > 0 {
		log.Println("Admin user already exists, skipping...")
		return nil
	}

	admin := models.User{
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
		Password:  "password",
		Role:      "ADMIN",
		Active:    true,
	}

	if err := db.Create(&admin).Error; err != nil {
		return fmt.Errorf("creating admin user: %w", err)
	}

	log.Println("Created admin user: admin@example.com / password")
	return nil
}

// seedDemoUsers creates sample user accounts for development.
func seedDemoUsers(db *gorm.DB) error {
	users := []models.User{
		{FirstName: "Jane", LastName: "Cooper", Email: "jane@example.com", Password: "password", Role: "EDITOR", Active: true},
		{FirstName: "Robert", LastName: "Fox", Email: "robert@example.com", Password: "password", Role: "USER", Active: true},
		{FirstName: "Emily", LastName: "Davis", Email: "emily@example.com", Password: "password", Role: "USER", Active: true},
		{FirstName: "Michael", LastName: "Chen", Email: "michael@example.com", Password: "password", Role: "USER", Active: false},
	}

	for _, u := range users {
		var count int64
		db.Model(&models.User{}).Where("email = ?", u.Email).Count(&count)
		if count > 0 {
			continue
		}

		if err := db.Create(&u).Error; err != nil {
			log.Printf("Warning: failed to create user %s: %v", u.Email, err)
			continue
		}
		log.Printf("Created user: %s / password", u.Email)
	}

	return nil
}

// seedBlogs creates sample blog posts for development.
func seedBlogs(db *gorm.DB) error {
	var count int64
	db.Model(&models.Blog{}).Count(&count)
	if count > 0 {
		log.Println("Blogs already seeded, skipping...")
		return nil
	}

	now := time.Now()
	past1 := now.Add(-72 * time.Hour)
	past2 := now.Add(-48 * time.Hour)
	past3 := now.Add(-24 * time.Hour)

	blogs := []models.Blog{
		{
			Title:       "Getting Started with Grit",
			Slug:        "getting-started-with-grit",
			Excerpt:     "Learn how to build full-stack applications with Grit, the Go + React meta-framework.",
			Content:     ` + "`" + `<h2>What is Grit?</h2><p>Grit is a full-stack meta-framework that combines <strong>Go</strong> (Gin + GORM) on the backend with <strong>Next.js</strong> (React + TypeScript) on the frontend. It gives you everything you need to build modern web applications — authentication, admin panel, file uploads, background jobs, and more.</p><h2>Quick Start</h2><p>Getting started is as simple as running a single command:</p><pre><code>grit new my-project</code></pre><p>This scaffolds a complete monorepo with:</p><ul><li>A Go API with JWT auth, RBAC, and CRUD endpoints</li><li>A Next.js web app with login, register, and dashboard</li><li>An admin panel with data tables, forms, and widgets</li><li>Docker setup for PostgreSQL, Redis, and MinIO</li></ul><h2>Why Grit?</h2><p>Most full-stack frameworks force you into a single language. Grit lets you use the <em>best tool for each job</em> — Go for performance-critical backend logic, React for rich interactive UIs, and shared TypeScript types to keep everything in sync.</p><blockquote><p>Build with confidence. Ship with Grit.</p></blockquote>` + "`" + `,
			Image:       "",
			Published:   true,
			PublishedAt: &past1,
		},
		{
			Title:       "Building Modern APIs with Go and Gin",
			Slug:        "building-modern-apis-with-go-and-gin",
			Excerpt:     "A deep dive into structuring Go APIs with the Gin framework, GORM, and clean architecture patterns.",
			Content:     ` + "`" + `<h2>API Architecture</h2><p>A well-structured API separates concerns into distinct layers. In Grit, we use a three-layer architecture:</p><ol><li><strong>Handlers</strong> — Parse requests, validate input, return responses</li><li><strong>Services</strong> — Business logic, data transformation, orchestration</li><li><strong>Models</strong> — Database schema, relationships, hooks</li></ol><h2>Example: Creating a Resource</h2><p>When you run <code>grit generate resource post</code>, Grit creates all three layers automatically. Each handler is thin — it delegates to the service layer:</p><pre><code>func (h *PostHandler) Create(c *gin.Context) {
    var input CreatePostInput
    if err := c.ShouldBindJSON(&amp;input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    post, err := h.Service.Create(&amp;input)
    // ...
}</code></pre><h2>Middleware Stack</h2><p>Grit includes production-ready middleware out of the box:</p><ul><li><strong>Auth</strong> — JWT token validation with automatic refresh</li><li><strong>CORS</strong> — Configurable cross-origin resource sharing</li><li><strong>Logger</strong> — Structured request logging</li><li><strong>Cache</strong> — Redis-backed response caching</li></ul>` + "`" + `,
			Image:       "",
			Published:   true,
			PublishedAt: &past2,
		},
		{
			Title:       "The Power of Monorepo Architecture",
			Slug:        "the-power-of-monorepo-architecture",
			Excerpt:     "Why monorepos work great for full-stack apps, and how Grit uses Turborepo to keep things fast.",
			Content:     ` + "`" + `<h2>Why Monorepo?</h2><p>A monorepo keeps your entire application — backend, frontend, shared types, and tooling — in a single repository. This has several advantages:</p><ul><li><strong>Shared types</strong> — Define your API types once, use them everywhere</li><li><strong>Atomic changes</strong> — Update backend and frontend in one commit</li><li><strong>Simplified CI/CD</strong> — One pipeline to build and deploy everything</li><li><strong>Better developer experience</strong> — No context switching between repos</li></ul><h2>Turborepo Integration</h2><p>Grit uses <strong>Turborepo</strong> for task orchestration. When you run <code>turbo dev</code>, it starts all your apps in parallel with intelligent caching:</p><pre><code>turbo dev
# Starts:
#   → apps/api    (Go API on :8080)
#   → apps/web    (Next.js on :3000)
#   → apps/admin  (Admin panel on :3001)</code></pre><h2>Shared Package</h2><p>The <code>packages/shared</code> directory contains Zod schemas, TypeScript types, and API route constants that are shared between the web app and admin panel. When the API changes, update the shared types and both frontends stay in sync.</p>` + "`" + `,
			Image:       "",
			Published:   true,
			PublishedAt: &past3,
		},
		{
			Title:       "Advanced RBAC Patterns",
			Slug:        "advanced-rbac-patterns",
			Excerpt:     "How to implement role-based access control with custom roles, middleware guards, and UI visibility.",
			Content:     ` + "`" + `<h2>Role-Based Access Control</h2><p>Grit ships with a flexible RBAC system that works across both the API and the admin panel. By default, you get three roles: <strong>ADMIN</strong>, <strong>EDITOR</strong>, and <strong>USER</strong>.</p><h2>Adding Custom Roles</h2><p>Need more roles? Use the CLI:</p><pre><code>grit add role MODERATOR</code></pre><p>This automatically updates 7 locations across your codebase — the Go constants, Zod schemas, TypeScript types, sidebar visibility, and form options.</p><h2>Route Protection</h2><p>Protect API routes with the RequireRole middleware:</p><pre><code>admin.Use(middleware.RequireRole("ADMIN", "EDITOR"))</code></pre><p>On the frontend, the sidebar automatically hides navigation items based on the user's role. Admin and Editor users see all resources, while regular users only see their profile.</p>` + "`" + `,
			Image:       "",
			Published:   false,
			PublishedAt: nil,
		},
	}

	for _, b := range blogs {
		if err := db.Create(&b).Error; err != nil {
			log.Printf("Warning: failed to create blog %q: %v", b.Title, err)
			continue
		}
		status := "draft"
		if b.Published {
			status = "published"
		}
		log.Printf("Created blog: %q (%s)", b.Title, status)
	}

	return nil
}
`
}
