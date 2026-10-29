# Docs audit ledger

What has been checked against the software, when, and what was found. Run by
`.claude/skills/docs-audit`.

A page listed here as **clean** was read claim by claim and each claim settled
against the thing that settles it. A page not listed here has not been audited,
whatever it looks like.

192 pages under `docs/app/docs`, plus courses and tutorials. This is the record
of which ones are done.

---

## 2026-10-29

**Mechanical sweep, whole site (284 files).** Three checkers, all of them now
clean except where noted.

| Check | Before | After |
|---|---|---|
| Documented `grit` commands that do not exist | 0 | 0 |
| Documented flags that do not exist | 1 | 0 |
| Dead internal links | 5 | 1 |
| Metadata titles the layout double-suffixes | 6 | 0 |
| Removed features still documented | 0 | 0 |
| Em dashes | 3461 | 3460 |

The one dead link left is `/login` in
`courses/grit-mobile/auth-navigation/page.tsx`, inside a JSX sample of the
reader's own app rather than a link on the page. Correct as written; the
checker cannot see that it is a sample because it is not in a CodeBlock.

**Em dashes are the one large thing still outstanding: 3460 of them across 208
files.** Banned in this project's user-facing writing. Not swept here because
it is not a character substitution: an em dash joins two clauses that each
want different punctuation, so every one is a sentence rewritten by hand. The
changelog alone holds 677.

### Pages read claim by claim

| Page | Verdict |
|---|---|
| `frontend/ui-components` | **was badly stale, now clean** |
| `concepts/architecture-modes/double` | **one false claim, fixed** |
| `plugins/stripe` | clean; steps verified against `grit plugin add stripe` |

#### `frontend/ui-components`

Every checkable claim on the page was wrong, and every copy-pasteable block
failed:

- "a registry of 100 ready-made React components": there are **141**.
- five categories with counts (marketing 20, saas 30, ecommerce 20, layout 20,
  auth 10). There are **three**: marketing 93, ecommerce 34, application-ui 14.
  `saas`, `layout` and `auth` do not exist, so `grit ui list --category saas`
  answers "no components in category".
- `billing-card-01`, `login-card-01`, `otp-input-01`, `hero-split-01`: none of
  these exist, and neither does the `-01` convention they imply. Names are
  `<category>-<group>-<name>`.

Counted from `grit ui list`; the replacement names were checked against the
live registry, which answers 200 for each.

**This also found a bug in the binary.** `grit ui list --category`'s help text
named the same five categories. Fixed in `cmd/grit/ui_cmd.go` by not naming
any: the flag points at `grit ui list`, which prints the real ones, so the same
copy cannot go stale again.

#### `concepts/architecture-modes/double`

Claimed `grit upgrade --triple` scaffolds `apps/admin/` into an existing
project and updates `grit.json`. No such flag: `--triple` is a `grit new`
shorthand, and `upgrade` reads the architecture off the directory layout
(`detectArchitecture`), so it cannot be told to change it. `grit.json` is not
consulted. Replaced with what is actually true, including the part that was
missing: the admin route group inside `apps/web/` does not transfer, which is
worth knowing before choosing double rather than after.

### Checked and fine

Recorded so nobody pays to investigate them twice.

- **16 pages have no entry in `config/docs-metadata.ts`**, so the search
  indexer titles them from their path and says so on every run. The titles it
  derives are good ones: `/docs/plugins/stripe` becomes "Stripe",
  `/docs/blueprints` becomes "Blueprints". Not worth curating entries for.
- **Both indexers agree** after this round of edits: 245 pages in the website
  search index, 238 in the CLI's offline index, all 238 with real prose. The
  two counts differ by scope, not by drift.

---

## 2026-10-29, second round

Countable claims, which is where this kind of page rots first. A number in
prose has an expiry date and nothing tells you when it passed.

### The FormBuilder takes 31 field types, not 8

`prerequisites/nextjs` said **8**, and named them: text, textarea, number,
select, date, toggle, checkbox, radio. `concepts/architecture-modes/triple`
said **10+**, which is true the way "2+" is true.

The real number is **31**, counted from the `FieldType` union the scaffold
writes into every project's `lib/resource.ts`, and confirmed in a generated
project rather than from the template.

This is not a rounding error in a sentence. A reader deciding whether the
admin can render their form counts on that number, and eight is a reason to
go and write the form by hand.

`internal/docscheck/field_types_test.go` compares every "N field types" claim
in the docs to that union, so the count cannot drift again.

### A plugin nobody could read about

`grit plugin add webhooks` has worked for months with no page under
`/docs/plugins`. The other nine plugins each have one. It was findable only
by running `grit plugin list`, which nobody does before deciding whether a
thing exists.

Written: `/docs/plugins/webhooks`, covering the Standard Webhooks signature
and why it signs `{id}.{timestamp}.{body}` rather than the body alone, the
retry and dead-letter path, and the delivery log.

`internal/docscheck/plugin_pages_test.go` compares `plugin.All()` to the
directories under `docs/app/docs/plugins`, so the next plugin cannot ship
undocumented.

### Still clean

Commands, flags, colour tokens, toolchain versions, dead links, double
titles, removed features. The one dead link left is the `/login` inside a
mobile course's code sample, which is the reader's own app.

---

## 2026-10-29, third round: paths

`check_paths.py` reads every project path the docs name and looks for it in
a scaffolded project. Not in the templates: the templates are what Grit
means to write and a project is what it wrote, and the gap between them is
what this is for.

A reader who opens `apps/api/internal/handler/` and finds nothing does not
conclude the page is wrong. They conclude they scaffolded it incorrectly.

### Fixed

| The docs said | It is actually |
|---|---|
| `internal/handler/` | `internal/handlers/` |
| `internal/service/auth_service.go` | `internal/services/auth.go` |
| `internal/worker/image_worker.go` | `internal/media/image.go` |
| `internal/handler/upload_handler.go` | `internal/handlers/upload.go` |
| `internal/config/database.go` | `internal/config/config.go` |
| `apps/web/lib/api-client.ts` | `apps/web/lib/api.ts` |
| `apps/web/lib/query-client.ts` | `apps/web/components/providers.tsx` |
| `apps/admin/components/data-table.tsx` | `apps/admin/components/tables/data-table.tsx` |

And one that was not in the wrong place but does not exist:
**`apps/web/tailwind.config.ts`**. Scaffolded projects are on Tailwind 4,
which is configured in CSS. The page showed a Tailwind 3 default export with
content globs and `require("tailwindcss-animate")`, so a reader following it
created a file Tailwind 4 ignores. Replaced with the real `@theme` block out
of `globals.css`.

### The checker had to be taught three things first

Its first run reported 143 misses and nearly all were its own fault, which
is worth recording because an unreliable checker gets switched off:

- **Generated files.** `models/post.go` is a fine example on a page about
  generating a Post; it is simply absent from a project where nobody
  generated one. Those directories are excluded.
- **Ordered alternation.** The extension pattern listed `ts` before `tsx`,
  so every `.tsx` path matched as `.ts` and came back missing. Longest
  extension first.
- **Apps the reference lacks.** A path under `apps/expo` is not wrong
  because the reference project has no Expo app.

### Left alone, deliberately

29 paths remain and most are tutorials telling the reader to **create** a
file: `components/widgets/revenue-stats.tsx` in the ecommerce tutorial,
`jobs/order_confirmation.go`, the chat components. Rewriting those to point
at existing framework files would be wrong.

Telling those apart from a stale claim needs the surrounding prose read, one
page at a time. That is the next pass, and the list above is where to start.

---

## Not yet audited

Everything else. The sections most likely to have rotted, in order, are the
ones whose subject was replaced or renamed, and the ones with counts in them.
`tech-kits`, `batteries` and `concepts` are the obvious places to go next:
they are old, they are long, and they are full of numbers.
