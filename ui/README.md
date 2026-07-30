# Grit UI

A shadcn-compatible registry of React **UI blocks** — full page sections built with
stock Tailwind CSS, organised by category and subcategory. Deployed at
**https://ui.gritframework.dev**.

Blocks install into **any** React project with Tailwind. You do not need to use the
Grit framework:

```bash
npx shadcn@latest add https://ui.gritframework.dev/r/marketing-hero-sections-simple-centered.json
```

Inside a Grit project the CLI wraps the same registry and picks the right app for
your architecture:

```bash
grit ui add marketing-hero-sections-simple-centered
```

Registry names are flat and carry the full path — `<category>-<subcategory>-<block>`
— because a block slug alone is not unique. "stats" exists under both Marketing and
Application UI, and two different blocks cannot share one registry name.

See [ROADMAP.md](./ROADMAP.md) for what is built and what is next.

---

## Layout

```
ui/
├── registry/
│   ├── catalog.ts               THE library definition: categories,
│   │                            subcategories, blocks. Single source of truth.
│   └── <category>/<subcategory>/<block>.tsx    the block sources
├── lib/
│   ├── blocks.ts                reads block source, counts, base URL
│   ├── block-map.ts             GENERATED — run `pnpm gen`
│   ├── highlight.ts             build-time Shiki, both themes
│   ├── analytics.ts             Zenith config, public + dashboard halves
│   ├── install-counts.ts        reads per-block install counts back out
│   └── format-count.ts          1234 -> "1.2k"
├── app/
│   ├── (site)/                  landing, category, subcategory pages
│   ├── (preview)/preview/[name]/  bare render, embedded by the viewer iframes
│   ├── r/[name]/                the registry endpoints
│   ├── zenith/                  the analytics dashboard (password gated)
│   └── global-not-found.tsx     the 404
├── scripts/generate-block-map.mjs
└── Dockerfile                   what Dokploy builds
```

Two root layouts, on purpose. `(site)` carries the theme provider and the
analytics tracker; `(preview)` carries neither, so a preview iframe renders in the
theme its query string asks for instead of reading the site's shared localStorage,
and embedding twelve previews does not send twelve pageviews.

## Adding a block

1. Add the block to its subcategory in `registry/catalog.ts`.
2. Create `registry/<category>/<subcategory>/<slug>.tsx`. Export a **default**
   function.
3. Run `pnpm gen` to regenerate the block map.
4. `pnpm dev` and check `/<category>/<subcategory>`.

Rules the whole registry depends on:

- **Mark interactive blocks `"use client"`.** Anything using `useState`,
  `useEffect` or an event handler needs it, or it breaks in every consumer's App
  Router project — not just here.
- **Give every prop a default.** An installed block should render immediately with
  sample content the user can replace. A block that needs five props before it
  shows anything is a block nobody evaluates.
- **Stock Tailwind only.** No design tokens, no config to patch. A block has to
  look right in a project that has never heard of Grit.
- **Check both themes.** Light and dark are both first-class; the viewer has a
  toggle and reviewers use it.

## Endpoints

| Path | Purpose |
|---|---|
| `/r/registry.json` | index of every block |
| `/r/<name>.json` | one registry item, with its source inlined |
| `/r/<name>` | same, extension optional |

Each item carries the block source **inlined** in `files[0].content`. This is not
optional for a remote registry: with a bare `path` and no content, `shadcn add`
writes an empty file and reports success.

## Local development

```bash
pnpm install
pnpm dev          # http://localhost:3100
pnpm build        # production build (runs pnpm gen first)
```

Two assertion scripts, run directly — one pure function each is not worth a test
runner:

```bash
npx tsx lib/format-count.test.ts
npx tsx lib/install-counts.test.ts
```

To test an install against your local registry, point `shadcn` at it:

```bash
npx shadcn@latest add http://127.0.0.1:3100/r/marketing-hero-sections-simple-centered.json
```

---

## Deploying with Dokploy

The site is a standalone Next.js container. Nothing else — no database, no Redis.
The registry is generated at build time and served as static files.

**1. Create the application**

In Dokploy: *Create Application* → source **GitHub** → repo `MUKE-coder/grit`,
branch `main`.

**2. Point it at this directory**

| Setting | Value |
|---|---|
| Build type | `Dockerfile` |
| Docker file | `ui/Dockerfile` |
| Docker context path | `ui` |
| Watch paths | `ui/**` |

The **watch path matters**: without it every framework commit rebuilds the site,
and this repo ships several releases a week.

**3. Environment**

| Variable | Value |
|---|---|
| `NEXT_PUBLIC_SITE_URL` | `https://ui.gritframework.dev` |
| `PORT` | `3100` |

`NEXT_PUBLIC_SITE_URL` is baked into the install commands the registry hands out,
so it must be the public origin. Set it wrong and every copied command points at
the wrong host.

**4. Analytics (optional)**

Add these under **Build-time Arguments** — *not* only Environment. They are read
while the pages are prerendered, and set only at runtime the tracker never
appears in the HTML no matter what the container says:

| Build argument | Value |
|---|---|
| `ZENITH_URL` | `https://analytics.gritframework.dev` |
| `ZENITH_SITE_KEY` | `zk_o84KgGIkU3QNkH4bC4pURnxsoeFwbmq3OfL1PeIowr4` |

Both are public: they ship in the snippet on every page and the site key can only
write events. The Dockerfile also re-declares them in the runtime stage, because
`/zenith` and the install-count refresh read `ZENITH_URL` at request time.

Then add the secrets under **Environment** (runtime only — never build args):

| Variable | Purpose |
|---|---|
| `ZENITH_API_KEY` | Reads analytics. Powers `/zenith` and the install counts. |
| `ZENITH_PW_HASH` | bcrypt hash gating `/zenith`. |
| `ZENITH_JWT_SECRET` | Signs the `/zenith` session cookie. |

With all three absent the site is unaffected: `/zenith` answers 503 and the
install-count badges stay hidden rather than showing zeros. The counts appear
within 15 minutes of `ZENITH_API_KEY` being set — the read is cached with a
revalidate window, so no redeploy is needed.

**5. Domain**

Add domain `ui.gritframework.dev`, container port `3100`, HTTPS on with Let's
Encrypt. In Cloudflare, point the record at your Dokploy host. If you use the
orange-cloud proxy, set SSL/TLS mode to **Full (strict)** so Cloudflare and
Traefik agree — Flexible causes a redirect loop.

**6. Verify the deploy**

```bash
curl -s https://ui.gritframework.dev/r/registry.json | head -c 200
```

The container's health check hits `/r/registry.json` rather than `/`, because a
site that renders while serving a broken registry is not healthy in any way that
matters to the people using it.

---

MIT licensed. Part of the [Grit Framework](https://gritframework.dev).
