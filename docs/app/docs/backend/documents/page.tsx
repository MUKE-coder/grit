import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/documents')

// The primitives, as data, so `grit docs search` can read them. The index
// reads a const array a page maps over and does not read the text inside a
// hand-written table cell.
const PRIMITIVES = [
  { call: 'New(style)', does: 'A new document with style as its default: font, size, line spacing, justification, margins, portrait or landscape.' },
  { call: 'Heading(level, text)', does: 'A heading at level 1 to 3, as a named Word style rather than bold text. This is what fills the navigation pane and an automatic table of contents.' },
  { call: 'Chapter(title)', does: 'A level 1 heading that starts a new page, using pageBreakBefore so the break moves with the heading when the text above it grows.' },
  { call: 'Para(text)', does: 'A paragraph in the document default.' },
  { call: 'Centered(text)', does: 'A centred paragraph, whatever the default alignment is. A title justified across the page is stretched into something nobody asked for.' },
  { call: 'Rich(runs...)', does: 'A paragraph whose emphasis changes mid-sentence. Each Run carries Bold, Italic and Underline.' },
  { call: 'Bullet(text, level)', does: 'A bulleted list item, nine levels.' },
  { call: 'Numbered(text, level)', does: 'A numbered list item, nine levels.' },
  { call: 'Table(headers, rows)', does: 'A bordered table. Pass no headers for a table with no header row.' },
  { call: 'PageBreak()', does: 'An explicit page break. Prefer Chapter when the next thing is a heading.' },
  { call: 'Bytes() / WriteTo(w)', does: 'Render. Nothing is written until one of these is called.' },
]

const STYLE_FIELDS = [
  { field: 'Font', type: 'string', note: 'A font name, "Times New Roman". Empty means Calibri.' },
  { field: 'Size', type: 'float64', note: 'Points. Zero means 11.' },
  { field: 'LineSpacing', type: 'float64', note: '1, 1.5, 2. Zero means single.' },
  { field: 'Justify', type: 'bool', note: 'Justify body text. Headings, titles and table cells are never justified.' },
  { field: 'Margin', type: 'Margin', note: 'Inches on each side. Zero means one inch all round. Inches because every house style is written in inches or centimetres and nobody has been asked for a 1440 twip margin.' },
  { field: 'Landscape', type: 'bool', note: 'Rotate the page.' },
]

export default function DocumentsPage() {
  return (
    <DocsShell>
      <p className="tag-mono text-accent/80 mb-3">Backend</p>
      <h1 className="text-4xl font-bold tracking-tight mb-4">Word documents</h1>
      <p className="text-lg text-muted-foreground leading-relaxed mb-6">
        <code>internal/docx</code> writes .docx files: a page size, a named font at a
        given size, justified text, a line spacing, and a page break before each
        chapter. It is in every project, and it travels on{' '}
        <code>grit upgrade</code>.
      </p>

      <div className="mb-10 rounded-lg border border-border bg-muted/30 p-4 text-sm text-muted-foreground leading-relaxed">
        <strong className="text-foreground">This is not the list export.</strong> A
        resource already exports its filtered list as CSV, XLSX and JSON, and a Word
        table of rows is that again in a format that holds it worse. What a list
        export cannot produce is a document, and the formatting is usually not a
        preference: a university department or a government office specifies Times
        New Roman 12pt, justified, 1.5 spacing and a page break per chapter, and
        marks against it. Getting those wrong costs the person submitting it.
      </div>

      <h2 className="text-2xl font-bold tracking-tight mb-4">The document a scheme asks for</h2>
      <p className="text-muted-foreground leading-relaxed mb-4">
        <code>Submission()</code> returns that combination, because it is specified
        often enough that every project otherwise writes it out again:
      </p>
      <CodeBlock
        language="go"
        filename="internal/handlers/proposal_handler.go"
        code={`import "myapp/internal/docx"

func (h *ProposalHandler) Download(c *gin.Context) {
    proposal, err := h.service().One(h.ctx(c), c.Param("id"))
    if err != nil {
        respond.NotFound(c, "Proposal not found")
        return
    }

    doc := docx.New(docx.Submission()) // Times New Roman 12pt, justified, 1.5, 1in
    doc.Heading(1, proposal.Title)
    doc.Centered(proposal.Author)

    for _, chapter := range proposal.Chapters {
        doc.Chapter(chapter.Title)   // starts a new page
        doc.Para(chapter.Body)
    }

    b, err := doc.Bytes()
    if err != nil {
        respond.Internal(c, err)
        return
    }
    c.Header("Content-Disposition", docx.Attachment(proposal.Title+".docx"))
    c.Data(http.StatusOK, docx.ContentType, b)
}`}
      />
      <p className="text-muted-foreground leading-relaxed mb-10">
        <code>Attachment</code> is not decoration. A filename holding a quote or a
        newline either truncates the header or injects another one, and a filename is
        the first thing in a response to come from user input. It strips what is
        unsafe and also sends the RFC 5987 form, so an accented or non-Latin title
        survives.
      </p>

      <h2 className="text-2xl font-bold tracking-tight mb-4">Style</h2>
      <div className="overflow-x-auto rounded-xl border border-border/50 mb-10">
        <table className="w-full text-sm">
          <thead className="bg-card/60">
            <tr className="text-left">
              <th className="px-4 py-3 font-semibold">Field</th>
              <th className="px-4 py-3 font-semibold">Type</th>
              <th className="px-4 py-3 font-semibold">What it does</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/40">
            {STYLE_FIELDS.map((f) => (
              <tr key={f.field}>
                <td className="px-4 py-3 font-mono text-[13px] align-top whitespace-nowrap">{f.field}</td>
                <td className="px-4 py-3 font-mono text-[13px] align-top text-muted-foreground">{f.type}</td>
                <td className="px-4 py-3 text-muted-foreground leading-relaxed align-top">{f.note}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <h2 className="text-2xl font-bold tracking-tight mb-4">Primitives</h2>
      <div className="overflow-x-auto rounded-xl border border-border/50 mb-10">
        <table className="w-full text-sm">
          <thead className="bg-card/60">
            <tr className="text-left">
              <th className="px-4 py-3 font-semibold">Call</th>
              <th className="px-4 py-3 font-semibold">What it does</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/40">
            {PRIMITIVES.map((p) => (
              <tr key={p.call}>
                <td className="px-4 py-3 font-mono text-[13px] align-top whitespace-nowrap">{p.call}</td>
                <td className="px-4 py-3 text-muted-foreground leading-relaxed align-top">{p.does}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <h2 className="text-2xl font-bold tracking-tight mb-4">A whole report in one call</h2>
      <p className="text-muted-foreground leading-relaxed mb-4">
        <code>RenderReport</code> is the pre-built template, the same arrangement{' '}
        <code>internal/pdf</code> has with <code>RenderInvoice</code>. Copy{' '}
        <code>report.go</code> as a starting point for a proposal, a dissertation, a
        tender response or a board paper:
      </p>
      <CodeBlock
        language="go"
        code={`b, err := docx.RenderReport(docx.Report{
    Title:  "An Assessment of Water Access in Mukono District",
    Author: "A Student",
    Date:   time.Now(),
    Chapters: []docx.ReportChapter{
        {Title: "Introduction", Paragraphs: []string{intro}},
        {Title: "Methodology", Sections: []docx.ReportSection{{
            Heading: "Sampling",
            Bullets: []string{"Forty households", "Two sub-counties"},
            Table:   &docx.ReportTable{Headers: []string{"Sub-county", "N"}, Rows: rows},
        }}},
    },
}, docx.Submission())`}
      />

      <h2 className="text-2xl font-bold tracking-tight mb-4">Why there is no dependency</h2>
      <p className="text-muted-foreground leading-relaxed mb-4">
        A .docx is a zip of XML parts, and this writes them with{' '}
        <code>archive/zip</code> and <code>encoding/xml</code> from the standard
        library. The Go libraries for this are either thinly maintained or AGPL, and
        an AGPL dependency inside a framework other people build commercial products
        on is a licence decision nobody would choose deliberately. The subset a report
        needs is small enough that owning it costs less than depending on one.
      </p>
      <p className="text-muted-foreground leading-relaxed mb-10">
        The package ships its own tests, which a project can run. They were written
        against an independent OOXML reader rather than against the writer&apos;s own
        output, and that caught two faults before any of this shipped: a missing{' '}
        <code>&lt;w:tblGrid&gt;</code>, which the schema requires and Word silently
        tolerates, and the font living only in <code>docDefaults</code>, so a reader
        asking the Normal style what font it was got nothing. The second is exactly
        the claim a marking scheme is checked against.
      </p>

      <h2 className="text-2xl font-bold tracking-tight mb-4">Citations</h2>
      <p className="text-muted-foreground leading-relaxed mb-10">
        A document with a bibliography needs a CSL processor, and that belongs on the
        client today. See{' '}
        <Link href="/docs/frontend/citations" className="underline underline-offset-4 hover:text-foreground">
          citations
        </Link>{' '}
        for why, and for the one rule that matters: render once, persist the result,
        and have the server read it rather than derive it a second time.
      </p>

      <div className="flex items-center justify-between border-t border-border/40 pt-6">
        <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
          <Link href="/docs/backend/invoices" className="gap-1.5">
            <ArrowLeft className="h-3.5 w-3.5" />
            Invoices &amp; Line Items
          </Link>
        </Button>
        <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
          <Link href="/docs/frontend/citations" className="gap-1.5">
            Citations
            <ArrowRight className="h-3.5 w-3.5" />
          </Link>
        </Button>
      </div>
    </DocsShell>
  )
}
