---
title: "Why I built Grit"
subtitle: "I kept rebuilding the same backend on every project — auth, roles, audit trails, rate limiting, storage, OTP, email, cache, jobs — before ever touching the actual product. So I combined all of it into one thing, on Go."
series: "The Daily Grit"
edition: 1
date: 2026-07-01
readingTime: "5 min"
author: "Muke JohnBaptist"
tags: [grit, golang, react, full-stack, meta-framework, founder-story]
canonical: "https://gritframework.dev/blog/why-i-built-grit"
---

## The same setup, every single time

Before Grit, I built full-stack apps the way most people do. Sometimes it was the
**MERN** stack — Express + React. Sometimes **Hono + React**. Sometimes **Next.js**
doing the whole thing. Different tools, same story.

Because no matter which stack I reached for, every serious project needed the exact
same scaffolding before it could do anything useful:

- **Auth with roles.** And the moment you have roles, you need an **audit trail** —
  a record of who did what.
- Once you have a login, you're a target — so you need **rate limiting** and
  brute-force protection.
- **99% of real apps deal with files**, so you need **file storage**.
- Auth wants **OTP**, and OTP means you need **email**.
- Then you start optimizing: **Redis cache**, then **async/background jobs**.
- On the frontend, **React Query** to keep server state from turning into chaos.

And here's the part that wore me down: after *all of that*, I still hadn't written
a single line of my actual product. I was just… setting up.

Then came the real work — and it was more of the same grind. Hand-rolling **CRUD**
for every model. Building an **admin panel** from scratch. And **syncing types**
between the backend and the frontend so they didn't silently drift apart. Every
project. Every time. It was a lot, and it ate weeks I never got back.

## Getting off the JavaScript-backend treadmill

For a while I chased the "this one's faster" promise. Every few months another
JavaScript backend showed up claiming to be *the* fast one. I tried them all. And
every single time, I was juggling the exact same problems — just with new syntax.

I wanted off that treadmill. I also wanted something the churn couldn't give me: a
language that **big, serious enterprises respect and trust** — not the framework of
the month. That's why I chose **Go**.

## I tried starter kits first (they weren't the answer)

Before I landed on Go, I went down the starter-kit road. I built a *lot* of them —
my own boilerplates to copy-paste from at the start of each project. But starter
kits have their own trap: they rot. Every project forks into a slightly different
snowflake, they drift apart, and fixing something in one never fixes it in the
others. A starter kit is a **snapshot**. What I actually wanted was a **framework**.

## So I built Grit

Grit is everything on that list — combined into one thing that genuinely feels like
a superpower. It fuses a **Go** backend (Gin + GORM) with a **React** frontend
(Next.js or Vite + TanStack Router) and a Filament-style **admin panel** — in one
monorepo with shared types, wired together and hardened for production.

Everything I used to rebuild by hand is already there on day one:

- **Auth** — JWT, OAuth2 social login, 2FA — plus **RBAC** roles
- **Audit trail** — tamper-evident, SHA-256 hash-chained
- **Rate limiting** + brute-force protection
- **File storage** — S3 / R2 / MinIO
- **Email + OTP**
- **Redis cache**, **background jobs**, and **cron**
- **React Query hooks**, generated for you
- **Types synced** between Go and TypeScript, automatically

## The part that feels like magic

Scaffold an entire project in one command:

```bash
grit new my-app --triple
```

Then describe a resource, and Grit writes the whole vertical slice — backend *and*
frontend — at once:

```bash
grit generate resource Product --fields "name:string,price:float,stock:int"
```

Out comes the Go model, service, handler, and routes; the Zod schema and TypeScript
types; the React Query hooks; and an admin page. The CRUD grind — gone. The admin
panel — generated. The types — already in sync.

## Install it and see for yourself

```bash
# macOS / Linux
curl -fsSL https://gritframework.dev/install.sh | sh

# Windows (PowerShell)
iwr -useb https://gritframework.dev/install.ps1 | iex
```

Prefer Go? `go install github.com/MUKE-coder/grit/v3/cmd/grit@latest`. Then:

```bash
grit new my-app --triple
cd my-app && docker compose up -d
grit start
```

Five minutes later you have a Go API, a web app, and an admin panel running — auth,
storage, jobs, and all the rest already wired.

That's the whole idea, and the reason Grit exists: **stop rebuilding the plumbing.
Start on the product.**

Tomorrow in **The Daily Grit**: your first Grit app, step by step.

*Go + React. Built with Grit.*

---

### Thumbnail prompt (Gemini)

> Generate a 1200×630 landscape social card, premium developer-tool aesthetic.
> **Background:** a bold gradient from sky-blue (#0EA5E9) at top-left to deep navy
> (#0B1120) at bottom-right, overlaid with a subtle, evenly-spaced dot-grid pattern.
> **Headline:** the words **"Stop rebuilding\nthe backend."** in a heavy geometric
> sans-serif, pure white, set inside a solid deep-navy highlight block anchored to
> the lower-left third. **Right third:** a clean, glossy 3D stack of metallic
> chip-like layers labeled (top to bottom) AUTH, AUDIT, STORAGE, JOBS, CACHE — the
> top "AUTH" layer highlighted with a small checkmark, the rest slightly recessed,
> soft studio shadows, implying "it's all already included." **Top-left:** place the
> attached Grit logo mark with the wordmark "Grit Framework" next to it, small and
> high-contrast, matching a Neon-style logo lockup. High contrast, lots of breathing
> room, no people, no lens flare, no 3D text.
