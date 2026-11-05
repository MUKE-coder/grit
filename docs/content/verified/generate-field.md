---
page: /docs/concepts/field-types
title: grit generate field adds a column to a resource that already exists
---

The field-types page lists the field types and says a column can be added to an
existing resource in place. "In place" is the hard part: the field has to reach
the model, the Zod schema, the TypeScript type, the admin's form and table, and
the service's column whitelist. A field that reaches the model and not the
whitelist is a column the API will not write, and nothing says so.

Three types here rather than one, because they take different paths: a plain
string, a `select` (options the admin draws, and a constraint), and a `toggle`
(a bool, which has its own trap with database defaults).

The select is written in the readable form on purpose. The positional form needs
quoting, because an unquoted `|` is a pipe: the first version of this flow had
`state:select:draft=Draft|active=Active|retired=Retired` bare, the shell cut it
at the first pipe, and the flow passed having built a one-option select. A flow
that passes while asserting less than it says is the failure mode to watch for
in this directory.

```bash
grit generate resource Part --fields "name:string"
grit generate field Part sku:string
grit generate field Part "state:select[draft, active, retired]"
grit generate field Part tracked:toggle
(cd apps/api && go build ./... && go vet ./...)
grit migrate
grit doctor
```
