# Plan — Permissions, module toggles, and multi-tenancy

> Response to issue #71 (@rahulcodepython). Reference implementation: **Shoppleet**
> (`D:\PROJECTS\Shoppleet\business-management-platform`), which already ships the
> roles + permissions system we want, and is multi-tenant.
>
> Status: **draft for review.** Nothing here is built yet.

---

## 0. What was verified before planning

| Claim | Finding |
|---|---|
| `grit remove resource Blog` lets you drop the blog | **Broken.** Deletes `models/blog.go`, leaves `blog_handler.go`, `blog_service.go`, `blogs_seeder.go` → `undefined: models.Blog` ×10, build fails |
| Endpoints are guarded by roles, not permissions | True. Only `middleware.RequireRole`; no permission concept anywhere |
| `User.Role` | A single `string` column — a user cannot hold two roles at all |
| Smallest scaffold (`--api`) | 27 packages, 16 models, 108 Go files, ~17k LOC |
| Plugin system | **Does not exist.** `/docs/plugins` lists names; no install command. `grit add` exists but only has `web-auth` |

---

## 1 & 2 — "Blank project" / "how do I build something basic"

**Position (agreed): don't add a blank mode.** Grit targets mature projects; the batteries
are the product. Small/marketing sites are better served by plain Next.js.

**But the position only holds if the escape hatch works**, and today it doesn't.

### Workstream A — fix `grit remove resource` (blocks the answer to #71)

`RemoveResource` deletes the model and reverses marker injections, but never deletes the
sibling files it generated.

- [ ] Delete the full artefact set, mirroring exactly what `generate resource` writes:
      `models/`, `handlers/`, `services/`, `database/*_seeder.go`, plus frontend
      (`hooks/use-*.ts`, `types`, Zod schema block, admin resource definition, admin page,
      TanStack route + page, Expo screens if present).
- [ ] Reverse every injection: AutoMigrate, GORM Studio, routes, resource registry,
      sidebar nav, seeder registry, **and (new) the permission catalog block** (§4).
- [ ] Guard: refuse to remove a resource another one references (belongs_to), with a clear
      error rather than a broken build.
- [ ] **Acceptance test:** `grit new x --api && grit remove resource Blog --force &&
      go build ./...` exits 0. Add to CI so it can't regress.

This is small, unblocks the honest answer to #71, and is worth shipping on its own.

---

## 3 — Module bloat: disable via env

**Agreed approach: ship everything, let `.env` switch modules off.**

### Workstream B — module toggles

- [ ] `internal/config`: a `Modules` struct read from env —
      `MODULE_AI`, `MODULE_JOBS`, `MODULE_CRON`, `MODULE_BACKUP`, `MODULE_WEBHOOKS`,
      `MODULE_REALTIME`, `MODULE_PDF`, `MODULE_EXPORT`, `MODULE_TOTP`, `MODULE_FLAGS`,
      `MODULE_SYNC`, `MODULE_AUDIT`. **Default `true`** — existing apps unchanged.
- [ ] `routes.go`: mount each module's route group only when enabled.
- [ ] Workers/schedulers: asynq workers and cron entries only registered when enabled.
- [ ] AutoMigrate: skip a disabled module's tables (so a no-jobs app has no jobs tables).
- [ ] **Admin nav must respect it** — a disabled module cannot leave a dead sidebar link.
      Expose the flags via an existing endpoint (`/api/system/modules`) and gate nav on it.
- [ ] `.env.example` documents every flag in one block.
- [ ] Docs page: "Turning modules off".

### Honest limitation to state in the reply

Env toggles stop a module *running*; they don't delete its code. Part of #71's complaint is
about **reading** a 17k-LOC codebase. Toggles fix the runtime surface and the mental
surface (nothing in the UI, no tables, no workers) — they do not shrink the repo. If we want
that too, the follow-up is a real module system (`grit add jobs` / `grit remove jobs`),
which is a bigger piece of work and should be its own decision.

---

## 4 — Roles **and** permissions (the enterprise feature)

**Goal:** ship RBAC where routes check *permissions*, roles are named bags of permissions,
every generated resource contributes its permissions automatically, and operators manage it
all from a real UI. Port the design from Shoppleet, minus its known bugs.

### 4.1 How Shoppleet does it (the parts worth copying)

- Permission = a flat dotted string, `module.submodule.feature.action`, action ∈
  `create|view|edit|delete`. **Wildcards** (`sales.*`) supported in stored grants.
- A **catalog** (Go literal) declares the 3-level tree and, per feature, *which actions are
  meaningful* (`allActions` / `viewOnly` / `viewEdit`). The UI renders `—` for actions not
  in that set. 9 modules → 175 leaves.
- `Role` row stores grants as a **JSON array in a `text` column**; unique `(business_id, name)`;
  `IsSystem` marks presets.
- Guard is **dual-mode**: `RequireRole("OWNER","MANAGER","perm:sales.transactions.sales.create")`
  — any arg matching passes. This is the backwards-compatibility bridge.
- Frontend `can(code)` hook gates sidebar items (`requires: "sales.*"`) and buttons.
- Presets: Owner (`*`), Manager, Cashier, Sales Rep, Stock Keeper; clone-from-preset in the
  editor; tri-state checkboxes; live "N / 175 granted" counter; text filter.

### 4.2 Bugs in Shoppleet — fix in the port, don't copy

1. `matchSegmented` returns true at the **first** `*`, so `auditing.user_activity.*.view`
   silently grants **all four actions**, not just view. → implement true segment-wise
   matching with a trailing-`*` rule.
2. `SeedDefaultRolesForBusiness` runs on **every** `GET /api/roles` and overwrites preset
   permissions → operator edits to a preset are **silently reverted** on next page load.
   → seed once (migration/first-boot), never on read; make "reset to factory" an explicit
   button.
3. Backend `Update` doesn't block renaming an `IsSystem` role (UI-only lock). → enforce
   server-side.
4. The wildcard matcher is hand-duplicated in Go and twice in TS, already drifting.
   → **one** implementation; server returns expanded leaves so the client never re-implements
   matching.
5. `AssignUser` must write both `role_id` and the legacy `role` string or `RequireRole`
   breaks. → in Grit, one code path resolves a user's grants; no dual writes.
6. The editor saves **expanded leaves**, so a role stops inheriting permissions added later.
   → **preserve wildcards as authored**: ticking a whole module stores `products.*`, so new
   permissions in that module are inherited. (Biggest functional improvement over Shoppleet.)

### 4.3 Design for Grit

**Key format.** Grit resources are flat, so the 4-segment scheme is awkward for generated
code. Proposal: **`<resource>.<action>`** (`products.create`), with the catalog carrying a
display-only `module` for grouping in the UI. `grit generate resource Product --module=Inventory`
sets the group; default group "Resources". Short, predictable keys; the tree stays for UI.

**Catalog is generated, not hand-maintained.** `internal/authz/catalog.go` with marker
injection (`// grit:perms:auto-start` / `-end`), matching Grit's existing pattern, so
`generate resource` adds a block and `remove resource` removes it.

**Storage.** Follow Shoppleet: JSON array of grants in a `text` column on `roles`. Simple,
one query, no join tables — plus wildcard-preserving save (fix #6) so it stays inheritable.

**Grant resolution + caching.** Shoppleet does up to 3 DB queries *per request*. Grit should
resolve once and cache role→grants in-process, invalidated by a `Version` bump on role
update (Redis when configured). Never put grants in the JWT — role changes must take effect
without re-login.

**One indirection point for "what can this user do".** Everything goes through
`authz.GrantsFor(ctx, user)`. This is the seam the multi-tenant plugin (§5) overrides to
return *per-organization* grants — designing it now avoids a breaking migration later.

### 4.4 Workstream C — build order

- [ ] **C1 Catalog** — types (`PermissionModule/Group/Feature`, `Action`), action presets
      (`allActions`/`viewOnly`/`viewEdit`), matcher + `Expand()`, unit tests incl. the
      `a.*.c.view` case from bug #1.
- [ ] **C2 Model** — `Role` (name, description, grants JSON, is_system, version), user→role
      assignment, `authz.GrantsFor()`, cache + invalidation. Migration seeds default roles
      **once**.
- [ ] **C3 Enforcement** — dual-mode `RequireRole("ADMIN","perm:products.create")`, keeping
      existing `RequireRole("ADMIN")` call sites working untouched. Owner/admin short-circuit.
- [ ] **C4 Codegen** — `generate resource` injects catalog block + emits `perm:` args on the
      routes it writes; `remove resource` reverses it. Default roles get the new perms.
- [ ] **C5 API** — `GET /api/permissions/catalog`, `GET/POST/PUT/DELETE /api/roles`,
      assign-user, `GET /api/auth/me` returns **expanded** grants (so the client never
      re-implements matching).
- [ ] **C6 Admin UI** — roles list (preset pill, user count) + permission tree editor
      (tri-state checkboxes, CRUD matrix, `—` for unsupported actions, N/total counter,
      filter, clone-from-preset). **Built once as a shared page component**, wired into
      Next.js admin, TanStack admin and desktop — the drift pattern that caused #69/v3.62.0
      must not repeat.
- [ ] **C7 Frontend gating** — `usePermissions()` / `can()`, `requires:` on nav items,
      hide-while-loading (no privilege flash).
- [ ] **C8 Docs** — authorization guide: concepts, key format, adding permissions, guarding
      routes, the UI, migrating from role-only. Update the security page + changelog.

**Backwards compatibility:** default roles ship with **all** permissions, so an upgraded app
behaves exactly as before until an operator narrows a role.

---

## 5 — Multi-tenancy as an opt-in plugin

**Position (agreed): not in core.** Delivered as `grit add multitenant`.

Note Shoppleet is **not** the `user × org × role` model #71 asks for: it is one business per
user (`users.business_id`) with a *branch* switcher. Its tenancy plumbing is worth copying;
its role-assignment shape is not.

### Workstream D

- [ ] **D1** `grit add multitenant` command (first real module-add; sets the pattern).
- [ ] **D2** Models: `Organization`, `OrganizationMember(user_id, org_id, role_id)` — this is
      what gives "Editor in org A, Moderator in org B". Override `authz.GrantsFor()` (§4.3)
      to resolve grants for the **active** org.
- [ ] **D3** Active-org resolution middleware (header / subdomain / session) → context.
- [ ] **D4** **Automatic query scoping via a GORM callback/global scope.** Shoppleet
      hand-writes `Where("business_id = ?")` **447 times across 33 files** — the single
      biggest weakness in that codebase. Grit must make scoping structural so it cannot be
      forgotten, with an explicit opt-out for cross-tenant admin queries.
- [ ] **D5** Retrofit `OrgID` onto existing models + migration for apps adopting it later.
- [ ] **D6** UI: org switcher, members + per-org role management, invitations.
- [ ] **D7** Tenant-isolation tests: cross-tenant read/write must fail (IDOR).
- [ ] **D8** Docs.

---

## Sequencing

```
A (fix remove)  ── small, unblocks the #71 reply, ship first
B (module env toggles) ── independent, medium
C (permissions) ── the big one; C1→C2→C3 are the foundation
D (multitenant) ── AFTER C. D2 plugs into the authz.GrantsFor seam from C2.
```

Doing **D before C** would mean building org-scoped role assignment twice. Do C first.

---

## Decisions taken (MUKE said "proceed" on the recommendations)

1. **Key format** — `<resource>.<action>`, e.g. `products.create`. Two segments.
   Module/group in the catalog is **display metadata for the UI only**, never part of
   the key. Grit resources are flat; a 4-segment scheme would be noise for generated code.
2. **Storage** — JSON grants column on `roles`, as Shoppleet — *but wildcards are
   preserved as authored*. Ticking a whole resource stores `products.*`, so the role
   inherits actions added later. (Shoppleet expands to leaves on save, which is why its
   roles silently stop inheriting.)
3. **User ↔ role** — **many-to-many** via a `user_roles` join table. Costs almost nothing
   now; the multi-tenant plugin later just adds `org_id` to that same table instead of
   forcing a breaking migration off a single `role_id` column.
4. **Module toggles** — env-based (Workstream B), as agreed.
5. **Catalog scope** — small core (users, roles, uploads, system/jobs/backups/audit);
   generated resources grow it. Shoppleet's 175 are its own domain, not every Grit app's.

### Progress

- [x] **C1 Catalog + matcher** — `internal/authz/permissions.go` in the scaffold:
      `Action`/`Feature`/`Group`/`Module`, action presets, `Catalog()` (core +
      `generatedModules()` between `grit:perms:auto-*` markers), `Keys()`, `Granted()`,
      `Expand()`, `HasAll()`. Shipped with tests that pin the semantics — including the
      middle-wildcard case Shoppleet gets wrong (`products.*.view` must NOT grant delete).
- [ ] C2 Role model + `user_roles` + `authz.GrantsFor()` + cache
- [ ] C3 Dual-mode guard
- [ ] C4 Codegen into the catalog markers
- [ ] C5 API endpoints
- [ ] C6 Admin UI (shared component across Next/TanStack/desktop)
- [ ] C7 Frontend `can()` + nav gating
- [ ] C8 Docs

## Superseded questions (answered by the decisions above)

1. **Key format** — `products.create` (proposed, short) or Shoppleet's full
   `module.submodule.feature.action`?
2. **Grant storage** — JSON column (Shoppleet, simple) or normalized
   `permissions` + `role_permissions` tables (queryable: "who can delete invoices?")?
3. **One role per user, or several?** Shoppleet and Grit are both single-role today.
   Multi-role in core costs little now and is painful to retrofit.
4. **Module toggles** — is "doesn't run, hidden from UI, no tables" enough for #71, or do we
   also want a real `grit add/remove <module>` that adds/removes the code?
5. **Scope of the port** — do we copy Shoppleet's exact 175-permission business catalog
   (sales/inventory/IMEI…), or ship a small core catalog (auth, users, uploads, system) and
   let generated resources grow it? (I recommend the latter — the former is Shoppleet's
   domain, not every Grit app's.)
