import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/models')

export default function ModelsPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend (Go API)</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Models & Database
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit uses GORM as its ORM layer. Models are Go structs that map directly to database tables,
                with struct tags controlling schema behavior, JSON serialization, and request validation.
              </p>
            </div>

            <div className="prose-grit">
              {/* ── GORM Model Pattern ─────────────────────────────── */}
              <h2 id="model-pattern">GORM Model Pattern</h2>
              <p>
                Every model in Grit is a plain Go struct decorated with three categories of struct tags:
              </p>
              <ul>
                <li><strong><code>gorm:"..."</code></strong> -- controls the database schema (column type, indexes, constraints).</li>
                <li><strong><code>json:"..."</code></strong> -- controls how the field is serialized in API responses.</li>
                <li><strong><code>binding:"..."</code></strong> -- controls request validation when Gin binds JSON input.</li>
              </ul>
              <p>
                Here is a minimal model that demonstrates all three:
              </p>
              <CodeBlock language="go" filename="apps/api/internal/models/post.go" code={`package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Post struct {
    ID        string         \`gorm:"primaryKey;size:36" json:"id"\`
    Title     string         \`gorm:"size:255;not null" json:"title" binding:"required,min=3"\`
    Slug      string         \`gorm:"size:255;uniqueIndex;not null" json:"slug"\`
    Body      string         \`gorm:"type:text" json:"body"\`
    Published bool           \`gorm:"default:false" json:"published"\`
    AuthorID  string         \`gorm:"size:36;not null;index" json:"author_id"\`
    Author    User           \`gorm:"foreignKey:AuthorID" json:"author,omitempty"\`
    CreatedAt time.Time      \`json:"created_at"\`
    UpdatedAt time.Time      \`json:"updated_at"\`
    DeletedAt gorm.DeletedAt \`gorm:"index" json:"-"\`
}

// Every Grit model uses a UUID string primary key, generated in a
// BeforeCreate hook. The generator writes this for you on each resource.
func (p *Post) BeforeCreate(tx *gorm.DB) error {
    if p.ID == "" {
        p.ID = uuid.New().String()
    }
    return nil
}`} />

              {/* ── Built-in User Model ─────────────────────────────── */}
              <h2 id="user-model">Built-in User Model</h2>
              <p>
                Every Grit project ships with a <code>User</code> model out of the box. It includes
                authentication fields, role-based access, soft deletes, and a password-hashing hook.
              </p>
              <CodeBlock language="go" filename="apps/api/internal/models/user.go" code={`package models

import (
    "time"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

// Role constants (UPPERCASE — role checks compare against these exactly)
const (
    RoleAdmin  = "ADMIN"
    RoleEditor = "EDITOR"
    RoleUser   = "USER"
    // grit:roles
)

type User struct {
    ID              string         \`gorm:"primaryKey;size:36" json:"id"\`
    FirstName       string         \`gorm:"size:255;not null" json:"first_name" binding:"required"\`
    LastName        string         \`gorm:"size:255;not null" json:"last_name" binding:"required"\`
    Email           string         \`gorm:"size:255;uniqueIndex;not null" json:"email" binding:"required,email"\`
    Password        string         \`gorm:"size:255" json:"-"\`
    Role            string         \`gorm:"size:20;default:USER" json:"role"\`
    Avatar          string         \`gorm:"size:500" json:"avatar"\`
    JobTitle        string         \`gorm:"size:255" json:"job_title"\`
    Bio             string         \`gorm:"type:text" json:"bio"\`
    Active          bool           \`gorm:"default:true" json:"active"\`
    Provider        string         \`gorm:"size:50;default:'local'" json:"provider"\` // local | google | github
    GoogleID        string         \`gorm:"size:255" json:"-"\`
    GithubID        string         \`gorm:"size:255" json:"-"\`
    EmailVerifiedAt *time.Time     \`json:"email_verified_at"\`
    Version         int            \`gorm:"not null;default:1" json:"version"\` // bumped on update for offline sync
    CreatedAt       time.Time      \`json:"created_at"\`
    UpdatedAt       time.Time      \`json:"updated_at"\`
    DeletedAt       gorm.DeletedAt \`gorm:"index" json:"-"\`
}`} />
              <p>Key things to note about the User model:</p>
              <ul>
                <li><code>Password</code> uses <code>json:"-"</code> so it is <strong>never</strong> included in API responses.</li>
                <li><code>DeletedAt</code> also uses <code>json:"-"</code> and enables GORM soft deletes.</li>
                <li><code>Email</code> has a <code>uniqueIndex</code> constraint to prevent duplicate accounts.</li>
                <li>The primary key is a <strong>UUID string</strong> (<code>size:36</code>), generated in <code>BeforeCreate</code> — not an auto-increment integer.</li>
                <li>Three built-in roles are defined as <strong>uppercase</strong> constants: <code>ADMIN</code>, <code>EDITOR</code>, <code>USER</code> (default <code>USER</code>). Role checks compare against these exact strings.</li>
                <li><code>Version</code> increments on every update so offline clients can detect a record has moved on (pairs with the sync engine).</li>
              </ul>

              {/* ── User Model Hooks ─────────────────────────────── */}
              <h3 id="user-hooks">Password Hashing Hook</h3>
              <p>
                The User model uses a GORM <code>BeforeCreate</code> hook to automatically hash the password
                with bcrypt before it is saved to the database. A <code>CheckPassword</code> method
                is provided for login verification.
              </p>
              <CodeBlock language="go" filename="apps/api/internal/models/user.go" code={`// BeforeCreate generates a UUID and hashes the password before saving.
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == "" {
        u.ID = uuid.New().String()
    }
    if u.Version == 0 {
        u.Version = 1
    }
    if u.Password != "" {
        hashed, err := bcrypt.GenerateFromPassword(
            []byte(u.Password), bcrypt.DefaultCost,
        )
        if err != nil {
            return err
        }
        u.Password = string(hashed)
    }
    return nil
}

// BeforeUpdate bumps Version so offline clients can detect a stale edit.
func (u *User) BeforeUpdate(tx *gorm.DB) error {
    u.Version++
    return nil
}

// CheckPassword compares the given password with the stored hash.
func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword(
        []byte(u.Password), []byte(password),
    )
    return err == nil
}`} />

              {/* ── Upload Model ─────────────────────────────── */}
              <h2 id="upload-model">Upload Model</h2>
              <p>
                Grit includes an <code>Upload</code> model for tracking files stored in S3-compatible storage.
                Each upload is associated with a user and records the file metadata.
              </p>
              <CodeBlock language="go" filename="apps/api/internal/models/upload.go" code={`package models

import (
    "time"
    "gorm.io/gorm"
)

type Upload struct {
    ID           string         \`gorm:"primaryKey;size:36" json:"id"\`
    Filename     string         \`gorm:"size:255;not null" json:"filename"\`
    OriginalName string         \`gorm:"size:255;not null" json:"original_name"\`
    MimeType     string         \`gorm:"size:100;not null" json:"mime_type"\`
    Size         int64          \`gorm:"not null" json:"size"\`
    Path         string         \`gorm:"size:500;not null" json:"path"\`
    URL          string         \`gorm:"size:500" json:"url"\`
    ThumbnailURL string         \`gorm:"size:500" json:"thumbnail_url"\`
    UserID       string         \`gorm:"size:36;index;not null" json:"user_id"\`
    User         User           \`gorm:"foreignKey:UserID" json:"-"\`
    CreatedAt    time.Time      \`json:"created_at"\`
    UpdatedAt    time.Time      \`json:"updated_at"\`
    DeletedAt    gorm.DeletedAt \`gorm:"index" json:"-"\`
}`} />

              {/* ── Creating Custom Models ─────────────────────────────── */}
              <h2 id="custom-models">Creating Custom Models</h2>
              <p>
                To create a new model, add a new file in <code>apps/api/internal/models/</code>.
                Follow the naming convention: one model per file, file name in <code>snake_case</code>,
                struct name in <code>PascalCase</code>.
              </p>
              <CodeBlock language="go" filename="apps/api/internal/models/invoice.go" code={`package models

import (
    "time"
    "gorm.io/gorm"
)

type Invoice struct {
    ID         string         \`gorm:"primaryKey;size:36" json:"id"\`
    Number     string         \`gorm:"size:50;uniqueIndex;not null" json:"number" binding:"required"\`
    CustomerID string         \`gorm:"size:36;not null;index" json:"customer_id" binding:"required"\`
    Customer   User           \`gorm:"foreignKey:CustomerID" json:"customer,omitempty"\`
    Amount     float64        \`gorm:"not null" json:"amount" binding:"required,gt=0"\`
    Status     string         \`gorm:"size:20;default:pending" json:"status"\`
    DueDate    time.Time      \`json:"due_date"\`
    Notes      string         \`gorm:"type:text" json:"notes"\`
    CreatedAt  time.Time      \`json:"created_at"\`
    UpdatedAt  time.Time      \`json:"updated_at"\`
    DeletedAt  gorm.DeletedAt \`gorm:"index" json:"-"\`
}`} />
              <p>
                When using <code>grit generate resource Invoice</code>, this file is created automatically
                along with the handler, hooks, Zod schema, and admin page.
              </p>

              {/* ── Field Types & GORM Tags ─────────────────────────────── */}
              <h2 id="field-types">Field Types & GORM Tags</h2>
              <p>
                GORM maps Go types to database column types automatically. Use struct tags
                to fine-tune the schema.
              </p>

              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Go Type</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">PostgreSQL</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Example Tag</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">string</td>
                      <td className="px-4 py-2.5 font-mono text-xs">varchar(256)</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`gorm:"size:255"`}</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">string</td>
                      <td className="px-4 py-2.5 font-mono text-xs">text</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`gorm:"type:text"`}</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">int / int64</td>
                      <td className="px-4 py-2.5 font-mono text-xs">bigint</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`gorm:"not null"`}</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">uint</td>
                      <td className="px-4 py-2.5 font-mono text-xs">bigint unsigned</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`gorm:"primarykey"`}</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">float64</td>
                      <td className="px-4 py-2.5 font-mono text-xs">double precision</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`gorm:"not null"`}</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">bool</td>
                      <td className="px-4 py-2.5 font-mono text-xs">boolean</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`gorm:"default:false"`}</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">time.Time</td>
                      <td className="px-4 py-2.5 font-mono text-xs">timestamp</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`json:"created_at"`}</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">*time.Time</td>
                      <td className="px-4 py-2.5 font-mono text-xs">timestamp (nullable)</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`json:"verified_at"`}</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">gorm.DeletedAt</td>
                      <td className="px-4 py-2.5 font-mono text-xs">timestamp (nullable)</td>
                      <td className="px-4 py-2.5 font-mono text-xs">{`gorm:"index" json:"-"`}</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <h3 id="common-tags">Common GORM Tags</h3>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Tag</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Description</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">primarykey</td>
                      <td className="px-4 py-2.5">Marks the field as the primary key</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">size:255</td>
                      <td className="px-4 py-2.5">Sets varchar length</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">not null</td>
                      <td className="px-4 py-2.5">Adds NOT NULL constraint</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">uniqueIndex</td>
                      <td className="px-4 py-2.5">Creates a unique index on the column</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">index</td>
                      <td className="px-4 py-2.5">Creates a regular index on the column</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">{`default:value`}</td>
                      <td className="px-4 py-2.5">Sets the default value for the column</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">type:text</td>
                      <td className="px-4 py-2.5">Overrides the column type (e.g., text, jsonb)</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">{`foreignKey:FieldName`}</td>
                      <td className="px-4 py-2.5">Specifies the foreign key for a relationship</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              {/* ── Relationships ─────────────────────────────── */}
              <h2 id="relationships">Relationships</h2>
              <p>
                GORM supports all standard relationship types. In Grit, the two most common are
                <code>belongsTo</code> and <code>hasMany</code>.
              </p>

              <h3 id="belongs-to">Belongs To</h3>
              <p>
                A <code>belongsTo</code> relationship means the current model holds the foreign key.
                For example, a Post belongs to a User (the author):
              </p>
              <CodeBlock language="go" filename="belongs_to.go" code={`type Post struct {
    ID       string \`gorm:"primaryKey;size:36" json:"id"\`
    Title    string \`gorm:"size:255;not null" json:"title"\`
    // Foreign key field (UUID string, matching the parent's PK)
    AuthorID string \`gorm:"size:36;not null;index" json:"author_id"\`
    // Relationship -- GORM loads the User when you Preload("Author")
    Author   User \`gorm:"foreignKey:AuthorID" json:"author,omitempty"\`
}

// Usage in a handler or service:
var post Post
db.Preload("Author").First(&post, "id = ?", id)`} />
              <p>
                The <code>json:"author,omitempty"</code> tag means the author object is only included
                in the JSON response when it has been preloaded. This avoids empty nested objects.
              </p>

              <h3 id="has-many">Has Many</h3>
              <p>
                A <code>hasMany</code> relationship means another model holds a foreign key pointing
                back to this model. For example, a User has many Posts:
              </p>
              <CodeBlock language="go" filename="has_many.go" code={`type User struct {
    ID        string \`gorm:"primaryKey;size:36" json:"id"\`
    FirstName string \`gorm:"size:255;not null" json:"first_name"\`
    Email     string \`gorm:"size:255;uniqueIndex;not null" json:"email"\`
    // Has many posts
    Posts     []Post \`gorm:"foreignKey:AuthorID" json:"posts,omitempty"\`
}

// Usage: load a user with all their posts
var user User
db.Preload("Posts").First(&user, "id = ?", id)`} />

              <h3 id="many-to-many">Many to Many</h3>
              <p>
                For many-to-many relationships, GORM creates a join table automatically:
              </p>
              <CodeBlock language="go" filename="many_to_many.go" code={`type Post struct {
    ID    string \`gorm:"primaryKey;size:36" json:"id"\`
    Title string \`gorm:"size:255" json:"title"\`
    Tags  []Tag  \`gorm:"many2many:post_tags;" json:"tags,omitempty"\`
}

type Tag struct {
    ID   string \`gorm:"primaryKey;size:36" json:"id"\`
    Name string \`gorm:"size:100;uniqueIndex" json:"name"\`
}

// GORM auto-creates a "post_tags" join table with post_id and tag_id columns.`} />

              {/* ── Soft Deletes ─────────────────────────────── */}
              <h2 id="soft-deletes">Soft Deletes</h2>
              <p>
                All Grit models include a <code>DeletedAt</code> field of type <code>gorm.DeletedAt</code>.
                When you call <code>db.Delete(&record)</code>, GORM does not actually remove the row.
                Instead, it sets <code>deleted_at</code> to the current timestamp.
              </p>
              <CodeBlock language="go" filename="soft_delete.go" code={`// Add this field to any model for soft deletes:
DeletedAt gorm.DeletedAt \`gorm:"index" json:"-"\`

// Soft-delete a record (sets deleted_at, does NOT remove the row)
db.Delete(&user)

// Normal queries automatically exclude soft-deleted records
db.Find(&users) // only returns records where deleted_at IS NULL

// To include soft-deleted records:
db.Unscoped().Find(&users)

// To permanently delete a record:
db.Unscoped().Delete(&user)`} />
              <p>
                The <code>json:"-"</code> tag on <code>DeletedAt</code> ensures the field is never
                exposed in API responses.
              </p>

              {/* ── Migrations ─────────────────────────────── */}
              <h2 id="automigrate">Migrations &amp; the Model Registry</h2>
              <p>
                Models are registered in a <code>Models()</code> function and migrated by a{' '}
                <code>Migrate()</code> runner &mdash; both in <code>models/user.go</code>. Migrations run
                as a separate <code>grit migrate</code> command, <strong>not</strong> on server startup, so
                you control exactly when the schema changes. The generator appends each new model at the{' '}
                <code>// grit:models</code> marker.
              </p>
              <CodeBlock language="go" filename="apps/api/internal/models/user.go" code={`// Models returns the ordered list of all models for migration.
// Parents (no foreign keys) come first so FK constraints resolve.
func Models() []interface{} {
    return []interface{}{
        &User{},
        &Upload{},
        // grit:models
    }
}

// Migrate runs AutoMigrate for every model, creating missing tables and
// ALTERing existing ones to add new columns — with a created/altered/
// unchanged diff logged so nothing migrates silently.
func Migrate(db *gorm.DB) error {
    for _, model := range Models() {
        if err := db.AutoMigrate(model); err != nil {
            return fmt.Errorf("migrating %T: %w", model, err)
        }
    }
    return nil
}`} />
              <p>
                It&apos;s wired to the migrate entrypoint, run on demand:
              </p>
              <CodeBlock terminal code="grit migrate" />
              <p>
                See <a href="/docs/backend/migrations" className="text-primary hover:underline">Migrations</a> for
                the full runner (column diffing, <code>--fresh</code>) and{' '}
                <a href="/docs/backend/seeders" className="text-primary hover:underline">Seeders</a> for filling
                the tables afterwards.
              </p>

              <h3 id="automigrate-behavior">How AutoMigrate Works</h3>
              <p>AutoMigrate is <strong>safe to run repeatedly</strong>. It will:</p>
              <ul>
                <li>Create tables that do not exist yet.</li>
                <li>Add new columns if you add fields to a struct.</li>
                <li>Create indexes defined in struct tags.</li>
              </ul>
              <p>AutoMigrate will <strong>not</strong>:</p>
              <ul>
                <li>Delete columns that no longer exist in the struct.</li>
                <li>Change a column&apos;s type if you change the Go type.</li>
                <li>Drop tables.</li>
              </ul>
              <p>
                For destructive changes, use <code>grit migrate:fresh</code> in development
                (this drops all tables and re-migrates). In production, write manual SQL migrations.
              </p>

              <CodeBlock terminal code={`grit migrate        # run AutoMigrate
grit migrate:fresh  # drop all + re-migrate (dev only!)`} />

              {/* ── Adding a New Model Step by Step ─────────────────────────────── */}
              <h2 id="adding-models">Adding a New Model (Step by Step)</h2>
              <p>
                The fastest way to add a model is with the code generator. But if you want to do it
                manually, follow these steps:
              </p>
              <ol>
                <li>
                  Create a new file in <code>apps/api/internal/models/</code> (e.g., <code>product.go</code>).
                </li>
                <li>
                  Define your struct with <code>gorm</code>, <code>json</code>, and <code>binding</code> tags.
                </li>
                <li>
                  Add <code>&amp;Product{'{}'}</code> to the <code>AutoMigrate</code> call in <code>user.go</code>.
                </li>
                <li>
                  Restart the server -- GORM will create the table automatically.
                </li>
              </ol>
              <CodeBlock language="go" filename="apps/api/internal/models/product.go" code={`package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Product struct {
    ID          string         \`gorm:"primaryKey;size:36" json:"id"\`
    Name        string         \`gorm:"size:255;not null" json:"name" binding:"required"\`
    Description string         \`gorm:"type:text" json:"description"\`
    Price       float64        \`gorm:"not null" json:"price" binding:"required,gt=0"\`
    SKU         string         \`gorm:"size:100;uniqueIndex" json:"sku" binding:"required"\`
    Stock       int            \`gorm:"default:0" json:"stock"\`
    Active      bool           \`gorm:"default:true" json:"active"\`
    CategoryID  string         \`gorm:"size:36;index" json:"category_id"\`
    Category    Category       \`gorm:"foreignKey:CategoryID" json:"category,omitempty"\`
    CreatedAt   time.Time      \`json:"created_at"\`
    UpdatedAt   time.Time      \`json:"updated_at"\`
    DeletedAt   gorm.DeletedAt \`gorm:"index" json:"-"\`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
    if p.ID == "" {
        p.ID = uuid.New().String()
    }
    return nil
}`} />

              <p>
                Or use the generator and let Grit do the work for you:
              </p>
              <CodeBlock terminal code="grit generate resource Product" />
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 mt-10 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/concepts/styles" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Style Variants
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/backend/handlers" className="gap-1.5">
                  Handlers
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
