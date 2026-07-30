# Grit UI — Roadmap

Grit UI is being rebuilt as a **UI Blocks** library in the shape of Tailwind Plus:
top-level categories, subcategories of blocks, and a viewer where you preview a
block, flip to its code, and copy it.

The 100 components shipped in v3.112.0 are **retired**. They were recovered from
git history and repaired enough to serve, but they were never designed to a
consistent spec and do not meet the bar for this site. They remain in git at tag
`v3.112.0` if any are ever worth reviving.

---

## Structure

Three levels:

```
/                          overview — every category and its subcategories
/[category]                e.g. /marketing — subcategory cards with counts
/[category]/[subcategory]  e.g. /marketing/hero-sections — every block, stacked
```

The subcategory page is the product. Each block on it gets:

- a **live preview** in a resizable frame
- **Preview / Code** toggle
- **breakpoint switcher** — mobile, tablet, desktop
- **light / dark** preview toggle, independent of the site theme
- **copy source**, and the `grit ui add` / `npx shadcn add` command
- a name in the Tailwind idiom — *"Simple centered"*, *"With background image"*

---

## Phase 0 — Foundation

Nothing else can start until this lands.

| Item | Detail |
|---|---|
| Design system | **Light by default**, dark via toggle. Neutral surfaces, hairline borders, monospace section labels, dot-grid edges. |
| Theme | `next-themes`, class strategy, no flash on load |
| Data model | `registry/<category>/<subcategory>/<block>.tsx` + a manifest with title, description, order |
| Routing | the three levels above, all statically generated |
| Viewer | preview frame, breakpoint switcher, code panel, copy, per-block theme |
| Registry | keep `/r/[name].json` and `grit ui add` working — the install path does not change |
| Retire | remove the 100 old blocks and their metadata |

---

## Phase 1 — Marketing

23 subcategories. **3 blocks each to start**, expanded later.

**Page Sections** — Hero Sections ✅ (12) · Feature Sections ✅ (9) · CTA Sections ·
Bento Grids · Pricing Sections · Header Sections · Newsletter Sections · Stats ·
Testimonials · Blog Sections · Contact Sections · Team Sections · Content Sections ·
Logo Clouds · FAQs · Footers

**Elements** — Headers ✅ (6) · Flyout Menus · Banners

**Feedback** — 404 Pages

**Page Examples** — Landing Pages · Pricing Pages · About Pages

## Phase 2 — Application UI

Confirmed from the reference so far:

**Application Shells** — Stacked Layouts · Sidebar Layouts · Multi-Column Layouts
**Headings** — Page Headings · Card Headings · Section Headings
**Data Display** — Description Lists · Stats · Calendars
**Lists** — to confirm

Remaining groups (Forms, Overlays, Elements, Navigation, Layout, Page Examples)
to be confirmed against the reference before they are scheduled, rather than
guessed at now.

## Phase 3 — Ecommerce

Deferred until Marketing and Application UI are complete.

---

## Working loop

Per subcategory:

1. You send the preview image(s).
2. I build 3 blocks to match, each responsive, accessible, and light/dark clean.
3. They are added to the manifest and appear on the subcategory page.
4. I verify: build passes, previews render at all three breakpoints in both
   themes, and the registry serves each block with its source inlined.

## Rules every block follows

These exist because the retired set violated all of them, and each violation
reached users:

- **`"use client"` on anything interactive.** Missing it breaks the block in
  every App Router project, not just here.
- **Default export**, matching the file name.
- **Every prop optional, with realistic defaults.** A block must render the
  moment it is installed — no required props, no placeholder text.
- **Light and dark both correct.** Not "works in dark, untested in light".
- **Self-contained.** No imports outside `react`, `lucide-react`, and the
  block's own file.
- **Registry-verified.** If it is listed, its source is served — an index that
  advertises what it cannot deliver is worse than a shorter index.
- **Interactive blocks must be driven, not screenshotted.** A dropdown looks
  correct in every screenshot taken while it is closed. Open it, at both themes
  and all three breakpoints, before calling it done.
- **Dismissable by Escape and by an outside click.** A menu that only closes by
  clicking its own trigger is a trap on touch and unusable by keyboard.
- **A short block sets `previewHeight`.** A header in the default 660px frame is
  mostly empty space, and padding the block itself with a fake page stub would
  ship that stub into everyone's project.
