import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/seeders')

export default function SeedersPage() {
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
                Seeders
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Seeders fill your database with starter data &mdash; the admin account,
                demo users, sample catalogue rows, anything you want on a fresh
                install. In Grit, <strong>every resource gets its own seeder file</strong>,
                you can generate one in a single command, and <code>--faker</code> fills
                it with realistic rows (relationships included).
              </p>
            </div>

            <div className="prose-grit">
              {/* Mental model / diagram */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  How it fits together
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  There is one thin <code>Seed()</code> runner that calls a{' '}
                  <code>Seed&lt;Resource&gt;</code> function per resource. Each of those
                  lives in its own file under <code>internal/database/</code>, so a
                  seeder is always easy to find and edit &mdash; including the built-in
                  users and blogs.
                </p>

                <CodeBlock
                  language="text"
                  filename="apps/api/internal/database/"
                  code={`  seed.go                 ← Seed(db): the runner
   │
   ├─ SeedUsers(db)        → users_seeder.go     (admin + demo users)
   ├─ SeedBlogs(db)        → blogs_seeder.go     (sample posts)
   ├─ SeedCategories(db)   → categories_seeder.go
   └─ SeedProducts(db)     → products_seeder.go
                              ▲
                              └─ grit generate seeder / --seed adds these
   seed_helpers.go         ← pickID / firstID (relationship helpers)`}
                />

                <p className="text-muted-foreground leading-relaxed mt-4">
                  When you generate a seeder, Grit writes the{' '}
                  <code>&lt;resource&gt;_seeder.go</code> file <em>and</em> registers
                  its call in <code>seed.go</code> at the <code>// grit:seeders</code>{' '}
                  marker. You never wire anything by hand.
                </p>
              </div>

              {/* Running */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Running seeders
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  After migrating, run every seeder with one command from anywhere in
                  the project:
                </p>

                <CodeBlock terminal code="grit seed" />

                <p className="text-muted-foreground leading-relaxed mb-4">
                  Seeders are <strong>idempotent</strong> &mdash; each checks whether its
                  table already has rows and skips if so, so re-running never
                  duplicates data.
                </p>

                <CodeBlock
                  language="bash"
                  filename="output"
                  code={`Seeding database...
Created admin user: admin@example.com / admin123
Created user: jane@example.com / admin123
Created blog: "Getting Started with Grit" (published)
Seeded 8 category
Seeded 60 product
Database seeded successfully.`}
                />
              </div>

              {/* Generating */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Generating a seeder
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Add a seeder to a resource you already generated &mdash; it reads the
                  model to pre-fill one example record with the right field types:
                </p>

                <CodeBlock terminal code="grit generate seeder Customer" />

                <p className="text-muted-foreground leading-relaxed mb-4">
                  Pass more than one, or emit the seeder at the same time you scaffold
                  the resource with <code>--seed</code>:
                </p>

                <CodeBlock
                  terminal
                  code={`grit generate seeder Customer Order Product

grit generate resource Tag --fields "name:string" --seed`}
                />
              </div>

              {/* Faker */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Filling rows with faker
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Without a flag you get <strong>one editable example row</strong>. Add{' '}
                  <code>--faker</code> (and <code>--count N</code>, default 10) to instead
                  generate a loop that fills many rows with{' '}
                  <a href="https://github.com/brianvoe/gofakeit" className="text-primary hover:underline">gofakeit</a>.
                  It ships inside the API, so this works offline.
                </p>

                <CodeBlock terminal code="grit generate seeder Product --faker --count 60" />

                <p className="text-muted-foreground leading-relaxed mb-4">
                  Values are chosen from each field&apos;s <strong>name and type</strong>:
                </p>

                <div className="overflow-x-auto mb-4">
                  <table className="w-full text-sm border border-border rounded-lg">
                    <thead>
                      <tr className="border-b border-border bg-muted/30 text-left">
                        <th className="px-4 py-2 font-medium">Field</th>
                        <th className="px-4 py-2 font-medium">Faker value</th>
                      </tr>
                    </thead>
                    <tbody className="font-mono text-[13px]">
                      <tr className="border-b border-border/50"><td className="px-4 py-2">name</td><td className="px-4 py-2">gofakeit.Name()</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2">email</td><td className="px-4 py-2">gofakeit.Email()</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2">phone / city / company</td><td className="px-4 py-2">gofakeit.Phone() / City() / Company()</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2">float (price)</td><td className="px-4 py-2">gofakeit.Price(1, 1000)</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2">int / uint</td><td className="px-4 py-2">gofakeit.Number(1, 100)</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2">bool</td><td className="px-4 py-2">gofakeit.Bool()</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2">date / datetime</td><td className="px-4 py-2">gofakeit.Date()</td></tr>
                      <tr><td className="px-4 py-2">file:image / files:image</td><td className="px-4 py-2">a sample picsum image URL</td></tr>
                    </tbody>
                  </table>
                </div>
                <p className="text-muted-foreground leading-relaxed">
                  Anything the guesser doesn&apos;t recognise falls back to{' '}
                  <code>gofakeit.Word()</code>. It&apos;s just Go &mdash; open the file
                  and swap in your own calls.
                </p>
              </div>

              {/* Relationships */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Relationships
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  This is the part most seeders get wrong. A{' '}
                  <code>belongs_to</code> field (a Product&apos;s Category, say) needs a{' '}
                  <em>real</em> parent id, not a random string. Grit handles it: the
                  seeder loads the parent ids once and links each row to one of them
                  &mdash; a random parent for faker, the first parent for the static
                  example.
                </p>

                <CodeBlock
                  language="go"
                  filename="products_seeder.go (faker)"
                  code={`func SeedProducts(db *gorm.DB) error {
    // ... skip if already seeded ...

    // Link each row to an existing parent (loaded once).
    var categoryIDs []string
    db.Model(&models.Category{}).Pluck("id", &categoryIDs)

    const n = 60
    for i := 0; i < n; i++ {
        r := models.Product{
            Name:       gofakeit.Name(),
            Price:      gofakeit.Price(1, 1000),
            CategoryID: pickID(categoryIDs), // ← a real, existing category
        }
        db.Create(&r)
    }
    return nil
}`}
                />

                <div className="rounded-lg border border-primary/20 bg-primary/5 p-4 mt-4">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-foreground">Seed order matters.</strong> A
                    child can only link to a parent that already exists, so seed{' '}
                    <strong>parents first</strong>. The runner calls seeders in the order
                    you generated the resources &mdash; generate <code>Category</code>{' '}
                    before <code>Product</code> and you&apos;re set. Need a different
                    order? Reorder the calls in <code>seed.go</code>.
                  </p>
                </div>
              </div>

              {/* Editing */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Editing a seeder
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  A static seeder is a plain slice of model structs &mdash; edit the
                  values, add rows, done:
                </p>

                <CodeBlock
                  language="go"
                  filename="categories_seeder.go"
                  code={`func SeedCategories(db *gorm.DB) error {
    var count int64
    db.Model(&models.Category{}).Count(&count)
    if count > 0 {
        return nil // already seeded
    }

    records := []models.Category{
        {Name: "Sample Name"},          // ← edit these
        // {Name: "Phones"},            // ← or add your own
        // {Name: "Accessories"},
    }

    for _, r := range records {
        db.Create(&r)
    }
    return nil
}`}
                />
              </div>

              {/* Commands recap */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Command reference
                </h2>
                <div className="overflow-x-auto">
                  <table className="w-full text-sm border border-border rounded-lg">
                    <thead>
                      <tr className="border-b border-border bg-muted/30 text-left">
                        <th className="px-4 py-2 font-medium">Command</th>
                        <th className="px-4 py-2 font-medium">Does</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr className="border-b border-border/50"><td className="px-4 py-2 font-mono text-[13px]">grit seed</td><td className="px-4 py-2 text-muted-foreground">Run every seeder</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2 font-mono text-[13px]">grit generate seeder X [Y…]</td><td className="px-4 py-2 text-muted-foreground">Add a seeder to existing resource(s)</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2 font-mono text-[13px]">grit generate resource X … --seed</td><td className="px-4 py-2 text-muted-foreground">Emit the seeder while scaffolding</td></tr>
                      <tr><td className="px-4 py-2 font-mono text-[13px]">… --faker --count N</td><td className="px-4 py-2 text-muted-foreground">Fill N rows with gofakeit instead of one example</td></tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Prev / Next */}
              <div className="flex items-center justify-between border-t border-border pt-8 mt-12">
                <Button variant="ghost" asChild>
                  <Link href="/docs/backend/migrations" className="gap-2">
                    <ArrowLeft className="h-4 w-4" />
                    Migrations
                  </Link>
                </Button>
                <Button variant="ghost" asChild>
                  <Link href="/docs/backend/models" className="gap-2">
                    Models
                    <ArrowRight className="h-4 w-4" />
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
