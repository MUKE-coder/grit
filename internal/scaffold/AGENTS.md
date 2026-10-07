# internal/scaffold

## Context

The project writer. 650-odd Go files whose job is to emit another project's
source: Go for the API, TypeScript and TSX for the frontends, YAML, Dockerfiles,
`.env`. Almost every file here is a function returning a **string template** (a
Go raw string), plus maps that say where each string lands.

`grit new` and `grit upgrade` both run through here, which is the single most
important fact about this package.

## Source of truth

- **Where a file goes**: the maps in `api_files.go`, `web_files.go`,
  `admin_files.go`, `admin_tanstack_files.go`, `desktop_client_files.go`.
- **What travels on upgrade**: `api_framework_owned_files.go` and the
  `write*Files(root, opts) error` writers called from both `scaffold.go` and
  `upgrade.go`.
- **Which shape a project is**: the predicates on `Options` in `scaffold.go`.
- **Where the admin panel's files live**: `adminPath` and friends in
  `admin_layout_paths.go`. Three shapes, not two.
- **Where a frontend lives**: `webAppRoot(root, opts)` in
  `admin_embedded_single.go`.

## Boundary rules

1. **Delivery has two lists, and a framework-owned file needs both.** The
   new-project map reaches `grit new`; a writer reaches `grit upgrade`. A file in
   only one reaches half the projects, and the failure is silent: the file is
   simply absent, which in a Vite app is a blank screen rather than a build
   error. **Prefer one `writeXFiles` writer called from both** over two lists.

2. **Never spell a path.** `apps/web` is wrong for a single, whose frontend is
   the project root. `apps/admin` is wrong for a double and a single, whose panel
   is embedded. Use the helpers.

3. **Ask the predicate you mean.** `HasAdminPanel()` for anything that produces a
   screen; `ShouldIncludeAdmin()` only for the separate application's
   `package.json`, port and container.

4. **Build an emitted import block as a list, never with `strings.Replace` on an
   anchor line.** A `Replace` that matches nothing fails silently and ships a
   file missing its import.

5. **No backticks inside a string template.** The template is a Go raw string, so
   a backtick ends it. A struct tag needs `` ` + "`" + `` ; a JS template literal
   is usually better rewritten as concatenation.

6. **Never `gofmt` a string template's contents**, and never run `gofmt` over a
   file to tidy an emitted Go snippet. Compare generated projects **by bytes**,
   with the same project name and both binaries; different names and decoded text
   each give a false answer.

7. **`writeFile` is manifest-guarded.** A file the project has edited comes back
   as a conflict rather than being overwritten. Rely on that; do not add your own
   existence check.

8. **A repair is idempotent or it is a bug.** `*_repair.go` runs on every
   upgrade, so it has to recognise its own output and do nothing.

9. **Compilation is not verification.** A template that compiles can still write
   a Tailwind class the app does not define (which renders as nothing), a label
   pointing at no control, or an import the project cannot resolve.

## Validation workflow

```bash
go build ./... && go test ./internal/scaffold/
```

Then, and this is the part that matters:

1. Rebuild **every** installed `grit.exe` before testing a template edit.
   `go build ./...` does not ship a scaffold change to the binary on your PATH.
2. `grit new` a fresh project in `D:\LEARNING\Grit Framework\Framework Tests`,
   in the architecture you touched **and** one you did not.
3. Install, migrate, and **stand it up on spare ports** (`REDIS_URL=`,
   `MAIL_MAILER=log`). Open the pages you changed and look at them.
4. If you touched a form: tab through it. Twelve of thirteen fields once
   shipped with a label attached to nothing.
5. If you touched something framework-owned: run `grit upgrade` on a project
   scaffolded by the **previous** release and confirm the file arrived.

Tests that match strings prove the template contains a string. They do not prove
the app works.
