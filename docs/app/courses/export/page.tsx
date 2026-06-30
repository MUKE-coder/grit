import Link from 'next/link'
import { FileSpreadsheet, Download, Gauge, Wand2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'CSV / Excel Export — Per-Resource Data Export in Grit',
  description:
    'Every generated Grit resource ships a CSV and XLSX export endpoint. Learn streaming CSV, chunked XLSX, constant-memory exports, and how filters carry through to the file.',
}

const learn = [
  { icon: Wand2, title: 'Generated for free', body: 'Every resource from grit generate gets /export with no extra code.' },
  { icon: FileSpreadsheet, title: 'CSV & XLSX', body: 'Pick the format with ?format=csv or ?format=xlsx.' },
  { icon: Gauge, title: 'Constant memory', body: 'Streaming CSV and chunked XLSX export millions of rows without blowing up RAM.' },
  { icon: Download, title: 'Filters carry through', body: 'The same query params that filter the list filter the export.' },
]

export default function ExportCourse() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">CSV / Excel Export</span>
        </div>

        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">4 challenges</span>
          </div>
          <h1 className="font-display text-3xl md:text-4xl font-bold text-foreground mb-4">
            CSV / Excel Export
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Data export is the feature every business app eventually needs and nobody enjoys
            building. In Grit it&apos;s automatic: every resource you scaffold with
            <Code>grit generate resource</Code> ships a streaming <strong className="text-foreground">CSV</strong> and
            chunked <strong className="text-foreground">XLSX</strong> export at <Code>/export</Code> —
            built for constant memory.
          </p>
        </div>

        <div className="grid sm:grid-cols-2 rounded-xl border border-foreground/15 overflow-hidden mb-12">
          {learn.map(({ icon: Icon, title, body }) => (
            <div key={title} className="border-b border-r border-foreground/15 p-5">
              <Icon className="h-5 w-5 text-primary mb-3" />
              <h3 className="font-semibold text-foreground text-sm mb-1">{title}</h3>
              <p className="text-sm text-muted-foreground leading-relaxed">{body}</p>
            </div>
          ))}
        </div>

        {/* 1 */}
        <section className="mb-12">
          <h2 className="font-display text-2xl font-bold text-foreground mb-4">The export endpoint</h2>
          <p className="text-muted-foreground leading-relaxed mb-4">
            When you generate a resource, Grit injects an export route alongside the usual CRUD
            handlers. You choose the format with a query parameter.
          </p>
          <CodeBlock filename="exporting a Product resource">
{`# CSV (default) — streamed row-by-row
GET /api/products/export?format=csv

# Excel workbook — written in chunks
GET /api/products/export?format=xlsx

# Filters from the list endpoint apply to the export too
GET /api/products/export?format=xlsx&status=active&sort=created_at`}
          </CodeBlock>
          <Note>
            Because export reuses the list query, whatever the user is looking at on screen is
            exactly what lands in the file — same filters, same sort, same scope.
          </Note>
        </section>

        {/* 2 */}
        <section className="mb-12">
          <h2 className="font-display text-2xl font-bold text-foreground mb-4">Why streaming matters</h2>
          <p className="text-muted-foreground leading-relaxed mb-4">
            The naive approach — load every row into a slice, build the whole file in memory,
            then send it — falls over on large tables. A 1-million-row export can eat gigabytes
            of RAM and time out. Grit avoids this by writing as it reads.
          </p>
          <Definition term="Streaming (CSV)">
            Rows are fetched in batches and written straight to the HTTP response as they arrive.
            Memory stays flat regardless of row count because no full copy of the data is ever held.
          </Definition>
          <Definition term="Chunked (XLSX)">
            The XLSX format is a zipped XML package, so it can&apos;t be streamed as plainly as CSV.
            Grit writes rows in chunks using a streaming writer, keeping peak memory bounded even
            for very large workbooks.
          </Definition>
          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Format</th>
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Best for</th>
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Memory profile</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-medium text-foreground">CSV</td><td className="px-3 py-2">Huge dumps, piping into other tools</td><td className="px-3 py-2">Flat (true stream)</td></tr>
                <tr><td className="px-3 py-2 font-medium text-foreground">XLSX</td><td className="px-3 py-2">Hand-off to non-technical users</td><td className="px-3 py-2">Bounded (chunked)</td></tr>
              </tbody>
            </table>
          </div>
          <Tip>
            For exports that take a while, dispatch the export as a background job and email a
            download link when it&apos;s ready — Grit&apos;s jobs + email batteries make that a few lines.
          </Tip>
        </section>

        {/* 3 */}
        <section className="mb-12">
          <h2 className="font-display text-2xl font-bold text-foreground mb-4">From the admin panel</h2>
          <p className="text-muted-foreground leading-relaxed mb-4">
            The admin DataTable includes an export menu wired to these endpoints, so end users
            can download the current view as CSV or Excel with one click — no API knowledge needed.
          </p>
        </section>

        {/* challenges */}
        <section className="mb-12">
          <h2 className="font-display text-2xl font-bold text-foreground mb-4">Practice</h2>
          <Challenge number={1} title="Generate and export">
            <p>Run <Code>grit generate resource Product --fields &quot;name:string,price:float,stock:int&quot;</Code>,
            seed a few rows, then hit <Code>/api/products/export?format=csv</Code>. Open the file —
            do the columns match your fields?</p>
          </Challenge>
          <Challenge number={2} title="Switch to Excel">
            <p>Request the same data with <Code>format=xlsx</Code>. Does it open cleanly in Excel,
            Numbers, or Google Sheets? What&apos;s different from the CSV?</p>
          </Challenge>
          <Challenge number={3} title="Carry a filter through">
            <p>Add a filter query param (e.g. <Code>?status=active</Code>) to the export URL. Does
            the file contain only the filtered rows? Why is reusing the list query the right design?</p>
          </Challenge>
          <Challenge number={4} title="Reason about scale">
            <p>You need to export 5 million rows. Which format do you choose and why? What would
            you change so the user isn&apos;t left waiting on an open HTTP connection?</p>
          </Challenge>
        </section>

        <CourseFooter />

        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses', label: 'All Courses' }}
            next={{ href: '/courses/batteries', label: 'Batteries Included' }}
          />
        </div>
      </main>
    </div>
  )
}
