import fs from 'node:fs'
import path from 'node:path'

// The GopherCon Africa deck, read from content/talk/<slug> at build time so
// /talk is fully static. See that directory's README for where the files come
// from and why this transform exists.

const TALK_DIR = path.join(process.cwd(), 'content', 'talk')

export interface TalkSlide {
  id: string
  /** 1-based position in the deck, which is what the counter shows. */
  number: number
  /** The <section>, ready to inject: notes removed, styles adjusted. */
  html: string
  /** Speaker notes, one entry per paragraph. */
  notes: string[]
  /** Set on the slide that opens a run, from data-section. */
  sectionLabel?: string
  /** True when the slide is on black, so the chrome can invert around it. */
  dark: boolean
}

export interface TalkSection {
  /** The id of the slide the run starts on. */
  start: string
  /** 1-based number of that slide. */
  number: number
  description: string
}

export interface TalkDeck {
  slug: string
  title: string
  slides: TalkSlide[]
  sections: TalkSection[]
}

interface DeckIndex {
  title?: string
  order: string[]
  sections?: Record<string, { description?: string; start?: string }>
}

/**
 * Three things the artifact runtime does for a slide that a browser does not.
 *
 * The font swap is the one that is easy to miss: next/font serves Geist under a
 * generated family name, so the literal 'Geist' in a slide's inline style
 * matches nothing and silently falls through to Arial. Every slide in the deck
 * is set in it, so that failure is the whole deck in the wrong typeface.
 */
function toBrowserHtml(raw: string): string {
  let html = raw

  // 1. The canvas. A <section> in a plain div is as tall as its content, so a
  //    slide that distributes its children down the page (the cover, anything
  //    using space-between) collapses to the top and the layout is silently not
  //    the one that was designed. Pinning also needs position on the section:
  //    the corner mark and the cover's accent bar are absolute against it.
  html = html.replace(
    /^(<section\b[^>]*?\bstyle=")/,
    '$1position:relative;overflow:hidden;box-sizing:border-box;width:1920px;height:1080px;',
  )

  // 2. The scoped font families. 'Geist Mono' first: 'Geist' is a prefix of it.
  html = html.split("'Geist Mono'").join('var(--font-geist-mono)')
  html = html.split("'Geist'").join('var(--font-geist)')

  // 3. <x-shape> belongs to the artifact runtime. Only the three dots on a
  //    terminal window use it, but map the shapes that have a CSS equivalent
  //    rather than only the one in use today.
  const RADIUS: Record<string, string> = {
    ellipse: '9999px',
    rounded: '8px',
    rect: '0px',
  }
  html = html.replace(
    /<x-shape\b([^>]*?)>\s*<\/x-shape>/g,
    (whole, attrs: string) => {
      const kind = /kind="([^"]*)"/.exec(attrs)?.[1] ?? 'rect'
      const radius = RADIUS[kind]
      // A shape this transform has no CSS for would render as an empty box in
      // the wrong place, which is worse than rendering as nothing.
      if (!radius) return ''
      const style = /style="([^"]*)"/.exec(attrs)?.[1] ?? ''
      return `<span style="${style};border-radius:${radius};display:block;flex:none"></span>`
    },
  )

  return html
}

function readSlide(dir: string, id: string, number: number): TalkSlide {
  const raw = fs.readFileSync(path.join(dir, 'slides', `${id}.html`), 'utf8')

  const aside = /<aside>([\s\S]*?)<\/aside>/.exec(raw)
  const notes = (aside?.[1] ?? '')
    .split(/\n\s*\n/)
    .map((p) => p.replace(/\s+/g, ' ').trim())
    .filter(Boolean)

  const body = raw.replace(/<aside>[\s\S]*?<\/aside>/, '').trim()
  const open = body.slice(0, body.indexOf('>') + 1)

  return {
    id,
    number,
    html: toBrowserHtml(body),
    notes,
    sectionLabel: /data-section="([^"]*)"/.exec(open)?.[1],
    // The deck has exactly two backgrounds. Read it rather than keep a list of
    // which slides are dark, which is a list that goes stale on the first edit.
    dark: /background:#000000/.test(open),
  }
}

export function getTalkDeck(slug: string): TalkDeck {
  const dir = path.join(TALK_DIR, slug)
  const index = JSON.parse(fs.readFileSync(path.join(dir, 'deck.json'), 'utf8')) as DeckIndex

  const slides = index.order.map((id, i) => readSlide(dir, id, i + 1))
  const numberOf = new Map(slides.map((s) => [s.id, s.number]))

  const sections: TalkSection[] = Object.values(index.sections ?? {})
    .map((s) => ({
      start: s.start ?? '',
      number: numberOf.get(s.start ?? '') ?? 0,
      description: s.description ?? '',
    }))
    // A section whose start slide was removed points at nothing, and a jump
    // target that does not exist is worse than a missing row in the outline.
    .filter((s) => s.number > 0)
    .sort((a, b) => a.number - b.number)

  return { slug, title: index.title ?? 'Talk', slides, sections }
}
