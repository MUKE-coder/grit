# The GopherCon Africa deck, as the website serves it

74 slides on a fixed 1920x1080 canvas, every style inline. `/talk` reads them
from here at build time, so this directory is what the site shows: nothing is
fetched at runtime.

## What is in here

| Path | What it is |
|------|------------|
| `deck.json` | `order` (the slide sequence), `sections` (the outline), `title` |
| `slides/<id>.html` | one `<section id="<id>">`, with the speaker notes in a trailing `<aside>` |

A slide not listed in `order` is not shown. The file name is the slide id and
the `<section>` carries the same id, because the viewer links to it by that id
(`/talk#d7`).

## These files came from a Slides artifact

The deck is authored as a Claude Slides artifact and copied here. That makes
this a second copy, and copies drift, so:

- **Edit the artifact, then re-copy.** Editing a file here changes the website
  and not the deck that gets presented, which is the worse of the two to have
  wrong.
- The copy is a plain file copy of `project/deck.json` and `project/slides/*`.
  Nothing is rewritten on the way in.

## The transform the site applies

`lib/talk-deck.ts` reads these files and adjusts three things that only matter
in a browser, which is why they are not baked in here:

1. `position:relative` on the section, so the pinned children (the corner mark,
   the cover's accent bar) have something to pin to.
2. `'Geist'` and `'Geist Mono'` become the CSS variables the site loads them
   under, because `next/font` gives a font a scoped family name rather than its
   real one.
3. `<x-shape kind="ellipse">` becomes a rounded `<div>`. That element is part of
   the artifact runtime and does not exist here.

If a slide starts using another `x-` element, add it to that transform or it
renders as nothing.
