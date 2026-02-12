import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'

export default function StoragePage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Batteries</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                File Storage
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit includes an S3-compatible file storage service that works with AWS S3, Cloudflare R2,
                Backblaze B2, and MinIO (for local development). Upload files, generate thumbnails, and
                serve signed URLs -- all out of the box.
              </p>
            </div>

            <div className="prose-grit">
              {/* Supported Providers */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Supported Providers
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The storage service uses the AWS SDK v2 under the hood, which means it works with any
                  S3-compatible object storage provider. Switch providers by changing environment variables --
                  no code changes required.
                </p>

                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Provider</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Use Case</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Endpoint</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium">AWS S3</td>
                        <td className="px-4 py-2.5">Production (AWS)</td>
                        <td className="px-4 py-2.5 font-mono text-xs">https://s3.amazonaws.com</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium">Cloudflare R2</td>
                        <td className="px-4 py-2.5">Production (no egress fees)</td>
                        <td className="px-4 py-2.5 font-mono text-xs">https://&lt;id&gt;.r2.cloudflarestorage.com</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium">Backblaze B2</td>
                        <td className="px-4 py-2.5">Production (low cost)</td>
                        <td className="px-4 py-2.5 font-mono text-xs">https://s3.us-west-002.backblazeb2.com</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-medium">MinIO</td>
                        <td className="px-4 py-2.5">Local development</td>
                        <td className="px-4 py-2.5 font-mono text-xs">http://localhost:9000</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Configuration */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Configuration
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Storage is configured via environment variables. MinIO is included in the Docker Compose
                  setup, so local development works out of the box.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">.env</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`# Storage Configuration
STORAGE_ENDPOINT=http://localhost:9000   # MinIO for local dev
STORAGE_BUCKET=myapp-uploads
STORAGE_REGION=us-east-1
STORAGE_ACCESS_KEY=minioadmin
STORAGE_SECRET_KEY=minioadmin`}</pre>
                </div>
              </div>

              {/* Storage Service API */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Storage Service API
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The storage service lives at <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">internal/storage/storage.go</code> and
                  provides five core methods.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">internal/storage/storage.go</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`// New creates a Storage instance from config.
// Automatically creates the bucket if it doesn't exist.
func New(cfg config.StorageConfig) (*Storage, error)

// Upload stores a file in the bucket at the given key.
func (s *Storage) Upload(ctx context.Context, key string, reader io.Reader, contentType string) error

// Download retrieves a file from the bucket.
func (s *Storage) Download(ctx context.Context, key string) (io.ReadCloser, error)

// Delete removes a file from the bucket.
func (s *Storage) Delete(ctx context.Context, key string) error

// GetURL returns the public URL for a stored file.
func (s *Storage) GetURL(key string) string

// GetSignedURL returns a pre-signed URL valid for the given duration.
func (s *Storage) GetSignedURL(ctx context.Context, key string, duration time.Duration) (string, error)`}</pre>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">Usage Example</h3>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">example.go</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`// Upload a file
err := store.Upload(ctx, "uploads/2026/01/avatar.jpg", fileReader, "image/jpeg")

// Get the public URL
url := store.GetURL("uploads/2026/01/avatar.jpg")
// -> http://localhost:9000/myapp-uploads/uploads/2026/01/avatar.jpg

// Generate a signed URL (expires in 1 hour)
signedURL, err := store.GetSignedURL(ctx, "uploads/2026/01/avatar.jpg", time.Hour)

// Download a file
reader, err := store.Download(ctx, "uploads/2026/01/avatar.jpg")
defer reader.Close()

// Delete a file
err = store.Delete(ctx, "uploads/2026/01/avatar.jpg")`}</pre>
                </div>
              </div>

              {/* Image Processing */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Image Processing
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">internal/storage/image.go</code> file
                  provides image processing utilities. Images are automatically resized if they exceed
                  the maximum width, and thumbnails can be generated via background jobs.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">internal/storage/image.go</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`const MaxImageWidth = 1920   // Maximum width for processed images
const ThumbnailSize = 300    // Size of generated thumbnails

// ProcessImage resizes an image if it exceeds MaxImageWidth.
// Preserves aspect ratio. Returns processed bytes.
func ProcessImage(reader io.Reader, mimeType string) ([]byte, error)

// GenerateThumbnail creates a square thumbnail (300x300).
// Uses center-crop with Lanczos resampling.
func GenerateThumbnail(reader io.Reader, mimeType string) ([]byte, error)

// IsImageMimeType returns true for supported image formats.
// Supports: image/jpeg, image/png, image/gif
func IsImageMimeType(mimeType string) bool`}</pre>
                </div>
              </div>

              {/* Upload Handler */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Upload Handler
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The upload handler at <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">POST /api/uploads</code> handles
                  multipart file uploads with MIME type validation, file size limits, and automatic
                  thumbnail generation for images (via background jobs).
                </p>

                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Endpoint</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Method</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Description</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">/api/uploads</td>
                        <td className="px-4 py-2.5 font-mono text-xs">POST</td>
                        <td className="px-4 py-2.5">Upload a file (multipart/form-data)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">/api/uploads</td>
                        <td className="px-4 py-2.5 font-mono text-xs">GET</td>
                        <td className="px-4 py-2.5">List uploads (paginated)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">/api/uploads/:id</td>
                        <td className="px-4 py-2.5 font-mono text-xs">GET</td>
                        <td className="px-4 py-2.5">Get upload details</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">/api/uploads/:id</td>
                        <td className="px-4 py-2.5 font-mono text-xs">DELETE</td>
                        <td className="px-4 py-2.5">Delete upload and stored file</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3">Allowed MIME Types</h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  By default, the following file types are allowed, including images, videos, and documents.
                  You can customize the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">AllowedMimeTypes</code> map
                  in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">internal/handlers/upload.go</code>.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">internal/handlers/upload.go</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`// MaxUploadSize is the maximum file size (50 MB).
const MaxUploadSize = 50 << 20

var AllowedMimeTypes = map[string]bool{
    "image/jpeg":      true,
    "image/png":       true,
    "image/gif":       true,
    "image/webp":      true,
    "video/mp4":       true,
    "video/webm":      true,
    "video/quicktime": true,
    "application/pdf": true,
    "text/plain":      true,
    "text/csv":        true,
    "application/json": true,
    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
}`}</pre>
                </div>
              </div>

              {/* Upload Model */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Upload Model
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every uploaded file is tracked in the database via the Upload GORM model.
                  When a file is deleted, both the database record and the stored file are removed.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">internal/models/upload.go</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`type Upload struct {
    ID            uint           \`gorm:"primarykey" json:"id"\`
    Filename      string         \`gorm:"size:255;not null" json:"filename"\`
    OriginalName  string         \`gorm:"size:255;not null" json:"original_name"\`
    MimeType      string         \`gorm:"size:100;not null" json:"mime_type"\`
    Size          int64          \`gorm:"not null" json:"size"\`
    Path          string         \`gorm:"size:500;not null" json:"path"\`
    URL           string         \`gorm:"size:500;not null" json:"url"\`
    ThumbnailURL  string         \`gorm:"size:500" json:"thumbnail_url"\`
    UserID        uint           \`gorm:"not null;index" json:"user_id"\`
    CreatedAt     time.Time      \`json:"created_at"\`
    UpdatedAt     time.Time      \`json:"updated_at"\`
    DeletedAt     gorm.DeletedAt \`gorm:"index" json:"-"\`
}`}</pre>
                </div>
              </div>

              {/* Upload Flow */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Upload Flow
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  When a file is uploaded, the handler validates it, uploads to S3, saves the record,
                  and (for images) enqueues a background job to generate a thumbnail.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">upload-flow.txt</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`POST /api/uploads (multipart/form-data)
  1. Validate file size (max 50 MB)
  2. Validate MIME type (AllowedMimeTypes)
  3. Generate unique filename: {timestamp}-{original_name}.{ext}
  4. Upload to S3 at key: uploads/{year}/{month}/{filename}
  5. Save Upload record to database
  6. If image -> enqueue ProcessImage background job
  7. Return { data: upload, message: "File uploaded successfully" }`}</pre>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">cURL Example</h3>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">curl -X POST http://localhost:8080/api/uploads \</span>
                    <br />
                    <span className="text-foreground/80">{'  '}-H &quot;Authorization: Bearer $TOKEN&quot; \</span>
                    <br />
                    <span className="text-foreground/80">{'  '}-F &quot;file=@photo.jpg&quot;</span>
                  </div>
                </div>
              </div>

              {/* Admin File Browser */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Admin File Browser
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The admin panel includes a file browser page where administrators can view all uploaded
                  files in a grid layout, filter by MIME type, view file details, and delete files.
                  Image thumbnails are displayed inline when available.
                </p>

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5">
                  <p className="text-sm text-foreground/80">
                    <strong>Thumbnail generation:</strong> When an image is uploaded, a background job
                    (via asynq) generates a 300x300 thumbnail and stores it at the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">thumbnails/</code> prefix.
                    The Upload record is updated with the thumbnail URL once processing is complete.
                    See the <Link href="/docs/batteries/jobs" className="text-primary hover:underline">Background Jobs</Link> page for details.
                  </p>
                </div>
              </div>
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
