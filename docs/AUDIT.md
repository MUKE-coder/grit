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

## Not yet audited

Everything else. The sections most likely to have rotted, in order, are the
ones whose subject was replaced or renamed, and the ones with counts in them.
`tech-kits`, `batteries` and `concepts` are the obvious places to go next:
they are old, they are long, and they are full of numbers.
