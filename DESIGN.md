# DESIGN.md — the visual contract

> TERMINOLOGY.md settles the words. API_DESIGN.md settles the HTTP contract.
> This settles the pixels.
>
> It is written to be **checked**, not admired. Every rule below is either a
> number a script can count or a question with a yes/no answer. Run
> `/grit-design-review` on a page before you call it done.

---

## Why this exists, with the measurement

The marketing site was audited on 2026-10-28. One page, the landing page:

| What | Found | Should be |
|---|---|---|
| Distinct type styles (size + weight) | **40** | ≤ 12 |
| Distinct border radii | **8** | 4 |
| Distinct section paddings | **6** | 3 |
| `<section>` elements | **36** | ≤ 12 |
| Top-level headings (h1 + h2) | **30** | ≤ 12 |

Nothing on that page is ugly. It is *unresolved*: forty type styles means
there is no scale, eight radii means there is no radius token, and thirty
headings competing for attention means none of them wins. That is what "we
have not found our sweet spot" feels like from the inside, and it is fixable
by subtraction rather than by more design.

The same audit found `.glow-purple` still in the stylesheet while the theme's
accent is sky blue. That single line is the whole problem in miniature: two
identities in one file, with nothing comparing them.

---

## 1. One identity

**The website looks like the product.** Somebody who clicks from the landing
page into a generated admin panel should not feel a seam. Today there is one:
the site is sky blue (`#0ea5e9`) and the product's design system is purple
(`#6c5ce7`).

**Decision: the site adopts the product's palette.** Two reasons. The site's
job is to make people feel what the product feels like, and sky blue is the
single most common accent in developer tooling, so it says nothing about us.

### Colour

Defined once, as CSS custom properties. Nothing hard-codes a hex.

| Token | Dark | Light | For |
|---|---|---|---|
| `--bg` | `#0a0a0f` | `#fcfcfd` | Page |
| `--bg-raised` | `#111118` | `#ffffff` | Cards, panels |
| `--bg-sunken` | `#07070b` | `#f4f4f6` | Code blocks, wells |
| `--line` | `#22222e` | `#e6e6ea` | Every border |
| `--line-strong` | `#2f2f3e` | `#d4d4da` | Hover, focus ring track |
| `--ink` | `#e8e8f0` | `#13131a` | Primary text |
| `--ink-muted` | `#9090a8` | `#5b5b6b` | Secondary text |
| `--ink-faint` | `#7c7c96` | `#76768a` | Metadata. AA floor, do not go lighter |
| `--accent` | `#6c5ce7` | `#5b4bd6` | One per screen, see below |
| `--accent-ink` | `#ffffff` | `#ffffff` | Label on accent |
| `--ok` | `#00b894` | `#00866b` | |
| `--warn` | `#fdcb6e` | `#8a6200` | |
| `--bad` | `#ff6b6b` | `#c42b2b` | |

**The accent is rationed.** At most **three** accent-coloured things per
viewport: the primary action, the active nav item, and one emphasis. A page
where every card has a purple border has no accent, it has a background.

**Contrast is not negotiable.** Body text ≥ 4.5:1, large text (≥ 24px) and UI
borders ≥ 3:1. `--ink-faint` is at the floor on purpose; anything lighter
fails. The review skill measures this and will not take a judgement call.

---

## 2. Type

Three faces, already loaded, no more:

- **Geist** — everything by default.
- **Ubuntu** (`font-display`) — headings at `text-3xl` and above only.
- **JetBrains Mono** (`font-mono`) — code, and metadata (eyebrows, stats,
  labels, counts).

### The scale

Twelve styles. If a design needs a thirteenth, it is wrong, not special.

| Name | Size / line | Weight | Use |
|---|---|---|---|
| `display` | `clamp(2.75rem, 6vw, 5rem)` / 0.95 | 700 | One per page. The h1. |
| `h1` | `clamp(2rem, 4vw, 3rem)` / 1.05 | 700 | Section openers |
| `h2` | `1.5rem` / 1.2 | 600 | Sub-sections |
| `h3` | `1.125rem` / 1.3 | 600 | Card titles |
| `lead` | `1.125rem` / 1.6 | 400 | The one paragraph under a heading |
| `body` | `0.9375rem` / 1.65 | 400 | Everything else |
| `body-sm` | `0.875rem` / 1.55 | 400 | Dense UI, table cells |
| `meta` | `0.8125rem` / 1.4 | 400 | Captions, helper text |
| `eyebrow` | `0.6875rem` / 1 | 500 | Mono, uppercase, `0.12em` tracking |
| `code` | `0.875rem` / 1.6 | 400 | Mono |
| `num` | `clamp(1.75rem, 3vw, 2.5rem)` / 1 | 700 | Stats. Tabular figures. |
| `button` | `0.875rem` / 1 | 500 | |

**Hierarchy comes from the ratio, not from the count.** The benchmark sites
get their presence from a display size five times the body size with nothing
in between competing. Adding a size between two sizes always makes a page
read flatter, never sharper.

Headings over three words set `tracking-tight` (`-0.02em`). Body never does.

---

## 3. Space

**One rhythm, repeated.** Three section paddings, not six:

| Token | Value | For |
|---|---|---|
| `section` | `py-24` (96px) | Every normal section |
| `section-lg` | `py-32` (128px) | Hero, final CTA |
| `section-sm` | `py-16` (64px) | Two sections that belong together |

Inside a section the vertical scale is **4 / 8 / 12 / 16 / 24 / 32 / 48 / 64**.
Nothing else. No `mt-[18px]`.

Container is `max-w-6xl` (1152px) for marketing, `max-w-3xl` (768px) for
reading. A line of body text never exceeds **72 characters**; that is what the
reading container is for.

**Gutters are 24px at every breakpoint**, so nothing touches the edge on a
phone.

---

## 4. Surface

This is where "crisp" actually lives, and it is mostly about restraint.

**Structure comes from lines, not elevation.** Every benchmark in
`best-ui-benchmark/` builds its structure from 1px borders. Shadows are for
things that genuinely float: a dropdown, a dialog, a drag preview.

- **Border**: `1px solid var(--line)`. Always 1px. Never 2px for emphasis;
  use `--line-strong` instead.
- **Radius**: four values and no others.
  | Token | Value | For |
  |---|---|---|
  | `rounded-sm` | 6px | Badges, inputs, small controls |
  | `rounded-md` | 10px | Cards, panels, code blocks |
  | `rounded-lg` | 16px | Hero panels, feature slabs |
  | `rounded-full` | 9999px | Pills, avatars, icon buttons |
- **Shadow**: exactly two, and only on floating things.
  - `shadow-pop`: `0 8px 24px -8px rgb(0 0 0 / 0.4)`
  - `shadow-lift`: `0 2px 8px -2px rgb(0 0 0 / 0.25)` for a dragged item.
- **Glow**: one, `glow-accent`, reserved for the hero and the final CTA. Two
  glows on a screen cancel each other out.

### The signature

Grit has one visual motif, and consistency is what makes it a signature rather
than decoration:

**Corner ticks.** A 10px L-shaped mark in `--line-strong` at the corners of a
framed region, as on the Filament panels in the benchmark set. Used on: the
hero panel, featured cards, and the final CTA. Nothing else.

**Dot grid.** `--line` at 8% opacity, 32px pitch, behind the hero and the
final CTA only. It is a texture, not a background; if you can see it clearly
it is too strong.

---

## 5. Motion

GSAP for scroll-driven work, CSS transitions for state.

- **State** (hover, focus, open): `150ms`, `cubic-bezier(0.4, 0, 0.2, 1)`.
- **Entrance** (on scroll): `500ms`, `cubic-bezier(0.16, 1, 0.3, 1)`, with a
  `12px` rise and a fade. Stagger `60ms`.
- **Nothing animates more than once.** A section that re-animates every time
  it scrolls back into view is a page that will not sit still to be read.
- **Nothing moves more than 24px.** Large travel reads as a loading failure.

```css
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
    scroll-behavior: auto !important;
  }
}
```

That block is not optional and the review skill checks for it. A GSAP timeline
also has to check `gsap.matchMedia()` with `(prefers-reduced-motion: reduce)`,
because the CSS override does not reach a JS tween.

---

## 6. Components

### Button

Three variants, and the hierarchy is the point: **one primary per view.**

| Variant | Dark | Use |
|---|---|---|
| `primary` | `--accent` fill, `--accent-ink` label | The one action |
| `secondary` | `--bg-raised` fill, `--line` border | Everything else |
| `ghost` | transparent, `--ink-muted` label | Tertiary, toolbars |

Height 40px (`sm`: 32px, `lg`: 48px). Radius `full` for marketing, `sm` for
app UI; pick one per surface and never mix them in a row. Icon 16px, 8px gap.
Every icon-only button carries an `aria-label`.

### Card

```
border: 1px solid var(--line)
radius: 10px
background: var(--bg-raised)
padding: 20px
hover: border-color → var(--accent) / 40%, 150ms. Nothing moves.
```

A card does **not** lift, scale or shadow on hover. The border is the whole
signal, which is why it has to be 1px everywhere else.

### Eyebrow

Mono, uppercase, 11px, `0.12em` tracking, `--ink-faint`. Sits 8px above a
heading. Optional, but when present every sibling section has one.

### Code block

`--bg-sunken`, radius 10px, `1px solid var(--line)`, 16px padding, mono 14px.
Language label in the top-right as an eyebrow. A copy button is an icon button
that confirms in place, never a toast.

---

## 7. Page shapes

### A marketing page

A reader gives a landing page about eight seconds and three scrolls. So:

1. **Hero** — one claim, one paragraph, one primary action, one secondary.
   The install command is visible without scrolling.
2. **Proof** — the thing working. A screenshot or a terminal, not a promise.
3. **Three to five sections**, each making exactly one point.
4. **Social proof** — stars, users, what people built.
5. **Final CTA** — the same action as the hero, said once more.

**Twelve sections is the ceiling.** The landing page has 36. Sections past
about twelve are not read; they are scrolled past on the way to the footer,
and their cost is that they bury the ones that would have been read.

When cutting, the test is not "is this true" but "does a reader who stops here
know less than one who stops at the section above".

### A documentation page

1. Eyebrow (the section), h1, one lead paragraph.
2. The shortest thing that works, in a code block, above the fold.
3. Prose in `max-w-3xl`. Never wider.
4. Every h2 in the right rail. If a page has more than eight h2s it is two
   pages.
5. Ends with two or three next links, each saying what is through it.

---

## 8. Accessibility, as a floor

Not a section at the end. These are checked:

- Contrast as in §1.
- Every interactive element reachable by Tab, in visual order.
- `:focus-visible` on everything focusable: 2px `--accent` ring, 2px offset.
  Never `outline: none` without a replacement.
- Targets ≥ 24×24px, 44×44 for anything touched on a phone.
- Icon-only controls have an accessible name.
- One `h1`. Heading levels do not skip.
- Images: real `alt`, or `alt=""` and `aria-hidden` if decorative.

---

## 9. The budgets

What `/grit-design-review` counts. A page over budget is not shipped.

| Budget | Limit |
|---|---|
| Distinct type styles (size + weight) | 12 |
| Distinct border radii | 4 |
| Distinct section paddings | 3 |
| Distinct shadows | 2 |
| Sections per marketing page | 12 |
| h1 + h2 per page | 12 |
| Accent-coloured elements per viewport | 3 |
| Body line length | 72ch |
| Contrast failures | 0 |
| Hard-coded hex outside the token file | 0 |

---

## 10. Never

- A second accent colour to make something stand out. Use size or space.
- A shadow to separate two things that are next to each other. Use a line.
- A radius that is not one of the four.
- Centre-aligned body text. Headings may centre; paragraphs do not.
- An animation that repeats on every scroll into view.
- Text over an image without a scrim that keeps it at 4.5:1.
- A new font. There are three.
- `!important`, except in the reduced-motion block.
