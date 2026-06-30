import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
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
  columns: [
    { key: "name", label: "Name", sortable: true, searchable: true },
    { key: "category.name", label: "Category", sortable: false },
    { key: "price", label: "Price", format: "number", sortable: true },
    { key: "published", label: "Published", format: "boolean" },
  ],
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
