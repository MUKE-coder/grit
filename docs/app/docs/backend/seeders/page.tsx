import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'
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
                <LaneFlow
                  id="seeders"
                  lanes={['Seed() runner', 'Per-resource seeders', 'Database']}
                  nodes={[
                    { id: 'runner', lane: 0, row: 1, title: 'Seed(db)', sub: 'the runner', tone: 'primary' },
                    { id: 'users', lane: 1, row: 0, title: 'SeedUsers', sub: 'built-in', tone: 'cyan' },
                    { id: 'blogs', lane: 1, row: 1, title: 'SeedBlogs', sub: 'built-in', tone: 'cyan' },
                    { id: 'yours', lane: 1, row: 2, title: 'SeedProduct', sub: 'yours', tone: 'amber' },
                    { id: 'db', lane: 2, row: 1, title: 'PostgreSQL', sub: 'idempotent upserts', tone: 'green' },
                  ]}
                  edges={[
                    { from: 'runner', to: 'users', tone: 'cyan' },
                    { from: 'runner', to: 'blogs', label: 'calls', tone: 'cyan' },
                    { from: 'runner', to: 'yours', tone: 'amber' },
                    { from: 'blogs', to: 'db', label: 'upsert', tone: 'green' },
                  ]}
                  legend={[
                    { tone: 'primary', label: 'Runner' },
                    { tone: 'cyan', label: 'Built-in seeders' },
                    { tone: 'green', label: 'Database' },
                  ]}
                  caption="One runner calls a Seed<Resource> per file — safe to run repeatedly"
                />

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
                  generate a seeder that fills many rows with{' '}
                  <a href="https://github.com/brianvoe/gofakeit" className="text-primary hover:underline">gofakeit</a>,
                  inserted in batches. It ships inside the API, so this works offline.
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
                  <code>gofakeit.Word()</code>. A column marked unique seeds from the
                  row&apos;s number instead (<code>SKU-0000001</code>, <code>SKU-0000002</code>),
                  so it never collides however many rows you ask for. It&apos;s just Go:
                  open the file and swap in your own calls.
                </p>
              </div>

              {/* At scale */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Seeding a million rows
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Give <code>grit seed</code> a resource and a count to top that table up to
                  exactly that many rows:
                </p>

                <CodeBlock terminal code="grit seed Contact --count 1000000" />

                <ul className="list-disc pl-6 space-y-2 text-muted-foreground mb-4">
                  <li>
                    <strong className="text-foreground">It tops up.</strong> The seeder counts
                    what is in the table and inserts only the rest. Run it twice and the
                    second run does nothing. Stop it halfway and the next run carries on from
                    where it stopped, to exactly the count you asked for.
                  </li>
                  <li>
                    <strong className="text-foreground">It inserts in batches.</strong> Rows go
                    in with multi-row inserts, one transaction per batch, sized to what the
                    database accepts for the table&apos;s column count (up to 1,000 rows). Model
                    hooks still run for every row, so IDs, slugs and auto-numbers work as usual.
                  </li>
                  <li>
                    <strong className="text-foreground">Memory stays flat.</strong> Rows are built
                    in chunks by a few goroutines and written as they are ready, never held all
                    at once. A million rows used about 15 MB of heap.
                  </li>
                  <li>
                    <strong className="text-foreground">It fails loudly.</strong> Progress prints
                    every two seconds with rows per second and time left, and the first batch
                    that fails stops the run with an error, instead of logging and reporting
                    success over a half-empty table.
                  </li>
                </ul>

                <p className="text-muted-foreground leading-relaxed mb-4">
                  Measured on one development machine, same resource (name, email, phone and a
                  unique code), before and after this change:
                </p>
                <div className="overflow-x-auto mb-4">
                  <table className="w-full text-sm border border-border rounded-lg">
                    <thead>
                      <tr className="border-b border-border bg-muted/30 text-left">
                        <th className="px-4 py-2 font-medium">Rows</th>
                        <th className="px-4 py-2 font-medium">SQLite, before</th>
                        <th className="px-4 py-2 font-medium">SQLite, now</th>
                        <th className="px-4 py-2 font-medium">Postgres, before</th>
                        <th className="px-4 py-2 font-medium">Postgres, now</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr className="border-b border-border/50"><td className="px-4 py-2">10,000</td><td className="px-4 py-2">96.6 s</td><td className="px-4 py-2">3.7 s</td><td className="px-4 py-2">31.9 s</td><td className="px-4 py-2">2.4 s</td></tr>
                      <tr><td className="px-4 py-2">1,000,000</td><td className="px-4 py-2">about 2 h 45 min</td><td className="px-4 py-2">3 min 14 s</td><td className="px-4 py-2">about 53 min</td><td className="px-4 py-2">about 50 s</td></tr>
                    </tbody>
                  </table>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The before figures for a million rows are estimated from the measured rate;
                  the after figures are measured. SQLite takes one writer at a time, so it is
                  slower than Postgres, which writes three batches at once.
                </p>

                <div className="rounded-lg border border-primary/20 bg-primary/5 p-4">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-foreground">A seeder from before --count.</strong>{' '}
                    Seeders generated by an older Grit still run with <code>grit seed</code>,
                    but not with <code>--count</code>. Run <code>grit upgrade</code>, then
                    regenerate the seeder with <code>grit generate seeder Contact --faker</code>{' '}
                    to switch it over.
                  </p>
                </div>
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
                  code={`func SeedProductsTo(db *gorm.DB, target int64) error {
    // Link each row to an existing parent (loaded once).
    var categoryIDs []string
    db.Model(&models.Category{}).Pluck("id", &categoryIDs)

    return SeedTopUp(db, "products", SeedPlan[models.Product]{
        Target: target,
        Make: func(n int64) models.Product {
            return models.Product{
                Name:       gofakeit.Name(),
                Price:      gofakeit.Price(1, 1000),
                CategoryID: pickID(categoryIDs), // a real, existing category
            }
        },
    })
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
                      <tr className="border-b border-border/50"><td className="px-4 py-2 font-mono text-[13px]">grit seed X --count N</td><td className="px-4 py-2 text-muted-foreground">Top X up to N rows, in batches; resumes if stopped</td></tr>
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
