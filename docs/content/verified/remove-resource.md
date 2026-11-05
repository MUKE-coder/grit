---
page: /docs/cli
title: grit remove resource undoes every injection
---

`grit generate resource` writes files and injects into about a dozen others: the
model registry, the routes, the admin sidebar, the permission catalogue, the
access table. `grit remove resource` has to reverse all of it, and the failure
mode is a project that no longer compiles because one injection was left behind
referring to a file that is gone.

That has happened more than once, which is why this is a flow rather than a
claim. The build and the vet at the end are the assertion, and `go vet` matters
as much as `go build`: `go test` runs a subset of vet as an error, so a project
that builds can still have a package whose tests refuse to.

```bash
grit generate resource Widget --fields "name:string,price:money"
grit migrate
(cd apps/api && go build ./... && go vet ./...)

grit remove resource Widget --force
(cd apps/api && go build ./... && go vet ./...)
```
