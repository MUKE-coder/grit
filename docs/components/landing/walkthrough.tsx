import { VideoFacade } from '@/components/video-facade'
import { WALKTHROUGH } from '@/config/videos'

/*
 * The walkthrough video on the landing page.
 *
 * The privacy behaviour and the markup live in VideoFacade, which the video
 * tutorials page uses for the same reason: a facade rather than an iframe, so
 * nothing is requested from YouTube until somebody presses play. This file is
 * now the landing page's choice of which video goes there, and nothing else.
 *
 * The video itself is declared in config/videos.ts, beside every other Grit
 * video, so the home page and /docs/videos cannot disagree about its title.
 */

export function Walkthrough() {
  return (
    <VideoFacade
      videoId={WALKTHROUGH.videoId}
      title={WALKTHROUGH.title}
      poster={WALKTHROUGH.poster}
      className="mx-auto max-w-4xl"
    />
  )
}
