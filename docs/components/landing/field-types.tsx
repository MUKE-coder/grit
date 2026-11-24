import Link from 'next/link'
import { ArrowRight } from 'lucide-react'

/*
 * The field types people actually ask about, on the home page.
 *
 * The reference has all thirty-two and that is the right place for all
 * thirty-two. It is the wrong place to find out how to add a photo, which is
 * the single most asked question about the generator: the answer was one row
 * of a table inside a 741-line page called "Field Types Reference", and
 * searching "image field" returned a battery page and a course lesson above
 * it.
 *
 * So the six that get asked about are here, on the page everybody lands on,
 * each one as the thing you type and the thing you get. Not a table of
 * thirty-two, which is a reference and nobody browses a reference.
 *
 * Every row was run through the real generator at v3.574.0. The Go types and
 * the admin controls below are its output.
 */

interface FieldRow {
  /** What you type in --fields. */
  syntax: string
  /** What the question sounds like when somebody asks it. */
  asked: string
  /** The Go the generator writes. */
  go: string
  /** The control the admin gets, with no further work. */
  control: string
}

const FIELDS: FieldRow[] = [
  {
    syntax: 'cover:file[image]',
    asked: 'An image, a photo, an avatar, a logo',
    go: '*files.FileRef',
    control: 'Dropzone with a preview, image processing and EXIF stripped',
  },
  {
    syntax: 'gallery:files[image]',
    asked: 'Several files at once',
    go: 'files.FileRefs',
    control: 'Multi-file dropzone, five by default',
  },
  {
    syntax: 'author:belongs_to:User',
    asked: 'A link to another table',
    go: 'AuthorID string + Author *User',
    control: 'Searchable relationship picker',
  },
  {
    syntax: 'status:select:draft|review|live',
    asked: 'One of a fixed set',
    go: 'string, with an enum tag',
    control: 'Searchable dropdown',
  },
  {
    syntax: 'plan:radio:free|pro',
    asked: 'One of a few, all visible',
    go: 'string, with an enum tag',
    control: 'Radio group',
  },
  {
    syntax: 'channels:check:email=Email|sms=SMS',
    asked: 'Several of a fixed set',
    go: 'datatypes.JSONSlice[string]',
    control: 'Checkbox group',
  },
]

export function FieldTypes() {
  return (
    <>
      <ul className="mt-12 grid list-none grid-cols-1 gap-4 p-0 md:grid-cols-2 lg:grid-cols-3">
        {FIELDS.map((field) => (
          <li key={field.syntax} className="g-card h-full">
            <p className="g-eyebrow mb-2.5">{field.asked}</p>
            <code className="block overflow-x-auto whitespace-nowrap rounded-lg border border-border bg-code-bg px-3 py-2 font-mono text-[12.5px] text-code-fg">
              {field.syntax}
            </code>
            <dl className="mt-3.5 space-y-1.5 text-[13px]">
              <div className="flex gap-2">
                <dt className="shrink-0 text-fg-subtle">Go</dt>
                <dd className="font-mono text-[12px] text-fg-muted">{field.go}</dd>
              </div>
              <div className="flex gap-2">
                <dt className="shrink-0 text-fg-subtle">Admin</dt>
                <dd className="text-fg-muted">{field.control}</dd>
              </div>
            </dl>
          </li>
        ))}
      </ul>

      <p className="mt-10 text-center text-[14px] leading-relaxed text-fg-muted">
        Six of thirty-two, and the ones most often asked about. Each writes the Go field,
        the TypeScript type, the Zod schema, the admin form control and the table column
        from that one declaration.
      </p>

      <div className="mt-6 flex flex-wrap items-center justify-center gap-x-6 gap-y-3 text-[14px]">
        <Link href="/docs/concepts/file-fields" className="g-link inline-flex items-center gap-1.5">
          Image and file fields in full
          <ArrowRight className="h-3.5 w-3.5" />
        </Link>
        <Link href="/docs/concepts/field-types" className="g-link">
          All thirty-two field types
        </Link>
        <Link href="/tools/field-builder" className="g-link">
          Build a --fields string
        </Link>
      </div>
    </>
  )
}
