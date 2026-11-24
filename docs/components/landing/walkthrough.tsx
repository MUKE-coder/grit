'use client'

import { useState } from 'react'
import Image from 'next/image'
import { Play } from 'lucide-react'

/*
 * The walkthrough video, as a facade rather than an iframe.
 *
 * A YouTube iframe on the landing page is roughly a megabyte of third-party
 * JavaScript and a set of cookies, paid for by every visitor including the
 * large majority who never press play. On a page whose whole argument is that
 * the framework is fast, that would be an odd thing to ship.
 *
 * So the section is a local thumbnail and a button. The iframe is created on
 * the first click and not before, which means:
 *
 *   - no third-party request at all until somebody asks for the video
 *   - the poster is served from this origin, so it is covered by the same
 *     caching and the same CSP as everything else on the page
 *   - youtube-nocookie.com once it does load, and autoplay, because a click
 *     on a play button is a request to play
 *
 * The title and the channel are the real ones, read from YouTube's oEmbed
 * endpoint rather than typed from memory.
 */

const VIDEO_ID = 'YGeuuSFVezk'
const TITLE = 'Grit: The Full-Stack Framework for Go & React (Full Walkthrough)'
const WATCH_URL = `https://www.youtube.com/watch?v=${VIDEO_ID}`

export function Walkthrough() {
  const [playing, setPlaying] = useState(false)

  return (
    <div className="mx-auto max-w-4xl">
      <div className="relative overflow-hidden rounded-2xl border border-border bg-surface shadow-[0_32px_80px_-24px_rgba(2,6,23,0.6)]">
        {/* 16:9, reserved before anything loads, so the page does not jump. */}
        <div className="relative aspect-video">
          {playing ? (
            <iframe
              className="absolute inset-0 h-full w-full"
              src={`https://www.youtube-nocookie.com/embed/${VIDEO_ID}?autoplay=1&rel=0&modestbranding=1`}
              title={TITLE}
              allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
              allowFullScreen
            />
          ) : (
            <button
              type="button"
              onClick={() => setPlaying(true)}
              className="group absolute inset-0 h-full w-full cursor-pointer"
              aria-label={`Play the walkthrough: ${TITLE}`}
            >
              <Image
                src="/video/walkthrough.jpg"
                alt=""
                fill
                sizes="(max-width: 896px) 100vw, 896px"
                className="object-cover"
                priority={false}
              />
              {/* Just enough for the play control to be legible. The
                  thumbnail is artwork with its own title in it, and dimming a
                  third of it to print that title again was the first thing
                  wrong with this section. */}
              <span
                aria-hidden="true"
                className="absolute inset-0 bg-black/15 transition-colors duration-200 group-hover:bg-black/5"
              />
              <span
                aria-hidden="true"
                className="absolute left-1/2 top-1/2 grid h-[72px] w-[72px] -translate-x-1/2 -translate-y-1/2 place-items-center rounded-full bg-white/95 shadow-lg transition-transform duration-200 ease-[cubic-bezier(0.17,0.55,0.55,1)] group-hover:scale-105 group-focus-visible:scale-105 motion-reduce:transition-none motion-reduce:group-hover:scale-100"
              >
                <Play className="ml-1 h-7 w-7 fill-black text-black" />
              </span>
            </button>
          )}
        </div>
      </div>

      <p className="mt-4 text-center text-[13.5px] text-fg-muted">
        Nothing is requested from YouTube until you press play.{' '}
        <a href={WATCH_URL} target="_blank" rel="noreferrer" className="g-link">
          Watch it on YouTube instead
        </a>
        .
      </p>
    </div>
  )
}
