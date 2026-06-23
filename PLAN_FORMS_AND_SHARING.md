# Forms, Tables & Sharing — Implementation Plan

> **Status**: Active. Began 2026-06-21. Targets v3.31.16 → v3.32.0.
>
> Driven by post-v3.31.15 feedback while building a real contact-app.
> Spans 8 features grouped into 4 phases. Each phase is shippable on its
> own — the plan exists so we never end up stuck halfway through a
> half-built abstraction.

## Locked design decisions (from clarification round)

1. **Public form sharing** — Token in URL + optional bcrypt-stored
   password on the share record. Records carry a `source_share_id`
   FK so the admin can filter by share.
2. **Sync auto-form-update** — Marker-based partial regen: append new
   fields to `apps/admin/resources/<plural>.ts` at a known marker, never
   touch existing field config. Removed fields are flagged but not
   deleted.
3. **Steps vs groups** — One mechanism. `form.groups: [...]` defined
   once on the resource. Renders as a stepped wizard inside a
   sheet/modal create flow; renders as side-by-side cards (each with
   its own Save button) on the update page. Each group's Save calls
   PATCH with only that group's fields.

## Phase 1 — Form mechanics + sync drift fix (this session)

Goal: make resources feel like a complete framework, not a one-shot
generator. Highest user-impact items first.

### 1.1 Sync auto-add admin fields (marker-based partial regen)

**Why**: The "I added a Salute field and ran grit sync but it's not in
the admin form" pain. Highest-frequency annoyance.

**Files**:
- `internal/generate/sync.go` — extend `Sync()` with a per-resource
  diff + injection pass
- `internal/generate/sync_resource_admin.go` (new) — owns parsing the
  admin resource file, computing the diff, and writing additions
- `internal/scaffold/admin_resource_files.go` (generator template) —
  emit explicit marker comments around the columns and form-fields
  arrays so we have stable injection points
- Tests in `internal/generate/sync_resource_admin_test.go`

**Marker shape**:
```ts
columns: [
  // grit:cols:auto-start
  { key: "name", ... },
  // grit:cols:auto-end
],
fields: [
  // grit:fields:auto-start
  { key: "name", ... },
  // grit:fields:auto-end
],
```

Injection appends *inside* the marker block when adding fields. Fields
removed from the Go model just print a warning — operator decides
whether to delete them.

### 1.2 Form render mode (sheet | modal | page)

**Config** (new on `FormDefinition`):
```ts
form: {
  render?: "sheet" | "modal" | "page";  // default "sheet"
  // ...rest unchanged
}
```

**Files**:
- `internal/scaffold/admin_resource_files.go` — type def + comment
- `internal/scaffold/admin_form_files.go` (or new file) —
  `<ResourceCreateModal>` component (centered Dialog) and a generated
  `/resources/<plural>/new/page.tsx` route for "page" mode
- `internal/scaffold/admin_tanstack_files.go` — `ResourcePage` switches
  on `form.render`

**Behavior**:
- `sheet` (default) — existing right-anchored sheet on desktop,
  bottom-anchored on mobile
- `modal` — centered dialog with backdrop, fixed max width
- `page` — clicking "New X" navigates to `/resources/<plural>/new`

### 1.3 Form groups (steps in create, cards in update)

**Config** (new on `FormDefinition`):
```ts
form: {
  fields: [...],   // existing flat layout (default)
  groups?: [
    { title: "Basics",  description?: string, fields: FieldDefinition[], scope?: "create" | "update" | "both" },
    { title: "Address", fields: [...] },
  ],
  // groups: present + render: sheet/modal → render as stepped wizard
  // groups: present + render: page (update)  → render as cards w/ per-group Save
}
```

`scope: "create"` = group only shows during Create; `"update"` = only
during Edit; `"both"` (default) = both. This is the **create-and-update**
flow the user described — put title/price in a `scope: "create"` group
and the rest in `scope: "update"` groups.

**Files**:
- `internal/scaffold/admin_form_files.go` — `<FormBuilder>` learns
  about `groups` + wizard navigation (Next/Back/Save)
- `internal/scaffold/admin_resource_files.go` — types
- `internal/scaffold/admin_tanstack_files.go` — update page renders
  per-group cards when `groups` are defined

**Backend support — PATCH endpoint**:
- `internal/scaffold/api_files.go` — generated handler gets a `Patch`
  method that binds a `map[string]interface{}`, validates the keys
  against the model's allowed fields, and calls `db.Model(&x).Updates(m)`
- Routes get `PATCH /api/<plural>/:id` registered alongside PUT
- React-side: each update-page group's Save calls
  `apiClient.patch('/api/contacts/' + id, groupOnlyFields)`

### 1.4 Column-pack default — pack heuristically at generation time

**Why**: Make the generator emit a tidier table out of the box. Existing
resources stay one-field-per-column unless the operator runs the
upgrade command.

**Heuristic**:
- If model has both `name` + `email`: pack into a "Contact" column with
  stacked render (`cell: (row) => <NameEmail row={row} />`)
- If model has a money-shaped float + a related currency string: pack
- Otherwise: one column per field as today

**Implementation**:
- `internal/generate/templates.go` (where columns get built) — add a
  `packGroups(fields)` step that detects pack candidates
- A new helper component `<StackedCell>` scaffolded into
  `apps/admin/components/tables/stacked-cell.tsx` so the generated
  `cell:` callbacks reference one helper, not raw JSX everywhere
- `grit pack table <Resource>` command — re-runs the pack heuristic
  against an existing `resources/<plural>.ts` (idempotent — skips
  already-packed columns)

## Phase 2 — Public form sharing (next session)

### 2.1 FormShare model + admin CRUD

- Model: `{ id, resource_name, token (unique, 32-char random),
  password_hash? (bcrypt), enabled bool default true, submission_count,
  created_by_user_id, created_at, updated_at }`
- Generated as part of the base scaffold (not per-resource): every
  Grit project ships with a `FormShare` table and a `/resources/shares`
  admin page.

### 2.2 Public API surface

- `GET /api/public/forms/:token` — returns `{ resource_name,
  has_password, schema }` so the public page can render the form
  without auth.
- `POST /api/public/forms/:token/check-password` — accepts `{ password }`,
  returns a short-lived JWT in a cookie that the submit endpoint accepts
- `POST /api/public/forms/:token/submit` — validates the JWT
  (when password-protected) and the body (via the shared Zod schema),
  calls the resource service's `Create`, records a `FormSubmission`
  audit row, returns 201.

### 2.3 Submission audit trail

- Add `source_share_id` nullable column to scaffolded models. Set by the
  public-submit handler.
- Admin resource list pages get a filter "Submitted via:
  Admin / Public link X / All". One generic filter, server-side joined
  to `form_shares`.

### 2.4 Public web page

- `apps/web/app/forms/[token]/page.tsx` — renders the form using the
  shared Zod schema. Reads `GET /api/public/forms/:token` for the
  shape, posts to `/submit`.

## Phase 3 — Expose form/table to web (next session)

### 3.1 `grit expose form <Resource>`

```bash
grit expose form Contact --to apps/web/app/contact-us/page.tsx
grit expose form Contact --to apps/web/app/contact-us/page.tsx --public-share
```

- Generates a Next.js page that imports the shared Zod schema and
  the generated `useCreateContact` hook.
- `--public-share` flag generates a page that creates a FormShare and
  posts to the public endpoint instead (needs Phase 2).

### 3.2 `grit expose table <Resource>`

```bash
grit expose table Contact --to apps/web/app/contacts/page.tsx
```

- Generates a Next.js page with a web-friendly table (pagination,
  search, basic styling — not the heavy admin DataTable).
- Honors the admin's `cell` packed columns.
- `--public` flag generates a page that calls a `GET /api/public/contacts`
  endpoint instead (operator-defined route, optional).

## Phase 4 — Web auth scaffolding (next session)

The web app already ships with login/sign-up/forgot-password pages.
What's missing is a way to **protect web pages** the way the admin
auto-protects its dashboard routes.

### 4.1 `grit add web-auth`

If web auth helpers are missing:
- Scaffold `apps/web/middleware.ts` with a config-driven route matcher.
- Scaffold `apps/web/components/ProtectedWebRoute.tsx` — a client
  wrapper component for in-page guarding.
- Scaffold `apps/web/hooks/use-me.ts` (web-side; thinner than admin).

### 4.2 Documentation

- New lesson under chapter 4 or 5: "Protecting web pages".

---

## Open questions deferred

- For Phase 1.3 (groups + PATCH): how does `Version` (optimistic
  concurrency) interact with partial updates? Current plan: PATCH still
  increments Version. Add `If-Match: <version>` semantics later if
  needed.
- For Phase 2.4 (public web page): styling. Use the admin's
  FormBuilder components or a stripped web-styled form? **Decision
  needed before Phase 2 starts.**
- For Phase 4: do we expose web auth via cookies (current admin
  pattern) or a different scheme? Keep cookies for consistency.

---

## Sequence + commit cadence

Phase 1 ships in **one tagged release** (v3.31.16). Phases 2/3/4 each
get their own minor (v3.31.17, .18, .19) or roll up into v3.32.0 once
the migration story is settled.

Each phase is self-contained: a learner who upgrades mid-phase
shouldn't see broken state.
