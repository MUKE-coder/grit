import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export default function SeedersPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Seeders
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Seeders populate your database with initial data &mdash; admin accounts, demo users,
                default categories, or any test data you need during development. Grit scaffolds a
                ready-to-use seed system with a dedicated command and extensible seeder functions.
              </p>
            </div>

            <div className="prose-grit">
              {/* Running Seeders */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Running Seeders
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  After running migrations, seed your database with initial data:
                </p>

                <CodeBlock
                  terminal
                  code="cd apps/api && go run cmd/seed/main.go"
                />

                <p className="text-muted-foreground leading-relaxed mb-4">
                  The seed command first ensures all tables exist (by running migrations), then
                  executes each seeder. Seeders are idempotent &mdash; running them multiple times
                  won&apos;t create duplicate records:
                </p>

                <CodeBlock
                  language="bash"
                  filename="output"
                  code={`Database connected successfully
Running migrations...
All tables are up to date — nothing to migrate.
Seeding database...
Created admin user: admin@example.com / password
Created user: jane@example.com / password
Created user: robert@example.com / password
Created user: emily@example.com / password
Created user: michael@example.com / password
Database seeded successfully.`}
                />
              </div>

              {/* Default Seeders */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Default Seeders
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every Grit project ships with two seeders: an admin account and a set of demo users.
                  These are defined in{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">internal/database/seed.go</code>.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">Admin User</h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Creates the default admin account for accessing the admin panel:
                </p>

                <CodeBlock
                  language="go"
                  filename="apps/api/internal/database/seed.go"
                  code={`func seedAdminUser(db *gorm.DB) error {
    var count int64
    db.Model(&models.User{}).Where("email = ?", "admin@example.com").Count(&count)
    if count > 0 {
        log.Println("Admin user already exists, skipping...")
        return nil
    }

    admin := models.User{
        Name:     "Admin",
        Email:    "admin@example.com",
        Password: "password",
        Role:     "admin",
        Active:   true,
    }

    if err := db.Create(&admin).Error; err != nil {
        return fmt.Errorf("creating admin user: %w", err)
    }

    log.Println("Created admin user: admin@example.com / password")
    return nil
}`}
                />

                <div className="p-4 rounded-lg border border-yellow-500/20 bg-yellow-500/5 mb-6">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-yellow-500/90">Important:</strong> Change the default admin
                    password immediately in production. The password is hashed via bcrypt in the{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">BeforeCreate</code> hook.
                  </p>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">Demo Users</h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Creates sample user accounts with different roles for testing:
                </p>

                <CodeBlock
                  language="go"
                  filename="apps/api/internal/database/seed.go"
                  code={`func seedDemoUsers(db *gorm.DB) error {
    users := []models.User{
        {Name: "Jane Cooper", Email: "jane@example.com", Password: "password", Role: "editor", Active: true},
        {Name: "Robert Fox", Email: "robert@example.com", Password: "password", Role: "user", Active: true},
        {Name: "Emily Davis", Email: "emily@example.com", Password: "password", Role: "user", Active: true},
        {Name: "Michael Chen", Email: "michael@example.com", Password: "password", Role: "user", Active: false},
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
}`}
                />
              </div>

              {/* Seed Entrypoint */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  The Seed Entrypoint
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The seed command lives at{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">cmd/seed/main.go</code>.
                  It runs migrations first to ensure all tables exist, then executes the{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Seed()</code> function:
                </p>

                <CodeBlock
                  language="go"
                  filename="apps/api/cmd/seed/main.go"
                  code={`package main

import (
    "fmt"
    "log"
    "os"

    "myapp/apps/api/internal/config"
    "myapp/apps/api/internal/database"
    "myapp/apps/api/internal/models"
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
}`}
                />
              </div>

              {/* Creating Custom Seeders */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Creating Custom Seeders
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Add your own seeders by creating new functions and registering them in the{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Seed()</code> function.
                  Each seeder should be idempotent &mdash; safe to run multiple times without creating duplicates.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">Step 1: Write the seeder function</h3>

                <CodeBlock
                  language="go"
                  filename="apps/api/internal/database/seed.go"
                  code={`func seedCategories(db *gorm.DB) error {
    categories := []models.Category{
        {Name: "Electronics", Slug: "electronics"},
        {Name: "Clothing", Slug: "clothing"},
        {Name: "Books", Slug: "books"},
        {Name: "Home & Garden", Slug: "home-garden"},
    }

    for _, cat := range categories {
        var count int64
        db.Model(&models.Category{}).Where("slug = ?", cat.Slug).Count(&count)
        if count > 0 {
            continue
        }

        if err := db.Create(&cat).Error; err != nil {
            log.Printf("Warning: failed to create category %s: %v", cat.Name, err)
            continue
        }
        log.Printf("Created category: %s", cat.Name)
    }

    return nil
}`}
                />

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">Step 2: Register it in Seed()</h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Add your seeder call to the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Seed()</code> function.
                  If you used <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code>,
                  the generator automatically adds it via the{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">// grit:seeders</code> marker:
                </p>

                <CodeBlock
                  language="go"
                  filename="apps/api/internal/database/seed.go"
                  code={`func Seed(db *gorm.DB) error {
    if err := seedAdminUser(db); err != nil {
        return fmt.Errorf("seeding admin user: %w", err)
    }

    if err := seedDemoUsers(db); err != nil {
        return fmt.Errorf("seeding demo users: %w", err)
    }

    if err := seedCategories(db); err != nil {
        return fmt.Errorf("seeding categories: %w", err)
    }

    // grit:seeders

    return nil
}`}
                />
              </div>

              {/* Seeder Pattern */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Best Practices
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Follow these patterns when writing seeders:
                </p>

                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-4">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Pattern</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Why</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium">Check before creating</td>
                        <td className="px-4 py-2.5">Prevents duplicate records when running seeders multiple times</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium">Use unique fields for lookups</td>
                        <td className="px-4 py-2.5">Check by email, slug, or other unique identifiers &mdash; not by name</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium">Log what you create</td>
                        <td className="px-4 py-2.5">Makes it easy to verify seeding worked correctly</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium">Wrap errors with context</td>
                        <td className="px-4 py-2.5">Use <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">fmt.Errorf(&quot;seeding X: %w&quot;, err)</code> for clear error messages</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-medium">Seed parent records first</td>
                        <td className="px-4 py-2.5">If Products need Categories, seed categories before products</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Fresh + Seed Workflow */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Reset &amp; Reseed
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  A common development workflow is to drop everything and start fresh with seed data.
                  Combine the fresh migration with seeding:
                </p>

                <CodeBlock
                  terminal
                  code={`go run cmd/migrate/main.go --fresh
go run cmd/seed/main.go`}
                />

                <p className="text-muted-foreground leading-relaxed mb-4">
                  This drops all tables, recreates them from your models, and populates them with seed data.
                  Perfect for resetting your local development environment.
                </p>
              </div>

              {/* Nav */}
              <div className="flex items-center justify-between pt-6 border-t border-border/30">
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/backend/migrations" className="gap-1.5">
                    <ArrowLeft className="h-3.5 w-3.5" />
                    Migrations
                  </Link>
                </Button>
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/admin/overview" className="gap-1.5">
                    Admin Overview
                    <ArrowRight className="h-3.5 w-3.5" />
                  </Link>
                </Button>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
