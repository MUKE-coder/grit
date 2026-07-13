import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/field-types')

interface Row {
  type: string
  syntax: string
  go: string
  ts: string
  zod: string
}

// Source of truth: internal/generate/field.go (GoType/TSType/ZodType).
const ROWS: Row[] = [
  { type: 'string', syntax: 'name:string', go: 'string', ts: 'string', zod: 'z.string()' },
  { type: 'text', syntax: 'bio:text', go: 'string', ts: 'string', zod: 'z.string()' },
  { type: 'richtext', syntax: 'body:richtext', go: 'string', ts: 'string', zod: 'z.string()' },
  { type: 'int', syntax: 'qty:int', go: 'int', ts: 'number', zod: 'z.number().int()' },
  { type: 'uint', syntax: 'stock:uint', go: 'uint', ts: 'number', zod: 'z.number().int().nonnegative()' },
  { type: 'float', syntax: 'price:float', go: 'float64', ts: 'number', zod: 'z.number()' },
  { type: 'bool', syntax: 'active:bool', go: 'bool', ts: 'boolean', zod: 'z.boolean()' },
  { type: 'datetime', syntax: 'published_at:datetime', go: '*time.Time', ts: 'string | null', zod: 'z.string().nullable()' },
  { type: 'date', syntax: 'due:date', go: '*time.Time', ts: 'string | null', zod: 'z.string().nullable()' },
  { type: 'slug', syntax: 'slug:slug:title', go: 'string', ts: 'string', zod: 'z.string()' },
  { type: 'belongs_to', syntax: 'category:belongs_to', go: 'string (FK)', ts: 'string', zod: 'z.string().uuid()' },
  { type: 'many_to_many', syntax: 'tags:many_to_many:Tag', go: '[]string', ts: 'string[]', zod: 'z.array(z.string().uuid())' },
  { type: 'string_array', syntax: 'sizes:string_array', go: 'datatypes.JSONSlice[string]', ts: 'string[]', zod: 'z.array(z.string())' },
  { type: 'file', syntax: 'image:file:image', go: '*files.FileRef', ts: 'FileRef | null', zod: 'FileRefSchema.nullable()' },
  { type: 'files', syntax: 'gallery:files:image', go: 'files.FileRefs', ts: 'FileRef[]', zod: 'z.array(FileRefSchema)' },
]

// Default admin rendering per type (form control + table column format).
const RENDER: { type: string; form: string; table: string }[] = [
  { type: 'string', form: 'Text input', table: 'text' },
  { type: 'text', form: 'Textarea', table: 'text' },
  { type: 'richtext', form: 'Rich-text editor (Tiptap)', table: 'richtext' },
  { type: 'int / uint / float', form: 'Number input', table: 'text' },
  { type: 'bool', form: 'Switch', table: 'boolean' },
  { type: 'datetime / date', form: 'Date / time picker', table: 'relative' },
  { type: 'slug', form: 'Auto-filled (read-only)', table: 'text' },
  { type: 'belongs_to', form: 'Relationship select (searchable)', table: 'text (dotted key)' },
  { type: 'many_to_many', form: 'Multi-select', table: '—' },
  { type: 'string_array', form: 'Tag / chips input', table: 'text' },
  { type: 'file', form: 'Dropzone (single)', table: 'file' },
  { type: 'files', form: 'Dropzone (multiple)', table: 'files' },
]

export default function FieldTypesPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Core Concepts · Reference</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Field Types</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                The complete set of field types you can pass to{' '}
                <code>grit generate resource --fields</code>. Each one drives the Go struct
                field, the shared TypeScript type, the Zod schema, and the default admin form
                control and table column &mdash; from one declaration.
              </p>
            </div>

            <div className="prose-grit">
              {/* Syntax */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Field syntax</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  A field is <code>name:type</code>, with optional extra segments. Separate
                  multiple fields with commas inside <code>--fields</code>.
                </p>
                <CodeBlock
                  terminal
                  code={`grit generate resource Product --fields "name:string,price:float,category:belongs_to:Category,gallery:files:image"`}
                />
                <ul className="space-y-2.5 mt-4 mb-4">
                  {[
                    ['Modifiers', 'name:string:unique — append :unique, :required, or :optional. String fields default to required; everything else defaults to optional.'],
                    ['Slug source', 'slug:slug:title — the 3rd segment is the field to slugify. Slugs are auto-unique and generated on save.'],
                    ['belongs_to', 'category:belongs_to (model inferred → Category) or author:belongs_to:User (explicit). Creates a <name>_id UUID foreign-key column.'],
                    ['many_to_many', 'tags:many_to_many:Tag — the related model is required. GORM builds the join table.'],
                    ['File accepts', 'image:file:image, doc:file:all, or att:file:[pdf,doc,image,video]. Aliases: image, video, audio, pdf, doc, all — or a bracketed list.'],
                  ].map(([k, v]) => (
                    <li key={k} className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="text-primary mt-1 font-mono text-xs shrink-0">{k}</span>
                      <span>{v}</span>
                    </li>
                  ))}
                </ul>
              </div>

              {/* Master table */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Type mapping</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  How each field type resolves across the stack. Source of truth:{' '}
                  <code>internal/generate/field.go</code>.
                </p>
                <LaneFlow
                  id="field-types"
                  lanes={['One field declaration', 'Resolves across the stack']}
                  nodes={[
                    { id: 'decl', lane: 0, row: 2, title: 'price:float', sub: '--fields', tone: 'primary' },
                    { id: 'go', lane: 1, row: 0, title: 'Go struct field', sub: 'float64', tone: 'cyan' },
                    { id: 'ts', lane: 1, row: 1, title: 'TypeScript type', sub: 'number', tone: 'blue' },
                    { id: 'zod', lane: 1, row: 2, title: 'Zod rule', sub: 'z.number()', tone: 'violet' },
                    { id: 'form', lane: 1, row: 3, title: 'Admin form', sub: 'number input', tone: 'amber' },
                    { id: 'col', lane: 1, row: 4, title: 'Table column', sub: 'right-aligned', tone: 'green' },
                  ]}
                  edges={[
                    { from: 'decl', to: 'go', tone: 'cyan' },
                    { from: 'decl', to: 'ts', label: 'maps to', tone: 'blue' },
                    { from: 'decl', to: 'zod', tone: 'violet' },
                    { from: 'decl', to: 'form', tone: 'amber' },
                    { from: 'decl', to: 'col', tone: 'green' },
                  ]}
                  legend={[
                    { tone: 'primary', label: 'Declaration' },
                    { tone: 'cyan', label: 'Go' },
                    { tone: 'blue', label: 'TypeScript' },
                    { tone: 'amber', label: 'Admin UI' },
                  ]}
                  caption="Declare a field once — Grit resolves it into Go, TypeScript, Zod, and the admin UI"
                />
                <div className="overflow-x-auto mb-2">
                  <table className="w-full text-sm border border-border rounded-lg">
                    <thead>
                      <tr className="border-b border-border bg-muted/30 text-left">
                        <th className="px-3 py-2 font-medium">Type</th>
                        <th className="px-3 py-2 font-medium">Example</th>
                        <th className="px-3 py-2 font-medium">Go</th>
                        <th className="px-3 py-2 font-medium">TypeScript</th>
                        <th className="px-3 py-2 font-medium">Zod</th>
                      </tr>
                    </thead>
                    <tbody className="font-mono text-[12px]">
                      {ROWS.map((r) => (
                        <tr key={r.type} className="border-b border-border/50 align-top">
                          <td className="px-3 py-2 text-primary">{r.type}</td>
                          <td className="px-3 py-2 text-muted-foreground">{r.syntax}</td>
                          <td className="px-3 py-2">{r.go}</td>
                          <td className="px-3 py-2">{r.ts}</td>
                          <td className="px-3 py-2 text-muted-foreground">{r.zod}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
                <p className="text-xs text-muted-foreground/70">
                  Primary keys are always UUID <code>string</code> (never <code>uint</code>); a{' '}
                  <code>belongs_to</code> FK column is a matching UUID string.
                </p>
              </div>

              {/* Rendering */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Default admin rendering</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Each type also picks a default form control and DataTable column format in the
                  generated admin resource. You can override either in the resource definition.
                </p>
                <div className="overflow-x-auto">
                  <table className="w-full text-sm border border-border rounded-lg">
                    <thead>
                      <tr className="border-b border-border bg-muted/30 text-left">
                        <th className="px-4 py-2 font-medium">Type</th>
                        <th className="px-4 py-2 font-medium">Form control</th>
                        <th className="px-4 py-2 font-medium">Table format</th>
                      </tr>
                    </thead>
                    <tbody>
                      {RENDER.map((r) => (
                        <tr key={r.type} className="border-b border-border/50">
                          <td className="px-4 py-2 font-mono text-[12px] text-primary">{r.type}</td>
                          <td className="px-4 py-2 text-muted-foreground text-[13px]">{r.form}</td>
                          <td className="px-4 py-2 font-mono text-[12px] text-muted-foreground">{r.table}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              <Callout type="escape" title="Escape hatch">
                The generated types are a starting point. Edit the Go struct tags, tune the Zod
                rules in <code>packages/shared/schemas</code>, or override the form/table rendering
                in the admin resource definition &mdash; then run{' '}
                <Link href="/docs/concepts/type-system">grit sync</Link> to keep TypeScript in step
                with Go. Grit generates opinions, not a cage.
              </Callout>

              <div className="flex items-center justify-between border-t border-border pt-8 mt-12">
                <Button variant="ghost" asChild>
                  <Link href="/docs/concepts/type-system" className="gap-2">
                    <ArrowLeft className="h-4 w-4" />
                    Type System
                  </Link>
                </Button>
                <Button variant="ghost" asChild>
                  <Link href="/docs/concepts/generated-files" className="gap-2">
                    Generated File Map
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
