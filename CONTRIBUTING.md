# Contributing to Grit

Thank you for being here. This file is the short version of everything you need
to send a change that gets merged.

## Before you write code

**Open an issue first for anything beyond a fix.** Grit generates code into
other people's projects, so a change here lands in every project built after
it and, through `grit upgrade`, in many built before. A feature that seems
small in the CLI can be a migration for a thousand applications. An issue
costs you ten minutes and can save you a weekend.

Fixes need no issue. If something is wrong, send the fix.

**Questions and ideas go to [Discussions](https://github.com/MUKE-coder/grit/discussions).**
The issue tracker is for defects and agreed work.

**Security problems do not go in an issue.** Read
[SECURITY.md](SECURITY.md) and report privately.

## Getting set up

You need Go 1.21+, Node 20+, pnpm and Docker.

```bash
git clone https://github.com/MUKE-coder/grit.git
cd grit
go build ./...
go test ./...
```

That is the whole setup. The CLI is a single Go binary with no code generation
step and no vendored dependencies.

To try your change against a real project:

```bash
go build -o /tmp/grit ./cmd/grit
cd /tmp && /tmp/grit new scratch --triple --next --db sqlite
cd scratch/apps/api && go build ./...
```

## How the repository is laid out

| Path | What lives there |
|------|------------------|
| `cmd/grit/` | The CLI: every command, and nothing else |
| `internal/scaffold/` | The templates. Every file a generated project contains is a Go function returning a string here |
| `internal/generate/` | `grit generate resource`: parsing a field spec into a model, service, handler, schema, types, hooks and an admin screen |
| `internal/prompt/` | The interactive `grit new` |
| `docs/` | The documentation site (Next.js), deployed at gritframework.dev |
| `demo/` | Grit Motors, the application behind demo.gritframework.dev |
| `ui/` | The Grit UI block registry |
| `examples/` | Complete applications, each with a guide |

**Generated code lives in `internal/scaffold`, not in a template directory.**
A template is a Go function returning a backtick string, with `{{MODULE}}`
where the project's module path goes. This is deliberate: the templates
compile, so a typo in a Go template is a build failure here rather than a
broken project for somebody else.

## Writing the change

**Match the surrounding code.** Comment density, naming, error handling and
idiom are already established in whatever file you are editing. A change that
reads like the file it is in is easier to review than a better change that
does not.

**Go:** handle every error, wrap with `fmt.Errorf("context: %w", err)`, keep
handlers thin and put logic in services. Run `gofmt`.

**TypeScript and React:** function components, hooks, React Query for data,
Zod for validation, no `any`. Tailwind utilities, and only colours the theme
defines: an unknown utility compiles to nothing and fails silently.

**Comments say why, not what.** The interesting comment in this codebase is
the one that explains the failure a piece of code exists to prevent.

## Tests

Every behavioural change needs a test. The bar is not coverage, it is this:
**if somebody deletes your change, does a test fail?**

Most tests in `internal/scaffold` assert on the emitted template strings. That
catches a removed guard but not a TypeScript error, so for a frontend change,
generate a project and run `npx tsc --noEmit` against it before you send the
pull request.

```bash
go test ./...                       # everything
go test ./internal/scaffold/ -run TestNameOfYours -v
```

## Sending it

- One change per pull request. Two unrelated fixes are two pull requests.
- Write the commit message for somebody reading it in two years with no
  context: what was wrong, and why the fix is the fix. Conventional commit
  prefixes (`feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`).
- `go test ./...` passes and `gofmt -l .` is quiet.
- Do not bump the version or edit the changelog. Releases are cut from `main`
  by the maintainer.

## What gets merged quickly

Fixes with a test that fails without them. Documentation that corrects
something wrong. A generated-project bug reproduced in a scratch project.
Accessibility fixes. Anything that removes code without removing behaviour.

## What takes longer

A new dependency, in the CLI or in a generated project. Every package Grit adds
is one that every project built with it carries, so the question is never "is
this a good library" but "is this worth its weight in a thousand projects".

A new architecture, database or frontend framework. These multiply the matrix
that has to keep working.

A change to a generated project's structure. It is a migration for everybody
on `grit upgrade`.

None of these are refused on principle. They need a conversation first, which
is what the issue is for.

## Code of Conduct

Taking part in this project means agreeing to the
[Code of Conduct](CODE_OF_CONDUCT.md).
