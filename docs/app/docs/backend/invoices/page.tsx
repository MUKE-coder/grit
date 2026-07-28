import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/invoices')

export default function InvoicesPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <p className="tag-mono text-primary/80 mb-3">Backend</p>
            <h1 className="text-4xl font-bold tracking-tight mb-4">Invoices &amp; Line Items</h1>
            <p className="text-lg text-muted-foreground leading-relaxed mb-6">
              An invoice is the canonical parent-with-children resource: an{' '}
              <code>Invoice</code> that owns many <code>InvoiceItem</code> rows. This guide covers
              generating both in one command, generating them separately, auto-numbering the
              invoice, and printing it.
            </p>
            <div className="mb-10 rounded-lg border border-border bg-muted/30 p-4 text-sm text-muted-foreground leading-relaxed">
              <strong className="text-foreground">Invoice is just the example.</strong> Every
              technique here is generic. The same <code>--items</code> shape models any
              parent-with-children &mdash; orders / order-items, purchase-orders / lines,
              surveys / questions, playlists / tracks. <code>grit generate sequence</code> numbers
              any resource (orders, receipts, tickets), and the print view is on every generated
              detail page. Read &ldquo;Invoice&rdquo; as &ldquo;whatever you&apos;re modeling.&rdquo;
            </div>

            {/* D1 — combined --items */}
            <h2 id="one-command" className="text-2xl font-semibold mb-4">
              Generate both in one command
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              <code>--items</code> generates the child resource <em>and</em> wires it to the
              parent as an inline, editable line-items table inside the parent&apos;s form. Rows
              are saved atomically with the invoice (GORM has-many in one transaction).
            </p>
            <CodeBlock
              terminal
              code={`grit g resource Invoice --fields "number:string,status:string" \\
  --items "InvoiceItem:description:string,qty:int,unit_rate:float"`}
            />
            <p className="text-muted-foreground leading-relaxed mb-4 mt-6">
              Read the two arguments like this:
            </p>
            <CodeBlock
              language="text"
              filename="anatomy of --items"
              code={`--fields "number:string,status:string"
          └──────────────┬──────────────┘
             the PARENT (Invoice) columns

--items "InvoiceItem : description:string,qty:int,unit_rate:float"
          └────┬─────┘   └──────────────────┬─────────────────────┘
          child model            child fields (same name:type
             name              grammar as --fields), comma-separated

# Result:
#   Invoice   → number, status,  + a line-items field holding InvoiceItems
#   InvoiceItem → description, qty, unit_rate, + invoice_id (belongs_to Invoice)
#
# The child gets a belongs_to back to the parent automatically; the parent's
# form renders the children as an add/remove table saved with the invoice.`}
            />

            {/* D2 — separate */}
            <h2 id="separately" className="text-2xl font-semibold mb-4 mt-12">
              Generate them separately
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Prefer to build the pieces one at a time? <code>--items</code> is just a shortcut
              for two <code>generate resource</code> calls plus a relationship. Do it by hand:
            </p>
            <CodeBlock
              terminal
              code={`# 1. the parent
grit g resource Invoice --fields "number:string,status:string"

# 2. the child, with a belongs_to back to the parent
grit g resource InvoiceItem \\
  --fields "description:string,qty:int,unit_rate:float,invoice:belongs_to:Invoice"`}
            />
            <p className="text-muted-foreground leading-relaxed mb-4 mt-6">
              That gives you two independent, fully-routed resources linked by{' '}
              <code>invoice_id</code>. The difference from <code>--items</code>: the child is its
              own top-level resource (its own page, list, and endpoints) rather than an inline
              table on the invoice form. Add the inline table later by adding a{' '}
              <code>line-items</code> field to the invoice&apos;s resource definition, or forget a
              column and add it with{' '}
              <Link href="/docs/backend/migrations" className="text-primary hover:underline">
                grit g field
              </Link>
              :
            </p>
            <CodeBlock terminal code={`grit g field Invoice due_date:date`} />

            {/* D3 — auto number */}
            <h2 id="auto-number" className="text-2xl font-semibold mb-4 mt-12">
              Auto-number the invoice
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              You rarely want users typing invoice numbers by hand. The one-liner is the{' '}
              <code>auto</code> field modifier — declare it on the field and Grit wires up the
              whole atomic-counter machinery for you:
            </p>
            <CodeBlock
              terminal
              code={`grit g resource Invoice --fields "number:string:auto:INV,status:string,total:float"`}
            />
            <p className="text-muted-foreground leading-relaxed mb-4 mt-6">
              <code>number:string:auto:INV</code> reads as{' '}
              <em>&ldquo;a string column named <code>number</code>, auto-generated with the prefix
              INV.&rdquo;</em> The prefix is optional (<code>number:string:auto</code> derives one
              from the model name). That single modifier does four things so you don&apos;t have to:
            </p>
            <ul className="list-disc pl-6 space-y-2 text-muted-foreground leading-relaxed mb-4">
              <li>
                stands up the shared <code>internal/sequence</code> package and registers its
                counter table with AutoMigrate (once per project);
              </li>
              <li>
                generates the model&apos;s <code>BeforeCreate</code> hook to fill the field from an
                atomic, gap-free counter — <code>INV-202607-0001</code>, <code>-0002</code>, …;
              </li>
              <li>
                makes the column <strong>optional</strong> (the server fills it, so nothing is
                required at the API boundary);
              </li>
              <li>
                <strong>hides it from the create/edit form</strong> — no empty &ldquo;Number&rdquo;
                box for users to puzzle over — while keeping it on the table and detail page.
              </li>
            </ul>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Create an invoice with no number and it comes back numbered; import one that already
              has a number and that number is kept (the hook only fills blanks). By default{' '}
              <code>auto</code> resets the counter monthly with a 4-digit width — want yearly, never,
              or a different width? Reach for the explicit <code>grit generate sequence</code> route
              below, which exposes all three knobs.
            </p>

            <h3 className="text-xl font-semibold mb-3 mt-8">
              The explicit route: <code>grit generate sequence</code>
            </h3>
            <p className="text-muted-foreground leading-relaxed mb-4">
              <code>auto</code> is a shortcut over this command. Use it directly when you want to
              control the reset cadence or width, number an <em>existing</em> field, or call the
              counter from your own handler code. It creates the same{' '}
              <strong>atomic, gap-free</strong> counter plus a typed helper you can call from the
              create path.
            </p>
            <CodeBlock
              terminal
              code={`grit generate sequence Invoice --prefix INV --reset monthly --width 4`}
            />
            <p className="text-muted-foreground leading-relaxed mb-4 mt-6">
              That writes <code>internal/sequence/</code> (a generic counter package) and{' '}
              <code>internal/services/invoice_sequence.go</code> exposing{' '}
              <code>NextInvoiceNumber(db, t)</code>. <code>--reset</code> controls when the
              counter rolls over (<code>monthly</code>, <code>yearly</code>, <code>never</code>)
              and <code>--width</code> the zero-padding — so the numbers look like:
            </p>
            <CodeBlock language="text" filename="pattern" code={`INV-202607-0001
INV-202607-0002
INV-202608-0001   ← monthly reset rolls the counter over`} />
            <p className="text-muted-foreground leading-relaxed mb-4 mt-6">
              Wire it into the invoice&apos;s <code>BeforeCreate</code> hook so every new invoice
              is numbered automatically — set it only when blank, so an imported invoice keeps its
              original number. Call the generic <code>sequence.Next</code> directly: the{' '}
              <code>services.NextInvoiceNumber</code> wrapper lives in the <code>services</code>{' '}
              package (which imports <code>models</code>), so calling it from a model would be an
              import cycle. The <code>sequence</code> package imports no models, so a model can call
              it — use the wrapper from <em>handlers</em> instead.
            </p>
            <CodeBlock
              language="go"
              filename="internal/models/invoice.go"
              code={`import "yourapp/apps/api/internal/sequence"

func (m *Invoice) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	if m.Number == "" {
		number, err := sequence.Next(tx, sequence.Config{
			Name: "invoice", Prefix: "INV", Reset: sequence.ResetMonthly, Width: 4,
		}, time.Now())
		if err != nil {
			return err
		}
		m.Number = number
	}
	return nil
}`}
            />
            <p className="text-muted-foreground leading-relaxed mb-4 mt-6">
              Because the counter lives in a row that&apos;s locked and incremented in the same
              transaction as the insert, this is safe under concurrent load — no two invoices ever
              get the same number, and there are no gaps. Want a different shape entirely (say{' '}
              <code>2026/Q3/0001</code>)? The helper is plain Go you can edit; the sequence package
              just hands you the next integer.
            </p>
            <p className="text-muted-foreground leading-relaxed mb-4 mt-6">
              <strong>The number fills in on the server, not in the form.</strong> The hook runs at
              create time and only when the field is blank — nothing pre-fills the browser. When you
              wire the sequence by hand, declare the field <code>number:string:optional</code>{' '}
              (string fields are required by default) so the create form won&apos;t demand it, and
              drop the <code>number</code> entry from the resource&apos;s <code>form.fields</code> if
              you&apos;d rather not show an empty box; it stays in the table and detail either way.
              The <code>auto</code> modifier above does both of these for you — optional column,
              hidden from the form — which is exactly why it exists.
            </p>

            {/* line-item totals */}
            <h2 id="line-item-totals" className="text-2xl font-semibold mb-4 mt-12">
              Line-item totals
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              The inline line-items table shows a per-row <strong>Total</strong> and a grand total
              as you type — driven by <em>column names</em>, not configuration. It looks for one
              column matching a quantity pattern (<code>qty</code> / <code>quantity</code>) and one
              matching a money pattern (<code>unit_rate</code>, <code>unit_price</code>,{' '}
              <code>rate</code>, <code>price</code>, or <code>amount</code>); if both exist it
              renders <code>Total = quantity × money</code> and sums the rows. The default{' '}
              <code>qty:int</code> + <code>unit_rate:float</code> match, so you get it for free.
            </p>
            <p className="text-muted-foreground leading-relaxed mb-4">
              It&apos;s <strong>display-only</strong> — computed in the browser to help data entry;
              only the declared columns are submitted, so no total is stored unless you add an{' '}
              <code>amount</code> column and compute it in the item&apos;s <code>BeforeSave</code>
              hook. And because the trigger is the column name, renaming <code>qty</code> to{' '}
              <code>count</code> or <code>unit_rate</code> to <code>cost</code> simply stops the
              auto-total from showing (nothing breaks) — keep a name in those patterns to keep it.
            </p>

            {/* D4 — PDF */}
            <h2 id="pdf" className="text-2xl font-semibold mb-4 mt-12">
              The PDF endpoint
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Every generated resource exposes a server-rendered PDF, and the detail page has a{' '}
              <strong>PDF</strong> button that opens it:
            </p>
            <CodeBlock terminal code={`GET /api/invoices/:id/pdf   →   application/pdf`} />
            <p className="text-muted-foreground leading-relaxed mb-4 mt-6">
              It is rendered in <strong>Go</strong>, not by the browser, which is the important
              part: the same bytes come back for everyone, so you can attach the PDF to an email,
              push it to S3, or hand it to a background job — none of which is possible with a
              print dialog. The layout is built from the record itself:
            </p>
            <ul className="list-disc pl-6 space-y-2 text-muted-foreground leading-relaxed mb-4">
              <li>
                a <strong>repeating header</strong> (your <code>APP_NAME</code> plus the
                record&apos;s identifier) and a <strong>repeating footer</strong> with{' '}
                <strong>Page N of M</strong> — on every page, so a long invoice stays navigable;
              </li>
              <li>
                a title block, then the resource&apos;s fields as a <strong>two-up grid</strong>{' '}
                (dates formatted, booleans as Yes/No, empty values as an em dash, and a{' '}
                <code>belongs_to</code> shown by the related record&apos;s name rather than its
                UUID);
              </li>
              <li>
                each set of <strong>line items as a table</strong> with numeric columns
                right-aligned, then totals and notes.
              </li>
            </ul>
            <p className="text-muted-foreground leading-relaxed mb-4">
              The handler is ordinary generated Go in{' '}
              <code>internal/handlers/invoice.go</code> — it builds a <code>pdf.Record</code> and
              hands it to <code>pdf.RenderRecord</code>. Restyle by editing that struct (reorder
              fields, add a total, change the title); go deeper by editing{' '}
              <code>internal/pdf/record.go</code>, or drop to the underlying{' '}
              <code>go-pdf/fpdf</code> document for full control:
            </p>
            <CodeBlock
              language="go"
              filename="internal/handlers/invoice.go (generated — yours to edit)"
              code={`rec := pdf.Record{
	Title:      "INVOICE",
	Subtitle:   pdf.Value(item.Number),
	Brand:      appName,
	FooterNote: appName + " · generated " + time.Now().Format("2 Jan 2006 15:04"),
	Fields: []pdf.Field{
		{Label: "Number",   Value: pdf.Value(item.Number)},
		{Label: "Status",   Value: pdf.Value(item.Status)},
		{Label: "Customer", Value: pdf.Display(item.Customer)},
	},
}

// Line items become a table section, with totals underneath.
rec.Sections = append(rec.Sections, pdf.Section{
	Title:   "Invoice Items",
	Headers: []string{"Description", "Qty", "Unit Rate"},
	Aligns:  []string{"L", "R", "R"},
	Rows:    itemRows,
})
rec.Totals = []pdf.TotalLine{{Label: "Total", Value: "6,000,000", Bold: true}}

out, err := pdf.RenderRecord(rec)`}
            />

            {/* D5 — print */}
            <h2 id="printing" className="text-2xl font-semibold mb-4 mt-12">
              Printing an invoice
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              The <strong>Print</strong> button next to it calls the browser&apos;s print dialog
              against a print-optimized layout: a print stylesheet hides everything except the
              record&apos;s detail card and its line-items table, so the sidebar, navbar, and the
              Edit/Delete/Back controls never reach the paper. Related resources are excluded too.
              The stylesheet sets <code>@page</code> margins, forces ink-friendly colors, repeats
              table headers across pages, and avoids splitting a row down the middle. Reach for{' '}
              <strong>PDF</strong> when you need a file to keep or send; reach for{' '}
              <strong>Print</strong> when you just want paper.
            </p>
            <p className="text-muted-foreground leading-relaxed mb-4">
              It works the moment the resource exists, with no per-resource code. Under the hood
              the detail content is wrapped in <code>#print-area</code> and the{' '}
              <code>@media print</code> block in the admin&apos;s <code>globals.css</code> makes
              only that subtree visible — so a custom invoice template is just a matter of styling
              that one container.
            </p>

            {/* nav */}
            <div className="flex items-center justify-between mt-16 pt-8 border-t border-border">
              <Button variant="ghost" asChild>
                <Link href="/docs/backend/migrations">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  Migrations
                </Link>
              </Button>
              <Button variant="ghost" asChild>
                <Link href="/docs/backend/seeders">
                  Seeders
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
