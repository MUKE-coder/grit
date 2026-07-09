---
title: "Build a desktop inventory app with Grit: generate resource → shippable .exe"
subtitle: "grit new-desktop scaffolds a Wails app that's secretly a hybrid — a real REST API embedded in the binary. Point grit generate resource at it and you get typed CRUD screens, native file uploads, and a searchable relationship picker. We build a Category → Product inventory system, photos and all, offline-first on SQLite."
series: "The Daily Grit"
edition: 8
date: 2026-07-08
readingTime: "10 min"
author: "Muke JohnBaptist"
tags: [grit, desktop, wails, go, react, code-generation, inventory, tutorial]
canonical: "https://gritframework.dev/blog/build-desktop-app-with-grit"
---

We've built a **web** store and a **mobile** store with Grit. Same command each
time — `grit generate resource` — describe a model and get the Go API, the
hooks, and the screens.

Today: a **desktop** app. A real native window you can hand someone as an
installer, that runs offline on a local SQLite file — or points at Postgres when
you want it to. We'll build a **Category → Product inventory system**: browse
products, add them with a **photo**, pick their category from a dropdown, and
watch stock. And the interesting part is what `grit new-desktop` quietly gives
you: a Wails shell with a **full REST API embedded in the same binary**, which is
what makes file uploads work without a server.

Let's build.

## 1. Install or update Grit

Desktop file uploads and the hybrid API landed in the **v3.31.79–v3.31.83** line,
so grab the latest.

```bash
# Install (macOS / Linux)
curl -fsSL https://gritframework.dev/install.sh | sh

# Install (Windows PowerShell)
iwr -useb https://gritframework.dev/install.ps1 | iex

# Already have Grit? Update in place:
grit update
```

Confirm you're on **v3.31.83 or newer**:

```bash
grit version
```

You'll also need the [Wails](https://wails.io) toolchain (`wails doctor` should
be green), Go 1.21+, and Node 18+.

**One more, for the installer step:** to package a Windows `.exe` installer in
step 7 you need **NSIS**. Install it now so the build doesn't stop halfway:

```powershell
winget install NSIS.NSIS      # or: choco install nsis  /  scoop install nsis
```

After installing, add the NSIS folder (usually `C:\Program Files (x86)\NSIS`) to
your PATH so `makensis` resolves — `grit package` checks for it up front and will
tell you if it's missing. You don't need NSIS for `wails dev` or for
`grit package --no-installer` (raw binary, no installer).

## 2. Scaffold the desktop app

```bash
grit new-desktop inventory
cd inventory
```

`new-desktop` scaffolds a standalone [Wails](https://wails.io) app: a Go backend,
a React frontend (TanStack Router + Tailwind), local auth, and **GORM Studio** for
poking at the database. By default it runs on **SQLite** (`app.db`, zero setup) —
set `DB_DRIVER=postgres` and `DB_DSN=...` in `.env` when you want Postgres instead.

### The part that isn't obvious: it's hybrid

A Wails app talks to Go through generated bindings — great for CRUD, useless for
file uploads (there's no `<input type="file">` story, no URL to load an image
from). So a Grit desktop app **also embeds a real Gin REST API in the same
binary**. You'll find three new packages:

```
internal/api/       # Gin router: POST /api/uploads, GET /uploads/*, /api/health
internal/storage/   # writes uploads into the OS app-data dir
internal/files/     # FileRef — the same JSON shape web & mobile store
```

The router is mounted **twice** (see `main.go`):

- as the Wails **asset-server handler**, so the webview can `fetch("/api/uploads")`
  and render `<img src="/uploads/photo.jpg">` **same-origin** — no port to find,
  no CORS to configure;
- on **`127.0.0.1:34999`**, so `curl` or any other client can hit the same API.
  (The Wails dev window itself runs on `localhost:34115` during `wails dev` —
  the embedded API deliberately uses a different port so they never collide.)

That's the whole trick behind uploads below.

## 3. Generate the models

A product **belongs to** a category, so create the parent first.

**Category:**

```bash
grit generate resource Category --fields "name:string,slug:slug"
```

**Product** — with a photo and a category link:

```bash
grit generate resource Product --fields "name:string,sku:string,quantity:int,price:float,photo:file:image,category:belongs_to:Category"
```

### What each command generated

```
✓ internal/models/<name>.go                       # GORM model (photo is a files.FileRef column)
✓ internal/service/<name>.go                       # CRUD service (list/get/create/update/delete + search)
✓ frontend/src/routes/_layout/<plural>.index.tsx   # list table (with photo thumbnails)
✓ frontend/src/routes/_layout/<plural>.new.tsx     # create form
✓ frontend/src/routes/_layout/<plural>.$id.edit.tsx# edit form
```

It also **injects** the Wails bindings into `app.go`, the input type into
`types.go`, and a link into the sidebar — so the resource is reachable the moment
you restart. TanStack Router is file-based, so those route files *are* the routes.

Two field types earned their keep here:

- **`photo:file:image`** → the form renders a native file picker that uploads to
  the embedded `/api/uploads`, previews the image, and the list shows a
  thumbnail. The model column is a `files.FileRef` (url + name + mime + size),
  stored as JSON.
- **`category:belongs_to:Category`** → the form renders a **searchable dropdown**
  populated from `GetCategories()`, so you pick a category by name instead of
  typing an id.

> **Prices, honestly:** `price:float` stores exactly what you enter. No cents
> trickery — `19.99` is `19.99`.

## 4. Run it

```bash
wails dev
```

Wails compiles the Go binary, generates the TypeScript bindings, and opens the
window with hot reload. Migrations run on boot (GORM AutoMigrate), so the
`categories` and `products` tables are created for you. On first run a default
admin is seeded — log in with **admin@example.com / admin123** — and
**Categories** + **Products** are already in the sidebar.

## 5. How an upload actually flows

When you pick a photo in the Product form, the generated code does this:

```tsx
const fd = new FormData();
fd.append("file", file);
const res = await fetch("/api/uploads", { method: "POST", body: fd });
const json = await res.json();
if (res.ok) setPhoto(json.data);   // json.data is a FileRef
```

That `fetch` is same-origin — it hits the Gin router mounted on the Wails asset
server. On the Go side, `POST /api/uploads` validates the type (JPEG/PNG/GIF/WebP),
writes the file into the **OS app-data directory**…

```go
base, _ := os.UserConfigDir()             // %APPDATA% on Windows, ~/Library/… on macOS
dir := filepath.Join(base, appName, "uploads")
```

…and returns a `FileRef` whose `url` is a **relative** path like
`/uploads/1751-abc.jpg`. The form stores that ref; the list renders
`<img src={item.photo.url}>`, which the same router serves straight off disk. No
S3, no network, works on a plane.

Why app-data and not next to the binary? An installed app often lives somewhere
read-only (Program Files, /Applications). App-data is per-user and writable, so a
photo attached today survives the next app update.

## 6. A low-stock touch

The generated list is a normal React component, so small business logic is easy.
Open `frontend/src/routes/_layout/products.index.tsx` and tint low-stock rows —
add this to the quantity cell:

```tsx
<td className="px-4 py-3 text-sm">
  <span className={Number(item.quantity) <= 5 ? "text-danger font-semibold" : "text-text-secondary"}>
    {item.quantity ?? 0}
    {Number(item.quantity) <= 5 ? " · low" : ""}
  </span>
</td>
```

Everything around it — pagination, search, the CSV/PDF export buttons, delete
confirmation — is generated. You're just decorating one cell.

## 7. Ship it

One command turns the app into something you can hand to a user:

```bash
grit package
```

On Windows that produces an **NSIS installer** (a single `*-installer.exe` in
`build/bin/`) — the file a shopkeeper double-clicks to install your app. On
macOS/Linux it produces the platform binary/app bundle. `grit package` is a
friendly wrapper over `wails build`: it checks you have the toolchain
(`wails`, plus `makensis` for the Windows installer), builds, and tells you where
the artifact landed. Skip the installer with `--no-installer`, or cross-compile
with `--platform darwin/arm64`.

```bash
grit package                    # installer for the current OS
grit package --no-installer     # just the raw binary
grit package --platform windows/amd64
```

The result is a single artifact that carries its own SQLite database and uploads
folder — no Docker, no server. Point it at Postgres later by shipping a different
`.env`; the code doesn't change. (For a full versioned release with branded
installer art and a GitHub release, the scaffold also ships
`scripts/release-desktop.sh <version>`.)

## The takeaway

`grit new-desktop` isn't just a Wails starter. It's a hybrid app — native window
*and* embedded REST API — so `grit generate resource` can give you the same thing
it gives web and mobile: typed CRUD, real file uploads, and relationship pickers,
with almost no plumbing. The inventory system you just built is mostly generated
code; the storefront polish is the thin, fun layer on top.

*Go + React. Built with Grit.*
