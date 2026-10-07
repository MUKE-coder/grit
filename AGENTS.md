# AGENTS.md

Start here if you are an AI agent, or a person, about to change this
repository. Grit is a full-stack meta-framework: a Go CLI that writes Go (Gin,
GORM) and Next.js or TanStack projects.

The one thing to understand before anything else: **almost nothing in
`internal/` is code that runs for a user. It is code that writes code.** A
change here compiles cleanly and reaches people as a file in their project that
is wrong, missing, or subtly different from the file next to it. That is why the
verification steps below are not optional, and why every directory has its own
rules.

## Read in this order

| Document | What it gives you |
|---|---|
| [GRIT.md](GRIT.md) | The master specification: what Grit is, its architecture and its folder structure. |
| [PHASES.md](PHASES.md) | The live remaining-work list. Read it before answering "what's left". |
| [TERMINOLOGY.md](TERMINOLOGY.md) | One name per thing. Settles the words with two meanings here, and the predicate traps that cost the most. |
| [API_DESIGN.md](API_DESIGN.md) | The HTTP contract the generated API honours, as a contract. |

Then the file for the directory you are editing:

| File | Covers |
|---|---|
| [internal/scaffold/AGENTS.md](internal/scaffold/AGENTS.md) | The project writer. Delivery's two lists, the path helpers, string templates. |
| [internal/generate/AGENTS.md](internal/generate/AGENTS.md) | The resource generator. Injection anchors, and why every injection needs a removal. |
| [cmd/grit/AGENTS.md](cmd/grit/AGENTS.md) | The CLI surface, the three version spots, and the release ritual. |
| [internal/ui/AGENTS.md](internal/ui/AGENTS.md) | The CLI design system: seven colour roles, ANSI fallbacks, `NO_COLOR`. |
| [docs/AGENTS.md](docs/AGENTS.md) | The docs site. Generated tutorials, the changelog, no em dashes. |
| [examples/AGENTS.md](examples/AGENTS.md) | The checked-in applications. They are output: fix the template, not the example. |

## The rules that hold everywhere

1. **Delivery has two lists.** A framework-owned file needs the new-project map
   *and* the upgrade writer. A file in only one reaches half the projects, and
   the failure is silent. Prefer one `writeXFiles(root, opts)` writer called
   from both.

2. **Compilation is not verification.** `go build` passing says nothing about
   whether the generated app runs. Scaffold a project, install it, stand it up,
   and open the page you changed.

3. **Rebuild the binary before testing a template change.** There are three
   copies of `grit.exe` on this machine, and `go build ./...` updates none of
   the ones on your PATH. Check with `grit version`.

4. **Stage new files explicitly** and verify with `git diff --cached --stat`
   before committing. A ` M` in `git status` means *unstaged*, and that shipped
   one release empty.

5. **No em dashes** in anything a person reads: docs, UI copy, commit messages,
   comments.

6. **Bugs in GORM Studio, Sentinel or Pulse get an issue on their repository**,
   not a workaround here. They are separate projects of ours.

7. **Comments explain why, not what.** The codebase is dense with the reasons
   behind decisions, because the decisions are mostly about failures that were
   invisible. Keep that up: a comment saying what a line does is noise, one
   saying which bug it prevents is the only record of that bug.

## Before you say you are done

```bash
go build ./... && go test ./...
cd docs && npx next build        # if you touched docs/
```

Then the project-level check appropriate to what you changed, as described in
the directory's own AGENTS.md. For anything user-facing, that means a generated
project you actually looked at.
