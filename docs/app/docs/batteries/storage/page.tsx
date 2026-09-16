import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/batteries/storage')

const code = 'text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded'

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="mb-12">
      <h2 className="text-2xl font-semibold tracking-tight mb-4">{title}</h2>
      {children}
    </div>
  )
}

function P({ children }: { children: React.ReactNode }) {
  return <p className="text-muted-foreground leading-relaxed mb-4">{children}</p>
}

function Table({ head, rows }: { head: string[]; rows: React.ReactNode[][] }) {
  return (
    <div className="rounded-lg border border-border/30 bg-card/30 overflow-x-auto mb-6">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-border/30 bg-accent/20">
            {head.map((h) => (
              <th key={h} className="text-left px-4 py-2.5 font-medium text-foreground/80">{h}</th>
            ))}
          </tr>
        </thead>
        <tbody className="text-muted-foreground">
          {rows.map((row, i) => (
            <tr key={i} className={i < rows.length - 1 ? 'border-b border-border/20' : ''}>
              {row.map((cell, j) => (
                <td key={j} className="px-4 py-2.5 align-top">{cell}</td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function Note({ children }: { children: React.ReactNode }) {
  return (
    <div className="p-4 rounded-lg border border-primary/20 bg-primary/5 mb-6">
      <p className="text-sm text-foreground/80 leading-relaxed">{children}</p>
    </div>
  )
}

export default function StoragePage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Batteries</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">File Storage</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                One <code className={code}>Disk</code> interface over a directory on your machine or any
                S3-compatible bucket: AWS S3, Cloudflare R2, Backblaze B2 and MinIO. The same code stores,
                serves and deletes files on all of them, so a new project uploads before Docker is running
                and moves to a bucket in production by changing one environment variable.
              </p>
              <LaneFlow
                id="bat-storage"
                lanes={['Your code', 'storage package', 'Driver']}
                nodes={[
                  { id: 'up', lane: 0, row: 0, title: 'Upload handler', sub: 'multipart or presign', tone: 'blue' },
                  { id: 'backup', lane: 0, row: 1, title: 'Backups', sub: 'storage.Disks.Get', tone: 'blue' },
                  { id: 'disk', lane: 1, row: 0, title: 'Disk', sub: 'Store, ServeFile, TemporaryURL', tone: 'primary' },
                  { id: 'media', lane: 1, row: 1, title: 'media.Transform', sub: 'orient, resize, thumb', tone: 'cyan' },
                  { id: 'local', lane: 2, row: 0, title: 'local', sub: 'storage/app', tone: 'green' },
                  { id: 's3', lane: 2, row: 1, title: 'S3 · R2 · B2 · MinIO', sub: 'a bucket', tone: 'green' },
                ]}
                edges={[
                  { from: 'up', to: 'disk', label: 'Put', tone: 'blue' },
                  { from: 'up', to: 'media', label: 'images', tone: 'cyan' },
                  { from: 'backup', to: 'disk', tone: 'blue' },
                  { from: 'disk', to: 'local', tone: 'green' },
                  { from: 'disk', to: 's3', tone: 'green' },
                ]}
                legend={[
                  { tone: 'primary', label: 'One interface' },
                  { tone: 'green', label: 'Any driver' },
                ]}
                caption="Handlers and jobs talk to a Disk; STORAGE_DRIVER decides what is behind it"
              />
            </div>

            <div className="prose-grit">
              <Section title="Drivers">
                <P>
                  <code className={code}>STORAGE_DRIVER</code> picks the driver, and only that driver&apos;s settings are
                  read. Every driver passes the same test suite in{' '}
                  <code className={code}>internal/storage/disk_test.go</code>: put, get, exists, stat, copy, move, list,
                  delete many, public URL, temporary URL (and its expiry), and refused path traversal.
                </P>
                <Table
                  head={['Driver', 'Use it for', 'Settings']}
                  rows={[
                    [<code key="l" className={code}>local</code>, 'Development without Docker, or one server with the directory on a volume', 'STORAGE_LOCAL_ROOT, STORAGE_URL_SECRET'],
                    [<code key="m" className={code}>minio</code>, 'Development with Docker (the compose file runs it)', 'MINIO_*'],
                    [<code key="s" className={code}>s3</code>, 'AWS. No endpoint needed; with no key, an IAM role supplies credentials', 'S3_BUCKET, S3_REGION, S3_ACCESS_KEY, S3_SECRET_KEY, S3_PUBLIC_URL'],
                    [<code key="r" className={code}>r2</code>, 'Cloudflare R2, no egress fees', 'R2_* (R2_PUBLIC_URL is required for images to display)'],
                    [<code key="b" className={code}>b2</code>, 'Backblaze B2', 'B2_*'],
                  ]}
                />
                <CodeBlock language="bash" filename=".env" code={`STORAGE_DRIVER=minio        # local, minio, s3, r2 or b2

# AWS S3: leave S3_ENDPOINT empty for the regional default. With no
# S3_ACCESS_KEY the SDK's credential chain (AWS_* variables, a profile,
# or the IAM role of the EC2, ECS or Lambda it runs on) is used.
S3_BUCKET=myapp-uploads
S3_REGION=eu-west-1

# Cloudflare R2: the S3 endpoint only answers signed requests, so
# images need the bucket's public origin (r2.dev or a custom domain).
R2_ENDPOINT=https://<account>.r2.cloudflarestorage.com
R2_BUCKET=myapp-uploads
R2_PUBLIC_URL=https://pub-abc123.r2.dev`} />
              </Section>

              <Section title="Local development">
                <P>
                  Outside production, <code className={code}>STORAGE_DRIVER=minio</code> with no{' '}
                  <code className={code}>MINIO_ACCESS_KEY</code> becomes <code className={code}>local</code>, so a fresh
                  project stores uploads, shows thumbnails and deletes files with nothing running but the API. Files
                  live in <code className={code}>storage/app</code> (ignored by git) and the API serves them from{' '}
                  <code className={code}>/files/*</code>.
                </P>
                <P>
                  Keys under the public prefixes are served to anyone. Every other key, backups and private originals
                  included, is served only through a signed temporary URL: an HMAC over the key and its expiry, signed with{' '}
                  <code className={code}>STORAGE_URL_SECRET</code> or, when that is unset, <code className={code}>JWT_SECRET</code>.
                  Writes go to a temporary file and are renamed into place, and a key that climbs out of the directory is refused.
                  Each file is sent with a sandboxing Content-Security-Policy, so an uploaded page cannot run script as your API.
                </P>
                <Note>
                  <strong>Production refuses <code className={code}>local</code></strong> unless{' '}
                  <code className={code}>ALLOW_LOCAL_STORAGE_IN_PRODUCTION=true</code>. Files on one machine are invisible
                  to a second replica and lost when a container is replaced, so set it only for a single server whose{' '}
                  <code className={code}>STORAGE_LOCAL_ROOT</code> is on a volume you back up.
                </Note>
              </Section>

              <Section title="The Disk interface">
                <P>
                  Every driver implements <code className={code}>storage.Disk</code>. A missing file is{' '}
                  <code className={code}>ErrNotFound</code>, a key that could name something outside the store is{' '}
                  <code className={code}>ErrInvalidKey</code>, and deleting a key that is already gone is not an error, on every driver.
                </P>
                <CodeBlock language="go" filename="internal/storage/disk.go" code={`type Disk interface {
    Put(ctx context.Context, key string, r io.Reader, opts PutOptions) error
    Get(ctx context.Context, key string) (io.ReadCloser, error)
    Exists(ctx context.Context, key string) (bool, error)
    Stat(ctx context.Context, key string) (Object, error) // Size, ContentType, LastModified
    Delete(ctx context.Context, keys ...string) error
    Copy(ctx context.Context, from, to string) error
    Move(ctx context.Context, from, to string) error
    List(ctx context.Context, prefix string) ([]Object, error)
    URL(key string) string
    TemporaryURL(ctx context.Context, key string, ttl time.Duration) (string, error)
}

type PutOptions struct {
    ContentType string
    Visibility  Visibility // "", VisibilityPublic or VisibilityPrivate
}`} />
                <P>
                  Handlers receive a <code className={code}>*storage.Storage</code>, which wraps one Disk and keeps the
                  method names generated code has always called (<code className={code}>Upload</code>,{' '}
                  <code className={code}>Download</code>, <code className={code}>Delete</code>,{' '}
                  <code className={code}>DeleteMany</code>, <code className={code}>GetURL</code>,{' '}
                  <code className={code}>GetSignedURL</code>, <code className={code}>Stat</code>). New code can take the
                  Disk itself with <code className={code}>Disk()</code>, and <code className={code}>storage.Wrap(disk)</code>{' '}
                  puts a driver of your own, or a fake in a test, behind the same type.
                </P>
              </Section>

              <Section title="Helpers">
                <P>
                  <code className={code}>storage.Store</code> takes a file from a form and returns its key. The key is
                  generated, <code className={code}>{'<dir>/<yyyy>/<mm>/<uuid><ext>'}</code>, so the name the file
                  arrived with never reaches storage. The content type is sniffed from the bytes: HTML and SVG are refused
                  whatever they claim, a declared image must be one, and the size limit is checked before anything is read.
                </P>
                <CodeBlock language="go" filename="handlers/avatar.go" code={`func (h *AvatarHandler) Upload(c *gin.Context) {
    header, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_FILE", "message": "No file provided"}})
        return
    }
    key, err := storage.Store(c.Request.Context(), h.Storage.Disk(), "avatars", header, storage.StoreOptions{
        MaxSize: 2 << 20,
        Allow:   func(contentType string) bool { return strings.HasPrefix(contentType, "image/") },
    })
    switch {
    case errors.Is(err, storage.ErrFileTooLarge):
        c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "FILE_TOO_LARGE", "message": "Avatars are 2 MB at most"}})
        return
    case errors.Is(err, storage.ErrFileTypeNotAllowed), errors.Is(err, storage.ErrContentMismatch):
        c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_FILE_TYPE", "message": "Upload an image"}})
        return
    case err != nil:
        c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "UPLOAD_FAILED", "message": "Failed to store the file"}})
        return
    }
    // key is avatars/2026/09/5b1c...e2.png; keep it, and the original name, on your row.
}

// StoreAs picks the name: avatars/user-42.png
key, err := storage.StoreAs(ctx, disk, "avatars", header, "user-42.png", storage.StoreOptions{})`} />
                <P>
                  <code className={code}>storage.ServeFile(c, disk, key, disposition)</code> streams a file through the
                  API with <code className={code}>Content-Type</code>, <code className={code}>Content-Length</code>,{' '}
                  <code className={code}>Last-Modified</code> and <code className={code}>Content-Disposition</code>. It
                  answers <code className={code}>If-Modified-Since</code> with 304 and a byte range with 206, so a video
                  seeks and a large download resumes. On S3 a range is one ranged GET. It is how a private file reaches a
                  user you have already authorised, on any driver.
                </P>
                <CodeBlock language="go" filename="handlers/invoice.go" code={`func (h *InvoiceHandler) PDF(c *gin.Context) {
    invoice, err := h.Service.GetForUser(c.Request.Context(), c.Param("id"), c.GetString("user_id"))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "Invoice not found"}})
        return
    }
    // storage.Inline shows it in the browser; storage.Attachment saves it.
    storage.ServeFileAs(c, h.Storage.Disk(), invoice.PDFKey, storage.Inline, invoice.Number+".pdf")
}`} />
              </Section>

              <Section title="Named disks">
                <P>
                  <code className={code}>STORAGE_DISKS</code> names more disks next to the default one. Each takes{' '}
                  <code className={code}>{'STORAGE_DISK_<NAME>_*'}</code> settings (<code className={code}>DRIVER</code>,{' '}
                  <code className={code}>BUCKET</code>, <code className={code}>ENDPOINT</code>,{' '}
                  <code className={code}>ACCESS_KEY</code>, <code className={code}>SECRET_KEY</code>,{' '}
                  <code className={code}>REGION</code>, <code className={code}>PUBLIC_URL</code>,{' '}
                  <code className={code}>ROOT</code>), and whatever it leaves unset comes from its driver&apos;s own settings.
                  The API opens them at startup and code reaches them through <code className={code}>storage.Disks</code>.
                </P>
                <CodeBlock language="bash" filename=".env" code={`STORAGE_DRIVER=s3
S3_BUCKET=myapp-uploads

# A private bucket for backups, on the same account
STORAGE_DISKS=backups
STORAGE_DISK_BACKUPS_BUCKET=myapp-backups

# Or on another provider altogether
# STORAGE_DISK_BACKUPS_DRIVER=b2
# STORAGE_DISK_BACKUPS_BUCKET=myapp-backups`} />
                <CodeBlock language="go" filename="example.go" code={`backups := storage.Disks.Get("backups") // nil when STORAGE_DISKS has no backups
def := storage.Disks.Default()`} />
                <P>
                  <code className={code}>Get</code> never falls back to the default disk: code asking for a private bucket
                  should not quietly write to the public one. Backups are the built-in user. The Data &amp; Backup page, the
                  scheduled backup and <code className={code}>grit backup</code> all write archives to the{' '}
                  <code className={code}>backups</code> disk when it exists and to the default disk when it does not, and an
                  archive taken before the disk was configured still downloads and is still pruned. A named{' '}
                  <code className={code}>local</code> disk keeps its files in <code className={code}>{'storage/<name>'}</code>{' '}
                  and is served under <code className={code}>{'/files/_disks/<name>/'}</code> by the default local disk.
                </P>
              </Section>

              <Section title="Visibility">
                <P>
                  Visibility is decided by key prefix. Keys under <code className={code}>STORAGE_PUBLIC_PREFIXES</code>{' '}
                  (default <code className={code}>uploads/,thumbnails/</code>) are readable by anyone with the URL; every
                  other key is private and read through <code className={code}>TemporaryURL</code> or{' '}
                  <code className={code}>ServeFile</code>. On S3 and MinIO the API writes a bucket policy allowing anonymous
                  reads on those prefixes and nothing else, each time it connects. On the local driver the file route enforces it.
                </P>
                <P>
                  <code className={code}>PutOptions.Visibility</code> states what you expect, and{' '}
                  <code className={code}>Put</code> refuses a key whose prefix disagrees with{' '}
                  <code className={code}>ErrVisibilityMismatch</code>, so a file you meant to keep private never lands
                  where anyone can read it, even after someone edits the prefixes.
                </P>
                <CodeBlock language="go" filename="example.go" code={`// Refused: uploads/ is public
err := disk.Put(ctx, "uploads/2026/09/contract.pdf", r, storage.PutOptions{Visibility: storage.VisibilityPrivate})
// errors.Is(err, storage.ErrVisibilityMismatch) == true

// Fine: contracts/ is not a public prefix
err = disk.Put(ctx, "contracts/2026/09/contract.pdf", r, storage.PutOptions{
    ContentType: "application/pdf",
    Visibility:  storage.VisibilityPrivate,
})
link, err := disk.TemporaryURL(ctx, "contracts/2026/09/contract.pdf", 15*time.Minute)`} />
                <Note>
                  <strong>Cloudflare R2 and Backblaze B2 do not enforce visibility.</strong> Neither has bucket policies: a
                  bucket with a public domain (r2.dev, a custom domain, a public B2 bucket) serves every key in it,
                  whatever the prefix. Keep private files in a private bucket of their own, a named disk such as{' '}
                  <code className={code}>backups</code>, and serve them with <code className={code}>TemporaryURL</code>{' '}
                  or <code className={code}>ServeFile</code>. The API logs a warning when it cannot set a bucket policy.
                </Note>
              </Section>

              <Section title="Uploads and the presign fallback">
                <P>
                  The admin and web apps upload in two steps. They optimise an image in the browser, ask for a presigned
                  PUT URL, send the bytes straight to the bucket, then record the file. The size is signed into the URL,
                  and <code className={code}>complete</code> asks the bucket what actually arrived and records only a key
                  that was presigned for the caller, once.
                </P>
                <P>
                  A driver that cannot presign (<code className={code}>local</code>) answers the presign request with{' '}
                  <code className={code}>{'{"method": "multipart"}'}</code>, and the upload client sends the file to{' '}
                  <code className={code}>POST /uploads</code> instead. Nothing in the UI changes. A multipart upload runs
                  the media pipeline on the server before storing: the primary image, its renditions, and the untouched
                  original under the private <code className={code}>originals/</code> prefix.
                </P>
                <Table
                  head={['Endpoint', 'Description']}
                  rows={[
                    [<code key="1" className={code}>POST /api/v1/uploads</code>, 'Multipart upload. Query: accepts, max_size, profile. Returns a FileRef'],
                    [<code key="2" className={code}>POST /api/v1/uploads/presign</code>, 'A presigned PUT URL under uploads/<user_id>/, or {"method": "multipart"}'],
                    [<code key="3" className={code}>POST /api/v1/uploads/complete</code>, 'Record a presigned upload after checking what the bucket holds'],
                    [<code key="4" className={code}>GET /api/v1/uploads</code>, 'Your uploads, paginated; everyone’s with uploads.view'],
                    [<code key="5" className={code}>GET /api/v1/uploads/stats</code>, 'Count and bytes, by kind (image, video, audio, pdf, spreadsheet, document)'],
                    [<code key="6" className={code}>GET /api/v1/uploads/:id</code>, 'One upload record'],
                    [<code key="7" className={code}>GET /api/v1/uploads/:id/download</code>, 'The file through the API, with Range support; ?inline=true shows it'],
                    [<code key="8" className={code}>DELETE /api/v1/uploads/:id</code>, 'Remove the record and the stored file'],
                    [<code key="9" className={code}>GET /api/v1/media/profiles</code>, 'The image optimisation profiles, so the browser uses the server’s numbers'],
                  ]}
                />
                <CodeBlock terminal code={`curl -X POST "http://localhost:8080/api/v1/uploads?accepts=image" \\
  -H "Authorization: Bearer $TOKEN" \\
  -F "file=@photo.jpg"

curl -OJ "http://localhost:8080/api/v1/uploads/$ID/download" -H "Authorization: Bearer $TOKEN"`} />
                <CodeBlock language="go" filename="internal/models/upload.go" code={`type Upload struct {
    ID           string    \`gorm:"primarykey;size:36" json:"id"\`      // a UUID
    Filename     string    \`json:"filename"\`                            // <uuid>.<ext>, as stored
    OriginalName string    \`json:"original_name"\`                       // the name it arrived with
    MimeType     string    \`json:"mime_type"\`                           // of the stored file
    Size         int64     \`json:"size"\`                                // of the stored file
    Path         string    \`json:"path"\`                                // the key
    URL          string    \`json:"url"\`
    ThumbnailURL string    \`json:"thumbnail_url"\`
    UserID       string    \`gorm:"index;size:36" json:"user_id"\`      // a UUID, not a number
    CreatedAt    time.Time \`json:"created_at"\`
}`} />
              </Section>

              <Section title="Images">
                <P>
                  Images go through one pipeline, <code className={code}>media.Transform</code>: EXIF orientation
                  applied, metadata (GPS included) stripped, resized to the profile and re-encoded, with a 400 by 400{' '}
                  <code className={code}>thumb</code> rendition. A multipart upload runs it inline. A presigned upload, or a
                  type the pipeline skipped, is handed to the <code className={code}>image:process</code> job, whose{' '}
                  <code className={code}>storage.GenerateThumbnail</code> runs the same pipeline, so a portrait phone photo
                  gets an upright thumbnail whichever way it arrived. A decompression bomb is refused from its header before
                  any pixels are decoded, and the job does not retry it. See{' '}
                  <Link href="/docs/batteries/jobs" className="text-primary hover:underline">Background Jobs</Link>.
                </P>
              </Section>

              <Section title="File lifecycle">
                <P>
                  An upload is unclaimed until a record points at it. Resource fields store a{' '}
                  <code className={code}>files.FileRef</code>, and saving the record claims the file. Uploads nobody claimed
                  within 24 hours are removed by the orphan cleanup, a page at a time with one batched delete per page, so
                  a form abandoned halfway does not leave files in the bucket forever. Deleting an upload removes its
                  stored file too.
                </P>
              </Section>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/frontend/shared-package" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Shared Package
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/batteries/email" className="gap-1.5">
                  Email System
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
