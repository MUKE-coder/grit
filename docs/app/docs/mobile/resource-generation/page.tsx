import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/mobile/resource-generation')

export default function MobileResourceGenerationPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Mobile</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Resource Generation
              </h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                When you run <code>grit generate resource</code> in a mobile project, Grit emits
                six Expo files per resource on top of the Go backend and shared types — a typed
                hook, four screens, and a shared form. This page is a deep dive into each one and
                how field types shape what gets generated.
              </p>
              <LaneFlow
                id="mob-resgen"
                lanes={['grit generate', 'Expo files per resource']}
                nodes={[
                  { id: 'cmd', lane: 0, row: 2, title: 'grit generate', sub: 'resource Product', tone: 'primary' },
                  { id: 'hook', lane: 1, row: 0, title: 'use-products', sub: 'typed hook', tone: 'cyan' },
                  { id: 'list', lane: 1, row: 1, title: 'List screen', sub: 'index', tone: 'blue' },
                  { id: 'detail', lane: 1, row: 2, title: 'Detail screen', sub: '[id]', tone: 'blue' },
                  { id: 'form', lane: 1, row: 3, title: 'Create / Edit', sub: 'shared form', tone: 'green' },
                  { id: 'more', lane: 1, row: 4, title: '+ 2 screens', sub: 'edit · etc.', tone: 'amber' },
                ]}
                edges={[
                  { from: 'cmd', to: 'hook', tone: 'cyan' },
                  { from: 'cmd', to: 'list', tone: 'blue' },
                  { from: 'cmd', to: 'detail', label: 'writes', tone: 'blue' },
                  { from: 'cmd', to: 'form', tone: 'green' },
                  { from: 'cmd', to: 'more', tone: 'amber' },
                ]}
                legend={[{ tone: 'primary', label: 'One command' }, { tone: 'cyan', label: 'Six Expo files' }]}
                caption="One command writes a typed hook, four screens, and a shared form — six Expo files per resource"
              />
            </div>

            {/* Fan-out diagram */}
            <div className="prose-grit mb-4">
              <h2>The resource → files fan-out</h2>
              <p>
                A single command produces a working native CRUD flow. Grit detects the Expo app
                automatically (it looks for <code>apps/expo</code>) and generates the mobile
                files alongside everything else:
              </p>
            </div>
            <CodeBlock language="text" filename="grit generate resource Product" code={`grit generate resource Product --fields "name:string,price:float,category:belongs_to"
   │
   ├─ Go backend        → model + service + handler + routes
   ├─ packages/shared   → schemas/product.ts + types/product.ts
   │
   └─ apps/expo/        (6 files)
      ├─ hooks/use-products.ts                       ← typed React Query hook
      ├─ app/products/index.tsx                      ← list screen
      ├─ app/products/[id].tsx                       ← detail screen
      ├─ app/products/new.tsx                        ← create screen
      ├─ app/products/edit/[id].tsx                  ← edit screen
      └─ components/resource-forms/products-form.tsx ← shared create/edit form
                        │
                        └─ + a "Products" card injected into
                           app/(tabs)/explore.tsx  (// grit:mobile-resources)`} className="mb-10" />

            {/* Files table */}
            <div className="prose-grit mb-4">
              <h2>The six generated files</h2>
            </div>
            <div className="overflow-x-auto mb-10">
              <table className="w-full text-sm border border-border rounded-lg">
                <thead>
                  <tr className="border-b border-border bg-muted/50">
                    <th className="text-left p-3 font-medium">File</th>
                    <th className="text-left p-3 font-medium">What it is</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border text-muted-foreground">
                  {[
                    ['hooks/use-<plural>.ts', 'Typed React Query hook: infinite-scroll list, single-record query, and create / update / delete mutations. Exports the resource interface + a Page type.'],
                    ['app/<plural>/index.tsx', 'List screen — search, sortable columns, pull-to-refresh, infinite scroll, CSV export, bulk import, and a create sheet.'],
                    ['app/<plural>/[id].tsx', 'Detail screen — one labelled Row per field, optional hero image, plus Edit and Delete actions.'],
                    ['app/<plural>/new.tsx', 'Create screen — renders the shared form and calls the create mutation.'],
                    ['app/<plural>/edit/[id].tsx', 'Edit screen — loads the record, pre-fills the shared form, calls the update mutation.'],
                    ['components/resource-forms/<plural>-form.tsx', 'Shared form rendered by both new and edit (and the list’s create sheet).'],
                  ].map(([file, desc]) => (
                    <tr key={file}>
                      <td className="p-3 font-mono text-xs text-primary align-top whitespace-nowrap">{file}</td>
                      <td className="p-3">{desc}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* The hook */}
            <div className="prose-grit mb-4">
              <h2>The typed hook</h2>
              <p>
                <code>use-&lt;plural&gt;.ts</code> is the single source of data access for the
                resource. It declares the resource&apos;s TypeScript interface (IDs are UUID
                strings; <code>created_at</code> / <code>updated_at</code> are always present)
                and exports five hooks:
              </p>
              <ul>
                <li><code>use&lt;Plural&gt;()</code> — an infinite-scroll list built on <code>useInfiniteQuery</code>. It accepts a search string, equality filters, sort key and order, and a page size. It accumulates pages and stops when the last page is reached.</li>
                <li><code>use&lt;Singular&gt;(id)</code> — fetches one record with <code>useQuery</code>.</li>
                <li><code>useCreate&lt;Singular&gt;()</code>, <code>useUpdate&lt;Singular&gt;()</code>, <code>useDelete&lt;Singular&gt;()</code> — mutations that invalidate the list cache on success.</li>
              </ul>
              <p>
                The filters argument is what powers <code>belongs_to</code> scoping — e.g.{' '}
                <code>useProducts(&quot;&quot;, {'{'} category_id: id {'}'})</code> lists only the
                products in one category. The API already supports the matching query parameter.
              </p>
            </div>

            {/* The form */}
            <div className="prose-grit mb-4">
              <h2>The shared form component</h2>
              <p>
                <code>&lt;plural&gt;-form.tsx</code> is a single component used by the create
                screen, the edit screen, and the list&apos;s create sheet. It pre-fills from an
                optional <code>initial</code> record and calls <code>onSubmit(values)</code> — so
                the parent owns the mutation and the navigation. Each field type maps to a native
                control:
              </p>
            </div>
            <div className="overflow-x-auto mb-10">
              <table className="w-full text-sm border border-border rounded-lg">
                <thead>
                  <tr className="border-b border-border bg-muted/50">
                    <th className="text-left p-3 font-medium">Field type</th>
                    <th className="text-left p-3 font-medium">Rendered control</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border text-muted-foreground">
                  {[
                    ['string, date, datetime', 'Single-line TextInput'],
                    ['text, richtext', 'Multi-line TextInput (textarea)'],
                    ['int, uint, float', 'Numeric TextInput with format/parse helpers'],
                    ['bool', 'Native Switch toggle'],
                    ['belongs_to', 'RelationSelect chip picker (see below)'],
                    ['file', 'Single native image upload with live preview'],
                    ['files', 'Multi-image grid with removable thumbnails'],
                    ['slug, many_to_many, string_array', 'Skipped in the form (slug is derived server-side)'],
                  ].map(([type, control]) => (
                    <tr key={type}>
                      <td className="p-3 font-mono text-xs text-primary align-top whitespace-nowrap">{type}</td>
                      <td className="p-3">{control}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* belongs_to */}
            <div className="prose-grit mb-4">
              <h2>How belongs_to fields render</h2>
              <p>
                A <code>belongs_to</code> field pulls in two things. In the <strong>form</strong>,
                it renders a <code>RelationSelect</code> — a chip picker whose options come from
                the related resource&apos;s own generated hook (loaded with a large first page so
                the picker isn&apos;t capped at the list default). Submitting sends the selected
                foreign key, and the field is validated as required.
              </p>
              <p>
                On the <strong>list screen</strong>, the same relations become filter pills in a
                filter sheet — tap a parent to scope the list to its children. The detail screen
                shows the parent&apos;s <code>name</code> or <code>title</code> (preloaded by the
                API) instead of the raw foreign key.
              </p>
            </div>
            <CodeBlock language="typescript" filename="products-form.tsx (belongs_to excerpt)" code={`const categoriesQuery = useCategories("", {}, "created_at", "desc", 500);
const categoriesOpts = categoriesQuery.data?.pages.flatMap((p) => p.data) ?? [];

// …in the form body:
<RelationSelect
  label="Category"
  value={categoryId}
  onChange={setCategoryId}
  options={categoriesOpts}
/>`} className="mb-10" />

            {/* file fields */}
            <div className="prose-grit mb-4">
              <h2>How file fields render</h2>
              <p>
                A <code>file</code> or <code>files</code> field renders a native image picker.
                Tapping it opens an <code>ImagePickerSheet</code>; picked images show{' '}
                <em>instantly</em> as a local preview, upload in the background via{' '}
                <code>uploadLocalFile()</code>, and the returned stored URLs go into the payload
                as <code>FileRef</code> objects. A single <code>file</code> is one framed
                preview; <code>files</code> is a grid of removable thumbnails with an add tile.
              </p>
              <div className="rounded-lg border border-primary/20 bg-primary/5 p-4 mt-4">
                <p className="text-sm text-muted-foreground leading-relaxed">
                  <strong className="text-foreground">Images that just work.</strong> Uploaded
                  images are stored through the API&apos;s S3-compatible storage. The scaffold
                  ships <code className="text-xs font-mono">lib/images.ts</code> with a{' '}
                  <code className="text-xs font-mono">resolveImageUrl()</code> helper that
                  rewrites the host so stored images load on a device (a phone can&apos;t reach{' '}
                  <code className="text-xs font-mono">localhost</code>) — list rows, detail hero
                  images, and form previews all run through it.
                </p>
              </div>
            </div>

            {/* Idempotency + re-run */}
            <div className="prose-grit mb-4 mt-10">
              <h2>Re-running is safe</h2>
              <p>
                Generation is idempotent where it matters. The resource card injected into the
                More tab is added at the <code>// grit:mobile-resources</code> marker and is
                skipped if it already exists, so regenerating a resource never duplicates the
                link. Shared helpers (the screen header, form sheet, pickers, upload and image
                helpers) are only written if they&apos;re missing — older scaffolds get them
                backfilled on the next generate.
              </p>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/mobile/first-app" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Your First Mobile App
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/mobile/building" className="gap-1.5">
                  Building &amp; Publishing
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
