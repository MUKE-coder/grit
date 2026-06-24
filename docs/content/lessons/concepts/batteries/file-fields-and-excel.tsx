import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        The previous lesson covered the <strong>S3 battery</strong> — the
        low-level upload handler, signed URLs, and the bucket abstraction.
        This lesson is about everything Grit layers on top so you almost
        never write that code by hand. You get <code>:file:</code> /{' '}
        <code>:files:</code> syntax in <code>grit generate resource</code>,
        FileRef columns with automatic lifecycle cleanup, drag-and-drop
        dropzones with five visual variants, and — new in v3.31.35 —
        client-side Excel + CSV + JSON import and export on every resource
        page.
      </p>

      <h2>Why this layer exists</h2>
      <p>
        The S3 battery gives you <code>POST /api/uploads</code> and a{' '}
        <code>storage.URL(key)</code> helper. Useful, but every resource
        that has a file column still needs:
      </p>
      <ul>
        <li>A model field that stores more than just a key (name, mime, size).</li>
        <li>A dropzone in the admin form.</li>
        <li>A preview cell in the table.</li>
        <li>Logic to delete the old file when the user picks a new one.</li>
        <li>Logic to delete orphan uploads if the form gets abandoned.</li>
        <li>A way to round-trip the column through Excel for bulk edits.</li>
      </ul>
      <p>
        Writing that six times across six resources is the work this
        layer removes.
      </p>

      <h2>:file: and :files: in the resource generator</h2>
      <p>
        The CLI parses two new field-type tokens. <code>:file:</code>{' '}
        means one uploaded file; <code>:files:</code> means a gallery.
        Both accept an optional comma-separated accept-list after a
        slash.
      </p>
      <CodeBlock
        language="bash"
        code={`# A Product with a single hero image and a gallery of PDFs.
grit generate resource Product \\
  title:string \\
  hero:file:image \\
  spec_sheets:files:pdf,doc`}
      />
      <p>
        Three things happen for each token:
      </p>
      <ul>
        <li>
          The Go model gets a <code>FileRef</code> (single) or{' '}
          <code>FileRefs</code> (multi) column stored as JSON.
        </li>
        <li>
          The Zod schema and TypeScript type emit the matching shape, so
          the admin form picks up the right field component automatically.
        </li>
        <li>
          The resource definition&apos;s <code>fields</code> array gets a{' '}
          <code>{`{ type: "file" | "files", accepts: [...] }`}</code>{' '}
          entry — that&apos;s what the FileField/FilesField components read.
        </li>
      </ul>

      <h2>FileRef — the JSON shape</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/files/types.go (generated)"
        code={`// FileRef is what gets stored on a model's JSON column. The
// 'url' is built from the storage service so it survives a CDN
// swap; mime + size let the table cell pick the right preview.
type FileRef struct {
  Key       string  ` + '`json:"key"`' + `
  URL       string  ` + '`json:"url"`' + `
  Name      string  ` + '`json:"name"`' + `
  Mime      string  ` + '`json:"mime"`' + `
  Size      int64   ` + '`json:"size"`' + `
  Thumbnail *string ` + '`json:"thumbnail_url,omitempty"`' + `
}

// FileRefs is just []FileRef — same shape, stored in a JSON
// array column.
type FileRefs []FileRef`}
      />
      <p>
        Why a struct, not just a key? Two reasons:
      </p>
      <ul>
        <li>
          <strong>Rendering without a join.</strong> The table cell needs
          mime to decide between an image thumbnail and a generic file
          icon. Storing it inline avoids hitting the <code>uploads</code>{' '}
          table on every row.
        </li>
        <li>
          <strong>CDN-resilient.</strong> The URL is rebuilt by the
          backend on read using the current <code>STORAGE_BASE_URL</code>,
          so swapping CDN domains doesn&apos;t require a data migration.
        </li>
      </ul>

      <h2>The lifecycle — what happens on save and delete</h2>
      <p>
        v3.31.30-33 closed the loop on file lifecycle. Three helpers in{' '}
        <code>internal/files</code> handle the cases that used to leak
        objects:
      </p>
      <CodeBlock
        language="go"
        filename="apps/api/internal/files/lifecycle.go (sketch)"
        code={`// Called inside the resource handler's Update path. Diffs the
// old struct against the new one, finds every *FileRef and
// FileRefs field via reflection, and S3-deletes any keys that
// were dropped. One line in the handler regardless of how many
// file columns the resource has.
files.CleanupRemoved(ctx, storage, oldProduct, newProduct)

// Called after Create + Update. Stamps claimed_at = now() on
// every Upload row referenced by the model's FileRef columns,
// so the orphan-cleanup cron knows these uploads are in use.
files.ClaimRefs(ctx, db, product)

// Daily cron job. Deletes Upload rows (and their S3 objects)
// where claimed_at IS NULL AND created_at < now() - 24h.
files.RunOrphanCleanup(ctx, db, storage, 24*time.Hour)`}
      />
      <p>
        Concretely, the three leaks this closes:
      </p>
      <ul>
        <li>
          <strong>Replace.</strong> User picks file A, saves, then
          picks file B and saves again. Without{' '}
          <code>CleanupRemoved</code>, file A stays in the bucket
          forever — the DB row points at B but A is unreachable.
        </li>
        <li>
          <strong>Gallery prune.</strong> User uploads four images, then
          removes the third in the dropzone. <code>CleanupRemoved</code>{' '}
          diffs the slice and deletes only the dropped key.
        </li>
        <li>
          <strong>Abandoned form.</strong> User uploads a file, closes
          the tab without saving. The Upload row exists; no model
          references it. The cron picks it up 24h later.
        </li>
      </ul>

      <TipBox tone="info">
        <strong>Why reflection, not codegen.</strong> Every generated
        handler is identical at the lifecycle layer —{' '}
        <code>CleanupRemoved(ctx, st, old, new)</code>. The reflection
        walk happens once per save, and the cost is well under the
        S3 round trip. Codegen would make every resource handler 30
        lines longer for no real win.
      </TipBox>

      <h2>The dropzone — five visual variants</h2>
      <p>
        The admin form picks a <code>FileField</code> or{' '}
        <code>FilesField</code> based on the resource definition. Both
        ship with five <code>dropzone</code> variants you can pick per
        field, so the same upload primitive renders as a banner uploader
        on the create form and a tiny avatar puck on the profile page.
      </p>
      <CodeBlock
        language="ts"
        filename="apps/admin/resources/users.ts"
        code={`{
  key: 'avatar',
  type: 'file',
  label: 'Profile photo',
  accepts: ['image'],
  maxSizeMB: 2,
  // Visual variant — same data shape, different chrome.
  // "default"  → boxed dashed (forms)
  // "compact"  → single inline row (modals)
  // "minimal"  → small text link (inline edit)
  // "avatar"   → circular target (profile pages)
  // "inline"   → tag-style with mini progress (lists)
  dropzone: 'avatar',
  progress: 'circular',  // 'bar' | 'circular' | 'pulse'
}`}
      />
      <p>
        Five variants because file upload is one of the UI affordances
        that has to look right in three or four different contexts.
        Forcing every form to use the same big dashed box looks
        amateurish; reskinning the same component each time is a waste.
      </p>

      <h2>v3.31.35 — Excel + CSV + JSON, in the browser</h2>
      <p>
        Every resource list page now has a split Export button (default
        format on click, chevron opens the menu) and an Import button
        next to it. All three formats — CSV, Excel (.xlsx), JSON — are
        written in the browser via{' '}
        <a
          href="https://sheetjs.com/"
          target="_blank"
          rel="noopener noreferrer"
          className="text-primary hover:underline"
        >
          SheetJS
        </a>
        . No new API routes; no <code>excelize</code> on the server; no
        async {`"`}your file is ready{`"`} email when the dataset is
        big.
      </p>

      <h3>Why client-side</h3>
      <p>
        The earlier v3.31.35 plan put export and import on the server,
        with an <code>asynq</code> + Resend cutoff for sheets over 5,000
        rows. Doing it in the browser collapses all of that:
      </p>
      <ul>
        <li>
          <strong>No new endpoints.</strong> Export reuses the existing
          paginated list route; import reuses the existing POST route.
        </li>
        <li>
          <strong>No async wiring.</strong> The browser does the work
          synchronously; large sheets show a progress bar instead of
          firing a job + email.
        </li>
        <li>
          <strong>Tenant data stays in the session.</strong> Rows
          don&apos;t round-trip through a server-side renderer just to
          build a file.
        </li>
      </ul>
      <p>
        Trade-off: very large datasets (~50k+ rows of wide records) are
        gated by browser memory, not server RAM. For the dataset sizes
        most admin panels deal with, it&apos;s fine — and the export
        loop streams every page from the API before building the file,
        so the output represents the entire filtered dataset, not just
        what&apos;s on screen.
      </p>

      <h3>Export — the menu</h3>
      <CodeBlock
        language="ts"
        filename="apps/admin/lib/excel-utils.ts (the surface)"
        code={`// Project rows onto the visible columns, write the file.
exportToFile(rows, columns, 'products', 'xlsx')

// Loop the resource endpoint until every row is in hand. The
// search params come from the active filter / sort / date range,
// so the export honours the user's current view.
const rows = await fetchAllPages(
  '/api/products',
  apiSearchParams,
  (loaded, total) => setProgress({ loaded, total }),
)`}
      />
      <p>
        The export menu lives on the toolbar of every resource page.
        Default click triggers the resource&apos;s default format (Excel
        when enabled, else CSV); the chevron opens the menu for the
        other formats. Hidden columns are excluded from the file, so
        what you see is what you get.
      </p>

      <h3>Import — three stages</h3>
      <p>
        The Import button opens a modal with three stages.
      </p>
      <ol>
        <li>
          <strong>Pick.</strong> Drag-and-drop or click to browse. A
          {` "`}Download template{`"`} link generates a blank workbook
          keyed by the resource&apos;s field names plus an example row.
        </li>
        <li>
          <strong>Preview.</strong> SheetJS parses the file in-browser,
          maps headers to fields (case + space + underscore + hyphen
          insensitive), coerces each cell to the right JS type via the
          field definition, and surfaces per-row errors. Unknown header
          columns are flagged so users notice a typo before submitting.
        </li>
        <li>
          <strong>Submit.</strong> Each valid row is POSTed to the
          resource endpoint at concurrency 4 with a live progress bar.
          React Query invalidates on success so the list refreshes.
          Failed rows get a per-row reason in the summary screen.
        </li>
      </ol>

      <h3>Header matching is loose</h3>
      <p>
        Real spreadsheets come from finance and marketing, not
        engineers. Headers like <code>First Name</code>,{' '}
        <code>first_name</code>, <code>firstname</code>, and{' '}
        <code>First-Name</code> all map to the same field — Grit
        normalises whitespace, underscores, hyphens, and case before
        lookup. The template you download is keyed by the wire field
        name, but the import accepts anything that round-trips through
        that normaliser.
      </p>

      <h3>Per-resource opt-outs</h3>
      <CodeBlock
        language="ts"
        filename="apps/admin/resources/products.ts"
        code={`export const productsResource = defineResource({
  // ...
  table: {
    columns: [/* ... */],
    // Hide a format from the export menu (default: all on).
    export: { csv: true, excel: true, json: false },
    // Or disable export entirely:
    // export: false,

    // Restrict importable fields to a subset (defaults to every
    // form field). Useful when you want to exclude generated
    // columns or fields that need server-side calculation.
    import: { fields: ['title', 'price', 'stock'] },
    // Or disable import entirely:
    // import: false,
  },
})`}
      />
      <p>
        File fields (<code>:file:</code> / <code>:files:</code>) are
        automatically excluded from imports — a spreadsheet can&apos;t
        carry a binary blob, so it makes no sense to accept them
        through this path. Users upload files the normal way and pair
        them with metadata via the import.
      </p>

      <h2>Putting it together — a product catalog flow</h2>
      <p>
        The typical end-to-end shape this layer enables:
      </p>
      <ol>
        <li>
          <strong>Scaffold the resource:</strong>{' '}
          <code>grit generate resource Product title:string price:number hero:file:image</code>.
        </li>
        <li>
          <strong>Bulk-create from a spreadsheet:</strong> click
          Import, drop a 500-row .xlsx with title and price columns,
          confirm validation, submit. 500 products land in seconds.
        </li>
        <li>
          <strong>Add hero images one by one:</strong> open each
          product, drop the hero image into the avatar-variant
          dropzone, save. The S3 battery handles the upload; the
          lifecycle helper claims the upload row.
        </li>
        <li>
          <strong>Edit a batch in Excel:</strong> filter the list to{' '}
          {`"`}last 30 days{`"`}, click Export → Excel. Edit prices in
          the spreadsheet. Re-import via a separate sheet, or update
          individual rows. (Update via import is on the roadmap; for
          now imports create new rows.)
        </li>
        <li>
          <strong>Delete a product:</strong> the lifecycle helper
          deletes the hero image from S3 in the same request.
        </li>
      </ol>

      <TipBox tone="warning">
        <strong>Excel currently creates, not updates.</strong>{' '}
        v3.31.35&apos;s import POSTs each row, which means the API
        treats it as a new record. Update-by-ID via import is on the
        roadmap; until then, hand-edit the few rows you need to change
        through the admin form, or write a one-off script for batch
        updates.
      </TipBox>

      <h2>Files in this layer — for when you need to dig</h2>
      <CodeBlock
        language="text"
        code={`apps/api/internal/files/
├── types.go           ← FileRef, FileRefs
├── lifecycle.go       ← CleanupRemoved, ClaimRefs, RunOrphanCleanup
└── orphan_cleanup.go  ← The daily cron job

apps/admin/lib/
├── excel-utils.ts     ← SheetJS wrappers: export, import, template

apps/admin/components/tables/
├── export-menu.tsx    ← Split button + format dropdown
├── import-modal.tsx   ← Drop → preview → submit
├── table-toolbar.tsx  ← Mounts both above

apps/admin/components/forms/fields/
├── file-field.tsx     ← Single FileRef field
└── files-field.tsx    ← FileRefs gallery field

apps/admin/components/ui/
└── dropzone.tsx       ← The five-variant primitive both fields wrap`}
      />

      <KnowledgeCheck
        question="A user uploads a hero image to a Product, saves, then changes the hero to a different image and saves again. What happens to the first image's S3 object?"
        choices={[
          {
            label: 'Both objects stay in S3 — there is no automatic cleanup',
            feedback:
              "That was true before v3.31.33. Now CleanupRemoved diffs the old and new model values and S3-deletes any FileRef key that's no longer referenced.",
          },
          {
            label: 'CleanupRemoved diffs old vs new on the Update handler and deletes the dropped key from S3',
            correct: true,
            feedback:
              "Right — one line in the handler regardless of how many file columns the resource has. Reflection finds every *FileRef and FileRefs field, computes the diff, deletes what's gone.",
          },
          {
            label: 'The daily orphan-cleanup cron picks it up 24 hours later',
            feedback:
              "The orphan cron handles abandoned uploads (claimed_at IS NULL), not replaced ones. Replacements are caught immediately by CleanupRemoved.",
          },
          {
            label: 'The frontend sends a DELETE request for the old key before saving',
            feedback:
              "That would leak: if the save fails after the delete, the user is left without their original image. Lifecycle has to be transactional with the save, which is why it lives on the backend.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Scaffold a small catalog and exercise the whole flow:</p>
            <ol>
              <li>
                <code>grit generate resource Product title:string price:number hero:file:image gallery:files:image</code>.
              </li>
              <li>
                In the admin, click Import, download the template,
                fill in 10 rows of titles and prices, drop it back,
                submit. Confirm 10 products land in the list.
              </li>
              <li>
                Pick one product. Upload a hero image. Save. Replace
                it with a different image. Save. Open MinIO console
                at <code>localhost:9001</code> and confirm the bucket
                has only one hero, not two.
              </li>
              <li>
                Filter the list to {`"`}Last 7 days{`"`}, then Export
                → Excel. Confirm only the 10 you created come down,
                not whatever was there before.
              </li>
              <li>
                On the Product resource, set{' '}
                <code>{`table: { export: { json: false } }`}</code>{' '}
                and confirm the JSON entry disappears from the menu.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            For step 3 the bucket name is <code>uploads</code> by
            default; sort by <code>Last Modified</code> in the console
            to confirm the timestamp matches your latest save and the
            older object is gone.
          </>
        }
        solution={
          <>
            <p>
              You should have a working catalog where bulk-create
              comes from a spreadsheet, single-file edits happen in
              the form, and the bucket stays clean. That&apos;s the
              loop most products run on.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — <strong>Mail (Resend)</strong>. Transactional
        email with editable templates — the most-needed external
        service after the DB. Coming after Mail: PDF export via{' '}
        <code>@react-pdf/renderer</code> (v3.31.36), which lets the
        same Export menu offer styled PDF receipts and reports
        alongside Excel.
      </p>
    </>
  )
}
