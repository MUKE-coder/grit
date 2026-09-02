import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/tutorials/custom-endpoints')

/* Every command, file and response on this page was run against a real
 * project before it was written. The endpoint tables are read from the
 * generated <resource>_routes.go files, not from memory. */

type Row = { method: string; path: string; what: string }

const PUBLIC_ROUTES: Row[] = [
  { method: 'GET', path: '/api/v1/public/categories', what: 'Paginated list. Read-only.' },
  { method: 'GET', path: '/api/v1/public/categories/:key', what: 'One row, by id or slug.' },
  { method: 'GET', path: '/api/v1/public/categories/:key/related', what: 'Rows that point at this one.' },
  { method: 'GET', path: '/api/v1/public/categories/tree', what: 'The whole hierarchy, one query. Only with --tree.' },
]

const PROTECTED_ROUTES: Row[] = [
  { method: 'GET', path: '/api/v1/categories', what: 'List, with search, filters, sort and pagination.' },
  { method: 'GET', path: '/api/v1/categories/:id', what: 'One row.' },
  { method: 'POST', path: '/api/v1/categories', what: 'Create.' },
  { method: 'PUT', path: '/api/v1/categories/:id', what: 'Replace.' },
  { method: 'PATCH', path: '/api/v1/categories/:id', what: 'Update the fields you send.' },
  { method: 'GET', path: '/api/v1/categories/export', what: 'CSV or XLSX of the current filter.' },
  { method: 'POST', path: '/api/v1/categories/import', what: 'Bulk import from a spreadsheet.' },
  { method: 'GET', path: '/api/v1/categories/import/template', what: 'A blank import file with the right headers.' },
  { method: 'GET', path: '/api/v1/categories/:id/pdf', what: 'The record as a PDF.' },
]

const ADMIN_ROUTES: Row[] = [
  { method: 'DELETE', path: '/api/v1/categories/:id', what: 'Soft delete. ADMIN role required.' },
  { method: 'POST', path: '/api/v1/categories/bulk', what: 'Act on many rows at once, delete included.' },
]

const TREE_ROUTES: Row[] = [
  { method: 'GET', path: '/api/v1/categories/tree', what: 'Nested rows.' },
  { method: 'GET', path: '/api/v1/categories/:id/breadcrumbs', what: 'The path from the root to this row.' },
  { method: 'PATCH', path: '/api/v1/categories/:id/move', what: 'Reparent, refusing to make a row its own ancestor.' },
  { method: 'POST', path: '/api/v1/categories/reorder', what: 'Reorder siblings.' },
  { method: 'POST', path: '/api/v1/categories/rebuild-tree', what: 'Recompute every path and depth.' },
]

const GORM_ROWS: { go: string; sql: string; note: string }[] = [
  {
    go: 'db.Find(&items)',
    sql: 'SELECT * FROM products',
    note: 'Many rows. No error when nothing matches: an empty slice is an answer.',
  },
  {
    go: 'db.First(&item, "id = ?", id)',
    sql: 'SELECT * FROM products WHERE id = ? LIMIT 1',
    note: 'One row. Returns gorm.ErrRecordNotFound when there is none, which is the 404.',
  },
  {
    go: 'db.Where("stock > ?", 0).Find(&items)',
    sql: 'SELECT * FROM products WHERE stock > 0',
    note: 'Always ? placeholders. String concatenation here is SQL injection.',
  },
  {
    go: 'db.Order("created_at DESC").Limit(10).Find(&items)',
    sql: 'ORDER BY created_at DESC LIMIT 10',
    note: 'Chainable. Nothing runs until Find, First or Count.',
  },
  {
    go: 'db.Create(&item)',
    sql: 'INSERT INTO products (...) VALUES (...)',
    note: 'Fills item.ID and the timestamps in place. Pass a pointer or it cannot.',
  },
  {
    go: 'db.Model(&item).Updates(map[string]any{...})',
    sql: 'UPDATE products SET ... WHERE id = ?',
    note: 'A map updates zero and false; a struct skips them, because it cannot tell zero from unset.',
  },
  {
    go: 'db.Delete(&models.Product{}, "id = ?", id)',
    sql: 'UPDATE products SET deleted_at = now() WHERE id = ?',
    note: 'A soft delete, because the model has gorm.DeletedAt. The row stays and stops being found.',
  },
  {
    go: 'db.Model(&models.Product{}).Count(&n)',
    sql: 'SELECT count(*) FROM products',
    note: 'Model, not Find. Counting into a slice loads every row to throw it away.',
  },
  {
    go: 'db.Preload("Category").Find(&items)',
    sql: 'two queries: products, then categories WHERE id IN (...)',
    note: 'Fills item.Category. Without it that field is nil and the JSON shows null.',
  },
  {
    go: 'db.Transaction(func(tx *gorm.DB) error { ... })',
    sql: 'BEGIN ... COMMIT / ROLLBACK',
    note: 'Return an error and everything rolls back. Use tx inside, never db.',
  },
]

function RouteTable({ rows, caption }: { rows: Row[]; caption?: string }) {
  return (
    <div className="overflow-x-auto rounded-lg border border-border mb-4">
      <table className="w-full text-sm">
        <thead className="bg-muted/40">
          <tr>
            <th className="text-left font-medium px-4 py-2.5 w-20">Method</th>
            <th className="text-left font-medium px-4 py-2.5">Path</th>
            <th className="text-left font-medium px-4 py-2.5">What it does</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((r) => (
            <tr key={r.method + r.path} className="border-t border-border align-top">
              <td className="px-4 py-2.5 font-mono text-xs text-primary whitespace-nowrap">{r.method}</td>
              <td className="px-4 py-2.5 font-mono text-xs whitespace-nowrap">{r.path}</td>
              <td className="px-4 py-2.5 text-muted-foreground text-xs">{r.what}</td>
            </tr>
          ))}
        </tbody>
      </table>
      {caption && <p className="text-xs text-muted-foreground px-4 py-2 border-t border-border">{caption}</p>}
    </div>
  )
}

function H2({ id, children }: { id: string; children: React.ReactNode }) {
  return (
    <h2 id={id} className="text-2xl font-semibold tracking-tight mb-4 scroll-mt-24">
      {children}
    </h2>
  )
}

function H3({ children }: { children: React.ReactNode }) {
  return <h3 className="text-lg font-semibold tracking-tight mb-3 mt-8">{children}</h3>
}

function P({ children }: { children: React.ReactNode }) {
  return <p className="text-muted-foreground leading-relaxed mb-4">{children}</p>
}

export default function CustomEndpointsPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Tutorial · Long read</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Custom API endpoints, end to end
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Build a small shop, see exactly which endpoints Grit gives you for nothing, then
                write four of your own: three public, one that only signed-in users can call. Then
                call all four from Next.js and from TanStack Start.
              </p>
              <p className="text-sm text-muted-foreground leading-relaxed mt-4">
                Written for someone who has never touched a handler or a service. Every command
                and every response on this page was run against a real project.
              </p>
            </div>

            <div className="prose-grit">
              {/* ---------------------------------------------------------- */}
              <div className="mb-12">
                <H2 id="what-we-build">What we are building</H2>
                <P>
                  A shop with three resources. They are chosen so that each one teaches something
                  the last one did not.
                </P>
                <ul className="space-y-2.5 mb-4">
                  {[
                    ['Category', 'A name, a slug, an image. Hierarchical, so Laptops can sit under Computers, and public, so the storefront can read it without a login.'],
                    ['Product', 'Belongs to one category. This is the simplest relationship, and the one that changes the least.'],
                    ['Campaign', 'Has many products, and a product can be in many campaigns. This is the relationship that changes the most.'],
                  ].map(([k, v]) => (
                    <li key={k} className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="text-primary mt-0.5 font-mono text-xs shrink-0 w-20">{k}</span>
                      <span>{v}</span>
                    </li>
                  ))}
                </ul>
                <P>
                  Along the way: seeding four categories from a file you control, adding a column
                  to a resource that already exists, and then the part this page is really about,
                  which is writing endpoints Grit did not generate.
                </P>
              </div>

              {/* ---------------------------------------------------------- */}
              <div className="mb-12">
                <H2 id="step-1">Step 1: the project</H2>
                <CodeBlock
                  terminal
                  code={`grit new shop --triple --next
cd shop`}
                />
                <P>
                  <code>--triple</code> is three apps in one repository: the Go API, a public web
                  app, and an admin panel. <code>--next</code> makes both frontends Next.js. They
                  share one <code>packages/shared</code>, so a type generated from a Go struct is
                  the same type in all three.
                </P>
                <P>
                  Open <code>.env</code> and point it at a database. SQLite needs nothing
                  installed and is fine for this whole tutorial:
                </P>
                <CodeBlock language="bash" filename=".env" code={`DATABASE_URL=sqlite:./app.db`} />
              </div>

              {/* ---------------------------------------------------------- */}
              <div className="mb-12">
                <H2 id="step-2">Step 2: the Category resource</H2>
                <CodeBlock
                  terminal
                  code={`grit generate resource Category \\
  --fields "name:string,slug:slug:name,image:file:image" \\
  --tree --public`}
                />
                <P>Three things there are worth slowing down on.</P>
                <ul className="space-y-2.5 mb-4">
                  {[
                    ['slug:slug:name', 'A URL-safe version of another field. The third segment names the source, so this slugifies name. Filled on save, unique, and never shown in the form: it is not something a person types.'],
                    ['--tree', 'Makes the resource hierarchical. Adds parent_id, plus path, depth and position, and five endpoints for moving rows around. A category can now sit under another category.'],
                    ['--public', 'Adds a second, read-only copy of the list and detail endpoints under /api/v1/public/, outside the login. This is what a storefront reads.'],
                  ].map(([k, v]) => (
                    <li key={k} className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="text-primary mt-0.5 font-mono text-xs shrink-0">{k}</span>
                      <span>{v}</span>
                    </li>
                  ))}
                </ul>
                <P>The command prints what it wrote. The files that matter for this page:</P>
                <CodeBlock
                  language="text"
                  code={`apps/api/internal/models/category.go        the struct, and the database table
apps/api/internal/services/category.go      business logic (yours to grow)
apps/api/internal/handlers/category.go      HTTP in, JSON out
apps/api/internal/routes/category_routes.go every URL this resource answers
packages/shared/types/category.ts           the TypeScript type
apps/web/hooks/use-categories.ts            React Query hooks
apps/admin/resources/categories/            the admin screen`}
                />
                <Callout type="note" title="One resource, one routes file">
                  Since v3.185.0 each resource owns{' '}
                  <code>internal/routes/&lt;resource&gt;_routes.go</code>. Creating that file
                  mounts the resource and deleting it unmounts it, so{' '}
                  <code>routes.go</code> does not grow as you add resources. You will add your own
                  routes file later, and it works exactly the same way.
                </Callout>
              </div>

              {/* ---------------------------------------------------------- */}
              <div className="mb-12">
                <H2 id="step-3">Step 3: seeding four categories</H2>
                <P>
                  An empty database makes everything harder to check. Grit can fill one with
                  gofakeit via <code>--faker</code>, but random words are no good here: we want four
                  categories we can name in a test. So generate the seeder on its own and edit it.
                </P>
                <CodeBlock terminal code={`grit generate seeder Category`} />
                <P>
                  That writes <code>apps/api/internal/database/categories_seeder.go</code> with one
                  example row, and registers it in <code>seed.go</code>. Replace the row with four:
                </P>
                <CodeBlock
                  language="go"
                  filename="apps/api/internal/database/categories_seeder.go"
                  code={`func SeedCategories(db *gorm.DB) error {
    // Idempotent: running it twice does not double the rows. Every generated
    // seeder starts with this, and yours should keep it.
    var count int64
    db.Model(&models.Category{}).Count(&count)
    if count > 0 {
        log.Println("Categories already seeded, skipping...")
        return nil
    }

    records := []models.Category{
        {Name: "Laptops", IsFeatured: true,
            Image: &files.FileRef{URL: "https://picsum.photos/seed/laptops/600/400", Name: "laptops.jpg", MIME: "image/jpeg"}},
        {Name: "Phones", IsFeatured: true,
            Image: &files.FileRef{URL: "https://picsum.photos/seed/phones/600/400", Name: "phones.jpg", MIME: "image/jpeg"}},
        {Name: "Audio", IsFeatured: false,
            Image: &files.FileRef{URL: "https://picsum.photos/seed/audio/600/400", Name: "audio.jpg", MIME: "image/jpeg"}},
        {Name: "Accessories", IsFeatured: false,
            Image: &files.FileRef{URL: "https://picsum.photos/seed/accessories/600/400", Name: "accessories.jpg", MIME: "image/jpeg"}},
    }

    for _, r := range records {
        if err := db.Create(&r).Error; err != nil {
            log.Printf("Warning: failed to seed category: %v", err)
        }
    }
    log.Printf("Seeded %d category(s)", len(records))
    return nil
}`}
                />
                <P>
                  Notice what is <em>not</em> in that list: no <code>ID</code>, no{' '}
                  <code>Slug</code>, no <code>Path</code> or <code>Depth</code>. The id is
                  generated on save, the slug is derived from the name, and the tree columns are
                  computed from the parent. Setting them by hand would fight the model.
                </P>
                <CodeBlock
                  terminal
                  code={`grit migrate   # create the tables
grit seed      # run every seeder

# Seeded 4 category(s)`}
                />
                <P>
                  <code>IsFeatured</code> in that seeder does not exist yet. Add it next, then come
                  back and run <code>grit seed</code>.
                </P>
              </div>

              {/* ---------------------------------------------------------- */}
              <div className="mb-12">
                <H2 id="step-4">Step 4: adding a field to a resource that already exists</H2>
                <P>
                  You do not regenerate the resource. Regenerating rewrites the model, the handler
                  and the admin screen, and takes any edit you made with it.
                </P>
                <CodeBlock terminal code={`grit generate field Category is_featured:toggle`} />
                <P>One command, five files:</P>
                <ul className="space-y-2.5 mb-4">
                  {[
                    ['The Go model', 'IsFeatured bool, added after your other fields.'],
                    ['The Zod schemas', 'Create and update both accept it.'],
                    ['The TypeScript type', 'is_featured: boolean.'],
                    ['The admin form', 'A switch.'],
                    ['The admin table', 'A yes/no column.'],
                  ].map(([k, v]) => (
                    <li key={k} className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="text-primary mt-0.5 text-xs shrink-0 w-32">{k}</span>
                      <span>{v}</span>
                    </li>
                  ))}
                </ul>
                <Callout type="note" title="No migration file">
                  The model is the source of truth for the schema. <code>grit migrate</code> reads
                  the struct and adds the column. There is nothing to write and nothing to check
                  in.
                </Callout>
                <CodeBlock terminal code={`grit migrate
grit seed`} />
                <P>
                  <code>toggle</code> is a bool that renders as a switch. For anything richer, the
                  same command takes <code>status:select:draft=Draft|live=Live</code>,{' '}
                  <code>notes:text</code>, and the other scalar types. Relationships, files and
                  slugs are the exception: those change too much, so regenerate the resource
                  instead.
                </P>
              </div>

              {/* ---------------------------------------------------------- */}
              <div className="mb-12">
                <H2 id="step-5">Step 5: Product, and what a relationship changes</H2>
                <CodeBlock
                  terminal
                  code={`grit generate resource Product \\
  --fields "name:string,slug:slug:name,description:text,price:money,stock:int,image:file:image,category:belongs_to:Category" \\
  --public`}
                />
                <P>
                  <code>category:belongs_to:Category</code> is the whole relationship. Everything
                  else follows from it.
                </P>
                <H3>What belongs_to actually does</H3>
                <P>In the model, one field becomes two:</P>
                <CodeBlock
                  language="go"
                  filename="apps/api/internal/models/product.go"
                  code={`// The column. A UUID string, indexed, and this is what the database stores.
CategoryID string \`gorm:"size:36;index" json:"category_id"\`

// The relation. NOT a column. GORM fills this when you ask it to, and
// leaves it nil when you do not.
Category *Category \`gorm:"foreignKey:CategoryID" json:"category,omitempty"\``}
                />
                <P>
                  That distinction is the single most useful thing to understand about
                  relationships in Grit. <code>CategoryID</code> is data. <code>Category</code> is a
                  convenience that costs a query. If you list a thousand products and never
                  <code> Preload</code>, you get a thousand <code>category_id</code> strings and a
                  thousand <code>null</code> categories, and the page is fast. If you do preload,
                  you get the names, and it costs one extra query for the whole page.
                </P>
                <CodeBlock
                  language="go"
                  code={`// No preload: category is null in the JSON.
db.Find(&products)

// One extra query for the whole page, not one per row.
db.Preload("Category").Find(&products)`}
                />
                <P>Elsewhere, the same field shows up as:</P>
                <ul className="space-y-2.5 mb-4">
                  {[
                    ['In the admin form', 'A searchable dropdown that loads categories from the API, rather than a text box you type a UUID into.'],
                    ['In the list endpoint', '?category_id=<id> filters by it, because the generated handler whitelists the foreign key column.'],
                    ['In TypeScript', 'category_id: string, and category?: Category.'],
                    ['In the seeder', 'Grit looks up real category ids and picks one, rather than inventing a UUID that points at nothing.'],
                  ].map(([k, v]) => (
                    <li key={k} className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="text-primary mt-0.5 text-xs shrink-0 w-36">{k}</span>
                      <span>{v}</span>
                    </li>
                  ))}
                </ul>
                <Callout type="warning" title="Order matters when seeding">
                  A product cannot point at a category that does not exist yet. Grit&apos;s
                  generated seeder loads the parent ids first and fails with a readable message if
                  there are none, rather than writing an empty foreign key and letting the database
                  answer with a constraint error nobody can read.
                </Callout>
              </div>

              {/* ---------------------------------------------------------- */}
              <div className="mb-12">
                <H2 id="step-6">Step 6: Campaign, and many-to-many</H2>
                <CodeBlock
                  terminal
                  code={`grit generate resource Campaign \\
  --fields "title:string,subtitle:string,slug:slug:title,description:text,image:file:image,products:many_to_many:Product" \\
  --public`}
                />
                <P>
                  A product can be in several campaigns and a campaign holds several products, so
                  neither table can hold the other&apos;s id. GORM creates a third table with
                  nothing in it but the two keys:
                </P>
                <CodeBlock
                  language="sql"
                  code={`campaign_products
  campaign_id  varchar(36)
  product_id   varchar(36)`}
                />
                <CodeBlock
                  language="go"
                  filename="apps/api/internal/models/campaign.go"
                  code={`// No CampaignID on Product, and no ProductID on Campaign. The join table
// is the relationship, and many2many names it.
Products []Product \`gorm:"many2many:campaign_products;" json:"products,omitempty"\``}
                />
                <H3>How it differs from belongs_to, in practice</H3>
                <div className="overflow-x-auto rounded-lg border border-border mb-4">
                  <table className="w-full text-sm">
                    <thead className="bg-muted/40">
                      <tr>
                        <th className="text-left font-medium px-4 py-2.5"> </th>
                        <th className="text-left font-medium px-4 py-2.5">belongs_to</th>
                        <th className="text-left font-medium px-4 py-2.5">many_to_many</th>
                      </tr>
                    </thead>
                    <tbody>
                      {[
                        ['Stored as', 'A column on this table', 'A separate join table'],
                        ['Set on create', 'category_id: "..."', 'product_ids: ["id", "id"]'],
                        ['Reading it', 'Preload("Category")', 'Preload("Products")'],
                        ['Changing it', 'Update the column', 'Replace the whole set'],
                        ['Admin control', 'Searchable dropdown', 'Multi-select picker'],
                        ['Can be empty', 'Yes, the column is nullable', 'Yes, no join rows'],
                      ].map((r) => (
                        <tr key={r[0]} className="border-t border-border align-top">
                          <td className="px-4 py-2.5 text-xs text-muted-foreground">{r[0]}</td>
                          <td className="px-4 py-2.5 font-mono text-xs">{r[1]}</td>
                          <td className="px-4 py-2.5 font-mono text-xs">{r[2]}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
                <P>
                  &quot;Replace the whole set&quot; is the one that surprises people. Sending{' '}
                  <code>product_ids: [&quot;a&quot;, &quot;b&quot;]</code> does not add two
                  products, it makes the campaign hold exactly those two. Sending{' '}
                  <code>product_ids: []</code> empties it. That is what a multi-select does, and it
                  is what the admin picker sends.
                </P>
                <Callout type="note" title="You can create a campaign with no products">
                  Nothing requires the relationship to be filled in at creation time. Create the
                  campaign, create the products, link them later with PUT or PATCH. The join table
                  simply has no rows until you do.
                </Callout>

                <P>
                  A belongs_to is the opposite by default: <code>category_id</code> is generated
                  with <code>binding:&quot;required&quot;</code>, so a product without one is a 422
                  rather than a row with a dangling reference. Drop the <code>binding</code> tag in
                  the request struct if you want products that are not filed anywhere yet.
                </P>
              </div>

              {/* ---------------------------------------------------------- */}
              <div className="mb-12">
                <H2 id="free-endpoints">What you already have, without writing an endpoint</H2>
                <P>
                  This is the part worth knowing before you write anything. Each of the three
                  resources above already answers every URL below. Nothing was hand-written, and
                  the list is read straight out of{' '}
                  <code>internal/routes/&lt;resource&gt;_routes.go</code>, which you can open.
                </P>

                <H3>Public: no login, API key only</H3>
                <P>
                  Only exists because the resource was generated with <code>--public</code>. Read
                  only: there is no POST, PUT or DELETE here, by design. A storefront reads; it
                  does not write.
                </P>
                <RouteTable
                  rows={PUBLIC_ROUTES}
                  caption=":key accepts an id or a slug, so /public/categories/laptops works."
                />

                <H3>Protected: a signed-in user, or an API key</H3>
                <P>
                  The everyday CRUD surface. The middleware on this group rejects anonymous
                  requests before your code runs, so a handler here never has to check.
                </P>
                <RouteTable rows={PROTECTED_ROUTES} />

                <H3>Admin: the ADMIN role as well</H3>
                <P>
                  Destructive operations sit behind a second check. Bulk is here rather than with
                  PATCH because bulk can delete, and a route is only as protected as its most
                  destructive branch.
                </P>
                <RouteTable rows={ADMIN_ROUTES} />

                <H3>Tree: only because Category used --tree</H3>
                <RouteTable rows={TREE_ROUTES} />

                <P>
                  Thirteen protected endpoints, four public ones and five tree ones for Category.
                  Product and Campaign get the same, minus the tree. Swap{' '}
                  <code>categories</code> for <code>products</code> or <code>campaigns</code> in
                  any path above and it exists.
                </P>
                <Callout type="tip" title="See them all in the running app">
                  <code>grit start</code> and open{' '}
                  <code>http://localhost:8080/docs</code>. Every endpoint, with request and
                  response shapes, generated from the router itself, so it cannot drift from what
                  the server actually serves.
                </Callout>
              </div>

              {/* ---------------------------------------------------------- */}
              <div className="mb-12">
                <H2 id="handlers-services">Handlers and services, from scratch</H2>
                <P>
                  If you have never written either, this is the section to read twice. Everything
                  after it is an application of these four paragraphs.
                </P>

                <H3>The journey of one request</H3>
                <CodeBlock
                  language="text"
                  code={`Browser
   │  GET /api/v1/public/categories/featured?limit=8
   ▼
Router          routes/catalog_routes.go
   │            matches the URL, picks the handler
   ▼
Middleware      checks the API key, the rate limit, CORS
   │            rejects here, before your code runs
   ▼
Handler         handlers/catalog.go
   │            reads ?limit, calls one service method,
   │            turns the answer into a status code
   ▼
Service         services/catalog.go
   │            the actual question: which categories are featured
   ▼
GORM            builds SQL, binds parameters, scans rows into structs
   ▼
Database`}
                />

                <H3>The handler: everything about HTTP, nothing about the business</H3>
                <P>
                  A handler is a function that takes a <code>*gin.Context</code> and returns
                  nothing. The context is the request and the response together: you read the URL,
                  the query string, the body and the signed-in user out of it, and you write the
                  status and the JSON back into it.
                </P>
                <P>A good handler does exactly four things, in order:</P>
                <ol className="space-y-2 mb-4 text-[14px] text-muted-foreground list-decimal pl-5">
                  <li>Reads what the request is asking for.</li>
                  <li>Calls one service method.</li>
                  <li>Turns an error into a status code.</li>
                  <li>Writes JSON.</li>
                </ol>
                <P>
                  If a handler is doing anything else, particularly if it contains the word{' '}
                  <code>Where</code>, that logic has nowhere to be reused from. A background job
                  cannot call a handler. Neither can a CLI command, nor another handler, nor a
                  test, without constructing a fake HTTP request first.
                </P>

                <H3>The service: everything about the business, nothing about HTTP</H3>
                <P>
                  A service is a plain struct holding a <code>*gorm.DB</code>, with methods on it.
                  It does not import gin. It does not know what a status code is. It takes
                  arguments and returns values and errors, which means anything can call it.
                </P>
                <CodeBlock
                  language="go"
                  code={`// The service says what happened.
func (s *CatalogService) CampaignBySlug(slug string) (*CampaignDetail, error) {
    ...
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, ErrNotFound
    }
}

// The handler decides what that means over HTTP.
if errors.Is(err, services.ErrNotFound) {
    fail(c, http.StatusNotFound, "NOT_FOUND", "Campaign not found")
    return
}`}
                />
                <P>
                  That split is the whole idea. &quot;Not found&quot; is a fact about the data.
                  &quot;404&quot; is a fact about HTTP. Returning a 404 from the service would make
                  the same method useless to a scheduled job, which has no response to write to.
                </P>

                <Callout type="note" title="Grit's generated handlers talk to GORM directly">
                  Worth knowing so the code does not confuse you: for the generated CRUD, the
                  handler queries the database itself through the shared{' '}
                  <code>paginate</code> helper, and the generated{' '}
                  <code>services/&lt;resource&gt;.go</code> is a starting point nothing calls yet.
                  Plain CRUD has no business logic to put in a service, and a method that only
                  forwards a call is a layer that costs a file and pays nothing. The moment you
                  have a real rule, that file is where it goes, and everything below shows how.
                </Callout>

                <H3>Where the database connection comes from</H3>
                <P>
                  You never open a connection. <code>cmd/server/main.go</code> opens one pool at
                  boot, and every request borrows a connection from it and gives it back. Opening
                  one per request would exhaust the database in minutes.
                </P>
                <CodeBlock
                  language="go"
                  code={`// cmd/server/main.go, at boot, once.
db, err := database.Connect(cfg)

// routes/routes.go hands it to every resource's routes file.
mountResources(&Mount{ DB: db, ... })

// your routes file puts it in the handler.
h := handlers.NewCatalogHandler(m.DB)

// the handler passes it to the service.
&services.CatalogService{DB: db}`}
                />
                <P>
                  One <code>*gorm.DB</code>, passed down. That is the whole dependency chain, and
                  it is why a handler can be constructed in a test with an in-memory SQLite
                  database and no server at all.
                </P>
              </div>

              {/* ---------------------------------------------------------- */}
              <div className="mb-12">
                <H2 id="gorm-cheatsheet">GORM cheat sheet</H2>
                <P>
                  GORM turns Go structs into SQL. Everything below is chainable, and nothing
                  executes until a finisher: <code>Find</code>, <code>First</code>,{' '}
                  <code>Count</code>, <code>Create</code>, <code>Updates</code>,{' '}
                  <code>Delete</code>.
                </P>
                <div className="overflow-x-auto rounded-lg border border-border mb-4">
                  <table className="w-full text-sm">
                    <thead className="bg-muted/40">
                      <tr>
                        <th className="text-left font-medium px-4 py-2.5">Go</th>
                        <th className="text-left font-medium px-4 py-2.5">SQL</th>
                        <th className="text-left font-medium px-4 py-2.5">Worth knowing</th>
                      </tr>
                    </thead>
                    <tbody>
                      {GORM_ROWS.map((r) => (
                        <tr key={r.go} className="border-t border-border align-top">
                          <td className="px-4 py-2.5 font-mono text-[11px] whitespace-nowrap">{r.go}</td>
                          <td className="px-4 py-2.5 font-mono text-[11px] text-muted-foreground">{r.sql}</td>
                          <td className="px-4 py-2.5 text-muted-foreground text-xs">{r.note}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                <Callout type="warning" title="The three that catch everyone">
                  <strong>Pass pointers.</strong> <code>db.Create(item)</code> cannot write the
                  generated id back; <code>db.Create(&amp;item)</code> can.{' '}
                  <strong>Updates with a struct skips zero values</strong>, so setting a bool to
                  false or a count to 0 silently does nothing: use a map.{' '}
                  <strong>Find never errors on an empty result</strong>, so checking{' '}
                  <code>err != nil</code> after it will not tell you the list was empty; check the
                  length.
                </Callout>

                <H3>N+1, the one that ruins a page</H3>
                <CodeBlock
                  language="go"
                  code={`// 1 query, then 1 more per product. 200 products is 201 queries.
db.Find(&products)
for i := range products {
    db.First(&products[i].Category, "id = ?", products[i].CategoryID)
}

// 2 queries, whatever the number of products.
db.Preload("Category").Find(&products)`}
                />
                <P>
                  It is invisible with ten rows in development and obvious with ten thousand in
                  production. If a loop contains a query, that is the shape.
                </P>
              </div>

              {/* ---------------------------------------------------------- */}
              <div className="mb-12">
                <H2 id="custom-public">Writing your own endpoints</H2>
                <P>
                  Four of them: three public, one that requires a login. Same pattern each time,
                  three files.
                </P>
                <ol className="space-y-2 mb-4 text-[14px] text-muted-foreground list-decimal pl-5">
                  <li><code>services/catalog.go</code>: the question and its answer.</li>
                  <li><code>handlers/catalog.go</code>: request in, status and JSON out.</li>
                  <li><code>routes/catalog_routes.go</code>: the URLs.</li>
                </ol>
                <P>
                  Name them after what they do, not after a model. These queries all answer
                  &quot;what does the shop front page need&quot;, so one <code>catalog</code>{' '}
                  triple holds all four rather than scattering one method into each of three
                  files.
                </P>

                <H3>1. Featured categories, newest first</H3>
                <CodeBlock
                  language="go"
                  filename="apps/api/internal/services/catalog.go"
                  code={`package services

import (
    "errors"
    "fmt"

    "gorm.io/gorm"

    "shop/apps/api/internal/models"
)

// CatalogService is the business logic behind the storefront's own endpoints.
type CatalogService struct {
    DB *gorm.DB
}

// FeaturedCategories returns the categories marked featured, newest first.
//
// limit is capped rather than trusted. It arrives from a query string, and
// ?limit=100000 on a public endpoint is a cheap way to make the database do
// expensive work.
func (s *CatalogService) FeaturedCategories(limit int) ([]models.Category, error) {
    if limit < 1 || limit > 50 {
        limit = 12
    }

    var items []models.Category
    err := s.DB.
        Where("is_featured = ?", true).
        Where("archived_at IS NULL").
        Order("created_at DESC").
        Limit(limit).
        Find(&items).Error
    if err != nil {
        return nil, fmt.Errorf("fetching featured categories: %w", err)
    }
    return items, nil
}`}
                />
                <P>
                  <code>archived_at IS NULL</code> is not optional. Grit gives every resource an
                  archive that hides rows from the admin&apos;s default list, and a custom query
                  that forgets it shows the storefront things somebody deliberately put away.
                </P>

                <H3>2. Categories with a product count</H3>
                <CodeBlock
                  language="go"
                  filename="apps/api/internal/services/catalog.go"
                  code={`// CategoryWithCount is a category plus how many products point at it.
//
// Its own type rather than models.Category, because product_count is not a
// column. Adding a non-column field to the model would make it appear in the
// admin table, the generated TypeScript type and the API reference as though
// it were stored.
type CategoryWithCount struct {
    models.Category
    ProductCount int64 \`json:"product_count"\`
}

// CategoriesWithProductCounts returns every live category and its product
// count, in one query.
//
// The obvious version loads the categories and then counts products for each
// one: one query plus one per row, the N+1 above. A LEFT JOIN with GROUP BY
// asks the database the whole question once.
//
// LEFT, not INNER: an inner join drops categories with no products, and a
// category with nothing in it is exactly the one an admin is looking for.
func (s *CatalogService) CategoriesWithProductCounts() ([]CategoryWithCount, error) {
    var rows []CategoryWithCount
    err := s.DB.
        Model(&models.Category{}).
        Select("categories.*, COUNT(products.id) AS product_count").
        Joins("LEFT JOIN products ON products.category_id = categories.id AND products.deleted_at IS NULL").
        Where("categories.archived_at IS NULL").
        Group("categories.id").
        Order("product_count DESC").
        Find(&rows).Error
    if err != nil {
        return nil, fmt.Errorf("counting products per category: %w", err)
    }
    return rows, nil
}`}
                />
                <P>
                  <code>products.deleted_at IS NULL</code> lives in the JOIN condition, not the
                  WHERE. In the WHERE it would drop categories whose only products are deleted,
                  turning the LEFT JOIN back into an INNER one.
                </P>

                <H3>3. A campaign detail page, by slug</H3>
                <CodeBlock
                  language="go"
                  filename="apps/api/internal/services/catalog.go"
                  code={`// ErrNotFound is what a handler turns into a 404.
//
// The service says "not found"; the handler decides that means 404. Returning
// a gin status from here would tie the business logic to HTTP, and the same
// method is then unusable from a background job or a CLI command.
var ErrNotFound = errors.New("not found")

type CampaignDetail struct {
    Campaign models.Campaign  \`json:"campaign"\`
    Products []models.Product \`json:"products"\`
    Likes    int64            \`json:"likes"\`
    Liked    bool             \`json:"liked"\`
}

// CampaignBySlug loads a campaign's detail page.
//
// By slug, not id: /campaigns/summer-sale is the URL a person can read and a
// search engine can index.
//
// viewerID is empty for a signed-out visitor, and Liked is then false without
// a query. Asking "has nobody liked this" is a question with a known answer.
func (s *CatalogService) CampaignBySlug(slug, viewerID string) (*CampaignDetail, error) {
    var campaign models.Campaign

    // Preload pulls the many-to-many in a second query keyed by the ids found
    // in the first, rather than a join that repeats every campaign column once
    // per product.
    err := s.DB.
        Preload("Products").
        Where("slug = ?", slug).
        Where("archived_at IS NULL").
        First(&campaign).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, ErrNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("fetching campaign %q: %w", slug, err)
    }

    detail := CampaignDetail{Campaign: campaign, Products: campaign.Products}

    if err := s.DB.Model(&models.CampaignLike{}).
        Where("campaign_id = ?", campaign.ID).
        Count(&detail.Likes).Error; err != nil {
        return nil, fmt.Errorf("counting likes: %w", err)
    }

    if viewerID != "" {
        var mine int64
        if err := s.DB.Model(&models.CampaignLike{}).
            Where("campaign_id = ? AND user_id = ?", campaign.ID, viewerID).
            Count(&mine).Error; err != nil {
            return nil, fmt.Errorf("checking like: %w", err)
        }
        detail.Liked = mine > 0
    }

    return &detail, nil
}`}
                />

                <H3>4. Liking a campaign, for signed-in users only</H3>
                <P>First a model, because a like has to be stored somewhere:</P>
                <CodeBlock
                  language="go"
                  filename="apps/api/internal/models/campaign_like.go"
                  code={`package models

import (
    "time"

    "gorm.io/gorm"

    "shop/apps/api/internal/ids"
)

// CampaignLike is one user liking one campaign.
//
// A join row rather than a counter on Campaign. A counter cannot answer "has
// this user already liked it", which is the first thing the button needs to
// know, and two people liking at once would race on a read-modify-write.
type CampaignLike struct {
    ID string \`gorm:"primarykey;size:36" json:"id"\`

    // The pair is unique, so liking twice is refused by the database rather
    // than by a check that two concurrent requests can both pass.
    UserID     string \`gorm:"size:36;not null;uniqueIndex:idx_campaign_like,priority:1" json:"user_id"\`
    CampaignID string \`gorm:"size:36;not null;uniqueIndex:idx_campaign_like,priority:2" json:"campaign_id"\`

    CreatedAt time.Time \`json:"created_at"\`
}

func (m *CampaignLike) BeforeCreate(tx *gorm.DB) error {
    if m.ID == "" {
        m.ID = ids.New()
    }
    return nil
}`}
                />
                <P>
                  Register it so <code>grit migrate</code> creates the table. Open{' '}
                  <code>internal/models/user.go</code> and add one line above the marker:
                </P>
                <CodeBlock
                  language="go"
                  filename="apps/api/internal/models/user.go"
                  code={`func Models() []interface{} {
    return []interface{}{
        &User{},
        // ...
        &CampaignLike{},
        // grit:models
    }
}`}
                />
                <Callout type="warning" title="Keep the marker">
                  <code>// grit:models</code> is where <code>grit generate resource</code> inserts
                  the next model. Delete it and future resources are generated but never
                  migrated, and the failure is a missing table at runtime rather than an error at
                  generation time.
                </Callout>
                <CodeBlock
                  language="go"
                  filename="apps/api/internal/services/catalog.go"
                  code={`// LikeCampaign records that a user likes a campaign, and returns the new total.
//
// Idempotent: liking twice leaves one row and is not an error. The button can
// be double-clicked and the request can be retried by a flaky network, and
// neither should produce a second like or a red toast.
//
// FirstOrCreate rather than "check, then insert", because two requests can
// both pass the check. The unique index is what actually enforces this;
// FirstOrCreate just avoids hitting it in the common case.
func (s *CatalogService) LikeCampaign(campaignID, userID string) (int64, error) {
    var campaign models.Campaign
    err := s.DB.Select("id").Where("id = ?", campaignID).First(&campaign).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return 0, ErrNotFound
    }
    if err != nil {
        return 0, fmt.Errorf("loading campaign: %w", err)
    }

    like := models.CampaignLike{CampaignID: campaignID, UserID: userID}
    if err := s.DB.
        Where("campaign_id = ? AND user_id = ?", campaignID, userID).
        FirstOrCreate(&like).Error; err != nil {
        return 0, fmt.Errorf("liking campaign: %w", err)
    }

    return s.campaignLikes(campaignID)
}

// UnlikeCampaign removes a like. Removing one that is not there is not an
// error either, for the same reason.
func (s *CatalogService) UnlikeCampaign(campaignID, userID string) (int64, error) {
    if err := s.DB.
        Where("campaign_id = ? AND user_id = ?", campaignID, userID).
        Delete(&models.CampaignLike{}).Error; err != nil {
        return 0, fmt.Errorf("unliking campaign: %w", err)
    }
    return s.campaignLikes(campaignID)
}

func (s *CatalogService) campaignLikes(campaignID string) (int64, error) {
    var n int64
    err := s.DB.Model(&models.CampaignLike{}).
        Where("campaign_id = ?", campaignID).
        Count(&n).Error
    return n, err
}`}
                />

                <H3>The handler</H3>
                <CodeBlock
                  language="go"
                  filename="apps/api/internal/handlers/catalog.go"
                  code={`package handlers

import (
    "errors"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "shop/apps/api/internal/services"
)

// CatalogHandler holds a service, not a *gorm.DB. Every method below reads the
// request, calls one service method, and turns the result into a status code.
type CatalogHandler struct {
    Catalog *services.CatalogService
}

func NewCatalogHandler(db *gorm.DB) *CatalogHandler {
    return &CatalogHandler{Catalog: &services.CatalogService{DB: db}}
}

// fail() is not defined here. It already exists in this package, in
// recovery.go, and writes the { "error": { "code", "message" } } envelope every
// Grit endpoint uses. Handlers share one package, so a second copy is a
// compile error rather than a quiet inconsistency.

// GET /api/v1/public/categories/featured?limit=8
func (h *CatalogHandler) FeaturedCategories(c *gin.Context) {
    // Atoi returns 0 on anything unparseable, and the service treats 0 as
    // "use the default". A missing or nonsense limit is not worth a 400.
    limit, _ := strconv.Atoi(c.Query("limit"))

    items, err := h.Catalog.FeaturedCategories(limit)
    if err != nil {
        fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch featured categories")
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": items})
}

// GET /api/v1/public/categories/with-counts
func (h *CatalogHandler) CategoriesWithCounts(c *gin.Context) {
    rows, err := h.Catalog.CategoriesWithProductCounts()
    if err != nil {
        fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch categories")
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": rows})
}

// GET /api/v1/public/campaigns/:key/detail
//
// The parameter is :key, not :slug, and that is not a preference. Gin builds
// one routing tree, so two routes sharing a path position must name the
// wildcard there identically. The generated route is /campaigns/:key, so this
// one has to be :key too, and registering :slug panics the router at boot.
func (h *CatalogHandler) CampaignDetail(c *gin.Context) {
    detail, err := h.Catalog.CampaignBySlug(c.Param("key"), c.GetString("user_id"))
    if errors.Is(err, services.ErrNotFound) {
        fail(c, http.StatusNotFound, "NOT_FOUND", "Campaign not found")
        return
    }
    if err != nil {
        fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch campaign")
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": detail})
}

// POST /api/v1/campaigns/:id/like
//
// The user id comes from the token, never from the body: a user_id in the
// request would let anyone like on somebody else's behalf.
func (h *CatalogHandler) LikeCampaign(c *gin.Context) {
    userID := c.GetString("user_id")
    if userID == "" {
        fail(c, http.StatusUnauthorized, "UNAUTHORIZED", "Sign in to like a campaign")
        return
    }

    likes, err := h.Catalog.LikeCampaign(c.Param("id"), userID)
    if errors.Is(err, services.ErrNotFound) {
        fail(c, http.StatusNotFound, "NOT_FOUND", "Campaign not found")
        return
    }
    if err != nil {
        fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to like campaign")
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "data":    gin.H{"likes": likes, "liked": true},
        "message": "Campaign liked",
    })
}

// DELETE /api/v1/campaigns/:id/like
//
// The mirror image, and worth having: a like button that cannot be undone is a
// like button people stop pressing.
func (h *CatalogHandler) UnlikeCampaign(c *gin.Context) {
    userID := c.GetString("user_id")
    if userID == "" {
        fail(c, http.StatusUnauthorized, "UNAUTHORIZED", "Sign in to like a campaign")
        return
    }

    likes, err := h.Catalog.UnlikeCampaign(c.Param("id"), userID)
    if err != nil {
        fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to unlike campaign")
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "data":    gin.H{"likes": likes, "liked": false},
        "message": "Campaign unliked",
    })
}`}
                />

                <H3>The routes</H3>
                <CodeBlock
                  language="go"
                  filename="apps/api/internal/routes/catalog_routes.go"
                  code={`package routes

import (
    "shop/apps/api/internal/handlers"
)

// Storefront routes: the endpoints this project added by hand.
//
// Its own file, next to the generated <resource>_routes.go files, and it works
// the same way: an init() that registers itself, so nothing else has to know
// this file exists. "grit generate resource" will never touch it, and
// "grit remove resource" will never delete it.
func init() {
    RegisterRoutes(func(m *Mount) {
        h := handlers.NewCatalogHandler(m.DB)

        // Public: no login. The group already requires an API key.
        //
        // "/categories/featured" is registered alongside the generated
        // "/categories/:key". Gin routes a static segment ahead of a
        // parameter, so /categories/featured reaches this and
        // /categories/laptops reaches the generated handler.
        m.Public.GET("/categories/featured", h.FeaturedCategories)
        m.Public.GET("/categories/with-counts", h.CategoriesWithCounts)
        m.Public.GET("/campaigns/:key/detail", h.CampaignDetail)

        // Protected: a valid JWT or API key. The middleware on this group has
        // already rejected everyone else before the handler runs.
        m.Protected.POST("/campaigns/:id/like", h.LikeCampaign)
        m.Protected.DELETE("/campaigns/:id/like", h.UnlikeCampaign)
    })
}`}
                />
                <P>
                  Which group you register on <em>is</em> the security decision. There is no check
                  inside <code>FeaturedCategories</code> because <code>m.Public</code> already
                  said what it is, and no check inside <code>LikeCampaign</code> beyond reading
                  the user, because <code>m.Protected</code> already rejected anyone without a
                  token.
                </P>
                <div className="overflow-x-auto rounded-lg border border-border mb-4">
                  <table className="w-full text-sm">
                    <thead className="bg-muted/40">
                      <tr>
                        <th className="text-left font-medium px-4 py-2.5">Group</th>
                        <th className="text-left font-medium px-4 py-2.5">Prefix</th>
                        <th className="text-left font-medium px-4 py-2.5">Who gets through</th>
                      </tr>
                    </thead>
                    <tbody>
                      {[
                        ['m.Public', '/api/v1/public/', 'Anyone with a valid API key. No user.'],
                        ['m.Protected', '/api/v1/', 'A valid JWT, or an API key.'],
                        ['m.Admin', '/api/v1/', 'The above, and the ADMIN role.'],
                        ['m.V1', '/api/v1/', 'No middleware at all. You add your own.'],
                      ].map((r) => (
                        <tr key={r[0]} className="border-t border-border align-top">
                          <td className="px-4 py-2.5 font-mono text-xs text-primary whitespace-nowrap">{r[0]}</td>
                          <td className="px-4 py-2.5 font-mono text-xs whitespace-nowrap">{r[1]}</td>
                          <td className="px-4 py-2.5 text-muted-foreground text-xs">{r[2]}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
                <CodeBlock
                  terminal
                  code={`grit migrate            # creates campaign_likes
cd apps/api && go build ./...
grit start`}
                />
                <P>Then check they answer:</P>
                <CodeBlock
                  terminal
                  code={`# public: needs an API key, no login
curl -H "X-API-Key: $KEY" \\
  http://localhost:8080/api/v1/public/categories/featured?limit=8
# {"data":[{"name":"Phones",...},{"name":"Laptops",...}]}

curl -H "X-API-Key: $KEY" \\
  http://localhost:8080/api/v1/public/categories/with-counts
# {"data":[{"name":"Accessories","product_count":9}, ...]}

# protected: needs a token
curl -X POST -H "Authorization: Bearer $TOKEN" \\
  http://localhost:8080/api/v1/campaigns/$ID/like
# {"data":{"liked":true,"likes":1},"message":"Campaign liked"}

# without one
curl -X POST http://localhost:8080/api/v1/campaigns/$ID/like
# 401 {"error":{"code":"UNAUTHORIZED","message":"Authentication required"}}`}
                />
                <Callout type="tip" title="Where the API key comes from">
                  The admin panel, under System &rarr; API Keys, or{' '}
                  <code>POST /api/v1/api-keys</code> while signed in. Use a{' '}
                  <strong>publishable</strong> key for a storefront: it is the one meant to be
                  visible in a browser. The secret kind belongs on a server only.
                </Callout>
              </div>

              {/* ---------------------------------------------------------- */}
              <div className="mb-12">
                <H2 id="nextjs">Calling them from Next.js</H2>
                <P>
                  Two shapes, and picking the right one is most of the work. Public data is
                  fetched on the server, where the API key never reaches the browser. Anything
                  tied to the signed-in user is fetched in the browser, where the session lives.
                </P>

                <H3>Public data, in a server component</H3>
                <CodeBlock
                  language="typescript"
                  filename="apps/web/lib/catalog.ts"
                  code={`import type { Campaign, Category, Product } from "@repo/shared/types";

// No NEXT_PUBLIC_ prefix. That prefix is what ships a variable to the browser,
// and this key must not go there.
const API = process.env.API_URL ?? "http://localhost:8080";
const KEY = process.env.GRIT_API_KEY!;

export type CategoryWithCount = Category & { product_count: number };

// Mirrors the CampaignDetail struct the service returns. Category, Campaign
// and Product all come from @repo/shared/types, generated from the Go models,
// so those three cannot drift from the API.
export type CampaignDetail = {
  campaign: Campaign;
  products: Product[];
  likes: number;
  liked: boolean;
};

async function publicGet<T>(path: string, revalidate = 60): Promise<T> {
  const res = await fetch(\`\${API}/api/v1/public\${path}\`, {
    headers: { "X-API-Key": KEY },
    // Cache the response for a minute. Without this, Next fetches on every
    // request and the API sees your traffic rather than your cache.
    next: { revalidate },
  });
  if (!res.ok) {
    throw new Error(\`\${path} failed: \${res.status}\`);
  }
  const body = await res.json();
  return body.data as T;
}

export const getFeaturedCategories = (limit = 8) =>
  publicGet<Category[]>(\`/categories/featured?limit=\${limit}\`);

export const getCategoriesWithCounts = () =>
  publicGet<CategoryWithCount[]>("/categories/with-counts");

export const getCampaign = (slug: string) =>
  publicGet<CampaignDetail>(\`/campaigns/\${slug}/detail\`);`}
                />
                <CodeBlock
                  language="tsx"
                  filename="apps/web/app/page.tsx"
                  code={`import Image from "next/image";
import Link from "next/link";
import { getFeaturedCategories } from "@/lib/catalog";

// A server component. It runs on the server, so the API key stays there and
// the browser is sent finished HTML.
export default async function HomePage() {
  const categories = await getFeaturedCategories(8);

  return (
    <section className="grid grid-cols-2 md:grid-cols-4 gap-4">
      {categories.map((c) => (
        <Link key={c.id} href={\`/categories/\${c.slug}\`} className="group">
          {c.image && (
            <Image
              src={c.image.url}
              alt={c.name}
              width={600}
              height={400}
              className="rounded-lg object-cover aspect-[3/2]"
            />
          )}
          <h3 className="mt-2 font-medium group-hover:underline">{c.name}</h3>
        </Link>
      ))}
    </section>
  );
}`}
                />
                <CodeBlock
                  language="tsx"
                  filename="apps/web/app/campaigns/[slug]/page.tsx"
                  code={`import { notFound } from "next/navigation";
import { getCampaign } from "@/lib/catalog";
import { LikeButton } from "./like-button";

export default async function CampaignPage({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await params;

  let detail;
  try {
    detail = await getCampaign(slug);
  } catch {
    // The service returned ErrNotFound, the handler made it a 404, and the
    // fetch helper threw. This page turns that into Next's own 404.
    notFound();
  }

  return (
    <article>
      <h1 className="text-3xl font-bold">{detail.campaign.title}</h1>
      <p className="text-muted-foreground">{detail.campaign.subtitle}</p>

      {/* Server-rendered count, then a client component takes over. */}
      <LikeButton
        campaignId={detail.campaign.id}
        initialLikes={detail.likes}
        initialLiked={detail.liked}
      />

      <ul className="mt-8 grid grid-cols-3 gap-4">
        {detail.products.map((p) => (
          <li key={p.id}>{p.name}</li>
        ))}
      </ul>
    </article>
  );
}`}
                />

                <H3>The protected call, in a client component</H3>
                <P>
                  Liking needs the signed-in user, so it happens in the browser. The generated{' '}
                  <code>lib/api.ts</code> already attaches the session cookie and the CSRF header,
                  so there is no token to pass by hand.
                </P>
                <CodeBlock
                  language="tsx"
                  filename="apps/web/app/campaigns/[slug]/like-button.tsx"
                  code={`"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { apiClient } from "@/lib/api";

export function LikeButton({
  campaignId,
  initialLikes,
  initialLiked,
}: {
  campaignId: string;
  initialLikes: number;
  initialLiked: boolean;
}) {
  const [likes, setLikes] = useState(initialLikes);
  const [liked, setLiked] = useState(initialLiked);
  const qc = useQueryClient();

  const toggle = useMutation({
    mutationFn: async () => {
      // apiClient rewrites /api/... to /api/v1/... and sends the session
      // cookie, so this is the protected endpoint, not the public one.
      const { data } = liked
        ? await apiClient.delete(\`/api/campaigns/\${campaignId}/like\`)
        : await apiClient.post(\`/api/campaigns/\${campaignId}/like\`);
      return data.data as { likes: number; liked: boolean };
    },
    onSuccess: (d) => {
      setLikes(d.likes);
      setLiked(d.liked);
      qc.invalidateQueries({ queryKey: ["campaigns", campaignId] });
    },
    onError: (err: any) => {
      // 401 means signed out. Everything else is a real failure.
      if (err?.response?.status === 401) {
        window.location.href = "/login";
      }
    },
  });

  return (
    <button
      onClick={() => toggle.mutate()}
      disabled={toggle.isPending}
      aria-pressed={liked}
      className="inline-flex items-center gap-2 rounded-lg border px-4 py-2 disabled:opacity-50"
    >
      <span aria-hidden>{liked ? "♥" : "♡"}</span>
      {likes} {likes === 1 ? "like" : "likes"}
    </button>
  );
}`}
                />
                <Callout type="warning" title="Do not call the protected endpoint from a server component">
                  A server component has no browser cookies unless you forward them by hand, so the
                  request arrives anonymous and the API answers 401. Server components read public
                  data; the browser does anything that depends on who is signed in.
                </Callout>
              </div>

              {/* ---------------------------------------------------------- */}
              <div className="mb-12">
                <H2 id="tanstack">Calling them from TanStack Start</H2>
                <P>
                  Same two shapes. A server function replaces the server component, and the client
                  half is identical because both frontends use React Query.
                </P>
                <CodeBlock
                  language="typescript"
                  filename="src/lib/catalog.ts"
                  code={`import { createServerFn } from "@tanstack/react-start";
import type { Category } from "@repo/shared/types";

// createServerFn keeps this on the server, so the API key is never bundled.
export const getFeaturedCategories = createServerFn({ method: "GET" })
  .validator((limit: number) => limit)
  .handler(async ({ data: limit }) => {
    const res = await fetch(
      \`\${process.env.API_URL}/api/v1/public/categories/featured?limit=\${limit}\`,
      { headers: { "X-API-Key": process.env.GRIT_API_KEY! } },
    );
    if (!res.ok) throw new Error(\`featured failed: \${res.status}\`);
    const body = await res.json();
    return body.data as Category[];
  });`}
                />
                <CodeBlock
                  language="tsx"
                  filename="src/routes/index.tsx"
                  code={`import { createFileRoute } from "@tanstack/react-router";
import { getFeaturedCategories } from "@/lib/catalog";

export const Route = createFileRoute("/")({
  // The loader runs before the component renders, so there is no spinner and
  // no layout shift on first paint.
  loader: () => getFeaturedCategories({ data: 8 }),
  component: Home,
});

function Home() {
  const categories = Route.useLoaderData();

  return (
    <section className="grid grid-cols-2 md:grid-cols-4 gap-4">
      {categories.map((c) => (
        <a key={c.id} href={\`/categories/\${c.slug}\`}>
          {c.image && <img src={c.image.url} alt={c.name} className="rounded-lg" />}
          <h3 className="mt-2 font-medium">{c.name}</h3>
        </a>
      ))}
    </section>
  );
}`}
                />
                <CodeBlock
                  language="tsx"
                  filename="src/components/like-button.tsx"
                  code={`import { useMutation } from "@tanstack/react-query";
import { useState } from "react";
import { apiClient } from "@/lib/api";

// Identical to the Next.js version. Both apps use React Query and the same
// generated apiClient, so anything you write for one works in the other.
export function LikeButton({ campaignId, initialLikes, initialLiked }: {
  campaignId: string; initialLikes: number; initialLiked: boolean;
}) {
  const [likes, setLikes] = useState(initialLikes);
  const [liked, setLiked] = useState(initialLiked);

  const toggle = useMutation({
    mutationFn: async () => {
      const { data } = liked
        ? await apiClient.delete(\`/api/campaigns/\${campaignId}/like\`)
        : await apiClient.post(\`/api/campaigns/\${campaignId}/like\`);
      return data.data as { likes: number; liked: boolean };
    },
    onSuccess: (d) => { setLikes(d.likes); setLiked(d.liked); },
  });

  return (
    <button onClick={() => toggle.mutate()} disabled={toggle.isPending} aria-pressed={liked}>
      {liked ? "♥" : "♡"} {likes}
    </button>
  );
}`}
                />
              </div>

              {/* ---------------------------------------------------------- */}
              <div className="mb-12">
                <H2 id="best-practices">Best practices, and the mistakes to skip</H2>
                <div className="space-y-4 mb-6">
                  {[
                    ['Put queries in services, not handlers', 'A handler is unreachable from a job, a command or a test without faking an HTTP request. The rule of thumb: if a handler contains the word Where, the logic is in the wrong file.'],
                    ['Return errors, not status codes, from services', 'The service says ErrNotFound; the handler decides that is a 404. Keeps the same method usable from a cron job.'],
                    ['Never interpolate user input into SQL', 'Always ?. Column names cannot be parameters, so if one has to come from a request, check it against a whitelist first. Grit’s generated handlers do exactly this for sort_by.'],
                    ['Cap anything a caller can size', 'limit, page_size, and the number of ids in a bulk request. ?limit=100000 on a public endpoint is free for the caller and expensive for you.'],
                    ['Respect archived_at in custom queries', 'The generated endpoints hide archived rows. A custom query that forgets shows the storefront things somebody put away on purpose.'],
                    ['Wrap multi-write operations in a transaction', 'db.Transaction(func(tx *gorm.DB) error {...}) and use tx inside. Returning an error rolls back everything.'],
                    ['Take the user id from the token, never the body', 'c.GetString("user_id") came from a signature this server checked. A user_id field in JSON came from the caller.'],
                    ['Make write endpoints idempotent where you can', 'Liking twice should leave one like, not error. Buttons get double-clicked and networks retry.'],
                  ].map(([title, body]) => (
                    <div key={title} className="rounded-lg border border-border p-4">
                      <p className="font-medium text-sm mb-1">{title}</p>
                      <p className="text-xs text-muted-foreground leading-relaxed">{body}</p>
                    </div>
                  ))}
                </div>

                <H3>Three that cost an afternoon</H3>
                <Callout type="warning" title="Wildcard names must match">
                  Gin has one routing tree. If the generated route is{' '}
                  <code>/campaigns/:key</code>, your custom route on the same segment must also say{' '}
                  <code>:key</code>. Using <code>:slug</code> panics the router at boot with a
                  message about conflicting wildcards, and reading <code>c.Param(&quot;slug&quot;)</code>{' '}
                  when the route said <code>:key</code> silently gives you an empty string and a
                  404 you cannot explain.
                </Callout>
                <Callout type="warning" title="Many-to-many is set through *_ids">
                  The field is <code>products</code> but the JSON key is{' '}
                  <code>product_ids</code>, singularised, and it holds ids rather than objects.
                  POST, PUT and PATCH all accept it, and all three replace the whole set rather
                  than adding to it.
                </Callout>
                <Callout type="warning" title="A new file in handlers/ shares one package">
                  Every file under <code>internal/handlers</code> is package{' '}
                  <code>handlers</code>, so a helper called <code>fail</code> that already exists
                  in another file is a compile error, not a shadow. That is the good outcome: the
                  alternative is two error formats in one API.
                </Callout>
              </div>

              {/* Nav */}
              <div className="flex flex-col sm:flex-row gap-3 pt-6 border-t border-border">
                <Button asChild variant="outline" className="justify-start">
                  <Link href="/docs/concepts/generated-files">
                    <ArrowLeft className="mr-2 h-4 w-4" />
                    Generated File Map
                  </Link>
                </Button>
                <Button asChild variant="outline" className="justify-start sm:ml-auto">
                  <Link href="/docs/admin/relationships">
                    Relationships in depth
                    <ArrowRight className="ml-2 h-4 w-4" />
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
