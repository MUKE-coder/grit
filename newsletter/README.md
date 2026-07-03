# The Daily Grit — Newsletter Kit

The official daily newsletter of the **Grit framework**. One short, practical read
every morning at **5:00 AM**, cross-posted to **LinkedIn**, **Medium**,
**daily.dev**, and the **docs blog**.

> **Name:** **The Daily Grit** — daily cadence, on-brand, and a wink at "the daily
> grind." (Alternates if you want to swap: *Built with Grit*, *Ship with Grit*.)

---

## LinkedIn newsletter setup

Paste these into the **Create a newsletter** dialog:

| Field | Value |
|-------|-------|
| **Newsletter title** | `The Daily Grit` |
| **How often** | `Daily` |
| **Description** | `The official newsletter of the Grit framework — Go + React, batteries included. Every morning, a 5-minute hands-on read: getting started, full-stack CRUD in one command, framework comparisons, and production tips. Ship something before your coffee's cold.` |
| **Logo (300×300)** | Grit "G" mark on the sky-blue → navy brand gradient |

---

## Editorial focus (first phase = awareness)

Order of priority while the audience is small: **awareness → getting started →
simple CRUD → comparisons → batteries → deploy**. Keep every edition ≤ 5-minute
read (≈700–1000 words), one concrete takeaway, one copy-pasteable command.

## Week 1 schedule

| # | Date | Article | Angle | Status |
|---|------|---------|-------|--------|
| 1 | 2026-07-01 | [Why I built Grit](2026-07-01-why-i-built-grit.md) | Awareness | ✅ published |
| 2 | 2026-07-02 | [Build your first Grit app](2026-07-02-build-your-first-grit-app.md) — a Category → Product store | Getting started · CRUD · files · relationships | draft |
| — | — | *Render the store on the web app* (teased at the end of #2) | Frontend | planned |
| — | — | [Grit vs. Next.js](2026-07-04-grit-vs-nextjs.md) | Comparison | draft |
| — | — | [Grit vs. Laravel](2026-07-05-grit-vs-laravel.md) | Comparison | draft |
| — | — | [Auth that just works](2026-07-06-auth-that-just-works.md) | Batteries | draft |
| — | — | [Idea to production with grit deploy](2026-07-07-idea-to-production.md) | Deploy | draft |

> Editions 2 & 3 were merged into a single hands-on build (#2). The remaining drafts
> keep their files; renumber their `edition:` frontmatter when each is published.

---

## Thumbnail house style (Gemini prompt guide)

Every edition ships a 1.91:1 social thumbnail in the **Neon style**. For each
article, paste that article's **Thumbnail prompt** into Gemini **and upload the
Grit logo PNG** — the prompt always instructs Gemini to place the attached logo +
"Grit Framework" wordmark in the top-left, exactly like Neon's "◈ NEON" lockup.

**Shared rules (baked into every prompt):**
- **Ratio:** 1200×630 (or 2400×1260), landscape social card.
- **Background:** a bold Grit-blue gradient (sky `#0EA5E9` / `#0284C7` → deep navy
  `#0B1120`) with a subtle, evenly-spaced **dot-grid** pattern across it.
- **Headline:** short, bold, white, set in a **solid highlight block** (navy or
  sky-blue) anchored bottom-left — the way Neon sets white text on a black/blue block.
  Bold geometric sans-serif.
- **Right third:** one clean hero element — a stylised **code window**, a **3D
  metallic object**, or **stacked layer chips** — never clutter.
- **Top-left:** the **attached Grit logo mark + "Grit Framework"** wordmark, small,
  clean, high-contrast.
- **Feel:** premium developer-tool aesthetic; high contrast; lots of breathing room;
  no stock-photo people, no lens flare, no 3D text.
- Accent color is varied per edition (within the brand family) so the feed looks
  cohesive but not repetitive.
