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
              You rarely want users typing invoice numbers by hand. <code>grit generate
              sequence</code> creates an <strong>atomic, gap-free</strong> counter backed by a
              database row (so two concurrent creates can never collide) and a typed helper you
              call from the create path.
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
              create time and only when the field is blank — nothing pre-fills the browser. So
              declare the field <code>number:string:optional</code> (string fields are required by
              default) and the create form won&apos;t demand it; the value appears on the detail
              page and list right after you save. Prefer not to show an empty Number box at all?
              Remove the <code>number</code> entry from the resource&apos;s <code>form.fields</code>{' '}
              — it stays in the table and detail, just not the create form.
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

            {/* D4 — print */}
            <h2 id="printing" className="text-2xl font-semibold mb-4 mt-12">
              Printing an invoice
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Every generated resource detail page ships with a <strong>Print</strong> button (open
              a row from the table with its <em>View</em> action, then hit Print). It calls the
              browser&apos;s print dialog against a print-optimized layout: a print stylesheet
              hides everything except the record&apos;s detail card and its line-items table, so
              the sidebar, navbar, and the Edit/Delete/Back controls never reach the paper. Related
              resources are excluded too — the printout is just the invoice and its items.
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
