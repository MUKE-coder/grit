import Link from 'next/link'
import { ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/changelog')

export default function ChangelogPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Release History</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Changelog
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                All notable changes to Grit are documented here. Each release includes new features,
                bug fixes, and any breaking changes you need to be aware of.
              </p>
            </div>

            {/* v3.32.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.32.0
                </span>
                <span className="text-sm text-muted-foreground">July 9, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Offline-first desktop, and a deep security &amp; correctness pass.</strong>{' '}
                  This release makes the monorepo desktop client a true
                  online/offline hybrid, and fixes a batch of issues found in a
                  full audit of the generated code.
                </p>
                <p>
                  <strong>New — offline-hybrid desktop.</strong> The{' '}
                  <code>apps/desktop</code> client (from <code>--full</code> or{' '}
                  <code>--desktop</code>) now works online by default,{' '}
                  <em>continuously mirroring</em> server data into a local SQLite
                  copy in the background. A <strong>Work offline</strong> toggle in
                  the dashboard&apos;s Settings lets you keep working against that
                  local copy with no connection; every edit queues, and the moment
                  you switch back online it <strong>auto-reconciles</strong> — pushes
                  your changes (with the existing per-field conflict merge) and pulls
                  anything new. Deletes now propagate to offline clients via
                  tombstones. <code>grit generate resource</code> registers each new
                  model for offline sync automatically.
                </p>
                <p>
                  <strong>Security.</strong> Closed a SQL-injection vector in the
                  shared paginator&apos;s <code>date_field</code> parameter
                  (reachable on every generated list endpoint) and whitelisted the
                  generated service&apos;s <code>ORDER BY</code>. Uploads now sniff
                  real content type instead of trusting the client header, cap the
                  request body, and reject HTML/SVG payloads. The seeder refuses the
                  default <code>admin123</code> password in production, and token
                  refresh re-checks that the account still exists and is active.
                </p>
                <p>
                  <strong>Correctness.</strong> Fixed generated <strong>desktop CRUD</strong>{' '}
                  (models now assign their UUID and use string IDs end-to-end — the
                  old code could store only one record and silently no-op updates and
                  deletes); the desktop embedded API moved off port <code>34115</code>{' '}
                  so it no longer collides with <code>wails dev</code>; date fields,
                  a <code>belongs_to</code> CSV-import build breaker, and a mobile{' '}
                  <code>Bearer undefined</code> token-refresh bug are all fixed.
                  Backups now stream (no more loading the whole database into memory),
                  include many-to-many join tables, and no longer corrupt values
                  containing <code>--</code>. CSV import streams and batches instead
                  of buffering the whole file, and stalled import jobs are reaped.
                </p>
              </div>
            </div>

            {/* v3.31.83 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.83
                </span>
                <span className="text-sm text-muted-foreground">July 8, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>NSIS made discoverable for desktop installers.</strong>{' '}
                  Building a Windows installer with <code>grit package</code> needs
                  NSIS, and it was easy to miss. The desktop app&apos;s README now
                  lists NSIS as a prerequisite (with{' '}
                  <code>winget install NSIS.NSIS</code> and friends), and when
                  <code>makensis</code> is missing <code>grit package</code> now
                  prints the exact install commands and a PATH hint instead of a
                  bare link — or points you at <code>--no-installer</code>. No
                  behaviour change, just fewer dead ends.
                </p>
              </div>
            </div>

            {/* v3.31.82 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.82
                </span>
                <span className="text-sm text-muted-foreground">July 8, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>New: <code>grit package</code> — build a distributable
                  desktop installer.</strong> Run it inside a{' '}
                  <code>grit new-desktop</code> app and it produces the artifact you
                  hand to a user: on Windows an NSIS installer (the single{' '}
                  <code>*-installer.exe</code> in <code>build/bin/</code>), on
                  macOS/Linux the platform binary/app bundle. It wraps{' '}
                  <code>wails build</code>, checks the toolchain (<code>wails</code>,
                  plus <code>makensis</code> for the installer) up front with a
                  clear error, and prints where the artifact landed.{' '}
                  <code>--no-installer</code> builds the raw binary only;{' '}
                  <code>--platform</code> cross-compiles. For a full versioned
                  release, <code>scripts/release-desktop.sh &lt;version&gt;</code>{' '}
                  still ships.
                </p>
              </div>
            </div>

            {/* v3.31.81 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.81
                </span>
                <span className="text-sm text-muted-foreground">July 8, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Desktop relationship fields are now a real dropdown.</strong>{' '}
                  A <code>belongs_to</code> field in a generated desktop form used
                  to render as a plain text box where you had to paste the related
                  row&apos;s id. It now loads the related records via their list
                  binding and renders a proper <code>&lt;select&gt;</code> of
                  names — pick a Category from the list instead of typing a UUID.
                </p>
              </div>
            </div>

            {/* v3.31.80 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.80
                </span>
                <span className="text-sm text-muted-foreground">July 8, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Fix: desktop list crashed on file fields.</strong> A
                  generated desktop list rendered a <code>file</code> field&apos;s
                  FileRef object directly into a table cell, which React refuses
                  (&quot;Objects are not valid as a React child&quot;). File
                  columns now render a thumbnail (and <code>files</code> columns a
                  small stack), so an inventory list with a product photo displays
                  instead of white-screening. Completes the desktop upload support
                  from v3.31.79.
                </p>
              </div>
            </div>

            {/* v3.31.79 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.79
                </span>
                <span className="text-sm text-muted-foreground">July 8, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Desktop apps are now hybrid — and support file uploads.</strong>{' '}
                  A <code>grit new-desktop</code> app still runs on Wails + SQLite/Postgres,
                  but it now also embeds a real Gin REST API in the same binary.
                  The router is mounted twice: as the Wails asset-server handler
                  (so the webview calls <code>/api/…</code> and loads{' '}
                  <code>&lt;img src="/uploads/x.jpg"&gt;</code> same-origin, no port,
                  no CORS) and on <code>127.0.0.1:34115</code> for curl / other
                  clients.
                </p>
                <p>
                  <strong>File uploads work end-to-end.</strong>{' '}
                  <code>grit generate resource ... photo:file:image</code> now
                  produces a working image field: a native file picker that uploads
                  to <code>POST /api/uploads</code>, files stored under the OS
                  app-data dir (writable even when the app is installed in Program
                  Files), a preview in the form, and <code>files:image</code> for
                  multi-image galleries. New <code>internal/files</code>,{' '}
                  <code>internal/storage</code> and <code>internal/api</code>{' '}
                  packages back it.
                </p>
                <p>
                  <strong>Two codegen bugs fixed along the way:</strong> a{' '}
                  <code>file:</code> field used to emit <code>*files.FileRef</code>{' '}
                  with no import (and no <code>files</code> package at all), breaking
                  the build; and a <code>slug</code> field called{' '}
                  <code>slugify()</code> from a package that didn&apos;t define it.
                  Both now compile.
                </p>
              </div>
            </div>

            {/* v3.31.78 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.78
                </span>
                <span className="text-sm text-muted-foreground">July 8, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Grit UI is no longer baked into generated apps.</strong> It
                  lives on as a standalone library, so a new project starts lean.
                  Scaffolded apps no longer include the <code>UIComponent</code>{' '}
                  model, the registry handler (<code>/r.json</code>,{' '}
                  <code>/r/:name</code>, <code>/ui-components</code>, admin CRUD), the
                  91-component seeder, <code>packages/grit-ui/</code>, or the web{' '}
                  <code>/components</code> browser.
                </p>
                <p>
                  A fresh <code>--triple</code> project now registers{' '}
                  <strong>19 models instead of 20</strong>, and seeding no longer
                  plants 100 component rows — your first backup drops from 109 rows to
                  9. Existing projects are untouched; delete those files yourself if
                  you want the same trim.
                </p>
              </div>
            </div>

            {/* v3.31.77 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.77
                </span>
                <span className="text-sm text-muted-foreground">July 8, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Automatic weekly database backups.</strong> Every Grit API
                  now takes a full-database backup <strong>every Sunday at 02:00
                  UTC</strong> and uploads it to your object storage (R2 / S3 /
                  MinIO). The four most recent are kept; older ones are purged from
                  storage but their rows survive as an audit trail.
                </p>
                <p>
                  Each archive is a ZIP: one <strong>CSV per table</strong> (opens in
                  any spreadsheet), a <code>dump.sql</code> of INSERTs in
                  parent&rarr;child order wrapped in BEGIN/COMMIT, and a{' '}
                  <code>metadata.json</code> manifest of row counts. It&apos;s pure Go
                  &mdash; no <code>pg_dump</code> binary &mdash; so it works on
                  Postgres and SQLite alike. The table list is derived from{' '}
                  <code>models.Models()</code>, so every{' '}
                  <code>grit generate resource</code> is included automatically and a
                  table name can never be injected.
                </p>
                <p>
                  Four surfaces: a <strong>Backups</strong> page in the admin panel
                  (list, back up now, download), REST endpoints (
                  <code>GET /backups</code>, <code>POST /backups/generate</code>,{' '}
                  <code>GET /backups/:id/download</code> &mdash; which mints a
                  15-minute pre-signed URL so the browser pulls straight from
                  storage), a mobile <strong>Backups</strong> screen, and the CLI:
                </p>
                <pre><code>{`grit backup                 # dump + upload to object storage
grit backup -o backup.zip   # write a local archive (no storage needed)
grit restore backup.zip     # migrate, then replay in ONE transaction`}</code></pre>
                <p>
                  Restore is a first-class command, not a doc page &mdash; a backup you
                  have never restored is a rumour. Manual backups are rate-limited to
                  one per 24h; the weekly cron bypasses it and uses{' '}
                  <code>asynq.Unique</code> so a rolling deploy can&apos;t enqueue it
                  twice. Without object storage configured (typical in dev) the weekly
                  job skips silently.
                </p>
              </div>
            </div>

            {/* v3.31.76 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.76
                </span>
                <span className="text-sm text-muted-foreground">July 6, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Numeric inputs format as you type — and never rescale.</strong>{' '}
                  Number fields in generated mobile forms now show thousands
                  separators while typing (<code>1000</code> renders{' '}
                  <code>1,000</code>) and submit the plain number. What you type is
                  what&apos;s stored: enter <code>100</code> and the record holds{' '}
                  <code>100</code> — no cents conversion, no divide-by-100.
                  New <code>lib/format.ts</code> exposes{' '}
                  <code>formatNumberInput()</code> / <code>parseNumberInput()</code>;
                  <code>float</code> fields keep up to two decimals.
                </p>
              </div>
            </div>

            {/* v3.31.75 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.75
                </span>
                <span className="text-sm text-muted-foreground">July 6, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Multi-image fields on mobile.</strong> A{' '}
                  <code>files</code> field (<code>name:files:image</code>) now
                  renders a proper multi-picker: select several photos from the
                  gallery at once, see them as a grid of removable thumbnails,
                  each uploaded in the background, and the payload carries an
                  array of file references. Single <code>file</code> fields stay
                  single-select. The picker sheet already supported multi-select;
                  generated forms now use it for array fields.
                </p>
              </div>
            </div>

            {/* v3.31.74 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.74
                </span>
                <span className="text-sm text-muted-foreground">July 6, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Searchable select for relationships on mobile.</strong>{' '}
                  A <code>belongs_to</code> field in a generated form used to
                  render every related record as a horizontal row of pills, which
                  falls apart once there are more than a handful. It now uses a new{' '}
                  <code>RelationSelect</code> component: a tidy select that opens a
                  bottom sheet with a pinned search box and a scrollable, filtered
                  list — pick one and it fills in. Wired into every generated
                  resource form; regenerating a resource with a relationship picks
                  it up.
                </p>
              </div>
            </div>

            {/* v3.31.73 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.73
                </span>
                <span className="text-sm text-muted-foreground">July 6, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>MinIO is now reachable from mobile devices.</strong> The
                  dev <code>docker-compose.yml</code> published MinIO on{' '}
                  <code>127.0.0.1:9002</code> (localhost only), so a phone or
                  emulator couldn&apos;t load uploaded images even though the API
                  (bound to all interfaces) worked fine — list and detail
                  thumbnails stayed blank. It now binds <code>9002:9000</code> on
                  all interfaces, so devices on your LAN can fetch stored images.
                  Pairs with <code>resolveImageUrl()</code> (v3.31.72), which
                  rewrites the <code>localhost</code> host to your dev IP. Existing
                  projects: change the minio <code>ports</code> to{' '}
                  <code>&quot;9002:9000&quot;</code> /{' '}
                  <code>&quot;9003:9001&quot;</code> and{' '}
                  <code>docker compose up -d minio</code>.
                </p>
              </div>
            </div>

            {/* v3.31.72 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.72
                </span>
                <span className="text-sm text-muted-foreground">July 6, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Image previews everywhere on mobile.</strong> Generated
                  resources now show pictures throughout: an instant local preview
                  in the create/edit form the moment you pick a photo (with an
                  upload spinner overlay), a thumbnail column in the list table,
                  and the hero image on the detail screen.
                </p>
                <p>
                  New <code>lib/images.ts</code> → <code>resolveImageUrl()</code>{' '}
                  fixes the classic dev gotcha: MinIO hands back{' '}
                  <code>http://localhost:9002/...</code> URLs that a device or
                  emulator can&apos;t reach (localhost = the device itself). It
                  rewrites the host to the same dev host the app already uses for
                  the API, so stored images actually load — while real S3/R2
                  public URLs pass through untouched. Every generated list, detail
                  and form image runs through it; regenerating a resource adds it.
                </p>
              </div>
            </div>

            {/* v3.31.71 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-danger/15 px-3 py-1 text-sm font-semibold text-danger">
                  v3.31.71 · critical fix
                </span>
                <span className="text-sm text-muted-foreground">July 6, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Fix: request bodies larger than 4 KB were silently
                  truncated.</strong> Pulse&apos;s error-tracking middleware
                  captures a request-body snippet for error context, but it
                  restored <em>only</em> that snippet to the request — discarding
                  everything past <code>MaxBodySize</code> (default 4096 bytes).
                  Every request that sends a <code>Content-Length</code> (mobile
                  apps, native/CLI clients, <code>curl</code>) reached handlers
                  with a body cut to 4 KB, so <strong>file uploads and any large
                  JSON POST failed</strong> with confusing &quot;no file&quot; /
                  parse errors. Browsers were unaffected because <code>fetch</code>
                  sends multipart <em>chunked</em> (no Content-Length), which
                  skipped the capture — which is why the web dropzone always
                  worked while mobile never did.
                </p>
                <p>
                  The scaffold now mounts Pulse with{' '}
                  <code>WithRequestBodyCaptureDisabled()</code>, so the full body
                  always reaches your handlers. This affects every generated API;
                  regenerate or add that option to your Pulse mount. Mobile image
                  uploads (avatar, resource forms, imports) now work end-to-end.
                </p>
              </div>
            </div>

            {/* v3.31.70 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.70
                </span>
                <span className="text-sm text-muted-foreground">July 6, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>A proper image picker for mobile forms.</strong> Tapping
                  an image field used to jump straight into Android&apos;s system
                  crop screen — whose only button was CROP, with no clear
                  &quot;use this photo&quot; and no permission prompt of our own.
                </p>
                <p>
                  Now it opens a clean, themed <strong>picker sheet</strong>:
                  choose <em>Library</em> or <em>Camera</em> (with a friendly
                  permission prompt and an &quot;Open Settings&quot; fallback if
                  access is off), then <strong>preview</strong> the selection and
                  decide — <strong>Use photo</strong>, <strong>Crop</strong>{' '}
                  (the native editor, only when you ask for it), or choose a
                  different one. The dropzone shows a spinner while the upload
                  runs. The generated resource form uses it for every image
                  field; regenerating any resource upgrades an existing app.
                </p>
              </div>
            </div>

            {/* v3.31.69 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.69
                </span>
                <span className="text-sm text-muted-foreground">July 6, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Fix: image uploads from the Expo app.</strong> Uploads
                  were failing with <code>400 &quot;No file provided&quot;</code>{' '}
                  because <code>expo-file-system</code>&apos;s{' '}
                  <code>uploadAsync</code> sends an empty body under the New
                  Architecture on SDK 54 (the request landed in a few ms with no
                  file). The mobile upload helper now uses{' '}
                  <code>fetch</code> + <code>FormData</code> with a React Native
                  file descriptor and — crucially — never sets{' '}
                  <code>Content-Type</code> by hand, so fetch keeps the multipart
                  boundary intact. Fixes avatar, blog and every generated
                  resource-form image field.
                </p>
                <p>
                  The <code>/uploads</code> handler is also more robust: it falls
                  back to the first file part under any field name and logs the
                  request&apos;s content-type + fields when a file is genuinely
                  missing, so client-side multipart problems are diagnosable from
                  the server terminal. Re-run <code>pnpm i</code> in{' '}
                  <code>apps/expo</code> after updating (the helper no longer
                  needs <code>expo-file-system</code> for uploads).
                </p>
              </div>
            </div>

            {/* v3.31.68 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.68
                </span>
                <span className="text-sm text-muted-foreground">July 6, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Background CSV import — imports now run server-side and
                  survive leaving the screen.</strong> The import endpoint no
                  longer blocks: it reads the upload, creates a job, processes
                  rows in a goroutine and returns <code>202</code> with a job id.
                </p>

                <h3>Backend (every architecture)</h3>
                <p>
                  A shared <code>ImportJob</code> table tracks every resource&apos;s
                  imports. <code>POST /&lt;plural&gt;/import</code> kicks the work
                  off and returns immediately; a new shared{' '}
                  <code>GET /imports/:id</code> reports live{' '}
                  <code>processed / total</code> plus the final
                  created / skipped / failed counts and per-row errors, so a
                  large file never times the request out.
                </p>

                <h3>Mobile</h3>
                <p>
                  Imports run in a module-level store, so they keep uploading and
                  polling even after you close the import sheet. A persistent
                  progress <strong>banner</strong> shows every in-flight import
                  across navigation, then the result — tap &quot;Continue in
                  background&quot; and carry on using the app. Regenerating any
                  resource upgrades an existing app to the new flow.
                </p>
              </div>
            </div>

            {/* v3.31.67 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.67
                </span>
                <span className="text-sm text-muted-foreground">July 3, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>CSV import — bulk-create records from a spreadsheet.</strong>{' '}
                  <code>grit generate resource</code> now generates a bulk import endpoint (all
                  architectures) plus a full mobile import flow.
                </p>

                <h3>Backend (every architecture)</h3>
                <p>
                  Each resource gets <code>POST /&lt;plural&gt;/import</code> (upload a CSV → typed
                  bulk-create) and <code>GET /&lt;plural&gt;/import/template</code> (a ready-to-fill
                  header CSV). Everything is optional except model-required fields:{' '}
                  <code>file</code> columns are skipped; a <code>belongs_to</code> is given by{' '}
                  <strong>name</strong> (a <code>category</code> column, not <code>category_id</code>)
                  and the related record is looked up and <strong>created if missing</strong>; rows
                  that hit a unique constraint are <strong>skipped</strong> (safe to re-import) and
                  other failures are reported per-row.
                </p>

                <h3>Mobile</h3>
                <p>
                  The resource list gains an import action that opens a sheet: <strong>download the
                  template</strong>, pick a CSV, <strong>preview</strong> the parsed rows, import
                  with a <strong>progress bar</strong>, then a summary of created / skipped /
                  failed with per-row errors. Adds <code>expo-document-picker</code> — run{' '}
                  <code>pnpm i</code>. Background (async) import is the final step.
                </p>
              </div>
            </div>

            {/* v3.31.66 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.66
                </span>
                <span className="text-sm text-muted-foreground">July 3, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Mobile: relationship filters.</strong> Resource lists with a{' '}
                  <code>belongs_to</code> field gain a funnel action that opens a filter sheet.
                </p>
                <p>
                  The sheet shows a picker per relationship (loaded from the related resource);
                  pick a value to scope the table (<code>?&lt;fk&gt;=&lt;id&gt;</code>, which the
                  API already supports), with an <strong>All</strong> chip and a{' '}
                  <strong>Clear all</strong>. The funnel shows a dot while filters are active, and
                  export respects them. Resources without a relationship simply don&apos;t show the
                  funnel. Re-run <code>grit generate resource</code> to pick it up. CSV import is
                  the last piece.
                </p>
              </div>
            </div>

            {/* v3.31.65 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.65
                </span>
                <span className="text-sm text-muted-foreground">July 3, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Mobile: CSV export.</strong> Every generated resource list gains a
                  download action that exports the data as CSV and opens the native share sheet.
                </p>
                <p>
                  Tapping export downloads <code>/&lt;plural&gt;/export</code> (honouring the
                  current search) to a file and hands it to the OS share sheet — mail it, save it,
                  open it in Sheets. New <code>lib/export.ts</code> helper built on{' '}
                  <code>expo-file-system</code> + <code>expo-sharing</code> (run <code>pnpm i</code>).
                  Filters and CSV import land next.
                </p>
              </div>
            </div>

            {/* v3.31.64 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.64
                </span>
                <span className="text-sm text-muted-foreground">July 3, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Mobile: scrollable data table.</strong> Generated resource lists now
                  render as a <strong>horizontally-scrollable table</strong> — a column per field
                  with tap-to-sort headers — instead of cards.
                </p>
                <p>
                  The title field leads (bold), followed by a column for every scalar and{' '}
                  <code>belongs_to</code> field; dates, numbers and booleans format per cell. Tap a
                  sortable header to sort (wired to the API&apos;s <code>sort_by</code> /{' '}
                  <code>sort_order</code>, which the list hook now accepts), tap a row to open the
                  detail. Search, infinite-scroll pagination and the quick-create sheet are
                  unchanged. Re-run <code>grit generate resource</code> to switch a list to the
                  table.
                </p>
              </div>
            </div>

            {/* v3.31.63 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.63
                </span>
                <span className="text-sm text-muted-foreground">July 3, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Mobile: quick-create bottom sheet.</strong> Adding a record no longer
                  always means a full-screen navigation — the resource list&apos;s <code>+</code>{' '}
                  now opens a slide-up <strong>sheet</strong> with the form, while detailed edits
                  stay a full <strong>page</strong>.
                </p>
                <p>
                  New shared <code>FormSheet</code> component — a themed, keyboard-aware bottom sheet
                  built on React Native&apos;s <code>Modal</code> (no extra dependencies). The
                  generated list renders the same <code>&lt;Name&gt;Form</code> inside it for a
                  fast add; the detail screen&apos;s Edit still opens the full page for longer
                  records. One form, two containers. Re-run{' '}
                  <code>grit generate resource</code> to pick it up.
                </p>
              </div>
            </div>

            {/* v3.31.62 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.62
                </span>
                <span className="text-sm text-muted-foreground">July 3, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Mobile: full CRUD on generated resources.</strong> Generated resources
                  now support <strong>edit, update and delete</strong>, not just create + read.
                </p>
                <p>
                  <code>grit generate resource</code> now emits a shared{' '}
                  <code>&lt;Name&gt;Form</code> component (in{' '}
                  <code>components/resource-forms/</code>) that both the create and edit screens
                  render — the create page and a new{' '}
                  <code>app/&lt;plural&gt;/edit/[id].tsx</code> screen pre-fill from the record and
                  drive <code>useCreate</code> / <code>useUpdate</code>. The detail screen gains{' '}
                  <strong>Edit</strong> and <strong>Delete</strong> (with a confirm) actions. The
                  shared form is container-agnostic, ready to drop into a bottom sheet next. Because
                  generation is now idempotent, re-run{' '}
                  <code>grit generate resource &lt;Name&gt; --fields …</code> to add CRUD to an
                  existing resource.
                </p>
              </div>
            </div>

            {/* v3.31.61 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.61
                </span>
                <span className="text-sm text-muted-foreground">July 3, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Mobile: built-in Blog &amp; User resources.</strong> The scaffolded Expo
                  app now surfaces the framework&apos;s built-in Blog feature and lets you add
                  users — no admin panel required.
                </p>
                <p>
                  The More tab&apos;s <strong>Resources</strong> section now leads with{' '}
                  <strong>Users</strong> and <strong>Blogs</strong> (alongside your generated
                  resources). Blogs get a paginated, searchable list plus a create screen (title,
                  excerpt, content, cover image upload, publish toggle) backed by{' '}
                  <code>/admin/blogs</code> — the posts <code>grit seed</code> already creates now
                  have a home. The Users list gained a <code>+</code> to create a user (name,
                  email, password, role, active) via <code>/admin/users</code>.
                </p>
              </div>
            </div>

            {/* v3.31.60 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.60
                </span>
                <span className="text-sm text-muted-foreground">July 3, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>generate resource is now idempotent.</strong> Re-running{' '}
                  <code>grit generate resource &lt;Name&gt;</code> for a resource that already
                  exists used to append duplicate injections — duplicate switch cases, routes and
                  exports — which broke the API build. Now every injection is skipped when it&apos;s
                  already present.
                </p>
                <p>
                  The two low-level inject helpers gained a whitespace-insensitive
                  &quot;already there?&quot; guard, so a second run finds each of its injections in
                  place and does nothing (files are also just overwritten). Safe to re-run to pick
                  up regenerated files, or after editing a resource&apos;s fields.
                </p>
              </div>
            </div>

            {/* v3.31.59 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.59
                </span>
                <span className="text-sm text-muted-foreground">July 3, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Mobile: create forms + a More tab hub.</strong> Generated mobile
                  resources can now <em>add</em> data, not just browse it — so a{' '}
                  <code>--mobile</code> project no longer needs an admin panel to get started.
                </p>

                <h3>Create screen per resource</h3>
                <p>
                  <code>grit generate resource</code> now also scaffolds{' '}
                  <code>app/&lt;plural&gt;/new.tsx</code> — a form wired to the generated{' '}
                  <code>useCreate&lt;X&gt;</code> mutation, with an input per field: text / number /
                  textarea / toggle for scalars, an image picker (upload → FileRef) for{' '}
                  <code>file</code> fields, and a chip picker for <code>belongs_to</code>{' '}
                  relationships. The list screen gains a <code>+</code> button to reach it.
                </p>

                <h3>&quot;More&quot; tab</h3>
                <p>
                  The mobile Explore tab is now <strong>More</strong> (with an ellipsis icon) and
                  acts as the app hub: a <strong>Resources</strong> section that{' '}
                  <code>grit generate resource</code> injects each new resource into, plus the
                  Users / Storage / Analytics / Notifications tools. Restart nothing — reload the
                  Expo app.
                </p>
              </div>
            </div>

            {/* v3.31.58 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.58
                </span>
                <span className="text-sm text-muted-foreground">July 3, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Fix: mobile file uploads (&quot;No file provided&quot;).</strong> Avatar /
                  image uploads from the Expo app failed with a <code>400</code> even though a file
                  was selected.
                </p>
                <p>
                  Two causes, both fixed. On the client, React Native&apos;s{' '}
                  <code>fetch</code> + <code>FormData</code> (RN&nbsp;0.81 / Expo SDK&nbsp;54) can
                  drop the file part entirely; the upload helper now uses{' '}
                  <code>expo-file-system</code>&apos;s native <code>uploadAsync</code>, which streams
                  the file reliably. On the server, the audit-log middleware read the entire request
                  body to digest it — wasteful for a binary upload and enough to leave{' '}
                  <code>ParseMultipartForm</code> with nothing to parse; it now skips{' '}
                  <code>multipart/form-data</code> bodies. Upload errors also surface the server&apos;s
                  actual reason now instead of a generic failure. Adds{' '}
                  <code>expo-file-system</code> — run <code>pnpm i</code>, restart the API, and{' '}
                  <code>npx expo start -c</code>.
                </p>
              </div>
            </div>

            {/* v3.31.57 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.57
                </span>
                <span className="text-sm text-muted-foreground">July 3, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Fix: mobile mutations blocked by CSRF (403).</strong> Native clients
                  could read data but every write — file uploads, generated-resource
                  create/update/delete, profile updates — failed with a{' '}
                  <code>403 CSRF_INVALID</code>.
                </p>
                <p>
                  React Native&apos;s <code>fetch</code> (and Android&apos;s OkHttp) transparently
                  store and resend the <code>grit_access</code> cookie the API sets at login. The
                  <code>AutoCSRF</code> guard saw that stray cookie and treated a{' '}
                  <strong>bearer</strong>-authenticated request as cookie-authenticated, demanding a
                  CSRF token the app never sends. The guard now skips CSRF whenever an{' '}
                  <code>Authorization: Bearer</code> header is present — an explicitly-authenticated
                  request can&apos;t be forged cross-site, so it&apos;s CSRF-immune regardless of a
                  tag-along cookie. On an existing project, restart the API after pulling the fix.
                </p>
              </div>
            </div>

            {/* v3.31.56 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.56
                </span>
                <span className="text-sm text-muted-foreground">July 3, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Mobile code generation.</strong> <code>grit generate resource</code>{' '}
                  now scaffolds the mobile app too, not just the backend, web hooks, and admin.
                </p>

                <h3>Generated Expo screens &amp; hooks</h3>
                <p>
                  When a project has an <code>apps/expo</code> app, generating a resource also
                  writes a typed React Query hook (infinite-scroll list, single item, and
                  create/update/delete mutations), a <strong>paginated list screen</strong>, and a{' '}
                  <strong>detail screen</strong> — all field-aware (images become thumbnails, a{' '}
                  <code>belongs_to</code> renders its related record, dates/bools/files format
                  sensibly). A shared safe-area <code>ScreenHeader</code> with a back button ships
                  with the scaffold.
                </p>

                <h3>Relationship filtering</h3>
                <p>
                  Every <code>belongs_to</code> resource is now filterable by its foreign key —{' '}
                  <code>GET /products?category_id=…</code> returns just that parent&apos;s children,
                  and the generated hook takes the same filter. This is what powers a real
                  category → products browse flow on mobile.
                </p>

                <h3>Mobile app polish</h3>
                <p>
                  The scaffolded Expo app gained a full <strong>light/dark theme</strong> (default
                  light, with a working Settings toggle), the <strong>Grit logo</strong> on auth,
                  a floating glass tab bar lifted above the Android system nav, safe-area page
                  headers, wired-up <strong>Explore destinations</strong> (Users with pagination,
                  Notifications, Storage, Analytics, Content, Integrations), and{' '}
                  <strong>profile avatar upload + change password</strong>. Also fixed: the
                  physical-device API URL (derived from Expo&apos;s host), a splash-screen hang,
                  the auth token shape, and post-login navigation. Adds{' '}
                  <code>expo-image-picker</code>, <code>expo-linear-gradient</code>,{' '}
                  <code>expo-blur</code>, and <code>react-native-css-interop</code> — run{' '}
                  <code>pnpm i</code> then <code>npx expo start -c</code>.
                </p>
              </div>
            </div>

            {/* v3.31.55 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.55
                </span>
                <span className="text-sm text-muted-foreground">July 3, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Redesigned mobile auth &amp; navigation.</strong> The scaffolded
                  <code>--mobile</code> app now ships a premium, production-grade UI out of the box
                  instead of the plain starter screens.
                </p>

                <h3>Polished login &amp; register</h3>
                <p>
                  Both auth screens are rebuilt as a single elevated card on a faint architectural
                  grid: a gradient brand header, icon-prefixed inputs, a show/hide password toggle,
                  inline validation, a gradient primary CTA, and Google sign-in. Every action fires
                  haptic feedback and the card animates in with a spring <code>FadeInUp</code>.
                </p>

                <h3>Floating glass tab bar</h3>
                <p>
                  The bottom navigation is now a floating, rounded bar with a native frosted-blur
                  background on iOS (solid elevated surface on Android) and a selection haptic on
                  every tab switch. A new <code>PressableScale</code> primitive gives buttons the
                  tactile spring-press micro-interaction. New dependencies:
                  <code>expo-linear-gradient</code> and <code>expo-blur</code> — run
                  <code>pnpm i</code> then <code>npx expo start -c</code>.
                </p>
              </div>
            </div>

            {/* v3.31.54 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.54
                </span>
                <span className="text-sm text-muted-foreground">July 3, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Mobile styling fix (NativeWind).</strong> A scaffolded
                  <code>--mobile</code> app rendered completely unstyled — raw text on a black
                  screen — and hung on startup.
                </p>

                <h3>The Expo app was missing its Babel config</h3>
                <p>
                  NativeWind requires a <code>babel.config.js</code> with the
                  <code>jsxImportSource: &quot;nativewind&quot;</code> option and the
                  <code>nativewind/babel</code> preset — without it every <code>className</code> is
                  silently ignored. The scaffold never generated that file. It now ships one, plus
                  the <code>react-native-worklets</code> dependency and its Babel plugin (required
                  by Reanimated&nbsp;4, whose absence caused the startup hang), and
                  <code>web.bundler: &quot;metro&quot;</code> in <code>app.json</code>. On an
                  existing project, add <code>apps/expo/babel.config.js</code>, install
                  <code>react-native-worklets</code>, then restart Metro with a clear cache:
                  <code>npx expo start -c</code>.
                </p>
              </div>
            </div>

            {/* v3.31.53 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.53
                </span>
                <span className="text-sm text-muted-foreground">July 3, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Mobile (Expo) scaffold fixes.</strong> A fresh <code>--mobile</code>
                  app now starts cleanly with correct dependency versions and real app icons.
                </p>

                <h3>App icons &amp; splash now ship with the scaffold</h3>
                <p>
                  <code>app.json</code> referenced <code>./assets/icon.png</code> and
                  <code>./assets/splash.png</code>, but those files were never generated — so
                  Metro failed with <em>&ldquo;Unable to resolve asset&rdquo;</em>. The Grit logo
                  is now embedded in the CLI and written to <code>icon.png</code>,
                  <code>splash.png</code>, <code>adaptive-icon.png</code>, and
                  <code>favicon.png</code> on scaffold (with matching Android adaptive-icon and
                  web favicon entries in <code>app.json</code>). The same brand logo is also
                  dropped into the web, admin, and single-app <code>public/</code> folders.
                </p>

                <h3>Expo dependency versions aligned to SDK 54</h3>
                <p>
                  The Expo app pinned <code>expo</code> 54 but shipped SDK-53 versions of a few
                  packages, triggering compatibility warnings. Bumped
                  <code>expo-image</code> (~3.0.11), <code>expo-haptics</code> (~15.0.8),
                  <code>expo-web-browser</code> (~15.0.11), <code>react-native</code> (0.81.5),
                  and <code>typescript</code> (~5.9.2). On an existing project you can also run
                  <code>npx expo install --fix</code>.
                </p>
              </div>
            </div>

            {/* v3.31.52 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.52
                </span>
                <span className="text-sm text-muted-foreground">July 1, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Self-update fix.</strong> <code>grit update</code> could report
                  success but leave no <code>grit</code> on your PATH.
                </p>

                <h3>go install now targets grit&apos;s real location</h3>
                <p>
                  When grit was installed via the install script (into <code>~/.grit/bin</code>),
                  the update&apos;s <code>go install</code> step wrote the new binary to the Go
                  default <code>GOBIN</code> (<code>~/go/bin</code>) instead — while the running
                  binary had already been renamed aside. The result: <em>&ldquo;Updated to
                  vX&rdquo;</em> followed by <code>grit: No such file or directory</code>. The
                  updater now sets <code>GOBIN</code> to the directory grit actually lives in, so
                  the refreshed binary lands exactly where your PATH expects it. If you hit this on
                  an older version, recover with{' '}
                  <code>curl -fsSL https://gritframework.dev/install.sh | sh</code> (or the
                  PowerShell one-liner), then update as normal.
                </p>
              </div>
            </div>

            {/* v3.31.51 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.51
                </span>
                <span className="text-sm text-muted-foreground">July 1, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Windows dev fix.</strong> <code>grit start</code> no longer fails
                  to boot the frontend on machines where Turborepo&apos;s native binary can&apos;t load.
                </p>

                <h3>Dev no longer depends on the turbo binary</h3>
                <p>
                  On some Windows setups <code>turbo dev</code> exits with <code>0xC0000135</code>{' '}
                  (STATUS_DLL_NOT_FOUND) because turbo ships a platform-specific native binary that
                  needs the Visual C++ runtime — which blocked <code>grit start</code> on a fresh
                  machine. The scaffolded root <code>dev</code> script now uses pnpm&apos;s own
                  parallel runner (<code>pnpm --parallel --filter &quot;./apps/*&quot; --if-present run dev</code>),
                  which needs no native binary, so the web and admin dev servers always come up.
                  Turbo is kept for <code>build</code> / <code>lint</code> / <code>test</code> where
                  its caching helps.
                </p>
              </div>
            </div>

            {/* v3.31.50 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.50
                </span>
                <span className="text-sm text-muted-foreground">June 26, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Form-sharing polish.</strong> Four fixes
                  from one fresh-project test session, all on the{' '}
                  <code>/system/form-shares</code> surface.
                </p>

                <h3>1. Resource is now a dropdown, not a text input</h3>
                <p>
                  Typing <code>Catgeory</code> instead of{' '}
                  <code>Category</code> in the New Share modal
                  silently created a broken share -- the dispatcher
                  fell through to <code>default</code>, the public
                  form showed an empty state, and the operator
                  didn&apos;t find out until a customer hit it. The
                  modal now lists every registered resource in a{' '}
                  <code>&lt;select&gt;</code>, sourced from a new{' '}
                  <code>GET /api/admin/form-shares/resources</code>{' '}
                  endpoint.
                </p>

                <h3>2. Form preview with per-field hide toggles</h3>
                <p>
                  Once a resource is picked, the modal renders a
                  preview of the public form: every field with its
                  type, required/optional badge, and a{' '}
                  <em>Hide</em> checkbox for each optional one.
                  Required fields can&apos;t be hidden (the submit
                  would 422). The selected hidden keys persist with
                  the share and the public-form endpoint filters
                  them out server-side -- so anonymous visitors
                  never see a column the operator marked private.
                </p>

                <h3>3. Custom title + description on the public form</h3>
                <p>
                  The public form&apos;s heading used to be{' '}
                  <code>{`<resource> submission`}</code> + a
                  hardcoded <em>&ldquo;Fill out the form below to
                  submit a new <code>{`<resource>`}</code>.&rdquo;</em>{' '}
                  Both can now be customised per share in the New
                  Share modal. Title falls back through three
                  sources -- custom_title → label → resource name
                  -- so old shares keep working with the same
                  heading they had before.
                </p>

                <h3>4. The public form actually renders the right fields now</h3>
                <p>
                  <strong>This was the big one.</strong> v3.31.43
                  fixed{' '}
                  <code>services/form_share_dispatch.go</code> to
                  return the resource&apos;s real field schema
                  (via reflection) -- but the matching change to{' '}
                  <code>webPublicFormPage()</code> in the framework
                  scaffold never landed. Every project scaffolded
                  with v3.31.43 through v3.31.49 was shipping a
                  hardcoded name / email / phone / message contact
                  form, regardless of what the resource actually
                  looked like.
                </p>
                <p>
                  The scaffold&apos;s public form page is now the
                  fields-aware version: reads the new{' '}
                  <code>fields[]</code>{' '}, <code>custom_title</code>,
                  <code> custom_description</code> from the API
                  and renders one input per field with the right
                  HTML shape (text / email / tel / textarea /
                  number / checkbox / date / datetime / file).
                </p>
                <p>
                  For projects already scaffolded with the stale
                  page, the upgrade is a one-file copy:{' '}
                  <code>apps/web/app/forms/[token]/page.tsx</code>{' '}
                  from a fresh scaffold replaces the broken one.
                </p>

                <h3>Plus: Edit modal now scaffolds in fresh projects</h3>
                <p>
                  v3.31.43 added an Edit button to the form-shares
                  table -- but only as a hand-applied patch to the
                  ecom test project, never in the scaffold. v3.31.50
                  ships it properly: a <em>Pencil</em> button on
                  each row opens an Edit modal with the same
                  preview + hide toggles + title / description /
                  password controls.
                </p>

                <h3>Data model</h3>
                <p>
                  <code>models.FormShare</code> gains three columns:
                </p>
                <ul>
                  <li>
                    <code>custom_title</code> /{' '}
                    <code>custom_description</code> -- short strings
                  </li>
                  <li>
                    <code>hidden_fields</code> -- JSON array of
                    field keys to omit
                  </li>
                </ul>
                <p>
                  GORM AutoMigrate adds them on next boot. The new{' '}
                  <code>{`// grit:form-share:registered`}</code>{' '}
                  marker in <code>RegisteredResources()</code>{' '}
                  gets injected by the generator on each{' '}
                  <code>grit generate resource</code>; pre-v3.31.50
                  projects warn instead of failing.
                </p>
              </div>
            </div>

            {/* v3.31.49 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.49
                </span>
                <span className="text-sm text-muted-foreground">June 25, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Three small ergonomic wins from real
                  operator feedback.</strong> All three landed
                  together in v3.31.49.
                </p>

                <h3>1. Activity log shows the operator&apos;s real IP, not &quot;::1&quot;</h3>
                <p>
                  Local-dev activity rows showed{' '}
                  <code>::1</code> in the IP column because
                  gin&apos;s <code>ClientIP()</code> correctly
                  reports the IPv6 loopback for same-machine
                  traffic. Operators expect to see their actual
                  public IP.
                </p>
                <p>
                  The admin / web axios clients now fetch the
                  operator&apos;s public IP once per session
                  (cached in <code>sessionStorage</code>; sourced
                  from <code>api.ipify.org</code>) and attach it
                  as <code>X-Public-IP-Hint</code> on every API
                  call. The new{' '}
                  <code>services.ResolveClientIP</code> helper
                  honours the hint <em>only</em> when the TCP peer
                  is loopback, so production traffic from real
                  proxies (which sets <code>X-Forwarded-For</code>{' '}
                  for gin to consume) keeps using the trusted path
                  and can&apos;t be spoofed by a client header.
                </p>
                <p>
                  When the lookup fails (offline / ad-blocker), the
                  feed falls back to the prior behaviour and
                  renders <code>localhost (::1)</code> with the
                  raw value tucked next to it so the origin stays
                  inspectable.
                </p>

                <h3>2. Web navbar gets an Admin CTA back</h3>
                <p>
                  v3.31.42 replaced the navbar&apos;s Admin link
                  with the v3.31.42 UserMenu (Login / Sign up +
                  avatar dropdown). Operators landing on the
                  marketing site lost the one-click bounce to the
                  admin app and had to type the URL by hand.
                </p>
                <p>
                  v3.31.49 puts an Admin button back in the
                  navbar, both in the base scaffold (no-auth, post
                  v3.31.48) and in the auth-aware variant. Points
                  at <code>NEXT_PUBLIC_ADMIN_URL</code> (defaults
                  to <code>http://localhost:3001</code> for dev;
                  set to your prod admin origin before shipping).
                </p>

                <h3>3. Landing page surfaces all the dev URLs</h3>
                <p>
                  The <code>grit new</code> welcome banner prints
                  every URL the scaffold ships with: API, API
                  Docs, GORM Studio, Sentinel, Pulse, Admin,
                  MinIO, Mailhog. Once the terminal scrolls past,
                  operators have to dig back through history to
                  find the right one.
                </p>
                <p>
                  A new <code>{`<DevLinks />`}</code> component
                  renders all of them as a clickable grid at the
                  bottom of the web landing page, grouped by
                  function (App / API / Data / Ops) and colour-
                  coded. The whole section is wrapped in a{' '}
                  <code>NODE_ENV !== &quot;production&quot;</code>{' '}
                  check at module level so production marketing
                  pages never leak the internal port map -- the
                  section disappears from the prod bundle
                  entirely, not just hidden behind a class.
                </p>

                <h3>Files changed</h3>
                <ul>
                  <li>
                    Backend: new{' '}
                    <code>services/clientip.go</code>{' '}
                    (ResolveClientIP); inline mirror in{' '}
                    <code>middleware/activity.go</code> (HTTP
                    audit logger); CORS
                    Access-Control-Allow-Headers extended to
                    allow the hint.
                  </li>
                  <li>
                    Admin: <code>lib/api-client.ts</code> fetches
                    + caches the public IP, attaches the hint;
                    activity page&apos;s{' '}
                    <code>prettyIP</code> helper renders
                    &quot;localhost&quot; for loopback so
                    fall-back rows still read cleanly.
                  </li>
                  <li>
                    Web: <code>components/navbar.tsx</code> +{' '}
                    auth variant gain the Admin button;{' '}
                    <code>components/dev-links.tsx</code> new
                    file; <code>app/page.tsx</code> renders{' '}
                    <code>{`<DevLinks />`}</code> at the bottom.
                  </li>
                  <li>
                    Env: <code>NEXT_PUBLIC_ADMIN_URL</code>{' '}
                    documented in <code>.env</code> with the
                    default + a prod-deployment note.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.31.48 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.48
                </span>
                <span className="text-sm text-muted-foreground">June 25, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Two bug fixes from real user feedback on
                  v3.31.47.</strong> Both happened on fresh{' '}
                  <code>grit new</code> scaffolds: a broken Go file
                  from a marker collision, and web auth shipping by
                  default when it should be opt-in.
                </p>

                <h3>1. injectBefore now matches marker as a standalone line</h3>
                <p>
                  The generator&apos;s <code>injectBefore</code>{' '}
                  did a raw <code>strings.Index</code> on the
                  marker. Markers like{' '}
                  <code>// grit:form-share:fields</code> sometimes
                  appear inside the docstring of the function they
                  precede (&ldquo;...at the marker comment...&rdquo;).
                  The substring match landed there <em>first</em>{' '}
                  and the case got injected into the comment, not
                  the function body -- producing a syntax error.
                </p>
                <p>
                  The matcher now walks line-by-line and requires
                  the trimmed line content to equal the marker
                  exactly. Docstrings that mention the marker by
                  name are safe again, and the form-share
                  scaffolded comments stay readable.
                </p>

                <h3>2. Web auth is now opt-in via grit add web-auth</h3>
                <p>
                  The base web scaffold has been quietly shipping
                  the full auth surface since v3.28.1: login /
                  register / forgot-password / callback pages, the
                  five themed AuthShells, the useAuth hook,
                  AuthProvider, UserMenu, the web-session marker,
                  and Login / Sign up buttons in the navbar. That
                  was always intended to be{' '}
                  <strong>opt-in</strong> via{' '}
                  <code>grit add web-auth</code> -- the base
                  scaffold should be a clean marketing site with no
                  auth UI.
                </p>
                <p>
                  v3.31.48 moves all of those files into{' '}
                  <code>grit add web-auth</code>:
                </p>
                <ul>
                  <li>
                    Base scaffold: navbar shows Home / Blog / Docs
                    / GitHub only, no Login / Sign up. AppChrome
                    keeps <code>/forms/&lt;token&gt;</code> as the
                    only chromeless prefix (for public form-share).
                  </li>
                  <li>
                    <code>grit add web-auth</code> now writes
                    everything: hooks/use-auth.ts,
                    lib/auth-provider.tsx, lib/web-session.ts, the
                    four (auth) pages, the five themed shells,
                    UserMenu, middleware.ts, ProtectedWebRoute.tsx
                    -- and REPLACES the navbar + AppChrome with
                    their auth-aware variants (which add the
                    (auth) chromeless prefixes and the UserMenu in
                    the navbar). Replacement requires{' '}
                    <code>--force</code> for safety.
                  </li>
                </ul>

                <h3>Migrating existing projects</h3>
                <p>
                  Projects scaffolded with v3.31.x before this
                  release already have the auth files. They keep
                  working unchanged -- no removal happens
                  automatically. Future <code>grit new</code> calls
                  produce the clean base scaffold; if you need auth
                  on a fresh project, run{' '}
                  <code>grit add web-auth</code> right after{' '}
                  <code>grit new</code>.
                </p>

                <h3>Both bugs were reported same day</h3>
                <p>
                  The user spun up a fresh ecom-app, hit the
                  form-share dispatch syntax error on
                  <code>grit dev</code>, then noticed the web
                  shipping auth they didn&apos;t want. Both shipped
                  fixed in v3.31.48 within a few hours.
                </p>
              </div>
            </div>

            {/* v3.31.47 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.47
                </span>
                <span className="text-sm text-muted-foreground">June 25, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>The Preset Chart builder.</strong> Operators
                  can now build custom charts straight from Dashboard
                  Settings -- pick a resource, pick a preset, pick a
                  visualization. The charts render in the Charts
                  section alongside the system Activity + Severity
                  widgets. No SQL involved.
                </p>

                <h3>Four presets</h3>
                <p>
                  The presets cover the bulk of admin-dashboard needs
                  without introducing a query plane:
                </p>
                <ul>
                  <li>
                    <strong>Count over time</strong> -- daily count
                    of new records (no field needed)
                  </li>
                  <li>
                    <strong>Group by field</strong> -- top-N counts
                    grouped by a categorical column (e.g. orders by
                    status, products by category)
                  </li>
                  <li>
                    <strong>Sum over time</strong> -- daily sum of a
                    numeric column (e.g. revenue per day)
                  </li>
                  <li>
                    <strong>Avg over time</strong> -- daily average
                    of a numeric column (e.g. average order value)
                  </li>
                </ul>

                <h3>Five visualizations</h3>
                <p>
                  Each chart renders as <code>bar</code>,{' '}
                  <code>line</code>, <code>area</code>,{' '}
                  <code>pie</code>, or <code>donut</code> -- using
                  Recharts. The builder dims out incompatible
                  combinations (pie for a time-series, line for
                  group_by) so users see why a choice doesn&apos;t
                  make sense rather than picking a broken combo
                  and getting a flat chart.
                </p>

                <h3>How it works under the hood</h3>
                <p>
                  Same dispatch pattern as the v3.31.44 resource
                  stats. A new service file{' '}
                  <code>chart_dispatch.go</code> ships with a
                  switch over resource name + a reflective helper
                  that runs the right SQL for each preset:
                </p>
                <ul>
                  <li>
                    <code>count_over_time</code>: pulls timestamps
                    + buckets in-memory (portable across SQLite + Postgres)
                  </li>
                  <li>
                    <code>group_by</code>: SQL{' '}
                    <code>GROUP BY field ORDER BY COUNT(*) DESC LIMIT N</code>
                  </li>
                  <li>
                    <code>sum_over_time</code> /{' '}
                    <code>avg_over_time</code>: pulls (created_at, field)
                    pairs + aggregates in-memory
                  </li>
                </ul>
                <p>
                  Field whitelisting is the security boundary: the
                  helper reflects on the model to build two sets
                  (string/bool columns valid for group_by, numeric
                  columns valid for sum/avg) and rejects any field
                  not in the right set. The same dispatch marker
                  used by v3.31.44 (<code>{`// grit:resource-stats:dispatch`}</code>)
                  is reused, so one generator injection covers both
                  the sparkline + the chart presets for a new
                  resource.
                </p>

                <h3>New endpoint</h3>
                <p>
                  <code>GET /api/admin/dashboard/chart/:resource?preset=group_by&amp;field=status&amp;limit=10</code>
                  {' '}returns{' '}
                  <code>{`{ data: { preset, rows: [{x, y}], meta } }`}</code>.
                  The frontend ChartCard renders the right Recharts
                  component based on the saved viz; the{' '}
                  <code>{`{x, y}`}</code> shape works for all four
                  presets without a discriminator.
                </p>

                <h3>Data model</h3>
                <p>
                  <code>models.DashboardLayout</code> gains one new
                  JSON column:
                </p>
                <ul>
                  <li>
                    <code>custom_charts</code> -- array of user-defined
                    chart configs. The PUT handler validates each entry
                    on write (drops malformed rows individually rather
                    than rejecting the whole save).
                  </li>
                </ul>
                <p>
                  GORM AutoMigrate adds the column on next boot. No
                  manual migration; existing saved layouts continue to
                  work (empty array = no custom charts).
                </p>

                <h3>Frontend pieces</h3>
                <p>
                  Three new files in the admin scaffold:
                </p>
                <ul>
                  <li>
                    <code>components/dashboard/CustomChartCard.tsx</code>
                    {' '}-- renders one chart with Recharts. Loading +
                    error states inline so a broken chart never
                    blanks the section.
                  </li>
                  <li>
                    <code>components/dashboard/ChartBuilderForm.tsx</code>
                    {' '}-- the inline builder. Resource picker, preset
                    picker (with tiles), field picker (filtered by
                    preset), viz picker (with grey-out for
                    incompatible). Used in the new{' '}
                    <strong>Custom charts</strong> section on
                    Dashboard Settings.
                  </li>
                  <li>
                    Settings page <em>Custom charts</em> panel --
                    lists saved charts with Edit + Delete; click{' '}
                    <em>Add chart</em> to open the inline builder.
                  </li>
                </ul>

                <h3>Research note</h3>
                <p>
                  The design is the &ldquo;Preset Charts&rdquo;
                  pattern from the research pass (Metabase, Grafana,
                  Looker Studio, Superset, Power BI). It&apos;s the
                  table-first pattern with curated presets rather
                  than the dimension/metric drag-drop pattern --
                  fits Grit&apos;s convention-over-configuration
                  audience and slots straight into the v3.31.45
                  settings page. The dimension/metric builder
                  (Design B from the research) stays available as
                  a future v3.32 upgrade if users start asking for
                  one more axis of freedom.
                </p>
              </div>
            </div>

            {/* v3.31.46 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.46
                </span>
                <span className="text-sm text-muted-foreground">June 25, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Polished By-Resource Latest tables + per-resource layout toggle.</strong>{' '}
                  The v3.31.44 Latest list rendered a single
                  &ldquo;Name: X · Status: Y&rdquo; line per row,
                  which turned any FileRef column into a JSON blob
                  visible to the user. v3.31.46 swaps that for the
                  same column-driven table layout the resource list
                  page uses, with proper FileRef thumbnails, badges,
                  date formatting, and currency rendering.
                </p>

                <h3>The Latest table now uses renderCell</h3>
                <p>
                  Every cell in the dashboard&apos;s Latest table
                  now goes through the same{' '}
                  <code>renderCell</code> dispatch the resource list
                  pages already use. That means columns with{' '}
                  <code>format: &quot;image&quot;</code> render
                  thumbnails, columns with{' '}
                  <code>format: &quot;badge&quot;</code> render the
                  configured pill, dates and currency get their
                  normal formatting. The column picker still uses
                  the v3.31.44 heuristics (prefer{' '}
                  <code>name</code>/<code>title</code>/<code>email</code>/
                  <code>status</code>/<code>price</code>) but now
                  always reserves a slot for any image / FileRef
                  column so visual rows always have a thumbnail when
                  the model defines one.
                </p>

                <h3>Per-resource layout toggle: Split vs Tabs</h3>
                <p>
                  The v3.31.44 layout was hardcoded:{' '}
                  <strong>Split</strong> (Total card ~33% on the
                  left, Latest table ~67% on the right). That ratio
                  breaks down for resources with many columns or
                  long string values -- the Latest table never has
                  room to breathe. v3.31.46 adds a per-resource
                  layout mode:
                </p>
                <ul>
                  <li>
                    <strong>Split</strong> (default) -- the v3.31.44
                    side-by-side layout.
                  </li>
                  <li>
                    <strong>Tabs</strong> -- both widgets render
                    full-width inside a tabbed container. Two tabs
                    (<em>Total &lt;Resource&gt;</em> / <em>Latest
                    &lt;Resource&gt;</em>) with the Latest tab
                    opened by default since that&apos;s the widget
                    that benefits most from the extra width.
                  </li>
                </ul>
                <p>
                  Picked per resource in Dashboard Settings under
                  the new <em>Resource layout</em> panel (below the
                  By Resource checkboxes). Only resources with at
                  least one widget enabled show up in the picker --
                  the choice is moot otherwise.
                </p>

                <h3>Data model</h3>
                <p>
                  <code>models.DashboardLayout</code> gains one new
                  JSON column:
                </p>
                <ul>
                  <li>
                    <code>resource_layouts</code> -- a string-keyed
                    object (<code>{`{ "products": "tabs", "orders": "split" }`}</code>).
                    Only non-default (<code>tabs</code>) entries are
                    persisted; missing slugs fall back to{' '}
                    <code>split</code> at render time. The PUT
                    handler validates incoming values and silently
                    drops anything that isn&apos;t{' '}
                    <code>split</code> or <code>tabs</code>.
                  </li>
                </ul>
                <p>
                  GORM AutoMigrate adds the column on next boot. No
                  manual migration; existing saved layouts continue
                  to work (an empty map means &ldquo;every resource
                  uses split&rdquo; -- the v3.31.44 behaviour).
                </p>

                <h3>Coming in v3.31.47</h3>
                <p>
                  Next release will tackle the &ldquo;build a custom
                  chart&rdquo; ask -- give operators a way to pick a
                  resource + group-by field + aggregation +
                  visualization (bar/line/pie/donut) without writing
                  SQL. The design landed on the &ldquo;Preset
                  Charts&rdquo; pattern (count over time,{' '}
                  <em>group by field</em>, sum/avg over time, top-N)
                  -- those four cover the bulk of admin-dashboard
                  needs without introducing a query plane, and slot
                  straight into the same Dashboard Settings page as
                  the v3.31.45 toggles.
                </p>
              </div>
            </div>

            {/* v3.31.45 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.45
                </span>
                <span className="text-sm text-muted-foreground">June 25, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Per-resource dashboard customisation + section reordering.</strong>{' '}
                  The v3.31.44 &ldquo;By Resource&rdquo; band was
                  uncustomisable — it always rendered the Total +
                  Latest pair for every resource. Dashboard Settings
                  now exposes both halves per resource, and the four
                  top-level dashboard sections (Cards, Charts,
                  Tables, By Resource) can be reordered.
                </p>

                <h3>Per-resource toggles in Dashboard Settings</h3>
                <p>
                  A new <em>By Resource</em> section appears at the
                  bottom of <code>/settings/dashboard</code>, grouped
                  by resource. Each resource exposes two checkboxes:
                </p>
                <ul>
                  <li>
                    <strong>Total <em>&lt;Resource&gt;</em></strong>{' '}
                    — the stat card with the 30-day sparkline.
                  </li>
                  <li>
                    <strong>Latest <em>&lt;Resource&gt;</em></strong>{' '}
                    — the newest-N records table.
                  </li>
                </ul>
                <p>
                  Toggling either one off hides just that widget; the
                  row stretches the visible half to fill the
                  available width. Resources with both halves
                  unchecked don&apos;t render at all. The
                  resource-level <code>{`dashboard: { enabled: false }`}</code>{' '}
                  opt-out still exists for resources that should
                  never appear on the dashboard, even as catalog
                  entries.
                </p>

                <h3>Section reorder</h3>
                <p>
                  A new <strong>Section order</strong> panel sits at
                  the top of Dashboard Settings, showing the four
                  sections as a numbered list with up/down chevrons.
                  The saved order persists on the existing{' '}
                  <code>DashboardLayout</code> row (new{' '}
                  <code>section_order</code> column). The dashboard
                  page renders the sections in that order using CSS{' '}
                  <code>order</code> on a flex container — no JSX
                  restructure was needed.
                </p>

                <h3>Data model — two new JSON columns</h3>
                <p>
                  <code>models.DashboardLayout</code> gains two
                  fields, both JSON arrays:
                </p>
                <ul>
                  <li>
                    <code>resources</code> — enabled keys for the
                    By Resource band, formatted as{' '}
                    <code>{`"<slug>:total"`}</code> /{' '}
                    <code>{`"<slug>:latest"`}</code>. Same
                    presence-vs-absence semantics as the existing{' '}
                    <code>cards</code> / <code>charts</code> /{' '}
                    <code>tables</code> arrays: an empty list on a
                    saved row means &ldquo;hide everything&rdquo;;
                    a missing row means &ldquo;show defaults.&rdquo;
                  </li>
                  <li>
                    <code>section_order</code> — section keys in
                    render order. Default empty (= built-in order).
                    Unknown keys are silently dropped at render
                    time; missing default keys get appended to the
                    end so a saved layout from before a new section
                    was added still renders the new section.
                  </li>
                </ul>

                <h3>Backward compatibility</h3>
                <p>
                  Pre-v3.31.45 projects don&apos;t have the new
                  columns. GORM AutoMigrate adds them on next boot.
                  Existing saved layouts continue to work — both new
                  arrays default to empty, which means &ldquo;use
                  built-in defaults&rdquo; (all resource widgets
                  shown, default section order). Frontend{' '}
                  <code>SavedLayout</code> gains the two fields as
                  required (TypeScript-side); the wire shape allows
                  them to be omitted on the PUT body, treated as
                  empty.
                </p>
              </div>
            </div>

            {/* v3.31.44 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.44
                </span>
                <span className="text-sm text-muted-foreground">June 25, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Per-resource dashboard widgets, scoped by DateFilter.</strong>{' '}
                  Every newly generated resource now ships with two
                  preset widgets on the main dashboard: a{' '}
                  <strong>Total stat card with a 30-day sparkline</strong>{' '}
                  on the left and a <strong>Latest 5 records</strong>{' '}
                  preview on the right. Both honor the existing
                  dashboard <code>DateFilter</code> so the count
                  obeys whichever range the operator has selected.
                </p>

                <h3>How it works</h3>
                <p>
                  Three pieces ship together:
                </p>
                <ul>
                  <li>
                    <strong>Service</strong>:{' '}
                    <code>services.ComputeResourceStats</code> in{' '}
                    <code>apps/api/internal/services/resource_stats_dispatch.go</code>{' '}
                    — a generator-driven switch over resource name,
                    each case calls a single reflective helper that
                    counts rows in the active range, builds a
                    30-day sparkline, and lists the newest N (JSON
                    round-tripped so <code>json:&quot;-&quot;</code>{' '}
                    columns like <code>PasswordHash</code> never leak).
                  </li>
                  <li>
                    <strong>Endpoint</strong>:{' '}
                    <code>GET /api/admin/dashboard/resource-stats/:resource</code>{' '}
                    — accepts the same{' '}
                    <code>created_since</code> /{' '}
                    <code>created_from</code> /{' '}
                    <code>created_to</code> params the resource list
                    pages already use, so the wire shape matches.
                  </li>
                  <li>
                    <strong>Widgets</strong>:{' '}
                    <code>ResourceStatCard</code>,{' '}
                    <code>ResourceLatestTable</code>, and a thin{' '}
                    <code>ResourceWidgetsRow</code> wrapper. The
                    dashboard page maps over registered resources
                    and renders one row per resource below the
                    existing Quick Access section.
                  </li>
                </ul>

                <h3>Sparkline window is always 30 days</h3>
                <p>
                  The sparkline ignores the active date filter on
                  purpose — under the &ldquo;Today&rdquo; preset
                  it would collapse to a single bar, which carries
                  no information. The total + latest list still
                  obey the filter; only the trend chart is fixed.
                </p>

                <h3>Opt-out per resource</h3>
                <p>
                  Resources can hide their widgets by setting{' '}
                  <code>dashboard: {`{ enabled: false }`}</code>{' '}
                  in the resource definition. The flag is opt-out
                  by design: a new resource is more often than not
                  worth showing on the dashboard.
                </p>

                <h3>Backward compatibility</h3>
                <p>
                  The generator injects a switch case into{' '}
                  <code>resource_stats_dispatch.go</code> on each{' '}
                  <code>grit generate</code> run, at the marker{' '}
                  <code>{`// grit:resource-stats:dispatch`}</code>.
                  Projects scaffolded before v3.31.44 don&apos;t
                  have the file or the marker — the generator
                  detects this and prints a one-line warning
                  instead of failing. Patch existing projects by
                  copying the scaffold file from the framework
                  repo, then re-running{' '}
                  <code>grit generate</code> to populate the cases
                  (or hand-edit them).
                </p>
              </div>
            </div>

            {/* v3.31.43 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.43
                </span>
                <span className="text-sm text-muted-foreground">June 25, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Form-share polish: matching public form + editable shares.</strong>{' '}
                  Two small but high-impact fixes on top of the
                  v3.31.41 form-share generator. Both ship through
                  the framework scaffold and the generator, so new
                  resources pick them up automatically and existing
                  projects get the imports added lazily on the
                  next <code>grit generate</code>.
                </p>

                <h3>1. Public form renders the resource&apos;s actual fields</h3>
                <p>
                  Before this release the public share page at{' '}
                  <code>/forms/[token]</code> rendered a hardcoded{' '}
                  Name + Email + Phone + Message contact form for
                  every resource. Creating a share for a Category
                  with <code>name</code> + <code>image</code>{' '}
                  fields still showed the contact form — and the
                  Name field happened to line up purely by
                  coincidence. Submitting any other field shape was
                  effectively impossible.
                </p>
                <p>
                  The dispatcher now exports{' '}
                  <code>services.PublicFields(resourceName)</code>,
                  a per-resource switch that reflects the model
                  struct and returns a typed{' '}
                  <code>PublicFieldInfo[]</code> with one entry per
                  user-facing column (framework + auto fields are
                  skipped). The HTTP type for each field is
                  inferred from the Go type:{' '}
                  <code>FileRef → file</code>,{' '}
                  <code>time.Time → datetime</code>,{' '}
                  <code>bool → checkbox</code>,{' '}
                  <code>int/float → number</code>, and{' '}
                  <code>string</code> with name heuristics for
                  email / phone / textarea fields.
                </p>
                <p>
                  <code>GET /api/public/forms/:token</code> now
                  returns a <code>fields[]</code> array alongside{' '}
                  <code>resource_name</code> and{' '}
                  <code>has_password</code>. The web page consumes
                  it and renders one input per field, with proper
                  shapes for checkbox, number, textarea, and
                  date/datetime. File fields render an inline
                  &ldquo;File uploads aren&apos;t supported on
                  public-share forms&rdquo; explainer instead of an
                  unusable input — file uploads require the
                  auth-gated <code>/api/uploads</code> endpoint and
                  aren&apos;t supported on anonymous shares yet.
                </p>

                <h3>2. Admin can edit existing shares</h3>
                <p>
                  The admin form-shares page already had Audit /
                  Copy / Open / Delete buttons but no way to change
                  a share&apos;s label or password protection after
                  creation. Want to add a password to an existing
                  link? Delete + recreate, and re-distribute the
                  new token to every recipient.
                </p>
                <p>
                  A new <strong>Edit</strong> button opens a modal
                  with three controls:
                </p>
                <ul>
                  <li>
                    <strong>Label</strong> — free text, optional.
                  </li>
                  <li>
                    <strong>Password mode</strong> — three pills:{' '}
                    <em>Keep current</em>, <em>Set password</em>,{' '}
                    <em>Remove password</em>. &ldquo;Remove&rdquo;
                    is disabled when the share has no password.
                  </li>
                  <li>
                    <strong>New password</strong> — shown only when
                    mode is &ldquo;Set password&rdquo;.
                  </li>
                </ul>
                <p>
                  The backend handler at{' '}
                  <code>PATCH /api/admin/form-shares/:id</code>{' '}
                  already supported the full payload (it accepts{' '}
                  <code>password: &quot;-&quot;</code> as the
                  sentinel for &ldquo;remove&rdquo;); this release
                  just adds the missing UI to call it.
                </p>

                <h3>Backward compatibility</h3>
                <p>
                  Projects scaffolded before v3.31.43 don&apos;t
                  have the <code>{`// grit:form-share:fields`}</code>{' '}
                  marker or the <code>reflect</code> +{' '}
                  <code>strings</code> imports the new code
                  depends on. The generator now adds the imports
                  lazily on the first generated resource and prints
                  a one-line warning when the marker is missing,
                  pointing operators at a manual patch. Existing
                  shares keep working — they just continue to show
                  the hardcoded form until the project is
                  re-scaffolded or the dispatcher is patched.
                </p>
              </div>
            </div>

            {/* v3.31.42 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.42
                </span>
                <span className="text-sm text-muted-foreground">June 25, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Three web-auth fixes shipped together.</strong>{' '}
                  All surfaces — admin scaffold, web scaffold,{' '}
                  <code>grit add web-auth</code>, and the generator —
                  picked up the changes so existing and new projects
                  both benefit.
                </p>

                <h3>1. Auth pages render full-bleed (no navbar / footer)</h3>
                <p>
                  Until now the web app&apos;s <code>(auth)/login</code>,{' '}
                  <code>(auth)/register</code>, <code>(auth)/forgot-password</code>,{' '}
                  <code>(auth)/callback</code>, and <code>forms/[token]</code>{' '}
                  pages all rendered inside the root layout, which
                  pinned the <code>&lt;Navbar /&gt;</code> +{' '}
                  <code>&lt;Footer /&gt;</code> to the top and bottom.
                  The auth pages already supply their own{' '}
                  <code>AuthShell</code> chrome (same template the
                  admin uses) so the result was visually doubled.
                </p>
                <p>
                  New <code>components/AppChrome.tsx</code> is a tiny
                  client wrapper that conditionally renders Navbar +
                  Footer based on the pathname. The root layout drops
                  to a thin server component again. Auth + public
                  form-share pages render full-bleed; everything else
                  is unchanged.
                </p>

                <h3>2. grit_web_session marker stops admin sessions from unlocking web pages</h3>
                <p>
                  <code>grit_access</code> is set by the API on the
                  API origin (<code>localhost:8080</code>). That same
                  cookie is also used by <code>apps/admin</code> — so
                  an operator who signed in via the admin app could
                  open <code>apps/web/account</code> in the same
                  browser and walk straight past the web middleware.
                  The API call from <code>useMe()</code> succeeded
                  (the cookie is valid), and{' '}
                  <code>ProtectedWebRoute</code> happily rendered the
                  page.
                </p>
                <p>
                  New <code>grit_web_session</code> marker cookie is
                  set by the web app&apos;s own login/register flow
                  on the WEB origin (<code>localhost:3000</code>) via
                  client JS — non-HttpOnly, intentionally — and
                  cleared by logout. Middleware reads{' '}
                  <code>grit_web_session</code> instead of{' '}
                  <code>grit_access</code>. Admin-only sessions never
                  stamp the marker, so the web gates bounce them. The
                  real session security is unchanged: useMe() still
                  validates the API JWT; the marker is just a fast
                  edge check.
                </p>
                <p>
                  Mechanically: <code>lib/web-session.ts</code> with{' '}
                  <code>setWebSessionMarker</code> /{' '}
                  <code>clearWebSessionMarker</code> /{' '}
                  <code>hasWebSessionMarker</code>; called from{' '}
                  <code>useLogin</code> + <code>useRegister</code>{' '}
                  (onSuccess), <code>useLogout</code> (onSettled),
                  and the direct-submit{' '}
                  <code>(auth)/login</code> + <code>(auth)/register</code>{' '}
                  pages.
                </p>

                <h3>3. Navbar UserMenu replaces the Admin button</h3>
                <p>
                  The web navbar used to ship with a single{' '}
                  <em>Admin</em> link that punted everyone to the
                  admin app&apos;s login. Replaced with{' '}
                  <code>components/UserMenu.tsx</code>:
                </p>
                <ul>
                  <li>
                    <strong>Signed out</strong> — <em>Log in</em> +{' '}
                    <em>Sign up</em> buttons that link to the web
                    app&apos;s own auth flow.
                  </li>
                  <li>
                    <strong>Signed in</strong> — avatar dropdown with
                    name + email at top, <em>Account</em> link,{' '}
                    <em>Sign out</em> button.
                  </li>
                  <li>
                    <strong>Loading</strong> — a placeholder of the
                    same width as the signed-out CTA pair so the
                    navbar doesn&apos;t shift on{' '}
                    <code>useMe()</code> resolve.
                  </li>
                </ul>

                <h3>Bonus: web-hook generator now imports FileRef</h3>
                <p>
                  Same TS2304 fix shipped in v3.31.37 for{' '}
                  <code>writeTSTypes</code>, applied to{' '}
                  <code>writeReactQueryHooks</code>. Generated
                  <code> apps/web/hooks/use-X.ts</code> files with{' '}
                  <code>:file:</code> / <code>:files:</code> fields now
                  emit{' '}
                  <code>{`import type { FileRef } from "@repo/shared/schemas"`}</code>{' '}
                  at the top.
                </p>

                <h3>Migration</h3>
                <p>
                  Run <code>grit upgrade</code> or hand-patch:
                </p>
                <ul>
                  <li>
                    Drop <code>apps/web/components/AppChrome.tsx</code>{' '}
                    + <code>UserMenu.tsx</code> + <code>lib/web-session.ts</code>{' '}
                    in.
                  </li>
                  <li>
                    Update <code>app/layout.tsx</code> to render{' '}
                    <code>&lt;AppChrome&gt;</code> instead of{' '}
                    <code>&lt;Navbar /&gt; ... &lt;Footer /&gt;</code>.
                  </li>
                  <li>
                    Update <code>middleware.ts</code> to read{' '}
                    <code>grit_web_session</code>.
                  </li>
                  <li>
                    Add the marker calls to{' '}
                    <code>hooks/use-auth.ts</code> +{' '}
                    <code>(auth)/login/page.tsx</code> +{' '}
                    <code>(auth)/register/page.tsx</code>.
                  </li>
                  <li>
                    Replace the Admin button in{' '}
                    <code>components/navbar.tsx</code> with{' '}
                    <code>&lt;UserMenu /&gt;</code>.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.31.41 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.41
                </span>
                <span className="text-sm text-muted-foreground">June 25, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>System Hub tile for Public form sharing +
                  two course lessons rewritten to match what actually
                  happens.</strong> No code-behavior changes — this
                  release closes the gap between the documentation
                  and the running system.
                </p>

                <h3>System Hub tile</h3>
                <p>
                  The <code>/system/form-shares</code> page has
                  existed since v3.31.20 but wasn&apos;t reachable
                  from the <code>/system</code> tile grid — operators
                  had to know the URL by heart. The course lesson
                  pointed users at &quot;System → Public form
                  sharing&quot; which only existed as a sidebar
                  shortcut, not as a Hub tile. Now both surfaces
                  carry it. The tile shows up in all four System Hub
                  variants (default <code>admin_v331_files.go</code>{' '}
                  + minimal/modern/glass via{' '}
                  <code>admin_v3315_pages.go</code>).
                </p>

                <h3>public-form-sharing lesson</h3>
                <p>
                  Fixed the navigation reference (System Hub tile vs
                  direct route) and added a new section on file
                  upload limitations. The default{' '}
                  <code>/forms/[token]</code> template renders text
                  inputs only — <code>:file:</code> /{' '}
                  <code>:files:</code> / <code>:image:</code>{' '}
                  columns are silently skipped because{' '}
                  <code>/api/uploads</code> is auth-gated. The lesson
                  now walks through three production-shaped
                  workarounds (presigned URLs, external links,
                  magic-link auth).
                </p>

                <h3>grit-expose lesson</h3>
                <p>
                  Added a deep &quot;Behind the scenes: how the auth
                  bypass actually works&quot; section that walks the
                  four security layers operators should understand
                  before shipping a <code>--public-share</code> form
                  to production:
                </p>
                <ul>
                  <li>
                    The routes live in a separate Gin group
                    (<code>publicForms := r.Group(&quot;/api/public/forms&quot;)</code>) with no{' '}
                    <code>middleware.Auth</code>.
                  </li>
                  <li>
                    Sentinel rate-limits each token aggressively by IP.
                  </li>
                  <li>
                    The dispatcher service is the real security
                    boundary — only resources with an explicit{' '}
                    <code>case</code> in the switch can be created;
                    unknown keys in the request body are silently
                    dropped at the typed-struct decode.
                  </li>
                  <li>
                    The optional bcrypt password on the FormShare row
                    is the fourth layer.
                  </li>
                </ul>
                <p>
                  Plus the same file-upload caveat with a copy of the
                  three workaround paths.
                </p>

                <h3>Migration</h3>
                <p>
                  Run <code>grit upgrade</code> to pull the System
                  Hub tile into an existing project. The
                  <code> tile</code> is purely additive — no
                  behaviour change for anything that was already
                  working.
                </p>
              </div>
            </div>

            {/* v3.31.40 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.40
                </span>
                <span className="text-sm text-muted-foreground">June 25, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Per-user dashboard customisation + dashboard
                  date filter.</strong> A new Settings → Dashboard page
                  lets each operator pick which stat cards, charts, and
                  tables show up on their dashboard, grouped by module.
                  The dashboard also gets a date-window filter that
                  scopes every stat and chart to the selected range.
                </p>

                <h3>Backend</h3>
                <ul>
                  <li>
                    New model <code>DashboardLayout</code> (one row
                    per user, unique on user_id) with three{' '}
                    <code>JSONSlice[string]</code> columns for cards /
                    charts / tables plus a <code>date_preset</code>{' '}
                    text column for the persisted window.
                  </li>
                  <li>
                    New handler exposing{' '}
                    <code>GET /api/dashboard-layout</code> (returns
                    the current user&apos;s saved layout or a
                    zero-valued struct if none) and{' '}
                    <code>PUT /api/dashboard-layout</code> (whole-row
                    replace).
                  </li>
                  <li>
                    Empty layout (id === &quot;&quot;) = &quot;show
                    all widgets&quot;. Saved layout with{' '}
                    <code>cards: []</code> = &quot;hide every stat
                    card&quot;. The frontend distinguishes the two by
                    checking <code>layout.id</code>.
                  </li>
                </ul>

                <h3>Widget catalog</h3>
                <p>
                  New <code>lib/dashboard-catalog.ts</code> aggregates
                  the catalog of pickable widgets from two sources:
                </p>
                <ul>
                  <li>
                    <strong>System</strong> widgets (Users, Events 24h,
                    Notifications, Resources count, Activity 7-day
                    chart, Severity mix, Recent activity, Quick
                    access) — these are the legacy hard-coded
                    dashboard tiles, now opt-out-able per user.
                  </li>
                  <li>
                    <strong>Per-resource</strong> widgets — every
                    entry in a <code>ResourceDefinition.dashboard.widgets</code>{' '}
                    array contributes one catalog entry, grouped
                    under the resource&apos;s module name.
                  </li>
                </ul>
                <p>
                  Widget keys are stable strings (<code>system:users</code>,{' '}
                  <code>products:total-products</code>, etc.) so the
                  saved layout doesn&apos;t break when widget order
                  changes in the definition.
                </p>

                <h3>Settings page</h3>
                <p>
                  At <code>/settings/dashboard</code>: three sections
                  (Cards / Charts / Tables), each with checkbox lists
                  grouped by module. Per-section{' '}
                  <em>Select all</em> /{' '}
                  <em>Deselect all</em>; per-module{' '}
                  <em>All</em> / <em>None</em> for fine-grained
                  configuration. Sidebar nav gets a new{' '}
                  &quot;Dashboard settings&quot; entry under System
                  (no admin gate; every user can customise their own
                  view).
                </p>

                <h3>Dashboard date filter</h3>
                <p>
                  The existing <code>DateFilter</code> component from
                  v3.31.34 is now on the dashboard too. URL-persisted
                  via <code>?date=preset</code> /{' '}
                  <code>?date_from</code> / <code>?date_to</code>;
                  initial value falls back to the saved{' '}
                  <code>date_preset</code> so a refresh keeps the
                  window. Every system widget query keys on the
                  active <code>dateParams</code> so changing the
                  filter retriggers a refetch with{' '}
                  <code>?created_since=7d</code> /{' '}
                  <code>?created_from=...&amp;created_to=...</code>{' '}
                  appended.
                </p>

                <h3>Migration</h3>
                <p>
                  New scaffolded projects ship everything wired up. To
                  add to an existing project, run <code>grit upgrade</code>{' '}
                  (which writes the new files), then hand-add the
                  model to <code>models/user.go</code>&apos;s{' '}
                  <code>Models()</code> list and the routes to{' '}
                  <code>routes.go</code> — or rerun the upgrade with{' '}
                  <code>--force</code> if you haven&apos;t customised
                  those files. The existing hand-coded dashboard
                  page keeps working with no changes; the new
                  filtering activates only after you replace
                  <code>app/(dashboard)/dashboard/page.tsx</code>{' '}
                  with the v3.31.40 template (it reads{' '}
                  <code>useDashboardLayout()</code> and{' '}
                  <code>resolveEnabledKeys()</code>).
                </p>
                <p>
                  Heads-up for early adopters: the v3.31.40 framework
                  scaffold ships the building blocks (model + handler
                  + catalog + hook + Settings page + sidebar entry)
                  but does <em>not</em> auto-rewrite the dashboard
                  page — that lands in a follow-up release once the
                  multi-style dashboard variants (default / modern /
                  minimal / glass) are all refactored to use the new
                  layout reader. Until then, the example dashboard
                  refactor in the docs walks you through the changes
                  by hand.
                </p>
              </div>
            </div>

            {/* v3.31.39 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.39
                </span>
                <span className="text-sm text-muted-foreground">June 24, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>CUD activity logging on every generated
                  resource.</strong> Until this release the Activity
                  feed only carried sign-ins and sign-outs — every
                  Create / Update / Delete on a generated resource
                  went unrecorded. Now each one writes a row with a
                  human-readable summary using a fixed format
                  convention.
                </p>

                <h3>Format convention</h3>
                <pre><code>{`{verb} {entityType} {identifier}[: {detail}]`}</code></pre>
                <p>
                  <code>identifier</code> is the human-readable label
                  (name, title, slug, sku); it must never be blank —
                  the helper falls back to <code>(unnamed)</code> if
                  the caller hands it an empty string. <code>detail</code>
                  is optional extra context (price for Create, diff
                  for Update) and only renders when non-empty.
                </p>
                <p>Example rows in the feed:</p>
                <ul>
                  <li><code>Created Product Desktop: KES 340,000</code></li>
                  <li><code>Updated Product Desktop: changed name, price</code></li>
                  <li><code>Updated Category Phones: image changed</code></li>
                  <li><code>Deleted Blog &quot;Welcome to the new site&quot;</code></li>
                </ul>

                <h3>Three new helpers in services/activity.go</h3>
                <ul>
                  <li>
                    <code>LogCreate(db, c, entityType, identifier, resourceID, detail)</code>
                  </li>
                  <li>
                    <code>LogUpdate(db, c, entityType, identifier, resourceID, detail)</code>
                  </li>
                  <li>
                    <code>LogDelete(db, c, entityType, identifier, resourceID)</code>
                  </li>
                </ul>
                <p>
                  Plus <code>DiffSummary(updates)</code> for rendering
                  a GORM Updates() map as a sorted, deterministic
                  diff string (1 field → <code>field changed</code>;
                  2–3 fields → <code>changed a, b, c</code>; 4+ →{' '}
                  <code>N fields changed (a, b, c, ...)</code>).
                  Errors are logged, never returned — losing an audit
                  row should not fail a real request.
                </p>

                <h3>Generator emits the calls automatically</h3>
                <p>
                  Every <code>grit generate resource</code> from
                  v3.31.39 onward inserts:
                </p>
                <ul>
                  <li>
                    <code>services.LogCreate(...)</code> after a
                    successful <code>Create</code>
                  </li>
                  <li>
                    <code>services.LogUpdate(...)</code> after{' '}
                    <code>Update</code> and <code>Patch</code> (the
                    grouped Save handler from v3.31.18 is logged the
                    same way)
                  </li>
                  <li>
                    <code>services.LogDelete(...)</code> after{' '}
                    <code>Delete</code>
                  </li>
                </ul>
                <p>
                  Identifier expression is picked at generation time
                  from the model&apos;s fields, in priority order:{' '}
                  <code>Name</code>, <code>Title</code>,{' '}
                  <code>Slug</code>, <code>SKU</code>,{' '}
                  <code>Subject</code>, <code>Label</code>,{' '}
                  <code>Email</code>. Falls back to{' '}
                  <code>item.ID</code> so the log line is never blank.
                </p>

                <h3>Migration</h3>
                <p>
                  Re-run <code>grit generate resource X</code> for
                  each existing resource (the file rewrite picks up
                  the new log calls), or hand-patch each handler:
                </p>
                <ul>
                  <li>
                    Add{' '}
                    <code>{`"<module>/internal/services"`}</code> to
                    the imports.
                  </li>
                  <li>
                    Drop{' '}
                    <code>services.LogCreate/Update/Delete(...)</code>{' '}
                    calls right after each success path, before the
                    <code>c.JSON(...)</code>.
                  </li>
                  <li>
                    Use <code>services.DiffSummary(updates)</code> for
                    the diff string in Update / Patch.
                  </li>
                </ul>
                <p>
                  Existing auth helpers (<code>LogLogin</code>,{' '}
                  <code>LogRegister</code>, <code>LogLogout</code>)
                  are unchanged — they keep their semantic
                  <code>auth.X</code> action names.
                </p>
              </div>
            </div>

            {/* v3.31.38 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.38
                </span>
                <span className="text-sm text-muted-foreground">June 24, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Auto comma-formatting on number inputs.</strong>{' '}
                  Typing <code>3000</code> in a price field now reads{' '}
                  <code>3,000</code> on screen the moment the fourth
                  digit lands. Helps catch zero-count mistakes (the
                  &quot;is that 30k or 300k?&quot; problem) without
                  changing the wire shape — the form still submits a
                  plain JS number, and the API still receives an int
                  or float as before.
                </p>

                <h3>How it works</h3>
                <p>
                  <code>NumberField</code> switched from{' '}
                  <code>type=&quot;number&quot;</code> to{' '}
                  <code>type=&quot;text&quot;</code> with{' '}
                  <code>inputMode=&quot;decimal&quot;</code> (or{' '}
                  <code>&quot;numeric&quot;</code> for int/uint columns).
                  The visible value is the comma-formatted string; the
                  field also keeps a parsed JS number in form state.
                  Mobile keyboards still pop up correctly thanks to{' '}
                  <code>inputMode</code>; the comma can render literally
                  because text inputs don&apos;t strip non-digits.
                </p>
                <p>
                  Cursor position is preserved across the reformat by
                  counting non-comma characters before the caret and
                  restoring after the same count in the new value, so
                  editing in the middle of a number doesn&apos;t fling
                  the caret to the end.
                </p>

                <h3>numberKind hint</h3>
                <p>
                  New <code>FieldDefinition</code> knob:
                </p>
                <pre><code>{`numberKind?: "int" | "uint" | "float"`}</code></pre>
                <p>
                  Tells the input which characters to accept:
                </p>
                <ul>
                  <li><code>int</code> — negatives yes, decimals no</li>
                  <li><code>uint</code> — neither negatives nor decimals</li>
                  <li><code>float</code> — both (legacy permissive default when unset)</li>
                </ul>
                <p>
                  The generator now emits the right{' '}
                  <code>numberKind</code> for every number field based on
                  the Go column type. <code>grit sync</code> also adds
                  it when injecting newly-added Go fields into existing
                  admin resource files. Hand-written resources that
                  don&apos;t set <code>numberKind</code> stay
                  permissive — no breaking change to your existing
                  forms.
                </p>

                <h3>Edge cases handled</h3>
                <ul>
                  <li><strong>Paste &quot;$1,234.56&quot;</strong> — strips the dollar, keeps the value</li>
                  <li><strong>Mid-typing &quot;3000.&quot;</strong> — preserves the trailing dot so the user can finish the decimal</li>
                  <li><strong>Backspace through a comma</strong> — the comma re-inserts after the digit is removed; caret tracks digit count, not column position</li>
                  <li><strong>Leading zeros</strong> — &quot;0123&quot; collapses to &quot;123&quot;; &quot;0&quot; stays &quot;0&quot; so &quot;0.5&quot; is reachable</li>
                  <li><strong>External value sync</strong> — opening Edit on an existing record formats the loaded number; subsequent typing skips the sync to avoid stomping mid-edit state</li>
                </ul>

                <h3>Migration</h3>
                <p>
                  Replace <code>apps/admin/components/forms/fields/number-field.tsx</code>{' '}
                  with the rewritten file and add the optional{' '}
                  <code>numberKind</code> knob to{' '}
                  <code>apps/admin/lib/resource.ts</code> on{' '}
                  <code>FieldDefinition</code>. Existing resource files
                  keep working — <code>numberKind</code> defaults to{' '}
                  <code>float</code> when unset, which gives the legacy
                  permissive behaviour. Run{' '}
                  <code>grit sync</code> on any project to backfill
                  <code>numberKind</code> for new fields the generator
                  finds; existing fields aren&apos;t touched.
                </p>
              </div>
            </div>

            {/* v3.31.37 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.37
                </span>
                <span className="text-sm text-muted-foreground">June 24, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Bug fix: opening a Create form on a resource
                  with a <code>:files:</code> column no longer crashes
                  into the global error boundary.</strong> Same release
                  also tidies up three companion TS errors that were
                  red-squiggling in IDEs even though SWC stripped them
                  at runtime.
                </p>

                <h3>What was broken</h3>
                <p>
                  <code>buildDefaults</code> in <code>form-builder.tsx</code>{' '}
                  seeded every non-toggle field to <code>&quot;&quot;</code>{' '}
                  (empty string). For <code>files</code> /{' '}
                  <code>images</code> / <code>videos</code> types,
                  react-hook-form&apos;s initial state was therefore a
                  string. The matching field component (<code>FilesField</code>,{' '}
                  <code>ImagesField</code>, <code>VideosField</code>)
                  immediately called <code>.map()</code> on that
                  &quot;array&quot; — strings have no <code>.map</code>{' '}
                  → TypeError → the parent <code>FormSheet</code> blew
                  up into Next.js&apos; error boundary. The user saw
                  &quot;Something went wrong&quot;.
                </p>

                <h3>The fix</h3>
                <p>
                  <code>buildDefaults</code> now branches by field type:
                  arrays default to <code>[]</code>, nullable objects
                  (<code>file</code> / <code>image</code> /{' '}
                  <code>video</code>) default to <code>null</code>,
                  toggles stay <code>false</code>, everything else
                  stays <code>&quot;&quot;</code>.
                </p>
                <pre><code>{`const ARRAY_FIELD_TYPES = new Set([
  "files", "images", "videos", "multi-relationship-select",
]);
const NULLABLE_OBJECT_FIELD_TYPES = new Set([
  "file", "image", "video",
]);
// ...
} else if (ARRAY_FIELD_TYPES.has(field.type)) {
  defaults[field.key] = [];
} else if (NULLABLE_OBJECT_FIELD_TYPES.has(field.type)) {
  defaults[field.key] = null;
}`}</code></pre>
                <p>
                  All three field components also got a defensive{' '}
                  <code>Array.isArray</code> guard so a stale form or
                  a deserialised-wrong API response can&apos;t crash
                  the dropzone — the field just renders empty and the
                  user can still upload.
                </p>

                <h3>TypeScript clean-up</h3>
                <ul>
                  <li>
                    <code>ColumnFormat</code> union now includes{' '}
                    <code>&quot;file&quot;</code> and{' '}
                    <code>&quot;files&quot;</code> — both renderers
                    were already implemented but the type union was
                    stale, so resource definitions emitted by the
                    v3.31.30+ generator flagged a TS2322 on every
                    file column.
                  </li>
                  <li>
                    Generated <code>packages/shared/types/&lt;model&gt;.ts</code>{' '}
                    files now import <code>FileRef</code> from{' '}
                    <code>schemas/file-ref</code> when any field is{' '}
                    <code>:file:</code> or <code>:files:</code> —
                    previously the type was referenced without an
                    import, fine at runtime but red squiggles in the
                    IDE.
                  </li>
                  <li>
                    <code>ImportModal</code>&apos;s narrowing check
                    against <code>resource.table.import</code> dropped
                    the redundant <code>!== false</code> comparison
                    (TS2367) — a plain truthy check correctly handles
                    all three values of the union.
                  </li>
                </ul>

                <h3>Migration</h3>
                <p>
                  Run <code>grit upgrade</code>, or hand-patch:
                </p>
                <ul>
                  <li>
                    <code>apps/admin/components/forms/form-builder.tsx</code>{' '}
                    — update <code>buildDefaults</code> + the file /
                    files renderer fallbacks.
                  </li>
                  <li>
                    <code>apps/admin/components/forms/fields/{`{files,images,videos}`}-field.tsx</code>{' '}
                    — replace the <code>(value ?? []).map()</code>{' '}
                    with the <code>Array.isArray</code> guard.
                  </li>
                  <li>
                    <code>apps/admin/lib/resource.ts</code> — add{' '}
                    <code>&quot;file&quot;</code> and{' '}
                    <code>&quot;files&quot;</code> to{' '}
                    <code>ColumnFormat</code>.
                  </li>
                  <li>
                    <code>apps/admin/components/tables/import-modal.tsx</code>{' '}
                    — simplify the <code>importCfg</code> check.
                  </li>
                  <li>
                    For models with file columns, add{' '}
                    <code>{`import type { FileRef } from "../schemas/file-ref";`}</code>{' '}
                    at the top of <code>packages/shared/types/&lt;model&gt;.ts</code>.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.31.36 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.36
                </span>
                <span className="text-sm text-muted-foreground">June 24, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Bug fix: FileRef inserts now succeed on
                  Postgres.</strong> Single-file (<code>:file:</code>)
                  and multi-file (<code>:files:</code>) columns failed
                  to insert on Postgres with{' '}
                  <code>ERROR: invalid input syntax for type json
                  (SQLSTATE 22P02)</code>. SQLite and MySQL projects
                  weren&apos;t affected. This release fixes the framework
                  scaffold; existing projects get a one-file patch.
                </p>

                <h3>What was broken</h3>
                <p>
                  <code>FileRef.Value()</code> and{' '}
                  <code>FileRefs.Value()</code> returned the{' '}
                  <code>[]byte</code> from <code>json.Marshal()</code>{' '}
                  directly:
                </p>
                <pre><code>{`func (f FileRef) Value() (driver.Value, error) {
  return json.Marshal(f)   // returns []byte
}`}</code></pre>
                <p>
                  Go&apos;s <code>database/sql</code> accepts{' '}
                  <code>[]byte</code> as a valid driver.Value — and
                  lib/pq (the standard Postgres driver) encodes{' '}
                  <code>[]byte</code> as <code>bytea</code>, Postgres&apos;
                  binary type. Postgres then tries to insert the bytea
                  blob into a <code>json</code> column, fails to parse
                  the framing, and rejects with SQLSTATE 22P02.
                </p>

                <h3>The fix</h3>
                <p>
                  Both <code>Value()</code> implementations now convert
                  the JSON bytes to a Go string before returning:
                </p>
                <pre><code>{`func (f FileRef) Value() (driver.Value, error) {
  b, err := json.Marshal(f)
  if err != nil {
    return nil, err
  }
  return string(b), nil   // text, not bytea
}`}</code></pre>
                <p>
                  lib/pq sends string values as plain text, which
                  Postgres parses as JSON cleanly. SQLite and MySQL
                  are tolerant of both shapes; only Postgres was
                  strict.
                </p>

                <h3>Regression guards</h3>
                <p>
                  Two new tests in the scaffolded{' '}
                  <code>file_ref_test.go</code> assert that{' '}
                  <code>Value()</code> returns a <code>string</code>{' '}
                  type — so a future contributor can&apos;t silently
                  revert to <code>[]byte</code> without CI catching it.
                  The tests fail with a clear message pointing at the
                  Postgres bytea-vs-json issue.
                </p>

                <h3>Migration</h3>
                <p>
                  Replace <code>apps/api/internal/files/file_ref.go</code>{' '}
                  with the regenerated copy:
                </p>
                <pre><code>{`grit upgrade --files`}</code></pre>
                <p>
                  Or hand-patch both Value() methods to wrap the
                  json.Marshal result in <code>string(b)</code>{' '}
                  before returning. Postgres-on-prod users running
                  existing projects should ship this immediately.
                  SQLite-dev or MySQL projects have no urgency.
                </p>
              </div>
            </div>

            {/* v3.31.35 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.35
                </span>
                <span className="text-sm text-muted-foreground">June 24, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Excel import + export, fully client-side via
                  SheetJS.</strong> Continues the data ops arc. Every
                  resource list page now ships with a three-format
                  download menu (CSV / Excel / JSON) and a drag-and-drop
                  Excel import that previews, validates, and submits
                  rows without a single new API route.
                </p>

                <h3>Why client-side</h3>
                <p>
                  The original v3.31.35 plan put export and import on
                  the server: <code>excelize</code> for writing,{' '}
                  <code>asynq</code> + Resend for the
                  &gt;5000-row async cutoff, a new{' '}
                  <code>/import</code> endpoint with template +
                  validation. Doing it in the browser via SheetJS{' '}
                  (<code>xlsx ^0.18.5</code>) collapses all of that —
                  no new routes, no async wiring, no &quot;your file
                  is ready&quot; email loop, and tenant row data
                  never leaves the user&apos;s session just to build
                  a file.
                </p>
                <p>
                  Trade-off: very large datasets (~50k+ rows) are
                  gated by the browser&apos;s memory ceiling, not
                  server RAM. The export menu still streams every
                  page from the API before building the file, so the
                  output represents the whole filtered dataset —
                  not just what&apos;s on screen.
                </p>

                <h3>New lib/excel-utils.ts</h3>
                <ul>
                  <li><code>exportToFile(rows, columns, name, format)</code> — writes CSV / XLSX / JSON, auto-sizing columns up to 60 chars.</li>
                  <li><code>fetchAllPages(endpoint, params, onProgress)</code> — loops the resource API at <code>page_size=200</code> until every row is in hand.</li>
                  <li><code>downloadImportTemplate(resource, allowedFields?)</code> — blank workbook keyed by form field keys with a placeholder example row.</li>
                  <li><code>parseImportFile(file, resource, allowedFields?)</code> — coerces each cell to the right JS type via the field definition, returns per-row errors.</li>
                  <li><code>submitImport(endpoint, rows, onProgress)</code> — POSTs each valid row at concurrency 4 with live progress.</li>
                </ul>

                <h3>ExportMenu</h3>
                <p>
                  Split button in the toolbar: clicking the main
                  half exports in the default format (Excel when
                  enabled, else CSV); the chevron opens a menu with
                  the other formats. Uses the active search, sort,
                  filters, and date range so an export honours the
                  view the user is looking at.
                </p>

                <h3>ImportModal</h3>
                <p>
                  Three stages — file pick → validation preview →
                  submit with progress bar. The preview surfaces
                  per-row errors with field+reason, flags unknown
                  header columns, and disables Import when nothing
                  is valid. On submit, React Query invalidates the
                  resource list so the table reflects the new rows.
                </p>
                <p>
                  Header matching is loose: spaces, underscores,
                  hyphens, and case are normalised before lookup, so
                  the same template works whether a user&apos;s spreadsheet
                  has <code>first_name</code>, <code>First Name</code>,
                  or <code>firstname</code>.
                </p>

                <h3>Per-resource opt-out</h3>
                <pre><code>{`table: {
  // Hide a format from the export menu (default: all on).
  export: { csv: true, excel: true, json: false },
  // Or disable export entirely.
  // export: false,

  // Restrict importable fields to a subset.
  import: { fields: ['title', 'price', 'stock'] },
  // Or disable import entirely.
  // import: false,
}`}</code></pre>

                <h3>Migration</h3>
                <p>
                  Three new files in your scaffolded admin app:{' '}
                  <code>apps/admin/lib/excel-utils.ts</code>,{' '}
                  <code>apps/admin/components/tables/export-menu.tsx</code>,{' '}
                  <code>apps/admin/components/tables/import-modal.tsx</code>.
                  Three refreshed files:{' '}
                  <code>apps/admin/lib/resource.ts</code>,{' '}
                  <code>apps/admin/components/tables/table-toolbar.tsx</code>,{' '}
                  <code>apps/admin/components/resource/resource-page.tsx</code>.
                  Run <code>grit upgrade</code> to pull them in. The{' '}
                  <code>xlsx</code> dependency was already declared
                  in <code>package.json</code> from v3.31.34, so no
                  install step is needed.
                </p>

                <h3>Coming next</h3>
                <p>
                  v3.31.36: PDF export via @react-pdf/renderer.
                </p>
              </div>
            </div>

            {/* v3.31.34 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.34
                </span>
                <span className="text-sm text-muted-foreground">June 24, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Date filter end-to-end + stats now actually
                  reflect the filtered window.</strong> Begins the
                  data ops arc and fixes a latent bug where{' '}
                  &quot;This Week&quot; / &quot;This Month&quot; stat
                  cards were showing the total count instead of the
                  windowed count.
                </p>

                <h3>The latent stats bug</h3>
                <p>
                  Auto-default stat cards have been emitting endpoints
                  like{' '}
                  <code>/api/products?page_size=1&amp;created_since=7d</code>{' '}
                  for a while, but the API ignored{' '}
                  <code>created_since</code> — so the &quot;This
                  Week&quot; card returned the same total as the
                  &quot;Total&quot; card. v3.31.34 makes the backend
                  honour the param.
                </p>

                <h3>Server-side (paginate package)</h3>
                <p>
                  <code>Bind(c)</code> now parses four query params:
                </p>
                <ul>
                  <li><code>?created_from=2026-01-01</code> — inclusive lower bound</li>
                  <li><code>?created_to=2026-12-31</code> — inclusive upper (snapped to 23:59:59.999)</li>
                  <li><code>?created_since=7d</code> — relative shortcut (h / d / w / m units)</li>
                  <li><code>?date_field=published_at</code> — override the default <code>created_at</code> target column</li>
                </ul>
                <p>
                  Explicit <code>created_from</code> /{' '}
                  <code>created_to</code> win over{' '}
                  <code>created_since</code> so a stat-card link
                  doesn&apos;t clobber a user&apos;s picked range.
                  Applied as a single WHERE clause in{' '}
                  <code>List[T]</code>; both offset and cursor
                  pagination paths inherit.
                </p>

                <h3>Resource def</h3>
                <pre><code>{`table: {
  dateFilter: { enabled: true, field: 'created_at', label: 'Created' }
}`}</code></pre>
                <p>
                  Enabled by default. Set <code>enabled: false</code>{' '}
                  to hide. Override <code>field</code> for resources
                  where the meaningful date isn&apos;t{' '}
                  <code>created_at</code> (e.g. a <code>Booking</code>{' '}
                  resource filtering by <code>scheduled_for</code>).
                </p>

                <h3>DateFilter component</h3>
                <p>
                  New <code>&lt;DateFilter&gt;</code> in{' '}
                  <code>components/tables/date-filter.tsx</code>:
                </p>
                <ul>
                  <li>Four presets — Today, Last 7 days, Last 30 days, This month</li>
                  <li>Custom range with two date inputs + Apply button</li>
                  <li>Active state shows the current selection as a toolbar pill; X clears</li>
                  <li>Close-on-outside-click popover</li>
                  <li>URL-persisted via <code>?date=preset</code> + <code>?date_from</code> / <code>?date_to</code> so refresh + shared links rehydrate</li>
                </ul>

                <h3>Stats reflect the filter</h3>
                <p>
                  When the user picks a date range,{' '}
                  <code>ResourceListView</code> appends the resolved
                  query params to every stat card&apos;s endpoint. The
                  card labels stay fixed (&quot;Total&quot;, &quot;This
                  Week&quot;, etc.) but their numbers now match the
                  table below. No more &quot;Total: 10,000; list shows
                  142&quot; mismatch.
                </p>

                <h3>Migration</h3>
                <p>
                  Four files refreshed:{' '}
                  <code>apps/api/internal/paginate/paginate.go</code>,{' '}
                  <code>apps/admin/components/tables/table-toolbar.tsx</code>,{' '}
                  <code>apps/admin/components/resource/resource-page.tsx</code>,{' '}
                  <code>apps/admin/hooks/use-resource.ts</code>,
                  plus the new{' '}
                  <code>apps/admin/components/tables/date-filter.tsx</code>.
                </p>

                <h3>Coming next</h3>
                <p>
                  v3.31.35: Excel import + async cutoff for export
                  (&gt;5000 rows = asynq job + Resend email) +
                  per-resource opt-out. v3.31.36: PDF export via{' '}
                  @react-pdf/renderer.
                </p>
              </div>
            </div>

            {/* v3.31.33 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.33
                </span>
                <span className="text-sm text-muted-foreground">June 24, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>File lifecycle — immediate S3 delete on
                  replacement + daily orphan cleanup cron.</strong>{' '}
                  Closes the loop on the file-fields work from
                  v3.31.30-32: bucket no longer accumulates dead
                  objects when files get swapped or forms get
                  abandoned.
                </p>

                <h3>internal/files lifecycle helpers</h3>
                <ul>
                  <li>
                    <code>DiffSingle(old, new)</code> — returns the
                    key removed when a single-file column is replaced
                    or cleared.
                  </li>
                  <li>
                    <code>DiffMulti(old, new)</code> — returns keys
                    present in old but missing from new (gallery
                    pruning).
                  </li>
                  <li>
                    <code>CleanupRemoved(ctx, st, old, new)</code> —
                    reflection-based: walks both struct values,
                    finds <code>*FileRef</code> + <code>FileRefs</code>{' '}
                    fields, computes the diff, deletes the removed S3
                    objects. One line in the handler regardless of
                    how many file columns the resource has.
                  </li>
                  <li>
                    <code>ClaimRefs(ctx, db, record)</code> — walks
                    the same FileRef columns and stamps{' '}
                    <code>claimed_at = now()</code> on the underlying
                    Upload rows so the orphan cleanup cron knows the
                    upload is in use.
                  </li>
                  <li>
                    <code>RunOrphanCleanup(ctx, db, st, minAge)</code>{' '}
                    — finds Upload rows with{' '}
                    <code>claimed_at IS NULL</code> older than{' '}
                    <code>minAge</code> (24h), deletes them from S3
                    and the DB. Best-effort S3 delete: if it fails
                    we still drop the DB row so the same orphan
                    isn&apos;t retried forever.
                  </li>
                </ul>

                <h3>Upload.ClaimedAt column</h3>
                <p>
                  New nullable timestamp on the Upload model.
                  Auto-migration adds it; existing rows start as NULL
                  and get claimed the next time their parent record
                  is updated. The 24h grace period before orphan
                  cleanup means a fresh deploy won&apos;t purge
                  historical uploads — the cron only catches uploads
                  truly created in the past 24h that never got
                  claimed.
                </p>

                <h3>Daily cron job</h3>
                <p>
                  New <code>uploads:cleanup_orphans</code> asynq task
                  runs at 03:15 daily (low-traffic window). Registered
                  in <code>internal/cron/cron.go</code> and handled by{' '}
                  <code>handleUploadsOrphanCleanup</code> in{' '}
                  <code>internal/jobs/jobs.go</code>.
                </p>

                <h3>Generated handler injection</h3>
                <p>
                  Resources with{' '}
                  <code>:file:</code> / <code>:files:</code> fields
                  now get:
                </p>
                <ul>
                  <li>
                    <code>Storage *storage.Storage</code> field on
                    the Handler struct, wired in routes via{' '}
                    <code>Storage: svc.Storage</code>.
                  </li>
                  <li>
                    Create handler: <code>files.ClaimRefs</code> call
                    after successful save.
                  </li>
                  <li>
                    Update handler: snapshots the old record,
                    diff-deletes removed S3 objects, then claims the
                    new refs.
                  </li>
                </ul>
                <p>
                  Resources without file fields stay exactly as before
                  — no Storage field, no extra imports, no dead code.
                  The injection is conditional on the generator
                  detecting at least one file/files field.
                </p>

                <h3>Migration</h3>
                <p>
                  Existing projects: the Upload model needs the
                  <code> ClaimedAt</code> column. GORM auto-migration
                  in <code>cmd/server/main.go</code> handles it on
                  next boot. Re-run{' '}
                  <code>grit generate resource &lt;Name&gt;</code> on
                  any resource with file fields to pick up the
                  cleanup-aware handler template.
                </p>

                <h3>Coming next</h3>
                <p>
                  v3.31.34 begins the data ops arc — server-side
                  Excel export via excelize, bulk Excel import with
                  template generation + row-by-row validation,
                  React-PDF rendering, per-page date filter, and
                  per-resource opt-out for export / import.
                </p>
              </div>
            </div>

            {/* v3.31.32 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.32
                </span>
                <span className="text-sm text-muted-foreground">June 24, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Storage admin page surfaces FileRef
                  totals.</strong> The original Files page was a flat
                  uploads grid — useful for browsing but offered no
                  sense of how much storage you were actually using,
                  or what was eating it.
                </p>

                <h3>New API endpoint</h3>
                <p>
                  <code>GET /api/uploads/stats</code> returns:
                </p>
                <ul>
                  <li>
                    <code>total_count</code> — how many uploads
                  </li>
                  <li>
                    <code>total_size</code> — sum of bytes
                  </li>
                  <li>
                    <code>by_kind</code> — count + size grouped by
                    MIME bucket (image / video / audio / pdf /
                    document / spreadsheet / other). Single SQL
                    GROUP BY with portable CASE expression — works on
                    Postgres and SQLite without engine-specific JSON
                    functions.
                  </li>
                </ul>

                <h3>Storage stats panel</h3>
                <p>
                  The Files admin page now shows three big numbers up
                  top (Total files / Total storage / Avg file size),
                  then a per-kind breakdown with proportional progress
                  bars sorted by largest consumer. Image-heavy
                  projects can see at a glance whether to migrate to a
                  CDN; CSV-heavy projects can spot a runaway export
                  pipeline.
                </p>

                <h3>Dropzone variant standardisation</h3>
                <p>
                  Default + Compact variants now route their uploading
                  state through the unified{' '}
                  <code>&lt;UploadProgress&gt;</code> component so the
                  per-field <code>progress</code> prop (bar / circular
                  / pulse) actually takes effect on both. Minimal,
                  Avatar, and Inline variants are space-constrained
                  by design and keep their bespoke single-spinner
                  treatment.
                </p>

                <h3>Migration</h3>
                <p>
                  Three files refreshed:{' '}
                  <code>apps/api/internal/handlers/upload.go</code>{' '}
                  (Stats handler),{' '}
                  <code>apps/api/internal/routes.go</code> (new route),{' '}
                  <code>apps/admin/app/(dashboard)/system/files/page.tsx</code>{' '}
                  (stats panel) and{' '}
                  <code>apps/admin/hooks/use-system.ts</code>{' '}
                  (useUploadStats hook). Re-run{' '}
                  <code>grit generate resource</code> for any resource
                  to pull the updates.
                </p>

                <h3>Coming next</h3>
                <p>
                  v3.31.33 ships the file lifecycle work: immediate
                  S3 delete when a record swaps its file, plus a
                  daily orphan-cleanup cron that purges Upload rows
                  whose key is referenced nowhere. v3.31.34 begins
                  the data ops arc — date filter, Excel
                  import/export, PDF render via @react-pdf/renderer.
                </p>
              </div>
            </div>

            {/* v3.31.31 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.31
                </span>
                <span className="text-sm text-muted-foreground">June 24, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>File fields polish — progress variants,
                  type-aware previews, reorder.</strong>
                </p>

                <h3>Bug fix from v3.31.30</h3>
                <p>
                  The Dropzone was still reading{' '}
                  <code>data.original_name</code> and{' '}
                  <code>data.mime_type</code> from the upload response,
                  but v3.31.30 changed{' '}
                  <code>POST /api/uploads</code> to return a FileRef
                  shape with <code>data.name</code> and{' '}
                  <code>data.mime</code>. Files uploaded after v3.31.30
                  appeared with the generic File client-side fallback
                  name instead of the real filename from the server.
                  Now reads both shapes (FileRef first, legacy
                  fallback) so cross-version compatibility holds.
                </p>
                <p>
                  UploadedFile also carries the explicit <code>key</code>{' '}
                  field now, so the FileField bridge round-trips the
                  S3 key losslessly instead of recomputing it from
                  the URL pathname.
                </p>

                <h3>Three progress variants</h3>
                <ul>
                  <li>
                    <strong>bar</strong> (default) — linear progress
                    bar with spinner + percentage label.
                  </li>
                  <li>
                    <strong>circular</strong> — donut with the % inside.
                    SVG, no extra dependency.
                  </li>
                  <li>
                    <strong>pulse</strong> — three pulsing dots + %.
                    Minimal chrome for compact contexts.
                  </li>
                </ul>
                <p>
                  Pick a variant per field:
                </p>
                <pre><code>{`{ key: "avatar", type: "file", accepts: ["image"], progress: "circular" }`}</code></pre>
                <p>
                  The Default dropzone variant routes its uploading
                  state through the new <code>&lt;UploadProgress&gt;</code>{' '}
                  component. The other four dropzone variants (compact /
                  minimal / avatar / inline) keep their bespoke inline
                  progress UI for now — v3.31.32 standardises them.
                </p>

                <h3>Type-aware FilePreview</h3>
                <p>
                  Single image preview stays as a thumbnail. Video gets
                  a play badge over a dark thumb. Audio shows a music
                  icon. PDF / Word / Excel / CSV render format-specific
                  glyphs with colour-coded tints (PDF red, Word blue,
                  Excel green). Everything else falls back to the
                  generic File icon.
                </p>

                <h3>Reorder by up/down arrows</h3>
                <p>
                  Multi-file (<code>:files:</code>) preview rows now
                  show small up/down arrow buttons when{' '}
                  <code>reorderable</code> is true (default). Adjacent
                  swap; first row&apos;s up button is disabled, last
                  row&apos;s down button is disabled. No new
                  dependencies — drag-reorder via dnd-kit is a future
                  polish.
                </p>

                <h3>Resource def knobs</h3>
                <p>
                  Three new optional props on file/files{' '}
                  <code>FieldDefinition</code>:
                </p>
                <ul>
                  <li>
                    <code>dropzone</code>: <code>&quot;default&quot;</code>{' '}
                    | <code>&quot;compact&quot;</code> |{' '}
                    <code>&quot;minimal&quot;</code> |{' '}
                    <code>&quot;avatar&quot;</code> |{' '}
                    <code>&quot;inline&quot;</code>
                  </li>
                  <li>
                    <code>progress</code>:{' '}
                    <code>&quot;bar&quot;</code> |{' '}
                    <code>&quot;circular&quot;</code> |{' '}
                    <code>&quot;pulse&quot;</code>
                  </li>
                  <li>
                    <code>reorderable</code>: <code>boolean</code>{' '}
                    (default true; multi-file only)
                  </li>
                </ul>
                <p>
                  These are pure overrides — the CLI doesn&apos;t emit
                  them automatically; hand-edit the resource def to
                  customise.
                </p>

                <h3>Migration</h3>
                <p>
                  Existing scaffolded projects need three files
                  refreshed:{' '}
                  <code>components/ui/dropzone.tsx</code>,{' '}
                  <code>components/forms/fields/file-field.tsx</code>,
                  and{' '}
                  <code>components/forms/fields/files-field.tsx</code>.
                  Plus add{' '}
                  <code>FileSpreadsheet</code> and{' '}
                  <code>Music</code> to the export block in{' '}
                  <code>lib/icons.ts</code>. Re-run{' '}
                  <code>grit generate resource</code> for any resource
                  to pull the updates — the files live once, not per
                  resource.
                </p>
              </div>
            </div>

            {/* v3.31.30 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.30
                </span>
                <span className="text-sm text-muted-foreground">June 24, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>File fields — first-class file + files
                  types in <code>grit generate resource</code>.</strong>{' '}
                  Replaces the awkward old pattern of treating uploads
                  as <code>string</code> URLs.
                </p>

                <h3>New CLI syntax</h3>
                <p>
                  Single file:{' '}
                  <code>grit generate resource Product --fields &quot;image:file:image&quot;</code>{' '}
                  scaffolds a single-image field that accepts
                  jpg / png / gif / webp / avif / svg.
                </p>
                <p>
                  Multiple files:{' '}
                  <code>gallery:files:image</code> for a multi-image
                  gallery.
                </p>
                <p>
                  Bracketed accept-list for mixed types:{' '}
                  <code>attachment:file:[pdf,doc,image,video,zip]</code>.
                  Bare commas don&apos;t work because the top-level
                  field separator is also <code>,</code>; the parser
                  is bracket-aware so the inner list stays glued
                  together.
                </p>
                <p>
                  Accept aliases: <code>image</code>,{' '}
                  <code>video</code>, <code>audio</code>,{' '}
                  <code>pdf</code>, <code>doc</code>,{' '}
                  <code>excel</code>, <code>csv</code>,{' '}
                  <code>zip</code>, <code>archive</code>,{' '}
                  <code>all</code>.
                </p>

                <h3>What gets generated</h3>
                <ul>
                  <li>
                    <strong>Go model:</strong> field typed as{' '}
                    <code>*files.FileRef</code> (single) or{' '}
                    <code>files.FileRefs</code> (multi), stored as JSON
                    via GORM Value/Scan adapters in the new{' '}
                    <code>internal/files</code> package.
                  </li>
                  <li>
                    <strong>Zod schema:</strong> imports{' '}
                    <code>FileRefSchema</code> from the shared package
                    — a single source of truth for the JSON shape.
                  </li>
                  <li>
                    <strong>Admin resource def:</strong> auto-emits{' '}
                    <code>accepts</code> and <code>maxSizeMB</code> so
                    the form&apos;s upload endpoint enforces the
                    per-field validation.
                  </li>
                  <li>
                    <strong>FormBuilder:</strong> dispatches{' '}
                    <code>file</code> / <code>files</code> types to
                    the FileRef-aware FileField / FilesField
                    components.
                  </li>
                  <li>
                    <strong>DataTable:</strong> file columns render as
                    thumbnails for images, MIME-typed icons for
                    everything else. Multi-file columns stack the
                    first three thumbnails with a +N overflow chip.
                  </li>
                </ul>

                <h3>API changes</h3>
                <p>
                  <code>POST /api/uploads</code> now accepts{' '}
                  <code>?accepts=&lt;aliases&gt;&amp;max_size=&lt;bytes&gt;</code>{' '}
                  query params so the server validates against the
                  per-field accept set (not just a global allowlist).
                  Response shape changed to return a{' '}
                  <code>FileRef</code> directly under{' '}
                  <code>data</code> — drop-in for form state.
                </p>

                <h3>Defaults</h3>
                <ul>
                  <li>Single file max: 5MB (300MB for video).</li>
                  <li>Multi-file count: 5.</li>
                  <li>
                    Dropzone variant: the existing default boxed-dashed
                    style. v3.31.31 adds 4 more variants (minimal,
                    card, avatar, inline) + 3 progress variants + dnd-kit
                    reorder.
                  </li>
                </ul>

                <h3>Migration</h3>
                <p>
                  Existing scaffolded projects need three things to
                  pick up file fields:{' '}
                  <code>apps/api/internal/files/</code> (new package),
                  the updated <code>handlers/upload.go</code>, and the
                  refactored{' '}
                  <code>components/forms/fields/file-field.tsx</code> +{' '}
                  <code>files-field.tsx</code> + the new{' '}
                  <code>lib/file-accepts.ts</code>. Re-run{' '}
                  <code>grit generate resource</code> for any resource
                  to get the updated templates — the new code lives
                  once per project, not per resource.
                </p>
              </div>
            </div>

            {/* v3.31.29 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.29
                </span>
                <span className="text-sm text-muted-foreground">June 24, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Stats cards now refetch after create / update /
                  delete.</strong> Total / This Week / This Month no
                  longer go stale until manual reload.
                </p>

                <h3>The bug</h3>
                <p>
                  Resource mutations all called{' '}
                  <code>invalidateQueries(&#123; queryKey: [endpoint] &#125;)</code>{' '}
                  on success, expecting React Query&apos;s
                  prefix-matching to invalidate every query under that
                  resource. But the stat-card query in{' '}
                  <code>PageHeader</code> was keyed with{' '}
                  <code>[&quot;stat&quot;, endpoint, field]</code> —
                  starting with the literal string{' '}
                  <code>&quot;stat&quot;</code>, not the endpoint. The
                  invalidation never matched it, and a{' '}
                  <code>staleTime: 30_000</code> meant the value didn&apos;t
                  even auto-refetch for 30 seconds.
                </p>
                <p>
                  Stats also use endpoints with query-string suffixes
                  (e.g.{' '}
                  <code>/api/products?page_size=1&amp;created_since=7d</code>),
                  so even if the key had started with the endpoint
                  string, it wouldn&apos;t have matched the bare
                  <code> /api/products</code> the mutation invalidates.
                </p>

                <h3>The fix</h3>
                <p>
                  Stat queryKey now starts with the <em>base</em> endpoint
                  (no query string):{' '}
                  <code>[endpoint.split(&quot;?&quot;)[0], &quot;stat&quot;, endpoint, field]</code>.
                  Mutation invalidation prefix-matches it, and the staleTime
                  is gone so the cards refetch immediately on success.
                </p>

                <h3>Migration</h3>
                <p>
                  Existing projects: copy the new <code>StatCardItem</code>{' '}
                  hook body from{' '}
                  <code>components/layout/page-header.tsx</code>.
                  One-function change.
                </p>
              </div>
            </div>

            {/* v3.31.28 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.28
                </span>
                <span className="text-sm text-muted-foreground">June 24, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Colored toasters.</strong> Success toasts
                  are now green, errors red, warnings amber, info
                  blue — instead of the previous neutral grey-on-grey
                  that made every toast look identical.
                </p>

                <h3>What changed</h3>
                <p>
                  The scaffolded <code>&lt;Toaster&gt;</code> in{' '}
                  <code>components/shared/providers.tsx</code> now
                  passes <code>richColors</code>, and{' '}
                  <code>app/globals.css</code> bridges Grit&apos;s
                  theme tokens (<code>--success</code>,{' '}
                  <code>--danger</code>, <code>--warning</code>,{' '}
                  <code>--info</code>) into sonner&apos;s palette
                  slots via <code>color-mix()</code>. Each theme
                  (atlas / aurora / pulse / midnight) already
                  redefines those four tokens, so toasters
                  automatically pick up the active brand colors
                  — no per-theme overrides needed.
                </p>

                <h3>Migration</h3>
                <p>
                  Existing projects: copy the new{' '}
                  <code>Toaster</code> mount from{' '}
                  <code>providers.tsx</code> and the{' '}
                  <code>[data-sonner-toaster]</code> CSS block from{' '}
                  <code>globals.css</code> (right after the
                  scrollbar rules). All toast call sites in the
                  scaffold already use{' '}
                  <code>toast.success()</code> /{' '}
                  <code>toast.error()</code> etc., so they pick up
                  the new colors with zero code changes.
                </p>
              </div>
            </div>

            {/* v3.31.27 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.27
                </span>
                <span className="text-sm text-muted-foreground">June 24, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Fix React Rules of Hooks violation when{' '}
                  <code>formView: &apos;page&apos;</code> resources
                  switch between list and form views.</strong>{' '}
                  Reported by a learner who scaffolded a Category
                  resource with <code>formView: &apos;page&apos;</code>{' '}
                  and clicked &quot;New Category&quot;.
                </p>

                <h3>The bug</h3>
                <p>
                  <code>ResourcePage</code> declared a few hooks at
                  the top (<code>useRouter</code>,{' '}
                  <code>useSearchParams</code>), then performed{' '}
                  <em>early returns</em> for the form-page case
                  (<code>action=create</code> or <code>action=edit</code>),
                  then declared ~20 more hooks below
                  (<code>useState</code>, <code>useResource</code>,{' '}
                  <code>useMemo</code>, <code>useCallback</code> x many).
                  When the URL changed and the component switched
                  between list mode and form mode, the hook count
                  changed between renders — React 19 throws &quot;Rendered
                  fewer hooks than expected.&quot;
                </p>

                <h3>The fix</h3>
                <p>
                  Split <code>ResourcePage</code> into a thin router
                  shell + a separate <code>ResourceListView</code>{' '}
                  component. The router only calls{' '}
                  <code>useSearchParams</code> and the routing
                  helpers, then either renders one of the form
                  variants or delegates to <code>ResourceListView</code>.
                  The list view owns all 20+ list-mode hooks. Each
                  function now has a stable hook count across
                  renders, and the form path never mounts the
                  list-mode hooks (so it doesn&apos;t spawn an
                  unnecessary <code>useResource</code> fetch either).
                </p>

                <h3>Migration</h3>
                <p>
                  Existing projects using <code>formView: &apos;page&apos;</code>{' '}
                  or <code>&apos;page-steps&apos;</code> need to update{' '}
                  <code>apps/admin/components/resource/resource-page.tsx</code>.
                  Re-run <code>grit generate resource &lt;Name&gt;</code> on
                  any resource (the file lives once, not per-resource)
                  or copy the new structure from the scaffold output.
                </p>
              </div>
            </div>

            {/* v3.31.14 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.14
                </span>
                <span className="text-sm text-muted-foreground">June 21, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>CLI prompt cleanup, Sentinel/Pulse links go to the API, and the Security + Performance dashboards finally show real data.</strong>
                </p>

                <h3>CLI: one form instead of three selects</h3>
                <p>
                  <code>grit new</code>&apos;s architecture / frontend /
                  theme prompts were running as three back-to-back{' '}
                  <code>huh.NewSelect</code> calls. On Git Bash (MINGW64)
                  the lack of full ANSI cursor-up support meant each
                  re-render stacked into scrollback, producing the
                  &quot;same prompt printed twice&quot; effect. Combined
                  into a single <code>huh.NewForm</code> with conditional
                  <code> WithHideFunc</code> groups — one tidy block of
                  output, atomic submit.
                </p>

                <h3>Sentinel / Pulse links point at the API origin</h3>
                <p>
                  <code>/system/security</code> &quot;Open Sentinel&quot;
                  used a Next.js <code>&lt;Link href=&quot;/sentinel/ui&quot;&gt;</code>,
                  which resolves relative to the admin host (<code>:3001</code>).
                  Both Sentinel and Pulse are mounted on the Go API
                  (<code>:8080</code>), so the links 404&apos;d. Replaced
                  with a plain <code>&lt;a&gt;</code> using
                  <code> NEXT_PUBLIC_API_URL</code>.
                </p>

                <h3>Security + Performance dashboards return real data</h3>
                <p>
                  Both pages were calling endpoints that either didn&apos;t
                  exist (<code>/api/admin/performance/summary</code>) or
                  returned a wrapped <code>{`{data: {...}}`}</code> envelope
                  with raw Sentinel/Pulse internals under unfamiliar keys
                  (<code>{`{summary, score, threats, ...}`}</code>). The
                  React queries unwrapped axios&apos; <code>.data</code>{' '}
                  and looked for <code>data.banned_ips_now</code> /
                  <code> data.latency.p50</code> which didn&apos;t exist
                  in the response — every KPI rendered as 0 or em-dash.
                </p>
                <ul>
                  <li>
                    <code>handlers/observability.go</code> rewritten: hits
                    Pulse&apos;s <code>/overview</code>, <code>/runtime/current</code>,
                    <code> /database/n1/ranked</code>, and <code>/errors</code>{' '}
                    endpoints, then reshapes the responses into a flat
                    {' '}<code>{`{latency, traffic, errors, saturation, slowest_routes, n1_detections, recent_errors}`}</code>{' '}
                    envelope. No more <code>{`{data: ...}`}</code> wrapper.
                  </li>
                  <li>
                    <code>handlers/security.go</code> rewritten: hits
                    Sentinel&apos;s <code>/ip/blocked</code>,
                    <code> /analytics/summary?window=24h</code>, and
                    <code> /threats</code> (the prior <code>/dashboard/summary</code>{' '}
                    endpoint doesn&apos;t exist in this Sentinel version);
                    returns the flat <code>{`{banned_ips_now, auto_bans_24h, active_bans, recent_threats, ...}`}</code> shape
                    the page expects.
                  </li>
                  <li>
                    Performance page corrected to hit{' '}
                    <code>/api/admin/observability/summary</code> instead
                    of the nonexistent <code>/performance/summary</code>.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.31.26 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.26
                </span>
                <span className="text-sm text-muted-foreground">June 21, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Fix two bugs in the FormShare dispatcher
                  template</strong> reported by a learner who created
                  a fresh project and ran{' '}
                  <code>grit generate resource Category</code> +{' '}
                  <code>Product</code>. The API failed to build with{' '}
                  <code>syntax error: non-declaration statement
                  outside function body</code>.
                </p>

                <h3>Bug 1 — marker collision with doc comment</h3>
                <p>
                  The scaffolded{' '}
                  <code>services/form_share_dispatch.go</code> doc
                  comment literally contained the string{' '}
                  <code>// grit:form-share:dispatch marker</code>.
                  When the generator ran{' '}
                  <code>injectBefore</code> for a new resource, it
                  found that occurrence first (above the function,
                  outside any function body) and inserted every case
                  there. The function&apos;s switch stayed empty,
                  and the cases sat in package scope where they
                  produced a syntax error.
                </p>
                <p>
                  Fix: rephrased the doc comment to describe the
                  marker without containing the marker string.
                </p>

                <h3>Bug 2 — function param named "body", inject uses "fields"</h3>
                <p>
                  The dispatcher&apos;s third parameter was{' '}
                  <code>body map[string]interface{`{}`}</code>, but
                  every injected case uses{' '}
                  <code>json.Marshal(fields)</code>. Even if Bug 1
                  hadn&apos;t hit first, the cases would have failed
                  to compile with{' '}
                  <code>undefined: fields</code>.
                </p>
                <p>
                  Fix: renamed the parameter to <code>fields</code>{' '}
                  so it matches what the inject template produces.
                </p>

                <h3>Migration</h3>
                <p>
                  Existing projects that ran{' '}
                  <code>grit generate resource X</code> on or after
                  v3.31.20 may have a broken{' '}
                  <code>form_share_dispatch.go</code>. To fix:
                </p>
                <ol>
                  <li>
                    Open{' '}
                    <code>apps/api/internal/services/form_share_dispatch.go</code>.
                  </li>
                  <li>
                    Move any <code>case "X":</code> blocks that
                    landed above the function back inside the{' '}
                    <code>switch</code> below.
                  </li>
                  <li>
                    Rename the function parameter from <code>body</code>{' '}
                    to <code>fields</code> if needed.
                  </li>
                  <li>
                    Rephrase the doc comment so it doesn&apos;t
                    contain the literal marker string.
                  </li>
                </ol>
              </div>
            </div>

            {/* v3.31.25 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.25
                </span>
                <span className="text-sm text-muted-foreground">June 21, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Audit trail for public form submissions.</strong>{' '}
                  The last deferred item from PLAN_FORMS_AND_SHARING.md.
                  Operators can now see every submission that came in
                  through each share — with timestamp, IP, and
                  User-Agent.
                </p>

                <h3>What landed</h3>
                <ul>
                  <li>
                    <strong>New <code>FormSubmission</code> model</strong>{' '}
                    — one row per successful public submission. Captures
                    share_id, resource_name, record_id, IP, User-Agent,
                    timestamp. Soft-deletable for retention.
                  </li>
                  <li>
                    <strong><code>PublicSubmit</code> writes the audit
                    row</strong> after a successful dispatch. Best-effort:
                    failure to write the audit row does NOT roll back
                    the user&apos;s submission. They still get their
                    record; the admin just misses one line in the
                    trail.
                  </li>
                  <li>
                    <strong>New admin endpoint:</strong>{' '}
                    <code>GET /api/admin/form-submissions?share_id=&amp;resource_name=</code>{' '}
                    — paginated audit log, filterable by share or
                    resource.
                  </li>
                  <li>
                    <strong>Admin UI:</strong> the /system/form-shares
                    page gains an <em>Audit</em> button per share.
                    Click → modal listing the 100 most recent
                    submissions for that share with timestamp, record
                    ID, IP, and a truncated UA tooltip.
                  </li>
                </ul>

                <h3>Why a separate table, not a column</h3>
                <p>
                  An earlier draft considered adding{' '}
                  <code>source_share_id</code> as a column on every
                  scaffolded model. That approach is invasive — every
                  existing project would need a migration to add the
                  column to Contact / Application / Lead / etc. The
                  audit-table approach is purely additive: new project
                  or existing,{' '}
                  <code>grit migrate</code> creates the new{' '}
                  <code>form_submissions</code> table and existing
                  models stay untouched.
                </p>
                <p>
                  Bonus: the audit table captures richer data than a
                  column could (IP + User-Agent), which is useful for
                  spam triage and compliance.
                </p>

                <h3>Phase recap, complete</h3>
                <p>
                  Every numbered item on PLAN_FORMS_AND_SHARING.md
                  has shipped:
                </p>
                <ul>
                  <li>v3.31.16 — sync auto-add admin fields</li>
                  <li>v3.31.17 — formView sheet / modal / page</li>
                  <li>v3.31.18 — form groups + per-group PATCH</li>
                  <li>v3.31.19 — column-pack auto-detection</li>
                  <li>v3.31.20 — public form sharing</li>
                  <li>v3.31.21 — grit expose form / table</li>
                  <li>v3.31.22 — grit add web-auth</li>
                  <li>v3.31.23 — course lessons + tests</li>
                  <li>v3.31.24 — --public-share + --token flags</li>
                  <li>v3.31.25 — audit trail (this release)</li>
                </ul>
              </div>
            </div>

            {/* v3.31.24 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.24
                </span>
                <span className="text-sm text-muted-foreground">June 21, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>
                    grit expose form gains <code>--public-share</code> +{' '}
                    <code>--token</code> flags
                  </strong>{' '}
                  — the deferred public-form variant from the v3.31.21
                  changelog now ships. Scaffold a public-facing form at
                  any URL of your choosing that posts to a FormShare
                  endpoint instead of the authenticated hook.
                </p>

                <h3>Usage</h3>
                <pre className="overflow-x-auto rounded-lg bg-bg-elevated p-3 text-xs"><code>{`# Hard-code the token into the page
grit expose form Contact \\
  --to apps/web/app/contact-us/page.tsx \\
  --public-share \\
  --token 9CkLh7gJZQrPeNwMo3F8x_iVjA8U2nXt

# Or omit --token and let the page read NEXT_PUBLIC_FORM_TOKEN at runtime
grit expose form Contact \\
  --to apps/web/app/contact-us/page.tsx \\
  --public-share`}</code></pre>

                <h3>What the generated page does</h3>
                <ul>
                  <li>
                    Posts to{' '}
                    <code>/api/public/forms/&lt;token&gt;/submit</code> —
                    no auth required, no useCreate hook imported.
                  </li>
                  <li>
                    Probes <code>/api/public/forms/&lt;token&gt;</code>{' '}
                    on mount to confirm the share is enabled and to
                    learn whether to render a password gate.
                  </li>
                  <li>
                    Shows an amber "Form unavailable" card when the
                    token is missing, disabled, or invalid — instead
                    of a blank form.
                  </li>
                  <li>
                    Token resolution: literal from <code>--token</code>{' '}
                    when set; otherwise{' '}
                    <code>process.env.NEXT_PUBLIC_FORM_TOKEN</code>{' '}
                    at module load. Pick whichever fits your env model.
                  </li>
                </ul>

                <h3>When to use this</h3>
                <p>
                  Use <code>--public-share</code> when you want a
                  branded public form at your own URL
                  (<code>/contact-us</code>, <code>/apply</code>,{' '}
                  <code>/leads</code>) instead of the default{' '}
                  <code>/forms/[token]</code> page. The dispatcher,
                  rate limits, and password gate behave identically;
                  only the URL and styling are yours to control.
                </p>

                <h3>Lesson update</h3>
                <p>
                  The grit-expose lesson now has a "<code>--public-share</code>:
                  a public form on YOUR url" section with both
                  embed-token and env-token examples.
                </p>
              </div>
            </div>

            {/* v3.31.23 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.23
                </span>
                <span className="text-sm text-muted-foreground">June 21, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Docs + tests follow-up to the
                  PLAN_FORMS_AND_SHARING.md arc.</strong> Phases 2-4
                  shipped without dedicated course lessons; this
                  release closes that gap.
                </p>

                <h3>Three new course lessons</h3>
                <p>
                  Chapter 4 (&quot;Code Generation &amp; Type Sync&quot;)
                  picks up a new module —{' '}
                  <strong>Going public</strong> — with three lessons
                  covering the post-resource lifecycle:
                </p>
                <ul>
                  <li>
                    <strong>grit expose form / table</strong> — when
                    to use each, anatomy of the commands, field
                    filtering, combining with form sharing.
                  </li>
                  <li>
                    <strong>Public form sharing</strong> — the dispatch
                    pattern, password gating, sharing the link,
                    disabling and regenerating, what it can&apos;t
                    do (yet).
                  </li>
                  <li>
                    <strong>Protecting web pages</strong> — middleware
                    vs ProtectedWebRoute, when each one fits, how
                    they layer.
                  </li>
                </ul>

                <h3>Unit tests for the expose package</h3>
                <p>
                  <code>internal/expose</code> now has 10 unit tests
                  covering the security-critical field filter
                  (autoFields drops framework columns, pointer + value
                  associations, slice associations; keeps all 7
                  primitive types), label generation (acronym
                  handling), and the pluralisation helpers
                  (pluralPascal, pluralKebab). Plus path validation
                  for resolveTarget. All passing.
                </p>

                <h3>Course chapter 4 now has 13 lessons</h3>
                <p>
                  Up from 10 in the previous release. The chapter
                  covers the full resource lifecycle from initial
                  generation through customisation, sharing,
                  exposure, and protection — end to end.
                </p>
              </div>
            </div>

            {/* v3.31.22 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.22
                </span>
                <span className="text-sm text-muted-foreground">June 21, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>grit add web-auth</strong> — Phase 4 of
                  PLAN_FORMS_AND_SHARING.md, the final phase. With this
                  release, every phase of the forms / sharing
                  initiative has shipped (v3.31.16 → v3.31.22).
                </p>

                <p>
                  The web app already shipped with login / register /
                  forgot-password pages and a <code>useMe()</code>{' '}
                  hook. What was missing: a way to mark <em>which</em>{' '}
                  customer-facing pages require sign-in.{' '}
                  <code>grit add web-auth</code> closes that gap with
                  two complementary patterns.
                </p>

                <h3>Files scaffolded</h3>
                <ul>
                  <li>
                    <strong><code>apps/web/middleware.ts</code></strong>{' '}
                    — SSR cookie redirect. Runs on every Next.js
                    request, checks for the{' '}
                    <code>grit_access</code> HttpOnly cookie, redirects
                    to <code>/login?next=…</code> when missing on a
                    protected path. Also bounces already-signed-in
                    visitors off the login/register pages so they
                    don&apos;t see a form they don&apos;t need. Edit
                    the <code>PROTECTED_PATHS</code> and{' '}
                    <code>AUTH_PATHS</code> arrays to customise.
                  </li>
                  <li>
                    <strong>
                      <code>apps/web/components/ProtectedWebRoute.tsx</code>
                    </strong>{' '}
                    — client-side wrapper. Wraps a page with{' '}
                    <code>&lt;ProtectedWebRoute&gt;{`{children}`}
                    &lt;/ProtectedWebRoute&gt;</code> to enforce auth
                    in cases where middleware can&apos;t help — e.g.
                    role-gated content (the cookie doesn&apos;t carry
                    the role; <code>useMe()</code> returns the full
                    user). Supports an optional{' '}
                    <code>roles</code> prop.
                  </li>
                </ul>

                <h3>The two patterns</h3>
                <ul>
                  <li>
                    <strong>Middleware (SSR)</strong> — fast, no
                    network round-trip per request, no flash of
                    unauthorized content. Use for "is the visitor
                    signed in?" pages: account dashboards, checkout,
                    member-only content.
                  </li>
                  <li>
                    <strong>ProtectedWebRoute (client)</strong> —
                    makes a real <code>/api/auth/me</code> probe.
                    Catches expired-but-present cookies and supports
                    role checks. Use it when middleware isn&apos;t
                    enough.
                  </li>
                </ul>

                <h3>Behavior</h3>
                <p>
                  Both files are idempotent — re-running{' '}
                  <code>grit add web-auth</code> without{' '}
                  <code>--force</code> skips existing files. The
                  scaffold prints a clear notice so operators know
                  what was created and what was preserved.
                </p>

                <h3>Phase recap (v3.31.16 → v3.31.22)</h3>
                <ul>
                  <li><strong>v3.31.16</strong> — sync auto-adds new model fields to admin resource files</li>
                  <li><strong>v3.31.17</strong> — formView sheet / modal / page</li>
                  <li><strong>v3.31.18</strong> — form groups + per-group PATCH save</li>
                  <li><strong>v3.31.19</strong> — column-pack auto-detection (name + email)</li>
                  <li><strong>v3.31.20</strong> — public form sharing (token + bcrypt password)</li>
                  <li><strong>v3.31.21</strong> — grit expose form / grit expose table</li>
                  <li><strong>v3.31.22</strong> — grit add web-auth (this release)</li>
                </ul>
              </div>
            </div>

            {/* v3.31.21 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.21
                </span>
                <span className="text-sm text-muted-foreground">June 21, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>grit expose form / grit expose table</strong> — Phase 3 of PLAN_FORMS_AND_SHARING.md.
                </p>

                <p>
                  Two new CLI commands that scaffold a Next.js page in{' '}
                  <code>apps/web/</code> for an existing resource. The
                  page consumes the auto-generated React Query hook
                  directly, so you get list/create flows on a
                  customer-facing site without re-implementing
                  anything.
                </p>

                <h3>Commands</h3>
                <pre className="overflow-x-auto rounded-lg bg-bg-elevated p-3 text-xs"><code>{`grit expose form Contact --to apps/web/app/contact-us/page.tsx
grit expose table Contact --to apps/web/app/contacts/page.tsx`}</code></pre>

                <ul>
                  <li>
                    Each command parses{' '}
                    <code>apps/api/internal/models/&lt;snake&gt;.go</code>{' '}
                    to determine the resource's primitive fields.
                    Relationship pointers (<code>Group *Group</code> or{' '}
                    <code>Group Group</code>) and slices
                    (<code>Tags []Tag</code>) are filtered out — only
                    fields that can render as one <code>&lt;input&gt;</code>{' '}
                    or one table cell make it through.
                  </li>
                  <li>
                    Both commands refuse to overwrite an existing file
                    unless you pass <code>--force</code> — protects
                    hand-customised pages from accidental loss.
                  </li>
                  <li>
                    Generated pages are plain Tailwind (no admin
                    chrome), suitable for embedding on a marketing
                    site or a customer dashboard.
                  </li>
                </ul>

                <h3>Supporting fixes</h3>
                <ul>
                  <li>
                    <strong>Web hook imports</strong>: the generator
                    now branches its <code>apiClient</code> import
                    path by app — <code>@/lib/api-client</code> for
                    admin, <code>@/lib/api</code> for web. Resolves
                    a pre-existing "Cannot find module" error in
                    web-side resource hooks.
                  </li>
                  <li>
                    <strong>Scaffolded <code>apps/web/lib/api.ts</code></strong>{' '}
                    now re-exports <code>apiClient = api</code> so
                    generated hooks resolve symmetrically across both
                    apps.
                  </li>
                  <li>
                    <strong>Web package.json</strong> gains{' '}
                    <code>@hookform/resolvers</code> as a dep (was
                    only in admin before).
                  </li>
                  <li>
                    <strong>ParseGoStructs</strong> exported from the{' '}
                    <code>internal/generate</code> package for reuse
                    by the new <code>internal/expose</code> package.
                  </li>
                </ul>

                <h3>Known limitations</h3>
                <ul>
                  <li>
                    Generated forms don't use the shared Zod schema for
                    validation — the schema's camelCase field names
                    don't match the API's snake_case JSON keys. Forms
                    submit snake-case keys directly; server-side
                    validation is the source of truth. Add client-side
                    validation by hand if you need it.
                  </li>
                  <li>
                    Forms have one field per primitive column. Custom
                    widgets (rich text, image uploaders, relationship
                    dropdowns) need manual additions after generation.
                  </li>
                  <li>
                    <code>--public-share</code> / <code>--public</code>{' '}
                    flags (post via the public form-share endpoint
                    instead of via the auth'd hook) are still on the
                    roadmap.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.31.20 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.20
                </span>
                <span className="text-sm text-muted-foreground">June 21, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Phase 2 of PLAN_FORMS_AND_SHARING.md — public form sharing.</strong>{' '}
                  Generate a token-protected link for any of your resources and anyone with the link can submit the form, no admin login required. Optional bcrypt password on the share for an extra gate.
                </p>

                <h3>What landed (end-to-end)</h3>
                <ul>
                  <li>
                    <strong>FormShare model</strong> (token, optional
                    bcrypt PasswordHash, enabled, submission count,
                    label). Auto-migrated on{' '}
                    <code>grit migrate</code>.
                  </li>
                  <li>
                    <strong>Admin handler + routes</strong>:{' '}
                    <code>GET/POST/PATCH/DELETE /api/admin/form-shares</code>.
                  </li>
                  <li>
                    <strong>Public handler + routes</strong> (no auth, no
                    CSRF):
                    {' '}<code>GET /api/public/forms/:token</code> +{' '}
                    <code>POST /api/public/forms/:token/submit</code>.
                    Both paths are listed in Sentinel's{' '}
                    <code>ExcludeRoutes</code> so the WAF doesn&apos;t
                    block public JSON bodies.
                  </li>
                  <li>
                    <strong>Marker-driven resource dispatch</strong>:{' '}
                    every <code>grit generate resource</code> appends a
                    case to{' '}
                    <code>services/form_share_dispatch.go</code> that
                    JSON-decodes the public payload into the resource's
                    model + calls <code>db.Create()</code>. Whitelisted
                    by name — unknown resources can&apos;t be
                    submitted publicly.
                  </li>
                  <li>
                    <strong>Admin page</strong> at{' '}
                    <code>/system/form-shares</code>: list shares,
                    create new (with password), toggle enabled, copy
                    public URL, delete.
                  </li>
                  <li>
                    <strong>Public web page</strong> at{' '}
                    <code>apps/web/app/forms/[token]/page.tsx</code> —
                    a minimal name/email/phone/message form that posts
                    to the public submit endpoint. Tailored forms for
                    other resource shapes come via{' '}
                    <code>grit expose form</code> in Phase 3.
                  </li>
                </ul>

                <h3>Threat model</h3>
                <ul>
                  <li>
                    <strong>Resource whitelisting</strong>: the dispatch
                    service&apos;s switch statement is the gate. A
                    share token can&apos;t conjure a record for a
                    resource that hasn&apos;t been explicitly added.
                  </li>
                  <li>
                    <strong>Field whitelisting</strong>: each resource
                    case JSON-decodes onto its typed model. Unknown
                    JSON keys are silently ignored; private fields
                    (<code>id</code>, <code>created_at</code>, …) are
                    untouched.
                  </li>
                  <li>
                    <strong>Rate limiting</strong>: Sentinel still
                    rate-limits the public path by IP — the WAF body
                    inspection is the only thing skipped.
                  </li>
                  <li>
                    <strong>Password (optional)</strong>: bcrypt cost
                    10. Submitted as <code>_password</code> alongside
                    the fields; rejected with 401 if mismatched.
                  </li>
                </ul>

                <h3>Deferred to v3.31.21</h3>
                <ul>
                  <li>
                    <strong>Audit trail</strong>: a{' '}
                    <code>source_share_id</code> column on each
                    submitted record so admins can filter "show me
                    public submissions" per resource.
                  </li>
                  <li>
                    <strong>Per-resource public form pages</strong>:
                    Phase 3's <code>grit expose form &lt;Resource&gt;</code>{' '}
                    will scaffold a tailored public page with the
                    exact field shape, replacing the generic
                    name+email+phone+message default.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.31.19 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.19
                </span>
                <span className="text-sm text-muted-foreground">June 21, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Column-pack auto-detection.</strong>{' '}
                  Phase 1.4 of <code>PLAN_FORMS_AND_SHARING.md</code>.
                </p>

                <p>
                  Generate a resource with both <code>name</code> and{' '}
                  <code>email</code> (or both <code>first_name</code>{' '}
                  and <code>last_name</code>) and the table now ships
                  with those fields packed into a single stacked
                  column — name on top, email muted below. No
                  hand-written <code>cell:</code> callback needed.
                </p>

                <h3>How it works</h3>
                <ul>
                  <li>
                    New helper:{' '}
                    <code>apps/admin/components/tables/stacked-cell.tsx</code>.
                    Exports a <code>StackedCell({`{ top, bottom }`})</code>{' '}
                    function returning two-line JSX. Called as a
                    function (not JSX syntax) so resource files stay{' '}
                    <code>.ts</code>.
                  </li>
                  <li>
                    Generator now runs a pack-detector over the
                    resource&apos;s field list. When a pattern matches,
                    the absorbed fields are silently skipped and the
                    packed line is emitted in their primary&apos;s
                    slot.
                  </li>
                  <li>
                    Import of <code>StackedCell</code> is conditional —
                    resources without a pack stay clean.
                  </li>
                </ul>

                <h3>Patterns recognised today</h3>
                <ul>
                  <li>
                    <code>name + email</code> → &quot;Contact&quot;
                    column
                  </li>
                  <li>
                    <code>first_name + last_name</code> →
                    &quot;Name&quot; column
                  </li>
                </ul>
                <p>
                  Both are easy to extend in{' '}
                  <code>internal/generate/column_packs.go</code>. Money
                  + currency badge, status + relative date, and a few
                  others are roadmap.
                </p>

                <h3>Existing resources</h3>
                <p>
                  Pre-v3.31.19 resources don&apos;t auto-pack. Either
                  add the pack by hand (the customising-tables lesson
                  has the recipe), or wait for{' '}
                  <code>grit pack table &lt;Resource&gt;</code> in a
                  future release.
                </p>
              </div>
            </div>

            {/* v3.31.18 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.18
                </span>
                <span className="text-sm text-muted-foreground">June 21, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Form groups + per-group PATCH save.</strong>{' '}
                  Phase 1.3 (partial) of PLAN_FORMS_AND_SHARING.md.
                </p>

                <p>
                  Long Update views with 10+ fields used to save the
                  whole record on every click — risky when two operators
                  edit different sections at once, slow because the
                  payload is large, and tedious because every field had
                  to be re-validated.{' '}
                  <strong>Define <code>form.groups</code></strong> and
                  each group renders as its own Card on the Update page,
                  with its own Save button that PATCHes only that
                  group&apos;s fields.
                </p>

                <h3>What landed</h3>
                <ul>
                  <li>
                    <strong>New <code>Patch</code> handler</strong> on
                    every generated resource. Whitelists writable columns
                    so the partial endpoint can&apos;t be tricked into
                    setting <code>id</code> / <code>created_at</code> /
                    <code>deleted_at</code> / <code>version</code> from
                    the client.
                  </li>
                  <li>
                    <strong>PATCH /api/&lt;plural&gt;/:id</strong> route
                    registered alongside PUT for every resource (both
                    standard and role-restricted route groups).
                  </li>
                  <li>
                    <strong><code>usePatchResource(endpoint)</code></strong>{' '}
                    hook in the admin&apos;s use-resource module.
                    Same shape as <code>useUpdateResource</code> but
                    calls PATCH and toasts &quot;Saved&quot; on
                    success.
                  </li>
                  <li>
                    <strong><code>GroupDefinition</code></strong> type
                    on <code>FormDefinition.groups</code>. Each group is{' '}
                    <code>{`{ title, description?, fields: string[], scope?: "create" | "update" | "both" }`}</code>.
                  </li>
                  <li>
                    <strong><code>&lt;UpdateGroups&gt;</code></strong>{' '}
                    component renders each <code>scope: "update"</code>{' '}
                    or <code>"both"</code> group as a separate Card on
                    the Update page (when{' '}
                    <code>formView: "page"</code> +{' '}
                    <code>form.groups</code> are defined).
                  </li>
                  <li>
                    ResourcePage dispatcher: when editing + groups
                    present, route to <code>UpdateGroups</code>;
                    otherwise fall back to the single-form FormPage.
                  </li>
                </ul>

                <h3>The "create-and-update" pattern</h3>
                <p>
                  Use <code>scope: "create"</code> on the minimal
                  required fields and <code>scope: "update"</code> on
                  the rest. Operators get a frictionless Create form;
                  detailed editing happens on the Update page as cards
                  with partial saves.
                </p>

                <h3>Deferred to v3.31.19</h3>
                <p>
                  Group rendering on the Create flow as a multi-step
                  wizard. The existing <code>steps</code> field still
                  works for that. v3.31.19 unifies them so groups drive
                  both contexts.
                </p>
              </div>
            </div>

            {/* v3.31.17 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.17
                </span>
                <span className="text-sm text-muted-foreground">June 21, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Form render modes — sheet / modal / page.</strong>{' '}
                  Phase 1.2 of <code>PLAN_FORMS_AND_SHARING.md</code>.
                </p>

                <p>
                  The <code>formView</code> field on{' '}
                  <code>defineResource</code> now accepts{' '}
                  <code>"sheet"</code> as an explicit value, and{' '}
                  <code>"modal"</code> renders as a proper centered
                  dialog instead of a sheet. Defaults are unchanged —
                  resources without an explicit <code>formView</code>{' '}
                  still get the long-form-friendly drawer.
                </p>

                <h3>The three rendering choices</h3>
                <ul>
                  <li>
                    <strong><code>"sheet"</code></strong> (default) —
                    right drawer on desktop, bottom sheet on mobile.
                    Best for long forms and multi-line fields.
                  </li>
                  <li>
                    <strong><code>"modal"</code></strong> — centered
                    dialog over a backdrop. Best for short focused
                    forms (1–6 fields).
                  </li>
                  <li>
                    <strong><code>"page"</code></strong> — dedicated
                    route via <code>?action=create|edit</code>. Best
                    for very long forms or anything that needs
                    shareable URLs.
                  </li>
                </ul>

                <h3>What shipped</h3>
                <ul>
                  <li>
                    New <code>FormSheet</code> component
                    (apps/admin/components/forms/form-sheet.tsx) — the
                    long-form-friendly drawer, formerly the
                    implementation of FormModal.
                  </li>
                  <li>
                    <code>FormModal</code> rewritten as a proper
                    centered dialog (max-w-md, backdrop blur, padding).
                  </li>
                  <li>
                    <code>ResourcePage</code> dispatcher picks the
                    right component based on <code>resource.formView</code>.
                  </li>
                  <li>
                    <code>ResourceDefinition</code> type union expanded
                    to include <code>"sheet"</code>.
                  </li>
                </ul>

                <h3>Migrating</h3>
                <p>
                  If you previously set <code>formView: "modal"</code>{' '}
                  explicitly and want the old sheet behavior, change
                  it to <code>"sheet"</code>. Resources that left{' '}
                  <code>formView</code> undefined stay on the sheet —
                  no migration needed.
                </p>
              </div>
            </div>

            {/* v3.31.16 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.16
                </span>
                <span className="text-sm text-muted-foreground">June 21, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>grit sync now auto-adds new model fields to
                  admin resource files.</strong> Phase 1.1 of the{' '}
                  <code>PLAN_FORMS_AND_SHARING.md</code> roadmap.
                </p>

                <p>
                  Add a column to a Go model, run{' '}
                  <code>grit migrate</code> + <code>grit sync</code>,
                  and the field now appears in{' '}
                  <code>apps/admin/resources/&lt;plural&gt;.ts</code>{' '}
                  as both a column and a form input — with a sensible
                  default type inferred from the Go type. Customised
                  entries (labels, helper text, badges, custom{' '}
                  <code>cell</code> renderers) are never touched.
                </p>

                <h3>How it works</h3>
                <p>
                  The generator now emits marker comments around the
                  auto-managed columns + form fields:
                </p>
                <pre className="overflow-x-auto rounded-lg bg-bg-elevated p-3 text-xs"><code>{`columns: [
  // grit:cols:auto-start
  { key: "name", ... },
  // grit:cols:auto-end
],
form: {
  fields: [
    // grit:fields:auto-start
    { key: "name", ... },
    // grit:fields:auto-end
  ],
},`}</code></pre>
                <p>
                  Sync diffs Go model fields against the file. For each
                  field with a <code>key:</code> not found anywhere in
                  the file, it inserts a default entry above the{' '}
                  <code>auto-end</code> marker. Sync is{' '}
                  <strong>insert-only</strong> — it never modifies or
                  removes existing entries.
                </p>

                <h3>Backward compatibility</h3>
                <p>
                  Resources scaffolded before v3.31.16 don&apos;t carry
                  the marker comments. Sync prints a per-resource
                  warning and skips them. To enable auto-add on an
                  existing resource, hand-edit the file to wrap its
                  columns array and form fields array with the four
                  marker lines once. After that, future syncs pick the
                  file up.
                </p>

                <h3>What&apos;s next (Phase 1 continued)</h3>
                <p>
                  v3.31.17+ ships the rest of Phase 1 per{' '}
                  <code>PLAN_FORMS_AND_SHARING.md</code>:
                </p>
                <ul>
                  <li>
                    Form render mode (<code>sheet | modal | page</code>)
                  </li>
                  <li>
                    Form groups (steps in create, cards in update) +
                    PATCH endpoint for per-group saves
                  </li>
                  <li>
                    Column-pack default heuristic +{' '}
                    <code>grit pack table &lt;Resource&gt;</code>
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.31.15 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.15
                </span>
                <span className="text-sm text-muted-foreground">June 21, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Auth UX overhaul + admin polish from a real
                  app-building session.</strong> Seven concrete fixes
                  driven by feedback while building a contact-app on the
                  prior release.
                </p>

                <h3>Framework fixes</h3>
                <ul>
                  <li>
                    <strong>Protected admin routes no longer flash a
                    blank white page</strong> when the session expires
                    or the API restarts. The admin layout now redirects
                    to <code>/login</code> on both network errors AND
                    null user (401), and shows a spinner while the
                    redirect fires. (<code>internal/scaffold/admin_layout_files.go</code>)
                  </li>
                  <li>
                    <strong>Login page bounces to /dashboard</strong>{' '}
                    when the session cookie is still valid — no more
                    seeing the login form while you&apos;re already
                    signed in.
                  </li>
                  <li>
                    <strong>New SessionWatchdog component</strong>{' '}
                    surfaces a modal at 14:30 of idle time with a 30s
                    countdown — &quot;Stay signed in&quot; refreshes
                    via <code>/api/auth/refresh</code>, &quot;Sign
                    out&quot; or timeout calls{' '}
                    <code>useLogout()</code>. Configurable via{' '}
                    <code>NEXT_PUBLIC_SESSION_IDLE_MS</code> and{' '}
                    <code>NEXT_PUBLIC_SESSION_COUNTDOWN_MS</code>.
                  </li>
                  <li>
                    <strong>Sentinel WAF no longer blocks richtext
                    admin POSTs.</strong> The 64 KB body cap is now
                    1 MB (richtext + embedded inline images need it),
                    and admin write endpoints with HTML payloads
                    (<code>/api/blogs</code>, <code>/api/posts</code>,{' '}
                    <code>/api/articles</code>, <code>/api/uploads</code>)
                    are listed under{' '}
                    <code>ExcludeRoutes</code> so the WAF&apos;s XSS
                    detection stops flagging every{' '}
                    <code>&lt;p&gt;</code> tag.
                  </li>
                  <li>
                    <strong>Generated resource tables drop the ID
                    column by default.</strong> UUIDs are noisy and
                    rarely scanned by eye — operators who need it can
                    add it back manually.
                  </li>
                  <li>
                    <strong>
                      ColumnDefinition gains an optional{' '}
                      <code>cell?: (row) =&gt; ReactNode</code>{' '}
                      renderer
                    </strong>{' '}
                    — pack multiple fields into one column (name +
                    email stacked, price + currency badge, status pill
                    + relative date) without dropping out to a
                    hand-written page.tsx. Takes precedence over{' '}
                    <code>format</code> and <code>badge</code> when
                    set.
                  </li>
                  <li>
                    <strong>grit sync prints a heads-up</strong> that
                    it does NOT update the admin resource definition,
                    pointing operators at{' '}
                    <code>apps/admin/resources/&lt;plural&gt;.ts</code>{' '}
                    when a new model field doesn&apos;t show up in the
                    admin form.
                  </li>
                </ul>

                <h3>Chapter 4 — 4 new lessons</h3>
                <p>
                  Chapter 4 now has <strong>11 lessons</strong> (up
                  from 7) — fully covering the post-generation flow:
                </p>
                <ul>
                  <li>
                    <strong>grit remove resource</strong> — the
                    rollback half of the lifecycle.
                  </li>
                  <li>
                    <strong>Customising admin forms</strong> — all 17
                    field types, helper text, multi-step flows.
                  </li>
                  <li>
                    <strong>Customising admin tables</strong> — formats,
                    badges, filters, and the new{' '}
                    <code>cell()</code> render function with three
                    column-packing recipes.
                  </li>
                  <li>
                    <strong>Using the generated API from the web
                    app</strong> — list, search, detail, create form
                    with the auto-generated React Query hook and
                    shared Zod schemas.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.31.13 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.13
                </span>
                <span className="text-sm text-muted-foreground">June 21, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Root <code>.env</code> is now the single source of truth for THEME + SOCIAL_AUTH_ENABLED.</strong>{' '}
                  Setting <code>SOCIAL_AUTH_ENABLED=false</code> in the
                  monorepo&apos;s root <code>.env</code> didn&apos;t hide the
                  Google / GitHub buttons even after a server restart — Next.js
                  only auto-loads <code>.env</code> from its own package
                  directory (<code>apps/admin/</code>, <code>apps/web/</code>),
                  so <code>process.env.SOCIAL_AUTH_ENABLED</code> inside{' '}
                  <code>next.config.ts</code> was <code>undefined</code> and
                  the <code>|| &quot;true&quot;</code> fallback always won.
                </p>
                <p>
                  Both scaffolded <code>next.config.ts</code> files now read
                  the root <code>.env</code> directly via a tiny inline
                  parser before the <code>env</code> block is evaluated. Shell
                  env still wins (only unset keys are filled in), so CI / Docker
                  overrides are unaffected. After upgrading, restart{' '}
                  <code>pnpm dev</code> (Next.js reads env at boot, not on
                  file-watch).
                </p>
              </div>
            </div>

            {/* v3.31.12 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.12
                </span>
                <span className="text-sm text-muted-foreground">June 21, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>System Health / Security / Performance pages build clean.</strong>{' '}
                  The scaffolded <code>/system/health</code>,{' '}
                  <code>/system/security</code>, and <code>/system/performance</code>{' '}
                  pages import five lucide icons (<code>CheckCircle</code>,{' '}
                  <code>Server</code>, <code>HardDrive</code>, <code>Clock</code>,{' '}
                  <code>Gauge</code>) that weren&apos;t re-exported from{' '}
                  <code>apps/admin/lib/icons.ts</code>. A fresh{' '}
                  <code>pnpm --filter ./apps/admin build</code> failed with{' '}
                  <code>Export &lt;Name&gt; doesn&apos;t exist in target module</code>.
                  Added the five names to both the <code>lucide-react</code>{' '}
                  import block and the named re-export block. All 24 admin
                  routes now prerender on a fresh scaffold.
                </p>
              </div>
            </div>

            {/* v3.31.11 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.11
                </span>
                <span className="text-sm text-muted-foreground">June 21, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>air entrypoint points at the built binary, not the source dir.</strong>{' '}
                  v3.31.10 fixed the <code>.exe</code> extension but the
                  scaffolded <code>.air.toml</code> still set{' '}
                  <code>entrypoint = &quot;./cmd/server&quot;</code> — and{' '}
                  <code>air</code> tries to exec the entrypoint as the binary,
                  so Windows hit{' '}
                  <code>CMD will not recognize non .exe file</code>. Per{' '}
                  air&apos;s docs, <code>entrypoint</code> names the built
                  binary (the same role <code>build.bin</code> plays). Fixed
                  to <code>entrypoint = &quot;./tmp/server.exe&quot;</code>.
                  Verified with a live <code>grit start server</code> in a
                  fresh scaffold plus a <code>/api/health</code> curl
                  returning the full database/redis/jobs/email shape.
                </p>
              </div>
            </div>

            {/* v3.31.10 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.31.10
                </span>
                <span className="text-sm text-muted-foreground">June 21, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Scaffolded <code>.air.toml</code> uses an <code>.exe</code> binary on Windows.</strong>{' '}
                  v3.31.9 shipped <code>grit start</code> with{' '}
                  <code>air</code>-backed hot reload, but the generated{' '}
                  <code>.air.toml</code> used{' '}
                  <code>bin = &quot;./tmp/server&quot;</code>. Windows refuses
                  to <code>CreateProcess</code> an extension-less file, so
                  starting the dev loop popped a{' '}
                  <em>&quot;Select an app to open &apos;server&apos;&quot;</em>{' '}
                  dialog instead of running the API. Switched to{' '}
                  <code>cmd = &quot;go build -o ./tmp/server.exe ./cmd/server&quot;</code>{' '}
                  in the scaffolded template so Windows can execute the
                  output directly.
                </p>
              </div>
            </div>

            {/* v3.27.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.27.0
                </span>
                <span className="text-sm text-muted-foreground">June 20, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Admin auth sweep + fresh-scaffold type-clean.</strong>{' '}
                  Closes the last gap left by v3.26.0&apos;s HttpOnly cookie
                  story: the admin app now uses cookies end-to-end too,{' '}
                  <code>js-cookie</code> is gone from both frontends, OAuth
                  no longer leaks tokens via URL params, and both{' '}
                  <code>apps/web</code> + <code>apps/admin</code> return
                  zero type errors on a fresh scaffold for the first time.
                </p>

                <h3>Admin uses HttpOnly cookies</h3>
                <ul>
                  <li>
                    <code>apps/admin/lib/api-client.ts</code> drops{' '}
                    <code>js-cookie</code>, adds{' '}
                    <code>withCredentials: true</code>, and echoes the
                    <code> grit_csrf</code> cookie into{' '}
                    <code>X-CSRF-Token</code> on every mutation.
                  </li>
                  <li>
                    <code>apps/admin/hooks/use-auth.ts</code> imports{' '}
                    <code>User</code>, <code>LoginRequest</code>,{' '}
                    <code>RegisterRequest</code>, <code>AuthResponse</code>,{' '}
                    <code>ApiResponse</code> from <code>@repo/shared/types</code>{' '}
                    instead of declaring them inline. <code>useMe</code>{' '}
                    returns <code>null</code> on 401 instead of throwing.
                    <code> useLogout</code> doesn&apos;t clear tokens locally
                    — the API does it via <code>Set-Cookie</code>.
                  </li>
                  <li>
                    The admin root redirect page (<code>app/page.tsx</code>)
                    drops the <code>Cookies.get(&apos;access_token&apos;)</code>{' '}
                    check (which couldn&apos;t see HttpOnly cookies anyway)
                    in favour of a <code>useMe()</code> probe.
                  </li>
                  <li>
                    The 401-refresh interceptor now POSTs{' '}
                    <code>/api/auth/refresh</code> with an empty body — the
                    API reads <code>grit_refresh</code> from the cookie and
                    issues a new <code>grit_access</code> via{' '}
                    <code>Set-Cookie</code>.
                  </li>
                  <li>
                    <code>profile delete</code> drops <code>Cookies.remove</code>{' '}
                    calls — the Go <code>DeleteProfile</code> handler now
                    calls <code>ClearAuthCookies</code> as part of its
                    response.
                  </li>
                  <li>
                    <code>js-cookie</code> and{' '}
                    <code>@types/js-cookie</code> dropped from{' '}
                    <code>apps/admin/package.json</code>.
                  </li>
                </ul>

                <h3>OAuth without URL leakage</h3>
                <p>
                  The Go OAuth callback handler now calls{' '}
                  <code>SetAuthCookies</code> BEFORE the 307 redirect to{' '}
                  <code>/auth/callback</code>. The cookies travel on the
                  redirect response itself, so the callback page no longer
                  needs to read <code>access_token</code> and{' '}
                  <code>refresh_token</code> from the URL. Tokens never
                  appear in browser history, server access logs, or Referer
                  headers. Closes the gap left when v3.26.5 fixed
                  email/password.
                </p>

                <h3>UUID vs number ID drift cleaned up</h3>
                <ul>
                  <li>
                    <code>useBulkDeleteResource</code> signature{' '}
                    <code>ids: number[]</code> → <code>ids: string[]</code>{' '}
                    (Grit&apos;s models all use UUID primary keys).
                  </li>
                  <li>
                    <code>form-modal.tsx</code> + <code>form-page.tsx</code>{' '}
                    + their <code>-steps</code> variants stop casting{' '}
                    <code>item.id</code> / <code>editId</code> to{' '}
                    <code>Number</code> — they pass the IDs through as
                    strings, matching{' '}
                    <code>useUpdateResource</code> /{' '}
                    <code>useResourceItem</code>.
                  </li>
                  <li>
                    <code>relationship-select-field.tsx</code> +{' '}
                    <code>multi-relationship-select-field.tsx</code> use{' '}
                    <code>String(item.id)</code> instead of asserting{' '}
                    <code>as number</code>.
                  </li>
                  <li>
                    <code>hooks/use-system.ts</code> dropped its inline{' '}
                    <code>Upload</code> interface and imports from{' '}
                    <code>@repo/shared/types</code>;{' '}
                    <code>UploadListResponse</code> is now an alias for{' '}
                    <code>PaginatedResponse&lt;Upload&gt;</code>.
                  </li>
                  <li>
                    Admin icon map: <code>Cpu</code>, <code>Zap</code>,{' '}
                    <code>Globe</code> added; <code>Shield</code> was
                    imported but not re-exported — fixed. System
                    observability / security pages corrected from{' '}
                    <code>@/lib/api</code> to <code>@/lib/api-client</code>.
                  </li>
                </ul>

                <p>
                  Net effect: a fresh <code>grit new</code> →{' '}
                  <code>pnpm install</code> →{' '}
                  <code>pnpm exec tsc --noEmit</code> on{' '}
                  <code>apps/web</code> AND <code>apps/admin</code> returns
                  zero type errors for the first time. The Go API builds
                  and template tests pass unchanged.
                </p>
              </div>
            </div>

            {/* v3.26.5 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.26.5
                </span>
                <span className="text-sm text-muted-foreground">June 20, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Web hooks finally consume <code>packages/shared</code> + use the v3.26.0 HttpOnly cookie auth.</strong>{' '}
                  Closes a contradiction a learner spotted: the docs teach
                  &quot;shared types live in <code>packages/shared</code>&quot; but
                  the scaffolded <code>use-blogs</code> and <code>use-auth</code>{' '}
                  hooks duplicated <code>User</code> and <code>Blog</code> inline.
                </p>

                <h3>Type imports</h3>
                <ul>
                  <li>
                    <code>use-blogs.ts</code> now imports{' '}
                    <code>Blog</code> + <code>PaginatedResponse</code> from{' '}
                    <code>@repo/shared/types</code>.
                  </li>
                  <li>
                    <code>use-auth.ts</code> imports <code>User</code>,{' '}
                    <code>LoginRequest</code>, <code>RegisterRequest</code>,{' '}
                    <code>AuthResponse</code>, <code>ApiResponse</code> from
                    the same barrel.
                  </li>
                  <li>
                    <code>lib/auth-provider.tsx</code> imports <code>User</code>{' '}
                    from <code>@repo/shared/types</code> instead of a 10-line
                    local copy.
                  </li>
                  <li>
                    Web app now has <code>@repo/shared: workspace:*</code> in
                    its <code>package.json</code> (was missing).{' '}
                    <code>next.config.ts</code> gets{' '}
                    <code>transpilePackages: [&apos;@repo/shared&apos;]</code>{' '}
                    so SWC picks up the TS source.
                  </li>
                </ul>

                <h3>Auth flow uses HttpOnly cookies end-to-end</h3>
                <ul>
                  <li>
                    Axios client gets <code>withCredentials: true</code> so
                    the browser actually attaches the{' '}
                    <code>grit_access</code> / <code>grit_refresh</code>{' '}
                    cookies the API issues.
                  </li>
                  <li>
                    A request interceptor echoes the <code>grit_csrf</code>{' '}
                    cookie into <code>X-CSRF-Token</code> on every
                    state-changing request — required by the AutoCSRF
                    middleware that v3.26.0 wired in.
                  </li>
                  <li>
                    <code>use-auth.ts</code> dropped <code>js-cookie</code>,{' '}
                    <code>storeTokens</code> / <code>clearTokens</code> /{' '}
                    <code>getAccessToken</code>, and every{' '}
                    <code>Authorization: Bearer</code> header attachment.{' '}
                    <code>useMe</code> returns <code>null</code> on 401
                    instead of throwing.
                  </li>
                  <li>
                    Login + register pages no longer call{' '}
                    <code>Cookies.set(&apos;access_token&apos;)</code>. The
                    API sets HttpOnly cookies via{' '}
                    <code>Set-Cookie</code>; the browser stores them; JS
                    never touches tokens.
                  </li>
                </ul>

                <h3>Known gaps deferred to a follow-up release</h3>
                <ul>
                  <li>
                    OAuth callback page still reads tokens from URL params
                    and sets them via <code>Cookies.set</code>. Proper fix
                    requires the Go-side OAuth handler to set cookies before
                    redirecting.
                  </li>
                  <li>
                    Admin app still uses <code>js-cookie</code> + Bearer
                    header auth throughout. Bigger refactor (TOTP, OAuth
                    begin/callback, profile page, multi-step token refresh)
                    that warrants its own release.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.26.4 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.26.4
                </span>
                <span className="text-sm text-muted-foreground">June 15, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong><code>grit start</code> now actually starts both, like the help said it would.</strong>
                </p>

                <h3>The bug</h3>
                <p>
                  In a web project, <code>grit start</code> (no subcommand) was
                  falling through to <code>cmd.Help()</code> — printing the
                  available subcommands and exiting. The command&apos;s own
                  <code>Long</code> description said it would start both the
                  API and the client, but only <code>grit start server</code>{' '}
                  and <code>grit start client</code> actually did anything.
                </p>

                <h3>What changed</h3>
                <ul>
                  <li>
                    <code>grit start</code> in a web project now spawns{' '}
                    <code>go run cmd/server/main.go</code> (in{' '}
                    <code>apps/api/</code>) and <code>pnpm dev</code> (at the
                    project root) in parallel.
                  </li>
                  <li>
                    Output from both processes is streamed to the same
                    terminal with a coloured <code>[api]</code> /{' '}
                    <code>[web]</code> prefix per line so a developer can tell
                    whose log is whose without splitting panes.
                  </li>
                  <li>
                    Ctrl+C (and SIGTERM) is forwarded to both children. If
                    either child exits on its own, the other is shut down too
                    — no zombie processes left behind.
                  </li>
                  <li>
                    Desktop projects are unchanged — <code>grit start</code>{' '}
                    still calls <code>wails dev</code>. The subcommands{' '}
                    <code>grit start server</code> and{' '}
                    <code>grit start client</code> still work if you only want
                    one side.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.26.3 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.26.3
                </span>
                <span className="text-sm text-muted-foreground">June 15, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Redis + MinIO host ports moved to dodge native-install collisions.</strong>{' '}
                  Same pattern as v3.26.2 did for Postgres, applied to the two
                  other ports learners actually clash on.
                </p>

                <h3>What changed</h3>
                <ul>
                  <li>
                    <strong>Redis:</strong> dev host port <code>6379 → 6380</code>.
                    Native installs (Memurai on Windows, <code>brew install redis</code>,{' '}
                    <code>apt install redis-server</code>, WSL Redis) all bind 6379.
                  </li>
                  <li>
                    <strong>MinIO S3 API:</strong> dev host port <code>9000 → 9002</code>.
                    Portainer&apos;s admin UI defaults to 9000; SonarQube and a
                    handful of monitoring stacks grab it too.
                  </li>
                  <li>
                    <strong>MinIO console:</strong> dev host port <code>9001 → 9003</code>.
                    Less common collision but kept in sync with the API port shift.
                  </li>
                  <li>
                    <strong>Mailhog kept on 1025 / 8025.</strong> Almost zero dev
                    machines have anything on those — shifting them just adds
                    learner confusion without preventing real failures.
                  </li>
                </ul>

                <h3>Inside the Docker network</h3>
                <p>
                  Containers still listen on the canonical ports — Redis on 6379,
                  MinIO on 9000 + 9001. The host-port shifts only affect how you
                  reach them from your laptop. Prod compose is unchanged because
                  inter-container traffic uses the docker network hostnames
                  (<code>redis</code>, <code>minio</code>) and container ports.
                </p>

                <h3>Where else this surfaces</h3>
                <ul>
                  <li>
                    <code>.env</code>: <code>REDIS_URL=redis://localhost:6380</code>{' '}
                    and <code>MINIO_ENDPOINT=http://localhost:9002</code>.
                  </li>
                  <li>
                    The CLI &quot;next steps&quot; banner after <code>grit new</code> now
                    prints the new host ports.
                  </li>
                  <li>
                    Scaffolded README&apos;s services table and the docs lessons
                    (Docker primer, dev-servers, batteries/redis-cache,
                    batteries/s3-storage) updated to match.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.26.2 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.26.2
                </span>
                <span className="text-sm text-muted-foreground">June 15, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Postgres host port 5432 → 5434 to dodge Windows
                  WinNAT reservations.</strong>{' '}
                  Closes the &quot;An attempt was made to access a socket in
                  a way forbidden by its access permissions&quot; bind error
                  that hit Windows users on a fresh{' '}
                  <code>docker compose up -d</code>.
                </p>

                <h3>What was wrong</h3>
                <p>
                  Even with no process visibly holding port 5432, Docker
                  Desktop on Windows would refuse to bind it. Cause: the
                  WinNAT service / Hyper-V Virtual Switch silently reserves
                  TCP port ranges at boot, and 5432 sits inside one of the
                  common reservations on Docker Desktop + WSL2 default
                  installs.
                </p>

                <h3>What changed</h3>
                <ul>
                  <li>
                    Dev <code>docker-compose.yml</code> now publishes Postgres
                    on host port <code>5434</code> (not <code>5432</code>).
                    Container port stays <code>5432</code> inside the Docker
                    network.
                  </li>
                  <li>
                    <code>.env</code> sets <code>POSTGRES_PORT=5434</code> so
                    the Go API connects to the same host port.
                  </li>
                  <li>
                    <code>docker-compose.prod.yml</code> pins{' '}
                    <code>POSTGRES_PORT=5432</code> in the api service
                    environment because inter-container traffic uses the
                    container port, not the dev host port.
                  </li>
                  <li>
                    Docker primer lesson gains error 2a, covering the
                    Windows-specific bind error with the{' '}
                    <code>netsh int ipv4 show excludedportrange</code>{' '}
                    diagnostic and three fix paths.
                  </li>
                </ul>

                <p>
                  Net effect: a fresh <code>grit new</code> →{' '}
                  <code>docker compose up -d</code> →{' '}
                  <code>grit migrate</code> now succeeds on Windows machines
                  whose Hyper-V reservation overlaps 5432 — without any
                  user-side intervention.
                </p>
              </div>
            </div>

            {/* v3.26.1 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.26.1
                </span>
                <span className="text-sm text-muted-foreground">June 15, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Single source of truth for Postgres credentials — fresh scaffolds Just Work.</strong>{' '}
                  Closes the &quot;SQLSTATE 28P01: password authentication failed&quot;
                  trap that bit every learner whose <code>.env</code> and{' '}
                  <code>docker-compose.yml</code> drifted apart.
                </p>

                <h3>The bug</h3>
                <p>
                  v3.25.x and v3.26.0 scaffolds wrote three disagreeing copies
                  of the DB credentials: <code>docker-compose.yml</code> hardcoded{' '}
                  <code>grit:grit</code>, <code>.env</code>&apos;s{' '}
                  <code>DATABASE_URL</code> used <code>grit:grit</code>, and{' '}
                  <code>.env</code>&apos;s <code>POSTGRES_PASSWORD</code> said{' '}
                  <code>change-me-in-production</code>. The moment a learner
                  edited one, the others were out of sync and <code>grit migrate</code>{' '}
                  failed.
                </p>

                <h3>The fix</h3>
                <ul>
                  <li>
                    One canonical <code>POSTGRES_*</code> block in{' '}
                    <code>.env</code> — <code>POSTGRES_USER</code> /{' '}
                    <code>POSTGRES_PASSWORD</code> / <code>POSTGRES_DB</code> /{' '}
                    <code>POSTGRES_HOST</code> / <code>POSTGRES_PORT</code>.
                  </li>
                  <li>
                    <code>POSTGRES_PASSWORD</code> is now generated at scaffold
                    time as a 48-hex-char random string (alongside JWT_SECRET,
                    PULSE_PASSWORD, etc.) — even <code>APP_ENV=production</code>{' '}
                    is safe on first boot.
                  </li>
                  <li>
                    <code>docker-compose.yml</code> reads from <code>.env</code>{' '}
                    via <code>${'${VAR:-grit}'}</code> substitution. No hardcoded
                    credentials anywhere.
                  </li>
                  <li>
                    The Go API builds <code>DATABASE_URL</code> from the same
                    <code> POSTGRES_*</code> parts at startup. Set{' '}
                    <code>DATABASE_URL</code> only if you want to point at
                    external Postgres (Neon, Supabase, RDS) or SQLite — it&apos;s
                    the explicit escape hatch.
                  </li>
                  <li>
                    Prod compose now sets <code>POSTGRES_HOST=postgres</code> in
                    the api environment so the Go binary finds the postgres
                    container on the docker network. No embedded <code>DATABASE_URL</code>{' '}
                    in the compose YAML anymore.
                  </li>
                </ul>

                <p>
                  A fresh <code>grit new</code> → <code>docker compose up -d</code>{' '}
                  → <code>grit migrate</code> now succeeds without any editing.
                </p>
              </div>
            </div>

            {/* v3.26.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.26.0
                </span>
                <span className="text-sm text-muted-foreground">June 14, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Concepts ch.2 expansion, Docker hardening, and the long-overdue grit.json fix.</strong>{' '}
                  A teaching + security release driven by real student feedback.
                </p>

                <h3>Concepts course chapter 2</h3>
                <ul>
                  <li>
                    <strong>&quot;Tour of your project&quot; lesson rewritten end-to-end.</strong>{' '}
                    The original 30-second map was missing most of what the scaffold
                    actually produces. The new lesson walks every folder + file
                    matched against a real fresh scaffold: every package under
                    <code> apps/api/internal</code> (25+ packages in a reference table),
                    the full <code>apps/web</code> + <code>apps/admin</code> layouts,
                    <code> packages/shared</code> + <code>packages/grit-ui</code>,
                    <code> tests/k6</code> (6 scripts), <code>e2e/</code> (Playwright),
                    <code> .claude/</code>, <code>.github/</code>, and every root
                    config file.
                  </li>
                  <li>
                    <strong>New &quot;A Docker primer&quot; lesson</strong> inserted
                    between project-tour and dev-servers. Most learners stall at{' '}
                    <code>docker compose up -d</code> because nobody taught Docker
                    first. The new lesson covers what Docker is, image/container/volume,
                    install per OS, how Grit uses Docker, the 12 commands you&apos;ll
                    actually type, the 6 errors learners hit + their fixes, plus a
                    full <strong>&quot;Run Grit without Docker&quot;</strong> path
                    (Neon + Upstash + Cloudflare R2 + Resend).
                  </li>
                </ul>

                <h3>Docker scaffold hardening</h3>
                <ul>
                  <li>
                    <strong><code>docker-compose.yml</code> binds every port to{' '}
                    <code>127.0.0.1</code></strong>, not the Docker default{' '}
                    <code>0.0.0.0</code>. Coffee-shop wifi can no longer reach
                    Postgres with <code>grit:grit</code> credentials.
                  </li>
                  <li>
                    <strong><code>docker-compose.prod.yml</code> documented as
                    reverse-proxy-first.</strong> A top-of-file comment block
                    spells out the security posture: nothing uses{' '}
                    <code>ports:</code>, only <code>expose:</code>. Postgres + Redis
                    have NO host binding at all in prod. Traffic must arrive via
                    Traefik / Caddy / nginx / Dokploy on the same Docker network.
                  </li>
                </ul>

                <h3>Scaffold fixes</h3>
                <ul>
                  <li>
                    <strong><code>grit.json</code> now writes the real CLI version.</strong>{' '}
                    Previously hardcoded to <code>3.3.0</code> (a leftover placeholder
                    from when the grit.json schema was at 3.3). Fresh projects now
                    show the actual scaffolding CLI version (3.26.0 today). Closes
                    the &quot;is my project version really 3.3?&quot; confusion.
                  </li>
                </ul>

                <h3>Docs site</h3>
                <ul>
                  <li>
                    Single-source-of-truth version constant in{' '}
                    <code>config/site.ts</code>. The header badge, install lesson
                    example output, verify-install lesson, animated terminal, and
                    changelog all read from one place — no more drift between the
                    CLI version and what the website shows.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.25.2 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.25.2
                </span>
                <span className="text-sm text-muted-foreground">May 31, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Smarter <code>grit update</code> + docs sweep.</strong>{' '}
                  Two follow-ups to the v3.25 install/update flow.
                </p>

                <h3>Update command</h3>
                <ul>
                  <li>
                    <strong>Short-circuits when already on latest.</strong> The
                    version check that previously only ran on the GitHub-binary
                    path is now lifted to the top of <code>grit update</code>{' '}
                    — both the Go-install and GitHub-binary strategies skip
                    their work when there&apos;s nothing to do. One HTTP round-trip,
                    then exit. (Was: always rename + go install + cleanup,
                    even when already current.)
                  </li>
                  <li>
                    <strong>Unix path no longer deletes before installing.</strong>{' '}
                    POSIX keeps the running process&apos;s inode alive when the
                    file at the same path is overwritten, so <code>go install</code>{' '}
                    can write straight on top. We previously did{' '}
                    <code>os.Remove</code> first, which left the user stranded
                    if <code>go install</code> failed.
                  </li>
                  <li>
                    <strong>Windows rename now rolls back on failure.</strong>{' '}
                    The .exe-locked-while-running dance (rename current to{' '}
                    <code>.old</code> → write new → delete .old) now restores
                    the original binary if <code>go install</code> fails before
                    writing the new one, so a flaky network or proxy issue
                    can&apos;t leave the user with no usable <code>grit</code>.
                  </li>
                </ul>

                <h3>Docs sweep</h3>
                <ul>
                  <li>
                    Replaced every <code>go install github.com/MUKE-coder/grit/...</code>{' '}
                    install reference across docs pages, tutorials, courses,
                    and the structured-data FAQ schema with the v3.25 one-line
                    install script (with <code>go install</code> kept as a
                    secondary option for power users with Go installed).
                  </li>
                  <li>
                    Hero terminal animation now opens with{' '}
                    <code>curl -fsSL https://gritframework.dev/install.sh | sh</code>{' '}
                    instead of <code>go install</code>.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.25.1 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.25.1
                </span>
                <span className="text-sm text-muted-foreground">May 31, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Scaffold now generates real secrets — and SQLite uncommenting is clean.</strong>{' '}
                  Two papercuts that turned into roadblocks the moment anyone flipped{' '}
                  <code>APP_ENV=production</code> on a fresh project.
                </p>

                <h3>Fixes</h3>
                <ul>
                  <li>
                    <strong>Random secrets at scaffold time.</strong>{' '}
                    <code>grit new</code> now generates cryptographically random values
                    (<code>crypto/rand</code> → hex) for <code>JWT_SECRET</code>,{' '}
                    <code>SENTINEL_PASSWORD</code>, <code>SENTINEL_SECRET_KEY</code>, and{' '}
                    <code>PULSE_PASSWORD</code> when writing <code>.env</code>. Previously
                    these shipped as <code>your-super-secret-...</code> and{' '}
                    <code>admin/sentinel</code> / <code>admin/pulse</code>, which Sentinel v2
                    and Pulse v1 explicitly refuse to start with in release mode. A fresh
                    scaffold now boots cleanly in production mode with both dashboards
                    mounted — no manual <code>openssl rand</code> step required.
                  </li>
                  <li>
                    <strong>Clean SQLite uncomment line.</strong> The commented-out SQLite
                    DSN previously had multiple leading spaces (<code># &nbsp;&nbsp;DATABASE_URL=sqlite:...</code>),
                    so removing the <code>#&nbsp;</code> left a line with leading
                    whitespace. godotenv tolerated it, but it was ugly and confusing.
                    Now single-<code>#</code>-prefixed for a clean uncomment.
                  </li>
                  <li>
                    <strong>k6 tutorial: turn Sentinel + Pulse off for the bench.</strong>{' '}
                    Both sit in the request middleware chain. Leaving them on while
                    load-testing means measuring <em>them</em>, not Gin. Step 3 now sets{' '}
                    <code>SENTINEL_ENABLED=false</code> +{' '}
                    <code>PULSE_ENABLED=false</code> alongside the SQLite switch.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.25.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.25.0
                </span>
                <span className="text-sm text-muted-foreground">May 31, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>One command to update: <code>grit update</code>.</strong>{' '}
                  The CLI now checks GitHub for the latest version, compares
                  against the running binary, and — depending on whether Go is
                  on your PATH — either runs <code>go install ...@latest</code> or
                  downloads the matching binary from the GitHub release and
                  swaps it in place. Atomic swap is handled by{' '}
                  <code>inconshreveable/go-update</code>, so it works correctly
                  even on Windows where you can&apos;t overwrite a running binary.
                </p>

                <h3>What changed</h3>
                <ul>
                  <li>
                    <strong>Smart <code>grit update</code></strong> — first
                    checks GitHub releases. If you&apos;re already on latest, exits
                    in a single round-trip with{' '}
                    <code>Already on the latest version</code>. No more wasted{' '}
                    <code>go install</code> runs.
                  </li>
                  <li>
                    <strong>No-Go-toolchain mode.</strong> If the{' '}
                    <code>go</code> binary isn&apos;t on PATH, <code>grit update</code>{' '}
                    falls back to the GitHub-binary path automatically. Means
                    grit can keep itself current even for users who installed
                    from the prebuilt archive and never touched Go.
                  </li>
                  <li>
                    <strong><code>--from-release</code> flag</strong> forces the
                    GitHub-binary path even when Go is installed. Useful if
                    you&apos;re behind a corporate proxy that blocks the module
                    proxy but allows github.com.
                  </li>
                  <li>
                    <strong>Alias <code>grit self-update</code></strong> for
                    discoverability — same command, clearer intent than{' '}
                    <code>update</code>.
                  </li>
                </ul>

                <h3>How to get it</h3>
                <p>
                  This is the bootstrap release — you need to install v3.25.0
                  manually one time, then future versions are a single command:
                </p>
                <pre><code>{`# from any directory, with Go on PATH:
go install github.com/MUKE-coder/grit/v3/cmd/grit@v3.25.0

# from then on:
grit update`}</code></pre>
              </div>
            </div>

            {/* v3.24.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.24.0
                </span>
                <span className="text-sm text-muted-foreground">May 31, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>SQLite, AI prompts everywhere, and a Learnings journal.</strong>{' '}
                  Three shipping threads: (1) the scaffolded API now speaks SQLite, not
                  just Postgres — flip <code>DATABASE_URL=sqlite:./app.db</code> in{' '}
                  <code>.env</code> and you skip Docker entirely; (2) every tech-kit page
                  now ships a copyable starter prompt for{' '}
                  <Link href="https://claude.ai" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">claude.ai</Link>{' '}
                  plus a new four-step{' '}
                  <Link href="/docs/ai-integration" className="text-primary hover:underline">AI Integration wizard</Link>{' '}
                  that generates a tailored planning prompt; (3) a new{' '}
                  <Link href="/docs/learnings" className="text-primary hover:underline">Learnings</Link>{' '}
                  section opens an engineering journal — first entry walks a stateless
                  service + k6 load test from <code>grit new --api</code> all the way
                  to a committed p50/p95/p99 latency chart.
                </p>

                <h3>Scaffold &amp; framework</h3>
                <ul>
                  <li>
                    <strong>SQLite support.</strong>{' '}
                    <code>internal/database/database.Connect</code> now branches on
                    DSN prefix: <code>sqlite://path</code>, <code>sqlite:path</code>,
                    <code> sqlite::memory:</code>, or Postgres for anything else. Uses{' '}
                    <code>github.com/glebarez/sqlite</code> (pure Go, no CGO) so it
                    works on Windows without a C toolchain. Existing Postgres setups
                    are unchanged.
                  </li>
                  <li>
                    <strong>.env documentation.</strong> Both <code>.env</code> and{' '}
                    <code>.env.example</code> now show all three DSN shapes inline.
                  </li>
                  <li>
                    <strong>Demo DEMO_MODE bypass.</strong>{' '}
                    Sentinel v2 + Pulse v1 refuse to start in release mode with default
                    credentials. The public Grit demo now opts in via{' '}
                    <code>DEMO_MODE=true</code> in <code>demo/internal/routes/routes.go</code>{' '}
                    — production deploys still get the gate; the publicly pokeable demo
                    skips it.
                  </li>
                </ul>

                <h3>Docs site</h3>
                <ul>
                  <li>
                    <strong>Per-kit starter prompts.</strong> All 7{' '}
                    <Link href="/docs/tech-kits" className="text-primary hover:underline">tech-kit pages</Link>{' '}
                    (single, single-vite, double, triple, api, mobile, desktop) now
                    have a &quot;Plan this kit with an AI&quot; section. One copy
                    button gives you a prompt to paste into claude.ai with your idea —
                    Claude returns project-description.md, project-phases.md,
                    design-style-guide.md, and prompt.md, the four planning files
                    you feed to Claude Code.
                  </li>
                  <li>
                    <strong>New <Link href="/docs/ai-integration" className="text-primary hover:underline">AI Integration page</Link>.</strong>{' '}
                    A four-step wizard (Platform → Tech Kit → Use case → Your prompt)
                    that customizes the prompt per project shape. Modelled on the
                    DGateway integration helper.
                  </li>
                  <li>
                    <strong>New Learnings section.</strong>{' '}
                    <Link href="/docs/learnings" className="text-primary hover:underline">/docs/learnings</Link>{' '}
                    is an engineering journal. The first entry — a stateless service
                    + k6 load test walkthrough — covers everything end to end:
                    scaffold, install k6, smoke test, average-load profile,
                    JSON output, three charting options, percentile interpretation
                    table, and committing the milestone.
                  </li>
                  <li>
                    <strong>Navbar trim.</strong> Dropped Stack Selector and Tutorials
                    from the top nav (still reachable from the sidebar and search).
                    Replaced with the AI Integration link as a top-level highlighted item.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.23.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.23.0
                </span>
                <span className="text-sm text-muted-foreground">May 29, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Security &amp; deploy hardening release.</strong> Three threads
                  in one ship: (1) the deploy-day fixes uncovered by a real production
                  build on <code>--single --vite</code>; (2) the React/Vite UI primitives
                  every Grit project ends up writing by hand; (3) a full pass against
                  the OWASP Top 10:2025 with code-level defences and a documented
                  testing methodology. Two new docs pages — <code>/docs/security</code>{' '}
                  and <code>/docs/testing</code> — are the audit checklist clients will
                  walk with you.
                </p>

                <h3>Deploy reliability</h3>
                <ul>
                  <li>
                    <strong><code>cron.Start(cfg, cache)</code> ships.</strong>{' '}
                    The single-app <code>main.go</code> was importing a function that
                    didn&apos;t exist; project wouldn&apos;t compile out of the box.
                    The new helper wraps the existing <code>Scheduler</code> and
                    returns <code>(*Scheduler, error)</code> so callers can stop it on
                    shutdown.
                  </li>
                  <li>
                    <strong>asynq worker now actually starts.</strong>{' '}
                    The single-app <code>main.go</code> was queueing jobs without ever
                    starting <code>jobs.StartWorker</code>, so the token-cleanup task
                    and every email/SMS/cron job sat in Redis forever. Worker startup +
                    graceful shutdown wired in.
                  </li>
                  <li>
                    <strong>SPA fallback no longer loops behind reverse proxies.</strong>{' '}
                    <code>c.FileFromFS(&quot;index.html&quot;, ...)</code> triggered{' '}
                    <code>http.FileServer</code>&apos;s canonical-URL 301 rule and
                    ping-ponged forever behind Traefik / Cloudflare
                    (<code>ERR_TOO_MANY_REDIRECTS</code>). The scaffold now pre-reads{' '}
                    <code>index.html</code> once and serves via <code>c.Data()</code>.
                  </li>
                  <li>
                    <strong>UUID primary keys no longer 401 every request.</strong>{' '}
                    GORM&apos;s <code>db.First(&amp;user, id)</code> shorthand assumes
                    an integer PK; with the scaffold&apos;s UUID-string PK Postgres
                    rejected it with{' '}
                    <em>&quot;trailing junk after numeric literal&quot;</em>. All eight
                    call sites (auth middleware + UserHandler + TOTP) switched to{' '}
                    <code>db.Where(&quot;id = ?&quot;, id).First(...)</code>.
                  </li>
                  <li>
                    <strong>Single-app layout cleaned up.</strong>{' '}
                    <code>main.go</code> moved to the project root so{' '}
                    <code>//go:embed all:frontend/dist</code> resolves on a fresh
                    clone; a placeholder <code>frontend/dist/index.html</code> ships so{' '}
                    <code>go build</code> works before <code>pnpm build</code>;{' '}
                    <code>apps/api/Dockerfile</code> + the multi-app{' '}
                    <code>docker-compose.prod.yml</code> are no longer generated in{' '}
                    <code>--single</code> mode (Dokploy was auto-detecting them and
                    failing). A new root Dockerfile pins pnpm to 9.15.0,{' '}
                    <code>chown</code>s before <code>USER</code> (fixes Sentinel/Pulse{' '}
                    <em>&quot;out of memory (14)&quot;</em> SQLite errors).
                  </li>
                  <li>
                    <strong>Auto-migrate + first-boot seed.</strong>{' '}
                    Single-app <code>main.go</code> now runs{' '}
                    <code>models.Migrate(db)</code> on startup (gate via{' '}
                    <code>AUTO_MIGRATE=false</code>) and seeds when the users table is
                    empty (gate via <code>AUTO_SEED</code>; off by default in
                    production). Fresh-deploy → working-login is a single command.
                  </li>
                  <li>
                    <strong>Vite scaffold fixes</strong> — <code>postcss.config.cjs</code>{' '}
                    (was <code>.js</code>, broke ESM <code>package.json</code>);{' '}
                    <code>api.ts</code> uses <code>import.meta.env.VITE_API_URL</code>{' '}
                    instead of <code>process.env.NEXT_PUBLIC_API_URL</code>;{' '}
                    <code>vite-env.d.ts</code> declares the type so{' '}
                    <code>tsc --noEmit</code> stops erroring; navbar/footer use TanStack
                    Router&apos;s <code>Link</code> + <code>useRouterState</code>{' '}
                    instead of <code>next/link</code> / <code>usePathname</code>;{' '}
                    <code>@tanstack/router-cli</code> in devDeps + a postinstall hook
                    so <code>routeTree.gen.ts</code> exists on a fresh clone.
                  </li>
                </ul>

                <h3>Vite UI primitives (every project ends up writing these)</h3>
                <ul>
                  <li>
                    <strong><code>lib/auth.ts</code></strong> — handles the actual{' '}
                    <code>&#123;data:&#123;user, tokens:&#123;access_token, refresh_token, expires_at&#125;&#125;&#125;</code>{' '}
                    envelope, persists tokens, exports{' '}
                    <code>login</code> / <code>register</code> / <code>me</code> /{' '}
                    <code>refresh</code> / <code>logout</code> / <code>clearAuth</code>.
                    Handles the TOTP-challenge response shape too.
                  </li>
                  <li>
                    <strong><code>lib/api.ts</code></strong> — auto-attaches{' '}
                    <code>Authorization: Bearer</code>, transparently retries once on
                    401 via <code>/api/auth/refresh</code>. Single-flight refresh so a
                    burst of 401s doesn&apos;t fan out into N refresh calls.
                  </li>
                  <li>
                    <strong><code>ConfirmDialog</code> + <code>useConfirm</code></strong>{' '}
                    — Promise-based confirm:{' '}
                    <code>const ok = await confirm(&#123; title, message, tone: &apos;danger&apos; &#125;)</code>.
                    Esc cancels, Enter confirms.
                  </li>
                  <li>
                    <strong><code>MoneyInput</code></strong> — Intl.NumberFormat
                    thousands separators, prefix slot,{' '}
                    <code>value: number | null</code>.
                  </li>
                  <li>
                    <strong><code>Combobox</code></strong> — keyboard-friendly
                    (Arrow Up/Down/Enter/Esc), filter on label + sublabel.
                  </li>
                  <li>
                    <strong><code>SessionExpiryMonitor</code></strong> — decodes JWT{' '}
                    <code>exp</code>, shows a Stay / Logout modal 30s before expiry.
                    Stay calls <code>refresh()</code>.
                  </li>
                  <li>
                    <strong><code>StatusBadge&lt;TStatus&gt;</code></strong> — typed
                    status → tone (success / warning / danger / info / neutral) with a
                    default tone map for paid / pending / overdue / etc.
                  </li>
                  <li>
                    <strong><code>StatsRow</code></strong> — list-page stat cards
                    (label, value, sub, icon, tone).
                  </li>
                </ul>

                <h3>OWASP Top 10:2025 hardening</h3>
                <p>
                  Every fresh <code>grit new</code> project now ships defences for
                  every category by default. The new{' '}
                  <Link href="/docs/security" className="text-primary hover:underline">
                    Security Guide
                  </Link>{' '}
                  walks each one category-by-category.
                </p>
                <ul>
                  <li>
                    <strong>A01 IDOR — <code>internal/authz</code></strong> ·{' '}
                    <code>authz.MustOwn(c, db, dest, id)</code> returns 404 (not 403) on
                    every failure so existence isn&apos;t leaked through error-message
                    differences. <code>authz.RequireRoles(&quot;admin&quot;)</code>{' '}
                    middleware for admin routes.
                  </li>
                  <li>
                    <strong>A01 SSRF — <code>internal/safefetch</code></strong> · drop-in{' '}
                    <code>safefetch.Get(ctx, url)</code> validates scheme + host
                    pre-flight AND re-checks the resolved IP at TCP-connect time via{' '}
                    <code>net.Dialer.Control</code> — closes the DNS-rebind TOCTOU.
                    Blocks loopback, RFC1918, link-local, CGNAT (100.64/10), AWS IMDS
                    (169.254.169.254 + <code>fd00:ec2::/32</code>),{' '}
                    <code>metadata.google.internal</code>.
                  </li>
                  <li>
                    <strong>A02 — <code>SecurityHeaders</code> middleware extended</strong>{' '}
                    · strict <code>Content-Security-Policy</code>{' '}
                    (<code>default-src &apos;self&apos;</code> + script allowlist +{' '}
                    <code>frame-ancestors &apos;none&apos;</code> + <code>object-src &apos;none&apos;</code>),
                    plus <code>Cross-Origin-Opener-Policy</code> and{' '}
                    <code>Cross-Origin-Resource-Policy</code>. Skipped on{' '}
                    <code>/docs</code>, <code>/studio</code>, <code>/sentinel</code>,{' '}
                    <code>/pulse</code> which serve vendored UIs.
                  </li>
                  <li>
                    <strong>A03 Supply chain</strong> ·{' '}
                    <code>.github/dependabot.yml</code> (Go modules + npm + GitHub
                    Actions, weekly) and{' '}
                    <code>.github/workflows/security.yml</code> running{' '}
                    <code>govulncheck</code> + <code>pnpm audit</code> (high+) +
                    CodeQL Go/JS on every PR and weekly.
                  </li>
                  <li>
                    <strong>A01-adjacent CSRF — <code>middleware.CSRF</code></strong>{' '}
                    · double-submit-cookie defence for cookie-auth routes (OAuth flow).
                    SameSite=Strict on the token cookie.
                  </li>
                  <li>
                    <strong>A09 — <code>middleware.LogSecurityEvent</code></strong>{' '}
                    · typed event constants for login success/failure, logout, password
                    change, TOTP enable/disable, role change, account lock, authZ
                    denial. Rides the existing tamper-evident ActivityLog hash chain.
                  </li>
                </ul>

                <h3>Performance &amp; security testing methodology</h3>
                <ul>
                  <li>
                    <strong>k6 suite</strong> in <code>tests/k6/</code> — all six
                    test types from the testing course (smoke / average-load / stress /
                    spike / soak / breakpoint) share one user journey in{' '}
                    <code>lib/common.js</code>. SLO-aligned thresholds; smoke +
                    average-load suitable as a CI regression gate.
                  </li>
                  <li>
                    <strong>New{' '}
                    <Link href="/docs/testing" className="text-primary hover:underline">
                      /docs/testing
                    </Link>{' '}
                    page</strong> — k6 install + reading results, the 5-phase pentest
                    methodology, the attack catalogue with curl one-liners against a
                    Grit app (IDOR / SQLi / XSS / SSRF / brute-force / misconfig),
                    CVSS scoring + audit-report structure.
                  </li>
                  <li>
                    <strong>Sentinel and Pulse issues</strong> filed for the
                    cross-project improvements this release uncovered: Sentinel{' '}
                    <a href="https://github.com/MUKE-coder/sentinel/issues/2" className="text-primary hover:underline">#2</a>{' '}
                    CSP report endpoint,{' '}
                    <a href="https://github.com/MUKE-coder/sentinel/issues/3" className="text-primary hover:underline">#3</a>{' '}
                    SSRF guard,{' '}
                    <a href="https://github.com/MUKE-coder/sentinel/issues/4" className="text-primary hover:underline">#4</a>{' '}
                    user-scoped rate limit + CAPTCHA,{' '}
                    <a href="https://github.com/MUKE-coder/sentinel/issues/5" className="text-primary hover:underline">#5</a>{' '}
                    CVSS finding model + alerts. Pulse{' '}
                    <a href="https://github.com/MUKE-coder/pulse/issues/1" className="text-primary hover:underline">#1</a>{' '}
                    p50/p95/p99 + SLO alerts,{' '}
                    <a href="https://github.com/MUKE-coder/pulse/issues/2" className="text-primary hover:underline">#2</a>{' '}
                    N+1 detector,{' '}
                    <a href="https://github.com/MUKE-coder/pulse/issues/3" className="text-primary hover:underline">#3</a>{' '}
                    USE method dashboard,{' '}
                    <a href="https://github.com/MUKE-coder/pulse/issues/4" className="text-primary hover:underline">#4</a>{' '}
                    k6 timeline + flame graphs.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.22.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.22.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Performance hardening release.</strong> A senior-level audit of
                  every scaffold template found 27 issues — the 10 critical and high-impact
                  ones are fixed. Apps built with Grit should now show materially lower CPU
                  burn under sustained load.
                </p>

                <h3>Critical fixes</h3>
                <ul>
                  <li>
                    <strong>ActivityLogger middleware</strong> — was spawning a fresh goroutine
                    per request, each blocking on a row-level <code>FOR UPDATE</code> lock for
                    the audit hash chain. At 10k req/s the old design created 10k goroutines
                    all serializing on the same lock. Replaced with a bounded channel (4096) +
                    single writer goroutine. The single-writer design eliminates the lock
                    entirely (chain ordering is sequential by construction); the bounded
                    channel caps memory + goroutine count under traffic spikes. Drops on
                    overflow rather than OOM, with a new <code>AuditDroppedCount()</code>{' '}
                    helper for monitoring saturation.
                  </li>
                  <li>
                    <strong><code>audit.VerifyChain</code></strong> — was loading the entire{' '}
                    <code>activity_log</code> table into memory before scanning. At 1M rows
                    that's 250MB+ heap, instant OOM at 100M. Now walks in chunks of 1000 rows
                    with a cursor on <code>(created_at, id)</code>, honours{' '}
                    <code>context</code> cancellation, and the integrity endpoint passes a
                    60-second deadline so a runaway scan can't hold the connection forever.
                  </li>
                </ul>

                <h3>High-priority fixes</h3>
                <ul>
                  <li>
                    <strong><code>flags.Engine.evaluate</code></strong> — copies the flag
                    struct under <code>RLock</code> then releases <em>before</em> doing all
                    decision logic (date checks, allowlist scans, bucketing, JSON parsing).
                    Cuts lock-hold time from milliseconds to nanoseconds on the flag-check
                    hot path.
                  </li>
                  <li>
                    <strong>Cache middleware</strong> — SHA-256 cache keys swapped for FNV-1a.
                    ~50× faster on the hot path of every cacheable request, no correctness
                    loss (cache keys don't need cryptographic strength).{' '}
                    <code>responseCapture</code> switches <code>[]byte</code> append to{' '}
                    <code>bytes.Buffer</code> — 3 allocations instead of one per Write chunk.
                  </li>
                  <li>
                    <strong>Generated service queries</strong> — <code>Update</code> dropped
                    the redundant third <code>First()</code> after <code>Updates()</code>{' '}
                    (Updates mutates the loaded struct in place); <code>Delete</code> dropped
                    the preflight <code>First()</code> (GORM's Delete is atomic + RowsAffected
                    reveals existence). 2 queries saved per generated CRUD op.
                  </li>
                  <li>
                    <strong>Generated Export handler</strong> — was loading every matching
                    row with <code>Find(&items)</code>; now uses{' '}
                    <code>FindInBatches</code> in chunks of 1000. CSV exports stream directly
                    to the response writer (true streaming, constant memory). XLSX still
                    buffers because excelize has no streaming API, but the scan is chunked
                    so we don't hold the entire result set in one slice. New{' '}
                    <code>export.CSVRows()</code> helper for header-less subsequent batches.
                  </li>
                </ul>

                <h3>Medium fixes</h3>
                <ul>
                  <li>
                    <strong>Webhook Replay</strong> — <code>retry_count</code> increment is
                    now atomic via <code>gorm.Expr("retry_count + ?", 1)</code>. Two
                    concurrent replays of the same event no longer race to write the same
                    +1 result.
                  </li>
                  <li>
                    <strong>Flags <code>bucketFor</code> for anonymous users</strong> —{' '}
                    <code>crypto/rand.Read</code> instead of{' '}
                    <code>time.Now().UnixNano() % 100</code>. The old approach was biased
                    toward recent buckets under high QPS.
                  </li>
                </ul>

                <h3>Skill file: Performance &amp; Production Hygiene section</h3>
                <p>
                  The <code>.claude/skills/grit/SKILL.md</code> that{' '}
                  <code>grit new</code> generates now includes a dedicated{' '}
                  <strong>Performance &amp; Production Hygiene</strong> section. AI assistants
                  helping users build apps will see explicit hot-path rules, DB query rules,
                  background job rules, logging rules, and memory rules — including which
                  framework primitives are already audited (so they know not to reintroduce
                  the patterns this release just fixed). Examples:
                </p>
                <ul>
                  <li>
                    Never spawn unbounded goroutines per-request — use a buffered channel +
                    fixed worker pool (the ActivityLogger pattern).
                  </li>
                  <li>
                    Never hold a mutex across slow operations — read shared state, copy what
                    you need, release, then do the work (the flags.evaluate pattern).
                  </li>
                  <li>
                    Never load a whole table into memory — use{' '}
                    <code>paginate.List</code>, <code>FindInBatches</code>, or cursor walks
                    (the VerifyChain pattern).
                  </li>
                  <li>
                    Never use <code>time.Now().UnixNano() % N</code> for randomness — biased
                    by call frequency; use <code>crypto/rand</code>.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.21.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.21.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Feature flags + A/B testing baked into the framework (
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/46" target="_blank" rel="noopener noreferrer">#46</a>).
                  No LaunchDarkly bolt-on, no PostHog SaaS dependency — the engine,
                  the model, the admin endpoints, and the realtime push all ship in
                  every scaffolded API.
                </p>

                <h3>Usage</h3>
                <pre><code>{`if flags.IsEnabled(c, "new_dashboard") {
    // … render the new dashboard
}

switch flags.Variant(c, "checkout_redesign") {
case "control":   /* old flow */
case "variant_a": /* new flow */
case "variant_b": /* alternate new flow */
}`}</code></pre>

                <h3>Mechanics</h3>
                <ul>
                  <li>
                    <code>FeatureFlag.Rules</code> JSON holds <code>rollout_percentage</code>,{' '}
                    <code>allowlist_user_ids</code>, <code>blocklist_user_ids</code>,{' '}
                    <code>enabled_from</code>, <code>enabled_until</code>,{' '}
                    <code>variants</code>.
                  </li>
                  <li>
                    All flags load into an in-memory cache at boot. A background
                    goroutine refreshes every 30s; admin writes trigger an immediate
                    refresh. Flag checks never hit the DB.
                  </li>
                  <li>
                    Sticky bucketing: <code>SHA-256(user_id || ":" || flag_name) % 100</code>.
                    A user always lands in the same bucket for a given flag — no
                    flicker between sessions.
                  </li>
                  <li>
                    Allowlist always passes (bypasses the percentage roll). Blocklist
                    always denies. Both run before the percentage check.
                  </li>
                  <li>
                    A/B mode kicks in when <code>Rules.Variants</code> is non-empty.{' '}
                    <code>Variant()</code> returns the bucket-mapped variant string.
                    Sticky per (user, flag).
                  </li>
                </ul>

                <h3>Realtime updates</h3>
                <p>
                  When a flag is created / updated / deleted, the engine refreshes
                  its cache and broadcasts a <code>"flag.updated"</code> realtime
                  event over the v3.12 WebSocket hub. Frontend subscribers can
                  invalidate their cache and refetch — flag changes propagate in &lt;1s
                  across all connected clients.
                </p>

                <h3>Admin endpoints</h3>
                <ul>
                  <li>
                    <code>GET /api/admin/flags</code> — paginated list (searchable on
                    name + description, sortable on name / created_at / enabled).
                  </li>
                  <li>
                    <code>POST /api/admin/flags</code> — create. Name is unique +
                    immutable.
                  </li>
                  <li>
                    <code>PUT /api/admin/flags/:id</code> — update description / enabled
                    / rules. Bumps Version (the v3.14 optimistic-lock column).
                  </li>
                  <li>
                    <code>DELETE /api/admin/flags/:id</code> — remove + invalidate cache.
                  </li>
                  <li>
                    <code>GET /api/admin/flags/:id/exposures</code> — variant counts
                    for the rollout-health view: <code>{`[{ "variant": "enabled", "count": 4231 }, ...]`}</code>.
                  </li>
                </ul>

                <h3>Fail-closed semantics</h3>
                <ul>
                  <li>
                    Unknown flags return <code>false</code>. A typo in a flag name never
                    accidentally enables a feature.
                  </li>
                  <li>
                    Misconfigured <code>Rules</code> JSON also returns <code>false</code>{' '}
                    — the engine never panics on bad data.
                  </li>
                  <li>
                    Anonymous users (empty user_id) get a random bucket per request +
                    are not exposure-tracked. For sticky anonymous flags, pass a stable
                    identifier (session ID, device ID).
                  </li>
                </ul>

                <p className="text-sm text-muted-foreground">
                  Pairs with the v3.16 activity log + v3.19 hash chain — every flag
                  change is auditable, signed, and tamper-evident. SOC2-ish flag
                  governance for free.
                </p>
              </div>
            </div>

            {/* v3.20.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.20.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Webhook receiver framework (
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/57" target="_blank" rel="noopener noreferrer">#57</a>).
                  Wiring up Stripe / GitHub / WhatsApp / any HMAC-signed inbound webhook
                  is now {'<'}10 lines of app code. Signature verification, idempotency,
                  failed-handler replay — all framework concerns now.
                </p>

                <h3>The shape</h3>
                <pre><code>{`// In your app boot (e.g. internal/webhooks/handlers.go)
func init() {
    webhooks.Register("stripe", webhooks.Provider{
        SecretEnv: "STRIPE_WEBHOOK_SECRET",
        Verify:    webhooks.StripeVerifier,
        Extract:   webhooks.StripeExtractor,
    })

    webhooks.On("stripe", "invoice.paid", func(ctx context.Context, e *models.WebhookEvent) error {
        // … process the event
        return nil
    })
}`}</code></pre>
                <p>
                  The framework already mounted <code>POST /webhooks/:provider</code> in
                  routes — the path param picks the registered provider. No per-provider
                  routing code in your app.
                </p>

                <h3>Pipeline</h3>
                <ol>
                  <li>Route hits → look up provider (404 if unknown).</li>
                  <li>Read raw body + headers.</li>
                  <li>
                    <code>provider.Verify(secret, body, headers)</code> — 401 on signature
                    mismatch.
                  </li>
                  <li>
                    <code>provider.Extract(body, headers)</code> returns{' '}
                    <code>(eventType, externalID)</code>.
                  </li>
                  <li>
                    Insert into <code>webhook_events</code> — UNIQUE on{' '}
                    <code>(provider, external_id)</code> means duplicate deliveries
                    become <code>status=skipped</code> no-ops.
                  </li>
                  <li>
                    <code>webhooks.Dispatch(ctx, event)</code> runs the registered handler
                    for <code>(provider, eventType)</code>, falling back to a catch-all{' '}
                    <code>""</code> handler if no specific match.
                  </li>
                  <li>
                    Handler success → <code>status=processed</code>; handler error →{' '}
                    <code>status=failed</code> + <code>handler_error</code> recorded.
                    Provider always gets <code>200</code> once we persisted the event,
                    so retries don't hammer.
                  </li>
                </ol>

                <h3>Shipped verifiers</h3>
                <ul>
                  <li>
                    <code>HMACVerifier(header)</code> — generic hex HMAC-SHA256 in a
                    named header. Most simple partners use this.
                  </li>
                  <li>
                    <code>StripeVerifier</code> — Stripe's <code>t=...,v1=...</code>{' '}
                    scheme with 5-minute replay tolerance.
                  </li>
                  <li>
                    <code>GitHubVerifier</code> — GitHub's <code>X-Hub-Signature-256: sha256=...</code>{' '}
                    header.
                  </li>
                  <li>
                    Roll your own <code>VerifyFunc</code> for anything else — it's just{' '}
                    <code>func(secret string, body []byte, headers map[string]string) error</code>.
                  </li>
                </ul>

                <h3>Shipped extractors</h3>
                <ul>
                  <li>
                    <code>JSONFieldExtractor("type", "id")</code> — pulls top-level fields
                    from the JSON body. The most common shape (Stripe-style envelopes).
                  </li>
                  <li>
                    <code>GitHubExtractor</code> — reads <code>X-GitHub-Event</code> +{' '}
                    <code>X-GitHub-Delivery</code> headers.
                  </li>
                </ul>

                <h3>Admin endpoints</h3>
                <ul>
                  <li>
                    <code>GET /api/admin/webhooks?provider=stripe&status=failed</code> —
                    paginated list with the standard envelope. Filters: provider, status.
                  </li>
                  <li>
                    <code>POST /api/admin/webhooks/:id/replay</code> — re-runs the
                    handler for an existing event. Increments <code>retry_count</code>{' '}
                    and records the new outcome. Use this after a deploy fixes a
                    handler bug.
                  </li>
                </ul>

                <p className="text-sm text-muted-foreground">
                  Pairs naturally with the v3.10 idempotency middleware — both are
                  &quot;safe replay&quot; primitives, just on different sides of the
                  network. Outbound retries reuse <code>Idempotency-Key</code>; inbound
                  duplicates dedupe on <code>(provider, external_id)</code>.
                </p>
              </div>
            </div>

            {/* v3.19.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.19.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Tamper-evident audit log via append-only hash chain (
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/48" target="_blank" rel="noopener noreferrer">#48</a>).
                  Builds on the v3.16 <code>ActivityLog</code> — every row now carries{' '}
                  <code>PrevHash</code> + <code>Hash</code> columns where{' '}
                  <code>Hash = SHA-256(PrevHash || canonical(row))</code>. Mutating any
                  row breaks the chain on the next verification pass.
                </p>

                <h3>The chain</h3>
                <ul>
                  <li>
                    <strong>Genesis row</strong> has <code>PrevHash = ""</code>; every
                    subsequent row references the previous row's <code>Hash</code>.
                  </li>
                  <li>
                    Hash input is the <em>stable canonical form</em> of the
                    audit-relevant fields (user_id, method, path, status, payload digest,
                    IP, UA, duration, created_at unix-nano). ID + PrevHash + Hash
                    themselves are <em>not</em> in the canonical form — they're either
                    random (ID) or derived (Hash, PrevHash).
                  </li>
                  <li>
                    Insert uses <code>FOR UPDATE</code> lock on the latest row inside
                    the same transaction, so concurrent writes serialize cleanly without
                    forking the chain.
                  </li>
                </ul>

                <h3>The package</h3>
                <p>
                  New <code>internal/audit</code> ships these:
                </p>
                <ul>
                  <li>
                    <code>audit.Canonical(entry)</code> — stable JSON bytes for hashing.
                  </li>
                  <li>
                    <code>audit.ComputeHash(prevHash, canonical)</code> — runs SHA-256
                    over <code>prevHash || canonical</code>; returns hex.
                  </li>
                  <li>
                    <code>audit.AppendChained(db, entry)</code> — atomic insert with
                    chain lock. The middleware uses this; you can call it from anywhere.
                  </li>
                  <li>
                    <code>audit.VerifyChain(db)</code> — walks every row in{' '}
                    <code>(created_at, id)</code> order and recomputes hashes. Returns{' '}
                    <code>ChainStatus</code> with the first mismatch (broken_at_id +
                    expected vs got + message).
                  </li>
                </ul>

                <h3>The endpoint</h3>
                <pre><code>{`GET /api/admin/activity/integrity

→ { "valid": true, "total_entries": 12345 }

→ { "valid": false, "broken_at": 47, "broken_at_id": "uuid",
    "expected": "abc123...", "got": "def456...",
    "message": "hash mismatch — row was modified, deleted, or inserted out of order" }`}</code></pre>
                <p>
                  Wire this to a nightly cron + alerting webhook for free SOC2-ish audit
                  monitoring. Run it on-demand from a settings page when staff need the
                  current chain state.
                </p>

                <h3>What this defends against</h3>
                <ul>
                  <li>
                    Direct SQL <code>UPDATE</code> / <code>DELETE</code> on{' '}
                    <code>activity_logs</code> — the most common attack vector
                    (DBA covering tracks).
                  </li>
                  <li>
                    Out-of-band insertion of forged history.
                  </li>
                </ul>

                <h3>What it does NOT defend against</h3>
                <ul>
                  <li>
                    Compromise of the running server itself — an attacker with code
                    execution can rewrite the entire chain.
                  </li>
                  <li>
                    External anchoring (publishing the daily root hash to a public
                    ledger like a tweet, a transaction, or a Sigstore log) is the
                    follow-up — flagged in #48 as bonus material, not shipped here.
                  </li>
                </ul>

                <p className="text-sm text-muted-foreground">
                  Verification cost: O(n) — about 2–3 seconds per million rows on a warm
                  cache. The middleware insert is still fire-and-forget so audit DB
                  latency never blocks the response path; chain failures log instead of
                  cascading.
                </p>
              </div>
            </div>

            {/* v3.18.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.18.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  PDF generation module (
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/13" target="_blank" rel="noopener noreferrer">#13</a>).
                  Every scaffolded API ships <code>internal/pdf/</code> with Grit-styled
                  section helpers + a worked <code>RenderInvoice</code> template. Pure Go,
                  no Chromium / wkhtmltopdf native dependencies.
                </p>

                <h3>The Doc primitives</h3>
                <p>
                  <code>pdf.New()</code> returns a <code>*Doc</code> preconfigured with
                  Helvetica + 20mm margins + A4 portrait + Grit blue accent. Embeds the
                  underlying <code>*fpdf.Fpdf</code> so the full library is available when
                  helpers don't fit.
                </p>
                <ul>
                  <li>
                    <code>Header(title, subtitle)</code> — accent-colored 22pt title +
                    muted-gray subtitle line.
                  </li>
                  <li>
                    <code>KV(label, value)</code> + <code>TwoColumnKV(...)</code> — small-caps
                    label + body value pairs.
                  </li>
                  <li>
                    <code>Table(headers, rows, widths, aligns)</code> — light gray header
                    row, plain data rows, configurable widths + alignment per column.
                  </li>
                  <li>
                    <code>Totals([]TotalLine)</code> — right-aligned totals stack; the bold
                    line gets accent coloring + slightly larger size for the grand total.
                  </li>
                  <li>
                    <code>Notes(text)</code> — labeled multiline section, skipped when empty.
                  </li>
                  <li>
                    <code>Footer(text)</code> — centered italic 25mm above the page bottom.
                  </li>
                  <li>
                    <code>d.Bytes()</code> finalizes and returns the PDF byte slice ready
                    to stream to <code>c.Data(200, "application/pdf", b)</code>.
                  </li>
                </ul>

                <h3>RenderInvoice — worked example</h3>
                <pre><code>{`pdf.RenderInvoice(pdf.Invoice{
    Number:    "INV-202605-0001",
    IssueDate: time.Now(),
    DueDate:   time.Now().Add(14 * 24 * time.Hour),
    BillTo:    pdf.Party{Name: "Abu Seal", Contact: "abu@example.com"},
    Items: []pdf.LineItem{
        {Description: "Office rent — June", Quantity: 1, UnitPrice: 1500000, Total: 1500000},
        {Description: "Service charge",      Quantity: 1, UnitPrice:  120000, Total:  120000},
    },
    Subtotal: 1620000, Total: 1620000,
    Currency: "UGX",
    Notes:    "Pay by mobile money: +256...",
})`}</code></pre>
                <p>
                  Returns <code>([]byte, error)</code> — wire it to a handler:
                </p>
                <pre><code>{`func (h *InvoiceHandler) PDF(c *gin.Context) {
    inv, _ := h.Service.GetByID(c.Param("id"))
    bytes, err := pdf.RenderInvoice(toInvoice(inv))
    if err != nil { respond.Internal(c, err); return }
    c.Header("Content-Disposition", \`attachment; filename="\` + inv.Number + \`.pdf"\`)
    c.Data(200, "application/pdf", bytes)
}`}</code></pre>
                <p className="text-sm text-muted-foreground">
                  Copy <code>invoice.go</code> as a starting point for receipts, leases,
                  statements, quotes — the same primitives compose all of them. Add
                  <code> github.com/go-pdf/fpdf v0.9.0</code> dependency lands automatically
                  in scaffolded <code>go.mod</code>.
                </p>
              </div>
            </div>

            {/* v3.17.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.17.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Quality-of-life bundle. Four GitHub issues closed:
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/12" target="_blank" rel="noopener noreferrer"> #12</a>,{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/31" target="_blank" rel="noopener noreferrer">#31</a>,{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/35" target="_blank" rel="noopener noreferrer">#35</a>,{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/43" target="_blank" rel="noopener noreferrer">#43</a>.
                </p>

                <h3><code>grit init</code> — #35</h3>
                <p>
                  New CLI command writes <code>CLAUDE.md</code> + <code>AGENTS.md</code>{' '}
                  to the current directory. Both files carry the framework's hard rules
                  (Forms / Frontend stdlib / Data / Backend / Resources / Sync / Auth)
                  so contributors and AI assistants get the conventions right on first
                  PR. Skips existing files unless <code>--force</code> is passed; re-run
                  with <code>--force</code> after a major framework upgrade to refresh.
                </p>

                <h3>Verbose AutoMigrate — #31</h3>
                <p>
                  <code>Migrate()</code> now snapshots <code>ColumnTypes</code> before
                  and after each <code>AutoMigrate</code> call and logs a diff:
                </p>
                <pre><code>{`================================================================
DATABASE MIGRATION — 8 model(s) registered
================================================================
  + created models.Building
  ~ models.User — added 2 column(s): is_vip, vip_notes
----------------------------------------------------------------
Migration done — 1 created, 1 altered (+2 column), 6 unchanged.
================================================================`}</code></pre>
                <p>
                  Silent migrations are gone. Also fixes a pre-existing bug where{' '}
                  <code>Migrate</code> <em>skipped</em> already-existing tables — so
                  columns added to a model never actually landed in the DB. Now they do.
                </p>

                <h3>Cursor-based pagination — #43</h3>
                <ul>
                  <li>
                    <code>paginate.List</code> gains opt-in cursor mode via{' '}
                    <code>Config.CursorMode: true</code>. Response carries{' '}
                    <code>Meta.NextCursor</code> + <code>Meta.HasMore</code>{' '}
                    instead of <code>Page</code>/<code>Pages</code>.
                  </li>
                  <li>
                    Detects <code>HasMore</code> by fetching <code>PageSize + 1</code>{' '}
                    rows — no separate count query needed.
                  </li>
                  <li>
                    Cursor is opaque base64 of <code>(sort_value, id)</code> so pages
                    stay stable when rows insert mid-pagination. Works with any sort
                    field; extracts the value via reflection on the last row.
                  </li>
                  <li>
                    Total count opt-in via <code>Config.IncludeTotal</code> — costs an
                    extra <code>COUNT(*)</code>, leave off unless your UI shows a
                    &quot;X of Y&quot; indicator.
                  </li>
                  <li>
                    Offset mode stays the default for back-compat; new resources can
                    flip the flag.
                  </li>
                </ul>

                <h3>Generator quality — #12</h3>
                <p>The remaining tag-default heuristics from issue #12:</p>
                <ul>
                  <li>
                    <strong>URL fields</strong> (suffix <code>_url</code> + named{' '}
                    <code>url</code> / <code>image</code> / <code>avatar</code> /{' '}
                    <code>thumbnail</code> / <code>logo</code> / <code>cover</code> /{' '}
                    <code>icon</code> / <code>banner</code> / <code>photo</code>) get{' '}
                    <code>size:500</code> instead of <code>size:255</code>. UTM-tagged
                    links and signed S3 URLs blow past 255 in the wild.
                  </li>
                  <li>
                    <strong>Long-text fields</strong> named <code>description</code> /{' '}
                    <code>notes</code> / <code>content</code> / <code>body</code> /{' '}
                    <code>summary</code> / <code>bio</code> / <code>details</code> /{' '}
                    <code>comment</code> / <code>comments</code> / <code>message</code>{' '}
                    get <code>type:text</code>.
                  </li>
                  <li>
                    <strong>Money fields</strong> on <code>float</code> type (suffix{' '}
                    <code>_amount</code> / <code>_price</code> / <code>_total</code> /{' '}
                    <code>_cost</code> / <code>_fee</code> / <code>_balance</code> /{' '}
                    <code>_rent</code> / <code>_salary</code> / <code>_wage</code> /{' '}
                    <code>_value</code> / <code>_revenue</code> / <code>_deposit</code>{' '}
                    + named <code>amount</code> / <code>price</code> / <code>total</code> /{' '}
                    <code>cost</code> / <code>fee</code> / <code>balance</code> /{' '}
                    <code>subtotal</code>) get <code>type:decimal(12,2)</code> for
                    fixed-precision storage. No more <code>1.99 + 0.01 = 1.9999999</code>.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.16.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.16.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Three coherent admin-operations features at once:
                  CSV/Excel export per resource (
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/15" target="_blank" rel="noopener noreferrer">#15</a>),
                  activity audit log middleware (
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/32" target="_blank" rel="noopener noreferrer">#32</a>),
                  and the <code>apiErrorMessage</code> frontend helper (
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/27" target="_blank" rel="noopener noreferrer">#27</a>).
                </p>

                <h3>CSV / Excel export per resource — #15</h3>
                <ul>
                  <li>
                    New <code>internal/export</code> package: <code>CSV(w, items, opts)</code>{' '}
                    and <code>XLSX(w, items, opts)</code> with a typed{' '}
                    <code>Column{'{Header, Field, Format}'}</code> config. Field uses
                    dot-notation for associations (<code>"Tenant.Name"</code>).
                  </li>
                  <li>
                    Format strings: <code>"currency:UGX"</code>, <code>"date:2006-01-02"</code>,{' '}
                    <code>"datetime"</code>, <code>"bool"</code>. Empty string falls back to{' '}
                    <code>fmt.Sprintf("%v")</code>.
                  </li>
                  <li>
                    Resource generator now emits an <code>Export(c *gin.Context)</code>{' '}
                    handler method on every new resource, with columns derived from the
                    field list. Routes inject <code>GET /api/{'<plural>'}/export</code>{' '}
                    automatically.
                  </li>
                  <li>
                    Honours the same <code>search</code> param as List, so users can
                    export a filtered subset.
                  </li>
                  <li>
                    Adds <code>github.com/xuri/excelize/v2 v2.8.1</code> to scaffolded{' '}
                    <code>go.mod</code>.
                  </li>
                </ul>

                <h3>Activity audit log — #32</h3>
                <ul>
                  <li>
                    New <code>models.ActivityLog</code> with user_id + method + path +
                    status + payload digest (sha256, not raw body) + IP + user-agent +
                    duration. UUID PK; <code>created_at</code> indexed for time-range queries.
                  </li>
                  <li>
                    New <code>middleware.ActivityLogger(db)</code> mounted on every
                    protected mutation route. Skips safe methods + non-2xx responses +
                    unauthenticated requests.
                  </li>
                  <li>
                    Insert is fire-and-forget (goroutine). Audit DB latency never
                    blocks the response path; if the DB is down the entry drops
                    rather than failing the request.
                  </li>
                  <li>
                    New endpoint <code>GET /api/admin/activity</code> (admin-only) with{' '}
                    <code>paginate.List</code> filtering by <code>user_id</code>,{' '}
                    <code>method</code>, and <code>path</code> prefix. Drop in any audit-log UI.
                  </li>
                </ul>

                <h3>apiErrorMessage helper — #27</h3>
                <ul>
                  <li>
                    Three helpers in <code>packages/shared/types/api.ts</code>:{' '}
                    <code>apiErrorMessage(err, fallback?)</code>,{' '}
                    <code>apiErrorCode(err)</code>,{' '}
                    <code>apiErrorFields(err)</code>.
                  </li>
                  <li>
                    Walks the standard envelope chain
                    (<code>response.data.error.message</code>) plus axios{' '}
                    <code>err.message</code> plus a fallback so{' '}
                    <code>toast.error(apiErrorMessage(err))</code> is always meaningful.
                  </li>
                  <li>
                    <code>apiErrorCode</code> returns the envelope's <code>code</code>{' '}
                    string (<code>"VALIDATION_ERROR"</code>,{' '}
                    <code>"VERSION_CONFLICT"</code>, etc.) for branching logic.
                  </li>
                  <li>
                    <code>apiErrorFields</code> surfaces per-field validation details so
                    forms can highlight specific inputs.
                  </li>
                  <li>
                    New <code>internal/respond</code> package on the server side too:{' '}
                    <code>respond.NotFound / Validation / Forbidden / Conflict /
                    Internal</code> for handlers, replacing ad-hoc inline{' '}
                    <code>c.JSON(500, gin.H{'{...}'})</code>.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.15.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.15.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Frontend stdlib + form primitives. Closes seven GitHub issues at once
                  (<a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/19" target="_blank" rel="noopener noreferrer">#19</a>,{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/20" target="_blank" rel="noopener noreferrer">#20</a>,{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/21" target="_blank" rel="noopener noreferrer">#21</a>,{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/22" target="_blank" rel="noopener noreferrer">#22</a>,{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/23" target="_blank" rel="noopener noreferrer">#23</a>,{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/33" target="_blank" rel="noopener noreferrer">#33</a>,{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/34" target="_blank" rel="noopener noreferrer">#34</a>).
                  Every primitive lifted from real Grit-built business apps.
                </p>

                <h3>Format helpers (<code>lib/format.ts</code>) — #33</h3>
                <ul>
                  <li>
                    <code>formatCurrency(amount, currency?)</code> — locale-aware,
                    no-decimal mode for UGX / JPY / KRW / RWF / TZS / VND.
                  </li>
                  <li>
                    <code>formatDate(value, fmt?)</code> — token formatter (yyyy / MMMM /
                    MMM / MM / dd / HH / mm / ss). Default <code>"MMM d, yyyy"</code>.
                  </li>
                  <li>
                    <code>formatDateTime(value)</code> — &quot;May 2, 2026 · 2:30 PM&quot;.
                  </li>
                  <li>
                    <code>humanize("checked_in")</code> → &quot;Checked in&quot;.
                  </li>
                  <li>
                    <code>initials("Abu Seal")</code> → &quot;AS&quot;.
                  </li>
                  <li>
                    <code>setFormatConfig({"{ locale, currency }"})</code> at boot to
                    override.
                  </li>
                </ul>

                <h3><code>{'<CurrencyField>'}</code> — #19</h3>
                <p>
                  Live comma formatting as the user types
                  (&quot;3000&quot; → &quot;3,000&quot;), paste-friendly
                  (&quot;$1,234.56&quot; works), emits raw <code>number</code> to{' '}
                  <code>onChange</code>. Optional <code>prefix</code> slot for currency
                  code. Auto-toggles between formatted display (blur) and raw digits
                  (focus) so editing isn't a fight.
                </p>

                <h3><code>{'<SearchableSelect>'}</code> — #20</h3>
                <p>
                  Combobox with typeahead, ↑/↓/Enter/Esc keyboard nav, portaled
                  dropdown (escapes overflow:hidden ancestors), optional clear button.
                  Replaces native <code>{'<select>'}</code> for FK fields and any enum
                  with &gt; 5 values.
                </p>

                <h3><code>{'<DateField>'}</code> + <code>{'<DateRangeFilter>'}</code> — #21</h3>
                <ul>
                  <li>
                    <code>{'<DateField>'}</code> wraps native{' '}
                    <code>{'<input type="date">'}</code> with the standard label/hint/
                    error chrome — picked the native one for a11y, RTL, and i18n.
                  </li>
                  <li>
                    <code>{'<DateRangeFilter>'}</code> — preset chip bar (Today / Last 7 /
                    Last 30 / This month / Last month / Last 90 / This year / All) +
                    custom-range fallback.
                  </li>
                  <li>
                    <code>presetRange("last90")</code> helper exposed for non-UI uses.
                  </li>
                </ul>

                <h3><code>{'<Drawer>'}</code> — #22</h3>
                <p>
                  Right-edge slide-in panel. Closes on Esc + backdrop + X button.
                  Configurable widths (<code>sm</code>/<code>md</code>/<code>lg</code>/<code>xl</code>).
                  Optional sticky footer slot for the typical Cancel/Save row. Pair
                  with <code>{'<FormGrid>'}</code> + <code>{'<FormActions>'}</code>{' '}
                  from v3.11 for the standard create/edit experience.
                </p>

                <h3><code>{'<StatusBadge>'}</code> — #34</h3>
                <p>
                  Status string → coloured pill. Default map covers
                  paid / active / completed / pending / overdue / cancelled / draft / archived /
                  checked_in / in_progress and friends. Override or extend per app:
                </p>
                <pre><code>{`setStatusVariants({
  shipped: "info",
  on_hold: "warning",
});`}</code></pre>

                <h3><code>{'<AppShell>'}</code> + grouped sidebar — #23</h3>
                <ul>
                  <li>
                    New <code>lib/nav-config.ts</code> — single source of truth for
                    sidebar sections. Adding a new section is a one-line config edit.
                  </li>
                  <li>
                    <code>components/layout/sidebar.tsx</code> rewritten as a
                    config-driven grouped sidebar (section title + items with icons +
                    optional badges).
                  </li>
                  <li>
                    New <code>components/layout/app-shell.tsx</code> — bundles
                    TitleBar + Sidebar + Topbar + scrollable content + Cmd/Ctrl-K
                    command palette in one component. Wrap your dashboard{' '}
                    <code>Outlet</code> with this.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.14.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.14.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Offline-first foundation.</strong> Git-style sync model — work
                  locally, click Sync explicitly, resolve conflicts per-field, push
                  one-by-one. Every scaffolded API now has Version-tracked rows + the
                  POST /api/sync/push and GET /api/sync/pull endpoints; every desktop
                  scaffold ships a local SQLite mirror, an outbox with squash semantics,
                  and a title-bar Sync button + conflict-resolution dialog.
                </p>

                <h3>Server: versioning + sync endpoints</h3>
                <ul>
                  <li>
                    <code>Version int</code> column added to User, Upload, Blog. A{' '}
                    <code>BeforeUpdate</code> GORM hook auto-increments on every
                    server-side write. The resource generator emits both on every new
                    model.
                  </li>
                  <li>
                    <code>POST /api/sync/push</code> accepts a batch of changes; each
                    entry includes the version the client believes the server has. On
                    mismatch the response contains <code>VERSION_CONFLICT</code> + the
                    current server state, so the client can drive a merge UI.
                  </li>
                  <li>
                    <code>GET /api/sync/pull?model=X&since=cursor</code> returns every
                    row in the table updated after the cursor, paginated, with a new
                    cursor in the response.
                  </li>
                  <li>
                    New <code>internal/sync/registry.go</code> maps logical table names
                    (e.g. <code>"buildings"</code>) to <code>reflect.Type</code> so the
                    handler decodes dynamic payloads. New resources auto-register via{' '}
                    <code>// grit:sync</code> marker.
                  </li>
                </ul>

                <h3>Desktop: sync engine</h3>
                <ul>
                  <li>
                    New <code>apps/desktop/sync/</code> Go package. Opens a local SQLite
                    file under the OS user-config dir on app boot.
                  </li>
                  <li>
                    Three tables: <code>sync_records</code> (local mirror — reads come
                    from here), <code>sync_outbox</code> (pending changes; UNIQUE on
                    (model, entity_id) for squash), <code>sync_cursors</code>{' '}
                    (incremental pull positions).
                  </li>
                  <li>
                    Squash semantics: edit a record three times offline → one outbox
                    entry with the final state. delete-after-create cancels both
                    locally without ever hitting the network.
                  </li>
                  <li>
                    <code>Engine.Sync()</code> runs Pull then Push. Push posts the whole
                    outbox in one HTTP call; the response drives per-entry result
                    handling — successes clear from the outbox, conflicts get the
                    server state stashed for the merge UI.
                  </li>
                </ul>

                <h3>Wails bindings</h3>
                <p>The frontend talks to the engine through these Wails-bound methods on App:</p>
                <ul>
                  <li>
                    <code>LocalCreate</code> / <code>LocalUpdate</code> / <code>LocalDelete</code> —
                    write-through to local SQLite + outbox.
                  </li>
                  <li>
                    <code>LocalGet</code> / <code>LocalList</code> — read from the local mirror.
                  </li>
                  <li>
                    <code>Sync(tables)</code> — pull listed tables then push the outbox.
                    Returns counts.
                  </li>
                  <li>
                    <code>PendingCount</code>, <code>GetPendingChanges</code> — drive the
                    title-bar badge and the review panel.
                  </li>
                  <li>
                    <code>ResolveConflict(table, entityID, mergedData, serverVersion)</code> —
                    accepts the user's merge for a conflicted entry.
                  </li>
                </ul>

                <h3>UI</h3>
                <ul>
                  <li>
                    <strong>Title-bar Sync button</strong> with a pending-count badge.
                    Green refresh icon when clean; amber alert + count when there are
                    pending changes.
                  </li>
                  <li>
                    <strong>PendingChangesPanel</strong> — right-edge drawer listing
                    every outbox entry, split into &quot;Needs review&quot; (conflicts)
                    and &quot;Ready to push&quot;. <code>Sync now</code> button at the bottom.
                  </li>
                  <li>
                    <strong>ConflictDialog</strong> — field-level merge UI. Three columns
                    (Field / Local / Server v_N), per-field click to choose. Apply
                    builds the merged record and calls <code>ResolveConflict</code>.
                  </li>
                </ul>

                <h3>React hooks</h3>
                <ul>
                  <li>
                    <code>usePendingCount()</code> — polls every 2s for the badge.
                  </li>
                  <li>
                    <code>usePendingChanges()</code> — full outbox + refresh function.
                  </li>
                  <li>
                    <code>useSyncMutation(tables)</code> — kicks off a Sync, exposes
                    running/result/error state.
                  </li>
                  <li>
                    <code>useResolveConflict()</code> — applies one merge and refreshes.
                  </li>
                </ul>

                <h3>Wire format</h3>
                <pre><code>{`POST /api/sync/push
{ "changes": [
    { "op": "create", "model": "buildings", "id": "uuid", "version": 0, "data": {...} },
    { "op": "update", "model": "tenants",   "id": "uuid", "version": 5, "data": {...} },
    { "op": "delete", "model": "leases",    "id": "uuid", "version": 3 }
] }

→ { "results": [
    { "ok": true, "new_version": 1 },
    { "ok": false, "code": "VERSION_CONFLICT", "server_version": 7, "server_data": {...} },
    { "ok": true }
] }`}</code></pre>

                <p className="text-sm text-muted-foreground">
                  Deferred to v3.14.1: React Query offline-aware data hooks
                  (<code>useOfflineList</code>, <code>useOfflineGet</code>,{' '}
                  <code>useOfflineMutation</code>) and the resource generator emitting
                  offline-aware frontend hooks. The engine and primitives ship now;
                  ergonomics layer next.
                </p>
              </div>
            </div>

            {/* v3.13.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.13.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  New <code>grit generate sequence</code> command produces atomic, gap-free
                  sequential numbers like <code>INV-202605-0001</code>. Pattern lifted from a
                  real Grit-built rental management app — invoice / receipt / order numbering
                  is now a one-liner.
                </p>

                <h3>What it generates</h3>
                <ul>
                  <li>
                    <strong>First invocation only:</strong>{' '}
                    <code>internal/sequence/sequence.go</code> — a generic counter package
                    with <code>Counter</code> (the GORM-backed row),{' '}
                    <code>Config</code> (name + prefix + reset + width), and an atomic{' '}
                    <code>Next(db, cfg, t)</code> helper.
                  </li>
                  <li>
                    <strong>Every invocation:</strong>{' '}
                    <code>internal/services/&lt;name&gt;_sequence.go</code> — a typed
                    convenience wrapper. Handlers call e.g.{' '}
                    <code>services.NextInvoiceNumber(h.DB, time.Now())</code> without
                    knowing the prefix or reset cadence.
                  </li>
                  <li>
                    Auto-injects <code>&sequence.Counter{'{}'}</code> into the{' '}
                    <code>Models()</code> migration slice (idempotent).
                  </li>
                </ul>

                <h3>Mechanics</h3>
                <ul>
                  <li>
                    Counter rows keyed by <code>(name, bucket)</code> where bucket is{' '}
                    <code>"YYYYMM"</code> for monthly resets, <code>"YYYY"</code> for yearly,
                    or empty for never. So a monthly counter automatically restarts at 1
                    on the first call of each new month.
                  </li>
                  <li>
                    Atomic via row-level <code>SELECT FOR UPDATE</code> on Postgres
                    (concurrent callers serialize on the counter row). SQLite serializes
                    writes globally so it's also safe.
                  </li>
                </ul>

                <h3>Usage</h3>
                <pre><code>{`grit generate sequence Invoice
grit generate sequence Order --prefix ORD --reset yearly --width 6
grit generate sequence Receipt --reset never`}</code></pre>

                <p>Flags:</p>
                <ul>
                  <li>
                    <code>--prefix</code> — alphabetic prefix (default: first 3 chars of
                    the name, uppercased)
                  </li>
                  <li>
                    <code>--reset</code> — when the counter resets:{' '}
                    <code>monthly</code> (default), <code>yearly</code>, or <code>never</code>
                  </li>
                  <li>
                    <code>--width</code> — zero-padded width of the numeric portion
                    (default 4)
                  </li>
                </ul>

                <p className="text-sm text-muted-foreground">
                  The <code>grit generate report</code> generator (Recharts tabs page +
                  Go ReportService) is deferred to a future release — it needs more design
                  work for the React chart layer than fits a same-day release.
                </p>
              </div>
            </div>

            {/* v3.12.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.12.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Realtime WebSocket hub baked into every API + a desktop client + hooks
                  for subscribing. And a sweep of every remaining numeric ID — UUIDs are
                  now the canonical ID type everywhere in the framework.
                </p>

                <h3>Realtime hub (API)</h3>
                <ul>
                  <li>
                    New package: <code>internal/realtime/hub.go</code>. One <code>Hub</code> per
                    process; each user can have multiple connections (desktop + mobile + web).
                  </li>
                  <li>
                    <code>Hub.SendToUser(userID, evt)</code>, <code>SendToUsers(ids, evt)</code>,
                    and <code>Broadcast(evt)</code> let any handler or service push events.
                  </li>
                  <li>
                    Slow-client safe: per-connection 32-message send buffer; when full, that
                    one client's message is dropped — never blocks the entire hub. Slow
                    clients resync on their next REST refetch.
                  </li>
                  <li>
                    New handler: <code>internal/handlers/realtime.go</code> upgrades the
                    request to a WebSocket and registers the client with the hub.
                  </li>
                  <li>
                    Mounted at <code>GET /api/ws?token=&lt;jwt&gt;</code> — query-string auth
                    because browsers can't set custom headers on the WS handshake.
                  </li>
                  <li>
                    Wire format: <code>{'{ type: "<topic>", payload: {...} }'}</code>. Suggested
                    topics: <code>chat.message.new</code>, <code>notification.new</code>,{' '}
                    <code>system.connected</code>, or your own{' '}
                    <code>{'resource.<name>.<verb>'}</code> namespace.
                  </li>
                  <li>
                    Dependency added: <code>github.com/gorilla/websocket v1.5.3</code>.
                  </li>
                </ul>

                <h3>Realtime client (desktop)</h3>
                <ul>
                  <li>
                    New file: <code>frontend/src/lib/realtime.ts</code>. Singleton client with
                    auto-reconnect via exponential backoff (1s, 2s, 4s, 8s, capped at 15s).
                  </li>
                  <li>
                    Global <code>realtimeBus</code> EventTarget — any component can subscribe.
                  </li>
                  <li>
                    Start from <code>AuthProvider</code> after tokens land, stop on logout.
                  </li>
                  <li>
                    New hook: <code>useRealtimeEvent&lt;T&gt;(type, callback)</code> subscribes to a
                    typed topic and unsubscribes on unmount. Plus{' '}
                    <code>useRealtimeAny()</code> for catch-all handlers (debug, toast bar).
                  </li>
                </ul>

                <h3>ID consistency sweep — UUIDs everywhere</h3>
                <p>
                  v3.9.1 standardized the User model on string UUID PKs but a long tail of
                  numeric IDs remained in the framework. v3.12.0 cleans them all up.
                </p>
                <ul>
                  <li>
                    <strong>Go scaffold:</strong> the prebuilt <code>Blog</code> model in{' '}
                    <code>api_blog_files.go</code> swaps from <code>gorm.Model</code> (auto-incr
                    uint) to a string UUID PK with a <code>BeforeCreate</code> hook. Service
                    signatures (<code>GetByID</code>, <code>Update</code>, <code>Delete</code>) and
                    handler param parsing all switch from <code>uint</code> to <code>string</code>.
                  </li>
                  <li>
                    <strong>Standalone desktop scaffold</strong> (<code>grit new-desktop</code>):
                    User, Blog, and Contact models all switch to string UUID PKs with{' '}
                    <code>BeforeCreate</code> hooks. All Wails-bound App methods and underlying
                    service signatures use <code>id string</code>. Frontend mutation typings
                    follow.
                  </li>
                  <li>
                    <strong>Shared TS types:</strong> <code>User</code>, <code>Upload</code>, and{' '}
                    <code>Blog</code> interfaces all use <code>id: string</code>. URL builders in{' '}
                    <code>API_ROUTES</code> take <code>id: string</code>. The <code>BlogSchema</code>{' '}
                    Zod schema uses <code>z.string()</code>.
                  </li>
                  <li>
                    <strong>Admin TS:</strong> <code>DataTable</code> selection state, generic{' '}
                    <code>useResourceItem</code> / <code>useUpdateResource</code> /{' '}
                    <code>useDeleteResource</code> mutation typings,{' '}
                    <code>RelationshipSelectField</code> single value,{' '}
                    <code>MultiRelationshipSelectField</code> array values, and <code>handleDelete</code>{' '}
                    callbacks all switch from <code>number</code> to <code>string</code>.
                  </li>
                </ul>

                <p className="text-sm text-muted-foreground">
                  Net effect: UUID is the canonical ID type across the entire framework. Any
                  resource generator output, any scaffolded type, any Go signature — all
                  string UUIDs. No more <code>id: number</code> hiding in some corner.
                </p>
              </div>
            </div>

            {/* v3.11.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.11.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Three new desktop primitive files ship with every <code>--desktop</code>{' '}
                  scaffold, lifted from a real Grit-built rental management app. They cover
                  the master-detail layout, form chrome, and filter chips that every CRUD
                  page reinvents — saving ~200 LOC per resource.
                </p>

                <h3><code>components/two-pane.tsx</code> — master-detail layout</h3>
                <ul>
                  <li>
                    <code>TwoPane</code> — outer flex container with overflow handling.
                  </li>
                  <li>
                    <code>ListPane</code> — fixed-width (352px) left pane with title +
                    count + new button + searchbar + optional filters slot + scrollable body
                    + optional footer. Toolbar slot for refresh buttons or other actions.
                  </li>
                  <li>
                    <code>ListRow</code> — icon/avatar + title + subtitle + right-side
                    meta. Selected state shows a 2px accent bar on the left edge.
                  </li>
                  <li>
                    <code>DetailPane</code> — right pane with optional header + scrollable
                    content. <code>empty=true</code> renders an <code>EmptyState</code>{' '}
                    with the configured title/hint instead.
                  </li>
                  <li>
                    <code>EmptyState</code>, <code>DetailSection</code> (small caps section
                    header), and <code>DetailField</code> (labelled value rows for read views).
                  </li>
                </ul>

                <h3><code>components/form.tsx</code> — form chrome</h3>
                <ul>
                  <li>
                    <code>TextField</code>, <code>TextAreaField</code>, <code>SelectField</code> —
                    forwarded refs, consistent label/hint/error layout, focus ring,
                    disabled styling. Plug straight into <code>react-hook-form</code>.
                  </li>
                  <li>
                    <code>FormGrid</code> — 1, 2, or 3 columns on <code>{'>'}=sm</code>; stacks
                    on small screens.
                  </li>
                  <li>
                    <code>FormSection</code> — small caps title + optional description
                    over a stack of fields.
                  </li>
                  <li>
                    <code>FormActions</code> — Cancel + Submit pair with{' '}
                    <code>isPending</code> support (button disables, label flips to
                    &quot;Saving...&quot;).
                  </li>
                </ul>

                <h3><code>components/filter-chip.tsx</code> — filter chips</h3>
                <ul>
                  <li>
                    <code>FilterChip</code> — toggleable pill, active state shows accent
                    background; optional <code>onClear</code> renders an X to clear a single
                    filter; optional <code>count</code> renders a small count badge.
                  </li>
                  <li>
                    <code>FilterBar</code> — horizontal scrollable wrapper. Drop into
                    <code>ListPane</code>&apos;s <code>filters</code> slot.
                  </li>
                </ul>

                <h3>Tailwind tokens</h3>
                <ul>
                  <li>
                    Added <code>listpane</code> spacing token (<code>22rem</code> / 352px)
                    to the desktop Tailwind config so <code>w-listpane</code> works.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.10.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.10.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Foundation release for upcoming offline-first work. Every scaffolded API
                  now ships with idempotent-retry semantics; every scaffolded client
                  auto-attaches an <code>Idempotency-Key</code> on mutations; and the
                  desktop scaffold gains a connection-status indicator backed by an
                  API heartbeat.
                </p>

                <h3>Idempotency middleware (API)</h3>
                <ul>
                  <li>
                    New file: <code>internal/middleware/idempotency.go</code> — wired into{' '}
                    <code>routes.Setup</code> as a global middleware.
                  </li>
                  <li>
                    Activates only when the request carries an <code>Idempotency-Key</code> header
                    and the method is <code>POST</code>/<code>PUT</code>/<code>PATCH</code>/<code>DELETE</code>.
                  </li>
                  <li>
                    First 2xx response is cached in Redis for <strong>24 hours</strong>, keyed by{' '}
                    <code>(method, path, key)</code>. Subsequent requests with the same key replay
                    the cached response instead of re-executing the handler.
                  </li>
                  <li>
                    Errors (4xx/5xx) are intentionally <strong>not</strong> cached — clients can
                    retry transient failures with the same key.
                  </li>
                  <li>
                    Sets <code>Idempotent-Replayed: true</code> response header on cache hits so
                    clients can distinguish replays from fresh executions.
                  </li>
                </ul>

                <h3>Client-side header injection</h3>
                <ul>
                  <li>
                    Desktop, Expo, web, and admin clients all auto-attach a UUIDv4{' '}
                    <code>Idempotency-Key</code> on unsafe methods via the request interceptor.
                  </li>
                  <li>
                    The 401-refresh path now reuses the same key when re-issuing a request after a
                    token refresh — so a token expiring mid-write can never double-create.
                  </li>
                </ul>

                <h3>Online-status hook (desktop)</h3>
                <ul>
                  <li>
                    New hook: <code>useOnlineStatus()</code> at{' '}
                    <code>frontend/src/hooks/use-online-status.ts</code>.
                  </li>
                  <li>
                    Combines <code>navigator.onLine</code> (cheap pre-check) with a 15-second
                    heartbeat to <code>/api/health</code> (the truth signal). Returns{' '}
                    <code>{'{ isOnline, lastCheckedAt }'}</code>.
                  </li>
                  <li>
                    Heartbeat times out after 5s so a sleeping laptop surfaces as offline
                    instantly on wake.
                  </li>
                  <li>
                    The title-bar gains a <code>ConnectionIndicator</code> — small green/amber
                    dot reflecting API reachability. Hover for last-checked timestamp.
                  </li>
                </ul>

                <p className="text-sm text-muted-foreground">
                  This is the foundation for the offline-first scaffold landing in a later
                  v3.x release — write-queues, optimistic updates, and last-write-wins
                  conflict resolution all need stable idempotency keys to be safe.
                </p>
              </div>
            </div>

            {/* v3.9.2 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.9.2
                </span>
                <span className="text-sm text-muted-foreground">April 25, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Every <code>grit generate resource</code> run now emits a List handler that is
                  ~15 lines instead of ~55. The page / sort / search boilerplate moved into a shared{' '}
                  <code>internal/paginate</code> package that ships with every scaffolded API —
                  one source of truth for clamping, whitelisting, and search. Addresses{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/14" target="_blank" rel="noopener noreferrer">issue #14</a>.
                </p>

                <h3>New <code>paginate</code> package</h3>
                <ul>
                  <li>
                    <code>paginate.List[T](query, paginate.Bind(c), paginate.Config{'{'}...{'}'})</code> —
                    typed, generic helper that runs search, sort, filter, and pagination against
                    any <code>*gorm.DB</code> query.
                  </li>
                  <li>
                    <code>paginate.Bind(c)</code> reads <code>page</code>, <code>page_size</code>,
                    {' '}<code>search</code>, <code>sort_by</code>, <code>sort_order</code> from the
                    Gin query, clamps <code>page</code> to ≥ 1 and <code>page_size</code> to
                    {' '}<code>[1, 100]</code>.
                  </li>
                  <li>
                    <code>paginate.Config</code> whitelists sortable columns and declares the
                    searchable column set — requests for columns outside the whitelist fall back
                    to <code>created_at desc</code>.
                  </li>
                  <li>
                    <code>paginate.Result[T]</code> returns the canonical{' '}
                    <code>{'{'} data, meta: {'{'} total, page, page_size, pages {'}'} {'}'}</code>{' '}
                    envelope — matches the existing API response format exactly.
                  </li>
                </ul>

                <h3>Generator update</h3>
                <ul>
                  <li>
                    The emitted List handler now delegates to <code>paginate.List</code>. Every
                    generated resource gets the same clamping, whitelisting, and UUID-safe search
                    behavior — no per-resource drift.
                  </li>
                  <li>
                    Searchable column selection uses <code>IsSearchable()</code> (text / string /
                    slug / richtext only), so FK UUID columns are no longer accidentally included
                    in <code>ILIKE</code> search — a leftover rough edge from{' '}
                    <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/12" target="_blank" rel="noopener noreferrer">issue #12</a>.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.9.1 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.9.1
                </span>
                <span className="text-sm text-muted-foreground">April 24, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Patch release fixing compilation and consistency bugs in v3.9.0. Every
                  freshly scaffolded project (including <code>--mobile --desktop</code>)
                  and every <code>grit generate resource</code> run now produces Go code that builds
                  cleanly on the first try. Thanks to <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/9" target="_blank" rel="noopener noreferrer">issue #9</a>,
                  {' '}<a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/10" target="_blank" rel="noopener noreferrer">#10</a>,
                  {' '}<a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/11" target="_blank" rel="noopener noreferrer">#11</a>,
                  {' '}and <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/12" target="_blank" rel="noopener noreferrer">#12</a>.
                </p>

                <h3>Scaffold fixes</h3>
                <ul>
                  <li>
                    <strong>Missing imports</strong>: added <code>&quot;log&quot;</code> to <code>config.go</code>,{' '}
                    <code>&quot;gorm.io/gorm/logger&quot;</code> to <code>user.go</code>,{' '}
                    <code>&quot;net/http&quot;</code> to <code>middleware/logger.go</code>.
                  </li>
                  <li>
                    <strong>Stray package prefix</strong>: removed <code>handlers.</code> qualifier on{' '}
                    <code>IsTrustedDevice</code> (same-package call).
                  </li>
                  <li>
                    <strong>User ID type consistency</strong>: normalized <code>UserID</code> and{' '}
                    <code>UploadID</code> to <code>string</code> UUIDs across 2FA models, auth service,
                    TOTP handler (<code>c.GetString(&quot;user_id&quot;)</code> replaces{' '}
                    <code>c.GetUint</code>), jobs package, and upload handler.
                  </li>
                </ul>

                <h3>Desktop scaffold fixes</h3>
                <ul>
                  <li>
                    <code>keychain.go</code> moved from <code>internal/</code> to the top level (the
                    subdirectory file was declaring <code>package main</code>, which Go rejects).
                  </li>
                  <li>
                    <code>go.mod</code> module path fixed from{' '}
                    <code>&lt;project&gt;/apps/api/apps/desktop</code> to{' '}
                    <code>&lt;project&gt;/apps/desktop</code>.
                  </li>
                </ul>

                <h3>Resource generator fixes</h3>
                <ul>
                  <li>
                    Service signatures take <code>id string</code> instead of <code>id uint</code> --
                    matches the UUID string PK the models have always emitted.
                  </li>
                  <li>
                    Handler FK fields, handler M2M arrays, TS interface FK fields, and TanStack hook
                    ID types all switched to <code>string</code> (were <code>uint</code> /{' '}
                    <code>number</code>).
                  </li>
                  <li>
                    <strong>Initialism-aware</strong>{' '}
                    <code>toPascalCase</code> / <code>toSnakeCase</code>:{' '}
                    <code>owner_id</code> → <code>OwnerID</code> (was <code>OwnerId</code>),{' '}
                    <code>image_url</code> → <code>ImageURL</code> (was <code>ImageUrl</code>),{' '}
                    <code>api_key</code> → <code>APIKey</code>. Round-trips correctly
                    (snake → pascal → snake).
                  </li>
                  <li>
                    <strong>Zod schemas</strong> now emit snake_case field names matching the Go
                    handler&apos;s JSON tags (previously emitted camelCase, causing validation and
                    <code>ShouldBindJSON</code> mismatches).
                  </li>
                  <li>
                    Zod FK and M2M validators use <code>z.string().uuid()</code> instead of{' '}
                    <code>z.number().int()</code>.
                  </li>
                  <li>
                    FK columns generate with <code>gorm:&quot;size:36;index&quot;</code> (matches UUID PK width).
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.9.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.9.0
                </span>
                <span className="text-sm text-muted-foreground">April 15, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>New: <code>--desktop</code> flag</h3>
                <ul>
                  <li>
                    <strong>Desktop + mobile + API in one monorepo</strong> &mdash;{' '}
                    <code>grit new myapp --mobile --desktop</code> scaffolds a complete multi-client
                    SaaS: Go API shared by an Expo mobile app AND a Wails desktop app. All three
                    share the same <code>packages/shared</code> types and schemas.
                  </li>
                  <li>
                    <strong>Wails as a thin client</strong> &mdash; The new desktop app is a
                    frameless Wails window that calls the shared API over HTTP. No embedded Go
                    business logic, no local SQLite. Wails bindings are used only for native OS
                    features: window controls, file dialogs, and OS keychain (macOS Keychain,
                    Windows Credential Manager, Linux Secret Service) for JWT storage.
                  </li>
                  <li>
                    <strong>Distinct from <code>grit new-desktop</code></strong> &mdash; The
                    standalone offline-first desktop scaffold (<code>grit new-desktop</code>) is
                    unchanged. <code>--desktop</code> is a new, separate capability for
                    always-online multi-client apps.
                  </li>
                </ul>

                <h3>Premium Desktop UX</h3>
                <ul>
                  <li>
                    <strong>Platform-aware window chrome</strong> &mdash; macOS traffic lights on
                    the left, Windows/Linux controls on the right. Detected at runtime via
                    <code>GetPlatform()</code> Wails binding.
                  </li>
                  <li>
                    <strong>Command palette</strong> (<code>⌘K</code>) &mdash; every scaffolded
                    desktop app ships with a Raycast/Linear-style command palette. Searchable
                    navigation + actions with keyboard-first UX.
                  </li>
                  <li>
                    <strong>Fixed 240px sidebar</strong> &mdash; not collapsible. Desktop windows
                    are wide enough; collapse toggles are a web pattern.
                  </li>
                  <li>
                    <strong>Global keyboard shortcuts</strong> &mdash;{' '}
                    <code>useShortcuts()</code> hook with defaults: <code>⌘K</code> palette,{' '}
                    <code>⌘,</code> settings, <code>⌘L</code> logout, <code>Esc</code> to close.
                  </li>
                  <li>
                    <strong>More negative space</strong> &mdash; content padding is <code>32px</code>
                    (vs web&apos;s <code>24px</code>) for long focus sessions. Subtler shadows
                    (OS chrome already provides elevation).
                  </li>
                </ul>

                <h3>Style Guide</h3>
                <ul>
                  <li>
                    New <strong>§14.5 Desktop App Patterns</strong> section in
                    <code>GRIT_STYLE_GUIDE.md</code> covering window chrome, sidebar (not
                    collapsible), topbar, command palette, keyboard shortcuts, OS keychain
                    integration, typography (tighter than web), and do&apos;s &amp; don&apos;ts
                    (no breadcrumbs, no header banners, no web-style autoplay).
                  </li>
                </ul>

                <h3>Usage</h3>
                <pre><code>{`grit new myapp --mobile --desktop --next
# apps/api + apps/web + apps/expo + apps/desktop

grit new myapp --desktop --triple
# apps/api + apps/web + apps/admin + apps/desktop

grit new myapp --api --desktop
# apps/api + apps/desktop (minimal)`}</code></pre>
              </div>
            </div>

            {/* v3.8.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.8.0
                </span>
                <span className="text-sm text-muted-foreground">April 11, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Design System</h3>
                <ul>
                  <li>
                    <strong>GRIT_STYLE_GUIDE.md</strong> &mdash; First official style guide for all Grit-scaffolded
                    projects. Premium Minimal aesthetic (Linear / Vercel school), Grit purple{' '}
                    <code>#6C5CE7</code> primary, Onest font. Covers typography, color palette, spacing,
                    shadows, every component spec (buttons, inputs, cards, tables, modals), auth page rules,
                    CLI scaffolding design, admin panel patterns, email templates.
                  </li>
                </ul>

                <h3>Admin Layout</h3>
                <ul>
                  <li>
                    <strong>Topbar refactor</strong> &mdash; Moved sidebar collapse toggle to top-left of the
                    topbar (next to mobile menu button). Moved theme toggle, notifications bell, and enhanced
                    user menu to the top-right cluster alongside search. The sidebar now contains only
                    navigation. Matches modern dashboard patterns (Linear, Vercel, Raycast).
                  </li>
                  <li>
                    <strong>Enhanced user menu</strong> &mdash; Dropdown now shows User Activity, Settings,
                    Billing, and Log out sections with a user name/email header.
                  </li>
                </ul>

                <h3>PageHeader Component</h3>
                <ul>
                  <li>
                    <strong>Consistent page headers</strong> &mdash; New <code>&lt;PageHeader /&gt;</code>{' '}
                    component at <code>components/layout/page-header.tsx</code> with title, description,
                    breadcrumbs, actions slot, and a 4-card stats grid. Every generated resource page
                    auto-includes it.
                  </li>
                  <li>
                    <strong>Auto-generated stats cards</strong> &mdash; Resource pages now ship with 4 default
                    stat cards (Total, This Week, This Month, Updated Recently) fetched from the API. Override
                    via <code>defineResource({`{ stats: { cards: [...] } }`})</code> or disable with{' '}
                    <code>stats: false</code>.
                  </li>
                </ul>

                <h3>Auth Pages</h3>
                <ul>
                  <li>
                    <strong>New centered auth variant</strong> &mdash; <code>grit new myapp --style centered</code>{' '}
                    scaffolds Linear-school single-card auth pages (login, sign-up, forgot-password). ~420px
                    card on a subtle radial gradient background. The original split-screen design remains the
                    default (unchanged).
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.7.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.7.0
                </span>
                <span className="text-sm text-muted-foreground">April 3, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Security</h3>
                <ul>
                  <li>
                    <strong>Security headers middleware</strong> &mdash; New <code>SecurityHeaders()</code>{' '}
                    middleware adds X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, Referrer-Policy,
                    Permissions-Policy, and HSTS (when HTTPS detected) on every response.
                  </li>
                  <li>
                    <strong>Max body size middleware</strong> &mdash; 10MB default limit returns 413 on exceed.
                  </li>
                  <li>
                    <strong>JWT secret validation</strong> &mdash; Warns if <code>JWT_SECRET</code> is shorter
                    than 32 characters.
                  </li>
                  <li>
                    <strong>Sentinel WAF</strong> &mdash; Now runs in <code>ModeBlock</code> in production{' '}
                    (was always <code>ModeLog</code>). Development keeps <code>ModeLog</code>.
                  </li>
                </ul>

                <h3>Performance</h3>
                <ul>
                  <li>
                    <strong>GORM AutoMigrate silence</strong> &mdash; Migration now uses a{' '}
                    <code>logger.Silent</code> session to suppress schema inspection SQL noise. Fixes{' '}
                    <a href="https://github.com/MUKE-coder/grit/issues/8" className="text-primary hover:underline"
                       target="_blank" rel="noopener noreferrer">issue #8</a>.
                  </li>
                </ul>

                <h3>Web App Auth</h3>
                <ul>
                  <li>
                    <strong>Auth pages for the web app</strong> &mdash; The web app (<code>apps/web</code>) now
                    ships with its own auth pages: login, register, forgot-password, OAuth callback.
                    Previously only the admin panel had auth. This is critical for e-commerce and SaaS where
                    end users log in on the web app, not the admin.
                  </li>
                  <li>
                    <strong><code>useAuth()</code> hook</strong> &mdash; React Query + js-cookie token
                    management with <code>AuthProvider</code> context wrapping the web app.
                  </li>
                </ul>

                <h3>Mobile (Expo)</h3>
                <ul>
                  <li>
                    <strong>Major Expo scaffold upgrade</strong> &mdash; 4 tabs (Home, Explore, Profile,
                    Settings) instead of 2. All forms use react-hook-form + zod. Home screen with stat cards
                    and pull-to-refresh. Explore screen with search and category discovery. Settings with
                    SectionList. Profile with display/edit mode.
                  </li>
                  <li>
                    <strong>OAuth in mobile</strong> &mdash; Google OAuth via <code>expo-web-browser</code>{' '}
                    with deep-link callback handling.
                  </li>
                  <li>
                    <strong>New Expo dependencies</strong> &mdash; react-hook-form, @hookform/resolvers, zod,
                    expo-image, expo-haptics, expo-web-browser. Splash screen config in{' '}
                    <code>app.json</code>.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.6.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.6.0
                </span>
                <span className="text-sm text-muted-foreground">March 27, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Scaffold into current directory</strong> &mdash; <code>grit new .</code> and{' '}
                    <code>grit new ./</code> now scaffold into the current directory instead of creating a
                    subfolder. Infers the project name from the folder name. Also auto-detects when the
                    current directory name matches the project name.
                  </li>
                  <li>
                    <strong><code>--force</code> flag</strong> &mdash; Allows scaffolding into non-empty
                    directories. Useful when a repo was cloned first (with README, .git, LICENSE) before
                    scaffolding: <code>grit new . --triple --vite --force</code>.
                  </li>
                  <li>
                    <strong><code>--here</code> flag</strong> &mdash; Explicit alternative to{' '}
                    <code>grit new .</code> for in-place scaffolding.
                  </li>
                  <li>
                    <strong>30 standalone courses</strong> &mdash; Added 20 new courses to the learning
                    platform (42 total across 3 tracks + 20 standalone). Topics include testing, GORM
                    mastery, WebSockets, Stripe payments, blog/CMS, CI/CD, middleware, and the 100-component
                    UI registry.
                  </li>
                </ul>

                <h3>Bug Fixes</h3>
                <ul>
                  <li>
                    <strong>Flags now skip interactive prompt</strong> &mdash; Running{' '}
                    <code>grit new myapp --triple --vite</code> no longer shows the architecture/frontend
                    selection prompt. Flags act as true shortcuts for non-interactive setup.
                  </li>
                  <li>
                    <strong>Module path upgrade to /v3</strong> &mdash; Fixed{' '}
                    <code>go install ...@latest</code> downloading v2.9.0 instead of v3.x.{' '}
                    All import paths updated from <code>grit/v2</code> to <code>grit/v3</code>.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.5.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v3.5.0
                </span>
                <span className="text-sm text-muted-foreground">March 26, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Documentation</h3>
                <ul>
                  <li>
                    <strong>Full docs redesign</strong> &mdash; Rebuilt the documentation site with a
                    Tailwind CSS-inspired aesthetic. New dark theme (<code>#0b1120</code>), sky-blue accents,
                    cleaner header with backdrop blur, redesigned code blocks with file tabs and line
                    highlighting, and new <code>StepWithCode</code> component for two-column step-by-step
                    guides (text left, code right).
                  </li>
                  <li>
                    <strong>Installation page redesigned</strong> &mdash; Step-numbered sections (01-04)
                    with the new two-column layout, system requirements table, architecture shortcuts,
                    and services grid.
                  </li>
                  <li>
                    <strong>Architecture Modes page</strong> &mdash; Visual cards for all 5 architectures
                    (single, double, triple, API only, mobile) with directory structure trees, features
                    list, ideal use cases, and frontend framework comparison.
                  </li>
                  <li>
                    <strong>TanStack Router guide</strong> &mdash; Complete guide for the TanStack Router
                    frontend option: project structure, routing patterns, comparison table with Next.js,
                    route examples, and admin panel auth guards.
                  </li>
                  <li>
                    <strong>New CLI Commands page</strong> &mdash; Documents <code>grit routes</code>,
                    <code>grit down/up</code> (maintenance mode), and <code>grit deploy</code>. Includes
                    complete command reference table for all 21 CLI commands.
                  </li>
                  <li>
                    <strong>Deploy Command guide</strong> &mdash; Step-by-step deployment pipeline with
                    systemd service unit and Caddyfile examples, flags table.
                  </li>
                </ul>

                <h3>Improvements</h3>
                <ul>
                  <li>Updated skill file with all v3.x architecture modes, frontend options, and new CLI commands.</li>
                  <li>Updated sidebar with new pages: Architecture Modes, New CLI Commands, TanStack Router, Deploy Command.</li>
                  <li>Frontend sidebar section renamed from &ldquo;Frontend (Next.js)&rdquo; to &ldquo;Frontend&rdquo; to reflect multi-framework support.</li>
                </ul>
              </div>
            </div>

            {/* v3.4.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v3.4.0
                </span>
                <span className="text-sm text-muted-foreground">March 26, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Multi-architecture code generator</strong> &mdash; <code>grit generate resource</code>
                    now works for all 5 architecture modes and both frontend frameworks. Generates Go model,
                    service, and handler at the correct path (<code>internal/</code> for single app,
                    <code>apps/api/internal/</code> for monorepo). Generates React Query hooks and admin
                    resource pages for both Next.js and TanStack Router.
                  </li>
                  <li>
                    <strong><code>grit.json</code> project manifest</strong> &mdash; Every scaffolded project
                    now includes a <code>grit.json</code> file at the root with <code>architecture</code> and
                    <code>frontend</code> fields. The generator reads this to determine correct file paths
                    and template variants, eliminating fragile filesystem heuristics.
                  </li>
                  <li>
                    <strong>TanStack Router resource generation</strong> &mdash; When generating resources
                    in a TanStack Router project, creates route files at
                    <code>src/routes/_dashboard/resources/</code> using <code>createFileRoute</code> instead
                    of Next.js <code>app/(dashboard)/resources/</code> page convention.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.3.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v3.3.0
                </span>
                <span className="text-sm text-muted-foreground">March 26, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features (Goravel-Inspired)</h3>
                <ul>
                  <li>
                    <strong><code>grit routes</code></strong> &mdash; List all registered API routes in a
                    formatted table. Parses <code>routes.go</code> and shows method, path, handler, and
                    middleware group (public/protected/admin). Works for both monorepo and single app projects.
                  </li>
                  <li>
                    <strong><code>grit down</code> / <code>grit up</code></strong> &mdash; Maintenance mode.
                    <code>grit down</code> creates a <code>.maintenance</code> file that triggers the new
                    maintenance middleware, returning 503 for all requests. <code>grit up</code> removes
                    it and resumes normal operation.
                  </li>
                  <li>
                    <strong><code>grit deploy</code></strong> &mdash; One-command production deployment.
                    Cross-compiles for Linux, builds frontend, uploads binary via SCP, configures a systemd
                    service, and optionally sets up Caddy reverse proxy with auto-TLS. Supports
                    <code>--host</code>, <code>--domain</code>, <code>--key</code> flags or
                    <code>DEPLOY_HOST</code>/<code>DEPLOY_DOMAIN</code>/<code>DEPLOY_KEY_FILE</code> env vars.
                  </li>
                  <li>
                    <strong>Maintenance middleware</strong> &mdash; All scaffolded projects now include a
                    <code>Maintenance()</code> Gin middleware that checks for a <code>.maintenance</code>
                    file on every request. Runs as the first global middleware.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.2.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v3.2.0
                </span>
                <span className="text-sm text-muted-foreground">March 26, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Single app architecture</strong> &mdash; <code>grit new my-app --single</code> creates
                    a single Go binary that serves both the API and an embedded React SPA. Uses <code>go:embed</code>
                    to bake the built frontend into the binary at compile time. One file to deploy. Dev mode runs
                    Go on <code>:8080</code> and Vite on <code>:5173</code> with API proxy.
                  </li>
                  <li>
                    <strong>Parameterized API paths</strong> &mdash; All Go API file generators now use
                    <code>opts.APIRoot()</code> and <code>opts.Module()</code> helpers, enabling the same
                    template functions to generate files for both monorepo (<code>apps/api/</code>) and
                    single app (project root) architectures.
                  </li>
                </ul>

                <h3>Single App Structure</h3>
                <ul>
                  <li><code>cmd/server/main.go</code> &mdash; Entry point with <code>go:embed frontend/dist/*</code> and SPA fallback routing</li>
                  <li><code>internal/</code> &mdash; Full Go backend (same as monorepo API)</li>
                  <li><code>frontend/</code> &mdash; React + Vite + TanStack Router SPA</li>
                  <li><code>Makefile</code> &mdash; <code>make dev</code> (parallel servers), <code>make build</code> (single binary)</li>
                </ul>
              </div>
            </div>

            {/* v3.1.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v3.1.0
                </span>
                <span className="text-sm text-muted-foreground">March 26, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>TanStack Router frontend scaffold</strong> &mdash; When selecting
                    TanStack Router (Vite) as your frontend, both the web app and admin panel are
                    now fully scaffolded with Vite + TanStack Router + React Query + Tailwind CSS.
                    Includes file-based routing via <code>@tanstack/router-vite-plugin</code>,
                    API proxy in dev mode, and all the same features as the Next.js scaffold.
                  </li>
                  <li>
                    <strong>TanStack Router admin panel</strong> &mdash; Complete admin panel with
                    TanStack Router: auth pages (login, sign-up, forgot password), dashboard layout
                    with sidebar, resource management (users, blogs) via ResourcePage component,
                    system pages (jobs, files, cron, mail, security), profile page. All existing
                    React components (DataTable, FormBuilder, widgets) are reused with automatic
                    <code>&quot;use client&quot;</code> directive stripping.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.0.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v3.0.0
                </span>
                <span className="text-sm text-muted-foreground">March 26, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Interactive project creation</strong> &mdash; <code>grit new my-app</code> now
                    launches an interactive prompt to select your architecture and frontend framework.
                    Power users can skip with flags: <code>--single --vite</code>, <code>--triple --next</code>,
                    <code>--api</code>, etc.
                  </li>
                  <li>
                    <strong>5 architecture modes</strong> &mdash; Choose the project structure that fits your team:
                    <strong>Single</strong> (Go API + embedded React SPA, one binary),
                    <strong>Double</strong> (Web + API Turborepo),
                    <strong>Triple</strong> (Web + Admin + API Turborepo),
                    <strong>API Only</strong> (Go backend, no frontend),
                    <strong>Mobile</strong> (API + Expo React Native).
                  </li>
                  <li>
                    <strong>Frontend framework choice</strong> &mdash; Pick between <strong>Next.js</strong> (SSR,
                    App Router) and <strong>TanStack Router</strong> (Vite, fast builds, small bundle, SPA).
                    Available for all architecture modes that include a frontend.
                  </li>
                </ul>

                <h3>Breaking Changes</h3>
                <ul>
                  <li>
                    <strong>Options struct refactored</strong> &mdash; The internal <code>Options</code> struct
                    now uses <code>Architecture</code> and <code>Frontend</code> enum fields instead of boolean
                    flags. Legacy flags (<code>--api</code>, <code>--mobile</code>, <code>--full</code>) still
                    work via the <code>Normalize()</code> migration layer.
                  </li>
                </ul>
              </div>
            </div>

            {/* v2.9.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v2.9.0
                </span>
                <span className="text-sm text-muted-foreground">March 26, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Two-Factor Authentication (TOTP)</strong> &mdash; Every <code>grit new</code> project
                    now includes a complete 2FA system with authenticator app support (Google Authenticator,
                    Authy, 1Password, etc.). Zero-dependency RFC 6238 implementation with HMAC-SHA1.
                    Includes setup flow with QR code URI generation, 6-digit code verification with
                    &plusmn;1 window clock skew tolerance, and seamless integration with the existing
                    JWT login flow.
                  </li>
                  <li>
                    <strong>Backup Codes</strong> &mdash; 10 one-time-use recovery codes generated when
                    enabling 2FA. Each code is individually bcrypt-hashed for storage. Codes can be
                    regenerated at any time (invalidates previous set). Use during login as an alternative
                    to the authenticator app.
                  </li>
                  <li>
                    <strong>Trusted Devices</strong> &mdash; &ldquo;Remember this device&rdquo; option
                    during TOTP verification. Sets an HttpOnly cookie with a SHA-256 hashed token stored
                    in the database. Trusted devices last 30 days with sliding expiry (refreshed on each use).
                    Users can revoke all trusted devices from their account.
                  </li>
                </ul>

                <h3>New Endpoints</h3>
                <ul>
                  <li><code>POST /api/auth/totp/setup</code> &mdash; Generate TOTP secret + QR URI (authenticated)</li>
                  <li><code>POST /api/auth/totp/enable</code> &mdash; Verify initial code and activate 2FA</li>
                  <li><code>POST /api/auth/totp/verify</code> &mdash; Verify TOTP code during login (public, uses pending token)</li>
                  <li><code>POST /api/auth/totp/backup-codes/verify</code> &mdash; Use backup code during login</li>
                  <li><code>POST /api/auth/totp/disable</code> &mdash; Disable 2FA (requires password)</li>
                  <li><code>GET /api/auth/totp/status</code> &mdash; Check 2FA status, remaining backup codes, trusted device count</li>
                  <li><code>POST /api/auth/totp/backup-codes</code> &mdash; Regenerate backup codes</li>
                  <li><code>DELETE /api/auth/totp/trusted-devices</code> &mdash; Revoke all trusted devices</li>
                </ul>
              </div>
            </div>

            {/* v2.8.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v2.8.0
                </span>
                <span className="text-sm text-muted-foreground">March 16, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Vercel AI Gateway integration</strong> &mdash; Replaced the multi-provider AI service
                    (Claude, OpenAI, Gemini with separate API implementations) with{' '}
                    <a href="https://vercel.com/ai-gateway" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Vercel AI Gateway</a>.
                    One API key now gives access to hundreds of models from all major providers through a single
                    OpenAI-compatible endpoint. Models use the <code>provider/model</code> format
                    (e.g. <code>anthropic/claude-sonnet-4-6</code>, <code>openai/gpt-5.4</code>,{' '}
                    <code>google/gemini-2.5-pro</code>). Includes automatic retries, fallbacks,
                    spend monitoring, and zero markup on tokens.
                  </li>
                </ul>

                <h3>Breaking Changes</h3>
                <ul>
                  <li>
                    <strong>AI environment variables</strong> &mdash; <code>AI_PROVIDER</code>,{' '}
                    <code>AI_API_KEY</code>, and <code>AI_MODEL</code> have been replaced with{' '}
                    <code>AI_GATEWAY_API_KEY</code>, <code>AI_GATEWAY_MODEL</code>, and{' '}
                    <code>AI_GATEWAY_URL</code>. Update your <code>.env</code> file accordingly.
                    Get your API key from{' '}
                    <a href="https://vercel.com/ai-gateway" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">vercel.com/ai-gateway</a>.
                  </li>
                </ul>
              </div>
            </div>

            {/* v2.7.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v2.7.0
                </span>
                <span className="text-sm text-muted-foreground">March 10, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>10 Official Plugins</strong> &mdash; New <code>grit-plugins</code> ecosystem with
                    drop-in Go packages for common functionality: WebSockets (<code>grit-websockets</code>),
                    Stripe payments (<code>grit-stripe</code>), OAuth social login (<code>grit-oauth</code>),
                    notifications (<code>grit-notifications</code>), full-text search (<code>grit-search</code>),
                    video processing (<code>grit-video</code>), WebRTC conferencing (<code>grit-conference</code>),
                    outgoing webhooks (<code>grit-webhooks</code>), i18n translations (<code>grit-i18n</code>),
                    and PDF/Excel/CSV export (<code>grit-export</code>). Each plugin includes a Claude Code
                    skill file for AI-assisted integration.
                  </li>
                  <li>
                    <strong>Claude Code Skills format</strong> &mdash; Updated the scaffolded AI skill file
                    from a monolithic <code>GRIT_SKILL.md</code> to the official Claude Code skills directory
                    structure (<code>.claude/skills/grit/SKILL.md</code> + <code>reference.md</code>) with
                    YAML frontmatter. AI assistants can now discover and use Grit conventions automatically.
                  </li>
                  <li>
                    <strong>Grit UI component registry (100 components)</strong> &mdash; Expanded from 91 to
                    100 pre-built components across 5 categories: marketing (21), auth (10), SaaS (30),
                    ecommerce (20), and layout (20).
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/plugins" className="text-primary hover:underline">Plugins</Link>{' '}
                    page &mdash; overview of all 10 plugins with installation, environment setup,
                    quick start code, features, and use cases for each.
                  </li>
                </ul>
              </div>
            </div>

            {/* v2.6.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v2.6.0
                </span>
                <span className="text-sm text-muted-foreground">March 6, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Fixes</h3>
                <ul>
                  <li>
                    <strong>GORM Studio (Desktop)</strong> &mdash; Replaced the broken custom HTML studio with the real{' '}
                    <code>gorm-studio</code> package. Desktop studio now runs on port 8080 at <code>/studio</code> using
                    Gin + gorm-studio, matching the web scaffold. Auto-opens browser on launch.
                  </li>
                </ul>
              </div>
            </div>

            {/* v2.5.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v2.5.0
                </span>
                <span className="text-sm text-muted-foreground">March 6, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>GRIT_SKILL.md</strong> &mdash; Desktop scaffolds now include a <code>GRIT_SKILL.md</code> file
                    in the project root. This is a comprehensive AI reference (12 sections) covering architecture,
                    CLI commands, resource generation, field types, code markers, golden rules, and common LLM mistakes
                    &mdash; so AI assistants can work with the project correctly out of the box.
                  </li>
                  <li>
                    <strong>Comprehensive README</strong> &mdash; The scaffolded <code>README.md</code> now includes a
                    full project walkthrough, &ldquo;Adding a New Module&rdquo; guide, supported field types table,
                    customization section (window size, title bar, database, app name), code markers reference,
                    and a ready-to-use AI prompt for building a Task Manager app.
                  </li>
                </ul>

                <h3>Fixes</h3>
                <ul>
                  <li>
                    <strong>Dashboard stats cache</strong> &mdash; Dashboard statistics now update immediately after
                    creating a blog or contact. Changed query keys from <code>[&quot;blogs-stats&quot;]</code> to{' '}
                    <code>[&quot;blogs&quot;, &quot;stats&quot;]</code> so TanStack Query{`'`}s prefix matching
                    invalidates dashboard queries when resources are created or deleted.
                  </li>
                </ul>
              </div>
            </div>

            {/* v2.4.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v2.4.0
                </span>
                <span className="text-sm text-muted-foreground">March 5, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Window controls on auth pages</strong> &mdash; Login and register pages now include
                    minimize, maximize, and close buttons with a draggable title area, so users can move and
                    manage the window before signing in.
                  </li>
                  <li>
                    <strong>Show/hide password toggle</strong> &mdash; All password fields on login and register
                    pages now have an eye icon toggle to reveal or hide the password text.
                  </li>
                </ul>

                <h3>Fixes</h3>
                <ul>
                  <li>
                    <strong>Desktop build script</strong> &mdash; Removed <code>tsc</code> from the frontend
                    build script. TanStack Router{`'`}s Vite plugin generates <code>routeTree.gen.ts</code> during
                    the Vite build, so running <code>tsc</code> before Vite caused{' '}
                    <code>Cannot find module {`'`}./routeTree.gen{`'`}</code> errors.
                  </li>
                  <li>
                    <strong>Title bar import path</strong> &mdash; Fixed the Wails binding import in{' '}
                    <code>title-bar.tsx</code> from a 2-level to 3-level relative path.
                  </li>
                  <li>
                    <strong>Auth hook file extension</strong> &mdash; Renamed <code>use-auth.ts</code> to{' '}
                    <code>use-auth.tsx</code> so TypeScript handles the JSX correctly.
                  </li>
                  <li>
                    <strong>Create resource cache refresh</strong> &mdash; Blog and contact create pages now
                    invalidate the React Query cache before navigating back, so new records appear in the
                    table immediately.
                  </li>
                </ul>
              </div>
            </div>

            {/* v2.2.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v2.2.0
                </span>
                <span className="text-sm text-muted-foreground">March 4, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Fixes</h3>
                <ul>
                  <li>
                    <strong>Desktop auth hook file extension</strong> &mdash; Renamed the scaffolded{' '}
                    <code>use-auth.ts</code> to <code>use-auth.tsx</code> so TypeScript correctly handles
                    the JSX in <code>&lt;AuthContext.Provider&gt;</code>. Previously, <code>grit new-desktop</code>{' '}
                    projects would fail to compile with <code>TS1005: {`'>'`} expected</code> errors.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    Added Desktop Handbook PDF download links to all 8 desktop documentation pages.
                  </li>
                </ul>
              </div>
            </div>

            {/* v2.1.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v2.1.0
                </span>
                <span className="text-sm text-muted-foreground">March 4, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>TanStack Router for desktop</strong> &mdash; Migrated the desktop frontend from
                    React Router to{' '}
                    <a href="https://tanstack.com/router" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">TanStack Router</a>{' '}
                    with file-based routing. Routes are auto-discovered by the Vite plugin &mdash; no centralized
                    route registry. Uses <code>createHashHistory()</code> for Wails compatibility and{' '}
                    <code>Route.useParams()</code> for type-safe params. Resource generation now creates 5 files
                    (list, new, edit routes + model + service) and performs 10 injections (down from 12).
                  </li>
                  <li>
                    <strong>Mobile navigation</strong> &mdash; Added a hamburger menu to the docs site header,
                    visible below the <code>lg</code> breakpoint. Opens a Sheet sidebar with all navigation links.
                    Auto-closes on link click.
                  </li>
                  <li>
                    <strong>CGO-free SQLite</strong> &mdash; Replaced <code>gorm.io/driver/sqlite</code> (requires
                    CGO) with <code>github.com/glebarez/sqlite</code> (pure Go) in all scaffold templates. Desktop
                    apps now build and run without CGO or a C compiler.
                  </li>
                  <li>
                    <strong>20 Desktop Project Ideas</strong> &mdash; New{' '}
                    <Link href="/docs/desktop/project-ideas" className="text-primary hover:underline">project ideas page</Link>{' '}
                    with 20 ready-to-build desktop app ideas across business, education, healthcare, logistics,
                    and more. Each includes resources, field definitions, and <code>grit generate</code> commands.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    Added TanStack Router explanations to all desktop doc pages: overview, getting started,
                    first app, resource generation, and POS app.
                  </li>
                  <li>
                    Updated{' '}
                    <Link href="/docs/desktop/llm-reference" className="text-primary hover:underline">LLM Reference</Link>,{' '}
                    <Link href="/docs/ai-skill" className="text-primary hover:underline">GRIT_SKILL.md</Link>, and
                    database docs to reflect TanStack Router and CGO-free SQLite changes.
                  </li>
                </ul>
              </div>
            </div>

            {/* v2.0.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v2.0.0
                </span>
                <span className="text-sm text-muted-foreground">March 4, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Native desktop apps (Wails)</strong> &mdash; New <code>grit new-desktop</code> command
                    scaffolds a complete desktop application with Go backend, React frontend (Vite + TanStack Router +
                    TanStack Query), SQLite database, JWT authentication, blog and contact CRUD, PDF/Excel export,
                    custom title bar, dark theme, and GORM Studio. Compiles to a single native executable for
                    Windows, macOS, and Linux. See{' '}
                    <Link href="/docs/desktop" className="text-primary hover:underline">Desktop docs</Link>.
                  </li>
                  <li>
                    <strong>Desktop resource generation</strong> &mdash; <code>grit generate resource</code> now
                    works inside desktop projects. Generates Go model, service, and TanStack Router route files
                    (list, new, edit), then injects code into 10 locations (db.go, main.go, app.go, types.go,
                    sidebar.tsx, studio/main.go) using <code>grit:</code> markers. See{' '}
                    <Link href="/docs/desktop/resource-generation" className="text-primary hover:underline">Desktop Resource Generation</Link>.
                  </li>
                  <li>
                    <strong>Project type auto-detection</strong> &mdash; All CLI commands now auto-detect whether
                    you are inside a web (Turborepo) or desktop (Wails) project. No flags needed.
                  </li>
                  <li>
                    <strong><code>grit start</code> for desktop</strong> &mdash; Running <code>grit start</code>{' '}
                    inside a desktop project launches <code>wails dev</code> with hot-reload for both Go and React.
                  </li>
                  <li>
                    <strong><code>grit compile</code></strong> &mdash; New command that runs <code>wails build</code>{' '}
                    to produce a distributable native binary.
                  </li>
                  <li>
                    <strong><code>grit studio</code></strong> &mdash; New command that launches GORM Studio. For
                    desktop projects it starts a standalone server on port 4000. For web projects it opens the
                    browser to the embedded Studio route.
                  </li>
                  <li>
                    <strong><code>grit remove resource</code> for desktop</strong> &mdash; Removes a previously
                    generated desktop resource, deleting files and reversing all 10 marker injections.
                  </li>
                  <li>
                    <strong>Grit UI component registry (91 components)</strong> &mdash; Every scaffolded web
                    project now includes a shadcn-compatible component registry with 91 pre-built components
                    across 5 categories: marketing (14), auth (10), SaaS (30), ecommerce (20), and layout (18).
                    Install via <code>npx shadcn@latest add</code> from <code>/r</code> endpoints.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/desktop" className="text-primary hover:underline">Desktop (Wails)</Link>{' '}
                    section &mdash; 8 pages covering overview, getting started, first app tutorial, POS app
                    tutorial, resource generation, building/distribution, project ideas, and LLM reference.
                  </li>
                  <li>
                    Updated{' '}
                    <Link href="/docs/ai-skill/llm-guide" className="text-primary hover:underline">LLM Reference</Link>{' '}
                    with complete desktop section: project structure, CLI commands, markers, and architecture comparison.
                  </li>
                </ul>
              </div>
            </div>

            {/* v1.4.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v1.4.0
                </span>
                <span className="text-sm text-muted-foreground">March 2, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Gzip response compression</strong> &mdash; All API responses are now
                    compressed automatically via a custom <code>Gzip()</code> middleware using the Go
                    standard library <code>compress/gzip</code> at <code>BestSpeed</code>.
                    JSON payloads shrink by 60–80%, reducing bandwidth on paginated list endpoints
                    with zero external dependencies.
                  </li>
                  <li>
                    <strong>Request ID tracing</strong> &mdash; A <code>RequestID()</code> middleware
                    injects a unique <code>X-Request-ID</code> header on every request (echoes the
                    upstream header or generates a nanosecond-based ID). The ID is stored in Gin
                    context and included in every structured log line for end-to-end request tracing.
                  </li>
                  <li>
                    <strong>Database connection pool tuning</strong> &mdash; The scaffold now sets
                    four GORM pool parameters: <code>MaxIdleConns(10)</code>,{' '}
                    <code>MaxOpenConns(100)</code>, <code>ConnMaxLifetime(30m)</code>, and{' '}
                    <code>ConnMaxIdleTime(10m)</code>. This prevents stale connections after network
                    interruptions and avoids connection exhaustion under load.
                  </li>
                  <li>
                    <strong>Cache-Control headers on public blog endpoints</strong> &mdash; The{' '}
                    <code>ListPublished</code> handler now returns{' '}
                    <code>Cache-Control: public, max-age=300</code> (5 minutes) and{' '}
                    <code>GetBySlug</code> returns <code>Cache-Control: public, max-age=3600</code>{' '}
                    (1 hour). CDNs and edge caches can now serve public blog content without hitting
                    the Go API.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/concepts/performance" className="text-primary hover:underline">Performance</Link>{' '}
                    page &mdash; comprehensive guide to all backend (Go/API) and frontend (Next.js)
                    performance optimisations that ship with every Grit project out of the box.
                    Covers Gzip, Request ID, connection pool, Cache-Control, presigned uploads,
                    background jobs, Redis caching, Server Components, ISR, React Query, next/image,
                    Turborepo, and code splitting.
                  </li>
                  <li>
                    New{' '}
                    <Link href="/docs/ai-skill/llm-guide" className="text-primary hover:underline">Complete LLM Reference</Link>{' '}
                    page &mdash; a dedicated machine-readable guide that teaches AI assistants
                    everything about Grit: project structure, all CLI commands, every field type,
                    code patterns, API response format, code markers, naming conventions, all
                    batteries, performance features, and the golden rules that must never be broken.
                  </li>
                </ul>
              </div>
            </div>

            {/* v1.3.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v1.3.0
                </span>
                <span className="text-sm text-muted-foreground">February 26, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Presigned URL uploads</strong> &mdash; File uploads now bypass the API server entirely.
                    The browser gets a presigned PUT URL, uploads directly to S3/R2/MinIO, then records the upload
                    in the database. This fixes file uploads breaking behind reverse proxies (Dokploy/Traefik/Nginx)
                    due to request body size limits and timeouts. Includes progress tracking via XHR.
                  </li>
                  <li>
                    <strong>Error pages for scaffolded apps</strong> &mdash; New <code>grit new</code> projects now include{' '}
                    <code>error.tsx</code>, <code>not-found.tsx</code>, and <code>global-error.tsx</code> for both
                    admin and web apps. Errors are displayed with styled UI instead of the default Next.js error page.
                  </li>
                  <li>
                    <strong>Production-ready Docker config</strong> &mdash; <code>docker-compose.prod.yml</code> now uses{' '}
                    <code>expose</code> instead of <code>ports</code>, <code>env_file</code> for secrets, MinIO service,
                    named bridge network, build args for <code>NEXT_PUBLIC_API_URL</code>, and Go 1.24.
                  </li>
                  <li>
                    <strong>Sentinel ExcludePaths</strong> &mdash; Pulse, GORM Studio, Sentinel, and API docs paths are
                    now excluded from rate limiting by default, fixing Pulse health checks triggering rate limits.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/getting-started/create-without-docker" className="text-primary hover:underline">Create without Docker</Link>{' '}
                    guide &mdash; set up a Grit project using Neon, Upstash, Cloudflare R2, and Resend instead of Docker.
                  </li>
                </ul>

                <h3>Infrastructure</h3>
                <ul>
                  <li>
                    Scaffold Dockerfile updated from Go 1.23 to Go 1.24
                  </li>
                  <li>
                    Next.js Dockerfile now accepts <code>NEXT_PUBLIC_API_URL</code> as a build argument
                  </li>
                  <li>
                    <code>.env</code> template includes Docker Compose production variables (<code>POSTGRES_USER</code>,{' '}
                    <code>POSTGRES_PASSWORD</code>, <code>POSTGRES_DB</code>, <code>API_URL</code>)
                  </li>
                </ul>
              </div>
            </div>

            {/* v1.1.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v1.1.0
                </span>
                <span className="text-sm text-muted-foreground">February 25, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Default font changed to Onest</strong> &mdash; New projects scaffolded with{' '}
                    <code>grit new</code> now use the{' '}
                    <a href="https://fonts.google.com/specimen/Onest" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">Onest</a>{' '}
                    Google Font for all UI text instead of DM Sans. JetBrains Mono remains the code font.
                    The font is loaded via <code>next/font/google</code> with weights 400, 500, 600, and 700.
                  </li>
                  <li>
                    <strong>Hire Us page</strong> &mdash; New{' '}
                    <Link href="/hire" className="text-primary hover:underline">/hire</Link>{' '}
                    page for professional Grit development services. Includes service offerings,
                    tech stack overview, and contact CTA.
                  </li>
                  <li>
                    <strong>Monetization banners</strong> &mdash; Docs sidebar now shows promotional cards for{' '}
                    <a href="https://gritcms.com" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">GritCMS</a>,
                    developer hiring services, and{' '}
                    <Link href="/donate" className="text-primary hover:underline">donations</Link>{' '}
                    &mdash; visible on every documentation page.
                  </li>
                  <li>
                    <strong>Grit Fullstack Course page</strong> &mdash; New{' '}
                    <Link href="/course" className="text-primary hover:underline">/course</Link>{' '}
                    page with a 10-module curriculum covering Go, React, Next.js, and the full Grit stack.
                  </li>
                </ul>

                <h3>Improvements</h3>
                <ul>
                  <li>
                    Top navigation now includes GritCMS, Hire Us, and a Sponsor heart icon for quick access
                    to all revenue channels.
                  </li>
                  <li>
                    <code>richtext</code> added to the FieldType union for better type safety in the code generator.
                  </li>
                </ul>

                <h3>Bug Fixes</h3>
                <ul>
                  <li>
                    <strong>OAuth callback fix</strong> &mdash; Fixed <code>TokenPair</code> struct field access
                    in the social login callback handler (was using map indexing instead of struct fields).
                  </li>
                  <li>
                    <strong>Course waitlist fix</strong> &mdash; Fixed Google Sheets submission to use
                    form-encoded data instead of JSON.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/getting-started/cli-cheatsheet" className="text-primary hover:underline">CLI Cheatsheet</Link>{' '}
                    page &mdash; complete reference for all Grit CLI commands with flags, field types,
                    generated files, common workflows, and full command tree.
                  </li>
                  <li>
                    New{' '}
                    <Link href="/docs/backend/oauth" className="text-primary hover:underline">Social Login (OAuth2)</Link>{' '}
                    setup guide for Google and GitHub authentication.
                  </li>
                  <li>
                    Updated Docker Cheat Sheet with force remove commands for containers and volumes.
                  </li>
                  <li>
                    Updated AI skill guide with social login (OAuth2) section.
                  </li>
                </ul>
              </div>
            </div>

            {/* v1.0.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v1.0.0
                </span>
                <span className="text-sm text-muted-foreground">February 24, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Social Login (Google + GitHub)</strong> &mdash; Every <code>grit new</code> project now
                    includes OAuth2 social authentication via{' '}
                    <a href="https://github.com/markbates/goth" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">Gothic</a>.
                    Users can sign in with Google or GitHub on all auth pages (login, register, admin).
                    Accounts are linked by email &mdash; existing users who sign in with a social provider are automatically connected.
                    Configurable via <code>GOOGLE_CLIENT_ID</code>, <code>GITHUB_CLIENT_ID</code> environment variables.
                  </li>
                  <li>
                    <strong>GORM Studio v1.0.1</strong> &mdash; Updated to the first stable tagged release of GORM Studio.
                  </li>
                </ul>

                <h3>Improvements</h3>
                <ul>
                  <li>
                    User model now includes <code>Provider</code>, <code>GoogleID</code>, and <code>GithubID</code> fields
                    for social account linking. Password field is now nullable to support OAuth-only accounts.
                  </li>
                  <li>
                    Admin users table shows Provider column with badges (Email, Google, GitHub) and new filter option.
                  </li>
                  <li>
                    Social login buttons (Google + GitHub) appear on all 4 admin style variants (default, modern, minimal, glass).
                  </li>
                </ul>
              </div>
            </div>

            {/* v0.19.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.19.0
                </span>
                <span className="text-sm text-muted-foreground">February 24, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Fixes</h3>
                <ul>
                  <li>
                    <strong>gin-docs AuthConfig</strong> &mdash; Updated scaffold template to use the new{' '}
                    <code>gindocs.AuthConfig</code> struct instead of the deprecated <code>gindocs.AuthBearer</code> constant,
                    fixing compilation errors in newly scaffolded projects.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/tutorials/contact-app" className="text-primary hover:underline">Your First App</Link>{' '}
                    tutorial &mdash; step-by-step Contact Manager guide covering project setup, resource generation, and CRUD
                  </li>
                  <li>
                    New{' '}
                    <Link href="/docs/deployment/dokploy" className="text-primary hover:underline">Dokploy Deployment</Link>{' '}
                    guide with Dockerfile examples
                  </li>
                  <li>
                    Improved terminal blocks across all tutorials with copy buttons and horizontal scroll
                  </li>
                  <li>
                    Updated{' '}
                    <Link href="/docs/backend/api-docs" className="text-primary hover:underline">API Documentation</Link>{' '}
                    page to reflect the new <code>AuthConfig</code> struct format
                  </li>
                </ul>
              </div>
            </div>

            {/* v0.18.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.18.0
                </span>
                <span className="text-sm text-muted-foreground">February 22, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Pulse (Observability)</strong> &mdash; Every <code>grit new</code> project now includes{' '}
                    <a href="https://github.com/MUKE-coder/pulse" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">Pulse</a>,
                    a self-hosted observability SDK. Provides request tracing, database monitoring, runtime metrics,
                    error tracking, health checks, alerting, Prometheus export, and an embedded React dashboard
                    at <code>/pulse</code>. Enabled by default, configurable via <code>PULSE_ENABLED</code>.
                    See <Link href="/docs/backend/pulse" className="text-primary hover:underline">Pulse docs</Link>.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/backend/pulse" className="text-primary hover:underline">Pulse (Observability)</Link> page
                    covering configuration, endpoints, health checks, alerting, Prometheus metrics, and data storage
                  </li>
                </ul>
              </div>
            </div>

            {/* v0.17.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.17.0
                </span>
                <span className="text-sm text-muted-foreground">February 22, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>API Documentation (gin-docs)</strong> &mdash; Replaced hand-written Scalar/OpenAPI
                    spec with{' '}
                    <a href="https://github.com/MUKE-coder/gin-docs" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">gin-docs</a>,
                    a zero-annotation API documentation generator. Routes and GORM models are introspected
                    automatically to produce an OpenAPI 3.1 spec with interactive Scalar or Swagger UI,
                    plus Postman and Insomnia export.
                  </li>
                  <li>
                    <strong>Dark/Light mode for Go Playground</strong> &mdash; The playground now follows the
                    site-wide theme toggle, switching between VS Code dark and light CodeMirror themes.
                  </li>
                  <li>
                    <strong>Umami Analytics</strong> &mdash; Optional visitor analytics via self-hosted Umami,
                    configured with <code>NEXT_PUBLIC_UMAMI_WEBSITE_ID</code> environment variable.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/backend/api-docs" className="text-primary hover:underline">API Documentation</Link> page
                    covering gin-docs configuration, GORM model schemas, route customization, UI switching, and spec export
                  </li>
                  <li>Full SEO + AEO implementation: sitemap, robots.txt, JSON-LD structured data, per-page metadata</li>
                </ul>

                <h3>Infrastructure</h3>
                <ul>
                  <li>Added Dockerfile for docs site deployment (Next.js standalone output)</li>
                  <li>Google Search Console verification</li>
                </ul>
              </div>
            </div>

            {/* v0.16.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.16.0
                </span>
                <span className="text-sm text-muted-foreground">February 21, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Go Playground</strong> &mdash; Interactive code editor at{' '}
                    <Link href="/playground" className="text-primary hover:underline">/playground</Link> with
                    Go syntax highlighting, code execution via the official Go Playground API, example snippets,
                    share links, and keyboard shortcuts (Ctrl+Enter to run).
                  </li>
                  <li>
                    <strong>GORM Studio updated</strong> &mdash; Updated to latest version with raw SQL editor,
                    schema export (SQL/JSON/YAML/DBML/ERD), data import/export (JSON/CSV/SQL/XLSX),
                    and Go model generation from database schema.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    <Link href="/docs/prerequisites/golang" className="text-primary hover:underline">Go for Grit Developers</Link> &mdash;
                    comprehensive rewrite with 22 sections covering methods, Gin routing, middleware, CORS,
                    handler/service architecture, GORM CRUD, migrations, seeding, JWT auth flow, and RBAC
                  </li>
                  <li>Fixed right-side table of contents for the Go prerequisites page</li>
                  <li>New Middleware and CORS sections added to Go guide</li>
                </ul>
              </div>
            </div>

            {/* v0.15.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.15.0
                </span>
                <span className="text-sm text-muted-foreground">February 20, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Security (Sentinel)</strong> &mdash; Every <code>grit new</code> project now ships with
                    a production-grade security suite powered by{' '}
                    <Link href="https://github.com/MUKE-coder/sentinel" className="text-primary hover:underline">Sentinel</Link>.
                    Includes WAF, rate limiting, brute-force protection, anomaly detection, IP geolocation,
                    security headers, and a real-time threat dashboard at <code>/sentinel/ui</code>.
                    See <Link href="/docs/batteries/security" className="text-primary hover:underline">Security docs</Link>.
                  </li>
                  <li>
                    <strong>Admin security page</strong> &mdash; New System &rarr; Security page in the admin panel
                    embeds the Sentinel dashboard for monitoring threats without leaving the admin UI.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>New: <Link href="/docs/batteries/security" className="text-primary hover:underline">Security (Sentinel)</Link> documentation page</li>
                  <li>Migrated getting-started pages (Installation, Quick Start, Troubleshooting) to use CodeBlock component</li>
                  <li>Added prerequisite learning pages for Go, Next.js, and Docker</li>
                </ul>
              </div>
            </div>

            {/* v0.14.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.14.0
                </span>
                <span className="text-sm text-muted-foreground">February 18, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Multi-step forms</strong> &mdash; New <code>formView: &quot;modal-steps&quot;</code> and{' '}
                    <code>&quot;page-steps&quot;</code> variants with horizontal/vertical step indicators,
                    per-step validation, progress bar, and clickable step navigation.
                    See <Link href="/docs/admin/multi-step-forms" className="text-primary hover:underline">Multi-Step Forms</Link>.
                  </li>
                  <li>
                    <strong>Standalone component usage</strong> &mdash; FormBuilder, FormStepper, and DataTable
                    can now be used on any page in both web and admin apps without the resource system.
                    See <Link href="/docs/admin/standalone-usage" className="text-primary hover:underline">Standalone Usage</Link>.
                  </li>
                  <li>
                    <strong>Richtext field type</strong> &mdash; New <code>richtext</code> field with Tiptap WYSIWYG
                    editor (bold, italic, headings, lists, code blocks, links, undo/redo).
                  </li>
                  <li>
                    <strong><code>string_array</code> field type</strong> &mdash; Store arrays of strings
                    using <code>datatypes.JSONSlice[string]</code>. Works with PostgreSQL and SQLite.
                    Maps to <code>string[]</code> in TypeScript and <code>z.array(z.string())</code> in Zod.
                  </li>
                  <li>
                    <strong>Built-in blog example</strong> &mdash; <code>grit new</code> now scaffolds a complete
                    blog with model, service, handler, seed data, public web pages, and admin resource definition.
                  </li>
                  <li>
                    <strong>Sidebar user avatar</strong> &mdash; Admin sidebar shows the current user&apos;s avatar
                    with a dropdown menu for profile and logout.
                  </li>
                  <li>
                    <strong>Profile avatar upload</strong> &mdash; Profile page now supports avatar image upload.
                  </li>
                  <li>
                    <strong><code>react-hook-form</code> in web app</strong> &mdash; Web app scaffold now includes{' '}
                    <code>react-hook-form</code> as a dependency, enabling standalone FormBuilder usage out of the box.
                  </li>
                </ul>

                <h3>Bug Fixes</h3>
                <ul>
                  <li>
                    <strong>Scalar API docs crash</strong> &mdash; Fixed <code>c.String</code> treating HTML as
                    a format string. Now uses <code>c.Data</code> to avoid panics when Scalar HTML
                    contains <code>%</code> characters in CSS/JS.
                  </li>
                  <li>
                    <strong>Blog route conflict</strong> &mdash; Admin blog CRUD routes moved
                    from <code>/api/blogs</code> to <code>/api/admin/blogs</code> to avoid conflict
                    with public blog routes.
                  </li>
                  <li>
                    <strong>Select dropdown styling</strong> &mdash; Fixed relationship select dropdown
                    rendering behind modals using portal-based positioning.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>New: <Link href="/docs/tutorials/product-catalog" className="text-primary hover:underline">Build a Product Catalog</Link> tutorial &mdash; resource generation, multi-step forms, standalone DataTable &amp; FormBuilder</li>
                  <li>New: <Link href="/docs/admin/multi-step-forms" className="text-primary hover:underline">Multi-Step Forms</Link> guide</li>
                  <li>New: <Link href="/docs/admin/standalone-usage" className="text-primary hover:underline">Standalone Usage</Link> guide</li>
                  <li>New: Changelog page</li>
                  <li>Updated CLI Commands, Code Generation, Quick Start, Resources, Shared Package, Web App, Seeders, and Forms pages</li>
                </ul>
              </div>
            </div>

            {/* v0.12.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.12.0
                </span>
                <span className="text-sm text-muted-foreground">February 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Relationship support</strong> &mdash; New <code>belongs_to</code> and{' '}
                    <code>many_to_many</code> field types for the code generator. Automatically creates
                    foreign keys, junction tables, and relationship-aware form fields.
                  </li>
                  <li>
                    <strong>Relationship select fields</strong> &mdash; New <code>relationship-select</code> and{' '}
                    <code>multi-relationship-select</code> form field components with search, portal-based
                    dropdowns, and tag-based multi-select.
                  </li>
                  <li>
                    <strong>Beginner tutorial</strong> &mdash; &quot;Learn Grit Step by Step&quot; tutorial
                    walking through building a full-stack app from scratch.
                  </li>
                </ul>
              </div>
            </div>

            {/* v0.11.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.11.0
                </span>
                <span className="text-sm text-muted-foreground">February 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Full-page form view</strong> &mdash; New <code>formView: &quot;page&quot;</code> option
                    renders forms as dedicated pages instead of modals.
                  </li>
                  <li>
                    <strong><code>slug</code> field type</strong> &mdash; Auto-generates URL-friendly slugs with
                    unique suffixes. Excluded from create/update forms and Zod schemas.
                  </li>
                  <li>
                    <strong>DataTable column customization</strong> &mdash; Hide/show columns, column visibility
                    toggle in table toolbar.
                  </li>
                  <li>
                    <strong><code>grit start</code> commands</strong> &mdash; <code>grit start client</code> and{' '}
                    <code>grit start server</code> for running frontend and API separately.
                  </li>
                </ul>
              </div>
            </div>

            {/* v0.10.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.10.0
                </span>
                <span className="text-sm text-muted-foreground">January 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Style variants</strong> &mdash; <code>--style</code> flag for <code>grit new</code> with
                    4 admin panel styles: default, modern, minimal, and glass.
                  </li>
                  <li>
                    <strong>Air hot reloading</strong> &mdash; Go API development with automatic rebuild on file
                    changes using Air.
                  </li>
                  <li>
                    <strong><code>grit remove resource</code></strong> &mdash; Remove a generated resource and
                    clean up all injected code (model, handler, routes, schemas, types, hooks, admin pages).
                  </li>
                  <li>
                    <strong>AI workflow docs</strong> &mdash; Guides for using Grit with Claude and Antigravity AI assistants.
                  </li>
                </ul>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Introduction
                </Link>
              </Button>
              <div />
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
