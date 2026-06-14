import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Files don&apos;t belong in your database (rows balloon, backups
        get expensive). They belong in object storage — S3 or
        S3-compatible. Grit ships a storage battery that abstracts the
        provider, so MinIO locally and R2 or AWS in production feels
        the same to your code.
      </p>

      <h2>Why we need it</h2>
      <p>
        Three jobs object storage handles better than a DB:
      </p>
      <ul>
        <li>
          <strong>Big binary data</strong> — images, PDFs, videos.
          Cheap per-GB, infinite scale, no DB bloat.
        </li>
        <li>
          <strong>Direct serving</strong> — clients fetch files via
          CDN-backed URLs, not through your API.
        </li>
        <li>
          <strong>Signed URLs</strong> — give time-limited download
          access without sharing credentials.
        </li>
      </ul>

      <h2>Where it lives</h2>
      <pre className="not-prose text-xs leading-relaxed bg-bg-elevated border border-border rounded-lg p-4 overflow-x-auto">{`apps/api/internal/storage/
├── storage.go        ← Service: Upload, URL, Delete
├── handler.go        ← HTTP upload handler (multipart/form-data)
└── image.go          ← Resize, crop, thumbnail (uses imaging)`}</pre>

      <h2>How it&apos;s implemented</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/storage/storage.go (simplified)"
        code={`type Service struct {
  client *minio.Client    // works with AWS S3, R2, MinIO — they all speak S3
  bucket string
  baseURL string          // for serving (CDN domain in prod)
}

func New(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*Service, error) {
  c, err := minio.New(endpoint, &minio.Options{
    Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
    Secure: useSSL,
  })
  if err != nil { return nil, err }
  return &Service{client: c, bucket: bucket}, nil
}

func (s *Service) Upload(ctx context.Context, key string, r io.Reader, size int64, mime string) error {
  _, err := s.client.PutObject(ctx, s.bucket, key, r, size, minio.PutObjectOptions{
    ContentType: mime,
  })
  return err
}

func (s *Service) URL(key string) string {
  return s.baseURL + "/" + key
}

func (s *Service) Delete(ctx context.Context, key string) error {
  return s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
}

func (s *Service) SignedURL(ctx context.Context, key string, ttl time.Duration) (string, error) {
  u, err := s.client.PresignedGetObject(ctx, s.bucket, key, ttl, nil)
  if err != nil { return "", err }
  return u.String(), nil
}`}
      />
      <p>
        Same interface, three deploy targets:
      </p>
      <ul>
        <li>
          <strong>Local dev:</strong> MinIO (via{' '}
          <code>docker-compose</code>) — set <code>STORAGE_DRIVER=minio</code>.
        </li>
        <li>
          <strong>Big prod (default for most teams):</strong> AWS S3 —
          set <code>STORAGE_DRIVER=s3</code>. Leave{' '}
          <code>S3_ENDPOINT</code> empty so the SDK uses the regional
          default. If your API runs on EC2/ECS/Lambda with an IAM role,
          leave the access/secret keys empty too — the SDK auto-discovers
          role credentials.
        </li>
        <li>
          <strong>Cheap prod (zero egress fees):</strong> Cloudflare R2 —
          set <code>STORAGE_DRIVER=r2</code> and the{' '}
          <code>R2_*</code> block.
        </li>
        <li>
          <strong>Cheapest per-GB:</strong> Backblaze B2 — set{' '}
          <code>STORAGE_DRIVER=b2</code> and the <code>B2_*</code> block.
        </li>
      </ul>
      <p>
        Same code; only the env vars change.
      </p>

      <h2>The upload handler</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/storage/handler.go"
        code={`func (h *Handler) Upload(c *gin.Context) {
  fh, err := c.FormFile("file")
  if err != nil {
    c.JSON(400, gin.H{"error": gin.H{"code": "no_file"}})
    return
  }
  if fh.Size > 10*1024*1024 {  // 10 MB cap
    c.JSON(400, gin.H{"error": gin.H{"code": "too_large"}})
    return
  }
  src, err := fh.Open()
  if err != nil { c.JSON(500, gin.H{"error": gin.H{"code": "open_failed"}}); return }
  defer src.Close()

  key := uuid.NewString() + filepath.Ext(fh.Filename)
  mime := fh.Header.Get("Content-Type")

  if err := h.storage.Upload(c.Request.Context(), key, src, fh.Size, mime); err != nil {
    c.JSON(500, gin.H{"error": gin.H{"code": "upload_failed", "message": err.Error()}})
    return
  }

  // Record in DB so we can list / delete later
  upload := models.Upload{
    Key: key, OriginalName: fh.Filename, Size: fh.Size, Mime: mime,
    UserID: c.GetUint("user_id"),
  }
  h.db.Create(&upload)

  c.JSON(201, gin.H{"data": gin.H{
    "id":   upload.ID,
    "url":  h.storage.URL(key),
    "size": fh.Size,
  }})
}`}
      />
      <p>
        Three details worth absorbing:
      </p>
      <ul>
        <li>
          <strong>Random key</strong> (<code>uuid + ext</code>) — never
          trust the client filename. <code>../../../etc/passwd</code>{' '}
          is a real attack.
        </li>
        <li>
          <strong>Size cap</strong> upfront. Otherwise a malicious user
          can DOS your storage budget.
        </li>
        <li>
          <strong>DB record</strong>. The S3 object is the bytes; the
          <code>uploads</code> table is the metadata. You need both to
          let users list / delete.
        </li>
      </ul>

      <h2>The frontend upload</h2>
      <CodeBlock
        language="ts"
        code={`async function uploadFile(file: File) {
  const form = new FormData()
  form.append('file', file)
  const res = await fetch('/api/uploads', {
    method: 'POST',
    body: form,                // do NOT set Content-Type — let the browser set it
    credentials: 'include',
  })
  return res.json()
}`}
      />
      <p>
        Critical gotcha: <strong>don&apos;t set Content-Type
        manually</strong>. The browser must set it with the correct
        boundary string for multipart parsing. Override it and Gin
        can&apos;t read the file.
      </p>

      <h2>Image processing</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/storage/image.go (sketch)"
        code={`func (s *Service) UploadImageWithThumbnail(ctx context.Context, src io.Reader, key string) error {
  img, _, err := image.Decode(src)
  if err != nil { return err }

  // Full size
  var full bytes.Buffer
  imaging.Encode(&full, img, imaging.JPEG, imaging.JPEGQuality(85))
  if err := s.Upload(ctx, key, &full, int64(full.Len()), "image/jpeg"); err != nil { return err }

  // 200px thumbnail
  thumb := imaging.Resize(img, 200, 0, imaging.Lanczos)
  var tBuf bytes.Buffer
  imaging.Encode(&tBuf, thumb, imaging.JPEG, imaging.JPEGQuality(75))
  return s.Upload(ctx, "thumbs/" + key, &tBuf, int64(tBuf.Len()), "image/jpeg")
}`}
      />
      <p>
        Cheap thumbnails: resize once on upload, serve thumbs in lists,
        full size on detail pages. Saves the user&apos;s bandwidth.
      </p>

      <TipBox tone="warning">
        <strong>Set CORS on your bucket if you serve directly.</strong>{' '}
        If the frontend fetches images from{' '}
        <code>uploads.example.com</code>, that domain needs the right{' '}
        <code>Access-Control-Allow-Origin</code> on its responses.
        MinIO has it open in dev; R2/S3 you configure once.
      </TipBox>

      <h2>Signed URLs — when you need access control</h2>
      <p>
        Public files (avatars, blog images): just use{' '}
        <code>storage.URL(key)</code>. No auth, served via CDN.
      </p>
      <p>
        Private files (user&apos;s personal documents): the bucket
        should reject anonymous access. Use{' '}
        <code>storage.SignedURL(key, 5*time.Minute)</code> — that
        returns a long URL with a query-string signature. Valid for 5
        minutes, then dies. Hand the URL to the frontend; the user&apos;s
        browser fetches the file directly without going through your
        API.
      </p>

      <h2>How to modify this battery</h2>
      <ul>
        <li>
          <strong>Increase the size cap</strong> — edit the <code>fh.Size</code>{' '}
          check in <code>handler.go</code>.
        </li>
        <li>
          <strong>Restrict MIME types</strong> — add an allow-list:{' '}
          <code>allowedMimes := map[string]bool&#123;&quot;image/jpeg&quot;: true, ...&#125;</code>.
        </li>
        <li>
          <strong>Add virus scanning</strong> — pipe the upload through
          ClamAV before saving. Heavier; only do it if your threat model
          requires it.
        </li>
        <li>
          <strong>Per-user quota</strong> — sum{' '}
          <code>uploads.size</code> by user_id, reject if over.
        </li>
      </ul>

      <h2>Local dev wiring</h2>
      <CodeBlock
        language="yaml"
        filename="docker-compose.yml (excerpt)"
        code={`minio:
  image: minio/minio:latest
  command: server /data --console-address ":9001"
  environment:
    MINIO_ROOT_USER: minioadmin
    MINIO_ROOT_PASSWORD: minioadmin
  ports:
    - "9000:9000"   # API
    - "9001:9001"   # web console`}
      />
      <CodeBlock
        language="env"
        code={`STORAGE_ENDPOINT=localhost:9000
STORAGE_ACCESS_KEY=minioadmin
STORAGE_SECRET_KEY=minioadmin
STORAGE_BUCKET=uploads
STORAGE_USE_SSL=false`}
      />
      <p>
        For prod, swap to{' '}
        <code>STORAGE_ENDPOINT=&lt;account-id&gt;.r2.cloudflarestorage.com</code>{' '}
        + Cloudflare R2 credentials. Zero code changes.
      </p>

      <KnowledgeCheck
        question="A user uploads `report.pdf`. Where does the bytes go, and where does the metadata go?"
        choices={[
          {
            label: 'Both the bytes and metadata go into the database',
            feedback:
              "Bytes don't belong in the DB — backups explode, queries slow down. Storage handles bytes; DB handles metadata.",
          },
          {
            label: 'Bytes → S3/MinIO bucket; metadata (id, filename, size, mime, user_id) → uploads table in the DB',
            correct: true,
            feedback:
              "Right — separation of concerns. The bucket is cheap, scalable storage; the DB tracks who owns what, when, and lets you join with other tables.",
          },
          {
            label: 'Both go to S3',
            feedback:
              "Then you can't query &quot;all uploads by user X&quot; with SQL. The DB is your search/join layer.",
          },
          {
            label: 'Bytes go to local disk',
            feedback:
              "Local disk means every API server has its own files — disaster for horizontal scaling. S3-style storage is the standard.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Add file upload to the User profile (avatar):</p>
            <ol>
              <li>
                Add an <code>AvatarKey string</code> field to the User
                model.
              </li>
              <li>
                Add a route:{' '}
                <code>POST /api/me/avatar</code> that accepts a multipart
                file, uploads via the storage service, sets{' '}
                <code>user.AvatarKey</code>, deletes the old one if
                present.
              </li>
              <li>
                Return the URL via{' '}
                <code>storage.URL(user.AvatarKey)</code>.
              </li>
              <li>
                Display it in the admin user list.
              </li>
              <li>
                Try uploading a 20MB file — confirm you get 400 due to
                size cap.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            For testing the delete-old-on-replace, upload an avatar,
            then upload another, then look in the MinIO console at{' '}
            <code>localhost:9001</code> — there should be only one
            object for the user, not two.
          </>
        }
        solution={
          <>
            <p>
              You should have a working avatar upload that displays in
              the admin panel and cleans up old files. That&apos;s the
              storage battery doing what every SaaS needs.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — <strong>Mail (Resend)</strong>. Transactional
        email with editable templates — the most-needed external
        service after the DB.
      </p>
    </>
  )
}
