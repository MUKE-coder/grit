---
title: "Grit plugins: what they are, the default ones, and how to build your own"
subtitle: "A Grit plugin doesn't hide in node_modules, it generates real code into your repo: models, routes, pages, migrations. And because every install is recorded, it's the one thing most plugin systems can't do: cleanly uninstall. Here's the model, the first-party plugins, and how to write one."
series: "The Daily Grit"
edition: 11
date: 2026-07-21
readingTime: "11 min"
author: "Muke JohnBaptist"
tags: [grit, plugins, extensibility, codegen, multitenancy, go, react, architecture]
canonical: "https://gritframework.dev/blog/grit-plugins-what-they-are-and-how-to-build-one"
thumbnail: "/blog/grit_plugins.png"
---

Most frameworks let you extend them with a package: you `install` it, it lives in
your dependency folder, and at runtime it reaches in and registers itself. That's
fine until you want to *see* what it did, or *remove* it. The code isn't yours;
it's behind a version number.

Grit plugins work the other way around.

## What a Grit plugin is

**A Grit plugin generates code into your repo.** Run `grit plugin add multitenant`
and it writes real files (Go models, handlers, routes, admin pages, migrations)
directly into `apps/api`, `apps/admin`, and the rest of your project. After it
runs, the code is *yours*. You can read it, edit it, step through it in a
debugger. There's no hidden runtime library deciding your app's behavior from
inside `vendor/`.

This is the same philosophy as the rest of Grit. `grit generate resource` doesn't
give you an abstract `Resource<T>`: it writes a model, a service, a handler,
typed hooks and an admin page you own. Plugins are that idea, extended: a plugin
is just a bigger, reusable code generation.

### The part almost nobody else can do: uninstall

Here's the trick. When a plugin runs, Grit records **everything** it did, every
file it created and every injection it made into an existing file, into a
lockfile:

```
.grit/plugins.lock.json
```

So removal isn't a second pile of code the author has to write and keep in sync
with the installer. It's derived. `grit plugin remove multitenant` reads the
lockfile and **replays the install backwards**: it deletes the files it added and
reverts the injections it made, in reverse order. Commit the lockfile, and your
plugins are as reversible as a git branch.

That's the structural advantage. A framework whose plugins publish migrations
into your app can't cleanly reverse them, publishing is a one-way door. Grit's
can, because the install is a recorded, replayable transaction.

## The commands

```bash
grit plugin list                 # what's available
grit plugin info multitenant     # what a plugin does before you run it
grit plugin add multitenant      # install (refuses to run on a dirty git tree)
grit plugin remove multitenant   # replay the install backwards
```

`grit plugin add` refuses to run with uncommitted changes, so the diff it
produces is always clean and reviewable. Read the diff, commit it, and it's part
of your app like anything else.

## The default plugin: multi-tenancy

The first first-party plugin is **`multitenant`**, and it's a good example of
why some things are plugins instead of core.

Multi-tenancy is a big, opinionated architectural decision. Not every app needs
it, and the ones that do want it woven through their data layer. So it's not
baked into every `grit new`; you opt in:

```bash
grit plugin add multitenant
```

What you get:

- **Organization** and **OrganizationMember** models: one user belongs to many
  organizations, holding a different role in each.
- The active organization resolved from a request header or the session: **no
  subdomains**, so it works the same in dev, mobile, and desktop.
- **Automatic org scoping on every query**, via a GORM callback. A forgotten
  `WHERE organization_id = ?` can't leak data across tenants, because you don't
  write the `WHERE`, the callback adds it, and it **fails closed**: no active
  org, no rows.
- **Per-organization permissions** that reuse the roles system you already have.

It touches models, migrations, middleware and admin UI at once: exactly the kind
of cross-cutting change that's painful to do by hand and perfect for a recorded,
reversible plugin.

## The first-party catalog

Alongside multi-tenancy, three more first-party plugins ship today, the small,
high-value features that the best admin ecosystems treat as add-ons:

- **impersonate**: an admin signs in *as* another user to reproduce a bug or
  check their access, with a full audit trail and a one-click return. The swap
  is server-side through HttpOnly cookies, so the admin never handles a token.
- **command-palette**: ⌘K / Ctrl-K navigation and quick actions across the
  admin, built from your resource registry. Frontend-only: it proves a Grit
  plugin can be pure client code, touching no Go at all.
- **saved-views**: save a table's filters, sort, search and date range as a
  named view you can return to, per user, per resource. In other ecosystems this
  is a *paid* add-on, which tells you what it's worth. Here it's a small plugin,
  because Grit's tables already keep their state in the URL, so a saved view is
  just a saved query string.

Each is deliberately contained: real value, a small blast radius, and a different
corner of the mechanism exercised: one touches auth and cookies, one touches no
server at all, one adds a model and a per-user API.

```bash
grit plugin add impersonate
grit plugin add command-palette
grit plugin add saved-views
```

## How to build your own plugin

A Grit plugin is a Go value that describes what to generate. The contract is
small: you declare the **files** to write and the **injections** to make, and
Grit handles recording, ordering, and reversal for you.

```go
package plugin

func init() { Register(webhookPlugin) }

var webhookPlugin = Plugin{
    Name:    "audit-webhook",
    Version: "1.0.0",
    Summary: "POST an event to a webhook on every write.",

    // Files returns path (relative to the project root) -> content. It's a
    // function of the install Context, so a plugin adapts to the project:
    // a triple app gets admin pages an --api project doesn't.
    Files: func(ctx Context) map[string]string {
        return map[string]string{
            ctx.APIRoot + "/internal/hooks/webhook.go": webhookGo(ctx.Module),
        }
    },

    // Injections edit existing files at a named marker. Grit records the exact
    // text in the lockfile so removal takes out precisely what went in.
    Injections: func(ctx Context) []Injection {
        return []Injection{
            {
                File:   "apps/api/internal/routes/routes.go",
                Marker: "// grit:routes:protected",
                Code:   `protected.POST("/webhook-test", hooks.Test)`,
            },
        }
    },

    // Migrations to run, env vars to set. Printed after a successful install.
    NextSteps: []string{"Set WEBHOOK_URL in .env"},
}
```

The rules the installer enforces for you:

1. **No overwrite.** A plugin can't clobber a file that already exists: it
   creates or it fails loudly.
2. **A missing marker is an error.** If an injection's marker isn't found, the
   install stops rather than silently doing nothing.
3. **No double install.** Installing an already-installed plugin is refused.
4. **Requirements both ways.** A plugin's dependencies must be present to
   install it, and can't be removed while it's still installed.

You never write uninstall code. Because every file and injection is recorded in
`.grit/plugins.lock.json`, `grit plugin remove` derives the reversal. Write the
install; get the uninstall for free.

The full authoring guide (every marker you can inject at, the `Context` fields,
optional injections, and how to test a plugin's install/remove round-trip) is in
the docs under **Plugins → Writing a plugin.**

## Where this goes

Filament, in the PHP world, has 900-plus community plugins and never built a
registry: it just published a contract and a directory, and let distribution
happen over the package manager. That's the model: **built-in plugins first to
prove the shape, then open it up.** The target is:

```bash
grit plugin add github.com/you/grit-something
```

Install a community plugin straight from a git URL, get a clean reviewable diff,
and remove it just as cleanly if it's not for you.

Plugins are how Grit stays small at the core and still grows. The batteries every
app needs are built in. The bigger, opinionated pieces (multi-tenancy today,
more soon) are plugins: real code, in your repo, that you can always take back
out.

---

**Read more in the docs:** [What are plugins](/docs/plugins/overview) ·
[Multi-tenancy](/docs/plugins/multitenant) ·
[Impersonate](/docs/plugins/impersonate) ·
[Command palette](/docs/plugins/command-palette) ·
[Saved views](/docs/plugins/saved-views) ·
[Writing a plugin](/docs/plugins/authoring)
