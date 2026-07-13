import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'PDF & Excel Export — Grit Desktop Course',
  description: 'Generate PDFs, Excel spreadsheets, and CSV files from your Grit desktop app. Learn the export service, download flow, and custom report templates.',
}

export default function ExportCourse() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <Link href="/courses/grit-desktop" className="hover:text-foreground transition-colors">Grit Desktop</Link>
          <span>/</span>
          <span className="text-foreground">PDF & Excel Export</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Course 4 of 5</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">10 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            PDF & Excel Export
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            In this course, you will learn how to export data from your desktop app into PDF documents,
            Excel spreadsheets, and CSV files. You{"'"}ll understand the export service, the download flow,
            and how to customize report templates.
          </p>
        </div>

        <div className="my-4 rounded-lg border border-primary/20 bg-primary/5 px-4 py-3">
          <p className="text-sm text-muted-foreground">
            <strong className="text-foreground">Reference docs:</strong>{' '}
            <a href="/docs/admin/datatable" className="text-primary hover:underline">DataTable & Export →</a>
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Why Export? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Why Export?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Users need to get data <strong className="text-foreground">out</strong> of your application.
            Common use cases include:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Invoices</strong> — generate a PDF invoice to email or print for a client</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Reports</strong> — export monthly data summaries for management</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Data backups</strong> — export all records to CSV for archiving or migration</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Spreadsheets</strong> — export data to Excel for further analysis in tools like Google Sheets</li>
          </ul>

          <Definition term="Export">
            Converting application data into a downloadable file format. The three most common
            formats are PDF (for human-readable documents), Excel (for structured data that can be
            edited), and CSV (for universal data interchange).
          </Definition>
        </section>

        {/* ═══ The Export Service ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The Export Service</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit desktop apps include a built-in export service in Go with three main methods. The
            service lives in the <Code>internal/service/</Code> directory:
          </p>

          <CodeBlock filename="internal/service/export.go (overview)">
{`type ExportService struct {
    db *gorm.DB
}

func NewExportService(db *gorm.DB) *ExportService {
    return &ExportService{db: db}
}

// GeneratePDF creates a formatted PDF document from records
func (s *ExportService) GeneratePDF(resource string) (string, error) {
    // Queries data, formats into PDF, saves to temp file
    // Returns the file path
}

// GenerateExcel creates an .xlsx spreadsheet from records
func (s *ExportService) GenerateExcel(resource string) (string, error) {
    // Queries data, creates Excel workbook with headers and rows
    // Returns the file path
}

// GenerateCSV creates a .csv file from records
func (s *ExportService) GenerateCSV(resource string) (string, error) {
    // Queries data, writes CSV with headers and rows
    // Returns the file path
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Each method queries the database, formats the data, writes it to a temporary file, and returns
            the file path. The React frontend then uses this path to trigger a download.
          </p>

          <Challenge number={1} title="Find the Export Service">
            <p>Open the export service file in your project (look in <Code>internal/service/</Code>). What methods does it have? What Go libraries does it import for PDF and Excel generation?</p>
          </Challenge>
        </section>

        {/* ═══ Generating a PDF ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Generating a PDF</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The PDF export creates a formatted document with a title, headers, and a data table.
            The Go code uses a PDF library to build the document programmatically:
          </p>

          <CodeBlock filename="PDF generation (simplified)">
{`func (s *ExportService) GeneratePDF(resource string) (string, error) {
    // 1. Query all records from the database
    var blogs []models.Blog
    s.db.Find(&blogs)

    // 2. Create a new PDF document
    pdf := gopdf.New()
    pdf.AddPage()

    // 3. Add a title
    pdf.SetFont("Helvetica", "B", 18)
    pdf.Cell(nil, "Blog Export Report")

    // 4. Add a table with headers
    pdf.SetFont("Helvetica", "B", 10)
    headers := []string{"ID", "Title", "Created At", "Status"}
    for _, h := range headers {
        pdf.Cell(nil, h)
    }

    // 5. Add data rows
    pdf.SetFont("Helvetica", "", 10)
    for _, blog := range blogs {
        pdf.Cell(nil, blog.Title)
        pdf.Cell(nil, blog.CreatedAt.Format("2006-01-02"))
    }

    // 6. Save to temp file and return path
    path := filepath.Join(os.TempDir(), "blog-export.pdf")
    pdf.WritePdf(path)
    return path, nil
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The generated PDF includes proper formatting: bold headers, alternating row colors,
            page numbers, and timestamps.
          </p>

          <Challenge number={2} title="Export Blogs to PDF">
            <p>Make sure you have at least 5 blog posts in your app. Find the export button in the blog list page and click it to generate a PDF. Does the file download? Open it — are all your blog posts listed?</p>
          </Challenge>
        </section>

        {/* ═══ Generating Excel Files ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Generating Excel Files</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Excel exports create <Code>.xlsx</Code> files with proper sheets, headers, and data rows.
            The Go code uses the <Code>excelize</Code> library:
          </p>

          <CodeBlock filename="Excel generation (simplified)">
{`func (s *ExportService) GenerateExcel(resource string) (string, error) {
    // 1. Query all records
    var blogs []models.Blog
    s.db.Find(&blogs)

    // 2. Create a new Excel workbook
    f := excelize.NewFile()
    sheet := "Blogs"
    f.SetSheetName("Sheet1", sheet)

    // 3. Write headers (row 1)
    headers := []string{"ID", "Title", "Content", "Created At"}
    for i, h := range headers {
        cell := fmt.Sprintf("%c1", 'A'+i)
        f.SetCellValue(sheet, cell, h)
    }

    // 4. Style headers (bold, background color)
    style, _ := f.NewStyle(&excelize.Style{
        Font: &excelize.Font{Bold: true},
        Fill: excelize.Fill{Type: "pattern", Color: []string{"#6c5ce7"}, Pattern: 1},
    })
    f.SetCellStyle(sheet, "A1", "D1", style)

    // 5. Write data rows
    for i, blog := range blogs {
        row := i + 2
        f.SetCellValue(sheet, fmt.Sprintf("A%d", row), blog.ID)
        f.SetCellValue(sheet, fmt.Sprintf("B%d", row), blog.Title)
        f.SetCellValue(sheet, fmt.Sprintf("C%d", row), blog.Content)
        f.SetCellValue(sheet, fmt.Sprintf("D%d", row), blog.CreatedAt.Format("2006-01-02"))
    }

    // 6. Save to temp file
    path := filepath.Join(os.TempDir(), "blog-export.xlsx")
    f.SaveAs(path)
    return path, nil
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The Excel file includes styled headers, proper column widths, and typed cells (numbers
            stay as numbers, dates as dates).
          </p>

          <Challenge number={3} title="Export to Excel">
            <p>Export your blog data to Excel. Open the <Code>.xlsx</Code> file in Excel, Google Sheets, or LibreOffice. Does it have proper column headers? Are the dates formatted correctly?</p>
          </Challenge>
        </section>

        {/* ═══ CSV Export ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">CSV Export</h2>

          <Definition term="CSV (Comma-Separated Values)">
            The simplest data format — plain text with values separated by commas. Each line is a row,
            the first line is headers. CSV files can be opened by any spreadsheet app, text editor,
            or programming language. They are the universal format for data exchange.
          </Definition>

          <CodeBlock filename="CSV generation (simplified)">
{`func (s *ExportService) GenerateCSV(resource string) (string, error) {
    var blogs []models.Blog
    s.db.Find(&blogs)

    path := filepath.Join(os.TempDir(), "blog-export.csv")
    file, _ := os.Create(path)
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Write headers
    writer.Write([]string{"ID", "Title", "Content", "Created At"})

    // Write data rows
    for _, blog := range blogs {
        writer.Write([]string{
            fmt.Sprintf("%d", blog.ID),
            blog.Title,
            blog.Content,
            blog.CreatedAt.Format("2006-01-02"),
        })
    }

    return path, nil
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            CSV is the simplest format and works everywhere. The output looks like this:
          </p>

          <CodeBlock filename="blog-export.csv">
{`ID,Title,Content,Created At
1,My First Post,Hello world...,2026-03-15
2,Second Post,More content...,2026-03-16
3,Third Post,Even more...,2026-03-17`}
          </CodeBlock>

          <Challenge number={4} title="Export to CSV">
            <p>Export your blog data to CSV. Open the file in a text editor (not a spreadsheet app). Can you read the raw data? Each line should be one blog post with comma-separated values.</p>
          </Challenge>
        </section>

        {/* ═══ Download Flow in Desktop ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Download Flow in Desktop</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The export flow in a desktop app works differently from a web app. There is no browser
            download dialog — instead, Go generates the file on disk and the React frontend opens it:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">1.</span> React calls the Go export method via Wails binding</li>
            <li className="flex gap-2"><span className="text-primary">2.</span> Go queries the database and generates the file in a temp directory</li>
            <li className="flex gap-2"><span className="text-primary">3.</span> Go returns the file path to React</li>
            <li className="flex gap-2"><span className="text-primary">4.</span> React uses the Wails runtime to open the file or show a save dialog</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            On the Go side, the bound method looks like:
          </p>

          <CodeBlock filename="app.go (export methods)">
{`// Export blogs to PDF — returns the file path
func (a *App) ExportBlogsPDF() (string, error) {
    return a.exportService.GeneratePDF("blogs")
}

// Export blogs to Excel
func (a *App) ExportBlogsExcel() (string, error) {
    return a.exportService.GenerateExcel("blogs")
}

// Export blogs to CSV
func (a *App) ExportBlogsCSV() (string, error) {
    return a.exportService.GenerateCSV("blogs")
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            On the React side:
          </p>

          <CodeBlock filename="React export button">
{`import { ExportBlogsPDF } from '../../../wailsjs/go/main/App'
import { BrowserOpenURL } from '../../../wailsjs/runtime'

async function handleExportPDF() {
  const filePath = await ExportBlogsPDF()
  // Open the file with the system default app
  BrowserOpenURL("file://" + filePath)
}`}
          </CodeBlock>

          <Challenge number={5} title="Trace the Export Flow">
            <p>Click an export button in the UI. Where does the file save? Check your system{"'"}s temp directory. Can you find the exported file there? Try opening it directly from that location.</p>
          </Challenge>

          <Challenge number={6} title="Export All Three Formats">
            <p>Export your blog data to all three formats: PDF, Excel, and CSV. Open each file. Compare them — which format is best for printing? Which is best for editing? Which is best for importing into another app?</p>
          </Challenge>
        </section>

        {/* ═══ Custom Report Templates ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Custom Report Templates</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You can customize the PDF layout to match your brand. Common customizations include:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Header</strong> — add your app name, logo, or company information</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Footer</strong> — page numbers, generation date, confidentiality notice</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Colors</strong> — match your brand{"'"}s accent color for headers and borders</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Layout</strong> — landscape vs portrait, margins, font sizes</li>
          </ul>

          <CodeBlock filename="Custom PDF header example">
{`func (s *ExportService) addHeader(pdf *gopdf.GoPdf, title string) {
    // Company name
    pdf.SetFont("Helvetica", "B", 24)
    pdf.SetTextColor(108, 92, 231) // Grit purple
    pdf.Cell(nil, "My Company")

    // Report title
    pdf.SetFont("Helvetica", "", 14)
    pdf.SetTextColor(144, 144, 168) // Muted text
    pdf.Cell(nil, title)

    // Date
    pdf.SetFont("Helvetica", "", 10)
    pdf.Cell(nil, "Generated: " + time.Now().Format("January 2, 2006"))

    // Divider line
    pdf.Line(20, pdf.GetY()+5, 575, pdf.GetY()+5)
}`}
          </CodeBlock>

          <Tip>
            For complex PDF layouts (invoices, receipts, certificates), consider using HTML-to-PDF conversion
            instead of building the layout with Go code. You can write the template in HTML/CSS and convert
            it to PDF, giving you much more control over the design.
          </Tip>

          <Challenge number={7} title="Customize the PDF Header">
            <p>Open the export service and modify the PDF generation code. Add a custom header with your app{"'"}s name and today{"'"}s date. Change the header color to match your theme. Export again — does the new header appear?</p>
          </Challenge>
        </section>

        {/* ═══ Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What You Learned</h2>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Why export matters — invoices, reports, data backups, spreadsheets</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> The export service with GeneratePDF, GenerateExcel, and GenerateCSV methods</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How PDF generation works — headers, tables, formatting with Go libraries</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How Excel generation works — workbooks, sheets, styled headers, typed cells</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How CSV export works — the simplest, most universal data format</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> The desktop download flow — Go generates file, returns path, React opens it</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How to customize PDF templates with headers, footers, and brand colors</li>
          </ul>

          <Challenge number={8} title="Build an Invoice System">
            <p>Generate an Invoice resource with fields: <Code>customer:string,amount:float,date:date,status:string,paid:bool</Code>. Add 10 invoices with different customers and amounts.</p>
          </Challenge>

          <Challenge number={9} title="Export Invoices">
            <p>Export your invoices to PDF, Excel, and CSV. The PDF should look like a real invoice report. The Excel file should have proper number formatting. The CSV should be importable into Google Sheets.</p>
          </Challenge>

          <Challenge number={10} title="Compare the Three Formats">
            <p>Open all three exported files side by side. Which format preserves the formatting best (PDF)? Which is easiest to edit (Excel)? Which is the smallest file size (CSV)? When would you use each one?</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses/grit-desktop/custom-ui', label: 'Prev: Custom UI & Theming' }}
            next={{ href: '/courses/grit-desktop/build', label: 'Next: Build & Distribution' }}
          />
        </div>
      </main>
    </div>
  )
}
