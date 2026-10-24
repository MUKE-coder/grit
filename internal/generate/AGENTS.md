# internal/generate

## Context

The resource generator: `grit generate resource Post title:string body:text`,
and its inverse `grit remove Post`. It runs against a project that **already
exists**, so unlike `internal/scaffold` it cannot simply write files. It writes
some and **injects into others**: `routes.go` gains a route group, the sidebar
gains a link, `packages/shared` gains a schema, the admin gains a resource
definition and a page.

Also here: `grit sync` (Go structs → TypeScript and Zod), `grit add field`, and
the generators for jobs, mail, factories and column packs.

## Source of truth

- **Anchors**: the `// grit:*` comments. `grep -rn '"// grit:' internal/generate`
  lists every one. The anchor is a contract with the scaffold, which has to emit
  it.
- **Injection**: `inject.go` (`injectBefore`, `injectInline`).
- **Removal**: `remove.go`, and the per-feature `*_remove.go` files.
- **Go type → TS and Zod**: `sync.go` and `field.go`.
- **The end-to-end shape**: `generator.go`, `Generator.Run()`.

## Boundary rules

1. **Every injection needs a matching removal.** `grit generate` and
   `grit remove` are a pair. A new injection with no removal leaves a dangling
   import or a sidebar link to a deleted page, and the project stops compiling
   on an operation that is supposed to clean up. Add both in the same change.

   This rule was here, in these words, while policies broke it for months.
   `generate` wrote `internal/policies/<snake>.go` **and** a `register<Name>s()`
   call in `policies.go`; `remove` took back neither, so it printed
   `Resource Widget removed successfully!` and left a tree that would not
   build (fixed in v3.411.1). Two things to take from it. A feature that
   injects in *two* places needs *two* removals, and taking back one of them
   only swaps one compile error for another. And the round-trip test that
   should have caught it reversed each injection by hand rather than calling
   the remover, so it passed throughout: a test of removal has to call the
   code `grit remove` calls.

   Reading the generator against the remover file by file then found four more
   (v3.411.2): the seeder, two `--tree` files and the tree test, and the
   workflow service a status field writes. One of those gives the rule its
   second half. **A path built from a name variant has to be removed using the
   same variant.** `seeder.go` writes `<plural-kebab>_seeder.go` and `remove.go`
   deleted `<plural_snake>_seeder.go`, which is the same string for a one-word
   resource and a different one for every other, so `grit remove Widget` was
   clean and `grit remove BlogPost` was not. Every fixture used a one-word name.
   **A one-word fixture cannot tell two name variants apart**, so test a removal
   with a two-word resource.

   Three releases of finding these one at a time is enough, so there is now a
   test that asks the question the command claims to answer:
   `TestGenerateThenRemoveLeavesNothingBehind` generates a two-word resource
   with `--tree`, `--seed` and a workflow, removes it with the real remover,
   and fails on **anything left in the project that still names it**, file or
   line. It found the tree routes (swept by `<camel>Handler.`, written as
   `<camel>TreeHandler.`), `internal/workflow/<name>.go`, and the manifest
   entries nothing had ever called `Forget` for. **Add your injection, then run
   that test**: it is cheaper than the compiler output in somebody's project.

2. **Build an emitted import block as a list, never with `strings.Replace` on an
   anchor line.** A `Replace` whose anchor has moved matches nothing, returns the
   input unchanged and reports success. See `go_import.go` for the shape to use.

3. **An injection is idempotent or it is a bug.** Generating the same resource
   twice must not produce two routes. Check for your own output first.

4. **Anchors are load-bearing and shared.** Renaming one means the scaffold stops
   emitting a hook the generator needs, and the generator silently injects
   nothing. If you add an anchor, add it to the scaffold's template in the same
   commit.

5. **Never call a service from a model.** `models` → `services` is an import
   cycle. Auto-numbering in `BeforeCreate` calls `sequence.Next` directly, never
   `services.NextXNumber`.

6. **Generate every field type before trusting a change to type mapping.** One
   resource with all 30 field types, then compile it. A mapping that is wrong for
   `money` or `json` is invisible in a resource that uses neither.

7. **The response shape, the Zod schema and the TypeScript type change
   together**, in one commit. That is what `grit sync` is for.

## Validation workflow

```bash
go test ./internal/generate/
```

The unit tests here are good and fast; use them. Then, against a real project:

1. `grit generate resource Thing ...` in a fresh project.
2. `cd apps/api && go build ./...`: the Go injections have to compile.
3. `pnpm typecheck`: so do the TypeScript ones.
4. Open the admin list page, the create form and the detail page.
5. `grit remove Thing`, then build and typecheck **again**. This is the step
   that catches a missing removal, and it is the step most often skipped.

For a change to type mapping, generate a resource carrying **all** field types
and compile that.
