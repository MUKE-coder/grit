/**
 * Every Grit video, in one place.
 *
 * The home page and /docs/videos both show the walkthrough, and before this
 * they would each have carried their own copy of the id, the title and the
 * poster path. Two copies of a title is one title that goes stale.
 *
 * ── Only videos that exist ────────────────────────────────────────────────
 *
 * Nothing goes in this list until it is published and watchable. A planned
 * video belongs in the issue that plans it, not on a page that offers it: a
 * "coming soon" card is a page telling somebody to come back, which is worse
 * than a short list.
 *
 * ── Posters are local on purpose ──────────────────────────────────────────
 *
 * `poster` is a path under docs/public, never img.youtube.com. VideoFacade
 * exists so that nothing is requested from YouTube until somebody presses
 * play, and a remote thumbnail would make that request on page load, which
 * would defeat the whole arrangement while looking like it worked.
 *
 * To add a video: publish it, save a 1280x720 still into docs/public/video/,
 * and add an entry here.
 *
 * ── A note for whoever makes this list long ───────────────────────────────
 *
 * Both documentation indexers (`scripts/docs-index.py` for `grit docs`, and
 * `docs/scripts/search-index.py` for the site's search box) read the prose
 * out of page.tsx files. They do not read this file. With a handful of
 * videos that is fine, because /docs/videos names in its own prose what the
 * videos cover. Past that, give both indexers a reader for this config, the
 * way they have one for the system design pages, or the library becomes
 * unsearchable without anything saying so.
 */

export type Video = {
  /** The YouTube id, the part after `v=`. */
  videoId: string
  /** The real title, as published. */
  title: string
  /** One line on what it covers, in the words somebody would search for. */
  blurb: string
  /** A still served from this origin. See the note above. */
  poster: string
  /** Runtime, as a person would say it. */
  length: string
  /** The Grit version it was recorded against. */
  recordedOn: string
  /** Who it is for, so a senior can skip a setup video. */
  level: 'Beginner' | 'Intermediate' | 'Advanced'
}

export const WALKTHROUGH: Video = {
  videoId: 'YGeuuSFVezk',
  title: 'Grit: The Full-Stack Framework for Go & React (Full Walkthrough)',
  blurb:
    'The whole framework end to end: scaffolding a project, generating a resource, the admin panel it produces, and what you get without writing any of it.',
  poster: '/video/walkthrough.jpg',
  length: 'Full walkthrough',
  recordedOn: 'v3.575.0',
  level: 'Beginner',
}

/**
 * The library, newest first.
 *
 * The walkthrough is the first entry and the home page's video is the same
 * object, so the two cannot drift.
 */
export const VIDEOS: Video[] = [WALKTHROUGH]
