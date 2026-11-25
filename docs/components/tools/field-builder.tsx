'use client'

import { useMemo, useState } from 'react'
import { Check, Copy, Plus, Trash2 } from 'lucide-react'
import { FIELD_PREVIEWS, shapeFor } from '@/config/field-preview.generated'

/**
 * Build a --fields string without memorising the syntax.
 *
 * Every type and every spelling here is taken from the generator's own
 * parser: the 32 constants in internal/generate/field.go, the modifiers
 * applyFieldModifiers accepts, and the third-position arguments that are NOT
 * modifiers (a related model, a dimension, an accept list, a slug source).
 *
 * Each shape this emits was run through ParseInlineFields before the tool was
 * written, because a builder that emits a string the CLI rejects is worse
 * than no builder: it is wrong with a confident interface on top. The test
 * that keeps the documentation honest about this lives in
 * internal/generate/docs_fields_test.go.
 */

type Extra = 'none' | 'options' | 'model' | 'dimension' | 'accept' | 'source'

interface TypeDef {
  value: string
  label: string
  group: string
  /** What the third position means for this type, if anything. */
  extra: Extra
  hint?: string
}

const TYPES: TypeDef[] = [
  // Text
  { value: 'string', label: 'string', group: 'Text', extra: 'none', hint: 'A single line' },
  { value: 'text', label: 'text', group: 'Text', extra: 'none', hint: 'Several lines' },
  { value: 'richtext', label: 'richtext', group: 'Text', extra: 'none', hint: 'A formatted editor' },
  { value: 'slug', label: 'slug', group: 'Text', extra: 'source', hint: 'Derived from another field' },
  { value: 'email', label: 'email', group: 'Text', extra: 'none' },
  { value: 'url', label: 'url', group: 'Text', extra: 'none' },
  { value: 'domain', label: 'domain', group: 'Text', extra: 'none' },
  { value: 'tel', label: 'tel', group: 'Text', extra: 'none' },
  { value: 'color', label: 'color', group: 'Text', extra: 'none' },
  { value: 'country', label: 'country', group: 'Text', extra: 'none' },

  // Numbers
  { value: 'int', label: 'int', group: 'Numbers', extra: 'none' },
  { value: 'uint', label: 'uint', group: 'Numbers', extra: 'none', hint: 'Never negative' },
  { value: 'float', label: 'float', group: 'Numbers', extra: 'none' },
  { value: 'money', label: 'money', group: 'Numbers', extra: 'none', hint: 'Minor units and a currency' },
  { value: 'percent', label: 'percent', group: 'Numbers', extra: 'none' },
  { value: 'rating', label: 'rating', group: 'Numbers', extra: 'none' },

  // Truth and time
  { value: 'bool', label: 'bool', group: 'Truth and time', extra: 'none' },
  { value: 'toggle', label: 'toggle', group: 'Truth and time', extra: 'none', hint: 'A bool, shown as a switch' },
  { value: 'date', label: 'date', group: 'Truth and time', extra: 'none' },
  { value: 'datetime', label: 'datetime', group: 'Truth and time', extra: 'none' },
  { value: 'time', label: 'time', group: 'Truth and time', extra: 'none' },

  // Choices
  { value: 'select', label: 'select', group: 'Choices', extra: 'options', hint: 'One of a list' },
  { value: 'radio', label: 'radio', group: 'Choices', extra: 'options', hint: 'One of a few' },
  { value: 'check', label: 'check', group: 'Choices', extra: 'options', hint: 'Several of a list' },

  // Relationships
  { value: 'belongs_to', label: 'belongs_to', group: 'Relationships', extra: 'model', hint: 'This row points at one other' },
  { value: 'one_to_one', label: 'one_to_one', group: 'Relationships', extra: 'model' },
  { value: 'many_to_many', label: 'many_to_many', group: 'Relationships', extra: 'model', hint: 'A join table' },

  // Files and structured data
  { value: 'file', label: 'file', group: 'Files and data', extra: 'accept', hint: 'One upload' },
  { value: 'files', label: 'files', group: 'Files and data', extra: 'accept', hint: 'Several uploads' },
  { value: 'string_array', label: 'string_array', group: 'Files and data', extra: 'none' },
  { value: 'json', label: 'json', group: 'Files and data', extra: 'none' },
  { value: 'vector', label: 'vector', group: 'Files and data', extra: 'dimension', hint: 'An embedding' },
]

const GROUPS = Array.from(new Set(TYPES.map((t) => t.group)))

const ACCEPTS = ['image', 'pdf', 'doc', 'video', 'audio', 'zip', 'all']

interface Row {
  key: number
  name: string
  type: string
  extra: string
  required: boolean
  unique: boolean
}

/**
 * snake_case to PascalCase, the way the generator names a struct field.
 *
 * Only for the Go line in the preview. It is the one place the field's own
 * name appears in generated code on this page, and "due_at string" reads as
 * something the model would never contain.
 */
function pascal(name: string): string {
  return name
    .trim()
    .split(/[_\s-]+/)
    .filter(Boolean)
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join('')
}

let nextKey = 1
function row(name = '', type = 'string', extra = ''): Row {
  return { key: nextKey++, name, type, extra, required: false, unique: false }
}

const PRESETS: { label: string; resource: string; rows: Row[] }[] = [
  {
    label: 'Blog post',
    resource: 'Post',
    rows: [
      { ...row('title', 'string'), required: true },
      row('slug', 'slug', 'title'),
      row('body', 'richtext'),
      { ...row('author', 'belongs_to', 'User'), required: true },
      row('tags', 'many_to_many', 'Tag'),
      row('cover', 'file', 'image'),
      row('published_at', 'datetime'),
    ],
  },
  {
    label: 'Invoice',
    resource: 'Invoice',
    rows: [
      { ...row('number', 'string'), required: true, unique: true },
      { ...row('customer', 'belongs_to', 'Customer'), required: true },
      { ...row('total', 'money'), required: true },
      row('status', 'select', 'draft, sent, paid'),
      row('due_at', 'date'),
      row('notes', 'text'),
    ],
  },
  {
    label: 'Product',
    resource: 'Product',
    rows: [
      { ...row('name', 'string'), required: true },
      row('sku', 'string'),
      { ...row('price', 'money'), required: true },
      row('photos', 'files', 'image'),
      row('in_stock', 'toggle'),
      row('category', 'belongs_to', 'Category'),
    ],
  },
]

/** The field spec for one row, in the generator's syntax. */
function spec(r: Row): string {
  const def = TYPES.find((t) => t.value === r.type)
  const name = r.name.trim()
  if (!name) return ''

  let out = `${name}:${r.type}`
  const extra = r.extra.trim()

  switch (def?.extra) {
    case 'options':
      // select[draft, sent, paid]: the readable bracket form, which the
      // parser rewrites into the positional one.
      if (extra) out = `${name}:${r.type}[${extra}]`
      break
    case 'dimension':
      // Required, and part of the column type rather than a modifier.
      if (extra) out = `${name}:${r.type}[${extra}]`
      break
    case 'model':
    case 'accept':
    case 'source':
      // A third positional segment, not a modifier.
      if (extra) out = `${name}:${r.type}:${extra}`
      break
    default:
      break
  }

  if (r.required) out += ':required'
  if (r.unique) out += ':unique'
  return out
}

function problemsFor(r: Row): string | null {
  const def = TYPES.find((t) => t.value === r.type)
  const name = r.name.trim()
  if (!name) return 'This field has no name yet.'
  if (!/^[a-z][a-z0-9_]*$/.test(name)) {
    return 'Use snake_case: a lowercase letter first, then letters, digits or underscores.'
  }
  if (def?.extra === 'dimension' && !r.extra.trim()) {
    return 'A vector needs its dimension, such as 1536. It is part of the column type, so changing it later means a migration and re-embedding every row.'
  }
  if (def?.extra === 'model' && !r.extra.trim()) {
    return 'Name the related model, such as User.'
  }
  if (def?.extra === 'options' && !r.extra.trim()) {
    return 'List the choices, separated by commas.'
  }
  if (def?.extra === 'source' && !r.extra.trim()) {
    return 'Name the field to build the slug from, such as title.'
  }
  return null
}


/**
 * What one field turns into, across the stack.
 *
 * The question the builder could not answer before was the one people
 * actually have, which is not "what is the syntax" but "what do I get". A
 * money field is two columns and never a float; a slug is filled by the
 * server and never appears on the form; a check is a JSON column rather than
 * a join table. Those are the things somebody finds out after generating,
 * and they are why the command gets run twice.
 *
 * Every value here comes from config/field-preview.generated.ts, which is
 * built from Field.GoType, GORMTag, TSType, ZodType and FormFieldType in
 * internal/generate. It is not a table somebody wrote next to the generator
 * and hoped would stay true: TestFieldPreviewDataIsCurrent fails when it
 * drifts.
 */
function Produces({ row }: { row: Row }) {
  const p = FIELD_PREVIEWS[row.type]
  if (!p) return null

  // The modifiers change the answer, so the preview has to follow them.
  // required drops the .optional() and adds .min(1), unique adds
  // uniqueIndex to the column. Showing the unmodified shape beside a
  // ticked box is a confident interface being wrong.
  const shape = shapeFor(p, { required: row.required, unique: row.unique })

  const cells: { label: string; value: string }[] = [
    {
      label: 'Go',
      value:
        (row.name.trim() !== '' ? pascal(row.name) + ' ' : '') +
        shape.go +
        (shape.gorm !== '' ? '  `gorm:"' + shape.gorm + '"`' : ''),
    },
    { label: 'TypeScript', value: shape.ts },
    { label: 'Zod', value: shape.zod },
    // A slug has no control because the server fills it, and saying so is
    // more useful than leaving a gap the reader has to interpret.
    { label: 'Admin', value: shape.control !== '' ? shape.control : 'not on the form' },
  ]

  return (
    <div className="mt-3 rounded-lg border border-border/40 bg-background/50 p-3">
      <p className="mb-2 text-[11px] font-medium uppercase tracking-wider text-muted-foreground/80">
        What this produces
      </p>
      <dl className="grid gap-x-6 gap-y-1.5 sm:grid-cols-2">
        {cells.map((c) => (
          <div key={c.label} className="flex min-w-0 items-baseline gap-2">
            <dt className="w-[5.5rem] shrink-0 text-[11px] text-muted-foreground">{c.label}</dt>
            <dd className="min-w-0 break-all font-mono text-[12px] text-foreground">{c.value}</dd>
          </div>
        ))}
      </dl>
      {p.note && (
        <p className="mt-2 border-t border-border/40 pt-2 text-[12px] leading-relaxed text-muted-foreground">
          {p.note}
        </p>
      )}
    </div>
  )
}

export function FieldBuilder() {
  const [resource, setResource] = useState('Invoice')
  const [rows, setRows] = useState<Row[]>(PRESETS[1].rows)
  const [copied, setCopied] = useState(false)

  const update = (key: number, patch: Partial<Row>) =>
    setRows((rs) => rs.map((r) => (r.key === key ? { ...r, ...patch } : r)))

  const fields = useMemo(() => rows.map(spec).filter(Boolean).join(','), [rows])
  const problems = useMemo(
    () => rows.map((r) => ({ key: r.key, message: problemsFor(r) })).filter((p) => p.message),
    [rows],
  )

  const command = useMemo(() => {
    const name = resource.trim() || 'Resource'
    return `grit generate resource ${name} \\\n  --fields "${fields}"`
  }, [resource, fields])

  async function copy() {
    try {
      await navigator.clipboard.writeText(command)
      setCopied(true)
      setTimeout(() => setCopied(false), 1800)
    } catch {
      // A clipboard write can be refused (no permission, an insecure origin).
      // The command is on screen and selectable, so there is nothing to
      // recover: say nothing rather than showing a tick that is not true.
      setCopied(false)
    }
  }

  return (
    <div className="grid gap-10 lg:grid-cols-[minmax(0,1fr)_380px]">
      <div>
        {/* Presets */}
      <div className="mb-8 flex flex-wrap items-center gap-2">
        <span className="mr-1 text-sm text-muted-foreground">Start from:</span>
        {PRESETS.map((p) => (
          <button
            key={p.label}
            type="button"
            onClick={() => {
              setResource(p.resource)
              setRows(p.rows.map((r) => ({ ...r, key: nextKey++ })))
            }}
            className="rounded-lg border border-border/60 px-3 py-1.5 text-sm transition-colors hover:border-accent/60 hover:text-foreground"
          >
            {p.label}
          </button>
        ))}
        <button
          type="button"
          onClick={() => {
            setResource('')
            setRows([row()])
          }}
          className="rounded-lg border border-border/60 px-3 py-1.5 text-sm text-muted-foreground transition-colors hover:border-accent/60 hover:text-foreground"
        >
          Empty
        </button>
      </div>

      {/* Resource name */}
      <div className="mb-6">
        <label htmlFor="resource" className="mb-1.5 block text-sm font-medium text-foreground">
          Resource name
        </label>
        <input
          id="resource"
          value={resource}
          onChange={(e) => setResource(e.target.value)}
          placeholder="Invoice"
          className="w-full max-w-xs rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent/30"
        />
        <p className="mt-1.5 text-xs text-muted-foreground">
          Singular and PascalCase. The generator produces the plural table, route and type names
          itself.
        </p>
      </div>

      {/* Fields */}
      <div className="mb-4 space-y-3">
        {rows.map((r) => {
          const def = TYPES.find((t) => t.value === r.type)
          const problem = problemsFor(r)
          return (
            <div
              key={r.key}
              className="rounded-xl border border-border/60 bg-card/30 p-4"
            >
              <div className="flex flex-wrap items-end gap-3">
                <div className="min-w-[9rem] flex-1">
                  <label
                    htmlFor={`name-${r.key}`}
                    className="mb-1 block text-xs font-medium text-muted-foreground"
                  >
                    Field name
                  </label>
                  <input
                    id={`name-${r.key}`}
                    value={r.name}
                    onChange={(e) => update(r.key, { name: e.target.value })}
                    placeholder="title"
                    className="w-full rounded-lg border border-border bg-background px-3 py-2 font-mono text-sm outline-none focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent/30"
                  />
                </div>

                <div className="min-w-[9rem]">
                  <label
                    htmlFor={`type-${r.key}`}
                    className="mb-1 block text-xs font-medium text-muted-foreground"
                  >
                    Type
                  </label>
                  <select
                    id={`type-${r.key}`}
                    value={r.type}
                    onChange={(e) => update(r.key, { type: e.target.value, extra: '' })}
                    className="w-full rounded-lg border border-border bg-background px-3 py-2 font-mono text-sm outline-none focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent/30"
                  >
                    {GROUPS.map((g) => (
                      <optgroup key={g} label={g}>
                        {TYPES.filter((t) => t.group === g).map((t) => (
                          <option key={t.value} value={t.value}>
                            {t.label}
                          </option>
                        ))}
                      </optgroup>
                    ))}
                  </select>
                </div>

                {def && def.extra !== 'none' && (
                  <div className="min-w-[11rem] flex-1">
                    <label
                      htmlFor={`extra-${r.key}`}
                      className="mb-1 block text-xs font-medium text-muted-foreground"
                    >
                      {def.extra === 'options' && 'Choices, comma separated'}
                      {def.extra === 'model' && 'Related model'}
                      {def.extra === 'dimension' && 'Dimensions'}
                      {def.extra === 'accept' && 'Accepts'}
                      {def.extra === 'source' && 'Built from'}
                    </label>
                    {def.extra === 'accept' ? (
                      <select
                        id={`extra-${r.key}`}
                        value={r.extra || 'all'}
                        onChange={(e) => update(r.key, { extra: e.target.value })}
                        className="w-full rounded-lg border border-border bg-background px-3 py-2 font-mono text-sm outline-none focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent/30"
                      >
                        {ACCEPTS.map((a) => (
                          <option key={a} value={a}>
                            {a}
                          </option>
                        ))}
                      </select>
                    ) : (
                      <input
                        id={`extra-${r.key}`}
                        value={r.extra}
                        onChange={(e) => update(r.key, { extra: e.target.value })}
                        placeholder={
                          def.extra === 'options'
                            ? 'draft, sent, paid'
                            : def.extra === 'model'
                              ? 'User'
                              : def.extra === 'dimension'
                                ? '1536'
                                : 'title'
                        }
                        className="w-full rounded-lg border border-border bg-background px-3 py-2 font-mono text-sm outline-none focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent/30"
                      />
                    )}
                  </div>
                )}

                <fieldset className="flex items-center gap-4 pb-2">
                  <legend className="sr-only">Modifiers for {r.name || 'this field'}</legend>
                  <label
                    htmlFor={`req-${r.key}`}
                    className="flex cursor-pointer items-center gap-1.5 text-sm text-muted-foreground"
                  >
                    <input
                      id={`req-${r.key}`}
                      type="checkbox"
                      checked={r.required}
                      onChange={(e) => update(r.key, { required: e.target.checked })}
                      className="h-4 w-4 accent-[var(--accent)]"
                    />
                    required
                  </label>
                  <label
                    htmlFor={`uniq-${r.key}`}
                    className="flex cursor-pointer items-center gap-1.5 text-sm text-muted-foreground"
                  >
                    <input
                      id={`uniq-${r.key}`}
                      type="checkbox"
                      checked={r.unique}
                      onChange={(e) => update(r.key, { unique: e.target.checked })}
                      className="h-4 w-4 accent-[var(--accent)]"
                    />
                    unique
                  </label>
                </fieldset>

                <button
                  type="button"
                  onClick={() => setRows((rs) => rs.filter((x) => x.key !== r.key))}
                  disabled={rows.length === 1}
                  aria-label={`Remove ${r.name || 'this field'}`}
                  className="mb-1 rounded-lg border border-border/60 p-2 text-muted-foreground transition-colors hover:border-red-500/50 hover:text-red-400 disabled:cursor-not-allowed disabled:opacity-40"
                >
                  <Trash2 className="h-4 w-4" aria-hidden="true" />
                </button>
              </div>

              {def?.hint && !problem && (
                <p className="mt-2 text-xs text-muted-foreground/70">{def.hint}</p>
              )}
              {problem && (
                <p className="mt-2 text-xs text-amber-500" role="status">
                  {problem}
                </p>
              )}
              <Produces row={r} />
            </div>
          )
        })}
      </div>

      <button
        type="button"
        onClick={() => setRows((rs) => [...rs, row()])}
        className="mb-10 inline-flex items-center gap-2 rounded-lg border border-border/60 px-3 py-2 text-sm transition-colors hover:border-accent/60"
      >
        <Plus className="h-4 w-4" aria-hidden="true" />
        Add a field
      </button>

      </div>

      {/* The command. A column of its own rather than a bar stuck to the
          bottom of the viewport: a bottom bar covers the first rows on a
          short window, which is exactly where somebody is typing. */}
      <aside className="lg:sticky lg:top-24 lg:self-start">
      <div className="rounded-xl border border-accent/30 bg-card/60 p-5">
        <div className="mb-3 flex items-center justify-between gap-4">
          <h2 className="font-display text-sm font-semibold text-foreground">Your command</h2>
          <button
            type="button"
            onClick={copy}
            className="inline-flex items-center gap-1.5 rounded-lg border border-border/60 px-3 py-1.5 text-xs transition-colors hover:border-accent/60"
          >
            {copied ? (
              <>
                <Check className="h-3.5 w-3.5 text-emerald-400" aria-hidden="true" />
                Copied
              </>
            ) : (
              <>
                <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                Copy
              </>
            )}
          </button>
        </div>
        <pre className="overflow-x-auto whitespace-pre-wrap break-all rounded-lg bg-background/80 p-4 font-mono text-[13px] leading-relaxed text-foreground">
          {command}
        </pre>
        {problems.length > 0 && (
          <p className="mt-3 text-xs text-amber-500" role="status">
            {problems.length} field{problems.length > 1 ? 's' : ''} above {problems.length > 1 ? 'are' : 'is'}{' '}
            incomplete, so this command is not ready to run yet.
          </p>
        )}
      </div>
      </aside>
    </div>
  )
}
