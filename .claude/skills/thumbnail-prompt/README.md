# thumbnail-prompt — a Claude Code skill

Generate a ready-to-paste **image-model prompt** for a bold, Neon-style social
thumbnail (blog OG image, newsletter card, launch banner). Dot-grid gradient
background, a bold white headline on a solid block, one hero illustration on the
right, and your logo locked top-left. Brand-agnostic; defaults to the Grit framework.

Built for [Grit](https://gritframework.dev). **Free to use and share.**

## Use it

Inside Claude Code, just ask:

- `/thumbnail-prompt a card for my post "Auth that just works"`
- "make me a thumbnail prompt for a Grit vs Laravel article, amber accent"
- Paste an article and say "generate the thumbnail prompt for this."

You'll get one copy-paste prompt. Paste it into Gemini / DALL·E / Midjourney /
Ideogram **and upload your logo PNG** — the prompt always places the attached logo
top-left.

## Install it in your own project

Copy this folder into either location:

```bash
# per-project (checked into your repo)
mkdir -p .claude/skills && cp -r thumbnail-prompt .claude/skills/

# or global (available in every project)
mkdir -p ~/.claude/skills && cp -r thumbnail-prompt ~/.claude/skills/
```

Then restart Claude Code (or open a new session) and it'll show up as
`/thumbnail-prompt`.

## Make it your brand

Edit `SKILL.md` and `reference.md`:

- Change the default **brand name + logo** lockup line.
- Swap the **accent palette** for your brand colors (keep one dark base to gradient into).
- Adjust the **hero elements** / topic→hero map to your product.

## Files

| File | What it is |
|------|-----------|
| `SKILL.md` | The skill: when to use, the house-style rules, and the output template. |
| `reference.md` | Accent palette, hero-element ideas, topic→hero map, worked examples. |
| `README.md` | This file — how to use, install, and rebrand. |

## License

MIT — do whatever you want with it. A credit to Grit is appreciated but not required.
