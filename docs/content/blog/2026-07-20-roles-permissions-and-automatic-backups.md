---
title: "Grit ships roles, permissions, and automatic database backups, by default"
subtitle: "Every new Grit app now comes with a real authorization model (roles, granular per-resource permissions, and a permission editor that works on web, mobile and desktop) plus scheduled database backups you can restore with one click. No wiring. It's just there when you run grit new."
series: "The Daily Grit"
edition: 10
date: 2026-07-20
readingTime: "9 min"
author: "Muke JohnBaptist"
tags: [grit, rbac, roles, permissions, authorization, backups, security, go, react]
canonical: "https://gritframework.dev/blog/roles-permissions-and-automatic-backups"
---

Two things every real app needs, and almost nobody wants to build twice:
**who's allowed to do what**, and **a way to get your data back when something
goes wrong.**

They're unglamorous. They're also the two you regret skipping. So Grit now ships
both in the box: the moment you run `grit new`, before you've written a line of
your own code.

Let me show you exactly what you get.

## Roles and permissions, out of the box

Every generated app starts with three roles seeded into the database:

- **ADMIN**: holds the `*` grant: every permission, including any you add
  later. Nothing to maintain.
- **EDITOR**: content and uploads, no user or system administration.
- **USER**: a standard account with no administrative permissions.

These aren't just labels on a string column. Behind them is a real model:

- a `roles` table and a `user_roles` join, so a user can hold real, editable roles
- a **permission catalog**: every capability expressed as a
  `resource.action` key: `users.view`, `blogs.create`, `uploads.delete`,
  `roles.edit`, and so on
- **wildcards** that stay wildcards: grant `products.*` and the role
  automatically gains any action you add to Products later: you never re-tick a
  box because you shipped a new feature

### The permission editor

Open **Roles & permissions** in the admin and you get a proper editor, not a
textarea of comma-separated strings:

- permissions grouped by module → feature, **collapsed by default** so you expand
  one area at a time instead of scrolling a wall of checkboxes
- a CRUD matrix per feature, create / view / edit / delete, with an em dash for
  actions a feature doesn't have, so you never tick something that does nothing
- a live "N / total granted" counter and a filter
- **copy-permissions-from**, to start a new role from an existing one

Built-in roles are safe to edit, their permissions are fully yours to change,
but their names are locked and they can't be deleted, because routes and the
upgrade path resolve them by name. Custom roles you can delete freely… unless
they're still assigned to someone. Try to delete a role a user still holds and
Grit stops you:

> Role is assigned to 3 user(s); reassign them before deleting it.

That's deliberate. Deleting a role out from under its users would silently strip
their access: exactly the kind of quiet permission bug you don't find until
someone can't do their job.

### It's the same everywhere

The part I'm proudest of: this isn't an admin-only feature. The **same** roles
and permissions model runs on all three clients.

- **Web admin** (Next.js or the Vite/TanStack admin): the full editor above.
- **Mobile** (Expo): a native roles list and permission editor, same endpoints,
  same wildcard rules.
- **Desktop** (Wails): the same, in the native window.

Each client ships a `usePermissions()` hook so your own screens can gate on the
real answer:

```tsx
const { can, isSuper } = usePermissions();

// exact permission
{can("products.delete") && <BulkDeleteButton />}

// any action on the resource
{can("products.*") && <ProductToolbar />}
```

And it is never *only* a UI convenience. Every route is enforced server-side in
Go, hiding a button is courtesy, the API rejects the request whether or not the
button rendered. The middleware is dual-mode, so a route can require a role name
or a permission key:

```go
// either works: a role name, or a fine-grained permission
admin.DELETE("/users/:id", middleware.RequireRole("ADMIN"), userHandler.Delete)
posts.POST("",           middleware.RequireRole("perm:blogs.create"), h.Create)
```

Assigning a role to a user is wired end to end too: pick a role when you create
or edit a user (the dropdown lists **every** role in your database, including
ones you just made), and the user immediately resolves to that role's grants.

## Automatic database backups

The second half: your data is worth nothing if you can't get it back.

Every Grit app ships with a backup system, no add-on, no cron file to write:

- **Backup now**: one click in the admin creates a full database dump, streamed
  straight to your object storage (MinIO in dev, S3/R2/B2 in production).
- **A schedule**: daily, weekly, monthly or yearly, at a time you choose. The
  scheduler picks it up on its next tick; no restart, no redeploy.
- **Download**: backups are pulled with a short-lived pre-signed URL, so the
  archive comes straight from storage and never proxies through your API.
- **Restore**: bring a backup back when you need it.

The backup dashboard lives under **System Hub → Data & Backup**. Set a daily
schedule once and forget it, the same way the rest of Grit's batteries work:
present, wired, and out of your way until you need them.

## Why "by default" matters

You could add all of this yourself. You've probably added a version of it, more
than once. That's the point.

The first day of a project is the day you have the least time for roles tables,
permission catalogs, and backup schedulers, and the most temptation to leave
them for "later." Grit's bet is that the batteries every serious app ends up
needing should be there on day one, already correct, so "later" never turns into
a 2 a.m. restore you can't perform.

Run `grit new`, and roles, permissions, and backups are already in the box.

```bash
grit new my-app --triple
# roles seeded, permission editor live, backup schedule ready
```

Ship the product. The boring, important parts are handled.

---

**Read more in the docs:** [Roles &amp; Permissions](/docs/security/authorization) ·
[Roles &amp; Permissions UI](/docs/admin/roles) ·
[Data &amp; Backup](/docs/batteries/backups)
