# Thumbnail prompt — reference

## Accent palette

Each accent gradients **into deep navy `#0B1120`**. Keep the brand's primary as the
anchor and rotate the rest across a series.

| Name | Hex | Feels like |
|------|-----|-----------|
| Sky blue *(Grit primary)* | `#0EA5E9` | default, product, announcements |
| Cyan | `#22D3EE` | getting started, speed |
| Indigo | `#6366F1` | code, generation, deep dives |
| Violet | `#8B5CF6` | AI, advanced |
| Slate | `#334155` | comparisons, neutral/serious |
| Amber | `#F59E0B` | performance, "vs", energy |
| Emerald | `#10B981` | security, auth, "it works" |
| Teal | `#14B8A6` | data, storage |
| Rose | `#F43F5E` | warnings, hot takes, breaking |

## Hero elements (right third — pick ONE)

- **Code window** — a glossy dark editor (macOS traffic-light dots) showing 3–6
  lines of the relevant language; faint accent glow behind it.
- **Terminal window** — dark terminal with a `$ command` and a green `✓` result
  line. Great for CLI / install / deploy topics.
- **Stacked layer chips** — a 3D stack of metallic, labeled slabs (e.g. AUTH,
  STORAGE, JOBS) — the top one checkmarked, the rest recessed. Says "all included."
- **3D metallic object** — a single glossy object that *is* the metaphor: a padlock
  (security), gears (performance), a database cylinder (data), a rocket (deploy).
- **Split composition** — canvas split by a thin glowing seam; brand on one side,
  the alternative on the other; a big bold `vs.` on the seam. For comparisons.
- **Floating config/JSON panel** — a small light panel with checkmarked lines
  ("Services provisioned", "Types synced"). For config / IaC / "it's done" beats.

## Topic → hero map

| Article type | Hero |
|--------------|------|
| Awareness / "what is X" | Stacked layer chips |
| Getting started / install | Terminal window |
| CRUD / codegen | Two code windows with a connector line |
| Config / API / types | Code window (or floating config panel) |
| Comparison ("X vs Y") | Split composition with `vs.` |
| Auth / security | 3D padlock or shield + chips (JWT, OAuth, 2FA) |
| Performance / speed | Gears / speedometer |
| Deploy / ship | Terminal + rising rocket/arrow, green "✓ Live" |

## Worked examples

**Input:** topic = "Full-stack CRUD in one command", brand = Grit Framework.

```
Generate a 1200×630 landscape social card, premium developer-tool aesthetic.
Background: a bold gradient from indigo (#6366F1) into deep navy (#0B1120),
overlaid with a subtle, evenly-spaced dot-grid pattern.
Headline: the words "Full-stack CRUD\nin one command" in a heavy geometric
sans-serif, pure white, set inside a solid deep-navy highlight block anchored to
the lower-left third.
Right third: two overlapping glossy code windows fanned out — the front titled
product.go showing Go, the back titled use-products.ts showing TypeScript, a thin
glowing connector line between them, soft studio shadows with a faint indigo glow.
Top-left: place the attached Grit logo mark with the wordmark "Grit Framework"
next to it, small and high-contrast, matching a Neon-style logo lockup.
High contrast, lots of breathing room, no people, no lens flare, no 3D text.
```

> Attach your Grit logo PNG when you run this prompt.

---

**Input:** topic = "Grit vs. Next.js", brand = Grit Framework.

```
Generate a 1200×630 landscape social card, premium developer-tool aesthetic.
Background: a bold split composition — left deep navy (#0B1120), right muted dark
slate — with a subtle evenly-spaced dot-grid across the whole canvas and a thin
glowing sky-blue (#0EA5E9) seam down the divide.
Headline: a big bold white "vs." on the seam, "Grit" in white on the left and
"Next.js API routes" in grey on the right, heavy geometric sans-serif.
Right third: on the Grit side a small glossy stack of layer chips (API, WEB,
ADMIN); on the Next side a single lone code-window chip — implying "one API vs.
many pieces", soft studio shadows.
Top-left: place the attached Grit logo mark with the wordmark "Grit Framework"
next to it, small and high-contrast, matching a Neon-style logo lockup.
High contrast, lots of breathing room, no people, no lens flare, no 3D text.
```

> Attach your Grit logo PNG when you run this prompt.
