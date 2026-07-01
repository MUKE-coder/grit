---
name: thumbnail-prompt
description: Generate a ready-to-paste image-model prompt (Gemini, DALL·E, Midjourney, etc.) for a bold, Neon-style social thumbnail — a blog OG image, newsletter card, or launch banner. Use whenever the user asks for a "thumbnail", "social card", "cover image", "OG image", or "banner" prompt in the dev-tool aesthetic: a dot-grid gradient background, a bold white headline on a solid color block, one hero illustration or code window on the right, and a brand logo locked to the top-left. Works for any brand; defaults to the Grit framework.
---

# Neon-style thumbnail prompt generator

Produce a single, copy-paste-ready prompt for an image model that renders a
premium developer-tool social card — the look pioneered by Neon, Vercel, and
Linear blog cards. The user pastes the prompt into their image model **and
uploads their logo PNG**; the prompt always instructs the model to place that
logo top-left.

## When to use

The user wants a thumbnail / social card / cover / OG image / banner and either
gives you a topic, a headline, or a whole article. Turn it into one polished prompt.

## Inputs (ask only if missing something essential)

- **Headline or topic** (required) — the article title, or the point of the image.
- **Brand** (default: `Grit Framework`, logo = the Grit "G" mark) — name + logo lockup.
- **Accent color** (optional) — pick one from `reference.md` if not given.
- **Hero element** (optional) — the right-side visual; pick a fitting one from
  `reference.md` if not given.
- **Ratio** (default `1200×630`, i.e. 1.91:1 social/OG) — use `1080×1080` if they
  want square, `1600×900` for a wide banner.

If the user gives a full article, read it, distill the single sharpest idea, and
write the headline yourself.

## The house style (invariants — never drop these)

1. **Ratio:** landscape 1.91:1 (`1200×630`) unless told otherwise.
2. **Background:** a bold gradient from the chosen **accent color** into **deep navy
   `#0B1120`**, overlaid with a **subtle, evenly-spaced dot-grid** pattern.
3. **Headline:** short (aim ≤ 6 words, ≤ 2 lines), **pure white**, **heavy geometric
   sans-serif**, set inside a **solid deep-navy highlight block anchored lower-left**.
   Use `\n` to control the line break.
4. **Right third:** exactly **one** clean hero element (code window, 3D metallic
   object, stacked layer chips, terminal, or a split comparison) with soft studio
   shadows or a faint accent glow. Never clutter.
5. **Top-left:** "place the **attached** [Brand] logo mark with the wordmark
   '[Brand Name]' next to it, small and high-contrast, matching a Neon-style logo
   lockup." (Always say *attached* — the user uploads the real logo.)
6. **Global:** high contrast, generous breathing room, **no people, no lens flare,
   no 3D text, no stock photos.**

## How to build the prompt

1. **Headline** — compress the topic to a punchy hook. Prefer verbs and tension
   ("Stop rebuilding the backend", "Auth that just works", "CRUD in one command").
2. **Accent** — pick from the palette in `reference.md`. Vary it across a series so
   a feed looks cohesive but not repetitive; keep the brand's primary as the anchor.
3. **Hero** — choose the element that best *shows* the idea (see the topic→hero map
   in `reference.md`): a comparison → split composition; a CLI → terminal window;
   "batteries included" → stacked layer chips; a config/API → code window.
4. **Assemble** using the template below, filling every `{...}`.

## Output template

Return this, filled in, inside a single fenced block, followed by the one-line
upload reminder — nothing else:

```
Generate a {RATIO} landscape social card, premium developer-tool aesthetic.
Background: a bold gradient from {ACCENT_NAME} ({ACCENT_HEX}) into deep navy
(#0B1120), overlaid with a subtle, evenly-spaced dot-grid pattern.
Headline: the words "{HEADLINE}" in a heavy geometric sans-serif, pure white,
set inside a solid deep-navy highlight block anchored to the lower-left third.
Right third: {HERO_DESCRIPTION}, soft studio shadows with a faint {ACCENT_NAME}
glow.
Top-left: place the attached {BRAND} logo mark with the wordmark "{BRAND_NAME}"
next to it, small and high-contrast, matching a Neon-style logo lockup.
High contrast, lots of breathing room, no people, no lens flare, no 3D text.
```

> Then add: **"Attach your {BRAND} logo PNG when you run this prompt."**

If the user asks for several (a series), output one prompt per topic, each with a
different accent from the palette, and keep the logo lockup identical across all.

See `reference.md` for the accent palette, hero ideas, the topic→hero map, and two
worked examples.
