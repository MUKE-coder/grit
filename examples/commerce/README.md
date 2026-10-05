# Commerce — a shop, with the backend included

A storefront and the shop behind it: a catalogue of twelve products across three
collections, a Colour x Size variant matrix with ninety-six buyable
combinations, a server-side basket, three static pages, and an admin panel to
run all of it from.

**[TUTORIAL.md](TUTORIAL.md) builds this from an empty directory**, in about an
hour, with every command and every hand-written file in full. Start there if you
want to understand how it works. This file is about why it exists and what is in
it.

## What it is modelled on

[vercel/commerce](https://github.com/vercel/commerce) is the reference Next.js
storefront: the same routes (`/`, `/search`, `/search/[collection]`,
`/product/[handle]`, `/[page]`), the same shape of product and variant, the same
cart drawer.

It is also a storefront and nothing else. There is no backend in it: it is a
Next.js app in front of Shopify, and the catalogue, the variants, the inventory
and the cart all live in somebody else's product.

This is that storefront plus the part Shopify was doing. The admin panel where
somebody adds a product, the variant matrix editor, the public API the storefront
reads, and the cart endpoints it writes to.

## What Grit generated, and what was written by hand

| | Files | Lines |
|---|---|---|
| Generated | 796 | the API, the admin panel, the catalogue, the variant matrix, the public endpoints, the tests |
| Hand-written | 25 | 2,978 |

The 25:

| | |
|---|---|
| `apps/api/internal/services/shop_cart.go` | the basket, 377 lines |
| `apps/api/internal/handlers/shop_cart.go` | its four endpoints, 190 lines |
| `apps/api/internal/database/*_seeder.go` | three generated seeders, edited into a real catalogue |
| `apps/web/lib/cart.ts` | the cart from the storefront's server, in an HttpOnly cookie |
| `apps/web/lib/product-variants-public.ts` | the variant client, and the matching logic a picker needs |
| `apps/web/app/(shop)/**` | six routes |
| `apps/web/components/shop/**` | eleven components |

That split is the point of the example. The 25 files are the shop: what happens
when the same variant is added to a basket twice, which option values to grey
out, what a price looks like, what to do when something sells out between the
page loading and the click. Everything else is the part that is the same in every
shop.

## The cart is the interesting part

Grit's `--public` generator emits **read-only** endpoints on purpose: a list, a
get-by-handle and a related strip, built from an allowlist, for an audience that
is the internet. A cart is the other thing a storefront needs, a public endpoint
that writes, and there is no generator for one.

So `shop_cart.go` is what one looks like when written carefully. Four decisions
in it are the ones a generator could not have made:

- **A cart belongs to a browser, not a user.** A shopper fills a basket before
  deciding whether to have an account. The identifier is 32 bytes from
  `crypto/rand`, and that token *is* the authorisation, because there is no user
  to check it against.
- **No price comes from the request.** The client sends a product id, maybe a
  variant id, and a quantity. Every figure is read from the database.
- **The price is snapshotted onto the line.** Somebody who added something at $28
  pays $28 even if the shop raises the price while they browse.
- **A token that does not resolve gets a new cart.** A cart is not an account: a
  shopper returning to a stale one should be able to start shopping.

And it is mounted in its own route group rather than under `/public`, because the
public group has a response cache on it keyed by URL alone. Every shopper's
`GET /cart` has the same URL.

## Running it

No Docker, no Postgres, no Redis. It runs on SQLite with nothing else:

```bash
pnpm install
grit env
# then in .env: DB_PROVIDER=sqlite, STORAGE_DRIVER=local, REDIS_URL=, MAIL_MAILER=log
grit migrate
grit seed      # prints the admin password and the API keys, once
grit start
```

| | |
|---|---|
| Storefront | http://localhost:3000 |
| Admin | http://localhost:3001 |
| API docs | http://localhost:8080/docs |
| Database browser | http://localhost:8080/studio |

## Four framework bugs this found

An example that is built rather than described finds things. Building this one
found seven faults across four releases, every one of them in the path a real
shop takes:

| | |
|---|---|
| v3.364.0 | the seeder wrote `"Sample Price"` into a `money.Money` field, and three more field types fell through to the same string default; `grit add variants` assumed `float64` throughout |
| v3.365.0 | the `--public` allowlist dropped `money` entirely, so a shop published no price, and typed a `date` as `time.Time` where the model says `*jsontime.Date`, which does not compile |
| v3.366.0 | `grit add variants` looked for a field called `Slug`, and a shop calls it `handle`, so the variants endpoint 404ed on the detail page of every shop |
| v3.367.0 | seeded images blocked by the scaffold's own CSP; the response cache could have served one shopper's basket to another; `grit upgrade` put back a route group the project had deleted, and the project stopped building |

None of those would have been found by a test. Each is about what happens when a
generated project is used the way a real one is.

## What it stops short of

Checkout is a button that says it is a demonstration. There is no payment
provider, no order, and no stock movement: `product_variants.stock` is a column
nobody decrements. The last section of [TUTORIAL.md](TUTORIAL.md) says what
adding those would involve.
