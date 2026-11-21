import Link from 'next/link'
import type { Metadata } from 'next'
import {
  GuidePage,
  Section,
  Step,
  H3,
  Outcome,
  Note,
  Warn,
  Mistake,
  Checklist,
  Faq,
  Related,
} from '@/components/guide'
import { CodeBlock } from '@/components/code-block'
import { Diagram, DiagramBox, DiagramArrow } from '@/components/diagram'

export const metadata: Metadata = {
  title: 'Upload files and images in Go to S3, R2 or MinIO',
  description:
    'Presigned uploads explained, why the browser should send the file directly to storage, and the validation that has to happen on the server anyway.',
  alternates: { canonical: 'https://gritframework.dev/guides/go-file-upload-s3-presigned' },
  openGraph: {
    title: 'Upload files and images in Go to S3, R2 or MinIO',
    description: 'The presigned flow, drawn out, with the security checks that matter.',
    url: 'https://gritframework.dev/guides/go-file-upload-s3-presigned',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="go-file-upload-s3-presigned"
      lead={
        <>
          <p className="mb-4">
            There are two ways to accept a file. The browser can send it to your API, which
            forwards it to storage, or your API can hand the browser a time-limited permission slip
            and let it upload directly. The second is almost always right, and the reason is not
            performance.
          </p>
          <p>
            Both are below, with the checks that are easy to skip and are the ones that matter: a
            file upload is the most direct route a stranger has into your infrastructure.
          </p>
        </>
      }
    >
      <Outcome>
        <p>Direct browser uploads to S3, Cloudflare R2, Backblaze B2 or MinIO, from Go.</p>
        <p>Validation that works even though your server never sees the bytes.</p>
        <p>Private files served through short-lived links rather than a public bucket.</p>
      </Outcome>

      <Section id="why" title="Why not just POST it to your API">
        <p>It works, and for small files it is fine. It stops being fine like this:</p>
        <Checklist
          items={[
            <>
              <strong>Your server holds the file while it uploads.</strong> A 200MB video over a
              slow connection occupies a request, a goroutine and memory or disk for several
              minutes. Twenty of those at once is your whole service.
            </>,
            <>
              <strong>Proxies have body limits.</strong> nginx defaults to 1MB and returns 413
              before your handler runs. So does almost every load balancer and platform, each with
              a different number.
            </>,
            <>
              <strong>A redeploy kills the upload.</strong> Restart during a ten-minute upload and
              the user starts again.
            </>,
            <>
              <strong>You pay for the bandwidth twice,</strong> once inbound and once outbound to
              storage.
            </>,
          ]}
        />
        <p>
          A presigned upload removes all four, because the bytes never touch your service. What it
          costs is that your server cannot inspect the file as it arrives, which is the thing the
          rest of this page is about.
        </p>
      </Section>

      <Section id="presigned" title="What a presigned URL actually is">
        <p>
          Not a secret link somebody generated. It is an ordinary storage URL with a signature in
          the query string, computed from your credentials, the exact operation, the exact object
          key and an expiry. Storage verifies the signature itself.
        </p>
        <CodeBlock
          language="text"
          code={`https://bucket.s3.eu-west-1.amazonaws.com/uploads/a1b2c3.jpg
  ?X-Amz-Algorithm=AWS4-HMAC-SHA256
  &X-Amz-Credential=AKIA.../20261121/eu-west-1/s3/aws4_request
  &X-Amz-Date=20261121T143022Z
  &X-Amz-Expires=300
  &X-Amz-Signature=5d41402abc4b2a76b9719d911017c592...`}
        />
        <p>
          Three consequences worth holding onto. The permission is for{' '}
          <strong>one operation on one key</strong>: a URL signed for{' '}
          <code>PUT uploads/a1b2c3.jpg</code> cannot write anywhere else. It expires, so a leaked
          URL is useless shortly afterwards. And your credentials are never sent to the browser,
          only a signature derived from them.
        </p>
        <Diagram>
          <DiagramBox title="1. Browser asks your API" sub="POST /uploads/presign" />
          <DiagramArrow label="filename, size, type" />
          <DiagramBox
            title="2. API checks and signs"
            sub="authorise, validate, generate a key"
            tone="primary"
          />
          <DiagramArrow label="url + key" />
          <DiagramBox title="3. Browser PUTs the file to storage" sub="your API is not involved" />
          <DiagramArrow label="key" />
          <DiagramBox title="4. Browser tells your API it is done" sub="POST /uploads/confirm" tone="primary" />
        </Diagram>
        <p>
          Step four is the one people leave out, and without it you have files in a bucket that
          nothing in your database knows about.
        </p>
      </Section>

      <Step n={1} title="Sign the URL">
        <CodeBlock
          filename="internal/services/storage.go"
          code={`package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Only what we are prepared to accept. An allow-list, not a block-list:
// a block-list is a list of the attacks somebody already thought of.
var allowedTypes = map[string]string{
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"image/webp":      ".webp",
	"application/pdf": ".pdf",
}

const maxUploadBytes = 10 << 20 // 10 MiB

func (s *StorageService) PresignUpload(
	ctx context.Context, userID uint, filename, contentType string, size int64,
) (url, key string, err error) {
	ext, ok := allowedTypes[contentType]
	if !ok {
		return "", "", ErrUnsupportedType
	}
	if size <= 0 || size > maxUploadBytes {
		return "", "", ErrFileTooLarge
	}

	// Generate the key. NEVER build it from the filename the client sent:
	// "../../etc/passwd" and "index.html" are both filenames, and the
	// second one is how a bucket starts serving somebody else's page.
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", "", fmt.Errorf("generate key: %w", err)
	}
	key = fmt.Sprintf("uploads/%d/%s%s", userID, hex.EncodeToString(buf), ext)

	req, err := s.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
		// Binding the content type into the signature means the upload is
		// rejected if the browser sends something else.
		ContentType:   &contentType,
		ContentLength: &size,
	}, s3.WithPresignExpires(5*time.Minute))
	if err != nil {
		return "", "", fmt.Errorf("presign: %w", err)
	}

	return req.URL, key, nil
}`}
        />
        <p>
          The original filename is still worth keeping, as a <em>column</em>, so the download can
          be called <code>invoice-july.pdf</code>. It is just never part of the path.
        </p>
        <Warn title="Everything the client tells you is a claim">
          <p>
            <code>Content-Type</code> is set by the browser from the file extension, and by{' '}
            <code>curl</code> from whatever you type. <code>size</code> is a number in a JSON body.
            Checking them is worth doing, because it catches honest mistakes early, and it is not
            security. The real check happens after the bytes exist.
          </p>
        </Warn>
      </Step>

      <Step n={2} title="Upload from the browser">
        <CodeBlock
          language="ts"
          code={`async function upload(file: File) {
  // 1. Ask for permission
  const { data } = await api<{ data: { url: string; key: string } }>(
    '/uploads/presign',
    {
      method: 'POST',
      body: JSON.stringify({
        filename: file.name,
        content_type: file.type,
        size: file.size,
      }),
    },
  )

  // 2. Send the file straight to storage. Note: no credentials, no cookies,
  // and the Content-Type must match what was signed, exactly.
  const put = await fetch(data.url, {
    method: 'PUT',
    body: file,
    headers: { 'Content-Type': file.type },
  })
  if (!put.ok) throw new Error('upload failed')

  // 3. Tell the API the object exists
  return api('/uploads/confirm', {
    method: 'POST',
    body: JSON.stringify({ key: data.key, filename: file.name }),
  })
}`}
        />
        <Note title="A progress bar needs XMLHttpRequest">
          <p>
            <code>fetch</code> has no upload progress event. For a visible percentage, which is
            worth having for anything over a few megabytes, use <code>XMLHttpRequest</code> and its{' '}
            <code>upload.onprogress</code>. It is the one place the older API is still the right
            tool.
          </p>
        </Note>
      </Step>

      <Step n={3} title="Confirm, and verify what actually arrived">
        <p>
          This is the step that makes the flow safe. Your server did not see the bytes, so it asks
          storage what it got.
        </p>
        <CodeBlock
          filename="internal/handlers/upload.go"
          code={`func (h *UploadHandler) Confirm(c *gin.Context) {
	var req ConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.ValidationError(c, err)
		return
	}
	userID := c.GetUint("user_id")

	// The key must be one WE issued to THIS user. Without this check, a
	// caller can confirm any key they can guess, and attach somebody
	// else's file to their own record.
	if !strings.HasPrefix(req.Key, fmt.Sprintf("uploads/%d/", userID)) {
		respond.Fail(c, http.StatusForbidden, "FORBIDDEN", "not your upload")
		return
	}

	// Ask storage what is actually there: the real size and the real type,
	// rather than what the client said five minutes ago.
	head, err := h.storage.Head(c.Request.Context(), req.Key)
	if err != nil {
		respond.Fail(c, http.StatusBadRequest, "UPLOAD_MISSING", "no file at that key")
		return
	}
	if head.ContentLength > maxUploadBytes {
		_ = h.storage.Delete(c.Request.Context(), req.Key)
		respond.Fail(c, http.StatusRequestEntityTooLarge, "TOO_LARGE", "file is too large")
		return
	}

	upload, err := h.uploads.Record(c.Request.Context(), userID, req.Key,
		req.Filename, head.ContentLength, head.ContentType)
	if err != nil {
		respond.Error(c, err)
		return
	}
	respond.Created(c, upload)
}`}
        />
        <H3>Check the bytes, not the extension</H3>
        <p>
          A file called <code>photo.jpg</code> with <code>Content-Type: image/jpeg</code> can
          contain anything at all. The only honest check reads the first few bytes:
        </p>
        <CodeBlock
          code={`// The first 512 bytes are enough for http.DetectContentType, so download
// only those rather than the whole file.
head := make([]byte, 512)
if _, err := io.ReadFull(obj.Body, head); err != nil && err != io.ErrUnexpectedEOF {
	return fmt.Errorf("read header: %w", err)
}
if !allowed[http.DetectContentType(head)] {
	_ = s.Delete(ctx, key)
	return ErrUnsupportedType
}`}
        />
        <p>
          For images, decoding is a stronger check still, and you were going to decode anyway to
          make thumbnails. A file that is not a real image fails to decode, which answers the
          question completely.
        </p>
      </Step>

      <Section id="serving" title="Serving them back">
        <p>
          A public bucket is right for avatars and product photos: they are cacheable, they are
          served by the CDN, and your service is not in the path.
        </p>
        <p>
          It is wrong for anything a user uploaded that another user should not see. Private bucket
          plus a short-lived download URL:
        </p>
        <CodeBlock
          code={`func (h *UploadHandler) Download(c *gin.Context) {
	upload, err := h.uploads.ByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		respond.Error(c, err)
		return
	}
	// The authorisation check happens HERE, in your code, which is the
	// entire reason the bucket is private.
	if upload.UserID != c.GetUint("user_id") {
		// 404, not 403: a 403 confirms the file exists.
		respond.Fail(c, http.StatusNotFound, "NOT_FOUND", "not found")
		return
	}

	url, err := h.storage.PresignGet(c.Request.Context(), upload.Key, 5*time.Minute,
		// Makes the browser save it under its original name rather than
		// the random key.
		"attachment; filename=\\""+upload.Filename+"\\"")
	if err != nil {
		respond.Error(c, err)
		return
	}
	c.Redirect(http.StatusFound, url)
}`}
        />
        <Warn title="Guessable keys are not access control">
          <p>
            A random key makes a URL hard to guess, which is worth something and is not a
            permission. Links get pasted into chat, forwarded, and indexed when somebody puts one
            on a public page. If the answer to &quot;may this person see this file&quot; is ever
            no, the bucket is private and your code decides.
          </p>
        </Warn>
      </Section>

      <Section id="images" title="Images need a second pass">
        <p>
          Serving the 4000 by 3000 photograph somebody took on their phone as a 48-pixel avatar
          wastes their bandwidth and your money, on every page load.
        </p>
        <CodeBlock
          code={`// In a background job, not in the request: resizing is CPU work and
// nobody should wait for it.
func (w *ImageWorker) MakeRenditions(ctx context.Context, key string) error {
	src, err := w.storage.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("get original: %w", err)
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		// Not an image, whatever it was called. That is now known.
		return fmt.Errorf("decode %s: %w", key, err)
	}

	for _, size := range []int{128, 512, 1280} {
		resized := imaging.Fit(img, size, size, imaging.Lanczos)
		var buf bytes.Buffer
		if err := webp.Encode(&buf, resized, nil); err != nil {
			return fmt.Errorf("encode %d: %w", size, err)
		}
		if err := w.storage.Put(ctx, renditionKey(key, size), &buf, "image/webp"); err != nil {
			return fmt.Errorf("store %d: %w", size, err)
		}
	}
	return nil
}`}
        />
        <Note title="A stored file is several objects">
          <p>
            Once there are renditions, deleting &quot;the upload&quot; means deleting the original
            and every derived file. Keep the list, on the record, rather than reconstructing it
            from a naming rule: the rule will change, and the old objects will stay in the bucket
            costing money forever.
          </p>
        </Note>
      </Section>

      <Section id="providers" title="S3, R2, B2 or MinIO">
        <p>
          All four speak the S3 API, so the same SDK and the same code work against each. What
          differs is the bill and the endpoint.
        </p>
        <Checklist
          items={[
            <>
              <strong>AWS S3.</strong> The reference. Egress is charged per gigabyte and is what
              makes an image-heavy product expensive.
            </>,
            <>
              <strong>Cloudflare R2.</strong> Same API, no egress charge. For anything serving
              media this is usually the decision that matters.
            </>,
            <>
              <strong>Backblaze B2.</strong> Cheap storage, free egress through the Cloudflare
              partnership.
            </>,
            <>
              <strong>MinIO.</strong> Runs on your machine in a container. For development, so
              that nothing is uploaded to a real bucket from a laptop.
            </>,
          ]}
        />
        <CodeBlock
          language="bash"
          code={`# One set of variables, four providers.
STORAGE_ENDPOINT=https://ACCOUNT.r2.cloudflarestorage.com   # R2
# STORAGE_ENDPOINT=http://localhost:9002                     # MinIO locally
STORAGE_REGION=auto
STORAGE_BUCKET=myapp
STORAGE_ACCESS_KEY=...
STORAGE_SECRET_KEY=...`}
        />
        <p>
          Non-AWS endpoints generally need path-style addressing, which is one option on the
          client and the cause of most &quot;works against S3, fails against MinIO&quot; reports:
        </p>
        <CodeBlock
          code={`s3.NewFromConfig(cfg, func(o *s3.Options) {
	o.BaseEndpoint = aws.String(endpoint)
	o.UsePathStyle = true   // bucket in the path, not in the hostname
})`}
        />
      </Section>

      <Section id="mistakes" title="When it goes wrong">
        <Mistake error="SignatureDoesNotMatch">
          <p>
            The request did not match what was signed. Almost always a header: a{' '}
            <code>Content-Type</code> the browser changed, or an extra header added by a fetch
            wrapper. Send exactly the headers that were signed and nothing else. Also check the
            clock: a server more than fifteen minutes out produces signatures everything rejects.
          </p>
        </Mistake>

        <Mistake error="403 Forbidden on the PUT, with a valid URL">
          <p>
            Usually the URL expired, which is five minutes by default and is easy to exceed while
            debugging. Otherwise the credentials lack <code>s3:PutObject</code> on that prefix: a
            presigned URL carries your permissions, so it can never grant more than the signer
            has.
          </p>
        </Mistake>

        <Mistake error="CORS error on the PUT to the storage endpoint">
          <p>
            The bucket has its own CORS configuration, separate from your API&apos;s. It must
            allow <code>PUT</code> from your origin.
          </p>
          <CodeBlock
            language="json"
            code={`[{
  "AllowedOrigins": ["https://yourdomain.com"],
  "AllowedMethods": ["PUT", "GET"],
  "AllowedHeaders": ["content-type"],
  "ExposeHeaders": ["etag"],
  "MaxAgeSeconds": 3000
}]`}
          />
        </Mistake>

        <Mistake error="413 Request Entity Too Large">
          <p>
            A proxy rejected the body before your handler saw it. If you are uploading through
            your API, raise the limit on nginx with <code>client_max_body_size</code>, or whatever
            your platform calls it. If you are using presigned URLs and still see this, you are
            not using them for this particular upload.
          </p>
        </Mistake>

        <Mistake error="The image shows locally and not on a phone">
          <p>
            Stored URLs contain <code>localhost</code>, which on a phone means the phone. Store the
            key and build the URL at read time from configuration, so one record works from the
            browser, the emulator and a real device. If you want MinIO reachable from a phone on
            the same network, publish it on the machine&apos;s LAN address rather than{' '}
            <code>127.0.0.1</code>.
          </p>
        </Mistake>
      </Section>

      <Section id="with-grit" title="As a field type">
        <p>Grit treats this as a property of a field rather than a subsystem to assemble:</p>
        <CodeBlock
          terminal
          code={`grit generate resource Product \\
  --fields "name:string:required,photo:image,manual:file,gallery:images"`}
        />
        <p>
          That writes the presign and confirm endpoints, the type and size checks, the upload
          component with its progress and drag-and-drop, the background job that makes renditions,
          and the admin preview. <code>STORAGE_DRIVER</code> selects S3, R2, B2, MinIO or the local
          disk, and local disk means a new project uploads files before you have any cloud account
          at all.
        </p>
        <p>
          <Link href="/docs/batteries/storage" className="text-accent hover:underline">
            The file storage reference
          </Link>{' '}
          covers the drivers and the field types.
        </p>
      </Section>

      <Faq
        slug="go-file-upload-s3-presigned"
        items={[
          {
            q: 'Presigned URLs or upload through my API?',
            a: (
              <>
                Presigned for anything a user chooses from their device, which is nearly
                everything. Through your API when the file is small and you must inspect or
                transform it before it is stored at all, or when you are generating the file
                yourself server-side, where there is no browser involved.
              </>
            ),
            text:
              'Presigned for anything a user picks from their device, which is most cases. Through your API when the file is small and must be inspected or transformed before being stored, or when the server generates the file itself.',
          },
          {
            q: 'How do I stop someone uploading something enormous?',
            a: (
              <>
                Bind <code>ContentLength</code> into the signature, which makes storage itself
                reject a larger body. Then check the real size at confirm time and delete the
                object if it is wrong. The client-side check is a courtesy to honest users; the
                signature is the enforcement.
              </>
            ),
            text:
              'Bind ContentLength into the signature so storage rejects a larger body, and check the real size at confirm time, deleting the object if it is wrong. The client-side check is a courtesy to honest users; the signature is the enforcement.',
          },
          {
            q: 'Is a random object key enough to keep a file private?',
            a: (
              <>
                No. It makes the URL unguessable, which is worth something, and it is not
                permission: links are pasted into chats, forwarded and occasionally indexed. If
                the answer to &quot;may this person see this&quot; can be no, the bucket is private
                and your handler decides on each request.
              </>
            ),
            text:
              'No. A random key makes the URL unguessable, which has some value, but it is not permission: links get pasted, forwarded and sometimes indexed. If access can ever be denied, the bucket must be private and the handler must decide per request.',
          },
          {
            q: 'Where should thumbnails be generated?',
            a: (
              <>
                In a background job. Resizing a large photograph takes real CPU time, and doing it
                in the request means the user waits and your handler holds a worker. Record the
                upload immediately, queue the work, and let the interface show the original or a
                placeholder until the renditions appear.
              </>
            ),
            text:
              'In a background job. Resizing a large photo takes real CPU time, so doing it in the request makes the user wait and ties up a worker. Record the upload, queue the work, and show the original or a placeholder until renditions exist.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/go-background-jobs-asynq',
            label: 'Background jobs in Go',
            blurb: 'Where the image processing belongs, and what a retry costs.',
          },
          {
            href: '/guides/docker-compose-postgres-redis',
            label: 'MinIO locally with Docker Compose',
            blurb: 'A local S3 so nothing is uploaded to a real bucket from a laptop.',
          },
          {
            href: '/docs/batteries/storage',
            label: 'File storage reference',
            blurb: 'The drivers, the field types and the generated endpoints.',
          },
          {
            href: '/docs/concepts/field-types',
            label: 'Field types',
            blurb: 'What file, image and images each produce.',
          },
        ]}
      />
    </GuidePage>
  )
}
