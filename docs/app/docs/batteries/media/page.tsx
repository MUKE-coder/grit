import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Diagram, DiagramBox, DiagramRow, DiagramArrow } from '@/components/diagram'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/batteries/media')

const C = 'text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded'

export default function MediaPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Batteries</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Image Optimisation</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Somebody uploads a 6 MB photograph straight off a phone. What you want stored
                is about 150 KB, correctly oriented, with the GPS coordinates removed and a
                thumbnail beside it. That happens by default, before the file is written, with
                no configuration at all.
              </p>
            </div>

            <div className="prose-grit">
              <h2 id="default">You do not have to configure anything</h2>
              <p>
                Every image upload goes through the pipeline using{' '}
                <code className={C}>DefaultProfile()</code>. Measured on a real 6.08 MB camera
                photo:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                language="bash"
                code={`6.08 MB  4000x3000 JPEG        the upload
 141 KB   1600x1200 JPEG q82    what gets stored
   9 KB    400x400  JPEG        the thumbnail, alongside it
------------------------------------------------------------
41x smaller, and the 6 MB file never reaches your public bucket`}
              />
            </div>

            <div className="prose-grit">
              <p>The defaults, and why each number is what it is:</p>
              <ul>
                <li>
                  <strong>Fit inside 1600x1600.</strong> Covers a 2x retina display at a
                  typical content width. Storing more is paying to keep pixels no browser will
                  draw.
                </li>
                <li>
                  <strong>Quality 0.82.</strong> The point at which JPEG is visually
                  indistinguishable from the source. Below about 0.75, artifacts appear on
                  gradients and skin.
                </li>
                <li>
                  <strong>A 400x400 thumbnail.</strong> What an admin table row or a card grid
                  needs at 2x.
                </li>
                <li>
                  <strong>The original is kept</strong>, under a private prefix, so a profile
                  change can be replayed later. Private because it is for reprocessing, not
                  for serving.
                </li>
                <li>
                  <strong>EXIF is oriented, then stripped.</strong> The one nobody thinks of.
                  Orientation first, or a portrait photo comes out sideways. Stripping second,
                  because a phone photo carries GPS coordinates, and a shop publishing product
                  photos would otherwise publish the seller&apos;s home address with them.
                </li>
              </ul>

              <h2 id="format">The format is chosen per image, not configured</h2>
              <p>
                Asking a developer to pick an output format is asking them to get it wrong
                once. The decision is made from the pixels:
              </p>
            </div>

            <Diagram>
              <DiagramRow>
                <DiagramBox
                  title="decode + auto-orient"
                  sub="EXIF applied, then dropped"
                  tone="violet"
                />
              </DiagramRow>
              <DiagramArrow label="resize into the profile's box, never upscaling" />
              <DiagramRow>
                <DiagramBox title="has alpha?" sub="any pixel not opaque" tone="cyan" />
              </DiagramRow>
              <DiagramArrow />
              <DiagramRow>
                <DiagramBox
                  title="WebP lossless"
                  sub="yes: transparency survives"
                  tone="primary"
                />
                <DiagramBox
                  title="JPEG at quality"
                  sub="no: ~20x smaller than lossless"
                  tone="green"
                />
              </DiagramRow>
            </Diagram>

            <div className="prose-grit">
              <p>
                The failure this prevents is a transparent logo encoded as JPEG, which
                silently gains a black box behind it. Under <code className={C}>Auto</code>{' '}
                that cannot happen.
              </p>
              <p>
                It also explains what the default backend will not do.{' '}
                <strong>No lossy WebP and no AVIF</strong>, because there is no pure-Go
                encoder for either. Measured on a photograph, pure-Go lossless WebP produced
                778 KB where JPEG q82 produced 35 KB: it is not a substitute for lossy
                encoding, it is a PNG replacement, which is exactly the job it is given here.
                Both formats are available on the libvips backend below.
              </p>

              <h2 id="backends">Two backends, and what swapping costs</h2>
              <p>
                The pipeline has a swappable backend. The default needs no system
                libraries; build with <code className={C}>-tags vips</code> and it uses
                libvips instead. Measured on the same 6.08 MB photograph, in the same
                container:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                language="bash"
                code={`                        pure Go (default)        libvips (-tags vips)
default profile         149.9 KB  JPEG  1001ms    33.9 KB  lossy WebP  1853ms
JPEG, forced            141.0 KB        950ms    123.7 KB              1366ms
AVIF                    downgrades to JPEG        79.9 KB              7544ms`}
              />
            </div>

            <div className="prose-grit">
              <p>
                <strong>libvips produces files about 4x smaller and is not faster.</strong>{' '}
                It is slower here, on both the default path and a like-for-like JPEG
                comparison. The often-quoted 4-8x speedup is libvips against ImageMagick,
                not against Go&apos;s native image package, and it did not reproduce. The
                reason to want libvips is bandwidth and storage, which is the thing that
                actually costs money.
              </p>
              <p>
                Note what AVIF costs: 7.5 seconds for one image, and on this photograph it
                came out <em>larger</em> than lossy WebP. It is worth having for the cases
                where it wins, but it is not a default and it is not viable on a synchronous
                upload.
              </p>
              <p>
                The catch is cgo. <code className={C}>grit deploy</code> cross-compiles to
                linux/amd64 with <code className={C}>CGO_ENABLED=0</code> from whatever
                machine you run it on, and cgo cannot cross-compile without a target
                toolchain, so the vips build has to happen where it will run. Docker is that
                place:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                language="bash"
                code={`docker build --build-arg IMAGE_BACKEND=vips -f Dockerfile.api .`}
              />
            </div>

            <div className="prose-grit">
              <p>
                The same profiles drive both backends. What changes is what{' '}
                <code className={C}>Auto</code> resolves to: lossy WebP under libvips, JPEG
                or lossless WebP under pure Go. A profile asking for AVIF on the pure-Go
                backend is downgraded rather than refused, so one binary still serves a
                project whose profiles assume libvips, and the ref records what was really
                produced. That is what <code className={C}>format</code> on the FileRef is
                for, and the backend is named in the upload log line.
              </p>

              <h2 id="profiles">Profiles, when a field wants something different</h2>
              <p>
                Profiles live in <code className={C}>internal/media/profiles.go</code>, a file
                written once and never regenerated, so what you put there survives{' '}
                <code className={C}>grit generate</code> and{' '}
                <code className={C}>grit upgrade</code>.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                filename="apps/api/internal/media/profiles.go"
                code={`func init() {
    media.Define("product-image", media.Profile{
        Max:     media.Fit(1000, 1000),
        Quality: 0.8,
        Renditions: map[string]media.Size{
            "thumb": media.Fill(300, 300),
            "card":  media.Fit(600, 600),
        },
    })

    media.Define("avatar", media.Profile{
        // Fill, not Fit: a portrait shown in a round frame should crop,
        // not letterbox.
        Max:             media.Fill(400, 400),
        Renditions:      map[string]media.Size{"thumb": media.Fill(80, 80)},
        DiscardOriginal: true,
    })
}`}
              />
            </div>

            <div className="prose-grit">
              <p>Name it from the upload:</p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="bash" code={`POST /api/v1/uploads?profile=product-image`} />
            </div>

            <div className="prose-grit">
              <p>
                A field you leave zero keeps the default, so a profile that only wants a
                different size says only that. An unknown or misspelled name falls back to the
                default rather than failing the upload, because a stale profile name in a
                deployed client build should degrade, not break.
              </p>
              <p>
                Note <code className={C}>DiscardOriginal</code> rather than a{' '}
                <code className={C}>KeepOriginal</code> that defaults to true. A Go bool cannot
                distinguish false from not-set, so a keep-flag would have silently discarded
                originals for every profile that did not mention it, while the documentation
                promised the opposite. The negative phrasing makes the zero value the
                recommended behaviour.
              </p>

              <h2 id="ref">What lands on the record</h2>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                language="json"
                code={`{
  "url": ".../uploads/2026/08/photo-178761052.jpg",
  "name": "big-photo.jpg",
  "mime": "image/jpeg",
  "size": 144333,
  "width": 1600, "height": 1200,
  "format": "jpeg",
  "optimised": true,
  "thumbnail_url": ".../uploads/2026/08/photo-178761052-thumb.jpg",
  "original_key": "originals/2026/08/178761052-big-photo.jpg",
  "original_size": 6375170,
  "renditions": {
    "thumb": { "url": "...", "width": 400, "height": 400, "size": 9135 }
  }
}`}
              />
            </div>

            <div className="prose-grit">
              <p>
                <code className={C}>format</code> is recorded rather than inferred from the
                URL, so a client never has to guess what it actually received.{' '}
                <code className={C}>optimised</code> is false when the transform failed and the
                file was stored as it arrived, which is the default failure policy: losing
                somebody&apos;s upload because an encoder choked is worse than storing a large
                file. Set <code className={C}>OnError: media.Reject</code> on a profile to
                refuse those instead.
              </p>

              <h2 id="client">Client-side, which is where it belongs</h2>
              <p>
                Uploads go from the browser straight to storage through a presigned URL and
                never pass through the API. So the optimisation happens on the client, in{' '}
                <code className={C}>@repo/upload</code>, before the bytes leave the device.
                That is not a compromise. Measured in real Chromium on the same 5 MB
                photograph:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                language="bash"
                code={`pure-Go server backend      149.9 KB
libvips server backend       33.9 KB   (needs cgo)
browser, client-side         35.0 KB   <- 147x smaller, no server involved`}
              />
            </div>

            <div className="prose-grit">
              <p>
                The browser matches libvips, because it has a lossy WebP encoder built in.
                That is the one thing pure Go could not do without cgo, and it turns out to
                have been on the client the whole time. Nothing is spent on server CPU or
                server bandwidth, and on a phone the 5 MB never leaves the handset.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                filename="apps/admin/lib/uploader.ts"
                code={`import { createUploader, createAxiosTransport } from "@repo/upload";
import { optimizeImage } from "@repo/upload/web";
import { apiClient } from "@/lib/api-client";

export const uploader = createUploader({
  transport: createAxiosTransport(apiClient),
  optimize: optimizeImage,
});`}
              />
            </div>

            <div className="prose-grit">
              <p>
                The optimiser is injected rather than imported, so the same uploader works on
                Next.js, a Vite SPA and Expo without any bundler resolving a platform. Expo
                swaps in <code className={C}>@repo/upload/expo</code>, which uses
                expo-image-manipulator because React Native has no canvas. React apps get{' '}
                <code className={C}>useUpload</code> from{' '}
                <code className={C}>@repo/upload/react</code>, with per-file progress and{' '}
                <code className={C}>describeSaving()</code> for the &quot;6.1 MB to 41 KB&quot;
                label.
              </p>
              <p>
                Profiles come from <code className={C}>GET /api/v1/media/profiles</code>, so
                the client uses the server&apos;s numbers rather than its own copy of them.
                Two copies drift the first time one changes.
              </p>

              <h2 id="trust">What the server does once it stops doing the work</h2>
              <p>
                A presigned URL is a capability handed to a browser, so the server can no
                longer guarantee what landed in the bucket. Its job becomes constraining and
                verifying rather than transforming, and two things do that.
              </p>
              <p>
                <strong>The exact byte count is signed into the URL.</strong> The client
                optimises first, so it knows the size before it asks, and S3 rejects a PUT of
                any other length. Without that the URL is an unbounded write capability: ask
                to upload two megabytes, send five gigabytes, and nothing on the server side
                ever sees it happen.
              </p>
              <p>
                <strong>The completion call re-reads the object from storage.</strong> Every
                number in that request is a claim, since the bytes never came through the
                API. Believing the reported size would make every storage total in the admin
                fiction. It now asks the bucket, and deletes anything over the limit.
              </p>

              <h2 id="sync">The server pipeline is still there</h2>
              <p>
                The transform is synchronous, which costs roughly half a second to two seconds
                on a large photograph, most of it decoding and resampling rather than encoding.
                That cost buys two things.
              </p>
              <p>
                The 6 MB file never lands in your public bucket, which is the entire point. And
                the record is only ever written with final URLs.
              </p>
              <p>
                The version this replaces did it the other way around. It stored the original,
                queued a job, and returned a reference whose thumbnail field was still empty
                because the worker had not run yet. That reference is what got written into the
                record, so <strong>every thumbnail Grit generated for a resource file field was
                orphaned</strong>: produced, paid for, and referenced by nothing. Doing the
                primary transform inline is what fixes it.
              </p>
              <p>
                Presigned uploads go straight from the browser to S3 and never pass through
                your server, so they cannot be transformed this way. Use the multipart endpoint
                for fields that want optimisation.
              </p>

              <h2 id="bombs">Decompression bombs are refused</h2>
              <p>
                A solid-colour PNG compresses to almost nothing whatever its dimensions, so
                an upload that passes every file-size check on the way in can still be
                enormous once decoded. Measured: a <strong>165 KB file at 12000x12000
                allocated 224 MB</strong>, and ten concurrent uploads of it would have been
                2.2 GB.
              </p>
              <p>
                The dimensions are read from the header before any pixels are allocated, and
                anything over <code className={C}>MaxPixels</code> is refused. The default is
                50 megapixels, which passes a 48 MP professional camera frame and refuses the
                bomb.
              </p>

              <h2 id="not">What is not optimised</h2>
              <ul>
                <li>
                  <strong>GIF</strong>, deliberately. Decoding one keeps the first frame only,
                  so optimising an animation would silently throw it away.
                </li>
                <li>
                  <strong>SVG</strong> is rejected at upload rather than optimised. It is a
                  stored-XSS vector and always was.
                </li>
                <li>
                  <strong>PDF and video.</strong> The shape generalises, the implementation
                  does not: both need external binaries, and neither fits in a static Go
                  binary.
                </li>
              </ul>
            </div>

            <div className="mt-16 flex items-center justify-between border-t border-border/50 pt-8">
              <Link href="/docs/batteries/storage">
                <Button variant="ghost" className="gap-2">
                  <ArrowLeft className="h-4 w-4" />
                  File Storage
                </Button>
              </Link>
              <Link href="/docs/backend/variants">
                <Button variant="ghost" className="gap-2">
                  Product Variants
                  <ArrowRight className="h-4 w-4" />
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
