---
page: /docs/testing
title: Not a flow. This file explains the directory.
---

# Declared flows

A file in this directory is a sequence of commands the docs claim works, run by
CI against a real scaffolded project on a real Postgres.

```bash
# Every flow, or one of them.
tests/docs/run-flows.sh
tests/docs/run-flows.sh remove-resource

# What is declared, and whether each names a page that exists.
go run ./tools/docsamples -list
go run ./tools/docsamples -check
```

## Why a file rather than a block in the page

A page can also mark a block `verify="<flow>"`, and four of them do. That stayed
at four while the docs grew to 193 pages, because adding one means editing a
React component: finding the right place in the JSX, and getting a template
literal inside it right. The two numbers were the finding in grit#168.

A file here is ten lines. The page it checks is named in the frontmatter, and
`-check` fails if that page has moved, so the flow does not quietly stop being
evidence about anything a reader can find.

## The format

```
---
page: /docs/backend/seeders
title: One line, shown by -list
---

Prose, for whoever reads this file next. Say what the flow proves and what it
would catch, not what the commands are: the commands are below.

```bash
grit generate resource Invoice --fields "total:money"
grit migrate
```
```

Rules worth knowing:

- The flow's name is the file name. `_`-prefixed files, like this one, are
  skipped.
- Every `bash` or `sh` block runs, in file order, under `set -euo pipefail`. A
  block in any other language is documentation.
- The working directory is a freshly scaffolded `grit new app --api` that has
  already been migrated, which is where a reader of the page is. `DATABASE_URL`
  points at a database of this flow's own.
- A flow named `fresh-*` starts in an empty directory instead, because its first
  command is the one that creates the project.
- A flow with no commands is refused. Passing without running anything is the
  one outcome worse than having no flow.

## Two ways a flow lies, both found writing the first four

**`cd` without a subshell.** `cd apps/api && go build ./... && cd ..` from the
project root leaves you in `apps/`, not where you started, and the next block
fails for a reason that has nothing to do with what it was testing. Use a
subshell, which cannot drift:

```bash
(cd apps/api && go build ./... && go vet ./...)
```

**An unquoted `|`.** A positional select spec, `state:select:a=A|b=B`, is a
pipe to the shell. The first version of `generate-field.md` had one bare: the
shell cut it at the first pipe, `grit` accepted a one-option select, and the
flow passed having asserted less than it claimed. Quote the spec, or use the
readable form, `"state:select[a, b]"`.

A flow that passes while asserting less than it says is worse than no flow,
because it is counted. Read the output of a new flow once before trusting it.
