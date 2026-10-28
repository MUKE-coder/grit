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
  { type: 'select', syntax: 'status:select[draft, paid]', go: 'string', ts: '"draft" | "paid"', zod: 'z.enum([...])' },
  { type: 'radio', syntax: 'plan:radio[free, pro]', go: 'string', ts: '"free" | "pro"', zod: 'z.enum([...])' },
  { type: 'check', syntax: 'tags:check[news, ops]', go: 'datatypes.JSONSlice[string]', ts: '("news" | "ops")[]', zod: 'z.array(z.enum([...]))' },
  { type: 'datetime', syntax: 'published_at:datetime', go: '*time.Time', ts: 'string | null', zod: 'z.string().nullable()' },
  { type: 'date', syntax: 'due:date', go: '*time.Time', ts: 'string | null', zod: 'z.string().nullable()' },
  { type: 'slug', syntax: 'slug:slug[title]', go: 'string', ts: 'string', zod: 'z.string()' },
  { type: 'belongs_to', syntax: 'category:belongs_to[Category]', go: 'string (FK)', ts: 'string', zod: 'z.string().uuid()' },
  { type: 'many_to_many', syntax: 'tags:many_to_many[Tag]', go: '[]string', ts: 'string[]', zod: 'z.array(z.string().uuid())' },
  { type: 'string_array', syntax: 'sizes:string_array', go: 'datatypes.JSONSlice[string]', ts: 'string[]', zod: 'z.array(z.string())' },
  { type: 'file', syntax: 'image:file[image]', go: '*files.FileRef', ts: 'FileRef | null', zod: 'FileRefSchema.nullable()' },
  { type: 'files', syntax: 'gallery:files[image]', go: 'files.FileRefs', ts: 'FileRef[]', zod: 'z.array(FileRefSchema)' },
  { type: 'email', syntax: 'email:email:unique', go: 'string (size 254)', ts: 'string', zod: 'EmailSchema' },
  { type: 'url', syntax: 'website:url', go: 'string (size 2048)', ts: 'string', zod: 'UrlSchema' },
  { type: 'domain', syntax: 'host:domain', go: 'string (size 253)', ts: 'string', zod: 'DomainSchema' },
  { type: 'tel', syntax: 'phone:tel[UG]', go: 'string (size 32, E.164)', ts: 'string', zod: 'PhoneSchema' },
  { type: 'country', syntax: 'country:country[UG]', go: 'string (size 2)', ts: 'string', zod: 'CountrySchema' },
  { type: 'color', syntax: 'brand:color', go: 'string (size 7)', ts: 'string', zod: 'ColorSchema' },
  { type: 'percent', syntax: 'discount:percent', go: 'float64 (decimal 5,2)', ts: 'number', zod: 'PercentSchema' },
  { type: 'rating', syntax: 'score:rating[10]', go: 'int', ts: 'number', zod: 'ratingSchema(10)' },
  { type: 'time', syntax: 'opens_at:time', go: 'string (HH:MM)', ts: 'string', zod: 'TimeSchema' },
  { type: 'json', syntax: 'settings:json', go: 'datatypes.JSON', ts: 'unknown', zod: 'JsonValueSchema' },
]

// What each type becomes in the database, and the thing about it that is not
// obvious from the Go type. Source of truth: internal/generate/field.go
// (GORMTag) and internal/generate/templates.go (the tags beside it).
const COLUMNS: { type: string; tag: string; note: string }[] = [
  { type: 'string', tag: 'size:255', note: 'A name containing "url" or "link" gets size:500, and a long-form one such as "description" gets type:text. Guessed from the field name, because a 255-character limit on a tracking URL fails in production and not in testing.' },
  { type: 'text', tag: 'type:text', note: 'No length limit. Rendered as a textarea.' },
  { type: 'richtext', tag: 'type:text', note: 'Also gets sanitize:"html": it is rendered as HTML wherever it is shown, so every write is cleaned before it is stored.' },
  { type: 'int', tag: '(none)', note: 'A signed integer. Negative values are allowed.' },
  { type: 'uint', tag: '(none)', note: 'Refuses a negative at the Zod and binding layers, not only in the column, so the API answers 422 rather than storing a wrapped value.' },
  { type: 'float', tag: 'type:decimal(12,2)', note: 'Only when the name looks like money. A field called price or amount gets fixed precision, because 1.99 + 0.01 in float64 is not 2.00. Prefer the money type, which does not guess.' },
  { type: 'money', tag: '(two columns)', note: 'price_amount BIGINT in minor units plus price_currency. One column holding "19.99 USD" is not something you can SUM, and a float is not something you can trust.' },
  { type: 'bool / toggle', tag: '(none)', note: 'Same column. toggle is a switch in the form, bool is a checkbox.' },
  { type: 'select / radio', tag: 'size:255 + enum:"a|b"', note: 'The enum tag records the allowed values on the model, so a seeder generated later fills the column from the list instead of with dictionary words.' },
  { type: 'check', tag: 'type:json', note: 'Several of the allowed values, stored as a JSON array.' },
  { type: 'date', tag: 'type:date', note: 'A day with no time and no zone. datetime keeps both, as *time.Time.' },
  { type: 'slug', tag: 'size:255 uniqueIndex', note: 'Always unique: a slug that repeats is two URLs for one thing. Generated from the named field on write.' },
  { type: 'belongs_to', tag: 'size:36 index', note: 'Matches the parent UUID primary key, and indexed because every list filtered by it would otherwise scan.' },
  { type: 'one_to_one', tag: 'size:36 uniqueIndex', note: 'uniqueIndex rather than index is the whole difference. Without it the database accepts a second row pointing at the same parent and "one to one" is a comment rather than a constraint.' },
  { type: 'many_to_many', tag: '(join table)', note: 'GORM builds and owns the join table. No column on this model.' },
  { type: 'string_array', tag: 'type:json', note: 'Free text values rather than a closed set. Use check when the set is closed.' },
  { type: 'file / files', tag: 'type:json', note: 'A FileRef carries the key, name, size and type, not just a path. type:json means jsonb on Postgres, so it stays queryable.' },
  { type: 'email', tag: 'size:254 + format:"email"', note: 'The longest address SMTP allows. The format tag is what normalises and validates it on write.' },
  { type: 'url', tag: 'size:2048 + format:"url"', note: 'http and https only.' },
  { type: 'domain', tag: 'size:253 + format:"domain"', note: 'The longest name DNS allows. Stored as lowercase punycode.' },
  { type: 'tel', tag: 'size:32 + format:"tel"', note: 'E.164 with the plus, which is at most 16 characters; the rest is headroom.' },
  { type: 'country', tag: 'size:2 + format:"country"', note: 'ISO 3166-1 alpha-2.' },
  { type: 'color', tag: 'size:7 + format:"color"', note: '#rrggbb.' },
  { type: 'percent', tag: 'type:decimal(5,2)', note: '0.00 to 100.00 exactly, which a float column does not promise.' },
  { type: 'rating', tag: '(none)', note: 'An int bounded by the argument, enforced in Zod rather than in the column.' },
  { type: 'time', tag: 'size:5', note: 'HH:MM. A clock time with no date, for opening hours and the like.' },
  { type: 'json', tag: 'datatypes.JSON', note: 'Anything. Nothing validates its shape, which is the point and the risk.' },
]

// Modifiers that change the tags rather than the type.
const MODIFIER_TAGS: { modifier: string; tag: string; note: string }[] = [
  { modifier: ':unique', tag: 'uniqueIndex', note: 'On a slug it is already implied.' },
  { modifier: ':required', tag: 'binding:"required" + not null', note: 'String fields are required by default; everything else is optional by default.' },
  { modifier: ':optional', tag: 'default:null', note: 'Drops the binding tag and allows null.' },
  { modifier: 'if x == y', tag: 'requiredif:"x=y"', note: 'Required only when the other field has that value, checked where the row is, not at binding. Binding runs before anything has looked at the row, so it cannot evaluate a condition.' },
  { modifier: ':encrypted', tag: 'type:text', note: 'Ciphertext is longer than its plaintext, so a sized column would reject a value that fits.' },
  { modifier: ':searchable', tag: 'search:"trigram"', note: 'A trigram index, so a LIKE %term% stops being a scan.' },
]

// The flags, from `grit generate resource --help`.
const FLAGS: { flag: string; what: string }[] = [
  { flag: '--fields', what: 'The fields, inline. The alternative is --from with a YAML file, which is the same grammar one field per line and survives being read six months later.' },
  { flag: '--faker', what: 'Also writes a seeder that fills many rows with believable fake data, and implies --seed. Use it on anything you want to look at a real list of; an empty table hides every layout problem. --count sets how many (10 by default).' },
  { flag: '--seed', what: 'Writes a seeder with one example record. Enough to prove the resource works, not enough to see a table under load.' },
  { flag: '--public', what: 'Also exposes read-only list and detail endpoints under /api/v1/public/, behind an API key. Read-only is the whole contract: a public endpoint that writes is hand-written, and must not go on this group, which is URL-cached.' },
  { flag: '--items', what: 'An inline has-many child, rendered as a line-items table inside this resource\u2019s form rather than as its own screen. An invoice with its lines, an order with its items: the parent is meaningless without them and nobody wants two screens to enter one document.' },
  { flag: '--owned-by', what: 'Scopes the resource to its owner: the named belongs_to-User field decides who may read or write each row, and the field is added if it is absent. ADMIN is exempt. This is enforced in the service, not in the handler, so it holds for every route.' },
  { flag: '--append-only', what: 'Rows are created and read, never changed or deleted. No update or delete routes, a GORM guard, and a database trigger installed by grit migrate, so the rule survives somebody with a SQL client. For ledgers, audit trails and consent records.' },
  { flag: '--audit-reads', what: 'Records every read in the tamper-evident activity log: who read which rows, when, from where. For data where access itself must be accountable.' },
  { flag: '--roles', what: 'Restricts the routes to named roles.' },
  { flag: '--tree', what: 'Makes the resource hierarchical: a self-referential parent, a materialized path, depth and sibling order, the tree queries and a move endpoint.' },
  { flag: '--tenant-owned', what: 'Scopes the resource to an organization, with the multitenant plugin: queries are filtered by the active org and OrgID is stamped on insert.' },
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
              <span className="tag-mono text-accent/80 mb-3 block">Core Concepts · Reference</span>
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
                  code={`grit generate resource Product --fields "name:string,price:money,category:belongs_to[Category],gallery:files[image]"`}
                />

                <p className="text-muted-foreground leading-relaxed mt-6 mb-4">
                  A type that takes an argument takes it in brackets, and a field that only
                  applies sometimes says so in words. Read one left to right and it tells you
                  everything about itself.
                </p>
                <CodeBlock
                  code={`cycle : select[monthly, annually] if kind == recurring
  │         │             │                    │
  │         │             │                    └─ when this field applies
  │         │             └─ this type's argument: the options
  │         └─ the type
  └─ the name: the column, and the JSON key`}
                />
                <p className="text-muted-foreground leading-relaxed mt-4 mb-4">
                  Labels are derived from the values, so <code>monthly</code> renders as
                  &quot;Monthly&quot; and <code>in_progress</code> as &quot;In Progress&quot;.
                  Give one explicitly only when it differs:{' '}
                  <code>[monthly, annually=&quot;Every year&quot;]</code>.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Only some types take an argument: <code>select</code>, <code>radio</code> and{' '}
                  <code>check</code> take their options, <code>slug</code> takes the field to
                  slugify, <code>belongs_to</code>, <code>one_to_one</code> and{' '}
                  <code>many_to_many</code> take the related model, <code>file</code> and{' '}
                  <code>files</code> take what they accept, <code>tel</code> and{' '}
                  <code>country</code> take a country code, and <code>rating</code> takes a number
                  of stars. Every other type goes straight from its name to its modifiers, which
                  are appended with colons: <code>email:email:unique</code>.
                </p>

                <CodeBlock
                  terminal
                  code={`grit generate resource Project --fields "name:string, kind:select[one_off, recurring], cycle:select[monthly, annually] if kind == recurring, owner:belongs_to[User]"`}
                />

                <p className="text-muted-foreground leading-relaxed mt-4 mb-4">
                  <strong>Name the argument, even when it can be guessed.</strong>{' '}
                  <code>client:belongs_to</code> infers the model from the field name and works;{' '}
                  <code>client:belongs_to[Client]</code> says which model, and keeps saying it
                  when somebody renames the field. The same goes for the upload types:{' '}
                  <code>photo:file[image]</code> and <code>gallery:files[image]</code> declare
                  what is accepted, where a bare <code>file</code> accepts anything. Every
                  example in these docs is written out in full for that reason.
                </p>

                <Callout type="note" title="The older positional spelling">
                  It still parses, and older tutorials use it:{' '}
                  <code>cycle:select:monthly=Monthly|annually=Annually:when=kind=recurring</code>.
                  The two produce identical fields, and there is one parser: the readable form is
                  rewritten into the positional one before anything else sees it.
                  <br />
                  <br />
                  Prefer the readable form in anything new. In the positional one the colon means
                  four different things, and the pipe separating options is a shell pipe, so an
                  unquoted example silently loses half of them.
                </Callout>

                <ul className="space-y-2.5 mt-4 mb-4">
                  {[
                    ['Modifiers', 'name:string:unique — append :unique, :required, :optional, :encrypted or :searchable. They always come last, after whatever argument the type takes: phone:tel[UG]:unique, category:belongs_to[Category]:optional. A condition is written if kind == recurring and goes last of all; :when=kind=recurring is the older spelling of the same thing. String fields default to required; everything else defaults to optional.'],
                    ['Conditional', 'cycle:select[monthly, annually] if kind == recurring — the field applies only while another field of the same resource holds a particular value. kind is a sibling field, named as it appears in the spec; recurring is one of its stored values, not a label. The admin hides the input and drops its required rule when the condition is false; the model gets requiredif:"kind=recurring", which the API enforces on every write. See conditional fields.'],
                    ['Auto-number', 'number:string:auto:INV (string only) — the server fills it from an atomic, gap-free counter (INV-202607-0001) in BeforeCreate. Optional and hidden from the form; the prefix is optional. Shorthand for grit generate sequence; use that command directly for yearly/never resets or a custom width.'],
                    ['Encrypted', 'notes:text:encrypted (string, text and richtext only). Stored with AES-256-GCM at rest, plaintext in code and over the API; needs FIELD_ENCRYPTION_KEY. Cannot be :unique, and is left out of search, sorting and filters, because ciphertext changes on every write.'],
                    ['Options', 'status:select[draft, sent, paid] (select / radio / check). Labels are derived by capitalising each word (in_progress → In Progress); give one only when it differs: [draft, sent="Sent to client"]. The positional spelling status:select:draft=Draft|sent=Sent is the older form of the same thing.'],
                    ['Slug source', 'slug:slug[title] — the argument is the field to slugify. Slugs are auto-unique and generated on save.'],
                    ['belongs_to', 'author:belongs_to[User] names the model, and is the form to prefer. category:belongs_to also works, inferring Category from the field name, until somebody renames the field. Creates a <name>_id UUID foreign-key column.'],
                    ['many_to_many', 'tags:many_to_many[Tag] — the related model is required. GORM builds the join table.'],
                    ['File accepts', 'image:file[image], doc:file[all], or att:files[pdf, doc, image]. Aliases: image, video, audio, pdf, doc, all — or a bracketed list.'],
                    ['Money', 'price:money — stores an integer count of minor units plus an ISO 4217 currency code, in two columns (price_amount, price_currency). Use it for anything you will add up. See the money page.'],
                    ['tel option', 'phone:tel[UG] sets the country the picker starts on and a local number is read in. Without it the picker starts on the browser\'s country. Modifiers follow the option: phone:tel[UG]:unique.'],
                    ['country option', 'country:country[UG] makes UG the value a new record starts with.'],
                    ['rating option', 'score:rating[10] gives the field 10 stars, from 2 to 10. The default is 5.'],
                  ].map(([k, v]) => (
                    <li key={k} className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="text-accent mt-1 font-mono text-xs shrink-0">{k}</span>
                      <span>{v}</span>
                    </li>
                  ))}
                </ul>

                <h3 className="text-lg font-semibold tracking-tight mt-8 mb-3">
                  When is a field conditional, rather than just optional?
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Optional means the answer may be blank. Conditional means the question does not
                  apply. The test: if you asked it about a record where the condition is false,
                  would the answer be <em>meaningless</em> rather than merely empty? A billing
                  cycle on a one-off project is meaningless, so it is conditional. A phone number
                  nobody gave you is empty, so it is optional.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Two things to get right in <code>when=kind=recurring</code>. The field name is a
                  sibling on the same resource, written as it appears in the spec, and it must
                  exist, or the command is refused. The value is the <em>stored</em> value, which
                  for a select is the left-hand side of <code>recurring=Recurring</code>, not the
                  label; it is compared as text, so case matters.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  It works on any field type, including a relationship: a plan that only applies to
                  a recurring project is{' '}
                  <code>plan:belongs_to[Plan]:when=kind=recurring</code>. For a condition that is
                  more than one field equalling one value, write the predicate by hand in the
                  resource definition and the rule in the service.
                </p>
              </div>

              {/* Master table */}
              <div className="mb-12">
                <h2 id="columns-and-tags" className="text-2xl font-semibold tracking-tight mb-4">
                  Columns and tags
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  What each type becomes in the database, and the thing about it that the Go type
                  does not tell you. These are the tags the generator writes; a model you edit by
                  hand keeps whatever you put there.
                </p>
                <div className="overflow-x-auto mb-8">
                  <table className="w-full text-sm table-fixed">
                    <colgroup>
                      <col className="w-[13%]" />
                      <col className="w-[25%]" />
                      <col className="w-[62%]" />
                    </colgroup>
                    <thead>
                      <tr className="border-b border-border/40 text-left">
                        <th className="px-3 py-2 font-medium">Type</th>
                        <th className="px-3 py-2 font-medium">Tag</th>
                        <th className="px-3 py-2 font-medium">Note</th>
                      </tr>
                    </thead>
                    <tbody>
                      {COLUMNS.map((row) => (
                        <tr key={row.type} className="border-b border-border/20 align-top">
                          <td className="px-3 py-2 font-mono text-xs text-accent break-words">{row.type}</td>
                          <td className="px-3 py-2 font-mono text-xs text-muted-foreground break-words">{row.tag}</td>
                          <td className="px-3 py-2 text-muted-foreground leading-relaxed">{row.note}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                <h3 className="text-xl font-semibold tracking-tight mb-3">Modifiers, and the tags they add</h3>
                <div className="overflow-x-auto">
                  <table className="w-full text-sm table-fixed">
                    <colgroup>
                      <col className="w-[18%]" />
                      <col className="w-[24%]" />
                      <col className="w-[58%]" />
                    </colgroup>
                    <thead>
                      <tr className="border-b border-border/40 text-left">
                        <th className="px-3 py-2 font-medium">Modifier</th>
                        <th className="px-3 py-2 font-medium">Tag</th>
                        <th className="px-3 py-2 font-medium">Note</th>
                      </tr>
                    </thead>
                    <tbody>
                      {MODIFIER_TAGS.map((row) => (
                        <tr key={row.modifier} className="border-b border-border/20 align-top">
                          <td className="px-3 py-2 font-mono text-xs text-accent break-words">{row.modifier}</td>
                          <td className="px-3 py-2 font-mono text-xs text-muted-foreground break-words">{row.tag}</td>
                          <td className="px-3 py-2 text-muted-foreground leading-relaxed">{row.note}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              <div className="mb-12">
                <h2 id="flags" className="text-2xl font-semibold tracking-tight mb-4">
                  The flags
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The fields say what the resource holds. These say what it is for, and most of
                  them change the generated code in more than one place at once.
                </p>
                <div className="overflow-x-auto">
                  <table className="w-full text-sm table-fixed">
                    <colgroup>
                      <col className="w-[22%]" />
                      <col className="w-[78%]" />
                    </colgroup>
                    <thead>
                      <tr className="border-b border-border/40 text-left">
                        <th className="px-3 py-2 font-medium">Flag</th>
                        <th className="px-3 py-2 font-medium">What it does</th>
                      </tr>
                    </thead>
                    <tbody>
                      {FLAGS.map((row) => (
                        <tr key={row.flag} className="border-b border-border/20 align-top">
                          <td className="px-3 py-2 font-mono text-xs text-accent break-words">{row.flag}</td>
                          <td className="px-3 py-2 text-muted-foreground leading-relaxed">{row.what}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

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
                          <td className="px-3 py-2 text-accent">{r.type}</td>
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
                          <td className="px-4 py-2 font-mono text-[12px] text-accent">{r.type}</td>
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
                  code={`grit generate resource Company --fields "name:string,email:email:unique,website:url,host:domain,phone:tel[UG],country:country,brand:color,discount:percent,score:rating[10],opens_at:time,settings:json" --seed --faker`}
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
                          <td className="px-3 py-2 font-mono text-[12px] text-accent">{r.type}</td>
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
