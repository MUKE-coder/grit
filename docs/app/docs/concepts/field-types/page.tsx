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
  { type: 'float', syntax: 'weight:float', go: 'float64', ts: 'number', zod: 'z.number()' },
  { type: 'money', syntax: 'price:money', go: 'money.Money', ts: 'Money', zod: 'MoneySchema' },
  { type: 'bool', syntax: 'active:bool', go: 'bool', ts: 'boolean', zod: 'z.boolean()' },
  { type: 'toggle', syntax: 'active:toggle', go: 'bool', ts: 'boolean', zod: 'z.boolean()' },
  { type: 'select', syntax: 'status:select:draft=Draft|paid=Paid', go: 'string', ts: '"draft" | "paid"', zod: 'z.enum([...])' },
  { type: 'radio', syntax: 'plan:radio:free=Free|pro=Pro', go: 'string', ts: '"free" | "pro"', zod: 'z.enum([...])' },
  { type: 'check', syntax: 'tags:check:news=News|ops=Ops', go: 'datatypes.JSONSlice[string]', ts: '("news" | "ops")[]', zod: 'z.array(z.enum([...]))' },
  { type: 'datetime', syntax: 'published_at:datetime', go: '*time.Time', ts: 'string | null', zod: 'z.string().nullable()' },
  { type: 'date', syntax: 'due:date', go: '*time.Time', ts: 'string | null', zod: 'z.string().nullable()' },
  { type: 'slug', syntax: 'slug:slug:title', go: 'string', ts: 'string', zod: 'z.string()' },
  { type: 'belongs_to', syntax: 'category:belongs_to', go: 'string (FK)', ts: 'string', zod: 'z.string().uuid()' },
  { type: 'many_to_many', syntax: 'tags:many_to_many:Tag', go: '[]string', ts: 'string[]', zod: 'z.array(z.string().uuid())' },
  { type: 'string_array', syntax: 'sizes:string_array', go: 'datatypes.JSONSlice[string]', ts: 'string[]', zod: 'z.array(z.string())' },
  { type: 'file', syntax: 'image:file:image', go: '*files.FileRef', ts: 'FileRef | null', zod: 'FileRefSchema.nullable()' },
  { type: 'files', syntax: 'gallery:files:image', go: 'files.FileRefs', ts: 'FileRef[]', zod: 'z.array(FileRefSchema)' },
  { type: 'email', syntax: 'email:email:unique', go: 'string (size 254)', ts: 'string', zod: 'EmailSchema' },
  { type: 'url', syntax: 'website:url', go: 'string (size 2048)', ts: 'string', zod: 'UrlSchema' },
  { type: 'domain', syntax: 'host:domain', go: 'string (size 253)', ts: 'string', zod: 'DomainSchema' },
  { type: 'tel', syntax: 'phone:tel:UG', go: 'string (size 32, E.164)', ts: 'string', zod: 'PhoneSchema' },
  { type: 'country', syntax: 'country:country:UG', go: 'string (size 2)', ts: 'string', zod: 'CountrySchema' },
  { type: 'color', syntax: 'brand:color', go: 'string (size 7)', ts: 'string', zod: 'ColorSchema' },
  { type: 'percent', syntax: 'discount:percent', go: 'float64 (decimal 5,2)', ts: 'number', zod: 'PercentSchema' },
  { type: 'rating', syntax: 'score:rating:10', go: 'int', ts: 'number', zod: 'ratingSchema(10)' },
  { type: 'time', syntax: 'opens_at:time', go: 'string (HH:MM)', ts: 'string', zod: 'TimeSchema' },
  { type: 'json', syntax: 'settings:json', go: 'datatypes.JSON', ts: 'unknown', zod: 'JsonValueSchema' },
]

// The formatted types: what the API stores, and what it refuses.
const FORMATTED: { type: string; stored: string; refused: string }[] = [
  { type: 'email', stored: 'Lowercased and trimmed: ada@example.com', refused: 'Anything that is not a bare address, including "Ada <ada@example.com>"' },
  { type: 'url', stored: 'As given', refused: 'Any scheme but http and https, or no host' },
  { type: 'domain', stored: 'A bare host in lowercase punycode: a pasted https://www.Example.co.ug/about becomes www.example.co.ug, and münchen.de becomes xn--mnchen-3ya.de', refused: 'A name DNS could not resolve, such as one without a dot' },
  { type: 'tel', stored: 'E.164: 0772 123456 in Uganda becomes +256772123456', refused: 'A number that is not valid for its country, by libphonenumber\'s metadata' },
  { type: 'country', stored: 'ISO 3166-1 alpha-2 in upper case: UG', refused: 'A code ISO has not assigned' },
  { type: 'color', stored: '#rrggbb in lower case: #ABC becomes #aabbcc', refused: 'Anything that is not a hex colour' },
  { type: 'percent', stored: 'Rounded to two decimals', refused: 'Below 0 or above 100' },
  { type: 'rating', stored: 'Whole stars; 0 means not rated', refused: 'A fraction, or more stars than the field has' },
  { type: 'time', stored: 'HH:MM on a 24-hour clock; seconds are dropped', refused: 'A time that does not exist, such as 25:00' },
  { type: 'json', stored: 'jsonb on Postgres, JSON on MySQL, text on SQLite', refused: 'Text that does not parse as JSON' },
]

// Default admin rendering per type (form control + table column format).
const RENDER: { type: string; form: string; table: string }[] = [
  { type: 'string', form: 'Text input', table: 'text' },
  { type: 'text', form: 'Textarea', table: 'text' },
  { type: 'richtext', form: 'Rich-text editor (Tiptap)', table: 'richtext' },
  { type: 'int / uint / float', form: 'Number input', table: 'text' },
  { type: 'money', form: 'Amount input + currency picker', table: 'money' },
  { type: 'bool', form: 'Switch', table: 'boolean' },
  { type: 'toggle', form: 'Switch', table: 'boolean' },
  { type: 'select', form: 'Searchable dropdown (value=Label options)', table: 'text' },
  { type: 'radio', form: 'Radio-button group (single choice)', table: 'text' },
  { type: 'check', form: 'Checkbox group (multi)', table: '—' },
  { type: 'datetime / date', form: 'Date / time picker', table: 'relative' },
  { type: 'slug', form: 'Auto-filled (read-only)', table: 'text' },
  { type: 'belongs_to', form: 'Relationship select (searchable)', table: 'text (dotted key)' },
  { type: 'many_to_many', form: 'Multi-select', table: '—' },
  { type: 'string_array', form: 'Tag / chips input', table: 'text' },
  { type: 'file', form: 'Dropzone (single)', table: 'file' },
  { type: 'files', form: 'Dropzone (multiple)', table: 'files' },
  { type: 'email', form: 'Email input, lowercased on blur', table: 'email (mailto: link)' },
  { type: 'url', form: 'URL input', table: 'link (opens in a new tab)' },
  { type: 'domain', form: 'Domain input that strips a pasted scheme and path', table: 'domain (link to https://)' },
  { type: 'tel', form: 'Searchable country picker with flags and dial codes, number formatted as you type', table: 'tel (formatted, tel: link)' },
  { type: 'country', form: 'Searchable country picker', table: 'country (flag and name)' },
  { type: 'color', form: 'Colour picker plus hex box', table: 'color (swatch and hex)' },
  { type: 'percent', form: 'Number input with a % suffix', table: 'percent (12.5%)' },
  { type: 'rating', form: 'Stars (a keyboard radio group)', table: 'rating (stars)' },
  { type: 'time', form: 'Time input', table: 'time (2:30 PM in the viewer\'s clock)' },
  { type: 'json', form: 'JSON editor that checks as you type', table: 'json (collapsible)' },
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
                  code={`grit generate resource Product --fields "name:string,price:money,category:belongs_to:Category,gallery:files:image"`}
                />
                <ul className="space-y-2.5 mt-4 mb-4">
                  {[
                    ['Modifiers', 'name:string:unique — append :unique, :required, or :optional. String fields default to required; everything else defaults to optional.'],
                    ['Auto-number', 'number:string:auto:INV (string only) — the server fills it from an atomic, gap-free counter (INV-202607-0001) in BeforeCreate. Optional and hidden from the form; the prefix is optional. Shorthand for grit generate sequence; use that command directly for yearly/never resets or a custom width.'],
                    ['Encrypted', 'notes:text:encrypted (string, text and richtext only). Stored with AES-256-GCM at rest, plaintext in code and over the API; needs FIELD_ENCRYPTION_KEY. Cannot be :unique, and is left out of search, sorting and filters, because ciphertext changes on every write.'],
                    ['Options', 'status:select:draft=Draft|sent=Sent (select / radio / check) — pipe-separated choices. Labels are OPTIONAL: status:select:draft|sent|paid generates the labels by capitalizing each value (in_progress → In Progress). Use value=Label only when the label differs from the stored value.'],
                    ['Slug source', 'slug:slug:title — the 3rd segment is the field to slugify. Slugs are auto-unique and generated on save.'],
                    ['belongs_to', 'category:belongs_to (model inferred → Category) or author:belongs_to:User (explicit). Creates a <name>_id UUID foreign-key column.'],
                    ['many_to_many', 'tags:many_to_many:Tag — the related model is required. GORM builds the join table.'],
                    ['File accepts', 'image:file:image, doc:file:all, or att:file:[pdf,doc,image,video]. Aliases: image, video, audio, pdf, doc, all — or a bracketed list.'],
                    ['Money', 'price:money — stores an integer count of minor units plus an ISO 4217 currency code, in two columns (price_amount, price_currency). Use it for anything you will add up. See the money page.'],
                    ['tel option', 'phone:tel:UG sets the country the picker starts on and a local number is read in. Without it the picker starts on the browser\'s country. Modifiers follow the option: phone:tel:UG:unique.'],
                    ['country option', 'country:country:UG makes UG the value a new record starts with.'],
                    ['rating option', 'score:rating:10 gives the field 10 stars, from 2 to 10. The default is 5.'],
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
                    { id: 'decl', lane: 0, row: 2, title: 'price:money', sub: '--fields', tone: 'primary' },
                    { id: 'go', lane: 1, row: 0, title: 'Go struct field', sub: 'money.Money', tone: 'cyan' },
                    { id: 'ts', lane: 1, row: 1, title: 'TypeScript type', sub: 'Money', tone: 'blue' },
                    { id: 'zod', lane: 1, row: 2, title: 'Zod rule', sub: 'MoneySchema', tone: 'violet' },
                    { id: 'form', lane: 1, row: 3, title: 'Admin form', sub: 'amount + currency', tone: 'amber' },
                    { id: 'col', lane: 1, row: 4, title: 'Table column', sub: 'formatted, right-aligned', tone: 'green' },
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

              {/* Formatted types */}
              <div className="mb-12" id="formatted-types">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Formatted types</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  email, url, domain, tel, country, color, percent, rating, time and json are
                  stored as plain columns with a known shape. The API checks and normalises them
                  on every write that goes through a model: create, update, PATCH, bulk edit, CSV
                  import and sync push. The check lives in <code>internal/fieldtypes</code>, which
                  reads the <code>format</code> tag on the model field.
                </p>
                <CodeBlock
                  terminal
                  code={`grit generate resource Company --fields "name:string,email:email:unique,website:url,host:domain,phone:tel:UG,country:country,brand:color,discount:percent,score:rating:10,opens_at:time,settings:json" --seed --faker`}
                />
                <div className="overflow-x-auto mt-4 mb-4">
                  <table className="w-full text-sm border border-border rounded-lg">
                    <thead>
                      <tr className="border-b border-border bg-muted/30 text-left">
                        <th className="px-3 py-2 font-medium">Type</th>
                        <th className="px-3 py-2 font-medium">Stored as</th>
                        <th className="px-3 py-2 font-medium">Refused</th>
                      </tr>
                    </thead>
                    <tbody className="text-[13px]">
                      {FORMATTED.map((r) => (
                        <tr key={r.type} className="border-b border-border/50 align-top">
                          <td className="px-3 py-2 font-mono text-[12px] text-primary">{r.type}</td>
                          <td className="px-3 py-2 text-muted-foreground">{r.stored}</td>
                          <td className="px-3 py-2 text-muted-foreground">{r.refused}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  A refused value is a 422 in the usual error shape, with the field named in{' '}
                  <code>details</code>. The admin form checks the same rules before it sends, and
                  shows the API&apos;s message if one gets through:
                </p>
                <CodeBlock
                  code={`{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Phone is not a valid phone number for UG",
    "details": { "phone": "Phone is not a valid phone number for UG" }
  }
}`}
                />
                <p className="text-muted-foreground leading-relaxed mt-4 mb-4">
                  Phone numbers are validated with libphonenumber&apos;s metadata on both sides:{' '}
                  <code>github.com/nyaruka/phonenumbers</code> in the API and{' '}
                  <code>libphonenumber-js</code> in the browser, so the form and the API agree
                  about which numbers are valid. The Go module is added to a project the first
                  time it generates a tel field, not before. The country picker, shared by tel and
                  country, is a Base UI combobox: type a country&apos;s name, its code (UG) or its
                  calling code (256), move with the arrow keys and press Enter.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  With <code>--faker</code>, every formatted field seeds values its own rule
                  accepts, with the row number mixed in so a unique column seeds thousands of rows
                  without a collision: valid mobile numbers from twenty countries, addresses such as{' '}
                  <code>ada.lovelace.12@example.com</code>, ratings weighted towards four and five
                  stars, and times in business hours.
                </p>
                <p className="text-muted-foreground leading-relaxed">
                  <code>grit g field</code> adds any of them to an existing resource except json,
                  which needs <code>gorm.io/datatypes</code> in the model: regenerate the resource
                  for that one. The Expo form uses the phone keyboard and checks E.164 for tel,
                  and the matching keyboard for email, url and domain; the desktop forms use the
                  browser&apos;s own email, url, tel, time and colour inputs.
                </p>
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
