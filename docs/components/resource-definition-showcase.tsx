'use client'

import { useState } from 'react'
import { CodeBlock } from '@/components/code-block'

/**
 * The resource definition, as an actual reference.
 *
 * Everything here comes from the ResourceDefinition interfaces the scaffold
 * emits into apps/admin/lib/resource.ts. The field-type and column-format
 * lists are complete rather than illustrative — a partial list on a page
 * titled "every option" is worse than no list, because readers stop looking
 * once they trust it.
 *
 * When the interfaces change, update the counts in the headings too. They are
 * hardcoded on purpose: a reader can check them against the list beside them.
 */

const TABS = [
  { key: 'shape', label: 'The shape' },
  { key: 'table', label: 'Table' },
  { key: 'form', label: 'Form' },
  { key: 'uploads', label: 'Upload fields' },
] as const

type TabKey = (typeof TABS)[number]['key']

const SHAPE_CODE = `import { defineResource } from "@/lib/resource";

export const productResource = defineResource({
  name: "Product",
  slug: "products",
  endpoint: "/api/products",
  icon: "Package",
  label: { singular: "Product", plural: "Products" },

  // sheet · modal · page · modal-steps · page-steps
  formView: "sheet",

  // Sidebar placement
  group: "Catalog",     // section heading to sit under
  adminOnly: false,     // hide from non-admins
  hidden: false,        // keep routable, drop from the sidebar

  table: { /* columns, filters, actions, export, import */ },
  form:  { /* fields, layout, steps, groups */ },

  // Four auto stat cards, or your own
  stats: true,
  dashboard: { enabled: true },
});`

const TABLE_CODE = `table: {
  columns: [
    { key: "name",  label: "Name",  sortable: true, searchable: true,
      onClick: "link" },
    { key: "price", label: "Price", format: "currency",
      currencyPrefix: "$" },
    { key: "status", label: "Status", format: "badge",
      badge: {
        active: { color: "emerald", label: "Active" },
        draft:  { color: "amber",   label: "Draft" },
      } },
    { key: "cover", label: "Cover", format: "image" },
    { key: "owner", label: "Owner", format: "user" },
    // Or take over the cell entirely:
    { key: "score", label: "Score",
      cell: (row) => <Sparkline value={row.score} /> },
  ],

  filters: [
    { key: "status", label: "Status", type: "select", options: [...] },
    { key: "created_at", label: "Created", type: "date-range" },
    { key: "price", label: "Price", type: "number-range" },
  ],

  rowActions: [
    { label: "Duplicate", onClick: (row) => duplicate(row) },
    { label: "Archive", variant: "danger",
      visible: (row) => !row.archived },
  ],

  actions: ["create", "view", "edit", "delete", "export"],
  bulkActions: ["delete", "export"],
  defaultSort: { key: "created_at", direction: "desc" },
  pageSize: 20,
  searchable: true,
  dateFilter: { enabled: true, field: "created_at", label: "Created" },
  export: { csv: true, json: true, excel: true, allPages: true },
  import: { excel: true, fields: ["name", "price", "status"] },
}`

const FORM_CODE = `form: {
  layout: "two-column",        // or "single"
  sheetWidth: "wide",          // "half" (50%) or "wide" (80%)

  fields: [
    { key: "name", label: "Name", type: "text", required: true,
      colSpan: 2, placeholder: "Aurora 14\\" Studio" },

    { key: "category", label: "Category",
      type: "relationship-select",
      relatedEndpoint: "/api/categories",
      displayField: "name",
      allowCreate: true },           // adds "New Category"

    { key: "price", label: "Price", type: "number",
      numberKind: "float", prefix: "$", min: 0, step: 0.01 },

    { key: "launch_on", label: "Launch", type: "date",
      minDate: "2020-01-01", maxDate: "2030-12-31" },

    { key: "items", label: "Line items", type: "line-items",
      itemFields: [...], foreignKey: "invoice_id" },
  ],

  // Break it into a wizard
  steps: [
    { title: "Basics",  fields: ["name", "category", "price"] },
    { title: "Media",   fields: ["cover", "gallery"] },
  ],
  stepVariant: "horizontal",
  perStepSave: true,   // on edit, each step saves on its own
}`

const UPLOAD_CODE = `// Every knob a file field takes

{ key: "gallery", label: "Gallery", type: "files",

  // What it takes. Aliases, not MIME strings.
  //   image · video · audio · pdf · doc · excel · csv · zip · archive · all
  accepts: ["image"],

  // Cap per file. Defaults to 5 MB — 300 MB when the field takes video.
  maxSizeMB: 10,

  // How the dropzone looks
  //   default · compact · minimal · avatar · inline
  dropzone: "default",

  // How progress reads while uploading
  //   bar (default) · circular · pulse
  progress: "circular",

  // Multi-file only: drag to reorder. On by default.
  reorderable: true,
}

// Single PDF, tight layout, no reordering to worry about
{ key: "spec_sheet", label: "Spec sheet", type: "file",
  accepts: ["pdf"], dropzone: "compact" }

// Avatar: a round target that replaces rather than appends
{ key: "avatar", label: "Avatar", type: "image",
  accepts: ["image"], dropzone: "avatar", maxSizeMB: 2 }

// Mixed downloads
{ key: "downloads", label: "Downloads", type: "files",
  accepts: ["zip", "doc"], progress: "bar" }`

const FIELD_TYPES = [
  'text', 'textarea', 'number', 'select', 'date', 'datetime', 'toggle', 'checkbox',
  'checkbox-group', 'radio', 'richtext', 'image', 'images', 'video', 'videos', 'file',
  'files', 'relationship-select', 'multi-relationship-select', 'line-items',
]

const COLUMN_FORMATS = [
  'text', 'badge', 'currency', 'date', 'relative', 'boolean', 'image', 'video', 'file',
  'files', 'link', 'email', 'color', 'richtext', 'user',
]

const DROPZONE_VARIANTS = ['default', 'compact', 'minimal', 'avatar', 'inline']
const PROGRESS_STYLES = ['bar', 'circular', 'pulse']
const ACCEPT_ALIASES = ['image', 'video', 'audio', 'pdf', 'doc', 'excel', 'csv', 'zip', 'archive', 'all']

function Chips({ title, items, mono = true }: { title: string; items: string[]; mono?: boolean }) {
  return (
    <div className="mb-5 last:mb-0">
      <div className="text-[10.5px] font-mono uppercase tracking-wider text-muted-foreground/70 mb-2.5">
        {title} <span className="text-muted-foreground/50">({items.length})</span>
      </div>
      <div className="flex flex-wrap gap-1.5">
        {items.map((i) => (
          <span
            key={i}
            className={`rounded-md border border-border/60 bg-background/60 px-2 py-1 text-[11px] text-foreground/80 ${
              mono ? 'font-mono' : ''
            }`}
          >
            {i}
          </span>
        ))}
      </div>
    </div>
  )
}

export function ResourceDefinitionShowcase() {
  const [active, setActive] = useState<TabKey>('shape')

  const code =
    active === 'shape' ? SHAPE_CODE :
    active === 'table' ? TABLE_CODE :
    active === 'form' ? FORM_CODE : UPLOAD_CODE

  const file =
    active === 'uploads'
      ? 'apps/admin/resources/products.ts — form.fields'
      : 'apps/admin/resources/products.ts'

  return (
    <div>
      <div role="tablist" aria-label="Resource definition" className="flex flex-wrap gap-2 mb-8">
        {TABS.map((t) => {
          const selected = t.key === active
          return (
            <button
              key={t.key}
              type="button"
              role="tab"
              aria-selected={selected}
              onClick={() => setActive(t.key)}
              className={`rounded-full border px-4 py-2 text-sm font-medium transition-colors ${
                selected
                  ? 'border-primary/40 bg-primary/10 text-foreground'
                  : 'border-border/60 text-muted-foreground hover:text-foreground hover:border-border'
              }`}
            >
              {t.label}
            </button>
          )
        })}
      </div>

      <div className="grid lg:grid-cols-[1fr_19rem] gap-8 lg:gap-10 items-start">
        <div className="rounded-xl overflow-hidden border border-border bg-card/40">
          <div className="flex items-center gap-2 px-4 py-2.5 bg-card/60 border-b border-border/60">
            <span className="text-[11.5px] font-mono text-muted-foreground truncate">{file}</span>
            <span className="ml-auto text-[10px] font-mono uppercase tracking-wider text-muted-foreground/60">
              generated, then yours
            </span>
          </div>
          <CodeBlock
            key={active}
            code={code}
            language="tsx"
            className="!border-0 !rounded-none !shadow-none !bg-transparent dark:!bg-transparent !m-0"
          />
        </div>

        <div className="rounded-xl border border-border/60 bg-card/40 p-4">
          {active === 'table' && (
            <>
              <Chips title="Column formats" items={COLUMN_FORMATS} />
              <Chips title="Filter types" items={['select', 'date-range', 'number-range', 'boolean']} />
              <p className="text-[11.5px] text-muted-foreground/80 leading-relaxed">
                A column can also take <code className="text-foreground/70">cell</code> and render
                anything you like &mdash; <code className="text-foreground/70">format</code> is the
                shortcut, not the ceiling.
              </p>
            </>
          )}
          {active === 'form' && (
            <>
              <Chips title="Field types" items={FIELD_TYPES} />
              <p className="text-[11.5px] text-muted-foreground/80 leading-relaxed">
                Fields carry the usual validation and layout props, plus{' '}
                <code className="text-foreground/70">numberKind</code> for int/uint/float,{' '}
                <code className="text-foreground/70">optionsUrl</code> for remote selects, and{' '}
                <code className="text-foreground/70">generate</code> for a fill-this-for-me button.
              </p>
            </>
          )}
          {active === 'uploads' && (
            <>
              <Chips title="Dropzone variants" items={DROPZONE_VARIANTS} />
              <Chips title="Progress styles" items={PROGRESS_STYLES} />
              <Chips title="Accept aliases" items={ACCEPT_ALIASES} />
              <p className="text-[11.5px] text-muted-foreground/80 leading-relaxed">
                Accept aliases are enforced on the server as well as in the browser, so a field
                declared <code className="text-foreground/70">pdf</code> refuses a PNG even if the
                request is hand-made.
              </p>
            </>
          )}
          {active === 'shape' && (
            <>
              <Chips
                title="Form views"
                items={['sheet', 'modal', 'page', 'modal-steps', 'page-steps']}
              />
              <Chips title="Table actions" items={['create', 'view', 'edit', 'delete', 'export']} />
              <p className="text-[11.5px] text-muted-foreground/80 leading-relaxed">
                <code className="text-foreground/70">grit generate resource</code> writes this file
                for you. Everything past that point is ordinary TypeScript you own &mdash; there is
                no regeneration step that overwrites your edits.
              </p>
            </>
          )}
        </div>
      </div>
    </div>
  )
}
