import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/migrations')

export default function MigrationsPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Migrations
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit uses GORM&apos;s migration system with a smart wrapper that only creates tables
                that don&apos;t exist yet. Migrations run as a separate command &mdash; not on server startup &mdash; giving
                you full control over when your database schema changes.
              </p>
            </div>

            <div className="prose-grit">
              {/* Mental model / diagram */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  The lifecycle at a glance
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Migrations create the <em>tables</em>; seeders fill them with{' '}
                  <em>rows</em>. Both are explicit commands you run &mdash; never on
                  server startup &mdash; so you always control when the schema and data
                  change. The usual first-run order is migrate, then seed, then serve.
                </p>

                <CodeBlock
                  language="text"
                  filename="first run"
                  code={`  models.Models()          grit migrate            grit seed
   (the registry)     →     creates missing   →     fills tables
                            tables only              (idempotent)
        │                        │                        │
   &User{}, &Upload{},     ✓ User    (skip)        SeedUsers  → admin + demo
   &Category{},            ✓ Upload  (skip)        SeedBlogs  → sample posts
   &Product{}, …           ✓ Category (created)    SeedCategories / SeedProducts
   // grit:models          ✓ Product  (created)    // grit:seeders
        │                        │                        │
        └──── generate resource adds an entry to BOTH lists automatically ────┘`}
                />

                <p className="text-muted-foreground leading-relaxed mt-4">
                  This page covers the left two boxes &mdash; the model registry and the
                  migrate command. For the third box (filling tables with data, faker,
                  and relationships) see{' '}
                  <Link href="/docs/backend/seeders" className="text-primary hover:underline">Seeders</Link>.
                </p>
              </div>

              {/* Running Migrations */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Running Migrations
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Before starting the API server for the first time (or after adding new models),
                  run the migrate command:
                </p>

                <CodeBlock
                  terminal
                  code="cd apps/api && go run cmd/migrate/main.go"
                />

                <p className="text-muted-foreground leading-relaxed mb-4">
                  The migrate command connects to your database, checks which tables already exist,
                  and only creates the ones that are missing. You&apos;ll see output like:
                </p>

                <CodeBlock
                  language="bash"
                  filename="output"
                  code={`Database connected successfully
Running migrations...
  ✓ models.User — already exists, skipping
  ✓ models.Upload — already exists, skipping
  ✓ models.Category — created
Migrated 1 table(s).
Migrations completed successfully.`}
                />
              </div>

              {/* How It Works */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  How It Works
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The migration system is built on two functions in{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">internal/models/user.go</code>:
                  a <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Models()</code> registry
                  and a <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Migrate()</code> runner.
                </p>

                <CodeBlock
                  language="go"
                  filename="apps/api/internal/models/user.go"
                  code={`// Models returns the ordered list of all models for migration.
// Models with no foreign key dependencies come first.
func Models() []interface{} {
    return []interface{}{
        &User{},
        &Upload{},
        // grit:models
    }
}

// Migrate runs database migrations only for tables that don't exist yet.
func Migrate(db *gorm.DB) error {
    models := Models()
    migrated := 0

    for _, model := range models {
        if db.Migrator().HasTable(model) {
            log.Printf("  ✓ %T — already exists, skipping", model)
            continue
        }

        if err := db.AutoMigrate(model); err != nil {
            return fmt.Errorf("migrating %T: %w", model, err)
        }
        log.Printf("  ✓ %T — created", model)
        migrated++
    }

    if migrated == 0 {
        log.Println("All tables are up to date — nothing to migrate.")
    } else {
        log.Printf("Migrated %d table(s).", migrated)
    }

    return nil
}`}
                />

                <p className="text-muted-foreground leading-relaxed mb-4">
                  For each model, <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Migrate()</code> calls{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">db.Migrator().HasTable()</code> first.
                  If the table already exists, it skips it entirely. Only new tables get created via{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">AutoMigrate</code>.
                </p>
              </div>

              {/* The Migrate Command */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  The Migrate Entrypoint
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The migrate command lives at{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">cmd/migrate/main.go</code>.
                  It loads your config, connects to the database, and runs the migration:
                </p>

                <CodeBlock
                  language="go"
                  filename="apps/api/cmd/migrate/main.go"
                  code={`package main

import (
    "flag"
    "fmt"
    "log"
    "os"

    "myapp/apps/api/internal/config"
    "myapp/apps/api/internal/database"
    "myapp/apps/api/internal/models"
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
}`}
                />
              </div>

              {/* Fresh Migrations */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Fresh Migrations
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  When you need to start from scratch &mdash; during development or testing &mdash; use
                  the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--fresh</code> flag.
                  This drops all tables before re-running migrations:
                </p>

                <CodeBlock
                  terminal
                  code="go run cmd/migrate/main.go --fresh"
                />

                <div className="p-4 rounded-lg border border-red-500/20 bg-red-500/5 mb-4">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-red-500/90">Warning:</strong> The{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--fresh</code> flag
                    permanently deletes all data. Never use it in production.
                  </p>
                </div>

                <p className="text-muted-foreground leading-relaxed mb-4">
                  Fresh migrations are useful when:
                </p>
                <ul className="space-y-2.5 mb-4">
                  {[
                    'You changed column types or removed fields from a model',
                    'You need to reset your local development database',
                    'You want to re-seed with fresh test data',
                  ].map((item) => (
                    <li key={item} className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="text-primary mt-1">&#10003;</span>
                      {item}
                    </li>
                  ))}
                </ul>

                <p className="text-muted-foreground leading-relaxed mb-4">
                  The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">DropAll</code> helper
                  uses raw SQL to drop all public tables:
                </p>

                <CodeBlock
                  language="go"
                  filename="apps/api/internal/database/migrate.go"
                  code={`// DropAll drops all tables in the database.
// Used by the migrate --fresh command.
func DropAll(db *gorm.DB) error {
    var tables []string
    if err := db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables).Error; err != nil {
        return fmt.Errorf("failed to list tables: %w", err)
    }

    if len(tables) == 0 {
        return nil
    }

    for _, table := range tables {
        if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %q CASCADE", table)).Error; err != nil {
            return fmt.Errorf("failed to drop table %s: %w", table, err)
        }
    }

    return nil
}`}
                />
              </div>

              {/* Adding New Models */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Adding New Models
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  When you generate a new resource with{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code>,
                  the model is automatically registered in the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Models()</code> function
                  via the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">// grit:models</code> marker.
                </p>

                <CodeBlock
                  terminal
                  code="grit generate resource Category"
                />

                <p className="text-muted-foreground leading-relaxed mb-4">
                  This adds <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">&amp;Category&#123;&#125;</code> to
                  the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Models()</code> list:
                </p>

                <CodeBlock
                  language="go"
                  filename="apps/api/internal/models/user.go"
                  code={`func Models() []interface{} {
    return []interface{}{
        &User{},
        &Upload{},
        &Category{},
        // grit:models
    }
}`}
                />

                <p className="text-muted-foreground leading-relaxed mb-4">
                  After generating the resource, run migrations to create the new table:
                </p>

                <CodeBlock
                  terminal
                  code="go run cmd/migrate/main.go"
                />

                <p className="text-muted-foreground leading-relaxed mb-4">
                  The output confirms that only the new table was created:
                </p>

                <CodeBlock
                  language="bash"
                  filename="output"
                  code={`Running migrations...
  ✓ models.User — already exists, skipping
  ✓ models.Upload — already exists, skipping
  ✓ models.Category — created
Migrated 1 table(s).`}
                />
              </div>

              {/* Foreign Key Ordering */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Foreign Key Ordering
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  When models have foreign key relationships, the order in{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Models()</code> matters.
                  Parent tables must come before child tables so foreign key constraints can be created.
                </p>

                <CodeBlock
                  language="go"
                  filename="Correct ordering"
                  code={`func Models() []interface{} {
    return []interface{}{
        &User{},     // ← No dependencies (parent)
        &Upload{},   // ← Depends on User (has UserID FK)
        &Category{}, // ← No dependencies
        &Product{},  // ← Depends on Category (has CategoryID FK)
        &Order{},    // ← Depends on User (has UserID FK)
        // grit:models
    }
}`}
                />

                <p className="text-muted-foreground leading-relaxed mb-4">
                  The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code> command
                  always appends new models at the end (before the marker). If a new model depends on
                  another table, make sure the parent model is listed first. You can safely reorder
                  the entries in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Models()</code> &mdash;
                  just keep the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">// grit:models</code> marker
                  as the last line.
                </p>

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-primary/90">Tip:</strong> If you see a
                    &quot;foreign key constraint&quot; error during migration, check that
                    the parent model appears before the child model in{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Models()</code>.
                  </p>
                </div>
              </div>

              {/* Typical Workflow */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Typical Workflow
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Here&apos;s the recommended workflow when starting or extending a Grit project:
                </p>

                <div className="space-y-4 mb-4">
                  {[
                    { step: '1', title: 'Start infrastructure', cmd: 'docker compose up -d' },
                    { step: '2', title: 'Run migrations', cmd: 'cd apps/api && go run cmd/migrate/main.go' },
                    { step: '3', title: 'Seed the database (optional)', cmd: 'go run cmd/seed/main.go' },
                    { step: '4', title: 'Start the server', cmd: 'go run cmd/server/main.go' },
                  ].map((item) => (
                    <div key={item.step} className="flex gap-4 items-start">
                      <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full border border-primary/30 bg-primary/10 text-xs font-medium text-primary">
                        {item.step}
                      </div>
                      <div className="flex-1">
                        <p className="text-sm font-medium text-foreground mb-1">{item.title}</p>
                        <CodeBlock terminal code={item.cmd} className="mb-0" />
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* What AutoMigrate Does */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  What GORM AutoMigrate Does
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Under the hood, Grit&apos;s <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Migrate()</code> function
                  calls GORM&apos;s <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">AutoMigrate</code> for
                  each missing table. AutoMigrate will:
                </p>
                <ul className="space-y-2.5 mb-4">
                  {[
                    'Create the table with columns matching your struct fields',
                    'Add indexes and constraints from struct tags (index, uniqueIndex)',
                    'Create foreign key constraints from relationship fields',
                    'Never delete existing columns or tables (safe by design)',
                    'Never change existing column types automatically',
                  ].map((item) => (
                    <li key={item} className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="text-primary mt-1">&#10003;</span>
                      {item}
                    </li>
                  ))}
                </ul>

                <div className="p-4 rounded-lg border border-yellow-500/20 bg-yellow-500/5">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-yellow-500/90">Note:</strong> AutoMigrate is great for development and simple schemas.
                    For production systems that need column renaming, type changes, or data migrations,
                    consider using a dedicated migration tool like{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">golang-migrate</code> or{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">goose</code> alongside GORM.
                    Use <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--fresh</code> during
                    development if you need to change column types.
                  </p>
                </div>
              </div>

              {/* Nav */}
              <div className="flex items-center justify-between pt-6 border-t border-border/30">
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/backend/response-format" className="gap-1.5">
                    <ArrowLeft className="h-3.5 w-3.5" />
                    API Response Format
                  </Link>
                </Button>
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/backend/seeders" className="gap-1.5">
                    Seeders
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
