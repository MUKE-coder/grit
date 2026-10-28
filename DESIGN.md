# DESIGN.md — the visual contract

> TERMINOLOGY.md settles the words. API_DESIGN.md settles the HTTP contract.
> This settles the pixels.
>
> It is written to be **checked**, not admired. Every rule below is either a
> number a script can count or a question with a yes/no answer. Run
> `/grit-design-review` on a page before you call it done.

---

## Where this comes from

A designer produced the system this document describes: tokens, a component
CSS bundle and two page designs. Those are the source, in this order:

1. **`docs/app/grit.css`** — the component layer, used as published. It is
   the designer's `bundle.css` byte for byte, so it can be replaced wholesale
   on an update rather than hand-translated. **Do not edit it to fix a page.**
2. **`docs/app/globals.css`** — the tokens, declared once.
3. **`docs/tailwind.config.ts`** — Tailwind reading those tokens. No
   arithmetic, no second source of truth.
4. **`docs-redesign/extracted/{landing,quickstart}`** — the mockups, for
   anything the three above do not pin down.

When two of them disagree, the one further up this list wins, with one
exception recorded in §4: the radius scale follows the component CSS rather
than the token file, because what the components draw is what the site looks
like.

## Why it exists, with the measurement

The marketing site was audited on 2026-10-28, before any of this landed:

| What | Then | Now | Budget |
|---|---|---|---|
| Distinct type styles (size + weight) | **40** | 30 | ≤ 12 |
| Border radii, against no declared scale | **8** | 0 off-scale | on-scale only |
| Distinct section paddings | **6** | 1 | ≤ 3 |
| `<section>` elements | **36** | 12 | ≤ 12 |
| Top-level headings (h1 + h2) | **30** | 12 | ≤ 12 |
| Contrast failures | **28** | 0 | 0 |

Nothing on that page was ugly. It was *unresolved*: forty type styles means
there is no scale, eight radii with no token behind them means every corner
was picked by hand, and thirty headings competing for attention means none of
them wins. That is what "we have not found our sweet spot" feels like from
the inside, and it was fixed by subtraction rather than by more design.

The same audit found `.glow-purple` in the stylesheet while the theme's
accent was sky blue. One line, two identities, nothing comparing them. That
is the failure this document exists to make impossible.

---

## 1. One identity

**The website looks like the product.** Somebody who clicks from the landing
page into a generated admin panel should not feel a seam.

### Colour

Declared once in `globals.css`, read by Tailwind, never hard-coded.

| Token | Dark | Light | For |
|---|---|---|---|
| `--bg` | `#08090B` | `#FFFFFF` | Page |
| `--surface-1` | `#0F1114` | `#F7F7F8` | Cards, panels |
| `--surface-2` | `#15181C` | `#F0F0F2` | Hover surface, wells |
| `--surface-3` | `#1C2025` | `#E8E8EB` | A surface on a surface |
| `--border` | `#23272D` | `#E4E4E7` | Every hairline |
| `--border-strong` | `#343A42` | `#C8CAD0` | A control's edge, hover, focus track |
| `--fg` | `#EDEEF0` | `#0B0C0E` | Primary text |
| `--fg-muted` | `#A1A7B0` | `#555B64` | Secondary text |
| `--fg-subtle` | `#7D848E` | `#6B717A` | Metadata, eyebrows. AA floor |
| `--accent` | `#A99BFF` | `#5B47E0` | The brand. Rationed, see below |
| `--accent-fg` | `#08090B` | `#FFFFFF` | The label **on** an accent fill |
| `--accent-soft` | `#A99BFF1A` | `rgba(91,71,224,.07)` | A wash, never a fill |
| `--accent-line` | `#A99BFF52` | `rgba(91,71,224,.30)` | The border of a washed panel |
| `--success` | `#3DD68C` | `#0B7A45` | |
| `--warning` | `#F5A524` | `#A35F00` | |
| `--danger` | `#FF6B6B` | `#C93636` | |

Code colours (`--code-bg`, `--code-fg`, `--code-keyword`, `--code-string`,
`--code-fn`, `--code-number`, `--code-comment`) **stay dark in both themes**,
so a terminal looks identical everywhere and a screenshot of one is stable.

**The accent is rationed, and what counts is the fill.** At most **three
accent-filled** elements per viewport. A fill is a call to action, and two of
those compete for the same click.

A *wash* is not a fill. `--accent-soft` is a 10% tint behind a panel, and ten
of those on a page read as one family rather than ten calls to action.
Accent-tinted **text** and icons are the same: a grid of accent-coloured
plugin icons is one set. The rule exists to stop two primary buttons sitting
next to each other.

A page where every card has a violet border has no accent; it has a
background.

**Contrast is not negotiable.** Body text ≥ 4.5:1, large text (≥ 24px) and UI
borders ≥ 3:1. `--fg-subtle` is at the floor deliberately; anything lighter
fails. The review measures this and does not take a judgement call.

**shadcn's `--accent` is a trap.** In shadcn it means *hover surface*; here it
means the brand. The alias is `--accent-surface`, pointed at `--surface-2`.
Anything written as `hover:bg-accent/50` is a hover state that turned violet,
and 74 of them had.

---

## 2. Type

**Two faces.** Geist for everything, Geist Mono for code and for metadata
(eyebrows, stats, labels, counts). `font-display` is kept as a name because
the docs use it in dozens of places, but it resolves to Geist.

### The twelve, by name

They are **named in `tailwind.config.ts`**, not described here and left to
be reproduced by hand. Use the name. A design that needs a thirteenth style
is wrong, not special.

| Class | Size / line | Weight | Use |
|---|---|---|---|
| `text-display-xl` | 80 / 84 | 600 | The hero h1, at the largest breakpoint |
| `text-display` | 56 / 60 | 600 | A section opener |
| `text-stat` | 64 / 64 | 500 | A single number that carries a point |
| `text-h1` | 40 / 48 | 600 | A docs page title |
| `text-h2` | 28 / 36 | 600 | Sub-sections |
| `text-h3` | 20 / 28 | 600 | Card titles |
| `text-h4` | 16 / 24 | 600 | The smallest heading |
| `text-lead` | 20 / 32 | 400 | The one paragraph under a heading |
| `text-body` | 16 / 28 | 400 | Everything else |
| `text-body-sm` | 14 / 22 | 400 | Dense UI, table cells |
| `text-caption` | 13 / 20 | 400 | Captions, helper text |
| `text-eyebrow` | 12 / 16 | 500 | Mono, uppercase, `0.08em` |
| `text-code` | 13.5 / 24 | 400 | Mono |

**Never invent a type class.** Tailwind emits nothing for a class it does not
know, and nothing is not an error: `text-lead` was written in two components
and existed in no stylesheet, so every section lead on the landing page
rendered at the inherited 16px instead of 20, silently, until somebody
measured it. If a name is not in the table above, it does not exist.

**Hierarchy comes from the ratio, not from the count.** Adding a size between
two sizes always makes a page read flatter, never sharper.

### The measure

**`max-width: 72ch` holds 104 characters in Geist. Never use `ch`.**

The `ch` unit is the advance of the digit `0`, and Geist's digits are tabular
and wide. Measured in the browser at 14, 15, 16 and 19px, `ch` is `0.663em`
where a real prose character averages `0.461em`: a 44% overstatement, the
same at every size.

There is one token, **`--measure: 33em`**, which is 72 characters times the
measured width of one, in a unit that tracks the font size. Prose uses
`max-w-[var(--measure)]`. The design system's own token says `72ch` and has
this bug; ours does not.

---

## 3. Space

**One rhythm, repeated.** The scale is
**4 / 8 / 12 / 16 / 24 / 32 / 48 / 64 / 96 / 128**. Nothing else. No
`mt-[18px]`.

Three section paddings:

| Value | For |
|---|---|
| 64px | Two sections that belong together |
| 96px | Every normal section |
| 128px | Hero, final CTA |

Containers: **1200px** for marketing, **1440px** for the docs shell, and the
measure above for prose. Gutters are 24px at every breakpoint, so nothing
touches the edge on a phone.

---

## 4. Surface

This is where "crisp" actually lives, and it is mostly restraint.

**Structure comes from lines, not elevation.** A shadow is for something that
genuinely floats: a dropdown, a dialog, a drag preview.

- **Border**: `1px solid var(--border)`. Always 1px. Never 2px for emphasis;
  use `--border-strong`.
- **Radius**: the eight below, and no others.

  | Value | Token | Tailwind | For |
  |---|---|---|---|
  | 4px | `--radius-xs` | `rounded-xs` | Inline code, ticks, the smallest chips |
  | 6px | `--radius-sm` | `rounded`, `rounded-sm` | Badges, inputs, small controls |
  | 8px | `--radius-btn-sm` | `rounded-btn-sm` | A 32px button, an icon button |
  | 10px | `--radius-md` | `rounded-md`, `rounded-lg` | Cards, panels, code blocks |
  | 12px | `--radius-btn-lg` | `rounded-btn` | A 48px button, a command row |
  | 16px | `--radius-lg` | `rounded-xl` | Hero panels, feature slabs |
  | 20px | `--radius-xl` | `rounded-2xl`, `rounded-3xl` | Bento cards, showcase panels |
  | full | `--radius-full` | `rounded-full` | Pills, avatars, icon buttons |

  Declared once as `--radius-*` and read by Tailwind. The old mapping derived
  6, 8 and 10 from a single `--radius`, and the 8px it produced was on no
  scale and on every button.

  Eight values, where this document used to say four, because the published
  component CSS draws with eight and that file is the source. The budget that
  mattered was never the count: it was whether a corner is a value somebody
  chose. So the check is **every radius is on the scale**, and a page may
  legitimately use every step.

  **A nested radius is derived, never chosen.** A child inside a padded
  parent takes the parent's radius minus the inset, so the curves stay
  concentric: `.g-tab` is 7px because it sits in a 10px `.g-tabs` with a 3px
  pad. 7 is not a ninth value and nothing else may use it.

  **A drawing is not an interface.** `data-decorative` opts a subtree out of
  the radius check, for geometry that represents a real object rather than a
  control: the phone frame's 28px bezel and 22px screen are the handset's own
  corners. Use it there and nowhere else. Without it the audit reports the
  same known-good value on every page showing a device, which is how a check
  stops being read.

- **Shadow**: two.
  - `--shadow-pop` for something that floats.
  - `--shadow-focus` for the focus ring: `0 0 0 2px var(--bg), 0 0 0 4px var(--accent)`.
- **Wash**: one radial `--accent-soft` glow per section at most, and only on
  the hero, the batteries, the benchmarks, the deploy section and the final
  CTA. Two glows on a screen cancel each other out.

### The signature

Grit has one visual motif, and consistency is what makes it a signature
rather than decoration.

**The rails and the corner mark.** A full-bleed 1px rule above each section, a
1200px frame with 1px rails down both sides running the whole page, and a
**7×7px accent square** where a rail meets a rule. In CSS: `g-frame`,
`g-cnr g-cnr--tl`, `g-cnr g-cnr--tr`.

The only thing that makes the squares read as a signature rather than as
decoration is that **every section has them**. A section without them looks
like a rendering bug.

**The 56px line grid** (`g-gridbg`) sits behind the hero, the stats band and
the closing CTA only. It is a texture: if you can see it clearly it is too
strong.

**One phrase per heading in the violet gradient** (`g-hl`). Never a whole
heading and never body copy, because emphasis applied evenly is no emphasis.

---

## 5. Motion

GSAP for scroll-driven work, CSS transitions for state.

| Token | Value | For |
|---|---|---|
| `--duration-fast` | 150ms | Hover, focus, a state flip |
| `--duration-base` | 240ms | A panel, a sliding indicator |
| `--duration-slow` | 480ms | An entrance |
| `--ease-out` | `cubic-bezier(.2,.8,.2,1)` | Anything arriving |
| `--ease-in-out` | `cubic-bezier(.65,0,.35,1)` | Anything moving between two states |

- **Nothing animates more than once.** A section that re-animates every time
  it scrolls back into view is a page that will not sit still to be read.
  `<Reveal>` disconnects its observer after it fires.
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

That block is not optional and the review checks for it. A GSAP timeline also
has to check `gsap.matchMedia()`, because the CSS override does not reach a
JS tween. And an element that starts at `opacity: 0` waiting for an observer
is invisible content, not a disabled animation: under reduced motion it must
render in place with no observer at all.

---

## 6. Components

Use `docs/components/grit/*`. They wrap `grit.css` and carry the behaviour
that a class cannot. `/docs/design` renders every one of them beside its
token values, and is the page to check a change against.

### Button

**The primary is inverted.** A near-white fill with a dark label, not the
brand violet. That is the whole reason violet can mean "brand" everywhere
else instead of meaning "there is a button here".

| Variant | Fill | Label | Use |
|---|---|---|---|
| `primary` | `--fg` | `--bg` | The one action on the view |
| `secondary` | transparent, `--border-strong` edge | `--fg` | Everything else |
| `ghost` | transparent | `--fg-muted` | Tertiary, toolbars |
| `accent` | `--accent` | `--accent-fg` | Rare: one highlighted action inside a dark panel |

Heights 40 / 32 (`sm`) / 48 (`lg`), radius 10 / 8 / 12 to match. Icon 16px,
8px gap. Every icon-only button carries an `aria-label`.

A primary button **does not change hue on hover**; it darkens. 27 of them
flipped from near-white to violet under the cursor, which reads as a
different button rather than the same one.

### Card

`g-card`: 1px `--border`, radius 16, `--surface-1`, 24px padding. On hover
the border goes to `--border-strong` and **nothing moves**. No lift, no
scale, no shadow. The border is the whole signal, which is why it has to be
1px everywhere else.

`g-bento`: the same idea at radius 20, for a bento cell.

### Eyebrow

Mono, uppercase, 12px, `0.08em`, `--fg-subtle`, **wrapped in braces**:
`{ BENCHMARKS }`. It is the only uppercase on the site, which is what makes
it read as a label rather than as shouting. When one section has one, every
sibling section has one.

### Terminal and code

`g-terminal`: dark in both themes. Three dots, a title, a copy button that
**confirms in place for 1.6s** and never raises a toast. Inline code is
`g-code-inline` at radius 4.

---

## 7. Page shapes

### A marketing page

A reader gives a landing page about eight seconds and three scrolls.

1. **Hero** — one claim, one paragraph, one primary action, one secondary.
   The install command is visible without scrolling.
2. **Proof** — the thing working. A screenshot or a terminal, not a promise.
3. **Three to five sections**, each making exactly one point.
4. **Social proof** — stars, users, what people built.
5. **Final CTA** — the same action as the hero, said once more.

Every section uses `<Section>`, so it gets the rule, the rails and the corner
marks without being asked.

**Twelve sections is the ceiling.** Sections past about twelve are not read;
they are scrolled past on the way to the footer, and their cost is that they
bury the ones that would have been read.

When cutting, the test is not "is this true" but "does a reader who stops
here know less than one who stops at the section above".

### A documentation page

`<DocsShell>`: three columns inside 1440 — 264 sidebar, 760 content, 232
on-this-page. The TOC hides below 1180 and the sidebar becomes a drawer below
1024.

1. `<DocsPageHeader>`: breadcrumb, h1, one lead paragraph, then the meta row
   (reading time, the version it is current for, Copy page, Edit on GitHub).
2. The shortest thing that works, in a code block, above the fold.
3. Prose at `--measure`. Never the full column.
4. The page's spine in the right rail: every h2, every h3 inside prose, and
   anything marked `data-toc`. If a page has more than eight, it is two
   pages.
5. Ends with two or three next links, each saying what is through it.

---

## 8. Accessibility, as a floor

Not a section at the end. These are checked:

- Contrast as in §1.
- Every interactive element reachable by Tab, in visual order.
- `:focus-visible` on everything focusable, using `--shadow-focus`. Never
  `outline: none` without a replacement.
- Targets ≥ 24×24px, 44×44 for anything touched on a phone. A text link in a
  block of text is exempt (WCAG 2.5.8) and the audit knows the difference: no
  vertical padding, no fill, no border, height within its line box.
- Icon-only controls have an accessible name.
- One `h1`. **Heading levels do not skip** — a card title that is an `h3`
  above the page's first `h2` is a missing section to anyone navigating by
  heading.
- A `role="tablist"` promises arrow keys. A tablist that only answers clicks
  has told a screen reader to expect something it does not do.
- Images: real `alt`, or `alt=""` for decoration.

---

## 9. The budgets

What `/grit-design-review` counts, and what `docs/public/design-audit.js`
measures. A page over budget is not shipped.

| Budget | Limit |
|---|---|
| Distinct type styles (size + weight) | 12 |
| Border radii off the scale | 0 |
| Distinct section paddings | 3 |
| Distinct shadows | 2 |
| Sections per marketing page | 12 |
| h1 + h2 per page | 12 |
| Accent-**filled** elements per viewport | 3 |
| Paragraphs over the measure | 0 |
| Contrast failures | 0 |
| Heading level skips | 0 |
| Targets under 24×24 that are not text links | 0 |
| Hard-coded hex outside the token file | 0 |

Three of these were measuring the wrong thing until 2026-10-28, and each
failure mode is worth knowing because the next check will have the same
shape:

- **The accent fill** compared red, green and blue and ignored alpha, so
  `bg-accent/10` counted as a call to action. A wash is not a fill.
- **The target size** exempted only `<a>` inside `<p>`, so breadcrumbs and
  meta rows were flagged, which teaches people to ignore the check.
- **The line length** assumed `0.5em` per character. It now renders a
  sentence in the paragraph's own font and divides, and it reports characters
  rather than `ch`, because `ch` is the unit that caused the problem.

A check that cries wolf is worse than no check, and a check that passes a
real fault is worse again.

---

## 10. Never

- A second accent colour to make something stand out. Use size or space.
- A violet primary button. The primary is inverted; that is the system.
- A shadow to separate two things that are next to each other. Use a line.
- A radius off the scale, or a nested radius that is not the parent's minus
  the inset.
- `ch` as a measure. Use `var(--measure)`.
- A type class that is not one of the twelve. Tailwind emits nothing for a
  class it does not know, and nothing is not an error.
- `hsl(var(…))`. The tokens are hex; mixing the two silently produces black.
- Editing `grit.css` to fix a page. It is the designer's file, replaced
  wholesale on an update.
- Centre-aligned body text. Headings may centre; paragraphs do not.
- An animation that repeats on every scroll into view.
- Text over an image without a scrim that keeps it at 4.5:1.
- A new font. There are two.
- `!important`, except in the reduced-motion block.
