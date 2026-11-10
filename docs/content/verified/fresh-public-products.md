---
page: /docs/frontend/fetching-data
title: A public Product resource, and the two files the web app reads it with
---

The page claims that one `--public` flag gives a storefront everything it needs
to read a resource: the two endpoints, an allowlist that is not the model, and a
typed read client in the web app that already sends the API key. Each of those
is a file the generator has to write into the right app, and the web half is the
part that was missing for a while (grit#91): the API side worked, the hooks
beside it called the authenticated routes, and a storefront built the obvious way
got 401s from every request.

So this checks the files exist and say what the page shows, rather than only that
the command exits 0. The greps are the assertions:

- the allowlist is a `publicProduct` struct in the handler, which is what makes a
  column added next month private by default
- `stock` is **not** in it, because the page tells the reader it was held back and
  names it
- the read client exists in `apps/web/lib/`, sends `X-API-Key`, and exports both
  functions the two pages call

One block, not several, because a `fresh-*` flow starts in an empty directory and
every later command has to run inside the project the first one created. A `cd`
that leaks into the next block is the first of the two ways a flow lies.

```bash
grit new shop --triple --next
cd shop

grit generate resource Product \
  --fields "name:string,slug:slug[name],description:text,price:money,image:file[image],featured:bool,stock:uint" \
  --public --faker --count 24

(cd apps/api && go build ./... && go vet ./...)

# The allowlist, and the field the page says is held back.
grep -q "type publicProduct struct" apps/api/internal/handlers/product_public.go
grep -q "Featured    bool" apps/api/internal/handlers/product_public.go
! grep -q "Stock" apps/api/internal/handlers/product_public.go

# The web app's read client: the key, and the two functions the pages import.
grep -q "X-API-Key" apps/web/lib/products-public.ts
grep -q "export const getPublicProducts" apps/web/lib/products-public.ts
grep -q "export const getPublicProduct " apps/web/lib/products-public.ts
grep -q "/api/v1/public/products" apps/web/lib/products-public.ts

# And the money type reaches the published shape, because the page tells the
# reader to format it with formatMoney rather than dividing by 100.
grep -q "price: Money" apps/web/lib/products-public.ts

grit migrate
grit seed
grit doctor
```
