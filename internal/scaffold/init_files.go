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

- Always use ` + "`<CurrencyField>`" + ` for monetary inputs — auto-formats commas, accepts a currency prefix, emits raw ` + "`number`" + ` to onChange.
- Always use ` + "`<SearchableSelect>`" + ` for FK fields and any enum with > 5 options. Native ` + "`<select>`" + ` for booleans / 2–4 options is fine.
- Always use ` + "`<DateField>`" + ` for dates and ` + "`<DateRangeFilter>`" + ` for filter bars. Don't roll your own date picker.
- Mount create/edit forms inside ` + "`<Drawer>`" + ` (right-edge slide-in). Modals are reserved for confirmations; full pages for top-level workflows like checkout.
- Compose drawer-mounted forms with ` + "`<FormGrid>`" + ` + ` + "`<FormSection>`" + ` + ` + "`<FormActions>`" + `. The ` + "`isPending`" + ` prop on FormActions wires up the loading state.
- Status pills: ` + "`<StatusBadge status={value} />`" + `. Extend the colour map per app via ` + "`setStatusVariants({ shipped: \"info\" })`" + ` at boot.

## Frontend stdlib

- Format helpers live in ` + "`@/lib/format`" + `: ` + "`formatCurrency`" + `, ` + "`formatDate`" + `, ` + "`formatDateTime`" + `, ` + "`humanize`" + `, ` + "`initials`" + `. Configure locale + default currency once at boot via ` + "`setFormatConfig({ locale, currency })`" + ` — never instantiate ` + "`Intl.NumberFormat`" + ` inline.
- Error toasts: ` + "`toast.error(apiErrorMessage(err))`" + `. The helper walks the standard envelope chain so you never see ` + "`AxiosError: Request failed with status code 422`" + `.
- Branch on error codes via ` + "`apiErrorCode(err)`" + `; surface per-field validation errors via ` + "`apiErrorFields(err)`" + `.

## Data

- All list endpoints return the standard envelope: ` + "`{ data, meta: { total, page, page_size, pages } }`" + `. Use ` + "`paginate.List[T]`" + ` from the API side — don't hand-roll page math.
- Search is ILIKE across the columns declared in ` + "`Config.Searchable`" + `. Only text-like fields are searchable by default — never include FK UUID columns.
- Sort whitelist via ` + "`Config.Sortable`" + `. Out-of-whitelist values fall back to ` + "`created_at desc`" + `.

## Backend

- Errors: ` + "`respond.NotFound(c, msg)`" + ` / ` + "`respond.Validation(c, msg, fields)`" + ` / ` + "`respond.Forbidden(c, msg)`" + ` / ` + "`respond.Conflict(c, msg)`" + ` / ` + "`respond.Internal(c, err)`" + `. Never write ` + "`c.JSON(500, gin.H{\"error\": err.Error()})`" + ` inline.
- Success: ` + "`respond.OK(c, data, msg?)`" + ` / ` + "`respond.Created(c, data, msg?)`" + ` for the standard ` + "`{ data, message? }`" + ` shape.
- Activity log fires automatically on every successful authenticated mutation. Don't hand-roll audit writes.
- Realtime: ` + "`hub.SendToUser(userID, event)`" + ` / ` + "`hub.Broadcast(event)`" + ` to push events to connected WebSocket clients. Topic naming: ` + "`<resource>.<verb>`" + ` (e.g. ` + "`building.created`" + `).

## Resources

- Generate full-stack resources with ` + "`grit generate resource Building --fields \"name:string,description:text,owner_id:belongs_to:User\"`" + `. Don't hand-write models — the generator emits the model, service, handler, Zod schema, TS types, and admin page in a consistent shape.
- Field modifiers: ` + "`:optional`" + ` makes a field optional in the request schema. Default is required for string fields.
- Belongs-to: ` + "`owner_id:belongs_to:User`" + ` produces both the FK column and the association struct. Returns ` + "`*string`" + ` UUIDs, not numeric IDs.
- The auto-emitted ` + "`Export(c)`" + ` handler streams CSV (default) or XLSX at ` + "`GET /api/<plural>/export?format=csv|xlsx`" + ` — re-uses the searchable column set.

## Sync (offline-first desktop apps)

- Local writes go through Wails-bound ` + "`LocalCreate`" + ` / ` + "`LocalUpdate`" + ` / ` + "`LocalDelete`" + ` — they hit the local SQLite mirror + outbox, never HTTP directly.
- Reads via ` + "`LocalGet`" + ` / ` + "`LocalList`" + ` come from the local mirror, kept fresh by the pull phase of ` + "`Sync()`" + `.
- Push triggers manually via the title-bar Sync button. Conflicts surface in ` + "`<ConflictDialog>`" + ` for per-field merge — defaults to local-wins.
- Every server model carries a ` + "`Version int`" + ` column with a ` + "`BeforeUpdate`" + ` hook that auto-increments. The sync engine uses this for optimistic-lock checks.

## Auth

- Gate routes with the auth middleware on the ` + "`protected`" + ` group. Role-restrict via ` + "`middleware.RequireRole(\"ADMIN\")`" + `.
- Refresh-token interceptor on the client skips ` + "`/auth/login`" + ` / ` + "`/auth/register`" + ` / ` + "`/auth/refresh`" + ` so a wrong-password 401 doesn't loop into a session wipe.
- Idempotency: client auto-attaches ` + "`Idempotency-Key`" + ` on POST/PUT/PATCH/DELETE. Server caches the first 2xx response for 24h so retries are safe.

## What NOT to do

- Don't import GORM directly in handlers — go through services.
- Don't write inline ` + "`<input type=\"number\">`" + ` for money. Use ` + "`<CurrencyField>`" + `.
- Don't compose role checks via if statements in handlers — use ` + "`RequireRole`" + ` middleware or the ` + "`--roles`" + ` generator flag.
- Don't add new auth flows — extend the existing ones via the auth service.
- Don't bypass ` + "`paginate.List`" + ` for list endpoints.
- Don't commit ` + "`apps/desktop/frontend/wailsjs/`" + ` — Wails regenerates it on every build.

---

*Generated by ` + "`grit init`" + `. Re-run with ` + "`--force`" + ` to refresh after a framework upgrade. The framework version that produced this file: see ` + "`grit version`" + `.*
`
}
