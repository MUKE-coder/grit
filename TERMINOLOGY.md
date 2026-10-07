# Terminology

One name per thing. This file is the arbiter when two words are in circulation,
and the words in the **Canonical** column are the ones to use in code, in
comments, in docs, in commit messages and in the CLI's own output.

It exists because several of these words have two meanings in this repository,
and picking the wrong one does not produce an error. It produces a file written
into a directory the project does not have, a predicate that answers the
question next to the one you asked, or a doc page describing a layout nobody
ships. Those failures are silent, which is why they are worth a table.

Read this with [CLAUDE.md](CLAUDE.md) (the project's working rules) and
[API_DESIGN.md](API_DESIGN.md) (the HTTP contract).

---

## The shape of a project

| Canonical | Means | Avoid | Use instead |
|---|---|---|---|
| **architecture** | The choice between `single`, `double`, `triple`, `api` and `mobile`. It is the name of the field (`Options.Architecture`) and of the flag. | "tier" as a noun: "which tier is this?" | "which architecture?" |
| **triple**, **double**, **single**, **api**, **mobile** | The five values. `ArchTriple`, `ArchDouble`, `ArchSingle`, `ArchAPI`, `ArchMobile`. | "3-tier", "full-stack mode" | the value's own name |
| **triple-tier** (prose only) | Acceptable in prose and tutorials for `triple`, where the reader is counting apps. Never in code. | "tier" on its own, which does not say which of the five | "a triple project" |
| **frontend** | Which router and bundler the React apps use: `next` or `tanstack`. `Options.Frontend`. | "the UI framework" (ambiguous with Tailwind and the component set) | "the frontend" |
| **TanStack** / **Vite** | The same frontend, named by its router and by its bundler. Both are fine; `tanstack` is the constant, `--vite` is the flag, and docs say "Vite" because that is what the user types. | "SPA" as a synonym, because a Next single is also one app | "the Vite frontend" |
| **project** | One thing `grit new` produces, whatever its architecture. The root directory. | "app" when you mean the whole project | "project" |
| **app** | One deployable unit inside a project: `apps/api`, `apps/web`, `apps/admin`, `apps/expo`, `apps/desktop`. | "service", "package" | "app" |
| **package** | A workspace library under `packages/`: `packages/shared`, `packages/upload`. Not an app. | "module" (that is Go's word) | "package" |

### The trap worth naming

A **single** project's frontend is at the **project root**, not `apps/web` and
not `frontend/`. Use `webAppRoot(root, opts)`; a literal `apps/web` writes into
a directory that does not exist, and nothing fails. Only a project scaffolded
before v3.380.0 has `frontend/`, and `Options.LegacySingleFlat` is what says so.

---

## The admin panel

The panel ships in three shapes, and this is the distinction that has cost the
most. **Always ask which question you mean.**

| Canonical | Means |
|---|---|
| **admin app** | `apps/admin`: its own `package.json`, port and container. Only a **triple** has one. |
| **embedded panel** | The panel as a route group inside the web app, at `/admin`. What a **double** has. |
| **SPA section** | The panel as a section of a single project's Vite router. What a **single --vite** has. |
| **admin panel** | All three, collectively. The thing a person sees. |

| Predicate | Answers |
|---|---|
| `ShouldIncludeAdmin()` | "Is there a separate admin **application**?" Only true for a triple. |
| `HasAdminPanel()` | "Is there an admin **panel** at all, in any shape?" |
| `ShouldEmbedAdmin()` / `ShouldEmbedAdminInSPA()` | Which of the two embedded shapes. |

**Every writer that produces a screen wants `HasAdminPanel()`.** Asking the
other question is how a double shipped a user menu linking to an
account-security page that was never written.

Paths come from `adminPath`, `adminComponent`, `adminHook`, `adminLib`, which
know all three shapes. Spelling out `apps/admin` does not.

---

## Making code

| Canonical | Means | Avoid | Use instead |
|---|---|---|---|
| **scaffold** | What `grit new` does: write a whole project. The `internal/scaffold` package. | "generate" for this | "scaffold" |
| **generate** | What `grit generate resource` does: add one resource to a project that exists, by writing new files and **injecting** into existing ones. The `internal/generate` package. | "scaffold a resource" | "generate a resource" |
| **upgrade** | What `grit upgrade` does: bring an existing project to this CLI's version. | "migrate" (that is the database) | "upgrade" |
| **repair** | One targeted, idempotent rewrite of a file an upgrade cannot deliver whole, because the project may have edited it. `*_repair.go`. | "patch", "fix up" | "repair" |
| **migrate** | Database schema only. `grit migrate`. | using it for the CLI or the code | "migrate" |
| **framework-owned** | A file the framework writes and the project must not edit: nothing in a resource definition describes it, and it only makes sense alongside its siblings. It **travels on upgrade**. | "internal file", "core file" | "framework-owned" |
| **delivery** | Getting a file into a project. There are two paths, and a framework-owned file needs **both**: the new-project map, and the upgrade writer. | (none) | "delivery" |
| **writer** | A `writeXFiles(root, opts) error` function called from **both** `grit new` and `grit upgrade`. The preferred shape, because one list cannot drift from the other. | (none) | "writer" |
| **template** | Two different things, so always qualify. A **string template** is a Go raw string inside `internal/scaffold/*.go`. A **template file** is a real file under `internal/scaffold/templates/` read by `tmpl()`. | bare "template" | "string template" / "template file" |
| **injection** | Adding a line to a file the project already has, at a named anchor, so a generated resource appears in `routes.go` and the sidebar. `internal/generate/inject.go`. | "patching routes" | "injection" |

---

## Products, and what is not ours to fix

| Canonical | What it is |
|---|---|
| **Grit** | This framework and its CLI. |
| **Grit UI** | The standalone component library. **Not** part of the scaffold, and removed from it in v3.31.78: do not re-add the registry, the `UIComponent` model or `packages/grit-ui/`. |
| **Grit blueprints** | The paid starter applications built *with* Grit, in their own repositories. |
| **Grit Cloud** | The hosting product. |
| **GORM Studio** | A separate project of ours, embedded in the generated API at `/studio`. |
| **Sentinel** | A separate project of ours: rate limiting. |
| **Pulse** | A separate project of ours: request observability. |

A bug in **GORM Studio, Sentinel or Pulse gets a GitHub issue on that
repository**, not a fix here. They are dependencies, and working around them in
the scaffold hides the defect from their other users.

---

## Words with no agreed meaning yet

These are in circulation and not settled. If you need one, define it where you
use it rather than assuming the reader shares your sense of it.

- **"the site"**: sometimes the marketing pages (`app/(marketing)`), sometimes
  the whole web app. Say which.
- **"the dashboard"**: the admin's landing screen, or its whole authenticated
  area. Say which.
- **"tier"**: see above. Fine in prose with a number attached, misleading on
  its own.
