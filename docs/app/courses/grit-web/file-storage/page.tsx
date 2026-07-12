import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'File Storage & Uploads — Grit Web Course',
  description: 'Learn how to handle file uploads with S3-compatible storage using presigned URLs. Covers MinIO for local development, Cloudflare R2, AWS S3, image processing, and admin file management.',
}

export default function FileStorageCourse() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <Link href="/courses/grit-web" className="hover:text-foreground transition-colors">Grit Web</Link>
          <span>/</span>
          <span className="text-foreground">File Storage & Uploads</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Course 5 of 8</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">14 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            File Storage & Uploads
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Most applications need to store files — profile photos, documents, product images, exports.
            In this course, you will learn how Grit handles file uploads using S3-compatible storage
            and presigned URLs. You{"'"}ll start with MinIO for local development, then learn how to switch
            to cloud providers like Cloudflare R2 and AWS S3 with zero code changes.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What is File Storage? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is File Storage?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When users upload files to your application — photos, PDFs, spreadsheets, videos — you need
            somewhere to store them. You <strong className="text-foreground">cannot</strong> store files
            in your database. Databases are designed for structured data (rows and columns), not large
            binary files. Instead, files are stored in <strong className="text-foreground">object storage</strong>.
          </p>

          <Definition term="Object Storage">
            A storage system designed for files (called {'"'}objects{'"'}). Each object has a unique key
            (like a file path), the file data itself, and metadata (content type, size, upload date).
            Unlike a filesystem with folders, object storage uses a flat namespace with key prefixes
            that look like paths: <Code>uploads/users/42/photo.jpg</Code>.
          </Definition>

          <Definition term="S3 (Simple Storage Service)">
            Amazon{"'"}s object storage service, and the industry standard API for file storage. When we
            say {'"'}S3-compatible,{'"'} we mean any storage service that speaks the same API — including
            AWS S3, Cloudflare R2, MinIO, DigitalOcean Spaces, and Backblaze B2. Write your code once,
            switch providers by changing environment variables.
          </Definition>

          <Definition term="Bucket">
            A top-level container for objects in S3 storage. Think of it like a hard drive or volume.
            Your application typically uses one bucket (e.g., <Code>my-app-uploads</Code>) and organizes
            files inside it using key prefixes like <Code>uploads/</Code>, <Code>avatars/</Code>,
            or <Code>exports/</Code>.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit uses <strong className="text-foreground">S3-compatible storage</strong> for all file
            operations. This means the same code works with AWS S3 in production, Cloudflare R2 for
            cost savings, or MinIO running locally in Docker during development. Your files live
            outside the database, outside your API server — in dedicated storage built for the job.
          </p>

          <Tip>
            Why not just save files to disk? In production, your API might run on multiple servers
            behind a load balancer. Files saved to one server{"'"}s disk wouldn{"'"}t be accessible from
            the others. Object storage solves this by providing a shared, durable, and scalable
            file system accessible from anywhere.
          </Tip>
        </section>

        {/* ═══ Section 2: How Presigned URLs Work ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">How Presigned URLs Work</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit uses <strong className="text-foreground">presigned URLs</strong> for file uploads.
            This is the modern, scalable way to handle file uploads — and it{"'"}s very different from
            the traditional approach of uploading files through your API.
          </p>

          <Definition term="Presigned URL">
            A temporary, signed URL that grants permission to upload (or download) a specific file
            directly to/from S3 storage. The URL contains a cryptographic signature that expires
            after a set time (usually 15 minutes). Anyone with the URL can perform the allowed
            operation — no additional authentication needed.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s the flow, step by step:
          </p>

          <ol className="space-y-3 text-muted-foreground mb-6">
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">1.</span>
              <span><strong className="text-foreground">Frontend asks API for an upload URL.</strong> The frontend sends the filename and content type to your API: {'"'}I want to upload photo.jpg (image/jpeg).{'"'}</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">2.</span>
              <span><strong className="text-foreground">API generates a presigned URL.</strong> Your Go API creates a time-limited, signed URL that allows uploading one specific file to one specific location in S3.</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">3.</span>
              <span><strong className="text-foreground">Frontend uploads directly to S3.</strong> The frontend sends the file directly to the presigned URL using a PUT request. The file goes straight to S3 — it never passes through your API server.</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">4.</span>
              <span><strong className="text-foreground">Frontend confirms the upload.</strong> After a successful upload, the frontend tells your API: {'"'}The upload to key uploads/abc123/photo.jpg is complete.{'"'} The API saves a record in the uploads table.</span>
            </li>
          </ol>

          <CodeBlock filename="Upload Flow Diagram">
{`┌──────────┐     1. Request URL      ┌──────────┐
│          │ ─────────────────────>  │          │
│ Frontend │                         │  Go API  │
│          │ <─────────────────────  │          │
└──────────┘  2. Presigned URL       └──────────┘
     │
     │  3. PUT file directly
     v
┌──────────┐
│    S3    │  (MinIO / R2 / AWS S3)
│ Storage  │
└──────────┘
     │
     │  4. Frontend confirms upload
     v
┌──────────┐     Save upload record   ┌──────────┐
│ Frontend │ ─────────────────────>  │  Go API  │
└──────────┘                         └──────────┘`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">Why presigned URLs?</strong> Three big reasons:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">No request body limits.</strong> Your API doesn{"'"}t handle the file data, so there{"'"}s no 10MB or 50MB upload limit to configure. Users can upload files of any size directly to S3.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">No API bottleneck.</strong> Large files don{"'"}t consume your API{"'"}s memory or bandwidth. A 500MB video goes straight to S3 while your API keeps serving other requests.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Progress tracking works.</strong> Because the frontend sends the file directly via XHR, you get real upload progress events. This lets you show a proper progress bar to the user.</li>
          </ul>

          <Note>
            The traditional approach — uploading the file to your API as a multipart form, then
            having the API forward it to S3 — works for small files but breaks down at scale. Your
            API becomes a bottleneck, memory usage spikes, and you lose progress tracking. Presigned
            URLs solve all of these problems.
          </Note>

          <Challenge number={1} title="Explain Presigned URLs">
            <p>In your own words, explain why presigned URLs are better than uploading files
            through the API. Write down at least 3 advantages. Bonus: can you think of a scenario
            where uploading through the API might be acceptable instead?</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Storage Configuration ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Storage Configuration</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit{"'"}s storage configuration is driven entirely by environment variables. This means you
            can switch from local MinIO to Cloudflare R2 or AWS S3 without changing any code — just
            update your <Code>.env</Code> file.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here are all the storage-related environment variables:
          </p>

          <CodeBlock filename=".env">
{`# Storage Configuration
STORAGE_DRIVER=minio          # minio | s3 | r2 | b2
MINIO_ENDPOINT=http://localhost:9002
MINIO_BUCKET=uploads
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_REGION=us-east-1
MINIO_USE_SSL=false

# For AWS S3 set STORAGE_DRIVER=s3 and use S3_* vars;
# for Cloudflare R2 use R2_*; for Backblaze B2 use B2_*.`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Let{"'"}s break down each variable:
          </p>

          <ul className="space-y-3 text-muted-foreground mb-6">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">STORAGE_DRIVER</strong> — The storage backend: <Code>minio</Code> (local dev), <Code>s3</Code>, <Code>r2</Code>, or <Code>b2</Code>. This selects which set of provider variables Grit reads and which client it initializes.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">MINIO_ENDPOINT</strong> — The MinIO URL for local dev: <Code>http://localhost:9002</Code>. (Cloud drivers use their own endpoint var — <Code>S3_ENDPOINT</Code>, <Code>R2_ENDPOINT</Code>, <Code>B2_ENDPOINT</Code>; leave S3 empty for the AWS regional default.)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">MINIO_BUCKET</strong> — The bucket where files are stored. Grit uses a single bucket and organizes files with key prefixes.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">MINIO_ACCESS_KEY</strong> — The access key ID. For local MinIO: <Code>minioadmin</Code>. For cloud providers: your IAM access key (via <Code>S3_ACCESS_KEY</Code> / <Code>R2_ACCESS_KEY</Code>).</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">MINIO_SECRET_KEY</strong> — The secret access key. For local MinIO: <Code>minioadmin</Code>. Keep production secrets out of Git.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">MINIO_REGION</strong> — Required by the S3 protocol. For MinIO any value works (<Code>us-east-1</Code>); R2 uses <Code>auto</Code>; AWS S3 uses your real region.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">MINIO_USE_SSL</strong> — HTTPS on/off. <Code>false</Code> for local MinIO (HTTP), <Code>true</Code> for cloud providers.</li>
          </ul>

          <Tip>
            The default <Code>.env</Code> file generated by <Code>grit new</Code> is pre-configured
            for local MinIO development. You don{"'"}t need to change anything to get started — just
            run <Code>docker compose up -d</Code> and MinIO will be ready.
          </Tip>

          <Challenge number={2} title="Find Your Storage Variables">
            <p>Open your project{"'"}s <Code>.env</Code> file and find the <Code>STORAGE_DRIVER</Code> and <Code>MINIO_</Code> variables.
            What driver is configured by default? What bucket name is used? What endpoint does it
            point to?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: MinIO for Local Development ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">MinIO for Local Development</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            In development, you don{"'"}t want to use a cloud storage service for every file upload test.
            Instead, Grit uses <strong className="text-foreground">MinIO</strong> — an S3-compatible
            object storage server that runs locally in Docker.
          </p>

          <Definition term="MinIO">
            An open-source, S3-compatible object storage server. It implements the full S3 API, so
            any code written for AWS S3 works with MinIO without changes. It runs as a single binary
            or Docker container, making it perfect for local development and testing.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            MinIO is already included in your project{"'"}s <Code>docker-compose.yml</Code>. When you
            run <Code>docker compose up -d</Code>, MinIO starts alongside PostgreSQL, Redis, and
            Mailhog. It exposes two ports:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Port 9000</strong> — The S3 API endpoint. Your Go API connects here to generate presigned URLs and manage files.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Port 9001</strong> — The MinIO Console (web UI). Open <Code>localhost:9003</Code> in your browser to manage buckets and browse files visually.</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            To access the MinIO Console, open <Code>http://localhost:9003</Code> and log in with
            the default credentials:
          </p>

          <CodeBlock filename="MinIO Console Login">
{`Username: minioadmin
Password: minioadmin`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Once logged in, you can:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><span className="text-primary">•</span> Browse existing buckets and their contents</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Create new buckets</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Upload and download files manually</li>
            <li className="flex gap-2"><span className="text-primary">•</span> View file metadata (size, content type, last modified)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Set access policies on buckets</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Your project should already have an <Code>uploads</Code> bucket. If it doesn{"'"}t exist yet,
            you can create it through the console:
          </p>

          <CodeBlock filename="Creating a Bucket">
{`1. Open http://localhost:9003
2. Log in with minioadmin / minioadmin
3. Click "Buckets" in the sidebar
4. Click "Create Bucket"
5. Enter the name: uploads
6. Click "Create Bucket"

The bucket is now ready to receive files.`}
          </CodeBlock>

          <Note>
            MinIO stores data in a Docker volume, so your files persist between container restarts.
            If you run <Code>docker compose down -v</Code> (with the <Code>-v</Code> flag), the
            volumes are deleted and you{"'"}ll lose all stored files. Without <Code>-v</Code>,
            your files are safe.
          </Note>

          <Challenge number={3} title="Explore MinIO Console">
            <p>Make sure Docker is running with <Code>docker compose up -d</Code>. Open{' '}
            <Code>localhost:9003</Code> in your browser and log in with{' '}
            <Code>minioadmin</Code> / <Code>minioadmin</Code>. Can you see the{' '}
            <Code>uploads</Code> bucket? Create a test bucket called {'"'}images{'"'} and verify
            it appears in the list.</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Uploading a File ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Uploading a File</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Now let{"'"}s see how file uploads actually work in Grit. The process has two steps:
            (1) get a presigned URL from your API, and (2) upload the file directly to S3 using that URL.
          </p>

          <h3 className="text-xl font-semibold text-foreground mb-3">Step 1: Get a Presigned URL</h3>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The frontend sends a POST request to the presign endpoint with the filename and content type:
          </p>

          <CodeBlock filename="POST /api/uploads/presign">
{`// Request
{
  "filename": "photo.jpg",
  "content_type": "image/jpeg"
}

// Response
{
  "data": {
    "upload_url": "http://localhost:9002/uploads/abc123-photo.jpg?X-Amz-Algorithm=...",
    "key": "uploads/abc123-photo.jpg"
  }
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The API generates a unique key for the file (to avoid naming collisions) and returns a
            presigned URL valid for 15 minutes. The <Code>key</Code> is the file{"'"}s path inside the
            bucket — you{"'"}ll need it to reference the file later.
          </p>

          <h3 className="text-xl font-semibold text-foreground mb-3">Step 2: Upload to S3</h3>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The frontend then PUTs the file directly to the presigned URL. Here{"'"}s the Go handler
            that generates the presigned URL:
          </p>

          <CodeBlock filename="apps/api/internal/handler/upload_handler.go">
{`func (h *UploadHandler) Presign(c *gin.Context) {
    var req struct {
        Filename    string \`json:"filename" binding:"required"\`
        ContentType string \`json:"content_type" binding:"required"\`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": gin.H{
            "code":    "VALIDATION_ERROR",
            "message": "filename and content_type are required",
        }})
        return
    }

    // Generate a unique key to avoid collisions
    key := fmt.Sprintf("uploads/%s-%s", uuid.New().String()[:8], req.Filename)

    // Generate presigned PUT URL (15 min expiry)
    uploadURL, err := h.storage.PresignPut(key, req.ContentType, 15*time.Minute)
    if err != nil {
        c.JSON(500, gin.H{"error": gin.H{
            "code":    "STORAGE_ERROR",
            "message": "Failed to generate upload URL",
        }})
        return
    }

    c.JSON(200, gin.H{"data": gin.H{
        "upload_url": uploadURL,
        "key":        key,
    }})
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            And here{"'"}s the frontend upload component with progress tracking using XHR:
          </p>

          <CodeBlock filename="Frontend Upload with Progress">
{`async function uploadFile(file: File) {
  // Step 1: Get presigned URL from API
  const res = await fetch("/api/uploads/presign", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": "Bearer " + token,
    },
    body: JSON.stringify({
      filename: file.name,
      content_type: file.type,
    }),
  });
  const { data } = await res.json();

  // Step 2: Upload file directly to S3 with progress
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    xhr.open("PUT", data.upload_url);
    xhr.setRequestHeader("Content-Type", file.type);

    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable) {
        const percent = Math.round((e.loaded / e.total) * 100);
        setProgress(percent); // Update progress bar
      }
    };

    xhr.onload = () => {
      if (xhr.status === 200) {
        resolve(data.key); // Return the file key
      } else {
        reject(new Error("Upload failed"));
      }
    };

    xhr.onerror = () => reject(new Error("Upload failed"));
    xhr.send(file);
  });
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Notice how we use <Code>XMLHttpRequest</Code> instead of <Code>fetch</Code> for the
            actual upload. This is because XHR provides <Code>upload.onprogress</Code> events, which
            let us show a real-time progress bar. The <Code>fetch</Code> API doesn{"'"}t support
            upload progress tracking.
          </p>

          <Tip>
            The presigned URL includes all the authentication information in the query string.
            You don{"'"}t need to send any authorization headers when uploading to the presigned URL —
            the signature in the URL itself grants permission.
          </Tip>

          <Challenge number={4} title="Call the Presign Endpoint">
            <p>With your project running (<Code>grit start</Code>), use the API docs or a tool like
            curl to call the presign endpoint:</p>
            <CodeBlock filename="Terminal">
{`curl -X POST http://localhost:8080/api/uploads/presign \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer YOUR_TOKEN" \\
  -d '{"filename": "test.txt", "content_type": "text/plain"}'`}
            </CodeBlock>
            <p className="mt-2">Copy the <Code>upload_url</Code> from the response. What does it look like?
            Can you spot the expiration time in the URL parameters?</p>
          </Challenge>

          <Challenge number={5} title="Upload Through the Admin">
            <p>Generate a resource with an image field:</p>
            <CodeBlock filename="Terminal">
{`grit generate resource Product --fields "name:string,price:float,image:string"`}
            </CodeBlock>
            <p className="mt-2">Open the admin panel, go to the Products page, and create a new product.
            Upload an image file through the form. Check the MinIO console at{' '}
            <Code>localhost:9003</Code> — can you find your uploaded file in the{' '}
            <Code>uploads</Code> bucket?</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Image Processing ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Image Processing</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When a user uploads an image, you often need more than just the original file. A 5MB
            profile photo is too large for a 40x40 avatar in a sidebar. A 4000px product image is
            too heavy for a thumbnail grid. Grit handles this with <strong className="text-foreground">background
            image processing</strong>.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s how it works:
          </p>

          <ol className="space-y-3 text-muted-foreground mb-6">
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">1.</span>
              <span><strong className="text-foreground">Upload completes.</strong> The frontend uploads the original image to S3 via a presigned URL and notifies the API.</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">2.</span>
              <span><strong className="text-foreground">API enqueues a background job.</strong> The upload handler dispatches an image processing job to the asynq task queue (backed by Redis).</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">3.</span>
              <span><strong className="text-foreground">Worker processes the image.</strong> A background worker picks up the job, downloads the original from S3, generates thumbnail and medium-sized versions, and uploads them back to S3.</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">4.</span>
              <span><strong className="text-foreground">Multiple sizes available.</strong> The original, thumbnail, and medium versions are all stored in S3 with predictable keys.</span>
            </li>
          </ol>

          <CodeBlock filename="Image Processing Flow">
{`Original Upload:  uploads/abc123-photo.jpg       (original)
   │
   └─> Background Job (asynq worker)
       │
       ├─> uploads/abc123-photo_thumb.jpg    (150x150)
       └─> uploads/abc123-photo_medium.jpg   (800x600)`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The image processing worker in your Go API looks like this:
          </p>

          <CodeBlock filename="apps/api/internal/worker/image_worker.go">
{`func (w *ImageWorker) ProcessImage(ctx context.Context, t *asynq.Task) error {
    var payload struct {
        Key         string \`json:"key"\`
        ContentType string \`json:"content_type"\`
    }
    if err := json.Unmarshal(t.Payload(), &payload); err != nil {
        return fmt.Errorf("unmarshal payload: %w", err)
    }

    // Download original from S3
    original, err := w.storage.Get(payload.Key)
    if err != nil {
        return fmt.Errorf("download original: %w", err)
    }

    // Generate thumbnail (150x150)
    thumb, err := resize(original, 150, 150)
    if err != nil {
        return fmt.Errorf("resize thumbnail: %w", err)
    }
    thumbKey := strings.TrimSuffix(payload.Key, filepath.Ext(payload.Key)) + "_thumb" + filepath.Ext(payload.Key)
    if err := w.storage.Put(thumbKey, thumb, payload.ContentType); err != nil {
        return fmt.Errorf("upload thumbnail: %w", err)
    }

    // Generate medium (800x600)
    medium, err := resize(original, 800, 600)
    if err != nil {
        return fmt.Errorf("resize medium: %w", err)
    }
    mediumKey := strings.TrimSuffix(payload.Key, filepath.Ext(payload.Key)) + "_medium" + filepath.Ext(payload.Key)
    if err := w.storage.Put(mediumKey, medium, payload.ContentType); err != nil {
        return fmt.Errorf("upload medium: %w", err)
    }

    return nil
}`}
          </CodeBlock>

          <Note>
            Image processing happens <strong className="text-foreground">asynchronously</strong> in
            a background worker. The user{"'"}s upload request completes immediately — they don{"'"}t
            wait for thumbnails to be generated. This keeps the UI responsive even for large images.
            You{"'"}ll learn more about background jobs in the next course.
          </Note>

          <Challenge number={6} title="Check for Thumbnails">
            <p>Upload a JPG or PNG image through the admin panel (a product image, profile photo,
            etc.). Wait a few seconds for the background worker to process it, then open the MinIO
            console at <Code>localhost:9003</Code>. Browse the <Code>uploads</Code> bucket. Can you
            find the original file and its thumbnail versions (<Code>_thumb</Code> and{' '}
            <Code>_medium</Code> suffixes)?</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Switching to Cloud Storage ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Switching to Cloud Storage</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When you{"'"}re ready to deploy, you{"'"}ll switch from local MinIO to a cloud storage provider.
            The beauty of S3-compatible storage is that <strong className="text-foreground">your code
            doesn{"'"}t change at all</strong> — you only update the <Code>.env</Code> variables.
          </p>

          <h3 className="text-xl font-semibold text-foreground mb-3">Cloudflare R2</h3>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Cloudflare R2 is a popular choice because it has <strong className="text-foreground">no
            egress fees</strong> — you don{"'"}t pay for downloading files, only for storage and write
            operations. This can save a lot of money compared to AWS S3.
          </p>

          <CodeBlock filename=".env (Cloudflare R2)">
{`STORAGE_DRIVER=r2
R2_ENDPOINT=https://your-account-id.r2.cloudflarestorage.com
R2_BUCKET=my-uploads
R2_ACCESS_KEY=your-r2-access-key
R2_SECRET_KEY=your-r2-secret-key
R2_REGION=auto`}
          </CodeBlock>

          <h3 className="text-xl font-semibold text-foreground mb-3">AWS S3</h3>

          <p className="text-muted-foreground leading-relaxed mb-4">
            AWS S3 is the original and most widely used object storage. If you{"'"}re already on AWS,
            it{"'"}s the natural choice:
          </p>

          <CodeBlock filename=".env (AWS S3)">
{`STORAGE_DRIVER=s3
S3_ENDPOINT=                 # empty = AWS regional default
S3_BUCKET=my-app-uploads
S3_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE
S3_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
S3_REGION=us-west-2`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Notice: <strong className="text-foreground">the code is identical in both cases.</strong> Only
            the environment variables change. This is the power of the S3-compatible abstraction.
            You develop locally with MinIO, deploy to R2 for cost savings, or use AWS S3 if you
            need AWS-specific features — all without touching your code.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s how the storage service decides which driver to use:
          </p>

          <CodeBlock filename="apps/api/internal/storage/storage.go">
{`func NewStorage(cfg *config.Config) (*Storage, error) {
    switch cfg.StorageDriver {
    case "s3":
        client, err := minio.New(cfg.StorageEndpoint, &minio.Options{
            Creds:  credentials.NewStaticV4(cfg.StorageAccessKey, cfg.StorageSecretKey, ""),
            Secure: cfg.StorageUseSSL,
            Region: cfg.StorageRegion,
        })
        if err != nil {
            return nil, fmt.Errorf("init s3 client: %w", err)
        }
        return &Storage{client: client, bucket: cfg.StorageBucket}, nil
    default:
        return nil, fmt.Errorf("unsupported storage driver: %s", cfg.StorageDriver)
    }
}`}
          </CodeBlock>

          <Tip>
            Start with MinIO locally. When you deploy, switch to Cloudflare R2 for the best
            price-to-performance ratio. Only use AWS S3 if you need specific features like S3
            event notifications, S3 Select, or you{"'"}re already deep in the AWS ecosystem.
          </Tip>

          <Challenge number={7} title="Read the Storage Service">
            <p>Open your project{"'"}s storage service code at{' '}
            <Code>apps/api/internal/storage/</Code>. Read through the files. How does the service
            decide which driver to use? What methods does the storage service expose (e.g.,{' '}
            <Code>PresignPut</Code>, <Code>Get</Code>, <Code>Delete</Code>)? List all the public
            methods you can find.</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: The Upload Model ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The Upload Model</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every file uploaded through Grit is tracked in the database using the <Code>Upload</Code>{' '}
            model. This gives you a record of every file — who uploaded it, when, what type it is,
            and where it{"'"}s stored in S3.
          </p>

          <CodeBlock filename="apps/api/internal/models/upload.go">
{`type Upload struct {
    ID          string    \`gorm:"primaryKey;size:36" json:"id"\`
    Filename    string    \`gorm:"not null" json:"filename"\`
    Key         string    \`gorm:"not null;uniqueIndex" json:"key"\`
    ContentType string    \`gorm:"not null" json:"content_type"\`
    Size        int64     \`json:"size"\`
    UserID      string    \`gorm:"size:36" json:"user_id"\`
    User        User      \`gorm:"foreignKey:UserID" json:"user,omitempty"\`
    CreatedAt   time.Time \`json:"created_at"\`
    UpdatedAt   time.Time \`json:"updated_at"\`
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Let{"'"}s break down each field:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Filename</strong> — The original filename the user uploaded (e.g., <Code>photo.jpg</Code>). Used for display purposes and when the user downloads the file.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Key</strong> — The unique storage key in S3 (e.g., <Code>uploads/abc123-photo.jpg</Code>). This is how you locate the file in the bucket. It has a unique index to prevent duplicates.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">ContentType</strong> — The MIME type (e.g., <Code>image/jpeg</Code>, <Code>application/pdf</Code>, <Code>text/csv</Code>). Used to set the correct headers when serving the file.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Size</strong> — The file size in bytes. Useful for showing {'"'}2.4 MB{'"'} in the UI or enforcing storage quotas per user.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">UserID</strong> — Which user uploaded the file. This enables per-user file management and access control.</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When the frontend confirms a successful upload, the API creates an Upload record:
          </p>

          <CodeBlock filename="apps/api/internal/handler/upload_handler.go">
{`func (h *UploadHandler) ConfirmUpload(c *gin.Context) {
    var req struct {
        Key         string \`json:"key" binding:"required"\`
        Filename    string \`json:"filename" binding:"required"\`
        ContentType string \`json:"content_type" binding:"required"\`
        Size        int64  \`json:"size"\`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": gin.H{
            "code":    "VALIDATION_ERROR",
            "message": "key, filename, and content_type are required",
        }})
        return
    }

    userID := c.GetString("userID") // from auth middleware

    upload := models.Upload{
        Filename:    req.Filename,
        Key:         req.Key,
        ContentType: req.ContentType,
        Size:        req.Size,
        UserID:      userID,
    }
    if err := h.db.Create(&upload).Error; err != nil {
        c.JSON(500, gin.H{"error": gin.H{
            "code":    "DATABASE_ERROR",
            "message": "Failed to save upload record",
        }})
        return
    }

    c.JSON(201, gin.H{
        "data":    upload,
        "message": "Upload confirmed",
    })
}`}
          </CodeBlock>

          <Challenge number={8} title="Explore the Uploads Table">
            <p>Open GORM Studio at <Code>localhost:8080/studio</Code> and find the{' '}
            <Code>uploads</Code> table. What columns does it have? Upload a few files through the
            admin panel, then refresh the table. Can you see the records with their filenames,
            keys, content types, and sizes?</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Admin File Management ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Admin File Management</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The admin panel includes a <strong className="text-foreground">Files</strong> system page
            where administrators can view and manage all uploaded files. This page is available in
            the admin sidebar under the System section.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The Files page provides:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">File list</strong> — A DataTable showing all uploads with filename, content type, size, uploader, and upload date</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Preview</strong> — Image files show a thumbnail preview directly in the table</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Search and filter</strong> — Search by filename, filter by content type or uploader</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Download</strong> — Generate a presigned download URL to retrieve any file</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Delete</strong> — Remove files from both S3 storage and the database</li>
          </ul>

          <CodeBlock filename="Admin Files Page">
{`// The Files page uses the standard DataTable component
// with columns configured for upload records:

const columns = [
  { key: "filename", label: "Filename", sortable: true },
  { key: "content_type", label: "Type", sortable: true },
  { key: "size", label: "Size", sortable: true,
    render: (value: number) => formatBytes(value) },
  { key: "user", label: "Uploaded By", sortable: true,
    render: (value: any) => value?.name || "Unknown" },
  { key: "created_at", label: "Uploaded", sortable: true,
    render: (value: string) => formatDate(value) },
];`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When you delete a file from the admin panel, Grit performs a two-step cleanup:
          </p>

          <ol className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">1.</span>
              <span>Delete the file (and any thumbnails) from S3 storage</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">2.</span>
              <span>Delete the Upload record from the database</span>
            </li>
          </ol>

          <Tip>
            The Files page is one of several <strong className="text-foreground">system pages</strong> built
            into the admin panel. Others include Jobs (background task monitoring), Cron (scheduled
            task management), and Mail Preview (email template testing). These pages are always
            available regardless of what resources you generate.
          </Tip>

          <Challenge number={9} title="Manage Files in Admin">
            <p>Upload 3 different files through your application — try an image (JPG/PNG), a
            document (PDF), and a spreadsheet (CSV or XLSX). Then open the admin panel and navigate
            to the Files page under the System section in the sidebar. Can you see all 3 files?
            Try sorting by size, searching by filename, and deleting one of them.</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-6">
            You{"'"}ve learned how Grit handles file storage — from the concepts of object storage and
            presigned URLs to hands-on uploading with MinIO and switching to cloud providers. Here{"'"}s
            what you covered:
          </p>

          <ul className="space-y-3 text-muted-foreground mb-8">
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Object storage</strong> — files stored outside the database in S3-compatible storage (MinIO, R2, AWS S3)</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Presigned URLs</strong> — time-limited signed URLs that let the frontend upload directly to S3, bypassing the API</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Storage configuration</strong> — 7 environment variables that control where and how files are stored</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">MinIO</strong> — a local S3-compatible server for development, with a web console at port 9001</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Upload flow</strong> — presign request, direct S3 upload with XHR progress, then confirm with the API</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Image processing</strong> — background jobs that generate thumbnail and medium-sized versions automatically</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Cloud switching</strong> — change .env variables to switch between MinIO, Cloudflare R2, and AWS S3 with zero code changes</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Upload model</strong> — tracks every file in the database with filename, key, content type, size, and user association</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Admin Files page</strong> — system page for viewing, searching, downloading, and deleting uploaded files</span>
            </li>
          </ul>

          <Challenge number={10} title="Presigned URL Quiz">
            <p>Answer these questions without looking back at the course:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>What HTTP method does the frontend use to upload a file to a presigned URL?</li>
              <li>Why do we use XHR instead of fetch for the upload step?</li>
              <li>What happens if you try to use a presigned URL after 15 minutes?</li>
              <li>What is the default MinIO endpoint in development?</li>
              <li>What R2_REGION value should you use for Cloudflare R2?</li>
            </ol>
          </Challenge>

          <Challenge number={11} title="Storage Driver Deep Dive">
            <p>Open <Code>apps/api/internal/storage/storage.go</Code> and read the{' '}
            <Code>NewStorage</Code> function. Then answer:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>What Go library does Grit use for S3 operations? (Hint: look at the imports)</li>
              <li>How is the <Code>Secure</Code> option set? What does it control?</li>
              <li>What happens if <Code>STORAGE_DRIVER</Code> is set to an unsupported value?</li>
            </ol>
          </Challenge>

          <Challenge number={12} title="Compare Cloud Providers">
            <p>Research the differences between Cloudflare R2, AWS S3, and Backblaze B2. For each
            provider, find out:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>How much does storage cost per GB/month?</li>
              <li>Are there egress (download) fees?</li>
              <li>Is it S3-compatible (can Grit use it without code changes)?</li>
            </ol>
          </Challenge>

          <Challenge number={13} title="Upload Lifecycle">
            <p>Trace the full lifecycle of a file upload by finding and reading these files in
            your project:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li><Code>apps/api/internal/handler/upload_handler.go</Code> — the Presign and ConfirmUpload handlers</li>
              <li><Code>apps/api/internal/storage/storage.go</Code> — the PresignPut method</li>
              <li><Code>apps/api/internal/worker/image_worker.go</Code> — the background image processor</li>
              <li><Code>apps/api/internal/models/upload.go</Code> — the Upload model</li>
            </ol>
            <p className="mt-2">Write down the complete journey of a file from the user{"'"}s browser
            to S3 storage to thumbnail generation to database record.</p>
          </Challenge>

          <Challenge number={14} title="Final Challenge: Build a Photo Gallery">
            <p>Put everything together. Generate a Photo resource:</p>
            <CodeBlock filename="Terminal">
{`grit generate resource Photo --fields "title:string,description:text:optional,image:string"`}
            </CodeBlock>
            <p className="mt-3">Then complete these tasks:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Open the admin panel and find the Photos page</li>
              <li>Upload 5 photos with different titles and descriptions</li>
              <li>Open GORM Studio at <Code>localhost:8080/studio</Code> and find the photos table — verify all 5 records exist</li>
              <li>Open the MinIO console at <Code>localhost:9003</Code> and browse the uploads bucket — find all 5 original images and their thumbnails</li>
              <li>Check the uploads table in GORM Studio — are there corresponding upload records for each photo?</li>
              <li>Delete one photo from the admin panel and verify the file is removed from both the database and MinIO</li>
            </ol>
          </Challenge>
        </section>

        {/* ═══ Footer ═══ */}
        <CourseFooter />

        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses/grit-web/admin-panel', label: 'Admin Panel' }}
            next={{ href: '/courses/grit-web/jobs-email', label: 'Jobs & Email' }}
          />
        </div>
      </main>
    </div>
  )
}
