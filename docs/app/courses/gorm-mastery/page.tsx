import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Database Mastery with GORM: Models, Migrations & Queries',
  description: 'Master GORM — Go\'s most popular ORM. Learn models, migrations, CRUD operations, relationships, scopes, hooks, raw SQL, and performance optimization.',
}

export default function GORMMasteryCourse() {
  return (
    <div className="min-h-screen bg-[#0b1120]">
      <SiteHeader />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Database Mastery with GORM</span>
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
            Database Mastery with GORM: Models, Migrations & Queries
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Every application stores data. The question is how elegantly you interact with it. GORM lets
            you work with databases using Go structs instead of raw SQL. In this course, you{"'"}ll master
            models, migrations, relationships, advanced queries, scopes, hooks, and performance optimization —
            everything you need to build data-heavy applications with Grit.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What is GORM? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is GORM?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Databases speak SQL. Go speaks structs. An ORM bridges the gap — it translates between your
            Go code and the database, so you can create, read, update, and delete records using familiar
            Go syntax instead of writing SQL strings.
          </p>

          <Definition term="ORM (Object-Relational Mapping)">
            A programming technique that maps database tables to objects (or structs) in your code. Instead
            of writing <Code>SELECT * FROM users WHERE id = 1</Code>, you write <Code>db.First(&user, 1)</Code>.
            The ORM generates the SQL for you. This makes code more readable, portable across databases,
            and less prone to SQL injection.
          </Definition>

          <Definition term="GORM">
            Go{"'"}s most popular ORM library. It supports PostgreSQL, MySQL, SQLite, and SQL Server. Features
            include auto-migration, associations (relationships), hooks (callbacks), preloading, scopes,
            soft deletes, and a powerful query builder. Grit uses GORM for all database operations.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s the same operation in raw SQL vs GORM:
          </p>

          <CodeBlock filename="Raw SQL vs GORM">
{`// Raw SQL
db.Exec("INSERT INTO users (name, email, role) VALUES (?, ?, ?)", "John", "john@test.com", "USER")

// GORM — same result, but type-safe and readable
user := User{Name: "John", Email: "john@test.com", Role: "USER"}
db.Create(&user)
// user.ID is now set automatically`}
          </CodeBlock>

          <Challenge number={1} title="ORM vs Raw SQL">
            <p>What are 3 advantages of using an ORM like GORM over writing raw SQL? Think about
            readability, safety, and portability. Can you think of a situation where raw SQL would
            be better than an ORM?</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: Defining Models ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Defining Models</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            A GORM model is a Go struct that represents a database table. Each field becomes a column.
            You control column behavior with struct tags — constraints, indexes, defaults, and JSON
            serialization names.
          </p>

          <CodeBlock filename="models/product.go">
{`type Product struct {
    gorm.Model
    // gorm.Model adds: ID, CreatedAt, UpdatedAt, DeletedAt

    Name        string   // tag: not null; json name
    Slug        string   // tag: uniqueIndex, not null; json slug
    Price       float64  // tag: not null, default:0; json price
    Description string   // tag: type:text; json description
    Active      bool     // tag: default:true; json active
    CategoryID  uint     // tag: index; json category_id

    // Relationship
    Category    Category // GORM auto-loads via CategoryID
}`}
          </CodeBlock>

          <Note>
            GORM struct tags use the format <Code>gorm:{'"'}tagvalue{'"'}</Code>. Common tags:
            <Code>not null</Code> (required), <Code>uniqueIndex</Code> (no duplicates),
            <Code>default:value</Code> (default value), <Code>type:text</Code> (column type),
            <Code>index</Code> (database index for faster queries).
          </Note>

          <Definition term="gorm.Model">
            A base struct that embeds 4 fields into your model: <Code>ID</Code> (uint, auto-incrementing
            primary key), <Code>CreatedAt</Code> (timestamp, set on create), <Code>UpdatedAt</Code> (timestamp,
            updated on every save), and <Code>DeletedAt</Code> (nullable timestamp for soft deletes). Every
            Grit model embeds this.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Common field types and their database column types:
          </p>

          <CodeBlock filename="Field Type Mapping">
{`string    → VARCHAR(255)     // Short text (names, emails)
string    → TEXT              // Long text (with gorm type:text tag)
int       → INTEGER           // Whole numbers
uint      → INTEGER UNSIGNED  // Positive whole numbers (IDs)
float64   → DOUBLE PRECISION  // Decimal numbers (prices)
bool      → BOOLEAN           // true/false
time.Time → TIMESTAMP         // Dates and times`}
          </CodeBlock>

          <Challenge number={2} title="Define a Blog Post Model">
            <p>Define a GORM model for a Blog Post with these fields: Title (required string),
            Slug (unique string), Content (long text), Published (boolean, defaults to false),
            AuthorID (foreign key to User). Don{"'"}t forget to embed <Code>gorm.Model</Code>.</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Auto-Migration ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Auto-Migration</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Once you define a model, you need to create the corresponding database table. GORM{"'"}s
            auto-migration reads your struct and creates (or updates) the table to match. No manual
            SQL needed.
          </p>

          <CodeBlock filename="database/database.go">
{`func Migrate(db *gorm.DB) {
    db.AutoMigrate(
        &models.User{},
        &models.Product{},
        &models.Category{},
        &models.Blog{},
    )
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            In a Grit project, you run migrations with:
          </p>

          <CodeBlock filename="Terminal">
{`# Run migrations (creates/updates tables)
grit migrate

# Fresh migration (drops all tables, re-creates from scratch)
grit migrate --fresh`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            What auto-migration does:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> Creates tables that don{"'"}t exist</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Adds columns that are missing</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Creates indexes and constraints</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Does <strong className="text-foreground">NOT</strong> delete columns (to prevent data loss)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Does <strong className="text-foreground">NOT</strong> change column types</li>
          </ul>

          <Tip>
            Use <Code>grit migrate --fresh</Code> during development when you change column types or
            remove fields. In production, never use <Code>--fresh</Code> — it drops all your data.
            For production column changes, use manual SQL migrations.
          </Tip>

          <Challenge number={3} title="Add a Field and Migrate">
            <p>Add a new <Code>Stock int</Code> field to your Product model. Run <Code>grit migrate</Code>.
            Open GORM Studio and check the products table. Did the new column appear? What{"'"}s its
            default value?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: CRUD Operations ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">CRUD Operations</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            CRUD stands for Create, Read, Update, Delete — the four fundamental database operations.
            GORM provides clean methods for each one. Let{"'"}s see them all:
          </p>

          <CodeBlock filename="Create">
{`// Create a single record
product := models.Product{
    Name:  "Widget",
    Price: 29.99,
    Active: true,
}
result := db.Create(&product)
// product.ID is now set (auto-increment)
// result.Error contains any error
// result.RowsAffected is 1

// Create multiple records
products := []models.Product{
    {Name: "Widget A", Price: 10.0},
    {Name: "Widget B", Price: 20.0},
    {Name: "Widget C", Price: 30.0},
}
db.Create(&products)`}
          </CodeBlock>

          <CodeBlock filename="Read">
{`// Find by primary key
var product models.Product
db.First(&product, 1)          // SELECT * FROM products WHERE id = 1

// Find by condition
db.Where("slug = ?", "widget").First(&product)

// Find all (with pagination)
var products []models.Product
db.Offset(0).Limit(20).Find(&products)

// Find all matching a condition
db.Where("active = ? AND price > ?", true, 10.0).Find(&products)`}
          </CodeBlock>

          <CodeBlock filename="Update">
{`// Update a single field
db.Model(&product).Update("price", 39.99)

// Update multiple fields
db.Model(&product).Updates(map[string]interface{}{
    "name":  "Super Widget",
    "price": 39.99,
})

// Update with struct (only non-zero fields)
db.Model(&product).Updates(models.Product{Name: "Super Widget", Price: 39.99})`}
          </CodeBlock>

          <CodeBlock filename="Delete">
{`// Soft delete (sets DeletedAt timestamp)
db.Delete(&product, id)
// Row still exists in database, but hidden from queries

// Permanently delete (skip soft delete)
db.Unscoped().Delete(&product, id)

// Restore a soft-deleted record
db.Model(&product).Unscoped().Update("deleted_at", nil)`}
          </CodeBlock>

          <Definition term="Soft Delete">
            Instead of removing the row from the database, GORM sets the <Code>DeletedAt</Code> timestamp.
            All subsequent queries automatically exclude soft-deleted records (by adding <Code>WHERE
            deleted_at IS NULL</Code>). The data is still in the database and can be recovered. This is
            the default behavior when your model embeds <Code>gorm.Model</Code>.
          </Definition>

          <Challenge number={4} title="CRUD with GORM Studio">
            <p>Using GORM Studio, perform these operations: (1) Create 3 products with different names
            and prices, (2) Query all products, (3) Update one product{"'"}s price, (4) Soft-delete a
            product, (5) Query all products again — is the deleted one still visible? (6) Use
            <Code>db.Unscoped().Find(&products)</Code> to see all products including deleted ones.</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Relationships ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Relationships</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Real applications have related data. A product belongs to a category. A category has many
            products. A post has many tags, and a tag has many posts. GORM handles all three relationship
            types: belongs_to, has_many, and many_to_many.
          </p>

          <CodeBlock filename="Belongs To (Product → Category)">
{`type Product struct {
    gorm.Model
    Name       string
    Price      float64
    CategoryID uint       // Foreign key
    Category   Category   // GORM loads this automatically
}

type Category struct {
    gorm.Model
    Name string
}

// Query: find a product with its category
var product models.Product
db.Preload("Category").First(&product, 1)
// product.Category.Name is now available`}
          </CodeBlock>

          <CodeBlock filename="Has Many (Category → Products)">
{`type Category struct {
    gorm.Model
    Name     string
    Products []Product  // GORM loads all products for this category
}

// Query: find a category with all its products
var category models.Category
db.Preload("Products").First(&category, 1)
// category.Products is now a slice of Product structs`}
          </CodeBlock>

          <CodeBlock filename="Many to Many (Post ↔ Tag)">
{`type Post struct {
    gorm.Model
    Title string
    Tags  []Tag  // gorm many2many:post_tags
}

type Tag struct {
    gorm.Model
    Name  string
    Posts []Post  // gorm many2many:post_tags
}

// GORM auto-creates a "post_tags" join table
// Query: find a post with its tags
var post models.Post
db.Preload("Tags").First(&post, 1)`}
          </CodeBlock>

          <Definition term="Preloading">
            Loading related data in a single query instead of making separate queries. Without preloading,
            accessing <Code>product.Category</Code> would be empty. With <Code>db.Preload({'"'}Category{'"'}).Find(&products)</Code>,
            GORM runs a second query to load all categories for the found products and attaches them
            to the struct. This solves the N+1 query problem.
          </Definition>

          <Note>
            The many_to_many relationship requires a GORM struct tag on the field. GORM auto-creates
            the join table (in this case <Code>post_tags</Code> with columns <Code>post_id</Code> and
            <Code>tag_id</Code>). You never need to define the join table manually.
          </Note>

          <Challenge number={5} title="Build and Query Relationships">
            <p>Create a Category model with a Name field. Add a <Code>CategoryID</Code> and
            <Code>Category</Code> field to your Product model. Run <Code>grit migrate</Code>. Then:
            (1) Create a category called {'"'}Electronics{'"'}, (2) Create 5 products in that category,
            (3) Query products with <Code>db.Preload({'"'}Category{'"'}).Find(&products)</Code>.
            Does the category name appear on each product?</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Query Building ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Query Building</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            GORM{"'"}s query builder lets you chain conditions to build complex queries. It{"'"}s like
            writing SQL, but in Go — type-safe, composable, and readable.
          </p>

          <CodeBlock filename="Where, Order, Limit">
{`// Find active products over $10, ordered by price (highest first)
var products []models.Product
db.Where("price > ? AND active = ?", 10.0, true).
    Order("price DESC").
    Limit(10).
    Find(&products)

// Multiple Where conditions (AND)
db.Where("category_id = ?", 1).
    Where("price BETWEEN ? AND ?", 10.0, 50.0).
    Find(&products)

// OR condition
db.Where("name LIKE ?", "%widget%").
    Or("name LIKE ?", "%gadget%").
    Find(&products)`}
          </CodeBlock>

          <CodeBlock filename="Count, Select, Group">
{`// Count records
var count int64
db.Model(&models.Product{}).
    Where("active = ?", true).
    Count(&count)

// Select specific columns (faster — less data transferred)
db.Select("name", "price").Find(&products)

// Group by with aggregate
type CategoryStats struct {
    CategoryID uint
    Total      int64
    AvgPrice   float64
}
var stats []CategoryStats
db.Model(&models.Product{}).
    Select("category_id, COUNT(*) as total, AVG(price) as avg_price").
    Group("category_id").
    Scan(&stats)`}
          </CodeBlock>

          <CodeBlock filename="Not, Distinct, Pluck">
{`// Not — exclude records
db.Not("active = ?", false).Find(&products)

// Distinct — unique values only
var names []string
db.Model(&models.Product{}).Distinct("name").Pluck("name", &names)

// Pluck — extract a single column into a slice
var prices []float64
db.Model(&models.Product{}).Where("active = ?", true).Pluck("price", &prices)`}
          </CodeBlock>

          <Challenge number={6} title="Write Complex Queries">
            <p>Write these queries using GORM{"'"}s query builder: (1) Find the 5 most expensive active
            products, ordered by price descending, (2) Count the total number of products per category,
            (3) Find all products whose name contains {'"'}Pro{'"'} and cost more than $50,
            (4) Get the average price of all active products.</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Scopes ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Scopes</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When you find yourself writing the same <Code>Where</Code> clause in multiple places, it{"'"}s
            time to create a scope. A scope is a reusable query function that you can chain onto any query.
          </p>

          <Definition term="Scope">
            A reusable query condition packaged as a function. The function takes a <Code>*gorm.DB</Code> and
            returns a <Code>*gorm.DB</Code> with additional conditions applied. Scopes can be composed —
            you can chain multiple scopes together to build complex queries from simple, named pieces.
          </Definition>

          <CodeBlock filename="scopes.go">
{`// Define reusable scopes
func Active(db *gorm.DB) *gorm.DB {
    return db.Where("active = ?", true)
}

func Expensive(db *gorm.DB) *gorm.DB {
    return db.Where("price > ?", 100)
}

func InCategory(categoryID uint) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("category_id = ?", categoryID)
    }
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        offset := (page - 1) * pageSize
        return db.Offset(offset).Limit(pageSize)
    }
}

// Usage: compose scopes together
var products []models.Product
db.Scopes(Active, Expensive).Find(&products)

// With parameterized scopes
db.Scopes(Active, InCategory(3), Paginate(1, 20)).Find(&products)`}
          </CodeBlock>

          <Tip>
            Scopes make your code self-documenting. Compare <Code>db.Where({'"'}active = ?{'"'}, true).Where({'"'}price
            &gt; ?{'"'}, 100)</Code> with <Code>db.Scopes(Active, Expensive)</Code>. The second version tells
            you exactly what the query does without reading the conditions.
          </Tip>

          <Challenge number={7} title="Create a Published Scope">
            <p>Create a <Code>Published</Code> scope for blog posts that filters to only published posts
            (<Code>published = true</Code>). Create a <Code>Recent</Code> scope that orders by
            <Code>created_at DESC</Code> and limits to 10. Use them together:
            <Code>db.Scopes(Published, Recent).Find(&posts)</Code>. Does it return the 10 most
            recent published posts?</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Hooks ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Hooks</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Sometimes you need code to run automatically before or after a database operation. For example,
            auto-generating a slug from the title when creating a blog post, or sending a welcome email
            after creating a user. GORM hooks let you do this.
          </p>

          <Definition term="Hook (Callback)">
            A method on your model that GORM calls automatically at specific points in the lifecycle of
            a database operation. Available hooks: <Code>BeforeCreate</Code>, <Code>AfterCreate</Code>,
            <Code>BeforeUpdate</Code>, <Code>AfterUpdate</Code>, <Code>BeforeSave</Code> (both create
            and update), <Code>AfterSave</Code>, <Code>BeforeDelete</Code>, <Code>AfterDelete</Code>.
          </Definition>

          <CodeBlock filename="hooks.go">
{`import "strings"

// BeforeCreate: auto-generate slug from name
func (p *Product) BeforeCreate(tx *gorm.DB) error {
    if p.Slug == "" {
        p.Slug = strings.ToLower(strings.ReplaceAll(p.Name, " ", "-"))
    }
    return nil
}

// AfterCreate: log the creation
func (p *Product) AfterCreate(tx *gorm.DB) error {
    fmt.Printf("Product created: %s (ID: %d)\n", p.Name, p.ID)
    return nil
}

// BeforeSave: validate price
func (p *Product) BeforeSave(tx *gorm.DB) error {
    if p.Price < 0 {
        return fmt.Errorf("price cannot be negative")
    }
    return nil
}

// BeforeDelete: prevent deleting products with orders
func (p *Product) BeforeDelete(tx *gorm.DB) error {
    var count int64
    tx.Model(&Order{}).Where("product_id = ?", p.ID).Count(&count)
    if count > 0 {
        return fmt.Errorf("cannot delete product with existing orders")
    }
    return nil
}`}
          </CodeBlock>

          <Note>
            If a hook returns an error, the database operation is cancelled and the error is returned
            to the caller. This makes hooks perfect for validation — if a product{"'"}s price is negative,
            the <Code>BeforeSave</Code> hook rejects the save before it reaches the database.
          </Note>

          <Challenge number={8} title="Add a Slug Hook">
            <p>Add a <Code>BeforeCreate</Code> hook to your Blog model that auto-generates a slug from
            the title. {'"'}My First Post{'"'} should become {'"'}my-first-post{'"'}. Test it: create a blog post with
            only a title, then check if the slug was set automatically. Also add a <Code>BeforeSave</Code>
            hook that ensures the title is not empty.</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Raw SQL ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Raw SQL</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            GORM{"'"}s query builder handles 90% of your database needs. But sometimes you need the full
            power of SQL — complex aggregations, window functions, database-specific features, or
            performance-critical queries. GORM lets you run raw SQL when needed.
          </p>

          <CodeBlock filename="Raw Queries">
{`// Raw SELECT — scan results into a struct
type CategoryRevenue struct {
    Category  string
    Total     int64
    AvgPrice  float64
    Revenue   float64
}

var results []CategoryRevenue
db.Raw("SELECT c.name as category, COUNT(*) as total, " +
    "AVG(p.price) as avg_price, SUM(p.price * p.stock) as revenue " +
    "FROM products p JOIN categories c ON p.category_id = c.id " +
    "GROUP BY c.name ORDER BY revenue DESC").
    Scan(&results)

// Raw UPDATE
db.Exec("UPDATE products SET price = price * 1.1 WHERE category_id = ?", categoryID)

// Raw DELETE
db.Exec("DELETE FROM sessions WHERE expires_at < NOW()")

// Check if raw query had an error
result := db.Exec("UPDATE products SET price = 0 WHERE id = ?", id)
if result.Error != nil {
    // Handle error
}
fmt.Println(result.RowsAffected) // Number of rows affected`}
          </CodeBlock>

          <Tip>
            Always use parameterized queries (<Code>?</Code> placeholders) even in raw SQL. Never
            concatenate user input into SQL strings. <Code>db.Raw({'"'}... WHERE id = ?{'"'}, id)</Code> is
            safe. <Code>db.Raw({'"'}... WHERE id = {'"'} + id)</Code> is a SQL injection vulnerability.
          </Tip>

          <Challenge number={9} title="Write a Revenue Query">
            <p>Write a raw SQL query that returns: (1) the average price per category, (2) the total
            number of products per category, and (3) the most expensive product in each category. Define
            a result struct to hold the data. Run it and check the results in GORM Studio.</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Performance Tips ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Performance Tips</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            A slow database is a slow application. Here are the most impactful performance optimizations
            you can apply to your GORM queries:
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">1. Use Indexes</strong> — An index is like a book{"'"}s table of
            contents. Without it, the database scans every row to find matches. With it, lookups are
            nearly instant.
          </p>

          <CodeBlock filename="Indexes">
{`type Product struct {
    gorm.Model
    Name       string  // gorm:"index"          — faster lookups by name
    Slug       string  // gorm:"uniqueIndex"    — fast + no duplicates
    CategoryID uint    // gorm:"index"          — faster joins and filters
    Email      string  // gorm:"uniqueIndex"    — fast + unique constraint
}

// Composite index (search by category AND active together)
// gorm:"index:idx_category_active"`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">2. Avoid N+1 Queries</strong> — The most common performance
            killer. If you load 100 products and each one triggers a separate query to load its category,
            that{"'"}s 101 queries. Use <Code>Preload</Code> to load everything in 2 queries.
          </p>

          <CodeBlock filename="N+1 Problem">
{`// BAD: N+1 queries (1 for products + N for categories)
var products []models.Product
db.Find(&products)
for _, p := range products {
    // This triggers a separate query for each product!
    fmt.Println(p.Category.Name)
}

// GOOD: 2 queries total (1 for products, 1 for all categories)
db.Preload("Category").Find(&products)
for _, p := range products {
    fmt.Println(p.Category.Name) // Already loaded
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">3. Paginate Large Results</strong> — Never load all records
            at once. Always use <Code>Offset</Code> and <Code>Limit</Code>.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">4. Select Only Needed Columns</strong> — If you only need the
            name and price, don{"'"}t load the entire row.
          </p>

          <CodeBlock filename="Efficient Queries">
{`// Only load what you need
db.Select("id", "name", "price").Find(&products)

// Paginate — never load all records
db.Scopes(Paginate(page, 20)).Find(&products)

// Count without loading data
var count int64
db.Model(&Product{}).Where("active = ?", true).Count(&count)`}
          </CodeBlock>

          <Challenge number={10} title="Audit Your Queries">
            <p>Review your product queries from the previous challenges. For each one, check:
            (1) Are you using Preload for relationships (avoiding N+1)? (2) Are you paginating
            large result sets? (3) Are you using Select to load only needed columns? (4) Do your
            frequently-filtered columns have indexes? Fix any issues you find.</p>
          </Challenge>
        </section>

        {/* ═══ Section 11: Summary + Final Challenge ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You{"'"}ve mastered GORM — from basic models to advanced query optimization:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Models</strong> — Go structs with GORM tags define your database schema</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Migrations</strong> — <Code>grit migrate</Code> creates/updates tables automatically</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">CRUD</strong> — Create, Read (First, Find, Where), Update, Delete (soft + hard)</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Relationships</strong> — belongs_to, has_many, many_to_many with Preload</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Query Building</strong> — Where, Or, Not, Order, Select, Group, Count</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Scopes</strong> — reusable, composable query conditions</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Hooks</strong> — automatic logic before/after database operations</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Raw SQL</strong> — full SQL power when GORM isn{"'"}t enough</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Performance</strong> — indexes, Preload, pagination, Select</span></li>
          </ul>

          <Challenge number={11} title="Final Challenge: Product Analytics (GORM Query Builder)">
            <p>Using GORM{"'"}s query builder (no raw SQL), write queries that return:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Total number of products per category (using Group + Count)</li>
              <li>Average price per category (using Group + Select with AVG)</li>
              <li>The most expensive product overall (using Order + First)</li>
              <li>All products with stock below 10 (using Where + Find)</li>
            </ol>
          </Challenge>

          <Challenge number={12} title="Final Challenge: Revenue Report (Raw SQL)">
            <p>Using raw SQL, write a single query that returns a revenue report:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Category name</li>
              <li>Total products in category</li>
              <li>Average price in category</li>
              <li>Total revenue (SUM of price * stock)</li>
              <li>Most expensive product name in each category</li>
            </ol>
            <p className="mt-2">Sort by total revenue descending. Define a result struct to hold the data.</p>
          </Challenge>

          <Challenge number={13} title="Final Challenge: Scopes + Hooks Combo">
            <p>Create a complete Product model with: (1) A <Code>BeforeCreate</Code> hook that generates
            a slug and validates the price, (2) An <Code>Active</Code> scope, a <Code>InPriceRange</Code> scope,
            and a <Code>Paginate</Code> scope, (3) A query that uses all 3 scopes together with Preload
            for the Category relationship.</p>
          </Challenge>

          <Challenge number={14} title="Final Challenge: Performance Optimization">
            <p>Take the queries from challenges 11 and 12. For each one: (1) Ensure all filtered columns
            have indexes, (2) Use Select to load only needed columns, (3) Use Preload instead of
            separate queries for relationships, (4) Add pagination to any query that could return
            many results. Measure the difference with GORM{"'"}s debug mode: <Code>db.Debug().Find(&products)</Code>.</p>
          </Challenge>
        </section>

        {/* ═══ Footer ═══ */}
        <CourseFooter />

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
