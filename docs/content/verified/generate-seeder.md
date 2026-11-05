---
page: /docs/backend/seeders
title: A generated seeder is registered, not just written
---

The seeders page says a resource generated with `--seed --faker --count N`
arrives with a seeder that `grit seed` runs. Three things have to line up: the
seeder file, its registration in `internal/database/seed.go`, and the faker
knowing what to put in each column from the field types.

A seeder that is written and not registered is the quiet failure. `grit seed`
succeeds, says nothing about it, and the table stays empty. So this flow names
the resource: `grit seed Supplier --count 20` has to find a seeder registered
under that name, and tops the table up to twenty rows rather than taking the
run's word for it.

```bash
grit generate resource Supplier --fields "name:string,email:email" --seed --faker --count 20
grit migrate
(cd apps/api && go build ./...)

grit seed
grit seed Supplier --count 20
```
