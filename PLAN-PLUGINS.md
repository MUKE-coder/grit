# Plugin system — research and plan

Research into Filament's plugin ecosystem (921 plugins, 433 authors) and what it
means for Grit's own mechanism.

---

## 1. How Filament plugins actually work

**Distribution is Composer.** A plugin is an ordinary PHP package installed with
`composer require vendor/name`. Filament runs no registry of its own — the
plugin *directory* on filamentphp.com is a catalogue that links to packages;
installation goes through the language's existing package manager.

**The contract is tiny.** A panel plugin implements three methods:

```php
interface Plugin {
    public function getId(): string;          // unique id among plugins
    public function register(Panel $panel): void;  // register resources, pages, hooks
    public function boot(Panel $panel): void;      // runs only when the panel is in use
}
```

Enabled by the user in one line:

```php
return $panel->plugin(new BlogPlugin());
// convention: a static make() so it reads BlogPlugin::make()
```

**Two plugin kinds**, which can live in one package:
- *Panel plugins* — whole feature sets (Blog, User Management, dashboards)
- *Standalone plugins* — a single form field or table column, usable outside a panel

**Assets** (CSS/JS/Alpine) register in the service provider's `packageBooted()`.

**Migrations, config and translations are published, not injected**:
`php artisan vendor:publish --tag="advanced-tables-migrations"`, then
`php artisan migrate`. The user opts in to owning a copy.

**There is no uninstall story.** Removal is `composer remove` plus whatever the
user remembers to undo by hand — published migrations, config, traits added to
their User model.

### The structural difference from Grit

Filament plugins are **runtime libraries**. Code stays in `vendor/`, the plugin
object hooks into a running panel, and upgrades arrive via Composer.

Grit is a **code generator**. There is no runtime framework object to hook into —
`grit new` emits a plain Go + React app that has no dependency on Grit at all
after generation. So a Grit plugin cannot be a library that registers itself; it
has to *write code into the project*.

That difference is the whole design, and it cuts both ways:

| | Filament (runtime) | Grit (codegen) |
|---|---|---|
| Upgrades | `composer update` | user owns the code; no upgrade path |
| Customising | extend/override the package | edit the generated code directly |
| Uninstall | delete the dep, undo by hand | **can be exact** — we recorded what we wrote |
| Reading the code | dig through `vendor/` | it's in your repo |

**The uninstall gap is the opening.** Filament's model can't cleanly reverse a
plugin, because publishing a migration is a one-way door. Grit's can, because
install records every file and every injection in `.grit/plugins.lock.json` and
removal replays it backwards. That's a genuine advantage, and it's already built.

**The lesson to steal is distribution.** Filament didn't build a registry — it
leaned on Composer and kept a catalogue. Grit's equivalent is a git URL plus a
directory page on gritframework.dev. No infrastructure to run.

---

## 2. Top Filament plugins vs what Grit already ships

Sorted by what they'd add to Grit today.

| Filament plugin | Capability | Grit status |
|---|---|---|
| **Shield** (most-starred) | Roles + permissions UI over policies | ✅ **shipped** v3.66–3.70 — catalog, wildcard grants, tri-state tree editor |
| **Breezy** | Login, register, reset, verify, profile | ✅ **core** — scaffolded auth + profile |
| **Spatie Media Library** | Media management | ✅ **core** — uploads, S3/R2/MinIO, image processing |
| **Filament Excel** | Export resources to Excel | ✅ **core** — export-menu + xlsx |
| **Advanced Filter** | Clause-based table filters | ✅ **core** — table filters, date filter |
| **Spatie Settings** | Settings UI | ~ partial — dashboard settings exist, no general settings store |
| **Environment Indicator** | Show which env you're in | ~ trivial; not present |
| **Custom Dashboards** | Drag-and-drop dashboard builder | ~ partial — dashboard settings + widgets, not user-composed |
| **Advanced Tables** | **Saved table views**, multi-sort, quick filters | ❌ **gap — strong candidate** |
| **Custom Fields** | User-defined fields, no migration | ❌ gap — big feature, high value for CMS-ish apps |
| **Impersonate** | Admin logs in as another user | ❌ gap — small, high utility, needs care |
| **Spotlight / command palette** | ⌘K navigation | ❌ gap — small, high polish |
| **Quick Create** | Header dropdown to create anything | ❌ gap — small |
| **Flowforge** | Kanban board from a model | ❌ gap — visual, demos well |
| **Spatie Tags** | Tagging | ❌ gap — small |
| **Peek** | Full-screen preview modal | ❌ gap — small |

### What this says

1. **Grit already ships natively what Filament's *most popular* plugins add.**
   Shield, Breezy, Media Library, Excel and Advanced Filter are the top of that
   ecosystem, and all five are core Grit. That's the "batteries included"
   position paying off — and worth saying out loud in marketing.

2. **The gaps are mostly small, self-contained UI features** — impersonation,
   command palette, quick create, tags, kanban. Each is a good first plugin:
   real value, contained blast radius, exercises a different part of the
   mechanism.

3. **Two are large and genuinely differentiated**: saved table views (Advanced
   Tables is a *paid* plugin, which tells you what it's worth) and user-defined
   custom fields.

---

## 3. Recommended plugin roadmap

**First-party, to prove the mechanism** (in order):

1. `multitenant` — orgs + membership + automatic query scoping. Already
   committed to; exercises models, migrations, middleware and admin UI at once.
2. `impersonate` — small, exercises auth + a UI affordance + an audit trail.
3. `command-palette` — frontend-only; proves a plugin that touches no Go.
4. `saved-views` — the Advanced Tables equivalent; the highest-value gap.

**Then open it up.** Once two or three exist and the shape is proven, publish the
authoring guide and a directory page, and accept community plugins by git URL.

---

## 4. Decisions taken

- **A Grit plugin generates code**, it is not a runtime dependency. This follows
  from Grit being a generator; there's nothing to hook into at runtime.
- **Removal is lockfile-driven.** Install records every file and injection in
  `.grit/plugins.lock.json`; removal replays it. No second list for an author to
  keep in sync — the exact drift that broke `grit remove resource`.
- **Distribution: built-in first, git URL next.** Follow Filament in not running
  a registry. A catalogue page on the docs site plus
  `grit plugin add github.com/user/grit-foo` is the target; built-in plugins
  ship in the CLI until the manifest format is settled.
- **No subdomains for multitenancy** (per maintainer). Tenant resolution is
  header/session-based.
