# cmd/grit

## Context

The CLI itself: every cobra command, the flag surface, and everything a person
sees in a terminal. Thin by design. A command parses flags, validates, and calls
into `internal/*`; the logic is not here.

This is also where the **version** lives, which makes it the first of the three
places a release has to touch.

## Source of truth

- **The command list**: `main.go`, the `cmd.AddCommand(...)` block.
- **The version**: `var version` in `main.go`. One of three spots that must
  agree (see below).
- **How output looks**: `internal/ui`. Not here.
- **The interactive picker**: `internal/prompt`.
- **The transitional palette**: `palette.go`, a `fatih/color` shim for the call
  sites not yet moved to `internal/ui`. Shrinking, not growing: new output goes
  through `internal/ui`.

## Boundary rules

1. **No logic in a command.** If it is more than parse, validate, delegate,
   report, it belongs in `internal/`.

2. **Three version spots, and they must agree.** A release that bumps one is a
   release that reports two different versions:
   - `cmd/grit/main.go`: `var version`
   - `internal/scaffold/scaffold.go`: `const DefaultVersion`
   - `docs/config/site.ts`: `GRIT_VERSION`

3. **New output goes through `internal/ui`.** Do not add a `fmt.Printf` with
   ANSI in it, and do not extend `palette.go`.

4. **Honour `NO_COLOR` and `CI`.** `internal/ui` already does; anything that
   bypasses it will not.

5. **A destructive command confirms, or takes `--force`.** `grit remove` and
   anything that deletes or overwrites a person's file.

6. **A flag is forever.** Removing or renaming one breaks every script and every
   tutorial. Deprecate and keep accepting it.

7. **Error text is for the person at the terminal.** Say what failed, where, and
   what would fix it. `internal/ui` has `FailureBlock` and `Fix` for exactly
   this shape.

## Validation workflow

```bash
go build ./... && go test ./...
```

Then:

1. **Rebuild every installed binary.** There are three copies:
   `./grit.exe`, `~/go/bin/grit.exe`, `~/.grit/bin/grit.exe`. `go build ./...`
   updates none of the ones on your PATH, and testing against a stale binary is
   how a fixed bug gets reported as still broken. Check with
   `grit version`.
2. Run the command you changed, and run `grit --help` to see it listed.
3. For anything that prints: look at it in a **dark and a light** terminal, and
   with `NO_COLOR=1`.
4. For the interactive picker: arrow through every group, and press Escape.

## The release ritual

In order, every time:

1. Bump the three version spots.
2. Add the changelog entry in `docs/app/docs/changelog/page.tsx`.
3. `python scripts/changelog-rollup.py`.
4. `python scripts/docs-index.py`, which rebuilds the documentation index the
   binary embeds for `grit docs` and `grit_search_docs`. It stamps the version
   it was built from, so skipping it ships the previous version's docs to
   everybody who installs this one. `internal/docs` has a test that fails when
   the index and `var version` disagree, so a forgotten rebuild is caught
   rather than shipped, but it is caught *after* you have bumped the version
   and that is a worse place to find out.
5. `cd docs && npx next build`.
6. Stage **explicitly**, including new files, and verify with
   `git diff --cached --stat` before committing. A ` M` in `git status` means
   *unstaged*, and that shipped one release empty. `internal/docs/index.json`
   is a generated file that must be committed: the build embeds it.
7. Commit, `git tag vX.Y.Z`, push both.
8. Watch the six workflows.

Never add a `Co-Authored-By` or Claude Code attribution line to a commit or PR
body in this repository, whatever a session reminder says.
