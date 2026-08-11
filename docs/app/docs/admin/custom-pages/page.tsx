import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/admin/custom-pages')

/**
 * The porting page. Written for the person who has bought an admin template and
 * wants its pages inside Grit — which is a different job from "customising the
 * admin", and the docs previously only answered the second one.
 */
export default function CustomPagesPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Admin Panel</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Custom pages and tables
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Bring your own table, your own form, your own page shell — and keep the URL-synced
                sorting, paging, filters, selection, bulk delete and toasts that the stock page
                already has.
              </p>
            </div>

            <div className="prose-grit mb-10">
              <h2 id="the-short-version">The short version</h2>
              <p>
                A generated resource page is three lines, and you are allowed to replace them:
              </p>
            </div>

            <CodeBlock
              language="tsx"
              filename="apps/admin/app/(dashboard)/resources/products/page.tsx"
              code={`"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { productsResource } from "@/resources/products";

export default function ProductsPage() {
  return <ResourcePage resource={productsResource} />;
}`}
            />

            <div className="prose-grit mb-10">
              <p>
                The reason people did not replace it was everything they would lose by doing so. The
                data was never the problem — <code>useResource</code> has always been a plain hook
                that takes an endpoint. The problem was the rest of the page: search, sort, page and
                filters kept in the address bar so a refresh or a shared link rehydrates the same
                view, row selection, bulk delete behind a confirm, toasts, cache invalidation, and
                stat cards that follow the active date range instead of contradicting the table
                under them.
              </p>
              <p>
                That is now a hook. <code>useResourceController</code> returns all of it and renders
                nothing.
              </p>
            </div>

            <CodeBlock
              language="tsx"
              filename="apps/admin/app/(dashboard)/resources/products/page.tsx"
              code={`"use client";

import { useResourceController } from "@/hooks/use-resource-controller";
import { productsResource } from "@/resources/products";
import { TemplateShell, TemplateTable, TemplatePager } from "@/components/template";

export default function ProductsPage() {
  const c = useResourceController(productsResource);

  return (
    <TemplateShell title={c.pluralName} onAdd={c.create}>
      <TemplateTable
        rows={c.rows}
        columns={c.columns}
        loading={c.isLoading}
        sortKey={c.sortBy}
        sortDir={c.sortOrder}
        onSort={c.setSort}
        selected={c.selection}
        onSelect={c.setSelection}
        onRowClick={c.edit}
      />
      <TemplatePager
        page={c.page}
        pages={c.totalPages}
        total={c.total}
        onChange={c.setPage}
      />
    </TemplateShell>
  );
}`}
            />

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-10">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">
                Why this is safe
              </h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                The stock <code>ResourcePage</code> is built on the same hook and contains no state
                of its own. If the controller could not rebuild the default page, it would be
                missing something — so anything the default page can do, yours can too.
              </p>
            </div>

            <div className="prose-grit mb-10">
              <h2 id="what-you-get">What the controller returns</h2>
              <p>
                Every setter that changes the query resets to page one, because a search run from
                page seven should not land on an empty page seven of two results.{' '}
                <code>setSort</code> toggles direction when you pass the same key twice.
              </p>
            </div>

            <CodeBlock
              language="ts"
              filename="the shape"
              code={`const c = useResourceController<Product>(productsResource)

// data
c.rows          // Product[]
c.meta          // { total, page, page_size, pages } | undefined
c.total         // number
c.totalPages    // number
c.isLoading     // boolean

// query state — sort/page/filters, and the date range round-trips through the URL
c.page          c.setPage(n)
c.pageSize      c.setPageSize(n)
c.search        c.setSearch(value)
c.sortBy        c.sortOrder      c.setSort(key)    // toggles direction
c.filters       c.setFilter(key, value)
c.dateRange     c.setDateRange(range)

// columns — visible ones, ready to render
c.columns       c.allColumns     c.hiddenColumns    c.toggleColumn(key)

// selection
c.selection     c.setSelection(ids)                 c.clearSelection()

// actions
c.actions       c.can("create" | "view" | "edit" | "delete")
c.create()      c.edit(row)      c.view(row)
c.remove(id)    // opens the confirm dialog
c.bulkRemove()  // opens the bulk confirm dialog
c.isDeleting    c.isBulkDeleting

// dialog state, if you render your own
c.form              // { open, item, close() }
c.confirmDelete     // { open, confirm(), cancel() }
c.confirmBulkDelete // { open, confirm(), cancel() }
c.importer          // { open, setOpen(open) }

// odds and ends
c.apiSearchParams   // the same query the table ran — use it for exports
c.stats             // stat cards, already scoped to the active date range
c.singularName      c.pluralName`}
            />

            <div className="prose-grit mb-10">
              <h2 id="one-piece-at-a-time">Replacing one piece at a time</h2>
              <p>
                You do not have to take the whole page. Keep the stock layout and swap only the
                table, because <code>c.columns</code> and <code>c.rows</code> are ordinary values:
              </p>
            </div>

            <CodeBlock
              language="tsx"
              filename="keeping the toolbar, replacing the table"
              code={`const c = useResourceController(productsResource);

return (
  <div>
    <PageHeader title={c.pluralName} stats={c.stats} />

    <TableToolbar
      resource={productsResource}
      search={c.search}
      onSearch={c.setSearch}
      selectedCount={c.selection.length}
      onBulkDelete={c.bulkRemove}
      allColumns={c.allColumns}
      hiddenColumns={c.hiddenColumns}
      onToggleColumn={c.toggleColumn}
      data={c.rows}
      dateRange={c.dateRange}
      onDateRangeChange={c.setDateRange}
      apiSearchParams={c.apiSearchParams}
    />

    {/* your table, Grit's everything else */}
    <TemplateTable rows={c.rows} onSort={c.setSort} />
  </div>
);`}
            />

            <div className="prose-grit mb-10">
              <h2 id="the-custom-file">Registering it once: the .custom.tsx file</h2>
              <p>
                Editing the route file works, but it only customises that one route. The detail
                page, a relationship picker and anything else rendering the resource still get the
                stock components. To set it once and have every route pick it up, use the
                customisation file that sits next to the resource:
              </p>
            </div>

            <CodeBlock
              language="bash"
              filename="apps/admin/resources/"
              code={`products.ts          # generated — rewritten on every grit generate
products.custom.tsx  # yours — created once, never touched again`}
            />

            <div className="prose-grit mb-10">
              <p>
                The split is what makes both halves safe. The config half can be regenerated freely
                because nothing of yours is in it. The custom half can hold components because it is
                a <code>.tsx</code> file and the generator will not overwrite it — it checks whether
                the file exists and leaves it alone if it does.
              </p>
            </div>

            <CodeBlock
              language="tsx"
              filename="apps/admin/resources/products.custom.tsx"
              code={`import type { ResourceCustomisation } from "@/lib/resource";
import { DataTable } from "@/components/tables/data-table";
import { StatusPill, TemplateTable } from "@/components/template";

const custom: ResourceCustomisation = {
  // 1. Override a single cell, keep everything else
  columns: {
    status: { cell: (row) => <StatusPill value={String(row.status)} /> },
    price: { cell: (row) => <b>{"$" + Number(row.price).toFixed(2)}</b> },
  },

  components: {
    // 2. Replace the table. Same props DataTable takes, so this is a drop-in —
    //    header, toolbar, filters and pagination all keep working.
    Table: (props) => <TemplateTable rows={props.data} onSort={props.onSort} />,

    // 3. Or wrap the original instead of replacing it
    // Table: (props) => <TemplateCard><DataTable {...props} /></TemplateCard>,

    // 4. Replace the whole page — call useResourceController inside it
    // Page: MyProductsPage,
  },
};

export default custom;`}
            />

            <div className="prose-grit mb-10">
              <p>
                <code>columns</code> and <code>fields</code> are patched <strong>by key</strong>,
                not replaced wholesale. That is deliberate: <code>grit sync</code> keeps adding new
                columns as you add fields to the Go model, and your renderers survive it. A key that
                does not match any generated column is simply ignored.
              </p>

              <h3 id="the-slots">The slots</h3>
              <ul>
                <li>
                  <code>Table</code> — receives exactly <code>DataTable</code>&apos;s props:{' '}
                  <code>columns</code>, <code>data</code>, <code>isLoading</code>,{' '}
                  <code>sortBy</code>, <code>sortOrder</code>, <code>onSort</code>,{' '}
                  <code>selectedRows</code>, <code>onSelectRows</code>, <code>onView</code>,{' '}
                  <code>onEdit</code>, <code>onDelete</code>, <code>rowActions</code>.
                </li>
                <li>
                  <code>Form</code> — receives <code>resource</code>, <code>item</code> (the record
                  being edited, or <code>null</code> for create) and <code>onClose</code>. Replaces
                  whichever container <code>formView</code> would have opened.
                </li>
                <li>
                  <code>EmptyState</code> — rendered instead of the table when the query has
                  finished and returned nothing.
                </li>
                <li>
                  <code>Page</code> — replaces the entire list view. Checked before anything else,
                  so a page slot owns its own routing.
                </li>
              </ul>

              <h3 id="wrapping">Wrapping instead of replacing</h3>
              <p>
                Because a slot receives the stock component&apos;s own props, you can render the
                original inside yours. That is the cheap way to restyle a shell or add something
                around a table without reimplementing sorting and selection:
              </p>
            </div>

            <CodeBlock
              language="tsx"
              filename="wrapping the default"
              code={`import { DataTable } from "@/components/tables/data-table";

components: {
  Table: (props) => (
    <div className="rounded-2xl border border-dashed p-2">
      <p className="mb-2 text-xs text-muted-foreground">
        {props.data.length} rows on this page
      </p>
      <DataTable {...props} />
    </div>
  ),
}`}
            />

            <div className="prose-grit mb-10">
              <h2 id="pages-that-are-not-resources">Pages that are not resources</h2>
              <p>
                Porting a whole template means analytics, settings and billing screens that are not
                CRUD over a table. Those do not need the controller at all — use the data hooks
                directly against any endpoint your API exposes:
              </p>
            </div>

            <CodeBlock
              language="tsx"
              filename="a page with no resource behind it"
              code={`import { useResource } from "@/hooks/use-resource";

export default function RevenuePage() {
  const { data, isLoading } = useResource<Invoice>("/api/invoices", {
    pageSize: 100,
    filters: { status: "paid" },
  });

  if (isLoading) return <TemplateSkeleton />;

  return <TemplateRevenueChart rows={data?.data ?? []} />;
}`}
            />

            <div className="prose-grit mb-10">
              <h2 id="regeneration">Will the generator overwrite this?</h2>
              <p>
                No. <code>grit generate resource</code> writes{' '}
                <code>resources/&lt;name&gt;.ts</code> and the thin page wrapper. Once you have
                replaced the wrapper with your own component, re-running the generator for a{' '}
                <em>new</em> resource does not touch it. What the generator does keep maintaining is
                the resource definition — and <code>grit sync</code> only ever inserts into it,
                between the <code>grit:cols:auto-start</code> and <code>grit:fields:auto-start</code>{' '}
                fences, so hand-edited labels and formats survive.
              </p>
            </div>

            <div className="flex items-center justify-between border-t border-border/40 pt-6">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/admin/datatable" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  DataTable
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/admin/forms" className="gap-1.5">
                  Form Builder
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
