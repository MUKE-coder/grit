import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { MaturityBadge } from '@/components/maturity-badge'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/file-fields')

/**
 * How to write an image, file or upload field.
 *
 * This is the most asked question about the generator, and until now the
 * answer was one row of a table and one line of prose inside a 741-line
 * reference page titled "Field Types Reference". Searching "image field" put
 * a battery page and a course lesson above it, which is what a title made of
 * words nobody types does.
 *
 * So: its own page, titled and keyworded with the words people actually
 * search (image, photo, avatar, logo, upload, attachment, gallery), and long
 * enough to answer the follow-up questions rather than only the first one.
 *
 * Every example on this page was run through the real generator at v3.574.0
 * and the Go types, the admin field definitions and the upload query string
 * below are copied from its output, not written from memory.
 */

const ALIASES = [
  { alias: 'image', what: 'jpg, jpeg, png, gif, webp, avif, svg' },
  { alias: 'video', what: 'the common web video container types' },
  { alias: 'audio', what: 'the common web audio types' },
  { alias: 'pdf', what: 'application/pdf' },
  { alias: 'doc', what: 'Word, Excel, PowerPoint, CSV and plain text' },
  { alias: 'all', what: 'anything the server allows. The same as leaving the brackets off.' },
]

export default function FileFieldsPage() {
  return (
    <DocsShell>
      <div className="mb-10">
        <span className="tag-mono text-accent/80 mb-3 block">Concepts</span>
        <h1 className="text-4xl font-bold tracking-tight mb-4">
          Image, file and upload fields
        </h1>
        <p className="text-lg text-muted-foreground leading-relaxed">
          How to declare a photo, an avatar, a logo, a document or a gallery on a resource,
          what each one stores, and what you get in the admin without writing any of it.
        </p>
        <MaturityBadge areas={['crud']} />
      </div>

      <div className="prose-grit">
        <h2 id="the-short-answer">The short answer</h2>
      </div>

      <CodeBlock
        language="bash"
        code={`grit generate resource Product --fields "cover:file[image],gallery:files[image]"`}
      />

      <div className="prose-grit">
        <p>
          <code>file</code> is one. <code>files</code> is many. The brackets say what the
          field accepts, and leaving them off accepts anything the server allows. That is
          the whole syntax; everything below is what it does.
        </p>

        <h2 id="what-you-get">What that writes</h2>
        <p>Copied from a real run, not paraphrased. The Go model:</p>
      </div>

      <CodeBlock
        language="go"
        filename="apps/api/internal/models/product.go"
        code={`Cover   *files.FileRef  \`gorm:"type:json" json:"cover"\`
Gallery files.FileRefs  \`gorm:"type:json" json:"gallery"\``}
      />

      <div className="prose-grit">
        <p>
          Note what it is not: a <code>string</code> holding a path. A{' '}
          <code>FileRef</code> is a JSON column carrying the key, the original filename,
          the size, the content type and the renditions, which is why the admin can show a
          thumbnail and a file size without a second request, and why{' '}
          <Link href="/docs/batteries/media">deleting a record can delete its files</Link>.
        </p>
        <p>And the admin field definition, which is where the accept list goes:</p>
      </div>

      <CodeBlock
        language="ts"
        filename="apps/admin/resources/products/products.ts"
        code={`{ key: "cover",   label: "Cover",   type: "file",  accepts: ["image"], maxSizeMB: 5 },
{ key: "gallery", label: "Gallery", type: "files", accepts: ["image"], maxSizeMB: 5, max: 5 },`}
      />

      <div className="prose-grit">
        <p>
          You get a dropzone for each, single or multiple, with drag and drop, a preview, a
          progress bar and the accept list wired in. <code>max: 5</code> is the default cap
          on how many files a <code>files</code> field takes; change it in that line.
        </p>

        <h2 id="accepts">What it accepts</h2>
        <p>
          An alias, or a bracketed list of them. Both of these are the same field with a
          different vocabulary, and both parse:
        </p>
      </div>

      <CodeBlock
        language="bash"
        code={`cover:file[image]              # the bracket form, used everywhere in these docs
cover:file:image               # the colon form, also accepted
docs:files[pdf, doc]           # a list. Spaces after the commas are fine
attachment:file                # no brackets: anything the server allows`}
      />

      <div className="my-6 overflow-hidden rounded-xl border border-border/40">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40 bg-card/50">
              <th className="px-4 py-2.5 text-left font-medium">Alias</th>
              <th className="px-4 py-2.5 text-left font-medium">Lets through</th>
            </tr>
          </thead>
          <tbody>
            {ALIASES.map((row) => (
              <tr key={row.alias} className="border-b border-border/20 last:border-0">
                <td className="px-4 py-2.5 font-mono text-[13px] text-accent align-top">
                  {row.alias}
                </td>
                <td className="px-4 py-2.5 text-[13px] text-muted-foreground leading-relaxed">
                  {row.what}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="prose-grit">
        <h2 id="enforced">The accept list is not only a browser hint</h2>
        <p>
          This is the part worth knowing, because in most stacks it is not true. The
          dropzone uses the list to filter the file picker, and it also travels to the
          server:
        </p>
      </div>

      <CodeBlock
        language="text"
        code={`POST /api/uploads?accepts=image&max_size=5242880`}
      />

      <div className="prose-grit">
        <p>
          The upload handler resolves those aliases itself through{' '}
          <code>internal/files/accepts.go</code> and refuses anything outside them, so a
          client that edits the request, or one that never used your form, gets the same
          answer. The size cap is enforced in the same place.
        </p>
        <p>
          <code>UPLOAD_ALLOWED_MIME</code> in <code>.env</code> is the outer boundary for
          the whole project. A field&apos;s accept list can narrow it and cannot widen it.
        </p>

        <h2 id="images">Images get more than storage</h2>
        <p>
          An upload that is an image goes through the processing pipeline on the way in,
          with no configuration: it is resized, re-encoded and stripped of its EXIF, which
          matters more than it sounds because a phone photo carries the GPS coordinates of
          where it was taken. The renditions land in the same <code>FileRef</code>.{' '}
          <Link href="/docs/batteries/media">Media and image processing</Link> has the
          profiles and the measured numbers.
        </p>

        <h2 id="storage">Where the bytes go</h2>
        <p>
          Wherever <code>STORAGE_DRIVER</code> says: <code>local</code> needs nothing at
          all, <code>minio</code> is what <code>docker compose up</code> starts, and{' '}
          <code>s3</code>, <code>r2</code> and <code>b2</code> need credentials. The field
          declaration is the same in every case, which is the point of the{' '}
          <code>FileRef</code> being a reference rather than a path.{' '}
          <Link href="/docs/batteries/storage">File storage</Link> covers the drivers.
        </p>

        <h2 id="gotchas">Three things that catch people</h2>
        <ul>
          <li>
            <strong>A file field is nullable and a files field is not.</strong>{' '}
            <code>*files.FileRef</code> can be <code>null</code>; <code>files.FileRefs</code>{' '}
            is an empty list when there is nothing in it. Write your frontend against that
            rather than against <code>undefined</code>.
          </li>
          <li>
            <strong>Marking it required does not mean the upload succeeded.</strong> The
            form sends a reference, and the file arrived through a separate request before
            that. If you need the two to be atomic, check the reference in your service.
          </li>
          <li>
            <strong>Renaming the field does not move the files.</strong> The key stored in
            the <code>FileRef</code> is where the object actually is. Changing the column
            name changes nothing in the bucket, which is usually what you want.
          </li>
        </ul>

        <h2 id="everything-else">Every other field type</h2>
        <p>
          This page is only about files because it is the question asked most. The{' '}
          <Link href="/docs/concepts/field-types">field types reference</Link> has all
          thirty-two, with the Go type, the TypeScript type, the Zod schema and the admin
          control each one produces.
        </p>
      </div>

      <div className="mt-12 flex flex-wrap gap-3">
        <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
          <Link href="/docs/concepts/field-types" className="gap-1.5">
            <ArrowLeft className="h-3.5 w-3.5" />
            All field types
          </Link>
        </Button>
        <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
          <Link href="/docs/batteries/media" className="gap-1.5">
            Media and image processing
            <ArrowRight className="h-3.5 w-3.5" />
          </Link>
        </Button>
      </div>
    </DocsShell>
  )
}
