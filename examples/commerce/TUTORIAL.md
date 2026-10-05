# Build a shop with Grit

This builds the application in this folder, from an empty directory, in about an
hour. At the end you have a storefront with a catalogue, a variant matrix, a
working basket and an admin panel to run it from, and you will have written
about 3,000 lines of the roughly 124,000 in the result.

It assumes you are comfortable with Go and with React, and that you have never
used Grit. Every command is given in full, every file you write is given in
full, and where something went wrong while this was being built the first time,
it says so, because those are the parts you are most likely to hit.

**What you are building:** [vercel/commerce](https://github.com/vercel/commerce)
is a storefront with no backend: it is a Next.js app in front of Shopify. This
builds the same storefront and the backend it talks to, which is the part
Shopify was doing.

---

## Contents

1. [Before you start](#1-before-you-start)
2. [Scaffold the project](#2-scaffold-the-project)
3. [Generate the catalogue](#3-generate-the-catalogue)
4. [Add the variant matrix](#4-add-the-variant-matrix)
5. [Run it for the first time](#5-run-it-for-the-first-time)
6. [Write the seed data](#6-write-the-seed-data)
7. [The cart, which Grit does not generate](#7-the-cart-which-grit-does-not-generate)
8. [The storefront](#8-the-storefront)
9. [What to do next](#9-what-to-do-next)

---

## 1. Before you start

You need:

| | |
|---|---|
| Go | 1.21 or newer |
| Node | 20 or newer |
| pnpm | `npm install -g pnpm` |
| Grit | `go install github.com/MUKE-coder/grit/v3/cmd/grit@latest` |

Check the CLI is on your path:

```bash
grit version
```

This tutorial was written against **v3.367.0**. Earlier versions have four bugs
this build found, all of which you would hit: seeded money prices, the public
API's price field, the variant lookup by handle, and placeholder images blocked
by the project's own security policy. Use v3.367.0 or later.

You do **not** need Docker, Postgres or Redis. The project runs on SQLite with
no services at all, which is how this tutorial runs it.

---

## 2. Scaffold the project

```bash
grit new commerce --triple --next
cd commerce
```

`--triple` gives three applications in one monorepo:

| | |
|---|---|
| `apps/api` | Go, Gin and GORM |
| `apps/web` | the storefront, Next.js |
| `apps/admin` | the admin panel, Next.js |

`--next` picks Next.js for the admin rather than Vite. A shop's admin is used by
a handful of people on good connections, so either works; Next.js keeps both
frontends on the same framework, which is one less thing to remember.

That writes about 680 files. Most of them you will never open: authentication,
sessions, API keys, file storage, background jobs, mail, rate limiting, audit
logging, health checks, OpenAPI, an admin panel, and the tests for all of it.

### Point it at SQLite, and at a folder on disk

Create the environment file:

```bash
grit env
```

Then edit `.env` and change two lines:

```ini
DB_PROVIDER=sqlite
REDIS_URL=
```

`REDIS_URL=` set to nothing means "run without Redis": the cache, background
jobs and cron are skipped rather than failing to connect. Add two more lines at
the end:

```ini
MAIL_MAILER=log
STORAGE_DRIVER=local
```

`MAIL_MAILER=log` prints password-reset mails to the log instead of sending
them. The second line is the one worth understanding, because it is where a
shop meets its first wall.

#### Where the pictures go

A shop is photographs. The moment you add a product in the admin, something has
to store a JPEG, make a thumbnail, and hand back a URL.

Grit can put those in four places, and `STORAGE_DRIVER` picks which:

| | |
|---|---|
| `local` | a folder on this machine, served by the API from `/files`. No bucket, no Docker, nothing else to run. |
| `minio` | the MinIO container in `docker-compose.yml`, which is S3 in a box for local development. The default. |
| `s3`, `r2`, `b2` | AWS S3, Cloudflare R2, Backblaze B2. What a deployed shop uses. |

**`local` is the one to pick if you have no object storage.** Files land in
`apps/api/storage/app/`, a directory you can open and look inside:

```
apps/api/storage/app/
├── originals/2026/10/<uuid>.png      the file as it was uploaded
└── uploads/2026/10/<uuid>.jpg        the optimised version
    └── <uuid>-thumb.jpg              and its thumbnail
```

Everything else works the same: the API optimises the image, makes the
renditions, signs temporary URLs for private files, and serves public ones from
`APP_URL/files/...`. A `FileRef` from the local disk and one from an S3 bucket
are the same shape, which is why moving to a bucket later changes one line of
`.env` and nothing in your code.

Keep the directory out of git:

```bash
echo "apps/api/storage/app/" >> .gitignore
```

> **You may not need the line at all.** Outside production, a `minio` that is
> not answering falls back to the local disk by itself and says so in the log:
>
> ```
> MinIO at http://localhost:9000 is not answering, so files are being kept on
> the local disk at storage/app and uploads work.
> ```
>
> Setting `STORAGE_DRIVER=local` makes it your decision rather than a fallback,
> which is worth doing if you have no intention of running Docker. Production
> never falls back: a server that has lost its bucket says so rather than
> quietly writing to a disk nobody is backing up.

> **If an upload says `File storage is not configured`**, you are on a version
> before v3.370.0. The fallback above tested for a missing `MINIO_ACCESS_KEY`,
> and `grit` writes that key into `.env`, so it could never fire. Either
> `grit upgrade`, or set `STORAGE_DRIVER=local` by hand, which works on every
> version.

#### One thing not to change

Leave `APP_PORT` alone. `APP_URL` is written into `.env` as
`http://localhost:8080`, and the two have to agree: `APP_URL` is what every
generated link is built from, uploaded files included, so a server listening on
one port while `APP_URL` names another serves images from an address where
nothing is answering. If you do need a different port, change both.

---

## 3. Generate the catalogue

A shop is three things: products, the collections they are grouped into, and the
static pages in the footer.

### Collections

```bash
grit generate resource Collection \
  --fields "title:string,handle:slug,description:text,image:file:image" \
  --public --seed
```

### Products

```bash
grit generate resource Product \
  --fields "title:string,handle:slug,description:richtext,price:money,featured_image:file:image,images:files,available:bool,collection:belongs_to:Collection" \
  --public --seed
```

### Pages

```bash
grit generate resource Page \
  --fields "title:string,handle:slug,body:richtext" \
  --public --seed
```

Three flags and three field types are doing real work here.

**`handle:slug`** is the URL-safe name. The generated model fills it in from the
title on insert, keeps it unique, and the public endpoint looks a row up by it,
so `/product/heavyweight-cotton-tee` works with no routing code. Call it
`handle` rather than `slug` if you like; the generator reads the name you chose.

**`price:money`** is `{ amount, currency }` with the amount in minor units.
`2800 USD` is $28.00. Nothing in the shop ever stores a price as a float, which
is the entire reason the type exists: `0.1 + 0.2` is not `0.3`, and a shop that
adds up floats is a shop with a subtotal that is one cent out often enough to
notice. The admin shows a currency field, the API publishes both halves, and the
shared TypeScript package has a `formatMoney` that knows which currencies have
no decimal places at all.

**`--public`** adds a second, narrower API: a read-only list and get-by-handle
under `/api/v1/public/products`, guarded by a publishable API key, built from an
allowlist rather than from the model. A column added next month is private until
somebody adds it to that allowlist, which is the right default when the audience
is the internet. The staff CRUD the admin panel uses is separate and behind
authentication.

**`--seed`** writes a seeder with one sample row in it, which you will replace in
step 6.

---

## 4. Add the variant matrix

A tee comes in four colours and four sizes. That is sixteen things you can buy,
each with its own stock level and possibly its own price, and it is the part of
a shop that is tedious to build and easy to get wrong.

```bash
grit add variants --resource Product
```

That installs five tables and the machinery over them:

| | |
|---|---|
| `options` | Colour, Size: shop-wide, shared between products |
| `option_values` | Black, White, S, M, L, XL, each with a swatch colour and a price delta |
| `product_options` | which options this product offers |
| `product_variants` | one row per buyable combination, with SKU, stock and an optional price override |
| `variant_option_values` | which values each variant is |

It also adds a matrix editor to the product's page in the admin, a shared option
library under Options in the sidebar, and one public endpoint:

```
GET /api/v1/public/products/:handle/variants
```

That endpoint is worth understanding, because the storefront is built on it. It
returns three things in one response:

```jsonc
{
  "options":  [ /* the axes to draw, with values, swatches and deltas */ ],
  "variants": [ /* every combination: id, sku, price, in_stock, option_value_ids */ ],
  "price_range": { "low": …, "high": …, "single": false }
}
```

One request, not three. A picker that fetched the options, then the variants,
then a price on every swatch click would make three round trips per interaction
on the busiest page in the shop, and the figure it landed on would still be the
browser's arithmetic rather than the server's.

A product with no variants gets empty lists and a range of its own price, which
is what lets the storefront render one component whether or not options were
ever set up. Most of a real shop is mugs and tote bags.

---

## 5. Run it for the first time

```bash
grit migrate
grit seed
```

`migrate` creates the tables from the models. `seed` fills them, creates an admin
user and mints two API keys, writing the publishable one into
`apps/web/.env.local` and `apps/admin/.env.local` so the frontends can reach the
public API without you copying anything.

Read the output of `grit seed`. It prints the admin email and password once, and
both API keys once.

Now start everything:

```bash
pnpm install
grit start
```

| | |
|---|---|
| API | http://localhost:8080 |
| API docs | http://localhost:8080/docs |
| Database browser | http://localhost:8080/studio |
| Storefront | http://localhost:3000 |
| Admin | http://localhost:3001 |

Sign in to the admin with the credentials `grit seed` printed. You have a working
product catalogue with a variant matrix editor, and a storefront that is still
the scaffold's landing page.

> **If every page says `INVALID_API_KEY`:** you ran `grit migrate --fresh` at
> some point, which drops the keys table. The seeder minted a new key and
> deliberately did not overwrite your `.env.local`, because overwriting
> somebody's environment file is how you lose an afternoon. Since v3.367.0 it
> tells you, and prints the key to paste in. Before that it was silent, and this
> cost half an hour the first time.

---

## 6. Write the seed data

`--seed` wrote three seeders with one row each, called "Sample Title". They exist
to be edited: the generator knows the shape of your model and nothing about your
shop.

Open `apps/api/internal/database/collections_seeder.go` and replace the records
with three real collections. The version in this folder is 72 lines and adds one
helper, `shopImage`, which builds a placeholder URL seeded by a string you choose
so the same product keeps the same picture across a reseed.

Then `products_seeder.go` with twelve products, and `pages_seeder.go` with three
pages. Copy them from this folder, or write your own; the shapes are the only
part that matters:

```go
records = append(records, models.Product{
    Title:         e.title,
    Handle:        e.handle,
    Description:   e.description,
    Price:         money.New(e.cents, "USD"),   // minor units: 2800 is $28.00
    FeaturedImage: shopImage(e.handle, e.handle+".jpg"),
    Images:        files.FileRefs{ /* three more */ },
    Available:     true,
    CollectionID:  collectionID,
})
```

Two things about the product seeder are deliberate and worth copying.

**It looks collections up by handle**, not by taking whichever row happened to be
inserted first. A seeder that depends on insertion order produces a different
shop every time you reseed it.

**The first six products are the apparel.** The generated variant seeder attaches
a Colour x Size matrix to the first few products and leaves the rest plain,
because a catalogue where every row has sixteen combinations is a catalogue where
nothing tests the plain case. A mug does not come in a medium.

Reseed:

```bash
grit migrate --fresh
grit seed
```

You should see:

```
Seeded 3 collection(s)
Seeded 12 product(s)
Seeded the Colour option with 4 values
Seeded the Size option with 4 values
Seeded 96 product variants across 6 products
Seeded 3 page(s)
```

Ninety-six variants, from twelve products and two options, and you wrote none of
them.

Check the public API has what the storefront will need:

```bash
KEY=$(grep NEXT_PUBLIC_API_KEY apps/web/.env.local | cut -d= -f2)
curl -s -H "X-API-Key: $KEY" "http://localhost:8080/api/v1/public/products?sort_by=price_amount&sort_order=asc" | head -40
```

Note `sort_by=price_amount`, not `price`. A money column is stored as two
columns, `price_amount` and `price_currency`, and the sortable one is the
amount.

---

## 7. The cart, which Grit does not generate

Everything so far was a flag. This part is not.

`--public` generates **read-only** endpoints, on purpose: a list, a get-by-handle
and a related strip, for an audience that is the internet. A cart is the other
thing a storefront needs, a public endpoint that *writes*, and there is no
generator for one. So you write it, and it is about 550 lines of Go.

### The models

Generate the tables, because that part is still a flag:

```bash
grit generate resource Cart --fields "token:string,currency:string"

grit generate resource CartItem \
  --fields "cart:belongs_to:Cart,product:belongs_to:Product,variant_id:string,quantity:int,unit_price:money,title:string,variant_label:string,image_url:string"

grit migrate
```

That gives you models, migrations, admin screens (a shop owner can look at
abandoned baskets) and staff CRUD. What it does not give you is the shopper's
side.

### The service

`apps/api/internal/services/shop_cart.go`, 377 lines. Copy it from this folder.
Four decisions in it are the ones that matter, and all four are the kind a
generator cannot make for you:

**A cart belongs to a browser, not to a user.** A shopper fills a basket before
deciding whether to have an account. The cart is identified by a token:

```go
func NewToken() (string, error) {
    raw := make([]byte, 32)
    if _, err := rand.Read(raw); err != nil {
        return "", fmt.Errorf("generating a cart token: %w", err)
    }
    return base64.RawURLEncoding.EncodeToString(raw), nil
}
```

32 bytes from `crypto/rand`. The token *is* the authorisation to read and change
that cart, because there is no user to check it against, so a sequential id or a
hash of something known would let anyone read the basket next door.

**The price comes from the database, never from the request.** The client sends a
product id, maybe a variant id, and a quantity. That is all. A price in a request
is a price the customer chose.

**The price is snapshotted onto the line.** A shopper who put something in the
basket at $28 is charged $28 at checkout even if you raise the price while they
browse. The alternative is a total that changes between the basket and the card.

**A token that does not resolve gets a new cart, not an error.** A cart is not an
account. It expires, it gets cleaned up, and somebody returning to a stale one
should be able to start shopping rather than be told something broke.

#### The one that bit

The variant label ("Black / XL") was first built with a three-table join, and the
join column was named `product_variant_id`. The table calls it `variant_id`, so
every add-to-basket with a variant returned a 500. The fix was not to correct the
column name but to delete the SQL: the variant's option values are already
preloaded, and the options are already loaded to resolve the price.

```go
func (s *ShopCartService) variantLabel(variant models.ProductVariant, options []models.Option) string {
    chosen := make(map[string]string, len(variant.OptionValues))
    for _, value := range variant.OptionValues {
        chosen[value.OptionID] = value.Label
    }
    parts := make([]string, 0, len(options))
    for _, option := range options {
        if label, ok := chosen[option.ID]; ok {
            parts = append(parts, label)
        }
    }
    return strings.Join(parts, " / ")
}
```

Ordered by the option's position, because "Black / XL" and "XL / Black" are the
same variant and only one of them reads like a label.

**The lesson generalises:** reaching into a generated schema with raw SQL is the
thing that breaks when the generator changes. The association was right there.

### The handler

`apps/api/internal/handlers/shop_cart.go`, 190 lines. Four endpoints, each
returning the whole cart so that one click updates the badge, the drawer and the
subtotal from one round trip.

Two things in it are about making a *public write surface* safe:

```go
const cartTokenHeader = "X-Cart-Token"

func private(c *gin.Context) {
    c.Header("Cache-Control", "private, no-store")
}
```

The token is a header, not a cookie, because the caller is your own Next.js
server: the browser talks to a server action, the action keeps the cookie and
talks to this API. The token is then out of reach of any script on your pages,
and the endpoint needs no CORS allowance and no CSRF token, because no browser
ever calls it directly.

### Mount the routes, and where

```go
// in apps/api/internal/routes/routes.go, after the publicAPI group
shopAPI := v1.Group("/shop")
shopAPI.Use(middleware.RequireAPIKey(db, svc.Cache))
{
    cart := handlers.NewShopCartHandler(db)
    shopAPI.GET("/cart", cart.Get)
    shopAPI.POST("/cart/items", cart.Add)
    shopAPI.PATCH("/cart/items/:id", cart.SetQuantity)
    shopAPI.DELETE("/cart/items/:id", cart.Remove)
}
```

**Its own group, not under `publicAPI`.** The public group has a response cache
on it, keyed on the URL and nothing else. That is exactly right for a catalogue,
where every caller gets the same answer, and exactly wrong here: every shopper's
`GET /cart` has the same URL, so one stored copy would be served to all of them.

Since v3.367.0 the cache middleware refuses to store a response that says
`Cache-Control: private`, so this would be safe either way. Mount it outside
anyway: the correct reading of that group is "the same answer for everybody", and
a cart is not that.

Check it:

```bash
grit routes | grep shop/cart
```

```
GET     /api/v1/shop/cart              cart.Get           api-key
POST    /api/v1/shop/cart/items        cart.Add           api-key
PATCH   /api/v1/shop/cart/items/:id    cart.SetQuantity   api-key
DELETE  /api/v1/shop/cart/items/:id    cart.Remove        api-key
```

And exercise it. Adding the same variant twice should raise the quantity rather
than make a second line, and another shopper's token must not reach your line:

```bash
B=http://localhost:8080/api/v1
T=$(curl -s -H "X-API-Key: $KEY" "$B/shop/cart" | jq -r .data.token)
P=$(curl -s -H "X-API-Key: $KEY" "$B/public/products/stoneware-mug" | jq -r .data.id)
curl -s -X POST -H "X-API-Key: $KEY" -H "X-Cart-Token: $T" \
  -H "Content-Type: application/json" -d "{\"product_id\":\"$P\",\"quantity\":2}" \
  "$B/shop/cart/items" | jq '.data | {count, subtotal}'
```

---

## 8. The storefront

`apps/web`. Nine files of components, five routes, two data-layer files.

### Replace the scaffold's route group

The scaffold's homepage lives in `app/(marketing)/`. Your shop's homepage is also
`/`, and Next.js refuses to build two route groups that resolve to the same path,
so the marketing group goes:

```bash
cd apps/web
mkdir -p "app/(shop)"
mv "app/(marketing)/blog" "app/(shop)/blog"
rm -rf "app/(marketing)"
```

> Before v3.367.0, the next `grit upgrade` wrote that group back and the project
> stopped building. An upgrade now leaves a deleted page deleted, because a
> deleted page is a routing decision and not damage.

### The data layer

Two files, both hand-written, both typed against the Go side.

**`lib/product-variants-public.ts`** (160 lines) is the client for the variant
endpoint, plus the two functions the picker needs:

```ts
// The variant a selection identifies, or null while the selection is partial.
export function matchVariant(payload, selected): PublicVariant | null

// Whether choosing this value leaves anything in stock, given the other choices.
export function valueIsAvailable(payload, optionId, valueId, selected): boolean
```

`valueIsAvailable` is what greys out "XL" once the shopper picks a colour XL was
never made in, rather than letting them select it and then showing an error.

**`lib/cart.ts`** (163 lines) is the cart, from the storefront's *server*. Every
function in it runs on the Next.js server and none in the browser:

```ts
const COOKIE = "shop_cart";

jar.set(COOKIE, cart.token, {
  httpOnly: true,
  sameSite: "lax",
  secure: process.env.NODE_ENV === "production",
  path: "/",
  maxAge: 60 * 60 * 24 * 30,
});
```

One subtlety: `getCart` does not create a cart.

```ts
export async function getCart(): Promise<Cart> {
  const jar = await cookies();
  if (!jar.get(COOKIE)?.value) return EMPTY;
  return call("GET", "/cart", undefined, false);
}
```

The API answers a read with a fresh token when the request had none, which is
correct for it and wrong to write down here: otherwise a crawler loading the
homepage gets a `Set-Cookie` and a row in `carts`. A cart is created when
something is put in it.

`app/(shop)/cart-actions.ts` (86 lines) wraps the four operations as server
actions, each returning the whole cart. They return errors rather than throwing
them, because an action that throws shows the error boundary, and replacing the
entire page with an apology is the wrong response to "that size just sold out".

### The components

| File | Lines | |
|---|---|---|
| `cart-provider.tsx` | 110 | the cart's client state, which is one thing: the last cart the server sent |
| `cart-drawer.tsx` | 252 | the basket, with Escape, focus return and a focus trap written out |
| `variant-picker.tsx` | 169 | the options, the price and the add button |
| `shop-header.tsx` | 134 | collections nav, basket badge, skip link |
| `shop-footer.tsx` | 121 | footer links, from the `pages` table |
| `gallery.tsx` | 104 | a big picture and thumbnails |
| `sort-links.tsx` | 68 | sort options, as links |
| `product-card.tsx` | 65 | one product in a grid |
| `price.tsx` | 51 | every price in the shop goes through this |
| `search-box.tsx` | 45 | a plain GET form |
| `product-grid.tsx` | 37 | a responsive grid |

Three of these are worth reading closely.

**`cart-provider.tsx` does no local arithmetic.** Adding something does not
increment a counter and then reconcile; it calls the server, and the server's
answer replaces what we had. That costs a round trip per click and buys the thing
shoppers care about: the number next to the basket is the number of things in the
basket, always, including after a reload, in a second tab, and when a variant
sold out between the page loading and the click.

**`variant-picker.tsx` computes no prices.** The price shown is the selected
variant's, from the server, and a range until enough is selected to name one. The
API already resolved every combination's price; a second implementation of that
arithmetic in the browser is a second answer to "what does this cost".

Its unavailable values are `aria-disabled`, not `disabled`:

```tsx
aria-pressed={isSelected}
aria-disabled={!available}
```

A `disabled` button cannot be focused, so a keyboard user cannot discover that
the combination exists and is unavailable. This one is reachable and refuses.

**`sort-links.tsx` uses links, not a `<select>`.** A link changes the URL, so the
sorted view is shareable, bookmarkable and in the back button, and it works
before any JavaScript has run.

### The routes

| | |
|---|---|
| `app/(shop)/layout.tsx` | header, footer, cart provider, drawer |
| `app/(shop)/page.tsx` | homepage: three tiles, collections, a grid |
| `app/(shop)/search/page.tsx` | everything, searchable and sortable |
| `app/(shop)/search/[collection]/page.tsx` | one collection |
| `app/(shop)/product/[handle]/page.tsx` | gallery, picker, details, related |
| `app/(shop)/[page]/page.tsx` | About, Shipping, Terms: one route, three rows |

The layout reads the navigation out of the database on every render, so adding a
collection in the admin adds it to the header and adding a page adds it to the
footer, neither needing a deploy. It also reads the cart once and hands it to the
provider, so the drawer opens with its contents already there rather than showing
a spinner on the click that matters most.

The collection page is the one place the allowlist shows through. The `collection`
relation is deliberately *not* published, because publishing a relation publishes
a whole related record nobody vetted. So the page looks the collection up by its
own handle and filters products by the id it gets back:

```ts
const collection = await getPublicCollection(handle);
if (!collection) notFound();
const products = await getPublicProducts({
  collection_id: collection.id,
  sort_by: sort.sort_by,
  sort_order: sort.sort_order,
});
```

Filtering by an id is not the same as publishing the relation: the id identifies
a row the endpoint was already willing to return.

### Run it

```bash
pnpm install
npx next build     # type-checks everything
grit start
```

Open http://localhost:3000. Click a tee, pick a colour, pick a size, watch the
price move, add it to the basket, reload the page, and the badge still says 1.

> **If the images are broken and the console is full of CSP violations:** you are
> on a version before v3.367.0. `picsum.photos` answers with a 302 to
> `fastly.picsum.photos`, and a Content-Security-Policy is checked against the
> host a redirect lands on, so a policy allowing only the first allowed nothing.
> `grit upgrade` fixes it.

---

## 9. What to do next

The shop in this folder stops at the basket. Checkout is a button that says so.
Three directions from here, in the order a real shop needs them:

**Checkout.** Add an `Order` and `OrderItem` resource, a public endpoint that
turns a cart into an order, and a payment provider. The pattern is the one in
step 7: generate the tables, hand-write the shopper's side. Prices come off the
cart lines, which already hold the snapshot taken when the item was added.

**Accounts.** `grit add web-auth` puts sign-in, registration and a customer area
in the storefront. Then give `Cart` a nullable `user_id` and adopt the
cookie's cart on sign-in, so a basket filled before signing in survives it.

**Stock that moves.** `product_variants.stock` is a column nobody decrements. An
order should take stock, a cancellation should give it back, and both belong in a
transaction with the order's own write.

And two things to try on what you have:

```bash
grit routes          # every endpoint, its handler, and who may call it
grit doctor          # security and configuration mistakes, before they ship
```

---

## What Grit wrote, and what you did

| | Files | |
|---|---|---|
| Generated by Grit | 796 | the API, the admin, the catalogue, the variant matrix, the public endpoints, the tests |
| Written by you | 25 | the cart (567 lines of Go), the storefront (2,104 lines of TS and TSX), the seed data (307) |

2,978 lines out of 124,280. The 25 are the shop: what a cart does
when the same variant is added twice, which option values are unavailable, what a
price looks like, what happens when something sells out between the page loading
and the click.

The other 796 are the part that is the same in every shop, and the reason this
took an hour.
