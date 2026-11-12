import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/frontend/citations')

// Measured on one Windows machine with go1.27.3, so read the ratio rather than
// the absolute numbers. The ratio is the part that reproduces anywhere.
const ENGINE_COST = [
  { style: 'harvard-cite-them-right', size: '10.5 KB', goja: '0.44 s', node: '' },
  { style: 'ieee', size: '17.6 KB', goja: '0.30 s', node: '' },
  { style: 'modern-language-association', size: '44 KB', goja: '1.1 s', node: '' },
  { style: 'apa', size: '86 KB', goja: '44 to 80 s', node: '0.61 s' },
  { style: 'chicago-author-date', size: '168 KB', goja: '209 s', node: '' },
]

export default function CitationsPage() {
  return (
    <DocsShell>
      <p className="tag-mono text-accent/80 mb-3">Frontend</p>
      <h1 className="text-4xl font-bold tracking-tight mb-4">Citations and bibliographies</h1>
      <p className="text-lg text-muted-foreground leading-relaxed mb-6">
        Formatting a citation in APA, Harvard, IEEE or Chicago means running a CSL
        processor. Grit does not ship one, and this page is the reason why, plus the
        arrangement that works today and the one rule that will cost you a day if you
        get it wrong.
      </p>

      <div className="mb-10 rounded-lg border border-border bg-muted/30 p-4 text-sm text-muted-foreground leading-relaxed">
        <strong className="text-foreground">Do not write your own.</strong> A citation
        style is not a format string. APA alone is 86 KB of XML describing
        disambiguation, name particles, et al. collapsing, ordering and
        substitution, and a style is revised. Take a processor.
      </div>

      <h2 className="text-2xl font-bold tracking-tight mb-4">Render on the client</h2>
      <p className="text-muted-foreground leading-relaxed mb-4">
        The only complete CSL implementation is citeproc-js, the engine behind Zotero
        and Mendeley. Run it in the browser, lazily, so the pages that do not cite
        anything never pay for it:
      </p>
      <CodeBlock
        language="ts"
        filename="apps/web/lib/citations.ts"
        code={`// citeproc is ~967 KB of ES5 and a style is another 10 to 170 KB, so this
// is a dynamic import and never a top-level one. For comparison, the whole
// rich text editor in a Grit project is 173 KB raw and 18.6 KB gzipped: a
// static import here would dwarf everything else the app ships, on a phone,
// on a slow connection.
let enginePromise: Promise<CiteprocEngine> | null = null

export async function engine(styleId: string) {
  if (!enginePromise) {
    enginePromise = (async () => {
      const { default: CSL } = await import('citeproc')
      const [style, locale] = await Promise.all([
        fetch(\`/csl/\${styleId}.csl\`).then((r) => r.text()),
        fetch('/csl/locales-en-GB.xml').then((r) => r.text()),
      ])
      return new CSL.Engine(
        {
          retrieveLocale: () => locale,
          retrieveItem: (id: string) => itemsById[id],
        },
        style,
      )
    })()
  }
  return enginePromise
}`}
      />

      <h2 className="text-2xl font-bold tracking-tight mb-4">Persist what it rendered</h2>
      <p className="text-muted-foreground leading-relaxed mb-4">
        This is the rule. When the document is later exported as{' '}
        <Link href="/docs/backend/documents" className="underline underline-offset-4 hover:text-foreground">
          a Word file
        </Link>{' '}
        or a PDF, the server needs the same citation text the person saw on screen.
        Store what the client rendered and have the server read it. Do not give the
        server a second implementation that derives it again.
      </p>
      <p className="text-muted-foreground leading-relaxed mb-4">
        Two implementations of one rule in two languages drift, and they drift
        quietly. This exact shape has already cost a day in a Grit project: an FNV
        hash written once in TypeScript and once in Go stopped agreeing, nothing
        failed, and every passage simply reported as unrecorded until a golden-file
        fixture that ran both sides caught it.
      </p>
      <CodeBlock
        language="bash"
        code={`grit generate resource Reference \\
  --fields "csl_json:json,rendered:text,style:string,rendered_at:datetime,document_id:belongs_to:Document"`}
      />
      <p className="text-muted-foreground leading-relaxed mb-10">
        <code>csl_json</code> is the source data in CSL-JSON, which is the
        interchange format every reference manager exports.{' '}
        <code>rendered</code> is what the engine produced, with the{' '}
        <code>style</code> it was produced in, so a change of style is a
        re-render and not a silent mismatch. The export reads{' '}
        <code>rendered</code>.
      </p>

      <h2 className="text-2xl font-bold tracking-tight mb-4">Why not on the server</h2>
      <p className="text-muted-foreground leading-relaxed mb-4">
        Because it was tried. Running citeproc-js under goja, the pure-Go JavaScript
        runtime, produces correct output, and the cost of building the engine, which
        is parsing the style XML, is this:
      </p>
      <div className="overflow-x-auto rounded-xl border border-border/50 mb-4">
        <table className="w-full text-sm">
          <thead className="bg-card/60">
            <tr className="text-left">
              <th className="px-4 py-3 font-semibold">Style</th>
              <th className="px-4 py-3 font-semibold">Size</th>
              <th className="px-4 py-3 font-semibold">goja</th>
              <th className="px-4 py-3 font-semibold">node</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/40">
            {ENGINE_COST.map((row) => (
              <tr key={row.style}>
                <td className="px-4 py-3 font-mono text-[13px] align-top whitespace-nowrap">{row.style}</td>
                <td className="px-4 py-3 text-muted-foreground align-top">{row.size}</td>
                <td className="px-4 py-3 text-muted-foreground align-top">{row.goja}</td>
                <td className="px-4 py-3 text-muted-foreground align-top">{row.node}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <p className="text-muted-foreground leading-relaxed mb-4">
        It is a cliff, not a slope: everything under about 45 KB is sub-second, and
        APA and Chicago fall off it. goja is 70 to 130 times slower than V8 at
        parsing a style, and a goja VM is not safe for concurrent use, so a server
        would need one engine per style per worker and the warm-up multiplies by
        both. APA and Chicago are the two styles most people ask for. Rendering,
        once an engine exists, is fine in either: roughly 200 ms for a bibliography.
      </p>
      <p className="text-muted-foreground leading-relaxed mb-10">
        <code>v8go</code> would fix the speed and brings <code>CGO_ENABLED=1</code>,
        a C toolchain in the image, no more cross-compiling and about 30 MB of image
        growth. Grit has a precedent for exactly that trade in the{' '}
        <Link href="/docs/batteries/media" className="underline underline-offset-4 hover:text-foreground">
          vips image backend
        </Link>
        , which is pure Go by default and cgo behind a build tag. If server-side
        rendering becomes necessary, that is the shape it should take. It is not the
        default.
      </p>

      <h2 className="text-2xl font-bold tracking-tight mb-4">The licence</h2>
      <p className="text-muted-foreground leading-relaxed mb-10">
        citeproc-js is dual licensed CPAL-1.0 or AGPL-1.0, which is why it is a
        dependency your project adds deliberately rather than one Grit puts in every
        scaffolded app.{' '}
        <code>@citation-js/plugin-csl</code> is MIT but depends on{' '}
        <code>citeproc</code>, so it does not escape it. The style files themselves
        are CC-BY-SA, which is fine as data. Read the terms against what you are
        building before you ship it.
      </p>

      <h2 className="text-2xl font-bold tracking-tight mb-4">What would change this</h2>
      <p className="text-muted-foreground leading-relaxed mb-10">
        A native Go CSL processor.{' '}
        <a
          href="https://github.com/front-matter/citeproc"
          className="underline underline-offset-4 hover:text-foreground"
        >
          front-matter/citeproc
        </a>{' '}
        is MIT, native Go and under active development, and its README says basic
        functionality is not yet working with a first working release expected in
        December 2026. When it renders, server-side citations become a few lines and
        this page gets shorter. Until then it is worth watching and not worth
        depending on.
      </p>

      <div className="flex items-center justify-between border-t border-border/40 pt-6">
        <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
          <Link href="/docs/backend/documents" className="gap-1.5">
            <ArrowLeft className="h-3.5 w-3.5" />
            Word documents
          </Link>
        </Button>
        <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
          <Link href="/docs/backend/invoices" className="gap-1.5">
            Invoices &amp; Line Items
            <ArrowRight className="h-3.5 w-3.5" />
          </Link>
        </Button>
      </div>
    </DocsShell>
  )
}
