import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import { LaneFlow } from '@/components/lane-flow'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Build an Invoice Generator: Desktop App with Wails — Grit Course',
  description: 'Build a real desktop invoice generator with Grit + Wails. Create clients, invoices with line items, generate PDF invoices, track payment status, and export to Excel.',
}

export default function InvoiceDesktopCourse() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Invoice Generator Desktop App</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">14 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Build an Invoice Generator: Desktop App with Wails
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            In this course, you will build a real, useful desktop application — an invoice generator.
            Create clients, generate invoices with line items, export professional PDF invoices, and
            track payment status. Everything runs offline with SQLite, no internet required.
          </p>
          <LaneFlow id="c-invoice" lanes={['Invoice app', 'Modules', 'Offline']}
            nodes={[
              { id: 'app', lane: 0, row: 1, title: 'Invoice app', sub: 'native desktop', tone: 'primary' },
              { id: 'clients', lane: 1, row: 0, title: 'Clients', sub: 'CRM', tone: 'cyan' },
              { id: 'invoices', lane: 1, row: 1, title: 'Invoices', sub: 'line items', tone: 'blue' },
              { id: 'pdf', lane: 1, row: 2, title: 'PDF export', sub: 'professional', tone: 'green' },
              { id: 'sqlite', lane: 2, row: 1, title: 'SQLite', sub: 'offline, no internet', tone: 'amber' },
            ]}
            edges={[{ from: 'app', to: 'clients', tone: 'cyan' }, { from: 'app', to: 'invoices', label: 'build', tone: 'blue' }, { from: 'invoices', to: 'pdf', tone: 'green' }, { from: 'app', to: 'sqlite', label: 'stores', tone: 'amber' }]}
            legend={[{ tone: 'primary', label: 'Native app' }, { tone: 'amber', label: 'Fully offline' }]}
            caption="Clients, invoices with line items, and PDF export — all running offline on SQLite" />
        </div>

        <div className="my-4 rounded-lg border border-primary/20 bg-primary/5 px-4 py-3">
          <p className="text-sm text-muted-foreground">
            <strong className="text-foreground">Reference docs:</strong>{' '}
            <a href="/docs/desktop/first-app" className="text-primary hover:underline">Desktop Apps →</a>
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What We're Building ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What We{"'"}re Building</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            An invoice generator is one of the most practical apps a freelancer or small business can
            use. Instead of paying $20/month for a SaaS tool, you build your own — and it runs entirely
            on your computer with no subscription:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Client management</strong> — store client names, emails, addresses, companies</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Invoice creation</strong> — auto-numbered invoices with line items, tax calculation</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">PDF export</strong> — professional invoices ready to email to clients</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Payment tracking</strong> — Draft, Sent, Paid, Overdue status workflow</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Excel export</strong> — export all invoices for accounting and tax filing</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Dashboard</strong> — total revenue, outstanding balance, overdue count</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Everything is stored locally in SQLite. No server, no cloud, no account needed. Open the
            app and start invoicing.
          </p>

          <Challenge number={1} title="List Invoice Features">
            <p>List 5 features an invoice generator needs beyond what{"'"}s listed above. Think about: recurring invoices, multiple currencies, late fees, payment reminders, client portal. Which would you prioritize?</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: Scaffold ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Scaffold the App</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Use Grit{"'"}s desktop scaffold command to create the project:
          </p>

          <CodeBlock filename="Terminal">
{`grit new-desktop invoicer`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This creates a Wails v2 project with Go backend, React + TypeScript frontend, local auth,
            SQLite database, and export capabilities (PDF + Excel). Start the development server:
          </p>

          <CodeBlock filename="Terminal">
{`cd invoicer && grit start`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The app opens as a native window. You should see the login page — create an account and
            explore the default dashboard. This is your starting point.
          </p>

          <Challenge number={2} title="Scaffold and Run">
            <p>Scaffold the project and run <Code>grit start</Code>. Create an account and log in. Explore the default pages — dashboard, blog, contacts. These are the starter templates you will customize into an invoice app.</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Design the Data Model ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Design the Data Model</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The invoice app needs three resources. Notice how they relate to each other: a Client has
            many Invoices, and an Invoice has many LineItems.
          </p>

          <CodeBlock filename="Terminal">
{`# Resource 1: Client
grit generate resource Client \\
  name:string email:string phone:string:optional \\
  address:text:optional company:string:optional

# Resource 2: Invoice
grit generate resource Invoice \\
  invoice_number:string:unique \\
  client_id:belongs_to:Client \\
  issue_date:date due_date:date status:string \\
  subtotal:float tax:float total:float \\
  notes:text:optional

# Resource 3: LineItem
grit generate resource LineItem \\
  description:string quantity:int \\
  unit_price:float total:float \\
  invoice_id:belongs_to:Invoice`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The relationships work like this:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Client</strong> — standalone entity. Has name, email, optional phone/address/company.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Invoice</strong> — belongs to a Client. Has dates, status, financial totals.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">LineItem</strong> — belongs to an Invoice. Each line is one product/service with qty and price.</li>
          </ul>

          <Challenge number={3} title="Generate All 3 Resources">
            <p>Generate all three resources. Check the Go models — are the <Code>belongs_to</Code> foreign keys correct? Can you see the relationships in the model structs?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Auto-Generate Invoice Numbers ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Auto-Generate Invoice Numbers</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every invoice needs a unique, sequential number. Instead of making the user type one,
            auto-generate it. Add this method to the Invoice service:
          </p>

          <CodeBlock filename="internal/services/invoice_service.go">
{`func (s *InvoiceService) NextNumber() string {
    var count int64
    s.db.Model(&models.Invoice{}).Count(&count)
    return fmt.Sprintf("INV-%03d", count+1)
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This counts existing invoices and generates the next number: <Code>INV-001</Code>,
            <Code> INV-002</Code>, <Code>INV-003</Code>, and so on. The <Code>%03d</Code> format pads
            with leading zeros up to 3 digits.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Call this method when creating a new invoice. The Wails binding exposes it to the frontend,
            so you can pre-fill the invoice number field automatically.
          </p>

          <Tip>
            For production apps with concurrent users, you{"'"}d want a database sequence or a transaction
            lock to avoid duplicate numbers. For a single-user desktop app, counting rows is perfectly safe.
          </Tip>

          <Challenge number={4} title="Add Auto-Numbering">
            <p>Add the <Code>NextNumber()</Code> method to your invoice service. Create 3 invoices — are they numbered <Code>INV-001</Code>, <Code>INV-002</Code>, <Code>INV-003</Code>? Delete the second one and create another — what number does it get?</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Build the Client List ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Build the Client List</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The client list is your first custom page. It shows all clients in a searchable, sortable
            DataTable — the same component Grit{"'"}s desktop scaffold provides:
          </p>

          <CodeBlock filename="Frontend — Client List Page">
{`// The DataTable component handles:
// - Column sorting (click headers)
// - Search/filter
// - Pagination
// - Row selection
// - Actions (edit, delete)

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'company', label: 'Company' },
  { key: 'email', label: 'Email' },
  { key: 'phone', label: 'Phone' },
]

// Fetch clients from the Go backend via Wails binding
const clients = await GetAllClients()`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The client list is the foundation. Every invoice needs a client, so build this page first.
            Add a {'"'}New Client{'"'} button that opens the FormBuilder with fields for name, email,
            company, phone, and address.
          </p>

          <Challenge number={5} title="Create 5 Clients">
            <p>Build the client list page and the create client form. Add 5 clients with different companies. Can you search by name? Sort by company? Edit a client{"'"}s email?</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Build the Invoice Form ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Build the Invoice Form</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The invoice form is the most complex page in the app. It has multiple steps and dynamic
            line items that calculate totals automatically:
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">Step 1: Invoice Details</strong>
          </p>
          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> Select client from dropdown (populated from your client list)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Invoice number (auto-filled, but editable)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Issue date (defaults to today)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Due date (defaults to 30 days from today)</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">Step 2: Line Items</strong>
          </p>
          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> Dynamic rows: description, quantity, unit price</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Line total auto-calculated: quantity x unit price</li>
            <li className="flex gap-2"><span className="text-primary">•</span> {'"'}Add Line Item{'"'} button to add more rows</li>
            <li className="flex gap-2"><span className="text-primary">•</span> {'"'}Remove{'"'} button on each row</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">Step 3: Review</strong>
          </p>
          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> Subtotal (sum of all line totals)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Tax rate input (percentage, e.g., 10%)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Tax amount (subtotal x tax rate)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Grand total (subtotal + tax)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Optional notes field</li>
          </ul>

          <CodeBlock filename="Line Item Calculation Logic">
{`// Each line item: total = quantity * unitPrice
const lineTotal = quantity * unitPrice

// Invoice totals
const subtotal = lineItems.reduce((sum, item) => sum + item.total, 0)
const taxAmount = subtotal * (taxRate / 100)
const grandTotal = subtotal + taxAmount`}
          </CodeBlock>

          <Challenge number={6} title="Create an Invoice with 4 Line Items">
            <p>Build the invoice form. Create an invoice with 4 line items (e.g., {'"'}Web Design — 40 hours at $75{'"'}, {'"'}Logo Design — 1 at $500{'"'}, etc.). Verify: Are the line totals correct? Is the subtotal the sum of all lines? Does the tax calculate properly?</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Invoice Detail View ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Invoice Detail View</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The invoice detail page shows a complete, formatted view of an invoice. This is what the
            user sees when they click on an invoice from the list:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Header</strong> — invoice number, status badge, issue date, due date</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Client info</strong> — name, company, email, address</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Line items table</strong> — description, quantity, unit price, line total</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Totals section</strong> — subtotal, tax, grand total</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Status badge</strong> — color-coded: Draft (gray), Sent (blue), Paid (green), Overdue (red)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Action buttons</strong> — Export PDF, Mark as Sent, Mark as Paid, Edit, Delete</li>
          </ul>

          <CodeBlock filename="Status Badge Colors">
{`const statusColors = {
  draft:   'bg-gray-500/10 text-gray-400 border-gray-500/20',
  sent:    'bg-blue-500/10 text-blue-400 border-blue-500/20',
  paid:    'bg-green-500/10 text-green-400 border-green-500/20',
  overdue: 'bg-red-500/10 text-red-400 border-red-500/20',
}`}
          </CodeBlock>

          <Challenge number={7} title="View an Invoice">
            <p>Build the invoice detail page. Navigate to one of your invoices. Is all the information displayed correctly — client name, line items, totals? Does the status badge use the right color?</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: PDF Export ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">PDF Export</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The killer feature of an invoice app: generating a professional PDF that you can email to
            clients. Grit{"'"}s desktop scaffold includes an export service that generates PDFs from Go:
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The PDF should include:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Company header</strong> — your business name, address, logo (optional)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Invoice number and dates</strong> — prominently displayed</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Bill To section</strong> — client name, company, address</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Line items table</strong> — with headers, alternating row colors</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Totals</strong> — subtotal, tax, grand total (bold)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Payment terms</strong> — {'"'}Payment due within 30 days{'"'}</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Notes</strong> — any additional notes from the invoice</li>
          </ul>

          <CodeBlock filename="Export Button (Frontend)">
{`// Call the Go export function via Wails binding
const handleExportPDF = async () => {
  try {
    const filePath = await ExportInvoicePDF(invoice.id)
    // filePath is where the PDF was saved
    alert("PDF saved to: " + filePath)
  } catch (error) {
    console.error("Export failed:", error)
  }
}`}
          </CodeBlock>

          <Tip>
            The Go backend handles PDF generation using a library like <Code>gofpdf</Code> or
            <Code> go-wkhtmltopdf</Code>. The Wails binding makes it callable from the React frontend
            as a simple async function. No REST API needed — it{"'"}s a direct function call.
          </Tip>

          <Challenge number={8} title="Export a PDF Invoice">
            <p>Add a {'"'}Export PDF{'"'} button to the invoice detail page. Generate a PDF for one of your invoices. Open the PDF — does it look professional? Is all the data correct: client info, line items, totals?</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Mark as Paid ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Payment Status Workflow</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Invoices follow a status workflow. Each transition represents a real business action:
          </p>

          <CodeBlock filename="Status Workflow">
{`Draft → Sent → Paid
  ↓              ↑
  └──→ Overdue ──┘

Draft:   Invoice created but not sent to client
Sent:    Invoice emailed/delivered to client
Paid:    Client has paid the invoice
Overdue: Past due date and not yet paid`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The status transitions should be explicit actions — buttons on the invoice detail page:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">{'"'}Mark as Sent{'"'}</strong> — visible when status is Draft</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">{'"'}Mark as Paid{'"'}</strong> — visible when status is Sent or Overdue</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Auto-Overdue</strong> — a background check: if status is Sent and due_date has passed, mark as Overdue</li>
          </ul>

          <CodeBlock filename="internal/services/invoice_service.go">
{`func (s *InvoiceService) UpdateStatus(id uint, status string) error {
    validTransitions := map[string][]string{
        "draft":   {"sent"},
        "sent":    {"paid", "overdue"},
        "overdue": {"paid"},
    }

    var invoice models.Invoice
    if err := s.db.First(&invoice, id).Error; err != nil {
        return err
    }

    allowed := validTransitions[invoice.Status]
    for _, s := range allowed {
        if s == status {
            return s.db.Model(&invoice).Update("status", status).Error
        }
    }
    return fmt.Errorf("cannot transition from %s to %s", invoice.Status, status)
}`}
          </CodeBlock>

          <Challenge number={9} title="Test the Status Workflow">
            <p>Create an invoice. It should start as Draft. Mark it as Sent — does the badge change to blue? Mark it as Paid — does the badge change to green? Try marking a Paid invoice as Draft — does the validation prevent it?</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Dashboard ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Dashboard</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The dashboard gives you a financial overview at a glance:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Total Revenue</strong> — sum of all Paid invoices</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Outstanding</strong> — sum of Sent invoices (money owed to you)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Overdue</strong> — count and total of overdue invoices (needs attention)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Recent Invoices</strong> — last 5 invoices with status, client, and amount</li>
          </ul>

          <CodeBlock filename="Go Service — Dashboard Stats">
{`func (s *InvoiceService) DashboardStats() (*DashboardStats, error) {
    var stats DashboardStats

    // Total revenue (paid invoices)
    s.db.Model(&models.Invoice{}).
        Where("status = ?", "paid").
        Select("COALESCE(SUM(total), 0)").
        Scan(&stats.TotalRevenue)

    // Outstanding (sent invoices)
    s.db.Model(&models.Invoice{}).
        Where("status = ?", "sent").
        Select("COALESCE(SUM(total), 0)").
        Scan(&stats.Outstanding)

    // Overdue count
    s.db.Model(&models.Invoice{}).
        Where("status = ?", "overdue").
        Count(&stats.OverdueCount)

    return &stats, nil
}`}
          </CodeBlock>

          <Challenge number={10} title="Verify Dashboard Accuracy">
            <p>Create 10+ invoices across all statuses (3 Draft, 3 Sent, 2 Paid, 2 Overdue). Check the dashboard — does Total Revenue only count Paid invoices? Does Outstanding only count Sent? Is the Overdue count correct?</p>
          </Challenge>
        </section>

        {/* ═══ Section 11: Excel Export ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Excel Export</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For accounting and tax purposes, you need to export all invoices to a spreadsheet.
            Grit{"'"}s desktop scaffold includes Excel export via the <Code>excelize</Code> Go library:
          </p>

          <CodeBlock filename="Excel Export Columns">
{`Invoice Number | Client | Issue Date | Due Date | Status | Subtotal | Tax | Total
INV-001        | Acme   | 2026-01-15 | 2026-02-14 | Paid | 3000.00 | 300.00 | 3300.00
INV-002        | Globex | 2026-01-20 | 2026-02-19 | Sent | 1500.00 | 150.00 | 1650.00
...`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The export includes:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">All invoices</strong> with key columns (number, client, dates, status, amounts)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Summary row</strong> at the bottom with total revenue and total outstanding</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Formatted headers</strong> — bold, with column width auto-fit</li>
          </ul>

          <Challenge number={11} title="Export to Excel">
            <p>Add an {'"'}Export All to Excel{'"'} button on the invoices list page. Export your invoices. Open the file in Excel, Google Sheets, or LibreOffice Calc. Are all the columns present? Is the data accurate?</p>
          </Challenge>
        </section>

        {/* ═══ Section 12: Build for Production ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Build for Production</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When your invoice app is complete, compile it into a standalone native binary:
          </p>

          <CodeBlock filename="Terminal">
{`grit compile`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This produces a single executable file that includes:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> The Go backend (compiled to native machine code)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> The React frontend (bundled, minified, and embedded)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> SQLite (embedded database, no external server needed)</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The end user just double-clicks the file. No installation, no dependencies, no internet required.
            The SQLite database is created automatically on first launch in the app{"'"}s data directory.
          </p>

          <Tip>
            The compiled binary is typically 20-40 MB. You can reduce it with <Code>-ldflags {'"'}-s -w{'"'}</Code> to
            strip debug symbols, bringing it down to 15-25 MB. That{"'"}s your entire app — backend, frontend,
            and database engine — in a single file.
          </Tip>

          <Challenge number={12} title="Build and Test the Binary">
            <p>Run <Code>grit compile</Code>. Find the output binary. Run it outside of dev mode — does the app open? Create an invoice, export it as PDF. Does everything work the same as in development?</p>
          </Challenge>
        </section>

        {/* ═══ Section 13: Summary + Final Challenge ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What You Learned</h2>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How to scaffold a desktop app with <Code>grit new-desktop</Code></li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Designing related data models (Client, Invoice, LineItem)</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Auto-generating sequential invoice numbers</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Building a multi-step form with dynamic line items</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Implementing a status workflow with valid transitions</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Exporting professional PDF invoices</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Building a financial dashboard with aggregated stats</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Exporting data to Excel for accounting</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Compiling into a standalone native binary</li>
          </ul>

          <Challenge number={13} title="Complete Invoicing Workflow">
            <p>Build the complete workflow:</p>
            <ol className="list-decimal list-inside space-y-1 mt-2">
              <li>Create 3 clients (different companies)</li>
              <li>Create 10 invoices: 2 Draft, 3 Sent, 3 Paid, 2 Overdue</li>
              <li>Each invoice should have 2-5 line items</li>
              <li>Export 2 invoices as PDF</li>
              <li>Export all invoices to Excel</li>
              <li>Check the dashboard — are the stats accurate?</li>
            </ol>
          </Challenge>

          <Challenge number={14} title="Share Your App">
            <p>Compile the app with <Code>grit compile</Code>. Send the binary to a friend or test on a different computer. Does it run without installing Go, Node, or Wails? Does the SQLite database get created automatically on first launch? Create an invoice on the new machine and export it as PDF.</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses', label: 'All Courses' }}
            next={{ href: '/courses', label: 'More Courses' }}
          />
        </div>
      </main>
    </div>
  )
}
