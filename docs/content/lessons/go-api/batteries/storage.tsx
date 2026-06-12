import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        File uploads — avatars, attachments, exports. Grit&apos;s storage
        package abstracts S3-compatible backends: <strong>MinIO</strong> in
        dev (ships with docker-compose), <strong>Cloudflare R2</strong> or{' '}
        <strong>AWS S3</strong> in production. One interface, swap the
        backend via env.
      </p>

      <h2>The Storage interface</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/storage/storage.go"
        code={`type Storage interface {
    Put(ctx context.Context, key string, body io.Reader, contentType string) (string, error)
    Get(ctx context.Context, key string) (io.ReadCloser, error)
    Delete(ctx context.Context, key string) error
    PresignedURL(ctx context.Context, key string, expires time.Duration) (string, error)
}`}
      />
      <p>
        Same four operations regardless of backend. The handler doesn&apos;t
        know if it&apos;s talking to MinIO or R2.
      </p>

      <h2>Switching backends</h2>
      <CodeBlock
        language="dotenv"
        filename=".env"
        code={`# Choose backend
STORAGE_DRIVER=minio        # minio | r2 | s3

# MinIO (dev — bundled in docker-compose)
MINIO_ENDPOINT=http://localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=bench-api-uploads

# Cloudflare R2 (prod)
R2_ACCOUNT_ID=...
R2_ACCESS_KEY=...
R2_SECRET_KEY=...
R2_BUCKET=acme-prod-uploads`}
      />
      <p>
        Change <code>STORAGE_DRIVER</code>, restart. Same handler code,
        different bucket.
      </p>

      <h2>The Upload handler</h2>
      <p>Grit ships a generic upload handler at <code>POST /api/uploads</code>:</p>
      <CodeBlock
        language="text"
        code={`POST /api/uploads
Headers: Authorization: Bearer <token>
Body:    multipart/form-data with file=<file>

Response:
{ "data": { "id": "9b...", "url": "/api/uploads/9b...", "size": 24832 } }`}
      />
      <p>
        Behind the scenes: file goes into the storage backend keyed by its
        SHA-256 hash (deduplication by content), an Upload row is created
        in the DB. The returned URL is a Grit proxy URL that streams the
        file with auth + access checks.
      </p>

      <h2>Direct uploads (presigned URLs) — when files get large</h2>
      <p>
        For files &gt;5 MB, you don&apos;t want them to pass through your
        API process. Use presigned URLs:
      </p>
      <CodeBlock
        language="go"
        code={`// API generates a one-time URL the client uploads directly to
url, err := h.storage.PresignedURL(ctx, key, 15*time.Minute)
// Return url to the client; they PUT the file directly to S3/R2`}
      />
      <p>
        Bandwidth doesn&apos;t hit your API. Useful for image uploads, video,
        backups, exports.
      </p>

      <TipBox tone="warning">
        <strong>Don&apos;t accept user URLs to fetch (SSRF)!</strong> A common
        bug: &quot;upload from URL&quot; lets users send a URL and your server
        fetches it. Without filtering, attackers can target your AWS metadata
        endpoint or internal services. Use{' '}
        <code>internal/safefetch.Get(url)</code> instead — it blocks
        loopback, RFC1918, and cloud-metadata IPs. The Defender&apos;s Handbook
        page covers SSRF in depth.
      </TipBox>

      <h2>Content type validation</h2>
      <p>
        The handler validates content-type from the uploaded bytes (sniffing),
        not from the user-supplied <code>Content-Type</code> header. An
        attacker who renames <code>shell.php</code> to <code>image.png</code>{' '}
        gets caught — the sniffer sees PHP source, the handler rejects.
      </p>

      <h2>Image processing</h2>
      <p>
        For avatars and thumbnails, Grit ships an image processing pipeline:
      </p>
      <CodeBlock
        language="go"
        code={`processed, err := h.imageProcessor.Resize(ctx, body, 256, 256)
// Outputs JPEG, 256x256, ~20 KB
storage.Put(ctx, key, processed, "image/jpeg")`}
      />

      <KnowledgeCheck
        question="A user submits a 50 MB video to /api/uploads. What's the right pattern?"
        choices={[
          {
            label: 'Stream through the API to storage — same handler, just bigger',
            feedback:
              "Possible but ties up a request worker for the whole upload, eats API bandwidth, and risks proxy timeouts. The cleaner pattern is presigned URLs.",
          },
          {
            label: 'Reject large uploads — set a max body size',
            feedback:
              "Defensible if your product genuinely doesn't accept large files. But if it does, you need a real solution.",
          },
          {
            label: 'Return a presigned URL; the client uploads directly to R2/S3',
            correct: true,
            feedback:
              "Right — for large files, bypass the API. PresignedURL gives the client a one-time URL with a short expiry; the file goes straight to storage. API stays fast.",
          },
          {
            label: 'Background-job the upload after returning 200 to the client',
            feedback:
              "Doesn't help — the client still uploaded through the API. Background-jobbing only delays the storage write, not the slow transfer.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Upload a real file and read it back:
            </p>
            <ol>
              <li>
                With your access token, POST a file (a small PNG) to{' '}
                <code>/api/uploads</code> with curl:
                <CodeBlock
                  terminal
                  code={`curl -X POST http://localhost:8080/api/uploads \\
  -H "Authorization: Bearer YOUR_TOKEN" \\
  -F "file=@avatar.png"`}
                />
              </li>
              <li>Note the <code>url</code> in the response.</li>
              <li>
                Open MinIO console (<code>http://localhost:9001</code>, login{' '}
                <code>minioadmin/minioadmin</code>) — you should see the
                file in the bucket.
              </li>
              <li>
                Fetch the file via the API URL with your token — should
                stream the PNG back.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            If MinIO console says &quot;bucket not found&quot;, run{' '}
            <code>docker compose restart minio</code> and try again — Grit
            auto-creates the bucket on startup.
          </>
        }
        solution={
          <>
            <p>The upload response shape:</p>
            <CodeBlock
              language="json"
              code={`{
  "data": {
    "id": "9b4d-...",
    "key": "uploads/sha256/abc123...",
    "url": "/api/uploads/9b4d-...",
    "size": 24832,
    "content_type": "image/png"
  }
}`}
            />
            <p>
              File is in MinIO under the content-hashed key, indexed in the
              DB by ID, accessible via a clean Grit URL.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Last battery — <strong>AI Gateway</strong>. Stream from Claude,
        GPT-4, and ~98 other models with one API key. Built for products
        that wrap AI as a feature.
      </p>
    </>
  )
}
