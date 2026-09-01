import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/admin/relationships')

export default function RelationshipsPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Admin Panel</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Relationships
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Generate resources with relationships &mdash; <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">belongs_to</code> for
                foreign keys and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">many_to_many</code> for
                junction tables. The code generator handles the Go model, API handlers with eager
                loading, Zod schemas, TypeScript types, and admin form components automatically.
              </p>
              <LaneFlow
                id="admin-rel"
                lanes={['Relationship', 'Generated across the stack']}
                nodes={[
                  { id: 'bt', lane: 0, row: 0, title: 'belongs_to', sub: 'category:belongs_to', tone: 'cyan' },
                  { id: 'mtm', lane: 0, row: 3, title: 'many_to_many', sub: 'tags:many_to_many', tone: 'violet' },
                  { id: 'fk', lane: 1, row: 0, title: 'FK column', sub: '<name>_id', tone: 'green' },
                  { id: 'eager', lane: 1, row: 1, title: 'Eager loading', sub: 'Preload', tone: 'blue' },
                  { id: 'join', lane: 1, row: 2, title: 'Join table', sub: 'GORM', tone: 'green' },
                  { id: 'picker', lane: 1, row: 3, title: 'Admin picker', sub: 'searchable select', tone: 'amber' },
                ]}
                edges={[
                  { from: 'bt', to: 'fk', tone: 'green' },
                  { from: 'bt', to: 'eager', label: 'preload', tone: 'blue' },
                  { from: 'mtm', to: 'join', tone: 'green' },
                  { from: 'mtm', to: 'picker', label: 'picker', tone: 'amber' },
                ]}
                legend={[
                  { tone: 'cyan', label: 'belongs_to' },
                  { tone: 'violet', label: 'many_to_many' },
                  { tone: 'amber', label: 'Admin UI picker' },
                ]}
                caption="Declare a relationship — Grit builds the columns, eager loading, and the searchable picker"
              />
            </div>

            <div className="prose-grit">
              {/* belongs_to */}
              <h2>belongs_to</h2>
              <p>
                The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">belongs_to</code> field
                type creates a foreign key relationship. When you add a <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">belongs_to</code> field
                to a resource, the code generator automatically creates:
              </p>
              <ul>
                <li>A foreign key column (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">category_id</code>) with a GORM index</li>
                <li>A GORM association struct field with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">foreignKey</code> tag</li>
                <li><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Preload</code> calls in all handler queries for eager loading</li>
                <li>A searchable relationship select dropdown in admin forms</li>
                <li>Dot notation column display in the DataTable</li>
              </ul>

              <h3>Syntax</h3>
              <p>
                You can either let the generator infer the related model from the field name, or
                specify it explicitly when the field name differs from the model:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                </div>
                <div className="p-5 font-mono text-sm">
                  <div className="mb-1 text-muted-foreground/50 text-xs"># Infer related model from field name</div>
                  <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit generate resource Product --fields &quot;name:string,category:belongs_to,price:float&quot;</span></div>
                  <div className="mt-4 mb-1 text-muted-foreground/50 text-xs"># Explicit related model (when FK name differs)</div>
                  <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit generate resource Post --fields &quot;title:string,author:belongs_to:User,content:text&quot;</span></div>
                </div>
              </div>
            </div>

            <div className="prose-grit">
              <p>
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">category:belongs_to</code> &mdash; infers
                the related model <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Category</code> from
                the field name.
              </p>
              <p>
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">author:belongs_to:User</code> &mdash; explicitly
                sets the related model to <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">User</code>,
                since &quot;author&quot; doesn&apos;t match a model name directly.
              </p>

              <h3>Generated Go Model</h3>
              <p>
                The code generator produces a Go struct with both the foreign key column and the
                association field:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="apps/api/internal/models/product.go" code={`type Product struct {
    ID         uint           \`gorm:"primarykey" json:"id"\`
    Name       string         \`gorm:"size:255" json:"name" binding:"required"\`
    CategoryID uint           \`gorm:"index" json:"category_id" binding:"required"\`
    Category   Category       \`gorm:"foreignKey:CategoryID" json:"category"\`
    Price      float64        \`json:"price"\`
    CreatedAt  time.Time      \`json:"created_at"\`
    UpdatedAt  time.Time      \`json:"updated_at"\`
    DeletedAt  gorm.DeletedAt \`gorm:"index" json:"-"\`
}`} />
            </div>

            <div className="prose-grit">
              <h3>Handler with Preload</h3>
              <p>
                The generated handler uses GORM&apos;s <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Preload</code> to
                automatically eager-load the related model in every query. This means the API
                response always includes the full related object, not just the foreign key ID:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="apps/api/internal/handlers/product_handler.go" code={`// List with eager loading
db.Preload("Category").Find(&products)

// Get by ID
db.Preload("Category").First(&product, id)

// After create/update — reload to include related data in response
db.Preload("Category").First(&product, product.ID)`} />
            </div>

            <div className="prose-grit">
              <h3>Admin Form &mdash; Relationship Select</h3>
              <p>
                The form generates a <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">relationship-select</code> field
                that fetches options from the related resource&apos;s API endpoint:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Relationship select field definition" code={`{
  key: "category_id",
  label: "Category",
  type: "relationship-select",
  required: true,
  relatedEndpoint: "/api/categories",
  displayField: "name",
}`} />
            </div>

            <div className="prose-grit">
              <ul>
                <li>Auto-fetches all categories via React Query</li>
                <li>Searchable dropdown with loading state</li>
                <li>Displays the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">name</code> field (configurable via <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">displayField</code>)</li>
              </ul>

              <h3>DataTable &mdash; Dot Notation</h3>
              <p>
                In the resource definition, the table column uses dot notation to display the
                related model&apos;s name:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Column definition with dot notation" code={`{ key: "category.name", label: "Category", sortable: false }`} />
            </div>

            <div className="prose-grit">
              <p>
                This accesses <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">row.category.name</code> from
                the API response, which includes the Preloaded data. Because the related data comes
                from a join, sorting on dot notation columns is disabled by default.
              </p>
            </div>

            <div className="prose-grit">
              {/* many_to_many */}
              <h2>many_to_many</h2>
              <p>
                The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">many_to_many</code> field
                type creates a junction table relationship. GORM handles the junction table
                automatically &mdash; you don&apos;t need to create or manage it yourself. The
                code generator produces the Go model annotation, association management in
                handlers, and a multi-select component in the admin form.
              </p>

              <h3>Syntax</h3>
              <p>
                For <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">many_to_many</code>,
                the related model is always required (unlike <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">belongs_to</code> where
                it can be inferred):
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                </div>
                <div className="p-5 font-mono text-sm">
                  <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit generate resource Product --fields &quot;name:string,category:belongs_to,tags:many_to_many:Tag,price:float&quot;</span></div>
                </div>
              </div>
            </div>

            <div className="prose-grit">
              <h3>Generated Go Model</h3>
              <p>
                The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">many2many</code> GORM
                tag tells GORM to create and manage the junction table automatically. The table
                name follows the convention <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">model_field</code> (e.g., <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">product_tags</code>):
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="apps/api/internal/models/product.go" code={`type Product struct {
    // ...other fields...
    Tags []Tag \`gorm:"many2many:product_tags" json:"tags"\`
}`} />
            </div>

            <div className="prose-grit">
              <h3>Handler &mdash; Association Management</h3>
              <p>
                Many-to-many associations require special handling in create and update operations.
                The generated handler uses GORM&apos;s <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Association</code> API
                to attach and replace related records by their IDs:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="apps/api/internal/handlers/product_handler.go" code={`// Create — attach tags by IDs
if len(req.TagIDs) > 0 {
    var tags []models.Tag
    h.DB.Where("id IN ?", req.TagIDs).Find(&tags)
    h.DB.Model(&item).Association("Tags").Replace(tags)
}

// Update — replace tags (pointer to detect omission)
if req.TagIDs != nil {
    var tags []models.Tag
    h.DB.Where("id IN ?", *req.TagIDs).Find(&tags)
    h.DB.Model(&item).Association("Tags").Replace(tags)
}`} />
            </div>

            <div className="prose-grit">
              <p>
                The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Replace</code> method
                removes any existing associations and replaces them with the new set. In the update
                handler, a pointer (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">*req.TagIDs</code>)
                is used to distinguish between &quot;not provided&quot; (nil) and &quot;explicitly
                set to empty&quot; (empty slice), enabling partial updates.
              </p>

              <h3>Admin Form &mdash; Multi-Select</h3>
              <p>
                The form generates a <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">multi-relationship-select</code> field
                that allows selecting multiple related records:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Multi-relationship select field definition" code={`{
  key: "tag_ids",
  label: "Tags",
  type: "multi-relationship-select",
  relatedEndpoint: "/api/tags",
  displayField: "name",
  relationshipKey: "tags",
}`} />
            </div>

            <div className="prose-grit">
              <ul>
                <li>Shows removable badge chips for selected items</li>
                <li>Searchable dropdown with multi-select support</li>
                <li><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">relationshipKey</code> maps to the API response field for extracting existing selections in edit mode</li>
              </ul>
            </div>

            <div className="prose-grit">
              {/* one_to_one */}
              <h2 id="one_to_one">one_to_one</h2>
              <p>
                A <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">belongs_to</code> whose foreign key is unique. Declared on the side
                that holds the key, because that is the only side the constraint can live on: a
                passport declares its user, not the other way round.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="bash" code={`grit g resource Passport --fields "user:one_to_one:User,number:string,expires_at:date"`} />
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="apps/api/internal/models/passport.go" code={`type Passport struct {
    UserID string \`gorm:"size:36;uniqueIndex" json:"user_id" binding:"required"\`
    User   User   \`gorm:"foreignKey:UserID" json:"user"\`
    // ...
}`} />
            </div>

            <div className="prose-grit">
              <p>
                <strong>uniqueIndex is the whole difference.</strong> Everything else, the foreign
                key column, the eager loading, the searchable picker in the admin, the CSV import,
                is identical to <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">belongs_to</code>, and that is deliberate: one_to_one is
                the same relationship with a constraint, so it reuses the same machinery rather
                than a parallel implementation.
              </p>
              <p>
                Without the unique index, &quot;one to one&quot; would be a comment. The database
                would accept a second passport pointing at the same user, and nothing would notice
                until somebody asked which one was real. With it, the second insert is refused:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="bash" code={`UNIQUE constraint failed: passports.user_id`} />
            </div>

            <div className="prose-grit">
              <p>
                Worth planning for in the UI: the second attempt fails at the database, so a form
                that lets somebody pick an already-taken parent will surface a constraint error
                rather than a friendly one. Filter the picker, or catch the duplicate and say
                which record already holds it.
              </p>

              {/* has_one & has_many */}
              <h2>has_one &amp; has_many (Inverse Side)</h2>
              <p>
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">has_one</code> and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">has_many</code> are
                the <strong>inverse</strong> of <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">belongs_to</code>.
                They don&apos;t need generator field types because:
              </p>
              <ul>
                <li>The foreign key lives on the <strong>child</strong> model (the one with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">belongs_to</code>)</li>
                <li>When you generate <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Product</code> with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">category:belongs_to</code>, the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Category</code> model automatically <strong>has many</strong> Products via GORM conventions</li>
                <li>You can add the association manually to your parent model if you need to query from the parent side</li>
              </ul>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="apps/api/internal/models/category.go" code={`// Add to your Category model manually
type Category struct {
    // ...existing fields...
    Products []Product \`gorm:"foreignKey:CategoryID" json:"products,omitempty"\`
}`} />
            </div>

            <div className="prose-grit">
              <p>
                This is a manual step &mdash; the generator does not add inverse associations
                automatically, since not every parent model needs to query its children. Add
                the field when you need it, and GORM will handle the rest.
              </p>
            </div>

            <div className="prose-grit">
              {/* Inline items */}
              <h2 id="inline-items">Inline items (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--items</code>)</h2>
              <p>
                <strong>Category / Product</strong> is the default relationship shape: two
                resources, two pages, two forms, two tables, linked by a{" "}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">belongs_to</code>.
                But some pairs &mdash; <strong>Invoice / InvoiceItem</strong>, Order / OrderLine,
                Survey / Question &mdash; want the child created <em>inside</em> the parent&apos;s
                form: you build the invoice and its line items in one go, and they save together
                or not at all.
              </p>
              <p>
                Generate that shape with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--items</code>:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                <div className="p-5 font-mono text-sm">
                  <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit generate resource Invoice \</span></div>
                  <div className="pl-4"><span className="text-foreground/80">--fields &quot;number:string,status:string&quot; \</span></div>
                  <div className="pl-4"><span className="text-foreground/80">--items &quot;InvoiceItem:description:string,qty:int,unit_rate:float&quot;</span></div>
                </div>
              </div>
            </div>

            <div className="prose-grit">
              <p>That one command:</p>
              <ul>
                <li>Generates <strong>InvoiceItem</strong> as a full resource (model, handler, routes) with a <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">invoice:belongs_to:Invoice</code> back-reference &mdash; so it&apos;s filterable by <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">?invoice_id=</code> &mdash; but marked <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">hidden</code>, so it stays out of the sidebar.</li>
                <li>Gives <strong>Invoice</strong> a has-many <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Items []InvoiceItem</code> and a <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">line-items</code> field in its form &mdash; an editable table with add/remove rows and a live per-row and grand total.</li>
                <li>Makes the parent&apos;s <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Create</code>/<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Update</code> handler accept an <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">items</code> array and persist the parent + children in <strong>one GORM transaction</strong> &mdash; atomic, no orphans.</li>
              </ul>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="apps/admin/resources/invoices.ts — the generated line-items field" code={`{
  key: "items",
  label: "Invoice Items",
  type: "line-items",
  colSpan: 2,
  itemEndpoint: "/api/invoice_items",  // child list, for the detail page
  foreignKey: "invoice_id",            // child FK back to the parent
  itemFields: [                        // the editable row columns
    { key: "description", label: "Description", type: "text" },
    { key: "qty",        label: "Qty",        type: "number", numberKind: "int" },
    { key: "unit_rate",  label: "Unit Rate",  type: "number", numberKind: "float" },
  ],
}`} />
            </div>

            <div className="prose-grit">
              <p>
                If a row&apos;s columns include a quantity and a rate/price, the table shows a
                derived <strong>Total</strong> column and a grand total automatically. The
                parent&apos;s <a href="/docs/admin/resources" className="text-primary hover:underline">detail page</a> renders
                the same items as a related table (fetched by the foreign key), so you see them
                after saving without any extra wiring. You can hand-edit this field like any other
                &mdash; add columns, change types, point it at a different child.
              </p>
            </div>

            <div className="prose-grit">
              {/* Hierarchies */}
              <h2 id="hierarchies">Hierarchies (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--tree</code>)</h2>
              <p>
                Everything above relates two <em>different</em> resources. A hierarchy relates
                a resource to <strong>itself</strong>: Electronics contains Cameras contains
                Lenses. That is one table with a parent pointing at another row in the same
                table, and it is the one relationship a plain <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">belongs_to</code> could
                not express, because a Go struct cannot contain itself by value.
              </p>
              <p><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--tree</code> handles it:</p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="bash" code={`grit generate resource Category \\
  --fields "name:string,slug:slug,description:text" \\
  --tree --public`} />
            </div>

            <div className="prose-grit">
              <p>It adds four columns and a service that knows how to use them:</p>
              <ul>
                <li><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">parent_id</code> &mdash; the link upwards, empty for a root</li>
                <li><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">path</code> &mdash; <code>&quot;/id/id/id/&quot;</code>, this row&apos;s id last</li>
                <li><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">depth</code> &mdash; 0 for a root, 1 for its children</li>
                <li><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">position</code> &mdash; the order among siblings</li>
              </ul>
              <p><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">path</code> is the one to understand, because every useful question
                about a hierarchy becomes a string comparison on it. &quot;Everything under
                Electronics&quot; is <code>WHERE path LIKE &apos;/electronics-id/%&apos;</code>:
                one indexed comparison, no recursion, no joins, at any depth. A materialized
                path rather than a recursive CTE because Grit runs on Postgres, MySQL and
                SQLite, and CTE support differs across all three while a path is identical
                everywhere.
              </p>

              <h3>Two levels, and the two questions a category page asks</h3>
              <p>
                Say you have Electronics with Cameras and Laptops under it. A category page
                almost always needs <strong>both</strong> of these, and they have different
                answers:
              </p>
              <ul>
                <li>
                  <strong>Which categories sit under this one?</strong> Those are the tiles you
                  render. Use <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">children</code> from the tree endpoint.
                </li>
                <li>
                  <strong>Which products belong here?</strong> Products are filed under Cameras,
                  not under Electronics, so filtering by the one id returns nothing and the page
                  looks broken while the data is perfect. Use <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">descendant_ids</code> from
                  the detail endpoint.
                </li>
              </ul>
              <p>
                Confusing the two is the usual first bug. <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">descendant_ids</code> is a flat
                list of ids for filtering products; it is not a shape you can render a menu
                from.
              </p>

              <h3>The easy way to fetch both: one call</h3>
              <p><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--tree</code> with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--public</code> mounts an endpoint that
                returns the whole published hierarchy, already nested, in a single query:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="bash" code={`GET /api/v1/public/categories/tree`} />
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="json" filename="response (real, ids trimmed)" code={`{
  "data": [
    {
      "id": "01a01d37-4afc...",
      "parent_id": "",
      "depth": 0,
      "name": "Clothing",
      "slug": "clothing",
      "children": null
    },
    {
      "id": "01a01d37-0e7f...",
      "parent_id": "",
      "depth": 0,
      "name": "Electronics",
      "slug": "electronics",
      "children": [
        {
          "id": "01a01d37-485e...",
          "parent_id": "01a01d37-0e7f...",
          "depth": 1,
          "name": "Cameras",
          "slug": "cameras",
          "children": null
        },
        {
          "id": "01a01d37-49ad...",
          "parent_id": "01a01d37-0e7f...",
          "depth": 1,
          "name": "Laptops",
          "slug": "laptops",
          "children": null
        }
      ]
    }
  ]
}`} />
            </div>

            <div className="prose-grit">
              <p>
                That one response serves the category index page (the roots) and every level-1
                page (each root&apos;s <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">children</code>), so a whole navigation tree
                costs one request. It sits in the public group, which has response caching
                mounted, so it is also among the cheapest things on the page.
              </p>
              <p>
                <strong>
                  A leaf&apos;s <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">children</code> is <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">null</code>, not <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">[]</code>.
                </strong>{" "}
                Go marshals an empty slice as null, so <code>node.children.map(...)</code>
                throws on Cameras. Guard it once, in the helper below, rather than at every
                render site.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="apps/web/hooks/use-categories.ts" code={`export interface CategoryNode {
  id: string
  parent_id: string
  depth: number
  name: string
  slug: string
  description?: string
  /** null on a leaf, not an empty array. */
  children: CategoryNode[] | null
}

/** The whole hierarchy, one request, cached hard because it rarely changes. */
export function useCategoryTree() {
  return useQuery({
    queryKey: ["category-tree"],
    staleTime: 5 * 60 * 1000,
    queryFn: () => get<{ data: CategoryNode[] }>("categories/tree"),
  })
}

/** Depth-first lookup by slug. A shop tree is tens of nodes, not thousands. */
export function findNode(nodes: CategoryNode[], slug: string): CategoryNode | undefined {
  for (const node of nodes) {
    if (node.slug === slug) return node
    const hit = node.children ? findNode(node.children, slug) : undefined
    if (hit) return hit
  }
  return undefined
}

/** Children as an array, whatever the API sent. The null is guarded once, here. */
export function childrenOf(node?: CategoryNode): CategoryNode[] {
  return node?.children ?? []
}`} />
            </div>

            <div className="prose-grit">
              <p>
                The level-1 page then renders its children with no extra request, and the index
                page reads the roots off the same cached response:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="apps/web/app/categories/[slug]/page.tsx" code={`const { data: tree } = useCategoryTree()
const category = findNode(tree?.data ?? [], slug)
const subCategories = childrenOf(category)

return (
  <>
    <h1>{category?.name}</h1>

    {/* Level 2: the tiles. Nothing renders on a leaf, which is correct. */}
    {subCategories.length > 0 && (
      <nav>
        {subCategories.map((child) => (
          <Link key={child.id} href={\`/categories/\${child.slug}\`}>
            {child.name}
          </Link>
        ))}
      </nav>
    )}

    {/* Products in this category AND everything under it. */}
    <ProductGrid slug={slug} />
  </>
)`} />
            </div>

            <div className="prose-grit">
              <p>
                For the products half, the detail endpoint hands back the subtree so you never
                walk the tree yourself:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="bash" code={`GET /api/v1/public/categories/electronics
  -> { "descendant_ids": ["<electronics>", "<cameras>", "<laptops>"] }

GET /api/v1/public/products?category_id=<electronics>,<cameras>,<laptops>
  -> every product in the branch`} />
            </div>

            <div className="prose-grit">
              <p><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">descendant_ids</code> includes the category itself, so the same code
                works unchanged on a leaf. The comma-separated filter is opt-in per column on
                the server and only ever enabled for id columns: splitting on commas is right
                for ids and wrong for anything a person types, where &quot;Smith, John&quot; is
                one value rather than two.
              </p>

              <h3>The rest of the endpoints</h3>
              <ul>
                <li><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">GET /api/v1/categories/tree</code> &mdash; the same tree behind auth, for the admin.</li>
                <li><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">GET /api/v1/categories/:id/breadcrumbs</code> &mdash; ancestors read
                  straight out of the stored path, so it costs one query at any depth.</li>
                <li><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">PATCH /api/v1/categories/:id/move</code> &mdash; reparent, carrying the
                  subtree, with a cycle refused as 422.</li>
                <li><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">POST /api/v1/categories/reorder</code> &mdash; sibling order.</li>
                <li><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">POST /api/v1/categories/rebuild-tree</code> &mdash; recompute every path
                  from <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">parent_id</code> alone. This is what you need after adding <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--tree</code>
                  to a resource that already had rows: those rows have no path, so the tree
                  renders flat until you rebuild.</li>
              </ul>
              <p>
                The admin gets a Tree / Table toggle on the list page: drag onto a row to nest,
                between rows to reorder, onto the bar at the top to promote back to a root.
                Dragging a node into its own subtree is refused before the request is made,
                because a branch moved inside itself detaches from the tree and no query ever
                finds it again.
              </p>
              <p>
                The <Link href="/blog/build-a-storefront-with-grit">storefront guide</Link>{" "}
                builds all of this against a real catalogue in Step 4e.
              </p>
            </div>

            <div className="prose-grit">
              {/* Full Example */}
              <h2>Full Example &mdash; E-Commerce</h2>
              <p>
                Here is a complete workflow that demonstrates both relationship types in an
                e-commerce scenario. Generate the parent models first, then the child model
                with relationships:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                </div>
                <div className="p-5 font-mono text-sm">
                  <div className="mb-1 text-muted-foreground/50 text-xs"># Step 1: Generate Category (the parent)</div>
                  <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit generate resource Category --fields &quot;name:string,slug:slug,description:text&quot;</span></div>
                  <div className="mt-4 mb-1 text-muted-foreground/50 text-xs"># Step 2: Generate Tag</div>
                  <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit generate resource Tag --fields &quot;name:string:unique&quot;</span></div>
                  <div className="mt-4 mb-1 text-muted-foreground/50 text-xs"># Step 3: Generate Product with relationships</div>
                  <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit generate resource Product --fields &quot;name:string,category:belongs_to,tags:many_to_many:Tag,price:float,published:bool&quot;</span></div>
                </div>
              </div>
            </div>

            <div className="prose-grit">
              <p>
                The Product resource definition generated by the commands above includes both
                relationship types in the columns and form fields:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="typescript" filename="apps/admin/resources/products.ts" code={`export default defineResource({
  name: "Product",
  endpoint: "/api/products",
  table: {
    columns: [
      { key: "name", label: "Name", sortable: true, searchable: true },
      { key: "category.name", label: "Category", sortable: false },
      { key: "price", label: "Price", format: "currency", sortable: true },
      { key: "published", label: "Published", format: "boolean" },
    ],
  },
  form: {
    fields: [
      { key: "name", label: "Name", type: "text", required: true },
      {
        key: "category_id",
        label: "Category",
        type: "relationship-select",
        required: true,
        relatedEndpoint: "/api/categories",
        displayField: "name",
      },
      {
        key: "tag_ids",
        label: "Tags",
        type: "multi-relationship-select",
        relatedEndpoint: "/api/tags",
        displayField: "name",
        relationshipKey: "tags",
      },
      { key: "price", label: "Price", type: "number" },
      { key: "published", label: "Published", type: "toggle" },
    ],
  },
})`} />
            </div>

            <div className="prose-grit">
              {/* Customizing Relationships */}
              <h2>Customizing Relationships</h2>
              <p>
                The generated relationship configuration works out of the box, but you can
                customize it to fit your needs. Here are the most common adjustments:
              </p>

              <h3>displayField</h3>
              <p>
                Defaults to <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">&quot;name&quot;</code>.
                Change it to display a different field in the dropdown and table. For example,
                if your related model uses <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">title</code> instead
                of <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">name</code>:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Custom displayField" code={`// Show user email instead of name
{
  key: "author_id",
  label: "Author",
  type: "relationship-select",
  relatedEndpoint: "/api/users",
  displayField: "email",
}

// Show article title
{
  key: "article_id",
  label: "Article",
  type: "relationship-select",
  relatedEndpoint: "/api/articles",
  displayField: "title",
}`} />
            </div>

            <div className="prose-grit">
              <h3>relatedEndpoint</h3>
              <p>
                Auto-generated as <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">/api/&lt;plural&gt;</code>.
                Change it if your API uses a different path or if you need to hit a filtered
                endpoint:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Custom relatedEndpoint" code={`// Custom API path
{
  key: "category_id",
  label: "Category",
  type: "relationship-select",
  relatedEndpoint: "/api/v2/product-categories",
  displayField: "name",
}

// Filtered endpoint — only active users
{
  key: "assignee_id",
  label: "Assignee",
  type: "relationship-select",
  relatedEndpoint: "/api/users?active=true",
  displayField: "name",
}`} />
            </div>

            <div className="prose-grit">
              <h3>Table Display</h3>
              <p>
                The dot notation in column definitions (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">category.name</code>)
                can be changed to access any nested field from the Preloaded response. For example,
                you might want to display a category&apos;s slug instead of its name:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Custom dot notation columns" code={`columns: [
  // Display category slug instead of name
  { key: "category.slug", label: "Category Slug", sortable: false },

  // Display author email
  { key: "author.email", label: "Author Email", sortable: false },
]`} />
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/admin/forms" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Form Builder
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/admin/widgets" className="gap-1.5">
                  Dashboard &amp; Widgets
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
