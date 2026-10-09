package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteAgentsDoc writes CLAUDE.md and AGENTS.md to the project root.
// Both files have the same content — different LLM tooling looks for
// different filenames, and shipping both means the conventions are
// found regardless of which assistant the developer is using.
//
// If a file already exists, it's skipped unless force is true. Returns
// the list of paths written for the CLI to print back.
func WriteAgentsDoc(root string, force bool) ([]string, error) {
	content := agentsDocContent()
	written := []string{}
	for _, name := range []string{"CLAUDE.md", "AGENTS.md"} {
		path := filepath.Join(root, name)
		if _, err := os.Stat(path); err == nil && !force {
			continue
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return written, fmt.Errorf("writing %s: %w", path, err)
		}
		written = append(written, name)
	}
	return written, nil
}

// agentsDocContent returns the canonical Grit conventions doc. Update
// this when framework defaults change so projects regenerating the
// file get the latest rules.
func agentsDocContent() string {
	return `# Grit conventions

> One file. The hard rules every contributor (human or AI) needs before
> their first PR. Re-run ` + "`grit init --force`" + ` to refresh after a major
> framework upgrade.

---

## Forms

- Always use ` + "`<CurrencyField>`" + ` for monetary inputs, auto-formats commas, accepts a currency prefix, emits raw ` + "`number`" + ` to onChange.
- Always use ` + "`<SearchableSelect>`" + ` for FK fields and any enum with > 5 options. Native ` + "`<select>`" + ` for booleans / 2–4 options is fine.
- Always use ` + "`<DateField>`" + ` for dates and ` + "`<DateRangeFilter>`" + ` for filter bars. Don't roll your own date picker.
- Mount create/edit forms inside ` + "`<Drawer>`" + ` (right-edge slide-in). Modals are reserved for confirmations; full pages for top-level workflows like checkout.
- Compose drawer-mounted forms with ` + "`<FormGrid>`" + ` + ` + "`<FormSection>`" + ` + ` + "`<FormActions>`" + `. The ` + "`isPending`" + ` prop on FormActions wires up the loading state.
- Status pills: ` + "`<StatusBadge status={value} />`" + `. Extend the colour map per app via ` + "`setStatusVariants({ shipped: \"info\" })`" + ` at boot.

## Frontend stdlib

- Format helpers live in ` + "`@/lib/format`" + `: ` + "`formatCurrency`" + `, ` + "`formatDate`" + `, ` + "`formatDateTime`" + `, ` + "`humanize`" + `, ` + "`initials`" + `. Configure locale + default currency once at boot via ` + "`setFormatConfig({ locale, currency })`" + `, never instantiate ` + "`Intl.NumberFormat`" + ` inline.
- Error toasts: ` + "`toast.error(apiErrorMessage(err))`" + `. The helper walks the standard envelope chain so you never see ` + "`AxiosError: Request failed with status code 422`" + `.
- Branch on error codes via ` + "`apiErrorCode(err)`" + `; surface per-field validation errors via ` + "`apiErrorFields(err)`" + `.

## Data

- All list endpoints return the standard envelope: ` + "`{ data, meta: { total, page, page_size, pages } }`" + `. Use ` + "`paginate.List[T]`" + ` from the API side, don't hand-roll page math.
- Search is a portable case-insensitive LIKE across the columns declared in ` + "`Config.Searchable`" + `. Only text-like fields are searchable by default, never include FK UUID columns.
- Sort whitelist via ` + "`Config.Sortable`" + `. Out-of-whitelist values fall back to ` + "`created_at desc`" + `.

## Backend

- Errors: ` + "`respond.NotFound(c, msg)`" + ` / ` + "`respond.Validation(c, msg, fields)`" + ` / ` + "`respond.Forbidden(c, msg)`" + ` / ` + "`respond.Conflict(c, msg)`" + ` / ` + "`respond.Internal(c, err)`" + `. Never write ` + "`c.JSON(500, gin.H{\"error\": err.Error()})`" + ` inline.
- Success: ` + "`respond.OK(c, data, msg?)`" + ` / ` + "`respond.Created(c, data, msg?)`" + ` for the standard ` + "`{ data, message? }`" + ` shape.
- Activity log fires automatically on every successful authenticated mutation. Don't hand-roll audit writes.
- Realtime: ` + "`hub.SendToUser(userID, event)`" + ` / ` + "`hub.Broadcast(event)`" + ` to push events to connected WebSocket clients. Topic naming: ` + "`<resource>.<verb>`" + ` (e.g. ` + "`building.created`" + `).

## Resources

- Generate full-stack resources with ` + "`grit generate resource Building --fields \"name:string,description:text,owner_id:belongs_to:User\"`" + `. Don't hand-write models, the generator emits the model, service, handler, Zod schema, TS types, and admin page in a consistent shape.
- Field modifiers: ` + "`:optional`" + ` makes a field optional in the request schema. Default is required for string fields.
- Belongs-to: ` + "`owner_id:belongs_to:User`" + ` produces both the FK column and the association struct. Returns ` + "`*string`" + ` UUIDs, not numeric IDs.
- The auto-emitted ` + "`Export(c)`" + ` handler streams CSV (default) or XLSX at ` + "`GET /api/<plural>/export?format=csv|xlsx`" + `, re-uses the searchable column set.
- Entering many rows at once is already there: every admin list page has **Bulk Create** and **Bulk Edit**, both spreadsheet-shaped grids over the same resource definition. Do not build an import screen or a repeating form. Hide the button with ` + "`table: { bulkCreate: false }`" + `.

## Money

- A ` + "`price:money`" + ` field is ` + "`money.Money`" + `: ` + "`{ Amount int64, Currency string }`" + ` with the amount in MINOR units. 2800 USD is $28.00. Never store a price in a float, and never compare or total prices by converting to one.
- Totals: ` + "`a.Add(b)`" + ` returns ` + "`(Money, error)`" + ` and errors on mixed currencies; ` + "`m.MulInt(n)`" + ` for a line total. Check the error even when the app has one currency, or the day it has two it is silently wrong.
- GORM stores it as two columns with the field's prefix: ` + "`price_amount`" + ` and ` + "`price_currency`" + `. So the sortable and filterable column is ` + "`price_amount`" + `, not ` + "`price`" + `: ` + "`?sort_by=price_amount&sort_order=asc`" + ` and ` + "`?price_amount_min=5000`" + `.
- On the frontend use ` + "`formatMoney(m)`" + ` from ` + "`@repo/shared/types`" + `. It knows that UGX and JPY have no minor unit at all, which a hardcoded ` + "`amount / 100`" + ` does not: that shows a 50,000 shilling price as 500 and nobody notices until the first Ugandan customer complains.

## The public API

- ` + "`grit generate resource X --public`" + ` adds a second, narrower surface: ` + "`GET /api/v1/public/<plural>`" + ` and ` + "`/:key`" + `, guarded by a publishable API key, built from an allowlist in ` + "`internal/handlers/<x>_public.go`" + `. It is READ-ONLY by design, because the audience is the internet.
- That file is NOT overwritten when the resource is regenerated. A column you want published goes in the struct and in ` + "`toPublicX`" + ` by hand; relations and anything that looks like cost, margin, stock or personal data are held back on purpose.
- ` + "`/:key`" + ` resolves by the resource's slug field, whatever you named it (` + "`handle`" + ` is the usual name in a shop), falling back to nothing: a resource with no slug is looked up by id.
- A public endpoint that WRITES (a cart, a newsletter sign-up, a review) is hand-written, and must not be mounted on the ` + "`publicAPI`" + ` group: that group has a response cache keyed on the URL and nothing else, so one caller's answer is served to every caller. Mount it on its own group with the same ` + "`RequireAPIKey`" + `, and set ` + "`Cache-Control: private`" + ` on the response.
- Never take a price, a total or a discount from a request body. Take an id and a quantity, and read every figure from the database.

## Variants

- ` + "`grit add variants --resource Product`" + ` installs options, option values, per-product offered options, variants and the join between them, plus a matrix editor in the admin and one public endpoint.
- ` + "`GET /api/v1/public/<plural>/:key/variants`" + ` returns everything a picker needs in ONE response: the options to draw, every combination with its resolved price and stock, and the price range. Do not fetch options, then variants, then a price per click.
- The price in that payload is the server's. Do not re-apply option deltas in the browser: a second implementation of that arithmetic is a second answer to "what does this cost".
- A product with no variants gets empty lists and a range of its own price, so one component renders both cases. Most of a real catalogue has no options at all.

## Sync (offline-first desktop apps)

- Local writes go through Wails-bound ` + "`LocalCreate`" + ` / ` + "`LocalUpdate`" + ` / ` + "`LocalDelete`" + `, they hit the local SQLite mirror + outbox, never HTTP directly.
- Reads via ` + "`LocalGet`" + ` / ` + "`LocalList`" + ` come from the local mirror, kept fresh by the pull phase of ` + "`Sync()`" + `.
- Push triggers manually via the title-bar Sync button. Conflicts surface in ` + "`<ConflictDialog>`" + ` for per-field merge, defaults to local-wins.
- Every server model carries a ` + "`Version int`" + ` column with a ` + "`BeforeUpdate`" + ` hook that auto-increments. The sync engine uses this for optimistic-lock checks.

## Auth

- Gate routes with the auth middleware on the ` + "`protected`" + ` group. Role-restrict via ` + "`middleware.RequireRole(\"ADMIN\")`" + `.
- Refresh-token interceptor on the client skips ` + "`/auth/login`" + ` / ` + "`/auth/register`" + ` / ` + "`/auth/refresh`" + ` so a wrong-password 401 doesn't loop into a session wipe.
- Idempotency: client auto-attaches ` + "`Idempotency-Key`" + ` on POST/PUT/PATCH/DELETE. Server caches the first 2xx response for 24h so retries are safe.

## The ones that cost an afternoon

Every rule here was learned by something breaking, and each has the same shape:
it fails quietly, or it fails somewhere other than where the mistake was.

### An icon in the admin map is not an export

Adding an icon to the admin ` + "`iconMap`" + ` record does not make it a named export from
` + "`@/lib/icons`" + `. A new page that imports it by name fails to build, and the error
names the import rather than the map. Add it to both.

### A change to login is six changes

The admin ships six auth styles, and the login and sign-up screens exist once
per style. Changing one leaves five behind, and nobody notices until a customer
uses a different style. Search for the screen, not for the component, and expect
six hits.

The shared ` + "`RegisterSchema`" + ` is camelCase: ` + "`firstName`" + `, not ` + "`first_name`" + `. The Go side is
snake_case and the boundary is the schema, so a form wired to the Go field name
type-checks and sends nothing.

### Never call a service from a model

Auto-numbering belongs in ` + "`BeforeCreate`" + `, and it has to call ` + "`sequence.Next`" + `
directly. Reaching for a helper in ` + "`services`" + ` makes models depend on services,
which already depends on models, and Go refuses the import cycle at a file that
had nothing to do with it.

### FindInBatches hands you a transaction with no query

Inside the callback, the ` + "`tx`" + ` has no query attached: read the destination slice
instead. A second ` + "`Order`" + ` inside the callback competes with the one driving the
batching and silently reorders the work, and dropping the callback error loses
every failure in the batch.

### GORM turns initialisms into something you would not guess

` + "`AllowIDPInitiated`" + ` becomes the column ` + "`allow_id_p_initiated`" + `. Pin any unusual
initialism with an explicit ` + "`column:`" + ` tag, and key map updates by the real column
name rather than by what the field looks like it should be.

### The web app and the admin have different palettes

` + "`apps/web`" + ` has ` + "`bg-bg-secondary`" + `, ` + "`text-text-secondary`" + ` and ` + "`text-accent`" + `. The admin has
` + "`muted`" + `, ` + "`card`" + ` and ` + "`primary`" + `. They are not interchangeable, and Tailwind renders an
unknown utility as nothing at all: the element is there, the colour is missing,
and no tool reports it. Copying a component between the two apps means
translating its classes.

### A relative path in .env means two databases

` + "`grit migrate`" + ` runs from ` + "`apps/api`" + `. A binary you built runs from wherever you
launched it. A relative SQLite path therefore resolves to two different files,
and the symptom is a migration that ran against a database the application
cannot see. Use an absolute path, or accept that only the CLI sees your data.

### Changing APP_PORT without APP_URL breaks uploads

` + "`APP_URL`" + ` is what the API hands clients for a stored file. Move the port and
leave the URL behind and every image loads from a port nothing is listening on.
They are one setting in two places.

### The refresh cookie is scoped to /api/auth

` + "`grit_refresh`" + ` is path-scoped, so a handler outside ` + "`/api/auth`" + ` cannot see it and
cannot identify the current session. Anything that needs the session identity
belongs under that prefix or needs the session passed to it.

### A --public resource is read-only

` + "`--public`" + ` generates a read surface. A public endpoint that writes has to be
hand-written, and it must not go on the URL-cached public group: the cache will
serve a stale answer to the next caller, and a write that appears to succeed
and does nothing is the worst available outcome.

### Hash inputs have to survive the column

A value hashed before it is stored and verified after it is read has to be
normalised to what the column keeps. Postgres truncates timestamps to
microseconds and MySQL to milliseconds, so a chain that verifies on SQLite can
fail on Postgres with no code change. Verify hash chains against the database
you deploy to.

### One React, when Expo is in the monorepo

With ` + "`apps/expo`" + ` present, every web app pins Expo's React. Two copies of React in
a hoisted ` + "`node_modules`" + ` break every component test at once, with an error that
points at a component rather than at the duplication.

### pnpm ignores .npmrc when a workspace file exists

` + "`node-linker`" + ` in ` + "`.npmrc`" + ` is ignored by pnpm 10 when there is a
` + "`pnpm-workspace.yaml`" + `, and the setting has to be in the workspace file instead.
Without ` + "`hoisted`" + `, Next cannot resolve ` + "`@swc/helpers`" + ` and every page throws
` + "`_interop_require_wildcard`" + `.

### Next is pinned exactly, on purpose

The version is exact rather than a caret range. A caret once adopted a
one-day-old release whose layout router threw on every navigation, and the
project had changed nothing.

---

## What NOT to do

- Don't import GORM directly in handlers, go through services.
- Don't write inline ` + "`<input type=\"number\">`" + ` for money. Use ` + "`<CurrencyField>`" + `.
- Don't compose role checks via if statements in handlers, use ` + "`RequireRole`" + ` middleware or the ` + "`--roles`" + ` generator flag.
- Don't add new auth flows, extend the existing ones via the auth service.
- Don't bypass ` + "`paginate.List`" + ` for list endpoints.
- Don't write raw SQL against a generated join table. The names are the generator's (` + "`variant_option_values.variant_id`" + `, not ` + "`product_variant_id`" + `), and the association is already there: preload it.
- Don't delete a generated page and expect it gone if you are on a version before v3.367.0; ` + "`grit upgrade`" + ` used to write it back, and two route groups resolving to the same path do not build.
- Don't commit ` + "`apps/desktop/frontend/wailsjs/`" + `, Wails regenerates it on every build.

---

*Generated by ` + "`grit init`" + `. Re-run with ` + "`--force`" + ` to refresh after a framework upgrade. The framework version that produced this file: see ` + "`grit version`" + `.*
`
}
