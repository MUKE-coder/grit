import Link from 'next/link'
import { ArrowLeft, ArrowRight, Film } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export const metadata = {
  title: 'Video plugin | Grit',
  description:
    'Video uploads converted by ffmpeg to an MP4 every browser and phone plays, with a poster frame, in a background queue that needs no Redis. Players for the web app and the Expo app.',
  alternates: { canonical: 'https://gritframework.dev/docs/plugins/video' },
}

export default function VideoPluginPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <Film className="h-5 w-5 text-primary" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">Plugin</span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">Video</h1>
          <p className="mb-8 text-lg leading-relaxed text-muted-foreground">
            Take a clip from a phone or a browser and turn it into something every client plays: one H.264 MP4 that
            starts before it has downloaded, a poster frame, and the size and length a feed needs to lay it out.
            ffmpeg does the work in the background, and the upload request never waits for it.
          </p>

          <CodeBlock language="bash" code={`grit plugin add video
grit migrate
pnpm install`} />

          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            The API&apos;s Docker image gets ffmpeg. To convert while developing, install it on your machine:{' '}
            <code>winget install ffmpeg</code>, <code>brew install ffmpeg</code> or <code>apt install ffmpeg</code>.
            Without it, uploaded videos wait and the log says why.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Upload, convert, play</h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            Upload the file as usual with <code>accepts=video</code>, which allows up to 300 MB, then hand its key to{' '}
            <code>POST /videos</code>. The answer is <code>202</code> with the video, pending.
          </p>
          <CodeBlock language="tsx" code={`import { createVideo } from "@/lib/video";
import { VideoPlayer } from "@/components/video-player";

const ref = await uploader.upload(file, file.name, { accepts: ["video"] });
const video = await createVideo(ref.key);

// Polls every two seconds while it converts, then plays it.
<VideoPlayer id={video.id} />

// In a feed: muted, looping, playing.
<VideoPlayer video={post.video} autoPlay />`} />
          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            The Expo app has the same pair, playing with <code>expo-video</code>, and{' '}
            <code>useUploadVideo()</code> takes a clip straight from <code>expo-image-picker</code>. The owner also
            hears <code>video.ready</code> or <code>video.failed</code> on their realtime channel, for an app that
            would rather listen than poll.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">What you get</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border text-left">
                  <th className="py-2 pr-4 font-medium text-foreground">Endpoint</th>
                  <th className="py-2 font-medium text-foreground">Does</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/50">
                  <td className="py-2 pr-4 font-mono text-xs">POST /api/v1/videos</td>
                  <td className="py-2">Queues the conversion of an upload you own, by its key.</td>
                </tr>
                <tr className="border-b border-border/50">
                  <td className="py-2 pr-4 font-mono text-xs">GET /api/v1/videos/:id</td>
                  <td className="py-2">
                    pending, processing, ready or failed, with the URL, poster, size and length once ready. Anyone
                    signed in sees a ready video; one still converting, only its owner.
                  </td>
                </tr>
                <tr>
                  <td className="py-2 pr-4 font-mono text-xs">DELETE /api/v1/videos/:id</td>
                  <td className="py-2">Removes your video and the files it made.</td>
                </tr>
              </tbody>
            </table>
          </div>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">The details that matter</h2>
          <ul className="mb-8 list-disc space-y-3 pl-6 leading-relaxed text-muted-foreground">
            <li>
              <strong className="text-foreground">One format, every client.</strong> H.264 and AAC in an MP4, capped
              at 720 on the short side, with the index at the front. Every browser, iOS and Android plays it, so
              there is no HLS to serve and no player to choose per platform. A portrait clip stays portrait.
            </li>
            <li>
              <strong className="text-foreground">The table is the queue.</strong> A worker claims a pending video
              with a conditional update, so a conversion survives a restart, several replicas share the work without
              Redis, and a claim held by a replica that died is taken over after fifteen minutes. One conversion at a
              time per replica, because ffmpeg uses every core it is given.
            </li>
            <li>
              <strong className="text-foreground">A crafted file cannot read the server.</strong> ffmpeg follows what a
              file says it is, and a playlist posing as a video can name other files for it to open. Only MP4, MOV and
              WebM are converted, each read with a forced demuxer and local files only, and a playlist is refused
              before ffmpeg sees it as input.
            </li>
            <li>
              <strong className="text-foreground">Only your own uploads.</strong> A video is made from an upload by its
              storage key, and a key that is not yours answers as not found, so it never confirms that someone
              else&apos;s file exists.
            </li>
            <li>
              <strong className="text-foreground">Failures say why.</strong> Too long, not a video, or a file that
              would not convert: the video is marked failed with a sentence for the person who uploaded it, and a
              failure that could be temporary is tried three times first. Limits are{' '}
              <code>videoService.Options.MaxEdge</code> and <code>MaxDuration</code>, set in <code>routes.go</code>:
              720 and three minutes.
            </li>
          </ul>

          <p className="mb-8 leading-relaxed text-muted-foreground">
            Built for the Instagram blueprint, and checked end to end on a new project: a 1080p clip uploaded over
            HTTP came back as a 720p H.264 MP4 with a poster within a second, served from storage, and another
            user sending the same key was told it does not exist.
          </p>

          <div className="mt-16 flex items-center justify-between border-t border-border/50 pt-8">
            <Link href="/docs/plugins/push">
              <Button variant="ghost" className="gap-2">
                <ArrowLeft className="h-4 w-4" />
                Push notifications
              </Button>
            </Link>
            <Link href="/docs/plugins/stripe">
              <Button variant="ghost" className="gap-2">
                Stripe payments
                <ArrowRight className="h-4 w-4" />
              </Button>
            </Link>
          </div>
        </div>
      </main>
    </div>
  )
}
