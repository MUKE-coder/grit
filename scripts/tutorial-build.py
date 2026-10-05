#!/usr/bin/env python3
"""Build examples/commerce/TUTORIAL.md, with every file's real contents in it.

The prose lives here; the code does not. Each `file(path, note)` call reads the
file out of examples/commerce and embeds it whole, under a heading that names
the exact path to create. A tutorial that paraphrases its code is a tour, and a
reader following it reaches the cart and stops, which is what happened.

    python scripts/tutorial-build.py            # write the tutorial
    python scripts/tutorial-build.py --check    # fail if it is out of date

The --check form runs in CI, so a change to the example that is not reflected
in the tutorial fails the build rather than going quietly stale.
"""
import io
import os
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
EXAMPLE = os.path.join(ROOT, 'examples', 'commerce')
OUT = os.path.join(EXAMPLE, 'TUTORIAL.md')

NL = '\n'
LANGS = {'.go': 'go', '.ts': 'ts', '.tsx': 'tsx', '.md': 'markdown',
         '.json': 'json', '.css': 'css'}

parts = []
embedded = []


def md(text):
    """A block of prose, de-indented and stripped."""
    parts.append(text.strip(NL))


def file(path, note='', verb='Create'):
    """Embed a file from the example, whole, labelled with its path."""
    full = os.path.join(EXAMPLE, path.replace('/', os.sep))
    with io.open(full, encoding='utf-8', newline='') as fh:
        body = fh.read().replace('\r\n', NL).rstrip(NL)
    lang = LANGS.get(os.path.splitext(path)[1], '')
    fence = '````' if '```' in body else '```'
    embedded.append(path)

    block = '**' + verb + ' `' + path + '`**'
    if note:
        block += NL + NL + note.strip(NL)
    block += NL + NL + fence + lang + NL + body + NL + fence
    parts.append(block)


# ───────────────────────────────────────────────────────────────────────────
md('''
# Build a shop with Grit

This builds the application in this folder, from an empty directory. Every
command is given in full, and **every file you write is given in full, with the
path to put it at**, so you can follow it start to finish without opening
anything else.

It assumes you are comfortable with Go and with React, and that you have never
used Grit.

**What you are building:** [vercel/commerce](https://github.com/vercel/commerce)
is a storefront with no backend: it is a Next.js app in front of Shopify. This
builds the same storefront and the backend it talks to, which is the part
Shopify was doing. A catalogue, a Colour x Size variant matrix, a server-side
basket, static pages, and an admin panel to run it from.

At the end you will have written 25 files, about 3,000 lines, on top of roughly
124,000 that Grit generated.
''')

md('''
---

## Contents

- [In a hurry? Clone it and run it](#in-a-hurry-clone-it-and-run-it)
- [1. Before you start](#1-before-you-start)
- [2. Scaffold the project](#2-scaffold-the-project)
- [3. Generate the catalogue](#3-generate-the-catalogue)
- [4. Add the variant matrix](#4-add-the-variant-matrix)
- [5. Run it for the first time](#5-run-it-for-the-first-time)
- [6. Write the seed data](#6-write-the-seed-data) *(3 files)*
- [7. The cart, which Grit does not generate](#7-the-cart-which-grit-does-not-generate) *(2 files + routes)*
- [8. The storefront](#8-the-storefront) *(20 files)*
- [9. What to do next](#9-what-to-do-next)
- [Every file you wrote](#every-file-you-wrote)
''')

md('''
---

## In a hurry? Clone it and run it

The finished shop is in this repository. To have it running in about five
minutes, with no Docker, no Postgres and no object storage:

```bash
# 1. Get it. --depth 1 because you want the shop, not the history.
git clone --depth 1 https://github.com/MUKE-coder/grit.git
cd grit/examples/commerce

# 2. Install the CLI, if you have not already.
go install github.com/MUKE-coder/grit/v3/cmd/grit@latest

# 3. Dependencies, and an .env with fresh secrets.
pnpm install
grit env
```

Then open `.env` and change one line and add three:

```ini
DB_PROVIDER=sqlite
```

```ini
REDIS_URL=
MAIL_MAILER=log
STORAGE_DRIVER=local
```

That is SQLite instead of Postgres, no Redis, mail printed to the log, and
uploads kept in a folder. Then:

```bash
grit migrate     # creates the tables
grit seed        # 12 products, 96 variants, and your admin password, printed once
grit start       # API, storefront and admin together
```

| | |
|---|---|
| Storefront | http://localhost:3000 |
| Admin | http://localhost:3001 |
| API docs | http://localhost:8080/docs |
| Database browser | http://localhost:8080/studio |

Sign in to the admin with `admin@example.com` and the password `grit seed`
printed. Read `grit seed`'s output before you lose it: the password and both API
keys are shown once.

### Only want the example, not the whole repository?

```bash
git clone --depth 1 --filter=blob:none --sparse https://github.com/MUKE-coder/grit.git
cd grit
git sparse-checkout set examples/commerce
cd examples/commerce
```

Everything else below builds the same thing from nothing, which is the part
worth doing if you want to understand it.
''')

md('''
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

This tutorial was written against **v3.371.0**. Earlier versions have bugs this
build found, and you would hit all of them: seeded money prices, the public
API's price field, the variant lookup by handle, placeholder images blocked by
the project's own security policy, and uploads refused on a machine without
Docker. Use v3.371.0 or later.

You do **not** need Docker, Postgres, Redis or any object storage.
''')

md('''
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

```bash
grit env
```

Open `.env`. Change this line:

```ini
DB_PROVIDER=sqlite
```

And add these three at the end:

```ini
REDIS_URL=
MAIL_MAILER=log
STORAGE_DRIVER=local
```

`REDIS_URL=` set to nothing means "run without Redis": the cache, background
jobs and cron are skipped rather than failing to connect. `MAIL_MAILER=log`
prints password-reset mails to the log. The third line is the one worth
understanding, because it is where a shop meets its first wall.

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

> **You may not need the line at all.** Outside production, a `minio` that is
> not answering falls back to the local disk by itself and says so in the log.
> Setting it makes it your decision rather than a fallback. Production never
> falls back: a server that has lost its bucket says so rather than quietly
> writing to a disk nobody is backing up.

#### One thing not to change

Leave `APP_PORT` alone. `APP_URL` is written into `.env` as
`http://localhost:8080`, and the two have to agree: `APP_URL` is what every
generated link is built from, uploaded files included, so a server listening on
one port while `APP_URL` names another serves images from an address where
nothing is answering. If you do need a different port, change both.
''')

md('''
---

## 3. Generate the catalogue

A shop is three things: products, the collections they are grouped into, and the
static pages in the footer. Run these three commands from the project root.

### Collections

```bash
grit generate resource Collection \\
  --fields "title:string,handle:slug,description:text,image:file:image" \\
  --public --seed
```

### Products

```bash
grit generate resource Product \\
  --fields "title:string,handle:slug,description:richtext,price:money,featured_image:file:image,images:files,available:bool,collection:belongs_to:Collection" \\
  --public --seed
```

### Pages

```bash
grit generate resource Page \\
  --fields "title:string,handle:slug,body:richtext" \\
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
notice.

**`--public`** adds a second, narrower API: a read-only list and get-by-handle
under `/api/v1/public/products`, guarded by a publishable API key, built from an
allowlist rather than from the model. A column added next month is private until
somebody adds it to that allowlist, which is the right default when the audience
is the internet. The staff CRUD the admin panel uses is separate and behind
authentication.

**`--seed`** writes a seeder with one sample row in it, which you replace in
step 6.
''')

md('''
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
''')

md('''
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

**Read the output of `grit seed`.** It prints the admin email and password once,
and both API keys once.

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

Sign in to the admin with the credentials `grit seed` printed. You have a
working product catalogue with a variant matrix editor, and a storefront that is
still the scaffold's landing page.

> **If every page says `INVALID_API_KEY`:** you ran `grit migrate --fresh` at
> some point, which drops the keys table. The seeder minted a new key and
> deliberately did not overwrite your `.env.local`, because overwriting
> somebody's environment file is how you lose an afternoon. Since v3.367.0 it
> tells you, and prints the key to paste in.
''')

# ─── 6. seed data ──────────────────────────────────────────────────────────
md('''
---

## 6. Write the seed data

`--seed` wrote three seeders with one row each, called "Sample Title". They exist
to be edited: the generator knows the shape of your model and nothing about your
shop. Replace all three with what follows.
''')

file('apps/api/internal/database/collections_seeder.go', verb='Replace', note='''
Three collections. One collection is indistinguishable from no collections on a
storefront, because nothing has to be filtered; three is enough for the
navigation, the collection pages and the related-items strip to all have
something to show.

It also adds a `shopImage` helper, which the products seeder below uses.
''')

file('apps/api/internal/database/products_seeder.go', verb='Replace', note='''
Twelve products. Two things here are deliberate and worth copying into your own:

**It looks collections up by handle**, not by taking whichever row happened to
be inserted first. A seeder that depends on insertion order produces a different
shop every time you reseed it.

**The first six products are the apparel.** The generated variant seeder
attaches a Colour x Size matrix to the first few products and leaves the rest
plain, because a catalogue where every row has sixteen combinations is a
catalogue where nothing tests the plain case. A mug does not come in a medium.
''')

file('apps/api/internal/database/pages_seeder.go', verb='Replace', note='''
Three pages. The footer links to them by handle and the storefront renders them
from one dynamic route, so adding a fourth needs no code at all.
''')

md('''
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
curl -s -H "X-API-Key: $KEY" \\
  "http://localhost:8080/api/v1/public/products?sort_by=price_amount&sort_order=asc" | head -40
```

Note `sort_by=price_amount`, not `price`. A money column is stored as two
columns, `price_amount` and `price_currency`, and the sortable one is the
amount.
''')

# ─── 7. the cart ───────────────────────────────────────────────────────────
md('''
---

## 7. The cart, which Grit does not generate

Everything so far was a flag. This part is not.

`--public` generates **read-only** endpoints, on purpose: a list, a get-by-handle
and a related strip, for an audience that is the internet. A cart is the other
thing a storefront needs, a public endpoint that *writes*, and there is no
generator for one. So you write it: two files, and four lines in a third.

### 7a. Generate the tables

The models are still a flag. Run these from the project root:

```bash
grit generate resource Cart --fields "token:string,currency:string"

grit generate resource CartItem \\
  --fields "cart:belongs_to:Cart,product:belongs_to:Product,variant_id:string,quantity:int,unit_price:money,title:string,variant_label:string,image_url:string"

grit migrate
```

That gives you models, migrations, admin screens (a shop owner can look at
abandoned baskets) and staff CRUD. What it does not give you is the shopper's
side, which is the next two files.

### 7b. The service
''')

file('apps/api/internal/services/shop_cart.go', note='''
This is a **new file**. The resource generator above wrote
`apps/api/internal/services/cart.go` and `cart_item.go` for the admin's CRUD;
this one is the shopper's side of the same tables, and they are deliberately
separate services rather than one with two audiences. That is also why the type
is `ShopCartService` and not `CartService`, which is taken.

Four decisions in it are the ones that matter, and all four are the kind a
generator cannot make for you:

1. **A cart belongs to a browser, not to a user.** A shopper fills a basket
   before deciding whether to have an account. The identifier is 32 bytes from
   `crypto/rand`, and that token *is* the authorisation, because there is no
   user to check it against: a sequential id or a hash of something known would
   let anyone read the basket next door.
2. **The price comes from the database, never from the request.** The client
   sends a product id, maybe a variant id, and a quantity. That is all. A price
   in a request is a price the customer chose.
3. **The price is snapshotted onto the line.** A shopper who put something in
   the basket at $28 is charged $28 at checkout even if you raise the price
   while they browse. The alternative is a total that changes between the basket
   and the card.
4. **A token that does not resolve gets a new cart, not an error.** A cart is
   not an account. It expires, it gets cleaned up, and somebody returning to a
   stale one should be able to start shopping rather than be told something
   broke.

One thing to notice in `variantLabel` at the bottom: it reads the preloaded
association rather than running a join. The first version of it ran a
three-table query and named the join column `product_variant_id`, which the
table does not have (it is `variant_id`), so every add-to-basket with a variant
returned a 500. Reaching into a generated schema with raw SQL is exactly what
breaks when the generator changes, and the association was right there.
''')

md('### 7c. The handler')

file('apps/api/internal/handlers/shop_cart.go', note='''
Also a **new file**, beside the generated `cart.go` and `cart_item.go`.

Four endpoints, each returning the whole cart so that one click updates the
badge, the drawer and the subtotal from one round trip.

Two things in it are about making a *public write surface* safe. The cart token
travels in a header, not a cookie, because the caller is your own Next.js
server: the browser talks to a server action, the action keeps the cookie and
talks to this API. The token is then out of reach of any script on your pages,
and the endpoint needs no CORS allowance and no CSRF token, because no browser
ever calls it directly. And every response says `Cache-Control: private`.
''')

md('''
### 7d. Mount the routes

Open **`apps/api/internal/routes/routes.go`** and find the public group, which
looks like this:

```go
	{
		productPublicVariants := handlers.NewProductVariantHandler(db)
		publicAPI.GET("/products/:key/variants", productPublicVariants.ListPublic)
		// grit:routes:public
	}
```

Add this **directly after that closing brace**:

```go
	// The storefront's cart: a public surface that writes.
	//
	// Its own group, NOT under publicAPI, and the reason is the response cache
	// mounted on that group. The cache key is the URL and nothing else, which
	// is exactly right for a catalogue and exactly wrong here: every shopper's
	// GET /cart has the same URL, so one stored copy would be served to all of
	// them.
	//
	// Same API key, so a cart gets the same identification, rate-limit bucket
	// and per-origin narrowing as the rest of the public surface.
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

Check it compiles and that the routes are there:

```bash
cd apps/api && go build ./... && cd ../..
grit routes | grep shop/cart
```

```
GET     /api/v1/shop/cart              cart.Get           api-key
POST    /api/v1/shop/cart/items        cart.Add           api-key
PATCH   /api/v1/shop/cart/items/:id    cart.SetQuantity   api-key
DELETE  /api/v1/shop/cart/items/:id    cart.Remove        api-key
```

### 7e. Exercise it

Restart the API, then add two mugs to a basket. Adding the same variant twice
should raise the quantity rather than make a second line, and another shopper's
token must not reach your line:

```bash
B=http://localhost:8080/api/v1
KEY=$(grep NEXT_PUBLIC_API_KEY apps/web/.env.local | cut -d= -f2)

# A cart token, and a product id.
T=$(curl -s -H "X-API-Key: $KEY" "$B/shop/cart" | jq -r .data.token)
P=$(curl -s -H "X-API-Key: $KEY" "$B/public/products/stoneware-mug" | jq -r .data.id)

curl -s -X POST -H "X-API-Key: $KEY" -H "X-Cart-Token: $T" \\
  -H "Content-Type: application/json" \\
  -d "{\\"product_id\\":\\"$P\\",\\"quantity\\":2}" \\
  "$B/shop/cart/items" | jq '.data | {count, subtotal}'
```

```json
{ "count": 2, "subtotal": { "amount": 4800, "currency": "USD" } }
```
''')

# ─── 8. the storefront ─────────────────────────────────────────────────────
md('''
---

## 8. The storefront

Twenty files in `apps/web`: two for the data layer, eleven components, six
routes, and one server-actions file. Every one is given in full below.

### 8a. Replace the scaffold's route group

The scaffold's homepage lives in `app/(marketing)/`. Your shop's homepage is also
`/`, and Next.js refuses to build two route groups that resolve to the same
path, so the marketing group goes:

```bash
cd apps/web
mkdir -p "app/(shop)"
mv "app/(marketing)/blog" "app/(shop)/blog"
rm -rf "app/(marketing)"
```

Then make the directories the rest of this section fills:

```bash
mkdir -p "app/(shop)/search/[collection]" \\
         "app/(shop)/product/[handle]" \\
         "app/(shop)/[page]" \\
         components/shop
cd ../..
```

> Before v3.367.0, the next `grit upgrade` wrote that group back and the project
> stopped building. An upgrade now leaves a deleted page deleted, because a
> deleted page is a routing decision and not damage.

### 8b. The data layer

Two files, both typed against the Go side.
''')

file('apps/web/lib/product-variants-public.ts', note='''
The client for the variant endpoint, plus the two functions the picker needs.

`grit add variants` generates the Go endpoint but no TypeScript client for it,
so this is the one part of the storefront's data layer that is not generated.

`matchVariant` returns the variant a selection identifies, or null while the
selection is partial: a tee in black is not a thing you can buy until it is also
a size. `valueIsAvailable` is what greys out "XL" once the shopper picks a
colour XL was never made in, rather than letting them select it and then showing
an error.
''')

file('apps/web/lib/cart.ts', note='''
The cart, from the storefront's **server**. Every function here runs on the
Next.js server and none in the browser: the browser talks to a server action,
the action keeps the cart token in an HttpOnly cookie, and the action talks to
the Go API.

One subtlety worth reading: `getCart` does not create a cart. The API answers a
read with a fresh token when the request had none, which is correct for it and
wrong to write down here, because otherwise a crawler loading the homepage gets
a `Set-Cookie` and a row in `carts`. A cart is created when something is put in
it.
''')

file('apps/web/app/(shop)/cart-actions.ts', note='''
The four cart operations as server actions, each returning the whole cart.

They return errors rather than throwing them. An action that throws shows the
error boundary, and replacing the entire page with an apology is the wrong
response to "that size just sold out".
''')

md('''
### 8c. The components

Eleven files, all in `apps/web/components/shop/`. Three of them are worth
reading closely and the notes say why.
''')

file('apps/web/components/shop/price.tsx', note='''
Every price in the shop goes through this. `formatMoney` in the shared package
knows which currencies have two decimal places and which have none: a hardcoded
`amount / 100` shows a 50,000 shilling price as 500, and nobody notices until
the first Ugandan customer complains.
''')

file('apps/web/components/shop/product-card.tsx', note='''
One product in a grid. The whole card is one link, so there is one tab stop per
product rather than two competing ones, and the accessible name is the title.
''')

file('apps/web/components/shop/product-grid.tsx', note='''
A `<ul>` with one `<li>` per product, because a screen reader then announces
"list, 12 items" before reading any of them.
''')

file('apps/web/components/shop/gallery.tsx', note='''
A big picture and a row of thumbnails, with arrows for anyone who would rather
not aim at a 72px target.
''')

file('apps/web/components/shop/cart-provider.tsx', note='''
The cart's client state, which is one thing: **the last cart the server sent.**

There is no local arithmetic here. Adding something does not increment a counter
and then reconcile; it calls the server, and the server's answer replaces what
we had. That costs a round trip per click and buys the thing shoppers care
about: the number next to the basket is the number of things in the basket,
always, including after a reload, in a second tab, and when a variant sold out
between the page loading and the click.
''')

file('apps/web/components/shop/cart-drawer.tsx', note='''
The basket. Hand-written rather than reached for from a primitive library, which
is the convention in this project, so the three things a dialog has to do are
here explicitly: Escape closes it, focus moves into it and back to the trigger,
and Tab cannot leave it while it is open. Those are the parts a hand-rolled
dialog usually skips, and skipping them fails exactly the people who are not in
the room when it is demonstrated.
''')

file('apps/web/components/shop/variant-picker.tsx', note='''
The options, the price and the add button, which have to be one component
because they all depend on the same selection.

**It computes no prices.** The price shown is the selected variant's, from the
server, and a range until enough is selected to name one. The API already
resolved every combination's price; a second implementation of that arithmetic
in the browser is a second answer to "what does this cost".

Unavailable values are `aria-disabled`, not `disabled`. A `disabled` button
cannot be focused, so a keyboard user cannot discover that the combination
exists and is unavailable. This one is reachable and refuses.
''')

file('apps/web/components/shop/sort-links.tsx', note='''
The sort options, as links rather than a `<select>`. A link changes the URL, so
the sorted view is shareable, bookmarkable and in the back button, and it works
before any JavaScript has run.

Note `sort_by: "price_amount"`. A money field is two columns, and the sortable
one is the amount.
''')

file('apps/web/components/shop/search-box.tsx', note='''
A plain GET form, so it needs no JavaScript and no client component.
''')

file('apps/web/components/shop/shop-header.tsx', note='''
Collections nav, basket badge and a skip link. The badge count is in the
button's accessible name, so it is announced rather than only drawn.
''')

file('apps/web/components/shop/shop-footer.tsx', note='''
The footer, whose "Information" links come out of the `pages` table. Adding a
fourth page needs no deploy and no code.
''')

md('''
### 8d. The routes

Six files under `apps/web/app/(shop)/`.
''')

file('apps/web/app/(shop)/layout.tsx', note='''
The shop's chrome, and the one place the cart is read.

The navigation and the footer links come out of the database on every render, so
adding a collection in the admin adds it to the header and adding a page adds it
to the footer, neither needing a deploy. The cart is read once here and handed to
the provider, so the drawer opens with its contents already there rather than
showing a spinner on the click that matters most.
''')

file('apps/web/app/(shop)/page.tsx', note='''
The homepage: three tiles, the collections, then a grid. That is the shape
vercel/commerce uses and it uses it for a reason: a shop's homepage has one job,
which is to get somebody onto a product page, and a hero with no products in it
does not do that job.
''')

file('apps/web/app/(shop)/search/page.tsx', note='''
Everything, searchable and sortable. Both the search term and the sort are in
the URL and both are handled by the public API, so the page works with 12
products and with 12,000.
''')

file('apps/web/app/(shop)/search/[collection]/page.tsx', note='''
One collection's products, and the one place the allowlist shows through.

The `collection` relation is deliberately *not* published, because publishing a
relation publishes a whole related record nobody vetted. So this page looks the
collection up by its own handle and filters products by the id it gets back.
Filtering by an id is not the same as publishing the relation: the id identifies
a row the endpoint was already willing to return.
''')

file('apps/web/app/(shop)/product/[handle]/page.tsx', note='''
The product page: gallery, picker, details, related strip, and structured data
so the price shows in search results.

Three requests, in parallel. The variant payload is one request rather than one
per swatch, which matters because this is the page with the most traffic and the
most interaction in any shop.
''')

file('apps/web/app/(shop)/[page]/page.tsx', note='''
About, Shipping and Terms: three rows in `pages`, one route, no code each.

This is the last route Next.js tries, because a static segment beats a dynamic
one, so `/search` is the search page and is never looked up here.
''')

md('''
### 8e. Run it

```bash
pnpm install
cd apps/web && npx next build && cd ../..   # type-checks everything
grit start
```

Open http://localhost:3000. Click a tee, pick a colour, pick a size, watch the
price move, add it to the basket, reload the page, and the badge still says 1.

> **If the images are broken and the console is full of CSP violations:** you
> are on a version before v3.367.0. `picsum.photos` answers with a 302 to
> `fastly.picsum.photos`, and a Content-Security-Policy is checked against the
> host a redirect lands on, so a policy allowing only the first allowed nothing.
> `grit upgrade` fixes it.
''')

md('''
---

## 9. What to do next

The shop stops at the basket. Checkout is a button that says so. Three
directions from here, in the order a real shop needs them:

**Checkout.** Add an `Order` and `OrderItem` resource, a public endpoint that
turns a cart into an order, and a payment provider. The pattern is the one in
step 7: generate the tables, hand-write the shopper's side. Prices come off the
cart lines, which already hold the snapshot taken when the item was added.

**Accounts.** `grit add web-auth` puts sign-in, registration and a customer area
in the storefront. Then give `Cart` a nullable `user_id` and adopt the cookie's
cart on sign-in, so a basket filled before signing in survives it.

**Stock that moves.** `product_variants.stock` is a column nobody decrements. An
order should take stock, a cancellation should give it back, and both belong in
a transaction with the order's own write.

And two things to try on what you have:

```bash
grit routes          # every endpoint, its handler, and who may call it
grit doctor          # security and configuration mistakes, before they ship
```
''')

# ─── appendix ──────────────────────────────────────────────────────────────
rows = NL.join('| `' + p + '` |' for p in embedded)
md('''
---

## Every file you wrote

Twenty-five files, in the order this tutorial creates them. Everything else in
the project came from `grit new`, `grit generate resource` and
`grit add variants`.

| File |
|---|
''' + rows + '''

Plus four lines mounted in `apps/api/internal/routes/routes.go`.

| | Files | |
|---|---|---|
| Generated by Grit | 796 | the API, the admin, the catalogue, the variant matrix, the public endpoints, the tests |
| Written by you | 25 | the cart, the storefront, the seed data |

2,978 lines out of 124,280. The 25 are the shop: what a cart does when the same
variant is added twice, which option values to grey out, what a price looks
like, what happens when something sells out between the page loading and the
click.

The other 796 are the part that is the same in every shop, and the reason this
took an afternoon.
''')

# ───────────────────────────────────────────────────────────────────────────
body = (NL + NL).join(parts) + NL

if '--check' in sys.argv:
    with io.open(OUT, encoding='utf-8', newline='') as fh:
        have = fh.read().replace('\r\n', NL)
    if have != body:
        print('TUTORIAL.md is out of date. Run: python scripts/tutorial-build.py')
        sys.exit(1)
    print('TUTORIAL.md matches the example (' + str(len(embedded)) + ' files embedded)')
else:
    with io.open(OUT, 'w', encoding='utf-8', newline='') as fh:
        fh.write(body)
    print('wrote ' + OUT)
    print(str(len(embedded)) + ' files embedded, ' + str(body.count(NL)) + ' lines')
