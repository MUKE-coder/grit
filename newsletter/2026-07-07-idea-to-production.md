---
title: "From idea to production with one command: grit deploy"
subtitle: "Build, upload, systemd, and Caddy with automatic HTTPS — or push to git and ship the bundled Dockerfile. Your call."
series: "The Daily Grit"
edition: 7
date: 2026-07-07
readingTime: "5 min"
author: "Muke JohnBaptist"
tags: [grit, deploy, devops, caddy, docker]
canonical: "https://gritframework.dev/blog/idea-to-production"
---

Week one ends where every project wants to end: in production. Grit gives you two
honest paths — a self-hosted one-liner, or a container you run anywhere.

## Option A — `grit deploy` to your own server

Point it at a box you own and Grit does the rest:

```bash
grit deploy --host deploy@server.com --domain myapp.com
```

Under the hood it:

```
→ Cross-compiles the Go binary for Linux (CGO_ENABLED=0)
→ Builds the frontend if present (pnpm build)
→ Uploads the binary via SCP
→ Creates a systemd service with auto-restart
→ Configures Caddy as a reverse proxy with automatic Let's Encrypt TLS
✓ Live at: https://myapp.com
```

No YAML, no control plane, no per-seat pricing. One long-lived Go process behind
Caddy — cheap and boring, in the best way. Set `DEPLOY_HOST` / `DEPLOY_DOMAIN` /
`DEPLOY_KEY_FILE` env vars and it's just `grit deploy`.

## Option B — the bundled Dockerfile

Prefer a platform (Railway, Render, Fly, a k8s cluster)? Every Grit project ships
a production `Dockerfile` and `docker-compose.prod.yml`. Push to git and let your
platform build it — the same artifact, your workflow.

## Single-binary mode is a deploy superpower

If you scaffolded with `--single`, your React SPA is embedded in the Go binary via
`go:embed`. That means **one file** to ship — no separate node server, no static
host to coordinate. Copy the binary, run it, done.

## What's already production-grade

You're not bolting on reliability at the end — it shipped with the scaffold:

- **Security:** OWASP-2025 headers, rate limiting, CSRF, tamper-evident audit log
- **Observability:** request tracing, DB monitoring, p50/p95/p99 latency
- **Performance:** gzip, connection pooling, graceful shutdown
- **Health checks** for load balancers and uptime monitors

## The whole week, in one arc

In seven mornings you went from *never heard of Grit* to a shipped app:

1. **What Grit is** — Go + React, batteries included
2. **First app in 5 minutes** — scaffold and run
3. **Full-stack CRUD** — one command, both sides
4. **vs. Next.js** — when you want a real backend
5. **vs. Laravel** — Rails-style DX on Go
6. **Auth** — JWT, OAuth, 2FA out of the box
7. **Deploy** — idea to production, one command

That's the pitch in practice: **spend your time on product, not plumbing.**

Start your own: `curl -fsSL https://gritframework.dev/install.sh | sh`

*Go + React. Built with Grit.*

---

### Thumbnail prompt (Gemini)

> Generate a 1200×630 landscape social card, premium developer-tool aesthetic.
> **Background:** a bold gradient from sky-blue (#0EA5E9) into deep navy (#0B1120),
> overlaid with a subtle evenly-spaced dot-grid pattern. **Headline:** the words
> **"Idea to production,\none command"** in a heavy geometric sans-serif, pure
> white, inside a solid deep-navy highlight block anchored lower-left. **Right
> third:** a glossy dark terminal window showing `$ grit deploy --domain myapp.com`
> with a green "✓ Live at https://myapp.com" line, and a subtle 3D rocket or
> upward arrow rising out of it, soft sky-blue glow. **Top-left:** place the
> attached Grit logo mark with the wordmark "Grit Framework" beside it, small,
> high-contrast, Neon-style lockup. High contrast, generous whitespace, no people,
> no lens flare, no 3D text.
