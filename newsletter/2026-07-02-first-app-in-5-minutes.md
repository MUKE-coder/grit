---
title: "Your first Grit app in 5 minutes"
subtitle: "Install the CLI, scaffold a full monorepo, and have a Go API + React apps running locally — in one sitting."
series: "The Daily Grit"
edition: 2
date: 2026-07-02
readingTime: "5 min"
author: "Muke JohnBaptist"
tags: [grit, getting-started, golang, react, tutorial]
canonical: "https://gritframework.dev/blog/first-app-in-5-minutes"
---

Yesterday we met Grit. Today we run it. By the end of this read you'll have a Go
API, a web app, and an admin panel running on your machine.

## 1. Install the CLI

```bash
# macOS / Linux
curl -fsSL https://gritframework.dev/install.sh | sh

# Windows (PowerShell)
iwr -useb https://gritframework.dev/install.ps1 | iex
```

Prefer Go? `go install github.com/MUKE-coder/grit/v3/cmd/grit@latest`. Verify:

```bash
grit version
```

You'll want **Go 1.21+**, **Node 18+**, **pnpm**, and **Docker** installed.

## 2. Scaffold a project

```bash
grit new my-app --triple
```

`--triple` is the kitchen-sink layout: a **web** app, an **admin** panel, and the
**Go API**, in one monorepo with a shared types package. (Other modes:
`--single` embeds a SPA in the Go binary, `--double` is web + API, `--api` is
API-only, `--mobile` adds Expo. Pick later — the generators are identical.)

## 3. Start the infrastructure

Grit ships a Docker Compose file with Postgres, Redis, MinIO, and Mailhog:

```bash
cd my-app
docker compose up -d
```

## 4. Run everything

```bash
grit start
```

That boots the **Go API and the frontends in parallel**, with color-prefixed logs
so you can tell who said what. Ctrl+C stops both. (If you have
[`air`](https://github.com/air-verse/air) installed, the API hot-reloads on `.go`
changes too.)

Now open:

| Surface | URL |
|---------|-----|
| Web app | http://localhost:3000 |
| Admin panel | http://localhost:3001 |
| API | http://localhost:8080 |
| API docs | http://localhost:8080/docs |
| GORM Studio (DB browser) | http://localhost:8080/studio |

## 5. Seed a login and sign in

```bash
grit migrate   # run GORM AutoMigrate
grit seed      # create an admin + demo users
```

Head to the admin panel, log in with the seeded admin, and you're looking at a
working dashboard — auth, sessions, and a data table, all already wired.

## What you just got

Without writing a line of code you have: JWT auth, an admin panel, a typed API
client, Docker infra, migrations, and a visual database browser. That's the "first
month" from yesterday's edition — done in five minutes.

**Tomorrow:** we add a real feature. One command generates a full-stack CRUD
resource — Go handler, React hook, and admin page — all at once.

*Go + React. Built with Grit.*

---

### Thumbnail prompt (Gemini)

> Generate a 1200×630 landscape social card, premium developer-tool aesthetic.
> **Background:** a bold gradient from cyan (#22D3EE) to deep navy (#0B1120),
> overlaid with a subtle evenly-spaced dot-grid pattern. **Headline:** the words
> **"Your first app\nin 5 minutes"** in a heavy geometric sans-serif, pure white,
> inside a solid deep-navy highlight block anchored lower-left. **Right third:** a
> glossy dark terminal window (macOS traffic-light dots) showing three
> monospaced lines: `$ grit new my-app --triple`, `$ docker compose up -d`,
> `$ grit start`, with a soft cyan glow behind it. A small green checkmark badge
> in the corner of the terminal. **Top-left:** place the attached Grit logo mark
> with the wordmark "Grit Framework" beside it, small and high-contrast, Neon-style
> lockup. High contrast, generous whitespace, no people, no lens flare, no 3D text.
