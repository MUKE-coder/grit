import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/infrastructure/database')

export default function DatabasePage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Infrastructure</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Database &amp; Migrations
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit uses GORM as its ORM and PostgreSQL as the primary database.
                Define your models as Go structs and run migrations with a dedicated command.
              </p>
            </div>

            <div className="prose-grit">
              {/* PostgreSQL Setup */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  PostgreSQL Setup
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  PostgreSQL runs via Docker Compose. After starting the containers, your database
                  is ready at <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">localhost:5432</code>.
                </p>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">docker compose up -d postgres</span>
                  </div>
                </div>
              </div>

              {/* Database URL */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Connection String
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The database connection is configured via the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">DATABASE_URL</code> environment
                  variable in your <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.env</code> file. The format follows the standard PostgreSQL connection string:
                </p>
                <CodeBlock language="bash" filename=".env" code={`DATABASE_URL=postgres://grit:grit@localhost:5432/myapp?sslmode=disable`} />
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The URL format breakdown:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Part</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Value</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Description</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">protocol</td>
                        <td className="px-4 py-2.5 font-mono text-xs">postgres://</td>
                        <td className="px-4 py-2.5">PostgreSQL driver</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">user:password</td>
                        <td className="px-4 py-2.5 font-mono text-xs">grit:grit</td>
                        <td className="px-4 py-2.5">Auth credentials (change in production)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">host:port</td>
                        <td className="px-4 py-2.5 font-mono text-xs">localhost:5432</td>
                        <td className="px-4 py-2.5">Database server address</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">database</td>
                        <td className="px-4 py-2.5 font-mono text-xs">myapp</td>
                        <td className="px-4 py-2.5">Database name (matches your project name)</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">sslmode</td>
                        <td className="px-4 py-2.5 font-mono text-xs">disable</td>
                        <td className="px-4 py-2.5">Use &quot;require&quot; in production</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* GORM Connection Code */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  GORM Database Connection
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit generates a database connection module at{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/api/internal/database/database.go</code>.
                  It opens a GORM connection with PostgreSQL and configures connection pooling:
                </p>
                <CodeBlock filename="apps/api/internal/database/database.go" code={`package database

import (
    "fmt"
    "log"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// Connect establishes a database connection using the provided DSN.
func Connect(dsn string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
    }

    // Connection pool settings
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)

    log.Println("Database connected successfully")
    return db, nil
}`} />
              </div>

              {/* Connection Pooling */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Connection Pooling
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  GORM uses Go&apos;s built-in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">database/sql</code> connection
                  pool under the hood. Grit configures these defaults:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Setting</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Default</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Description</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">MaxIdleConns</td>
                        <td className="px-4 py-2.5 font-mono text-xs">10</td>
                        <td className="px-4 py-2.5">Maximum idle connections kept open</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">MaxOpenConns</td>
                        <td className="px-4 py-2.5 font-mono text-xs">100</td>
                        <td className="px-4 py-2.5">Maximum total open connections</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <p className="text-sm text-muted-foreground/60 mt-3">
                  You can also set <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">ConnMaxLifetime</code> and{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">ConnMaxIdleTime</code> for
                  long-running production applications. See the Go{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">database/sql</code> docs for details.
                </p>
              </div>

              {/* Migrations */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Migrations
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit uses a smart migration system that only creates tables which don&apos;t exist yet.
                  Migrations run as a separate command before starting the server:
                </p>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm mb-4">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">cd apps/api &amp;&amp; go run cmd/migrate/main.go</span>
                  </div>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  For full details on how migrations work, fresh migrations, foreign key ordering,
                  and the typical workflow, see the{' '}
                  <Link href="/docs/backend/migrations" className="text-primary hover:underline">
                    Migrations guide
                  </Link>.
                </p>
              </div>

              {/* Defining Models */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Defining Models
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Models live in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/api/internal/models/</code>.
                  Each model is a Go struct with GORM tags that define the database schema. Here is the
                  User model that ships with every Grit project:
                </p>
                <CodeBlock filename="apps/api/internal/models/user.go" code={`package models

import (
    "time"

    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

const (
    RoleAdmin  = "admin"
    RoleEditor = "editor"
    RoleUser   = "user"
)

type User struct {
    ID              uint           \`gorm:"primarykey" json:"id"\`
    Name            string         \`gorm:"size:255;not null" json:"name" binding:"required"\`
    Email           string         \`gorm:"size:255;uniqueIndex;not null" json:"email" binding:"required,email"\`
    Password        string         \`gorm:"size:255;not null" json:"-"\`
    Role            string         \`gorm:"size:20;default:user" json:"role"\`
    Avatar          string         \`gorm:"size:500" json:"avatar"\`
    Active          bool           \`gorm:"default:true" json:"active"\`
    EmailVerifiedAt *time.Time     \`json:"email_verified_at"\`
    CreatedAt       time.Time      \`json:"created_at"\`
    UpdatedAt       time.Time      \`json:"updated_at"\`
    DeletedAt       gorm.DeletedAt \`gorm:"index" json:"-"\`
}

// BeforeCreate hashes the password before saving.
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.Password != "" {
        hashedPassword, err := bcrypt.GenerateFromPassword(
            []byte(u.Password), bcrypt.DefaultCost,
        )
        if err != nil {
            return err
        }
        u.Password = string(hashedPassword)
    }
    return nil
}`} />
              </div>

              {/* Struct Tags */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  GORM Struct Tags
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Common GORM struct tags used in Grit models:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Tag</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Example</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Effect</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      {[
                        { tag: 'primarykey', example: 'gorm:"primarykey"', effect: 'Marks as primary key' },
                        { tag: 'size', example: 'gorm:"size:255"', effect: 'Sets VARCHAR length' },
                        { tag: 'not null', example: 'gorm:"not null"', effect: 'Adds NOT NULL constraint' },
                        { tag: 'uniqueIndex', example: 'gorm:"uniqueIndex"', effect: 'Creates unique index' },
                        { tag: 'index', example: 'gorm:"index"', effect: 'Creates regular index' },
                        { tag: 'default', example: 'gorm:"default:user"', effect: 'Sets default column value' },
                        { tag: 'type', example: 'gorm:"type:text"', effect: 'Sets explicit column type' },
                        { tag: 'foreignKey', example: 'gorm:"foreignKey:UserID"', effect: 'Defines foreign key relationship' },
                      ].map((item) => (
                        <tr key={item.tag} className="border-b border-border/20 last:border-b-0">
                          <td className="px-4 py-2.5 font-mono text-xs text-primary/80">{item.tag}</td>
                          <td className="px-4 py-2.5 font-mono text-xs">{item.example}</td>
                          <td className="px-4 py-2.5">{item.effect}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Common GORM Operations */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Common GORM Operations
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  These are the most common database operations you will use in your Grit handlers and services.
                  GORM provides a fluent, chainable API.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">Create</h3>
                <CodeBlock filename="Create a record" code={`user := models.User{
    Name:  "John Doe",
    Email: "john@example.com",
    Password: "secret123",
}

result := db.Create(&user)
if result.Error != nil {
    return fmt.Errorf("creating user: %w", result.Error)
}
// user.ID is now populated`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">Find (Single Record)</h3>
                <CodeBlock filename="Find by ID or condition" code={`// Find by primary key
var user models.User
db.First(&user, 1) // SELECT * FROM users WHERE id = 1

// Find by condition
db.Where("email = ?", "john@example.com").First(&user)

// Check if record exists
if errors.Is(result.Error, gorm.ErrRecordNotFound) {
    // User not found
}`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">Find (Multiple Records)</h3>
                <CodeBlock filename="Query lists with pagination" code={`var users []models.User

// All records
db.Find(&users)

// With conditions
db.Where("active = ?", true).Find(&users)

// With pagination
db.Offset(0).Limit(20).Order("created_at DESC").Find(&users)

// Count total for pagination
var count int64
db.Model(&models.User{}).Where("active = ?", true).Count(&count)`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">Update</h3>
                <CodeBlock filename="Update records" code={`// Update single field
db.Model(&user).Update("name", "Jane Doe")

// Update multiple fields
db.Model(&user).Updates(models.User{
    Name: "Jane Doe",
    Role: "admin",
})

// Update with map (includes zero-value fields)
db.Model(&user).Updates(map[string]interface{}{
    "active": false,
    "name":   "Jane Doe",
})`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">Delete</h3>
                <CodeBlock filename="Delete records" code={`// Soft delete (sets deleted_at, record still in DB)
db.Delete(&user, 1)

// Hard delete (permanently removes from DB)
db.Unscoped().Delete(&user, 1)

// Delete by condition
db.Where("active = ? AND created_at < ?", false, cutoffDate).Delete(&models.User{})`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">Preload (Relationships)</h3>
                <CodeBlock filename="Eager-load relationships" code={`// Define a Post model with relationship
type Post struct {
    ID        uint           \`gorm:"primarykey" json:"id"\`
    Title     string         \`gorm:"size:255;not null" json:"title"\`
    Body      string         \`gorm:"type:text" json:"body"\`
    UserID    uint           \`json:"user_id"\`
    User      User           \`json:"user"\`
    CreatedAt time.Time      \`json:"created_at"\`
}

// Preload the User relationship
var posts []Post
db.Preload("User").Find(&posts)

// Nested preload
db.Preload("User").Preload("Comments").Find(&posts)`} />
              </div>

              {/* Indexing */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Indexing
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Proper indexing is critical for query performance. GORM creates indexes from struct tags
                  during AutoMigrate:
                </p>
                <CodeBlock filename="Index examples" code={`type Product struct {
    ID       uint   \`gorm:"primarykey"\`
    Name     string \`gorm:"size:255;index"\`           // Regular index
    SKU      string \`gorm:"size:100;uniqueIndex"\`     // Unique index
    Category string \`gorm:"size:100;index:idx_cat_price"\` // Composite index
    Price    float64 \`gorm:"index:idx_cat_price"\`     // Same composite index
    DeletedAt gorm.DeletedAt \`gorm:"index"\`           // Soft delete index
}`} />
                <p className="text-sm text-muted-foreground/60">
                  Add indexes to columns you frequently filter, sort, or join on. The{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">DeletedAt</code> field
                  should always have an index since GORM adds a{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">WHERE deleted_at IS NULL</code> condition
                  to every query on soft-deletable models.
                </p>
              </div>

              {/* SQLite for Testing */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  SQLite for Quick Testing
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  If you want to prototype without Docker or PostgreSQL, GORM supports SQLite as a drop-in
                  replacement. Add the SQLite driver and swap the connection:
                </p>
                <CodeBlock filename="SQLite connection" code={`import (
    "github.com/glebarez/sqlite"
    "gorm.io/gorm"
)

func ConnectSQLite(dbPath string) (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to SQLite: %w", err)
    }
    return db, nil
}

// Usage:
// db, err := ConnectSQLite("test.db")     // file-based
// db, err := ConnectSQLite(":memory:")     // in-memory (tests)`} />
                <div className="p-4 rounded-lg border border-yellow-500/20 bg-yellow-500/5">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-yellow-500/90">Note:</strong> SQLite is great for prototyping and unit tests,
                    but always test with PostgreSQL before deploying. Some PostgreSQL-specific features
                    (like JSONB columns, array types, and certain index types) are not available in SQLite.
                  </p>
                </div>
              </div>

              {/* GORM Studio */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  GORM Studio
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every Grit project includes{' '}
                  <Link href="https://github.com/MUKE-coder/gorm-studio" target="_blank" className="text-primary hover:underline">
                    GORM Studio
                  </Link>{' '}
                  &mdash; a full-featured visual database browser and editor embedded directly into your API at{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">/studio</code>.
                </p>

                {/* Feature grid */}
                <div className="grid gap-3 sm:grid-cols-2 mb-6">
                  {[
                    { title: 'Data Browser', desc: 'Paginated grid with sorting, full-text search, column filtering, and relationship navigation' },
                    { title: 'CRUD Operations', desc: 'Create, edit, and delete records through modal forms. Bulk deletion support' },
                    { title: 'Raw SQL Editor', desc: 'Execute queries with read/write detection and DDL blocking for safety' },
                    { title: 'Schema Export', desc: 'Export schemas as SQL, JSON, YAML, DBML, or PNG/PDF entity-relationship diagrams' },
                    { title: 'Data Export', desc: 'Export data as JSON, CSV (ZIP), or SQL INSERT statements' },
                    { title: 'Data Import', desc: 'Import data from JSON, CSV, SQL, or Excel (.xlsx) files into existing tables' },
                    { title: 'Schema Import', desc: 'Import schemas from SQL, JSON, YAML, or DBML files to create tables' },
                    { title: 'Go Model Generation', desc: 'Generate Go model structs with proper GORM tags from your database schema' },
                  ].map((item) => (
                    <div key={item.title} className="rounded-lg border border-border/30 bg-card/30 p-4">
                      <p className="text-sm font-medium text-foreground/90 mb-1">{item.title}</p>
                      <p className="text-xs text-muted-foreground/70 leading-relaxed">{item.desc}</p>
                    </div>
                  ))}
                </div>

                <p className="text-muted-foreground leading-relaxed mb-4">
                  Enable or disable it in your <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.env</code>:
                </p>
                <CodeBlock language="bash" filename=".env" code={`GORM_STUDIO_ENABLED=true`} />

                <h3 className="text-lg font-semibold mt-6 mb-3">Configuration</h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  GORM Studio is mounted in{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">routes.go</code> with
                  all registered models. When you generate a new resource, the CLI automatically injects
                  the model using the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">/* grit:studio */</code> marker.
                </p>
                <CodeBlock language="go" filename="internal/routes/routes.go" code={`studio.Mount(router, db, []interface{}{
    &models.User{},
    &models.Post{},  // auto-injected by grit generate
    /* grit:studio */
}, studio.Config{
    Prefix:     "/studio",
    ReadOnly:   false,        // set true to disable mutations
    DisableSQL: false,        // set true to hide SQL editor
})`} />

                <h3 className="text-lg font-semibold mt-6 mb-3">Security</h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  GORM Studio includes built-in security: table name validation against registered models,
                  parameterized queries, DDL blocking (DROP, ALTER, TRUNCATE, CREATE), CSV formula injection
                  prevention, and SRI hashes for embedded assets.
                </p>

                <p className="text-sm text-muted-foreground/60 mt-3">
                  Access GORM Studio at{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">http://localhost:8080/studio</code>{' '}
                  when the API is running. Disable it in production by setting{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">GORM_STUDIO_ENABLED=false</code>.
                </p>
              </div>

              {/* Nav */}
              <div className="flex items-center justify-between pt-6 border-t border-border/30">
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/infrastructure/docker" className="gap-1.5">
                    <ArrowLeft className="h-3.5 w-3.5" />
                    Docker Setup
                  </Link>
                </Button>
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/infrastructure/deployment" className="gap-1.5">
                    Deployment
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
