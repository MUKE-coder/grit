# Grit UI

A shadcn-compatible registry of 100 React components — marketing, SaaS, ecommerce,
auth and app layout. Deployed at **https://ui.gritframework.dev**.

Components install into **any** React project with Tailwind. You do not need to use
the Grit framework:

```bash
npx shadcn@latest add https://ui.gritframework.dev/r/hero-split-01.json
```

Inside a Grit project the CLI wraps the same registry and picks the right app for
your architecture:

```bash
grit ui add hero-split-01
```

---

## Layout

```
ui/
├── registry/            100 component sources (.tsx) — the source of truth
├── registry.meta.json   name, title, description, category, dependencies
├── lib/
│   ├── registry.ts      builds the shadcn registry payloads
│   ├── tokens.ts        the palette, shared by the site and every registry item
│   └── component-map.ts GENERATED — run `pnpm gen`
├── app/
│   ├── page.tsx         gallery
│   ├── c/[name]/        per-component page: preview, source, install command
│   ├── preview/[name]/  bare render, embedded by the gallery iframes
│   ├── install/         setup guide
│   └── r/[name]/        the registry endpoints
└── Dockerfile           what Dokploy builds
```

## Adding a component

1. Drop `my-component-01.tsx` into `registry/`. Export a **default** function.
2. Add an entry to `registry.meta.json` (`name`, `title`, `description`,
   `category`, `dependencies`).
3. Run `pnpm gen` — it regenerates the component map and warns about any
   mismatch between the metadata and what is on disk.
4. `pnpm dev` and check `/c/my-component-01`.

Two rules the whole registry depends on:

- **Mark interactive components `"use client"`.** Anything using `useState`,
  `useEffect` or an event handler needs it, or it breaks in every consumer's App
  Router project — not just here.
- **Give every prop a default.** An installed component should render immediately
  with sample content the user can replace. A component that needs five props
  before it shows anything is a component nobody evaluates.

## Endpoints

| Path | Purpose |
|---|---|
| `/r/registry.json` | index of every component |
| `/r/<name>.json` | one registry item, with its source inlined |
| `/r/<name>` | same, extension optional |

Each item carries the component source in `files[0].content`, plus `cssVars` and a
`tailwind` colour scale. Both are needed: `cssVars` supplies the palette, and the
Tailwind scale is what makes classes like `text-text-muted` and `bg-bg-elevated`
resolve. `shadcn add` merges both automatically.

## Local development

```bash
pnpm install
pnpm dev          # http://localhost:3100
pnpm build        # production build (runs pnpm gen first)
```

To test an install against your local registry, point `shadcn` at it:

```bash
npx shadcn@latest add http://127.0.0.1:3100/r/hero-split-01.json
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

**4. Domain**

Add domain `ui.gritframework.dev`, container port `3100`, HTTPS on with Let's
Encrypt. In Cloudflare, point the record at your Dokploy host. If you use the
orange-cloud proxy, set SSL/TLS mode to **Full (strict)** so Cloudflare and
Traefik agree — Flexible causes a redirect loop.

**5. Verify the deploy**

```bash
curl -s https://ui.gritframework.dev/r/registry.json | head -c 200
```

The container's health check hits `/r/registry.json` rather than `/`, because a
site that renders while serving a broken registry is not healthy in any way that
matters to the people using it.

---

MIT licensed. Part of the [Grit Framework](https://gritframework.dev).
